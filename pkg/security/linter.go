package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SecurityLinter provides automated security validation and linting
type SecurityLinter struct {
	scanner *Scanner
	fixer   *Fixer
	rules   []LintRule
}

// LintRule defines a security linting rule
type LintRule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Pattern     *regexp.Regexp `json:"-"`
	PatternStr  string         `json:"pattern"`
	Category    string         `json:"category"`
	Severity    SeverityLevel  `json:"severity"`
	Description string         `json:"description"`
	Message     string         `json:"message"`
	Suggestion  string         `json:"suggestion"`
	Enabled     bool           `json:"enabled"`
	Tags        []string       `json:"tags"`
}

// LintResult contains the results of security linting
type LintResult struct {
	Issues       []LintIssue `json:"issues"`
	Summary      LintSummary `json:"summary"`
	ScannedFiles int         `json:"scanned_files"`
	RulesApplied int         `json:"rules_applied"`
}

// LintIssue represents a security issue found by the linter
type LintIssue struct {
	RuleID      string        `json:"rule_id"`
	FilePath    string        `json:"file_path"`
	LineNumber  int           `json:"line_number"`
	Column      int           `json:"column"`
	Severity    SeverityLevel `json:"severity"`
	Category    string        `json:"category"`
	Message     string        `json:"message"`
	Suggestion  string        `json:"suggestion"`
	LineContent string        `json:"line_content"`
	Tags        []string      `json:"tags"`
}

// LintSummary provides a summary of linting results
type LintSummary struct {
	TotalIssues     int                   `json:"total_issues"`
	BySeverity      map[SeverityLevel]int `json:"by_severity"`
	ByCategory      map[string]int        `json:"by_category"`
	CriticalFiles   []string              `json:"critical_files"`
	MostCommonRules []string              `json:"most_common_rules"`
}

// NewSecurityLinter creates a new security linter with predefined rules
func NewSecurityLinter() *SecurityLinter {
	return &SecurityLinter{
		scanner: NewScanner(),
		fixer:   NewFixer(),
		rules:   getSecurityLintRules(),
	}
}

