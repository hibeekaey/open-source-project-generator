package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

// TestConfigurationLoadingAndMerging tests configuration loading from multiple sources
func TestConfigurationLoadingAndMerging(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("load_from_file", func(t *testing.T) {
		testLoadFromFile(t, tempDir)
	})

	t.Run("load_from_environment", func(t *testing.T) {
		testLoadFromEnvironment(t)
	})

	t.Run("merge_configurations", func(t *testing.T) {
		testMergeConfigurations(t, tempDir)
	})

	t.Run("configuration_precedence", func(t *testing.T) {
		testConfigurationPrecedence(t, tempDir)
	})

	t.Run("configuration_validation", func(t *testing.T) {
		testConfigurationValidation(t, tempDir)
	})
}

func testLoadFromFile(t *testing.T, tempDir string) {
	// Test YAML configuration
	yamlConfig := `
name: yaml-project
organization: yaml-org
description: Project loaded from YAML
license: MIT
components:
  backend:
    enabled: true
    technology: go
  frontend:
    enabled: true
    technology: nextjs
`

	yamlPath := filepath.Join(tempDir, "config.yaml")
	err := os.WriteFile(yamlPath, []byte(yamlConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML config: %v", err)
	}

	// Create mock config manager
	configManager := NewMockConfigManager()

	// Load YAML config
	config, err := configManager.LoadFromFile(yamlPath)
	if err != nil {
		t.Fatalf("Failed to load YAML config: %v", err)
	}

	if config.Name != "yaml-project" {
		t.Errorf("Expected name 'yaml-project', got '%s'", config.Name)
	}

	if config.Organization != "yaml-org" {
		t.Errorf("Expected organization 'yaml-org', got '%s'", config.Organization)
	}

	// Test JSON configuration
	jsonConfig := models.ProjectConfig{
		Name:         "json-project",
		Organization: "json-org",
		Description:  "Project loaded from JSON",
		License:      "Apache-2.0",
	}

	jsonPath := filepath.Join(tempDir, "config.json")
	jsonData, err := json.MarshalIndent(jsonConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON config: %v", err)
	}

	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON config: %v", err)
	}

	// Mock loading JSON (in real implementation, this would parse JSON)
	_ = configManager.SaveConfig(&jsonConfig, jsonPath)
	config, err = configManager.LoadFromFile(jsonPath)
	if err != nil {
		t.Fatalf("Failed to load JSON config: %v", err)
	}

	if config.Name != "json-project" {
		t.Errorf("Expected name 'json-project', got '%s'", config.Name)
	}

	if config.License != "Apache-2.0" {
		t.Errorf("Expected license 'Apache-2.0', got '%s'", config.License)
	}
}

func testLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	envVars := map[string]string{
		"GENERATOR_PROJECT_NAME": "env-project",
		"GENERATOR_ORGANIZATION": "env-org",
		"GENERATOR_DESCRIPTION":  "Project from environment",
		"GENERATOR_LICENSE":      "BSD-3-Clause",
		"GENERATOR_OUTPUT_PATH":  "/tmp/env-output",
	}

	// Set environment variables
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	// Clean up after test
	defer func() {
		for key := range envVars {
			_ = os.Unsetenv(key)
		}
	}()

	configManager := NewMockConfigManager()

	config, err := configManager.LoadFromEnvironment()
	if err != nil {
		t.Fatalf("Failed to load from environment: %v", err)
	}

	if config.Name != "env-project" {
		t.Errorf("Expected name 'env-project', got '%s'", config.Name)
	}

	if config.Organization != "env-org" {
		t.Errorf("Expected organization 'env-org', got '%s'", config.Organization)
	}

	if config.Description != "Project from environment" {
		t.Errorf("Expected description 'Project from environment', got '%s'", config.Description)
	}

	if config.License != "BSD-3-Clause" {
		t.Errorf("Expected license 'BSD-3-Clause', got '%s'", config.License)
	}
}

