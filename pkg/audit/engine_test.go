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
		t.Fatal("NewEngine() returned nil")
	}

	// Check that default rules are loaded
	rules := engine.GetAuditRules()
	if len(rules) == 0 {
		t.Error("Expected default audit rules to be loaded")
	}

	// Verify some expected default rules exist
	expectedRules := []string{"security-001", "security-002", "quality-001", "license-001", "performance-001"}
	ruleMap := make(map[string]bool)
	for _, rule := range rules {
		ruleMap[rule.ID] = true
	}

	for _, expectedID := range expectedRules {
		if !ruleMap[expectedID] {
			t.Errorf("Expected default rule %s not found", expectedID)
		}
	}
}

func TestAuditRuleManagement(t *testing.T) {
	engine := NewEngine()

	// Test adding a rule
	newRule := interfaces.AuditRule{
		ID:          "test-001",
		Name:        "Test Rule",
		Description: "A test rule",
		Category:    "test",
		Type:        "test",
		Severity:    "low",
		Enabled:     true,
	}

	err := engine.AddAuditRule(newRule)
	if err != nil {
		t.Fatalf("Failed to add audit rule: %v", err)
	}

	// Verify rule was added
	rules := engine.GetAuditRules()
	found := false
	for _, rule := range rules {
		if rule.ID == "test-001" {
			found = true
			if rule.Name != "Test Rule" {
				t.Errorf("Expected rule name 'Test Rule', got '%s'", rule.Name)
			}
			break
		}
	}
	if !found {
		t.Error("Added rule not found in rule list")
	}

	// Test removing a rule
	err = engine.RemoveAuditRule("test-001")
	if err != nil {
		t.Fatalf("Failed to remove audit rule: %v", err)
	}

	// Verify rule was removed
	rules = engine.GetAuditRules()
	for _, rule := range rules {
		if rule.ID == "test-001" {
			t.Error("Rule should have been removed but still exists")
		}
	}

	// Test removing non-existent rule
	err = engine.RemoveAuditRule("non-existent")
	if err == nil {
		t.Error("Expected error when removing non-existent rule")
	}
}

func TestProjectExists(t *testing.T) {
	engine := &Engine{}

	// Test with current directory (should exist)
	err := engine.projectExists(".")
	if err != nil {
		t.Errorf("Current directory should exist: %v", err)
	}

	// Test with non-existent directory
	err = engine.projectExists("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Test with a file instead of directory
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	_ = tempFile.Close()

	err = engine.projectExists(tempFile.Name())
	if err == nil {
		t.Error("Expected error when path is a file, not directory")
	}
}

func TestAuditProject(t *testing.T) {
	engine := NewEngine()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "audit_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create some test files
	testFiles := map[string]string{
		"README.md":    "# Test Project\nThis is a test project.",
		"package.json": `{"name": "test", "version": "1.0.0", "description": "Test package"}`,
		"main.js":      `console.log("Hello, world!");`,
		"test.js":      `describe("test", () => { it("should work", () => {}); });`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test audit with default options
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
		t.Fatal("AuditProject returned nil result")
	}

	if result.ProjectPath != tempDir {
		t.Errorf("Expected project path %s, got %s", tempDir, result.ProjectPath)
	}

	if result.AuditTime.IsZero() {
		t.Error("AuditTime should be set")
	}

	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("OverallScore should be between 0 and 100, got %f", result.OverallScore)
	}

	// Verify that audit sections are present when requested
	if result.Security == nil {
		t.Error("Security audit result should be present")
	}

	if result.Quality == nil {
		t.Error("Quality audit result should be present")
	}

	if result.Licenses == nil {
		t.Error("License audit result should be present")
	}

	if result.Performance == nil {
		t.Error("Performance audit result should be present")
	}
}

