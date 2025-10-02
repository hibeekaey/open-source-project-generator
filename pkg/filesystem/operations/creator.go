package operations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Creator handles file and directory creation operations
type Creator struct {
	fsOps FileSystemOperationsInterface
}

// NewCreator creates a new file creator
func NewCreator(fsOps FileSystemOperationsInterface) *Creator {
	return &Creator{
		fsOps: fsOps,
	}
}

// CreateProjectRoot creates the root project directory
func (c *Creator) CreateProjectRoot(config *models.ProjectConfig, outputPath string) (string, error) {
	return c.fsOps.CreateProjectRoot(config, outputPath)
}

// CreateFile creates a file with the given content
func (c *Creator) CreateFile(path string, content []byte, perm os.FileMode) error {
	return c.fsOps.WriteFile(path, content, perm)
}

// CreateConfigurationFile creates a configuration file with validation
func (c *Creator) CreateConfigurationFile(projectPath, filename, content string) error {
	filePath := filepath.Join(projectPath, filename)
	if err := c.fsOps.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create %s: %w", filename, err)
	}
	return nil
}
