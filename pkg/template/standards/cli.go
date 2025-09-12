package standards

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StandardsManager manages frontend template standardization
type StandardsManager struct {
	generator *TemplateGenerator
	validator *TemplateValidator
}

// NewStandardsManager creates a new standards manager
func NewStandardsManager() *StandardsManager {
	return &StandardsManager{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
	}
}

// ValidateAllFrontendTemplates validates all frontend templates against standards
func (sm *StandardsManager) ValidateAllFrontendTemplates(templatesDir string) (map[string]*ValidationResult, error) {
	results := make(map[string]*ValidationResult)

	frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	for _, templateType := range frontendTemplates {
		templatePath := filepath.Join(templatesDir, "frontend", templateType)

		result, err := sm.validator.ValidateTemplate(templatePath, templateType)
		if err != nil {
			return nil, fmt.Errorf("failed to validate template %s: %w", templateType, err)
		}

		results[templateType] = result
	}

	return results, nil
}

// GenerateValidationReport generates a comprehensive validation report
func (sm *StandardsManager) GenerateValidationReport(results map[string]*ValidationResult) string {
	var report strings.Builder

	report.WriteString("# Frontend Template Validation Report\n\n")

	totalTemplates := len(results)
	validTemplates := 0
	totalErrors := 0
	totalWarnings := 0
	totalSuggestions := 0

	for _, result := range results {
		if result.IsValid {
			validTemplates++
		}
		totalErrors += len(result.Errors)
		totalWarnings += len(result.Warnings)
		totalSuggestions += len(result.Suggestions)
	}

	// Summary
	report.WriteString("## Summary\n\n")
	report.WriteString(fmt.Sprintf("- **Total Templates**: %d\n", totalTemplates))
	report.WriteString(fmt.Sprintf("- **Valid Templates**: %d\n", validTemplates))
	report.WriteString(fmt.Sprintf("- **Invalid Templates**: %d\n", totalTemplates-validTemplates))
	report.WriteString(fmt.Sprintf("- **Total Errors**: %d\n", totalErrors))
	report.WriteString(fmt.Sprintf("- **Total Warnings**: %d\n", totalWarnings))
	report.WriteString(fmt.Sprintf("- **Total Suggestions**: %d\n\n", totalSuggestions))

	// Detailed results
	report.WriteString("## Detailed Results\n\n")

	for templateName, result := range results {
		report.WriteString(fmt.Sprintf("### %s\n\n", templateName))

		if result.IsValid {
			report.WriteString("✅ **Status**: Valid\n\n")
		} else {
			report.WriteString("❌ **Status**: Invalid\n\n")
		}

		// Errors
		if len(result.Errors) > 0 {
			report.WriteString("#### Errors\n\n")
			for _, error := range result.Errors {
				report.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", error.Type, error.File, error.Description))
				if error.Expected != "" && error.Actual != "" {
					report.WriteString(fmt.Sprintf("  - Expected: `%s`\n", error.Expected))
					report.WriteString(fmt.Sprintf("  - Actual: `%s`\n", error.Actual))
				}
			}
			report.WriteString("\n")
		}

		// Warnings
		if len(result.Warnings) > 0 {
			report.WriteString("#### Warnings\n\n")
			for _, warning := range result.Warnings {
				report.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", warning.Type, warning.File, warning.Description))
			}
			report.WriteString("\n")
		}

		// Suggestions
		if len(result.Suggestions) > 0 {
			report.WriteString("#### Suggestions\n\n")
			for _, suggestion := range result.Suggestions {
				report.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", suggestion.Type, suggestion.File, suggestion.Description))
			}
			report.WriteString("\n")
		}
	}

	return report.String()
}

