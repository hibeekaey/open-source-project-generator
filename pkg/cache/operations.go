// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// operationCallbacks holds callback functions for cache events
type operationCallbacks struct {
	onHit      func(key string)
	onMiss     func(key string)
	onEviction func(key string, reason string)
}

// CacheOperations handles cache entry operations and management.
type CacheOperations struct {
	config    *interfaces.CacheConfig
	callbacks *operationCallbacks
}

// NewCacheOperations creates a new cache operations instance.
func NewCacheOperations(config *interfaces.CacheConfig) *CacheOperations {
	return &CacheOperations{
		config:    config,
		callbacks: &operationCallbacks{},
	}
}

// GetEntry retrieves a cache entry and handles hit/miss tracking.
func (co *CacheOperations) GetEntry(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (*interfaces.CacheEntry, bool) {
	entry, exists := entries[key]
	if !exists {
		metrics.Misses++
		metrics.Gets++
		if co.callbacks.onMiss != nil {
			co.callbacks.onMiss(key)
		}
		return nil, false
	}

	// Check if entry is expired
	now := time.Now()
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
		delete(entries, key)
		metrics.Misses++
		metrics.Gets++
		if co.callbacks.onMiss != nil {
			co.callbacks.onMiss(key)
		}
		return nil, false
	}

	// Update access information
	entry.AccessedAt = now
	entry.AccessCount++

	metrics.Hits++
	metrics.Gets++
	if co.callbacks.onHit != nil {
		co.callbacks.onHit(key)
	}

	return entry, true
}

// CreateEntry creates a new cache entry with proper metadata.
func (co *CacheOperations) CreateEntry(key string, value any, ttl time.Duration) *interfaces.CacheEntry {
	now := time.Now()
	size := co.calculateSize(value)

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
	} else if co.config.DefaultTTL > 0 {
		expiresAt := now.Add(co.config.DefaultTTL)
		entry.ExpiresAt = &expiresAt
		entry.TTL = co.config.DefaultTTL
	}

	// Compress if enabled and value is large enough
	if co.config.EnableCompression && size > 1024 {
		if err := co.compressEntry(entry); err != nil {
			// Log warning but continue without compression
			fmt.Printf("Warning: Failed to compress cache entry: %v\n", err)
		}
	}

	return entry
}

// DeleteEntry removes an entry and updates metrics.
func (co *CacheOperations) DeleteEntry(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	entry, exists := entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	delete(entries, key)
	metrics.CurrentSize -= entry.Size
	metrics.CurrentEntries = len(entries)
	metrics.Deletes++

	if co.callbacks.onEviction != nil {
		co.callbacks.onEviction(key, interfaces.EvictionReasonManual)
	}

	return nil
}

// CleanExpiredEntries removes expired entries from the cache.
func (co *CacheOperations) CleanExpiredEntries(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) []string {
	now := time.Now()
	expiredKeys := make([]string, 0)

	// Find expired entries
	for key, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired entries
	for _, key := range expiredKeys {
		entry := entries[key]
		delete(entries, key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if co.callbacks.onEviction != nil {
			co.callbacks.onEviction(key, interfaces.EvictionReasonTTL)
		}
	}

	metrics.CurrentEntries = len(entries)
	return expiredKeys
}

// EvictIfNeeded evicts entries if cache limits are exceeded.
func (co *CacheOperations) EvictIfNeeded(newEntrySize int64, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// Check size limit
	if co.config.MaxSize > 0 && metrics.CurrentSize+newEntrySize > co.config.MaxSize {
		if err := co.evictBySize(newEntrySize, entries, metrics); err != nil {
			return err
		}
	}

	// Check entry count limit
	if co.config.MaxEntries > 0 && len(entries) >= co.config.MaxEntries {
		if err := co.evictByCount(entries, metrics); err != nil {
			return err
		}
	}

	return nil
}

// evictBySize evicts entries to make room for new entry
func (co *CacheOperations) evictBySize(newEntrySize int64, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	targetSize := int64(float64(co.config.MaxSize) * (1.0 - co.config.EvictionRatio))

	// Get entries sorted by eviction policy
	sortedEntries := co.getSortedEntriesForEviction(entries)

	for _, entry := range sortedEntries {
		if metrics.CurrentSize+newEntrySize <= targetSize {
			break
		}

		delete(entries, entry.Key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if co.callbacks.onEviction != nil {
			co.callbacks.onEviction(entry.Key, interfaces.EvictionReasonSize)
		}
	}

	return nil
}

// evictByCount evicts entries to make room for new entry
func (co *CacheOperations) evictByCount(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	targetCount := int(float64(co.config.MaxEntries) * (1.0 - co.config.EvictionRatio))

	// Get entries sorted by eviction policy
	sortedEntries := co.getSortedEntriesForEviction(entries)

	for _, entry := range sortedEntries {
		if len(entries) <= targetCount {
			break
		}

		delete(entries, entry.Key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if co.callbacks.onEviction != nil {
			co.callbacks.onEviction(entry.Key, interfaces.EvictionReasonCapacity)
		}
	}

	return nil
}

// getSortedEntriesForEviction returns entries sorted by eviction policy
func (co *CacheOperations) getSortedEntriesForEviction(entries map[string]*interfaces.CacheEntry) []*interfaces.CacheEntry {
	entryList := make([]*interfaces.CacheEntry, 0, len(entries))
	for _, entry := range entries {
		entryList = append(entryList, entry)
	}

	switch co.config.EvictionPolicy {
	case "lru": // Least Recently Used
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessedAt.Before(entryList[j].AccessedAt)
		})
	case "lfu": // Least Frequently Used
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessCount < entryList[j].AccessCount
		})
	case "fifo": // First In, First Out
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].CreatedAt.Before(entryList[j].CreatedAt)
		})
	case "ttl": // Shortest TTL first
		sort.Slice(entryList, func(i, j int) bool {
			if entryList[i].ExpiresAt == nil && entryList[j].ExpiresAt == nil {
				return entryList[i].CreatedAt.Before(entryList[j].CreatedAt)
			}
			if entryList[i].ExpiresAt == nil {
				return false
			}
			if entryList[j].ExpiresAt == nil {
				return true
			}
			return entryList[i].ExpiresAt.Before(*entryList[j].ExpiresAt)
		})
	default: // Default to LRU
		sort.Slice(entryList, func(i, j int) bool {
			return entryList[i].AccessedAt.Before(entryList[j].AccessedAt)
		})
	}

	return entryList
}

