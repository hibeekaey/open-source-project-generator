package validation

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportGenerator(t *testing.T) {
	generator := NewReportGenerator()

	assert.NotNil(t, generator)
	assert.NotNil(t, generator.templates)
}

func TestReportGenerator_GenerateReport(t *testing.T) {
	generator := NewReportGenerator()

	// Create test validation result
	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: interfaces.ValidationSeverityError,
				Message:  "Test error message",
				File:     "/test/file.go",
				Line:     10,
				Column:   5,
				Rule:     "test.rule",
				Fixable:  true,
			},
			{
				Type:     "warning",
				Severity: interfaces.ValidationSeverityWarning,
				Message:  "Test warning message",
				File:     "/test/another.go",
				Line:     20,
				Column:   15,
				Rule:     "test.warning",
				Fixable:  false,
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   10,
			ValidFiles:   8,
			ErrorCount:   1,
			WarningCount: 1,
		},
	}

	options := ReportOptions{
		Title:           "Test Validation Report",
		IncludeMetadata: true,
		ShowStatistics:  true,
	}

	tests := []struct {
		name           string
		format         string
		expectError    bool
		validateOutput func([]byte) error
	}{
		{
			name:        "JSON format",
			format:      "json",
			expectError: false,
			validateOutput: func(output []byte) error {
				var report map[string]interface{}
				return json.Unmarshal(output, &report)
			},
		},
		{
			name:        "Markdown format",
			format:      "markdown",
			expectError: false,
			validateOutput: func(output []byte) error {
				content := string(output)
				if !strings.Contains(content, "# Test Validation Report") {
					return assert.AnError
				}
				if !strings.Contains(content, "Test error message") {
					return assert.AnError
				}
				return nil
			},
		},
		{
			name:        "HTML format",
			format:      "html",
			expectError: false,
			validateOutput: func(output []byte) error {
				content := string(output)
				if !strings.Contains(content, "<html>") {
					return assert.AnError
				}
				if !strings.Contains(content, "Test Validation Report") {
					return assert.AnError
				}
				return nil
			},
		},
		{
			name:        "XML format",
			format:      "xml",
			expectError: false,
			validateOutput: func(output []byte) error {
				content := string(output)
				if !strings.Contains(content, "<?xml") {
					return assert.AnError
				}
				if !strings.Contains(content, "<validation-report>") {
					return assert.AnError
				}
				return nil
			},
		},
		{
			name:        "CSV format",
			format:      "csv",
			expectError: false,
			validateOutput: func(output []byte) error {
				content := string(output)
				if !strings.Contains(content, "Type,Severity,Message") {
					return assert.AnError
				}
				if !strings.Contains(content, "error,error,Test error message") {
					return assert.AnError
				}
				return nil
			},
		},
		{
			name:        "unsupported format",
			format:      "unsupported",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := generator.GenerateReport(result, tt.format, options)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				if tt.validateOutput != nil {
					assert.NoError(t, tt.validateOutput(output))
				}
			}
		})
	}
}

func TestReportGenerator_generateJSONReport(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: true,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "info",
				Severity: interfaces.ValidationSeverityInfo,
				Message:  "Test info message",
				File:     "/test/file.go",
				Rule:     "test.info",
				Fixable:  false,
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   5,
			ValidFiles:   5,
			ErrorCount:   0,
			WarningCount: 0,
		},
	}

	tests := []struct {
		name    string
		options ReportOptions
	}{
		{
			name: "basic JSON report",
			options: ReportOptions{
				Title: "Basic Report",
			},
		},
		{
			name: "JSON report with metadata",
			options: ReportOptions{
				Title:           "Report with Metadata",
				IncludeMetadata: true,
			},
		},
		{
			name: "JSON report with statistics",
			options: ReportOptions{
				Title:          "Report with Statistics",
				ShowStatistics: true,
			},
		},
		{
			name: "JSON report with all options",
			options: ReportOptions{
				Title:           "Complete Report",
				IncludeMetadata: true,
				ShowStatistics:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := generator.generateJSONReport(result, tt.options)
			require.NoError(t, err)
			require.NotNil(t, output)

			// Validate JSON structure
			var report map[string]interface{}
			err = json.Unmarshal(output, &report)
			require.NoError(t, err)

			// Check required fields
			assert.Contains(t, report, "validation_result")
			assert.Contains(t, report, "metadata")

			if tt.options.IncludeMetadata {
				assert.Contains(t, report, "options")
			}

			if tt.options.ShowStatistics {
				assert.Contains(t, report, "statistics")
			}
		})
	}
}

