package config

import (
	"fmt"
	"strconv"
	"strings"
)

// IOSConfigValidator validates iOS component configurations
type IOSConfigValidator struct{}

// NewIOSConfigValidator creates a new iOS config validator
func NewIOSConfigValidator() *IOSConfigValidator {
	return &IOSConfigValidator{}
}

// Validate validates the iOS configuration
func (v *IOSConfigValidator) Validate(config map[string]interface{}) error {
	// Validate name
	if name, exists := config["name"]; exists {
		if err := validateProjectName(name); err != nil {
			return NewFieldError("name", err.Error())
		}
	}

	// Validate bundle_id (required)
	if bundleID, exists := config["bundle_id"]; exists {
		if err := v.ValidateBundleID(bundleID); err != nil {
			return NewFieldError("bundle_id", err.Error())
		}
	}

	// Validate deployment_target
	if deploymentTarget, exists := config["deployment_target"]; exists {
		if err := v.ValidateDeploymentTarget(deploymentTarget); err != nil {
			return NewFieldError("deployment_target", err.Error())
		}
	}

	// Validate language
	if language, exists := config["language"]; exists {
		if err := v.ValidateLanguage(language); err != nil {
			return NewFieldError("language", err.Error())
		}
	}

	return nil
}

// GetRequiredFields returns required configuration fields
func (v *IOSConfigValidator) GetRequiredFields() []string {
	return []string{"name", "bundle_id"}
}

// GetOptionalFields returns optional configuration fields
func (v *IOSConfigValidator) GetOptionalFields() []string {
	return []string{"deployment_target", "language"}
}

// GetFieldDescription returns description for a field
func (v *IOSConfigValidator) GetFieldDescription(field string) string {
	descriptions := map[string]string{
		"name":              "The name of the iOS project",
		"bundle_id":         "The bundle identifier (e.g., com.example.app)",
		"deployment_target": "The minimum iOS version (e.g., 15.0)",
		"language":          "The programming language (swift or objective-c)",
	}

	if desc, exists := descriptions[field]; exists {
		return desc
	}
	return ""
}

// ValidateBundleID validates an iOS bundle identifier
func (v *IOSConfigValidator) ValidateBundleID(value interface{}) error {
	bundleID, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(bundleID) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Must contain at least one dot
	if !strings.Contains(bundleID, ".") {
		return fmt.Errorf("must be a valid bundle identifier (e.g., com.example.app)")
	}

	// Check each segment
	segments := strings.Split(bundleID, ".")
	if len(segments) < 2 {
		return fmt.Errorf("must have at least two segments (e.g., com.example)")
	}

	for _, segment := range segments {
		if len(segment) == 0 {
			return fmt.Errorf("segments cannot be empty")
		}

		// Segments must start with a letter
		if !isLetter(rune(segment[0])) {
			return fmt.Errorf("segment '%s' must start with a letter", segment)
		}

		// Segments can only contain letters, digits, and hyphens
		for _, char := range segment {
			if !isAlphanumeric(char) && char != '-' {
				return fmt.Errorf("segment '%s' contains invalid character: %c", segment, char)
			}
		}
	}

	return nil
}

// ValidateDeploymentTarget validates an iOS deployment target version
func (v *IOSConfigValidator) ValidateDeploymentTarget(value interface{}) error {
	var versionStr string

	// Handle both string and numeric types
	switch v := value.(type) {
	case string:
		versionStr = v
	case float64:
		versionStr = fmt.Sprintf("%.1f", v)
	case int:
		versionStr = fmt.Sprintf("%d.0", v)
	default:
		return fmt.Errorf("must be a string or number (e.g., '15.0' or 15.0)")
	}

	if len(versionStr) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Parse version string (e.g., "15.0", "16.4")
	parts := strings.Split(versionStr, ".")
	if len(parts) < 1 || len(parts) > 3 {
		return fmt.Errorf("must be a valid iOS version (e.g., 15.0, 16.4)")
	}

	// Validate major version
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid version format: %s", versionStr)
	}

	// iOS versions range from 1 to current (as of 2024, iOS 17 is latest)
	// We'll allow up to iOS 20 for future-proofing
	if major < 1 || major > 20 {
		return fmt.Errorf("major version must be between 1 and 20")
	}

	// Validate minor version if present
	if len(parts) > 1 {
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid version format: %s", versionStr)
		}
		if minor < 0 || minor > 99 {
			return fmt.Errorf("minor version must be between 0 and 99")
		}
	}

	// Validate patch version if present
	if len(parts) > 2 {
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid version format: %s", versionStr)
		}
		if patch < 0 || patch > 99 {
			return fmt.Errorf("patch version must be between 0 and 99")
		}
	}

	return nil
}

// ValidateLanguage validates the programming language choice
func (v *IOSConfigValidator) ValidateLanguage(value interface{}) error {
	language, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	supportedLanguages := []string{"swift", "objective-c", "objc"}
	language = strings.ToLower(language)

	for _, supported := range supportedLanguages {
		if language == supported {
			return nil
		}
	}

	return fmt.Errorf("must be one of: swift, objective-c")
}
