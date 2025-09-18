// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Manager implements the CacheManager interface for cache operations.
type Manager struct {
	cacheDir    string
	config      *interfaces.CacheConfig
	entries     map[string]*interfaces.CacheEntry
	metrics     *interfaces.CacheMetrics
	offlineMode bool
	mutex       sync.RWMutex
	callbacks   *cacheCallbacks
	lastCleanup time.Time
	initialized bool
}

// cacheCallbacks holds callback functions for cache events
type cacheCallbacks struct {
	onHit      func(key string)
	onMiss     func(key string)
	onEviction func(key string, reason string)
}

// cacheFile represents the structure of the cache file
type cacheFile struct {
	Version   string                            `json:"version"`
	CreatedAt time.Time                         `json:"created_at"`
	UpdatedAt time.Time                         `json:"updated_at"`
	Config    *interfaces.CacheConfig           `json:"config"`
	Entries   map[string]*interfaces.CacheEntry `json:"entries"`
	Metrics   *interfaces.CacheMetrics          `json:"metrics"`
}

// NewManager creates a new cache manager instance.
func NewManager(cacheDir string) interfaces.CacheManager {
	manager := &Manager{
		cacheDir:    cacheDir,
		config:      interfaces.DefaultCacheConfig(),
		entries:     make(map[string]*interfaces.CacheEntry),
		metrics:     &interfaces.CacheMetrics{},
		offlineMode: false,
		callbacks:   &cacheCallbacks{},
		lastCleanup: time.Now(),
		initialized: false,
	}

	// Set the cache location in config
	manager.config.Location = cacheDir

	// Initialize the cache
	if err := manager.initialize(); err != nil {
		// Log error but don't fail - cache will work in memory only
		fmt.Printf("Warning: Failed to initialize cache: %v\n", err)
	}

	return manager
}

// initialize sets up the cache directory and loads existing cache data
func (m *Manager) initialize() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(m.cacheDir, 0750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Load existing cache data
	if err := m.loadCache(); err != nil {
		// If loading fails, start with empty cache
		m.entries = make(map[string]*interfaces.CacheEntry)
		m.metrics = &interfaces.CacheMetrics{
			CurrentSize:    0,
			MaxSize:        m.config.MaxSize,
			CurrentEntries: 0,
			MaxEntries:     m.config.MaxEntries,
		}
	}

	m.initialized = true
	return nil
}

// loadCache loads cache data from disk
func (m *Manager) loadCache() error {
	cacheFilePath := filepath.Join(m.cacheDir, "cache.json")

	// #nosec G304 - cacheFilePath is constructed internally and safe
	file, err := os.Open(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing cache file
		}
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close cache file: %v\n", err)
		}
	}()

	var cacheData cacheFile
	if err := json.NewDecoder(file).Decode(&cacheData); err != nil {
		return fmt.Errorf("failed to decode cache file: %w", err)
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

	m.entries = validEntries
	if cacheData.Metrics != nil {
		m.metrics = cacheData.Metrics
	}
	if cacheData.Config != nil {
		m.config = cacheData.Config
	}

	return nil
}

