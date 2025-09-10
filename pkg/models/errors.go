package models

import "fmt"

// ErrorType represents the type of error
type ErrorType int

const (
	ValidationErrorType ErrorType = iota
	TemplateErrorType
	FileSystemErrorType
	NetworkErrorType
	ConfigurationErrorType
)

// GeneratorError represents a custom error type for the generator
type GeneratorError struct {
	Type    ErrorType              `json:"type"`
	Message string                 `json:"message"`
	Cause   error                  `json:"cause,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (e *GeneratorError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *GeneratorError) Unwrap() error {
	return e.Cause
}

// NewGeneratorError creates a new GeneratorError
func NewGeneratorError(errorType ErrorType, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (e *GeneratorError) WithContext(key string, value interface{}) *GeneratorError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}
