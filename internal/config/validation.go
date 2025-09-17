package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ConfigValidator provides comprehensive configuration validation
type ConfigValidator struct {
	schema *interfaces.ConfigSchema
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(schema *interfaces.ConfigSchema) *ConfigValidator {
	return &ConfigValidator{
		schema: schema,
	}
}

// ValidateProjectConfig validates a complete project configuration
func (v *ConfigValidator) ValidateProjectConfig(config *models.ProjectConfig) (*interfaces.ConfigValidationResult, error) {
	if config == nil {
		return &interfaces.ConfigValidationResult{
			Valid: false,
			Errors: []interfaces.ConfigValidationError{
				{
					Field:    "config",
					Type:     "null",
					Message:  "configuration cannot be null",
					Severity: "error",
					Rule:     "required",
				},
			},
			Summary: interfaces.ConfigValidationSummary{
				ErrorCount: 1,
			},
		}, nil
	}

	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary:  interfaces.ConfigValidationSummary{},
	}

	// Validate basic fields
	v.validateBasicFields(config, result)

	// Validate components
	v.validateComponents(&config.Components, result)

	// Validate versions
	if config.Versions != nil {
		v.validateVersions(config.Versions, result)
	}

	// Validate output path
	v.validateOutputPath(config.OutputPath, result)

	// Calculate summary
	v.calculateSummary(result)

	// Set overall validity
	result.Valid = result.Summary.ErrorCount == 0

	return result, nil
}

// validateBasicFields validates basic project configuration fields
func (v *ConfigValidator) validateBasicFields(config *models.ProjectConfig, result *interfaces.ConfigValidationResult) {
	// Validate name
	if err := v.validateField("name", config.Name, "string", true); err != nil {
		result.Errors = append(result.Errors, *err)
	} else if config.Name != "" {
		// Additional name validation
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(config.Name) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "name",
				Value:      config.Name,
				Type:       "pattern",
				Message:    "project name can only contain letters, numbers, underscores, and hyphens",
				Suggestion: "use only alphanumeric characters, underscores, and hyphens",
				Severity:   "error",
				Rule:       "pattern",
			})
		}
	}

	// Validate organization
	if err := v.validateField("organization", config.Organization, "string", true); err != nil {
		result.Errors = append(result.Errors, *err)
	}

	// Validate description (optional)
	if config.Description != "" {
		if len(config.Description) > 500 {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "description",
				Value:      config.Description,
				Type:       "length",
				Message:    "description is very long, consider shortening it",
				Suggestion: "keep description under 500 characters",
				Severity:   "warning",
				Rule:       "max_length",
			})
		}
	}

	// Validate license
	if config.License != "" {
		validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"}
		valid := false
		for _, license := range validLicenses {
			if config.License == license {
				valid = true
				break
			}
		}
		if !valid {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "license",
				Value:      config.License,
				Type:       "enum",
				Message:    fmt.Sprintf("license '%s' is not in the list of recommended licenses", config.License),
				Suggestion: fmt.Sprintf("consider using one of: %s", strings.Join(validLicenses, ", ")),
				Severity:   "warning",
				Rule:       "enum",
			})
		}
	}

	// Validate email format
	if config.Email != "" {
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		if !regexp.MustCompile(emailPattern).MatchString(config.Email) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "email",
				Value:      config.Email,
				Type:       "pattern",
				Message:    "invalid email format",
				Suggestion: "provide a valid email address (e.g., user@example.com)",
				Severity:   "error",
				Rule:       "pattern",
			})
		}
	}

	// Validate repository URL
	if config.Repository != "" {
		urlPattern := `^https?://.*$`
		if !regexp.MustCompile(urlPattern).MatchString(config.Repository) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "repository",
				Value:      config.Repository,
				Type:       "pattern",
				Message:    "repository must be a valid HTTP/HTTPS URL",
				Suggestion: "provide a full URL starting with http:// or https://",
				Severity:   "error",
				Rule:       "pattern",
			})
		}
	}
}

