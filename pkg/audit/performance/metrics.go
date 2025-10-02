package performance

import (
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// MetricsAnalyzer provides performance metrics analysis functionality.
type MetricsAnalyzer struct {
	bundleAnalyzer *BundleAnalyzer

	// Performance thresholds
	maxLoadTime       time.Duration
	maxFirstPaint     time.Duration
	maxTimeToInteract time.Duration
	maxLayoutShift    float64
}

// NewMetricsAnalyzer creates a new metrics analyzer instance.
func NewMetricsAnalyzer() *MetricsAnalyzer {
	return &MetricsAnalyzer{
		bundleAnalyzer:    NewBundleAnalyzer(),
		maxLoadTime:       2 * time.Second,
		maxFirstPaint:     1 * time.Second,
		maxTimeToInteract: 3 * time.Second,
		maxLayoutShift:    0.1,
	}
}

// CheckPerformanceMetrics analyzes performance metrics for a project.
func (ma *MetricsAnalyzer) CheckPerformanceMetrics(path string) (*interfaces.PerformanceMetricsResult, error) {
	result := &interfaces.PerformanceMetricsResult{
		LoadTime:              0,
		FirstPaint:            0,
		FirstContentful:       0,
		LargestContentful:     0,
		TimeToInteractive:     0,
		CumulativeLayoutShift: 0.0,
		Issues:                []interfaces.PerformanceIssue{},
		Summary: interfaces.PerformanceMetricsSummary{
			OverallScore:     100.0,
			PerformanceGrade: "A",
			IssueCount:       0,
			Recommendations:  []string{},
		},
	}

	// This is a simplified implementation that estimates metrics
	// In a real implementation, you would integrate with tools like:
	// - Lighthouse for web performance
	// - WebPageTest
	// - Custom performance monitoring

	// Estimate metrics based on bundle size and project structure
	bundleResult, err := ma.bundleAnalyzer.AnalyzeBundleSize(path)
	if err == nil {
		ma.estimateMetricsFromBundle(bundleResult, result)
	}

	// Analyze performance issues
	issues, err := ma.bundleAnalyzer.AnalyzePerformanceIssues(path)
	if err == nil {
		result.Issues = issues
	}

	// Calculate performance score and grade
	ma.calculatePerformanceScore(result)

	// Generate recommendations
	ma.generatePerformanceRecommendations(result)

	result.Summary.IssueCount = len(result.Issues)

	return result, nil
}

// estimateMetricsFromBundle estimates performance metrics based on bundle analysis.
func (ma *MetricsAnalyzer) estimateMetricsFromBundle(bundleResult *interfaces.BundleAnalysisResult, result *interfaces.PerformanceMetricsResult) {
	// Rough estimates based on bundle size
	sizeInMB := float64(bundleResult.TotalSize) / (1024 * 1024)

	// Estimate load time based on bundle size
	// Assumes average connection speed and processing time
	baseLoadTime := time.Duration(sizeInMB*100) * time.Millisecond

	// Add network latency estimate
	networkLatency := 200 * time.Millisecond
	result.LoadTime = baseLoadTime + networkLatency

	// Estimate paint metrics
	result.FirstPaint = result.LoadTime / 3
	result.FirstContentful = result.LoadTime / 2
	result.LargestContentful = result.LoadTime * 4 / 5
	result.TimeToInteractive = result.LoadTime * 3 / 2

	// Estimate layout shift based on asset types
	imageAssets := 0
	for _, asset := range bundleResult.Assets {
		if asset.Type == "image" {
			imageAssets++
		}
	}

	// More images typically mean higher chance of layout shift
	if imageAssets > 10 {
		result.CumulativeLayoutShift = 0.15
	} else if imageAssets > 5 {
		result.CumulativeLayoutShift = 0.08
	} else {
		result.CumulativeLayoutShift = 0.03
	}
}

// calculatePerformanceScore calculates the overall performance score and grade.
func (ma *MetricsAnalyzer) calculatePerformanceScore(result *interfaces.PerformanceMetricsResult) {
	score := 100.0

	// Deduct points based on load time
	if result.LoadTime > ma.maxLoadTime {
		excess := result.LoadTime - ma.maxLoadTime
		deduction := float64(excess.Milliseconds()) / 100 // 1 point per 100ms over limit
		score -= deduction
	}

	// Deduct points based on first paint
	if result.FirstPaint > ma.maxFirstPaint {
		excess := result.FirstPaint - ma.maxFirstPaint
		deduction := float64(excess.Milliseconds()) / 50 // 1 point per 50ms over limit
		score -= deduction
	}

	// Deduct points based on time to interactive
	if result.TimeToInteractive > ma.maxTimeToInteract {
		excess := result.TimeToInteractive - ma.maxTimeToInteract
		deduction := float64(excess.Milliseconds()) / 200 // 1 point per 200ms over limit
		score -= deduction
	}

	// Deduct points based on layout shift
	if result.CumulativeLayoutShift > ma.maxLayoutShift {
		excess := result.CumulativeLayoutShift - ma.maxLayoutShift
		deduction := excess * 100 // 100 points per 0.01 CLS over limit
		score -= deduction
	}

	// Deduct points for performance issues
	for _, issue := range result.Issues {
		switch issue.Severity {
		case "critical":
			score -= 20
		case "high":
			score -= 10
		case "medium":
			score -= 5
		case "low":
			score -= 2
		}
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	result.Summary.OverallScore = score

	// Assign grade based on score
	switch {
	case score >= 90:
		result.Summary.PerformanceGrade = "A"
	case score >= 80:
		result.Summary.PerformanceGrade = "B"
	case score >= 70:
		result.Summary.PerformanceGrade = "C"
	case score >= 60:
		result.Summary.PerformanceGrade = "D"
	default:
		result.Summary.PerformanceGrade = "F"
	}
}

// generatePerformanceRecommendations generates performance optimization recommendations.
func (ma *MetricsAnalyzer) generatePerformanceRecommendations(result *interfaces.PerformanceMetricsResult) {
	var recommendations []string

	// Load time recommendations
	if result.LoadTime > ma.maxLoadTime {
		recommendations = append(recommendations,
			fmt.Sprintf("Load time (%v) exceeds recommended %v. Consider optimizing bundle size and enabling code splitting.",
				result.LoadTime, ma.maxLoadTime))
	}

	// First paint recommendations
	if result.FirstPaint > ma.maxFirstPaint {
		recommendations = append(recommendations,
			"First paint time is slow. Consider inlining critical CSS and optimizing above-the-fold content.")
	}

	// Time to interactive recommendations
	if result.TimeToInteractive > ma.maxTimeToInteract {
		recommendations = append(recommendations,
			"Time to interactive is slow. Consider reducing JavaScript execution time and deferring non-critical scripts.")
	}

	// Layout shift recommendations
	if result.CumulativeLayoutShift > ma.maxLayoutShift {
		recommendations = append(recommendations,
			"High cumulative layout shift detected. Ensure images and ads have defined dimensions.")
	}

	// Issue-specific recommendations
	hasLargeAssets := false
	hasNestedLoops := false
	hasDOMManipulation := false

	for _, issue := range result.Issues {
		switch issue.Type {
		case "large_asset":
			if !hasLargeAssets {
				recommendations = append(recommendations,
					"Large assets detected. Consider image optimization, lazy loading, and modern formats (WebP, AVIF).")
				hasLargeAssets = true
			}
		case "nested_loops":
			if !hasNestedLoops {
				recommendations = append(recommendations,
					"Nested loops detected. Consider algorithm optimization and data structure improvements.")
				hasNestedLoops = true
			}
		case "dom_manipulation":
			if !hasDOMManipulation {
				recommendations = append(recommendations,
					"Direct DOM manipulation detected. Consider using virtual DOM or batching DOM updates.")
				hasDOMManipulation = true
			}
		}
	}

	// General recommendations based on grade
	switch result.Summary.PerformanceGrade {
	case "D", "F":
		recommendations = append(recommendations,
			"Performance is poor. Consider a comprehensive performance audit and optimization strategy.")
	case "C":
		recommendations = append(recommendations,
			"Performance needs improvement. Focus on the highest impact optimizations first.")
	case "B":
		recommendations = append(recommendations,
			"Good performance with room for improvement. Consider advanced optimization techniques.")
	}

	result.Summary.Recommendations = recommendations
}

// EstimateLoadTime estimates load time based on various factors.
func (ma *MetricsAnalyzer) EstimateLoadTime(bundleSize int64, connectionType string) time.Duration {
	// Connection speed estimates (bytes per second)
	connectionSpeeds := map[string]int64{
		"slow-3g": 50 * 1024,        // 50 KB/s
		"3g":      100 * 1024,       // 100 KB/s
		"4g":      1024 * 1024,      // 1 MB/s
		"wifi":    5 * 1024 * 1024,  // 5 MB/s
		"cable":   10 * 1024 * 1024, // 10 MB/s
	}

	speed, exists := connectionSpeeds[connectionType]
	if !exists {
		speed = connectionSpeeds["4g"] // Default to 4G
	}

	// Calculate download time
	downloadTime := time.Duration(bundleSize/speed) * time.Second

	// Add processing time (estimated 10% of download time)
	processingTime := downloadTime / 10

	// Add network latency
	latency := 200 * time.Millisecond

	return downloadTime + processingTime + latency
}

// AnalyzeWebVitals analyzes Core Web Vitals metrics.
func (ma *MetricsAnalyzer) AnalyzeWebVitals(result *interfaces.PerformanceMetricsResult) map[string]string {
	vitals := make(map[string]string)

	// Largest Contentful Paint (LCP)
	if result.LargestContentful <= 2500*time.Millisecond {
		vitals["LCP"] = "good"
	} else if result.LargestContentful <= 4000*time.Millisecond {
		vitals["LCP"] = "needs-improvement"
	} else {
		vitals["LCP"] = "poor"
	}

	// First Input Delay (FID) - estimated based on time to interactive
	fid := result.TimeToInteractive - result.FirstContentful
	if fid <= 100*time.Millisecond {
		vitals["FID"] = "good"
	} else if fid <= 300*time.Millisecond {
		vitals["FID"] = "needs-improvement"
	} else {
		vitals["FID"] = "poor"
	}

	// Cumulative Layout Shift (CLS)
	if result.CumulativeLayoutShift <= 0.1 {
		vitals["CLS"] = "good"
	} else if result.CumulativeLayoutShift <= 0.25 {
		vitals["CLS"] = "needs-improvement"
	} else {
		vitals["CLS"] = "poor"
	}

	return vitals
}

// GetPerformanceBudget returns recommended performance budgets.
func (ma *MetricsAnalyzer) GetPerformanceBudget() map[string]interface{} {
	return map[string]interface{}{
		"max_bundle_size_mb":      5,
		"max_load_time_ms":        2000,
		"max_first_paint_ms":      1000,
		"max_time_interactive_ms": 3000,
		"max_layout_shift":        0.1,
		"max_asset_size_mb":       1,
		"max_js_bundle_mb":        2,
		"max_css_bundle_kb":       100,
		"max_image_size_kb":       500,
	}
}