func TestGenerateAuditReport(t *testing.T) {
	engine := NewEngine()

	// Create a sample audit result
	result := &interfaces.AuditResult{
		ProjectPath:  "/test/project",
		AuditTime:    time.Now(),
		OverallScore: 85.5,
		Security: &interfaces.SecurityAuditResult{
			Score:           90.0,
			Vulnerabilities: []interfaces.Vulnerability{},
		},
		Quality: &interfaces.QualityAuditResult{
			Score:        80.0,
			TestCoverage: 75.0,
		},
		Recommendations: []string{"Improve test coverage", "Update dependencies"},
	}

	// Test JSON report generation
	jsonReport, err := engine.GenerateAuditReport(result, "json")
	if err != nil {
		t.Fatalf("Failed to generate JSON report: %v", err)
	}
	if len(jsonReport) == 0 {
		t.Error("JSON report should not be empty")
	}

	// Test HTML report generation
	htmlReport, err := engine.GenerateAuditReport(result, "html")
	if err != nil {
		t.Fatalf("Failed to generate HTML report: %v", err)
	}
	if len(htmlReport) == 0 {
		t.Error("HTML report should not be empty")
	}

	// Test Markdown report generation
	markdownReport, err := engine.GenerateAuditReport(result, "markdown")
	if err != nil {
		t.Fatalf("Failed to generate Markdown report: %v", err)
	}
	if len(markdownReport) == 0 {
		t.Error("Markdown report should not be empty")
	}

	// Test unsupported format
	_, err = engine.GenerateAuditReport(result, "unsupported")
	if err == nil {
		t.Error("Expected error for unsupported report format")
	}
}

func TestShouldSkipFile(t *testing.T) {
	engine := &Engine{}

	testCases := []struct {
		filename string
		expected bool
	}{
		{"main.go", false},
		{"test.js", false},
		{"README.md", false},
		{"image.jpg", true},
		{"binary.exe", true},
		{".hidden", true},
		{"node_modules/package/index.js", true},
		{".git/config", true},
		{"build/output.js", true},
	}

	for _, tc := range testCases {
		result := engine.shouldSkipFile(tc.filename)
		if result != tc.expected {
			t.Errorf("shouldSkipFile(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestIsSourceCodeFile(t *testing.T) {
	engine := &Engine{}

	testCases := []struct {
		filename string
		expected bool
	}{
		{"main.go", true},
		{"app.js", true},
		{"component.tsx", true},
		{"style.css", false},
		{"README.md", false},
		{"image.png", false},
		{"data.json", false},
		{"script.py", true},
		{"Main.java", true},
	}

	for _, tc := range testCases {
		result := engine.isSourceCodeFile(tc.filename)
		if result != tc.expected {
			t.Errorf("isSourceCodeFile(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestIsTestFile(t *testing.T) {
	engine := &Engine{}

	testCases := []struct {
		filename string
		expected bool
	}{
		{"main_test.go", true},
		{"app.test.js", true},
		{"component.spec.tsx", true},
		{"test_helper.py", true},
		{"main.go", false},
		{"app.js", false},
		{"README.md", false},
		{"project/tests/integration.js", true},
		{"project/spec/unit.rb", true},
	}

	for _, tc := range testCases {
		result := engine.isTestFile(tc.filename)
		if result != tc.expected {
			t.Errorf("isTestFile(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestGetAuditSummary(t *testing.T) {
	engine := NewEngine()

	// Create sample audit results
	results := []*interfaces.AuditResult{
		{
			OverallScore: 85.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 90.0,
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "high"},
					{Severity: "medium"},
				},
			},
			Quality: &interfaces.QualityAuditResult{
				Score: 80.0,
			},
		},
		{
			OverallScore: 75.0,
			Security: &interfaces.SecurityAuditResult{
				Score: 70.0,
				Vulnerabilities: []interfaces.Vulnerability{
					{Severity: "critical"},
				},
			},
			Quality: &interfaces.QualityAuditResult{
				Score: 80.0,
			},
		},
	}

	summary, err := engine.GetAuditSummary(results)
	if err != nil {
		t.Fatalf("GetAuditSummary failed: %v", err)
	}

	if summary.TotalProjects != 2 {
		t.Errorf("Expected 2 total projects, got %d", summary.TotalProjects)
	}

	expectedAverage := (85.0 + 75.0) / 2
	if summary.AverageScore != expectedAverage {
		t.Errorf("Expected average score %f, got %f", expectedAverage, summary.AverageScore)
	}

	expectedSecurityAverage := (90.0 + 70.0) / 2
	if summary.SecurityScore != expectedSecurityAverage {
		t.Errorf("Expected security score %f, got %f", expectedSecurityAverage, summary.SecurityScore)
	}

	// Test with empty results
	emptySummary, err := engine.GetAuditSummary([]*interfaces.AuditResult{})
	if err != nil {
		t.Fatalf("GetAuditSummary with empty results failed: %v", err)
	}

	if emptySummary.TotalProjects != 0 {
		t.Errorf("Expected 0 total projects for empty results, got %d", emptySummary.TotalProjects)
	}
}
