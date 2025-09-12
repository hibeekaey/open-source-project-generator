package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/models"
)

// VercelValidator handles Vercel deployment validation
type VercelValidator struct {
	standardConfig map[string]interface{}
}

// NewVercelValidator creates a new Vercel validator
func NewVercelValidator() *VercelValidator {
	return &VercelValidator{
		standardConfig: make(map[string]interface{}),
	}
}

// ValidateVercelCompatibility validates Vercel deployment compatibility
func (vv *VercelValidator) ValidateVercelCompatibility(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Check for package.json
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if err := vv.validatePackageJSONForVercel(packageJSONPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate package.json for Vercel: %w", err)
	}

	// Check for Vercel configuration
	vercelConfigPath := filepath.Join(projectPath, "vercel.json")
	vercelResult, err := vv.validateVercelConfig(vercelConfigPath, result)
	if err != nil {
		return nil, fmt.Errorf("failed to validate vercel.json: %w", err)
	}
	// Merge results
	result.Valid = result.Valid && vercelResult.Valid
	result.Errors = append(result.Errors, vercelResult.Errors...)
	result.Warnings = append(result.Warnings, vercelResult.Warnings...)

	// Check for Next.js configuration
	nextConfigPath := filepath.Join(projectPath, "next.config.js")
	if err := vv.validateNextConfigForVercel(nextConfigPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate next.config.js for Vercel: %w", err)
	}

	// Validate build configuration
	if err := vv.validateBuildConfiguration(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate build configuration: %w", err)
	}

	// Validate environment variables setup
	if err := vv.validateEnvironmentVariables(projectPath, result); err != nil {
		return nil, fmt.Errorf("failed to validate environment variables: %w", err)
	}

	return result, nil
}

// ValidateVercelConfig validates a vercel.json configuration file
func (vv *VercelValidator) ValidateVercelConfig(vercelConfigPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	return vv.validateVercelConfig(vercelConfigPath, result)
}

// ValidateEnvironmentVariablesConsistency validates environment variables across templates
func (vv *VercelValidator) ValidateEnvironmentVariablesConsistency(templatesPath string) (*models.ValidationResult, error) {
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

	templateDirs, err := vv.getFrontendTemplateDirs(frontendPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get template directories: %w", err)
	}

	// Validate environment variables consistency
	if err := vv.validateEnvVarsConsistency(templateDirs, result); err != nil {
		return nil, fmt.Errorf("failed to validate environment variables consistency: %w", err)
	}

	return result, nil
}

// validatePackageJSONForVercel validates package.json for Vercel compatibility
func (vv *VercelValidator) validatePackageJSONForVercel(packageJSONPath string, result *models.ValidationResult) error {
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "PackageJSON",
			Tag:     "exists",
			Value:   packageJSONPath,
			Message: "package.json is required for Vercel deployment",
		})
		return nil
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
		return nil
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
		return nil
	}

	// Validate required scripts for Vercel
	requiredScripts := []string{"build", "start"}
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for _, script := range requiredScripts {
			if _, exists := scripts[script]; !exists {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "Scripts",
					Tag:     "required",
					Value:   script,
					Message: fmt.Sprintf("Missing required script for Vercel deployment: %s", script),
				})
			}
		}

		// Validate build script content
		if buildScript, exists := scripts["build"]; exists {
			if buildStr, ok := buildScript.(string); ok {
				if !strings.Contains(buildStr, "next build") && !strings.Contains(buildStr, "npm run build") {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "Scripts",
						Message: "Build script should use 'next build' for optimal Vercel deployment",
					})
				}
			}
		}
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "Scripts",
			Tag:     "missing",
			Value:   "",
			Message: "Scripts section is required for Vercel deployment",
		})
	}

	// Validate engines for Node.js version
	if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
		if nodeVersion, exists := engines["node"]; exists {
			if nodeStr, ok := nodeVersion.(string); ok {
				if !vv.isValidNodeVersionForVercel(nodeStr) {
					result.Warnings = append(result.Warnings, models.ValidationWarning{
						Field:   "Engines",
						Message: fmt.Sprintf("Node.js version %s may not be supported by Vercel", nodeStr),
					})
				}
			}
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "Engines",
				Message: "Node.js version not specified - Vercel will use default version",
			})
		}
	}

	// Check for Next.js dependency
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		if _, exists := deps["next"]; !exists {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "Dependencies",
				Message: "Next.js dependency not found - ensure framework is correctly specified",
			})
		}
	}

	return nil
}

