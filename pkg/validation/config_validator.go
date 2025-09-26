package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"gopkg.in/yaml.v3"
)

// ConfigValidator provides specialized configuration file validation
type ConfigValidator struct {
	schemas map[string]*interfaces.ConfigSchema
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	validator := &ConfigValidator{
		schemas: make(map[string]*interfaces.ConfigSchema),
	}
	validator.initializeDefaultSchemas()
	return validator
}

// ValidateConfigurationFiles validates all configuration files in a project
func (cv *ConfigValidator) ValidateConfigurationFiles(projectPath string) ([]*interfaces.ConfigValidationResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(projectPath); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
	}

	var results []*interfaces.ConfigValidationResult

	// Find and validate configuration files
	configFiles, err := cv.findConfigurationFiles(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find configuration files: %w", err)
	}

	for _, configFile := range configFiles {
		result, err := cv.validateConfigurationFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s: %w", configFile, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// ValidateConfigurationFile validates a single configuration file
func (cv *ConfigValidator) validateConfigurationFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	// Determine file type and validate accordingly
	ext := strings.ToLower(filepath.Ext(filePath))
	fileName := filepath.Base(filePath)

	switch ext {
	case ".json":
		return cv.validateJSONFile(filePath)
	case ".yaml", ".yml":
		return cv.validateYAMLFile(filePath)
	case ".env":
		return cv.validateEnvFile(filePath)
	case ".toml":
		return cv.validateTOMLFile(filePath)
	default:
		// Check for specific configuration files without extensions
		switch fileName {
		case "Dockerfile":
			return cv.validateDockerfile(filePath)
		case "Makefile":
			return cv.validateMakefile(filePath)
		case ".gitignore":
			return cv.validateGitignore(filePath)
		case ".dockerignore":
			return cv.validateDockerignore(filePath)
		default:
			// Unknown configuration file type
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:    "file_type",
				Value:    fileName,
				Type:     "unknown",
				Message:  "Unknown configuration file type",
				Severity: interfaces.ValidationSeverityInfo,
				Rule:     "config.unknown_type",
			})
			result.Summary.WarningCount++
		}
	}

	return result, nil
}

// validateJSONFile validates a JSON configuration file
func (cv *ConfigValidator) validateJSONFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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
	if schema, exists := cv.schemas[fileName]; exists {
		if err := cv.validateAgainstSchema(data, schema, result); err != nil {
			return nil, fmt.Errorf("schema validation failed: %w", err)
		}
	}

	// Perform specific validations based on file name
	switch fileName {
	case "package.json":
		cv.validatePackageJSON(data, result)
	case "tsconfig.json":
		cv.validateTSConfig(data, result)
	case ".eslintrc.json":
		cv.validateESLintConfig(data, result)
	}

	return result, nil
}

// validateYAMLFile validates a YAML configuration file
func (cv *ConfigValidator) validateYAMLFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	// Read and parse YAML
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
	if err := yaml.Unmarshal(content, &data); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "syntax",
			Value:    string(content),
			Type:     "syntax_error",
			Message:  fmt.Sprintf("Invalid YAML syntax: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.yaml.syntax",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	// Validate against schema if available
	fileName := filepath.Base(filePath)
	if schema, exists := cv.schemas[fileName]; exists {
		if err := cv.validateAgainstSchema(data, schema, result); err != nil {
			return nil, fmt.Errorf("schema validation failed: %w", err)
		}
	}

	// Perform specific validations based on file name
	switch fileName {
	case "docker-compose.yml", "docker-compose.yaml":
		cv.validateDockerCompose(data, result)
	case ".github/workflows/ci.yml", ".github/workflows/ci.yaml":
		cv.validateGitHubWorkflow(data, result)
	}

	return result, nil
}

