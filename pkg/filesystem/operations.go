package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// FileSystemOperations handles safe file and directory operations
type FileSystemOperations struct {
	dryRun bool
}

// NewFileSystemOperations creates a new filesystem operations handler
func NewFileSystemOperations() *FileSystemOperations {
	return &FileSystemOperations{
		dryRun: false,
	}
}

// NewDryRunFileSystemOperations creates a new filesystem operations handler in dry-run mode
func NewDryRunFileSystemOperations() *FileSystemOperations {
	return &FileSystemOperations{
		dryRun: true,
	}
}

// CreateDirectory creates a directory with proper permissions and validation
func (fso *FileSystemOperations) CreateDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Validate that the path doesn't contain dangerous elements before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	// Clean the path after validation
	cleanPath := filepath.Clean(path)

	if fso.dryRun {
		return nil
	}

	// Create directory with secure permissions (0750)
	if err := os.MkdirAll(cleanPath, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", cleanPath, err)
	}

	return nil
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func (fso *FileSystemOperations) EnsureDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Validate path before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	cleanPath := filepath.Clean(path)

	if fso.dryRun {
		return nil
	}

	// Check if directory already exists
	if fso.FileExists(cleanPath) {
		// Verify it's actually a directory
		info, err := os.Stat(cleanPath)
		if err != nil {
			return fmt.Errorf("failed to stat existing path %s: %w", cleanPath, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path %s exists but is not a directory", cleanPath)
		}
		return nil
	}

	// Create directory with proper permissions
	return fso.CreateDirectory(cleanPath)
}

// WriteFile writes content to a file with specified permissions and validation
func (fso *FileSystemOperations) WriteFile(path string, content []byte, perm os.FileMode) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if content == nil {
		return fmt.Errorf("file content cannot be nil")
	}

	// Validate that the path doesn't contain dangerous elements before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	// Clean the path after validation
	cleanPath := filepath.Clean(path)

	// Validate permissions (must be between 0000 and 0777)
	if perm > 0777 {
		return fmt.Errorf("invalid file permissions: %o exceeds maximum allowed permissions", perm)
	}

	if fso.dryRun {
		return nil
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(cleanPath)
	if err := fso.EnsureDirectory(dir); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", cleanPath, err)
	}

	// Write the file with specified permissions
	if err := os.WriteFile(cleanPath, content, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", cleanPath, err)
	}

	return nil
}

// FileExists checks if a file or directory exists at the given path
func (fso *FileSystemOperations) FileExists(path string) bool {
	if path == "" {
		return false
	}

	cleanPath := filepath.Clean(path)
	_, err := os.Stat(cleanPath)
	return err == nil
}

