package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// TemplateValidator handles template consistency validation
type TemplateValidator struct {
	standardConfigs map[string]any
	embeddedFS      fs.FS
	useEmbedded     bool
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	// Get embedded filesystem from template package
	embeddedFS := template.GetEmbeddedFS()

	return &TemplateValidator{
		standardConfigs: make(map[string]any),
		embeddedFS:      embeddedFS,
		useEmbedded:     embeddedFS != nil,
	}
}

// NewTemplateValidatorWithFS creates a new template validator with a specific filesystem
func NewTemplateValidatorWithFS(embeddedFS fs.FS) *TemplateValidator {
	return &TemplateValidator{
		standardConfigs: make(map[string]any),
		embeddedFS:      embeddedFS,
		useEmbedded:     embeddedFS != nil,
	}
}

// ValidateTemplateConsistency validates consistency across frontend templates
func (tv *TemplateValidator) ValidateTemplateConsistency(templatesPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Template consistency validation completed",
	}

	if tv.useEmbedded {
		// Use embedded filesystem for validation
		return tv.validateEmbeddedTemplateConsistency(result)
	}

	// Fallback to filesystem-based validation
	return tv.validateFilesystemTemplateConsistency(templatesPath, result)
}

// validateEmbeddedTemplateConsistency validates templates using embedded filesystem
func (tv *TemplateValidator) validateEmbeddedTemplateConsistency(result *models.ValidationResult) (*models.ValidationResult, error) {
	// Check if frontend templates exist in embedded filesystem
	frontendPath := filepath.Join(constants.TemplateBaseDir, "frontend")

	// Check if frontend directory exists in embedded filesystem
	if _, err := fs.Stat(tv.embeddedFS, frontendPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:    "error",
				Message: "Frontend templates directory not found in embedded filesystem",
				File:    frontendPath,
			})
			return result, nil
		}
		return nil, fmt.Errorf("failed to check embedded frontend templates: %w", err)
	}

	// Validate embedded template files
	if err := tv.validateEmbeddedTemplateFiles(frontendPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate embedded template files: %w", err)
	}

	return result, nil
}

// validateFilesystemTemplateConsistency validates templates using filesystem (fallback)
func (tv *TemplateValidator) validateFilesystemTemplateConsistency(templatesPath string, result *models.ValidationResult) (*models.ValidationResult, error) {
	// Find all frontend template directories
	frontendPath := filepath.Join(templatesPath, "frontend")
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:    "error",
			Message: "Frontend templates directory not found",
			File:    frontendPath,
		})
		return result, nil
	}

	// Basic template validation
	if err := tv.validateTemplateFiles(frontendPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate template files: %w", err)
	}

	return result, nil
}

// validateEmbeddedTemplateFiles performs validation on embedded template files
func (tv *TemplateValidator) validateEmbeddedTemplateFiles(templatePath string, result *models.ValidationResult) error {
	return fs.WalkDir(tv.embeddedFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check for template files
		if strings.HasSuffix(path, ".tmpl") {
			if err := tv.validateEmbeddedTemplateFile(path, result); err != nil {
				result.Issues = append(result.Issues, models.ValidationIssue{
					Type:    "warning",
					Message: fmt.Sprintf("Embedded template validation warning: %v", err),
					File:    path,
				})
			}
		}

		return nil
	})
}

// validateTemplateFiles performs basic validation on template files (filesystem fallback)
func (tv *TemplateValidator) validateTemplateFiles(templatePath string, result *models.ValidationResult) error {
	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check for template files
		if strings.HasSuffix(path, ".tmpl") {
			if err := tv.validateTemplateFile(path, result); err != nil {
				result.Issues = append(result.Issues, models.ValidationIssue{
					Type:    "warning",
					Message: fmt.Sprintf("Template validation warning: %v", err),
					File:    path,
				})
			}
		}

		return nil
	})
}

// validateEmbeddedTemplateFile performs validation on a single embedded template file
func (tv *TemplateValidator) validateEmbeddedTemplateFile(filePath string, result *models.ValidationResult) error {
	// Read content from embedded filesystem
	content, err := fs.ReadFile(tv.embeddedFS, filePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded template file: %w", err)
	}

	return tv.validateTemplateContent(string(content), filePath, result)
}

// validateTemplateFile performs basic validation on a single template file (filesystem fallback)
func (tv *TemplateValidator) validateTemplateFile(filePath string, result *models.ValidationResult) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	return tv.validateTemplateContent(string(content), filePath, result)
}

