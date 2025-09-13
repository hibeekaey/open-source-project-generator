package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectTypeValidator_ValidateFrontendProject(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid frontend project",
			setupProject: func(projectPath string) error {
				// Create required directories
				dirs := []string{"src/app", "src/components", "public"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"package.json": `{
						"name": "test-frontend",
						"version": "1.0.0",
						"scripts": {
							"dev": "next dev",
							"build": "next build",
							"start": "next start"
						},
						"dependencies": {
							"next": "^14.0.0",
							"react": "^18.0.0",
							"react-dom": "^18.0.0"
						}
					}`,
					"next.config.js": `/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
}

module.exports = nextConfig`,
					"tailwind.config.js": `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`,
					"tsconfig.json": `{
						"compilerOptions": {
							"target": "es5",
							"lib": ["dom", "dom.iterable", "es6"],
							"allowJs": true,
							"skipLibCheck": true,
							"strict": true,
							"forceConsistentCasingInFileNames": true,
							"noEmit": true,
							"esModuleInterop": true,
							"module": "esnext",
							"moduleResolution": "node",
							"resolveJsonModule": true,
							"isolatedModules": true,
							"jsx": "preserve",
							"incremental": true,
							"plugins": [
								{
									"name": "next"
								}
							],
							"paths": {
								"@/*": ["./src/*"]
							}
						},
						"include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
						"exclude": ["node_modules"]
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
			name: "missing required files",
			setupProject: func(projectPath string) error {
				// Only create package.json
				packageJSON := `{
					"name": "test-frontend",
					"version": "1.0.0",
					"scripts": {}
				}`
				return os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644)
			},
			expectedValid:  false,
			expectedErrors: 3, // missing next.config.js, tailwind.config.js, tsconfig.json
		},
		{
			name: "invalid package.json",
			setupProject: func(projectPath string) error {
				// Create all required files but with invalid package.json
				files := map[string]string{
					"package.json":       `{"name": "test"}`, // missing required fields
					"next.config.js":     `module.exports = {}`,
					"tailwind.config.js": `module.exports = {}`,
					"tsconfig.json":      `{}`,
				}

				for filePath, content := range files {
					if err := os.WriteFile(filepath.Join(projectPath, filePath), []byte(content), 0644); err != nil {
						return err
					}
				}

				return nil
			},
			expectedValid:  true, // Structure is valid, just missing recommended dependencies
			expectedErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-frontend-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate project
			result, err := validator.ValidateFrontendProject(tempDir)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestProjectTypeValidator_ValidateBackendProject(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid backend project",
			setupProject: func(projectPath string) error {
				// Create recommended directories
				dirs := []string{
					"internal/handlers",
					"internal/models",
					"internal/services",
					"pkg",
					"cmd",
				}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"go.mod": `module github.com/example/test-backend

go 1.24

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.0
)`,
					"main.go": `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}`,
					"Dockerfile": `FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]`,
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
			name: "missing required files",
			setupProject: func(projectPath string) error {
				// Only create go.mod
				goMod := `module test

go 1.24`
				return os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goMod), 0644)
			},
			expectedValid:  false,
			expectedErrors: 2, // missing main.go and Dockerfile
		},
		{
			name: "invalid main.go",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"go.mod": `module test

go 1.24`,
					"main.go": `package notmain

func notmain() {
	println("hello")
}`,
					"Dockerfile": `FROM golang:1.24-alpine