// validateEnvFile validates environment configuration files
func (cv *ConfigValidator) validateEnvFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.Summary.TotalProperties++

		// Validate environment variable format
		if !strings.Contains(line, "=") {
			result.Valid = false
			result.Errors = append(result.Errors, interfaces.ConfigValidationError{
				Field:    fmt.Sprintf("line_%d", lineNum+1),
				Value:    line,
				Type:     "format_error",
				Message:  "Environment variable must be in KEY=VALUE format",
				Severity: interfaces.ValidationSeverityError,
				Rule:     "config.env.format",
			})
			result.Summary.ErrorCount++
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]

		// Validate key format
		if err := cv.validateEnvKey(key); err != nil {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      key,
				Type:       "key_format",
				Message:    fmt.Sprintf("Invalid environment variable name: %v", err),
				Suggestion: "Use uppercase letters, numbers, and underscores only",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.env.key_format",
			})
			result.Summary.WarningCount++
		}

		// Check for potential secrets
		if cv.isPotentialSecret(key, value) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      key,
				Type:       "security",
				Message:    "Potential secret detected in environment file",
				Suggestion: "Consider using a secrets management system",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.env.secrets",
			})
			result.Summary.WarningCount++
		}

		result.Summary.ValidProperties++
	}

	return result, nil
}

// validateTOMLFile validates TOML configuration files
func (cv *ConfigValidator) validateTOMLFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	// Basic TOML syntax validation (simplified)
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

	// Simple TOML validation - check for basic syntax issues
	lines := strings.Split(string(content), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.Summary.TotalProperties++

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			result.Summary.ValidProperties++
			continue
		}

		// Check for key-value pairs
		if strings.Contains(line, "=") {
			result.Summary.ValidProperties++
			continue
		}

		// Invalid line format
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:    fmt.Sprintf("line_%d", lineNum+1),
			Value:    line,
			Type:     "format_warning",
			Message:  "Potentially invalid TOML syntax",
			Severity: interfaces.ValidationSeverityWarning,
			Rule:     "config.toml.syntax",
		})
		result.Summary.WarningCount++
	}

	return result, nil
}

// validateDockerfile validates Dockerfile configuration
func (cv *ConfigValidator) validateDockerfile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	lines := strings.Split(string(content), "\n")
	hasFrom := false

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.Summary.TotalProperties++

		// Check for FROM instruction
		if strings.HasPrefix(strings.ToUpper(line), "FROM ") {
			hasFrom = true
			result.Summary.ValidProperties++

			// Check for latest tag usage
			if strings.Contains(line, ":latest") {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      fmt.Sprintf("line_%d", lineNum+1),
					Value:      line,
					Type:       "best_practice",
					Message:    "Avoid using 'latest' tag in production",
					Suggestion: "Use specific version tags for better reproducibility",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker.latest_tag",
				})
				result.Summary.WarningCount++
			}
			continue
		}

		// Check for USER instruction
		if strings.HasPrefix(strings.ToUpper(line), "USER ") {
			if strings.Contains(line, "root") {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      fmt.Sprintf("line_%d", lineNum+1),
					Value:      line,
					Type:       "security",
					Message:    "Running as root user is a security risk",
					Suggestion: "Create and use a non-root user",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker.root_user",
				})
				result.Summary.WarningCount++
			}
		}

		result.Summary.ValidProperties++
	}

	// Check if FROM instruction exists
	if !hasFrom {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "dockerfile",
			Value:    filePath,
			Type:     "missing_instruction",
			Message:  "Dockerfile must contain a FROM instruction",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "docker.from_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	return result, nil
}

// validateMakefile validates Makefile configuration
func (cv *ConfigValidator) validateMakefile(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	lines := strings.Split(string(content), "\n")

	for lineNum, line := range lines {
		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		result.Summary.TotalProperties++

		// Check for target definitions (lines ending with :)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "\t") {
			result.Summary.ValidProperties++
			continue
		}

		// Check for command lines (should start with tab)
		if strings.HasPrefix(line, "\t") {
			result.Summary.ValidProperties++
			continue
		}

		// Check for variable assignments
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "\t") {
			result.Summary.ValidProperties++
			continue
		}

		// Potential syntax issue
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:    fmt.Sprintf("line_%d", lineNum+1),
			Value:    line,
			Type:     "syntax_warning",
			Message:  "Potential Makefile syntax issue",
			Severity: interfaces.ValidationSeverityWarning,
			Rule:     "makefile.syntax",
		})
		result.Summary.WarningCount++
	}

	return result, nil
}