// saveCache saves cache data to disk
func (m *Manager) saveCache() error {
	if !m.config.PersistToDisk {
		return nil
	}

	cacheFilePath := filepath.Join(m.cacheDir, "cache.json")

	cacheData := cacheFile{
		Version:   "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    m.config,
		Entries:   m.entries,
		Metrics:   m.metrics,
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

// Get retrieves a value from the cache
func (m *Manager) Get(key string) (any, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, exists := m.entries[key]
	if !exists {
		m.metrics.Misses++
		m.metrics.Gets++
		if m.callbacks.onMiss != nil {
			m.callbacks.onMiss(key)
		}
		return nil, fmt.Errorf("key not found: %s", key)
	}

	// Check if entry is expired
	now := time.Now()
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
		delete(m.entries, key)
		m.metrics.Misses++
		m.metrics.Gets++
		if m.callbacks.onMiss != nil {
			m.callbacks.onMiss(key)
		}
		return nil, fmt.Errorf("key expired: %s", key)
	}

	// Update access information
	entry.AccessedAt = now
	entry.AccessCount++

	m.metrics.Hits++
	m.metrics.Gets++
	if m.callbacks.onHit != nil {
		m.callbacks.onHit(key)
	}

	return entry.Value, nil
}

// Set stores a value in the cache with TTL
func (m *Manager) Set(key string, value any, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()

	// Calculate entry size (approximate)
	size := m.calculateSize(value)

	// Check if we need to evict entries
	if err := m.evictIfNeeded(size); err != nil {
		return fmt.Errorf("failed to evict entries: %w", err)
	}

	// Create cache entry
	entry := &interfaces.CacheEntry{
		Key:         key,
		Value:       value,
		Size:        size,
		CreatedAt:   now,
		UpdatedAt:   now,
		AccessedAt:  now,
		TTL:         ttl,
		AccessCount: 0,
		Compressed:  false,
		Metadata:    make(map[string]any),
	}

	// Set expiration time if TTL is provided
	if ttl > 0 {
		expiresAt := now.Add(ttl)
		entry.ExpiresAt = &expiresAt
	} else if m.config.DefaultTTL > 0 {
		expiresAt := now.Add(m.config.DefaultTTL)
		entry.ExpiresAt = &expiresAt
		entry.TTL = m.config.DefaultTTL
	}

	// Compress if enabled and value is large enough
	if m.config.EnableCompression && size > 1024 {
		if err := m.compressEntry(entry); err != nil {
			// Log warning but continue without compression
			fmt.Printf("Warning: Failed to compress cache entry: %v\n", err)
		}
	}

	// Store the entry
	m.entries[key] = entry

	// Update metrics
	m.metrics.Sets++
	m.metrics.CurrentEntries = len(m.entries)
	m.metrics.CurrentSize += size

	// Save to disk if persistence is enabled
	if m.config.PersistToDisk {
		go func() {
			if err := m.saveCache(); err != nil {
				fmt.Printf("Warning: Failed to save cache: %v\n", err)
			}
		}()
	}

	return nil
}

// calculateSize estimates the size of a value in bytes
func (m *Manager) calculateSize(value any) int64 {
	// This is a rough estimation
	switch v := value.(type) {
	case string:
		return int64(len(v))
	case []byte:
		return int64(len(v))
	case int, int32, int64, float32, float64:
		return 8
	case bool:
		return 1
	default:
		// For complex types, use JSON marshaling to estimate size
		if data, err := json.Marshal(v); err == nil {
			return int64(len(data))
		}
		return 100 // Default estimate
	}
}

// evictIfNeeded evicts entries if cache limits are exceeded
func (m *Manager) evictIfNeeded(newEntrySize int64) error {
	// Check size limit
	if m.config.MaxSize > 0 && m.metrics.CurrentSize+newEntrySize > m.config.MaxSize {
		if err := m.evictBySize(newEntrySize); err != nil {
			return err
		}
	}

	// Check entry count limit
	if m.config.MaxEntries > 0 && len(m.entries) >= m.config.MaxEntries {
		if err := m.evictByCount(); err != nil {
			return err
		}
	}

	return nil
}

// evictBySize evicts entries to make room for new entry
func (m *Manager) evictBySize(newEntrySize int64) error {
	targetSize := int64(float64(m.config.MaxSize) * (1.0 - m.config.EvictionRatio))

	// Get entries sorted by eviction policy
	entries := m.getSortedEntriesForEviction()

	for _, entry := range entries {
		if m.metrics.CurrentSize+newEntrySize <= targetSize {
			break
		}

		delete(m.entries, entry.Key)
		m.metrics.CurrentSize -= entry.Size
		m.metrics.Evictions++

		if m.callbacks.onEviction != nil {
			m.callbacks.onEviction(entry.Key, interfaces.EvictionReasonSize)
		}
	}

	return nil
}

// evictByCount evicts entries to make room for new entry
func (m *Manager) evictByCount() error {
	targetCount := int(float64(m.config.MaxEntries) * (1.0 - m.config.EvictionRatio))

	// Get entries sorted by eviction policy
	entries := m.getSortedEntriesForEviction()

	for _, entry := range entries {
		if len(m.entries) <= targetCount {
			break
		}

		delete(m.entries, entry.Key)
		m.metrics.CurrentSize -= entry.Size
		m.metrics.Evictions++

		if m.callbacks.onEviction != nil {
			m.callbacks.onEviction(entry.Key, interfaces.EvictionReasonCapacity)
		}
	}

	return nil
}

// getSortedEntriesForEviction returns entries sorted by eviction policy
func (m *Manager) getSortedEntriesForEviction() []*interfaces.CacheEntry {
	entries := make([]*interfaces.CacheEntry, 0, len(m.entries))
	for _, entry := range m.entries {
		entries = append(entries, entry)
	}

	switch m.config.EvictionPolicy {
	case "lru": // Least Recently Used
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].AccessedAt.Before(entries[j].AccessedAt)
		})
	case "lfu": // Least Frequently Used
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].AccessCount < entries[j].AccessCount
		})
	case "fifo": // First In, First Out
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].CreatedAt.Before(entries[j].CreatedAt)
		})
	case "ttl": // Shortest TTL first
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].ExpiresAt == nil && entries[j].ExpiresAt == nil {
				return entries[i].CreatedAt.Before(entries[j].CreatedAt)
			}
			if entries[i].ExpiresAt == nil {
				return false
			}
			if entries[j].ExpiresAt == nil {
				return true
			}
			return entries[i].ExpiresAt.Before(*entries[j].ExpiresAt)
		})
	default: // Default to LRU
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].AccessedAt.Before(entries[j].AccessedAt)
		})
	}

	return entries
}

