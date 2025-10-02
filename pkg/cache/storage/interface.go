// Package storage provides cache storage backend implementations.
package storage

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheFile represents the structure of the cache file used by storage backends.
type CacheFile struct {
	Version   string                            `json:"version"`
	CreatedAt time.Time                         `json:"created_at"`
	UpdatedAt time.Time                         `json:"updated_at"`
	Config    *interfaces.CacheConfig           `json:"config"`
	Entries   map[string]*interfaces.CacheEntry `json:"entries"`
	Metrics   *interfaces.CacheMetrics          `json:"metrics"`
}

// StorageBackend defines the interface for cache storage backends.
type StorageBackend interface {
	// Initialize sets up the storage backend
	Initialize() error

	// Load loads cache data from storage
	Load() (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error)

	// Save saves cache data to storage
	Save(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error

	// Clear removes all cache data from storage
	Clear() error

	// Backup creates a backup of the cache data
	Backup(path string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error

	// Restore restores cache data from a backup
	Restore(path string) (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error)

	// Validate validates the storage integrity
	Validate() error

	// GetLocation returns the storage location
	GetLocation() string

	// SetConfig updates the storage configuration
	SetConfig(config *interfaces.CacheConfig)
}

// StorageManager manages different storage backends.
type StorageManager struct {
	backend StorageBackend
	config  *interfaces.CacheConfig
}

// NewStorageManager creates a new storage manager.
func NewStorageManager(cacheDir string, config *interfaces.CacheConfig) *StorageManager {
	var backend StorageBackend

	// Choose storage backend based on configuration
	if config.PersistToDisk {
		backend = NewFilesystemStorage(cacheDir, config)
	} else {
		backend = NewMemoryStorage(config)
	}

	return &StorageManager{
		backend: backend,
		config:  config,
	}
}

// Initialize initializes the storage backend.
func (sm *StorageManager) Initialize() error {
	return sm.backend.Initialize()
}

// Load loads cache data from the storage backend.
func (sm *StorageManager) Load() (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	return sm.backend.Load()
}

// Save saves cache data to the storage backend.
func (sm *StorageManager) Save(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	return sm.backend.Save(entries, metrics)
}

// Clear clears all cache data from the storage backend.
func (sm *StorageManager) Clear() error {
	return sm.backend.Clear()
}

// Backup creates a backup of the cache data.
func (sm *StorageManager) Backup(path string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) error {
	return sm.backend.Backup(path, entries, metrics)
}

// Restore restores cache data from a backup.
func (sm *StorageManager) Restore(path string) (map[string]*interfaces.CacheEntry, *interfaces.CacheMetrics, error) {
	return sm.backend.Restore(path)
}

// Validate validates the storage integrity.
func (sm *StorageManager) Validate() error {
	return sm.backend.Validate()
}

// GetLocation returns the storage location.
func (sm *StorageManager) GetLocation() string {
	return sm.backend.GetLocation()
}

// SetConfig updates the storage configuration and switches backends if needed.
func (sm *StorageManager) SetConfig(config *interfaces.CacheConfig) {
	sm.config = config
	sm.backend.SetConfig(config)

	// Switch storage backend if persistence setting changed
	currentLocation := sm.backend.GetLocation()
	if config.PersistToDisk && currentLocation == "memory" {
		// Switch to filesystem storage
		sm.backend = NewFilesystemStorage(config.Location, config)
	} else if !config.PersistToDisk && currentLocation != "memory" {
		// Switch to memory storage
		sm.backend = NewMemoryStorage(config)
	}
}
