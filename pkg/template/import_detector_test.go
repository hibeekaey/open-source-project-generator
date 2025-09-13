package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewImportDetector(t *testing.T) {
	detector := NewImportDetector()

	if detector == nil {
		t.Fatal("NewImportDetector returned nil")
	}

	if detector.functionPackageMap == nil {
		t.Fatal("functionPackageMap not initialized")
	}

	if detector.fileSet == nil {
		t.Fatal("fileSet not initialized")
	}

	// Test that some expected mappings exist
	expectedMappings := map[string]string{
		"time.Now":         "time",
		"fmt.Printf":       "fmt",
		"strings.Contains": "strings",
		"json.Marshal":     "encoding/json",
	}

	for function, expectedPackage := range expectedMappings {
		if actualPackage, exists := detector.functionPackageMap[function]; !exists {
			t.Errorf("Expected function %s not found in mapping", function)
		} else if actualPackage != expectedPackage {
			t.Errorf("Function %s mapped to %s, expected %s", function, actualPackage, expectedPackage)
		}
	}
}

func TestPreprocessTemplateContent(t *testing.T) {
	detector := NewImportDetector()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace template variables",
			input:    `package {{.Name}}\n\nfunc main() {\n\tfmt.Println("{{.ProjectName}}")\n}`,
			expected: `package TemplateName\n\nfunc main() {\n\tfmt.Println("TemplateProjectName")\n}`,
		},
		{
			name:     "Remove template directives",
			input:    `{{ if .EnableAuth }}\nfunc auth() {}\n{{ end }}`,
			expected: `\nfunc auth() {}\n`,
		},
		{
			name:     "Replace remaining template expressions",
			input:    `const version = "{{.Version}}"`,
			expected: `const version = "template_placeholder"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := detector.preprocessTemplateContent(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}

func TestIsStandardLibrary(t *testing.T) {
	detector := NewImportDetector()

	testCases := []struct {
		pkg      string
		expected bool
	}{
		{"fmt", true},
		{"time", true},
		{"encoding/json", true},
		{"net/http", true},
		{"path/filepath", true},
		{"github.com/gin-gonic/gin", false},
		{"golang.org/x/crypto", false},
		{"custom/package", false},
	}

	for _, tc := range testCases {
		t.Run(tc.pkg, func(t *testing.T) {
			result := detector.isStandardLibrary(tc.pkg)
			if result != tc.expected {
				t.Errorf("Package %s: expected %v, got %v", tc.pkg, tc.expected, result)
			}
		})
	}
}

func TestExtractFunctionName(t *testing.T) {
	detector := NewImportDetector()

	// This would require creating AST nodes, which is complex for unit tests
	// We'll test this through integration tests instead
	_ = detector
}

func TestFindMissingImports(t *testing.T) {
	detector := NewImportDetector()

	usages := []FunctionUsage{
		{Function: "time.Now", RequiredPackage: "time"},
		{Function: "fmt.Printf", RequiredPackage: "fmt"},
		{Function: "strings.Contains", RequiredPackage: "strings"},
	}

	currentImports := []ImportStatement{
		{Package: "fmt", IsStdLib: true},
	}

	missing := detector.findMissingImports(usages, currentImports)

	expectedMissing := []string{"time", "strings"}
	if len(missing) != len(expectedMissing) {
		t.Errorf("Expected %d missing imports, got %d", len(expectedMissing), len(missing))
	}

	for _, expected := range expectedMissing {
		found := false
		for _, actual := range missing {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected missing import %s not found", expected)
		}
	}
}

func TestAnalyzeTemplateFileIntegration(t *testing.T) {
	// Create a temporary template file for testing
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "test.go.tmpl")

	templateContent := `package {{.Name}}

import (
	"fmt"
)

func main() {
	fmt.Printf("Hello %s\n", "{{.Name}}")
	now := time.Now()
	result := strings.Contains("test", "es")
	_ = now
	_ = result
}
`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	detector := NewImportDetector()
	report, err := detector.AnalyzeTemplateFile(templateFile)

	if err != nil {
		t.Fatalf("AnalyzeTemplateFile failed: %v", err)
	}

	if report == nil {
		t.Fatal("Report is nil")
	}

	// Check that missing imports were detected
	expectedMissing := []string{"time", "strings"}
	if len(report.MissingImports) != len(expectedMissing) {
		t.Errorf("Expected %d missing imports, got %d: %v",
			len(expectedMissing), len(report.MissingImports), report.MissingImports)
	}

	// Check that function usages were detected
	if len(report.UsedFunctions) == 0 {
		t.Error("No function usages detected")
	}

	// Check that current imports were detected
	if len(report.CurrentImports) != 1 || report.CurrentImports[0].Package != "fmt" {
		t.Errorf("Expected 1 current import (fmt), got %v", report.CurrentImports)
	}
}

func TestAnalyzeTemplateFileWithErrors(t *testing.T) {
	detector := NewImportDetector()

	// Test with non-existent file
	report, err := detector.AnalyzeTemplateFile("non-existent-file.go.tmpl")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if report == nil {
		t.Fatal("Report should not be nil even on error")
	}

	if len(report.Errors) == 0 {
		t.Error("Expected errors to be recorded in report")
	}
}

func TestBuildFunctionPackageMap(t *testing.T) {
	mapping := buildFunctionPackageMap()

	if len(mapping) == 0 {
		t.Fatal("Function package mapping is empty")
	}

	// Test some critical mappings
	criticalMappings := map[string]string{
		"time.Now":         "time",
		"fmt.Printf":       "fmt",
		"strings.Contains": "strings",
		"json.Marshal":     "encoding/json",
		"http.Get":         "net/http",
		"os.Getenv":        "os",
		"errors.New":       "errors",
	}

	for function, expectedPackage := range criticalMappings {
		if actualPackage, exists := mapping[function]; !exists {
			t.Errorf("Critical function %s not found in mapping", function)
		} else if actualPackage != expectedPackage {
			t.Errorf("Function %s mapped to %s, expected %s", function, actualPackage, expectedPackage)
		}
	}

	// Ensure no empty mappings
	for function, pkg := range mapping {
		if pkg == "" {
			t.Errorf("Function %s has empty package mapping", function)
		}
		if function == "" {
			t.Error("Found empty function name in mapping")
		}
	}
}
