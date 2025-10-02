// Package errors provides debug mode functionality with detailed tracing
package errors

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DebugManager manages debug mode functionality and detailed tracing
type DebugManager struct {
	enabled       bool
	traceLevel    TraceLevel
	logger        interfaces.Logger
	tracer        *DetailedTracer
	profiler      *PerformanceProfiler
	memoryTracker *MemoryTracker
	config        *DebugConfig
	mutex         sync.RWMutex
}

// DebugConfig contains configuration for debug mode
type DebugConfig struct {
	EnableTracing        bool          `json:"enable_tracing"`
	EnableProfiling      bool          `json:"enable_profiling"`
	EnableMemoryTracking bool          `json:"enable_memory_tracking"`
	TraceLevel           TraceLevel    `json:"trace_level"`
	MaxTraceDepth        int           `json:"max_trace_depth"`
	SampleRate           float64       `json:"sample_rate"`
	OutputPath           string        `json:"output_path"`
	FlushInterval        time.Duration `json:"flush_interval"`
}

// TraceLevel represents different levels of tracing detail
type TraceLevel int

const (
	TraceLevelNone TraceLevel = iota
	TraceLevelBasic
	TraceLevelDetailed
	TraceLevelVerbose
)

// String returns the string representation of trace level
func (tl TraceLevel) String() string {
	switch tl {
	case TraceLevelNone:
		return "none"
	case TraceLevelBasic:
		return "basic"
	case TraceLevelDetailed:
		return "detailed"
	case TraceLevelVerbose:
		return "verbose"
	default:
		return "unknown"
	}
}

// DefaultDebugConfig returns default debug configuration
func DefaultDebugConfig() *DebugConfig {
	homeDir, _ := os.UserHomeDir()
	return &DebugConfig{
		EnableTracing:        true,
		EnableProfiling:      true,
		EnableMemoryTracking: true,
		TraceLevel:           TraceLevelDetailed,
		MaxTraceDepth:        10,
		SampleRate:           1.0,
		OutputPath:           fmt.Sprintf("%s/.generator/debug", homeDir),
		FlushInterval:        5 * time.Second,
	}
}

// NewDebugManager creates a new debug manager
func NewDebugManager(config *DebugConfig, logger interfaces.Logger) *DebugManager {
	if config == nil {
		config = DefaultDebugConfig()
	}

	dm := &DebugManager{
		enabled:    true,
		traceLevel: config.TraceLevel,
		logger:     logger,
		config:     config,
	}

	// Initialize components based on configuration
	if config.EnableTracing {
		dm.tracer = NewDetailedTracer(config, logger)
	}

	if config.EnableProfiling {
		dm.profiler = NewPerformanceProfiler(config, logger)
	}

	if config.EnableMemoryTracking {
		dm.memoryTracker = NewMemoryTracker(config, logger)
	}

	return dm
}

// StartTrace starts a new trace for an operation
func (dm *DebugManager) StartTrace(ctx context.Context, operation string, details map[string]interface{}) context.Context {
	if !dm.enabled || dm.tracer == nil {
		return ctx
	}

	return dm.tracer.StartTrace(ctx, operation, details)
}

// EndTrace ends a trace for an operation
func (dm *DebugManager) EndTrace(ctx context.Context, result map[string]interface{}) {
	if !dm.enabled || dm.tracer == nil {
		return
	}

	dm.tracer.EndTrace(ctx, result)
}

// TraceStep adds a step to the current trace
func (dm *DebugManager) TraceStep(ctx context.Context, step string, details map[string]interface{}) {
	if !dm.enabled || dm.tracer == nil {
		return
	}

	dm.tracer.TraceStep(ctx, step, details)
}

// ProfileOperation profiles the performance of an operation
func (dm *DebugManager) ProfileOperation(operation string, fn func() error) error {
	if !dm.enabled || dm.profiler == nil {
		return fn()
	}

	return dm.profiler.ProfileOperation(operation, fn)
}

// TrackMemory tracks memory usage at a specific point
func (dm *DebugManager) TrackMemory(label string, details map[string]interface{}) {
	if !dm.enabled || dm.memoryTracker == nil {
		return
	}

	dm.memoryTracker.TrackMemory(label, details)
}

