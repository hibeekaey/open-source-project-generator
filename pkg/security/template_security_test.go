package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateSecurityManager(t *testing.T) {
	tsm := NewTemplateSecurityManager()

	assert.NotNil(t, tsm)
	assert.NotEmpty(t, tsm.allowedFunctions)
	assert.NotEmpty(t, tsm.blockedPatterns)
	assert.Equal(t, int64(1024*1024), tsm.maxTemplateSize)
	assert.True(t, tsm.sandboxMode)

	// Verify some expected allowed functions
	expectedFunctions := []string{"upper", "lower", "trim", "replace", "default", "join"}
	for _, fn := range expectedFunctions {
		assert.Contains(t, tsm.allowedFunctions, fn)
	}
}

func TestTemplateSecurityManager_ValidateTemplateContent(t *testing.T) {
	tsm := NewTemplateSecurityManager()

	tests := []struct {
		name           string
		content        string
		filePath       string
		expectSecure   bool
		expectIssues   int
		expectWarnings int
		issueContains  string
	}{
		{
			name:         "safe template",
			content:      "Hello {{.Name}}! Your email is {{.Email | lower}}.",
			filePath:     "safe.tmpl",
			expectSecure: true,
			expectIssues: 0,
		},
		{
			name:          "template with exec",
			content:       "{{exec \"rm -rf /\"}}",
			filePath:      "dangerous.tmpl",
			expectSecure:  false,
			expectIssues:  1,
			issueContains: "dangerous pattern",
		},
		{
			name:          "template with system call",
			content:       "{{system \"whoami\"}}",
			filePath:      "system.tmpl",
			expectSecure:  false,
			expectIssues:  1,
			issueContains: "dangerous pattern",
		},
		{
			name:           "template with path traversal",
			content:        "{{template \"../../etc/passwd\"}}",
			filePath:       "traversal.tmpl",
			expectSecure:   false,
			expectIssues:   1,
			expectWarnings: 1, // Also expect a warning about external templates
			issueContains:  "dangerous pattern",
		},
		{
			name:          "template with script tag",
			content:       "<script>alert('xss')</script>",
			filePath:      "xss.tmpl",
			expectSecure:  false,
			expectIssues:  1,
			issueContains: "dangerous pattern",
		},
		{
			name:          "template with javascript protocol",
			content:       "<a href=\"javascript:alert('xss')\">Click</a>",
			filePath:      "js.tmpl",
			expectSecure:  false,
			expectIssues:  1,
			issueContains: "dangerous pattern",
		},
		{
			name:           "template with unsafe actions",
			content:        "{{printf \"Hello %s\" .Name}}",
			filePath:       "printf.tmpl",
			expectSecure:   true,
			expectIssues:   0,
			expectWarnings: 1,
		},
		{
			name:           "template with external includes",
			content:        "{{template \"external.tmpl\"}}",
			filePath:       "include.tmpl",
			expectSecure:   true,
			expectIssues:   0,
			expectWarnings: 1,
		},
		{
			name:           "template with variable assignments",
			content:        "{{$var := .Name}}Hello {{$var}}",
			filePath:       "assign.tmpl",
			expectSecure:   true,
			expectIssues:   0,
			expectWarnings: 1,
		},
		{
			name:          "oversized template",
			content:       strings.Repeat("a", 2*1024*1024), // 2MB
			filePath:      "large.tmpl",
			expectSecure:  false,
			expectIssues:  1,
			issueContains: "exceeds maximum allowed size",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tsm.ValidateTemplateContent(tt.content, tt.filePath)

			assert.Equal(t, tt.expectSecure, result.IsSecure)
			assert.Equal(t, tt.filePath, result.FilePath)
			assert.Equal(t, int64(len(tt.content)), result.FileSize)
			assert.Len(t, result.SecurityIssues, tt.expectIssues)
			assert.Len(t, result.Warnings, tt.expectWarnings)

			if tt.issueContains != "" && len(result.SecurityIssues) > 0 {
				assert.Contains(t, result.SecurityIssues[0], tt.issueContains)
			}
		})
	}
}

