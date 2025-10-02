package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProcessingEngineIntegration tests the processing engine with various template types
func TestProcessingEngineIntegration(t *testing.T) {
	engine := NewProcessingEngine()

	// Create a temporary directory for templates
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a realistic project template structure
	templates := map[string]string{
		"README.md.tmpl": `# {{.Name | title}}

{{.Description}}

## Components

{{if hasFrontend .}}### Frontend
- Next.js App: {{.Components.Frontend.NextJS.App}}
- Admin Panel: {{.Components.Frontend.NextJS.Admin}}
{{end}}

{{if hasBackend .}}### Backend
- Go/Gin API: {{.Components.Backend.GoGin}}
{{end}}

## Installation

` + "```" + `bash
npm install {{.Name | kebabCase}}
` + "```" + `

## Configuration

- Node Version: {{nodeVersion .}}
- Go Version: {{goVersion .}}

## License

{{.License}}
`,
		"package.json.tmpl": `{
  "name": "{{.Name | kebabCase}}",
  "version": "1.0.0",
  "description": "{{.Description}}",
  "main": "index.js",
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "{{nextjsVersion .}}",
    "react": "{{reactVersion .}}"
  },
  "engines": {
    "node": "{{nodeVersion .}}"
  }
}`,
		"src/config.go.tmpl": `package main

import "fmt"

// Config represents the application configuration
type Config struct {
	Name         string
	Organization string
	Version      string
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return &Config{
		Name:         "{{.Name}}",
		Organization: "{{.Organization}}",
		Version:      "1.0.0",
	}
}

func main() {
	config := NewConfig()
	fmt.Printf("Starting %s by %s\n", config.Name, config.Organization)
}`,
		"docker-compose.yml.tmpl": `version: '3.8'

services:
  {{.Name | kebabCase}}-app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - NODE_VERSION={{nodeVersion .}}
    {{if hasBackend .}}
  {{.Name | kebabCase}}-api:
    build: ./api
    ports:
      - "8080:8080"
    environment:
      - GO_VERSION={{goVersion .}}
    {{end}}`,
		"static-file.txt": "This is a static file that should be copied as-is",
	}

	// Create template files
	for path, content := range templates {
		fullPath := filepath.Join(templateDir, path)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)

		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "AwesomeProject",
		Organization: "TechCorp",
		Description:  "An awesome project generated from templates",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Admin: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "18.17.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"next":  "13.4.0",
				"react": "18.2.0",
			},
		},
	}

	// Process the template directory
	err := engine.ProcessDirectory(templateDir, outputDir, config)
	require.NoError(t, err)

	// Verify README.md was processed correctly
	readmeContent, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	require.NoError(t, err)
	readmeStr := string(readmeContent)

	assert.Contains(t, readmeStr, "# Awesomeproject")
	assert.Contains(t, readmeStr, "An awesome project generated from templates")
	assert.Contains(t, readmeStr, "### Frontend")
	assert.Contains(t, readmeStr, "Next.js App: true")
	assert.Contains(t, readmeStr, "Admin Panel: true")
	assert.Contains(t, readmeStr, "### Backend")
	assert.Contains(t, readmeStr, "Go/Gin API: true")
	assert.Contains(t, readmeStr, "npm install awesome-project")
	assert.Contains(t, readmeStr, "Node Version: 18.17.0")
	assert.Contains(t, readmeStr, "Go Version: 1.21.0")
	assert.Contains(t, readmeStr, "MIT")

	// Verify package.json was processed correctly
	packageContent, err := os.ReadFile(filepath.Join(outputDir, "package.json"))
	require.NoError(t, err)
	packageStr := string(packageContent)

	assert.Contains(t, packageStr, `"name": "awesome-project"`)
	assert.Contains(t, packageStr, `"description": "An awesome project generated from templates"`)
	assert.Contains(t, packageStr, `"next": "13.4.0"`)
	assert.Contains(t, packageStr, `"react": "18.2.0"`)
	assert.Contains(t, packageStr, `"node": "18.17.0"`)

	// Verify Go config was processed correctly
	configContent, err := os.ReadFile(filepath.Join(outputDir, "src", "config.go"))
	require.NoError(t, err)
	configStr := string(configContent)

	assert.Contains(t, configStr, `Name:         "AwesomeProject"`)
	assert.Contains(t, configStr, `Organization: "TechCorp"`)
	assert.Contains(t, configStr, `fmt.Printf("Starting %s by %s\n", config.Name, config.Organization)`)

	// Verify docker-compose.yml was processed correctly
	dockerContent, err := os.ReadFile(filepath.Join(outputDir, "docker-compose.yml"))
	require.NoError(t, err)
	dockerStr := string(dockerContent)

	assert.Contains(t, dockerStr, "awesome-project-app:")
	assert.Contains(t, dockerStr, "NODE_VERSION=18.17.0")
	assert.Contains(t, dockerStr, "awesome-project-api:")
	assert.Contains(t, dockerStr, "GO_VERSION=1.21.0")

	// Verify static file was copied as-is
	staticContent, err := os.ReadFile(filepath.Join(outputDir, "static-file.txt"))
	require.NoError(t, err)
	assert.Equal(t, "This is a static file that should be copied as-is", string(staticContent))
}

