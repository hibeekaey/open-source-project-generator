package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

// Manager implements the ConfigManager interface
type Manager struct {
	cacheDir     string
	defaultsPath string
	configPath   string
	settings     map[string]interface{}
	schema       *interfaces.ConfigSchema
	envPrefix    string
}

// NewManager creates a new configuration manager
func NewManager(cacheDir, defaultsPath string) interfaces.ConfigManager {
	manager := &Manager{
		cacheDir:     cacheDir,
		defaultsPath: defaultsPath,
		settings:     make(map[string]interface{}),
		envPrefix:    "GENERATOR",
	}

	// Initialize schema
	manager.schema = manager.createConfigSchema()

	// Set default config path
	if cacheDir != "" {
		manager.configPath = filepath.Join(cacheDir, "config.yaml")
	}

	return manager
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
		Name:         "my-project",
		Organization: "my-org",
		License:      "MIT",
		OutputPath:   "./output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   false,
					Admin:  false,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
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
			Node: "20.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"react":      "18.2.0",
				"next":       "13.4.0",
				"typescript": "5.0.0",
			},
		},
	}, nil
}

// LoadConfig loads configuration from a file
func (m *Manager) LoadConfig(configPath string) (*models.ProjectConfig, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(configPath); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	content, err := utils.SafeReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("🚫 %s %s", 
			"Unable to read configuration file.", 
			"Check if the file exists and has proper permissions")
	}

	var config models.ProjectConfig
	ext := strings.ToLower(filepath.Ext(configPath))

	switch ext {
	case ".json":
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("🚫 %s %s", 
				"Invalid JSON configuration file.", 
				"Check the file syntax and fix any JSON formatting errors")
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("🚫 %s %s", 
				"Invalid YAML configuration file.", 
				"Check the file syntax and fix any YAML formatting errors")
		}
	default:
		return nil, fmt.Errorf("🚫 Unsupported configuration file format '%s'. Use .json, .yaml, or .yml files", ext)
	}

	return &config, nil
}

// SaveConfig saves configuration to a file
func (m *Manager) SaveConfig(config *models.ProjectConfig, configPath string) error {
	ext := strings.ToLower(filepath.Ext(configPath))
	var content []byte
	var err error

	switch ext {
	case ".json":
		content, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}
	case ".yaml", ".yml":
		content, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err := utils.SafeWriteFile(configPath, content); err != nil {
		return fmt.Errorf("🚫 %s %s", 
			"Unable to write configuration file.", 
			"Check file permissions and available disk space")
	}

	return nil
}

// ValidateConfig validates a project configuration
func (m *Manager) ValidateConfig(config *models.ProjectConfig) error {
	validator := NewConfigValidator(m.schema)
	result, err := validator.ValidateProjectConfig(config)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid {
		var errorMessages []string
		for _, validationError := range result.Errors {
			errorMessages = append(errorMessages, validationError.Message)
		}
		return fmt.Errorf("🚫 Configuration validation failed: %s", strings.Join(errorMessages, "; "))
	}

	// Normalize license - set default if empty or unsupported
	m.normalizeLicense(config)

	return nil
}

// normalizeLicense normalizes the license field, setting default for unsupported types
func (m *Manager) normalizeLicense(config *models.ProjectConfig) {
	supportedLicenses := []string{
		"MIT",
		"Apache-2.0",
		"GPL-3.0",
		"BSD-3-Clause",
	}

	// Set default if empty
	if config.License == "" {
		config.License = "MIT"
		return
	}

	// Check if license is supported
	for _, supported := range supportedLicenses {
		if config.License == supported {
			return // License is supported
		}
	}

	// License not supported - set to default
	config.License = "MIT"
}

// Helper function to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetSetting gets a configuration setting
func (m *Manager) GetSetting(key string) (any, error) {
	if value, exists := m.settings[key]; exists {
		return value, nil
	}

	// Try to get from environment variables
	envKey := m.envPrefix + "_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue, nil
	}

	// Try to get from schema defaults
	if m.schema != nil {
		if prop, exists := m.schema.Properties[key]; exists && prop.Default != nil {
			return prop.Default, nil
		}
	}

	return nil, fmt.Errorf("setting '%s' not found", key)
}

// SetSetting sets a configuration setting
func (m *Manager) SetSetting(key string, value any) error {
	// Validate the setting against schema if available
	if m.schema != nil {
		if err := m.validateSetting(key, value); err != nil {
			return fmt.Errorf("invalid setting value: %w", err)
		}
	}

	m.settings[key] = value
	return nil
}

