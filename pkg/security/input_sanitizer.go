// Package security provides comprehensive security validation and sanitization
// for all user inputs and template processing operations.
package security

import (
	"fmt"
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// InputSanitizer provides comprehensive input sanitization and validation
type InputSanitizer struct {
	// Configuration for sanitization behavior
	allowHTML         bool
	maxStringLength   int
	allowedFileExts   []string
	blockedPatterns   []*regexp.Regexp
	allowedURLSchemes []string
}

// NewInputSanitizer creates a new input sanitizer with secure defaults
func NewInputSanitizer() *InputSanitizer {
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

	return &InputSanitizer{
		allowHTML:         false,
		maxStringLength:   10000,
		allowedFileExts:   []string{".go", ".js", ".ts", ".jsx", ".tsx", ".json", ".yaml", ".yml", ".md", ".txt", ".html", ".css", ".scss", ".sql"},
		blockedPatterns:   dangerousPatterns,
		allowedURLSchemes: []string{"http", "https", "ftp", "ftps"},
	}
}

// SanitizationResult contains the result of input sanitization
type SanitizationResult struct {
	Original    string   `json:"original"`
	Sanitized   string   `json:"sanitized"`
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	WasModified bool     `json:"was_modified"`
}

// SanitizeString performs comprehensive string sanitization
func (s *InputSanitizer) SanitizeString(input string, fieldName string) *SanitizationResult {
	result := &SanitizationResult{
		Original:  input,
		Sanitized: input,
		IsValid:   true,
		Errors:    []string{},
		Warnings:  []string{},
	}

	// Check for empty input
	if strings.TrimSpace(input) == "" {
		return result
	}

	// Check length limits
	if len(input) > s.maxStringLength {
		result.Errors = append(result.Errors, fmt.Sprintf("%s exceeds maximum length of %d characters", fieldName, s.maxStringLength))
		result.IsValid = false
		return result
	}

	// Check for dangerous patterns
	for _, pattern := range s.blockedPatterns {
		if pattern.MatchString(input) {
			result.Errors = append(result.Errors, fmt.Sprintf("%s contains potentially dangerous content", fieldName))
			result.IsValid = false
			return result
		}
	}

	// Sanitize HTML if not allowed
	if !s.allowHTML {
		sanitized := html.EscapeString(input)
		if sanitized != input {
			result.Sanitized = sanitized
			result.WasModified = true
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s contained HTML characters that were escaped", fieldName))
		}
	}

	// Normalize whitespace
	normalized := strings.TrimSpace(result.Sanitized)
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")
	if normalized != result.Sanitized {
		result.Sanitized = normalized
		result.WasModified = true
	}

	// Check for control characters
	for _, r := range result.Sanitized {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			result.Errors = append(result.Errors, fmt.Sprintf("%s contains invalid control characters", fieldName))
			result.IsValid = false
			return result
		}
	}

	return result
}

// SanitizeProjectName performs specialized sanitization for project names
func (s *InputSanitizer) SanitizeProjectName(name string) *SanitizationResult {
	result := s.SanitizeString(name, "project_name")
	if !result.IsValid {
		return result
	}

	// Additional project name validation
	sanitized := result.Sanitized

	// Convert to lowercase and replace spaces with hyphens
	sanitized = strings.ToLower(sanitized)
	sanitized = regexp.MustCompile(`\s+`).ReplaceAllString(sanitized, "-")

	// Remove invalid characters for project names
	sanitized = regexp.MustCompile(`[^a-z0-9\-_]`).ReplaceAllString(sanitized, "")

	// Remove multiple consecutive hyphens/underscores
	sanitized = regexp.MustCompile(`[-_]+`).ReplaceAllString(sanitized, "-")

	// Trim leading/trailing hyphens/underscores
	sanitized = strings.Trim(sanitized, "-_")

	if sanitized != result.Sanitized {
		result.Sanitized = sanitized
		result.WasModified = true
		result.Warnings = append(result.Warnings, "Project name was normalized to follow naming conventions")
	}

	// Validate final result
	if sanitized == "" {
		result.Errors = append(result.Errors, "Project name cannot be empty after sanitization")
		result.IsValid = false
		return result
	}

	// Check for reserved names
	reservedNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9", "node_modules", "vendor", "build", "dist", "tmp", "temp"}
	for _, reserved := range reservedNames {
		if sanitized == reserved {
			result.Errors = append(result.Errors, fmt.Sprintf("Project name '%s' is reserved and cannot be used", sanitized))
			result.IsValid = false
			return result
		}
	}

	// Must start with letter or number
	if len(sanitized) > 0 && !unicode.IsLetter(rune(sanitized[0])) && !unicode.IsDigit(rune(sanitized[0])) {
		result.Errors = append(result.Errors, "Project name must start with a letter or number")
		result.IsValid = false
		return result
	}

	return result
}

