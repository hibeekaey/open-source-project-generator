package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MetadataValidator handles validation of template metadata files and content.
type MetadataValidator struct {
	parser *MetadataParser
}

// NewMetadataValidator creates a new metadata validator instance.
func NewMetadataValidator(parser *MetadataParser) *MetadataValidator {
	return &MetadataValidator{
		parser: parser,
	}
}

// ValidateTemplateMetadataFile validates template metadata file if present.
func (v *MetadataValidator) ValidateTemplateMetadataFile(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check for metadata files
	metadataFiles := []string{"template.yaml", "template.yml", "metadata.yaml", "metadata.yml"}
	var foundMetadata bool

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if _, err := os.Stat(metadataPath); err == nil {
			foundMetadata = true
			// Validate metadata file content
			if validationIssues := v.ValidateMetadataFileContent(metadataPath); len(validationIssues) > 0 {
				issues = append(issues, validationIssues...)
			}
			break
		}
	}

	if !foundMetadata {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "No metadata file found (template.yaml, template.yml, metadata.yaml, or metadata.yml)",
			Rule:     "has-metadata",
			Fixable:  true,
		})
	}

	return issues
}

// ValidateMetadataFileContent validates the content of a metadata file.
func (v *MetadataValidator) ValidateMetadataFileContent(metadataPath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Validate file path to prevent path traversal attacks
	if err := v.validateFilePath(metadataPath); err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid metadata file path: %v", err),
			File:     metadataPath,
			Rule:     "path-validation",
			Fixable:  false,
		})
		return issues
	}

	// Read metadata file
	content, err := os.ReadFile(metadataPath) // #nosec G304 - path is validated above
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot read metadata file: %v", err),
			File:     metadataPath,
			Rule:     "metadata-readable",
			Fixable:  false,
		})
		return issues
	}

	// Validate YAML syntax by attempting to parse
	templateName := filepath.Base(filepath.Dir(metadataPath))
	if _, err := v.parser.ParseTemplateYAML(content, templateName); err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid YAML syntax: %v", err),
			File:     metadataPath,
			Rule:     "yaml-syntax",
			Fixable:  false,
		})
		return issues
	}

	// Basic YAML syntax check (simplified)
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for basic YAML key-value format
		if !strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  "Line does not appear to be valid YAML",
				File:     metadataPath,
				Line:     i + 1,
				Rule:     "yaml-syntax",
				Fixable:  false,
			})
		}
	}

	// Check for required metadata fields
	requiredFields := []string{"name:", "description:", "version:", "author:"}
	for _, field := range requiredFields {
		if !strings.Contains(contentStr, field) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Missing recommended field: %s", strings.TrimSuffix(field, ":")),
				File:     metadataPath,
				Rule:     "required-fields",
				Fixable:  true,
			})
		}
	}

	return issues
}

// ValidateMetadata validates template metadata structure and content.
func (v *MetadataValidator) ValidateMetadata(metadata *models.TemplateMetadata) []models.ValidationIssue {
	var issues []models.ValidationIssue

	if metadata == nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Metadata cannot be nil",
			Rule:     "metadata-exists",
			Fixable:  false,
		})
		return issues
	}

	// Validate required fields
	if metadata.Name == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template name is required",
			Rule:     "name-required",
			Fixable:  false,
		})
	}

	if metadata.Description == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template description is recommended",
			Rule:     "description-recommended",
			Fixable:  true,
		})
	}

	if metadata.Version == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template version is recommended",
			Rule:     "version-recommended",
			Fixable:  true,
		})
	}

	if metadata.Author == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template author is recommended",
			Rule:     "author-recommended",
			Fixable:  true,
		})
	}

	if metadata.License == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template license is recommended",
			Rule:     "license-recommended",
			Fixable:  true,
		})
	}

	// Validate name format (should be kebab-case)
	if metadata.Name != "" && !v.isValidTemplateName(metadata.Name) {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template name should be in kebab-case format",
			Rule:     "name-format",
			Fixable:  true,
		})
	}

	// Validate category
	if metadata.Category != "" {
		validCategories := []string{"backend", "frontend", "mobile", "infrastructure", "base"}
		if !v.contains(validCategories, metadata.Category) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Invalid category: %s. Valid categories: %v", metadata.Category, validCategories),
				Rule:     "category-valid",
				Fixable:  true,
			})
		}
	}

	// Validate version format (basic semver check)
	if metadata.Version != "" && !v.isValidVersion(metadata.Version) {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Version should follow semantic versioning (e.g., 1.0.0)",
			Rule:     "version-format",
			Fixable:  true,
		})
	}

	// Validate variables
	for varName, templateVar := range metadata.Variables {
		if varIssues := v.validateTemplateVariable(varName, templateVar); len(varIssues) > 0 {
			issues = append(issues, varIssues...)
		}
	}

	return issues
}

// validateTemplateVariable validates a single template variable.
func (v *MetadataValidator) validateTemplateVariable(varName string, templateVar models.TemplateVar) []models.ValidationIssue {
	var issues []models.ValidationIssue

	if templateVar.Name == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  fmt.Sprintf("Variable '%s' should have a name field", varName),
			Rule:     "variable-name",
			Fixable:  true,
		})
	}

	if templateVar.Type == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  fmt.Sprintf("Variable '%s' should have a type field", varName),
			Rule:     "variable-type",
			Fixable:  true,
		})
	} else {
		// Validate type is one of the supported types
		validTypes := []string{"string", "number", "boolean", "array", "object"}
		if !v.contains(validTypes, templateVar.Type) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Variable '%s' has unsupported type '%s'. Valid types: %v", varName, templateVar.Type, validTypes),
				Rule:     "variable-type-valid",
				Fixable:  true,
			})
		}
	}

	if templateVar.Description == "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "info",
			Severity: "info",
			Message:  fmt.Sprintf("Variable '%s' should have a description", varName),
			Rule:     "variable-description",
			Fixable:  true,
		})
	}

	return issues
}

// validateFilePath validates file path to prevent path traversal attacks.
func (v *MetadataValidator) validateFilePath(path string) error {
	// Clean the path to resolve any .. or . components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts (only if it contains .. after cleaning)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected")
	}

	// Allow absolute paths for temp directories and valid file operations
	// Only reject paths that are clearly malicious
	if strings.Contains(path, "../") && strings.Contains(path, "etc") {
		return fmt.Errorf("suspicious path detected")
	}

	return nil
}

// isValidTemplateName checks if template name follows kebab-case convention.
func (v *MetadataValidator) isValidTemplateName(name string) bool {
	// Empty names are invalid
	if name == "" {
		return false
	}

	// Basic kebab-case validation: lowercase letters, numbers, and hyphens only
	for _, char := range name {
		if char < 'a' || char > 'z' {
			if char < '0' || char > '9' {
				if char != '-' {
					return false
				}
			}
		}
	}

	// Should not start or end with hyphen
	return !strings.HasPrefix(name, "-") && !strings.HasSuffix(name, "-")
}

// isValidVersion checks if version follows basic semantic versioning.
func (v *MetadataValidator) isValidVersion(version string) bool {
	// Basic semver pattern: X.Y.Z where X, Y, Z are numbers
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
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

// contains checks if a slice contains a specific string.
func (v *MetadataValidator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
