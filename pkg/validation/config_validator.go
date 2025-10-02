package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation/formats"
)

// ConfigValidator provides specialized configuration file validation
type ConfigValidator struct {
	schemaManager     *SchemaManager
	jsonValidator     *formats.JSONValidator
	yamlValidator     *formats.YAMLValidator
	envValidator      *formats.EnvValidator
	dockerValidator   *formats.DockerValidator
	makefileValidator *formats.MakefileValidator
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		schemaManager:     NewSchemaManager(),
		jsonValidator:     formats.NewJSONValidator(),
		yamlValidator:     formats.NewYAMLValidator(),
		envValidator:      formats.NewEnvValidator(),
		dockerValidator:   formats.NewDockerValidator(),
		makefileValidator: formats.NewMakefileValidator(),
	}
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
	// Determine file type and delegate to appropriate format validator
	ext := strings.ToLower(filepath.Ext(filePath))
	fileName := filepath.Base(filePath)

	switch ext {
	case ".json":
		return cv.jsonValidator.ValidateJSONFile(filePath)
	case ".yaml", ".yml":
		return cv.yamlValidator.ValidateYAMLFile(filePath)
	case ".env":
		return cv.envValidator.ValidateEnvFile(filePath)
	default:
		// Check for specific configuration files without extensions
		switch fileName {
		case "Dockerfile":
			return cv.dockerValidator.ValidateDockerfile(filePath)
		case "Makefile":
			return cv.makefileValidator.ValidateMakefile(filePath)
		case ".gitignore":
			return cv.validateGitignore(filePath)
		case ".dockerignore":
			return cv.validateDockerignore(filePath)
		default:
			// Unknown configuration file type - return basic validation
			return cv.validateUnknownFile(filePath, fileName)
		}
	}
}

// validateUnknownFile validates files with unknown extensions
func (cv *ConfigValidator) validateUnknownFile(filePath, fileName string) (*interfaces.ConfigValidationResult, error) {
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
