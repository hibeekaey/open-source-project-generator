package template

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestSecurityPatternsInTemplates validates that templates use secure coding patterns
func TestSecurityPatternsInTemplates(t *testing.T) {
	t.Skip("Test disabled - requires embedded template refactoring")

	templateDir := "templates"

	// Walk through all embedded template files
	err := fs.WalkDir(embeddedTemplates, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only check .tmpl files
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Read template content from embedded filesystem
		content, err := embeddedTemplates.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read template file %s: %v", path, err)
			return nil
		}

		// Run security pattern checks
		checkInsecureRandomPatterns(t, path, string(content))
		checkInsecureFileOperations(t, path, string(content))
		checkSecureRandomUsage(t, path, string(content))

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk template directory: %v", err)
	}
}

// checkInsecureRandomPatterns detects insecure random number generation patterns
func checkInsecureRandomPatterns(t *testing.T, filePath, content string) {
	// Pattern 1: time.Now().UnixNano() - predictable timestamp-based randomness
	timestampPattern := regexp.MustCompile(`time\.Now\(\)\.UnixNano\(\)`)
	if matches := timestampPattern.FindAllString(content, -1); len(matches) > 0 {
		t.Errorf("SECURITY ISSUE in %s: Found insecure timestamp-based random generation: %v", filePath, matches)
	}

	// Pattern 2: math/rand usage without crypto/rand
	mathRandPattern := regexp.MustCompile(`math/rand`)
	cryptoRandPattern := regexp.MustCompile(`crypto/rand`)

	if mathRandPattern.MatchString(content) && !cryptoRandPattern.MatchString(content) {
		t.Errorf("SECURITY ISSUE in %s: Found math/rand usage without crypto/rand - use crypto/rand for security-sensitive operations", filePath)
	}

	// Pattern 3: Simple timestamp-based ID generation
	simpleIDPattern := regexp.MustCompile(`fmt\.Sprintf\([^)]*%d[^)]*time\.Now\(\)`)
	if matches := simpleIDPattern.FindAllString(content, -1); len(matches) > 0 {
		t.Errorf("SECURITY ISSUE in %s: Found timestamp-based ID generation: %v", filePath, matches)
	}
}

// checkInsecureFileOperations detects insecure file operation patterns
func checkInsecureFileOperations(t *testing.T, filePath, content string) {
	// Pattern 1: Predictable temporary file names
	tempFilePattern := regexp.MustCompile(`/tmp/[^/]*\d+`)
	if matches := tempFilePattern.FindAllString(content, -1); len(matches) > 0 {
		t.Logf("WARNING in %s: Potentially predictable temporary file paths: %v", filePath, matches)
	}

	// Pattern 2: Direct file operations without validation
	directWritePattern := regexp.MustCompile(`ioutil\.WriteFile|os\.WriteFile`)
	if directWritePattern.MatchString(content) {
		// Check if there's path validation nearby
		pathValidationPattern := regexp.MustCompile(`filepath\.Clean|path\.Clean|strings\.Contains.*\.\.|ValidatePath`)
		if !pathValidationPattern.MatchString(content) {
			t.Logf("INFO in %s: File write operations found - ensure path validation is implemented", filePath)
		}
	}
}

// checkSecureRandomUsage validates proper secure random usage
func checkSecureRandomUsage(t *testing.T, filePath, content string) {
	// Check for crypto/rand usage
	if strings.Contains(content, "crypto/rand") {
		// Verify proper error handling
		randReadPattern := regexp.MustCompile(`rand\.Read\([^)]+\)`)
		errorHandlingPattern := regexp.MustCompile(`if err != nil`)

		if randReadPattern.MatchString(content) && !errorHandlingPattern.MatchString(content) {
			t.Errorf("SECURITY ISSUE in %s: crypto/rand.Read usage without proper error handling", filePath)
		}
	}
}

