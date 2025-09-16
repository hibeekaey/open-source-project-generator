package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// TemplateValidator handles template consistency validation
type TemplateValidator struct {
	standardConfigs map[string]any
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{
		standardConfigs: make(map[string]any),
	}
}

// ValidateTemplateConsistency validates consistency across frontend templates
func (tv *TemplateValidator) ValidateTemplateConsistency(templatesPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:   true,
		Issues:  []models.ValidationIssue{},
		Summary: "Template consistency validation completed",
	}

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

// validateTemplateFiles performs basic validation on template files
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

// validateTemplateFile performs basic validation on a single template file
func (tv *TemplateValidator) validateTemplateFile(filePath string, result *models.ValidationResult) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	contentStr := string(content)

	// Basic template syntax validation
	if strings.Contains(contentStr, "{{") && !strings.Contains(contentStr, "}}") {
		return fmt.Errorf("unclosed template expression")
	}

	// Check for basic template structure
	if strings.Contains(contentStr, "{{") && strings.Contains(contentStr, "}}") {
		// Basic validation passed
		return nil
	}

	_ = result // Suppress unused parameter warning
	return nil
}

// ValidatePackageJSON validates package.json consistency
func (tv *TemplateValidator) ValidatePackageJSON(packageJsonPath string) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(packageJsonPath); err != nil {
		return fmt.Errorf("invalid package.json path: %w", err)
	}

	content, err := utils.SafeReadFile(packageJsonPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]any
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("invalid JSON syntax: %w", err)
	}

	// Check required fields
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		if _, exists := pkg[field]; !exists {
			return fmt.Errorf("missing required field '%s'", field)
		}
	}

	return nil
}
