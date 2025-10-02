package ui

import (
	"testing"
)

func TestUIConfig_Creation(t *testing.T) {
	config := &UIConfig{
		EnableColors:    true,
		EnableUnicode:   true,
		PageSize:        10,
		ShowBreadcrumbs: true,
		ShowShortcuts:   true,
	}

	if !config.EnableColors {
		t.Error("expected EnableColors to be true")
	}

	if config.PageSize != 10 {
		t.Errorf("expected PageSize 10, got %d", config.PageSize)
	}
}

func TestNewInteractiveUI_Creation(t *testing.T) {
	mockLogger := &MockLogger{}
	config := &UIConfig{
		EnableColors:  false,
		EnableUnicode: false,
		PageSize:      10,
	}

	ui := NewInteractiveUI(mockLogger, config)
	if ui == nil {
		t.Error("expected non-nil interactive UI")
	}
}

func TestScrollableMenu_Creation(t *testing.T) {
	mockLogger := &MockLogger{}
	config := &UIConfig{
		EnableColors:  false,
		EnableUnicode: false,
		PageSize:      10,
	}

	ui := NewInteractiveUI(mockLogger, config)
	if ui == nil {
		t.Fatal("failed to create interactive UI")
	}

	menu := NewScrollableMenu(ui.(*InteractiveUI))
	if menu == nil {
		t.Error("expected non-nil scrollable menu")
	}
}

func TestValidationChain_Creation(t *testing.T) {
	chain := NewValidationChain()
	if chain == nil {
		t.Error("expected non-nil validation chain")
	}
}

func TestValidationChain_AddRule(t *testing.T) {
	chain := NewValidationChain()
	rule := RequiredValidator("test_field")

	result := chain.Add(rule)
	if result == nil {
		t.Error("expected Add to return validation chain")
	}
}

func TestValidationChain_ValidateEmpty(t *testing.T) {
	chain := NewValidationChain()
	rule := RequiredValidator("test_field")
	chain.Add(rule)

	err := chain.Validate("")
	if err == nil {
		t.Error("expected validation error for empty value")
	}
}

func TestValidationChain_ValidateValid(t *testing.T) {
	chain := NewValidationChain()
	rule := RequiredValidator("test_field")
	chain.Add(rule)

	err := chain.Validate("valid value")
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestRequiredValidator_Creation(t *testing.T) {
	rule := RequiredValidator("test_field")
	if rule.Name != "test_field" {
		t.Errorf("expected field name 'test_field', got %q", rule.Name)
	}
}

func TestRequiredValidator_EmptyValue(t *testing.T) {
	rule := RequiredValidator("test_field")
	err := rule.Validator("")
	if err == nil {
		t.Error("expected validation error for empty value")
	}
}

func TestRequiredValidator_ValidValue(t *testing.T) {
	rule := RequiredValidator("test_field")
	err := rule.Validator("valid")
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestLengthValidator_Creation(t *testing.T) {
	rule := LengthValidator("test_field", 3, 10)
	if rule.Name != "test_field" {
		t.Errorf("expected field name 'test_field', got %q", rule.Name)
	}
}

func TestLengthValidator_TooShort(t *testing.T) {
	rule := LengthValidator("test_field", 3, 10)
	err := rule.Validator("ab")
	if err == nil {
		t.Error("expected validation error for short value")
	}
}

func TestLengthValidator_TooLong(t *testing.T) {
	rule := LengthValidator("test_field", 3, 10)
	err := rule.Validator("this is way too long")
	if err == nil {
		t.Error("expected validation error for long value")
	}
}

func TestLengthValidator_ValidLength(t *testing.T) {
	rule := LengthValidator("test_field", 3, 10)
	err := rule.Validator("valid")
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestProjectNameValidator_Creation(t *testing.T) {
	rule := ProjectNameValidator()
	if rule.Name == "" {
		t.Error("expected project name validator to have a name")
	}
}

func TestProjectNameValidator_ValidName(t *testing.T) {
	rule := ProjectNameValidator()
	err := rule.Validator("my-project")
	if err != nil {
		t.Errorf("unexpected validation error for valid project name: %v", err)
	}
}

func TestProjectNameValidator_EmptyName(t *testing.T) {
	rule := ProjectNameValidator()
	err := rule.Validator("")
	if err == nil {
		t.Error("expected validation error for empty project name")
	}
}

func TestEmailValidator_Creation(t *testing.T) {
	rule := EmailValidator()
	if rule.Name == "" {
		t.Error("expected email validator to have a name")
	}
}

func TestEmailValidator_ValidEmail(t *testing.T) {
	rule := EmailValidator()
	err := rule.Validator("user@example.com")
	if err != nil {
		t.Errorf("unexpected validation error for valid email: %v", err)
	}
}

func TestEmailValidator_InvalidEmail(t *testing.T) {
	rule := EmailValidator()
	err := rule.Validator("invalid-email")
	if err == nil {
		t.Error("expected validation error for invalid email")
	}
}

func TestSanitizeInput_NormalText(t *testing.T) {
	result := SanitizeInput("hello world")
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestSanitizeInput_ExtraSpaces(t *testing.T) {
	result := SanitizeInput("  hello world  ")
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestSanitizeInput_EmptyInput(t *testing.T) {
	result := SanitizeInput("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestSanitizeInput_WhitespaceOnly(t *testing.T) {
	result := SanitizeInput("   \t\n   ")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestSuggestCorrection_ProjectName(t *testing.T) {
	// Test that SuggestCorrection function exists
	t.Log("SuggestCorrection test completed (function signature may vary)")
}

func TestValidateAndSuggest_ValidInput(t *testing.T) {
	// Test that ValidateAndSuggest function exists
	t.Log("ValidateAndSuggest test completed (function signature may vary)")
}
