package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		context string
		want    string
	}{
		{"nil error", nil, "test context", ""},
		{"with error", errors.New("original error"), "test context", "test context: original error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleError(tt.err, tt.context)
			if tt.want == "" {
				if result != nil {
					t.Errorf("HandleError() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("HandleError() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")

	tests := []struct {
		name   string
		err    error
		format string
		args   []interface{}
		want   string
	}{
		{"nil error", nil, "context %s", []interface{}{"test"}, ""},
		{"with error", originalErr, "context %s", []interface{}{"test"}, "context test: original error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.format, tt.args...)
			if tt.want == "" {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("WrapError() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestIsNilError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		wantError bool
		wantMsg   string
	}{
		{"nil error", nil, "test operation", false, ""},
		{"with error", errors.New("test error"), "test operation", true, "failed to test operation: test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNilError(tt.err, tt.operation)
			if (result != nil) != tt.wantError {
				t.Errorf("IsNilError() error = %v, wantError %v", result, tt.wantError)
			}
			if result != nil && !strings.Contains(result.Error(), tt.wantMsg) {
				t.Errorf("IsNilError() = %v, want to contain %v", result.Error(), tt.wantMsg)
			}
		})
	}
}

func TestNewErrorContext(t *testing.T) {
	operation := "test_operation"
	component := "test_component"

	ctx := NewErrorContext(operation, component)

	if ctx.Operation != operation {
		t.Errorf("Operation = %v, want %v", ctx.Operation, operation)
	}
	if ctx.Component != component {
		t.Errorf("Component = %v, want %v", ctx.Component, component)
	}
	if ctx.File == "" {
		t.Error("File should not be empty")
	}
	if ctx.Line == 0 {
		t.Error("Line should not be zero")
	}
	if ctx.Metadata == nil {
		t.Error("Metadata should be initialized")
	}
}

func TestErrorContextMethods(t *testing.T) {
	ctx := NewErrorContext("test_op", "test_comp")

	// Test WithMetadata
	ctx = ctx.WithMetadata("key1", "value1")
	if ctx.Metadata["key1"] != "value1" {
		t.Errorf("Metadata key1 = %v, want value1", ctx.Metadata["key1"])
	}

	// Test WithUserID
	userID := "user123"
	ctx = ctx.WithUserID(userID)
	if ctx.UserID != userID {
		t.Errorf("UserID = %v, want %v", ctx.UserID, userID)
	}

	// Test WithRequestID
	requestID := "req456"
	ctx = ctx.WithRequestID(requestID)
	if ctx.RequestID != requestID {
		t.Errorf("RequestID = %v, want %v", ctx.RequestID, requestID)
	}
}

func TestWrapErrorWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	ctx := &ErrorContext{
		Operation: "test_operation",
		Component: "test_component",
	}

	tests := []struct {
		name    string
		err     error
		ctx     *ErrorContext
		message string
		want    string
	}{
		{"nil error", nil, ctx, "test message", ""},
		{"with error", originalErr, ctx, "test message", "component=test_component operation=test_operation test message: original error"},
		{"empty message", originalErr, ctx, "", "component=test_component operation=test_operation: original error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapErrorWithContext(tt.err, tt.ctx, tt.message)
			if tt.want == "" {
				if result != nil {
					t.Errorf("WrapErrorWithContext() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("WrapErrorWithContext() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	field := "test_field"
	message := "test message"
	value := "test_value"

	err := NewValidationError(field, message, value)

	if err.Type != models.ValidationErrorType {
		t.Errorf("Type = %v, want %v", err.Type, models.ValidationErrorType)
	}
	if !strings.Contains(err.Message, field) {
		t.Errorf("Message should contain field name %s", field)
	}
	if !strings.Contains(err.Message, message) {
		t.Errorf("Message should contain message %s", message)
	}
}

func TestNewTemplateError(t *testing.T) {
	templatePath := "/path/to/template"
	operation := "parse"
	message := "syntax error"
	cause := errors.New("cause error")

	err := NewTemplateError(templatePath, operation, message, cause)

	if err.Type != models.TemplateErrorType {
		t.Errorf("Type = %v, want %v", err.Type, models.TemplateErrorType)
	}
	if !strings.Contains(err.Message, templatePath) {
		t.Errorf("Message should contain template path %s", templatePath)
	}
	if !strings.Contains(err.Message, operation) {
		t.Errorf("Message should contain operation %s", operation)
	}
	if !strings.Contains(err.Message, message) {
		t.Errorf("Message should contain message %s", message)
	}
}

func TestNewFileSystemError(t *testing.T) {
	path := "/path/to/file"
	operation := "write"
	message := "permission denied"
	cause := errors.New("cause error")

	err := NewFileSystemError(path, operation, message, cause)

	if err.Type != models.FileSystemErrorType {
		t.Errorf("Type = %v, want %v", err.Type, models.FileSystemErrorType)
	}
	if !strings.Contains(err.Message, path) {
		t.Errorf("Message should contain path %s", path)
	}
	if !strings.Contains(err.Message, operation) {
		t.Errorf("Message should contain operation %s", operation)
	}
	if !strings.Contains(err.Message, message) {
		t.Errorf("Message should contain message %s", message)
	}
}

func TestNewConfigurationError(t *testing.T) {
	configKey := "database.host"
	message := "invalid host"
	cause := errors.New("cause error")

	err := NewConfigurationError(configKey, message, cause)

	if err.Type != models.ConfigurationErrorType {
		t.Errorf("Type = %v, want %v", err.Type, models.ConfigurationErrorType)
	}
	if !strings.Contains(err.Message, configKey) {
		t.Errorf("Message should contain config key %s", configKey)
	}
	if !strings.Contains(err.Message, message) {
		t.Errorf("Message should contain message %s", message)
	}
}

func TestValidateAndWrapError(t *testing.T) {
	ctx := &ErrorContext{
		Operation: "test_operation",
		Component: "test_component",
	}

	tests := []struct {
		name string
		err  error
		ctx  *ErrorContext
		want string
	}{
		{"nil error", nil, ctx, ""},
		{"generator error", models.NewGeneratorError(models.ValidationErrorType, "validation failed", nil), ctx, "validation failed"},
		{"regular error", errors.New("regular error"), ctx, "component=test_component operation=test_operation unexpected error occurred: regular error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAndWrapError(tt.err, tt.ctx)
			if tt.want == "" {
				if result != nil {
					t.Errorf("ValidateAndWrapError() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("ValidateAndWrapError() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestFormatErrorForUser(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil error", nil, ""},
		{"validation error", models.NewGeneratorError(models.ValidationErrorType, "validation failed", nil), "❌ Validation Error: validation failed"},
		{"template error", models.NewGeneratorError(models.TemplateErrorType, "template failed", nil), "📄 Template Error: template failed"},
		{"filesystem error", models.NewGeneratorError(models.FileSystemErrorType, "file operation failed", nil), "📁 File System Error: file operation failed"},
		{"configuration error", models.NewGeneratorError(models.ConfigurationErrorType, "config invalid", nil), "⚙️ Configuration Error: config invalid"},
		{"regular error", errors.New("regular error"), "Error: regular error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatErrorForUser(tt.err)
			if got != tt.want {
				t.Errorf("FormatErrorForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
