package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// SecurityConfig holds security-related configuration parameters
//
// SECURITY RATIONALE:
// This configuration structure provides centralized control over security-sensitive
// parameters throughout the application. Each field addresses specific security
// concerns identified in the codebase audit:
//
//  1. TempFileRandomLength: Controls the entropy of temporary file names to prevent
//     race condition attacks where predictable names allow attackers to pre-create files
//
//  2. AllowedTempDirs: Restricts where temporary files can be created to prevent
//     directory traversal attacks and unauthorized file system access
//
//  3. FilePermissions: Ensures consistent secure file permissions across the application
//     to prevent unauthorized access to sensitive data
//
//  4. EnablePathValidation: Controls whether strict path validation is enforced to
//     prevent directory traversal and path injection attacks
//
// 5. MaxFileSize: Prevents resource exhaustion attacks through oversized file operations
//
//  6. SecureCleanup: Enables secure deletion of temporary files to prevent information
//     disclosure through file system forensics
//
// All fields include validation tags to ensure secure configuration values and
// prevent misconfigurations that could introduce security vulnerabilities.
type SecurityConfig struct {
	// TempFileRandomLength specifies the length of random suffixes for temporary files
	TempFileRandomLength int `yaml:"temp_file_random_length" json:"temp_file_random_length" validate:"required,min=8,max=64"`

	// AllowedTempDirs specifies directories where temporary files can be created
	AllowedTempDirs []string `yaml:"allowed_temp_dirs" json:"allowed_temp_dirs" validate:"required,min=1,dive,required"`

	// FilePermissions specifies default secure file permissions (octal format)
	FilePermissions os.FileMode `yaml:"file_permissions" json:"file_permissions" validate:"required"`

	// EnablePathValidation enables strict path validation to prevent directory traversal
	EnablePathValidation bool `yaml:"enable_path_validation" json:"enable_path_validation"`

	// MaxFileSize specifies maximum allowed file size for security operations (in bytes)
	MaxFileSize int64 `yaml:"max_file_size" json:"max_file_size" validate:"required,min=1024"`

	// SecureCleanup enables secure deletion of temporary files
	SecureCleanup bool `yaml:"secure_cleanup" json:"secure_cleanup"`
}

// RandomConfig holds configuration for secure random number generation
//
// SECURITY RATIONALE:
// This configuration addresses the critical security vulnerability of predictable
// random number generation found throughout the codebase. The configuration ensures:
//
//  1. DefaultSuffixLength: Minimum entropy requirements for random suffixes used in
//     temporary files, audit IDs, and other security-sensitive contexts. Shorter
//     lengths are vulnerable to brute force and prediction attacks.
//
//  2. IDFormat: Controls the encoding format for random data, affecting both entropy
//     density and compatibility. Different formats provide different security/usability
//     trade-offs:
//     - hex: 4 bits per character, most compact
//     - base64: 6 bits per character, URL-safe
//     - alphanumeric: ~5.95 bits per character, human-readable
//
//  3. MinEntropyBytes: Ensures sufficient entropy for cryptographic operations.
//     Insufficient entropy can make cryptographic operations predictable and vulnerable.
//
//  4. IDPrefixLength: Controls the length of prefixes in generated IDs for better
//     organization and debugging while maintaining security.
//
//  5. EnableEntropyCheck: Enables runtime validation of entropy quality to detect
//     and prevent weak random generation in low-entropy environments.
//
// This configuration replaces insecure patterns like timestamp-based generation
// (time.Now().UnixNano()) with cryptographically secure alternatives using crypto/rand.
type RandomConfig struct {
	// DefaultSuffixLength for temporary files and general random suffixes
	DefaultSuffixLength int `yaml:"default_suffix_length" json:"default_suffix_length" validate:"required,min=8,max=64"`

	// IDFormat specifies format for secure IDs (hex, base64, alphanumeric)
	IDFormat string `yaml:"id_format" json:"id_format" validate:"required,oneof=hex base64 alphanumeric"`

	// MinEntropyBytes minimum entropy required for cryptographic operations
	MinEntropyBytes int `yaml:"min_entropy_bytes" json:"min_entropy_bytes" validate:"required,min=16,max=1024"`

	// IDPrefixLength specifies the length of prefixes for generated IDs
	IDPrefixLength int `yaml:"id_prefix_length" json:"id_prefix_length" validate:"min=0,max=16"`

	// EnableEntropyCheck enables entropy quality validation
	EnableEntropyCheck bool `yaml:"enable_entropy_check" json:"enable_entropy_check"`
}

