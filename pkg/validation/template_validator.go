package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/open-source-template-generator/pkg/models"
)

// TemplateValidator handles template consistency validation
type TemplateValidator struct {
	standardConfigs map[string]interface{}
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{
		standardConfigs: make(map[string]interface{}),
	}
}

// ValidateTemplateConsistency validates consistency across frontend templates
func (tv *TemplateValidator) ValidateTemplateConsistency(templatesPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Find all frontend template directories
	frontendPath := filepath.Join(templatesPath, "frontend")
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "TemplatesPath",
			Tag:     "exists",
			Value:   frontendPath,
			Message: "Frontend templates directory does not exist",
		})
		return result, nil
	}

	// Get all frontend template directories
	templateDirs, err := tv.getFrontendTemplateDirs(frontendPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get template directories: %w", err)
	}

	if len(templateDirs) == 0 {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Templates",
			Message: "No frontend templates found for validation",
		})
		return result, nil
	}

	// Validate package.json consistency
	if err := tv.validatePackageJSONConsistency(templateDirs, result); err != nil {
		return nil, fmt.Errorf("failed to validate package.json consistency: %w", err)
	}

	// Validate TypeScript configuration consistency
	if err := tv.validateTypeScriptConsistency(templateDirs, result); err != nil {
		return nil, fmt.Errorf("failed to validate TypeScript consistency: %w", err)
	}

	// Validate build configuration consistency
	if err := tv.validateBuildConfigConsistency(templateDirs, result); err != nil {
		return nil, fmt.Errorf("failed to validate build config consistency: %w", err)
	}

	return result, nil
}

// ValidatePackageJSONStructure validates a single package.json against standards
func (tv *TemplateValidator) ValidatePackageJSONStructure(packageJSONPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "PackageJSON",
			Tag:     "read",
			Value:   packageJSONPath,
			Message: fmt.Sprintf("Failed to read package.json: %s", err.Error()),
		})
		return result, nil
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "PackageJSON",
			Tag:     "syntax",
			Value:   packageJSONPath,
			Message: fmt.Sprintf("Invalid JSON format: %s", err.Error()),
		})
		return result, nil
	}

	// Validate required fields
	requiredFields := []string{"name", "version", "scripts", "dependencies", "devDependencies", "engines"}
	for _, field := range requiredFields {
		if _, exists := packageJSON[field]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   field,
				Tag:     "required",
				Value:   "",
				Message: fmt.Sprintf("Missing required field: %s", field),
			})
		}
	}

	// Validate scripts structure
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		tv.validateScripts(scripts, result)
	}

	// Validate engines
	if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
		tv.validateEngines(engines, result)
	}

	// Validate dependencies structure
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		tv.validateDependencies(deps, "dependencies", result)
	}

	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		tv.validateDependencies(devDeps, "devDependencies", result)
	}

	return result, nil
}

// ValidateTypeScriptConfig validates TypeScript configuration
func (tv *TemplateValidator) ValidateTypeScriptConfig(tsconfigPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	data, err := os.ReadFile(tsconfigPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "TSConfig",
			Tag:     "read",
			Value:   tsconfigPath,
			Message: fmt.Sprintf("Failed to read tsconfig.json: %s", err.Error()),
		})
		return result, nil
	}

	var tsconfig map[string]interface{}
	if err := json.Unmarshal(data, &tsconfig); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "TSConfig",
			Tag:     "syntax",
			Value:   tsconfigPath,
			Message: fmt.Sprintf("Invalid JSON format: %s", err.Error()),
		})
		return result, nil
	}

	// Validate required sections
	requiredSections := []string{"compilerOptions", "include"}
	for _, section := range requiredSections {
		if _, exists := tsconfig[section]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   section,
				Tag:     "required",
				Value:   "",
				Message: fmt.Sprintf("Missing required section: %s", section),
			})
		}
	}

	// Validate compiler options
	if compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{}); ok {
		tv.validateCompilerOptions(compilerOptions, result)
	}

	return result, nil
}

// getFrontendTemplateDirs returns all frontend template directories
func (tv *TemplateValidator) getFrontendTemplateDirs(frontendPath string) ([]string, error) {
	var templateDirs []string

	entries, err := os.ReadDir(frontendPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "nextjs-") {
			templateDirs = append(templateDirs, filepath.Join(frontendPath, entry.Name()))
		}
	}

	return templateDirs, nil
}

