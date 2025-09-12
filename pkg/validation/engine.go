package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

// Engine implements the ValidationEngine interface
type Engine struct {
	configValidator   *models.ConfigValidator
	templateValidator *TemplateValidator
	vercelValidator   *VercelValidator
	securityValidator *SecurityValidator
}

// NewEngine creates a new validation engine
func NewEngine() interfaces.ValidationEngine {
	return &Engine{
		configValidator:   models.NewConfigValidator(),
		templateValidator: NewTemplateValidator(),
		vercelValidator:   NewVercelValidator(),
		securityValidator: NewSecurityValidator(),
	}
}

// ValidateProject validates the entire generated project structure
func (e *Engine) ValidateProject(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check if project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "ProjectPath",
			Tag:     "exists",
			Value:   projectPath,
			Message: "Project directory does not exist",
		})
		return result, nil
	}

	// Validate project structure
	if err := e.validateProjectStructure(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate project structure: %w", err)
	}

	// Validate configuration files (includes package files)
	if err := e.validateConfigurationFiles(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate configuration files: %w", err)
	}

	// Validate Docker files
	if err := e.validateDockerFiles(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate Docker files: %w", err)
	}

	// Validate dependency compatibility
	if err := e.validateDependencyCompatibility(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate dependency compatibility: %w", err)
	}

	return result, nil
}

// ValidatePackageJSON validates a package.json file
func (e *Engine) ValidatePackageJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate required fields
	requiredFields := []string{"name", "version", "scripts"}
	for _, field := range requiredFields {
		if _, exists := packageJSON[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate name format
	if name, ok := packageJSON["name"].(string); ok {
		if strings.Contains(name, " ") || strings.Contains(name, "_") {
			return fmt.Errorf("package name should use kebab-case format")
		}
	}

	// Validate version format
	if version, ok := packageJSON["version"].(string); ok {
		if !e.isValidSemVer(version) {
			return fmt.Errorf("invalid semantic version format: %s", version)
		}
	}

	return nil
}

// ValidateGoMod validates a go.mod file
func (e *Engine) ValidateGoMod(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Check for module declaration
	hasModule := false
	hasGoVersion := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			hasModule = true
		}
		if strings.HasPrefix(line, "go ") {
			hasGoVersion = true
			// Validate Go version format
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				version := parts[1]
				if !e.isValidGoVersion(version) {
					return fmt.Errorf("invalid Go version format: %s", version)
				}
			}
		}
	}

	if !hasModule {
		return fmt.Errorf("missing module declaration")
	}

	if !hasGoVersion {
		return fmt.Errorf("missing Go version declaration")
	}

	return nil
}

// ValidateDockerfile validates a Dockerfile
func (e *Engine) ValidateDockerfile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	hasFrom := false
	hasWorkdir := false
	hasCopy := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "FROM ") {
			hasFrom = true
		}
		if strings.HasPrefix(line, "WORKDIR ") {
			hasWorkdir = true
		}
		if strings.HasPrefix(line, "COPY ") || strings.HasPrefix(line, "ADD ") {
			hasCopy = true
		}
	}

	if !hasFrom {
		return fmt.Errorf("missing FROM instruction")
	}

	if !hasWorkdir {
		return fmt.Errorf("missing WORKDIR instruction")
	}

	if !hasCopy {
		return fmt.Errorf("missing COPY or ADD instruction")
	}

	return nil
}

// ValidateYAML validates a YAML configuration file
func (e *Engine) ValidateYAML(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var yamlData interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}

	return nil
}

// ValidateJSON validates a JSON configuration file
func (e *Engine) ValidateJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

// validateProjectStructure validates the overall project directory structure
func (e *Engine) validateProjectStructure(projectPath string, result *models.ValidationResult) error {
	// Define expected directories based on common project structures
	expectedDirs := map[string]bool{
		"frontend": false,
		"backend":  false,
		"mobile":   false,
		"deploy":   false,
		"docs":     false,
		".github":  false,
	}

	// Check which directories exist
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return fmt.Errorf("failed to read project directory: %w", err)
	}

	foundDirs := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			foundDirs[entry.Name()] = true
			if _, expected := expectedDirs[entry.Name()]; expected {
				expectedDirs[entry.Name()] = true
			}
		}
	}

	// Check for required files
	requiredFiles := []string{"README.md", "Makefile"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "ProjectStructure",
				Message: fmt.Sprintf("Missing recommended file: %s", file),
			})
		}
	}

	// Validate that at least one main component directory exists
	hasMainComponent := foundDirs["frontend"] || foundDirs["backend"] || foundDirs["mobile"]
	if !hasMainComponent {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "ProjectStructure",
			Tag:     "required",
			Value:   "",
			Message: "Project must contain at least one main component directory (frontend, backend, or mobile)",
		})
	}

	return nil
}

