// Package errors provides comprehensive error handling for the CLI application
package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// CLIError represents a comprehensive CLI error with categorization and context
type CLIError struct {
	Type        string                 `json:"type"`
	Message     string                 `json:"message"`
	Code        int                    `json:"code"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Context     *ErrorContext          `json:"context,omitempty"`
	Cause       error                  `json:"cause,omitempty"`
	Stack       string                 `json:"stack,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    Severity               `json:"severity"`
	Recoverable bool                   `json:"recoverable"`
}

// ErrorContext provides detailed context information for errors
type ErrorContext struct {
	Command     string            `json:"command"`
	Arguments   []string          `json:"arguments"`
	Flags       map[string]string `json:"flags"`
	WorkingDir  string            `json:"working_dir"`
	Environment string            `json:"environment,omitempty"`
	CI          *CIEnvironment    `json:"ci,omitempty"`
	Operation   string            `json:"operation,omitempty"`
	Component   string            `json:"component,omitempty"`
	File        string            `json:"file,omitempty"`
	Line        int               `json:"line,omitempty"`
}

// CIEnvironment contains CI-specific information
type CIEnvironment struct {
	IsCI     bool   `json:"is_ci"`
	Provider string `json:"provider,omitempty"`
	JobID    string `json:"job_id,omitempty"`
	BuildID  string `json:"build_id,omitempty"`
}

// Severity represents the severity level of an error
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Error type constants
const (
	ErrorTypeValidation    = "validation"
	ErrorTypeConfiguration = "configuration"
	ErrorTypeTemplate      = "template"
	ErrorTypeNetwork       = "network"
	ErrorTypeFileSystem    = "filesystem"
	ErrorTypePermission    = "permission"
	ErrorTypeCache         = "cache"
	ErrorTypeVersion       = "version"
	ErrorTypeAudit         = "audit"
	ErrorTypeGeneration    = "generation"
	ErrorTypeInternal      = "internal"
	ErrorTypeUser          = "user"
	ErrorTypeDependency    = "dependency"
	ErrorTypeSecurity      = "security"
)

// Exit code constants
const (
	ExitCodeSuccess              = 0
	ExitCodeGeneral              = 1
	ExitCodeValidationFailed     = 2
	ExitCodeConfigurationInvalid = 3
	ExitCodeTemplateNotFound     = 4
	ExitCodeNetworkError         = 5
	ExitCodeFileSystemError      = 6
	ExitCodePermissionDenied     = 7
	ExitCodeCacheError           = 8
	ExitCodeVersionError         = 9
	ExitCodeAuditFailed          = 10
	ExitCodeGenerationFailed     = 11
	ExitCodeDependencyError      = 12
	ExitCodeSecurityError        = 13
	ExitCodeUserError            = 14
	ExitCodeInternalError        = 99
)

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for error wrapping
func (e *CLIError) Unwrap() error {
	return e.Cause
}

// ExitCode returns the appropriate exit code for the error
func (e *CLIError) ExitCode() int {
	return e.Code
}

