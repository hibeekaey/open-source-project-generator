package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/models"
)

// ProjectTypeValidator provides validation for specific project types
type ProjectTypeValidator struct {
	engine *Engine
}

// NewProjectTypeValidator creates a new project type validator
func NewProjectTypeValidator() *ProjectTypeValidator {
	return &ProjectTypeValidator{
		engine: NewEngine().(*Engine),
	}
}

// ValidateFrontendProject validates a frontend project structure
func (v *ProjectTypeValidator) ValidateFrontendProject(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check for required frontend files
	requiredFiles := []string{
		"package.json",
		"next.config.js",
		"tailwind.config.js",
		"tsconfig.json",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "FrontendStructure",
				Tag:     "required",
				Value:   file,
				Message: fmt.Sprintf("Required frontend file missing: %s", file),
			})
		}
	}

	// Check for required directories
	requiredDirs := []string{
		"src",
		"src/app",
		"src/components",
		"public",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "FrontendStructure",
				Message: fmt.Sprintf("Recommended frontend directory missing: %s", dir),
			})
		}
	}

	// Validate package.json for frontend-specific dependencies
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if err := v.validateFrontendPackageJSON(packageJSONPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate frontend package.json: %w", err)
	}

	// Validate Next.js configuration
	nextConfigPath := filepath.Join(projectPath, "next.config.js")
	if _, err := os.Stat(nextConfigPath); err == nil {
		if err := v.validateNextConfig(nextConfigPath, result); err != nil {
			return nil, fmt.Errorf("failed to validate Next.js config: %w", err)
		}
	}

	// Validate TypeScript configuration
	tsConfigPath := filepath.Join(projectPath, "tsconfig.json")
	if _, err := os.Stat(tsConfigPath); err == nil {
		if err := v.validateTSConfig(tsConfigPath, result); err != nil {
			return nil, fmt.Errorf("failed to validate TypeScript config: %w", err)
		}
	}

	return result, nil
}

// ValidateBackendProject validates a backend project structure
func (v *ProjectTypeValidator) ValidateBackendProject(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check for required backend files
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"Dockerfile",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "BackendStructure",
				Tag:     "required",
				Value:   file,
				Message: fmt.Sprintf("Required backend file missing: %s", file),
			})
		}
	}

	// Check for recommended directories
	recommendedDirs := []string{
		"internal",
		"internal/handlers",
		"internal/models",
		"internal/services",
		"internal/repositories",
		"pkg",
		"cmd",
	}

	for _, dir := range recommendedDirs {
		dirPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "BackendStructure",
				Message: fmt.Sprintf("Recommended backend directory missing: %s", dir),
			})
		}
	}

	// Validate go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if err := v.engine.ValidateGoMod(goModPath); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "go.mod",
			Tag:     "syntax",
			Value:   "go.mod",
			Message: fmt.Sprintf("Go module validation failed: %s", err.Error()),
		})
	}

	// Validate main.go
	mainGoPath := filepath.Join(projectPath, "main.go")
	if err := v.validateMainGo(mainGoPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate main.go: %w", err)
	}

	return result, nil
}

// ValidateMobileProject validates a mobile project structure
func (v *ProjectTypeValidator) ValidateMobileProject(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check for Android project
	androidPath := filepath.Join(projectPath, "android")
	if _, err := os.Stat(androidPath); err == nil {
		if err := v.validateAndroidProject(androidPath, result); err != nil {
			return nil, fmt.Errorf("failed to validate Android project: %w", err)
		}
	}

	// Check for iOS project
	iosPath := filepath.Join(projectPath, "ios")
	if _, err := os.Stat(iosPath); err == nil {
		if err := v.validateIOSProject(iosPath, result); err != nil {
			return nil, fmt.Errorf("failed to validate iOS project: %w", err)
		}
	}

	// Check if at least one mobile platform exists
	hasAndroid := false
	hasIOS := false
	if _, err := os.Stat(androidPath); err == nil {
		hasAndroid = true
	}
	if _, err := os.Stat(iosPath); err == nil {
		hasIOS = true
	}

	if !hasAndroid && !hasIOS {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "MobileStructure",
			Tag:     "required",
			Value:   "",
			Message: "Mobile project must contain at least one platform (android or ios)",
		})
	}

	return result, nil
}

