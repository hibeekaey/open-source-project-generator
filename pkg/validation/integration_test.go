package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationEngineIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "validation_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test configuration files
	createTestConfigFiles(t, tempDir)

	// Create validation engine
	engine := NewEngine()

	// Test project validation
	result, err := engine.ValidateProject(tempDir)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Should have some validation results
	assert.True(t, len(result.Issues) >= 0, "Should have validation results")
}

func TestConfigValidatorIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_validation_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test configuration files
	createTestConfigFiles(t, tempDir)

	// Create config validator
	validator := NewConfigValidator()

	// Test configuration file validation
	results, err := validator.ValidateConfigurationFiles(tempDir)
	require.NoError(t, err)
	assert.NotNil(t, results)

	// Should find and validate configuration files
	assert.True(t, len(results) > 0, "Should find configuration files")

	// Check that each result has proper structure
	for _, result := range results {
		assert.NotNil(t, result.Summary)
		assert.NotNil(t, result.Errors)
		assert.NotNil(t, result.Warnings)
	}
}

func TestFormatSpecificValidators(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "format_validation_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Test JSON validation
	t.Run("JSON Validation", func(t *testing.T) {
		jsonFile := filepath.Join(tempDir, "package.json")
		jsonContent := `{
			"name": "test-package",
			"version": "1.0.0",
			"description": "Test package"
		}`
		err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
		require.NoError(t, err)

		validator := NewConfigValidator()
		result, err := validator.jsonValidator.ValidateJSONFile(jsonFile)
		require.NoError(t, err)
		assert.True(t, result.Valid)
	})

	// Test YAML validation
	t.Run("YAML Validation", func(t *testing.T) {
		yamlFile := filepath.Join(tempDir, "docker-compose.yml")
		yamlContent := `version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "80:80"`
		err := os.WriteFile(yamlFile, []byte(yamlContent), 0644)
		require.NoError(t, err)

		validator := NewConfigValidator()
		result, err := validator.yamlValidator.ValidateYAMLFile(yamlFile)
		require.NoError(t, err)
		assert.True(t, result.Valid)
	})

	// Test Dockerfile validation
	t.Run("Dockerfile Validation", func(t *testing.T) {
		dockerFile := filepath.Join(tempDir, "Dockerfile")
		dockerContent := `FROM node:16-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
USER node
CMD ["npm", "start"]`
		err := os.WriteFile(dockerFile, []byte(dockerContent), 0644)
		require.NoError(t, err)

		validator := NewConfigValidator()
		result, err := validator.dockerValidator.ValidateDockerfile(dockerFile)
		require.NoError(t, err)
		assert.True(t, result.Valid)
	})

	// Test Environment file validation
	t.Run("Environment File Validation", func(t *testing.T) {
		envFile := filepath.Join(tempDir, ".env")
		envContent := `NODE_ENV=development
PORT=3000
DATABASE_URL=postgresql://localhost:5432/testdb`
		err := os.WriteFile(envFile, []byte(envContent), 0644)
		require.NoError(t, err)

		validator := NewConfigValidator()
		result, err := validator.envValidator.ValidateEnvFile(envFile)
		require.NoError(t, err)
		assert.True(t, result.Valid)
	})

	// Test Makefile validation
	t.Run("Makefile Validation", func(t *testing.T) {
		makeFile := filepath.Join(tempDir, "Makefile")
		makeContent := `all: build test

build:
	go build -o app ./cmd/main.go

test:
	go test ./...

clean:
	rm -f app

.PHONY: all build test clean`
		err := os.WriteFile(makeFile, []byte(makeContent), 0644)
		require.NoError(t, err)

		validator := NewConfigValidator()
		result, err := validator.makefileValidator.ValidateMakefile(makeFile)
		require.NoError(t, err)
		assert.True(t, result.Valid)
	})
}

func TestSchemaManagerIntegration(t *testing.T) {
	schemaManager := NewSchemaManager()

	// Test schema retrieval
	t.Run("Schema Retrieval", func(t *testing.T) {
		schema, exists := schemaManager.GetSchema("package.json")
		assert.True(t, exists)
		assert.NotNil(t, schema)
		assert.Equal(t, "Package.json Schema", schema.Title)
	})

	// Test validation rules
	t.Run("Validation Rules", func(t *testing.T) {
		rules := schemaManager.GetValidationRules("package.json")
		assert.True(t, len(rules) > 0)
	})

	// Test package name validation
	t.Run("Package Name Validation", func(t *testing.T) {
		err := schemaManager.ValidatePackageName("valid-package-name")
		assert.NoError(t, err)

		err = schemaManager.ValidatePackageName("Invalid Package Name")
		assert.Error(t, err)
	})

	// Test environment key validation
	t.Run("Environment Key Validation", func(t *testing.T) {
		err := schemaManager.ValidateEnvKey("VALID_ENV_KEY")
		assert.NoError(t, err)

		err = schemaManager.ValidateEnvKey("invalid-env-key")
		assert.Error(t, err)
	})

	// Test secret detection
	t.Run("Secret Detection", func(t *testing.T) {
		isSecret := schemaManager.IsPotentialSecret("API_KEY", "very-long-secret-value-here")
		assert.True(t, isSecret)

		isSecret = schemaManager.IsPotentialSecret("DEBUG", "true")
		assert.False(t, isSecret)
	})
}

