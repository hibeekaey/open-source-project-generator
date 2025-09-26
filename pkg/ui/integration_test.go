// Package ui provides integration tests for the interactive UI framework.
//
// This file contains end-to-end tests that verify the complete interactive UI
// system works correctly with all components integrated together.
package ui

import (
	"bufio"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TestCompleteInteractiveFlow tests a complete interactive workflow
func TestCompleteInteractiveFlow(t *testing.T) {
	logger := &MockLogger{}
	config := &UIConfig{
		EnableColors:    false, // Disable colors for testing
		EnableUnicode:   false,
		PageSize:        5,
		Timeout:         1 * time.Minute,
		AutoSave:        false,
		ShowBreadcrumbs: false,
		ShowShortcuts:   false,
		ConfirmOnQuit:   false,
	}

	ui := NewInteractiveUI(logger, config)
	ctx := context.Background()

	// Test session management
	sessionConfig := interfaces.SessionConfig{
		Title:       "Integration Test Session",
		Description: "Testing complete interactive flow",
		Timeout:     30 * time.Second,
		AutoSave:    false,
	}

	session, err := ui.StartSession(ctx, sessionConfig)
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}

	defer func() {
		if err := ui.EndSession(ctx, session); err != nil {
			t.Errorf("Failed to end session: %v", err)
		}
	}()

	// Test validation chain
	validationChain := ProjectConfigValidation()

	// Test valid input
	validInput := "my-awesome-project"
	sanitized, suggestions, err := ValidateAndSuggest(validInput, validationChain)
	if err != nil {
		t.Errorf("Expected valid input to pass validation, got error: %v", err)
	}
	if sanitized != validInput {
		t.Errorf("Expected sanitized input '%s', got '%s'", validInput, sanitized)
	}
	if suggestions != nil {
		t.Errorf("Expected no suggestions for valid input, got: %v", suggestions)
	}

	// Test invalid input
	invalidInput := "my project with spaces!"
	_, suggestions, err = ValidateAndSuggest(invalidInput, validationChain)
	if err == nil {
		t.Error("Expected invalid input to fail validation")
	}
	if len(suggestions) == 0 {
		t.Error("Expected suggestions for invalid input")
	}

	// Test error handling
	validationErr := interfaces.NewValidationError(
		"project_name",
		invalidInput,
		"Invalid project name format",
		"validation_failed",
	)
	if err := validationErr.WithSuggestions("Use only letters, numbers, hyphens, and underscores"); err != nil {
		t.Logf("Warning: Failed to add suggestions: %v", err)
	}

	errorConfig := interfaces.ErrorConfig{
		Title:       "Validation Error",
		Message:     validationErr.Message,
		ErrorType:   "Validation",
		Suggestions: validationErr.Suggestions,
		AllowRetry:  true,
		AllowBack:   true,
		AllowQuit:   false,
	}

	// This would normally require user interaction, so we just test the config
	if errorConfig.Title != "Validation Error" {
		t.Errorf("Expected error title 'Validation Error', got '%s'", errorConfig.Title)
	}
}