// ValidateSettings validates configuration settings
func (m *Manager) ValidateSettings() error {
	if m.schema == nil {
		return nil // No schema to validate against
	}

	var errors []string

	// Check required fields
	for _, required := range m.schema.Required {
		if _, exists := m.settings[required]; !exists {
			// Check if it has a default value
			if prop, propExists := m.schema.Properties[required]; propExists && prop.Default != nil {
				m.settings[required] = prop.Default
			} else {
				errors = append(errors, fmt.Sprintf("required setting '%s' is missing", required))
			}
		}
	}

	// Validate each setting
	for key, value := range m.settings {
		if err := m.validateSetting(key, value); err != nil {
			errors = append(errors, fmt.Sprintf("setting '%s': %v", key, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

// LoadFromEnvironment loads configuration from environment variables
func (m *Manager) LoadFromEnvironment() (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{}

	// Load basic project information from environment
	if name := os.Getenv(m.envPrefix + "_NAME"); name != "" {
		config.Name = name
	}
	if org := os.Getenv(m.envPrefix + "_ORGANIZATION"); org != "" {
		config.Organization = org
	}
	if desc := os.Getenv(m.envPrefix + "_DESCRIPTION"); desc != "" {
		config.Description = desc
	}
	if license := os.Getenv(m.envPrefix + "_LICENSE"); license != "" {
		config.License = license
	}
	if author := os.Getenv(m.envPrefix + "_AUTHOR"); author != "" {
		config.Author = author
	}
	if email := os.Getenv(m.envPrefix + "_EMAIL"); email != "" {
		config.Email = email
	}
	if repo := os.Getenv(m.envPrefix + "_REPOSITORY"); repo != "" {
		config.Repository = repo
	}
	if output := os.Getenv(m.envPrefix + "_OUTPUT_PATH"); output != "" {
		config.OutputPath = output
	}

	// Load component settings from environment
	m.loadComponentsFromEnv(&config.Components)

	// Load version settings from environment
	config.Versions = m.loadVersionsFromEnv()

	return config, nil
}

// MergeConfigurations merges multiple configurations with proper precedence
func (m *Manager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
	if len(configs) == 0 {
		return nil
	}

	if len(configs) == 1 {
		return configs[0]
	}

	// Start with the first config as base
	result := *configs[0]

	// Merge each subsequent config
	for i := 1; i < len(configs); i++ {
		if configs[i] == nil {
			continue
		}

		m.mergeConfig(&result, configs[i])
	}

	return &result
}

// GetConfigSchema returns the configuration schema
func (m *Manager) GetConfigSchema() *interfaces.ConfigSchema {
	return m.schema
}

// GetConfigSources returns configuration sources
func (m *Manager) GetConfigSources() ([]interfaces.ConfigSource, error) {
	var sources []interfaces.ConfigSource

	// Check defaults file
	if m.defaultsPath != "" {
		valid := fileExists(m.defaultsPath)
		sources = append(sources, interfaces.ConfigSource{
			Type:     "defaults",
			Location: m.defaultsPath,
			Priority: 1,
			Valid:    valid,
		})
	}

	// Check config file
	if m.configPath != "" {
		valid := fileExists(m.configPath)
		sources = append(sources, interfaces.ConfigSource{
			Type:     "file",
			Location: m.configPath,
			Priority: 2,
			Valid:    valid,
		})
	}

	// Check environment variables
	envVars := m.getRelevantEnvVars()
	if len(envVars) > 0 {
		sources = append(sources, interfaces.ConfigSource{
			Type:     "environment",
			Location: fmt.Sprintf("%d environment variables", len(envVars)),
			Priority: 3,
			Valid:    true,
		})
	}

	return sources, nil
}

// GetConfigLocation returns the configuration location
func (m *Manager) GetConfigLocation() string {
	return m.defaultsPath
}

// CreateDefaultConfig creates a default configuration file
func (m *Manager) CreateDefaultConfig(path string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load default configuration
	defaultConfig, err := m.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	// Save to specified path
	return m.SaveConfig(defaultConfig, path)
}

// BackupConfig backs up the configuration
func (m *Manager) BackupConfig(path string) error {
	if !fileExists(path) {
		return fmt.Errorf("config file does not exist: %s", path)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", path, timestamp)

	// Read original file
	content, err := utils.SafeReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Write backup file
	if err := utils.SafeWriteFile(backupPath, content); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// RestoreConfig restores configuration from backup
func (m *Manager) RestoreConfig(backupPath string) error {
	if !fileExists(backupPath) {
		return fmt.Errorf("backup file does not exist: %s", backupPath)
	}

	// Determine original path (remove .backup_timestamp suffix)
	originalPath := regexp.MustCompile(`\.backup_\d{8}_\d{6}$`).ReplaceAllString(backupPath, "")

	// Read backup file
	content, err := utils.SafeReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Write to original location
	if err := utils.SafeWriteFile(originalPath, content); err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	return nil
}

// LoadEnvironmentVariables loads environment variables
func (m *Manager) LoadEnvironmentVariables() map[string]string {
	return m.getRelevantEnvVars()
}

// SetEnvironmentDefaults sets environment defaults
func (m *Manager) SetEnvironmentDefaults() error {
	defaults := map[string]string{
		m.envPrefix + "_LICENSE":     "MIT",
		m.envPrefix + "_OUTPUT_PATH": "./output",
	}

	for key, value := range defaults {
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("failed to set environment default %s: %w", key, err)
			}
		}
	}

	return nil
}

// GetEnvironmentPrefix returns the environment prefix
func (m *Manager) GetEnvironmentPrefix() string {
	return m.envPrefix
}

// createConfigSchema creates the configuration schema
func (m *Manager) createConfigSchema() *interfaces.ConfigSchema {
	return &interfaces.ConfigSchema{
		Version:     "1.0.0",
		Title:       "Project Generator Configuration",
		Description: "Configuration schema for the Open Source Project Generator",
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Project name",
				Required:    true,
				Pattern:     "^[a-zA-Z0-9_-]+$",
				MinLength:   &[]int{1}[0],
				MaxLength:   &[]int{100}[0],
				Examples:    []string{"my-project", "awesome_app"},
			},
			"organization": {
				Type:        "string",
				Description: "Organization or author name",
				Required:    true,
				MinLength:   &[]int{1}[0],
				MaxLength:   &[]int{100}[0],
				Examples:    []string{"my-org", "john-doe"},
			},
			"description": {
				Type:        "string",
				Description: "Project description",
				Required:    false,
				MaxLength:   &[]int{500}[0],
				Examples:    []string{"A sample project", "My awesome application"},
			},
			"license": {
				Type:        "string",
				Description: "Project license",
				Default:     "MIT",
				Enum:        []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"},
				Examples:    []string{"MIT", "Apache-2.0"},
			},
			"author": {
				Type:        "string",
				Description: "Project author",
				Required:    false,
				MaxLength:   &[]int{100}[0],
				Examples:    []string{"John Doe", "jane.smith"},
			},
			"email": {
				Type:        "string",
				Description: "Author email",
				Required:    false,
				Pattern:     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
				Examples:    []string{"john@example.com", "jane.smith@company.org"},
			},
			"repository": {
				Type:        "string",
				Description: "Repository URL",
				Required:    false,
				Pattern:     "^https?://.*$",
				Examples:    []string{"https://github.com/user/repo", "https://gitlab.com/user/repo"},
			},
			"output_path": {
				Type:        "string",
				Description: "Output directory path",
				Default:     "./output",
				Required:    false,
				Examples:    []string{"./output", "/tmp/projects", "../my-project"},
			},
		},
		Required: []string{"name", "organization"},
	}
}

// validateSetting validates a single setting against the schema
func (m *Manager) validateSetting(key string, value interface{}) error {
	if m.schema == nil {
		return nil
	}

	prop, exists := m.schema.Properties[key]
	if !exists {
		// Allow unknown settings for flexibility
		return nil
	}

	// Type validation
	if err := m.validateType(value, prop.Type); err != nil {
		return err
	}

	// String-specific validations
	if prop.Type == "string" {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}

		// Length validation
		if prop.MinLength != nil && len(str) < *prop.MinLength {
			return fmt.Errorf("minimum length is %d, got %d", *prop.MinLength, len(str))
		}
		if prop.MaxLength != nil && len(str) > *prop.MaxLength {
			return fmt.Errorf("maximum length is %d, got %d", *prop.MaxLength, len(str))
		}

		// Pattern validation
		if prop.Pattern != "" {
			matched, err := regexp.MatchString(prop.Pattern, str)
			if err != nil {
				return fmt.Errorf("invalid pattern: %w", err)
			}
			if !matched {
				return fmt.Errorf("value does not match pattern %s", prop.Pattern)
			}
		}

		// Enum validation
		if len(prop.Enum) > 0 {
			valid := false
			for _, enumValue := range prop.Enum {
				if str == enumValue {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("value must be one of: %s", strings.Join(prop.Enum, ", "))
			}
		}
	}

	// Number-specific validations
	if prop.Type == "number" || prop.Type == "integer" {
		num, err := m.convertToFloat64(value)
		if err != nil {
			return err
		}

		if prop.Minimum != nil && num < *prop.Minimum {
			return fmt.Errorf("minimum value is %f, got %f", *prop.Minimum, num)
		}
		if prop.Maximum != nil && num > *prop.Maximum {
			return fmt.Errorf("maximum value is %f, got %f", *prop.Maximum, num)
		}
	}

	return nil
}

// validateType validates the type of a value
func (m *Manager) validateType(value interface{}, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "number":
		if _, err := m.convertToFloat64(value); err != nil {
			return fmt.Errorf("expected number, got %T", value)
		}
	case "integer":
		if _, err := m.convertToInt(value); err != nil {
			return fmt.Errorf("expected integer, got %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		if reflect.TypeOf(value).Kind() != reflect.Slice {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if reflect.TypeOf(value).Kind() != reflect.Map && reflect.TypeOf(value).Kind() != reflect.Struct {
			return fmt.Errorf("expected object, got %T", value)
		}
	}
	return nil
}

// convertToFloat64 converts various numeric types to float64
func (m *Manager) convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to number", value)
	}
}

// convertToInt converts various types to int
func (m *Manager) convertToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to integer", value)
	}
}

