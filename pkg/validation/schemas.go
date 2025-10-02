package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Pre-compiled regular expressions for schema validation
var (
	invalidCharsRegex = regexp.MustCompile(`[^a-z0-9\-._~]`)
	envKeyRegex       = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
)

// SchemaManager manages validation schemas and rules for configuration files
type SchemaManager struct {
	schemas         map[string]*interfaces.ConfigSchema
	validationRules map[string][]ValidationSchemaRule
}

// ValidationSchemaRule defines a validation rule for schema validation
type ValidationSchemaRule struct {
	ID          string
	Name        string
	Description string
	Category    string
	Severity    string
	Enabled     bool
	Pattern     string
	FileTypes   []string
	Validator   func(interface{}) error
}

// NewSchemaManager creates a new schema manager with default schemas
func NewSchemaManager() *SchemaManager {
	sm := &SchemaManager{
		schemas:         make(map[string]*interfaces.ConfigSchema),
		validationRules: make(map[string][]ValidationSchemaRule),
	}

	sm.initializeDefaultSchemas()
	sm.initializeValidationRules()

	return sm
}

// GetSchema returns a schema by name
func (sm *SchemaManager) GetSchema(name string) (*interfaces.ConfigSchema, bool) {
	schema, exists := sm.schemas[name]
	return schema, exists
}

// AddSchema adds a new schema to the manager
func (sm *SchemaManager) AddSchema(name string, schema *interfaces.ConfigSchema) {
	sm.schemas[name] = schema
}

// RemoveSchema removes a schema from the manager
func (sm *SchemaManager) RemoveSchema(name string) {
	delete(sm.schemas, name)
}

// ListSchemas returns all available schema names
func (sm *SchemaManager) ListSchemas() []string {
	names := make([]string, 0, len(sm.schemas))
	for name := range sm.schemas {
		names = append(names, name)
	}
	return names
}

// ValidateAgainstSchema validates data against a configuration schema
func (sm *SchemaManager) ValidateAgainstSchema(data interface{}, schema *interfaces.ConfigSchema, result *interfaces.ConfigValidationResult) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("data must be an object")
	}

	// Check required properties
	for _, required := range schema.Required {
		result.Summary.TotalProperties++
		if _, exists := dataMap[required]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    required,
				Value:    "",
				Type:     "missing_required",
				Message:  fmt.Sprintf("Required property '%s' is missing", required),
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.required_property",
			})
			result.Summary.ErrorCount++
			result.Summary.MissingRequired++
		} else {
			result.Summary.ValidProperties++
		}
	}

	// Validate each property
	for key, value := range dataMap {
		if propSchema, exists := schema.Properties[key]; exists {
			if err := sm.validatePropertyAgainstSchema(key, value, propSchema, result); err != nil {
				return fmt.Errorf("validation failed for property '%s': %w", key, err)
			}
		}
	}

	return nil
}

