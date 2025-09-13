//go:build !ci

package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCORSSecurityFixes tests all CORS-related security fixes
func TestCORSSecurityFixes(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedFix    bool
		expectedOutput []string
		framework      string
	}{
		{
			name:        "CORS null origin - Go Gin",
			input:       `    c.Header("Access-Control-Allow-Origin", "null")`,
			expectedFix: true,
			expectedOutput: []string{
				"SECURITY FIX",
				"Removed Access-Control-Allow-Origin: null header",
				"omit the header entirely",
			},
			framework: "gin",
		},
		{
			name:        "CORS null origin - Node.js Express",
			input:       `  res.setHeader('Access-Control-Allow-Origin', 'null');`,
			expectedFix: true,
			expectedOutput: []string{
				"SECURITY FIX",
				"Removed Access-Control-Allow-Origin: null header",
				"omit the header entirely",
			},
			framework: "express",
		},
		{
			name:        "CORS wildcard - Go Gin",
			input:       `    c.Header("Access-Control-Allow-Origin", "*")`,
			expectedFix: true,
			expectedOutput: []string{
				"SECURITY FIX",
				"isAllowedOrigin",
				"c.Header",
			},
			framework: "gin",
		},
		{
			name:        "CORS wildcard - Node.js Express",
			input:       `  res.setHeader('Access-Control-Allow-Origin', '*');`,
			expectedFix: true,
			expectedOutput: []string{
				"SECURITY FIX",
				"isAllowedOrigin",
				"res.setHeader",
			},
			framework: "express",
		},
		{
			name:        "Safe CORS configuration",
			input:       `    c.Header("Access-Control-Allow-Origin", "https://example.com")`,
			expectedFix: false,
			framework:   "gin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string

			// Apply appropriate fix function based on the vulnerability type
			if strings.Contains(tt.input, "null") {
				result = fixCORSNullOrigin(tt.input)
			} else if strings.Contains(tt.input, "*") {
				result = fixCORSWildcard(tt.input)
			} else {
				result = tt.input
			}

			if tt.expectedFix {
				// Verify the fix was applied
				if result == tt.input {
					t.Errorf("Expected fix to be applied, but input remained unchanged")
				}

				// Verify expected content is present
				for _, expected := range tt.expectedOutput {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected output to contain %q, got %q", expected, result)
					}
				}
			} else {
				// Verify no fix was applied for safe configurations
				if result != tt.input {
					t.Errorf("Expected no fix for safe configuration, but got %q", result)
				}
			}
		})
	}
}

// TestSecurityHeaderImplementation tests security header fixes
func TestSecurityHeaderImplementation(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedHeaders []string
		framework       string
	}{
		{
			name:  "Go Gin content-type header",
			input: `    c.Header("Content-Type", "application/json")`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
				"1; mode=block",
			},
			framework: "gin",
		},
		{
			name:  "Node.js Express content-type header",
			input: `  res.setHeader('Content-Type', 'application/json');`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
				"1; mode=block",
			},
			framework: "express",
		},
		{
			name:  "HTML content-type header",
			input: `    c.Header("Content-Type", "text/html")`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"nosniff",
			},
			framework: "gin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddSecurityHeaders(tt.input)

			// Verify original header is preserved
			if !strings.Contains(result, tt.input) {
				t.Errorf("Original header should be preserved in result")
			}

			// Verify security headers were added
			for _, header := range tt.expectedHeaders {
				if !strings.Contains(result, header) {
					t.Errorf("Expected security header %q to be added, got %q", header, result)
				}
			}

			// Verify proper indentation is maintained
			lines := strings.Split(result, "\n")
			if len(lines) > 1 {
				originalIndent := getIndentation(tt.input)
				for i := 1; i < len(lines); i++ {
					if strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(lines[i], originalIndent) {
						t.Errorf("Security header line should maintain original indentation: %q", lines[i])
					}
				}
			}
		})
	}
}