func TestTemplateSecurityManager_ValidateTemplateFile(t *testing.T) {
	tempDir := t.TempDir()
	tsm := NewTemplateSecurityManager()

	// Test with safe template file
	safeFile := filepath.Join(tempDir, "safe.tmpl")
	safeContent := "Hello {{.Name}}!"
	err := os.WriteFile(safeFile, []byte(safeContent), 0644)
	require.NoError(t, err)

	result, err := tsm.ValidateTemplateFile(safeFile)
	assert.NoError(t, err)
	assert.True(t, result.IsSecure)
	assert.Equal(t, safeFile, result.FilePath)
	assert.Equal(t, int64(len(safeContent)), result.FileSize)

	// Test with dangerous template file
	dangerousFile := filepath.Join(tempDir, "dangerous.tmpl")
	dangerousContent := "{{exec \"rm -rf /\"}}"
	err = os.WriteFile(dangerousFile, []byte(dangerousContent), 0644)
	require.NoError(t, err)

	result, err = tsm.ValidateTemplateFile(dangerousFile)
	assert.NoError(t, err)
	assert.False(t, result.IsSecure)
	assert.Greater(t, len(result.SecurityIssues), 0)

	// Test with path traversal in file path
	traversalPath := "../../../etc/passwd"
	result, err = tsm.ValidateTemplateFile(traversalPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsafe file path")
	assert.False(t, result.IsSecure)
	assert.Contains(t, result.SecurityIssues[0], "path traversal")

	// Test with non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.tmpl")
	result, err = tsm.ValidateTemplateFile(nonExistentFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat template file")

	// Test with oversized file
	largeFile := filepath.Join(tempDir, "large.tmpl")
	largeContent := strings.Repeat("a", 2*1024*1024) // 2MB
	err = os.WriteFile(largeFile, []byte(largeContent), 0644)
	require.NoError(t, err)

	result, err = tsm.ValidateTemplateFile(largeFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template file too large")
	assert.False(t, result.IsSecure)
	assert.Contains(t, result.SecurityIssues[0], "exceeds maximum allowed size")
}

func TestTemplateSecurityManager_CreateSecureTemplate(t *testing.T) {
	tsm := NewTemplateSecurityManager()

	// Test creating safe template
	safeContent := "Hello {{.Name | upper}}! {{default \"World\" .Greeting}}"
	tmpl, err := tsm.CreateSecureTemplate("safe", safeContent)
	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
	assert.Equal(t, "safe", tmpl.Name())

	// Test creating dangerous template
	dangerousContent := "{{exec \"rm -rf /\"}}"
	tmpl, err = tsm.CreateSecureTemplate("dangerous", dangerousContent)
	assert.Error(t, err)
	assert.Nil(t, tmpl)
	assert.Contains(t, err.Error(), "template security validation failed")

	// Test creating template with invalid syntax
	invalidContent := "{{.Name"
	tmpl, err = tsm.CreateSecureTemplate("invalid", invalidContent)
	assert.Error(t, err)
	assert.Nil(t, tmpl)
	assert.Contains(t, err.Error(), "template parsing failed")
}

func TestTemplateSecurityManager_ProcessTemplateSecurely(t *testing.T) {
	tempDir := t.TempDir()
	tsm := NewTemplateSecurityManager()

	// Create a safe template
	safeContent := "Hello {{.Name | upper}}!"
	tmpl, err := tsm.CreateSecureTemplate("safe", safeContent)
	require.NoError(t, err)

	// Test processing to valid output path
	outputFile := filepath.Join(tempDir, "output.txt")
	data := map[string]interface{}{
		"Name": "world",
	}

	err = tsm.ProcessTemplateSecurely(tmpl, data, outputFile)
	assert.NoError(t, err)

	// Verify output file was created with correct content
	outputContent, err := os.ReadFile(outputFile)
	assert.NoError(t, err)
	assert.Equal(t, "Hello WORLD!", string(outputContent))

	// Test processing to path with traversal
	traversalPath := "../../../etc/passwd"
	err = tsm.ProcessTemplateSecurely(tmpl, data, traversalPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsafe output path")

	// Test processing to nested directory (should create directories)
	nestedOutput := filepath.Join(tempDir, "nested", "dir", "output.txt")
	err = tsm.ProcessTemplateSecurely(tmpl, data, nestedOutput)
	assert.NoError(t, err)

	// Verify nested directories were created
	nestedContent, err := os.ReadFile(nestedOutput)
	assert.NoError(t, err)
	assert.Equal(t, "Hello WORLD!", string(nestedContent))
}

func TestTemplateSecurityManager_ScanTemplateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tsm := NewTemplateSecurityManager()

	// Create test template files
	templates := map[string]string{
		"safe.tmpl":          "Hello {{.Name}}!",
		"dangerous.tmpl":     "{{exec \"rm -rf /\"}}",
		"warning.tmpl":       "{{printf \"Hello %s\" .Name}}",
		"not_template.txt":   "This is not a template",
		"subdir/nested.tmpl": "Nested {{.Value}}",
	}

	for filePath, content := range templates {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Scan directory
	results, err := tsm.ScanTemplateDirectory(tempDir)
	assert.NoError(t, err)

	// Should only scan template files (not .txt files)
	assert.Len(t, results, 4) // safe.tmpl, dangerous.tmpl, warning.tmpl, subdir/nested.tmpl

	// Verify results
	safeResult := results[filepath.Join(tempDir, "safe.tmpl")]
	assert.NotNil(t, safeResult)
	assert.True(t, safeResult.IsSecure)

	dangerousResult := results[filepath.Join(tempDir, "dangerous.tmpl")]
	assert.NotNil(t, dangerousResult)
	assert.False(t, dangerousResult.IsSecure)
	assert.Greater(t, len(dangerousResult.SecurityIssues), 0)

	warningResult := results[filepath.Join(tempDir, "warning.tmpl")]
	assert.NotNil(t, warningResult)
	assert.True(t, warningResult.IsSecure)
	assert.Greater(t, len(warningResult.Warnings), 0)

	nestedResult := results[filepath.Join(tempDir, "subdir", "nested.tmpl")]
	assert.NotNil(t, nestedResult)
	assert.True(t, nestedResult.IsSecure)

	// Test scanning non-existent directory
	results, err = tsm.ScanTemplateDirectory("/nonexistent")
	assert.Error(t, err)
	assert.Empty(t, results)
}

func TestGetTemplateSecuritySummary(t *testing.T) {
	// Create mock results
	results := map[string]*TemplateValidationResult{
		"safe1.tmpl": {
			IsSecure:       true,
			SecurityIssues: []string{},
			Warnings:       []string{},
		},
		"safe2.tmpl": {
			IsSecure:       true,
			SecurityIssues: []string{},
			Warnings:       []string{"Minor warning"},
		},
		"dangerous1.tmpl": {
			IsSecure:       false,
			SecurityIssues: []string{"Dangerous pattern detected"},
			Warnings:       []string{},
		},
		"dangerous2.tmpl": {
			IsSecure:       false,
			SecurityIssues: []string{"Dangerous pattern detected", "Another issue"},
			Warnings:       []string{"Warning"},
		},
		"warning_only.tmpl": {
			IsSecure:       true,
			SecurityIssues: []string{},
			Warnings:       []string{"Common warning", "Another warning"},
		},
	}

	summary := GetTemplateSecuritySummary(results)

	assert.Equal(t, 5, summary["total_templates"])
	assert.Equal(t, 3, summary["secure_templates"])
	assert.Equal(t, 2, summary["insecure_templates"])
	assert.Equal(t, 3, summary["templates_with_warnings"])
	assert.Equal(t, 3, summary["total_issues"])
	assert.Equal(t, 4, summary["total_warnings"])

	criticalIssues := summary["critical_issues"].([]string)
	assert.Len(t, criticalIssues, 2) // Two unique issues

	commonWarnings := summary["common_warnings"].([]string)
	assert.Len(t, commonWarnings, 0) // No warnings appear more than once
}

func TestTemplateSecurityManager_ApplySecurityConfig(t *testing.T) {
	tsm := NewTemplateSecurityManager()

	// Test applying valid configuration
	config := &TemplateSecurityConfig{
		MaxTemplateSize:       2 * 1024 * 1024, // 2MB
		AllowedFunctions:      []string{"upper", "lower"},
		BlockedPatterns:       []string{`test_pattern`},
		SandboxMode:           false,
		AllowExternalIncludes: true,
		CustomFunctions: map[string]interface{}{
			"custom": func(s string) string { return "custom_" + s },
		},
	}

	err := tsm.ApplySecurityConfig(config)
	assert.NoError(t, err)

	// Verify configuration was applied
	assert.Equal(t, int64(2*1024*1024), tsm.maxTemplateSize)
	assert.False(t, tsm.sandboxMode)
	assert.Len(t, tsm.allowedFunctions, 3) // upper, lower, custom
	assert.Contains(t, tsm.allowedFunctions, "upper")
	assert.Contains(t, tsm.allowedFunctions, "lower")
	assert.Contains(t, tsm.allowedFunctions, "custom")
	assert.Len(t, tsm.blockedPatterns, 1)

	// Test applying configuration with invalid regex
	invalidConfig := &TemplateSecurityConfig{
		BlockedPatterns: []string{`[invalid_regex`},
	}

	err = tsm.ApplySecurityConfig(invalidConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid regex pattern")
}

func TestTemplateSecurityManager_EdgeCases(t *testing.T) {
	tsm := NewTemplateSecurityManager()

	// Test with empty template content
	result := tsm.ValidateTemplateContent("", "empty.tmpl")
	assert.True(t, result.IsSecure)
	assert.Equal(t, int64(0), result.FileSize)

	// Test with whitespace-only content
	result = tsm.ValidateTemplateContent("   \n\t  ", "whitespace.tmpl")
	assert.True(t, result.IsSecure)

	// Test with very long file path
	longPath := strings.Repeat("a", 1000) + ".tmpl"
	result = tsm.ValidateTemplateContent("safe content", longPath)
	assert.True(t, result.IsSecure)
	assert.Equal(t, longPath, result.FilePath)

	// Test with multiple dangerous patterns in one template
	multiDangerous := "{{exec \"cmd1\"}} and {{system \"cmd2\"}} and <script>alert('xss')</script>"
	result = tsm.ValidateTemplateContent(multiDangerous, "multi.tmpl")
	assert.False(t, result.IsSecure)
	assert.Greater(t, len(result.SecurityIssues), 1)
	assert.Greater(t, len(result.BlockedPatterns), 1)
}

func TestTemplateSecurityManager_Integration(t *testing.T) {
	tempDir := t.TempDir()
	tsm := NewTemplateSecurityManager()

	// Test complete workflow: validate, create, process

	// 1. Create and validate a template file
	templateFile := filepath.Join(tempDir, "integration.tmpl")
	templateContent := `# {{.ProjectName | upper}}

Welcome to {{.ProjectName}}!

Author: {{.Author}}
Version: {{default "1.0.0" .Version}}

## Features
{{range .Features}}
- {{.}}
{{end}}

Generated on: {{.Timestamp}}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	require.NoError(t, err)

	// 2. Validate template file
	result, err := tsm.ValidateTemplateFile(templateFile)
	require.NoError(t, err)
	assert.True(t, result.IsSecure)

	// 3. Create secure template
	tmpl, err := tsm.CreateSecureTemplate("integration", templateContent)
	require.NoError(t, err)

	// 4. Process template with data
	outputFile := filepath.Join(tempDir, "README.md")
	data := map[string]interface{}{
		"ProjectName": "my-awesome-project",
		"Author":      "John Doe",
		"Version":     "2.0.0",
		"Features":    []string{"Feature 1", "Feature 2", "Feature 3"},
		"Timestamp":   "2023-01-01",
	}

	err = tsm.ProcessTemplateSecurely(tmpl, data, outputFile)
	require.NoError(t, err)

	// 5. Verify output
	outputContent, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	output := string(outputContent)
	assert.Contains(t, output, "# MY-AWESOME-PROJECT")
	assert.Contains(t, output, "Welcome to my-awesome-project!")
	assert.Contains(t, output, "Author: John Doe")
	assert.Contains(t, output, "Version: 2.0.0")
	assert.Contains(t, output, "- Feature 1")
	assert.Contains(t, output, "- Feature 2")
	assert.Contains(t, output, "- Feature 3")
	assert.Contains(t, output, "Generated on: 2023-01-01")

	// 6. Scan directory for templates
	results, err := tsm.ScanTemplateDirectory(tempDir)
	require.NoError(t, err)
	assert.Len(t, results, 1) // Should find the integration.tmpl file

	templateResult := results[templateFile]
	assert.NotNil(t, templateResult)
	assert.True(t, templateResult.IsSecure)

	// 7. Get security summary
	summary := GetTemplateSecuritySummary(results)
	assert.Equal(t, 1, summary["total_templates"])
	assert.Equal(t, 1, summary["secure_templates"])
	assert.Equal(t, 0, summary["insecure_templates"])
}
