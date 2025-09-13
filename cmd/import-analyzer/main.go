package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/open-source-template-generator/pkg/template"
)

func main() {
	var (
		templateDir = flag.String("dir", "templates", "Directory containing template files to analyze")
		outputFile  = flag.String("output", "", "Output file for the report (default: stdout)")
		jsonOutput  = flag.Bool("json", false, "Output report in JSON format")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
	)
	flag.Parse()

	// Validate template directory exists
	if _, err := os.Stat(*templateDir); os.IsNotExist(err) {
		log.Fatalf("Template directory does not exist: %s", *templateDir)
	}

	// Create detector
	detector := template.NewImportDetector()

	if *verbose {
		fmt.Printf("Analyzing template files in: %s\n", *templateDir)
	}

	// Analyze directory
	analysis, err := detector.AnalyzeDirectory(*templateDir)
	if err != nil {
		log.Fatalf("Failed to analyze directory: %v", err)
	}

	// Generate output
	var output string
	if *jsonOutput {
		jsonData, err := json.MarshalIndent(analysis, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		output = string(jsonData)
	} else {
		output = detector.GenerateTextReport(analysis)
	}

	// Write output
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Report written to: %s\n", *outputFile)
	} else {
		fmt.Print(output)
	}

	// Print summary to stderr if writing to file
	if *outputFile != "" && *verbose {
		fmt.Fprintf(os.Stderr, "Analysis complete:\n")
		fmt.Fprintf(os.Stderr, "  Files analyzed: %d\n", analysis.Summary.TotalFiles)
		fmt.Fprintf(os.Stderr, "  Files with issues: %d\n", analysis.Summary.FilesWithIssues)
		fmt.Fprintf(os.Stderr, "  Total missing imports: %d\n", analysis.Summary.TotalMissingImports)
	}

	// Exit with error code if issues found
	if analysis.Summary.FilesWithIssues > 0 {
		os.Exit(1)
	}
}
