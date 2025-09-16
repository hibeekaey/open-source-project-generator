package models

import (
	"errors"
	"fmt"
)

// ErrorType represents the type of error
type ErrorType int

const (
	ValidationErrorType ErrorType = iota
	TemplateErrorType
	FileSystemErrorType
	ConfigurationErrorType
)

// GeneratorError represents a custom error type for the generator
type GeneratorError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
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
	}
}

// Predefined error types
var (
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrTemplateNotFound = errors.New("template not found")
	ErrFileSystemError  = errors.New("file system error")
	ErrValidationFailed = errors.New("validation failed")
)
