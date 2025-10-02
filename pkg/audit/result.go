// Package audit provides result aggregation functionality for the audit engine.
//
// This file contains the ResultAggregator struct and related functionality for
// aggregating audit results, calculating scores, and generating recommendations.
package audit

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ResultAggregator manages audit result aggregation and report generation.
//
// It provides functionality for:
//   - Result merging and aggregation
//   - Score calculation for different audit types
//   - Recommendation generation based on audit results
//   - Report generation in multiple formats
type ResultAggregator struct {
	// Configuration for scoring weights
	securityWeight    float64
	qualityWeight     float64
	licenseWeight     float64
	performanceWeight float64
}

// NewResultAggregator creates a new result aggregator with default weights.
func NewResultAggregator() *ResultAggregator {
	return &ResultAggregator{
		securityWeight:    0.3, // 30% weight for security
		qualityWeight:     0.3, // 30% weight for quality
		licenseWeight:     0.2, // 20% weight for license
		performanceWeight: 0.2, // 20% weight for performance
	}
}

// NewResultAggregatorWithWeights creates a new result aggregator with custom weights.
func NewResultAggregatorWithWeights(security, quality, license, performance float64) *ResultAggregator {
	return &ResultAggregator{
		securityWeight:    security,
		qualityWeight:     quality,
		licenseWeight:     license,
		performanceWeight: performance,
	}
}

// AggregateResults aggregates multiple audit results into a summary.
func (ra *ResultAggregator) AggregateResults(results []*interfaces.AuditResult) (*interfaces.AuditSummary, error) {
	if len(results) == 0 {
		return &interfaces.AuditSummary{}, nil
	}

	summary := &interfaces.AuditSummary{
		TotalProjects: len(results),
	}

	var totalScore, securityScore, qualityScore, licenseScore, performanceScore float64
	var securityCount, qualityCount, licenseCount, performanceCount int

	issueFrequency := make(map[string]int)

	for _, result := range results {
		totalScore += result.OverallScore

		if result.Security != nil {
			securityScore += result.Security.Score
			securityCount++

			// Count security issues
			for _, vuln := range result.Security.Vulnerabilities {
				key := fmt.Sprintf("security-%s", vuln.Severity)
				issueFrequency[key]++
			}
		}

		if result.Quality != nil {
			qualityScore += result.Quality.Score
			qualityCount++

			// Count quality issues
			for _, smell := range result.Quality.CodeSmells {
				key := fmt.Sprintf("quality-%s", smell.Type)
				issueFrequency[key]++
			}
		}

		if result.Licenses != nil {
			licenseScore += result.Licenses.Score
			licenseCount++

			// Count license conflicts
			for range result.Licenses.Conflicts {
				issueFrequency["license-conflict"]++
			}
		}

		if result.Performance != nil {
			performanceScore += result.Performance.Score
			performanceCount++

			// Count performance issues
			for _, issue := range result.Performance.Issues {
				key := fmt.Sprintf("performance-%s", issue.Type)
				issueFrequency[key]++
			}
		}
	}

	// Calculate average scores
	summary.AverageScore = totalScore / float64(len(results))

	if securityCount > 0 {
		summary.SecurityScore = securityScore / float64(securityCount)
	}
	if qualityCount > 0 {
		summary.QualityScore = qualityScore / float64(qualityCount)
	}
	if licenseCount > 0 {
		summary.LicenseScore = licenseScore / float64(licenseCount)
	}
	if performanceCount > 0 {
		summary.PerformanceScore = performanceScore / float64(performanceCount)
	}

	// Generate common issues list (only for multiple projects)
	if len(results) > 1 {
		for issueType, frequency := range issueFrequency {
			if frequency > len(results)/4 { // Issues that appear in more than 25% of projects
				summary.CommonIssues = append(summary.CommonIssues, interfaces.CommonIssue{
					Type:        issueType,
					Description: ra.getIssueDescription(issueType),
					Frequency:   frequency,
					Severity:    ra.getIssueSeverity(issueType),
					Category:    ra.getIssueCategory(issueType),
				})
			}
		}
	}

	return summary, nil
}