// validatePackageJSONConsistency validates consistency across package.json files
func (tv *TemplateValidator) validatePackageJSONConsistency(templateDirs []string, result *models.ValidationResult) error {
	packageJSONs := make(map[string]map[string]interface{})

	// Read all package.json files
	for _, dir := range templateDirs {
		packageJSONPath := filepath.Join(dir, "package.json.tmpl")
		if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "PackageJSON",
				Message: fmt.Sprintf("Missing package.json.tmpl in %s", filepath.Base(dir)),
			})
			continue
		}

		data, err := os.ReadFile(packageJSONPath)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "PackageJSON",
				Tag:     "read",
				Value:   packageJSONPath,
				Message: fmt.Sprintf("Failed to read package.json: %s", err.Error()),
			})
			continue
		}

		var packageJSON map[string]interface{}
		if err := json.Unmarshal(data, &packageJSON); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "PackageJSON",
				Tag:     "syntax",
				Value:   packageJSONPath,
				Message: fmt.Sprintf("Invalid JSON format: %s", err.Error()),
			})
			continue
		}

		packageJSONs[filepath.Base(dir)] = packageJSON
	}

	// Compare scripts across templates
	tv.compareScriptsConsistency(packageJSONs, result)

	// Compare engines consistency
	tv.compareEnginesConsistency(packageJSONs, result)

	// Compare core dependencies consistency
	tv.compareCoreDependenciesConsistency(packageJSONs, result)

	return nil
}

// validateTypeScriptConsistency validates TypeScript configuration consistency
func (tv *TemplateValidator) validateTypeScriptConsistency(templateDirs []string, result *models.ValidationResult) error {
	tsconfigs := make(map[string]map[string]interface{})

	// Read all tsconfig.json files
	for _, dir := range templateDirs {
		tsconfigPath := filepath.Join(dir, "tsconfig.json.tmpl")
		if _, err := os.Stat(tsconfigPath); os.IsNotExist(err) {
			// Not all templates may have tsconfig.json
			continue
		}

		data, err := os.ReadFile(tsconfigPath)
		if err != nil {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "TSConfig",
				Message: fmt.Sprintf("Failed to read tsconfig.json in %s: %s", filepath.Base(dir), err.Error()),
			})
			continue
		}

		var tsconfig map[string]interface{}
		if err := json.Unmarshal(data, &tsconfig); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "TSConfig",
				Tag:     "syntax",
				Value:   tsconfigPath,
				Message: fmt.Sprintf("Invalid JSON format: %s", err.Error()),
			})
			continue
		}

		tsconfigs[filepath.Base(dir)] = tsconfig
	}

	// Compare compiler options consistency
	tv.compareCompilerOptionsConsistency(tsconfigs, result)

	return nil
}

// validateBuildConfigConsistency validates build configuration consistency
func (tv *TemplateValidator) validateBuildConfigConsistency(templateDirs []string, result *models.ValidationResult) error {
	// Check for consistent build configurations across templates
	for _, dir := range templateDirs {
		templateName := filepath.Base(dir)

		// Check for Next.js config
		nextConfigPath := filepath.Join(dir, "next.config.js.tmpl")
		if _, err := os.Stat(nextConfigPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "NextConfig",
				Message: fmt.Sprintf("Missing next.config.js.tmpl in %s", templateName),
			})
		}

		// Check for Tailwind config
		tailwindConfigPath := filepath.Join(dir, "tailwind.config.js.tmpl")
		if _, err := os.Stat(tailwindConfigPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "TailwindConfig",
				Message: fmt.Sprintf("Missing tailwind.config.js.tmpl in %s", templateName),
			})
		}
	}

	return nil
}

// compareScriptsConsistency compares scripts across package.json files
func (tv *TemplateValidator) compareScriptsConsistency(packageJSONs map[string]map[string]interface{}, result *models.ValidationResult) {
	requiredScripts := []string{"dev", "build", "start", "lint", "type-check", "test", "format", "clean"}

	for templateName, packageJSON := range packageJSONs {
		scripts, ok := packageJSON["scripts"].(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "Scripts",
				Tag:     "type",
				Value:   templateName,
				Message: fmt.Sprintf("Scripts section is not an object in %s", templateName),
			})
			continue
		}

		// Check for required scripts
		for _, requiredScript := range requiredScripts {
			if _, exists := scripts[requiredScript]; !exists {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "Scripts",
					Tag:     "required",
					Value:   fmt.Sprintf("%s.%s", templateName, requiredScript),
					Message: fmt.Sprintf("Missing required script '%s' in %s", requiredScript, templateName),
				})
			}
		}
	}
}

