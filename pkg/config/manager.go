// Package config provides comprehensive configuration management functionality.
//
// This package implements configuration export, import, validation, transformation,
// and security features for the Open Source Project Generator.
package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Manager implements comprehensive configuration management
type Manager struct {
	configDir     string
	encryptionKey []byte
	logger        interfaces.Logger
	validator     *ConfigValidator
	transformer   *ConfigTransformer
	encryptor     *ConfigEncryptor
	persistence   *ConfigurationPersistence
	auditLogger   *ConfigAuditLogger
}

// ConfigValidator provides configuration validation
type ConfigValidator struct {
	schema *interfaces.ConfigSchema
}

// ValidateConfiguration validates a saved configuration
func (v *ConfigValidator) ValidateConfiguration(config *SavedConfiguration) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate basic metadata
	if config.Name == "" {
		return fmt.Errorf("configuration name is required")
	}

	if config.Version == "" {
		return fmt.Errorf("configuration version is required")
	}

	// Validate project configuration if present
	if config.ProjectConfig != nil {
		result, err := v.ValidateProjectConfig(config.ProjectConfig)
		if err != nil {
			return fmt.Errorf("failed to validate project config: %w", err)
		}
		if !result.Valid {
			var errors []string
			for _, validationError := range result.Errors {
				errors = append(errors, validationError.Message)
			}
			return fmt.Errorf("project configuration validation failed: %s", strings.Join(errors, "; "))
		}
	}

	// Validate template selections
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

// ValidateProjectConfig validates a project configuration
func (v *ConfigValidator) ValidateProjectConfig(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	if config == nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:    "config",
					Type:     "null",
					Message:  "configuration cannot be null",
					Severity: "error",
					Rule:     "required",
				},
			},
			Summary: interfaces.ConfigValidationSummary{
				ErrorCount: 1,
			},
		}, nil
	}

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary:  interfaces.ConfigValidationSummary{},
	}

	// Validate required fields
	if config.Name == "" {
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:      "name",
			Type:       "required",
			Message:    "project name is required",
			Suggestion: "provide a project name",
			Severity:   "error",
			Rule:       "required",
		})
	}

	if config.Organization == "" {
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:      "organization",
			Type:       "required",
			Message:    "organization is required",
			Suggestion: "provide an organization name",
			Severity:   "error",
			Rule:       "required",
		})
	}

	// Validate field formats
	if config.Name != "" {
		if !isValidProjectName(config.Name) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "name",
				Value:      config.Name,
				Type:       "pattern",
				Message:    "project name contains invalid characters",
				Suggestion: "use only letters, numbers, hyphens, and underscores",
				Severity:   "error",
				Rule:       "pattern",
			})
		}
	}

	if config.Email != "" {
		if !isValidEmail(config.Email) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "email",
				Value:      config.Email,
				Type:       "pattern",
				Message:    "invalid email format",
				Suggestion: "provide a valid email address",
				Severity:   "error",
				Rule:       "pattern",
			})
		}
	}

	// Calculate summary
	result.Summary.ErrorCount = len(result.Errors)
	result.Summary.WarningCount = len(result.Warnings)
	result.Summary.TotalProperties = 8 // Basic count
	result.Summary.ValidProperties = result.Summary.TotalProperties - result.Summary.ErrorCount

	// Set overall validity
	result.Valid = result.Summary.ErrorCount == 0

	return result, nil
}

// ConfigTransformer handles configuration transformations
type ConfigTransformer struct {
	logger interfaces.Logger
}

// ConfigEncryptor handles configuration encryption
type ConfigEncryptor struct {
	key    []byte
	logger interfaces.Logger
}

// ConfigAuditLogger handles configuration audit logging
type ConfigAuditLogger struct {
	logFile string
	logger  interfaces.Logger
}

// ConfigExportOptions defines options for configuration export
type ConfigExportOptions struct {
	Format         string   `json:"format"`          // yaml, json
	IncludeMeta    bool     `json:"include_meta"`    // include metadata
	EncryptSecrets bool     `json:"encrypt_secrets"` // encrypt sensitive data
	ExcludeFields  []string `json:"exclude_fields"`  // fields to exclude
	Minify         bool     `json:"minify"`          // minify output
}

// ConfigImportOptions defines options for configuration import
type ConfigImportOptions struct {
	Format         string            `json:"format"`          // yaml, json, auto
	ValidateSchema bool              `json:"validate_schema"` // validate against schema
	MergeStrategy  string            `json:"merge_strategy"`  // replace, merge, append
	FieldMappings  map[string]string `json:"field_mappings"`  // field name mappings
	Transform      bool              `json:"transform"`       // apply transformations
}

