// Package app provides the core application logic for the Open Source Template Generator.
//
// This package implements the main application structure, CLI command handling,
// and orchestrates the interaction between different components like template
// processing, validation, and project generation.
//
// The application follows clean architecture principles with dependency injection
// to ensure testability and maintainability.
package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-source-template-generator/internal/config"
	"github.com/open-source-template-generator/pkg/cli"
	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
	"github.com/open-source-template-generator/pkg/validation"
	"github.com/open-source-template-generator/pkg/version"
	"github.com/spf13/cobra"
)

// App represents the main application instance that orchestrates all CLI operations.
// It manages the basic components, CLI interface, and error handling.
//
// The App struct serves as the central coordinator for:
//   - CLI command processing and routing
//   - Basic component initialization
//   - Project generation workflows
//   - Basic validation operations
//   - Configuration management
type App struct {
	configManager  interfaces.ConfigManager
	validator      interfaces.ValidationEngine
	cli            interfaces.CLIInterface
	generator      interfaces.FileSystemGenerator
	templateEngine interfaces.TemplateEngine
	versionManager interfaces.VersionManager
}

// NewApp creates a new application instance with all required dependencies.
//
// Returns:
//   - *App: New application instance ready for use
//   - error: Any error that occurred during initialization
func NewApp() (*App, error) {
	// Initialize basic components
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	generator := filesystem.NewGenerator()
	templateEngine := template.NewEngine()
	versionManager := version.NewManager()

	// Initialize CLI
	cli := cli.NewCLI(configManager, validator)

	return &App{
		configManager:  configManager,
		validator:      validator,
		cli:            cli,
		generator:      generator,
		templateEngine: templateEngine,
		versionManager: versionManager,
	}, nil
}

// Run starts the application and processes command-line arguments.
//
// Parameters:
//   - args: Command-line arguments (typically os.Args[1:])
//
// Returns:
//   - error: Any error that occurred during execution
func (a *App) Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:   "generator",
		Short: "Open Source Template Generator",
		Long:  "A tool for generating production-ready project templates",
	}

	// Add generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from templates",
		Long:  "Generate a new project using interactive prompts to configure components and settings",
		RunE:  a.runGenerate,
	}
	rootCmd.AddCommand(generateCmd)

	// Add help command
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Show help information",
		Long:  "Display help information for the generator tool",
		RunE:  a.runHelp,
	}
	rootCmd.AddCommand(helpCmd)

	// Add version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the version of the generator tool",
		RunE:  a.runVersion,
	}
	rootCmd.AddCommand(versionCmd)

	// Set the arguments
	rootCmd.SetArgs(args)

	// Execute the command
	return rootCmd.Execute()
}

// runGenerate handles the generate command
func (a *App) runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Println("Starting project generation...")

	// Use CLI to collect project configuration
	config, err := a.cli.PromptProjectDetails()
	if err != nil {
		return fmt.Errorf("failed to collect project configuration: %w", err)
	}

	// Generate the project
	if err := a.generateProject(config); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Println("Project generated successfully!")
	return nil
}

// runHelp handles the help command
func (a *App) runHelp(cmd *cobra.Command, args []string) error {
	fmt.Println("Open Source Template Generator")
	fmt.Println("A tool for generating production-ready project templates")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  generate    Generate a new project from templates")
	fmt.Println("  help        Show this help information")
	fmt.Println("  version     Show version information")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  generator generate    # Start interactive project generation")
	fmt.Println("  generator help        # Show this help")
	fmt.Println("  generator version     # Show version")
	return nil
}

// runVersion handles the version command
func (a *App) runVersion(cmd *cobra.Command, args []string) error {
	fmt.Println("Open Source Template Generator v1.0.0")
	return nil
}

// generateProject generates a project based on the provided configuration
func (a *App) generateProject(config *models.ProjectConfig) error {
	// Create the project directory structure
	if err := a.generator.CreateProject(config, config.OutputPath); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Process templates
	if err := a.processTemplates(config); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}

	// Basic validation
	if err := a.validateProject(config.OutputPath); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	return nil
}

// processTemplates processes all templates for the project
func (a *App) processTemplates(config *models.ProjectConfig) error {
	// Get template directory (assuming templates are in ./templates)
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
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		components = append(components, "mobile")
	}
	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		components = append(components, "infrastructure")
	}

	for _, component := range components {
		componentDir := filepath.Join(templateDir, component)
		if _, err := os.Stat(componentDir); os.IsNotExist(err) {
			continue // Skip if component directory doesn't exist
		}

		outputDir := filepath.Join(config.OutputPath, config.Name)
		if err := a.templateEngine.ProcessDirectory(componentDir, outputDir, config); err != nil {
			return fmt.Errorf("failed to process %s templates: %w", component, err)
		}
	}

	return nil
}

// validateProject performs basic validation on the generated project
func (a *App) validateProject(projectPath string) error {
	result, err := a.validator.ValidateProject(projectPath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid {
		fmt.Println("Project validation found issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s: %s\n", issue.Type, issue.Message)
		}
	}

	return nil
}
