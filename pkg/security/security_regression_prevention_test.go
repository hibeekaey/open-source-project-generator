//go:build !ci

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSecurityRegressionPreventionSuite implements comprehensive regression tests to prevent future vulnerabilities
func TestSecurityRegressionPreventionSuite(t *testing.T) {
	t.Run("KnownVulnerabilityPatterns", func(t *testing.T) {
		testKnownVulnerabilityPatterns(t)
	})

	t.Run("SecurityFixRegression", func(t *testing.T) {
		testSecurityFixRegression(t)
	})

	t.Run("SafeCodePreservation", func(t *testing.T) {
		testSafeCodePreservation(t)
	})

	t.Run("SecurityMetricsConsistency", func(t *testing.T) {
		testSecurityMetricsConsistency(t)
	})

	t.Run("FixIdempotency", func(t *testing.T) {
		testFixIdempotency(t)
	})
}

// testKnownVulnerabilityPatterns ensures all known vulnerability patterns are still detected
func testKnownVulnerabilityPatterns(t *testing.T) {
	// Define critical vulnerability patterns that must always be detected
	criticalVulnerabilities := []struct {
		name         string
		code         string
		expectedType SecurityIssueType
		severity     SeverityLevel
		description  string
	}{
		{
			name:         "CORS null origin vulnerability",
			code:         `c.Header("Access-Control-Allow-Origin", "null")`,
			expectedType: CORSVulnerability,
			severity:     SeverityCritical,
			description:  "CORS null origin must always be detected as critical",
		},
		{
			name:         "JWT none algorithm vulnerability",
			code:         `token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)`,
			expectedType: WeakAuthentication,
			severity:     SeverityCritical,
			description:  "JWT none algorithm must always be detected as critical",
		},
		{
			name:         "SQL injection via string concatenation",
			code:         `query := "SELECT * FROM users WHERE id = " + userID`,
			expectedType: SQLInjectionRisk,
			severity:     SeverityCritical,
			description:  "SQL string concatenation must always be detected as critical",
		},
		{
			name:         "SQL injection via format string",
			code:         `query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)`,
			expectedType: SQLInjectionRisk,
			severity:     SeverityHigh,
			description:  "SQL format string must always be detected as high severity",
		},
		{
			name:         "Weak JWT secret",
			code:         `jwt.sign(payload, "secret")`,
			expectedType: WeakAuthentication,
			severity:     SeverityHigh,
			description:  "Weak JWT secrets must always be detected as high severity",
		},
	}

	scanner := NewScanner()

	for _, vuln := range criticalVulnerabilities {
		t.Run(vuln.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "vulnerability_test.go.tmpl")

			err := os.WriteFile(testFile, []byte(vuln.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Vulnerability scan failed: %v", err)
			}

			// Must detect the vulnerability
			if len(issues) == 0 {
				t.Errorf("REGRESSION: Known vulnerability not detected: %s", vuln.description)
				return
			}

			// Verify correct classification
			found := false
			for _, issue := range issues {
				if issue.IssueType == vuln.expectedType && issue.Severity == vuln.severity {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("REGRESSION: Vulnerability detected but with wrong classification. Expected: %s/%s", vuln.expectedType, vuln.severity)
				for _, issue := range issues {
					t.Logf("Found: %s/%s - %s", issue.IssueType, issue.Severity, issue.Description)
				}
			}
		})
	}
}

// testSecurityFixRegression ensures security fixes continue to work correctly
func testSecurityFixRegression(t *testing.T) {
	fixRegressionTests := []struct {
		name           string
		vulnerableCode string
		fixFunction    func(string) string
		expectedFix    []string
		shouldChange   bool
		description    string
	}{
		{
			name:           "CORS null origin fix regression",
			vulnerableCode: `c.Header("Access-Control-Allow-Origin", "null")`,
			fixFunction:    fixCORSNullOrigin,
			expectedFix:    []string{"SECURITY FIX", "omit the header entirely"},
			shouldChange:   true,
			description:    "CORS null origin fix must continue to work",
		},
		{
			name:           "CORS wildcard fix regression",
			vulnerableCode: `res.setHeader('Access-Control-Allow-Origin', '*')`,
			fixFunction:    fixCORSWildcard,
			expectedFix:    []string{"isAllowedOrigin", "SECURITY FIX"},
			shouldChange:   true,
			description:    "CORS wildcard fix must continue to work",
		},
		{
			name:           "JWT none algorithm fix regression",
			vulnerableCode: `algorithm: "none"`,
			fixFunction:    fixJWTNoneAlgorithm,
			expectedFix:    []string{"HS256", "SECURITY FIX"},
			shouldChange:   true,
			description:    "JWT none algorithm fix must continue to work",
		},
		{
			name:           "Security headers fix regression",
			vulnerableCode: `c.Header("Content-Type", "application/json")`,
			fixFunction:    AddSecurityHeaders,
			expectedFix:    []string{"X-Content-Type-Options", "X-Frame-Options", "X-XSS-Protection"},
			shouldChange:   true,
			description:    "Security headers fix must continue to work",
		},
	}

	for _, test := range fixRegressionTests {
		t.Run(test.name, func(t *testing.T) {
			result := test.fixFunction(test.vulnerableCode)

			if test.shouldChange {
				// Verify fix was applied
				if result == test.vulnerableCode {
					t.Errorf("REGRESSION: %s - Expected fix to be applied, but code remained unchanged", test.description)
				}

				// Verify expected fix content
				for _, expected := range test.expectedFix {
					if !strings.Contains(result, expected) {
						t.Errorf("REGRESSION: %s - Expected fix result to contain %q, got %q", test.description, expected, result)
					}
				}
			} else {
				// Verify safe code was not modified
				if result != test.vulnerableCode {
					t.Errorf("REGRESSION: %s - Safe code should not be modified, but got %q", test.description, result)
				}
			}
		})
	}
}

