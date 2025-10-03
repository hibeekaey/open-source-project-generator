// Package performance provides monitoring dashboard and reporting capabilities
package performance

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Dashboard provides monitoring dashboard functionality
type Dashboard struct {
	monitor     *SystemMonitor
	metrics     *MetricsCollector
	refreshRate time.Duration
}

// DashboardData contains all data needed for the monitoring dashboard
type DashboardData struct {
	Timestamp       time.Time             `json:"timestamp"`
	SystemHealth    *SystemHealthSnapshot `json:"system_health"`
	PerformanceData *PerformanceReport    `json:"performance_data"`
	Diagnostics     *DiagnosticsReport    `json:"diagnostics"`
	Alerts          []HealthAlert         `json:"alerts"`
	Trends          *TrendAnalysis        `json:"trends"`
}

// TrendAnalysis provides trend analysis for monitoring data
type TrendAnalysis struct {
	HealthTrend      string                 `json:"health_trend"`      // "improving", "stable", "degrading"
	PerformanceTrend string                 `json:"performance_trend"` // "improving", "stable", "degrading"
	ErrorTrend       string                 `json:"error_trend"`       // "improving", "stable", "degrading"
	MemoryTrend      string                 `json:"memory_trend"`      // "improving", "stable", "degrading"
	TrendPeriod      time.Duration          `json:"trend_period"`
	TrendData        map[string][]float64   `json:"trend_data"`
	Predictions      map[string]interface{} `json:"predictions"`
}

// NewDashboard creates a new monitoring dashboard
func NewDashboard(monitor *SystemMonitor, metrics *MetricsCollector) *Dashboard {
	return &Dashboard{
		monitor:     monitor,
		metrics:     metrics,
		refreshRate: 30 * time.Second,
	}
}

// GetDashboardData returns comprehensive dashboard data
func (d *Dashboard) GetDashboardData() *DashboardData {
	data := &DashboardData{
		Timestamp: time.Now(),
	}

	// Get current system health
	if d.monitor != nil {
		data.SystemHealth = d.monitor.GetCurrentHealth(context.TODO())
		data.Diagnostics = d.monitor.GetDiagnostics()

		// Collect alerts from recent health snapshots
		data.Alerts = d.collectRecentAlerts()

		// Generate trend analysis
		data.Trends = d.generateTrendAnalysis()
	}

	// Get performance data
	if d.metrics != nil {
		metricsReport := d.metrics.GenerateReport()
		// Convert MetricsReport to PerformanceReport for compatibility
		data.PerformanceData = &PerformanceReport{
			GeneratedAt:     metricsReport.GeneratedAt,
			TotalBenchmarks: int(metricsReport.TotalCommands),
			Summary: &PerformanceSummary{
				OverallStatus: "active",
			},
		}
	}

	return data
}

// collectRecentAlerts collects alerts from recent health snapshots
func (d *Dashboard) collectRecentAlerts() []HealthAlert {
	if d.monitor == nil {
		return []HealthAlert{}
	}

	history := d.monitor.GetHealthHistory()
	alerts := make([]HealthAlert, 0)

	// Collect alerts from the last hour
	cutoff := time.Now().Add(-1 * time.Hour)

	for _, snapshot := range history {
		if snapshot.Timestamp.After(cutoff) {
			alerts = append(alerts, snapshot.Alerts...)
		}
	}

	// Sort alerts by timestamp (newest first)
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].Timestamp.After(alerts[j].Timestamp)
	})

	// Return only the most recent 50 alerts
	if len(alerts) > 50 {
		alerts = alerts[:50]
	}

	return alerts
}

