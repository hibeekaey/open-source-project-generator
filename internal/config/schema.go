package config

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Schema defines the validation schema for project configuration
type Schema struct {
	// SupportedComponentTypes lists all valid component types
	SupportedComponentTypes []string

	// RequiredFields defines which fields are required in the configuration
	RequiredFields []string

	// ComponentTypeSchemas defines validation rules for each component type
	ComponentTypeSchemas map[string]*ComponentSchema
}

// ComponentSchema defines validation rules for a specific component type
type ComponentSchema struct {
	Type           string
	RequiredConfig []string                            // Required config keys
	OptionalConfig []string                            // Optional config keys
	ConfigDefaults map[string]interface{}              // Default values for config
	Validators     map[string]ConfigFieldValidatorFunc // Custom validators for config fields
}

// ConfigFieldValidatorFunc is a function that validates a specific config field
type ConfigFieldValidatorFunc func(value interface{}) error

// DefaultSchema returns the default validation schema
func DefaultSchema() *Schema {
	return &Schema{
		SupportedComponentTypes: []string{
			"nextjs",
			"go-backend",
			"android",
			"ios",
			"docker",
			"kubernetes",
			"terraform",
		},
		RequiredFields: []string{
			"name",
			"output_dir",
		},
		ComponentTypeSchemas: map[string]*ComponentSchema{
			"nextjs": {
				Type:           "nextjs",
				RequiredConfig: []string{"name"},
				OptionalConfig: []string{"typescript", "tailwind", "app_router", "eslint"},
				ConfigDefaults: map[string]interface{}{
					"typescript": true,
					"tailwind":   true,
					"app_router": true,
					"eslint":     true,
				},
				Validators: map[string]ConfigFieldValidatorFunc{
					"name": validateProjectName,
				},
			},
			"go-backend": {
				Type:           "go-backend",
				RequiredConfig: []string{"name", "module"},
				OptionalConfig: []string{"framework", "port"},
				ConfigDefaults: map[string]interface{}{
					"framework": "gin",
					"port":      8080,
				},
				Validators: map[string]ConfigFieldValidatorFunc{
					"name":   validateProjectName,
					"module": validateGoModule,
					"port":   validatePort,
				},
			},
			"android": {
				Type:           "android",
				RequiredConfig: []string{"name", "package"},
				OptionalConfig: []string{"min_sdk", "target_sdk", "language"},
				ConfigDefaults: map[string]interface{}{
					"min_sdk":    24,
					"target_sdk": 34,
					"language":   "kotlin",
				},
				Validators: map[string]ConfigFieldValidatorFunc{
					"name":    validateProjectName,
					"package": validateAndroidPackage,
				},
			},
			"ios": {
				Type:           "ios",
				RequiredConfig: []string{"name", "bundle_id"},
				OptionalConfig: []string{"deployment_target", "language"},
				ConfigDefaults: map[string]interface{}{
					"deployment_target": "15.0",
					"language":          "swift",
				},
				Validators: map[string]ConfigFieldValidatorFunc{
					"name":      validateProjectName,
					"bundle_id": validateBundleID,
				},
			},
			"docker": {
				Type:           "docker",
				RequiredConfig: []string{},
				OptionalConfig: []string{"compose_version"},
				ConfigDefaults: map[string]interface{}{
					"compose_version": "3.8",
				},
			},
			"kubernetes": {
				Type:           "kubernetes",
				RequiredConfig: []string{},
				OptionalConfig: []string{"namespace", "replicas"},
				ConfigDefaults: map[string]interface{}{
					"namespace": "default",
					"replicas":  1,
				},
			},
			"terraform": {
				Type:           "terraform",
				RequiredConfig: []string{},
				OptionalConfig: []string{"provider", "region"},
				ConfigDefaults: map[string]interface{}{
					"provider": "aws",
					"region":   "us-east-1",
				},
			},
		},
	}
}

// GetComponentSchema returns the schema for a specific component type
func (s *Schema) GetComponentSchema(componentType string) (*ComponentSchema, error) {
	schema, exists := s.ComponentTypeSchemas[componentType]
	if !exists {
		return nil, fmt.Errorf("no schema defined for component type: %s", componentType)
	}
	return schema, nil
}

// IsComponentTypeSupported checks if a component type is supported
func (s *Schema) IsComponentTypeSupported(componentType string) bool {
	for _, supported := range s.SupportedComponentTypes {
		if supported == componentType {
			return true
		}
	}
	return false
}

// ApplyDefaults applies default values to a component configuration
func (cs *ComponentSchema) ApplyDefaults(config map[string]interface{}) map[string]interface{} {
	if config == nil {
		config = make(map[string]interface{})
	}

	// Apply defaults for missing keys
	for key, defaultValue := range cs.ConfigDefaults {
		if _, exists := config[key]; !exists {
			config[key] = defaultValue
		}
	}

	return config
}

