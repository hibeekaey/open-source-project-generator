package app

import (
	"errors"
	"strings"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "error with cause",
			appError: &AppError{
				Type:    ErrorTypeValidation,
				Message: "validation failed",
				Cause:   errors.New("field is required"),
			},
			expected: "validation failed: field is required",
		},
		{
			name: "error without cause",
			appError: &AppError{
				Type:    ErrorTypeTemplate,
				Message: "template processing failed",
				Cause:   nil,
			},
			expected: "template processing failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.appError.Error(); got != tt.expected {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	appError := &AppError{
		Type:    ErrorTypeInternal,
		Message: "internal error",
		Cause:   cause,
	}

	if unwrapped := appError.Unwrap(); !errors.Is(unwrapped, cause) {
		t.Errorf("AppError.Unwrap() = %v, want %v", unwrapped, cause)
	}

	// Test with no cause
	appErrorNoCause := &AppError{
		Type:    ErrorTypeInternal,
		Message: "internal error",
		Cause:   nil,
	}

	if unwrapped := appErrorNoCause.Unwrap(); unwrapped != nil {
		t.Errorf("AppError.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestNewAppError(t *testing.T) {
	cause := errors.New("test cause")
	appError := NewAppError(ErrorTypeValidation, "test message", cause)

	if appError.Type != ErrorTypeValidation {
		t.Errorf("Expected type %v, got %v", ErrorTypeValidation, appError.Type)
	}

	if appError.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", appError.Message)
	}

	if !errors.Is(appError.Cause, cause) {
		t.Errorf("Expected cause %v, got %v", cause, appError.Cause)
	}

	if appError.Context == nil {
		t.Error("Expected context to be initialized")
	}

	if appError.Stack == "" {
		t.Error("Expected stack trace to be captured")
	}
}

func TestAppError_WithContext(t *testing.T) {
	appError := NewAppError(ErrorTypeValidation, "test message", nil)
	result := appError.WithContext("key", "value")

	if result != appError {
		t.Error("WithContext should return the same instance")
	}

	if appError.Context["key"] != "value" {
		t.Errorf("Expected context key 'key' to have value 'value', got %v", appError.Context["key"])
	}
}

func TestAppError_ErrorTypeString(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrorTypeValidation, "Validation Error"},
		{ErrorTypeTemplate, "Template Error"},
		{ErrorTypeFileSystem, "File System Error"},
		{ErrorTypeNetwork, "Network Error"},
		{ErrorTypeConfiguration, "Configuration Error"},
		{ErrorTypeGeneration, "Generation Error"},
		{ErrorTypeInternal, "Internal Error"},
		{ErrorType(999), "Unknown Error"}, // Test unknown type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			appError := &AppError{Type: tt.errorType}
			if got := appError.ErrorTypeString(); got != tt.expected {
				t.Errorf("ErrorTypeString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewErrorHandler(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)
	if handler == nil {
		t.Fatal("NewErrorHandler returned nil")
	}

	if handler.logger != logger {
		t.Error("ErrorHandler logger not set correctly")
	}
}

func TestErrorHandler_Handle(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)

	// Test handling AppError
	appError := NewAppError(ErrorTypeValidation, "validation failed", nil)
	handler.Handle(appError) // Should not panic

	// Test handling generic error
	genericError := errors.New("generic error")
	handler.Handle(genericError) // Should not panic
}

func TestErrorHandler_handleAppError(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)

	tests := []struct {
		name      string
		errorType ErrorType
	}{
		{"validation error", ErrorTypeValidation},
		{"configuration error", ErrorTypeConfiguration},
		{"network error", ErrorTypeNetwork},
		{"template error", ErrorTypeTemplate},
		{"filesystem error", ErrorTypeFileSystem},
		{"generation error", ErrorTypeGeneration},
		{"internal error", ErrorTypeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := NewAppError(tt.errorType, "test error", nil)
			_ = appError.WithContext("test_key", "test_value")
			handler.handleAppError(appError) // Should not panic
		})
	}
}

func TestErrorHandler_formatStructuredLogMessage(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)
	appError := NewAppError(ErrorTypeValidation, "test message", nil)
	_ = appError.WithContext("key1", "value1")
	_ = appError.WithContext("key2", 42)

	message := handler.formatStructuredLogMessage(appError)

	if !strings.Contains(message, "type=Validation Error") {
		t.Error("Message should contain error type")
	}

	if !strings.Contains(message, "message=\"test message\"") {
		t.Error("Message should contain error message")
	}

	if !strings.Contains(message, "key1=value1") {
		t.Error("Message should contain context key1")
	}

	if !strings.Contains(message, "key2=42") {
		t.Error("Message should contain context key2")
	}
}

func TestErrorHandler_logErrorDetails(t *testing.T) {
	logger, err := NewLogger(LogLevelDebug, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)
	cause := errors.New("underlying cause")
	appError := NewAppError(ErrorTypeInternal, "test error", cause)

	// Should not panic
	handler.logErrorDetails(appError, LogLevelDebug)
	handler.logErrorDetails(appError, LogLevelInfo)
	handler.logErrorDetails(appError, LogLevelError)
}