// validateGitignore validates .gitignore file
func (cv *ConfigValidator) validateGitignore(filePath string) (*interfaces.ConfigValidationResult, error) {
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

	lines := strings.Split(string(content), "\n")
	commonPatterns := []string{"node_modules/", "*.log", ".env", "dist/", "build/"}
	foundPatterns := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.Summary.TotalProperties++
		result.Summary.ValidProperties++

		// Check for common patterns
		for _, pattern := range commonPatterns {
			if line == pattern {
				foundPatterns[pattern] = true
			}
		}
	}

	// Suggest missing common patterns
	for _, pattern := range commonPatterns {
		if !foundPatterns[pattern] {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "patterns",
				Value:      pattern,
				Type:       "missing_pattern",
				Message:    fmt.Sprintf("Consider adding common pattern: %s", pattern),
				Suggestion: fmt.Sprintf("Add '%s' to ignore common files", pattern),
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "gitignore.common_patterns",
			})
			result.Summary.WarningCount++
		}
	}

	return result, nil
}

// validateDockerignore validates .dockerignore file
func (cv *ConfigValidator) validateDockerignore(filePath string) (*interfaces.ConfigValidationResult, error) {
	// Similar to gitignore validation but with Docker-specific patterns
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

	lines := strings.Split(string(content), "\n")
	dockerPatterns := []string{"node_modules", ".git", "*.md", "Dockerfile", ".dockerignore"}
	foundPatterns := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result.Summary.TotalProperties++
		result.Summary.ValidProperties++

		// Check for common Docker patterns
		for _, pattern := range dockerPatterns {
			if line == pattern {
				foundPatterns[pattern] = true
			}
		}
	}

	// Suggest missing Docker-specific patterns
	for _, pattern := range dockerPatterns {
		if !foundPatterns[pattern] {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      "patterns",
				Value:      pattern,
				Type:       "missing_pattern",
				Message:    fmt.Sprintf("Consider adding Docker pattern: %s", pattern),
				Suggestion: fmt.Sprintf("Add '%s' to optimize Docker builds", pattern),
				Severity:   interfaces.ValidationSeverityInfo,
				Rule:       "dockerignore.common_patterns",
			})
			result.Summary.WarningCount++
		}
	}

	return result, nil
}

// Helper methods for specific file validations

// validatePackageJSON validates package.json specific structure
func (cv *ConfigValidator) validatePackageJSON(data interface{}, result *interfaces.ConfigValidationResult) {
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
			if err := cv.validatePackageName(nameStr); err != nil {
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
}

// validateTSConfig validates TypeScript configuration
func (cv *ConfigValidator) validateTSConfig(data interface{}, result *interfaces.ConfigValidationResult) {
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
func (cv *ConfigValidator) validateESLintConfig(data interface{}, result *interfaces.ConfigValidationResult) {
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
}

// validateDockerCompose validates Docker Compose configuration
func (cv *ConfigValidator) validateDockerCompose(data interface{}, result *interfaces.ConfigValidationResult) {
	compose, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for version
	result.Summary.TotalProperties++
	if version, exists := compose["version"]; exists {
		result.Summary.ValidProperties++
		if versionStr, ok := version.(string); ok {
			if versionStr == "2" || strings.HasPrefix(versionStr, "2.") {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      "version",
					Value:      versionStr,
					Type:       "deprecated",
					Message:    "Docker Compose version 2.x is deprecated",
					Suggestion: "Consider upgrading to version 3.x",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker_compose.version",
				})
				result.Summary.WarningCount++
			}
		}
	}

	// Check for services
	result.Summary.TotalProperties++
	if services, exists := compose["services"]; exists {
		result.Summary.ValidProperties++

		if servicesMap, ok := services.(map[string]interface{}); ok {
			for serviceName, service := range servicesMap {
				if serviceMap, ok := service.(map[string]interface{}); ok {
					// Check for privileged mode
					if privileged, exists := serviceMap["privileged"]; exists {
						if privilegedBool, ok := privileged.(bool); ok && privilegedBool {
							result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
								Field:      fmt.Sprintf("services.%s.privileged", serviceName),
								Value:      "true",
								Type:       "security",
								Message:    "Privileged mode is a security risk",
								Suggestion: "Avoid using privileged mode unless absolutely necessary",
								Severity:   interfaces.ValidationSeverityWarning,
								Rule:       "docker_compose.privileged",
							})
							result.Summary.WarningCount++
						}
					}
				}
			}
		}
	}
}

