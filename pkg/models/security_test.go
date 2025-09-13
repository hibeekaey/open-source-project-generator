package models

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockFileSystem implements FileSystemInterface for testing
type MockFileSystem struct {
	StatFunc  func(name string) (os.FileInfo, error)
	IsAbsFunc func(path string) bool
	CleanFunc func(path string) string
}

func (mfs *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if mfs.StatFunc != nil {
		return mfs.StatFunc(name)
	}
	return &MockFileInfo{name: name, isDir: true}, nil
}

func (mfs *MockFileSystem) IsAbs(path string) bool {
	if mfs.IsAbsFunc != nil {
		return mfs.IsAbsFunc(path)
	}
	return strings.HasPrefix(path, "/")
}

func (mfs *MockFileSystem) Clean(path string) string {
	if mfs.CleanFunc != nil {
		return mfs.CleanFunc(path)
	}
	// Simple mock clean implementation that behaves like filepath.Clean
	if path == "/" {
		return "/"
	}
	return strings.TrimSuffix(path, "/")
}

// MockFileInfo implements os.FileInfo for testing
type MockFileInfo struct {
	name  string
	isDir bool
}

func (mfi *MockFileInfo) Name() string       { return mfi.name }
func (mfi *MockFileInfo) Size() int64        { return 1024 }
func (mfi *MockFileInfo) Mode() os.FileMode  { return 0o755 }
func (mfi *MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi *MockFileInfo) IsDir() bool        { return mfi.isDir }
func (mfi *MockFileInfo) Sys() interface{}   { return nil }

// CreateStandardMockFS creates a mock file system with standard behavior
func CreateStandardMockFS() *MockFileSystem {
	return &MockFileSystem{
		StatFunc: func(name string) (os.FileInfo, error) {
			return &MockFileInfo{name: name, isDir: true}, nil
		},
		IsAbsFunc: func(path string) bool {
			return strings.HasPrefix(path, "/")
		},
		CleanFunc: func(path string) string {
			if path == "/" {
				return "/"
			}
			return strings.TrimSuffix(path, "/")
		},
	}
}

// SecurityTestValidator creates a security validator with standardized mock dependencies
func SecurityTestValidator() *SecurityValidator {
	return NewSecurityValidatorWithFS(CreateStandardMockFS())
}

func TestSecurityConfig_Validation(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	// Get standard test cases
	tests := testSuite.StandardSecurityConfigTestCases()

	// Add some additional edge cases
	additionalTests := []SecurityTestCase{
		{
			Name:   "missing required fields",
			Config: &SecurityConfig{
				// Missing required fields
			},
			ExpectValid: false,
			ExpectError: "TempFileRandomLength",
			Description: "Empty configuration should fail validation",
		},
		{
			Name: "temp file random length too large",
			Config: &SecurityConfig{
				TempFileRandomLength: 128, // Too large
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			ExpectValid: false,
			ExpectError: "max",
			Description: "Configuration with excessive random length should fail",
		},
		{
			Name: "max file size too small",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          512, // Too small
			},
			ExpectValid: false,
			ExpectError: "min",
			Description: "Configuration with insufficient max file size should fail",
		},
	}

	tests = append(tests, additionalTests...)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := validator.ValidateSecurityConfig(tt.Config.(*SecurityConfig))

			assert.Equal(t, tt.ExpectValid, result.Valid, "Validation result should match expected: %s", tt.Description)

			if !tt.ExpectValid && tt.ExpectError != "" {
				assert.NotEmpty(t, result.Errors, "Should have validation errors")

				// Check if expected error is present
				found := false
				for _, err := range result.Errors {
					if err.Tag == tt.ExpectError || err.Field == tt.ExpectError {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error '%s' not found in validation errors", tt.ExpectError)
			}
		})
	}
}

func TestRandomConfig_Validation(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	// Get standard test cases
	tests := testSuite.StandardRandomConfigTestCases()

	// Add some additional edge cases
	additionalTests := []SecurityTestCase{
		{
			Name:   "missing required fields",
			Config: &RandomConfig{
				// Missing required fields
			},
			ExpectValid: false,
			ExpectError: "DefaultSuffixLength",
			Description: "Empty random configuration should fail validation",
		},
		{
			Name: "suffix length too large",
			Config: &RandomConfig{
				DefaultSuffixLength: 128, // Too large
				IDFormat:            "hex",
				MinEntropyBytes:     32,
			},
			ExpectValid: false,
			ExpectError: "max",
			Description: "Configuration with excessive suffix length should fail",
		},
		{
			Name: "entropy bytes too large",
			Config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     2048, // Too large
			},
			ExpectValid: false,
			ExpectError: "max",
			Description: "Configuration with excessive entropy should fail",
		},
		{
			Name: "prefix length too large",
			Config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     32,
				IDPrefixLength:      32, // Too large
			},
			ExpectValid: false,
			ExpectError: "max",
			Description: "Configuration with excessive prefix length should fail",
		},
	}

	tests = append(tests, additionalTests...)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := validator.ValidateRandomConfig(tt.Config.(*RandomConfig))

			assert.Equal(t, tt.ExpectValid, result.Valid, "Validation result should match expected: %s", tt.Description)

			if !tt.ExpectValid && tt.ExpectError != "" {
				assert.NotEmpty(t, result.Errors, "Should have validation errors")

				// Check if expected error is present
				found := false
				for _, err := range result.Errors {
					if err.Tag == tt.ExpectError || err.Field == tt.ExpectError {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error '%s' not found in validation errors", tt.ExpectError)
			}
		})
	}
}

