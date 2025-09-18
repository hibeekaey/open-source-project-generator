package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("Expected engine to be created, got nil")
	}

	// Verify default rules are loaded
	auditEngine := engine.(*Engine)
	if len(auditEngine.rules) == 0 {
		t.Error("Expected default audit rules to be loaded")
	}
}

func TestEngine_AuditSecurity(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create a test file with potential security issues
	testFile := `
const config = {
	apiKey: "sk-1234567890abcdef",
	password: "secret123",
	token: "ghp_xxxxxxxxxxxxxxxxxxxx"
};
`
	err := os.WriteFile(filepath.Join(tempDir, "config.js"), []byte(testFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := engine.AuditSecurity(tempDir)
	if err != nil {
		t.Fatalf("AuditSecurity failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected security audit result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}

	// Test with non-existent directory
	_, err = engine.AuditSecurity("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestEngine_ScanVulnerabilities(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with package.json
	tempDir := t.TempDir()
	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"dependencies": {
			"lodash": "4.17.20",
			"express": "4.17.1"
		}
	}`

	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	result, err := engine.ScanVulnerabilities(tempDir)
	if err != nil {
		t.Fatalf("ScanVulnerabilities failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected vulnerability report, got nil")
	}

	if result.ScanTime.IsZero() {
		t.Error("Expected scan time to be set")
	}

	if result.Summary.Total < 0 {
		t.Error("Expected non-negative total vulnerabilities")
	}
}

func TestEngine_CheckSecurityPolicies(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	result, err := engine.CheckSecurityPolicies(tempDir)
	if err != nil {
		t.Fatalf("CheckSecurityPolicies failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected policy compliance result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}

	if len(result.Policies) == 0 {
		t.Error("Expected security policies to be checked")
	}
}

func TestEngine_AuditCodeQuality(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with test code
	tempDir := t.TempDir()

	// Create a Go file with quality issues
	goFile := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

// This is a duplicate function
func duplicate() {
	x := 1
	y := 2
	z := x + y
	fmt.Println(z)
}

// This is another duplicate function
func anotherDuplicate() {
	x := 1
	y := 2
	z := x + y
	fmt.Println(z)
}

func complexFunction() {
	if true {
		if true {
			if true {
				if true {
					if true {
						if true {
							fmt.Println("Too complex")
						}
					}
				}
			}
		}
	}
}
`

	err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	result, err := engine.AuditCodeQuality(tempDir)
	if err != nil {
		t.Fatalf("AuditCodeQuality failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected quality audit result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}
}

func TestEngine_CheckBestPractices(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	// Create README file (good practice)
	readme := "# Test Project\n\nThis is a test project."
	err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	result, err := engine.CheckBestPractices(tempDir)
	if err != nil {
		t.Fatalf("CheckBestPractices failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected best practices result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}

	if len(result.Practices) == 0 {
		t.Error("Expected best practices to be checked")
	}
}

func TestEngine_AnalyzeDependencies(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with multiple dependency files
	tempDir := t.TempDir()

	// Create package.json
	packageJSON := `{
		"name": "test-project",
		"version": "1.0.0",
		"dependencies": {
			"react": "^18.0.0",
			"lodash": "^4.17.21"
		},
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`

	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create go.mod
	goMod := `module test-project

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/stretchr/testify v1.8.4
)
`

	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	result, err := engine.AnalyzeDependencies(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeDependencies failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected dependency analysis result, got nil")
	}

	if result.Summary.TotalDependencies == 0 {
		t.Error("Expected dependencies to be found")
	}
}

func TestEngine_AuditLicenses(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	// Create LICENSE file
	license := `MIT License

Copyright (c) 2024 Test Project

Permission is hereby granted...`

	err := os.WriteFile(filepath.Join(tempDir, "LICENSE"), []byte(license), 0644)
	if err != nil {
		t.Fatalf("Failed to create LICENSE: %v", err)
	}

	result, err := engine.AuditLicenses(tempDir)
	if err != nil {
		t.Fatalf("AuditLicenses failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected license audit result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}
}

func TestEngine_CheckLicenseCompatibility(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	result, err := engine.CheckLicenseCompatibility(tempDir)
	if err != nil {
		t.Fatalf("CheckLicenseCompatibility failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected license compatibility result, got nil")
	}

	if result.Summary.RiskLevel == "" {
		t.Error("Expected risk level to be set")
	}
}

func TestEngine_AuditPerformance(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	result, err := engine.AuditPerformance(tempDir)
	if err != nil {
		t.Fatalf("AuditPerformance failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected performance audit result, got nil")
	}

	if result.Score < 0 || result.Score > 100 {
		t.Errorf("Expected score between 0-100, got %f", result.Score)
	}
}

func TestEngine_AnalyzeBundleSize(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with build output
	tempDir := t.TempDir()
	distDir := filepath.Join(tempDir, "dist")
	err := os.MkdirAll(distDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create dist directory: %v", err)
	}

	// Create some build files
	jsFile := "console.log('Hello, World!');"
	err = os.WriteFile(filepath.Join(distDir, "main.js"), []byte(jsFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create JS file: %v", err)
	}

	cssFile := "body { margin: 0; }"
	err = os.WriteFile(filepath.Join(distDir, "styles.css"), []byte(cssFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create CSS file: %v", err)
	}

	result, err := engine.AnalyzeBundleSize(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeBundleSize failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected bundle analysis result, got nil")
	}

	if result.TotalSize == 0 {
		t.Error("Expected total size to be greater than 0")
	}

	if len(result.Assets) == 0 {
		t.Error("Expected assets to be found")
	}
}

func TestEngine_DetectSecrets(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with files containing secrets
	tempDir := t.TempDir()

	// Create file with API key
	secretFile := `
const config = {
	apiKey: "sk-1234567890abcdef1234567890abcdef",
	password: "mySecretPassword123",
	token: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
};
`

	err := os.WriteFile(filepath.Join(tempDir, "config.js"), []byte(secretFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create secret file: %v", err)
	}

	result, err := engine.DetectSecrets(tempDir)
	if err != nil {
		t.Fatalf("DetectSecrets failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected secret scan result, got nil")
	}

	if result.Summary.FilesScanned == 0 {
		t.Error("Expected files to be scanned")
	}

	// Should detect some secrets
	if result.Summary.TotalSecrets == 0 {
		t.Error("Expected secrets to be detected")
	}
}

func TestEngine_MeasureComplexity(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory with code files
	tempDir := t.TempDir()

	// Create a complex Go function
	complexFile := `package main

func simpleFunction() {
	return
}

func complexFunction(x int) int {
	if x > 10 {
		if x > 20 {
			if x > 30 {
				if x > 40 {
					if x > 50 {
						return x * 2
					}
					return x + 10
				}
				return x + 5
			}
			return x + 2
		}
		return x + 1
	}
	return x
}
`

	err := os.WriteFile(filepath.Join(tempDir, "complex.go"), []byte(complexFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create complex file: %v", err)
	}

	result, err := engine.MeasureComplexity(tempDir)
	if err != nil {
		t.Fatalf("MeasureComplexity failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected complexity analysis result, got nil")
	}

	if result.Summary.TotalFiles == 0 {
		t.Error("Expected files to be analyzed")
	}
}

func TestEngine_AuditProject(t *testing.T) {
	engine := NewEngine()

	// Create temporary directory
	tempDir := t.TempDir()

	// Create basic project structure
	readme := "# Test Project"
	err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte(readme), 0644)
	if err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	// Test with default options
	result, err := engine.AuditProject(tempDir, nil)
	if err != nil {
		t.Fatalf("AuditProject failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected audit result, got nil")
	}

	if result.ProjectPath != tempDir {
		t.Errorf("Expected project path '%s', got '%s'", tempDir, result.ProjectPath)
	}

	if result.AuditTime.IsZero() {
		t.Error("Expected audit time to be set")
	}

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Expected overall score between 0-100, got %f", result.OverallScore)
	}

	// Test with specific options
	options := &interfaces.AuditOptions{
		Security:    true,
		Quality:     false,
		Licenses:    true,
		Performance: false,
	}

	result, err = engine.AuditProject(tempDir, options)
	if err != nil {
		t.Fatalf("AuditProject with options failed: %v", err)
	}

	if result.Security == nil {
		t.Error("Expected security audit result")
	}

	if result.Quality != nil {
		t.Error("Expected no quality audit result")
	}

	if result.Licenses == nil {
		t.Error("Expected license audit result")
	}

	if result.Performance != nil {
		t.Error("Expected no performance audit result")
	}
}

func TestEngine_GenerateAuditReport(t *testing.T) {
	engine := NewEngine()

	// Create sample audit result
	result := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Now(),
		OverallScore: 85.5,
		Security: &interfaces.SecurityAuditResult{
			Score: 90.0,
		},
		Quality: &interfaces.QualityAuditResult{
			Score: 80.0,
		},
		Recommendations: []string{
			"Update dependencies",
			"Add more tests",
		},
	}

	// Test JSON format
	jsonReport, err := engine.GenerateAuditReport(result, "json")
	if err != nil {
		t.Fatalf("GenerateAuditReport JSON failed: %v", err)
	}

	if len(jsonReport) == 0 {
		t.Error("Expected JSON report content")
	}

	// Test HTML format
	htmlReport, err := engine.GenerateAuditReport(result, "html")
	if err != nil {
		t.Fatalf("GenerateAuditReport HTML failed: %v", err)
	}

	if len(htmlReport) == 0 {
		t.Error("Expected HTML report content")
	}

	// Test Markdown format
	mdReport, err := engine.GenerateAuditReport(result, "markdown")
	if err != nil {
		t.Fatalf("GenerateAuditReport Markdown failed: %v", err)
	}

	if len(mdReport) == 0 {
		t.Error("Expected Markdown report content")
	}

	// Test unsupported format
	_, err = engine.GenerateAuditReport(result, "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported format")
	}
}

func TestEngine_GetAuditSummary(t *testing.T) {
	engine := NewEngine()

	// Create sample audit results
	results := []*interfaces.AuditResult{
		{
			ProjectPath:  "/test/project1",
			OverallScore: 85.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 90.0,
			},
		},
		{
			ProjectPath:  "/test/project2",
			OverallScore: 75.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 80.0,
			},
		},
	}

	summary, err := engine.GetAuditSummary(results)
	if err != nil {
		t.Fatalf("GetAuditSummary failed: %v", err)
	}

	if summary == nil {
		t.Fatal("Expected audit summary, got nil")
	}

	if summary.TotalProjects != len(results) {
		t.Errorf("Expected %d projects, got %d", len(results), summary.TotalProjects)
	}

	expectedAverage := (85.0 + 75.0) / 2
	if summary.AverageScore != expectedAverage {
		t.Errorf("Expected average score %f, got %f", expectedAverage, summary.AverageScore)
	}

	// Test with empty results
	emptySummary, err := engine.GetAuditSummary([]*interfaces.AuditResult{})
	if err != nil {
		t.Fatalf("GetAuditSummary with empty results failed: %v", err)
	}

	if emptySummary.TotalProjects != 0 {
		t.Error("Expected 0 projects for empty results")
	}
}

func TestEngine_SetAndGetAuditRules(t *testing.T) {
	engine := NewEngine()

	// Test getting default rules
	rules := engine.GetAuditRules()
	if len(rules) == 0 {
		t.Error("Expected default audit rules")
	}

	// Test setting custom rules
	customRules := []interfaces.AuditRule{
		{
			ID:          "custom-1",
			Name:        "Custom Rule 1",
			Description: "A custom audit rule",
			Category:    interfaces.AuditCategorySecurity,
			Severity:    interfaces.AuditSeverityHigh,
			Enabled:     true,
		},
	}

	err := engine.SetAuditRules(customRules)
	if err != nil {
		t.Fatalf("SetAuditRules failed: %v", err)
	}

	updatedRules := engine.GetAuditRules()
	if len(updatedRules) != len(customRules) {
		t.Errorf("Expected %d rules, got %d", len(customRules), len(updatedRules))
	}

	// Test adding a rule
	newRule := interfaces.AuditRule{
		ID:          "custom-2",
		Name:        "Custom Rule 2",
		Description: "Another custom rule",
		Category:    interfaces.AuditCategoryQuality,
		Severity:    interfaces.AuditSeverityMedium,
		Enabled:     true,
	}

	err = engine.AddAuditRule(newRule)
	if err != nil {
		t.Fatalf("AddAuditRule failed: %v", err)
	}

	finalRules := engine.GetAuditRules()
	if len(finalRules) != len(customRules)+1 {
		t.Error("Expected rule count to increase by 1")
	}

	// Test removing a rule
	err = engine.RemoveAuditRule("custom-2")
	if err != nil {
		t.Fatalf("RemoveAuditRule failed: %v", err)
	}

	afterRemoval := engine.GetAuditRules()
	if len(afterRemoval) != len(customRules) {
		t.Error("Expected rule count to decrease by 1")
	}
}

// Test helper methods

func TestEngine_projectExists(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Test with existing directory
	tempDir := t.TempDir()
	err := engine.projectExists(tempDir)
	if err != nil {
		t.Errorf("Expected no error for existing directory, got: %v", err)
	}

	// Test with non-existent directory
	err = engine.projectExists("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestEngine_shouldSkipFile(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		filename   string
		shouldSkip bool
	}{
		{"test.go", false},
		{"test.js", false},
		{"test.py", false},
		{"binary", true},
		{".git/config", true},
		{"node_modules/package/index.js", true},
		{"test.exe", true},
		{"test.so", true},
		{"image.png", true},
		{"document.pdf", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := engine.shouldSkipFile(tt.filename)
			if result != tt.shouldSkip {
				t.Errorf("Expected shouldSkip=%v for %s, got %v", tt.shouldSkip, tt.filename, result)
			}
		})
	}
}

func TestEngine_isSourceCodeFile(t *testing.T) {
	engine := NewEngine().(*Engine)

	tests := []struct {
		filename     string
		isSourceCode bool
	}{
		{"main.go", true},
		{"app.js", true},
		{"script.py", true},
		{"Component.tsx", true},
		{"styles.css", true},
		{"config.json", false},
		{"README.md", false},
		{"image.png", false},
		{"binary", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := engine.isSourceCodeFile(tt.filename)
			if result != tt.isSourceCode {
				t.Errorf("Expected isSourceCode=%v for %s, got %v", tt.isSourceCode, tt.filename, result)
			}
		})
	}
}

func TestEngine_hasLicenseFile(t *testing.T) {
	engine := NewEngine().(*Engine)

	// Test directory without license
	tempDir := t.TempDir()
	if engine.hasLicenseFile(tempDir) {
		t.Error("Expected no license file in empty directory")
	}

	// Test directory with LICENSE file
	license := "MIT License\n\nCopyright..."
	err := os.WriteFile(filepath.Join(tempDir, "LICENSE"), []byte(license), 0644)
	if err != nil {
		t.Fatalf("Failed to create LICENSE file: %v", err)
	}

	if !engine.hasLicenseFile(tempDir) {
		t.Error("Expected to find LICENSE file")
	}
}

// Benchmark tests
func BenchmarkEngine_AuditSecurity(b *testing.B) {
	engine := NewEngine()
	tempDir := b.TempDir()

	// Create test file
	testFile := `const config = { apiKey: "test-key" };`
	err := os.WriteFile(filepath.Join(tempDir, "config.js"), []byte(testFile), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.AuditSecurity(tempDir)
		if err != nil {
			b.Fatalf("AuditSecurity failed: %v", err)
		}
	}
}

func BenchmarkEngine_DetectSecrets(b *testing.B) {
	engine := NewEngine()
	tempDir := b.TempDir()

	// Create test file with secrets
	secretFile := `const config = { apiKey: "sk-1234567890abcdef", password: "secret123" };`
	err := os.WriteFile(filepath.Join(tempDir, "config.js"), []byte(secretFile), 0644)
	if err != nil {
		b.Fatalf("Failed to create secret file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.DetectSecrets(tempDir)
		if err != nil {
			b.Fatalf("DetectSecrets failed: %v", err)
		}
	}
}

func BenchmarkEngine_AnalyzeDependencies(b *testing.B) {
	engine := NewEngine()
	tempDir := b.TempDir()

	// Create package.json
	packageJSON := `{
		"name": "benchmark-project",
		"dependencies": {
			"react": "^18.0.0",
			"lodash": "^4.17.21"
		}
	}`

	err := os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		b.Fatalf("Failed to create package.json: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.AnalyzeDependencies(tempDir)
		if err != nil {
			b.Fatalf("AnalyzeDependencies failed: %v", err)
		}
	}
}
