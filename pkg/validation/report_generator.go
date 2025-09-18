package validation

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ReportGenerator generates validation reports in various formats
type ReportGenerator struct {
	templates map[string]*template.Template
}

// NewReportGenerator creates a new report generator
func NewReportGenerator() *ReportGenerator {
	generator := &ReportGenerator{
		templates: make(map[string]*template.Template),
	}
	generator.initializeTemplates()
	return generator
}

// ReportOptions contains options for report generation
type ReportOptions struct {
	Title           string
	IncludeMetadata bool
	IncludeFixes    bool
	GroupBySeverity bool
	GroupByCategory bool
	ShowOnlyErrors  bool
	ShowStatistics  bool
	CustomCSS       string
	OutputPath      string
}

// ValidationSummary contains aggregated validation statistics
type ValidationSummary struct {
	TotalResults     int                            `json:"total_results"`
	TotalFiles       int                            `json:"total_files"`
	ValidFiles       int                            `json:"valid_files"`
	ErrorCount       int                            `json:"error_count"`
	WarningCount     int                            `json:"warning_count"`
	FixableCount     int                            `json:"fixable_count"`
	IssuesByCategory map[string]int                 `json:"issues_by_category"`
	IssuesBySeverity map[string]int                 `json:"issues_by_severity"`
	IssuesByRule     map[string]int                 `json:"issues_by_rule"`
	Results          []*interfaces.ValidationResult `json:"results"`
	GeneratedAt      time.Time                      `json:"generated_at"`
}

// GenerateReport generates a validation report in the specified format
func (rg *ReportGenerator) GenerateReport(result *interfaces.ValidationResult, format string, options ReportOptions) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return rg.generateJSONReport(result, options)
	case "html":
		return rg.generateHTMLReport(result, options)
	case "markdown", "md":
		return rg.generateMarkdownReport(result, options)
	case "xml":
		return rg.generateXMLReport(result, options)
	case "csv":
		return rg.generateCSVReport(result, options)
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}
}

// generateJSONReport generates a JSON validation report
func (rg *ReportGenerator) generateJSONReport(result *interfaces.ValidationResult, options ReportOptions) ([]byte, error) {
	report := map[string]interface{}{
		"validation_result": result,
		"metadata": map[string]interface{}{
			"generated_at": time.Now().UTC(),
			"format":       "json",
			"version":      "1.0",
		},
	}

	if options.IncludeMetadata {
		report["options"] = options
	}

	if options.ShowStatistics {
		report["statistics"] = rg.calculateStatistics(result)
	}

	return json.MarshalIndent(report, "", "  ")
}

// generateMarkdownReport generates a Markdown validation report
func (rg *ReportGenerator) generateMarkdownReport(result *interfaces.ValidationResult, options ReportOptions) ([]byte, error) {
	var buf strings.Builder

	// Header
	title := options.Title
	if title == "" {
		title = "Validation Report"
	}
	buf.WriteString(fmt.Sprintf("# %s\n\n", title))

	if options.IncludeMetadata {
		buf.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05 UTC")))
	}

	// Status
	if result.Valid {
		buf.WriteString("✅ **Status: VALID**\n\n")
	} else {
		buf.WriteString("❌ **Status: INVALID**\n\n")
	}

	// Statistics
	if options.ShowStatistics {
		stats := rg.calculateStatistics(result)
		buf.WriteString("## Summary\n\n")
		buf.WriteString(fmt.Sprintf("- **Total Files:** %d\n", stats["total_files"]))
		buf.WriteString(fmt.Sprintf("- **Valid Files:** %d\n", stats["valid_files"]))
		buf.WriteString(fmt.Sprintf("- **Total Issues:** %d\n", stats["total_issues"]))
		buf.WriteString(fmt.Sprintf("- **Errors:** %d\n", stats["error_count"]))
		buf.WriteString(fmt.Sprintf("- **Warnings:** %d\n", stats["warning_count"]))
		buf.WriteString(fmt.Sprintf("- **Info:** %d\n", stats["info_count"]))
		buf.WriteString(fmt.Sprintf("- **Fixable Issues:** %d\n\n", stats["fixable_issues"]))
	}

	// Issues
	if len(result.Issues) > 0 {
		buf.WriteString("## Issues\n\n")
		for _, issue := range result.Issues {
			icon := rg.getSeverityIcon(issue.Severity)
			buf.WriteString(fmt.Sprintf("%s **%s** (%s): %s\n", icon, strings.ToUpper(issue.Type[:1])+issue.Type[1:], strings.ToUpper(issue.Severity[:1])+issue.Severity[1:], issue.Message))
			if issue.File != "" {
				buf.WriteString(fmt.Sprintf("   - File: `%s`\n", issue.File))
			}
			if issue.Line > 0 {
				buf.WriteString(fmt.Sprintf("   - Line: %d\n", issue.Line))
			}
			if issue.Rule != "" {
				buf.WriteString(fmt.Sprintf("   - Rule: `%s`\n", issue.Rule))
			}
			if issue.Fixable {
				buf.WriteString("   - ✅ Fixable\n")
			}
			buf.WriteString("\n")
		}
	}

	return []byte(buf.String()), nil
}

