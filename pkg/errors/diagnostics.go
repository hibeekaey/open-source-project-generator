// Package errors provides comprehensive diagnostic collection for error analysis
package errors

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DiagnosticCollector collects comprehensive diagnostic information for error analysis
type DiagnosticCollector struct {
	config          *EnhancedErrorConfig
	logger          interfaces.Logger
	systemCollector *SystemInfoCollector
	envCollector    *EnvironmentCollector
	logCollector    *LogCollector
	mutex           sync.RWMutex
	stats           *DiagnosticStatistics
}

// DiagnosticInfo contains comprehensive diagnostic information
type DiagnosticInfo struct {
	CollectedAt     time.Time              `json:"collected_at"`
	Operation       string                 `json:"operation"`
	SystemInfo      *DiagnosticSystemInfo  `json:"system_info,omitempty"`
	EnvironmentInfo *EnvironmentInfo       `json:"environment_info,omitempty"`
	ProcessInfo     *ProcessInfo           `json:"process_info,omitempty"`
	RelevantLogs    []LogEntry             `json:"relevant_logs,omitempty"`
	StackTrace      string                 `json:"stack_trace,omitempty"`
	Context         map[string]interface{} `json:"context,omitempty"`
	Metrics         *DiagnosticMetrics     `json:"metrics,omitempty"`
}

// DiagnosticSystemInfo contains system information for diagnostics
type DiagnosticSystemInfo struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	ProcessID   int    `json:"process_id"`
	ParentPID   int    `json:"parent_pid"`
	MemoryUsage int64  `json:"memory_usage"`
	DiskSpace   int64  `json:"disk_space"`
	CPUCount    int    `json:"cpu_count"`
}

// DiagnosticMetrics contains performance and resource metrics
type DiagnosticMetrics struct {
	MemoryUsage    int64         `json:"memory_usage"`
	CPUUsage       float64       `json:"cpu_usage"`
	DiskUsage      int64         `json:"disk_usage"`
	NetworkLatency time.Duration `json:"network_latency"`
	FileHandles    int           `json:"file_handles"`
	Goroutines     int           `json:"goroutines"`
}

// ProcessInfo contains information about the current process
type ProcessInfo struct {
	PID         int               `json:"pid"`
	ParentPID   int               `json:"parent_pid"`
	CommandLine []string          `json:"command_line"`
	WorkingDir  string            `json:"working_dir"`
	Environment map[string]string `json:"environment"`
	StartTime   time.Time         `json:"start_time"`
	Runtime     time.Duration     `json:"runtime"`
	MemoryStats *MemoryStats      `json:"memory_stats,omitempty"`
}

// MemoryStats contains detailed memory statistics
type MemoryStats struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	Lookups      uint64 `json:"lookups"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"frees"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`
	StackInuse   uint64 `json:"stack_inuse"`
	StackSys     uint64 `json:"stack_sys"`
	MSpanInuse   uint64 `json:"mspan_inuse"`
	MSpanSys     uint64 `json:"mspan_sys"`
	MCacheInuse  uint64 `json:"mcache_inuse"`
	MCacheSys    uint64 `json:"mcache_sys"`
	GCSys        uint64 `json:"gc_sys"`
	OtherSys     uint64 `json:"other_sys"`
	NextGC       uint64 `json:"next_gc"`
	LastGC       uint64 `json:"last_gc"`
	PauseTotalNs uint64 `json:"pause_total_ns"`
	NumGC        uint32 `json:"num_gc"`
	NumForcedGC  uint32 `json:"num_forced_gc"`
}

// DiagnosticStatistics tracks diagnostic collection statistics
type DiagnosticStatistics struct {
	TotalCollections      int            `json:"total_collections"`
	CollectionsByType     map[string]int `json:"collections_by_type"`
	AverageCollectionTime time.Duration  `json:"average_collection_time"`
	LastCollection        *time.Time     `json:"last_collection,omitempty"`
	FailedCollections     int            `json:"failed_collections"`
	CollectionErrors      []string       `json:"collection_errors"`
}

// NewDiagnosticCollector creates a new diagnostic collector
func NewDiagnosticCollector(config *EnhancedErrorConfig, logger interfaces.Logger) *DiagnosticCollector {
	return &DiagnosticCollector{
		config:          config,
		logger:          logger,
		systemCollector: NewSystemInfoCollector(),
		envCollector:    NewEnvironmentCollector(),
		logCollector:    NewLogCollector(logger),
		stats: &DiagnosticStatistics{
			CollectionsByType: make(map[string]int),
			CollectionErrors:  make([]string, 0),
		},
	}
}

