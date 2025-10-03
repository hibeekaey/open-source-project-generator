// Package commands provides command-specific implementations for the CLI.
//
// This package contains individual command modules that handle specific CLI operations,
// separating command logic from the main CLI orchestration.
package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// ValidateCLI defines the CLI methods needed by ValidateCommand.
type ValidateCLI interface {
	// Validation methods
	ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error)

	// Output methods
	QuietOutput(format string, args ...interface{})
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	Error(text string) string
	Warning(text string) string
	Info(text string) string
	IsQuietMode() bool

	// Error handling
	CreateValidationError(message string, details map[string]interface{}) error
}

// ValidateCommand handles the validate command functionality.
//
// The ValidateCommand struct provides comprehensive project validation capabilities including:
//   - Project structure validation
//   - Configuration file validation
//   - Dependency validation
//   - Security validation
//   - Auto-fix capabilities for common issues
//   - Detailed validation reporting
//
// It supports multiple output formats and can generate detailed reports for CI/CD integration.
type ValidateCommand struct {
	cli ValidateCLI
}

// NewValidateCommand creates a new ValidateCommand instance.
//
// Parameters:
//   - cli: The CLI interface for accessing validation functionality
//
// Returns:
//   - *ValidateCommand: New ValidateCommand instance ready for use
func NewValidateCommand(cli ValidateCLI) *ValidateCommand {
	return &ValidateCommand{
		cli: cli,
	}
}

// Execute handles the validate command execution with comprehensive validation logic.
//
// This method performs the following operations:
//   - Parses command flags and arguments
//   - Validates the target project path
//   - Executes validation with specified options
//   - Formats and outputs validation results
//   - Generates reports if requested
//   - Returns appropriate exit codes based on validation results
//
// Parameters:
//   - cmd: The cobra command instance
//   - args: Command line arguments
//
// Returns:
//   - error: Any error encountered during validation
func (vc *ValidateCommand) Execute(cmd *cobra.Command, args []string) error {
	// Get path from args or use current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Parse command flags
	options, err := vc.parseFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	// Execute validation
	result, err := vc.cli.ValidateProject(path, *options)
	if err != nil {
		return vc.handleValidationError(err)
	}

	// Output results
	if err := vc.outputResults(cmd, result, path, options); err != nil {
		return fmt.Errorf("failed to output validation results: %w", err)
	}

	// Generate report if requested
	if options.Report && options.OutputFile != "" {
		if err := vc.generateReport(result, options.ReportFormat, options.OutputFile); err != nil {
			return fmt.Errorf("failed to generate validation report: %w", err)
		}
		vc.cli.QuietOutput("📄 Validation report saved: %s", vc.cli.Info(options.OutputFile))
	}

	// Return appropriate exit code
	return vc.handleValidationResult(result, path)
}

// parseFlags extracts and validates command flags into ValidationOptions.
//
// This method handles all validate command flags including:
//   - Basic validation options (fix, report, rules)
//   - Output formatting options
//   - Advanced validation options (strict mode, exclusions)
//   - Global flags (verbose, non-interactive)
//
// Parameters:
//   - cmd: The cobra command instance
//
// Returns:
//   - *interfaces.ValidationOptions: Parsed validation options
//   - error: Any error encountered during flag parsing
func (vc *ValidateCommand) parseFlags(cmd *cobra.Command) (*interfaces.ValidationOptions, error) {
	// Basic validation flags
	fix, _ := cmd.Flags().GetBool("fix")
	report, _ := cmd.Flags().GetBool("report")
	reportFormat, _ := cmd.Flags().GetString("report-format")
	rules, _ := cmd.Flags().GetStringSlice("rules")
	ignoreWarnings, _ := cmd.Flags().GetBool("ignore-warnings")
	outputFile, _ := cmd.Flags().GetString("output-file")
	output, _ := cmd.Flags().GetString("output")

	// Use --output if provided, otherwise use --output-file
	if output != "" {
		outputFile = output
	}

	// Advanced validation flags
	strict, _ := cmd.Flags().GetBool("strict")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")
	excludeRules, _ := cmd.Flags().GetStringSlice("exclude-rules")
	showFixes, _ := cmd.Flags().GetBool("show-fixes")

	// Global flags
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Log advanced options for debugging
	if strict {
		vc.cli.DebugOutput("Using strict validation mode")
	}
	if len(excludeRules) > 0 {
		vc.cli.DebugOutput("Excluding rules: %v", excludeRules)
	}

	// Store advanced flags for future use
	_ = summaryOnly
	_ = showFixes

	return &interfaces.ValidationOptions{
		Verbose:        verbose,
		Fix:            fix,
		Report:         report,
		ReportFormat:   reportFormat,
		Rules:          rules,
		IgnoreWarnings: ignoreWarnings,
		OutputFile:     outputFile,
	}, nil
}

