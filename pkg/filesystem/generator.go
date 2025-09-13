package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// Generator implements the FileSystemGenerator interface
type Generator struct {
	// dryRun indicates whether to actually create files or just simulate
	dryRun bool
}

// NewGenerator creates a new filesystem generator
func NewGenerator() interfaces.FileSystemGenerator {
	return &Generator{
		dryRun: false,
	}
}

// NewDryRunGenerator creates a new filesystem generator in dry-run mode
func NewDryRunGenerator() interfaces.FileSystemGenerator {
	return &Generator{
		dryRun: true,
	}
}

// CreateProject creates the complete project structure based on configuration
func (g *Generator) CreateProject(config *models.ProjectConfig, outputPath string) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Ensure the output path is absolute
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", outputPath, err)
	}

	// Create the root project directory
	projectPath := filepath.Join(absOutputPath, config.Name)
	if err := g.EnsureDirectory(projectPath); err != nil {
		return fmt.Errorf("failed to create project directory %s: %w", projectPath, err)
	}

	return nil
}

// CreateDirectory creates a directory with proper permissions
func (g *Generator) CreateDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Validate that the path doesn't contain dangerous elements before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	// Clean the path after validation
	cleanPath := filepath.Clean(path)

	if g.dryRun {
		return nil
	}

	// Create directory with proper permissions (0755)
	if err := os.MkdirAll(cleanPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", cleanPath, err)
	}

	return nil
}

// WriteFile writes content to a file with specified permissions
func (g *Generator) WriteFile(path string, content []byte, perm os.FileMode) error {
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

	if g.dryRun {
		return nil
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(cleanPath)
	if err := g.EnsureDirectory(dir); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", cleanPath, err)
	}

	// Write the file with specified permissions
	if err := os.WriteFile(cleanPath, content, perm); err != nil {
		return fmt.Errorf("failed to write file %s: %w", cleanPath, err)
	}

	return nil
}

// CopyAssets copies binary assets from source to destination
func (g *Generator) CopyAssets(srcDir, destDir string) error {
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

	// Clean paths after validation
	cleanSrcDir := filepath.Clean(srcDir)
	cleanDestDir := filepath.Clean(destDir)

	// Check if source directory exists
	if !g.FileExists(cleanSrcDir) {
		return fmt.Errorf("source directory does not exist: %s", cleanSrcDir)
	}

	if g.dryRun {
		return nil
	}

	// Ensure destination directory exists
	if err := g.EnsureDirectory(cleanDestDir); err != nil {
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
			return g.CreateDirectory(destPath)
		}

		// Copy file
		return g.copyFile(srcPath, destPath, info.Mode())
	})
}

// CreateSymlink creates a symbolic link
func (g *Generator) CreateSymlink(target, link string) error {
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

	if g.dryRun {
		return nil
	}

	// Ensure parent directory of link exists
	linkDir := filepath.Dir(cleanLink)
	if err := g.EnsureDirectory(linkDir); err != nil {
		return fmt.Errorf("failed to create parent directory for symlink %s: %w", cleanLink, err)
	}

	// Create symbolic link
	if err := os.Symlink(cleanTarget, cleanLink); err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", cleanTarget, cleanLink, err)
	}

	return nil
}

// FileExists checks if a file exists at the given path
func (g *Generator) FileExists(path string) bool {
	if path == "" {
		return false
	}

	cleanPath := filepath.Clean(path)
	_, err := os.Stat(cleanPath)
	return err == nil
}

// EnsureDirectory ensures a directory exists, creating it if necessary
func (g *Generator) EnsureDirectory(path string) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Validate path before cleaning
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: path traversal detected in %s", path)
	}

	cleanPath := filepath.Clean(path)

	if g.dryRun {
		return nil
	}

	// Check if directory already exists
	if g.FileExists(cleanPath) {
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
	return g.CreateDirectory(cleanPath)
}

// copyFile copies a single file from source to destination using optimized I/O
func (g *Generator) copyFile(srcPath, destPath string, perm os.FileMode) error {
	if g.dryRun {
		return nil
	}

	// Use optimized copy for better performance
	return OptimizedCopy(srcPath, destPath, perm)
}