// loadComponentsFromEnv loads component configuration from environment variables
func (m *Manager) loadComponentsFromEnv(components *models.Components) {
	// Frontend components
	if nextjsApp := os.Getenv(m.envPrefix + "_NEXTJS_APP"); nextjsApp != "" {
		if val, err := strconv.ParseBool(nextjsApp); err == nil {
			components.Frontend.NextJS.App = val
		}
	}
	if nextjsHome := os.Getenv(m.envPrefix + "_NEXTJS_HOME"); nextjsHome != "" {
		if val, err := strconv.ParseBool(nextjsHome); err == nil {
			components.Frontend.NextJS.Home = val
		}
	}
	if nextjsAdmin := os.Getenv(m.envPrefix + "_NEXTJS_ADMIN"); nextjsAdmin != "" {
		if val, err := strconv.ParseBool(nextjsAdmin); err == nil {
			components.Frontend.NextJS.Admin = val
		}
	}
	if nextjsShared := os.Getenv(m.envPrefix + "_NEXTJS_SHARED"); nextjsShared != "" {
		if val, err := strconv.ParseBool(nextjsShared); err == nil {
			components.Frontend.NextJS.Shared = val
		}
	}

	// Backend components
	if goGin := os.Getenv(m.envPrefix + "_GO_GIN"); goGin != "" {
		if val, err := strconv.ParseBool(goGin); err == nil {
			components.Backend.GoGin = val
		}
	}

	// Mobile components
	if android := os.Getenv(m.envPrefix + "_ANDROID"); android != "" {
		if val, err := strconv.ParseBool(android); err == nil {
			components.Mobile.Android = val
		}
	}
	if ios := os.Getenv(m.envPrefix + "_IOS"); ios != "" {
		if val, err := strconv.ParseBool(ios); err == nil {
			components.Mobile.IOS = val
		}
	}

	// Infrastructure components
	if docker := os.Getenv(m.envPrefix + "_DOCKER"); docker != "" {
		if val, err := strconv.ParseBool(docker); err == nil {
			components.Infrastructure.Docker = val
		}
	}
	if kubernetes := os.Getenv(m.envPrefix + "_KUBERNETES"); kubernetes != "" {
		if val, err := strconv.ParseBool(kubernetes); err == nil {
			components.Infrastructure.Kubernetes = val
		}
	}
	if terraform := os.Getenv(m.envPrefix + "_TERRAFORM"); terraform != "" {
		if val, err := strconv.ParseBool(terraform); err == nil {
			components.Infrastructure.Terraform = val
		}
	}
}