// SanitizeFilePath performs comprehensive file path sanitization
func (s *InputSanitizer) SanitizeFilePath(path string, fieldName string) *SanitizationResult {
	result := &SanitizationResult{
		Original:  path,
		Sanitized: path,
		IsValid:   true,
		Errors:    []string{},
		Warnings:  []string{},
	}

	if strings.TrimSpace(path) == "" {
		return result
	}

	// Clean the path
	cleaned := filepath.Clean(path)
	if cleaned != path {
		result.Sanitized = cleaned
		result.WasModified = true
		result.Warnings = append(result.Warnings, fmt.Sprintf("%s was normalized", fieldName))
	}

	// Check for path traversal attempts
	if strings.Contains(cleaned, "..") {
		result.Errors = append(result.Errors, fmt.Sprintf("%s contains path traversal attempts", fieldName))
		result.IsValid = false
		return result
	}

	// Check for absolute paths where relative expected
	if strings.Contains(fieldName, "relative") && filepath.IsAbs(cleaned) {
		result.Errors = append(result.Errors, fmt.Sprintf("%s should be a relative path", fieldName))
		result.IsValid = false
		return result
	}

	// Validate file extension if it's a file path
	if strings.Contains(fieldName, "file") && filepath.Ext(cleaned) != "" {
		ext := strings.ToLower(filepath.Ext(cleaned))
		allowed := false
		for _, allowedExt := range s.allowedFileExts {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%s has extension '%s' which may not be supported", fieldName, ext))
		}
	}

	// Check path length (Windows MAX_PATH limit)
	if len(cleaned) > 260 {
		result.Errors = append(result.Errors, fmt.Sprintf("%s exceeds maximum path length of 260 characters", fieldName))
		result.IsValid = false
		return result
	}

	// Check for invalid characters in path
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(cleaned, char) {
			result.Errors = append(result.Errors, fmt.Sprintf("%s contains invalid character '%s'", fieldName, char))
			result.IsValid = false
			return result
		}
	}

	result.Sanitized = cleaned
	return result
}

// SanitizeURL performs URL sanitization and validation
func (s *InputSanitizer) SanitizeURL(urlStr string, fieldName string) *SanitizationResult {
	result := &SanitizationResult{
		Original:  urlStr,
		Sanitized: urlStr,
		IsValid:   true,
		Errors:    []string{},
		Warnings:  []string{},
	}

	if strings.TrimSpace(urlStr) == "" {
		return result
	}

	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("%s is not a valid URL: %v", fieldName, err))
		result.IsValid = false
		return result
	}

	// Check scheme
	if parsedURL.Scheme == "" {
		result.Errors = append(result.Errors, fmt.Sprintf("%s must include a scheme (http:// or https://)", fieldName))
		result.IsValid = false
		return result
	}

	// Validate allowed schemes
	schemeAllowed := false
	for _, allowedScheme := range s.allowedURLSchemes {
		if parsedURL.Scheme == allowedScheme {
			schemeAllowed = true
			break
		}
	}
	if !schemeAllowed {
		result.Errors = append(result.Errors, fmt.Sprintf("%s uses unsupported scheme '%s'", fieldName, parsedURL.Scheme))
		result.IsValid = false
		return result
	}

	// Check host
	if parsedURL.Host == "" {
		result.Errors = append(result.Errors, fmt.Sprintf("%s must include a host", fieldName))
		result.IsValid = false
		return result
	}

	// Normalize URL
	normalized := parsedURL.String()
	if normalized != urlStr {
		result.Sanitized = normalized
		result.WasModified = true
		result.Warnings = append(result.Warnings, fmt.Sprintf("%s was normalized", fieldName))
	}

	return result
}

// SanitizeEmail performs email sanitization and validation
func (s *InputSanitizer) SanitizeEmail(email string, fieldName string) *SanitizationResult {
	result := s.SanitizeString(email, fieldName)
	if !result.IsValid || strings.TrimSpace(result.Sanitized) == "" {
		return result
	}

	// Normalize email
	normalized := strings.ToLower(strings.TrimSpace(result.Sanitized))
	if normalized != result.Sanitized {
		result.Sanitized = normalized
		result.WasModified = true
		result.Warnings = append(result.Warnings, fmt.Sprintf("%s was normalized to lowercase", fieldName))
	}

	// Validate email format
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(result.Sanitized) {
		result.Errors = append(result.Errors, fmt.Sprintf("%s is not a valid email format", fieldName))
		result.IsValid = false
		return result
	}

	return result
}

// ValidateAndSanitizeMap sanitizes all string values in a map
func (s *InputSanitizer) ValidateAndSanitizeMap(data map[string]interface{}, prefix string) map[string]*SanitizationResult {
	results := make(map[string]*SanitizationResult)

	for key, value := range data {
		fieldName := key
		if prefix != "" {
			fieldName = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			results[fieldName] = s.SanitizeString(v, fieldName)
		case map[string]interface{}:
			// Recursively sanitize nested maps
			nestedResults := s.ValidateAndSanitizeMap(v, fieldName)
			for nestedKey, nestedResult := range nestedResults {
				results[nestedKey] = nestedResult
			}
		case []interface{}:
			// Sanitize slice elements if they are strings
			for i, item := range v {
				if str, ok := item.(string); ok {
					itemFieldName := fmt.Sprintf("%s[%d]", fieldName, i)
					results[itemFieldName] = s.SanitizeString(str, itemFieldName)
				}
			}
		}
	}

	return results
}

// GetSanitizationSummary returns a summary of sanitization results
func GetSanitizationSummary(results map[string]*SanitizationResult) map[string]interface{} {
	summary := map[string]interface{}{
		"total_fields":    len(results),
		"valid_fields":    0,
		"invalid_fields":  0,
		"modified_fields": 0,
		"errors":          []string{},
		"warnings":        []string{},
	}

	for fieldName, result := range results {
		if result.IsValid {
			summary["valid_fields"] = summary["valid_fields"].(int) + 1
		} else {
			summary["invalid_fields"] = summary["invalid_fields"].(int) + 1
			for _, err := range result.Errors {
				summary["errors"] = append(summary["errors"].([]string), fmt.Sprintf("%s: %s", fieldName, err))
			}
		}

		if result.WasModified {
			summary["modified_fields"] = summary["modified_fields"].(int) + 1
		}

		for _, warning := range result.Warnings {
			summary["warnings"] = append(summary["warnings"].([]string), fmt.Sprintf("%s: %s", fieldName, warning))
		}
	}

	return summary
}
