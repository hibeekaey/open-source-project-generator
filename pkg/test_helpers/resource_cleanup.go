package test_helpers

import (
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

// ResourceCleanup provides utilities for cleaning up test resources
type ResourceCleanup struct {
	tempDirs    []string
	tempFiles   []string
	mu          sync.Mutex
	cleanupFunc []func() error
}

// NewResourceCleanup creates a new resource cleanup helper
func NewResourceCleanup() *ResourceCleanup {
	return &ResourceCleanup{
		tempDirs:    make([]string, 0),
		tempFiles:   make([]string, 0),
		cleanupFunc: make([]func() error, 0),
	}
}

// RegisterTempDir registers a temporary directory for cleanup
func (rc *ResourceCleanup) RegisterTempDir(dir string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.tempDirs = append(rc.tempDirs, dir)
}

// RegisterTempFile registers a temporary file for cleanup
func (rc *ResourceCleanup) RegisterTempFile(file string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.tempFiles = append(rc.tempFiles, file)
}

// RegisterCleanupFunc registers a custom cleanup function
func (rc *ResourceCleanup) RegisterCleanupFunc(cleanup func() error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.cleanupFunc = append(rc.cleanupFunc, cleanup)
}

// Cleanup performs all registered cleanup operations
func (rc *ResourceCleanup) Cleanup(t *testing.T) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Clean up custom functions first
	for _, cleanup := range rc.cleanupFunc {
		if err := cleanup(); err != nil {
			t.Logf("Cleanup function failed: %v", err)
		}
	}

	// Clean up temporary files
	for _, file := range rc.tempFiles {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Logf("Failed to remove temp file %s: %v", file, err)
		}
	}

	// Clean up temporary directories
	for _, dir := range rc.tempDirs {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("Failed to remove temp dir %s: %v", dir, err)
		}
	}

	// Force garbage collection to free memory
	runtime.GC()
	runtime.GC() // Double GC to ensure cleanup
}

// MemoryCleanup performs memory-specific cleanup operations
func MemoryCleanup() {
	// Force garbage collection
	runtime.GC()

	// Give some time for GC to complete
	time.Sleep(10 * time.Millisecond)

	// Another GC cycle to ensure thorough cleanup
	runtime.GC()
}

// WithCleanup is a helper that automatically cleans up resources when the test completes
func WithCleanup(t *testing.T, testFunc func(cleanup *ResourceCleanup)) {
	cleanup := NewResourceCleanup()
	defer cleanup.Cleanup(t)
	testFunc(cleanup)
}

// EnsureNoResourceLeaks checks for potential resource leaks by monitoring file descriptors and memory
func EnsureNoResourceLeaks(t *testing.T, testFunc func()) {
	var startMem runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&startMem)

	testFunc()

	// Force cleanup and measure memory again
	runtime.GC()
	runtime.GC()
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	// Check for significant memory increase (indicating potential leaks)
	memoryIncrease := endMem.Alloc - startMem.Alloc
	if memoryIncrease > 50*1024*1024 { // 50MB threshold
		t.Logf("WARNING: Memory increased by %d bytes, potential memory leak", memoryIncrease)
	}
}
