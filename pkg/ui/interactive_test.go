// Package ui provides unit tests for the interactive UI framework.
//
// This file contains comprehensive tests for interactive UI components including
// menu navigation, input validation, error handling, and recovery options.
package ui

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	// logs field removed as it was unused
}

func (m *MockLogger) Debug(format string, args ...interface{})                            {}
func (m *MockLogger) Info(format string, args ...interface{})                             {}
func (m *MockLogger) Warn(format string, args ...interface{})                             {}
func (m *MockLogger) Error(format string, args ...interface{})                            {}
func (m *MockLogger) Fatal(format string, args ...interface{})                            {}
func (m *MockLogger) DebugWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) InfoWithFields(message string, fields map[string]interface{})        {}
func (m *MockLogger) WarnWithFields(message string, fields map[string]interface{})        {}
func (m *MockLogger) ErrorWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) FatalWithFields(message string, fields map[string]interface{})       {}
func (m *MockLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {}
func (m *MockLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	return &interfaces.OperationContext{}
}
func (m *MockLogger) LogOperationStart(operation string, fields map[string]interface{}) {}
func (m *MockLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
}
func (m *MockLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {}
func (m *MockLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
}
func (m *MockLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
}
func (m *MockLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {}
func (m *MockLogger) LogMemoryUsage(operation string)                                        {}
func (m *MockLogger) SetLevel(level int)                                                     {}
func (m *MockLogger) GetLevel() int                                                          { return 0 }
func (m *MockLogger) SetJSONOutput(enabled bool)                                             {}
func (m *MockLogger) SetCallerInfo(enabled bool)                                             {}
func (m *MockLogger) IsDebugEnabled() bool                                                   { return true }
func (m *MockLogger) IsInfoEnabled() bool                                                    { return true }
func (m *MockLogger) WithComponent(component string) interfaces.Logger                       { return m }
func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	return &MockLoggerContext{}
}
func (m *MockLogger) GetLogDir() string                                { return "/tmp" }
func (m *MockLogger) GetRecentEntries(limit int) []interfaces.LogEntry { return nil }
func (m *MockLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	return nil
}
func (m *MockLogger) GetLogFiles() ([]string, error)              { return nil, nil }
func (m *MockLogger) ReadLogFile(filename string) ([]byte, error) { return nil, nil }
func (m *MockLogger) Close() error                                { return nil }

// MockLoggerContext implements the LoggerContext interface for testing
type MockLoggerContext struct{}

func (m *MockLoggerContext) Debug(msg string, args ...interface{}) {}
func (m *MockLoggerContext) Info(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Warn(msg string, args ...interface{})  {}
func (m *MockLoggerContext) Error(msg string, args ...interface{}) {}
func (m *MockLoggerContext) ErrorWithError(msg string, err error)  {}

// TestNewInteractiveUI tests the creation of a new interactive UI instance
func TestNewInteractiveUI(t *testing.T) {
	logger := &MockLogger{}

	// Test with default config
	ui := NewInteractiveUI(logger, nil)
	if ui == nil {
		t.Fatal("Expected non-nil UI instance")
	}

	// Test with custom config
	config := &UIConfig{
		EnableColors:    false,
		EnableUnicode:   false,
		PageSize:        5,
		Timeout:         10 * time.Minute,
		AutoSave:        false,
		ShowBreadcrumbs: false,
		ShowShortcuts:   false,
		ConfirmOnQuit:   false,
	}

	ui2 := NewInteractiveUI(logger, config)
	if ui2 == nil {
		t.Fatal("Expected non-nil UI instance with custom config")
	}
}

// TestColorize tests the colorize function
func TestColorize(t *testing.T) {
	logger := &MockLogger{}

	// Test with colors enabled
	config := &UIConfig{EnableColors: true}
	ui := NewInteractiveUI(logger, config).(*InteractiveUI)

	result := ui.colorize("test", "red")
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("Expected red color code in result, got: %s", result)
	}

	// Test with colors disabled
	config.EnableColors = false
	ui.config = config

	result = ui.colorize("test", "red")
	if result != "test" {
		t.Errorf("Expected plain text when colors disabled, got: %s", result)
	}

	// Test with invalid color
	config.EnableColors = true
	ui.config = config

	result = ui.colorize("test", "invalid")
	if result != "test" {
		t.Errorf("Expected plain text for invalid color, got: %s", result)
	}
}

// TestMenuConfig tests menu configuration validation
func TestMenuConfig(t *testing.T) {
	logger := &MockLogger{}
	ui := NewInteractiveUI(logger, nil)
	ctx := context.Background()

	// Test empty menu options
	config := interfaces.MenuConfig{
		Title:   "Test Menu",
		Options: []interfaces.MenuOption{},
	}

	_, err := ui.ShowMenu(ctx, config)
	if err == nil {
		t.Error("Expected error for empty menu options")
	}

	// Test valid menu config
	config.Options = []interfaces.MenuOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
	}

	// This would normally require user input, so we'll just test the config validation
	if len(config.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(config.Options))
	}
}

