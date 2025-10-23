package interfaces

import (
	"context"
	"io"

	"github.com/cuesoftinc/open-source-project-generator/internal/generator/bootstrap"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// BootstrapExecutorInterface defines the enhanced interface for executing bootstrap tools
type BootstrapExecutorInterface interface {
	// Execute runs the bootstrap tool
	Execute(ctx context.Context, spec *bootstrap.BootstrapSpec) (*models.ExecutionResult, error)

	// ExecuteWithStreaming runs with real-time output
	ExecuteWithStreaming(ctx context.Context, spec *bootstrap.BootstrapSpec, output io.Writer) (*models.ExecutionResult, error)

	// SupportsComponent checks if this executor handles the component type
	SupportsComponent(componentType string) bool

	// GetDefaultFlags returns default CLI flags for the component type
	GetDefaultFlags(componentType string) []string

	// ValidateConfig validates component-specific configuration
	ValidateConfig(config map[string]interface{}) error
}
