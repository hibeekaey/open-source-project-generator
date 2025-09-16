package interfaces

import "github.com/open-source-template-generator/pkg/models"

// ValidationEngine defines the contract for basic project validation operations.
//
// The ValidationEngine interface provides essential validation capabilities for generated projects,
// templates, and configurations. It ensures that generated projects meet basic quality standards
// and are free from common configuration issues.
//
// The validation engine covers:
//   - Basic project structure validation
//   - Configuration file syntax validation
//   - Essential template validation
//
// Implementations should provide:
//   - Basic validation results with actionable feedback
//   - Simple validation rules for core functionality
type ValidationEngine interface {
	// ValidateProject validates the basic project structure.
	//
	// This method performs basic validation of a generated project including:
	//   - Directory structure and file organization
	//   - Configuration file syntax
	//   - Essential file presence
	//
	// Parameters:
	//   - projectPath: Path to the root directory of the generated project
	//
	// Returns:
	//   - *models.ValidationResult: Basic validation results with issues and suggestions
	//   - error: Any error that occurred during validation process
	ValidateProject(projectPath string) (*models.ValidationResult, error)

	// ValidatePackageJSON validates a package.json file for Node.js projects.
	//
	// This method checks:
	//   - JSON syntax correctness
	//   - Required fields presence (name, version, etc.)
	//   - Basic script definitions
	//
	// Parameters:
	//   - path: Path to the package.json file
	//
	// Returns:
	//   - error: Any validation error found in the package.json file
	ValidatePackageJSON(path string) error

	// ValidateGoMod validates a go.mod file for Go projects.
	//
	// This method checks:
	//   - Module declaration syntax
	//   - Basic dependency declarations
	//   - Module path correctness
	//
	// Parameters:
	//   - path: Path to the go.mod file
	//
	// Returns:
	//   - error: Any validation error found in the go.mod file
	ValidateGoMod(path string) error

	// ValidateDockerfile validates a Dockerfile for containerization.
	//
	// This method checks:
	//   - Dockerfile syntax and instruction validity
	//   - Basic structure and best practices
	//
	// Parameters:
	//   - path: Path to the Dockerfile
	//
	// Returns:
	//   - error: Any validation error found in the Dockerfile
	ValidateDockerfile(path string) error

	// ValidateYAML validates a YAML configuration file.
	//
	// This method checks:
	//   - YAML syntax correctness
	//   - Basic structure validation
	//
	// Parameters:
	//   - path: Path to the YAML file
	//
	// Returns:
	//   - error: Any validation error found in the YAML file
	ValidateYAML(path string) error

	// ValidateJSON validates a JSON configuration file.
	//
	// This method checks:
	//   - JSON syntax correctness
	//   - Basic structure validation
	//
	// Parameters:
	//   - path: Path to the JSON file
	//
	// Returns:
	//   - error: Any validation error found in the JSON file
	ValidateJSON(path string) error

	// ValidateTemplate validates a template file.
	//
	// This method checks:
	//   - Template syntax correctness
	//   - Basic template structure
	//
	// Parameters:
	//   - path: Path to the template file
	//
	// Returns:
	//   - error: Any validation error found in the template file
	ValidateTemplate(path string) error
}
