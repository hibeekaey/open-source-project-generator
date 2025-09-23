package cli

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// setupAuditCommand sets up the audit command
func (c *CLI) setupAuditCommand() {
	auditCmd := &cobra.Command{
		Use:   "audit [path] [flags]",
		Short: "Comprehensive security, quality, and compliance auditing",
		Long: `Audit project security, code quality, license compliance, and performance.

Provides detailed reports with scores and recommendations for improvement.`,
		RunE: c.runAudit,
		Example: `  # Audit current directory
  generator audit
  
  # Audit specific project
  generator audit ./my-project
  
  # Security audit only
  generator audit --security
  
  # Generate detailed report
  generator audit --detailed --output-format html`,
	}

	auditCmd.Flags().Bool("security", true, "Perform security vulnerability scanning")
	auditCmd.Flags().Bool("quality", true, "Perform code quality analysis")
	auditCmd.Flags().Bool("licenses", true, "Perform license compliance checking")
	auditCmd.Flags().Bool("performance", true, "Perform performance analysis")
	auditCmd.Flags().String("output-format", "text", "Output format (text, json, html)")
	auditCmd.Flags().String("output-file", "", "Save audit report to file")
	auditCmd.Flags().Bool("detailed", false, "Generate detailed audit report")

	// Additional audit flags
	auditCmd.Flags().Bool("fail-on-high", false, "Fail if high severity issues are found")
	auditCmd.Flags().Bool("fail-on-medium", false, "Fail if medium or higher severity issues are found")
	auditCmd.Flags().Float64("min-score", 0.0, "Minimum acceptable audit score (0.0-10.0)")
	auditCmd.Flags().StringSlice("exclude-categories", []string{}, "Exclude specific audit categories")
	auditCmd.Flags().Bool("summary-only", false, "Show only audit summary")

	c.rootCmd.AddCommand(auditCmd)
}

// runAudit handles the audit command
func (c *CLI) runAudit(cmd *cobra.Command, args []string) error {
	// Get path from args or use current directory
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Get flags
	security, _ := cmd.Flags().GetBool("security")
	quality, _ := cmd.Flags().GetBool("quality")
	licenses, _ := cmd.Flags().GetBool("licenses")
	performance, _ := cmd.Flags().GetBool("performance")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	outputFile, _ := cmd.Flags().GetString("output-file")
	detailed, _ := cmd.Flags().GetBool("detailed")
	failOnHigh, _ := cmd.Flags().GetBool("fail-on-high")
	failOnMedium, _ := cmd.Flags().GetBool("fail-on-medium")
	minScore, _ := cmd.Flags().GetFloat64("min-score")
	// Additional audit flags (for future implementation)
	excludeCategories, _ := cmd.Flags().GetStringSlice("exclude-categories")
	summaryOnly, _ := cmd.Flags().GetBool("summary-only")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	globalOutputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = c.isNonInteractiveMode()
	}

	// Use global output format if command-specific format not set
	if outputFormat == "text" && globalOutputFormat != "text" {
		outputFormat = globalOutputFormat
	}

	// Log additional options for debugging
	if len(excludeCategories) > 0 {
		c.DebugOutput("Excluding audit categories: %v", excludeCategories)
	}
	if summaryOnly {
		c.DebugOutput("Showing summary only")
	}

	// Create audit options
	options := interfaces.AuditOptions{
		Security:     security,
		Quality:      quality,
		Licenses:     licenses,
		Performance:  performance,
		OutputFormat: outputFormat,
		OutputFile:   outputFile,
		Detailed:     detailed,
	}

	// Perform audit using the existing method
	result, err := c.AuditProject(path, options)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			c.outputFormatter.Error("Audit failed."),
			c.outputFormatter.Info(err.Error()))
	}

	// Display results based on format
	switch outputFormat {
	case "json":
		// JSON output would be implemented here
		c.QuietOutput("JSON output not yet implemented")
	case "html":
		// HTML output would be implemented here
		c.QuietOutput("HTML output not yet implemented")
	default:
		// Beautiful text output
		c.QuietOutput("🔍 Audit Results")
		c.QuietOutput("═══════════════════════════════════════════════════")
		c.QuietOutput("")
		c.QuietOutput("📂 Project: %s", path)
		c.QuietOutput("📊 Overall Score: %.1f/100", result.OverallScore)
		c.QuietOutput("")

		// Show security score if available
		if result.Security != nil {
			c.QuietOutput("🔒 Security Score: %.1f/100", result.Security.Score)
		}

		// Show quality score if available
		if result.Quality != nil {
			c.QuietOutput("✨ Quality Score: %.1f/100", result.Quality.Score)
		}

		// Show license info if available
		if result.Licenses != nil {
			c.QuietOutput("📄 License Compatible: %t", result.Licenses.Compatible)
		}

		// Show performance score if available
		if result.Performance != nil {
			c.QuietOutput("⚡ Performance Score: %.1f/100", result.Performance.Score)
		}

		// Show recommendations if available
		if len(result.Recommendations) > 0 {
			c.QuietOutput("")
			c.QuietOutput("💡 Recommendations:")
			for i, rec := range result.Recommendations {
				c.QuietOutput("    %d. %s", i+1, rec)
			}
		}
	}

	// Check if audit should fail based on criteria
	if failOnHigh && result.OverallScore < 70 {
		return fmt.Errorf("🚫 %s %s",
			c.outputFormatter.Error("Audit failed due to high severity issues."),
			c.outputFormatter.Info("Use --fail-on-high=false to allow high severity issues"))
	}

	if failOnMedium && result.OverallScore < 50 {
		return fmt.Errorf("🚫 %s %s",
			c.outputFormatter.Error("Audit failed due to medium or high severity issues."),
			c.outputFormatter.Info("Use --fail-on-medium=false to allow medium severity issues"))
	}

	if minScore > 0 && result.OverallScore < minScore {
		return fmt.Errorf("🚫 %s %s",
			c.outputFormatter.Error(fmt.Sprintf("Audit score %.1f is below minimum required score %.1f", result.OverallScore, minScore)),
			c.outputFormatter.Info("Use --min-score to adjust the minimum score requirement"))
	}

	return nil
}
