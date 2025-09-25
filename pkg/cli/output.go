// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
//
// This package handles comprehensive CLI operations including:
//   - Advanced project generation with multiple options
//   - Project validation and auditing
//   - Template management and discovery
//   - Configuration management
//   - Version and update management
//   - Cache management
//   - Logging and debugging support
package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Color constants for beautiful CLI output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// Color helper functions for beautiful output
func (c *CLI) colorize(color, text string) string {
	if c.quietMode {
		return text // No colors in quiet mode
	}
	return color + text + ColorReset
}

func (c *CLI) success(text string) string {
	return c.colorize(ColorGreen+ColorBold, text)
}

func (c *CLI) info(text string) string {
	return c.colorize(ColorBlue, text)
}

func (c *CLI) warning(text string) string {
	return c.colorize(ColorYellow, text)
}

func (c *CLI) error(text string) string {
	return c.colorize(ColorRed+ColorBold, text)
}

func (c *CLI) highlight(text string) string {
	return c.colorize(ColorCyan+ColorBold, text)
}

func (c *CLI) dim(text string) string {
	return c.colorize(ColorDim, text)
}

// Verbose output methods for enhanced debugging and user feedback

// VerboseOutput prints verbose information if verbose mode is enabled
func (c *CLI) VerboseOutput(format string, args ...interface{}) {
	if c.verboseMode && !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	if c.logger != nil && c.verboseMode && c.logger.IsInfoEnabled() {
		c.logger.Info(format, args...)
	}
}

// DebugOutput prints debug information if debug mode is enabled
func (c *CLI) DebugOutput(format string, args ...interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("🐛 "+format+"\n", args...)
	}
	if c.logger != nil && c.logger.IsDebugEnabled() {
		c.logger.Debug(format, args...)
	}
}

// QuietOutput prints information only if not in quiet mode
func (c *CLI) QuietOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for QuietOutput - it's meant for user-facing messages
}

// ErrorOutput prints error information (always shown unless completely silent)
func (c *CLI) ErrorOutput(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	// Only log to logger in verbose or debug mode to avoid cluttering output
	if c.logger != nil && (c.verboseMode || c.debugMode) {
		c.logger.Error(format, args...)
	}
}

// WarningOutput prints warning information if not in quiet mode
func (c *CLI) WarningOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf("⚠️  "+format+"\n", args...)
	}
	if c.logger != nil {
		c.logger.Warn(format, args...)
	}
}

// SuccessOutput prints success information if not in quiet mode
func (c *CLI) SuccessOutput(format string, args ...interface{}) {
	if !c.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for SuccessOutput - it's meant for user-facing messages
}

// PerformanceOutput prints performance metrics if debug mode is enabled
func (c *CLI) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	if c.debugMode && !c.quietMode {
		fmt.Printf("⚡ %s completed in %v\n", operation, duration)
		if len(metrics) > 0 {
			fmt.Printf("📊 Performance Metrics:\n")
			for k, v := range metrics {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}
	if c.logger != nil {
		allMetrics := make(map[string]interface{})
		allMetrics["duration_ms"] = duration.Milliseconds()
		allMetrics["duration_human"] = duration.String()
		for k, v := range metrics {
			allMetrics[k] = v
		}
		c.logger.LogPerformanceMetrics(operation, allMetrics)
	}
}

// StartOperationWithOutput starts an operation with verbose output
func (c *CLI) StartOperationWithOutput(operation string, description string) *interfaces.OperationContext {
	c.VerboseOutput("%s", description)

	var ctx *interfaces.OperationContext
	if c.logger != nil && c.verboseMode {
		ctx = c.logger.StartOperation(operation, map[string]interface{}{
			"description": description,
		})
	}

	return ctx
}

// FinishOperationWithOutput completes an operation with verbose output
func (c *CLI) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	if ctx != nil && c.logger != nil && c.verboseMode {
		c.logger.FinishOperation(ctx, map[string]interface{}{
			"description": description,
		})
	}
	c.VerboseOutput("%s", description)
}

// FinishOperationWithError completes an operation with error output
func (c *CLI) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	if ctx != nil && c.logger != nil {
		c.logger.FinishOperationWithError(ctx, err, nil)
	}
	c.ErrorOutput("❌ %s failed: %v", operation, err)
}

// GetExitCode returns the current exit code
func (c *CLI) GetExitCode() int {
	return c.exitCode
}

// SetExitCode sets the exit code for the CLI
func (c *CLI) SetExitCode(code int) {
	c.exitCode = code
}

// GetVersionManager returns the version manager instance
func (c *CLI) GetVersionManager() interfaces.VersionManager {
	return c.versionManager
}
