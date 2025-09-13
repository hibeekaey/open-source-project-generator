// Package validation provides comprehensive validation capabilities for generated
// projects, templates, and configurations in the Open Source Template Generator.
//
// This package implements the ValidationEngine interface and provides:
//   - Project structure and file organization validation
//   - Configuration file syntax and semantic validation
//   - Cross-template consistency and compatibility checks
//   - Security vulnerability scanning and best practices validation
//   - Version compatibility validation across different technologies
//   - Platform-specific deployment validation (Vercel, Docker, Kubernetes)
//
// The validation engine ensures that generated projects meet quality standards,
// follow best practices, and are free from common configuration issues.
//
// Key Features:
//   - Detailed validation results with actionable feedback
//   - Performance optimization for large projects
//   - Extensible validation rules and custom validators
//   - Integration with external validation tools and services
//   - Security-focused validation with vulnerability scanning
//   - Cross-platform compatibility validation
//
// Validation Types:
//   - Syntax validation for JSON, YAML, and other configuration files
//   - Semantic validation for package.json, go.mod, Dockerfile
//   - Cross-file consistency validation
//   - Version compatibility validation
//   - Security best practices validation
//   - Platform deployment readiness validation
//
// Usage:
//
//	validator := validation.NewEngine()
//	result, err := validator.ValidateProject("/path/to/project")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if !result.Valid {
//	    for _, issue := range result.Issues {
//	        fmt.Printf("Issue: %s\n", issue.Message)
//	    }
//	}
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/open-source-template-generator/pkg/constants"
	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

// Engine implements the ValidationEngine interface
type Engine struct {
	configValidator   *models.ConfigValidator
	templateValidator *TemplateValidator
	vercelValidator   *VercelValidator
	securityValidator *SecurityValidator
	versionValidator  *VersionValidator
}

// NewEngine creates a new validation engine
func NewEngine() interfaces.ValidationEngine {
	return &Engine{
		configValidator:   models.NewConfigValidator(),
		templateValidator: NewTemplateValidator(),
		vercelValidator:   NewVercelValidator(),
		securityValidator: NewSecurityValidator(),
		versionValidator:  NewVersionValidator(),
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

// ValidatePackageJSON validates a package.json file with comprehensive validation
func (e *Engine) ValidatePackageJSON(path string) error {
	// Validate file path first
	validator := utils.NewValidator()
	validator.ValidateFilePath(path, "package_json_path")

	if validator.HasErrors() {
		return fmt.Errorf("invalid package.json path: %s", utils.FormatValidationErrors(validator.GetErrors()))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return utils.NewFileSystemError(path, "read", "failed to read package.json file", err)
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return utils.NewValidationError("package_json_content", "invalid JSON format in package.json", err)
	}

	// Use enhanced validation for package.json structure
	return e.validatePackageJSONStructure(packageJSON, validator)
}

// validatePackageJSONStructure validates the structure and content of package.json
func (e *Engine) validatePackageJSONStructure(packageJSON map[string]interface{}, validator *utils.Validator) error {
	// Validate required fields
	requiredFields := []string{"name", "version", "scripts"}
	for _, field := range requiredFields {
		if _, exists := packageJSON[field]; !exists {
			validator.AddError(field, fmt.Sprintf("Required field '%s' is missing from package.json", field), "required", nil)
		}
	}

	// Validate package name format
	if name, exists := packageJSON["name"]; exists {
		if nameStr, ok := name.(string); ok {
			validator.ValidateStringLength(nameStr, "name", 1, 214)
			validator.ValidateStringPattern(nameStr, "name", `^[a-z0-9]([a-z0-9\-_.])*$`, "valid npm package name")
		} else {
			validator.AddError("name", "Package name must be a string", "invalid_type", name)
		}
	}

	// Validate version format
	if version, exists := packageJSON["version"]; exists {
		if versionStr, ok := version.(string); ok {
			validator.ValidateStringPattern(versionStr, "version", `^\d+\.\d+\.\d+`, "semantic version")
		} else {
			validator.AddError("version", "Version must be a string", "invalid_type", version)
		}
	}

	// Validate description if present
	if description, exists := packageJSON["description"]; exists {
		if descStr, ok := description.(string); ok {
			validator.ValidateStringLength(descStr, "description", 0, 500)
			validator.ValidateSecureString(descStr, "description")
		}
	}

	// Validate author email if present
	if author, exists := packageJSON["author"]; exists {
		if authorMap, ok := author.(map[string]interface{}); ok {
			if email, emailExists := authorMap["email"]; emailExists {
				if emailStr, ok := email.(string); ok {
					validator.ValidateEmail(emailStr, "author.email")
				}
			}
		}
	}

	// Validate repository URL if present
	if repository, exists := packageJSON["repository"]; exists {
		if repoMap, ok := repository.(map[string]interface{}); ok {
			if url, urlExists := repoMap["url"]; urlExists {
				if urlStr, ok := url.(string); ok {
					validator.ValidateURL(urlStr, "repository.url")
				}
			}
		}
	}

	if validator.HasErrors() {
		return fmt.Errorf("package.json validation failed: %s", utils.FormatValidationErrors(validator.GetErrors()))
	}

	return nil
}

// ValidateGoMod validates a go.mod file with comprehensive validation
func (e *Engine) ValidateGoMod(path string) error {
	// Validate file path first
	validator := utils.NewValidator()
	validator.ValidateFilePath(path, "go_mod_path")

	if validator.HasErrors() {
		return fmt.Errorf("invalid go.mod path: %s", utils.FormatValidationErrors(validator.GetErrors()))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return utils.NewFileSystemError(path, "read", "failed to read go.mod file", err)
	}

	content := string(data)

	// Validate content security
	validator.ValidateSecureString(content, "go_mod_content")

	if validator.HasErrors() {
		return fmt.Errorf("go.mod security validation failed: %s", utils.FormatValidationErrors(validator.GetErrors()))
	}

	return e.validateGoModStructure(content, validator)
}

// validateGoModStructure validates the structure and content of go.mod
func (e *Engine) validateGoModStructure(content string, validator *utils.Validator) error {
	lines := strings.Split(content, "\n")

	// Check for required declarations
	hasModule := false
	hasGoVersion := false
	var moduleName string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "module ") {
			hasModule = true
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				moduleName = parts[1]
				// Validate module name format
				validator.ValidateStringPattern(moduleName, "module_name",
					`^[a-zA-Z0-9._/-]+$`, "valid Go module name")

				// Check for common security issues in module names
				validator.ValidateSecureString(moduleName, "module_name")
			} else {
				validator.AddError("module_declaration", "Module declaration is incomplete", "incomplete", line)
			}
		}

		if strings.HasPrefix(line, "go ") {
			hasGoVersion = true
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				version := parts[1]
				if !e.isValidGoVersion(version) {
					validator.AddError("go_version", fmt.Sprintf("Invalid Go version format: %s", version), "invalid_format", version)
				}

				// Validate minimum Go version for security
				if e.isGoVersionTooOld(version) {
					validator.AddError("go_version", fmt.Sprintf("Go version %s is too old and may have security vulnerabilities", version), "security_risk", version)
				}
			} else {
				validator.AddError("go_declaration", "Go version declaration is incomplete", "incomplete", line)
			}
		}

		// Validate require statements for security
		if strings.HasPrefix(line, "require ") || strings.Contains(line, "require(") {
			e.validateGoRequireStatement(line, validator)
		}
	}

	// Check for required declarations
	if !hasModule {
		validator.AddError("module_declaration", "Missing module declaration in go.mod", "required", nil)
	}

	if !hasGoVersion {
		validator.AddError("go_version", "Missing Go version declaration in go.mod", "required", nil)
	}

	if validator.HasErrors() {
		return fmt.Errorf("go.mod validation failed: %s", utils.FormatValidationErrors(validator.GetErrors()))
	}

	return nil
}

