// Package performance provides performance analysis functionality for the audit engine.
//
// This package contains specialized analyzers for different aspects of performance
// auditing including bundle size analysis and performance metrics evaluation.
package performance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// BundleAnalyzer provides bundle size analysis functionality.
type BundleAnalyzer struct {
	// Configuration options for bundle analysis
	maxBundleSize      int64   // Maximum recommended bundle size in bytes
	compressionRatio   float64 // Expected compression ratio for estimates
	assetSizeThreshold int64   // Threshold for flagging large assets
}

// NewBundleAnalyzer creates a new bundle analyzer instance.
func NewBundleAnalyzer() *BundleAnalyzer {
	return &BundleAnalyzer{
		maxBundleSize:      5 * 1024 * 1024, // 5MB default
		compressionRatio:   0.3,             // 30% compression ratio
		assetSizeThreshold: 1024 * 1024,     // 1MB threshold for large assets
	}
}

// AnalyzeBundleSize analyzes bundle size and performance characteristics.
func (ba *BundleAnalyzer) AnalyzeBundleSize(path string) (*interfaces.BundleAnalysisResult, error) {
	result := &interfaces.BundleAnalysisResult{
		TotalSize:   0,
		GzippedSize: 0,
		Assets:      []interfaces.BundleAsset{},
		Chunks:      []interfaces.BundleChunk{},
		Summary: interfaces.BundleAnalysisSummary{
			TotalAssets:      0,
			TotalChunks:      0,
			CompressionRatio: 0.0,
			LargestAsset:     "",
			Recommendations:  []string{},
		},
	}

	// Look for common build output directories
	buildDirs := []string{"dist", "build", "public", "static", "assets"}

	for _, buildDir := range buildDirs {
		buildPath := filepath.Join(path, buildDir)
		if _, err := os.Stat(buildPath); err == nil {
			err := ba.analyzeBuildDirectory(buildPath, result)
			if err != nil {
				continue
			}
			break
		}
	}

	// If no build directory found, analyze source files
	if result.TotalSize == 0 {
		ba.analyzeSourceFiles(path, result)
	}

	// Calculate summary statistics
	ba.calculateSummary(result)

	// Generate recommendations
	ba.generateRecommendations(result)

	return result, nil
}

// analyzeBuildDirectory analyzes a build directory for bundle information.
func (ba *BundleAnalyzer) analyzeBuildDirectory(buildPath string, result *interfaces.BundleAnalysisResult) error {
	return filepath.Walk(buildPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		result.TotalSize += info.Size()

		// Estimate gzipped size based on file type
		gzippedSize := ba.estimateGzippedSize(filePath, info.Size())
		result.GzippedSize += gzippedSize

		// Create asset entry
		asset := interfaces.BundleAsset{
			Name:        filepath.Base(filePath),
			Size:        info.Size(),
			GzippedSize: gzippedSize,
			Type:        ba.getAssetType(filePath),
			Percentage:  0, // Will be calculated later
		}

		result.Assets = append(result.Assets, asset)

		return nil
	})
}

// analyzeSourceFiles analyzes source files when no build directory is found.
func (ba *BundleAnalyzer) analyzeSourceFiles(path string, result *interfaces.BundleAnalysisResult) {
	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ba.shouldSkipFile(filePath) || !ba.isSourceCodeFile(filePath) {
			return nil
		}

		result.TotalSize += info.Size()
		gzippedSize := ba.estimateGzippedSize(filePath, info.Size())
		result.GzippedSize += gzippedSize

		asset := interfaces.BundleAsset{
			Name:        filepath.Base(filePath),
			Size:        info.Size(),
			GzippedSize: gzippedSize,
			Type:        ba.getAssetType(filePath),
			Percentage:  0,
		}

		result.Assets = append(result.Assets, asset)

		return nil
	})
}

