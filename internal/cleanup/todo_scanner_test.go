package cleanup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTODOScanner_ScanProject(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "todo_scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files with various TODO patterns
	testFiles := map[string]string{
		"main.go": `package main

import "fmt"

// TODO: Add proper error handling
func main() {
	fmt.Println("Hello World")
	// FIXME: This is a security vulnerability
	password := "hardcoded"
	// HACK: Quick fix for performance issue
	for i := 0; i < 1000; i++ {
		// Do something
	}
}`,
		"security.go": `package security

// TODO: Implement proper authentication
func Authenticate() {
	// XXX: This needs immediate attention
}

// BUG: Memory leak in this function
func ProcessData() {
	// OPTIMIZE: This could be faster
}`,
		"docs/README.md": `# Project

TODO: Add installation instructions
FIXME: Update API documentation`,
		"vendor/external.go": `// TODO: This should be ignored`,
	}

	// Write test files
	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	// Create scanner
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	// Scan project
	report, err := scanner.ScanProject(tempDir)
	if err != nil {
		t.Fatalf("Failed to scan project: %v", err)
	}

	// Verify results
	if report.Summary.TotalTODOs == 0 {
		t.Error("Expected to find TODO comments, but found none")
	}

	// Should find TODOs in main.go, security.go, and README.md
	// Should NOT find TODOs in vendor/ directory
	expectedFiles := []string{"main.go", "security.go", "docs/README.md"}
	foundFiles := make(map[string]bool)

	for _, todo := range report.TODOs {
		// Get relative path from temp dir for comparison
		relPath, _ := filepath.Rel(tempDir, todo.File)
		foundFiles[relPath] = true

		// Verify vendor files are skipped
		if strings.Contains(todo.File, "vendor/") {
			t.Errorf("Found TODO in vendor directory, should be skipped: %s", todo.File)
		}
	}

	for _, expectedFile := range expectedFiles {
		if !foundFiles[expectedFile] {
			t.Errorf("Expected to find TODOs in %s, but didn't", expectedFile)
		}
	}

	// Verify categories are assigned correctly
	hasSecurityTODO := false
	hasPerformanceTODO := false
	for _, todo := range report.TODOs {
		if todo.Category == CategorySecurity {
			hasSecurityTODO = true
		}
		if todo.Category == CategoryPerformance {
			hasPerformanceTODO = true
		}
	}

	if !hasSecurityTODO {
		t.Error("Expected to find security-related TODO")
	}
	if !hasPerformanceTODO {
		t.Error("Expected to find performance-related TODO")
	}

	// Verify priorities are assigned correctly
	hasCriticalTODO := false
	for _, todo := range report.TODOs {
		if todo.Priority == PriorityCritical {
			hasCriticalTODO = true
			break
		}
	}

	if !hasCriticalTODO {
		t.Error("Expected to find critical priority TODO")
	}
}

func TestTODOScanner_DeterminePriority(t *testing.T) {
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	tests := []struct {
		todoType string
		message  string
		context  string
		expected Priority
	}{
		{"FIXME", "security vulnerability", "", PriorityCritical},
		{"TODO", "security issue here", "", PriorityCritical},
		{"BUG", "memory leak", "", PriorityCritical},
		{"HACK", "performance problem", "", PriorityHigh},
		{"TODO", "optimize this function", "", PriorityMedium},
		{"TODO", "add feature", "", PriorityLow},
		{"NOTE", "remember to update", "", PriorityLow},
	}

	for _, test := range tests {
		priority := scanner.determinePriority(test.todoType, test.message, test.context)
		if priority != test.expected {
			t.Errorf("For TODO type '%s' with message '%s', expected priority %v, got %v",
				test.todoType, test.message, test.expected, priority)
		}
	}
}

func TestTODOScanner_DetermineCategory(t *testing.T) {
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	tests := []struct {
		message  string
		context  string
		filePath string
		expected Category
	}{
		{"security vulnerability", "", "", CategorySecurity},
		{"authentication issue", "", "", CategorySecurity},
		{"", "// TODO: Fix security hole", "", CategorySecurity},
		{"", "", "pkg/security/auth.go", CategorySecurity},
		{"optimize performance", "", "", CategoryPerformance},
		{"memory usage", "", "", CategoryPerformance},
		{"add documentation", "", "", CategoryDocumentation},
		{"", "", "docs/README.md", CategoryDocumentation},
		{"fix bug", "", "", CategoryBug},
		{"refactor code", "", "", CategoryRefactor},
		{"add new feature", "", "", CategoryFeature},
	}

	for _, test := range tests {
		category := scanner.determineCategory(test.message, test.context, test.filePath)
		if category != test.expected {
			t.Errorf("For message '%s', context '%s', file '%s', expected category %v, got %v",
				test.message, test.context, test.filePath, test.expected, category)
		}
	}
}