// validateGitHubWorkflow validates GitHub Actions workflow
func (cv *ConfigValidator) validateGitHubWorkflow(data interface{}, result *interfaces.ConfigValidationResult) {
	workflow, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	// Check for name
	result.Summary.TotalProperties++
	if _, exists := workflow["name"]; exists {
		result.Summary.ValidProperties++
	} else {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "name",
			Value:      "",
			Type:       "missing_recommended",
			Message:    "Workflow name is recommended for clarity",
			Suggestion: "Add a descriptive name for the workflow",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "github_workflow.name",
		})
		result.Summary.WarningCount++
	}

	// Check for on (triggers)
	result.Summary.TotalProperties++
	if _, exists := workflow["on"]; exists {
		result.Summary.ValidProperties++
	} else {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "on",
			Value:    "",
			Type:     "missing_required",
			Message:  "Workflow triggers ('on') are required",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "github_workflow.triggers",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}
}

// Helper methods

// findConfigurationFiles finds all configuration files in a project
func (cv *ConfigValidator) findConfigurationFiles(projectPath string) ([]string, error) {
	var configFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is a configuration file
		if cv.isConfigurationFile(info.Name()) {
			configFiles = append(configFiles, path)
		}

		return nil
	})

	return configFiles, err
}

// isConfigurationFile checks if a file is a configuration file
func (cv *ConfigValidator) isConfigurationFile(fileName string) bool {
	configExtensions := []string{".json", ".yaml", ".yml", ".toml", ".env"}
	configFiles := []string{"Dockerfile", "Makefile", ".gitignore", ".dockerignore"}

	// Check extensions
	for _, ext := range configExtensions {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return true
		}
	}

	// Check specific file names
	for _, configFile := range configFiles {
		if fileName == configFile {
			return true
		}
	}

	return false
}

// validateAgainstSchema validates data against a configuration schema
func (cv *ConfigValidator) validateAgainstSchema(data interface{}, schema *interfaces.ConfigSchema, result *interfaces.ConfigValidationResult) error {
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
			if err := cv.validatePropertyAgainstSchema(key, value, propSchema, result); err != nil {
				return fmt.Errorf("validation failed for property '%s': %w", key, err)
			}
		}
	}

	return nil
}

// validatePropertyAgainstSchema validates a single property against its schema
func (cv *ConfigValidator) validatePropertyAgainstSchema(key string, value interface{}, schema interfaces.PropertySchema, result *interfaces.ConfigValidationResult) error {
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
	}

	return nil
}

// validatePackageName validates NPM package name format
func (cv *ConfigValidator) validatePackageName(name string) error {
	// NPM package name rules
	if len(name) > 214 {
		return fmt.Errorf("name too long")
	}

	if strings.ToLower(name) != name {
		return fmt.Errorf("name must be lowercase")
	}

	// Check for invalid characters
	invalidChars := regexp.MustCompile(`[^a-z0-9\-._~]`)
	if invalidChars.MatchString(name) {
		return fmt.Errorf("name contains invalid characters")
	}

	return nil
}

// validateEnvKey validates environment variable key format
func (cv *ConfigValidator) validateEnvKey(key string) error {
	// Environment variable names should be uppercase with underscores
	envKeyRegex := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	if !envKeyRegex.MatchString(key) {
		return fmt.Errorf("environment variable names should be uppercase with underscores")
	}
	return nil
}

// isPotentialSecret checks if a key-value pair might contain a secret
func (cv *ConfigValidator) isPotentialSecret(key, value string) bool {
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
func (cv *ConfigValidator) initializeDefaultSchemas() {
	// Package.json schema
	cv.schemas["package.json"] = &interfaces.ConfigSchema{
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
		},
	}
}
