package config

import (
	"fmt"
)

// NextJSConfigValidator validates NextJS component configurations
type NextJSConfigValidator struct{}

// NewNextJSConfigValidator creates a new NextJS config validator
func NewNextJSConfigValidator() *NextJSConfigValidator {
	return &NextJSConfigValidator{}
}

// Validate validates the NextJS configuration
func (v *NextJSConfigValidator) Validate(config map[string]interface{}) error {
	// Validate name
	if name, exists := config["name"]; exists {
		if err := validateProjectName(name); err != nil {
			return NewFieldError("name", err.Error())
		}
	}

	// Validate typescript (must be boolean)
	if typescript, exists := config["typescript"]; exists {
		if _, ok := typescript.(bool); !ok {
			return NewFieldError("typescript", "must be a boolean value (true or false)")
		}
	}

	// Validate tailwind (must be boolean)
	if tailwind, exists := config["tailwind"]; exists {
		if _, ok := tailwind.(bool); !ok {
			return NewFieldError("tailwind", "must be a boolean value (true or false)")
		}
	}

	// Validate app_router (must be boolean)
	if appRouter, exists := config["app_router"]; exists {
		if _, ok := appRouter.(bool); !ok {
			return NewFieldError("app_router", "must be a boolean value (true or false)")
		}
	}

	// Validate eslint (must be boolean)
	if eslint, exists := config["eslint"]; exists {
		if _, ok := eslint.(bool); !ok {
			return NewFieldError("eslint", "must be a boolean value (true or false)")
		}
	}

	return nil
}

// GetRequiredFields returns required configuration fields
func (v *NextJSConfigValidator) GetRequiredFields() []string {
	return []string{"name"}
}

// GetOptionalFields returns optional configuration fields
func (v *NextJSConfigValidator) GetOptionalFields() []string {
	return []string{"typescript", "tailwind", "app_router", "eslint"}
}

// GetFieldDescription returns description for a field
func (v *NextJSConfigValidator) GetFieldDescription(field string) string {
	descriptions := map[string]string{
		"name":       "The name of the Next.js project",
		"typescript": "Enable TypeScript support (boolean)",
		"tailwind":   "Enable Tailwind CSS support (boolean)",
		"app_router": "Use Next.js App Router instead of Pages Router (boolean)",
		"eslint":     "Enable ESLint configuration (boolean)",
	}

	if desc, exists := descriptions[field]; exists {
		return desc
	}
	return ""
}

// ValidateTypescript validates the typescript field
func (v *NextJSConfigValidator) ValidateTypescript(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("typescript must be a boolean value")
	}
	return nil
}

// ValidateTailwind validates the tailwind field
func (v *NextJSConfigValidator) ValidateTailwind(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("tailwind must be a boolean value")
	}
	return nil
}

// ValidateAppRouter validates the app_router field
func (v *NextJSConfigValidator) ValidateAppRouter(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("app_router must be a boolean value")
	}
	return nil
}

// ValidateESLint validates the eslint field
func (v *NextJSConfigValidator) ValidateESLint(value interface{}) error {
	if _, ok := value.(bool); !ok {
		return fmt.Errorf("eslint must be a boolean value")
	}
	return nil
}