// GetDebugInfo returns comprehensive debug information
func (dm *DebugManager) GetDebugInfo() *DebugInfo {
	info := &DebugInfo{
		Enabled:    dm.enabled,
		TraceLevel: dm.traceLevel,
		Timestamp:  time.Now(),
	}

	if dm.tracer != nil {
		info.TracingInfo = dm.tracer.GetTracingInfo()
	}

	if dm.profiler != nil {
		info.ProfilingInfo = dm.profiler.GetProfilingInfo()
	}

	if dm.memoryTracker != nil {
		info.MemoryInfo = dm.memoryTracker.GetMemoryInfo()
	}

	return info
}

// DebugInfo contains comprehensive debug information
type DebugInfo struct {
	Enabled       bool           `json:"enabled"`
	TraceLevel    TraceLevel     `json:"trace_level"`
	Timestamp     time.Time      `json:"timestamp"`
	TracingInfo   *TracingInfo   `json:"tracing_info,omitempty"`
	ProfilingInfo *ProfilingInfo `json:"profiling_info,omitempty"`
	MemoryInfo    *MemoryInfo    `json:"memory_info,omitempty"`
}

// DetailedTracer provides detailed tracing functionality
type DetailedTracer struct {
	config       *DebugConfig
	logger       interfaces.Logger
	activeTraces map[string]*TraceContext
	mutex        sync.RWMutex
}

// TraceContext represents a single trace context
type TraceContext struct {
	ID        string                 `json:"id"`
	Operation string                 `json:"operation"`
	StartTime time.Time              `json:"start_time"`
	Steps     []TraceStep            `json:"steps"`
	Details   map[string]interface{} `json:"details"`
	Parent    *TraceContext          `json:"parent,omitempty"`
	Children  []*TraceContext        `json:"children,omitempty"`
	mutex     sync.RWMutex
}

