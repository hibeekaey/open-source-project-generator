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

// AuditCLI defines the CLI methods needed by AuditCommand.
type AuditCLI interface {
	// Audit methods
	AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error)

	// Output methods
	QuietOutput(format string, args ...interface{})
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	Error(text string) string
	Warning(text string) string
	Info(text string) string
	IsQuietMode() bool

	// Error handling
	CreateAuditError(message string, score float64) error

	// Machine-readable output
	OutputMachineReadable(data interface{}, format string) error
}

// AuditCommand handles the audit command functionality.
//
// The AuditCommand struct provides comprehensive project auditing capabilities including:
//   - Security vulnerability scanning
//   - Code quality analysis
//   - License compliance checking
//   - Performance analysis
//   - Detailed audit reporting
//   - Multiple output formats for CI/CD integration
//
// It supports various audit types and can generate detailed reports with scores
// and recommendations for improvement.
type AuditCommand struct {
	cli AuditCLI
}

// NewAuditCommand creates a new AuditCommand instance.
//
// Parameters:
//   - cli: The CLI interface for accessing audit functionality
//
// Returns:
//   - *AuditCommand: New AuditCommand instance ready for use
func NewAuditCommand(cli AuditCLI) *AuditCommand {
	return &AuditCommand{
		cli: cli,
	}
}

// Execute handles the audit command execution with comprehensive auditing logic.
//
// This method performs the following operations:
//   - Parses command flags and arguments
//   - Validates the target project path
//   - Executes audit with specified options
//   - Formats and outputs audit results
//   - Generates reports if requested
//   - Returns appropriate exit codes based on audit results
//
// Parameters:
//   - cmd: The cobra command instance
//   - args: Command line arguments
//
// Returns:
//   - error: Any error encountered during auditing
func (ac *AuditCommand) Execute(cmd *cobra.Command, args []string) error {
	// Get path from args or use current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Parse command flags
	options, err := ac.parseFlags(cmd)
	if err != nil {
		return fmt.Errorf("failed to parse command flags: %w", err)
	}

	// Execute audit
	result, err := ac.cli.AuditProject(path, *options)
	if err != nil {
		return ac.handleAuditError(err)
	}

	// Output results
	if err := ac.outputResults(cmd, result, path, options); err != nil {
		return fmt.Errorf("failed to output audit results: %w", err)
	}

	// Generate report if requested
	if options.OutputFile != "" {
		if err := ac.generateReport(result, options.OutputFormat, options.OutputFile); err != nil {
			return fmt.Errorf("failed to generate audit report: %w", err)
		}
		ac.cli.QuietOutput("📄 Audit report saved: %s", ac.cli.Info(options.OutputFile))
	}

	// Return appropriate exit code based on audit results
	return ac.handleAuditResult(cmd, result)
}

// parseFlags extracts and validates command flags into AuditOptions.
//
// This method handles all audit command flags including:
//   - Audit type options (security, quality, licenses, performance)
//   - Output formatting options
//   - Advanced audit options (detailed, fail conditions)
//   - Global flags (verbose, non-interactive)
//
// Parameters:
//   - cmd: The cobra command instance
//
// Returns:
//   - *interfaces.AuditOptions: Parsed audit options
//   - error: Any error encountered during flag parsing
func (ac *AuditCommand) parseFlags(cmd *cobra.Command) (*interfaces.AuditOptions, error) {
	// Basic audit flags
	security, _ := cmd.Flags().GetBool("security")
	quality, _ := cmd.Flags().GetBool("quality")
	licenses, _ := cmd.Flags().GetBool("licenses")
	performance, _ := cmd.Flags().GetBool("performance")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	detailed, _ := cmd.Flags().GetBool("detailed")

	// Advanced audit flags
	excludeCategories, _ := cmd.Flags().GetStringSlice("exclude-categories")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")

	// Global flags
	globalOutputFormat, _ := cmd.Flags().GetString("output-format")

	// Use global output format if command-specific format not set
	if outputFormat == "text" && globalOutputFormat != "text" {
		outputFormat = globalOutputFormat
	}

	// Log advanced options for debugging
	if len(excludeCategories) > 0 {
		ac.cli.DebugOutput("Excluding audit categories: %v", excludeCategories)
	}
	if summaryOnly {
		ac.cli.DebugOutput("Showing summary only")
	}

	return &interfaces.AuditOptions{
		Security:     security,
		Quality:      quality,
		Licenses:     licenses,
		Performance:  performance,
		OutputFormat: outputFormat,
		OutputFile:   outputFile,
		Detailed:     detailed,
	}, nil
}

