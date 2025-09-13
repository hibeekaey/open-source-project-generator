package security

import "regexp"

// getSecurityPatterns returns all predefined security patterns
func getSecurityPatterns() []SecurityPattern {
	return []SecurityPattern{
		// CORS Vulnerabilities
		{
			Name:           "CORS Null Origin",
			Pattern:        regexp.MustCompile(`Access-Control-Allow-Origin.*["']null["']`),
			IssueType:      CORSVulnerability,
			Severity:       SeverityCritical,
			Description:    "Setting Access-Control-Allow-Origin to 'null' can allow bypass attacks",
			Recommendation: "Omit the Access-Control-Allow-Origin header entirely for disallowed origins instead of setting it to 'null'",
			FixAvailable:   true,
		},
		{
			Name:           "CORS Wildcard with Credentials",
			Pattern:        regexp.MustCompile(`Access-Control-Allow-Origin.*\*.*Access-Control-Allow-Credentials.*true|Access-Control-Allow-Credentials.*true.*Access-Control-Allow-Origin.*\*`),
			IssueType:      CORSVulnerability,
			Severity:       SeverityHigh,
			Description:    "Using wildcard (*) for Access-Control-Allow-Origin with credentials enabled is insecure",
			Recommendation: "Use specific origins instead of wildcard when credentials are allowed",
			FixAvailable:   true,
		},
		{
			Name:           "CORS Overly Permissive",
			Pattern:        regexp.MustCompile(`Access-Control-Allow-Origin.*["']\*["']`),
			IssueType:      CORSVulnerability,
			Severity:       SeverityMedium,
			Description:    "Wildcard CORS policy allows requests from any origin",
			Recommendation: "Specify explicit allowed origins instead of using wildcard",
			FixAvailable:   true,
		},

		// Missing Security Headers
		{
			Name:           "Content-Type Header Set",
			Pattern:        regexp.MustCompile(`(?i)\.Header\(["']Content-Type["'].*["'](?:application/json|text/html)["']`),
			IssueType:      MissingSecurityHeader,
			Severity:       SeverityLow,
			Description:    "Content-Type header set - ensure X-Content-Type-Options: nosniff is also set",
			Recommendation: "Add X-Content-Type-Options: nosniff header to prevent MIME type sniffing",
			FixAvailable:   true,
		},
		{
			Name:           "HTTP Response Headers",
			Pattern:        regexp.MustCompile(`(?i)response\.Header\(|\.setHeader\(`),
			IssueType:      MissingSecurityHeader,
			Severity:       SeverityLow,
			Description:    "HTTP headers being set - ensure security headers are included",
			Recommendation: "Add security headers: X-Frame-Options, X-XSS-Protection, X-Content-Type-Options",
			FixAvailable:   true,
		},

		// Authentication Issues
		{
			Name:           "JWT None Algorithm",
			Pattern:        regexp.MustCompile(`(?i)jwt\.SigningMethodNone|algorithm.*["']none["']`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityCritical,
			Description:    "JWT 'none' algorithm allows token forgery",
			Recommendation: "Use secure algorithms like HS256, RS256, or ES256 for JWT tokens",
			FixAvailable:   true,
		},
		{
			Name:           "Weak JWT Secret",
			Pattern:        regexp.MustCompile(`(?i)jwt.*["'](?:secret|password|123|test)["']`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityHigh,
			Description:    "Weak or default JWT secret detected",
			Recommendation: "Use a strong, randomly generated secret for JWT signing",
			FixAvailable:   false,
		},
		{
			Name:           "JWT Token Signing",
			Pattern:        regexp.MustCompile(`(?i)jwt\.sign|jwt\.Sign`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityLow,
			Description:    "JWT token signing detected - ensure expiration time is set",
			Recommendation: "Set appropriate expiration time for JWT tokens (e.g., 15 minutes for access tokens)",
			FixAvailable:   true,
		},

		// SQL Injection Risks
		{
			Name:           "String Concatenation in SQL",
			Pattern:        regexp.MustCompile(`(?i)["'].*(?:select|insert|update|delete).*["'].*\+|(?:select|insert|update|delete).*\+`),
			IssueType:      SQLInjectionRisk,
			Severity:       SeverityCritical,
			Description:    "SQL query uses string concatenation which may lead to SQL injection",
			Recommendation: "Use parameterized queries or prepared statements instead of string concatenation",
			FixAvailable:   true,
		},
		{
			Name:           "Direct Variable in SQL",
			Pattern:        regexp.MustCompile(`(?i)(?:select|insert|update|delete).*\$\{.*\}|(?:select|insert|update|delete).*%s`),
			IssueType:      SQLInjectionRisk,
			Severity:       SeverityHigh,
			Description:    "SQL query directly interpolates variables which may lead to SQL injection",
			Recommendation: "Use parameterized queries with placeholders ($1, $2, etc.) instead of direct variable interpolation",
			FixAvailable:   true,
		},

		// Information Leakage
		{
			Name:           "Detailed Error Messages",
			Pattern:        regexp.MustCompile(`(?i)fmt\.Errorf.*(?:database error|sql error|internal error)`),
			IssueType:      InformationLeakage,
			Severity:       SeverityMedium,
			Description:    "Detailed error messages may leak sensitive information",
			Recommendation: "Use generic error messages in production and log detailed errors securely",
			FixAvailable:   true,
		},
		{
			Name:           "Debug Information Exposure",
			Pattern:        regexp.MustCompile(`(?i)(?:debug|trace|stack).*(?:true|enabled|on)`),
			IssueType:      InformationLeakage,
			Severity:       SeverityMedium,
			Description:    "Debug information may be exposed in production",
			Recommendation: "Disable debug information in production environments",
			FixAvailable:   true,
		},

		// Insecure Configurations
		{
			Name:           "Cookie Configuration",
			Pattern:        regexp.MustCompile(`(?i)cookie|setcookie`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityLow,
			Description:    "Cookie configuration detected - ensure HttpOnly and Secure flags are set",
			Recommendation: "Set HttpOnly and Secure flags on sensitive cookies",
			FixAvailable:   true,
		},

		// Insecure Random Generation Patterns
		{
			Name:           "Timestamp-based Random Generation",
			Pattern:        regexp.MustCompile(`time\.Now\(\)\.UnixNano\(\)`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityCritical,
			Description:    "Timestamp-based random generation is predictable and insecure",
			Recommendation: "Use crypto/rand for cryptographically secure random generation",
			FixAvailable:   true,
		},
		{
			Name:           "Math/rand Usage",
			Pattern:        regexp.MustCompile(`math/rand`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityHigh,
			Description:    "math/rand provides predictable pseudorandom numbers unsuitable for security",
			Recommendation: "Use crypto/rand for security-sensitive random number generation",
			FixAvailable:   true,
		},
		{
			Name:           "Predictable Random Functions",
			Pattern:        regexp.MustCompile(`rand\.(?:Int|Int31|Int63|Intn|Float32|Float64)\(`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityHigh,
			Description:    "Predictable random function from math/rand package",
			Recommendation: "Use crypto/rand.Read() or security.SecureRandom interface",
			FixAvailable:   true,
		},
		{
			Name:           "Timestamp-based ID Generation",
			Pattern:        regexp.MustCompile(`fmt\.Sprintf\([^)]*time\.Now\(\)\.Unix(?:Nano)?\(\)`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityHigh,
			Description:    "ID generation using timestamps is predictable",
			Recommendation: "Use secure random ID generation with crypto/rand",
			FixAvailable:   true,
		},
		{
			Name:           "Predictable Temporary File Names",
			Pattern:        regexp.MustCompile(`\.tmp\..*time\.Now\(\)|time\.Now\(\).*\.tmp`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityCritical,
			Description:    "Temporary file names using timestamps are predictable and vulnerable to race conditions",
			Recommendation: "Use secure random suffixes for temporary file names",
			FixAvailable:   true,
		},

		// Additional Security Patterns
		{
			Name:           "Hardcoded Secrets",
			Pattern:        regexp.MustCompile(`(?i)(?:password|secret|key|token).*=.*["'][^"']{8,}["']`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityHigh,
			Description:    "Potential hardcoded secret or password detected",
			Recommendation: "Use environment variables or secure configuration management for secrets",
			FixAvailable:   false,
		},
		{
			Name:           "Insecure HTTP URLs",
			Pattern:        regexp.MustCompile(`http://[^/\s]+`),
			IssueType:      WeakAuthentication,
			Severity:       SeverityMedium,
			Description:    "Insecure HTTP URL detected - verify if HTTPS should be used",
			Recommendation: "Use HTTPS URLs for all external communications",
			FixAvailable:   true,
		},
		{
			Name:           "Missing Input Validation",
			Pattern:        regexp.MustCompile(`(?i)(?:request\.body|req\.body|c\.bind).*without.*validation`),
			IssueType:      SQLInjectionRisk,
			Severity:       SeverityMedium,
			Description:    "Input processing without explicit validation detected",
			Recommendation: "Implement proper input validation and sanitization",
			FixAvailable:   true,
		},
	}
}