// outputResults formats and displays validation results based on output mode and format.
//
// This method handles different output scenarios:
//   - Machine-readable output for automation (JSON/YAML)
//   - Human-readable output with colors and formatting
//   - Summary-only output for quick validation checks
//   - Detailed output with issue descriptions and fix suggestions
//
// Parameters:
//   - cmd: The cobra command instance
//   - result: The validation result to output
//   - path: The validated project path
//   - options: The validation options used
//
// Returns:
//   - error: Any error encountered during output formatting
func (vc *ValidateCommand) outputResults(cmd *cobra.Command, result *interfaces.ValidationResult, path string, options *interfaces.ValidationOptions) error {
	// Get output format flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")
	showFixes, _ := cmd.Flags().GetBool("show-fixes")

	// Auto-detect non-interactive mode if not explicitly set
	// TODO: Implement auto-detection of non-interactive mode
	// if !nonInteractive {
	//     nonInteractive = vc.cli.IsNonInteractiveMode()
	// }

	// Handle machine-readable output for automation
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return vc.outputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	return vc.outputHumanReadable(result, path, summaryOnly, showFixes, options.IgnoreWarnings)
}

// outputMachineReadable outputs validation results in machine-readable format.
//
// This method supports JSON and YAML output formats for CI/CD integration
// and automation workflows.
//
// Parameters:
//   - result: The validation result to output
//   - format: The output format (json or yaml)
//
// Returns:
//   - error: Any error encountered during output formatting
func (vc *ValidateCommand) outputMachineReadable(result *interfaces.ValidationResult, format string) error {
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	case "yaml":
		// YAML output would require a YAML library
		// For now, fall back to JSON
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// outputHumanReadable outputs validation results in human-readable format.
//
// This method provides colorized, formatted output with:
//   - Validation status and summary
//   - Detailed issue listings with file locations
//   - Warning listings (unless ignored)
//   - Fix suggestions (if requested)
//   - Appropriate emoji and color coding
//
// Parameters:
//   - result: The validation result to output
//   - path: The validated project path
//   - summaryOnly: Whether to show only summary information
//   - showFixes: Whether to show fix suggestions
//   - ignoreWarnings: Whether to ignore warnings in output
//
// Returns:
//   - error: Any error encountered during output formatting
func (vc *ValidateCommand) outputHumanReadable(result *interfaces.ValidationResult, path string, summaryOnly, showFixes, ignoreWarnings bool) error {
	if vc.cli.IsQuietMode() {
		return nil
	}

	// Output validation summary
	vc.cli.QuietOutput("🔍 Validation completed for: %s", path)
	if result.Valid {
		vc.cli.QuietOutput("✅ Project looks good!")
	} else {
		vc.cli.QuietOutput("%s %s",
			vc.cli.Error("❌ Found some issues that need attention."),
			vc.cli.Info("See details below"))
	}

	// Output statistics
	vc.cli.QuietOutput("📊 Issues: %s", vc.cli.Error(fmt.Sprintf("%d", len(result.Issues))))
	vc.cli.QuietOutput("⚠️  Warnings: %s", vc.cli.Warning(fmt.Sprintf("%d", len(result.Warnings))))

	// Output detailed issues if not summary-only
	if len(result.Issues) > 0 && !summaryOnly {
		vc.cli.QuietOutput("\n🚨 Issues that need fixing:")
		for _, issue := range result.Issues {
			vc.cli.QuietOutput("  - %s: %s", issue.Severity, issue.Message)
			if issue.File != "" {
				vc.cli.VerboseOutput("    File: %s:%d:%d", issue.File, issue.Line, issue.Column)
			}
		}
	}

	// Output warnings if not ignored and not summary-only
	if len(result.Warnings) > 0 && !ignoreWarnings && !summaryOnly {
		vc.cli.QuietOutput("\n%s", vc.cli.Warning("⚠️  Things to consider:"))
		for _, warning := range result.Warnings {
			vc.cli.QuietOutput("  - %s: %s", warning.Severity, warning.Message)
			if warning.File != "" {
				vc.cli.VerboseOutput("    File: %s:%d:%d", warning.File, warning.Line, warning.Column)
			}
		}
	}

	// Output fix suggestions if requested
	if showFixes && len(result.FixSuggestions) > 0 {
		vc.cli.QuietOutput("\nSuggested fixes:")
		for _, suggestion := range result.FixSuggestions {
			vc.cli.QuietOutput("  - %s", suggestion.Description)
			if suggestion.AutoFixable {
				vc.cli.QuietOutput("    (Auto-fixable with --fix flag)")
			}
		}
	}

	return nil
}

// generateReport creates a validation report in the specified format.
//
// This method supports multiple report formats:
//   - JSON: Machine-readable structured data
//   - Text: Human-readable plain text report
//   - HTML: Formatted HTML report (future implementation)
//   - Markdown: Markdown-formatted report (future implementation)
//
// Parameters:
//   - result: The validation result to include in the report
//   - format: The report format to generate
//   - outputFile: The file path to save the report
//
// Returns:
//   - error: Any error encountered during report generation
func (vc *ValidateCommand) generateReport(result *interfaces.ValidationResult, format, outputFile string) error {
	var content []byte
	var err error

	switch format {
	case "json":
		content, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON report: %w", err)
		}
	case "text":
		content = vc.generateTextReport(result)
	case "html":
		// HTML report generation would be implemented here
		content = vc.generateTextReport(result) // Fallback to text for now
	case "markdown":
		// Markdown report generation would be implemented here
		content = vc.generateTextReport(result) // Fallback to text for now
	default:
		content = vc.generateTextReport(result)
	}

	// Write report to file
	if err := os.WriteFile(outputFile, content, 0600); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

// generateTextReport creates a plain text validation report.
//
// This method generates a comprehensive text report including:
//   - Validation status and summary
//   - Issue counts and statistics
//   - Detailed issue listings
//   - Fix suggestions
//
// Parameters:
//   - result: The validation result to format
//
// Returns:
//   - []byte: The formatted text report content
func (vc *ValidateCommand) generateTextReport(result *interfaces.ValidationResult) []byte {
	status := "✅ Looks good!"
	if !result.Valid {
		status = "❌ Needs attention"
	}

	report := fmt.Sprintf(`🔍 Validation Report
===================

Status: %s
📊 Issues: %d
⚠️  Warnings: %d

`, status, len(result.Issues), len(result.Warnings))

	// Add detailed issues
	if len(result.Issues) > 0 {
		report += "Issues:\n"
		for i, issue := range result.Issues {
			report += fmt.Sprintf("%d. %s: %s\n", i+1, issue.Severity, issue.Message)
			if issue.File != "" {
				report += fmt.Sprintf("   File: %s:%d:%d\n", issue.File, issue.Line, issue.Column)
			}
			if issue.Rule != "" {
				report += fmt.Sprintf("   Rule: %s\n", issue.Rule)
			}
			report += "\n"
		}
	}

	// Add warnings
	if len(result.Warnings) > 0 {
		report += "Warnings:\n"
		for i, warning := range result.Warnings {
			report += fmt.Sprintf("%d. %s: %s\n", i+1, warning.Severity, warning.Message)
			if warning.File != "" {
				report += fmt.Sprintf("   File: %s:%d:%d\n", warning.File, warning.Line, warning.Column)
			}
			report += "\n"
		}
	}

	// Add fix suggestions
	if len(result.FixSuggestions) > 0 {
		report += "Fix Suggestions:\n"
		for i, suggestion := range result.FixSuggestions {
			report += fmt.Sprintf("%d. %s\n", i+1, suggestion.Description)
			if suggestion.AutoFixable {
				report += "   (Auto-fixable with --fix flag)\n"
			}
			report += "\n"
		}
	}

	return []byte(report)
}

// handleValidationError formats validation errors with helpful context.
//
// This method provides user-friendly error messages with suggestions
// for common validation issues.
//
// Parameters:
//   - err: The validation error to handle
//
// Returns:
//   - error: Formatted error with helpful context
func (vc *ValidateCommand) handleValidationError(err error) error {
	return fmt.Errorf("🚫 %s %s",
		vc.cli.Error("Project validation encountered an issue."),
		vc.cli.Info("Try running with --verbose to see more details: "+err.Error()))
}

// handleValidationResult processes validation results and returns appropriate exit codes.
//
// This method determines the appropriate exit code based on validation results:
//   - Success (0): No issues found
//   - Failure (1): Issues found that need attention
//   - Error (2): Validation process failed
//
// Parameters:
//   - result: The validation result to process
//   - path: The validated project path
//
// Returns:
//   - error: Validation error if issues were found, nil if validation passed
func (vc *ValidateCommand) handleValidationResult(result *interfaces.ValidationResult, path string) error {
	if result.Valid {
		return nil
	}

	// Create detailed error information
	details := map[string]interface{}{
		"issues_count":   len(result.Issues),
		"warnings_count": len(result.Warnings),
		"path":           path,
	}

	var message string
	if len(result.Issues) > 0 {
		message = fmt.Sprintf("🚫 Found %s that need your attention",
			vc.cli.Error(fmt.Sprintf("%d validation issues", len(result.Issues))))
	} else if len(result.Warnings) > 0 {
		message = fmt.Sprintf("⚠️  Found %s that should be addressed",
			vc.cli.Warning(fmt.Sprintf("%d warnings", len(result.Warnings))))
	} else {
		message = fmt.Sprintf("🚫 %s %s",
			vc.cli.Error("Validation failed."),
			vc.cli.Info("Please check your project structure and configuration"))
	}

	return vc.cli.CreateValidationError(message, details)
}
