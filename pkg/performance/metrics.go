// Package performance provides metrics collection and analysis capabilities
package performance

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// MetricsCollector collects and analyzes performance metrics
type MetricsCollector struct {
	commands      map[string][]*CommandMetrics
	systemMetrics *SystemMetrics
	mutex         sync.RWMutex
	startTime     time.Time
	totalCommands int64
	totalDuration time.Duration
	totalMemory   int64
}

// SystemMetrics tracks system-level performance metrics
type SystemMetrics struct {
	CPUUsage       float64                `json:"cpu_usage"`
	MemoryUsage    int64                  `json:"memory_usage"`
	DiskUsage      int64                  `json:"disk_usage"`
	NetworkLatency time.Duration          `json:"network_latency"`
	CacheHitRate   float64                `json:"cache_hit_rate"`
	LastUpdated    time.Time              `json:"last_updated"`
	CustomMetrics  map[string]interface{} `json:"custom_metrics"`
}

// MetricsReport provides comprehensive performance analysis
type MetricsReport struct {
	GeneratedAt       time.Time                  `json:"generated_at"`
	CollectionPeriod  time.Duration              `json:"collection_period"`
	TotalCommands     int64                      `json:"total_commands"`
	AverageDuration   time.Duration              `json:"average_duration"`
	TotalMemoryUsage  int64                      `json:"total_memory_usage"`
	CommandSummary    map[string]*CommandSummary `json:"command_summary"`
	SystemMetrics     *SystemMetrics             `json:"system_metrics"`
	TopSlowCommands   []*CommandMetrics          `json:"top_slow_commands"`
	TopMemoryCommands []*CommandMetrics          `json:"top_memory_commands"`
	Recommendations   []string                   `json:"recommendations"`
}

