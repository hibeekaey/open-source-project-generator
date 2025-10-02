// Package cli provides output management functionality for the CLI interface.
//
// This module handles all output formatting, color management, and different
// output modes (verbose, quiet, debug) for the CLI application.
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

// ColorManager handles color formatting for CLI output
type ColorManager struct {
	enabled bool
}

// NewColorManager creates a new ColorManager instance
func NewColorManager(enabled bool) *ColorManager {
	return &ColorManager{
		enabled: enabled,
	}
}

// Colorize applies color formatting to text if colors are enabled
func (cm *ColorManager) Colorize(color, text string) string {
	if !cm.enabled {
		return text
	}
	return color + text + ColorReset
}

// Success formats text with green color for success messages
func (cm *ColorManager) Success(text string) string {
	return cm.Colorize(ColorGreen+ColorBold, text)
}

// Info formats text with blue color for informational messages
func (cm *ColorManager) Info(text string) string {
	return cm.Colorize(ColorBlue, text)
}

// Warning formats text with yellow color for warning messages
func (cm *ColorManager) Warning(text string) string {
	return cm.Colorize(ColorYellow, text)
}

// Error formats text with red color for error messages
func (cm *ColorManager) Error(text string) string {
	return cm.Colorize(ColorRed+ColorBold, text)
}

// Highlight formats text with cyan color for highlighting
func (cm *ColorManager) Highlight(text string) string {
	return cm.Colorize(ColorCyan+ColorBold, text)
}

// Dim formats text with dim color for less important information
func (cm *ColorManager) Dim(text string) string {
	return cm.Colorize(ColorDim, text)
}

// OutputManager handles different types of CLI output based on mode settings
type OutputManager struct {
	verboseMode bool
	quietMode   bool
	debugMode   bool
	colorizer   *ColorManager
	logger      interfaces.Logger
}

// NewOutputManager creates a new OutputManager instance
func NewOutputManager(verboseMode, quietMode, debugMode bool, logger interfaces.Logger) *OutputManager {
	// Disable colors in quiet mode
	colorsEnabled := !quietMode

	return &OutputManager{
		verboseMode: verboseMode,
		quietMode:   quietMode,
		debugMode:   debugMode,
		colorizer:   NewColorManager(colorsEnabled),
		logger:      logger,
	}
}

// GetColorManager returns the color manager instance
func (om *OutputManager) GetColorManager() *ColorManager {
	return om.colorizer
}

// VerboseOutput prints verbose information if verbose mode is enabled
func (om *OutputManager) VerboseOutput(format string, args ...interface{}) {
	if om.verboseMode && !om.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	if om.logger != nil && om.verboseMode && om.logger.IsInfoEnabled() {
		om.logger.Info(format, args...)
	}
}

// DebugOutput prints debug information if debug mode is enabled
func (om *OutputManager) DebugOutput(format string, args ...interface{}) {
	if om.debugMode && !om.quietMode {
		fmt.Printf("🐛 "+format+"\n", args...)
	}
	if om.logger != nil && om.logger.IsDebugEnabled() {
		om.logger.Debug(format, args...)
	}
}

// QuietOutput prints information only if not in quiet mode
func (om *OutputManager) QuietOutput(format string, args ...interface{}) {
	if !om.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for QuietOutput - it's meant for user-facing messages
}

// ErrorOutput prints error information (always shown unless completely silent)
func (om *OutputManager) ErrorOutput(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	// Only log to logger in verbose or debug mode to avoid cluttering output
	if om.logger != nil && (om.verboseMode || om.debugMode) {
		om.logger.Error(format, args...)
	}
}

// WarningOutput prints warning information if not in quiet mode
func (om *OutputManager) WarningOutput(format string, args ...interface{}) {
	if !om.quietMode {
		fmt.Printf("⚠️  "+format+"\n", args...)
	}
	if om.logger != nil {
		om.logger.Warn(format, args...)
	}
}

// SuccessOutput prints success information if not in quiet mode
func (om *OutputManager) SuccessOutput(format string, args ...interface{}) {
	if !om.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for SuccessOutput - it's meant for user-facing messages
}

// PerformanceOutput prints performance metrics if debug mode is enabled
func (om *OutputManager) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	if om.debugMode && !om.quietMode {
		fmt.Printf("⚡ %s completed in %v\n", operation, duration)
		if len(metrics) > 0 {
			fmt.Printf("📊 Performance Metrics:\n")
			for k, v := range metrics {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}
	if om.logger != nil {
		allMetrics := make(map[string]interface{})
		allMetrics["duration_ms"] = duration.Milliseconds()
		allMetrics["duration_human"] = duration.String()
		for k, v := range metrics {
			allMetrics[k] = v
		}
		om.logger.LogPerformanceMetrics(operation, allMetrics)
	}
}

// StartOperationWithOutput starts an operation with verbose output
func (om *OutputManager) StartOperationWithOutput(operation string, description string) *interfaces.OperationContext {
	om.VerboseOutput("%s", description)

	var ctx *interfaces.OperationContext
	if om.logger != nil && om.verboseMode {
		ctx = om.logger.StartOperation(operation, map[string]interface{}{
			"description": description,
		})
	}

	return ctx
}

// FinishOperationWithOutput completes an operation with verbose output
func (om *OutputManager) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	if ctx != nil && om.logger != nil && om.verboseMode {
		om.logger.FinishOperation(ctx, map[string]interface{}{
			"description": description,
		})
	}
	om.VerboseOutput("%s", description)
}

// FinishOperationWithError completes an operation with error output
func (om *OutputManager) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	if ctx != nil && om.logger != nil {
		om.logger.FinishOperationWithError(ctx, err, nil)
	}
	om.ErrorOutput("❌ %s failed: %v", operation, err)
}

// SetVerboseMode updates the verbose mode setting
func (om *OutputManager) SetVerboseMode(enabled bool) {
	om.verboseMode = enabled
}

// SetQuietMode updates the quiet mode setting
func (om *OutputManager) SetQuietMode(enabled bool) {
	om.quietMode = enabled
	// Update color manager when quiet mode changes
	om.colorizer.enabled = !enabled
}

// SetDebugMode updates the debug mode setting
func (om *OutputManager) SetDebugMode(enabled bool) {
	om.debugMode = enabled
}

// IsVerboseMode returns whether verbose mode is enabled
func (om *OutputManager) IsVerboseMode() bool {
	return om.verboseMode
}

// IsQuietMode returns whether quiet mode is enabled
func (om *OutputManager) IsQuietMode() bool {
	return om.quietMode
}

// IsDebugMode returns whether debug mode is enabled
func (om *OutputManager) IsDebugMode() bool {
	return om.debugMode
}
