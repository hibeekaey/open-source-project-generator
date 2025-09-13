package interfaces

import "github.com/open-source-template-generator/pkg/models"

// ValidationEngine defines the contract for comprehensive project validation operations.
//
// The ValidationEngine interface provides validation capabilities for generated projects,
// templates, and configurations. It ensures that generated projects meet quality standards,
// follow best practices, and are free from common issues.
//
// The validation engine covers:
//   - Project structure and file organization validation
//   - Configuration file syntax and semantic validation
//   - Cross-template consistency and compatibility checks
//   - Security vulnerability scanning
//   - Version compatibility validation
//   - Platform-specific deployment validation
//
// Implementations should provide:
//   - Detailed validation results with actionable feedback
//   - Performance optimization for large projects
//   - Extensible validation rules and custom validators
//   - Integration with external validation tools and services
type ValidationEngine interface {
	// ValidateProject validates the entire generated project structure.
	//
	// This method performs comprehensive validation of a generated project including:
	//   - Directory structure and file organization
	//   - Configuration file syntax and semantics
	//   - Cross-file consistency and dependencies
	//   - Build system configuration
	//   - Security best practices compliance
	//
	// Parameters:
	//   - projectPath: Path to the root directory of the generated project
	//
	// Returns:
	//   - *models.ValidationResult: Detailed validation results with issues and suggestions
	//   - error: Any error that occurred during validation process
	ValidateProject(projectPath string) (*models.ValidationResult, error)

	// ValidatePackageJSON validates a package.json file for Node.js projects.
	//
	// This method checks:
	//   - JSON syntax correctness
	//   - Required fields presence (name, version, etc.)
	//   - Dependency version compatibility
	//   - Script definitions and validity
	//   - Security vulnerability scanning
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
	//   - Go version compatibility
	//   - Dependency declarations and versions
	//   - Replace directives validity
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
	//   - Base image security and best practices
	//   - Layer optimization and caching strategies
	//   - Security practices (non-root user, etc.)
	//   - Multi-stage build configuration
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
	//   - Schema validation against expected structure
	//   - Value type and format validation
	//   - Required field presence
	//   - Cross-reference consistency
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
	//   - Schema validation against expected structure
	//   - Value type and format validation
	//   - Required field presence
	//   - Nested object consistency
	//
	// Parameters:
	//   - path: Path to the JSON file
	//
	// Returns:
	//   - error: Any validation error found in the JSON file
	ValidateJSON(path string) error

	// ValidateTemplateConsistency validates consistency across frontend templates
	ValidateTemplateConsistency(templatesPath string) (*models.ValidationResult, error)

	// ValidatePackageJSONStructure validates a single package.json against standards
	ValidatePackageJSONStructure(packageJSONPath string) (*models.ValidationResult, error)

	// ValidateTypeScriptConfig validates TypeScript configuration
	ValidateTypeScriptConfig(tsconfigPath string) (*models.ValidationResult, error)

	// ValidateVercelCompatibility validates Vercel deployment compatibility
	ValidateVercelCompatibility(projectPath string) (*models.ValidationResult, error)

	// ValidateVercelConfig validates a vercel.json configuration file
	ValidateVercelConfig(vercelConfigPath string) (*models.ValidationResult, error)

	// ValidateEnvironmentVariablesConsistency validates environment variables across templates
	ValidateEnvironmentVariablesConsistency(templatesPath string) (*models.ValidationResult, error)

	// ValidateSecurityVulnerabilities validates packages for security vulnerabilities
	ValidateSecurityVulnerabilities(projectPath string) (*models.ValidationResult, error)

	// ValidatePreGeneration performs comprehensive pre-generation validation for a single template
	ValidatePreGeneration(config *models.ProjectConfig, templatePath string) (*models.ValidationResult, error)

	// ValidatePreGenerationDirectory performs pre-generation validation for an entire template directory
	ValidatePreGenerationDirectory(config *models.ProjectConfig, templateDir string) (*models.ValidationResult, error)

	// ValidateNodeJSVersionCompatibility validates Node.js version compatibility across templates
	ValidateNodeJSVersionCompatibility(projectPath string) (*models.ValidationResult, error)

	// ValidateCrossTemplateVersionConsistency validates version consistency across different template types
	ValidateCrossTemplateVersionConsistency(templatesPath string) (*models.ValidationResult, error)

	// ValidateNodeJSVersionConfiguration validates a Node.js version configuration
	ValidateNodeJSVersionConfiguration(config *models.NodeVersionConfig) (*models.ValidationResult, error)
}