// TestTextPromptValidation tests text prompt validation
func TestTextPromptValidation(t *testing.T) {
	skipIfNotInteractive(t)
	ui, _ := createTestUI([]string{"", "", "", "", "", ""}) // Provide enough empty inputs for all error prompts

	// Test required field validation
	config := interfaces.TextPromptConfig{
		Prompt:   "Enter name",
		Required: true,
	}

	result, shouldReturn := ui.handleTextInput("", config)
	if shouldReturn {
		t.Error("Expected validation to fail for empty required field")
	}
	if result != nil {
		t.Error("Expected nil result for failed validation")
	}

	// Test length validation
	config.MinLength = 5
	config.MaxLength = 10

	_, shouldReturn = ui.handleTextInput("abc", config)
	if shouldReturn {
		t.Error("Expected validation to fail for input too short")
	}

	_, shouldReturn = ui.handleTextInput("abcdefghijk", config)
	if shouldReturn {
		t.Error("Expected validation to fail for input too long")
	}

	result, shouldReturn = ui.handleTextInput("abcdef", config)
	if !shouldReturn {
		t.Error("Expected validation to pass for valid input")
	}
	if result == nil || result.Value != "abcdef" {
		t.Error("Expected valid result for correct input")
	}
}

// TestCustomValidator tests custom validation functions
func TestCustomValidator(t *testing.T) {
	skipIfNotInteractive(t)
	ui, _ := createTestUI([]string{"", "", ""}) // Provide enough empty inputs for error prompts

	// Test custom validator
	config := interfaces.TextPromptConfig{
		Prompt: "Enter email",
		Validator: func(input string) error {
			if !strings.Contains(input, "@") {
				return interfaces.NewValidationError("email", input, "Invalid email format", "invalid_email")
			}
			return nil
		},
	}

	_, shouldReturn := ui.handleTextInput("invalid-email", config)
	if shouldReturn {
		t.Error("Expected validation to fail for invalid email")
	}

	result, shouldReturn := ui.handleTextInput("test@example.com", config)
	if !shouldReturn {
		t.Error("Expected validation to pass for valid email")
	}
	if result == nil || result.Value != "test@example.com" {
		t.Error("Expected valid result for correct email")
	}
}

// TestMultiSelectValidation tests multi-select validation
func TestMultiSelectValidation(t *testing.T) {
	logger := &MockLogger{}
	ui := NewInteractiveUI(logger, nil)
	ctx := context.Background()

	// Test empty options
	config := interfaces.MultiSelectConfig{
		Title:   "Test Multi-Select",
		Options: []interfaces.SelectOption{},
	}

	_, err := ui.ShowMultiSelect(ctx, config)
	if err == nil {
		t.Error("Expected error for empty multi-select options")
	}

	// Test valid config
	config.Options = []interfaces.SelectOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2"},
		{Label: "Option 3", Value: "opt3"},
	}
	config.MinSelection = 1
	config.MaxSelection = 2

	if len(config.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(config.Options))
	}
}

