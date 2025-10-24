package orchestrator

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	// ErrCategoryToolNotFound indicates a required tool is not available
	ErrCategoryToolNotFound ErrorCategory = "TOOL_NOT_FOUND"

	// ErrCategoryToolExecution indicates a tool execution failed
	ErrCategoryToolExecution ErrorCategory = "TOOL_EXECUTION"

	// ErrCategoryInvalidConfig indicates configuration validation failed
	ErrCategoryInvalidConfig ErrorCategory = "INVALID_CONFIG"

	// ErrCategoryFileSystem indicates a file system operation failed
	ErrCategoryFileSystem ErrorCategory = "FILE_SYSTEM"

	// ErrCategorySecurity indicates a security validation failed
	ErrCategorySecurity ErrorCategory = "SECURITY"

	// ErrCategoryIntegration indicates component integration failed
	ErrCategoryIntegration ErrorCategory = "INTEGRATION"

	// ErrCategoryValidation indicates structure validation failed
	ErrCategoryValidation ErrorCategory = "VALIDATION"

	// ErrCategoryTimeout indicates an operation timed out
	ErrCategoryTimeout ErrorCategory = "TIMEOUT"

	// ErrCategoryUnknown indicates an unknown error
	ErrCategoryUnknown ErrorCategory = "UNKNOWN"
)

// GenerationError represents an error that occurred during project generation
type GenerationError struct {
	Category    ErrorCategory
	Message     string
	Cause       error
	Component   string
	Recoverable bool
	Suggestions []string
}

// Error implements the error interface
func (e *GenerationError) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("[%s]", e.Category))

	if e.Component != "" {
		builder.WriteString(fmt.Sprintf(" Component '%s':", e.Component))
	}

	builder.WriteString(fmt.Sprintf(" %s", e.Message))

	if e.Cause != nil {
		builder.WriteString(fmt.Sprintf(" (caused by: %v)", e.Cause))
	}

	return builder.String()
}

// Unwrap returns the underlying cause error
func (e *GenerationError) Unwrap() error {
	return e.Cause
}

// WithSuggestions adds suggestions to the error
func (e *GenerationError) WithSuggestions(suggestions ...string) *GenerationError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// GetSuggestions returns formatted suggestions for resolving the error
func (e *GenerationError) GetSuggestions() string {
	if len(e.Suggestions) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\nSuggestions:\n")
	for i, suggestion := range e.Suggestions {
		builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
	}

	return builder.String()
}

// NewGenerationError creates a new generation error
func NewGenerationError(category ErrorCategory, message string, cause error) *GenerationError {
	return &GenerationError{
		Category:    category,
		Message:     message,
		Cause:       cause,
		Recoverable: isRecoverable(category),
		Suggestions: make([]string, 0),
	}
}

// NewToolNotFoundError creates a tool not found error
func NewToolNotFoundError(toolName string, component string) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryToolNotFound,
		Message:     fmt.Sprintf("required tool '%s' is not available", toolName),
		Component:   component,
		Recoverable: true,
		Suggestions: []string{
			fmt.Sprintf("Install %s on your system", toolName),
			"Use --no-external-tools flag to force fallback generation",
			"Check tool installation instructions in the documentation",
		},
	}
}

// NewToolExecutionError creates a tool execution error
func NewToolExecutionError(toolName string, component string, cause error) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryToolExecution,
		Message:     fmt.Sprintf("tool '%s' execution failed", toolName),
		Cause:       cause,
		Component:   component,
		Recoverable: true,
		Suggestions: []string{
			"Check tool output for specific error messages",
			"Verify tool is properly installed and configured",
			"Try running the tool manually to diagnose the issue",
			"Use --verbose flag for detailed execution logs",
		},
	}
}

// NewConfigValidationError creates a configuration validation error
func NewConfigValidationError(field string, message string) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryInvalidConfig,
		Message:     fmt.Sprintf("configuration validation failed for field '%s': %s", field, message),
		Recoverable: false,
		Suggestions: []string{
			"Check your configuration file for errors",
			"Refer to the configuration schema documentation",
			"Use --init-config to generate a valid configuration template",
		},
	}
}

// NewFileSystemError creates a file system error
func NewFileSystemError(operation string, path string, cause error) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryFileSystem,
		Message:     fmt.Sprintf("file system operation '%s' failed for path '%s'", operation, path),
		Cause:       cause,
		Recoverable: true,
		Suggestions: []string{
			"Check file system permissions",
			"Verify disk space is available",
			"Ensure the path is accessible",
		},
	}
}

// NewSecurityError creates a security validation error
func NewSecurityError(message string, cause error) *GenerationError {
	return &GenerationError{
		Category:    ErrCategorySecurity,
		Message:     message,
		Cause:       cause,
		Recoverable: false,
		Suggestions: []string{
			"Review input for potentially dangerous patterns",
			"Ensure paths do not contain traversal attempts",
			"Check that all inputs are properly sanitized",
		},
	}
}

// NewIntegrationError creates an integration error
func NewIntegrationError(message string, cause error) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryIntegration,
		Message:     message,
		Cause:       cause,
		Recoverable: true,
		Suggestions: []string{
			"Check that all components were generated successfully",
			"Verify component configurations are compatible",
			"Review integration logs for specific errors",
		},
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryValidation,
		Message:     message,
		Cause:       cause,
		Recoverable: true,
		Suggestions: []string{
			"Check generated project structure",
			"Verify all required files were created",
			"Review validation logs for specific issues",
		},
	}
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(operation string, component string) *GenerationError {
	return &GenerationError{
		Category:    ErrCategoryTimeout,
		Message:     fmt.Sprintf("operation '%s' timed out", operation),
		Component:   component,
		Recoverable: true,
		Suggestions: []string{
			"Increase timeout duration in configuration",
			"Check network connectivity if downloading dependencies",
			"Verify system resources are not constrained",
		},
	}
}