// validateVercelConfig validates vercel.json configuration
func (vv *VercelValidator) validateVercelConfig(vercelConfigPath string, result *models.ValidationResult) (*models.ValidationResult, error) {
	if _, err := os.Stat(vercelConfigPath); os.IsNotExist(err) {
		// vercel.json is optional, but we can provide recommendations
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "VercelConfig",
			Message: "vercel.json not found - consider adding for custom deployment configuration",
		})
		return result, nil
	}

	data, err := os.ReadFile(vercelConfigPath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "VercelConfig",
			Tag:     "read",
			Value:   vercelConfigPath,
			Message: fmt.Sprintf("Failed to read vercel.json: %s", err.Error()),
		})
		return result, nil
	}

	var vercelConfig map[string]interface{}
	if err := json.Unmarshal(data, &vercelConfig); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "VercelConfig",
			Tag:     "syntax",
			Value:   vercelConfigPath,
			Message: fmt.Sprintf("Invalid JSON format: %s", err.Error()),
		})
		return result, nil
	}

	// Validate framework specification
	if framework, exists := vercelConfig["framework"]; exists {
		if frameworkStr, ok := framework.(string); ok {
			validFrameworks := []string{"nextjs", "react", "vue", "nuxtjs", "gatsby", "svelte"}
			if !vv.contains(validFrameworks, frameworkStr) {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Framework",
					Message: fmt.Sprintf("Unusual framework specified: %s", frameworkStr),
				})
			}
		}
	} else {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Framework",
			Message: "Framework not specified - Vercel will auto-detect",
		})
	}

	// Validate build configuration
	if buildCommand, exists := vercelConfig["buildCommand"]; exists {
		if buildStr, ok := buildCommand.(string); ok {
			if strings.TrimSpace(buildStr) == "" {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "BuildCommand",
					Message: "Empty build command specified",
				})
			}
		}
	}

	// Validate security headers
	if headers, exists := vercelConfig["headers"]; exists {
		if headersArray, ok := headers.([]interface{}); ok {
			vv.validateSecurityHeaders(headersArray, result)
		}
	} else {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Headers",
			Message: "No security headers configured - consider adding for better security",
		})
	}

	// Validate functions configuration
	if functions, exists := vercelConfig["functions"]; exists {
		if functionsMap, ok := functions.(map[string]interface{}); ok {
			vv.validateFunctionsConfig(functionsMap, result)
		}
	}

	return result, nil
}

// validateNextConfigForVercel validates next.config.js for Vercel compatibility
func (vv *VercelValidator) validateNextConfigForVercel(nextConfigPath string, result *models.ValidationResult) error {
	if _, err := os.Stat(nextConfigPath); os.IsNotExist(err) {
		// next.config.js is optional
		return nil
	}

	// Read the file content (basic validation since it's JavaScript)
	data, err := os.ReadFile(nextConfigPath)
	if err != nil {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "NextConfig",
			Message: fmt.Sprintf("Failed to read next.config.js: %s", err.Error()),
		})
		return nil
	}

	content := string(data)

	// Check for common Vercel-incompatible configurations
	if strings.Contains(content, "output: 'export'") {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "NextConfig",
			Message: "Static export mode detected - ensure this is intended for Vercel deployment",
		})
	}

	if strings.Contains(content, "trailingSlash: true") {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "NextConfig",
			Message: "Trailing slash configuration may affect Vercel routing",
		})
	}

	// Check for experimental features that might not be supported
	if strings.Contains(content, "experimental") {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "NextConfig",
			Message: "Experimental features detected - verify Vercel compatibility",
		})
	}

	return nil
}

// validateBuildConfiguration validates build-related configuration
func (vv *VercelValidator) validateBuildConfiguration(projectPath string, result *models.ValidationResult) error {
	// Check for .vercelignore file
	vercelIgnorePath := filepath.Join(projectPath, ".vercelignore")
	if _, err := os.Stat(vercelIgnorePath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "VercelIgnore",
			Message: "Consider adding .vercelignore to exclude unnecessary files from deployment",
		})
	}

	// Check for public directory
	publicPath := filepath.Join(projectPath, "public")
	if _, err := os.Stat(publicPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "PublicDirectory",
			Message: "Public directory not found - static assets may not be served correctly",
		})
	}

	// Check for pages or app directory (Next.js structure)
	pagesPath := filepath.Join(projectPath, "pages")
	appPath := filepath.Join(projectPath, "src", "app")
	srcPagesPath := filepath.Join(projectPath, "src", "pages")

	hasPages := false
	if _, err := os.Stat(pagesPath); err == nil {
		hasPages = true
	}
	if _, err := os.Stat(appPath); err == nil {
		hasPages = true
	}
	if _, err := os.Stat(srcPagesPath); err == nil {
		hasPages = true
	}

	if !hasPages {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "ProjectStructure",
			Message: "No pages or app directory found - verify Next.js project structure",
		})
	}

	return nil
}