// TestValidationChainIntegration tests validation chain integration
func TestValidationChainIntegration(t *testing.T) {
	// Test project name validation chain
	projectChain := NewValidationChain().
		Add(RequiredValidator("project_name")).
		Add(LengthValidator("project_name", 2, 50)).
		Add(ProjectNameValidator())

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"", false, "empty input"},
		{"a", false, "too short"},
		{"my-project", true, "valid project name"},
		{"my_awesome_project_123", true, "valid with underscores and numbers"},
		{"project-with-very-long-name-that-exceeds-fifty-characters-limit", false, "too long"},
		{"project with spaces", false, "contains spaces"},
		{"project@special", false, "contains special characters"},
		{"-invalid-start", false, "starts with hyphen"},
		{"invalid-end-", false, "ends with hyphen"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := projectChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected input '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected input '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestEmailValidationIntegration tests email validation integration
func TestEmailValidationIntegration(t *testing.T) {
	emailChain := EmailConfigValidation(false) // Not required

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"", true, "empty input (optional)"},
		{"user@example.com", true, "valid email"},
		{"test.email+tag@domain.co.uk", true, "complex valid email"},
		{"invalid-email", false, "missing @ symbol"},
		{"@domain.com", false, "missing username"},
		{"user@", false, "missing domain"},
		{"user@domain", true, "basic domain (valid)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := emailChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected email '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected email '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestURLValidationIntegration tests URL validation integration
func TestURLValidationIntegration(t *testing.T) {
	urlChain := URLConfigValidation(false) // Not required

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"", true, "empty input (optional)"},
		{"https://example.com", true, "valid HTTPS URL"},
		{"http://example.com", true, "valid HTTP URL"},
		{"https://sub.example.com/path", true, "valid URL with subdomain and path"},
		{"ftp://example.com", false, "invalid protocol"},
		{"example.com", false, "missing protocol"},
		{"https://", false, "incomplete URL"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := urlChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected URL '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected URL '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestVersionValidationIntegration tests version validation integration
func TestVersionValidationIntegration(t *testing.T) {
	versionChain := VersionConfigValidation(false) // Not required

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"", true, "empty input (optional)"},
		{"1.0.0", true, "valid semantic version"},
		{"v1.0.0", true, "valid semantic version with v prefix"},
		{"2.1.3-alpha", true, "valid pre-release version"},
		{"1.0.0+build.1", true, "valid version with build metadata"},
		{"1.0.0-alpha.1+build.1", true, "valid complex version"},
		{"1.0", false, "incomplete version"},
		{"v1", false, "incomplete version with prefix"},
		{"1.0.0.0", false, "too many version parts"},
		{"invalid", false, "non-numeric version"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := versionChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected version '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected version '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestNumericValidationIntegration tests numeric validation integration
func TestNumericValidationIntegration(t *testing.T) {
	numericChain := NewValidationChain().
		Add(NumericValidator("port", 1, 65535))

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"", true, "empty input (optional)"},
		{"80", true, "valid port number"},
		{"8080", true, "valid port number"},
		{"65535", true, "maximum valid port"},
		{"1", true, "minimum valid port"},
		{"0", false, "below minimum"},
		{"65536", false, "above maximum"},
		{"abc", false, "non-numeric"},
		{"80.5", true, "decimal number (valid for numeric)"},
		{"-80", false, "negative number"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := numericChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected number '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected number '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestCustomValidatorIntegration tests custom validator integration
func TestCustomValidatorIntegration(t *testing.T) {
	// Create a custom validator for organization names
	orgValidator := CustomValidator(
		"organization",
		"Organization name must be between 2 and 30 characters and contain only letters, numbers, and spaces",
		func(input string) error {
			input = SanitizeInput(input)
			if len(input) < 2 || len(input) > 30 {
				return fmt.Errorf("length must be between 2 and 30 characters")
			}
			for _, r := range input {
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != ' ' {
					return fmt.Errorf("contains invalid characters")
				}
			}
			return nil
		},
		[]string{
			"Use only letters, numbers, and spaces",
			"Length should be between 2 and 30 characters",
		},
		[]interfaces.RecoveryOption{
			CreateDefaultValueRecovery("My Organization"),
		},
	)

	orgChain := NewValidationChain().Add(orgValidator)

	testCases := []struct {
		input       string
		shouldPass  bool
		description string
	}{
		{"My Company", true, "valid organization name"},
		{"Tech Corp 123", true, "valid with numbers"},
		{"A", false, "too short"},
		{"Very Long Organization Name That Exceeds Thirty Characters", false, "too long"},
		{"Company@Inc", false, "contains special characters"},
		{"Company-Inc", false, "contains hyphen"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := orgChain.Validate(tc.input)
			if tc.shouldPass && err != nil {
				t.Errorf("Expected organization '%s' to pass validation, got error: %v", tc.input, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected organization '%s' to fail validation, but it passed", tc.input)
			}
		})
	}
}

// TestSanitizeInputIntegration tests input sanitization integration
func TestSanitizeInputIntegration(t *testing.T) {
	testCases := []struct {
		input       string
		expected    string
		description string
	}{
		{"  hello world  ", "hello world", "trim whitespace"},
		{"hello\r\nworld", "hello\nworld", "normalize line endings"},
		{"hello\rworld", "hello\nworld", "convert CR to LF"},
		{"hello\tworld", "hello\tworld", "preserve tabs"},
		{"hello\x00world", "helloworld", "remove null characters"},
		{"hello\x1bworld", "helloworld", "remove escape characters"},
		{"normal text", "normal text", "preserve normal text"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := SanitizeInput(tc.input)
			if result != tc.expected {
				t.Errorf("Expected sanitized input '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

// TestSuggestCorrectionIntegration tests correction suggestion integration
func TestSuggestCorrectionIntegration(t *testing.T) {
	testCases := []struct {
		input       string
		pattern     string
		expected    string
		description string
	}{
		{"My Project Name", "project_name", "my-project-name", "convert spaces to hyphens"},
		{"Project@Name!", "project_name", "projectname", "remove special characters"},
		{"UPPERCASE", "project_name", "uppercase", "convert to lowercase"},
		{"mixed_Case-Name", "project_name", "mixed_case-name", "preserve valid characters"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := SuggestCorrection(tc.input, tc.pattern)
			if result != tc.expected {
				t.Errorf("Expected suggestion '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

// TestRecoveryOptionsIntegration tests recovery options integration
func TestRecoveryOptionsIntegration(t *testing.T) {
	// Test default value recovery
	defaultRecovery := CreateDefaultValueRecovery("default-project")
	if defaultRecovery.Label != "Use default: default-project" {
		t.Errorf("Expected default recovery label, got '%s'", defaultRecovery.Label)
	}
	if !defaultRecovery.Safe {
		t.Error("Expected default recovery to be safe")
	}

	// Test suggestion recovery
	suggestionRecovery := CreateSuggestionRecovery("suggested-name")
	if suggestionRecovery.Label != "Apply suggestion: suggested-name" {
		t.Errorf("Expected suggestion recovery label, got '%s'", suggestionRecovery.Label)
	}
	if !suggestionRecovery.Safe {
		t.Error("Expected suggestion recovery to be safe")
	}

	// Test skip field recovery
	skipRecovery := CreateSkipFieldRecovery("optional_field")
	if skipRecovery.Label != "Skip optional_field" {
		t.Errorf("Expected skip recovery label, got '%s'", skipRecovery.Label)
	}
	if !skipRecovery.Safe {
		t.Error("Expected skip recovery to be safe")
	}
}

// TestProgressTrackerIntegration tests progress tracker integration
func TestProgressTrackerIntegration(t *testing.T) {
	skipIfNotInteractive(t)
	// Create a test UI with mock input to handle the "Press Enter" prompt
	mockReader := NewMockReader([]string{"", "", ""}) // Provide multiple empty inputs for any prompts
	mockWriter := &MockWriter{}
	mockLogger := &MockLogger{}

	config := &UIConfig{
		EnableColors:    false,
		EnableUnicode:   false,
		PageSize:        10,
		Timeout:         30 * time.Second,
		AutoSave:        false,
		ShowBreadcrumbs: false,
		ShowShortcuts:   false,
		ConfirmOnQuit:   false,
	}

	ui := &InteractiveUI{
		reader:    bufio.NewReader(mockReader),
		writer:    bufio.NewWriter(mockWriter),
		shortcuts: make(map[string]interfaces.KeyboardShortcut),
		logger:    mockLogger,
		config:    config,
	}
	ui.setupDefaultShortcuts()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	progressConfig := interfaces.ProgressConfig{
		Title:       "Integration Test Progress",
		Description: "Testing progress tracking integration",
		Steps:       []string{"Initialize", "Process", "Complete"},
		ShowPercent: true,
		ShowETA:     true,
		Cancellable: true,
	}

	tracker, err := ui.ShowProgress(ctx, progressConfig)
	if err != nil {
		t.Fatalf("Failed to create progress tracker: %v", err)
	}

	// Test progress updates
	steps := []struct {
		progress float64
		step     int
		desc     string
	}{
		{0.0, 0, "Starting initialization"},
		{0.3, 0, "Initializing components"},
		{0.5, 1, "Processing data"},
		{0.8, 1, "Processing complete"},
		{1.0, 2, "Finalizing"},
	}

	for _, step := range steps {
		if err := tracker.SetProgress(step.progress); err != nil {
			t.Errorf("Failed to set progress to %f: %v", step.progress, err)
		}

		if err := tracker.SetCurrentStep(step.step, step.desc); err != nil {
			t.Errorf("Failed to set step %d: %v", step.step, err)
		}

		if err := tracker.AddLog(fmt.Sprintf("Step %d: %s", step.step, step.desc)); err != nil {
			t.Errorf("Failed to add log: %v", err)
		}

		// Small delay to simulate work
		time.Sleep(10 * time.Millisecond)
	}

	// Test completion
	if err := tracker.Complete(); err != nil {
		t.Errorf("Failed to complete progress: %v", err)
	}

	// Test closing
	if err := tracker.Close(); err != nil {
		t.Errorf("Failed to close progress tracker: %v", err)
	}
}

// BenchmarkValidationChain benchmarks validation chain performance
func BenchmarkValidationChain(b *testing.B) {
	chain := NewValidationChain().
		Add(RequiredValidator("test")).
		Add(LengthValidator("test", 2, 50)).
		Add(ProjectNameValidator())

	testInput := "my-test-project"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Validate(testInput)
	}
}

// BenchmarkSanitizeInput benchmarks input sanitization
func BenchmarkSanitizeInput(b *testing.B) {
	testInput := "  My Project Name with\r\nspecial\tcharacters\x00  "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SanitizeInput(testInput)
	}
}
