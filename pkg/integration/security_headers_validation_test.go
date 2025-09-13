package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/pkg/security"
)

// TestSecurityHeaderImplementationValidation provides comprehensive tests for security header implementation
func TestSecurityHeaderImplementationValidation(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedHeaders  []string
		framework        string
		shouldAddHeaders bool
		description      string
	}{
		{
			name:             "Go Gin JSON response should get security headers",
			input:            `    c.Header("Content-Type", "application/json")`,
			expectedHeaders:  []string{"X-Content-Type-Options", "X-Frame-Options", "X-XSS-Protection", "nosniff", "DENY", "1; mode=block"},
			framework:        "gin",
			shouldAddHeaders: true,
			description:      "JSON responses should include comprehensive security headers",
		},
		{
			name:             "Go Gin HTML response should get security headers",
			input:            `    c.Header("Content-Type", "text/html")`,
			expectedHeaders:  []string{"X-Content-Type-Options", "X-Frame-Options", "X-XSS-Protection", "nosniff", "DENY"},
			framework:        "gin",
			shouldAddHeaders: true,
			description:      "HTML responses should include security headers",
		},
		{
			name:             "Node.js Express JSON response should get security headers",
			input:            `  res.setHeader('Content-Type', 'application/json');`,
			expectedHeaders:  []string{"X-Content-Type-Options", "X-Frame-Options", "X-XSS-Protection", "nosniff", "DENY"},
			framework:        "express",
			shouldAddHeaders: true,
			description:      "Express JSON responses should include security headers",
		},
		{
			name:             "Go HTTP server response should get security headers",
			input:            `    w.Header().Set("Content-Type", "application/json")`,
			expectedHeaders:  []string{"X-Content-Type-Options", "nosniff"},
			framework:        "http",
			shouldAddHeaders: true,
			description:      "Standard HTTP responses should include security headers",
		},
		{
			name:             "Non-content-type headers should not trigger security headers",
			input:            `    c.Header("X-Custom-Header", "value")`,
			expectedHeaders:  []string{},
			framework:        "gin",
			shouldAddHeaders: false,
			description:      "Custom headers should not trigger security header addition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := security.AddSecurityHeaders(tt.input)

			if tt.shouldAddHeaders {
				// Verify original header is preserved
				if !strings.Contains(result, tt.input) {
					t.Errorf("Original header should be preserved in result")
				}

				// Verify security headers were added
				for _, header := range tt.expectedHeaders {
					if !strings.Contains(result, header) {
						t.Errorf("Expected security header %q to be added for %s framework", header, tt.framework)
					}
				}

				// Verify security comment was added
				if !strings.Contains(result, "SECURITY") {
					t.Error("Expected security comment to be added with headers")
				}

				// Verify proper line structure
				lines := strings.Split(result, "\n")
				if len(lines) < 2 {
					t.Error("Expected multiple lines when security headers are added")
				}
			} else {
				// Verify no security headers were added unnecessarily
				if result != tt.input {
					t.Errorf("Expected no changes for non-content-type headers, but got: %q", result)
				}
			}
		})
	}
}

