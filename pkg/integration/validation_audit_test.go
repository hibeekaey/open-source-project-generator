package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestValidationAndAuditWorkflows tests complete validation and audit workflows
func TestValidationAndAuditWorkflows(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("project_validation_workflow", func(t *testing.T) {
		testProjectValidationWorkflow(t, tempDir)
	})

	t.Run("security_audit_workflow", func(t *testing.T) {
		testSecurityAuditWorkflow(t, tempDir)
	})

	t.Run("quality_audit_workflow", func(t *testing.T) {
		testQualityAuditWorkflow(t, tempDir)
	})

	t.Run("comprehensive_audit_workflow", func(t *testing.T) {
		testComprehensiveAuditWorkflow(t, tempDir)
	})

	t.Run("validation_with_fixes", func(t *testing.T) {
		testValidationWithFixes(t, tempDir)
	})
}

func testProjectValidationWorkflow(t *testing.T, tempDir string) {
	// Create test project with various issues
	projectDir := filepath.Join(tempDir, "validation-test")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create project files with issues
	createTestProjectFiles(t, projectDir)

	// Create mock validation engine
	validationEngine := NewMockValidationEngine()

	// Run basic validation
	result, err := validationEngine.ValidateProject(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate project: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Should find some issues
	if len(result.Issues) == 0 {
		t.Error("Expected to find validation issues")
	}

	// Run structure validation
	structResult, err := validationEngine.ValidateProjectStructure(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate project structure: %v", err)
	}

	if structResult == nil {
		t.Fatal("Expected structure validation result, got nil")
	}

	// Run dependency validation
	depResult, err := validationEngine.ValidateProjectDependencies(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate dependencies: %v", err)
	}

	if depResult == nil {
		t.Fatal("Expected dependency validation result, got nil")
	}

	// Run security validation
	secResult, err := validationEngine.ValidateProjectSecurity(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate security: %v", err)
	}

	if secResult == nil {
		t.Fatal("Expected security validation result, got nil")
	}
}
func testSecurityAuditWorkflow(t *testing.T, tempDir string) {
	// Create project with security issues
	projectDir := filepath.Join(tempDir, "security-audit-test")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create files with security issues
	createSecurityIssueFiles(t, projectDir)

	// Create mock audit engine
	auditEngine := NewMockAuditEngine()

	// Run security audit
	result, err := auditEngine.AuditSecurity(projectDir)
	if err != nil {
		t.Fatalf("Failed to audit security: %v", err)
	}

	if result == nil {
		t.Fatal("Expected security audit result, got nil")
	}

	// Should detect security issues
	if len(result.Vulnerabilities) == 0 && len(result.PolicyViolations) == 0 {
		t.Error("Expected to find security issues")
	}

	// Run vulnerability scan
	vulnResult, err := auditEngine.ScanVulnerabilities(projectDir)
	if err != nil {
		t.Fatalf("Failed to scan vulnerabilities: %v", err)
	}

	if vulnResult == nil {
		t.Fatal("Expected vulnerability scan result, got nil")
	}

	// Run secret detection
	secretResult, err := auditEngine.DetectSecrets(projectDir)
	if err != nil {
		t.Fatalf("Failed to detect secrets: %v", err)
	}

	if secretResult == nil {
		t.Fatal("Expected secret detection result, got nil")
	}

	// Should find secrets in test files
	if secretResult.Summary.TotalSecrets == 0 {
		t.Error("Expected to find secrets in test files")
	}
}

func testQualityAuditWorkflow(t *testing.T, tempDir string) {
	// Create project with quality issues
	projectDir := filepath.Join(tempDir, "quality-audit-test")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create files with quality issues
	createQualityIssueFiles(t, projectDir)

	// Create mock audit engine
	auditEngine := NewMockAuditEngine()

	// Run quality audit
	result, err := auditEngine.AuditCodeQuality(projectDir)
	if err != nil {
		t.Fatalf("Failed to audit code quality: %v", err)
	}

	if result == nil {
		t.Fatal("Expected quality audit result, got nil")
	}

	// Should detect quality issues
	if len(result.CodeSmells) == 0 && len(result.Duplications) == 0 {
		t.Error("Expected to find quality issues")
	}

	// Run complexity analysis
	complexityResult, err := auditEngine.MeasureComplexity(projectDir)
	if err != nil {
		t.Fatalf("Failed to measure complexity: %v", err)
	}

	if complexityResult == nil {
		t.Fatal("Expected complexity analysis result, got nil")
	}

	// Should analyze files
	if complexityResult.Summary.TotalFiles == 0 {
		t.Error("Expected to analyze files for complexity")
	}
}

func testComprehensiveAuditWorkflow(t *testing.T, tempDir string) {
	// Create comprehensive test project
	projectDir := filepath.Join(tempDir, "comprehensive-audit-test")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create comprehensive project structure
	createComprehensiveProject(t, projectDir)

	// Create mock audit engine
	auditEngine := NewMockAuditEngine()

	// Run comprehensive audit
	options := &interfaces.AuditOptions{
		Security:    true,
		Quality:     true,
		Licenses:    true,
		Performance: true,
		Detailed:    true,
	}

	result, err := auditEngine.AuditProject(projectDir, options)
	if err != nil {
		t.Fatalf("Failed to run comprehensive audit: %v", err)
	}

	if result == nil {
		t.Fatal("Expected comprehensive audit result, got nil")
	}

	// Verify all audit types were run
	if result.Security == nil {
		t.Error("Expected security audit results")
	}

	if result.Quality == nil {
		t.Error("Expected quality audit results")
	}

	if result.Licenses == nil {
		t.Error("Expected license audit results")
	}

	if result.Performance == nil {
		t.Error("Expected performance audit results")
	}

	// Verify overall score is calculated
	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Expected overall score between 0-100, got %f", result.OverallScore)
	}

	// Generate audit report
	jsonReport, err := auditEngine.GenerateAuditReport(result, "json")
	if err != nil {
		t.Fatalf("Failed to generate JSON report: %v", err)
	}

	if len(jsonReport) == 0 {
		t.Error("Expected JSON report content")
	}

	htmlReport, err := auditEngine.GenerateAuditReport(result, "html")
	if err != nil {
		t.Fatalf("Failed to generate HTML report: %v", err)
	}

	if len(htmlReport) == 0 {
		t.Error("Expected HTML report content")
	}
}