// CopyFile copies a single file from source to destination with validation
func (fso *FileSystemOperations) CopyFile(srcPath, destPath string, perm os.FileMode) error {
	if srcPath == "" {
		return fmt.Errorf("source path cannot be empty")
	}

	if destPath == "" {
		return fmt.Errorf("destination path cannot be empty")
	}

	// Validate paths before cleaning
	if strings.Contains(srcPath, "..") || strings.Contains(destPath, "..") {
		return fmt.Errorf("invalid path: path traversal detected")
	}

	cleanSrcPath := filepath.Clean(srcPath)
	cleanDestPath := filepath.Clean(destPath)

	// Validate permissions
	if perm > 0777 {
		return fmt.Errorf("invalid file permissions: %o exceeds maximum allowed permissions", perm)
	}

	if fso.dryRun {
		return nil
	}

	// Check if source file exists
	if !fso.FileExists(cleanSrcPath) {
		return fmt.Errorf("source file does not exist: %s", cleanSrcPath)
	}

	// Open source file with path validation
	srcFile, err := utils.SafeOpen(cleanSrcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(cleanDestPath)
	if err := fso.EnsureDirectory(destDir); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file with secure permissions
	destFile, err := utils.SafeCreate(cleanDestPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	// Copy file content
	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Set file permissions
	return os.Chmod(cleanDestPath, perm)
}

// CopyDirectory recursively copies a directory from source to destination
func (fso *FileSystemOperations) CopyDirectory(srcDir, destDir string) error {
	if srcDir == "" {
		return fmt.Errorf("source directory cannot be empty")
	}

	if destDir == "" {
		return fmt.Errorf("destination directory cannot be empty")
	}

	// Validate paths before cleaning
	if strings.Contains(srcDir, "..") || strings.Contains(destDir, "..") {
		return fmt.Errorf("invalid path: path traversal detected")
	}

	cleanSrcDir := filepath.Clean(srcDir)
	cleanDestDir := filepath.Clean(destDir)

	// Check if source directory exists
	if !fso.FileExists(cleanSrcDir) {
		return fmt.Errorf("source directory does not exist: %s", cleanSrcDir)
	}

	if fso.dryRun {
		return nil
	}

	// Ensure destination directory exists
	if err := fso.EnsureDirectory(cleanDestDir); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", cleanDestDir, err)
	}

	// Walk through source directory and copy all files
	return filepath.Walk(cleanSrcDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking source directory: %w", err)
		}

		// Calculate relative path from source directory
		relPath, err := filepath.Rel(cleanSrcDir, srcPath)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}

		// Calculate destination path
		destPath := filepath.Join(cleanDestDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return fso.CreateDirectory(destPath)
		}

		// Copy file
		return fso.CopyFile(srcPath, destPath, info.Mode())
	})
}

// CreateSymlink creates a symbolic link with validation
func (fso *FileSystemOperations) CreateSymlink(target, link string) error {
	if target == "" {
		return fmt.Errorf("symlink target cannot be empty")
	}

	if link == "" {
		return fmt.Errorf("symlink path cannot be empty")
	}

	// Validate paths before cleaning
	if strings.Contains(link, "..") {
		return fmt.Errorf("invalid path: path traversal detected in link path %s", link)
	}

	// Clean paths after validation
	cleanTarget := filepath.Clean(target)
	cleanLink := filepath.Clean(link)

	if fso.dryRun {
		return nil
	}

	// Ensure parent directory of link exists
	linkDir := filepath.Dir(cleanLink)
	if err := fso.EnsureDirectory(linkDir); err != nil {
		return fmt.Errorf("failed to create parent directory for symlink %s: %w", cleanLink, err)
	}

	// Create symbolic link
	if err := os.Symlink(cleanTarget, cleanLink); err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", cleanTarget, cleanLink, err)
	}

	return nil
}

// ValidateFileContent validates that a file exists and is readable
func (fso *FileSystemOperations) ValidateFileContent(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	cleanPath := filepath.Clean(path)

	if !fso.FileExists(cleanPath) {
		return fmt.Errorf("file does not exist: %s", cleanPath)
	}

	// Try to read the file to ensure it's accessible
	if _, err := utils.SafeReadFile(cleanPath); err != nil {
		return fmt.Errorf("failed to read file %s: %w", cleanPath, err)
	}

	return nil
}

// ValidateDirectoryStructure validates that required directories exist
func (fso *FileSystemOperations) ValidateDirectoryStructure(basePath string, requiredDirs []string) error {
	if basePath == "" {
		return fmt.Errorf("base path cannot be empty")
	}

	if len(requiredDirs) == 0 {
		return fmt.Errorf("required directories list cannot be empty")
	}

	cleanBasePath := filepath.Clean(basePath)

	if !fso.FileExists(cleanBasePath) {
		return fmt.Errorf("base directory does not exist: %s", cleanBasePath)
	}

	for _, dir := range requiredDirs {
		if dir == "" {
			continue // Skip empty directory names
		}

		dirPath := filepath.Join(cleanBasePath, dir)
		if !fso.FileExists(dirPath) {
			return fmt.Errorf("required directory missing: %s", dir)
		}

		// Verify it's actually a directory
		info, err := os.Stat(dirPath)
		if err != nil {
			return fmt.Errorf("failed to stat directory %s: %w", dir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path %s exists but is not a directory", dir)
		}
	}

	return nil
}

// ValidateFileStructure validates that required files exist
func (fso *FileSystemOperations) ValidateFileStructure(basePath string, requiredFiles []string) error {
	if basePath == "" {
		return fmt.Errorf("base path cannot be empty")
	}

	if len(requiredFiles) == 0 {
		return fmt.Errorf("required files list cannot be empty")
	}

	cleanBasePath := filepath.Clean(basePath)

	if !fso.FileExists(cleanBasePath) {
		return fmt.Errorf("base directory does not exist: %s", cleanBasePath)
	}

	for _, file := range requiredFiles {
		if file == "" {
			continue // Skip empty file names
		}

		filePath := filepath.Join(cleanBasePath, file)
		if !fso.FileExists(filePath) {
			return fmt.Errorf("required file missing: %s", file)
		}

		// Verify it's actually a file (not a directory)
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", file, err)
		}
		if info.IsDir() {
			return fmt.Errorf("path %s exists but is a directory, expected file", file)
		}
	}

	return nil
}

