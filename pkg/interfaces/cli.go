// Package interfaces defines the core contracts and interfaces for the
// Open Source Template Generator components.
//
// This package contains interface definitions that enable dependency injection,
// testing, and modular architecture throughout the application.
package interfaces

import "github.com/open-source-template-generator/pkg/models"

// CLIInterface defines the contract for command-line interface operations.
//
// This interface abstracts CLI functionality to enable testing and different
// CLI implementations. It covers the complete user interaction workflow from
// project configuration collection to generation confirmation.
//
// Implementations should provide:
//   - Interactive prompts for project configuration
//   - Component selection with validation
//   - Configuration preview and confirmation
//   - Progress indication and user feedback
type CLIInterface interface {
	// Run executes the CLI application with command-line arguments.
	//
	// This method starts the CLI application and handles the complete
	// user interaction workflow. It should process command-line flags,
	// route to appropriate subcommands, and manage the overall CLI lifecycle.
	//
	// Returns:
	//   - error: Any error that occurred during CLI execution
	Run() error

	// PromptProjectDetails collects comprehensive project configuration from user input.
	//
	// This method guides users through interactive prompts to collect all
	// necessary project information including basic details, component selection,
	// and configuration options. It should validate input and ensure the
	// configuration is complete and valid.
	//
	// Returns:
	//   - *models.ProjectConfig: Complete project configuration
	//   - error: Any error that occurred during configuration collection
	PromptProjectDetails() (*models.ProjectConfig, error)

	// SelectComponents allows users to choose which components to include in their project.
	//
	// This method presents available components (frontend, backend, mobile,
	// infrastructure) and allows users to select which ones to include.
	// It should validate component dependencies and provide warnings for
	// potentially problematic combinations.
	//
	// Returns:
	//   - []string: List of selected component identifiers
	//   - error: Any error that occurred during component selection
	SelectComponents() ([]string, error)

	// ConfirmGeneration shows a configuration preview and asks for user confirmation.
	//
	// This method displays a comprehensive summary of the project configuration
	// including selected components, versions, and output path. It should give
	// users a final opportunity to review and confirm before generation begins.
	//
	// Parameters:
	//   - config: Complete project configuration to preview
	//
	// Returns:
	//   - bool: true if user confirms generation, false if cancelled
	ConfirmGeneration(*models.ProjectConfig) bool
}
