package orchestrator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ToolCache manages persistent caching of tool availability and version information
type ToolCache struct {
	cache        map[string]*models.CachedTool
	cacheMu      sync.RWMutex
	cacheFile    string
	ttl          time.Duration
	logger       *logger.Logger
	autoSave     bool
	lastSaved    time.Time
	saveInterval time.Duration
	offlineMode  bool
}

// ToolCacheConfig holds configuration for the tool cache
type ToolCacheConfig struct {
	CacheDir     string        // Directory to store cache file
	TTL          time.Duration // Time-to-live for cache entries
	AutoSave     bool          // Whether to auto-save cache periodically
	SaveInterval time.Duration // Interval for auto-save
}

// DefaultToolCacheConfig returns default cache configuration
func DefaultToolCacheConfig() *ToolCacheConfig {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "project-generator")

	return &ToolCacheConfig{
		CacheDir:     cacheDir,
		TTL:          5 * time.Minute,
		AutoSave:     true,
		SaveInterval: 30 * time.Second,
	}
}

// NewToolCache creates a new tool cache instance
func NewToolCache(config *ToolCacheConfig, log *logger.Logger) (*ToolCache, error) {
	if config == nil {
		config = DefaultToolCacheConfig()
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(config.CacheDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheFile := filepath.Join(config.CacheDir, "tool_cache.json")

	tc := &ToolCache{
		cache:        make(map[string]*models.CachedTool),
		cacheFile:    cacheFile,
		ttl:          config.TTL,
		logger:       log,
		autoSave:     config.AutoSave,
		saveInterval: config.SaveInterval,
		lastSaved:    time.Now(),
		offlineMode:  false,
	}

	// Load existing cache from disk
	if err := tc.Load(); err != nil {
		if log != nil {
			log.Debug(fmt.Sprintf("Could not load cache from disk: %v", err))
		}
		// Not a fatal error, just start with empty cache
	}

	return tc, nil
}

// Get retrieves a cached tool entry if it exists and is not expired
func (tc *ToolCache) Get(toolName string) (*models.CachedTool, bool) {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	cached, exists := tc.cache[toolName]
	if !exists {
		return nil, false
	}

	// Check if cache entry is expired
	if time.Since(cached.CachedAt) > cached.TTL {
		return nil, false
	}

	return cached, true
}

// Set stores a tool entry in the cache
func (tc *ToolCache) Set(toolName string, available bool, version string) {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	tc.cache[toolName] = &models.CachedTool{
		Available: available,
		Version:   version,
		CachedAt:  time.Now(),
		TTL:       tc.ttl,
	}

	// Auto-save if enabled and interval has passed
	if tc.autoSave && time.Since(tc.lastSaved) > tc.saveInterval {
		go func() {
			if err := tc.Save(); err != nil && tc.logger != nil {
				tc.logger.Debug(fmt.Sprintf("Failed to auto-save cache: %v", err))
			}
		}()
	}
}

// Clear removes all entries from the cache
func (tc *ToolCache) Clear() {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	tc.cache = make(map[string]*models.CachedTool)

	if tc.logger != nil {
		tc.logger.Debug("Tool cache cleared")
	}
}

// ClearExpired removes expired entries from the cache
func (tc *ToolCache) ClearExpired() int {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	removed := 0
	for toolName, cached := range tc.cache {
		if time.Since(cached.CachedAt) > cached.TTL {
			delete(tc.cache, toolName)
			removed++
		}
	}

	if tc.logger != nil && removed > 0 {
		tc.logger.Debug(fmt.Sprintf("Removed %d expired cache entries", removed))
	}

	return removed
}

// SetTTL updates the time-to-live for cache entries
func (tc *ToolCache) SetTTL(ttl time.Duration) {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	tc.ttl = ttl
}

// GetTTL returns the current time-to-live setting
func (tc *ToolCache) GetTTL() time.Duration {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	return tc.ttl
}

// Size returns the number of entries in the cache
func (tc *ToolCache) Size() int {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	return len(tc.cache)
}

// Save persists the cache to disk
func (tc *ToolCache) Save() error {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	// Create a serializable version of the cache
	data, err := json.MarshalIndent(tc.cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Write to temporary file first
	tempFile := tc.cacheFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, tc.cacheFile); err != nil {
		// Clean up temp file, log error if cleanup fails
		if removeErr := os.Remove(tempFile); removeErr != nil && tc.logger != nil {
			tc.logger.Debug(fmt.Sprintf("Failed to remove temp file %s: %v", tempFile, removeErr))
		}
		return fmt.Errorf("failed to rename cache file: %w", err)
	}

	tc.lastSaved = time.Now()

	if tc.logger != nil {
		tc.logger.Debug(fmt.Sprintf("Tool cache saved to %s", tc.cacheFile))
	}

	return nil
}

// Load reads the cache from disk
func (tc *ToolCache) Load() error {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	// Check if cache file exists
	if _, err := os.Stat(tc.cacheFile); os.IsNotExist(err) {
		return fmt.Errorf("cache file does not exist")
	}

	// Read cache file
	data, err := os.ReadFile(tc.cacheFile)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	// Unmarshal cache data
	cache := make(map[string]*models.CachedTool)
	if err := json.Unmarshal(data, &cache); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	tc.cache = cache

	if tc.logger != nil {
		tc.logger.Debug(fmt.Sprintf("Tool cache loaded from %s (%d entries)", tc.cacheFile, len(cache)))
	}

	return nil
}

// GetCacheFile returns the path to the cache file
func (tc *ToolCache) GetCacheFile() string {
	return tc.cacheFile
}

// GetStats returns statistics about the cache
func (tc *ToolCache) GetStats() map[string]interface{} {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	expired := 0
	available := 0
	unavailable := 0

	for _, cached := range tc.cache {
		if time.Since(cached.CachedAt) > cached.TTL {
			expired++
		} else if cached.Available {
			available++
		} else {
			unavailable++
		}
	}

	return map[string]interface{}{
		"total":       len(tc.cache),
		"expired":     expired,
		"available":   available,
		"unavailable": unavailable,
		"cache_file":  tc.cacheFile,
		"ttl":         tc.ttl.String(),
		"last_saved":  tc.lastSaved,
	}
}

// SetOfflineMode sets whether the cache is in offline mode
func (tc *ToolCache) SetOfflineMode(offline bool) {
	tc.cacheMu.Lock()
	defer tc.cacheMu.Unlock()

	tc.offlineMode = offline

	if tc.logger != nil {
		if offline {
			tc.logger.Debug("Cache set to offline mode - will not expire entries")
		} else {
			tc.logger.Debug("Cache set to online mode")
		}
	}
}

// IsOfflineMode returns whether the cache is in offline mode
func (tc *ToolCache) IsOfflineMode() bool {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	return tc.offlineMode
}

// GetWithOfflineSupport retrieves a cached tool entry, respecting offline mode
func (tc *ToolCache) GetWithOfflineSupport(toolName string) (*models.CachedTool, bool) {
	tc.cacheMu.RLock()
	defer tc.cacheMu.RUnlock()

	cached, exists := tc.cache[toolName]
	if !exists {
		return nil, false
	}

	// In offline mode, don't check expiration
	if tc.offlineMode {
		if tc.logger != nil {
			tc.logger.Debug(fmt.Sprintf("Tool '%s' retrieved from cache (offline mode)", toolName))
		}
		return cached, true
	}

	// Check if cache entry is expired in online mode
	if time.Since(cached.CachedAt) > cached.TTL {
		return nil, false
	}

	return cached, true
}
