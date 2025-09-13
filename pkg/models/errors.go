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
	NetworkErrorType
	ConfigurationErrorType
	SecurityErrorType
	CryptographicErrorType
	FileSecurityErrorType
	PathValidationErrorType
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

// SecurityOperationError represents security-specific errors with enhanced context
type SecurityOperationError struct {
	*GeneratorError
	Severity    SecuritySeverity `json:"severity"`
	Component   string           `json:"component"`
	Operation   string           `json:"operation"`
	Remediation string           `json:"remediation,omitempty"`
}

// SecuritySeverity represents the severity level of security errors
type SecuritySeverity int

const (
	SecuritySeverityLow SecuritySeverity = iota
	SecuritySeverityMedium
	SecuritySeverityHigh
	SecuritySeverityCritical
)

// String returns the string representation of SecuritySeverity
func (s SecuritySeverity) String() string {
	switch s {
	case SecuritySeverityLow:
		return "low"
	case SecuritySeverityMedium:
		return "medium"
	case SecuritySeverityHigh:
		return "high"
	case SecuritySeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// NewSecurityError creates a new SecurityOperationError with sanitized messages
func NewSecurityError(errorType ErrorType, severity SecuritySeverity, component, operation, message string, cause error) *SecurityOperationError {
	// Sanitize the message to prevent information leakage
	sanitizedMessage := sanitizeSecurityMessage(message)

	baseError := NewGeneratorError(errorType, sanitizedMessage, cause)

	return &SecurityOperationError{
		GeneratorError: baseError,
		Severity:       severity,
		Component:      component,
		Operation:      operation,
	}
}

// WithRemediation adds remediation guidance to the security error
func (e *SecurityOperationError) WithRemediation(remediation string) *SecurityOperationError {
	e.Remediation = remediation
	return e
}

// IsCritical returns true if the error is critical severity
func (e *SecurityOperationError) IsCritical() bool {
	return e.Severity == SecuritySeverityCritical
}

// WithContext adds context information to the security error
func (e *SecurityOperationError) WithContext(key string, value interface{}) *SecurityOperationError {
	_ = e.GeneratorError.WithContext(key, value)
	return e
}

// sanitizeSecurityMessage removes potentially sensitive information from error messages
func sanitizeSecurityMessage(message string) string {
	// Remove file paths that might contain sensitive information
	// This is a basic implementation - in production, you might want more sophisticated sanitization
	if len(message) > 200 {
		return message[:200] + "... [message truncated for security]"
	}
	return message
}

// Predefined security errors with proper context
var (
	// Cryptographic errors
	ErrInsufficientEntropy = NewSecurityError(
		CryptographicErrorType,
		SecuritySeverityCritical,
		"crypto",
		"random_generation",
		"insufficient entropy for secure random generation",
		nil,
	).WithRemediation("Ensure crypto/rand is available and system has sufficient entropy")

	ErrCryptographicFailure = NewSecurityError(
		CryptographicErrorType,
		SecuritySeverityHigh,
		"crypto",
		"operation",
		"cryptographic operation failed",
		nil,
	).WithRemediation("Check system cryptographic capabilities and retry operation")

	// File security errors
	ErrInvalidPath = NewSecurityError(
		PathValidationErrorType,
		SecuritySeverityHigh,
		"filesystem",
		"path_validation",
		"path validation failed: potential security risk",
		nil,
	).WithRemediation("Use only validated, sanitized file paths within allowed directories")

	ErrTempFileCreation = NewSecurityError(
		FileSecurityErrorType,
		SecuritySeverityMedium,
		"filesystem",
		"temp_file_creation",
		"failed to create secure temporary file",
		nil,
	).WithRemediation("Ensure temporary directory is writable and has sufficient space")

	ErrAtomicWrite = NewSecurityError(
		FileSecurityErrorType,
		SecuritySeverityMedium,
		"filesystem",
		"atomic_write",
		"atomic write operation failed",
		nil,
	).WithRemediation("Check file permissions and available disk space")

	ErrInsecurePermissions = NewSecurityError(
		FileSecurityErrorType,
		SecuritySeverityHigh,
		"filesystem",
		"permission_check",
		"file permissions are not secure",
		nil,
	).WithRemediation("Set appropriate file permissions (e.g., 0600 for sensitive files)")

	ErrDangerousDirectory = NewSecurityError(
		PathValidationErrorType,
		SecuritySeverityHigh,
		"filesystem",
		"directory_validation",
		"directory path is potentially dangerous",
		nil,
	).WithRemediation("Use only trusted directory paths and avoid user-controlled input")

	// General security errors
	ErrSecurityViolation = NewSecurityError(
		SecurityErrorType,
		SecuritySeverityHigh,
		"security",
		"validation",
		"security policy violation detected",
		nil,
	).WithRemediation("Review operation against security policies and requirements")

	ErrUnauthorizedOperation = NewSecurityError(
		SecurityErrorType,
		SecuritySeverityHigh,
		"security",
		"authorization",
		"unauthorized operation attempted",
		nil,
	).WithRemediation("Ensure proper authorization before performing sensitive operations")
)

// IsSecurityError checks if an error is a security-related error
func IsSecurityError(err error) bool {
	if err == nil {
		return false
	}

	var secErr *SecurityOperationError
	return errors.As(err, &secErr)
}

// GetSecuritySeverity extracts the security severity from an error
func GetSecuritySeverity(err error) SecuritySeverity {
	var secErr *SecurityOperationError
	if errors.As(err, &secErr) {
		return secErr.Severity
	}
	return SecuritySeverityLow
}

// WrapWithSecurityContext wraps an existing error with security context
func WrapWithSecurityContext(err error, severity SecuritySeverity, component, operation string) *SecurityOperationError {
	if err == nil {
		return nil
	}

	return NewSecurityError(
		SecurityErrorType,
		severity,
		component,
		operation,
		err.Error(),
		err,
	)
}
