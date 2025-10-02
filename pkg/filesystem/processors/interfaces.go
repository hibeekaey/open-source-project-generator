package processors

import (
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FileSystemOperationsInterface defines the interface for file system operations
type FileSystemOperationsInterface interface {
	CreateProjectRoot(config *models.ProjectConfig, outputPath string) (string, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	FileExists(path string) bool
	ValidateFileContent(path string) error
	ValidateProjectRoot(projectPath string, config *models.ProjectConfig) error
	ValidateFileStructure(projectPath string, requiredFiles []string) error
}
