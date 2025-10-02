// Package interfaces defines cache management interfaces.
//
// This file contains interfaces for comprehensive cache management
// including caching operations, offline mode support, and cache statistics.
package interfaces

import (
	"time"
)

// CacheManager defines the interface for cache management operations.
//
// This interface provides comprehensive cache management including:
//   - Cache operations (get, set, delete, clear)
//   - Cache statistics and monitoring
//   - Cache validation and repair
//   - Offline mode support
//   - TTL (Time To Live) management
type CacheManager interface {
	// Basic cache operations
	Get(key string) (any, error)
	Set(key string, value any, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) bool
	Clear() error
	Clean() error

	// Cache information
	GetStats() (*CacheStats, error)
	GetSize() (int64, error)
	GetLocation() string
	GetKeys() ([]string, error)
	GetKeysByPattern(pattern string) ([]string, error)

	// Cache validation and maintenance
	ValidateCache() error
	RepairCache() error
	CompactCache() error
	BackupCache(path string) error
	RestoreCache(path string) error

	// Offline support
	EnableOfflineMode() error
	DisableOfflineMode() error
	IsOfflineMode() bool
	SyncCache() error

	// TTL management
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	RefreshTTL(key string) error
	GetExpiredKeys() ([]string, error)

	// Cache configuration
	SetCacheConfig(config *CacheConfig) error
	GetCacheConfig() (*CacheConfig, error)
	SetMaxSize(size int64) error
	SetDefaultTTL(ttl time.Duration) error

	// Cache events and monitoring
	OnCacheHit(callback func(key string))
	OnCacheMiss(callback func(key string))
	OnCacheEviction(callback func(key string, reason string))
	GetHitRate() float64
	GetMissRate() float64
}

// CacheStats contains cache statistics and information
type CacheStats struct {
	TotalEntries int       `json:"total_entries"`
	TotalSize    int64     `json:"total_size"`
	HitCount     int64     `json:"hit_count"`
	MissCount    int64     `json:"miss_count"`
	HitRate      float64   `json:"hit_rate"`
	LastAccessed time.Time `json:"last_accessed"`
	LastModified time.Time `json:"last_modified"`
	CreatedAt    time.Time `json:"created_at"`
}

// CacheConfig defines configuration options for the cache
type CacheConfig struct {
	// Storage configuration
	Location   string        `json:"location"`
	MaxSize    int64         `json:"max_size"`
	MaxEntries int           `json:"max_entries"`
	DefaultTTL time.Duration `json:"default_ttl"`

	// Eviction policy
	EvictionPolicy string  `json:"eviction_policy"` // lru, lfu, fifo, ttl
	EvictionRatio  float64 `json:"eviction_ratio"`

	// Compression
	EnableCompression bool   `json:"enable_compression"`
	CompressionLevel  int    `json:"compression_level"`
	CompressionType   string `json:"compression_type"` // gzip, lz4, snappy

	// Persistence
	PersistToDisk   bool          `json:"persist_to_disk"`
	SyncInterval    time.Duration `json:"sync_interval"`
	BackupInterval  time.Duration `json:"backup_interval"`
	BackupRetention int           `json:"backup_retention"`

	// Performance
	ConcurrencyLevel int  `json:"concurrency_level"`
	EnableMetrics    bool `json:"enable_metrics"`
	EnableProfiling  bool `json:"enable_profiling"`

	// Offline mode
	OfflineMode    bool          `json:"offline_mode"`
	OfflineTTL     time.Duration `json:"offline_ttl"`
	OfflineMaxSize int64         `json:"offline_max_size"`
	PreloadKeys    []string      `json:"preload_keys"`
}

