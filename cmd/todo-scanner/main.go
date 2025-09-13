package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/internal/cleanup"
)

func main() {
	var (
		projectRoot = flag.String("root", ".", "Project root directory to scan")
		outputFile  = flag.String("output", "", "Output file path (default: stdout)")
		format      = flag.String("format", "markdown", "Output format: markdown, text, json")
		verbose     = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	// Resolve absolute path
	absRoot, err := filepath.Abs(*projectRoot)
	if err != nil {
		log.Fatalf("Failed to resolve project root: %v", err)
	}

	if *verbose {
		fmt.Printf("Scanning project: %s\n", absRoot)
	}

	// Create scanner configuration
	config := &cleanup.TODOScanConfig{
		SkipPatterns:    []string{"vendor/", ".git/", "node_modules/", ".cleanup-backups/"},
		IncludePatterns: []string{"*.go", "*.md", "*.yaml", "*.yml", "*.json"},
		CustomKeywords:  []string{"TODO", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},
		OutputFormat:    *format,
	}

	// Create scanner
	scanner := cleanup.NewTODOScanner(config)

	// Scan project
	if *verbose {
		fmt.Println("Scanning for TODO/FIXME comments...")
	}

	report, err := scanner.ScanProject(absRoot)
	if err != nil {
		log.Fatalf("Failed to scan project: %v", err)
	}

	if *verbose {
		fmt.Printf("Found %d TODO/FIXME comments in %d files\n",
			report.Summary.TotalTODOs, report.FilesScanned)
	}

	// Generate report
	reportContent, err := scanner.GenerateReport(report)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Output report
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(reportContent), 0644)
		if err != nil {
			log.Fatalf("Failed to write report to file: %v", err)
		}
		if *verbose {
			fmt.Printf("Report written to: %s\n", *outputFile)
		}
	} else {
		fmt.Print(reportContent)
	}

	// Print summary to stderr if outputting to file
	if *outputFile != "" {
		fmt.Fprintf(os.Stderr, "TODO Analysis Summary:\n")
		fmt.Fprintf(os.Stderr, "  Total TODOs: %d\n", report.Summary.TotalTODOs)
		fmt.Fprintf(os.Stderr, "  Critical: %d, High: %d, Medium: %d, Low: %d\n",
			report.Summary.CriticalTODOs, report.Summary.HighTODOs,
			report.Summary.MediumTODOs, report.Summary.LowTODOs)
		fmt.Fprintf(os.Stderr, "  Security: %d, Performance: %d, Features: %d, Bugs: %d\n",
			report.Summary.SecurityTODOs, report.Summary.PerformanceTODOs,
			report.Summary.FeatureTODOs, report.Summary.BugTODOs)
	}
}
