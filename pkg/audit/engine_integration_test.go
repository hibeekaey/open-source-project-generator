package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TestEngineIntegration tests the complete audit functionality
func TestEngineIntegration(t *testing.T) {
	// Create a temporary test project
	tempDir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create test files
	createTestProject(t, tempDir)

	// Create audit engine
	engine := NewEngine()

	// Test complete audit
	t.Run("Complete Audit", func(t *testing.T) {
		options := &interfaces.AuditOptions{
			Security:    true,
			Quality:     true,
			Licenses:    true,
			Performance: true,
		}

		result, err := engine.AuditProject(tempDir, options)
		if err != nil {
			t.Fatalf("AuditProject failed: %v", err)
		}

		// Verify result structure
		if result == nil {
			t.Fatal("Expected audit result, got nil")
		}

		if result.ProjectPath != tempDir {
			t.Errorf("Expected project path %s, got %s", tempDir, result.ProjectPath)
		}

		if result.AuditTime.IsZero() {
			t.Error("Expected audit time to be set")
		}

		// Verify all audit types were performed
		if result.Security == nil {
			t.Error("Expected security audit result")
		}

		if result.Quality == nil {
			t.Error("Expected quality audit result")
		}

		if result.Licenses == nil {
			t.Error("Expected license audit result")
		}

		if result.Performance == nil {
			t.Error("Expected performance audit result")
		}

		// Verify overall score is calculated
		if result.OverallScore < 0 || result.OverallScore > 100 {
			t.Errorf("Expected overall score between 0-100, got %f", result.OverallScore)
		}

		// Verify recommendations are generated
		if len(result.Recommendations) == 0 {
			t.Error("Expected recommendations to be generated")
		}
	})

	// Test individual audit components
	t.Run("Security Audit", func(t *testing.T) {
		result, err := engine.AuditSecurity(tempDir)
		if err != nil {
			t.Fatalf("AuditSecurity failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected security audit result, got nil")
		}

		if result.Score < 0 || result.Score > 100 {
			t.Errorf("Expected security score between 0-100, got %f", result.Score)
		}
	})

	t.Run("Quality Audit", func(t *testing.T) {
		result, err := engine.AuditCodeQuality(tempDir)
		if err != nil {
			t.Fatalf("AuditCodeQuality failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected quality audit result, got nil")
		}

		if result.Score < 0 || result.Score > 100 {
			t.Errorf("Expected quality score between 0-100, got %f", result.Score)
		}
	})

	t.Run("License Audit", func(t *testing.T) {
		result, err := engine.AuditLicenses(tempDir)
		if err != nil {
			t.Fatalf("AuditLicenses failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected license audit result, got nil")
		}

		if result.Score < 0 || result.Score > 100 {
			t.Errorf("Expected license score between 0-100, got %f", result.Score)
		}
	})

	t.Run("Performance Audit", func(t *testing.T) {
		result, err := engine.AuditPerformance(tempDir)
		if err != nil {
			t.Fatalf("AuditPerformance failed: %v", err)
		}

		if result == nil {
			t.Fatal("Expected performance audit result, got nil")
		}

		if result.Score < 0 || result.Score > 100 {
			t.Errorf("Expected performance score between 0-100, got %f", result.Score)
		}
	})
}