// validateGoRequireStatement validates individual require statements
func (e *Engine) validateGoRequireStatement(line string, validator *utils.Validator) {
	// Basic validation for require statements
	if strings.Contains(line, "//") {
		// Remove comments for validation
		line = strings.Split(line, "//")[0]
	}

	line = strings.TrimSpace(line)

	// Check for suspicious patterns in dependencies
	suspiciousPatterns := []string{
		"localhost",
		"127.0.0.1",
		"file://",
		"../",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(line), pattern) {
			validator.AddError("require_statement",
				fmt.Sprintf("Potentially unsafe dependency pattern detected: %s", pattern),
				"security_risk", line)
		}
	}
}

// isGoVersionTooOld checks if a Go version is too old for security
func (e *Engine) isGoVersionTooOld(version string) bool {
	// Consider versions older than 1.19 as potentially risky
	// This is a simplified check - in production you might want more sophisticated version comparison
	if strings.HasPrefix(version, "1.1") && len(version) >= 4 {
		if version[3] < '9' {
			return true
		}
	}
	return false
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
		case constants.FilePackageJSON:
			if err := e.ValidatePackageJSON(path); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   relativePath,
					Tag:     "syntax",
					Value:   fileName,
					Message: fmt.Sprintf("Package.json validation failed: %s", err.Error()),
				})
			}
		case constants.FileGoMod:
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
		case constants.FileDockerfile:
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

		if fileName == constants.FileDockerfile || strings.HasSuffix(fileName, ".Dockerfile") {
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

		if !info.IsDir() && filepath.Base(path) == constants.FilePackageJSON {
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

// ValidatePreGeneration performs comprehensive pre-generation validation
func (e *Engine) ValidatePreGeneration(config *models.ProjectConfig, templatePath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Validate that configuration is not nil
	if config == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "config",
			Tag:     "required",
			Value:   "nil",
			Message: "Project configuration cannot be nil",
		})
		return result, nil
	}

	// Validate version configuration exists
	if config.Versions == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions",
			Tag:     "required",
			Value:   "nil",
			Message: "Version configuration is required for template generation",
		})
		return result, nil
	}

	// Validate Node.js versions for frontend templates
	if e.isFrontendTemplate(templatePath) {
		if err := e.validateNodeJSPreGeneration(config, result); err != nil {
			return nil, fmt.Errorf("node.js pre-generation validation failed: %w", err)
		}
	}

	// Validate Go versions for backend templates
	if e.isBackendTemplate(templatePath) {
		if err := e.validateGoPreGeneration(config, result); err != nil {
			return nil, fmt.Errorf("go pre-generation validation failed: %w", err)
		}
	}

	// Validate template-specific version requirements
	if err := e.validateTemplateSpecificVersions(config, templatePath, result); err != nil {
		return nil, fmt.Errorf("template-specific validation failed: %w", err)
	}

	return result, nil
}

