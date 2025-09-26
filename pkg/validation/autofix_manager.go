package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// AutoFixManager manages automatic fixes for validation issues
type AutoFixManager struct {
	fixStrategies map[string]FixStrategy
	dryRun        bool
	backupEnabled bool
}

// FixStrategy defines how to fix a specific type of validation issue
type FixStrategy struct {
	Name        string
	Description string
	Automatic   bool
	Handler     func(issue interfaces.ValidationIssue) (*interfaces.Fix, error)
}

// NewAutoFixManager creates a new auto-fix manager
func NewAutoFixManager() *AutoFixManager {
	manager := &AutoFixManager{
		fixStrategies: make(map[string]FixStrategy),
		dryRun:        false,
		backupEnabled: true,
	}
	manager.initializeFixStrategies()
	return manager
}

// SetDryRun enables or disables dry-run mode
func (afm *AutoFixManager) SetDryRun(enabled bool) {
	afm.dryRun = enabled
}

// SetBackupEnabled enables or disables backup creation
func (afm *AutoFixManager) SetBackupEnabled(enabled bool) {
	afm.backupEnabled = enabled
}

// FixIssues applies fixes for a list of validation issues
func (afm *AutoFixManager) FixIssues(projectPath string, issues []interfaces.ValidationIssue) (*interfaces.FixResult, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(projectPath); err != nil {
		return nil, fmt.Errorf("invalid project path: %w", err)
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

		result.Summary.TotalFixes++

		// Find appropriate fix strategy
		strategy, exists := afm.fixStrategies[issue.Rule]
		if !exists {
			// Try to find a generic strategy based on issue type
			strategy, exists = afm.findGenericStrategy(issue)
			if !exists {
				result.Skipped = append(result.Skipped, interfaces.Fix{
					ID:          fmt.Sprintf("skip_%s", issue.Rule),
					Type:        "skip",
					Description: fmt.Sprintf("No fix strategy available for rule: %s", issue.Rule),
					File:        issue.File,
				})
				result.Summary.SkippedFixes++
				continue
			}
		}

		// Generate fix
		fix, err := strategy.Handler(issue)
		if err != nil {
			result.Failed = append(result.Failed, interfaces.FixFailure{
				Fix: interfaces.Fix{
					ID:          fmt.Sprintf("failed_%s", issue.Rule),
					Type:        "failed",
					Description: fmt.Sprintf("Failed to generate fix for rule: %s", issue.Rule),
					File:        issue.File,
				},
				Error: err.Error(),
			})
			result.Summary.FailedFixes++
			continue
		}

		if fix == nil {
			result.Skipped = append(result.Skipped, interfaces.Fix{
				ID:          fmt.Sprintf("skip_%s", issue.Rule),
				Type:        "skip",
				Description: fmt.Sprintf("No fix generated for rule: %s", issue.Rule),
				File:        issue.File,
			})
			result.Summary.SkippedFixes++
			continue
		}

		// Apply fix
		if afm.dryRun {
			result.Applied = append(result.Applied, *fix)
			result.Summary.AppliedFixes++
		} else {
			if err := afm.applyFix(projectPath, *fix); err != nil {
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
	}

	result.Summary.FilesModified = len(modifiedFiles)

	return result, nil
}

// PreviewFixes shows what fixes would be applied without actually applying them
func (afm *AutoFixManager) PreviewFixes(projectPath string, issues []interfaces.ValidationIssue) (*interfaces.FixPreview, error) {
	// Temporarily enable dry-run mode
	originalDryRun := afm.dryRun
	afm.dryRun = true
	defer func() { afm.dryRun = originalDryRun }()

	fixResult, err := afm.FixIssues(projectPath, issues)
	if err != nil {
		return nil, fmt.Errorf("failed to preview fixes: %w", err)
	}

	preview := &interfaces.FixPreview{
		Fixes:   fixResult.Applied,
		Changes: []interfaces.FileChange{},
		Summary: fixResult.Summary,
	}

	// Generate file change previews
	for _, fix := range fixResult.Applied {
		change, err := afm.createFileChangePreview(fix)
		if err != nil {
			continue // Skip if we can't create preview
		}
		preview.Changes = append(preview.Changes, *change)
	}

	return preview, nil
}

// GetFixableIssues returns only the issues that can be automatically fixed
func (afm *AutoFixManager) GetFixableIssues(issues []interfaces.ValidationIssue) []interfaces.ValidationIssue {
	var fixableIssues []interfaces.ValidationIssue

	for _, issue := range issues {
		if issue.Fixable {
			// Check if we have a fix strategy for this issue
			if _, exists := afm.fixStrategies[issue.Rule]; exists {
				fixableIssues = append(fixableIssues, issue)
			} else if _, exists := afm.findGenericStrategy(issue); exists {
				fixableIssues = append(fixableIssues, issue)
			}
		}
	}

	return fixableIssues
}

// applyFix applies a single fix
func (afm *AutoFixManager) applyFix(projectPath string, fix interfaces.Fix) error {
	// Create backup if enabled
	if afm.backupEnabled && fix.Action != interfaces.FixActionCreate {
		if err := afm.createBackup(fix.File); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	switch fix.Action {
	case interfaces.FixActionCreate:
		return afm.applyCreateFix(fix)
	case interfaces.FixActionReplace:
		return afm.applyReplaceFix(fix)
	case interfaces.FixActionInsert:
		return afm.applyInsertFix(fix)
	case interfaces.FixActionDelete:
		return afm.applyDeleteFix(fix)
	case interfaces.FixActionRename:
		return afm.applyRenameFix(fix)
	case interfaces.FixActionMove:
		return afm.applyMoveFix(fix)
	default:
		return fmt.Errorf("unsupported fix action: %s", fix.Action)
	}
}

// createBackup creates a backup of the file before modification
func (afm *AutoFixManager) createBackup(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupPath := filePath + ".backup"
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read original file: %w", err)
	}

	if err := os.WriteFile(backupPath, content, 0o600); err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}

	return nil
}

// Fix application methods

// applyCreateFix creates a new file
func (afm *AutoFixManager) applyCreateFix(fix interfaces.Fix) error {
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
	if err := os.WriteFile(fix.File, []byte(fix.Content), 0o600); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	return nil
}

// applyReplaceFix replaces content in a file
func (afm *AutoFixManager) applyReplaceFix(fix interfaces.Fix) error {
	content, err := utils.SafeReadFile(fix.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if fix.Line > 0 && fix.Line <= len(lines) {
		lines[fix.Line-1] = fix.Content
	} else {
		return fmt.Errorf("invalid line number: %d", fix.Line)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyInsertFix inserts content into a file
func (afm *AutoFixManager) applyInsertFix(fix interfaces.Fix) error {
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
	} else {
		return fmt.Errorf("invalid line number: %d", fix.Line)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyDeleteFix deletes content from a file
func (afm *AutoFixManager) applyDeleteFix(fix interfaces.Fix) error {
	content, err := utils.SafeReadFile(fix.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	if fix.Line > 0 && fix.Line <= len(lines) {
		// Delete the specified line
		lines = append(lines[:fix.Line-1], lines[fix.Line:]...)
	} else {
		return fmt.Errorf("invalid line number: %d", fix.Line)
	}

	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(fix.File, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// applyRenameFix renames a file
func (afm *AutoFixManager) applyRenameFix(fix interfaces.Fix) error {
	if err := os.Rename(fix.File, fix.Content); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	return nil
}

// applyMoveFix moves a file to a new location
func (afm *AutoFixManager) applyMoveFix(fix interfaces.Fix) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(fix.Content)
	if err := os.MkdirAll(destDir, 0o750); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err := os.Rename(fix.File, fix.Content); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}

// createFileChangePreview creates a preview of file changes
func (afm *AutoFixManager) createFileChangePreview(fix interfaces.Fix) (*interfaces.FileChange, error) {
	change := &interfaces.FileChange{
		File:   fix.File,
		Action: fix.Action,
	}

	switch fix.Action {
	case interfaces.FixActionCreate:
		change.LinesAfter = strings.Count(fix.Content, "\n") + 1
		change.Preview = fmt.Sprintf("Create file with %d lines", change.LinesAfter)

	case interfaces.FixActionReplace:
		if _, err := os.Stat(fix.File); err == nil {
			content, err := utils.SafeReadFile(fix.File)
			if err != nil {
				return nil, err
			}
			lines := strings.Split(string(content), "\n")
			change.LinesBefore = len(lines)
			change.LinesAfter = len(lines) // Same number of lines, just replacing content
			change.Preview = fmt.Sprintf("Replace line %d: %s", fix.Line, fix.Content)
		}

	case interfaces.FixActionInsert:
		if _, err := os.Stat(fix.File); err == nil {
			content, err := utils.SafeReadFile(fix.File)
			if err != nil {
				return nil, err
			}
			lines := strings.Split(string(content), "\n")
			change.LinesBefore = len(lines)
			change.LinesAfter = len(lines) + 1
			change.Preview = fmt.Sprintf("Insert at line %d: %s", fix.Line, fix.Content)
		}

	case interfaces.FixActionDelete:
		if _, err := os.Stat(fix.File); err == nil {
			content, err := utils.SafeReadFile(fix.File)
			if err != nil {
				return nil, err
			}
			lines := strings.Split(string(content), "\n")
			change.LinesBefore = len(lines)
			change.LinesAfter = len(lines) - 1
			change.Preview = fmt.Sprintf("Delete line %d", fix.Line)
		}

	case interfaces.FixActionRename:
		change.Preview = fmt.Sprintf("Rename to: %s", fix.Content)

	case interfaces.FixActionMove:
		change.Preview = fmt.Sprintf("Move to: %s", fix.Content)
	}

	return change, nil
}

// findGenericStrategy finds a generic fix strategy based on issue type
func (afm *AutoFixManager) findGenericStrategy(issue interfaces.ValidationIssue) (FixStrategy, bool) {
	// Try to match by issue type or message patterns
	if strings.Contains(issue.Message, "missing") && strings.Contains(issue.Message, "file") {
		return afm.fixStrategies["generic.create_missing_file"], true
	}

	if strings.Contains(issue.Message, "space") && strings.Contains(issue.Message, "name") {
		return afm.fixStrategies["generic.fix_naming"], true
	}

	return FixStrategy{}, false
}

// initializeFixStrategies initializes the available fix strategies
func (afm *AutoFixManager) initializeFixStrategies() {
	// README file creation
	afm.fixStrategies["structure.readme.required"] = FixStrategy{
		Name:        "Create README",
		Description: "Creates a basic README.md file",
		Automatic:   true,
		Handler:     afm.createReadmeFix,
	}

	// LICENSE file creation
	afm.fixStrategies["structure.license.required"] = FixStrategy{
		Name:        "Create LICENSE",
		Description: "Creates a basic LICENSE file",
		Automatic:   true,
		Handler:     afm.createLicenseFix,
	}

	// Naming convention fixes
	afm.fixStrategies["quality.naming.conventions"] = FixStrategy{
		Name:        "Fix Naming Conventions",
		Description: "Fixes file and directory naming issues",
		Automatic:   true,
		Handler:     afm.fixNamingConventions,
	}

	// Template file extension fix
	afm.fixStrategies["template.file.extension"] = FixStrategy{
		Name:        "Add Template Extension",
		Description: "Adds .tmpl extension to template files",
		Automatic:   true,
		Handler:     afm.fixTemplateExtension,
	}

	// Gitignore creation
	afm.fixStrategies["structure.gitignore.recommended"] = FixStrategy{
		Name:        "Create Gitignore",
		Description: "Creates a basic .gitignore file",
		Automatic:   true,
		Handler:     afm.createGitignoreFix,
	}

	// Permission fixes
	afm.fixStrategies["security.permissions"] = FixStrategy{
		Name:        "Fix File Permissions",
		Description: "Fixes overly permissive file permissions",
		Automatic:   true,
		Handler:     afm.fixFilePermissions,
	}

	// Generic strategies
	afm.fixStrategies["generic.create_missing_file"] = FixStrategy{
		Name:        "Create Missing File",
		Description: "Creates a missing file with basic content",
		Automatic:   false,
		Handler:     afm.createMissingFile,
	}

	afm.fixStrategies["generic.fix_naming"] = FixStrategy{
		Name:        "Fix Naming",
		Description: "Fixes naming convention issues",
		Automatic:   true,
		Handler:     afm.fixNamingConventions,
	}
}

// Fix strategy handlers

// createReadmeFix creates a README.md file
func (afm *AutoFixManager) createReadmeFix(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	readmeContent := `# Project Name

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

	return &interfaces.Fix{
		ID:          fmt.Sprintf("create_readme_%s", issue.File),
		Type:        "create_file",
		Description: "Create README.md file",
		File:        filepath.Join(filepath.Dir(issue.File), "README.md"),
		Action:      interfaces.FixActionCreate,
		Content:     readmeContent,
		Automatic:   true,
	}, nil
}

// createLicenseFix creates a LICENSE file
func (afm *AutoFixManager) createLicenseFix(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	licenseContent := `MIT License

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

	return &interfaces.Fix{
		ID:          fmt.Sprintf("create_license_%s", issue.File),
		Type:        "create_file",
		Description: "Create LICENSE file",
		File:        filepath.Join(filepath.Dir(issue.File), "LICENSE"),
		Action:      interfaces.FixActionCreate,
		Content:     licenseContent,
		Automatic:   true,
	}, nil
}

// fixNamingConventions fixes naming convention issues
func (afm *AutoFixManager) fixNamingConventions(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	if strings.Contains(issue.Message, "space") {
		// Replace spaces with underscores
		newName := strings.ReplaceAll(filepath.Base(issue.File), " ", "_")
		newPath := filepath.Join(filepath.Dir(issue.File), newName)

		return &interfaces.Fix{
			ID:          fmt.Sprintf("fix_naming_%s", issue.File),
			Type:        "rename_file",
			Description: fmt.Sprintf("Rename file to follow naming conventions: %s", newName),
			File:        issue.File,
			Action:      interfaces.FixActionRename,
			Content:     newPath,
			Automatic:   true,
		}, nil
	}

	return nil, fmt.Errorf("unsupported naming convention issue")
}

// fixTemplateExtension adds .tmpl extension to template files
func (afm *AutoFixManager) fixTemplateExtension(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	return &interfaces.Fix{
		ID:          fmt.Sprintf("fix_template_ext_%s", issue.File),
		Type:        "rename_file",
		Description: "Add .tmpl extension to template file",
		File:        issue.File,
		Action:      interfaces.FixActionRename,
		Content:     issue.File + ".tmpl",
		Automatic:   true,
	}, nil
}

// createGitignoreFix creates a .gitignore file
func (afm *AutoFixManager) createGitignoreFix(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	gitignoreContent := `# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
*.exe
*.dll
*.so
*.dylib

# Logs
*.log
logs/

# Environment variables
.env
.env.local
.env.*.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Temporary files
*.tmp
*.temp
`

	return &interfaces.Fix{
		ID:          fmt.Sprintf("create_gitignore_%s", issue.File),
		Type:        "create_file",
		Description: "Create .gitignore file",
		File:        filepath.Join(filepath.Dir(issue.File), ".gitignore"),
		Action:      interfaces.FixActionCreate,
		Content:     gitignoreContent,
		Automatic:   true,
	}, nil
}

// fixFilePermissions fixes file permission issues
func (afm *AutoFixManager) fixFilePermissions(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	// This would typically use chmod, but we'll create a fix that documents the change
	return &interfaces.Fix{
		ID:          fmt.Sprintf("fix_permissions_%s", issue.File),
		Type:        "fix_permissions",
		Description: "Fix file permissions to be more secure",
		File:        issue.File,
		Action:      "chmod", // Custom action for permission changes
		Content:     "644",   // New permission mode
		Automatic:   true,
	}, nil
}

// createMissingFile creates a missing file with basic content
func (afm *AutoFixManager) createMissingFile(issue interfaces.ValidationIssue) (*interfaces.Fix, error) {
	// Determine file type and create appropriate content
	ext := strings.ToLower(filepath.Ext(issue.File))
	var content string

	switch ext {
	case ".md":
		content = fmt.Sprintf("# %s\n\nThis file was automatically generated.\n", filepath.Base(issue.File))
	case ".txt":
		content = "This file was automatically generated.\n"
	case ".json":
		content = "{}\n"
	case ".yaml", ".yml":
		content = "# This file was automatically generated\n"
	default:
		content = "# This file was automatically generated\n"
	}

	return &interfaces.Fix{
		ID:          fmt.Sprintf("create_missing_%s", issue.File),
		Type:        "create_file",
		Description: fmt.Sprintf("Create missing file: %s", filepath.Base(issue.File)),
		File:        issue.File,
		Action:      interfaces.FixActionCreate,
		Content:     content,
		Automatic:   false, // Require user confirmation for generic file creation
	}, nil
}
