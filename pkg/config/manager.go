// Package config provides a unified configuration management system for the entire application.
//
// This package consolidates all configuration management logic from various packages into a single,
// comprehensive configuration system that can be used across the application.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

// intPtr returns a pointer to an int
func intPtr(i int) *int {
	return &i
}

// UnifiedConfigManager provides comprehensive configuration management
type UnifiedConfigManager struct {
	configDir    string
	defaultsPath string
	settings     map[string]interface{}
	schema       *interfaces.ConfigSchema
	envPrefix    string
	persistence  *UnifiedConfigurationPersistence
	validator    *UnifiedConfigValidator
}

// UnifiedConfigValidator provides unified configuration validation
type UnifiedConfigValidator struct {
	rules map[string]ConfigValidationRule
}

// ConfigValidationRule defines a configuration validation rule
type ConfigValidationRule struct {
	Name        string
	Validator   func(interface{}) error
	Message     string
	Suggestions []string
	Required    bool
}

// UnifiedConfigurationPersistence handles saving and loading configurations
type UnifiedConfigurationPersistence struct {
	configDir   string
	compression bool
	maxConfigs  int
}

// NewUnifiedConfigManager creates a new unified configuration manager
func NewUnifiedConfigManager(configDir, defaultsPath string) *UnifiedConfigManager {
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".generator")
	}

	manager := &UnifiedConfigManager{
		configDir:    configDir,
		defaultsPath: defaultsPath,
		settings:     make(map[string]interface{}),
		envPrefix:    "GENERATOR",
		persistence: &UnifiedConfigurationPersistence{
			configDir:   filepath.Join(configDir, "configs"),
			compression: false,
			maxConfigs:  100,
		},
		validator: &UnifiedConfigValidator{
			rules: make(map[string]ConfigValidationRule),
		},
	}

	// Initialize schema
	manager.schema = manager.createConfigSchema()

	// Initialize validation rules
	manager.initializeValidationRules()

	return manager
}

// createConfigSchema creates the configuration schema
func (m *UnifiedConfigManager) createConfigSchema() *interfaces.ConfigSchema {
	return &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Project Configuration",
		Description: "Configuration schema for project generation",
		Required:    []string{"name", "organization", "author"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Project name",
				MinLength:   intPtr(1),
				MaxLength:   intPtr(100),
			},
			"organization": {
				Type:        "string",
				Description: "Organization name",
				MinLength:   intPtr(1),
				MaxLength:   intPtr(100),
			},
			"description": {
				Type:        "string",
				Description: "Project description",
				MaxLength:   intPtr(500),
			},
			"author": {
				Type:        "string",
				Description: "Author name",
				MinLength:   intPtr(1),
				MaxLength:   intPtr(100),
			},
			"email": {
				Type:        "string",
				Description: "Author email",
				Pattern:     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
			},
			"repository": {
				Type:        "string",
				Description: "Repository URL",
				Pattern:     "^https?://.*",
			},
			"license": {
				Type:        "string",
				Description: "License type",
				Enum:        []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "Unlicense"},
			},
			"output_path": {
				Type:        "string",
				Description: "Output directory path",
			},
		},
	}
}

// initializeValidationRules sets up configuration validation rules
func (m *UnifiedConfigManager) initializeValidationRules() {
	m.validator.rules = map[string]ConfigValidationRule{
		"required": {
			Name: "required",
			Validator: func(value interface{}) error {
				if value == nil {
					return fmt.Errorf("field is required")
				}
				if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
					return fmt.Errorf("field cannot be empty")
				}
				return nil
			},
			Message:     "This field is required",
			Suggestions: []string{"Please provide a value for this field"},
			Required:    true,
		},
		"project_name": {
			Name: "project_name",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					// Project name validation: alphanumeric, hyphens, underscores
					if !isValidProjectName(str) {
						return fmt.Errorf("invalid project name format")
					}
				}
				return nil
			},
			Message:     "Project name must contain only alphanumeric characters, hyphens, and underscores",
			Suggestions: []string{"my-awesome-project", "project_name"},
			Required:    false,
		},
		"package_name": {
			Name: "package_name",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					// Package name validation: lowercase, hyphens only
					if !isValidPackageName(str) {
						return fmt.Errorf("invalid package name format")
					}
				}
				return nil
			},
			Message:     "Package name must be lowercase with hyphens only",
			Suggestions: []string{"my-package", "awesome-package"},
			Required:    false,
		},
		"email": {
			Name: "email",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !isValidEmail(str) {
						return fmt.Errorf("invalid email format")
					}
				}
				return nil
			},
			Message:     "Please enter a valid email address",
			Suggestions: []string{"example@domain.com"},
			Required:    false,
		},
		"url": {
			Name: "url",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !isValidURL(str) {
						return fmt.Errorf("invalid URL format")
					}
				}
				return nil
			},
			Message:     "Please enter a valid URL",
			Suggestions: []string{"https://example.com"},
			Required:    false,
		},
		"license": {
			Name: "license",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !isValidLicense(str) {
						return fmt.Errorf("invalid license type")
					}
				}
				return nil
			},
			Message:     "Please select a valid license type",
			Suggestions: []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"},
			Required:    false,
		},
		"path": {
			Name: "path",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !isValidPath(str) {
						return fmt.Errorf("invalid path format")
					}
				}
				return nil
			},
			Message:     "Please enter a valid path",
			Suggestions: []string{"/path/to/directory", "./relative/path"},
			Required:    false,
		},
	}
}

