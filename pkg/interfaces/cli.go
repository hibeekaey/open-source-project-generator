package interfaces

import "github.com/open-source-template-generator/pkg/models"

// CLIInterface defines the contract for CLI operations
type CLIInterface interface {
	// Run executes the CLI application
	Run() error

	// PromptProjectDetails collects project configuration from user input
	PromptProjectDetails() (*models.ProjectConfig, error)

	// SelectComponents allows user to choose which components to include
	SelectComponents() ([]string, error)

	// ConfirmGeneration shows configuration preview and asks for confirmation
	ConfirmGeneration(*models.ProjectConfig) bool
}
