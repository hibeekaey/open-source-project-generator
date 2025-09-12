package standards

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult represents the result of a template validation
type ValidationResult struct {
	IsValid      bool                   `json:"is_valid"`
	Errors       []ValidationError      `json:"errors"`
	Warnings     []ValidationWarning    `json:"warnings"`
	TemplateName string                 `json:"template_name"`
	Suggestions  []ValidationSuggestion `json:"suggestions"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Field       string `json:"field"`
	Expected    string `json:"expected"`
	Actual      string `json:"actual"`
	Description string `json:"description"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Field       string `json:"field"`
	Description string `json:"description"`
}

// ValidationSuggestion represents a suggestion for improvement
type ValidationSuggestion struct {
	Type        string `json:"type"`
	File        string `json:"file"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// TemplateValidator validates frontend templates against standards
type TemplateValidator struct {
	standards *FrontendStandards
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{
		standards: GetFrontendStandards(),
	}
}

// ValidateTemplate validates a frontend template against standards
func (v *TemplateValidator) ValidateTemplate(templatePath, templateType string) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:      true,
		Errors:       []ValidationError{},
		Warnings:     []ValidationWarning{},
		Suggestions:  []ValidationSuggestion{},
		TemplateName: templateType,
	}

	// Validate package.json
	if err := v.validatePackageJSON(templatePath, templateType, result); err != nil {
		return nil, fmt.Errorf("failed to validate package.json: %w", err)
	}

	// Validate tsconfig.json
	if err := v.validateTSConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate tsconfig.json: %w", err)
	}

	// Validate .eslintrc.json
	if err := v.validateESLintConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate .eslintrc.json: %w", err)
	}

	// Validate .prettierrc
	if err := v.validatePrettierConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate .prettierrc: %w", err)
	}

	// Validate vercel.json
	if err := v.validateVercelConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate vercel.json: %w", err)
	}

	// Validate tailwind.config.js
	if err := v.validateTailwindConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate tailwind.config.js: %w", err)
	}

	// Validate next.config.js
	if err := v.validateNextConfig(templatePath, result); err != nil {
		return nil, fmt.Errorf("failed to validate next.config.js: %w", err)
	}

	// Set overall validation status
	result.IsValid = len(result.Errors) == 0

	return result, nil
}

// validatePackageJSON validates package.json against standards
func (v *TemplateValidator) validatePackageJSON(templatePath, templateType string, result *ValidationResult) error {
	packagePath := filepath.Join(templatePath, "package.json.tmpl")

	if !fileExists(packagePath) {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_file",
			File:        "package.json.tmpl",
			Description: "package.json.tmpl file is missing",
		})
		return nil
	}

	content, err := os.ReadFile(packagePath)
	if err != nil {
		return fmt.Errorf("failed to read package.json.tmpl: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "invalid_json",
			File:        "package.json.tmpl",
			Description: "Invalid JSON format",
		})
		return nil
	}

	// Validate scripts
	v.validatePackageScripts(pkg, templateType, result)

	// Validate engines
	v.validatePackageEngines(pkg, result)

	// Validate dependencies
	v.validatePackageDependencies(pkg, templateType, result)

	// Validate dev dependencies
	v.validatePackageDevDependencies(pkg, result)

	return nil
}

// validatePackageScripts validates package.json scripts
func (v *TemplateValidator) validatePackageScripts(pkg map[string]interface{}, templateType string, result *ValidationResult) {
	scripts, ok := pkg["scripts"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_field",
			File:        "package.json.tmpl",
			Field:       "scripts",
			Description: "scripts field is missing or invalid",
		})
		return
	}

	expectedScripts := v.standards.PackageJSON.Scripts

	// Adjust expected scripts for template-specific ports
	if port, exists := v.standards.PackageJSON.Ports[templateType]; exists && templateType != "nextjs-app" {
		expectedScripts = make(map[string]string)
		for k, v := range v.standards.PackageJSON.Scripts {
			expectedScripts[k] = v
		}
		expectedScripts["dev"] = fmt.Sprintf("next dev -p %d", port)
		expectedScripts["start"] = fmt.Sprintf("next start -p %d", port)
	}

	for scriptName, expectedCommand := range expectedScripts {
		if actualCommand, exists := scripts[scriptName]; exists {
			if actualCommand != expectedCommand {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "script_mismatch",
					File:        "package.json.tmpl",
					Field:       fmt.Sprintf("scripts.%s", scriptName),
					Description: fmt.Sprintf("Script '%s' should be '%s' but is '%s'", scriptName, expectedCommand, actualCommand),
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Type:        "missing_script",
				File:        "package.json.tmpl",
				Field:       fmt.Sprintf("scripts.%s", scriptName),
				Expected:    expectedCommand,
				Description: fmt.Sprintf("Required script '%s' is missing", scriptName),
			})
		}
	}
}

// validatePackageEngines validates package.json engines
func (v *TemplateValidator) validatePackageEngines(pkg map[string]interface{}, result *ValidationResult) {
	engines, ok := pkg["engines"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_field",
			File:        "package.json.tmpl",
			Field:       "engines",
			Description: "engines field is missing or invalid",
		})
		return
	}

	for engineName, expectedVersion := range v.standards.PackageJSON.Engines {
		if actualVersion, exists := engines[engineName]; exists {
			if actualVersion != expectedVersion {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "engine_version_mismatch",
					File:        "package.json.tmpl",
					Field:       fmt.Sprintf("engines.%s", engineName),
					Description: fmt.Sprintf("Engine '%s' should be '%s' but is '%s'", engineName, expectedVersion, actualVersion),
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Type:        "missing_engine",
				File:        "package.json.tmpl",
				Field:       fmt.Sprintf("engines.%s", engineName),
				Expected:    expectedVersion,
				Description: fmt.Sprintf("Required engine '%s' is missing", engineName),
			})
		}
	}
}

// validatePackageDependencies validates package.json dependencies
func (v *TemplateValidator) validatePackageDependencies(pkg map[string]interface{}, templateType string, result *ValidationResult) {
	dependencies, ok := pkg["dependencies"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_field",
			File:        "package.json.tmpl",
			Field:       "dependencies",
			Description: "dependencies field is missing or invalid",
		})
		return
	}

	// Check base dependencies
	for depName, expectedVersion := range v.standards.PackageJSON.Dependencies {
		if actualVersion, exists := dependencies[depName]; exists {
			if actualVersion != expectedVersion {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "dependency_version_mismatch",
					File:        "package.json.tmpl",
					Field:       fmt.Sprintf("dependencies.%s", depName),
					Description: fmt.Sprintf("Dependency '%s' should be '%s' but is '%s'", depName, expectedVersion, actualVersion),
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Type:        "missing_dependency",
				File:        "package.json.tmpl",
				Field:       fmt.Sprintf("dependencies.%s", depName),
				Expected:    expectedVersion,
				Description: fmt.Sprintf("Required dependency '%s' is missing", depName),
			})
		}
	}

	// Check template-specific dependencies
	templateDeps := GetTemplateSpecificDependencies(templateType)
	for depName, expectedVersion := range templateDeps {
		if actualVersion, exists := dependencies[depName]; exists {
			if actualVersion != expectedVersion {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "template_dependency_version_mismatch",
					File:        "package.json.tmpl",
					Field:       fmt.Sprintf("dependencies.%s", depName),
					Description: fmt.Sprintf("Template-specific dependency '%s' should be '%s' but is '%s'", depName, expectedVersion, actualVersion),
				})
			}
		} else {
			result.Suggestions = append(result.Suggestions, ValidationSuggestion{
				Type:        "missing_template_dependency",
				File:        "package.json.tmpl",
				Description: fmt.Sprintf("Consider adding template-specific dependency '%s': '%s'", depName, expectedVersion),
				Action:      "add_dependency",
			})
		}
	}
}

// validatePackageDevDependencies validates package.json devDependencies
func (v *TemplateValidator) validatePackageDevDependencies(pkg map[string]interface{}, result *ValidationResult) {
	devDependencies, ok := pkg["devDependencies"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_field",
			File:        "package.json.tmpl",
			Field:       "devDependencies",
			Description: "devDependencies field is missing or invalid",
		})
		return
	}

	for depName, expectedVersion := range v.standards.PackageJSON.DevDeps {
		if actualVersion, exists := devDependencies[depName]; exists {
			if actualVersion != expectedVersion {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "dev_dependency_version_mismatch",
					File:        "package.json.tmpl",
					Field:       fmt.Sprintf("devDependencies.%s", depName),
					Description: fmt.Sprintf("Dev dependency '%s' should be '%s' but is '%s'", depName, expectedVersion, actualVersion),
				})
			}
		} else {
			result.Errors = append(result.Errors, ValidationError{
				Type:        "missing_dev_dependency",
				File:        "package.json.tmpl",
				Field:       fmt.Sprintf("devDependencies.%s", depName),
				Expected:    expectedVersion,
				Description: fmt.Sprintf("Required dev dependency '%s' is missing", depName),
			})
		}
	}
}

// validateTSConfig validates tsconfig.json against standards
func (v *TemplateValidator) validateTSConfig(templatePath string, result *ValidationResult) error {
	tsconfigPath := filepath.Join(templatePath, "tsconfig.json.tmpl")

	if !fileExists(tsconfigPath) {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_file",
			File:        "tsconfig.json.tmpl",
			Description: "tsconfig.json.tmpl file is missing",
		})
		return nil
	}

	content, err := os.ReadFile(tsconfigPath)
	if err != nil {
		return fmt.Errorf("failed to read tsconfig.json.tmpl: %w", err)
	}

	var tsconfig map[string]interface{}
	if err := json.Unmarshal(content, &tsconfig); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "invalid_json",
			File:        "tsconfig.json.tmpl",
			Description: "Invalid JSON format",
		})
		return nil
	}

	// Validate compiler options
	compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{})
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "missing_field",
			File:        "tsconfig.json.tmpl",
			Field:       "compilerOptions",
			Description: "compilerOptions field is missing or invalid",
		})
		return nil
	}

	// Check key compiler options
	expectedOptions := v.standards.TypeScript.CompilerOptions
	for optionName, expectedValue := range expectedOptions {
		if optionName == "paths" {
			continue // Handle paths separately
		}

		if actualValue, exists := compilerOptions[optionName]; exists {
			if !compareValues(actualValue, expectedValue) {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "compiler_option_mismatch",
					File:        "tsconfig.json.tmpl",
					Field:       fmt.Sprintf("compilerOptions.%s", optionName),
					Description: fmt.Sprintf("Compiler option '%s' differs from standard", optionName),
				})
			}
		} else {
			result.Suggestions = append(result.Suggestions, ValidationSuggestion{
				Type:        "missing_compiler_option",
				File:        "tsconfig.json.tmpl",
				Description: fmt.Sprintf("Consider adding compiler option '%s'", optionName),
				Action:      "add_compiler_option",
			})
		}
	}

	return nil
}

// validateESLintConfig validates .eslintrc.json against standards
func (v *TemplateValidator) validateESLintConfig(templatePath string, result *ValidationResult) error {
	eslintPath := filepath.Join(templatePath, ".eslintrc.json.tmpl")

	if !fileExists(eslintPath) {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:        "missing_file",
			File:        ".eslintrc.json.tmpl",
			Description: ".eslintrc.json.tmpl file is missing",
		})
		return nil
	}

	content, err := os.ReadFile(eslintPath)
	if err != nil {
		return fmt.Errorf("failed to read .eslintrc.json.tmpl: %w", err)
	}

	var eslintConfig map[string]interface{}
	if err := json.Unmarshal(content, &eslintConfig); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "invalid_json",
			File:        ".eslintrc.json.tmpl",
			Description: "Invalid JSON format",
		})
		return nil
	}

	// Validate extends
	if extends, ok := eslintConfig["extends"].([]interface{}); ok {
		expectedExtends := v.standards.ESLint.Extends
		for _, expectedExtend := range expectedExtends {
			found := false
			for _, actualExtend := range extends {
				if actualExtend == expectedExtend {
					found = true
					break
				}
			}
			if !found {
				result.Suggestions = append(result.Suggestions, ValidationSuggestion{
					Type:        "missing_eslint_extend",
					File:        ".eslintrc.json.tmpl",
					Description: fmt.Sprintf("Consider adding ESLint extend '%s'", expectedExtend),
					Action:      "add_eslint_extend",
				})
			}
		}
	}

	return nil
}

// validatePrettierConfig validates .prettierrc against standards
func (v *TemplateValidator) validatePrettierConfig(templatePath string, result *ValidationResult) error {
	prettierPath := filepath.Join(templatePath, ".prettierrc.tmpl")

	if !fileExists(prettierPath) {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:        "missing_file",
			File:        ".prettierrc.tmpl",
			Description: ".prettierrc.tmpl file is missing",
		})
		return nil
	}

	content, err := os.ReadFile(prettierPath)
	if err != nil {
		return fmt.Errorf("failed to read .prettierrc.tmpl: %w", err)
	}

	var prettierConfig map[string]interface{}
	if err := json.Unmarshal(content, &prettierConfig); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "invalid_json",
			File:        ".prettierrc.tmpl",
			Description: "Invalid JSON format",
		})
		return nil
	}

	// Check key prettier options
	expectedConfig := map[string]interface{}{
		"semi":          v.standards.Prettier.Semi,
		"trailingComma": v.standards.Prettier.TrailingComma,
		"singleQuote":   v.standards.Prettier.SingleQuote,
		"printWidth":    v.standards.Prettier.PrintWidth,
		"tabWidth":      v.standards.Prettier.TabWidth,
		"useTabs":       v.standards.Prettier.UseTabs,
	}

	for optionName, expectedValue := range expectedConfig {
		if actualValue, exists := prettierConfig[optionName]; exists {
			if actualValue != expectedValue {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:        "prettier_option_mismatch",
					File:        ".prettierrc.tmpl",
					Field:       optionName,
					Description: fmt.Sprintf("Prettier option '%s' should be %v but is %v", optionName, expectedValue, actualValue),
				})
			}
		} else {
			result.Suggestions = append(result.Suggestions, ValidationSuggestion{
				Type:        "missing_prettier_option",
				File:        ".prettierrc.tmpl",
				Description: fmt.Sprintf("Consider adding Prettier option '%s': %v", optionName, expectedValue),
				Action:      "add_prettier_option",
			})
		}
	}

	return nil
}

// validateVercelConfig validates vercel.json against standards
func (v *TemplateValidator) validateVercelConfig(templatePath string, result *ValidationResult) error {
	vercelPath := filepath.Join(templatePath, "vercel.json.tmpl")

	if !fileExists(vercelPath) {
		result.Suggestions = append(result.Suggestions, ValidationSuggestion{
			Type:        "missing_file",
			File:        "vercel.json.tmpl",
			Description: "Consider adding vercel.json.tmpl for Vercel deployment configuration",
			Action:      "add_vercel_config",
		})
		return nil
	}

	content, err := os.ReadFile(vercelPath)
	if err != nil {
		return fmt.Errorf("failed to read vercel.json.tmpl: %w", err)
	}

	var vercelConfig map[string]interface{}
	if err := json.Unmarshal(content, &vercelConfig); err != nil {
		result.Errors = append(result.Errors, ValidationError{
			Type:        "invalid_json",
			File:        "vercel.json.tmpl",
			Description: "Invalid JSON format",
		})
		return nil
	}

	// Check framework
	if framework, exists := vercelConfig["framework"]; exists {
		if framework != v.standards.Vercel.Framework {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:        "vercel_framework_mismatch",
				File:        "vercel.json.tmpl",
				Field:       "framework",
				Description: fmt.Sprintf("Framework should be '%s' but is '%s'", v.standards.Vercel.Framework, framework),
			})
		}
	}

	// Check build command
	if buildCommand, exists := vercelConfig["buildCommand"]; exists {
		if buildCommand != v.standards.Vercel.BuildCommand {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:        "vercel_build_command_mismatch",
				File:        "vercel.json.tmpl",
				Field:       "buildCommand",
				Description: fmt.Sprintf("Build command should be '%s' but is '%s'", v.standards.Vercel.BuildCommand, buildCommand),
			})
		}
	}

	return nil
}

// validateTailwindConfig validates tailwind.config.js against standards
func (v *TemplateValidator) validateTailwindConfig(templatePath string, result *ValidationResult) error {
	tailwindPath := filepath.Join(templatePath, "tailwind.config.js.tmpl")

	if !fileExists(tailwindPath) {
		result.Suggestions = append(result.Suggestions, ValidationSuggestion{
			Type:        "missing_file",
			File:        "tailwind.config.js.tmpl",
			Description: "Consider adding tailwind.config.js.tmpl for consistent Tailwind CSS configuration",
			Action:      "add_tailwind_config",
		})
		return nil
	}

	// For now, just check if the file exists and is readable
	content, err := os.ReadFile(tailwindPath)
	if err != nil {
		return fmt.Errorf("failed to read tailwind.config.js.tmpl: %w", err)
	}

	// Basic validation - check if it contains expected content patterns
	contentStr := string(content)
	if !strings.Contains(contentStr, "tailwindcss-animate") {
		result.Suggestions = append(result.Suggestions, ValidationSuggestion{
			Type:        "missing_tailwind_plugin",
			File:        "tailwind.config.js.tmpl",
			Description: "Consider adding tailwindcss-animate plugin for consistent animations",
			Action:      "add_tailwind_plugin",
		})
	}

	return nil
}

// validateNextConfig validates next.config.js against standards
func (v *TemplateValidator) validateNextConfig(templatePath string, result *ValidationResult) error {
	nextConfigPath := filepath.Join(templatePath, "next.config.js.tmpl")

	if !fileExists(nextConfigPath) {
		result.Suggestions = append(result.Suggestions, ValidationSuggestion{
			Type:        "missing_file",
			File:        "next.config.js.tmpl",
			Description: "Consider adding next.config.js.tmpl for consistent Next.js configuration",
			Action:      "add_next_config",
		})
		return nil
	}

	// For now, just check if the file exists and is readable
	_, err := os.ReadFile(nextConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read next.config.js.tmpl: %w", err)
	}

	return nil
}

// Helper functions

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// compareValues compares two values for equality (handles different types)
func compareValues(a, b interface{}) bool {
	switch va := a.(type) {
	case string:
		if vb, ok := b.(string); ok {
			return va == vb
		}
	case bool:
		if vb, ok := b.(bool); ok {
			return va == vb
		}
	case float64:
		if vb, ok := b.(float64); ok {
			return va == vb
		}
	case []interface{}:
		if vb, ok := b.([]interface{}); ok {
			if len(va) != len(vb) {
				return false
			}
			for i, item := range va {
				if !compareValues(item, vb[i]) {
					return false
				}
			}
			return true
		}
	case map[string]interface{}:
		if vb, ok := b.(map[string]interface{}); ok {
			if len(va) != len(vb) {
				return false
			}
			for key, value := range va {
				if vbValue, exists := vb[key]; exists {
					if !compareValues(value, vbValue) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}