// ValidatePreGenerationDirectory performs pre-generation validation for an entire template directory
func (e *Engine) ValidatePreGenerationDirectory(config *models.ProjectConfig, templateDir string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Validate that configuration is not nil
	if config == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "config",
			Tag:     "required",
			Value:   "nil",
			Message: "Project configuration cannot be nil",
		})
		return result, nil
	}

	// Validate version configuration exists
	if config.Versions == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions",
			Tag:     "required",
			Value:   "nil",
			Message: "Version configuration is required for template generation",
		})
		return result, nil
	}

	// Collect all template files for comprehensive validation
	templateFiles, err := e.collectTemplateFiles(templateDir)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "templateDir",
			Tag:     "accessible",
			Value:   templateDir,
			Message: fmt.Sprintf("Failed to collect template files: %s", err.Error()),
		})
		return result, nil
	}

	// Validate Node.js versions if frontend templates are present
	if e.hasFrontendTemplates(templateFiles) {
		if err := e.validateNodeJSPreGeneration(config, result); err != nil {
			return nil, fmt.Errorf("node.js pre-generation validation failed: %w", err)
		}
	}

	// Validate Go versions if backend templates are present
	if e.hasBackendTemplates(templateFiles) {
		if err := e.validateGoPreGeneration(config, result); err != nil {
			return nil, fmt.Errorf("go pre-generation validation failed: %w", err)
		}
	}

	// Validate cross-template consistency
	if err := e.validateCrossTemplateConsistency(config, templateFiles, result); err != nil {
		return nil, fmt.Errorf("cross-template consistency validation failed: %w", err)
	}

	return result, nil
}

// validateNodeJSPreGeneration validates Node.js configuration for pre-generation
func (e *Engine) validateNodeJSPreGeneration(config *models.ProjectConfig, result *models.ValidationResult) error {
	if config.Versions.NodeJS == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.nodejs",
			Tag:     "required",
			Value:   "nil",
			Message: "Node.js version configuration is required for frontend templates",
		})
		return nil
	}

	// Use version validator to validate Node.js configuration
	versionValidator := NewVersionValidator()
	versionResult := versionValidator.ValidateNodeVersionConfig(config.Versions.NodeJS)

	// Convert version validation errors to validation result errors
	for _, vErr := range versionResult.Errors {
		severity := "error"
		if vErr.Severity == "critical" {
			severity = "critical"
			result.Valid = false
		}

		result.Errors = append(result.Errors, models.ValidationError{
			Field:   fmt.Sprintf("versions.nodejs.%s", vErr.Field),
			Tag:     "version_validation",
			Value:   vErr.Value,
			Message: fmt.Sprintf("[%s] %s", severity, vErr.Message),
		})
	}

	// Convert version validation warnings to validation result warnings
	for _, vWarn := range versionResult.Warnings {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   fmt.Sprintf("versions.nodejs.%s", vWarn.Field),
			Message: vWarn.Message,
		})
	}

	return nil
}

// validateGoPreGeneration validates Go configuration for pre-generation
func (e *Engine) validateGoPreGeneration(config *models.ProjectConfig, result *models.ValidationResult) error {
	if config.Versions.Go == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.go",
			Tag:     "required",
			Value:   "",
			Message: "Go version is required for backend templates",
		})
		return nil
	}

	// Validate Go version format
	if !e.isValidGoVersion(config.Versions.Go) {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.go",
			Tag:     "format",
			Value:   config.Versions.Go,
			Message: fmt.Sprintf("Invalid Go version format: %s", config.Versions.Go),
		})
		return nil
	}

	// Check minimum Go version requirement
	majorMinor, err := e.extractGoMajorMinor(config.Versions.Go)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.go",
			Tag:     "parse",
			Value:   config.Versions.Go,
			Message: fmt.Sprintf("Failed to parse Go version: %s", err.Error()),
		})
		return nil
	}

	if majorMinor < 1.20 {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.go",
			Tag:     "minimum_version",
			Value:   config.Versions.Go,
			Message: fmt.Sprintf("Go version %s is not supported, minimum required version is 1.20", config.Versions.Go),
		})
	}

	return nil
}

