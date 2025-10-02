package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestManager_Integration tests the complete template functionality with coordinated components
func TestManager_Integration(t *testing.T) {
	// Create template engine and manager
	templateEngine := NewEmbeddedEngine()
	manager := NewManager(templateEngine)

	t.Run("ListTemplates", func(t *testing.T) {
		templates, err := manager.ListTemplates(interfaces.TemplateFilter{})
		require.NoError(t, err)
		assert.NotEmpty(t, templates)

		// Verify we have expected templates
		templateNames := make(map[string]bool)
		for _, tmpl := range templates {
			templateNames[tmpl.Name] = true
		}

		// Check for some expected templates
		expectedTemplates := []string{"go-gin", "nextjs-app", "nextjs-admin", "nextjs-home"}
		for _, expected := range expectedTemplates {
			assert.True(t, templateNames[expected], "Expected template %s not found", expected)
		}
	})

	t.Run("GetTemplateInfo", func(t *testing.T) {
		templateInfo, err := manager.GetTemplateInfo("go-gin")
		require.NoError(t, err)
		assert.NotNil(t, templateInfo)
		assert.Equal(t, "go-gin", templateInfo.Name)
		assert.Equal(t, "backend", templateInfo.Category)
		assert.NotEmpty(t, templateInfo.DisplayName)
	})

	t.Run("SearchTemplates", func(t *testing.T) {
		// Search for Go templates
		templates, err := manager.SearchTemplates("go")
		require.NoError(t, err)
		assert.NotEmpty(t, templates)

		// Verify all results contain "go" in some form
		for _, tmpl := range templates {
			found := false
			searchFields := []string{
				tmpl.Name,
				tmpl.DisplayName,
				tmpl.Description,
				tmpl.Technology,
			}
			for _, field := range searchFields {
				if containsIgnoreCaseIntegration(field, "go") {
					found = true
					break
				}
			}
			if !found {
				// Check tags
				for _, tag := range tmpl.Tags {
					if containsIgnoreCaseIntegration(tag, "go") {
						found = true
						break
					}
				}
			}
			assert.True(t, found, "Template %s doesn't match search query 'go'", tmpl.Name)
		}
	})

	t.Run("GetTemplatesByCategory", func(t *testing.T) {
		frontendTemplates, err := manager.GetTemplatesByCategory("frontend")
		require.NoError(t, err)
		assert.NotEmpty(t, frontendTemplates)

		// Verify all templates are frontend
		for _, tmpl := range frontendTemplates {
			assert.Equal(t, "frontend", tmpl.Category)
		}
	})

	t.Run("GetTemplatesByTechnology", func(t *testing.T) {
		// First check what technologies are available
		allTemplates, err := manager.ListTemplates(interfaces.TemplateFilter{})
		require.NoError(t, err)

		// Find a technology that exists
		var testTechnology string
		for _, tmpl := range allTemplates {
			if tmpl.Technology != "" && tmpl.Technology != "Unknown" {
				testTechnology = tmpl.Technology
				break
			}
		}

		if testTechnology == "" {
			t.Skip("No templates with specific technology found")
		}

		techTemplates, err := manager.GetTemplatesByTechnology(testTechnology)
		require.NoError(t, err)
		assert.NotEmpty(t, techTemplates)

		// Verify all templates have the expected technology
		for _, tmpl := range techTemplates {
			assert.Equal(t, testTechnology, tmpl.Technology)
		}
	})

	t.Run("ValidateTemplate", func(t *testing.T) {
		// Create a temporary template directory for testing
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "test-template")
		err := os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		// Create a simple template file
		templateFile := filepath.Join(templateDir, "test.txt.tmpl")
		err = os.WriteFile(templateFile, []byte("Hello {{.Name}}!"), 0644)
		require.NoError(t, err)

		// Create metadata file
		metadataFile := filepath.Join(templateDir, "template.yaml")
		metadataContent := `name: test-template
description: A test template
version: 1.0.0
author: Test Author
license: MIT
`
		err = os.WriteFile(metadataFile, []byte(metadataContent), 0644)
		require.NoError(t, err)

		// Validate the template
		result, err := manager.ValidateTemplate(templateDir)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Valid, "Template should be valid")
	})

	t.Run("ProcessCustomTemplate", func(t *testing.T) {
		// Create a temporary template directory
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "test-template")
		err := os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		// Create a simple template file
		templateFile := filepath.Join(templateDir, "README.md.tmpl")
		templateContent := `# {{.Name}}

{{.Description}}

Author: {{.Organization}}
License: {{.License}}
`
		err = os.WriteFile(templateFile, []byte(templateContent), 0644)
		require.NoError(t, err)

		// Create metadata file
		metadataFile := filepath.Join(templateDir, "template.yaml")
		metadataContent := `name: test-template
description: A test template
version: 1.0.0
author: Test Author
license: MIT
`
		err = os.WriteFile(metadataFile, []byte(metadataContent), 0644)
		require.NoError(t, err)

		// Create output directory
		outputDir := filepath.Join(tempDir, "output")

		// Create project config
		config := &models.ProjectConfig{
			Name:         "TestProject",
			Description:  "A test project",
			Organization: "TestOrg",
			License:      "MIT",
			OutputPath:   outputDir,
		}

		// Process the template
		err = manager.ProcessCustomTemplate(templateDir, config, outputDir)
		require.NoError(t, err)

		// Verify output file was created
		outputFile := filepath.Join(outputDir, "README.md")
		assert.FileExists(t, outputFile)

		// Verify content was processed correctly
		content, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "# TestProject")
		assert.Contains(t, contentStr, "A test project")
		assert.Contains(t, contentStr, "Author: TestOrg")
		assert.Contains(t, contentStr, "License: MIT")
	})

	t.Run("GetTemplateMetadata", func(t *testing.T) {
		metadata, err := manager.GetTemplateMetadata("go-gin")
		require.NoError(t, err)
		assert.NotNil(t, metadata)
		// Metadata might be empty for embedded templates without explicit metadata
	})

	t.Run("GetTemplateVariables", func(t *testing.T) {
		variables, err := manager.GetTemplateVariables("go-gin")
		require.NoError(t, err)
		assert.NotNil(t, variables)

		// Should have at least some variables (either from template metadata or defaults)
		assert.NotEmpty(t, variables)

		// Check if we have default variables or template-specific variables
		hasDefaults := false
		if _, ok := variables["Name"]; ok {
			hasDefaults = true
		}

		if hasDefaults {
			assert.Contains(t, variables, "Name")
			assert.Contains(t, variables, "Organization")
			assert.Contains(t, variables, "Description")
			assert.Contains(t, variables, "License")
		}
	})

	t.Run("PreviewTemplate", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "TestProject",
			Description:  "A test project",
			Organization: "TestOrg",
			License:      "MIT",
			OutputPath:   "/tmp/test",
		}

		preview, err := manager.PreviewTemplate("go-gin", config)
		require.NoError(t, err)
		assert.NotNil(t, preview)
		assert.Equal(t, "go-gin", preview.TemplateName)
		assert.NotEmpty(t, preview.Files)
		assert.Greater(t, preview.Summary.TotalFiles, 0)
	})

	t.Run("CacheOperations", func(t *testing.T) {
		// Test cache template
		err := manager.CacheTemplate("go-gin")
		require.NoError(t, err)

		// Test get cached templates
		cachedTemplates, err := manager.GetCachedTemplates()
		require.NoError(t, err)
		assert.NotEmpty(t, cachedTemplates)

		// Test refresh cache
		err = manager.RefreshTemplateCache()
		require.NoError(t, err)

		// Test clear cache
		err = manager.ClearTemplateCache()
		require.NoError(t, err)
	})

	t.Run("TemplateCompatibility", func(t *testing.T) {
		compatibility, err := manager.GetTemplateCompatibility("go-gin")
		require.NoError(t, err)
		assert.NotNil(t, compatibility)
		assert.NotEmpty(t, compatibility.SupportedPlatforms)
	})

	t.Run("TemplateLocation", func(t *testing.T) {
		location, err := manager.GetTemplateLocation("go-gin")
		require.NoError(t, err)
		assert.NotEmpty(t, location)
		assert.Contains(t, location, "embedded:")
	})

	t.Run("TemplateDependencies", func(t *testing.T) {
		dependencies, err := manager.GetTemplateDependencies("go-gin")
		require.NoError(t, err)
		// Dependencies might be empty, but should not error
		assert.NotNil(t, dependencies)
	})
}

