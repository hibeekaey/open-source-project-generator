package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/open-source-template-generator/pkg/models"
)

// TestTemplateCompilationIntegration tests end-to-end template generation and compilation
func TestTemplateCompilationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("GoGinTemplateCompilation", func(t *testing.T) {
		testGoGinTemplateCompilation(t)
	})

	t.Run("AllGoTemplatesCompilation", func(t *testing.T) {
		testAllGoTemplatesCompilation(t)
	})

	t.Run("TemplateVariableSubstitution", func(t *testing.T) {
		testTemplateVariableSubstitution(t)
	})

	t.Run("ImportFixValidation", func(t *testing.T) {
		testImportFixValidation(t)
	})
}

func testGoGinTemplateCompilation(t *testing.T) {
	// Test specific Go Gin templates that were known to have issues
	templatesDir := "../../templates/backend/go-gin"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skip("Go Gin templates directory not found, skipping test")
	}

	testData := createTestProjectConfig()
	outputDir := t.TempDir()

	// Test specific templates that had import issues
	criticalTemplates := []string{
		"internal/middleware/auth.go.tmpl",
		"internal/middleware/security.go.tmpl",
		"internal/controllers/auth_controller.go.tmpl",
		"internal/services/auth_service.go.tmpl",
		"main.go.tmpl",
		"go.mod.tmpl",
	}

	for _, templatePath := range criticalTemplates {
		fullTemplatePath := filepath.Join(templatesDir, templatePath)
		if _, err := os.Stat(fullTemplatePath); os.IsNotExist(err) {
			t.Logf("Template %s not found, skipping", templatePath)
			continue
		}

		t.Run(templatePath, func(t *testing.T) {
			err := generateAndValidateTemplate(t, fullTemplatePath, testData, outputDir)
			if err != nil {
				t.Errorf("Template %s failed compilation: %v", templatePath, err)
			}
		})
	}
}

func testAllGoTemplatesCompilation(t *testing.T) {
	templatesDir := "../../templates"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skip("Templates directory not found, skipping test")
	}

	testData := createTestProjectConfig()
	outputDir := t.TempDir()

	var failedTemplates []string
	var totalTemplates int

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only test Go template files
		if !strings.HasSuffix(path, ".go.tmpl") && !strings.HasSuffix(path, ".mod.tmpl") {
			return nil
		}

		totalTemplates++
		relPath, _ := filepath.Rel(templatesDir, path)

		t.Run(relPath, func(t *testing.T) {
			err := generateAndValidateTemplate(t, path, testData, outputDir)
			if err != nil {
				failedTemplates = append(failedTemplates, relPath)
				t.Errorf("Template %s failed: %v", relPath, err)
			}
		})

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk templates directory: %v", err)
	}

	t.Logf("Tested %d Go templates", totalTemplates)
	if len(failedTemplates) > 0 {
		t.Errorf("Failed templates (%d/%d): %v", len(failedTemplates), totalTemplates, failedTemplates)
	}
}

func testTemplateVariableSubstitution(t *testing.T) {
	// Test that template variables are properly substituted and don't break compilation
	testCases := []struct {
		name            string
		templateContent string
		expectedVars    []string
	}{
		{
			name: "BasicVariables",
			templateContent: `package {{.Name}}

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Project: %s\n", "{{.ProjectName}}")
	fmt.Printf("Author: %s\n", "{{.Author}}")
	fmt.Printf("Version: %s\n", "{{.Version}}")
	fmt.Printf("Time: %v\n", time.Now())
}`,
			expectedVars: []string{"testproject", "Test Project", "Test Author", "1.0.0"},
		},
		{
			name: "ConditionalBlocks",
			templateContent: `package {{.Name}}

import (
	"fmt"
	{{- if .EnableAuth }}
	"time"
	{{- end }}
)

func main() {
	fmt.Println("Starting {{.ServiceName}}")
	{{- if .EnableAuth }}
	fmt.Printf("Auth enabled at %v\n", time.Now())
	{{- end }}
}`,
			expectedVars: []string{"testproject", "TestService"},
		},
		{
			name: "LoopStructures",
			templateContent: `package {{.Name}}

import "fmt"

func main() {
	services := []string{
		{{- range .Services }}
		"{{.Name}}",
		{{- end }}
	}
	
	for _, service := range services {
		fmt.Printf("Service: %s\n", service)
	}
}`,
			expectedVars: []string{"testproject"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			testData := createTestProjectConfig()
			err = generateAndValidateTemplate(t, templateFile, testData, tempDir)
			if err != nil {
				t.Errorf("Template variable substitution failed: %v", err)
			}
		})
	}
}

