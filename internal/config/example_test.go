package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExampleProjectConfig tests the example project configuration file
func TestExampleProjectConfig(t *testing.T) {
	parser := NewParser()
	validator := NewValidator()

	// Path to example config (relative to this test file)
	configPath := filepath.Join("..", "..", "configs", "example-project.yaml")

	// Parse the example configuration
	config, err := parser.ParseFile(configPath)
	require.NoError(t, err, "Failed to parse example-project.yaml")
	assert.NotNil(t, config)

	// Validate the configuration
	err = validator.ValidateAndApplyDefaults(config)
	require.NoError(t, err, "Example configuration should be valid")

	// Verify basic structure
	assert.Equal(t, "example-fullstack-app", config.Name)
	assert.NotEmpty(t, config.Description)
	assert.Equal(t, "./example-fullstack-app", config.OutputDir)

	// Verify components
	assert.Len(t, config.Components, 4)

	// Check Next.js component
	nextjsComp := config.Components[0]
	assert.Equal(t, "nextjs", nextjsComp.Type)
	assert.Equal(t, "web-app", nextjsComp.Name)
	assert.True(t, nextjsComp.Enabled)

	// Check Go backend component
	goComp := config.Components[1]
	assert.Equal(t, "go-backend", goComp.Type)
	assert.Equal(t, "api-server", goComp.Name)
	assert.True(t, goComp.Enabled)

	// Check Android component (disabled)
	androidComp := config.Components[2]
	assert.Equal(t, "android", androidComp.Type)
	assert.False(t, androidComp.Enabled)

	// Check iOS component (disabled)
	iosComp := config.Components[3]
	assert.Equal(t, "ios", iosComp.Type)
	assert.False(t, iosComp.Enabled)

	// Verify integration config
	assert.True(t, config.Integration.GenerateDockerCompose)
	assert.True(t, config.Integration.GenerateScripts)
	assert.NotEmpty(t, config.Integration.APIEndpoints)
	assert.NotEmpty(t, config.Integration.SharedEnvironment)

	// Verify options
	assert.True(t, config.Options.UseExternalTools)
	assert.False(t, config.Options.DryRun)
	assert.True(t, config.Options.CreateBackup)
	assert.False(t, config.Options.ForceOverwrite)

	// Get validation report
	report := validator.ValidateWithReport(config)
	assert.True(t, report.Valid)
	assert.False(t, report.HasErrors())
}
