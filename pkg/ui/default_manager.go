// Package ui provides default value management for project configuration.
//
// This file implements the DefaultManager which handles sensible default values
// for optional metadata fields and provides comprehensive validation for all
// project configuration fields.
package ui

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// DefaultManager manages default values for project configuration
type DefaultManager struct {
	logger         interfaces.Logger
	configPath     string
	systemDefaults *SystemDefaults
	userDefaults   *UserDefaults
}

// SystemDefaults contains system-level default values
type SystemDefaults struct {
	License      string
	Organization string
	OutputPath   string
}

// UserDefaults contains user-specific default values
type UserDefaults struct {
	Author       string
	Email        string
	Organization string
	License      string
	Repository   string
}

// DefaultSource represents where a default value comes from
type DefaultSource string

const (
	DefaultSourceSystem      DefaultSource = "system"
	DefaultSourceUser        DefaultSource = "user"
	DefaultSourceEnvironment DefaultSource = "environment"
	DefaultSourceGit         DefaultSource = "git"
	DefaultSourceInferred    DefaultSource = "inferred"
)

// DefaultValue represents a default value with its source
type DefaultValue struct {
	Value  string        `json:"value"`
	Source DefaultSource `json:"source"`
	Reason string        `json:"reason"`
}

// NewDefaultManager creates a new default manager
func NewDefaultManager(logger interfaces.Logger) *DefaultManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "project-generator", "defaults.json")

	dm := &DefaultManager{
		logger:     logger,
		configPath: configPath,
		systemDefaults: &SystemDefaults{
			License:    "MIT",
			OutputPath: "output/generated",
		},
		userDefaults: &UserDefaults{},
	}

	// Load defaults from various sources
	dm.loadSystemDefaults()
	dm.loadUserDefaults()
	dm.loadEnvironmentDefaults()
	dm.loadGitDefaults()

	return dm
}

// GetDefaultsForProject returns default values for a new project
func (dm *DefaultManager) GetDefaultsForProject(projectName string) *ProjectConfigDefaults {
	defaults := &ProjectConfigDefaults{
		License:      dm.getDefaultLicense().Value,
		Author:       dm.getDefaultAuthor().Value,
		Email:        dm.getDefaultEmail().Value,
		Organization: dm.getDefaultOrganization().Value,
	}

	// Infer repository URL if possible
	if repoDefault := dm.inferRepositoryURL(projectName); repoDefault.Value != "" {
		// We don't set this as default since it's often project-specific
		if dm.logger != nil {
			dm.logger.DebugWithFields("Inferred repository URL", map[string]interface{}{
				"project": projectName,
				"url":     repoDefault.Value,
				"source":  repoDefault.Source,
			})
		}
	}

	return defaults
}

// getDefaultLicense returns the default license
func (dm *DefaultManager) getDefaultLicense() DefaultValue {
	// Priority: User preference > Environment > System default
	if dm.userDefaults.License != "" {
		return DefaultValue{
			Value:  dm.userDefaults.License,
			Source: DefaultSourceUser,
			Reason: "User-configured default license",
		}
	}

	if envLicense := os.Getenv("PROJECT_DEFAULT_LICENSE"); envLicense != "" {
		return DefaultValue{
			Value:  envLicense,
			Source: DefaultSourceEnvironment,
			Reason: "Environment variable PROJECT_DEFAULT_LICENSE",
		}
	}

	return DefaultValue{
		Value:  dm.systemDefaults.License,
		Source: DefaultSourceSystem,
		Reason: "System default license",
	}
}

// getDefaultAuthor returns the default author
func (dm *DefaultManager) getDefaultAuthor() DefaultValue {
	// Priority: User preference > Git config > Environment > System user
	if dm.userDefaults.Author != "" {
		return DefaultValue{
			Value:  dm.userDefaults.Author,
			Source: DefaultSourceUser,
			Reason: "User-configured default author",
		}
	}

	if gitAuthor := dm.getGitConfig("user.name"); gitAuthor != "" {
		return DefaultValue{
			Value:  gitAuthor,
			Source: DefaultSourceGit,
			Reason: "Git configuration user.name",
		}
	}

	if envAuthor := os.Getenv("PROJECT_DEFAULT_AUTHOR"); envAuthor != "" {
		return DefaultValue{
			Value:  envAuthor,
			Source: DefaultSourceEnvironment,
			Reason: "Environment variable PROJECT_DEFAULT_AUTHOR",
		}
	}

	// Fallback to system user
	if currentUser, err := user.Current(); err == nil && currentUser.Name != "" {
		return DefaultValue{
			Value:  currentUser.Name,
			Source: DefaultSourceSystem,
			Reason: "System user name",
		}
	}

	return DefaultValue{
		Value:  "",
		Source: DefaultSourceSystem,
		Reason: "No default author available",
	}
}

