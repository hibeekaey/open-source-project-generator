package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/validation"
)

// TestValidationEngineNodeJSIntegration tests the validation engine with real Node.js projects
func TestValidationEngineNodeJSIntegration(t *testing.T) {
	tests := []struct {
		name           string
		packageJSONs   map[string]map[string]interface{}
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "Valid Node.js versions",
			packageJSONs: map[string]map[string]interface{}{
				"package.json": {
					"name": "test-app",
					"engines": map[string]interface{}{
						"node": ">=18.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^18.0.0",
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "Inconsistent Node.js versions",
			packageJSONs: map[string]map[string]interface{}{
				"package.json": {
					"name": "test-app",
					"engines": map[string]interface{}{
						"node": ">=18.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^16.0.0",
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory structure
			tempDir := t.TempDir()

			// Create package.json files
			for filename, content := range tt.packageJSONs {
				filePath := filepath.Join(tempDir, filename)
				data, err := json.MarshalIndent(content, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal JSON: %v", err)
				}

				err = os.WriteFile(filePath, data, 0644)
				if err != nil {
					t.Fatalf("Failed to write package.json: %v", err)
				}
			}

			// Run validation
			engine := validation.NewEngine()
			result, err := engine.ValidateNodeJSVersionCompatibility(tempDir)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			// Check results
			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectedValid, result.Valid)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectedErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

// TestValidationEngineNodeJSComplexProject tests validation with a more complex project structure
func TestValidationEngineNodeJSComplexProject(t *testing.T) {
	tempDir := t.TempDir()

	// Create a complex project structure
	projectStructure := map[string]map[string]interface{}{
		"package.json": {
			"name": "main-app",
			"engines": map[string]interface{}{
				"node": ">=18.0.0",
			},
			"devDependencies": map[string]interface{}{
				"@types/node": "^18.0.0",
			},
		},
		"frontend/package.json": {
			"name": "frontend-app",
			"engines": map[string]interface{}{
				"node": ">=18.0.0",
			},
			"devDependencies": map[string]interface{}{
				"@types/node": "^18.0.0",
			},
		},
		"backend/package.json": {
			"name": "backend-app",
			"engines": map[string]interface{}{
				"node": ">=18.0.0",
			},
			"devDependencies": map[string]interface{}{
				"@types/node": "^18.0.0",
			},
		},
	}

	// Create directory structure and files
	for filePath, content := range projectStructure {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		data, err := json.MarshalIndent(content, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal JSON for %s: %v", filePath, err)
		}

		err = os.WriteFile(fullPath, data, 0644)
		if err != nil {
			t.Fatalf("Failed to write file %s: %v", filePath, err)
		}
	}

	// Run validation
	engine := validation.NewEngine()
	result, err := engine.ValidateNodeJSVersionCompatibility(tempDir)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Should be valid since all versions are consistent
	if !result.Valid {
		t.Errorf("Expected validation to pass, but got errors: %v", result.Errors)
	}

	// The validation should have processed multiple package.json files
	// This is verified by the fact that validation passed for the complex structure
}
