package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSecurityRegressionSuite runs comprehensive regression tests to prevent future vulnerabilities
func TestSecurityRegressionSuite(t *testing.T) {
	// Define known secure patterns that should never be flagged as vulnerabilities
	securePatterns := []struct {
		name        string
		code        string
		description string
	}{
		{
			name:        "Secure CORS with specific origin",
			code:        `c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")`,
			description: "Specific trusted origin should not be flagged",
		},
		{
			name:        "Secure CORS with environment variable",
			code:        `c.Header("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))`,
			description: "Environment-based origin should not be flagged",
		},
		{
			name:        "Parameterized SQL query",
			code:        `db.Query("SELECT * FROM users WHERE id = $1", userID)`,
			description: "Parameterized queries should not be flagged",
		},
		{
			name:        "Named SQL parameters",
			code:        `db.Query("SELECT * FROM users WHERE id = :id", sql.Named("id", userID))`,
			description: "Named parameters should not be flagged",
		},
		{
			name:        "Secure JWT algorithm",
			code:        `token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`,
			description: "Secure JWT algorithms should not be flagged",
		},
		{
			name:        "JWT with expiration",
			code:        `jwt.sign(payload, secret, { expiresIn: '15m' })`,
			description: "JWT with expiration should not be flagged",
		},
		{
			name:        "Generic error message",
			code:        `return errors.New("operation failed")`,
			description: "Generic errors should not be flagged",
		},
		{
			name:        "Environment-based debug config",
			code:        `debug: os.Getenv("DEBUG") == "true"`,
			description: "Environment-based debug should not be flagged",
		},
		{
			name:        "Secure cookie with flags",
			code:        `http.SetCookie(w, &http.Cookie{Name: "session", Value: id, HttpOnly: true, Secure: true})`,
			description: "Secure cookies should not be flagged",
		},
		{
			name:        "HTTPS URL",
			code:        `const API_URL = "https://api.example.com"`,
			description: "HTTPS URLs should not be flagged",
		},
	}

	scanner := NewScanner()

	for _, pattern := range securePatterns {
		t.Run(pattern.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "secure.go.tmpl")

			err := os.WriteFile(testFile, []byte(pattern.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Secure patterns should not generate security issues
			if len(issues) > 0 {
				t.Errorf("Secure pattern flagged as vulnerability: %s", pattern.description)
				for _, issue := range issues {
					t.Logf("False positive: %s - %s", issue.IssueType, issue.Description)
				}
			}
		})
	}
}