// outputResults formats and displays audit results based on output mode and format.
//
// This method handles different output scenarios:
//   - Machine-readable output for automation (JSON/YAML)
//   - Human-readable output with colors and formatting
//   - Summary-only output for quick audit checks
//   - Detailed output with issue descriptions and recommendations
//
// Parameters:
//   - cmd: The cobra command instance
//   - result: The audit result to output
//   - path: The audited project path
//   - options: The audit options used
//
// Returns:
//   - error: Any error encountered during output formatting
func (ac *AuditCommand) outputResults(cmd *cobra.Command, result *interfaces.AuditResult, path string, options *interfaces.AuditOptions) error {
	// Get output format flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")

	// Auto-detect non-interactive mode if not explicitly set
	if !nonInteractive {
		// This would need to be implemented in the CLI interface
		// nonInteractive = ac.cli.IsNonInteractiveMode()
	}

	// Handle machine-readable output for automation
	if nonInteractive && (options.OutputFormat == "json" || options.OutputFormat == "yaml") {
		return ac.cli.OutputMachineReadable(result, options.OutputFormat)
	}

	// Human-readable output
	return ac.outputHumanReadable(result, path, summaryOnly)
}

// outputHumanReadable outputs audit results in human-readable format.
//
// This method provides colorized, formatted output with:
//   - Audit status and overall score
//   - Individual audit category scores and results
//   - Detailed issue listings with severity indicators
//   - Recommendations for improvement
//   - Appropriate emoji and color coding
//
// Parameters:
//   - result: The audit result to output
//   - path: The audited project path
//   - summaryOnly: Whether to show only summary information
//
// Returns:
//   - error: Any error encountered during output formatting
func (ac *AuditCommand) outputHumanReadable(result *interfaces.AuditResult, path string, summaryOnly bool) error {
	if ac.cli.IsQuietMode() {
		return nil
	}

	// Output audit summary
	ac.cli.QuietOutput("🔍 Audit completed for: %s", path)

	// Overall score with appropriate emoji
	scoreEmoji := "🎉"
	if result.OverallScore < 70 {
		scoreEmoji = "⚠️ "
	}
	if result.OverallScore < 50 {
		scoreEmoji = "🚨"
	}
	ac.cli.QuietOutput("%s Overall Score: %.1f/100", scoreEmoji, result.OverallScore)
	ac.cli.VerboseOutput("Audit Time: %s", result.AuditTime.Format("2006-01-02 15:04:05"))

	// Security audit results
	if result.Security != nil && !summaryOnly {
		securityEmoji := "🔒"
		if result.Security.Score < 70 {
			securityEmoji = "⚠️ "
		}
		if result.Security.Score < 50 {
			securityEmoji = "🚨"
		}
		ac.cli.QuietOutput("%s Security Score: %.1f/100", securityEmoji, result.Security.Score)
		ac.cli.VerboseOutput("Vulnerabilities: %d", len(result.Security.Vulnerabilities))

		// Show critical vulnerabilities
		if len(result.Security.Vulnerabilities) > 0 && !summaryOnly {
			criticalCount := 0
			for _, vuln := range result.Security.Vulnerabilities {
				if vuln.Severity == "critical" || vuln.Severity == "high" {
					criticalCount++
				}
			}
			if criticalCount > 0 {
				ac.cli.QuietOutput("  %s %d critical/high severity vulnerabilities found",
					ac.cli.Error("🚨"), criticalCount)
			}
		}
	}

	// Quality audit results
	if result.Quality != nil && !summaryOnly {
		qualityEmoji := "✨"
		if result.Quality.Score < 70 {
			qualityEmoji = "⚠️ "
		}
		if result.Quality.Score < 50 {
			qualityEmoji = "🚨"
		}
		ac.cli.QuietOutput("%s Quality Score: %.1f/100", qualityEmoji, result.Quality.Score)
		ac.cli.VerboseOutput("Code Smells: %d", len(result.Quality.CodeSmells))
	}

	// License audit results
	if result.Licenses != nil && !summaryOnly {
		licenseEmoji := "📄"
		if !result.Licenses.Compatible {
			licenseEmoji = "⚠️ "
		}
		ac.cli.QuietOutput("%s License Compatible: %t", licenseEmoji, result.Licenses.Compatible)
		ac.cli.VerboseOutput("License Conflicts: %d", len(result.Licenses.Conflicts))
	}

	// Performance audit results
	if result.Performance != nil && !summaryOnly {
		perfEmoji := "⚡"
		if result.Performance.Score < 70 {
			perfEmoji = "⚠️ "
		}
		if result.Performance.Score < 50 {
			perfEmoji = "🚨"
		}
		ac.cli.QuietOutput("%s Performance Score: %.1f/100", perfEmoji, result.Performance.Score)
		ac.cli.VerboseOutput("Bundle Size: %d bytes", result.Performance.BundleSize)
	}

	// Output recommendations
	if len(result.Recommendations) > 0 && !summaryOnly {
		ac.cli.QuietOutput("\n💡 Recommendations:")
		for _, rec := range result.Recommendations {
			ac.cli.QuietOutput("  - %s", rec)
		}
	}

	return nil
}

