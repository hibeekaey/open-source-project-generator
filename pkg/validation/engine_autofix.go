package validation

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// FixValidationIssues applies fixes for validation issues
func (e *Engine) FixValidationIssues(path string, issues []interfaces.ValidationIssue) (*interfaces.FixResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	result := &interfaces.FixResult{
		Applied: []interfaces.Fix{},
		Failed:  []interfaces.FixFailure{},
		Skipped: []interfaces.Fix{},
		Summary: interfaces.FixSummary{
			TotalFixes:    0,
			AppliedFixes:  0,
			FailedFixes:   0,
			SkippedFixes:  0,
			FilesModified: 0,
		},
	}

	modifiedFiles := make(map[string]bool)

	for _, issue := range issues {
		if !issue.Fixable {
			continue
		}

		fix := e.createFixForIssue(issue)
		if fix == nil {
			continue
		}

		result.Summary.TotalFixes++

		if err := e.ApplyFix(path, *fix); err != nil {
			result.Failed = append(result.Failed, interfaces.FixFailure{
				Fix:   *fix,
				Error: err.Error(),
			})
			result.Summary.FailedFixes++
		} else {
			result.Applied = append(result.Applied, *fix)
			result.Summary.AppliedFixes++
			modifiedFiles[fix.File] = true
		}
	}

	result.Summary.FilesModified = len(modifiedFiles)

	return result, nil
}

// GetFixableIssues returns only the issues that can be automatically fixed
func (e *Engine) GetFixableIssues(issues []interfaces.ValidationIssue) []interfaces.ValidationIssue {
	fixableIssues := []interfaces.ValidationIssue{}

	for _, issue := range issues {
		if issue.Fixable {
			fixableIssues = append(fixableIssues, issue)
		}
	}

	return fixableIssues
}

// PreviewFixes shows what fixes would be applied without actually applying them
func (e *Engine) PreviewFixes(path string, issues []interfaces.ValidationIssue) (*interfaces.FixPreview, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	preview := &interfaces.FixPreview{
		Fixes:   []interfaces.Fix{},
		Changes: []interfaces.FileChange{},
		Summary: interfaces.FixSummary{
			TotalFixes:    0,
			AppliedFixes:  0,
			FailedFixes:   0,
			SkippedFixes:  0,
			FilesModified: 0,
		},
	}

	modifiedFiles := make(map[string]bool)

	for _, issue := range issues {
		if !issue.Fixable {
			continue
		}

		fix := e.createFixForIssue(issue)
		if fix == nil {
			continue
		}

		preview.Fixes = append(preview.Fixes, *fix)
		preview.Summary.TotalFixes++

		// Create file change preview
		change, err := e.createFileChangePreview(*fix)
		if err != nil {
			continue
		}

		preview.Changes = append(preview.Changes, *change)
		modifiedFiles[fix.File] = true
	}

	preview.Summary.FilesModified = len(modifiedFiles)

	return preview, nil
}

// ApplyFix applies a single fix
func (e *Engine) ApplyFix(path string, fix interfaces.Fix) error {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(fix.File); err != nil {
		return fmt.Errorf("invalid file path in fix: %w", err)
	}

	switch fix.Action {
	case interfaces.FixActionCreate:
		return e.applyCreateFix(fix)
	case interfaces.FixActionReplace:
		return e.applyReplaceFix(fix)
	case interfaces.FixActionInsert:
		return e.applyInsertFix(fix)
	case interfaces.FixActionDelete:
		return e.applyDeleteFix(fix)
	case interfaces.FixActionRename:
		return e.applyRenameFix(fix)
	default:
		return fmt.Errorf("unsupported fix action: %s", fix.Action)
	}
}

// GenerateValidationReport generates a validation report in the specified format
func (e *Engine) GenerateValidationReport(result *interfaces.ValidationResult, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return e.generateJSONReport(result)
	case "html":
		return e.generateHTMLReport(result)
	case "markdown", "md":
		return e.generateMarkdownReport(result)
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}
}

// GetValidationSummary creates a summary from multiple validation results
func (e *Engine) GetValidationSummary(results []*interfaces.ValidationResult) (*interfaces.ValidationSummary, error) {
	if len(results) == 0 {
		return &interfaces.ValidationSummary{}, nil
	}

	summary := &interfaces.ValidationSummary{
		TotalFiles:   0,
		ValidFiles:   0,
		ErrorCount:   0,
		WarningCount: 0,
		FixableCount: 0,
	}

	for _, result := range results {
		summary.TotalFiles += result.Summary.TotalFiles
		summary.ValidFiles += result.Summary.ValidFiles
		summary.ErrorCount += result.Summary.ErrorCount
		summary.WarningCount += result.Summary.WarningCount
		summary.FixableCount += result.Summary.FixableCount
	}

	return summary, nil
}

