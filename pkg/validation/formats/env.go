package formats

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// Pre-compiled regular expressions for environment validation
var (
	envKeyStartRegex    = regexp.MustCompile(`^[A-Za-z_]`)
	envKeyValidRegex    = regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	ipAddressRegex      = regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	base64PatternRegex  = regexp.MustCompile(`^[A-Za-z0-9+/]+=*$`)
	hexPatternRegex     = regexp.MustCompile(`^[a-fA-F0-9]+$`)
	numericPatternRegex = regexp.MustCompile(`^\d+$`)
)

// EnvValidator provides specialized environment file validation
type EnvValidator struct {
	secretPatterns []SecretPattern
}

// SecretPattern defines patterns for detecting potential secrets
type SecretPattern struct {
	Name        string
	Pattern     *regexp.Regexp
	Description string
}

// NewEnvValidator creates a new environment file validator
func NewEnvValidator() *EnvValidator {
	validator := &EnvValidator{}
	validator.initializeSecretPatterns()
	return validator
}

// ValidateEnvFile validates environment configuration files
func (ev *EnvValidator) ValidateEnvFile(filePath string) (*interfaces.ConfigValidationResult, error) {
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
	envVars := make(map[string]string)

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

		// Store for duplicate checking
		if existingValue, exists := envVars[key]; exists {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      key,
				Type:       "duplicate_key",
				Message:    fmt.Sprintf("Duplicate environment variable '%s' (previous value: %s)", key, existingValue),
				Suggestion: "Remove duplicate environment variable definitions",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.env.duplicate_key",
			})
			result.Summary.WarningCount++
		}
		envVars[key] = value

		// Validate key format
		if err := ev.validateEnvKey(key); err != nil {
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

		// Validate value format
		ev.validateEnvValue(key, value, lineNum+1, result)

		// Check for potential secrets
		if ev.isPotentialSecret(key, value) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum+1),
				Value:      key,
				Type:       "security",
				Message:    "Potential secret detected in environment file",
				Suggestion: "Consider using a secrets management system or .env.example file",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.env.secrets",
			})
			result.Summary.WarningCount++
		}

		// Check for hardcoded URLs and IPs
		ev.validateURLsAndIPs(key, value, lineNum+1, result)

		result.Summary.ValidProperties++
	}

	// Check for common missing environment variables
	ev.checkCommonMissingVars(envVars, result)

	return result, nil
}

// validateEnvKey validates environment variable key format
func (ev *EnvValidator) validateEnvKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	// Environment variable names should start with a letter or underscore
	if !envKeyStartRegex.MatchString(key) {
		return fmt.Errorf("key must start with a letter or underscore")
	}

	// Check for invalid characters using pre-compiled regex
	if !envKeyValidRegex.MatchString(key) {
		return fmt.Errorf("key contains invalid characters (use only letters, numbers, and underscores)")
	}

	// Recommend uppercase convention
	if strings.ToUpper(key) != key {
		return fmt.Errorf("consider using uppercase for environment variable names")
	}

	// Check for reserved names
	reservedNames := []string{"PATH", "HOME", "USER", "PWD", "SHELL"}
	for _, reserved := range reservedNames {
		if key == reserved {
			return fmt.Errorf("'%s' is a reserved environment variable name", key)
		}
	}

	return nil
}

// validateEnvValue validates environment variable value
func (ev *EnvValidator) validateEnvValue(key, value string, lineNum int, result *interfaces.ConfigValidationResult) {
	// Check for unquoted values with spaces
	if strings.Contains(value, " ") && !ev.isQuoted(value) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      value,
			Type:       "format_warning",
			Message:    "Value contains spaces but is not quoted",
			Suggestion: "Quote values containing spaces: KEY=\"value with spaces\"",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "config.env.unquoted_spaces",
		})
		result.Summary.WarningCount++
	}

	// Check for empty values
	if value == "" || value == "\"\"" || value == "''" {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      key,
			Type:       "empty_value",
			Message:    "Environment variable has empty value",
			Suggestion: "Provide a default value or remove if not needed",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "config.env.empty_value",
		})
		result.Summary.WarningCount++
	}

	// Check for boolean-like values
	if ev.isBooleanLike(value) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      value,
			Type:       "format_suggestion",
			Message:    "Boolean-like value detected",
			Suggestion: "Consider using 'true'/'false' or '1'/'0' for consistency",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "config.env.boolean_format",
		})
		result.Summary.WarningCount++
	}

	// Check for numeric values that might need quotes
	if ev.isNumericButShouldBeString(key, value) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      value,
			Type:       "format_suggestion",
			Message:    "Numeric value might need quotes if treated as string",
			Suggestion: "Quote numeric values if they should be treated as strings",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "config.env.numeric_string",
		})
		result.Summary.WarningCount++
	}
}