func TestErrorHandler_handleGenericError(t *testing.T) {
	logger, err := NewLogger(LogLevelInfo, false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	handler := NewErrorHandler(logger)
	genericError := errors.New("generic error")

	// Should not panic
	handler.handleGenericError(genericError)
}

func TestWrapperFunctions(t *testing.T) {
	cause := errors.New("test cause")

	tests := []struct {
		name     string
		wrapper  func(string, error) *AppError
		expected ErrorType
	}{
		{"WrapValidationError", WrapValidationError, ErrorTypeValidation},
		{"WrapTemplateError", WrapTemplateError, ErrorTypeTemplate},
		{"WrapFileSystemError", WrapFileSystemError, ErrorTypeFileSystem},
		{"WrapNetworkError", WrapNetworkError, ErrorTypeNetwork},
		{"WrapConfigurationError", WrapConfigurationError, ErrorTypeConfiguration},
		{"WrapGenerationError", WrapGenerationError, ErrorTypeGeneration},
		{"WrapInternalError", WrapInternalError, ErrorTypeInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := tt.wrapper("test message", cause)
			if appError.Type != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, appError.Type)
			}
			if appError.Message != "test message" {
				t.Errorf("Expected message 'test message', got '%s'", appError.Message)
			}
			if !errors.Is(appError.Cause, cause) {
				t.Errorf("Expected cause %v, got %v", cause, appError.Cause)
			}
		})
	}
}

func TestContextualErrorCreation(t *testing.T) {
	tests := []struct {
		name     string
		creator  func() *AppError
		expected ErrorType
	}{
		{
			name: "NewValidationErrorWithContext",
			creator: func() *AppError {
				return NewValidationErrorWithContext("field1", "is required", "value1")
			},
			expected: ErrorTypeValidation,
		},
		{
			name: "NewTemplateErrorWithContext",
			creator: func() *AppError {
				return NewTemplateErrorWithContext("/path/template", "render", "failed", nil)
			},
			expected: ErrorTypeTemplate,
		},
		{
			name: "NewFileSystemErrorWithContext",
			creator: func() *AppError {
				return NewFileSystemErrorWithContext("/path/file", "write", "permission denied", nil)
			},
			expected: ErrorTypeFileSystem,
		},
		{
			name: "NewNetworkErrorWithContext",
			creator: func() *AppError {
				return NewNetworkErrorWithContext("http://example.com", "GET", "timeout", nil)
			},
			expected: ErrorTypeNetwork,
		},
		{
			name: "NewConfigurationErrorWithContext",
			creator: func() *AppError {
				return NewConfigurationErrorWithContext("database.host", "invalid value", nil)
			},
			expected: ErrorTypeConfiguration,
		},
		{
			name: "NewGenerationErrorWithContext",
			creator: func() *AppError {
				return NewGenerationErrorWithContext("frontend", "scaffold", "template not found", nil)
			},
			expected: ErrorTypeGeneration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appError := tt.creator()
			if appError.Type != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, appError.Type)
			}
			if len(appError.Context) == 0 {
				t.Error("Expected context to be populated")
			}
		})
	}
}

func TestPropagateError(t *testing.T) {
	// Test with nil error
	result := PropagateError(nil, "test context")
	if result != nil {
		t.Error("PropagateError should return nil for nil input")
	}

	// Test with AppError
	appError := NewAppError(ErrorTypeValidation, "test error", nil)
	result = PropagateError(appError, "test context")
	if !errors.Is(result, appError) {
		t.Error("PropagateError should return the same AppError instance")
	}
	if appError.Context["propagation_context"] != "test context" {
		t.Error("PropagateError should add propagation context")
	}

	// Test with generic error
	genericError := errors.New("generic error")
	result = PropagateError(genericError, "test context")
	if result == nil {
		t.Fatal("PropagateError should not return nil for generic error")
	}
	appErr := &AppError{}
	ok := errors.As(result, &appErr)
	if !ok {
		t.Fatal("PropagateError should return AppError for generic error")
	}
	if appErr.Type != ErrorTypeInternal {
		t.Error("PropagateError should wrap generic error as internal error")
	}
}

func TestChainErrors(t *testing.T) {
	// Test with empty slice
	result := ChainErrors([]error{}, "test operation")
	if result != nil {
		t.Error("ChainErrors should return nil for empty slice")
	}

	// Test with single error
	singleError := errors.New("single error")
	result = ChainErrors([]error{singleError}, "test operation")
	if result == nil {
		t.Fatal("ChainErrors should not return nil for single error")
	}

	// Test with multiple errors
	errors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		nil, // Should be skipped
		errors.New("error 3"),
	}
	result = ChainErrors(errors, "test operation")
	if result == nil {
		t.Fatal("ChainErrors should not return nil for multiple errors")
	}
	if result.Type != ErrorTypeInternal {
		t.Error("ChainErrors should create internal error")
	}
	if !strings.Contains(result.Message, "multiple errors during test operation") {
		t.Error("ChainErrors should include operation in message")
	}
	if result.Context["error_count"] != 4 {
		t.Errorf("Expected error_count to be 4, got %v", result.Context["error_count"])
	}
}
