package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecureFileOperations_ValidateFilePath(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	tests := []struct {
		name      string
		path      string
		operation string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "empty path",
			path:      "",
			operation: "read",
			wantError: true,
			errorMsg:  "file path cannot be empty",
		},
		{
			name:      "path traversal",
			path:      "../../../etc/passwd",
			operation: "read",
			wantError: true,
			errorMsg:  "path traversal detected",
		},
		{
			name:      "blocked path - /etc/passwd",
			path:      "/etc/passwd",
			operation: "read",
			wantError: true,
			errorMsg:  "access to blocked path denied",
		},
		{
			name:      "blocked path - Windows System32",
			path:      "C:\\Windows\\System32\\config\\SAM",
			operation: "read",
			wantError: true,
			errorMsg:  "access to blocked path denied",
		},
		{
			name:      "path outside allowed directories",
			path:      "/tmp/outside.txt",
			operation: "read",
			wantError: true,
			errorMsg:  "path not within allowed directories",
		},
		{
			name:      "disallowed file extension",
			path:      filepath.Join(tempDir, "malicious.exe"),
			operation: "write",
			wantError: true,
			errorMsg:  "file extension '.exe' is not allowed",
		},
		{
			name:      "valid path within allowed directory",
			path:      filepath.Join(tempDir, "valid.txt"),
			operation: "read",
			wantError: false,
		},
		{
			name:      "valid path with allowed extension",
			path:      filepath.Join(tempDir, "config.yaml"),
			operation: "write",
			wantError: false,
		},
		{
			name:      "valid path without extension",
			path:      filepath.Join(tempDir, "Dockerfile"),
			operation: "write",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.ValidateFilePath(tt.path, tt.operation)
			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecureFileOperations_SecureReadFile(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content for reading"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Test successful read
	result, data, err := sfo.SecureReadFile(testFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, testContent, string(data))
	assert.Equal(t, "read", result.Operation)
	assert.Equal(t, int64(len(testContent)), result.BytesRead)
	assert.NotEmpty(t, result.Checksum)
	assert.NotEmpty(t, result.Permissions)

	// Test reading non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	result, data, err = sfo.SecureReadFile(nonExistentFile)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Nil(t, data)
	assert.Contains(t, result.Error, "failed to stat file")

	// Test reading file that's too large
	sfo.maxFileSize = 5 // Set very small limit
	result, data, err = sfo.SecureReadFile(testFile)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Nil(t, data)
	assert.Contains(t, result.Error, "exceeds maximum allowed size")
}

func TestSecureFileOperations_SecureWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	testFile := filepath.Join(tempDir, "write_test.txt")
	testContent := []byte("test content for writing")

	// Test successful write
	result, err := sfo.SecureWriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "write", result.Operation)
	assert.Equal(t, int64(len(testContent)), result.BytesWritten)
	assert.NotEmpty(t, result.Checksum)
	assert.Equal(t, "0644", result.Permissions)

	// Verify file was written
	writtenContent, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, writtenContent)

	// Test writing with overly permissive permissions
	result, err = sfo.SecureWriteFile(testFile, testContent, 0777)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Warnings[0], "may be too permissive")
	assert.Equal(t, "0644", result.Permissions) // Should be corrected

	// Test writing data that's too large
	sfo.maxFileSize = 5
	result, err = sfo.SecureWriteFile(testFile, testContent, 0644)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "exceeds maximum allowed size")

	// Test writing to invalid path
	invalidFile := "/etc/passwd"
	sfo.maxFileSize = 100 * 1024 * 1024 // Reset size limit
	result, err = sfo.SecureWriteFile(invalidFile, testContent, 0644)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "path validation failed")
}

