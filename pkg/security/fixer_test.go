package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixer_FixCORSNullOrigin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Go Gin CORS null origin",
			input:    `    c.Header("Access-Control-Allow-Origin", "null")`,
			expected: `    // SECURITY FIX: Removed Access-Control-Allow-Origin: null header` + "\n" + `    // For disallowed origins, omit the header entirely`,
		},
		{
			name:     "Node.js CORS null origin",
			input:    `  res.setHeader('Access-Control-Allow-Origin', 'null');`,
			expected: `  // SECURITY FIX: Removed Access-Control-Allow-Origin: null header` + "\n" + `  // For disallowed origins, omit the header entirely`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixCORSNullOrigin(tt.input)
			if result != tt.expected {
				t.Errorf("fixCORSNullOrigin() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFixer_FixCORSWildcard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "Go Gin CORS wildcard",
			input: `    c.Header("Access-Control-Allow-Origin", "*")`,
			contains: []string{
				"SECURITY FIX",
				"isAllowedOrigin",
				"c.Header",
			},
		},
		{
			name:  "Node.js CORS wildcard",
			input: `  res.setHeader('Access-Control-Allow-Origin', '*');`,
			contains: []string{
				"SECURITY FIX",
				"isAllowedOrigin",
				"res.setHeader",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixCORSWildcard(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("fixCORSWildcard() result should contain %q, got %q", expected, result)
				}
			}
		})
	}
}

func TestFixer_AddSecurityHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "Go Gin content-type header",
			input: `    c.Header("Content-Type", "application/json")`,
			contains: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
			},
		},
		{
			name:  "Node.js content-type header",
			input: `  res.setHeader('Content-Type', 'application/json');`,
			contains: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddSecurityHeaders(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("AddSecurityHeaders() result should contain %q, got %q", expected, result)
				}
			}
		})
	}
}

func TestFixer_FixJWTNoneAlgorithm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "JWT none algorithm",
			input:    `  algorithm: "none"`,
			expected: `  algorithm: "HS256" // SECURITY FIX: Replaced 'none' with secure HS256 algorithm`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixJWTNoneAlgorithm(tt.input)
			if result != tt.expected {
				t.Errorf("fixJWTNoneAlgorithm() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFixer_FixFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "security-fixer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with security issues
	testFile := filepath.Join(tempDir, "test.go.tmpl")
	testContent := `package main

import "github.com/gin-gonic/gin"

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "null")
        c.Header("Content-Type", "application/json")
        c.Next()
    }
}

func jwtConfig() {
    algorithm: "none"
}`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create fixer and apply fixes
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      false,
		FixType:      "all",
		CreateBackup: true,
	}

	result, err := fixer.FixFile(testFile, options)
	if err != nil {
		t.Fatalf("FixFile() error = %v", err)
	}

	// Verify fixes were applied
	if len(result.FixedIssues) == 0 {
		t.Error("Expected security issues to be fixed, but none were found")
	}

	// Verify backup was created
	if result.BackupsCreated != 1 {
		t.Errorf("Expected 1 backup to be created, got %d", result.BackupsCreated)
	}

	// Read the fixed file and verify changes
	fixedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read fixed file: %v", err)
	}

	fixedStr := string(fixedContent)

	// Verify CORS null origin was fixed
	if strings.Contains(fixedStr, `"null"`) {
		t.Error("CORS null origin should have been fixed")
	}

	// Verify security headers were added
	if !strings.Contains(fixedStr, "X-Content-Type-Options") {
		t.Error("Security headers should have been added")
	}

	// Verify JWT algorithm was fixed
	if strings.Contains(fixedStr, `"none"`) {
		t.Error("JWT none algorithm should have been fixed")
	}
}

func TestFixer_DryRun(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "security-fixer-dryrun-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with security issues
	testFile := filepath.Join(tempDir, "test.go.tmpl")
	testContent := `c.Header("Access-Control-Allow-Origin", "null")`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create fixer and run in dry-run mode
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       true,
		Verbose:      false,
		FixType:      "all",
		CreateBackup: false,
	}

	result, err := fixer.FixFile(testFile, options)
	if err != nil {
		t.Fatalf("FixFile() error = %v", err)
	}

	// Verify issues were detected but not fixed
	if len(result.FixedIssues) == 0 {
		t.Error("Expected security issues to be detected in dry-run mode")
	}

	// Verify file was not modified
	currentContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(currentContent) != testContent {
		t.Error("File should not be modified in dry-run mode")
	}

	// Verify no backup was created
	if result.BackupsCreated != 0 {
		t.Errorf("Expected 0 backups in dry-run mode, got %d", result.BackupsCreated)
	}
}

func TestFixer_FixTypeFilter(t *testing.T) {
	fixer := NewFixer()

	tests := []struct {
		name      string
		issueType SecurityIssueType
		fixType   string
		expected  bool
	}{
		{"All fixes", CORSVulnerability, "all", true},
		{"CORS filter matches", CORSVulnerability, "cors", true},
		{"CORS filter doesn't match headers", MissingSecurityHeader, "cors", false},
		{"Headers filter matches", MissingSecurityHeader, "headers", true},
		{"Auth filter matches", WeakAuthentication, "auth", true},
		{"SQL filter matches", SQLInjectionRisk, "sql", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixer.shouldApplyFix(tt.issueType, tt.fixType)
			if result != tt.expected {
				t.Errorf("shouldApplyFix(%v, %q) = %v, want %v", tt.issueType, tt.fixType, result, tt.expected)
			}
		})
	}
}

func TestGetIndentation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"No indentation", "code", ""},
		{"Space indentation", "    code", "    "},
		{"Tab indentation", "\t\tcode", "\t\t"},
		{"Mixed indentation", " \t code", " \t "},
		{"All whitespace", "    ", "    "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIndentation(tt.input)
			if result != tt.expected {
				t.Errorf("getIndentation(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