// ToJSON converts the error to JSON format
func (e *CLIError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// IsRecoverable returns whether the error is recoverable
func (e *CLIError) IsRecoverable() bool {
	return e.Recoverable
}

// GetSeverity returns the error severity
func (e *CLIError) GetSeverity() Severity {
	return e.Severity
}

// NewCLIError creates a new CLI error with comprehensive information
func NewCLIError(errorType, message string, code int) *CLIError {
	// Capture stack trace
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	return &CLIError{
		Type:        errorType,
		Message:     message,
		Code:        code,
		Details:     make(map[string]interface{}),
		Suggestions: []string{},
		Timestamp:   time.Now(),
		Stack:       stack,
		Severity:    determineSeverity(errorType, code),
		Recoverable: determineRecoverable(errorType, code),
	}
}

// WithDetails adds details to the CLI error
func (e *CLIError) WithDetails(key string, value interface{}) *CLIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithSuggestions adds actionable suggestions to the CLI error
func (e *CLIError) WithSuggestions(suggestions ...string) *CLIError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithContext adds context information to the CLI error
func (e *CLIError) WithContext(ctx *ErrorContext) *CLIError {
	e.Context = ctx
	return e
}

// WithCause adds the underlying cause to the CLI error
func (e *CLIError) WithCause(cause error) *CLIError {
	e.Cause = cause
	return e
}

// WithSeverity sets the error severity
func (e *CLIError) WithSeverity(severity Severity) *CLIError {
	e.Severity = severity
	return e
}

// WithRecoverable sets whether the error is recoverable
func (e *CLIError) WithRecoverable(recoverable bool) *CLIError {
	e.Recoverable = recoverable
	return e
}

// determineSeverity determines the severity based on error type and code
func determineSeverity(errorType string, code int) Severity {
	switch errorType {
	case ErrorTypeSecurity:
		return SeverityCritical
	case ErrorTypeInternal:
		return SeverityHigh
	case ErrorTypeFileSystem, ErrorTypePermission:
		return SeverityHigh
	case ErrorTypeNetwork, ErrorTypeCache:
		return SeverityMedium
	case ErrorTypeValidation, ErrorTypeConfiguration:
		return SeverityMedium
	case ErrorTypeTemplate, ErrorTypeVersion:
		return SeverityLow
	case ErrorTypeUser:
		return SeverityLow
	default:
		if code >= ExitCodeInternalError {
			return SeverityCritical
		} else if code >= ExitCodeAuditFailed {
			return SeverityHigh
		} else if code >= ExitCodeNetworkError {
			return SeverityMedium
		}
		return SeverityLow
	}
}

// determineRecoverable determines if an error is recoverable
func determineRecoverable(errorType string, code int) bool {
	switch errorType {
	case ErrorTypeNetwork, ErrorTypeCache:
		return true
	case ErrorTypeValidation, ErrorTypeConfiguration:
		return true
	case ErrorTypeTemplate:
		return true
	case ErrorTypeUser:
		return true
	case ErrorTypeVersion:
		return true
	case ErrorTypeFileSystem:
		return false // Usually requires manual intervention
	case ErrorTypePermission:
		return false // Requires permission changes
	case ErrorTypeSecurity:
		return false // Security issues need careful handling
	case ErrorTypeInternal:
		return false // Internal errors are not recoverable
	default:
		return code < ExitCodeFileSystemError
	}
}

// Error factory functions for common error types

// NewValidationError creates a validation error with standard suggestions
func NewValidationError(message string, field string, value interface{}) *CLIError {
	err := NewCLIError(ErrorTypeValidation, message, ExitCodeValidationFailed)
	err = err.WithDetails("field", field).WithDetails("value", value)
	err = err.WithSuggestions(
		"Check the project structure and configuration files",
		"Run with --verbose flag for more details",
		"Use --fix flag to automatically fix common issues",
		"Validate configuration with 'generator config validate'",
	)
	return err
}

// NewConfigurationError creates a configuration error with standard suggestions
func NewConfigurationError(message string, configPath string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeConfiguration, message, ExitCodeConfigurationInvalid)
	err = err.WithDetails("config_path", configPath).WithCause(cause)
	err = err.WithSuggestions(
		"Check the configuration file syntax",
		"Validate the configuration with 'generator config validate'",
		"Use 'generator config show' to see current configuration",
		"Check file permissions and accessibility",
	)
	return err
}

// NewTemplateError creates a template error with standard suggestions
func NewTemplateError(message string, templateName string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeTemplate, message, ExitCodeTemplateNotFound)
	err = err.WithDetails("template_name", templateName).WithCause(cause)
	err = err.WithSuggestions(
		"List available templates with 'generator list-templates'",
		"Check template name spelling",
		"Use 'generator template info <name>' for template details",
		"Verify template exists and is accessible",
	)
	return err
}

