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
		projectRoot = flag.String("root", ".", "Project root directory")
		outputFile  = flag.String("output", "todo-resolution-report.md", "Output file for the resolution report")
		dryRun      = flag.Bool("dry-run", false, "Perform a dry run without making changes")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
	)
	flag.Parse()

	// Get absolute path
	absRoot, err := filepath.Abs(*projectRoot)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if *verbose {
		fmt.Printf("Resolving TODOs in: %s\n", absRoot)
		fmt.Printf("Output file: %s\n", *outputFile)
		fmt.Printf("Dry run: %v\n", *dryRun)
	}

	// Create resolver config
	config := &cleanup.TODOResolverConfig{
		DryRun:  *dryRun,
		Verbose: *verbose,
	}

	// Create TODO resolver
	resolver := cleanup.NewTODOResolver(absRoot, config)

	// Resolve remaining TODOs
	fmt.Println("Analyzing and resolving remaining TODOs...")
	report, err := resolver.ResolveRemainingTODOs()
	if err != nil {
		log.Fatalf("Failed to resolve TODOs: %v", err)
	}

	// Generate and save report
	reportContent := report.GenerateReport()

	err = os.WriteFile(*outputFile, []byte(reportContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}

	// Print summary
	fmt.Printf("\nTODO Resolution Complete!\n")
	fmt.Printf("Total TODOs found: %d\n", report.TotalFound)
	fmt.Printf("Resolved: %d\n", len(report.Resolved))
	fmt.Printf("Documented for future work: %d\n", len(report.Documented))
	fmt.Printf("Removed (obsolete): %d\n", len(report.Removed))
	fmt.Printf("False positives: %d\n", len(report.FalsePositives))
	fmt.Printf("\nDetailed report saved to: %s\n", *outputFile)

	if *dryRun {
		fmt.Println("\nNote: This was a dry run. No files were modified.")
	}
}
