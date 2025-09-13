package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSecurityLinter(t *testing.T) {
	t.Run("NewSecurityLinter", testNewSecurityLinter)
	t.Run("LintFile", testLintFile)
	t.Run("LintDirectory", testLintDirectory)
	t.Run("ExportResults", testExportResults)
	t.Run("FilterResults", testFilterResults)
}

func testNewSecurityLinter(t *testing.T) {
	linter := NewSecurityLinter()

	if linter == nil {
		t.Fatal("NewSecurityLinter returned nil")
	}

	if linter.scanner == nil {
		t.Error("Scanner not initialized")
	}

	if linter.fixer == nil {
		t.Error("Fixer not initialized")
	}

	if len(linter.rules) == 0 {
		t.Error("No linting rules loaded")
	}

	// Verify specific rules are present
	ruleIDs := make(map[string]bool)
	for _, rule := range linter.rules {
		ruleIDs[rule.ID] = true
	}

	expectedRules := []string{"SEC001", "SEC002", "SEC003", "SEC004", "SEC005"}
	for _, expectedRule := range expectedRules {
		if !ruleIDs[expectedRule] {
			t.Errorf("Expected rule %s not found", expectedRule)
		}
	}
}

func testLintFile(t *testing.T) {
	linter := NewSecurityLinter()

	// Create a temporary file with security issues
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	testContent := `package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	// SEC001: Timestamp-based random generation
	id := fmt.Sprintf("id_%d", time.Now().UnixNano())
	
	// SEC002: Math/rand usage
	randomNum := rand.Intn(100)
	
	// SEC004: Timestamp-based ID generation
	eventID := fmt.Sprintf("event_%d", time.Now().Unix())
	
	fmt.Println(id, randomNum, eventID)
}
`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Lint the file
	issues, err := linter.LintFile(testFile)
	if err != nil {
		t.Fatalf("LintFile failed: %v", err)
	}

	// Verify issues were found
	if len(issues) == 0 {
		t.Error("Expected security issues to be found, but none were detected")
	}

	// Check for specific rule violations
	foundRules := make(map[string]bool)
	for _, issue := range issues {
		foundRules[issue.RuleID] = true

		// Verify issue structure
		if issue.FilePath != testFile {
			t.Errorf("Expected file path %s, got %s", testFile, issue.FilePath)
		}

		if issue.LineNumber <= 0 {
			t.Errorf("Invalid line number: %d", issue.LineNumber)
		}

		if issue.Message == "" {
			t.Error("Issue message is empty")
		}

		if issue.Suggestion == "" {
			t.Error("Issue suggestion is empty")
		}
	}

	// Verify expected rules were triggered
	expectedRules := []string{"SEC001", "SEC003", "SEC004"}
	for _, expectedRule := range expectedRules {
		if !foundRules[expectedRule] {
			t.Errorf("Expected rule %s to be triggered", expectedRule)
		}
	}
}

