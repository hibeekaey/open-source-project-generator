package ui

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
)

// Re-export common validation patterns from comprehensive validation
var (
	// ProjectNamePattern validates project names (alphanumeric, hyphens, underscores)
	ProjectNamePattern = validation.NewValidator().GetPattern("project_name")

	// PackageNamePattern validates package names (lowercase, hyphens)
	PackageNamePattern = validation.NewValidator().GetPattern("package_name")

	// VersionPattern validates semantic version numbers
	VersionPattern = validation.NewValidator().GetPattern("version")

	// URLPattern validates HTTP/HTTPS URLs
	URLPattern = validation.NewValidator().GetPattern("url")

	// GitHubRepoPattern validates GitHub repository names
	GitHubRepoPattern = validation.NewValidator().GetPattern("github_repo")
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

// Common Validators that delegate to comprehensive validation

// RequiredValidator validates that input is not empty
func RequiredValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			validator := validation.NewValidator()
			return validator.ValidateNonEmptyString(input, fieldName)
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
			validator := validation.NewValidator()
			issues := validator.ValidateString(fieldName, input, min, max)
			if len(issues) > 0 {
				return fmt.Errorf("%s", issues[0].Message)
			}
			return nil
		},
		Message: fmt.Sprintf("%s must be between %d and %d characters", fieldName, min, max),
		Suggestions: []string{
			fmt.Sprintf("Enter a value between %d and %d characters", min, max),
		},
	}
}

// EmailValidator validates email format
func EmailValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			validator := validation.NewValidator()
			return validator.ValidateEmail(input)
		},
		Message: fmt.Sprintf("%s must be a valid email address", fieldName),
		Suggestions: []string{
			"example@domain.com",
			"user@company.org",
		},
	}
}

// URLValidator validates URL format
func URLValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			if input == "" {
				return nil // Empty is allowed
			}
			// More strict URL validation
			if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
				return fmt.Errorf("%s must be a valid URL", fieldName)
			}
			// Check for incomplete URLs
			if input == "https://" || input == "http://" {
				return fmt.Errorf("%s must be a valid URL", fieldName)
			}
			// Check for invalid protocols
			if strings.HasPrefix(input, "ftp://") {
				return fmt.Errorf("%s must be a valid URL", fieldName)
			}
			return nil
		},
		Message: fmt.Sprintf("%s must be a valid URL", fieldName),
		Suggestions: []string{
			"https://example.com",
			"https://github.com/user/repo",
		},
	}
}

// ProjectNameValidator validates project name format
func ProjectNameValidator(fieldName ...string) ValidationRule {
	name := "project_name"
	if len(fieldName) > 0 {
		name = fieldName[0]
	}
	return ValidationRule{
		Name: name,
		Validator: func(input string) error {
			validator := validation.NewValidator()
			issues := validator.ValidateString(name, input, 1, 100, "required", "project_name")
			if len(issues) > 0 {
				return fmt.Errorf("%s", issues[0].Message)
			}
			return nil
		},
		Message: fmt.Sprintf("%s must contain only alphanumeric characters, hyphens, and underscores", name),
		Suggestions: []string{
			"my-awesome-project",
			"project_name",
			"awesome-project-123",
		},
	}
}

// PackageNameValidator validates package name format
func PackageNameValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			validator := validation.NewValidator()
			issues := validator.ValidateString(fieldName, input, 1, 100, "required", "package_name")
			if len(issues) > 0 {
				return fmt.Errorf("%s", issues[0].Message)
			}
			return nil
		},
		Message: fmt.Sprintf("%s must be lowercase with hyphens only", fieldName),
		Suggestions: []string{
			"my-package",
			"awesome-package",
			"package-name",
		},
	}
}

// VersionValidator validates semantic version format
func VersionValidator(fieldName string) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			validator := validation.NewValidator()
			return validator.ValidateVersion(input)
		},
		Message: fmt.Sprintf("%s must follow semantic versioning", fieldName),
		Suggestions: []string{
			"1.0.0",
			"v1.2.3",
			"2.0.0-beta.1",
		},
	}
}

// ValidateAndSuggest validates input and provides suggestions
func ValidateAndSuggest(input string, chain *ValidationChain) (string, []string, error) {
	if err := chain.Validate(input); err != nil {
		// Extract suggestions from the error if available
		var suggestions []string
		if validationErr, ok := err.(*interfaces.ValidationError); ok {
			suggestions = validationErr.Suggestions
		}
		return input, suggestions, err
	}
	return input, nil, nil
}

// CreateValidationChain creates a validation chain for common use cases
func CreateValidationChain(fieldName string, required bool, minLength, maxLength int, fieldType string) *ValidationChain {
	chain := NewValidationChain()

	if required {
		chain.Add(RequiredValidator(fieldName))
	}

	if minLength > 0 || maxLength > 0 {
		chain.Add(LengthValidator(fieldName, minLength, maxLength))
	}

	switch fieldType {
	case "email":
		chain.Add(EmailValidator(fieldName))
	case "url":
		chain.Add(URLValidator(fieldName))
	case "project_name":
		chain.Add(ProjectNameValidator(fieldName))
	case "package_name":
		chain.Add(PackageNameValidator(fieldName))
	case "version":
		chain.Add(VersionValidator(fieldName))
	}

	return chain
}

// ProjectConfigValidation creates a validation chain for project configuration
func ProjectConfigValidation() *ValidationChain {
	return NewValidationChain().
		Add(RequiredValidator("project_name")).
		Add(ProjectNameValidator("project_name"))
}

// EmailConfigValidation creates a validation chain for email configuration
func EmailConfigValidation(required ...bool) *ValidationChain {
	chain := NewValidationChain()
	if len(required) == 0 || required[0] {
		chain.Add(RequiredValidator("email"))
	}
	chain.Add(EmailValidator("email"))
	return chain
}

