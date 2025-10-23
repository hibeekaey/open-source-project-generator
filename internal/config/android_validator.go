package config

import (
	"fmt"
	"strings"
)

// AndroidConfigValidator validates Android component configurations
type AndroidConfigValidator struct{}

// NewAndroidConfigValidator creates a new Android config validator
func NewAndroidConfigValidator() *AndroidConfigValidator {
	return &AndroidConfigValidator{}
}

// Validate validates the Android configuration
func (v *AndroidConfigValidator) Validate(config map[string]interface{}) error {
	// Validate name
	if name, exists := config["name"]; exists {
		if err := validateProjectName(name); err != nil {
			return NewFieldError("name", err.Error())
		}
	}

	// Validate package (required)
	if pkg, exists := config["package"]; exists {
		if err := v.ValidatePackage(pkg); err != nil {
			return NewFieldError("package", err.Error())
		}
	}

	// Validate min_sdk
	if minSDK, exists := config["min_sdk"]; exists {
		if err := v.ValidateMinSDK(minSDK); err != nil {
			return NewFieldError("min_sdk", err.Error())
		}
	}

	// Validate target_sdk
	if targetSDK, exists := config["target_sdk"]; exists {
		if err := v.ValidateTargetSDK(targetSDK); err != nil {
			return NewFieldError("target_sdk", err.Error())
		}
	}

	// Validate language
	if language, exists := config["language"]; exists {
		if err := v.ValidateLanguage(language); err != nil {
			return NewFieldError("language", err.Error())
		}
	}

	// Cross-field validation: min_sdk should not be greater than target_sdk
	if minSDK, minExists := config["min_sdk"]; minExists {
		if targetSDK, targetExists := config["target_sdk"]; targetExists {
			if err := v.ValidateSDKRange(minSDK, targetSDK); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetRequiredFields returns required configuration fields
func (v *AndroidConfigValidator) GetRequiredFields() []string {
	return []string{"name", "package"}
}

// GetOptionalFields returns optional configuration fields
func (v *AndroidConfigValidator) GetOptionalFields() []string {
	return []string{"min_sdk", "target_sdk", "language"}
}

// GetFieldDescription returns description for a field
func (v *AndroidConfigValidator) GetFieldDescription(field string) string {
	descriptions := map[string]string{
		"name":       "The name of the Android project",
		"package":    "The Java/Kotlin package name (e.g., com.example.app)",
		"min_sdk":    "The minimum Android SDK API level (e.g., 24)",
		"target_sdk": "The target Android SDK API level (e.g., 34)",
		"language":   "The programming language (kotlin or java)",
	}

	if desc, exists := descriptions[field]; exists {
		return desc
	}
	return ""
}

// ValidatePackage validates an Android package name
func (v *AndroidConfigValidator) ValidatePackage(value interface{}) error {
	pkg, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	if len(pkg) == 0 {
		return fmt.Errorf("cannot be empty")
	}

	// Must contain at least one dot
	if !strings.Contains(pkg, ".") {
		return fmt.Errorf("must be a valid package name (e.g., com.example.app)")
	}

	// Check each segment
	segments := strings.Split(pkg, ".")
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

		// Segments can only contain letters, digits, and underscores
		for _, char := range segment {
			if !isAlphanumeric(char) && char != '_' {
				return fmt.Errorf("segment '%s' contains invalid character: %c", segment, char)
			}
		}
	}

	return nil
}

// ValidateMinSDK validates the minimum SDK API level
func (v *AndroidConfigValidator) ValidateMinSDK(value interface{}) error {
	// Handle both int and float64 (JSON unmarshaling can produce float64)
	var minSDK int
	switch v := value.(type) {
	case int:
		minSDK = v
	case float64:
		minSDK = int(v)
	default:
		return fmt.Errorf("must be a number")
	}

	// Android API levels range from 1 to current (as of 2024, API 34 is latest)
	// We'll allow up to API 40 for future-proofing
	if minSDK < 1 || minSDK > 40 {
		return fmt.Errorf("must be a valid Android API level (1-40)")
	}

	// Warn if using very old API levels (< 21)
	if minSDK < 21 {
		// This is just validation, we don't fail but could add to warnings
	}

	return nil
}

// ValidateTargetSDK validates the target SDK API level
func (v *AndroidConfigValidator) ValidateTargetSDK(value interface{}) error {
	// Handle both int and float64 (JSON unmarshaling can produce float64)
	var targetSDK int
	switch v := value.(type) {
	case int:
		targetSDK = v
	case float64:
		targetSDK = int(v)
	default:
		return fmt.Errorf("must be a number")
	}

	// Android API levels range from 1 to current (as of 2024, API 34 is latest)
	// We'll allow up to API 40 for future-proofing
	if targetSDK < 1 || targetSDK > 40 {
		return fmt.Errorf("must be a valid Android API level (1-40)")
	}

	return nil
}

// ValidateLanguage validates the programming language choice
func (v *AndroidConfigValidator) ValidateLanguage(value interface{}) error {
	language, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	supportedLanguages := []string{"kotlin", "java"}
	language = strings.ToLower(language)

	for _, supported := range supportedLanguages {
		if language == supported {
			return nil
		}
	}

	return fmt.Errorf("must be one of: %s", strings.Join(supportedLanguages, ", "))
}

// ValidateSDKRange validates that min_sdk is not greater than target_sdk
func (v *AndroidConfigValidator) ValidateSDKRange(minSDKValue, targetSDKValue interface{}) error {
	var minSDK, targetSDK int

	// Convert minSDK
	switch v := minSDKValue.(type) {
	case int:
		minSDK = v
	case float64:
		minSDK = int(v)
	default:
		return NewFieldError("min_sdk", "must be a number")
	}

	// Convert targetSDK
	switch v := targetSDKValue.(type) {
	case int:
		targetSDK = v
	case float64:
		targetSDK = int(v)
	default:
		return NewFieldError("target_sdk", "must be a number")
	}

	if minSDK > targetSDK {
		return NewFieldError("min_sdk", fmt.Sprintf("cannot be greater than target_sdk (%d > %d)", minSDK, targetSDK))
	}

	return nil
}
