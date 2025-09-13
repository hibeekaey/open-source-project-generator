package security

import (
	"testing"
)

func TestSecurityPatterns(t *testing.T) {
	patterns := getSecurityPatterns()

	testCases := []struct {
		name        string
		code        string
		shouldMatch bool
		issueType   SecurityIssueType
		severity    SeverityLevel
	}{
		// CORS Vulnerabilities
		{
			name:        "CORS Null Origin",
			code:        `c.Header("Access-Control-Allow-Origin", "null")`,
			shouldMatch: true,
			issueType:   CORSVulnerability,
			severity:    SeverityCritical,
		},
		{
			name:        "CORS Wildcard",
			code:        `c.Header("Access-Control-Allow-Origin", "*")`,
			shouldMatch: true,
			issueType:   CORSVulnerability,
			severity:    SeverityMedium,
		},
		{
			name:        "Safe CORS",
			code:        `c.Header("Access-Control-Allow-Origin", "https://example.com")`,
			shouldMatch: false,
		},

		// JWT Authentication Issues
		{
			name:        "JWT None Algorithm",
			code:        `token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)`,
			shouldMatch: true,
			issueType:   WeakAuthentication,
			severity:    SeverityCritical,
		},
		{
			name:        "Weak JWT Secret",
			code:        `token := jwt.Sign(claims, "secret")`,
			shouldMatch: true,
			issueType:   WeakAuthentication,
			severity:    SeverityHigh,
		},
		{
			name:        "JWT Signing",
			code:        `token := jwt.sign(payload, secret)`,
			shouldMatch: true,
			issueType:   WeakAuthentication,
			severity:    SeverityLow,
		},

		// SQL Injection Risks
		{
			name:        "SQL String Concatenation",
			code:        `query := "SELECT * FROM users WHERE id = " + userID`,
			shouldMatch: true,
			issueType:   SQLInjectionRisk,
			severity:    SeverityCritical,
		},
		{
			name:        "SQL Variable Interpolation",
			code:        `query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)`,
			shouldMatch: true,
			issueType:   SQLInjectionRisk,
			severity:    SeverityHigh,
		},
		{
			name:        "Safe SQL Query",
			code:        `query := "SELECT * FROM users WHERE id = $1"`,
			shouldMatch: false,
		},

		// Information Leakage
		{
			name:        "Detailed Error Message",
			code:        `return fmt.Errorf("database error: %v", err)`,
			shouldMatch: true,
			issueType:   InformationLeakage,
			severity:    SeverityMedium,
		},
		{
			name:        "Debug Enabled",
			code:        `debug: true`,
			shouldMatch: true,
			issueType:   InformationLeakage,
			severity:    SeverityMedium,
		},
		{
			name:        "Safe Error Handling",
			code:        `return errors.New("internal server error")`,
			shouldMatch: false,
		},

		// Security Headers
		{
			name:        "Content-Type Header",
			code:        `c.Header("Content-Type", "application/json")`,
			shouldMatch: true,
			issueType:   MissingSecurityHeader,
			severity:    SeverityLow,
		},
		{
			name:        "HTTP Header Setting",
			code:        `response.Header("X-Custom", "value")`,
			shouldMatch: true,
			issueType:   MissingSecurityHeader,
			severity:    SeverityLow,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched := false
			var matchedPattern SecurityPattern

			for _, pattern := range patterns {
				if pattern.Pattern.MatchString(tc.code) {
					matched = true
					matchedPattern = pattern
					break
				}
			}

			if matched != tc.shouldMatch {
				if tc.shouldMatch {
					t.Errorf("Expected pattern to match code: %s", tc.code)
				} else {
					t.Errorf("Pattern should not match code: %s", tc.code)
				}
				return
			}

			if tc.shouldMatch {
				if matchedPattern.IssueType != tc.issueType {
					t.Errorf("Expected issue type %s, got %s", tc.issueType, matchedPattern.IssueType)
				}
				if matchedPattern.Severity != tc.severity {
					t.Errorf("Expected severity %s, got %s", tc.severity, matchedPattern.Severity)
				}
			}
		})
	}
}

func TestPatternCoverage(t *testing.T) {
	patterns := getSecurityPatterns()

	// Ensure we have patterns for all major security issue types
	issueTypes := make(map[SecurityIssueType]bool)
	for _, pattern := range patterns {
		issueTypes[pattern.IssueType] = true
	}

	expectedTypes := []SecurityIssueType{
		CORSVulnerability,
		MissingSecurityHeader,
		WeakAuthentication,
		SQLInjectionRisk,
		InformationLeakage,
	}

	for _, expectedType := range expectedTypes {
		if !issueTypes[expectedType] {
			t.Errorf("Missing patterns for issue type: %s", expectedType)
		}
	}
}

func TestPatternQuality(t *testing.T) {
	patterns := getSecurityPatterns()

	for _, pattern := range patterns {
		// Ensure all patterns have required fields
		if pattern.Name == "" {
			t.Error("Pattern missing name")
		}
		if pattern.Pattern == nil {
			t.Error("Pattern missing regex")
		}
		if pattern.Description == "" {
			t.Error("Pattern missing description")
		}
		if pattern.Recommendation == "" {
			t.Error("Pattern missing recommendation")
		}

		// Ensure severity levels are valid
		validSeverities := map[SeverityLevel]bool{
			SeverityCritical: true,
			SeverityHigh:     true,
			SeverityMedium:   true,
			SeverityLow:      true,
		}

		if !validSeverities[pattern.Severity] {
			t.Errorf("Invalid severity level: %s", pattern.Severity)
		}
	}
}
