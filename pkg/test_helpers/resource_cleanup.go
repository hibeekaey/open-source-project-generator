package test_helpers

import (
	"os"
	"runtime"
	"sync"
	"testing"
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

	// Execute custom cleanup functions
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

	// Force garbage collection
	runtime.GC()
	runtime.GC()
}

// MemoryCleanup performs memory-specific cleanup operations
func MemoryCleanup() {
	runtime.GC()
}
