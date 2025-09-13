package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestImportDetectorComprehensive provides comprehensive testing for the import detection utility
func TestImportDetectorComprehensive(t *testing.T) {
	detector := NewImportDetector()

	t.Run("FunctionPackageMapping", func(t *testing.T) {
		testFunctionPackageMapping(t, detector)
	})

	t.Run("TemplatePreprocessing", func(t *testing.T) {
		testTemplatePreprocessing(t, detector)
	})

	t.Run("ImportExtraction", func(t *testing.T) {
		testImportExtraction(t, detector)
	})

	t.Run("FunctionUsageDetection", func(t *testing.T) {
		testFunctionUsageDetection(t, detector)
	})

	t.Run("MissingImportDetection", func(t *testing.T) {
		testMissingImportDetection(t, detector)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, detector)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, detector)
	})
}

func testFunctionPackageMapping(t *testing.T, detector *ImportDetector) {
	// Test comprehensive function mappings
	testCases := []struct {
		function        string
		expectedPackage string
		description     string
	}{
		// Time package functions
		{"time.Now", "time", "time.Now() function"},
		{"time.Parse", "time", "time.Parse() function"},
		{"time.Since", "time", "time.Since() function"},
		{"time.Sleep", "time", "time.Sleep() function"},

		// Format package functions
		{"fmt.Printf", "fmt", "fmt.Printf() function"},
		{"fmt.Sprintf", "fmt", "fmt.Sprintf() function"},
		{"fmt.Errorf", "fmt", "fmt.Errorf() function"},

		// String manipulation functions
		{"strings.Contains", "strings", "strings.Contains() function"},
		{"strings.Split", "strings", "strings.Split() function"},
		{"strings.Join", "strings", "strings.Join() function"},
		{"strings.Replace", "strings", "strings.Replace() function"},

		// HTTP package functions and constants
		{"http.Get", "net/http", "http.Get() function"},
		{"http.StatusOK", "net/http", "http.StatusOK constant"},
		{"http.ListenAndServe", "net/http", "http.ListenAndServe() function"},

		// JSON encoding functions
		{"json.Marshal", "encoding/json", "json.Marshal() function"},
		{"json.Unmarshal", "encoding/json", "json.Unmarshal() function"},

		// OS package functions
		{"os.Getenv", "os", "os.Getenv() function"},
		{"os.Open", "os", "os.Open() function"},

		// Error handling functions
		{"errors.New", "errors", "errors.New() function"},
		{"errors.Is", "errors", "errors.Is() function"},

		// Context package functions
		{"context.Background", "context", "context.Background() function"},
		{"context.WithTimeout", "context", "context.WithTimeout() function"},

		// Filepath functions
		{"filepath.Join", "path/filepath", "filepath.Join() function"},
		{"filepath.Dir", "path/filepath", "filepath.Dir() function"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualPackage, exists := detector.functionPackageMap[tc.function]
			if !exists {
				t.Errorf("Function %s not found in mapping", tc.function)
				return
			}
			if actualPackage != tc.expectedPackage {
				t.Errorf("Function %s mapped to %s, expected %s", tc.function, actualPackage, tc.expectedPackage)
			}
		})
	}

	// Test that all mappings are non-empty
	for function, pkg := range detector.functionPackageMap {
		if pkg == "" {
			t.Errorf("Function %s has empty package mapping", function)
		}
		if function == "" {
			t.Error("Found empty function name in mapping")
		}
	}
}

