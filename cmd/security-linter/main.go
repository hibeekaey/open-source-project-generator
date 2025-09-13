package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/security"
)

// CLI configuration
type Config struct {
	Directory      string
	OutputFormat   string
	OutputFile     string
	Severity       string
	Categories     string
	Verbose        bool
	FailOnHigh     bool
	FailOnCritical bool
	DryRun         bool
	ConfigFile     string
}

func main() {
	config := parseFlags()

	if config.Verbose {
		fmt.Printf("Security Linter v1.0.0\n")
		fmt.Printf("Scanning directory: %s\n", config.Directory)
	}

	// Create security linter
	linter := security.NewSecurityLinter()

	// Perform linting
	result, err := linter.LintDirectory(config.Directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during linting: %v\n", err)
		os.Exit(1)
	}

	// Filter results based on configuration
	filteredResult := filterResults(result, config)

	// Print summary
	printSummary(filteredResult, config)

	// Export results if requested
	if config.OutputFile != "" {
		if err := linter.ExportResults(filteredResult, config.OutputFormat, config.OutputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting results: %v\n", err)
			os.Exit(1)
		}
		if config.Verbose {
			fmt.Printf("Results exported to: %s\n", config.OutputFile)
		}
	}

	// Determine exit code based on findings
	exitCode := determineExitCode(filteredResult, config)
	if exitCode != 0 {
		if config.Verbose {
			fmt.Printf("Exiting with code %d due to security issues\n", exitCode)
		}
	}

	os.Exit(exitCode)
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.Directory, "dir", ".", "Directory to scan for security issues")
	flag.StringVar(&config.OutputFormat, "format", "json", "Output format (json, sarif, junit)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file path (stdout if not specified)")
	flag.StringVar(&config.Severity, "severity", "all", "Minimum severity level (low, medium, high, critical, all)")
	flag.StringVar(&config.Categories, "categories", "all", "Categories to check (comma-separated or 'all')")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.FailOnHigh, "fail-on-high", false, "Exit with non-zero code on high severity issues")
	flag.BoolVar(&config.FailOnCritical, "fail-on-critical", true, "Exit with non-zero code on critical severity issues")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Perform dry run without making changes")
	flag.StringVar(&config.ConfigFile, "config", "", "Configuration file path")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Security Linter - Automated security validation tool\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -dir ./src -format sarif -output security-report.sarif\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -severity high -fail-on-high -verbose\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -categories cryptography,sql-injection -format junit -output results.xml\n", os.Args[0])
	}

	flag.Parse()

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	return config
}

