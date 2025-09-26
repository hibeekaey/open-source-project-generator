package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MockConfigManager implements interfaces.ConfigManager for testing
type MockConfigManager struct {
	configs  map[string]*models.ProjectConfig
	settings map[string]any
	schema   *interfaces.ConfigSchema
}

// NewMockConfigManager creates a new mock config manager
func NewMockConfigManager() *MockConfigManager {
	return &MockConfigManager{
		configs:  make(map[string]*models.ProjectConfig),
		settings: make(map[string]any),
		schema: &interfaces.ConfigSchema{
			Properties: map[string]interfaces.PropertySchema{
				"name": {
					Type:        "string",
					Description: "Project name",
					Required:    true,
				},
				"organization": {
					Type:        "string",
					Description: "Organization name",
					Required:    false,
				},
			},
			Required: []string{"name"},
		},
	}
}

// LoadDefaults loads default configuration
func (m *MockConfigManager) LoadDefaults() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         "default-project",
		Organization: "default-org",
		Description:  "Default project description",
		License:      "MIT",
	}, nil
}

// ValidateConfig validates project configuration
func (m *MockConfigManager) ValidateConfig(config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

// SaveConfig saves configuration to file
func (m *MockConfigManager) SaveConfig(config *models.ProjectConfig, path string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	m.configs[path] = config
	return nil
}

// LoadConfig loads configuration from file
func (m *MockConfigManager) LoadConfig(path string) (*models.ProjectConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	config, exists := m.configs[path]
	if !exists {
		return nil, fmt.Errorf("config file not found: %s", path)
	}
	return config, nil
}

// GetSetting gets a configuration setting
func (m *MockConfigManager) GetSetting(key string) (any, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	value, exists := m.settings[key]
	if !exists {
		return nil, fmt.Errorf("setting not found: %s", key)
	}
	return value, nil
}

// SetSetting sets a configuration setting
func (m *MockConfigManager) SetSetting(key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	m.settings[key] = value
	return nil
}

// ValidateSettings validates all settings
func (m *MockConfigManager) ValidateSettings() error {
	for key, value := range m.settings {
		if key == "" {
			return fmt.Errorf("empty key found in settings")
		}
		if value == nil {
			return fmt.Errorf("nil value found for key: %s", key)
		}
	}
	return nil
}

// LoadFromFile loads configuration from file
func (m *MockConfigManager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

// LoadFromEnvironment loads configuration from environment variables
func (m *MockConfigManager) LoadFromEnvironment() (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{}

	if name := os.Getenv("GENERATOR_PROJECT_NAME"); name != "" {
		config.Name = name
	}
	if org := os.Getenv("GENERATOR_ORGANIZATION"); org != "" {
		config.Organization = org
	}
	if desc := os.Getenv("GENERATOR_DESCRIPTION"); desc != "" {
		config.Description = desc
	}
	if license := os.Getenv("GENERATOR_LICENSE"); license != "" {
		config.License = license
	}

	return config, nil
}

// MergeConfigurations merges multiple configurations
func (m *MockConfigManager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
	if len(configs) == 0 {
		return &models.ProjectConfig{}
	}

	result := &models.ProjectConfig{}

	for _, config := range configs {
		if config == nil {
			continue
		}

		if config.Name != "" {
			result.Name = config.Name
		}
		if config.Organization != "" {
			result.Organization = config.Organization
		}
		if config.Description != "" {
			result.Description = config.Description
		}
		if config.License != "" {
			result.License = config.License
		}
	}

	return result
}

// GetConfigSchema returns the configuration schema
func (m *MockConfigManager) GetConfigSchema() *interfaces.ConfigSchema {
	return m.schema
}

// ValidateConfigFromFile validates configuration from file
func (m *MockConfigManager) ValidateConfigFromFile(path string) (*interfaces.ConfigValidationResult, error) {
	config, err := m.LoadConfig(path)
	if err != nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:   "file",
					Message: err.Error(),
					Type:    "file_error",
				},
			},
		}, nil
	}

	if err := m.ValidateConfig(config); err != nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:   "config",
					Message: err.Error(),
					Type:    "validation_error",
				},
			},
		}, nil
	}

	return &interfaces.ConfigValidationResult{
		Valid:  true,
		Errors: []interfaces.ConfigValidationError{},
	}, nil
}

// GetConfigSources returns configuration sources
func (m *MockConfigManager) GetConfigSources() ([]interfaces.ConfigSource, error) {
	return []interfaces.ConfigSource{
		{
			Type:     "file",
			Location: "/path/to/config.yaml",
			Priority: 1,
		},
		{
			Type:     "environment",
			Location: "env",
			Priority: 2,
		},
	}, nil
}

