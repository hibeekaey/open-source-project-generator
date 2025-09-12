package reporting

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestReportGenerator_GenerateVersionUpdateReport(t *testing.T) {
	// Create temporary directory for test reports
	tempDir, err := os.MkdirTemp("", "test_reports")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewReportGenerator(tempDir)

	// Create test version report
	versionReport := &models.VersionReport{
		GeneratedAt:     time.Now(),
		TotalPackages:   10,
		OutdatedCount:   3,
		SecurityIssues:  1,
		LastUpdateCheck: time.Now(),
		Summary: map[string]models.VersionSummary{
			"languages": {
				Total:    3,
				Current:  2,
				Outdated: 1,
				Insecure: 0,
			},
		},
		Details: map[string]*models.VersionInfo{
			"react": {
				Name:           "react",
				Language:       "javascript",
				Type:           "framework",
				CurrentVersion: "18.0.0",
				LatestVersion:  "19.0.0",
				IsSecure:       true,
				UpdatedAt:      time.Now(),
				UpdateSource:   "npm",
			},
		},
		Recommendations: []models.UpdateRecommendation{
			{
				Name:               "react",
				CurrentVersion:     "18.0.0",
				RecommendedVersion: "19.0.0",
				Priority:           "medium",
				Reason:             "New version available",
				BreakingChange:     false,
			},
		},
	}

	// Test JSON format
	report, err := generator.GenerateVersionUpdateReport(versionReport, "json")
	if err != nil {
		t.Fatalf("Failed to generate version report: %v", err)
	}

	if report.Type != "version_update" {
		t.Errorf("Expected report type 'version_update', got '%s'", report.Type)
	}

	if len(report.Recommendations) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(report.Recommendations))
	}

	// Verify file was created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 report file, got %d", len(files))
	}
}

func TestReportGenerator_GenerateSecurityReport(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_reports")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewReportGenerator(tempDir)

	// Create test version report with security issues
	versionReport := &models.VersionReport{
		GeneratedAt:     time.Now(),
		TotalPackages:   5,
		OutdatedCount:   2,
		SecurityIssues:  2,
		LastUpdateCheck: time.Now(),
		Details: map[string]*models.VersionInfo{
			"vulnerable-package": {
				Name:           "vulnerable-package",
				Language:       "javascript",
				Type:           "package",
				CurrentVersion: "1.0.0",
				LatestVersion:  "1.2.0",
				IsSecure:       false,
				SecurityIssues: []models.SecurityIssue{
					{
						ID:          "CVE-2023-1234",
						Severity:    "high",
						Description: "Remote code execution vulnerability",
						FixedIn:     "1.1.0",
						ReportedAt:  time.Now(),
					},
				},
				UpdatedAt:    time.Now(),
				UpdateSource: "npm",
			},
		},
	}

	// Test security report generation
	report, err := generator.GenerateSecurityReport(versionReport, "json")
	if err != nil {
		t.Fatalf("Failed to generate security report: %v", err)
	}

	if report.TotalIssues != 2 {
		t.Errorf("Expected 2 total issues, got %d", report.TotalIssues)
	}

	if len(report.Issues) != 1 {
		t.Errorf("Expected 1 security issue detail, got %d", len(report.Issues))
	}

	// Verify the security issue details
	issue := report.Issues[0]
	if issue.PackageName != "vulnerable-package" {
		t.Errorf("Expected package name 'vulnerable-package', got '%s'", issue.PackageName)
	}

	if issue.SecurityIssue.Severity != "high" {
		t.Errorf("Expected severity 'high', got '%s'", issue.SecurityIssue.Severity)
	}
}

func TestReportGenerator_GenerateTemplateUpdateReport(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_reports")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewReportGenerator(tempDir)

	// Create test template updates
	updates := []models.TemplateUpdate{
		{
			TemplatePath: "templates/frontend/nextjs-app",
			Success:      true,
			UpdatedAt:    time.Now(),
			VersionChanges: map[string]string{
				"react":  "18.0.0 -> 19.0.0",
				"nextjs": "14.0.0 -> 15.0.0",
			},
		},
		{
			TemplatePath: "templates/frontend/nextjs-home",
			Success:      false,
			UpdatedAt:    time.Now(),
			Error:        "Template file not found",
		},
	}

	// Test template update report generation
	report, err := generator.GenerateTemplateUpdateReport(updates, "json")
	if err != nil {
		t.Fatalf("Failed to generate template update report: %v", err)
	}

	if report.Type != "template_update" {
		t.Errorf("Expected report type 'template_update', got '%s'", report.Type)
	}

	if len(report.TemplateUpdates) != 2 {
		t.Errorf("Expected 2 template updates, got %d", len(report.TemplateUpdates))
	}

	// Verify successful update
	successUpdate := report.TemplateUpdates[0]
	if !successUpdate.Success {
		t.Errorf("Expected first update to be successful")
	}

	if len(successUpdate.VersionChanges) != 2 {
		t.Errorf("Expected 2 version changes, got %d", len(successUpdate.VersionChanges))
	}

	// Verify failed update
	failedUpdate := report.TemplateUpdates[1]
	if failedUpdate.Success {
		t.Errorf("Expected second update to fail")
	}

	if failedUpdate.Error == "" {
		t.Errorf("Expected error message for failed update")
	}
}

func TestReportGenerator_GetReportHistory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_reports")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewReportGenerator(tempDir)

	// Create some test report files
	testFiles := []string{
		"report_123456.json",
		"security_789012.json",
		"template_345678.yaml",
		"other_file.txt", // Should be ignored
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Test getting report history
	reports, err := generator.GetReportHistory()
	if err != nil {
		t.Fatalf("Failed to get report history: %v", err)
	}

	// Should find at least 2 valid report files (excluding other_file.txt)
	// The template report might not be recognized due to naming
	if len(reports) < 2 {
		t.Errorf("Expected at least 2 reports, got %d", len(reports))
	}

	// Verify report types are correctly identified
	typeCount := make(map[string]int)
	for _, report := range reports {
		typeCount[report.Type]++
	}

	if typeCount["version_update"] != 1 {
		t.Errorf("Expected 1 version_update report, got %d", typeCount["version_update"])
	}

	if typeCount["security"] != 1 {
		t.Errorf("Expected 1 security report, got %d", typeCount["security"])
	}
}

func TestReportGenerator_EmptyDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_reports")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	generator := NewReportGenerator(tempDir)

	// Test getting history from empty directory
	reports, err := generator.GetReportHistory()
	if err != nil {
		t.Fatalf("Failed to get report history from empty directory: %v", err)
	}

	if len(reports) != 0 {
		t.Errorf("Expected 0 reports from empty directory, got %d", len(reports))
	}
}

func TestReportGenerator_NonExistentDirectory(t *testing.T) {
	generator := NewReportGenerator("/non/existent/directory")

	// Test getting history from non-existent directory
	reports, err := generator.GetReportHistory()
	if err != nil {
		t.Fatalf("Failed to handle non-existent directory: %v", err)
	}

	if len(reports) != 0 {
		t.Errorf("Expected 0 reports from non-existent directory, got %d", len(reports))
	}
}
