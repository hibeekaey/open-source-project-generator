package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// TemplateSecurityManager provides secure template processing capabilities
type TemplateSecurityManager struct {
	allowedFunctions map[string]interface{}
	blockedPatterns  []*regexp.Regexp
	maxTemplateSize  int64
	sandboxMode      bool
}

// NewTemplateSecurityManager creates a new template security manager
func NewTemplateSecurityManager() *TemplateSecurityManager {
	// Define safe template functions
	allowedFunctions := map[string]interface{}{
		// String functions
		"upper":     strings.ToUpper,
		"lower":     strings.ToLower,
		"trim":      strings.TrimSpace,
		"replace":   strings.ReplaceAll,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,

		// Safe utility functions
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"split": func(sep, str string) []string {
			return strings.Split(str, sep)
		},
	}

	// Define dangerous patterns to block
	dangerousPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i){{.*exec.*}}`),
		regexp.MustCompile(`(?i){{.*system.*}}`),
		regexp.MustCompile(`(?i){{.*cmd.*}}`),
		regexp.MustCompile(`(?i){{.*shell.*}}`),
		regexp.MustCompile(`(?i){{.*eval.*}}`),
		regexp.MustCompile(`(?i){{.*import.*}}`),
		regexp.MustCompile(`(?i){{.*require.*}}`),
		regexp.MustCompile(`(?i){{.*\.\..*}}`), // Path traversal
		regexp.MustCompile(`(?i)<script[^>]*>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)vbscript:`),
	}

	return &TemplateSecurityManager{
		allowedFunctions: allowedFunctions,
		blockedPatterns:  dangerousPatterns,
		maxTemplateSize:  1024 * 1024, // 1MB max template size
		sandboxMode:      true,
	}
}

// TemplateValidationResult contains the result of template security validation
type TemplateValidationResult struct {
	IsSecure        bool     `json:"is_secure"`
	SecurityIssues  []string `json:"security_issues"`
	Warnings        []string `json:"warnings"`
	BlockedPatterns []string `json:"blocked_patterns"`
	FilePath        string   `json:"file_path"`
	FileSize        int64    `json:"file_size"`
}

// ValidateTemplateContent validates template content for security issues
func (tsm *TemplateSecurityManager) ValidateTemplateContent(content string, filePath string) *TemplateValidationResult {
	result := &TemplateValidationResult{
		IsSecure:        true,
		SecurityIssues:  []string{},
		Warnings:        []string{},
		BlockedPatterns: []string{},
		FilePath:        filePath,
		FileSize:        int64(len(content)),
	}

	// Check template size
	if result.FileSize > tsm.maxTemplateSize {
		result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("Template size (%d bytes) exceeds maximum allowed size (%d bytes)", result.FileSize, tsm.maxTemplateSize))
		result.IsSecure = false
	}

	// Check for dangerous patterns
	for _, pattern := range tsm.blockedPatterns {
		if matches := pattern.FindAllString(content, -1); len(matches) > 0 {
			result.SecurityIssues = append(result.SecurityIssues, fmt.Sprintf("Template contains dangerous pattern: %s", pattern.String()))
			result.BlockedPatterns = append(result.BlockedPatterns, matches...)
			result.IsSecure = false
		}
	}

	// Check for potentially unsafe template actions
	unsafeActions := []string{
		"call", "html", "js", "urlquery", "printf", "print", "println",
	}

	for _, action := range unsafeActions {
		actionPattern := regexp.MustCompile(fmt.Sprintf(`{{\s*%s\s+`, action))
		if actionPattern.MatchString(content) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Template uses potentially unsafe action: %s", action))
		}
	}

	// Check for external file includes
	includePattern := regexp.MustCompile(`{{\s*template\s+"[^"]*"\s*}}`)
	if matches := includePattern.FindAllString(content, -1); len(matches) > 0 {
		result.Warnings = append(result.Warnings, "Template includes external templates - ensure they are also validated")
	}

	// Check for variable assignments that might be dangerous
	assignPattern := regexp.MustCompile(`{{\s*\$\w+\s*:=.*}}`)
	if matches := assignPattern.FindAllString(content, -1); len(matches) > 0 {
		result.Warnings = append(result.Warnings, "Template contains variable assignments - review for security implications")
	}

	return result
}

// ValidateTemplateFile validates a template file for security issues
func (tsm *TemplateSecurityManager) ValidateTemplateFile(filePath string) (*TemplateValidationResult, error) {
	// Validate file path
	cleanPath := filepath.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		return &TemplateValidationResult{
			IsSecure:       false,
			SecurityIssues: []string{"Template file path contains path traversal attempts"},
			FilePath:       filePath,
		}, fmt.Errorf("unsafe file path: %s", filePath)
	}

	// Check file size before reading
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat template file: %w", err)
	}

	if fileInfo.Size() > tsm.maxTemplateSize {
		return &TemplateValidationResult{
			IsSecure:       false,
			SecurityIssues: []string{fmt.Sprintf("Template file size (%d bytes) exceeds maximum allowed size (%d bytes)", fileInfo.Size(), tsm.maxTemplateSize)},
			FilePath:       filePath,
			FileSize:       fileInfo.Size(),
		}, fmt.Errorf("template file too large: %d bytes", fileInfo.Size())
	}

	// Read and validate content
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return tsm.ValidateTemplateContent(string(content), filePath), nil
}