// NewNetworkError creates a network error with standard suggestions
func NewNetworkError(message string, url string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeNetwork, message, ExitCodeNetworkError)
	err = err.WithDetails("url", url).WithCause(cause)
	err = err.WithSuggestions(
		"Check your internet connection",
		"Use --offline flag to work with cached data",
		"Check if the URL is accessible",
		"Verify proxy settings if applicable",
		"Try again later if the service is temporarily unavailable",
	)
	return err
}

// NewFileSystemError creates a filesystem error with standard suggestions
func NewFileSystemError(message string, path string, operation string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeFileSystem, message, ExitCodeFileSystemError)
	err = err.WithDetails("path", path).WithDetails("operation", operation).WithCause(cause)
	err = err.WithSuggestions(
		"Check if the path exists and is accessible",
		"Verify file permissions",
		"Ensure sufficient disk space",
		"Check if the file is locked by another process",
	)
	return err
}

// NewPermissionError creates a permission error with standard suggestions
func NewPermissionError(message string, path string, requiredPermission string, cause error) *CLIError {
	err := NewCLIError(ErrorTypePermission, message, ExitCodePermissionDenied)
	err = err.WithDetails("path", path).WithDetails("required_permission", requiredPermission).WithCause(cause)
	err = err.WithSuggestions(
		"Check file/directory permissions",
		"Run with appropriate user privileges",
		"Verify ownership of the target directory",
		fmt.Sprintf("Ensure %s permission is granted", requiredPermission),
	)
	return err
}

// NewCacheError creates a cache error with standard suggestions
func NewCacheError(message string, cacheKey string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeCache, message, ExitCodeCacheError)
	err = err.WithDetails("cache_key", cacheKey).WithCause(cause)
	err = err.WithSuggestions(
		"Clear cache with 'generator cache clear'",
		"Check cache directory permissions",
		"Verify sufficient disk space for cache",
		"Try running without cache (--no-cache flag)",
	)
	return err
}

// NewVersionError creates a version error with standard suggestions
func NewVersionError(message string, component string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeVersion, message, ExitCodeVersionError)
	err = err.WithDetails("component", component).WithCause(cause)
	err = err.WithSuggestions(
		"Check for updates with 'generator version --check-updates'",
		"Verify internet connection for version checking",
		"Use --offline flag to skip version checks",
		"Check component compatibility requirements",
	)
	return err
}

// NewAuditError creates an audit error with standard suggestions
func NewAuditError(message string, auditType string, score float64, cause error) *CLIError {
	err := NewCLIError(ErrorTypeAudit, message, ExitCodeAuditFailed)
	err = err.WithDetails("audit_type", auditType).WithDetails("score", score).WithCause(cause)
	err = err.WithSuggestions(
		"Review audit recommendations",
		"Fix high-priority security issues",
		"Improve code quality metrics",
		"Check project dependencies for vulnerabilities",
	)
	return err
}

// NewGenerationError creates a generation error with standard suggestions
func NewGenerationError(message string, component string, operation string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeGeneration, message, ExitCodeGenerationFailed)
	err = err.WithDetails("component", component).WithDetails("operation", operation).WithCause(cause)
	err = err.WithSuggestions(
		"Check template and configuration validity",
		"Verify output directory permissions",
		"Use --dry-run to preview generation",
		"Check for conflicting files in output directory",
	)
	return err
}

// NewDependencyError creates a dependency error with standard suggestions
func NewDependencyError(message string, dependency string, version string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeDependency, message, ExitCodeDependencyError)
	err = err.WithDetails("dependency", dependency).WithDetails("version", version).WithCause(cause)
	err = err.WithSuggestions(
		"Check dependency version compatibility",
		"Update dependencies to compatible versions",
		"Review dependency documentation",
		"Use --update-versions flag to get latest versions",
	)
	return err
}

