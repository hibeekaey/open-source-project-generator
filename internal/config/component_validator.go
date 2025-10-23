package config

import (
	"fmt"
)

// ConfigValidator defines the interface for component config validation
type ConfigValidator interface {
	// Validate validates the configuration
	Validate(config map[string]interface{}) error

	// GetRequiredFields returns required configuration fields
	GetRequiredFields() []string

	// GetOptionalFields returns optional configuration fields
	GetOptionalFields() []string

	// GetFieldDescription returns description for a field
	GetFieldDescription(field string) string
}

// ComponentConfigValidator validates component-specific configurations
type ComponentConfigValidator struct {
	validators map[string]ConfigValidator
}

// NewComponentConfigValidator creates a new component config validator
func NewComponentConfigValidator() *ComponentConfigValidator {
	return &ComponentConfigValidator{
		validators: make(map[string]ConfigValidator),
	}
}

// RegisterValidator registers a validator for a component type
func (ccv *ComponentConfigValidator) RegisterValidator(componentType string, validator ConfigValidator) {
	ccv.validators[componentType] = validator
}

// Validate validates configuration for a specific component type
func (ccv *ComponentConfigValidator) Validate(componentType string, config map[string]interface{}) error {
	validator, exists := ccv.validators[componentType]
	if !exists {
		// No specific validator registered, skip validation
		return nil
	}

	return validator.Validate(config)
}

// GetValidator returns the validator for a specific component type
func (ccv *ComponentConfigValidator) GetValidator(componentType string) (ConfigValidator, bool) {
	validator, exists := ccv.validators[componentType]
	return validator, exists
}

// HasValidator checks if a validator is registered for a component type
func (ccv *ComponentConfigValidator) HasValidator(componentType string) bool {
	_, exists := ccv.validators[componentType]
	return exists
}

// GetRequiredFields returns required fields for a component type
func (ccv *ComponentConfigValidator) GetRequiredFields(componentType string) []string {
	validator, exists := ccv.validators[componentType]
	if !exists {
		return []string{}
	}
	return validator.GetRequiredFields()
}

// GetOptionalFields returns optional fields for a component type
func (ccv *ComponentConfigValidator) GetOptionalFields(componentType string) []string {
	validator, exists := ccv.validators[componentType]
	if !exists {
		return []string{}
	}
	return validator.GetOptionalFields()
}

// GetFieldDescription returns description for a field in a component type
func (ccv *ComponentConfigValidator) GetFieldDescription(componentType, field string) string {
	validator, exists := ccv.validators[componentType]
	if !exists {
		return ""
	}
	return validator.GetFieldDescription(field)
}

// ValidateWithDetails validates and returns detailed error information
func (ccv *ComponentConfigValidator) ValidateWithDetails(componentType string, config map[string]interface{}) *ComponentValidationResult {
	result := &ComponentValidationResult{
		ComponentType: componentType,
		Valid:         true,
		Errors:        []FieldError{},
	}

	validator, exists := ccv.validators[componentType]
	if !exists {
		// No validator registered, consider valid
		return result
	}

	// Check required fields
	for _, field := range validator.GetRequiredFields() {
		if _, exists := config[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, FieldError{
				Field:       field,
				Message:     "required field is missing",
				Description: validator.GetFieldDescription(field),
			})
		}
	}

	// Validate the configuration
	if err := validator.Validate(config); err != nil {
		result.Valid = false
		// Try to extract field-specific errors
		if fieldErr, ok := err.(*FieldError); ok {
			result.Errors = append(result.Errors, *fieldErr)
		} else {
			result.Errors = append(result.Errors, FieldError{
				Field:   "",
				Message: err.Error(),
			})
		}
	}

	return result
}

// ComponentValidationResult contains detailed validation results for a component
type ComponentValidationResult struct {
	ComponentType string
	Valid         bool
	Errors        []FieldError
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field       string
	Message     string
	Description string
}

// Error implements the error interface
func (fe *FieldError) Error() string {
	if fe.Field != "" {
		return fmt.Sprintf("field '%s': %s", fe.Field, fe.Message)
	}
	return fe.Message
}

// NewFieldError creates a new field error
func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// NewFieldErrorWithDescription creates a new field error with description
func NewFieldErrorWithDescription(field, message, description string) *FieldError {
	return &FieldError{
		Field:       field,
		Message:     message,
		Description: description,
	}
}