// validatePropertyAgainstSchema validates a single property against its schema
func (sm *SchemaManager) validatePropertyAgainstSchema(key string, value interface{}, schema interfaces.PropertySchema, result *interfaces.ConfigValidationResult) error {
	// Type validation
	switch schema.Type {
	case "string":
		if strValue, ok := value.(string); ok {
			// Length validation
			if schema.MinLength != nil && len(strValue) < *schema.MinLength {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    key,
					Value:    strValue,
					Type:     "validation_error",
					Message:  fmt.Sprintf("String too short, minimum length is %d", *schema.MinLength),
					Severity: interfaces.ValidationSeverityError,
					Rule:     "schema.min_length",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}

			if schema.MaxLength != nil && len(strValue) > *schema.MaxLength {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    key,
					Value:    strValue,
					Type:     "validation_error",
					Message:  fmt.Sprintf("String too long, maximum length is %d", *schema.MaxLength),
					Severity: interfaces.ValidationSeverityError,
					Rule:     "schema.max_length",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}

			// Pattern validation
			if schema.Pattern != "" {
				regex, err := regexp.Compile(schema.Pattern)
				if err != nil {
					return fmt.Errorf("invalid pattern in schema: %w", err)
				}
				if !regex.MatchString(strValue) {
					result.Errors = append(result.Errors, interfaces.ConfigValidationError{
						Field:    key,
						Value:    strValue,
						Type:     "validation_error",
						Message:  fmt.Sprintf("String does not match pattern %s", schema.Pattern),
						Severity: interfaces.ValidationSeverityError,
						Rule:     "schema.pattern",
					})
					result.Summary.ErrorCount++
					result.Valid = false
				}
			}

			// Enum validation
			if len(schema.Enum) > 0 {
				valid := false
				for _, enumValue := range schema.Enum {
					if strValue == enumValue {
						valid = true
						break
					}
				}
				if !valid {
					result.Errors = append(result.Errors, interfaces.ConfigValidationError{
						Field:    key,
						Value:    strValue,
						Type:     "validation_error",
						Message:  fmt.Sprintf("Value must be one of: %v", schema.Enum),
						Severity: interfaces.ValidationSeverityError,
						Rule:     "schema.enum",
					})
					result.Summary.ErrorCount++
					result.Valid = false
				}
			}
		} else {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    key,
				Value:    fmt.Sprintf("%v", value),
				Type:     "type_error",
				Message:  "Expected string type",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.type",
			})
			result.Summary.ErrorCount++
			result.Valid = false
		}
	case "number":
		if numValue, ok := value.(float64); ok {
			// Range validation
			if schema.Minimum != nil && numValue < *schema.Minimum {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    key,
					Value:    fmt.Sprintf("%v", numValue),
					Type:     "validation_error",
					Message:  fmt.Sprintf("Number too small, minimum is %v", *schema.Minimum),
					Severity: interfaces.ValidationSeverityError,
					Rule:     "schema.minimum",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}

			if schema.Maximum != nil && numValue > *schema.Maximum {
				result.Errors = append(result.Errors, interfaces.ConfigValidationError{
					Field:    key,
					Value:    fmt.Sprintf("%v", numValue),
					Type:     "validation_error",
					Message:  fmt.Sprintf("Number too large, maximum is %v", *schema.Maximum),
					Severity: interfaces.ValidationSeverityError,
					Rule:     "schema.maximum",
				})
				result.Summary.ErrorCount++
				result.Valid = false
			}
		} else {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    key,
				Value:    fmt.Sprintf("%v", value),
				Type:     "type_error",
				Message:  "Expected number type",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.type",
			})
			result.Summary.ErrorCount++
			result.Valid = false
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    key,
				Value:    fmt.Sprintf("%v", value),
				Type:     "type_error",
				Message:  "Expected boolean type",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.type",
			})
			result.Summary.ErrorCount++
			result.Valid = false
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    key,
				Value:    fmt.Sprintf("%v", value),
				Type:     "type_error",
				Message:  "Expected array type",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.type",
			})
			result.Summary.ErrorCount++
			result.Valid = false
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    key,
				Value:    fmt.Sprintf("%v", value),
				Type:     "type_error",
				Message:  "Expected object type",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "schema.type",
			})
			result.Summary.ErrorCount++
			result.Valid = false
		}
	}

	return nil
}

// GetValidationRules returns validation rules for a specific file type
func (sm *SchemaManager) GetValidationRules(fileType string) []ValidationSchemaRule {
	return sm.validationRules[fileType]
}

// AddValidationRule adds a new validation rule
func (sm *SchemaManager) AddValidationRule(fileType string, rule ValidationSchemaRule) {
	if sm.validationRules[fileType] == nil {
		sm.validationRules[fileType] = []ValidationSchemaRule{}
	}
	sm.validationRules[fileType] = append(sm.validationRules[fileType], rule)
}

// RemoveValidationRule removes a validation rule by ID
func (sm *SchemaManager) RemoveValidationRule(fileType, ruleID string) {
	rules := sm.validationRules[fileType]
	for i, rule := range rules {
		if rule.ID == ruleID {
			sm.validationRules[fileType] = append(rules[:i], rules[i+1:]...)
			break
		}
	}
}

