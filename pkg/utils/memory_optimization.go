package utils

import (
	"runtime"
	"sync"
	"time"
)

// MemoryPool provides pooled byte slices to reduce allocations
type MemoryPool struct {
	pools map[int]*sync.Pool
	mutex sync.RWMutex
}

// NewMemoryPool creates a new memory pool
func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		pools: make(map[int]*sync.Pool),
	}
}

// Get retrieves a byte slice of the specified size from the pool
func (mp *MemoryPool) Get(size int) []byte {
	// Round up to the nearest power of 2 for better pooling
	poolSize := nextPowerOf2(size)

	mp.mutex.RLock()
	pool, exists := mp.pools[poolSize]
	mp.mutex.RUnlock()

	if !exists {
		mp.mutex.Lock()
		// Double-check after acquiring write lock
		if pool, exists = mp.pools[poolSize]; !exists {
			pool = &sync.Pool{
				New: func() interface{} {
					return make([]byte, poolSize)
				},
			}
			mp.pools[poolSize] = pool
		}
		mp.mutex.Unlock()
	}

	buf := pool.Get().([]byte)
	return buf[:size] // Return slice with requested size
}

// Put returns a byte slice to the pool
func (mp *MemoryPool) Put(buf []byte) {
	if cap(buf) == 0 {
		return
	}

	poolSize := cap(buf)

	mp.mutex.RLock()
	pool, exists := mp.pools[poolSize]
	mp.mutex.RUnlock()

	if exists {
		// Clear the slice before returning to pool
		for i := range buf[:cap(buf)] {
			buf[i] = 0
		}
		pool.Put(buf[:cap(buf)])
	}
}

// nextPowerOf2 returns the next power of 2 greater than or equal to n
func nextPowerOf2(n int) int {
	if n <= 0 {
		return 1
	}

	// Handle small sizes with fixed buckets
	if n <= 64 {
		return 64
	}
	if n <= 128 {
		return 128
	}
	if n <= 256 {
		return 256
	}
	if n <= 512 {
		return 512
	}
	if n <= 1024 {
		return 1024
	}
	if n <= 2048 {
		return 2048
	}
	if n <= 4096 {
		return 4096
	}
	if n <= 8192 {
		return 8192
	}
	if n <= 16384 {
		return 16384
	}
	if n <= 32768 {
		return 32768
	}
	if n <= 65536 {
		return 65536
	}

	// For larger sizes, use actual power of 2 calculation
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++

	return n
}

// ResourceManager manages memory and other resources efficiently
type ResourceManager struct {
	memoryPool    *MemoryPool
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
	gcThreshold   int64 // Memory threshold in bytes to trigger GC
	lastGCTime    time.Time
	gcInterval    time.Duration
}

// NewResourceManager creates a new resource manager
func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		memoryPool:  NewMemoryPool(),
		stopCleanup: make(chan struct{}),
		gcThreshold: 100 * 1024 * 1024, // 100MB threshold
		gcInterval:  5 * time.Minute,   // Minimum 5 minutes between forced GC
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

// ForceGC forces garbage collection if conditions are met
func (rm *ResourceManager) ForceGC() {
	now := time.Now()
	if now.Sub(rm.lastGCTime) > rm.gcInterval {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if int64(m.Alloc) > rm.gcThreshold {
			runtime.GC()
			rm.lastGCTime = now
		}
	}
}

// cleanupRoutine runs periodic cleanup tasks
func (rm *ResourceManager) cleanupRoutine() {
	for {
		select {
		case <-rm.cleanupTicker.C:
			rm.ForceGC()
		case <-rm.stopCleanup:
			rm.cleanupTicker.Stop()
			return
		}
	}
}

// Close stops the resource manager
func (rm *ResourceManager) Close() {
	close(rm.stopCleanup)
}

// MemoryStats provides memory usage statistics
type MemoryStats struct {
	Alloc         uint64  // Bytes allocated and not yet freed
	TotalAlloc    uint64  // Total bytes allocated (even if freed)
	Sys           uint64  // Bytes obtained from system
	Lookups       uint64  // Number of pointer lookups
	Mallocs       uint64  // Number of mallocs
	Frees         uint64  // Number of frees
	HeapAlloc     uint64  // Bytes allocated and not yet freed (same as Alloc)
	HeapSys       uint64  // Bytes obtained from system
	HeapIdle      uint64  // Bytes in idle spans
	HeapInuse     uint64  // Bytes in non-idle span
	HeapReleased  uint64  // Bytes released to the OS
	HeapObjects   uint64  // Total number of allocated objects
	GCCPUFraction float64 // Fraction of CPU time used by GC
	NumGC         uint32  // Number of completed GC cycles
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:         m.Alloc,
		TotalAlloc:    m.TotalAlloc,
		Sys:           m.Sys,
		Lookups:       m.Lookups,
		Mallocs:       m.Mallocs,
		Frees:         m.Frees,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		HeapIdle:      m.HeapIdle,
		HeapInuse:     m.HeapInuse,
		HeapReleased:  m.HeapReleased,
		HeapObjects:   m.HeapObjects,
		GCCPUFraction: m.GCCPUFraction,
		NumGC:         m.NumGC,
	}
}

// ObjectPool provides generic object pooling
type ObjectPool struct {
	pool sync.Pool
	new  func() interface{}
}

// NewObjectPool creates a new object pool with the given constructor function
func NewObjectPool(newFunc func() interface{}) *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{New: newFunc},
		new:  newFunc,
	}
}

// Get retrieves an object from the pool
func (op *ObjectPool) Get() interface{} {
	return op.pool.Get()
}

// Put returns an object to the pool
func (op *ObjectPool) Put(obj interface{}) {
	op.pool.Put(obj)
}

// Global resource manager instance
var globalResourceManager = NewResourceManager()

// GetGlobalBuffer gets a buffer from the global resource manager
func GetGlobalBuffer(size int) []byte {
	return globalResourceManager.GetBuffer(size)
}

// PutGlobalBuffer returns a buffer to the global resource manager
func PutGlobalBuffer(buf []byte) {
	globalResourceManager.PutBuffer(buf)
}

// ForceGlobalGC forces garbage collection using the global resource manager
func ForceGlobalGC() {
	globalResourceManager.ForceGC()
}

// CloseGlobalResourceManager closes the global resource manager
func CloseGlobalResourceManager() {
	globalResourceManager.Close()
}
