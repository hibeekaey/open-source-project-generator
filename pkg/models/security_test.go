package models

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityConfig_Validation(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name        string
		config      *SecurityConfig
		expectValid bool
		expectError string
	}{
		{
			name: "valid security config",
			config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp", "/var/tmp"},
				FilePermissions:      0o600,
				EnablePathValidation: true,
				MaxFileSize:          10 * 1024 * 1024,
				SecureCleanup:        true,
			},
			expectValid: true,
		},
		{
			name:   "missing required fields",
			config: &SecurityConfig{
				// Missing required fields
			},
			expectValid: false,
			expectError: "TempFileRandomLength",
		},
		{
			name: "temp file random length too small",
			config: &SecurityConfig{
				TempFileRandomLength: 4, // Too small
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			expectValid: false,
			expectError: "min",
		},
		{
			name: "temp file random length too large",
			config: &SecurityConfig{
				TempFileRandomLength: 128, // Too large
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			expectValid: false,
			expectError: "max",
		},
		{
			name: "empty allowed temp dirs",
			config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{}, // Empty
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			expectValid: false,
			expectError: "min",
		},
		{
			name: "max file size too small",
			config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          512, // Too small
			},
			expectValid: false,
			expectError: "min",
		},
		{
			name: "dangerous directory path",
			config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/etc"}, // Dangerous
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			expectValid: false,
			expectError: "dangerous_path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateSecurityConfig(tt.config)

			assert.Equal(t, tt.expectValid, result.Valid, "Validation result should match expected")

			if !tt.expectValid && tt.expectError != "" {
				assert.NotEmpty(t, result.Errors, "Should have validation errors")

				// Check if expected error is present
				found := false
				for _, err := range result.Errors {
					if err.Tag == tt.expectError || err.Field == tt.expectError {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error '%s' not found in validation errors", tt.expectError)
			}
		})
	}
}

func TestRandomConfig_Validation(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name        string
		config      *RandomConfig
		expectValid bool
		expectError string
	}{
		{
			name: "valid random config",
			config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     32,
				IDPrefixLength:      4,
				EnableEntropyCheck:  true,
			},
			expectValid: true,
		},
		{
			name:   "missing required fields",
			config: &RandomConfig{
				// Missing required fields
			},
			expectValid: false,
			expectError: "DefaultSuffixLength",
		},
		{
			name: "invalid ID format",
			config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "invalid", // Invalid format
				MinEntropyBytes:     32,
			},
			expectValid: false,
			expectError: "oneof",
		},
		{
			name: "suffix length too small",
			config: &RandomConfig{
				DefaultSuffixLength: 4, // Too small
				IDFormat:            "hex",
				MinEntropyBytes:     32,
			},
			expectValid: false,
			expectError: "min",
		},
		{
			name: "suffix length too large",
			config: &RandomConfig{
				DefaultSuffixLength: 128, // Too large
				IDFormat:            "hex",
				MinEntropyBytes:     32,
			},
			expectValid: false,
			expectError: "max",
		},
		{
			name: "entropy bytes too small",
			config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     8, // Too small
			},
			expectValid: false,
			expectError: "min",
		},
		{
			name: "entropy bytes too large",
			config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     2048, // Too large
			},
			expectValid: false,
			expectError: "max",
		},
		{
			name: "prefix length too large",
			config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     32,
				IDPrefixLength:      32, // Too large
			},
			expectValid: false,
			expectError: "max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateRandomConfig(tt.config)

			assert.Equal(t, tt.expectValid, result.Valid, "Validation result should match expected")

			if !tt.expectValid && tt.expectError != "" {
				assert.NotEmpty(t, result.Errors, "Should have validation errors")

				// Check if expected error is present
				found := false
				for _, err := range result.Errors {
					if err.Tag == tt.expectError || err.Field == tt.expectError {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected error '%s' not found in validation errors", tt.expectError)
			}
		})
	}
}

func TestCombinedConfig_Validation(t *testing.T) {
	validator := NewSecurityValidator()

	secConfig := &SecurityConfig{
		TempFileRandomLength: 16,
		AllowedTempDirs:      []string{"/tmp"},
		FilePermissions:      0o600,
		EnablePathValidation: true,
		MaxFileSize:          1024,
		SecureCleanup:        true,
	}

	randConfig := &RandomConfig{
		DefaultSuffixLength: 16, // Matches security config
		IDFormat:            "hex",
		MinEntropyBytes:     32,
		IDPrefixLength:      4,
		EnableEntropyCheck:  true,
	}

	t.Run("compatible configs", func(t *testing.T) {
		result := validator.ValidateCombinedConfig(secConfig, randConfig)
		assert.True(t, result.Valid, "Compatible configs should be valid")
	})

	t.Run("mismatched random lengths", func(t *testing.T) {
		mismatchedRandConfig := &RandomConfig{
			DefaultSuffixLength: 12, // Different from security config
			IDFormat:            "hex",
			MinEntropyBytes:     32,
			IDPrefixLength:      4,
			EnableEntropyCheck:  true,
		}

		result := validator.ValidateCombinedConfig(secConfig, mismatchedRandConfig)
		assert.True(t, result.Valid, "Should still be valid but with warnings")
		assert.NotEmpty(t, result.Warnings, "Should have compatibility warnings")
	})
}

func TestSecurityValidator_FilePermissions(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name        string
		permissions os.FileMode
		expectValid bool
		expectWarn  bool
	}{
		{
			name:        "secure permissions (600)",
			permissions: 0o600,
			expectValid: true,
			expectWarn:  false,
		},
		{
			name:        "secure permissions (644)",
			permissions: 0o644,
			expectValid: true,
			expectWarn:  true, // Warning for group/other read
		},
		{
			name:        "insecure permissions (666)",
			permissions: 0o666,
			expectValid: true,
			expectWarn:  true, // Warning for group/other write
		},
		{
			name:        "no owner write (400)",
			permissions: 0o400,
			expectValid: false, // Error for no owner write
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      tt.permissions,
				MaxFileSize:          1024,
			}

			result := validator.ValidateSecurityConfig(config)

			assert.Equal(t, tt.expectValid, result.Valid, "Validation result should match expected")

			if tt.expectWarn {
				assert.NotEmpty(t, result.Warnings, "Should have warnings for insecure permissions")
			}
		})
	}
}

func TestSecurityValidator_DangerousDirectories(t *testing.T) {
	validator := NewSecurityValidator()

	dangerousDirs := []string{
		"/",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/etc",
		"/root",
		"/boot",
	}

	for _, dir := range dangerousDirs {
		t.Run("dangerous_dir_"+dir, func(t *testing.T) {
			config := &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{dir},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			}

			result := validator.ValidateSecurityConfig(config)
			assert.False(t, result.Valid, "Dangerous directory should be invalid: %s", dir)
			assert.NotEmpty(t, result.Errors, "Should have errors for dangerous directory")
		})
	}
}

func TestSecurityValidator_SecurityWarnings(t *testing.T) {
	validator := NewSecurityValidator()

	t.Run("weak random length warning", func(t *testing.T) {
		config := &SecurityConfig{
			TempFileRandomLength: 12, // Less than recommended 16
			AllowedTempDirs:      []string{"/tmp"},
			FilePermissions:      0o600,
			MaxFileSize:          1024,
		}

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
	validator := NewSecurityValidator()

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
	validator := NewSecurityValidator()

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
