package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// VersionValidator handles Node.js version validation logic
type VersionValidator struct {
	compatibilityMatrix *models.VersionCompatibilityMatrix
}

// NewVersionValidator creates a new version validator
func NewVersionValidator() *VersionValidator {
	return &VersionValidator{
		compatibilityMatrix: getDefaultCompatibilityMatrix(),
	}
}

// ValidateNodeVersionConfig validates a Node.js version configuration
func (v *VersionValidator) ValidateNodeVersionConfig(config *models.NodeVersionConfig) *models.VersionValidationResult {
	result := &models.VersionValidationResult{
		Valid:       true,
		Errors:      []models.VersionValidationError{},
		Warnings:    []models.VersionValidationWarning{},
		Suggestions: []models.VersionSuggestion{},
		ValidatedAt: time.Now(),
	}

	// Validate runtime version format
	if err := v.validateVersionFormat(config.Runtime, "runtime"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "runtime",
			Value:    config.Runtime,
			Message:  err.Error(),
			Severity: "error",
			Code:     "INVALID_VERSION_FORMAT",
		})
	}

	// Validate types package version format
	if err := v.validateVersionFormat(config.TypesPackage, "types"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "types",
			Value:    config.TypesPackage,
			Message:  err.Error(),
			Severity: "error",
			Code:     "INVALID_VERSION_FORMAT",
		})
	}

	// Validate NPM version format
	if err := v.validateVersionFormat(config.NPMVersion, "npm"); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "npm",
			Value:    config.NPMVersion,
			Message:  err.Error(),
			Severity: "error",
			Code:     "INVALID_VERSION_FORMAT",
		})
	}

	// Validate Docker image format
	if err := v.validateDockerImageFormat(config.DockerImage); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "docker",
			Value:    config.DockerImage,
			Message:  err.Error(),
			Severity: "error",
			Code:     "INVALID_DOCKER_IMAGE",
		})
	}

	// Validate compatibility between runtime and types versions
	if compatErr := v.validateRuntimeTypesCompatibility(config.Runtime, config.TypesPackage); compatErr != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "compatibility",
			Value:    fmt.Sprintf("runtime: %s, types: %s", config.Runtime, config.TypesPackage),
			Message:  compatErr.Error(),
			Severity: "critical",
			Code:     "VERSION_COMPATIBILITY_MISMATCH",
		})
	}

	// Add suggestions for improvements
	suggestions := v.generateVersionSuggestions(config)
	result.Suggestions = append(result.Suggestions, suggestions...)

	return result
}

// ValidateVersionCompatibility validates compatibility between multiple version configurations
func (v *VersionValidator) ValidateVersionCompatibility(configs []*models.NodeVersionConfig) *models.VersionValidationResult {
	result := &models.VersionValidationResult{
		Valid:       true,
		Errors:      []models.VersionValidationError{},
		Warnings:    []models.VersionValidationWarning{},
		Suggestions: []models.VersionSuggestion{},
		ValidatedAt: time.Now(),
	}

	if len(configs) < 2 {
		return result
	}

	// Check for consistency across configurations
	baseConfig := configs[0]
	for i, config := range configs[1:] {
		if config.Runtime != baseConfig.Runtime {
			result.Warnings = append(result.Warnings, models.VersionValidationWarning{
				Field:   fmt.Sprintf("config[%d].runtime", i+1),
				Value:   config.Runtime,
				Message: fmt.Sprintf("Runtime version mismatch: expected %s, got %s", baseConfig.Runtime, config.Runtime),
				Code:    "RUNTIME_VERSION_MISMATCH",
			})
		}

		if config.TypesPackage != baseConfig.TypesPackage {
			result.Warnings = append(result.Warnings, models.VersionValidationWarning{
				Field:   fmt.Sprintf("config[%d].types", i+1),
				Value:   config.TypesPackage,
				Message: fmt.Sprintf("Types package version mismatch: expected %s, got %s", baseConfig.TypesPackage, config.TypesPackage),
				Code:    "TYPES_VERSION_MISMATCH",
			})
		}
	}

	return result
}

// ValidateAgainstLTS validates if the configuration uses LTS versions
func (v *VersionValidator) ValidateAgainstLTS(config *models.NodeVersionConfig) *models.VersionValidationResult {
	result := &models.VersionValidationResult{
		Valid:       true,
		Errors:      []models.VersionValidationError{},
		Warnings:    []models.VersionValidationWarning{},
		Suggestions: []models.VersionSuggestion{},
		ValidatedAt: time.Now(),
	}

	// Extract major version from runtime requirement
	majorVersion, err := v.extractMajorVersion(config.Runtime)
	if err != nil {
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "runtime",
			Value:    config.Runtime,
			Message:  fmt.Sprintf("Cannot extract major version: %s", err.Error()),
			Severity: "error",
			Code:     "VERSION_PARSE_ERROR",
		})
		result.Valid = false
		return result
	}

	// Check if it's an LTS version (even major versions are LTS)
	if majorVersion%2 != 0 {
		result.Warnings = append(result.Warnings, models.VersionValidationWarning{
			Field:   "runtime",
			Value:   config.Runtime,
			Message: fmt.Sprintf("Node.js %d is not an LTS version. Consider using an LTS version for production stability.", majorVersion),
			Code:    "NON_LTS_VERSION",
		})

		// Suggest the nearest LTS version
		suggestedLTS := majorVersion - 1
		if majorVersion < 18 {
			suggestedLTS = 18 // Minimum supported LTS
		}

		result.Suggestions = append(result.Suggestions, models.VersionSuggestion{
			Field:          "runtime",
			CurrentValue:   config.Runtime,
			SuggestedValue: fmt.Sprintf(">=%d.0.0", suggestedLTS),
			Reason:         fmt.Sprintf("Node.js %d is the nearest LTS version with long-term support", suggestedLTS),
			Priority:       "medium",
			BreakingChange: false,
		})
	}

	return result
}

