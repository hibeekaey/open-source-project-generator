package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/pkg/cli"
	"github.com/open-source-template-generator/pkg/template"
)

// runAnalyzeCommand runs the template analysis command
func runAnalyzeCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: analyze-templates <template-directory> [output-file]")
	}

	templateDir := args[0]

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template directory does not exist: %s", templateDir)
	}

	// Create CLI instance for display functions
	cliInstance := &cli.CLI{}

	scanner := template.NewTemplateScanner(templateDir)
	_, err := scanner.ScanFrontendTemplates()
	if err != nil {
		return fmt.Errorf("failed to analyze templates: %w", err)
	}

	// Display the analysis report
	if err := cliInstance.AnalyzeTemplatesCommand([]string{templateDir}); err != nil {
		return err
	}

	// If output file is specified, save the report
	if len(args) > 1 {
		outputFile := args[1]

		// Create output directory if it doesn't exist
		outputDir := filepath.Dir(outputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := cliInstance.SaveAnalysisReportCommand([]string{templateDir, outputFile}); err != nil {
			return err
		}
	}

	return nil
}
