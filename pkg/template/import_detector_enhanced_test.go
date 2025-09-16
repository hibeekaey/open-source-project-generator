package template

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportDetector_AnalyzeDirectory(t *testing.T) {
	// Create temporary directory with test templates
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test template files
	testFiles := map[string]string{
		"good.go.tmpl": `package main
import (
	"fmt"
	"time"
)
func main() {
	fmt.Println(time.Now())
}`,
		"bad.go.tmpl": `package main
import (
	"fmt"
)
func main() {
	fmt.Println(time.Now())
	strings.ToLower("test")
}`,
		"non_template.go": `package main
func main() {
	// This should be ignored
}`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	detector := NewImportDetector()
	analysis, err := detector.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze directory: %v", err)
	}

	// Should only analyze .tmpl files and only report files with issues
	if analysis.Summary.TotalFiles != 1 {
		t.Errorf("Expected 1 file with issues, got %d", analysis.Summary.TotalFiles)
	}

	if analysis.Summary.FilesWithIssues != 1 {
		t.Errorf("Expected 1 file with issues, got %d", analysis.Summary.FilesWithIssues)
	}

	// Check that the bad template was found
	found := false
	for _, report := range analysis.Reports {
		if strings.Contains(report.FilePath, "bad.go.tmpl") {
			found = true
			if len(report.MissingImports) == 0 {
				t.Error("Expected missing imports in bad.go.tmpl")
			}
		}
	}

	if !found {
		t.Error("Expected to find bad.go.tmpl in analysis results")
	}
}

func TestImportDetector_GenerateTextReport(t *testing.T) {
	detector := NewImportDetector()

	// Create sample analysis data
	analysis := &AnalysisReport{
		Reports: []MissingImportReport{
			{
				FilePath:       "test.go.tmpl",
				MissingImports: []string{"time", "strings"},
				UsedFunctions: []FunctionUsage{
					{Function: "time.Now", Line: 5, RequiredPackage: "time"},
					{Function: "strings.ToLower", Line: 6, RequiredPackage: "strings"},
				},
				CurrentImports: []ImportStatement{
					{Package: "fmt", IsStdLib: true},
				},
			},
		},
		Summary: AnalysisSummary{
			TotalFiles:          1,
			FilesWithIssues:     1,
			TotalMissingImports: 2,
			MostCommonMissing:   map[string]int{"time": 1, "strings": 1},
		},
		GeneratedAt: "2023-01-01T00:00:00Z",
	}

	report := detector.GenerateTextReport(analysis)

	// Check that report contains expected sections
	expectedSections := []string{
		"Template Import Analysis Report",
		"Generated: 2023-01-01T00:00:00Z",
		"Total Files Analyzed: 1",
		"Files with Issues: 1",
		"Most Common Missing Imports:",
		"Detailed Results:",
		"test.go.tmpl",
		"Missing Imports:",
		"time",
		"strings",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing expected section: %s", section)
		}
	}
}

func TestImportDetector_JSONSerialization(t *testing.T) {
	// Create sample analysis data
	analysis := &AnalysisReport{
		Reports: []MissingImportReport{
			{
				FilePath:       "test.go.tmpl",
				MissingImports: []string{"time"},
				UsedFunctions: []FunctionUsage{
					{Function: "time.Now", Line: 5, Column: 10, RequiredPackage: "time"},
				},
				CurrentImports: []ImportStatement{
					{Package: "fmt", IsStdLib: true},
				},
			},
		},
		Summary: AnalysisSummary{
			TotalFiles:          1,
			FilesWithIssues:     1,
			TotalMissingImports: 1,
			MostCommonMissing:   map[string]int{"time": 1},
		},
		GeneratedAt: "2023-01-01T00:00:00Z",
	}

	// Test JSON serialization
	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Test JSON deserialization
	var deserializedAnalysis AnalysisReport
	err = json.Unmarshal(jsonData, &deserializedAnalysis)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify data integrity
	if deserializedAnalysis.Summary.TotalFiles != analysis.Summary.TotalFiles {
		t.Errorf("JSON serialization failed: TotalFiles mismatch")
	}

	if len(deserializedAnalysis.Reports) != len(analysis.Reports) {
		t.Errorf("JSON serialization failed: Reports count mismatch")
	}
}

func TestImportDetector_AddFunctionMapping(t *testing.T) {
	detector := NewImportDetector()

	// Add custom function mapping
	detector.AddFunctionMapping("custom.Function", "github.com/example/custom")

	// Test that custom mapping works
	templateContent := `package main

func main() {
	custom.Function()
}`

	// Create temp file for testing
	tempFile, err := os.CreateTemp("", "test*.go.tmpl")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()

	_, err = tempFile.WriteString(templateContent)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	_ = tempFile.Close()

	report, err := detector.AnalyzeTemplateFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to analyze template file: %v", err)
	}

	expectedMissing := "github.com/example/custom"
	found := false
	for _, missing := range report.MissingImports {
		if missing == expectedMissing {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected missing import '%s' not found in %v", expectedMissing, report.MissingImports)
	}
}

func TestImportDetector_GetFunctionMappings(t *testing.T) {
	detector := NewImportDetector()

	mappings := detector.GetFunctionMappings()

	// Check that we have expected standard library mappings
	expectedMappings := map[string]string{
		"time.Now":         "time",
		"fmt.Printf":       "fmt",
		"strings.Contains": "strings",
		"json.Marshal":     "encoding/json",
		"http.Get":         "net/http",
	}

	for function, expectedPackage := range expectedMappings {
		if pkg, exists := mappings[function]; !exists {
			t.Errorf("Function %s not found in mappings", function)
		} else if pkg != expectedPackage {
			t.Errorf("Function %s mapped to %s, expected %s", function, pkg, expectedPackage)
		}
	}

	// Test that modifying returned map doesn't affect original
	originalCount := len(mappings)
	mappings["test.Function"] = "test"

	newMappings := detector.GetFunctionMappings()
	if len(newMappings) != originalCount {
		t.Error("GetFunctionMappings should return a copy, not the original map")
	}
}

func TestImportDetector_EmptyDirectory(t *testing.T) {
	// Create empty temporary directory
	tempDir, err := os.MkdirTemp("", "empty_template_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	detector := NewImportDetector()
	analysis, err := detector.AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("Failed to analyze empty directory: %v", err)
	}

	if analysis.Summary.TotalFiles != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", analysis.Summary.TotalFiles)
	}

	if analysis.Summary.FilesWithIssues != 0 {
		t.Errorf("Expected 0 files with issues in empty directory, got %d", analysis.Summary.FilesWithIssues)
	}

	report := detector.GenerateTextReport(analysis)
	if !strings.Contains(report, "No issues found") {
		t.Error("Expected 'No issues found' message for empty directory")
	}
}