// ValidateConfig validates a component configuration against the schema
func (cs *ComponentSchema) ValidateConfig(config map[string]interface{}) error {
	// Check required fields
	for _, required := range cs.RequiredConfig {
		if _, exists := config[required]; !exists {
			return fmt.Errorf("required config field missing: %s", required)
		}
	}

	// Validate each field with custom validators
	for key, value := range config {
		if validator, exists := cs.Validators[key]; exists {
			if err := validator(value); err != nil {
				return fmt.Errorf("validation failed for field '%s': %w", key, err)
			}
		}
	}

	return nil
}

// Built-in validators

// validateProjectName validates a project name
func validateProjectName(value interface{}) error {
	name, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(name) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("cannot exceed 100 characters")
	}

	// Check for valid characters (alphanumeric, dash, underscore)
	for _, char := range name {
		if !isAlphanumeric(char) && char != '-' && char != '_' {
			return fmt.Errorf("contains invalid character: %c (only alphanumeric, dash, and underscore allowed)", char)
		}
	}

	// Cannot start or end with special characters
	if name[0] == '-' || name[0] == '_' || name[len(name)-1] == '-' || name[len(name)-1] == '_' {
		return fmt.Errorf("cannot start or end with dash or underscore")
	}

	return nil
}

// validateGoModule validates a Go module name
func validateGoModule(value interface{}) error {
	module, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(module) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Basic validation for Go module path format
	if !strings.Contains(module, "/") {
		return fmt.Errorf("must be a valid module path (e.g., github.com/user/project)")
	}

	// Check for invalid characters
	invalidChars := []string{" ", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(module, char) {
			return fmt.Errorf("contains invalid character: %s", char)
		}
	}

	return nil
}

// validatePort validates a port number
func validatePort(value interface{}) error {
	// Handle both int and float64 (JSON unmarshaling can produce float64)
	var port int
	switch v := value.(type) {
	case int:
		port = v
	case float64:
		port = int(v)
	default:
		return fmt.Errorf("must be a number")
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("must be between 1 and 65535")
	}

	// Warn about privileged ports (< 1024) but don't fail
	// This is just validation, actual usage might require sudo

	return nil
}

// validateAndroidPackage validates an Android package name
func validateAndroidPackage(value interface{}) error {
	pkg, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(pkg) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Must contain at least one dot
	if !strings.Contains(pkg, ".") {
		return fmt.Errorf("must be a valid package name (e.g., com.example.app)")
	}

	// Check each segment
	segments := strings.Split(pkg, ".")
	if len(segments) < 2 {
		return fmt.Errorf("must have at least two segments (e.g., com.example)")
	}

	for _, segment := range segments {
		if len(segment) == 0 {
			return fmt.Errorf("segments cannot be empty")
		}

		// Segments must start with a letter
		if !isLetter(rune(segment[0])) {
			return fmt.Errorf("segment '%s' must start with a letter", segment)
		}

		// Segments can only contain letters, digits, and underscores
		for _, char := range segment {
			if !isAlphanumeric(char) && char != '_' {
				return fmt.Errorf("segment '%s' contains invalid character: %c", segment, char)
			}
		}
	}

	return nil
}

// validateBundleID validates an iOS bundle identifier
func validateBundleID(value interface{}) error {
	bundleID, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(bundleID) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Must contain at least one dot
	if !strings.Contains(bundleID, ".") {
		return fmt.Errorf("must be a valid bundle identifier (e.g., com.example.app)")
	}

	// Check each segment
	segments := strings.Split(bundleID, ".")
	if len(segments) < 2 {
		return fmt.Errorf("must have at least two segments (e.g., com.example)")
	}

	for _, segment := range segments {
		if len(segment) == 0 {
			return fmt.Errorf("segments cannot be empty")
		}

		// Segments must start with a letter
		if !isLetter(rune(segment[0])) {
			return fmt.Errorf("segment '%s' must start with a letter", segment)
		}

		// Segments can only contain letters, digits, and hyphens
		for _, char := range segment {
			if !isAlphanumeric(char) && char != '-' {
				return fmt.Errorf("segment '%s' contains invalid character: %c", segment, char)
			}
		}
	}

	return nil
}

// Helper functions

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func isLetter(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidateProjectConfig performs high-level validation on the entire project configuration
func ValidateProjectConfig(config *models.ProjectConfig, schema *Schema) error {
	// Check required fields
	for _, field := range schema.RequiredFields {
		switch field {
		case "name":
			if config.Name == "" {
				return NewValidationError("name", "project name is required")
			}
		case "output_dir":
			if config.OutputDir == "" {
				return NewValidationError("output_dir", "output directory is required")
			}
		}
	}

	// Validate project name
	if err := validateProjectName(config.Name); err != nil {
		return NewValidationError("name", err.Error())
	}

	// Validate output directory path
	if err := validateOutputDir(config.OutputDir); err != nil {
		return NewValidationError("output_dir", err.Error())
	}

	// Validate that at least one component is enabled
	hasEnabledComponent := false
	for _, comp := range config.Components {
		if comp.Enabled {
			hasEnabledComponent = true
			break
		}
	}
	if !hasEnabledComponent {
		return NewValidationError("components", "at least one component must be enabled")
	}

	return nil
}

// validateOutputDir validates an output directory path
func validateOutputDir(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for invalid characters
	invalidChars := []string{"*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("contains invalid character: %s", char)
		}
	}

	return nil
}