// CreateProjectRoot creates the root project directory with validation
func (fso *FileSystemOperations) CreateProjectRoot(config *models.ProjectConfig, outputPath string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("project config cannot be nil")
	}

	if outputPath == "" {
		return "", fmt.Errorf("output path cannot be empty")
	}

	// Validate required config fields
	if config.Name == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}

	// Ensure the output path is absolute
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("unable to resolve output path '%s': %w", outputPath, err)
	}

	// Create the root project directory
	projectPath := filepath.Join(absOutputPath, config.Name)
	if err := fso.EnsureDirectory(projectPath); err != nil {
		return "", fmt.Errorf("failed to create project root directory '%s': %w", projectPath, err)
	}

	return projectPath, nil
}

// ValidateProjectRoot validates that a project root directory exists and is valid
func (fso *FileSystemOperations) ValidateProjectRoot(projectPath string, config *models.ProjectConfig) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	cleanProjectPath := filepath.Clean(projectPath)

	// Check if project directory exists
	if !fso.FileExists(cleanProjectPath) {
		return fmt.Errorf("project directory does not exist: %s", cleanProjectPath)
	}

	// Verify it's actually a directory
	info, err := os.Stat(cleanProjectPath)
	if err != nil {
		return fmt.Errorf("failed to stat project directory %s: %w", cleanProjectPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("project path %s exists but is not a directory", cleanProjectPath)
	}

	return nil
}

// RemoveFile safely removes a file with validation
func (fso *FileSystemOperations) RemoveFile(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Validate path before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	cleanPath := filepath.Clean(path)

	if fso.dryRun {
		return nil
	}

	// Check if file exists
	if !fso.FileExists(cleanPath) {
		return fmt.Errorf("file does not exist: %s", cleanPath)
	}

	// Remove the file
	if err := os.Remove(cleanPath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", cleanPath, err)
	}

	return nil
}

// RemoveDirectory safely removes a directory and all its contents with validation
func (fso *FileSystemOperations) RemoveDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Validate path before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	cleanPath := filepath.Clean(path)

	if fso.dryRun {
		return nil
	}

	// Check if directory exists
	if !fso.FileExists(cleanPath) {
		return fmt.Errorf("directory does not exist: %s", cleanPath)
	}

	// Verify it's actually a directory
	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to stat directory %s: %w", cleanPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %s exists but is not a directory", cleanPath)
	}

	// Remove the directory and all its contents
	if err := os.RemoveAll(cleanPath); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", cleanPath, err)
	}

	return nil
}

// GetFileInfo returns file information for the given path
func (fso *FileSystemOperations) GetFileInfo(path string) (os.FileInfo, error) {
	if path == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", cleanPath, err)
	}

	return info, nil
}

// IsDryRun returns whether the operations handler is in dry-run mode
func (fso *FileSystemOperations) IsDryRun() bool {
	return fso.dryRun
}

// SetDryRun sets the dry-run mode for the operations handler
func (fso *FileSystemOperations) SetDryRun(dryRun bool) {
	fso.dryRun = dryRun
}
