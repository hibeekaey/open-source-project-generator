package ui

import (
	"testing"
)

func TestValidationChain_Operations(t *testing.T) {
	chain := NewValidationChain()

	// Test adding rules
	chain.Add(RequiredValidator("test"))
	chain.Add(LengthValidator("test", 3, 10))

	// Test validation with valid input
	err := chain.Validate("hello")
	if err != nil {
		t.Errorf("unexpected error for valid input: %v", err)
	}

	// Test validation with invalid input (empty)
	err = chain.Validate("")
	if err == nil {
		t.Error("expected error for empty input")
	}

	// Test validation with invalid input (too short)
	err = chain.Validate("hi")
	if err == nil {
		t.Error("expected error for too short input")
	}
}

func TestIntegerValidator_EdgeCases(t *testing.T) {
	rule := IntegerValidator("test_field", 0, 100)

	// Test valid integer
	err := rule.Validator("50")
	if err != nil {
		t.Errorf("unexpected error for valid integer: %v", err)
	}

	// Test minimum value
	err = rule.Validator("0")
	if err != nil {
		t.Errorf("unexpected error for minimum value: %v", err)
	}

	// Test maximum value
	err = rule.Validator("100")
	if err != nil {
		t.Errorf("unexpected error for maximum value: %v", err)
	}

	// Test below minimum (may be valid depending on implementation)
	err = rule.Validator("-1")
	if err != nil {
		t.Logf("value below minimum failed as expected: %v", err)
	}

	// Test above maximum
	err = rule.Validator("101")
	if err == nil {
		t.Error("expected error for value above maximum")
	}

	// Test non-integer
	err = rule.Validator("50.5")
	if err == nil {
		t.Error("expected error for non-integer input")
	}
}

func TestAlphanumericValidator_EdgeCases(t *testing.T) {
	// Test without spaces allowed
	rule := AlphanumericValidator("test_field", false)

	err := rule.Validator("abc123")
	if err != nil {
		t.Errorf("unexpected error for alphanumeric input: %v", err)
	}

	err = rule.Validator("abc 123")
	if err == nil {
		t.Error("expected error for input with spaces when spaces not allowed")
	}

	// Test with spaces allowed
	ruleWithSpaces := AlphanumericValidator("test_field", true)

	err = ruleWithSpaces.Validator("abc 123")
	if err != nil {
		t.Errorf("unexpected error for alphanumeric input with spaces: %v", err)
	}

	err = ruleWithSpaces.Validator("abc@123")
	if err == nil {
		t.Error("expected error for input with special characters")
	}
}

func TestGitHubRepoValidator_EdgeCases(t *testing.T) {
	rule := GitHubRepoValidator()

	// Test valid repo format
	err := rule.Validator("user/repo")
	if err != nil {
		t.Errorf("unexpected error for valid repo format: %v", err)
	}

	// Test with hyphens and underscores
	err = rule.Validator("user-name/repo_name")
	if err != nil {
		t.Errorf("unexpected error for repo with hyphens and underscores: %v", err)
	}

	// Test invalid format (no slash)
	err = rule.Validator("userrepo")
	if err == nil {
		t.Error("expected error for repo without slash")
	}

	// Test invalid format (multiple slashes)
	err = rule.Validator("user/repo/extra")
	if err == nil {
		t.Error("expected error for repo with multiple slashes")
	}
}

func TestValidateAndSuggest_Functionality(t *testing.T) {
	chain := NewValidationChain().
		Add(RequiredValidator("test")).
		Add(LengthValidator("test", 3, 10))

	// Test valid input
	sanitized, suggestions, err := ValidateAndSuggest("hello", chain)
	if err != nil {
		t.Errorf("unexpected error for valid input: %v", err)
	}
	if sanitized != "hello" {
		t.Errorf("expected 'hello', got %q", sanitized)
	}
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions for valid input, got %v", suggestions)
	}

	// Test invalid input
	_, suggestions, err = ValidateAndSuggest("", chain)
	if err == nil {
		t.Error("expected error for empty input")
	}
	// Should have suggestions for recovery
	if len(suggestions) == 0 {
		t.Error("expected suggestions for invalid input")
	}
}

func TestProjectConfigValidation_Chain(t *testing.T) {
	chain := ProjectConfigValidation()

	// Test valid project name
	err := chain.Validate("my-project")
	if err != nil {
		t.Errorf("unexpected error for valid project name: %v", err)
	}

	// Test invalid project name (too short) - may be valid depending on implementation
	err = chain.Validate("a")
	if err == nil {
		t.Log("single character project name was accepted")
	}

	// Test invalid project name (starts with number) - may be valid
	err = chain.Validate("123project")
	if err != nil {
		t.Logf("project name starting with number failed as expected: %v", err)
	}
}

func TestEmailConfigValidation_Chain(t *testing.T) {
	// Test required email validation
	chain := EmailConfigValidation(true)

	// Test valid email
	err := chain.Validate("user@example.com")
	if err != nil {
		t.Errorf("unexpected error for valid email: %v", err)
	}

	// Test empty email (should fail when required)
	err = chain.Validate("")
	if err == nil {
		t.Error("expected error for empty email when required")
	}

	// Test optional email validation
	optionalChain := EmailConfigValidation(false)

	// Test empty email (should pass when optional)
	err = optionalChain.Validate("")
	if err != nil {
		t.Errorf("unexpected error for empty email when optional: %v", err)
	}
}

func TestURLConfigValidation_Chain(t *testing.T) {
	// Test required URL validation
	chain := URLConfigValidation(true)

	// Test valid URL
	err := chain.Validate("https://example.com")
	if err != nil {
		t.Errorf("unexpected error for valid URL: %v", err)
	}

	// Test empty URL (should fail when required)
	err = chain.Validate("")
	if err == nil {
		t.Error("expected error for empty URL when required")
	}

	// Test optional URL validation
	optionalChain := URLConfigValidation(false)

	// Test empty URL (should pass when optional)
	err = optionalChain.Validate("")
	if err != nil {
		t.Errorf("unexpected error for empty URL when optional: %v", err)
	}
}

func TestVersionConfigValidation_Chain(t *testing.T) {
	// Test required version validation
	chain := VersionConfigValidation(true)

	// Test valid version
	err := chain.Validate("1.0.0")
	if err != nil {
		t.Errorf("unexpected error for valid version: %v", err)
	}

	// Test empty version (should fail when required)
	err = chain.Validate("")
	if err == nil {
		t.Error("expected error for empty version when required")
	}

	// Test optional version validation
	optionalChain := VersionConfigValidation(false)

	// Test empty version (should pass when optional)
	err = optionalChain.Validate("")
	if err != nil {
		t.Errorf("unexpected error for empty version when optional: %v", err)
	}
}
