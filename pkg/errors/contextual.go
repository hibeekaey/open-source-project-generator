// Package errors provides contextual error messages with actionable suggestions
package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ContextualErrorBuilder helps build errors with rich context and suggestions
type ContextualErrorBuilder struct {
	errorType   string
	message     string
	code        int
	context     *ErrorContext
	details     map[string]interface{}
	suggestions []string
	cause       error
	severity    Severity
	recoverable bool
}

// NewContextualError creates a new contextual error builder
func NewContextualError(errorType, message string, code int) *ContextualErrorBuilder {
	return &ContextualErrorBuilder{
		errorType:   errorType,
		message:     message,
		code:        code,
		details:     make(map[string]interface{}),
		suggestions: []string{},
		severity:    determineSeverity(errorType, code),
		recoverable: determineRecoverable(errorType, code),
	}
}

// WithOperation adds operation context to the error
func (b *ContextualErrorBuilder) WithOperation(operation string) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.Operation = operation
	return b
}

// WithComponent adds component context to the error
func (b *ContextualErrorBuilder) WithComponent(component string) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.Component = component
	return b
}

// WithFile adds file context to the error
func (b *ContextualErrorBuilder) WithFile(file string, line int) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.File = file
	b.context.Line = line
	return b
}

// WithCommand adds command context to the error
func (b *ContextualErrorBuilder) WithCommand(command string, args []string, flags map[string]string) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.Command = command
	b.context.Arguments = args
	b.context.Flags = flags
	return b
}

// WithEnvironment adds environment context to the error
func (b *ContextualErrorBuilder) WithEnvironment(env string, ci *CIEnvironment) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.Environment = env
	b.context.CI = ci
	return b
}

// WithWorkingDirectory adds working directory context
func (b *ContextualErrorBuilder) WithWorkingDirectory(dir string) *ContextualErrorBuilder {
	if b.context == nil {
		b.context = &ErrorContext{}
	}
	b.context.WorkingDir = dir
	return b
}

// WithDetail adds a detail to the error
func (b *ContextualErrorBuilder) WithDetail(key string, value interface{}) *ContextualErrorBuilder {
	b.details[key] = value
	return b
}

// WithSuggestion adds a suggestion to the error
func (b *ContextualErrorBuilder) WithSuggestion(suggestion string) *ContextualErrorBuilder {
	b.suggestions = append(b.suggestions, suggestion)
	return b
}

// WithCause adds the underlying cause to the error
func (b *ContextualErrorBuilder) WithCause(cause error) *ContextualErrorBuilder {
	b.cause = cause
	return b
}

// WithSeverity sets the error severity
func (b *ContextualErrorBuilder) WithSeverity(severity Severity) *ContextualErrorBuilder {
	b.severity = severity
	return b
}

// WithRecoverable sets whether the error is recoverable
func (b *ContextualErrorBuilder) WithRecoverable(recoverable bool) *ContextualErrorBuilder {
	b.recoverable = recoverable
	return b
}

// Build creates the final CLIError
func (b *ContextualErrorBuilder) Build() *CLIError {
	err := NewCLIError(b.errorType, b.message, b.code)

	if b.context != nil {
		err = err.WithContext(b.context)
	}

	for key, value := range b.details {
		err = err.WithDetails(key, value)
	}

	if len(b.suggestions) > 0 {
		err = err.WithSuggestions(b.suggestions...)
	}

	if b.cause != nil {
		err = err.WithCause(b.cause)
	}

	err = err.WithSeverity(b.severity)
	err = err.WithRecoverable(b.recoverable)

	return err
}

// ContextualSuggestionGenerator generates contextual suggestions based on error details
type ContextualSuggestionGenerator struct{}

// NewContextualSuggestionGenerator creates a new suggestion generator
func NewContextualSuggestionGenerator() *ContextualSuggestionGenerator {
	return &ContextualSuggestionGenerator{}
}

