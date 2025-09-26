package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine_ValidateProject(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedIssues int
	}{
		{
			name: "valid project structure",
			setupProject: func(projectPath string) error {
				// Create basic project structure
				dirs := []string{"frontend", "backend", "docs"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"README.md": "# Test Project",
					"LICENSE":   "MIT License\n\nCopyright (c) 2024 Test Project",
					"go.mod":    "module test\n\ngo 1.21",
					"frontend/package.json": `{
						"name": "test-frontend",
						"version": "1.0.0"
					}`,
					"Dockerfile": "FROM node:20\nWORKDIR /app\nCOPY . .",
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
						return err
					}
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
		{
			name: "missing project directory",
			setupProject: func(projectPath string) error {
				// Don't create the directory
				return nil
			},
			expectedValid:  false,
			expectedIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")

			if err := tt.setupProject(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result, err := engine.ValidateProject(projectPath)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Issues, tt.expectedIssues)
		})
	}
}

func TestEngine_ValidatePackageJSON(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid package.json",
			content: `{
				"name": "test-package",
				"version": "1.0.0"
			}`,
			expectError: false,
		},
		{
			name: "missing required fields",
			content: `{
				"description": "test package"
			}`,
			expectError: true,
		},
		{
			name: "invalid JSON",
			content: `{
				"name": "test-package"
				"version": "1.0.0"
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "package.json")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidatePackageJSON(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEngine_ValidateGoMod(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid go.mod",
			content: `module test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)`,
			expectError: false,
		},
		{
			name:        "missing module declaration",
			content:     `go 1.21`,
			expectError: true,
		},
		{
			name:        "missing go version",
			content:     `module test-project`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "go.mod")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidateGoMod(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEngine_ValidateDockerfile(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid Dockerfile",
			content: `FROM node:20
WORKDIR /app
COPY package.json .
RUN npm install
COPY . .
CMD ["npm", "start"]`,
			expectError: false,
		},
		{
			name: "missing FROM instruction",
			content: `WORKDIR /app
COPY . .`,
			expectError: true,
		},
		{
			name: "missing WORKDIR instruction",
			content: `FROM node:20
COPY . .`,
			expectError: true,
		},
		{
			name: "missing COPY instruction",
			content: `FROM node:20
WORKDIR /app`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "Dockerfile")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidateDockerfile(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEngine_ValidateYAML(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid YAML",
			content: `name: test
version: 1.0.0
config:
  enabled: true`,
			expectError: false,
		},
		{
			name: "invalid YAML",
			content: `name: test
version: 1.0.0
config:
  enabled: true
  invalid: [`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "test.yaml")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidateYAML(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEngine_ValidateJSON(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid JSON",
			content: `{
				"name": "test",
				"version": "1.0.0",
				"config": {
					"enabled": true
				}
			}`,
			expectError: false,
		},
		{
			name: "invalid JSON",
			content: `{
				"name": "test",
				"version": "1.0.0",
				"config": {
					"enabled": true
				}
			`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "test.json")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidateJSON(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEngine_ValidateTemplate(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name: "valid template",
			content: `Hello {{.Name}}!
Version: {{.Version}}
{{if .Enabled}}Feature is enabled{{end}}`,
			expectError: false,
		},
		{
			name: "invalid template - no template syntax",
			content: `Hello World!
This is a regular text file.`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "test.tmpl")

			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			err := engine.ValidateTemplate(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateProjectStructure(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedIssues int
	}{
		{
			name: "valid structure with all components",
			setupProject: func(projectPath string) error {
				dirs := []string{"frontend", "backend", "mobile", "infrastructure"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"README.md": "# Test Project",
					"LICENSE":   "MIT License\n\nCopyright (c) 2024 Test Project",
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
		{
			name: "minimal valid structure",
			setupProject: func(projectPath string) error {
				// Create the directory
				if err := os.MkdirAll(projectPath, 0755); err != nil {
					return err
				}

				// Create required files
				files := map[string]string{
					"README.md": "# Test Project",
					"LICENSE":   "MIT License\n\nCopyright (c) 2024 Test Project",
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  true,
			expectedIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")

			if err := tt.setupProject(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			result, err := engine.ValidateProject(projectPath)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Issues, tt.expectedIssues)
		})
	}
}

// Enhanced tests from engine_enhanced_test.go

func TestEngine_ValidateProjectStructure(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Test with valid directory
	result, err := engine.ValidateProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("ValidateProjectStructure failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Test with non-existent directory
	_, err = engine.ValidateProjectStructure("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestEngine_ValidateProjectDependencies(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with package.json
	tempDir := t.TempDir()
	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"dependencies": {
			"react": "^18.0.0",
			"lodash": "^4.17.21"
		}
	}`

	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	result, err := engine.ValidateProjectDependencies(tempDir)
	if err != nil {
		t.Fatalf("ValidateProjectDependencies failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	if !result.Valid {
		t.Error("Expected valid result for basic package.json")
	}
}

func TestEngine_ValidateProjectSecurity(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with test files
	tempDir := t.TempDir()

	// Create a file with potential secret
	secretFile := `
const config = {
	apiKey: "sk-1234567890abcdef",
	password: "secret123"
};
`
	err := os.WriteFile(filepath.Join(tempDir, "config.js"), []byte(secretFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create config.js: %v", err)
	}

	result, err := engine.ValidateProjectSecurity(tempDir)
	if err != nil {
		t.Fatalf("ValidateProjectSecurity failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected security validation result, got nil")
	}

	// Should detect secrets
	if result.Summary.SecretsFound == 0 {
		t.Error("Expected to find secrets in test file")
	}
}

func TestEngine_ValidateProjectQuality(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with test files
	tempDir := t.TempDir()

	// Create a simple Go file
	goFile := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

func complexFunction() {
	if true {
		if true {
			if true {
				if true {
					if true {
						fmt.Println("Very nested")
					}
				}
			}
		}
	}
}
`
	err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	result, err := engine.ValidateProjectQuality(tempDir)
	if err != nil {
		t.Fatalf("ValidateProjectQuality failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected quality validation result, got nil")
	}

	if result.Summary.QualityScore < 0 || result.Summary.QualityScore > 100 {
		t.Errorf("Expected quality score between 0-100, got %f", result.Summary.QualityScore)
	}
}

func TestEngine_ValidateConfiguration(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectValid bool
	}{
		{
			name: "valid config",
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "Test description",
				License:      "MIT",
			},
			expectValid: true,
		},
		{
			name:        "nil config",
			config:      nil,
			expectValid: false,
		},
		{
			name: "empty name",
			config: &models.ProjectConfig{
				Name:         "",
				Organization: "test-org",
				License:      "MIT",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ValidateConfiguration(tt.config)
			if err != nil {
				t.Fatalf("ValidateConfiguration failed: %v", err)
			}

			if result == nil {
				t.Fatal("Expected validation result, got nil")
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}
		})
	}
}
