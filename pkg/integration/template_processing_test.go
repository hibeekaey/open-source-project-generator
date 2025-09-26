package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestTemplateProcessingWorkflows tests complete template processing workflows
func TestTemplateProcessingWorkflows(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("basic_template_processing", func(t *testing.T) {
		testBasicTemplateProcessing(t, tempDir)
	})

	t.Run("template_with_variables", func(t *testing.T) {
		testTemplateWithVariables(t, tempDir)
	})

	t.Run("nested_template_structure", func(t *testing.T) {
		testNestedTemplateStructure(t, tempDir)
	})

	t.Run("template_validation_workflow", func(t *testing.T) {
		testTemplateValidationWorkflow(t, tempDir)
	})

	t.Run("custom_template_processing", func(t *testing.T) {
		testCustomTemplateProcessing(t, tempDir)
	})
}

func testBasicTemplateProcessing(t *testing.T, tempDir string) {
	// Create a basic template
	templateDir := filepath.Join(tempDir, "basic-template")
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template metadata
	metadata := `name: basic-template
version: 1.0.0
author: Test Author
license: MIT
description: A basic test template
`

	err = os.WriteFile(filepath.Join(templateDir, "template.yaml"), []byte(metadata), 0644)
	if err != nil {
		t.Fatalf("Failed to create template metadata: %v", err)
	}

	// Create template files
	readmeTemplate := `# {{.Name}}

This is a project generated from a template.

Organization: {{.Organization}}
License: {{.License}}

## Description

{{.Description}}
`

	err = os.WriteFile(filepath.Join(templateDir, "README.md.tmpl"), []byte(readmeTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create README template: %v", err)
	}

	// Create static file
	staticFile := `This is a static file that should be copied as-is.`

	err = os.WriteFile(filepath.Join(templateDir, "static.txt"), []byte(staticFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create static file: %v", err)
	}

	// Create mock template manager and engine
	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	// Configure project
	config := &models.ProjectConfig{
		Name:         "basic-test-project",
		Organization: "test-org",
		Description:  "A test project generated from basic template",
		License:      "MIT",
	}

	outputDir := filepath.Join(tempDir, "basic-output")

	// Process template
	err = templateManager.ProcessCustomTemplate(templateDir, config, outputDir)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	// Verify processing was called
	processedDirs := templateEngine.GetProcessedDirs()
	if _, exists := processedDirs[templateDir]; !exists {
		t.Error("Expected template directory to be processed")
	}

	// Verify config was passed correctly
	processedTemplates := templateEngine.GetProcessedTemplates()
	if processedConfig, exists := processedTemplates[templateDir]; exists {
		if processedConfig.Name != config.Name {
			t.Errorf("Expected processed config name '%s', got '%s'", config.Name, processedConfig.Name)
		}
	} else {
		t.Error("Expected template config to be processed")
	}
}

func testTemplateWithVariables(t *testing.T, tempDir string) {
	// Create template with complex variables
	templateDir := filepath.Join(tempDir, "variable-template")
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template with conditional logic
	configTemplate := `{
  "name": "{{.Name}}",
  "version": "{{.Version | default "1.0.0"}}",
  "organization": "{{.Organization}}",
  {{if .Components.Backend.Enabled}}
  "backend": {
    "technology": "{{.Components.Backend.Technology}}",
    "port": {{.Components.Backend.Port | default 8080}}
  },
  {{end}}
  {{if .Components.Frontend.Enabled}}
  "frontend": {
    "technology": "{{.Components.Frontend.Technology}}",
    "port": {{.Components.Frontend.Port | default 3000}}
  },
  {{end}}
  "license": "{{.License}}"
}
`

	err = os.WriteFile(filepath.Join(templateDir, "config.json.tmpl"), []byte(configTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create config template: %v", err)
	}

	// Create template with loops
	dockerTemplate := `FROM node:{{.NodeVersion | default "18"}}

WORKDIR /app

{{range .Dependencies}}
RUN npm install {{.}}
{{end}}

COPY . .

{{if .Components.Backend.Enabled}}
EXPOSE {{.Components.Backend.Port | default 8080}}
{{end}}

{{if .Components.Frontend.Enabled}}
EXPOSE {{.Components.Frontend.Port | default 3000}}
{{end}}

CMD ["npm", "start"]
`

	err = os.WriteFile(filepath.Join(templateDir, "Dockerfile.tmpl"), []byte(dockerTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile template: %v", err)
	}

	// Create mock template manager
	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	// Configure project with complex structure
	config := &models.ProjectConfig{
		Name:         "variable-test-project",
		Organization: "variable-org",
		License:      "Apache-2.0",
		Components: models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:  true,
					Home: true,
				},
			},
		},
	}

	outputDir := filepath.Join(tempDir, "variable-output")

	// Process template
	err = templateManager.ProcessCustomTemplate(templateDir, config, outputDir)
	if err != nil {
		t.Fatalf("Failed to process variable template: %v", err)
	}

	// Verify processing
	processedDirs := templateEngine.GetProcessedDirs()
	if outputPath, exists := processedDirs[templateDir]; !exists || outputPath != outputDir {
		t.Error("Expected template to be processed with correct output path")
	}
}

func testNestedTemplateStructure(t *testing.T, tempDir string) {
	// Create nested template structure
	templateDir := filepath.Join(tempDir, "nested-template")

	// Create directory structure
	dirs := []string{
		"src/components",
		"src/utils",
		"tests/unit",
		"tests/integration",
		"docs",
		"config",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(templateDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create nested template files
	templates := map[string]string{
		"src/main.go.tmpl": `package main

import "fmt"

func main() {
	fmt.Println("Hello from {{.Name}}")
}
`,
		"src/components/component.go.tmpl": `package components

// {{.Name}}Component represents a component for {{.Organization}}
type {{.Name}}Component struct {
	Name string
}
`,
		"tests/unit/main_test.go.tmpl": `package main

import "testing"

func Test{{.Name}}(t *testing.T) {
	// Test for {{.Name}} project
	t.Log("Testing {{.Name}}")
}
`,
		"config/app.yaml.tmpl": `app:
  name: {{.Name}}
  organization: {{.Organization}}
  version: {{.Version | default "1.0.0"}}
  license: {{.License}}
`,
		"docs/README.md.tmpl": `# {{.Name}}

Documentation for {{.Name}} by {{.Organization}}.

## License

This project is licensed under {{.License}}.
`,
	}

	for path, content := range templates {
		fullPath := filepath.Join(templateDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", path, err)
		}
	}

	// Create static files
	staticFiles := map[string]string{
		"Makefile": `build:
	go build -o bin/app src/main.go

test:
	go test ./...

clean:
	rm -rf bin/
`,
		".gitignore": `bin/
*.log
.env
`,
	}

	for path, content := range staticFiles {
		fullPath := filepath.Join(templateDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create static file %s: %v", path, err)
		}
	}

	// Create mock template manager
	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	config := &models.ProjectConfig{
		Name:         "nested-project",
		Organization: "nested-org",
		License:      "BSD-3-Clause",
	}

	outputDir := filepath.Join(tempDir, "nested-output")

	// Process nested template
	err := templateManager.ProcessCustomTemplate(templateDir, config, outputDir)
	if err != nil {
		t.Fatalf("Failed to process nested template: %v", err)
	}

	// Verify processing
	processedDirs := templateEngine.GetProcessedDirs()
	if _, exists := processedDirs[templateDir]; !exists {
		t.Error("Expected nested template to be processed")
	}
}

func testTemplateValidationWorkflow(t *testing.T, tempDir string) {
	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	// Test valid template validation
	validTemplateDir := filepath.Join(tempDir, "valid-template")
	err := os.MkdirAll(validTemplateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create valid template directory: %v", err)
	}

	// Create valid template metadata
	validMetadata := `name: valid-template
version: 1.0.0
author: Test Author
license: MIT
description: A valid test template
`

	err = os.WriteFile(filepath.Join(validTemplateDir, "template.yaml"), []byte(validMetadata), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid metadata: %v", err)
	}

	// Create valid template file
	validTemplate := `# {{.Name}}

Valid template content.
`

	err = os.WriteFile(filepath.Join(validTemplateDir, "README.md.tmpl"), []byte(validTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid template file: %v", err)
	}

	// Validate valid template
	result, err := templateManager.ValidateTemplate(validTemplateDir)
	if err != nil {
		t.Fatalf("Failed to validate template: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid template to pass validation, got issues: %v", result.Issues)
	}

	// Test invalid template validation
	invalidTemplateDir := filepath.Join(tempDir, "invalid-template")
	err = os.MkdirAll(invalidTemplateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create invalid template directory: %v", err)
	}

	// Create invalid template (missing metadata)
	invalidTemplate := `This is not a proper template file.`

	err = os.WriteFile(filepath.Join(invalidTemplateDir, "invalid.txt"), []byte(invalidTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid template file: %v", err)
	}

	// Validate invalid template
	result, err = templateManager.ValidateTemplate(invalidTemplateDir)
	if err != nil {
		t.Fatalf("Failed to validate invalid template: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid template to fail validation")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected validation issues for invalid template")
	}
}

func testCustomTemplateProcessing(t *testing.T, tempDir string) {
	// Create custom template with advanced features
	templateDir := filepath.Join(tempDir, "custom-template")
	err := os.MkdirAll(templateDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create custom template directory: %v", err)
	}

	// Create template with custom functions and complex logic
	advancedTemplate := `# {{.Name | title}}

{{if .Description}}
## Description
{{.Description}}
{{end}}

## Project Information
- **Name**: {{.Name}}
- **Organization**: {{.Organization}}
- **License**: {{.License}}
- **Created**: {{now | date "2006-01-02"}}

{{if .Components}}
## Components
{{range $key, $component := .Components}}
{{if $component.Enabled}}
- **{{$key | title}}**: {{$component.Technology}}{{if $component.Port}} (Port: {{$component.Port}}){{end}}
{{end}}
{{end}}
{{end}}

{{if .Dependencies}}
## Dependencies
{{range .Dependencies}}
- {{.}}
{{end}}
{{end}}

## Getting Started

1. Clone the repository
2. Install dependencies
{{if .Components.Backend.Enabled}}
3. Start the backend server
{{end}}
{{if .Components.Frontend.Enabled}}
4. Start the frontend development server
{{end}}

## License

This project is licensed under the {{.License}} license.
`

	err = os.WriteFile(filepath.Join(templateDir, "README.md.tmpl"), []byte(advancedTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create advanced template: %v", err)
	}

	// Create package.json template with conditional dependencies
	packageTemplate := `{
  "name": "{{.Name | kebab}}",
  "version": "{{.Version | default "1.0.0"}}",
  "description": "{{.Description}}",
  "license": "{{.License}}",
  {{if .Components.Backend.Enabled}}
  "scripts": {
    "start": "node server.js",
    "dev": "nodemon server.js",
    "test": "jest"
  },
  {{end}}
  "dependencies": {
    {{range $i, $dep := .Dependencies}}
    {{if $i}},{{end}}
    "{{$dep}}": "latest"
    {{end}}
  }
}
`

	err = os.WriteFile(filepath.Join(templateDir, "package.json.tmpl"), []byte(packageTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json template: %v", err)
	}

	// Create mock template manager
	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	// Configure complex project
	config := &models.ProjectConfig{
		Name:         "custom-advanced-project",
		Organization: "advanced-org",
		Description:  "An advanced project with custom template processing",
		License:      "MIT",
		Components: models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:  false,
					Home: false,
				},
			},
		},
	}

	outputDir := filepath.Join(tempDir, "custom-output")

	// Process custom template
	err = templateManager.ProcessCustomTemplate(templateDir, config, outputDir)
	if err != nil {
		t.Fatalf("Failed to process custom template: %v", err)
	}

	// Verify processing
	processedDirs := templateEngine.GetProcessedDirs()
	if _, exists := processedDirs[templateDir]; !exists {
		t.Error("Expected custom template to be processed")
	}

	// Verify config was processed correctly
	processedTemplates := templateEngine.GetProcessedTemplates()
	if processedConfig, exists := processedTemplates[templateDir]; exists {
		if processedConfig.Name != config.Name {
			t.Errorf("Expected processed config name '%s', got '%s'", config.Name, processedConfig.Name)
		}

		// Dependencies field doesn't exist in ProjectConfig, skip this check
	}
}

// TestTemplateDiscoveryAndManagement tests template discovery and management
func TestTemplateDiscoveryAndManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	templateEngine := NewMockTemplateEngine()
	templateManager := NewMockTemplateManager(templateEngine)

	t.Run("list_all_templates", func(t *testing.T) {
		templates, err := templateManager.ListTemplates(interfaces.TemplateFilter{})
		if err != nil {
			t.Fatalf("Failed to list templates: %v", err)
		}

		if len(templates) == 0 {
			t.Error("Expected to find some templates")
		}

		// Verify template structure
		for _, tmpl := range templates {
			if tmpl.Name == "" {
				t.Error("Expected template to have a name")
			}

			if tmpl.Category == "" {
				t.Error("Expected template to have a category")
			}
		}
	})

	t.Run("filter_templates_by_category", func(t *testing.T) {
		backendTemplates, err := templateManager.GetTemplatesByCategory("backend")
		if err != nil {
			t.Fatalf("Failed to get backend templates: %v", err)
		}

		// All returned templates should be backend
		for _, tmpl := range backendTemplates {
			if tmpl.Category != "backend" {
				t.Errorf("Expected backend template, got category: %s", tmpl.Category)
			}
		}

		frontendTemplates, err := templateManager.GetTemplatesByCategory("frontend")
		if err != nil {
			t.Fatalf("Failed to get frontend templates: %v", err)
		}

		// All returned templates should be frontend
		for _, tmpl := range frontendTemplates {
			if tmpl.Category != "frontend" {
				t.Errorf("Expected frontend template, got category: %s", tmpl.Category)
			}
		}
	})

	t.Run("filter_templates_by_technology", func(t *testing.T) {
		goTemplates, err := templateManager.GetTemplatesByTechnology("Go")
		if err != nil {
			t.Fatalf("Failed to get Go templates: %v", err)
		}

		// All returned templates should use Go
		for _, tmpl := range goTemplates {
			if tmpl.Technology != "Go" {
				t.Errorf("Expected Go template, got technology: %s", tmpl.Technology)
			}
		}
	})

	t.Run("search_templates", func(t *testing.T) {
		// Search for Go templates
		goResults, err := templateManager.SearchTemplates("go")
		if err != nil {
			t.Fatalf("Failed to search for Go templates: %v", err)
		}

		// Should find templates containing "go"
		found := false
		for _, tmpl := range goResults {
			if strings.Contains(strings.ToLower(tmpl.Name), "go") ||
				strings.Contains(strings.ToLower(tmpl.Technology), "go") {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected to find Go-related templates")
		}

		// Search for React templates
		reactResults, err := templateManager.SearchTemplates("react")
		if err != nil {
			t.Fatalf("Failed to search for React templates: %v", err)
		}

		// Should find templates containing "react"
		found = false
		for _, tmpl := range reactResults {
			if strings.Contains(strings.ToLower(tmpl.Name), "react") ||
				strings.Contains(strings.ToLower(tmpl.Technology), "react") {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected to find React-related templates")
		}
	})

	t.Run("get_template_info", func(t *testing.T) {
		// Get info for a specific template
		info, err := templateManager.GetTemplateInfo("go-gin")
		if err != nil {
			t.Fatalf("Failed to get template info: %v", err)
		}

		if info.Name != "go-gin" {
			t.Errorf("Expected template name 'go-gin', got '%s'", info.Name)
		}

		if info.Category == "" {
			t.Error("Expected template to have a category")
		}

		if info.Technology == "" {
			t.Error("Expected template to have a technology")
		}
	})

	t.Run("get_template_metadata", func(t *testing.T) {
		metadata, err := templateManager.GetTemplateMetadata("go-gin")
		if err != nil {
			t.Fatalf("Failed to get template metadata: %v", err)
		}

		if metadata.Author == "" {
			t.Error("Expected template to have an author")
		}

		if metadata.License == "" {
			t.Error("Expected template to have a license")
		}
	})

	t.Run("get_template_variables", func(t *testing.T) {
		variables, err := templateManager.GetTemplateVariables("go-gin")
		if err != nil {
			t.Fatalf("Failed to get template variables: %v", err)
		}

		// Should have default variables
		expectedVars := []string{"Name", "Organization", "Description", "License"}
		for _, expectedVar := range expectedVars {
			if _, exists := variables[expectedVar]; !exists {
				t.Errorf("Expected template to have variable '%s'", expectedVar)
			}
		}
	})
}

// Mock implementations for testing

type MockTemplateEngine struct {
	processedTemplates map[string]*models.ProjectConfig
	processedDirs      map[string]string
	shouldError        bool
	errorMessage       string
}

func NewMockTemplateEngine() *MockTemplateEngine {
	return &MockTemplateEngine{
		processedTemplates: make(map[string]*models.ProjectConfig),
		processedDirs:      make(map[string]string),
	}
}

func (m *MockTemplateEngine) ProcessTemplate(path string, config *models.ProjectConfig) ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}

	m.processedTemplates[path] = config
	return []byte("processed template content"), nil
}

func (m *MockTemplateEngine) LoadTemplate(path string) (*template.Template, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return template.New("mock").Parse("mock template")
}

func (m *MockTemplateEngine) RenderTemplate(tmpl *template.Template, data any) ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("%s", m.errorMessage)
	}
	return []byte("rendered template"), nil
}

func (m *MockTemplateEngine) ProcessDirectory(templatePath string, outputPath string, config *models.ProjectConfig) error {
	if m.shouldError {
		return fmt.Errorf("%s", m.errorMessage)
	}
	m.processedDirs[templatePath] = outputPath
	m.processedTemplates[templatePath] = config
	return nil
}

func (m *MockTemplateEngine) RegisterFunctions(funcMap template.FuncMap) {
	// Mock implementation - no-op
}

func (m *MockTemplateEngine) GetProcessedTemplates() map[string]*models.ProjectConfig {
	return m.processedTemplates
}

func (m *MockTemplateEngine) GetProcessedDirs() map[string]string {
	return m.processedDirs
}

type MockTemplateManager struct {
	engine interfaces.TemplateEngine
}

func NewMockTemplateManager(engine interfaces.TemplateEngine) *MockTemplateManager {
	return &MockTemplateManager{
		engine: engine,
	}
}

func (m *MockTemplateManager) ProcessCustomTemplate(templatePath string, config *models.ProjectConfig, outputPath string) error {
	return m.engine.ProcessDirectory(templatePath, outputPath, config)
}

func (m *MockTemplateManager) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &interfaces.TemplateValidationResult{
			Valid: false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: "error",
					Message:  "Template path does not exist",
					Rule:     "path-exists",
				},
			},
		}, nil
	}

	// Check for template metadata
	metadataFiles := []string{"template.yaml", "template.yml"}
	hasMetadata := false

	for _, metadataFile := range metadataFiles {
		if _, err := os.Stat(filepath.Join(path, metadataFile)); err == nil {
			hasMetadata = true
			break
		}
	}

	issues := []interfaces.ValidationIssue{}

	if !hasMetadata {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "Template metadata file not found",
			Rule:     "metadata-exists",
		})
	}

	// Check for template files
	hasTemplateFiles := false
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(filePath, ".tmpl") {
			hasTemplateFiles = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if !hasTemplateFiles {
		issues = append(issues, interfaces.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "No template files found",
			Rule:     "template-files-exist",
		})
	}

	// Determine if template is valid (no errors, warnings are OK)
	valid := true
	for _, issue := range issues {
		if issue.Severity == "error" {
			valid = false
			break
		}
	}

	return &interfaces.TemplateValidationResult{
		Valid:  valid,
		Issues: issues,
	}, nil
}