WORKDIR /app
COPY . .
CMD ["go", "run", "main.go"]`,
				}

				for filePath, content := range files {
					if err := os.WriteFile(filepath.Join(projectPath, filePath), []byte(content), 0644); err != nil {
						return err
					}
				}

				return nil
			},
			expectedValid:  false,
			expectedErrors: 2, // invalid package declaration and missing main function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-backend-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate project
			result, err := validator.ValidateBackendProject(tempDir)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestProjectTypeValidator_ValidateMobileProject(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid mobile project with both platforms",
			setupProject: func(projectPath string) error {
				// Create Android structure
				androidDirs := []string{
					"android/app/src/main",
					"android/gradle/wrapper",
				}
				for _, dir := range androidDirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create iOS structure
				iosDirs := []string{
					"ios/TestApp.xcodeproj",
					"ios/TestApp",
				}
				for _, dir := range iosDirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create required files
				files := map[string]string{
					"android/build.gradle.kts":                 `// Top-level build file`,
					"android/settings.gradle.kts":              `rootProject.name = "TestApp"`,
					"android/app/build.gradle.kts":             `plugins { id("com.android.application") }`,
					"android/app/src/main/AndroidManifest.xml": `<?xml version="1.0" encoding="utf-8"?><manifest></manifest>`,
					"ios/TestApp.xcodeproj/project.pbxproj":    `// Xcode project file`,
					"ios/Podfile":                              `platform :ios, '14.0'`,
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
			name: "valid mobile project with Android only",
			setupProject: func(projectPath string) error {
				// Create Android structure only
				androidDirs := []string{
					"android/app/src/main",
				}
				for _, dir := range androidDirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				files := map[string]string{
					"android/build.gradle.kts":                 `// Top-level build file`,
					"android/settings.gradle.kts":              `rootProject.name = "TestApp"`,
					"android/app/build.gradle.kts":             `plugins { id("com.android.application") }`,
					"android/app/src/main/AndroidManifest.xml": `<?xml version="1.0" encoding="utf-8"?><manifest></manifest>`,
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
			name: "no mobile platforms",
			setupProject: func(projectPath string) error {
				// Create empty directory
				return nil
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "incomplete Android project",
			setupProject: func(projectPath string) error {
				// Create Android directory but missing required files
				if err := os.MkdirAll(filepath.Join(projectPath, "android"), 0755); err != nil {
					return err
				}

				// Only create one required file
				buildGradle := `// Top-level build file`
				return os.WriteFile(filepath.Join(projectPath, "android", "build.gradle.kts"), []byte(buildGradle), 0644)
			},
			expectedValid:  false,
			expectedErrors: 3, // missing settings.gradle.kts, app/build.gradle.kts, AndroidManifest.xml
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-mobile-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate project
			result, err := validator.ValidateMobileProject(tempDir)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestProjectTypeValidator_ValidateInfrastructureProject(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid infrastructure project with all components",
			setupProject: func(projectPath string) error {
				// Create Kubernetes directory
				if err := os.MkdirAll(filepath.Join(projectPath, "k8s"), 0755); err != nil {
					return err
				}

				files := map[string]string{
					"Dockerfile": `FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 3000
CMD ["npm", "start"]`,
					"docker-compose.yml": `version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"`,
					"main.tf": `terraform {
  required_version = ">= 1.0"
}

provider "aws" {
  region = "us-west-2"
}`,
					"variables.tf": `variable "environment" {
  description = "Environment name"
  type        = string
  default     = "dev"
}`,
					"outputs.tf": `output "app_url" {
  description = "Application URL"
  value       = "https://example.com"
}`,
					"k8s/deployment.yaml": `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
      - name: app
        image: test-app:latest
        ports:
        - containerPort: 3000`,
					"k8s/service.yaml": `apiVersion: v1
kind: Service
metadata:
  name: test-app-service
spec:
  selector:
    app: test-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
  type: LoadBalancer`,
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
			name: "valid infrastructure project with Docker only",
			setupProject: func(projectPath string) error {
				files := map[string]string{
					"Dockerfile": `FROM node:18-alpine
WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
					"docker-compose.yml": `version: '3.8'
services:
  app:
    build: .`,
				}

				for filePath, content := range files {
					if err := os.WriteFile(filepath.Join(projectPath, filePath), []byte(content), 0644); err != nil {
						return err
					}
				}

				return nil
			},
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "no infrastructure components",
			setupProject: func(projectPath string) error {
				// Create empty directory
				return nil
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "invalid Kubernetes YAML",
			setupProject: func(projectPath string) error {
				if err := os.MkdirAll(filepath.Join(projectPath, "k8s"), 0755); err != nil {
					return err
				}

				// Create invalid YAML
				invalidYAML := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: test-app
  template:
    metadata:
      labels:
        app: test-app
    spec:
      containers:
      - name: app
        image: test-app:latest
        ports:
        - containerPort: 3000
      - invalid: yaml: structure`

				return os.WriteFile(filepath.Join(projectPath, "k8s", "deployment.yaml"), []byte(invalidYAML), 0644)
			},
			expectedValid:  false,
			expectedErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "test-infra-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup project
			err = tt.setupProject(tempDir)
			require.NoError(t, err)

			// Validate project
			result, err := validator.ValidateInfrastructureProject(tempDir)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}

func TestValidateFrontendPackageJSON(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name             string
		packageJSON      string
		expectedWarnings int
	}{
		{
			name: "complete frontend package.json",
			packageJSON: `{
				"name": "test-frontend",
				"version": "1.0.0",
				"scripts": {
					"dev": "next dev",
					"build": "next build",
					"start": "next start"
				},
				"dependencies": {
					"next": "^14.0.0",
					"react": "^18.0.0",
					"react-dom": "^18.0.0"
				}
			}`,
			expectedWarnings: 0,
		},
		{
			name: "missing frontend dependencies",
			packageJSON: `{
				"name": "test-frontend",
				"version": "1.0.0",
				"scripts": {
					"dev": "next dev"
				},
				"dependencies": {
					"lodash": "^4.17.21"
				}
			}`,
			expectedWarnings: 5, // 3 missing deps + 2 missing scripts
		},
		{
			name: "missing scripts",
			packageJSON: `{
				"name": "test-frontend",
				"version": "1.0.0",
				"scripts": {},
				"dependencies": {
					"next": "^14.0.0",
					"react": "^18.0.0",
					"react-dom": "^18.0.0"
				}
			}`,
			expectedWarnings: 3, // 3 missing scripts
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "package-*.json")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.packageJSON)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err = validator.validateFrontendPackageJSON(tempFile.Name(), result)
			require.NoError(t, err)

			assert.Len(t, result.Warnings, tt.expectedWarnings)
		})
	}
}

func TestValidateMainGo(t *testing.T) {
	validator := NewProjectTypeValidator()

	tests := []struct {
		name           string
		content        string
		expectedValid  bool
		expectedErrors int
	}{
		{
			name: "valid main.go",
			content: `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			expectedValid:  true,
			expectedErrors: 0,
		},
		{
			name: "missing package main",
			content: `package notmain

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "missing main function",
			content: `package main

import "fmt"

func notmain() {
	fmt.Println("Hello, World!")
}`,
			expectedValid:  false,
			expectedErrors: 1,
		},
		{
			name: "missing both package and main function",
			content: `package notmain

import "fmt"

func notmain() {
	fmt.Println("Hello, World!")
}`,
			expectedValid:  false,
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempFile, err := os.CreateTemp("", "main-*.go")
			require.NoError(t, err)
			defer os.Remove(tempFile.Name())

			// Write content
			_, err = tempFile.WriteString(tt.content)
			require.NoError(t, err)
			tempFile.Close()

			// Validate
			result := &models.ValidationResult{
				Valid:    true,
				Errors:   []models.ValidationError{},
				Warnings: []models.ValidationWarning{},
			}

			err = validator.validateMainGo(tempFile.Name(), result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedValid, result.Valid)
			assert.Len(t, result.Errors, tt.expectedErrors)
		})
	}
}