// ValidatePackageName validates NPM package name format
func (sm *SchemaManager) ValidatePackageName(name string) error {
	// NPM package name rules
	if len(name) > 214 {
		return fmt.Errorf("name too long")
	}

	if strings.ToLower(name) != name {
		return fmt.Errorf("name must be lowercase")
	}

	// Check for invalid characters using pre-compiled regex
	if invalidCharsRegex.MatchString(name) {
		return fmt.Errorf("name contains invalid characters")
	}

	return nil
}

// ValidateEnvKey validates environment variable key format
func (sm *SchemaManager) ValidateEnvKey(key string) error {
	// Environment variable names should be uppercase with underscores
	if !envKeyRegex.MatchString(key) {
		return fmt.Errorf("environment variable names should be uppercase with underscores")
	}
	return nil
}

// IsPotentialSecret checks if a key-value pair might contain a secret
func (sm *SchemaManager) IsPotentialSecret(key, value string) bool {
	secretKeywords := []string{"password", "secret", "key", "token", "api", "auth"}
	keyLower := strings.ToLower(key)

	for _, keyword := range secretKeywords {
		if strings.Contains(keyLower, keyword) && len(value) > 10 {
			return true
		}
	}

	return false
}

// initializeDefaultSchemas initializes default configuration schemas
func (sm *SchemaManager) initializeDefaultSchemas() {
	// Package.json schema
	sm.schemas["package.json"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Package.json Schema",
		Description: "Schema for Node.js package.json files",
		Required:    []string{"name", "version"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Package name",
				Pattern:     `^[a-z0-9\-._~]+$`,
				MaxLength:   &[]int{214}[0],
			},
			"version": {
				Type:        "string",
				Description: "Package version",
				Pattern:     `^\d+\.\d+\.\d+`,
			},
			"description": {
				Type:        "string",
				Description: "Package description",
				MaxLength:   &[]int{500}[0],
			},
			"main": {
				Type:        "string",
				Description: "Entry point file",
			},
			"scripts": {
				Type:        "object",
				Description: "NPM scripts",
			},
			"dependencies": {
				Type:        "object",
				Description: "Production dependencies",
			},
			"devDependencies": {
				Type:        "object",
				Description: "Development dependencies",
			},
			"keywords": {
				Type:        "array",
				Description: "Package keywords",
			},
			"author": {
				Type:        "string",
				Description: "Package author",
			},
			"license": {
				Type:        "string",
				Description: "Package license",
			},
		},
	}

	// TypeScript config schema
	sm.schemas["tsconfig.json"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "TypeScript Configuration Schema",
		Description: "Schema for TypeScript configuration files",
		Required:    []string{},
		Properties: map[string]interfaces.PropertySchema{
			"compilerOptions": {
				Type:        "object",
				Description: "TypeScript compiler options",
				Required:    true,
			},
			"include": {
				Type:        "array",
				Description: "Files to include in compilation",
			},
			"exclude": {
				Type:        "array",
				Description: "Files to exclude from compilation",
			},
			"extends": {
				Type:        "string",
				Description: "Base configuration to extend",
			},
		},
	}

	// ESLint config schema
	sm.schemas[".eslintrc.json"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "ESLint Configuration Schema",
		Description: "Schema for ESLint configuration files",
		Required:    []string{},
		Properties: map[string]interfaces.PropertySchema{
			"extends": {
				Type:        "array",
				Description: "Base configurations to extend",
			},
			"rules": {
				Type:        "object",
				Description: "ESLint rules configuration",
			},
			"env": {
				Type:        "object",
				Description: "Environment settings",
			},
			"parserOptions": {
				Type:        "object",
				Description: "Parser options",
			},
			"plugins": {
				Type:        "array",
				Description: "ESLint plugins",
			},
		},
	}

	// Docker Compose schema
	sm.schemas["docker-compose.yml"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Docker Compose Schema",
		Description: "Schema for Docker Compose files",
		Required:    []string{"services"},
		Properties: map[string]interfaces.PropertySchema{
			"version": {
				Type:        "string",
				Description: "Docker Compose file format version",
				Enum:        []string{"3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "3.8", "3.9"},
			},
			"services": {
				Type:        "object",
				Description: "Service definitions",
				Required:    true,
			},
			"networks": {
				Type:        "object",
				Description: "Network definitions",
			},
			"volumes": {
				Type:        "object",
				Description: "Volume definitions",
			},
		},
	}

	// GitHub Actions workflow schema
	sm.schemas["workflow.yml"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "GitHub Actions Workflow Schema",
		Description: "Schema for GitHub Actions workflow files",
		Required:    []string{"on", "jobs"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Workflow name",
			},
			"on": {
				Type:        "object",
				Description: "Workflow triggers",
				Required:    true,
			},
			"jobs": {
				Type:        "object",
				Description: "Job definitions",
				Required:    true,
			},
			"env": {
				Type:        "object",
				Description: "Environment variables",
			},
		},
	}
}

