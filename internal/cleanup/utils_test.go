package cleanup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCleanupUtils_ScanProjectFiles(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "utils_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"main.go":          "package main",
		"internal/app.go":  "package internal",
		"pkg/utils.go":     "package pkg",
		"vendor/lib.go":    "package vendor", // Should be skipped
		"test.txt":         "not a go file",  // Should be skipped
		"cmd/tool/main.go": "package main",
	}

	for filename, content := range testFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", filename, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	utils := NewCleanupUtils()
	skipPatterns := []string{"vendor/"}

	files, err := utils.ScanProjectFiles(tempDir, skipPatterns)
	if err != nil {
		t.Fatalf("Failed to scan project files: %v", err)
	}

	// Should find 4 Go files (excluding vendor and non-Go files)
	expectedCount := 4
	if len(files) != expectedCount {
		t.Errorf("Expected %d Go files, got %d", expectedCount, len(files))
	}

	// Check that vendor files are excluded
	for _, file := range files {
		if strings.Contains(file, "vendor/") {
			t.Errorf("Vendor file should be excluded: %s", file)
		}
		if !strings.HasSuffix(file, ".go") {
			t.Errorf("Non-Go file should be excluded: %s", file)
		}
	}
}

func TestCleanupUtils_GenerateCleanupReport(t *testing.T) {
	utils := NewCleanupUtils()

	// Create test analysis
	analysis := &ProjectAnalysis{
		ProjectRoot: "/test/project",
		Timestamp:   time.Now(),
		TODOs: []TODOItem{
			{
				File:     "main.go",
				Line:     10,
				Type:     "TODO",
				Message:  "Implement feature",
				Priority: PriorityMedium,
				Category: CategoryFeature,
			},
			{
				File:     "security.go",
				Line:     25,
				Type:     "FIXME",
				Message:  "Security vulnerability",
				Priority: PriorityHigh,
				Category: CategorySecurity,
			},
		},
		Duplicates: []DuplicateCodeBlock{
			{
				Files:      []string{"file1.go", "file2.go"},
				Similarity: 0.85,
				Suggestion: "Extract to common function",
			},
		},
		UnusedCode: []UnusedCodeItem{
			{
				File:   "utils.go",
				Line:   50,
				Type:   "function",
				Name:   "UnusedFunc",
				Reason: "Never called",
			},
		},
		ImportIssues: []ImportIssue{
			{
				File:       "main.go",
				Line:       5,
				Type:       "grouping",
				Import:     "fmt",
				Suggestion: "Group standard library imports",
			},
		},
	}

	// Create test result
	result := &CleanupResult{
		Success:       true,
		FilesModified: []string{"main.go", "security.go"},
		IssuesFixed: []FixedIssue{
			{
				Type:        "todo",
				File:        "main.go",
				Line:        10,
				Description: "Implemented feature",
				Action:      "Added implementation",
			},
		},
		IssuesRemaining: []RemainingIssue{
			{
				Type:        "security",
				File:        "security.go",
				Line:        25,
				Description: "Security issue needs manual review",
				Reason:      "Complex security fix required",
				Suggestion:  "Review with security team",
			},
		},
		Duration: 5 * time.Minute,
	}

	report := utils.GenerateCleanupReport(analysis, result)

	// Verify report contains expected sections
	expectedSections := []string{
		"# Cleanup Report",
		"## Summary",
		"## Analysis Results",
		"### TODO/FIXME Comments",
		"#### Security",
		"#### Feature",
		"### Duplicate Code Blocks",
		"### Unused Code Items",
		"### Import Organization Issues",
		"## Recommendations",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing expected section: %s", section)
		}
	}

	// Verify specific content
	if !strings.Contains(report, "Files Modified: 2") {
		t.Error("Report should contain files modified count")
	}

	if !strings.Contains(report, "Security vulnerability") {
		t.Error("Report should contain security TODO message")
	}

	if !strings.Contains(report, "Extract to common function") {
		t.Error("Report should contain duplicate code suggestion")
	}
}