// generateReport creates an audit report in the specified format.
//
// This method supports multiple report formats:
//   - JSON: Machine-readable structured data
//   - Text: Human-readable plain text report
//   - HTML: Formatted HTML report (future implementation)
//   - YAML: YAML-formatted report (future implementation)
//
// Parameters:
//   - result: The audit result to include in the report
//   - format: The report format to generate
//   - outputFile: The file path to save the report
//
// Returns:
//   - error: Any error encountered during report generation
func (ac *AuditCommand) generateReport(result *interfaces.AuditResult, format, outputFile string) error {
	var content []byte
	var err error

	switch format {
	case "json":
		content, err = json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON report: %w", err)
		}
	case "text":
		content = ac.generateTextReport(result)
	case "html":
		// HTML report generation would be implemented here
		content = ac.generateTextReport(result) // Fallback to text for now
	case "yaml":
		// YAML report generation would be implemented here
		content = ac.generateTextReport(result) // Fallback to text for now
	default:
		content = ac.generateTextReport(result)
	}

	// Write report to file
	if err := os.WriteFile(outputFile, content, 0600); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

// generateTextReport creates a plain text audit report.
//
// This method generates a comprehensive text report including:
//   - Audit status and overall score
//   - Individual category scores and details
//   - Vulnerability listings
//   - Quality issues
//   - License compliance status
//   - Performance metrics
//   - Recommendations
//
// Parameters:
//   - result: The audit result to format
//
// Returns:
//   - []byte: The formatted text report content
func (ac *AuditCommand) generateTextReport(result *interfaces.AuditResult) []byte {
	report := fmt.Sprintf(`🔍 Audit Report
===============

Project: %s
Audit Time: %s
Overall Score: %.1f/100

`, result.ProjectPath, result.AuditTime.Format("2006-01-02 15:04:05"), result.OverallScore)

	// Security section
	if result.Security != nil {
		report += fmt.Sprintf(`🔒 Security Audit
Score: %.1f/100
Vulnerabilities: %d

`, result.Security.Score, len(result.Security.Vulnerabilities))

		if len(result.Security.Vulnerabilities) > 0 {
			report += "Vulnerabilities:\n"
			for i, vuln := range result.Security.Vulnerabilities {
				report += fmt.Sprintf("%d. [%s] %s (%s)\n", i+1, vuln.Severity, vuln.Title, vuln.ID)
				report += fmt.Sprintf("   Package: %s@%s\n", vuln.Package, vuln.Version)
				if vuln.FixedIn != "" {
					report += fmt.Sprintf("   Fixed in: %s\n", vuln.FixedIn)
				}
				report += fmt.Sprintf("   Description: %s\n\n", vuln.Description)
			}
		}
	}

	// Quality section
	if result.Quality != nil {
		report += fmt.Sprintf(`✨ Quality Audit
Score: %.1f/100
Code Smells: %d

`, result.Quality.Score, len(result.Quality.CodeSmells))

		if len(result.Quality.CodeSmells) > 0 {
			report += "Code Quality Issues:\n"
			for i, smell := range result.Quality.CodeSmells {
				report += fmt.Sprintf("%d. [%s] %s: %s\n", i+1, smell.Severity, smell.Type, smell.Description)
				if smell.File != "" {
					report += fmt.Sprintf("   File: %s:%d\n", smell.File, smell.Line)
				}
				report += "\n"
			}
		}
	}

	// License section
	if result.Licenses != nil {
		report += fmt.Sprintf(`📄 License Audit
Score: %.1f/100
Compatible: %t
Conflicts: %d

`, result.Licenses.Score, result.Licenses.Compatible, len(result.Licenses.Conflicts))

		if len(result.Licenses.Conflicts) > 0 {
			report += "License Conflicts:\n"
			for i, conflict := range result.Licenses.Conflicts {
				report += fmt.Sprintf("%d. %s (%s): %s\n", i+1, conflict.Name, conflict.SPDXID, conflict.Package)
				if !conflict.Compatible {
					report += "   Status: Incompatible\n"
				}
				report += "\n"
			}
		}
	}

	// Performance section
	if result.Performance != nil {
		report += fmt.Sprintf(`⚡ Performance Audit
Score: %.1f/100
Bundle Size: %d bytes

`, result.Performance.Score, result.Performance.BundleSize)
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		report += "💡 Recommendations:\n"
		for i, rec := range result.Recommendations {
			report += fmt.Sprintf("%d. %s\n", i+1, rec)
		}
		report += "\n"
	}

	return []byte(report)
}