func validateConfig(config *Config) error {
	// Validate directory exists
	if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", config.Directory)
	}

	// Validate output format
	validFormats := []string{"json", "sarif", "junit"}
	formatValid := false
	for _, format := range validFormats {
		if config.OutputFormat == format {
			formatValid = true
			break
		}
	}
	if !formatValid {
		return fmt.Errorf("invalid output format: %s (valid: %s)", config.OutputFormat, strings.Join(validFormats, ", "))
	}

	// Validate severity level
	validSeverities := []string{"low", "medium", "high", "critical", "all"}
	severityValid := false
	for _, severity := range validSeverities {
		if config.Severity == severity {
			severityValid = true
			break
		}
	}
	if !severityValid {
		return fmt.Errorf("invalid severity level: %s (valid: %s)", config.Severity, strings.Join(validSeverities, ", "))
	}

	// Create output directory if needed
	if config.OutputFile != "" {
		outputDir := filepath.Dir(config.OutputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	return nil
}

func filterResults(result *security.LintResult, config Config) *security.LintResult {
	filtered := &security.LintResult{
		Issues:       make([]security.LintIssue, 0),
		Summary:      result.Summary,
		ScannedFiles: result.ScannedFiles,
		RulesApplied: result.RulesApplied,
	}

	// Parse categories filter
	var allowedCategories map[string]bool
	if config.Categories != "all" {
		allowedCategories = make(map[string]bool)
		for _, category := range strings.Split(config.Categories, ",") {
			allowedCategories[strings.TrimSpace(category)] = true
		}
	}

	// Parse severity filter
	minSeverity := parseSeverityLevel(config.Severity)

	// Filter issues
	for _, issue := range result.Issues {
		// Check category filter
		if allowedCategories != nil && !allowedCategories[issue.Category] {
			continue
		}

		// Check severity filter
		if !meetsSeverityThreshold(issue.Severity, minSeverity) {
			continue
		}

		filtered.Issues = append(filtered.Issues, issue)
	}

	// Recalculate summary for filtered results
	recalculateSummary(filtered)

	return filtered
}

func parseSeverityLevel(severity string) security.SeverityLevel {
	switch strings.ToLower(severity) {
	case "low":
		return security.SeverityLow
	case "medium":
		return security.SeverityMedium
	case "high":
		return security.SeverityHigh
	case "critical":
		return security.SeverityCritical
	default:
		return security.SeverityLow
	}
}

func meetsSeverityThreshold(issueSeverity, minSeverity security.SeverityLevel) bool {
	severityOrder := map[security.SeverityLevel]int{
		security.SeverityLow:      1,
		security.SeverityMedium:   2,
		security.SeverityHigh:     3,
		security.SeverityCritical: 4,
	}

	return severityOrder[issueSeverity] >= severityOrder[minSeverity]
}

func recalculateSummary(result *security.LintResult) {
	result.Summary.TotalIssues = len(result.Issues)
	result.Summary.BySeverity = make(map[security.SeverityLevel]int)
	result.Summary.ByCategory = make(map[string]int)

	for _, issue := range result.Issues {
		result.Summary.BySeverity[issue.Severity]++
		result.Summary.ByCategory[issue.Category]++
	}
}

func printSummary(result *security.LintResult, config Config) {
	fmt.Printf("\n=== Security Linting Summary ===\n")
	fmt.Printf("Scanned Files: %d\n", result.ScannedFiles)
	fmt.Printf("Rules Applied: %d\n", result.RulesApplied)
	fmt.Printf("Total Issues: %d\n", result.Summary.TotalIssues)

	if result.Summary.TotalIssues > 0 {
		fmt.Printf("\nIssues by Severity:\n")
		if count := result.Summary.BySeverity[security.SeverityCritical]; count > 0 {
			fmt.Printf("  Critical: %d\n", count)
		}
		if count := result.Summary.BySeverity[security.SeverityHigh]; count > 0 {
			fmt.Printf("  High: %d\n", count)
		}
		if count := result.Summary.BySeverity[security.SeverityMedium]; count > 0 {
			fmt.Printf("  Medium: %d\n", count)
		}
		if count := result.Summary.BySeverity[security.SeverityLow]; count > 0 {
			fmt.Printf("  Low: %d\n", count)
		}

		fmt.Printf("\nIssues by Category:\n")
		for category, count := range result.Summary.ByCategory {
			fmt.Printf("  %s: %d\n", category, count)
		}

		if len(result.Summary.CriticalFiles) > 0 {
			fmt.Printf("\nFiles with Critical Issues:\n")
			for _, file := range result.Summary.CriticalFiles {
				fmt.Printf("  %s\n", file)
			}
		}

		if config.Verbose && len(result.Issues) > 0 {
			fmt.Printf("\nDetailed Issues:\n")
			for _, issue := range result.Issues {
				fmt.Printf("  %s:%d [%s] %s - %s\n",
					issue.FilePath, issue.LineNumber, issue.Severity, issue.RuleID, issue.Message)
				if issue.Suggestion != "" {
					fmt.Printf("    Suggestion: %s\n", issue.Suggestion)
				}
			}
		}
	} else {
		fmt.Printf("✅ No security issues found!\n")
	}
}

func determineExitCode(result *security.LintResult, config Config) int {
	if config.FailOnCritical && result.Summary.BySeverity[security.SeverityCritical] > 0 {
		return 2 // Critical issues found
	}

	if config.FailOnHigh && result.Summary.BySeverity[security.SeverityHigh] > 0 {
		return 1 // High severity issues found
	}

	return 0 // No issues or only low/medium severity
}
