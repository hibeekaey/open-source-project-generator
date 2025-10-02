// Package template provides template validation functionality for the
// Open Source Project Generator.
package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateValidator handles template validation operations
type TemplateValidator struct {
	// Dependencies for enhanced validation
	templateManager interfaces.TemplateManager
}

// NewTemplateValidator creates a new template validator instance
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{}
}

// NewTemplateValidatorWithManager creates a new template validator instance with template manager
func NewTemplateValidatorWithManager(templateManager interfaces.TemplateManager) *TemplateValidator {
	return &TemplateValidator{
		templateManager: templateManager,
	}
}

// ValidateTemplate validates a template structure and metadata
func (tv *TemplateValidator) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &interfaces.TemplateValidationResult{
			Valid: false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: "error",
					Message:  fmt.Sprintf("Template path does not exist: %s", path),
					Rule:     "path-exists",
					Fixable:  false,
				},
			},
		}, nil
	}

	var issues []models.ValidationIssue
	var warnings []models.ValidationIssue

	// Validate template structure
	structureIssues := tv.validateTemplateStructure(path)
	issues = append(issues, structureIssues...)

	// Validate metadata if present
	metadataIssues := tv.validateTemplateMetadataFile(path)
	issues = append(issues, metadataIssues...)

	// Validate template files
	templateIssues := tv.validateTemplateFiles(path)
	issues = append(issues, templateIssues...)

	// Separate errors from warnings
	var errors []models.ValidationIssue
	for _, issue := range issues {
		if issue.Severity == "error" {
			errors = append(errors, issue)
		} else {
			warnings = append(warnings, issue)
		}
	}

	// Convert models.ValidationIssue to interfaces.ValidationIssue
	interfaceErrors := make([]interfaces.ValidationIssue, len(errors))
	for i, issue := range errors {
		interfaceErrors[i] = interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
			Metadata: issue.Metadata,
		}
	}

	interfaceWarnings := make([]interfaces.ValidationIssue, len(warnings))
	for i, issue := range warnings {
		interfaceWarnings[i] = interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
			Metadata: issue.Metadata,
		}
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(errors) == 0,
		Issues:   interfaceErrors,
		Warnings: interfaceWarnings,
	}, nil
}

// ValidateTemplateStructure validates template structure
func (tv *TemplateValidator) ValidateTemplateStructure(template *interfaces.TemplateInfo) error {
	// Validate required fields
	if template.Name == "" {
		return fmt.Errorf("🚫 template name is required")
	}
	if template.Category == "" {
		return fmt.Errorf("🚫 template category is required")
	}
	if template.Version == "" {
		return fmt.Errorf("🚫 template version is required")
	}

	// Validate name format (should be kebab-case)
	if !tv.isValidTemplateName(template.Name) {
		return fmt.Errorf("🚫 template name must be in kebab-case format")
	}

	// Validate category
	validCategories := []string{"backend", "frontend", "mobile", "infrastructure", "base"}
	if !tv.contains(validCategories, template.Category) {
		return fmt.Errorf("🚫 invalid category: %s. Valid categories: %v", template.Category, validCategories)
	}

	return nil
}

// ValidateTemplateMetadata validates template metadata
func (tv *TemplateValidator) ValidateTemplateMetadata(metadata *interfaces.TemplateMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	if metadata.Author == "" {
		return fmt.Errorf("metadata author is required")
	}
	if metadata.License == "" {
		return fmt.Errorf("metadata license is required")
	}

	return nil
}

// ValidateCustomTemplate validates custom template
func (tv *TemplateValidator) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Use the same validation logic as ValidateTemplate
	return tv.ValidateTemplate(path)
}

