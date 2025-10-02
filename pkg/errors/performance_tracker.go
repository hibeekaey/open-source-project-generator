// Package errors provides performance tracking for error handling operations
package errors

import (
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// PerformanceTracker tracks performance metrics for error handling operations
type PerformanceTracker struct {
	config           *EnhancedErrorConfig
	logger           interfaces.Logger
	stats            *PerformanceStatistics
	operationHistory []OperationRecord
	mutex            sync.RWMutex
}

// PerformanceInfo contains performance information for an operation
type PerformanceInfo struct {
	Operation       string                 `json:"operation"`
	Duration        time.Duration          `json:"duration"`
	Threshold       time.Duration          `json:"threshold,omitempty"`
	IsSlowOperation bool                   `json:"is_slow_operation"`
	Suggestions     []string               `json:"suggestions,omitempty"`
	Metrics         map[string]interface{} `json:"metrics,omitempty"`
}

// PerformanceStatistics tracks overall performance statistics
type PerformanceStatistics struct {
	TotalOperations      int            `json:"total_operations"`
	SlowOperations       int            `json:"slow_operations"`
	AverageDuration      time.Duration  `json:"average_duration"`
	MedianDuration       time.Duration  `json:"median_duration"`
	P95Duration          time.Duration  `json:"p95_duration"`
	P99Duration          time.Duration  `json:"p99_duration"`
	OperationsByType     map[string]int `json:"operations_by_type"`
	SlowOperationsByType map[string]int `json:"slow_operations_by_type"`
	ThresholdViolations  map[string]int `json:"threshold_violations"`
	LastUpdated          time.Time      `json:"last_updated"`
}

// OperationRecord records details about a single operation
type OperationRecord struct {
	Operation string                 `json:"operation"`
	StartTime time.Time              `json:"start_time"`
	Duration  time.Duration          `json:"duration"`
	Context   map[string]interface{} `json:"context,omitempty"`
	WasSlow   bool                   `json:"was_slow"`
	Threshold time.Duration          `json:"threshold"`
}

// NewPerformanceTracker creates a new performance tracker
func NewPerformanceTracker(config *EnhancedErrorConfig, logger interfaces.Logger) *PerformanceTracker {
	return &PerformanceTracker{
		config: config,
		logger: logger,
		stats: &PerformanceStatistics{
			OperationsByType:     make(map[string]int),
			SlowOperationsByType: make(map[string]int),
			ThresholdViolations:  make(map[string]int),
		},
		operationHistory: make([]OperationRecord, 0, 1000), // Keep last 1000 operations
	}
}

// TrackOperation tracks the performance of an operation
func (pt *PerformanceTracker) TrackOperation(operation string, duration time.Duration, context map[string]interface{}) *PerformanceInfo {
	if !pt.config.EnablePerformanceTracking {
		return nil
	}

	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	// Get threshold for this operation
	threshold := pt.getThreshold(operation)
	isSlowOperation := duration > threshold

	// Record the operation
	record := OperationRecord{
		Operation: operation,
		StartTime: time.Now().Add(-duration),
		Duration:  duration,
		Context:   context,
		WasSlow:   isSlowOperation,
		Threshold: threshold,
	}

	// Add to history (keep only last 1000)
	pt.operationHistory = append(pt.operationHistory, record)
	if len(pt.operationHistory) > 1000 {
		pt.operationHistory = pt.operationHistory[1:]
	}

	// Update statistics
	pt.updateStatistics(operation, duration, isSlowOperation, threshold)

	// Create performance info
	info := &PerformanceInfo{
		Operation:       operation,
		Duration:        duration,
		Threshold:       threshold,
		IsSlowOperation: isSlowOperation,
		Metrics:         pt.collectMetrics(operation, context),
	}

	// Generate suggestions for slow operations
	if isSlowOperation {
		info.Suggestions = pt.generatePerformanceSuggestions(operation, duration, threshold, context)
	}

	// Log performance information
	if pt.logger != nil {
		if isSlowOperation {
			pt.logger.Warn("Slow operation detected: %s took %v (threshold: %v)", operation, duration, threshold)
		} else if pt.config.VerboseMode {
			pt.logger.Debug("Operation completed: %s took %v", operation, duration)
		}
	}

	return info
}

// getThreshold returns the performance threshold for an operation
func (pt *PerformanceTracker) getThreshold(operation string) time.Duration {
	// Check for specific operation threshold
	if threshold, exists := pt.config.PerformanceThresholds[operation]; exists {
		return threshold
	}

	// Check for operation type threshold
	operationType := pt.getOperationType(operation)
	if threshold, exists := pt.config.PerformanceThresholds[operationType]; exists {
		return threshold
	}

	// Default threshold
	return 5 * time.Second
}

// getOperationType extracts the operation type from operation name
func (pt *PerformanceTracker) getOperationType(operation string) string {
	// Map operations to types
	operationTypes := map[string]string{
		"validation":    "validation",
		"validate":      "validation",
		"configuration": "configuration",
		"config":        "configuration",
		"template":      "template",
		"generation":    "generation",
		"generate":      "generation",
		"network":       "network",
		"download":      "network",
		"fetch":         "network",
		"filesystem":    "file_operation",
		"file":          "file_operation",
		"directory":     "file_operation",
		"cache":         "cache",
		"audit":         "audit",
	}

	for key, opType := range operationTypes {
		if operation == key || (len(operation) > len(key) && operation[:len(key)] == key) {
			return opType
		}
	}

	return "general"
}

// updateStatistics updates performance statistics
func (pt *PerformanceTracker) updateStatistics(operation string, duration time.Duration, isSlowOperation bool, threshold time.Duration) {
	pt.stats.TotalOperations++
	pt.stats.OperationsByType[operation]++
	pt.stats.LastUpdated = time.Now()

	if isSlowOperation {
		pt.stats.SlowOperations++
		pt.stats.SlowOperationsByType[operation]++
		pt.stats.ThresholdViolations[operation]++
	}

	// Update average duration
	if pt.stats.TotalOperations == 1 {
		pt.stats.AverageDuration = duration
	} else {
		pt.stats.AverageDuration = (pt.stats.AverageDuration*time.Duration(pt.stats.TotalOperations-1) + duration) / time.Duration(pt.stats.TotalOperations)
	}

	// Update percentiles (simplified calculation)
	pt.updatePercentiles()
}

// updatePercentiles updates duration percentiles
func (pt *PerformanceTracker) updatePercentiles() {
	if len(pt.operationHistory) == 0 {
		return
	}

	// Sort durations for percentile calculation
	durations := make([]time.Duration, len(pt.operationHistory))
	for i, record := range pt.operationHistory {
		durations[i] = record.Duration
	}

	// Simple bubble sort for small datasets
	for i := 0; i < len(durations)-1; i++ {
		for j := 0; j < len(durations)-i-1; j++ {
			if durations[j] > durations[j+1] {
				durations[j], durations[j+1] = durations[j+1], durations[j]
			}
		}
	}

	// Calculate percentiles
	n := len(durations)
	if n > 0 {
		pt.stats.MedianDuration = durations[n/2]
		pt.stats.P95Duration = durations[int(float64(n)*0.95)]
		pt.stats.P99Duration = durations[int(float64(n)*0.99)]
	}
}

// collectMetrics collects additional performance metrics
func (pt *PerformanceTracker) collectMetrics(operation string, context map[string]interface{}) map[string]interface{} {
	metrics := make(map[string]interface{})

	// Add context metrics
	for key, value := range context {
		metrics[key] = value
	}

	// Add operation-specific metrics
	switch operation {
	case "validation", "validate":
		metrics["validation_type"] = "project_structure"
	case "generation", "generate":
		if templateName, ok := context["template_name"]; ok {
			metrics["template"] = templateName
		}
	case "network", "download":
		if url, ok := context["url"]; ok {
			metrics["target_url"] = url
		}
	}

	return metrics
}

// generatePerformanceSuggestions generates suggestions for improving performance
func (pt *PerformanceTracker) generatePerformanceSuggestions(operation string, duration, threshold time.Duration, context map[string]interface{}) []string {
	var suggestions []string

	// General suggestions
	suggestions = append(suggestions, fmt.Sprintf("Operation took %v, which exceeds the %v threshold", duration, threshold))

	// Operation-specific suggestions
	switch operation {
	case "validation", "validate":
		suggestions = append(suggestions, []string{
			"Consider using --minimal flag for faster validation",
			"Check if large files are slowing down validation",
			"Use --parallel flag if available for concurrent validation",
		}...)
	case "generation", "generate":
		suggestions = append(suggestions, []string{
			"Use --dry-run to test generation without creating files",
			"Consider using a smaller template for faster generation",
			"Check available disk space and I/O performance",
		}...)
	case "network", "download":
		suggestions = append(suggestions, []string{
			"Check network connectivity and bandwidth",
			"Use --offline flag to work with cached data",
			"Consider using a different mirror or CDN",
		}...)
	case "filesystem", "file":
		suggestions = append(suggestions, []string{
			"Check disk I/O performance and available space",
			"Verify file permissions are not causing delays",
			"Consider using SSD storage for better performance",
		}...)
	case "cache":
		suggestions = append(suggestions, []string{
			"Clear cache if it has grown too large",
			"Check cache directory permissions",
			"Consider increasing cache size limits",
		}...)
	default:
		suggestions = append(suggestions, []string{
			"Use --verbose flag to identify bottlenecks",
			"Check system resources (CPU, memory, disk)",
			"Consider running with --debug for detailed timing",
		}...)
	}

	// Context-specific suggestions
	if fileCount, ok := context["file_count"].(int); ok && fileCount > 1000 {
		suggestions = append(suggestions, "Large number of files detected - consider processing in batches")
	}

	if size, ok := context["size"].(int64); ok && size > 100*1024*1024 { // 100MB
		suggestions = append(suggestions, "Large file size detected - consider streaming or chunked processing")
	}

	return suggestions
}

// GetStatistics returns performance statistics
func (pt *PerformanceTracker) GetStatistics() *PerformanceStatistics {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	// Return a copy to avoid race conditions
	stats := *pt.stats
	return &stats
}

// GetOperationHistory returns recent operation history
func (pt *PerformanceTracker) GetOperationHistory(limit int) []OperationRecord {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	if limit <= 0 || limit > len(pt.operationHistory) {
		limit = len(pt.operationHistory)
	}

	// Return the last 'limit' operations
	start := len(pt.operationHistory) - limit
	history := make([]OperationRecord, limit)
	copy(history, pt.operationHistory[start:])

	return history
}

// GetSlowOperations returns operations that exceeded their thresholds
func (pt *PerformanceTracker) GetSlowOperations(limit int) []OperationRecord {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	var slowOps []OperationRecord
	for _, record := range pt.operationHistory {
		if record.WasSlow {
			slowOps = append(slowOps, record)
		}
	}

	// Return the most recent slow operations
	if limit > 0 && limit < len(slowOps) {
		slowOps = slowOps[len(slowOps)-limit:]
	}

	return slowOps
}

// GeneratePerformanceReport generates a comprehensive performance report
func (pt *PerformanceTracker) GeneratePerformanceReport() *PerformanceReport {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	report := &PerformanceReport{
		GeneratedAt:     time.Now(),
		Statistics:      *pt.stats,
		SlowOperations:  pt.GetSlowOperations(10),
		Recommendations: pt.generateRecommendations(),
		TrendAnalysis:   pt.analyzeTrends(),
	}

	return report
}

// PerformanceReport contains a comprehensive performance analysis
type PerformanceReport struct {
	GeneratedAt     time.Time             `json:"generated_at"`
	Statistics      PerformanceStatistics `json:"statistics"`
	SlowOperations  []OperationRecord     `json:"slow_operations"`
	Recommendations []string              `json:"recommendations"`
	TrendAnalysis   *TrendAnalysis        `json:"trend_analysis"`
}

// TrendAnalysis contains performance trend information
type TrendAnalysis struct {
	PerformanceTrend string     `json:"performance_trend"` // "improving", "degrading", "stable"
	TrendConfidence  float64    `json:"trend_confidence"`  // 0.0 to 1.0
	LastImprovement  *time.Time `json:"last_improvement,omitempty"`
	LastDegradation  *time.Time `json:"last_degradation,omitempty"`
}

// generateRecommendations generates performance improvement recommendations
func (pt *PerformanceTracker) generateRecommendations() []string {
	var recommendations []string

	// Analyze slow operation patterns
	if pt.stats.SlowOperations > 0 {
		slowPercentage := float64(pt.stats.SlowOperations) / float64(pt.stats.TotalOperations) * 100

		if slowPercentage > 20 {
			recommendations = append(recommendations, "High percentage of slow operations detected - consider system optimization")
		}

		// Find most problematic operation types
		for opType, count := range pt.stats.SlowOperationsByType {
			if count > 3 {
				recommendations = append(recommendations, fmt.Sprintf("Operation type '%s' frequently slow - investigate bottlenecks", opType))
			}
		}
	}

	// Analyze average performance
	if pt.stats.AverageDuration > 3*time.Second {
		recommendations = append(recommendations, "Average operation duration is high - consider performance optimization")
	}

	// Analyze P95 performance
	if pt.stats.P95Duration > 10*time.Second {
		recommendations = append(recommendations, "95th percentile duration is very high - investigate worst-case scenarios")
	}

	// General recommendations
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance is within acceptable ranges")
	}

	return recommendations
}

