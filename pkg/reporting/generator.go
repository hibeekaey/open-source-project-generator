package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

// ReportGenerator handles generation of version update reports
type ReportGenerator struct {
	outputDir string
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(outputDir string) *ReportGenerator {
	return &ReportGenerator{
		outputDir: outputDir,
	}
}

// GenerateVersionUpdateReport creates a comprehensive version update report
func (r *ReportGenerator) GenerateVersionUpdateReport(report *models.VersionReport, format string) (*models.UpdateReport, error) {
	updateReport := &models.UpdateReport{
		GeneratedAt:     time.Now(),
		ReportID:        r.generateReportID(),
		Type:            "version_update",
		Summary:         r.generateSummary(report),
		Details:         r.generateDetails(report),
		Recommendations: report.Recommendations,
		Metadata:        make(map[string]string),
	}

	// Add metadata
	updateReport.Metadata["total_packages"] = fmt.Sprintf("%d", report.TotalPackages)
	updateReport.Metadata["outdated_count"] = fmt.Sprintf("%d", report.OutdatedCount)
	updateReport.Metadata["security_issues"] = fmt.Sprintf("%d", report.SecurityIssues)
	updateReport.Metadata["format"] = format

	// Save report to file
	if err := r.saveReport(updateReport, format); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return updateReport, nil
}

// GenerateTemplateUpdateReport creates a report for template updates
func (r *ReportGenerator) GenerateTemplateUpdateReport(updates []models.TemplateUpdate, format string) (*models.UpdateReport, error) {
	updateReport := &models.UpdateReport{
		GeneratedAt:     time.Now(),
		ReportID:        r.generateReportID(),
		Type:            "template_update",
		Summary:         r.generateTemplateUpdateSummary(updates),
		TemplateUpdates: updates,
		Metadata:        make(map[string]string),
	}

	// Add metadata
	updateReport.Metadata["total_templates"] = fmt.Sprintf("%d", len(updates))
	updateReport.Metadata["format"] = format

	// Save report to file
	if err := r.saveReport(updateReport, format); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return updateReport, nil
}

// GenerateSecurityReport creates a security-focused report
func (r *ReportGenerator) GenerateSecurityReport(report *models.VersionReport, format string) (*models.SecurityReport, error) {
	securityReport := &models.SecurityReport{
		GeneratedAt:    time.Now(),
		ReportID:       r.generateReportID(),
		TotalIssues:    report.SecurityIssues,
		CriticalIssues: r.countCriticalIssues(report),
		HighIssues:     r.countHighIssues(report),
		Issues:         r.extractSecurityIssues(report),
		Metadata:       make(map[string]string),
	}

	// Add metadata
	securityReport.Metadata["scan_time"] = report.LastUpdateCheck.Format(time.RFC3339)
	securityReport.Metadata["format"] = format

	// Save report to file
	if err := r.saveSecurityReport(securityReport, format); err != nil {
		return nil, fmt.Errorf("failed to save security report: %w", err)
	}

	return securityReport, nil
}

// generateReportID creates a unique report identifier
func (r *ReportGenerator) generateReportID() string {
	return fmt.Sprintf("report_%d", time.Now().Unix())
}

// generateSummary creates a summary from version report
func (r *ReportGenerator) generateSummary(report *models.VersionReport) string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Version Update Report - %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	summary.WriteString(strings.Repeat("=", 50) + "\n\n")

	summary.WriteString(fmt.Sprintf("Total Packages: %d\n", report.TotalPackages))
	summary.WriteString(fmt.Sprintf("Outdated Packages: %d\n", report.OutdatedCount))
	summary.WriteString(fmt.Sprintf("Security Issues: %d\n", report.SecurityIssues))
	summary.WriteString(fmt.Sprintf("Last Check: %s\n\n", report.LastUpdateCheck.Format("2006-01-02 15:04:05")))

	if len(report.Summary) > 0 {
		summary.WriteString("Summary by Category:\n")
		for category, categorySum := range report.Summary {
			summary.WriteString(fmt.Sprintf("  %s: %d total, %d outdated, %d insecure\n",
				strings.Title(category), categorySum.Total, categorySum.Outdated, categorySum.Insecure))
		}
		summary.WriteString("\n")
	}

	if len(report.Recommendations) > 0 {
		summary.WriteString("Priority Updates:\n")
		criticalCount := 0
		highCount := 0
		for _, rec := range report.Recommendations {
			if rec.Priority == "critical" {
				criticalCount++
			} else if rec.Priority == "high" {
				highCount++
			}
		}
		summary.WriteString(fmt.Sprintf("  Critical: %d\n", criticalCount))
		summary.WriteString(fmt.Sprintf("  High: %d\n", highCount))
		summary.WriteString(fmt.Sprintf("  Total Recommendations: %d\n", len(report.Recommendations)))
	}

	return summary.String()
}