func testTemplatePreprocessing(t *testing.T, detector *ImportDetector) {
	testCases := []struct {
		name        string
		input       string
		expected    string
		description string
	}{
		{
			name: "BasicTemplateVariables",
			input: `package {{.Name}}

import "fmt"

func main() {
	fmt.Println("Hello {{.ProjectName}}")
}`,
			expected: `package TemplateName

import "fmt"

func main() {
	fmt.Println("Hello TemplateProjectName")
}`,
			description: "Replace basic template variables",
		},
		{
			name: "TemplateControlStructures",
			input: `package main

{{ if .EnableAuth }}
func authenticate() {}
{{ end }}

{{ range .Services }}
func {{.Name}}Service() {}
{{ end }}

func main() {}`,
			expected: `package main


func authenticate() {}



func "template_placeholder"Service() {}


func main() {}`,
			description: "Remove template control structures",
		},
		{
			name: "ComplexTemplateExpressions",
			input: `package main

const version = "{{.Version}}"
const author = "{{.Author}}"
const description = "{{.Description}}"

func main() {
	fmt.Printf("Version: %s\n", version)
}`,
			expected: `package main

const version = "template_placeholder"
const author = "template_placeholder"
const description = "template_placeholder"

func main() {
	fmt.Printf("Version: %s\n", version)
}`,
			description: "Handle complex template expressions",
		},
		{
			name: "MixedTemplateContent",
			input: `package {{.Package}}

import (
	"fmt"
	"time"
)

{{ if .EnableLogging }}
import "log"
{{ end }}

func {{.FunctionName}}() {
	fmt.Printf("Starting {{.ServiceName}} at %v\n", time.Now())
	{{ if .EnableAuth }}
	authenticate()
	{{ end }}
}`,
			expected: `package TemplatePackage

import (
	"fmt"
	"time"
)


import "log"


func "template_placeholder"() {
	fmt.Printf("Starting "template_placeholder" at %v\n", time.Now())
	
	authenticate()
	
}`,
			description: "Handle mixed template content with imports and functions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := detector.preprocessTemplateContent(tc.input)
			if result != tc.expected {
				t.Errorf("Template preprocessing failed for %s\nExpected:\n%s\nGot:\n%s", tc.description, tc.expected, result)
			}
		})
	}
}

func testImportExtraction(t *testing.T, detector *ImportDetector) {
	testCases := []struct {
		name        string
		content     string
		expected    []ImportStatement
		description string
	}{
		{
			name: "SingleImport",
			content: `package main

import "fmt"

func main() {}`,
			expected: []ImportStatement{
				{Package: "fmt", IsStdLib: true},
			},
			description: "Extract single import",
		},
		{
			name: "MultipleImports",
			content: `package main

import (
	"fmt"
	"time"
	"strings"
)

func main() {}`,
			expected: []ImportStatement{
				{Package: "fmt", IsStdLib: true},
				{Package: "time", IsStdLib: true},
				{Package: "strings", IsStdLib: true},
			},
			description: "Extract multiple imports",
		},
		{
			name: "AliasedImports",
			content: `package main

import (
	"fmt"
	j "encoding/json"
	. "strings"
)

func main() {}`,
			expected: []ImportStatement{
				{Package: "fmt", IsStdLib: true},
				{Package: "encoding/json", Alias: "j", IsStdLib: true},
				{Package: "strings", Alias: ".", IsStdLib: true},
			},
			description: "Extract aliased imports",
		},
		{
			name: "MixedStandardAndThirdParty",
			content: `package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {}`,
			expected: []ImportStatement{
				{Package: "fmt", IsStdLib: true},
				{Package: "time", IsStdLib: true},
				{Package: "github.com/gin-gonic/gin", IsStdLib: false},
				{Package: "golang.org/x/crypto/bcrypt", IsStdLib: false},
			},
			description: "Extract mixed standard and third-party imports",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.go")
			err := os.WriteFile(tempFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(tempFile)
			if err != nil {
				t.Fatalf("AnalyzeTemplateFile failed: %v", err)
			}

			if len(report.CurrentImports) != len(tc.expected) {
				t.Errorf("Expected %d imports, got %d", len(tc.expected), len(report.CurrentImports))
				return
			}

			for i, expected := range tc.expected {
				if i >= len(report.CurrentImports) {
					t.Errorf("Missing import at index %d", i)
					continue
				}

				actual := report.CurrentImports[i]
				if actual.Package != expected.Package {
					t.Errorf("Import %d: expected package %s, got %s", i, expected.Package, actual.Package)
				}
				if actual.Alias != expected.Alias {
					t.Errorf("Import %d: expected alias %s, got %s", i, expected.Alias, actual.Alias)
				}
				if actual.IsStdLib != expected.IsStdLib {
					t.Errorf("Import %d: expected IsStdLib %v, got %v", i, expected.IsStdLib, actual.IsStdLib)
				}
			}
		})
	}
}