// CalculateOverallScore calculates the overall score for an audit result.
func (ra *ResultAggregator) CalculateOverallScore(result *interfaces.AuditResult) float64 {
	var totalScore float64
	var totalWeight float64

	if result.Security != nil {
		totalScore += result.Security.Score * ra.securityWeight
		totalWeight += ra.securityWeight
	}

	if result.Quality != nil {
		totalScore += result.Quality.Score * ra.qualityWeight
		totalWeight += ra.qualityWeight
	}

	if result.Licenses != nil {
		totalScore += result.Licenses.Score * ra.licenseWeight
		totalWeight += ra.licenseWeight
	}

	if result.Performance != nil {
		totalScore += result.Performance.Score * ra.performanceWeight
		totalWeight += ra.performanceWeight
	}

	if totalWeight > 0 {
		return totalScore / totalWeight
	}

	return 0
}

// CalculateSecurityScore calculates the overall security score.
func (ra *ResultAggregator) CalculateSecurityScore(secResult *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) float64 {
	score := 100.0

	// Deduct points for vulnerabilities
	for _, vuln := range secResult.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			score -= 20
		case "high":
			score -= 10
		case "medium":
			score -= 5
		case "low":
			score -= 2
		}
	}

	// Deduct points for policy violations
	for _, violation := range secResult.PolicyViolations {
		switch violation.Severity {
		case "critical":
			score -= 15
		case "high":
			score -= 8
		case "medium":
			score -= 4
		case "low":
			score -= 1
		}
	}

	// Deduct points for high-confidence secrets
	if secretResult != nil {
		for _, secret := range secretResult.Secrets {
			if secret.Confidence >= 0.8 {
				score -= 10
			} else if secret.Confidence >= 0.6 {
				score -= 5
			}
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// CalculateQualityScore calculates the overall quality score.
func (ra *ResultAggregator) CalculateQualityScore(result *interfaces.QualityAuditResult) float64 {
	score := 100.0

	// Deduct points for code smells
	for _, smell := range result.CodeSmells {
		switch smell.Severity {
		case "high":
			score -= 5
		case "medium":
			score -= 3
		case "low":
			score -= 1
		}
	}

	// Deduct points for duplications
	for _, dup := range result.Duplications {
		if dup.Percentage > 10 {
			score -= 10
		} else if dup.Percentage > 5 {
			score -= 5
		}
	}

	// Adjust score based on test coverage
	if result.TestCoverage < 50 {
		score -= 20
	} else if result.TestCoverage < 80 {
		score -= 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

// CalculatePerformanceScore calculates the performance score.
func (ra *ResultAggregator) CalculatePerformanceScore(result *interfaces.PerformanceAuditResult) float64 {
	score := 100.0

	// Deduct points for bundle size
	sizeInMB := float64(result.BundleSize) / (1024 * 1024)
	if sizeInMB > 10 {
		score -= 30
	} else if sizeInMB > 5 {
		score -= 20
	} else if sizeInMB > 2 {
		score -= 10
	}

	// Deduct points for load time
	if result.LoadTime > 5*time.Second {
		score -= 40
	} else if result.LoadTime > 3*time.Second {
		score -= 25
	} else if result.LoadTime > 1*time.Second {
		score -= 10
	}

	// Deduct points for performance issues
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "high":
			score -= 10
		case "medium":
			score -= 5
		case "low":
			score -= 2
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

// GenerateRecommendations generates recommendations based on audit results.
func (ra *ResultAggregator) GenerateRecommendations(result *interfaces.AuditResult) []string {
	var recommendations []string

	// Security recommendations
	if result.Security != nil {
		secRecs := ra.GenerateSecurityRecommendations(result.Security, nil)
		recommendations = append(recommendations, secRecs...)
	}

	// Quality recommendations
	if result.Quality != nil {
		qualRecs := ra.GenerateQualityRecommendations(result.Quality)
		recommendations = append(recommendations, qualRecs...)
	}

	// License recommendations
	if result.Licenses != nil {
		licRecs := ra.GenerateLicenseRecommendations(result.Licenses)
		recommendations = append(recommendations, licRecs...)
	}

	// Performance recommendations
	if result.Performance != nil {
		perfRecs := ra.GeneratePerformanceRecommendations(result.Performance)
		recommendations = append(recommendations, perfRecs...)
	}

	// Overall recommendations
	if result.OverallScore < 70 {
		recommendations = append(recommendations, "Overall project health is below 70% - consider addressing the most critical issues first")
	}

	return recommendations
}

// GenerateSecurityRecommendations generates security recommendations.
func (ra *ResultAggregator) GenerateSecurityRecommendations(secResult *interfaces.SecurityAuditResult, secretResult *interfaces.SecretScanResult) []string {
	var recommendations []string

	if len(secResult.Vulnerabilities) > 0 {
		recommendations = append(recommendations, "Update vulnerable dependencies to secure versions")
		recommendations = append(recommendations, "Implement automated dependency scanning in CI/CD pipeline")
	}

	if len(secResult.PolicyViolations) > 0 {
		recommendations = append(recommendations, "Address security policy violations")
		recommendations = append(recommendations, "Implement security linting in development workflow")
	}

	if secretResult != nil && secretResult.Summary.HighConfidence > 0 {
		recommendations = append(recommendations, "Remove hardcoded secrets from source code")
		recommendations = append(recommendations, "Use environment variables or secret management systems")
	}

	if secResult.Score < 70 {
		recommendations = append(recommendations, "Consider improving security practices - score is below 70%")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Security posture is good - continue following security best practices")
	}

	return recommendations
}

// GenerateQualityRecommendations generates quality recommendations.
func (ra *ResultAggregator) GenerateQualityRecommendations(result *interfaces.QualityAuditResult) []string {
	var recommendations []string

	if len(result.CodeSmells) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d code quality issues", len(result.CodeSmells)))
	}

	if len(result.Duplications) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Refactor %d code duplications", len(result.Duplications)))
	}

	if result.TestCoverage < 80 {
		recommendations = append(recommendations, fmt.Sprintf("Increase test coverage from %.1f%% to at least 80%%", result.TestCoverage))
	}

	if result.Score < 70 {
		recommendations = append(recommendations, "Consider implementing code quality gates in CI/CD pipeline")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Code quality is good - maintain current practices")
	}

	return recommendations
}

// GenerateLicenseRecommendations generates license recommendations.
func (ra *ResultAggregator) GenerateLicenseRecommendations(result *interfaces.LicenseAuditResult) []string {
	var recommendations []string

	if len(result.Conflicts) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Resolve %d license conflicts", len(result.Conflicts)))
		recommendations = append(recommendations, "Consider using a license compatibility tool")
	}

	if !result.Compatible {
		recommendations = append(recommendations, "Review project license compatibility with dependencies")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "License compliance is good")
	}

	return recommendations
}

// GeneratePerformanceRecommendations generates performance recommendations.
func (ra *ResultAggregator) GeneratePerformanceRecommendations(result *interfaces.PerformanceAuditResult) []string {
	var recommendations []string

	sizeInMB := float64(result.BundleSize) / (1024 * 1024)
	if sizeInMB > 5 {
		recommendations = append(recommendations, fmt.Sprintf("Reduce bundle size from %.1fMB", sizeInMB))
		recommendations = append(recommendations, "Consider code splitting and lazy loading")
	}

	if result.LoadTime > 2*time.Second {
		recommendations = append(recommendations, fmt.Sprintf("Optimize load time from %v", result.LoadTime))
		recommendations = append(recommendations, "Enable compression and caching")
	}

	if len(result.Issues) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d performance issues", len(result.Issues)))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance is good - maintain current optimizations")
	}

	return recommendations
}

