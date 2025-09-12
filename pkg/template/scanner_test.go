package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateScanner_ScanFrontendTemplates(t *testing.T) {
	// Use the actual templates directory for testing
	templateDir := "../../templates"

	// Check if templates directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skip("Templates directory not found, skipping integration test")
	}

	scanner := NewTemplateScanner(templateDir)
	analysis, err := scanner.ScanFrontendTemplates()

	if err != nil {
		t.Fatalf("Failed to scan frontend templates: %v", err)
	}

	// Verify we found the expected templates
	expectedTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	if len(analysis.Templates) != len(expectedTemplates) {
		t.Errorf("Expected %d templates, got %d", len(expectedTemplates), len(analysis.Templates))
	}

	// Verify each template has basic information
	for _, template := range analysis.Templates {
		if template.Name == "" {
			t.Error("Template name should not be empty")
		}

		if template.Type == "" {
			t.Error("Template type should not be empty")
		}

		if len(template.ConfigFiles) == 0 {
			t.Errorf("Template %s should have configuration files", template.Name)
		}

		// Verify package.json exists
		hasPackageJSON := false
		for _, file := range template.ConfigFiles {
			if file == "package.json.tmpl" {
				hasPackageJSON = true
				break
			}
		}

		if !hasPackageJSON {
			t.Errorf("Template %s should have package.json.tmpl", template.Name)
		}
	}

	// Verify inconsistencies are detected
	t.Logf("Found %d inconsistencies", len(analysis.Inconsistencies))
	for _, inconsistency := range analysis.Inconsistencies {
		t.Logf("Inconsistency: %s - %s", inconsistency.Type, inconsistency.Description)
	}

	// Verify missing files are detected
	t.Logf("Found %d missing files", len(analysis.MissingFiles))
	for _, missing := range analysis.MissingFiles {
		t.Logf("Missing file: %s in template %s", missing.File, missing.Template)
	}

	// Verify version references are found
	if len(analysis.VersionReferences) == 0 {
		t.Error("Should find version references")
	}

	// Verify dependency patterns are analyzed
	if len(analysis.DependencyPatterns) == 0 {
		t.Error("Should find dependency patterns")
	}
}

func TestTemplateScanner_analyzeTemplate(t *testing.T) {
	templateDir := "../../templates"

	// Check if templates directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skip("Templates directory not found, skipping integration test")
	}

	scanner := NewTemplateScanner(templateDir)

	templatePath := filepath.Join(templateDir, "frontend", "nextjs-app")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("nextjs-app template not found, skipping test")
	}

	templateInfo, err := scanner.analyzeTemplate("nextjs-app", templatePath)
	if err != nil {
		t.Fatalf("Failed to analyze template: %v", err)
	}

	if templateInfo.Name != "nextjs-app" {
		t.Errorf("Expected name 'nextjs-app', got '%s'", templateInfo.Name)
	}

	if templateInfo.Type != "application" {
		t.Errorf("Expected type 'application', got '%s'", templateInfo.Type)
	}

	if len(templateInfo.ConfigFiles) == 0 {
		t.Error("Should find configuration files")
	}

	// Check for specific config files
	expectedFiles := []string{"package.json.tmpl", "next.config.js.tmpl", "tailwind.config.js.tmpl"}
	for _, expectedFile := range expectedFiles {
		found := false
		for _, file := range templateInfo.ConfigFiles {
			if file == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find config file: %s", expectedFile)
		}
	}
}

func TestTemplateScanner_isConfigFile(t *testing.T) {
	scanner := &TemplateScanner{}

	testCases := []struct {
		fileName string
		expected bool
	}{
		{"package.json.tmpl", true},
		{"tsconfig.json.tmpl", true},
		{".eslintrc.json.tmpl", true},
		{".prettierrc.tmpl", true},
		{"next.config.js.tmpl", true},
		{"tailwind.config.js.tmpl", true},
		{"vercel.json.tmpl", true},
		{"jest.config.js.tmpl", true},
		{"postcss.config.js.tmpl", true},
		{".gitignore.tmpl", true},
		{".env.local.example.tmpl", true},
		{"README.md.tmpl", false},
		{"index.tsx.tmpl", false},
		{"random.txt", false},
	}

	for _, tc := range testCases {
		result := scanner.isConfigFile(tc.fileName)
		if result != tc.expected {
			t.Errorf("isConfigFile(%s) = %v, expected %v", tc.fileName, result, tc.expected)
		}
	}
}

func TestTemplateScanner_determineTemplateType(t *testing.T) {
	scanner := &TemplateScanner{}

	testCases := []struct {
		name     string
		expected string
	}{
		{"nextjs-app", "application"},
		{"nextjs-home", "landing"},
		{"nextjs-admin", "dashboard"},
		{"unknown-template", "unknown"},
	}

	for _, tc := range testCases {
		result := scanner.determineTemplateType(tc.name)
		if result != tc.expected {
			t.Errorf("determineTemplateType(%s) = %s, expected %s", tc.name, result, tc.expected)
		}
	}
}

func TestTemplateScanner_extractPortFromScript(t *testing.T) {
	scanner := &TemplateScanner{}

	testCases := []struct {
		script   string
		expected string
	}{
		{"next dev", "3000"},
		{"next dev -p 3001", "3001"},
		{"next dev -p 3002", "3002"},
		{"next start -p 4000", "4000"},
		{"npm run dev", "3000"},
	}

	for _, tc := range testCases {
		result := scanner.extractPortFromScript(tc.script)
		if result != tc.expected {
			t.Errorf("extractPortFromScript(%s) = %s, expected %s", tc.script, result, tc.expected)
		}
	}
}
