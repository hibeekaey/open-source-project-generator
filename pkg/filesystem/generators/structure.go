package generators

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// StructureGenerator handles project directory structure generation
type StructureGenerator struct {
	fsOps FileSystemOperationsInterface
}

// NewStructureGenerator creates a new structure generator
func NewStructureGenerator(fsOps FileSystemOperationsInterface) *StructureGenerator {
	return &StructureGenerator{
		fsOps: fsOps,
	}
}

// GenerateProjectStructure creates the complete project directory structure
func (sg *StructureGenerator) GenerateProjectStructure(config *models.ProjectConfig, outputPath string) error {
	// Validate organization is provided
	if config.Organization == "" {
		return fmt.Errorf("organization cannot be empty")
	}

	// Create the root project directory using FileSystemOperations
	_, err := sg.fsOps.CreateProjectRoot(config, outputPath)
	if err != nil {
		return fmt.Errorf("failed to create project root: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates the generated project structure
func (sg *StructureGenerator) ValidateProjectStructure(projectPath string, config *models.ProjectConfig) error {
	// Validate project root using FileSystemOperations
	if err := sg.fsOps.ValidateProjectRoot(projectPath, config); err != nil {
		return fmt.Errorf("project root validation failed: %w", err)
	}

	// Validate that required configuration files exist
	requiredFiles := []string{
		"Makefile",
		"README.md",
		"docker-compose.yml",
		".gitignore",
	}

	if err := sg.fsOps.ValidateFileStructure(projectPath, requiredFiles); err != nil {
		return fmt.Errorf("required root files validation failed: %w", err)
	}

	return nil
}