// CreateSecureTemplate creates a template with security restrictions
func (tsm *TemplateSecurityManager) CreateSecureTemplate(name, content string) (*template.Template, error) {
	// Validate content first
	validation := tsm.ValidateTemplateContent(content, name)
	if !validation.IsSecure {
		return nil, fmt.Errorf("template security validation failed: %v", validation.SecurityIssues)
	}

	// Create template with restricted function map
	tmpl := template.New(name)

	// Add only allowed functions
	tmpl = tmpl.Funcs(tsm.allowedFunctions)

	// Parse template
	parsedTemplate, err := tmpl.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %w", err)
	}

	return parsedTemplate, nil
}

// ProcessTemplateSecurely processes a template with security restrictions
func (tsm *TemplateSecurityManager) ProcessTemplateSecurely(tmpl *template.Template, data interface{}, outputPath string) error {
	// Validate output path
	cleanOutputPath := filepath.Clean(outputPath)
	if strings.Contains(cleanOutputPath, "..") {
		return fmt.Errorf("unsafe output path: %s", outputPath)
	}

	// Ensure output directory exists and is safe
	outputDir := filepath.Dir(cleanOutputPath)
	if err := os.MkdirAll(outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file with secure permissions
	outputFile, err := os.OpenFile(cleanOutputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// Execute template
	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// ScanTemplateDirectory recursively scans a directory for template security issues
func (tsm *TemplateSecurityManager) ScanTemplateDirectory(dirPath string) (map[string]*TemplateValidationResult, error) {
	results := make(map[string]*TemplateValidationResult)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is a template (common extensions)
		ext := strings.ToLower(filepath.Ext(path))
		templateExts := []string{".tmpl", ".tpl", ".template", ".gotmpl"}
		isTemplate := false
		for _, templateExt := range templateExts {
			if ext == templateExt {
				isTemplate = true
				break
			}
		}

		if !isTemplate {
			return nil
		}

		// Validate template file
		result, err := tsm.ValidateTemplateFile(path)
		if err != nil {
			// Create error result
			result = &TemplateValidationResult{
				IsSecure:       false,
				SecurityIssues: []string{fmt.Sprintf("Failed to validate template: %v", err)},
				FilePath:       path,
			}
		}

		results[path] = result
		return nil
	})

	return results, err
}

// GetSecuritySummary returns a summary of template security scan results
func GetTemplateSecuritySummary(results map[string]*TemplateValidationResult) map[string]interface{} {
	summary := map[string]interface{}{
		"total_templates":         len(results),
		"secure_templates":        0,
		"insecure_templates":      0,
		"templates_with_warnings": 0,
		"total_issues":            0,
		"total_warnings":          0,
		"critical_issues":         []string{},
		"common_warnings":         []string{},
	}

	issueCount := make(map[string]int)
	warningCount := make(map[string]int)

	for filePath, result := range results {
		if result.IsSecure {
			summary["secure_templates"] = summary["secure_templates"].(int) + 1
		} else {
			summary["insecure_templates"] = summary["insecure_templates"].(int) + 1
		}

		if len(result.Warnings) > 0 {
			summary["templates_with_warnings"] = summary["templates_with_warnings"].(int) + 1
		}

		summary["total_issues"] = summary["total_issues"].(int) + len(result.SecurityIssues)
		summary["total_warnings"] = summary["total_warnings"].(int) + len(result.Warnings)

		// Track critical issues
		for _, issue := range result.SecurityIssues {
			issueCount[issue]++
			if issueCount[issue] == 1 { // First occurrence
				summary["critical_issues"] = append(summary["critical_issues"].([]string), fmt.Sprintf("%s (in %s)", issue, filePath))
			}
		}

		// Track common warnings
		for _, warning := range result.Warnings {
			warningCount[warning]++
		}
	}

	// Add most common warnings to summary
	for warning, count := range warningCount {
		if count > 1 {
			summary["common_warnings"] = append(summary["common_warnings"].([]string), fmt.Sprintf("%s (found in %d templates)", warning, count))
		}
	}

	return summary
}

// TemplateSecurityConfig allows customization of security settings
type TemplateSecurityConfig struct {
	MaxTemplateSize       int64                  `json:"max_template_size"`
	AllowedFunctions      []string               `json:"allowed_functions"`
	BlockedPatterns       []string               `json:"blocked_patterns"`
	SandboxMode           bool                   `json:"sandbox_mode"`
	AllowExternalIncludes bool                   `json:"allow_external_includes"`
	CustomFunctions       map[string]interface{} `json:"-"` // Not serializable
}

// ApplySecurityConfig applies custom security configuration
func (tsm *TemplateSecurityManager) ApplySecurityConfig(config *TemplateSecurityConfig) error {
	if config.MaxTemplateSize > 0 {
		tsm.maxTemplateSize = config.MaxTemplateSize
	}

	tsm.sandboxMode = config.SandboxMode

	// Update allowed functions if specified
	if len(config.AllowedFunctions) > 0 {
		// Reset to empty map and add only specified functions
		newFunctions := make(map[string]interface{})
		for _, funcName := range config.AllowedFunctions {
			if fn, exists := tsm.allowedFunctions[funcName]; exists {
				newFunctions[funcName] = fn
			}
		}
		tsm.allowedFunctions = newFunctions
	}

	// Add custom functions if provided
	if len(config.CustomFunctions) > 0 {
		for name, fn := range config.CustomFunctions {
			tsm.allowedFunctions[name] = fn
		}
	}

	// Update blocked patterns if specified
	if len(config.BlockedPatterns) > 0 {
		tsm.blockedPatterns = make([]*regexp.Regexp, 0, len(config.BlockedPatterns))
		for _, pattern := range config.BlockedPatterns {
			compiled, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern '%s': %w", pattern, err)
			}
			tsm.blockedPatterns = append(tsm.blockedPatterns, compiled)
		}
	}

	return nil
}