func TestCombinedConfig_Validation(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	secConfig := testSuite.Fixtures.StandardSecurityConfig()
	randConfig := testSuite.Fixtures.StandardRandomConfig()

	t.Run("compatible configs", func(t *testing.T) {
		result := validator.ValidateCombinedConfig(secConfig, randConfig)
		assert.True(t, result.Valid, "Compatible configs should be valid")
	})

	t.Run("mismatched random lengths", func(t *testing.T) {
		mismatchedRandConfig := testSuite.Fixtures.MinimalRandomConfig()
		mismatchedRandConfig.DefaultSuffixLength = 12 // Different from security config

		result := validator.ValidateCombinedConfig(secConfig, mismatchedRandConfig)
		assert.True(t, result.Valid, "Should still be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have compatibility warnings")
	})
}

func TestSecurityValidator_FilePermissions(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	// Use standardized permission test cases
	tests := testSuite.PermissionTestCases()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := validator.ValidateSecurityConfig(tt.Config.(*SecurityConfig))

			assert.Equal(t, tt.ExpectValid, result.Valid, "Validation result should match expected: %s", tt.Description)

			if !tt.ExpectValid && tt.ExpectError != "" {
				assert.NotEmpty(t, result.Errors, "Should have validation errors")

				// Check if expected error is present
				found := false
				for _, err := range result.Errors {
					if err.Tag == tt.ExpectError || err.Field == tt.ExpectError {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error '%s' not found in validation errors", tt.ExpectError)
			}
		})
	}
}

func TestSecurityValidator_DangerousDirectories(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	dangerousDirs := testSuite.Fixtures.DangerousDirectories()

	for _, dir := range dangerousDirs {
		t.Run("dangerous_dir_"+strings.ReplaceAll(dir, "/", "_"), func(t *testing.T) {
			config := testSuite.Fixtures.MinimalSecurityConfig()
			config.AllowedTempDirs = []string{dir}

			result := validator.ValidateSecurityConfig(config)
			assert.False(t, result.Valid, "Dangerous directory should be invalid: %s", dir)
			assert.NotEmpty(t, result.Errors, "Should have errors for dangerous directory")
		})
	}
}

