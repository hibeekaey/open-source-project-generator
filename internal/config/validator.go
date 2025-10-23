package config

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Validator validates project configurations
type Validator struct {
	schema *Schema
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{
		schema: DefaultSchema(),
	}
}

// NewValidatorWithSchema creates a validator with a custom schema
func NewValidatorWithSchema(schema *Schema) *Validator {
	return &Validator{
		schema: schema,
	}
}

// Validate validates a complete project configuration
func (v *Validator) Validate(config *models.ProjectConfig) error {
	// Validate high-level project configuration
	if err := ValidateProjectConfig(config, v.schema); err != nil {
		return fmt.Errorf("project configuration validation failed: %w", err)
	}

	// Validate each component
	for i, comp := range config.Components {
		if err := v.ValidateComponent(&comp); err != nil {
			return fmt.Errorf("component %d (%s) validation failed: %w", i, comp.Name, err)
		}
	}

	// Validate integration configuration
	if err := v.ValidateIntegration(&config.Integration); err != nil {
		return fmt.Errorf("integration configuration validation failed: %w", err)
	}

	// Validate options
	if err := v.ValidateOptions(&config.Options); err != nil {
		return fmt.Errorf("options validation failed: %w", err)
	}

	// Check for duplicate component names
	if err := v.checkDuplicateNames(config.Components); err != nil {
		return err
	}

	return nil
}

// ValidateComponent validates a single component configuration
func (v *Validator) ValidateComponent(comp *models.ComponentConfig) error {
	// Check if component type is supported
	if !v.schema.IsComponentTypeSupported(comp.Type) {
		return fmt.Errorf("unsupported component type: %s (supported types: %s)",
			comp.Type, strings.Join(v.schema.SupportedComponentTypes, ", "))
	}

	// Validate component name
	if comp.Name == "" {
		return fmt.Errorf("component name is required")
	}

	if err := validateProjectName(comp.Name); err != nil {
		return fmt.Errorf("invalid component name: %w", err)
	}

	// Get component schema
	compSchema, err := v.schema.GetComponentSchema(comp.Type)
	if err != nil {
		return err
	}

	// Validate component-specific configuration
	if err := compSchema.ValidateConfig(comp.Config); err != nil {
		return fmt.Errorf("component configuration validation failed: %w", err)
	}

	return nil
}

// ValidateIntegration validates integration configuration
func (v *Validator) ValidateIntegration(integration *models.IntegrationConfig) error {
	// Validate API endpoints
	for name, endpoint := range integration.APIEndpoints {
		if err := validateEndpoint(name, endpoint); err != nil {
			return fmt.Errorf("invalid API endpoint '%s': %w", name, err)
		}
	}

	// Validate shared environment variables
	for key, value := range integration.SharedEnvironment {
		if err := validateEnvironmentVariable(key, value); err != nil {
			return fmt.Errorf("invalid environment variable '%s': %w", key, err)
		}
	}

	return nil
}

// ValidateOptions validates project options
func (v *Validator) ValidateOptions(options *models.ProjectOptions) error {
	// Validate option combinations
	if options.DryRun && options.CreateBackup {
		// This is fine, but backup won't actually be created in dry-run
		// Just a logical check, not an error
	}

	if options.ForceOverwrite && !options.CreateBackup {
		// Warn but don't fail - user explicitly chose this
	}

	return nil
}

// checkDuplicateNames checks for duplicate component names
func (v *Validator) checkDuplicateNames(components []models.ComponentConfig) error {
	names := make(map[string]bool)

	for _, comp := range components {
		if !comp.Enabled {
			continue
		}

		if names[comp.Name] {
			return fmt.Errorf("duplicate component name: %s", comp.Name)
		}
		names[comp.Name] = true
	}

	return nil
}

// validateEndpoint validates an API endpoint URL
func validateEndpoint(name, endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	// Basic URL validation
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("endpoint must start with http:// or https://")
	}

	// Check for invalid characters
	invalidChars := []string{" ", "\"", "<", ">", "{", "}"}
	for _, char := range invalidChars {
		if strings.Contains(endpoint, char) {
			return fmt.Errorf("contains invalid character: %s", char)
		}
	}

	return nil
}

