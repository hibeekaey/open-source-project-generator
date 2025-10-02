package security

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupManager_BackupDirectory(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	sourceDir := filepath.Join(tempDir, "source")

	// Create source directory with test files
	err := os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)

	testFiles := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(sourceDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	bm := NewBackupManager(backupDir)
	results, err := bm.BackupDirectory(sourceDir)

	assert.NoError(t, err)
	assert.Len(t, results, 3) // Should backup all 3 files

	// Verify all files were backed up successfully
	for filePath, result := range results {
		assert.True(t, result.Success, "Backup should succeed for %s", filePath)
		assert.NotEmpty(t, result.BackupPath, "Backup path should be set for %s", filePath)
		assert.Greater(t, result.FileSize, int64(0), "File size should be positive for %s", filePath)
		assert.NotEmpty(t, result.Checksum, "Checksum should be calculated for %s", filePath)

		// Verify backup file exists
		_, err := os.Stat(result.BackupPath)
		assert.NoError(t, err, "Backup file should exist for %s", filePath)
	}
}

func TestBackupManager_RestoreFile(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "original.txt")
	originalContent := "original content"

	// Create original file
	err := os.WriteFile(originalFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)

	// Create backup
	backupResult, err := bm.BackupFile(originalFile)
	require.NoError(t, err)
	require.True(t, backupResult.Success)

	// Modify original file
	modifiedContent := "modified content"
	err = os.WriteFile(originalFile, []byte(modifiedContent), 0644)
	require.NoError(t, err)

	// Restore from backup
	restoreResult, err := bm.RestoreFile(originalFile, backupResult.Timestamp)
	assert.NoError(t, err)
	assert.True(t, restoreResult.Success)
	assert.Equal(t, backupResult.BackupPath, restoreResult.BackupPath)

	// Verify file was restored
	restoredContent, err := os.ReadFile(originalFile)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, string(restoredContent))
}

func TestBackupManager_RestoreFile_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "nonexistent.txt")

	bm := NewBackupManager(backupDir)

	// Try to restore non-existent backup
	pastTime := time.Now().Add(-1 * time.Hour)
	result, err := bm.RestoreFile(originalFile, pastTime)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backup not found")
	assert.Contains(t, result.Error, "backup file not found")
}

func TestBackupManager_SetMaxBackups(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")

	bm := NewBackupManager(backupDir)

	// Test setting valid max backups
	bm.SetMaxBackups(5)
	assert.Equal(t, 5, bm.maxBackups)

	// Test setting invalid max backups (should be ignored)
	bm.SetMaxBackups(0)
	assert.Equal(t, 5, bm.maxBackups) // Should remain unchanged

	bm.SetMaxBackups(-1)
	assert.Equal(t, 5, bm.maxBackups) // Should remain unchanged
}

func TestBackupManager_GetBackupDirectory(t *testing.T) {
	backupDir := "/test/backup/dir"
	bm := NewBackupManager(backupDir)

	assert.Equal(t, backupDir, bm.GetBackupDirectory())
}

func TestBackupManager_CleanupAllBackups(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "test.txt")

	// Create original file
	err := os.WriteFile(originalFile, []byte("test content"), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)

	// Create multiple backups
	for i := 0; i < 3; i++ {
		_, err := bm.BackupFile(originalFile)
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	}

	// Verify backups exist
	entries, err := os.ReadDir(backupDir)
	require.NoError(t, err)
	assert.Greater(t, len(entries), 0)

	// Cleanup all backups
	err = bm.CleanupAllBackups()
	assert.NoError(t, err)

	// Verify all backups are removed
	entries, err = os.ReadDir(backupDir)
	require.NoError(t, err)

	backupCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".backup" {
			backupCount++
		}
	}
	assert.Equal(t, 0, backupCount)
}

func TestBackupManager_GetBackupStats(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "test.txt")

	bm := NewBackupManager(backupDir)

	// Test stats when disabled
	bm.SetEnabled(false)
	stats, err := bm.GetBackupStats()
	assert.NoError(t, err)
	assert.Equal(t, false, stats["enabled"])
	assert.Equal(t, 0, stats["total_backups"])

	// Enable and create backups
	bm.SetEnabled(true)

	// Create original file
	err = os.WriteFile(originalFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// Create backups with longer delays to ensure different timestamps
	for i := 0; i < 2; i++ {
		_, err := bm.BackupFile(originalFile)
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond) // Increased delay
	}

	// Get stats
	stats, err = bm.GetBackupStats()
	assert.NoError(t, err)
	assert.Equal(t, true, stats["enabled"])
	assert.Equal(t, backupDir, stats["backup_directory"])
	assert.Equal(t, 10, stats["max_backups"])
	assert.Equal(t, 2, stats["total_backups"])
	assert.Greater(t, stats["total_size"].(int64), int64(0))
	assert.NotNil(t, stats["oldest_backup"])
	assert.NotNil(t, stats["newest_backup"])
}

