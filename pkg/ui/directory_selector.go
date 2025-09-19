// Package ui provides interactive user interface components for directory selection and validation.
//
// This file implements the DirectorySelector which handles interactive selection
// of output directories with path validation, existence checking, and conflict handling.
package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DirectorySelector handles interactive selection and validation of output directories
type DirectorySelector struct {
	ui        interfaces.InteractiveUIInterface
	validator *DirectoryValidator
	logger    interfaces.Logger
}

// DirectoryValidator provides validation for directory paths and operations
type DirectoryValidator struct {
	// Add any validation configuration here
}

// DirectorySelectionResult contains the result of directory selection
type DirectorySelectionResult struct {
	Path               string
	Exists             bool
	RequiresCreation   bool
	ConflictResolution string // "overwrite", "backup", "merge", "cancel"
	BackupPath         string
	Cancelled          bool
}

// NewDirectorySelector creates a new directory selector
func NewDirectorySelector(ui interfaces.InteractiveUIInterface, logger interfaces.Logger) *DirectorySelector {
	validator := &DirectoryValidator{}

	return &DirectorySelector{
		ui:        ui,
		validator: validator,
		logger:    logger,
	}
}

// SelectOutputDirectory interactively selects and validates an output directory
func (ds *DirectorySelector) SelectOutputDirectory(ctx context.Context, defaultPath string) (*DirectorySelectionResult, error) {
	if defaultPath == "" {
		defaultPath = "output/generated"
	}

	// Step 1: Collect directory path
	path, err := ds.collectDirectoryPath(ctx, defaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to collect directory path: %w", err)
	}

	if path == "" {
		return &DirectorySelectionResult{Cancelled: true}, nil
	}

	// Step 2: Validate and check path
	result, err := ds.validateAndCheckPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to validate path: %w", err)
	}

	// Step 3: Handle conflicts if directory exists and is not empty
	if result.Exists {
		// Check if directory is empty
		isEmpty, err := ds.isDirectoryEmpty(result.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to check if directory is empty: %w", err)
		}

		// Only handle conflicts for non-empty directories
		if !isEmpty {
			conflictResult, err := ds.handleDirectoryConflict(ctx, result)
			if err != nil {
				return nil, fmt.Errorf("failed to handle directory conflict: %w", err)
			}
			result = conflictResult
		}
	}

	return result, nil
}

// collectDirectoryPath collects the output directory path from user
func (ds *DirectorySelector) collectDirectoryPath(ctx context.Context, defaultPath string) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:       "Output Directory",
		Description:  "Enter the path where your project should be generated",
		DefaultValue: defaultPath,
		Required:     true,
		Validator:    ds.validator.ValidateDirectoryPath,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Output Directory Guidelines:
