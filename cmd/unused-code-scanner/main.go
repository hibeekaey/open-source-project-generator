package main

import (
	"fmt"
	"log"
	"os"

	"github.com/open-source-template-generator/internal/cleanup"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <project-root>")
	}

	projectRoot := os.Args[1]

	analyzer := cleanup.NewCodeAnalyzer()
	unused, err := analyzer.IdentifyUnusedCode(projectRoot)
	if err != nil {
		log.Fatalf("Failed to analyze unused code: %v", err)
	}

	fmt.Printf("Found %d unused code items:\n\n", len(unused))

	// Group by type
	byType := make(map[string][]cleanup.UnusedCodeItem)
	for _, item := range unused {
		byType[item.Type] = append(byType[item.Type], item)
	}

	for itemType, items := range byType {
		fmt.Printf("=== %s (%d items) ===\n", itemType, len(items))
		for i, item := range items {
			if i >= 10 { // Limit output for readability
				fmt.Printf("... and %d more\n", len(items)-10)
				break
			}
			fmt.Printf("- %s:%d - %s (%s)\n", item.File, item.Line, item.Name, item.Reason)
		}
		fmt.Println()
	}
}
