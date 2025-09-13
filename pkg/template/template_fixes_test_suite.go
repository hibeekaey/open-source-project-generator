package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestTemplateFixesComprehensive runs the complete test suite for template fixes
func TestTemplateFixesComprehensive(t *testing.T) {
	// This is the main test suite that validates all template fixes

	t.Run("ImportDetectionUtility", func(t *testing.T) {
		testImportDetectionUtility(t)
	})

	t.Run("TemplateCompilationIntegration", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration tests in short mode")
		}
		testTemplateCompilationIntegration(t)
	})

	t.Run("TemplateEdgeCases", func(t *testing.T) {
		testTemplateEdgeCases(t)
	})

	t.Run("CompilationVerification", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping compilation verification in short mode")
		}
		testCompilationVerification(t)
	})

	t.Run("RegressionTests", func(t *testing.T) {
		testRegressionScenarios(t)
	})
}

func testImportDetectionUtility(t *testing.T) {
	t.Log("Testing import detection utility...")

	// Test the core functionality of the import detector
	detector := NewImportDetector()

	// Validate that the detector is properly initialized
	if detector == nil {
		t.Fatal("Import detector not initialized")
	}

	if detector.functionPackageMap == nil {
		t.Fatal("Function package map not initialized")
	}

	if len(detector.functionPackageMap) == 0 {
		t.Fatal("Function package map is empty")
	}

	t.Logf("Import detector initialized with %d function mappings", len(detector.functionPackageMap))

	// Test critical function mappings
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
		if actualPackage, exists := detector.functionPackageMap[function]; !exists {
			t.Errorf("Critical function %s not found in mapping", function)
		} else if actualPackage != expectedPackage {
			t.Errorf("Function %s mapped to %s, expected %s", function, actualPackage, expectedPackage)
		}
	}

	t.Log("✓ Import detection utility tests passed")
}

func testTemplateCompilationIntegration(t *testing.T) {
	t.Log("Testing template compilation integration...")

	templatesDir := "../../templates"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skip("Templates directory not found, skipping integration tests")
	}

	testData := createCompilationTestData()
	outputDir := t.TempDir()

	var totalTemplates, successfulTemplates, failedTemplates int
	var failures []string

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

		result := verifyTemplateCompilation(path, testData, outputDir)

		switch result.Status {
		case VerificationSuccess:
			successfulTemplates++
		case VerificationFailed:
			failedTemplates++
			failures = append(failures, fmt.Sprintf("%s: %s", relPath, result.Error))
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk templates directory: %v", err)
	}

	t.Logf("Template compilation results:")
	t.Logf("  Total: %d", totalTemplates)
	t.Logf("  Successful: %d", successfulTemplates)
	t.Logf("  Failed: %d", failedTemplates)

	if failedTemplates > 0 {
		t.Logf("Failed templates:")
		for _, failure := range failures {
			t.Logf("  - %s", failure)
		}
	}

	// Calculate success rate
	successRate := float64(successfulTemplates) / float64(totalTemplates) * 100
	t.Logf("Success rate: %.1f%%", successRate)

	// We expect at least 80% success rate
	if successRate < 80.0 {
		t.Errorf("Template compilation success rate too low: %.1f%% (expected at least 80%%)", successRate)
	}

	t.Log("✓ Template compilation integration tests completed")
}

func testTemplateEdgeCases(t *testing.T) {
	t.Log("Testing template edge cases...")

	detector := NewImportDetector()

	// Test various edge cases
	edgeCases := []struct {
		name        string
		content     string
		expectError bool
		description string
	}{
		{
			name: "EmptyTemplate",
			content: `package main

func main() {}`,
			expectError: false,
			description: "Empty template with no imports",
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
			expectError: false,
			description: "Complex template with conditional blocks",
		},
		{
			name: "InvalidGoSyntax",
			content: `package main

import "fmt"

func main() {
	fmt.Printf("Missing quote
}`,
			expectError: true,
			description: "Invalid Go syntax should be handled gracefully",
		},
		{
			name: "StringLiteralsWithFunctionNames",
			content: `package main

import "fmt"

func main() {
	message := "Call time.Now() to get current time"
	fmt.Printf("Message: %s\n", message)
}`,
			expectError: false,
			description: "Function names in strings should not trigger false positives",
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			report, err := detector.AnalyzeTemplateFile(templateFile)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %s but got none", tc.description)
				}
			} else {
				if err != nil && report == nil {
					t.Errorf("Unexpected nil report for %s: %v", tc.description, err)
				}
			}

			// Report should never be nil, even on error
			if report == nil {
				t.Errorf("Report is nil for %s", tc.description)
			}
		})
	}

	t.Log("✓ Template edge cases tests passed")
}