// TestEngineReportGeneration tests report generation functionality
func TestEngineReportGeneration(t *testing.T) {
	engine := NewEngine()

	// Create a mock audit result
	result := &interfaces.AuditResult{
		ProjectPath:     "/test/project",
		AuditTime:       time.Now(),
		OverallScore:    85.5,
		Security:        &interfaces.SecurityAuditResult{Score: 90.0},
		Quality:         &interfaces.QualityAuditResult{Score: 80.0},
		Licenses:        &interfaces.LicenseAuditResult{Score: 85.0},
		Performance:     &interfaces.PerformanceAuditResult{Score: 87.0},
		Recommendations: []string{"Test recommendation"},
	}

	t.Run("JSON Report", func(t *testing.T) {
		report, err := engine.GenerateAuditReport(result, "json")
		if err != nil {
			t.Fatalf("GenerateAuditReport failed: %v", err)
		}

		if len(report) == 0 {
			t.Error("Expected non-empty JSON report")
		}
	})

	t.Run("HTML Report", func(t *testing.T) {
		report, err := engine.GenerateAuditReport(result, "html")
		if err != nil {
			t.Fatalf("GenerateAuditReport failed: %v", err)
		}

		if len(report) == 0 {
			t.Error("Expected non-empty HTML report")
		}
	})

	t.Run("Markdown Report", func(t *testing.T) {
		report, err := engine.GenerateAuditReport(result, "markdown")
		if err != nil {
			t.Fatalf("GenerateAuditReport failed: %v", err)
		}

		if len(report) == 0 {
			t.Error("Expected non-empty Markdown report")
		}
	})

	t.Run("Unsupported Format", func(t *testing.T) {
		_, err := engine.GenerateAuditReport(result, "unsupported")
		if err == nil {
			t.Error("Expected error for unsupported format")
		}
	})
}

// TestEngineRuleManagement tests rule management functionality
func TestEngineRuleManagement(t *testing.T) {
	engine := NewEngine()

	t.Run("Get Default Rules", func(t *testing.T) {
		rules := engine.GetAuditRules()
		if len(rules) == 0 {
			t.Error("Expected default rules to be loaded")
		}
	})

	t.Run("Add Rule", func(t *testing.T) {
		rule := interfaces.AuditRule{
			ID:          "test-001",
			Name:        "Test Rule",
			Description: "Test rule description",
			Category:    interfaces.AuditCategorySecurity,
			Type:        interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityMedium,
			Enabled:     true,
		}

		err := engine.AddAuditRule(rule)
		if err != nil {
			t.Fatalf("AddAuditRule failed: %v", err)
		}

		// Verify rule was added
		rules := engine.GetAuditRules()
		found := false
		for _, r := range rules {
			if r.ID == "test-001" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected rule to be added")
		}
	})

	t.Run("Remove Rule", func(t *testing.T) {
		err := engine.RemoveAuditRule("test-001")
		if err != nil {
			t.Fatalf("RemoveAuditRule failed: %v", err)
		}

		// Verify rule was removed
		rules := engine.GetAuditRules()
		for _, r := range rules {
			if r.ID == "test-001" {
				t.Error("Expected rule to be removed")
			}
		}
	})
}

// TestEngineAuditSummary tests audit summary functionality
func TestEngineAuditSummary(t *testing.T) {
	engine := NewEngine()

	// Create mock audit results
	results := []*interfaces.AuditResult{
		{
			ProjectPath:  "/test/project1",
			AuditTime:    time.Now(),
			OverallScore: 85.0,
			Security:     &interfaces.SecurityAuditResult{Score: 90.0},
			Quality:      &interfaces.QualityAuditResult{Score: 80.0},
		},
		{
			ProjectPath:  "/test/project2",
			AuditTime:    time.Now(),
			OverallScore: 75.0,
			Security:     &interfaces.SecurityAuditResult{Score: 80.0},
			Quality:      &interfaces.QualityAuditResult{Score: 70.0},
		},
	}

	summary, err := engine.GetAuditSummary(results)
	if err != nil {
		t.Fatalf("GetAuditSummary failed: %v", err)
	}

	if summary == nil {
		t.Fatal("Expected audit summary, got nil")
	}

	if summary.TotalProjects != 2 {
		t.Errorf("Expected 2 projects, got %d", summary.TotalProjects)
	}

	expectedAverage := (85.0 + 75.0) / 2
	if summary.AverageScore != expectedAverage {
		t.Errorf("Expected average score %f, got %f", expectedAverage, summary.AverageScore)
	}
}