// compressEntry compresses the value in a cache entry
func (m *Manager) compressEntry(entry *interfaces.CacheEntry) error {
	// Convert value to bytes
	var data []byte
	var err error

	switch v := entry.Value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value for compression: %w", err)
		}
	}

	// Compress the data
	var compressed []byte
	switch m.config.CompressionType {
	case "gzip":
		compressed, err = m.compressGzip(data)
	default:
		compressed, err = m.compressGzip(data)
	}

	if err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}

	// Only use compression if it actually reduces size
	if len(compressed) < len(data) {
		entry.Value = compressed
		entry.Size = int64(len(compressed))
		entry.Compressed = true
		entry.Metadata["original_size"] = len(data)
		entry.Metadata["compression_type"] = m.config.CompressionType
	}

	return nil
}

// compressGzip compresses data using gzip
func (m *Manager) compressGzip(data []byte) ([]byte, error) {
	var buf strings.Builder
	writer, err := gzip.NewWriterLevel(&buf, m.config.CompressionLevel)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

// Delete removes a value from the cache
func (m *Manager) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry, exists := m.entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	delete(m.entries, key)
	m.metrics.CurrentSize -= entry.Size
	m.metrics.CurrentEntries = len(m.entries)
	m.metrics.Deletes++

	if m.callbacks.onEviction != nil {
		m.callbacks.onEviction(key, interfaces.EvictionReasonManual)
	}

	// Save to disk if persistence is enabled
	if m.config.PersistToDisk {
		go func() {
			if err := m.saveCache(); err != nil {
				fmt.Printf("Warning: Failed to save cache: %v\n", err)
			}
		}()
	}

	return nil
}

// Clear removes all values from the cache
func (m *Manager) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear all entries
	m.entries = make(map[string]*interfaces.CacheEntry)
	m.metrics.CurrentSize = 0
	m.metrics.CurrentEntries = 0

	// Remove cache file if it exists
	cacheFilePath := filepath.Join(m.cacheDir, "cache.json")
	if err := os.Remove(cacheFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}

	return nil
}

// Clean removes expired entries from the cache
func (m *Manager) Clean() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	// Find expired entries
	for key, entry := range m.entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired entries
	for _, key := range expiredKeys {
		entry := m.entries[key]
		delete(m.entries, key)
		m.metrics.CurrentSize -= entry.Size
		m.metrics.Evictions++

		if m.callbacks.onEviction != nil {
			m.callbacks.onEviction(key, interfaces.EvictionReasonTTL)
		}
	}

	m.metrics.CurrentEntries = len(m.entries)
	m.lastCleanup = now

	// Save to disk if persistence is enabled
	if m.config.PersistToDisk {
		go func() {
			if err := m.saveCache(); err != nil {
				fmt.Printf("Warning: Failed to save cache: %v\n", err)
			}
		}()
	}

	return nil
}

// GetStats returns cache statistics
func (m *Manager) GetStats() (*interfaces.CacheStats, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Count expired entries
	now := time.Now()
	expiredCount := 0
	for _, entry := range m.entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredCount++
		}
	}

	// Calculate hit rate
	hitRate := 0.0
	if m.metrics.Gets > 0 {
		hitRate = float64(m.metrics.Hits) / float64(m.metrics.Gets)
	}

	// Determine cache health
	health := "healthy"
	if expiredCount > len(m.entries)/2 {
		health = "degraded"
	}
	if !m.initialized {
		health = "unhealthy"
	}

	return &interfaces.CacheStats{
		TotalEntries:   len(m.entries),
		TotalSize:      m.metrics.CurrentSize,
		HitRate:        hitRate,
		ExpiredEntries: expiredCount,
		LastCleanup:    m.lastCleanup,
		CacheLocation:  m.cacheDir,
		OfflineMode:    m.offlineMode,
		CacheHealth:    health,
	}, nil
}

// GetSize returns the total cache size
func (m *Manager) GetSize() (int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.metrics.CurrentSize, nil
}