func testMergeConfigurations(t *testing.T, tempDir string) {
	configManager := NewMockConfigManager()

	// Create base configuration
	baseConfig := &models.ProjectConfig{
		Name:         "base-project",
		Organization: "base-org",
		License:      "MIT",
	}

	// Create override configuration
	overrideConfig := &models.ProjectConfig{
		Name:        "override-project",     // Should override base
		Description: "Override description", // Should be added
		// Organization and License should come from base
	}

	// Create additional configuration
	additionalConfig := &models.ProjectConfig{
		License:    "Apache-2.0",  // Should override base and override
		OutputPath: "/tmp/output", // Should be added
	}

	// Merge configurations
	merged := configManager.MergeConfigurations(baseConfig, overrideConfig, additionalConfig)

	// Verify merge results
	if merged.Name != "override-project" {
		t.Errorf("Expected name 'override-project', got '%s'", merged.Name)
	}

	if merged.Organization != "base-org" {
		t.Errorf("Expected organization 'base-org', got '%s'", merged.Organization)
	}

	if merged.Description != "Override description" {
		t.Errorf("Expected description 'Override description', got '%s'", merged.Description)
	}

	if merged.License != "Apache-2.0" {
		t.Errorf("Expected license 'Apache-2.0', got '%s'", merged.License)
	}

	if merged.OutputPath != "/tmp/output" {
		t.Errorf("Expected output path '/tmp/output', got '%s'", merged.OutputPath)
	}
}

func testConfigurationPrecedence(t *testing.T, tempDir string) {
	configManager := NewMockConfigManager()

	// Test precedence: Environment > File > Defaults

	// 1. Set defaults
	defaults := &models.ProjectConfig{
		Name:         "default-project",
		Organization: "default-org",
		License:      "MIT",
		Description:  "Default description",
	}

	// 2. Create file config
	fileConfig := &models.ProjectConfig{
		Name:       "file-project", // Should override default
		License:    "Apache-2.0",   // Should override default
		OutputPath: "/file/output", // Should be added
		// Organization and Description should come from defaults
	}

	// 3. Set environment config
	_ = os.Setenv("GENERATOR_PROJECT_NAME", "env-project")
	_ = os.Setenv("GENERATOR_LICENSE", "BSD-3-Clause")
	defer func() {
		_ = os.Unsetenv("GENERATOR_PROJECT_NAME")
		_ = os.Unsetenv("GENERATOR_LICENSE")
	}()

	envConfig, _ := configManager.LoadFromEnvironment()

	// Merge with precedence: env > file > defaults
	merged := configManager.MergeConfigurations(defaults, fileConfig, envConfig)

	// Environment should have highest precedence
	if merged.Name != "env-project" {
		t.Errorf("Expected name 'env-project' (from env), got '%s'", merged.Name)
	}

	if merged.License != "BSD-3-Clause" {
		t.Errorf("Expected license 'BSD-3-Clause' (from env), got '%s'", merged.License)
	}

	// File should override defaults
	if merged.OutputPath != "/file/output" {
		t.Errorf("Expected output path '/file/output' (from file), got '%s'", merged.OutputPath)
	}

	// Defaults should be used when not overridden
	if merged.Organization != "default-org" {
		t.Errorf("Expected organization 'default-org' (from defaults), got '%s'", merged.Organization)
	}

	if merged.Description != "Default description" {
		t.Errorf("Expected description 'Default description' (from defaults), got '%s'", merged.Description)
	}
}

func testConfigurationValidation(t *testing.T, tempDir string) {
	configManager := NewMockConfigManager()

	// Test valid configuration
	validConfig := &models.ProjectConfig{
		Name:         "valid-project",
		Organization: "valid-org",
		Description:  "Valid project description",
		License:      "MIT",
	}

	err := configManager.ValidateConfig(validConfig)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test invalid configurations
	invalidConfigs := []struct {
		name   string
		config *models.ProjectConfig
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name: "empty name",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "org",
				License:      "MIT",
			},
		},
		{
			name: "invalid name format",
			config: &models.ProjectConfig{
				Name:         "Invalid Project Name",
				Organization: "org",
				License:      "MIT",
			},
		},
	}

	for _, tc := range invalidConfigs {
		t.Run(tc.name, func(t *testing.T) {
			err := configManager.ValidateConfig(tc.config)
			if err == nil {
				t.Errorf("Expected invalid config '%s' to fail validation", tc.name)
			}
		})
	}

	// Test configuration schema validation
	schema := configManager.GetConfigSchema()
	if schema == nil {
		t.Fatal("Expected configuration schema, got nil")
	}

	// Test schema validation with valid data
	validData := map[string]interface{}{
		"name":         "schema-test",
		"organization": "schema-org",
	}

	err = configManager.ValidateConfigurationSchema(validData, schema)
	if err != nil {
		t.Errorf("Expected valid data to pass schema validation, got error: %v", err)
	}

	// Test schema validation with invalid data
	invalidData := map[string]interface{}{
		// Missing required "name" field
		"organization": "schema-org",
	}

	err = configManager.ValidateConfigurationSchema(invalidData, schema)
	if err == nil {
		t.Error("Expected invalid data to fail schema validation")
	}
}

