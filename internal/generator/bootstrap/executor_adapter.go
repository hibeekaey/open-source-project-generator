package bootstrap

import (
	"context"
	"fmt"
	"io"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ExecutorAdapter wraps a BaseExecutor to implement the BootstrapExecutorInterface
// This adapter handles the interface{} conversions required by the interface
type ExecutorAdapter struct {
	executor interface {
		Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error)
		ExecuteWithStreaming(ctx context.Context, spec *BootstrapSpec, output io.Writer) (*models.ExecutionResult, error)
		SupportsComponent(componentType string) bool
		GetDefaultFlags(componentType string) []string
		ValidateConfig(config map[string]interface{}) error
	}
}

// NewExecutorAdapter creates a new executor adapter
func NewExecutorAdapter(executor interface {
	Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error)
	ExecuteWithStreaming(ctx context.Context, spec *BootstrapSpec, output io.Writer) (*models.ExecutionResult, error)
	SupportsComponent(componentType string) bool
	GetDefaultFlags(componentType string) []string
	ValidateConfig(config map[string]interface{}) error
}) *ExecutorAdapter {
	return &ExecutorAdapter{executor: executor}
}

// Execute implements the BootstrapExecutorInterface with interface{} parameters
func (ea *ExecutorAdapter) Execute(ctx context.Context, spec interface{}) (interface{}, error) {
	bootstrapSpec, ok := spec.(*BootstrapSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type: expected *BootstrapSpec, got %T", spec)
	}

	return ea.executor.Execute(ctx, bootstrapSpec)
}

// ExecuteWithStreaming implements the BootstrapExecutorInterface with interface{} parameters
func (ea *ExecutorAdapter) ExecuteWithStreaming(ctx context.Context, spec interface{}, output interface{}) (interface{}, error) {
	bootstrapSpec, ok := spec.(*BootstrapSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type: expected *BootstrapSpec, got %T", spec)
	}

	writer, ok := output.(io.Writer)
	if !ok {
		return nil, fmt.Errorf("invalid output type: expected io.Writer, got %T", output)
	}

	return ea.executor.ExecuteWithStreaming(ctx, bootstrapSpec, writer)
}

// SupportsComponent delegates to the wrapped executor
func (ea *ExecutorAdapter) SupportsComponent(componentType string) bool {
	return ea.executor.SupportsComponent(componentType)
}

// GetDefaultFlags delegates to the wrapped executor
func (ea *ExecutorAdapter) GetDefaultFlags(componentType string) []string {
	return ea.executor.GetDefaultFlags(componentType)
}

// ValidateConfig delegates to the wrapped executor
func (ea *ExecutorAdapter) ValidateConfig(config map[string]interface{}) error {
	return ea.executor.ValidateConfig(config)
}
