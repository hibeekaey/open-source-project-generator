// Package cli provides integration between CLI and enhanced error handling
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/errors"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ErrorIntegration provides enhanced error handling integration for CLI
type ErrorIntegration struct {
	enhancedHandler *errors.EnhancedErrorHandler
	cli             *CLI
	config          *errors.EnhancedErrorConfig
	logger          interfaces.Logger
}

// NewErrorIntegration creates a new error integration for CLI
func NewErrorIntegration(cli *CLI, logger interfaces.Logger) (*ErrorIntegration, error) {
	// Create enhanced error configuration based on CLI settings
	config := errors.DefaultEnhancedErrorConfig()

	// Configure based on CLI mode
	if cli != nil {
		config.VerboseMode = cli.verboseMode
		config.QuietMode = cli.quietMode
		config.InteractiveMode = !cli.helper.DetectNonInteractiveMode(cli.rootCmd)

		// Configure logging paths
		homeDir, _ := os.UserHomeDir()
		config.LogPath = fmt.Sprintf("%s/.generator/logs/cli.log", homeDir)
		config.ReportPath = fmt.Sprintf("%s/.generator/reports", homeDir)
	}

	// Create enhanced error handler
	enhancedHandler, err := errors.NewEnhancedErrorHandler(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create enhanced error handler: %w", err)
	}

	return &ErrorIntegration{
		enhancedHandler: enhancedHandler,
		cli:             cli,
		config:          config,
		logger:          logger,
	}, nil
}

// HandleError handles an error with enhanced error handling and CLI integration
func (ei *ErrorIntegration) HandleError(err error, operation string, contextData map[string]interface{}) int {
	if err == nil {
		return 0
	}

	ctx := context.Background()

	// Add CLI context information
	if ei.cli != nil {
		if contextData == nil {
			contextData = make(map[string]interface{})
		}

		contextData["cli_mode"] = ei.getCLIMode()
		contextData["verbose"] = ei.cli.verboseMode
		contextData["quiet"] = ei.cli.quietMode
		contextData["debug"] = ei.cli.debugMode

		if ei.cli.rootCmd != nil {
			contextData["command"] = ei.cli.rootCmd.Name()
			contextData["args"] = ei.cli.rootCmd.Flags().Args()
		}
	}

	// Handle error with enhanced handler
	result := ei.enhancedHandler.HandleError(ctx, err, operation, contextData)

	// Display CLI-specific error information
	ei.displayCLIError(result)

	// Return appropriate exit code
	return result.ExitCode
}

// getCLIMode determines the current CLI mode
func (ei *ErrorIntegration) getCLIMode() string {
	if ei.cli == nil {
		return "unknown"
	}

	if ei.cli.quietMode {
		return "quiet"
	}
	if ei.cli.debugMode {
		return "debug"
	}
	if ei.cli.verboseMode {
		return "verbose"
	}

	return "normal"
}

// displayCLIError displays error information in CLI-appropriate format
func (ei *ErrorIntegration) displayCLIError(result *errors.EnhancedErrorResult) {
	if result == nil || ei.config.QuietMode {
		return
	}

	// Use CLI's color manager if available
	if ei.cli != nil {
		ei.displayColorizedError(result)
	} else {
		ei.displayPlainError(result)
	}
}