// analyzeTrends analyzes performance trends over time
func (pt *PerformanceTracker) analyzeTrends() *TrendAnalysis {
	if len(pt.operationHistory) < 10 {
		return &TrendAnalysis{
			PerformanceTrend: "insufficient_data",
			TrendConfidence:  0.0,
		}
	}

	// Simple trend analysis based on recent vs older operations
	recentCount := len(pt.operationHistory) / 4 // Last 25%
	if recentCount < 5 {
		recentCount = 5
	}

	recentOps := pt.operationHistory[len(pt.operationHistory)-recentCount:]
	olderOps := pt.operationHistory[:len(pt.operationHistory)-recentCount]

	// Calculate average durations
	var recentAvg, olderAvg time.Duration
	for _, op := range recentOps {
		recentAvg += op.Duration
	}
	recentAvg /= time.Duration(len(recentOps))

	for _, op := range olderOps {
		olderAvg += op.Duration
	}
	olderAvg /= time.Duration(len(olderOps))

	// Determine trend
	analysis := &TrendAnalysis{}

	if recentAvg < olderAvg {
		analysis.PerformanceTrend = "improving"
		analysis.TrendConfidence = 0.7
	} else if recentAvg > olderAvg {
		analysis.PerformanceTrend = "degrading"
		analysis.TrendConfidence = 0.7
	} else {
		analysis.PerformanceTrend = "stable"
		analysis.TrendConfidence = 0.8
	}

	return analysis
}

// Reset resets performance tracking statistics
func (pt *PerformanceTracker) Reset() {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	pt.stats = &PerformanceStatistics{
		OperationsByType:     make(map[string]int),
		SlowOperationsByType: make(map[string]int),
		ThresholdViolations:  make(map[string]int),
	}
	pt.operationHistory = make([]OperationRecord, 0, 1000)
}

// Close closes the performance tracker and releases resources
func (pt *PerformanceTracker) Close() error {
	// Generate final report if needed
	if pt.logger != nil && pt.config.VerboseMode {
		report := pt.GeneratePerformanceReport()
		pt.logger.Info("Performance tracking summary: %d operations, %d slow operations (%.1f%%)",
			report.Statistics.TotalOperations,
			report.Statistics.SlowOperations,
			float64(report.Statistics.SlowOperations)/float64(report.Statistics.TotalOperations)*100)
	}

	return nil
}
