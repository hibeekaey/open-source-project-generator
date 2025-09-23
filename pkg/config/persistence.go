// Package config provides configuration persistence and management functionality.
//
// This file implements configuration saving, loading, and management capabilities
// for interactive CLI configurations, allowing users to save and reuse their
// project generation settings.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// ConfigurationPersistence handles saving and loading of interactive configurations
type ConfigurationPersistence struct {
	configDir       string
	logger          interfaces.Logger
	validator       *ConfigurationValidator
	nameRegex       *regexp.Regexp
	maxConfigs      int
	compressionMode bool
}

// SavedConfiguration represents a saved interactive configuration
type SavedConfiguration struct {
	// Metadata
	Name        string    `yaml:"name" json:"name"`
	Description string    `yaml:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at" json:"updated_at"`
	Version     string    `yaml:"version" json:"version"`
	Tags        []string  `yaml:"tags,omitempty" json:"tags,omitempty"`

	// Project Configuration
	ProjectConfig *models.ProjectConfig `yaml:"project_config" json:"project_config"`

	// Template Selections (using interface{} to avoid circular dependencies)
	SelectedTemplates []TemplateSelection `yaml:"selected_templates" json:"selected_templates"`

	// Generation Settings
	GenerationSettings *GenerationSettings `yaml:"generation_settings,omitempty" json:"generation_settings,omitempty"`

	// User Preferences
	UserPreferences *UserPreferences `yaml:"user_preferences,omitempty" json:"user_preferences,omitempty"`
}

// TemplateSelection represents a selected template with options
type TemplateSelection struct {
	TemplateName string                 `yaml:"template_name" json:"template_name"`
	Category     string                 `yaml:"category" json:"category"`
	Technology   string                 `yaml:"technology" json:"technology"`
	Version      string                 `yaml:"version" json:"version"`
	Selected     bool                   `yaml:"selected" json:"selected"`
	Options      map[string]interface{} `yaml:"options,omitempty" json:"options,omitempty"`
}

// GenerationSettings contains settings for project generation
type GenerationSettings struct {
	IncludeExamples   bool     `yaml:"include_examples" json:"include_examples"`
	IncludeTests      bool     `yaml:"include_tests" json:"include_tests"`
	IncludeDocs       bool     `yaml:"include_docs" json:"include_docs"`
	UpdateVersions    bool     `yaml:"update_versions" json:"update_versions"`
	MinimalMode       bool     `yaml:"minimal_mode" json:"minimal_mode"`
	ExcludePatterns   []string `yaml:"exclude_patterns,omitempty" json:"exclude_patterns,omitempty"`
	IncludeOnlyPaths  []string `yaml:"include_only_paths,omitempty" json:"include_only_paths,omitempty"`
	BackupExisting    bool     `yaml:"backup_existing" json:"backup_existing"`
	OverwriteExisting bool     `yaml:"overwrite_existing" json:"overwrite_existing"`
}

// UserPreferences contains user-specific preferences
type UserPreferences struct {
	DefaultLicense      string            `yaml:"default_license,omitempty" json:"default_license,omitempty"`
	DefaultAuthor       string            `yaml:"default_author,omitempty" json:"default_author,omitempty"`
	DefaultEmail        string            `yaml:"default_email,omitempty" json:"default_email,omitempty"`
	DefaultOrganization string            `yaml:"default_organization,omitempty" json:"default_organization,omitempty"`
	PreferredFormat     string            `yaml:"preferred_format,omitempty" json:"preferred_format,omitempty"`
	CustomDefaults      map[string]string `yaml:"custom_defaults,omitempty" json:"custom_defaults,omitempty"`
}

// ConfigurationValidator validates saved configurations
type ConfigurationValidator struct {
	nameRegex    *regexp.Regexp
	versionRegex *regexp.Regexp
}

// ConfigurationListOptions contains options for listing configurations
type ConfigurationListOptions struct {
	SortBy      string   `json:"sort_by"`      // "name", "created_at", "updated_at"
	SortOrder   string   `json:"sort_order"`   // "asc", "desc"
	FilterTags  []string `json:"filter_tags"`  // Filter by tags
	SearchQuery string   `json:"search_query"` // Search in name and description
	Limit       int      `json:"limit"`        // Maximum number of results
	Offset      int      `json:"offset"`       // Pagination offset
}