func TestBackupManager_GetBackupStats_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")

	bm := NewBackupManager(backupDir)

	// Create empty backup directory
	err := os.MkdirAll(backupDir, 0755)
	require.NoError(t, err)

	stats, err := bm.GetBackupStats()
	assert.NoError(t, err)
	assert.Equal(t, 0, stats["total_backups"])
	assert.Equal(t, int64(0), stats["total_size"])
	assert.Nil(t, stats["oldest_backup"])
	assert.Nil(t, stats["newest_backup"])
}

func TestBackupManager_CleanupOldBackups(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "test.txt")

	// Create original file
	err := os.WriteFile(originalFile, []byte("test content"), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)
	bm.SetMaxBackups(2) // Keep only 2 backups

	// Create 4 backups
	for i := 0; i < 4; i++ {
		_, err := bm.BackupFile(originalFile)
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	}

	// List backups - should only have 2 (cleanup happens automatically)
	backups, err := bm.ListBackups(originalFile)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(backups), 2)
}

func TestBackupManager_BackupFile_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "test.txt")

	// Create original file
	err := os.WriteFile(originalFile, []byte("test content"), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)
	bm.SetEnabled(false)

	result, err := bm.BackupFile(originalFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "backup disabled", result.BackupPath)
}

func TestBackupManager_BackupFile_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")

	bm := NewBackupManager(backupDir)

	result, err := bm.BackupFile(nonExistentFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "file does not exist", result.BackupPath)
}

func TestBackupManager_BackupDirectory_WithErrors(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	sourceDir := filepath.Join(tempDir, "source")

	// Create source directory
	err := os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)

	// Create a file
	testFile := filepath.Join(sourceDir, "test.txt")
	err = os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)

	// Test with non-existent source directory
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	results, err := bm.BackupDirectory(nonExistentDir)
	assert.Error(t, err)
	assert.Empty(t, results)

	// Test with valid directory
	results, err = bm.BackupDirectory(sourceDir)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestBackupManager_ListBackups_InvalidTimestamp(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")

	// Create backup directory
	err := os.MkdirAll(backupDir, 0755)
	require.NoError(t, err)

	// Create a file with invalid timestamp format
	invalidBackupFile := filepath.Join(backupDir, "test.txt_invalid_timestamp.backup")
	err = os.WriteFile(invalidBackupFile, []byte("content"), 0644)
	require.NoError(t, err)

	bm := NewBackupManager(backupDir)
	backups, err := bm.ListBackups("test.txt")

	assert.NoError(t, err)
	assert.Empty(t, backups) // Should skip files with invalid timestamps
}

func TestBackupManager_Integration(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	originalFile := filepath.Join(tempDir, "integration_test.txt")

	bm := NewBackupManager(backupDir)
	bm.SetMaxBackups(3)

	// Test complete workflow
	originalContent := "original content"
	err := os.WriteFile(originalFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	// 1. Create initial backup
	result1, err := bm.BackupFile(originalFile)
	require.NoError(t, err)
	assert.True(t, result1.Success)

	// 2. Modify file and create another backup
	time.Sleep(100 * time.Millisecond)
	modifiedContent := "modified content"
	err = os.WriteFile(originalFile, []byte(modifiedContent), 0644)
	require.NoError(t, err)

	result2, err := bm.BackupFile(originalFile)
	require.NoError(t, err)
	assert.True(t, result2.Success)

	// 3. List backups
	backups, err := bm.ListBackups(originalFile)
	assert.NoError(t, err)
	assert.Len(t, backups, 2)

	// 4. Restore from first backup
	_, err = bm.RestoreFile(originalFile, result1.Timestamp)
	assert.NoError(t, err)

	// 5. Verify restoration
	restoredContent, err := os.ReadFile(originalFile)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, string(restoredContent))

	// 6. Get stats
	stats, err := bm.GetBackupStats()
	assert.NoError(t, err)
	assert.Equal(t, 2, stats["total_backups"])

	// 7. Cleanup all backups
	err = bm.CleanupAllBackups()
	assert.NoError(t, err)

	// 8. Verify cleanup
	backups, err = bm.ListBackups(originalFile)
	assert.NoError(t, err)
	assert.Empty(t, backups)
}
