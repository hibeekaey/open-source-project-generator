package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-source-template-generator/pkg/security"
)

func main() {
	var (
		templateDir = flag.String("dir", "templates", "Directory containing template files to fix")
		dryRun      = flag.Bool("dry-run", false, "Show what would be fixed without making changes")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		fixType     = flag.String("fix-type", "all", "Type of fixes to apply: all, cors, headers, auth, sql")
		backup      = flag.Bool("backup", true, "Create backup files before applying fixes")
	)
	flag.Parse()

	fixer := security.NewFixer()

	if *verbose {
		fmt.Printf("Scanning and fixing templates in directory: %s\n", *templateDir)
		if *dryRun {
			fmt.Println("Running in dry-run mode - no changes will be made")
		}
	}

	// Configure fixer options
	options := security.FixerOptions{
		DryRun:       *dryRun,
		Verbose:      *verbose,
		FixType:      *fixType,
		CreateBackup: *backup,
	}

	result, err := fixer.FixDirectory(*templateDir, options)
	if err != nil {
		log.Fatalf("Error fixing directory: %v", err)
	}

	// Print results
	printFixResults(result, *verbose)

	if result.HasErrors() {
		os.Exit(1)
	}
}

func printFixResults(result *security.FixResult, verbose bool) {
	fmt.Printf("Security Fix Report\n")
	fmt.Printf("==================\n\n")

	if len(result.FixedIssues) == 0 && len(result.Errors) == 0 {
		fmt.Println("No security issues found to fix!")
		return
	}

	if len(result.FixedIssues) > 0 {
		fmt.Printf("Fixed Issues (%d):\n", len(result.FixedIssues))
		fmt.Println("------------------")
		for _, fix := range result.FixedIssues {
			fmt.Printf("✓ %s (Line %d): %s\n", fix.FilePath, fix.LineNumber, fix.Description)
			if verbose {
				fmt.Printf("  Fix: %s\n", fix.FixDescription)
			}
		}
		fmt.Println()
	}

	if len(result.Errors) > 0 {
		fmt.Printf("Errors (%d):\n", len(result.Errors))
		fmt.Println("-------------")
		for _, err := range result.Errors {
			fmt.Printf("✗ %s: %s\n", err.FilePath, err.Error)
		}
		fmt.Println()
	}

	fmt.Printf("Summary: %d issues fixed, %d errors\n", len(result.FixedIssues), len(result.Errors))
	if result.BackupsCreated > 0 {
		fmt.Printf("Backup files created: %d\n", result.BackupsCreated)
	}
}
