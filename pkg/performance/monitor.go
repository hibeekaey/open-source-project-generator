// Package performance provides system health monitoring and diagnostics capabilities
package performance

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// SystemMonitor provides comprehensive system health monitoring
type SystemMonitor struct {
	cache           interfaces.CacheManager
	logger          interfaces.Logger
	metrics         *MetricsCollector
	healthCheckers  map[string]HealthChecker
	diagnostics     *DiagnosticsCollector
	monitoring      bool
	monitorInterval time.Duration
	mutex           sync.RWMutex
	stopChan        chan struct{}
	healthHistory   []*SystemHealthSnapshot
	maxHistorySize  int
}

// HealthChecker defines interface for component health checking
type HealthChecker interface {
	CheckHealth(ctx context.Context) *ComponentHealth
	GetComponentName() string
}

// ComponentHealth represents the health status of a system component
type ComponentHealth struct {
	ComponentName   string                 `json:"component_name"`
	Status          string                 `json:"status"` // "healthy", "degraded", "unhealthy"
	LastCheck       time.Time              `json:"last_check"`
	ResponseTime    time.Duration          `json:"response_time"`
	ErrorCount      int                    `json:"error_count"`
	WarningCount    int                    `json:"warning_count"`
	Details         map[string]interface{} `json:"details"`
	Issues          []string               `json:"issues"`
	Recommendations []string               `json:"recommendations"`
}

// SystemHealthSnapshot represents a point-in-time system health status
type SystemHealthSnapshot struct {
	Timestamp       time.Time                   `json:"timestamp"`
	OverallStatus   string                      `json:"overall_status"`
	Components      map[string]*ComponentHealth `json:"components"`
	SystemMetrics   *SystemMetrics              `json:"system_metrics"`
	PerformanceData *MetricsReport              `json:"performance_data"`
	Alerts          []HealthAlert               `json:"alerts"`
}