func TestSecurityValidator_SecurityWarnings(t *testing.T) {
	// Use standardized test suite
	testSuite := NewSecurityValidationTestSuite()
	validator := SecurityTestValidator()

	t.Run("weak random length warning", func(t *testing.T) {
		config := testSuite.Fixtures.MinimalSecurityConfig()
		config.TempFileRandomLength = 12 // Less than recommended

		result := validator.ValidateSecurityConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have warnings for weak random length")
	})

	t.Run("path validation disabled warning", func(t *testing.T) {
		config := &SecurityConfig{
			TempFileRandomLength: 16,
			AllowedTempDirs:      []string{"/tmp"},
			FilePermissions:      0o600,
			EnablePathValidation: false, // Disabled
			MaxFileSize:          1024,
		}

		result := validator.ValidateSecurityConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")

		// Check for high impact warning
		found := false
		for _, warning := range result.Warnings {
			if warning.Field == "EnablePathValidation" && warning.Impact == "high" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have high impact warning for disabled path validation")
	})

	t.Run("secure cleanup disabled warning", func(t *testing.T) {
		config := &SecurityConfig{
			TempFileRandomLength: 16,
			AllowedTempDirs:      []string{"/tmp"},
			FilePermissions:      0o600,
			SecureCleanup:        false, // Disabled
			MaxFileSize:          1024,
		}

		result := validator.ValidateSecurityConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have warnings for disabled secure cleanup")
	})
}

func TestRandomConfig_SecurityWarnings(t *testing.T) {
	// Use mock file system for faster, more reliable tests
	mockFS := &MockFileSystem{}
	validator := NewSecurityValidatorWithFS(mockFS)

	t.Run("low entropy warning", func(t *testing.T) {
		config := &RandomConfig{
			DefaultSuffixLength: 16,
			IDFormat:            "hex",
			MinEntropyBytes:     16, // Less than recommended 32
		}

		result := validator.ValidateRandomConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have warnings for low entropy")
	})

	t.Run("alphanumeric format warning", func(t *testing.T) {
		config := &RandomConfig{
			DefaultSuffixLength: 16,
			IDFormat:            "alphanumeric", // Lower entropy density
			MinEntropyBytes:     32,
		}

		result := validator.ValidateRandomConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have warnings for alphanumeric format")
	})

	t.Run("entropy check disabled warning", func(t *testing.T) {
		config := &RandomConfig{
			DefaultSuffixLength: 16,
			IDFormat:            "hex",
			MinEntropyBytes:     32,
			EnableEntropyCheck:  false, // Disabled
		}

		result := validator.ValidateRandomConfig(config)
		assert.True(t, result.Valid, "Should be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have warnings for disabled entropy check")
	})
}

func TestDefaultConfigs(t *testing.T) {
	// Use mock file system for faster, more reliable tests
	mockFS := &MockFileSystem{}
	validator := NewSecurityValidatorWithFS(mockFS)

	t.Run("default security config is valid", func(t *testing.T) {
		config := DefaultSecurityConfig()
		result := validator.ValidateSecurityConfig(config)
		assert.True(t, result.Valid, "Default security config should be valid")
	})

	t.Run("default random config is valid", func(t *testing.T) {
		config := DefaultRandomConfig()
		result := validator.ValidateRandomConfig(config)
		assert.True(t, result.Valid, "Default random config should be valid")
	})

	t.Run("default configs are compatible", func(t *testing.T) {
		secConfig := DefaultSecurityConfig()
		randConfig := DefaultRandomConfig()
		result := validator.ValidateCombinedConfig(secConfig, randConfig)
		assert.True(t, result.Valid, "Default configs should be compatible")
	})
}

func TestSecurityCustomValidationFunctions(t *testing.T) {
	t.Run("validateSecurePath", func(t *testing.T) {
		tests := []struct {
			path     string
			expected bool
		}{
			{"/tmp/safe", true},
			{"/var/tmp/safe", true},
			{"../dangerous", false},
			{"./relative", false},
			{"/tmp/file\x00null", false},
			{"", true}, // Empty is allowed
		}

		for _, tt := range tests {
			// Note: This is testing the logic, not the validator framework integration
			result := !containsDangerousPatterns(tt.path)
			assert.Equal(t, tt.expected, result, "Path validation for: %s", tt.path)
		}
	})
}

// Helper function to test path validation logic
func containsDangerousPatterns(path string) bool {
	if path == "" {
		return false
	}
	return strings.Contains(path, "..") ||
		strings.Contains(path, "./") ||
		strings.Contains(path, ".\\") ||
		strings.Contains(path, "\x00")
}

func TestSecurityErrors(t *testing.T) {
	t.Run("security errors are defined", func(t *testing.T) {
		securityErrors := []*SecurityOperationError{
			ErrInsufficientEntropy,
			ErrInvalidPath,
			ErrTempFileCreation,
			ErrAtomicWrite,
			ErrInsecurePermissions,
			ErrDangerousDirectory,
		}

		for _, err := range securityErrors {
			assert.NotNil(t, err, "Security error should be defined")
			assert.NotEmpty(t, err.Error(), "Security error should have message")
			assert.NotEmpty(t, err.Component, "Security error should have component")
			assert.NotEmpty(t, err.Operation, "Security error should have operation")
			assert.NotEmpty(t, err.Remediation, "Security error should have remediation")
			assert.True(t, IsSecurityError(err), "Should be identified as security error")
		}
	})
}

func TestSecurityValidationResult(t *testing.T) {
	t.Run("validation result structure", func(t *testing.T) {
		result := &SecurityValidationResult{
			Valid: true,
			Errors: []SecurityError{
				{
					Field:    "TestField",
					Tag:      "required",
					Value:    nil,
					Message:  "Test message",
					Severity: "error",
				},
			},
			Warnings: []SecurityWarning{
				{
					Field:   "TestField",
					Message: "Test warning",
					Impact:  "medium",
				},
			},
		}

		assert.True(t, result.Valid)
		assert.Len(t, result.Errors, 1)
		assert.Len(t, result.Warnings, 1)
		assert.Equal(t, "TestField", result.Errors[0].Field)
		assert.Equal(t, "error", result.Errors[0].Severity)
		assert.Equal(t, "medium", result.Warnings[0].Impact)
	})
}