// generateDetails creates detailed information from version report
func (r *ReportGenerator) generateDetails(report *models.VersionReport) string {
	var details strings.Builder

	details.WriteString("Detailed Version Information:\n")
	details.WriteString(strings.Repeat("-", 40) + "\n\n")

	if len(report.Recommendations) > 0 {
		details.WriteString("Update Recommendations:\n")
		for _, rec := range report.Recommendations {
			priority := rec.Priority
			if rec.Priority == "critical" {
				priority = "🔴 CRITICAL"
			} else if rec.Priority == "high" {
				priority = "🟠 HIGH"
			} else if rec.Priority == "medium" {
				priority = "🟡 MEDIUM"
			} else {
				priority = "🟢 LOW"
			}

			details.WriteString(fmt.Sprintf("%s %s:\n", priority, rec.Name))
			details.WriteString(fmt.Sprintf("  Current: %s\n", rec.CurrentVersion))
			details.WriteString(fmt.Sprintf("  Recommended: %s\n", rec.RecommendedVersion))
			details.WriteString(fmt.Sprintf("  Reason: %s\n", rec.Reason))
			if rec.BreakingChange {
				details.WriteString("  ⚠️  Breaking Change\n")
			}
			details.WriteString("\n")
		}
	}

	return details.String()
}

// generateTemplateUpdateSummary creates summary for template updates
func (r *ReportGenerator) generateTemplateUpdateSummary(updates []models.TemplateUpdate) string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Template Update Report - %s\n", time.Now().Format("2006-01-02 15:04:05")))
	summary.WriteString(strings.Repeat("=", 50) + "\n\n")

	summary.WriteString(fmt.Sprintf("Total Templates Updated: %d\n", len(updates)))

	successCount := 0
	failureCount := 0
	for _, update := range updates {
		if update.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	summary.WriteString(fmt.Sprintf("Successful Updates: %d\n", successCount))
	summary.WriteString(fmt.Sprintf("Failed Updates: %d\n", failureCount))

	return summary.String()
}

// countCriticalIssues counts critical security issues
func (r *ReportGenerator) countCriticalIssues(report *models.VersionReport) int {
	count := 0
	for _, info := range report.Details {
		for _, issue := range info.SecurityIssues {
			if issue.Severity == "critical" {
				count++
			}
		}
	}
	return count
}

// countHighIssues counts high severity security issues
func (r *ReportGenerator) countHighIssues(report *models.VersionReport) int {
	count := 0
	for _, info := range report.Details {
		for _, issue := range info.SecurityIssues {
			if issue.Severity == "high" {
				count++
			}
		}
	}
	return count
}

// extractSecurityIssues extracts all security issues from version report
func (r *ReportGenerator) extractSecurityIssues(report *models.VersionReport) []models.SecurityIssueDetail {
	var issues []models.SecurityIssueDetail

	for packageName, info := range report.Details {
		for _, issue := range info.SecurityIssues {
			issues = append(issues, models.SecurityIssueDetail{
				PackageName:    packageName,
				CurrentVersion: info.CurrentVersion,
				SecurityIssue:  issue,
				RecommendedFix: issue.FixedIn,
			})
		}
	}

	return issues
}

// saveReport saves the update report to file
func (r *ReportGenerator) saveReport(report *models.UpdateReport, format string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("%s.%s", report.ReportID, format)
	filepath := filepath.Join(r.outputDir, filename)

	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(report, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(report)
	default:
		// Default to text format
		data = []byte(report.Summary + "\n" + report.Details)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

// saveSecurityReport saves the security report to file
func (r *ReportGenerator) saveSecurityReport(report *models.SecurityReport, format string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("security_%s.%s", report.ReportID, format)
	filepath := filepath.Join(r.outputDir, filename)

	var data []byte
	var err error

	switch format {
	case "json":
		data, err = json.MarshalIndent(report, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(report)
	default:
		// Default to text format
		var summary strings.Builder
		summary.WriteString(fmt.Sprintf("Security Report - %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
		summary.WriteString(strings.Repeat("=", 50) + "\n\n")
		summary.WriteString(fmt.Sprintf("Total Issues: %d\n", report.TotalIssues))
		summary.WriteString(fmt.Sprintf("Critical: %d\n", report.CriticalIssues))
		summary.WriteString(fmt.Sprintf("High: %d\n", report.HighIssues))
		summary.WriteString("\nDetailed Issues:\n")

		for _, issue := range report.Issues {
			summary.WriteString(fmt.Sprintf("\n%s (%s):\n", issue.PackageName, issue.CurrentVersion))
			summary.WriteString(fmt.Sprintf("  Issue: %s\n", issue.SecurityIssue.Description))
			summary.WriteString(fmt.Sprintf("  Severity: %s\n", issue.SecurityIssue.Severity))
			if issue.RecommendedFix != "" {
				summary.WriteString(fmt.Sprintf("  Fix: Update to %s\n", issue.RecommendedFix))
			}
		}

		data = []byte(summary.String())
	}

	if err != nil {
		return fmt.Errorf("failed to marshal security report: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write security report file: %w", err)
	}

	return nil
}

// GetReportHistory returns a list of previously generated reports
func (r *ReportGenerator) GetReportHistory() ([]models.ReportInfo, error) {
	var reports []models.ReportInfo

	if _, err := os.Stat(r.outputDir); os.IsNotExist(err) {
		return reports, nil // No reports directory exists yet
	}

	entries, err := os.ReadDir(r.outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read reports directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Parse filename to extract report info
		name := entry.Name()
		var reportType string
		if strings.HasPrefix(name, "security_") {
			reportType = "security"
		} else if strings.HasPrefix(name, "report_") {
			reportType = "version_update"
		} else {
			continue
		}

		reports = append(reports, models.ReportInfo{
			Filename:    name,
			Type:        reportType,
			GeneratedAt: info.ModTime(),
			Size:        info.Size(),
		})
	}

	return reports, nil
}
