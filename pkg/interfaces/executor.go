package interfaces

import (
	"context"
	"io"
)

// ToolDiscoveryInterface defines the contract for discovering and validating bootstrap tools
type ToolDiscoveryInterface interface {
	// IsAvailable checks if a tool is available in the system PATH
	IsAvailable(toolName string) (bool, error)

	// GetVersion retrieves the installed version of a tool
	GetVersion(toolName string) (string, error)

	// CheckRequirements validates that all required tools are available
	CheckRequirements(tools []string) (interface{}, error)

	// GetInstallInstructions returns OS-specific installation instructions
	GetInstallInstructions(toolName string, os string) string
}

// ToolExecutorInterface defines the contract for executing external tools
type ToolExecutorInterface interface {
	// Execute runs a tool with the given arguments
	Execute(ctx context.Context, toolName string, args []string, workDir string) (interface{}, error)

	// ExecuteWithStreaming runs a tool and streams output to the provided writer
	ExecuteWithStreaming(ctx context.Context, toolName string, args []string, workDir string, output io.Writer) error

	// ValidateTool checks if a tool is whitelisted and safe to execute
	ValidateTool(toolName string, args []string) error
}

// IntegrationManagerInterface defines the contract for integrating generated components
type IntegrationManagerInterface interface {
	// Integrate configures generated components to work together
	Integrate(ctx context.Context, components interface{}) error

	// GenerateDockerCompose creates a Docker Compose file for all components
	GenerateDockerCompose(components interface{}) (string, error)

	// ConfigureEnvironment sets up environment variables and configuration
	ConfigureEnvironment(components interface{}) error

	// GenerateScripts creates build and run scripts for the project
	GenerateScripts(components interface{}) error
}
