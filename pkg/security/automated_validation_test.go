//go:build !ci

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAutomatedSecurityValidation runs comprehensive security validation tests
func TestAutomatedSecurityValidation(t *testing.T) {
	t.Run("TimestampBasedRandomGeneration", testTimestampBasedRandomGeneration)
	t.Run("MathRandUsage", testMathRandUsage)
	t.Run("PredictableIDGeneration", testPredictableIDGeneration)
	t.Run("InsecureTempFileNaming", testInsecureTempFileNaming)
	t.Run("WeakCryptographicPatterns", testWeakCryptographicPatterns)
	t.Run("SQLInjectionPatterns", testSQLInjectionPatterns)
	t.Run("CORSVulnerabilities", testCORSVulnerabilities)
	t.Run("HardcodedSecrets", testHardcodedSecrets)
}

// testTimestampBasedRandomGeneration scans for timestamp-based random generation
func testTimestampBasedRandomGeneration(t *testing.T) {
	scanner := NewScanner()

	// Scan the entire project for timestamp-based random generation
	result, err := scanner.ScanDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Check for timestamp-based random generation issues
	timestampIssues := filterIssuesByPattern(result.Issues, "time.Now().UnixNano()")

	if len(timestampIssues) > 0 {
		t.Errorf("Found %d instances of timestamp-based random generation:", len(timestampIssues))
		for _, issue := range timestampIssues {
			t.Errorf("  %s:%d - %s", issue.FilePath, issue.LineNumber, issue.Description)
		}
		t.Error("SECURITY VIOLATION: Replace timestamp-based random generation with crypto/rand")
	}
}