// validateTemplateSpecificVersions validates version requirements specific to template types
func (e *Engine) validateTemplateSpecificVersions(config *models.ProjectConfig, templatePath string, result *models.ValidationResult) error {
	// Validate package.json template versions
	if strings.Contains(templatePath, "package.json.tmpl") {
		return e.validatePackageJSONTemplateVersions(config, templatePath, result)
	}

	// Validate go.mod template versions
	if strings.Contains(templatePath, "go.mod.tmpl") {
		return e.validateGoModTemplateVersions(config, templatePath, result)
	}

	// Validate Dockerfile template versions
	if strings.Contains(templatePath, "Dockerfile.tmpl") {
		return e.validateDockerfileTemplateVersions(config, templatePath, result)
	}

	return nil
}

// validatePackageJSONTemplateVersions validates versions for package.json templates
func (e *Engine) validatePackageJSONTemplateVersions(config *models.ProjectConfig, templatePath string, result *models.ValidationResult) error {
	if config.Versions.NodeJS == nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.nodejs",
			Tag:     "required",
			Value:   "nil",
			Message: fmt.Sprintf("Node.js version configuration required for package.json template: %s", templatePath),
		})
		return nil
	}

	// Validate Node.js runtime and @types/node compatibility
	runtime := config.Versions.NodeJS.Runtime
	types := config.Versions.NodeJS.TypesPackage

	if runtime == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.nodejs.runtime",
			Tag:     "required",
			Value:   "",
			Message: fmt.Sprintf("Node.js runtime version required for template: %s", templatePath),
		})
	}

	if types == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.nodejs.types",
			Tag:     "required",
			Value:   "",
			Message: fmt.Sprintf("@types/node version required for template: %s", templatePath),
		})
	}

	// Validate compatibility if both are present
	if runtime != "" && types != "" {
		runtimeMajor, err1 := e.extractMajorVersion(runtime)
		typesMajor, err2 := e.extractMajorVersion(types)

		if err1 != nil || err2 != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "versions.nodejs.compatibility",
				Tag:     "parse",
				Value:   fmt.Sprintf("runtime: %s, types: %s", runtime, types),
				Message: fmt.Sprintf("Failed to parse version numbers for compatibility check in template: %s", templatePath),
			})
		} else if typesMajor < runtimeMajor || typesMajor > runtimeMajor+2 {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "versions.nodejs.compatibility",
				Tag:     "mismatch",
				Value:   fmt.Sprintf("runtime: %d, types: %d", runtimeMajor, typesMajor),
				Message: fmt.Sprintf("@types/node version %d incompatible with Node.js runtime version %d in template: %s", typesMajor, runtimeMajor, templatePath),
			})
		}
	}

	return nil
}

// validateGoModTemplateVersions validates versions for go.mod templates
func (e *Engine) validateGoModTemplateVersions(config *models.ProjectConfig, templatePath string, result *models.ValidationResult) error {
	if config.Versions.Go == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "versions.go",
			Tag:     "required",
			Value:   "",
			Message: fmt.Sprintf("Go version required for go.mod template: %s", templatePath),
		})
	}

	return nil
}

// validateDockerfileTemplateVersions validates versions for Dockerfile templates
func (e *Engine) validateDockerfileTemplateVersions(config *models.ProjectConfig, templatePath string, result *models.ValidationResult) error {
	// Validate Docker image configuration for Node.js templates
	if e.isFrontendTemplate(templatePath) && config.Versions.NodeJS != nil {
		dockerImage := config.Versions.NodeJS.DockerImage
		if dockerImage == "" {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "versions.nodejs.docker",
				Message: fmt.Sprintf("Docker image not specified for Node.js template: %s", templatePath),
			})
		} else if !strings.Contains(dockerImage, "node:") {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "versions.nodejs.docker",
				Tag:     "invalid_image",
				Value:   dockerImage,
				Message: fmt.Sprintf("Docker image %s is not a Node.js image for template: %s", dockerImage, templatePath),
			})
		}
	}

	return nil
}

// validateCrossTemplateConsistency validates version consistency across multiple templates
func (e *Engine) validateCrossTemplateConsistency(config *models.ProjectConfig, templateFiles []string, result *models.ValidationResult) error {
	// Group templates by type
	packageJSONTemplates := []string{}
	dockerTemplates := []string{}

	for _, file := range templateFiles {
		if strings.Contains(file, "package.json.tmpl") {
			packageJSONTemplates = append(packageJSONTemplates, file)
		} else if strings.Contains(file, "Dockerfile.tmpl") && e.isFrontendTemplate(file) {
			dockerTemplates = append(dockerTemplates, file)
		}
	}

	// Validate consistency across package.json templates
	if len(packageJSONTemplates) > 1 {
		if config.Versions.NodeJS == nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "versions.nodejs",
				Tag:     "consistency",
				Value:   fmt.Sprintf("%d templates", len(packageJSONTemplates)),
				Message: "Node.js version configuration required for consistent package.json generation across multiple templates",
			})
		} else {
			// All package.json templates will use the same versions - this is good
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field: "versions.nodejs.consistency",
				Message: fmt.Sprintf("All %d package.json templates will use Node.js runtime: %s, @types/node: %s",
					len(packageJSONTemplates), config.Versions.NodeJS.Runtime, config.Versions.NodeJS.TypesPackage),
			})
		}
	}

	// Validate consistency across Docker templates
	if len(dockerTemplates) > 0 {
		if config.Versions.NodeJS == nil || config.Versions.NodeJS.DockerImage == "" {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "versions.nodejs.docker",
				Message: fmt.Sprintf("Docker image not specified for %d Docker templates", len(dockerTemplates)),
			})
		}
	}

	return nil
}