// validateURLsAndIPs validates URLs and IP addresses in environment values
func (ev *EnvValidator) validateURLsAndIPs(key, value string, lineNum int, result *interfaces.ConfigValidationResult) {
	// Check for localhost URLs in production-like environments
	if strings.Contains(strings.ToLower(key), "prod") || strings.Contains(strings.ToLower(key), "production") {
		if strings.Contains(value, "localhost") || strings.Contains(value, "127.0.0.1") {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", lineNum),
				Value:      value,
				Type:       "configuration_warning",
				Message:    "Localhost URL detected in production-like environment variable",
				Suggestion: "Use production URLs for production environments",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "config.env.localhost_in_prod",
			})
			result.Summary.WarningCount++
		}
	}

	// Check for HTTP URLs (should be HTTPS)
	if strings.HasPrefix(value, "http://") && !strings.Contains(value, "localhost") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      value,
			Type:       "security",
			Message:    "HTTP URL detected (not secure)",
			Suggestion: "Use HTTPS URLs for security",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "config.env.http_url",
		})
		result.Summary.WarningCount++
	}

	// Check for hardcoded IP addresses using pre-compiled regex
	if ipAddressRegex.MatchString(value) && !strings.Contains(value, "127.0.0.1") && !strings.Contains(value, "0.0.0.0") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", lineNum),
			Value:      value,
			Type:       "configuration_warning",
			Message:    "Hardcoded IP address detected",
			Suggestion: "Consider using domain names instead of IP addresses",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "config.env.hardcoded_ip",
		})
		result.Summary.WarningCount++
	}
}

// checkCommonMissingVars checks for commonly expected environment variables
func (ev *EnvValidator) checkCommonMissingVars(envVars map[string]string, result *interfaces.ConfigValidationResult) {
	commonVars := map[string]string{
		"NODE_ENV":     "Specify the Node.js environment (development, production, test)",
		"PORT":         "Specify the application port",
		"DATABASE_URL": "Database connection string",
		"LOG_LEVEL":    "Logging level (debug, info, warn, error)",
	}

	// Only suggest if we have some environment variables (not an empty file)
	if len(envVars) > 0 {
		for varName, description := range commonVars {
			if _, exists := envVars[varName]; !exists {
				// Only suggest if there are related variables
				if ev.hasRelatedVars(varName, envVars) {
					result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
						Field:      "missing_vars",
						Value:      varName,
						Type:       "missing_recommended",
						Message:    fmt.Sprintf("Consider adding %s", varName),
						Suggestion: description,
						Severity:   interfaces.ValidationSeverityInfo,
						Rule:       "config.env.missing_common",
					})
					result.Summary.WarningCount++
				}
			}
		}
	}
}

// hasRelatedVars checks if there are related environment variables
func (ev *EnvValidator) hasRelatedVars(varName string, envVars map[string]string) bool {
	relatedPatterns := map[string][]string{
		"NODE_ENV":     {"NODE_", "NPM_", "YARN_"},
		"PORT":         {"HOST", "SERVER", "APP_"},
		"DATABASE_URL": {"DB_", "DATABASE_", "POSTGRES_", "MYSQL_", "MONGO_"},
		"LOG_LEVEL":    {"LOG_", "DEBUG", "VERBOSE"},
	}

	patterns, exists := relatedPatterns[varName]
	if !exists {
		return false
	}

	for existingVar := range envVars {
		for _, pattern := range patterns {
			if strings.Contains(existingVar, pattern) {
				return true
			}
		}
	}

	return false
}