// generateTrendAnalysis generates trend analysis from historical data
func (d *Dashboard) generateTrendAnalysis() *TrendAnalysis {
	if d.monitor == nil {
		return &TrendAnalysis{
			HealthTrend:      "unknown",
			PerformanceTrend: "unknown",
			ErrorTrend:       "unknown",
			MemoryTrend:      "unknown",
			TrendPeriod:      1 * time.Hour,
			TrendData:        make(map[string][]float64),
			Predictions:      make(map[string]interface{}),
		}
	}

	history := d.monitor.GetHealthHistory()
	if len(history) < 3 {
		return &TrendAnalysis{
			HealthTrend:      "insufficient_data",
			PerformanceTrend: "insufficient_data",
			ErrorTrend:       "insufficient_data",
			MemoryTrend:      "insufficient_data",
			TrendPeriod:      1 * time.Hour,
			TrendData:        make(map[string][]float64),
			Predictions:      make(map[string]interface{}),
		}
	}

	analysis := &TrendAnalysis{
		TrendPeriod: time.Since(history[0].Timestamp),
		TrendData:   make(map[string][]float64),
		Predictions: make(map[string]interface{}),
	}

	// Analyze health trend
	healthScores := make([]float64, 0)
	memoryUsage := make([]float64, 0)
	errorCounts := make([]float64, 0)

	for _, snapshot := range history {
		// Calculate health score (healthy=1.0, degraded=0.5, unhealthy=0.0)
		var healthScore float64
		switch snapshot.OverallStatus {
		case "healthy":
			healthScore = 1.0
		case "degraded":
			healthScore = 0.5
		case "unhealthy":
			healthScore = 0.0
		}
		healthScores = append(healthScores, healthScore)

		// Collect memory usage
		if snapshot.SystemMetrics != nil {
			memoryUsage = append(memoryUsage, float64(snapshot.SystemMetrics.MemoryUsage))
		}

		// Count errors in alerts
		errorCount := 0.0
		for _, alert := range snapshot.Alerts {
			if alert.Severity == "error" || alert.Severity == "critical" {
				errorCount++
			}
		}
		errorCounts = append(errorCounts, errorCount)
	}

	// Store trend data
	analysis.TrendData["health_scores"] = healthScores
	analysis.TrendData["memory_usage"] = memoryUsage
	analysis.TrendData["error_counts"] = errorCounts

	// Analyze trends
	analysis.HealthTrend = d.analyzeTrend(healthScores)
	analysis.MemoryTrend = d.analyzeTrend(memoryUsage)
	analysis.ErrorTrend = d.analyzeInverseTrend(errorCounts) // Lower error count is better

	// Generate performance trend from metrics
	if d.metrics != nil {
		report := d.metrics.GenerateReport()
		if report.TotalCommands > 0 {
			// Simple performance trend based on average duration
			avgDurationMs := float64(report.AverageDuration.Milliseconds())
			if avgDurationMs < 1000 {
				analysis.PerformanceTrend = "good"
			} else if avgDurationMs < 5000 {
				analysis.PerformanceTrend = "acceptable"
			} else {
				analysis.PerformanceTrend = "poor"
			}
		} else {
			analysis.PerformanceTrend = "no_data"
		}
	}

	// Generate simple predictions
	if len(healthScores) >= 3 {
		analysis.Predictions["health_prediction"] = d.predictTrend(healthScores)
	}
	if len(memoryUsage) >= 3 {
		analysis.Predictions["memory_prediction"] = d.predictTrend(memoryUsage)
	}

	return analysis
}

// analyzeTrend analyzes a trend from a series of values (higher is better)
func (d *Dashboard) analyzeTrend(values []float64) string {
	if len(values) < 3 {
		return "insufficient_data"
	}

	// Calculate simple linear trend
	n := len(values)
	recent := values[n-3:] // Last 3 values

	if len(recent) < 3 {
		return "insufficient_data"
	}

	// Simple trend analysis: compare first and last of recent values
	if recent[2] > recent[0]*1.1 {
		return "improving"
	} else if recent[2] < recent[0]*0.9 {
		return "degrading"
	} else {
		return "stable"
	}
}

// analyzeInverseTrend analyzes a trend where lower values are better
func (d *Dashboard) analyzeInverseTrend(values []float64) string {
	if len(values) < 3 {
		return "insufficient_data"
	}

	n := len(values)
	recent := values[n-3:]

	if len(recent) < 3 {
		return "insufficient_data"
	}

	// For inverse trends, decreasing values are improving
	if recent[2] < recent[0]*0.9 {
		return "improving"
	} else if recent[2] > recent[0]*1.1 {
		return "degrading"
	} else {
		return "stable"
	}
}

// predictTrend provides simple trend prediction
func (d *Dashboard) predictTrend(values []float64) map[string]interface{} {
	if len(values) < 3 {
		return map[string]interface{}{
			"prediction": "insufficient_data",
		}
	}

	// Simple linear regression for prediction
	n := len(values)
	recent := values[n-3:]

	// Calculate simple slope
	slope := (recent[2] - recent[0]) / 2.0

	// Predict next value
	nextValue := recent[2] + slope

	return map[string]interface{}{
		"next_value":     nextValue,
		"trend_slope":    slope,
		"confidence":     "low", // Simple prediction has low confidence
		"prediction_for": "next_measurement",
	}
}

