package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/performance"
)

func main() {
	var (
		baselineFile   = flag.String("baseline", "benchmark_results/baseline_results.txt", "Path to baseline benchmark results")
		currentFile    = flag.String("current", "benchmark_results/current_results.txt", "Path to current benchmark results")
		outputFile     = flag.String("output", "benchmark_results/analysis_report.md", "Path to output analysis report")
		format         = flag.String("format", "markdown", "Output format (markdown, json, text)")
		thresholdTime  = flag.Float64("threshold-time", 10.0, "Time regression threshold percentage")
		thresholdMem   = flag.Float64("threshold-memory", 15.0, "Memory regression threshold percentage")
		thresholdAlloc = flag.Float64("threshold-allocs", 20.0, "Allocation regression threshold percentage")
		verbose        = flag.Bool("verbose", false, "Enable verbose output")
		help           = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Starting performance analysis...")
	}

	// Create analyzer
	analyzer := performance.NewPerformanceAnalyzer()

	// Set custom thresholds if provided
	thresholds := &performance.PerformanceThresholds{
		MaxTimeRegression:   *thresholdTime,
		MaxMemoryRegression: *thresholdMem,
		MaxAllocRegression:  *thresholdAlloc,
		MinIterations:       100,
	}
	analyzer.SetThresholds(thresholds)

	// Load baseline results
	if *verbose {
		log.Printf("Loading baseline results from: %s", *baselineFile)
	}

	if err := analyzer.LoadBaselineResults(*baselineFile); err != nil {
		log.Fatalf("Failed to load baseline results: %v", err)
	}

	// Load current results
	if *verbose {
		log.Printf("Loading current results from: %s", *currentFile)
	}

	if err := analyzer.LoadCurrentResults(*currentFile); err != nil {
		log.Fatalf("Failed to load current results: %v", err)
	}

	// Analyze performance
	if *verbose {
		log.Println("Analyzing performance differences...")
	}

	report, err := analyzer.AnalyzePerformance()
	if err != nil {
		log.Fatalf("Failed to analyze performance: %v", err)
	}

	// Generate output
	var output string
	switch *format {
	case "markdown", "md":
		output = analyzer.GenerateReport(report)
	case "json":
		output = generateJSONReport(report)
	case "text":
		output = generateTextReport(report)
	default:
		log.Fatalf("Unsupported output format: %s", *format)
	}

	// Write output
	if *outputFile == "-" {
		fmt.Print(output)
	} else {
		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(*outputFile), 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		if err := os.WriteFile(*outputFile, []byte(output), 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}

		if *verbose {
			log.Printf("Analysis report written to: %s", *outputFile)
		}
	}

	// Print summary to stderr for CI/CD integration
	printSummary(report)

	// Exit with appropriate code
	if report.Summary.CriticalRegressions > 0 {
		os.Exit(1) // Critical regressions found
	} else if report.Summary.MinorRegressions > 0 {
		os.Exit(2) // Minor regressions found
	}
	// Exit 0 for no regressions or improvements
}

func showHelp() {
	fmt.Println("Performance Analyzer - Compare benchmark results and detect regressions")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  performance-analyzer [options]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Basic analysis with default files")
	fmt.Println("  performance-analyzer")
	fmt.Println()
	fmt.Println("  # Custom files and thresholds")
	fmt.Println("  performance-analyzer -baseline old.txt -current new.txt -threshold-time 5.0")
	fmt.Println()
	fmt.Println("  # Output to stdout in JSON format")
	fmt.Println("  performance-analyzer -output - -format json")
	fmt.Println()
	fmt.Println("Exit Codes:")
	fmt.Println("  0 - No regressions or improvements found")
	fmt.Println("  1 - Critical regressions found")
	fmt.Println("  2 - Minor regressions found")
}

