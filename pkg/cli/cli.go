// Package cli provides basic command-line interface functionality for the
// Open Source Template Generator.
//
// This package handles essential user interactions including:
//   - Basic project configuration collection
//   - Simple component selection
//   - Configuration confirmation
//   - Basic user feedback
//
// The CLI provides a simple, focused experience for configuring and generating projects.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// CLI implements the CLIInterface for basic user operations.
//
// The CLI struct provides methods for:
//   - Collecting basic project configuration through simple prompts
//   - Validating essential user input
//   - Displaying basic progress and feedback
//   - Confirming project configuration before generation
type CLI struct {
	configManager    interfaces.ConfigManager    // Manages configuration and defaults
	validator        interfaces.ValidationEngine // Validates user input and project structure
	generatorVersion string                      // Generator version for project metadata
}

// NewCLI creates a new CLI instance with the provided dependencies.
//
// Parameters:
//   - configManager: Handles configuration loading and validation
//   - validator: Provides input validation
//
// Returns:
//   - *CLI: New CLI instance ready for use
func NewCLI(configManager interfaces.ConfigManager, validator interfaces.ValidationEngine) interfaces.CLIInterface {
	return &CLI{
		configManager:    configManager,
		validator:        validator,
		generatorVersion: "1.0.0", // Default version
	}
}

// Run executes the CLI application with command-line arguments.
func (c *CLI) Run() error {
	// Basic CLI implementation - just collect config and generate
	config, err := c.PromptProjectDetails()
	if err != nil {
		return fmt.Errorf("failed to collect project details: %w", err)
	}

	if !c.ConfirmGeneration(config) {
		fmt.Println("Project generation canceled.")
		return nil
	}

	fmt.Println("Project generation completed successfully!")
	return nil
}

// PromptProjectDetails collects basic project configuration from user input.
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	fmt.Println("Welcome to the Open Source Template Generator!")
	fmt.Println("Let's set up your project...")

	config := &models.ProjectConfig{}

	// Prompt for project name
	var projectName string
	prompt := &survey.Input{
		Message: "What is your project name?",
	}
	if err := survey.AskOne(prompt, &projectName); err != nil {
		return nil, fmt.Errorf("failed to get project name: %w", err)
	}
	config.Name = projectName

	// Prompt for organization
	var organization string
	prompt = &survey.Input{
		Message: "What is your organization name?",
		Default: "my-org",
	}
	if err := survey.AskOne(prompt, &organization); err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	config.Organization = organization

	// Prompt for components
	var selectedComponents []string
	componentPrompt := &survey.MultiSelect{
		Message: "Which components would you like to include?",
		Options: []string{"frontend", "backend", "mobile", "infrastructure"},
		Default: []string{"frontend", "backend"},
	}
	if err := survey.AskOne(componentPrompt, &selectedComponents); err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	// Set components
	config.Components = models.Components{}
	for _, component := range selectedComponents {
		switch component {
		case "frontend":
			config.Components.Frontend = models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   true,
					Admin:  false,
					Shared: true,
				},
			}
		case "backend":
			config.Components.Backend = models.BackendComponents{
				GoGin: true,
			}
		case "mobile":
			config.Components.Mobile = models.MobileComponents{
				Android: true,
				IOS:     true,
			}
		case "infrastructure":
			config.Components.Infrastructure = models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			}
		}
	}

	// Prompt for output path
	var outputPath string
	prompt = &survey.Input{
		Message: "Where should the project be generated?",
		Default: "./output",
	}
	if err := survey.AskOne(prompt, &outputPath); err != nil {
		return nil, fmt.Errorf("failed to get output path: %w", err)
	}

	// Validate and create output path
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve output path: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0750); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	config.OutputPath = absPath

	// Set basic version information
	config.Versions = &models.VersionConfig{
		Node: "20.0.0",
		Go:   "1.21.0",
		Packages: map[string]string{
			"react":      "18.2.0",
			"next":       "13.4.0",
			"typescript": "5.0.0",
		},
	}

	return config, nil
}

// ConfirmGeneration shows a basic configuration preview and asks for user confirmation.
func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Project Configuration Preview")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Project Name: %s\n", config.Name)
	fmt.Printf("Organization: %s\n", config.Organization)
	fmt.Printf("Output Path: %s\n", config.OutputPath)

	fmt.Println("\nSelected Components:")
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		fmt.Println("  - Frontend (Next.js)")
	}
	if config.Components.Backend.GoGin {
		fmt.Println("  - Backend (Go + Gin)")
	}
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		fmt.Println("  - Mobile (Android + iOS)")
	}
	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		fmt.Println("  - Infrastructure (Docker + K8s + Terraform)")
	}

	fmt.Println(strings.Repeat("=", 50))

	var confirm bool
	prompt := &survey.Confirm{
		Message: "Do you want to generate this project?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		fmt.Printf("Error getting confirmation: %v\n", err)
		return false
	}

	return confirm
}
