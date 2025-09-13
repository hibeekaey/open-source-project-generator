package app

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/open-source-template-generator/pkg/utils"
)

// ResourceManager manages application resources and memory usage
type ResourceManager struct {
	memoryPool      *utils.MemoryPool
	stringOps       *utils.OptimizedStringOperations
	resourceManager *utils.ResourceManager
	cleanupTicker   *time.Ticker
	stopCleanup     chan struct{}
	mutex           sync.RWMutex
	activeResources map[string]*Resource
	maxMemoryUsage  int64
	gcThreshold     int64
	lastGCTime      time.Time
}

// Resource represents a managed resource
type Resource struct {
	ID         string
	Type       string
	Size       int64
	CreatedAt  time.Time
	LastAccess time.Time
	RefCount   int32
	Data       interface{}
}

// NewResourceManager creates a new resource manager
func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		memoryPool:      utils.NewMemoryPool(),
		stringOps:       utils.NewOptimizedStringOperations(),
		resourceManager: utils.NewResourceManager(),
		stopCleanup:     make(chan struct{}),
		activeResources: make(map[string]*Resource),
		maxMemoryUsage:  500 * 1024 * 1024, // 500MB limit
		gcThreshold:     100 * 1024 * 1024, // 100MB threshold for GC
	}

	// Start cleanup routine
	rm.cleanupTicker = time.NewTicker(30 * time.Second)
	go rm.cleanupRoutine()

	return rm
}

// GetBuffer gets a buffer from the memory pool
func (rm *ResourceManager) GetBuffer(size int) []byte {
	return rm.memoryPool.Get(size)
}

// PutBuffer returns a buffer to the memory pool
func (rm *ResourceManager) PutBuffer(buf []byte) {
	rm.memoryPool.Put(buf)
}

// GetStringOperations returns the optimized string operations instance
func (rm *ResourceManager) GetStringOperations() *utils.OptimizedStringOperations {
	return rm.stringOps
}

// RegisterResource registers a resource for management
func (rm *ResourceManager) RegisterResource(id, resourceType string, size int64, data interface{}) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Check memory limits
	currentUsage := rm.calculateCurrentUsage()
	if currentUsage+size > rm.maxMemoryUsage {
		// Try to free some resources first
		rm.freeUnusedResources()

		// Check again after cleanup
		currentUsage = rm.calculateCurrentUsage()
		if currentUsage+size > rm.maxMemoryUsage {
			return fmt.Errorf("memory limit exceeded: current=%d, requested=%d, limit=%d",
				currentUsage, size, rm.maxMemoryUsage)
		}
	}

	now := time.Now()
	resource := &Resource{
		ID:         id,
		Type:       resourceType,
		Size:       size,
		CreatedAt:  now,
		LastAccess: now,
		RefCount:   1,
		Data:       data,
	}

	rm.activeResources[id] = resource
	return nil
}

// GetResource retrieves a managed resource
func (rm *ResourceManager) GetResource(id string) (*Resource, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	resource, exists := rm.activeResources[id]
	if exists {
		resource.LastAccess = time.Now()
		resource.RefCount++
	}

	return resource, exists
}

// ReleaseResource releases a reference to a resource
func (rm *ResourceManager) ReleaseResource(id string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if resource, exists := rm.activeResources[id]; exists {
		resource.RefCount--
		if resource.RefCount <= 0 {
			delete(rm.activeResources, id)
		}
	}
}

// calculateCurrentUsage calculates current memory usage (must be called with lock held)
func (rm *ResourceManager) calculateCurrentUsage() int64 {
	var total int64
	for _, resource := range rm.activeResources {
		total += resource.Size
	}
	return total
}

// freeUnusedResources frees resources that haven't been accessed recently
func (rm *ResourceManager) freeUnusedResources() {
	now := time.Now()
	threshold := 5 * time.Minute // Free resources not accessed in 5 minutes

	for id, resource := range rm.activeResources {
		if resource.RefCount <= 0 && now.Sub(resource.LastAccess) > threshold {
			delete(rm.activeResources, id)
		}
	}
}