// Helper methods for template analysis

// isFrontendTemplate checks if a template is for frontend components
func (e *Engine) isFrontendTemplate(templatePath string) bool {
	frontendIndicators := []string{
		"frontend/",
		"nextjs-",
		"package.json.tmpl",
		"next.config.js.tmpl",
		"tailwind.config.js.tmpl",
		"tsconfig.json.tmpl",
	}

	for _, indicator := range frontendIndicators {
		if strings.Contains(templatePath, indicator) {
			return true
		}
	}

	return false
}

// isBackendTemplate checks if a template is for backend components
func (e *Engine) isBackendTemplate(templatePath string) bool {
	backendIndicators := []string{
		"backend/",
		"go-gin/",
		"go.mod.tmpl",
		"main.go.tmpl",
		"internal/",
		"pkg/",
	}

	for _, indicator := range backendIndicators {
		if strings.Contains(templatePath, indicator) {
			return true
		}
	}

	return false
}

// collectTemplateFiles collects all template files in a directory
func (e *Engine) collectTemplateFiles(templateDir string) ([]string, error) {
	var templateFiles []string

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			templateFiles = append(templateFiles, path)
		}

		return nil
	})

	return templateFiles, err
}

// hasFrontendTemplates checks if any frontend templates are present
func (e *Engine) hasFrontendTemplates(templateFiles []string) bool {
	for _, file := range templateFiles {
		if e.isFrontendTemplate(file) {
			return true
		}
	}
	return false
}

// hasBackendTemplates checks if any backend templates are present
func (e *Engine) hasBackendTemplates(templateFiles []string) bool {
	for _, file := range templateFiles {
		if e.isBackendTemplate(file) {
			return true
		}
	}
	return false
}

// extractMajorVersion extracts major version number from a version string
func (e *Engine) extractMajorVersion(version string) (int, error) {
	// Remove version operators (>=, ^, ~, etc.)
	re := regexp.MustCompile(`^[>=<~^]*`)
	cleanVersion := re.ReplaceAllString(version, "")

	// Split by dots and get the first part
	parts := strings.Split(cleanVersion, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	return major, nil
}

// extractGoMajorMinor extracts major.minor version from Go version string
func (e *Engine) extractGoMajorMinor(version string) (float64, error) {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid Go version format: %s", version)
	}

	major := parts[0]
	minor := parts[1]

	// Convert to float for comparison
	majorFloat := 0.0
	minorFloat := 0.0

	for _, char := range major {
		if char >= '0' && char <= '9' {
			majorFloat = majorFloat*10 + float64(char-'0')
		}
	}

	for _, char := range minor {
		if char >= '0' && char <= '9' {
			minorFloat = minorFloat*10 + float64(char-'0')
		}
	}

	return majorFloat + minorFloat/100, nil
}

// ValidateNodeJSVersionCompatibility validates Node.js version compatibility across templates
func (e *Engine) ValidateNodeJSVersionCompatibility(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Collect all package.json files
	packageJSONFiles, err := e.collectPackageJSONFiles(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to collect package.json files: %w", err)
	}

	if len(packageJSONFiles) == 0 {
		return result, nil // No package.json files to validate
	}

	// Parse Node.js versions from each package.json
	nodeVersions := make(map[string]NodeJSVersionInfo)
	for _, filePath := range packageJSONFiles {
		versionInfo, err := e.extractNodeJSVersionFromPackageJSON(filePath)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   filePath,
				Tag:     "parse_error",
				Value:   "",
				Message: fmt.Sprintf("Failed to parse Node.js version from %s: %s", filePath, err.Error()),
			})
			continue
		}
		nodeVersions[filePath] = *versionInfo
	}

	// Validate consistency across all package.json files
	if err := e.validateNodeJSVersionConsistency(nodeVersions, result); err != nil {
		return nil, fmt.Errorf("failed to validate Node.js version consistency: %w", err)
	}

	// Validate individual version configurations
	for filePath, versionInfo := range nodeVersions {
		if err := e.validateIndividualNodeJSVersion(filePath, versionInfo, result); err != nil {
			return nil, fmt.Errorf("failed to validate individual Node.js version for %s: %w", filePath, err)
		}
	}

	return result, nil
}

