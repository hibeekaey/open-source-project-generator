package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// ExitCode represents CLI exit codes
type ExitCode int

const (
	// ExitSuccess indicates successful completion
	ExitSuccess ExitCode = 0

	// ExitConfigError indicates configuration validation failure
	ExitConfigError ExitCode = 1

	// ExitToolsMissing indicates required tools are not available
	ExitToolsMissing ExitCode = 2

	// ExitGenerationFailed indicates component generation failure
	ExitGenerationFailed ExitCode = 3

	// ExitFileSystemError indicates file system operation failure
	ExitFileSystemError ExitCode = 4

	// ExitUserCancelled indicates user cancelled the operation
	ExitUserCancelled ExitCode = 5
)

// ExitCodeHandler manages exit codes and error categorization
type ExitCodeHandler struct {
	logger *logger.Logger
}

// NewExitCodeHandler creates a new exit code handler
func NewExitCodeHandler(log *logger.Logger) *ExitCodeHandler {
	return &ExitCodeHandler{
		logger: log,
	}
}

// DetermineExitCode analyzes an error and returns appropriate exit code
func (ech *ExitCodeHandler) DetermineExitCode(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}

	// Check for user cancellation
	if errors.Is(err, ErrUserCancelled) {
		return ExitUserCancelled
	}

	// Check for CLIError types
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr.ExitCode
	}

	// Check for GenerationError types
	var genErr *orchestrator.GenerationError
	if errors.As(err, &genErr) {
		return ech.exitCodeFromGenerationError(genErr)
	}

	// Check for specific error messages
	errMsg := err.Error()

	// Configuration errors
	if containsAny(errMsg, []string{"validation failed", "invalid configuration", "config error"}) {
		return ExitConfigError
	}

	// Tool errors
	if containsAny(errMsg, []string{"tool not found", "tool not available", "missing tool", "not whitelisted"}) {
		return ExitToolsMissing
	}

	// Generation errors
	if containsAny(errMsg, []string{"generation failed", "component generation", "bootstrap failed", "fallback failed"}) {
		return ExitGenerationFailed
	}

	// File system errors
	if containsAny(errMsg, []string{"file system", "permission denied", "no such file", "directory", "failed to create", "failed to write"}) {
		return ExitFileSystemError
	}

	// Default to generation failed for unknown errors
	return ExitGenerationFailed
}

// exitCodeFromGenerationError converts a GenerationError to an exit code
func (ech *ExitCodeHandler) exitCodeFromGenerationError(genErr *orchestrator.GenerationError) ExitCode {
	switch genErr.Category {
	case orchestrator.ErrCategoryInvalidConfig:
		return ExitConfigError
	case orchestrator.ErrCategoryToolNotFound:
		return ExitToolsMissing
	case orchestrator.ErrCategoryToolExecution:
		return ExitGenerationFailed
	case orchestrator.ErrCategoryFileSystem:
		return ExitFileSystemError
	case orchestrator.ErrCategoryIntegration:
		return ExitGenerationFailed
	case orchestrator.ErrCategoryValidation:
		return ExitGenerationFailed
	case orchestrator.ErrCategoryTimeout:
		return ExitGenerationFailed
	case orchestrator.ErrCategorySecurity:
		return ExitConfigError
	default:
		return ExitGenerationFailed
	}
}

// ExitWithCode logs the error and exits with the appropriate code
func (ech *ExitCodeHandler) ExitWithCode(err error, code ExitCode) {
	if err != nil {
		ech.logger.Error(fmt.Sprintf("Error: %v", err))

		// Display suggestions if available
		var cliErr *CLIError
		if errors.As(err, &cliErr) {
			if suggestions := cliErr.GetSuggestions(); suggestions != "" {
				fmt.Fprintf(os.Stderr, "%s", suggestions)
			}
		}

		// Display suggestions from GenerationError
		var genErr *orchestrator.GenerationError
		if errors.As(err, &genErr) {
			if suggestions := genErr.GetSuggestions(); suggestions != "" {
				fmt.Fprintf(os.Stderr, "%s", suggestions)
			}
		}
	}

	ech.logExitCode(code)
	os.Exit(int(code))
}

// ExitWithMessage logs a message and exits with the specified code
func (ech *ExitCodeHandler) ExitWithMessage(message string, code ExitCode) {
	if message != "" {
		if code == ExitSuccess {
			ech.logger.Success(message)
		} else {
			ech.logger.Error(message)
		}
	}

	ech.logExitCode(code)
	os.Exit(int(code))
}

// logExitCode logs the exit code for debugging
func (ech *ExitCodeHandler) logExitCode(code ExitCode) {
	reason := ech.getExitCodeReason(code)
	ech.logger.Debug(fmt.Sprintf("Exiting with code %d: %s", code, reason))
}

