// Package metrics provides cache metrics collection and reporting.
package metrics

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Reporter handles cache metrics reporting and statistics generation.
type Reporter struct {
	collector *Collector
}

// NewReporter creates a new metrics reporter.
func NewReporter(collector *Collector) *Reporter {
	return &Reporter{
		collector: collector,
	}
}

// GenerateStats generates comprehensive cache statistics.
func (r *Reporter) GenerateStats(entries map[string]*interfaces.CacheEntry, cacheLocation string, offlineMode bool) *interfaces.CacheStats {
	metrics := r.collector.GetMetrics()

	// Count expired entries
	now := time.Now()
	expiredCount := 0
	for _, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredCount++
		}
	}

	// Calculate hit rate
	hitRate := r.collector.GetHitRate()

	return &interfaces.CacheStats{
		TotalEntries: len(entries),
		TotalSize:    metrics.CurrentSize,
		HitCount:     metrics.Hits,
		MissCount:    metrics.Misses,
		HitRate:      hitRate,
		LastAccessed: time.Now(), // Use current time as placeholder
		LastModified: time.Now(), // Use current time as placeholder
		CreatedAt:    time.Now(), // Use current time as placeholder
	}
}

// GenerateReport generates a detailed cache report.
func (r *Reporter) GenerateReport(entries map[string]*interfaces.CacheEntry, cacheLocation string, offlineMode bool) *CacheReport {
	metrics := r.collector.GetMetrics()
	stats := r.GenerateStats(entries, cacheLocation, offlineMode)

	return &CacheReport{
		Timestamp:       time.Now(),
		Uptime:          r.collector.GetUptime(),
		Stats:           stats,
		Metrics:         metrics,
		Performance:     r.generatePerformanceReport(metrics),
		Health:          r.generateHealthReport(entries, metrics),
		Recommendations: r.generateRecommendations(entries, metrics),
	}
}

// CacheReport represents a comprehensive cache report.
type CacheReport struct {
	Timestamp       time.Time                `json:"timestamp"`
	Uptime          time.Duration            `json:"uptime"`
	Stats           *interfaces.CacheStats   `json:"stats"`
	Metrics         *interfaces.CacheMetrics `json:"metrics"`
	Performance     *PerformanceReport       `json:"performance"`
	Health          *HealthReport            `json:"health"`
	Recommendations []string                 `json:"recommendations"`
}

// PerformanceReport represents cache performance metrics.
type PerformanceReport struct {
	HitRate             float64 `json:"hit_rate"`
	MissRate            float64 `json:"miss_rate"`
	TotalOperations     int64   `json:"total_operations"`
	OperationsPerSecond float64 `json:"operations_per_second"`
	AverageEntrySize    float64 `json:"average_entry_size"`
	CacheUtilization    float64 `json:"cache_utilization"`
}

// HealthReport represents cache health status.
type HealthReport struct {
	Status           string    `json:"status"`
	Issues           []string  `json:"issues"`
	Warnings         []string  `json:"warnings"`
	ExpiredEntries   int       `json:"expired_entries"`
	CorruptedEntries int       `json:"corrupted_entries"`
	LastCheck        time.Time `json:"last_check"`
}

// generatePerformanceReport generates performance metrics.
func (r *Reporter) generatePerformanceReport(metrics *interfaces.CacheMetrics) *PerformanceReport {
	uptime := r.collector.GetUptime()
	totalOps := r.collector.GetTotalOperations()

	var opsPerSecond float64
	if uptime.Seconds() > 0 {
		opsPerSecond = float64(totalOps) / uptime.Seconds()
	}

	var avgEntrySize float64
	if metrics.CurrentEntries > 0 {
		avgEntrySize = float64(metrics.CurrentSize) / float64(metrics.CurrentEntries)
	}

	var utilization float64
	if metrics.MaxSize > 0 {
		utilization = float64(metrics.CurrentSize) / float64(metrics.MaxSize)
	}

	return &PerformanceReport{
		HitRate:             r.collector.GetHitRate(),
		MissRate:            r.collector.GetMissRate(),
		TotalOperations:     totalOps,
		OperationsPerSecond: opsPerSecond,
		AverageEntrySize:    avgEntrySize,
		CacheUtilization:    utilization,
	}
}