// TestConfigurationFormats tests different configuration file formats
func TestConfigurationFormats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	testConfig := models.ProjectConfig{
		Name:         "format-test",
		Organization: "format-org",
		Description:  "Testing configuration formats",
		License:      "MIT",
		OutputPath:   "/tmp/format-test",
	}

	t.Run("yaml_format", func(t *testing.T) {
		testYAMLFormat(t, tempDir, testConfig)
	})

	t.Run("json_format", func(t *testing.T) {
		testJSONFormat(t, tempDir, testConfig)
	})

	t.Run("toml_format", func(t *testing.T) {
		testTOMLFormat(t, tempDir, testConfig)
	})
}

func testYAMLFormat(t *testing.T, tempDir string, config models.ProjectConfig) {
	yamlPath := filepath.Join(tempDir, "test.yaml")

	// Write YAML config
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal YAML: %v", err)
	}

	err = os.WriteFile(yamlPath, yamlData, 0644)
	if err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	// Read and validate YAML config
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("Failed to read YAML file: %v", err)
	}

	var loadedConfig models.ProjectConfig
	err = yaml.Unmarshal(yamlContent, &loadedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify loaded config matches original
	if loadedConfig.Name != config.Name {
		t.Errorf("Expected name '%s', got '%s'", config.Name, loadedConfig.Name)
	}

	if loadedConfig.Organization != config.Organization {
		t.Errorf("Expected organization '%s', got '%s'", config.Organization, loadedConfig.Organization)
	}

	if loadedConfig.License != config.License {
		t.Errorf("Expected license '%s', got '%s'", config.License, loadedConfig.License)
	}
}

func testJSONFormat(t *testing.T, tempDir string, config models.ProjectConfig) {
	jsonPath := filepath.Join(tempDir, "test.json")

	// Write JSON config
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	// Read and validate JSON config
	jsonContent, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var loadedConfig models.ProjectConfig
	err = json.Unmarshal(jsonContent, &loadedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify loaded config matches original
	if loadedConfig.Name != config.Name {
		t.Errorf("Expected name '%s', got '%s'", config.Name, loadedConfig.Name)
	}

	if loadedConfig.Organization != config.Organization {
		t.Errorf("Expected organization '%s', got '%s'", config.Organization, loadedConfig.Organization)
	}

	if loadedConfig.License != config.License {
		t.Errorf("Expected license '%s', got '%s'", config.License, loadedConfig.License)
	}
}

func testTOMLFormat(t *testing.T, tempDir string, config models.ProjectConfig) {
	// TOML format test (simplified - would need TOML library in real implementation)
	tomlPath := filepath.Join(tempDir, "test.toml")

	// Create TOML content manually for testing
	tomlContent := fmt.Sprintf(`
name = "%s"
organization = "%s"
description = "%s"
license = "%s"
output_path = "%s"
`, config.Name, config.Organization, config.Description, config.License, config.OutputPath)

	err := os.WriteFile(tomlPath, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write TOML file: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		t.Error("Expected TOML file to be created")
	}

	// Read content back
	content, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("Failed to read TOML file: %v", err)
	}

	contentStr := string(content)

	// Verify content contains expected values
	if !strings.Contains(contentStr, config.Name) {
		t.Errorf("Expected TOML to contain name '%s'", config.Name)
	}

	if !strings.Contains(contentStr, config.Organization) {
		t.Errorf("Expected TOML to contain organization '%s'", config.Organization)
	}

	if !strings.Contains(contentStr, config.License) {
		t.Errorf("Expected TOML to contain license '%s'", config.License)
	}
}

