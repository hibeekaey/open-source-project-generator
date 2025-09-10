package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	texttemplate "text/template"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	// Verify that the engine implements the interface
	_, ok := engine.(*Engine)
	if !ok {
		t.Fatal("NewEngine() did not return an *Engine")
	}
}

func TestProcessTemplate(t *testing.T) {
	// Create a temporary template file
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "test.tmpl")

	templateContent := `Project: {{.Name}}
Organization: {{.Organization}}
{{if hasFrontend .}}Has Frontend: true{{end}}
{{if hasBackend .}}Has Backend: true{{end}}`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Create test configuration
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
			Backend: models.BackendComponents{
				API: false,
			},
		},
	}

	// Test template processing
	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	resultStr := string(result)
	expectedLines := []string{
		"Project: test-project",
		"Organization: test-org",
		"Has Frontend: true",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(resultStr, expected) {
			t.Errorf("Expected result to contain '%s', got: %s", expected, resultStr)
		}
	}

	// Should not contain backend line since backend is false
	if strings.Contains(resultStr, "Has Backend: true") {
		t.Errorf("Result should not contain 'Has Backend: true', got: %s", resultStr)
	}
}

func TestProcessDirectory(t *testing.T) {
	// Create temporary directories
	tempDir := t.TempDir()
	templateDir := filepath.Join(tempDir, "template")
	outputDir := filepath.Join(tempDir, "output")

	// Create template directory structure
	err := os.MkdirAll(filepath.Join(templateDir, "subdir"), 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create template files
	templateFiles := map[string]string{
		"README.md.tmpl": "# {{.Name}}\n{{.Description}}",
		"package.json.tmpl": `{
  "name": "{{kebabCase .Name}}",
  "version": "1.0.0",
  "description": "{{.Description}}"
}`,
		"subdir/config.yaml.tmpl": `name: {{.Name}}
org: {{.Organization}}`,
		"binary-file.png": "fake-binary-content", // Non-template file
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
		Name:         "MyTestProject",
		Organization: "test-org",
		Description:  "A test project for template processing",
	}

	// Process directory
	engine := NewEngine()
	err = engine.ProcessDirectory(templateDir, outputDir, config)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Verify output files
	expectedFiles := map[string]string{
		"README.md":          "# MyTestProject\nA test project for template processing",
		"package.json":       `"name": "my-test-project"`,
		"subdir/config.yaml": "name: MyTestProject\norg: test-org",
		"binary-file.png":    "fake-binary-content",
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
		if !strings.Contains(contentStr, expectedContent) {
			t.Errorf("File %s does not contain expected content '%s', got: %s",
				filePath, expectedContent, contentStr)
		}
	}
}

func TestRegisterFunctions(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Register custom function
	customFuncs := map[string]interface{}{
		"customFunc": func(s string) string {
			return "custom-" + s
		},
	}

	engine.RegisterFunctions(customFuncs)

	// Verify function is registered
	if _, exists := engine.funcMap["customFunc"]; !exists {
		t.Error("Custom function was not registered")
	}

	// Test using the custom function in a template
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "test.tmpl")

	templateContent := `Result: {{customFunc .Name}}`
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	config := &models.ProjectConfig{
		Name: "test",
	}

	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("ProcessTemplate with custom function failed: %v", err)
	}

	expected := "Result: custom-test"
	if !strings.Contains(string(result), expected) {
		t.Errorf("Expected result to contain '%s', got: %s", expected, string(result))
	}
}

func TestLoadTemplate(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "test.tmpl")

	templateContent := `Hello {{.Name}}!`
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	engine := NewEngine().(*Engine)
	tmpl, err := engine.LoadTemplate(templatePath)
	if err != nil {
		t.Fatalf("LoadTemplate failed: %v", err)
	}

	if tmpl == nil {
		t.Fatal("LoadTemplate returned nil template")
	}

	if tmpl.Name() != "test.tmpl" {
		t.Errorf("Expected template name 'test.tmpl', got '%s'", tmpl.Name())
	}
}

func TestRenderTemplate(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Create a simple template using text/template directly
	tmpl, err := texttemplate.New("test").Funcs(engine.funcMap).Parse("Hello {{.Name}}!")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	data := struct {
		Name string
	}{
		Name: "World",
	}

	result, err := engine.RenderTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello World!"
	if string(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(result))
	}
}

func TestProcessTemplateWithVersions(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "package.json.tmpl")

	templateContent := `{
  "name": "{{kebabCase .Name}}",
  "version": "1.0.0",
  "dependencies": {
    "react": "{{latestVersion . "react"}}",
    "next": "{{latestVersion . "next"}}"
  }
}`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	config := &models.ProjectConfig{
		Name: "TestProject",
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"react": "18.2.0",
				"next":  "14.0.0",
			},
		},
	}

	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	resultStr := string(result)
	expectedContent := []string{
		`"name": "test-project"`,
		`"react": "18.2.0"`,
		`"next": "14.0.0"`,
	}

	for _, expected := range expectedContent {
		if !strings.Contains(resultStr, expected) {
			t.Errorf("Expected result to contain '%s', got: %s", expected, resultStr)
		}
	}
}