// compareEnginesConsistency compares engines across package.json files
func (tv *TemplateValidator) compareEnginesConsistency(packageJSONs map[string]map[string]interface{}, result *models.ValidationResult) {
	var referenceEngines map[string]interface{}
	var referenceTemplate string

	// Find a reference template with engines
	for templateName, packageJSON := range packageJSONs {
		if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
			referenceEngines = engines
			referenceTemplate = templateName
			break
		}
	}

	if referenceEngines == nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Engines",
			Message: "No engines configuration found in any template",
		})
		return
	}

	// Compare all other templates against reference
	for templateName, packageJSON := range packageJSONs {
		if templateName == referenceTemplate {
			continue
		}

		engines, ok := packageJSON["engines"].(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "Engines",
				Tag:     "missing",
				Value:   templateName,
				Message: fmt.Sprintf("Missing engines configuration in %s", templateName),
			})
			continue
		}

		// Check for inconsistencies
		for engine, version := range referenceEngines {
			if currentVersion, exists := engines[engine]; !exists {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "Engines",
					Tag:     "missing",
					Value:   fmt.Sprintf("%s.%s", templateName, engine),
					Message: fmt.Sprintf("Missing engine '%s' in %s", engine, templateName),
				})
			} else if !reflect.DeepEqual(version, currentVersion) {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field: "Engines",
					Message: fmt.Sprintf("Engine '%s' version mismatch: %s has '%v', %s has '%v'",
						engine, referenceTemplate, version, templateName, currentVersion),
				})
			}
		}
	}
}

// compareCoreDependenciesConsistency compares core dependencies across templates
func (tv *TemplateValidator) compareCoreDependenciesConsistency(packageJSONs map[string]map[string]interface{}, result *models.ValidationResult) {
	coreDependencies := []string{"next", "react", "react-dom", "typescript", "tailwindcss"}

	for templateName, packageJSON := range packageJSONs {
		deps, ok := packageJSON["dependencies"].(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "Dependencies",
				Tag:     "type",
				Value:   templateName,
				Message: fmt.Sprintf("Dependencies section is not an object in %s", templateName),
			})
			continue
		}

		// Check for core dependencies
		for _, coreDep := range coreDependencies {
			if _, exists := deps[coreDep]; !exists {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Dependencies",
					Message: fmt.Sprintf("Missing core dependency '%s' in %s", coreDep, templateName),
				})
			}
		}
	}
}

// compareCompilerOptionsConsistency compares TypeScript compiler options
func (tv *TemplateValidator) compareCompilerOptionsConsistency(tsconfigs map[string]map[string]interface{}, result *models.ValidationResult) {
	if len(tsconfigs) < 2 {
		return // Need at least 2 configs to compare
	}

	var referenceOptions map[string]interface{}
	var referenceTemplate string

	// Find a reference template with compiler options
	for templateName, tsconfig := range tsconfigs {
		if compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{}); ok {
			referenceOptions = compilerOptions
			referenceTemplate = templateName
			break
		}
	}

	if referenceOptions == nil {
		return
	}

	// Core compiler options that should be consistent
	coreOptions := []string{"target", "lib", "strict", "esModuleInterop", "moduleResolution", "jsx"}

	// Compare all other templates against reference
	for templateName, tsconfig := range tsconfigs {
		if templateName == referenceTemplate {
			continue
		}

		compilerOptions, ok := tsconfig["compilerOptions"].(map[string]interface{})
		if !ok {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "CompilerOptions",
				Message: fmt.Sprintf("Missing compilerOptions in %s", templateName),
			})
			continue
		}

		// Check core options consistency
		for _, option := range coreOptions {
			refValue, refExists := referenceOptions[option]
			currentValue, currentExists := compilerOptions[option]

			if refExists && !currentExists {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "CompilerOptions",
					Message: fmt.Sprintf("Missing compiler option '%s' in %s", option, templateName),
				})
			} else if refExists && currentExists && !reflect.DeepEqual(refValue, currentValue) {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field: "CompilerOptions",
					Message: fmt.Sprintf("Compiler option '%s' mismatch: %s has '%v', %s has '%v'",
						option, referenceTemplate, refValue, templateName, currentValue),
				})
			}
		}
	}
}