// testSafeCodePreservation ensures safe code patterns are never flagged as vulnerabilities
func testSafeCodePreservation(t *testing.T) {
	safePatterns := []struct {
		name        string
		code        string
		description string
	}{
		{
			name:        "Secure CORS with specific origin",
			code:        `c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")`,
			description: "Specific trusted origins should never be flagged",
		},
		{
			name:        "Parameterized SQL query",
			code:        `db.Query("SELECT * FROM users WHERE id = $1", userID)`,
			description: "Parameterized queries should never be flagged",
		},
		{
			name:        "Secure JWT algorithm",
			code:        `token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`,
			description: "Secure JWT algorithms should never be flagged",
		},
		{
			name:        "Environment-based configuration",
			code:        `debug: os.Getenv("DEBUG") == "true"`,
			description: "Environment-based config should never be flagged",
		},
		{
			name:        "Generic error message",
			code:        `return errors.New("operation failed")`,
			description: "Generic errors should never be flagged",
		},
		{
			name:        "HTTPS URL",
			code:        `const API_URL = "https://api.example.com"`,
			description: "HTTPS URLs should never be flagged",
		},
		{
			name:        "Secure cookie with flags",
			code:        `http.SetCookie(w, &http.Cookie{Name: "session", Value: id, HttpOnly: true, Secure: true})`,
			description: "Secure cookies should never be flagged",
		},
	}

	scanner := NewScanner()
	fixer := NewFixer()

	for _, pattern := range safePatterns {
		t.Run(pattern.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "safe_code.go.tmpl")

			err := os.WriteFile(testFile, []byte(pattern.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test scanning - should not detect issues
			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			if len(issues) > 0 {
				t.Errorf("REGRESSION: Safe pattern flagged as vulnerability: %s", pattern.description)
				for _, issue := range issues {
					t.Logf("False positive: %s - %s", issue.IssueType, issue.Description)
				}
			}

			// Test fixing - should not modify safe code
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "all",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Security fix failed: %v", err)
			}

			if len(result.FixedIssues) > 0 {
				t.Errorf("REGRESSION: Safe code was modified: %s", pattern.description)
				for _, fix := range result.FixedIssues {
					t.Logf("Unnecessary fix: %s - %s", fix.IssueType, fix.Description)
				}
			}

			// Verify file content remains unchanged
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			if string(content) != pattern.code {
				t.Errorf("REGRESSION: Safe code was modified: expected %q, got %q", pattern.code, string(content))
			}
		})
	}
}

