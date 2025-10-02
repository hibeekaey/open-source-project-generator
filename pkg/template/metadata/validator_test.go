package metadata

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetadataValidator(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)
	assert.NotNil(t, validator)
	assert.Equal(t, parser, validator.parser)
}

func TestMetadataValidator_ValidateTemplateMetadataFile(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "validator_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	t.Run("valid metadata file", func(t *testing.T) {
		// Create test template directory with valid metadata
		templateDir := filepath.Join(tempDir, "valid-template")
		err = os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		yamlContent := `
name: valid-template
description: Valid template
version: 1.0.0
author: Test Author
`
		yamlPath := filepath.Join(templateDir, "template.yaml")
		err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		issues := validator.ValidateTemplateMetadataFile(templateDir)
		assert.Empty(t, issues)
	})

	t.Run("no metadata file", func(t *testing.T) {
		// Create test template directory without metadata
		templateDir := filepath.Join(tempDir, "no-metadata-template")
		err = os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		issues := validator.ValidateTemplateMetadataFile(templateDir)
		require.Len(t, issues, 1)
		assert.Equal(t, "warning", issues[0].Type)
		assert.Equal(t, "has-metadata", issues[0].Rule)
		assert.Contains(t, issues[0].Message, "No metadata file found")
	})

	t.Run("invalid yaml syntax", func(t *testing.T) {
		// Create test template directory with invalid YAML
		templateDir := filepath.Join(tempDir, "invalid-yaml-template")
		err = os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		invalidYaml := `
name: invalid-template
description: [invalid yaml syntax
`
		yamlPath := filepath.Join(templateDir, "template.yaml")
		err = os.WriteFile(yamlPath, []byte(invalidYaml), 0644)
		require.NoError(t, err)

		issues := validator.ValidateTemplateMetadataFile(templateDir)
		require.NotEmpty(t, issues)

		// Should have an error for invalid YAML syntax
		hasYamlError := false
		for _, issue := range issues {
			if issue.Rule == "yaml-syntax" && issue.Type == "error" {
				hasYamlError = true
				break
			}
		}
		assert.True(t, hasYamlError, "Should have YAML syntax error")
	})
}

func TestMetadataValidator_ValidateMetadataFileContent(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "validator_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	t.Run("valid metadata content", func(t *testing.T) {
		yamlContent := `
name: test-template
description: Test template
version: 1.0.0
author: Test Author
`
		yamlPath := filepath.Join(tempDir, "valid.yaml")
		err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		issues := validator.ValidateMetadataFileContent(yamlPath)
		assert.Empty(t, issues)
	})

	t.Run("missing required fields", func(t *testing.T) {
		yamlContent := `
# Minimal metadata missing recommended fields
name: minimal-template
`
		yamlPath := filepath.Join(tempDir, "minimal.yaml")
		err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
		require.NoError(t, err)

		issues := validator.ValidateMetadataFileContent(yamlPath)

		// Should have warnings for missing recommended fields
		missingFields := []string{"description", "version", "author"}
		for _, field := range missingFields {
			found := false
			for _, issue := range issues {
				if issue.Rule == "required-fields" && issue.Type == "warning" {
					if strings.Contains(issue.Message, "Missing recommended field: "+field) {
						found = true
						break
					}
				}
			}
			assert.True(t, found, "Should warn about missing field: "+field)
		}
	})

	t.Run("unreadable file", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "nonexistent.yaml")

		issues := validator.ValidateMetadataFileContent(nonExistentPath)
		require.NotEmpty(t, issues)
		assert.Equal(t, "error", issues[0].Type)
		assert.Equal(t, "metadata-readable", issues[0].Rule)
		assert.Contains(t, issues[0].Message, "Cannot read metadata file")
	})

	t.Run("path traversal attempt", func(t *testing.T) {
		maliciousPath := "../../../etc/passwd"

		issues := validator.ValidateMetadataFileContent(maliciousPath)
		require.NotEmpty(t, issues)
		assert.Equal(t, "error", issues[0].Type)
		assert.Equal(t, "path-validation", issues[0].Rule)
		assert.Contains(t, issues[0].Message, "Invalid metadata file path")
	})
}