// TestAuthenticationSecurityImprovements tests authentication-related security fixes
func TestAuthenticationSecurityImprovements(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		fixFunction    func(string) string
		expectedOutput []string
		shouldChange   bool
	}{
		{
			name:        "JWT none algorithm fix",
			input:       `  algorithm: "none"`,
			fixFunction: fixJWTNoneAlgorithm,
			expectedOutput: []string{
				"HS256",
				"SECURITY FIX",
				"secure HS256 algorithm",
			},
			shouldChange: true,
		},
		{
			name:        "JWT signing with expiration",
			input:       `token := jwt.Sign(claims, secret)`,
			fixFunction: addJWTExpiration,
			expectedOutput: []string{
				"SECURITY",
				"expiration",
			},
			shouldChange: true,
		},
		{
			name:        "Cookie security flags",
			input:       `http.SetCookie(w, &http.Cookie{Name: "session", Value: sessionID})`,
			fixFunction: addSecureCookieFlags,
			expectedOutput: []string{
				"SECURITY",
				"HttpOnly",
				"Secure",
			},
			shouldChange: true,
		},
		{
			name:         "Safe JWT configuration",
			input:        `algorithm: "HS256"`,
			fixFunction:  fixJWTNoneAlgorithm,
			shouldChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			if tt.shouldChange {
				// Verify the fix was applied
				if result == tt.input {
					t.Errorf("Expected authentication fix to be applied, but input remained unchanged")
				}

				// Verify expected security improvements
				for _, expected := range tt.expectedOutput {
					if !strings.Contains(result, expected) {
						t.Errorf("Expected output to contain %q, got %q", expected, result)
					}
				}
			} else {
				// Verify no unnecessary changes for safe configurations
				if result != tt.input {
					t.Errorf("Expected no changes for safe configuration, but got %q", result)
				}
			}
		})
	}
}

// TestSecurityRegressionPrevention tests that security fixes don't introduce new vulnerabilities
func TestSecurityRegressionPrevention(t *testing.T) {
	// Test cases that should never be "fixed" as they are already secure
	safeCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Secure CORS with specific origin",
			input: `c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")`,
		},
		{
			name:  "Parameterized SQL query",
			input: `db.Query("SELECT * FROM users WHERE id = $1", userID)`,
		},
		{
			name:  "Secure JWT algorithm",
			input: `token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`,
		},
		{
			name:  "Generic error message",
			input: `return errors.New("operation failed")`,
		},
		{
			name:  "Environment-based debug config",
			input: `debug: os.Getenv("DEBUG") == "true"`,
		},
	}

	fixer := NewFixer()

	for _, tc := range safeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file with safe code
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "safe.go.tmpl")

			err := os.WriteFile(testFile, []byte(tc.input), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Apply security fixes
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "all",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("FixFile failed: %v", err)
			}

			// Verify no fixes were applied to safe code
			if len(result.FixedIssues) > 0 {
				t.Errorf("Safe code should not be modified, but %d fixes were applied", len(result.FixedIssues))
				for _, fix := range result.FixedIssues {
					t.Logf("Unexpected fix: %s - %s", fix.IssueType, fix.Description)
				}
			}

			// Verify file content remains unchanged
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			if string(content) != tc.input {
				t.Errorf("Safe code was modified: expected %q, got %q", tc.input, string(content))
			}
		})
	}
}

// TestIntegratedSecurityValidation tests end-to-end security validation
func TestIntegratedSecurityValidation(t *testing.T) {
	// Create a comprehensive test template with multiple security issues
	testTemplate := `package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CORS vulnerability - should be fixed
		c.Header("Access-Control-Allow-Origin", "null")
		
		// Missing security headers - should be enhanced
		c.Header("Content-Type", "application/json")
		
		c.Next()
	}
}

func authHandler(c *gin.Context) {
	// JWT vulnerability - should be fixed
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"user": "test",
	})
	
	// Missing expiration - should be added
	tokenString, _ := jwt.Sign(token, []byte("secret"))
	
	c.JSON(200, gin.H{"token": tokenString})
}

func getUserByID(c *gin.Context) {
	userID := c.Param("id")
	
	// SQL injection vulnerability - should be fixed
	query := "SELECT * FROM users WHERE id = " + userID
	
	// Information leakage - should be fixed
	if err := db.Query(query); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("database error: %v", err)})
		return
	}
}

func safeHandler(c *gin.Context) {
	// This is already secure - should not be modified
	c.Header("Access-Control-Allow-Origin", "https://trusted-domain.com")
	
	// Parameterized query - should not be modified
	db.Query("SELECT * FROM users WHERE id = $1", userID)
	
	// Generic error - should not be modified
	c.JSON(500, gin.H{"error": "operation failed"})
}`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "comprehensive.go.tmpl")

	err := os.WriteFile(testFile, []byte(testTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// First, scan for security issues
	scanner := NewScanner()
	report, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Security scan failed: %v", err)
	}

	// Verify security issues were detected
	if len(report) == 0 {
		t.Error("Expected security issues to be detected")
	}

	// Verify we have issues of different types
	issueTypes := make(map[SecurityIssueType]int)
	for _, issue := range report {
		issueTypes[issue.IssueType]++
	}

	expectedTypes := []SecurityIssueType{
		CORSVulnerability,
		MissingSecurityHeader,
		WeakAuthentication,
		SQLInjectionRisk,
		InformationLeakage,
	}

	for _, expectedType := range expectedTypes {
		if issueTypes[expectedType] == 0 {
			t.Errorf("Expected to find %s issues, but none were detected", expectedType)
		}
	}

	// Apply security fixes
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      true,
		FixType:      "all",
		CreateBackup: true,
	}

	fixResult, err := fixer.FixFile(testFile, options)
	if err != nil {
		t.Fatalf("Security fix failed: %v", err)
	}

	// Verify fixes were applied
	if len(fixResult.FixedIssues) == 0 {
		t.Error("Expected security fixes to be applied")
	}

	// Verify backup was created
	if fixResult.BackupsCreated != 1 {
		t.Errorf("Expected 1 backup to be created, got %d", fixResult.BackupsCreated)
	}

	// Read the fixed content
	fixedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read fixed file: %v", err)
	}

	fixedStr := string(fixedContent)

	// Verify specific fixes were applied
	securityChecks := []struct {
		name        string
		shouldExist bool
		pattern     string
	}{
		{"CORS null origin removed", false, `"null"`},
		{"Security headers added", true, "X-Content-Type-Options"},
		{"JWT none algorithm fixed", false, "SigningMethodNone"},
		{"Security comments added", true, "SECURITY FIX"},
		{"Safe code preserved", true, "https://trusted-domain.com"},
		{"Parameterized query preserved", true, "$1"},
	}

	for _, check := range securityChecks {
		exists := strings.Contains(fixedStr, check.pattern)
		if exists != check.shouldExist {
			if check.shouldExist {
				t.Errorf("%s: Expected pattern %q to exist in fixed content", check.name, check.pattern)
			} else {
				t.Errorf("%s: Pattern %q should have been removed from fixed content", check.name, check.pattern)
			}
		}
	}

	// Re-scan the fixed file to verify issues were resolved
	newReport, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Post-fix security scan failed: %v", err)
	}

	// Verify critical issues were resolved
	criticalIssues := 0
	for _, issue := range newReport {
		if issue.Severity == SeverityCritical {
			criticalIssues++
		}
	}

	if criticalIssues > 0 {
		t.Errorf("Expected critical security issues to be resolved, but %d remain", criticalIssues)
	}
}