// calculateSummary calculates summary statistics for the bundle analysis.
func (ba *BundleAnalyzer) calculateSummary(result *interfaces.BundleAnalysisResult) {
	result.Summary.TotalAssets = len(result.Assets)
	result.Summary.TotalChunks = len(result.Chunks)

	if result.TotalSize > 0 && result.GzippedSize > 0 {
		result.Summary.CompressionRatio = float64(result.GzippedSize) / float64(result.TotalSize)
	}

	// Find largest asset
	var largestSize int64
	for _, asset := range result.Assets {
		if asset.Size > largestSize {
			largestSize = asset.Size
			result.Summary.LargestAsset = asset.Name
		}
	}

	// Calculate percentages for each asset
	if result.TotalSize > 0 {
		for i := range result.Assets {
			result.Assets[i].Percentage = float64(result.Assets[i].Size) / float64(result.TotalSize) * 100
		}
	}
}

// generateRecommendations generates optimization recommendations based on analysis.
func (ba *BundleAnalyzer) generateRecommendations(result *interfaces.BundleAnalysisResult) {
	// Check total bundle size
	if result.TotalSize > ba.maxBundleSize {
		result.Summary.Recommendations = append(result.Summary.Recommendations,
			fmt.Sprintf("Bundle size (%d MB) exceeds recommended limit (%d MB). Consider code splitting.",
				result.TotalSize/(1024*1024), ba.maxBundleSize/(1024*1024)))
	}

	// Check compression ratio
	if result.Summary.CompressionRatio > 0.8 {
		result.Summary.Recommendations = append(result.Summary.Recommendations,
			"Poor compression ratio detected. Enable gzip compression on server.")
	}

	// Check for large individual assets
	for _, asset := range result.Assets {
		if asset.Size > ba.assetSizeThreshold {
			result.Summary.Recommendations = append(result.Summary.Recommendations,
				fmt.Sprintf("Large asset detected: %s (%d KB). Consider optimization or lazy loading.",
					asset.Name, asset.Size/1024))
		}
	}

	// Check asset distribution
	jsSize := int64(0)
	cssSize := int64(0)
	imageSize := int64(0)

	for _, asset := range result.Assets {
		switch asset.Type {
		case "javascript":
			jsSize += asset.Size
		case "stylesheet":
			cssSize += asset.Size
		case "image":
			imageSize += asset.Size
		}
	}

	if jsSize > result.TotalSize/2 {
		result.Summary.Recommendations = append(result.Summary.Recommendations,
			"JavaScript assets comprise more than 50% of bundle. Consider code splitting and tree shaking.")
	}

	if imageSize > result.TotalSize/3 {
		result.Summary.Recommendations = append(result.Summary.Recommendations,
			"Images comprise more than 33% of bundle. Consider image optimization and modern formats (WebP, AVIF).")
	}
}

// estimateGzippedSize estimates the gzipped size of a file based on its type and size.
func (ba *BundleAnalyzer) estimateGzippedSize(filePath string, size int64) int64 {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".js", ".jsx", ".ts", ".tsx", ".css", ".scss", ".sass", ".less", ".html", ".htm", ".xml", ".json":
		// Text files compress well
		return int64(float64(size) * ba.compressionRatio)
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".mp3", ".mp4", ".pdf":
		// Already compressed formats
		return size
	default:
		// Default compression ratio
		return int64(float64(size) * 0.6)
	}
}

// getAssetType returns the type of an asset based on its extension.
func (ba *BundleAnalyzer) getAssetType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".js", ".jsx", ".ts", ".tsx":
		return "javascript"
	case ".css", ".scss", ".sass", ".less":
		return "stylesheet"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "image"
	case ".woff", ".woff2", ".ttf", ".eot":
		return "font"
	case ".json":
		return "data"
	case ".html", ".htm":
		return "markup"
	default:
		return "other"
	}
}

// shouldSkipFile determines if a file should be skipped during analysis.
func (ba *BundleAnalyzer) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files for source analysis
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".zip", ".tar", ".gz", ".rar", ".7z",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Skip hidden files and directories
	if strings.HasPrefix(filepath.Base(filePath), ".") {
		return true
	}

	// Skip common directories
	skipDirs := []string{
		"node_modules", "vendor", ".git", ".svn", ".hg",
		"target", "bin", "obj",
	}

	for _, skipDir := range skipDirs {
		if strings.Contains(filePath, skipDir) {
			return true
		}
	}

	return false
}