// validateComponents validates component configuration
func (v *ConfigValidator) validateComponents(components *models.Components, result *interfaces.ConfigValidationResult) {
	// Check if at least one component is selected
	hasAnyComponent := components.Frontend.NextJS.App ||
		components.Frontend.NextJS.Home ||
		components.Frontend.NextJS.Admin ||
		components.Frontend.NextJS.Shared ||
		components.Backend.GoGin ||
		components.Mobile.Android ||
		components.Mobile.IOS ||
		components.Infrastructure.Docker ||
		components.Infrastructure.Kubernetes ||
		components.Infrastructure.Terraform

	if !hasAnyComponent {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "components",
			Type:       "selection",
			Message:    "no components selected, project will be minimal",
			Suggestion: "select at least one component to generate a functional project",
			Severity:   "warning",
			Rule:       "min_selection",
		})
	}

	// Validate component dependencies
	if components.Frontend.NextJS.Admin && !components.Frontend.NextJS.Shared {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "components.frontend.nextjs.admin",
			Type:       "dependency",
			Message:    "admin component works best with shared components",
			Suggestion: "consider enabling shared components for better code reuse",
			Severity:   "warning",
			Rule:       "dependency",
		})
	}

	if components.Infrastructure.Kubernetes && !components.Infrastructure.Docker {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "components.infrastructure.kubernetes",
			Type:       "dependency",
			Message:    "Kubernetes deployment typically requires Docker",
			Suggestion: "consider enabling Docker for containerization",
			Severity:   "warning",
			Rule:       "dependency",
		})
	}
}

// validateVersions validates version configuration
func (v *ConfigValidator) validateVersions(versions *models.VersionConfig, result *interfaces.ConfigValidationResult) {
	// Validate Node.js version format
	if versions.Node != "" {
		if !v.isValidSemVer(versions.Node) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "versions.node",
				Value:      versions.Node,
				Type:       "format",
				Message:    "invalid Node.js version format",
				Suggestion: "use semantic versioning format (e.g., 20.0.0)",
				Severity:   "error",
				Rule:       "semver",
			})
		}
	}

	// Validate Go version format
	if versions.Go != "" {
		if !v.isValidSemVer(versions.Go) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      "versions.go",
				Value:      versions.Go,
				Type:       "format",
				Message:    "invalid Go version format",
				Suggestion: "use semantic versioning format (e.g., 1.21.0)",
				Severity:   "error",
				Rule:       "semver",
			})
		}
	}

	// Validate package versions
	for pkg, version := range versions.Packages {
		if version != "" && !v.isValidSemVer(version) {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("versions.packages.%s", pkg),
				Value:      version,
				Type:       "format",
				Message:    fmt.Sprintf("invalid version format for package '%s'", pkg),
				Suggestion: "use semantic versioning format (e.g., 1.0.0)",
				Severity:   "error",
				Rule:       "semver",
			})
		}
	}
}

// validateOutputPath validates the output path
func (v *ConfigValidator) validateOutputPath(outputPath string, result *interfaces.ConfigValidationResult) {
	if outputPath == "" {
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:      "output_path",
			Type:       "required",
			Message:    "output path is required",
			Suggestion: "specify where to generate the project (e.g., ./my-project)",
			Severity:   "error",
			Rule:       "required",
		})
		return
	}

	// Check for potentially dangerous paths
	// Allow safe temporary directories and common development paths
	if v.isDangerousPath(outputPath) {
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:      "output_path",
			Value:      outputPath,
			Type:       "security",
			Message:    "output path points to a system directory",
			Suggestion: "use a safe directory like ./my-project or ~/projects/my-project",
			Severity:   "error",
			Rule:       "security",
		})
	}
}