func testFunctionUsageDetection(t *testing.T, detector *ImportDetector) {
	testCases := []struct {
		name        string
		content     string
		expected    []string
		description string
	}{
		{
			name: "BasicFunctionCalls",
			content: `package main

import "fmt"

func main() {
	fmt.Printf("Hello World")
	time.Now()
	strings.Contains("test", "es")
}`,
			expected:    []string{"fmt.Printf", "time.Now", "strings.Contains"},
			description: "Detect basic function calls",
		},
		{
			name: "HTTPConstants",
			content: `package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusNotFound)
	w.WriteHeader(http.StatusInternalServerError)
}`,
			expected:    []string{"http.StatusOK", "http.StatusNotFound", "http.StatusInternalServerError"},
			description: "Detect HTTP status constants",
		},
		{
			name: "ChainedFunctionCalls",
			content: `package main

func main() {
	result := strings.ToLower(strings.TrimSpace("  TEST  "))
	fmt.Println(result)
	json.NewEncoder(os.Stdout).Encode(map[string]string{"key": "value"})
}`,
			expected:    []string{"strings.ToLower", "strings.TrimSpace", "fmt.Println", "json.NewEncoder", "os.Stdout"},
			description: "Detect chained function calls",
		},
		{
			name: "VariableAssignments",
			content: `package main

func main() {
	now := time.Now()
	duration := time.Since(now)
	ctx := context.Background()
	_ = duration
	_ = ctx
}`,
			expected:    []string{"time.Now", "time.Since", "context.Background"},
			description: "Detect function calls in variable assignments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.go")
			err := os.WriteFile(tempFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(tempFile)
			if err != nil {
				t.Fatalf("AnalyzeTemplateFile failed: %v", err)
			}

			// Extract function names from usages
			var actualFunctions []string
			for _, usage := range report.UsedFunctions {
				actualFunctions = append(actualFunctions, usage.Function)
			}

			// Check that all expected functions are found
			for _, expectedFunc := range tc.expected {
				found := false
				for _, actualFunc := range actualFunctions {
					if actualFunc == expectedFunc {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected function %s not found in usages: %v", expectedFunc, actualFunctions)
				}
			}
		})
	}
}