// LoadDefaults loads default configuration values
func (m *UnifiedConfigManager) LoadDefaults() (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{
		Name:             "my-project",
		Organization:     "my-organization",
		Author:           "My Name",
		Description:      "A generated project",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
		License:          "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"npm": "10.0.0",
			},
		},
	}

	// Load from defaults file if it exists
	if m.defaultsPath != "" {
		if err := m.LoadFromFile(m.defaultsPath, config); err != nil {
			// Log warning but continue with empty config
			fmt.Printf("Warning: Could not load defaults from %s: %v\n", m.defaultsPath, err)
		}
	}

	// Load from environment variables
	envConfig := m.LoadFromEnvironment()
	config = m.MergeConfigurations(config, envConfig)

	return config, nil
}

// LoadFromFile loads configuration from a file
func (m *UnifiedConfigManager) LoadFromFile(path string, config *models.ProjectConfig) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Validate and sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid path: %s", path)
	}
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Determine file type and parse accordingly
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, config)
	case ".json":
		err = json.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

// SaveToFile saves configuration to a file
func (m *UnifiedConfigManager) SaveToFile(config *models.ProjectConfig, path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Determine file type and marshal accordingly
	ext := strings.ToLower(filepath.Ext(path))
	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// LoadFromEnvironment loads configuration from environment variables
func (m *UnifiedConfigManager) LoadFromEnvironment() *models.ProjectConfig {
	config := &models.ProjectConfig{}

	// Map environment variables to config fields
	envMap := map[string]*string{
		"GENERATOR_NAME":         &config.Name,
		"GENERATOR_ORGANIZATION": &config.Organization,
		"GENERATOR_DESCRIPTION":  &config.Description,
		"GENERATOR_AUTHOR":       &config.Author,
		"GENERATOR_EMAIL":        &config.Email,
		"GENERATOR_REPOSITORY":   &config.Repository,
		"GENERATOR_LICENSE":      &config.License,
		"GENERATOR_OUTPUT_PATH":  &config.OutputPath,
	}

	for envVar, field := range envMap {
		if value := os.Getenv(envVar); value != "" {
			*field = value
		}
	}

	return config
}

// MergeConfigurations merges multiple configurations with precedence
func (m *UnifiedConfigManager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
	result := &models.ProjectConfig{}

	for _, config := range configs {
		if config == nil {
			continue
		}

		// Merge fields with non-empty values taking precedence
		if config.Name != "" {
			result.Name = config.Name
		}
		if config.Organization != "" {
			result.Organization = config.Organization
		}
		if config.Description != "" {
			result.Description = config.Description
		}
		if config.Author != "" {
			result.Author = config.Author
		}
		if config.Email != "" {
			result.Email = config.Email
		}
		if config.Repository != "" {
			result.Repository = config.Repository
		}
		if config.License != "" {
			result.License = config.License
		}
		if config.OutputPath != "" {
			result.OutputPath = config.OutputPath
		}

		// Merge components if they exist
		if config.Components != (models.Components{}) {
			result.Components = config.Components
		}

		// Merge versions if they exist
		if config.Versions != nil {
			result.Versions = config.Versions
		}

		// Merge features if they exist
		if len(config.Features) > 0 {
			result.Features = config.Features
		}
	}

	return result
}

// ValidateConfig validates a configuration using unified validation rules
func (m *UnifiedConfigManager) ValidateConfig(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary:  interfaces.ConfigValidationSummary{},
	}

	if config == nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "config",
			Type:     "null",
			Message:  "Configuration cannot be null",
			Severity: "error",
			Rule:     "required",
		})
		return result, nil
	}

	// Validate each field according to schema
	for fieldName, fieldSchema := range m.schema.Properties {
		value := m.getFieldValue(config, fieldName)

		// Check if required field is missing
		isRequired := false
		for _, reqField := range m.schema.Required {
			if reqField == fieldName {
				isRequired = true
				break
			}
		}

		if isRequired && (value == nil || value == "") {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    fieldName,
				Type:     "missing",
				Message:  fmt.Sprintf("%s is required", fieldName),
				Severity: "error",
				Rule:     "required",
			})
			continue
		}

		// Skip validation for empty optional fields
		if !isRequired && (value == nil || value == "") {
			continue
		}

		// Apply validation rules based on field schema
		if strValue, ok := value.(string); ok {
			// Check min/max length
			if fieldSchema.MinLength != nil && len(strValue) < *fieldSchema.MinLength {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    fieldName,
					Type:     "validation",
					Message:  fmt.Sprintf("%s must be at least %d characters long", fieldName, *fieldSchema.MinLength),
					Severity: "error",
					Rule:     "min_length",
				})
			}

			if fieldSchema.MaxLength != nil && len(strValue) > *fieldSchema.MaxLength {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    fieldName,
					Type:     "validation",
					Message:  fmt.Sprintf("%s must be no more than %d characters long", fieldName, *fieldSchema.MaxLength),
					Severity: "error",
					Rule:     "max_length",
				})
			}

			// Check pattern
			if fieldSchema.Pattern != "" {
				if rule, exists := m.validator.rules["pattern"]; exists {
					if err := rule.Validator(value); err != nil {
						error := interfaces.ConfigValidationError{
							Field:    fieldName,
							Type:     "validation",
							Message:  rule.Message,
							Severity: "error",
							Rule:     "pattern",
						}
						result.Errors = append(result.Errors, error)
					}
				}
			}

			// Check enum values
			if len(fieldSchema.Enum) > 0 {
				valid := false
				for _, enumValue := range fieldSchema.Enum {
					if strValue == enumValue {
						valid = true
						break
					}
				}
				if !valid {
					result.Errors = append(result.Errors, interfaces.ConfigValidationError{
						Field:    fieldName,
						Type:     "validation",
						Message:  fmt.Sprintf("%s must be one of: %s", fieldName, strings.Join(fieldSchema.Enum, ", ")),
						Severity: "error",
						Rule:     "enum",
					})
				}
			}
		}
	}

	// Calculate summary
	result.Summary.ErrorCount = len(result.Errors)
	result.Summary.WarningCount = len(result.Warnings)
	result.Valid = result.Summary.ErrorCount == 0

	return result, nil
}