// initializeValidationRules initializes validation rules for different file types
func (sm *SchemaManager) initializeValidationRules() {
	// Package.json validation rules
	sm.validationRules["package.json"] = []ValidationSchemaRule{
		{
			ID:          "package_json_name_format",
			Name:        "Package Name Format",
			Description: "Validates NPM package name format",
			Category:    "format",
			Severity:    "error",
			Enabled:     true,
			FileTypes:   []string{"package.json"},
		},
		{
			ID:          "package_json_version_format",
			Name:        "Version Format",
			Description: "Validates semantic version format",
			Category:    "format",
			Severity:    "error",
			Enabled:     true,
			Pattern:     `^\d+\.\d+\.\d+`,
			FileTypes:   []string{"package.json"},
		},
		{
			ID:          "package_json_license",
			Name:        "License Field",
			Description: "Recommends including license field",
			Category:    "best_practice",
			Severity:    "warning",
			Enabled:     true,
			FileTypes:   []string{"package.json"},
		},
	}

	// TypeScript config validation rules
	sm.validationRules["tsconfig.json"] = []ValidationSchemaRule{
		{
			ID:          "tsconfig_strict_mode",
			Name:        "Strict Mode",
			Description: "Recommends enabling TypeScript strict mode",
			Category:    "best_practice",
			Severity:    "warning",
			Enabled:     true,
			FileTypes:   []string{"tsconfig.json"},
		},
		{
			ID:          "tsconfig_compiler_options",
			Name:        "Compiler Options",
			Description: "Validates presence of compiler options",
			Category:    "structure",
			Severity:    "error",
			Enabled:     true,
			FileTypes:   []string{"tsconfig.json"},
		},
	}

	// Docker Compose validation rules
	sm.validationRules["docker-compose.yml"] = []ValidationSchemaRule{
		{
			ID:          "docker_compose_version",
			Name:        "Version Deprecation",
			Description: "Warns about deprecated Docker Compose versions",
			Category:    "deprecation",
			Severity:    "warning",
			Enabled:     true,
			FileTypes:   []string{"docker-compose.yml", "docker-compose.yaml"},
		},
		{
			ID:          "docker_compose_privileged",
			Name:        "Privileged Mode",
			Description: "Warns about security risks of privileged mode",
			Category:    "security",
			Severity:    "warning",
			Enabled:     true,
			FileTypes:   []string{"docker-compose.yml", "docker-compose.yaml"},
		},
	}

	// Environment file validation rules
	sm.validationRules[".env"] = []ValidationSchemaRule{
		{
			ID:          "env_key_format",
			Name:        "Environment Key Format",
			Description: "Validates environment variable naming conventions",
			Category:    "format",
			Severity:    "warning",
			Enabled:     true,
			Pattern:     `^[A-Z][A-Z0-9_]*$`,
			FileTypes:   []string{".env"},
		},
		{
			ID:          "env_secrets_detection",
			Name:        "Secrets Detection",
			Description: "Detects potential secrets in environment files",
			Category:    "security",
			Severity:    "warning",
			Enabled:     true,
			FileTypes:   []string{".env"},
		},
	}
}