// HealthAlert represents a system health alert
type HealthAlert struct {
	ID         string                 `json:"id"`
	Severity   string                 `json:"severity"` // "info", "warning", "error", "critical"
	Component  string                 `json:"component"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Resolved   bool                   `json:"resolved"`
	ResolvedAt *time.Time             `json:"resolved_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// DiagnosticsCollector collects diagnostic information
type DiagnosticsCollector struct {
	systemInfo     *SystemInfo
	errorHistory   []ErrorRecord
	performanceLog []PerformanceRecord
	mutex          sync.RWMutex
	maxRecords     int
}

// SystemInfo contains static system information
type SystemInfo struct {
	OS             string    `json:"os"`
	Architecture   string    `json:"architecture"`
	GoVersion      string    `json:"go_version"`
	NumCPU         int       `json:"num_cpu"`
	StartTime      time.Time `json:"start_time"`
	WorkingDir     string    `json:"working_dir"`
	ExecutablePath string    `json:"executable_path"`
}

// RuntimeMetrics contains runtime performance metrics
type RuntimeMetrics struct {
	Timestamp     time.Time `json:"timestamp"`
	MemoryAlloc   uint64    `json:"memory_alloc"`
	MemoryTotal   uint64    `json:"memory_total"`
	MemorySys     uint64    `json:"memory_sys"`
	NumGoroutines int       `json:"num_goroutines"`
	NumGC         uint32    `json:"num_gc"`
	GCPauseTotal  uint64    `json:"gc_pause_total"`
	HeapObjects   uint64    `json:"heap_objects"`
	StackInUse    uint64    `json:"stack_in_use"`
}

// ErrorRecord represents an error occurrence
type ErrorRecord struct {
	Timestamp time.Time              `json:"timestamp"`
	Component string                 `json:"component"`
	Operation string                 `json:"operation"`
	ErrorType string                 `json:"error_type"`
	Message   string                 `json:"message"`
	Severity  string                 `json:"severity"`
	Context   map[string]interface{} `json:"context"`
	Resolved  bool                   `json:"resolved"`
}

// PerformanceRecord represents a performance measurement
type PerformanceRecord struct {
	Timestamp  time.Time              `json:"timestamp"`
	Operation  string                 `json:"operation"`
	Duration   time.Duration          `json:"duration"`
	MemoryUsed int64                  `json:"memory_used"`
	Success    bool                   `json:"success"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// NewSystemMonitor creates a new system monitor instance
func NewSystemMonitor(cache interfaces.CacheManager, logger interfaces.Logger) *SystemMonitor {
	monitor := &SystemMonitor{
		cache:           cache,
		logger:          logger,
		metrics:         NewMetricsCollector(),
		healthCheckers:  make(map[string]HealthChecker),
		diagnostics:     NewDiagnosticsCollector(),
		monitorInterval: 30 * time.Second,
		stopChan:        make(chan struct{}),
		healthHistory:   make([]*SystemHealthSnapshot, 0),
		maxHistorySize:  100,
	}

	// Register default health checkers
	monitor.registerDefaultHealthCheckers()

	return monitor
}

// NewDiagnosticsCollector creates a new diagnostics collector
func NewDiagnosticsCollector() *DiagnosticsCollector {
	return &DiagnosticsCollector{
		systemInfo:     collectSystemInfo(),
		errorHistory:   make([]ErrorRecord, 0),
		performanceLog: make([]PerformanceRecord, 0),
		maxRecords:     1000,
	}
}

// collectSystemInfo collects static system information
func collectSystemInfo() *SystemInfo {
	return &SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		StartTime:    time.Now(),
	}
}

// registerDefaultHealthCheckers registers built-in health checkers
func (sm *SystemMonitor) registerDefaultHealthCheckers() {
	// Cache health checker
	if sm.cache != nil {
		sm.RegisterHealthChecker("cache", &CacheHealthChecker{cache: sm.cache})
	}

	// Memory health checker
	sm.RegisterHealthChecker("memory", &MemoryHealthChecker{})

	// Performance health checker
	sm.RegisterHealthChecker("performance", &PerformanceHealthChecker{metrics: sm.metrics})
}

// RegisterHealthChecker registers a health checker for a component
func (sm *SystemMonitor) RegisterHealthChecker(name string, checker HealthChecker) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.healthCheckers[name] = checker
}

// StartMonitoring starts continuous system monitoring
func (sm *SystemMonitor) StartMonitoring(ctx context.Context) error {
	sm.mutex.Lock()
	if sm.monitoring {
		sm.mutex.Unlock()
		return fmt.Errorf("monitoring already started")
	}
	sm.monitoring = true
	sm.mutex.Unlock()

	go sm.monitoringLoop(ctx)

	if sm.logger != nil {
		sm.logger.Info("System monitoring started", "interval", sm.monitorInterval)
	}

	return nil
}

// StopMonitoring stops system monitoring
func (sm *SystemMonitor) StopMonitoring() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if !sm.monitoring {
		return fmt.Errorf("monitoring not started")
	}

	close(sm.stopChan)
	sm.monitoring = false

	if sm.logger != nil {
		sm.logger.Info("System monitoring stopped")
	}

	return nil
}

// monitoringLoop runs the continuous monitoring loop
func (sm *SystemMonitor) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(sm.monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-sm.stopChan:
			return
		case <-ticker.C:
			snapshot := sm.collectHealthSnapshot(ctx)
			sm.addHealthSnapshot(snapshot)
			sm.processAlerts(snapshot)
		}
	}
}

// collectHealthSnapshot collects a complete system health snapshot
func (sm *SystemMonitor) collectHealthSnapshot(ctx context.Context) *SystemHealthSnapshot {
	sm.mutex.RLock()
	checkers := make(map[string]HealthChecker)
	for k, v := range sm.healthCheckers {
		checkers[k] = v
	}
	sm.mutex.RUnlock()

	snapshot := &SystemHealthSnapshot{
		Timestamp:       time.Now(),
		Components:      make(map[string]*ComponentHealth),
		SystemMetrics:   sm.collectSystemMetrics(),
		PerformanceData: sm.metrics.GenerateReport(),
		Alerts:          make([]HealthAlert, 0),
	}

	// Check health of all registered components
	var healthyCount, degradedCount, unhealthyCount int

	for name, checker := range checkers {
		health := checker.CheckHealth(ctx)
		snapshot.Components[name] = health

		switch health.Status {
		case "healthy":
			healthyCount++
		case "degraded":
			degradedCount++
		case "unhealthy":
			unhealthyCount++
		}
	}

	// Determine overall status
	if unhealthyCount > 0 {
		snapshot.OverallStatus = "unhealthy"
	} else if degradedCount > 0 {
		snapshot.OverallStatus = "degraded"
	} else {
		snapshot.OverallStatus = "healthy"
	}

	return snapshot
}

// collectSystemMetrics collects current system metrics
func (sm *SystemMonitor) collectSystemMetrics() *SystemMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &SystemMetrics{
		CPUUsage:       getCPUUsage(),
		MemoryUsage:    safeUint64ToInt64(memStats.Alloc),
		DiskUsage:      getDiskUsage(),
		NetworkLatency: getNetworkLatency(),
		CacheHitRate:   sm.getCacheHitRate(),
		LastUpdated:    time.Now(),
		CustomMetrics:  make(map[string]interface{}),
	}
}

// addHealthSnapshot adds a health snapshot to history
func (sm *SystemMonitor) addHealthSnapshot(snapshot *SystemHealthSnapshot) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.healthHistory = append(sm.healthHistory, snapshot)

	// Keep only the most recent snapshots
	if len(sm.healthHistory) > sm.maxHistorySize {
		sm.healthHistory = sm.healthHistory[1:]
	}
}

// processAlerts processes health snapshot and generates alerts
func (sm *SystemMonitor) processAlerts(snapshot *SystemHealthSnapshot) {
	alerts := make([]HealthAlert, 0)

	// Check for component health issues
	for componentName, health := range snapshot.Components {
		switch health.Status {
		case "unhealthy":
			alert := HealthAlert{
				ID:        fmt.Sprintf("health_%s_%d", componentName, snapshot.Timestamp.Unix()),
				Severity:  "error",
				Component: componentName,
				Message:   fmt.Sprintf("Component %s is unhealthy", componentName),
				Timestamp: snapshot.Timestamp,
				Metadata: map[string]interface{}{
					"error_count":   health.ErrorCount,
					"warning_count": health.WarningCount,
					"response_time": health.ResponseTime,
				},
			}
			alerts = append(alerts, alert)
		case "degraded":
			alert := HealthAlert{
				ID:        fmt.Sprintf("degraded_%s_%d", componentName, snapshot.Timestamp.Unix()),
				Severity:  "warning",
				Component: componentName,
				Message:   fmt.Sprintf("Component %s is degraded", componentName),
				Timestamp: snapshot.Timestamp,
				Metadata: map[string]interface{}{
					"warning_count": health.WarningCount,
					"response_time": health.ResponseTime,
				},
			}
			alerts = append(alerts, alert)
		}
	}

	// Check system metrics for issues
	if snapshot.SystemMetrics != nil {
		if snapshot.SystemMetrics.CPUUsage > 90.0 {
			alerts = append(alerts, HealthAlert{
				ID:        fmt.Sprintf("cpu_high_%d", snapshot.Timestamp.Unix()),
				Severity:  "warning",
				Component: "system",
				Message:   fmt.Sprintf("High CPU usage: %.1f%%", snapshot.SystemMetrics.CPUUsage),
				Timestamp: snapshot.Timestamp,
				Metadata:  map[string]interface{}{"cpu_usage": snapshot.SystemMetrics.CPUUsage},
			})
		}

		if snapshot.SystemMetrics.MemoryUsage > 1024*1024*1024 { // 1GB
			alerts = append(alerts, HealthAlert{
				ID:        fmt.Sprintf("memory_high_%d", snapshot.Timestamp.Unix()),
				Severity:  "warning",
				Component: "system",
				Message:   fmt.Sprintf("High memory usage: %d MB", snapshot.SystemMetrics.MemoryUsage/(1024*1024)),
				Timestamp: snapshot.Timestamp,
				Metadata:  map[string]interface{}{"memory_usage": snapshot.SystemMetrics.MemoryUsage},
			})
		}
	}

	snapshot.Alerts = alerts

	// Log alerts
	for _, alert := range alerts {
		if sm.logger != nil {
			sm.logger.Warn("Health alert generated", "alert_id", alert.ID, "severity", alert.Severity, "message", alert.Message)
		}
	}
}

// GetCurrentHealth returns the current system health status
func (sm *SystemMonitor) GetCurrentHealth(ctx context.Context) *SystemHealthSnapshot {
	return sm.collectHealthSnapshot(ctx)
}

// GetHealthHistory returns the health history
func (sm *SystemMonitor) GetHealthHistory() []*SystemHealthSnapshot {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// Return a copy to prevent external modification
	history := make([]*SystemHealthSnapshot, len(sm.healthHistory))
	copy(history, sm.healthHistory)
	return history
}

// GetDiagnostics returns comprehensive diagnostic information
func (sm *SystemMonitor) GetDiagnostics() *DiagnosticsReport {
	return sm.diagnostics.GenerateReport()
}

// RecordError records an error for diagnostics
func (sm *SystemMonitor) RecordError(component, operation, errorType, message, severity string, context map[string]interface{}) {
	sm.diagnostics.RecordError(component, operation, errorType, message, severity, context)
}

// RecordPerformance records a performance measurement
func (sm *SystemMonitor) RecordPerformance(operation string, duration time.Duration, memoryUsed int64, success bool, metadata map[string]interface{}) {
	sm.diagnostics.RecordPerformance(operation, duration, memoryUsed, success, metadata)
}

// SetMonitorInterval sets the monitoring interval
func (sm *SystemMonitor) SetMonitorInterval(interval time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.monitorInterval = interval
}

// safeUint64ToInt64 safely converts uint64 to int64 with bounds checking
func safeUint64ToInt64(value uint64) int64 {
	if value > 0x7FFFFFFFFFFFFFFF { // Max int64 value
		return 0x7FFFFFFFFFFFFFFF
	}
	return int64(value)
}

// Helper functions for system metrics collection
func getCPUUsage() float64 {
	// Simplified CPU usage calculation
	// In a real implementation, this would use system-specific APIs
	return 0.0
}

func getDiskUsage() int64 {
	// Simplified disk usage calculation
	// In a real implementation, this would check actual disk usage
	return 0
}

func getNetworkLatency() time.Duration {
	// Simplified network latency calculation
	// In a real implementation, this would ping a known endpoint
	return 0
}

func (sm *SystemMonitor) getCacheHitRate() float64 {
	if sm.cache == nil {
		return 0.0
	}

	stats, err := sm.cache.GetStats()
	if err != nil {
		return 0.0
	}

	return stats.HitRate
}

// DiagnosticsReport contains comprehensive diagnostic information
type DiagnosticsReport struct {
	GeneratedAt        time.Time           `json:"generated_at"`
	SystemInfo         *SystemInfo         `json:"system_info"`
	RuntimeMetrics     *RuntimeMetrics     `json:"runtime_metrics"`
	ErrorSummary       *ErrorSummary       `json:"error_summary"`
	PerformanceSummary *MetricsReport      `json:"performance_summary"`
	RecentErrors       []ErrorRecord       `json:"recent_errors"`
	RecentPerformance  []PerformanceRecord `json:"recent_performance"`
	Recommendations    []string            `json:"recommendations"`
}

// ErrorSummary provides summary statistics for errors
type ErrorSummary struct {
	TotalErrors       int            `json:"total_errors"`
	ErrorsByType      map[string]int `json:"errors_by_type"`
	ErrorsByComponent map[string]int `json:"errors_by_component"`
	ErrorsBySeverity  map[string]int `json:"errors_by_severity"`
	RecentErrorRate   float64        `json:"recent_error_rate"`
}

// GenerateReport generates a comprehensive diagnostics report
func (dc *DiagnosticsCollector) GenerateReport() *DiagnosticsReport {
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	report := &DiagnosticsReport{
		GeneratedAt:       time.Now(),
		SystemInfo:        dc.systemInfo,
		RuntimeMetrics:    dc.collectRuntimeMetrics(),
		ErrorSummary:      dc.generateErrorSummary(),
		RecentErrors:      dc.getRecentErrors(50),
		RecentPerformance: dc.getRecentPerformance(50),
		Recommendations:   dc.generateRecommendations(),
	}

	return report
}

// collectRuntimeMetrics collects current runtime metrics
func (dc *DiagnosticsCollector) collectRuntimeMetrics() *RuntimeMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &RuntimeMetrics{
		Timestamp:     time.Now(),
		MemoryAlloc:   memStats.Alloc,
		MemoryTotal:   memStats.TotalAlloc,
		MemorySys:     memStats.Sys,
		NumGoroutines: runtime.NumGoroutine(),
		NumGC:         memStats.NumGC,
		GCPauseTotal:  memStats.PauseTotalNs,
		HeapObjects:   memStats.HeapObjects,
		StackInUse:    memStats.StackInuse,
	}
}

// RecordError records an error occurrence
func (dc *DiagnosticsCollector) RecordError(component, operation, errorType, message, severity string, context map[string]interface{}) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	record := ErrorRecord{
		Timestamp: time.Now(),
		Component: component,
		Operation: operation,
		ErrorType: errorType,
		Message:   message,
		Severity:  severity,
		Context:   context,
		Resolved:  false,
	}

	dc.errorHistory = append(dc.errorHistory, record)

	// Keep only recent records
	if len(dc.errorHistory) > dc.maxRecords {
		dc.errorHistory = dc.errorHistory[1:]
	}
}

// RecordPerformance records a performance measurement
func (dc *DiagnosticsCollector) RecordPerformance(operation string, duration time.Duration, memoryUsed int64, success bool, metadata map[string]interface{}) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	record := PerformanceRecord{
		Timestamp:  time.Now(),
		Operation:  operation,
		Duration:   duration,
		MemoryUsed: memoryUsed,
		Success:    success,
		Metadata:   metadata,
	}

	dc.performanceLog = append(dc.performanceLog, record)

	// Keep only recent records
	if len(dc.performanceLog) > dc.maxRecords {
		dc.performanceLog = dc.performanceLog[1:]
	}
}

// generateErrorSummary generates error summary statistics
func (dc *DiagnosticsCollector) generateErrorSummary() *ErrorSummary {
	summary := &ErrorSummary{
		ErrorsByType:      make(map[string]int),
		ErrorsByComponent: make(map[string]int),
		ErrorsBySeverity:  make(map[string]int),
	}

	recentCutoff := time.Now().Add(-1 * time.Hour)
	var recentErrors int

	for _, record := range dc.errorHistory {
		summary.TotalErrors++
		summary.ErrorsByType[record.ErrorType]++
		summary.ErrorsByComponent[record.Component]++
		summary.ErrorsBySeverity[record.Severity]++

		if record.Timestamp.After(recentCutoff) {
			recentErrors++
		}
	}

	if summary.TotalErrors > 0 {
		summary.RecentErrorRate = float64(recentErrors) / float64(summary.TotalErrors)
	}

	return summary
}

// getRecentErrors returns the most recent error records
func (dc *DiagnosticsCollector) getRecentErrors(limit int) []ErrorRecord {
	if len(dc.errorHistory) == 0 {
		return []ErrorRecord{}
	}

	start := len(dc.errorHistory) - limit
	if start < 0 {
		start = 0
	}

	result := make([]ErrorRecord, len(dc.errorHistory[start:]))
	copy(result, dc.errorHistory[start:])
	return result
}

// getRecentPerformance returns the most recent performance records
func (dc *DiagnosticsCollector) getRecentPerformance(limit int) []PerformanceRecord {
	if len(dc.performanceLog) == 0 {
		return []PerformanceRecord{}
	}

	start := len(dc.performanceLog) - limit
	if start < 0 {
		start = 0
	}

	result := make([]PerformanceRecord, len(dc.performanceLog[start:]))
	copy(result, dc.performanceLog[start:])
	return result
}

// generateRecommendations generates diagnostic recommendations
func (dc *DiagnosticsCollector) generateRecommendations() []string {
	recommendations := make([]string, 0)

	// Analyze error patterns
	errorSummary := dc.generateErrorSummary()
	if errorSummary.RecentErrorRate > 0.1 {
		recommendations = append(recommendations, "High recent error rate detected. Consider investigating error patterns.")
	}

	// Analyze memory usage
	runtimeMetrics := dc.collectRuntimeMetrics()
	if runtimeMetrics.MemoryAlloc > 500*1024*1024 { // 500MB
		recommendations = append(recommendations, "High memory usage detected. Consider optimizing memory allocation.")
	}

	if runtimeMetrics.NumGoroutines > 1000 {
		recommendations = append(recommendations, "High number of goroutines detected. Check for goroutine leaks.")
	}

	return recommendations
}
