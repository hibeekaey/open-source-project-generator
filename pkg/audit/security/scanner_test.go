package security

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSecurityScanner(t *testing.T) {
	scanner := NewSecurityScanner()
	assert.NotNil(t, scanner)
	assert.NotNil(t, scanner.secretDetector)
}

func TestSecurityScanner_ProjectExists(t *testing.T) {
	scanner := NewSecurityScanner()

	// Test with existing directory
	tempDir, err := os.MkdirTemp("", "test-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = scanner.projectExists(tempDir)
	assert.NoError(t, err)

	// Test with non-existing directory
	err = scanner.projectExists("/non/existing/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project path does not exist")

	// Test with file instead of directory
	tempFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(tempFile, []byte("test"), 0644)
	require.NoError(t, err)

	err = scanner.projectExists(tempFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project path is not a directory")
}

func TestSecurityScanner_ShouldSkipFile(t *testing.T) {
	scanner := NewSecurityScanner()

	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"Skip binary file", "test.exe", true},
		{"Skip image file", "image.jpg", true},
		{"Skip node_modules", "node_modules/package/file.js", true},
		{"Skip .git directory", ".git/config", true},
		{"Don't skip source file", "src/main.go", false},
		{"Don't skip config file", "config.yaml", false},
		{"Skip zip file", "archive.zip", true},
		{"Skip vendor directory", "vendor/package/file.go", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := scanner.shouldSkipFile(tc.filePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityScanner_MaskSecret(t *testing.T) {
	scanner := NewSecurityScanner()

	testCases := []struct {
		name     string
		secret   string
		expected string
	}{
		{"Short secret", "abc", "***"},
		{"Medium secret", "abcdef", "ab**ef"},
		{"Long secret", "abcdefghijklmnop", "ab************op"},
		{"Very short", "ab", "**"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := scanner.maskSecret(tc.secret)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityScanner_DetectSecrets(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "test-secrets-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test file with secrets
	testFile := filepath.Join(tempDir, "config.js")
	testContent := `
const config = {
  apiKey: "AKIAIOSFODNN7EXAMPLE",
  secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  githubToken: "ghp_1234567890abcdef1234567890abcdef12345678",
  password: "mySecretPassword123"
};
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Run secret detection
	result, err := scanner.DetectSecrets(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.NotNil(t, result)
	assert.True(t, result.Summary.FilesScanned > 0)
	assert.True(t, result.Summary.TotalSecrets > 0)
	assert.True(t, len(result.Secrets) > 0)

	// Check that secrets are properly masked
	for _, secret := range result.Secrets {
		assert.NotEmpty(t, secret.Masked)
		assert.Contains(t, secret.Masked, "*")
		assert.NotEqual(t, secret.Secret, secret.Masked)
	}
}

func TestSecurityScanner_ScanVulnerabilities(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory with vulnerable package.json
	tempDir, err := os.MkdirTemp("", "test-vulns-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create package.json with vulnerable dependency
	packageJSON := filepath.Join(tempDir, "package.json")
	packageContent := `{
  "name": "test-project",
  "dependencies": {
    "lodash": "4.17.10",
    "express": "4.16.0"
  }
}`
	err = os.WriteFile(packageJSON, []byte(packageContent), 0644)
	require.NoError(t, err)

	// Run vulnerability scan
	result, err := scanner.ScanVulnerabilities(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.NotNil(t, result)
	assert.NotNil(t, result.ScanTime)
	assert.True(t, result.Summary.Total >= 0)
	assert.NotEmpty(t, result.Recommendations)
}

func TestSecurityScanner_CheckSecurityPolicies(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "test-policies-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test file with potential secret
	testFile := filepath.Join(tempDir, "config.go")
	testContent := `package main

const (
	APIKey = "AKIAIOSFODNN7EXAMPLE"
	Secret = "mySecretPassword123"
)
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Run security policy check
	result, err := scanner.CheckSecurityPolicies(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.NotNil(t, result)
	assert.True(t, result.Summary.TotalPolicies > 0)
	assert.True(t, len(result.Policies) > 0)
	assert.True(t, result.Score >= 0 && result.Score <= 100)

	// Check that policies are properly defined
	for _, policy := range result.Policies {
		assert.NotEmpty(t, policy.ID)
		assert.NotEmpty(t, policy.Name)
		assert.NotEmpty(t, policy.Description)
		assert.NotEmpty(t, policy.Category)
		assert.NotEmpty(t, policy.Severity)
	}
}

func TestSecurityScanner_AuditSecurity(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "test-audit-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFile := filepath.Join(tempDir, "main.go")
	testContent := `package main

import "fmt"

const APIKey = "AKIAIOSFODNN7EXAMPLE"

func main() {
	fmt.Println("Hello, World!")
}
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create package.json
	packageJSON := filepath.Join(tempDir, "package.json")
	packageContent := `{
  "name": "test-project",
  "dependencies": {
    "lodash": "4.17.10"
  }
}`
	err = os.WriteFile(packageJSON, []byte(packageContent), 0644)
	require.NoError(t, err)

	// Run security audit
	result, err := scanner.AuditSecurity(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.NotNil(t, result)
	assert.True(t, result.Score >= 0 && result.Score <= 100)
	assert.NotNil(t, result.Vulnerabilities)
	assert.NotNil(t, result.PolicyViolations)
	assert.NotEmpty(t, result.Recommendations)
}

func TestSecurityScanner_CheckHardcodedSecrets(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "test-hardcoded-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test file with hardcoded secrets
	testFile := filepath.Join(tempDir, "config.py")
	testContent := `
API_KEY = "AKIAIOSFODNN7EXAMPLE"
SECRET_KEY = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
GITHUB_TOKEN = "ghp_1234567890abcdef1234567890abcdef12345678"
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Run hardcoded secrets check
	violations, err := scanner.checkHardcodedSecrets(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.True(t, len(violations) > 0)
	for _, violation := range violations {
		assert.Equal(t, "SEC-001", violation.Policy)
		assert.Equal(t, "critical", violation.Severity)
		assert.NotEmpty(t, violation.Description)
		assert.NotEmpty(t, violation.File)
		assert.True(t, violation.Line > 0)
	}
}

func TestSecurityScanner_CheckSecureConfiguration(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "test-config-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create insecure docker-compose.yml
	dockerCompose := filepath.Join(tempDir, "docker-compose.yml")
	dockerContent := `version: '3'
services:
  app:
    image: myapp
    privileged: true
    security_opt:
      - seccomp:unconfined
`
	err = os.WriteFile(dockerCompose, []byte(dockerContent), 0644)
	require.NoError(t, err)

	// Create insecure Dockerfile
	dockerfile := filepath.Join(tempDir, "Dockerfile")
	dockerfileContent := `FROM ubuntu:20.04
USER root
COPY . /
`
	err = os.WriteFile(dockerfile, []byte(dockerfileContent), 0644)
	require.NoError(t, err)

	// Run secure configuration check
	violations, err := scanner.checkSecureConfiguration(tempDir)
	require.NoError(t, err)

	// Verify results
	assert.True(t, len(violations) > 0)
	for _, violation := range violations {
		assert.Equal(t, "SEC-003", violation.Policy)
		assert.Equal(t, "medium", violation.Severity)
		assert.NotEmpty(t, violation.Description)
		assert.NotEmpty(t, violation.File)
		assert.True(t, violation.Line > 0)
	}
}

func TestSecurityScanner_CalculateSecurityScore(t *testing.T) {
	scanner := NewSecurityScanner()

	testCases := []struct {
		name         string
		result       *interfaces.SecurityAuditResult
		secretResult *interfaces.SecretScanResult
		expectedMin  float64
		expectedMax  float64
	}{
		{
			name: "Perfect score",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{
					HighConfidence:   0,
					MediumConfidence: 0,
					LowConfidence:    0,
				},
			},
			expectedMin: 100.0,
			expectedMax: 100.0,
		},
		{
			name: "Critical vulnerability",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "critical"},
				},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{},
			},
			expectedMin: 70.0,
			expectedMax: 80.0,
		},
		{
			name: "High confidence secret",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{
					HighConfidence: 1,
				},
			},
			expectedMin: 80.0,
			expectedMax: 90.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := scanner.calculateSecurityScore(tc.result, tc.secretResult)
			assert.True(t, score >= tc.expectedMin && score <= tc.expectedMax,
				"Score %f should be between %f and %f", score, tc.expectedMin, tc.expectedMax)
		})
	}
}

func TestSecurityScanner_GenerateSecurityRecommendations(t *testing.T) {
	scanner := NewSecurityScanner()

	testCases := []struct {
		name         string
		result       *interfaces.SecurityAuditResult
		secretResult *interfaces.SecretScanResult
		expectedMin  int
	}{
		{
			name: "With vulnerabilities",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high"},
				},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{},
			},
			expectedMin: 2,
		},
		{
			name: "With secrets",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{
					TotalSecrets: 1,
				},
			},
			expectedMin: 3,
		},
		{
			name: "Clean project",
			result: &interfaces.SecurityAuditResult{
				Vulnerabilities:  []interfaces.Vulnerability{},
				PolicyViolations: []interfaces.PolicyViolation{},
			},
			secretResult: &interfaces.SecretScanResult{
				Summary: interfaces.SecretScanSummary{},
			},
			expectedMin: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recommendations := scanner.generateSecurityRecommendations(tc.result, tc.secretResult)
			assert.True(t, len(recommendations) >= tc.expectedMin,
				"Expected at least %d recommendations, got %d", tc.expectedMin, len(recommendations))
			for _, rec := range recommendations {
				assert.NotEmpty(t, rec)
			}
		})
	}
}

