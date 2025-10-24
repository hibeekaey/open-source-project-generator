package fallback

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// GeneratorInterface defines the interface for fallback project generation
type GeneratorInterface interface {
	// Generate creates a project component using custom templates
	Generate(ctx context.Context, spec *models.FallbackSpec) (*models.ComponentResult, error)

	// SupportsComponent checks if this generator can handle the given component type
	SupportsComponent(componentType string) bool

	// GetRequiredManualSteps returns manual steps needed after generation
	GetRequiredManualSteps(componentType string) []string
}

// Registry manages fallback generators for different component types
type Registry struct {
	generators map[string]GeneratorInterface
}

// NewRegistry creates a new fallback generator registry
func NewRegistry() *Registry {
	return &Registry{
		generators: make(map[string]GeneratorInterface),
	}
}

// Register adds a generator for a specific component type
func (r *Registry) Register(componentType string, generator GeneratorInterface) {
	r.generators[componentType] = generator
}

// Get retrieves a generator for the specified component type
func (r *Registry) Get(componentType string) (GeneratorInterface, error) {
	gen, exists := r.generators[componentType]
	if !exists {
		return nil, fmt.Errorf("no fallback generator registered for component type: %s", componentType)
	}
	return gen, nil
}

// Supports checks if a fallback generator exists for the component type
func (r *Registry) Supports(componentType string) bool {
	_, exists := r.generators[componentType]
	return exists
}

// GetSupportedTypes returns all component types with registered fallback generators
func (r *Registry) GetSupportedTypes() []string {
	types := make([]string, 0, len(r.generators))
	for componentType := range r.generators {
		types = append(types, componentType)
	}
	return types
}

// DefaultRegistry returns a registry with all standard fallback generators registered
func DefaultRegistry() *Registry {
	registry := NewRegistry()

	// Register Android fallback generator
	registry.Register("android", NewAndroidGenerator())

	// Register iOS fallback generator
	registry.Register("ios", NewIOSGenerator())

	return registry
}
