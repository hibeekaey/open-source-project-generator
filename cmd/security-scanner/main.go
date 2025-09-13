package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-source-template-generator/pkg/security"
)

func main() {
	var (
		templateDir = flag.String("dir", "templates", "Directory to scan for template files")
		outputFile  = flag.String("output", "", "Output file for security report (JSON format)")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		minSeverity = flag.String("severity", "low", "Minimum severity level to report (low, medium, high, critical)")
		issueType   = flag.String("type", "", "Filter by issue type (cors_vulnerability, missing_security_header, weak_authentication, sql_injection_risk, information_leakage)")
		fixableOnly = flag.Bool("fixable-only", false, "Show only issues that have automated fixes available")
		showSummary = flag.Bool("summary", false, "Show only summary statistics")
	)
	flag.Parse()

	scanner := security.NewScanner()

	if *verbose {
		fmt.Printf("Scanning templates in directory: %s\n", *templateDir)
	}

	report, err := scanner.ScanDirectory(*templateDir)
	if err != nil {
		log.Fatalf("Error scanning directory: %v", err)
	}

	// Apply filters
	filteredReport := applyFilters(report, *minSeverity, *issueType, *fixableOnly)

	if *verbose {
		fmt.Printf("Found %d security issues (filtered: %d)\n", len(report.Issues), len(filteredReport.Issues))
	}

	// Output results
	if *outputFile != "" {
		if err := writeReportToFile(filteredReport, *outputFile); err != nil {
			log.Fatalf("Error writing report to file: %v", err)
		}
		fmt.Printf("Security report written to: %s\n", *outputFile)
	} else {
		if *showSummary {
			printSummary(filteredReport)
		} else {
			printReport(filteredReport)
		}
	}

	// Exit with error code if critical issues found
	if filteredReport.HasCriticalIssues() {
		os.Exit(1)
	}
}

func writeReportToFile(report *security.SecurityReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func applyFilters(report *security.SecurityReport, minSeverity, issueType string, fixableOnly bool) *security.SecurityReport {
	filtered := &security.SecurityReport{
		ScannedFiles: report.ScannedFiles,
		TotalLines:   report.TotalLines,
		Issues:       []security.SecurityIssue{},
	}

	// Convert string severity to SeverityLevel
	var minSev security.SeverityLevel
	switch minSeverity {
	case "critical":
		minSev = security.SeverityCritical
	case "high":
		minSev = security.SeverityHigh
	case "medium":
		minSev = security.SeverityMedium
	default:
		minSev = security.SeverityLow
	}

	// Apply severity filter
	severityFiltered := report.FilterBySeverity(minSev)

	// Apply other filters
	for _, issue := range severityFiltered {
		// Filter by issue type if specified
		if issueType != "" && string(issue.IssueType) != issueType {
			continue
		}

		// Filter by fixable only if specified
		if fixableOnly && !issue.FixAvailable {
			continue
		}

		filtered.Issues = append(filtered.Issues, issue)
	}

	return filtered
}

func printReport(report *security.SecurityReport) {
	fmt.Printf("Security Scan Report\n")
	fmt.Printf("===================\n\n")

	if len(report.Issues) == 0 {
		fmt.Println("No security issues found!")
		return
	}

	for _, issue := range report.Issues {
		fmt.Printf("File: %s (Line %d)\n", issue.FilePath, issue.LineNumber)
		fmt.Printf("Type: %s\n", issue.IssueType)
		fmt.Printf("Severity: %s\n", issue.Severity)
		fmt.Printf("Description: %s\n", issue.Description)
		if issue.Recommendation != "" {
			fmt.Printf("Recommendation: %s\n", issue.Recommendation)
		}
		if issue.FixAvailable {
			fmt.Printf("Fix Available: Yes\n")
		}
		fmt.Println()
	}

	printSummary(report)
}

func printSummary(report *security.SecurityReport) {
	fmt.Printf("Summary: %d issues found\n", len(report.Issues))
	fmt.Printf("Critical: %d, High: %d, Medium: %d, Low: %d\n",
		report.CountBySeverity("critical"),
		report.CountBySeverity("high"),
		report.CountBySeverity("medium"),
		report.CountBySeverity("low"))

	fmt.Printf("\nIssue Types:\n")
	fmt.Printf("CORS Vulnerabilities: %d\n", report.CountByType(security.CORSVulnerability))
	fmt.Printf("Missing Security Headers: %d\n", report.CountByType(security.MissingSecurityHeader))
	fmt.Printf("Weak Authentication: %d\n", report.CountByType(security.WeakAuthentication))
	fmt.Printf("SQL Injection Risks: %d\n", report.CountByType(security.SQLInjectionRisk))
	fmt.Printf("Information Leakage: %d\n", report.CountByType(security.InformationLeakage))

	fixableIssues := report.GetFixableIssues()
	fmt.Printf("\nFixable Issues: %d/%d (%.1f%%)\n",
		len(fixableIssues),
		len(report.Issues),
		float64(len(fixableIssues))/float64(len(report.Issues))*100)
}