// cleanupRoutine runs periodic cleanup tasks
func (rm *ResourceManager) cleanupRoutine() {
	for {
		select {
		case <-rm.cleanupTicker.C:
			rm.performCleanup()
		case <-rm.stopCleanup:
			rm.cleanupTicker.Stop()
			return
		}
	}
}

// performCleanup performs periodic cleanup
func (rm *ResourceManager) performCleanup() {
	rm.mutex.Lock()

	// Free unused resources
	rm.freeUnusedResources()

	// Check if we should trigger GC
	currentUsage := rm.calculateCurrentUsage()
	shouldGC := currentUsage > rm.gcThreshold &&
		time.Since(rm.lastGCTime) > 2*time.Minute

	rm.mutex.Unlock()

	if shouldGC {
		runtime.GC()
		rm.lastGCTime = time.Now()
	}
}

// GetMemoryStats returns current memory statistics
func (rm *ResourceManager) GetMemoryStats() ResourceStats {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	stats := ResourceStats{
		ActiveResources: len(rm.activeResources),
		TotalSize:       rm.calculateCurrentUsage(),
		MaxSize:         rm.maxMemoryUsage,
	}

	// Get runtime memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	stats.RuntimeStats = RuntimeMemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}

	return stats
}

// ResourceStats provides resource usage statistics
type ResourceStats struct {
	ActiveResources int
	TotalSize       int64
	MaxSize         int64
	RuntimeStats    RuntimeMemoryStats
}

// RuntimeMemoryStats provides runtime memory statistics
type RuntimeMemoryStats struct {
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
}

// SetMemoryLimit sets the maximum memory usage limit
func (rm *ResourceManager) SetMemoryLimit(limit int64) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.maxMemoryUsage = limit
}

// SetGCThreshold sets the garbage collection threshold
func (rm *ResourceManager) SetGCThreshold(threshold int64) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.gcThreshold = threshold
}

// Close shuts down the resource manager
func (rm *ResourceManager) Close() error {
	close(rm.stopCleanup)

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Clear all resources
	rm.activeResources = make(map[string]*Resource)

	// Close underlying resource manager
	if rm.resourceManager != nil {
		rm.resourceManager.Close()
	}

	return nil
}

// WithResourceManager creates a context with resource management
func (rm *ResourceManager) WithResourceManager(ctx context.Context) context.Context {
	return context.WithValue(ctx, "resourceManager", rm)
}

// FromContext retrieves resource manager from context
func FromContext(ctx context.Context) (*ResourceManager, bool) {
	rm, ok := ctx.Value("resourceManager").(*ResourceManager)
	return rm, ok
}

// MemoryEfficientProcessor provides memory-efficient processing capabilities
type MemoryEfficientProcessor struct {
	resourceManager *ResourceManager
	batchSize       int
	maxConcurrency  int
}

// NewMemoryEfficientProcessor creates a new memory-efficient processor
func NewMemoryEfficientProcessor(rm *ResourceManager) *MemoryEfficientProcessor {
	return &MemoryEfficientProcessor{
		resourceManager: rm,
		batchSize:       100, // Process 100 items at a time
		maxConcurrency:  4,   // Limit concurrent operations
	}
}

// ProcessInBatches processes items in memory-efficient batches
func (mep *MemoryEfficientProcessor) ProcessInBatches(items []interface{}, processor func(item interface{}) error) error {
	for i := 0; i < len(items); i += mep.batchSize {
		end := i + mep.batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]

		// Process batch
		for _, item := range batch {
			if err := processor(item); err != nil {
				return fmt.Errorf("failed to process item in batch %d-%d: %w", i, end-1, err)
			}
		}

		// Force GC after each batch to keep memory usage low
		if i%500 == 0 { // Every 5 batches
			runtime.GC()
		}
	}

	return nil
}

// ProcessConcurrently processes items concurrently with memory limits
func (mep *MemoryEfficientProcessor) ProcessConcurrently(items []interface{}, processor func(item interface{}) error) error {
	semaphore := make(chan struct{}, mep.maxConcurrency)
	errChan := make(chan error, len(items))
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := processor(item); err != nil {
				errChan <- err
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