func TestTODOScanner_GenerateMarkdownReport(t *testing.T) {
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	// Create test report
	report := &TODOReport{
		Timestamp:    time.Now(),
		ProjectRoot:  "/test/project",
		TotalFiles:   10,
		FilesScanned: 8,
		TODOs: []TODOItem{
			{
				File:     "main.go",
				Line:     10,
				Type:     "TODO",
				Message:  "Add error handling",
				Context:  "// TODO: Add error handling",
				Priority: PriorityMedium,
				Category: CategoryFeature,
			},
			{
				File:     "security.go",
				Line:     25,
				Type:     "FIXME",
				Message:  "Security vulnerability",
				Context:  "// FIXME: Security vulnerability",
				Priority: PriorityCritical,
				Category: CategorySecurity,
			},
		},
		Summary: &TODOSummary{
			TotalTODOs:    2,
			SecurityTODOs: 1,
			FeatureTODOs:  1,
			CriticalTODOs: 1,
			MediumTODOs:   1,
		},
		CategoryBreakdown: make(map[Category][]TODOItem),
		PriorityBreakdown: make(map[Priority][]TODOItem),
		FileBreakdown:     make(map[string][]TODOItem),
	}

	// Populate breakdowns
	for _, todo := range report.TODOs {
		report.CategoryBreakdown[todo.Category] = append(report.CategoryBreakdown[todo.Category], todo)
		report.PriorityBreakdown[todo.Priority] = append(report.PriorityBreakdown[todo.Priority], todo)
		report.FileBreakdown[todo.File] = append(report.FileBreakdown[todo.File], todo)
	}

	// Generate report
	reportContent, err := scanner.GenerateReport(report)
	if err != nil {
		t.Fatalf("Failed to generate report: %v", err)
	}

	// Verify report content
	expectedSections := []string{
		"# TODO/FIXME Analysis Report",
		"## Summary",
		"**Total TODOs:** 2",
		"**Critical:** 1",
		"**Security:** 1",
		"## Critical Priority TODOs",
		"### security.go:25 - FIXME",
		"**Message:** Security vulnerability",
	}

	for _, section := range expectedSections {
		if !strings.Contains(reportContent, section) {
			t.Errorf("Report should contain section: %s", section)
		}
	}
}

func TestTODOScanner_ShouldSkipFile(t *testing.T) {
	config := &TODOScanConfig{
		SkipPatterns: []string{"vendor/", ".git/", "node_modules/"},
	}
	scanner := NewTODOScanner(config)

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", false},
		{"pkg/security/auth.go", false},
		{"vendor/external/lib.go", true},
		{".git/config", true},
		{"node_modules/package/index.js", true},
		{"internal/app/app.go", false},
	}

	for _, test := range tests {
		result := scanner.shouldSkipFile(test.path)
		if result != test.expected {
			t.Errorf("For path '%s', expected skip=%v, got skip=%v", test.path, test.expected, result)
		}
	}
}

func TestTODOScanner_IsTextFile(t *testing.T) {
	config := &TODOScanConfig{
		IncludePatterns: []string{"*.go", "*.md", "*.yaml", "*.yml", "*.json"},
	}
	scanner := NewTODOScanner(config)

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"README.md", true},
		{"config.yaml", true},
		{"config.yml", true},
		{"package.json", true},
		{"binary.exe", false},
		{"image.png", false},
		{"data.txt", false},
	}

	for _, test := range tests {
		result := scanner.isTextFile(test.path)
		if result != test.expected {
			t.Errorf("For path '%s', expected isText=%v, got isText=%v", test.path, test.expected, result)
		}
	}
}

func TestTODOScanner_GenerateSummary(t *testing.T) {
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	todos := []TODOItem{
		{Category: CategorySecurity, Priority: PriorityCritical},
		{Category: CategorySecurity, Priority: PriorityHigh},
		{Category: CategoryPerformance, Priority: PriorityMedium},
		{Category: CategoryFeature, Priority: PriorityLow},
		{Category: CategoryBug, Priority: PriorityHigh},
	}

	summary := scanner.generateSummary(todos)

	if summary.TotalTODOs != 5 {
		t.Errorf("Expected 5 total TODOs, got %d", summary.TotalTODOs)
	}
	if summary.SecurityTODOs != 2 {
		t.Errorf("Expected 2 security TODOs, got %d", summary.SecurityTODOs)
	}
	if summary.CriticalTODOs != 1 {
		t.Errorf("Expected 1 critical TODO, got %d", summary.CriticalTODOs)
	}
	if summary.HighTODOs != 2 {
		t.Errorf("Expected 2 high priority TODOs, got %d", summary.HighTODOs)
	}
}