// loadVersionsFromEnv loads version configuration from environment variables
func (m *Manager) loadVersionsFromEnv() *models.VersionConfig {
	versions := &models.VersionConfig{
		Packages: make(map[string]string),
	}

	if node := os.Getenv(m.envPrefix + "_NODE_VERSION"); node != "" {
		versions.Node = node
	}
	if goVersion := os.Getenv(m.envPrefix + "_GO_VERSION"); goVersion != "" {
		versions.Go = goVersion
	}

	// Load package versions
	packagePrefixes := []string{"REACT", "NEXT", "TYPESCRIPT", "GIN", "GORM"}
	for _, prefix := range packagePrefixes {
		envKey := m.envPrefix + "_" + prefix + "_VERSION"
		if version := os.Getenv(envKey); version != "" {
			packageName := strings.ToLower(prefix)
			versions.Packages[packageName] = version
		}
	}

	return versions
}

// mergeConfig merges source config into target config
func (m *Manager) mergeConfig(target, source *models.ProjectConfig) {
	// Merge basic fields (source takes precedence if not empty)
	if source.Name != "" {
		target.Name = source.Name
	}
	if source.Organization != "" {
		target.Organization = source.Organization
	}
	if source.Description != "" {
		target.Description = source.Description
	}
	if source.License != "" {
		target.License = source.License
	}
	if source.Author != "" {
		target.Author = source.Author
	}
	if source.Email != "" {
		target.Email = source.Email
	}
	if source.Repository != "" {
		target.Repository = source.Repository
	}
	if source.OutputPath != "" {
		target.OutputPath = source.OutputPath
	}

	// Merge components
	m.mergeComponents(&target.Components, &source.Components)

	// Merge versions
	if source.Versions != nil {
		if target.Versions == nil {
			target.Versions = &models.VersionConfig{Packages: make(map[string]string)}
		}
		m.mergeVersions(target.Versions, source.Versions)
	}

	// Update metadata
	if !source.GeneratedAt.IsZero() {
		target.GeneratedAt = source.GeneratedAt
	}
	if source.GeneratorVersion != "" {
		target.GeneratorVersion = source.GeneratorVersion
	}
}

