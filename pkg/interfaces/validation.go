package interfaces

import "github.com/open-source-template-generator/pkg/models"

// ValidationEngine defines the contract for project validation operations
type ValidationEngine interface {
	// ValidateProject validates the entire generated project structure
	ValidateProject(projectPath string) (*models.ValidationResult, error)

	// ValidatePackageJSON validates a package.json file
	ValidatePackageJSON(path string) error

	// ValidateGoMod validates a go.mod file
	ValidateGoMod(path string) error

	// ValidateDockerfile validates a Dockerfile
	ValidateDockerfile(path string) error

	// ValidateYAML validates a YAML configuration file
	ValidateYAML(path string) error

	// ValidateJSON validates a JSON configuration file
	ValidateJSON(path string) error
}
