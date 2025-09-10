package app

import (
	"fmt"
	"runtime"
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

// handleAppError handles application-specific errors
func (h *ErrorHandler) handleAppError(err *AppError) {
	// Log the error with context
	logMsg := fmt.Sprintf("[%s] %s", err.ErrorTypeString(), err.Message)

	if len(err.Context) > 0 {
		logMsg += " Context:"
		for k, v := range err.Context {
			logMsg += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	switch err.Type {
	case ErrorTypeValidation, ErrorTypeConfiguration:
		h.logger.Warn(logMsg)
		if h.logger.level <= LogLevelDebug && err.Cause != nil {
			h.logger.Debug("Caused by: %v", err.Cause)
		}
	case ErrorTypeNetwork:
		h.logger.Warn(logMsg)
		if h.logger.level <= LogLevelDebug {
			h.logger.Debug("Stack trace: %s", err.Stack)
		}
	case ErrorTypeTemplate, ErrorTypeFileSystem, ErrorTypeGeneration:
		h.logger.Error(logMsg)
		if err.Cause != nil {
			h.logger.Error("Caused by: %v", err.Cause)
		}
		if h.logger.level <= LogLevelDebug {
			h.logger.Debug("Stack trace: %s", err.Stack)
		}
	case ErrorTypeInternal:
		h.logger.Error(logMsg)
		h.logger.Error("Stack trace: %s", err.Stack)
		if err.Cause != nil {
			h.logger.Error("Caused by: %v", err.Cause)
		}
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