// isRecoverable determines if an error category is recoverable
func isRecoverable(category ErrorCategory) bool {
	switch category {
	case ErrCategoryToolNotFound:
		return true // Can fallback to custom generation
	case ErrCategoryToolExecution:
		return true // Can retry or fallback
	case ErrCategoryInvalidConfig:
		return false // Must fix configuration
	case ErrCategoryFileSystem:
		return true // Can retry after fixing permissions/space
	case ErrCategorySecurity:
		return false // Security issues must be resolved
	case ErrCategoryIntegration:
		return true // Can continue without full integration
	case ErrCategoryValidation:
		return true // Can continue with warnings
	case ErrCategoryTimeout:
		return true // Can retry with longer timeout
	default:
		return false
	}
}

// ErrorRecoveryStrategy defines how to recover from an error
type ErrorRecoveryStrategy struct {
	Category    ErrorCategory
	Description string
	Actions     []RecoveryAction
}

// RecoveryAction defines a specific recovery action
type RecoveryAction struct {
	Description string
	Execute     func(ctx interface{}) error
}

// GetRecoveryStrategy returns the recovery strategy for an error
func GetRecoveryStrategy(err *GenerationError) *ErrorRecoveryStrategy {
	switch err.Category {
	case ErrCategoryToolNotFound:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryToolNotFound,
			Description: "Use fallback generation instead of bootstrap tool",
			Actions: []RecoveryAction{
				{
					Description: "Check if fallback generator is available",
				},
				{
					Description: "Switch to fallback generation mode",
				},
				{
					Description: "Retry component generation",
				},
			},
		}

	case ErrCategoryToolExecution:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryToolExecution,
			Description: "Retry execution or fallback to custom generation",
			Actions: []RecoveryAction{
				{
					Description: "Retry tool execution once",
				},
				{
					Description: "If retry fails, switch to fallback generation",
				},
			},
		}

	case ErrCategoryFileSystem:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryFileSystem,
			Description: "Retry file system operation after fixing issues",
			Actions: []RecoveryAction{
				{
					Description: "Check and fix permissions",
				},
				{
					Description: "Verify disk space",
				},
				{
					Description: "Retry operation",
				},
			},
		}

	case ErrCategoryIntegration:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryIntegration,
			Description: "Continue without full integration",
			Actions: []RecoveryAction{
				{
					Description: "Mark integration as incomplete",
				},
				{
					Description: "Generate manual integration instructions",
				},
				{
					Description: "Continue with project generation",
				},
			},
		}

	case ErrCategoryValidation:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryValidation,
			Description: "Continue with warnings",
			Actions: []RecoveryAction{
				{
					Description: "Log validation warnings",
				},
				{
					Description: "Continue with project generation",
				},
			},
		}

	case ErrCategoryTimeout:
		return &ErrorRecoveryStrategy{
			Category:    ErrCategoryTimeout,
			Description: "Retry with increased timeout",
			Actions: []RecoveryAction{
				{
					Description: "Increase timeout duration",
				},
				{
					Description: "Retry operation",
				},
			},
		}

	default:
		return &ErrorRecoveryStrategy{
			Category:    err.Category,
			Description: "No automatic recovery available",
			Actions:     []RecoveryAction{},
		}
	}
}

// ErrorContext provides context for error handling
type ErrorContext struct {
	Operation     string
	Component     string
	Phase         string
	AttemptNumber int
	CanRetry      bool
	CanFallback   bool
}

// ShouldRetry determines if an operation should be retried based on error and context
func ShouldRetry(err *GenerationError, ctx *ErrorContext) bool {
	if !ctx.CanRetry {
		return false
	}

	if ctx.AttemptNumber >= 2 {
		return false // Max 2 attempts
	}

	// Only retry recoverable errors
	if !err.Recoverable {
		return false
	}

	// Retry specific categories
	switch err.Category {
	case ErrCategoryToolExecution:
		return true
	case ErrCategoryFileSystem:
		return true
	case ErrCategoryTimeout:
		return true
	default:
		return false
	}
}

// ShouldFallback determines if fallback generation should be used
func ShouldFallback(err *GenerationError, ctx *ErrorContext) bool {
	if !ctx.CanFallback {
		return false
	}

	// Fallback for tool-related errors
	switch err.Category {
	case ErrCategoryToolNotFound:
		return true
	case ErrCategoryToolExecution:
		return ctx.AttemptNumber >= 2 // Fallback after retry fails
	default:
		return false
	}
}

// FormatError formats an error for display to the user
func FormatError(err error) string {
	genErr := &GenerationError{}
	if errors.As(err, &genErr) {
		var builder strings.Builder

		builder.WriteString(genErr.Error())
		builder.WriteString("\n")

		if suggestions := genErr.GetSuggestions(); suggestions != "" {
			builder.WriteString(suggestions)
		}

		return builder.String()
	}

	return err.Error()
}

// AggregateErrors combines multiple errors into a single error message
func AggregateErrors(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	if len(errors) == 1 {
		return errors[0]
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("multiple errors occurred (%d):\n", len(errors)))

	for i, err := range errors {
		builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}

	return fmt.Errorf("%s", builder.String())
}
