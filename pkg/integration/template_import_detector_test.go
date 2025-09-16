package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
)

// TestImportDetectorIntegration tests the import detector with real template files
func TestImportDetectorIntegration(t *testing.T) {
	detector := template.NewImportDetector()

	if detector == nil {
		t.Fatal("NewImportDetector returned nil")
	}

	// Test that function mappings are accessible through public interface
	mappings := detector.GetFunctionMappings()
	if mappings == nil {
		t.Fatal("GetFunctionMappings returned nil")
	}

	// Test that some expected mappings exist
	expectedMappings := map[string]string{
		"time.Now":         "time",
		"fmt.Printf":       "fmt",
		"strings.Contains": "strings",
		"json.Marshal":     "encoding/json",
	}

	for function, expectedPackage := range expectedMappings {
		if actualPackage, exists := mappings[function]; !exists {
			t.Errorf("Expected function %s not found in mapping", function)
		} else if actualPackage != expectedPackage {
			t.Errorf("Function %s mapped to %s, expected %s", function, actualPackage, expectedPackage)
		}
	}
}

// TestImportDetectorAnalyzeTemplateFile tests analyzing a real template file
func TestImportDetectorAnalyzeTemplateFile(t *testing.T) {
	detector := template.NewImportDetector()

	// Create a temporary template file for testing
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "test.go.tmpl")

	templateContent := `package {{.Name}}

import (
	"fmt"
)

func main() {
	fmt.Println("Hello {{.ProjectName}}")
	time.Now() // Missing import for time
	json.Marshal(data) // Missing import for encoding/json
}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	// Analyze the template file
	report, err := detector.AnalyzeTemplateFile(templateFile)
	if err != nil {
		t.Fatalf("Failed to analyze template file: %v", err)
	}

	if report == nil {
		t.Fatal("Analysis report is nil")
	}

	if report.FilePath != templateFile {
		t.Errorf("Expected file path %s, got %s", templateFile, report.FilePath)
	}

	// Check that missing imports were detected
	expectedMissingImports := []string{"time", "encoding/json"}
	for _, expectedImport := range expectedMissingImports {
		found := false
		for _, missingImport := range report.MissingImports {
			if missingImport == expectedImport {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected missing import %s not found in report", expectedImport)
		}
	}
}

// TestImportDetectorAnalyzeDirectory tests analyzing a directory of template files
func TestImportDetectorAnalyzeDirectory(t *testing.T) {
	detector := template.NewImportDetector()

	// Create a temporary directory with multiple template files
	tempDir := t.TempDir()

	// Create first template file
	template1 := filepath.Join(tempDir, "template1.go.tmpl")
	content1 := `package {{.Name}}
import "fmt"
func main() {
	fmt.Println("test")
	time.Sleep(1) // Missing time import
}`

	err := os.WriteFile(template1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create template1: %v", err)
	}

	// Create second template file
	template2 := filepath.Join(tempDir, "template2.go.tmpl")
	content2 := `package {{.Name}}
func process() {
	json.Marshal(data) // Missing encoding/json import
	http.Get("url") // Missing net/http import
}`

	err = os.WriteFile(template2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create template2: %v", err)
	}

	// Analyze the directory
	analysis, err := detector.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze directory: %v", err)
	}

	if analysis == nil {
		t.Fatal("Analysis result is nil")
	}

	if len(analysis.Reports) != 2 {
		t.Errorf("Expected 2 reports, got %d", len(analysis.Reports))
	}

	// Check that the summary contains expected information
	if analysis.Summary.TotalFiles != 2 {
		t.Errorf("Expected 2 total files in summary, got %d", analysis.Summary.TotalFiles)
	}

	if analysis.Summary.TotalMissingImports == 0 {
		t.Error("Expected missing imports to be detected")
	}
}

// TestImportDetectorCustomMapping tests adding custom function mappings
func TestImportDetectorCustomMapping(t *testing.T) {
	detector := template.NewImportDetector()

	// Add a custom function mapping
	detector.AddFunctionMapping("custom.Function", "github.com/example/custom")

	// Verify the mapping was added
	mappings := detector.GetFunctionMappings()
	if pkg, exists := mappings["custom.Function"]; !exists {
		t.Error("Custom function mapping not found")
	} else if pkg != "github.com/example/custom" {
		t.Errorf("Expected package 'github.com/example/custom', got '%s'", pkg)
	}
}

// TestImportDetectorGenerateReport tests report generation
func TestImportDetectorGenerateReport(t *testing.T) {
	detector := template.NewImportDetector()

	// Create a temporary template file
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "test.go.tmpl")

	templateContent := `package {{.Name}}
func main() {
	fmt.Println("test") // Missing fmt import
}`

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	// Analyze the directory
	analysis, err := detector.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze directory: %v", err)
	}

	// Generate text report
	report := detector.GenerateTextReport(analysis)
	if report == "" {
		t.Error("Generated report is empty")
	}

	// Check that the report contains expected information
	if !strings.Contains(report, "Import Analysis Report") {
		t.Error("Report should contain title")
	}

	if !strings.Contains(report, "fmt") {
		t.Error("Report should mention missing fmt import")
	}
}