// TestKnownVulnerabilityPatterns ensures all known vulnerability patterns are still detected
func TestKnownVulnerabilityPatterns(t *testing.T) {
	// Define known vulnerability patterns that should always be detected
	vulnerabilityPatterns := []struct {
		name         string
		code         string
		expectedType SecurityIssueType
		severity     SeverityLevel
		description  string
	}{
		{
			name:         "CORS null origin",
			code:         `c.Header("Access-Control-Allow-Origin", "null")`,
			expectedType: CORSVulnerability,
			severity:     SeverityCritical,
			description:  "CORS null origin should always be detected",
		},
		{
			name:         "CORS wildcard",
			code:         `res.setHeader('Access-Control-Allow-Origin', '*')`,
			expectedType: CORSVulnerability,
			severity:     SeverityMedium,
			description:  "CORS wildcard should always be detected",
		},
		{
			name:         "JWT none algorithm",
			code:         `token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)`,
			expectedType: WeakAuthentication,
			severity:     SeverityCritical,
			description:  "JWT none algorithm should always be detected",
		},
		{
			name:         "SQL string concatenation",
			code:         `query := "SELECT * FROM users WHERE id = " + userID`,
			expectedType: SQLInjectionRisk,
			severity:     SeverityCritical,
			description:  "SQL string concatenation should always be detected",
		},
		{
			name:         "SQL format string",
			code:         `query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)`,
			expectedType: SQLInjectionRisk,
			severity:     SeverityHigh,
			description:  "SQL format string should always be detected",
		},
		{
			name:         "Weak JWT secret",
			code:         `jwt.sign(payload, "secret")`,
			expectedType: WeakAuthentication,
			severity:     SeverityHigh,
			description:  "Weak JWT secret should always be detected",
		},
		{
			name:         "Detailed database error",
			code:         `return fmt.Errorf("database error: %v", err)`,
			expectedType: InformationLeakage,
			severity:     SeverityMedium,
			description:  "Detailed database errors should always be detected",
		},
		{
			name:         "Debug enabled in production",
			code:         `debug: true`,
			expectedType: InformationLeakage,
			severity:     SeverityMedium,
			description:  "Debug enabled should always be detected",
		},
		{
			name:         "Hardcoded password",
			code:         `password := "mypassword123"`,
			expectedType: WeakAuthentication,
			severity:     SeverityHigh,
			description:  "Hardcoded passwords should always be detected",
		},
		{
			name:         "HTTP URL in production",
			code:         `const API_URL = "http://api.example.com"`,
			expectedType: WeakAuthentication,
			severity:     SeverityMedium,
			description:  "HTTP URLs should always be detected",
		},
	}

	scanner := NewScanner()

	for _, vuln := range vulnerabilityPatterns {
		t.Run(vuln.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "vulnerable.go.tmpl")

			err := os.WriteFile(testFile, []byte(vuln.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Vulnerability patterns should always generate security issues
			if len(issues) == 0 {
				t.Errorf("Known vulnerability not detected: %s", vuln.description)
				return
			}

			// Verify the correct issue type and severity
			found := false
			for _, issue := range issues {
				if issue.IssueType == vuln.expectedType {
					found = true
					if issue.Severity != vuln.severity {
						t.Errorf("Expected severity %s for %s, got %s", vuln.severity, vuln.name, issue.Severity)
					}
					break
				}
			}

			if !found {
				t.Errorf("Expected issue type %s for %s, but not found", vuln.expectedType, vuln.name)
			}
		})
	}
}

// TestSecurityFixRegression ensures security fixes don't break over time
func TestSecurityFixRegression(t *testing.T) {
	// Test cases that verify fixes continue to work correctly
	fixRegressionTests := []struct {
		name           string
		vulnerableCode string
		fixFunction    func(string) string
		expectedFix    []string
		shouldChange   bool
	}{
		{
			name:           "CORS null origin fix regression",
			vulnerableCode: `c.Header("Access-Control-Allow-Origin", "null")`,
			fixFunction:    fixCORSNullOrigin,
			expectedFix:    []string{"SECURITY FIX", "omit the header entirely"},
			shouldChange:   true,
		},
		{
			name:           "CORS wildcard fix regression",
			vulnerableCode: `res.setHeader('Access-Control-Allow-Origin', '*')`,
			fixFunction:    fixCORSWildcard,
			expectedFix:    []string{"isAllowedOrigin", "SECURITY FIX"},
			shouldChange:   true,
		},
		{
			name:           "JWT none algorithm fix regression",
			vulnerableCode: `algorithm: "none"`,
			fixFunction:    fixJWTNoneAlgorithm,
			expectedFix:    []string{"HS256", "SECURITY FIX"},
			shouldChange:   true,
		},
		{
			name:           "Security headers fix regression",
			vulnerableCode: `c.Header("Content-Type", "application/json")`,
			fixFunction:    addSecurityHeaders,
			expectedFix:    []string{"X-Content-Type-Options", "X-Frame-Options", "X-XSS-Protection"},
			shouldChange:   true,
		},
		{
			name:           "Safe code should not be modified",
			vulnerableCode: `c.Header("Access-Control-Allow-Origin", "https://trusted.com")`,
			fixFunction:    fixCORSNullOrigin,
			shouldChange:   false,
		},
	}

	for _, test := range fixRegressionTests {
		t.Run(test.name, func(t *testing.T) {
			result := test.fixFunction(test.vulnerableCode)

			if test.shouldChange {
				// Verify fix was applied
				if result == test.vulnerableCode {
					t.Errorf("Expected fix to be applied, but code remained unchanged")
				}

				// Verify expected fix content
				for _, expected := range test.expectedFix {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected fix result to contain %q, got %q", expected, result)
					}
				}
			} else {
				// Verify safe code was not modified
				if result != test.vulnerableCode {
					t.Errorf("Safe code should not be modified, but got %q", result)
				}
			}
		})
	}
}