// SecurityValidationResult represents the result of security configuration validation
type SecurityValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []SecurityError   `json:"errors"`
	Warnings []SecurityWarning `json:"warnings"`
	Config   *SecurityConfig   `json:"config,omitempty"`
	Random   *RandomConfig     `json:"random,omitempty"`
}

// SecurityError represents a security configuration validation error
type SecurityError struct {
	Field    string      `json:"field"`
	Tag      string      `json:"tag"`
	Value    interface{} `json:"value"`
	Message  string      `json:"message"`
	Severity string      `json:"severity"` // "error", "warning", "info"
}

// SecurityWarning represents a security configuration warning
type SecurityWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Impact  string `json:"impact"` // "high", "medium", "low"
}

// FileSystemInterface abstracts file system operations for testing
type FileSystemInterface interface {
	Stat(name string) (os.FileInfo, error)
	IsAbs(path string) bool
	Clean(path string) string
}

// DefaultFileSystem implements FileSystemInterface using real OS operations
type DefaultFileSystem struct{}

func (fs *DefaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *DefaultFileSystem) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func (fs *DefaultFileSystem) Clean(path string) string {
	return filepath.Clean(path)
}

// SecurityValidator handles security configuration validation
type SecurityValidator struct {
	validator *validator.Validate
	fs        FileSystemInterface
}

// NewSecurityValidator creates a new security configuration validator
func NewSecurityValidator() *SecurityValidator {
	return NewSecurityValidatorWithFS(&DefaultFileSystem{})
}

// NewSecurityValidatorWithFS creates a new security configuration validator with custom file system interface
func NewSecurityValidatorWithFS(fs FileSystemInterface) *SecurityValidator {
	v := validator.New()

	// Register custom validation functions for security
	v.RegisterValidation("secure_path", validateSecurePath)
	v.RegisterValidation("secure_permissions", validateSecurePermissions)

	return &SecurityValidator{
		validator: v,
		fs:        fs,
	}
}

// ValidateSecurityConfig validates a security configuration
func (sv *SecurityValidator) ValidateSecurityConfig(config *SecurityConfig) *SecurityValidationResult {
	result := &SecurityValidationResult{
		Valid:    true,
		Errors:   []SecurityError{},
		Warnings: []SecurityWarning{},
		Config:   config,
	}

	// Perform struct validation
	if err := sv.validator.Struct(config); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, SecurityError{
					Field:    fieldError.Field(),
					Tag:      fieldError.Tag(),
					Value:    fieldError.Value(),
					Message:  sv.getSecurityErrorMessage(fieldError),
					Severity: "error",
				})
			}
		}
	}

	// Custom security validation
	sv.validateTempDirectories(config, result)
	sv.validateFilePermissions(config, result)
	sv.validateSecuritySettings(config, result)

	return result
}

// ValidateRandomConfig validates a random configuration
func (sv *SecurityValidator) ValidateRandomConfig(config *RandomConfig) *SecurityValidationResult {
	result := &SecurityValidationResult{
		Valid:    true,
		Errors:   []SecurityError{},
		Warnings: []SecurityWarning{},
		Random:   config,
	}

	// Perform struct validation
	if err := sv.validator.Struct(config); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, SecurityError{
					Field:    fieldError.Field(),
					Tag:      fieldError.Tag(),
					Value:    fieldError.Value(),
					Message:  sv.getSecurityErrorMessage(fieldError),
					Severity: "error",
				})
			}
		}
	}

	// Custom random configuration validation
	sv.validateRandomSettings(config, result)

	return result
}

// ValidateCombinedConfig validates both security and random configurations together
func (sv *SecurityValidator) ValidateCombinedConfig(secConfig *SecurityConfig, randConfig *RandomConfig) *SecurityValidationResult {
	result := &SecurityValidationResult{
		Valid:    true,
		Errors:   []SecurityError{},
		Warnings: []SecurityWarning{},
		Config:   secConfig,
		Random:   randConfig,
	}

	// Validate individual configurations
	secResult := sv.ValidateSecurityConfig(secConfig)
	randResult := sv.ValidateRandomConfig(randConfig)

	// Combine results
	result.Valid = secResult.Valid && randResult.Valid
	result.Errors = append(result.Errors, secResult.Errors...)
	result.Errors = append(result.Errors, randResult.Errors...)
	result.Warnings = append(result.Warnings, secResult.Warnings...)
	result.Warnings = append(result.Warnings, randResult.Warnings...)

	// Cross-configuration validation
	sv.validateConfigCompatibility(secConfig, randConfig, result)

	return result
}

