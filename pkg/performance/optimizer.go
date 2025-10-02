// Package performance provides performance optimization and monitoring capabilities
// for the Open Source Project Generator CLI commands and operations.
package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CommandOptimizer provides performance optimization for CLI commands
type CommandOptimizer struct {
	cache         interfaces.CacheManager
	logger        interfaces.Logger
	metrics       *MetricsCollector
	lazyLoaders   map[string]LazyLoader
	optimizations map[string]OptimizationConfig
	mutex         sync.RWMutex
	initialized   bool
}

// OptimizationConfig defines optimization settings for specific operations
type OptimizationConfig struct {
	EnableCaching     bool          `json:"enable_caching"`
	CacheTTL          time.Duration `json:"cache_ttl"`
	EnableLazyLoading bool          `json:"enable_lazy_loading"`
	MaxConcurrency    int           `json:"max_concurrency"`
	TimeoutDuration   time.Duration `json:"timeout_duration"`
	EnableProfiling   bool          `json:"enable_profiling"`
}

// LazyLoader defines interface for lazy loading expensive operations
type LazyLoader interface {
	Load(ctx context.Context) (interface{}, error)
	IsLoaded() bool
	GetCacheKey() string
}

// CommandMetrics tracks performance metrics for CLI commands
type CommandMetrics struct {
	CommandName   string                 `json:"command_name"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      time.Duration          `json:"duration"`
	MemoryUsage   int64                  `json:"memory_usage"`
	CacheHits     int                    `json:"cache_hits"`
	CacheMisses   int                    `json:"cache_misses"`
	LazyLoads     int                    `json:"lazy_loads"`
	Optimizations []string               `json:"optimizations"`
	CustomMetrics map[string]interface{} `json:"custom_metrics"`
}

// NewCommandOptimizer creates a new command optimizer instance
func NewCommandOptimizer(cache interfaces.CacheManager, logger interfaces.Logger) *CommandOptimizer {
	optimizer := &CommandOptimizer{
		cache:         cache,
		logger:        logger,
		metrics:       NewMetricsCollector(),
		lazyLoaders:   make(map[string]LazyLoader),
		optimizations: make(map[string]OptimizationConfig),
	}

	// Set default optimizations for common commands
	optimizer.setDefaultOptimizations()

	return optimizer
}

// setDefaultOptimizations configures default optimization settings
func (co *CommandOptimizer) setDefaultOptimizations() {
	// Generate command optimizations
	co.optimizations["generate"] = OptimizationConfig{
		EnableCaching:     true,
		CacheTTL:          30 * time.Minute,
		EnableLazyLoading: true,
		MaxConcurrency:    4,
		TimeoutDuration:   5 * time.Minute,
		EnableProfiling:   false,
	}

	// Version command optimizations
	co.optimizations["version"] = OptimizationConfig{
		EnableCaching:     true,
		CacheTTL:          1 * time.Hour,
		EnableLazyLoading: true,
		MaxConcurrency:    2,
		TimeoutDuration:   30 * time.Second,
		EnableProfiling:   false,
	}

	// Audit command optimizations
	co.optimizations["audit"] = OptimizationConfig{
		EnableCaching:     true,
		CacheTTL:          15 * time.Minute,
		EnableLazyLoading: true,
		MaxConcurrency:    6,
		TimeoutDuration:   10 * time.Minute,
		EnableProfiling:   true,
	}

	// Template list optimizations
	co.optimizations["list-templates"] = OptimizationConfig{
		EnableCaching:     true,
		CacheTTL:          2 * time.Hour,
		EnableLazyLoading: true,
		MaxConcurrency:    2,
		TimeoutDuration:   1 * time.Minute,
		EnableProfiling:   false,
	}

	// Cache operations optimizations
	co.optimizations["cache"] = OptimizationConfig{
		EnableCaching:     false, // Don't cache cache operations
		CacheTTL:          0,
		EnableLazyLoading: false,
		MaxConcurrency:    1,
		TimeoutDuration:   2 * time.Minute,
		EnableProfiling:   false,
	}
}

// OptimizeCommand applies performance optimizations to a command execution
func (co *CommandOptimizer) OptimizeCommand(ctx context.Context, commandName string, operation func(ctx context.Context) (interface{}, error)) (interface{}, *CommandMetrics, error) {
	co.mutex.Lock()
	config, exists := co.optimizations[commandName]
	if !exists {
		config = OptimizationConfig{
			EnableCaching:     false,
			EnableLazyLoading: false,
			MaxConcurrency:    1,
			TimeoutDuration:   5 * time.Minute,
			EnableProfiling:   false,
		}
	}
	co.mutex.Unlock()

	// Start metrics collection
	metrics := &CommandMetrics{
		CommandName:   commandName,
		StartTime:     time.Now(),
		CustomMetrics: make(map[string]interface{}),
		Optimizations: make([]string, 0),
	}

	// Record initial memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialMemory := memStats.Alloc

	// Create context with timeout
	if config.TimeoutDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.TimeoutDuration)
		defer cancel()
		metrics.Optimizations = append(metrics.Optimizations, "timeout")
	}

	var result interface{}
	var err error

	// Try cache first if enabled
	if config.EnableCaching && co.cache != nil {
		cacheKey := fmt.Sprintf("cmd:%s:%s", commandName, co.generateCacheKey(ctx))

		if cached, cacheErr := co.cache.Get(cacheKey); cacheErr == nil {
			metrics.CacheHits++
			metrics.Optimizations = append(metrics.Optimizations, "cache_hit")
			result = cached

			if co.logger != nil {
				co.logger.Debug("Command result served from cache", "command", commandName, "cache_key", cacheKey)
			}
		} else {
			metrics.CacheMisses++

			// Execute operation
			result, err = operation(ctx)

			// Cache result if successful
			if err == nil && result != nil {
				if cacheErr := co.cache.Set(cacheKey, result, config.CacheTTL); cacheErr == nil {
					metrics.Optimizations = append(metrics.Optimizations, "cache_store")
				}
			}
		}
	} else {
		// Execute operation without caching
		result, err = operation(ctx)
	}

	// Record final metrics
	metrics.EndTime = time.Now()
	metrics.Duration = metrics.EndTime.Sub(metrics.StartTime)

	// Record final memory usage
	runtime.ReadMemStats(&memStats)
	if memStats.Alloc >= initialMemory {
		metrics.MemoryUsage = int64(memStats.Alloc - initialMemory)
	} else {
		metrics.MemoryUsage = 0
	}

	// Record metrics
	co.metrics.RecordCommand(metrics)

	if co.logger != nil {
		co.logger.LogPerformanceMetrics(commandName, map[string]interface{}{
			"duration_ms":   metrics.Duration.Milliseconds(),
			"memory_bytes":  metrics.MemoryUsage,
			"cache_hits":    metrics.CacheHits,
			"cache_misses":  metrics.CacheMisses,
			"optimizations": metrics.Optimizations,
		})
	}

	return result, metrics, err
}

// RegisterLazyLoader registers a lazy loader for expensive operations
func (co *CommandOptimizer) RegisterLazyLoader(key string, loader LazyLoader) {
	co.mutex.Lock()
	defer co.mutex.Unlock()
	co.lazyLoaders[key] = loader
}

// GetLazyLoader retrieves and optionally loads a lazy loader
func (co *CommandOptimizer) GetLazyLoader(ctx context.Context, key string, autoLoad bool) (interface{}, error) {
	co.mutex.RLock()
	loader, exists := co.lazyLoaders[key]
	co.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("lazy loader not found: %s", key)
	}

	if !loader.IsLoaded() && autoLoad {
		return loader.Load(ctx)
	}

	if !loader.IsLoaded() {
		return nil, fmt.Errorf("lazy loader not loaded: %s", key)
	}

	// If already loaded, try to get from cache
	if co.cache != nil {
		if cached, err := co.cache.Get(loader.GetCacheKey()); err == nil {
			return cached, nil
		}
	}

	return loader.Load(ctx)
}

// SetOptimizationConfig sets optimization configuration for a command
func (co *CommandOptimizer) SetOptimizationConfig(commandName string, config OptimizationConfig) {
	co.mutex.Lock()
	defer co.mutex.Unlock()
	co.optimizations[commandName] = config
}

// GetOptimizationConfig gets optimization configuration for a command
func (co *CommandOptimizer) GetOptimizationConfig(commandName string) (OptimizationConfig, bool) {
	co.mutex.RLock()
	defer co.mutex.RUnlock()
	config, exists := co.optimizations[commandName]
	return config, exists
}

// generateCacheKey generates a cache key based on context
func (co *CommandOptimizer) generateCacheKey(ctx context.Context) string {
	// Simple implementation - in production, this would be more sophisticated
	return fmt.Sprintf("%d", time.Now().Unix()/300) // 5-minute buckets
}

// GetMetrics returns current performance metrics
func (co *CommandOptimizer) GetMetrics() *MetricsCollector {
	return co.metrics
}

// ClearCache clears optimization-related cache entries
func (co *CommandOptimizer) ClearCache() error {
	if co.cache == nil {
		return fmt.Errorf("cache manager not available")
	}

	// Get all keys with cmd: prefix
	keys, err := co.cache.GetKeysByPattern("cmd:*")
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	// Delete optimization cache entries
	for _, key := range keys {
		if err := co.cache.Delete(key); err != nil {
			if co.logger != nil {
				co.logger.Warn("Failed to delete cache key", "key", key, "error", err)
			}
		}
	}

	return nil
}

// OptimizeFileSystemOperations provides optimized file system operations
func (co *CommandOptimizer) OptimizeFileSystemOperations(ctx context.Context, operations []func() error) error {
	config := OptimizationConfig{
		MaxConcurrency:  4,
		TimeoutDuration: 2 * time.Minute,
	}

	// Use semaphore to limit concurrency
	semaphore := make(chan struct{}, config.MaxConcurrency)
	errChan := make(chan error, len(operations))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, config.TimeoutDuration)
	defer cancel()

	var wg sync.WaitGroup

	for _, op := range operations {
		wg.Add(1)
		go func(operation func() error) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}

			// Execute operation
			if err := operation(); err != nil {
				errChan <- err
			}
		}(op)
	}

	// Wait for all operations to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors
	var errors []error
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("file system operations failed: %v", errors)
	}

	return nil
}

// ProfileCommand profiles a command execution and returns profiling data
func (co *CommandOptimizer) ProfileCommand(ctx context.Context, commandName string, operation func(ctx context.Context) error) (*ProfilingData, error) {
	config, exists := co.optimizations[commandName]
	if !exists || !config.EnableProfiling {
		return nil, fmt.Errorf("profiling not enabled for command: %s", commandName)
	}

	profiler := NewProfiler()

	// Start profiling
	if err := profiler.Start(); err != nil {
		return nil, fmt.Errorf("failed to start profiler: %w", err)
	}

	// Execute operation
	startTime := time.Now()
	err := operation(ctx)
	duration := time.Since(startTime)

	// Stop profiling
	profilingData, profErr := profiler.Stop()
	if profErr != nil {
		if co.logger != nil {
			co.logger.Warn("Failed to stop profiler", "error", profErr)
		}
	}

	if profilingData != nil {
		profilingData.CommandName = commandName
		profilingData.Duration = duration
		profilingData.Success = err == nil
	}

	return profilingData, err
}

// GetPerformanceReport generates a comprehensive performance report
func (co *CommandOptimizer) GetPerformanceReport() *MetricsReport {
	return co.metrics.GenerateReport()
}

// EnableOptimization enables specific optimization for a command
func (co *CommandOptimizer) EnableOptimization(commandName, optimizationType string) error {
	co.mutex.Lock()
	defer co.mutex.Unlock()

	config, exists := co.optimizations[commandName]
	if !exists {
		config = OptimizationConfig{}
	}

	switch optimizationType {
	case "caching":
		config.EnableCaching = true
		if config.CacheTTL == 0 {
			config.CacheTTL = 30 * time.Minute
		}
	case "lazy_loading":
		config.EnableLazyLoading = true
	case "profiling":
		config.EnableProfiling = true
	default:
		return fmt.Errorf("unknown optimization type: %s", optimizationType)
	}

	co.optimizations[commandName] = config
	return nil
}

// DisableOptimization disables specific optimization for a command
func (co *CommandOptimizer) DisableOptimization(commandName, optimizationType string) error {
	co.mutex.Lock()
	defer co.mutex.Unlock()

	config, exists := co.optimizations[commandName]
	if !exists {
		return fmt.Errorf("no optimization config found for command: %s", commandName)
	}

	switch optimizationType {
	case "caching":
		config.EnableCaching = false
	case "lazy_loading":
		config.EnableLazyLoading = false
	case "profiling":
		config.EnableProfiling = false
	default:
		return fmt.Errorf("unknown optimization type: %s", optimizationType)
	}

	co.optimizations[commandName] = config
	return nil
}