// isPotentialSecret checks if a key-value pair might contain a secret
func (ev *EnvValidator) isPotentialSecret(key, value string) bool {
	keyLower := strings.ToLower(key)

	// Check against secret patterns
	for _, pattern := range ev.secretPatterns {
		if pattern.Pattern.MatchString(keyLower) && len(value) > 8 {
			return true
		}
	}

	// Check for long random-looking strings
	if len(value) > 20 && ev.looksLikeSecret(value) {
		return true
	}

	return false
}

// looksLikeSecret checks if a value looks like a secret (random string)
func (ev *EnvValidator) looksLikeSecret(value string) bool {
	// Remove quotes
	value = strings.Trim(value, "\"'")

	// Check for base64-like patterns using pre-compiled regex
	if base64PatternRegex.MatchString(value) && len(value) > 20 {
		return true
	}

	// Check for hex patterns using pre-compiled regex
	if hexPatternRegex.MatchString(value) && len(value) > 16 {
		return true
	}

	// Check for high entropy (random-looking strings)
	return ev.hasHighEntropy(value)
}

// hasHighEntropy checks if a string has high entropy (appears random)
func (ev *EnvValidator) hasHighEntropy(s string) bool {
	if len(s) < 16 {
		return false
	}

	charCount := make(map[rune]int)
	for _, char := range s {
		charCount[char]++
	}

	// Calculate entropy
	entropy := 0.0
	length := float64(len(s))
	for _, count := range charCount {
		if count > 0 {
			freq := float64(count) / length
			entropy -= freq * (float64(count) / length)
		}
	}

	// High entropy threshold (adjust as needed)
	return entropy > 3.5
}

// isQuoted checks if a value is properly quoted
func (ev *EnvValidator) isQuoted(value string) bool {
	return (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
		(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'"))
}

// isBooleanLike checks if a value looks like a boolean
func (ev *EnvValidator) isBooleanLike(value string) bool {
	lower := strings.ToLower(strings.Trim(value, "\"'"))
	booleanValues := []string{"yes", "no", "on", "off", "enabled", "disabled", "y", "n"}

	for _, boolVal := range booleanValues {
		if lower == boolVal {
			return true
		}
	}

	return false
}

// isNumericButShouldBeString checks if a numeric value might need to be quoted
func (ev *EnvValidator) isNumericButShouldBeString(key, value string) bool {
	// Remove quotes for checking
	unquoted := strings.Trim(value, "\"'")

	// Check if it's numeric using pre-compiled regex
	if !numericPatternRegex.MatchString(unquoted) {
		return false
	}

	// Check if the key suggests it should be a string
	stringLikeKeys := []string{"version", "id", "code", "zip", "postal"}
	keyLower := strings.ToLower(key)

	for _, stringKey := range stringLikeKeys {
		if strings.Contains(keyLower, stringKey) {
			return true
		}
	}

	return false
}

// initializeSecretPatterns initializes patterns for detecting secrets
func (ev *EnvValidator) initializeSecretPatterns() {
	patterns := []struct {
		name        string
		pattern     string
		description string
	}{
		{"password", `password|passwd|pwd`, "Password fields"},
		{"secret", `secret|private`, "Secret keys"},
		{"token", `token|jwt`, "Authentication tokens"},
		{"key", `key|api`, "API keys"},
		{"auth", `auth|oauth`, "Authentication credentials"},
		{"cert", `cert|certificate|crt`, "Certificates"},
		{"private_key", `private.*key|pkey`, "Private keys"},
		{"database", `db.*pass|database.*pass`, "Database passwords"},
		{"aws", `aws.*secret|aws.*key`, "AWS credentials"},
		{"github", `github.*token|gh.*token`, "GitHub tokens"},
	}

	for _, p := range patterns {
		compiled := regexp.MustCompile(p.pattern)
		ev.secretPatterns = append(ev.secretPatterns, SecretPattern{
			Name:        p.name,
			Pattern:     compiled,
			Description: p.description,
		})
	}
}