func TestReportGenerator_generateMarkdownReport(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: interfaces.ValidationSeverityError,
				Message:  "Critical error found",
				File:     "/src/main.go",
				Line:     42,
				Rule:     "critical.error",
				Fixable:  true,
			},
			{
				Type:     "warning",
				Severity: interfaces.ValidationSeverityWarning,
				Message:  "Potential issue detected",
				File:     "/src/utils.go",
				Line:     15,
				Rule:     "potential.issue",
				Fixable:  false,
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   10,
			ValidFiles:   8,
			ErrorCount:   1,
			WarningCount: 1,
		},
	}

	tests := []struct {
		name           string
		options        ReportOptions
		expectedTitle  string
		expectMetadata bool
		expectStats    bool
	}{
		{
			name: "basic markdown report",
			options: ReportOptions{
				Title: "Test Report",
			},
			expectedTitle:  "Test Report",
			expectMetadata: false,
			expectStats:    false,
		},
		{
			name: "markdown report with metadata",
			options: ReportOptions{
				Title:           "Report with Metadata",
				IncludeMetadata: true,
			},
			expectedTitle:  "Report with Metadata",
			expectMetadata: true,
			expectStats:    false,
		},
		{
			name: "markdown report with statistics",
			options: ReportOptions{
				Title:          "Report with Stats",
				ShowStatistics: true,
			},
			expectedTitle:  "Report with Stats",
			expectMetadata: false,
			expectStats:    true,
		},
		{
			name: "default title",
			options: ReportOptions{
				ShowStatistics: true,
			},
			expectedTitle:  "Validation Report",
			expectMetadata: false,
			expectStats:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := generator.generateMarkdownReport(result, tt.options)
			require.NoError(t, err)
			require.NotNil(t, output)

			content := string(output)

			// Check title
			assert.Contains(t, content, "# "+tt.expectedTitle)

			// Check status
			assert.Contains(t, content, "❌ **Status: INVALID**")

			// Check metadata
			if tt.expectMetadata {
				assert.Contains(t, content, "**Generated:**")
			}

			// Check statistics
			if tt.expectStats {
				assert.Contains(t, content, "## Summary")
				assert.Contains(t, content, "- **Total Files:** 10")
				assert.Contains(t, content, "- **Valid Files:** 8")
				assert.Contains(t, content, "- **Errors:** 1")
				assert.Contains(t, content, "- **Warnings:** 1")
			}

			// Check issues
			assert.Contains(t, content, "## Issues")
			assert.Contains(t, content, "Critical error found")
			assert.Contains(t, content, "Potential issue detected")
			assert.Contains(t, content, "✅ Fixable")
		})
	}
}

func TestReportGenerator_generateHTMLReport(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid:  true,
		Issues: []interfaces.ValidationIssue{},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   5,
			ValidFiles:   5,
			ErrorCount:   0,
			WarningCount: 0,
		},
	}

	options := ReportOptions{
		Title: "HTML Test Report",
	}

	output, err := generator.generateHTMLReport(result, options)
	require.NoError(t, err)
	require.NotNil(t, output)

	content := string(output)

	// Check HTML structure
	assert.Contains(t, content, "<!DOCTYPE html>")
	assert.Contains(t, content, "<html>")
	assert.Contains(t, content, "<head>")
	assert.Contains(t, content, "<body>")
	assert.Contains(t, content, "</html>")

	// Check title
	assert.Contains(t, content, "<title>HTML Test Report</title>")
	assert.Contains(t, content, "<h1>HTML Test Report</h1>")

	// Check status
	assert.Contains(t, content, "Status: VALID")
	assert.Contains(t, content, `class="valid"`)

	// Check summary
	assert.Contains(t, content, "Total Files: 5")
	assert.Contains(t, content, "Valid Files: 5")
}

func TestReportGenerator_generateXMLReport(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:    "error",
				Message: "Test error",
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   3,
			ValidFiles:   2,
			ErrorCount:   1,
			WarningCount: 0,
		},
	}

	options := ReportOptions{
		Title: "XML Test Report",
	}

	output, err := generator.generateXMLReport(result, options)
	require.NoError(t, err)
	require.NotNil(t, output)

	content := string(output)

	// Check XML structure
	assert.Contains(t, content, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, content, "<validation-report>")
	assert.Contains(t, content, "</validation-report>")

	// Check summary attributes
	assert.Contains(t, content, `valid="false"`)
	assert.Contains(t, content, `total-files="3"`)
	assert.Contains(t, content, `valid-files="2"`)
	assert.Contains(t, content, `total-issues="1"`)
}

func TestReportGenerator_generateCSVReport(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "error",
				Severity: interfaces.ValidationSeverityError,
				Message:  "Test error message",
				File:     "/test/file.go",
				Line:     10,
				Rule:     "test.rule",
				Fixable:  true,
			},
			{
				Type:     "warning",
				Severity: interfaces.ValidationSeverityWarning,
				Message:  "Test warning message",
				File:     "/test/another.go",
				Line:     20,
				Rule:     "test.warning",
				Fixable:  false,
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   2,
			ValidFiles:   0,
			ErrorCount:   1,
			WarningCount: 1,
		},
	}

	options := ReportOptions{
		Title: "CSV Test Report",
	}

	output, err := generator.generateCSVReport(result, options)
	require.NoError(t, err)
	require.NotNil(t, output)

	content := string(output)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	// Check header
	assert.Equal(t, "Type,Severity,Message,File,Line,Rule,Fixable", lines[0])

	// Check data rows
	assert.Len(t, lines, 3) // Header + 2 issues
	assert.Contains(t, lines[1], "error,error,Test error message,/test/file.go,10,test.rule,true")
	assert.Contains(t, lines[2], "warning,warning,Test warning message,/test/another.go,20,test.warning,false")
}