// displayColorizedError displays error with CLI colors
func (ei *ErrorIntegration) displayColorizedError(result *errors.EnhancedErrorResult) {
	if result.Error == nil {
		return
	}

	// Display main error with appropriate color
	var icon string
	var colorFunc func(string) string

	switch result.Error.Severity {
	case errors.SeverityLow:
		icon = "ℹ️"
		colorFunc = ei.cli.Info
	case errors.SeverityMedium:
		icon = "⚠️"
		colorFunc = ei.cli.Warning
	case errors.SeverityHigh:
		icon = "🚨"
		colorFunc = ei.cli.Error
	case errors.SeverityCritical:
		icon = "🔥"
		colorFunc = ei.cli.Error
	default:
		icon = "❌"
		colorFunc = ei.cli.Error
	}

	fmt.Fprintf(os.Stderr, "\n%s %s\n", icon, colorFunc(result.Error.Message))

	// Show additional context in verbose mode
	if ei.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "   %s: %s (%s: %d)\n",
			ei.cli.Dim("Type"), result.Error.Type,
			ei.cli.Dim("Code"), result.Error.Code)
		fmt.Fprintf(os.Stderr, "   %s: %s\n",
			ei.cli.Dim("Severity"), result.Error.Severity)

		if result.Error.Context != nil && result.Error.Context.Operation != "" {
			fmt.Fprintf(os.Stderr, "   %s: %s\n",
				ei.cli.Dim("Operation"), result.Error.Context.Operation)
		}
	}

	// Display suggestions with CLI formatting
	if len(result.Suggestions) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "💡", ei.cli.Highlight("Suggestions:"))
		for i, suggestion := range result.Suggestions {
			if i < 5 { // Limit suggestions in CLI
				fmt.Fprintf(os.Stderr, "   %s. %s\n",
					ei.cli.Info(fmt.Sprintf("%d", i+1)), suggestion)
			}
		}

		if len(result.Suggestions) > 5 {
			fmt.Fprintf(os.Stderr, "   %s\n",
				ei.cli.Dim(fmt.Sprintf("... and %d more (use --verbose for all)", len(result.Suggestions)-5)))
		}
	}

	// Display recovery information
	if result.Recovery != nil {
		ei.displayRecoveryInfo(result.Recovery)
	}

	// Display user experience enhancements
	if result.UserExperience != nil {
		ei.displayUserExperienceEnhancements(result.UserExperience)
	}

	// Display performance information in debug mode
	if ei.config.EnablePerformanceTracking && result.Performance != nil && result.Performance.IsSlowOperation {
		ei.displayPerformanceWarning(result.Performance)
	}
}

// displayPlainError displays error without colors
func (ei *ErrorIntegration) displayPlainError(result *errors.EnhancedErrorResult) {
	if result.Error == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\nError: %s\n", result.Error.Message)

	if ei.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "Type: %s (Code: %d)\n", result.Error.Type, result.Error.Code)
		fmt.Fprintf(os.Stderr, "Severity: %s\n", result.Error.Severity)
	}

	if len(result.Suggestions) > 0 {
		fmt.Fprintf(os.Stderr, "\nSuggestions:\n")
		for i, suggestion := range result.Suggestions {
			if i < 5 {
				fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, suggestion)
			}
		}
	}
}

// displayRecoveryInfo displays recovery information with CLI formatting
func (ei *ErrorIntegration) displayRecoveryInfo(recovery *errors.RecoveryResult) {
	if recovery == nil {
		return
	}

	if recovery.Success {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "✅",
			ei.cli.Success(fmt.Sprintf("Recovery: %s", recovery.Message)))
	} else {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "❌",
			ei.cli.Error(fmt.Sprintf("Recovery failed: %s", recovery.Message)))

		if len(recovery.Suggestions) > 0 {
			fmt.Fprintf(os.Stderr, "\n%s %s\n", "🔧", ei.cli.Info("Recovery suggestions:"))
			for _, suggestion := range recovery.Suggestions {
				fmt.Fprintf(os.Stderr, "   • %s\n", suggestion)
			}
		}
	}
}

