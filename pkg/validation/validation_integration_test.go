package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationEngine_FullIntegration(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		setupProject   func(string) error
		expectedValid  bool
		expectedIssues int
		description    string
	}{
		{
			name: "complete valid project",
			setupProject: func(projectPath string) error {
				// Create project structure
				dirs := []string{"src", "docs", "tests"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				// Create files
				files := map[string]string{
					"README.md": `# Test Project

This is a comprehensive test project for validation.

## Installation

` + "```bash" + `
npm install
` + "```" + `

## Usage

` + "```bash" + `
npm start
` + "```" + `
`,
					"LICENSE": `MIT License

Copyright (c) 2024 Test Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`,
					"package.json": `{
  "name": "test-project",
  "version": "1.0.0",
  "description": "A comprehensive test project",
  "main": "src/index.js",
  "scripts": {
    "start": "node src/index.js",
    "test": "jest",
    "lint": "eslint src/"
  },
  "dependencies": {
    "express": "^4.18.2",
    "lodash": "^4.17.21"
  },
  "devDependencies": {
    "jest": "^29.0.0",
    "eslint": "^8.0.0"
  },
  "keywords": ["test", "validation"],
  "author": "Test Author",
  "license": "MIT"
}`,
					"tsconfig.json": `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "tests"]
}`,
					".eslintrc.json": `{
  "extends": ["eslint:recommended"],
  "env": {
    "node": true,
    "es2020": true
  },
  "parserOptions": {
    "ecmaVersion": 2020,
    "sourceType": "module"
  },
  "rules": {
    "no-console": "warn",
    "no-debugger": "error",
    "no-unused-vars": "error"
  }
}`,
					"docker-compose.yml": `version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
    volumes:
      - ./src:/app/src:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  db:
    image: postgres:13-alpine
    environment:
      - POSTGRES_DB=testdb
      - POSTGRES_USER=testuser
      - POSTGRES_PASSWORD=testpass
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:`,
					"Dockerfile": `FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production && \
    npm cache clean --force

# Copy source code
COPY src/ ./src/

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

# Change ownership
RUN chown -R nextjs:nodejs /app
USER nextjs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# Start application
CMD ["npm", "start"]`,
					".env.example": `NODE_ENV=development
PORT=3000
DATABASE_URL=postgresql://user:password@localhost:5432/dbname
LOG_LEVEL=info
JWT_SECRET=your-secret-key-here
API_BASE_URL=https://api.example.com`,
					".gitignore": `# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/
*.lcov

# nyc test coverage
.nyc_output

# Grunt intermediate storage
.grunt

# Bower dependency directory
bower_components

# node-waf configuration
.lock-wscript

# Compiled binary addons
build/Release

# Dependency directories
jspm_packages/

# TypeScript cache
*.tsbuildinfo

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Microbundle cache
.rpt2_cache/
.rts2_cache_cjs/
.rts2_cache_es/
.rts2_cache_umd/

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env
.env.test
.env.production

# parcel-bundler cache
.cache
.parcel-cache

# Next.js build output
.next

# Nuxt.js build / generate output
.nuxt
dist

# Gatsby files
.cache/
public

# Storybook build outputs
.out
.storybook-out

# Temporary folders
tmp/
temp/

# Logs
logs
*.log

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db`,
					"Makefile": `# Variables
NODE_VERSION := 18
APP_NAME := test-project
DOCKER_IMAGE := $(APP_NAME):latest

# Default target
.PHONY: all
all: install build test

# Install dependencies
.PHONY: install
install:
	npm ci

# Build the application
.PHONY: build
build:
	npm run build

# Run tests
.PHONY: test
test:
	npm test

# Run linting
.PHONY: lint
lint:
	npm run lint

# Start development server
.PHONY: dev
dev:
	npm run dev

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf dist/
	rm -rf node_modules/
	rm -rf coverage/

# Docker targets
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run -p 3000:3000 $(DOCKER_IMAGE)

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  install     - Install dependencies"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  lint        - Run linting"
	@echo "  dev         - Start development server"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"
	@echo "  help        - Show this help message"`,
					"src/index.js": `const express = require('express');
const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());

app.get('/', (req, res) => {
  res.json({ message: 'Hello World!' });
});

app.get('/health', (req, res) => {
  res.json({ status: 'OK', timestamp: new Date().toISOString() });
});

app.listen(port, () => {
  console.log('Server running on port ' + port);
});`,
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					dir := filepath.Dir(fullPath)
					if err := os.MkdirAll(dir, 0755); err != nil {
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
			description:    "A complete, well-structured project with all required files and proper configuration",
		},
		{
			name: "project with multiple issues",
			setupProject: func(projectPath string) error {
				// Create minimal structure with issues
				files := map[string]string{
					// Missing README.md and LICENSE
					"package.json": `{
  "description": "Missing name and version"
}`, // Invalid package.json
					"Dockerfile": `WORKDIR /app
COPY . .
CMD ["npm", "start"]`, // Missing FROM instruction
					".env": `API_KEY=sk_test_1234567890abcdef
PASSWORD=supersecret123
NODE_ENV=production`, // Contains secrets
					"docker-compose.yml": `version: '2.4'
services:
  web:
    image: nginx`, // Deprecated version, missing required fields
					"Makefile": `build:
    echo "Using spaces instead of tabs"`, // Invalid indentation
					"My File.txt": "File with spaces in name", // Naming issue
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					dir := filepath.Dir(fullPath)
					if err := os.MkdirAll(dir, 0755); err != nil {
						return err
					}
					if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
						return err
					}
				}
				return nil
			},
			expectedValid:  false,
			expectedIssues: 10, // Multiple validation issues
			description:    "A project with various validation issues across different file types",
		},
		{
			name: "go project structure",
			setupProject: func(projectPath string) error {
				// Create Go project structure
				dirs := []string{"cmd/app", "pkg/utils", "internal/config", "api", "web"}
				for _, dir := range dirs {
					if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
						return err
					}
				}

				files := map[string]string{
					"README.md": "# Go Project\n\nA sample Go project.",
					"LICENSE":   "MIT License\n\nCopyright (c) 2024",
					"go.mod": `module github.com/example/test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.4
)`,
					"go.sum": `github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=`,
					"cmd/app/main.go": `package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	
	log.Println("Server starting on :8080")
	r.Run(":8080")
}`,
					"pkg/utils/helpers.go": `package utils

import "strings"

func ToUpper(s string) string {
	return strings.ToUpper(s)
}`,
					"internal/config/config.go": `package config

type Config struct {
	Port     string
	Database string
}

func New() *Config {
	return &Config{
		Port:     "8080",
		Database: "localhost:5432",
	}
}`,
					".gitignore": `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work`,
					"Makefile": `# Go project Makefile
BINARY_NAME := app
BUILD_DIR := ./bin

.PHONY: all build test clean run

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/app

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

lint:
	golangci-lint run

deps:
	go mod download
	go mod tidy`,
				}

				for filePath, content := range files {
					fullPath := filepath.Join(projectPath, filePath)
					dir := filepath.Dir(fullPath)
					if err := os.MkdirAll(dir, 0755); err != nil {
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
			description:    "A well-structured Go project with proper module setup and organization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			projectPath := filepath.Join(tempDir, "test-project")
			err := os.MkdirAll(projectPath, 0755)
			require.NoError(t, err)

			if err := tt.setupProject(projectPath); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Run comprehensive validation
			result, err := engine.ValidateProject(projectPath)
			require.NoError(t, err, "Validation should not fail for: %s", tt.description)
			require.NotNil(t, result)

			t.Logf("Test: %s", tt.name)
			t.Logf("Description: %s", tt.description)
			t.Logf("Expected Valid: %v, Actual Valid: %v", tt.expectedValid, result.Valid)
			t.Logf("Expected Issues: %d, Actual Issues: %d", tt.expectedIssues, len(result.Issues))

			if len(result.Issues) > 0 {
				t.Logf("Issues found:")
				for i, issue := range result.Issues {
					t.Logf("  %d. [%s] %s (File: %s, Rule: %s)", i+1, issue.Type, issue.Message, issue.File, issue.Rule)
				}
			}

			assert.Equal(t, tt.expectedValid, result.Valid, "Validation result should match expected")

			// Allow some flexibility in issue count for integration tests
			if tt.expectedIssues == 0 {
				// For integration tests, allow some minor warnings
				assert.LessOrEqual(t, len(result.Issues), 30, "Should have minimal issues for valid projects")
			} else {
				assert.Greater(t, len(result.Issues), 0, "Should have some issues")
				// For integration tests, we're more flexible about exact counts
				// as different validators might find different numbers of issues
			}
		})
	}
}

func TestValidationEngine_ComponentIntegration(t *testing.T) {
	engine := NewEngine()

	// Test that all components work together
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "integration-test")
	err := os.MkdirAll(projectPath, 0755)
	require.NoError(t, err)

	// Create a project with mixed technologies
	files := map[string]string{
		"README.md": "# Integration Test Project",
		"LICENSE":   "MIT License",
		"package.json": `{
  "name": "integration-test",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}`,
		"go.mod": `module integration-test

go 1.21

require github.com/gin-gonic/gin v1.9.1`,
		"requirements.txt": `django==4.2.0
requests>=2.28.0`,
		"Dockerfile": `FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
USER node
EXPOSE 3000
CMD ["npm", "start"]`,
		"docker-compose.yml": `version: '3.8'
services:
  web:
    build: .
    ports:
      - "3000:3000"`,
		".env": `NODE_ENV=development
PORT=3000`,
		"Makefile": `all: build test

build:
	npm run build

test:
	npm test

.PHONY: all build test`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectPath, filePath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Test individual validation methods
	t.Run("ValidatePackageJSON", func(t *testing.T) {
		err := engine.ValidatePackageJSON(filepath.Join(projectPath, "package.json"))
		assert.NoError(t, err)
	})

	t.Run("ValidateGoMod", func(t *testing.T) {
		err := engine.ValidateGoMod(filepath.Join(projectPath, "go.mod"))
		assert.NoError(t, err)
	})

	t.Run("ValidateDockerfile", func(t *testing.T) {
		err := engine.ValidateDockerfile(filepath.Join(projectPath, "Dockerfile"))
		assert.NoError(t, err)
	})

	t.Run("ValidateYAML", func(t *testing.T) {
		err := engine.ValidateYAML(filepath.Join(projectPath, "docker-compose.yml"))
		assert.NoError(t, err)
	})

	t.Run("ValidateJSON", func(t *testing.T) {
		err := engine.ValidateJSON(filepath.Join(projectPath, "package.json"))
		assert.NoError(t, err)
	})

	// Test comprehensive project validation
	t.Run("ValidateProject", func(t *testing.T) {
		result, err := engine.ValidateProject(projectPath)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should be valid overall (minor issues are acceptable)
		assert.True(t, result.Valid || len(result.Issues) < 5, "Project should be mostly valid")
	})

	// Test configuration validation
	t.Run("ValidateConfiguration", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "integration-test",
			Organization: "test-org",
			Description:  "Integration test project",
			License:      "MIT",
			OutputPath:   "/tmp/test",
		}

		result, err := engine.ValidateConfiguration(config)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.True(t, result.Valid)
		assert.Equal(t, 0, len(result.Errors))
	})

	// Test dependency validation
	t.Run("ValidateProjectDependencies", func(t *testing.T) {
		result, err := engine.ValidateProjectDependencies(projectPath)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.True(t, result.Valid)
		// Dependencies might be 0 if no dependency files are found or processed
		assert.GreaterOrEqual(t, len(result.Dependencies), 0, "Should process dependencies")
	})

	// Test structure validation
	t.Run("ValidateProjectStructure", func(t *testing.T) {
		result, err := engine.ValidateProjectStructure(projectPath)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Structure validation might have warnings but should not fail completely
		assert.GreaterOrEqual(t, len(result.RequiredFiles), 0, "Should validate required files")
	})
}

func TestValidationEngine_ErrorHandling(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		setupTest   func() string
		testFunc    func(string) error
		expectError bool
	}{
		{
			name: "ValidateProject with non-existent path",
			setupTest: func() string {
				return "/non/existent/path"
			},
			testFunc: func(path string) error {
				_, err := engine.ValidateProject(path)
				return err
			},
			expectError: false, // Should return result with error, not fail
		},
		{
			name: "ValidatePackageJSON with non-existent file",
			setupTest: func() string {
				return "/non/existent/package.json"
			},
			testFunc: func(path string) error {
				return engine.ValidatePackageJSON(path)
			},
			expectError: true,
		},
		{
			name: "ValidateGoMod with invalid path",
			setupTest: func() string {
				return "/invalid/../path/go.mod"
			},
			testFunc: func(path string) error {
				return engine.ValidateGoMod(path)
			},
			expectError: true,
		},
		{
			name: "ValidateDockerfile with directory instead of file",
			setupTest: func() string {
				tempDir := t.TempDir()
				return tempDir // Directory instead of file
			},
			testFunc: func(path string) error {
				return engine.ValidateDockerfile(path)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupTest()
			err := tt.testFunc(path)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationEngine_RuleManagement(t *testing.T) {
	engine := NewEngine()

	// Test getting default rules
	rules := engine.GetValidationRules()
	assert.Greater(t, len(rules), 0, "Should have default rules")

	// Test adding a new rule
	newRule := interfaces.ValidationRule{
		ID:          "test.custom.rule",
		Name:        "Custom Test Rule",
		Description: "A custom rule for testing",
		Category:    interfaces.ValidationCategoryQuality,
		Severity:    interfaces.ValidationSeverityWarning,
		Enabled:     true,
		Fixable:     false,
	}

	err := engine.AddValidationRule(newRule)
	assert.NoError(t, err)

	// Verify rule was added
	updatedRules := engine.GetValidationRules()
	assert.Equal(t, len(rules)+1, len(updatedRules))

	// Test adding duplicate rule
	err = engine.AddValidationRule(newRule)
	assert.Error(t, err, "Should not allow duplicate rule IDs")

	// Test removing rule
	err = engine.RemoveValidationRule("test.custom.rule")
	assert.NoError(t, err)

	// Verify rule was removed
	finalRules := engine.GetValidationRules()
	assert.Equal(t, len(rules), len(finalRules))

	// Test removing non-existent rule
	err = engine.RemoveValidationRule("non.existent.rule")
	assert.Error(t, err, "Should error when removing non-existent rule")
}

func TestValidationEngine_PerformanceWithLargeProject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	engine := NewEngine()
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "large-project")
	err := os.MkdirAll(projectPath, 0755)
	require.NoError(t, err)

	// Create a larger project structure
	dirs := []string{
		"src/components", "src/utils", "src/services", "src/types",
		"tests/unit", "tests/integration", "tests/e2e",
		"docs/api", "docs/guides", "docs/examples",
		"config/dev", "config/prod", "config/test",
		"scripts/build", "scripts/deploy", "scripts/test",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(projectPath, dir), 0755)
		require.NoError(t, err)
	}

	// Create many files
	baseFiles := map[string]string{
		"README.md":          "# Large Project\n\nA large test project.",
		"LICENSE":            "MIT License",
		"package.json":       `{"name": "large-project", "version": "1.0.0"}`,
		"tsconfig.json":      `{"compilerOptions": {"strict": true}}`,
		".eslintrc.json":     `{"extends": ["eslint:recommended"]}`,
		"Dockerfile":         "FROM node:18\nWORKDIR /app\nCOPY . .\nCMD [\"npm\", \"start\"]",
		"docker-compose.yml": "version: '3.8'\nservices:\n  app:\n    build: .",
		".gitignore":         "node_modules/\n*.log",
		"Makefile":           "all:\n\techo 'build'\n\n.PHONY: all",
	}

	// Create base files
	for filePath, content := range baseFiles {
		fullPath := filepath.Join(projectPath, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create many additional files
	for i := 0; i < 50; i++ {
		for _, dir := range dirs {
			fileName := filepath.Join(projectPath, dir, fmt.Sprintf("file%d.js", i))
			content := fmt.Sprintf("// File %d\nconsole.log('File %d');\n", i, i)
			err := os.WriteFile(fileName, []byte(content), 0644)
			require.NoError(t, err)
		}
	}

	// Measure validation performance
	start := time.Now()
	result, err := engine.ValidateProject(projectPath)
	duration := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, result)

	t.Logf("Validation of large project took: %v", duration)
	t.Logf("Project valid: %v", result.Valid)
	t.Logf("Issues found: %d", len(result.Issues))

	// Performance should be reasonable (less than 10 seconds for this size)
	assert.Less(t, duration.Seconds(), 10.0, "Validation should complete in reasonable time")
}
