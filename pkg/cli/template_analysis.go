package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/template"
)

// AnalyzeTemplatesCommand analyzes frontend template configurations
func (c *CLI) AnalyzeTemplatesCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: analyze-templates <template-directory>")
	}

	templateDir := args[0]

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template directory does not exist: %s", templateDir)
	}

	scanner := template.NewTemplateScanner(templateDir)
	analysis, err := scanner.ScanFrontendTemplates()
	if err != nil {
		return fmt.Errorf("failed to analyze templates: %w", err)
	}

	// Generate and display report
	return c.displayAnalysisReport(analysis)
}

// displayAnalysisReport displays the template analysis report
func (c *CLI) displayAnalysisReport(analysis *template.ConfigurationAnalysis) error {
	fmt.Println("=== Frontend Template Configuration Analysis ===")

	// Display template summary
	fmt.Printf("Found %d frontend templates:\n", len(analysis.Templates))
	for _, tmpl := range analysis.Templates {
		fmt.Printf("  - %s (%s) - %d config files", tmpl.Name, tmpl.Type, len(tmpl.ConfigFiles))
		if tmpl.Port != "" {
			fmt.Printf(" - Port: %s", tmpl.Port)
		}
		fmt.Println()
	}
	fmt.Println()

	// Display inconsistencies
	if len(analysis.Inconsistencies) > 0 {
		fmt.Printf("🚨 Found %d inconsistencies:\n", len(analysis.Inconsistencies))
		for i, inconsistency := range analysis.Inconsistencies {
			fmt.Printf("  %d. %s\n", i+1, inconsistency.Description)
			fmt.Printf("     Type: %s\n", inconsistency.Type)
			fmt.Printf("     Templates: %s\n", strings.Join(inconsistency.Templates, ", "))
			if inconsistency.Details != "" {
				fmt.Printf("     Details: %s\n", inconsistency.Details)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("✅ No inconsistencies found")
	}

	// Display missing files
	if len(analysis.MissingFiles) > 0 {
		fmt.Printf("📁 Missing configuration files:\n")
		for _, missing := range analysis.MissingFiles {
			fmt.Printf("  - %s missing from %s: %s\n", missing.File, missing.Template, missing.Reason)
		}
		fmt.Println()
	} else {
		fmt.Println("✅ All required configuration files present")
	}

	// Display version references
	if len(analysis.VersionReferences) > 0 {
		fmt.Printf("🔢 Version references found:\n")
		for version, templates := range analysis.VersionReferences {
			fmt.Printf("  - %s: used in %s\n", version, strings.Join(templates, ", "))
		}
		fmt.Println()
	}

	// Display dependency patterns
	if len(analysis.DependencyPatterns) > 0 {
		fmt.Printf("📦 Dependency patterns (showing first 10):\n")
		count := 0
		for _, dep := range analysis.DependencyPatterns {
			if count >= 10 {
				fmt.Printf("  ... and %d more dependencies\n", len(analysis.DependencyPatterns)-10)
				break
			}
			fmt.Printf("  - %s: used in %s\n", dep.Package, strings.Join(dep.Templates, ", "))
			count++
		}
		fmt.Println()
	}

	// Display detailed template information
	fmt.Println("=== Detailed Template Information ===")
	for _, tmpl := range analysis.Templates {
		fmt.Printf("Template: %s (%s)\n", tmpl.Name, tmpl.Type)
		fmt.Printf("  Path: %s\n", tmpl.Path)
		fmt.Printf("  Configuration Files:\n")
		for _, file := range tmpl.ConfigFiles {
			fmt.Printf("    - %s\n", file)
		}

		if len(tmpl.Scripts) > 0 {
			fmt.Printf("  NPM Scripts:\n")
			for script, command := range tmpl.Scripts {
				fmt.Printf("    - %s: %s\n", script, command)
			}
		}

		if len(tmpl.Dependencies) > 0 {
			fmt.Printf("  Dependencies (%d): %s\n", len(tmpl.Dependencies), strings.Join(tmpl.Dependencies[:min(5, len(tmpl.Dependencies))], ", "))
			if len(tmpl.Dependencies) > 5 {
				fmt.Printf("    ... and %d more\n", len(tmpl.Dependencies)-5)
			}
		}

		if len(tmpl.DevDependencies) > 0 {
			fmt.Printf("  Dev Dependencies (%d): %s\n", len(tmpl.DevDependencies), strings.Join(tmpl.DevDependencies[:min(5, len(tmpl.DevDependencies))], ", "))
			if len(tmpl.DevDependencies) > 5 {
				fmt.Printf("    ... and %d more\n", len(tmpl.DevDependencies)-5)
			}
		}

		fmt.Println()
	}

	return nil
}

// SaveAnalysisReportCommand saves the analysis report to a JSON file
func (c *CLI) SaveAnalysisReportCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: save-analysis-report <template-directory> <output-file>")
	}

	templateDir := args[0]
	outputFile := args[1]

	scanner := template.NewTemplateScanner(templateDir)
	analysis, err := scanner.ScanFrontendTemplates()
	if err != nil {
		return fmt.Errorf("failed to analyze templates: %w", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save analysis to JSON file
	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis to JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write analysis report: %w", err)
	}

	fmt.Printf("Analysis report saved to: %s\n", outputFile)
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
