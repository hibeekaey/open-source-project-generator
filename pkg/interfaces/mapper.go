package interfaces

import (
	"context"
)

// StructureMapperInterface defines the contract for mapping generated outputs to target structure
type StructureMapperInterface interface {
	// Map relocates generated files from source to target directory structure
	Map(ctx context.Context, source string, target string, componentType string) error

	// ValidateStructure verifies that the target structure is correct
	ValidateStructure(rootDir string) error

	// ValidateStructureWithComponents validates structure for specific components
	ValidateStructureWithComponents(rootDir string, componentTypes []string) error

	// GetTargetPath returns the target path for a given component type
	GetTargetPath(componentType string) string

	// UpdateReferences updates import paths and references after relocation
	UpdateReferences(ctx context.Context, rootDir string, componentType string) error
}

// ComponentMapperInterface defines the contract for component-to-directory mappings
type ComponentMapperInterface interface {
	// GetMapping returns the directory mapping for a component type
	GetMapping(componentType string) (string, error)

	// RegisterMapping adds a new component-to-directory mapping
	RegisterMapping(componentType string, targetPath string) error

	// ListMappings returns all registered component mappings
	ListMappings() map[string]string
}
