package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
	"github.com/open-source-template-generator/pkg/validation"
)

func main() {
	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "A test project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   false,
					Admin:  false,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"react":      "18.2.0",
				"next":       "13.4.0",
				"typescript": "5.0.0",
			},
		},
		OutputPath: "./test-output",
	}

	// Create output directory
	if err := os.MkdirAll(config.OutputPath, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Initialize components
	generator := filesystem.NewGenerator()
	templateEngine := template.NewEngine()
	validator := validation.NewEngine()

	fmt.Println("Testing core generator functionality...")

	// Test 1: Create project structure
	fmt.Println("1. Creating project structure...")
	if err := generator.CreateProject(config, config.OutputPath); err != nil {
		fmt.Printf("Failed to create project structure: %v\n", err)
		return
	}
	fmt.Println("✅ Project structure created")

	// Test 2: Process templates
	fmt.Println("2. Processing templates...")
	if err := processTemplates(templateEngine, config); err != nil {
		fmt.Printf("Failed to process templates: %v\n", err)
		return
	}
	fmt.Println("✅ Templates processed")

	// Test 3: Validate project
	fmt.Println("3. Validating project...")
	projectPath := filepath.Join(config.OutputPath, config.Name)
	result, err := validator.ValidateProject(projectPath)
	if err != nil {
		fmt.Printf("Failed to validate project: %v\n", err)
		return
	}

	if !result.Valid {
		fmt.Println("⚠️  Project validation found issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s: %s\n", issue.Type, issue.Message)
		}
	} else {
		fmt.Println("✅ Project validation passed")
	}

	fmt.Println("\n🎉 Generator test completed successfully!")
	fmt.Printf("Generated project at: %s\n", projectPath)
}

func processTemplates(templateEngine interfaces.TemplateEngine, config *models.ProjectConfig) error {
	// Get template directory
	templateDir := "./templates"
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found: %s", templateDir)
	}

	// Process each component's templates
	components := []string{"base"}

	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		components = append(components, "frontend")
	}
	if config.Components.Backend.GoGin {
		components = append(components, "backend")
	}

	for _, component := range components {
		componentDir := filepath.Join(templateDir, component)
		if _, err := os.Stat(componentDir); os.IsNotExist(err) {
			continue // Skip if component directory doesn't exist
		}

		outputDir := filepath.Join(config.OutputPath, config.Name)
		if err := templateEngine.ProcessDirectory(componentDir, outputDir, config); err != nil {
			return fmt.Errorf("failed to process %s templates: %w", component, err)
		}
	}

	return nil
}
