package testutils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// RandomGenerationHelpers provides utilities for generating random test data
type RandomGenerationHelpers struct{}

// NewRandomGenerationHelpers creates a new instance of random generation helpers
func NewRandomGenerationHelpers() *RandomGenerationHelpers {
	return &RandomGenerationHelpers{}
}

// GenerateRandomSuffix generates a random suffix for test files
func (r *RandomGenerationHelpers) GenerateRandomSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenerateHexString generates a random hex string of specified length
func (r *RandomGenerationHelpers) GenerateHexString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("test_hex_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// GenerateBytes generates random bytes of specified length
func (r *RandomGenerationHelpers) GenerateBytes(length int) []byte {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to deterministic bytes for testing
		for i := range bytes {
			bytes[i] = byte(i % 256)
		}
	}
	return bytes
}

// GenerateAlphanumeric generates a random alphanumeric string
func (r *RandomGenerationHelpers) GenerateAlphanumeric(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to deterministic generation
			result[i] = charset[i%len(charset)]
		} else {
			result[i] = charset[num.Int64()]
		}
	}
	return string(result)
}

// GenerateBase64String generates a random base64 encoded string
func (r *RandomGenerationHelpers) GenerateBase64String(length int) string {
	bytes := r.GenerateBytes(length)
	return base64.StdEncoding.EncodeToString(bytes)
}

// GenerateSecureID generates a secure random ID
func (r *RandomGenerationHelpers) GenerateSecureID() string {
	return r.GenerateHexString(32)
}

// VersionManagementHelpers provides utilities for version management testing
type VersionManagementHelpers struct{}

// NewVersionManagementHelpers creates a new instance of version management helpers
func NewVersionManagementHelpers() *VersionManagementHelpers {
	return &VersionManagementHelpers{}
}

// ValidateTemplate validates a template structure for testing
func (v *VersionManagementHelpers) ValidateTemplate(templatePath string) error {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template does not exist: %s", templatePath)
	}
	return nil
}

// SetVersion sets a version for testing purposes
func (v *VersionManagementHelpers) SetVersion(name, version string) *models.VersionInfo {
	return &models.VersionInfo{
		Name:           name,
		Language:       "go",
		Type:           "package",
		CurrentVersion: version,
		LatestVersion:  version,
		UpdatedAt:      time.Now(),
		CheckedAt:      time.Now(),
		UpdateSource:   "test",
		IsSecure:       true,
	}
}

// UpdateTemplate simulates updating a template for testing
func (v *VersionManagementHelpers) UpdateTemplate(templatePath, newVersion string) error {
	// Simulate template update
	return nil
}

// GetLatestVersion gets the latest version for testing
func (v *VersionManagementHelpers) GetLatestVersion(packageName string) string {
	// Return a mock latest version
	return "1.0.0"
}

// GetRegistryInfo gets registry information for testing
func (v *VersionManagementHelpers) GetRegistryInfo(packageName string) map[string]interface{} {
	return map[string]interface{}{
		"name":    packageName,
		"version": "1.0.0",
		"latest":  "1.0.0",
	}
}

// GetAffectedTemplates gets affected templates for testing
func (v *VersionManagementHelpers) GetAffectedTemplates(packageName string) []string {
	return []string{
		fmt.Sprintf("template1_%s", packageName),
		fmt.Sprintf("template2_%s", packageName),
	}
}

// UpdateAllTemplates updates all templates for testing
func (v *VersionManagementHelpers) UpdateAllTemplates(version string) error {
	// Simulate updating all templates
	return nil
}

// CheckSecurity checks security for testing
func (v *VersionManagementHelpers) CheckSecurity(packageName string) bool {
	// Mock security check - always return true for testing
	return true
}

// IsAvailable checks if a package is available for testing
func (v *VersionManagementHelpers) IsAvailable(packageName string) bool {
	// Mock availability check - always return true for testing
	return true
}

// BackupTemplates backs up templates for testing
func (v *VersionManagementHelpers) BackupTemplates(templateDir string) error {
	// Simulate backup operation
	backupDir := filepath.Join(templateDir, ".backup")
	return os.MkdirAll(backupDir, 0755)
}

// GetVersionHistory gets version history for testing
func (v *VersionManagementHelpers) GetVersionHistory(packageName string) []string {
	return []string{"1.0.0", "0.9.0", "0.8.0"}
}

// RestoreTemplates restores templates for testing
func (v *VersionManagementHelpers) RestoreTemplates(templateDir string) error {
	// Simulate restore operation
	return nil
}

// GetSupportedPackages gets supported packages for testing
func (v *VersionManagementHelpers) GetSupportedPackages() []string {
	return []string{"react", "vue", "angular", "node", "go"}
}

// FileHelpers provides utilities for file operations in tests
type FileHelpers struct{}

// NewFileHelpers creates a new instance of file helpers
func NewFileHelpers() *FileHelpers {
	return &FileHelpers{}
}

// CreateTempDir creates a temporary directory for testing
func (f *FileHelpers) CreateTempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// CreateTempFile creates a temporary file for testing
func (f *FileHelpers) CreateTempFile(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

// WriteTestFile writes content to a test file
func (f *FileHelpers) WriteTestFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// ReadTestFile reads content from a test file
func (f *FileHelpers) ReadTestFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// CleanupTestFiles removes test files and directories
func (f *FileHelpers) CleanupTestFiles(paths ...string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

// TestAssertions provides common test assertion utilities
type TestAssertions struct{}

// NewTestAssertions creates a new instance of test assertions
func NewTestAssertions() *TestAssertions {
	return &TestAssertions{}
}

// AssertStringContains checks if a string contains a substring
func (a *TestAssertions) AssertStringContains(t interface{}, str, substr, message string) {
	if !strings.Contains(str, substr) {
		if t, ok := t.(interface{ Errorf(string, ...interface{}) }); ok {
			t.Errorf("%s: expected string to contain %q, got %q", message, substr, str)
		}
	}
}

// AssertStringEquals checks if two strings are equal
func (a *TestAssertions) AssertStringEquals(t interface{}, expected, actual, message string) {
	if expected != actual {
		if t, ok := t.(interface{ Errorf(string, ...interface{}) }); ok {
			t.Errorf("%s: expected %q, got %q", message, expected, actual)
		}
	}
}

// AssertNoError checks that an error is nil
func (a *TestAssertions) AssertNoError(t interface{}, err error, message string) {
	if err != nil {
		if t, ok := t.(interface{ Fatalf(string, ...interface{}) }); ok {
			t.Fatalf("%s: unexpected error: %v", message, err)
		}
	}
}

// AssertError checks that an error is not nil
func (a *TestAssertions) AssertError(t interface{}, err error, message string) {
	if err == nil {
		if t, ok := t.(interface{ Errorf(string, ...interface{}) }); ok {
			t.Errorf("%s: expected error but got none", message)
		}
	}
}

// TestSuite provides a complete test utilities suite
type TestSuite struct {
	Random     *RandomGenerationHelpers
	Version    *VersionManagementHelpers
	File       *FileHelpers
	Assertions *TestAssertions
}

// NewTestSuite creates a new complete test utilities suite
func NewTestSuite() *TestSuite {
	return &TestSuite{
		Random:     NewRandomGenerationHelpers(),
		Version:    NewVersionManagementHelpers(),
		File:       NewFileHelpers(),
		Assertions: NewTestAssertions(),
	}
}