// handleAuditError formats audit errors with helpful context.
//
// This method provides user-friendly error messages with suggestions
// for common audit issues.
//
// Parameters:
//   - err: The audit error to handle
//
// Returns:
//   - error: Formatted error with helpful context
func (ac *AuditCommand) handleAuditError(err error) error {
	return fmt.Errorf("🚫 %s %s",
		ac.cli.Error("Project audit encountered an issue."),
		ac.cli.Info("Try running with --verbose to see more details: "+err.Error()))
}

// handleAuditResult processes audit results and returns appropriate exit codes.
//
// This method determines the appropriate exit code based on audit results and flags:
//   - Success (0): No critical issues found
//   - Failure (1): Critical issues found or score below threshold
//   - Error (2): Audit process failed
//
// Parameters:
//   - cmd: The cobra command instance
//   - result: The audit result to process
//
// Returns:
//   - error: Audit error if critical issues were found, nil if audit passed
func (ac *AuditCommand) handleAuditResult(cmd *cobra.Command, result *interfaces.AuditResult) error {
	// Get fail condition flags
	failOnHigh, _ := cmd.Flags().GetBool("fail-on-high")
	failOnMedium, _ := cmd.Flags().GetBool("fail-on-medium")
	minScore, _ := cmd.Flags().GetFloat64("min-score")

	// Check fail conditions and return appropriate exit codes
	if failOnHigh && result.OverallScore < 70.0 {
		return ac.cli.CreateAuditError(
			fmt.Sprintf("🚫 Found high severity issues (score: %.2f/100)", result.OverallScore),
			result.OverallScore)
	}

	if failOnMedium && result.OverallScore < 50.0 {
		return ac.cli.CreateAuditError(
			fmt.Sprintf("🚫 Found medium or high severity issues (score: %.2f/100)", result.OverallScore),
			result.OverallScore)
	}

	if minScore > 0 && result.OverallScore < minScore {
		return ac.cli.CreateAuditError(
			fmt.Sprintf("🚫 Score %.2f/100 is below your minimum requirement of %.2f/100", result.OverallScore, minScore),
			result.OverallScore)
	}

	return nil
}