// URLConfigValidation creates a validation chain for URL configuration
func URLConfigValidation(required ...bool) *ValidationChain {
	chain := NewValidationChain()
	if len(required) == 0 || required[0] {
		chain.Add(RequiredValidator("url"))
	}
	chain.Add(URLValidator("url"))
	return chain
}

// VersionConfigValidation creates a validation chain for version configuration
func VersionConfigValidation(required ...bool) *ValidationChain {
	chain := NewValidationChain()
	if len(required) == 0 || required[0] {
		chain.Add(RequiredValidator("version"))
	}
	chain.Add(VersionValidator("version"))
	return chain
}

// NumericValidator validates numeric input
func NumericValidator(fieldName string, min, max int) ValidationRule {
	return ValidationRule{
		Name: fieldName,
		Validator: func(input string) error {
			if input == "" {
				return nil // Empty is allowed
			}
			// Check if input is numeric (allow decimal points)
			hasDecimal := false
			for i, char := range input {
				if char == '.' {
					if hasDecimal || i == 0 || i == len(input)-1 {
						return fmt.Errorf("%s must be numeric", fieldName)
					}
					hasDecimal = true
				} else if char < '0' || char > '9' {
					return fmt.Errorf("%s must be numeric", fieldName)
				}
			}
			// Convert to float for range checking
			var value float64
			if hasDecimal {
				// Parse decimal number
				parts := strings.Split(input, ".")
				if len(parts) != 2 {
					return fmt.Errorf("%s must be numeric", fieldName)
				}
				// Basic parsing for decimal
				value = 0
				for _, char := range parts[0] {
					value = value*10 + float64(char-'0')
				}
				decimal := 0.0
				multiplier := 0.1
				for _, char := range parts[1] {
					decimal += float64(char-'0') * multiplier
					multiplier *= 0.1
				}
				value += decimal
			} else {
				// Parse integer
				value = 0
				for _, char := range input {
					value = value*10 + float64(char-'0')
				}
			}
			if min > 0 && value < float64(min) {
				return fmt.Errorf("%s must be at least %d", fieldName, min)
			}
			if max > 0 && value > float64(max) {
				return fmt.Errorf("%s must be at most %d", fieldName, max)
			}
			return nil
		},
		Message: fmt.Sprintf("%s must be numeric", fieldName),
		Suggestions: []string{
			"123",
			"456",
		},
	}
}

// CustomValidator creates a custom validation rule
func CustomValidator(fieldName string, message string, validatorFunc ValidatorFunc, suggestions []string, recoveryOptions ...[]interfaces.RecoveryOption) ValidationRule {
	rule := ValidationRule{
		Name:        fieldName,
		Validator:   validatorFunc,
		Message:     message,
		Suggestions: suggestions,
	}
	if len(recoveryOptions) > 0 {
		rule.Recovery = recoveryOptions[0]
	}
	return rule
}

// SanitizeInput sanitizes input by trimming whitespace and normalizing
func SanitizeInput(input string) string {
	// Trim whitespace
	result := strings.TrimSpace(input)

	// Normalize line endings (convert \r\n to \n)
	result = strings.ReplaceAll(result, "\r\n", "\n")

	// Convert \r to \n
	result = strings.ReplaceAll(result, "\r", "\n")

	// Remove null characters
	result = strings.ReplaceAll(result, "\x00", "")

	// Remove escape characters
	result = strings.ReplaceAll(result, "\x1b", "")

	return result
}

// CreateDefaultValueRecovery creates a recovery option with default value
func CreateDefaultValueRecovery(defaultValue string, fieldName ...string) interfaces.RecoveryOption {
	name := "field"
	if len(fieldName) > 0 {
		name = fieldName[0]
	}
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Use default: %s", defaultValue),
		Description: fmt.Sprintf("Use default value '%s' for %s", defaultValue, name),
		Action:      func() error { return nil }, // No-op action
		Safe:        true,
		Metadata: map[string]interface{}{
			"default_value": defaultValue,
			"field_name":    name,
		},
	}
}

// SuggestCorrection suggests a correction for input
func SuggestCorrection(input string, suggestion string) string {
	if suggestion == "" {
		return input
	}

	// Apply various corrections based on common patterns
	result := input

	// Convert spaces to hyphens
	result = strings.ReplaceAll(result, " ", "-")

	// Remove special characters except hyphens and underscores
	var cleaned strings.Builder
	for _, char := range result {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_' {
			cleaned.WriteRune(char)
		}
	}
	result = cleaned.String()

	// Convert to lowercase
	result = strings.ToLower(result)

	// Remove multiple consecutive hyphens
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}

	// Remove leading/trailing hyphens
	result = strings.Trim(result, "-")

	return result
}

// CreateSuggestionRecovery creates a recovery option with suggestions
func CreateSuggestionRecovery(suggestion string) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Apply suggestion: %s", suggestion),
		Description: "Use the suggested value",
		Action:      func() error { return nil },
		Safe:        true,
		Metadata: map[string]interface{}{
			"suggestion": suggestion,
		},
	}
}

// CreateSkipFieldRecovery creates a recovery option to skip a field
func CreateSkipFieldRecovery(fieldName string) interfaces.RecoveryOption {
	return interfaces.RecoveryOption{
		Label:       fmt.Sprintf("Skip %s", fieldName),
		Description: fmt.Sprintf("Skip the %s field", fieldName),
		Action:      func() error { return nil },
		Safe:        true,
		Metadata: map[string]interface{}{
			"field_name": fieldName,
			"action":     "skip",
		},
	}
}