// ValidateInfrastructureProject validates infrastructure project structure
func (v *ProjectTypeValidator) ValidateInfrastructureProject(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check for Docker files
	dockerFiles := []string{"Dockerfile", "docker-compose.yml", "docker-compose.yaml"}
	hasDocker := false
	for _, file := range dockerFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			hasDocker = true
			break
		}
	}

	// Check for Kubernetes files
	k8sPath := filepath.Join(projectPath, "k8s")
	hasK8s := false
	if _, err := os.Stat(k8sPath); err == nil {
		hasK8s = true
		if err := v.validateKubernetesProject(k8sPath, result); err != nil {
			return nil, fmt.Errorf("failed to validate Kubernetes project: %w", err)
		}
	}

	// Check for Terraform files
	terraformFiles := []string{"main.tf", "variables.tf", "outputs.tf"}
	hasTerraform := false
	for _, file := range terraformFiles {
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			hasTerraform = true
			break
		}
	}

	if !hasDocker && !hasK8s && !hasTerraform {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "InfrastructureStructure",
			Tag:     "required",
			Value:   "",
			Message: "Infrastructure project must contain at least one infrastructure component (Docker, Kubernetes, or Terraform)",
		})
	}

	return result, nil
}

// validateFrontendPackageJSON validates frontend-specific package.json requirements
func (v *ProjectTypeValidator) validateFrontendPackageJSON(packageJSONPath string, result *models.ValidationResult) error {
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return nil // Already handled by caller
	}

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Check for required frontend dependencies
	requiredDeps := []string{"next", "react", "react-dom"}
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		for _, dep := range requiredDeps {
			if _, exists := deps[dep]; !exists {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "dependencies",
					Message: fmt.Sprintf("Missing recommended frontend dependency: %s", dep),
				})
			}
		}
	}

	// Check for required scripts
	requiredScripts := []string{"dev", "build", "start"}
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for _, script := range requiredScripts {
			if _, exists := scripts[script]; !exists {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "scripts",
					Message: fmt.Sprintf("Missing recommended script: %s", script),
				})
			}
		}
	}

	return nil
}

// validateNextConfig validates Next.js configuration
func (v *ProjectTypeValidator) validateNextConfig(nextConfigPath string, result *models.ValidationResult) error {
	data, err := os.ReadFile(nextConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read next.config.js: %w", err)
	}

	content := string(data)

	// Basic validation - check for common configuration patterns
	if !strings.Contains(content, "module.exports") && !strings.Contains(content, "export default") {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "next.config.js",
			Message: "Next.js config should export configuration object",
		})
	}

	return nil
}

// validateTSConfig validates TypeScript configuration
func (v *ProjectTypeValidator) validateTSConfig(tsConfigPath string, result *models.ValidationResult) error {
	data, err := os.ReadFile(tsConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read tsconfig.json: %w", err)
	}

	var tsConfig map[string]interface{}
	if err := json.Unmarshal(data, &tsConfig); err != nil {
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "tsconfig.json",
			Tag:     "syntax",
			Value:   "tsconfig.json",
			Message: fmt.Sprintf("Invalid TypeScript config JSON: %s", err.Error()),
		})
		result.Valid = false
		return nil
	}

	// Check for required compiler options
	if compilerOptions, ok := tsConfig["compilerOptions"].(map[string]interface{}); ok {
		requiredOptions := map[string]interface{}{
			"target":       "es5",
			"lib":          []interface{}{"dom", "dom.iterable", "es6"},
			"allowJs":      true,
			"skipLibCheck": true,
			"strict":       true,
		}

		for option, expectedValue := range requiredOptions {
			if value, exists := compilerOptions[option]; !exists {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "compilerOptions",
					Message: fmt.Sprintf("Missing recommended TypeScript compiler option: %s", option),
				})
			} else if expectedValue != nil {
				// For array values, we'll just check if they exist rather than comparing values
				if option == "lib" {
					// Skip detailed comparison for lib array
					continue
				}
				// For simple values, compare directly
				if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", expectedValue) {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "compilerOptions",
						Message: fmt.Sprintf("TypeScript compiler option %s has unexpected value", option),
					})
				}
			}
		}
	}

	return nil
}