// GenerateSuggestions generates contextual suggestions for an error
func (g *ContextualSuggestionGenerator) GenerateSuggestions(err *CLIError) []string {
	var suggestions []string

	// Add type-specific suggestions
	suggestions = append(suggestions, g.getTypeSpecificSuggestions(err)...)

	// Add context-specific suggestions
	suggestions = append(suggestions, g.getContextSpecificSuggestions(err)...)

	// Add environment-specific suggestions
	suggestions = append(suggestions, g.getEnvironmentSpecificSuggestions(err)...)

	// Add cause-specific suggestions
	suggestions = append(suggestions, g.getCauseSpecificSuggestions(err)...)

	// Remove duplicates and return
	return removeDuplicateStrings(suggestions)
}

// getTypeSpecificSuggestions returns suggestions based on error type
func (g *ContextualSuggestionGenerator) getTypeSpecificSuggestions(err *CLIError) []string {
	switch err.Type {
	case ErrorTypeValidation:
		return []string{
			"Use --fix flag to automatically fix common validation issues",
			"Run 'generator config validate' to check configuration",
			"Check project structure against template requirements",
		}
	case ErrorTypeConfiguration:
		return []string{
			"Use 'generator config show' to see current configuration",
			"Validate configuration with 'generator config validate'",
			"Check configuration file syntax (YAML/JSON)",
		}
	case ErrorTypeTemplate:
		return []string{
			"List available templates with 'generator list-templates'",
			"Get template info with 'generator template info <name>'",
			"Check template name spelling and availability",
		}
	case ErrorTypeNetwork:
		return []string{
			"Check internet connection",
			"Use --offline flag to work with cached data",
			"Verify proxy settings if applicable",
		}
	case ErrorTypeFileSystem:
		return []string{
			"Check file and directory permissions",
			"Verify sufficient disk space",
			"Ensure target directory is writable",
		}
	case ErrorTypePermission:
		return []string{
			"Run with appropriate user privileges",
			"Check file/directory ownership",
			"Verify required permissions are granted",
		}
	case ErrorTypeCache:
		return []string{
			"Clear cache with 'generator cache clear'",
			"Check cache directory permissions",
			"Try running without cache",
		}
	case ErrorTypeVersion:
		return []string{
			"Check for updates with 'generator version --check-updates'",
			"Verify component compatibility",
			"Use --offline flag to skip version checks",
		}
	case ErrorTypeAudit:
		return []string{
			"Review audit recommendations",
			"Fix high-priority issues first",
			"Check project dependencies for vulnerabilities",
		}
	case ErrorTypeGeneration:
		return []string{
			"Use --dry-run to preview generation",
			"Check template and configuration validity",
			"Verify output directory permissions",
		}
	case ErrorTypeDependency:
		return []string{
			"Update dependencies to compatible versions",
			"Check dependency documentation",
			"Use --update-versions flag for latest versions",
		}
	case ErrorTypeSecurity:
		return []string{
			"Address security issues immediately",
			"Review security audit recommendations",
			"Update vulnerable dependencies",
		}
	case ErrorTypeUser:
		return []string{
			"Check command syntax and arguments",
			"Use --help flag for usage information",
			"Verify input format and values",
		}
	default:
		return []string{
			"Check error details for specific guidance",
			"Consult documentation for troubleshooting",
			"Use --verbose flag for more information",
		}
	}
}

// getContextSpecificSuggestions returns suggestions based on error context
func (g *ContextualSuggestionGenerator) getContextSpecificSuggestions(err *CLIError) []string {
	var suggestions []string

	if err.Context == nil {
		return suggestions
	}

	// Command-specific suggestions
	if err.Context.Command != "" {
		switch err.Context.Command {
		case "generate":
			suggestions = append(suggestions, []string{
				"Use --config flag to specify configuration file",
				"Try --minimal flag for basic project structure",
				"Use --dry-run to preview without creating files",
			}...)
		case "validate":
			suggestions = append(suggestions, []string{
				"Use --fix flag to automatically fix issues",
				"Generate report with --report flag",
				"Check specific rules with --rules flag",
			}...)
		case "audit":
			suggestions = append(suggestions, []string{
				"Focus on specific categories (--security, --quality)",
				"Generate detailed report with --detailed flag",
				"Set minimum score with --min-score flag",
			}...)
		}
	}

	// Operation-specific suggestions
	if err.Context.Operation != "" {
		switch {
		case strings.Contains(err.Context.Operation, "file"):
			suggestions = append(suggestions, "Check file path and permissions")
		case strings.Contains(err.Context.Operation, "network"):
			suggestions = append(suggestions, "Verify network connectivity")
		case strings.Contains(err.Context.Operation, "config"):
			suggestions = append(suggestions, "Validate configuration syntax")
		}
	}

	// File-specific suggestions
	if err.Context.File != "" {
		ext := filepath.Ext(err.Context.File)
		switch ext {
		case ".yaml", ".yml":
			suggestions = append(suggestions, "Check YAML syntax and indentation")
		case ".json":
			suggestions = append(suggestions, "Validate JSON syntax")
		case ".go":
			suggestions = append(suggestions, "Check Go syntax and imports")
		case ".js", ".ts":
			suggestions = append(suggestions, "Check JavaScript/TypeScript syntax")
		}

		if err.Context.Line > 0 {
			suggestions = append(suggestions, fmt.Sprintf("Check line %d in %s", err.Context.Line, err.Context.File))
		}
	}

	return suggestions
}