// TraceStep represents a single step in a trace
type TraceStep struct {
	Timestamp time.Time              `json:"timestamp"`
	Step      string                 `json:"step"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Caller    *CallerInfo            `json:"caller,omitempty"`
}

// TracingInfo contains tracing statistics and information
type TracingInfo struct {
	ActiveTraces int            `json:"active_traces"`
	TotalTraces  int            `json:"total_traces"`
	AverageDepth float64        `json:"average_depth"`
	LongestTrace time.Duration  `json:"longest_trace"`
	TracesByOp   map[string]int `json:"traces_by_operation"`
}

// NewDetailedTracer creates a new detailed tracer
func NewDetailedTracer(config *DebugConfig, logger interfaces.Logger) *DetailedTracer {
	return &DetailedTracer{
		config:       config,
		logger:       logger,
		activeTraces: make(map[string]*TraceContext),
	}
}

// StartTrace starts a new trace
func (dt *DetailedTracer) StartTrace(ctx context.Context, operation string, details map[string]interface{}) context.Context {
	traceID := fmt.Sprintf("trace_%d", time.Now().UnixNano())

	trace := &TraceContext{
		ID:        traceID,
		Operation: operation,
		StartTime: time.Now(),
		Steps:     make([]TraceStep, 0),
		Details:   details,
	}

	// Check for parent trace in context
	if parentTrace := GetTraceFromContext(ctx); parentTrace != nil {
		trace.Parent = parentTrace
		parentTrace.mutex.Lock()
		parentTrace.Children = append(parentTrace.Children, trace)
		parentTrace.mutex.Unlock()
	}

	dt.mutex.Lock()
	dt.activeTraces[traceID] = trace
	dt.mutex.Unlock()

	// Add trace to context
	ctx = SetTraceInContext(ctx, trace)

	// Log trace start
	if dt.logger != nil {
		dt.logger.DebugWithFields("Trace started", map[string]interface{}{
			"trace_id":  traceID,
			"operation": operation,
			"details":   details,
		})
	}

	return ctx
}

// EndTrace ends a trace
func (dt *DetailedTracer) EndTrace(ctx context.Context, result map[string]interface{}) {
	trace := GetTraceFromContext(ctx)
	if trace == nil {
		return
	}

	duration := time.Since(trace.StartTime)

	// Add final step
	dt.TraceStep(ctx, "trace_completed", map[string]interface{}{
		"duration": duration,
		"result":   result,
	})

	// Remove from active traces
	dt.mutex.Lock()
	delete(dt.activeTraces, trace.ID)
	dt.mutex.Unlock()

	// Log trace completion
	if dt.logger != nil {
		dt.logger.DebugWithFields("Trace completed", map[string]interface{}{
			"trace_id":  trace.ID,
			"operation": trace.Operation,
			"duration":  duration,
			"steps":     len(trace.Steps),
			"result":    result,
		})
	}
}

// TraceStep adds a step to the current trace
func (dt *DetailedTracer) TraceStep(ctx context.Context, step string, details map[string]interface{}) {
	trace := GetTraceFromContext(ctx)
	if trace == nil {
		return
	}

	traceStep := TraceStep{
		Timestamp: time.Now(),
		Step:      step,
		Details:   details,
	}

	// Add caller information
	if file, line := GetCallerInfo(2); file != "" {
		traceStep.Caller = &CallerInfo{
			File: file,
			Line: line,
		}
	}

	trace.mutex.Lock()
	trace.Steps = append(trace.Steps, traceStep)
	trace.mutex.Unlock()

	// Log trace step
	if dt.logger != nil && dt.config.TraceLevel >= TraceLevelVerbose {
		dt.logger.DebugWithFields("Trace step", map[string]interface{}{
			"trace_id": trace.ID,
			"step":     step,
			"details":  details,
		})
	}
}

// GetTracingInfo returns tracing information
func (dt *DetailedTracer) GetTracingInfo() *TracingInfo {
	dt.mutex.RLock()
	defer dt.mutex.RUnlock()

	info := &TracingInfo{
		ActiveTraces: len(dt.activeTraces),
		TracesByOp:   make(map[string]int),
	}

	// Calculate statistics
	var totalDepth int
	for _, trace := range dt.activeTraces {
		info.TracesByOp[trace.Operation]++
		depth := dt.calculateTraceDepth(trace)
		totalDepth += depth

		duration := time.Since(trace.StartTime)
		if duration > info.LongestTrace {
			info.LongestTrace = duration
		}
	}

	if len(dt.activeTraces) > 0 {
		info.AverageDepth = float64(totalDepth) / float64(len(dt.activeTraces))
	}

	return info
}

// calculateTraceDepth calculates the depth of a trace
func (dt *DetailedTracer) calculateTraceDepth(trace *TraceContext) int {
	depth := 1
	for _, child := range trace.Children {
		childDepth := dt.calculateTraceDepth(child)
		if childDepth+1 > depth {
			depth = childDepth + 1
		}
	}
	return depth
}

// PerformanceProfiler provides performance profiling functionality
type PerformanceProfiler struct {
	config   *DebugConfig
	logger   interfaces.Logger
	profiles map[string]*PerformanceProfile
	mutex    sync.RWMutex
}

// PerformanceProfile represents performance data for an operation
type PerformanceProfile struct {
	Operation     string        `json:"operation"`
	Count         int           `json:"count"`
	TotalDuration time.Duration `json:"total_duration"`
	MinDuration   time.Duration `json:"min_duration"`
	MaxDuration   time.Duration `json:"max_duration"`
	AvgDuration   time.Duration `json:"avg_duration"`
	LastExecution time.Time     `json:"last_execution"`
}

// ProfilingInfo contains profiling statistics
type ProfilingInfo struct {
	TotalOperations int                            `json:"total_operations"`
	Profiles        map[string]*PerformanceProfile `json:"profiles"`
	SlowestOp       string                         `json:"slowest_operation"`
	FastestOp       string                         `json:"fastest_operation"`
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler(config *DebugConfig, logger interfaces.Logger) *PerformanceProfiler {
	return &PerformanceProfiler{
		config:   config,
		logger:   logger,
		profiles: make(map[string]*PerformanceProfile),
	}
}

// ProfileOperation profiles the execution of an operation
func (pp *PerformanceProfiler) ProfileOperation(operation string, fn func() error) error {
	startTime := time.Now()

	// Execute the operation
	err := fn()

	duration := time.Since(startTime)

	// Update profile
	pp.updateProfile(operation, duration)

	// Log performance data
	if pp.logger != nil {
		pp.logger.DebugWithFields("Operation profiled", map[string]interface{}{
			"operation": operation,
			"duration":  duration,
			"success":   err == nil,
		})
	}

	return err
}

// updateProfile updates the performance profile for an operation
func (pp *PerformanceProfiler) updateProfile(operation string, duration time.Duration) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	profile, exists := pp.profiles[operation]
	if !exists {
		profile = &PerformanceProfile{
			Operation:   operation,
			MinDuration: duration,
			MaxDuration: duration,
		}
		pp.profiles[operation] = profile
	}

	profile.Count++
	profile.TotalDuration += duration
	profile.AvgDuration = profile.TotalDuration / time.Duration(profile.Count)
	profile.LastExecution = time.Now()

	if duration < profile.MinDuration {
		profile.MinDuration = duration
	}
	if duration > profile.MaxDuration {
		profile.MaxDuration = duration
	}
}

// GetProfilingInfo returns profiling information
func (pp *PerformanceProfiler) GetProfilingInfo() *ProfilingInfo {
	pp.mutex.RLock()
	defer pp.mutex.RUnlock()

	info := &ProfilingInfo{
		Profiles: make(map[string]*PerformanceProfile),
	}

	var slowestDuration, fastestDuration time.Duration
	var slowestOp, fastestOp string
	first := true

	for op, profile := range pp.profiles {
		info.TotalOperations += profile.Count
		info.Profiles[op] = profile

		if first || profile.AvgDuration > slowestDuration {
			slowestDuration = profile.AvgDuration
			slowestOp = op
		}
		if first || profile.AvgDuration < fastestDuration {
			fastestDuration = profile.AvgDuration
			fastestOp = op
		}
		first = false
	}

	info.SlowestOp = slowestOp
	info.FastestOp = fastestOp

	return info
}

// MemoryTracker tracks memory usage patterns
type MemoryTracker struct {
	config    *DebugConfig
	logger    interfaces.Logger
	snapshots []MemorySnapshot
	mutex     sync.RWMutex
}

// MemorySnapshot represents a memory usage snapshot
type MemorySnapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Label     string                 `json:"label"`
	Stats     runtime.MemStats       `json:"stats"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// MemoryInfo contains memory tracking information
type MemoryInfo struct {
	CurrentUsage uint64           `json:"current_usage"`
	PeakUsage    uint64           `json:"peak_usage"`
	TotalAllocs  uint64           `json:"total_allocs"`
	GCCount      uint32           `json:"gc_count"`
	Snapshots    []MemorySnapshot `json:"snapshots"`
}

// NewMemoryTracker creates a new memory tracker
func NewMemoryTracker(config *DebugConfig, logger interfaces.Logger) *MemoryTracker {
	return &MemoryTracker{
		config:    config,
		logger:    logger,
		snapshots: make([]MemorySnapshot, 0),
	}
}

// TrackMemory takes a memory usage snapshot
func (mt *MemoryTracker) TrackMemory(label string, details map[string]interface{}) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	snapshot := MemorySnapshot{
		Timestamp: time.Now(),
		Label:     label,
		Stats:     stats,
		Details:   details,
	}

	mt.mutex.Lock()
	mt.snapshots = append(mt.snapshots, snapshot)

	// Keep only the last 100 snapshots
	if len(mt.snapshots) > 100 {
		mt.snapshots = mt.snapshots[1:]
	}
	mt.mutex.Unlock()

	// Log memory usage
	if mt.logger != nil {
		mt.logger.DebugWithFields("Memory tracked", map[string]interface{}{
			"label":       label,
			"alloc":       stats.Alloc,
			"total_alloc": stats.TotalAlloc,
			"sys":         stats.Sys,
			"num_gc":      stats.NumGC,
			"details":     details,
		})
	}
}

