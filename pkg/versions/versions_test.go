package versions

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFrom(t *testing.T) {
	// Create a temporary test config
	tmpDir := t.TempDir()
	testConfigPath := filepath.Join(tmpDir, "versions.yaml")

	// Use clearly fake test versions to avoid confusion with real versions
	testNextJSVersion := "99.0.0"
	testReactVersion := "99.1.0"
	testGoVersion := "99.2.0"

	testConfig := fmt.Sprintf(`
frontend:
  nextjs:
    version: "%s"
    package: "create-next-app"
  react:
    version: "%s"
backend:
  go:
    version: "%s"
    docker_tag: "%s-alpine"
metadata:
  last_updated: "2024-01-01"
  schema_version: "1.0.0"
`, testNextJSVersion, testReactVersion, testGoVersion, testGoVersion)

	err := os.WriteFile(testConfigPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	// Test loading
	config, err := LoadFrom(testConfigPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify values
	assert.Equal(t, testNextJSVersion, config.Frontend.NextJS.Version)
	assert.Equal(t, "create-next-app", config.Frontend.NextJS.Package)
	assert.Equal(t, testReactVersion, config.Frontend.React.Version)
	assert.Equal(t, testGoVersion, config.Backend.Go.Version)
	assert.Equal(t, testGoVersion+"-alpine", config.Backend.Go.DockerTag)
	assert.Equal(t, "2024-01-01", config.Metadata.LastUpdated)
	assert.Equal(t, "1.0.0", config.Metadata.SchemaVersion)
}

func TestLoadFrom_FileNotFound(t *testing.T) {
	_, err := LoadFrom("/nonexistent/path/versions.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read versions config")
}

func TestLoadFrom_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	testConfigPath := filepath.Join(tmpDir, "invalid.yaml")

	invalidYAML := `
frontend:
  nextjs:
    version: "99.0.0"
  - invalid yaml structure
`

	err := os.WriteFile(testConfigPath, []byte(invalidYAML), 0644)
	require.NoError(t, err)

	_, err = LoadFrom(testConfigPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse versions config")
}

func TestGetAndReload(t *testing.T) {
	// Save original config path
	originalPath := configPath
	defer func() {
		SetConfigPath(originalPath)
	}()

	// Create a temporary test config
	tmpDir := t.TempDir()
	testConfigPath := filepath.Join(tmpDir, "versions.yaml")

	// Use clearly fake test versions to avoid confusion with real versions
	initialVersion := "99.0.0"
	updatedVersion := "99.1.0"

	testConfig := fmt.Sprintf(`
frontend:
  nextjs:
    version: "%s"
    package: "create-next-app"
metadata:
  last_updated: "2024-01-01"
  schema_version: "1.0.0"
`, initialVersion)

	err := os.WriteFile(testConfigPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	SetConfigPath(testConfigPath)

	// First Get should load the config
	config1, err := Get()
	require.NoError(t, err)
	assert.Equal(t, initialVersion, config1.Frontend.NextJS.Version)

	// Second Get should return cached config
	config2, err := Get()
	require.NoError(t, err)
	assert.Equal(t, config1, config2)

	// Update the file
	updatedConfig := fmt.Sprintf(`
frontend:
  nextjs:
    version: "%s"
    package: "create-next-app"
metadata:
  last_updated: "2024-01-02"
  schema_version: "1.0.0"
`, updatedVersion)

	err = os.WriteFile(testConfigPath, []byte(updatedConfig), 0644)
	require.NoError(t, err)

	// Reload should pick up the new version
	err = Reload()
	require.NoError(t, err)

	config3, err := Get()
	require.NoError(t, err)
	assert.Equal(t, updatedVersion, config3.Frontend.NextJS.Version)
}

func TestSetConfigPath(t *testing.T) {
	originalPath := configPath
	defer func() {
		SetConfigPath(originalPath)
	}()

	newPath := "/custom/path/versions.yaml"
	SetConfigPath(newPath)

	assert.Equal(t, newPath, configPath)
}

func TestCompleteConfigStructure(t *testing.T) {
	// Test with the actual versions.yaml file if it exists
	actualConfigPath := "../../configs/versions.yaml"
	if _, err := os.Stat(actualConfigPath); os.IsNotExist(err) {
		t.Skip("Skipping test: actual versions.yaml not found")
	}

	config, err := LoadFrom(actualConfigPath)
	require.NoError(t, err)

	// Verify all major sections are present
	assert.NotEmpty(t, config.Frontend.NextJS.Version)
	assert.NotEmpty(t, config.Frontend.React.Version)
	assert.NotEmpty(t, config.Backend.Go.Version)
	assert.NotEmpty(t, config.Backend.Frameworks.Gin.Version)
	assert.NotEmpty(t, config.Android.Kotlin.Version)
	assert.NotEmpty(t, config.Android.Gradle.Version)
	assert.NotEmpty(t, config.IOS.Swift.Version)
	assert.NotEmpty(t, config.Docker.Alpine.Version)
	assert.NotEmpty(t, config.Infrastructure.Terraform.Version)
	assert.NotEmpty(t, config.Metadata.LastUpdated)
	assert.NotEmpty(t, config.Metadata.SchemaVersion)
}
