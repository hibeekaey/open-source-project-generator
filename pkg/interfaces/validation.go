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
}
