package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCORSSecurityFixesValidation provides comprehensive unit tests for CORS security fixes
func TestCORSSecurityFixesValidation(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedFixed      bool
		expectedContent    []string
		notExpectedContent []string
		framework          string
		description        string
	}{
		{
			name:               "Go Gin CORS null origin vulnerability",
			input:              `    c.Header("Access-Control-Allow-Origin", "null")`,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY FIX", "Removed Access-Control-Allow-Origin: null header", "omit the header entirely"},
			notExpectedContent: []string{`"null"`},
			framework:          "gin",
			description:        "Should remove null origin header and add security comment",
		},
		{
			name:               "Node.js Express CORS null origin vulnerability",
			input:              `  res.setHeader('Access-Control-Allow-Origin', 'null');`,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY FIX", "Removed Access-Control-Allow-Origin: null header", "omit the header entirely"},
			notExpectedContent: []string{`'null'`},
			framework:          "express",
			description:        "Should remove null origin header for Express framework",
		},
		{
			name:               "Go Gin CORS wildcard vulnerability",
			input:              `    c.Header("Access-Control-Allow-Origin", "*")`,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY FIX", "isAllowedOrigin", "c.Header"},
			notExpectedContent: []string{},
			framework:          "gin",
			description:        "Should replace wildcard with origin validation logic",
		},
		{
			name:               "Node.js Express CORS wildcard vulnerability",
			input:              `  res.setHeader('Access-Control-Allow-Origin', '*');`,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY FIX", "isAllowedOrigin", "res.setHeader"},
			notExpectedContent: []string{},
			framework:          "express",
			description:        "Should replace wildcard with origin validation for Express",
		},
		{
			name:               "Safe CORS configuration should not be modified",
			input:              `    c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")`,
			expectedFixed:      false,
			expectedContent:    []string{},
			notExpectedContent: []string{"SECURITY FIX"},
			framework:          "gin",
			description:        "Safe CORS configurations should remain unchanged",
		},
		{
			name:               "Environment-based CORS should not be modified",
			input:              `    c.Header("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))`,
			expectedFixed:      false,
			expectedContent:    []string{},
			notExpectedContent: []string{"SECURITY FIX"},
			framework:          "gin",
			description:        "Environment-based CORS should not be flagged",
		},
		{
			name:               "CORS with credentials and wildcard",
			input:              `    c.Header("Access-Control-Allow-Origin", "*")\n    c.Header("Access-Control-Allow-Credentials", "true")`,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY FIX", "isAllowedOrigin"},
			notExpectedContent: []string{},
			framework:          "gin",
			description:        "CORS wildcard with credentials should be fixed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string

			// Apply appropriate CORS fix based on vulnerability type
			if strings.Contains(tt.input, "null") {
				result = fixCORSNullOrigin(tt.input)
			} else if strings.Contains(tt.input, "*") {
				result = fixCORSWildcard(tt.input)
			} else {
				result = tt.input
			}

			// Verify if fix was applied as expected
			wasFixed := result != tt.input
			if wasFixed != tt.expectedFixed {
				if tt.expectedFixed {
					t.Errorf("Expected CORS fix to be applied, but input remained unchanged")
				} else {
					t.Errorf("Expected no CORS fix for safe configuration, but got: %q", result)
				}
			}

			// Verify expected content is present
			for _, expected := range tt.expectedContent {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected CORS fix result to contain %q, got %q", expected, result)
				}
			}

			// Verify unwanted content is not present
			for _, notExpected := range tt.notExpectedContent {
				if strings.Contains(result, notExpected) {
					t.Errorf("CORS fix result should not contain %q, got %q", notExpected, result)
				}
			}
		})
	}
}

// TestCORSVulnerabilityDetection tests detection of CORS vulnerabilities
func TestCORSVulnerabilityDetection(t *testing.T) {
	scanner := NewScanner()

	vulnerabilityTests := []struct {
		name         string
		code         string
		shouldDetect bool
		expectedType SecurityIssueType
		severity     SeverityLevel
	}{
		{
			name:         "CORS null origin should be detected as critical",
			code:         `c.Header("Access-Control-Allow-Origin", "null")`,
			shouldDetect: true,
			expectedType: CORSVulnerability,
			severity:     SeverityCritical,
		},
		{
			name:         "CORS wildcard should be detected as medium",
			code:         `res.setHeader('Access-Control-Allow-Origin', '*')`,
			shouldDetect: true,
			expectedType: CORSVulnerability,
			severity:     SeverityMedium,
		},
		{
			name:         "Safe CORS should not be detected",
			code:         `c.Header("Access-Control-Allow-Origin", "https://example.com")`,
			shouldDetect: false,
		},
		{
			name:         "Environment CORS should not be detected",
			code:         `c.Header("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))`,
			shouldDetect: false,
		},
	}

	for _, tt := range vulnerabilityTests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "cors_test.go.tmpl")

			err := os.WriteFile(testFile, []byte(tt.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("CORS vulnerability scan failed: %v", err)
			}

			corsIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == CORSVulnerability {
					corsIssueFound = true
					if tt.shouldDetect && issue.Severity != tt.severity {
						t.Errorf("Expected CORS vulnerability severity %s, got %s", tt.severity, issue.Severity)
					}
					break
				}
			}

			if tt.shouldDetect && !corsIssueFound {
				t.Errorf("Expected CORS vulnerability to be detected, but none found")
			}

			if !tt.shouldDetect && corsIssueFound {
				t.Errorf("CORS vulnerability should not be detected for safe code")
			}
		})
	}
}