// ConfigListOptions defines options for listing configurations
type ConfigListOptions struct {
	SortBy      string   `json:"sort_by"`      // name, created_at, updated_at, size
	SortOrder   string   `json:"sort_order"`   // asc, desc
	FilterTags  []string `json:"filter_tags"`  // filter by tags
	SearchQuery string   `json:"search_query"` // search in name and description
	Format      string   `json:"format"`       // table, json, yaml
	Limit       int      `json:"limit"`        // maximum number of results
	Offset      int      `json:"offset"`       // pagination offset
}

// ConfigInfo contains basic information about a configuration
type ConfigInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     string    `json:"version"`
	Tags        []string  `json:"tags,omitempty"`
	Size        int64     `json:"size"`
	Encrypted   bool      `json:"encrypted"`
	Format      string    `json:"format"`
}

// NewManager creates a new configuration manager
func NewManager(configDir string, logger interfaces.Logger) (interfaces.ConfigManager, error) {
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".generator", "configs")
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate or load encryption key
	encryptionKey, err := getOrCreateEncryptionKey(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption: %w", err)
	}

	// Create schema
	schema := createConfigSchema()

	// Initialize components
	validator := &ConfigValidator{schema: schema}
	transformer := &ConfigTransformer{logger: logger}
	encryptor := &ConfigEncryptor{key: encryptionKey, logger: logger}
	persistence := NewConfigurationPersistence(configDir, logger)
	auditLogger := &ConfigAuditLogger{
		logFile: filepath.Join(configDir, "audit.log"),
		logger:  logger,
	}

	manager := &Manager{
		configDir:     configDir,
		encryptionKey: encryptionKey,
		logger:        logger,
		validator:     validator,
		transformer:   transformer,
		encryptor:     encryptor,
		persistence:   persistence,
		auditLogger:   auditLogger,
	}

	return manager, nil
}

// ExportConfiguration exports a configuration to a file or returns data
func (m *Manager) ExportConfiguration(name string, options *ConfigExportOptions) ([]byte, error) {
	if options == nil {
		options = &ConfigExportOptions{
			Format:         "yaml",
			IncludeMeta:    true,
			EncryptSecrets: false,
		}
	}

	// Load configuration
	config, err := m.persistence.LoadConfiguration(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Transform configuration for export
	exportData, err := m.transformer.TransformForExport(config, options)
	if err != nil {
		return nil, fmt.Errorf("failed to transform configuration: %w", err)
	}

	// Encrypt sensitive data if requested
	if options.EncryptSecrets {
		exportData, err = m.encryptor.EncryptSensitiveFields(exportData)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt sensitive data: %w", err)
		}
	}

	// Serialize to requested format
	var data []byte
	switch strings.ToLower(options.Format) {
	case "json":
		if options.Minify {
			data, err = json.Marshal(exportData)
		} else {
			data, err = json.MarshalIndent(exportData, "", "  ")
		}
	case "yaml", "yml":
		data, err = yaml.Marshal(exportData)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", options.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to serialize configuration: %w", err)
	}

	// Log export action
	m.auditLogger.LogAction("export", name, map[string]interface{}{
		"format":          options.Format,
		"encrypt_secrets": options.EncryptSecrets,
		"size":            len(data),
	})

	return data, nil
}