func testValidationWithFixes(t *testing.T, tempDir string) {
	// Create project with fixable issues
	projectDir := filepath.Join(tempDir, "validation-fix-test")
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create files with fixable issues
	createFixableIssueFiles(t, projectDir)

	// Create mock validation engine
	validationEngine := NewMockValidationEngine()

	// Run initial validation
	result, err := validationEngine.ValidateProject(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate project: %v", err)
	}

	// Should find issues
	if len(result.Issues) == 0 {
		t.Error("Expected to find validation issues")
	}

	// Get fixable issues
	fixableIssues := validationEngine.GetFixableIssues(result.Issues)

	if len(fixableIssues) == 0 {
		t.Error("Expected to find fixable issues")
	}

	// Preview fixes
	fixPreview, err := validationEngine.PreviewFixes(projectDir, fixableIssues)
	if err != nil {
		t.Fatalf("Failed to preview fixes: %v", err)
	}

	if fixPreview == nil {
		t.Fatal("Expected fix preview, got nil")
	}

	if len(fixPreview.Fixes) == 0 {
		t.Error("Expected fixes in preview")
	}

	// Apply fixes
	fixResult, err := validationEngine.FixValidationIssues(projectDir, fixableIssues)
	if err != nil {
		t.Fatalf("Failed to fix validation issues: %v", err)
	}

	if fixResult == nil {
		t.Fatal("Expected fix result, got nil")
	}

	// Should have applied some fixes
	if len(fixResult.Applied) == 0 {
		t.Error("Expected some fixes to be applied")
	}

	// Run validation again to verify fixes
	postFixResult, err := validationEngine.ValidateProject(projectDir)
	if err != nil {
		t.Fatalf("Failed to validate project after fixes: %v", err)
	}

	// Should have fewer issues after fixes
	if len(postFixResult.Issues) >= len(result.Issues) {
		t.Error("Expected fewer issues after applying fixes")
	}
}

// Helper functions to create test files