// ApplyStandardsToTemplate applies standardized configurations to a template
func (sm *StandardsManager) ApplyStandardsToTemplate(templatePath, templateType string) error {
	// Generate standardized files
	files := sm.generator.GenerateAllStandardizedFiles(templateType)
	filePaths := GetStandardizedFilePaths()

	// Write each standardized file
	fileMap := map[string]string{
		"PackageJSON":    files.PackageJSON,
		"TSConfig":       files.TSConfig,
		"ESLintConfig":   files.ESLintConfig,
		"PrettierConfig": files.PrettierConfig,
		"VercelConfig":   files.VercelConfig,
		"TailwindConfig": files.TailwindConfig,
		"NextConfig":     files.NextConfig,
		"PostCSSConfig":  files.PostCSSConfig,
		"JestConfig":     files.JestConfig,
		"JestSetup":      files.JestSetup,
	}

	for fileType, content := range fileMap {
		fileName := filePaths[fileType]
		filePath := filepath.Join(templatePath, fileName)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fileName, err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}

		fmt.Printf("✅ Updated %s\n", fileName)
	}

	return nil
}

// ApplyStandardsToAllTemplates applies standards to all frontend templates
func (sm *StandardsManager) ApplyStandardsToAllTemplates(templatesDir string) error {
	frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	for _, templateType := range frontendTemplates {
		templatePath := filepath.Join(templatesDir, "frontend", templateType)

		fmt.Printf("Applying standards to %s...\n", templateType)

		if err := sm.ApplyStandardsToTemplate(templatePath, templateType); err != nil {
			return fmt.Errorf("failed to apply standards to %s: %w", templateType, err)
		}

		fmt.Printf("✅ Successfully applied standards to %s\n\n", templateType)
	}

	return nil
}

// GenerateStandardsDocumentation generates documentation for the standards
func (sm *StandardsManager) GenerateStandardsDocumentation() string {
	var doc strings.Builder

	doc.WriteString("# Frontend Template Standards\n\n")
	doc.WriteString("This document describes the standardized configurations for all frontend templates.\n\n")

	// Package.json standards
	doc.WriteString("## Package.json Standards\n\n")
	doc.WriteString("### Required Scripts\n\n")

	standards := GetFrontendStandards()
	for scriptName, command := range standards.PackageJSON.Scripts {
		doc.WriteString(fmt.Sprintf("- `%s`: `%s`\n", scriptName, command))
	}

	doc.WriteString("\n### Required Dependencies\n\n")
	for depName, version := range standards.PackageJSON.Dependencies {
		doc.WriteString(fmt.Sprintf("- `%s`: `%s`\n", depName, version))
	}

	doc.WriteString("\n### Required Dev Dependencies\n\n")
	for depName, version := range standards.PackageJSON.DevDeps {
		doc.WriteString(fmt.Sprintf("- `%s`: `%s`\n", depName, version))
	}

	doc.WriteString("\n### Engine Requirements\n\n")
	for engine, version := range standards.PackageJSON.Engines {
		doc.WriteString(fmt.Sprintf("- `%s`: `%s`\n", engine, version))
	}

	// Template-specific configurations
	doc.WriteString("\n## Template-Specific Configurations\n\n")

	templateTypes := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	for _, templateType := range templateTypes {
		doc.WriteString(fmt.Sprintf("### %s\n\n", templateType))

		if port, exists := standards.PackageJSON.Ports[templateType]; exists {
			doc.WriteString(fmt.Sprintf("- **Port**: %d\n", port))
		}

		specificDeps := GetTemplateSpecificDependencies(templateType)
		if len(specificDeps) > 0 {
			doc.WriteString("- **Additional Dependencies**:\n")
			for depName, version := range specificDeps {
				doc.WriteString(fmt.Sprintf("  - `%s`: `%s`\n", depName, version))
			}
		}
		doc.WriteString("\n")
	}

	// Configuration files
	doc.WriteString("## Configuration Files\n\n")
	doc.WriteString("All frontend templates must include the following standardized configuration files:\n\n")

	configFiles := []string{
		"package.json.tmpl",
		"tsconfig.json.tmpl",
		".eslintrc.json.tmpl",
		".prettierrc.tmpl",
		"vercel.json.tmpl",
		"tailwind.config.js.tmpl",
		"next.config.js.tmpl",
		"postcss.config.js.tmpl",
		"jest.config.js.tmpl",
		"jest.setup.js.tmpl",
	}

	for _, file := range configFiles {
		doc.WriteString(fmt.Sprintf("- `%s`\n", file))
	}

	doc.WriteString("\n## Vercel Deployment Standards\n\n")
	doc.WriteString("All frontend templates are configured for Vercel deployment with:\n\n")
	doc.WriteString("- **Framework**: Next.js\n")
	doc.WriteString("- **Build Command**: `npm run build`\n")
	doc.WriteString("- **Dev Command**: `npm run dev`\n")
	doc.WriteString("- **Install Command**: `npm install`\n")
	doc.WriteString("- **Security Headers**: Configured for production security\n")
	doc.WriteString("- **Environment Variables**: Standardized naming conventions\n\n")

	doc.WriteString("## Validation Rules\n\n")
	doc.WriteString("Templates are validated against the following rules:\n\n")
	doc.WriteString("1. **Package.json Consistency**: All required scripts, dependencies, and engines must be present\n")
	doc.WriteString("2. **TypeScript Configuration**: Standardized compiler options and path mappings\n")
	doc.WriteString("3. **ESLint Configuration**: Consistent linting rules across all templates\n")
	doc.WriteString("4. **Prettier Configuration**: Uniform code formatting settings\n")
	doc.WriteString("5. **Vercel Compatibility**: Proper deployment configuration\n")
	doc.WriteString("6. **Security Standards**: Required security headers and configurations\n\n")

	return doc.String()
}

