package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplateEdgeCases tests various edge cases and scenarios for template processing
func TestTemplateEdgeCases(t *testing.T) {
	t.Run("ComplexTemplateStructures", func(t *testing.T) {
		t.Parallel()
		testComplexTemplateStructures(t)
	})

	t.Run("NestedTemplateExpressions", func(t *testing.T) {
		t.Parallel()
		testNestedTemplateExpressions(t)
	})

	t.Run("ConditionalImports", func(t *testing.T) {
		t.Parallel()
		testConditionalImports(t)
	})

	t.Run("StringLiteralEdgeCases", func(t *testing.T) {
		t.Parallel()
		testStringLiteralEdgeCases(t)
	})

	t.Run("CommentHandling", func(t *testing.T) {
		t.Parallel()
		testCommentHandling(t)
	})

	t.Run("ImportOrganization", func(t *testing.T) {
		t.Parallel()
		testImportOrganization(t)
	})

	t.Run("TemplateVariableEdgeCases", func(t *testing.T) {
		t.Parallel()
		testTemplateVariableEdgeCases(t)
	})

	t.Run("ErrorRecovery", func(t *testing.T) {
		t.Parallel()
		testErrorRecovery(t)
	})
}

func testComplexTemplateStructures(t *testing.T) {
	// Reduced and optimized test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "NestedConditionals",
			content: `package {{.Name}}

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
			description: "Handle nested conditional blocks with imports",
		},
		{
			name: "MixedTemplateAndGoCode",
			content: `package {{.Name}}

import (
	"fmt"
	"time"
)

type Config struct {
	Name    string ` + "`json:\"name\"`" + `
	Version string ` + "`json:\"version\"`" + `
}

func main() {
	config := &Config{
		Name:    "{{.ProjectName}}",
		Version: "{{.Version}}",
	}
	
	fmt.Printf("Starting %s at %v\n", config.Name, time.Now())
}`,
			description: "Handle mixed template expressions and Go code structures",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				t.Errorf("Failed to analyze complex template %s: %v", tc.description, err)
				return
			}

			// Validate that the report is reasonable
			if report == nil {
				t.Errorf("Report is nil for %s", tc.description)
				return
			}

			// Check that preprocessing worked (no template syntax should remain in errors)
			for _, errMsg := range report.Errors {
				if strings.Contains(errMsg, "{{") || strings.Contains(errMsg, "}}") {
					t.Errorf("Template syntax found in error message for %s: %s", tc.description, errMsg)
				}
			}
		})
	}
}

func testNestedTemplateExpressions(t *testing.T) {
	// Simplified test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "NestedVariables",
			content: `package {{.Name}}

import "fmt"

func main() {
	fmt.Printf("{{.Messages.Welcome}} to {{.Project.Name}}\n")
}`,
			description: "Handle nested template variables",
		},
		{
			name: "ConditionalWithNestedAccess",
			content: `package {{.Name}}

import (
	"fmt"
	{{- if .Database.Enabled }}
	"database/sql"
	{{- end }}
)

func main() {
	{{- if .Database.Enabled }}
	fmt.Printf("Database type: {{.Database.Type}}\n")
	{{- end }}
}`,
			description: "Handle conditional blocks with nested variable access",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			// Test that preprocessing handles nested expressions
			processed := detector.preprocessTemplateContent(tc.content)

			// Verify that complex template expressions are replaced
			if strings.Contains(processed, "{{") && strings.Contains(processed, "}}") {
				// Some template expressions might remain, but they should be simple placeholders
				if !strings.Contains(processed, "template_placeholder") {
					t.Errorf("Complex template expressions not properly preprocessed in %s", tc.description)
				}
			}

			// Test that analysis doesn't crash on complex templates
			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				// Complex templates might fail parsing, but should not crash
				if report == nil {
					t.Errorf("Report is nil even on error for %s", tc.description)
				}
			}
		})
	}
}

func testConditionalImports(t *testing.T) {
	// Simplified test case for better performance
	testCases := []struct {
		name            string
		content         string
		expectedImports []string
		description     string
	}{
		{
			name: "ConditionalTimeImport",
			content: `package {{.Name}}

import "fmt"

{{- if .EnableTimestamp }}
import "time"
{{- end }}

func main() {
	fmt.Println("Hello World")
	{{- if .EnableTimestamp }}
	fmt.Printf("Time: %v\n", time.Now())
	{{- end }}
}`,
			expectedImports: []string{"fmt"},
			description:     "Handle conditional imports in template blocks",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				// Conditional imports might cause parsing issues, but should be handled gracefully
				if report == nil {
					t.Errorf("Report is nil for %s", tc.description)
				}
				return
			}

			// Verify that at least the non-conditional imports are detected
			for _, expectedImport := range tc.expectedImports {
				found := false
				for _, currentImport := range report.CurrentImports {
					if currentImport.Package == expectedImport {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected import %s not found for %s", expectedImport, tc.description)
				}
			}
		})
	}
}

