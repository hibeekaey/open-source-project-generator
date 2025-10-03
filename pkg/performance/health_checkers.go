// Package performance provides health checker implementations for system monitoring
package performance

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CacheHealthChecker checks the health of the cache system
type CacheHealthChecker struct {
	cache interfaces.CacheManager
}

// CheckHealth performs a health check on the cache system
func (chc *CacheHealthChecker) CheckHealth(ctx context.Context) *ComponentHealth {
	startTime := time.Now()
	health := &ComponentHealth{
		ComponentName:   "cache",
		LastCheck:       startTime,
		Details:         make(map[string]interface{}),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	if chc.cache == nil {
		health.Status = "unhealthy"
		health.Issues = append(health.Issues, "Cache manager not initialized")
		health.ErrorCount = 1
		health.ResponseTime = time.Since(startTime)
		return health
	}

	// Test basic cache operations
	testKey := "health_check_test"
	testValue := "test_value"

	// Test cache write
	if err := chc.cache.Set(testKey, testValue, 1*time.Minute); err != nil {
		health.Status = "unhealthy"
		health.Issues = append(health.Issues, fmt.Sprintf("Cache write failed: %v", err))
		health.ErrorCount++
	}

	// Test cache read
	if _, err := chc.cache.Get(testKey); err != nil {
		health.Status = "degraded"
		health.Issues = append(health.Issues, fmt.Sprintf("Cache read failed: %v", err))
		health.WarningCount++
	}

	// Clean up test data
	_ = chc.cache.Delete(testKey)

	// Get cache statistics
	if stats, err := chc.cache.GetStats(); err == nil {
		health.Details["total_entries"] = stats.TotalEntries
		health.Details["total_size"] = stats.TotalSize
		health.Details["hit_rate"] = stats.HitRate
		health.Details["last_accessed"] = stats.LastAccessed
		health.Details["offline_mode"] = chc.cache.IsOfflineMode()

		// Check cache performance
		if stats.HitRate < 0.5 {
			health.WarningCount++
			health.Issues = append(health.Issues, "Low cache hit rate")
			health.Recommendations = append(health.Recommendations, "Consider adjusting cache TTL or cache strategy")
		}

		if stats.TotalEntries > 10000 {
			health.WarningCount++
			health.Issues = append(health.Issues, "High number of cache entries")
			health.Recommendations = append(health.Recommendations, "Consider running cache cleanup")
		}
	} else {
		health.WarningCount++
		health.Issues = append(health.Issues, fmt.Sprintf("Failed to get cache stats: %v", err))
	}

	// Determine overall status
	if health.ErrorCount > 0 {
		health.Status = "unhealthy"
	} else if health.WarningCount > 0 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	health.ResponseTime = time.Since(startTime)
	return health
}

// GetComponentName returns the component name
func (chc *CacheHealthChecker) GetComponentName() string {
	return "cache"
}

// MemoryHealthChecker checks system memory health
type MemoryHealthChecker struct{}

// CheckHealth performs a memory health check
func (mhc *MemoryHealthChecker) CheckHealth(ctx context.Context) *ComponentHealth {
	startTime := time.Now()
	health := &ComponentHealth{
		ComponentName:   "memory",
		LastCheck:       startTime,
		Details:         make(map[string]interface{}),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Collect memory statistics
	health.Details["alloc_bytes"] = memStats.Alloc
	health.Details["total_alloc_bytes"] = memStats.TotalAlloc
	health.Details["sys_bytes"] = memStats.Sys
	health.Details["num_gc"] = memStats.NumGC
	health.Details["heap_objects"] = memStats.HeapObjects
	health.Details["num_goroutines"] = runtime.NumGoroutine()

	// Check memory thresholds
	allocMB := memStats.Alloc / (1024 * 1024)
	sysMB := memStats.Sys / (1024 * 1024)

	if allocMB > 500 { // 500MB
		health.WarningCount++
		health.Issues = append(health.Issues, fmt.Sprintf("High memory allocation: %d MB", allocMB))
		health.Recommendations = append(health.Recommendations, "Consider optimizing memory usage or running garbage collection")
	}

	if sysMB > 1000 { // 1GB
		health.ErrorCount++
		health.Issues = append(health.Issues, fmt.Sprintf("Very high system memory usage: %d MB", sysMB))
		health.Recommendations = append(health.Recommendations, "Investigate memory leaks or reduce memory footprint")
	}

	// Check goroutine count
	numGoroutines := runtime.NumGoroutine()
	if numGoroutines > 1000 {
		health.WarningCount++
		health.Issues = append(health.Issues, fmt.Sprintf("High number of goroutines: %d", numGoroutines))
		health.Recommendations = append(health.Recommendations, "Check for goroutine leaks")
	}

	// Check GC frequency
	if memStats.NumGC > 0 {
		avgGCPause := memStats.PauseTotalNs / uint64(memStats.NumGC)
		if avgGCPause > 10*1000*1000 { // 10ms
			health.WarningCount++
			health.Issues = append(health.Issues, fmt.Sprintf("High GC pause time: %d ns", avgGCPause))
			health.Recommendations = append(health.Recommendations, "Consider optimizing memory allocation patterns")
		}
	}

	// Determine overall status
	if health.ErrorCount > 0 {
		health.Status = "unhealthy"
	} else if health.WarningCount > 0 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	health.ResponseTime = time.Since(startTime)
	return health
}

// GetComponentName returns the component name
func (mhc *MemoryHealthChecker) GetComponentName() string {
	return "memory"
}

// PerformanceHealthChecker checks performance metrics health
type PerformanceHealthChecker struct {
	metrics *MetricsCollector
}

// CheckHealth performs a performance health check
func (phc *PerformanceHealthChecker) CheckHealth(ctx context.Context) *ComponentHealth {
	startTime := time.Now()
	health := &ComponentHealth{
		ComponentName:   "performance",
		LastCheck:       startTime,
		Details:         make(map[string]interface{}),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	if phc.metrics == nil {
		health.Status = "degraded"
		health.WarningCount = 1
		health.Issues = append(health.Issues, "Performance metrics collector not available")
		health.ResponseTime = time.Since(startTime)
		return health
	}

	// Get performance report
	report := phc.metrics.GenerateReport()

	health.Details["total_commands"] = report.TotalCommands
	health.Details["average_duration"] = report.AverageDuration
	health.Details["collection_period"] = report.CollectionPeriod

	// Check for performance issues
	if report.AverageDuration > 5*time.Second {
		health.WarningCount++
		health.Issues = append(health.Issues, fmt.Sprintf("High average command duration: %v", report.AverageDuration))
		health.Recommendations = append(health.Recommendations, "Consider enabling performance optimizations")
	}

	// Check individual command performance
	slowCommands := 0
	for commandName, summary := range report.CommandSummary {
		if summary.AverageDuration > 10*time.Second {
			slowCommands++
			health.Issues = append(health.Issues, fmt.Sprintf("Slow command detected: %s (%v)", commandName, summary.AverageDuration))
		}

		if summary.CacheHitRate < 0.3 && summary.ExecutionCount > 5 {
			health.WarningCount++
			health.Issues = append(health.Issues, fmt.Sprintf("Low cache hit rate for %s: %.1f%%", commandName, summary.CacheHitRate*100))
			health.Recommendations = append(health.Recommendations, fmt.Sprintf("Optimize caching strategy for %s command", commandName))
		}
	}

	if slowCommands > 0 {
		health.WarningCount += slowCommands
		health.Recommendations = append(health.Recommendations, "Enable performance profiling for slow commands")
	}

	// Check system metrics if available
	if report.SystemMetrics != nil {
		health.Details["cpu_usage"] = report.SystemMetrics.CPUUsage
		health.Details["memory_usage"] = report.SystemMetrics.MemoryUsage
		health.Details["cache_hit_rate"] = report.SystemMetrics.CacheHitRate

		if report.SystemMetrics.CPUUsage > 80.0 {
			health.WarningCount++
			health.Issues = append(health.Issues, fmt.Sprintf("High CPU usage: %.1f%%", report.SystemMetrics.CPUUsage))
			health.Recommendations = append(health.Recommendations, "Consider reducing concurrent operations")
		}
	}

	// Determine overall status
	if health.ErrorCount > 0 {
		health.Status = "unhealthy"
	} else if health.WarningCount > 0 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	health.ResponseTime = time.Since(startTime)
	return health
}

// GetComponentName returns the component name
func (phc *PerformanceHealthChecker) GetComponentName() string {
	return "performance"
}

// TemplateHealthChecker checks template system health
type TemplateHealthChecker struct {
	templateManager interfaces.TemplateManager
}

// CheckHealth performs a template system health check
func (thc *TemplateHealthChecker) CheckHealth(ctx context.Context) *ComponentHealth {
	startTime := time.Now()
	health := &ComponentHealth{
		ComponentName:   "templates",
		LastCheck:       startTime,
		Details:         make(map[string]interface{}),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	if thc.templateManager == nil {
		health.Status = "unhealthy"
		health.ErrorCount = 1
		health.Issues = append(health.Issues, "Template manager not initialized")
		health.ResponseTime = time.Since(startTime)
		return health
	}

	// Test template listing
	templates, err := thc.templateManager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		health.Status = "degraded"
		health.WarningCount++
		health.Issues = append(health.Issues, fmt.Sprintf("Failed to list templates: %v", err))
	} else {
		health.Details["template_count"] = len(templates)

		if len(templates) == 0 {
			health.WarningCount++
			health.Issues = append(health.Issues, "No templates available")
			health.Recommendations = append(health.Recommendations, "Ensure template files are properly installed")
		}
	}

	// Test template validation for a few templates
	if len(templates) > 0 {
		validTemplates := 0
		for i, template := range templates {
			if i >= 3 { // Only check first 3 templates to avoid long health checks
				break
			}

			if _, err := thc.templateManager.ValidateTemplate(template.Name); err != nil {
				health.WarningCount++
				health.Issues = append(health.Issues, fmt.Sprintf("Template validation failed for %s: %v", template.Name, err))
			} else {
				validTemplates++
			}
		}

		health.Details["valid_templates_checked"] = validTemplates
	}

	// Determine overall status
	if health.ErrorCount > 0 {
		health.Status = "unhealthy"
	} else if health.WarningCount > 0 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	health.ResponseTime = time.Since(startTime)
	return health
}

// GetComponentName returns the component name
func (thc *TemplateHealthChecker) GetComponentName() string {
	return "templates"
}

// ConfigHealthChecker checks configuration system health
type ConfigHealthChecker struct {
	configManager interfaces.ConfigManager
}

// CheckHealth performs a configuration system health check
func (chc *ConfigHealthChecker) CheckHealth(ctx context.Context) *ComponentHealth {
	startTime := time.Now()
	health := &ComponentHealth{
		ComponentName:   "config",
		LastCheck:       startTime,
		Details:         make(map[string]interface{}),
		Issues:          make([]string, 0),
		Recommendations: make([]string, 0),
	}

	if chc.configManager == nil {
		health.Status = "unhealthy"
		health.ErrorCount = 1
		health.Issues = append(health.Issues, "Configuration manager not initialized")
		health.ResponseTime = time.Since(startTime)
		return health
	}

	// Test basic configuration operations
	health.Details["config_manager_available"] = true

	// Test configuration loading (basic check)
	if _, err := chc.configManager.LoadConfig(""); err != nil {
		// This is expected if no config file exists, so it's just informational
		health.Details["default_config_loadable"] = false
	} else {
		health.Details["default_config_loadable"] = true
	}

	// Determine overall status
	if health.ErrorCount > 0 {
		health.Status = "unhealthy"
	} else if health.WarningCount > 0 {
		health.Status = "degraded"
	} else {
		health.Status = "healthy"
	}

	health.ResponseTime = time.Since(startTime)
	return health
}

// GetComponentName returns the component name
func (chc *ConfigHealthChecker) GetComponentName() string {
	return "config"
}