// ExportValidationResults exports validation results to JSON
func (sm *StandardsManager) ExportValidationResults(results map[string]*ValidationResult, outputPath string) error {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal validation results: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write validation results to %s: %w", outputPath, err)
	}

	return nil
}

// CompareTemplateConfigurations compares configurations between templates
func (sm *StandardsManager) CompareTemplateConfigurations(templatesDir string) (string, error) {
	var report strings.Builder

	report.WriteString("# Template Configuration Comparison\n\n")

	frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	// Compare package.json files
	report.WriteString("## Package.json Comparison\n\n")

	packageConfigs := make(map[string]map[string]interface{})

	for _, templateType := range frontendTemplates {
		packagePath := filepath.Join(templatesDir, "frontend", templateType, "package.json.tmpl")

		if content, err := os.ReadFile(packagePath); err == nil {
			var pkg map[string]interface{}
			if err := json.Unmarshal(content, &pkg); err == nil {
				packageConfigs[templateType] = pkg
			}
		}
	}

	// Compare scripts
	if len(packageConfigs) > 0 {
		report.WriteString("### Scripts Comparison\n\n")
		report.WriteString("| Script | nextjs-app | nextjs-home | nextjs-admin |\n")
		report.WriteString("|--------|------------|-------------|---------------|\n")

		allScripts := make(map[string]bool)
		for _, pkg := range packageConfigs {
			if scripts, ok := pkg["scripts"].(map[string]interface{}); ok {
				for scriptName := range scripts {
					allScripts[scriptName] = true
				}
			}
		}

		for scriptName := range allScripts {
			report.WriteString(fmt.Sprintf("| %s |", scriptName))

			for _, templateType := range frontendTemplates {
				if pkg, exists := packageConfigs[templateType]; exists {
					if scripts, ok := pkg["scripts"].(map[string]interface{}); ok {
						if script, exists := scripts[scriptName]; exists {
							report.WriteString(fmt.Sprintf(" `%s` |", script))
						} else {
							report.WriteString(" ❌ |")
						}
					} else {
						report.WriteString(" ❌ |")
					}
				} else {
					report.WriteString(" ❌ |")
				}
			}
			report.WriteString("\n")
		}
		report.WriteString("\n")
	}

	return report.String(), nil
}
