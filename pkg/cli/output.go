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

// OutputFormatter handles CLI output formatting and styling
type OutputFormatter struct {
	verboseMode bool
	quietMode   bool
	debugMode   bool
	logger      interfaces.Logger
}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter(verboseMode, quietMode, debugMode bool, logger interfaces.Logger) *OutputFormatter {
	return &OutputFormatter{
		verboseMode: verboseMode,
		quietMode:   quietMode,
		debugMode:   debugMode,
		logger:      logger,
	}
}

// Colorize applies color formatting to text
func (of *OutputFormatter) Colorize(color, text string) string {
	if of.quietMode {
		return text // No colors in quiet mode
	}
	return color + text + ColorReset
}

// Success formats text as success message
func (of *OutputFormatter) Success(text string) string {
	return of.Colorize(ColorGreen+ColorBold, text)
}

// Info formats text as info message
func (of *OutputFormatter) Info(text string) string {
	return of.Colorize(ColorBlue, text)
}

// Warning formats text as warning message
func (of *OutputFormatter) Warning(text string) string {
	return of.Colorize(ColorYellow, text)
}

// Error formats text as error message
func (of *OutputFormatter) Error(text string) string {
	return of.Colorize(ColorRed+ColorBold, text)
}

// Highlight formats text as highlighted message
func (of *OutputFormatter) Highlight(text string) string {
	return of.Colorize(ColorCyan+ColorBold, text)
}

// Dim formats text as dimmed message
func (of *OutputFormatter) Dim(text string) string {
	return of.Colorize(ColorDim, text)
}

// VerboseOutput prints verbose information if verbose mode is enabled
func (of *OutputFormatter) VerboseOutput(format string, args ...interface{}) {
	if of.verboseMode && !of.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	if of.logger != nil && of.verboseMode && of.logger.IsInfoEnabled() {
		of.logger.Info(format, args...)
	}
}

// DebugOutput prints debug information if debug mode is enabled
func (of *OutputFormatter) DebugOutput(format string, args ...interface{}) {
	if of.debugMode && !of.quietMode {
		fmt.Printf("🐛 "+format+"\n", args...)
	}
	if of.logger != nil && of.logger.IsDebugEnabled() {
		of.logger.Debug(format, args...)
	}
}

// QuietOutput prints information only if not in quiet mode
func (of *OutputFormatter) QuietOutput(format string, args ...interface{}) {
	if !of.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for QuietOutput - it's meant for user-facing messages
}

// ErrorOutput prints error information (always shown unless completely silent)
func (of *OutputFormatter) ErrorOutput(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	// Only log to logger in verbose or debug mode to avoid cluttering output
	if of.logger != nil && (of.verboseMode || of.debugMode) {
		of.logger.Error(format, args...)
	}
}

// WarningOutput prints warning information if not in quiet mode
func (of *OutputFormatter) WarningOutput(format string, args ...interface{}) {
	if !of.quietMode {
		fmt.Printf("⚠️  "+format+"\n", args...)
	}
	if of.logger != nil {
		of.logger.Warn(format, args...)
	}
}

// SuccessOutput prints success information if not in quiet mode
func (of *OutputFormatter) SuccessOutput(format string, args ...interface{}) {
	if !of.quietMode {
		fmt.Printf(format+"\n", args...)
	}
	// Don't log to structured logger for SuccessOutput - it's meant for user-facing messages
}

// PerformanceOutput prints performance metrics if debug mode is enabled
func (of *OutputFormatter) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	if of.debugMode && !of.quietMode {
		fmt.Printf("⚡ %s completed in %v\n", operation, duration)
		if len(metrics) > 0 {
			fmt.Printf("📊 Performance Metrics:\n")
			for k, v := range metrics {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}
	if of.logger != nil {
		allMetrics := make(map[string]interface{})
		allMetrics["duration_ms"] = duration.Milliseconds()
		allMetrics["duration_human"] = duration.String()
		for k, v := range metrics {
			allMetrics[k] = v
		}
		of.logger.LogPerformanceMetrics(operation, allMetrics)
	}
}

// StartOperationWithOutput starts an operation with verbose output
func (of *OutputFormatter) StartOperationWithOutput(operation string, description string) *interfaces.OperationContext {
	of.VerboseOutput("%s", description)

	if of.logger != nil {
		return of.logger.StartOperation(operation, map[string]interface{}{
			"description": description,
		})
	}

	return &interfaces.OperationContext{}
}

// FinishOperationWithOutput finishes an operation with verbose output
func (of *OutputFormatter) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	of.VerboseOutput("✓ %s", description)

	if of.logger != nil {
		of.logger.FinishOperation(ctx, map[string]interface{}{
			"description": description,
		})
	}
}

// FinishOperationWithError finishes an operation with error output
func (of *OutputFormatter) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	of.ErrorOutput("✗ %s failed: %v", operation, err)

	if of.logger != nil {
		of.logger.FinishOperationWithError(ctx, err, map[string]interface{}{
			"operation": operation,
		})
	}
}