// TestManager_ComponentCoordination tests that components work together correctly
func TestManager_ComponentCoordination(t *testing.T) {
	templateEngine := NewEmbeddedEngine()
	manager := NewManager(templateEngine)

	t.Run("DiscoveryAndCacheCoordination", func(t *testing.T) {
		// First call should populate cache
		templates1, err := manager.ListTemplates(interfaces.TemplateFilter{})
		require.NoError(t, err)
		assert.NotEmpty(t, templates1)

		// Second call should use cache
		templates2, err := manager.ListTemplates(interfaces.TemplateFilter{})
		require.NoError(t, err)
		assert.Equal(t, len(templates1), len(templates2))
	})

	t.Run("ValidationAndProcessingCoordination", func(t *testing.T) {
		// Create a temporary template
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "test-template")
		err := os.MkdirAll(templateDir, 0755)
		require.NoError(t, err)

		// Create template file with validation issues
		templateFile := filepath.Join(templateDir, "test.txt.tmpl")
		err = os.WriteFile(templateFile, []byte("Hello {{.Name}!"), 0644) // Missing closing brace
		require.NoError(t, err)

		// Validation should catch the issue
		result, err := manager.ValidateTemplate(templateDir)
		require.NoError(t, err)
		assert.False(t, result.Valid, "Template should be invalid due to syntax error")
		assert.NotEmpty(t, result.Issues)

		// Fix the template
		err = os.WriteFile(templateFile, []byte("Hello {{.Name}}!"), 0644)
		require.NoError(t, err)

		// Now validation should pass
		result, err = manager.ValidateTemplate(templateDir)
		require.NoError(t, err)
		assert.True(t, result.Valid, "Template should be valid after fix")

		// And processing should work
		outputDir := filepath.Join(tempDir, "output")
		config := &models.ProjectConfig{
			Name:       "TestProject",
			OutputPath: outputDir,
		}

		err = manager.ProcessCustomTemplate(templateDir, config, outputDir)
		require.NoError(t, err)

		// Verify output
		outputFile := filepath.Join(outputDir, "test.txt")
		assert.FileExists(t, outputFile)
		content, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		assert.Equal(t, "Hello TestProject!", string(content))
	})

	t.Run("FilteringAndSearchCoordination", func(t *testing.T) {
		// Test that filtering and search work together
		allTemplates, err := manager.ListTemplates(interfaces.TemplateFilter{})
		require.NoError(t, err)

		// Filter by category
		frontendTemplates, err := manager.GetTemplatesByCategory("frontend")
		require.NoError(t, err)
		assert.Less(t, len(frontendTemplates), len(allTemplates))

		// Search within category should return subset
		searchResults, err := manager.SearchTemplates("nextjs")
		require.NoError(t, err)

		// All nextjs results should be frontend templates
		for _, result := range searchResults {
			if result.Category == "frontend" {
				found := false
				for _, frontend := range frontendTemplates {
					if frontend.Name == result.Name {
						found = true
						break
					}
				}
				assert.True(t, found, "Search result %s should be in frontend templates", result.Name)
			}
		}
	})
}

// TestManager_IntegrationErrorHandling tests error handling across components
func TestManager_IntegrationErrorHandling(t *testing.T) {
	templateEngine := NewEmbeddedEngine()
	manager := NewManager(templateEngine)

	t.Run("NonExistentTemplate", func(t *testing.T) {
		_, err := manager.GetTemplateInfo("non-existent-template")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("InvalidTemplatePath", func(t *testing.T) {
		result, err := manager.ValidateTemplate("/non/existent/path")
		require.NoError(t, err)
		assert.False(t, result.Valid)
		assert.NotEmpty(t, result.Issues)
	})

	t.Run("ProcessingInvalidTemplate", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:       "TestProject",
			OutputPath: "/tmp/test",
		}

		err := manager.ProcessTemplate("non-existent-template", config, "/tmp/output")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// Helper function to check if a string contains another string (case-insensitive)
func containsIgnoreCaseIntegration(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
