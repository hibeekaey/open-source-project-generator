// Package interfaces defines the core contracts and interfaces for the
// Open Source Project Generator components.
package interfaces

import "time"

// CacheManager defines the contract for cache management operations.
//
// This interface abstracts caching functionality to enable different
// cache implementations and storage backends.
type CacheManager interface {
	// Cache operations
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Clean() error

	// Cache information
	GetStats() (*CacheStats, error)
	GetSize() (int64, error)
	GetLocation() string

	// Cache validation
	ValidateCache() error
	RepairCache() error

	// Offline support
	EnableOfflineMode() error
	DisableOfflineMode() error
	IsOfflineMode() bool
}

// CacheStats contains cache statistics and information
type CacheStats struct {
	TotalEntries   int       `json:"total_entries"`
	TotalSize      int64     `json:"total_size"`
	HitRate        float64   `json:"hit_rate"`
	ExpiredEntries int       `json:"expired_entries"`
	LastCleanup    time.Time `json:"last_cleanup"`
	CacheLocation  string    `json:"cache_location"`
}
