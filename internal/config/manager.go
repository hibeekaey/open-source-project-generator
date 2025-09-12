package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	yaml "gopkg.in/yaml.v3"
)

// Manager implements the ConfigManager interface
type Manager struct {
	validator    *models.ConfigValidator
	cacheDir     string
	defaultsPath string
}

// NewManager creates a new configuration manager
func NewManager(cacheDir, defaultsPath string) interfaces.ConfigManager {
	return &Manager{
		validator:    models.NewConfigValidator(),
		cacheDir:     cacheDir,
		defaultsPath: defaultsPath,
	}
}

// LoadDefaults loads default configuration values
func (m *Manager) LoadDefaults() (*models.ProjectConfig, error) {
	// Try to load from defaults file first
	if m.defaultsPath != "" && fileExists(m.defaultsPath) {
		config, err := m.LoadConfig(m.defaultsPath)
		if err == nil {
			return config, nil
		}
		// If loading fails, fall back to hardcoded defaults
	}

	// Return hardcoded defaults
	return &models.ProjectConfig{
		License: "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   false,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: false,
				IOS:     false,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: false,
				Terraform:  false,
			},
		},
		Versions: &models.VersionConfig{
			Node:      "20.0.0",
			Go:        "1.22.0",
			Kotlin:    "2.0.0",
			Swift:     "5.9.0",
			NextJS:    "15.5.3",
			React:     "18.0.0",
			Packages:  make(map[string]string),
			UpdatedAt: time.Now(),
		},
		CustomVars:       make(map[string]string),
		OutputPath:       "./generated-project",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}, nil
}

// ValidateConfig validates the provided project configuration
func (m *Manager) ValidateConfig(config *models.ProjectConfig) error {
	result := m.validator.ValidateProjectConfig(config)
	if !result.Valid {
		var errorMessages []string
		for _, err := range result.Errors {
			errorMessages = append(errorMessages, err.Message)
		}
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errorMessages, "; "))
	}
	return nil
}

// GetLatestVersions fetches the latest versions of packages and frameworks
func (m *Manager) GetLatestVersions() (*models.VersionConfig, error) {
	// Try to load from cache first
	cacheFile := filepath.Join(m.cacheDir, "versions.json")
	if fileExists(cacheFile) {
		if versions, err := m.loadVersionsFromCache(cacheFile); err == nil {
			// Check if cache is still valid (less than 24 hours old)
			if time.Since(versions.UpdatedAt) < 24*time.Hour {
				return versions, nil
			}
		}
	}

	// For now, return default versions (in a real implementation, this would fetch from registries)
	versions := &models.VersionConfig{
		Node:   "20.0.0",
		Go:     "1.22.0",
		Kotlin: "2.0.0",
		Swift:  "5.9.0",
		NextJS: "15.5.3",
		React:  "18.0.0",
		Packages: map[string]string{
			"express":     "4.18.0",
			"lodash":      "4.17.21",
			"tailwindcss": "3.4.0",
			"typescript":  "5.3.0",
		},
		UpdatedAt: time.Now(),
	}

	// Cache the versions
	if err := m.saveVersionsToCache(versions, cacheFile); err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: failed to cache versions: %v\n", err)
	}

	return versions, nil
}

// MergeConfigs merges base configuration with override values
func (m *Manager) MergeConfigs(base, override *models.ProjectConfig) *models.ProjectConfig {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	// Create a copy of base config
	merged := *base

	// Merge basic fields
	if override.Name != "" {
		merged.Name = override.Name
	}
	if override.Organization != "" {
		merged.Organization = override.Organization
	}
	if override.Description != "" {
		merged.Description = override.Description
	}
	if override.License != "" {
		merged.License = override.License
	}
	if override.Author != "" {
		merged.Author = override.Author
	}
	if override.Email != "" {
		merged.Email = override.Email
	}
	if override.Repository != "" {
		merged.Repository = override.Repository
	}
	if override.OutputPath != "" {
		merged.OutputPath = override.OutputPath
	}
	if override.GeneratorVersion != "" {
		merged.GeneratorVersion = override.GeneratorVersion
	}

	// Merge components (override takes precedence for each component)
	merged.Components = m.mergeComponents(base.Components, override.Components)

	// Merge versions
	if override.Versions != nil {
		merged.Versions = m.mergeVersions(base.Versions, override.Versions)
	}

	// Merge custom variables
	if override.CustomVars != nil {
		if merged.CustomVars == nil {
			merged.CustomVars = make(map[string]string)
		}
		for k, v := range override.CustomVars {
			merged.CustomVars[k] = v
		}
	}

	// Update generation timestamp
	merged.GeneratedAt = time.Now()

	return &merged
}

// SaveConfig saves configuration to a file
func (m *Manager) SaveConfig(config *models.ProjectConfig, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Determine format based on file extension
	ext := strings.ToLower(filepath.Ext(path))

	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
	default:
		// Default to YAML
		data, err = yaml.Marshal(config)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration from a file
func (m *Manager) LoadConfig(path string) (*models.ProjectConfig, error) {
	if !fileExists(path) {
		return nil, fmt.Errorf("configuration file does not exist: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config models.ProjectConfig
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	default:
		// Try YAML first, then JSON
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			err = json.Unmarshal(data, &config)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &config, nil
}

// Helper methods

func (m *Manager) mergeComponents(base, override models.Components) models.Components {
	merged := base

	// Frontend components
	if override.Frontend.MainApp {
		merged.Frontend.MainApp = true
	}
	if override.Frontend.Home {
		merged.Frontend.Home = true
	}
	if override.Frontend.Admin {
		merged.Frontend.Admin = true
	}

	// Backend components
	if override.Backend.API {
		merged.Backend.API = true
	}

	// Mobile components
	if override.Mobile.Android {
		merged.Mobile.Android = true
	}
	if override.Mobile.IOS {
		merged.Mobile.IOS = true
	}

	// Infrastructure components
	if override.Infrastructure.Docker {
		merged.Infrastructure.Docker = true
	}
	if override.Infrastructure.Kubernetes {
		merged.Infrastructure.Kubernetes = true
	}
	if override.Infrastructure.Terraform {
		merged.Infrastructure.Terraform = true
	}

	return merged
}

func (m *Manager) mergeVersions(base, override *models.VersionConfig) *models.VersionConfig {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := *base

	if override.Node != "" {
		merged.Node = override.Node
	}
	if override.Go != "" {
		merged.Go = override.Go
	}
	if override.Kotlin != "" {
		merged.Kotlin = override.Kotlin
	}
	if override.Swift != "" {
		merged.Swift = override.Swift
	}
	if override.NextJS != "" {
		merged.NextJS = override.NextJS
	}
	if override.React != "" {
		merged.React = override.React
	}

	// Merge packages
	if override.Packages != nil {
		if merged.Packages == nil {
			merged.Packages = make(map[string]string)
		}
		for k, v := range override.Packages {
			merged.Packages[k] = v
		}
	}

	merged.UpdatedAt = time.Now()
	return &merged
}

func (m *Manager) loadVersionsFromCache(cacheFile string) (*models.VersionConfig, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var versions models.VersionConfig
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}

func (m *Manager) saveVersionsToCache(versions *models.VersionConfig, cacheFile string) error {
	// Ensure cache directory exists
	dir := filepath.Dir(cacheFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