// generateHTMLReport generates an HTML validation report
func (rg *ReportGenerator) generateHTMLReport(result *interfaces.ValidationResult, options ReportOptions) ([]byte, error) {
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f5f5f5; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .issues { margin: 20px 0; }
        .issue { margin: 10px 0; padding: 10px; border-left: 4px solid #ccc; }
        .error { border-left-color: #d32f2f; background-color: #ffebee; }
        .warning { border-left-color: #f57c00; background-color: #fff3e0; }
        .info { border-left-color: #1976d2; background-color: #e3f2fd; }
        .valid { color: #388e3c; }
        .invalid { color: #d32f2f; }
    </style>
</head>
<body>
    <div class="header">
        <h1>%s</h1>
        <p class="%s">Status: %s</p>
    </div>
    <div class="summary">
        <h2>Summary</h2>
        <ul>
            <li>Total Files: %d</li>
            <li>Valid Files: %d</li>
            <li>Total Issues: %d</li>
            <li>Fixable Issues: %d</li>
        </ul>
    </div>
</body>
</html>
`,
		options.Title,
		options.Title,
		map[bool]string{true: "valid", false: "invalid"}[result.Valid],
		map[bool]string{true: "VALID", false: "INVALID"}[result.Valid],
		result.Summary.TotalFiles,
		result.Summary.ValidFiles,
		result.Summary.ErrorCount+result.Summary.WarningCount,
		result.Summary.FixableCount,
	)

	return []byte(htmlContent), nil
}

// generateXMLReport generates an XML validation report
func (rg *ReportGenerator) generateXMLReport(result *interfaces.ValidationResult, options ReportOptions) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	buf.WriteString(`<validation-report>` + "\n")
	buf.WriteString(fmt.Sprintf(`  <summary valid="%t" total-files="%d" valid-files="%d" total-issues="%d"/>`+"\n",
		result.Valid, result.Summary.TotalFiles, result.Summary.ValidFiles, result.Summary.ErrorCount+result.Summary.WarningCount))
	buf.WriteString(`</validation-report>` + "\n")

	return []byte(buf.String()), nil
}

// generateCSVReport generates a CSV validation report
func (rg *ReportGenerator) generateCSVReport(result *interfaces.ValidationResult, options ReportOptions) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString("Type,Severity,Message,File,Line,Rule,Fixable\n")

	for _, issue := range result.Issues {
		buf.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%s,%t\n",
			issue.Type, issue.Severity, issue.Message, issue.File, issue.Line, issue.Rule, issue.Fixable))
	}

	return []byte(buf.String()), nil
}

// calculateStatistics calculates statistics for a validation result
func (rg *ReportGenerator) calculateStatistics(result *interfaces.ValidationResult) map[string]interface{} {
	return map[string]interface{}{
		"total_files":   result.Summary.TotalFiles,
		"valid_files":   result.Summary.ValidFiles,
		"total_issues":  result.Summary.ErrorCount + result.Summary.WarningCount,
		"error_count":   result.Summary.ErrorCount,
		"warning_count": result.Summary.WarningCount,
		"fixable_count": result.Summary.FixableCount,
	}
}

// getSeverityIcon returns an icon for the severity level
func (rg *ReportGenerator) getSeverityIcon(severity string) string {
	switch strings.ToLower(severity) {
	case "error":
		return "🔴"
	case "warning":
		return "🟡"
	case "info":
		return "🔵"
	default:
		return "⚪"
	}
}

// initializeTemplates initializes HTML templates
func (rg *ReportGenerator) initializeTemplates() {
	// Templates would be initialized here
}