// getExitCodeReason returns a human-readable reason for an exit code
func (ech *ExitCodeHandler) getExitCodeReason(code ExitCode) string {
	switch code {
	case ExitSuccess:
		return "Success"
	case ExitConfigError:
		return "Configuration validation failed"
	case ExitToolsMissing:
		return "Required tools are missing"
	case ExitGenerationFailed:
		return "Component generation failed"
	case ExitFileSystemError:
		return "File system operation failed"
	case ExitUserCancelled:
		return "User cancelled the operation"
	default:
		return "Unknown error"
	}
}

// GetExitCodeDescription returns a detailed description of an exit code
func (ech *ExitCodeHandler) GetExitCodeDescription(code ExitCode) string {
	switch code {
	case ExitSuccess:
		return "The operation completed successfully."
	case ExitConfigError:
		return "The configuration file is invalid or contains errors. Please check your configuration and try again."
	case ExitToolsMissing:
		return "One or more required tools are not installed or not available in PATH. Please install the missing tools and try again."
	case ExitGenerationFailed:
		return "Failed to generate one or more components. Check the error messages above for details."
	case ExitFileSystemError:
		return "A file system operation failed. This could be due to permissions, disk space, or other file system issues."
	case ExitUserCancelled:
		return "The operation was cancelled by the user."
	default:
		return "An unknown error occurred."
	}
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// ErrUserCancelled is returned when the user cancels an operation
var ErrUserCancelled = errors.New("operation cancelled by user")

// CLIError represents a CLI-specific error with exit code and suggestions
type CLIError struct {
	Category    string
	Message     string
	Cause       error
	ExitCode    ExitCode
	Suggestions []string
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Category, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Category, e.Message)
}

// Unwrap returns the underlying cause error
func (e *CLIError) Unwrap() error {
	return e.Cause
}

// WithSuggestions adds recovery suggestions to the error
func (e *CLIError) WithSuggestions(suggestions ...string) *CLIError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithExitCode sets the exit code for the error
func (e *CLIError) WithExitCode(code ExitCode) *CLIError {
	e.ExitCode = code
	return e
}

// GetSuggestions returns formatted suggestions for resolving the error
func (e *CLIError) GetSuggestions() string {
	if len(e.Suggestions) == 0 {
		return ""
	}

	result := "\nNext Steps:\n"
	for i, suggestion := range e.Suggestions {
		result += fmt.Sprintf("  %d. %s\n", i+1, suggestion)
	}
	return result
}

// NewCLIError creates a new CLI error
func NewCLIError(category string, message string, cause error) *CLIError {
	return &CLIError{
		Category:    category,
		Message:     message,
		Cause:       cause,
		ExitCode:    ExitGenerationFailed, // Default exit code
		Suggestions: make([]string, 0),
	}
}

// NewConfigError creates a configuration error
func NewConfigError(message string, cause error) *CLIError {
	return &CLIError{
		Category: "configuration",
		Message:  message,
		Cause:    cause,
		ExitCode: ExitConfigError,
		Suggestions: []string{
			"Check your configuration file for errors",
			"Refer to the configuration schema documentation",
			"Use 'generator init-config' to generate a valid configuration template",
		},
	}
}

// NewToolError creates a tool-related error
func NewToolError(toolName string, message string, cause error) *CLIError {
	return &CLIError{
		Category: "tool",
		Message:  fmt.Sprintf("tool '%s': %s", toolName, message),
		Cause:    cause,
		ExitCode: ExitToolsMissing,
		Suggestions: []string{
			fmt.Sprintf("Install %s on your system", toolName),
			"Use --no-external-tools flag to force fallback generation",
			"Run 'generator check-tools' to see installation instructions",
		},
	}
}

// NewGenerationError creates a generation error
func NewGenerationError(component string, message string, cause error) *CLIError {
	return &CLIError{
		Category: "generation",
		Message:  fmt.Sprintf("component '%s': %s", component, message),
		Cause:    cause,
		ExitCode: ExitGenerationFailed,
		Suggestions: []string{
			"Check the error messages above for details",
			"Try running with --verbose flag for more information",
			"Use --no-external-tools to try fallback generation",
		},
	}
}

// NewFileSystemError creates a file system error
func NewFileSystemError(operation string, path string, cause error) *CLIError {
	return &CLIError{
		Category: "filesystem",
		Message:  fmt.Sprintf("operation '%s' failed for path '%s'", operation, path),
		Cause:    cause,
		ExitCode: ExitFileSystemError,
		Suggestions: []string{
			"Check file system permissions",
			"Verify disk space is available",
			"Ensure the path is accessible",
		},
	}
}

// NewUserCancelledError creates a user cancellation error
func NewUserCancelledError(message string) *CLIError {
	return &CLIError{
		Category: "user",
		Message:  message,
		Cause:    nil,
		ExitCode: ExitUserCancelled,
		Suggestions: []string{
			"Run the command again when ready",
		},
	}
}
