package security

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewSecureFileOperations(t *testing.T) {
	sfo := NewSecureFileOperations()

	if sfo == nil {
		t.Fatal("NewSecureFileOperations returned nil")
	}

	if sfo.TempFileRandomLength != 16 {
		t.Errorf("Expected TempFileRandomLength to be 16, got %d", sfo.TempFileRandomLength)
	}

	if !sfo.EnablePathValidation {
		t.Error("Expected EnablePathValidation to be true")
	}

	if len(sfo.AllowedTempDirs) == 0 {
		t.Error("Expected at least one allowed temp directory")
	}
}

func TestNewSecureFileOperationsWithConfig(t *testing.T) {
	secureRandom := NewSecureRandom()
	allowedDirs := []string{"/tmp", "/var/tmp"}

	sfo := NewSecureFileOperationsWithConfig(secureRandom, 32, allowedDirs, false)

	if sfo.TempFileRandomLength != 32 {
		t.Errorf("Expected TempFileRandomLength to be 32, got %d", sfo.TempFileRandomLength)
	}

	if sfo.EnablePathValidation {
		t.Error("Expected EnablePathValidation to be false")
	}

	if len(sfo.AllowedTempDirs) != 2 {
		t.Errorf("Expected 2 allowed temp directories, got %d", len(sfo.AllowedTempDirs))
	}
}

func TestCreateSecureTempFile(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()

	// Test creating temp file in specific directory
	tempFile, err := sfo.CreateSecureTempFile(tempDir, "test.")
	if err != nil {
		t.Fatalf("Failed to create secure temp file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Verify file exists and is in correct directory
	if !strings.HasPrefix(tempFile.Name(), filepath.Join(tempDir, "test.")) {
		t.Errorf("Temp file not created with expected prefix: %s", tempFile.Name())
	}

	// Verify file has secure permissions (0600)
	fileInfo, err := tempFile.Stat()
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}

	expectedPerm := os.FileMode(0600)
	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, fileInfo.Mode().Perm())
	}

	// Test that multiple temp files have different names
	tempFile2, err := sfo.CreateSecureTempFile(tempDir, "test.")
	if err != nil {
		t.Fatalf("Failed to create second secure temp file: %v", err)
	}
	defer tempFile2.Close()
	defer os.Remove(tempFile2.Name())

	if tempFile.Name() == tempFile2.Name() {
		t.Error("Two temp files have the same name - not secure!")
	}
}

func TestCreateSecureTempFileDefaultDir(t *testing.T) {
	sfo := NewSecureFileOperations()

	// Test creating temp file with default directory
	tempFile, err := sfo.CreateSecureTempFile("", "test.")
	if err != nil {
		t.Fatalf("Failed to create secure temp file with default dir: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Verify file is in allowed temp directory
	tempDir := filepath.Dir(tempFile.Name())
	found := false
	for _, allowedDir := range sfo.AllowedTempDirs {
		// Get absolute paths for comparison
		absAllowedDir, err := filepath.Abs(allowedDir)
		if err != nil {
			continue
		}
		absTempDir, err := filepath.Abs(tempDir)
		if err != nil {
			continue
		}

		if strings.HasPrefix(absTempDir, absAllowedDir) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Temp file created outside allowed directories: %s (allowed: %v)", tempFile.Name(), sfo.AllowedTempDirs)
	}
}

func TestWriteFileAtomic(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "test_atomic.txt")
	testData := []byte("Hello, atomic world!")

	// Test atomic write
	err := sfo.WriteFileAtomic(targetFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to write file atomically: %v", err)
	}

	// Verify file exists and has correct content
	readData, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if !bytes.Equal(testData, readData) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testData, readData)
	}

	// Verify file has correct permissions
	fileInfo, err := os.Stat(targetFile)
	if err != nil {
		t.Fatalf("Failed to stat written file: %v", err)
	}

	expectedPerm := os.FileMode(0644)
	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, fileInfo.Mode().Perm())
	}
}

func TestWriteFileAtomicOverwrite(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "test_overwrite.txt")

	// Write initial content
	initialData := []byte("Initial content")
	err := sfo.WriteFileAtomic(targetFile, initialData, 0644)
	if err != nil {
		t.Fatalf("Failed to write initial file: %v", err)
	}

	// Overwrite with new content
	newData := []byte("New content that is longer than the initial content")
	err = sfo.WriteFileAtomic(targetFile, newData, 0644)
	if err != nil {
		t.Fatalf("Failed to overwrite file atomically: %v", err)
	}

	// Verify file has new content
	readData, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}

	if !bytes.Equal(newData, readData) {
		t.Errorf("File content mismatch after overwrite. Expected: %s, Got: %s", newData, readData)
	}
}

func TestValidatePath(t *testing.T) {
	sfo := NewSecureFileOperations()

	tests := []struct {
		name        string
		path        string
		allowedDirs []string
		expectError bool
	}{
		{
			name:        "Valid relative path",
			path:        "test/file.txt",
			allowedDirs: nil,
			expectError: false,
		},
		{
			name:        "Directory traversal attempt",
			path:        "../../../etc/passwd",
			allowedDirs: nil,
			expectError: true,
		},
		{
			name:        "Path with .. in middle",
			path:        "test/../../../etc/passwd",
			allowedDirs: nil,
			expectError: true,
		},
		{
			name:        "Empty path",
			path:        "",
			allowedDirs: nil,
			expectError: true,
		},
		{
			name:        "Dangerous system path",
			path:        "/etc/passwd",
			allowedDirs: nil,
			expectError: true,
		},
		{
			name:        "Path within allowed directory",
			path:        "/tmp/test/file.txt",
			allowedDirs: []string{"/tmp"},
			expectError: false,
		},
		{
			name:        "Path outside allowed directory",
			path:        "/home/user/file.txt",
			allowedDirs: []string{"/tmp"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sfo.ValidatePath(tt.path, tt.allowedDirs)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for path %s, but got none", tt.path)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for path %s: %v", tt.path, err)
			}
		})
	}
}

