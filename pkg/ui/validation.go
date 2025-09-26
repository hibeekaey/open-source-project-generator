// Package ui provides input validation and error recovery functionality.
//
// This file contains validation utilities, common validators, and error recovery
// mechanisms for interactive UI components.
package ui

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Common validation patterns
var (
	// ProjectNamePattern validates project names (alphanumeric, hyphens, underscores)
	ProjectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$`)

	// PackageNamePattern validates package names (lowercase, hyphens)
	PackageNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)

	// VersionPattern validates semantic version numbers
	VersionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?(?:\+([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?$`)

	// URLPattern validates HTTP/HTTPS URLs
	URLPattern = regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})?(?:/[^\s]*)?$`)

	// GitHubRepoPattern validates GitHub repository names
	GitHubRepoPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
)

// ValidatorFunc defines a function type for input validation
type ValidatorFunc func(input string) error

// ValidationRule defines a validation rule with error messages and recovery options
type ValidationRule struct {
	Name        string
	Validator   ValidatorFunc
	Message     string
	Suggestions []string
	Recovery    []interfaces.RecoveryOption
}

// ValidationChain allows chaining multiple validation rules
type ValidationChain struct {
	rules []ValidationRule
}

// NewValidationChain creates a new validation chain
func NewValidationChain() *ValidationChain {
	return &ValidationChain{
		rules: make([]ValidationRule, 0),
	}
}

// Add adds a validation rule to the chain
func (vc *ValidationChain) Add(rule ValidationRule) *ValidationChain {
	vc.rules = append(vc.rules, rule)
	return vc
}

// Validate runs all validation rules in the chain
func (vc *ValidationChain) Validate(input string) error {
	for _, rule := range vc.rules {
		if err := rule.Validator(input); err != nil {
			// Create a validation error with suggestions and recovery options
			validationErr := interfaces.NewValidationError(
				rule.Name,
				input,
				rule.Message,
				"validation_failed",
			)
			if suggestErr := validationErr.WithSuggestions(rule.Suggestions...); suggestErr != nil {
				// Log error but continue
				fmt.Printf("Warning: Failed to add suggestions to validation error: %v\n", suggestErr)
			}
			if recoveryErr := validationErr.WithRecoveryOptions(rule.Recovery...); recoveryErr != nil {
				// Log error but continue
				fmt.Printf("Warning: Failed to add recovery options to validation error: %v\n", recoveryErr)
			}
			return validationErr
		}
	}
	return nil
}

// Common Validators

// RequiredValidator validates that input is not empty
func RequiredValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("field is required")
			}
			return nil
		},
		Message: fmt.Sprintf("%s is required", fieldName),
		Suggestions: []string{
			"Enter a value for this field",
			"This field cannot be left empty",
		},
	}
}

// LengthValidator validates input length
func LengthValidator(fieldName string, min, max int) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			length := len(strings.TrimSpace(input))
			if min > 0 && length < min {
				return fmt.Errorf("must be at least %d characters", min)
			}
			if max > 0 && length > max {
				return fmt.Errorf("must be at most %d characters", max)
			}
			return nil
		},
		Message: fmt.Sprintf("%s length must be between %d and %d characters", fieldName, min, max),
		Suggestions: []string{
			fmt.Sprintf("Enter between %d and %d characters", min, max),
			"Check the length of your input",
		},
	}
}

// PatternValidator validates input against a regular expression
func PatternValidator(fieldName string, pattern *regexp.Regexp, message string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			if !pattern.MatchString(strings.TrimSpace(input)) {
				return fmt.Errorf("invalid format")
			}
			return nil
		},
		Message: message,
		Suggestions: []string{
			"Check the format of your input",
			"Refer to the help for format examples",
		},
	}
}

// ProjectNameValidator validates project names
func ProjectNameValidator() ValidationRule {
	return ValidationRule{
		Name: "project_name",
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if len(input) < 2 {
				return fmt.Errorf("project name must be at least 2 characters")
			}
			if len(input) > 50 {
				return fmt.Errorf("project name must be at most 50 characters")
			}
			if !ProjectNamePattern.MatchString(input) {
				return fmt.Errorf("invalid project name format")
			}
			return nil
		},
		Message: "Project name must contain only letters, numbers, hyphens, and underscores",
		Suggestions: []string{
			"Use only letters, numbers, hyphens (-), and underscores (_)",
			"Start and end with a letter or number",
			"Examples: my-project, awesome_app, project123",
		},
		Recovery: []interfaces.RecoveryOption{
			{
				Label:       "Suggest valid name",
				Description: "Generate a valid project name based on your input",
				Safe:        true,
				Action: func() error {
					// This would be implemented to suggest a valid name
					return nil
				},
			},
		},
	}
}

// EmailValidator validates email addresses
func EmailValidator() ValidationRule {
	return ValidationRule{
		Name: "email",
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}
			_, err := mail.ParseAddress(input)
			if err != nil {
				return fmt.Errorf("invalid email format")
			}
			return nil
		},
		Message: "Please enter a valid email address",
		Suggestions: []string{
			"Use format: user@domain.com",
			"Check for typos in the email address",
			"Ensure the domain is valid",
		},
		Recovery: []interfaces.RecoveryOption{
			{
				Label:       "Use default email",
				Description: "Use a default email address",
				Safe:        true,
			},
		},
	}
}

// URLValidator validates HTTP/HTTPS URLs
func URLValidator() ValidationRule {
	return ValidationRule{
		Name: "url",
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}
			if !URLPattern.MatchString(input) {
				return fmt.Errorf("invalid URL format")
			}
			return nil
		},
		Message: "Please enter a valid HTTP or HTTPS URL",
		Suggestions: []string{
			"Use format: https://example.com",
			"Include the protocol (http:// or https://)",
			"Check for typos in the URL",
		},
		Recovery: []interfaces.RecoveryOption{
			{
				Label:       "Add https:// prefix",
				Description: "Automatically add https:// to the beginning",
				Safe:        true,
			},
		},
	}
}

// VersionValidator validates semantic version numbers
func VersionValidator() ValidationRule {
	return ValidationRule{
		Name: "version",
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}
			if !VersionPattern.MatchString(input) {
				return fmt.Errorf("invalid version format")
			}
			return nil
		},
		Message: "Please enter a valid semantic version",
		Suggestions: []string{
			"Use format: 1.0.0 or v1.0.0",
			"Include major.minor.patch numbers",
			"Examples: 1.0.0, v2.1.3, 0.1.0-alpha",
		},
		Recovery: []interfaces.RecoveryOption{
			{
				Label:       "Use default version",
				Description: "Use version 1.0.0",
				Safe:        true,
			},
		},
	}
}

// GitHubRepoValidator validates GitHub repository names
func GitHubRepoValidator() ValidationRule {
	return ValidationRule{
		Name: "github_repo",
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}
			if !GitHubRepoPattern.MatchString(input) {
				return fmt.Errorf("invalid GitHub repository format")
			}
			return nil
		},
		Message: "Please enter a valid GitHub repository name",
		Suggestions: []string{
			"Use format: username/repository",
			"Examples: octocat/Hello-World, microsoft/vscode",
			"Check the repository name on GitHub",
		},
	}
}

// NumericValidator validates numeric input
func NumericValidator(fieldName string, min, max float64) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}

			value, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("must be a valid number")
			}

			if min != 0 && value < min {
				return fmt.Errorf("must be at least %g", min)
			}
			if max != 0 && value > max {
				return fmt.Errorf("must be at most %g", max)
			}

			return nil
		},
		Message: fmt.Sprintf("%s must be a number between %g and %g", fieldName, min, max),
		Suggestions: []string{
			"Enter a valid number",
			fmt.Sprintf("Value must be between %g and %g", min, max),
		},
	}
}

// IntegerValidator validates integer input
func IntegerValidator(fieldName string, min, max int) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}

			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("must be a valid integer")
			}

			if min != 0 && value < min {
				return fmt.Errorf("must be at least %d", min)
			}
			if max != 0 && value > max {
				return fmt.Errorf("must be at most %d", max)
			}

			return nil
		},
		Message: fmt.Sprintf("%s must be an integer between %d and %d", fieldName, min, max),
		Suggestions: []string{
			"Enter a whole number",
			fmt.Sprintf("Value must be between %d and %d", min, max),
		},
	}
}

// AlphanumericValidator validates alphanumeric input
func AlphanumericValidator(fieldName string, allowSpaces bool) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			input = strings.TrimSpace(input)
			if input == "" {
				return nil // Allow empty for optional fields
			}

			for _, r := range input {
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
					if !allowSpaces || r != ' ' {
						return fmt.Errorf("must contain only letters and numbers")
					}
				}
			}
			return nil
		},
		Message: func() string {
			if allowSpaces {
				return fmt.Sprintf("%s must contain only letters, numbers, and spaces", fieldName)
			}
			return fmt.Sprintf("%s must contain only letters and numbers", fieldName)
		}(),
		Suggestions: []string{
			"Use only letters and numbers",
			"Remove special characters",
		},
	}
}

// CustomValidator creates a custom validation rule
func CustomValidator(name, message string, validator ValidatorFunc, suggestions []string, recovery []interfaces.RecoveryOption) ValidationRule {
	return ValidationRule{
		Name:        name,
		Validator:   validator,
		Message:     message,
		Suggestions: suggestions,
		Recovery:    recovery,
	}
}

// Common validation chains for different input types

// ProjectConfigValidation creates a validation chain for project configuration
func ProjectConfigValidation() *ValidationChain {
	return NewValidationChain().
		Add(RequiredValidator("project_name")).
		Add(ProjectNameValidator())
}

// EmailConfigValidation creates a validation chain for email input
func EmailConfigValidation(required bool) *ValidationChain {
	chain := NewValidationChain()
	if required {
		chain.Add(RequiredValidator("email"))
	}
	return chain.Add(EmailValidator())
}

// URLConfigValidation creates a validation chain for URL input
func URLConfigValidation(required bool) *ValidationChain {
	chain := NewValidationChain()
	if required {
		chain.Add(RequiredValidator("url"))
	}
	return chain.Add(URLValidator())
}

// VersionConfigValidation creates a validation chain for version input
func VersionConfigValidation(required bool) *ValidationChain {
	chain := NewValidationChain()
	if required {
		chain.Add(RequiredValidator("version"))
	}
	return chain.Add(VersionValidator())
}

// Utility functions for validation

// SanitizeInput sanitizes user input by trimming whitespace and normalizing
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Normalize line endings
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SuggestCorrection suggests a correction for invalid input
func SuggestCorrection(input, pattern string) string {
	// This is a simplified suggestion mechanism
	// In a real implementation, this could use more sophisticated algorithms

	input = strings.ToLower(strings.TrimSpace(input))

	// Remove invalid characters for project names
	if pattern == "project_name" {
		var result strings.Builder
		for _, r := range input {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
				result.WriteRune(r)
			} else if r == ' ' {
				result.WriteRune('-')
			}
		}
		return result.String()
	}

	return input
}

// ValidateAndSuggest validates input and provides suggestions if invalid
func ValidateAndSuggest(input string, chain *ValidationChain) (string, []string, error) {
	sanitized := SanitizeInput(input)

	err := chain.Validate(sanitized)
	if err != nil {
		suggestions := []string{}

		// Add general suggestions
		suggestions = append(suggestions, "Check the input format")
		suggestions = append(suggestions, "Refer to the help for examples")

		// Add specific suggestions based on validation error
		if validationErr, ok := err.(*interfaces.ValidationError); ok {
			suggestions = append(suggestions, validationErr.Suggestions...)
		}

		return sanitized, suggestions, err
	}

	return sanitized, nil, nil
}

// Recovery action implementations

// CreateDefaultValueRecovery creates a recovery option that sets a default value
func CreateDefaultValueRecovery(defaultValue string) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Use default: %s", defaultValue),
		Description: fmt.Sprintf("Set the value to the default: %s", defaultValue),
		Safe:        true,
		Action: func() error {
			// This would set the value to the default
			// Implementation depends on the context
			return nil
		},
	}
}

// CreateSuggestionRecovery creates a recovery option that applies a suggestion
func CreateSuggestionRecovery(suggestion string) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Apply suggestion: %s", suggestion),
		Description: fmt.Sprintf("Use the suggested value: %s", suggestion),
		Safe:        true,
		Action: func() error {
			// This would apply the suggestion
			// Implementation depends on the context
			return nil
		},
	}
}

// CreateSkipFieldRecovery creates a recovery option that skips the field
func CreateSkipFieldRecovery(fieldName string) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Skip %s", fieldName),
		Description: fmt.Sprintf("Leave %s empty and continue", fieldName),
		Safe:        true,
		Action: func() error {
			// This would skip the field
			// Implementation depends on the context
			return nil
		},
	}
}