// ValidateCrossTemplateVersionConsistency validates version consistency across different template types
func (e *Engine) ValidateCrossTemplateVersionConsistency(templatesPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Collect all template files
	templateFiles, err := e.collectTemplateFiles(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to collect template files: %w", err)
	}

	// Group templates by type
	templateGroups := e.groupTemplatesByType(templateFiles)

	// Validate Node.js version consistency across frontend templates
	if len(templateGroups["frontend"]) > 0 {
		if err := e.validateFrontendTemplateVersionConsistency(templateGroups["frontend"], result); err != nil {
			return nil, fmt.Errorf("failed to validate frontend template consistency: %w", err)
		}
	}

	// Validate Docker image consistency
	if len(templateGroups["docker"]) > 0 {
		if err := e.validateDockerTemplateVersionConsistency(templateGroups["docker"], result); err != nil {
			return nil, fmt.Errorf("failed to validate Docker template consistency: %w", err)
		}
	}

	// Validate CI/CD template consistency
	if len(templateGroups["ci"]) > 0 {
		if err := e.validateCITemplateVersionConsistency(templateGroups["ci"], result); err != nil {
			return nil, fmt.Errorf("failed to validate CI template consistency: %w", err)
		}
	}

	return result, nil
}

// ValidateNodeJSVersionConfiguration validates a Node.js version configuration using the version validator
func (e *Engine) ValidateNodeJSVersionConfiguration(config *models.NodeVersionConfig) (*models.ValidationResult, error) {
	// Use the version validator to validate the configuration
	versionResult := e.versionValidator.ValidateNodeVersionConfig(config)

	// Convert version validation result to standard validation result
	result := &models.ValidationResult{
		Valid:    versionResult.Valid,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Convert version validation errors
	for _, vErr := range versionResult.Errors {
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   vErr.Field,
			Tag:     "version_validation",
			Value:   vErr.Value,
			Message: vErr.Message,
		})
	}

	// Convert version validation warnings
	for _, vWarn := range versionResult.Warnings {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   vWarn.Field,
			Message: vWarn.Message,
		})
	}

	return result, nil
}

// Helper types and methods for Node.js version validation

// NodeJSVersionInfo holds extracted Node.js version information from package.json
type NodeJSVersionInfo struct {
	EnginesNode  string
	TypesNode    string
	NPMVersion   string
	HasEngines   bool
	HasTypesNode bool
}

// collectPackageJSONFiles collects all package.json files in a project
func (e *Engine) collectPackageJSONFiles(projectPath string) ([]string, error) {
	var packageJSONFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Base(path) == constants.FilePackageJSON {
			packageJSONFiles = append(packageJSONFiles, path)
		}

		return nil
	})

	return packageJSONFiles, err
}

// extractNodeJSVersionFromPackageJSON extracts Node.js version information from a package.json file
func (e *Engine) extractNodeJSVersionFromPackageJSON(filePath string) (*NodeJSVersionInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	versionInfo := &NodeJSVersionInfo{}

	// Extract engines.node
	if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
		if node, ok := engines["node"].(string); ok {
			versionInfo.EnginesNode = node
			versionInfo.HasEngines = true
		}
		if npm, ok := engines["npm"].(string); ok {
			versionInfo.NPMVersion = npm
		}
	}

	// Extract @types/node from devDependencies
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		if typesNode, ok := devDeps["@types/node"].(string); ok {
			versionInfo.TypesNode = typesNode
			versionInfo.HasTypesNode = true
		}
	}

	// Also check dependencies (though @types/node is usually in devDependencies)
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		if typesNode, ok := deps["@types/node"].(string); ok && versionInfo.TypesNode == "" {
			versionInfo.TypesNode = typesNode
			versionInfo.HasTypesNode = true
		}
	}

	return versionInfo, nil
}

// validateNodeJSVersionConsistency validates consistency across multiple Node.js version configurations
func (e *Engine) validateNodeJSVersionConsistency(nodeVersions map[string]NodeJSVersionInfo, result *models.ValidationResult) error {
	if len(nodeVersions) <= 1 {
		return nil // Nothing to compare
	}

	// Collect all unique engine versions and types versions
	engineVersions := make(map[string][]string)
	typesVersions := make(map[string][]string)

	for filePath, versionInfo := range nodeVersions {
		if versionInfo.HasEngines && versionInfo.EnginesNode != "" {
			engineVersions[versionInfo.EnginesNode] = append(engineVersions[versionInfo.EnginesNode], filePath)
		}
		if versionInfo.HasTypesNode && versionInfo.TypesNode != "" {
			typesVersions[versionInfo.TypesNode] = append(typesVersions[versionInfo.TypesNode], filePath)
		}
	}

	// Check for inconsistencies in engines.node
	if len(engineVersions) > 1 {
		result.Valid = false
		var versions []string
		for version := range engineVersions {
			versions = append(versions, version)
		}
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "engines.node",
			Tag:     "consistency",
			Value:   strings.Join(versions, ", "),
			Message: fmt.Sprintf("Inconsistent Node.js engine versions found: %s", strings.Join(versions, ", ")),
		})
	}

	// Check for inconsistencies in @types/node
	if len(typesVersions) > 1 {
		result.Valid = false
		var versions []string
		for version := range typesVersions {
			versions = append(versions, version)
		}
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "@types/node",
			Tag:     "consistency",
			Value:   strings.Join(versions, ", "),
			Message: fmt.Sprintf("Inconsistent @types/node versions found: %s", strings.Join(versions, ", ")),
		})
	}

	// Warn about missing engines or types in some files
	filesWithEngines := 0
	filesWithTypes := 0
	for _, versionInfo := range nodeVersions {
		if versionInfo.HasEngines {
			filesWithEngines++
		}
		if versionInfo.HasTypesNode {
			filesWithTypes++
		}
	}

	totalFiles := len(nodeVersions)
	if filesWithEngines > 0 && filesWithEngines < totalFiles {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "engines.node",
			Message: fmt.Sprintf("Only %d out of %d package.json files specify engines.node", filesWithEngines, totalFiles),
		})
	}

	if filesWithTypes > 0 && filesWithTypes < totalFiles {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "@types/node",
			Message: fmt.Sprintf("Only %d out of %d package.json files specify @types/node", filesWithTypes, totalFiles),
		})
	}

	return nil
}

