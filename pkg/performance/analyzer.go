package performance

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Name        string
	Iterations  int64
	NsPerOp     float64
	BytesPerOp  int64
	AllocsPerOp int64
	MBPerSec    float64
	Timestamp   time.Time
}

// PerformanceAnalyzer analyzes benchmark results and detects regressions
type PerformanceAnalyzer struct {
	baselineResults map[string]*BenchmarkResult
	currentResults  map[string]*BenchmarkResult
	thresholds      *PerformanceThresholds
}

// PerformanceThresholds defines acceptable performance degradation limits
type PerformanceThresholds struct {
	MaxTimeRegression   float64 // Maximum acceptable time increase (percentage)
	MaxMemoryRegression float64 // Maximum acceptable memory increase (percentage)
	MaxAllocRegression  float64 // Maximum acceptable allocation increase (percentage)
	MinIterations       int64   // Minimum iterations for reliable results
}

// ComparisonResult represents the comparison between baseline and current results
type ComparisonResult struct {
	BenchmarkName       string
	BaselineTime        float64
	CurrentTime         float64
	TimeChange          float64
	TimeChangePercent   float64
	BaselineMemory      int64
	CurrentMemory       int64
	MemoryChange        int64
	MemoryChangePercent float64
	BaselineAllocs      int64
	CurrentAllocs       int64
	AllocsChange        int64
	AllocsChangePercent float64
	IsRegression        bool
	Severity            string
}

// PerformanceReport contains the complete analysis results
type PerformanceReport struct {
	GeneratedAt     time.Time
	TotalBenchmarks int
	Regressions     []*ComparisonResult
	Improvements    []*ComparisonResult
	NewBenchmarks   []*BenchmarkResult
	Summary         *PerformanceSummary
}

// PerformanceSummary provides high-level performance metrics
type PerformanceSummary struct {
	AverageTimeChange   float64
	AverageMemoryChange float64
	AverageAllocsChange float64
	CriticalRegressions int
	MinorRegressions    int
	Improvements        int
	OverallStatus       string
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer() *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		baselineResults: make(map[string]*BenchmarkResult),
		currentResults:  make(map[string]*BenchmarkResult),
		thresholds: &PerformanceThresholds{
			MaxTimeRegression:   10.0, // 10% time increase is concerning
			MaxMemoryRegression: 15.0, // 15% memory increase is concerning
			MaxAllocRegression:  20.0, // 20% allocation increase is concerning
			MinIterations:       100,  // Minimum iterations for reliable results
		},
	}
}

// LoadBaselineResults loads baseline benchmark results from file
func (pa *PerformanceAnalyzer) LoadBaselineResults(filename string) error {
	results, err := pa.parseBenchmarkFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load baseline results: %w", err)
	}

	pa.baselineResults = results
	return nil
}

// LoadCurrentResults loads current benchmark results from file
func (pa *PerformanceAnalyzer) LoadCurrentResults(filename string) error {
	results, err := pa.parseBenchmarkFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load current results: %w", err)
	}

	pa.currentResults = results
	return nil
}