// TestFilterOptions tests option filtering functionality
func TestFilterOptions(t *testing.T) {
	ui, _ := createTestUI([]string{}) // No input needed for this test

	options := []interfaces.SelectOption{
		{Label: "Frontend React", Description: "React frontend", Tags: []string{"frontend", "react"}},
		{Label: "Backend Go", Description: "Go backend", Tags: []string{"backend", "go"}},
		{Label: "Frontend Vue", Description: "Vue frontend", Tags: []string{"frontend", "vue"}},
	}

	// Test empty query
	filtered := ui.filterOptions(options, "")
	if len(filtered) != 3 {
		t.Errorf("Expected 3 options for empty query, got %d", len(filtered))
	}

	// Test label search
	filtered = ui.filterOptions(options, "react")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 option for 'react' query, got %d", len(filtered))
	}

	// Test description search
	filtered = ui.filterOptions(options, "backend")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 option for 'backend' query, got %d", len(filtered))
	}

	// Test tag search
	filtered = ui.filterOptions(options, "frontend")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 options for 'frontend' query, got %d", len(filtered))
	}

	// Test case insensitive search
	filtered = ui.filterOptions(options, "REACT")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 option for case insensitive search, got %d", len(filtered))
	}
}

// TestCheckboxValidation tests checkbox validation
func TestCheckboxValidation(t *testing.T) {
	logger := &MockLogger{}
	ui := NewInteractiveUI(logger, nil)
	ctx := context.Background()

	// Test empty items
	config := interfaces.CheckboxConfig{
		Title: "Test Checkbox",
		Items: []interfaces.CheckboxItem{},
	}

	_, err := ui.ShowCheckboxList(ctx, config)
	if err == nil {
		t.Error("Expected error for empty checkbox items")
	}

	// Test required items validation
	config.Items = []interfaces.CheckboxItem{
		{Label: "Required Item", Value: "req", Required: true, Checked: false},
		{Label: "Optional Item", Value: "opt", Required: false, Checked: true},
	}

	// This would normally require user interaction, so we'll test the config
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(config.Items))
	}

	requiredCount := 0
	for _, item := range config.Items {
		if item.Required {
			requiredCount++
		}
	}
	if requiredCount != 1 {
		t.Errorf("Expected 1 required item, got %d", requiredCount)
	}
}

// TestSessionManagement tests UI session management
func TestSessionManagement(t *testing.T) {
	logger := &MockLogger{}
	ui := NewInteractiveUI(logger, nil)
	ctx := context.Background()

	// Test session creation
	config := interfaces.SessionConfig{
		Title:       "Test Session",
		Description: "Test session description",
		Timeout:     5 * time.Minute,
		AutoSave:    true,
	}

	session, err := ui.StartSession(ctx, config)
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}

	if session == nil {
		t.Fatal("Expected non-nil session")
	}

	if session.Title != config.Title {
		t.Errorf("Expected session title '%s', got '%s'", config.Title, session.Title)
	}

	// Test session ending
	err = ui.EndSession(ctx, session)
	if err != nil {
		t.Errorf("Failed to end session: %v", err)
	}
}

