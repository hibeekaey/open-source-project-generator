package operations

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Validator handles file and project validation operations
type Validator struct {
	fsOps FileSystemOperationsInterface
}

// NewValidator creates a new validator
func NewValidator(fsOps FileSystemOperationsInterface) *Validator {
	return &Validator{
		fsOps: fsOps,
	}
}

// ValidateProjectRoot validates the project root directory
func (v *Validator) ValidateProjectRoot(projectPath string, config *models.ProjectConfig) error {
	return v.fsOps.ValidateProjectRoot(projectPath, config)
}

// ValidateRequiredFiles validates that required files exist
func (v *Validator) ValidateRequiredFiles(projectPath string, requiredFiles []string) error {
	return v.fsOps.ValidateFileStructure(projectPath, requiredFiles)
}

// ValidateFileContent validates file content
func (v *Validator) ValidateFileContent(filePath string) error {
	if !v.fsOps.FileExists(filePath) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	return v.fsOps.ValidateFileContent(filePath)
}

// FileExists checks if a file exists
func (v *Validator) FileExists(filePath string) bool {
	return v.fsOps.FileExists(filePath)
}
