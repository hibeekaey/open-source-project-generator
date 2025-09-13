//go:build !ci

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAuthenticationSecurityFixesValidation provides comprehensive tests for authentication security improvements
func TestAuthenticationSecurityFixesValidation(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		fixFunction        func(string) string
		expectedFixed      bool
		expectedContent    []string
		notExpectedContent []string
		authType           string
		description        string
	}{
		{
			name:               "JWT none algorithm vulnerability fix",
			input:              `  algorithm: "none"`,
			fixFunction:        fixJWTNoneAlgorithm,
			expectedFixed:      true,
			expectedContent:    []string{"HS256", "SECURITY FIX", "secure HS256 algorithm"},
			notExpectedContent: []string{`"none"`},
			authType:           "jwt_none",
			description:        "Should replace none algorithm with secure HS256",
		},
		{
			name:               "JWT signing method none vulnerability",
			input:              `token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)`,
			fixFunction:        fixJWTNoneAlgorithm,
			expectedFixed:      false, // This pattern isn't handled by fixJWTNoneAlgorithm
			expectedContent:    []string{},
			notExpectedContent: []string{},
			authType:           "jwt_method_none",
			description:        "JWT SigningMethodNone should be detected but may need different fix",
		},
		{
			name:               "JWT token signing should get expiration guidance",
			input:              `token := jwt.Sign(claims, secret)`,
			fixFunction:        addJWTExpiration,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY", "expiration"},
			notExpectedContent: []string{},
			authType:           "jwt_no_expiration",
			description:        "Should add expiration guidance for JWT signing",
		},
		{
			name:               "Node.js JWT signing should get expiration guidance",
			input:              `const token = jwt.sign(payload, secret)`,
			fixFunction:        addJWTExpiration,
			expectedFixed:      true,
			expectedContent:    []string{"expiresIn", "15m"},
			notExpectedContent: []string{},
			authType:           "nodejs_jwt",
			description:        "Should add expiration guidance for Node.js JWT",
		},
		{
			name:               "Cookie configuration should get security flags",
			input:              `http.SetCookie(w, &http.Cookie{Name: "session", Value: sessionID})`,
			fixFunction:        addSecureCookieFlags,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY", "HttpOnly", "Secure"},
			notExpectedContent: []string{},
			authType:           "insecure_cookie",
			description:        "Should add secure cookie flags guidance",
		},
		{
			name:               "Safe JWT algorithm should not be modified",
			input:              `algorithm: "HS256"`,
			fixFunction:        fixJWTNoneAlgorithm,
			expectedFixed:      false,
			expectedContent:    []string{},
			notExpectedContent: []string{"SECURITY FIX"},
			authType:           "safe_jwt",
			description:        "Safe JWT algorithms should remain unchanged",
		},
		{
			name:               "JWT with expiration gets additional guidance",
			input:              `jwt.sign(payload, secret, { expiresIn: '15m' })`,
			fixFunction:        addJWTExpiration,
			expectedFixed:      true,
			expectedContent:    []string{"SECURITY", "expiresIn"},
			notExpectedContent: []string{},
			authType:           "jwt_with_expiration",
			description:        "JWT with expiration gets additional security guidance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			// Verify if fix was applied as expected
			wasFixed := result != tt.input
			if wasFixed != tt.expectedFixed {
				if tt.expectedFixed {
					t.Errorf("Expected authentication fix to be applied, but input remained unchanged")
				} else {
					t.Errorf("Expected no authentication fix for safe configuration, but got: %q", result)
				}
			}

			// Verify expected content is present
			for _, expected := range tt.expectedContent {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected authentication fix result to contain %q, got %q", expected, result)
				}
			}

			// Verify unwanted content is not present
			for _, notExpected := range tt.notExpectedContent {
				if strings.Contains(result, notExpected) {
					t.Errorf("Authentication fix result should not contain %q, got %q", notExpected, result)
				}
			}
		})
	}
}

