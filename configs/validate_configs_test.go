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
		"minimal.yaml",
		"fullstack-complete.yaml",
		"mobile-app-with-backend.yaml",
		"frontend-only.yaml",
		"backend-only.yaml",
		"mobile-only.yaml",
		"advanced-options.yaml",
		"performance-optimized.yaml",
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