// parseBenchmarkFile parses a benchmark results file
func (pa *PerformanceAnalyzer) parseBenchmarkFile(filename string) (map[string]*BenchmarkResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := make(map[string]*BenchmarkResult)
	scanner := bufio.NewScanner(file)

	// Regex to match benchmark lines
	benchmarkRegex := regexp.MustCompile(`^(Benchmark\w+(?:/\w+)*)-\d+\s+(\d+)\s+([\d.]+)\s+ns/op(?:\s+([\d.]+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?(?:\s+([\d.]+)\s+MB/s)?`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		matches := benchmarkRegex.FindStringSubmatch(line)
		if len(matches) >= 4 {
			result := &BenchmarkResult{
				Name:      matches[1],
				Timestamp: time.Now(),
			}

			// Parse iterations
			if iterations, err := strconv.ParseInt(matches[2], 10, 64); err == nil {
				result.Iterations = iterations
			}

			// Parse ns/op
			if nsPerOp, err := strconv.ParseFloat(matches[3], 64); err == nil {
				result.NsPerOp = nsPerOp
			}

			// Parse B/op (optional)
			if len(matches) > 4 && matches[4] != "" {
				if bytesPerOp, err := strconv.ParseFloat(matches[4], 64); err == nil {
					result.BytesPerOp = int64(bytesPerOp)
				}
			}

			// Parse allocs/op (optional)
			if len(matches) > 5 && matches[5] != "" {
				if allocsPerOp, err := strconv.ParseInt(matches[5], 10, 64); err == nil {
					result.AllocsPerOp = allocsPerOp
				}
			}

			// Parse MB/s (optional)
			if len(matches) > 6 && matches[6] != "" {
				if mbPerSec, err := strconv.ParseFloat(matches[6], 64); err == nil {
					result.MBPerSec = mbPerSec
				}
			}

			results[result.Name] = result
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// AnalyzePerformance compares baseline and current results
func (pa *PerformanceAnalyzer) AnalyzePerformance() (*PerformanceReport, error) {
	if len(pa.baselineResults) == 0 {
		return nil, fmt.Errorf("no baseline results loaded")
	}

	if len(pa.currentResults) == 0 {
		return nil, fmt.Errorf("no current results loaded")
	}

	report := &PerformanceReport{
		GeneratedAt:     time.Now(),
		TotalBenchmarks: len(pa.currentResults),
		Regressions:     make([]*ComparisonResult, 0),
		Improvements:    make([]*ComparisonResult, 0),
		NewBenchmarks:   make([]*BenchmarkResult, 0),
	}

	var totalTimeChange, totalMemoryChange, totalAllocsChange float64
	var comparisonCount int

	// Compare each current result with baseline
	for name, current := range pa.currentResults {
		baseline, exists := pa.baselineResults[name]
		if !exists {
			// New benchmark
			report.NewBenchmarks = append(report.NewBenchmarks, current)
			continue
		}

		// Skip if iterations are too low for reliable comparison
		if current.Iterations < pa.thresholds.MinIterations || baseline.Iterations < pa.thresholds.MinIterations {
			continue
		}

		comparison := pa.compareBenchmarks(name, baseline, current)

		if comparison.IsRegression {
			report.Regressions = append(report.Regressions, comparison)
		} else if comparison.TimeChange < 0 || comparison.MemoryChange < 0 {
			report.Improvements = append(report.Improvements, comparison)
		}

		// Accumulate for averages
		totalTimeChange += comparison.TimeChangePercent
		totalMemoryChange += comparison.MemoryChangePercent
		totalAllocsChange += comparison.AllocsChangePercent
		comparisonCount++
	}

	// Generate summary
	report.Summary = &PerformanceSummary{
		Improvements: len(report.Improvements),
	}

	if comparisonCount > 0 {
		report.Summary.AverageTimeChange = totalTimeChange / float64(comparisonCount)
		report.Summary.AverageMemoryChange = totalMemoryChange / float64(comparisonCount)
		report.Summary.AverageAllocsChange = totalAllocsChange / float64(comparisonCount)
	}

	// Count regression severity
	for _, regression := range report.Regressions {
		if regression.Severity == "critical" {
			report.Summary.CriticalRegressions++
		} else {
			report.Summary.MinorRegressions++
		}
	}

	// Determine overall status
	if report.Summary.CriticalRegressions > 0 {
		report.Summary.OverallStatus = "CRITICAL_REGRESSIONS"
	} else if report.Summary.MinorRegressions > 0 {
		report.Summary.OverallStatus = "MINOR_REGRESSIONS"
	} else if report.Summary.Improvements > 0 {
		report.Summary.OverallStatus = "IMPROVEMENTS"
	} else {
		report.Summary.OverallStatus = "NO_CHANGE"
	}

	return report, nil
}

// compareBenchmarks compares two benchmark results
func (pa *PerformanceAnalyzer) compareBenchmarks(name string, baseline, current *BenchmarkResult) *ComparisonResult {
	comparison := &ComparisonResult{
		BenchmarkName:  name,
		BaselineTime:   baseline.NsPerOp,
		CurrentTime:    current.NsPerOp,
		BaselineMemory: baseline.BytesPerOp,
		CurrentMemory:  current.BytesPerOp,
		BaselineAllocs: baseline.AllocsPerOp,
		CurrentAllocs:  current.AllocsPerOp,
	}

	// Calculate time changes
	comparison.TimeChange = current.NsPerOp - baseline.NsPerOp
	if baseline.NsPerOp > 0 {
		comparison.TimeChangePercent = (comparison.TimeChange / baseline.NsPerOp) * 100
	}

	// Calculate memory changes
	comparison.MemoryChange = current.BytesPerOp - baseline.BytesPerOp
	if baseline.BytesPerOp > 0 {
		comparison.MemoryChangePercent = (float64(comparison.MemoryChange) / float64(baseline.BytesPerOp)) * 100
	}

	// Calculate allocation changes
	comparison.AllocsChange = current.AllocsPerOp - baseline.AllocsPerOp
	if baseline.AllocsPerOp > 0 {
		comparison.AllocsChangePercent = (float64(comparison.AllocsChange) / float64(baseline.AllocsPerOp)) * 100
	}

	// Determine if this is a regression
	isTimeRegression := comparison.TimeChangePercent > pa.thresholds.MaxTimeRegression
	isMemoryRegression := comparison.MemoryChangePercent > pa.thresholds.MaxMemoryRegression
	isAllocRegression := comparison.AllocsChangePercent > pa.thresholds.MaxAllocRegression

	comparison.IsRegression = isTimeRegression || isMemoryRegression || isAllocRegression

	// Determine severity
	if comparison.TimeChangePercent > pa.thresholds.MaxTimeRegression*2 ||
		comparison.MemoryChangePercent > pa.thresholds.MaxMemoryRegression*2 ||
		comparison.AllocsChangePercent > pa.thresholds.MaxAllocRegression*2 {
		comparison.Severity = "critical"
	} else if comparison.IsRegression {
		comparison.Severity = "minor"
	} else {
		comparison.Severity = "none"
	}

	return comparison
}

// GenerateReport generates a formatted performance report
func (pa *PerformanceAnalyzer) GenerateReport(report *PerformanceReport) string {
	var sb strings.Builder

	sb.WriteString("# Performance Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", report.GeneratedAt.Format(time.RFC3339)))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Benchmarks:** %d\n", report.TotalBenchmarks))
	sb.WriteString(fmt.Sprintf("- **Critical Regressions:** %d\n", report.Summary.CriticalRegressions))
	sb.WriteString(fmt.Sprintf("- **Minor Regressions:** %d\n", report.Summary.MinorRegressions))
	sb.WriteString(fmt.Sprintf("- **Improvements:** %d\n", report.Summary.Improvements))
	sb.WriteString(fmt.Sprintf("- **New Benchmarks:** %d\n", len(report.NewBenchmarks)))
	sb.WriteString(fmt.Sprintf("- **Overall Status:** %s\n\n", report.Summary.OverallStatus))

	// Average changes
	sb.WriteString("### Average Performance Changes\n\n")
	sb.WriteString(fmt.Sprintf("- **Time:** %.2f%%\n", report.Summary.AverageTimeChange))
	sb.WriteString(fmt.Sprintf("- **Memory:** %.2f%%\n", report.Summary.AverageMemoryChange))
	sb.WriteString(fmt.Sprintf("- **Allocations:** %.2f%%\n\n", report.Summary.AverageAllocsChange))

	// Critical regressions
	if len(report.Regressions) > 0 {
		sb.WriteString("## Performance Regressions\n\n")

		criticalRegressions := make([]*ComparisonResult, 0)
		minorRegressions := make([]*ComparisonResult, 0)

		for _, regression := range report.Regressions {
			if regression.Severity == "critical" {
				criticalRegressions = append(criticalRegressions, regression)
			} else {
				minorRegressions = append(minorRegressions, regression)
			}
		}

		if len(criticalRegressions) > 0 {
			sb.WriteString("### Critical Regressions ⚠️\n\n")
			sb.WriteString("| Benchmark | Time Change | Memory Change | Allocs Change |\n")
			sb.WriteString("|-----------|-------------|---------------|---------------|\n")

			for _, regression := range criticalRegressions {
				sb.WriteString(fmt.Sprintf("| %s | %.2f%% | %.2f%% | %.2f%% |\n",
					regression.BenchmarkName,
					regression.TimeChangePercent,
					regression.MemoryChangePercent,
					regression.AllocsChangePercent))
			}
			sb.WriteString("\n")
		}

		if len(minorRegressions) > 0 {
			sb.WriteString("### Minor Regressions\n\n")
			sb.WriteString("| Benchmark | Time Change | Memory Change | Allocs Change |\n")
			sb.WriteString("|-----------|-------------|---------------|---------------|\n")

			for _, regression := range minorRegressions {
				sb.WriteString(fmt.Sprintf("| %s | %.2f%% | %.2f%% | %.2f%% |\n",
					regression.BenchmarkName,
					regression.TimeChangePercent,
					regression.MemoryChangePercent,
					regression.AllocsChangePercent))
			}
			sb.WriteString("\n")
		}
	}

	// Improvements
	if len(report.Improvements) > 0 {
		sb.WriteString("## Performance Improvements ✅\n\n")
		sb.WriteString("| Benchmark | Time Improvement | Memory Improvement | Allocs Improvement |\n")
		sb.WriteString("|-----------|------------------|--------------------|--------------------|")

		for _, improvement := range report.Improvements {
			sb.WriteString(fmt.Sprintf("| %s | %.2f%% | %.2f%% | %.2f%% |\n",
				improvement.BenchmarkName,
				-improvement.TimeChangePercent, // Negative because it's an improvement
				-improvement.MemoryChangePercent,
				-improvement.AllocsChangePercent))
		}
		sb.WriteString("\n")
	}

	// New benchmarks
	if len(report.NewBenchmarks) > 0 {
		sb.WriteString("## New Benchmarks\n\n")
		sb.WriteString("| Benchmark | Time (ns/op) | Memory (B/op) | Allocs (allocs/op) |\n")
		sb.WriteString("|-----------|--------------|---------------|--------------------|\n")

		for _, benchmark := range report.NewBenchmarks {
			sb.WriteString(fmt.Sprintf("| %s | %.2f | %d | %d |\n",
				benchmark.Name,
				benchmark.NsPerOp,
				benchmark.BytesPerOp,
				benchmark.AllocsPerOp))
		}
		sb.WriteString("\n")
	}

	// Recommendations
	sb.WriteString("## Recommendations\n\n")

	if report.Summary.CriticalRegressions > 0 {
		sb.WriteString("### Critical Actions Required\n")
		sb.WriteString("- Investigate critical performance regressions immediately\n")
		sb.WriteString("- Consider reverting changes that caused significant slowdowns\n")
		sb.WriteString("- Profile affected components to identify bottlenecks\n\n")
	}

	if report.Summary.MinorRegressions > 0 {
		sb.WriteString("### Minor Optimizations\n")
		sb.WriteString("- Review minor regressions for optimization opportunities\n")
		sb.WriteString("- Consider if trade-offs are acceptable for functionality gains\n\n")
	}

	if len(report.Improvements) > 0 {
		sb.WriteString("### Positive Changes\n")
		sb.WriteString("- Document performance improvements for future reference\n")
		sb.WriteString("- Consider applying similar optimizations to other components\n\n")
	}

	sb.WriteString("---\n")
	sb.WriteString("*Generated by Performance Analysis Tool*\n")

	return sb.String()
}

// SetThresholds allows customizing performance regression thresholds
func (pa *PerformanceAnalyzer) SetThresholds(thresholds *PerformanceThresholds) {
	pa.thresholds = thresholds
}

// GetThresholds returns current performance thresholds
func (pa *PerformanceAnalyzer) GetThresholds() *PerformanceThresholds {
	return pa.thresholds
}