// getDefaultEmail returns the default email
func (dm *DefaultManager) getDefaultEmail() DefaultValue {
	// Priority: User preference > Git config > Environment
	if dm.userDefaults.Email != "" {
		return DefaultValue{
			Value:  dm.userDefaults.Email,
			Source: DefaultSourceUser,
			Reason: "User-configured default email",
		}
	}

	if gitEmail := dm.getGitConfig("user.email"); gitEmail != "" {
		return DefaultValue{
			Value:  gitEmail,
			Source: DefaultSourceGit,
			Reason: "Git configuration user.email",
		}
	}

	if envEmail := os.Getenv("PROJECT_DEFAULT_EMAIL"); envEmail != "" {
		return DefaultValue{
			Value:  envEmail,
			Source: DefaultSourceEnvironment,
			Reason: "Environment variable PROJECT_DEFAULT_EMAIL",
		}
	}

	return DefaultValue{
		Value:  "",
		Source: DefaultSourceSystem,
		Reason: "No default email available",
	}
}

// getDefaultOrganization returns the default organization
func (dm *DefaultManager) getDefaultOrganization() DefaultValue {
	// Priority: User preference > Environment > System default
	if dm.userDefaults.Organization != "" {
		return DefaultValue{
			Value:  dm.userDefaults.Organization,
			Source: DefaultSourceUser,
			Reason: "User-configured default organization",
		}
	}

	if envOrg := os.Getenv("PROJECT_DEFAULT_ORGANIZATION"); envOrg != "" {
		return DefaultValue{
			Value:  envOrg,
			Source: DefaultSourceEnvironment,
			Reason: "Environment variable PROJECT_DEFAULT_ORGANIZATION",
		}
	}

	if dm.systemDefaults.Organization != "" {
		return DefaultValue{
			Value:  dm.systemDefaults.Organization,
			Source: DefaultSourceSystem,
			Reason: "System default organization",
		}
	}

	return DefaultValue{
		Value:  "",
		Source: DefaultSourceSystem,
		Reason: "No default organization available",
	}
}

// inferRepositoryURL attempts to infer a repository URL based on project name and user info
func (dm *DefaultManager) inferRepositoryURL(projectName string) DefaultValue {
	// Try to infer from Git remote if in a Git repository
	if gitRemote := dm.getGitConfig("remote.origin.url"); gitRemote != "" {
		// Extract username/organization from existing remote
		if username := dm.extractUsernameFromGitURL(gitRemote); username != "" {
			repoName := strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))
			inferredURL := fmt.Sprintf("https://github.com/%s/%s", username, repoName)
			return DefaultValue{
				Value:  inferredURL,
				Source: DefaultSourceInferred,
				Reason: "Inferred from Git remote origin",
			}
		}
	}

	// Try to infer from environment variables
	if githubUser := os.Getenv("GITHUB_USERNAME"); githubUser != "" {
		repoName := strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))
		inferredURL := fmt.Sprintf("https://github.com/%s/%s", githubUser, repoName)
		return DefaultValue{
			Value:  inferredURL,
			Source: DefaultSourceEnvironment,
			Reason: "Inferred from GITHUB_USERNAME environment variable",
		}
	}

	return DefaultValue{
		Value:  "",
		Source: DefaultSourceSystem,
		Reason: "Could not infer repository URL",
	}
}

// loadSystemDefaults loads system-level defaults
func (dm *DefaultManager) loadSystemDefaults() {
	// System defaults are hardcoded for now
	// In a real implementation, these might come from a system configuration file
	dm.systemDefaults = &SystemDefaults{
		License:    "MIT",
		OutputPath: "output/generated",
	}
}

// loadUserDefaults loads user-specific defaults from configuration file
func (dm *DefaultManager) loadUserDefaults() {
	// For now, we'll just initialize empty user defaults
	// In a real implementation, this would load from a JSON/YAML file
	dm.userDefaults = &UserDefaults{}

	// Try to load from config file if it exists
	if _, err := os.Stat(dm.configPath); err == nil {
		if dm.logger != nil {
			dm.logger.DebugWithFields("Loading user defaults", map[string]interface{}{
				"config_path": dm.configPath,
			})
		}
		// TODO: Implement actual file loading
	}
}