func TestMetadataValidator_ValidateMetadata(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	t.Run("nil metadata", func(t *testing.T) {
		issues := validator.ValidateMetadata(nil)
		require.Len(t, issues, 1)
		assert.Equal(t, "error", issues[0].Type)
		assert.Equal(t, "metadata-exists", issues[0].Rule)
		assert.Contains(t, issues[0].Message, "Metadata cannot be nil")
	})

	t.Run("valid complete metadata", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Name:        "valid-template",
			Description: "A valid template",
			Version:     "1.0.0",
			Author:      "Test Author",
			License:     "MIT",
			Category:    "backend",
			Variables: map[string]models.TemplateVar{
				"ProjectName": {
					Name:        "ProjectName",
					Type:        "string",
					Description: "Name of the project",
					Required:    true,
				},
			},
		}

		issues := validator.ValidateMetadata(metadata)
		assert.Empty(t, issues)
	})

	t.Run("missing required name", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Description: "Template without name",
			Version:     "1.0.0",
		}

		issues := validator.ValidateMetadata(metadata)
		require.NotEmpty(t, issues)

		hasNameError := false
		for _, issue := range issues {
			if issue.Rule == "name-required" && issue.Type == "error" {
				hasNameError = true
				break
			}
		}
		assert.True(t, hasNameError, "Should have error for missing name")
	})

	t.Run("missing recommended fields", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Name: "minimal-template",
		}

		issues := validator.ValidateMetadata(metadata)

		// Should have warnings for missing recommended fields
		recommendedFields := []string{"description", "version", "author", "license"}
		for _, field := range recommendedFields {
			found := false
			for _, issue := range issues {
				if issue.Type == "warning" && contains([]string{issue.Message}, field) {
					found = true
					break
				}
			}
			assert.True(t, found, "Should warn about missing field: "+field)
		}
	})

	t.Run("invalid name format", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Name:        "Invalid_Template_Name",
			Description: "Template with invalid name",
			Version:     "1.0.0",
		}

		issues := validator.ValidateMetadata(metadata)

		hasNameFormatWarning := false
		for _, issue := range issues {
			if issue.Rule == "name-format" && issue.Type == "warning" {
				hasNameFormatWarning = true
				break
			}
		}
		assert.True(t, hasNameFormatWarning, "Should warn about invalid name format")
	})

	t.Run("invalid category", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Name:        "test-template",
			Description: "Template with invalid category",
			Category:    "invalid-category",
		}

		issues := validator.ValidateMetadata(metadata)

		hasCategoryWarning := false
		for _, issue := range issues {
			if issue.Rule == "category-valid" && issue.Type == "warning" {
				hasCategoryWarning = true
				break
			}
		}
		assert.True(t, hasCategoryWarning, "Should warn about invalid category")
	})

	t.Run("invalid version format", func(t *testing.T) {
		metadata := &models.TemplateMetadata{
			Name:        "test-template",
			Description: "Template with invalid version",
			Version:     "invalid-version",
		}

		issues := validator.ValidateMetadata(metadata)

		hasVersionWarning := false
		for _, issue := range issues {
			if issue.Rule == "version-format" && issue.Type == "warning" {
				hasVersionWarning = true
				break
			}
		}
		assert.True(t, hasVersionWarning, "Should warn about invalid version format")
	})
}