// validateField validates a single field against schema rules
func (v *ConfigValidator) validateField(fieldName, value, expectedType string, required bool) *interfaces.ConfigValidationError {
	// Check required fields
	if required && value == "" {
		return &interfaces.ConfigValidationError{
			Field:      fieldName,
			Type:       "required",
			Message:    fmt.Sprintf("%s is required", fieldName),
			Suggestion: fmt.Sprintf("provide a value for %s", fieldName),
			Severity:   "error",
			Rule:       "required",
		}
	}

	// If not required and empty, skip further validation
	if !required && value == "" {
		return nil
	}

	// Get schema property if available
	if v.schema != nil {
		if prop, exists := v.schema.Properties[fieldName]; exists {
			// Length validation
			if prop.MinLength != nil && len(value) < *prop.MinLength {
				return &interfaces.ConfigValidationError{
					Field:      fieldName,
					Value:      value,
					Type:       "length",
					Message:    fmt.Sprintf("%s is too short (minimum %d characters)", fieldName, *prop.MinLength),
					Suggestion: fmt.Sprintf("provide at least %d characters", *prop.MinLength),
					Severity:   "error",
					Rule:       "min_length",
				}
			}

			if prop.MaxLength != nil && len(value) > *prop.MaxLength {
				return &interfaces.ConfigValidationError{
					Field:      fieldName,
					Value:      value,
					Type:       "length",
					Message:    fmt.Sprintf("%s is too long (maximum %d characters)", fieldName, *prop.MaxLength),
					Suggestion: fmt.Sprintf("keep under %d characters", *prop.MaxLength),
					Severity:   "error",
					Rule:       "max_length",
				}
			}

			// Pattern validation
			if prop.Pattern != "" {
				if matched, _ := regexp.MatchString(prop.Pattern, value); !matched {
					return &interfaces.ConfigValidationError{
						Field:      fieldName,
						Value:      value,
						Type:       "pattern",
						Message:    fmt.Sprintf("%s does not match required pattern", fieldName),
						Suggestion: "check the format requirements",
						Severity:   "error",
						Rule:       "pattern",
					}
				}
			}
		}
	}

	return nil
}

// isValidSemVer checks if a version string follows semantic versioning
func (v *ConfigValidator) isValidSemVer(version string) bool {
	// Basic semantic versioning pattern: major.minor.patch
	pattern := `^(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9\-\.]+)?(\+[a-zA-Z0-9\-\.]+)?$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

// calculateSummary calculates validation summary statistics
func (v *ConfigValidator) calculateSummary(result *interfaces.ConfigValidationResult) {
	result.Summary.ErrorCount = len(result.Errors)
	result.Summary.WarningCount = len(result.Warnings)

	// Count missing required fields
	for _, err := range result.Errors {
		if err.Rule == "required" {
			result.Summary.MissingRequired++
		}
	}

	// Estimate total and valid properties
	result.Summary.TotalProperties = 10 // Basic count of main config properties
	result.Summary.ValidProperties = result.Summary.TotalProperties - result.Summary.ErrorCount
	if result.Summary.ValidProperties < 0 {
		result.Summary.ValidProperties = 0
	}
}

// isDangerousPath checks if a path is potentially dangerous for project generation
func (v *ConfigValidator) isDangerousPath(outputPath string) bool {
	// Exact matches for root directories
	if outputPath == "/" {
		return true
	}

	// Dangerous system directories
	dangerousPaths := []string{"/usr", "/etc", "/bin", "/sbin", "/boot"}
	for _, dangerous := range dangerousPaths {
		if strings.HasPrefix(outputPath, dangerous+"/") || outputPath == dangerous {
			return true
		}
	}

	// Special handling for /var - allow temp directories but block others
	if outputPath == "/var" || strings.HasPrefix(outputPath, "/var/") {
		// Allow macOS temp directories and common temp paths
		safePaths := []string{"/var/folders/", "/var/tmp/"}
		for _, safe := range safePaths {
			if strings.HasPrefix(outputPath, safe) {
				return false
			}
		}
		// Block other /var paths
		return true
	}

	return false
}