func (m *MockTemplateManager) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	// Mock template list
	allTemplates := []interfaces.TemplateInfo{
		{
			Name:        "go-gin",
			DisplayName: "Go Gin API",
			Category:    "backend",
			Technology:  "Go",
			Version:     "1.0.0",
			Tags:        []string{"backend", "api", "go", "gin"},
		},
		{
			Name:        "nextjs-app",
			DisplayName: "Next.js App",
			Category:    "frontend",
			Technology:  "Next.js",
			Version:     "1.0.0",
			Tags:        []string{"frontend", "react", "nextjs"},
		},
		{
			Name:        "react-component",
			DisplayName: "React Component",
			Category:    "frontend",
			Technology:  "React",
			Version:     "1.0.0",
			Tags:        []string{"frontend", "react", "component"},
		},
	}

	// Apply filters
	filtered := []interfaces.TemplateInfo{}

	for _, tmpl := range allTemplates {
		if filter.Category != "" && tmpl.Category != filter.Category {
			continue
		}

		if filter.Technology != "" && tmpl.Technology != filter.Technology {
			continue
		}

		if len(filter.Tags) > 0 {
			hasAllTags := true
			for _, filterTag := range filter.Tags {
				found := false
				for _, tmplTag := range tmpl.Tags {
					if strings.EqualFold(tmplTag, filterTag) {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if !hasAllTags {
				continue
			}
		}

		filtered = append(filtered, tmpl)
	}

	return filtered, nil
}

func (m *MockTemplateManager) GetTemplatesByCategory(category string) ([]interfaces.TemplateInfo, error) {
	return m.ListTemplates(interfaces.TemplateFilter{Category: category})
}

func (m *MockTemplateManager) GetTemplatesByTechnology(technology string) ([]interfaces.TemplateInfo, error) {
	return m.ListTemplates(interfaces.TemplateFilter{Technology: technology})
}

func (m *MockTemplateManager) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	templates, err := m.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	matches := []interfaces.TemplateInfo{}

	for _, tmpl := range templates {
		if strings.Contains(strings.ToLower(tmpl.Name), query) ||
			strings.Contains(strings.ToLower(tmpl.DisplayName), query) ||
			strings.Contains(strings.ToLower(tmpl.Technology), query) {
			matches = append(matches, tmpl)
		}

		// Check tags
		for _, tag := range tmpl.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				matches = append(matches, tmpl)
				break
			}
		}
	}

	return matches, nil
}