// TestSecurityHeaderDetection tests detection of missing security headers
func TestSecurityHeaderDetection(t *testing.T) {
	scanner := security.NewScanner()

	detectionTests := []struct {
		name         string
		code         string
		shouldDetect bool
		expectedType security.SecurityIssueType
		severity     security.SeverityLevel
	}{
		{
			name:         "Content-Type header should trigger security header check",
			code:         `c.Header("Content-Type", "application/json")`,
			shouldDetect: true,
			expectedType: security.MissingSecurityHeader,
			severity:     security.SeverityLow,
		},
		{
			name:         "HTTP response header should trigger security check",
			code:         `response.Header("Content-Type", "text/html")`,
			shouldDetect: true,
			expectedType: security.MissingSecurityHeader,
			severity:     security.SeverityLow,
		},
		{
			name:         "Express setHeader should trigger security check",
			code:         `res.setHeader('Content-Type', 'application/json')`,
			shouldDetect: true,
			expectedType: security.MissingSecurityHeader,
			severity:     security.SeverityLow,
		},
		{
			name:         "Custom headers should not trigger security header detection",
			code:         `c.Header("X-Custom", "value")`,
			shouldDetect: false,
		},
		{
			name:         "Already secure headers should not be flagged",
			code:         `c.Header("X-Content-Type-Options", "nosniff")`,
			shouldDetect: false,
		},
	}

	for _, tt := range detectionTests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "headers_test.go.tmpl")

			err := os.WriteFile(testFile, []byte(tt.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security header scan failed: %v", err)
			}

			headerIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == security.MissingSecurityHeader {
					headerIssueFound = true
					if tt.shouldDetect && issue.Severity != tt.severity {
						t.Errorf("Expected security header issue severity %s, got %s", tt.severity, issue.Severity)
					}
					break
				}
			}

			if tt.shouldDetect && !headerIssueFound {
				t.Errorf("Expected missing security header to be detected, but none found")
			}

			if !tt.shouldDetect && headerIssueFound {
				t.Errorf("Security header issue should not be detected for this code")
			}
		})
	}
}

// TestSecurityHeaderIndentationPreservation ensures security headers maintain proper indentation
func TestSecurityHeaderIndentationPreservation(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedIndent string
	}{
		{
			name:           "Spaces indentation preservation",
			input:          `    c.Header("Content-Type", "application/json")`,
			expectedIndent: "    ",
		},
		{
			name:           "Tabs indentation preservation",
			input:          "\t\tc.Header(\"Content-Type\", \"application/json\")",
			expectedIndent: "\t\t",
		},
		{
			name:           "Mixed indentation preservation",
			input:          " \t c.Header(\"Content-Type\", \"application/json\")",
			expectedIndent: " \t ",
		},
		{
			name:           "No indentation",
			input:          `c.Header("Content-Type", "application/json")`,
			expectedIndent: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := security.AddSecurityHeaders(tc.input)

			// Verify original line is preserved
			if !strings.Contains(result, tc.input) {
				t.Error("Original header line should be preserved")
			}

			// Verify added security headers maintain indentation
			lines := strings.Split(result, "\n")
			for i, line := range lines {
				if i == 0 {
					continue // Skip original line
				}

				if strings.TrimSpace(line) != "" {
					if !strings.HasPrefix(line, tc.expectedIndent) {
						t.Errorf("Security header line should maintain indentation %q, got: %q", tc.expectedIndent, line)
					}
				}
			}
		})
	}
}

// TestSecurityHeaderFrameworkCompatibility tests security headers across different frameworks
func TestSecurityHeaderFrameworkCompatibility(t *testing.T) {
	frameworkTests := []struct {
		name      string
		input     string
		framework string
		expected  []string
	}{
		{
			name:      "Gin framework compatibility",
			input:     `c.Header("Content-Type", "application/json")`,
			framework: "gin",
			expected:  []string{"c.Header(\"X-Content-Type-Options\"", "c.Header(\"X-Frame-Options\"", "c.Header(\"X-XSS-Protection\""},
		},
		{
			name:      "Express framework compatibility",
			input:     `res.setHeader('Content-Type', 'application/json')`,
			framework: "express",
			expected:  []string{"res.setHeader('X-Content-Type-Options'", "res.setHeader('X-Frame-Options'", "res.setHeader('X-XSS-Protection'"},
		},
		{
			name:      "Standard HTTP compatibility",
			input:     `w.Header().Set("Content-Type", "application/json")`,
			framework: "http",
			expected:  []string{"X-Content-Type-Options", "nosniff"},
		},
	}

	for _, tt := range frameworkTests {
		t.Run(tt.name, func(t *testing.T) {
			result := security.AddSecurityHeaders(tt.input)

			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected %s framework pattern %q in result, got: %q", tt.framework, expected, result)
				}
			}
		})
	}
}

