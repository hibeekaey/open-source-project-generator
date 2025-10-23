package orchestrator

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

// RollbackManager manages rollback operations for failed generations
type RollbackManager struct {
	backups       map[string]string // original path -> backup path
	tempDirs      []string          // temporary directories to clean up
	mu            sync.Mutex
	logger        *logger.Logger
	backupManager *filesystem.BackupManager
}

// NewRollbackManager creates a new rollback manager
func NewRollbackManager(log *logger.Logger) *RollbackManager {
	// Create backup manager with default backup directory
	backupMgr, err := filesystem.NewBackupManager(".backups")
	if err != nil {
		if log != nil {
			log.Warn(fmt.Sprintf("Failed to create backup manager: %v", err))
		}
		backupMgr = nil
	}

	return &RollbackManager{
		backups:       make(map[string]string),
		tempDirs:      make([]string, 0),
		logger:        log,
		backupManager: backupMgr,
	}
}

// RegisterBackup registers a backup for potential rollback
func (rm *RollbackManager) RegisterBackup(originalPath, backupPath string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.backups[originalPath] = backupPath

	if rm.logger != nil {
		rm.logger.Debug(fmt.Sprintf("Registered backup: %s -> %s", originalPath, backupPath))
	}
}

// CreateBackup creates a backup of the specified path and registers it for rollback
func (rm *RollbackManager) CreateBackup(path string) (string, error) {
	if rm.backupManager == nil {
		return "", fmt.Errorf("backup manager not initialized")
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Path doesn't exist, no backup needed
		if rm.logger != nil {
			rm.logger.Debug(fmt.Sprintf("Path does not exist, skipping backup: %s", path))
		}
		return "", nil
	}

	if rm.logger != nil {
		rm.logger.Info(fmt.Sprintf("Creating backup of: %s", path))
	}

	// Create backup
	backupPath, err := rm.backupManager.Backup(path)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	// Register backup for potential rollback
	rm.RegisterBackup(path, backupPath)

	if rm.logger != nil {
		rm.logger.Info(fmt.Sprintf("Backup created: %s", backupPath))
	}

	return backupPath, nil
}

// RegisterTempDir registers a temporary directory for cleanup
func (rm *RollbackManager) RegisterTempDir(tempDir string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.tempDirs = append(rm.tempDirs, tempDir)

	if rm.logger != nil {
		rm.logger.Debug(fmt.Sprintf("Registered temp directory: %s", tempDir))
	}
}

// Rollback performs rollback operations
func (rm *RollbackManager) Rollback(ctx context.Context) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.logger != nil {
		rm.logger.Info("Starting rollback...")
	}

	var errors []error

	// Clean up temporary directories
	for _, tempDir := range rm.tempDirs {
		if err := rm.cleanupTempDir(tempDir); err != nil {
			errors = append(errors, fmt.Errorf("failed to cleanup temp dir %s: %w", tempDir, err))
		}
	}

	// Restore from backups
	for originalPath, backupPath := range rm.backups {
		if err := rm.restoreBackup(originalPath, backupPath); err != nil {
			errors = append(errors, fmt.Errorf("failed to restore backup for %s: %w", originalPath, err))
		}
	}

	if len(errors) > 0 {
		if rm.logger != nil {
			rm.logger.Error(fmt.Sprintf("Rollback completed with %d errors", len(errors)))
		}
		return fmt.Errorf("rollback completed with errors: %v", errors)
	}

	if rm.logger != nil {
		rm.logger.Info("Rollback completed successfully")
	}

	// Clear registrations after successful rollback
	rm.backups = make(map[string]string)
	rm.tempDirs = make([]string, 0)

	return nil
}

// cleanupTempDir removes a temporary directory
func (rm *RollbackManager) cleanupTempDir(tempDir string) error {
	// Check if directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		// Already cleaned up
		return nil
	}

	// Remove directory
	if err := os.RemoveAll(tempDir); err != nil {
		return err
	}

	if rm.logger != nil {
		rm.logger.Debug(fmt.Sprintf("Cleaned up temp directory: %s", tempDir))
	}

	return nil
}

// restoreBackup restores a backup to its original location
func (rm *RollbackManager) restoreBackup(originalPath, backupPath string) error {
	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		if rm.logger != nil {
			rm.logger.Warn(fmt.Sprintf("Backup not found: %s", backupPath))
		}
		return nil
	}

	if rm.logger != nil {
		rm.logger.Info(fmt.Sprintf("Restoring backup from %s to %s", backupPath, originalPath))
	}

	// Use backup manager if available for more robust restoration
	if rm.backupManager != nil {
		if err := rm.backupManager.Restore(backupPath, originalPath); err != nil {
			return fmt.Errorf("failed to restore backup using backup manager: %w", err)
		}
	} else {
		// Fallback to simple restore
		// Remove current directory if it exists
		if _, err := os.Stat(originalPath); err == nil {
			if err := os.RemoveAll(originalPath); err != nil {
				return fmt.Errorf("failed to remove current directory: %w", err)
			}
		}

		// Restore backup
		if err := os.Rename(backupPath, originalPath); err != nil {
			return fmt.Errorf("failed to restore backup: %w", err)
		}
	}

	if rm.logger != nil {
		rm.logger.Info(fmt.Sprintf("Successfully restored backup: %s", originalPath))
	}

	return nil
}

// CleanupTempDirs cleans up all registered temporary directories without restoring backups
func (rm *RollbackManager) CleanupTempDirs() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	var errors []error

	for _, tempDir := range rm.tempDirs {
		if err := rm.cleanupTempDir(tempDir); err != nil {
			errors = append(errors, err)
		}
	}

	// Clear temp dirs after cleanup
	rm.tempDirs = make([]string, 0)

	if len(errors) > 0 {
		return fmt.Errorf("cleanup completed with errors: %v", errors)
	}

	return nil
}

// Clear clears all registered backups and temp directories without performing any operations
func (rm *RollbackManager) Clear() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.backups = make(map[string]string)
	rm.tempDirs = make([]string, 0)

	if rm.logger != nil {
		rm.logger.Debug("Rollback manager cleared")
	}
}

// HasBackups returns true if there are registered backups
func (rm *RollbackManager) HasBackups() bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	return len(rm.backups) > 0
}

// HasTempDirs returns true if there are registered temp directories
func (rm *RollbackManager) HasTempDirs() bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	return len(rm.tempDirs) > 0
}

// GetBackups returns a copy of the registered backups
func (rm *RollbackManager) GetBackups() map[string]string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	backups := make(map[string]string, len(rm.backups))
	for k, v := range rm.backups {
		backups[k] = v
	}

	return backups
}

// GetTempDirs returns a copy of the registered temp directories
func (rm *RollbackManager) GetTempDirs() []string {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	tempDirs := make([]string, len(rm.tempDirs))
	copy(tempDirs, rm.tempDirs)

	return tempDirs
}