// generateHealthReport generates health status report.
func (r *Reporter) generateHealthReport(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) *HealthReport {
	now := time.Now()
	expiredCount := 0
	corruptedCount := 0
	issues := make([]string, 0)
	warnings := make([]string, 0)

	// Check entries
	for key, entry := range entries {
		if entry == nil {
			corruptedCount++
			continue
		}

		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredCount++
		}

		// Check for corruption
		if entry.Size < 0 || entry.AccessCount < 0 || entry.Key != key {
			corruptedCount++
		}
	}

	// Determine status and issues
	status := "healthy"
	totalEntries := len(entries)

	if totalEntries > 0 {
		expiredRatio := float64(expiredCount) / float64(totalEntries)
		corruptedRatio := float64(corruptedCount) / float64(totalEntries)

		if expiredRatio > 0.5 {
			issues = append(issues, fmt.Sprintf("High expired entry ratio: %.1f%%", expiredRatio*100))
			status = "degraded"
		} else if expiredRatio > 0.2 {
			warnings = append(warnings, fmt.Sprintf("Moderate expired entry ratio: %.1f%%", expiredRatio*100))
		}

		if corruptedRatio > 0.1 {
			issues = append(issues, fmt.Sprintf("Corrupted entries detected: %.1f%%", corruptedRatio*100))
			status = "unhealthy"
		} else if corruptedRatio > 0.05 {
			warnings = append(warnings, fmt.Sprintf("Some corrupted entries detected: %.1f%%", corruptedRatio*100))
		}
	}

	// Check metrics health
	if metrics.CurrentSize > metrics.MaxSize && metrics.MaxSize > 0 {
		issues = append(issues, "Cache size exceeds maximum limit")
		if status == "healthy" {
			status = "degraded"
		}
	}

	if metrics.CurrentEntries > metrics.MaxEntries && metrics.MaxEntries > 0 {
		issues = append(issues, "Cache entry count exceeds maximum limit")
		if status == "healthy" {
			status = "degraded"
		}
	}

	// Check hit rate
	hitRate := r.collector.GetHitRate()
	if hitRate < 0.3 && metrics.Gets > 100 {
		warnings = append(warnings, fmt.Sprintf("Low cache hit rate: %.1f%%", hitRate*100))
	}

	return &HealthReport{
		Status:           status,
		Issues:           issues,
		Warnings:         warnings,
		ExpiredEntries:   expiredCount,
		CorruptedEntries: corruptedCount,
		LastCheck:        now,
	}
}

// generateRecommendations generates optimization recommendations.
func (r *Reporter) generateRecommendations(entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) []string {
	recommendations := make([]string, 0)

	// Check expired entries
	now := time.Now()
	expiredCount := 0
	for _, entry := range entries {
		if entry.ExpiresAt != nil && entry.ExpiresAt.Before(now) {
			expiredCount++
		}
	}

	totalEntries := len(entries)
	if totalEntries > 0 {
		expiredRatio := float64(expiredCount) / float64(totalEntries)
		if expiredRatio > 0.2 {
			recommendations = append(recommendations, "Consider running cache cleanup to remove expired entries")
		}
	}

	// Check hit rate
	hitRate := r.collector.GetHitRate()
	if hitRate < 0.5 && metrics.Gets > 100 {
		recommendations = append(recommendations, "Consider increasing cache size or reviewing cache strategy")
	}

	// Check size utilization
	if metrics.MaxSize > 0 {
		utilization := float64(metrics.CurrentSize) / float64(metrics.MaxSize)
		if utilization > 0.9 {
			recommendations = append(recommendations, "Consider increasing MaxSize - cache is nearly full")
		} else if utilization < 0.1 && metrics.CurrentSize > 0 {
			recommendations = append(recommendations, "Consider decreasing MaxSize - cache is underutilized")
		}
	}

	// Check entry count utilization
	if metrics.MaxEntries > 0 {
		entryUtilization := float64(metrics.CurrentEntries) / float64(metrics.MaxEntries)
		if entryUtilization > 0.9 {
			recommendations = append(recommendations, "Consider increasing MaxEntries - entry limit nearly reached")
		}
	}

	// Performance recommendations
	totalOps := r.collector.GetTotalOperations()
	uptime := r.collector.GetUptime()
	if uptime.Hours() > 1 && totalOps > 1000 {
		opsPerSecond := float64(totalOps) / uptime.Seconds()
		if opsPerSecond > 100 {
			recommendations = append(recommendations, "High operation rate detected - consider enabling compression")
		}
	}

	return recommendations
}