// getFieldValue gets the value of a field from a config struct using reflection
func (m *UnifiedConfigManager) getFieldValue(config *models.ProjectConfig, fieldName string) interface{} {
	v := reflect.ValueOf(config).Elem()
	t := reflect.TypeOf(config).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, fieldName) {
			return v.Field(i).Interface()
		}
	}

	return nil
}

// GetSetting gets a setting value
func (m *UnifiedConfigManager) GetSetting(key string) (interface{}, error) {
	if value, exists := m.settings[key]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("setting %s not found", key)
}

// SetSetting sets a setting value
func (m *UnifiedConfigManager) SetSetting(key string, value interface{}) error {
	m.settings[key] = value
	return nil
}

// GetConfigLocation returns the configuration directory
func (m *UnifiedConfigManager) GetConfigLocation() string {
	return m.configDir
}

// CreateDefaultConfig creates a default configuration file
func (m *UnifiedConfigManager) CreateDefaultConfig(path string) error {
	defaultConfig, err := m.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %v", err)
	}

	return m.SaveToFile(defaultConfig, path)
}

// Helper functions for validation

func isValidProjectName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Check for valid characters: alphanumeric, hyphens, underscores
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return false
		}
	}

	// Must start and end with alphanumeric
	first := name[0]
	last := name[len(name)-1]
	return ((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || (first >= '0' && first <= '9')) &&
		((last >= 'a' && last <= 'z') || (last >= 'A' && last <= 'Z') || (last >= '0' && last <= '9'))
}

func isValidPackageName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Check for valid characters: lowercase letters, numbers, hyphens
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}

	// Must start and end with alphanumeric
	first := name[0]
	last := name[len(name)-1]
	return ((first >= 'a' && first <= 'z') || (first >= '0' && first <= '9')) &&
		((last >= 'a' && last <= 'z') || (last >= '0' && last <= '9'))
}

func isValidEmail(email string) bool {
	if len(email) == 0 {
		return true // Empty email is allowed
	}

	// Basic email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	return true
}

func isValidURL(url string) bool {
	if len(url) == 0 {
		return true // Empty URL is allowed
	}

	// Basic URL validation
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func isValidLicense(license string) bool {
	validLicenses := []string{
		"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "Unlicense",
		"LGPL-3.0", "AGPL-3.0", "MPL-2.0", "EPL-2.0", "CC0-1.0",
	}

	for _, valid := range validLicenses {
		if license == valid {
			return true
		}
	}

	return false
}

func isValidPath(path string) bool {
	if len(path) == 0 {
		return true // Empty path is allowed
	}

	// Basic path validation
	return !strings.Contains(path, "..") && !strings.Contains(path, "//")
}