// TestAuthenticationVulnerabilityDetection tests detection of authentication vulnerabilities
func TestAuthenticationVulnerabilityDetection(t *testing.T) {
	scanner := NewScanner()

	vulnerabilityTests := []struct {
		name         string
		code         string
		shouldDetect bool
		expectedType SecurityIssueType
		severity     SeverityLevel
	}{
		{
			name:         "JWT none algorithm should be detected as critical",
			code:         `token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityCritical,
		},
		{
			name:         "JWT algorithm none in config should be detected as critical",
			code:         `algorithm: "none"`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityCritical,
		},
		{
			name:         "Weak JWT secret should be detected as high",
			code:         `jwt.sign(payload, "secret")`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityHigh,
		},
		{
			name:         "JWT signing without expiration should be detected as low",
			code:         `token := jwt.Sign(claims, secretKey)`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityLow,
		},
		{
			name:         "Cookie configuration should be detected as low",
			code:         `http.SetCookie(w, cookie)`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityLow,
		},
		{
			name:         "Hardcoded password should be detected as high",
			code:         `password := "mypassword123"`,
			shouldDetect: true,
			expectedType: WeakAuthentication,
			severity:     SeverityHigh,
		},
		{
			name:         "Safe JWT algorithm should not be detected",
			code:         `token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`,
			shouldDetect: false,
		},
		{
			name:         "Environment-based secret should not be detected",
			code:         `jwt.sign(payload, process.env.JWT_SECRET)`,
			shouldDetect: false,
		},
	}

	for _, tt := range vulnerabilityTests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "auth_test.go.tmpl")

			err := os.WriteFile(testFile, []byte(tt.code), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Authentication vulnerability scan failed: %v", err)
			}

			authIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == WeakAuthentication {
					authIssueFound = true
					if tt.shouldDetect && issue.Severity != tt.severity {
						t.Errorf("Expected authentication vulnerability severity %s, got %s", tt.severity, issue.Severity)
					}
					break
				}
			}

			if tt.shouldDetect && !authIssueFound {
				t.Errorf("Expected authentication vulnerability to be detected, but none found")
			}

			if !tt.shouldDetect && authIssueFound {
				t.Errorf("Authentication vulnerability should not be detected for safe code")
			}
		})
	}
}

// TestJWTSecurityImprovements tests specific JWT security enhancements
func TestJWTSecurityImprovements(t *testing.T) {
	jwtTests := []struct {
		name        string
		input       string
		fixType     string
		expected    []string
		notExpected []string
	}{
		{
			name:        "JWT none algorithm replacement",
			input:       `{ algorithm: "none" }`,
			fixType:     "algorithm_fix",
			expected:    []string{"HS256", "SECURITY FIX"},
			notExpected: []string{`"none"`},
		},
		{
			name:        "JWT expiration addition for Go",
			input:       `token := jwt.Sign(claims, secret)`,
			fixType:     "expiration_fix",
			expected:    []string{"SECURITY", "expiration", "time.Now().Add(15 * time.Minute).Unix()"},
			notExpected: []string{},
		},
		{
			name:        "JWT expiration addition for Node.js",
			input:       `jwt.sign(payload, secret)`,
			fixType:     "expiration_fix",
			expected:    []string{"expiresIn", "15m"},
			notExpected: []string{},
		},
	}

	for _, tt := range jwtTests {
		t.Run(tt.name, func(t *testing.T) {
			var result string

			switch tt.fixType {
			case "algorithm_fix":
				result = fixJWTNoneAlgorithm(tt.input)
			case "expiration_fix":
				result = addJWTExpiration(tt.input)
			default:
				result = tt.input
			}

			// Verify expected improvements
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected JWT security improvement to contain %q, got %q", expected, result)
				}
			}

			// Verify unwanted content is removed/not present
			for _, notExpected := range tt.notExpected {
				if strings.Contains(result, notExpected) {
					t.Errorf("JWT security improvement should not contain %q, got %q", notExpected, result)
				}
			}
		})
	}
}

// TestCookieSecurityImprovements tests cookie security enhancements
func TestCookieSecurityImprovements(t *testing.T) {
	cookieTests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Go HTTP cookie security",
			input:    `http.SetCookie(w, &http.Cookie{Name: "session", Value: id})`,
			expected: []string{"SECURITY", "HttpOnly", "Secure"},
		},
		{
			name:     "Generic cookie configuration",
			input:    `cookie := &Cookie{Name: "auth", Value: token}`,
			expected: []string{"SECURITY", "HttpOnly", "Secure"},
		},
		{
			name:     "Express cookie setting",
			input:    `res.cookie('session', sessionId)`,
			expected: []string{"SECURITY", "HttpOnly", "Secure"},
		},
	}

	for _, tt := range cookieTests {
		t.Run(tt.name, func(t *testing.T) {
			result := addSecureCookieFlags(tt.input)

			// Verify security guidance was added
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected cookie security improvement to contain %q, got %q", expected, result)
				}
			}
		})
	}
}

// TestAuthenticationFixIdempotency ensures authentication fixes can be applied multiple times safely
func TestAuthenticationFixIdempotency(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		fixFunction func(string) string
	}{
		{
			name:        "JWT none algorithm fix idempotency",
			input:       `algorithm: "none"`,
			fixFunction: fixJWTNoneAlgorithm,
		},
		{
			name:        "JWT expiration fix idempotency",
			input:       `jwt.sign(payload, secret)`,
			fixFunction: addJWTExpiration,
		},
		{
			name:        "Cookie security fix idempotency",
			input:       `http.SetCookie(w, cookie)`,
			fixFunction: addSecureCookieFlags,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Apply fix first time
			result1 := tc.fixFunction(tc.input)

			// Apply fix second time
			result2 := tc.fixFunction(result1)

			// Verify second application doesn't change the result
			if result1 != result2 {
				t.Errorf("Authentication fix should be idempotent, but second application changed result")
				t.Logf("First result:  %q", result1)
				t.Logf("Second result: %q", result2)
			}
		})
	}
}

// TestAuthenticationSecurityRegressionPrevention ensures auth fixes don't introduce new vulnerabilities
func TestAuthenticationSecurityRegressionPrevention(t *testing.T) {
	regressionTests := []struct {
		name        string
		input       string
		fixFunction func(string) string
		checkFor    []string
	}{
		{
			name:        "JWT fix doesn't introduce injection",
			input:       `algorithm: "none"`,
			fixFunction: fixJWTNoneAlgorithm,
			checkFor:    []string{"'; DROP", "UNION SELECT", "<script>"},
		},
		{
			name:        "Cookie fix doesn't introduce XSS",
			input:       `http.SetCookie(w, cookie)`,
			fixFunction: addSecureCookieFlags,
			checkFor:    []string{"<script>", "javascript:", "onload="},
		},
		{
			name:        "JWT expiration fix doesn't introduce dangerous patterns",
			input:       `jwt.sign(payload, secret)`,
			fixFunction: addJWTExpiration,
			checkFor:    []string{"eval(", "innerHTML", "document.write"},
		},
	}

	for _, tt := range regressionTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			// Verify fix doesn't introduce dangerous patterns
			for _, dangerousPattern := range tt.checkFor {
				if strings.Contains(strings.ToLower(result), strings.ToLower(dangerousPattern)) {
					t.Errorf("Authentication fix introduced dangerous pattern %q in result: %q", dangerousPattern, result)
				}
			}
		})
	}
}

// TestAuthenticationFixPreservesContext ensures fixes maintain code context and functionality
func TestAuthenticationFixPreservesContext(t *testing.T) {
	contextTests := []struct {
		name           string
		input          string
		fixFunction    func(string) string
		shouldPreserve []string
	}{
		{
			name:           "JWT algorithm fix preserves variable names",
			input:          `const jwtConfig = { algorithm: "none", issuer: "myapp" }`,
			fixFunction:    fixJWTNoneAlgorithm,
			shouldPreserve: []string{"jwtConfig", "issuer", "myapp"},
		},
		{
			name:           "Cookie fix preserves cookie properties",
			input:          `http.SetCookie(w, &http.Cookie{Name: "session", Value: sessionID, Path: "/"})`,
			fixFunction:    addSecureCookieFlags,
			shouldPreserve: []string{"session", "sessionID", "Path"},
		},
		{
			name:           "JWT signing fix preserves function call structure",
			input:          `token, err := jwt.Sign(userClaims, secretKey)`,
			fixFunction:    addJWTExpiration,
			shouldPreserve: []string{"token", "err", "userClaims", "secretKey"},
		},
	}

	for _, tt := range contextTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			// Verify important context is preserved
			for _, preserve := range tt.shouldPreserve {
				if !strings.Contains(result, preserve) {
					t.Errorf("Authentication fix should preserve %q, got result: %q", preserve, result)
				}
			}
		})
	}
}

// TestAuthenticationSecurityBestPractices validates that fixes follow security best practices
func TestAuthenticationSecurityBestPractices(t *testing.T) {
	bestPracticeTests := []struct {
		name        string
		input       string
		fixFunction func(string) string
		practices   []string
	}{
		{
			name:        "JWT algorithm fix uses secure algorithm",
			input:       `algorithm: "none"`,
			fixFunction: fixJWTNoneAlgorithm,
			practices:   []string{"HS256"}, // Should use secure algorithm
		},
		{
			name:        "JWT expiration follows time limits",
			input:       `jwt.sign(payload, secret)`,
			fixFunction: addJWTExpiration,
			practices:   []string{"15m"}, // Should suggest reasonable expiration
		},
		{
			name:        "Cookie security includes essential flags",
			input:       `http.SetCookie(w, cookie)`,
			fixFunction: addSecureCookieFlags,
			practices:   []string{"HttpOnly", "Secure"}, // Should mention both flags
		},
	}

	for _, tt := range bestPracticeTests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fixFunction(tt.input)

			// Verify security best practices are followed
			for _, practice := range tt.practices {
				if !strings.Contains(result, practice) {
					t.Errorf("Authentication fix should follow best practice %q, got result: %q", practice, result)
				}
			}
		})
	}
}