// LintDirectory performs security linting on all files in the given directory
func (sl *SecurityLinter) LintDirectory(dir string) (*LintResult, error) {
	result := &LintResult{
		Issues: make([]LintIssue, 0),
		Summary: LintSummary{
			BySeverity: make(map[SeverityLevel]int),
			ByCategory: make(map[string]int),
		},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-relevant files
		if info.IsDir() || !sl.shouldLintFile(path) {
			return nil
		}

		fileIssues, err := sl.LintFile(path)
		if err != nil {
			return fmt.Errorf("error linting file %s: %w", path, err)
		}

		result.Issues = append(result.Issues, fileIssues...)
		result.ScannedFiles++

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Generate summary
	sl.generateSummary(result)
	result.RulesApplied = len(sl.rules)

	return result, nil
}

// LintFile performs security linting on a single file
func (sl *SecurityLinter) LintFile(filePath string) ([]LintIssue, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var issues []LintIssue

	for lineNum, line := range lines {
		for _, rule := range sl.rules {
			if !rule.Enabled {
				continue
			}

			if matches := rule.Pattern.FindStringSubmatch(line); matches != nil {
				issue := LintIssue{
					RuleID:      rule.ID,
					FilePath:    filePath,
					LineNumber:  lineNum + 1,
					Column:      strings.Index(line, matches[0]) + 1,
					Severity:    rule.Severity,
					Category:    rule.Category,
					Message:     rule.Message,
					Suggestion:  rule.Suggestion,
					LineContent: strings.TrimSpace(line),
					Tags:        rule.Tags,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues, nil
}

// shouldLintFile determines if a file should be linted
func (sl *SecurityLinter) shouldLintFile(path string) bool {
	ext := filepath.Ext(path)
	lintableExts := []string{".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".cs", ".php", ".rb", ".tmpl", ".yaml", ".yml", ".json"}

	for _, lintExt := range lintableExts {
		if ext == lintExt {
			return true
		}
	}

	return false
}

// generateSummary generates a summary of the linting results
func (sl *SecurityLinter) generateSummary(result *LintResult) {
	result.Summary.TotalIssues = len(result.Issues)

	// Count by severity and category
	ruleCount := make(map[string]int)
	criticalFiles := make(map[string]bool)

	for _, issue := range result.Issues {
		result.Summary.BySeverity[issue.Severity]++
		result.Summary.ByCategory[issue.Category]++
		ruleCount[issue.RuleID]++

		if issue.Severity == SeverityCritical {
			criticalFiles[issue.FilePath] = true
		}
	}

	// Extract critical files
	for file := range criticalFiles {
		result.Summary.CriticalFiles = append(result.Summary.CriticalFiles, file)
	}

	// Find most common rules (top 5)
	type ruleFreq struct {
		rule  string
		count int
	}
	var frequencies []ruleFreq
	for rule, count := range ruleCount {
		frequencies = append(frequencies, ruleFreq{rule, count})
	}

	// Simple sort by count (descending)
	for i := 0; i < len(frequencies); i++ {
		for j := i + 1; j < len(frequencies); j++ {
			if frequencies[j].count > frequencies[i].count {
				frequencies[i], frequencies[j] = frequencies[j], frequencies[i]
			}
		}
	}

	// Take top 5
	maxRules := 5
	if len(frequencies) < maxRules {
		maxRules = len(frequencies)
	}
	for i := 0; i < maxRules; i++ {
		result.Summary.MostCommonRules = append(result.Summary.MostCommonRules, frequencies[i].rule)
	}
}

// ExportResults exports linting results to various formats
func (sl *SecurityLinter) ExportResults(result *LintResult, format, outputPath string) error {
	switch strings.ToLower(format) {
	case "json":
		return sl.exportJSON(result, outputPath)
	case "sarif":
		return sl.exportSARIF(result, outputPath)
	case "junit":
		return sl.exportJUnit(result, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportJSON exports results in JSON format
func (sl *SecurityLinter) exportJSON(result *LintResult, outputPath string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// exportSARIF exports results in SARIF format for GitHub Security tab
func (sl *SecurityLinter) exportSARIF(result *LintResult, outputPath string) error {
	sarif := map[string]interface{}{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": []map[string]interface{}{
			{
				"tool": map[string]interface{}{
					"driver": map[string]interface{}{
						"name":    "Security Linter",
						"version": "1.0.0",
						"rules":   sl.getSARIFRules(),
					},
				},
				"results": sl.getSARIFResults(result),
			},
		},
	}

	data, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SARIF: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// exportJUnit exports results in JUnit XML format for CI/CD integration
func (sl *SecurityLinter) exportJUnit(result *LintResult, outputPath string) error {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="Security Linter" tests="%d" failures="%d" errors="0" time="0">
%s
</testsuite>`

	var testCases strings.Builder
	failures := 0

	for _, issue := range result.Issues {
		if issue.Severity == SeverityCritical || issue.Severity == SeverityHigh {
			failures++
			testCases.WriteString(fmt.Sprintf(`  <testcase classname="%s" name="%s" time="0">
    <failure message="%s" type="%s">%s:%d - %s</failure>
  </testcase>
`, issue.Category, issue.RuleID, issue.Message, issue.Severity, issue.FilePath, issue.LineNumber, issue.Suggestion))
		} else {
			testCases.WriteString(fmt.Sprintf(`  <testcase classname="%s" name="%s" time="0"/>
`, issue.Category, issue.RuleID))
		}
	}

	finalXML := fmt.Sprintf(xml, len(result.Issues), failures, testCases.String())
	return os.WriteFile(outputPath, []byte(finalXML), 0644)
}

// getSARIFRules converts linting rules to SARIF format
func (sl *SecurityLinter) getSARIFRules() []map[string]interface{} {
	var rules []map[string]interface{}

	for _, rule := range sl.rules {
		sarifRule := map[string]interface{}{
			"id":   rule.ID,
			"name": rule.Name,
			"shortDescription": map[string]interface{}{
				"text": rule.Description,
			},
			"fullDescription": map[string]interface{}{
				"text": rule.Message,
			},
			"defaultConfiguration": map[string]interface{}{
				"level": sl.severityToSARIFLevel(rule.Severity),
			},
			"properties": map[string]interface{}{
				"category": rule.Category,
				"tags":     rule.Tags,
			},
		}
		rules = append(rules, sarifRule)
	}

	return rules
}

// getSARIFResults converts linting issues to SARIF format
func (sl *SecurityLinter) getSARIFResults(result *LintResult) []map[string]interface{} {
	var results []map[string]interface{}

	for _, issue := range result.Issues {
		sarifResult := map[string]interface{}{
			"ruleId": issue.RuleID,
			"level":  sl.severityToSARIFLevel(issue.Severity),
			"message": map[string]interface{}{
				"text": issue.Message,
			},
			"locations": []map[string]interface{}{
				{
					"physicalLocation": map[string]interface{}{
						"artifactLocation": map[string]interface{}{
							"uri": issue.FilePath,
						},
						"region": map[string]interface{}{
							"startLine":   issue.LineNumber,
							"startColumn": issue.Column,
						},
					},
				},
			},
		}
		results = append(results, sarifResult)
	}

	return results
}

// severityToSARIFLevel converts our severity levels to SARIF levels
func (sl *SecurityLinter) severityToSARIFLevel(severity SeverityLevel) string {
	switch severity {
	case SeverityCritical:
		return "error"
	case SeverityHigh:
		return "error"
	case SeverityMedium:
		return "warning"
	case SeverityLow:
		return "note"
	default:
		return "warning"
	}
}

// HasCriticalIssues returns true if the result contains critical security issues
func (result *LintResult) HasCriticalIssues() bool {
	return result.Summary.BySeverity[SeverityCritical] > 0
}

// HasHighSeverityIssues returns true if the result contains high severity issues
func (result *LintResult) HasHighSeverityIssues() bool {
	return result.Summary.BySeverity[SeverityCritical] > 0 || result.Summary.BySeverity[SeverityHigh] > 0
}

// GetIssuesByFile returns issues grouped by file path
func (result *LintResult) GetIssuesByFile() map[string][]LintIssue {
	byFile := make(map[string][]LintIssue)
	for _, issue := range result.Issues {
		byFile[issue.FilePath] = append(byFile[issue.FilePath], issue)
	}
	return byFile
}
