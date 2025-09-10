package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine_ValidateProject(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid project structure",
			setupProject: func(projectPath string) error {
				// Create basic project structure
				dirs := []string{"frontend", "backend", "docs", ".github/workflows"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"README.md": "# Test Project",
					"Makefile":  "all:\n\techo 'test'",
					"go.mod":    "module test\n\ngo 1.21",
					"frontend/package.json": `{
						"name": "test-frontend",
						"version": "1.0.0",
						"scripts": {
							"dev": "next dev"
						}
					}`,
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
			expectedErrors: 0,
		},
		{
			name: "missing project directory",
			setupProject: func(projectPath string) error {
				// Don't create the directory
				return nil
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "invalid package.json",
			setupProject: func(projectPath string) error {
				if err := os.MkdirAll(filepath.Join(projectPath, "frontend"), 0755); err != nil {
					return err
				}

				// Create invalid package.json
				invalidJSON := `{
					"name": "test frontend",
					"version": "invalid-version"
				}`
				return os.WriteFile(filepath.Join(projectPath, "frontend", "package.json"), []byte(invalidJSON), 0644)
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "invalid go.mod",
			setupProject: func(projectPath string) error {
				if err := os.MkdirAll(filepath.Join(projectPath, "backend"), 0755); err != nil {
					return err
				}

				// Create invalid go.mod (missing module declaration)
				invalidGoMod := "go 1.21"
				return os.WriteFile(filepath.Join(projectPath, "backend", "go.mod"), []byte(invalidGoMod), 0644)
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-project-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate project
			result, err := engine.ValidateProject(tempDir)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
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
				"name": "test-project",
				"version": "1.0.0",
				"scripts": {
					"dev": "next dev",
					"build": "next build"
				},
				"dependencies": {
					"next": "^14.0.0",
					"react": "^18.0.0"
				}
			}`,
			expectError: false,
		},
		{
			name: "missing required fields",
			content: `{
				"description": "A test project"
			}`,
			expectError: true,
		},
		{
			name: "invalid name format",
			content: `{
				"name": "Test Project",
				"version": "1.0.0",
				"scripts": {}
			}`,
			expectError: true,
		},
		{
			name: "invalid version format",
			content: `{
				"name": "test-project",
				"version": "invalid",
				"scripts": {}
			}`,
			expectError: true,
		},
		{
			name: "invalid JSON",
			content: `{
				"name": "test-project"
				"version": "1.0.0"
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "package-*.json")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			err = engine.ValidatePackageJSON(tempFile.Name())
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
			content: `module github.com/example/project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.0
)`,
			expectError: false,
		},
		{
			name: "missing module declaration",
			content: `go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)`,
			expectError: true,
		},
		{
			name: "missing go version",
			content: `module github.com/example/project

require (
	github.com/gin-gonic/gin v1.9.1
)`,
			expectError: true,
		},
		{
			name: "invalid go version",
			content: `module github.com/example/project

go invalid

require (
	github.com/gin-gonic/gin v1.9.1
)`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "go-*.mod")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			err = engine.ValidateGoMod(tempFile.Name())
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
			content: `FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

EXPOSE 3000
CMD ["npm", "start"]`,
			expectError: false,
		},
		{
			name: "missing FROM instruction",
			content: `WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
			expectError: true,
		},
		{
			name: "missing WORKDIR instruction",
			content: `FROM node:18-alpine
COPY . .
CMD ["npm", "start"]`,
			expectError: true,
		},
		{
			name: "missing COPY instruction",
			content: `FROM node:18-alpine
WORKDIR /app
CMD ["npm", "start"]`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "Dockerfile-*")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			err = engine.ValidateDockerfile(tempFile.Name())
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
			content: `name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'`,
			expectError: false,
		},
		{
			name: "invalid YAML",
			content: `name: CI
on:
  push:
    branches: [main
  pull_request:
    branches: [main]`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "test-*.yml")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			err = engine.ValidateYAML(tempFile.Name())
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
				"nested": {
					"array": [1, 2, 3],
					"boolean": true
				}
			}`,
			expectError: false,
		},
		{
			name: "invalid JSON",
			content: `{
				"name": "test",
				"version": "1.0.0"
				"invalid": true
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "test-*.json")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			err = engine.ValidateJSON(tempFile.Name())
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateProjectStructure(t *testing.T) {
	engine := &Engine{}

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid structure with all components",
			setupProject: func(projectPath string) error {
				dirs := []string{"frontend", "backend", "mobile", "deploy", "docs", ".github"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				files := []string{"README.md", "Makefile"}
				for _, file := range files {
					if err := os.WriteFile(filepath.Join(projectPath, file), []byte("test"), 0644); err != nil {
						return err
					}
				}

				return nil
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "missing main components",
			setupProject: func(projectPath string) error {
				// Only create docs and deploy directories
				dirs := []string{"docs", "deploy"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "minimal valid structure",
			setupProject: func(projectPath string) error {
				// Create only frontend directory
				if err := os.MkdirAll(filepath.Join(projectPath, "frontend"), 0755); err != nil {
					return err
				}
				return nil
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-structure-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate structure
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err = engine.validateProjectStructure(tempDir, result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestValidateDependencyCompatibility(t *testing.T) {
	engine := &Engine{}

	tests := []struct {
		name             string
		setupProject     func(string) error
		expectedWarnings int
	}{
		{
			name: "single package.json",
			setupProject: func(projectPath string) error {
				packageJSON := `{
					"name": "test-project",
					"version": "1.0.0",
					"scripts": {}
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedWarnings: 0,
		},
		{
			name: "multiple package.json files",
			setupProject: func(projectPath string) error {
				// Create frontend package.json
				if err := os.MkdirAll(filepath.Join(projectPath, "frontend"), 0755); err != nil {
					return err
				}
				frontendPackageJSON := `{
					"name": "frontend",
					"version": "1.0.0",
					"scripts": {}
				}`
				if err := os.WriteFile(filepath.Join(projectPath, "frontend", "package.json"), []byte(frontendPackageJSON), 0644); err != nil {
					return err
				}

				// Create admin package.json
				if err := os.MkdirAll(filepath.Join(projectPath, "admin"), 0755); err != nil {
					return err
				}
				adminPackageJSON := `{
					"name": "admin",
					"version": "1.0.0",
					"scripts": {}
				}`
				return os.WriteFile(filepath.Join(projectPath, "admin", "package.json"), []byte(adminPackageJSON), 0644)
			},
			expectedWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-deps-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate dependencies
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err = engine.validateDependencyCompatibility(tempDir, result)
			require.NoError(t, err)

			assert.Len(t, result.Warnings, tt.expectedWarnings)
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	engine := &Engine{}

	t.Run("isValidSemVer", func(t *testing.T) {
		tests := []struct {
			version string
			valid   bool
		}{
			{"1.0.0", true},
			{"0.1.2", true},
			{"10.20.30", true},
			{"1.0", false},
			{"1.0.0.0", false},
			{"1.0.a", false},
			{"", false},
			{"v1.0.0", false}, // This implementation doesn't handle 'v' prefix
		}

		for _, tt := range tests {
			result := engine.isValidSemVer(tt.version)
			assert.Equal(t, tt.valid, result, "version: %s", tt.version)
		}
	})

	t.Run("isValidGoVersion", func(t *testing.T) {
		tests := []struct {
			version string
			valid   bool
		}{
			{"1.21", true},
			{"1.21.0", true},
			{"1.20.5", true},
			{"1", false},
			{"1.21.0.0", false},
			{"1.a", false},
			{"", false},
		}

		for _, tt := range tests {
			result := engine.isValidGoVersion(tt.version)
			assert.Equal(t, tt.valid, result, "version: %s", tt.version)
		}
	})
}

// Helper function to create a test project structure
func createTestProject(tb testing.TB, projectPath string) {
	// Create directories
	dirs := []string{
		"frontend/src/components",
		"backend/internal/handlers",
		"mobile/android/app",
		"mobile/ios/App",
		"deploy/kubernetes",
		"docs",
		".github/workflows",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectPath, dir), 0755)
		require.NoError(tb, err)
	}

	// Create files
	files := map[string]string{
		"README.md": "# Test Project\n\nThis is a test project.",
		"Makefile":  "all:\n\techo 'Building project...'\n",
		"frontend/package.json": `{
			"name": "test-frontend",
			"version": "1.0.0",
			"scripts": {
				"dev": "next dev",
				"build": "next build"
			},
			"dependencies": {
				"next": "^14.0.0",
				"react": "^18.0.0"
			}
		}`,
		"backend/go.mod": `module github.com/example/test-backend

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.0
)`,
		"backend/Dockerfile": `FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080
CMD ["./main"]`,
		"docker-compose.yml": `version: '3.8'
services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
  backend:
    build: ./backend
    ports:
      - "8080:8080"`,
		".github/workflows/ci.yml": `name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectPath, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(tb, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(tb, err)
	}
}

// Benchmark tests
func BenchmarkValidateProject(b *testing.B) {
	engine := NewEngine()

	// Create a test project
	tempDir, err := os.MkdirTemp("", "benchmark-project-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	createTestProject(b, tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ValidateProject(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkValidatePackageJSON(b *testing.B) {
	engine := NewEngine()

	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"scripts": {
			"dev": "next dev",
			"build": "next build",
			"test": "jest"
		},
		"dependencies": {
			"next": "^14.0.0",
			"react": "^18.0.0",
			"typescript": "^5.0.0"
		},
		"devDependencies": {
			"jest": "^29.0.0",
			"@types/node": "^20.0.0"
		}
	}`

	tempFile, err := os.CreateTemp("", "benchmark-package-*.json")
	require.NoError(b, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(packageJSON)
	require.NoError(b, err)
	tempFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := engine.ValidatePackageJSON(tempFile.Name())
		require.NoError(b, err)
	}
}
