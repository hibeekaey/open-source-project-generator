package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// ValidationError represents a single validation error with context
type ValidationError struct {
	Field   string
	Message string
	Code    string
	Value   interface{}
}

// NewValidationErrorStruct creates a new validation error struct
func NewValidationErrorStruct(field, message, code string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
		Value:   value,
	}
}

// Validator provides comprehensive input validation
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message, code string, value interface{}) {
	v.errors = append(v.errors, NewValidationErrorStruct(field, message, code, value))
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// GetResult returns the validation result
func (v *Validator) GetResult() ValidationResult {
	return ValidationResult{
		IsValid: !v.HasErrors(),
		Errors:  v.errors,
	}
}

// Clear clears all validation errors
func (v *Validator) Clear() {
	v.errors = make([]ValidationError, 0)
}

// Convenience functions

// ValidateNonEmptyString validates that a string is not empty after trimming whitespace
func ValidateNonEmptyString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return nil // Empty email is allowed
	}

	// Basic email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateURL validates a URL
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil // Empty URL is allowed
	}

	// Basic URL validation
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// ValidateProjectName validates a project name
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}

	// Check for valid characters: alphanumeric, hyphens, underscores
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') && char != '-' && char != '_' {
			return fmt.Errorf("project name must contain only alphanumeric characters, hyphens, and underscores")
		}
	}

	// Must start and end with alphanumeric
	first := name[0]
	last := name[len(name)-1]
	if (first < 'a' || first > 'z') && (first < 'A' || first > 'Z') && (first < '0' || first > '9') {
		return fmt.Errorf("project name must start with alphanumeric character")
	}
	if (last < 'a' || last > 'z') && (last < 'A' || last > 'Z') && (last < '0' || last > '9') {
		return fmt.Errorf("project name must end with alphanumeric character")
	}

	return nil
}

// ValidatePackageName validates a package name
func ValidatePackageName(name string) error {
	if name == "" {
		return fmt.Errorf("package name is required")
	}

	// Check for valid characters: lowercase letters, numbers, hyphens
	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
			return fmt.Errorf("package name must contain only lowercase letters, numbers, and hyphens")
		}
	}

	// Must start and end with alphanumeric
	first := name[0]
	last := name[len(name)-1]
	if (first < 'a' || first > 'z') && (first < '0' || first > '9') {
		return fmt.Errorf("package name must start with alphanumeric character")
	}
	if (last < 'a' || last > 'z') && (last < '0' || last > '9') {
		return fmt.Errorf("package name must end with alphanumeric character")
	}

	return nil
}

// ValidateVersion validates a semantic version string
func ValidateVersion(version string) error {
	if version == "" {
		return nil
	}

	// Basic semantic version validation
	if !strings.Contains(version, ".") {
		return fmt.Errorf("invalid version format, expected semantic versioning (e.g., 1.0.0)")
	}

	return nil
}

// ValidateNonEmptySlice validates that a slice is not empty
func ValidateNonEmptySlice(slice interface{}, fieldName string) error {
	if slice == nil {
		return fmt.Errorf("%s is required", fieldName)
	}

	// Use reflection to check if it's a slice and if it's empty
	switch v := slice.(type) {
	case []string:
		if len(v) == 0 {
			return fmt.Errorf("%s cannot be empty", fieldName)
		}
	case []int:
		if len(v) == 0 {
			return fmt.Errorf("%s cannot be empty", fieldName)
		}
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("%s cannot be empty", fieldName)
		}
	default:
		// For non-slice types, return an error
		return fmt.Errorf("%s must be a slice", fieldName)
	}

	return nil
}

// ValidateNotNil validates that a value is not nil
func ValidateNotNil(value interface{}, fieldName string) error {
	if value == nil {
		return fmt.Errorf("%s is required", fieldName)
	}

	// Use reflection to check for nil pointers
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}

// Utility functions for error formatting

// FormatValidationErrors formats validation errors for user display
func FormatValidationErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, fmt.Sprintf("• %s", err.Message))
	}

	return fmt.Sprintf("Validation failed:\n%s", strings.Join(messages, "\n"))
}

// GetValidationErrorsByField groups validation errors by field
func GetValidationErrorsByField(errors []ValidationError) map[string][]ValidationError {
	result := make(map[string][]ValidationError)

	for _, err := range errors {
		result[err.Field] = append(result[err.Field], err)
	}

	return result
}
