package interactive

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validator defines a function that validates input
type Validator func(input string) error

// ValidateNotEmpty validates that input is not empty
func ValidateNotEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// ValidateProjectName validates a project name
func ValidateProjectName(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	// Project name should be alphanumeric with hyphens and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, input)
	if !matched {
		return fmt.Errorf("project name must contain only letters, numbers, hyphens, and underscores")
	}

	return nil
}

// ValidatePort validates a port number
func ValidatePort(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	port, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("port must be a number")
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}

// ValidateGoModule validates a Go module path
func ValidateGoModule(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	// Basic validation for Go module path
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._/-]+$`, input)
	if !matched {
		return fmt.Errorf("invalid Go module path format")
	}

	return nil
}

// ValidatePackageName validates a Java/Android package name
func ValidatePackageName(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	// Package name should be in format: com.example.app
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$`, input)
	if !matched {
		return fmt.Errorf("package name must be in format: com.example.app (lowercase, dot-separated)")
	}

	return nil
}

// ValidateBundleID validates an iOS bundle identifier
func ValidateBundleID(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	// Bundle ID should be in format: com.example.App
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9]*(\.[a-zA-Z][a-zA-Z0-9]*)+$`, input)
	if !matched {
		return fmt.Errorf("bundle ID must be in format: com.example.App (dot-separated)")
	}

	return nil
}

// ValidateAPILevel validates an Android API level
func ValidateAPILevel(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	level, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("API level must be a number")
	}

	if level < 21 || level > 36 {
		return fmt.Errorf("API level must be between 21 and 36")
	}

	return nil
}

// ValidateIOSVersion validates an iOS deployment target version
func ValidateIOSVersion(input string) error {
	if err := ValidateNotEmpty(input); err != nil {
		return err
	}

	// Basic validation for iOS version format (e.g., 15.0, 16.4)
	matched, _ := regexp.MatchString(`^\d+\.\d+$`, input)
	if !matched {
		return fmt.Errorf("iOS version must be in format: X.Y (e.g., 15.0)")
	}

	return nil
}

// InputWithValidation prompts for input with validation
func InputWithValidation(prompter PrompterInterface, message string, defaultValue string, validator Validator) (string, error) {
	for {
		input, err := prompter.Input(message, defaultValue)
		if err != nil {
			return "", err
		}

		if validator != nil {
			if err := validator(input); err != nil {
				fmt.Printf("  ✗ %v. Please try again.\n", err)
				continue
			}
		}

		return input, nil
	}
}