func TestReportGenerator_calculateStatistics(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: false,
		Issues: []interfaces.ValidationIssue{
			{Type: "error"},
			{Type: "warning"},
			{Type: "warning"},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   10,
			ValidFiles:   7,
			ErrorCount:   1,
			WarningCount: 2,
			FixableCount: 2,
		},
	}

	stats := generator.calculateStatistics(result)

	assert.Equal(t, 10, stats["total_files"])
	assert.Equal(t, 7, stats["valid_files"])
	assert.Equal(t, 3, stats["total_issues"])
	assert.Equal(t, 1, stats["error_count"])
	assert.Equal(t, 2, stats["warning_count"])
	assert.Equal(t, 2, stats["fixable_count"])
}

func TestReportGenerator_getSeverityIcon(t *testing.T) {
	generator := NewReportGenerator()

	tests := []struct {
		severity string
		expected string
	}{
		{"error", "🔴"},
		{"warning", "🟡"},
		{"info", "🔵"},
		{"unknown", "⚪"},
		{"ERROR", "🔴"},
		{"WARNING", "🟡"},
		{"INFO", "🔵"},
	}

	for _, tt := range tests {
		t.Run(tt.severity, func(t *testing.T) {
			icon := generator.getSeverityIcon(tt.severity)
			assert.Equal(t, tt.expected, icon)
		})
	}
}

func TestReportOptions_Variations(t *testing.T) {
	generator := NewReportGenerator()

	result := &interfaces.ValidationResult{
		Valid: true,
		Issues: []interfaces.ValidationIssue{
			{
				Type:     "info",
				Severity: interfaces.ValidationSeverityInfo,
				Message:  "Info message",
				File:     "/test/file.go",
				Rule:     "test.info",
				Fixable:  false,
			},
		},
		Summary: interfaces.ValidationSummary{
			TotalFiles:   1,
			ValidFiles:   1,
			ErrorCount:   0,
			WarningCount: 0,
		},
	}

	tests := []struct {
		name    string
		options ReportOptions
	}{
		{
			name: "all options enabled",
			options: ReportOptions{
				Title:           "Complete Report",
				IncludeMetadata: true,
				IncludeFixes:    true,
				GroupBySeverity: true,
				GroupByCategory: true,
				ShowOnlyErrors:  false,
				ShowStatistics:  true,
				CustomCSS:       "body { color: red; }",
				OutputPath:      "/tmp/report.html",
			},
		},
		{
			name: "minimal options",
			options: ReportOptions{
				Title: "Minimal Report",
			},
		},
		{
			name: "errors only",
			options: ReportOptions{
				Title:          "Errors Only",
				ShowOnlyErrors: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with different formats
			formats := []string{"json", "markdown", "html", "xml", "csv"}

			for _, format := range formats {
				output, err := generator.GenerateReport(result, format, tt.options)
				assert.NoError(t, err, "Format: %s", format)
				assert.NotNil(t, output, "Format: %s", format)
				assert.Greater(t, len(output), 0, "Format: %s", format)
			}
		})
	}
}

func TestValidationSummary_Structure(t *testing.T) {
	// Test the ValidationSummary struct
	now := time.Now()
	summary := ValidationSummary{
		TotalResults:     5,
		TotalFiles:       10,
		ValidFiles:       8,
		ErrorCount:       1,
		WarningCount:     1,
		FixableCount:     1,
		IssuesByCategory: map[string]int{"security": 1, "style": 1},
		IssuesBySeverity: map[string]int{"error": 1, "warning": 1},
		IssuesByRule:     map[string]int{"rule1": 1, "rule2": 1},
		Results:          []*interfaces.ValidationResult{},
		GeneratedAt:      now,
	}

	assert.Equal(t, 5, summary.TotalResults)
	assert.Equal(t, 10, summary.TotalFiles)
	assert.Equal(t, 8, summary.ValidFiles)
	assert.Equal(t, 1, summary.ErrorCount)
	assert.Equal(t, 1, summary.WarningCount)
	assert.Equal(t, 1, summary.FixableCount)
	assert.Equal(t, map[string]int{"security": 1, "style": 1}, summary.IssuesByCategory)
	assert.Equal(t, map[string]int{"error": 1, "warning": 1}, summary.IssuesBySeverity)
	assert.Equal(t, map[string]int{"rule1": 1, "rule2": 1}, summary.IssuesByRule)
	assert.Equal(t, now, summary.GeneratedAt)
}