• Specify where your project files will be created
• Can be relative (e.g., ./my-project) or absolute (e.g., /home/user/projects/my-project)
• Directory will be created if it doesn't exist
• Use forward slashes (/) for cross-platform compatibility
• Examples: 
  - output/my-app (relative to current directory)
  - ./projects/new-project (relative to current directory)
  - /home/user/workspace/my-project (absolute path)
  - ~/projects/my-app (home directory relative)`,
		MaxLength: 500,
		MinLength: 1,
	}

	result, err := ds.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", nil
	}

	// Clean and normalize the path
	path := strings.TrimSpace(result.Value)
	if path == "" {
		path = defaultPath
	}

	// Expand home directory if needed
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to convert to absolute path: %w", err)
	}

	return absPath, nil
}

// validateAndCheckPath validates the path and checks if it exists
func (ds *DirectorySelector) validateAndCheckPath(ctx context.Context, path string) (*DirectorySelectionResult, error) {
	result := &DirectorySelectionResult{
		Path: path,
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path doesn't exist - this is fine, we'll create it
			result.Exists = false
			result.RequiresCreation = true

			// Check if parent directory exists and is writable
			parentDir := filepath.Dir(path)
			if err := ds.validator.ValidateParentDirectory(parentDir); err != nil {
				return nil, fmt.Errorf("parent directory validation failed: %w", err)
			}

			return result, nil
		}
		return nil, fmt.Errorf("failed to check path: %w", err)
	}

	// Path exists - check if it's a directory
	if !info.IsDir() {
		return nil, fmt.Errorf("path exists but is not a directory: %s", path)
	}

	result.Exists = true
	result.RequiresCreation = false

	// Check if directory is empty
	isEmpty, err := ds.isDirectoryEmpty(path)
	if err != nil {
		return nil, fmt.Errorf("failed to check if directory is empty: %w", err)
	}

	if isEmpty {
		// Empty directory is fine to use
		return result, nil
	}

	// Directory exists and is not empty - will need conflict resolution
	return result, nil
}

// handleDirectoryConflict handles conflicts when the target directory exists and is not empty
func (ds *DirectorySelector) handleDirectoryConflict(ctx context.Context, result *DirectorySelectionResult) (*DirectorySelectionResult, error) {
	// First, show what's in the directory
	if err := ds.showDirectoryContents(ctx, result.Path); err != nil {
		ds.logger.WarnWithFields("Failed to show directory contents", map[string]interface{}{
			"path":  result.Path,
			"error": err.Error(),
		})
	}

	// Present conflict resolution options
	options := []interfaces.MenuOption{
		{
			Label:       "Overwrite",
			Description: "Replace all existing files (creates backup first)",
			Value:       "overwrite",
			Icon:        "⚠️",
		},
		{
			Label:       "Merge",
			Description: "Keep existing files, add new ones (may cause conflicts)",
			Value:       "merge",
			Icon:        "🔄",
		},
		{
			Label:       "Choose Different Directory",
			Description: "Select a different output directory",
			Value:       "different",
			Icon:        "📁",
		},
		{
			Label:       "Cancel",
			Description: "Cancel the generation process",
			Value:       "cancel",
			Icon:        "❌",
		},
	}

	config := interfaces.MenuConfig{
		Title:       "Directory Conflict",
		Description: fmt.Sprintf("The directory '%s' already exists and contains files.", result.Path),
		Options:     options,
		DefaultItem: 2, // Default to "Choose Different Directory"
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Directory Conflict Resolution:
• Overwrite: All existing files will be replaced. A backup will be created automatically.
• Merge: Existing files will be kept, new files will be added. This may cause conflicts if files have the same names.
• Choose Different Directory: Select a different location for your project.
• Cancel: Stop the generation process.

Recommendation: Choose a different directory to avoid potential issues.`,
	}

	menuResult, err := ds.ui.ShowMenu(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to show conflict resolution menu: %w", err)
	}

	if menuResult.Cancelled || menuResult.Action != "select" {
		result.Cancelled = true
		return result, nil
	}

	resolution := menuResult.SelectedValue.(string)
	result.ConflictResolution = resolution

	switch resolution {
	case "overwrite":
		// Confirm overwrite with backup
		confirmed, err := ds.confirmOverwriteWithBackup(ctx, result.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !confirmed {
			result.Cancelled = true
			return result, nil
		}

		// Generate backup path
		result.BackupPath = ds.generateBackupPath(result.Path)

	case "merge":
		// Confirm merge operation
		confirmed, err := ds.confirmMergeOperation(ctx, result.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm merge: %w", err)
		}
		if !confirmed {
			result.Cancelled = true
			return result, nil
		}

	case "different":
		// Recursively call directory selection
		newResult, err := ds.SelectOutputDirectory(ctx, filepath.Dir(result.Path))
		if err != nil {
			return nil, fmt.Errorf("failed to select different directory: %w", err)
		}
		return newResult, nil

	case "cancel":
		result.Cancelled = true
		return result, nil

	default:
		return nil, fmt.Errorf("unknown conflict resolution: %s", resolution)
	}

	return result, nil
}