// GetMemoryInfo returns memory tracking information
func (mt *MemoryTracker) GetMemoryInfo() *MemoryInfo {
	mt.mutex.RLock()
	defer mt.mutex.RUnlock()

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	info := &MemoryInfo{
		CurrentUsage: stats.Alloc,
		TotalAllocs:  stats.TotalAlloc,
		GCCount:      stats.NumGC,
		Snapshots:    make([]MemorySnapshot, len(mt.snapshots)),
	}

	copy(info.Snapshots, mt.snapshots)

	// Find peak usage
	for _, snapshot := range mt.snapshots {
		if snapshot.Stats.Alloc > info.PeakUsage {
			info.PeakUsage = snapshot.Stats.Alloc
		}
	}

	return info
}

// Context key for trace context
type traceContextKey struct{}

// SetTraceInContext sets a trace in the context
func SetTraceInContext(ctx context.Context, trace *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey{}, trace)
}

// GetTraceFromContext gets a trace from the context
func GetTraceFromContext(ctx context.Context) *TraceContext {
	if trace, ok := ctx.Value(traceContextKey{}).(*TraceContext); ok {
		return trace
	}
	return nil
}

// Close closes the debug manager and releases resources
func (dm *DebugManager) Close() error {
	// Generate final debug report
	if dm.logger != nil {
		info := dm.GetDebugInfo()
		dm.logger.InfoWithFields("Debug session completed", map[string]interface{}{
			"debug_info": info,
		})
	}

	return nil
}