// createFixForIssue creates a fix for a specific validation issue
func (e *Engine) createFixForIssue(issue interfaces.ValidationIssue) *interfaces.Fix {
	switch issue.Rule {
	case "structure.readme.required":
		return &interfaces.Fix{
			ID:          fmt.Sprintf("fix_%s_%d", issue.Rule, issue.Line),
			Type:        "create_file",
			Description: "Create README.md file",
			File:        filepath.Join(filepath.Dir(issue.File), "README.md"),
			Action:      interfaces.FixActionCreate,
			Content:     e.generateReadmeContent(),
			Automatic:   true,
		}

	case "structure.license.required":
		return &interfaces.Fix{
			ID:          fmt.Sprintf("fix_%s_%d", issue.Rule, issue.Line),
			Type:        "create_file",
			Description: "Create LICENSE file",
			File:        filepath.Join(filepath.Dir(issue.File), "LICENSE"),
			Action:      interfaces.FixActionCreate,
			Content:     e.generateLicenseContent(),
			Automatic:   true,
		}

	case "quality.naming.conventions":
		if strings.Contains(issue.Message, "space") {
			newName := strings.ReplaceAll(filepath.Base(issue.File), " ", "_")
			newPath := filepath.Join(filepath.Dir(issue.File), newName)
			return &interfaces.Fix{
				ID:          fmt.Sprintf("fix_%s_%d", issue.Rule, issue.Line),
				Type:        "rename_file",
				Description: fmt.Sprintf("Rename file to follow naming conventions: %s", newName),
				File:        issue.File,
				Action:      interfaces.FixActionRename,
				Content:     newPath,
				Automatic:   true,
			}
		}

	case "template.file.extension":
		return &interfaces.Fix{
			ID:          fmt.Sprintf("fix_%s_%d", issue.Rule, issue.Line),
			Type:        "rename_file",
			Description: "Add .tmpl extension to template file",
			File:        issue.File,
			Action:      interfaces.FixActionRename,
			Content:     issue.File + ".tmpl",
			Automatic:   true,
		}
	}

	return nil
}

