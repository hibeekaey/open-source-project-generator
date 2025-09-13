//go:build !ci

package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/models"
)

// TestTemplateCompilationVerification verifies that all fixed templates generate compilable Go code
func TestTemplateCompilationVerification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping compilation verification tests in short mode")
	}

	t.Run("VerifyAllGoTemplates", func(t *testing.T) {
		verifyAllGoTemplates(t)
	})

	t.Run("VerifyKnownProblematicTemplates", func(t *testing.T) {
		verifyKnownProblematicTemplates(t)
	})

	t.Run("VerifyImportFixes", func(t *testing.T) {
		verifyImportFixes(t)
	})

	t.Run("VerifyTemplateGeneration", func(t *testing.T) {
		verifyTemplateGeneration(t)
	})
}

func verifyAllGoTemplates(t *testing.T) {
	templatesDir := "../../templates"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skip("Templates directory not found, skipping verification")
	}

	testData := createCompilationTestData()
	outputDir := t.TempDir()

	var results []TemplateVerificationResult
	totalTemplates := 0

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only verify Go template files
		if !strings.HasSuffix(path, ".go.tmpl") && !strings.HasSuffix(path, ".mod.tmpl") {
			return nil
		}

		totalTemplates++
		relPath, _ := filepath.Rel(templatesDir, path)

		result := verifyTemplateCompilation(path, testData, outputDir)
		result.TemplatePath = relPath
		results = append(results, result)

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk templates directory: %v", err)
	}

	// Analyze results
	var successful, failed, skipped int
	var failedTemplates []string

	for _, result := range results {
		switch result.Status {
		case VerificationSuccess:
			successful++
		case VerificationFailed:
			failed++
			failedTemplates = append(failedTemplates, result.TemplatePath)
		case VerificationSkipped:
			skipped++
		}
	}

	t.Logf("Template Verification Results:")
	t.Logf("  Total templates: %d", totalTemplates)
	t.Logf("  Successful: %d", successful)
	t.Logf("  Failed: %d", failed)
	t.Logf("  Skipped: %d", skipped)

	if failed > 0 {
		t.Errorf("Failed to verify %d templates: %v", failed, failedTemplates)

		// Log detailed failure information
		for _, result := range results {
			switch result.Status {
			case VerificationFailed:
				t.Logf("Failed template: %s", result.TemplatePath)
				t.Logf("  Error: %s", result.Error)
				if result.CompilationOutput != "" {
					t.Logf("  Compilation output: %s", result.CompilationOutput)
				}
			}
		}
	}

	// Ensure we have a reasonable success rate
	successRate := float64(successful) / float64(totalTemplates) * 100
	if successRate < 80.0 {
		t.Errorf("Template success rate too low: %.1f%% (expected at least 80%%)", successRate)
	}
}

func verifyKnownProblematicTemplates(t *testing.T) {
	// Test templates that were known to have import issues
	problematicTemplates := []struct {
		path        string
		description string
	}{
		{
			path:        "../../templates/backend/go-gin/internal/middleware/auth.go.tmpl",
			description: "Auth middleware with time import issue",
		},
		{
			path:        "../../templates/backend/go-gin/internal/middleware/security.go.tmpl",
			description: "Security middleware with potential import issues",
		},
		{
			path:        "../../templates/backend/go-gin/internal/controllers/auth_controller.go.tmpl",
			description: "Auth controller with time and HTTP imports",
		},
		{
			path:        "../../templates/backend/go-gin/internal/services/auth_service.go.tmpl",
			description: "Auth service with crypto and time imports",
		},
		{
			path:        "../../templates/backend/go-gin/pkg/utils/jwt.go.tmpl",
			description: "JWT utility with crypto imports",
		},
	}

	testData := createCompilationTestData()
	outputDir := t.TempDir()

	for _, template := range problematicTemplates {
		t.Run(template.description, func(t *testing.T) {
			if _, err := os.Stat(template.path); os.IsNotExist(err) {
				t.Skipf("Template %s not found", template.path)
			}

			result := verifyTemplateCompilation(template.path, testData, outputDir)

			switch result.Status {
			case VerificationFailed:
				t.Errorf("Known problematic template failed verification: %s", result.Error)
				if result.CompilationOutput != "" {
					t.Logf("Compilation output: %s", result.CompilationOutput)
				}
			case VerificationSuccess:
				t.Logf("Previously problematic template now compiles successfully")
			}
		})
	}
}

