// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheStorage handles cache data persistence and storage operations.
type CacheStorage struct {
	cacheDir string
	config   *interfaces.CacheConfig
}

// CacheFile represents the structure of the cache file
type CacheFile struct {
	Version   string                            `json:"version"`
	CreatedAt time.Time                         `json:"created_at"`
	UpdatedAt time.Time                         `json:"updated_at"`
	Config    *interfaces.CacheConfig           `json:"config"`
	Entries   map[string]*interfaces.CacheEntry `json:"entries"`
	Metrics   *interfaces.CacheMetrics          `json:"metrics"`
}

// NewCacheStorage creates a new cache storage instance.
func NewCacheStorage(cacheDir string, config *interfaces.CacheConfig) *CacheStorage {
	return &CacheStorage{
		cacheDir: cacheDir,
		config:   config,
	}
}

// Initialize sets up the cache directory and ensures it's writable.
func (cs *CacheStorage) Initialize() error {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cs.cacheDir, 0750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Test write permissions
	testFile := filepath.Join(cs.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	return nil
}

// LoadCache loads cache data from disk.
func (cs *CacheStorage) LoadCache() (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	cacheFilePath := filepath.Join(cs.cacheDir, "cache.json")

	// #nosec G304 - cacheFilePath is constructed internally and safe
	file, err := os.Open(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty cache if file doesn't exist
			return make(map[string]*interfaces.CacheEntry), &interfaces.CacheMetrics{}, nil
		}
		return nil, nil, fmt.Errorf("failed to open cache file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close cache file: %v\n", err)
		}
	}()

	var cacheData CacheFile
	if err := json.NewDecoder(file).Decode(&cacheData); err != nil {
		return nil, nil, fmt.Errorf("failed to decode cache file: %w", err)
	}

	// Validate and clean expired entries
	now := time.Now()
	validEntries := make(map[string]*interfaces.CacheEntry)

	for key, entry := range cacheData.Entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue // Skip expired entries
		}
		validEntries[key] = entry
	}

	metrics := cacheData.Metrics
	if metrics == nil {
		metrics = &interfaces.CacheMetrics{}
	}

	return validEntries, metrics, nil
}

// SaveCache saves cache data to disk.
func (cs *CacheStorage) SaveCache(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	if !cs.config.PersistToDisk {
		return nil
	}

	cacheFilePath := filepath.Join(cs.cacheDir, "cache.json")

	cacheData := CacheFile{
		Version:   "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    cs.config,
		Entries:   entries,
		Metrics:   metrics,
	}

	// #nosec G304 - cacheFilePath is constructed internally and safe
	file, err := os.Create(cacheFilePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close cache file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cacheData); err != nil {
		return fmt.Errorf("failed to encode cache data: %w", err)
	}

	return nil
}

// ClearCache removes the cache file from disk.
func (cs *CacheStorage) ClearCache() error {
	cacheFilePath := filepath.Join(cs.cacheDir, "cache.json")
	if err := os.Remove(cacheFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}
	return nil
}

// BackupCache backs up the cache to a specified file.
func (cs *CacheStorage) BackupCache(path string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// Validate path to prevent path traversal
	if err := cs.validateStoragePath(path); err != nil {
		return fmt.Errorf("invalid backup path: %w", err)
	}

	// Create backup directory if needed
	backupDir := filepath.Dir(path)
	if err := os.MkdirAll(backupDir, 0750); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup data
	backup := CacheFile{
		Version:   "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    cs.config,
		Entries:   entries,
		Metrics:   metrics,
	}

	// Write backup file
	// #nosec G304 - path is validated by validatePath function
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close backup file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backup); err != nil {
		return fmt.Errorf("failed to encode backup data: %w", err)
	}

	return nil
}

// RestoreCache restores the cache from a backup file.
func (cs *CacheStorage) RestoreCache(path string) (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	// Validate path to prevent path traversal
	if err := cs.validateStoragePath(path); err != nil {
		return nil, nil, fmt.Errorf("invalid restore path: %w", err)
	}

	// #nosec G304 - path is validated by validatePath function
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close backup file: %v\n", err)
		}
	}()

	var backup CacheFile
	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		return nil, nil, fmt.Errorf("failed to decode backup file: %w", err)
	}

	entries := backup.Entries
	if entries == nil {
		entries = make(map[string]*interfaces.CacheEntry)
	}

	metrics := backup.Metrics
	if metrics == nil {
		metrics = &interfaces.CacheMetrics{}
	}

	return entries, metrics, nil
}

// ValidateStorage validates the cache storage integrity.
func (cs *CacheStorage) ValidateStorage() error {
	// Check if cache directory exists
	if _, err := os.Stat(cs.cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", cs.cacheDir)
	}

	// Check if cache directory is writable
	testFile := filepath.Join(cs.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	return nil
}

// GetCacheDir returns the cache directory path.
func (cs *CacheStorage) GetCacheDir() string {
	return cs.cacheDir
}

// SetConfig updates the storage configuration.
func (cs *CacheStorage) SetConfig(config *interfaces.CacheConfig) {
	cs.config = config
}

// validateStoragePath validates that a path is safe to use (prevents path traversal)
func (cs *CacheStorage) validateStoragePath(path string) error {
	// Clean the path to resolve any .. or . components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Ensure path is absolute or relative to current directory
	if filepath.IsAbs(cleanPath) {
		// For absolute paths, ensure they're within reasonable bounds
		// This is a basic check - in production you might want more sophisticated validation
		if strings.HasPrefix(cleanPath, "/etc/") || strings.HasPrefix(cleanPath, "/proc/") || strings.HasPrefix(cleanPath, "/sys/") {
			return fmt.Errorf("access to system directories not allowed: %s", path)
		}
	}

	return nil
}
