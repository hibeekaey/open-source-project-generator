// Package metrics provides cache metrics collection and reporting.
package metrics

import (
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Collector handles cache metrics collection and calculation.
type Collector struct {
	metrics     *interfaces.CacheMetrics
	mutex       sync.RWMutex
	startTime   time.Time
	lastCleanup time.Time
}

// NewCollector creates a new metrics collector.
func NewCollector() *Collector {
	now := time.Now()
	return &Collector{
		metrics: &interfaces.CacheMetrics{
			CurrentSize:    0,
			MaxSize:        0,
			CurrentEntries: 0,
			MaxEntries:     0,
		},
		startTime:   now,
		lastCleanup: now,
	}
}

// RecordHit records a cache hit.
func (c *Collector) RecordHit(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.Hits++
	c.metrics.Gets++
	c.updateRates()
}

// RecordMiss records a cache miss.
func (c *Collector) RecordMiss(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.Misses++
	c.metrics.Gets++
	c.updateRates()
}

// RecordSet records a cache set operation.
func (c *Collector) RecordSet(key string, size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.Sets++
	c.metrics.CurrentSize += size
}

// RecordDelete records a cache delete operation.
func (c *Collector) RecordDelete(key string, size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.Deletes++
	c.metrics.CurrentSize -= size
	if c.metrics.CurrentSize < 0 {
		c.metrics.CurrentSize = 0
	}
}

// RecordEviction records a cache eviction.
func (c *Collector) RecordEviction(key string, reason string, size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.Evictions++
	c.metrics.CurrentSize -= size
	if c.metrics.CurrentSize < 0 {
		c.metrics.CurrentSize = 0
	}
}

// UpdateSize updates the current cache size and entry count.
func (c *Collector) UpdateSize(size int64, entries int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.CurrentSize = size
	c.metrics.CurrentEntries = entries
}

// SetLimits sets the maximum size and entry limits.
func (c *Collector) SetLimits(maxSize int64, maxEntries int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.MaxSize = maxSize
	c.metrics.MaxEntries = maxEntries
}

// RecordCleanup records a cleanup operation.
func (c *Collector) RecordCleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.lastCleanup = time.Now()
	c.metrics.LastCleanup = c.lastCleanup
}

// RecordCompaction records a compaction operation.
func (c *Collector) RecordCompaction() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.LastCompaction = time.Now()
}

// RecordBackup records a backup operation.
func (c *Collector) RecordBackup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.metrics.LastBackup = time.Now()
}

// GetMetrics returns a copy of the current metrics.
func (c *Collector) GetMetrics() *interfaces.CacheMetrics {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Return a copy to prevent external modification
	metricsCopy := *c.metrics
	return &metricsCopy
}

// GetHitRate returns the current hit rate.
func (c *Collector) GetHitRate() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.metrics.Gets == 0 {
		return 0.0
	}
	return float64(c.metrics.Hits) / float64(c.metrics.Gets)
}

// GetMissRate returns the current miss rate.
func (c *Collector) GetMissRate() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.metrics.Gets == 0 {
		return 0.0
	}
	return float64(c.metrics.Misses) / float64(c.metrics.Gets)
}

// Reset resets all metrics to zero.
func (c *Collector) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	c.metrics = &interfaces.CacheMetrics{
		CurrentSize:    c.metrics.CurrentSize,    // Keep current size
		MaxSize:        c.metrics.MaxSize,        // Keep limits
		CurrentEntries: c.metrics.CurrentEntries, // Keep current entries
		MaxEntries:     c.metrics.MaxEntries,     // Keep limits
	}
	c.startTime = now
	c.lastCleanup = now
}

// GetUptime returns the cache uptime.
func (c *Collector) GetUptime() time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return time.Since(c.startTime)
}

// GetLastCleanup returns the time of the last cleanup.
func (c *Collector) GetLastCleanup() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.lastCleanup
}

// SetMetrics sets the metrics (used for loading from storage).
func (c *Collector) SetMetrics(metrics *interfaces.CacheMetrics) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if metrics != nil {
		c.metrics = metrics
		c.updateRates()
	}
}

// updateRates updates the hit and miss rates (must be called with lock held).
func (c *Collector) updateRates() {
	if c.metrics.Gets > 0 {
		c.metrics.HitRate = float64(c.metrics.Hits) / float64(c.metrics.Gets)
		c.metrics.MissRate = float64(c.metrics.Misses) / float64(c.metrics.Gets)
	} else {
		c.metrics.HitRate = 0.0
		c.metrics.MissRate = 0.0
	}
}

// GetTotalOperations returns the total number of operations.
func (c *Collector) GetTotalOperations() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.metrics.Gets + c.metrics.Sets + c.metrics.Deletes
}