// validateTemplateContent performs validation on template content
func (tv *TemplateValidator) validateTemplateContent(contentStr, filePath string, result *models.ValidationResult) error {
	// For .tmpl files, require template syntax
	if strings.HasSuffix(filePath, ".tmpl") {
		if !strings.Contains(contentStr, "{{") || !strings.Contains(contentStr, "}}") {
			return fmt.Errorf("template file must contain template syntax ({{ }})")
		}
	}

	// Basic template syntax validation
	if strings.Contains(contentStr, "{{") && !strings.Contains(contentStr, "}}") {
		return fmt.Errorf("unclosed template expression")
	}

	// Check for mismatched template delimiters
	openCount := strings.Count(contentStr, "{{")
	closeCount := strings.Count(contentStr, "}}")
	if openCount != closeCount {
		return fmt.Errorf("mismatched template delimiters: %d opening, %d closing", openCount, closeCount)
	}

	// Check for common template syntax issues
	if err := tv.validateTemplateSyntax(contentStr, filePath); err != nil {
		return err
	}

	// Check for security issues in templates
	if err := tv.validateTemplateSecurity(contentStr, filePath); err != nil {
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:    "warning",
			Message: fmt.Sprintf("Security warning: %v", err),
			File:    filePath,
		})
	}

	return nil
}

// validateTemplateSyntax validates template syntax for common issues
func (tv *TemplateValidator) validateTemplateSyntax(content, filePath string) error {
	// Check for common template function issues - basic validation passed

	// Check for potential infinite loops in templates
	if strings.Contains(content, "{{range") && !strings.Contains(content, "{{end}}") {
		return fmt.Errorf("unclosed range block")
	}

	// Check for potential undefined variables (basic check)
	if strings.Contains(content, "{{.") {
		// Extract variable references and check for common typos
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(line, "{{.") && strings.Contains(line, "}}") {
				// Basic validation - could be improved with more sophisticated parsing
				if strings.Contains(line, "{{. ") || strings.Contains(line, " .}}") {
					return fmt.Errorf("line %d: suspicious template syntax with spaces around dot", i+1)
				}
			}
		}
	}

	return nil
}

// validateTemplateSecurity validates template content for security issues
func (tv *TemplateValidator) validateTemplateSecurity(content, filePath string) error {
	// Check for potential code injection patterns
	dangerousPatterns := []string{
		"{{.Env",     // Environment variable access
		"{{.System",  // System access
		"{{.Exec",    // Command execution
		"{{.File",    // File system access
		"{{.Path",    // Path manipulation
		"{{.Shell",   // Shell access
		"{{.Command", // Command execution
		"{{.Process", // Process access
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(content, pattern) {
			return fmt.Errorf("potentially dangerous template pattern detected: %s", pattern)
		}
	}

	// Check for hardcoded secrets or sensitive data patterns
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"key",
		"credential",
	}

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lowerLine := strings.ToLower(line)
		for _, pattern := range sensitivePatterns {
			if strings.Contains(lowerLine, pattern) && strings.Contains(line, "=") {
				return fmt.Errorf("line %d: potential hardcoded sensitive data: %s", i+1, pattern)
			}
		}
	}

	return nil
}

// ValidatePackageJSON validates package.json consistency
func (tv *TemplateValidator) ValidatePackageJSON(packageJsonPath string) error {
	if tv.useEmbedded {
		return tv.validateEmbeddedPackageJSON(packageJsonPath)
	}
	return tv.validateFilesystemPackageJSON(packageJsonPath)
}

// validateEmbeddedPackageJSON validates package.json from embedded filesystem
func (tv *TemplateValidator) validateEmbeddedPackageJSON(packageJsonPath string) error {
	// Normalize path for embedded filesystem
	if !strings.HasPrefix(packageJsonPath, constants.TemplateBaseDir) {
		packageJsonPath = filepath.Join(constants.TemplateBaseDir, packageJsonPath)
	}

	content, err := fs.ReadFile(tv.embeddedFS, packageJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded package.json: %w", err)
	}

	return tv.validatePackageJSONContent(content, packageJsonPath)
}

// validateFilesystemPackageJSON validates package.json from filesystem (fallback)
func (tv *TemplateValidator) validateFilesystemPackageJSON(packageJsonPath string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(packageJsonPath); err != nil {
		return fmt.Errorf("invalid package.json path: %w", err)
	}

	content, err := utils.SafeReadFile(packageJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	return tv.validatePackageJSONContent(content, packageJsonPath)
}

// validatePackageJSONContent validates package.json content
func (tv *TemplateValidator) validatePackageJSONContent(content []byte, filePath string) error {
	var pkg map[string]any
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("invalid JSON syntax in %s: %w", filePath, err)
	}

	// Check required fields
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		if _, exists := pkg[field]; !exists {
			return fmt.Errorf("missing required field '%s' in %s", field, filePath)
		}
	}

	// Additional validation for package.json specific fields
	if name, ok := pkg["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("package name cannot be empty in %s", filePath)
		}
		// Check for valid npm package name format
		if strings.Contains(name, " ") {
			return fmt.Errorf("package name cannot contain spaces in %s", filePath)
		}
	}

	if version, ok := pkg["version"].(string); ok {
		if strings.TrimSpace(version) == "" {
			return fmt.Errorf("package version cannot be empty in %s", filePath)
		}
		// Basic semver validation
		if !isValidSemver(version) {
			return fmt.Errorf("invalid semver format '%s' in %s", version, filePath)
		}
	}

	return nil
}

