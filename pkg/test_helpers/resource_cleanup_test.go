package test_helpers

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewResourceCleanup(t *testing.T) {
	cleanup := NewResourceCleanup()

	if cleanup == nil {
		t.Fatal("NewResourceCleanup() returned nil")
	}

	if cleanup.tempDirs == nil {
		t.Error("tempDirs slice should be initialized")
	}

	if cleanup.tempFiles == nil {
		t.Error("tempFiles slice should be initialized")
	}

	if cleanup.cleanupFunc == nil {
		t.Error("cleanupFunc slice should be initialized")
	}
}

func TestRegisterTempDir(t *testing.T) {
	cleanup := NewResourceCleanup()
	testDir := "/tmp/test-dir"

	cleanup.RegisterTempDir(testDir)

	if len(cleanup.tempDirs) != 1 {
		t.Errorf("Expected 1 temp dir, got %d", len(cleanup.tempDirs))
	}

	if cleanup.tempDirs[0] != testDir {
		t.Errorf("Expected %s, got %s", testDir, cleanup.tempDirs[0])
	}
}

func TestRegisterTempFile(t *testing.T) {
	cleanup := NewResourceCleanup()
	testFile := "/tmp/test-file.txt"

	cleanup.RegisterTempFile(testFile)

	if len(cleanup.tempFiles) != 1 {
		t.Errorf("Expected 1 temp file, got %d", len(cleanup.tempFiles))
	}

	if cleanup.tempFiles[0] != testFile {
		t.Errorf("Expected %s, got %s", testFile, cleanup.tempFiles[0])
	}
}

func TestRegisterCleanupFunc(t *testing.T) {
	cleanup := NewResourceCleanup()
	called := false

	cleanupFunc := func() error {
		called = true
		return nil
	}

	cleanup.RegisterCleanupFunc(cleanupFunc)

	if len(cleanup.cleanupFunc) != 1 {
		t.Errorf("Expected 1 cleanup function, got %d", len(cleanup.cleanupFunc))
	}

	// Test that the function is actually called during cleanup
	cleanup.Cleanup(t)

	if !called {
		t.Error("Cleanup function was not called")
	}
}

func TestCleanupTempFiles(t *testing.T) {
	cleanup := NewResourceCleanup()

	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-file.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Register the file for cleanup
	cleanup.RegisterTempFile(testFile)

	// Verify file exists before cleanup
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Test file should exist before cleanup")
	}

	// Perform cleanup
	cleanup.Cleanup(t)

	// Verify file is removed after cleanup
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Test file should be removed after cleanup")
	}
}

func TestCleanupTempDirs(t *testing.T) {
	cleanup := NewResourceCleanup()

	// Create a temporary directory with content
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test-subdir")

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a file inside the directory
	testFile := filepath.Join(testDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file in directory: %v", err)
	}

	// Register the directory for cleanup
	cleanup.RegisterTempDir(testDir)

	// Verify directory exists before cleanup
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatal("Test directory should exist before cleanup")
	}

	// Perform cleanup
	cleanup.Cleanup(t)

	// Verify directory is removed after cleanup
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Error("Test directory should be removed after cleanup")
	}
}

func TestCleanupWithErrors(t *testing.T) {
	cleanup := NewResourceCleanup()

	// Register non-existent files and directories (should not cause failures)
	cleanup.RegisterTempFile("/non/existent/file.txt")
	cleanup.RegisterTempDir("/non/existent/directory")

	// Register a cleanup function that returns an error
	cleanup.RegisterCleanupFunc(func() error {
		return os.ErrNotExist
	})

	// Cleanup should not panic or fail the test
	cleanup.Cleanup(t)
}

func TestWithCleanup(t *testing.T) {
	var cleanupInstance *ResourceCleanup
	testExecuted := false

	WithCleanup(t, func(cleanup *ResourceCleanup) {
		cleanupInstance = cleanup
		testExecuted = true

		// Create a temporary file that should be cleaned up
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "with-cleanup-test.txt")

		err := os.WriteFile(testFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		cleanup.RegisterTempFile(testFile)
	})

	if !testExecuted {
		t.Error("Test function was not executed")
	}

	if cleanupInstance == nil {
		t.Error("Cleanup instance should not be nil")
	}
}

func TestMemoryCleanup(t *testing.T) {
	// This test mainly ensures MemoryCleanup doesn't panic
	// and completes within a reasonable time
	start := time.Now()

	MemoryCleanup()

	duration := time.Since(start)
	if duration > time.Second {
		t.Errorf("MemoryCleanup took too long: %v", duration)
	}
}

func TestEnsureNoResourceLeaks(t *testing.T) {
	// Test that EnsureNoResourceLeaks executes the test function
	testExecuted := false

	EnsureNoResourceLeaks(t, func() {
		testExecuted = true

		// Allocate some memory but not enough to trigger the leak warning
		data := make([]byte, 1024)
		_ = data
	})

	if !testExecuted {
		t.Error("Test function was not executed")
	}
}

func TestEnsureNoResourceLeaksWithLargeMem(t *testing.T) {
	// Test with a function that allocates significant memory
	// This should trigger a warning log but not fail the test
	EnsureNoResourceLeaks(t, func() {
		// Allocate a large amount of memory
		data := make([][]byte, 1000)
		for i := range data {
			data[i] = make([]byte, 100*1024) // 100KB each
		}
		// Keep reference to prevent immediate GC
		_ = data
	})
}

func TestConcurrentAccess(t *testing.T) {
	cleanup := NewResourceCleanup()

	// Test concurrent access to the cleanup instance
	done := make(chan bool, 3)

	// Goroutine 1: Register temp dirs
	go func() {
		for i := 0; i < 10; i++ {
			cleanup.RegisterTempDir("/tmp/test-" + string(rune(i)))
		}
		done <- true
	}()

	// Goroutine 2: Register temp files
	go func() {
		for i := 0; i < 10; i++ {
			cleanup.RegisterTempFile("/tmp/file-" + string(rune(i)))
		}
		done <- true
	}()

	// Goroutine 3: Register cleanup functions
	go func() {
		for i := 0; i < 10; i++ {
			cleanup.RegisterCleanupFunc(func() error { return nil })
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify all registrations were successful
	if len(cleanup.tempDirs) != 10 {
		t.Errorf("Expected 10 temp dirs, got %d", len(cleanup.tempDirs))
	}

	if len(cleanup.tempFiles) != 10 {
		t.Errorf("Expected 10 temp files, got %d", len(cleanup.tempFiles))
	}

	if len(cleanup.cleanupFunc) != 10 {
		t.Errorf("Expected 10 cleanup functions, got %d", len(cleanup.cleanupFunc))
	}
}

func TestMultipleCleanupCalls(t *testing.T) {
	cleanup := NewResourceCleanup()

	// Create and register a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "multiple-cleanup-test.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cleanup.RegisterTempFile(testFile)

	// Call cleanup multiple times - should not panic or cause errors
	cleanup.Cleanup(t)
	cleanup.Cleanup(t) // Second call should handle non-existent files gracefully
	cleanup.Cleanup(t) // Third call for good measure
}