// TestNewVulnerabilityDetection tests for detection of newly discovered vulnerability patterns
func TestNewVulnerabilityDetection(t *testing.T) {
	// Simulate newly discovered vulnerability patterns that should be detected
	newVulnerabilities := []struct {
		name        string
		code        string
		description string
	}{
		{
			name: "CORS with credentials and wildcard",
			code: `res.setHeader('Access-Control-Allow-Origin', '*');
res.setHeader('Access-Control-Allow-Credentials', 'true');`,
			description: "CORS wildcard with credentials should be detected",
		},
		{
			name:        "Template injection in SQL",
			code:        `query := "SELECT * FROM users WHERE name = '${userName}'"`,
			description: "Template literal injection should be detected",
		},
		{
			name:        "Eval-like functions",
			code:        `result := eval(userInput)`,
			description: "Eval functions should be detected as dangerous",
		},
		{
			name:        "Insecure random generation",
			code:        `sessionID := fmt.Sprintf("%d", rand.Int())`,
			description: "Weak random generation should be detected",
		},
	}

	scanner := NewScanner()

	for _, vuln := range newVulnerabilities {
		t.Run(vuln.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "new_vuln.go.tmpl")

			err := os.WriteFile(testFile, []byte(vuln.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Note: Some of these may not be detected by current patterns
			// This test serves as a reminder to add detection for new vulnerabilities
			if len(issues) == 0 {
				t.Logf("New vulnerability pattern not yet detected: %s", vuln.description)
				t.Logf("Consider adding detection pattern for: %s", vuln.code)
			} else {
				t.Logf("Successfully detected potential vulnerability: %s", vuln.description)
			}
		})
	}
}

// TestSecurityPatternEvolution ensures security patterns evolve with new threats
func TestSecurityPatternEvolution(t *testing.T) {
	// Test that our security patterns cover modern security concerns
	modernSecurityConcerns := []struct {
		category string
		patterns []string
	}{
		{
			category: "CORS Security",
			patterns: []string{
				"null origin handling",
				"wildcard with credentials",
				"overly permissive origins",
			},
		},
		{
			category: "Authentication Security",
			patterns: []string{
				"JWT none algorithm",
				"weak secrets",
				"missing expiration",
				"insecure cookies",
			},
		},
		{
			category: "Injection Prevention",
			patterns: []string{
				"SQL injection",
				"template injection",
				"command injection",
			},
		},
		{
			category: "Information Disclosure",
			patterns: []string{
				"detailed error messages",
				"debug information exposure",
				"stack trace leakage",
			},
		},
		{
			category: "Transport Security",
			patterns: []string{
				"HTTP URLs",
				"missing security headers",
				"insecure protocols",
			},
		},
	}

	patterns := getSecurityPatterns()

	for _, concern := range modernSecurityConcerns {
		t.Run(concern.category, func(t *testing.T) {
			// Verify we have patterns for each security concern category
			categoryPatterns := 0
			for _, pattern := range patterns {
				// Check if pattern addresses this security concern
				patternName := strings.ToLower(pattern.Name)
				patternDesc := strings.ToLower(pattern.Description)

				for _, expectedPattern := range concern.patterns {
					if strings.Contains(patternName, strings.ToLower(expectedPattern)) ||
						strings.Contains(patternDesc, strings.ToLower(expectedPattern)) {
						categoryPatterns++
						break
					}
				}
			}

			if categoryPatterns == 0 {
				t.Errorf("No security patterns found for category: %s", concern.category)
				t.Logf("Consider adding patterns for: %v", concern.patterns)
			} else {
				t.Logf("Found %d patterns for %s category", categoryPatterns, concern.category)
			}
		})
	}
}