// validateTemplateStructure validates the basic structure of a template directory
func (tv *TemplateValidator) validateTemplateStructure(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check if it's a directory
	info, err := os.Stat(templatePath)
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot access template path: %v", err),
			Rule:     "path-accessible",
			Fixable:  false,
		})
		return issues
	}

	if !info.IsDir() {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template path must be a directory",
			Rule:     "is-directory",
			Fixable:  false,
		})
		return issues
	}

	// Check for required files/directories
	requiredItems := []string{
		// At least one template file should exist
	}

	for _, item := range requiredItems {
		itemPath := filepath.Join(templatePath, item)
		if _, err := os.Stat(itemPath); os.IsNotExist(err) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Recommended item missing: %s", item),
				Rule:     "recommended-structure",
				Fixable:  true,
			})
		}
	}

	// Check for template files
	hasTemplateFiles, err := tv.hasTemplateFiles(templatePath)
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Error checking template files: %v", err),
			Rule:     "template-files-check",
			Fixable:  false,
		})
	} else if !hasTemplateFiles {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "No template files (.tmpl) found in template directory",
			Rule:     "has-template-files",
			Fixable:  false,
		})
	}

	return issues
}

// validateTemplateMetadataFile validates template metadata file if present
func (tv *TemplateValidator) validateTemplateMetadataFile(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check for metadata files
	metadataFiles := []string{"template.yaml", "template.yml", "metadata.yaml", "metadata.yml"}
	var foundMetadata bool

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if _, err := os.Stat(metadataPath); err == nil {
			foundMetadata = true
			// Validate metadata file content
			if validationIssues := tv.validateMetadataFileContent(metadataPath); len(validationIssues) > 0 {
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

// validateTemplateFiles validates individual template files
func (tv *TemplateValidator) validateTemplateFiles(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check template files
		if strings.HasSuffix(path, ".tmpl") {
			if fileIssues := tv.validateTemplateFile(path); len(fileIssues) > 0 {
				issues = append(issues, fileIssues...)
			}
		}

		return nil
	})

	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Error walking template directory: %v", err),
			Rule:     "directory-walk",
			Fixable:  false,
		})
	}

	return issues
}

// validateTemplateFile validates a single template file
func (tv *TemplateValidator) validateTemplateFile(filePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Validate file path to prevent path traversal attacks
	if err := tv.validateFilePath(filePath); err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid file path: %v", err),
			File:     filePath,
			Rule:     "path-validation",
			Fixable:  false,
		})
		return issues
	}

	// Read file content
	content, err := os.ReadFile(filePath) // #nosec G304 - path is validated above
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot read template file: %v", err),
			File:     filePath,
			Rule:     "file-readable",
			Fixable:  false,
		})
		return issues
	}

	// Basic template syntax validation
	contentStr := string(content)

	// Check for unmatched template delimiters
	openCount := strings.Count(contentStr, "{{")
	closeCount := strings.Count(contentStr, "}}")

	if openCount != closeCount {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Unmatched template delimiters {{ and }}",
			File:     filePath,
			Rule:     "template-syntax",
			Fixable:  false,
		})
	}

	// Check for common template variables
	commonVars := []string{"{{.Name}}", "{{.Organization}}", "{{.Description}}"}
	hasVars := false
	for _, variable := range commonVars {
		if strings.Contains(contentStr, variable) {
			hasVars = true
			break
		}
	}

	if !hasVars && openCount > 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "info",
			Severity: "info",
			Message:  "Template file contains template syntax but no common variables",
			File:     filePath,
			Rule:     "has-common-vars",
			Fixable:  false,
		})
	}

	return issues
}

// validateMetadataFileContent validates the content of a metadata file
func (tv *TemplateValidator) validateMetadataFileContent(metadataPath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Validate file path to prevent path traversal attacks
	if err := tv.validateFilePath(metadataPath); err != nil {
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

// hasTemplateFiles checks if directory contains template files
func (tv *TemplateValidator) hasTemplateFiles(templatePath string) (bool, error) {
	hasFiles := false

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			hasFiles = true
			return filepath.SkipDir // Stop walking once we find a template file
		}

		return nil
	})

	return hasFiles, err
}