// getEnvironmentSpecificSuggestions returns suggestions based on environment
func (g *ContextualSuggestionGenerator) getEnvironmentSpecificSuggestions(err *CLIError) []string {
	var suggestions []string

	if err.Context == nil {
		return suggestions
	}

	// CI environment suggestions
	if err.Context.CI != nil && err.Context.CI.IsCI {
		suggestions = append(suggestions, []string{
			"Use --non-interactive flag in CI environments",
			"Set appropriate environment variables",
			"Check CI-specific configuration",
		}...)

		switch err.Context.CI.Provider {
		case "github":
			suggestions = append(suggestions, "Check GitHub Actions workflow configuration")
		case "gitlab":
			suggestions = append(suggestions, "Check GitLab CI configuration")
		case "jenkins":
			suggestions = append(suggestions, "Check Jenkins pipeline configuration")
		}
	}

	// Working directory suggestions
	if err.Context.WorkingDir != "" {
		if !isValidProjectDirectory(err.Context.WorkingDir) {
			suggestions = append(suggestions, "Ensure you're in a valid project directory")
		}
	}

	return suggestions
}

// getCauseSpecificSuggestions returns suggestions based on the underlying cause
func (g *ContextualSuggestionGenerator) getCauseSpecificSuggestions(err *CLIError) []string {
	var suggestions []string

	if err.Cause == nil {
		return suggestions
	}

	causeMsg := err.Cause.Error()

	// Common error patterns
	switch {
	case strings.Contains(causeMsg, "permission denied"):
		suggestions = append(suggestions, []string{
			"Check file/directory permissions",
			"Run with appropriate user privileges",
			"Verify ownership of target files",
		}...)
	case strings.Contains(causeMsg, "no such file or directory"):
		suggestions = append(suggestions, []string{
			"Verify the file path exists",
			"Check for typos in file names",
			"Ensure required files are present",
		}...)
	case strings.Contains(causeMsg, "connection refused"):
		suggestions = append(suggestions, []string{
			"Check if the service is running",
			"Verify network connectivity",
			"Check firewall settings",
		}...)
	case strings.Contains(causeMsg, "timeout"):
		suggestions = append(suggestions, []string{
			"Increase timeout settings if possible",
			"Check network latency",
			"Retry the operation",
		}...)
	case strings.Contains(causeMsg, "invalid syntax"):
		suggestions = append(suggestions, []string{
			"Check file syntax",
			"Validate configuration format",
			"Review documentation for correct syntax",
		}...)
	case strings.Contains(causeMsg, "not found"):
		suggestions = append(suggestions, []string{
			"Check if the resource exists",
			"Verify spelling and case sensitivity",
			"Ensure required dependencies are installed",
		}...)
	}

	return suggestions
}

// Helper functions

// isValidProjectDirectory checks if a directory looks like a valid project directory
func isValidProjectDirectory(dir string) bool {
	// Check for common project files
	commonFiles := []string{
		"package.json", "go.mod", "Cargo.toml", "requirements.txt",
		"pom.xml", "build.gradle", "Makefile", ".git",
	}

	for _, file := range commonFiles {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return true
		}
	}

	return false
}