// validateIndividualNodeJSVersion validates an individual Node.js version configuration
func (e *Engine) validateIndividualNodeJSVersion(filePath string, versionInfo NodeJSVersionInfo, result *models.ValidationResult) error {
	relativePath := filepath.Base(filepath.Dir(filePath))

	// Validate engines.node format if present
	if versionInfo.HasEngines && versionInfo.EnginesNode != "" {
		if err := e.versionValidator.validateVersionFormat(versionInfo.EnginesNode, "runtime"); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   fmt.Sprintf("%s/engines.node", relativePath),
				Tag:     "format",
				Value:   versionInfo.EnginesNode,
				Message: fmt.Sprintf("Invalid engines.node format in %s: %s", filePath, err.Error()),
			})
		}
	}

	// Validate @types/node format if present
	if versionInfo.HasTypesNode && versionInfo.TypesNode != "" {
		if err := e.versionValidator.validateVersionFormat(versionInfo.TypesNode, "types"); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   fmt.Sprintf("%s/@types/node", relativePath),
				Tag:     "format",
				Value:   versionInfo.TypesNode,
				Message: fmt.Sprintf("Invalid @types/node format in %s: %s", filePath, err.Error()),
			})
		}
	}

	// Validate compatibility between engines.node and @types/node if both are present
	if versionInfo.HasEngines && versionInfo.HasTypesNode && versionInfo.EnginesNode != "" && versionInfo.TypesNode != "" {
		if err := e.versionValidator.validateRuntimeTypesCompatibility(versionInfo.EnginesNode, versionInfo.TypesNode); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   fmt.Sprintf("%s/compatibility", relativePath),
				Tag:     "compatibility",
				Value:   fmt.Sprintf("engines: %s, types: %s", versionInfo.EnginesNode, versionInfo.TypesNode),
				Message: fmt.Sprintf("Version compatibility issue in %s: %s", filePath, err.Error()),
			})
		}
	}

	return nil
}

// groupTemplatesByType groups template files by their type (frontend, docker, ci, etc.)
func (e *Engine) groupTemplatesByType(templateFiles []string) map[string][]string {
	groups := make(map[string][]string)

	for _, file := range templateFiles {
		if e.isFrontendTemplate(file) {
			groups["frontend"] = append(groups["frontend"], file)
		}
		if e.isDockerTemplate(file) {
			groups["docker"] = append(groups["docker"], file)
		}
		if e.isCITemplate(file) {
			groups["ci"] = append(groups["ci"], file)
		}
	}

	return groups
}

// isDockerTemplate checks if a template is Docker-related
func (e *Engine) isDockerTemplate(templatePath string) bool {
	dockerIndicators := []string{
		"Dockerfile.tmpl",
		"docker-compose.yml.tmpl",
		"docker-compose.yaml.tmpl",
		".dockerignore.tmpl",
	}

	fileName := filepath.Base(templatePath)
	for _, indicator := range dockerIndicators {
		if fileName == indicator {
			return true
		}
	}

	return false
}

// isCITemplate checks if a template is CI/CD-related
func (e *Engine) isCITemplate(templatePath string) bool {
	ciIndicators := []string{
		".github/workflows/",
		".gitlab-ci.yml.tmpl",
		"Jenkinsfile.tmpl",
		".travis.yml.tmpl",
	}

	for _, indicator := range ciIndicators {
		if strings.Contains(templatePath, indicator) {
			return true
		}
	}

	return false
}

