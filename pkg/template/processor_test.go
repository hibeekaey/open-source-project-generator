package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestNewDirectoryProcessor(t *testing.T) {
	t.Parallel()

	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)

	if processor == nil {
		t.Fatal("NewDirectoryProcessor returned nil")
	}

	if processor.engine != engine {
		t.Error("DirectoryProcessor engine not set correctly")
	}

	if processor.metadata == nil {
		t.Error("DirectoryProcessor metadata parser not initialized")
	}
}

func TestProcessTemplateDirectory(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory structure
	err := os.MkdirAll(filepath.Join(templateDir, "src", "components"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template files
	templateFiles := map[string]string{
		"template.yaml": `name: test-template
version: 1.0.0
conditions:
  - name: has_frontend
    component: frontend
    value: true
    operator: eq
files:
  - source: frontend/app.js.tmpl
    conditions:
      - name: frontend_enabled
        component: frontend
        value: true
        operator: eq
  - source: backend/server.go.tmpl
    conditions:
      - name: backend_enabled
        component: backend
        value: true
        operator: eq`,
		"README.md.tmpl": `# {{.Name}}

{{.Description}}

## Components
{{if hasFrontend .}}
- Frontend Application
{{end}}
{{if hasBackend .}}
- Backend API
{{end}}`,
		"frontend/app.js.tmpl": `// {{.Name}} Frontend Application
console.log('Hello from {{.Name}}!');

{{if .Components.Frontend.NextJS.App}}
// Main application code
{{end}}`,
		"backend/server.go.tmpl": `package main

// {{.Name}} Backend Server
func main() {
    println("Starting {{.Name}} server...")
}`,
		"assets/logo.png": "fake-binary-content",
		"src/components/Button.tsx.tmpl": `import React from 'react';

interface ButtonProps {
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({ children }) => {
  return <button className="btn">{{"{{"}}children{{"}}"}}</button>;
};`,
	}

	for filePath, content := range templateFiles {
		fullPath := filepath.Join(templateDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", filePath, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", filePath, err)
		}
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "MyTestApp",
		Organization: "test-org",
		Description:  "A test application for template processing",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
	}

	// Process directory
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify output files
	expectedFiles := map[string][]string{
		"README.md": {
			"# MyTestApp",
			"A test application for template processing",
			"- Frontend Application",
			"- Backend API",
		},
		"frontend/app.js": {
			"// MyTestApp Frontend Application",
			"console.log('Hello from MyTestApp!');",
			"// Main application code",
		},
		"backend/server.go": {
			"package main",
			"// MyTestApp Backend Server",
			`println("Starting MyTestApp server...")`,
		},
		"assets/logo.png": {
			"fake-binary-content",
		},
		"src/components/Button.tsx": {
			"import React from 'react';",
			"interface ButtonProps",
			"export const Button",
		},
	}

	for filePath, expectedContent := range expectedFiles {
		fullPath := filepath.Join(outputDir, filePath)

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", filePath)
			continue
		}

		// Read and verify content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", filePath, err)
			continue
		}

		contentStr := string(content)
		for _, expected := range expectedContent {
			if !strings.Contains(contentStr, expected) {
				t.Errorf("File %s does not contain expected content '%s', got: %s",
					filePath, expected, contentStr)
			}
		}
	}

	// Verify template.yaml is not copied
	metadataPath := filepath.Join(outputDir, "template.yaml")
	if _, err := os.Stat(metadataPath); !os.IsNotExist(err) {
		t.Error("Metadata file should not be copied to output")
	}
}

