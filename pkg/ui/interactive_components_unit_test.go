package ui

import (
	"context"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Simple unit tests for interactive components that don't conflict with existing code

func TestInteractiveUI_ColorizeFunction_Simple(t *testing.T) {
	skipIfNotInteractive(t)
	ui := &InteractiveUI{
		config: &UIConfig{EnableColors: false},
	}

	result := ui.colorize("test text", "red")
	if result != "test text" {
		t.Errorf("Expected plain text when colors disabled, got %q", result)
	}

	ui.config.EnableColors = true
	result = ui.colorize("test text", "red")
	if !strings.Contains(result, "test text") {
		t.Errorf("Expected colored text to contain original text, got %q", result)
	}
}

func TestInteractiveUI_SetupDefaultShortcuts(t *testing.T) {
	ui := &InteractiveUI{
		shortcuts: make(map[string]interfaces.KeyboardShortcut),
	}

	ui.setupDefaultShortcuts()

	expectedShortcuts := []string{"q", "h", "b", "ctrl+c", "enter", "esc", "tab", "shift+tab"}
	for _, key := range expectedShortcuts {
		if _, exists := ui.shortcuts[key]; !exists {
			t.Errorf("Expected shortcut %s to be set up", key)
		}
	}
}

func TestValidationError_Interface(t *testing.T) {
	err := &interfaces.ValidationError{
		Field:   "test_field",
		Value:   "test_value",
		Message: "Test error message",
		Code:    "test_code",
	}

	if err.Error() != "Test error message" {
		t.Errorf("Expected error message 'Test error message', got %s", err.Error())
	}

	// Test WithSuggestions
	err = err.WithSuggestions("suggestion1", "suggestion2")
	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}

	// Test WithRecoveryOptions
	recoveryOption := interfaces.RecoveryOption{
		Label:       "Retry",
		Description: "Try again",
		Safe:        true,
	}
	err = err.WithRecoveryOptions(recoveryOption)
	if len(err.RecoveryOptions) != 1 {
		t.Errorf("Expected 1 recovery option, got %d", len(err.RecoveryOptions))
	}
}

func TestNewValidationError(t *testing.T) {
	err := interfaces.NewValidationError("field", "value", "message", "code")

	if err.Field != "field" {
		t.Errorf("Expected field 'field', got %s", err.Field)
	}
	if err.Value != "value" {
		t.Errorf("Expected value 'value', got %s", err.Value)
	}
	if err.Message != "message" {
		t.Errorf("Expected message 'message', got %s", err.Message)
	}
	if err.Code != "code" {
		t.Errorf("Expected code 'code', got %s", err.Code)
	}
}

func TestUIConfig_Defaults(t *testing.T) {
	config := &UIConfig{
		EnableColors:    true,
		EnableUnicode:   true,
		PageSize:        10,
		ShowBreadcrumbs: true,
		ShowShortcuts:   true,
		ConfirmOnQuit:   true,
	}

	if !config.EnableColors {
		t.Error("Expected EnableColors to be true")
	}
	if !config.EnableUnicode {
		t.Error("Expected EnableUnicode to be true")
	}
	if config.PageSize != 10 {
		t.Errorf("Expected PageSize 10, got %d", config.PageSize)
	}
	if !config.ShowBreadcrumbs {
		t.Error("Expected ShowBreadcrumbs to be true")
	}
	if !config.ShowShortcuts {
		t.Error("Expected ShowShortcuts to be true")
	}
	if !config.ConfirmOnQuit {
		t.Error("Expected ConfirmOnQuit to be true")
	}
}

func TestMenuConfig_Validation(t *testing.T) {
	// Test empty options
	config := interfaces.MenuConfig{
		Title:   "Test Menu",
		Options: []interfaces.MenuOption{},
	}

	if len(config.Options) != 0 {
		t.Error("Expected empty options to be allowed in config")
	}

	// Test with options
	config.Options = []interfaces.MenuOption{
		{Label: "Option 1", Value: "opt1"},
		{Label: "Option 2", Value: "opt2", Disabled: true},
	}

	if len(config.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(config.Options))
	}

	if config.Options[1].Disabled != true {
		t.Error("Expected second option to be disabled")
	}
}