// TestEngineComponentDelegation tests that the engine properly delegates to components
func TestEngineComponentDelegation(t *testing.T) {
	engine := NewEngine()

	// Create a temporary test project
	tempDir, err := os.MkdirTemp("", "audit-delegation-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create minimal test files
	createTestProject(t, tempDir)

	t.Run("Security Scanner Delegation", func(t *testing.T) {
		// Test DetectSecrets
		_, err := engine.DetectSecrets(tempDir)
		if err != nil {
			t.Errorf("DetectSecrets failed: %v", err)
		}

		// Test ScanVulnerabilities
		_, err = engine.ScanVulnerabilities(tempDir)
		if err != nil {
			t.Errorf("ScanVulnerabilities failed: %v", err)
		}

		// Test CheckSecurityPolicies
		_, err = engine.CheckSecurityPolicies(tempDir)
		if err != nil {
			t.Errorf("CheckSecurityPolicies failed: %v", err)
		}
	})

	t.Run("Quality Analyzer Delegation", func(t *testing.T) {
		// Test MeasureComplexity
		_, err := engine.MeasureComplexity(tempDir)
		if err != nil {
			t.Errorf("MeasureComplexity failed: %v", err)
		}
	})

	t.Run("License Checker Delegation", func(t *testing.T) {
		// Test CheckLicenseCompatibility
		_, err := engine.CheckLicenseCompatibility(tempDir)
		if err != nil {
			t.Errorf("CheckLicenseCompatibility failed: %v", err)
		}

		// Test ScanLicenseViolations
		_, err = engine.ScanLicenseViolations(tempDir)
		if err != nil {
			t.Errorf("ScanLicenseViolations failed: %v", err)
		}
	})

	t.Run("Performance Analyzer Delegation", func(t *testing.T) {
		// Test AnalyzeBundleSize
		_, err := engine.AnalyzeBundleSize(tempDir)
		if err != nil {
			t.Errorf("AnalyzeBundleSize failed: %v", err)
		}

		// Test CheckPerformanceMetrics
		_, err = engine.CheckPerformanceMetrics(tempDir)
		if err != nil {
			t.Errorf("CheckPerformanceMetrics failed: %v", err)
		}
	})

	t.Run("Dependency Analysis Delegation", func(t *testing.T) {
		// Test AnalyzeDependencies
		_, err := engine.AnalyzeDependencies(tempDir)
		if err != nil {
			t.Errorf("AnalyzeDependencies failed: %v", err)
		}
	})

	t.Run("Best Practices Check", func(t *testing.T) {
		// Test CheckBestPractices
		_, err := engine.CheckBestPractices(tempDir)
		if err != nil {
			t.Errorf("CheckBestPractices failed: %v", err)
		}
	})
}

// createTestProject creates a minimal test project structure
func createTestProject(t *testing.T, dir string) {
	// Create a simple Go file
	goFile := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

func complexFunction(a, b, c int) int {
	if a > 0 {
		if b > 0 {
			if c > 0 {
				return a + b + c
			}
		}
	}
	return 0
}
`
	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Create a test file
	testFile := `package main

import "testing"

func TestMain(t *testing.T) {
	// Test implementation
}
`
	err = os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(testFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create main_test.go: %v", err)
	}

	// Create go.mod
	goMod := `module test-project

go 1.21

require (
	github.com/example/dependency v1.0.0
)
`
	err = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create package.json for JavaScript testing
	packageJSON := `{
	"name": "test-project",
	"version": "1.0.0",
	"description": "Test project",
	"dependencies": {
		"lodash": "^4.17.21"
	}
}
`
	err = os.WriteFile(filepath.Join(dir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create a JavaScript file
	jsFile := `function complexFunction(a, b, c) {
	if (a > 0) {
		if (b > 0) {
			if (c > 0) {
				return a + b + c;
			}
		}
	}
	return 0;
}

console.log("Debug statement");
`
	err = os.WriteFile(filepath.Join(dir, "index.js"), []byte(jsFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create index.js: %v", err)
	}

	// Create README
	readme := `# Test Project

This is a test project for audit functionality.
`
	err = os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Create LICENSE
	license := `MIT License

Copyright (c) 2024 Test Project

Permission is hereby granted...
`
	err = os.WriteFile(filepath.Join(dir, "LICENSE"), []byte(license), 0644)
	if err != nil {
		t.Fatalf("Failed to create LICENSE: %v", err)
	}

	// Create .gitignore
	gitignore := `node_modules/
*.log
.env
`
	err = os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignore), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}
}
