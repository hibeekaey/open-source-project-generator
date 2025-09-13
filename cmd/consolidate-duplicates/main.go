package main

import (
	"fmt"
	"log"
	"os"

	"github.com/open-source-template-generator/internal/cleanup"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <project-root> [--dry-run]")
	}

	projectRoot := os.Args[1]
	dryRun := len(os.Args) > 2 && os.Args[2] == "--dry-run"

	consolidator := cleanup.NewCodeConsolidator(projectRoot)

	fmt.Println("Creating consolidation plan...")
	plan, err := consolidator.CreateConsolidationPlan()
	if err != nil {
		log.Fatalf("Failed to create consolidation plan: %v", err)
	}

	fmt.Printf("Found %d duplicates to consolidate\n", len(plan.Duplicates))

	if dryRun {
		fmt.Println("Running in dry-run mode...")
	}

	if err := consolidator.ExecuteConsolidationPlan(plan, dryRun); err != nil {
		log.Fatalf("Failed to execute consolidation plan: %v", err)
	}

	if !dryRun {
		fmt.Println("Consolidation completed successfully!")
	}
}