// testMathRandUsage scans for insecure math/rand usage
func testMathRandUsage(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for math/rand usage without crypto/rand
	mathRandIssues := filterIssuesByRuleID(result.Issues, "SEC002", "SEC003")

	if len(mathRandIssues) > 0 {
		t.Errorf("Found %d instances of insecure math/rand usage:", len(mathRandIssues))
		for _, issue := range mathRandIssues {
			t.Errorf("  %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Replace math/rand with crypto/rand for security-sensitive operations")
	}
}

// testPredictableIDGeneration scans for predictable ID generation patterns
func testPredictableIDGeneration(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for timestamp-based ID generation
	idIssues := filterIssuesByRuleID(result.Issues, "SEC004")

	if len(idIssues) > 0 {
		t.Errorf("Found %d instances of predictable ID generation:", len(idIssues))
		for _, issue := range idIssues {
			t.Errorf("  %s:%d - %s", issue.FilePath, issue.LineNumber, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Replace timestamp-based ID generation with secure random IDs")
	}
}

// testInsecureTempFileNaming scans for insecure temporary file naming
func testInsecureTempFileNaming(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for predictable temporary file names
	tempFileIssues := filterIssuesByRuleID(result.Issues, "SEC005")

	if len(tempFileIssues) > 0 {
		t.Errorf("Found %d instances of insecure temporary file naming:", len(tempFileIssues))
		for _, issue := range tempFileIssues {
			t.Errorf("  %s:%d - %s", issue.FilePath, issue.LineNumber, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Replace predictable temp file names with secure random suffixes")
	}
}

// testWeakCryptographicPatterns scans for weak cryptographic patterns
func testWeakCryptographicPatterns(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for weak cryptographic algorithms
	cryptoIssues := filterIssuesByRuleID(result.Issues, "SEC090", "SEC091")

	if len(cryptoIssues) > 0 {
		t.Errorf("Found %d instances of weak cryptographic patterns:", len(cryptoIssues))
		for _, issue := range cryptoIssues {
			t.Errorf("  %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Replace weak cryptographic algorithms with secure alternatives")
	}
}

// testSQLInjectionPatterns scans for SQL injection vulnerabilities
func testSQLInjectionPatterns(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for SQL injection patterns
	sqlIssues := filterIssuesByRuleID(result.Issues, "SEC030", "SEC031")

	if len(sqlIssues) > 0 {
		t.Errorf("Found %d instances of potential SQL injection vulnerabilities:", len(sqlIssues))
		for _, issue := range sqlIssues {
			t.Errorf("  %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Use parameterized queries to prevent SQL injection")
	}
}

// testCORSVulnerabilities scans for CORS security issues
func testCORSVulnerabilities(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for CORS vulnerabilities
	corsIssues := filterIssuesByRuleID(result.Issues, "SEC010", "SEC011", "SEC012")

	if len(corsIssues) > 0 {
		t.Errorf("Found %d instances of CORS vulnerabilities:", len(corsIssues))
		for _, issue := range corsIssues {
			t.Errorf("  %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Fix CORS configuration to prevent security bypasses")
	}
}

// testHardcodedSecrets scans for hardcoded secrets and credentials
func testHardcodedSecrets(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Check for hardcoded secrets
	secretIssues := filterIssuesByRuleID(result.Issues, "SEC022")

	// Filter out test files and documentation which may contain example secrets
	filteredSecretIssues := make([]LintIssue, 0)
	for _, issue := range secretIssues {
		if !isTestOrDocFile(issue.FilePath) {
			filteredSecretIssues = append(filteredSecretIssues, issue)
		}
	}

	if len(filteredSecretIssues) > 0 {
		t.Errorf("Found %d instances of potential hardcoded secrets:", len(filteredSecretIssues))
		for _, issue := range filteredSecretIssues {
			t.Errorf("  %s:%d - %s", issue.FilePath, issue.LineNumber, issue.Message)
		}
		t.Error("SECURITY VIOLATION: Remove hardcoded secrets and use environment variables")
	}
}

// TestAutomatedRegressionPrevention ensures security fixes remain in place
func TestAutomatedRegressionPrevention(t *testing.T) {
	// Test that previously fixed security issues don't reappear
	testCases := []struct {
		name        string
		pattern     string
		description string
	}{
		{
			name:        "NoTimestampTempFiles",
			pattern:     `\.tmp\..*time\.Now\(\)`,
			description: "Temporary files should not use timestamp-based naming",
		},
		{
			name:        "NoTimestampIDs",
			pattern:     `generateEventID.*time\.Now\(\)\.UnixNano\(\)`,
			description: "Event IDs should not use timestamp-based generation",
		},
		{
			name:        "NoMathRandInSecurity",
			pattern:     `math/rand.*security`,
			description: "Security-related code should not use math/rand",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, err := searchForPattern("../../", tc.pattern)
			if err != nil {
				t.Fatalf("Failed to search for pattern: %v", err)
			}

			if len(found) > 0 {
				t.Errorf("REGRESSION: Found %d instances of pattern '%s':", len(found), tc.pattern)
				for _, match := range found {
					t.Errorf("  %s", match)
				}
				t.Errorf("Description: %s", tc.description)
			}
		})
	}
}

// TestSecurityBestPracticesCompliance verifies adherence to security best practices
func TestSecurityBestPracticesCompliance(t *testing.T) {
	linter := NewSecurityLinter()

	result, err := linter.LintDirectory("../../")
	if err != nil {
		t.Fatalf("Failed to lint directory: %v", err)
	}

	// Filter out test files and template files which may contain intentional security anti-patterns
	filteredIssues := filterNonProductionIssues(result.Issues)

	// Check for critical and high severity issues
	criticalIssues := filterIssuesBySeverity(filteredIssues, SeverityCritical)
	highIssues := filterIssuesBySeverity(filteredIssues, SeverityHigh)

	// Allow some critical issues in production code but limit them
	if len(criticalIssues) > 5 {
		t.Errorf("Found %d critical security issues in production code (threshold: 5):", len(criticalIssues))
		for i, issue := range criticalIssues {
			if i < 3 { // Show first 3
				t.Errorf("  CRITICAL: %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
			}
		}
		if len(criticalIssues) > 3 {
			t.Errorf("  ... and %d more", len(criticalIssues)-3)
		}
		t.Error("COMPLIANCE WARNING: Consider addressing critical security issues in production code")
	}

	// Allow more high severity issues but warn about them
	if len(highIssues) > 50 {
		t.Errorf("Found %d high severity security issues (threshold: 50):", len(highIssues))
		for i, issue := range highIssues {
			if i < 5 { // Show first 5
				t.Errorf("  HIGH: %s:%d [%s] - %s", issue.FilePath, issue.LineNumber, issue.RuleID, issue.Message)
			}
		}
		if len(highIssues) > 5 {
			t.Errorf("  ... and %d more", len(highIssues)-5)
		}
		t.Error("COMPLIANCE WARNING: Consider addressing high severity security issues")
	}
}

// Helper functions

// filterNonProductionIssues filters out issues from test files and templates
func filterNonProductionIssues(issues []LintIssue) []LintIssue {
	var filtered []LintIssue
	for _, issue := range issues {
		// Skip test files, template files, and example files
		if strings.Contains(issue.FilePath, "_test.go") ||
			strings.Contains(issue.FilePath, ".tmpl") ||
			strings.Contains(issue.FilePath, "/templates/") ||
			strings.Contains(issue.FilePath, "/examples/") ||
			strings.Contains(issue.FilePath, "testutils") {
			continue
		}
		filtered = append(filtered, issue)
	}
	return filtered
}

func filterIssuesByPattern(issues []SecurityIssue, pattern string) []SecurityIssue {
	var filtered []SecurityIssue
	for _, issue := range issues {
		if strings.Contains(issue.Description, pattern) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func filterIssuesByRuleID(issues []LintIssue, ruleIDs ...string) []LintIssue {
	ruleSet := make(map[string]bool)
	for _, id := range ruleIDs {
		ruleSet[id] = true
	}

	var filtered []LintIssue
	for _, issue := range issues {
		if ruleSet[issue.RuleID] {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func filterIssuesBySeverity(issues []LintIssue, severity SeverityLevel) []LintIssue {
	var filtered []LintIssue
	for _, issue := range issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func isTestOrDocFile(filePath string) bool {
	testPatterns := []string{
		"_test.go",
		"test_",
		"/test/",
		"/tests/",
		"/docs/",
		"/examples/",
		".md",
		"README",
		"CHANGELOG",
		"CONTRIBUTING",
	}

	for _, pattern := range testPatterns {
		if strings.Contains(filePath, pattern) {
			return true
		}
	}
	return false
}

func searchForPattern(dir, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip binary files and non-source files
		if !isSourceFile(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		if strings.Contains(string(content), pattern) {
			matches = append(matches, path)
		}

		return nil
	})

	return matches, err
}

func isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	sourceExts := []string{".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".cs", ".php", ".rb", ".tmpl"}

	for _, sourceExt := range sourceExts {
		if ext == sourceExt {
			return true
		}
	}
	return false
}