// showDirectoryContents displays the contents of the directory to help user make decision
func (ds *DirectorySelector) showDirectoryContents(ctx context.Context, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	if len(entries) == 0 {
		return nil // Directory is empty
	}

	// Limit to first 10 entries to avoid overwhelming output
	maxEntries := 10
	rows := make([][]string, 0, len(entries))

	for i, entry := range entries {
		if i >= maxEntries {
			rows = append(rows, []string{"...", fmt.Sprintf("and %d more items", len(entries)-maxEntries), ""})
			break
		}

		entryType := "File"
		size := ""
		if entry.IsDir() {
			entryType = "Directory"
			// Count items in subdirectory
			subPath := filepath.Join(path, entry.Name())
			if subEntries, err := os.ReadDir(subPath); err == nil {
				size = fmt.Sprintf("%d items", len(subEntries))
			}
		} else {
			// Get file size
			if info, err := entry.Info(); err == nil {
				size = formatFileSize(info.Size())
			}
		}

		rows = append(rows, []string{entry.Name(), entryType, size})
	}

	tableConfig := interfaces.TableConfig{
		Title:   fmt.Sprintf("Directory Contents: %s", path),
		Headers: []string{"Name", "Type", "Size"},
		Rows:    rows,
	}

	return ds.ui.ShowTable(ctx, tableConfig)
}

// confirmOverwriteWithBackup confirms overwrite operation with backup
func (ds *DirectorySelector) confirmOverwriteWithBackup(ctx context.Context, path string) (bool, error) {
	backupPath := ds.generateBackupPath(path)

	config := interfaces.ConfirmConfig{
		Prompt: "Confirm Overwrite with Backup",
		Description: fmt.Sprintf(`This will:
1. Create a backup of existing files at: %s
2. Remove all existing files from: %s
3. Generate new project files

This operation cannot be undone (except by restoring from backup).`, backupPath, path),
		DefaultValue: false, // Default to No for safety
		YesLabel:     "Overwrite with Backup",
		NoLabel:      "Cancel",
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Overwrite with Backup:
• A complete backup will be created before any files are modified
• The backup will include all existing files and directories
• You can restore from the backup if needed
• This is the safest option if you need to replace existing files`,
	}

	result, err := ds.ui.PromptConfirm(ctx, config)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

// confirmMergeOperation confirms merge operation
func (ds *DirectorySelector) confirmMergeOperation(ctx context.Context, path string) (bool, error) {
	config := interfaces.ConfirmConfig{
		Prompt: "Confirm Merge Operation",
		Description: fmt.Sprintf(`This will:
1. Keep all existing files in: %s
2. Add new project files alongside existing ones
3. Overwrite files with the same names (no backup)

Existing files with different names will be preserved.`, path),
		DefaultValue: false, // Default to No for safety
		YesLabel:     "Merge Files",
		NoLabel:      "Cancel",
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Merge Operation:
• Existing files will be preserved unless they have the same name as generated files
• Files with the same names will be overwritten without backup
• This may cause conflicts or unexpected behavior
• Consider using a different directory if you're unsure`,
	}

	result, err := ds.ui.PromptConfirm(ctx, config)
	if err != nil {
		return false, err
	}

	return result.Confirmed && result.Action == "confirm", nil
}

// generateBackupPath generates a unique backup path
func (ds *DirectorySelector) generateBackupPath(originalPath string) string {
	timestamp := fmt.Sprintf("%d", os.Getpid()) // Use PID for uniqueness in this session
	backupName := fmt.Sprintf("%s.backup.%s", filepath.Base(originalPath), timestamp)
	return filepath.Join(filepath.Dir(originalPath), backupName)
}

