// Package filesystem provides secure file system operations.
package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
)

// Operations provides secure file system operations with validation.
type Operations struct {
	validator *security.Validator
}

// NewOperations creates a new file system operations handler.
func NewOperations() *Operations {
	return &Operations{
		validator: security.NewValidator(),
	}
}

// CreateFile creates a file with secure permissions and path validation.
func (ops *Operations) CreateFile(path string, data []byte) error {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := ops.CreateDir(dir); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write file with secure permissions (0600)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// CreateDir creates a directory with secure permissions and path validation.
func (ops *Operations) CreateDir(path string) error {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return fmt.Errorf("invalid directory path: %w", err)
	}

	// Create directory with secure permissions (0750)
	if err := os.MkdirAll(path, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ReadFile reads a file with path validation.
func (ops *Operations) ReadFile(path string) ([]byte, error) {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Read file
	data, err := os.ReadFile(path) // #nosec G304 - Path validated by ValidatePathSecurity above
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// WriteFile writes a file with secure permissions and path validation.
func (ops *Operations) WriteFile(path string, data []byte) error {
	return ops.CreateFile(path, data)
}

// DeleteFile deletes a file with path validation.
func (ops *Operations) DeleteFile(path string) error {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	// Delete file
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// DeleteDir deletes a directory and its contents with path validation.
func (ops *Operations) DeleteDir(path string) error {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return fmt.Errorf("invalid directory path: %w", err)
	}

	// Delete directory
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete directory: %w", err)
	}

	return nil
}

// CopyFile copies a file from source to destination with path validation.
func (ops *Operations) CopyFile(src, dst string) error {
	// Validate paths
	if err := ops.validator.ValidatePathSecurity(src); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	if err := ops.validator.ValidatePathSecurity(dst); err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(src) // #nosec G304 - Path validated by ValidatePathSecurity above
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := ops.CreateDir(dstDir); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file with secure permissions
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600) // #nosec G304 - Path validated by ValidatePathSecurity above
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

// MoveFile moves a file from source to destination with path validation.
func (ops *Operations) MoveFile(src, dst string) error {
	// Validate paths
	if err := ops.validator.ValidatePathSecurity(src); err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	if err := ops.validator.ValidatePathSecurity(dst); err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	// Ensure destination directory exists
	dstDir := filepath.Dir(dst)
	if err := ops.CreateDir(dstDir); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Try to rename (atomic operation if on same filesystem)
	if err := os.Rename(src, dst); err != nil {
		// If rename fails, fall back to copy and delete
		if err := ops.CopyFile(src, dst); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
		if err := ops.DeleteFile(src); err != nil {
			return fmt.Errorf("failed to delete source file: %w", err)
		}
	}

	return nil
}

// FileExists checks if a file exists.
func (ops *Operations) FileExists(path string) (bool, error) {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return false, fmt.Errorf("invalid file path: %w", err)
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file existence: %w", err)
}

// DirExists checks if a directory exists.
func (ops *Operations) DirExists(path string) (bool, error) {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return false, fmt.Errorf("invalid directory path: %w", err)
	}

	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check directory existence: %w", err)
}

// ListDir lists all files and directories in a directory.
func (ops *Operations) ListDir(path string) ([]os.FileInfo, error) {
	// Validate path
	if err := ops.validator.ValidatePathSecurity(path); err != nil {
		return nil, fmt.Errorf("invalid directory path: %w", err)
	}

	// Open directory
	dir, err := os.Open(path) // #nosec G304 - Path validated by ValidatePathSecurity in calling function
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}
	defer func() {
		if err := dir.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	// Read directory contents
	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return entries, nil
}

// CleanupTemp removes temporary files and directories.
func (ops *Operations) CleanupTemp(paths []string) error {
	var errors []error

	for _, path := range paths {
		// Validate path
		if err := ops.validator.ValidatePathSecurity(path); err != nil {
			errors = append(errors, fmt.Errorf("invalid path %s: %w", path, err))
			continue
		}

		// Check if path exists
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Path doesn't exist, skip
			}
			errors = append(errors, fmt.Errorf("failed to stat %s: %w", path, err))
			continue
		}

		// Delete based on type
		if info.IsDir() {
			if err := ops.DeleteDir(path); err != nil {
				errors = append(errors, err)
			}
		} else {
			if err := ops.DeleteFile(path); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup failed with %d errors: %w", len(errors), errors[0])
	}

	return nil
}
