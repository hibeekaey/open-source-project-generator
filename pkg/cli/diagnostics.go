// Package cli provides CLI error types and utilities.
package cli

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// DiagnosticInfo contains detailed diagnostic information for errors
type DiagnosticInfo struct {
	Timestamp    time.Time
	Error        error
	StackTrace   string
	GoVersion    string
	OS           string
	Architecture string
	Context      map[string]interface{}
}

// DiagnosticsCollector collects and formats diagnostic information
type DiagnosticsCollector struct {
	logger  *logger.Logger
	verbose bool
}

// NewDiagnosticsCollector creates a new diagnostics collector
func NewDiagnosticsCollector(log *logger.Logger, verbose bool) *DiagnosticsCollector {
	return &DiagnosticsCollector{
		logger:  log,
		verbose: verbose,
	}
}

// CollectDiagnostics collects diagnostic information for an error
func (dc *DiagnosticsCollector) CollectDiagnostics(err error, context map[string]interface{}) *DiagnosticInfo {
	if err == nil {
		return nil
	}

	info := &DiagnosticInfo{
		Timestamp:    time.Now(),
		Error:        err,
		GoVersion:    runtime.Version(),
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Context:      context,
	}

	// Collect stack trace if verbose mode is enabled
	if dc.verbose {
		info.StackTrace = dc.captureStackTrace()
	}

	return info
}

