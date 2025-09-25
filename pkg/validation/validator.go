// Package validation provides a unified validation system for the entire application.
//
// This package consolidates all validation logic from various packages into a single,
// comprehensive validation system that can be used across the application.
package validation

import (
	"fmt"
	"net/mail"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ValidationLevel represents the severity of a validation issue
type ValidationLevel int

const (
	LevelInfo ValidationLevel = iota
	LevelWarning
	LevelError
	LevelFatal
)

// String returns the string representation of the validation level
func (l ValidationLevel) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// ValidationIssue represents a single validation issue
type ValidationIssue struct {
	Field       string                 `json:"field"`
	Level       ValidationLevel        `json:"level"`
	Message     string                 `json:"message"`
	Code        string                 `json:"code"`
	Value       interface{}            `json:"value,omitempty"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid   bool                         `json:"valid"`
	Issues  []ValidationIssue            `json:"issues"`
	Summary interfaces.ValidationSummary `json:"summary"`
}

// UnifiedValidator provides comprehensive validation for all application needs
type UnifiedValidator struct {
	patterns map[string]*regexp.Regexp
	rules    map[string]ValidationRule
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Name        string
	Validator   func(interface{}) error
	Message     string
	Suggestions []string
	Level       ValidationLevel
}

// NewUnifiedValidator creates a new unified validator with all common patterns
func NewUnifiedValidator() *UnifiedValidator {
	validator := &UnifiedValidator{
		patterns: make(map[string]*regexp.Regexp),
		rules:    make(map[string]ValidationRule),
	}

	validator.initializePatterns()
	validator.initializeRules()

	return validator
}

// initializePatterns sets up common regex patterns
func (v *UnifiedValidator) initializePatterns() {
	v.patterns = map[string]*regexp.Regexp{
		"project_name":    regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$`),
		"package_name":    regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`),
		"version":         regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?(?:\+([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?$`),
		"url":             regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})?(?:/[^\s]*)?$`),
		"github_repo":     regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`),
		"safe_string":     regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?()]+$`),
		"secure_password": regexp.MustCompile(`^.{8,}$`), // Simplified regex for Go compatibility
	}
}

// initializeRules sets up common validation rules
func (v *UnifiedValidator) initializeRules() {
	v.rules = map[string]ValidationRule{
		"required": {
			Name: "required",
			Validator: func(value interface{}) error {
				if value == nil {
					return fmt.Errorf("field is required")
				}
				if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
					return fmt.Errorf("field cannot be empty")
				}
				return nil
			},
			Message:     "This field is required",
			Suggestions: []string{"Please provide a value for this field"},
			Level:       LevelError,
		},
		"email": {
			Name: "email",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if _, err := mail.ParseAddress(str); err != nil {
						return fmt.Errorf("invalid email format")
					}
				}
				return nil
			},
			Message:     "Please enter a valid email address",
			Suggestions: []string{"example@domain.com"},
			Level:       LevelError,
		},
		"url": {
			Name: "url",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if _, err := url.ParseRequestURI(str); err != nil {
						return fmt.Errorf("invalid URL format")
					}
				}
				return nil
			},
			Message:     "Please enter a valid URL",
			Suggestions: []string{"https://example.com"},
			Level:       LevelError,
		},
		"project_name": {
			Name: "project_name",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !v.patterns["project_name"].MatchString(str) {
						return fmt.Errorf("invalid project name format")
					}
				}
				return nil
			},
			Message:     "Project name must contain only alphanumeric characters, hyphens, and underscores",
			Suggestions: []string{"my-awesome-project", "project_name"},
			Level:       LevelError,
		},
		"package_name": {
			Name: "package_name",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !v.patterns["package_name"].MatchString(str) {
						return fmt.Errorf("invalid package name format")
					}
				}
				return nil
			},
			Message:     "Package name must be lowercase with hyphens only",
			Suggestions: []string{"my-package", "awesome-package"},
			Level:       LevelError,
		},
		"version": {
			Name: "version",
			Validator: func(value interface{}) error {
				if str, ok := value.(string); ok && str != "" {
					if !v.patterns["version"].MatchString(str) {
						return fmt.Errorf("invalid version format")
					}
				}
				return nil
			},
			Message:     "Version must follow semantic versioning (e.g., 1.0.0)",
			Suggestions: []string{"1.0.0", "v1.2.3", "2.0.0-beta.1"},
			Level:       LevelError,
		},
		"min_length": {
			Name: "min_length",
			Validator: func(value interface{}) error {
				// This will be customized per use case
				return nil
			},
			Message:     "Value is too short",
			Suggestions: []string{"Please provide a longer value"},
			Level:       LevelError,
		},
		"max_length": {
			Name: "max_length",
			Validator: func(value interface{}) error {
				// This will be customized per use case
				return nil
			},
			Message:     "Value is too long",
			Suggestions: []string{"Please provide a shorter value"},
			Level:       LevelError,
		},
	}
}

// ValidateField validates a single field with specified rules
func (v *UnifiedValidator) ValidateField(field string, value interface{}, rules ...string) []ValidationIssue {
	var issues []ValidationIssue

	for _, ruleName := range rules {
		if rule, exists := v.rules[ruleName]; exists {
			if err := rule.Validator(value); err != nil {
				issue := ValidationIssue{
					Field:       field,
					Level:       rule.Level,
					Message:     rule.Message,
					Code:        ruleName,
					Value:       value,
					Suggestions: rule.Suggestions,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// ValidateString validates a string with length constraints
func (v *UnifiedValidator) ValidateString(field, value string, minLen, maxLen int, rules ...string) []ValidationIssue {
	var issues []ValidationIssue

	// Add length validation rules
	if minLen > 0 {
		if len(value) < minLen {
			issues = append(issues, ValidationIssue{
				Field:       field,
				Level:       LevelError,
				Message:     fmt.Sprintf("Value must be at least %d characters long", minLen),
				Code:        "min_length",
				Value:       value,
				Suggestions: []string{fmt.Sprintf("Please provide at least %d characters", minLen)},
			})
		}
	}

	if maxLen > 0 {
		if len(value) > maxLen {
			issues = append(issues, ValidationIssue{
				Field:       field,
				Level:       LevelError,
				Message:     fmt.Sprintf("Value must be no more than %d characters long", maxLen),
				Code:        "max_length",
				Value:       value,
				Suggestions: []string{fmt.Sprintf("Please provide no more than %d characters", maxLen)},
			})
		}
	}

	// Add other validation rules
	fieldIssues := v.ValidateField(field, value, rules...)
	issues = append(issues, fieldIssues...)

	return issues
}

// ValidateProjectConfig validates a complete project configuration
func (v *UnifiedValidator) ValidateProjectConfig(config *models.ProjectConfig) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Issues: []ValidationIssue{},
	}

	if config == nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Field:   "config",
			Level:   LevelFatal,
			Message: "Configuration cannot be null",
			Code:    "null_config",
		})
		result.Valid = false
		return result
	}

	// Validate basic fields
	result.Issues = append(result.Issues, v.ValidateString("name", config.Name, 1, 100, "required", "project_name")...)
	result.Issues = append(result.Issues, v.ValidateString("organization", config.Organization, 1, 100, "required", "package_name")...)
	result.Issues = append(result.Issues, v.ValidateString("description", config.Description, 0, 500)...)
	result.Issues = append(result.Issues, v.ValidateString("author", config.Author, 1, 100, "required")...)
	result.Issues = append(result.Issues, v.ValidateString("email", config.Email, 0, 100, "email")...)
	result.Issues = append(result.Issues, v.ValidateString("repository", config.Repository, 0, 200, "url")...)

	// Validate output path
	if config.OutputPath != "" {
		if !filepath.IsAbs(config.OutputPath) {
			result.Issues = append(result.Issues, ValidationIssue{
				Field:       "output_path",
				Level:       LevelWarning,
				Message:     "Output path should be absolute",
				Code:        "relative_path",
				Suggestions: []string{"Use an absolute path for better reliability"},
			})
		}
	}

	// Calculate summary
	v.calculateSummary(result)
	result.Valid = result.Summary.ErrorCount == 0

	return result
}

// calculateSummary calculates validation summary statistics
func (v *UnifiedValidator) calculateSummary(result *ValidationResult) {
	result.Summary.TotalFiles = 1 // Single config validation
	result.Summary.ValidFiles = 0
	if result.Valid {
		result.Summary.ValidFiles = 1
	}

	for _, issue := range result.Issues {
		switch issue.Level {
		case LevelError, LevelFatal:
			result.Summary.ErrorCount++
		case LevelWarning:
			result.Summary.WarningCount++
		}
	}
}

// AddCustomRule adds a custom validation rule
func (v *UnifiedValidator) AddCustomRule(name string, rule ValidationRule) {
	v.rules[name] = rule
}

// GetPattern returns a compiled regex pattern by name
func (v *UnifiedValidator) GetPattern(name string) *regexp.Regexp {
	return v.patterns[name]
}

// ValidateEmail validates an email address
func (v *UnifiedValidator) ValidateEmail(email string) error {
	if email == "" {
		return nil // Empty email is allowed
	}

	_, err := mail.ParseAddress(email)
	return err
}

// ValidateURL validates a URL
func (v *UnifiedValidator) ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil // Empty URL is allowed
	}

	_, err := url.ParseRequestURI(urlStr)
	return err
}

// ValidateNonEmptyString validates that a string is not empty
func (v *UnifiedValidator) ValidateNonEmptyString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidatePath validates a file path
func (v *UnifiedValidator) ValidatePath(path string, allowedBasePaths ...string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Check if path is within allowed base paths
	if len(allowedBasePaths) > 0 {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("invalid path: %v", err)
		}

		for _, basePath := range allowedBasePaths {
			absBasePath, err := filepath.Abs(basePath)
			if err != nil {
				continue
			}

			if strings.HasPrefix(absPath, absBasePath) {
				return nil
			}
		}

		return fmt.Errorf("path must be within allowed base paths")
	}

	return nil
}

// ValidatePassword validates a password for security requirements
func (v *UnifiedValidator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateVersion validates a semantic version string
func (v *UnifiedValidator) ValidateVersion(version string) error {
	if version == "" {
		return nil
	}

	pattern := v.patterns["version"]
	if !pattern.MatchString(version) {
		return fmt.Errorf("invalid version format, expected semantic versioning (e.g., 1.0.0)")
	}

	return nil
}

// ValidateSlice validates that a slice is not empty
func (v *UnifiedValidator) ValidateSlice(slice interface{}, fieldName string) error {
	if slice == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}

	// Use reflection to check if it's a slice and has elements
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("%s must be a slice", fieldName)
	}

	if sliceValue.Len() == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	return nil
}