// GenerateHTMLReport generates an HTML audit report.
func (ra *ResultAggregator) GenerateHTMLReport(result *interfaces.AuditResult) ([]byte, error) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Audit Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f5f5f5; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .score { font-size: 24px; font-weight: bold; }
        .high { color: #d32f2f; }
        .medium { color: #f57c00; }
        .low { color: #388e3c; }
        .recommendations { background-color: #e3f2fd; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Audit Report</h1>
        <p><strong>Project:</strong> %s</p>
        <p><strong>Audit Time:</strong> %s</p>
        <p><strong>Overall Score:</strong> <span class="score">%.1f%%</span></p>
    </div>
`, filepath.Base(result.ProjectPath), result.ProjectPath, result.AuditTime.Format(time.RFC3339), result.OverallScore)

	if result.Security != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Security Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Vulnerabilities:</strong> %d</p>
        <p><strong>Policy Violations:</strong> %d</p>
    </div>`, result.Security.Score, len(result.Security.Vulnerabilities), len(result.Security.PolicyViolations))
	}

	if result.Quality != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Quality Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Code Smells:</strong> %d</p>
        <p><strong>Test Coverage:</strong> %.1f%%</p>
    </div>`, result.Quality.Score, len(result.Quality.CodeSmells), result.Quality.TestCoverage)
	}

	if result.Licenses != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>License Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Compatible:</strong> %t</p>
        <p><strong>Conflicts:</strong> %d</p>
    </div>`, result.Licenses.Score, result.Licenses.Compatible, len(result.Licenses.Conflicts))
	}

	if result.Performance != nil {
		html += fmt.Sprintf(`
    <div class="section">
        <h2>Performance Audit</h2>
        <p><strong>Score:</strong> %.1f%%</p>
        <p><strong>Bundle Size:</strong> %d bytes</p>
        <p><strong>Issues:</strong> %d</p>
    </div>`, result.Performance.Score, result.Performance.BundleSize, len(result.Performance.Issues))
	}

	if len(result.Recommendations) > 0 {
		html += `
    <div class="recommendations">
        <h2>Recommendations</h2>
        <ul>`
		for _, rec := range result.Recommendations {
			html += fmt.Sprintf("<li>%s</li>", rec)
		}
		html += `
        </ul>
    </div>`
	}

	html += `
</body>
</html>`

	return []byte(html), nil
}