// isSourceCodeFile checks if a file is a source code file.
func (ba *BundleAnalyzer) isSourceCodeFile(filePath string) bool {
	sourceExtensions := []string{
		".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".cs", ".cpp", ".c", ".h",
		".rb", ".php", ".swift", ".kt", ".rs", ".scala", ".clj", ".hs", ".ml", ".fs",
		".css", ".scss", ".sass", ".less", ".html", ".htm", ".xml", ".yaml", ".yml",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sourceExt := range sourceExtensions {
		if ext == sourceExt {
			return true
		}
	}

	return false
}

// AnalyzePerformanceIssues analyzes performance issues in the project.
func (ba *BundleAnalyzer) AnalyzePerformanceIssues(path string) ([]interfaces.PerformanceIssue, error) {
	var issues []interfaces.PerformanceIssue

	// Check for large files
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ba.shouldSkipFile(filePath) {
			return nil
		}

		// Check for large source files
		if ba.isSourceCodeFile(filePath) && info.Size() > 100*1024 { // 100KB
			issues = append(issues, interfaces.PerformanceIssue{
				Type:        "large_file",
				Severity:    "medium",
				Description: fmt.Sprintf("Large source file: %s (%d bytes)", filepath.Base(filePath), info.Size()),
				Impact:      "May slow down compilation and IDE performance",
				File:        filePath,
			})
		}

		// Check for large assets
		if ba.isAssetFile(filePath) && info.Size() > ba.assetSizeThreshold {
			issues = append(issues, interfaces.PerformanceIssue{
				Type:        "large_asset",
				Severity:    "high",
				Description: fmt.Sprintf("Large asset file: %s (%d bytes)", filepath.Base(filePath), info.Size()),
				Impact:      "May slow down page load times",
				File:        filePath,
			})
		}

		return nil
	})

	// Check for performance anti-patterns in code
	codeIssues, codeErr := ba.analyzeCodePerformanceIssues(path)
	if codeErr == nil {
		issues = append(issues, codeIssues...)
	}

	return issues, err
}

// Pre-compiled regular expressions for performance patterns
var (
	nestedLoopsRegex     = regexp.MustCompile(`for\s*\(.*\)\s*\{`)
	domManipulationRegex = regexp.MustCompile(`\.innerHTML\s*=`)
	inefficientDOMRegex  = regexp.MustCompile(`document\.getElementById`)
	debugStatementsRegex = regexp.MustCompile(`console\.log`)
)

// performancePattern represents a compiled performance pattern
type performancePattern struct {
	regex       *regexp.Regexp
	type_       string
	severity    string
	description string
	impact      string
}

// analyzeCodePerformanceIssues analyzes code for performance anti-patterns.
func (ba *BundleAnalyzer) analyzeCodePerformanceIssues(path string) ([]interfaces.PerformanceIssue, error) {
	var issues []interfaces.PerformanceIssue

	// Pre-compiled performance patterns for better performance
	performancePatterns := []performancePattern{
		{nestedLoopsRegex, "nested_loops", "high", "Nested loops detected", "May cause performance issues"},
		{domManipulationRegex, "dom_manipulation", "medium", "Direct DOM manipulation", "May cause layout thrashing"},
		{inefficientDOMRegex, "inefficient_dom", "low", "Inefficient DOM query", "Consider caching DOM references"},
		{debugStatementsRegex, "debug_statements", "low", "Debug statements in production", "May impact performance"},
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ba.shouldSkipFile(filePath) || !ba.isSourceCodeFile(filePath) {
			return nil
		}

		// #nosec G304 - Audit tool legitimately reads files for analysis
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")

		for i, line := range lines {
			for _, pattern := range performancePatterns {
				if pattern.regex.MatchString(line) {
					issues = append(issues, interfaces.PerformanceIssue{
						Type:        pattern.type_,
						Severity:    pattern.severity,
						Description: pattern.description,
						Impact:      pattern.impact,
						File:        fmt.Sprintf("%s:%d", filePath, i+1),
					})
				}
			}
		}

		return nil
	})

	return issues, err
}

// isAssetFile checks if a file is an asset file.
func (ba *BundleAnalyzer) isAssetFile(filePath string) bool {
	assetExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp",
		".mp3", ".mp4", ".avi", ".mov", ".wmv", ".pdf", ".zip",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, assetExt := range assetExtensions {
		if ext == assetExt {
			return true
		}
	}

	return false
}