// validateScripts validates the scripts section of package.json
func (tv *TemplateValidator) validateScripts(scripts map[string]interface{}, result *models.ValidationResult) {
	requiredScripts := []string{"dev", "build", "start", "lint", "type-check", "test"}

	for _, script := range requiredScripts {
		if _, exists := scripts[script]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "Scripts",
				Tag:     "required",
				Value:   script,
				Message: fmt.Sprintf("Missing required script: %s", script),
			})
		}
	}

	// Validate script commands
	for scriptName, command := range scripts {
		if cmdStr, ok := command.(string); ok {
			if strings.TrimSpace(cmdStr) == "" {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Scripts",
					Message: fmt.Sprintf("Empty script command for '%s'", scriptName),
				})
			}
		}
	}
}

// validateEngines validates the engines section of package.json
func (tv *TemplateValidator) validateEngines(engines map[string]interface{}, result *models.ValidationResult) {
	requiredEngines := []string{"node", "npm"}

	for _, engine := range requiredEngines {
		if _, exists := engines[engine]; !exists {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "Engines",
				Tag:     "required",
				Value:   engine,
				Message: fmt.Sprintf("Missing required engine: %s", engine),
			})
		}
	}

	// Validate version formats
	for engine, version := range engines {
		if versionStr, ok := version.(string); ok {
			if !tv.isValidVersionConstraint(versionStr) {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Engines",
					Message: fmt.Sprintf("Invalid version constraint for %s: %s", engine, versionStr),
				})
			}
		}
	}
}

// validateDependencies validates dependencies or devDependencies section
func (tv *TemplateValidator) validateDependencies(deps map[string]interface{}, section string, result *models.ValidationResult) {
	for depName, version := range deps {
		if versionStr, ok := version.(string); ok {
			if strings.TrimSpace(versionStr) == "" {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   section,
					Message: fmt.Sprintf("Empty version for dependency '%s'", depName),
				})
			}
		}
	}
}

// validateCompilerOptions validates TypeScript compiler options
func (tv *TemplateValidator) validateCompilerOptions(compilerOptions map[string]interface{}, result *models.ValidationResult) {
	requiredOptions := []string{"target", "lib", "strict", "esModuleInterop", "moduleResolution"}

	for _, option := range requiredOptions {
		if _, exists := compilerOptions[option]; !exists {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "CompilerOptions",
				Message: fmt.Sprintf("Missing recommended compiler option: %s", option),
			})
		}
	}

	// Validate specific option values
	if target, exists := compilerOptions["target"]; exists {
		if targetStr, ok := target.(string); ok {
			validTargets := []string{"es5", "es6", "es2015", "es2017", "es2018", "es2019", "es2020", "esnext"}
			if !tv.contains(validTargets, targetStr) {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "CompilerOptions",
					Message: fmt.Sprintf("Unusual target value: %s", targetStr),
				})
			}
		}
	}

	if strict, exists := compilerOptions["strict"]; exists {
		if strictBool, ok := strict.(bool); ok && !strictBool {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "CompilerOptions",
				Message: "Strict mode is disabled - consider enabling for better type safety",
			})
		}
	}
}

// Helper functions

// isValidVersionConstraint checks if a version constraint is valid
func (tv *TemplateValidator) isValidVersionConstraint(version string) bool {
	// Basic validation for npm version constraints
	if strings.TrimSpace(version) == "" {
		return false
	}

	// Allow common patterns like "^1.0.0", "~1.0.0", ">=1.0.0", "1.0.0"
	validPrefixes := []string{"^", "~", ">=", "<=", ">", "<", "="}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(version, prefix) {
			return true
		}
	}

	// Check if it's a plain version number
	parts := strings.Split(version, ".")
	return len(parts) >= 2 && len(parts) <= 3
}

// contains checks if a slice contains a string
func (tv *TemplateValidator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