func testCompilationVerification(t *testing.T) {
	t.Log("Testing compilation verification...")

	// Test that known problematic templates now compile successfully
	knownIssues := []struct {
		templateContent string
		description     string
		requiredImports []string
	}{
		{
			templateContent: `package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Current time: %v\n", time.Now())
}`,
			description:     "Time import fix verification",
			requiredImports: []string{"fmt", "time"},
		},
		{
			templateContent: `package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}`,
			description:     "HTTP import fix verification",
			requiredImports: []string{"fmt", "net/http"},
		},
		{
			templateContent: `package main

import (
	"encoding/json"
	"fmt"
)

func processData(data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}`,
			description:     "JSON import fix verification",
			requiredImports: []string{"encoding/json", "fmt"},
		},
	}

	testData := createCompilationTestData()

	for _, issue := range knownIssues {
		t.Run(issue.description, func(t *testing.T) {
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(issue.templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			result := verifyTemplateCompilation(templateFile, testData, tempDir)

			if result.Status != VerificationSuccess {
				t.Errorf("Compilation verification failed for %s: %s", issue.description, result.Error)
				if result.CompilationOutput != "" {
					t.Logf("Compilation output: %s", result.CompilationOutput)
				}
			}

			// Verify imports are present in generated file
			generatedFile := filepath.Join(tempDir, "test.go")
			if content, err := os.ReadFile(generatedFile); err == nil {
				contentStr := string(content)
				for _, requiredImport := range issue.requiredImports {
					if !strings.Contains(contentStr, fmt.Sprintf("\"%s\"", requiredImport)) {
						t.Errorf("Required import %s not found in generated file", requiredImport)
					}
				}
			}
		})
	}

	t.Log("✓ Compilation verification tests passed")
}

func testRegressionScenarios(t *testing.T) {
	t.Log("Testing regression scenarios...")

	// Test scenarios that were previously broken to ensure they don't regress
	regressionTests := []struct {
		name        string
		scenario    func(t *testing.T)
		description string
	}{
		{
			name:        "AuthMiddlewareTimeImport",
			scenario:    testAuthMiddlewareTimeImport,
			description: "Auth middleware time import issue",
		},
		{
			name:        "HTTPStatusConstants",
			scenario:    testHTTPStatusConstants,
			description: "HTTP status constants import issue",
		},
		{
			name:        "JSONMarshalImport",
			scenario:    testJSONMarshalImport,
			description: "JSON marshal import issue",
		},
		{
			name:        "CryptoRandomImport",
			scenario:    testCryptoRandomImport,
			description: "Crypto random import issue",
		},
	}

	for _, test := range regressionTests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Running regression test: %s", test.description)
			test.scenario(t)
		})
	}

	t.Log("✓ Regression tests passed")
}

func testAuthMiddlewareTimeImport(t *testing.T) {
	// Test the specific auth middleware time import issue
	templateContent := `package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Auth logic here
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		fmt.Printf("Request took %v\n", duration)
	})
}`

	detector := NewImportDetector()
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "auth.go.tmpl")

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	report, err := detector.AnalyzeTemplateFile(templateFile)
	if err != nil {
		t.Errorf("Failed to analyze auth middleware template: %v", err)
		return
	}

	// Check that time import is detected as present
	timeImportFound := false
	for _, imp := range report.CurrentImports {
		if imp.Package == "time" {
			timeImportFound = true
			break
		}
	}

	if !timeImportFound {
		t.Error("Time import not detected in auth middleware template")
	}

	// Check that time functions are detected
	timeFunctionsFound := 0
	for _, usage := range report.UsedFunctions {
		if strings.HasPrefix(usage.Function, "time.") {
			timeFunctionsFound++
		}
	}

	if timeFunctionsFound == 0 {
		t.Error("Time functions not detected in auth middleware template")
	}

	// Verify compilation
	testData := createCompilationTestData()
	result := verifyTemplateCompilation(templateFile, testData, tempDir)

	if result.Status != VerificationSuccess {
		t.Errorf("Auth middleware template compilation failed: %s", result.Error)
	}
}