func TestSecureFileOperations_SecureCopyFile(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Create source file
	srcFile := filepath.Join(tempDir, "source.txt")
	srcContent := []byte("content to copy")
	err := os.WriteFile(srcFile, srcContent, 0644)
	require.NoError(t, err)

	dstFile := filepath.Join(tempDir, "destination.txt")

	// Test successful copy
	result, err := sfo.SecureCopyFile(srcFile, dstFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "copy", result.Operation)
	assert.Contains(t, result.FilePath, srcFile)
	assert.Contains(t, result.FilePath, dstFile)
	assert.Equal(t, int64(len(srcContent)), result.BytesWritten)
	assert.NotEmpty(t, result.Checksum)

	// Verify file was copied
	copiedContent, err := os.ReadFile(dstFile)
	assert.NoError(t, err)
	assert.Equal(t, srcContent, copiedContent)

	// Test copying non-existent source
	nonExistentSrc := filepath.Join(tempDir, "nonexistent.txt")
	result, err = sfo.SecureCopyFile(nonExistentSrc, dstFile)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "failed to stat source file")

	// Test copying file that's too large
	sfo.maxFileSize = 5
	result, err = sfo.SecureCopyFile(srcFile, dstFile)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "exceeds maximum allowed size")

	// Test copying to invalid destination
	sfo.maxFileSize = 100 * 1024 * 1024 // Reset size limit
	invalidDst := "/etc/passwd"
	result, err = sfo.SecureCopyFile(srcFile, invalidDst)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "destination path validation failed")
}

func TestSecureFileOperations_SecureCreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	testDir := filepath.Join(tempDir, "test_directory")

	// Test successful directory creation
	result, err := sfo.SecureCreateDirectory(testDir, 0755)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "mkdir", result.Operation)
	assert.Equal(t, testDir, result.FilePath)
	assert.Equal(t, "0755", result.Permissions)

	// Verify directory was created
	info, err := os.Stat(testDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test creating directory with overly permissive permissions
	testDir2 := filepath.Join(tempDir, "test_directory2")
	result, err = sfo.SecureCreateDirectory(testDir2, 0777)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Warnings[0], "may be too permissive")
	assert.Equal(t, "0755", result.Permissions) // Should be corrected

	// Test creating directory at invalid path
	invalidDir := "/etc/test"
	result, err = sfo.SecureCreateDirectory(invalidDir, 0755)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "path validation failed")
}

func TestSecureFileOperations_SecureRemoveFile(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Create test file
	testFile := filepath.Join(tempDir, "to_remove.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Test successful removal
	result, err := sfo.SecureRemoveFile(testFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "remove", result.Operation)
	assert.Equal(t, testFile, result.FilePath)

	// Verify file was removed
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))

	// Test removing non-existent file
	result, err = sfo.SecureRemoveFile(testFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Warnings[0], "File does not exist")

	// Test removing file at invalid path
	invalidFile := "/etc/passwd"
	result, err = sfo.SecureRemoveFile(invalidFile)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "path validation failed")
}

func TestSecureFileOperations_GetFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Create test file
	testFile := filepath.Join(tempDir, "permissions_test.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	require.NoError(t, err)

	// Test getting file permissions
	perms, err := sfo.GetFilePermissions(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, perms)
	assert.Contains(t, perms, "mode")
	assert.Contains(t, perms, "mode_octal")
	assert.Contains(t, perms, "is_dir")
	assert.Contains(t, perms, "size")
	assert.Contains(t, perms, "mod_time")
	assert.Contains(t, perms, "is_readable")
	assert.Contains(t, perms, "is_writable")
	assert.Contains(t, perms, "is_executable")

	assert.False(t, perms["is_dir"].(bool))
	assert.True(t, perms["is_readable"].(bool))
	assert.Greater(t, perms["size"].(int64), int64(0))

	// Test getting directory permissions
	testDir := filepath.Join(tempDir, "test_dir")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	perms, err = sfo.GetFilePermissions(testDir)
	assert.NoError(t, err)
	assert.True(t, perms["is_dir"].(bool))

	// Test getting permissions for non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	perms, err = sfo.GetFilePermissions(nonExistentFile)
	assert.Error(t, err)
	assert.Nil(t, perms)

	// Test getting permissions for invalid path
	invalidFile := "/etc/passwd"
	perms, err = sfo.GetFilePermissions(invalidFile)
	assert.Error(t, err)
	assert.Nil(t, perms)
}

func TestSecureFileOperations_ValidateFileIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Create test file
	testFile := filepath.Join(tempDir, "integrity_test.txt")
	testContent := []byte("content for integrity check")
	err := os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	// Calculate expected checksum
	_, data, err := sfo.SecureReadFile(testFile)
	require.NoError(t, err)
	expectedChecksum := sfo.calculateChecksum(data)

	// Test successful integrity validation
	result, err := sfo.ValidateFileIntegrity(testFile, expectedChecksum)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "integrity_check", result.Operation)
	assert.Equal(t, expectedChecksum, result.Checksum)

	// Test integrity validation with wrong checksum
	wrongChecksum := "wrong_checksum"
	result, err = sfo.ValidateFileIntegrity(testFile, wrongChecksum)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "integrity check failed")
	assert.Contains(t, result.Warnings[0], "may have been modified")

	// Test integrity validation for non-existent file
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	result, err = sfo.ValidateFileIntegrity(nonExistentFile, expectedChecksum)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "failed to read file")
}