// CollectDiagnostics collects comprehensive diagnostic information
func (dc *DiagnosticCollector) CollectDiagnostics(ctx context.Context, err error, operation string, context map[string]interface{}) *DiagnosticInfo {
	startTime := time.Now()

	dc.mutex.Lock()
	dc.stats.TotalCollections++
	dc.stats.CollectionsByType[operation]++
	now := time.Now()
	dc.stats.LastCollection = &now
	dc.mutex.Unlock()

	info := &DiagnosticInfo{
		CollectedAt: startTime,
		Operation:   operation,
		Context:     context,
	}

	// Collect system information if enabled
	if dc.config.CollectSystemInfo {
		if systemInfo, sysErr := dc.systemCollector.Collect(); sysErr == nil {
			info.SystemInfo = systemInfo
		} else {
			dc.recordError(fmt.Sprintf("Failed to collect system info: %v", sysErr))
		}
	}

	// Collect environment information if enabled
	if dc.config.CollectEnvironment {
		if envInfo, envErr := dc.envCollector.Collect(); envErr == nil {
			info.EnvironmentInfo = envInfo
		} else {
			dc.recordError(fmt.Sprintf("Failed to collect environment info: %v", envErr))
		}
	}

	// Collect process information
	if processInfo, procErr := dc.collectProcessInfo(); procErr == nil {
		info.ProcessInfo = processInfo
	} else {
		dc.recordError(fmt.Sprintf("Failed to collect process info: %v", procErr))
	}

	// Collect stack trace if enabled
	if dc.config.CollectStackTraces {
		info.StackTrace = dc.collectStackTrace()
	}

	// Collect relevant logs
	if logs, logErr := dc.logCollector.CollectRelevantLogs(operation, 10); logErr == nil {
		info.RelevantLogs = logs
	} else {
		dc.recordError(fmt.Sprintf("Failed to collect logs: %v", logErr))
	}

	// Collect performance metrics
	if metrics, metricsErr := dc.collectMetrics(); metricsErr == nil {
		info.Metrics = metrics
	} else {
		dc.recordError(fmt.Sprintf("Failed to collect metrics: %v", metricsErr))
	}

	// Update statistics
	duration := time.Since(startTime)
	dc.mutex.Lock()
	if dc.stats.TotalCollections > 0 {
		dc.stats.AverageCollectionTime = (dc.stats.AverageCollectionTime*time.Duration(dc.stats.TotalCollections-1) + duration) / time.Duration(dc.stats.TotalCollections)
	} else {
		dc.stats.AverageCollectionTime = duration
	}
	dc.mutex.Unlock()

	return info
}

// collectProcessInfo collects information about the current process
func (dc *DiagnosticCollector) collectProcessInfo() (*ProcessInfo, error) {
	info := &ProcessInfo{
		PID:       os.Getpid(),
		ParentPID: os.Getppid(),
		StartTime: time.Now(), // Approximation - would need process start time
	}

	// Get command line arguments
	info.CommandLine = os.Args

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		info.WorkingDir = wd
	}

	// Get environment variables (filtered for security)
	info.Environment = dc.getFilteredEnvironment()

	// Get memory statistics
	if memStats, err := dc.collectMemoryStats(); err == nil {
		info.MemoryStats = memStats
	}

	return info, nil
}

// getFilteredEnvironment returns environment variables filtered for security
func (dc *DiagnosticCollector) getFilteredEnvironment() map[string]string {
	env := make(map[string]string)

	// Safe environment variables to include
	safeVars := []string{
		"PATH", "HOME", "USER", "SHELL", "TERM", "LANG", "LC_ALL",
		"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "CIRCLECI", "TRAVIS",
		"GO_VERSION", "GOPATH", "GOROOT", "GOOS", "GOARCH",
		"NODE_VERSION", "NPM_VERSION", "PYTHON_VERSION",
		"GENERATOR_", // Our own environment variables
	}

	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Check if this is a safe variable
		for _, safe := range safeVars {
			if key == safe || strings.HasPrefix(key, safe) {
				env[key] = value
				break
			}
		}
	}

	return env
}

// collectMemoryStats collects detailed memory statistics
func (dc *DiagnosticCollector) collectMemoryStats() (*MemoryStats, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &MemoryStats{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		Lookups:      m.Lookups,
		Mallocs:      m.Mallocs,
		Frees:        m.Frees,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
		StackInuse:   m.StackInuse,
		StackSys:     m.StackSys,
		MSpanInuse:   m.MSpanInuse,
		MSpanSys:     m.MSpanSys,
		MCacheInuse:  m.MCacheInuse,
		MCacheSys:    m.MCacheSys,
		GCSys:        m.GCSys,
		OtherSys:     m.OtherSys,
		NextGC:       m.NextGC,
		LastGC:       m.LastGC,
		PauseTotalNs: m.PauseTotalNs,
		NumGC:        m.NumGC,
		NumForcedGC:  m.NumForcedGC,
	}, nil
}

// collectStackTrace collects the current stack trace
func (dc *DiagnosticCollector) collectStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// collectMetrics collects performance and resource metrics
func (dc *DiagnosticCollector) collectMetrics() (*DiagnosticMetrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := &DiagnosticMetrics{
		MemoryUsage: int64(m.Alloc & 0x7FFFFFFFFFFFFFFF), // Ensure positive int64
		Goroutines:  runtime.NumGoroutine(),
	}

	// Collect additional metrics if available
	if diskUsage, err := dc.getDiskUsage(); err == nil {
		metrics.DiskUsage = diskUsage
	}

	return metrics, nil
}

