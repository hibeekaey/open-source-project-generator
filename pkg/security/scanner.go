package security

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// SecurityIssueType represents the type of security issue
type SecurityIssueType string

const (
	CORSVulnerability     SecurityIssueType = "cors_vulnerability"
	MissingSecurityHeader SecurityIssueType = "missing_security_header"
	WeakAuthentication    SecurityIssueType = "weak_authentication"
	SQLInjectionRisk      SecurityIssueType = "sql_injection_risk"
	InformationLeakage    SecurityIssueType = "information_leakage"
)

// SeverityLevel represents the severity of a security issue
type SeverityLevel string

const (
	SeverityCritical SeverityLevel = "critical"
	SeverityHigh     SeverityLevel = "high"
	SeverityMedium   SeverityLevel = "medium"
	SeverityLow      SeverityLevel = "low"
)

// SecurityIssue represents a security vulnerability found in template files
type SecurityIssue struct {
	FilePath       string            `json:"file_path"`
	LineNumber     int               `json:"line_number"`
	IssueType      SecurityIssueType `json:"issue_type"`
	Severity       SeverityLevel     `json:"severity"`
	Description    string            `json:"description"`
	Recommendation string            `json:"recommendation"`
	FixAvailable   bool              `json:"fix_available"`
}

// SecurityReport contains the results of a security scan
type SecurityReport struct {
	Issues       []SecurityIssue `json:"issues"`
	ScannedFiles int             `json:"scanned_files"`
	TotalLines   int             `json:"total_lines"`
}

// HasCriticalIssues returns true if the report contains critical security issues
func (r *SecurityReport) HasCriticalIssues() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityCritical {
			return true
		}
	}
	return false
}

// CountBySeverity returns the number of issues with the given severity
func (r *SecurityReport) CountBySeverity(severity SeverityLevel) int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

// CountByType returns the number of issues with the given type
func (r *SecurityReport) CountByType(issueType SecurityIssueType) int {
	count := 0
	for _, issue := range r.Issues {
		if issue.IssueType == issueType {
			count++
		}
	}
	return count
}

// GetFixableIssues returns only issues that have fixes available
func (r *SecurityReport) GetFixableIssues() []SecurityIssue {
	var fixable []SecurityIssue
	for _, issue := range r.Issues {
		if issue.FixAvailable {
			fixable = append(fixable, issue)
		}
	}
	return fixable
}

// FilterBySeverity returns issues with severity at or above the specified level
func (r *SecurityReport) FilterBySeverity(minSeverity SeverityLevel) []SecurityIssue {
	severityOrder := map[SeverityLevel]int{
		SeverityLow:      1,
		SeverityMedium:   2,
		SeverityHigh:     3,
		SeverityCritical: 4,
	}

	minLevel := severityOrder[minSeverity]
	var filtered []SecurityIssue

	for _, issue := range r.Issues {
		if severityOrder[issue.Severity] >= minLevel {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// Scanner scans template files for security vulnerabilities
type Scanner struct {
	patterns []SecurityPattern
}

// SecurityPattern defines a pattern to match security issues
type SecurityPattern struct {
	Name           string
	Pattern        *regexp.Regexp
	IssueType      SecurityIssueType
	Severity       SeverityLevel
	Description    string
	Recommendation string
	FixAvailable   bool
}

// NewScanner creates a new security scanner with predefined patterns
func NewScanner() *Scanner {
	return &Scanner{
		patterns: getSecurityPatterns(),
	}
}

// ScanDirectory scans all template files in the given directory
func (s *Scanner) ScanDirectory(dir string) (*SecurityReport, error) {
	report := &SecurityReport{
		Issues: make([]SecurityIssue, 0),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || !isTemplateFile(path) {
			return nil
		}

		issues, err := s.ScanFile(path)
		if err != nil {
			return fmt.Errorf("error scanning file %s: %w", path, err)
		}

		report.Issues = append(report.Issues, issues...)
		report.ScannedFiles++

		return nil
	})

	return report, err
}

// ScanFile scans a single file for security issues
func (s *Scanner) ScanFile(filePath string) ([]SecurityIssue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issues []SecurityIssue
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		for _, pattern := range s.patterns {
			if pattern.Pattern.MatchString(line) {
				issue := SecurityIssue{
					FilePath:       filePath,
					LineNumber:     lineNumber,
					IssueType:      pattern.IssueType,
					Severity:       pattern.Severity,
					Description:    pattern.Description,
					Recommendation: pattern.Recommendation,
					FixAvailable:   pattern.FixAvailable,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues, scanner.Err()
}

// isTemplateFile checks if the file is a template file
func isTemplateFile(path string) bool {
	ext := filepath.Ext(path)
	templateExts := []string{".tmpl", ".go", ".js", ".ts", ".yaml", ".yml", ".json"}

	for _, templateExt := range templateExts {
		if ext == templateExt {
			return true
		}
	}

	return false
}