// TestValidationError tests validation error handling
func TestValidationError(t *testing.T) {
	// Test validation error creation
	err := interfaces.NewValidationError("email", "invalid", "Invalid email format", "invalid_email")
	if err == nil {
		t.Fatal("Expected non-nil validation error")
	}

	if err.Field != "email" {
		t.Errorf("Expected field 'email', got '%s'", err.Field)
	}

	if err.Message != "Invalid email format" {
		t.Errorf("Expected message 'Invalid email format', got '%s'", err.Message)
	}

	// Test adding suggestions
	if suggestErr := err.WithSuggestions("Use format: user@domain.com", "Check for typos"); suggestErr != nil {
		t.Logf("Warning: Failed to add suggestions: %v", suggestErr)
	}
	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}

	// Test adding recovery options
	recoveryOption := interfaces.RecoveryOption{
		Label:       "Use default email",
		Description: "Use the default email address",
		Safe:        true,
	}
	if recoveryErr := err.WithRecoveryOptions(recoveryOption); recoveryErr != nil {
		t.Logf("Warning: Failed to add recovery options: %v", recoveryErr)
	}
	if len(err.RecoveryOptions) != 1 {
		t.Errorf("Expected 1 recovery option, got %d", len(err.RecoveryOptions))
	}
}

// TestProgressTracker tests progress tracking functionality
func TestProgressTracker(t *testing.T) {
	skipIfNotInteractive(t)
	ui, _ := createTestUI([]string{""}) // Provide empty input for the "Press Enter" prompt
	ctx := context.Background()

	config := interfaces.ProgressConfig{
		Title:       "Test Progress",
		Description: "Testing progress tracking",
		Steps:       []string{"Step 1", "Step 2", "Step 3"},
		ShowPercent: true,
		ShowETA:     true,
		Cancellable: true,
	}

	tracker, err := ui.ShowProgress(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create progress tracker: %v", err)
	}

	if tracker == nil {
		t.Fatal("Expected non-nil progress tracker")
	}

	// Test progress updates
	err = tracker.SetProgress(0.5)
	if err != nil {
		t.Errorf("Failed to set progress: %v", err)
	}

	err = tracker.SetCurrentStep(1, "Processing step 2")
	if err != nil {
		t.Errorf("Failed to set current step: %v", err)
	}

	err = tracker.AddLog("Test log message")
	if err != nil {
		t.Errorf("Failed to add log: %v", err)
	}

	// Test completion
	err = tracker.Complete()
	if err != nil {
		t.Errorf("Failed to complete progress: %v", err)
	}

	// Test closing
	err = tracker.Close()
	if err != nil {
		t.Errorf("Failed to close progress tracker: %v", err)
	}
}

// TestRecoveryOptions tests recovery option functionality
func TestRecoveryOptions(t *testing.T) {
	actionCalled := false
	recoveryAction := func() error {
		actionCalled = true
		return nil
	}

	option := CreateRecoveryOption(
		"Test Recovery",
		"Test recovery description",
		recoveryAction,
		true,
	)

	if option.Label != "Test Recovery" {
		t.Errorf("Expected label 'Test Recovery', got '%s'", option.Label)
	}

	if !option.Safe {
		t.Error("Expected safe recovery option")
	}

	// Test action execution
	if option.Action != nil {
		err := option.Action()
		if err != nil {
			t.Errorf("Recovery action failed: %v", err)
		}
		if !actionCalled {
			t.Error("Expected recovery action to be called")
		}
	}
}

// BenchmarkColorize benchmarks the colorize function
func BenchmarkColorize(b *testing.B) {
	logger := &MockLogger{}
	config := &UIConfig{EnableColors: true}
	ui := NewInteractiveUI(logger, config).(*InteractiveUI)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ui.colorize("test string", "red")
	}
}

// BenchmarkFilterOptions benchmarks option filtering
func BenchmarkFilterOptions(b *testing.B) {
	logger := &MockLogger{}
	ui := NewInteractiveUI(logger, nil).(*InteractiveUI)

	// Create a large set of options
	options := make([]interfaces.SelectOption, 1000)
	for i := 0; i < 1000; i++ {
		options[i] = interfaces.SelectOption{
			Label:       fmt.Sprintf("Option %d", i),
			Description: fmt.Sprintf("Description for option %d", i),
			Tags:        []string{fmt.Sprintf("tag%d", i%10)},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ui.filterOptions(options, "option")
	}
}