// CommandSummary provides summary statistics for a specific command
type CommandSummary struct {
	CommandName        string        `json:"command_name"`
	ExecutionCount     int           `json:"execution_count"`
	TotalDuration      time.Duration `json:"total_duration"`
	AverageDuration    time.Duration `json:"average_duration"`
	MinDuration        time.Duration `json:"min_duration"`
	MaxDuration        time.Duration `json:"max_duration"`
	TotalMemoryUsage   int64         `json:"total_memory_usage"`
	AverageMemoryUsage int64         `json:"average_memory_usage"`
	CacheHitRate       float64       `json:"cache_hit_rate"`
	OptimizationsUsed  []string      `json:"optimizations_used"`
	LastExecution      time.Time     `json:"last_execution"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		commands: make(map[string][]*CommandMetrics),
		systemMetrics: &SystemMetrics{
			CustomMetrics: make(map[string]interface{}),
		},
		startTime: time.Now(),
	}
}

// RecordCommand records metrics for a command execution
func (mc *MetricsCollector) RecordCommand(metrics *CommandMetrics) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Add to command history
	if mc.commands[metrics.CommandName] == nil {
		mc.commands[metrics.CommandName] = make([]*CommandMetrics, 0)
	}
	mc.commands[metrics.CommandName] = append(mc.commands[metrics.CommandName], metrics)

	// Update totals
	mc.totalCommands++
	mc.totalDuration += metrics.Duration
	mc.totalMemory += metrics.MemoryUsage

	// Keep only last 100 entries per command to prevent memory bloat
	if len(mc.commands[metrics.CommandName]) > 100 {
		mc.commands[metrics.CommandName] = mc.commands[metrics.CommandName][1:]
	}
}

// RecordSystemMetrics records system-level metrics
func (mc *MetricsCollector) RecordSystemMetrics(metrics *SystemMetrics) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.systemMetrics = metrics
	mc.systemMetrics.LastUpdated = time.Now()
}

// GetCommandMetrics returns metrics for a specific command
func (mc *MetricsCollector) GetCommandMetrics(commandName string) []*CommandMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	metrics, exists := mc.commands[commandName]
	if !exists {
		return nil
	}

	// Return a copy to prevent external modification
	result := make([]*CommandMetrics, len(metrics))
	copy(result, metrics)
	return result
}

// GetSystemMetrics returns current system metrics
func (mc *MetricsCollector) GetSystemMetrics() *SystemMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// Return a copy
	result := *mc.systemMetrics
	result.CustomMetrics = make(map[string]interface{})
	for k, v := range mc.systemMetrics.CustomMetrics {
		result.CustomMetrics[k] = v
	}
	return &result
}

// GenerateReport generates a comprehensive performance report
func (mc *MetricsCollector) GenerateReport() *MetricsReport {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	report := &MetricsReport{
		GeneratedAt:       time.Now(),
		CollectionPeriod:  time.Since(mc.startTime),
		TotalCommands:     mc.totalCommands,
		CommandSummary:    make(map[string]*CommandSummary),
		SystemMetrics:     mc.GetSystemMetrics(),
		TopSlowCommands:   make([]*CommandMetrics, 0),
		TopMemoryCommands: make([]*CommandMetrics, 0),
		Recommendations:   make([]string, 0),
	}

	if mc.totalCommands > 0 {
		report.AverageDuration = time.Duration(mc.totalDuration.Nanoseconds() / mc.totalCommands)
		report.TotalMemoryUsage = mc.totalMemory
	}

	// Generate command summaries
	allMetrics := make([]*CommandMetrics, 0)
	for commandName, metrics := range mc.commands {
		summary := mc.generateCommandSummary(commandName, metrics)
		report.CommandSummary[commandName] = summary
		allMetrics = append(allMetrics, metrics...)
	}

	// Find top slow commands
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].Duration > allMetrics[j].Duration
	})
	if len(allMetrics) > 10 {
		report.TopSlowCommands = allMetrics[:10]
	} else {
		report.TopSlowCommands = allMetrics
	}

	// Find top memory consuming commands
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].MemoryUsage > allMetrics[j].MemoryUsage
	})
	if len(allMetrics) > 10 {
		report.TopMemoryCommands = allMetrics[:10]
	} else {
		report.TopMemoryCommands = allMetrics
	}

	// Generate recommendations
	report.Recommendations = mc.generateRecommendations(report)

	return report
}

// generateCommandSummary generates summary statistics for a command
func (mc *MetricsCollector) generateCommandSummary(commandName string, metrics []*CommandMetrics) *CommandSummary {
	if len(metrics) == 0 {
		return &CommandSummary{CommandName: commandName}
	}

	summary := &CommandSummary{
		CommandName:       commandName,
		ExecutionCount:    len(metrics),
		OptimizationsUsed: make([]string, 0),
	}

	var totalDuration time.Duration
	var totalMemory int64
	var totalCacheHits, totalCacheRequests int
	optimizationMap := make(map[string]bool)

	// Initialize min/max with first metric
	summary.MinDuration = metrics[0].Duration
	summary.MaxDuration = metrics[0].Duration
	summary.LastExecution = metrics[0].EndTime

	for _, metric := range metrics {
		totalDuration += metric.Duration
		totalMemory += metric.MemoryUsage
		totalCacheHits += metric.CacheHits
		totalCacheRequests += metric.CacheHits + metric.CacheMisses

		// Track min/max duration
		if metric.Duration < summary.MinDuration {
			summary.MinDuration = metric.Duration
		}
		if metric.Duration > summary.MaxDuration {
			summary.MaxDuration = metric.Duration
		}

		// Track latest execution
		if metric.EndTime.After(summary.LastExecution) {
			summary.LastExecution = metric.EndTime
		}

		// Collect unique optimizations
		for _, opt := range metric.Optimizations {
			optimizationMap[opt] = true
		}
	}

	summary.TotalDuration = totalDuration
	summary.AverageDuration = time.Duration(totalDuration.Nanoseconds() / int64(len(metrics)))
	summary.TotalMemoryUsage = totalMemory
	summary.AverageMemoryUsage = totalMemory / int64(len(metrics))

	if totalCacheRequests > 0 {
		summary.CacheHitRate = float64(totalCacheHits) / float64(totalCacheRequests)
	}

	// Convert optimization map to slice
	for opt := range optimizationMap {
		summary.OptimizationsUsed = append(summary.OptimizationsUsed, opt)
	}

	return summary
}

// generateRecommendations generates performance recommendations based on metrics
func (mc *MetricsCollector) generateRecommendations(report *MetricsReport) []string {
	recommendations := make([]string, 0)

	// Check for slow commands
	if len(report.TopSlowCommands) > 0 {
		slowestCommand := report.TopSlowCommands[0]
		if slowestCommand.Duration > 5*time.Second {
			recommendations = append(recommendations,
				fmt.Sprintf("Command '%s' is taking %v to execute. Consider enabling caching or lazy loading.",
					slowestCommand.CommandName, slowestCommand.Duration))
		}
	}

	// Check cache hit rates
	for commandName, summary := range report.CommandSummary {
		if summary.ExecutionCount > 5 && summary.CacheHitRate < 0.5 {
			recommendations = append(recommendations,
				fmt.Sprintf("Command '%s' has low cache hit rate (%.1f%%). Consider adjusting cache TTL or cache key strategy.",
					commandName, summary.CacheHitRate*100))
		}
	}

	// Check memory usage
	if len(report.TopMemoryCommands) > 0 {
		memoryHeavyCommand := report.TopMemoryCommands[0]
		if memoryHeavyCommand.MemoryUsage > 100*1024*1024 { // 100MB
			recommendations = append(recommendations,
				fmt.Sprintf("Command '%s' is using %d MB of memory. Consider implementing streaming or chunked processing.",
					memoryHeavyCommand.CommandName, memoryHeavyCommand.MemoryUsage/(1024*1024)))
		}
	}

	// Check system metrics
	if report.SystemMetrics != nil {
		if report.SystemMetrics.CPUUsage > 80.0 {
			recommendations = append(recommendations, "High CPU usage detected. Consider reducing concurrent operations or optimizing algorithms.")
		}
		if report.SystemMetrics.CacheHitRate < 0.7 {
			recommendations = append(recommendations, "System cache hit rate is below optimal. Consider increasing cache size or reviewing cache strategy.")
		}
	}

	// Check for commands without optimizations
	for commandName, summary := range report.CommandSummary {
		if summary.ExecutionCount > 3 && len(summary.OptimizationsUsed) == 0 {
			recommendations = append(recommendations,
				fmt.Sprintf("Command '%s' is not using any optimizations. Consider enabling caching or lazy loading.", commandName))
		}
	}

	return recommendations
}

// GetAverageDuration returns the average duration for a specific command
func (mc *MetricsCollector) GetAverageDuration(commandName string) time.Duration {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	metrics, exists := mc.commands[commandName]
	if !exists || len(metrics) == 0 {
		return 0
	}

	var total time.Duration
	for _, metric := range metrics {
		total += metric.Duration
	}

	return time.Duration(total.Nanoseconds() / int64(len(metrics)))
}

// GetCacheHitRate returns the cache hit rate for a specific command
func (mc *MetricsCollector) GetCacheHitRate(commandName string) float64 {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	metrics, exists := mc.commands[commandName]
	if !exists || len(metrics) == 0 {
		return 0.0
	}

	var totalHits, totalRequests int
	for _, metric := range metrics {
		totalHits += metric.CacheHits
		totalRequests += metric.CacheHits + metric.CacheMisses
	}

	if totalRequests == 0 {
		return 0.0
	}

	return float64(totalHits) / float64(totalRequests)
}

// Reset clears all collected metrics
func (mc *MetricsCollector) Reset() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.commands = make(map[string][]*CommandMetrics)
	mc.systemMetrics = &SystemMetrics{
		CustomMetrics: make(map[string]interface{}),
	}
	mc.startTime = time.Now()
	mc.totalCommands = 0
	mc.totalDuration = 0
	mc.totalMemory = 0
}

// GetCommandNames returns all command names that have metrics
func (mc *MetricsCollector) GetCommandNames() []string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	names := make([]string, 0, len(mc.commands))
	for name := range mc.commands {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// GetTotalExecutions returns the total number of command executions
func (mc *MetricsCollector) GetTotalExecutions() int64 {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.totalCommands
}

// GetCollectionPeriod returns how long metrics have been collected
func (mc *MetricsCollector) GetCollectionPeriod() time.Duration {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return time.Since(mc.startTime)
}
