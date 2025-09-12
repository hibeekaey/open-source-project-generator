package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/version"
)

func TestTemplateEngineWithVersionManager(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Create a simple template file
	templateContent := `{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "dependencies": {
    "next": "^{{nextjsVersion .}}",
    "react": "^{{reactVersion .}}",
    "typescript": "^{{packageVersion . "typescript"}}"
  },
  "engines": {
    "node": ">=18.0.0"
  }
}`

	templatePath := filepath.Join(tempDir, "package.json.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create version cache and manager
	cache := version.NewMemoryCache(time.Hour)
	versionManager := version.NewManager(cache)

	// Create template engine with version manager
	engine := NewEngineWithVersionManager(versionManager)

	// Create project config
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Versions: &models.VersionConfig{
			Node:   "20.11.0",
			NextJS: "15.0.0",
			React:  "18.2.0",
			Packages: map[string]string{
				"typescript": "5.3.0",
			},
			UpdatedAt: time.Now(),
		},
	}

	// Process the template
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	resultStr := string(result)

	// Verify that versions were correctly substituted
	if !strings.Contains(resultStr, `"next": "^15.0.0"`) {
		t.Errorf("Expected Next.js version 15.0.0, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, `"react": "^18.2.0"`) {
		t.Errorf("Expected React version 18.2.0, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, `"typescript": "^5.3.0"`) {
		t.Errorf("Expected TypeScript version 5.3.0, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, `"name": "test-project"`) {
		t.Errorf("Expected project name to be preserved, got: %s", resultStr)
	}
}

func TestTemplateEngineVersionFallbacks(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Create a template that uses version functions
	templateContent := `Node: {{nodeVersion .}}
Go: {{goVersion .}}
Next.js: {{nextjsVersion .}}
React: {{reactVersion .}}
Unknown Package: {{packageVersion . "unknown-package"}}`

	templatePath := filepath.Join(tempDir, "versions.txt.tmpl")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Create template engine without version manager
	engine := NewEngine()

	// Create minimal project config without versions
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
	}

	// Process the template
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("Failed to process template: %v", err)
	}

	resultStr := string(result)

	// Verify that fallback versions are used
	if !strings.Contains(resultStr, "Node: 20.11.0") {
		t.Errorf("Expected fallback Node version, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, "Go: 1.22.0") {
		t.Errorf("Expected fallback Go version, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, "Next.js: 15.0.0") {
		t.Errorf("Expected fallback Next.js version, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, "React: 18.2.0") {
		t.Errorf("Expected fallback React version, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, "Unknown Package: latest") {
		t.Errorf("Expected 'latest' for unknown package, got: %s", resultStr)
	}
}

func TestTemplateEngineDirectoryProcessing(t *testing.T) {
	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "templates")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory structure
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create a package.json template
	packageTemplate := `{
  "name": "{{.Name}}",
  "dependencies": {
    "next": "^{{nextjsVersion .}}",
    "react": "^{{reactVersion .}}"
  }
}`

	packagePath := filepath.Join(templateDir, "package.json.tmpl")
	if err := os.WriteFile(packagePath, []byte(packageTemplate), 0644); err != nil {
		t.Fatalf("Failed to write package template: %v", err)
	}

	// Create a non-template file
	readmeContent := "# {{.Name}}\n\nThis is a test project."
	readmePath := filepath.Join(templateDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatalf("Failed to write README: %v", err)
	}

	// Create version manager and template engine
	cache := version.NewMemoryCache(time.Hour)
	versionManager := version.NewManager(cache)
	engine := NewEngineWithVersionManager(versionManager)

	// Create project config
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Versions: &models.VersionConfig{
			NextJS: "15.0.0",
			React:  "18.2.0",
		},
	}

	// Process the directory
	err := engine.ProcessDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("Failed to process directory: %v", err)
	}

	// Verify package.json was processed
	packageOutput := filepath.Join(outputDir, "package.json")
	packageContent, err := os.ReadFile(packageOutput)
	if err != nil {
		t.Fatalf("Failed to read processed package.json: %v", err)
	}

	packageStr := string(packageContent)
	if !strings.Contains(packageStr, `"next": "^15.0.0"`) {
		t.Errorf("Expected processed Next.js version in package.json, got: %s", packageStr)
	}

	// Verify README.md was copied as-is (not processed as template)
	readmeOutput := filepath.Join(outputDir, "README.md")
	readmeOutputContent, err := os.ReadFile(readmeOutput)
	if err != nil {
		t.Fatalf("Failed to read copied README.md: %v", err)
	}

	if string(readmeOutputContent) != readmeContent {
		t.Errorf("README.md should be copied as-is, got: %s", string(readmeOutputContent))
	}
}