// removeDuplicateStrings removes duplicate strings from a slice
func removeDuplicateStrings(strings []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range strings {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// ErrorContextCollector collects context information for errors
type ErrorContextCollector struct{}

// NewErrorContextCollector creates a new context collector
func NewErrorContextCollector() *ErrorContextCollector {
	return &ErrorContextCollector{}
}

// CollectContext collects comprehensive context information
func (c *ErrorContextCollector) CollectContext() *ErrorContext {
	ctx := &ErrorContext{}

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		ctx.WorkingDir = wd
	}

	// Detect CI environment
	ctx.CI = c.detectCIEnvironment()

	// Set environment type
	if ctx.CI.IsCI {
		ctx.Environment = "ci"
	} else {
		ctx.Environment = "local"
	}

	return ctx
}

// CollectCommandContext collects context for a specific command
func (c *ErrorContextCollector) CollectCommandContext(command string, args []string, flags map[string]string) *ErrorContext {
	ctx := c.CollectContext()
	ctx.Command = command
	ctx.Arguments = args
	ctx.Flags = flags
	return ctx
}

// detectCIEnvironment detects if running in a CI environment
func (c *ErrorContextCollector) detectCIEnvironment() *CIEnvironment {
	ci := &CIEnvironment{}

	// Check common CI environment variables
	if os.Getenv("CI") == "true" || os.Getenv("CONTINUOUS_INTEGRATION") == "true" {
		ci.IsCI = true
	}

	// Detect specific CI providers
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ci.Provider = "github"
		ci.JobID = os.Getenv("GITHUB_RUN_ID")
		ci.BuildID = os.Getenv("GITHUB_RUN_NUMBER")
	} else if os.Getenv("GITLAB_CI") == "true" {
		ci.Provider = "gitlab"
		ci.JobID = os.Getenv("CI_JOB_ID")
		ci.BuildID = os.Getenv("CI_PIPELINE_ID")
	} else if os.Getenv("JENKINS_URL") != "" {
		ci.Provider = "jenkins"
		ci.JobID = os.Getenv("BUILD_ID")
		ci.BuildID = os.Getenv("BUILD_NUMBER")
	} else if os.Getenv("CIRCLECI") == "true" {
		ci.Provider = "circleci"
		ci.JobID = os.Getenv("CIRCLE_BUILD_NUM")
		ci.BuildID = os.Getenv("CIRCLE_BUILD_NUM")
	} else if os.Getenv("TRAVIS") == "true" {
		ci.Provider = "travis"
		ci.JobID = os.Getenv("TRAVIS_JOB_ID")
		ci.BuildID = os.Getenv("TRAVIS_BUILD_ID")
	}

	return ci
}

// GetCallerInfo gets information about the caller
func GetCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}

	// Get just the filename, not the full path
	return filepath.Base(file), line
}

// FormatErrorForUser formats an error message for user-friendly display
func FormatErrorForUser(err *CLIError, verbose bool) string {
	var parts []string

	// Add error type and message
	parts = append(parts, fmt.Sprintf("Error: %s", err.Message))

	// Add context in verbose mode
	if verbose && err.Context != nil {
		if err.Context.Command != "" {
			parts = append(parts, fmt.Sprintf("Command: %s", err.Context.Command))
		}
		if err.Context.Operation != "" {
			parts = append(parts, fmt.Sprintf("Operation: %s", err.Context.Operation))
		}
		if err.Context.File != "" {
			if err.Context.Line > 0 {
				parts = append(parts, fmt.Sprintf("File: %s:%d", err.Context.File, err.Context.Line))
			} else {
				parts = append(parts, fmt.Sprintf("File: %s", err.Context.File))
			}
		}
	}

	// Add suggestions
	if len(err.Suggestions) > 0 {
		parts = append(parts, "")
		parts = append(parts, "Suggestions:")
		for _, suggestion := range err.Suggestions {
			parts = append(parts, fmt.Sprintf("  • %s", suggestion))
		}
	}

	// Add details in verbose mode
	if verbose && len(err.Details) > 0 {
		parts = append(parts, "")
		parts = append(parts, "Details:")
		for key, value := range err.Details {
			parts = append(parts, fmt.Sprintf("  %s: %v", key, value))
		}
	}

	return strings.Join(parts, "\n")
}
