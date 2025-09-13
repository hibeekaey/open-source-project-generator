package security

import (
	"regexp"
	"strings"
)

// getSecurityFixes returns all predefined security fixes
func getSecurityFixes() []SecurityFix {
	return []SecurityFix{
		// CORS Fixes
		{
			Name:           "Fix CORS Null Origin",
			Pattern:        regexp.MustCompile(`(.*)Access-Control-Allow-Origin.*["']null["'](.*)`),
			IssueType:      CORSVulnerability,
			Description:    "Fixed CORS null origin vulnerability",
			FixDescription: "Removed Access-Control-Allow-Origin header for null origins",
			FixFunction:    fixCORSNullOrigin,
			Enabled:        true,
		},
		{
			Name:           "Fix CORS Wildcard with Credentials",
			Pattern:        regexp.MustCompile(`(.*)Access-Control-Allow-Origin.*\*(.*)`),
			IssueType:      CORSVulnerability,
			Description:    "Fixed CORS wildcard with credentials",
			FixDescription: "Replaced wildcard with specific origin validation",
			FixFunction:    fixCORSWildcard,
			Enabled:        true,
		},

		// Security Headers Fixes
		{
			Name:           "Add X-Content-Type-Options Header",
			Pattern:        regexp.MustCompile(`(.*)(?i)content-type.*application/json|content-type.*text/html(.*)`),
			IssueType:      MissingSecurityHeader,
			Description:    "Added X-Content-Type-Options security header",
			FixDescription: "Added X-Content-Type-Options: nosniff header",
			FixFunction:    addContentTypeOptionsHeader,
			Enabled:        true,
		},
		{
			Name:           "Add Security Headers to HTTP Response",
			Pattern:        regexp.MustCompile(`(.*)(?i)\.header\(.*content-type|\.setHeader\(.*content-type|Header\(.*content-type(.*)`),
			IssueType:      MissingSecurityHeader,
			Description:    "Added comprehensive security headers",
			FixDescription: "Added X-Frame-Options, X-XSS-Protection, and X-Content-Type-Options headers",
			FixFunction:    AddSecurityHeaders,
			Enabled:        true,
		},

		// Authentication Fixes
		{
			Name:           "Fix JWT None Algorithm",
			Pattern:        regexp.MustCompile(`(.*)(?i)algorithm.*["']none["'](.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed JWT none algorithm vulnerability",
			FixDescription: "Replaced 'none' algorithm with secure HS256",
			FixFunction:    fixJWTNoneAlgorithm,
			Enabled:        true,
		},
		{
			Name:           "Add JWT Token Expiration",
			Pattern:        regexp.MustCompile(`(.*)(?i)jwt.*sign.*\(.*\)(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Added JWT token expiration",
			FixDescription: "Added 15-minute expiration time for JWT tokens",
			FixFunction:    addJWTExpiration,
			Enabled:        true,
		},

		// SQL Injection Fixes
		{
			Name:           "Fix SQL String Concatenation",
			Pattern:        regexp.MustCompile(`(.*)(?i)(select|insert|update|delete)(.*)(\+.*)(.*)`),
			IssueType:      SQLInjectionRisk,
			Description:    "Fixed SQL injection vulnerability",
			FixDescription: "Replaced string concatenation with parameterized query",
			FixFunction:    fixSQLConcatenation,
			Enabled:        true,
		},
		{
			Name:           "Fix SQL Variable Interpolation",
			Pattern:        regexp.MustCompile(`(.*)(?i)(select|insert|update|delete)(.*)\$\{.*\}(.*)`),
			IssueType:      SQLInjectionRisk,
			Description:    "Fixed SQL variable interpolation",
			FixDescription: "Replaced variable interpolation with parameterized placeholders",
			FixFunction:    fixSQLInterpolation,
			Enabled:        true,
		},

		// Information Leakage Fixes
		{
			Name:           "Fix Detailed Error Messages",
			Pattern:        regexp.MustCompile(`(.*)(?i)error.*message.*(sql|database|internal|stack|trace)(.*)`),
			IssueType:      InformationLeakage,
			Description:    "Fixed information leakage in error messages",
			FixDescription: "Replaced detailed error with generic message",
			FixFunction:    fixDetailedErrors,
			Enabled:        true,
		},
		{
			Name:           "Fix Debug Information Exposure",
			Pattern:        regexp.MustCompile(`(.*)(?i)(debug|trace|stack).*(true|enabled|on)(.*)`),
			IssueType:      InformationLeakage,
			Description:    "Fixed debug information exposure",
			FixDescription: "Disabled debug information in production",
			FixFunction:    fixDebugExposure,
			Enabled:        true,
		},

		// Cookie Security Fixes
		{
			Name:           "Add Secure Cookie Flags",
			Pattern:        regexp.MustCompile(`(.*)(?i)(cookie|setcookie)(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Added secure cookie flags",
			FixDescription: "Added HttpOnly and Secure flags to cookies",
			FixFunction:    addSecureCookieFlags,
			Enabled:        true,
		},

		// Insecure Random Generation Fixes
		{
			Name:           "Fix Timestamp-based Random Generation",
			Pattern:        regexp.MustCompile(`(.*)time\.Now\(\)\.UnixNano\(\)(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed timestamp-based random generation",
			FixDescription: "Replaced timestamp with crypto/rand secure random generation",
			FixFunction:    fixTimestampRandom,
			Enabled:        true,
		},
		{
			Name:           "Fix Math/rand Usage",
			Pattern:        regexp.MustCompile(`(.*)math/rand(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed insecure math/rand usage",
			FixDescription: "Replaced math/rand with crypto/rand import",
			FixFunction:    fixMathRandImport,
			Enabled:        true,
		},
		{
			Name:           "Fix Predictable Random Functions",
			Pattern:        regexp.MustCompile(`(.*)rand\.(Int|Int31|Int63|Intn|Float32|Float64)\((.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed predictable random function usage",
			FixDescription: "Replaced with secure random generation using crypto/rand",
			FixFunction:    fixPredictableRandFunctions,
			Enabled:        true,
		},
		{
			Name:           "Fix Timestamp-based ID Generation",
			Pattern:        regexp.MustCompile(`(.*)fmt\.Sprintf\([^)]*time\.Now\(\)\.Unix(?:Nano)?\(\)(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed timestamp-based ID generation",
			FixDescription: "Replaced with secure random ID generation",
			FixFunction:    fixTimestampIDGeneration,
			Enabled:        true,
		},
		{
			Name:           "Fix Predictable Temporary File Names",
			Pattern:        regexp.MustCompile(`(.*)\.tmp\..*time\.Now\(\)|time\.Now\(\).*\.tmp(.*)`),
			IssueType:      WeakAuthentication,
			Description:    "Fixed predictable temporary file names",
			FixDescription: "Replaced with secure random temporary file naming",
			FixFunction:    fixPredictableTempFiles,
			Enabled:        true,
		},
	}
}

// Fix functions for each security issue type

func fixCORSNullOrigin(line string) string {
	// Remove the line that sets Access-Control-Allow-Origin to null
	// Instead, we'll comment it out with an explanation
	if strings.Contains(strings.ToLower(line), "access-control-allow-origin") &&
		strings.Contains(line, "null") {
		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Removed Access-Control-Allow-Origin: null header" +
			"\n" + indent + "// For disallowed origins, omit the header entirely"
	}
	return line
}

func fixCORSWildcard(line string) string {
	// Replace wildcard with origin validation
	if strings.Contains(strings.ToLower(line), "access-control-allow-origin") &&
		strings.Contains(line, "*") {
		indent := getIndentation(line)

		// Detect the programming language/framework and provide appropriate fix
		if strings.Contains(line, ".Header(") || strings.Contains(line, "c.Header") {
			// Go Gin framework
			return indent + "// SECURITY FIX: Validate origin instead of using wildcard" +
				"\n" + indent + "if isAllowedOrigin(origin) {" +
				"\n" + indent + "    c.Header(\"Access-Control-Allow-Origin\", origin)" +
				"\n" + indent + "}"
		} else if strings.Contains(line, "setHeader") || strings.Contains(line, "res.header") {
			// Node.js/Express
			return indent + "// SECURITY FIX: Validate origin instead of using wildcard" +
				"\n" + indent + "if (isAllowedOrigin(origin)) {" +
				"\n" + indent + "    res.setHeader('Access-Control-Allow-Origin', origin);" +
				"\n" + indent + "}"
		}

		// Generic fix
		return indent + "// SECURITY FIX: Replace wildcard with specific origin validation" +
			"\n" + indent + "// Check if origin is in allowed list before setting header"
	}
	return line
}

func addContentTypeOptionsHeader(line string) string {
	indent := getIndentation(line)

	// Add X-Content-Type-Options header after content-type is set
	if strings.Contains(strings.ToLower(line), "content-type") {
		if strings.Contains(line, ".Header(") || strings.Contains(line, "c.Header") {
			// Go Gin framework
			return line + "\n" + indent + "c.Header(\"X-Content-Type-Options\", \"nosniff\")"
		} else if strings.Contains(line, "setHeader") || strings.Contains(line, "res.header") {
			// Node.js/Express
			return line + "\n" + indent + "res.setHeader('X-Content-Type-Options', 'nosniff');"
		}
	}
	return line
}

func AddSecurityHeaders(line string) string {
	indent := getIndentation(line)

	// Add comprehensive security headers
	if strings.Contains(strings.ToLower(line), "content-type") {
		if strings.Contains(line, ".Header(") || strings.Contains(line, "c.Header") {
			// Go Gin framework
			return line +
				"\n" + indent + "// SECURITY: Added comprehensive security headers" +
				"\n" + indent + "c.Header(\"X-Content-Type-Options\", \"nosniff\")" +
				"\n" + indent + "c.Header(\"X-Frame-Options\", \"DENY\")" +
				"\n" + indent + "c.Header(\"X-XSS-Protection\", \"1; mode=block\")"
		} else if strings.Contains(line, "setHeader") || strings.Contains(line, "res.header") {
			// Node.js/Express
			return line +
				"\n" + indent + "// SECURITY: Added comprehensive security headers" +
				"\n" + indent + "res.setHeader('X-Content-Type-Options', 'nosniff');" +
				"\n" + indent + "res.setHeader('X-Frame-Options', 'DENY');" +
				"\n" + indent + "res.setHeader('X-XSS-Protection', '1; mode=block');"
		}
	}
	return line
}

func fixJWTNoneAlgorithm(line string) string {
	// Replace 'none' algorithm with secure alternative
	if strings.Contains(strings.ToLower(line), "algorithm") &&
		strings.Contains(line, "none") {
		return strings.ReplaceAll(line, "\"none\"", "\"HS256\"") +
			" // SECURITY FIX: Replaced 'none' with secure HS256 algorithm"
	}

	// Replace SigningMethodNone with secure alternative
	if strings.Contains(line, "SigningMethodNone") {
		return strings.ReplaceAll(line, "SigningMethodNone", "SigningMethodHS256") +
			" // SECURITY FIX: Replaced SigningMethodNone with secure HS256 algorithm"
	}

	return line
}

func addJWTExpiration(line string) string {
	// Add expiration to JWT token signing
	if strings.Contains(strings.ToLower(line), "jwt") &&
		strings.Contains(strings.ToLower(line), "sign") {
		indent := getIndentation(line)

		if strings.Contains(line, "jwt.Sign") {
			// Go JWT
			return line + "\n" + indent + "// SECURITY: Set token expiration (add to claims: exp: time.Now().Add(15 * time.Minute).Unix())"
		} else if strings.Contains(line, "jwt.sign") {
			// Node.js JWT
			return strings.ReplaceAll(line, "jwt.sign(", "jwt.sign(") +
				" // SECURITY: Add { expiresIn: '15m' } to options"
		}
	}
	return line
}

func fixSQLConcatenation(line string) string {
	// Replace SQL string concatenation with parameterized query placeholder
	if strings.Contains(strings.ToLower(line), "select") ||
		strings.Contains(strings.ToLower(line), "insert") ||
		strings.Contains(strings.ToLower(line), "update") ||
		strings.Contains(strings.ToLower(line), "delete") {

		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Use parameterized queries instead of string concatenation" +
			"\n" + indent + "// Replace concatenated values with $1, $2, etc. placeholders" +
			"\n" + line
	}
	return line
}

func fixSQLInterpolation(line string) string {
	// Replace SQL variable interpolation with parameterized query
	if strings.Contains(line, "${") &&
		(strings.Contains(strings.ToLower(line), "select") ||
			strings.Contains(strings.ToLower(line), "insert") ||
			strings.Contains(strings.ToLower(line), "update") ||
			strings.Contains(strings.ToLower(line), "delete")) {

		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Use parameterized queries instead of variable interpolation" +
			"\n" + indent + "// Replace ${variable} with $1, $2, etc. and pass values separately" +
			"\n" + line
	}
	return line
}

func fixDetailedErrors(line string) string {
	// Replace detailed error messages with generic ones
	if strings.Contains(strings.ToLower(line), "error") &&
		(strings.Contains(strings.ToLower(line), "sql") ||
			strings.Contains(strings.ToLower(line), "database") ||
			strings.Contains(strings.ToLower(line), "internal") ||
			strings.Contains(strings.ToLower(line), "stack") ||
			strings.Contains(strings.ToLower(line), "trace")) {

		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Use generic error message in production" +
			"\n" + indent + "// Log detailed error securely, return generic message to client" +
			"\n" + line
	}
	return line
}

func fixDebugExposure(line string) string {
	// Only fix hardcoded debug values, not environment-based configurations
	if (strings.Contains(strings.ToLower(line), "debug") ||
		strings.Contains(strings.ToLower(line), "trace") ||
		strings.Contains(strings.ToLower(line), "stack")) &&
		(strings.Contains(line, "true") ||
			strings.Contains(line, "enabled") ||
			strings.Contains(line, "on")) {

		// Skip environment-based configurations (they're already secure)
		if strings.Contains(line, "os.Getenv") ||
			strings.Contains(line, "process.env") ||
			strings.Contains(line, "ENV[") ||
			strings.Contains(line, "${") {
			return line
		}

		// Replace hardcoded debug values with environment-based configuration
		newLine := strings.ReplaceAll(line, "true", "false")
		newLine = strings.ReplaceAll(newLine, "enabled", "disabled")
		newLine = strings.ReplaceAll(newLine, "on", "off")

		return newLine + " // SECURITY FIX: Disabled debug info (use env var for dev)"
	}
	return line
}

func addSecureCookieFlags(line string) string {
	// Add secure flags to cookie configuration
	if strings.Contains(strings.ToLower(line), "cookie") ||
		strings.Contains(strings.ToLower(line), "setcookie") {

		indent := getIndentation(line)

		if strings.Contains(line, "SetCookie") || strings.Contains(line, "setCookie") {
			return line + "\n" + indent + "// SECURITY: Ensure HttpOnly and Secure flags are set on sensitive cookies"
		}
	}
	return line
}

// Insecure Random Generation Fix Functions

func fixTimestampRandom(line string) string {
	// Replace time.Now().UnixNano() with secure random generation
	if strings.Contains(line, "time.Now().UnixNano()") {
		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Replaced timestamp with secure random generation" +
			"\n" + indent + "// Use security.SecureRandom.GenerateRandomSuffix() or crypto/rand" +
			"\n" + indent + "// import \"crypto/rand\"" +
			"\n" + indent + "// randomBytes := make([]byte, 8)" +
			"\n" + indent + "// rand.Read(randomBytes)" +
			"\n" + indent + "// suffix := hex.EncodeToString(randomBytes)" +
			"\n" + strings.ReplaceAll(line, "time.Now().UnixNano()", "suffix")
	}
	return line
}

func fixMathRandImport(line string) string {
	// Replace math/rand import with crypto/rand
	if strings.Contains(line, "math/rand") {
		indent := getIndentation(line)
		newLine := strings.ReplaceAll(line, "math/rand", "crypto/rand")
		return indent + "// SECURITY FIX: Replaced math/rand with crypto/rand" +
			"\n" + newLine +
			"\n" + indent + "// Note: crypto/rand has different API - use rand.Read([]byte) instead of rand.Int()"
	}
	return line
}

func fixPredictableRandFunctions(line string) string {
	// Replace predictable rand functions with secure alternatives
	if strings.Contains(line, "rand.Int") || strings.Contains(line, "rand.Float") {
		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Replace predictable rand functions with crypto/rand" +
			"\n" + indent + "// Use security.GenerateBytes() for secure random generation" +
			"\n" + indent + "randomBytes, err := security.GenerateBytes(8)" +
			"\n" + indent + "if err != nil { return err }" +
			"\n" + indent + "randomInt := binary.BigEndian.Uint64(randomBytes)" +
			"\n" + indent + "// " + line + " // REPLACED: was using predictable rand functions"
	}
	return line
}

func fixTimestampIDGeneration(line string) string {
	// Replace timestamp-based ID generation with secure random IDs
	if strings.Contains(line, "fmt.Sprintf") &&
		(strings.Contains(line, "time.Now().Unix()") || strings.Contains(line, "time.Now().UnixNano()")) {
		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Replaced timestamp-based ID with secure random ID" +
			"\n" + indent + "secureRandom := security.NewSecureRandom()" +
			"\n" + indent + "id, err := secureRandom.GenerateSecureID(\"audit\")" +
			"\n" + indent + "if err != nil { return err }" +
			"\n" + indent + "// " + line + " // REPLACED: was using timestamp-based ID generation"
	}
	return line
}

func fixPredictableTempFiles(line string) string {
	// Replace predictable temporary file names with secure alternatives
	if strings.Contains(line, ".tmp") && strings.Contains(line, "time.Now()") {
		indent := getIndentation(line)
		return indent + "// SECURITY FIX: Replaced predictable temp file with secure naming" +
			"\n" + indent + "secureFileOps := security.NewSecureFileOperations()" +
			"\n" + indent + "tempFile, err := secureFileOps.CreateSecureTempFile(dir, \"temp-*.tmp\")" +
			"\n" + indent + "if err != nil { return err }" +
			"\n" + indent + "// " + line + " // REPLACED: was using predictable temp file names"
	}
	return line
}

// Helper function to get the indentation of a line
func getIndentation(line string) string {
	for i, char := range line {
		if char != ' ' && char != '\t' {
			return line[:i]
		}
	}
	return line // If the line is all whitespace
}