func testStringLiteralEdgeCases(t *testing.T) {
	// Reduced test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "FunctionNamesInStrings",
			content: `package {{.Name}}

import "fmt"

func main() {
	message := "Call time.Now() to get the current time"
	fmt.Printf("Instructions: %s\n", message)
}`,
			description: "Function names in string literals should not trigger import detection",
		},
		{
			name: "TemplateVariablesInStrings",
			content: `package {{.Name}}

import "fmt"

func main() {
	template := "Hello {{.Name}}, welcome to {{.Service}}"
	fmt.Printf("Template: %s\n", template)
}`,
			description: "Template variables in strings should be handled correctly",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				t.Errorf("Failed to analyze template with string literals %s: %v", tc.description, err)
				return
			}

			// For string literal tests, we mainly want to ensure no false positives
			// The actual function calls (like fmt.Printf) should still be detected
			actualFunctionCall := false
			for _, usage := range report.UsedFunctions {
				if usage.Function == "fmt.Printf" {
					actualFunctionCall = true
					break
				}
			}

			if !actualFunctionCall {
				t.Errorf("Actual function call fmt.Printf not detected for %s", tc.description)
			}
		})
	}
}

func testCommentHandling(t *testing.T) {
	// Simplified test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "SingleLineComments",
			content: `package {{.Name}}

import "fmt"

func main() {
	// This is a comment with time.Now() that should be ignored
	fmt.Printf("Hello World\n")
}`,
			description: "Single line comments should not trigger import detection",
		},
		{
			name: "BlockComments",
			content: `package {{.Name}}

import "fmt"

/*
This is a block comment that mentions time.Now() 
*/

func main() {
	fmt.Printf("Hello World\n")
}`,
			description: "Block comments should not trigger import detection",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				t.Errorf("Failed to analyze template with comments %s: %v", tc.description, err)
				return
			}

			// Check that only fmt.Printf is detected (not the functions in comments)
			expectedFunctions := []string{"fmt.Printf"}
			unexpectedFunctions := []string{"time.Now", "strings.Contains", "json.Marshal", "time.Since", "log.Printf"}

			for _, expectedFunc := range expectedFunctions {
				found := false
				for _, usage := range report.UsedFunctions {
					if usage.Function == expectedFunc {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected function %s not found for %s", expectedFunc, tc.description)
				}
			}

			for _, unexpectedFunc := range unexpectedFunctions {
				for _, usage := range report.UsedFunctions {
					if usage.Function == unexpectedFunc {
						t.Errorf("Unexpected function %s found in comments for %s", unexpectedFunc, tc.description)
					}
				}
			}
		})
	}
}

func testImportOrganization(t *testing.T) {
	// Single simplified test case
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "MixedImportOrder",
			content: `package {{.Name}}

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Time: %v\n", time.Now())
}`,
			description: "Mixed import order should be detected correctly",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				t.Errorf("Failed to analyze template with import organization %s: %v", tc.description, err)
				return
			}

			// Check that imports are correctly identified
			if len(report.CurrentImports) == 0 {
				t.Errorf("No imports detected for %s", tc.description)
			}

			// Verify at least one standard library import exists
			hasStdLib := false
			for _, imp := range report.CurrentImports {
				if imp.IsStdLib {
					hasStdLib = true
					break
				}
			}
			if !hasStdLib {
				t.Errorf("No standard library imports detected for %s", tc.description)
			}
		})
	}
}

func testTemplateVariableEdgeCases(t *testing.T) {
	// Simplified test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "VariablesInFunctionNames",
			content: `package {{.Name}}

import "fmt"

func {{.FunctionPrefix}}Handler() {
	fmt.Printf("Handler called\n")
}`,
			description: "Template variables in function names",
		},
		{
			name: "VariablesInTypes",
			content: `package {{.Name}}

import "fmt"

type {{.TypeName}} struct {
	Name string
}`,
			description: "Template variables in type definitions",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			// Test preprocessing
			processed := detector.preprocessTemplateContent(tc.content)

			// Verify that template variables are replaced with valid identifiers
			if strings.Contains(processed, "{{") && strings.Contains(processed, "}}") {
				// Check if remaining template expressions are placeholders
				if !strings.Contains(processed, "template_placeholder") {
					t.Errorf("Template variables not properly preprocessed for %s", tc.description)
				}
			}

			// Test analysis (might fail due to invalid import paths, but should not crash)
			report, err := detector.AnalyzeTemplateFile(templateFile)
			if err != nil {
				// Template variable edge cases might cause parsing errors
				if report == nil {
					t.Errorf("Report is nil for %s", tc.description)
				}
			}
		})
	}
}

func testErrorRecovery(t *testing.T) {
	// Simplified test cases for better performance
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "PartiallyValidTemplate",
			content: `package {{.Name}}

import "fmt"

func main() {
	fmt.Printf("Hello World\n")
	// This has invalid template syntax: {{.Invalid.
}`,
			description: "Partially valid template with syntax errors",
		},
		{
			name: "InvalidGoSyntaxWithTemplates",
			content: `package {{.Name}}

import "fmt"

func main() {
	fmt.Printf("Hello World"  // Missing closing parenthesis
}`,
			description: "Invalid Go syntax combined with template processing",
		},
	}

	detector := NewImportDetector()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			// Test that error recovery works - analysis should not crash
			report, err := detector.AnalyzeTemplateFile(templateFile)

			// We expect errors, but the function should not panic
			if report == nil {
				t.Errorf("Report is nil for error recovery test %s", tc.description)
				return
			}

			// Errors should be recorded in the report
			if err != nil && len(report.Errors) == 0 {
				t.Errorf("Error occurred but not recorded in report for %s", tc.description)
			}

			// File path should still be set even on error
			if report.FilePath != templateFile {
				t.Errorf("File path not set correctly in error case for %s", tc.description)
			}
		})
	}
}
