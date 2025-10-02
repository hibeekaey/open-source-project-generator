package ui

import (
	"fmt"
	"regexp"
	"testing"
)

func TestCustomValidator_Creation(t *testing.T) {
	validator := func(input string) error {
		if input == "invalid" {
			return fmt.Errorf("input is invalid")
		}
		return nil
	}

	rule := CustomValidator(
		"custom_field",
		"Custom validation failed",
		validator,
		[]string{"try 'valid'"},
		nil,
	)

	// Test rule structure
	if rule.Name != "custom_field" {
		t.Errorf("expected field name 'custom_field', got %q", rule.Name)
	}

	if rule.Message != "Custom validation failed" {
		t.Errorf("expected message 'Custom validation failed', got %q", rule.Message)
	}

	if len(rule.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(rule.Suggestions))
	}

	// Test validator function
	err := rule.Validator("valid")
	if err != nil {
		t.Errorf("unexpected error for valid input: %v", err)
	}

	err = rule.Validator("invalid")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestSuggestCorrection_Functionality(t *testing.T) {
	// Test basic correction suggestion
	suggestion := SuggestCorrection("test@gmial.com", "email")
	if suggestion == "" {
		t.Error("expected non-empty suggestion")
	}

	// Test with different input types
	suggestion = SuggestCorrection("123project", "project_name")
	if suggestion == "" {
		t.Error("expected non-empty suggestion for project name")
	}

	// Test with valid input (should still provide suggestion)
	suggestion = SuggestCorrection("valid-project", "project_name")
	// Should not be empty as it provides a suggestion mechanism
	if suggestion == "" {
		t.Error("expected suggestion even for valid input")
	}
}

func TestCreateDefaultValueRecovery_Options(t *testing.T) {
	recovery := CreateDefaultValueRecovery("default-value")

	if recovery.Label == "" {
		t.Error("expected non-empty label")
	}

	if recovery.Action == nil {
		t.Error("expected non-nil action function")
	}

	// Test the action function
	err := recovery.Action()
	if err != nil {
		t.Errorf("unexpected error from recovery action: %v", err)
	}
}

func TestCreateSuggestionRecovery_Options(t *testing.T) {
	recovery := CreateSuggestionRecovery("suggested-value")

	if recovery.Label == "" {
		t.Error("expected non-empty label")
	}

	if recovery.Action == nil {
		t.Error("expected non-nil action function")
	}

	// Test the action function
	err := recovery.Action()
	if err != nil {
		t.Errorf("unexpected error from recovery action: %v", err)
	}
}

func TestCreateSkipFieldRecovery_Options(t *testing.T) {
	recovery := CreateSkipFieldRecovery("test_field")

	if recovery.Label == "" {
		t.Error("expected non-empty label")
	}

	if recovery.Action == nil {
		t.Error("expected non-nil action function")
	}

	// Test the action function
	err := recovery.Action()
	if err != nil {
		t.Errorf("unexpected error from recovery action: %v", err)
	}
}

func TestValidationChain_AddMultiple(t *testing.T) {
	chain := NewValidationChain()

	// Add multiple rules
	chain.Add(RequiredValidator("test"))
	chain.Add(LengthValidator("test", 3, 10))
	chain.Add(AlphanumericValidator("test", false))

	// Test with valid input
	err := chain.Validate("hello123")
	if err != nil {
		t.Errorf("unexpected error for valid input: %v", err)
	}

	// Test with invalid input (empty - fails required)
	err = chain.Validate("")
	if err == nil {
		t.Error("expected error for empty input")
	}

	// Test with invalid input (too short - fails length)
	err = chain.Validate("hi")
	if err == nil {
		t.Error("expected error for too short input")
	}

	// Test with invalid input (special chars - fails alphanumeric)
	err = chain.Validate("hello@world")
	if err == nil {
		t.Error("expected error for input with special characters")
	}
}

func TestPatternValidator_CustomPattern(t *testing.T) {
	// Create a custom pattern for testing
	pattern := regexp.MustCompile(`^[A-Z][a-z]+$`)
	rule := PatternValidator("test_field", pattern, "Must start with uppercase letter followed by lowercase letters")

	// Test valid input
	err := rule.Validator("Hello")
	if err != nil {
		t.Errorf("unexpected error for valid input: %v", err)
	}

	// Test invalid input (starts with lowercase)
	err = rule.Validator("hello")
	if err == nil {
		t.Error("expected error for input starting with lowercase")
	}

	// Test invalid input (contains numbers)
	err = rule.Validator("Hello123")
	if err == nil {
		t.Error("expected error for input containing numbers")
	}

	// Test invalid input (all uppercase)
	err = rule.Validator("HELLO")
	if err == nil {
		t.Error("expected error for all uppercase input")
	}
}

func TestValidationRule_MessageAndSuggestions(t *testing.T) {
	rule := ProjectNameValidator()

	// Check that rule has proper message
	if rule.Message == "" {
		t.Error("expected non-empty validation message")
	}

	// Check that rule has suggestions
	if len(rule.Suggestions) == 0 {
		t.Error("expected validation suggestions")
	}

	// Check that rule has recovery options
	if len(rule.Recovery) == 0 {
		t.Error("expected recovery options")
	}
}

func TestValidationChain_EmptyChain(t *testing.T) {
	chain := NewValidationChain()

	// Empty chain should pass any input
	err := chain.Validate("anything")
	if err != nil {
		t.Errorf("unexpected error for empty validation chain: %v", err)
	}

	err = chain.Validate("")
	if err != nil {
		t.Errorf("unexpected error for empty input with empty chain: %v", err)
	}
}
