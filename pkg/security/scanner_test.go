package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_ScanFile(t *testing.T) {
	// Create a temporary test file with security issues
	testContent := `
// CORS vulnerability
c.Header("Access-Control-Allow-Origin", "null")

// Missing security headers
c.Header("Content-Type", "application/json")

// JWT with none algorithm
token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)

// SQL injection risk
query := "SELECT * FROM users WHERE id = " + userID

// Weak JWT secret
jwt_secret := "secret"
`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	scanner := NewScanner()
	issues, err := scanner.ScanFile(testFile)
	if err != nil {
		t.Fatalf("ScanFile failed: %v", err)
	}

	// Verify that security issues were detected
	if len(issues) == 0 {
		t.Error("Expected security issues to be found, but none were detected")
	}

	// Check for specific issue types
	foundCORS := false
	foundSQL := false
	foundJWT := false

	for _, issue := range issues {
		switch issue.IssueType {
		case CORSVulnerability:
			foundCORS = true
		case SQLInjectionRisk:
			foundSQL = true
		case WeakAuthentication:
			foundJWT = true
		}
	}

	if !foundCORS {
		t.Error("Expected CORS vulnerability to be detected")
	}
	if !foundSQL {
		t.Error("Expected SQL injection risk to be detected")
	}
	if !foundJWT {
		t.Error("Expected JWT authentication issue to be detected")
	}
}

func TestScanner_ScanDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with different extensions
	testFiles := map[string]string{
		"cors.go.tmpl": `c.Header("Access-Control-Allow-Origin", "null")`,
		"auth.js":      `const token = jwt.sign(payload, "secret", { algorithm: "none" });`,
		"db.go":        `query := "SELECT * FROM users WHERE id = " + id`,
		"safe.go":      `// This file has no security issues`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	scanner := NewScanner()
	report, err := scanner.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if report.ScannedFiles != 4 {
		t.Errorf("Expected 4 files to be scanned, got %d", report.ScannedFiles)
	}

	if len(report.Issues) == 0 {
		t.Error("Expected security issues to be found")
	}
}

func TestSecurityReport_HasCriticalIssues(t *testing.T) {
	report := &SecurityReport{
		Issues: []SecurityIssue{
			{Severity: SeverityMedium},
			{Severity: SeverityHigh},
		},
	}

	if report.HasCriticalIssues() {
		t.Error("Expected no critical issues")
	}

	report.Issues = append(report.Issues, SecurityIssue{Severity: SeverityCritical})

	if !report.HasCriticalIssues() {
		t.Error("Expected critical issues to be detected")
	}
}

func TestSecurityReport_CountBySeverity(t *testing.T) {
	report := &SecurityReport{
		Issues: []SecurityIssue{
			{Severity: SeverityCritical},
			{Severity: SeverityHigh},
			{Severity: SeverityHigh},
			{Severity: SeverityMedium},
		},
	}

	if count := report.CountBySeverity(SeverityCritical); count != 1 {
		t.Errorf("Expected 1 critical issue, got %d", count)
	}

	if count := report.CountBySeverity(SeverityHigh); count != 2 {
		t.Errorf("Expected 2 high severity issues, got %d", count)
	}

	if count := report.CountBySeverity(SeverityLow); count != 0 {
		t.Errorf("Expected 0 low severity issues, got %d", count)
	}
}

func TestIsTemplateFile(t *testing.T) {
	testCases := []struct {
		filename string
		expected bool
	}{
		{"test.go.tmpl", true},
		{"config.yaml", true},
		{"script.js", true},
		{"style.css", false},
		{"README.md", false},
		{"main.go", true},
	}

	for _, tc := range testCases {
		result := isTemplateFile(tc.filename)
		if result != tc.expected {
			t.Errorf("isTemplateFile(%s) = %v, expected %v", tc.filename, result, tc.expected)
		}
	}
}

func TestSecurityReport_CountByType(t *testing.T) {
	report := &SecurityReport{
		Issues: []SecurityIssue{
			{IssueType: CORSVulnerability},
			{IssueType: WeakAuthentication},
			{IssueType: WeakAuthentication},
			{IssueType: SQLInjectionRisk},
		},
	}

	if count := report.CountByType(CORSVulnerability); count != 1 {
		t.Errorf("Expected 1 CORS vulnerability, got %d", count)
	}

	if count := report.CountByType(WeakAuthentication); count != 2 {
		t.Errorf("Expected 2 weak authentication issues, got %d", count)
	}

	if count := report.CountByType(MissingSecurityHeader); count != 0 {
		t.Errorf("Expected 0 missing security header issues, got %d", count)
	}
}

func TestSecurityReport_GetFixableIssues(t *testing.T) {
	report := &SecurityReport{
		Issues: []SecurityIssue{
			{FixAvailable: true},
			{FixAvailable: false},
			{FixAvailable: true},
			{FixAvailable: false},
		},
	}

	fixable := report.GetFixableIssues()
	if len(fixable) != 2 {
		t.Errorf("Expected 2 fixable issues, got %d", len(fixable))
	}

	for _, issue := range fixable {
		if !issue.FixAvailable {
			t.Error("GetFixableIssues returned non-fixable issue")
		}
	}
}

func TestSecurityReport_FilterBySeverity(t *testing.T) {
	report := &SecurityReport{
		Issues: []SecurityIssue{
			{Severity: SeverityLow},
			{Severity: SeverityMedium},
			{Severity: SeverityHigh},
			{Severity: SeverityCritical},
		},
	}

	// Test filtering by medium severity (should include medium, high, critical)
	filtered := report.FilterBySeverity(SeverityMedium)
	if len(filtered) != 3 {
		t.Errorf("Expected 3 issues with medium+ severity, got %d", len(filtered))
	}

	// Test filtering by critical severity (should include only critical)
	filtered = report.FilterBySeverity(SeverityCritical)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 critical issue, got %d", len(filtered))
	}

	// Test filtering by low severity (should include all)
	filtered = report.FilterBySeverity(SeverityLow)
	if len(filtered) != 4 {
		t.Errorf("Expected 4 issues with low+ severity, got %d", len(filtered))
	}
}