func TestSecureFileOperations_SetSecurityConfig(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	config := map[string]interface{}{
		"max_file_size":            int64(50 * 1024 * 1024), // 50MB
		"allowed_extensions":       []string{".txt", ".md"},
		"blocked_paths":            []string{"/blocked/path"},
		"require_permission_check": false,
	}

	err := sfo.SetSecurityConfig(config)
	assert.NoError(t, err)

	// Verify configuration was applied
	assert.Equal(t, int64(50*1024*1024), sfo.maxFileSize)
	assert.Equal(t, []string{".txt", ".md"}, sfo.allowedExts)
	assert.Equal(t, []string{"/blocked/path"}, sfo.blockedPaths)
	assert.False(t, sfo.requirePermCheck)
}

func TestSecureFileOperations_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Test with empty allowed paths
	sfoEmpty := NewSecureFileOperations([]string{})
	testFile := filepath.Join(tempDir, "test.txt")
	err := sfoEmpty.ValidateFilePath(testFile, "read")
	assert.Error(t, err) // Should fail because no paths are allowed

	// Test with very long path
	longPath := filepath.Join(tempDir, strings.Repeat("a", 1000), "file.txt")
	err = sfo.ValidateFilePath(longPath, "read")
	assert.NoError(t, err) // Should pass validation (path length not restricted)

	// Test with special characters in path
	specialPath := filepath.Join(tempDir, "file with spaces & symbols!@#.txt")
	err = sfo.ValidateFilePath(specialPath, "read")
	assert.NoError(t, err) // Should pass validation

	// Test case sensitivity in blocked paths
	sfo.blockedPaths = []string{"/ETC/PASSWD"}
	err = sfo.ValidateFilePath("/etc/passwd", "read")
	assert.Error(t, err) // Should be blocked (case-insensitive check)
}

func TestSecureFileOperations_Integration(t *testing.T) {
	tempDir := t.TempDir()
	allowedPaths := []string{tempDir}
	sfo := NewSecureFileOperations(allowedPaths)

	// Test complete workflow: create directory, write file, read file, copy file, remove file

	// 1. Create directory
	testDir := filepath.Join(tempDir, "integration_test")
	result, err := sfo.SecureCreateDirectory(testDir, 0755)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 2. Write file
	testFile := filepath.Join(testDir, "test.txt")
	testContent := []byte("integration test content")
	result, err = sfo.SecureWriteFile(testFile, testContent, 0644)
	require.NoError(t, err)
	assert.True(t, result.Success)
	originalChecksum := result.Checksum

	// 3. Read file
	result, data, err := sfo.SecureReadFile(testFile)
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, testContent, data)

	// 4. Validate integrity
	result, err = sfo.ValidateFileIntegrity(testFile, originalChecksum)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 5. Copy file
	copyFile := filepath.Join(testDir, "copy.txt")
	result, err = sfo.SecureCopyFile(testFile, copyFile)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 6. Verify copy has same content
	result, copyData, err := sfo.SecureReadFile(copyFile)
	require.NoError(t, err)
	assert.Equal(t, testContent, copyData)

	// 7. Get file permissions
	perms, err := sfo.GetFilePermissions(testFile)
	require.NoError(t, err)
	assert.False(t, perms["is_dir"].(bool))
	assert.True(t, perms["is_readable"].(bool))

	// 8. Remove files
	result, err = sfo.SecureRemoveFile(testFile)
	require.NoError(t, err)
	assert.True(t, result.Success)

	result, err = sfo.SecureRemoveFile(copyFile)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// Verify files are removed
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(copyFile)
	assert.True(t, os.IsNotExist(err))
}

// Helper method to calculate checksum (would be part of SecureFileOperations in real implementation)
func (sfo *SecureFileOperations) calculateChecksum(data []byte) string {
	// This is a simplified version - the real implementation would use crypto/sha256
	return "mock_checksum_" + string(rune(len(data)))
}
