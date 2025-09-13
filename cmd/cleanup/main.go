package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/open-source-template-generator/internal/cleanup"
)

func main() {
	var (
		projectRoot = flag.String("root", ".", "Project root directory")
		dryRun      = flag.Bool("dry-run", false, "Perform dry run without making changes")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		backupDir   = flag.String("backup-dir", "", "Backup directory (default: .cleanup-backups)")
		action      = flag.String("action", "analyze", "Action to perform: analyze, validate, init")
	)
	flag.Parse()

	// Resolve absolute path
	absRoot, err := filepath.Abs(*projectRoot)
	if err != nil {
		log.Fatalf("Failed to resolve project root: %v", err)
	}

	// Create configuration
	config := &cleanup.Config{
		BackupDir: *backupDir,
		DryRun:    *dryRun,
		Verbose:   *verbose,
	}

	// Create cleanup manager
	manager, err := cleanup.NewManager(absRoot, config)
	if err != nil {
		log.Fatalf("Failed to create cleanup manager: %v", err)
	}
	defer manager.Shutdown()

	switch *action {
	case "init":
		if err := initializeCleanup(manager); err != nil {
			log.Fatalf("Failed to initialize cleanup: %v", err)
		}
	case "analyze":
		if err := analyzeProject(manager); err != nil {
			log.Fatalf("Failed to analyze project: %v", err)
		}
	case "validate":
		if err := validateProject(manager); err != nil {
			log.Fatalf("Failed to validate project: %v", err)
		}
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func initializeCleanup(manager *cleanup.Manager) error {
	fmt.Println("Initializing cleanup infrastructure...")

	if err := manager.Initialize(); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	fmt.Println("Cleanup infrastructure initialized successfully!")
	return nil
}

func analyzeProject(manager *cleanup.Manager) error {
	fmt.Println("Analyzing project...")

	// Initialize first
	if err := manager.Initialize(); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Perform analysis
	analysis, err := manager.AnalyzeProject()
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Print summary
	fmt.Println("\n" + analysis.GetSummary())

	// Print detailed results
	if len(analysis.TODOs) > 0 {
		fmt.Println("\nTODO/FIXME Comments:")
		for i, todo := range analysis.TODOs {
			if i >= 10 { // Limit output
				fmt.Printf("  ... and %d more\n", len(analysis.TODOs)-i)
				break
			}
			fmt.Printf("  %s:%d - %s: %s\n", todo.File, todo.Line, todo.Type, todo.Message)
		}
	}

	if len(analysis.Duplicates) > 0 {
		fmt.Println("\nPotential Duplicate Code:")
		for i, dup := range analysis.Duplicates {
			if i >= 5 { // Limit output
				fmt.Printf("  ... and %d more\n", len(analysis.Duplicates)-i)
				break
			}
			fmt.Printf("  Similar code in %d files: %v\n", len(dup.Files), dup.Files)
		}
	}

	if len(analysis.UnusedCode) > 0 {
		fmt.Println("\nUnused Code Items:")
		for i, unused := range analysis.UnusedCode {
			if i >= 10 { // Limit output
				fmt.Printf("  ... and %d more\n", len(analysis.UnusedCode)-i)
				break
			}
			fmt.Printf("  %s:%d - %s %s (%s)\n", unused.File, unused.Line, unused.Type, unused.Name, unused.Reason)
		}
	}

	if len(analysis.ImportIssues) > 0 {
		fmt.Println("\nImport Organization Issues:")
		for i, issue := range analysis.ImportIssues {
			if i >= 10 { // Limit output
				fmt.Printf("  ... and %d more\n", len(analysis.ImportIssues)-i)
				break
			}
			fmt.Printf("  %s:%d - %s: %s\n", issue.File, issue.Line, issue.Type, issue.Suggestion)
		}
	}

	return nil
}

func validateProject(manager *cleanup.Manager) error {
	fmt.Println("Validating project...")

	// Initialize first
	if err := manager.Initialize(); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Perform validation
	result, err := manager.ValidateProject()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Print results
	fmt.Printf("\nValidation Results:\n")
	fmt.Printf("  Success: %t\n", result.Success)
	fmt.Printf("  Build Success: %t\n", result.BuildSuccess)
	fmt.Printf("  Tests Passed: %t\n", result.TestsPassed)
	fmt.Printf("  Duration: %v\n", result.Duration)
	fmt.Printf("  Errors: %d\n", len(result.Errors))
	fmt.Printf("  Warnings: %d\n", len(result.Warnings))

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range result.Errors {
			fmt.Printf("  - %s: %s\n", err.Type, err.Message)
			if err.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", err.Suggestion)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  - %s\n", warning)
		}
	}

	return nil
}