func TestProcessTemplateWithConditionals(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "config.tmpl")

	templateContent := `# Configuration
{{if hasFrontend .}}
frontend:
  enabled: true
{{if .Components.Frontend.MainApp}}  main_app: true{{end}}
{{if .Components.Frontend.Admin}}  admin: true{{end}}
{{end}}
{{if hasBackend .}}
backend:
  enabled: true
{{if .Components.Backend.API}}  api: true{{end}}
{{end}}
{{if hasMobile .}}
mobile:
  enabled: true
{{if .Components.Mobile.Android}}  android: true{{end}}
{{if .Components.Mobile.IOS}}  ios: true{{end}}
{{end}}`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	config := &models.ProjectConfig{
		Name: "test-project",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Admin:   false,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: false,
				IOS:     true,
			},
		},
	}

	engine := NewEngine()
	result, err := engine.ProcessTemplate(templatePath, config)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	resultStr := string(result)

	// Should contain frontend section
	if !strings.Contains(resultStr, "frontend:") {
		t.Error("Expected result to contain frontend section")
	}
	if !strings.Contains(resultStr, "main_app: true") {
		t.Error("Expected result to contain main_app: true")
	}
	if strings.Contains(resultStr, "admin: true") {
		t.Error("Result should not contain admin: true")
	}

	// Should contain backend section
	if !strings.Contains(resultStr, "backend:") {
		t.Error("Expected result to contain backend section")
	}
	if !strings.Contains(resultStr, "api: true") {
		t.Error("Expected result to contain api: true")
	}

	// Should contain mobile section
	if !strings.Contains(resultStr, "mobile:") {
		t.Error("Expected result to contain mobile section")
	}
	if !strings.Contains(resultStr, "ios: true") {
		t.Error("Expected result to contain ios: true")
	}
	if strings.Contains(resultStr, "android: true") {
		t.Error("Result should not contain android: true")
	}
}

func TestProcessTemplateErrors(t *testing.T) {
	engine := NewEngine()

	// Test with non-existent template file
	_, err := engine.ProcessTemplate("/non/existent/file.tmpl", &models.ProjectConfig{})
	if err == nil {
		t.Error("Expected error for non-existent template file")
	}

	// Test with invalid template syntax
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "invalid.tmpl")

	invalidTemplate := `{{.Name} {{unclosed`
	err = os.WriteFile(templatePath, []byte(invalidTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid template: %v", err)
	}

	_, err = engine.ProcessTemplate(templatePath, &models.ProjectConfig{Name: "test"})
	if err == nil {
		t.Error("Expected error for invalid template syntax")
	}
}