// testSecurityMetricsConsistency ensures security metrics remain accurate over time
func testSecurityMetricsConsistency(t *testing.T) {
	// Create a comprehensive test template with known security issues
	testTemplate := `package main

// CORS vulnerabilities (2 critical issues)
c.Header("Access-Control-Allow-Origin", "null")
res.setHeader('Access-Control-Allow-Origin', '*')

// Authentication issues (3 issues: 1 critical, 1 high, 1 low)
token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
jwt_secret := "password"
jwt.sign(payload, secretKey)

// SQL injection risks (2 critical issues)
query := "SELECT * FROM users WHERE id = " + userID
query2 := fmt.Sprintf("DELETE FROM users WHERE id = %s", id)

// Information leakage (2 medium issues)
return fmt.Errorf("database error: %v", err)
debug: true

// Missing security headers (1 low issue)
c.Header("Content-Type", "application/json")`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "metrics_consistency.go.tmpl")

	err := os.WriteFile(testFile, []byte(testTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	scanner := NewScanner()
	issues, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Security scan failed: %v", err)
	}

	// Expected counts based on the test template
	expectedCounts := map[SecurityIssueType]int{
		CORSVulnerability:     2,
		WeakAuthentication:    3,
		SQLInjectionRisk:      2,
		InformationLeakage:    2,
		MissingSecurityHeader: 1,
	}

	expectedSeverityCounts := map[SeverityLevel]int{
		SeverityCritical: 5, // 2 CORS + 1 JWT none + 2 SQL
		SeverityHigh:     1, // 1 weak secret
		SeverityMedium:   2, // 2 info leakage
		SeverityLow:      2, // 1 JWT signing + 1 missing headers
	}

	// Verify issue type counts
	actualCounts := make(map[SecurityIssueType]int)
	for _, issue := range issues {
		actualCounts[issue.IssueType]++
	}

	for expectedType, expectedCount := range expectedCounts {
		actualCount := actualCounts[expectedType]
		if actualCount != expectedCount {
			t.Errorf("REGRESSION: Expected %d %s issues, got %d", expectedCount, expectedType, actualCount)
		}
	}

	// Verify severity distribution
	actualSeverityCounts := make(map[SeverityLevel]int)
	for _, issue := range issues {
		actualSeverityCounts[issue.Severity]++
	}

	for expectedSeverity, expectedCount := range expectedSeverityCounts {
		actualCount := actualSeverityCounts[expectedSeverity]
		if actualCount != expectedCount {
			t.Errorf("REGRESSION: Expected %d %s severity issues, got %d", expectedCount, expectedSeverity, actualCount)
		}
	}

	t.Logf("Security metrics consistency verified:")
	t.Logf("- Total issues: %d", len(issues))
	t.Logf("- Critical: %d, High: %d, Medium: %d, Low: %d",
		actualSeverityCounts[SeverityCritical],
		actualSeverityCounts[SeverityHigh],
		actualSeverityCounts[SeverityMedium],
		actualSeverityCounts[SeverityLow])
}

// testFixIdempotency ensures security fixes can be applied multiple times without issues
func testFixIdempotency(t *testing.T) {
	idempotencyTests := []struct {
		name    string
		content string
		fixType string
	}{
		{
			name:    "CORS fixes idempotency",
			content: `c.Header("Access-Control-Allow-Origin", "null")`,
			fixType: "cors",
		},
		{
			name:    "Authentication fixes idempotency",
			content: `algorithm: "none"`,
			fixType: "auth",
		},
		{
			name:    "Security headers idempotency",
			content: `c.Header("Content-Type", "application/json")`,
			fixType: "headers",
		},
		{
			name:    "All fixes idempotency",
			content: `c.Header("Access-Control-Allow-Origin", "null")\nalgorithm: "none"\nc.Header("Content-Type", "application/json")`,
			fixType: "all",
		},
	}

	fixer := NewFixer()

	for _, test := range idempotencyTests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "idempotency_test.go.tmpl")

			err := os.WriteFile(testFile, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      test.fixType,
				CreateBackup: false,
			}

			// Apply fixes first time
			_, err = fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("First fix application failed: %v", err)
			}

			// Read content after first fix
			content1, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file after first fix: %v", err)
			}

			// Apply fixes second time
			result2, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Second fix application failed: %v", err)
			}

			// Verify no additional fixes were applied
			if len(result2.FixedIssues) > 0 {
				t.Errorf("REGRESSION: Expected no fixes on second run, but %d were applied", len(result2.FixedIssues))
				for _, fix := range result2.FixedIssues {
					t.Logf("Unexpected fix: %s - %s", fix.IssueType, fix.Description)
				}
			}

			// Read content after second fix
			content2, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file after second fix: %v", err)
			}

			// Verify content is identical (idempotency)
			if string(content1) != string(content2) {
				t.Errorf("REGRESSION: Fix is not idempotent - content changed on second application")
				t.Logf("First result length: %d", len(content1))
				t.Logf("Second result length: %d", len(content2))
			}
		})
	}
}

