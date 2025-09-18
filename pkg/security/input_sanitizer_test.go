package security

import (
	"strings"
	"testing"
)

func TestInputSanitizer_SanitizeString(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name      string
		input     string
		fieldName string
		wantValid bool
		wantError bool
	}{
		{
			name:      "valid string",
			input:     "hello world",
			fieldName: "test",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			fieldName: "test",
			wantValid: true,
			wantError: false,
		},
		{
			name:      "string with script tag",
			input:     "<script>alert('xss')</script>",
			fieldName: "test",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "string with javascript protocol",
			input:     "javascript:alert('xss')",
			fieldName: "test",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "string with path traversal",
			input:     "../../../etc/passwd",
			fieldName: "test",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "string with SQL injection",
			input:     "'; DROP TABLE users; --",
			fieldName: "test",
			wantValid: false,
			wantError: true,
		},
		{
			name:      "very long string",
			input:     strings.Repeat("a", 20000),
			fieldName: "test",
			wantValid: false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeString(tt.input, tt.fieldName)

			if result.IsValid != tt.wantValid {
				t.Errorf("SanitizeString() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			hasErrors := len(result.Errors) > 0
			if hasErrors != tt.wantError {
				t.Errorf("SanitizeString() has errors = %v, want %v", hasErrors, tt.wantError)
			}
		})
	}
}

func TestInputSanitizer_SanitizeProjectName(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name         string
		input        string
		wantValid    bool
		wantModified bool
		expectedName string
	}{
		{
			name:         "valid project name",
			input:        "my-project",
			wantValid:    true,
			wantModified: false,
			expectedName: "my-project",
		},
		{
			name:         "project name with spaces",
			input:        "My Project Name",
			wantValid:    true,
			wantModified: true,
			expectedName: "my-project-name",
		},
		{
			name:         "project name with invalid characters",
			input:        "my@project#name!",
			wantValid:    true,
			wantModified: true,
			expectedName: "myprojectname",
		},
		{
			name:         "reserved name",
			input:        "con",
			wantValid:    false,
			wantModified: false,
			expectedName: "con",
		},
		{
			name:         "empty after sanitization",
			input:        "!@#$%",
			wantValid:    false,
			wantModified: true,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeProjectName(tt.input)

			if result.IsValid != tt.wantValid {
				t.Errorf("SanitizeProjectName() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}

			if result.WasModified != tt.wantModified {
				t.Errorf("SanitizeProjectName() WasModified = %v, want %v", result.WasModified, tt.wantModified)
			}

			if result.Sanitized != tt.expectedName {
				t.Errorf("SanitizeProjectName() Sanitized = %v, want %v", result.Sanitized, tt.expectedName)
			}
		})
	}
}

func TestInputSanitizer_SanitizeFilePath(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name      string
		input     string
		fieldName string
		wantValid bool
	}{
		{
			name:      "valid relative path",
			input:     "src/main.go",
			fieldName: "file_path",
			wantValid: true,
		},
		{
			name:      "path with traversal",
			input:     "../../../etc/passwd",
			fieldName: "file_path",
			wantValid: false,
		},
		{
			name:      "absolute path for relative field",
			input:     "/usr/bin/bash",
			fieldName: "relative_file_path",
			wantValid: false,
		},
		{
			name:      "path with invalid characters",
			input:     "file<name>.txt",
			fieldName: "file_path",
			wantValid: false,
		},
		{
			name:      "very long path",
			input:     strings.Repeat("a/", 200),
			fieldName: "file_path",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeFilePath(tt.input, tt.fieldName)

			if result.IsValid != tt.wantValid {
				t.Errorf("SanitizeFilePath() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
		})
	}
}

func TestInputSanitizer_SanitizeEmail(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "valid email",
			input:     "user@example.com",
			wantValid: true,
		},
		{
			name:      "email with uppercase",
			input:     "User@Example.COM",
			wantValid: true,
		},
		{
			name:      "invalid email format",
			input:     "not-an-email",
			wantValid: false,
		},
		{
			name:      "email with dangerous content",
			input:     "user+<script>@example.com",
			wantValid: false,
		},
		{
			name:      "empty email",
			input:     "",
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeEmail(tt.input, "email")

			if result.IsValid != tt.wantValid {
				t.Errorf("SanitizeEmail() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
		})
	}
}

func TestInputSanitizer_SanitizeURL(t *testing.T) {
	sanitizer := NewInputSanitizer()

	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "valid HTTP URL",
			input:     "http://example.com",
			wantValid: true,
		},
		{
			name:      "valid HTTPS URL",
			input:     "https://example.com/path",
			wantValid: true,
		},
		{
			name:      "URL without scheme",
			input:     "example.com",
			wantValid: false,
		},
		{
			name:      "javascript protocol",
			input:     "javascript:alert('xss')",
			wantValid: false,
		},
		{
			name:      "data protocol",
			input:     "data:text/html,<script>alert('xss')</script>",
			wantValid: false,
		},
		{
			name:      "URL without host",
			input:     "http://",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeURL(tt.input, "url")

			if result.IsValid != tt.wantValid {
				t.Errorf("SanitizeURL() IsValid = %v, want %v", result.IsValid, tt.wantValid)
			}
		})
	}
}

func TestInputSanitizer_ValidateAndSanitizeMap(t *testing.T) {
	sanitizer := NewInputSanitizer()

	testData := map[string]interface{}{
		"name":        "test-project",
		"description": "A test project",
		"email":       "user@example.com",
		"dangerous":   "<script>alert('xss')</script>",
		"nested": map[string]interface{}{
			"value": "nested-value",
			"bad":   "javascript:alert('xss')",
		},
		"list": []interface{}{
			"item1",
			"<script>bad</script>",
			123, // non-string item
		},
	}

	results := sanitizer.ValidateAndSanitizeMap(testData, "")

	// Check that dangerous content was detected
	if result, exists := results["dangerous"]; exists {
		if result.IsValid {
			t.Error("Expected dangerous content to be invalid")
		}
	}

	// Check nested map handling
	if result, exists := results["nested.bad"]; exists {
		if result.IsValid {
			t.Error("Expected nested dangerous content to be invalid")
		}
	}

	// Check that valid content passed
	if result, exists := results["name"]; exists {
		if !result.IsValid {
			t.Error("Expected valid name to be valid")
		}
	}
}

func TestGetSanitizationSummary(t *testing.T) {
	results := map[string]*SanitizationResult{
		"field1": {
			IsValid:     true,
			WasModified: false,
		},
		"field2": {
			IsValid:     false,
			Errors:      []string{"Invalid content"},
			WasModified: true,
		},
		"field3": {
			IsValid:     true,
			Warnings:    []string{"Minor issue"},
			WasModified: true,
		},
	}

	summary := GetSanitizationSummary(results)

	if summary["total_fields"] != 3 {
		t.Errorf("Expected total_fields = 3, got %v", summary["total_fields"])
	}

	if summary["valid_fields"] != 2 {
		t.Errorf("Expected valid_fields = 2, got %v", summary["valid_fields"])
	}

	if summary["invalid_fields"] != 1 {
		t.Errorf("Expected invalid_fields = 1, got %v", summary["invalid_fields"])
	}

	if summary["modified_fields"] != 2 {
		t.Errorf("Expected modified_fields = 2, got %v", summary["modified_fields"])
	}
}