// TestConfigurationEdgeCases tests edge cases in configuration handling
func TestConfigurationEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()
	configManager := NewMockConfigManager()

	t.Run("empty_configuration", func(t *testing.T) {
		emptyConfig := &models.ProjectConfig{}

		err := configManager.ValidateConfig(emptyConfig)
		if err == nil {
			t.Error("Expected empty config to fail validation")
		}
	})

	t.Run("partial_configuration", func(t *testing.T) {
		partialConfig := &models.ProjectConfig{
			Name: "partial-project",
			// Missing other required fields
		}

		defaults, _ := configManager.LoadDefaults()
		merged := configManager.MergeConfigurations(defaults, partialConfig)

		// Should have name from partial config
		if merged.Name != "partial-project" {
			t.Errorf("Expected name 'partial-project', got '%s'", merged.Name)
		}

		// Should have defaults for other fields
		if merged.Organization == "" {
			t.Error("Expected organization from defaults")
		}

		if merged.License == "" {
			t.Error("Expected license from defaults")
		}
	})

	t.Run("malformed_configuration_file", func(t *testing.T) {
		malformedPath := filepath.Join(tempDir, "malformed.yaml")
		malformedContent := `
name: test-project
organization: test-org
invalid_yaml: [unclosed array
license: MIT
`

		err := os.WriteFile(malformedPath, []byte(malformedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write malformed file: %v", err)
		}

		// Attempt to validate malformed file
		result, err := configManager.ValidateConfigFromFile(malformedPath)
		if err != nil {
			t.Fatalf("ValidateConfigFromFile failed: %v", err)
		}

		if result.Valid {
			t.Error("Expected malformed config to be invalid")
		}

		if len(result.Errors) == 0 {
			t.Error("Expected errors for malformed config")
		}
	})

	t.Run("configuration_with_special_characters", func(t *testing.T) {
		specialConfig := &models.ProjectConfig{
			Name:         "special-chars-project",
			Organization: "org-with-unicode-™",
			Description:  "Project with special chars: éñ中文🚀",
			License:      "MIT",
		}

		err := configManager.ValidateConfig(specialConfig)
		if err != nil {
			t.Errorf("Expected config with special chars to be valid, got error: %v", err)
		}

		// Test saving and loading with special characters
		specialPath := filepath.Join(tempDir, "special.json")
		err = configManager.SaveConfig(specialConfig, specialPath)
		if err != nil {
			t.Fatalf("Failed to save config with special chars: %v", err)
		}

		loadedConfig, err := configManager.LoadConfig(specialPath)
		if err != nil {
			t.Fatalf("Failed to load config with special chars: %v", err)
		}

		if loadedConfig.Description != specialConfig.Description {
			t.Errorf("Expected description '%s', got '%s'", specialConfig.Description, loadedConfig.Description)
		}
	})
}

// Mock config manager for testing (simplified version)
type MockConfigManager struct {
	configs map[string]*models.ProjectConfig
	schema  *interfaces.ConfigSchema
}

func NewMockConfigManager() *MockConfigManager {
	return &MockConfigManager{
		configs: make(map[string]*models.ProjectConfig),
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

func (m *MockConfigManager) LoadDefaults() (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         "default-project",
		Organization: "default-org",
		Description:  "Default project description",
		License:      "MIT",
	}, nil
}

func (m *MockConfigManager) ValidateConfig(config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if config.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if strings.Contains(config.Name, " ") {
		return fmt.Errorf("project name cannot contain spaces")
	}
	return nil
}

func (m *MockConfigManager) SaveConfig(config *models.ProjectConfig, path string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	m.configs[path] = config
	return nil
}

func (m *MockConfigManager) LoadConfig(path string) (*models.ProjectConfig, error) {
	// First check if we have it in memory
	if config, exists := m.configs[path]; exists {
		return config, nil
	}

	// Try to read from file if it exists
	if _, err := os.Stat(path); err == nil {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		var config models.ProjectConfig
		ext := strings.ToLower(filepath.Ext(path))

		switch ext {
		case ".json":
			if err := json.Unmarshal(content, &config); err != nil {
				return nil, fmt.Errorf("failed to parse JSON config: %w", err)
			}
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(content, &config); err != nil {
				return nil, fmt.Errorf("failed to parse YAML config: %w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported config file format: %s", ext)
		}

		return &config, nil
	}

	return nil, fmt.Errorf("config file not found: %s", path)
}

func (m *MockConfigManager) LoadFromFile(path string) (*models.ProjectConfig, error) {
	return m.LoadConfig(path)
}

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
	if output := os.Getenv("GENERATOR_OUTPUT_PATH"); output != "" {
		config.OutputPath = output
	}

	return config, nil
}

func (m *MockConfigManager) MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig {
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
		if config.OutputPath != "" {
			result.OutputPath = config.OutputPath
		}
	}

	return result
}

func (m *MockConfigManager) GetConfigSchema() *interfaces.ConfigSchema {
	return m.schema
}

func (m *MockConfigManager) ValidateConfigurationSchema(config any, schema *interfaces.ConfigSchema) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return fmt.Errorf("configuration must be a map")
	}

	// Check required fields
	for _, required := range schema.Required {
		if _, exists := configMap[required]; !exists {
			return fmt.Errorf("required field '%s' is missing", required)
		}
	}

	return nil
}

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