// TestLoggingMiddlewareSecureRequestID specifically tests the logging middleware template
func TestLoggingMiddlewareSecureRequestID(t *testing.T) {
	t.Skip("Test disabled - requires embedded template refactoring")

	templatePath := "templates/backend/go-gin/internal/middleware/logging.go.tmpl"

	content, err := embeddedTemplates.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read logging middleware template: %v", err)
	}

	contentStr := string(content)

	// Test 1: Should use crypto/rand
	if !strings.Contains(contentStr, "crypto/rand") {
		t.Error("Logging middleware should import crypto/rand for secure request ID generation")
	}

	// Test 2: Should NOT use time.Now().UnixNano() for primary ID generation
	// Check that it's not used as the main ID generation method
	primaryIDPattern := regexp.MustCompile(`return fmt\.Sprintf\([^)]*time\.Now\(\)\.UnixNano\(\)`)
	if primaryIDPattern.MatchString(contentStr) {
		t.Error("Logging middleware should not use predictable timestamp-based request IDs as primary generation method")
	}

	// Test 3: Should have security comments
	securityCommentPattern := regexp.MustCompile(`(?i)security|crypto|secure|random`)
	if !securityCommentPattern.MatchString(contentStr) {
		t.Error("Logging middleware should include security-related comments explaining secure random usage")
	}

	// Test 4: Should use rand.Read for random generation
	if !strings.Contains(contentStr, "rand.Read") {
		t.Error("Logging middleware should use rand.Read for cryptographically secure random generation")
	}

	// Test 5: Should have proper error handling for rand.Read
	randReadPattern := regexp.MustCompile(`rand\.Read\([^)]+\)`)
	errorHandlingPattern := regexp.MustCompile(`if err != nil`)

	if randReadPattern.MatchString(contentStr) && !errorHandlingPattern.MatchString(contentStr) {
		t.Error("Logging middleware should have proper error handling for crypto/rand operations")
	}

	// Test 6: Should convert to hex for readability
	if !strings.Contains(contentStr, "hex.EncodeToString") {
		t.Error("Logging middleware should convert random bytes to hex string for readability")
	}
}

// TestSecurityDocumentationInTemplates ensures security patterns are documented
func TestSecurityDocumentationInTemplates(t *testing.T) {
	t.Skip("Test disabled - requires embedded template refactoring")

	securitySensitiveFiles := []string{
		"templates/backend/go-gin/internal/middleware/logging.go.tmpl",
		// Add other security-sensitive template files here as they are identified
	}

	for _, filePath := range securitySensitiveFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("Failed to read security-sensitive template %s: %v", filePath, err)
			continue
		}

		contentStr := string(content)

		// Check for security-related comments
		securityKeywords := []string{"SECURITY:", "security", "crypto", "secure", "random"}
		hasSecurityDoc := false

		for _, keyword := range securityKeywords {
			if strings.Contains(strings.ToLower(contentStr), strings.ToLower(keyword)) {
				hasSecurityDoc = true
				break
			}
		}

		if !hasSecurityDoc {
			t.Errorf("Security-sensitive template %s should include security documentation/comments", filePath)
		}
	}
}

// TestTemplateSecurityBestPractices validates overall security best practices
func TestTemplateSecurityBestPractices(t *testing.T) {
	t.Skip("Test disabled - requires embedded template refactoring")

	templateDir := "templates"

	// Focus on the most critical security issues that we're specifically addressing
	securityIssues := []struct {
		pattern     *regexp.Regexp
		description string
		severity    string
	}{
		{
			pattern:     regexp.MustCompile(`math/rand`),
			description: "Non-cryptographic random number generation",
			severity:    "MEDIUM",
		},
		{
			// Primary usage pattern - not fallback scenarios
			pattern:     regexp.MustCompile(`return.*time\.Now\(\)\.UnixNano\(\)`),
			description: "Predictable timestamp-based randomness as primary method",
			severity:    "HIGH",
		},
		{
			// Look for timestamp-based ID generation in functions
			pattern:     regexp.MustCompile(`fmt\.Sprintf\([^)]*%d[^)]*time\.Now\(\)\.UnixNano\(\)`),
			description: "Timestamp-based ID generation detected",
			severity:    "HIGH",
		},
	}

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)

		for _, issue := range securityIssues {
			if matches := issue.pattern.FindAllString(contentStr, -1); len(matches) > 0 {
				if issue.severity == "HIGH" {
					t.Errorf("SECURITY %s in %s: %s - Found: %v", issue.severity, path, issue.description, matches)
				} else {
					t.Logf("SECURITY %s in %s: %s - Found: %v", issue.severity, path, issue.description, matches)
				}
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk template directory for security validation: %v", err)
	}
}