// isValidTemplateName checks if template name follows kebab-case convention
func (tv *TemplateValidator) isValidTemplateName(name string) bool {
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

// contains checks if slice contains string
func (tv *TemplateValidator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// validateFilePath validates file path to prevent path traversal attacks
func (tv *TemplateValidator) validateFilePath(filePath string) error {
	// Clean the path to resolve any .. or . elements
	cleanPath := filepath.Clean(filePath)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in file path")
	}

	// Ensure path is absolute or relative to current directory
	if filepath.IsAbs(cleanPath) {
		// For absolute paths, ensure they don't access system directories
		systemDirs := []string{"/etc", "/proc", "/sys", "/dev", "/root"}
		for _, sysDir := range systemDirs {
			if strings.HasPrefix(cleanPath, sysDir) {
				return fmt.Errorf("access to system directory not allowed")
			}
		}
	}

	return nil
}

// ValidateTemplateStructureAdvanced performs advanced template structure validation
func (tv *TemplateValidator) ValidateTemplateStructureAdvanced(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	var issues []interfaces.ValidationIssue
	var warnings []interfaces.ValidationIssue

	// Validate required fields
	if templateInfo.Name == "" {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template name is required",
			Rule:     "required-name",
			Fixable:  false,
		})
	}

	if templateInfo.Category == "" {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template category is required",
			Rule:     "required-category",
			Fixable:  false,
		})
	}

	if templateInfo.Version == "" {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template version is not specified",
			Rule:     "missing-version",
			Fixable:  true,
		})
	}

	// Validate name format (should be kebab-case)
	if !tv.isValidTemplateName(templateInfo.Name) {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template name must be in kebab-case format (lowercase letters, numbers, and hyphens only)",
			Rule:     "name-format",
			Fixable:  false,
		})
	}

	// Validate category
	validCategories := []string{"backend", "frontend", "mobile", "infrastructure", "base"}
	if !tv.contains(validCategories, templateInfo.Category) {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid category: %s. Valid categories: %v", templateInfo.Category, validCategories),
			Rule:     "valid-category",
			Fixable:  false,
		})
	}

	// Validate description
	if templateInfo.Description == "" {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template description is missing",
			Rule:     "missing-description",
			Fixable:  true,
		})
	} else if len(templateInfo.Description) < 10 {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template description is too short (should be at least 10 characters)",
			Rule:     "short-description",
			Fixable:  true,
		})
	}

	// Validate technology field
	if templateInfo.Technology == "" || templateInfo.Technology == "Unknown" {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template technology is not specified or unknown",
			Rule:     "missing-technology",
			Fixable:  true,
		})
	}

	// Validate tags
	if len(templateInfo.Tags) == 0 {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template has no tags for better discoverability",
			Rule:     "missing-tags",
			Fixable:  true,
		})
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(issues) == 0,
		Issues:   issues,
		Warnings: warnings,
		Summary: interfaces.ValidationSummary{
			ErrorCount:   len(issues),
			WarningCount: len(warnings),
			FixableCount: 0, // Would be calculated based on fixable issues
		},
	}, nil
}

// ValidateTemplateDependencies validates template dependencies
func (tv *TemplateValidator) ValidateTemplateDependencies(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	var issues []interfaces.ValidationIssue
	var warnings []interfaces.ValidationIssue

	// Check if template manager is available for dependency validation
	if tv.templateManager == nil {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template manager not available for dependency validation",
			Rule:     "no-template-manager",
			Fixable:  false,
		})
		return &interfaces.TemplateValidationResult{
			Valid:    true,
			Issues:   issues,
			Warnings: warnings,
			Summary: interfaces.ValidationSummary{
				ErrorCount:   0,
				WarningCount: len(warnings),
				FixableCount: 0,
			},
		}, nil
	}

	// Validate each dependency
	for _, dependency := range templateInfo.Dependencies {
		if dependency == "" {
			issues = append(issues, interfaces.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  "Empty dependency name found",
				Rule:     "empty-dependency",
				Fixable:  false,
			})
			continue
		}

		// Check if dependency exists (for template dependencies)
		if tv.isTemplateDependency(dependency) {
			_, err := tv.templateManager.GetTemplateInfo(dependency)
			if err != nil {
				issues = append(issues, interfaces.ValidationIssue{
					Type:     "error",
					Severity: "error",
					Message:  fmt.Sprintf("Template dependency '%s' not found", dependency),
					Rule:     "missing-template-dependency",
					Fixable:  false,
				})
			}
		}

		// Validate dependency format
		if !tv.isValidDependencyFormat(dependency) {
			warnings = append(warnings, interfaces.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Dependency '%s' may have invalid format", dependency),
				Rule:     "dependency-format",
				Fixable:  true,
			})
		}
	}

	// Check for circular dependencies
	circularDeps := tv.detectCircularDependencies(templateInfo)
	for _, dep := range circularDeps {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Circular dependency detected: %s", dep),
			Rule:     "circular-dependency",
			Fixable:  false,
		})
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(issues) == 0,
		Issues:   issues,
		Warnings: warnings,
		Summary: interfaces.ValidationSummary{
			FixableCount: 0,
			ErrorCount:   len(issues),
			WarningCount: len(warnings),
		},
	}, nil
}

