package cleanup

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles creating and managing backups of files before modifications
type BackupManager struct {
	backupDir string
}

// Backup represents a backup of files
type Backup struct {
	ID        string
	Timestamp time.Time
	Files     []BackupFile
	BasePath  string
}

// BackupFile represents a single backed up file
type BackupFile struct {
	OriginalPath string
	BackupPath   string
	Checksum     string
	Size         int64
}

// NewBackupManager creates a new backup manager
func NewBackupManager(backupDir string) *BackupManager {
	return &BackupManager{
		backupDir: backupDir,
	}
}

// CreateBackup creates a backup of the specified files
func (bm *BackupManager) CreateBackup(files []string) (*Backup, error) {
	timestamp := time.Now()
	backupID := fmt.Sprintf("cleanup_%d", timestamp.Unix())

	backup := &Backup{
		ID:        backupID,
		Timestamp: timestamp,
		Files:     make([]BackupFile, 0, len(files)),
		BasePath:  bm.backupDir,
	}

	// Create backup directory
	backupPath := filepath.Join(bm.backupDir, backupID)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Backup each file
	for _, file := range files {
		backupFile, err := bm.backupFile(file, backupPath)
		if err != nil {
			return nil, fmt.Errorf("failed to backup file %s: %w", file, err)
		}
		backup.Files = append(backup.Files, *backupFile)
	}

	return backup, nil
}

// backupFile creates a backup of a single file
func (bm *BackupManager) backupFile(originalPath, backupDir string) (*BackupFile, error) {
	// Open original file
	src, err := os.Open(originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	// Get file info
	info, err := src.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create backup file path
	relPath, err := filepath.Rel(".", originalPath)
	if err != nil {
		relPath = originalPath
	}
	backupPath := filepath.Join(backupDir, relPath)

	// Create backup directory structure
	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory structure: %w", err)
	}

	// Create backup file
	dst, err := os.Create(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	// Copy file content and calculate checksum
	hasher := sha256.New()
	multiWriter := io.MultiWriter(dst, hasher)

	if _, err := io.Copy(multiWriter, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	return &BackupFile{
		OriginalPath: originalPath,
		BackupPath:   backupPath,
		Checksum:     checksum,
		Size:         info.Size(),
	}, nil
}

// RestoreBackup restores files from a backup
func (bm *BackupManager) RestoreBackup(backup *Backup) error {
	for _, file := range backup.Files {
		if err := bm.restoreFile(file); err != nil {
			return fmt.Errorf("failed to restore file %s: %w", file.OriginalPath, err)
		}
	}
	return nil
}

// restoreFile restores a single file from backup
func (bm *BackupManager) restoreFile(backupFile BackupFile) error {
	// Open backup file
	src, err := os.Open(backupFile.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer src.Close()

	// Create original file directory if needed
	if err := os.MkdirAll(filepath.Dir(backupFile.OriginalPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create original file
	dst, err := os.Create(backupFile.OriginalPath)
	if err != nil {
		return fmt.Errorf("failed to create original file: %w", err)
	}
	defer dst.Close()

	// Copy content back
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	return nil
}

// ListBackups returns a list of available backups
func (bm *BackupManager) ListBackups() ([]*Backup, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Backup{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []*Backup
	for _, entry := range entries {
		if entry.IsDir() {
			// Parse backup directory name to extract timestamp
			// This is a simplified implementation
			backup := &Backup{
				ID:       entry.Name(),
				BasePath: filepath.Join(bm.backupDir, entry.Name()),
			}
			backups = append(backups, backup)
		}
	}

	return backups, nil
}

// CleanupOldBackups removes backups older than the specified duration
func (bm *BackupManager) CleanupOldBackups(maxAge time.Duration) error {
	backups, err := bm.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)

	for _, backup := range backups {
		if backup.Timestamp.Before(cutoff) {
			backupPath := filepath.Join(bm.backupDir, backup.ID)
			if err := os.RemoveAll(backupPath); err != nil {
				return fmt.Errorf("failed to remove old backup %s: %w", backup.ID, err)
			}
		}
	}

	return nil
}