// validateConfigurationFiles validates various configuration files
func (e *Engine) validateConfigurationFiles(projectPath string, result *models.ValidationResult) error {
	// Track processed files to avoid duplicates
	processedFiles := make(map[string]bool)

	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relativePath, _ := filepath.Rel(projectPath, path)
		fileName := filepath.Base(path)

		// Skip if already processed
		if processedFiles[relativePath] {
			return nil
		}
		processedFiles[relativePath] = true

		// Validate specific file types
		switch fileName {
		case "package.json":
			if err := e.ValidatePackageJSON(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Package.json validation failed: %s", err.Error()),
				})
			}
		case "go.mod":
			if err := e.ValidateGoMod(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Go.mod validation failed: %s", err.Error()),
				})
			}
		case "docker-compose.yml", "docker-compose.yaml":
			if err := e.ValidateYAML(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Docker Compose validation failed: %s", err.Error()),
				})
			}
		case "Dockerfile":
			if err := e.ValidateDockerfile(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Dockerfile validation failed: %s", err.Error()),
				})
			}
		}

		// Check workflow files
		if strings.Contains(relativePath, ".github/workflows/") && (strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml")) {
			if err := e.ValidateYAML(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Workflow file validation failed: %s", err.Error()),
				})
			}
		}

		return nil
	})

	return err
}

// validatePackageFiles validates package management files
func (e *Engine) validatePackageFiles(projectPath string, result *models.ValidationResult) error {
	// Find and validate package.json files
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		relativePath, _ := filepath.Rel(projectPath, path)

		switch fileName {
		case "package.json":
			if err := e.ValidatePackageJSON(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Package.json validation failed: %s", err.Error()),
				})
			}
		case "go.mod":
			if err := e.ValidateGoMod(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Go.mod validation failed: %s", err.Error()),
				})
			}
		}

		return nil
	})

	return err
}

// validateDockerFiles validates Docker-related files
func (e *Engine) validateDockerFiles(projectPath string, result *models.ValidationResult) error {
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		relativePath, _ := filepath.Rel(projectPath, path)

		if fileName == "Dockerfile" || strings.HasSuffix(fileName, ".Dockerfile") {
			if err := e.ValidateDockerfile(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Dockerfile validation failed: %s", err.Error()),
				})
			}
		}

		return nil
	})

	return err
}

// validateDependencyCompatibility validates dependency compatibility across components
func (e *Engine) validateDependencyCompatibility(projectPath string, result *models.ValidationResult) error {
	// This is a placeholder for more sophisticated dependency compatibility checking
	// In a real implementation, this would:
	// 1. Parse all package.json files and check for conflicting dependencies
	// 2. Validate Go module dependencies for compatibility
	// 3. Check Docker base image compatibility
	// 4. Validate Kubernetes resource compatibility

	// For now, we'll add a basic check for common issues
	packageJSONFiles := []string{}

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == "package.json" {
			packageJSONFiles = append(packageJSONFiles, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Check for potential dependency conflicts between package.json files
	if len(packageJSONFiles) > 1 {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Dependencies",
			Message: fmt.Sprintf("Found %d package.json files - ensure dependency versions are compatible", len(packageJSONFiles)),
		})
	}

	return nil
}

// Helper functions

// isValidSemVer checks if a version string is a valid semantic version
func (e *Engine) isValidSemVer(version string) bool {
	// Basic semantic version validation
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check if part contains only digits (basic check)
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

// isValidGoVersion checks if a Go version string is valid
func (e *Engine) isValidGoVersion(version string) bool {
	// Go versions can be like "1.21", "1.21.0", etc.
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

// ValidateTemplateConsistency validates consistency across frontend templates
func (e *Engine) ValidateTemplateConsistency(templatesPath string) (*models.ValidationResult, error) {
	return e.templateValidator.ValidateTemplateConsistency(templatesPath)
}

// ValidatePackageJSONStructure validates a single package.json against standards
func (e *Engine) ValidatePackageJSONStructure(packageJSONPath string) (*models.ValidationResult, error) {
	return e.templateValidator.ValidatePackageJSONStructure(packageJSONPath)
}

// ValidateTypeScriptConfig validates TypeScript configuration
func (e *Engine) ValidateTypeScriptConfig(tsconfigPath string) (*models.ValidationResult, error) {
	return e.templateValidator.ValidateTypeScriptConfig(tsconfigPath)
}

// ValidateVercelCompatibility validates Vercel deployment compatibility
func (e *Engine) ValidateVercelCompatibility(projectPath string) (*models.ValidationResult, error) {
	return e.vercelValidator.ValidateVercelCompatibility(projectPath)
}

// ValidateVercelConfig validates a vercel.json configuration file
func (e *Engine) ValidateVercelConfig(vercelConfigPath string) (*models.ValidationResult, error) {
	return e.vercelValidator.ValidateVercelConfig(vercelConfigPath)
}

// ValidateEnvironmentVariablesConsistency validates environment variables across templates
func (e *Engine) ValidateEnvironmentVariablesConsistency(templatesPath string) (*models.ValidationResult, error) {
	return e.vercelValidator.ValidateEnvironmentVariablesConsistency(templatesPath)
}

// ValidateSecurityVulnerabilities validates packages for security vulnerabilities
func (e *Engine) ValidateSecurityVulnerabilities(projectPath string) (*models.ValidationResult, error) {
	return e.securityValidator.ValidateSecurityVulnerabilities(projectPath)
}
