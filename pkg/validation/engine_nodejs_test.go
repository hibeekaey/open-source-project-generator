package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestValidateNodeJSVersionCompatibility(t *testing.T) {
	tests := []struct {
		name           string
		packageJSONs   map[string]map[string]interface{}
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "consistent versions",
			packageJSONs: map[string]map[string]interface{}{
				"frontend/package.json": {
					"name": "frontend-app",
					"engines": map[string]interface{}{
						"node": ">=20.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
				"admin/package.json": {
					"name": "admin-app",
					"engines": map[string]interface{}{
						"node": ">=20.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "inconsistent engine versions",
			packageJSONs: map[string]map[string]interface{}{
				"frontend/package.json": {
					"name": "frontend-app",
					"engines": map[string]interface{}{
						"node": ">=20.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
				"admin/package.json": {
					"name": "admin-app",
					"engines": map[string]interface{}{
						"node": ">=18.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "incompatible runtime and types versions",
			packageJSONs: map[string]map[string]interface{}{
				"frontend/package.json": {
					"name": "frontend-app",
					"engines": map[string]interface{}{
						"node": ">=20.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^18.17.0",
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "missing engines in some files",
			packageJSONs: map[string]map[string]interface{}{
				"frontend/package.json": {
					"name": "frontend-app",
					"engines": map[string]interface{}{
						"node": ">=20.0.0",
					},
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
				"admin/package.json": {
					"name": "admin-app",
					"devDependencies": map[string]interface{}{
						"@types/node": "^20.17.0",
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "nodejs-validation-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create package.json files
			for filePath, content := range tt.packageJSONs {
				fullPath := filepath.Join(tempDir, filePath)
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}

				data, err := json.MarshalIndent(content, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal JSON: %v", err)
				}

				if err := os.WriteFile(fullPath, data, 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", fullPath, err)
				}
			}

			// Run validation
			engine := NewEngine()
			result, err := engine.ValidateNodeJSVersionCompatibility(tempDir)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			// Check results
			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, result.Valid)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectedErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidateNodeJSVersionConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		config         *models.NodeVersionConfig
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid configuration",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "invalid runtime format",
			config: &models.NodeVersionConfig{
				Runtime:      "invalid-version",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectedValid:  false,
			expectedErrors: 2, // Format error + compatibility error
		},
		{
			name: "incompatible runtime and types",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^18.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "invalid docker image",
			config: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "python:3.9",
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "empty configuration",
			config: &models.NodeVersionConfig{
				Runtime:      "",
				TypesPackage: "",
				NPMVersion:   "",
				DockerImage:  "",
			},
			expectedValid:  false,
			expectedErrors: 5, // 4 empty field errors + 1 compatibility error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine()
			result, err := engine.ValidateNodeJSVersionConfiguration(tt.config)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, result.Valid)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectedErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidateCrossTemplateVersionConsistency(t *testing.T) {
	tests := []struct {
		name             string
		templateFiles    map[string]string
		expectedWarnings int
	}{
		{
			name: "consistent frontend templates",
			templateFiles: map[string]string{
				"frontend/nextjs-app/package.json.tmpl": `{
					"name": "test-app",
					"engines": {
						"node": ">=20.0.0"
					},
					"devDependencies": {
						"@types/node": "^20.17.0"
					}
				}`,
				"frontend/nextjs-admin/package.json.tmpl": `{
					"name": "test-admin",
					"engines": {
						"node": ">=20.0.0"
					},
					"devDependencies": {
						"@types/node": "^20.17.0"
					}
				}`,
			},
			expectedWarnings: 1, // Multiple version references detected (runtime and types)
		},
		{
			name: "inconsistent docker images",
			templateFiles: map[string]string{
				"frontend/Dockerfile.tmpl": `FROM node:20-alpine
WORKDIR /app
COPY package.json ./`,
				"admin/Dockerfile.tmpl": `FROM node:18-alpine
WORKDIR /app
COPY package.json ./`,
			},
			expectedWarnings: 1,
		},
		{
			name: "mixed template types",
			templateFiles: map[string]string{
				"frontend/package.json.tmpl": `{
					"engines": {
						"node": ">=20.0.0"
					}
				}`,
				"Dockerfile.tmpl": `FROM node:20-alpine`,
				".github/workflows/ci.yml.tmpl": `
name: CI
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-node@v3
        with:
          node-version: '20'`,
			},
			expectedWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "template-consistency-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create template files
			for filePath, content := range tt.templateFiles {
				fullPath := filepath.Join(tempDir, filePath)
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}

				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", fullPath, err)
				}
			}

			// Run validation
			engine := NewEngine()
			result, err := engine.ValidateCrossTemplateVersionConsistency(tempDir)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			// Check warnings count
			if len(result.Warnings) != tt.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d warnings: %v", tt.expectedWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

func TestExtractNodeJSVersionFromPackageJSON(t *testing.T) {
	tests := []struct {
		name          string
		packageJSON   map[string]interface{}
		expectedInfo  NodeJSVersionInfo
		expectedError bool
	}{
		{
			name: "complete package.json",
			packageJSON: map[string]interface{}{
				"name": "test-app",
				"engines": map[string]interface{}{
					"node": ">=20.0.0",
					"npm":  ">=10.0.0",
				},
				"devDependencies": map[string]interface{}{
					"@types/node": "^20.17.0",
					"typescript":  "^5.0.0",
				},
			},
			expectedInfo: NodeJSVersionInfo{
				EnginesNode:  ">=20.0.0",
				TypesNode:    "^20.17.0",
				NPMVersion:   ">=10.0.0",
				HasEngines:   true,
				HasTypesNode: true,
			},
			expectedError: false,
		},
		{
			name: "missing engines",
			packageJSON: map[string]interface{}{
				"name": "test-app",
				"devDependencies": map[string]interface{}{
					"@types/node": "^20.17.0",
				},
			},
			expectedInfo: NodeJSVersionInfo{
				EnginesNode:  "",
				TypesNode:    "^20.17.0",
				NPMVersion:   "",
				HasEngines:   false,
				HasTypesNode: true,
			},
			expectedError: false,
		},
		{
			name: "types in dependencies",
			packageJSON: map[string]interface{}{
				"name": "test-app",
				"engines": map[string]interface{}{
					"node": ">=20.0.0",
				},
				"dependencies": map[string]interface{}{
					"@types/node": "^20.17.0",
				},
			},
			expectedInfo: NodeJSVersionInfo{
				EnginesNode:  ">=20.0.0",
				TypesNode:    "^20.17.0",
				NPMVersion:   "",
				HasEngines:   true,
				HasTypesNode: true,
			},
			expectedError: false,
		},
		{
			name: "minimal package.json",
			packageJSON: map[string]interface{}{
				"name": "test-app",
			},
			expectedInfo: NodeJSVersionInfo{
				EnginesNode:  "",
				TypesNode:    "",
				NPMVersion:   "",
				HasEngines:   false,
				HasTypesNode: false,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "package-*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			// Write package.json content
			data, err := json.MarshalIndent(tt.packageJSON, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			if _, err := tempFile.Write(data); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tempFile.Close()

			// Test extraction
			engine := NewEngine().(*Engine)
			info, err := engine.extractNodeJSVersionFromPackageJSON(tempFile.Name())

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil {
				if info.EnginesNode != tt.expectedInfo.EnginesNode {
					t.Errorf("Expected EnginesNode=%s, got %s", tt.expectedInfo.EnginesNode, info.EnginesNode)
				}
				if info.TypesNode != tt.expectedInfo.TypesNode {
					t.Errorf("Expected TypesNode=%s, got %s", tt.expectedInfo.TypesNode, info.TypesNode)
				}
				if info.NPMVersion != tt.expectedInfo.NPMVersion {
					t.Errorf("Expected NPMVersion=%s, got %s", tt.expectedInfo.NPMVersion, info.NPMVersion)
				}
				if info.HasEngines != tt.expectedInfo.HasEngines {
					t.Errorf("Expected HasEngines=%v, got %v", tt.expectedInfo.HasEngines, info.HasEngines)
				}
				if info.HasTypesNode != tt.expectedInfo.HasTypesNode {
					t.Errorf("Expected HasTypesNode=%v, got %v", tt.expectedInfo.HasTypesNode, info.HasTypesNode)
				}
			}
		})
	}
}

func TestGroupTemplatesByType(t *testing.T) {
	engine := NewEngine().(*Engine)

	templateFiles := []string{
		"frontend/nextjs-app/package.json.tmpl",
		"frontend/nextjs-admin/package.json.tmpl",
		"frontend/shared-components/package.json.tmpl",
		"backend/go-gin/Dockerfile.tmpl",
		"frontend/nextjs-app/Dockerfile.tmpl",
		"infrastructure/docker/docker-compose.yml.tmpl",
		".github/workflows/ci.yml.tmpl",
		".github/workflows/deploy.yml.tmpl",
		"backend/go-gin/go.mod.tmpl",
		"docs/README.md.tmpl",
	}

	groups := engine.groupTemplatesByType(templateFiles)

	// Check frontend group
	expectedFrontend := 4 // 3 package.json + 1 Dockerfile in frontend
	if len(groups["frontend"]) != expectedFrontend {
		t.Errorf("Expected %d frontend templates, got %d: %v", expectedFrontend, len(groups["frontend"]), groups["frontend"])
	}

	// Check docker group
	expectedDocker := 3 // 2 Dockerfiles + 1 docker-compose
	if len(groups["docker"]) != expectedDocker {
		t.Errorf("Expected %d docker templates, got %d: %v", expectedDocker, len(groups["docker"]), groups["docker"])
	}

	// Check CI group
	expectedCI := 2 // 2 GitHub Actions workflows
	if len(groups["ci"]) != expectedCI {
		t.Errorf("Expected %d CI templates, got %d: %v", expectedCI, len(groups["ci"]), groups["ci"])
	}
}

func TestValidateNodeJSVersionConsistency(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name             string
		nodeVersions     map[string]NodeJSVersionInfo
		expectedValid    bool
		expectedErrors   int
		expectedWarnings int
	}{
		{
			name: "consistent versions",
			nodeVersions: map[string]NodeJSVersionInfo{
				"app1/package.json": {
					EnginesNode:  ">=20.0.0",
					TypesNode:    "^20.17.0",
					HasEngines:   true,
					HasTypesNode: true,
				},
				"app2/package.json": {
					EnginesNode:  ">=20.0.0",
					TypesNode:    "^20.17.0",
					HasEngines:   true,
					HasTypesNode: true,
				},
			},
			expectedValid:    true,
			expectedErrors:   0,
			expectedWarnings: 0,
		},
		{
			name: "inconsistent engine versions",
			nodeVersions: map[string]NodeJSVersionInfo{
				"app1/package.json": {
					EnginesNode:  ">=20.0.0",
					TypesNode:    "^20.17.0",
					HasEngines:   true,
					HasTypesNode: true,
				},
				"app2/package.json": {
					EnginesNode:  ">=18.0.0",
					TypesNode:    "^20.17.0",
					HasEngines:   true,
					HasTypesNode: true,
				},
			},
			expectedValid:    false,
			expectedErrors:   1,
			expectedWarnings: 0,
		},
		{
			name: "missing engines in some files",
			nodeVersions: map[string]NodeJSVersionInfo{
				"app1/package.json": {
					EnginesNode:  ">=20.0.0",
					TypesNode:    "^20.17.0",
					HasEngines:   true,
					HasTypesNode: true,
				},
				"app2/package.json": {
					EnginesNode:  "",
					TypesNode:    "^20.17.0",
					HasEngines:   false,
					HasTypesNode: true,
				},
			},
			expectedValid:    true,
			expectedErrors:   0,
			expectedWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err := engine.validateNodeJSVersionConsistency(tt.nodeVersions, result)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, result.Valid)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectedErrors, len(result.Errors), result.Errors)
			}

			if len(result.Warnings) != tt.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tt.expectedWarnings, len(result.Warnings), result.Warnings)
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkValidateNodeJSVersionCompatibility(b *testing.B) {
	// Create a temporary directory with multiple package.json files
	tempDir, err := os.MkdirTemp("", "nodejs-validation-benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create 10 package.json files
	for i := 0; i < 10; i++ {
		dir := filepath.Join(tempDir, "app"+string(rune('0'+i)))
		if err := os.MkdirAll(dir, 0755); err != nil {
			b.Fatalf("Failed to create directory: %v", err)
		}

		packageJSON := map[string]interface{}{
			"name": "test-app-" + string(rune('0'+i)),
			"engines": map[string]interface{}{
				"node": ">=20.0.0",
			},
			"devDependencies": map[string]interface{}{
				"@types/node": "^20.17.0",
			},
		}

		data, _ := json.MarshalIndent(packageJSON, "", "  ")
		filePath := filepath.Join(dir, "package.json")
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			b.Fatalf("Failed to write file: %v", err)
		}
	}

	engine := NewEngine().(*Engine)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ValidateNodeJSVersionCompatibility(tempDir)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}

func TestIntegrationValidatePreGenerationWithNodeJS(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		config         *models.ProjectConfig
		templatePath   string
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid frontend template with NodeJS config",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "Test project",
				License:      "MIT",
				OutputPath:   "/tmp/test",
				Versions: &models.VersionConfig{
					Node: "20.0.0",
					Go:   "1.21.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
					UpdatedAt: time.Now(),
				},
			},
			templatePath:   "frontend/nextjs-app/package.json.tmpl",
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "invalid NodeJS config for frontend template",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "Test project",
				License:      "MIT",
				OutputPath:   "/tmp/test",
				Versions: &models.VersionConfig{
					Node: "20.0.0",
					Go:   "1.21.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.17.0", // Incompatible with runtime
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
					UpdatedAt: time.Now(),
				},
			},
			templatePath:   "frontend/nextjs-app/package.json.tmpl",
			expectedValid:  false,
			expectedErrors: 2, // Version validator error + template-specific error
		},
		{
			name: "missing NodeJS config for frontend template",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "Test project",
				License:      "MIT",
				OutputPath:   "/tmp/test",
				Versions: &models.VersionConfig{
					Node:      "20.0.0",
					Go:        "1.21.0",
					NodeJS:    nil, // Missing NodeJS config
					UpdatedAt: time.Now(),
				},
			},
			templatePath:   "frontend/nextjs-app/package.json.tmpl",
			expectedValid:  false,
			expectedErrors: 2, // Pre-generation error + template-specific error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ValidatePreGeneration(tt.config, tt.templatePath)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			if result.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, result.Valid)
			}

			if len(result.Errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d: %v", tt.expectedErrors, len(result.Errors), result.Errors)
			}
		})
	}
}
