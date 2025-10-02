package metadata

import (
	"embed"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testFS embed.FS

func TestNewMetadataParser(t *testing.T) {
	parser := NewMetadataParser(testFS)
	assert.NotNil(t, parser)
	assert.Equal(t, testFS, parser.embeddedFS)
}

func TestMetadataParser_ParseTemplateYAML(t *testing.T) {
	parser := NewMetadataParser(nil)

	tests := []struct {
		name         string
		yamlContent  string
		templateName string
		expected     *models.TemplateMetadata
		expectError  bool
	}{
		{
			name: "complete metadata",
			yamlContent: `
name: test-template
display_name: Test Template
description: A test template
category: backend
technology: Go
version: 1.2.3
tags:
  - test
  - example
dependencies:
  - base
metadata:
  author: Test Author
  license: MIT
  repository: https://github.com/test/repo
  homepage: https://test.com
  keywords:
    - testing
    - template
  created: 2023-01-01T00:00:00Z
  updated: 2023-12-01T00:00:00Z
  min_version: 1.0.0
  max_version: 2.0.0
  variables:
    ProjectName: "Name of the project"
    Description: "Project description"
`,
			templateName: "test-template",
			expected: &models.TemplateMetadata{
				Name:         "test-template",
				DisplayName:  "Test Template",
				Description:  "A test template",
				Category:     "backend",
				Technology:   "Go",
				Version:      "1.2.3",
				Tags:         []string{"test", "example"},
				Dependencies: []string{"base"},
				Author:       "Test Author",
				License:      "MIT",
				Repository:   "https://github.com/test/repo",
				Homepage:     "https://test.com",
				Keywords:     []string{"testing", "template"},
				MinVersion:   "1.0.0",
				MaxVersion:   "2.0.0",
				Variables: map[string]models.TemplateVar{
					"ProjectName": {
						Name:        "ProjectName",
						Type:        "string",
						Description: "Name of the project",
						Required:    false,
					},
					"Description": {
						Name:        "Description",
						Type:        "string",
						Description: "Project description",
						Required:    false,
					},
				},
			},
		},
		{
			name: "minimal metadata with defaults",
			yamlContent: `
name: minimal-template
description: Minimal template
`,
			templateName: "minimal-template",
			expected: &models.TemplateMetadata{
				Name:         "minimal-template",
				DisplayName:  "Minimal Template",
				Description:  "Minimal template",
				Version:      "1.0.0",
				License:      "MIT",
				Author:       "Open Source Project Generator",
				Tags:         []string{},
				Dependencies: []string{},
				Keywords:     []string{},
				Variables:    map[string]models.TemplateVar{},
			},
		},
		{
			name: "empty metadata with template name fallback",
			yamlContent: `
# Empty metadata file
`,
			templateName: "fallback-template",
			expected: &models.TemplateMetadata{
				Name:         "fallback-template",
				DisplayName:  "Fallback Template",
				Version:      "1.0.0",
				License:      "MIT",
				Author:       "Open Source Project Generator",
				Tags:         []string{},
				Dependencies: []string{},
				Keywords:     []string{},
				Variables:    map[string]models.TemplateVar{},
			},
		},
		{
			name:         "invalid yaml",
			yamlContent:  `invalid: yaml: content: [`,
			templateName: "invalid",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseTemplateYAML([]byte(tt.yamlContent), tt.templateName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.DisplayName, result.DisplayName)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Category, result.Category)
			assert.Equal(t, tt.expected.Technology, result.Technology)
			assert.Equal(t, tt.expected.Version, result.Version)
			assert.Equal(t, tt.expected.Tags, result.Tags)
			assert.Equal(t, tt.expected.Dependencies, result.Dependencies)
			assert.Equal(t, tt.expected.Author, result.Author)
			assert.Equal(t, tt.expected.License, result.License)
			assert.Equal(t, tt.expected.Repository, result.Repository)
			assert.Equal(t, tt.expected.Homepage, result.Homepage)
			assert.Equal(t, tt.expected.Keywords, result.Keywords)
			assert.Equal(t, tt.expected.MinVersion, result.MinVersion)
			assert.Equal(t, tt.expected.MaxVersion, result.MaxVersion)
			assert.Equal(t, tt.expected.Variables, result.Variables)

			// Check that time fields are properly set for complete metadata
			if tt.name == "complete metadata" {
				expectedCreated, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
				expectedUpdated, _ := time.Parse(time.RFC3339, "2023-12-01T00:00:00Z")
				assert.Equal(t, expectedCreated, result.CreatedAt)
				assert.Equal(t, expectedUpdated, result.UpdatedAt)
			}
		})
	}
}

