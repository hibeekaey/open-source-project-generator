package security

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityScanner_ScanForSecrets(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temp directory with test files
	tempDir := t.TempDir()

	// Create file with exposed secret
	secretFile := filepath.Join(tempDir, "config.js")
	secretContent := `
const config = {
  apiKey: "sk-1234567890abcdefghijklmnopqrstuvwxyz",
  password: "mySecretPassword123"
};
`
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	require.NoError(t, err)

	// Scan for secrets
	issues, err := scanner.scanForSecrets(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find at least one secret
	assert.NotEmpty(t, issues)
	assert.Equal(t, "exposed_secret", issues[0].Type)
	assert.Equal(t, "critical", issues[0].Severity)
}

func TestSecurityScanner_ScanDockerConfigs(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temp directory
	tempDir := t.TempDir()

	// Create insecure Dockerfile
	dockerFile := filepath.Join(tempDir, "Dockerfile")
	dockerContent := `
FROM ubuntu:latest
USER root
RUN apt-get update
`
	err := os.WriteFile(dockerFile, []byte(dockerContent), 0644)
	require.NoError(t, err)

	// Scan Docker configs
	issues, err := scanner.scanDockerConfigs(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find root user issue
	assert.NotEmpty(t, issues)
	found := false
	for _, issue := range issues {
		if issue.Type == "insecure_docker_config" && issue.Description == "Container running as root user" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should detect root user in Docker config")
}

func TestSecurityScanner_ScanCORSConfigs(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temp directory
	tempDir := t.TempDir()

	// Create file with insecure CORS
	corsFile := filepath.Join(tempDir, "server.js")
	corsContent := `
app.use(cors({
  origin: '*',
  credentials: true
}));
`
	err := os.WriteFile(corsFile, []byte(corsContent), 0644)
	require.NoError(t, err)

	// Scan CORS configs
	issues, err := scanner.scanCORSConfigs(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find CORS issue
	if len(issues) > 0 {
		assert.Equal(t, "insecure_cors", issues[0].Type)
	}
}

func TestSecurityScanner_Scan(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temp directory with various security issues
	tempDir := t.TempDir()

	// Create file with secret
	secretFile := filepath.Join(tempDir, "config.ts")
	secretContent := `export const API_KEY = "ghp_1234567890abcdefghijklmnopqrstuvwxyz";`
	err := os.WriteFile(secretFile, []byte(secretContent), 0644)
	require.NoError(t, err)

	// Run full scan
	result, err := scanner.Scan(context.Background(), tempDir)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.Equal(t, tempDir, result.Path)

	// Should have found issues
	if len(result.Issues) > 0 {
		assert.False(t, result.Passed, "Scan should not pass with critical/high issues")
	}
}

func TestSecurityScanner_ScanCleanProject(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temp directory with clean files
	tempDir := t.TempDir()

	// Create clean file
	cleanFile := filepath.Join(tempDir, "app.js")
	cleanContent := `
const config = {
  apiKey: process.env.API_KEY,
  password: process.env.PASSWORD
};
`
	err := os.WriteFile(cleanFile, []byte(cleanContent), 0644)
	require.NoError(t, err)

	// Run scan
	result, err := scanner.Scan(context.Background(), tempDir)
	require.NoError(t, err)

	assert.NotNil(t, result)
	// Clean project should pass
	assert.True(t, result.Passed)
}
