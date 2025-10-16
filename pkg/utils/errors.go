package utils

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ErrorContext provides additional context for error handling
type ErrorContext struct {
	Operation string
	Component string
	File      string
	Line      int
	UserID    string
	RequestID string
	Metadata  map[string]interface{}
}

// NewErrorContext creates a new error context with caller information
func NewErrorContext(operation, component string) *ErrorContext {
	_, file, line, _ := runtime.Caller(1)
	return &ErrorContext{
		Operation: operation,
		Component: component,
		File:      file,
		Line:      line,
		Metadata:  make(map[string]interface{}),
	}
}

// WithMetadata adds metadata to the error context
func (ec *ErrorContext) WithMetadata(key string, value interface{}) *ErrorContext {
	ec.Metadata[key] = value
	return ec
}

// WithUserID adds user ID to the error context
func (ec *ErrorContext) WithUserID(userID string) *ErrorContext {
	ec.UserID = userID
	return ec
}

// WithRequestID adds request ID to the error context
func (ec *ErrorContext) WithRequestID(requestID string) *ErrorContext {
	ec.RequestID = requestID
	return ec
}

// HandleError provides standard error handling with context
func HandleError(err error, context string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}

// WrapError wraps an error with additional context using standardized format
func WrapError(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}

// WrapErrorWithContext wraps an error with structured context information
func WrapErrorWithContext(err error, ctx *ErrorContext, message string) error {
	if err == nil {
		return nil
	}

	// Create a standardized error message format
	var parts []string

	if ctx.Component != "" {
		parts = append(parts, fmt.Sprintf("component=%s", ctx.Component))
	}

	if ctx.Operation != "" {
		parts = append(parts, fmt.Sprintf("operation=%s", ctx.Operation))
	}

	if message != "" {
		parts = append(parts, message)
	}

	contextStr := strings.Join(parts, " ")
	return fmt.Errorf("%s: %w", contextStr, err)
}

// IsNilError checks if an error is nil and returns a formatted error if not
func IsNilError(err error, operation string) error {
	if err != nil {
		return fmt.Errorf("failed to %s: %w", operation, err)
	}
	return nil
}

// NewValidationError creates a standardized validation error
func NewValidationError(field, message string, value interface{}) *models.GeneratorError {
	return models.NewGeneratorError(
		models.ValidationErrorType,
		fmt.Sprintf("validation failed for field '%s': %s", field, message),
		nil,
	)
}

// NewTemplateError creates a standardized template processing error
func NewTemplateError(templatePath, operation, message string, cause error) *models.GeneratorError {
	return models.NewGeneratorError(
		models.TemplateErrorType,
		fmt.Sprintf("template %s failed during %s: %s", templatePath, operation, message),
		cause,
	)
}

// NewFileSystemError creates a standardized filesystem error
func NewFileSystemError(path, operation, message string, cause error) *models.GeneratorError {
	return models.NewGeneratorError(
		models.FileSystemErrorType,
		fmt.Sprintf("filesystem %s failed for path '%s': %s", operation, path, message),
		cause,
	)
}

// NewConfigurationError creates a standardized configuration error
func NewConfigurationError(configKey, message string, cause error) *models.GeneratorError {
	return models.NewGeneratorError(
		models.ConfigurationErrorType,
		fmt.Sprintf("configuration error for '%s': %s", configKey, message),
		cause,
	)
}

// ValidateAndWrapError validates an error and wraps it with appropriate context
func ValidateAndWrapError(err error, ctx *ErrorContext) error {
	if err == nil {
		return nil
	}

	// Check if it's already a GeneratorError
	genErr := &models.GeneratorError{}
	if errors.As(err, &genErr) {
		// Return the error as-is since WithContext method was removed
		return genErr
	}

	// Wrap as a generic error with context
	return WrapErrorWithContext(err, ctx, "unexpected error occurred")
}

// FormatErrorForUser formats an error message for user display
func FormatErrorForUser(err error) string {
	if err == nil {
		return ""
	}

	// Handle GeneratorError
	genErr := &models.GeneratorError{}
	if errors.As(err, &genErr) {
		return formatGeneratorErrorForUser(genErr)
	}

	// Note: SecurityOperationError was removed in simplified architecture

	// Default formatting for other errors
	return fmt.Sprintf("Error: %s", err.Error())
}

// formatGeneratorErrorForUser formats a GeneratorError for user display
func formatGeneratorErrorForUser(err *models.GeneratorError) string {
	var message strings.Builder

	switch err.Type {
	case models.ValidationErrorType:
		message.WriteString("❌ Validation Error: ")
	case models.TemplateErrorType:
		message.WriteString("📄 Template Error: ")
	case models.FileSystemErrorType:
		message.WriteString("📁 File System Error: ")
	case models.ConfigurationErrorType:
		message.WriteString("⚙️ Configuration Error: ")
	default:
		message.WriteString("❗ Error: ")
	}

	message.WriteString(err.Message)

	// Context field removed

	return message.String()
}

// Note: formatSecurityErrorForUser was removed in simplified architecture