func TestMetadataValidator_validateTemplateVariable(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	t.Run("valid variable", func(t *testing.T) {
		templateVar := models.TemplateVar{
			Name:        "ProjectName",
			Type:        "string",
			Description: "Name of the project",
			Required:    true,
		}

		issues := validator.validateTemplateVariable("ProjectName", templateVar)
		assert.Empty(t, issues)
	})

	t.Run("missing variable name", func(t *testing.T) {
		templateVar := models.TemplateVar{
			Type:        "string",
			Description: "Variable without name",
		}

		issues := validator.validateTemplateVariable("TestVar", templateVar)
		require.NotEmpty(t, issues)

		hasNameWarning := false
		for _, issue := range issues {
			if issue.Rule == "variable-name" && issue.Type == "warning" {
				hasNameWarning = true
				break
			}
		}
		assert.True(t, hasNameWarning, "Should warn about missing variable name")
	})

	t.Run("missing variable type", func(t *testing.T) {
		templateVar := models.TemplateVar{
			Name:        "TestVar",
			Description: "Variable without type",
		}

		issues := validator.validateTemplateVariable("TestVar", templateVar)
		require.NotEmpty(t, issues)

		hasTypeWarning := false
		for _, issue := range issues {
			if issue.Rule == "variable-type" && issue.Type == "warning" {
				hasTypeWarning = true
				break
			}
		}
		assert.True(t, hasTypeWarning, "Should warn about missing variable type")
	})

	t.Run("invalid variable type", func(t *testing.T) {
		templateVar := models.TemplateVar{
			Name:        "TestVar",
			Type:        "invalid-type",
			Description: "Variable with invalid type",
		}

		issues := validator.validateTemplateVariable("TestVar", templateVar)
		require.NotEmpty(t, issues)

		hasTypeValidWarning := false
		for _, issue := range issues {
			if issue.Rule == "variable-type-valid" && issue.Type == "warning" {
				hasTypeValidWarning = true
				break
			}
		}
		assert.True(t, hasTypeValidWarning, "Should warn about invalid variable type")
	})

	t.Run("missing variable description", func(t *testing.T) {
		templateVar := models.TemplateVar{
			Name: "TestVar",
			Type: "string",
		}

		issues := validator.validateTemplateVariable("TestVar", templateVar)
		require.NotEmpty(t, issues)

		hasDescriptionInfo := false
		for _, issue := range issues {
			if issue.Rule == "variable-description" && issue.Type == "info" {
				hasDescriptionInfo = true
				break
			}
		}
		assert.True(t, hasDescriptionInfo, "Should have info about missing variable description")
	})
}

func TestMetadataValidator_validateFilePath(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{"valid relative path", "template/metadata.yaml", false},
		{"valid simple path", "metadata.yaml", false},
		{"path traversal attempt", "../../../etc/passwd", true},
		{"absolute path", "/etc/passwd", false},
		{"clean relative path", "./template/metadata.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateFilePath(tt.path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetadataValidator_isValidTemplateName(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid kebab-case", "valid-template-name", true},
		{"valid single word", "template", true},
		{"valid with numbers", "template-v2", true},
		{"invalid uppercase", "Invalid-Template", false},
		{"invalid underscore", "invalid_template", false},
		{"invalid space", "invalid template", false},
		{"invalid start hyphen", "-invalid", false},
		{"invalid end hyphen", "invalid-", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isValidTemplateName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMetadataValidator_isValidVersion(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid semver", "1.0.0", true},
		{"valid major version", "2.1.3", true},
		{"valid with leading zeros", "01.02.03", true},
		{"invalid single number", "1", false},
		{"invalid two numbers", "1.0", false},
		{"invalid four numbers", "1.0.0.1", false},
		{"invalid with letters", "1.0.0-alpha", false},
		{"invalid empty", "", false},
		{"invalid with spaces", "1.0 .0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isValidVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMetadataValidator_contains(t *testing.T) {
	parser := NewMetadataParser(nil)
	validator := NewMetadataValidator(parser)

	slice := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		item     string
		expected bool
	}{
		{"item exists", "banana", true},
		{"item does not exist", "grape", false},
		{"empty item", "", false},
		{"case sensitive", "Apple", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.contains(slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for tests
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}