// loadEnvironmentDefaults loads defaults from environment variables
func (dm *DefaultManager) loadEnvironmentDefaults() {
	// Environment variables are checked in the getter methods
	// This method could pre-validate them
	envVars := []string{
		"PROJECT_DEFAULT_LICENSE",
		"PROJECT_DEFAULT_AUTHOR",
		"PROJECT_DEFAULT_EMAIL",
		"PROJECT_DEFAULT_ORGANIZATION",
		"GITHUB_USERNAME",
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			if dm.logger != nil {
				dm.logger.DebugWithFields("Found environment default", map[string]interface{}{
					"variable": envVar,
					"value":    value,
				})
			}
		}
	}
}

// loadGitDefaults loads defaults from Git configuration
func (dm *DefaultManager) loadGitDefaults() {
	gitConfigs := []string{"user.name", "user.email", "remote.origin.url"}

	for _, config := range gitConfigs {
		if value := dm.getGitConfig(config); value != "" {
			if dm.logger != nil {
				dm.logger.DebugWithFields("Found Git default", map[string]interface{}{
					"config": config,
					"value":  value,
				})
			}
		}
	}
}

// getGitConfig retrieves a Git configuration value
func (dm *DefaultManager) getGitConfig(key string) string {
	// This is a simplified implementation
	// In a real implementation, you would use git commands or a Git library
	// For now, we'll return empty strings
	return ""
}

// extractUsernameFromGitURL extracts username from a Git URL
func (dm *DefaultManager) extractUsernameFromGitURL(url string) string {
	// Handle GitHub URLs
	if strings.Contains(url, "github.com") {
		// Handle both SSH and HTTPS URLs
		if strings.HasPrefix(url, "git@github.com:") {
			// SSH format: git@github.com:username/repo.git
			parts := strings.Split(strings.TrimPrefix(url, "git@github.com:"), "/")
			if len(parts) > 0 {
				return parts[0]
			}
		} else if strings.Contains(url, "github.com/") {
			// HTTPS format: https://github.com/username/repo.git
			parts := strings.Split(url, "/")
			for i, part := range parts {
				if part == "github.com" && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}

	return ""
}

// SaveUserDefaults saves user defaults to configuration file
func (dm *DefaultManager) SaveUserDefaults(defaults *UserDefaults) error {
	dm.userDefaults = defaults

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(dm.configPath)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// TODO: Implement actual file saving (JSON/YAML)
	if dm.logger != nil {
		dm.logger.InfoWithFields("Saved user defaults", map[string]interface{}{
			"config_path": dm.configPath,
			"author":      defaults.Author,
			"email":       defaults.Email,
			"license":     defaults.License,
		})
	}

	return nil
}

// GetUserDefaults returns the current user defaults
func (dm *DefaultManager) GetUserDefaults() *UserDefaults {
	return dm.userDefaults
}

// ValidateDefaults validates all default values
func (dm *DefaultManager) ValidateDefaults() error {
	validator := &ProjectConfigValidator{
		projectNameRegex: utils.ProjectNamePattern,
		emailRegex:       utils.EmailPattern,
		urlRegex:         utils.URLPattern,
	}

	// Validate author
	if dm.userDefaults.Author != "" {
		if err := validator.ValidateAuthor(dm.userDefaults.Author); err != nil {
			return fmt.Errorf("invalid default author: %w", err)
		}
	}

	// Validate email
	if dm.userDefaults.Email != "" {
		if err := validator.ValidateEmail(dm.userDefaults.Email); err != nil {
			return fmt.Errorf("invalid default email: %w", err)
		}
	}

	// Validate license
	if dm.userDefaults.License != "" {
		if err := validator.ValidateLicense(dm.userDefaults.License); err != nil {
			return fmt.Errorf("invalid default license: %w", err)
		}
	}

	// Validate repository
	if dm.userDefaults.Repository != "" {
		if err := validator.ValidateRepository(dm.userDefaults.Repository); err != nil {
			return fmt.Errorf("invalid default repository: %w", err)
		}
	}

	return nil
}

// GetDefaultSources returns information about where defaults come from
func (dm *DefaultManager) GetDefaultSources() map[string]DefaultValue {
	return map[string]DefaultValue{
		"license":      dm.getDefaultLicense(),
		"author":       dm.getDefaultAuthor(),
		"email":        dm.getDefaultEmail(),
		"organization": dm.getDefaultOrganization(),
	}
}
