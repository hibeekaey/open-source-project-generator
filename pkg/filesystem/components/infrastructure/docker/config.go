package docker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FileSystemOperations interface for file operations
type FileSystemOperations interface {
	WriteFile(path string, content []byte, perm os.FileMode) error
	EnsureDirectory(path string) error
}

// ConfigGenerator handles Docker configuration generation
type ConfigGenerator struct {
	fsOps         FileSystemOperations
	dockerfileGen *DockerfileGenerator
	composeGen    *ComposeGenerator
}

// NewConfigGenerator creates a new Docker config generator
func NewConfigGenerator(fsOps FileSystemOperations) *ConfigGenerator {
	return &ConfigGenerator{
		fsOps:         fsOps,
		dockerfileGen: NewDockerfileGenerator(),
		composeGen:    NewComposeGenerator(),
	}
}

// GenerateDockerFiles creates all Docker configuration files
func (cg *ConfigGenerator) GenerateDockerFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate docker-compose.yml
	dockerComposeContent := cg.composeGen.GenerateDockerCompose(config)
	dockerComposePath := filepath.Join(projectPath, "Deploy/docker/docker-compose.yml")
	if err := cg.fsOps.WriteFile(dockerComposePath, []byte(dockerComposeContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Generate docker-compose.dev.yml
	dockerComposeDevContent := cg.composeGen.GenerateDockerComposeDev(config)
	dockerComposeDevPath := filepath.Join(projectPath, "Deploy/docker/docker-compose.dev.yml")
	if err := cg.fsOps.WriteFile(dockerComposeDevPath, []byte(dockerComposeDevContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.dev.yml: %w", err)
	}

	// Generate docker-compose.prod.yml
	dockerComposeProdContent := cg.composeGen.GenerateDockerComposeProd(config)
	dockerComposeProdPath := filepath.Join(projectPath, "Deploy/docker/docker-compose.prod.yml")
	if err := cg.fsOps.WriteFile(dockerComposeProdPath, []byte(dockerComposeProdContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.prod.yml: %w", err)
	}

	// Generate .dockerignore
	dockerIgnoreContent := cg.dockerfileGen.GenerateDockerIgnore(config)
	dockerIgnorePath := filepath.Join(projectPath, "Deploy/docker/.dockerignore")
	if err := cg.fsOps.WriteFile(dockerIgnorePath, []byte(dockerIgnoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .dockerignore: %w", err)
	}

	// Generate Dockerfile for frontend
	frontendDockerfileContent := cg.dockerfileGen.GenerateFrontendDockerfile(config)
	frontendDockerfilePath := filepath.Join(projectPath, "Deploy/docker/Dockerfile.frontend")
	if err := cg.fsOps.WriteFile(frontendDockerfilePath, []byte(frontendDockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile.frontend: %w", err)
	}

	// Generate Dockerfile for backend
	backendDockerfileContent := cg.dockerfileGen.GenerateBackendDockerfile(config)
	backendDockerfilePath := filepath.Join(projectPath, "Deploy/docker/Dockerfile.backend")
	if err := cg.fsOps.WriteFile(backendDockerfilePath, []byte(backendDockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile.backend: %w", err)
	}

	return nil
}
