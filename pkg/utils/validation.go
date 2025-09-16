package utils

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
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

// String validation methods

// ValidateNonEmptyString validates that a string is not empty after trimming whitespace
func ValidateNonEmptyString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateStringLength validates string length constraints
func (v *Validator) ValidateStringLength(value, field string, minLength, maxLength int) {
	trimmed := strings.TrimSpace(value)
	length := len(trimmed)

	if minLength > 0 && length < minLength {
		v.AddError(field, fmt.Sprintf("%s must be at least %d characters long", field, minLength), "min_length", value)
	}

	if maxLength > 0 && length > maxLength {
		v.AddError(field, fmt.Sprintf("%s must be no more than %d characters long", field, maxLength), "max_length", value)
	}
}

// ValidateStringPattern validates string against a regex pattern
func (v *Validator) ValidateStringPattern(value, field, pattern, patternName string) {
	if value == "" {
		return // Skip pattern validation for empty strings
	}

	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		v.AddError(field, fmt.Sprintf("%s pattern validation failed: %v", field, err), "pattern_error", value)
		return
	}

	if !matched {
		v.AddError(field, fmt.Sprintf("%s must match %s format", field, patternName), "pattern_mismatch", value)
	}
}

// ValidateProjectName validates project name with comprehensive rules
func (v *Validator) ValidateProjectName(name string) {
	field := "project_name"

	// Check if empty
	if strings.TrimSpace(name) == "" {
		v.AddError(field, "Project name is required", "required", name)
		return
	}

	// Length validation
	v.ValidateStringLength(name, field, 1, 100)

	// Pattern validation - alphanumeric, hyphens, underscores
	v.ValidateStringPattern(name, field, `^[a-zA-Z0-9_-]+$`, "alphanumeric with hyphens and underscores")

	// Must start with letter or number
	if len(name) > 0 && !unicode.IsLetter(rune(name[0])) && !unicode.IsDigit(rune(name[0])) {
		v.AddError(field, "Project name must start with a letter or number", "invalid_start", name)
	}

	// Reserved names check
	reservedNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	lowerName := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			v.AddError(field, fmt.Sprintf("Project name '%s' is reserved and cannot be used", name), "reserved_name", name)
			break
		}
	}
}

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email, field string) {
	if email == "" {
		return // Skip validation for empty emails if not required
	}

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	v.ValidateStringPattern(email, field, emailPattern, "valid email")
}

// ValidateURL validates URL format
func (v *Validator) ValidateURL(urlStr, field string) {
	if urlStr == "" {
		return // Skip validation for empty URLs if not required
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		v.AddError(field, fmt.Sprintf("%s must be a valid URL", field), "invalid_url", urlStr)
		return
	}

	if parsedURL.Scheme == "" {
		v.AddError(field, fmt.Sprintf("%s must include a scheme (http:// or https://)", field), "missing_scheme", urlStr)
	}

	if parsedURL.Host == "" {
		v.AddError(field, fmt.Sprintf("%s must include a host", field), "missing_host", urlStr)
	}
}

// ValidateLicenseType validates license type against supported licenses
func (v *Validator) ValidateLicenseType(license, field string) {
	if license == "" {
		// Use default license if not specified
		return
	}

	supportedLicenses := []string{
		"MIT",
		"Apache-2.0",
		"GPL-3.0",
		"BSD-3-Clause",
	}

	for _, supported := range supportedLicenses {
		if license == supported {
			return // Valid license found
		}
	}

	// License not supported - add error with suggestion
	v.AddError(field, fmt.Sprintf("License '%s' is not supported. Supported licenses are: %s. Will default to Apache-2.0", license, strings.Join(supportedLicenses, ", ")), "unsupported_license", license)
}

// Path validation methods

// ValidateDirectoryPath validates directory path format and security
func (v *Validator) ValidateDirectoryPath(path, field string) {
	if path == "" {
		v.AddError(field, fmt.Sprintf("%s is required", field), "required", path)
		return
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		v.AddError(field, fmt.Sprintf("%s contains invalid path traversal", field), "path_traversal", path)
	}

	// Only check for absolute paths if the field explicitly requires relative paths
	if strings.Contains(field, "relative_") && filepath.IsAbs(cleanPath) {
		v.AddError(field, fmt.Sprintf("%s should be a relative path", field), "absolute_path", path)
	}

	// Validate path length
	v.ValidateStringLength(path, field, 1, 260) // Windows MAX_PATH limit
}

// ValidateFilePath validates file path format and security
func (v *Validator) ValidateFilePath(path, field string) {
	v.ValidateDirectoryPath(path, field) // Use same validation as directory

	// Additional file-specific validation
	if path != "" {
		ext := filepath.Ext(path)
		if ext == "" {
			v.AddError(field, fmt.Sprintf("%s should have a file extension", field), "missing_extension", path)
		}
	}
}

// Numeric validation methods

// ValidateIntRange validates integer within a range
func (v *Validator) ValidateIntRange(value int, field string, min, max int) {
	if value < min {
		v.AddError(field, fmt.Sprintf("%s must be at least %d", field, min), "min_value", value)
	}

	if value > max {
		v.AddError(field, fmt.Sprintf("%s must be no more than %d", field, max), "max_value", value)
	}
}