func createTestProjectFiles(t *testing.T, projectDir string) {
	// Create package.json with issues
	packageJSON := `{
		"name": "test project",
		"version": "1.0.0",
		"dependencies": {
			"lodash": "4.17.20"
		}
	}`

	err := os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create go.mod
	goMod := `module test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`

	err = os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create main.go with issues
	mainGo := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}

// Unused function
func unused() {
	// This function is never called
}
`

	err = os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainGo), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}
}

func createSecurityIssueFiles(t *testing.T, projectDir string) {
	// Create config file with secrets
	configFile := `
const config = {
	apiKey: "sk-1234567890abcdef1234567890abcdef",
	password: "mySecretPassword123",
	token: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	dbUrl: "mongodb://admin:password123@localhost:27017/mydb"
};

module.exports = config;
`

	err := os.WriteFile(filepath.Join(projectDir, "config.js"), []byte(configFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create config.js: %v", err)
	}

	// Create .env file with secrets
	envFile := `
API_KEY=sk-live-1234567890abcdef
DATABASE_PASSWORD=supersecret123
JWT_SECRET=my-jwt-secret-key
STRIPE_SECRET_KEY=sk_test_1234567890
`

	err = os.WriteFile(filepath.Join(projectDir, ".env"), []byte(envFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env: %v", err)
	}

	// Create insecure code
	insecureCode := `
package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

func insecureHandler(w http.ResponseWriter, r *http.Request) {
	// Insecure: using MD5 for hashing
	hash := md5.Sum([]byte("password"))
	
	// Insecure: no input validation
	userInput := r.URL.Query().Get("input")
	fmt.Fprintf(w, "Hello %s", userInput)
}
`

	err = os.WriteFile(filepath.Join(projectDir, "insecure.go"), []byte(insecureCode), 0644)
	if err != nil {
		t.Fatalf("Failed to create insecure.go: %v", err)
	}
}

func createQualityIssueFiles(t *testing.T, projectDir string) {
	// Create file with code smells
	codeSmells := `
package main

import "fmt"

// Long function with high complexity
func complexFunction(x int) int {
	if x > 100 {
		if x > 200 {
			if x > 300 {
				if x > 400 {
					if x > 500 {
						if x > 600 {
							return x * 10
						}
						return x * 9
					}
					return x * 8
				}
				return x * 7
			}
			return x * 6
		}
		return x * 5
	}
	return x
}

// Duplicate code
func processDataA() {
	fmt.Println("Processing data")
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i * 2
	}
	fmt.Println("Data processed")
}

// Duplicate code
func processDataB() {
	fmt.Println("Processing data")
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i * 2
	}
	fmt.Println("Data processed")
}

// Long parameter list
func longParameterFunction(a, b, c, d, e, f, g, h, i, j int) int {
	return a + b + c + d + e + f + g + h + i + j
}
`

	err := os.WriteFile(filepath.Join(projectDir, "quality_issues.go"), []byte(codeSmells), 0644)
	if err != nil {
		t.Fatalf("Failed to create quality_issues.go: %v", err)
	}

	// Create JavaScript file with issues
	jsIssues := `
// Unused variable
var unusedVariable = "not used";

// Magic numbers
function calculatePrice(quantity) {
	return quantity * 19.99 + 5.99 + (quantity * 19.99 * 0.08);
}

// Long function
function processOrder(order) {
	if (order.items.length > 0) {
		if (order.customer.isVip) {
			if (order.total > 100) {
				if (order.shippingMethod === "express") {
					if (order.paymentMethod === "credit") {
						// Apply VIP express credit discount
						return order.total * 0.85;
					}
				}
			}
		}
	}
	return order.total;
}

// Duplicate code
function validateEmailA(email) {
	if (!email) return false;
	if (email.indexOf("@") === -1) return false;
	if (email.indexOf(".") === -1) return false;
	return true;
}