// GenerateTextReport generates a text-based dashboard report
func (d *Dashboard) GenerateTextReport() string {
	data := d.GetDashboardData()

	var report strings.Builder

	report.WriteString("=== System Monitoring Dashboard ===\n")
	report.WriteString(fmt.Sprintf("Generated: %s\n\n", data.Timestamp.Format(time.RFC3339)))

	// System Health Section
	if data.SystemHealth != nil {
		report.WriteString("🏥 SYSTEM HEALTH\n")
		report.WriteString(fmt.Sprintf("Overall Status: %s\n", strings.ToUpper(data.SystemHealth.OverallStatus)))

		if len(data.SystemHealth.Components) > 0 {
			report.WriteString("\nComponent Status:\n")
			for name, health := range data.SystemHealth.Components {
				status := health.Status
				if health.ErrorCount > 0 {
					status += fmt.Sprintf(" (%d errors)", health.ErrorCount)
				}
				if health.WarningCount > 0 {
					status += fmt.Sprintf(" (%d warnings)", health.WarningCount)
				}
				report.WriteString(fmt.Sprintf("  %s: %s\n", name, status))
			}
		}
		report.WriteString("\n")
	}

	// Performance Section
	if data.PerformanceData != nil {
		report.WriteString("⚡ PERFORMANCE METRICS\n")
		report.WriteString(fmt.Sprintf("Total Benchmarks: %d\n", data.PerformanceData.TotalBenchmarks))
		if data.PerformanceData.Summary != nil {
			report.WriteString(fmt.Sprintf("Overall Status: %s\n", data.PerformanceData.Summary.OverallStatus))
		}
		report.WriteString("\n")
	}

	// Alerts Section
	if len(data.Alerts) > 0 {
		report.WriteString("🚨 RECENT ALERTS\n")
		alertCount := len(data.Alerts)
		if alertCount > 10 {
			alertCount = 10 // Show only recent 10
		}

		for i := 0; i < alertCount; i++ {
			alert := data.Alerts[i]
			report.WriteString(fmt.Sprintf("  [%s] %s: %s (%s)\n",
				strings.ToUpper(alert.Severity),
				alert.Component,
				alert.Message,
				alert.Timestamp.Format("15:04:05")))
		}
		report.WriteString("\n")
	}

	// Trends Section
	if data.Trends != nil {
		report.WriteString("📈 TRENDS\n")
		report.WriteString(fmt.Sprintf("Health Trend: %s\n", data.Trends.HealthTrend))
		report.WriteString(fmt.Sprintf("Performance Trend: %s\n", data.Trends.PerformanceTrend))
		report.WriteString(fmt.Sprintf("Memory Trend: %s\n", data.Trends.MemoryTrend))
		report.WriteString(fmt.Sprintf("Error Trend: %s\n", data.Trends.ErrorTrend))
		report.WriteString("\n")
	}

	// Diagnostics Section
	if data.Diagnostics != nil && data.Diagnostics.ErrorSummary != nil {
		report.WriteString("🔍 DIAGNOSTICS SUMMARY\n")
		report.WriteString(fmt.Sprintf("Total Errors: %d\n", data.Diagnostics.ErrorSummary.TotalErrors))
		report.WriteString(fmt.Sprintf("Recent Error Rate: %.2f%%\n", data.Diagnostics.ErrorSummary.RecentErrorRate*100))

		if len(data.Diagnostics.Recommendations) > 0 {
			report.WriteString("\nRecommendations:\n")
			for _, rec := range data.Diagnostics.Recommendations {
				report.WriteString(fmt.Sprintf("  • %s\n", rec))
			}
		}
		report.WriteString("\n")
	}

	report.WriteString("=== End of Report ===\n")

	return report.String()
}

// GenerateJSONReport generates a JSON-based dashboard report
func (d *Dashboard) GenerateJSONReport() (string, error) {
	data := d.GetDashboardData()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal dashboard data: %w", err)
	}

	return string(jsonData), nil
}

// GetAlertSummary returns a summary of current alerts
func (d *Dashboard) GetAlertSummary() map[string]int {
	alerts := d.collectRecentAlerts()

	summary := map[string]int{
		"critical": 0,
		"error":    0,
		"warning":  0,
		"info":     0,
	}

	for _, alert := range alerts {
		if !alert.Resolved {
			summary[alert.Severity]++
		}
	}

	return summary
}

// GetHealthSummary returns a summary of component health
func (d *Dashboard) GetHealthSummary() map[string]int {
	if d.monitor == nil {
		return map[string]int{}
	}

	health := d.monitor.GetCurrentHealth(context.TODO())
	if health == nil {
		return map[string]int{}
	}

	summary := map[string]int{
		"healthy":   0,
		"degraded":  0,
		"unhealthy": 0,
	}

	for _, component := range health.Components {
		summary[component.Status]++
	}

	return summary
}

// SetRefreshRate sets the dashboard refresh rate
func (d *Dashboard) SetRefreshRate(rate time.Duration) {
	d.refreshRate = rate
}

// GetRefreshRate returns the current refresh rate
func (d *Dashboard) GetRefreshRate() time.Duration {
	return d.refreshRate
}