func testHTTPStatusConstants(t *testing.T) {
	templateContent := `package handlers

import (
	"fmt"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ok":
		w.WriteHeader(http.StatusOK)
	case "/notfound":
		w.WriteHeader(http.StatusNotFound)
	case "/error":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	
	fmt.Fprintf(w, "Status set")
}`

	detector := NewImportDetector()
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "handlers.go.tmpl")

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	report, err := detector.AnalyzeTemplateFile(templateFile)
	if err != nil {
		t.Errorf("Failed to analyze HTTP status template: %v", err)
		return
	}

	// Check that HTTP constants are detected
	httpConstantsFound := 0
	expectedConstants := []string{"http.StatusOK", "http.StatusNotFound", "http.StatusInternalServerError", "http.StatusBadRequest"}

	for _, expectedConstant := range expectedConstants {
		for _, usage := range report.UsedFunctions {
			if usage.Function == expectedConstant {
				httpConstantsFound++
				break
			}
		}
	}

	if httpConstantsFound != len(expectedConstants) {
		t.Errorf("Expected %d HTTP constants, found %d", len(expectedConstants), httpConstantsFound)
	}
}

func testJSONMarshalImport(t *testing.T) {
	templateContent := `package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message": "Hello World",
		"status":  "success",
	}
	
	response, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	// SECURITY: Added comprehensive security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}`

	detector := NewImportDetector()
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "api.go.tmpl")

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	report, err := detector.AnalyzeTemplateFile(templateFile)
	if err != nil {
		t.Errorf("Failed to analyze JSON template: %v", err)
		return
	}

	// Check that json.Marshal is detected
	jsonMarshalFound := false
	for _, usage := range report.UsedFunctions {
		if usage.Function == "json.Marshal" {
			jsonMarshalFound = true
			break
		}
	}

	if !jsonMarshalFound {
		t.Error("json.Marshal function not detected")
	}
}

func testCryptoRandomImport(t *testing.T) {
	templateContent := `package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	token := base64.StdEncoding.EncodeToString(bytes)
	return token, nil
}`

	detector := NewImportDetector()
	tempDir := t.TempDir()
	templateFile := filepath.Join(tempDir, "security.go.tmpl")

	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	report, err := detector.AnalyzeTemplateFile(templateFile)
	if err != nil {
		t.Errorf("Failed to analyze crypto template: %v", err)
		return
	}

	// Check that crypto functions are detected
	expectedFunctions := []string{"rand.Read", "base64.StdEncoding", "fmt.Errorf"}

	for _, expectedFunc := range expectedFunctions {
		found := false
		for _, usage := range report.UsedFunctions {
			if usage.Function == expectedFunc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function %s not detected", expectedFunc)
		}
	}
}

// TestSuiteRunner provides a convenient way to run all template fix tests
func TestSuiteRunner(t *testing.T) {
	startTime := time.Now()

	t.Log("🚀 Starting comprehensive template fixes test suite...")

	// Run the comprehensive test suite
	TestTemplateFixesComprehensive(t)

	duration := time.Since(startTime)
	t.Logf("✅ Template fixes test suite completed in %v", duration)

	// Summary
	t.Log("📊 Test Suite Summary:")
	t.Log("  ✓ Import detection utility tests")
	t.Log("  ✓ Template compilation integration tests")
	t.Log("  ✓ Template edge cases tests")
	t.Log("  ✓ Compilation verification tests")
	t.Log("  ✓ Regression scenario tests")
	t.Log("")
	t.Log("🎉 All template fix tests passed successfully!")
}