// TestSecurityFixCompatibility ensures fixes work across different template formats
func TestSecurityFixCompatibility(t *testing.T) {
	// Test security fixes across different file formats and frameworks
	compatibilityTests := []struct {
		name      string
		extension string
		framework string
		template  string
	}{
		{
			name:      "Go template",
			extension: ".go.tmpl",
			framework: "gin",
			template:  `c.Header("Access-Control-Allow-Origin", "null")`,
		},
		{
			name:      "JavaScript template",
			extension: ".js.tmpl",
			framework: "express",
			template:  `res.setHeader('Access-Control-Allow-Origin', 'null');`,
		},
		{
			name:      "TypeScript template",
			extension: ".ts.tmpl",
			framework: "express",
			template:  `response.setHeader('Access-Control-Allow-Origin', 'null');`,
		},
		{
			name:      "YAML configuration",
			extension: ".yaml.tmpl",
			framework: "config",
			template:  `cors:\n  origin: "null"`,
		},
		{
			name:      "JSON configuration",
			extension: ".json.tmpl",
			framework: "config",
			template:  `{"cors": {"origin": "null"}}`,
		},
	}

	fixer := NewFixer()

	for _, test := range compatibilityTests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "template"+test.extension)

			err := os.WriteFile(testFile, []byte(test.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "all",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Security fix failed for %s: %v", test.name, err)
			}

			// Verify fixes can be applied to different file formats
			// Note: Some fixes may not apply to all formats (e.g., YAML/JSON)
			if len(result.FixedIssues) > 0 {
				t.Logf("Successfully applied %d fixes to %s format", len(result.FixedIssues), test.extension)
			} else {
				t.Logf("No fixes applied to %s format (may be expected)", test.extension)
			}

			// Verify no errors occurred during fixing
			if len(result.Errors) > 0 {
				t.Errorf("Errors occurred while fixing %s format: %v", test.extension, result.Errors)
			}
		})
	}
}

// TestSecurityMetricsRegression ensures security metrics remain accurate
func TestSecurityMetricsRegression(t *testing.T) {
	// Create a test template with known security issues
	testTemplate := `package main

// CORS vulnerabilities (2 issues)
c.Header("Access-Control-Allow-Origin", "null")
res.setHeader('Access-Control-Allow-Origin', '*')

// Authentication issues (2 issues)  
token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
jwt_secret := "password"

// SQL injection risks (2 issues)
query := "SELECT * FROM users WHERE id = " + userID
query2 := fmt.Sprintf("DELETE FROM users WHERE id = %s", id)

// Information leakage (2 issues)
return fmt.Errorf("database error: %v", err)
debug: true

// Missing security headers (1 issue)
c.Header("Content-Type", "application/json")`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "metrics_test.go.tmpl")

	err := os.WriteFile(testFile, []byte(testTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Scan for security issues
	scanner := NewScanner()
	issues, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Security scan failed: %v", err)
	}

	// Verify expected number of issues by type
	expectedCounts := map[SecurityIssueType]int{
		CORSVulnerability:     2,
		WeakAuthentication:    2,
		SQLInjectionRisk:      2,
		InformationLeakage:    2,
		MissingSecurityHeader: 1,
	}

	actualCounts := make(map[SecurityIssueType]int)
	for _, issue := range issues {
		actualCounts[issue.IssueType]++
	}

	for expectedType, expectedCount := range expectedCounts {
		actualCount := actualCounts[expectedType]
		if actualCount != expectedCount {
			t.Errorf("Expected %d %s issues, got %d", expectedCount, expectedType, actualCount)
		}
	}

	// Verify severity distribution
	severityCounts := make(map[SeverityLevel]int)
	for _, issue := range issues {
		severityCounts[issue.Severity]++
	}

	// Should have at least some critical and high severity issues
	if severityCounts[SeverityCritical] == 0 {
		t.Error("Expected at least one critical severity issue")
	}

	if severityCounts[SeverityHigh] == 0 {
		t.Error("Expected at least one high severity issue")
	}

	t.Logf("Security metrics regression test completed:")
	t.Logf("- Total issues found: %d", len(issues))
	t.Logf("- Critical: %d, High: %d, Medium: %d, Low: %d",
		severityCounts[SeverityCritical],
		severityCounts[SeverityHigh],
		severityCounts[SeverityMedium],
		severityCounts[SeverityLow])
}
