package formats

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// JSONValidator provides specialized JSON configuration file validation
type JSONValidator struct {
	schemas map[string]*interfaces.ConfigSchema
}

// NewJSONValidator creates a new JSON configuration validator
func NewJSONValidator() *JSONValidator {
	validator := &JSONValidator{
		schemas: make(map[string]*interfaces.ConfigSchema),
	}
	validator.initializeJSONSchemas()
	return validator
}

// ValidateJSONFile validates a JSON configuration file
func (jv *JSONValidator) ValidateJSONFile(filePath string) (*interfaces.ConfigValidationResult, error) {
	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	// Read and parse JSON
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "file",
			Value:    filePath,
			Type:     "read_error",
			Message:  fmt.Sprintf("Failed to read file: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.file_access",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "syntax",
			Value:    string(content),
			Type:     "syntax_error",
			Message:  fmt.Sprintf("Invalid JSON syntax: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.json.syntax",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	// Validate against schema if available
	fileName := filepath.Base(filePath)
	if schema, exists := jv.schemas[fileName]; exists {
		if err := jv.validateAgainstSchema(data, schema, result); err != nil {
			return nil, fmt.Errorf("schema validation failed: %w", err)
		}
	}

	// Perform specific validations based on file name
	switch fileName {
	case "package.json":
		jv.validatePackageJSON(data, result)
	case "tsconfig.json":
		jv.validateTSConfig(data, result)
	case ".eslintrc.json":
		jv.validateESLintConfig(data, result)
	}

	return result, nil
}

// validatePackageJSON validates package.json specific structure
func (jv *JSONValidator) validatePackageJSON(data interface{}, result *interfaces.ConfigValidationResult) {
	pkg, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check required fields
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		result.Summary.TotalProperties++
		if _, exists := pkg[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    field,
				Value:    "",
				Type:     "missing_required",
				Message:  fmt.Sprintf("Required field '%s' is missing", field),
				Severity: interfaces.ValidationSeverityError,
				Rule:     "package_json.required_fields",
			})
			result.Summary.ErrorCount++
			result.Summary.MissingRequired++
		} else {
			result.Summary.ValidProperties++
		}
	}

	// Validate name format
	if name, exists := pkg["name"]; exists {
		if nameStr, ok := name.(string); ok {
			if err := jv.validatePackageName(nameStr); err != nil {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      "name",
					Value:      nameStr,
					Type:       "format_warning",
					Message:    fmt.Sprintf("Package name format issue: %v", err),
					Suggestion: "Use lowercase letters, numbers, hyphens, and dots only",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "package_json.name_format",
				})
				result.Summary.WarningCount++
			}
		}
	}

	// Check for security vulnerabilities in dependencies
	if deps, exists := pkg["dependencies"]; exists {
		if depsMap, ok := deps.(map[string]interface{}); ok {
			jv.validateDependencies(depsMap, "dependencies", result)
		}
	}

	if devDeps, exists := pkg["devDependencies"]; exists {
		if devDepsMap, ok := devDeps.(map[string]interface{}); ok {
			jv.validateDependencies(devDepsMap, "devDependencies", result)
		}
	}
}

// validateTSConfig validates TypeScript configuration
func (jv *JSONValidator) validateTSConfig(data interface{}, result *interfaces.ConfigValidationResult) {
	config, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for compilerOptions
	result.Summary.TotalProperties++
	if compilerOptions, exists := config["compilerOptions"]; exists {
		result.Summary.ValidProperties++

		if options, ok := compilerOptions.(map[string]interface{}); ok {
			// Check for strict mode
			if strict, exists := options["strict"]; exists {
				if strictBool, ok := strict.(bool); ok && !strictBool {
					result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
						Field:      "compilerOptions.strict",
						Value:      "false",
						Type:       "best_practice",
						Message:    "Consider enabling strict mode for better type safety",
						Suggestion: "Set 'strict': true in compilerOptions",
						Severity:   interfaces.ValidationSeverityWarning,
						Rule:       "tsconfig.strict_mode",
					})
					result.Summary.WarningCount++
				}
			}

			// Check for target version
			if target, exists := options["target"]; exists {
				if targetStr, ok := target.(string); ok {
					if strings.ToLower(targetStr) == "es5" {
						result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
							Field:      "compilerOptions.target",
							Value:      targetStr,
							Type:       "outdated",
							Message:    "ES5 target is outdated",
							Suggestion: "Consider using ES2020 or later for better performance",
							Severity:   interfaces.ValidationSeverityWarning,
							Rule:       "tsconfig.target_version",
						})
						result.Summary.WarningCount++
					}
				}
			}
		}
	} else {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "compilerOptions",
			Value:      "",
			Type:       "missing_recommended",
			Message:    "compilerOptions section is recommended",
			Suggestion: "Add compilerOptions to configure TypeScript compilation",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "tsconfig.compiler_options",
		})
		result.Summary.WarningCount++
	}
}