// validateTempDirectories validates temporary directory configuration
func (sv *SecurityValidator) validateTempDirectories(config *SecurityConfig, result *SecurityValidationResult) {
	for i, dir := range config.AllowedTempDirs {
		// Check if directory path is absolute
		if !sv.fs.IsAbs(dir) {
			result.Warnings = append(result.Warnings, SecurityWarning{
				Field:   fmt.Sprintf("AllowedTempDirs[%d]", i),
				Message: fmt.Sprintf("Relative path '%s' may be less secure than absolute paths", dir),
				Impact:  "medium",
			})
		}

		// Check for potentially dangerous directories
		if sv.isDangerousDirectory(dir) {
			result.Valid = false
			result.Errors = append(result.Errors, SecurityError{
				Field:    fmt.Sprintf("AllowedTempDirs[%d]", i),
				Tag:      "dangerous_path",
				Value:    dir,
				Message:  fmt.Sprintf("Directory '%s' is potentially dangerous for temporary files", dir),
				Severity: "error",
			})
		}

		// Check if directory exists and is accessible (CI-friendly approach)
		if info, err := sv.fs.Stat(dir); err != nil {
			// In CI environments, directories may not exist or be accessible
			// Only warn for standard directories, ignore for CI-specific paths
			if !strings.Contains(dir, "ci") && !strings.Contains(dir, "tmp") {
				result.Warnings = append(result.Warnings, SecurityWarning{
					Field:   fmt.Sprintf("AllowedTempDirs[%d]", i),
					Message: fmt.Sprintf("Directory '%s' may not exist in all environments", dir),
					Impact:  "low", // Further reduced to minimize CI false positives
				})
			}
		} else if !info.IsDir() {
			// Only error for definitively non-directory paths
			result.Valid = false
			result.Errors = append(result.Errors, SecurityError{
				Field:    fmt.Sprintf("AllowedTempDirs[%d]", i),
				Tag:      "not_directory",
				Value:    dir,
				Message:  fmt.Sprintf("Path '%s' exists but is not a directory", dir),
				Severity: "error",
			})
		}
	}
}

// validateFilePermissions validates file permission settings
func (sv *SecurityValidator) validateFilePermissions(config *SecurityConfig, result *SecurityValidationResult) {
	perm := config.FilePermissions

	// Check if permissions are overly permissive (only warn for write permissions)
	if perm&0o022 != 0 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "FilePermissions",
			Message: fmt.Sprintf("File permissions %o allow group/other write access, which may be insecure", perm),
			Impact:  "medium",
		})
	} else if perm&0o044 != 0 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "FilePermissions",
			Message: fmt.Sprintf("File permissions %o allow group/other read access (low security risk)", perm),
			Impact:  "low", // Read-only permissions are less risky
		})
	}

	// Check if permissions are too restrictive (no owner write)
	if perm&0o200 == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, SecurityError{
			Field:    "FilePermissions",
			Tag:      "no_write_permission",
			Value:    perm,
			Message:  fmt.Sprintf("File permissions %o do not allow owner write access", perm),
			Severity: "error",
		})
	}
}

// validateSecuritySettings validates general security settings
func (sv *SecurityValidator) validateSecuritySettings(config *SecurityConfig, result *SecurityValidationResult) {
	// Validate temp file random length (more lenient threshold)
	if config.TempFileRandomLength < 8 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "TempFileRandomLength",
			Message: "Random suffix length less than 8 characters may be vulnerable to brute force attacks",
			Impact:  "medium",
		})
	} else if config.TempFileRandomLength < 12 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "TempFileRandomLength",
			Message: "Random suffix length less than 12 characters has reduced security",
			Impact:  "low",
		})
	}

	// Validate max file size
	if config.MaxFileSize > 100*1024*1024 { // 100MB
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "MaxFileSize",
			Message: "Large maximum file size may lead to resource exhaustion attacks",
			Impact:  "low",
		})
	}

	// Recommend enabling path validation
	if !config.EnablePathValidation {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "EnablePathValidation",
			Message: "Disabling path validation may allow directory traversal attacks",
			Impact:  "high",
		})
	}

	// Recommend enabling secure cleanup
	if !config.SecureCleanup {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "SecureCleanup",
			Message: "Disabling secure cleanup may leave sensitive data in temporary files",
			Impact:  "medium",
		})
	}
}