// GetKeys returns all non-expired cache keys.
func (co *CacheOperations) GetKeys(entries map[string]*interfaces.CacheEntry) []string {
	keys := make([]string, 0, len(entries))
	now := time.Now()

	for key, entry := range entries {
		// Skip expired entries
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// GetKeysByPattern returns cache keys matching a pattern.
func (co *CacheOperations) GetKeysByPattern(pattern string, entries map[string]*interfaces.CacheEntry) ([]string, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	keys := make([]string, 0)
	now := time.Now()

	for key, entry := range entries {
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

// GetExpiredKeys returns all expired keys.
func (co *CacheOperations) GetExpiredKeys(entries map[string]*interfaces.CacheEntry) []string {
	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	sort.Strings(expiredKeys)
	return expiredKeys
}

// ValidateEntries validates cache entries for integrity.
func (co *CacheOperations) ValidateEntries(entries map[string]*interfaces.CacheEntry) []string {
	now := time.Now()
	issues := make([]string, 0)

	for key, entry := range entries {
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

	return issues
}

// RepairEntries repairs corrupted cache entries.
func (co *CacheOperations) RepairEntries(entries map[string]*interfaces.CacheEntry) map[string]*interfaces.CacheEntry {
	now := time.Now()
	repairedEntries := make(map[string]*interfaces.CacheEntry)

	for key, entry := range entries {
		if entry == nil {
			continue // Skip nil entries
		}

		// Fix key mismatch
		if entry.Key != key {
			entry.Key = key
		}

		// Fix negative size
		if entry.Size < 0 {
			entry.Size = co.calculateSize(entry.Value)
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

	return repairedEntries
}

// UpdateTTL updates the TTL for an entry.
func (co *CacheOperations) UpdateTTL(key string, ttl time.Duration, entries map[string]*interfaces.CacheEntry) error {
	entry, exists := entries[key]
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

// GetTTL gets the remaining TTL for an entry.
func (co *CacheOperations) GetTTL(key string, entries map[string]*interfaces.CacheEntry) (time.Duration, error) {
	entry, exists := entries[key]
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

// RefreshTTL refreshes the TTL for an entry.
func (co *CacheOperations) RefreshTTL(key string, entries map[string]*interfaces.CacheEntry) error {
	entry, exists := entries[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	if entry.TTL > 0 {
		expiresAt := time.Now().Add(entry.TTL)
		entry.ExpiresAt = &expiresAt
	}

	return nil
}

// SetCallbacks sets the callback functions for cache events.
func (co *CacheOperations) SetCallbacks(onHit, onMiss func(string), onEviction func(string, string)) {
	co.callbacks.onHit = onHit
	co.callbacks.onMiss = onMiss
	co.callbacks.onEviction = onEviction
}

// SetConfig updates the operations configuration.
func (co *CacheOperations) SetConfig(config *interfaces.CacheConfig) {
	co.config = config
}

// calculateSize estimates the size of a value in bytes
func (co *CacheOperations) calculateSize(value any) int64 {
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

// compressEntry compresses the value in a cache entry
func (co *CacheOperations) compressEntry(entry *interfaces.CacheEntry) error {
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
	switch co.config.CompressionType {
	case "gzip":
		compressed, err = co.compressGzip(data)
	default:
		compressed, err = co.compressGzip(data)
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
		entry.Metadata["compression_type"] = co.config.CompressionType
	}

	return nil
}

// compressGzip compresses data using gzip
func (co *CacheOperations) compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, co.config.CompressionLevel)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