func TestCleanupUtils_ValidateCleanupResult(t *testing.T) {
	utils := NewCleanupUtils()

	tests := []struct {
		name           string
		result         *CleanupResult
		expectedIssues int
	}{
		{
			name:           "nil result",
			result:         nil,
			expectedIssues: 1,
		},
		{
			name: "valid result",
			result: &CleanupResult{
				Success:       true,
				FilesModified: []string{"main.go"},
				IssuesFixed: []FixedIssue{
					{Type: "todo", Description: "Fixed TODO"},
				},
				Duration: 2 * time.Minute,
			},
			expectedIssues: 0,
		},
		{
			name: "inconsistent success status",
			result: &CleanupResult{
				Success: false,
				IssuesFixed: []FixedIssue{
					{Type: "todo", Description: "Fixed TODO"},
				},
				Duration: 2 * time.Minute,
			},
			expectedIssues: 2, // Both inconsistent status and no files modified
		},
		{
			name: "negative duration",
			result: &CleanupResult{
				Success:  true,
				Duration: -1 * time.Minute,
			},
			expectedIssues: 1,
		},
		{
			name: "unusually long duration",
			result: &CleanupResult{
				Success:  true,
				Duration: 25 * time.Hour,
			},
			expectedIssues: 1,
		},
		{
			name: "issues fixed but no files modified",
			result: &CleanupResult{
				Success:       true,
				FilesModified: []string{},
				IssuesFixed: []FixedIssue{
					{Type: "todo", Description: "Fixed TODO"},
				},
				Duration: 2 * time.Minute,
			},
			expectedIssues: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			issues := utils.ValidateCleanupResult(test.result)
			if len(issues) != test.expectedIssues {
				t.Errorf("Expected %d issues, got %d: %v",
					test.expectedIssues, len(issues), issues)
			}
		})
	}
}

func TestCleanupUtils_EstimateCleanupTime(t *testing.T) {
	utils := NewCleanupUtils()

	tests := []struct {
		name     string
		analysis *ProjectAnalysis
		minTime  time.Duration
		maxTime  time.Duration
	}{
		{
			name:     "nil analysis",
			analysis: nil,
			minTime:  0,
			maxTime:  0,
		},
		{
			name: "empty analysis",
			analysis: &ProjectAnalysis{
				TODOs:        []TODOItem{},
				Duplicates:   []DuplicateCodeBlock{},
				UnusedCode:   []UnusedCodeItem{},
				ImportIssues: []ImportIssue{},
			},
			minTime: 30 * time.Second, // Base time
			maxTime: 45 * time.Second, // Base time + buffer
		},
		{
			name: "analysis with critical TODOs",
			analysis: &ProjectAnalysis{
				TODOs: []TODOItem{
					{Priority: PriorityCritical},
					{Priority: PriorityHigh},
					{Priority: PriorityMedium},
					{Priority: PriorityLow},
				},
				Duplicates:   []DuplicateCodeBlock{{}},
				UnusedCode:   []UnusedCodeItem{{}, {}},
				ImportIssues: []ImportIssue{{}, {}, {}},
			},
			minTime: 8 * time.Minute,  // Should be substantial
			maxTime: 16 * time.Minute, // With buffer (increased to account for 25% buffer)
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			estimate := utils.EstimateCleanupTime(test.analysis)

			if estimate < test.minTime {
				t.Errorf("Estimate %v is less than minimum expected %v",
					estimate, test.minTime)
			}

			if estimate > test.maxTime {
				t.Errorf("Estimate %v is greater than maximum expected %v",
					estimate, test.maxTime)
			}
		})
	}
}