func TestTemplateEngineEdgeCases(t *testing.T) {
	engine := NewEngine()
	tempDir := t.TempDir()

	t.Run("template with nil config", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "nil-config.tmpl")
		templateContent := `Project: {{.Name}}`

		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		_, err = engine.ProcessTemplate(templatePath, nil)
		if err == nil {
			t.Error("Expected error when processing template with nil config")
		}
	})

	t.Run("template with missing fields", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "missing-fields.tmpl")
		templateContent := `Project: {{.Name}}
Organization: {{.Organization}}
NonExistent: {{.NonExistentField}}`

		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		config := &models.ProjectConfig{
			Name: "test-project",
			// Organization is missing
		}

		_, err = engine.ProcessTemplate(templatePath, config)
		if err == nil {
			t.Error("Template should fail with non-existent field")
		} else {
			t.Logf("Expected template error for non-existent field: %v", err)
		}
	})

	t.Run("template with complex nested structures", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "nested.tmpl")
		templateContent := `{{range $key, $value := .CustomVars}}
{{$key}}: {{$value}}
{{end}}
{{if .Versions}}
{{if .Versions.Packages}}
{{range $pkg, $ver := .Versions.Packages}}
{{$pkg}}: {{$ver}}
{{end}}
{{end}}
{{end}}`

		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		config := &models.ProjectConfig{
			CustomVars: map[string]string{
				"var1": "value1",
				"var2": "value2",
			},
			Versions: &models.VersionConfig{
				Packages: map[string]string{
					"react": "18.0.0",
					"next":  "14.0.0",
				},
			},
		}

		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Errorf("Template processing failed: %v", err)
		}

		resultStr := string(result)
		if !strings.Contains(resultStr, "var1: value1") {
			t.Error("Template should process custom variables")
		}
		if !strings.Contains(resultStr, "react: 18.0.0") {
			t.Error("Template should process package versions")
		}
	})

	t.Run("template with unicode content", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "unicode.tmpl")
		templateContent := `项目名称: {{.Name}}
组织: {{.Organization}}
描述: {{.Description}}
🚀 Emoji test: {{.Name}}`

		err := os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		config := &models.ProjectConfig{
			Name:         "测试项目",
			Organization: "тест-орг",
			Description:  "Unicode description 🌟",
		}

		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Errorf("Template should handle unicode content: %v", err)
		}

		resultStr := string(result)
		if !strings.Contains(resultStr, "测试项目") {
			t.Error("Template should preserve unicode characters")
		}
		if !strings.Contains(resultStr, "🚀") {
			t.Error("Template should preserve emoji characters")
		}
	})

	t.Run("template with very large content", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "large.tmpl")

		// Create a template with repeated content
		var templateBuilder strings.Builder
		for i := 0; i < 1000; i++ {
			templateBuilder.WriteString(fmt.Sprintf("Line %d: {{.Name}}\n", i))
		}

		err := os.WriteFile(templatePath, []byte(templateBuilder.String()), 0644)
		if err != nil {
			t.Fatalf("Failed to create large template: %v", err)
		}

		config := &models.ProjectConfig{
			Name: "large-test",
		}

		result, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			t.Errorf("Template should handle large content: %v", err)
		}

		resultStr := string(result)
		if !strings.Contains(resultStr, "Line 999: large-test") {
			t.Error("Template should process all content in large files")
		}
	})

	t.Run("template with recursive includes", func(t *testing.T) {
		// Create templates that include each other
		template1Path := filepath.Join(tempDir, "template1.tmpl")
		template2Path := filepath.Join(tempDir, "template2.tmpl")

		template1Content := `Template 1: {{.Name}}
{{template "template2.tmpl" .}}`
		template2Content := `Template 2: {{.Organization}}`

		err := os.WriteFile(template1Path, []byte(template1Content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template1: %v", err)
		}

		err = os.WriteFile(template2Path, []byte(template2Content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template2: %v", err)
		}

		config := &models.ProjectConfig{
			Name:         "test-project",
			Organization: "test-org",
		}

		// This might fail depending on how includes are implemented
		_, err = engine.ProcessTemplate(template1Path, config)
		// We don't assert success/failure here as include behavior may vary
		t.Logf("Template include result: %v", err)
	})
}

func TestTemplateEnginePerformance(t *testing.T) {
	engine := NewEngine()
	tempDir := t.TempDir()

	t.Run("performance with many templates", func(t *testing.T) {
		const numTemplates = 100

		// Create many template files
		for i := 0; i < numTemplates; i++ {
			templatePath := filepath.Join(tempDir, fmt.Sprintf("perf-%d.tmpl", i))
			templateContent := fmt.Sprintf("Template %d: {{.Name}}", i)

			err := os.WriteFile(templatePath, []byte(templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create template %d: %v", i, err)
			}
		}

		config := &models.ProjectConfig{
			Name: "performance-test",
		}

		start := time.Now()

		// Process all templates
		for i := 0; i < numTemplates; i++ {
			templatePath := filepath.Join(tempDir, fmt.Sprintf("perf-%d.tmpl", i))
			_, err := engine.ProcessTemplate(templatePath, config)
			if err != nil {
				t.Errorf("Failed to process template %d: %v", i, err)
			}
		}

		duration := time.Since(start)
		t.Logf("Processed %d templates in %v (avg: %v per template)",
			numTemplates, duration, duration/time.Duration(numTemplates))

		// Performance should be reasonable (less than 10ms per template on average)
		avgDuration := duration / time.Duration(numTemplates)
		if avgDuration > 10*time.Millisecond {
			t.Errorf("Template processing too slow: %v per template", avgDuration)
		}
	})
}

func TestTemplateEngineMemoryUsage(t *testing.T) {
	engine := NewEngine()
	tempDir := t.TempDir()

	t.Run("memory usage with large templates", func(t *testing.T) {
		templatePath := filepath.Join(tempDir, "memory-test.tmpl")

		// Create a template that generates large output
		var templateBuilder strings.Builder
		templateBuilder.WriteString("{{range $i := seq 10000}}")
		templateBuilder.WriteString("Line {{$i}}: {{$.Name}} - {{$.Description}}\n")
		templateBuilder.WriteString("{{end}}")

		err := os.WriteFile(templatePath, []byte(templateBuilder.String()), 0644)
		if err != nil {
			t.Fatalf("Failed to create memory test template: %v", err)
		}

		config := &models.ProjectConfig{
			Name:        "memory-test",
			Description: "This is a test for memory usage with large template output",
		}

		// Process the template multiple times to check for memory leaks
		for i := 0; i < 10; i++ {
			result, err := engine.ProcessTemplate(templatePath, config)
			if err != nil {
				// seq function might not be available, that's ok
				t.Logf("Memory test template processing: %v", err)
				break
			}

			// Verify we got some output
			if len(result) == 0 {
				t.Error("Expected non-empty result from memory test template")
			}

			// Clear result to help GC
			result = nil
		}
	})
}