// isValidSemver performs basic semantic version validation
func isValidSemver(version string) bool {
	// Basic semver pattern check (simplified)
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return false
		}
		// Check if part contains only digits (simplified check)
		for _, char := range part {
			if char < '0' || char > '9' {
				// Allow pre-release and build metadata for now
				if strings.Contains(part, "-") || strings.Contains(part, "+") {
					return true
				}
				return false
			}
		}
	}

	return true
}

// ValidateEmbeddedTemplateStructure validates the structure of embedded templates
func (tv *TemplateValidator) ValidateEmbeddedTemplateStructure() (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Embedded template structure validation completed",
	}

	if !tv.useEmbedded {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:    "error",
			Message: "Embedded filesystem not available",
		})
		return result, nil
	}

	// Check for required template directories
	requiredDirs := []string{
		"templates/frontend",
		"templates/backend",
		"templates/mobile",
		"templates/infrastructure",
		"templates/base",
	}

	for _, dir := range requiredDirs {
		if _, err := fs.Stat(tv.embeddedFS, dir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				result.Issues = append(result.Issues, models.ValidationIssue{
					Type:    "warning",
					Message: fmt.Sprintf("Optional template directory not found: %s", dir),
					File:    dir,
				})
			} else {
				result.Valid = false
				result.Issues = append(result.Issues, models.ValidationIssue{
					Type:    "error",
					Message: fmt.Sprintf("Error accessing template directory: %s", err.Error()),
					File:    dir,
				})
			}
		}
	}

	// Validate template metadata files
	if err := tv.validateTemplateMetadata(result); err != nil {
		return nil, fmt.Errorf("failed to validate template metadata: %w", err)
	}

	return result, nil
}

// validateTemplateMetadata validates template.yaml files in embedded templates
func (tv *TemplateValidator) validateTemplateMetadata(result *models.ValidationResult) error {
	return fs.WalkDir(tv.embeddedFS, constants.TemplateBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Look for template.yaml files
		if !d.IsDir() && d.Name() == "template.yaml" {
			if err := tv.validateTemplateYAML(path, result); err != nil {
				result.Issues = append(result.Issues, models.ValidationIssue{
					Type:    "warning",
					Message: fmt.Sprintf("Template metadata validation warning: %v", err),
					File:    path,
				})
			}
		}

		return nil
	})
}

// validateTemplateYAML validates a template.yaml file
func (tv *TemplateValidator) validateTemplateYAML(yamlPath string, result *models.ValidationResult) error {
	content, err := fs.ReadFile(tv.embeddedFS, yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read template.yaml: %w", err)
	}

	// Basic YAML structure validation
	contentStr := string(content)

	// Check for required fields in template.yaml
	requiredFields := []string{"name:", "description:", "version:"}
	for _, field := range requiredFields {
		if !strings.Contains(contentStr, field) {
			return fmt.Errorf("missing required field '%s'", strings.TrimSuffix(field, ":"))
		}
	}

	// Check for valid YAML syntax (basic check)
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Basic YAML syntax check
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if key == "" {
					return fmt.Errorf("line %d: empty key in YAML", i+1)
				}
			}
		}
	}

	return nil
}

// GetValidationSummary returns a summary of validation capabilities
func (tv *TemplateValidator) GetValidationSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"embedded_filesystem_available": tv.useEmbedded,
		"validation_methods": []string{
			"template_consistency",
			"package_json_validation",
			"template_structure_validation",
			"template_metadata_validation",
			"template_syntax_validation",
			"template_security_validation",
		},
	}

	if tv.useEmbedded {
		summary["embedded_template_base"] = constants.TemplateBaseDir
		summary["fallback_mode"] = false
	} else {
		summary["fallback_mode"] = true
		summary["filesystem_validation_only"] = true
	}

	return summary
}

// ValidateAllEmbeddedTemplates performs comprehensive validation of all embedded templates
func (tv *TemplateValidator) ValidateAllEmbeddedTemplates() (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Comprehensive embedded template validation completed",
	}

	if !tv.useEmbedded {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:    "error",
			Message: "Embedded filesystem not available for comprehensive validation",
		})
		return result, nil
	}

	// Validate template structure
	structureResult, err := tv.ValidateEmbeddedTemplateStructure()
	if err != nil {
		return nil, fmt.Errorf("failed to validate template structure: %w", err)
	}

	// Merge structure validation results
	result.Issues = append(result.Issues, structureResult.Issues...)
	if !structureResult.Valid {
		result.Valid = false
	}

	// Validate template consistency
	consistencyResult, err := tv.ValidateTemplateConsistency("")
	if err != nil {
		return nil, fmt.Errorf("failed to validate template consistency: %w", err)
	}

	// Merge consistency validation results
	result.Issues = append(result.Issues, consistencyResult.Issues...)
	if !consistencyResult.Valid {
		result.Valid = false
	}

	// Update summary
	totalIssues := len(result.Issues)
	errorCount := 0
	warningCount := 0

	for _, issue := range result.Issues {
		switch issue.Type {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		}
	}

	result.Summary = fmt.Sprintf("Validation completed: %d total issues (%d errors, %d warnings)",
		totalIssues, errorCount, warningCount)

	return result, nil
}