func TestCleanupUtils_CreateCleanupPlan(t *testing.T) {
	utils := NewCleanupUtils()

	analysis := &ProjectAnalysis{
		ProjectRoot: "/test/project",
		TODOs: []TODOItem{
			{
				Type:     "FIXME",
				Message:  "Security issue",
				Priority: PriorityHigh,
				Category: CategorySecurity,
			},
			{
				Type:     "TODO",
				Message:  "Add feature",
				Priority: PriorityMedium,
				Category: CategoryFeature,
			},
		},
		Duplicates: []DuplicateCodeBlock{
			{
				Files:      []string{"file1.go", "file2.go"},
				Suggestion: "Extract common function",
			},
		},
		UnusedCode: []UnusedCodeItem{
			{
				Type: "function",
				Name: "UnusedFunc",
			},
		},
		ImportIssues: []ImportIssue{
			{
				Type:       "grouping",
				Suggestion: "Group imports properly",
			},
		},
	}

	plan := utils.CreateCleanupPlan(analysis)

	// Verify plan structure
	if plan == nil {
		t.Fatal("Plan should not be nil")
	}

	if plan.ProjectRoot != analysis.ProjectRoot {
		t.Errorf("Expected project root %s, got %s",
			analysis.ProjectRoot, plan.ProjectRoot)
	}

	// Should have 5 phases (security, unused, imports, remaining TODOs, duplicates)
	expectedPhases := 5
	if len(plan.Phases) != expectedPhases {
		t.Errorf("Expected %d phases, got %d", expectedPhases, len(plan.Phases))
	}

	// Verify phase priorities are in order
	for i := 1; i < len(plan.Phases); i++ {
		if plan.Phases[i].Priority <= plan.Phases[i-1].Priority {
			t.Error("Phases should be ordered by priority")
		}
	}

	// Verify total task count
	totalTasks := plan.GetTaskCount()
	expectedTasks := 5 // 2 TODOs + 1 duplicate + 1 unused + 1 import
	if totalTasks != expectedTasks {
		t.Errorf("Expected %d total tasks, got %d", expectedTasks, totalTasks)
	}

	// Verify estimated time is reasonable
	totalTime := plan.GetTotalEstimatedTime()
	if totalTime <= 0 {
		t.Error("Total estimated time should be positive")
	}

	// First phase should be security
	if len(plan.Phases) > 0 && plan.Phases[0].Name != "Security Issues" {
		t.Error("First phase should be security issues")
	}
}

func TestCleanupUtils_CategoryAndPriorityNames(t *testing.T) {
	utils := NewCleanupUtils()

	// Test category names
	categoryTests := []struct {
		category Category
		expected string
	}{
		{CategorySecurity, "Security"},
		{CategoryPerformance, "Performance"},
		{CategoryFeature, "Feature"},
		{CategoryBug, "Bug"},
		{CategoryDocumentation, "Documentation"},
		{CategoryRefactor, "Refactor"},
		{Category(999), "Other"}, // Unknown category
	}

	for _, test := range categoryTests {
		name := utils.getCategoryName(test.category)
		if name != test.expected {
			t.Errorf("Expected category name %s, got %s", test.expected, name)
		}
	}

	// Test priority names
	priorityTests := []struct {
		priority Priority
		expected string
	}{
		{PriorityCritical, "Critical"},
		{PriorityHigh, "High"},
		{PriorityMedium, "Medium"},
		{PriorityLow, "Low"},
		{Priority(999), "Unknown"}, // Unknown priority
	}

	for _, test := range priorityTests {
		name := utils.getPriorityName(test.priority)
		if name != test.expected {
			t.Errorf("Expected priority name %s, got %s", test.expected, name)
		}
	}
}

func TestCleanupUtils_FilterTodos(t *testing.T) {
	utils := NewCleanupUtils()

	todos := []TODOItem{
		{Category: CategorySecurity, Message: "Security issue"},
		{Category: CategoryFeature, Message: "Feature request"},
		{Category: CategorySecurity, Message: "Another security issue"},
		{Category: CategoryBug, Message: "Bug fix"},
	}

	// Test filtering by category
	securityTodos := utils.filterTodosByCategory(todos, CategorySecurity)
	if len(securityTodos) != 2 {
		t.Errorf("Expected 2 security TODOs, got %d", len(securityTodos))
	}

	// Test filtering excluding category
	nonSecurityTodos := utils.filterTodosExcludeCategory(todos, CategorySecurity)
	if len(nonSecurityTodos) != 2 {
		t.Errorf("Expected 2 non-security TODOs, got %d", len(nonSecurityTodos))
	}

	// Verify no overlap
	if len(securityTodos)+len(nonSecurityTodos) != len(todos) {
		t.Error("Filtered TODOs should account for all original TODOs")
	}
}
