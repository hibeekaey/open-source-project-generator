package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/pkg/template/standards"
)

func main() {
	var (
		templatesDir = flag.String("templates", "templates", "Path to templates directory")
		action       = flag.String("action", "validate", "Action to perform: validate, apply, compare, docs")
		outputFile   = flag.String("output", "", "Output file for reports (optional)")
		templateType = flag.String("template", "", "Specific template type to process (optional)")
	)
	flag.Parse()

	manager := standards.NewStandardsManager()

	switch *action {
	case "validate":
		if err := validateTemplates(manager, *templatesDir, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error validating templates: %v\n", err)
			os.Exit(1)
		}

	case "apply":
		if err := applyStandards(manager, *templatesDir, *templateType); err != nil {
			fmt.Fprintf(os.Stderr, "Error applying standards: %v\n", err)
			os.Exit(1)
		}

	case "compare":
		if err := compareTemplates(manager, *templatesDir, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error comparing templates: %v\n", err)
			os.Exit(1)
		}

	case "docs":
		if err := generateDocs(manager, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating documentation: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", *action)
		fmt.Fprintf(os.Stderr, "Available actions: validate, apply, compare, docs\n")
		os.Exit(1)
	}
}

func validateTemplates(manager *standards.StandardsManager, templatesDir, outputFile string) error {
	fmt.Println("🔍 Validating frontend templates...")

	results, err := manager.ValidateAllFrontendTemplates(templatesDir)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Generate report
	report := manager.GenerateValidationReport(results)

	if outputFile != "" {
		// Write report to file
		if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
			return fmt.Errorf("failed to write report to %s: %w", outputFile, err)
		}
		fmt.Printf("📄 Validation report written to %s\n", outputFile)

		// Also export JSON results
		jsonFile := outputFile + ".json"
		if err := manager.ExportValidationResults(results, jsonFile); err != nil {
			return fmt.Errorf("failed to export JSON results: %w", err)
		}
		fmt.Printf("📊 JSON results exported to %s\n", jsonFile)
	} else {
		// Print report to stdout
		fmt.Println(report)
	}

	// Summary
	totalTemplates := len(results)
	validTemplates := 0
	for _, result := range results {
		if result.IsValid {
			validTemplates++
		}
	}

	if validTemplates == totalTemplates {
		fmt.Printf("✅ All %d templates are valid!\n", totalTemplates)
	} else {
		fmt.Printf("❌ %d out of %d templates have issues\n", totalTemplates-validTemplates, totalTemplates)
	}

	return nil
}

func applyStandards(manager *standards.StandardsManager, templatesDir, templateType string) error {
	// Create template updater
	updater := standards.NewTemplateUpdater()

	// Example version map (in real usage, this would come from version management system)
	versions := map[string]string{
		"NextJS": "15.5.2",
		"React":  "19.1.0",
	}

	if templateType != "" {
		// Apply to specific template
		fmt.Printf("🔧 Applying standards to %s template...\n", templateType)
		templatePath := filepath.Join(templatesDir, "frontend", templateType)

		if err := updater.UpdateTemplate(templatePath, templateType, versions); err != nil {
			return fmt.Errorf("failed to apply standards to %s: %w", templateType, err)
		}

		fmt.Printf("✅ Successfully applied standards to %s\n", templateType)
	} else {
		// Apply to all templates
		fmt.Println("🔧 Applying standards to all frontend templates...")

		if err := updater.UpdateAllTemplates(templatesDir, versions); err != nil {
			return fmt.Errorf("failed to apply standards: %w", err)
		}

		fmt.Println("✅ Successfully applied standards to all templates")
	}

	return nil
}

func compareTemplates(manager *standards.StandardsManager, templatesDir, outputFile string) error {
	fmt.Println("📊 Comparing template configurations...")

	comparison, err := manager.CompareTemplateConfigurations(templatesDir)
	if err != nil {
		return fmt.Errorf("comparison failed: %w", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(comparison), 0644); err != nil {
			return fmt.Errorf("failed to write comparison to %s: %w", outputFile, err)
		}
		fmt.Printf("📄 Comparison report written to %s\n", outputFile)
	} else {
		fmt.Println(comparison)
	}

	return nil
}

func generateDocs(manager *standards.StandardsManager, outputFile string) error {
	fmt.Println("📚 Generating standards documentation...")

	docs := manager.GenerateStandardsDocumentation()

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(docs), 0644); err != nil {
			return fmt.Errorf("failed to write documentation to %s: %w", outputFile, err)
		}
		fmt.Printf("📄 Documentation written to %s\n", outputFile)
	} else {
		fmt.Println(docs)
	}

	return nil
}