// ImportConfiguration imports a configuration from data
func (m *Manager) ImportConfiguration(name string, data []byte, options *ConfigImportOptions) error {
	if options == nil {
		options = &ConfigImportOptions{
			Format:         "auto",
			ValidateSchema: true,
			MergeStrategy:  "replace",
			Transform:      true,
		}
	}

	// Detect format if auto
	format := options.Format
	if format == "auto" {
		format = m.detectFormat(data)
	}

	// Parse configuration data
	var importData interface{}
	var err error

	switch strings.ToLower(format) {
	case "json":
		err = json.Unmarshal(data, &importData)
	case "yaml", "yml":
		err = yaml.Unmarshal(data, &importData)
	default:
		return fmt.Errorf("unsupported import format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to parse configuration data: %w", err)
	}

	// Decrypt sensitive data if encrypted
	importData, err = m.encryptor.DecryptSensitiveFields(importData)
	if err != nil {
		return fmt.Errorf("failed to decrypt sensitive data: %w", err)
	}

	// Transform configuration for import
	config, err := m.transformer.TransformForImport(importData, options)
	if err != nil {
		return fmt.Errorf("failed to transform configuration: %w", err)
	}

	// Validate configuration if requested
	if options.ValidateSchema {
		if err := m.validator.ValidateConfiguration(config); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	// Handle merge strategy
	if options.MergeStrategy != "replace" && m.persistence.ConfigurationExists(name) {
		existing, err := m.persistence.LoadConfiguration(name)
		if err != nil {
			return fmt.Errorf("failed to load existing configuration: %w", err)
		}

		config, err = m.mergeConfigurations(existing, config, options.MergeStrategy)
		if err != nil {
			return fmt.Errorf("failed to merge configurations: %w", err)
		}
	}

	// Save configuration
	if err := m.persistence.SaveConfiguration(name, config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	// Log import action
	m.auditLogger.LogAction("import", name, map[string]interface{}{
		"format":          format,
		"merge_strategy":  options.MergeStrategy,
		"validate_schema": options.ValidateSchema,
		"size":            len(data),
	})

	return nil
}

// ListConfigurations lists all configurations with filtering and sorting
func (m *Manager) ListConfigurations(options *ConfigListOptions) ([]*ConfigInfo, error) {
	if options == nil {
		options = &ConfigListOptions{
			SortBy:    "updated_at",
			SortOrder: "desc",
			Format:    "table",
			Limit:     50,
		}
	}

	// Get configurations from persistence layer
	persistenceOptions := &ConfigurationListOptions{
		SortBy:      options.SortBy,
		SortOrder:   options.SortOrder,
		FilterTags:  options.FilterTags,
		SearchQuery: options.SearchQuery,
		Limit:       options.Limit,
		Offset:      options.Offset,
	}

	configs, err := m.persistence.ListConfigurations(persistenceOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}

	// Convert to ConfigInfo
	var configInfos []*ConfigInfo
	for _, config := range configs {
		info := &ConfigInfo{
			Name:        config.Name,
			Description: config.Description,
			CreatedAt:   config.CreatedAt,
			UpdatedAt:   config.UpdatedAt,
			Version:     config.Version,
			Tags:        config.Tags,
			Format:      "yaml", // Default format
		}

		// Get file size
		filename := m.persistence.generateConfigFilename(config.Name)
		filePath := filepath.Join(m.configDir, filename)
		if fileInfo, err := os.Stat(filePath); err == nil {
			info.Size = fileInfo.Size()
		}

		// Check if encrypted (simplified check)
		info.Encrypted = m.isConfigurationEncrypted(config)

		configInfos = append(configInfos, info)
	}

	return configInfos, nil
}

// ValidateConfiguration validates a configuration
func (m *Manager) ValidateConfiguration(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	return m.validator.ValidateProjectConfig(config)
}

// DeleteConfiguration deletes a configuration
func (m *Manager) DeleteConfiguration(name string) error {
	// Log deletion action before deleting
	m.auditLogger.LogAction("delete", name, map[string]interface{}{
		"timestamp": time.Now(),
	})

	return m.persistence.DeleteConfiguration(name)
}

// GetConfiguration retrieves a configuration
func (m *Manager) GetConfiguration(name string) (*SavedConfiguration, error) {
	return m.persistence.LoadConfiguration(name)
}

// SaveConfiguration saves a configuration
func (m *Manager) SaveConfiguration(name string, config *SavedConfiguration) error {
	// Validate configuration
	if err := m.validator.ValidateConfiguration(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Log save action
	m.auditLogger.LogAction("save", name, map[string]interface{}{
		"version": config.Version,
	})

	return m.persistence.SaveConfiguration(name, config)
}

// GetConfigurationInfo gets basic information about a configuration
func (m *Manager) GetConfigurationInfo(name string) (*ConfigInfo, error) {
	info, err := m.persistence.GetConfigurationInfo(name)
	if err != nil {
		return nil, err
	}

	// Convert to our ConfigInfo type
	return &ConfigInfo{
		Name:        info.Name,
		Description: info.Description,
		CreatedAt:   info.CreatedAt,
		UpdatedAt:   info.UpdatedAt,
		Version:     info.Version,
		Tags:        info.Tags,
		Size:        info.FileSize,
		Format:      "yaml",
		Encrypted:   false, // TODO: Implement encryption detection
	}, nil
}

// Helper methods

// detectFormat detects the format of configuration data
func (m *Manager) detectFormat(data []byte) string {
	// Try to parse as JSON first
	var jsonData interface{}
	if json.Unmarshal(data, &jsonData) == nil {
		return "json"
	}

	// Try to parse as YAML
	var yamlData interface{}
	if yaml.Unmarshal(data, &yamlData) == nil {
		return "yaml"
	}

	// Default to YAML
	return "yaml"
}

// isConfigurationEncrypted checks if a configuration contains encrypted data
func (m *Manager) isConfigurationEncrypted(config *SavedConfiguration) bool {
	// Simple check for encrypted fields
	// In a real implementation, this would check for encrypted field markers
	return false
}

// mergeConfigurations merges two configurations based on strategy
func (m *Manager) mergeConfigurations(existing, new *SavedConfiguration, strategy string) (*SavedConfiguration, error) {
	switch strategy {
	case "merge":
		// Merge configurations, new values override existing
		merged := *existing
		if new.Description != "" {
			merged.Description = new.Description
		}
		if len(new.Tags) > 0 {
			merged.Tags = append(merged.Tags, new.Tags...)
		}
		if new.ProjectConfig != nil {
			merged.ProjectConfig = new.ProjectConfig
		}
		if len(new.SelectedTemplates) > 0 {
			merged.SelectedTemplates = new.SelectedTemplates
		}
		if new.GenerationSettings != nil {
			merged.GenerationSettings = new.GenerationSettings
		}
		if new.UserPreferences != nil {
			merged.UserPreferences = new.UserPreferences
		}
		merged.UpdatedAt = time.Now()
		return &merged, nil

	case "append":
		// Append new data to existing
		merged := *existing
		merged.Tags = append(merged.Tags, new.Tags...)
		merged.SelectedTemplates = append(merged.SelectedTemplates, new.SelectedTemplates...)
		merged.UpdatedAt = time.Now()
		return &merged, nil

	case "replace":
		fallthrough
	default:
		// Replace existing with new
		new.CreatedAt = existing.CreatedAt
		new.UpdatedAt = time.Now()
		return new, nil
	}
}

// getOrCreateEncryptionKey gets or creates an encryption key
func getOrCreateEncryptionKey(configDir string) ([]byte, error) {
	keyFile := filepath.Join(configDir, ".encryption_key")

	// Try to read existing key
	if data, err := os.ReadFile(keyFile); err == nil {
		key, err := base64.StdEncoding.DecodeString(string(data))
		if err == nil && len(key) == 32 {
			return key, nil
		}
	}

	// Generate new key
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Save key
	keyData := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(keyFile, []byte(keyData), 0600); err != nil {
		return nil, fmt.Errorf("failed to save encryption key: %w", err)
	}

	return key, nil
}

// createConfigSchema creates the configuration schema
func createConfigSchema() *interfaces.ConfigSchema {
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
			},
			"organization": {
				Type:        "string",
				Description: "Organization or author name",
				Required:    true,
				MinLength:   &[]int{1}[0],
				MaxLength:   &[]int{100}[0],
			},
			"description": {
				Type:        "string",
				Description: "Project description",
				Required:    false,
				MaxLength:   &[]int{500}[0],
			},
			"license": {
				Type:        "string",
				Description: "Project license",
				Default:     "MIT",
				Enum:        []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"},
			},
		},
		Required: []string{"name", "organization"},
	}
}