// displayUserExperienceEnhancements displays UX enhancements with CLI formatting
func (ei *ErrorIntegration) displayUserExperienceEnhancements(ux *errors.UserExperienceInfo) {
	if ux == nil {
		return
	}

	// Display contextual help
	if ux.ContextualHelp != "" {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "📖",
			ei.cli.Info(fmt.Sprintf("Help: %s", ux.ContextualHelp)))
	}

	// Display quick fixes
	if len(ux.QuickFixes) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "🔧", ei.cli.Highlight("Quick fixes:"))
		for i, fix := range ux.QuickFixes {
			if i < 3 { // Limit quick fixes in CLI
				fmt.Fprintf(os.Stderr, "   %s. %s",
					ei.cli.Info(fmt.Sprintf("%d", i+1)), fix.Description)

				if fix.Command != "" && ei.config.VerboseMode {
					fmt.Fprintf(os.Stderr, " (%s: %s)",
						ei.cli.Dim("Command"), ei.cli.Highlight(fix.Command))
				}
				fmt.Fprintf(os.Stderr, "\n")
			}
		}
	}

	// Display next steps
	if len(ux.NextSteps) > 0 {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "👉", ei.cli.Info("Next steps:"))
		for i, step := range ux.NextSteps {
			if i < 3 { // Limit next steps in CLI
				fmt.Fprintf(os.Stderr, "   %s. %s\n",
					ei.cli.Dim(fmt.Sprintf("%d", i+1)), step)
			}
		}
	}

	// Display documentation links in verbose mode
	if len(ux.RelatedDocs) > 0 && ei.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "\n%s %s\n", "📚", ei.cli.Info("Documentation:"))
		for i, doc := range ux.RelatedDocs {
			if i < 2 { // Limit docs in CLI
				fmt.Fprintf(os.Stderr, "   • %s: %s\n", doc.Title, ei.cli.Dim(doc.URL))
			}
		}
	}
}

// displayPerformanceWarning displays performance warnings
func (ei *ErrorIntegration) displayPerformanceWarning(perf *errors.PerformanceInfo) {
	if perf == nil || !perf.IsSlowOperation {
		return
	}

	fmt.Fprintf(os.Stderr, "\n%s %s\n", "⚡",
		ei.cli.Warning(fmt.Sprintf("Performance: Operation took %v (threshold: %v)",
			perf.Duration, perf.Threshold)))

	if len(perf.Suggestions) > 0 && ei.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "   %s:\n", ei.cli.Info("Performance suggestions"))
		for i, suggestion := range perf.Suggestions {
			if i < 2 { // Limit performance suggestions
				fmt.Fprintf(os.Stderr, "   • %s\n", suggestion)
			}
		}
	}
}

// HandleValidationError handles validation errors with enhanced context
func (ei *ErrorIntegration) HandleValidationError(field string, value interface{}, message string) int {
	ctx := context.Background()
	result := ei.enhancedHandler.NewEnhancedValidationError(ctx, message, field, value)
	ei.displayCLIError(result)
	return result.ExitCode
}

// HandleConfigurationError handles configuration errors with enhanced context
func (ei *ErrorIntegration) HandleConfigurationError(configPath string, err error, message string) int {
	ctx := context.Background()
	result := ei.enhancedHandler.NewEnhancedConfigurationError(ctx, message, configPath, err)
	ei.displayCLIError(result)
	return result.ExitCode
}

// HandleNetworkError handles network errors with enhanced context
func (ei *ErrorIntegration) HandleNetworkError(url string, err error, message string) int {
	ctx := context.Background()
	result := ei.enhancedHandler.NewEnhancedNetworkError(ctx, message, url, err)
	ei.displayCLIError(result)
	return result.ExitCode
}

// HandleFileSystemError handles filesystem errors with enhanced context
func (ei *ErrorIntegration) HandleFileSystemError(path string, operation string, err error, message string) int {
	ctx := context.Background()
	result := ei.enhancedHandler.NewEnhancedFileSystemError(ctx, message, path, operation, err)
	ei.displayCLIError(result)
	return result.ExitCode
}

// SetVerboseMode updates verbose mode for error handling
func (ei *ErrorIntegration) SetVerboseMode(verbose bool) {
	ei.config.VerboseMode = verbose
	if ei.enhancedHandler != nil {
		ei.enhancedHandler.SetVerboseMode(verbose)
	}
}

