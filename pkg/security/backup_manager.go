package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// BackupManager handles file backup operations for safety
type BackupManager struct {
	backupDir       string
	maxBackups      int
	backupExtension string
	enabled         bool
}

// NewBackupManager creates a new backup manager
func NewBackupManager(backupDir string) *BackupManager {
	return &BackupManager{
		backupDir:       backupDir,
		maxBackups:      10, // Keep last 10 backups by default
		backupExtension: ".backup",
		enabled:         true,
	}
}

// BackupResult contains information about a backup operation
type BackupResult struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	Success      bool      `json:"success"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	FileSize     int64     `json:"file_size"`
	Checksum     string    `json:"checksum,omitempty"`
}

// BackupFile creates a backup of the specified file
func (bm *BackupManager) BackupFile(filePath string) (*BackupResult, error) {
	result := &BackupResult{
		OriginalPath: filePath,
		Timestamp:    time.Now(),
	}

	if !bm.enabled {
		result.Success = true
		result.BackupPath = "backup disabled"
		return result, nil
	}

	// Check if original file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		result.Success = true
		result.BackupPath = "file does not exist"
		return result, nil
	}
	if err != nil {
		result.Error = fmt.Sprintf("failed to stat original file: %v", err)
		return result, err
	}

	result.FileSize = fileInfo.Size()

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(bm.backupDir, 0750); err != nil {
		result.Error = fmt.Sprintf("failed to create backup directory: %v", err)
		return result, err
	}

	// Generate backup filename with timestamp
	filename := filepath.Base(filePath)
	timestamp := result.Timestamp.Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s_%s%s", filename, timestamp, bm.backupExtension)
	backupPath := filepath.Join(bm.backupDir, backupFilename)

	// Use secure file operations for backup
	sfo := NewSecureFileOperations([]string{bm.backupDir, filepath.Dir(filePath)})
	copyResult, err := sfo.SecureCopyFile(filePath, backupPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to create backup: %v", err)
		return result, err
	}

	result.BackupPath = backupPath
	result.Success = copyResult.Success
	result.Checksum = copyResult.Checksum

	// Clean up old backups
	if err := bm.cleanupOldBackups(filename); err != nil {
		// Log warning but don't fail the backup operation
		fmt.Printf("Warning: failed to cleanup old backups: %v\n", err)
	}

	return result, nil
}

// BackupDirectory creates backups of all files in a directory
func (bm *BackupManager) BackupDirectory(dirPath string) (map[string]*BackupResult, error) {
	results := make(map[string]*BackupResult)

	if !bm.enabled {
		return results, nil
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create backup for each file
		result, backupErr := bm.BackupFile(path)
		results[path] = result

		// Continue walking even if individual backup fails
		if backupErr != nil {
			fmt.Printf("Warning: failed to backup %s: %v\n", path, backupErr)
		}

		return nil
	})

	return results, err
}

// RestoreFile restores a file from backup
func (bm *BackupManager) RestoreFile(originalPath string, backupTimestamp time.Time) (*BackupResult, error) {
	result := &BackupResult{
		OriginalPath: originalPath,
		Timestamp:    time.Now(),
	}

	// Find backup file
	filename := filepath.Base(originalPath)
	timestampStr := backupTimestamp.Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s_%s%s", filename, timestampStr, bm.backupExtension)
	backupPath := filepath.Join(bm.backupDir, backupFilename)

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		result.Error = fmt.Sprintf("backup file not found: %s", backupPath)
		return result, fmt.Errorf("backup not found")
	}

	// Use secure file operations for restore
	sfo := NewSecureFileOperations([]string{filepath.Dir(originalPath), bm.backupDir})
	copyResult, err := sfo.SecureCopyFile(backupPath, originalPath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to restore from backup: %v", err)
		return result, err
	}

	result.BackupPath = backupPath
	result.Success = copyResult.Success
	result.Checksum = copyResult.Checksum

	return result, nil
}

// ListBackups returns a list of available backups for a file
func (bm *BackupManager) ListBackups(originalPath string) ([]BackupInfo, error) {
	var backups []BackupInfo
	filename := filepath.Base(originalPath)

	// Read backup directory
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return backups, fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Find backups for this file
	prefix := filename + "_"
	suffix := bm.backupExtension

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
			// Extract timestamp from filename
			timestampPart := strings.TrimPrefix(name, prefix)
			timestampPart = strings.TrimSuffix(timestampPart, suffix)

			timestamp, err := time.Parse("20060102_150405", timestampPart)
			if err != nil {
				continue // Skip invalid timestamp formats
			}

			// Get file info
			backupPath := filepath.Join(bm.backupDir, name)
			fileInfo, err := os.Stat(backupPath)
			if err != nil {
				continue
			}

			backups = append(backups, BackupInfo{
				OriginalPath: originalPath,
				BackupPath:   backupPath,
				Timestamp:    timestamp,
				Size:         fileInfo.Size(),
			})
		}
	}

	return backups, nil
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	OriginalPath string    `json:"original_path"`
	BackupPath   string    `json:"backup_path"`
	Timestamp    time.Time `json:"timestamp"`
	Size         int64     `json:"size"`
}

// cleanupOldBackups removes old backups beyond the maximum limit
func (bm *BackupManager) cleanupOldBackups(filename string) error {
	backups, err := bm.ListBackups(filename)
	if err != nil {
		return err
	}

	// Sort backups by timestamp (newest first)
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].Timestamp.Before(backups[j].Timestamp) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}

	// Remove old backups beyond the limit
	if len(backups) > bm.maxBackups {
		for i := bm.maxBackups; i < len(backups); i++ {
			if err := os.Remove(backups[i].BackupPath); err != nil {
				return fmt.Errorf("failed to remove old backup %s: %w", backups[i].BackupPath, err)
			}
		}
	}

	return nil
}

// SetEnabled enables or disables backup functionality
func (bm *BackupManager) SetEnabled(enabled bool) {
	bm.enabled = enabled
}

// IsEnabled returns whether backup functionality is enabled
func (bm *BackupManager) IsEnabled() bool {
	return bm.enabled
}

// SetMaxBackups sets the maximum number of backups to keep
func (bm *BackupManager) SetMaxBackups(max int) {
	if max > 0 {
		bm.maxBackups = max
	}
}

// GetBackupDirectory returns the backup directory path
func (bm *BackupManager) GetBackupDirectory() string {
	return bm.backupDir
}

// CleanupAllBackups removes all backup files
func (bm *BackupManager) CleanupAllBackups() error {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), bm.backupExtension) {
			backupPath := filepath.Join(bm.backupDir, entry.Name())
			if err := os.Remove(backupPath); err != nil {
				return fmt.Errorf("failed to remove backup %s: %w", backupPath, err)
			}
		}
	}

	return nil
}

// GetBackupStats returns statistics about backups
func (bm *BackupManager) GetBackupStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"backup_directory": bm.backupDir,
		"enabled":          bm.enabled,
		"max_backups":      bm.maxBackups,
		"total_backups":    0,
		"total_size":       int64(0),
		"oldest_backup":    nil,
		"newest_backup":    nil,
	}

	if !bm.enabled {
		return stats, nil
	}

	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return stats, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var oldestTime, newestTime time.Time
	totalSize := int64(0)
	backupCount := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), bm.backupExtension) {
			continue
		}

		backupCount++

		// Get file info
		backupPath := filepath.Join(bm.backupDir, entry.Name())
		fileInfo, err := os.Stat(backupPath)
		if err != nil {
			continue
		}

		totalSize += fileInfo.Size()
		modTime := fileInfo.ModTime()

		if oldestTime.IsZero() || modTime.Before(oldestTime) {
			oldestTime = modTime
		}
		if newestTime.IsZero() || modTime.After(newestTime) {
			newestTime = modTime
		}
	}

	stats["total_backups"] = backupCount
	stats["total_size"] = totalSize

	if !oldestTime.IsZero() {
		stats["oldest_backup"] = oldestTime
	}
	if !newestTime.IsZero() {
		stats["newest_backup"] = newestTime
	}

	return stats, nil
}