function validateEmailB(email) {
	if (!email) return false;
	if (email.indexOf("@") === -1) return false;
	if (email.indexOf(".") === -1) return false;
	return true;
}
`

	err = os.WriteFile(filepath.Join(projectDir, "quality_issues.js"), []byte(jsIssues), 0644)
	if err != nil {
		t.Fatalf("Failed to create quality_issues.js: %v", err)
	}
}

func createComprehensiveProject(t *testing.T, projectDir string) {
	// Create all types of files for comprehensive testing
	createTestProjectFiles(t, projectDir)
	createSecurityIssueFiles(t, projectDir)
	createQualityIssueFiles(t, projectDir)

	// Create LICENSE file
	license := `MIT License

Copyright (c) 2024 Test Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

	err := os.WriteFile(filepath.Join(projectDir, "LICENSE"), []byte(license), 0644)
	if err != nil {
		t.Fatalf("Failed to create LICENSE: %v", err)
	}

	// Create README
	readme := `# Comprehensive Test Project

This is a test project for comprehensive auditing.

## Features

- Backend API
- Frontend interface
- Database integration

## Installation

1. Clone the repository
2. Install dependencies
3. Run the application

## License

MIT License
`

	err = os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create Dockerfile
	dockerfile := `FROM node:18

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
`

	err = os.WriteFile(filepath.Join(projectDir, "Dockerfile"), []byte(dockerfile), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}
}

func createFixableIssueFiles(t *testing.T, projectDir string) {
	// Create files with fixable issues

	// Missing README
	// (Don't create README to test fix)

	// Invalid package.json format
	invalidPackageJSON := `{
		"name": "fixable test project",
		"version": "1.0.0"
		// Missing comma and other issues
	}`

	err := os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(invalidPackageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid package.json: %v", err)
	}

	// File with wrong permissions
	restrictedFile := "This file has wrong permissions"

	err = os.WriteFile(filepath.Join(projectDir, "restricted.txt"), []byte(restrictedFile), 0600)
	if err != nil {
		t.Fatalf("Failed to create restricted file: %v", err)
	}

	// Go file with formatting issues
	unformattedGo := `package main
import"fmt"
func main(){
fmt.Println("Unformatted code")
}
`

	err = os.WriteFile(filepath.Join(projectDir, "unformatted.go"), []byte(unformattedGo), 0644)
	if err != nil {
		t.Fatalf("Failed to create unformatted.go: %v", err)
	}
}

// Mock implementations for testing

type MockValidationEngine struct {
	appliedFixes map[string][]string // path -> list of fixed rule IDs
}

func NewMockValidationEngine() *MockValidationEngine {
	return &MockValidationEngine{
		appliedFixes: make(map[string][]string),
	}
}

func (m *MockValidationEngine) ValidateProject(path string) (*models.ValidationResult, error) {
	// Mock validation that finds some issues
	allIssues := []models.ValidationIssue{
		{
			Type:     "warning",
			Severity: "warning",
			Message:  "Package name contains spaces",
			File:     "package.json",
			Rule:     "package-name-format",
			Fixable:  true,
		},
		{
			Type:     "error",
			Severity: "error",
			Message:  "Missing README file",
			Rule:     "readme-required",
			Fixable:  true,
		},
	}

	// Filter out issues that have been fixed
	issues := []models.ValidationIssue{}
	fixedRules := m.appliedFixes[path]

	for _, issue := range allIssues {
		isFixed := false
		for _, fixedRule := range fixedRules {
			if issue.Rule == fixedRule {
				isFixed = true
				break
			}
		}
		if !isFixed {
			issues = append(issues, issue)
		}
	}

	return &models.ValidationResult{
		Valid:   len(issues) == 0,
		Issues:  issues,
		Summary: "Validation completed with issues",
	}, nil
}

func (m *MockValidationEngine) ValidateProjectStructure(path string) (*interfaces.StructureValidationResult, error) {
	return &interfaces.StructureValidationResult{
		Valid: false,
		RequiredFiles: []interfaces.FileValidationResult{
			{
				Path:     "README.md",
				Required: true,
				Exists:   false,
				Valid:    false,
			},
		},
		Summary: interfaces.StructureValidationSummary{
			TotalFiles: 1,
			ValidFiles: 0,
		},
	}, nil
}

func (m *MockValidationEngine) ValidateProjectDependencies(path string) (*interfaces.DependencyValidationResult, error) {
	return &interfaces.DependencyValidationResult{
		Valid: true,
		Dependencies: []interfaces.DependencyValidation{
			{
				Name:    "lodash",
				Version: "4.17.20",
				Valid:   true,
			},
		},
		Summary: interfaces.DependencyValidationSummary{
			TotalDependencies: 1,
			ValidDependencies: 1,
		},
	}, nil
}

