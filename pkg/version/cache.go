package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
)

// CacheEntry represents a cached version with TTL
type CacheEntry struct {
	Version   string    `json:"version"`
	ExpiresAt time.Time `json:"expires_at"`
}

// FileCache implements the VersionCache interface with file-based persistence
type FileCache struct {
	mu       sync.RWMutex
	cache    map[string]CacheEntry
	filePath string
	ttl      time.Duration
}

// NewFileCache creates a new file-based version cache
func NewFileCache(cacheDir string, ttl time.Duration) (*FileCache, error) {
	if ttl <= 0 {
		ttl = 24 * time.Hour // Default TTL of 24 hours
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	filePath := filepath.Join(cacheDir, "version_cache.json")

	cache := &FileCache{
		cache:    make(map[string]CacheEntry),
		filePath: filePath,
		ttl:      ttl,
	}

	// Load existing cache from file
	if err := cache.load(); err != nil {
		// Log warning but don't fail - start with empty cache
		fmt.Printf("Warning: failed to load cache from %s: %v\n", filePath, err)
	}

	return cache, nil
}

// Get retrieves a cached version by key
func (c *FileCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return "", false
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		// Entry expired, remove it
		delete(c.cache, key)
		// Save updated cache asynchronously
		go func() {
			if err := c.save(); err != nil {
				fmt.Printf("Warning: failed to save cache after expiry cleanup: %v\n", err)
			}
		}()
		return "", false
	}

	return entry.Version, true
}

// Set stores a version in the cache
func (c *FileCache) Set(key, version string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheEntry{
		Version:   version,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	// Save to file asynchronously
	go func() {
		if err := c.save(); err != nil {
			fmt.Printf("Warning: failed to save cache: %v\n", err)
		}
	}()

	return nil
}

// Delete removes a version from the cache
func (c *FileCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)

	// Save to file asynchronously
	go func() {
		if err := c.save(); err != nil {
			fmt.Printf("Warning: failed to save cache: %v\n", err)
		}
	}()

	return nil
}

// Clear removes all cached versions
func (c *FileCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]CacheEntry)

	// Save to file asynchronously
	go func() {
		if err := c.save(); err != nil {
			fmt.Printf("Warning: failed to save cache: %v\n", err)
		}
	}()

	return nil
}

// Keys returns all cached keys
func (c *FileCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.cache))
	now := time.Now()

	for key, entry := range c.cache {
		// Only return non-expired keys
		if now.Before(entry.ExpiresAt) {
			keys = append(keys, key)
		}
	}

	return keys
}

// CleanExpired removes all expired entries from the cache
func (c *FileCache) CleanExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removed := 0

	for key, entry := range c.cache {
		if now.After(entry.ExpiresAt) {
			delete(c.cache, key)
			removed++
		}
	}

	if removed > 0 {
		// Save to file asynchronously
		go func() {
			if err := c.save(); err != nil {
				fmt.Printf("Warning: failed to save cache after cleanup: %v\n", err)
			}
		}()
	}

	return removed
}

// load reads the cache from the file
func (c *FileCache) load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, start with empty cache
			return nil
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	var cache map[string]CacheEntry
	if err := json.Unmarshal(data, &cache); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	c.cache = cache
	return nil
}

// save writes the cache to the file
func (c *FileCache) save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c.cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	// Write to temporary file first, then rename for atomic operation
	tempFile := c.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	if err := os.Rename(tempFile, c.filePath); err != nil {
		return fmt.Errorf("failed to rename cache file: %w", err)
	}

	return nil
}

// MemoryCache implements the VersionCache interface with in-memory storage
type MemoryCache struct {
	mu    sync.RWMutex
	cache map[string]CacheEntry
	ttl   time.Duration
}

// NewMemoryCache creates a new in-memory version cache
func NewMemoryCache(ttl time.Duration) *MemoryCache {
	if ttl <= 0 {
		ttl = 24 * time.Hour // Default TTL of 24 hours
	}

	return &MemoryCache{
		cache: make(map[string]CacheEntry),
		ttl:   ttl,
	}
}

// Get retrieves a cached version by key
func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return "", false
	}

	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		// Entry expired, remove it
		delete(c.cache, key)
		return "", false
	}

	return entry.Version, true
}

// Set stores a version in the cache
func (c *MemoryCache) Set(key, version string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CacheEntry{
		Version:   version,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	return nil
}

// Delete removes a version from the cache
func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
	return nil
}

// Clear removes all cached versions
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]CacheEntry)
	return nil
}

// Keys returns all cached keys
func (c *MemoryCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.cache))
	now := time.Now()

	for key, entry := range c.cache {
		// Only return non-expired keys
		if now.Before(entry.ExpiresAt) {
			keys = append(keys, key)
		}
	}

	return keys
}

// Ensure both implementations satisfy the interface
var _ interfaces.VersionCache = (*FileCache)(nil)
var _ interfaces.VersionCache = (*MemoryCache)(nil)