// isDirectoryEmpty checks if a directory is empty
func (ds *DirectorySelector) isDirectoryEmpty(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// formatFileSize formats file size in human-readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// Validation methods for DirectoryValidator

// ValidateDirectoryPath validates the directory path format and accessibility
func (v *DirectoryValidator) ValidateDirectoryPath(path string) error {
	if path == "" {
		return interfaces.NewValidationError("path", path, "Directory path is required", "required").
			WithSuggestions("Enter a valid directory path")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)
	if cleanPath != path && cleanPath != "." {
		return interfaces.NewValidationError("path", path, "Path contains invalid characters or sequences", "invalid_format").
			WithSuggestions(
				"Use forward slashes (/) for directory separators",
				"Avoid '..' or '.' in the middle of paths",
				fmt.Sprintf("Suggested clean path: %s", cleanPath),
			)
	}

	// Check for invalid characters (Windows-specific, but good to check everywhere)
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return interfaces.NewValidationError("path", path, fmt.Sprintf("Path contains invalid character: %s", char), "invalid_character").
				WithSuggestions(
					"Remove invalid characters from the path",
					"Use only letters, numbers, hyphens, underscores, and forward slashes",
				)
		}
	}

	// Check path length (reasonable limit)
	if len(path) > 500 {
		return interfaces.NewValidationError("path", path, "Path is too long (maximum 500 characters)", "max_length").
			WithSuggestions("Use a shorter path")
	}

	// Check for reserved names on Windows
	pathParts := strings.Split(filepath.Clean(path), string(filepath.Separator))
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	for _, part := range pathParts {
		upperPart := strings.ToUpper(part)
		for _, reserved := range reservedNames {
			if upperPart == reserved {
				return interfaces.NewValidationError("path", path, fmt.Sprintf("Path contains reserved name: %s", part), "reserved_name").
					WithSuggestions("Use a different directory name that doesn't conflict with system reserved words")
			}
		}
	}

	return nil
}

// ValidateParentDirectory validates that the parent directory exists and is writable
func (v *DirectoryValidator) ValidateParentDirectory(parentPath string) error {
	// Check if parent directory exists
	info, err := os.Stat(parentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return interfaces.NewValidationError("parent", parentPath, "Parent directory does not exist", "not_exists").
				WithSuggestions(
					"Create the parent directory first",
					"Use a path where the parent directory already exists",
					fmt.Sprintf("Parent directory: %s", parentPath),
				)
		}
		return interfaces.NewValidationError("parent", parentPath, fmt.Sprintf("Cannot access parent directory: %v", err), "access_error").
			WithSuggestions("Check directory permissions and path validity")
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return interfaces.NewValidationError("parent", parentPath, "Parent path exists but is not a directory", "not_directory").
			WithSuggestions("Use a different path where the parent is a directory")
	}

	// Check if we can write to the parent directory
	testFile := filepath.Join(parentPath, ".write_test_temp")
	// Validate the test file path to prevent directory traversal
	testFile = filepath.Clean(testFile)
	if !strings.HasPrefix(testFile, filepath.Clean(parentPath)) {
		return interfaces.NewValidationError("parent", parentPath, "Invalid test file path", "invalid_path")
	}
	file, err := os.Create(testFile)
	if err != nil {
		return interfaces.NewValidationError("parent", parentPath, "Cannot write to parent directory (permission denied)", "permission_denied").
			WithSuggestions(
				"Check directory permissions",
				"Use a directory where you have write access",
				"Run with appropriate permissions",
			)
	}
	_ = file.Close()
	_ = os.Remove(testFile) // Clean up test file

	return nil
}

// CreateDirectory creates the directory and any necessary parent directories
func (ds *DirectorySelector) CreateDirectory(path string) error {
	if err := os.MkdirAll(path, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	ds.logger.InfoWithFields("Created directory", map[string]interface{}{
		"path": path,
	})

	return nil
}

// CreateBackup creates a backup of the existing directory
func (ds *DirectorySelector) CreateBackup(sourcePath, backupPath string) error {
	// This is a simplified backup - in a real implementation, you might want to use
	// more sophisticated backup methods or external tools
	if err := ds.copyDirectory(sourcePath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	ds.logger.InfoWithFields("Created backup", map[string]interface{}{
		"source": sourcePath,
		"backup": backupPath,
	})

	return nil
}

// copyDirectory recursively copies a directory
func (ds *DirectorySelector) copyDirectory(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Read source directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := ds.copyDirectory(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy subdirectory %s: %w", entry.Name(), err)
			}
		} else {
			// Copy file
			if err := ds.copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// copyFile copies a single file
func (ds *DirectorySelector) copyFile(src, dst string) error {
	// Validate and clean paths to prevent directory traversal
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// Ensure paths are absolute to prevent traversal attacks
	if !filepath.IsAbs(src) {
		return fmt.Errorf("source path must be absolute: %s", src)
	}
	if !filepath.IsAbs(dst) {
		return fmt.Errorf("destination path must be absolute: %s", dst)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = dstFile.Close() }()

	// Copy file contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
