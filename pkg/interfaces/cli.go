// Package interfaces defines the core contracts and interfaces for the
// Open Source Template Generator components.
//
// This package contains interface definitions that enable dependency injection,
// testing, and modular architecture throughout the application.
package interfaces

import "github.com/cuesoftinc/open-source-project-generator/pkg/models"

// CLIInterface defines the contract for command-line interface operations.
//
// This interface abstracts CLI functionality to enable testing and different
// CLI implementations. It covers basic user interaction workflow for
// project configuration collection and generation.
//
// Implementations should provide:
//   - Basic prompts for project configuration
//   - Simple component selection
//   - Configuration confirmation
//   - Basic user feedback
type CLIInterface interface {
	// Run executes the CLI application with command-line arguments.
	//
	// This method starts the CLI application and handles the basic
	// user interaction workflow. It should process command-line flags,
	// route to appropriate subcommands, and manage the overall CLI lifecycle.
	//
	// Returns:
	//   - error: Any error that occurred during CLI execution
	Run() error

	// PromptProjectDetails collects basic project configuration from user input.
	//
	// This method guides users through simple prompts to collect essential
	// project information including basic details and component selection.
	// It should validate input and ensure the configuration is complete.
	//
	// Returns:
	//   - *models.ProjectConfig: Complete project configuration
	//   - error: Any error that occurred during configuration collection
	PromptProjectDetails() (*models.ProjectConfig, error)

	// ConfirmGeneration shows a basic configuration preview and asks for user confirmation.
	//
	// This method displays a simple summary of the project configuration
	// including selected components and output path. It should give
	// users a final opportunity to review and confirm before generation begins.
	//
	// Parameters:
	//   - config: Complete project configuration to preview
	//
	// Returns:
	//   - bool: true if user confirms generation, false if canceled
	ConfirmGeneration(*models.ProjectConfig) bool
}