func generateJSONReport(report *performance.PerformanceReport) string {
	// Simple JSON generation (in a real implementation, use json.Marshal)
	return fmt.Sprintf(`{
  "generated_at": "%s",
  "total_benchmarks": %d,
  "critical_regressions": %d,
  "minor_regressions": %d,
  "improvements": %d,
  "new_benchmarks": %d,
  "overall_status": "%s",
  "average_time_change": %.2f,
  "average_memory_change": %.2f,
  "average_allocs_change": %.2f
}`,
		report.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		report.TotalBenchmarks,
		report.Summary.CriticalRegressions,
		report.Summary.MinorRegressions,
		report.Summary.Improvements,
		len(report.NewBenchmarks),
		report.Summary.OverallStatus,
		report.Summary.AverageTimeChange,
		report.Summary.AverageMemoryChange,
		report.Summary.AverageAllocsChange)
}

func generateTextReport(report *performance.PerformanceReport) string {
	output := fmt.Sprintf("Performance Analysis Report\n")
	output += fmt.Sprintf("Generated: %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))

	output += fmt.Sprintf("Summary:\n")
	output += fmt.Sprintf("  Total Benchmarks: %d\n", report.TotalBenchmarks)
	output += fmt.Sprintf("  Critical Regressions: %d\n", report.Summary.CriticalRegressions)
	output += fmt.Sprintf("  Minor Regressions: %d\n", report.Summary.MinorRegressions)
	output += fmt.Sprintf("  Improvements: %d\n", report.Summary.Improvements)
	output += fmt.Sprintf("  New Benchmarks: %d\n", len(report.NewBenchmarks))
	output += fmt.Sprintf("  Overall Status: %s\n\n", report.Summary.OverallStatus)

	output += fmt.Sprintf("Average Changes:\n")
	output += fmt.Sprintf("  Time: %.2f%%\n", report.Summary.AverageTimeChange)
	output += fmt.Sprintf("  Memory: %.2f%%\n", report.Summary.AverageMemoryChange)
	output += fmt.Sprintf("  Allocations: %.2f%%\n\n", report.Summary.AverageAllocsChange)

	if len(report.Regressions) > 0 {
		output += fmt.Sprintf("Regressions:\n")
		for _, regression := range report.Regressions {
			output += fmt.Sprintf("  %s [%s]: Time %.2f%%, Memory %.2f%%, Allocs %.2f%%\n",
				regression.BenchmarkName,
				regression.Severity,
				regression.TimeChangePercent,
				regression.MemoryChangePercent,
				regression.AllocsChangePercent)
		}
		output += "\n"
	}

	if len(report.Improvements) > 0 {
		output += fmt.Sprintf("Improvements:\n")
		for _, improvement := range report.Improvements {
			output += fmt.Sprintf("  %s: Time %.2f%%, Memory %.2f%%, Allocs %.2f%%\n",
				improvement.BenchmarkName,
				-improvement.TimeChangePercent,
				-improvement.MemoryChangePercent,
				-improvement.AllocsChangePercent)
		}
		output += "\n"
	}

	return output
}

func printSummary(report *performance.PerformanceReport) {
	fmt.Fprintf(os.Stderr, "Performance Analysis Summary:\n")
	fmt.Fprintf(os.Stderr, "  Status: %s\n", report.Summary.OverallStatus)
	fmt.Fprintf(os.Stderr, "  Critical Regressions: %d\n", report.Summary.CriticalRegressions)
	fmt.Fprintf(os.Stderr, "  Minor Regressions: %d\n", report.Summary.MinorRegressions)
	fmt.Fprintf(os.Stderr, "  Improvements: %d\n", report.Summary.Improvements)

	if report.Summary.CriticalRegressions > 0 {
		fmt.Fprintf(os.Stderr, "❌ CRITICAL: Performance regressions detected!\n")
	} else if report.Summary.MinorRegressions > 0 {
		fmt.Fprintf(os.Stderr, "⚠️  WARNING: Minor performance regressions detected\n")
	} else if report.Summary.Improvements > 0 {
		fmt.Fprintf(os.Stderr, "✅ SUCCESS: Performance improvements detected\n")
	} else {
		fmt.Fprintf(os.Stderr, "ℹ️  INFO: No significant performance changes\n")
	}
}