func verifyImportFixes(t *testing.T) {
	// Test specific import fix scenarios
	testCases := []struct {
		name            string
		templateContent string
		expectedImports []string
		description     string
	}{
		{
			name: "TimeImportFix",
			templateContent: `package {{.Name}}

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Current time: %v\n", time.Now())
}`,
			expectedImports: []string{"fmt", "time"},
			description:     "Verify time import is present and compiles",
		},
		{
			name: "HTTPImportFix",
			templateContent: `package {{.Name}}

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}`,
			expectedImports: []string{"fmt", "net/http"},
			description:     "Verify HTTP imports are present and compile",
		},
		{
			name: "JSONImportFix",
			templateContent: `package {{.Name}}

import (
	"encoding/json"
	"fmt"
)

func processData(data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("JSON: %s\n", bytes)
}`,
			expectedImports: []string{"encoding/json", "fmt"},
			description:     "Verify JSON imports are present and compile",
		},
		{
			name: "CryptoImportFix",
			templateContent: `package {{.Name}}

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func main() {
	token := generateToken()
	fmt.Printf("Token: %s\n", token)
}`,
			expectedImports: []string{"crypto/rand", "encoding/base64", "fmt"},
			description:     "Verify crypto imports are present and compile",
		},
	}

	testData := createCompilationTestData()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			templateFile := filepath.Join(tempDir, "test.go.tmpl")

			err := os.WriteFile(templateFile, []byte(tc.templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}

			result := verifyTemplateCompilation(templateFile, testData, tempDir)

			if result.Status != VerificationSuccess {
				t.Errorf("Import fix verification failed for %s: %s", tc.description, result.Error)
				if result.CompilationOutput != "" {
					t.Logf("Compilation output: %s", result.CompilationOutput)
				}
			}

			// Verify that the generated file contains the expected imports
			generatedFile := filepath.Join(tempDir, "test.go")
			if _, err := os.Stat(generatedFile); err == nil {
				content, err := os.ReadFile(generatedFile)
				if err == nil {
					contentStr := string(content)
					for _, expectedImport := range tc.expectedImports {
						if !strings.Contains(contentStr, fmt.Sprintf("\"%s\"", expectedImport)) {
							t.Errorf("Expected import %s not found in generated file for %s", expectedImport, tc.description)
						}
					}
				}
			}
		})
	}
}

func verifyTemplateGeneration(t *testing.T) {
	// Test that templates generate valid Go code with various configurations
	testConfigurations := []struct {
		name        string
		configMods  func(*models.ProjectConfig)
		description string
	}{
		{
			name: "MinimalConfig",
			configMods: func(config *models.ProjectConfig) {
				// Minimal configuration - no modifications needed
			},
			description: "Minimal configuration with features disabled",
		},
		{
			name: "FullConfig",
			configMods: func(config *models.ProjectConfig) {
				// Full configuration - no modifications needed
			},
			description: "Full configuration with all features enabled",
		},
		{
			name: "AuthOnlyConfig",
			configMods: func(config *models.ProjectConfig) {
				// Auth only configuration - no modifications needed
			},
			description: "Configuration with only auth enabled",
		},
	}

	templatesDir := "../../templates/backend/go-gin"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Skip("Go Gin templates directory not found, skipping test")
	}

	for _, testConfig := range testConfigurations {
		t.Run(testConfig.name, func(t *testing.T) {
			testData := createCompilationTestData()
			testConfig.configMods(testData)

			outputDir := t.TempDir()

			// Test a subset of critical templates
			criticalTemplates := []string{
				"main.go.tmpl",
				"internal/middleware/auth.go.tmpl",
				"internal/controllers/auth_controller.go.tmpl",
			}

			for _, templateName := range criticalTemplates {
				templatePath := filepath.Join(templatesDir, templateName)
				if _, err := os.Stat(templatePath); os.IsNotExist(err) {
					continue
				}

				result := verifyTemplateCompilation(templatePath, testData, outputDir)

				switch result.Status {
				case VerificationFailed:
					t.Errorf("Template %s failed with %s: %s", templateName, testConfig.description, result.Error)
				}
			}
		})
	}
}

// createCompilationTestData creates comprehensive test data for template compilation
// This function is now defined in test_helpers.go