// GetLocation returns the cache directory location
func (m *Manager) GetLocation() string {
	return m.cacheDir
}

// Exists checks if a key exists in the cache
func (m *Manager) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, exists := m.entries[key]
	if !exists {
		return false
	}

	// Check if entry is expired
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
}

// GetKeys returns all cache keys
func (m *Manager) GetKeys() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]string, 0, len(m.entries))
	now := time.Now()

	for key, entry := range m.entries {
		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys, nil
}

// GetKeysByPattern returns cache keys matching a pattern
func (m *Manager) GetKeysByPattern(pattern string) ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	keys := make([]string, 0)
	now := time.Now()

	for key, entry := range m.entries {
		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}

		if regex.MatchString(key) {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	return keys, nil
}

// ValidateCache validates the cache integrity
func (m *Manager) ValidateCache() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if cache directory exists
	if _, err := os.Stat(m.cacheDir); os.IsNotExist(err) {
		return fmt.Errorf("cache directory does not exist: %s", m.cacheDir)
	}

	// Check if cache directory is writable
	testFile := filepath.Join(m.cacheDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	if err := os.Remove(testFile); err != nil {
		fmt.Printf("Warning: Failed to remove test file: %v\n", err)
	}

	// Validate cache entries
	now := time.Now()
	issues := make([]string, 0)

	for key, entry := range m.entries {
		if entry == nil {
			issues = append(issues, fmt.Sprintf("nil entry for key: %s", key))
			continue
		}

		if entry.Key != key {
			issues = append(issues, fmt.Sprintf("key mismatch: expected %s, got %s", key, entry.Key))
		}

		if entry.Size < 0 {
			issues = append(issues, fmt.Sprintf("negative size for key: %s", key))
		}

		if entry.CreatedAt.After(now) {
			issues = append(issues, fmt.Sprintf("future creation time for key: %s", key))
		}

		if entry.AccessCount < 0 {
			issues = append(issues, fmt.Sprintf("negative access count for key: %s", key))
		}
	}

	if len(issues) > 0 {
		return fmt.Errorf("cache validation failed: %s", strings.Join(issues, "; "))
	}

	return nil
}

// RepairCache repairs corrupted cache data
func (m *Manager) RepairCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(m.cacheDir, 0750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Clean up corrupted entries
	now := time.Now()
	repairedEntries := make(map[string]*interfaces.CacheEntry)

	for key, entry := range m.entries {
		if entry == nil {
			continue // Skip nil entries
		}

		// Fix key mismatch
		if entry.Key != key {
			entry.Key = key
		}

		// Fix negative size
		if entry.Size < 0 {
			entry.Size = m.calculateSize(entry.Value)
		}

		// Fix future timestamps
		if entry.CreatedAt.After(now) {
			entry.CreatedAt = now
		}
		if entry.UpdatedAt.After(now) {
			entry.UpdatedAt = now
		}
		if entry.AccessedAt.After(now) {
			entry.AccessedAt = now
		}

		// Fix negative access count
		if entry.AccessCount < 0 {
			entry.AccessCount = 0
		}

		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}

		repairedEntries[key] = entry
	}

	m.entries = repairedEntries

	// Recalculate metrics
	totalSize := int64(0)
	for _, entry := range m.entries {
		totalSize += entry.Size
	}

	m.metrics.CurrentSize = totalSize
	m.metrics.CurrentEntries = len(m.entries)

	// Save repaired cache
	if err := m.saveCache(); err != nil {
		return fmt.Errorf("failed to save repaired cache: %w", err)
	}

	return nil
}

// EnableOfflineMode enables offline mode
func (m *Manager) EnableOfflineMode() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.offlineMode = true
	m.config.OfflineMode = true

	// Preload essential keys if configured
	for _, key := range m.config.PreloadKeys {
		if !m.Exists(key) {
			// Log warning about missing preload key
			fmt.Printf("Warning: Preload key not found in cache: %s\n", key)
		}
	}

	return nil
}

// DisableOfflineMode disables offline mode
func (m *Manager) DisableOfflineMode() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.offlineMode = false
	m.config.OfflineMode = false

	return nil
}

// IsOfflineMode returns whether offline mode is enabled
func (m *Manager) IsOfflineMode() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.offlineMode
}

// CompactCache compacts the cache
func (m *Manager) CompactCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clean expired entries first
	now := time.Now()
	for key, entry := range m.entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			delete(m.entries, key)
			m.metrics.CurrentSize -= entry.Size
		}
	}

	// Recalculate metrics
	totalSize := int64(0)
	for _, entry := range m.entries {
		totalSize += entry.Size
	}

	m.metrics.CurrentSize = totalSize
	m.metrics.CurrentEntries = len(m.entries)

	// Save compacted cache
	if err := m.saveCache(); err != nil {
		return fmt.Errorf("failed to save compacted cache: %w", err)
	}

	return nil
}