// validateMainGo validates the main.go file
func (v *ProjectTypeValidator) validateMainGo(mainGoPath string, result *models.ValidationResult) error {
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return nil // Already handled by caller
	}

	data, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	content := string(data)

	// Check for package main declaration
	if !strings.Contains(content, "package main") {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "main.go",
			Tag:     "syntax",
			Value:   "package declaration",
			Message: "main.go must have 'package main' declaration",
		})
	}

	// Check for main function
	if !strings.Contains(content, "func main()") {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "main.go",
			Tag:     "syntax",
			Value:   "main function",
			Message: "main.go must have main() function",
		})
	}

	return nil
}

// validateAndroidProject validates Android project structure
func (v *ProjectTypeValidator) validateAndroidProject(androidPath string, result *models.ValidationResult) error {
	// Check for required Android files
	requiredFiles := []string{
		"build.gradle.kts",
		"settings.gradle.kts",
		"app/build.gradle.kts",
		"app/src/main/AndroidManifest.xml",
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(androidPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "AndroidStructure",
				Tag:     "required",
				Value:   file,
				Message: fmt.Sprintf("Required Android file missing: %s", file),
			})
		}
	}

	// Validate AndroidManifest.xml
	manifestPath := filepath.Join(androidPath, "app/src/main/AndroidManifest.xml")
	if _, err := os.Stat(manifestPath); err == nil {
		if err := v.engine.ValidateYAML(manifestPath); err != nil {
			// AndroidManifest.xml is XML, not YAML, but we can still check basic structure
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "AndroidManifest.xml",
				Message: "AndroidManifest.xml may have syntax issues",
			})
		}
	}

	return nil
}

// validateIOSProject validates iOS project structure
func (v *ProjectTypeValidator) validateIOSProject(iosPath string, result *models.ValidationResult) error {
	// Look for .xcodeproj directory
	entries, err := os.ReadDir(iosPath)
	if err != nil {
		return fmt.Errorf("failed to read iOS directory: %w", err)
	}

	hasXcodeProject := false
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcodeproj") {
			hasXcodeProject = true
			break
		}
	}

	if !hasXcodeProject {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "IOSStructure",
			Tag:     "required",
			Value:   ".xcodeproj",
			Message: "iOS project must contain .xcodeproj directory",
		})
	}

	// Check for Podfile
	podfilePath := filepath.Join(iosPath, "Podfile")
	if _, err := os.Stat(podfilePath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "IOSStructure",
			Message: "iOS project should have Podfile for dependency management",
		})
	}

	return nil
}

// validateKubernetesProject validates Kubernetes project structure
func (v *ProjectTypeValidator) validateKubernetesProject(k8sPath string, result *models.ValidationResult) error {
	foundFiles := 0
	entries, err := os.ReadDir(k8sPath)
	if err != nil {
		return fmt.Errorf("failed to read Kubernetes directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
			foundFiles++
			// Validate YAML syntax
			filePath := filepath.Join(k8sPath, entry.Name())
			if err := v.engine.ValidateYAML(filePath); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   entry.Name(),
					Tag:     "syntax",
					Value:   entry.Name(),
					Message: fmt.Sprintf("Kubernetes YAML validation failed: %s", err.Error()),
				})
			}
		}
	}

	if foundFiles == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "KubernetesStructure",
			Tag:     "required",
			Value:   "",
			Message: "Kubernetes directory must contain at least one YAML manifest file",
		})
	}

	return nil
}
