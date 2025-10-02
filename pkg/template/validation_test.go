package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemplateValidator(t *testing.T) {
	validator := NewTemplateValidator()
	assert.NotNil(t, validator)
}

func TestTemplateValidator_ValidateTemplate(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		setupTemplate  func(t *testing.T) string
		expectedValid  bool
		expectedIssues int
	}{
		{
			name: "non-existent path",
			setupTemplate: func(t *testing.T) string {
				return "/non/existent/path"
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
		{
			name: "valid template directory",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create template.yaml
				templateYAML := `name: test-template
description: Test template
version: 1.0.0
author: Test Author
`
				err := os.WriteFile(filepath.Join(tmpDir, "template.yaml"), []byte(templateYAML), 0644)
				require.NoError(t, err)

				// Create a template file
				templateContent := `# {{.Name}}

This is a test template for {{.Organization}}.

Description: {{.Description}}
`
				err = os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
				require.NoError(t, err)

				return tmpDir
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
		{
			name: "directory without template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create a regular file (not a template)
				err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Test"), 0644)
				require.NoError(t, err)

				return tmpDir
			},
			expectedValid:  true, // Should be valid but with warnings
			expectedIssues: 0,
		},
		{
			name: "template with syntax errors",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()

				// Create template file with unmatched delimiters
				templateContent := `# {{.Name}

This template has unmatched delimiters {{.Organization
`
				err := os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
				require.NoError(t, err)

				return tmpDir
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.setupTemplate(t)

			result, err := validator.ValidateTemplate(templatePath)
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Equal(t, tt.expectedIssues, len(result.Issues))
		})
	}
}