// BackupCache backs up the cache to a file
func (m *Manager) BackupCache(path string) error {
	// Validate path to prevent path traversal
	if err := validatePath(path); err != nil {
		return fmt.Errorf("invalid backup path: %w", err)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create backup directory if needed
	backupDir := filepath.Dir(path)
	if err := os.MkdirAll(backupDir, 0750); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup data
	backup := cacheFile{
		Version:   "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Config:    m.config,
		Entries:   m.entries,
		Metrics:   m.metrics,
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

// RestoreCache restores the cache from a backup file
func (m *Manager) RestoreCache(path string) error {
	// Validate path to prevent path traversal
	if err := validatePath(path); err != nil {
		return fmt.Errorf("invalid restore path: %w", err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// #nosec G304 - path is validated by validatePath function
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: Failed to close backup file: %v\n", err)
		}
	}()

	var backup cacheFile
	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		return fmt.Errorf("failed to decode backup file: %w", err)
	}

	// Restore data
	m.entries = backup.Entries
	if backup.Config != nil {
		m.config = backup.Config
	}
	if backup.Metrics != nil {
		m.metrics = backup.Metrics
	}

	// Save restored cache
	if err := m.saveCache(); err != nil {
		return fmt.Errorf("failed to save restored cache: %w", err)
	}

	return nil
}

// SyncCache synchronizes the cache
func (m *Manager) SyncCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.saveCache()
}

// SetTTL sets TTL for a key
func (m *Manager) SetTTL(key string, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry, exists := m.entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	entry.TTL = ttl
	if ttl > 0 {
		expiresAt := time.Now().Add(ttl)
		entry.ExpiresAt = &expiresAt
	} else {
		entry.ExpiresAt = nil
	}

	return nil
}

// GetTTL gets TTL for a key
func (m *Manager) GetTTL(key string) (time.Duration, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, exists := m.entries[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	if entry.ExpiresAt == nil {
		return 0, nil // No expiration
	}

	remaining := time.Until(*entry.ExpiresAt)
	if remaining < 0 {
		return 0, nil // Already expired
	}

	return remaining, nil
}

// RefreshTTL refreshes TTL for a key
func (m *Manager) RefreshTTL(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry, exists := m.entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	if entry.TTL > 0 {
		expiresAt := time.Now().Add(entry.TTL)
		entry.ExpiresAt = &expiresAt
	}

	return nil
}

// GetExpiredKeys returns expired keys
func (m *Manager) GetExpiredKeys() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range m.entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	sort.Strings(expiredKeys)
	return expiredKeys, nil
}

// SetCacheConfig sets cache configuration
func (m *Manager) SetCacheConfig(config *interfaces.CacheConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	m.config = config
	m.offlineMode = config.OfflineMode

	return nil
}

// GetCacheConfig gets cache configuration
func (m *Manager) GetCacheConfig() (*interfaces.CacheConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *m.config
	return &configCopy, nil
}

// SetMaxSize sets maximum cache size
func (m *Manager) SetMaxSize(size int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.config.MaxSize = size
	m.metrics.MaxSize = size

	// Evict entries if current size exceeds new limit
	if m.metrics.CurrentSize > size {
		return m.evictBySize(0)
	}

	return nil
}

// SetDefaultTTL sets default TTL
func (m *Manager) SetDefaultTTL(ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.config.DefaultTTL = ttl
	return nil
}

// OnCacheHit sets cache hit callback
func (m *Manager) OnCacheHit(callback func(key string)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.callbacks.onHit = callback
}

// OnCacheMiss sets cache miss callback
func (m *Manager) OnCacheMiss(callback func(key string)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.callbacks.onMiss = callback
}

// OnCacheEviction sets cache eviction callback
func (m *Manager) OnCacheEviction(callback func(key string, reason string)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.callbacks.onEviction = callback
}

// GetHitRate returns cache hit rate
func (m *Manager) GetHitRate() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.metrics.Gets == 0 {
		return 0.0
	}

	return float64(m.metrics.Hits) / float64(m.metrics.Gets)
}

// GetMissRate returns cache miss rate
func (m *Manager) GetMissRate() float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.metrics.Gets == 0 {
		return 0.0
	}

	return float64(m.metrics.Misses) / float64(m.metrics.Gets)
}

// validatePath validates that a path is safe to use (prevents path traversal)
func validatePath(path string) error {
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