// validateVersionFormat validates the format of version strings
func (v *VersionValidator) validateVersionFormat(version, field string) error {
	if version == "" {
		return fmt.Errorf("%s version cannot be empty", field)
	}

	// Define regex patterns for different version formats
	patterns := map[string]*regexp.Regexp{
		"runtime": regexp.MustCompile(`^(>=|>|<=|<|~|\^)?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9\-\.]+)?(\+[a-zA-Z0-9\-\.]+)?$`),
		"types":   regexp.MustCompile(`^(\^|~)?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9\-\.]+)?$`),
		"npm":     regexp.MustCompile(`^(>=|>|<=|<|~|\^)?(\d+)\.(\d+)\.(\d+)(-[a-zA-Z0-9\-\.]+)?$`),
	}

	pattern, exists := patterns[field]
	if !exists {
		pattern = patterns["runtime"] // Default pattern
	}

	if !pattern.MatchString(version) {
		return fmt.Errorf("invalid %s version format: %s", field, version)
	}

	return nil
}

// validateDockerImageFormat validates Docker image format
func (v *VersionValidator) validateDockerImageFormat(image string) error {
	if image == "" {
		return fmt.Errorf("docker image cannot be empty")
	}

	// Basic Docker image format validation
	// Format: [registry/]name[:tag]
	dockerPattern := regexp.MustCompile(`^([a-zA-Z0-9\-\.]+/)?[a-zA-Z0-9\-\.]+:[a-zA-Z0-9\-\.]+$`)
	if !dockerPattern.MatchString(image) {
		return fmt.Errorf("invalid Docker image format: %s", image)
	}

	// Check if it's a Node.js image
	if !strings.Contains(strings.ToLower(image), "node") {
		return fmt.Errorf("docker image should be a Node.js image: %s", image)
	}

	return nil
}

// validateRuntimeTypesCompatibility validates compatibility between runtime and types versions
func (v *VersionValidator) validateRuntimeTypesCompatibility(runtime, types string) error {
	runtimeMajor, err := v.extractMajorVersion(runtime)
	if err != nil {
		return fmt.Errorf("cannot extract runtime major version: %w", err)
	}

	typesMajor, err := v.extractMajorVersion(types)
	if err != nil {
		return fmt.Errorf("cannot extract types major version: %w", err)
	}

	// Types version should match or be close to runtime version
	// Allow types to be same major version or one version ahead
	if typesMajor < runtimeMajor || typesMajor > runtimeMajor+2 {
		return fmt.Errorf("types version %d is not compatible with runtime version %d", typesMajor, runtimeMajor)
	}

	return nil
}

// extractMajorVersion extracts the major version number from a version string
func (v *VersionValidator) extractMajorVersion(version string) (int, error) {
	// Remove version operators (>=, ^, ~, etc.)
	cleanVersion := regexp.MustCompile(`^[>=<~^]*`).ReplaceAllString(version, "")

	// Split by dots and get the first part
	parts := strings.Split(cleanVersion, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	return major, nil
}

// generateVersionSuggestions generates suggestions for version improvements
func (v *VersionValidator) generateVersionSuggestions(config *models.NodeVersionConfig) []models.VersionSuggestion {
	var suggestions []models.VersionSuggestion

	// Check if using latest LTS recommendations
	if v.compatibilityMatrix != nil {
		recommended := v.compatibilityMatrix.NodeJS

		if config.Runtime != recommended.Runtime {
			suggestions = append(suggestions, models.VersionSuggestion{
				Field:          "runtime",
				CurrentValue:   config.Runtime,
				SuggestedValue: recommended.Runtime,
				Reason:         "Use the recommended LTS version for better stability and support",
				Priority:       "medium",
				BreakingChange: false,
			})
		}

		if config.TypesPackage != recommended.TypesPackage {
			suggestions = append(suggestions, models.VersionSuggestion{
				Field:          "types",
				CurrentValue:   config.TypesPackage,
				SuggestedValue: recommended.TypesPackage,
				Reason:         "Use types version that matches the runtime version",
				Priority:       "high",
				BreakingChange: false,
			})
		}

		if config.DockerImage != recommended.DockerImage {
			suggestions = append(suggestions, models.VersionSuggestion{
				Field:          "docker",
				CurrentValue:   config.DockerImage,
				SuggestedValue: recommended.DockerImage,
				Reason:         "Use the recommended Docker image for consistency",
				Priority:       "low",
				BreakingChange: false,
			})
		}
	}

	return suggestions
}

// getDefaultCompatibilityMatrix returns the default compatibility matrix
func getDefaultCompatibilityMatrix() *models.VersionCompatibilityMatrix {
	return &models.VersionCompatibilityMatrix{
		NodeJS: models.NodeVersionConfig{
			Runtime:      ">=20.0.0",
			TypesPackage: "^20.17.0",
			NPMVersion:   ">=10.0.0",
			DockerImage:  "node:20-alpine",
			LTSStatus:    true,
			Description:  "Node.js 20 LTS - Recommended for production use",
		},
		UpdatedAt: time.Now(),
	}
}
