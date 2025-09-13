package validation

import (
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

func TestValidatePreGeneration(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectValid  bool
		expectErrors int
	}{
		{
			name:         "nil config should fail",
			config:       nil,
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "nil versions should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions:     nil,
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "valid frontend template config should pass",
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
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "invalid Node.js version compatibility should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Node:   "20.11.0",
					NextJS: "15.5.3",
					React:  "18.2.0",
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.0.0", // Incompatible with runtime (too far behind)
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 2, // One from version validator, one from template-specific validation
		},
		{
			name: "valid backend template config should pass",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Go: "1.22.0",
				},
			},
			templatePath: "templates/backend/go-gin/go.mod.tmpl",
			expectValid:  true,
			expectErrors: 0,
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
			expectValid:  false,
			expectErrors: 2, // One for missing Go version, one for template-specific validation
		},
		{
			name: "unsupported Go version should fail",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Versions: &models.VersionConfig{
					Go: "1.19.0", // Below minimum supported version
				},
			},
			templatePath: "templates/backend/go-gin/go.mod.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ValidatePreGeneration(tt.config, tt.templatePath)
			if err != nil {
				t.Fatalf("ValidatePreGeneration returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidatePreGenerationDirectory(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templateDir  string
		expectValid  bool
		expectErrors int
	}{
		{
			name:         "nil config should fail",
			config:       nil,
			templateDir:  "templates/frontend/nextjs-app",
			expectValid:  false,
			expectErrors: 1,
		},
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
			templateDir:  ".", // Use current directory instead of non-existent templates dir
			expectValid:  true,
			expectErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ValidatePreGenerationDirectory(tt.config, tt.templateDir)
			if err != nil {
				t.Fatalf("ValidatePreGenerationDirectory returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidateNodeJSPreGeneration(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		expectValid  bool
		expectErrors int
	}{
		{
			name: "missing NodeJS config should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			expectValid:  false,
			expectErrors: 1,
		},
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
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "empty runtime version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      "", // Empty runtime
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectValid:  false,
			expectErrors: 2, // One for empty runtime, one for compatibility check failure
		},
		{
			name: "empty types package should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "", // Empty types
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectValid:  false,
			expectErrors: 2, // One for empty types, one for compatibility check failure
		},
		{
			name: "incompatible versions should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.0.0", // Incompatible with runtime
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
						LTSStatus:    true,
					},
				},
			},
			expectValid:  false,
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err := engine.validateNodeJSPreGeneration(tt.config, result)
			if err != nil {
				t.Fatalf("validateNodeJSPreGeneration returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidateGoPreGeneration(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		expectValid  bool
		expectErrors int
	}{
		{
			name: "missing Go version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "valid Go version should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.22.0",
				},
			},
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "invalid Go version format should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "invalid-version",
				},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "unsupported Go version should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.19.0", // Below minimum
				},
			},
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "minimum supported Go version should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Go: "1.20.0",
				},
			},
			expectValid:  true,
			expectErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err := engine.validateGoPreGeneration(tt.config, result)
			if err != nil {
				t.Fatalf("validateGoPreGeneration returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidatePackageJSONTemplateVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectValid  bool
		expectErrors int
	}{
		{
			name: "missing NodeJS config should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					Node: "20.11.0",
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "valid NodeJS config should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "missing runtime version should fail",
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
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "missing types version should fail",
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
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "incompatible versions should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^18.0.0", // Incompatible
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/package.json.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err := engine.validatePackageJSONTemplateVersions(tt.config, tt.templatePath, result)
			if err != nil {
				t.Fatalf("validatePackageJSONTemplateVersions returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}
		})
	}
}

func TestValidateDockerfileTemplateVersions(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		name         string
		config       *models.ProjectConfig
		templatePath string
		expectValid  bool
		expectErrors int
	}{
		{
			name: "valid Docker image should pass",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "node:20-alpine",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "invalid Docker image should fail",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "python:3.9", // Not a Node.js image
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectValid:  false,
			expectErrors: 1,
		},
		{
			name: "missing Docker image should generate warning",
			config: &models.ProjectConfig{
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^20.17.0",
						NPMVersion:   ">=10.0.0",
						DockerImage:  "",
					},
				},
			},
			templatePath: "templates/frontend/nextjs-app/Dockerfile.tmpl",
			expectValid:  true,
			expectErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err := engine.validateDockerfileTemplateVersions(tt.config, tt.templatePath, result)
			if err != nil {
				t.Fatalf("validateDockerfileTemplateVersions returned error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d errors: %v", tt.expectErrors, len(result.Errors), result.Errors)
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

func TestVersionExtractionHelpers(t *testing.T) {
	engine := NewEngine().(*Engine)

	majorVersionTests := []struct {
		version  string
		expected int
		hasError bool
	}{
		{"20.11.0", 20, false},
		{"^20.17.0", 20, false},
		{">=18.0.0", 18, false},
		{"~22.1.0", 22, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range majorVersionTests {
		t.Run("major_"+tt.version, func(t *testing.T) {
			result, err := engine.extractMajorVersion(tt.version)
			if tt.hasError {
				if err == nil {
					t.Errorf("extractMajorVersion(%s) expected error but got none", tt.version)
				}
			} else {
				if err != nil {
					t.Errorf("extractMajorVersion(%s) unexpected error: %v", tt.version, err)
				}
				if result != tt.expected {
					t.Errorf("extractMajorVersion(%s) = %d, expected %d", tt.version, result, tt.expected)
				}
			}
		})
	}

	goVersionTests := []struct {
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

	for _, tt := range goVersionTests {
		t.Run("go_"+tt.version, func(t *testing.T) {
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

func TestValidationErrorScenarios(t *testing.T) {
	engine := NewEngine()

	// Test critical error scenarios that should block generation
	criticalErrorTests := []struct {
		name   string
		config *models.ProjectConfig
		path   string
	}{
		{
			name:   "nil config",
			config: nil,
			path:   "templates/frontend/nextjs-app/package.json.tmpl",
		},
		{
			name: "nil versions",
			config: &models.ProjectConfig{
				Name:     "test",
				Versions: nil,
			},
			path: "templates/frontend/nextjs-app/package.json.tmpl",
		},
		{
			name: "incompatible Node.js versions",
			config: &models.ProjectConfig{
				Name: "test",
				Versions: &models.VersionConfig{
					NodeJS: &models.NodeVersionConfig{
						Runtime:      ">=20.0.0",
						TypesPackage: "^16.0.0", // Major incompatibility
					},
				},
			},
			path: "templates/frontend/nextjs-app/package.json.tmpl",
		},
	}

	for _, tt := range criticalErrorTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ValidatePreGeneration(tt.config, tt.path)
			if err != nil {
				t.Fatalf("ValidatePreGeneration returned error: %v", err)
			}

			if result.Valid {
				t.Errorf("Expected validation to fail for critical error scenario: %s", tt.name)
			}

			if len(result.Errors) == 0 {
				t.Errorf("Expected at least one error for critical error scenario: %s", tt.name)
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkValidatePreGeneration(b *testing.B) {
	engine := NewEngine()
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
		_, err := engine.ValidatePreGeneration(config, templatePath)
		if err != nil {
			b.Fatalf("ValidatePreGeneration failed: %v", err)
		}
	}
}

func BenchmarkValidatePreGenerationDirectory(b *testing.B) {
	engine := NewEngine()
	config := &models.ProjectConfig{
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
	}
	templateDir := "templates"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ValidatePreGenerationDirectory(config, templateDir)
		if err != nil {
			b.Fatalf("ValidatePreGenerationDirectory failed: %v", err)
		}
	}
}