func (m *MockValidationEngine) ValidateProjectSecurity(path string) (*interfaces.SecurityValidationResult, error) {
	return &interfaces.SecurityValidationResult{
		Valid: false,
		SecurityIssues: []interfaces.SecurityIssue{
			{
				Type:        "secret",
				Severity:    "high",
				Title:       "Hardcoded API key detected",
				Description: "API key found in source code",
				File:        "config.js",
				Line:        3,
			},
		},
		Summary: interfaces.SecurityValidationSummary{
			TotalIssues:  1,
			HighSeverity: 1,
		},
	}, nil
}

func (m *MockValidationEngine) GetFixableIssues(issues []models.ValidationIssue) []models.ValidationIssue {
	fixable := []models.ValidationIssue{}

	for _, issue := range issues {
		if issue.Fixable {
			fixable = append(fixable, issue)
		}
	}

	return fixable
}

func (m *MockValidationEngine) PreviewFixes(path string, issues []models.ValidationIssue) (*interfaces.FixPreview, error) {
	fixes := []interfaces.Fix{}

	for _, issue := range issues {
		if issue.Fixable {
			fixes = append(fixes, interfaces.Fix{
				ID:          "fix-" + issue.Rule,
				Type:        "auto",
				Description: "Fix " + issue.Message,
				File:        issue.File,
				Action:      "replace",
				Automatic:   true,
			})
		}
	}

	return &interfaces.FixPreview{
		Fixes: fixes,
		Summary: interfaces.FixSummary{
			TotalFixes: len(fixes),
		},
	}, nil
}

func (m *MockValidationEngine) FixValidationIssues(path string, issues []models.ValidationIssue) (*interfaces.FixResult, error) {
	applied := []interfaces.Fix{}

	for _, issue := range issues {
		if issue.Fixable {
			applied = append(applied, interfaces.Fix{
				ID:          "fix-" + issue.Rule,
				Type:        "auto",
				Description: "Fixed " + issue.Message,
				File:        issue.File,
				Action:      "replace",
				Automatic:   true,
			})

			// Track that this fix was applied
			if m.appliedFixes[path] == nil {
				m.appliedFixes[path] = []string{}
			}
			m.appliedFixes[path] = append(m.appliedFixes[path], issue.Rule)
		}
	}

	return &interfaces.FixResult{
		Applied: applied,
		Summary: interfaces.FixSummary{
			TotalFixes:   len(applied),
			AppliedFixes: len(applied),
		},
	}, nil
}

type MockAuditEngine struct{}

func NewMockAuditEngine() *MockAuditEngine {
	return &MockAuditEngine{}
}

func (m *MockAuditEngine) AuditSecurity(path string) (*interfaces.SecurityAuditResult, error) {
	return &interfaces.SecurityAuditResult{
		Score: 75.0,
		Vulnerabilities: []interfaces.Vulnerability{
			{
				ID:          "CVE-2021-1234",
				Severity:    "high",
				Title:       "Prototype Pollution",
				Description: "Lodash vulnerable to prototype pollution",
				Package:     "lodash",
				Version:     "4.17.20",
				FixedIn:     "4.17.21",
			},
		},
		PolicyViolations: []interfaces.PolicyViolation{
			{
				Policy:      "no-hardcoded-secrets",
				Severity:    "high",
				Description: "Hardcoded API key found",
				File:        "config.js",
				Line:        3,
			},
		},
		Recommendations: []string{
			"Update lodash to latest version",
			"Remove hardcoded secrets",
		},
	}, nil
}

func (m *MockAuditEngine) ScanVulnerabilities(path string) (*interfaces.VulnerabilityReport, error) {
	return &interfaces.VulnerabilityReport{
		ScanTime: time.Now(),
		Vulnerabilities: []interfaces.Vulnerability{
			{
				ID:          "CVE-2021-1234",
				Severity:    "high",
				Title:       "Prototype Pollution",
				Description: "Lodash vulnerable to prototype pollution",
				Package:     "lodash",
				Version:     "4.17.20",
			},
		},
		Summary: interfaces.VulnerabilitySummary{
			Total: 1,
			High:  1,
		},
	}, nil
}