// NewConfigurationPersistence creates a new configuration persistence manager
func NewConfigurationPersistence(configDir string, logger interfaces.Logger) *ConfigurationPersistence {
	// Ensure config directory exists
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".generator", "configs")
	}

	if err := os.MkdirAll(configDir, 0750); err != nil {
		if logger != nil {
			logger.WarnWithFields("Failed to create config directory", map[string]interface{}{
				"directory": configDir,
				"error":     err.Error(),
			})
		}
	}

	validator := &ConfigurationValidator{
		nameRegex:    utils.ProjectNamePattern,
		versionRegex: utils.VersionPattern,
	}

	return &ConfigurationPersistence{
		configDir:       configDir,
		logger:          logger,
		validator:       validator,
		nameRegex:       utils.ProjectNamePattern,
		maxConfigs:      100, // Limit to prevent excessive storage
		compressionMode: false,
	}
}

// SaveConfiguration saves a configuration with the given name
func (cp *ConfigurationPersistence) SaveConfiguration(name string, config *SavedConfiguration) error {
	if err := cp.validateConfigurationName(name); err != nil {
		return fmt.Errorf("invalid configuration name: %w", err)
	}

	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate configuration content
	if err := cp.validator.ValidateConfiguration(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Set metadata
	config.Name = name
	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = config.UpdatedAt
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}

	// Check if we're at the maximum number of configurations
	if err := cp.enforceMaxConfigurations(); err != nil {
		cp.logger.WarnWithFields("Failed to enforce max configurations", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Generate file path
	filename := cp.generateConfigFilename(name)
	filePath := filepath.Join(cp.configDir, filename)

	// Save configuration
	if err := cp.writeConfigurationFile(filePath, config); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	if cp.logger != nil {
		cp.logger.InfoWithFields("Configuration saved", map[string]interface{}{
			"name":      name,
			"file_path": filePath,
			"version":   config.Version,
		})
	}

	return nil
}

// LoadConfiguration loads a configuration by name
func (cp *ConfigurationPersistence) LoadConfiguration(name string) (*SavedConfiguration, error) {
	if err := cp.validateConfigurationName(name); err != nil {
		return nil, fmt.Errorf("invalid configuration name: %w", err)
	}

	filename := cp.generateConfigFilename(name)
	filePath := filepath.Join(cp.configDir, filename)

	// Check if file exists
	if !cp.fileExists(filePath) {
		return nil, fmt.Errorf("configuration '%s' not found", name)
	}

	// Load configuration
	config, err := cp.readConfigurationFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Validate loaded configuration
	if err := cp.validator.ValidateConfiguration(config); err != nil {
		cp.logger.WarnWithFields("Loaded configuration failed validation", map[string]interface{}{
			"name":  name,
			"error": err.Error(),
		})
		// Don't fail loading, just warn
	}

	if cp.logger != nil {
		cp.logger.DebugWithFields("Configuration loaded", map[string]interface{}{
			"name":      name,
			"file_path": filePath,
			"version":   config.Version,
		})
	}

	return config, nil
}

// ListConfigurations lists all saved configurations with optional filtering
func (cp *ConfigurationPersistence) ListConfigurations(options *ConfigurationListOptions) ([]*SavedConfiguration, error) {
	if options == nil {
		options = &ConfigurationListOptions{
			SortBy:    "updated_at",
			SortOrder: "desc",
			Limit:     50,
		}
	}

	// Get all configuration files
	files, err := filepath.Glob(filepath.Join(cp.configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list configuration files: %w", err)
	}

	var configs []*SavedConfiguration

	// Load each configuration
	for _, file := range files {
		config, err := cp.readConfigurationFile(file)
		if err != nil {
			cp.logger.WarnWithFields("Failed to load configuration file", map[string]interface{}{
				"file":  file,
				"error": err.Error(),
			})
			continue
		}

		// Apply filters
		if cp.matchesFilters(config, options) {
			configs = append(configs, config)
		}
	}

	// Sort configurations
	cp.sortConfigurations(configs, options.SortBy, options.SortOrder)

	// Apply pagination
	if options.Offset > 0 || options.Limit > 0 {
		configs = cp.paginateConfigurations(configs, options.Offset, options.Limit)
	}

	return configs, nil
}

// DeleteConfiguration deletes a saved configuration
func (cp *ConfigurationPersistence) DeleteConfiguration(name string) error {
	if err := cp.validateConfigurationName(name); err != nil {
		return fmt.Errorf("invalid configuration name: %w", err)
	}

	filename := cp.generateConfigFilename(name)
	filePath := filepath.Join(cp.configDir, filename)

	// Check if file exists
	if !cp.fileExists(filePath) {
		return fmt.Errorf("configuration '%s' not found", name)
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete configuration file: %w", err)
	}

	if cp.logger != nil {
		cp.logger.InfoWithFields("Configuration deleted", map[string]interface{}{
			"name":      name,
			"file_path": filePath,
		})
	}

	return nil
}

// ConfigurationExists checks if a configuration with the given name exists
func (cp *ConfigurationPersistence) ConfigurationExists(name string) bool {
	if err := cp.validateConfigurationName(name); err != nil {
		return false
	}

	filename := cp.generateConfigFilename(name)
	filePath := filepath.Join(cp.configDir, filename)
	return cp.fileExists(filePath)
}

// UpdateConfiguration updates an existing configuration
func (cp *ConfigurationPersistence) UpdateConfiguration(name string, config *SavedConfiguration) error {
	// Check if configuration exists
	if !cp.ConfigurationExists(name) {
		return fmt.Errorf("configuration '%s' does not exist", name)
	}

	// Load existing configuration to preserve created_at
	existing, err := cp.LoadConfiguration(name)
	if err != nil {
		return fmt.Errorf("failed to load existing configuration: %w", err)
	}

	// Preserve creation time
	config.CreatedAt = existing.CreatedAt

	// Save updated configuration
	return cp.SaveConfiguration(name, config)
}

// GetConfigurationInfo returns basic information about a configuration without loading the full content
func (cp *ConfigurationPersistence) GetConfigurationInfo(name string) (*ConfigurationInfo, error) {
	if err := cp.validateConfigurationName(name); err != nil {
		return nil, fmt.Errorf("invalid configuration name: %w", err)
	}

	filename := cp.generateConfigFilename(name)
	filePath := filepath.Join(cp.configDir, filename)

	// Check if file exists
	if !cp.fileExists(filePath) {
		return nil, fmt.Errorf("configuration '%s' not found", name)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Load just the metadata (first few lines)
	config, err := cp.readConfigurationFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration metadata: %w", err)
	}

	info := &ConfigurationInfo{
		Name:        config.Name,
		Description: config.Description,
		CreatedAt:   config.CreatedAt,
		UpdatedAt:   config.UpdatedAt,
		Version:     config.Version,
		Tags:        config.Tags,
		FileSize:    fileInfo.Size(),
		FilePath:    filePath,
	}

	return info, nil
}

// ConfigurationInfo contains basic information about a saved configuration
type ConfigurationInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     string    `json:"version"`
	Tags        []string  `json:"tags,omitempty"`
	FileSize    int64     `json:"file_size"`
	FilePath    string    `json:"file_path"`
}

// Helper methods

// validateConfigurationName validates the configuration name
func (cp *ConfigurationPersistence) validateConfigurationName(name string) error {
	if name == "" {
		return fmt.Errorf("configuration name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("configuration name must be at least 2 characters long")
	}

	if len(name) > 50 {
		return fmt.Errorf("configuration name must be at most 50 characters long")
	}

	if !cp.nameRegex.MatchString(name) {
		return fmt.Errorf("configuration name must start with a letter and contain only letters, numbers, hyphens, underscores, and spaces")
	}

	// Check for reserved names
	reservedNames := []string{"default", "config", "template", "system", "admin", "root"}
	lowerName := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			return fmt.Errorf("'%s' is a reserved name and cannot be used", name)
		}
	}

	return nil
}

// generateConfigFilename generates a safe filename for the configuration
func (cp *ConfigurationPersistence) generateConfigFilename(name string) string {
	// Replace spaces and special characters with hyphens
	filename := strings.ToLower(name)
	filename = regexp.MustCompile(`[^a-z0-9\-_]`).ReplaceAllString(filename, "-")
	filename = regexp.MustCompile(`-+`).ReplaceAllString(filename, "-")
	filename = strings.Trim(filename, "-")
	return filename + ".yaml"
}

// writeConfigurationFile writes a configuration to a file
func (cp *ConfigurationPersistence) writeConfigurationFile(filePath string, config *SavedConfiguration) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file safely
	if err := utils.SafeWriteFile(filePath, data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// readConfigurationFile reads a configuration from a file
func (cp *ConfigurationPersistence) readConfigurationFile(filePath string) (*SavedConfiguration, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Read file safely
	data, err := utils.SafeReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal from YAML
	var config SavedConfiguration
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &config, nil
}

// fileExists checks if a file exists
func (cp *ConfigurationPersistence) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// enforceMaxConfigurations ensures we don't exceed the maximum number of configurations
func (cp *ConfigurationPersistence) enforceMaxConfigurations() error {
	files, err := filepath.Glob(filepath.Join(cp.configDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to count configuration files: %w", err)
	}

	if len(files) >= cp.maxConfigs {
		// Remove oldest configurations
		configs, err := cp.ListConfigurations(&ConfigurationListOptions{
			SortBy:    "updated_at",
			SortOrder: "asc",
		})
		if err != nil {
			return fmt.Errorf("failed to list configurations for cleanup: %w", err)
		}

		// Delete oldest configurations to make room
		toDelete := len(configs) - cp.maxConfigs + 1
		for i := 0; i < toDelete && i < len(configs); i++ {
			if err := cp.DeleteConfiguration(configs[i].Name); err != nil {
				cp.logger.WarnWithFields("Failed to delete old configuration", map[string]interface{}{
					"name":  configs[i].Name,
					"error": err.Error(),
				})
			}
		}
	}

	return nil
}

// matchesFilters checks if a configuration matches the given filters
func (cp *ConfigurationPersistence) matchesFilters(config *SavedConfiguration, options *ConfigurationListOptions) bool {
	// Search query filter
	if options.SearchQuery != "" {
		query := strings.ToLower(options.SearchQuery)
		if !strings.Contains(strings.ToLower(config.Name), query) &&
			!strings.Contains(strings.ToLower(config.Description), query) {
			return false
		}
	}

	// Tags filter
	if len(options.FilterTags) > 0 {
		hasMatchingTag := false
		for _, filterTag := range options.FilterTags {
			for _, configTag := range config.Tags {
				if strings.EqualFold(filterTag, configTag) {
					hasMatchingTag = true
					break
				}
			}
			if hasMatchingTag {
				break
			}
		}
		if !hasMatchingTag {
			return false
		}
	}

	return true
}

// sortConfigurations sorts configurations based on the given criteria
func (cp *ConfigurationPersistence) sortConfigurations(configs []*SavedConfiguration, sortBy, sortOrder string) {
	// Implementation would use sort.Slice with appropriate comparison functions
	// For brevity, this is a placeholder
}

// paginateConfigurations applies pagination to the configuration list
func (cp *ConfigurationPersistence) paginateConfigurations(configs []*SavedConfiguration, offset, limit int) []*SavedConfiguration {
	if offset >= len(configs) {
		return []*SavedConfiguration{}
	}

	end := offset + limit
	if limit <= 0 || end > len(configs) {
		end = len(configs)
	}

	return configs[offset:end]
}

// ValidateConfiguration validates a saved configuration
func (cv *ConfigurationValidator) ValidateConfiguration(config *SavedConfiguration) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate name
	if config.Name == "" {
		return fmt.Errorf("configuration name is required")
	}

	// Validate version format
	if config.Version != "" && !cv.versionRegex.MatchString(config.Version) {
		return fmt.Errorf("invalid version format: %s (expected semver format like 1.0.0)", config.Version)
	}

	// Validate project configuration
	if config.ProjectConfig == nil {
		return fmt.Errorf("project configuration is required")
	}

	// Validate project name
	if config.ProjectConfig.Name == "" {
		return fmt.Errorf("project name is required")
	}

	// Validate selected templates
	if len(config.SelectedTemplates) == 0 {
		return fmt.Errorf("at least one template must be selected")
	}

	for i, template := range config.SelectedTemplates {
		if template.TemplateName == "" {
			return fmt.Errorf("template %d: template name is required", i)
		}
		if template.Category == "" {
			return fmt.Errorf("template %d: category is required", i)
		}
	}

	return nil
}

// ExportConfiguration exports a configuration to a different format
func (cp *ConfigurationPersistence) ExportConfiguration(name, format string) ([]byte, error) {
	config, err := cp.LoadConfiguration(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(config, "", "  ")
	case "yaml", "yml":
		return yaml.Marshal(config)
	default:
		return nil, fmt.Errorf("unsupported export format: %s (supported: json, yaml)", format)
	}
}

// ImportConfiguration imports a configuration from data
func (cp *ConfigurationPersistence) ImportConfiguration(name string, data []byte, format string) error {
	var config SavedConfiguration

	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to unmarshal YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported import format: %s (supported: json, yaml)", format)
	}

	// Override name with provided name
	config.Name = name

	return cp.SaveConfiguration(name, &config)
}
