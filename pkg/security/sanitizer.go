// Package security provides input sanitization and validation for the CLI tool.
package security

import (
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides input sanitization functionality.
type Sanitizer struct {
	maxStringLength int
	allowedFileExts []string
	blockedPatterns []*regexp.Regexp
}

// NewSanitizer creates a new sanitizer with secure defaults.
func NewSanitizer() *Sanitizer {
	// Compile dangerous patterns once for performance
	dangerousPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)data:`),
		regexp.MustCompile(`(?i)vbscript:`),
		regexp.MustCompile(`(?i)on\w+\s*=`),
		regexp.MustCompile(`\.\./`),
		regexp.MustCompile(`\.\.\\`),
		regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter)\s`),
		regexp.MustCompile(`['";]`),
		regexp.MustCompile(`(?i)eval\s*\(`),
		regexp.MustCompile(`(?i)exec\s*\(`),
		regexp.MustCompile(`(?i)system\s*\(`),
	}

	return &Sanitizer{
		maxStringLength: 10000,
		allowedFileExts: []string{".go", ".js", ".ts", ".jsx", ".tsx", ".json", ".yaml", ".yml", ".md", ".txt", ".html", ".css", ".scss", ".sql"},
		blockedPatterns: dangerousPatterns,
	}
}

// SanitizeProjectName sanitizes and validates a project name.
// Returns the sanitized name and an error if validation fails.
func (s *Sanitizer) SanitizeProjectName(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}

	// Check length
	if len(name) > 100 {
		return "", fmt.Errorf("project name exceeds maximum length of 100 characters")
	}

	// Check for dangerous patterns
	for _, pattern := range s.blockedPatterns {
		if pattern.MatchString(name) {
			return "", fmt.Errorf("project name contains potentially dangerous content")
		}
	}

	// Convert to lowercase and replace spaces with hyphens
	sanitized := strings.ToLower(strings.TrimSpace(name))
	sanitized = regexp.MustCompile(`\s+`).ReplaceAllString(sanitized, "-")

	// Remove invalid characters for project names
	sanitized = regexp.MustCompile(`[^a-z0-9\-_]`).ReplaceAllString(sanitized, "")

	// Remove multiple consecutive hyphens/underscores
	sanitized = regexp.MustCompile(`[-_]+`).ReplaceAllString(sanitized, "-")

	// Trim leading/trailing hyphens/underscores
	sanitized = strings.Trim(sanitized, "-_")

	if sanitized == "" {
		return "", fmt.Errorf("project name cannot be empty after sanitization")
	}

	// Check for reserved names
	reservedNames := []string{
		"con", "prn", "aux", "nul",
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
		"node_modules", "vendor", "build", "dist", "tmp", "temp",
	}
	for _, reserved := range reservedNames {
		if sanitized == reserved {
			return "", fmt.Errorf("project name '%s' is reserved and cannot be used", sanitized)
		}
	}

	// Must start with letter or number
	if len(sanitized) > 0 && !unicode.IsLetter(rune(sanitized[0])) && !unicode.IsDigit(rune(sanitized[0])) {
		return "", fmt.Errorf("project name must start with a letter or number")
	}

	return sanitized, nil
}

// SanitizePath sanitizes and validates a file or directory path.
// Returns the sanitized path and an error if validation fails.
func (s *Sanitizer) SanitizePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path
	cleaned := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path contains path traversal attempts")
	}

	// Check for null bytes (security risk)
	if strings.Contains(cleaned, "\x00") {
		return "", fmt.Errorf("path contains null bytes")
	}

	// Check for invalid characters in path
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(cleaned, char) {
			return "", fmt.Errorf("path contains invalid character '%s'", char)
		}
	}

	// Check path length (Windows MAX_PATH limit)
	if len(cleaned) > 260 {
		return "", fmt.Errorf("path exceeds maximum path length of 260 characters")
	}

	// Check for reserved names on Windows
	baseName := filepath.Base(cleaned)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	upperBaseName := strings.ToUpper(baseName)
	for _, reserved := range reservedNames {
		if upperBaseName == reserved || strings.HasPrefix(upperBaseName, reserved+".") {
			return "", fmt.Errorf("path uses reserved name '%s'", reserved)
		}
	}

	return cleaned, nil
}

// SanitizeString performs general string sanitization.
// Returns the sanitized string and an error if validation fails.
func (s *Sanitizer) SanitizeString(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", nil
	}

	// Check length limits
	if len(input) > s.maxStringLength {
		return "", fmt.Errorf("string exceeds maximum length of %d characters", s.maxStringLength)
	}

	// Check for dangerous patterns
	for _, pattern := range s.blockedPatterns {
		if pattern.MatchString(input) {
			return "", fmt.Errorf("string contains potentially dangerous content")
		}
	}

	// Escape HTML
	sanitized := html.EscapeString(input)

	// Normalize whitespace
	sanitized = strings.TrimSpace(sanitized)
	sanitized = regexp.MustCompile(`\s+`).ReplaceAllString(sanitized, " ")

	// Check for control characters
	for _, r := range sanitized {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return "", fmt.Errorf("string contains invalid control characters")
		}
	}

	return sanitized, nil
}

// SanitizeConfigValue sanitizes configuration values.
// Returns the sanitized value and an error if validation fails.
func (s *Sanitizer) SanitizeConfigValue(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}

	// Check length
	if len(value) > 1000 {
		return "", fmt.Errorf("configuration value exceeds maximum length of 1000 characters")
	}

	// Check for dangerous patterns
	for _, pattern := range s.blockedPatterns {
		if pattern.MatchString(value) {
			return "", fmt.Errorf("configuration value contains potentially dangerous content")
		}
	}

	// Trim whitespace
	sanitized := strings.TrimSpace(value)

	// Check for control characters
	for _, r := range sanitized {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return "", fmt.Errorf("configuration value contains invalid control characters")
		}
	}

	return sanitized, nil
}

// SanitizeMap sanitizes all string values in a map.
// Returns a map of field names to errors for any validation failures.
func (s *Sanitizer) SanitizeMap(data map[string]interface{}) (map[string]interface{}, map[string]error) {
	sanitized := make(map[string]interface{})
	errors := make(map[string]error)

	for key, value := range data {
		switch v := value.(type) {
		case string:
			sanitizedValue, err := s.SanitizeConfigValue(v)
			if err != nil {
				errors[key] = err
			} else {
				sanitized[key] = sanitizedValue
			}
		case map[string]interface{}:
			// Recursively sanitize nested maps
			nestedSanitized, nestedErrors := s.SanitizeMap(v)
			sanitized[key] = nestedSanitized
			for nestedKey, nestedErr := range nestedErrors {
				errors[key+"."+nestedKey] = nestedErr
			}
		case []interface{}:
			// Sanitize slice elements if they are strings
			sanitizedSlice := make([]interface{}, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					sanitizedValue, err := s.SanitizeConfigValue(str)
					if err != nil {
						errors[fmt.Sprintf("%s[%d]", key, i)] = err
					} else {
						sanitizedSlice[i] = sanitizedValue
					}
				} else {
					sanitizedSlice[i] = item
				}
			}
			sanitized[key] = sanitizedSlice
		default:
			sanitized[key] = value
		}
	}

	return sanitized, errors
}