// validateEnvironmentVariables validates environment variables setup
func (vv *VercelValidator) validateEnvironmentVariables(projectPath string, result *models.ValidationResult) error {
	// Check for .env files
	envFiles := []string{".env.local", ".env.example", ".env"}
	foundEnvFiles := []string{}

	for _, envFile := range envFiles {
		envPath := filepath.Join(projectPath, envFile)
		if _, err := os.Stat(envPath); err == nil {
			foundEnvFiles = append(foundEnvFiles, envFile)
		}
	}

	if len(foundEnvFiles) == 0 {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "EnvironmentVariables",
			Message: "No environment variable files found - consider adding .env.example for documentation",
		})
	}

	// Check for .env.example if .env.local exists
	hasEnvLocal := vv.contains(foundEnvFiles, ".env.local")
	hasEnvExample := vv.contains(foundEnvFiles, ".env.example")

	if hasEnvLocal && !hasEnvExample {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "EnvironmentVariables",
			Message: "Consider adding .env.example to document required environment variables",
		})
	}

	return nil
}

// validateEnvVarsConsistency validates environment variables consistency across templates
func (vv *VercelValidator) validateEnvVarsConsistency(templateDirs []string, result *models.ValidationResult) error {
	envVars := make(map[string]map[string]bool) // template -> env vars

	// Debug logging
	result.Warnings = append(result.Warnings, models.ValidationWarning{
		Field:   "Debug",
		Message: fmt.Sprintf("Processing %d template directories", len(templateDirs)),
	})

	for _, dir := range templateDirs {
		templateName := filepath.Base(dir)
		envVars[templateName] = make(map[string]bool)

		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Debug",
			Message: fmt.Sprintf("Processing template: %s", templateName),
		})

		// Check vercel.json for environment variables
		vercelConfigPath := filepath.Join(dir, "vercel.json.tmpl")
		if _, err := os.Stat(vercelConfigPath); err == nil {
			if vars, err := vv.extractEnvVarsFromVercelConfig(vercelConfigPath); err == nil {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Debug",
					Message: fmt.Sprintf("Found %d env vars in %s vercel.json: %v", len(vars), templateName, vars),
				})
				for _, envVar := range vars {
					envVars[templateName][envVar] = true
				}
			} else {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Debug",
					Message: fmt.Sprintf("Error extracting env vars from %s vercel.json: %v", templateName, err),
				})
			}
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "Debug",
				Message: fmt.Sprintf("No vercel.json.tmpl found for %s", templateName),
			})
		}

		// Check .env.example files
		envExamplePath := filepath.Join(dir, ".env.example.tmpl")
		if _, err := os.Stat(envExamplePath); err == nil {
			if vars, err := vv.extractEnvVarsFromEnvFile(envExamplePath); err == nil {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Debug",
					Message: fmt.Sprintf("Found %d env vars in %s .env.example: %v", len(vars), templateName, vars),
				})
				for _, envVar := range vars {
					envVars[templateName][envVar] = true
				}
			} else {
				result.Warnings = append(result.Warnings, models.ValidationWarning{
					Field:   "Debug",
					Message: fmt.Sprintf("Error extracting env vars from %s .env.example: %v", templateName, err),
				})
			}
		} else {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "Debug",
				Message: fmt.Sprintf("No .env.example.tmpl found for %s", templateName),
			})
		}
	}

	// Compare environment variables across templates
	vv.compareEnvVarsConsistency(envVars, result)

	return nil
}

// validateSecurityHeaders validates security headers configuration
func (vv *VercelValidator) validateSecurityHeaders(headers []interface{}, result *models.ValidationResult) {
	requiredSecurityHeaders := []string{
		"X-Frame-Options",
		"X-Content-Type-Options",
		"Referrer-Policy",
	}

	foundHeaders := make(map[string]bool)

	for _, headerConfig := range headers {
		if headerMap, ok := headerConfig.(map[string]interface{}); ok {
			if headersArray, exists := headerMap["headers"]; exists {
				if headersList, ok := headersArray.([]interface{}); ok {
					for _, header := range headersList {
						if headerObj, ok := header.(map[string]interface{}); ok {
							if key, exists := headerObj["key"]; exists {
								if keyStr, ok := key.(string); ok {
									foundHeaders[keyStr] = true
								}
							}
						}
					}
				}
			}
		}
	}

	for _, requiredHeader := range requiredSecurityHeaders {
		if !foundHeaders[requiredHeader] {
			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field:   "SecurityHeaders",
				Message: fmt.Sprintf("Missing recommended security header: %s", requiredHeader),
			})
		}
	}
}

