package utils

import (
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestValidateNonEmptyString(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantError bool
	}{
		{"valid string", "hello", "field", false},
		{"empty string", "", "field", true},
		{"whitespace only", "   ", "field", true},
		{"string with content", "  hello  ", "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptyString(tt.value, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNonEmptyString() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && !strings.Contains(err.Error(), tt.fieldName) {
				t.Errorf("Error should contain field name %s, got %s", tt.fieldName, err.Error())
			}
		})
	}
}

func TestValidateNonEmptySlice(t *testing.T) {
	tests := []struct {
		name      string
		slice     interface{}
		fieldName string
		wantError bool
	}{
		{"valid slice", []string{"a", "b"}, "items", false},
		{"empty slice", []string{}, "items", true},
		{"nil slice", []string(nil), "items", true},
		{"not a slice", "not a slice", "items", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptySlice(tt.slice, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNonEmptySlice() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateNotNil(t *testing.T) {
	var nilPtr *string
	validPtr := new(string)

	tests := []struct {
		name      string
		value     interface{}
		fieldName string
		wantError bool
	}{
		{"valid pointer", validPtr, "ptr", false},
		{"nil pointer", nilPtr, "ptr", true},
		{"nil interface", nil, "interface", true},
		{"valid string", "hello", "str", false},
		{"valid int", 42, "num", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotNil(tt.value, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNotNil() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNewValidationErrorStruct(t *testing.T) {
	field := "test_field"
	message := "test message"
	code := "test_code"
	value := "test_value"

	err := NewValidationErrorStruct(field, message, code, value)

	if err.Field != field {
		t.Errorf("Field = %v, want %v", err.Field, field)
	}
	if err.Message != message {
		t.Errorf("Message = %v, want %v", err.Message, message)
	}
	if err.Code != code {
		t.Errorf("Code = %v, want %v", err.Code, code)
	}
	if err.Value != value {
		t.Errorf("Value = %v, want %v", err.Value, value)
	}
}

func TestValidator(t *testing.T) {
	v := NewValidator()

	// Test initial state
	if v.HasErrors() {
		t.Error("New validator should not have errors")
	}
	if len(v.GetErrors()) != 0 {
		t.Error("New validator should have empty error list")
	}

	// Test adding errors
	v.AddError("field1", "message1", "code1", "value1")
	if !v.HasErrors() {
		t.Error("Validator should have errors after adding one")
	}
	if len(v.GetErrors()) != 1 {
		t.Errorf("Error count = %v, want 1", len(v.GetErrors()))
	}

	// Test adding multiple errors
	v.AddError("field2", "message2", "code2", "value2")
	if len(v.GetErrors()) != 2 {
		t.Errorf("Error count = %v, want 2", len(v.GetErrors()))
	}

	// Test get result
	result := v.GetResult()
	if result.IsValid {
		t.Error("Result should not be valid when there are errors")
	}
	if len(result.Errors) != 2 {
		t.Errorf("Result error count = %v, want 2", len(result.Errors))
	}

	// Test clear
	v.Clear()
	if v.HasErrors() {
		t.Error("Validator should not have errors after clear")
	}
	if len(v.GetErrors()) != 0 {
		t.Error("Validator should have empty error list after clear")
	}

	// Test result after clear
	result = v.GetResult()
	if !result.IsValid {
		t.Error("Result should be valid when there are no errors")
	}
}

func TestValidateStringLength(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     string
		field     string
		minLength int
		maxLength int
		wantError bool
	}{
		{"valid length", "hello", "field", 3, 10, false},
		{"too short", "hi", "field", 3, 10, true},
		{"too long", "this is a very long string", "field", 3, 10, true},
		{"exact min", "abc", "field", 3, 10, false},
		{"exact max", "1234567890", "field", 3, 10, false},
		{"no min limit", "hi", "field", 0, 10, false},
		{"no max limit", "this is a very long string", "field", 3, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateStringLength(tt.value, tt.field, tt.minLength, tt.maxLength)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateStringLength() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateStringPattern(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name        string
		value       string
		field       string
		pattern     string
		patternName string
		wantError   bool
	}{
		{"valid pattern", "hello123", "field", `^[a-z0-9]+$`, "alphanumeric", false},
		{"invalid pattern", "Hello123!", "field", `^[a-z0-9]+$`, "alphanumeric", true},
		{"empty value", "", "field", `^[a-z0-9]+$`, "alphanumeric", false}, // Skip validation for empty
		{"invalid regex", "hello", "field", `[`, "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateStringPattern(tt.value, tt.field, tt.pattern, tt.patternName)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateStringPattern() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"valid name", "my-project_1", false},
		{"empty name", "", true},
		{"whitespace only", "   ", true},
		{"too long", strings.Repeat("a", 101), true},
		{"invalid characters", "my-project!", true},
		{"starts with special char", "-myproject", true},
		{"reserved name", "con", true},
		{"reserved name uppercase", "CON", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateProjectName(tt.value)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateProjectName() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		email     string
		field     string
		wantError bool
	}{
		{"valid email", "test@example.com", "email", false},
		{"empty email", "", "email", false}, // Skip validation for empty
		{"invalid email", "invalid-email", "email", true},
		{"missing domain", "test@", "email", true},
		{"missing @", "testexample.com", "email", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateEmail(tt.email, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateEmail() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		url       string
		field     string
		wantError bool
	}{
		{"valid URL", "https://example.com", "url", false},
		{"empty URL", "", "url", false}, // Skip validation for empty
		{"invalid URL", "not-a-url", "url", true},
		{"missing scheme", "example.com", "url", true},
		{"missing host", "https://", "url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateURL(tt.url, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateURL() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateLicenseType(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		license   string
		field     string
		wantError bool
	}{
		{"valid license", "MIT", "license", false},
		{"empty license", "", "license", false}, // Use default
		{"unsupported license", "CUSTOM", "license", true},
		{"Apache license", "Apache-2.0", "license", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateLicenseType(tt.license, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateLicenseType() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateDirectoryPath(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		path      string
		field     string
		wantError bool
	}{
		{"valid path", "src/main", "path", false},
		{"empty path", "", "path", true},
		{"path traversal", "../../../etc", "path", true},
		{"absolute path with relative field", "/absolute/path", "relative_path", true},
		{"long path", strings.Repeat("a", 261), "path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateDirectoryPath(tt.path, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateDirectoryPath() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		path      string
		field     string
		wantError bool
	}{
		{"valid file path", "src/main.go", "path", false},
		{"empty path", "", "path", true},
		{"no extension", "src/main", "path", true},
		{"path traversal", "../../../etc/passwd", "path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateFilePath(tt.path, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateFilePath() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateIntRange(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     int
		field     string
		min       int
		max       int
		wantError bool
	}{
		{"valid range", 5, "field", 1, 10, false},
		{"below min", 0, "field", 1, 10, true},
		{"above max", 11, "field", 1, 10, true},
		{"exact min", 1, "field", 1, 10, false},
		{"exact max", 10, "field", 1, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateIntRange(tt.value, tt.field, tt.min, tt.max)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateIntRange() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidatePositiveInt(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     int
		field     string
		wantError bool
	}{
		{"positive number", 5, "field", false},
		{"zero", 0, "field", true},
		{"negative number", -5, "field", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidatePositiveInt(tt.value, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidatePositiveInt() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateSliceLength(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		slice     interface{}
		field     string
		minLength int
		maxLength int
		wantError bool
	}{
		{"valid length", []string{"a", "b", "c"}, "field", 2, 5, false},
		{"too short", []string{"a"}, "field", 2, 5, true},
		{"too long", []string{"a", "b", "c", "d", "e", "f"}, "field", 2, 5, true},
		{"not a slice", "not a slice", "field", 2, 5, true},
		{"no min limit", []string{"a"}, "field", 0, 5, false},
		{"no max limit", []string{"a", "b", "c", "d", "e", "f"}, "field", 2, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateSliceLength(tt.slice, tt.field, tt.minLength, tt.maxLength)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateSliceLength() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateUniqueSlice(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		slice     interface{}
		field     string
		wantError bool
	}{
		{"unique elements", []string{"a", "b", "c"}, "field", false},
		{"duplicate elements", []string{"a", "b", "a"}, "field", true},
		{"not a slice", "not a slice", "field", true},
		{"empty slice", []string{}, "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateUniqueSlice(tt.slice, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateUniqueSlice() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     interface{}
		field     string
		wantError bool
	}{
		{"valid string", "hello", "field", false},
		{"empty string", "", "field", true},
		{"whitespace string", "   ", "field", true},
		{"nil value", nil, "field", true},
		{"valid slice", []string{"a"}, "field", false},
		{"empty slice", []string{}, "field", true},
		{"nil pointer", (*string)(nil), "field", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateRequired(tt.value, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateRequired() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateProjectConfig(t *testing.T) {
	v := NewValidator()

	validConfig := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test description",
		Email:        "test@example.com",
		License:      "MIT",
		OutputPath:   "output",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
	}

	invalidConfig := &models.ProjectConfig{
		Name:         "",                  // Invalid: empty name
		Organization: "test@org!",         // Invalid: special characters
		Email:        "invalid-email",     // Invalid: bad email format
		License:      "CUSTOM",            // Invalid: unsupported license
		OutputPath:   "../../../etc",      // Invalid: path traversal
		Components:   models.Components{}, // Invalid: no components selected
	}

	tests := []struct {
		name      string
		config    *models.ProjectConfig
		wantError bool
	}{
		{"nil config", nil, true},
		{"valid config", validConfig, false},
		{"invalid config", invalidConfig, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateProjectConfig(tt.config)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateProjectConfig() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateSecureString(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     string
		field     string
		wantError bool
	}{
		{"safe string", "hello world", "field", false},
		{"empty string", "", "field", false}, // Skip validation for empty
		{"script tag", "<script>alert('xss')</script>", "field", true},
		{"javascript protocol", "javascript:alert('xss')", "field", true},
		{"path traversal", "../../../etc/passwd", "field", true},
		{"excessive length", strings.Repeat("a", 10001), "field", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateSecureString(tt.value, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateSecureString() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateNoSQLInjection(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name      string
		value     string
		field     string
		wantError bool
	}{
		{"safe string", "hello world", "field", false},
		{"empty string", "", "field", false}, // Skip validation for empty
		{"single quote", "user's name", "field", true},
		{"sql comment", "user -- comment", "field", true},
		{"union attack", "1 UNION SELECT * FROM users", "field", true},
		{"drop table", "'; DROP TABLE users; --", "field", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v.Clear()
			v.ValidateNoSQLInjection(tt.value, tt.field)
			if v.HasErrors() != tt.wantError {
				t.Errorf("ValidateNoSQLInjection() hasErrors = %v, wantError %v", v.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestFormatValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		errors []ValidationError
		want   string
	}{
		{"no errors", []ValidationError{}, ""},
		{"single error", []ValidationError{
			{Field: "field1", Message: "Error 1", Code: "code1", Value: "value1"},
		}, "Validation failed:\n• Error 1"},
		{"multiple errors", []ValidationError{
			{Field: "field1", Message: "Error 1", Code: "code1", Value: "value1"},
			{Field: "field2", Message: "Error 2", Code: "code2", Value: "value2"},
		}, "Validation failed:\n• Error 1\n• Error 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatValidationErrors(tt.errors)
			if got != tt.want {
				t.Errorf("FormatValidationErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValidationErrorsByField(t *testing.T) {
	errors := []ValidationError{
		{Field: "field1", Message: "Error 1", Code: "code1", Value: "value1"},
		{Field: "field2", Message: "Error 2", Code: "code2", Value: "value2"},
		{Field: "field1", Message: "Error 3", Code: "code3", Value: "value3"},
	}

	result := GetValidationErrorsByField(errors)

	if len(result) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(result))
	}

	if len(result["field1"]) != 2 {
		t.Errorf("Expected 2 errors for field1, got %d", len(result["field1"]))
	}

	if len(result["field2"]) != 1 {
		t.Errorf("Expected 1 error for field2, got %d", len(result["field2"]))
	}
}
