package configs_test

//nolint:gosec // Test file - file inclusion via variable is expected

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestExampleConfigurations validates all example configuration files
func TestExampleConfigurations(t *testing.T) {
	// List of example configuration files to test
	configFiles := []string{
		"fullstack-complete.yaml",
		"frontend-only.yaml",
		"backend-only.yaml",
		"mobile-only.yaml",
		"minimal.yaml",
		"web-and-api.yaml",
		"advanced-options.yaml",
		"example-project.yaml",
	}

	validator := config.NewValidator()

	for _, filename := range configFiles {
		t.Run(filename, func(t *testing.T) {
			// Read configuration file
			configPath := filepath.Join(".", filename)
			data, err := os.ReadFile(configPath)
			require.NoError(t, err, "Failed to read config file: %s", filename)

			// Parse YAML
			var projectConfig models.ProjectConfig
			err = yaml.Unmarshal(data, &projectConfig)
			require.NoError(t, err, "Failed to parse YAML: %s", filename)

			// Validate configuration
			err = validator.ValidateAndApplyDefaults(&projectConfig)
			assert.NoError(t, err, "Configuration validation failed for: %s", filename)

			// Additional checks
			assert.NotEmpty(t, projectConfig.Name, "Project name should not be empty")
			assert.NotEmpty(t, projectConfig.OutputDir, "Output directory should not be empty")

			// Check that at least one component is enabled
			hasEnabledComponent := false
			for _, comp := range projectConfig.Components {
				if comp.Enabled {
					hasEnabledComponent = true
					break
				}
			}
			assert.True(t, hasEnabledComponent, "At least one component should be enabled")
		})
	}
}

// TestMinimalConfiguration tests the minimal configuration
func TestMinimalConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "minimal.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify minimal config has required fields
	assert.Equal(t, "minimal-project", projectConfig.Name)
	assert.Equal(t, "./minimal-project", projectConfig.OutputDir)
	assert.Len(t, projectConfig.Components, 1)
	assert.True(t, projectConfig.Components[0].Enabled)
}

// TestFullStackConfiguration tests the full-stack configuration
func TestFullStackConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "fullstack-complete.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify all components are present
	assert.Equal(t, "fullstack-app", projectConfig.Name)
	assert.Len(t, projectConfig.Components, 4)

	// Check component types
	componentTypes := make(map[string]bool)
	for _, comp := range projectConfig.Components {
		componentTypes[comp.Type] = true
	}
	assert.True(t, componentTypes["nextjs"])
	assert.True(t, componentTypes["go-backend"])
	assert.True(t, componentTypes["android"])
	assert.True(t, componentTypes["ios"])

	// Verify integration settings
	assert.True(t, projectConfig.Integration.GenerateDockerCompose)
	assert.True(t, projectConfig.Integration.GenerateScripts)
	assert.NotEmpty(t, projectConfig.Integration.APIEndpoints)
}

// TestFrontendOnlyConfiguration tests the frontend-only configuration
func TestFrontendOnlyConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "frontend-only.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify only frontend is enabled
	enabledComponents := 0
	for _, comp := range projectConfig.Components {
		if comp.Enabled {
			enabledComponents++
			assert.Equal(t, "nextjs", comp.Type)
		}
	}
	assert.Equal(t, 1, enabledComponents)
}

// TestBackendOnlyConfiguration tests the backend-only configuration
func TestBackendOnlyConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "backend-only.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify only backend is enabled
	enabledComponents := 0
	for _, comp := range projectConfig.Components {
		if comp.Enabled {
			enabledComponents++
			assert.Equal(t, "go-backend", comp.Type)
		}
	}
	assert.Equal(t, 1, enabledComponents)

	// Verify backend-specific configuration
	for _, comp := range projectConfig.Components {
		if comp.Type == "go-backend" && comp.Enabled {
			assert.NotEmpty(t, comp.Config["module"])
			assert.NotEmpty(t, comp.Config["framework"])
		}
	}
}

// TestMobileOnlyConfiguration tests the mobile-only configuration
func TestMobileOnlyConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "mobile-only.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify only mobile components are enabled
	enabledComponents := 0
	for _, comp := range projectConfig.Components {
		if comp.Enabled {
			enabledComponents++
			assert.True(t, comp.Type == "android" || comp.Type == "ios")
		}
	}
	assert.Equal(t, 2, enabledComponents)
}

// TestAdvancedOptionsConfiguration tests the advanced options configuration
func TestAdvancedOptionsConfiguration(t *testing.T) {
	configPath := filepath.Join(".", "advanced-options.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()
	err = validator.ValidateAndApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Verify all components are enabled
	for _, comp := range projectConfig.Components {
		assert.True(t, comp.Enabled)
	}

	// Verify advanced options
	assert.True(t, projectConfig.Options.UseExternalTools)
	assert.False(t, projectConfig.Options.DryRun)
	assert.True(t, projectConfig.Options.Verbose)
	assert.True(t, projectConfig.Options.CreateBackup)
	assert.False(t, projectConfig.Options.ForceOverwrite)

	// Verify extensive environment variables
	assert.NotEmpty(t, projectConfig.Integration.SharedEnvironment)
	assert.Contains(t, projectConfig.Integration.SharedEnvironment, "LOG_LEVEL")
	assert.Contains(t, projectConfig.Integration.SharedEnvironment, "API_TIMEOUT")
}

// TestConfigurationValidationReport tests validation reporting
func TestConfigurationValidationReport(t *testing.T) {
	configPath := filepath.Join(".", "fullstack-complete.yaml")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var projectConfig models.ProjectConfig
	err = yaml.Unmarshal(data, &projectConfig)
	require.NoError(t, err)

	validator := config.NewValidator()

	// Apply defaults first
	err = validator.ApplyDefaults(&projectConfig)
	require.NoError(t, err)

	// Get validation report
	report := validator.ValidateWithReport(&projectConfig)

	assert.True(t, report.Valid, "Configuration should be valid")
	assert.Empty(t, report.Errors, "Should have no errors")

	// Report may have warnings, which is fine
	if report.HasWarnings() {
		t.Logf("Validation warnings: %v", report.Warnings)
	}
}