// GenerateMarkdownReport generates a Markdown audit report.
func (ra *ResultAggregator) GenerateMarkdownReport(result *interfaces.AuditResult) ([]byte, error) {
	md := fmt.Sprintf(`# Audit Report

**Project:** %s  
**Audit Time:** %s  
**Overall Score:** %.1f%%

`, result.ProjectPath, result.AuditTime.Format(time.RFC3339), result.OverallScore)

	if result.Security != nil {
		md += fmt.Sprintf(`## Security Audit

- **Score:** %.1f%%
- **Vulnerabilities:** %d
- **Policy Violations:** %d

`, result.Security.Score, len(result.Security.Vulnerabilities), len(result.Security.PolicyViolations))
	}

	if result.Quality != nil {
		md += fmt.Sprintf(`## Quality Audit

- **Score:** %.1f%%
- **Code Smells:** %d
- **Test Coverage:** %.1f%%

`, result.Quality.Score, len(result.Quality.CodeSmells), result.Quality.TestCoverage)
	}

	if result.Licenses != nil {
		md += fmt.Sprintf(`## License Audit

- **Score:** %.1f%%
- **Compatible:** %t
- **Conflicts:** %d

`, result.Licenses.Score, result.Licenses.Compatible, len(result.Licenses.Conflicts))
	}

	if result.Performance != nil {
		md += fmt.Sprintf(`## Performance Audit

- **Score:** %.1f%%
- **Bundle Size:** %d bytes
- **Issues:** %d

`, result.Performance.Score, result.Performance.BundleSize, len(result.Performance.Issues))
	}

	if len(result.Recommendations) > 0 {
		md += "## Recommendations\n\n"
		for _, rec := range result.Recommendations {
			md += fmt.Sprintf("- %s\n", rec)
		}
		md += "\n"
	}

	return []byte(md), nil
}

// Helper methods for issue categorization

func (ra *ResultAggregator) getIssueDescription(issueType string) string {
	descriptions := map[string]string{
		"security-critical":       "Critical security vulnerabilities found",
		"security-high":           "High severity security vulnerabilities found",
		"security-medium":         "Medium severity security vulnerabilities found",
		"security-low":            "Low severity security vulnerabilities found",
		"quality-high":            "High severity code quality issues found",
		"quality-medium":          "Medium severity code quality issues found",
		"quality-low":             "Low severity code quality issues found",
		"license-conflict":        "License compatibility conflicts found",
		"performance-large_file":  "Large files affecting performance found",
		"performance-large_asset": "Large assets affecting load time found",
	}

	if desc, exists := descriptions[issueType]; exists {
		return desc
	}
	return fmt.Sprintf("Issue type: %s", issueType)
}

func (ra *ResultAggregator) getIssueSeverity(issueType string) string {
	if len(issueType) > 9 && issueType[:9] == "security-" {
		return issueType[9:] // Extract severity from "security-{severity}"
	}
	if len(issueType) > 8 && issueType[:8] == "quality-" {
		return issueType[8:] // Extract severity from "quality-{severity}"
	}
	if issueType == "license-conflict" {
		return "medium"
	}
	if len(issueType) > 12 && issueType[:12] == "performance-" {
		return "medium"
	}
	return "low"
}

func (ra *ResultAggregator) getIssueCategory(issueType string) string {
	if len(issueType) > 9 && issueType[:9] == "security-" {
		return "security"
	}
	if len(issueType) > 8 && issueType[:8] == "quality-" {
		return "quality"
	}
	if len(issueType) > 8 && issueType[:8] == "license-" {
		return "license"
	}
	if len(issueType) > 12 && issueType[:12] == "performance-" {
		return "performance"
	}
	return "unknown"
}
