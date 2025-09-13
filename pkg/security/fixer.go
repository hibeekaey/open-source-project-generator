package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FixerOptions configures the behavior of the security fixer
type FixerOptions struct {
	DryRun       bool
	Verbose      bool
	FixType      string
	CreateBackup bool
}

// FixedIssue represents a security issue that was successfully fixed
type FixedIssue struct {
	FilePath        string            `json:"file_path"`
	LineNumber      int               `json:"line_number"`
	IssueType       SecurityIssueType `json:"issue_type"`
	Description     string            `json:"description"`
	FixDescription  string            `json:"fix_description"`
	OriginalContent string            `json:"original_content"`
	FixedContent    string            `json:"fixed_content"`
}

// FixError represents an error that occurred during fixing
type FixError struct {
	FilePath string `json:"file_path"`
	Error    string `json:"error"`
}

// FixResult contains the results of applying security fixes
type FixResult struct {
	FixedIssues    []FixedIssue `json:"fixed_issues"`
	Errors         []FixError   `json:"errors"`
	BackupsCreated int          `json:"backups_created"`
}

// HasErrors returns true if the fix result contains errors
func (r *FixResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// Fixer applies automated security fixes to template files
type Fixer struct {
	fixes []SecurityFix
}

// SecurityFix defines how to automatically fix a security issue
type SecurityFix struct {
	Name           string
	Pattern        *regexp.Regexp
	IssueType      SecurityIssueType
	Description    string
	FixDescription string
	FixFunction    func(line string) string
	Enabled        bool
}

// NewFixer creates a new security fixer with predefined fixes
func NewFixer() *Fixer {
	return &Fixer{
		fixes: getSecurityFixes(),
	}
}

// FixDirectory applies security fixes to all template files in the given directory
func (f *Fixer) FixDirectory(dir string, options FixerOptions) (*FixResult, error) {
	result := &FixResult{
		FixedIssues: make([]FixedIssue, 0),
		Errors:      make([]FixError, 0),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-template files
		if info.IsDir() || !isTemplateFile(path) {
			return nil
		}

		fileResult, err := f.FixFile(path, options)
		if err != nil {
			result.Errors = append(result.Errors, FixError{
				FilePath: path,
				Error:    err.Error(),
			})
			return nil // Continue processing other files
		}

		result.FixedIssues = append(result.FixedIssues, fileResult.FixedIssues...)
		result.BackupsCreated += fileResult.BackupsCreated

		return nil
	})

	return result, err
}

// FixFile applies security fixes to a single file
func (f *Fixer) FixFile(filePath string, options FixerOptions) (*FixResult, error) {
	result := &FixResult{
		FixedIssues: make([]FixedIssue, 0),
		Errors:      make([]FixError, 0),
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	modified := false

	// Apply fixes line by line
	for lineNum, line := range lines {
		for _, fix := range f.fixes {
			if !fix.Enabled || !f.shouldApplyFix(fix.IssueType, options.FixType) {
				continue
			}

			if fix.Pattern.MatchString(line) {
				originalLine := line
				fixedLine := fix.FixFunction(line)

				if fixedLine != originalLine {
					lines[lineNum] = fixedLine
					modified = true

					fixedIssue := FixedIssue{
						FilePath:        filePath,
						LineNumber:      lineNum + 1,
						IssueType:       fix.IssueType,
						Description:     fix.Description,
						FixDescription:  fix.FixDescription,
						OriginalContent: originalLine,
						FixedContent:    fixedLine,
					}
					result.FixedIssues = append(result.FixedIssues, fixedIssue)

					if options.Verbose {
						fmt.Printf("Fixed %s:%d - %s\n", filePath, lineNum+1, fix.Description)
					}
				}
			}
		}
	}

	// Write the fixed content back to the file
	if modified && !options.DryRun {
		// Create backup if requested
		if options.CreateBackup {
			backupPath := filePath + ".backup." + time.Now().Format("20060102-150405")
			if err := os.WriteFile(backupPath, content, 0644); err != nil {
				return result, fmt.Errorf("failed to create backup: %w", err)
			}
			result.BackupsCreated++
		}

		// Write the fixed content
		fixedContent := strings.Join(lines, "\n")
		if err := os.WriteFile(filePath, []byte(fixedContent), 0644); err != nil {
			return result, fmt.Errorf("failed to write fixed content: %w", err)
		}
	}

	return result, nil
}

// shouldApplyFix determines if a fix should be applied based on the fix type filter
func (f *Fixer) shouldApplyFix(issueType SecurityIssueType, fixType string) bool {
	if fixType == "all" {
		return true
	}

	switch fixType {
	case "cors":
		return issueType == CORSVulnerability
	case "headers":
		return issueType == MissingSecurityHeader
	case "auth":
		return issueType == WeakAuthentication
	case "sql":
		return issueType == SQLInjectionRisk
	default:
		return true
	}
}