// GetConfigLocation returns configuration location
func (m *MockConfigManager) GetConfigLocation() string {
	return "/path/to/config.yaml"
}

// CreateDefaultConfig creates default configuration file
func (m *MockConfigManager) CreateDefaultConfig(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	defaultConfig := &models.ProjectConfig{
		Name:         "new-project",
		Organization: "my-org",
		Description:  "A new project",
		License:      "MIT",
	}

	return m.SaveConfig(defaultConfig, path)
}

// BackupConfig creates a backup of configuration
func (m *MockConfigManager) BackupConfig(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	config, exists := m.configs[path]
	if !exists {
		return fmt.Errorf("config not found at path: %s", path)
	}

	backupPath := path + ".backup." + time.Now().Format("20060102150405")
	return m.SaveConfig(config, backupPath)
}

// RestoreConfig restores configuration from backup
func (m *MockConfigManager) RestoreConfig(backupPath string) error {
	if backupPath == "" {
		return fmt.Errorf("backup path cannot be empty")
	}

	config, err := m.LoadConfig(backupPath)
	if err != nil {
		return err
	}

	// For simplicity, restore to the original path (remove .backup.timestamp)
	originalPath := strings.Split(backupPath, ".backup.")[0]
	m.configs[originalPath] = config
	return nil
}

// LoadEnvironmentVariables loads environment variables
func (m *MockConfigManager) LoadEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)

	if name := os.Getenv("GENERATOR_PROJECT_NAME"); name != "" {
		envVars["GENERATOR_PROJECT_NAME"] = name
	}
	if org := os.Getenv("GENERATOR_ORGANIZATION"); org != "" {
		envVars["GENERATOR_ORGANIZATION"] = org
	}

	return envVars
}