// applyCreateFix creates a new file
func (e *Engine) applyCreateFix(fix interfaces.Fix) error {
	// Check if file already exists
	if _, err := os.Stat(fix.File); err == nil {
		return fmt.Errorf("file already exists: %s", fix.File)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fix.File)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	file, err := os.Create(fix.File)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Write content
	if _, err := file.WriteString(fix.Content); err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// applyReplaceFix replaces content in a file
func (e *Engine) applyReplaceFix(fix interfaces.Fix) error {
	content, err := utils.SafeReadFile(fix.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if fix.Line > 0 && fix.Line <= len(lines) {
		lines[fix.Line-1] = fix.Content
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyInsertFix inserts content into a file
func (e *Engine) applyInsertFix(fix interfaces.Fix) error {
	content, err := utils.SafeReadFile(fix.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if fix.Line > 0 && fix.Line <= len(lines)+1 {
		// Insert at the specified line
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:fix.Line-1]...)
		newLines = append(newLines, fix.Content)
		newLines = append(newLines, lines[fix.Line-1:]...)
		lines = newLines
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyDeleteFix deletes content from a file
func (e *Engine) applyDeleteFix(fix interfaces.Fix) error {
	content, err := utils.SafeReadFile(fix.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if fix.Line > 0 && fix.Line <= len(lines) {
		// Delete the specified line
		lines = append(lines[:fix.Line-1], lines[fix.Line:]...)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyRenameFix renames a file
func (e *Engine) applyRenameFix(fix interfaces.Fix) error {
	if err := os.Rename(fix.File, fix.Content); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return nil
}

// createFileChangePreview creates a preview of file changes
func (e *Engine) createFileChangePreview(fix interfaces.Fix) (*interfaces.FileChange, error) {
	change := &interfaces.FileChange{
		File:   fix.File,
		Action: fix.Action,
	}

	switch fix.Action {
	case interfaces.FixActionCreate:
		change.LinesAfter = strings.Count(fix.Content, "\n") + 1
		change.Preview = fmt.Sprintf("Create file with %d lines", change.LinesAfter)

	case interfaces.FixActionReplace:
		content, err := utils.SafeReadFile(fix.File)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(content), "\n")
		change.LinesBefore = len(lines)
		change.LinesAfter = len(lines) // Same number of lines, just replacing content
		change.Preview = fmt.Sprintf("Replace line %d: %s", fix.Line, fix.Content)

	case interfaces.FixActionInsert:
		content, err := utils.SafeReadFile(fix.File)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(content), "\n")
		change.LinesBefore = len(lines)
		change.LinesAfter = len(lines) + 1
		change.Preview = fmt.Sprintf("Insert at line %d: %s", fix.Line, fix.Content)

	case interfaces.FixActionDelete:
		content, err := utils.SafeReadFile(fix.File)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(content), "\n")
		change.LinesBefore = len(lines)
		change.LinesAfter = len(lines) - 1
		change.Preview = fmt.Sprintf("Delete line %d", fix.Line)

	case interfaces.FixActionRename:
		change.Preview = fmt.Sprintf("Rename to: %s", fix.Content)
	}

	return change, nil
}

// generateJSONReport generates a JSON validation report
func (e *Engine) generateJSONReport(result *interfaces.ValidationResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

// generateHTMLReport generates an HTML validation report
func (e *Engine) generateHTMLReport(result *interfaces.ValidationResult) ([]byte, error) {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>Validation Report</title>
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
        <h1>Validation Report</h1>
        <p class="{{if .Valid}}valid{{else}}invalid{{end}}">
            Status: {{if .Valid}}VALID{{else}}INVALID{{end}}
        </p>
    </div>

    <div class="summary">
        <h2>Summary</h2>
        <ul>
            <li>Total Files: {{.Summary.TotalFiles}}</li>
            <li>Valid Files: {{.Summary.ValidFiles}}</li>
            <li>Errors: {{.Summary.ErrorCount}}</li>
            <li>Warnings: {{.Summary.WarningCount}}</li>
            <li>Fixable Issues: {{.Summary.FixableCount}}</li>
        </ul>
    </div>

    {{if .Issues}}
    <div class="issues">
        <h2>Issues</h2>
        {{range .Issues}}
        <div class="issue {{.Severity}}">
            <strong>{{.Type | title}}</strong> ({{.Severity | title}}): {{.Message}}
            {{if .File}}<br><em>File: {{.File}}</em>{{end}}
            {{if .Line}}<br><em>Line: {{.Line}}</em>{{end}}
            {{if .Rule}}<br><em>Rule: {{.Rule}}</em>{{end}}
            {{if .Fixable}}<br><em>✓ Fixable</em>{{end}}
        </div>
        {{end}}
    </div>
    {{end}}

    {{if .Warnings}}
    <div class="issues">
        <h2>Warnings</h2>
        {{range .Warnings}}
        <div class="issue warning">
            <strong>{{.Type | title}}</strong> ({{.Severity | title}}): {{.Message}}
            {{if .File}}<br><em>File: {{.File}}</em>{{end}}
            {{if .Line}}<br><em>Line: {{.Line}}</em>{{end}}
            {{if .Rule}}<br><em>Rule: {{.Rule}}</em>{{end}}
            {{if .Fixable}}<br><em>✓ Fixable</em>{{end}}
        </div>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, result); err != nil {
		return nil, fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return []byte(buf.String()), nil
}

// generateMarkdownReport generates a Markdown validation report
func (e *Engine) generateMarkdownReport(result *interfaces.ValidationResult) ([]byte, error) {
	var buf strings.Builder

	// Header
	buf.WriteString("# Validation Report\n\n")

	if result.Valid {
		buf.WriteString("✅ **Status: VALID**\n\n")
	} else {
		buf.WriteString("❌ **Status: INVALID**\n\n")
	}

	// Summary
	buf.WriteString("## Summary\n\n")
	buf.WriteString(fmt.Sprintf("- **Total Files:** %d\n", result.Summary.TotalFiles))
	buf.WriteString(fmt.Sprintf("- **Valid Files:** %d\n", result.Summary.ValidFiles))
	buf.WriteString(fmt.Sprintf("- **Errors:** %d\n", result.Summary.ErrorCount))
	buf.WriteString(fmt.Sprintf("- **Warnings:** %d\n", result.Summary.WarningCount))
	buf.WriteString(fmt.Sprintf("- **Fixable Issues:** %d\n\n", result.Summary.FixableCount))

	// Issues
	if len(result.Issues) > 0 {
		buf.WriteString("## Issues\n\n")
		for _, issue := range result.Issues {
			var icon string
			switch issue.Severity {
			case "warning":
				icon = "🟡"
			case "info":
				icon = "🔵"
			default:
				icon = "🔴"
			}

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

	// Warnings
	if len(result.Warnings) > 0 {
		buf.WriteString("## Warnings\n\n")
		for _, warning := range result.Warnings {
			buf.WriteString(fmt.Sprintf("🟡 **%s**: %s\n", strings.ToUpper(warning.Type[:1])+warning.Type[1:], warning.Message))
			if warning.File != "" {
				buf.WriteString(fmt.Sprintf("   - File: `%s`\n", warning.File))
			}
			if warning.Line > 0 {
				buf.WriteString(fmt.Sprintf("   - Line: %d\n", warning.Line))
			}
			if warning.Rule != "" {
				buf.WriteString(fmt.Sprintf("   - Rule: `%s`\n", warning.Rule))
			}
			if warning.Fixable {
				buf.WriteString("   - ✅ Fixable\n")
			}
			buf.WriteString("\n")
		}
	}

	return []byte(buf.String()), nil
}

// generateReadmeContent generates default README content
func (e *Engine) generateReadmeContent() string {
	return `# Project Name

## Description

Brief description of your project.

## Installation

Instructions on how to install and set up your project.

## Usage

Examples of how to use your project.

## Contributing

Guidelines for contributing to your project.

## License

Information about the project license.
`
}

// generateLicenseContent generates default LICENSE content (MIT License)
func (e *Engine) generateLicenseContent() string {
	return `MIT License

Copyright (c) 2024 Project Name

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`
}
