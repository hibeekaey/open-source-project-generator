// Package filesystem provides secure file system operations.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
)

// BackupManager manages backup and restore operations.
type BackupManager struct {
	validator       *security.Validator
	ops             *Operations
	backupDirectory string
}

// NewBackupManager creates a new backup manager.
func NewBackupManager(backupDir string) (*BackupManager, error) {
	validator := security.NewValidator()

	// Validate backup directory path
	if err := validator.ValidatePathSecurity(backupDir); err != nil {
		return nil, fmt.Errorf("invalid backup directory: %w", err)
	}

	ops := NewOperations()

	// Create backup directory if it doesn't exist
	if err := ops.CreateDir(backupDir); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return &BackupManager{
		validator:       validator,
		ops:             ops,
		backupDirectory: backupDir,
	}, nil
}

// Backup creates a backup of a file or directory.
// Returns the backup path.
func (bm *BackupManager) Backup(path string) (string, error) {
	// Validate path
	if err := bm.validator.ValidatePathSecurity(path); err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Check if path exists
	exists, err := bm.ops.FileExists(path)
	if err != nil {
		return "", fmt.Errorf("failed to check path existence: %w", err)
	}
	if !exists {
		// Check if it's a directory
		exists, err = bm.ops.DirExists(path)
		if err != nil {
			return "", fmt.Errorf("failed to check directory existence: %w", err)
		}
		if !exists {
			return "", fmt.Errorf("path does not exist: %s", path)
		}
	}

	// Generate backup path with timestamp
	timestamp := time.Now().Format("20060102-150405")
	baseName := filepath.Base(path)
	backupName := fmt.Sprintf("%s.%s.backup", baseName, timestamp)
	backupPath := filepath.Join(bm.backupDirectory, backupName)

	// Check if it's a file or directory
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() {
		// Backup directory
		if err := bm.backupDir(path, backupPath); err != nil {
			return "", fmt.Errorf("failed to backup directory: %w", err)
		}
	} else {
		// Backup file
		if err := bm.ops.CopyFile(path, backupPath); err != nil {
			return "", fmt.Errorf("failed to backup file: %w", err)
		}
	}

	return backupPath, nil
}

// Restore restores a file or directory from backup.
func (bm *BackupManager) Restore(backupPath, targetPath string) error {
	// Validate paths
	if err := bm.validator.ValidatePathSecurity(backupPath); err != nil {
		return fmt.Errorf("invalid backup path: %w", err)
	}
	if err := bm.validator.ValidatePathSecurity(targetPath); err != nil {
		return fmt.Errorf("invalid target path: %w", err)
	}

	// Check if backup exists
	exists, err := bm.ops.FileExists(backupPath)
	if err != nil {
		return fmt.Errorf("failed to check backup existence: %w", err)
	}
	if !exists {
		// Check if it's a directory
		exists, err = bm.ops.DirExists(backupPath)
		if err != nil {
			return fmt.Errorf("failed to check backup directory existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("backup does not exist: %s", backupPath)
		}
	}

	// Check if it's a file or directory
	info, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("failed to stat backup: %w", err)
	}

	if info.IsDir() {
		// Restore directory
		if err := bm.restoreDir(backupPath, targetPath); err != nil {
			return fmt.Errorf("failed to restore directory: %w", err)
		}
	} else {
		// Restore file
		if err := bm.ops.CopyFile(backupPath, targetPath); err != nil {
			return fmt.Errorf("failed to restore file: %w", err)
		}
	}

	return nil
}

// DeleteBackup deletes a backup.
func (bm *BackupManager) DeleteBackup(backupPath string) error {
	// Validate path
	if err := bm.validator.ValidatePathSecurity(backupPath); err != nil {
		return fmt.Errorf("invalid backup path: %w", err)
	}

	// Check if backup is in backup directory
	relPath, err := filepath.Rel(bm.backupDirectory, backupPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("backup path is not in backup directory")
	}

	// Check if it's a file or directory
	info, err := os.Stat(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to stat backup: %w", err)
	}

	if info.IsDir() {
		return bm.ops.DeleteDir(backupPath)
	}
	return bm.ops.DeleteFile(backupPath)
}

// ListBackups lists all backups in the backup directory.
func (bm *BackupManager) ListBackups() ([]string, error) {
	entries, err := bm.ops.ListDir(bm.backupDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		backupPath := filepath.Join(bm.backupDirectory, entry.Name())
		backups = append(backups, backupPath)
	}

	return backups, nil
}

// CleanupOldBackups removes backups older than the specified duration.
func (bm *BackupManager) CleanupOldBackups(maxAge time.Duration) error {
	entries, err := bm.ops.ListDir(bm.backupDirectory)
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	now := time.Now()
	var errors []error

	for _, entry := range entries {
		// Check if backup is older than maxAge
		if now.Sub(entry.ModTime()) > maxAge {
			backupPath := filepath.Join(bm.backupDirectory, entry.Name())
			if err := bm.DeleteBackup(backupPath); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup failed with %d errors: %w", len(errors), errors[0])
	}

	return nil
}

// backupDir recursively backs up a directory.
func (bm *BackupManager) backupDir(srcDir, dstDir string) error {
	// Create destination directory
	if err := bm.ops.CreateDir(dstDir); err != nil {
		return err
	}

	// List source directory contents
	entries, err := bm.ops.ListDir(srcDir)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// Recursively backup subdirectory
			if err := bm.backupDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := bm.ops.CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// restoreDir recursively restores a directory.
func (bm *BackupManager) restoreDir(srcDir, dstDir string) error {
	// Create destination directory
	if err := bm.ops.CreateDir(dstDir); err != nil {
		return err
	}

	// List source directory contents
	entries, err := bm.ops.ListDir(srcDir)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// Recursively restore subdirectory
			if err := bm.restoreDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := bm.ops.CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
