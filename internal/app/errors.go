package app

import (
	"fmt"
	"runtime"
	"strings"
)

// AppError represents an application-specific error
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
	Stack   string
}

// ErrorType represents the type of error
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota
	ErrorTypeTemplate
	ErrorTypeFileSystem
	ErrorTypeNetwork
	ErrorTypeConfiguration
	ErrorTypeGeneration
	ErrorTypeInternal
)

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error
func NewAppError(errorType ErrorType, message string, cause error) *AppError {
	// Capture stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
		Stack:   stack,
	}
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// ErrorTypeString returns a string representation of the error type
func (e *AppError) ErrorTypeString() string {
	switch e.Type {
	case ErrorTypeValidation:
		return "Validation Error"
	case ErrorTypeTemplate:
		return "Template Error"
	case ErrorTypeFileSystem:
		return "File System Error"
	case ErrorTypeNetwork:
		return "Network Error"
	case ErrorTypeConfiguration:
		return "Configuration Error"
	case ErrorTypeGeneration:
		return "Generation Error"
	case ErrorTypeInternal:
		return "Internal Error"
	default:
		return "Unknown Error"
	}
}

// ErrorHandler provides centralized error handling
type ErrorHandler struct {
	logger *Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Handle handles an error appropriately based on its type
func (h *ErrorHandler) Handle(err error) {
	if appErr, ok := err.(*AppError); ok {
		h.handleAppError(appErr)
	} else {
		h.handleGenericError(err)
	}
}

// handleAppError handles application-specific errors with standardized logging
func (h *ErrorHandler) handleAppError(err *AppError) {
	// Create structured log message
	logMsg := h.formatStructuredLogMessage(err)

	// Log based on error type and severity
	switch err.Type {
	case ErrorTypeValidation, ErrorTypeConfiguration:
		h.logger.Warn("%s", logMsg)
		h.logErrorDetails(err, LogLevelDebug)
	case ErrorTypeNetwork:
		h.logger.Warn("%s", logMsg)
		h.logErrorDetails(err, LogLevelDebug)
	case ErrorTypeTemplate, ErrorTypeFileSystem, ErrorTypeGeneration:
		h.logger.Error("%s", logMsg)
		h.logErrorDetails(err, LogLevelInfo)
	case ErrorTypeInternal:
		h.logger.Error("%s", logMsg)
		h.logErrorDetails(err, LogLevelError)
		// Always log stack trace for internal errors
		h.logger.Error("Stack trace: %s", err.Stack)
	default:
		h.logger.Error("%s", logMsg)
		h.logErrorDetails(err, LogLevelInfo)
	}
}

// formatStructuredLogMessage creates a consistent log message format
func (h *ErrorHandler) formatStructuredLogMessage(err *AppError) string {
	var parts []string

	// Add error type
	parts = append(parts, fmt.Sprintf("type=%s", err.ErrorTypeString()))

	// Add message
	parts = append(parts, fmt.Sprintf("message=\"%s\"", err.Message))

	// Add context information
	if len(err.Context) > 0 {
		for k, v := range err.Context {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
	}

	return strings.Join(parts, " ")
}

// logErrorDetails logs additional error details based on log level
func (h *ErrorHandler) logErrorDetails(err *AppError, minLevel LogLevel) {
	// Get current logger level
	currentLevel := h.logger.GetLevel()
	if currentLevel > int(minLevel) {
		return
	}

	// Log cause if available
	if err.Cause != nil {
		h.logger.Debug("caused_by=\"%v\"", err.Cause)
	}

	// Log stack trace for debug level
	if currentLevel <= int(LogLevelDebug) && err.Stack != "" {
		h.logger.Debug("stack_trace=\"%s\"", err.Stack)
	}
}

// handleGenericError handles generic errors
func (h *ErrorHandler) handleGenericError(err error) {
	h.logger.Error("Unexpected error: %v", err)
}

// WrapValidationError wraps a validation error
func WrapValidationError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeValidation, message, cause)
}