func TestProcessTemplateDirectoryWithConditionals(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory structure
	err := os.MkdirAll(filepath.Join(templateDir, "optional"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template files with conditions
	templateFiles := map[string]string{
		"template.yaml": `name: conditional-template
version: 1.0.0
files:
  - source: frontend-only.txt.tmpl
    conditions:
      - name: has_frontend
        component: frontend
        value: true
        operator: eq
  - source: backend-only.txt.tmpl
    conditions:
      - name: has_backend
        component: backend
        value: true
        operator: eq`,
		"always-included.txt.tmpl": "This file is always included: {{.Name}}",
		"frontend-only.txt.tmpl":   "Frontend file: {{.Name}}",
		"backend-only.txt.tmpl":    "Backend file: {{.Name}}",
	}

	for filePath, content := range templateFiles {
		fullPath := filepath.Join(templateDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", filePath, err)
		}
	}

	// Test with frontend only
	configFrontendOnly := &models.ProjectConfig{
		Name: "FrontendApp",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: false,
			},
		},
	}

	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err = processor.ProcessTemplateDirectory(templateDir, outputDir, configFrontendOnly)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify frontend-only file exists
	frontendFile := filepath.Join(outputDir, "frontend-only.txt")
	if _, err := os.Stat(frontendFile); os.IsNotExist(err) {
		t.Error("Frontend-only file should exist when frontend is enabled")
	}

	// Verify backend-only file does not exist
	backendFile := filepath.Join(outputDir, "backend-only.txt")
	if _, err := os.Stat(backendFile); !os.IsNotExist(err) {
		t.Error("Backend-only file should not exist when backend is disabled")
	}

	// Verify always-included file exists
	alwaysFile := filepath.Join(outputDir, "always-included.txt")
	if _, err := os.Stat(alwaysFile); os.IsNotExist(err) {
		t.Error("Always-included file should exist")
	}
}

func TestProcessTemplateWithInheritance(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create base template
	baseTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Name}}</title>
</head>
<body>
    <h1>{{.Name}}</h1>
    {{/* content */}}
    <footer>© {{.Organization}}</footer>
</body>
</html>`

	// Create child template that extends base
	childTemplate := `{{/* extends "base.html.tmpl" */}}
<main>
    <p>Welcome to {{.Name}}!</p>
    <p>{{.Description}}</p>
</main>`

	// Write templates
	err = os.WriteFile(filepath.Join(templateDir, "base.html.tmpl"), []byte(baseTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create base template: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "index.html.tmpl"), []byte(childTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create child template: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "TestApp",
		Organization: "TestOrg",
		Description:  "A test application",
	}

	// Process directory
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify output
	outputFile := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	expectedContent := []string{
		"<!DOCTYPE html>",
		"<title>TestApp</title>",
		"<h1>TestApp</h1>",
		"<main>",
		"<p>Welcome to TestApp!</p>",
		"<p>A test application</p>",
		"</main>",
		"<footer>© TestOrg</footer>",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Output does not contain expected content '%s', got: %s", expected, contentStr)
		}
	}

	// Verify base template is not copied
	baseFile := filepath.Join(outputDir, "base.html")
	if _, err := os.Stat(baseFile); !os.IsNotExist(err) {
		t.Error("Base template should not be copied to output")
	}
}

func TestProcessTemplateWithIncludes(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory with partials subdirectory
	err := os.MkdirAll(filepath.Join(templateDir, "partials"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create partial templates
	headerPartial := `<header>
    <h1>{{.Name}}</h1>
    <nav>Navigation</nav>
</header>`

	footerPartial := `<footer>
    <p>© {{.Organization}} - {{.Name}}</p>
