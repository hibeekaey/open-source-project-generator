package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestValidatePreGeneration(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectError  bool
	}{
		{
			name: "valid frontend template should pass",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  false,
		},
		{
			name: "missing NodeJS config for frontend should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  true,
		},
		{
			name: "valid backend template should pass",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Go: "1.22.0",
				},
			},
			templatePath: "templates/backend/go-gin/go.mod.tmpl",
			expectError:  false,
		},
		{
			name: "missing Go version for backend should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			templatePath: "templates/backend/go-gin/go.mod.tmpl",
			expectError:  true,
		},
		{
			name: "nil versions should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions:     nil,
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validatePreGeneration(tt.config, tt.templatePath)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePreGenerationDirectory(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create test template files
	frontendDir := filepath.Join(tempDir, "frontend", "nextjs-app")
	backendDir := filepath.Join(tempDir, "backend", "go-gin")

	err := os.MkdirAll(frontendDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create frontend dir: %v", err)
	}

	err = os.MkdirAll(backendDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backend dir: %v", err)
	}

	// Create test template files
	packageJSONPath := filepath.Join(frontendDir, "package.json.tmpl")
	goModPath := filepath.Join(backendDir, "go.mod.tmpl")

	err = os.WriteFile(packageJSONPath, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json.tmpl: %v", err)
	}

	err = os.WriteFile(goModPath, []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod.tmpl: %v", err)
	}

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		templateDir string
		expectError bool
	}{
		{
			name: "valid config with mixed templates should pass",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					Go:     "1.22.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			templateDir: tempDir,
			expectError: false,
		},
		{
			name: "missing NodeJS config should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Go: "1.22.0",
				},
			},
			templateDir: tempDir,
			expectError: true,
		},
		{
			name: "nil versions should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions:     nil,
			},
			templateDir: tempDir,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validatePreGenerationDirectory(tt.config, tt.templateDir)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateNodeJSVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectError bool
	}{
		{
			name: "valid NodeJS config should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing NodeJS config should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			expectError: true,
		},
		{
			name: "empty runtime version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      "",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
				},
			},
			expectError: true,
		},
		{
			name: "empty types package should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validateNodeJSVersions(tt.config)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateGoVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectError bool
	}{
		{
			name: "valid Go version should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.22.0",
				},
			},
			expectError: false,
		},
		{
			name: "missing Go version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			expectError: true,
		},
		{
			name: "invalid Go version format should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "invalid-version",
				},
			},
			expectError: true,
		},
		{
			name: "unsupported Go version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.19.0",
				},
			},
			expectError: true,
		},
		{
			name: "minimum supported Go version should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.20.0",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validateGoVersions(tt.config)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePackageJSONVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectError  bool
	}{
		{
			name: "compatible versions should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  false,
		},
		{
			name: "incompatible versions should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.0.0", // Incompatible
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  true,
		},
		{
			name: "major version difference should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^16.0.0", // Too far behind
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validatePackageJSONVersions(tt.config, tt.templatePath)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateDockerVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectError  bool
	}{
		{
			name: "valid Node.js Docker image should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						DockerImage: "node:20-alpine",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectError:  false,
		},
		{
			name: "invalid Docker image should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						DockerImage: "python:3.9",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectError:  true,
		},
		{
			name: "empty Docker image should pass with warning",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						DockerImage: "",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validateDockerVersions(tt.config, tt.templatePath)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTemplateTypeDetection(t *testing.T) {
	engine := NewEngine().(*Engine)

	frontendTests := []struct {
		path     string
		expected bool
	}{
		{"templates/frontend/nextjs-app/package.json.tmpl", true},
		{"templates/frontend/nextjs-admin/next.config.js.tmpl", true},
		{"templates/frontend/shared-components/tsconfig.json.tmpl", true},
		{"templates/backend/go-gin/go.mod.tmpl", false},
		{"templates/infrastructure/docker/Dockerfile.tmpl", false},
	}

	for _, tt := range frontendTests {
		t.Run("frontend_"+tt.path, func(t *testing.T) {
			result := engine.isFrontendTemplate(tt.path)
			if result != tt.expected {
				t.Errorf("isFrontendTemplate(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}

	backendTests := []struct {
		path     string
		expected bool
	}{
		{"templates/backend/go-gin/go.mod.tmpl", true},
		{"templates/backend/go-gin/main.go.tmpl", true},
		{"templates/backend/go-gin/internal/config/config.go.tmpl", true},
		{"templates/frontend/nextjs-app/package.json.tmpl", false},
		{"templates/infrastructure/kubernetes/deployment.yaml.tmpl", false},
	}

	for _, tt := range backendTests {
		t.Run("backend_"+tt.path, func(t *testing.T) {
			result := engine.isBackendTemplate(tt.path)
			if result != tt.expected {
				t.Errorf("isBackendTemplate(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCollectTemplateFiles(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create test template files
	frontendDir := filepath.Join(tempDir, "frontend")
	backendDir := filepath.Join(tempDir, "backend")

	err := os.MkdirAll(frontendDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create frontend dir: %v", err)
	}

	err = os.MkdirAll(backendDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create backend dir: %v", err)
	}

	// Create test template files
	templateFiles := []string{
		filepath.Join(frontendDir, "package.json.tmpl"),
		filepath.Join(frontendDir, "next.config.js.tmpl"),
		filepath.Join(backendDir, "go.mod.tmpl"),
		filepath.Join(backendDir, "main.go.tmpl"),
	}

	for _, file := range templateFiles {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", file, err)
		}
	}

	// Create non-template files (should be ignored)
	nonTemplateFiles := []string{
		filepath.Join(frontendDir, "README.md"),
		filepath.Join(backendDir, "config.yaml"),
	}

	for _, file := range nonTemplateFiles {
		err = os.WriteFile(file, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create non-template file %s: %v", file, err)
		}
	}

	// Test collecting template files
	collected, err := engine.collectTemplateFiles(tempDir)
	if err != nil {
		t.Fatalf("collectTemplateFiles failed: %v", err)
	}

	if len(collected) != len(templateFiles) {
		t.Errorf("Expected %d template files, got %d", len(templateFiles), len(collected))
	}

	// Verify all template files were collected
	for _, expectedFile := range templateFiles {
		found := false
		for _, collectedFile := range collected {
			if collectedFile == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Template file %s was not collected", expectedFile)
		}
	}
}

func TestHasTemplateTypes(t *testing.T) {
	engine := NewEngine().(*Engine)

	templateFiles := []string{
		"templates/frontend/nextjs-app/package.json.tmpl",
		"templates/backend/go-gin/go.mod.tmpl",
		"templates/infrastructure/docker/Dockerfile.tmpl",
	}

	if !engine.hasFrontendTemplates(templateFiles) {
		t.Error("Expected to find frontend templates")
	}

	if !engine.hasBackendTemplates(templateFiles) {
		t.Error("Expected to find backend templates")
	}

	frontendOnlyFiles := []string{
		"templates/frontend/nextjs-app/package.json.tmpl",
		"templates/frontend/nextjs-admin/next.config.js.tmpl",
	}

	if !engine.hasFrontendTemplates(frontendOnlyFiles) {
		t.Error("Expected to find frontend templates in frontend-only list")
	}

	if engine.hasBackendTemplates(frontendOnlyFiles) {
		t.Error("Did not expect to find backend templates in frontend-only list")
	}
}

func TestVersionValidationHelpers(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Test isValidGoVersion
	goVersionTests := []struct {
		version string
		valid   bool
	}{
		{"1.22.0", true},
		{"1.20", true},
		{"1.21.5", true},
		{"invalid", false},
		{"1", false},
		{"", false},
		{"1.22.0.1", false},
	}

	for _, tt := range goVersionTests {
		t.Run("go_version_"+tt.version, func(t *testing.T) {
			result := engine.isValidGoVersion(tt.version)
			if result != tt.valid {
				t.Errorf("isValidGoVersion(%s) = %v, expected %v", tt.version, result, tt.valid)
			}
		})
	}

	// Test extractGoMajorMinor
	goMajorMinorTests := []struct {
		version  string
		expected float64
		hasError bool
	}{
		{"1.22.0", 1.22, false},
		{"1.20", 1.20, false},
		{"1.21.5", 1.21, false},
		{"invalid", 0, true},
		{"1", 0, true},
	}

	for _, tt := range goMajorMinorTests {
		t.Run("go_major_minor_"+tt.version, func(t *testing.T) {
			result, err := engine.extractGoMajorMinor(tt.version)
			if tt.hasError {
				if err == nil {
					t.Errorf("extractGoMajorMinor(%s) expected error but got none", tt.version)
				}
			} else {
				if err != nil {
					t.Errorf("extractGoMajorMinor(%s) unexpected error: %v", tt.version, err)
				}
				if result != tt.expected {
					t.Errorf("extractGoMajorMinor(%s) = %f, expected %f", tt.version, result, tt.expected)
				}
			}
		})
	}
}

// Integration test for ProcessTemplate with pre-generation validation
func TestProcessTemplateWithValidation(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Create a temporary template file
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "package.json.tmpl")
	templateContent := `{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "engines": {
    "node": "{{.Versions.NodeJS.Runtime}}"
  },
  "devDependencies": {
    "@types/node": "{{.Versions.NodeJS.TypesPackage}}"
  }
}`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectError bool
	}{
		{
			name: "valid config should process successfully",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid config should fail validation",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.0.0", // Incompatible
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := engine.ProcessTemplate(templatePath, tt.config)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkValidatePreGeneration(b *testing.B) {
	engine := NewEngine().(*Engine)
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			NextJS: "15.5.3",
			React:  "18.2.0",
			NodeJS: &models.NodeVersionConfig{
				Runtime:      ">=20.0.0",
				TypesPackage: "^20.17.0",
				NPMVersion:   ">=10.0.0",
				DockerImage:  "node:20-alpine",
				LTSStatus:    true,
			},
		},
	}
	templatePath := "templates/frontend/nextjs-app/package.json.tmpl"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := engine.validatePreGeneration(config, templatePath)
		if err != nil {
			b.Fatalf("validatePreGeneration failed: %v", err)
		}
	}
}