// NewSecurityError creates a security error with standard suggestions
func NewSecurityError(message string, securityIssue string, severity string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeSecurity, message, ExitCodeSecurityError)
	err = err.WithDetails("security_issue", securityIssue).WithDetails("severity", severity).WithCause(cause)
	err = err.WithSuggestions(
		"Address security vulnerabilities immediately",
		"Review security audit recommendations",
		"Update vulnerable dependencies",
		"Implement recommended security practices",
	)
	return err
}

// NewUserError creates a user error with standard suggestions
func NewUserError(message string, userInput string, expectedFormat string) *CLIError {
	err := NewCLIError(ErrorTypeUser, message, ExitCodeUserError)
	err = err.WithDetails("user_input", userInput).WithDetails("expected_format", expectedFormat)
	err = err.WithSuggestions(
		"Check command syntax and arguments",
		"Use --help flag for command usage information",
		"Verify input format and values",
		"Check examples in documentation",
	)
	return err
}

// NewInternalError creates an internal error with standard suggestions
func NewInternalError(message string, component string, cause error) *CLIError {
	err := NewCLIError(ErrorTypeInternal, message, ExitCodeInternalError)
	err = err.WithDetails("component", component).WithCause(cause)
	err = err.WithSuggestions(
		"This is an internal error - please report it",
		"Include error details and steps to reproduce",
		"Try running with --debug flag for more information",
		"Check if the issue persists with latest version",
	)
	return err
}

// Error chaining and propagation utilities

// WrapError wraps an existing error with additional context
func WrapError(err error, errorType string, message string, code int) *CLIError {
	if err == nil {
		return nil
	}

	// If it's already a CLIError, enhance it
	cliErr := &CLIError{}
	if stderrors.As(err, &cliErr) {
		cliErr.Message = fmt.Sprintf("%s: %s", message, cliErr.Message)
		return cliErr
	}

	// Create new CLIError wrapping the original
	return NewCLIError(errorType, message, code).WithCause(err)
}

// ChainErrors combines multiple errors into a single error
func ChainErrors(errors []error, operation string) *CLIError {
	if len(errors) == 0 {
		return nil
	}

	// Filter out nil errors
	var validErrors []error
	for _, err := range errors {
		if err != nil {
			validErrors = append(validErrors, err)
		}
	}

	if len(validErrors) == 0 {
		return nil
	}

	if len(validErrors) == 1 {
		return WrapError(validErrors[0], ErrorTypeInternal, fmt.Sprintf("Error during %s", operation), ExitCodeGeneral)
	}

	// Create summary error for multiple errors
	var messages []string
	for i, err := range validErrors {
		messages = append(messages, fmt.Sprintf("error %d: %s", i+1, err.Error()))
	}

	chainedMessage := fmt.Sprintf("Multiple errors during %s", operation)
	chainedErr := NewCLIError(ErrorTypeInternal, chainedMessage, ExitCodeGeneral)
	chainedErr = chainedErr.WithDetails("error_count", len(validErrors))
	chainedErr = chainedErr.WithDetails("error_summary", strings.Join(messages, "; "))

	return chainedErr
}

// PropagateError propagates an error with additional context
func PropagateError(err error, context string) error {
	if err == nil {
		return nil
	}

	// If it's already a CLIError, add context
	cliErr := &CLIError{}
	if stderrors.As(err, &cliErr) {
		if cliErr.Context == nil {
			cliErr.Context = &ErrorContext{}
		}
		if cliErr.Context.Operation == "" {
			cliErr.Context.Operation = context
		} else {
			cliErr.Context.Operation = fmt.Sprintf("%s -> %s", cliErr.Context.Operation, context)
		}
		return cliErr
	}

	// Wrap as internal error
	return NewInternalError(fmt.Sprintf("%s: %s", context, err.Error()), context, err)
}