// Helper functions for validation

// isValidProjectName checks if a project name is valid
func isValidProjectName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Allow letters, numbers, hyphens, and underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}

	return true
}

// isValidEmail checks if an email address is valid
func isValidEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	// Simple email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]
	if len(local) == 0 || len(local) > 64 || len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Check for basic domain structure
	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

// Interface compliance check
var _ interfaces.ConfigManager = (*Manager)(nil)

// Additional methods to satisfy the ConfigManager interface

// LoadDefaults loads default configuration values
func (m *Manager) LoadDefaults() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         "my-project",
		Organization: "my-org",
		License:      "MIT",
		OutputPath:   "./output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Shared: true,
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
				"react":      "18.2.0",
				"next":       "13.4.0",
				"typescript": "5.0.0",
			},
		},
	}, nil
}

// ValidateConfig validates a project configuration
func (m *Manager) ValidateConfig(config *models.ProjectConfig) error {
	result, err := m.validator.ValidateProjectConfig(config)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid {
		var errorMessages []string
		for _, validationError := range result.Errors {
			errorMessages = append(errorMessages, validationError.Message)
		}
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

// SaveConfig saves configuration to a file
func (m *Manager) SaveConfig(config *models.ProjectConfig, path string) error {
	// Validate configuration first
	if err := m.ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Determine format from file extension
	ext := strings.ToLower(filepath.Ext(path))
	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
	default:
		return fmt.Errorf("unsupported file format: %s (use .json, .yaml, or .yml)", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file safely
	if err := utils.SafeWriteFile(path, data); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	// Log the action
	m.auditLogger.LogAction("save_config", path, map[string]interface{}{
		"format": ext,
		"size":   len(data),
	})

	return nil
}

// LoadConfig loads configuration from a file
func (m *Manager) LoadConfig(path string) (*models.ProjectConfig, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	content, err := utils.SafeReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config models.ProjectConfig
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		if err := json.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON configuration: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML configuration: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s (use .json, .yaml, or .yml)", ext)
	}

	// Log the action
	m.auditLogger.LogAction("load_config", path, map[string]interface{}{
		"format": ext,
		"size":   len(content),
	})

	return &config, nil
}

// GetSetting gets a configuration setting (placeholder implementation)
func (m *Manager) GetSetting(key string) (any, error) {
	return nil, fmt.Errorf("setting '%s' not found", key)
}

// SetSetting sets a configuration setting (placeholder implementation)
func (m *Manager) SetSetting(key string, value any) error {
	return nil
}

// ValidateSettings validates configuration settings (placeholder implementation)
func (m *Manager) ValidateSettings() error {
	return nil
}

// LoadFromFile loads configuration from a file
func (m *Manager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

// LoadFromEnvironment loads configuration from environment variables (placeholder implementation)
func (m *Manager) LoadFromEnvironment() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{}, nil
}

// MergeConfigurations merges multiple configurations
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

		// Merge basic fields (later configs override earlier ones)
		if configs[i].Name != "" {
			result.Name = configs[i].Name
		}
		if configs[i].Organization != "" {
			result.Organization = configs[i].Organization
		}
		if configs[i].Description != "" {
			result.Description = configs[i].Description
		}
		if configs[i].License != "" {
			result.License = configs[i].License
		}
		if configs[i].Author != "" {
			result.Author = configs[i].Author
		}
		if configs[i].Email != "" {
			result.Email = configs[i].Email
		}
		if configs[i].Repository != "" {
			result.Repository = configs[i].Repository
		}
		if configs[i].OutputPath != "" {
			result.OutputPath = configs[i].OutputPath
		}

		// Merge features
		if len(configs[i].Features) > 0 {
			result.Features = append(result.Features, configs[i].Features...)
		}

		// Merge versions
		if configs[i].Versions != nil {
			if result.Versions == nil {
				result.Versions = &models.VersionConfig{Packages: make(map[string]string)}
			}
			if configs[i].Versions.Node != "" {
				result.Versions.Node = configs[i].Versions.Node
			}
			if configs[i].Versions.Go != "" {
				result.Versions.Go = configs[i].Versions.Go
			}
			for pkg, version := range configs[i].Versions.Packages {
				if version != "" {
					result.Versions.Packages[pkg] = version
				}
			}
		}
	}

	return &result
}