// validateFrontendTemplateVersionConsistency validates version consistency across frontend templates
func (e *Engine) validateFrontendTemplateVersionConsistency(frontendTemplates []string, result *models.ValidationResult) error {
	// Extract Node.js version references from frontend templates
	nodeVersionRefs := make(map[string][]string)

	for _, templatePath := range frontendTemplates {
		versions, err := e.extractVersionReferencesFromTemplate(templatePath)
		if err != nil {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   templatePath,
				Message: fmt.Sprintf("Could not extract version references from template: %s", err.Error()),
			})
			continue
		}

		for version, refs := range versions {
			nodeVersionRefs[version] = append(nodeVersionRefs[version], refs...)
		}
	}

	// Check for inconsistencies
	if len(nodeVersionRefs) > 1 {
		var versions []string
		for version := range nodeVersionRefs {
			versions = append(versions, version)
		}
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "frontend_templates",
			Message: fmt.Sprintf("Multiple Node.js version references found in frontend templates: %s", strings.Join(versions, ", ")),
		})
	}

	return nil
}

// validateDockerTemplateVersionConsistency validates version consistency across Docker templates
func (e *Engine) validateDockerTemplateVersionConsistency(dockerTemplates []string, result *models.ValidationResult) error {
	// Extract Node.js base image references from Docker templates
	baseImages := make(map[string][]string)

	for _, templatePath := range dockerTemplates {
		images, err := e.extractDockerBaseImagesFromTemplate(templatePath)
		if err != nil {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   templatePath,
				Message: fmt.Sprintf("Could not extract Docker base images from template: %s", err.Error()),
			})
			continue
		}

		for image := range images {
			baseImages[image] = append(baseImages[image], templatePath)
		}
	}

	// Check for inconsistencies in Node.js base images
	nodeImages := make(map[string][]string)
	for image, templates := range baseImages {
		if strings.Contains(strings.ToLower(image), "node") {
			nodeImages[image] = templates
		}
	}

	if len(nodeImages) > 1 {
		var images []string
		for image := range nodeImages {
			images = append(images, image)
		}
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "docker_templates",
			Message: fmt.Sprintf("Multiple Node.js Docker base images found: %s", strings.Join(images, ", ")),
		})
	}

	return nil
}

// validateCITemplateVersionConsistency validates version consistency across CI/CD templates
func (e *Engine) validateCITemplateVersionConsistency(ciTemplates []string, result *models.ValidationResult) error {
	// Extract Node.js version references from CI templates
	nodeVersionRefs := make(map[string][]string)

	for _, templatePath := range ciTemplates {
		versions, err := e.extractNodeVersionsFromCITemplate(templatePath)
		if err != nil {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   templatePath,
				Message: fmt.Sprintf("Could not extract Node.js versions from CI template: %s", err.Error()),
			})
			continue
		}

		for version := range versions {
			nodeVersionRefs[version] = append(nodeVersionRefs[version], templatePath)
		}
	}

	// Check for inconsistencies
	if len(nodeVersionRefs) > 1 {
		var versions []string
		for version := range nodeVersionRefs {
			versions = append(versions, version)
		}
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "ci_templates",
			Message: fmt.Sprintf("Multiple Node.js versions found in CI templates: %s", strings.Join(versions, ", ")),
		})
	}

	return nil
}

// extractVersionReferencesFromTemplate extracts version references from a template file
func (e *Engine) extractVersionReferencesFromTemplate(templatePath string) (map[string][]string, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	content := string(data)
	versions := make(map[string][]string)

	// Look for Node.js version patterns in templates
	nodeVersionPatterns := []string{
		`"node":\s*"([^"]+)"`,             // engines.node in package.json
		`@types/node":\s*"([^"]+)"`,       // @types/node dependency
		`{{\.Versions\.NodeJS\.Runtime}}`, // Template variable reference
		`{{\.Versions\.NodeJS\.TypesPackage}}`,
	}

	for _, pattern := range nodeVersionPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				version := match[1]
				versions[version] = append(versions[version], templatePath)
			} else {
				// For template variables, use the pattern itself as the key
				versions[pattern] = append(versions[pattern], templatePath)
			}
		}
	}

	return versions, nil
}

// extractDockerBaseImagesFromTemplate extracts Docker base images from a template file
func (e *Engine) extractDockerBaseImagesFromTemplate(templatePath string) (map[string]bool, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	content := string(data)
	images := make(map[string]bool)

	// Look for FROM instructions in Dockerfiles
	fromPattern := regexp.MustCompile(`FROM\s+([^\s\n]+)`)
	matches := fromPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			images[match[1]] = true
		}
	}

	return images, nil
}

// extractNodeVersionsFromCITemplate extracts Node.js versions from CI template files
func (e *Engine) extractNodeVersionsFromCITemplate(templatePath string) (map[string]bool, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	content := string(data)
	versions := make(map[string]bool)

	// Look for Node.js version patterns in CI files
	nodeVersionPatterns := []string{
		`node-version:\s*['"]?([^'"\s\n]+)['"]?`,                  // GitHub Actions
		`node_js:\s*['"]?([^'"\s\n]+)['"]?`,                       // Travis CI
		`image:\s*node:([^\s\n]+)`,                                // Docker-based CI
		`setup-node@.*\n.*node-version:\s*['"]?([^'"\s\n]+)['"]?`, // GitHub Actions setup-node
	}

	for _, pattern := range nodeVersionPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				versions[match[1]] = true
			}
		}
	}

	return versions, nil
}
