package config

import (
	"fmt"
	"strings"
)

// GoConfigValidator validates Go backend component configurations
type GoConfigValidator struct{}

// NewGoConfigValidator creates a new Go backend config validator
func NewGoConfigValidator() *GoConfigValidator {
	return &GoConfigValidator{}
}

// Validate validates the Go backend configuration
func (v *GoConfigValidator) Validate(config map[string]interface{}) error {
	// Validate name
	if name, exists := config["name"]; exists {
		if err := validateProjectName(name); err != nil {
			return NewFieldError("name", err.Error())
		}
	}

	// Validate module (required)
	if module, exists := config["module"]; exists {
		if err := v.ValidateModule(module); err != nil {
			return NewFieldError("module", err.Error())
		}
	}

	// Validate port (must be valid port number)
	if port, exists := config["port"]; exists {
		if err := v.ValidatePort(port); err != nil {
			return NewFieldError("port", err.Error())
		}
	}

	// Validate framework (must be supported)
	if framework, exists := config["framework"]; exists {
		if err := v.ValidateFramework(framework); err != nil {
			return NewFieldError("framework", err.Error())
		}
	}

	return nil
}

// GetRequiredFields returns required configuration fields
func (v *GoConfigValidator) GetRequiredFields() []string {
	return []string{"name", "module"}
}

// GetOptionalFields returns optional configuration fields
func (v *GoConfigValidator) GetOptionalFields() []string {
	return []string{"framework", "port"}
}

// GetFieldDescription returns description for a field
func (v *GoConfigValidator) GetFieldDescription(field string) string {
	descriptions := map[string]string{
		"name":      "The name of the Go backend project",
		"module":    "The Go module path (e.g., github.com/user/project)",
		"framework": "The web framework to use (gin, echo, or fiber)",
		"port":      "The port number for the server (1-65535)",
	}

	if desc, exists := descriptions[field]; exists {
		return desc
	}
	return ""
}

// ValidateModule validates a Go module path
func (v *GoConfigValidator) ValidateModule(value interface{}) error {
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

	// Module path should not start or end with /
	if strings.HasPrefix(module, "/") || strings.HasSuffix(module, "/") {
		return fmt.Errorf("module path cannot start or end with /")
	}

	return nil
}

// ValidatePort validates a port number
func (v *GoConfigValidator) ValidatePort(value interface{}) error {
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

	return nil
}

// ValidateFramework validates the web framework choice
func (v *GoConfigValidator) ValidateFramework(value interface{}) error {
	framework, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	supportedFrameworks := []string{"gin", "echo", "fiber"}
	framework = strings.ToLower(framework)

	for _, supported := range supportedFrameworks {
		if framework == supported {
			return nil
		}
	}

	return fmt.Errorf("must be one of: %s", strings.Join(supportedFrameworks, ", "))
}
