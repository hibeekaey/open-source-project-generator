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
