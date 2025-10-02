// Package cache provides caching functionality for the
// Open Source Project Generator.
package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cache/metrics"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cache/operations"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cache/storage"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cache/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Manager implements the CacheManager interface for cache operations.
type Manager struct {
	cacheDir    string
	config      *interfaces.CacheConfig
	entries     map[string]*interfaces.CacheEntry
	offlineMode bool
	mutex       sync.RWMutex
	initialized bool

	// Component managers
	operations *operations.CacheOperations
	storage    *storage.StorageManager
	metrics    *metrics.Collector
	validator  *validation.Validator
	reporter   *metrics.Reporter

	// Callbacks
	onHit      func(key string)
	onMiss     func(key string)
	onEviction func(key string, reason string)
}

// NewManager creates a new cache manager instance.
func NewManager(cacheDir string) interfaces.CacheManager {
	config := interfaces.DefaultCacheConfig()
	config.Location = cacheDir

	manager := &Manager{
		cacheDir: cacheDir,
		config:   config,
		entries:  make(map[string]*interfaces.CacheEntry),
	}

	// Initialize components
	manager.operations = operations.NewCacheOperations(config)
	manager.storage = storage.NewStorageManager(cacheDir, config)
	manager.metrics = metrics.NewCollector()
	manager.validator = validation.NewValidator(cacheDir, config)
	manager.reporter = metrics.NewReporter(manager.metrics)

	// Set up callbacks
	manager.operations.SetCallbacks(
		manager.onCacheHit,
		manager.onCacheMiss,
		manager.onCacheEviction,
	)

	// Initialize the cache
	if err := manager.initialize(); err != nil {
		fmt.Printf("Warning: Failed to initialize cache: %v\n", err)
	}

	return manager
}

// initialize sets up the cache and loads existing data.
func (m *Manager) initialize() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.storage.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	entries, loadedMetrics, err := m.storage.Load()
	if err != nil {
		entries = make(map[string]*interfaces.CacheEntry)
		loadedMetrics = &interfaces.CacheMetrics{}
	}

	m.entries = entries
	m.metrics.SetMetrics(loadedMetrics)
	m.metrics.SetLimits(m.config.MaxSize, m.config.MaxEntries)
	m.initialized = true
	return nil
}

// updateMetricsAndSave updates metrics and saves to storage if needed.
func (m *Manager) updateMetricsAndSave() {
	totalSize := int64(0)
	for _, entry := range m.entries {
		totalSize += entry.Size
	}
	m.metrics.UpdateSize(totalSize, len(m.entries))

	if m.config.PersistToDisk {
		currentMetrics := m.metrics.GetMetrics()
		if err := m.storage.Save(m.entries, currentMetrics); err != nil {
			fmt.Printf("Warning: Failed to save cache: %v\n", err)
		}
	}
}

func (m *Manager) Get(key string) (any, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.operations.Get(key, m.entries, m.metrics.GetMetrics())
}

func (m *Manager) Set(key string, value any, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if err := m.operations.Set(key, value, ttl, m.entries, m.metrics.GetMetrics()); err != nil {
		return err
	}
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if err := m.operations.Delete(key, m.entries, m.metrics.GetMetrics()); err != nil {
		return err
	}
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.operations.Exists(key, m.entries)
}

func (m *Manager) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.operations.Clear(m.entries, m.metrics.GetMetrics())
	if err := m.storage.Clear(); err != nil {
		return err
	}
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) Clean() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.operations.Clean(m.entries, m.metrics.GetMetrics())
	m.metrics.RecordCleanup()
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) GetStats() (*interfaces.CacheStats, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.reporter.GenerateStats(m.entries, m.cacheDir, m.offlineMode), nil
}

func (m *Manager) GetSize() (int64, error) { return m.metrics.GetMetrics().CurrentSize, nil }
func (m *Manager) GetLocation() string     { return m.cacheDir }

func (m *Manager) GetKeys() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.operations.GetKeys(m.entries), nil
}

func (m *Manager) GetKeysByPattern(pattern string) ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.operations.GetKeysByPattern(pattern, m.entries)
}

func (m *Manager) ValidateCache() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.validator.ValidateCache(m.entries)
}

func (m *Manager) RepairCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.entries = m.validator.RepairEntries(m.entries)
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) CompactCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if err := m.operations.Compact(m.entries, m.metrics.GetMetrics()); err != nil {
		return err
	}
	m.metrics.RecordCompaction()
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) BackupCache(path string) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if err := m.storage.Backup(path, m.entries, m.metrics.GetMetrics()); err != nil {
		return err
	}
	m.metrics.RecordBackup()
	return nil
}

func (m *Manager) RestoreCache(path string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	entries, restoredMetrics, err := m.storage.Restore(path)
	if err != nil {
		return err
	}
	m.entries = entries
	m.metrics.SetMetrics(restoredMetrics)
	m.updateMetricsAndSave()
	return nil
}

func (m *Manager) EnableOfflineMode() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.offlineMode = true
	m.config.OfflineMode = true
	return nil
}
func (m *Manager) DisableOfflineMode() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.offlineMode = false
	m.config.OfflineMode = false
	return nil
}
func (m *Manager) IsOfflineMode() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.offlineMode
}
func (m *Manager) SyncCache() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.storage.Save(m.entries, m.metrics.GetMetrics())
}

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

func (m *Manager) GetTTL(key string) (time.Duration, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	entry, exists := m.entries[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}
	if entry.ExpiresAt == nil {
		return 0, nil
	}
	remaining := time.Until(*entry.ExpiresAt)
	if remaining < 0 {
		return 0, nil
	}
	return remaining, nil
}

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

func (m *Manager) GetExpiredKeys() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.operations.GetExpiredKeys(m.entries), nil
}

// Cache configuration
func (m *Manager) SetCacheConfig(config *interfaces.CacheConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	m.config = config
	m.offlineMode = config.OfflineMode
	m.operations.SetConfig(config)
	m.storage.SetConfig(config)
	m.validator.SetConfig(config)
	m.metrics.SetLimits(config.MaxSize, config.MaxEntries)
	return nil
}

func (m *Manager) GetCacheConfig() (*interfaces.CacheConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	configCopy := *m.config
	return &configCopy, nil
}

func (m *Manager) SetMaxSize(size int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.config.MaxSize = size
	m.metrics.SetLimits(size, m.config.MaxEntries)
	m.operations.SetConfig(m.config)
	return nil
}

func (m *Manager) SetDefaultTTL(ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.config.DefaultTTL = ttl
	m.operations.SetConfig(m.config)
	return nil
}

// Cache events and monitoring
func (m *Manager) OnCacheHit(callback func(key string))                     { m.onHit = callback }
func (m *Manager) OnCacheMiss(callback func(key string))                    { m.onMiss = callback }
func (m *Manager) OnCacheEviction(callback func(key string, reason string)) { m.onEviction = callback }
func (m *Manager) GetHitRate() float64                                      { return m.metrics.GetHitRate() }
func (m *Manager) GetMissRate() float64                                     { return m.metrics.GetMissRate() }

// Internal callback handlers
func (m *Manager) onCacheHit(key string) {
	m.metrics.RecordHit(key)
	if m.onHit != nil {
		m.onHit(key)
	}
}

func (m *Manager) onCacheMiss(key string) {
	m.metrics.RecordMiss(key)
	if m.onMiss != nil {
		m.onMiss(key)
	}
}

func (m *Manager) onCacheEviction(key string, reason string) {
	if m.onEviction != nil {
		m.onEviction(key, reason)
	}
}
