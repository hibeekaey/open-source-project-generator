package config
// Test file - gosec warnings suppressed for test utilities
//nolint:gosec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigurationWorkflow tests the complete configuration workflow
func TestConfigurationWorkflow(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	// Step 1: Generate a template configuration
	template, err := parser.GenerateTemplate("yaml")
	require.NoError(t, err)
	assert.NotEmpty(t, template)

	// Step 2: Parse the template
	config, err := parser.ParseString(template, "yaml")
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Step 3: Apply defaults
	err = validator.ApplyDefaults(config)
	require.NoError(t, err)

	// Step 4: Validate the configuration
	err = validator.Validate(config)
	require.NoError(t, err)

	// Step 5: Get validation report
	report := validator.ValidateWithReport(config)
	assert.True(t, report.Valid)
	assert.False(t, report.HasErrors())
}

// TestConfigurationFileWorkflow tests the file-based configuration workflow
func TestConfigurationFileWorkflow(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "project.yaml")

	// Step 1: Write a template configuration file
	err := parser.WriteTemplate(configPath, "yaml")
	require.NoError(t, err)

	// Step 2: Parse the configuration file
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Step 3: Validate and apply defaults
	err = validator.ValidateAndApplyDefaults(config)
	require.NoError(t, err)

	// Verify the configuration is valid
	assert.Equal(t, "my-project", config.Name)
	assert.NotEmpty(t, config.Components)
}

// TestConfigurationWithValidation tests parsing with validation in one step
func TestConfigurationWithValidation(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "project.yaml")

	// Create a valid configuration file
	validConfig := `
name: test-project
description: Test project
output_dir: ./output
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      name: web-app
      typescript: true
integration:
  generate_docker_compose: true
options:
  use_external_tools: true
`

	err := os.WriteFile(configPath, []byte(validConfig), 0644)
	require.NoError(t, err)

	// Parse and validate in one step
	config, err := parser.ParseFileWithValidation(configPath, validator)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-project", config.Name)

	// Verify defaults were applied
	assert.NotNil(t, config.Components[0].Config["tailwind"])
}

// TestConfigurationErrorHandling tests error handling in the workflow
func TestConfigurationErrorHandling(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	// Test 1: Invalid YAML
	invalidYAML := `
name: test
components:
  - type: nextjs
    invalid syntax here
`

	_, err := parser.ParseString(invalidYAML, "yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML")

	// Test 2: Valid YAML but invalid configuration
	invalidConfig := `
name: ""
output_dir: ./output
components: []
`

	config, err := parser.ParseString(invalidConfig, "yaml")
	require.NoError(t, err)

	err = validator.Validate(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Test 3: Validation report for invalid config
	report := validator.ValidateWithReport(config)
	assert.False(t, report.Valid)
	assert.True(t, report.HasErrors())
}

// TestMultiFormatSupport tests support for multiple configuration formats
func TestMultiFormatSupport(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	tmpDir := t.TempDir()

	formats := []string{"yaml", "yml", "json"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			// Generate template
			template, err := parser.GenerateTemplate(format)
			require.NoError(t, err)

			// Parse template
			config, err := parser.ParseString(template, format)
			require.NoError(t, err)

			// Validate
			err = validator.ValidateAndApplyDefaults(config)
			require.NoError(t, err)

			// Write to file
			ext := format
			if format == "yml" {
				ext = "yaml" // Use yaml extension for yml format
			}
			configPath := filepath.Join(tmpDir, "config."+ext)
			err = parser.WriteTemplateForce(configPath, format)
			require.NoError(t, err)

			// Read back and validate
			config2, err := parser.ParseFileWithValidation(configPath, validator)
			require.NoError(t, err)
			assert.Equal(t, config.Name, config2.Name)
		})
	}
}

// TestComponentTypeValidation tests validation for different component types
func TestComponentTypeValidation(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	componentTypes := []string{"nextjs", "go-backend", "android", "ios"}

	for _, componentType := range componentTypes {
		t.Run(componentType, func(t *testing.T) {
			// Get component schema
			schema, err := validator.GetComponentSchema(componentType)
			require.NoError(t, err)
			assert.Equal(t, componentType, schema.Type)

			// Create a minimal valid configuration
			configYAML := `
name: test-project
output_dir: ./output
components:
  - type: ` + componentType + `
    name: test-component
    enabled: true
    config:
      name: test-component
`

			// Add required fields based on component type
			switch componentType {
			case "go-backend":
				configYAML += `      module: github.com/user/project
`
			case "android":
				configYAML += `      package: com.example.test
`
			case "ios":
				configYAML += `      bundle_id: com.example.test
`
			}

			configYAML += `integration:
  generate_docker_compose: false
options:
  use_external_tools: true
`

			config, err := parser.ParseString(configYAML, "yaml")
			require.NoError(t, err)

			err = validator.ValidateAndApplyDefaults(config)
			require.NoError(t, err)

			// Verify defaults were applied
			assert.NotEmpty(t, config.Components[0].Config)
		})
	}
}

// TestValidationReportDetails tests detailed validation reporting
func TestValidationReportDetails(t *testing.T) {
	validator := NewValidator()

	// Create a configuration with multiple issues
	config := `
name: ""
output_dir: ./output
components:
  - type: nextjs
    name: invalid name!
    enabled: true
    config:
      name: invalid name!
  - type: go-backend
    name: api
    enabled: true
    config:
      name: api
      module: invalid
integration:
  api_endpoints:
    backend: invalid-url
options:
  force_overwrite: true
  create_backup: false
`

	parser := NewParser()
	parsedConfig, err := parser.ParseString(config, "yaml")
	require.NoError(t, err)

	// Get detailed validation report
	report := validator.ValidateWithReport(parsedConfig)
	assert.False(t, report.Valid)
	assert.True(t, report.HasErrors())
	assert.True(t, report.HasWarnings())

	// Verify report string contains useful information
	reportStr := report.String()
	assert.Contains(t, reportStr, "Validation failed")
	assert.Contains(t, reportStr, "Errors:")
	assert.Contains(t, reportStr, "Warnings:")
}