</footer>`

	// Create main template with includes
	mainTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Name}}</title>
</head>
<body>
    {{/* include "partials/header.tmpl" */}}
    
    <main>
        <p>{{.Description}}</p>
    </main>
    
    {{/* include "partials/footer.tmpl" */}}
</body>
</html>`

	// Write templates
	err = os.WriteFile(filepath.Join(templateDir, "partials", "header.tmpl"), []byte(headerPartial), 0644)
	if err != nil {
		t.Fatalf("Failed to create header partial: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "partials", "footer.tmpl"), []byte(footerPartial), 0644)
	if err != nil {
		t.Fatalf("Failed to create footer partial: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "page.html.tmpl"), []byte(mainTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create main template: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "IncludeTest",
		Organization: "TestOrg",
		Description:  "Testing template includes",
	}

	// Process directory
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify output
	outputFile := filepath.Join(outputDir, "page.html")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	expectedContent := []string{
		"<!DOCTYPE html>",
		"<title>IncludeTest</title>",
		"<header>",
		"<h1>IncludeTest</h1>",
		"<nav>Navigation</nav>",
		"</header>",
		"<main>",
		"<p>Testing template includes</p>",
		"</main>",
		"<footer>",
		"<p>© TestOrg - IncludeTest</p>",
		"</footer>",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Output does not contain expected content '%s', got: %s", expected, contentStr)
		}
	}

	// Verify partial templates are not copied
	headerFile := filepath.Join(outputDir, "partials", "header.tmpl")
	if _, err := os.Stat(headerFile); !os.IsNotExist(err) {
		t.Error("Partial templates should not be copied to output")
	}
}

func TestCopyAsset(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")

	err := os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create test asset files
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	textContent := []byte("This is a text file")

	srcBinary := filepath.Join(srcDir, "image.png")
	srcText := filepath.Join(srcDir, "readme.txt")

	err = os.WriteFile(srcBinary, binaryContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	err = os.WriteFile(srcText, textContent, 0755) // Different permissions
	if err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	// Test asset copying
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)

	destBinary := filepath.Join(destDir, "image.png")
	destText := filepath.Join(destDir, "readme.txt")

	// Copy binary file
	err = processor.copyAsset(srcBinary, destBinary)
	if err != nil {
		t.Fatalf("Failed to copy binary asset: %v", err)
	}

	// Copy text file
	err = processor.copyAsset(srcText, destText)
	if err != nil {
		t.Fatalf("Failed to copy text asset: %v", err)
	}

	// Verify binary file content
	copiedBinary, err := os.ReadFile(destBinary)
	if err != nil {
		t.Fatalf("Failed to read copied binary file: %v", err)
	}

	if string(copiedBinary) != string(binaryContent) {
		t.Error("Binary file content does not match")
	}

	// Verify text file content
	copiedText, err := os.ReadFile(destText)
	if err != nil {
		t.Fatalf("Failed to read copied text file: %v", err)
	}

	if string(copiedText) != string(textContent) {
		t.Error("Text file content does not match")
	}

	// Verify file permissions are preserved
	srcInfo, err := os.Stat(srcText)
	if err != nil {
		t.Fatalf("Failed to get source file info: %v", err)
	}

	destInfo, err := os.Stat(destText)
	if err != nil {
		t.Fatalf("Failed to get destination file info: %v", err)
	}

	if srcInfo.Mode() != destInfo.Mode() {
		t.Errorf("File permissions not preserved: expected %v, got %v", srcInfo.Mode(), destInfo.Mode())
	}
}
func TestProcessPathTemplate(t *testing.T) {
	t.Parallel()

	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)

	config := &models.ProjectConfig{
		Name:         "MyTestApp",
		Organization: "test-org",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name replacement",
			input:    "src/{{.Name}}/main.go",
			expected: "src/MyTestApp/main.go",
		},
		{
			name:     "kebab case name",
			input:    "{{kebabCase .Name}}/config.yaml",
			expected: "my-test-app/config.yaml",
		},
		{
			name:     "snake case name",
			input:    "{{snakeCase .Name}}/database.sql",
			expected: "my_test_app/database.sql",
		},
		{
			name:     "organization replacement",
			input:    "{{.Organization}}/{{.Name}}/app.js",
			expected: "test-org/MyTestApp/app.js",
		},
		{
			name:     "no template variables",
			input:    "static/assets/logo.png",
			expected: "static/assets/logo.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := processor.processPathTemplate(tt.input, config)
			if err != nil {
				t.Fatalf("processPathTemplate failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestShouldProcessPath(t *testing.T) {
	t.Parallel()

	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)

	metadata := &TemplateMetadata{
		Files: []TemplateFile{
			{
				Source: "frontend/app.js.tmpl",
				Conditions: []TemplateCondition{
					{
						Name:      "has_frontend",
						Component: "frontend",
						Value:     true,
						Operator:  "eq",
					},
				},
			},
			{
				Source: "backend/server.go.tmpl",
				Conditions: []TemplateCondition{
					{
						Name:      "has_backend",
						Component: "backend",
						Value:     true,
						Operator:  "eq",
					},
				},
			},
		},
		Conditions: []TemplateCondition{
			{
				Name:      "default_condition",
				Component: "always",
				Value:     true,
				Operator:  "eq",
			},
		},
	}

	configWithFrontend := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: false,
			},
		},
	}

	configWithoutFrontend := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: false,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
	}

	tests := []struct {
		name     string
		path     string
		config   *models.ProjectConfig
		expected bool
	}{
		{
			name:     "frontend file with frontend enabled",
			path:     "frontend/app.js.tmpl",
			config:   configWithFrontend,
			expected: true,
		},
		{
			name:     "frontend file with frontend disabled",
			path:     "frontend/app.js.tmpl",
			config:   configWithoutFrontend,
			expected: false,
		},
		{
			name:     "backend file with backend enabled",
			path:     "backend/server.go.tmpl",
			config:   configWithoutFrontend,
			expected: true,
		},
		{
			name:     "backend file with backend disabled",
			path:     "backend/server.go.tmpl",
			config:   configWithFrontend,
			expected: false,
		},
		{
			name:     "unspecified file uses global conditions",
			path:     "README.md.tmpl",
			config:   configWithFrontend,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := processor.shouldProcessPath(tt.path, metadata, tt.config)
			if err != nil {
				t.Fatalf("shouldProcessPath failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestComplexTemplateDirectoryProcessing(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Simplified template structure for better performance
	dirs := []string{
		"frontend/src",
		"backend/cmd",
		"docs",
		"partials",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(templateDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Simplified template metadata
	metadata := `name: simple-template
version: 1.0.0
description: Simplified template for testing
files:
  - source: frontend/src/App.tsx.tmpl
    conditions:
      - name: has_frontend
        component: frontend
        value: true
        operator: eq
  - source: backend/cmd/main.go.tmpl
    conditions:
      - name: has_backend
        component: backend
        value: true
        operator: eq`

	// Create base template for inheritance
	baseTemplate := `{{/* Base template for {{.Name}} */}}
# {{.Name}}

{{.Description}}

## Components
{{/* content */}}

---
Generated by {{.Organization}}`

	// Simplified templates for better performance
	templates := map[string]string{
		"template.yaml": metadata,
		"base.md.tmpl":  baseTemplate,
		"README.md.tmpl": `{{/* extends "base.md.tmpl" */}}
{{if hasFrontend .}}
- Frontend: React Application
{{end}}
{{if hasBackend .}}
- Backend: Go API Server
{{end}}`,
		"frontend/src/App.tsx.tmpl": `import React from 'react';

function App() {
  return (
    <div className="App">
      <h1>{{.Name}}</h1>
      <p>{{.Description}}</p>
    </div>
  );
}

export default App;`,
		"backend/cmd/main.go.tmpl": `package main

import "fmt"

func main() {
	fmt.Println("Starting {{.Name}} server...")
}`,
		"partials/header.tmpl": `<header>
  <h1>{{.Name}}</h1>
  <p>by {{.Organization}}</p>
</header>`,
		"docs/api.md.tmpl": `# {{.Name}} API Documentation

{{/* include "partials/header.tmpl" */}}

## Overview
{{.Description}}`,
	}

	// Write all template files
	for filePath, content := range templates {
		fullPath := filepath.Join(templateDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", filePath, err)
		}
	}

	// Simplified test configuration
	config := &models.ProjectConfig{
		Name:         "SimpleApp",
		Organization: "TechCorp",
		Description:  "A simple test application",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
	}

	// Process directory
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err := processor.ProcessTemplateDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify expected files exist and have correct content
	expectedFiles := map[string][]string{
		"README.md": {
			"# SimpleApp",
			"A simple test application",
			"- Frontend: React Application",
			"- Backend: Go API Server",
			"Generated by TechCorp",
		},
		"frontend/src/App.tsx": {
			"import React from 'react';",
			"<h1>SimpleApp</h1>",
			"<p>A simple test application</p>",
		},
		"backend/cmd/main.go": {
			"package main",
			"Starting SimpleApp server...",
		},
		"docs/api.md": {
			"# SimpleApp API Documentation",
			"<h1>SimpleApp</h1>",
			"<p>by TechCorp</p>",
			"## Overview",
		},
	}

	for filePath, expectedContent := range expectedFiles {
		fullPath := filepath.Join(outputDir, filePath)

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", filePath)
			continue
		}

		// Read and verify content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", filePath, err)
			continue
		}

		contentStr := string(content)
		for _, expected := range expectedContent {
			if !strings.Contains(contentStr, expected) {
				t.Errorf("File %s does not contain expected content '%s'", filePath, expected)
			}
		}
	}

	// Verify template files and partials are not copied
	templateFiles := []string{
		"template.yaml",
		"base.md.tmpl",
		"partials/header.tmpl",
	}

	for _, templateFile := range templateFiles {
		fullPath := filepath.Join(outputDir, templateFile)
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			t.Errorf("Template file %s should not be copied to output", templateFile)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	t.Parallel()

	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)

	t.Run("invalid template directory", func(t *testing.T) {
		t.Parallel()
		err := processor.ProcessTemplateDirectory("/nonexistent", "/tmp/output", &models.ProjectConfig{})
		if err == nil {
			t.Error("Expected error for nonexistent template directory")
		}
	})

	t.Run("invalid extends directive", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "template")
		outputDir := filepath.Join(tempDir, "output")

		err := os.MkdirAll(templateDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		// Create template with invalid extends
		invalidTemplate := `{{/* extends "nonexistent.tmpl" */}}
Content here`

		err = os.WriteFile(filepath.Join(templateDir, "invalid.tmpl"), []byte(invalidTemplate), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid template: %v", err)
		}

		config := &models.ProjectConfig{Name: "Test"}
		err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
		if err == nil {
			t.Error("Expected error for invalid extends directive")
		}
	})

	t.Run("invalid include directive", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "template")
		outputDir := filepath.Join(tempDir, "output")

		err := os.MkdirAll(templateDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		// Create template with invalid include
		invalidTemplate := `{{/* include "nonexistent.tmpl" */}}
Content here`

		err = os.WriteFile(filepath.Join(templateDir, "invalid.tmpl"), []byte(invalidTemplate), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid template: %v", err)
		}

		config := &models.ProjectConfig{Name: "Test"}
		err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
		if err == nil {
			t.Error("Expected error for invalid include directive")
		}
	})

	t.Run("malformed template syntax", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		templateDir := filepath.Join(tempDir, "template")
		outputDir := filepath.Join(tempDir, "output")

		err := os.MkdirAll(templateDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		// Create template with malformed syntax
		malformedTemplate := `{{.Name}
{{if .Components.Frontend.NextJS.App}
Content without closing tags`

		err = os.WriteFile(filepath.Join(templateDir, "malformed.tmpl"), []byte(malformedTemplate), 0644)
		if err != nil {
			t.Fatalf("Failed to create malformed template: %v", err)
		}

		config := &models.ProjectConfig{Name: "Test"}
		err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
		if err == nil {
			t.Error("Expected error for malformed template syntax")
		}
	})
}

func TestNestedTemplateInheritance(t *testing.T) {
	t.Parallel()

	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	err := os.MkdirAll(filepath.Join(templateDir, "layouts"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create base layout
	baseLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Name}} - {{.Description}}</title>
</head>
<body>
    <header><h1>{{.Name}}</h1></header>
    <main>
        {{/* content */}}
    </main>
    <footer>© {{.Organization}}</footer>
</body>
</html>`

	// Create page layout that extends base
	pageLayout := `{{/* extends "layouts/base.tmpl" */}}
<section>
    <h2>Page Layout</h2>
    {{/* content */}}
</section>`

	// Create specific page that extends page layout
	homePage := `{{/* extends "layouts/page.tmpl" */}}
<div>
    <h3>Welcome to {{.Name}}</h3>
    <p>{{.Description}}</p>
</div>`

	// Write templates
	err = os.WriteFile(filepath.Join(templateDir, "layouts", "base.tmpl"), []byte(baseLayout), 0644)
	if err != nil {
		t.Fatalf("Failed to create base layout: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "layouts", "page.tmpl"), []byte(pageLayout), 0644)
	if err != nil {
		t.Fatalf("Failed to create page layout: %v", err)
	}

	err = os.WriteFile(filepath.Join(templateDir, "home.html.tmpl"), []byte(homePage), 0644)
	if err != nil {
		t.Fatalf("Failed to create home page: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "NestedApp",
		Organization: "TestCorp",
		Description:  "Testing nested inheritance",
	}

	// Process directory
	engine := NewEngine().(*Engine)
	processor := NewDirectoryProcessor(engine)
	err = processor.ProcessTemplateDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessTemplateDirectory failed: %v", err)
	}

	// Verify output
	outputFile := filepath.Join(outputDir, "home.html")
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	expectedContent := []string{
		"<!DOCTYPE html>",
		"<title>NestedApp - Testing nested inheritance</title>",
		"<header><h1>NestedApp</h1></header>",
		"<section>",
		"<h2>Page Layout</h2>",
		"<h3>Welcome to NestedApp</h3>",
		"<p>Testing nested inheritance</p>",
		"</section>",
		"<footer>© TestCorp</footer>",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Output does not contain expected content '%s', got: %s", expected, contentStr)
		}
	}
}