// ValidatePositiveInt validates that an integer is positive
func (v *Validator) ValidatePositiveInt(value int, field string) {
	if value <= 0 {
		v.AddError(field, fmt.Sprintf("%s must be a positive number", field), "positive_required", value)
	}
}

// Collection validation methods

// ValidateNonEmptySlice validates that a slice is not empty
func ValidateNonEmptySlice(slice interface{}, fieldName string) error {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("%s must be a slice", fieldName)
	}
	if v.Len() == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateSliceLength validates slice length constraints
func (v *Validator) ValidateSliceLength(slice interface{}, field string, minLength, maxLength int) {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		v.AddError(field, fmt.Sprintf("%s must be a slice", field), "invalid_type", slice)
		return
	}

	length := val.Len()

	if minLength > 0 && length < minLength {
		v.AddError(field, fmt.Sprintf("%s must contain at least %d items", field, minLength), "min_length", slice)
	}

	if maxLength > 0 && length > maxLength {
		v.AddError(field, fmt.Sprintf("%s must contain no more than %d items", field, maxLength), "max_length", slice)
	}
}

// ValidateUniqueSlice validates that all slice elements are unique
func (v *Validator) ValidateUniqueSlice(slice interface{}, field string) {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		v.AddError(field, fmt.Sprintf("%s must be a slice", field), "invalid_type", slice)
		return
	}

	seen := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		if seen[item] {
			v.AddError(field, fmt.Sprintf("%s contains duplicate values", field), "duplicate_values", slice)
			return
		}
		seen[item] = true
	}
}

// Nil validation methods

// ValidateNotNil validates that a value is not nil
func ValidateNotNil(value interface{}, fieldName string) error {
	if value == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}

	// Check for nil pointers, slices, maps, etc.
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if v.IsNil() {
			return fmt.Errorf("%s cannot be nil", fieldName)
		}
	}

	return nil
}

// ValidateRequired validates that a value is not nil or empty
func (v *Validator) ValidateRequired(value interface{}, field string) {
	if value == nil {
		v.AddError(field, fmt.Sprintf("%s is required", field), "required", value)
		return
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.String:
		if strings.TrimSpace(val.String()) == "" {
			v.AddError(field, fmt.Sprintf("%s is required", field), "required", value)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			v.AddError(field, fmt.Sprintf("%s is required", field), "required", value)
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			v.AddError(field, fmt.Sprintf("%s is required", field), "required", value)
		}
	}
}

// Project-specific validation methods

// ValidateProjectConfig validates a complete project configuration
func (v *Validator) ValidateProjectConfig(config *models.ProjectConfig) {
	if config == nil {
		v.AddError("config", "Project configuration is required", "required", nil)
		return
	}

	// Validate project name
	v.ValidateProjectName(config.Name)

	// Validate organization
	if config.Organization != "" {
		v.ValidateStringLength(config.Organization, "organization", 1, 100)
		v.ValidateStringPattern(config.Organization, "organization", `^[a-zA-Z0-9_.-]+$`, "valid organization name")
	}

	// Validate description
	if config.Description != "" {
		v.ValidateStringLength(config.Description, "description", 0, 500)
	}

	// Validate author email if provided
	if config.Email != "" {
		v.ValidateEmail(config.Email, "email")
	}

	// Validate license type
	v.ValidateLicenseType(config.License, "license")

	// Repository field removed

	// Validate output path
	if config.OutputPath != "" {
		v.ValidateDirectoryPath(config.OutputPath, "output_path")
	}

	// Validate components - check if at least one component is selected
	hasAnyComponent := config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin ||
		config.Components.Backend.GoGin ||
		config.Components.Mobile.Android || config.Components.Mobile.IOS ||
		config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker

	if !hasAnyComponent {
		v.AddError("components", "At least one component must be selected", "required", config.Components)
	}
}

// Security validation methods

// ValidateSecureString validates string for security concerns
func (v *Validator) ValidateSecureString(value, field string) {
	if value == "" {
		return
	}

	// Check for potential injection patterns
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`data:`,
		`vbscript:`,
		`onload=`,
		`onerror=`,
		`../`,
		`..\\`,
	}

	lowerValue := strings.ToLower(value)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerValue, pattern) {
			v.AddError(field, fmt.Sprintf("%s contains potentially dangerous content", field), "security_risk", value)
			return
		}
	}

	// Check for excessive length (potential DoS)
	if len(value) > 10000 {
		v.AddError(field, fmt.Sprintf("%s is too long (potential security risk)", field), "excessive_length", value)
	}
}

// ValidateNoSQLInjection validates string for SQL injection patterns
func (v *Validator) ValidateNoSQLInjection(value, field string) {
	if value == "" {
		return
	}

	sqlPatterns := []string{
		`'`,
		`"`,
		`;`,
		`--`,
		`/*`,
		`*/`,
		`union`,
		`select`,
		`insert`,
		`update`,
		`delete`,
		`drop`,
		`create`,
		`alter`,
	}

	lowerValue := strings.ToLower(value)
	for _, pattern := range sqlPatterns {
		if strings.Contains(lowerValue, pattern) {
			v.AddError(field, fmt.Sprintf("%s contains potentially dangerous SQL patterns", field), "sql_injection_risk", value)
			return
		}
	}
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