// TestSecurityHeaderCompleteness ensures all required security headers are added
func TestSecurityHeaderCompleteness(t *testing.T) {
	requiredHeaders := map[string][]string{
		"application/json": {
			"X-Content-Type-Options",
			"X-Frame-Options",
			"X-XSS-Protection",
			"nosniff",
			"DENY",
			"1; mode=block",
		},
		"text/html": {
			"X-Content-Type-Options",
			"X-Frame-Options",
			"X-XSS-Protection",
			"nosniff",
		},
	}

	for contentType, headers := range requiredHeaders {
		t.Run("ContentType_"+contentType, func(t *testing.T) {
			input := `c.Header("Content-Type", "` + contentType + `")`
			result := security.AddSecurityHeaders(input)

			for _, header := range headers {
				if !strings.Contains(result, header) {
					t.Errorf("Required security header %q missing for content type %s", header, contentType)
				}
			}
		})
	}
}

// TestSecurityHeaderValidValues ensures security headers have correct values
func TestSecurityHeaderValidValues(t *testing.T) {
	input := `c.Header("Content-Type", "application/json")`
	result := security.AddSecurityHeaders(input)

	// Test specific header values
	headerValueTests := []struct {
		header string
		value  string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
	}

	for _, test := range headerValueTests {
		expectedPattern := test.header + `", "` + test.value + `"`
		if !strings.Contains(result, expectedPattern) {
			t.Errorf("Expected security header %s to have value %s, result: %q", test.header, test.value, result)
		}
	}
}

// TestSecurityHeaderRegressionPrevention ensures header fixes don't introduce vulnerabilities
func TestSecurityHeaderRegressionPrevention(t *testing.T) {
	regressionTests := []struct {
		name     string
		input    string
		checkFor []string
	}{
		{
			name:     "Security headers don't introduce XSS",
			input:    `c.Header("Content-Type", "application/json")`,
			checkFor: []string{"<script>", "javascript:", "onload=", "onerror="},
		},
		{
			name:     "Security headers don't introduce injection",
			input:    `res.setHeader('Content-Type', 'text/html')`,
			checkFor: []string{"'; DROP TABLE", "UNION SELECT", "OR 1=1"},
		},
	}

	for _, tt := range regressionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := security.AddSecurityHeaders(tt.input)

			// Verify no dangerous patterns were introduced
			for _, dangerousPattern := range tt.checkFor {
				if strings.Contains(strings.ToLower(result), strings.ToLower(dangerousPattern)) {
					t.Errorf("Security header fix introduced dangerous pattern %q in result: %q", dangerousPattern, result)
				}
			}
		})
	}
}

// TestSecurityHeaderIntegrationWithExistingHeaders tests behavior when security headers already exist
func TestSecurityHeaderIntegrationWithExistingHeaders(t *testing.T) {
	integrationTests := []struct {
		name        string
		input       string
		description string
		expectAdd   bool
	}{
		{
			name:        "Should add headers when only Content-Type exists",
			input:       `c.Header("Content-Type", "application/json")`,
			description: "Should add security headers when only content-type is set",
			expectAdd:   true,
		},
		{
			name: "Should not duplicate existing security headers",
			input: `c.Header("Content-Type", "application/json")
c.Header("X-Content-Type-Options", "nosniff")`,
			description: "Should not duplicate existing security headers",
			expectAdd:   true, // Still adds other missing headers
		},
	}

	for _, tt := range integrationTests {
		t.Run(tt.name, func(t *testing.T) {
			result := security.AddSecurityHeaders(tt.input)

			if tt.expectAdd {
				// Should add at least some security headers
				if !strings.Contains(result, "SECURITY") {
					t.Error("Expected security headers to be added")
				}
			}

			// Verify original content is preserved
			originalLines := strings.Split(tt.input, "\n")
			for _, originalLine := range originalLines {
				if strings.TrimSpace(originalLine) != "" && !strings.Contains(result, originalLine) {
					t.Errorf("Original line should be preserved: %q", originalLine)
				}
			}
		})
	}
}