// validateEnvironmentVariable validates an environment variable key and value
func validateEnvironmentVariable(key, value string) error {
	// Validate key
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Environment variable keys should be uppercase with underscores
	for _, char := range key {
		if !isUppercase(char) && !isDigit(char) && char != '_' {
			return fmt.Errorf("key '%s' contains invalid character: %c (use uppercase letters, digits, and underscores)", key, char)
		}
	}

	// Key cannot start with a digit
	if isDigit(rune(key[0])) {
		return fmt.Errorf("key '%s' cannot start with a digit", key)
	}

	// Value can be empty (for optional variables)
	// No specific validation needed for value

	return nil
}

// ApplyDefaults applies default values to component configurations
func (v *Validator) ApplyDefaults(config *models.ProjectConfig) error {
	for i := range config.Components {
		comp := &config.Components[i]

		// Get component schema
		compSchema, err := v.schema.GetComponentSchema(comp.Type)
		if err != nil {
			return fmt.Errorf("failed to get schema for component %s: %w", comp.Name, err)
		}

		// Apply defaults
		comp.Config = compSchema.ApplyDefaults(comp.Config)
	}

	return nil
}

// ValidateAndApplyDefaults validates configuration and applies defaults
func (v *Validator) ValidateAndApplyDefaults(config *models.ProjectConfig) error {
	// Apply defaults first
	if err := v.ApplyDefaults(config); err != nil {
		return fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Then validate
	if err := v.Validate(config); err != nil {
		return err
	}

	return nil
}

// GetSupportedComponentTypes returns all supported component types
func (v *Validator) GetSupportedComponentTypes() []string {
	return v.schema.SupportedComponentTypes
}

// GetComponentSchema returns the schema for a specific component type
func (v *Validator) GetComponentSchema(componentType string) (*ComponentSchema, error) {
	return v.schema.GetComponentSchema(componentType)
}

// Helper functions

func isUppercase(char rune) bool {
	return char >= 'A' && char <= 'Z'
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// ValidationReport contains detailed validation results
type ValidationReport struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []string
}

// AddError adds an error to the validation report
func (vr *ValidationReport) AddError(field, message string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// AddWarning adds a warning to the validation report
func (vr *ValidationReport) AddWarning(message string) {
	vr.Warnings = append(vr.Warnings, message)
}

// HasErrors returns true if there are any errors
func (vr *ValidationReport) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (vr *ValidationReport) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// String returns a formatted string representation of the validation report
func (vr *ValidationReport) String() string {
	var builder strings.Builder

	if vr.Valid {
		builder.WriteString("Validation passed")
	} else {
		builder.WriteString("Validation failed")
	}

	if len(vr.Errors) > 0 {
		builder.WriteString("\n\nErrors:\n")
		for i, err := range vr.Errors {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
		}
	}

	if len(vr.Warnings) > 0 {
		builder.WriteString("\nWarnings:\n")
		for i, warning := range vr.Warnings {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, warning))
		}
	}

	return builder.String()
}

// ValidateWithReport validates configuration and returns a detailed report
func (v *Validator) ValidateWithReport(config *models.ProjectConfig) *ValidationReport {
	report := &ValidationReport{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []string{},
	}

	// Validate and collect all errors
	if err := v.Validate(config); err != nil {
		report.Valid = false
		// Try to extract field-specific errors
		if valErr, ok := err.(*ValidationError); ok {
			report.Errors = append(report.Errors, *valErr)
		} else {
			report.Errors = append(report.Errors, ValidationError{
				Field:   "",
				Message: err.Error(),
			})
		}
	}

	// Add warnings for potentially problematic configurations
	if config.Options.ForceOverwrite && !config.Options.CreateBackup {
		report.AddWarning("Force overwrite is enabled without backup - existing files will be permanently lost")
	}

	if len(config.Components) > 10 {
		report.AddWarning(fmt.Sprintf("Large number of components (%d) may result in long generation time", len(config.Components)))
	}

	// Check for components without integration
	if len(config.Components) > 1 && !config.Integration.GenerateDockerCompose && !config.Integration.GenerateScripts {
		report.AddWarning("Multiple components without integration configuration - consider enabling Docker Compose or scripts generation")
	}

	return report
}