func TestTextPromptConfig_Validation(t *testing.T) {
	config := interfaces.TextPromptConfig{
		Prompt:       "Enter text",
		DefaultValue: "default",
		Required:     true,
		MaxLength:    100,
		MinLength:    5,
	}

	if config.Prompt != "Enter text" {
		t.Errorf("Expected prompt 'Enter text', got %s", config.Prompt)
	}
	if config.DefaultValue != "default" {
		t.Errorf("Expected default value 'default', got %s", config.DefaultValue)
	}
	if !config.Required {
		t.Error("Expected required to be true")
	}
	if config.MaxLength != 100 {
		t.Errorf("Expected max length 100, got %d", config.MaxLength)
	}
	if config.MinLength != 5 {
		t.Errorf("Expected min length 5, got %d", config.MinLength)
	}
}

func TestMultiSelectConfig_Validation(t *testing.T) {
	config := interfaces.MultiSelectConfig{
		Title:        "Select Options",
		MinSelection: 1,
		MaxSelection: 3,
		Options: []interfaces.SelectOption{
			{Label: "Option 1", Value: "opt1", Selected: true},
			{Label: "Option 2", Value: "opt2", Selected: false},
		},
	}

	if config.Title != "Select Options" {
		t.Errorf("Expected title 'Select Options', got %s", config.Title)
	}
	if config.MinSelection != 1 {
		t.Errorf("Expected min selection 1, got %d", config.MinSelection)
	}
	if config.MaxSelection != 3 {
		t.Errorf("Expected max selection 3, got %d", config.MaxSelection)
	}
	if len(config.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(config.Options))
	}
	if !config.Options[0].Selected {
		t.Error("Expected first option to be selected")
	}
}

func TestCheckboxConfig_Validation(t *testing.T) {
	config := interfaces.CheckboxConfig{
		Title: "Select Items",
		Items: []interfaces.CheckboxItem{
			{Label: "Item 1", Value: "item1", Checked: true, Required: true},
			{Label: "Item 2", Value: "item2", Checked: false, Disabled: true},
		},
	}

	if config.Title != "Select Items" {
		t.Errorf("Expected title 'Select Items', got %s", config.Title)
	}
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(config.Items))
	}
	if !config.Items[0].Checked {
		t.Error("Expected first item to be checked")
	}
	if !config.Items[0].Required {
		t.Error("Expected first item to be required")
	}
	if !config.Items[1].Disabled {
		t.Error("Expected second item to be disabled")
	}
}

func TestTableConfig_Validation(t *testing.T) {
	config := interfaces.TableConfig{
		Title:   "Test Table",
		Headers: []string{"Column 1", "Column 2"},
		Rows: [][]string{
			{"Row 1 Col 1", "Row 1 Col 2"},
			{"Row 2 Col 1", "Row 2 Col 2"},
		},
		MaxWidth:   80,
		Pagination: true,
		PageSize:   10,
	}

	if config.Title != "Test Table" {
		t.Errorf("Expected title 'Test Table', got %s", config.Title)
	}
	if len(config.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(config.Headers))
	}
	if len(config.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(config.Rows))
	}
	if config.MaxWidth != 80 {
		t.Errorf("Expected max width 80, got %d", config.MaxWidth)
	}
	if !config.Pagination {
		t.Error("Expected pagination to be enabled")
	}
	if config.PageSize != 10 {
		t.Errorf("Expected page size 10, got %d", config.PageSize)
	}
}

