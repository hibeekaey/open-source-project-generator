package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigValidator(t *testing.T) {
	cv := NewConfigValidator()

	assert.NotNil(t, cv)
	assert.NotNil(t, cv.schemaManager)

	// Verify that the schema manager has default schemas
	schemas := cv.schemaManager.ListSchemas()
	assert.Contains(t, schemas, "package.json")
	assert.Contains(t, schemas, "tsconfig.json")
}

func TestConfigValidator_ValidateConfigurationFiles(t *testing.T) {
	cv := NewConfigValidator()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "config_validator_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a valid package.json
	packageJSON := `{
		"name": "test-package",
		"version": "1.0.0",
		"description": "Test package"
	}`

	packageJSONPath := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(packageJSONPath, []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Validate the configuration files
	results, err := cv.ValidateConfigurationFiles(tempDir)
	require.NoError(t, err)
	assert.NotEmpty(t, results)

	// Should have at least one result
	assert.Len(t, results, 1)

	// The result should be valid
	result := results[0]
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)
}

func TestConfigValidator_ValidateJSONFile(t *testing.T) {
	cv := NewConfigValidator()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "config_validator_json_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test valid JSON
	validJSON := `{
		"name": "test-package",
		"version": "1.0.0"
	}`

	validJSONPath := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(validJSONPath, []byte(validJSON), 0644)
	require.NoError(t, err)

	result, err := cv.jsonValidator.ValidateJSONFile(validJSONPath)
	require.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)

	// Test invalid JSON
	invalidJSON := `{
		"name": "test-package"
		"version": "1.0.0"
	}` // Missing comma

	invalidJSONPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSONPath, []byte(invalidJSON), 0644)
	require.NoError(t, err)

	result, err = cv.jsonValidator.ValidateJSONFile(invalidJSONPath)
	require.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Greater(t, result.Summary.ErrorCount, 0)
}

func TestConfigValidator_ValidateEnvFile(t *testing.T) {
	cv := NewConfigValidator()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "config_validator_env_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test valid .env file
	validEnv := `NODE_ENV=production
DATABASE_URL=postgresql://localhost:5432/mydb
PORT=3000
API_KEY=sk-1234567890abcdef`

	envPath := filepath.Join(tempDir, ".env")
	err = os.WriteFile(envPath, []byte(validEnv), 0644)
	require.NoError(t, err)

	result, err := cv.envValidator.ValidateEnvFile(envPath)
	require.NoError(t, err)
	assert.True(t, result.Valid)

	// Should have warnings for potential secrets
	assert.Greater(t, result.Summary.WarningCount, 0)

	// Check that API_KEY was flagged as potential secret
	foundSecretWarning := false
	for _, warning := range result.Warnings {
		if warning.Type == "security" && warning.Rule == "config.env.secrets" {
			foundSecretWarning = true
			break
		}
	}
	assert.True(t, foundSecretWarning, "Expected to find secret warning for API_KEY")
}

func TestConfigValidator_ValidateDockerfile(t *testing.T) {
	cv := NewConfigValidator()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "config_validator_docker_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test valid Dockerfile
	validDockerfile := `FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
USER node
CMD ["npm", "start"]`

	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(validDockerfile), 0644)
	require.NoError(t, err)

	result, err := cv.dockerValidator.ValidateDockerfile(dockerfilePath)
	require.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.Summary.ErrorCount)

	// Test Dockerfile without FROM instruction
	invalidDockerfile := `WORKDIR /app
COPY . .
CMD ["npm", "start"]`

	invalidDockerfilePath := filepath.Join(tempDir, "Dockerfile.invalid")
	err = os.WriteFile(invalidDockerfilePath, []byte(invalidDockerfile), 0644)
	require.NoError(t, err)

	result, err = cv.dockerValidator.ValidateDockerfile(invalidDockerfilePath)
	require.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Greater(t, result.Summary.ErrorCount, 0)
	assert.Greater(t, result.Summary.MissingRequired, 0)
}