// TestSecurityFixIdempotency ensures that applying fixes multiple times doesn't cause issues
func TestSecurityFixIdempotency(t *testing.T) {
	testContent := `c.Header("Access-Control-Allow-Origin", "null")`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "idempotent.go.tmpl")

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      false,
		FixType:      "all",
		CreateBackup: false,
	}

	// Apply fixes first time
	result1, err := fixer.FixFile(testFile, options)
	if err != nil {
		t.Fatalf("First fix application failed: %v", err)
	}

	if len(result1.FixedIssues) == 0 {
		t.Error("Expected fixes to be applied on first run")
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
		t.Errorf("Expected no fixes on second run, but %d were applied", len(result2.FixedIssues))
	}

	// Read content after second fix
	content2, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file after second fix: %v", err)
	}

	// Verify content is identical
	if string(content1) != string(content2) {
		t.Error("File content should be identical after second fix application")
	}
}

// TestSecurityValidationPerformance tests that security validation performs well on large files
func TestSecurityValidationPerformance(t *testing.T) {
	// Create a large template file with mixed secure and insecure patterns
	var content strings.Builder

	// Add file header
	content.WriteString("package main\n\n")

	// Add many lines with mixed security patterns
	for i := 0; i < 1000; i++ {
		if i%10 == 0 {
			// Add some security issues
			content.WriteString(fmt.Sprintf(`// Line %d with CORS issue
c.Header("Access-Control-Allow-Origin", "null")
`, i))
		} else if i%15 == 0 {
			// Add some SQL injection risks
			content.WriteString(fmt.Sprintf(`// Line %d with SQL issue
query := "SELECT * FROM table WHERE id = " + userID
`, i))
		} else {
			// Add safe code
			content.WriteString(fmt.Sprintf(`// Line %d - safe code
c.Header("X-Custom-Header", "safe-value")
`, i))
		}
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.go.tmpl")

	err := os.WriteFile(testFile, []byte(content.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Test scanning performance
	scanner := NewScanner()
	report, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("Scanning large file failed: %v", err)
	}

	// Verify issues were found
	if len(report) == 0 {
		t.Error("Expected security issues to be found in large file")
	}

	// Test fixing performance
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      false,
		FixType:      "all",
		CreateBackup: false,
	}

	result, err := fixer.FixFile(testFile, options)
	if err != nil {
		t.Fatalf("Fixing large file failed: %v", err)
	}

	// Verify fixes were applied
	if len(result.FixedIssues) == 0 {
		t.Error("Expected fixes to be applied to large file")
	}

	t.Logf("Performance test completed: scanned and fixed %d issues in large file", len(result.FixedIssues))
}