// SetQuietMode updates quiet mode for error handling
func (ei *ErrorIntegration) SetQuietMode(quiet bool) {
	ei.config.QuietMode = quiet
	if ei.enhancedHandler != nil {
		ei.enhancedHandler.SetQuietMode(quiet)
	}
}

// SetInteractiveMode updates interactive mode for error handling
func (ei *ErrorIntegration) SetInteractiveMode(interactive bool) {
	ei.config.InteractiveMode = interactive
	if ei.enhancedHandler != nil {
		ei.enhancedHandler.SetInteractiveMode(interactive)
	}
}

// GetStatistics returns comprehensive error handling statistics
func (ei *ErrorIntegration) GetStatistics() *errors.EnhancedErrorStatistics {
	if ei.enhancedHandler == nil {
		return nil
	}
	return ei.enhancedHandler.GetStatistics()
}

// GenerateErrorReport generates a comprehensive error report
func (ei *ErrorIntegration) GenerateErrorReport() string {
	if ei.enhancedHandler == nil {
		return "Error handling not initialized"
	}

	stats := ei.enhancedHandler.GetStatistics()
	if stats == nil {
		return "No error statistics available"
	}

	var report strings.Builder

	report.WriteString("Enhanced Error Handling Report\n")
	report.WriteString("==============================\n\n")

	// Basic statistics
	report.WriteString(fmt.Sprintf("Total Errors: %d\n", stats.TotalErrors))
	report.WriteString(fmt.Sprintf("Recovery Rate: %.1f%%\n", stats.RecoveryRate))

	// Error categories
	if len(stats.ErrorsByCategory) > 0 {
		report.WriteString("\nError Categories:\n")
		for category, count := range stats.ErrorsByCategory {
			percentage := float64(count) / float64(stats.TotalErrors) * 100
			report.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", category, count, percentage))
		}
	}

	// Performance statistics
	if stats.PerformanceStats != nil {
		report.WriteString("\nPerformance:\n")
		report.WriteString(fmt.Sprintf("  Total Operations: %d\n", stats.PerformanceStats.TotalOperations))
		report.WriteString(fmt.Sprintf("  Slow Operations: %d\n", stats.PerformanceStats.SlowOperations))
		report.WriteString(fmt.Sprintf("  Average Duration: %v\n", stats.PerformanceStats.AverageDuration))
	}

	// User experience statistics
	if stats.UserExperienceStats != nil {
		report.WriteString("\nUser Experience:\n")
		report.WriteString(fmt.Sprintf("  Quick Fixes Provided: %d\n", stats.UserExperienceStats.QuickFixesProvided))
		report.WriteString(fmt.Sprintf("  Quick Fix Success Rate: %.1f%%\n", stats.UserExperienceStats.QuickFixSuccessRate))
	}

	return report.String()
}

// Close closes the error integration and releases resources
func (ei *ErrorIntegration) Close() error {
	if ei.enhancedHandler != nil {
		return ei.enhancedHandler.Close()
	}
	return nil
}

// Helper method to integrate with existing CLI error handling
func (c *CLI) handleError(err error, command string, args []string) int {
	// Create error integration if not exists
	if c.errorIntegration == nil {
		errorIntegration, integrationErr := NewErrorIntegration(c, c.logger)
		if integrationErr != nil {
			// Fallback to basic error handling
			c.ErrorOutput("Error: %v", err)
			return 1
		}
		c.errorIntegration = errorIntegration
	}

	// Handle error with enhanced error handling
	context := map[string]interface{}{
		"command": command,
		"args":    args,
	}

	return c.errorIntegration.HandleError(err, command, context)
}

// Add error integration field to CLI struct (this would be added to the CLI struct definition)
// errorIntegration *ErrorIntegration
