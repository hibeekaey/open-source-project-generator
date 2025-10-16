package ui

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestProjectConfigValidator_ValidateProjectName tests project name validation
func TestProjectConfigValidator_ValidateProjectName(t *testing.T) {
	validator := &ProjectConfigValidator{
		projectNameRegex: regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`),
	}

	tests := []struct {
		name        string
		projectName string
		expectError bool
		errorCode   string
	}{
		{
			name:        "valid project name",
			projectName: "my-awesome-project",
			expectError: false,
		},
		{
			name:        "valid project name with underscores",
			projectName: "my_awesome_project",
			expectError: false,
		},
		{
			name:        "valid project name with numbers",
			projectName: "project123",
			expectError: false,
		},
		{
			name:        "empty project name",
			projectName: "",
			expectError: true,
			errorCode:   "required",
		},
		{
			name:        "project name too short",
			projectName: "a",
			expectError: true,
			errorCode:   "min_length",
		},
		{
			name:        "project name too long",
			projectName: "this-is-a-very-long-project-name-that-exceeds-the-maximum-allowed-length",
			expectError: true,
			errorCode:   "max_length",
		},
		{
			name:        "project name starts with number",
			projectName: "123project",
			expectError: true,
			errorCode:   "invalid_format",
		},
		{
			name:        "project name with spaces",
			projectName: "my project",
			expectError: true,
			errorCode:   "invalid_format",
		},
		{
			name:        "reserved name",
			projectName: "con",
			expectError: true,
			errorCode:   "reserved_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateProjectName(tt.projectName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				validationErr := &interfaces.ValidationError{}
				if errors.As(err, &validationErr) {
					if validationErr.Code != tt.errorCode {
						t.Errorf("expected error code %q, got %q", tt.errorCode, validationErr.Code)
					}
				} else {
					t.Errorf("expected ValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestProjectConfigValidator_ValidateEmail(t *testing.T) {
	validator := &ProjectConfigValidator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}

	tests := []struct {
		name        string
		email       string
		expectError bool
		errorCode   string
	}{
		{
			name:        "valid email",
			email:       "user@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "empty email (optional)",
			email:       "",
			expectError: false,
		},
		{
			name:        "invalid email format",
			email:       "invalid-email",
			expectError: true,
			errorCode:   "invalid_format",
		},
		{
			name:        "email without domain",
			email:       "user@",
			expectError: true,
			errorCode:   "invalid_format",
		},
		{
			name:        "email too long",
			email:       "this-is-a-very-long-email-address-that-definitely-exceeds-the-maximum-allowed-length-of-one-hundred-characters@example.com",
			expectError: true,
			errorCode:   "max_length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateEmail(tt.email)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				validationErr := &interfaces.ValidationError{}
				if errors.As(err, &validationErr) {
					if validationErr.Code != tt.errorCode {
						t.Errorf("expected error code %q, got %q", tt.errorCode, validationErr.Code)
					}
				} else {
					t.Errorf("expected ValidationError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestProjectConfigCollector_generateMetadataPreview(t *testing.T) {
	// Create a simple mock UI that implements the required interface
	mockUI := &testUI{}
	mockLogger := &MockLogger{}
	collector := NewProjectConfigCollector(mockUI, mockLogger)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Description:  "A test project",
		Author:       "John Doe",
		Email:        "john@example.com",
		Organization: "Test Corp",
		License:      "MIT",
		Repository:   "https://github.com/user/test-project",
	}

	preview := collector.generateMetadataPreview(config)

	// Check that preview contains expected content
	expectedContent := []string{
		"# test-project",
		"A test project",
		"\"name\": \"test-project\"",
		"\"author\": \"John Doe <john@example.com>\"",
		"\"license\": \"MIT\"",
		"MIT License",
		"Copyright (c)",
		"Test Corp",
		"module github.com/user/test-project",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(preview, expected) {
			t.Errorf("preview missing expected content: %q", expected)
		}
	}
}

// testUI is a minimal implementation for testing
type testUI struct{}

func (t *testUI) PromptText(ctx context.Context, config interfaces.TextPromptConfig) (*interfaces.TextResult, error) {
	return &interfaces.TextResult{Value: "test", Action: "submit"}, nil
}

func (t *testUI) ShowMenu(ctx context.Context, config interfaces.MenuConfig) (*interfaces.MenuResult, error) {
	return &interfaces.MenuResult{SelectedIndex: 0, SelectedValue: "MIT", Action: "select"}, nil
}

func (t *testUI) PromptConfirm(ctx context.Context, config interfaces.ConfirmConfig) (*interfaces.ConfirmResult, error) {
	return &interfaces.ConfirmResult{Confirmed: true, Action: "confirm"}, nil
}

func (t *testUI) ShowTable(ctx context.Context, config interfaces.TableConfig) error {
	return nil
}

// Stub implementations for other required methods
func (t *testUI) ShowMultiSelect(ctx context.Context, config interfaces.MultiSelectConfig) (*interfaces.MultiSelectResult, error) {
	return nil, nil
}
func (t *testUI) ShowCheckboxList(ctx context.Context, config interfaces.CheckboxConfig) (*interfaces.CheckboxResult, error) {
	return nil, nil
}
func (t *testUI) PromptSelect(ctx context.Context, config interfaces.SelectConfig) (*interfaces.SelectResult, error) {
	return nil, nil
}
func (t *testUI) ShowTree(ctx context.Context, config interfaces.TreeConfig) error {
	return nil
}
func (t *testUI) ShowProgress(ctx context.Context, config interfaces.ProgressConfig) (interfaces.ProgressTracker, error) {
	return nil, nil
}
func (t *testUI) ShowBreadcrumb(ctx context.Context, path []string) error {
	return nil
}
func (t *testUI) ShowHelp(ctx context.Context, helpContext string) error {
	return nil
}
func (t *testUI) ShowError(ctx context.Context, config interfaces.ErrorConfig) (*interfaces.ErrorResult, error) {
	return nil, nil
}
func (t *testUI) StartSession(ctx context.Context, config interfaces.SessionConfig) (*interfaces.UISession, error) {
	return nil, nil
}
func (t *testUI) EndSession(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (t *testUI) SaveSessionState(ctx context.Context, session *interfaces.UISession) error {
	return nil
}
func (t *testUI) RestoreSessionState(ctx context.Context, sessionID string) (*interfaces.UISession, error) {
	return nil, nil
}
