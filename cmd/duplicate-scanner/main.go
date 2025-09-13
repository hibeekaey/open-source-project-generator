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
	duplicates, err := analyzer.FindDuplicateCode(projectRoot)
	if err != nil {
		log.Fatalf("Failed to analyze duplicate code: %v", err)
	}

	fmt.Printf("Found %d duplicate code blocks:\n\n", len(duplicates))

	for i, dup := range duplicates {
		fmt.Printf("=== Duplicate %d ===\n", i+1)
		fmt.Printf("Content: %s\n", dup.Content)
		fmt.Printf("Similarity: %.2f\n", dup.Similarity)
		fmt.Printf("Files:\n")
		for j, file := range dup.Files {
			if j < len(dup.StartLines) && j < len(dup.EndLines) {
				fmt.Printf("  - %s (lines %d-%d)\n", file, dup.StartLines[j], dup.EndLines[j])
			} else {
				fmt.Printf("  - %s\n", file)
			}
		}
		fmt.Printf("Suggestion: %s\n\n", dup.Suggestion)
	}
}