// TestSecurityRegressionTimeBasedChecks performs time-based regression checks
func TestSecurityRegressionTimeBasedChecks(t *testing.T) {
	// This test ensures that security patterns don't degrade over time
	// and that performance doesn't regress

	startTime := time.Now()

	// Create a large template with many security issues
	var largeTemplate strings.Builder
	largeTemplate.WriteString("package main\n\n")

	// Add 100 different security issues
	for i := 0; i < 100; i++ {
		largeTemplate.WriteString("// Security issue set " + string(rune(i)) + "\n")
		largeTemplate.WriteString(`c.Header("Access-Control-Allow-Origin", "null")` + "\n")
		largeTemplate.WriteString(`query := "SELECT * FROM users WHERE id = " + userID` + "\n")
		largeTemplate.WriteString(`algorithm: "none"` + "\n")
		largeTemplate.WriteString(`return fmt.Errorf("database error: %v", err)` + "\n")
		largeTemplate.WriteString(`c.Header("Content-Type", "application/json")` + "\n\n")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large_regression_test.go.tmpl")

	err := os.WriteFile(testFile, []byte(largeTemplate.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test scanning performance
	scanner := NewScanner()
	scanStart := time.Now()
	issues, err := scanner.ScanFile(testFile)
	scanDuration := time.Since(scanStart)

	if err != nil {
		t.Fatalf("Large file scan failed: %v", err)
	}

	// Should find approximately 500 issues (5 per iteration * 100 iterations)
	expectedIssues := 500
	tolerance := 50 // Allow some tolerance for pattern matching variations

	if len(issues) < expectedIssues-tolerance || len(issues) > expectedIssues+tolerance {
		t.Errorf("REGRESSION: Expected approximately %d issues, got %d", expectedIssues, len(issues))
	}

	// Performance regression check - scanning should complete within reasonable time
	maxScanTime := 5 * time.Second
	if scanDuration > maxScanTime {
		t.Errorf("REGRESSION: Scanning performance degraded - took %v, expected < %v", scanDuration, maxScanTime)
	}

	// Test fixing performance
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      false,
		FixType:      "all",
		CreateBackup: false,
	}

	fixStart := time.Now()
	result, err := fixer.FixFile(testFile, options)
	fixDuration := time.Since(fixStart)

	if err != nil {
		t.Fatalf("Large file fix failed: %v", err)
	}

	// Should apply many fixes
	if len(result.FixedIssues) == 0 {
		t.Error("REGRESSION: Expected fixes to be applied to large file")
	}

	// Performance regression check - fixing should complete within reasonable time
	maxFixTime := 10 * time.Second
	if fixDuration > maxFixTime {
		t.Errorf("REGRESSION: Fixing performance degraded - took %v, expected < %v", fixDuration, maxFixTime)
	}

	totalDuration := time.Since(startTime)
	t.Logf("Regression test completed in %v:", totalDuration)
	t.Logf("- Scan time: %v (%d issues found)", scanDuration, len(issues))
	t.Logf("- Fix time: %v (%d fixes applied)", fixDuration, len(result.FixedIssues))
}

// TestSecurityPatternEvolutionRegression ensures security patterns evolve appropriately
func TestSecurityPatternEvolutionRegression(t *testing.T) {
	// Test that security patterns cover modern security concerns
	// This test will fail if important security categories are missing

	patterns := getSecurityPatterns()

	// Required security categories that must be covered
	requiredCategories := map[string][]string{
		"CORS Security": {
			"null origin",
			"wildcard",
			"credentials",
		},
		"Authentication": {
			"jwt",
			"algorithm",
			"secret",
			"cookie",
		},
		"Injection Prevention": {
			"sql",
			"concatenation",
			"format",
		},
		"Information Disclosure": {
			"error",
			"debug",
			"trace",
		},
		"Transport Security": {
			"http",
			"header",
			"security",
		},
	}

	for category, keywords := range requiredCategories {
		t.Run("Category_"+category, func(t *testing.T) {
			categoryPatterns := 0

			for _, pattern := range patterns {
				patternText := strings.ToLower(pattern.Name + " " + pattern.Description)

				for _, keyword := range keywords {
					if strings.Contains(patternText, strings.ToLower(keyword)) {
						categoryPatterns++
						break
					}
				}
			}

			if categoryPatterns == 0 {
				t.Errorf("REGRESSION: No security patterns found for category: %s", category)
				t.Errorf("Consider adding patterns that include keywords: %v", keywords)
			} else {
				t.Logf("Found %d patterns for %s category", categoryPatterns, category)
			}
		})
	}

	// Verify minimum number of patterns exist
	minPatterns := 15
	if len(patterns) < minPatterns {
		t.Errorf("REGRESSION: Expected at least %d security patterns, got %d", minPatterns, len(patterns))
	}

	// Verify all patterns have required fields
	for i, pattern := range patterns {
		if pattern.Name == "" {
			t.Errorf("REGRESSION: Pattern %d missing name", i)
		}
		if pattern.Pattern == nil {
			t.Errorf("REGRESSION: Pattern %d missing regex", i)
		}
		if pattern.Description == "" {
			t.Errorf("REGRESSION: Pattern %d missing description", i)
		}
		if pattern.Recommendation == "" {
			t.Errorf("REGRESSION: Pattern %d missing recommendation", i)
		}
	}
}