// GetConfigSchema returns the configuration schema
func (m *Manager) GetConfigSchema() *interfaces.ConfigSchema {
	return m.validator.schema
}

// ValidateConfigFromFile validates a configuration file
func (m *Manager) ValidateConfigFromFile(path string) (*interfaces.ConfigValidationResult, error) {
	config, err := m.LoadConfig(path)
	if err != nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:    "file",
					Value:    path,
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

	return m.validator.ValidateProjectConfig(config)
}

// GetConfigSources returns configuration sources (placeholder implementation)
func (m *Manager) GetConfigSources() ([]interfaces.ConfigSource, error) {
	return []interfaces.ConfigSource{
		{
			Type:     "file",
			Location: m.configDir,
			Priority: 1,
			Valid:    true,
		},
	}, nil
}

// GetConfigLocation returns the configuration location
func (m *Manager) GetConfigLocation() string {
	return m.configDir
}

// CreateDefaultConfig creates a default configuration file
func (m *Manager) CreateDefaultConfig(path string) error {
	defaultConfig, err := m.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	return m.SaveConfig(defaultConfig, path)
}

// BackupConfig backs up a configuration file (placeholder implementation)
func (m *Manager) BackupConfig(path string) error {
	return fmt.Errorf("backup functionality not implemented")
}

// RestoreConfig restores a configuration from backup (placeholder implementation)
func (m *Manager) RestoreConfig(backupPath string) error {
	return fmt.Errorf("restore functionality not implemented")
}

// LoadEnvironmentVariables loads environment variables (placeholder implementation)
func (m *Manager) LoadEnvironmentVariables() map[string]string {
	return make(map[string]string)
}

// SetEnvironmentDefaults sets environment defaults (placeholder implementation)
func (m *Manager) SetEnvironmentDefaults() error {
	return nil
}

// GetEnvironmentPrefix returns the environment prefix
func (m *Manager) GetEnvironmentPrefix() string {
	return "GENERATOR"
}