func testLintDirectory(t *testing.T) {
	linter := NewSecurityLinter()

	// Create a temporary directory with multiple files
	tempDir := t.TempDir()

	// Create test files with different security issues
	testFiles := map[string]string{
		"secure.go": `package main
import "crypto/rand"
func main() {
	// This is secure code
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
}`,
		"insecure.go": `package main
import "time"
func main() {
	// SEC001: Timestamp-based random
	id := time.Now().UnixNano()
}`,
		"cors.js": `// SEC010: CORS null origin
res.setHeader('Access-Control-Allow-Origin', 'null');`,
		"sql.go": `package main
func query(userID string) {
	// SEC030: SQL concatenation
	sql := "SELECT * FROM users WHERE id = '" + userID + "'"
}`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Lint the directory
	result, err := linter.LintDirectory(tempDir)
	if err != nil {
		t.Fatalf("LintDirectory failed: %v", err)
	}

	// Verify results
	if result.ScannedFiles != 4 {
		t.Errorf("Expected 4 scanned files, got %d", result.ScannedFiles)
	}

	if len(result.Issues) == 0 {
		t.Error("Expected security issues to be found")
	}

	// Verify summary is populated
	if result.Summary.TotalIssues != len(result.Issues) {
		t.Errorf("Summary total issues (%d) doesn't match actual issues (%d)",
			result.Summary.TotalIssues, len(result.Issues))
	}

	// Check that issues are from different files
	fileIssues := make(map[string]int)
	for _, issue := range result.Issues {
		filename := filepath.Base(issue.FilePath)
		fileIssues[filename]++
	}

	if len(fileIssues) < 2 {
		t.Error("Expected issues from multiple files")
	}
}

func testExportResults(t *testing.T) {
	linter := NewSecurityLinter()

	// Create sample results
	result := &LintResult{
		Issues: []LintIssue{
			{
				RuleID:      "SEC001",
				FilePath:    "test.go",
				LineNumber:  10,
				Column:      5,
				Severity:    SeverityCritical,
				Category:    "cryptography",
				Message:     "Test message",
				Suggestion:  "Test suggestion",
				LineContent: "test content",
				Tags:        []string{"test"},
			},
		},
		Summary: LintSummary{
			TotalIssues: 1,
			BySeverity:  map[SeverityLevel]int{SeverityCritical: 1},
			ByCategory:  map[string]int{"cryptography": 1},
		},
		ScannedFiles: 1,
		RulesApplied: 10,
	}

	tempDir := t.TempDir()

	// Test JSON export
	jsonFile := filepath.Join(tempDir, "test.json")
	err := linter.ExportResults(result, "json", jsonFile)
	if err != nil {
		t.Errorf("JSON export failed: %v", err)
	}

	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("JSON file was not created")
	}

	// Test SARIF export
	sarifFile := filepath.Join(tempDir, "test.sarif")
	err = linter.ExportResults(result, "sarif", sarifFile)
	if err != nil {
		t.Errorf("SARIF export failed: %v", err)
	}

	if _, err := os.Stat(sarifFile); os.IsNotExist(err) {
		t.Error("SARIF file was not created")
	}

	// Test JUnit export
	junitFile := filepath.Join(tempDir, "test.xml")
	err = linter.ExportResults(result, "junit", junitFile)
	if err != nil {
		t.Errorf("JUnit export failed: %v", err)
	}

	if _, err := os.Stat(junitFile); os.IsNotExist(err) {
		t.Error("JUnit file was not created")
	}

	// Test invalid format
	err = linter.ExportResults(result, "invalid", "test.txt")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func testFilterResults(t *testing.T) {
	// Create sample results with various severities and categories
	result := &LintResult{
		Issues: []LintIssue{
			{RuleID: "SEC001", Severity: SeverityCritical, Category: "cryptography"},
			{RuleID: "SEC002", Severity: SeverityHigh, Category: "cryptography"},
			{RuleID: "SEC010", Severity: SeverityMedium, Category: "cors"},
			{RuleID: "SEC050", Severity: SeverityLow, Category: "security-headers"},
		},
		Summary: LintSummary{TotalIssues: 4},
	}

	// Test severity filtering
	highSeverityResult := filterResultsBySeverity(result, SeverityHigh)
	if len(highSeverityResult.Issues) != 2 { // Critical and High
		t.Errorf("Expected 2 high+ severity issues, got %d", len(highSeverityResult.Issues))
	}

	// Test category filtering
	cryptoResult := filterResultsByCategory(result, []string{"cryptography"})
	if len(cryptoResult.Issues) != 2 {
		t.Errorf("Expected 2 cryptography issues, got %d", len(cryptoResult.Issues))
	}
}

// Helper functions for testing

func filterResultsBySeverity(result *LintResult, minSeverity SeverityLevel) *LintResult {
	severityOrder := map[SeverityLevel]int{
		SeverityLow:      1,
		SeverityMedium:   2,
		SeverityHigh:     3,
		SeverityCritical: 4,
	}

	minLevel := severityOrder[minSeverity]
	filtered := &LintResult{
		Issues:       make([]LintIssue, 0),
		Summary:      result.Summary,
		ScannedFiles: result.ScannedFiles,
		RulesApplied: result.RulesApplied,
	}

	for _, issue := range result.Issues {
		if severityOrder[issue.Severity] >= minLevel {
			filtered.Issues = append(filtered.Issues, issue)
		}
	}

	return filtered
}

func filterResultsByCategory(result *LintResult, categories []string) *LintResult {
	categorySet := make(map[string]bool)
	for _, category := range categories {
		categorySet[category] = true
	}

	filtered := &LintResult{
		Issues:       make([]LintIssue, 0),
		Summary:      result.Summary,
		ScannedFiles: result.ScannedFiles,
		RulesApplied: result.RulesApplied,
	}

	for _, issue := range result.Issues {
		if categorySet[issue.Category] {
			filtered.Issues = append(filtered.Issues, issue)
		}
	}

	return filtered
}

// TestLintRules tests individual linting rules
func TestLintRules(t *testing.T) {
	rules := getSecurityLintRules()

	testCases := []struct {
		ruleID      string
		testCode    string
		shouldMatch bool
	}{
		{
			ruleID:      "SEC001",
			testCode:    "id := time.Now().UnixNano()",
			shouldMatch: true,
		},
		{
			ruleID:      "SEC001",
			testCode:    "id := secureRandom.Generate()",
			shouldMatch: false,
		},
		{
			ruleID:      "SEC002",
			testCode:    `import "math/rand"`,
			shouldMatch: true,
		},
		{
			ruleID:      "SEC002",
			testCode:    `import "crypto/rand"`,
			shouldMatch: false,
		},
		{
			ruleID:      "SEC003",
			testCode:    "num := rand.Intn(100)",
			shouldMatch: true,
		},
		{
			ruleID:      "SEC003",
			testCode:    "rand.Read(bytes)",
			shouldMatch: false,
		},
		{
			ruleID:      "SEC010",
			testCode:    `Access-Control-Allow-Origin: 'null'`,
			shouldMatch: true,
		},
		{
			ruleID:      "SEC030",
			testCode:    `sql := "SELECT * FROM users WHERE id = '" + userID + "'"`,
			shouldMatch: true,
		},
	}

	// Create rule lookup map
	ruleMap := make(map[string]LintRule)
	for _, rule := range rules {
		ruleMap[rule.ID] = rule
	}

	for _, tc := range testCases {
		t.Run(tc.ruleID+"_"+strings.ReplaceAll(tc.testCode, " ", "_"), func(t *testing.T) {
			rule, exists := ruleMap[tc.ruleID]
			if !exists {
				t.Fatalf("Rule %s not found", tc.ruleID)
			}

			matches := rule.Pattern.MatchString(tc.testCode)
			if matches != tc.shouldMatch {
				t.Errorf("Rule %s pattern match failed. Code: %s, Expected: %v, Got: %v",
					tc.ruleID, tc.testCode, tc.shouldMatch, matches)
			}
		})
	}
}

// TestLinterPerformance tests the performance of the linter
func TestLinterPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	linter := NewSecurityLinter()

	// Create a large test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large_test.go")

	var content strings.Builder
	content.WriteString("package main\n\n")

	// Generate 1000 lines of code with some security issues
	for i := 0; i < 1000; i++ {
		if i%100 == 0 {
			content.WriteString(fmt.Sprintf("// Line %d with security issue\n", i))
			content.WriteString("id := time.Now().UnixNano()\n")
		} else {
			content.WriteString(fmt.Sprintf("// Regular line %d\n", i))
			content.WriteString("fmt.Println(\"hello\")\n")
		}
	}

	err := os.WriteFile(testFile, []byte(content.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Measure linting performance
	start := time.Now()
	issues, err := linter.LintFile(testFile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("LintFile failed: %v", err)
	}

	t.Logf("Linted %d lines in %v, found %d issues", 1000*2, duration, len(issues))

	// Performance should be reasonable (less than 1 second for 2000 lines)
	if duration > time.Second {
		t.Errorf("Linting took too long: %v", duration)
	}

	// Should find the expected number of issues (10 security issues)
	if len(issues) < 10 {
		t.Errorf("Expected at least 10 issues, found %d", len(issues))
	}
}
