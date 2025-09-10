package models

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ConfigValidator handles configuration validation
type ConfigValidator struct {
	validator *validator.Validate
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	v := validator.New()

	// Register custom validation functions
	v.RegisterValidation("semver", validateSemVer)
	v.RegisterValidation("alphanum", validateAlphaNum)

	return &ConfigValidator{
		validator: v,
	}
}

// ValidateProjectConfig validates a project configuration
func (cv *ConfigValidator) ValidateProjectConfig(config *ProjectConfig) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	// Perform struct validation
	if err := cv.validator.Struct(config); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldError.Field(),
					Tag:     fieldError.Tag(),
					Value:   fmt.Sprintf("%v", fieldError.Value()),
					Message: cv.getErrorMessage(fieldError),
				})
			}
		}
	}

	// Custom validation logic
	cv.validateComponentDependencies(config, result)
	cv.validateVersionCompatibility(config, result)
	cv.addWarnings(config, result)

	return result
}

// ValidateVersionConfig validates version configuration
func (cv *ConfigValidator) ValidateVersionConfig(config *VersionConfig) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	if err := cv.validator.Struct(config); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldError.Field(),
					Tag:     fieldError.Tag(),
					Value:   fmt.Sprintf("%v", fieldError.Value()),
					Message: cv.getErrorMessage(fieldError),
				})
			}
		}
	}

	return result
}

// ValidateTemplateMetadata validates template metadata
func (cv *ConfigValidator) ValidateTemplateMetadata(metadata *TemplateMetadata) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	if err := cv.validator.Struct(metadata); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldError.Field(),
					Tag:     fieldError.Tag(),
					Value:   fmt.Sprintf("%v", fieldError.Value()),
					Message: cv.getErrorMessage(fieldError),
				})
			}
		}
	}

	// Validate template variables
	for i, variable := range metadata.Variables {
		if varResult := cv.ValidateTemplateVar(&variable); !varResult.Valid {
			for _, err := range varResult.Errors {
				err.Field = fmt.Sprintf("Variables[%d].%s", i, err.Field)
				result.Errors = append(result.Errors, err)
			}
			result.Valid = false
		}
	}

	return result
}

// ValidateTemplateVar validates a template variable
func (cv *ConfigValidator) ValidateTemplateVar(variable *TemplateVar) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
	}

	if err := cv.validator.Struct(variable); err != nil {
		result.Valid = false
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldError.Field(),
					Tag:     fieldError.Tag(),
					Value:   fmt.Sprintf("%v", fieldError.Value()),
					Message: cv.getErrorMessage(fieldError),
				})
			}
		}
	}

	// Validate default value type matches declared type
	if variable.Default != nil {
		if !cv.validateValueType(variable.Default, variable.Type) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "Default",
				Tag:     "type_mismatch",
				Value:   fmt.Sprintf("%v", variable.Default),
				Message: fmt.Sprintf("Default value type does not match declared type '%s'", variable.Type),
			})
		}
	}

	return result
}

// validateComponentDependencies validates component dependencies
func (cv *ConfigValidator) validateComponentDependencies(config *ProjectConfig, result *ValidationResult) {
	// Check if at least one component is selected
	hasAnyComponent := config.Components.Frontend.MainApp || config.Components.Frontend.Home ||
		config.Components.Frontend.Admin || config.Components.Backend.API ||
		config.Components.Mobile.Android || config.Components.Mobile.IOS ||
		config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform

	if !hasAnyComponent {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Components",
			Tag:     "required",
			Value:   "",
			Message: "At least one component must be selected",
		})
	}

	// Validate component dependencies
	if config.Components.Infrastructure.Kubernetes && !config.Components.Infrastructure.Docker {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Components.Infrastructure.Kubernetes",
			Message: "Kubernetes deployment typically requires Docker containers",
		})
	}

	if (config.Components.Frontend.MainApp || config.Components.Frontend.Home || config.Components.Frontend.Admin) &&
		!config.Components.Backend.API {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Components.Frontend",
			Message: "Frontend applications typically require a backend API",
		})
	}
}

// validateVersionCompatibility validates version compatibility
func (cv *ConfigValidator) validateVersionCompatibility(config *ProjectConfig, result *ValidationResult) {
	if config.Versions == nil {
		return
	}

	// Add version compatibility checks here
	// For now, just validate that versions are not empty if components are selected
	if config.Components.Frontend.MainApp || config.Components.Frontend.Home || config.Components.Frontend.Admin {
		if config.Versions.Node == "" || config.Versions.NextJS == "" || config.Versions.React == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   "Versions",
				Message: "Frontend components require Node.js, Next.js, and React versions",
			})
		}
	}

	if config.Components.Backend.API && config.Versions.Go == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Versions.Go",
			Message: "Backend API component requires Go version",
		})
	}

	if config.Components.Mobile.Android && config.Versions.Kotlin == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Versions.Kotlin",
			Message: "Android component requires Kotlin version",
		})
	}

	if config.Components.Mobile.IOS && config.Versions.Swift == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Versions.Swift",
			Message: "iOS component requires Swift version",
		})
	}
}

// addWarnings adds additional warnings based on configuration
func (cv *ConfigValidator) addWarnings(config *ProjectConfig, result *ValidationResult) {
	// Check for common issues
	if config.Name == config.Organization {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Name",
			Message: "Project name is the same as organization name",
		})
	}

	if len(config.Description) < 20 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Description",
			Message: "Description is quite short, consider adding more details",
		})
	}
}

// validateValueType validates that a value matches the expected type
func (cv *ConfigValidator) validateValueType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "int":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		}
		return false
	case "float":
		switch value.(type) {
		case float32, float64:
			return true
		}
		return false
	case "bool":
		_, ok := value.(bool)
		return ok
	case "array":
		switch value.(type) {
		case []interface{}, []string, []int, []float64, []bool:
			return true
		}
		return false
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return false
	}
}

// getErrorMessage returns a human-readable error message for validation errors
func (cv *ConfigValidator) getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "semver":
		return fmt.Sprintf("%s must be a valid semantic version", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", fe.Field())
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", fe.Field(), fe.Tag())
	}
}

// Custom validation functions

// validateSemVer validates semantic version format
func validateSemVer(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	if version == "" {
		return true // Allow empty for optional fields
	}

	// Basic semantic version regex pattern
	semverPattern := `^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`
	matched, _ := regexp.MatchString(semverPattern, version)
	return matched
}

// validateAlphaNum validates alphanumeric characters with hyphens and underscores
func validateAlphaNum(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	// Allow alphanumeric characters, hyphens, and underscores
	pattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}