// ValidateTemplateCompatibility validates template compatibility
func (tv *TemplateValidator) ValidateTemplateCompatibility(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	var issues []interfaces.ValidationIssue
	var warnings []interfaces.ValidationIssue

	// Get compatibility information if template manager is available
	if tv.templateManager != nil {
		compatInfo, err := tv.templateManager.GetTemplateCompatibility(templateInfo.Name)
		if err == nil {
			// Validate minimum generator version
			if compatInfo.MinGeneratorVersion != "" {
				// In a real implementation, you would compare with current generator version
				// For now, we'll just check if it's specified
				if !tv.isValidVersionFormat(compatInfo.MinGeneratorVersion) {
					warnings = append(warnings, interfaces.ValidationIssue{
						Type:     "warning",
						Severity: "warning",
						Message:  fmt.Sprintf("Invalid minimum generator version format: %s", compatInfo.MinGeneratorVersion),
						Rule:     "invalid-min-version",
						Fixable:  true,
					})
				}
			}

			// Validate maximum generator version
			if compatInfo.MaxGeneratorVersion != "" {
				if !tv.isValidVersionFormat(compatInfo.MaxGeneratorVersion) {
					warnings = append(warnings, interfaces.ValidationIssue{
						Type:     "warning",
						Severity: "warning",
						Message:  fmt.Sprintf("Invalid maximum generator version format: %s", compatInfo.MaxGeneratorVersion),
						Rule:     "invalid-max-version",
						Fixable:  true,
					})
				}
			}

			// Validate supported platforms
			if len(compatInfo.SupportedPlatforms) == 0 {
				warnings = append(warnings, interfaces.ValidationIssue{
					Type:     "warning",
					Severity: "warning",
					Message:  "No supported platforms specified",
					Rule:     "missing-platforms",
					Fixable:  true,
				})
			} else {
				validPlatforms := []string{"linux", "darwin", "windows", "freebsd", "openbsd"}
				for _, platform := range compatInfo.SupportedPlatforms {
					if !tv.contains(validPlatforms, platform) {
						warnings = append(warnings, interfaces.ValidationIssue{
							Type:     "warning",
							Severity: "warning",
							Message:  fmt.Sprintf("Unknown platform: %s", platform),
							Rule:     "unknown-platform",
							Fixable:  true,
						})
					}
				}
			}

			// Validate required features
			for _, feature := range compatInfo.RequiredFeatures {
				if feature == "" {
					warnings = append(warnings, interfaces.ValidationIssue{
						Type:     "warning",
						Severity: "warning",
						Message:  "Empty required feature found",
						Rule:     "empty-feature",
						Fixable:  true,
					})
				}
			}

			// Check for conflicts
			for _, conflict := range compatInfo.Conflicts {
				if conflict == templateInfo.Name {
					issues = append(issues, interfaces.ValidationIssue{
						Type:     "error",
						Severity: "error",
						Message:  "Template conflicts with itself",
						Rule:     "self-conflict",
						Fixable:  false,
					})
				}
			}
		}
	}

	// Validate metadata compatibility fields
	if templateInfo.Metadata.Created.IsZero() {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template creation date is not specified",
			Rule:     "missing-created-date",
			Fixable:  true,
		})
	}

	if templateInfo.Metadata.Updated.IsZero() {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template update date is not specified",
			Rule:     "missing-updated-date",
			Fixable:  true,
		})
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(issues) == 0,
		Issues:   issues,
		Warnings: warnings,
		Summary: interfaces.ValidationSummary{
			ErrorCount:   len(issues),
			WarningCount: len(warnings),
			FixableCount: 0,
		},
	}, nil
}

