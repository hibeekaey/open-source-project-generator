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

	fmt.Println("Analyzing unused code...")
	analyzer := cleanup.NewCodeAnalyzer()
	unused, err := analyzer.IdentifyUnusedCode(projectRoot)
	if err != nil {
		log.Fatalf("Failed to analyze unused code: %v", err)
	}

	fmt.Printf("Found %d unused code items\n", len(unused))

	// Filter to only safe-to-remove items for actual removal
	var safeToRemove []cleanup.UnusedCodeItem
	remover := cleanup.NewCodeRemover(dryRun)

	for _, item := range unused {
		// Only include imports and clearly unused items
		if item.Type == "import" || (item.Type != "import" && len(item.Name) > 0 && item.Name[0] >= 'a' && item.Name[0] <= 'z') {
			safeToRemove = append(safeToRemove, item)
		}
	}

	fmt.Printf("Safe to remove: %d items\n", len(safeToRemove))

	plan := remover.CreateRemovalPlan(safeToRemove)

	if dryRun {
		fmt.Println("Running in dry-run mode...")
	}

	if err := remover.ExecuteRemovalPlan(plan); err != nil {
		log.Fatalf("Failed to execute removal plan: %v", err)
	}

	if !dryRun {
		fmt.Println("Unused code removal completed successfully!")
	}
}