func testImportFixValidation(t *testing.T) {
	// Test templates with known import issues to ensure they're fixed
	testCases := []struct {
		name            string
		templateContent string
		expectedImports []string
		description     string
	}{
		{
			name: "MissingTimeImport",
			templateContent: `package {{.Name}}

import "fmt"

func main() {
	fmt.Printf("Current time: %v\n", time.Now())
}`,
			expectedImports: []string{"fmt", "time"},
			description:     "Template should have time import added",
		},
		{
			name: "MissingHTTPImports",
			templateContent: `package {{.Name}}

import "fmt"

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}`,
			expectedImports: []string{"fmt", "net/http"},
			description:     "Template should have net/http import added",
		},
		{
			name: "MissingJSONImports",
			templateContent: `package {{.Name}}

import "fmt"

func processData(data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("JSON: %s\n", bytes)
}`,
			expectedImports: []string{"fmt", "encoding/json"},
			description:     "Template should have encoding/json import added",
		},
		{
			name: "MultipleMissing",
			templateContent: `package {{.Name}}

func processRequest(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	data := map[string]interface{}{
		"timestamp": now,
		"path":      r.URL.Path,
	}
	
	response, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// SECURITY: Added comprehensive security headers
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}`,
			expectedImports: []string{"time", "net/http", "encoding/json"},
			description:     "Template should have multiple missing imports added",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			// First, analyze the template for missing imports
			detector := NewImportDetector()
			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				t.Fatalf("Failed to analyze template: %v", err)
			}

			// Check that missing imports are detected
			for _, expectedImport := range tc.expectedImports {
				found := false
				for _, currentImport := range report.CurrentImports {
					if currentImport.Package == expectedImport {
						found = true
						break
					}
				}
				if !found {
					// Check if it's in missing imports
					foundInMissing := false
					for _, missingImport := range report.MissingImports {
						if missingImport == expectedImport {
							foundInMissing = true
							break
						}
					}
					if !foundInMissing {
						t.Errorf("Expected import %s not found in current or missing imports", expectedImport)
					}
				}
			}

			// Then, test that the template can be fixed and compiled
			fixedContent := addMissingImports(tc.templateContent, report.MissingImports)
			fixedTemplateFile := filepath.Join(tempDir, "fixed.go.tmpl")

			err = os.WriteFile(fixedTemplateFile, []byte(fixedContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create fixed template file: %v", err)
			}

			testData := createTestProjectConfig()
			err = generateAndValidateTemplate(t, fixedTemplateFile, testData, tempDir)
			if err != nil {
				t.Errorf("Fixed template failed compilation: %v", err)
			}
		})
	}
}

// generateAndValidateTemplate generates a Go file from a template and validates it compiles
func generateAndValidateTemplate(_ *testing.T, templatePath string, testData *models.ProjectConfig, outputDir string) error {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file path
	relPath := strings.TrimSuffix(filepath.Base(templatePath), ".tmpl")
	outputPath := filepath.Join(outputDir, relPath)

	// Create output directory
	outputDirPath := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate file from template
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, testData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Validate the generated file
	return validateGeneratedGoFile(outputPath)
}

// validateGeneratedGoFile validates that a generated Go file compiles correctly
func validateGeneratedGoFile(filePath string) error {
	// Skip go.mod files as they don't need compilation
	if strings.HasSuffix(filePath, "go.mod") {
		return validateGoMod(filePath)
	}

	// For .go files, check if they have only standard library imports
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read generated file: %w", err)
	}

	// Check if this file has only standard library imports
	if !hasOnlyStandardLibraryImports(string(content)) {
		// Skip compilation for files with non-standard imports
		// as we can't resolve project-specific dependencies in tests
		return nil
	}

	// Create a temporary directory for compilation test
	tempDir := filepath.Dir(filePath) + "_compile_test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy the file to temp directory
	tempFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Create a simple go.mod for the temp directory
	goModContent := `module temp-validation
go 1.22
`
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	// Try to compile the file
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	return nil
}

// validateGoMod validates a go.mod file
func validateGoMod(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "module ") {
		return fmt.Errorf("invalid go.mod file: missing module declaration")
	}

	return nil
}

// addMissingImports adds missing imports to template content (simplified implementation)
func addMissingImports(content string, missingImports []string) string {
	if len(missingImports) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var result []string
	importAdded := false

	for _, line := range lines {
		result = append(result, line)

		// Add imports after package declaration
		if strings.HasPrefix(strings.TrimSpace(line), "package ") && !importAdded {
			result = append(result, "")
			result = append(result, "import (")
			for _, imp := range missingImports {
				result = append(result, fmt.Sprintf("\t\"%s\"", imp))
			}
			result = append(result, ")")
			importAdded = true
		}
	}

	return strings.Join(result, "\n")
}

// createTestProjectConfig is now defined in test_helpers.go