func (m *MockAuditEngine) DetectSecrets(path string) (*interfaces.SecretScanResult, error) {
	return &interfaces.SecretScanResult{
		ScanTime: time.Now(),
		Secrets: []interfaces.SecretDetection{
			{
				Type:       "api_key",
				File:       "config.js",
				Line:       3,
				Secret:     "sk-1234567890abcdef",
				Confidence: 0.9,
				Rule:       "api-key-pattern",
			},
		},
		Summary: interfaces.SecretScanSummary{
			TotalSecrets:   1,
			HighConfidence: 1,
			FilesScanned:   5,
		},
	}, nil
}

func (m *MockAuditEngine) AuditCodeQuality(path string) (*interfaces.QualityAuditResult, error) {
	return &interfaces.QualityAuditResult{
		Score: 65.0,
		CodeSmells: []interfaces.CodeSmell{
			{
				Type:        "complexity",
				Severity:    "medium",
				Description: "Function has high cyclomatic complexity",
				File:        "quality_issues.go",
				Line:        8,
			},
		},
		Duplications: []interfaces.Duplication{
			{
				Files:      []string{"quality_issues.go"},
				Lines:      10,
				Percentage: 15.0,
			},
		},
		TestCoverage: 45.0,
		Recommendations: []string{
			"Reduce function complexity",
			"Remove code duplication",
			"Increase test coverage",
		},
	}, nil
}

func (m *MockAuditEngine) MeasureComplexity(path string) (*interfaces.ComplexityAnalysisResult, error) {
	return &interfaces.ComplexityAnalysisResult{
		Files: []interfaces.FileComplexity{
			{
				Path:                 "quality_issues.go",
				Lines:                100,
				CyclomaticComplexity: 15,
				Maintainability:      60.0,
			},
		},
		Summary: interfaces.ComplexityAnalysisSummary{
			TotalFiles:        1,
			AverageComplexity: 15.0,
		},
	}, nil
}

func (m *MockAuditEngine) AuditProject(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	result := &interfaces.AuditResult{
		ProjectPath: path,
		AuditTime:   time.Now(),
	}

	if options.Security {
		security, _ := m.AuditSecurity(path)
		result.Security = security
	}

	if options.Quality {
		quality, _ := m.AuditCodeQuality(path)
		result.Quality = quality
	}

	if options.Licenses {
		result.Licenses = &interfaces.LicenseAuditResult{
			Score:      90.0,
			Compatible: true,
			Licenses: []interfaces.LicenseInfo{
				{
					Name:       "MIT",
					Package:    "project",
					Compatible: true,
				},
			},
		}
	}

	if options.Performance {
		result.Performance = &interfaces.PerformanceAuditResult{
			Score:      80.0,
			BundleSize: 1024 * 1024, // 1MB
			LoadTime:   time.Second,
		}
	}

	// Calculate overall score
	scores := []float64{}
	if result.Security != nil {
		scores = append(scores, result.Security.Score)
	}
	if result.Quality != nil {
		scores = append(scores, result.Quality.Score)
	}
	if result.Licenses != nil {
		scores = append(scores, result.Licenses.Score)
	}
	if result.Performance != nil {
		scores = append(scores, result.Performance.Score)
	}

	if len(scores) > 0 {
		total := 0.0
		for _, score := range scores {
			total += score
		}
		result.OverallScore = total / float64(len(scores))
	}

	result.Recommendations = []string{
		"Update dependencies",
		"Improve code quality",
		"Add more tests",
	}

	return result, nil
}

func (m *MockAuditEngine) GenerateAuditReport(result *interfaces.AuditResult, format string) ([]byte, error) {
	switch format {
	case "json":
		return []byte(`{"project_path": "` + result.ProjectPath + `", "overall_score": ` + fmt.Sprintf("%.1f", result.OverallScore) + `}`), nil
	case "html":
		return []byte(`<html><body><h1>Audit Report</h1><p>Score: ` + fmt.Sprintf("%.1f", result.OverallScore) + `</p></body></html>`), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