func TestValidatePathDisabled(t *testing.T) {
	// Test with path validation disabled
	sfo := NewSecureFileOperationsWithConfig(nil, 16, nil, false)

	// This should not validate paths even if they contain traversal
	err := sfo.WriteFileAtomic("../test.txt", []byte("test"), 0644)
	// We expect this to fail for other reasons (like directory creation), not path validation
	if err != nil && strings.Contains(err.Error(), "path validation failed") {
		t.Error("Path validation should be disabled but validation error occurred")
	}
}

func TestSecureDelete(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_delete.txt")
	testData := []byte("Sensitive data that should be securely deleted")

	// Create test file
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Test file should exist before deletion")
	}

	// Securely delete the file
	err = sfo.SecureDelete(testFile)
	if err != nil {
		t.Fatalf("Failed to securely delete file: %v", err)
	}

	// Verify file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("File should not exist after secure deletion")
	}
}

func TestSecureDeleteNonExistentFile(t *testing.T) {
	sfo := NewSecureFileOperations()
	nonExistentFile := "/tmp/non_existent_file_12345.txt"

	// Should not error when trying to delete non-existent file
	err := sfo.SecureDelete(nonExistentFile)
	if err != nil {
		t.Errorf("Secure delete of non-existent file should not error: %v", err)
	}
}

func TestSecureDeleteEmptyFile(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty_file.txt")

	// Create empty test file
	err := os.WriteFile(testFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty test file: %v", err)
	}

	// Securely delete the empty file
	err = sfo.SecureDelete(testFile)
	if err != nil {
		t.Fatalf("Failed to securely delete empty file: %v", err)
	}

	// Verify file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Empty file should not exist after secure deletion")
	}
}

func TestEnsureSecurePermissions(t *testing.T) {
	sfo := NewSecureFileOperations()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_permissions.txt")

	// Create test file with default permissions
	err := os.WriteFile(testFile, []byte("test"), 0666)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set secure permissions
	expectedPerm := os.FileMode(0600)
	err = sfo.EnsureSecurePermissions(testFile, expectedPerm)
	if err != nil {
		t.Fatalf("Failed to set secure permissions: %v", err)
	}

	// Verify permissions were set correctly
	fileInfo, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file after permission change: %v", err)
	}

	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Expected permissions %v, got %v", expectedPerm, fileInfo.Mode().Perm())
	}
}

func TestEnsureSecurePermissionsNonExistentFile(t *testing.T) {
	sfo := NewSecureFileOperations()
	nonExistentFile := "/tmp/non_existent_file_permissions.txt"

	// Should error when trying to set permissions on non-existent file
	err := sfo.EnsureSecurePermissions(nonExistentFile, 0600)
	if err == nil {
		t.Error("Setting permissions on non-existent file should error")
	}
}

// Test convenience functions
func TestConvenienceFunctions(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "convenience_test.txt")
	testData := []byte("Testing convenience functions")

	// Test WriteFileAtomic convenience function
	err := WriteFileAtomic(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("Convenience WriteFileAtomic failed: %v", err)
	}

	// Verify file was written correctly
	readData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file written by convenience function: %v", err)
	}

	if !bytes.Equal(testData, readData) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", testData, readData)
	}

	// Test CreateSecureTempFile convenience function
	tempFile, err := CreateSecureTempFile(tempDir, "convenience.")
	if err != nil {
		t.Fatalf("Convenience CreateSecureTempFile failed: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Test ValidatePath convenience function
	err = ValidatePath(testFile, nil)
	if err != nil {
		t.Errorf("Convenience ValidatePath failed: %v", err)
	}

	// Test EnsureSecurePermissions convenience function
	err = EnsureSecurePermissions(testFile, 0600)
	if err != nil {
		t.Fatalf("Convenience EnsureSecurePermissions failed: %v", err)
	}

	// Test SecureDelete convenience function
	err = SecureDelete(testFile)
	if err != nil {
		t.Fatalf("Convenience SecureDelete failed: %v", err)
	}
}

// Benchmark tests
func BenchmarkWriteFileAtomic(b *testing.B) {
	sfo := NewSecureFileOperations()
	tempDir := b.TempDir()
	testData := bytes.Repeat([]byte("benchmark test data "), 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tempDir, fmt.Sprintf("benchmark_%d.txt", i))
		err := sfo.WriteFileAtomic(testFile, testData, 0644)
		if err != nil {
			b.Fatalf("WriteFileAtomic failed: %v", err)
		}
		os.Remove(testFile) // Clean up
	}
}

func BenchmarkCreateSecureTempFile(b *testing.B) {
	sfo := NewSecureFileOperations()
	tempDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tempFile, err := sfo.CreateSecureTempFile(tempDir, "benchmark.")
		if err != nil {
			b.Fatalf("CreateSecureTempFile failed: %v", err)
		}
		tempFile.Close()
		os.Remove(tempFile.Name())
	}
}

func BenchmarkValidatePath(b *testing.B) {
	sfo := NewSecureFileOperations()
	testPath := "/tmp/test/file.txt"
	allowedDirs := []string{"/tmp"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sfo.ValidatePath(testPath, allowedDirs)
		if err != nil {
			b.Fatalf("ValidatePath failed: %v", err)
		}
	}
}