// mergeComponents merges component configurations
func (m *Manager) mergeComponents(target, source *models.Components) {
	// Merge frontend components
	if source.Frontend.NextJS.App {
		target.Frontend.NextJS.App = true
	}
	if source.Frontend.NextJS.Home {
		target.Frontend.NextJS.Home = true
	}
	if source.Frontend.NextJS.Admin {
		target.Frontend.NextJS.Admin = true
	}
	if source.Frontend.NextJS.Shared {
		target.Frontend.NextJS.Shared = true
	}

	// Merge backend components
	if source.Backend.GoGin {
		target.Backend.GoGin = true
	}

	// Merge mobile components
	if source.Mobile.Android {
		target.Mobile.Android = true
	}
	if source.Mobile.IOS {
		target.Mobile.IOS = true
	}

	// Merge infrastructure components
	if source.Infrastructure.Docker {
		target.Infrastructure.Docker = true
	}
	if source.Infrastructure.Kubernetes {
		target.Infrastructure.Kubernetes = true
	}
	if source.Infrastructure.Terraform {
		target.Infrastructure.Terraform = true
	}
}

// mergeVersions merges version configurations
func (m *Manager) mergeVersions(target, source *models.VersionConfig) {
	if source.Node != "" {
		target.Node = source.Node
	}
	if source.Go != "" {
		target.Go = source.Go
	}

	// Merge package versions
	if target.Packages == nil {
		target.Packages = make(map[string]string)
	}
	for pkg, version := range source.Packages {
		if version != "" {
			target.Packages[pkg] = version
		}
	}
}

// getRelevantEnvVars gets environment variables with the configured prefix
func (m *Manager) getRelevantEnvVars() map[string]string {
	envVars := make(map[string]string)
	prefix := m.envPrefix + "_"

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], prefix) {
			envVars[parts[0]] = parts[1]
		}
	}

	return envVars
}

// ValidateConfigDetailed validates a project configuration and returns detailed results
func (m *Manager) ValidateConfigDetailed(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	validator := NewConfigValidator(m.schema)
	return validator.ValidateProjectConfig(config)
}

// ValidateConfigFromFile validates a configuration file and returns detailed results
func (m *Manager) ValidateConfigFromFile(configPath string) (*interfaces.ConfigValidationResult, error) {
	config, err := m.LoadConfig(configPath)
	if err != nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:    "file",
					Value:    configPath,
					Type:     "file_error",
					Message:  fmt.Sprintf("failed to load configuration file: %v", err),
					Severity: "error",
					Rule:     "file_access",
				},
			},
			Summary: interfaces.ConfigValidationSummary{
				ErrorCount: 1,
			},
		}, nil
	}

	return m.ValidateConfigDetailed(config)
}

// GetValidationSchema returns the validation schema for external use
func (m *Manager) GetValidationSchema() *interfaces.ConfigSchema {
	return m.schema
}

// UpdateSchema updates the configuration schema
func (m *Manager) UpdateSchema(schema *interfaces.ConfigSchema) {
	m.schema = schema
}