// CacheEntry represents a cache entry with metadata
type CacheEntry struct {
	Key         string         `json:"key"`
	Value       any            `json:"value"`
	Size        int64          `json:"size"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	AccessedAt  time.Time      `json:"accessed_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	TTL         time.Duration  `json:"ttl"`
	AccessCount int64          `json:"access_count"`
	Compressed  bool           `json:"compressed"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// CacheMetrics contains cache performance metrics
type CacheMetrics struct {
	// Hit/Miss statistics
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	HitRate  float64 `json:"hit_rate"`
	MissRate float64 `json:"miss_rate"`

	// Operation statistics
	Gets      int64 `json:"gets"`
	Sets      int64 `json:"sets"`
	Deletes   int64 `json:"deletes"`
	Evictions int64 `json:"evictions"`

	// Size statistics
	CurrentSize    int64 `json:"current_size"`
	MaxSize        int64 `json:"max_size"`
	CurrentEntries int   `json:"current_entries"`
	MaxEntries     int   `json:"max_entries"`

	// Performance statistics
	AverageGetTime  time.Duration `json:"average_get_time"`
	AverageSetTime  time.Duration `json:"average_set_time"`
	TotalOperations int64         `json:"total_operations"`

	// Maintenance statistics
	LastCleanup    time.Time `json:"last_cleanup"`
	LastCompaction time.Time `json:"last_compaction"`
	LastBackup     time.Time `json:"last_backup"`
}

// CacheHealth represents the health status of the cache
type CacheHealth struct {
	Status          string         `json:"status"` // healthy, degraded, unhealthy
	LastCheck       time.Time      `json:"last_check"`
	Issues          []CacheIssue   `json:"issues"`
	Warnings        []CacheWarning `json:"warnings"`
	Recommendations []string       `json:"recommendations"`
}

// CacheIssue represents a cache issue
type CacheIssue struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	DetectedAt  time.Time `json:"detected_at"`
	Resolution  string    `json:"resolution"`
	Fixable     bool      `json:"fixable"`
}

// CacheWarning represents a cache warning
type CacheWarning struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Threshold   any    `json:"threshold"`
	Current     any    `json:"current"`
	Suggestion  string `json:"suggestion"`
}

// CacheOperation represents a cache operation for logging/monitoring
type CacheOperation struct {
	Type      string        `json:"type"` // get, set, delete, clear, clean
	Key       string        `json:"key"`
	Success   bool          `json:"success"`
	Duration  time.Duration `json:"duration"`
	Size      int64         `json:"size"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// CacheEvictionReason defines reasons for cache eviction
const (
	EvictionReasonTTL      = "ttl"
	EvictionReasonSize     = "size"
	EvictionReasonCapacity = "capacity"
	EvictionReasonManual   = "manual"
	EvictionReasonError    = "error"
)

// CacheStatus defines cache status values
const (
	CacheStatusHealthy   = "healthy"
	CacheStatusDegraded  = "degraded"
	CacheStatusUnhealthy = "unhealthy"
	CacheStatusOffline   = "offline"
)

// CacheIssueType defines types of cache issues
const (
	CacheIssueTypeCorruption    = "corruption"
	CacheIssueTypePermission    = "permission"
	CacheIssueTypeDiskSpace     = "disk_space"
	CacheIssueTypePerformance   = "performance"
	CacheIssueTypeConfiguration = "configuration"
)

// CacheWarningType defines types of cache warnings
const (
	CacheWarningTypeSize        = "size"
	CacheWarningTypeHitRate     = "hit_rate"
	CacheWarningTypePerformance = "performance"
	CacheWarningTypeExpiration  = "expiration"
)

// StorageBackend defines the interface for cache storage backends.
type StorageBackend interface {
	// Initialize sets up the storage backend
	Initialize() error

	// Basic operations
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) bool
	Clear() error

	// Batch operations
	GetMultiple(keys []string) (map[string]interface{}, error)
	SetMultiple(items map[string]interface{}, ttl time.Duration) error
	DeleteMultiple(keys []string) error

	// Metadata operations
	GetKeys() ([]string, error)
	GetSize() (int64, error)
	GetStats() (*CacheStats, error)

	// Maintenance operations
	Cleanup() error
	Compact() error
	Backup(path string) error
	Restore(path string) error

	// Configuration
	SetConfig(config *CacheConfig) error
	GetConfig() (*CacheConfig, error)

	// Health and monitoring
	HealthCheck() error
	GetMetrics() (*CacheMetrics, error)

	// Lifecycle
	Close() error
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Location:          "~/.generator/cache",
		MaxSize:           1024 * 1024 * 1024, // 1GB
		MaxEntries:        10000,
		DefaultTTL:        24 * time.Hour,
		EvictionPolicy:    "lru",
		EvictionRatio:     0.1,
		EnableCompression: true,
		CompressionLevel:  6,
		CompressionType:   "gzip",
		PersistToDisk:     true,
		SyncInterval:      5 * time.Minute,
		BackupInterval:    24 * time.Hour,
		BackupRetention:   7,
		ConcurrencyLevel:  4,
		EnableMetrics:     true,
		EnableProfiling:   false,
		OfflineMode:       false,
		OfflineTTL:        7 * 24 * time.Hour, // 7 days
		OfflineMaxSize:    512 * 1024 * 1024,  // 512MB
		PreloadKeys:       []string{},
	}
}