// WrapTemplateError wraps a template error
func WrapTemplateError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeTemplate, message, cause)
}

// WrapFileSystemError wraps a file system error
func WrapFileSystemError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeFileSystem, message, cause)
}

// WrapNetworkError wraps a network error
func WrapNetworkError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeNetwork, message, cause)
}

// WrapConfigurationError wraps a configuration error
func WrapConfigurationError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeConfiguration, message, cause)
}

// WrapGenerationError wraps a generation error
func WrapGenerationError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeGeneration, message, cause)
}

// WrapInternalError wraps an internal error
func WrapInternalError(message string, cause error) *AppError {
	return NewAppError(ErrorTypeInternal, message, cause)
}

// StandardizedErrorCreation provides consistent error creation patterns

// NewValidationErrorWithContext creates a validation error with standardized context
func NewValidationErrorWithContext(field, message string, value interface{}) *AppError {
	err := NewAppError(ErrorTypeValidation,
		fmt.Sprintf("validation failed for field '%s': %s", field, message), nil)
	return err.WithContext("field", field).WithContext("value", value)
}

// NewTemplateErrorWithContext creates a template error with standardized context
func NewTemplateErrorWithContext(templatePath, operation, message string, cause error) *AppError {
	err := NewAppError(ErrorTypeTemplate,
		fmt.Sprintf("template processing failed: %s", message), cause)
	return err.WithContext("template_path", templatePath).WithContext("operation", operation)
}

// NewFileSystemErrorWithContext creates a filesystem error with standardized context
func NewFileSystemErrorWithContext(path, operation, message string, cause error) *AppError {
	err := NewAppError(ErrorTypeFileSystem,
		fmt.Sprintf("filesystem operation failed: %s", message), cause)
	return err.WithContext("path", path).WithContext("operation", operation)
}

// NewNetworkErrorWithContext creates a network error with standardized context
func NewNetworkErrorWithContext(endpoint, operation, message string, cause error) *AppError {
	err := NewAppError(ErrorTypeNetwork,
		fmt.Sprintf("network operation failed: %s", message), cause)
	return err.WithContext("endpoint", endpoint).WithContext("operation", operation)
}

// NewConfigurationErrorWithContext creates a configuration error with standardized context
func NewConfigurationErrorWithContext(configKey, message string, cause error) *AppError {
	err := NewAppError(ErrorTypeConfiguration,
		fmt.Sprintf("configuration error: %s", message), cause)
	return err.WithContext("config_key", configKey)
}

// NewGenerationErrorWithContext creates a generation error with standardized context
func NewGenerationErrorWithContext(component, operation, message string, cause error) *AppError {
	err := NewAppError(ErrorTypeGeneration,
		fmt.Sprintf("generation failed: %s", message), cause)
	return err.WithContext("component", component).WithContext("operation", operation)
}

// ErrorPropagation provides standardized error propagation patterns

// PropagateError propagates an error with additional context while preserving the original error type
func PropagateError(err error, context string) error {
	if err == nil {
		return nil
	}

	// If it's already an AppError, add context and return
	if appErr, ok := err.(*AppError); ok {
		return appErr.WithContext("propagation_context", context)
	}

	// Otherwise, wrap as internal error
	return NewAppError(ErrorTypeInternal, fmt.Sprintf("%s: %s", context, err.Error()), err)
}

// ChainErrors chains multiple errors together with context
func ChainErrors(errors []error, operation string) *AppError {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		return PropagateError(errors[0], operation).(*AppError)
	}

	// Create a summary error for multiple errors
	var messages []string
	for i, err := range errors {
		if err != nil {
			messages = append(messages, fmt.Sprintf("error %d: %s", i+1, err.Error()))
		}
	}

	chainedMessage := fmt.Sprintf("multiple errors during %s: %s", operation, strings.Join(messages, "; "))
	chainedErr := NewAppError(ErrorTypeInternal, chainedMessage, nil)
	_ = chainedErr.WithContext("error_count", len(errors))

	return chainedErr
}
