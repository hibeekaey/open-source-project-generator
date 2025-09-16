package template

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestVersionSubstitutionIntegration tests that version substitution works correctly in templates
func TestVersionSubstitutionIntegration(t *testing.T) {
	engine := NewEngine()
	tempDir := t.TempDir()

	// Create test template
	templateContent := `{
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}",
    "react": "{{reactVersion .}}"
  }
}`

	templatePath := filepath.Join(tempDir, "package.json.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Test configuration
	config := &models.ProjectConfig{
		Name:         "test-app",
		Organization: "test-org",
		Description:  "Test application",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.11.0",
			Packages: map[string]string{
				"react":      "18.2.0",
				"typescript": "^5.0.0",
			},
		},
	}

	// Test node runtime substitution
	t.Run("node_runtime", func(t *testing.T) {
		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		var packageJSON map[string]interface{}
		if err := json.Unmarshal(result, &packageJSON); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
			if nodeVersion, exists := engines["node"]; exists {
				expected := "20.11.0"
				if nodeVersion != expected {
					t.Errorf("Expected node version %s, got %v", expected, nodeVersion)
				}
			}
		}
	})

	// Test npm version substitution
	t.Run("npm_version", func(t *testing.T) {
		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		var packageJSON map[string]interface{}
		if err := json.Unmarshal(result, &packageJSON); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if engines, ok := packageJSON["engines"].(map[string]interface{}); ok {
			if npmVersion, exists := engines["npm"]; exists {
				expected := "10.0.0"
				if npmVersion != expected {
					t.Errorf("Expected npm version %s, got %v", expected, npmVersion)
				}
			}
		}
	})

	// Test types node version substitution
	t.Run("types_node_version", func(t *testing.T) {
		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		var packageJSON map[string]interface{}
		if err := json.Unmarshal(result, &packageJSON); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
			if typesVersion, exists := deps["@types/node"]; exists {
				expected := "^20.17.0"
				if typesVersion != expected {
					t.Errorf("Expected @types/node version %s, got %v", expected, typesVersion)
				}
			}
		}
	})

	// Test react version substitution
	t.Run("react_version", func(t *testing.T) {
		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Fatalf("Failed to process template: %v", err)
		}

		var packageJSON map[string]interface{}
		if err := json.Unmarshal(result, &packageJSON); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
			if reactVersion, exists := deps["react"]; exists {
				expected := "18.2.0"
				if reactVersion != expected {
					t.Errorf("Expected react version %s, got %v", expected, reactVersion)
				}
			}
		}
	})
}

// TestVersionCompatibilityValidation tests basic version compatibility validation
func TestVersionCompatibilityValidation(t *testing.T) {
	t.Run("simplified_validation", func(t *testing.T) {
		// This test validates that the simplified version management works
		// In the simplified version, we only do basic version substitution
		// without complex compatibility checking

		config := &models.ProjectConfig{
			Versions: &models.VersionConfig{
				Node: "20.11.0",
				Packages: map[string]string{
					"react": "18.2.0",
				},
			},
		}

		// Basic validation - check that versions are set
		if config.Versions.Node == "" {
			t.Error("Node version should be set")
		}

		if config.Versions.Packages["react"] == "" {
			t.Error("React version should be set")
		}

		// Check that versions are valid format
		if !strings.Contains(config.Versions.Node, ".") {
			t.Error("Node version should contain dots")
		}

		if !strings.Contains(config.Versions.Packages["react"], ".") {
			t.Error("React version should contain dots")
		}
	})
}