// SetEnvironmentDefaults sets environment defaults
func (m *MockConfigManager) SetEnvironmentDefaults() error {
	defaults := map[string]string{
		"GENERATOR_LICENSE": "MIT",
		"GENERATOR_OUTPUT":  "./output",
	}

	for key, value := range defaults {
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

// GetEnvironmentPrefix returns environment prefix
func (m *MockConfigManager) GetEnvironmentPrefix() string {
	return "GENERATOR"
}

// Test functions

func TestMockConfigManager_LoadDefaults(t *testing.T) {
	manager := NewMockConfigManager()

	config, err := manager.LoadDefaults()
	if err != nil {
		t.Fatalf("LoadDefaults failed: %v", err)
	}

	if config.Name != "default-project" {
		t.Errorf("Expected name 'default-project', got '%s'", config.Name)
	}

	if config.License != "MIT" {
		t.Errorf("Expected license 'MIT', got '%s'", config.License)
	}
}

func TestMockConfigManager_ValidateConfig(t *testing.T) {
	manager := NewMockConfigManager()

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				License:      "MIT",
			},
			expectError: false,
		},
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name: "empty name",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "test-org",
				License:      "MIT",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateConfig(tt.config)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestMockConfigManager_SaveAndLoadConfig(t *testing.T) {
	manager := NewMockConfigManager()

	config := &models.ProjectConfig{
		Name:         "save-test",
		Organization: "save-org",
		Description:  "Save test description",
		License:      "Apache-2.0",
	}

	// Test save
	err := manager.SaveConfig(config, "/test/path/config.yaml")
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Test load
	loadedConfig, err := manager.LoadConfig("/test/path/config.yaml")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loadedConfig.Name != config.Name {
		t.Errorf("Expected name '%s', got '%s'", config.Name, loadedConfig.Name)
	}

	if loadedConfig.License != config.License {
		t.Errorf("Expected license '%s', got '%s'", config.License, loadedConfig.License)
	}
}

func TestMockConfigManager_Settings(t *testing.T) {
	manager := NewMockConfigManager()

	// Test set setting
	err := manager.SetSetting("test_key", "test_value")
	if err != nil {
		t.Fatalf("SetSetting failed: %v", err)
	}

	// Test get setting
	value, err := manager.GetSetting("test_key")
	if err != nil {
		t.Fatalf("GetSetting failed: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value)
	}

	// Test get non-existent setting
	_, err = manager.GetSetting("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent setting")
	}

	// Test validate settings
	err = manager.ValidateSettings()
	if err != nil {
		t.Fatalf("ValidateSettings failed: %v", err)
	}
}

func TestMockConfigManager_MergeConfigurations(t *testing.T) {
	manager := NewMockConfigManager()

	config1 := &models.ProjectConfig{
		Name:         "project1",
		Organization: "org1",
	}

	config2 := &models.ProjectConfig{
		Description: "description2",
		License:     "MIT",
	}

	config3 := &models.ProjectConfig{
		Name:    "project3",   // Should override config1
		License: "Apache-2.0", // Should override config2
	}

	merged := manager.MergeConfigurations(config1, config2, config3)

	if merged.Name != "project3" {
		t.Errorf("Expected name 'project3', got '%s'", merged.Name)
	}

	if merged.Organization != "org1" {
		t.Errorf("Expected organization 'org1', got '%s'", merged.Organization)
	}

	if merged.Description != "description2" {
		t.Errorf("Expected description 'description2', got '%s'", merged.Description)
	}

	if merged.License != "Apache-2.0" {
		t.Errorf("Expected license 'Apache-2.0', got '%s'", merged.License)
	}
}

func TestMockConfigManager_LoadFromEnvironment(t *testing.T) {
	manager := NewMockConfigManager()

	// Set environment variables
	_ = os.Setenv("GENERATOR_PROJECT_NAME", "env-project")
	_ = os.Setenv("GENERATOR_ORGANIZATION", "env-org")
	defer func() {
		_ = os.Unsetenv("GENERATOR_PROJECT_NAME")
		_ = os.Unsetenv("GENERATOR_ORGANIZATION")
	}()

	config, err := manager.LoadFromEnvironment()
	if err != nil {
		t.Fatalf("LoadFromEnvironment failed: %v", err)
	}

	if config.Name != "env-project" {
		t.Errorf("Expected name 'env-project', got '%s'", config.Name)
	}

	if config.Organization != "env-org" {
		t.Errorf("Expected organization 'env-org', got '%s'", config.Organization)
	}
}

func TestMockConfigManager_GetConfigSchema(t *testing.T) {
	manager := NewMockConfigManager()

	schema := manager.GetConfigSchema()
	if schema == nil {
		t.Fatal("Expected schema, got nil")
	}

	if len(schema.Properties) == 0 {
		t.Error("Expected properties in schema")
	}

	if len(schema.Required) == 0 {
		t.Error("Expected required fields in schema")
	}

	nameProperty, exists := schema.Properties["name"]
	if !exists {
		t.Error("Expected 'name' property in schema")
	}

	if nameProperty.Type != "string" {
		t.Errorf("Expected name type 'string', got '%s'", nameProperty.Type)
	}
}

func TestMockConfigManager_ValidateConfigFromFile(t *testing.T) {
	manager := NewMockConfigManager()

	// Set up a valid config
	validConfig := &models.ProjectConfig{
		Name:         "valid-project",
		Organization: "valid-org",
		License:      "MIT",
	}
	_ = manager.SaveConfig(validConfig, "/valid/path")

	// Test valid config
	result, err := manager.ValidateConfigFromFile("/valid/path")
	if err != nil {
		t.Fatalf("ValidateConfigFromFile failed: %v", err)
	}

	if !result.Valid {
		t.Error("Expected valid result")
	}

	if len(result.Errors) > 0 {
		t.Errorf("Expected no errors, got %d", len(result.Errors))
	}

	// Test invalid path
	result, err = manager.ValidateConfigFromFile("/invalid/path")
	if err != nil {
		t.Fatalf("ValidateConfigFromFile failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result for non-existent file")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors for non-existent file")
	}
}

func TestMockConfigManager_ConfigSources(t *testing.T) {
	manager := NewMockConfigManager()

	sources, err := manager.GetConfigSources()
	if err != nil {
		t.Fatalf("GetConfigSources failed: %v", err)
	}

	if len(sources) == 0 {
		t.Error("Expected config sources")
	}

	location := manager.GetConfigLocation()
	if location == "" {
		t.Error("Expected config location")
	}
}

func TestMockConfigManager_BackupAndRestore(t *testing.T) {
	manager := NewMockConfigManager()

	// Set up initial config
	originalConfig := &models.ProjectConfig{
		Name:         "original",
		Organization: "original-org",
		License:      "MIT",
	}
	_ = manager.SaveConfig(originalConfig, "/test/config")

	// Create backup
	err := manager.BackupConfig("/test/config")
	if err != nil {
		t.Fatalf("BackupConfig failed: %v", err)
	}

	// Modify config
	modifiedConfig := &models.ProjectConfig{
		Name:         "modified",
		Organization: "modified-org",
		License:      "Apache-2.0",
	}
	_ = manager.SaveConfig(modifiedConfig, "/test/config")

	// Restore from backup
	backupPath := "/test/config.backup.20240101120000" // Mock backup path
	_ = manager.SaveConfig(originalConfig, backupPath) // Simulate backup file

	err = manager.RestoreConfig(backupPath)
	if err != nil {
		t.Fatalf("RestoreConfig failed: %v", err)
	}

	// Verify restoration
	restoredConfig, err := manager.LoadConfig("/test/config")
	if err != nil {
		t.Fatalf("LoadConfig after restore failed: %v", err)
	}

	if restoredConfig.Name != "original" {
		t.Errorf("Expected restored name 'original', got '%s'", restoredConfig.Name)
	}
}

func TestMockConfigManager_EnvironmentVariables(t *testing.T) {
	manager := NewMockConfigManager()

	// Set test environment variables
	_ = os.Setenv("GENERATOR_PROJECT_NAME", "env-test")
	_ = os.Setenv("GENERATOR_ORGANIZATION", "env-test-org")
	defer func() {
		_ = os.Unsetenv("GENERATOR_PROJECT_NAME")
		_ = os.Unsetenv("GENERATOR_ORGANIZATION")
	}()

	envVars := manager.LoadEnvironmentVariables()

	if envVars["GENERATOR_PROJECT_NAME"] != "env-test" {
		t.Errorf("Expected GENERATOR_PROJECT_NAME 'env-test', got '%s'", envVars["GENERATOR_PROJECT_NAME"])
	}

	if envVars["GENERATOR_ORGANIZATION"] != "env-test-org" {
		t.Errorf("Expected GENERATOR_ORGANIZATION 'env-test-org', got '%s'", envVars["GENERATOR_ORGANIZATION"])
	}

	// Test environment defaults
	_ = os.Unsetenv("GENERATOR_LICENSE")
	err := manager.SetEnvironmentDefaults()
	if err != nil {
		t.Fatalf("SetEnvironmentDefaults failed: %v", err)
	}

	if os.Getenv("GENERATOR_LICENSE") != "MIT" {
		t.Errorf("Expected GENERATOR_LICENSE 'MIT', got '%s'", os.Getenv("GENERATOR_LICENSE"))
	}

	// Test environment prefix
	prefix := manager.GetEnvironmentPrefix()
	if prefix != "GENERATOR" {
		t.Errorf("Expected prefix 'GENERATOR', got '%s'", prefix)
	}
}

func TestMockConfigManager_ErrorHandling(t *testing.T) {
	manager := NewMockConfigManager()

	// Test empty key in SetSetting
	err := manager.SetSetting("", "value")
	if err == nil {
		t.Error("Expected error for empty key in SetSetting")
	}

	// Test empty key in GetSetting
	_, err = manager.GetSetting("")
	if err == nil {
		t.Error("Expected error for empty key in GetSetting")
	}

	// Test empty path in SaveConfig
	err = manager.SaveConfig(&models.ProjectConfig{Name: "test"}, "")
	if err == nil {
		t.Error("Expected error for empty path in SaveConfig")
	}

	// Test nil config in SaveConfig
	err = manager.SaveConfig(nil, "/test/path")
	if err == nil {
		t.Error("Expected error for nil config in SaveConfig")
	}

	// Test empty path in LoadConfig
	_, err = manager.LoadConfig("")
	if err == nil {
		t.Error("Expected error for empty path in LoadConfig")
	}

	// Test empty path in CreateDefaultConfig
	err = manager.CreateDefaultConfig("")
	if err == nil {
		t.Error("Expected error for empty path in CreateDefaultConfig")
	}

	// Test empty path in BackupConfig
	err = manager.BackupConfig("")
	if err == nil {
		t.Error("Expected error for empty path in BackupConfig")
	}

	// Test empty path in RestoreConfig
	err = manager.RestoreConfig("")
	if err == nil {
		t.Error("Expected error for empty path in RestoreConfig")
	}
}

// Benchmark tests
func BenchmarkMockConfigManager_LoadDefaults(b *testing.B) {
	manager := NewMockConfigManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.LoadDefaults()
		if err != nil {
			b.Fatalf("LoadDefaults failed: %v", err)
		}
	}
}

func BenchmarkMockConfigManager_ValidateConfig(b *testing.B) {
	manager := NewMockConfigManager()
	config := &models.ProjectConfig{
		Name:         "benchmark-project",
		Organization: "benchmark-org",
		License:      "MIT",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := manager.ValidateConfig(config)
		if err != nil {
			b.Fatalf("ValidateConfig failed: %v", err)
		}
	}
}

func BenchmarkMockConfigManager_MergeConfigurations(b *testing.B) {
	manager := NewMockConfigManager()

	config1 := &models.ProjectConfig{Name: "project1", Organization: "org1"}
	config2 := &models.ProjectConfig{Description: "desc2", License: "MIT"}
	config3 := &models.ProjectConfig{Name: "project3", License: "Apache-2.0"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = manager.MergeConfigurations(config1, config2, config3)
	}
}
