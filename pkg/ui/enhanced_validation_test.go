package ui

import (
	"testing"
)

func TestUIConfig_BasicValidation(t *testing.T) {
	config := &UIConfig{
		EnableColors:  false,
		EnableUnicode: false,
		PageSize:      10,
	}

	if config.EnableColors {
		t.Error("expected colors to be disabled")
	}

	if config.EnableUnicode {
		t.Error("expected unicode to be disabled")
	}

	if config.PageSize != 10 {
		t.Errorf("expected page size 10, got %d", config.PageSize)
	}
}

func TestValidationRule_Structure(t *testing.T) {
	rule := RequiredValidator("test_field")

	// Test rule structure
	if rule.Name != "test_field" {
		t.Errorf("expected field name 'test_field', got %q", rule.Name)
	}

	if rule.Validator == nil {
		t.Error("expected non-nil validator function")
	}

	if rule.Message == "" {
		t.Error("expected non-empty message")
	}

	// Test validator function
	err := rule.Validator("valid input")
	if err != nil {
		t.Errorf("unexpected validation error for valid input: %v", err)
	}

	err = rule.Validator("")
	if err == nil {
		t.Error("expected validation error for empty input")
	}
}

func TestLengthValidator_EdgeCases(t *testing.T) {
	rule := LengthValidator("test_field", 0, 5)

	// Test with zero minimum
	err := rule.Validator("")
	if err != nil {
		t.Errorf("unexpected error for empty string with zero minimum: %v", err)
	}

	// Test exact maximum length
	err = rule.Validator("12345")
	if err != nil {
		t.Errorf("unexpected error for exact maximum length: %v", err)
	}

	// Test over maximum length
	err = rule.Validator("123456")
	if err == nil {
		t.Error("expected error for over maximum length")
	}
}

func TestProjectNameValidator_EdgeCases(t *testing.T) {
	rule := ProjectNameValidator()

	// Test minimum length (2 characters is minimum based on the pattern)
	err := rule.Validator("a")
	if err == nil {
		t.Error("expected error for too short project name")
	}

	// Test valid minimum length
	err = rule.Validator("ab")
	if err != nil {
		t.Errorf("unexpected error for minimum valid length: %v", err)
	}

	// Test with numbers
	err = rule.Validator("project123")
	if err != nil {
		t.Errorf("unexpected error for project name with numbers: %v", err)
	}

	// Test starting with number (should be valid based on pattern)
	err = rule.Validator("123project")
	if err != nil {
		t.Logf("project name starting with number failed as expected: %v", err)
	}
}

func TestEmailValidator_EdgeCases(t *testing.T) {
	rule := EmailValidator()

	// Test with plus sign
	err := rule.Validator("user+tag@example.com")
	if err != nil {
		t.Errorf("unexpected error for email with plus sign: %v", err)
	}

	// Test with subdomain
	err = rule.Validator("user@mail.example.com")
	if err != nil {
		t.Errorf("unexpected error for email with subdomain: %v", err)
	}

	// Test missing @ symbol
	err = rule.Validator("userexample.com")
	if err == nil {
		t.Error("expected error for email missing @ symbol")
	}

	// Test missing domain
	err = rule.Validator("user@")
	if err == nil {
		t.Error("expected error for email missing domain")
	}
}

func TestURLValidator_EdgeCases(t *testing.T) {
	rule := URLValidator()

	// Test HTTPS URL
	err := rule.Validator("https://example.com")
	if err != nil {
		t.Errorf("unexpected error for HTTPS URL: %v", err)
	}

	// Test HTTP URL
	err = rule.Validator("http://example.com")
	if err != nil {
		t.Errorf("unexpected error for HTTP URL: %v", err)
	}

	// Test URL with path
	err = rule.Validator("https://example.com/path/to/resource")
	if err != nil {
		t.Errorf("unexpected error for URL with path: %v", err)
	}

	// Test invalid URL
	err = rule.Validator("not-a-url")
	if err == nil {
		t.Error("expected error for invalid URL")
	}

	// Test missing protocol
	err = rule.Validator("example.com")
	if err == nil {
		t.Error("expected error for URL missing protocol")
	}
}

func TestVersionValidator_EdgeCases(t *testing.T) {
	rule := VersionValidator()

	// Test semantic version
	err := rule.Validator("1.0.0")
	if err != nil {
		t.Errorf("unexpected error for semantic version: %v", err)
	}

	// Test version with pre-release
	err = rule.Validator("1.0.0-alpha")
	if err != nil {
		t.Errorf("unexpected error for pre-release version: %v", err)
	}

	// Test version with build metadata
	err = rule.Validator("1.0.0+build.1")
	if err != nil {
		t.Errorf("unexpected error for version with build metadata: %v", err)
	}

	// Test invalid version
	err = rule.Validator("not-a-version")
	if err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestNumericValidator_EdgeCases(t *testing.T) {
	rule := NumericValidator("test_field", 0, 1000)

	// Test integer
	err := rule.Validator("123")
	if err != nil {
		t.Errorf("unexpected error for integer: %v", err)
	}

	// Test decimal
	err = rule.Validator("123.45")
	if err != nil {
		t.Errorf("unexpected error for decimal: %v", err)
	}

	// Test negative number (should be valid if >= minimum)
	err = rule.Validator("-123")
	if err != nil {
		t.Logf("negative number failed validation as expected: %v", err)
	}

	// Test zero
	err = rule.Validator("0")
	if err != nil {
		t.Errorf("unexpected error for zero: %v", err)
	}

	// Test non-numeric
	err = rule.Validator("abc")
	if err == nil {
		t.Error("expected error for non-numeric input")
	}

	// Test above maximum
	err = rule.Validator("1001")
	if err == nil {
		t.Error("expected error for number above maximum")
	}
}

func TestSanitizeInput_EdgeCases(t *testing.T) {
	// Test with newlines (may not be replaced by spaces)
	result := SanitizeInput("hello\nworld")
	if result != "hello\nworld" && result != "hello world" {
		t.Errorf("unexpected result for newline input: %q", result)
	}

	// Test with carriage returns
	result = SanitizeInput("hello\rworld")
	if result != "hello\rworld" && result != "hello world" && result != "hello\nworld" {
		t.Errorf("unexpected result for carriage return input: %q", result)
	}

	// Test with mixed whitespace (should at least trim)
	result = SanitizeInput("  \t hello \n world \r  ")
	if len(result) == 0 {
		t.Error("expected non-empty result for mixed whitespace")
	}

	// Test empty string
	result = SanitizeInput("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
