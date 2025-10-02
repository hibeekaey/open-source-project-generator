package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Pre-compiled regular expressions for engine helpers validation
var (
	projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9]$`)
)

// validateProjectStructureBasic performs basic project structure validation
func (e *Engine) validateProjectStructureBasic(path string, result *models.ValidationResult) error {
	// Check for README file
	readmeFiles := []string{"README.md", "README.txt", "README.rst", "README"}
	readmeFound := false
	for _, readme := range readmeFiles {
		if _, err := os.Stat(filepath.Join(path, readme)); err == nil {
			readmeFound = true
			break
		}
	}
	if !readmeFound {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "README file is missing",
			File:     path,
			Rule:     "structure.readme.required",
			Fixable:  true,
		})
	}

	// Check for LICENSE file
	licenseFiles := []string{"LICENSE", "LICENSE.txt", "LICENSE.md", "COPYING"}
	licenseFound := false
	for _, license := range licenseFiles {
		if _, err := os.Stat(filepath.Join(path, license)); err == nil {
			licenseFound = true
			break
		}
	}
	if !licenseFound {
		result.Valid = false
		result.Issues = append(result.Issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "LICENSE file is missing",
			File:     path,
			Rule:     "structure.license.required",
			Fixable:  true,
		})
	}

	return nil
}

// validateProjectDependenciesBasic performs basic dependency validation
func (e *Engine) validateProjectDependenciesBasic(path string, result *models.ValidationResult) error {
	// Validate package.json if it exists
	packageJsonPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJsonPath); err == nil {
		if err := e.ValidatePackageJSON(packageJsonPath); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid package.json: %v", err),
				File:     packageJsonPath,
				Rule:     "dependencies.package_json.valid",
				Fixable:  false,
			})
		}
	}

	// Validate go.mod if it exists
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		if err := e.ValidateGoMod(goModPath); err != nil {
			result.Valid = false
			result.Issues = append(result.Issues, models.ValidationIssue{
				Type:     "error",
				Severity: "error",
				Message:  fmt.Sprintf("Invalid go.mod: %v", err),
				File:     goModPath,
				Rule:     "dependencies.go_mod.valid",
				Fixable:  false,
			})
		}
	}

	return nil
}

// Helper functions for configuration validation
func (e *Engine) validateRequiredConfigFields(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	result.Summary.TotalProperties++

	if config.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "name",
			Value:    "",
			Type:     "required",
			Message:  "Project name is required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.name.required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	} else {
		result.Summary.ValidProperties++
	}

	result.Summary.TotalProperties++
	if config.Organization == "" {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "organization",
			Value:    "",
			Type:     "required",
			Message:  "Organization is required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.organization.required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	} else {
		result.Summary.ValidProperties++
	}

	// OutputPath is optional for basic config validation
	// It's only required when actually generating a project
	result.Summary.TotalProperties++
	if config.OutputPath != "" {
		result.Summary.ValidProperties++
	}
	// Note: OutputPath validation can be added separately for generation-time validation
}

// validateConfigFieldFormats validates configuration field formats
func (e *Engine) validateConfigFieldFormats(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	// Validate project name format
	if config.Name != "" {
		if !projectNameRegex.MatchString(config.Name) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "name",
				Value:      config.Name,
				Type:       "format",
				Message:    "Project name should contain only alphanumeric characters, hyphens, and underscores",
				Suggestion: "Use a name like 'my-project' or 'my_project'",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.name.format",
			})
			result.Summary.WarningCount++
		}
	}

	// Validate output path format
	if config.OutputPath != "" {
		if strings.Contains(config.OutputPath, "..") {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    "output_path",
				Value:    config.OutputPath,
				Type:     "security",
				Message:  "Output path cannot contain '..' for security reasons",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "config.output_path.security",
			})
			result.Valid = false
			result.Summary.ErrorCount++
		}
	}
}

// validateComponentConfiguration validates component configuration
func (e *Engine) validateComponentConfiguration(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	// Basic component validation - check if any components are configured
	if !config.Components.Frontend.NextJS.App && !config.Components.Backend.GoGin {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "components",
			Value:      "",
			Type:       "missing",
			Message:    "No components are enabled in the configuration",
			Suggestion: "Enable at least one component (frontend or backend)",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "config.components.empty",
		})
		result.Summary.WarningCount++
	}
}
