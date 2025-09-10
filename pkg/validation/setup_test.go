package validation

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupEngine_SetupProject(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(30 * time.Second) // Shorter timeout for tests

	tests := []struct {
		name           string
		setupProject   func(string) error
		config         *models.ProjectConfig
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "frontend project setup",
			setupProject: func(projectPath string) error {
				// Create frontend directory with package.json
				frontendPath := filepath.Join(projectPath, "frontend")
				if err := os.MkdirAll(frontendPath, 0755); err != nil {
					return err
				}

				packageJSON := `{
					"name": "test-frontend",
					"version": "1.0.0",
					"scripts": {
						"build": "echo 'Building frontend...'",
						"dev": "echo 'Starting dev server...'"
					},
					"dependencies": {}
				}`
				return os.WriteFile(filepath.Join(frontendPath, "package.json"), []byte(packageJSON), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Frontend: models.FrontendComponents{
						MainApp: true,
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "backend project setup",
			setupProject: func(projectPath string) error {
				// Create backend directory with go.mod
				backendPath := filepath.Join(projectPath, "backend")
				if err := os.MkdirAll(backendPath, 0755); err != nil {
					return err
				}

				goMod := `module test-backend

go 1.21`
				if err := os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644); err != nil {
					return err
				}

				mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`
				return os.WriteFile(filepath.Join(backendPath, "main.go"), []byte(mainGo), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Backend: models.BackendComponents{
						API: true,
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "mobile project setup",
			setupProject: func(projectPath string) error {
				// Create mobile directories
				androidPath := filepath.Join(projectPath, "mobile", "android")
				if err := os.MkdirAll(androidPath, 0755); err != nil {
					return err
				}

				iosPath := filepath.Join(projectPath, "mobile", "ios")
				if err := os.MkdirAll(iosPath, 0755); err != nil {
					return err
				}

				// Create gradlew script
				gradlew := `#!/bin/bash
echo "Gradle wrapper executed"`
				if err := os.WriteFile(filepath.Join(androidPath, "gradlew"), []byte(gradlew), 0755); err != nil {
					return err
				}

				// Create Podfile
				podfile := `platform :ios, '14.0'
target 'TestApp' do
  use_frameworks!
end`
				return os.WriteFile(filepath.Join(iosPath, "Podfile"), []byte(podfile), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Mobile: models.MobileComponents{
						Android: true,
						IOS:     true,
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "infrastructure project setup",
			setupProject: func(projectPath string) error {
				// Create Terraform files
				mainTf := `terraform {
  required_version = ">= 1.0"
}

provider "aws" {
  region = "us-west-2"
}`
				if err := os.WriteFile(filepath.Join(projectPath, "main.tf"), []byte(mainTf), 0644); err != nil {
					return err
				}

				variablesTf := `variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}`
				return os.WriteFile(filepath.Join(projectPath, "variables.tf"), []byte(variablesTf), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Infrastructure: models.InfrastructureComponents{
						Terraform: true,
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
			tempDir, err := os.MkdirTemp("", "test-setup-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Run setup
			result, err := engine.SetupProject(tempDir, tt.config)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestSetupEngine_VerifyProject(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(30 * time.Second) // Shorter timeout for tests

	tests := []struct {
		name           string
		setupProject   func(string) error
		config         *models.ProjectConfig
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "frontend project verification",
			setupProject: func(projectPath string) error {
				// Create frontend directory with package.json and build script
				frontendPath := filepath.Join(projectPath, "frontend")
				if err := os.MkdirAll(frontendPath, 0755); err != nil {
					return err
				}

				packageJSON := `{
					"name": "test-frontend",
					"version": "1.0.0",
					"scripts": {
						"build": "echo 'Build successful' > build.log",
						"dev": "echo 'Starting dev server...'"
					},
					"dependencies": {}
				}`
				return os.WriteFile(filepath.Join(frontendPath, "package.json"), []byte(packageJSON), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Frontend: models.FrontendComponents{
						MainApp: true,
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "backend project verification",
			setupProject: func(projectPath string) error {
				// Create backend directory with valid Go project
				backendPath := filepath.Join(projectPath, "backend")
				if err := os.MkdirAll(backendPath, 0755); err != nil {
					return err
				}

				goMod := `module test-backend

go 1.21`
				if err := os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644); err != nil {
					return err
				}

				mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`
				return os.WriteFile(filepath.Join(backendPath, "main.go"), []byte(mainGo), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Backend: models.BackendComponents{
						API: true,
					},
				},
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "invalid backend project",
			setupProject: func(projectPath string) error {
				// Create backend directory with invalid Go code
				backendPath := filepath.Join(projectPath, "backend")
				if err := os.MkdirAll(backendPath, 0755); err != nil {
					return err
				}

				goMod := `module test-backend

go 1.21`
				if err := os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644); err != nil {
					return err
				}

				// Invalid Go code
				mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!"
	// Missing closing parenthesis
}`
				return os.WriteFile(filepath.Join(backendPath, "main.go"), []byte(mainGo), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Backend: models.BackendComponents{
						API: true,
					},
				},
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "infrastructure project verification",
			setupProject: func(projectPath string) error {
				// Create only Dockerfile (docker-compose might not be available in test environment)
				dockerfile := `FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 3000
CMD ["npm", "start"]`
				return os.WriteFile(filepath.Join(projectPath, "Dockerfile"), []byte(dockerfile), 0644)
			},
			config: &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Components: models.Components{
					Infrastructure: models.InfrastructureComponents{
						Docker: true,
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
			tempDir, err := os.MkdirTemp("", "test-verify-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Run verification
			result, err := engine.VerifyProject(tempDir, tt.config)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestSetupEngine_SetupFrontendComponents(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(10 * time.Second)

	tests := []struct {
		name           string
		setupProject   func(string) error
		config         *models.ProjectConfig
		expectedErrors int
	}{
		{
			name: "valid frontend setup",
			setupProject: func(projectPath string) error {
				frontendPath := filepath.Join(projectPath, "frontend")
				if err := os.MkdirAll(frontendPath, 0755); err != nil {
					return err
				}

				packageJSON := `{
					"name": "test-frontend",
					"version": "1.0.0",
					"scripts": {},
					"dependencies": {}
				}`
				return os.WriteFile(filepath.Join(frontendPath, "package.json"), []byte(packageJSON), 0644)
			},
			config: &models.ProjectConfig{
				Components: models.Components{
					Frontend: models.FrontendComponents{
						MainApp: true,
					},
				},
			},
			expectedErrors: 0,
		},
		{
			name: "missing frontend directory",
			setupProject: func(projectPath string) error {
				// Don't create frontend directory
				return nil
			},
			config: &models.ProjectConfig{
				Components: models.Components{
					Frontend: models.FrontendComponents{
						MainApp: true,
					},
				},
			},
			expectedErrors: 0, // Should only generate warnings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-frontend-setup-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Run frontend setup
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err = engine.setupFrontendComponents(tempDir, tt.config, result)
			require.NoError(t, err)

			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestSetupEngine_SetupBackendComponents(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(10 * time.Second)

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedErrors int
	}{
		{
			name: "valid backend setup",
			setupProject: func(projectPath string) error {
				backendPath := filepath.Join(projectPath, "backend")
				if err := os.MkdirAll(backendPath, 0755); err != nil {
					return err
				}

				goMod := `module test-backend

go 1.21`
				return os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedErrors: 0,
		},
		{
			name: "missing backend directory",
			setupProject: func(projectPath string) error {
				// Don't create backend directory
				return nil
			},
			expectedErrors: 0, // Should only generate warnings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-backend-setup-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Run backend setup
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			config := &models.ProjectConfig{
				Components: models.Components{
					Backend: models.BackendComponents{
						API: true,
					},
				},
			}

			err = engine.setupBackendComponents(tempDir, config, result)
			require.NoError(t, err)

			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestSetupEngine_RunCommand(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(5 * time.Second)

	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "successful command",
			command:     "echo",
			args:        []string{"hello"},
			expectError: false,
		},
		{
			name:        "failing command",
			command:     "false", // Command that always fails
			args:        []string{},
			expectError: true,
		},
		{
			name:        "non-existent command",
			command:     "non-existent-command-12345",
			args:        []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-command-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Run command
			err = engine.runCommand(tempDir, tt.command, tt.args...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetupEngine_Timeout(t *testing.T) {
	engine := NewSetupEngine()
	engine.SetTimeout(1 * time.Second) // Very short timeout

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "test-timeout-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Run a command that should timeout (sleep for longer than timeout)
	err = engine.runCommand(tempDir, "sleep", "5")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestSetupEngine_SetTimeout(t *testing.T) {
	engine := NewSetupEngine()

	// Test setting timeout
	newTimeout := 2 * time.Minute
	engine.SetTimeout(newTimeout)
	assert.Equal(t, newTimeout, engine.timeout)
}

// Benchmark tests
func BenchmarkSetupEngine_SetupProject(b *testing.B) {
	engine := NewSetupEngine()
	engine.SetTimeout(30 * time.Second)

	// Create a test project
	tempDir, err := os.MkdirTemp("", "benchmark-setup-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create frontend directory with package.json
	frontendPath := filepath.Join(tempDir, "frontend")
	err = os.MkdirAll(frontendPath, 0755)
	require.NoError(b, err)

	packageJSON := `{
		"name": "test-frontend",
		"version": "1.0.0",
		"scripts": {
			"build": "echo 'Building...'",
			"dev": "echo 'Dev server...'"
		},
		"dependencies": {}
	}`
	err = os.WriteFile(filepath.Join(frontendPath, "package.json"), []byte(packageJSON), 0644)
	require.NoError(b, err)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.SetupProject(tempDir, config)
		require.NoError(b, err)
	}
}

func BenchmarkSetupEngine_VerifyProject(b *testing.B) {
	engine := NewSetupEngine()
	engine.SetTimeout(30 * time.Second)

	// Create a test project
	tempDir, err := os.MkdirTemp("", "benchmark-verify-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create backend directory with Go project
	backendPath := filepath.Join(tempDir, "backend")
	err = os.MkdirAll(backendPath, 0755)
	require.NoError(b, err)

	goMod := `module test-backend

go 1.21`
	err = os.WriteFile(filepath.Join(backendPath, "go.mod"), []byte(goMod), 0644)
	require.NoError(b, err)

	mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`
	err = os.WriteFile(filepath.Join(backendPath, "main.go"), []byte(mainGo), 0644)
	require.NoError(b, err)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Backend: models.BackendComponents{
				API: true,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.VerifyProject(tempDir, config)
		require.NoError(b, err)
	}
}