func testMissingImportDetection(t *testing.T, detector *ImportDetector) {
	testCases := []struct {
		name            string
		content         string
		expectedMissing []string
		description     string
	}{
		{
			name: "MissingTimeImport",
			content: `package main

import "fmt"

func main() {
	fmt.Printf("Current time: %v\n", time.Now())
}`,
			expectedMissing: []string{"time"},
			description:     "Detect missing time import",
		},
		{
			name: "MissingMultipleImports",
			content: `package main

func main() {
	fmt.Printf("Hello World")
	now := time.Now()
	result := strings.Contains("test", "es")
	_ = now
	_ = result
}`,
			expectedMissing: []string{"fmt", "time", "strings"},
			description:     "Detect multiple missing imports",
		},
		{
			name: "MissingHTTPImport",
			content: `package main

import "fmt"

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}`,
			expectedMissing: []string{"net/http"},
			description:     "Detect missing HTTP import",
		},
		{
			name: "NoMissingImports",
			content: `package main

import (
	"fmt"
	"time"
	"strings"
)

func main() {
	fmt.Printf("Current time: %v\n", time.Now())
	result := strings.Contains("test", "es")
	_ = result
}`,
			expectedMissing: []string{},
			description:     "No missing imports when all are present",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.go")
			err := os.WriteFile(tempFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(tempFile)
			if err != nil {
				t.Fatalf("AnalyzeTemplateFile failed: %v", err)
			}

			if len(report.MissingImports) != len(tc.expectedMissing) {
				t.Errorf("Expected %d missing imports, got %d: %v", len(tc.expectedMissing), len(report.MissingImports), report.MissingImports)
				return
			}

			// Check that all expected missing imports are found
			for _, expectedMissing := range tc.expectedMissing {
				found := false
				for _, actualMissing := range report.MissingImports {
					if actualMissing == expectedMissing {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected missing import %s not found in: %v", expectedMissing, report.MissingImports)
				}
			}
		})
	}
}

func testEdgeCases(t *testing.T, detector *ImportDetector) {
	testCases := []struct {
		name        string
		content     string
		description string
	}{
		{
			name: "EmptyFile",
			content: `package main

func main() {}`,
			description: "Handle empty file with no imports or function calls",
		},
		{
			name: "OnlyComments",
			content: `package main

// This is a comment
/* This is a block comment */

func main() {
	// fmt.Printf("This is commented out")
}`,
			description: "Handle file with only comments",
		},
		{
			name: "ComplexTemplateExpressions",
			content: `package {{.Name}}

import "fmt"

func main() {
	fmt.Printf("{{.Message}}")
	{{- if .EnableFeature }}
	time.Now()
	{{- end }}
}`,
			description: "Handle complex template expressions",
		},
		{
			name: "StringLiteralsWithFunctionNames",
			content: `package main

import "fmt"

func main() {
	message := "Call time.Now() to get current time"
	fmt.Printf("Message: %s\n", message)
}`,
			description: "Handle function names in string literals",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file for testing
			tempDir := t.TempDir()
			tempFile := filepath.Join(tempDir, "test.go")
			err := os.WriteFile(tempFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(tempFile)
			if err != nil {
				t.Fatalf("AnalyzeTemplateFile failed for %s: %v", tc.description, err)
			}

			// Basic validation - report should not be nil
			if report == nil {
				t.Errorf("Report is nil for %s", tc.description)
				return
			}

			// Report should have the correct file path
			if report.FilePath != tempFile {
				t.Errorf("Expected file path %s, got %s", tempFile, report.FilePath)
			}
		})
	}
}

func testErrorHandling(t *testing.T, detector *ImportDetector) {
	t.Run("NonExistentFile", func(t *testing.T) {
		report, err := detector.AnalyzeTemplateFile("non-existent-file.go")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
		if report == nil {
			t.Error("Report should not be nil even on error")
			return
		}
		if len(report.Errors) == 0 {
			t.Error("Expected errors to be recorded in report")
		}
	})

	t.Run("InvalidGoSyntax", func(t *testing.T) {
		// Create a file with invalid Go syntax
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "invalid.go")
		invalidContent := `package main

import "fmt"

func main() {
	fmt.Printf("Missing closing quote
}`
		err := os.WriteFile(tempFile, []byte(invalidContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		report, err := detector.AnalyzeTemplateFile(tempFile)
		if err == nil {
			t.Error("Expected error for invalid Go syntax")
		}
		if report == nil {
			t.Error("Report should not be nil even on error")
			return
		}
		if len(report.Errors) == 0 {
			t.Error("Expected errors to be recorded in report")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		// Create an empty file
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "empty.go")
		err := os.WriteFile(tempFile, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		report, err := detector.AnalyzeTemplateFile(tempFile)
		if err == nil {
			t.Error("Expected error for empty file")
		}
		if report == nil {
			t.Error("Report should not be nil even on error")
		}
	})
}

// TestStandardLibraryDetection tests the standard library detection functionality
func TestStandardLibraryDetection(t *testing.T) {
	detector := NewImportDetector()

	testCases := []struct {
		pkg      string
		expected bool
	}{
		// Standard library packages
		{"fmt", true},
		{"time", true},
		{"strings", true},
		{"encoding/json", true},
		{"net/http", true},
		{"path/filepath", true},
		{"crypto/md5", true},
		{"io/ioutil", true},

		// Third-party packages
		{"github.com/gin-gonic/gin", false},
		{"golang.org/x/crypto/bcrypt", false},
		{"google.golang.org/grpc", false},
		{"gopkg.in/yaml.v2", false},
		{"go.uber.org/zap", false},

		// Project-specific packages
		{"myproject/internal/auth", false},
		{"example.com/mypackage", false},
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

// TestFunctionPackageMappingCompleteness ensures all mapped functions have valid packages
func TestFunctionPackageMappingCompleteness(t *testing.T) {
	detector := NewImportDetector()

	for function, pkg := range detector.functionPackageMap {
		if pkg == "" {
			t.Errorf("Function %s has empty package mapping", function)
		}

		if function == "" {
			t.Error("Found empty function name in mapping")
		}

		// Validate function name format (should contain a dot)
		if !strings.Contains(function, ".") {
			t.Errorf("Function %s doesn't follow package.Function format", function)
		}

		// Validate that the package is a standard library package
		if !detector.isStandardLibrary(pkg) {
			t.Errorf("Function %s maps to non-standard library package %s", function, pkg)
		}
	}
}