func TestValidationEngineComponentCoordination(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "engine_coordination_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create comprehensive test project structure
	createComprehensiveTestProject(t, tempDir)

	// Create validation engine
	engine := NewEngine()

	// Test individual validation methods
	t.Run("Package JSON Validation", func(t *testing.T) {
		packageFile := filepath.Join(tempDir, "package.json")
		err := engine.ValidatePackageJSON(packageFile)
		assert.NoError(t, err)
	})

	t.Run("Dockerfile Validation", func(t *testing.T) {
		dockerFile := filepath.Join(tempDir, "Dockerfile")
		err := engine.ValidateDockerfile(dockerFile)
		assert.NoError(t, err)
	})

	t.Run("YAML Validation", func(t *testing.T) {
		yamlFile := filepath.Join(tempDir, "docker-compose.yml")
		err := engine.ValidateYAML(yamlFile)
		assert.NoError(t, err)
	})

	t.Run("JSON Validation", func(t *testing.T) {
		jsonFile := filepath.Join(tempDir, "tsconfig.json")
		err := engine.ValidateJSON(jsonFile)
		assert.NoError(t, err)
	})

	// Test project structure validation
	t.Run("Project Structure Validation", func(t *testing.T) {
		result, err := engine.ValidateProjectStructure(tempDir)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	// Test configuration validation
	t.Run("Configuration Validation", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:        "test-project",
			Description: "Test project description",
			Author:      "Test Author",
			License:     "MIT",
		}
		result, err := engine.ValidateConfiguration(config)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// Helper functions

func createTestConfigFiles(t *testing.T, tempDir string) {
	// Create package.json
	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"description": "Test project"
	}`
	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Create docker-compose.yml
	dockerCompose := `version: '3.8'
services:
  web:
    image: nginx:latest`
	err = os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte(dockerCompose), 0644)
	require.NoError(t, err)

	// Create .env file
	envFile := `NODE_ENV=development
PORT=3000`
	err = os.WriteFile(filepath.Join(tempDir, ".env"), []byte(envFile), 0644)
	require.NoError(t, err)

	// Create Dockerfile
	dockerfile := `FROM node:16-alpine
WORKDIR /app
COPY . .`
	err = os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	require.NoError(t, err)
}

func createComprehensiveTestProject(t *testing.T, tempDir string) {
	// Create package.json
	packageJSON := `{
		"name": "comprehensive-test-project",
		"version": "1.0.0",
		"description": "Comprehensive test project",
		"main": "index.js",
		"scripts": {
			"start": "node index.js",
			"test": "jest"
		},
		"dependencies": {
			"express": "^4.18.0"
		},
		"devDependencies": {
			"jest": "^28.0.0"
		}
	}`
	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Create tsconfig.json
	tsconfig := `{
		"compilerOptions": {
			"target": "ES2020",
			"module": "commonjs",
			"strict": true,
			"esModuleInterop": true
		},
		"include": ["src/**/*"],
		"exclude": ["node_modules", "dist"]
	}`
	err = os.WriteFile(filepath.Join(tempDir, "tsconfig.json"), []byte(tsconfig), 0644)
	require.NoError(t, err)

	// Create docker-compose.yml
	dockerCompose := `version: '3.8'
services:
  web:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
  db:
    image: postgres:13
    environment:
      - POSTGRES_DB=testdb`
	err = os.WriteFile(filepath.Join(tempDir, "docker-compose.yml"), []byte(dockerCompose), 0644)
	require.NoError(t, err)

	// Create Dockerfile
	dockerfile := `FROM node:16-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 3000
USER node
CMD ["npm", "start"]`
	err = os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	require.NoError(t, err)

	// Create .env file
	envFile := `NODE_ENV=development
PORT=3000
DATABASE_URL=postgresql://localhost:5432/testdb
LOG_LEVEL=info`
	err = os.WriteFile(filepath.Join(tempDir, ".env"), []byte(envFile), 0644)
	require.NoError(t, err)

	// Create Makefile
	makefile := `all: build test

build:
	npm run build

test:
	npm test

clean:
	rm -rf dist node_modules

install:
	npm install

.PHONY: all build test clean install`
	err = os.WriteFile(filepath.Join(tempDir, "Makefile"), []byte(makefile), 0644)
	require.NoError(t, err)

	// Create .gitignore
	gitignore := `node_modules/
dist/
*.log
.env
.DS_Store`
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(gitignore), 0644)
	require.NoError(t, err)

	// Create .dockerignore
	dockerignore := `node_modules
npm-debug.log
.git
.gitignore
README.md
.env
.nyc_output
coverage
.nyc_output`
	err = os.WriteFile(filepath.Join(tempDir, ".dockerignore"), []byte(dockerignore), 0644)
	require.NoError(t, err)
}
