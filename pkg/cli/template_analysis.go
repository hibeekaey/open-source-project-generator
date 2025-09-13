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
	if _, statErr := os.Stat(templateDir); os.IsNotExist(statErr) {
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
	c.displayAnalysisHeader()
	c.displayTemplateSummary(analysis.Templates)
	c.displayInconsistencies(analysis.Inconsistencies)
	c.displayMissingFiles(analysis.MissingFiles)
	c.displayVersionReferences(analysis.VersionReferences)
	c.displayDependencyPatterns(analysis.DependencyPatterns)
	c.displayDetailedTemplateInfo(analysis.Templates)
	return nil
}

// displayAnalysisHeader displays the analysis report header
func (c *CLI) displayAnalysisHeader() {
	fmt.Println("=== Frontend Template Configuration Analysis ===")
}

// displayTemplateSummary displays a summary of found templates
func (c *CLI) displayTemplateSummary(templates []template.TemplateInfo) {
	fmt.Printf("Found %d frontend templates:\n", len(templates))
	for _, tmpl := range templates {
		fmt.Printf("  - %s (%s) - %d config files", tmpl.Name, tmpl.Type, len(tmpl.ConfigFiles))
		if tmpl.Port != "" {
			fmt.Printf(" - Port: %s", tmpl.Port)
		}
		fmt.Println()
	}
	fmt.Println()
}

// displayInconsistencies displays found inconsistencies
func (c *CLI) displayInconsistencies(inconsistencies []template.Inconsistency) {
	if len(inconsistencies) > 0 {
		fmt.Printf("🚨 Found %d inconsistencies:\n", len(inconsistencies))
		for i, inconsistency := range inconsistencies {
			c.displaySingleInconsistency(i+1, inconsistency)
		}
	} else {
		fmt.Println("✅ No inconsistencies found")
	}
}

// displaySingleInconsistency displays details of a single inconsistency
func (c *CLI) displaySingleInconsistency(index int, inconsistency template.Inconsistency) {
	fmt.Printf("  %d. %s\n", index, inconsistency.Description)
	fmt.Printf("     Type: %s\n", inconsistency.Type)
	fmt.Printf("     Templates: %s\n", strings.Join(inconsistency.Templates, ", "))
	if inconsistency.Details != "" {
		fmt.Printf("     Details: %s\n", inconsistency.Details)
	}
	fmt.Println()
}

// displayMissingFiles displays missing configuration files
func (c *CLI) displayMissingFiles(missingFiles []template.MissingFile) {
	if len(missingFiles) > 0 {
		fmt.Printf("📁 Missing configuration files:\n")
		for _, missing := range missingFiles {
			fmt.Printf("  - %s missing from %s: %s\n", missing.File, missing.Template, missing.Reason)
		}
		fmt.Println()
	} else {
		fmt.Println("✅ All required configuration files present")
	}
}

// displayVersionReferences displays version references found in templates
func (c *CLI) displayVersionReferences(versionRefs map[string][]string) {
	if len(versionRefs) > 0 {
		fmt.Printf("🔢 Version references found:\n")
		for version, templates := range versionRefs {
			fmt.Printf("  - %s: used in %s\n", version, strings.Join(templates, ", "))
		}
		fmt.Println()
	}
}

// displayDependencyPatterns displays dependency patterns (limited to first 10)
func (c *CLI) displayDependencyPatterns(patterns map[string]template.DependencyInfo) {
	if len(patterns) == 0 {
		return
	}

	fmt.Printf("📦 Dependency patterns (showing first 10):\n")
	count := 0
	for _, dep := range patterns {
		if count >= 10 {
			fmt.Printf("  ... and %d more dependencies\n", len(patterns)-10)
			break
		}
		fmt.Printf("  - %s: used in %s\n", dep.Package, strings.Join(dep.Templates, ", "))
		count++
	}
	fmt.Println()
}

// displayDetailedTemplateInfo displays detailed information for each template
func (c *CLI) displayDetailedTemplateInfo(templates []template.TemplateInfo) {
	fmt.Println("=== Detailed Template Information ===")
	for _, tmpl := range templates {
		c.displaySingleTemplateInfo(tmpl)
	}
}

// displaySingleTemplateInfo displays detailed information for a single template
func (c *CLI) displaySingleTemplateInfo(tmpl template.TemplateInfo) {
	fmt.Printf("Template: %s (%s)\n", tmpl.Name, tmpl.Type)
	fmt.Printf("  Path: %s\n", tmpl.Path)

	c.displayConfigFiles(tmpl.ConfigFiles)
	c.displayScripts(tmpl.Scripts)
	c.displayDependencies("Dependencies", tmpl.Dependencies)
	c.displayDependencies("Dev Dependencies", tmpl.DevDependencies)

	fmt.Println()
}

// displayConfigFiles displays configuration files for a template
func (c *CLI) displayConfigFiles(configFiles []string) {
	fmt.Printf("  Configuration Files:\n")
	for _, file := range configFiles {
		fmt.Printf("    - %s\n", file)
	}
}

// displayScripts displays NPM scripts for a template
func (c *CLI) displayScripts(scripts map[string]string) {
	if len(scripts) > 0 {
		fmt.Printf("  NPM Scripts:\n")
		for script, command := range scripts {
			fmt.Printf("    - %s: %s\n", script, command)
		}
	}
}

// displayDependencies displays dependencies with optional truncation
func (c *CLI) displayDependencies(label string, deps []string) {
	if len(deps) == 0 {
		return
	}

	maxDisplay := 5
	depsToShow := deps
	if len(deps) > maxDisplay {
		depsToShow = deps[:maxDisplay]
	}

	fmt.Printf("  %s (%d): %s\n", label, len(deps), strings.Join(depsToShow, ", "))
	if len(deps) > maxDisplay {
		fmt.Printf("    ... and %d more\n", len(deps)-maxDisplay)
	}
}

// SaveAnalysisReportCommand saves the analysis report to a JSON file
func (c *CLI) SaveAnalysisReportCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: save-analysis-report <template-directory> <output-file>")
	}

	templateDir := args[0]
	outputFile := args[1]

	scanner := template.NewTemplateScanner(templateDir)
	analysis, scanErr := scanner.ScanFrontendTemplates()
	if scanErr != nil {
		return fmt.Errorf("failed to analyze templates: %w", scanErr)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save analysis to JSON file
	jsonData, marshalErr := json.MarshalIndent(analysis, "", "  ")
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal analysis to JSON: %w", marshalErr)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write analysis report: %w", err)
	}

	fmt.Printf("Analysis report saved to: %s\n", outputFile)
	return nil
}