// getDiskUsage gets disk usage for the current working directory
func (dc *DiagnosticCollector) getDiskUsage() (int64, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	// This is a simplified implementation
	// In a real implementation, you'd use syscalls to get actual disk usage
	if stat, err := os.Stat(wd); err == nil {
		return stat.Size(), nil
	}

	return 0, fmt.Errorf("unable to get disk usage")
}

// recordError records a diagnostic collection error
func (dc *DiagnosticCollector) recordError(errorMsg string) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	dc.stats.FailedCollections++
	dc.stats.CollectionErrors = append(dc.stats.CollectionErrors, errorMsg)

	// Keep only the last 10 errors
	if len(dc.stats.CollectionErrors) > 10 {
		dc.stats.CollectionErrors = dc.stats.CollectionErrors[1:]
	}

	if dc.logger != nil {
		dc.logger.Warn("Diagnostic collection error: %s", errorMsg)
	}
}

// GetStatistics returns diagnostic collection statistics
func (dc *DiagnosticCollector) GetStatistics() *DiagnosticStatistics {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	// Return a copy to avoid race conditions
	stats := *dc.stats
	return &stats
}

// SetVerboseMode sets verbose mode for diagnostic collection
func (dc *DiagnosticCollector) SetVerboseMode(verbose bool) {
	// Verbose mode affects what diagnostics are collected
	if verbose {
		dc.config.CollectStackTraces = true
		dc.config.CollectSystemInfo = true
	}
}

// Close closes the diagnostic collector and releases resources
func (dc *DiagnosticCollector) Close() error {
	// Clean up any resources
	return nil
}

// SystemInfoCollector collects system information
type SystemInfoCollector struct{}

// NewSystemInfoCollector creates a new system info collector
func NewSystemInfoCollector() *SystemInfoCollector {
	return &SystemInfoCollector{}
}

// Collect collects system information
func (sic *SystemInfoCollector) Collect() (*DiagnosticSystemInfo, error) {
	info := &DiagnosticSystemInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		ProcessID: os.Getpid(),
		ParentPID: os.Getppid(),
	}

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		info.Hostname = hostname
	}

	// Get username
	if user := os.Getenv("USER"); user != "" {
		info.Username = user
	} else if user := os.Getenv("USERNAME"); user != "" {
		info.Username = user
	}

	// Get CPU count
	info.CPUCount = runtime.NumCPU()

	// Get memory usage (simplified)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info.MemoryUsage = int64(m.Alloc & 0x7FFFFFFFFFFFFFFF) // Ensure positive int64

	return info, nil
}

// EnvironmentCollector collects environment information
type EnvironmentCollector struct{}

// NewEnvironmentCollector creates a new environment collector
func NewEnvironmentCollector() *EnvironmentCollector {
	return &EnvironmentCollector{}
}

// Collect collects environment information
func (ec *EnvironmentCollector) Collect() (*EnvironmentInfo, error) {
	info := &EnvironmentInfo{
		Environment: make(map[string]string),
	}

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		info.WorkingDir = wd
	}

	// Collect relevant environment variables
	relevantVars := []string{
		"PATH", "HOME", "USER", "SHELL", "TERM",
		"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL",
		"GO_VERSION", "NODE_VERSION", "PYTHON_VERSION",
	}

	for _, envVar := range relevantVars {
		if value := os.Getenv(envVar); value != "" {
			info.Environment[envVar] = value
		}
	}

	// Detect CI environment
	info.CI = ec.detectCIEnvironment()

	return info, nil
}

// detectCIEnvironment detects CI environment information
func (ec *EnvironmentCollector) detectCIEnvironment() *CIEnvironment {
	ci := &CIEnvironment{}

	// Check common CI environment variables
	if os.Getenv("CI") == "true" || os.Getenv("CONTINUOUS_INTEGRATION") == "true" {
		ci.IsCI = true
	}

	// Detect specific CI providers
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ci.Provider = "github"
		ci.JobID = os.Getenv("GITHUB_RUN_ID")
		ci.BuildID = os.Getenv("GITHUB_RUN_NUMBER")
	} else if os.Getenv("GITLAB_CI") == "true" {
		ci.Provider = "gitlab"
		ci.JobID = os.Getenv("CI_JOB_ID")
		ci.BuildID = os.Getenv("CI_PIPELINE_ID")
	} else if os.Getenv("JENKINS_URL") != "" {
		ci.Provider = "jenkins"
		ci.JobID = os.Getenv("BUILD_ID")
		ci.BuildID = os.Getenv("BUILD_NUMBER")
	}

	return ci
}

// LogCollector collects relevant log entries
type LogCollector struct {
	logger interfaces.Logger
}

// NewLogCollector creates a new log collector
func NewLogCollector(logger interfaces.Logger) *LogCollector {
	return &LogCollector{
		logger: logger,
	}
}

// CollectRelevantLogs collects recent log entries relevant to the operation
func (lc *LogCollector) CollectRelevantLogs(operation string, limit int) ([]LogEntry, error) {
	// This is a simplified implementation
	// In a real implementation, you'd read from log files or log buffers

	logs := make([]LogEntry, 0, limit)

	// For now, return empty logs
	// This would be implemented to read from actual log sources

	return logs, nil
}
