package interfaces

import (
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FileSystemGenerator defines the contract for file system operations
type FileSystemGenerator interface {
	// CreateProject creates the complete project structure based on configuration
	CreateProject(config *models.ProjectConfig, outputPath string) error

	// CreateDirectory creates a directory with proper permissions
	CreateDirectory(path string) error

	// WriteFile writes content to a file with specified permissions
	WriteFile(path string, content []byte, perm os.FileMode) error

	// CopyAssets copies binary assets from source to destination
	CopyAssets(srcDir, destDir string) error

	// CreateSymlink creates a symbolic link
	CreateSymlink(target, link string) error

	// FileExists checks if a file exists at the given path
	FileExists(path string) bool

	// EnsureDirectory ensures a directory exists, creating it if necessary
	EnsureDirectory(path string) error
}

// StandardizedStructureGenerator defines the contract for standardized directory structure generation
type StandardizedStructureGenerator interface {
	// GenerateStandardizedStructure creates the complete standardized project structure
	GenerateStandardizedStructure(config *models.ProjectConfig, outputPath string) error

	// CreateFrontendDirectoryStructure creates App/ directory with main/, home/, admin/, shared-components/ subdirectories
	CreateFrontendDirectoryStructure(projectPath string, config *models.ProjectConfig) error

	// CreateBackendDirectoryStructure creates CommonServer/ directory with cmd/, internal/, pkg/, migrations/, docs/ structure
	CreateBackendDirectoryStructure(projectPath string, config *models.ProjectConfig) error

	// CreateMobileDirectoryStructure creates Mobile/ directory with android/, ios/, shared/ subdirectories
	CreateMobileDirectoryStructure(projectPath string, config *models.ProjectConfig) error

	// CreateInfrastructureDirectoryStructure creates Deploy/ directory with docker/, k8s/, terraform/, monitoring/ subdirectories
	CreateInfrastructureDirectoryStructure(projectPath string, config *models.ProjectConfig) error

	// CreateCommonDirectoryStructure creates Docs/, Scripts/, and .github/ directories with appropriate content
	CreateCommonDirectoryStructure(projectPath string, config *models.ProjectConfig) error

	// GenerateStandardProjectFiles generates standard project files (README.md, CONTRIBUTING.md, LICENSE, .gitignore, Makefile)
	GenerateStandardProjectFiles(projectPath string, config *models.ProjectConfig) error
}