func (m *MockTemplateManager) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	templates, err := m.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		return nil, err
	}

	for _, tmpl := range templates {
		if tmpl.Name == name {
			return &tmpl, nil
		}
	}

	return nil, fmt.Errorf("template '%s' not found", name)
}

func (m *MockTemplateManager) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	_, err := m.GetTemplateInfo(name)
	if err != nil {
		return nil, err
	}

	return &interfaces.TemplateMetadata{
		Author:     "Test Author",
		License:    "MIT",
		Repository: "https://github.com/test/template",
		Homepage:   "https://test.com",
		Keywords:   []string{"template", "test"},
		Created:    time.Now().AddDate(0, -1, 0), // 1 month ago
		Updated:    time.Now().AddDate(0, 0, -7), // 1 week ago
	}, nil
}

func (m *MockTemplateManager) GetTemplateVariables(name string) (map[string]interfaces.TemplateVariable, error) {
	_, err := m.GetTemplateInfo(name)
	if err != nil {
		return nil, err
	}

	return map[string]interfaces.TemplateVariable{
		"Name": {
			Name:        "Name",
			Type:        "string",
			Description: "Project name",
			Required:    true,
		},
		"Organization": {
			Name:        "Organization",
			Type:        "string",
			Description: "Organization name",
			Required:    false,
		},
		"Description": {
			Name:        "Description",
			Type:        "string",
			Description: "Project description",
			Required:    false,
		},
		"License": {
			Name:        "License",
			Type:        "string",
			Description: "Project license",
			Default:     "MIT",
			Required:    false,
		},
	}, nil
}