// ValidateTemplateMetadataAdvanced performs advanced metadata validation
func (tv *TemplateValidator) ValidateTemplateMetadataAdvanced(metadata *interfaces.TemplateMetadata) (*interfaces.TemplateValidationResult, error) {
	var issues []interfaces.ValidationIssue
	var warnings []interfaces.ValidationIssue

	// Validate author
	if metadata.Author == "" {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template author is not specified",
			Rule:     "missing-author",
			Fixable:  true,
		})
	}

	// Validate license
	if metadata.License == "" {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template license is not specified",
			Rule:     "missing-license",
			Fixable:  true,
		})
	} else {
		validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "MPL-2.0"}
		if !tv.contains(validLicenses, metadata.License) {
			warnings = append(warnings, interfaces.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("License '%s' is not a common open source license", metadata.License),
				Rule:     "uncommon-license",
				Fixable:  true,
			})
		}
	}

	// Validate repository URL
	if metadata.Repository != "" && !tv.isValidURL(metadata.Repository) {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Repository URL appears to be invalid",
			Rule:     "invalid-repository-url",
			Fixable:  true,
		})
	}

	// Validate homepage URL
	if metadata.Homepage != "" && !tv.isValidURL(metadata.Homepage) {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Homepage URL appears to be invalid",
			Rule:     "invalid-homepage-url",
			Fixable:  true,
		})
	}

	// Validate keywords
	if len(metadata.Keywords) == 0 {
		warnings = append(warnings, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "No keywords specified for better discoverability",
			Rule:     "missing-keywords",
			Fixable:  true,
		})
	} else {
		for _, keyword := range metadata.Keywords {
			if keyword == "" {
				warnings = append(warnings, interfaces.ValidationIssue{
					Type:     "warning",
					Severity: "warning",
					Message:  "Empty keyword found",
					Rule:     "empty-keyword",
					Fixable:  true,
				})
			}
		}
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(issues) == 0,
		Issues:   issues,
		Warnings: warnings,
		Summary: interfaces.ValidationSummary{
			FixableCount: 0,
			ErrorCount:   len(issues),
			WarningCount: len(warnings),
		},
	}, nil
}

// Helper methods for validation

// isTemplateDependency checks if a dependency is a template dependency
func (tv *TemplateValidator) isTemplateDependency(dependency string) bool {
	// Simple heuristic: if it doesn't contain version specifiers or URLs, it might be a template
	return !strings.Contains(dependency, "@") && !strings.Contains(dependency, "://") && !strings.Contains(dependency, ".")
}

// isValidDependencyFormat checks if dependency format is valid
func (tv *TemplateValidator) isValidDependencyFormat(dependency string) bool {
	// Basic validation - not empty and doesn't contain invalid characters
	if dependency == "" {
		return false
	}

	// Check for common invalid characters
	invalidChars := []string{" ", "\t", "\n", "\r"}
	for _, char := range invalidChars {
		if strings.Contains(dependency, char) {
			return false
		}
	}

	return true
}

// detectCircularDependencies detects circular dependencies (simplified implementation)
func (tv *TemplateValidator) detectCircularDependencies(templateInfo *interfaces.TemplateInfo) []string {
	var circular []string

	// Simple check: if template depends on itself
	for _, dep := range templateInfo.Dependencies {
		if dep == templateInfo.Name {
			circular = append(circular, dep)
		}
	}

	// In a full implementation, you would do a deep dependency graph analysis
	return circular
}

// isValidVersionFormat checks if version format is valid (simplified semver check)
func (tv *TemplateValidator) isValidVersionFormat(version string) bool {
	if version == "" {
		return false
	}

	// Basic semver pattern check
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	// Check if parts are numeric (simplified)
	for _, part := range parts {
		if part == "" {
			return false
		}
		// In a full implementation, you would use proper semver validation
	}

	return true
}

// isValidURL checks if URL format is valid (simplified)
func (tv *TemplateValidator) isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