// TestProcessingEngineWithCustomFunctions tests custom function registration
func TestProcessingEngineWithCustomFunctions(t *testing.T) {
	engine := NewProcessingEngine()

	// Register custom functions
	customFuncs := map[string]interface{}{
		"reverse": func(s string) string {
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		},
		"multiply": func(a, b int) int {
			return a * b
		},
	}

	engine.RegisterFunctions(customFuncs)

	// Create a template using custom functions
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	templateContent := `Project: {{.Name}}
Reversed: {{reverse .Name}}
Calculation: {{multiply 5 3}}
Combined: {{.Name | reverse | upper}}`

	templatePath := filepath.Join(templateDir, "test.tmpl")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(t, err)

	config := &models.ProjectConfig{
		Name: "TestProject",
	}

	// Process the template
	err = engine.ProcessDirectory(templateDir, outputDir, config)
	require.NoError(t, err)

	// Verify the output
	outputContent, err := os.ReadFile(filepath.Join(outputDir, "test"))
	require.NoError(t, err)
	outputStr := string(outputContent)

	assert.Contains(t, outputStr, "Project: TestProject")
	assert.Contains(t, outputStr, "Reversed: tcejorPtseT")
	assert.Contains(t, outputStr, "Calculation: 15")
	assert.Contains(t, outputStr, "Combined: TCEJORPTSET")
}

// TestProcessingEngineErrorHandling tests error handling scenarios
func TestProcessingEngineErrorHandling(t *testing.T) {
	engine := NewProcessingEngine()

	t.Run("Invalid template syntax", func(t *testing.T) {
		templateDir := t.TempDir()
		outputDir := t.TempDir()

		// Create template with invalid syntax
		templateContent := "{{.Name" // Missing closing brace
		templatePath := filepath.Join(templateDir, "invalid.tmpl")
		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		require.NoError(t, err)

		config := &models.ProjectConfig{Name: "Test"}

		// Should fail during processing
		err = engine.ProcessDirectory(templateDir, outputDir, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to process template")
	})

	t.Run("Template execution error", func(t *testing.T) {
		templateDir := t.TempDir()
		outputDir := t.TempDir()

		// Create template that will fail during execution
		templateContent := "{{.Name | divide 0}}" // Division by zero
		templatePath := filepath.Join(templateDir, "error.tmpl")
		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		require.NoError(t, err)

		config := &models.ProjectConfig{Name: "Test"}

		// Should fail during processing
		err = engine.ProcessDirectory(templateDir, outputDir, config)
		assert.Error(t, err)
	})
}

// BenchmarkProcessingEngineIntegration benchmarks the full processing pipeline
func BenchmarkProcessingEngineIntegration(b *testing.B) {
	engine := NewProcessingEngine()

	// Create a realistic template structure
	templateDir := b.TempDir()

	templates := map[string]string{
		"README.md.tmpl":      "# {{.Name | title}}\n{{.Description}}",
		"package.json.tmpl":   `{"name": "{{.Name | kebabCase}}", "version": "1.0.0"}`,
		"src/main.go.tmpl":    "package main\n// Project: {{.Name}}",
		"config/app.yml.tmpl": "name: {{.Name}}\nversion: 1.0.0",
	}

	for path, content := range templates {
		fullPath := filepath.Join(templateDir, path)
		dir := filepath.Dir(fullPath)
		_ = os.MkdirAll(dir, 0755)
		_ = os.WriteFile(fullPath, []byte(content), 0644)
	}

	config := &models.ProjectConfig{
		Name:        "BenchmarkProject",
		Description: "A project for benchmarking",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputDir := b.TempDir()
		err := engine.ProcessDirectory(templateDir, outputDir, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}