// validateFunctionsConfig validates functions configuration
func (vv *VercelValidator) validateFunctionsConfig(functions map[string]interface{}, result *models.ValidationResult) {
	for pattern, config := range functions {
		if configMap, ok := config.(map[string]interface{}); ok {
			// Validate maxDuration
			if maxDuration, exists := configMap["maxDuration"]; exists {
				if duration, ok := maxDuration.(float64); ok {
					if duration > 300 { // 5 minutes max for hobby plan
						result.Warnings = append(result.Warnings, models.ValidationWarning{
							Field:   "Functions",
							Message: fmt.Sprintf("Function %s has maxDuration > 300s, may require Pro plan", pattern),
						})
					}
				}
			}

			// Validate memory
			if memory, exists := configMap["memory"]; exists {
				if memoryVal, ok := memory.(float64); ok {
					if memoryVal > 1024 { // 1GB max for hobby plan
						result.Warnings = append(result.Warnings, models.ValidationWarning{
							Field:   "Functions",
							Message: fmt.Sprintf("Function %s has memory > 1024MB, may require Pro plan", pattern),
						})
					}
				}
			}
		}
	}
}

// Helper functions

// getFrontendTemplateDirs returns all frontend template directories
func (vv *VercelValidator) getFrontendTemplateDirs(frontendPath string) ([]string, error) {
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

// extractEnvVarsFromVercelConfig extracts environment variables from vercel.json
func (vv *VercelValidator) extractEnvVarsFromVercelConfig(vercelConfigPath string) ([]string, error) {
	data, err := os.ReadFile(vercelConfigPath)
	if err != nil {
		return nil, err
	}

	var vercelConfig map[string]interface{}
	if err := json.Unmarshal(data, &vercelConfig); err != nil {
		return nil, err
	}

	var envVars []string

	if env, exists := vercelConfig["env"]; exists {
		if envMap, ok := env.(map[string]interface{}); ok {
			for key := range envMap {
				envVars = append(envVars, key)
			}
		}
	}

	if build, exists := vercelConfig["build"]; exists {
		if buildMap, ok := build.(map[string]interface{}); ok {
			if buildEnv, exists := buildMap["env"]; exists {
				if buildEnvMap, ok := buildEnv.(map[string]interface{}); ok {
					for key := range buildEnvMap {
						envVars = append(envVars, key)
					}
				}
			}
		}
	}

	return envVars, nil
}

// extractEnvVarsFromEnvFile extracts environment variables from .env file
func (vv *VercelValidator) extractEnvVarsFromEnvFile(envFilePath string) ([]string, error) {
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, err
	}

	var envVars []string
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				envVar := strings.TrimSpace(parts[0])
				if envVar != "" {
					envVars = append(envVars, envVar)
				}
			}
		}
	}

	return envVars, nil
}

// compareEnvVarsConsistency compares environment variables across templates
func (vv *VercelValidator) compareEnvVarsConsistency(envVars map[string]map[string]bool, result *models.ValidationResult) {
	// Find common environment variables
	commonVars := make(map[string]int)
	totalTemplates := len(envVars)

	result.Warnings = append(result.Warnings, models.ValidationWarning{
		Field:   "Debug",
		Message: fmt.Sprintf("Comparing env vars across %d templates", totalTemplates),
	})

	for _, templateVars := range envVars {
		for envVar := range templateVars {
			commonVars[envVar]++
		}
	}

	result.Warnings = append(result.Warnings, models.ValidationWarning{
		Field:   "Debug",
		Message: fmt.Sprintf("Common vars found: %v", commonVars),
	})

	// Check for inconsistencies
	for envVar, count := range commonVars {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   "Debug",
			Message: fmt.Sprintf("Checking %s: count=%d, total=%d", envVar, count, totalTemplates),
		})

		if count > 0 && count < totalTemplates {
			// This env var exists in some but not all templates
			missingIn := []string{}
			for templateName, templateVars := range envVars {
				if !templateVars[envVar] {
					missingIn = append(missingIn, templateName)
				}
			}

			result.Warnings = append(result.Warnings, models.ValidationWarning{
				Field: "EnvironmentVariables",
				Message: fmt.Sprintf("Environment variable '%s' is missing in templates: %s",
					envVar, strings.Join(missingIn, ", ")),
			})
		}
	}
}

// isValidNodeVersionForVercel checks if Node.js version is supported by Vercel
func (vv *VercelValidator) isValidNodeVersionForVercel(version string) bool {
	// Vercel supports Node.js 18.x, 20.x, and 22.x as of 2024
	supportedVersions := []string{"18", "20", "22"}

	for _, supported := range supportedVersions {
		if strings.Contains(version, supported) {
			return true
		}
	}

	return false
}

// contains checks if a slice contains a string
func (vv *VercelValidator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