// FormatDiagnostics formats diagnostic information for display
func (dc *DiagnosticsCollector) FormatDiagnostics(info *DiagnosticInfo) string {
	if info == nil {
		return ""
	}

	var builder strings.Builder

	builder.WriteString("\n")
	builder.WriteString(separator("="))
	builder.WriteString("\nDIAGNOSTIC INFORMATION\n")
	builder.WriteString(separator("="))
	builder.WriteString("\n\n")

	// Basic information
	builder.WriteString(fmt.Sprintf("Timestamp: %s\n", info.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("Go Version: %s\n", info.GoVersion))
	builder.WriteString(fmt.Sprintf("OS: %s\n", info.OS))
	builder.WriteString(fmt.Sprintf("Architecture: %s\n", info.Architecture))
	builder.WriteString("\n")

	// Error information
	builder.WriteString("Error Details:\n")
	builder.WriteString(separator("-"))
	builder.WriteString("\n")
	builder.WriteString(dc.formatErrorDetails(info.Error))
	builder.WriteString("\n")

	// Context information
	if len(info.Context) > 0 {
		builder.WriteString("\nContext:\n")
		builder.WriteString(separator("-"))
		builder.WriteString("\n")
		for key, value := range info.Context {
			builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
		builder.WriteString("\n")
	}

	// Stack trace (only in verbose mode)
	if dc.verbose && info.StackTrace != "" {
		builder.WriteString("\nStack Trace:\n")
		builder.WriteString(separator("-"))
		builder.WriteString("\n")
		builder.WriteString(info.StackTrace)
		builder.WriteString("\n")
	}

	builder.WriteString(separator("="))
	builder.WriteString("\n")

	return builder.String()
}

// LogDiagnostics logs diagnostic information
func (dc *DiagnosticsCollector) LogDiagnostics(err error, context map[string]interface{}) {
	if !dc.verbose {
		return
	}

	info := dc.CollectDiagnostics(err, context)
	if info == nil {
		return
	}

	// Log to debug level
	dc.logger.Debug("Collecting diagnostic information...")

	// Log error details
	dc.logger.Debug(fmt.Sprintf("Error Type: %T", err))
	dc.logger.Debug(fmt.Sprintf("Error Message: %s", err.Error()))

	// Log context
	for key, value := range context {
		dc.logger.Debug(fmt.Sprintf("Context[%s]: %v", key, value))
	}

	// Log stack trace
	if info.StackTrace != "" {
		dc.logger.Debug("Stack Trace:")
		for _, line := range strings.Split(info.StackTrace, "\n") {
			if line != "" {
				dc.logger.Debug("  " + line)
			}
		}
	}
}

// formatErrorDetails formats detailed error information
func (dc *DiagnosticsCollector) formatErrorDetails(err error) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Type: %T\n", err))
	builder.WriteString(fmt.Sprintf("Message: %s\n", err.Error()))

	// Check for CLIError
	cliErr := &CLIError{}
	if errors.As(err, &cliErr) {
		builder.WriteString(fmt.Sprintf("Category: %s\n", cliErr.Category))
		builder.WriteString(fmt.Sprintf("Exit Code: %d\n", cliErr.ExitCode))

		if cliErr.Cause != nil {
			builder.WriteString(fmt.Sprintf("Cause: %v\n", cliErr.Cause))
			builder.WriteString(fmt.Sprintf("Cause Type: %T\n", cliErr.Cause))
		}

		if len(cliErr.Suggestions) > 0 {
			builder.WriteString("Suggestions:\n")
			for i, suggestion := range cliErr.Suggestions {
				builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
			}
		}
	}

	// Check for GenerationError
	genErr := &orchestrator.GenerationError{}
	if errors.As(err, &genErr) {
		builder.WriteString(fmt.Sprintf("Category: %s\n", genErr.Category))
		builder.WriteString(fmt.Sprintf("Component: %s\n", genErr.Component))
		builder.WriteString(fmt.Sprintf("Recoverable: %t\n", genErr.Recoverable))

		if genErr.Cause != nil {
			builder.WriteString(fmt.Sprintf("Cause: %v\n", genErr.Cause))
			builder.WriteString(fmt.Sprintf("Cause Type: %T\n", genErr.Cause))
		}

		if len(genErr.Suggestions) > 0 {
			builder.WriteString("Suggestions:\n")
			for i, suggestion := range genErr.Suggestions {
				builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
			}
		}
	}

	// Unwrap chain
	unwrapped := err
	depth := 0
	for unwrapped != nil {
		if unwrapper, ok := unwrapped.(interface{ Unwrap() error }); ok {
			unwrapped = unwrapper.Unwrap()
			if unwrapped != nil {
				depth++
				builder.WriteString(fmt.Sprintf("\nWrapped Error (depth %d):\n", depth))
				builder.WriteString(fmt.Sprintf("  Type: %T\n", unwrapped))
				builder.WriteString(fmt.Sprintf("  Message: %s\n", unwrapped.Error()))
			}
		} else {
			break
		}
	}

	return builder.String()
}

// captureStackTrace captures the current stack trace
func (dc *DiagnosticsCollector) captureStackTrace() string {
	// Get stack trace from debug package
	stack := debug.Stack()

	// Parse and format the stack trace
	lines := strings.Split(string(stack), "\n")

	var builder strings.Builder

	// Skip the first few lines (goroutine info and this function)
	skipLines := 7
	for i := skipLines; i < len(lines); i++ {
		line := lines[i]

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Format the line
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

// GetSystemInfo returns system information for diagnostics
func (dc *DiagnosticsCollector) GetSystemInfo() map[string]interface{} {
	info := make(map[string]interface{})

	info["go_version"] = runtime.Version()
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH
	info["num_cpu"] = runtime.NumCPU()
	info["num_goroutine"] = runtime.NumGoroutine()

	// Memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	info["memory_alloc_mb"] = memStats.Alloc / 1024 / 1024
	info["memory_total_alloc_mb"] = memStats.TotalAlloc / 1024 / 1024
	info["memory_sys_mb"] = memStats.Sys / 1024 / 1024
	info["num_gc"] = memStats.NumGC

	return info
}

// LogSystemInfo logs system information for diagnostics
func (dc *DiagnosticsCollector) LogSystemInfo() {
	if !dc.verbose {
		return
	}

	dc.logger.Debug("System Information:")

	info := dc.GetSystemInfo()
	for key, value := range info {
		dc.logger.Debug(fmt.Sprintf("  %s: %v", key, value))
	}
}

// FormatVerboseError formats an error with full diagnostic information
func (dc *DiagnosticsCollector) FormatVerboseError(err error, context map[string]interface{}) string {
	if err == nil {
		return ""
	}

	if !dc.verbose {
		// Non-verbose mode: just return the error message
		return err.Error()
	}

	// Verbose mode: collect and format full diagnostics
	info := dc.CollectDiagnostics(err, context)
	return dc.FormatDiagnostics(info)
}

// separator creates a separator line
func separator(char string) string {
	return strings.Repeat(char, 70)
}