func TestGetSecretDetectionRules(t *testing.T) {
	detector := NewSecretDetector()
	rules := detector.GetRules()
	assert.NotEmpty(t, rules)

	for _, rule := range rules {
		assert.NotEmpty(t, rule.Name)
		assert.NotEmpty(t, rule.Pattern)
		assert.True(t, rule.Confidence > 0 && rule.Confidence <= 1)
	}
}

func TestSecurityScanner_ScanFileForSecrets(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create temporary file with secrets
	tempDir, err := os.MkdirTemp("", "test-file-secrets-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "secrets.txt")
	testContent := `
AWS_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
GITHUB_TOKEN=ghp_1234567890abcdef1234567890abcdef12345678
API_KEY=sk-1234567890abcdef1234567890abcdef
PASSWORD=mySecretPassword123
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Scan file for secrets using the secret detector
	secrets, err := scanner.secretDetector.scanFileForSecrets(testFile)
	require.NoError(t, err)

	// Verify results
	assert.True(t, len(secrets) > 0)
	for _, secret := range secrets {
		assert.NotEmpty(t, secret.Type)
		assert.Equal(t, testFile, secret.File)
		assert.True(t, secret.Line > 0)
		assert.True(t, secret.Column > 0)
		assert.NotEmpty(t, secret.Secret)
		assert.True(t, secret.Confidence > 0)
		assert.NotEmpty(t, secret.Rule)
		assert.NotEmpty(t, secret.Pattern)
		assert.NotEmpty(t, secret.Masked)
		assert.Contains(t, secret.Masked, "*")
	}
}

// Benchmark tests
func BenchmarkSecurityScanner_DetectSecrets(b *testing.B) {
	scanner := NewSecurityScanner()

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "bench-secrets-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create multiple test files
	for i := 0; i < 10; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("file%d.js", i))
		testContent := `
const config = {
  apiKey: "AKIAIOSFODNN7EXAMPLE",
  secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  githubToken: "ghp_1234567890abcdef1234567890abcdef12345678"
};
`
		err = os.WriteFile(testFile, []byte(testContent), 0644)
		require.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.DetectSecrets(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkSecurityScanner_ScanVulnerabilities(b *testing.B) {
	scanner := NewSecurityScanner()

	// Create temporary directory with package.json
	tempDir, err := os.MkdirTemp("", "bench-vulns-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	packageJSON := filepath.Join(tempDir, "package.json")
	packageContent := `{
  "name": "test-project",
  "dependencies": {
    "lodash": "4.17.10",
    "express": "4.16.0"
  }
}`
	err = os.WriteFile(packageJSON, []byte(packageContent), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.ScanVulnerabilities(tempDir)
		require.NoError(b, err)
	}
}
