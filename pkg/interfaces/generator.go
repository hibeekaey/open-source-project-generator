package interfaces

import (
	"context"
)

// GeneratorInterface defines the contract for project generation components
type GeneratorInterface interface {
	// Generate creates a project component based on the provided specification
	Generate(ctx context.Context, spec interface{}) (interface{}, error)

	// SupportsComponent checks if this generator can handle the given component type
	SupportsComponent(componentType string) bool

	// Validate validates the generation specification before execution
	Validate(spec interface{}) error
}

// BootstrapExecutorInterface defines the contract for executing external CLI tools
type BootstrapExecutorInterface interface {
	// Execute runs an external bootstrap tool with the given specification
	Execute(ctx context.Context, spec interface{}) (interface{}, error)

	// SupportsComponent checks if this executor can handle the given component type
	SupportsComponent(componentType string) bool

	// GetDefaultFlags returns default CLI flags for the component type
	GetDefaultFlags(componentType string) []string
}

// FallbackGeneratorInterface defines the contract for custom project generation
type FallbackGeneratorInterface interface {
	// Generate creates a project component using custom templates
	Generate(ctx context.Context, spec interface{}) (interface{}, error)

	// SupportsComponent checks if this generator can handle the given component type
	SupportsComponent(componentType string) bool

	// GetRequiredManualSteps returns manual steps needed after generation
	GetRequiredManualSteps(componentType string) []string
}

// ProjectCoordinatorInterface defines the contract for orchestrating project generation
type ProjectCoordinatorInterface interface {
	// Generate orchestrates the complete project generation workflow
	Generate(ctx context.Context, config interface{}) (interface{}, error)

	// DryRun previews what would be generated without creating files
	DryRun(ctx context.Context, config interface{}) (interface{}, error)

	// Validate validates the project configuration
	Validate(config interface{}) error
}
