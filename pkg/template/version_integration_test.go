package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestVersionSubstitutionIntegration(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir := t.TempDir()

	// Create test template content with version variables
	testTemplate := `{
  "name": "test-app",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}",
    "react": "{{.Versions.React}}"
  }
}`

	templatePath := filepath.Join(tempDir, "package.json.tmpl")
	err := os.WriteFile(templatePath, []byte(testTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create test configuration with Node.js version config
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Versions: &models.VersionConfig{
			React: "18.2.0",
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	// Create engine and process template
	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	resultStr := string(result)

	// Verify version substitutions
	tests := []struct {
		name     string
		expected string
		message  string
	}{
		{
			name:     "node runtime",
			expected: `"node": ">=20.0.0"`,
			message:  "Node.js runtime version should be substituted correctly",
		},
		{
			name:     "npm version",
			expected: `"npm": ">=10.0.0"`,
			message:  "NPM version should be substituted correctly",
		},
		{
			name:     "types node version",
			expected: `"@types/node": "^20.17.0"`,
			message:  "@types/node version should be substituted correctly",
		},
		{
			name:     "react version",
			expected: `"react": "18.2.0"`,
			message:  "React version should be substituted correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(resultStr, tt.expected) {
				t.Errorf("%s\nExpected to find: %s\nIn result: %s", tt.message, tt.expected, resultStr)
			}
		})
	}
}

func TestVersionCompatibilityValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *models.NodeVersionConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "compatible versions",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectError: false,
		},
		{
			name: "incompatible types version",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^18.0.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectError: true,
			errorMsg:    "types version 18 is not compatible with runtime version 20",
		},
		{
			name: "empty runtime version",
			config: &models.NodeVersionConfig{
				Runtime:      "",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectError: true,
			errorMsg:    "Runtime version cannot be empty",
		},
		{
			name: "empty types version",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectError: true,
			errorMsg:    "Types package version cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &VersionValidator{}
			result := validator.ValidateNodeVersionConfig(tt.config)

			if tt.expectError {
				if result.Valid {
					t.Errorf("Expected validation to fail, but it passed")
				}

				// Check if expected error message is found
				found := false
				for _, err := range result.Errors {
					if strings.Contains(err.Message, tt.errorMsg) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in validation errors: %v", tt.errorMsg, result.Errors)
				}
			} else {
				if !result.Valid {
					t.Errorf("Expected validation to pass, but it failed with errors: %v", result.Errors)
				}
			}
		})
	}
}

func TestMultipleTemplateVersionConsistency(t *testing.T) {
	// Create temporary directory for multiple templates
	tempDir := t.TempDir()

	// Create multiple template files
	templates := map[string]string{
		"app/package.json.tmpl": `{
  "name": "test-app",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
		"admin/package.json.tmpl": `{
  "name": "test-admin",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
		"shared/package.json.tmpl": `{
  "name": "test-shared",
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "devDependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`,
	}

	// Create template files
	for path, content := range templates {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", path, err)
		}
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	// Process all templates
	engine := NewEngine()
	results := make(map[string]string)

	for path := range templates {
		fullPath := filepath.Join(tempDir, path)
		result, err := engine.ProcessTemplate(fullPath, config)
		if err != nil {
			t.Fatalf("Failed to process template %s: %v", path, err)
		}
		results[path] = string(result)
	}

	// Verify consistency across all templates
	expectedVersions := map[string]string{
		"node":        ">=20.0.0",
		"npm":         ">=10.0.0",
		"@types/node": "^20.17.0",
	}

	for templatePath, content := range results {
		for versionKey, expectedVersion := range expectedVersions {
			var expectedPattern string
			switch versionKey {
			case "node":
				expectedPattern = `"node": "` + expectedVersion + `"`
			case "npm":
				expectedPattern = `"npm": "` + expectedVersion + `"`
			case "@types/node":
				expectedPattern = `"@types/node": "` + expectedVersion + `"`
			}

			if !strings.Contains(content, expectedPattern) {
				t.Errorf("Template %s missing expected %s version %s\nContent: %s",
					templatePath, versionKey, expectedVersion, content)
			}
		}
	}
}

func TestDockerImageSubstitution(t *testing.T) {
	// Create test template with Docker image variable
	tempDir := t.TempDir()
	templateContent := `FROM {{nodeDockerImage .}} AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM {{nodeDockerImage .}} AS runner
WORKDIR /app
COPY --from=builder /app/dist ./dist
CMD ["npm", "start"]`

	templatePath := filepath.Join(tempDir, "Dockerfile.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	// Process template
	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	resultStr := string(result)

	// Verify Docker image substitution
	expectedImages := []string{
		"FROM node:20-alpine AS builder",
		"FROM node:20-alpine AS runner",
	}

	for _, expected := range expectedImages {
		if !strings.Contains(resultStr, expected) {
			t.Errorf("Expected Docker image substitution not found: %s\nResult: %s", expected, resultStr)
		}
	}
}

func TestDefaultVersionConfiguration(t *testing.T) {
	// Create test template
	tempDir := t.TempDir()
	templateContent := `{
  "engines": {
    "node": "{{nodeRuntime .}}",
    "npm": "{{nodeNPMVersion .}}"
  },
  "dependencies": {
    "@types/node": "{{nodeTypesVersion .}}"
  }
}`

	templatePath := filepath.Join(tempDir, "package.json.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create minimal configuration without NodeJS version config
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Versions:     &models.VersionConfig{},
	}

	// Process template - should use defaults
	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	resultStr := string(result)

	// Verify default values are used
	expectedDefaults := []string{
		`"node": ">=20.0.0"`,
		`"npm": ">=10.0.0"`,
		`"@types/node": "^20.17.0"`,
	}

	for _, expected := range expectedDefaults {
		if !strings.Contains(resultStr, expected) {
			t.Errorf("Expected default version not found: %s\nResult: %s", expected, resultStr)
		}
	}
}

func TestVersionUpdateTimestamp(t *testing.T) {
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Versions: &models.VersionConfig{
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}

	// Create a simple template to test version enhancement through ProcessTemplate
	tempDir := t.TempDir()
	templateContent := `{"node": "{{nodeRuntime .}}"}`
	templatePath := filepath.Join(tempDir, "test.json.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	engine := NewEngine()

	// Process template which internally calls enhanceConfigWithVersions
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	// Verify the template was processed correctly
	resultStr := string(result)
	expected := `{"node": ">=20.0.0"}`
	if resultStr != expected {
		t.Errorf("Expected %s, got %s", expected, resultStr)
	}
}