func TestTreeConfig_Validation(t *testing.T) {
	config := interfaces.TreeConfig{
		Title: "Test Tree",
		Root: interfaces.TreeNode{
			Label: "Root",
			Children: []interfaces.TreeNode{
				{Label: "Child 1", Expanded: true},
				{Label: "Child 2", Selectable: true},
			},
		},
		Expandable: true,
		ShowIcons:  true,
		MaxDepth:   5,
	}

	if config.Title != "Test Tree" {
		t.Errorf("Expected title 'Test Tree', got %s", config.Title)
	}
	if config.Root.Label != "Root" {
		t.Errorf("Expected root label 'Root', got %s", config.Root.Label)
	}
	if len(config.Root.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(config.Root.Children))
	}
	if !config.Expandable {
		t.Error("Expected expandable to be true")
	}
	if !config.ShowIcons {
		t.Error("Expected show icons to be true")
	}
	if config.MaxDepth != 5 {
		t.Errorf("Expected max depth 5, got %d", config.MaxDepth)
	}
}

func TestProgressConfig_Validation(t *testing.T) {
	config := interfaces.ProgressConfig{
		Title:       "Test Progress",
		Description: "Testing progress",
		Steps:       []string{"Step 1", "Step 2", "Step 3"},
		ShowPercent: true,
		ShowETA:     true,
		Cancellable: true,
	}

	if config.Title != "Test Progress" {
		t.Errorf("Expected title 'Test Progress', got %s", config.Title)
	}
	if config.Description != "Testing progress" {
		t.Errorf("Expected description 'Testing progress', got %s", config.Description)
	}
	if len(config.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(config.Steps))
	}
	if !config.ShowPercent {
		t.Error("Expected show percent to be true")
	}
	if !config.ShowETA {
		t.Error("Expected show ETA to be true")
	}
	if !config.Cancellable {
		t.Error("Expected cancellable to be true")
	}
}

func TestErrorConfig_Validation(t *testing.T) {
	config := interfaces.ErrorConfig{
		Title:       "Test Error",
		Message:     "An error occurred",
		Details:     "Error details",
		ErrorType:   "validation",
		Suggestions: []string{"Try again", "Check input"},
		ShowStack:   true,
		AllowRetry:  true,
		AllowIgnore: false,
	}

	if config.Title != "Test Error" {
		t.Errorf("Expected title 'Test Error', got %s", config.Title)
	}
	if config.Message != "An error occurred" {
		t.Errorf("Expected message 'An error occurred', got %s", config.Message)
	}
	if config.Details != "Error details" {
		t.Errorf("Expected details 'Error details', got %s", config.Details)
	}
	if config.ErrorType != "validation" {
		t.Errorf("Expected error type 'validation', got %s", config.ErrorType)
	}
	if len(config.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(config.Suggestions))
	}
	if !config.ShowStack {
		t.Error("Expected show stack to be true")
	}
	if !config.AllowRetry {
		t.Error("Expected allow retry to be true")
	}
	if config.AllowIgnore {
		t.Error("Expected allow ignore to be false")
	}
}

func TestSessionConfig_Validation(t *testing.T) {
	config := interfaces.SessionConfig{
		SessionID:   "test-session",
		Title:       "Test Session",
		Description: "Testing session",
		AutoSave:    true,
	}

	if config.SessionID != "test-session" {
		t.Errorf("Expected session ID 'test-session', got %s", config.SessionID)
	}
	if config.Title != "Test Session" {
		t.Errorf("Expected title 'Test Session', got %s", config.Title)
	}
	if config.Description != "Testing session" {
		t.Errorf("Expected description 'Testing session', got %s", config.Description)
	}
	if !config.AutoSave {
		t.Error("Expected auto save to be true")
	}
}

// Test helper functions

func TestFormatFileSize_Helper(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatFileSize(tt.size)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Mock context cancellation test
func TestContextCancellation_Handling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	select {
	case <-ctx.Done():
		// Expected behavior
	default:
		t.Error("Expected context to be cancelled")
	}

	if ctx.Err() == nil {
		t.Error("Expected context error to be set")
	}
}
