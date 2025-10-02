// Package operations provides cache operation implementations.
package operations

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// SetOperation handles cache set operations.
type SetOperation struct {
	config    *interfaces.CacheConfig
	callbacks *OperationCallbacks
}

// NewSetOperation creates a new set operation handler.
func NewSetOperation(config *interfaces.CacheConfig, callbacks *OperationCallbacks) *SetOperation {
	return &SetOperation{
		config:    config,
		callbacks: callbacks,
	}
}

// Execute performs a cache set operation.
func (s *SetOperation) Execute(key string, value any, ttl time.Duration, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	now := time.Now()

	// Calculate entry size (approximate)
	size := s.calculateSize(value)

	// Check if we need to evict entries
	if err := s.evictIfNeeded(size, entries, metrics); err != nil {
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
	} else if s.config.DefaultTTL > 0 {
		expiresAt := now.Add(s.config.DefaultTTL)
		entry.ExpiresAt = &expiresAt
		entry.TTL = s.config.DefaultTTL
	}

	// Compress if enabled and value is large enough
	if s.config.EnableCompression && size > 1024 {
		if err := s.compressEntry(entry); err != nil {
			// Log warning but continue without compression
			fmt.Printf("Warning: Failed to compress cache entry: %v\n", err)
		}
	}

	// Store the entry
	entries[key] = entry

	// Update metrics
	metrics.Sets++
	metrics.CurrentEntries = len(entries)
	metrics.CurrentSize += size

	return nil
}

// calculateSize estimates the size of a value in bytes
func (s *SetOperation) calculateSize(value any) int64 {
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
func (s *SetOperation) evictIfNeeded(newEntrySize int64, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// Check size limit
	if s.config.MaxSize > 0 && metrics.CurrentSize+newEntrySize > s.config.MaxSize {
		if err := s.evictBySize(newEntrySize, entries, metrics); err != nil {
			return err
		}
	}

	// Check entry count limit
	if s.config.MaxEntries > 0 && len(entries) >= s.config.MaxEntries {
		if err := s.evictByCount(entries, metrics); err != nil {
			return err
		}
	}

	return nil
}

// evictBySize evicts entries to make room for new entry
func (s *SetOperation) evictBySize(newEntrySize int64, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	targetSize := int64(float64(s.config.MaxSize) * (1.0 - s.config.EvictionRatio))

	// Get entries sorted by eviction policy
	sortedEntries := s.getSortedEntriesForEviction(entries)

	for _, entry := range sortedEntries {
		if metrics.CurrentSize+newEntrySize <= targetSize {
			break
		}

		delete(entries, entry.Key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if s.callbacks != nil && s.callbacks.OnEviction != nil {
			s.callbacks.OnEviction(entry.Key, interfaces.EvictionReasonSize)
		}
	}

	return nil
}

// evictByCount evicts entries to make room for new entry
func (s *SetOperation) evictByCount(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	targetCount := int(float64(s.config.MaxEntries) * (1.0 - s.config.EvictionRatio))

	// Get entries sorted by eviction policy
	sortedEntries := s.getSortedEntriesForEviction(entries)

	for _, entry := range sortedEntries {
		if len(entries) <= targetCount {
			break
		}

		delete(entries, entry.Key)
		metrics.CurrentSize -= entry.Size
		metrics.Evictions++

		if s.callbacks != nil && s.callbacks.OnEviction != nil {
			s.callbacks.OnEviction(entry.Key, interfaces.EvictionReasonCapacity)
		}
	}

	return nil
}

// getSortedEntriesForEviction returns entries sorted by eviction policy
func (s *SetOperation) getSortedEntriesForEviction(entries map[string]*interfaces.CacheEntry) []*interfaces.CacheEntry {
	// Use the existing sorting logic from operations.go
	ops := NewCacheOperations(s.config)
	return ops.getSortedEntriesForEviction(entries)
}

// compressEntry compresses the value in a cache entry
func (s *SetOperation) compressEntry(entry *interfaces.CacheEntry) error {
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
	switch s.config.CompressionType {
	case "gzip":
		compressed, err = s.compressGzip(data)
	default:
		compressed, err = s.compressGzip(data)
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
		entry.Metadata["compression_type"] = s.config.CompressionType
	}

	return nil
}

// compressGzip compresses data using gzip
func (s *SetOperation) compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, s.config.CompressionLevel)
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