// validateESLintConfig validates ESLint configuration
func (jv *JSONValidator) validateESLintConfig(data interface{}, result *interfaces.ConfigValidationResult) {
	config, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for extends
	result.Summary.TotalProperties++
	if _, exists := config["extends"]; exists {
		result.Summary.ValidProperties++
	} else {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "extends",
			Value:      "",
			Type:       "missing_recommended",
			Message:    "Consider extending from a base configuration",
			Suggestion: "Add 'extends' to inherit from standard configurations",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "eslint.extends",
		})
		result.Summary.WarningCount++
	}

	// Check for rules
	if rules, exists := config["rules"]; exists {
		if rulesMap, ok := rules.(map[string]interface{}); ok {
			jv.validateESLintRules(rulesMap, result)
		}
	}
}

// validateDependencies validates package dependencies for security issues
func (jv *JSONValidator) validateDependencies(deps map[string]interface{}, section string, result *interfaces.ConfigValidationResult) {
	knownVulnerableDeps := map[string][]string{
		"lodash": {"<4.17.19"},
		"axios":  {"<0.21.1"},
		"yargs":  {"<15.4.1"},
	}

	for depName, version := range deps {
		if versionStr, ok := version.(string); ok {
			if vulnerableVersions, exists := knownVulnerableDeps[depName]; exists {
				for _, vulnVersion := range vulnerableVersions {
					if strings.Contains(versionStr, strings.TrimPrefix(vulnVersion, "<")) {
						result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
							Field:      fmt.Sprintf("%s.%s", section, depName),
							Value:      versionStr,
							Type:       "security",
							Message:    fmt.Sprintf("Potentially vulnerable version of %s", depName),
							Suggestion: fmt.Sprintf("Update to a version >= %s", strings.TrimPrefix(vulnVersion, "<")),
							Severity:   interfaces.ValidationSeverityWarning,
							Rule:       "package_json.vulnerable_dependency",
						})
						result.Summary.WarningCount++
					}
				}
			}
		}
	}
}

// validateESLintRules validates ESLint rules configuration
func (jv *JSONValidator) validateESLintRules(rules map[string]interface{}, result *interfaces.ConfigValidationResult) {
	recommendedRules := map[string]string{
		"no-console":     "warn",
		"no-debugger":    "error",
		"no-unused-vars": "error",
	}

	for ruleName, recommendedValue := range recommendedRules {
		if _, exists := rules[ruleName]; !exists {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("rules.%s", ruleName),
				Value:      "",
				Type:       "missing_recommended",
				Message:    fmt.Sprintf("Consider adding rule '%s'", ruleName),
				Suggestion: fmt.Sprintf("Add '%s': '%s' to rules", ruleName, recommendedValue),
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "eslint.recommended_rules",
			})
			result.Summary.WarningCount++
		}
	}
}

// Helper methods

// validateAgainstSchema validates data against a configuration schema
func (jv *JSONValidator) validateAgainstSchema(data interface{}, schema *interfaces.ConfigSchema, result *interfaces.ConfigValidationResult) error {
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
			if err := jv.validatePropertyAgainstSchema(key, value, propSchema, result); err != nil {
				return fmt.Errorf("validation failed for property '%s': %w", key, err)
			}
		}
	}

	return nil
}

// validatePropertyAgainstSchema validates a single property against its schema
func (jv *JSONValidator) validatePropertyAgainstSchema(key string, value interface{}, schema interfaces.PropertySchema, result *interfaces.ConfigValidationResult) error {
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
	}

	return nil
}

// validatePackageName validates NPM package name format
func (jv *JSONValidator) validatePackageName(name string) error {
	// NPM package name rules
	if len(name) > 214 {
		return fmt.Errorf("name too long")
	}

	if strings.ToLower(name) != name {
		return fmt.Errorf("name must be lowercase")
	}

	// Check for invalid characters (simplified)
	if strings.ContainsAny(name, " !@#$%^&*()+=[]{}|\\:;\"'<>?,/") {
		return fmt.Errorf("name contains invalid characters")
	}

	return nil
}

// initializeJSONSchemas initializes default JSON configuration schemas
func (jv *JSONValidator) initializeJSONSchemas() {
	// Package.json schema
	jv.schemas["package.json"] = &interfaces.ConfigSchema{
		Version:     "1.0",
		Title:       "Package.json Schema",
		Description: "Schema for Node.js package.json files",
		Required:    []string{"name", "version"},
		Properties: map[string]interfaces.PropertySchema{
			"name": {
				Type:        "string",
				Description: "Package name",
				MaxLength:   &[]int{214}[0],
			},
			"version": {
				Type:        "string",
				Description: "Package version",
			},
			"description": {
				Type:        "string",
				Description: "Package description",
				MaxLength:   &[]int{500}[0],
			},
		},
	}
}
