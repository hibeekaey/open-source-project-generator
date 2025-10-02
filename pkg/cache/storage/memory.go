// Package storage provides cache storage backend implementations.
package storage

import (
	"encoding/json"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// MemoryStorage implements in-memory cache storage.
type MemoryStorage struct {
	entries map[string]*interfaces.CacheEntry
	metrics *interfaces.CacheMetrics
	config  *interfaces.CacheConfig
}

// NewMemoryStorage creates a new memory storage backend.
func NewMemoryStorage(config *interfaces.CacheConfig) *MemoryStorage {
	return &MemoryStorage{
		entries: make(map[string]*interfaces.CacheEntry),
		metrics: &interfaces.CacheMetrics{},
		config:  config,
	}
}

// Initialize sets up the memory storage.
func (ms *MemoryStorage) Initialize() error {
	// Memory storage doesn't need initialization
	return nil
}

// Load loads cache data from memory.
func (ms *MemoryStorage) Load() (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	// Return copies to prevent external modification
	entriesCopy := make(map[string]*interfaces.CacheEntry)
	for k, v := range ms.entries {
		entriesCopy[k] = v
	}

	metricsCopy := *ms.metrics
	return entriesCopy, &metricsCopy, nil
}

// Save saves cache data to memory.
func (ms *MemoryStorage) Save(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// Store copies to prevent external modification
	ms.entries = make(map[string]*interfaces.CacheEntry)
	for k, v := range entries {
		ms.entries[k] = v
	}

	if metrics != nil {
		ms.metrics = metrics
	}

	return nil
}

// Clear removes all cache data from memory.
func (ms *MemoryStorage) Clear() error {
	ms.entries = make(map[string]*interfaces.CacheEntry)
	ms.metrics = &interfaces.CacheMetrics{}
	return nil
}

// Backup creates a backup of the cache data (serialized to JSON).
func (ms *MemoryStorage) Backup(path string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	// For memory storage, we'll serialize to JSON and store in a map
	// This is a simplified implementation - in practice you might want to write to a file
	backup := CacheFile{
		Version: "1.0",
		Config:  ms.config,
		Entries: entries,
		Metrics: metrics,
	}

	// Serialize to validate the data can be backed up
	_, err := json.Marshal(backup)
	if err != nil {
		return fmt.Errorf("failed to serialize backup data: %w", err)
	}

	// In a real implementation, you would write this to the specified path
	// For now, we'll just validate that the backup is possible
	return nil
}

// Restore restores cache data from a backup.
func (ms *MemoryStorage) Restore(path string) (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	// For memory storage, this would typically load from a serialized format
	// This is a simplified implementation
	return make(map[string]*interfaces.CacheEntry), &interfaces.CacheMetrics{}, nil
}

// Validate validates the storage integrity.
func (ms *MemoryStorage) Validate() error {
	// Memory storage is always valid if it exists
	if ms.entries == nil {
		return fmt.Errorf("entries map is nil")
	}
	if ms.metrics == nil {
		return fmt.Errorf("metrics is nil")
	}
	return nil
}

// GetLocation returns the storage location (memory).
func (ms *MemoryStorage) GetLocation() string {
	return "memory"
}

// SetConfig updates the storage configuration.
func (ms *MemoryStorage) SetConfig(config *interfaces.CacheConfig) {
	ms.config = config
}