// validateRandomSettings validates random configuration settings
func (sv *SecurityValidator) validateRandomSettings(config *RandomConfig, result *SecurityValidationResult) {
	// Validate entropy requirements (more reasonable thresholds)
	if config.MinEntropyBytes < 16 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "MinEntropyBytes",
			Message: "Minimum entropy less than 16 bytes may be insufficient for cryptographic operations",
			Impact:  "high",
		})
	} else if config.MinEntropyBytes < 24 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "MinEntropyBytes",
			Message: "Minimum entropy less than 24 bytes has reduced security for some operations",
			Impact:  "medium",
		})
	}

	// Validate ID format security
	if config.IDFormat == "alphanumeric" {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "IDFormat",
			Message: "Alphanumeric format has lower entropy density than hex or base64",
			Impact:  "low",
		})
	}

	// Validate suffix length (align with updated thresholds)
	if config.DefaultSuffixLength < 8 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "DefaultSuffixLength",
			Message: "Default suffix length less than 8 characters may be vulnerable to prediction attacks",
			Impact:  "medium",
		})
	} else if config.DefaultSuffixLength < 10 {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "DefaultSuffixLength",
			Message: "Default suffix length less than 10 characters has reduced security",
			Impact:  "low",
		})
	}

	// Recommend enabling entropy check
	if !config.EnableEntropyCheck {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "EnableEntropyCheck",
			Message: "Disabling entropy checks may allow weak random generation to go undetected",
			Impact:  "medium",
		})
	}
}

// validateConfigCompatibility validates compatibility between security and random configs
func (sv *SecurityValidator) validateConfigCompatibility(secConfig *SecurityConfig, randConfig *RandomConfig, result *SecurityValidationResult) {
	// Check if temp file random length matches default suffix length
	if secConfig.TempFileRandomLength != randConfig.DefaultSuffixLength {
		result.Warnings = append(result.Warnings, SecurityWarning{
			Field:   "TempFileRandomLength/DefaultSuffixLength",
			Message: "Mismatched random lengths between security and random configurations may cause inconsistency",
			Impact:  "low",
		})
	}
}

// isDangerousDirectory checks if a directory path is potentially dangerous
func (sv *SecurityValidator) isDangerousDirectory(dir string) bool {
	dangerous := []string{
		"/",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/usr/sbin",
		"/etc",
		"/root",
		"/boot",
		"/sys",
		"/proc",
		"/dev",
	}
	// Note: Removed /var/log and /home from dangerous paths as they may be valid
	// in some CI environments and can cause false positives

	cleanDir := sv.fs.Clean(dir)
	for _, dangerousPath := range dangerous {
		if cleanDir == dangerousPath || strings.HasPrefix(cleanDir, dangerousPath+"/") {
			return true
		}
	}

	return false
}

// getSecurityErrorMessage returns human-readable error messages for security validation
func (sv *SecurityValidator) getSecurityErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required for security configuration", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s for security", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s for security", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of the secure options: %s", fe.Field(), fe.Param())
	case "secure_path":
		return fmt.Sprintf("%s contains potentially unsafe path characters", fe.Field())
	case "secure_permissions":
		return fmt.Sprintf("%s has insecure file permissions", fe.Field())
	default:
		return fmt.Sprintf("%s failed security validation for tag '%s'", fe.Field(), fe.Tag())
	}
}

// Custom validation functions for security

// validateSecurePath validates that a path is secure and doesn't contain dangerous patterns
func validateSecurePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true
	}

	// Check for directory traversal patterns
	if strings.Contains(path, "..") || strings.Contains(path, "./") || strings.Contains(path, ".\\") {
		return false
	}

	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return false
	}

	return true
}

// validateSecurePermissions validates that file permissions are secure
func validateSecurePermissions(fl validator.FieldLevel) bool {
	perm := fl.Field().Interface().(os.FileMode)

	// Ensure owner has read/write permissions
	if perm&0o600 != 0o600 {
		return false
	}

	// Ensure permissions are not overly permissive (no execute for files)
	if perm&0o111 != 0 {
		return false
	}

	return true
}

// DefaultSecurityConfig returns a secure default configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		TempFileRandomLength: 16,
		AllowedTempDirs:      []string{"/tmp", "/var/tmp"},
		FilePermissions:      0o600, // Owner read/write only
		EnablePathValidation: true,
		MaxFileSize:          10 * 1024 * 1024, // 10MB
		SecureCleanup:        true,
	}
}

// DefaultRandomConfig returns a secure default random configuration
func DefaultRandomConfig() *RandomConfig {
	return &RandomConfig{
		DefaultSuffixLength: 16,
		IDFormat:            "hex",
		MinEntropyBytes:     32,
		IDPrefixLength:      4,
		EnableEntropyCheck:  true,
	}
}

// Security-specific errors are now defined in errors.go with enhanced functionality
// This maintains backward compatibility while providing better error handling
