// Package cli provides error handling and exit code management for automation
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// StructuredError represents an error with structured information for automation
type StructuredError struct {
	Type        string                 `json:"type"`
	Code        int                    `json:"code"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
	Context     *ErrorContext          `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// ErrorContext provides context information for structured errors
type ErrorContext struct {
	Command     string            `json:"command"`
	Arguments   []string          `json:"arguments"`
	Flags       map[string]string `json:"flags"`
	WorkingDir  string            `json:"working_dir"`
	Environment string            `json:"environment,omitempty"`
	CI          *CIEnvironment    `json:"ci,omitempty"`
}

// Error implements the error interface
func (e *StructuredError) Error() string {
	return e.Message
}

// ExitCode returns the appropriate exit code for the error
func (e *StructuredError) ExitCode() int {
	return e.Code
}

// ToJSON converts the error to JSON format
func (e *StructuredError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// NewStructuredError creates a new structured error
func NewStructuredError(errorType, message string, code int) *StructuredError {
	return &StructuredError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// WithDetails adds details to the structured error
func (e *StructuredError) WithDetails(key string, value interface{}) *StructuredError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithSuggestions adds suggestions to the structured error
func (e *StructuredError) WithSuggestions(suggestions ...string) *StructuredError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithContext adds context to the structured error
func (e *StructuredError) WithContext(ctx *ErrorContext) *StructuredError {
	e.Context = ctx
	return e
}

// Exit code constants for different error types
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
	ExitCodeInternalError        = 99
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
)

// handleError processes errors and outputs them in the appropriate format
func (c *CLI) handleError(err error, cmd string, args []string) int {
	if err == nil {
		return ExitCodeSuccess
	}

	// Check if it's already a structured error
	if structErr, ok := err.(*StructuredError); ok {
		return c.outputStructuredError(structErr, cmd, args)
	}

	// Check if it's a CLI error from automation.go
	if cliErr, ok := err.(*CLIError); ok {
		structErr := NewStructuredError(ErrorTypeInternal, cliErr.Message, cliErr.Code)
		return c.outputStructuredError(structErr, cmd, args)
	}

	// Convert regular error to structured error
	structErr := NewStructuredError(ErrorTypeInternal, err.Error(), ExitCodeGeneral)
	return c.outputStructuredError(structErr, cmd, args)
}

// outputStructuredError outputs a structured error in the appropriate format
func (c *CLI) outputStructuredError(err *StructuredError, cmd string, args []string) int {
	// Add context if not already present
	if err.Context == nil {
		workingDir, _ := os.Getwd()
		err.Context = &ErrorContext{
			Command:    cmd,
			Arguments:  args,
			WorkingDir: workingDir,
			CI:         c.detectCIEnvironment(),
		}

		// Add environment type
		if err.Context.CI.IsCI {
			err.Context.Environment = "ci"
		} else {
			err.Context.Environment = "local"
		}
	}

	// Set exit code
	c.SetExitCode(err.Code)

	// Check if we should output in machine-readable format
	nonInteractive := c.isNonInteractiveMode()
	outputFormat := "text"
	if c.rootCmd != nil {
		if format, cmdErr := c.rootCmd.PersistentFlags().GetString("output-format"); cmdErr == nil {
			outputFormat = format
		}
	}

	// Output error in appropriate format
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return c.outputMachineReadableError(err, outputFormat)
	}

	// Human-readable error output
	c.ErrorOutput("%s", err.Message)

	// Show details in verbose mode
	if c.verboseMode || c.debugMode {
		if len(err.Details) > 0 {
			c.ErrorOutput("Details:")
			for key, value := range err.Details {
				c.ErrorOutput("  %s: %v", key, value)
			}
		}

		if len(err.Suggestions) > 0 {
			c.ErrorOutput("Suggestions:")
			for _, suggestion := range err.Suggestions {
				c.ErrorOutput("  - %s", suggestion)
			}
		}

		if err.Context != nil {
			c.ErrorOutput("Context:")
			c.ErrorOutput("  Command: %s", err.Context.Command)
			if len(err.Context.Arguments) > 0 {
				c.ErrorOutput("  Arguments: %v", err.Context.Arguments)
			}
			c.ErrorOutput("  Working Directory: %s", err.Context.WorkingDir)
			if err.Context.CI.IsCI {
				c.ErrorOutput("  CI Environment: %s", err.Context.CI.Provider)
			}
		}
	}

	return err.Code
}

// outputMachineReadableError outputs an error in machine-readable format
func (c *CLI) outputMachineReadableError(err *StructuredError, format string) int {
	switch format {
	case "json":
		if jsonData, jsonErr := err.ToJSON(); jsonErr == nil {
			fmt.Fprintln(os.Stderr, string(jsonData))
		} else {
			// Fallback to simple JSON
			fmt.Fprintf(os.Stderr, `{"type":"%s","code":%d,"message":"%s","timestamp":"%s"}%s`,
				err.Type, err.Code, err.Message, err.Timestamp.Format(time.RFC3339), "\n")
		}
	case "yaml":
		// For now, output as JSON since we don't have YAML library
		// This can be enhanced later with proper YAML support
		if jsonData, jsonErr := err.ToJSON(); jsonErr == nil {
			fmt.Fprintln(os.Stderr, string(jsonData))
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: %s (Code: %d)\n", err.Message, err.Code)
	}

	return err.Code
}

// createValidationError creates a structured validation error
func (c *CLI) createValidationError(message string, details map[string]interface{}) *StructuredError {
	err := NewStructuredError(ErrorTypeValidation, message, ExitCodeValidationFailed)
	for key, value := range details {
		err = err.WithDetails(key, value)
	}
	err = err.WithSuggestions(
		"Check the project structure and configuration files",
		"Run with --verbose flag for more details",
		"Use --fix flag to automatically fix common issues",
	)
	return err
}

// createConfigurationError creates a structured configuration error
func (c *CLI) createConfigurationError(message string, configPath string) *StructuredError {
	err := NewStructuredError(ErrorTypeConfiguration, message, ExitCodeConfigurationInvalid)
	err = err.WithDetails("config_path", configPath)
	err = err.WithSuggestions(
		"Check the configuration file syntax",
		"Validate the configuration with 'generator config validate'",
		"Use 'generator config show' to see current configuration",
	)
	return err
}

// createTemplateError creates a structured template error
func (c *CLI) createTemplateError(message string, templateName string) *StructuredError {
	err := NewStructuredError(ErrorTypeTemplate, message, ExitCodeTemplateNotFound)
	err = err.WithDetails("template_name", templateName)
	err = err.WithSuggestions(
		"List available templates with 'generator list-templates'",
		"Check template name spelling",
		"Use 'generator template info <name>' for template details",
	)
	return err
}

// createNetworkError creates a structured network error
//
//nolint:unused // May be used in future network operations
func (c *CLI) createNetworkError(message string, url string) *StructuredError {
	err := NewStructuredError(ErrorTypeNetwork, message, ExitCodeNetworkError)
	err = err.WithDetails("url", url)
	err = err.WithSuggestions(
		"Check your internet connection",
		"Use --offline flag to work with cached data",
		"Check if the URL is accessible",
	)
	return err
}

// createFileSystemError creates a structured filesystem error
//
//nolint:unused // May be used in future filesystem operations
func (c *CLI) createFileSystemError(message string, path string) *StructuredError {
	err := NewStructuredError(ErrorTypeFileSystem, message, ExitCodeFileSystemError)
	err = err.WithDetails("path", path)
	err = err.WithSuggestions(
		"Check if the path exists and is accessible",
		"Verify file permissions",
		"Ensure sufficient disk space",
	)
	return err
}

// createPermissionError creates a structured permission error
//
//nolint:unused // May be used in future permission-related operations
func (c *CLI) createPermissionError(message string, path string) *StructuredError {
	err := NewStructuredError(ErrorTypePermission, message, ExitCodePermissionDenied)
	err = err.WithDetails("path", path)
	err = err.WithSuggestions(
		"Check file/directory permissions",
		"Run with appropriate user privileges",
		"Verify ownership of the target directory",
	)
	return err
}

// createAuditError creates a structured audit error
func (c *CLI) createAuditError(message string, score float64) *StructuredError {
	err := NewStructuredError(ErrorTypeAudit, message, ExitCodeAuditFailed)
	err = err.WithDetails("score", score)
	err = err.WithSuggestions(
		"Review audit recommendations",
		"Fix high-priority security issues",
		"Improve code quality metrics",
	)
	return err
}

// Success response structure for machine-readable output
type SuccessResponse struct {
	Success   bool            `json:"success"`
	Message   string          `json:"message,omitempty"`
	Data      interface{}     `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Context   *SuccessContext `json:"context,omitempty"`
}

// SuccessContext provides context for successful operations
type SuccessContext struct {
	Command   string         `json:"command"`
	Arguments []string       `json:"arguments"`
	Duration  string         `json:"duration,omitempty"`
	CI        *CIEnvironment `json:"ci,omitempty"`
}

// outputSuccess outputs a success response in machine-readable format
func (c *CLI) outputSuccess(message string, data interface{}, cmd string, args []string) error {
	nonInteractive := c.isNonInteractiveMode()
	outputFormat := "text"
	if c.rootCmd != nil {
		if format, err := c.rootCmd.PersistentFlags().GetString("output-format"); err == nil {
			outputFormat = format
		}
	}

	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		response := &SuccessResponse{
			Success:   true,
			Message:   message,
			Data:      data,
			Timestamp: time.Now(),
			Context: &SuccessContext{
				Command:   cmd,
				Arguments: args,
				CI:        c.detectCIEnvironment(),
			},
		}

		return c.outputMachineReadable(response, outputFormat)
	}

	// Human-readable output
	if message != "" && !c.quietMode {
		c.SuccessOutput(message)
	}

	return nil
}