// TestCORSFixIdempotency ensures CORS fixes can be applied multiple times safely
func TestCORSFixIdempotency(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "CORS null origin fix idempotency",
			input: `c.Header("Access-Control-Allow-Origin", "null")`,
		},
		{
			name:  "CORS wildcard fix idempotency",
			input: `res.setHeader('Access-Control-Allow-Origin', '*')`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "cors_idempotent.go.tmpl")

			err := os.WriteFile(testFile, []byte(tc.input), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			fixer := NewFixer()
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "cors",
				CreateBackup: false,
			}

			// Apply fixes first time
			result1, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("First CORS fix application failed: %v", err)
			}

			if len(result1.FixedIssues) == 0 {
				t.Error("Expected CORS fixes to be applied on first run")
			}

			// Read content after first fix
			content1, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file after first fix: %v", err)
			}

			// Apply fixes second time
			result2, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Second CORS fix application failed: %v", err)
			}

			// Verify no additional fixes were applied
			if len(result2.FixedIssues) > 0 {
				t.Errorf("Expected no CORS fixes on second run, but %d were applied", len(result2.FixedIssues))
			}

			// Read content after second fix
			content2, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file after second fix: %v", err)
			}

			// Verify content is identical
			if string(content1) != string(content2) {
				t.Error("CORS fix should be idempotent - file content should be identical after second application")
			}
		})
	}
}

// TestCORSFixPreservesIndentation ensures CORS fixes maintain proper code formatting
func TestCORSFixPreservesIndentation(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedIndent string
	}{
		{
			name:           "Spaces indentation",
			input:          `    c.Header("Access-Control-Allow-Origin", "null")`,
			expectedIndent: "    ",
		},
		{
			name:           "Tabs indentation",
			input:          "\t\tc.Header(\"Access-Control-Allow-Origin\", \"null\")",
			expectedIndent: "\t\t",
		},
		{
			name:           "Mixed indentation",
			input:          " \t c.Header(\"Access-Control-Allow-Origin\", \"null\")",
			expectedIndent: " \t ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fixCORSNullOrigin(tc.input)

			// Verify original indentation is preserved in fix comments
			lines := strings.Split(result, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" && strings.HasPrefix(line, "//") {
					if !strings.HasPrefix(line, tc.expectedIndent+"//") {
						t.Errorf("CORS fix should preserve indentation %q, got line: %q", tc.expectedIndent, line)
					}
				}
			}
		})
	}
}

// TestCORSSecurityRegressionPrevention ensures CORS fixes don't introduce new vulnerabilities
func TestCORSSecurityRegressionPrevention(t *testing.T) {
	// Test that CORS fixes don't accidentally introduce new security issues
	regressionTests := []struct {
		name        string
		input       string
		fixFunction func(string) string
		checkFor    []string
	}{
		{
			name:        "CORS null fix doesn't introduce eval",
			input:       `c.Header("Access-Control-Allow-Origin", "null")`,
			fixFunction: fixCORSNullOrigin,
			checkFor:    []string{"eval", "innerHTML", "document.write"},
		},
		{
			name:        "CORS wildcard fix doesn't introduce SQL injection",
			input:       `res.setHeader('Access-Control-Allow-Origin', '*')`,
			fixFunction: fixCORSWildcard,
			checkFor:    []string{"SELECT *", "INSERT INTO", "DELETE FROM"},
		},
	}

	for _, tt := range regressionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			// Verify fix doesn't introduce dangerous patterns
			for _, dangerousPattern := range tt.checkFor {
				if strings.Contains(strings.ToLower(result), strings.ToLower(dangerousPattern)) {
					t.Errorf("CORS fix introduced dangerous pattern %q in result: %q", dangerousPattern, result)
				}
			}

			// Verify fix doesn't contain obvious security anti-patterns
			antiPatterns := []string{
				"password",
				"secret123",
				"admin",
				"root",
				"<script>",
				"javascript:",
			}

			for _, antiPattern := range antiPatterns {
				if strings.Contains(strings.ToLower(result), antiPattern) {
					t.Errorf("CORS fix introduced security anti-pattern %q in result: %q", antiPattern, result)
				}
			}
		})
	}
}