func TestMetadataParser_LoadTemplateMetadata(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "metadata_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test template directory
	templateDir := filepath.Join(tempDir, "test-template")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	// Create template.yaml file
	yamlContent := `
name: test-template
description: Test template
version: 1.0.0
metadata:
  author: Test Author
  license: MIT
`
	yamlPath := filepath.Join(templateDir, "template.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	parser := NewMetadataParser(nil)

	// Test loading from file system
	metadata, err := parser.LoadTemplateMetadata(templateDir)
	require.NoError(t, err)
	assert.Equal(t, "test-template", metadata.Name)
	assert.Equal(t, "Test template", metadata.Description)
	assert.Equal(t, "1.0.0", metadata.Version)
	assert.Equal(t, "Test Author", metadata.Author)
	assert.Equal(t, "MIT", metadata.License)

	// Test loading from non-existent directory
	_, err = parser.LoadTemplateMetadata(filepath.Join(tempDir, "nonexistent"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no metadata file found")
}

func TestMetadataParser_LoadTemplateMetadata_WithYmlExtension(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "metadata_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test template directory
	templateDir := filepath.Join(tempDir, "test-template")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	// Create template.yml file (not .yaml)
	yamlContent := `
name: yml-template
description: Template with yml extension
version: 2.0.0
`
	yamlPath := filepath.Join(templateDir, "template.yml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	parser := NewMetadataParser(nil)

	// Test loading from file system
	metadata, err := parser.LoadTemplateMetadata(templateDir)
	require.NoError(t, err)
	assert.Equal(t, "yml-template", metadata.Name)
	assert.Equal(t, "Template with yml extension", metadata.Description)
	assert.Equal(t, "2.0.0", metadata.Version)
}

func TestMetadataParser_LoadMetadataFromFile(t *testing.T) {
	// Create temporary file for test
	tempDir, err := os.MkdirTemp("", "metadata_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	yamlContent := `
name: file-template
description: Template loaded from specific file
version: 3.0.0
metadata:
  author: File Author
`
	yamlPath := filepath.Join(tempDir, "template.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err)

	parser := NewMetadataParser(nil)

	metadata, err := parser.LoadMetadataFromFile(yamlPath)
	require.NoError(t, err)
	assert.Equal(t, "file-template", metadata.Name)
	assert.Equal(t, "Template loaded from specific file", metadata.Description)
	assert.Equal(t, "3.0.0", metadata.Version)
	assert.Equal(t, "File Author", metadata.Author)

	// Test loading from non-existent file
	_, err = parser.LoadMetadataFromFile(filepath.Join(tempDir, "nonexistent.yaml"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read metadata file")
}

func TestMetadataParser_formatDisplayName(t *testing.T) {
	parser := NewMetadataParser(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "Simple"},
		{"kebab-case", "Kebab Case"},
		{"multi-word-template", "Multi Word Template"},
		{"single", "Single"},
		{"", ""},
		{"already-formatted", "Already Formatted"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.formatDisplayName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMetadataParser_ParseMetadataContent(t *testing.T) {
	parser := NewMetadataParser(nil)

	yamlContent := `
name: content-template
description: Template from content
version: 1.5.0
`

	metadata, err := parser.ParseMetadataContent([]byte(yamlContent), "content-template")
	require.NoError(t, err)
	assert.Equal(t, "content-template", metadata.Name)
	assert.Equal(t, "Template from content", metadata.Description)
	assert.Equal(t, "1.5.0", metadata.Version)
}

// Test with embedded filesystem would require actual embedded test files
// This is a placeholder for integration tests with embedded templates
func TestMetadataParser_LoadMetadataFromEmbedded(t *testing.T) {
	parser := NewMetadataParser(nil)

	// Test with nil embedded filesystem
	_, err := parser.LoadMetadataFromEmbedded("test/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "embedded filesystem not available")

	// Test with embedded filesystem but non-existent path
	parser = NewMetadataParser(testFS)
	_, err = parser.LoadMetadataFromEmbedded("nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no metadata file found")
}