func TestTemplateValidator_ValidateTemplateStructure(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name        string
		template    *interfaces.TemplateInfo
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid template info",
			template: &interfaces.TemplateInfo{
				Name:     "test-template",
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: false,
		},
		{
			name: "missing name",
			template: &interfaces.TemplateInfo{
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: true,
			errorMsg:    "template name is required",
		},
		{
			name: "missing category",
			template: &interfaces.TemplateInfo{
				Name:    "test-template",
				Version: "1.0.0",
			},
			expectError: true,
			errorMsg:    "template category is required",
		},
		{
			name: "missing version",
			template: &interfaces.TemplateInfo{
				Name:     "test-template",
				Category: "backend",
			},
			expectError: true,
			errorMsg:    "template version is required",
		},
		{
			name: "invalid name format",
			template: &interfaces.TemplateInfo{
				Name:     "TestTemplate",
				Category: "backend",
				Version:  "1.0.0",
			},
			expectError: true,
			errorMsg:    "template name must be in kebab-case format",
		},
		{
			name: "invalid category",
			template: &interfaces.TemplateInfo{
				Name:     "test-template",
				Category: "invalid",
				Version:  "1.0.0",
			},
			expectError: true,
			errorMsg:    "invalid category",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTemplateStructure(tt.template)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplateValidator_ValidateTemplateMetadata(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name        string
		metadata    *interfaces.TemplateMetadata
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid metadata",
			metadata: &interfaces.TemplateMetadata{
				Author:  "Test Author",
				License: "MIT",
			},
			expectError: false,
		},
		{
			name:        "nil metadata",
			metadata:    nil,
			expectError: true,
			errorMsg:    "metadata cannot be nil",
		},
		{
			name: "missing author",
			metadata: &interfaces.TemplateMetadata{
				License: "MIT",
			},
			expectError: true,
			errorMsg:    "metadata author is required",
		},
		{
			name: "missing license",
			metadata: &interfaces.TemplateMetadata{
				Author: "Test Author",
			},
			expectError: true,
			errorMsg:    "metadata license is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTemplateMetadata(tt.metadata)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplateValidator_ValidateCustomTemplate(t *testing.T) {
	validator := NewTemplateValidator()

	// Create a temporary template directory
	tmpDir := t.TempDir()

	// Create template.yaml
	templateYAML := `name: custom-template
description: Custom test template
version: 1.0.0
author: Test Author
`
	err := os.WriteFile(filepath.Join(tmpDir, "template.yaml"), []byte(templateYAML), 0644)
	require.NoError(t, err)

	// Create a template file
	templateContent := `# {{.Name}}

Custom template for {{.Organization}}.
`
	err = os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
	require.NoError(t, err)

	result, err := validator.ValidateCustomTemplate(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
}

func TestTemplateValidator_validateTemplateStructure(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		setupTemplate  func(t *testing.T) string
		expectedIssues int
	}{
		{
			name: "valid directory with template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "test.tmpl"), []byte("{{.Name}}"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 0,
		},
		{
			name: "directory without template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 1, // Warning about no template files
		},
		{
			name: "file instead of directory",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.txt")
				err := os.WriteFile(filePath, []byte("test"), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectedIssues: 1, // Error about not being a directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.setupTemplate(t)
			issues := validator.validateTemplateStructure(templatePath)
			assert.Len(t, issues, tt.expectedIssues)
		})
	}
}

func TestTemplateValidator_validateTemplateMetadataFile(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		setupTemplate  func(t *testing.T) string
		expectedIssues int
	}{
		{
			name: "valid template.yaml",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				templateYAML := `name: test-template
description: Test template
version: 1.0.0
author: Test Author
`
				err := os.WriteFile(filepath.Join(tmpDir, "template.yaml"), []byte(templateYAML), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 0,
		},
		{
			name: "no metadata file",
			setupTemplate: func(t *testing.T) string {
				return t.TempDir()
			},
			expectedIssues: 1, // Warning about missing metadata
		},
		{
			name: "invalid yaml syntax",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				invalidYAML := `name test-template
description Test template
invalid line without colon
version: 1.0.0
`
				err := os.WriteFile(filepath.Join(tmpDir, "template.yaml"), []byte(invalidYAML), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 6, // Multiple warnings about YAML syntax and missing fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.setupTemplate(t)
			issues := validator.validateTemplateMetadataFile(templatePath)
			assert.Len(t, issues, tt.expectedIssues)
		})
	}
}

func TestTemplateValidator_validateTemplateFiles(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		setupTemplate  func(t *testing.T) string
		expectedIssues int
	}{
		{
			name: "valid template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				templateContent := `# {{.Name}}

Description: {{.Description}}
Organization: {{.Organization}}
`
				err := os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 0,
		},
		{
			name: "template with syntax errors",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				templateContent := `# {{.Name}

Unmatched delimiter: {{.Organization
`
				err := os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 2, // Error about unmatched delimiters + info about no common vars
		},
		{
			name: "template without common variables",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				templateContent := `# {{.CustomVar}}

This template uses custom variables: {{.AnotherVar}}
`
				err := os.WriteFile(filepath.Join(tmpDir, "README.md.tmpl"), []byte(templateContent), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedIssues: 1, // Info about no common variables
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.setupTemplate(t)
			issues := validator.validateTemplateFiles(templatePath)
			assert.Len(t, issues, tt.expectedIssues)
		})
	}
}

func TestTemplateValidator_validateTemplateFile(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name           string
		setupFile      func(t *testing.T) string
		expectedIssues int
	}{
		{
			name: "valid template file",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.tmpl")
				content := `# {{.Name}}

Organization: {{.Organization}}
Description: {{.Description}}
`
				err := os.WriteFile(filePath, []byte(content), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectedIssues: 0,
		},
		{
			name: "unmatched delimiters",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.tmpl")
				content := `# {{.Name}

Missing closing delimiter
`
				err := os.WriteFile(filePath, []byte(content), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectedIssues: 2, // Error about unmatched delimiters + info about no common vars
		},
		{
			name: "no common variables",
			setupFile: func(t *testing.T) string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "test.tmpl")
				content := `# {{.CustomVar}}

Custom content: {{.AnotherVar}}
`
				err := os.WriteFile(filePath, []byte(content), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectedIssues: 1, // Info issue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile(t)
			issues := validator.validateTemplateFile(filePath)
			assert.Len(t, issues, tt.expectedIssues)
		})
	}
}

func TestTemplateValidator_hasTemplateFiles(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name          string
		setupTemplate func(t *testing.T) string
		expectedHas   bool
	}{
		{
			name: "directory with template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "test.tmpl"), []byte("{{.Name}}"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedHas: true,
		},
		{
			name: "directory without template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedHas: false,
		},
		{
			name: "nested template files",
			setupTemplate: func(t *testing.T) string {
				tmpDir := t.TempDir()
				subDir := filepath.Join(tmpDir, "subdir")
				err := os.MkdirAll(subDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(subDir, "nested.tmpl"), []byte("{{.Name}}"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedHas: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.setupTemplate(t)
			hasFiles, err := validator.hasTemplateFiles(templatePath)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedHas, hasFiles)
		})
	}
}

func TestTemplateValidator_isValidTemplateName(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid kebab-case", "test-template", true},
		{"valid single word", "template", true},
		{"valid with numbers", "template-v2", true},
		{"invalid uppercase", "Test-Template", false},
		{"invalid underscore", "test_template", false},
		{"invalid space", "test template", false},
		{"invalid starting hyphen", "-test-template", false},
		{"invalid ending hyphen", "test-template-", false},
		{"invalid special chars", "test@template", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isValidTemplateName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateValidator_contains(t *testing.T) {
	validator := NewTemplateValidator()

	slice := []string{"backend", "frontend", "mobile"}

	tests := []struct {
		name     string
		item     string
		expected bool
	}{
		{"contains backend", "backend", true},
		{"contains frontend", "frontend", true},
		{"does not contain invalid", "invalid", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.contains(slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateValidator_validateFilePath(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid relative path",
			path:        "templates/test/file.tmpl",
			expectError: false,
		},
		{
			name:        "valid absolute path",
			path:        "/tmp/templates/test.tmpl",
			expectError: false,
		},
		{
			name:        "path traversal attempt",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "system directory access",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "access to system directory not allowed",
		},
		{
			name:        "proc directory access",
			path:        "/proc/version",
			expectError: true,
			errorMsg:    "access to system directory not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateFilePath(tt.path)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
