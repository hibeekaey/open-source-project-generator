// Package quality provides code quality analysis functionality for the audit engine.
//
// This package contains analyzers for measuring test coverage and other
// quality metrics to help identify areas needing better test coverage.
package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// CoverageAnalyzer analyzes test coverage metrics for a project
type CoverageAnalyzer struct {
	minCoverage     float64
	coverageTargets map[string]float64 // file type -> target coverage
	testPatterns    []string
}

// CoverageResult contains test coverage analysis results
type CoverageResult struct {
	OverallCoverage float64                  `json:"overall_coverage"`
	FileCoverage    map[string]float64       `json:"file_coverage"`
	TypeCoverage    map[string]float64       `json:"type_coverage"`
	TestFiles       []string                 `json:"test_files"`
	SourceFiles     []string                 `json:"source_files"`
	UncoveredFiles  []string                 `json:"uncovered_files"`
	Summary         CoverageSummary          `json:"summary"`
	Recommendations []CoverageRecommendation `json:"recommendations"`
}

// CoverageSummary contains coverage statistics
type CoverageSummary struct {
	TotalSourceFiles int     `json:"total_source_files"`
	TotalTestFiles   int     `json:"total_test_files"`
	CoveredFiles     int     `json:"covered_files"`
	UncoveredFiles   int     `json:"uncovered_files"`
	TestRatio        float64 `json:"test_ratio"`
	CoverageGrade    string  `json:"coverage_grade"`
}

// CoverageRecommendation contains coverage improvement recommendations
type CoverageRecommendation struct {
	Type        string   `json:"type"`
	Priority    string   `json:"priority"`
	Description string   `json:"description"`
	Files       []string `json:"files,omitempty"`
	Action      string   `json:"action"`
}

// NewCoverageAnalyzer creates a new coverage analyzer instance
func NewCoverageAnalyzer() *CoverageAnalyzer {
	return &CoverageAnalyzer{
		minCoverage: 80.0, // Default minimum coverage threshold
		coverageTargets: map[string]float64{
			".go":   85.0,
			".js":   80.0,
			".ts":   80.0,
			".py":   85.0,
			".java": 80.0,
			".cs":   75.0,
		},
		testPatterns: []string{
			"*_test.go",
			"*.test.js",
			"*.test.ts",
			"*.spec.js",
			"*.spec.ts",
			"test_*.py",
			"*_test.py",
			"*Test.java",
			"*Tests.java",
			"*Test.cs",
			"*Tests.cs",
		},
	}
}

// AnalyzeTestCoverage analyzes test coverage for a project
func (ca *CoverageAnalyzer) AnalyzeTestCoverage(path string) (*CoverageResult, error) {
	if err := ca.projectExists(path); err != nil {
		return nil, err
	}

	result := &CoverageResult{
		FileCoverage:    make(map[string]float64),
		TypeCoverage:    make(map[string]float64),
		TestFiles:       []string{},
		SourceFiles:     []string{},
		UncoveredFiles:  []string{},
		Recommendations: []CoverageRecommendation{},
	}

	// Collect source and test files
	err := ca.collectFiles(path, result)
	if err != nil {
		return nil, fmt.Errorf("failed to collect files: %w", err)
	}

	// Analyze coverage by file type
	ca.analyzeCoverageByType(result)

	// Calculate overall coverage
	ca.calculateOverallCoverage(result)

	// Generate coverage summary
	ca.generateCoverageSummary(result)

	// Generate recommendations
	ca.generateRecommendations(result)

	return result, nil
}

// collectFiles collects source and test files from the project
func (ca *CoverageAnalyzer) collectFiles(path string, result *CoverageResult) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ca.shouldSkipFile(filePath) {
			return nil
		}

		if ca.isTestFile(filePath) {
			result.TestFiles = append(result.TestFiles, filePath)
		} else if ca.isSourceCodeFile(filePath) {
			result.SourceFiles = append(result.SourceFiles, filePath)

			// Check if this source file has corresponding test
			if !ca.hasCorrespondingTest(filePath, result.TestFiles) {
				result.UncoveredFiles = append(result.UncoveredFiles, filePath)
				result.FileCoverage[filePath] = 0.0
			} else {
				// Estimate coverage based on test file presence and quality
				coverage := ca.estimateFileCoverage(filePath, result.TestFiles)
				result.FileCoverage[filePath] = coverage
			}
		}

		return nil
	})
}

// analyzeCoverageByType analyzes coverage by file type
func (ca *CoverageAnalyzer) analyzeCoverageByType(result *CoverageResult) {
	typeFiles := make(map[string][]string)
	typeCoverage := make(map[string]float64)

	// Group files by type
	for _, file := range result.SourceFiles {
		ext := strings.ToLower(filepath.Ext(file))
		typeFiles[ext] = append(typeFiles[ext], file)
	}

	// Calculate coverage for each type
	for fileType, files := range typeFiles {
		var totalCoverage float64
		coveredFiles := 0

		for _, file := range files {
			if coverage, exists := result.FileCoverage[file]; exists && coverage > 0 {
				totalCoverage += coverage
				coveredFiles++
			}
		}

		if len(files) > 0 {
			typeCoverage[fileType] = (float64(coveredFiles) / float64(len(files))) * 100
		}
	}

	result.TypeCoverage = typeCoverage
}

// calculateOverallCoverage calculates the overall project coverage
func (ca *CoverageAnalyzer) calculateOverallCoverage(result *CoverageResult) {
	if len(result.SourceFiles) == 0 {
		result.OverallCoverage = 0.0
		return
	}

	coveredFiles := 0
	for _, coverage := range result.FileCoverage {
		if coverage > 0 {
			coveredFiles++
		}
	}

	result.OverallCoverage = (float64(coveredFiles) / float64(len(result.SourceFiles))) * 100
}

// generateCoverageSummary generates coverage summary statistics
func (ca *CoverageAnalyzer) generateCoverageSummary(result *CoverageResult) {
	result.Summary = CoverageSummary{
		TotalSourceFiles: len(result.SourceFiles),
		TotalTestFiles:   len(result.TestFiles),
		CoveredFiles:     len(result.SourceFiles) - len(result.UncoveredFiles),
		UncoveredFiles:   len(result.UncoveredFiles),
		CoverageGrade:    ca.getCoverageGrade(result.OverallCoverage),
	}

	if result.Summary.TotalSourceFiles > 0 {
		result.Summary.TestRatio = float64(result.Summary.TotalTestFiles) / float64(result.Summary.TotalSourceFiles)
	}
}

// generateRecommendations generates coverage improvement recommendations
func (ca *CoverageAnalyzer) generateRecommendations(result *CoverageResult) {
	// Recommend adding tests for uncovered files
	if len(result.UncoveredFiles) > 0 {
		result.Recommendations = append(result.Recommendations, CoverageRecommendation{
			Type:        "missing_tests",
			Priority:    "high",
			Description: fmt.Sprintf("Add tests for %d uncovered files", len(result.UncoveredFiles)),
			Files:       result.UncoveredFiles,
			Action:      "Create test files for the listed source files",
		})
	}

	// Recommend improving coverage for low-coverage file types
	for fileType, coverage := range result.TypeCoverage {
		target := ca.coverageTargets[fileType]
		if target == 0 {
			target = ca.minCoverage
		}

		if coverage < target {
			result.Recommendations = append(result.Recommendations, CoverageRecommendation{
				Type:        "low_type_coverage",
				Priority:    "medium",
				Description: fmt.Sprintf("Improve test coverage for %s files (current: %.1f%%, target: %.1f%%)", fileType, coverage, target),
				Action:      fmt.Sprintf("Add more comprehensive tests for %s files", fileType),
			})
		}
	}

	// Recommend improving overall coverage if below threshold
	if result.OverallCoverage < ca.minCoverage {
		result.Recommendations = append(result.Recommendations, CoverageRecommendation{
			Type:        "low_overall_coverage",
			Priority:    "high",
			Description: fmt.Sprintf("Overall test coverage is below target (current: %.1f%%, target: %.1f%%)", result.OverallCoverage, ca.minCoverage),
			Action:      "Focus on adding tests for critical functionality and uncovered code paths",
		})
	}

	// Recommend improving test ratio if too low
	if result.Summary.TestRatio < 0.5 {
		result.Recommendations = append(result.Recommendations, CoverageRecommendation{
			Type:        "low_test_ratio",
			Priority:    "medium",
			Description: fmt.Sprintf("Test-to-source file ratio is low (%.2f)", result.Summary.TestRatio),
			Action:      "Consider adding more test files to improve test coverage",
		})
	}
}

// isTestFile determines if a file is a test file
func (ca *CoverageAnalyzer) isTestFile(filePath string) bool {
	fileName := filepath.Base(filePath)

	for _, pattern := range ca.testPatterns {
		matched, _ := filepath.Match(pattern, fileName)
		if matched {
			return true
		}
	}

	// Additional checks for common test patterns
	lowerName := strings.ToLower(fileName)
	testIndicators := []string{"test", "spec", "tests", "specs"}

	for _, indicator := range testIndicators {
		if strings.Contains(lowerName, indicator) {
			return true
		}
	}

	return false
}

// hasCorrespondingTest checks if a source file has a corresponding test file
func (ca *CoverageAnalyzer) hasCorrespondingTest(sourceFile string, testFiles []string) bool {
	baseName := strings.TrimSuffix(filepath.Base(sourceFile), filepath.Ext(sourceFile))
	sourceDir := filepath.Dir(sourceFile)

	// Common test file naming patterns
	testPatterns := []string{
		baseName + "_test",
		baseName + ".test",
		baseName + "_spec",
		baseName + ".spec",
		"test_" + baseName,
		baseName + "Test",
		baseName + "Tests",
	}

	for _, testFile := range testFiles {
		testBaseName := strings.TrimSuffix(filepath.Base(testFile), filepath.Ext(testFile))
		testDir := filepath.Dir(testFile)

		// Check if test is in same directory or test subdirectory
		if testDir == sourceDir || strings.Contains(testDir, sourceDir) {
			for _, pattern := range testPatterns {
				if strings.EqualFold(testBaseName, pattern) {
					return true
				}
			}
		}
	}

	return false
}

// estimateFileCoverage estimates coverage for a file based on its test files
func (ca *CoverageAnalyzer) estimateFileCoverage(sourceFile string, testFiles []string) float64 {
	// This is a simplified estimation
	// In a real implementation, you would parse actual coverage reports

	baseName := strings.TrimSuffix(filepath.Base(sourceFile), filepath.Ext(sourceFile))
	sourceDir := filepath.Dir(sourceFile)

	var bestTestScore float64 = 0

	for _, testFile := range testFiles {
		testDir := filepath.Dir(testFile)
		testBaseName := strings.TrimSuffix(filepath.Base(testFile), filepath.Ext(testFile))

		// Score based on naming similarity and location
		score := ca.calculateTestScore(baseName, sourceDir, testBaseName, testDir)
		if score > bestTestScore {
			bestTestScore = score
		}
	}

	// Convert score to estimated coverage percentage
	if bestTestScore > 0.8 {
		return 85.0 // High confidence in good coverage
	} else if bestTestScore > 0.6 {
		return 70.0 // Medium confidence
	} else if bestTestScore > 0.3 {
		return 50.0 // Low confidence
	} else {
		return 25.0 // Very low confidence
	}
}

// calculateTestScore calculates a score for how well a test file covers a source file
func (ca *CoverageAnalyzer) calculateTestScore(sourceName, sourceDir, testName, testDir string) float64 {
	score := 0.0

	// Name similarity score
	if strings.Contains(strings.ToLower(testName), strings.ToLower(sourceName)) {
		score += 0.5
	}

	// Directory proximity score
	if testDir == sourceDir {
		score += 0.3
	} else if strings.Contains(testDir, sourceDir) {
		score += 0.2
	}

	// Test pattern matching score
	lowerTestName := strings.ToLower(testName)
	if strings.Contains(lowerTestName, "test") || strings.Contains(lowerTestName, "spec") {
		score += 0.2
	}

	return score
}

// getCoverageGrade returns a letter grade for coverage percentage
func (ca *CoverageAnalyzer) getCoverageGrade(coverage float64) string {
	if coverage >= 90 {
		return "A"
	} else if coverage >= 80 {
		return "B"
	} else if coverage >= 70 {
		return "C"
	} else if coverage >= 60 {
		return "D"
	} else {
		return "F"
	}
}

// projectExists checks if the project path exists and is a directory
func (ca *CoverageAnalyzer) projectExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("project path does not exist: %s", path)
		}
		return fmt.Errorf("unable to access project path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", path)
	}

	return nil
}

// shouldSkipFile determines if a file should be skipped during analysis
func (ca *CoverageAnalyzer) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
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
		"build", "dist", "target", "bin", "obj",
	}

	for _, skipDir := range skipDirs {
		if strings.Contains(filePath, skipDir) {
			return true
		}
	}

	return false
}

// isSourceCodeFile determines if a file is a source code file
func (ca *CoverageAnalyzer) isSourceCodeFile(filePath string) bool {
	sourceExtensions := []string{
		".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".c", ".cpp", ".h", ".hpp",
		".cs", ".php", ".rb", ".swift", ".kt", ".scala", ".rs", ".dart", ".vue",
		".html", ".css", ".scss", ".sass", ".less", ".sql", ".sh", ".bash", ".ps1",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sourceExt := range sourceExtensions {
		if ext == sourceExt {
			return true
		}
	}

	return false
}

// SetMinCoverage sets the minimum coverage threshold
func (ca *CoverageAnalyzer) SetMinCoverage(coverage float64) {
	ca.minCoverage = coverage
}

// GetMinCoverage returns the current minimum coverage threshold
func (ca *CoverageAnalyzer) GetMinCoverage() float64 {
	return ca.minCoverage
}

// SetCoverageTarget sets the coverage target for a specific file type
func (ca *CoverageAnalyzer) SetCoverageTarget(fileType string, target float64) {
	ca.coverageTargets[fileType] = target
}

// GetCoverageTargets returns the current coverage targets
func (ca *CoverageAnalyzer) GetCoverageTargets() map[string]float64 {
	return ca.coverageTargets
}

// AddTestPattern adds a new test file pattern
func (ca *CoverageAnalyzer) AddTestPattern(pattern string) {
	ca.testPatterns = append(ca.testPatterns, pattern)
}

// GetTestPatterns returns the current test file patterns
func (ca *CoverageAnalyzer) GetTestPatterns() []string {
	return ca.testPatterns
}

// AnalyzeCodeSmells analyzes code for quality issues
func (ca *CoverageAnalyzer) AnalyzeCodeSmells(path string) ([]interfaces.CodeSmell, error) {
	var codeSmells []interfaces.CodeSmell

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ca.shouldSkipFile(filePath) || !ca.isSourceCodeFile(filePath) {
			return nil
		}

		smells, err := ca.analyzeFileForCodeSmells(filePath)
		if err != nil {
			return nil // Continue on error
		}

		codeSmells = append(codeSmells, smells...)
		return nil
	})

	return codeSmells, err
}

// Pre-compiled regular expressions for code smell patterns
var (
	technicalDebtRegex     = regexp.MustCompile(`(?i)todo|fixme|hack`)
	debugCodeRegex         = regexp.MustCompile(`console\.log|print\(|println\(`)
	magicNumberRegex       = regexp.MustCompile(`\.length\s*>\s*\d{2,}`)
	longParameterListRegex = regexp.MustCompile(`function\s+\w+\s*\([^)]{50,}`)
	longMethodRegex        = regexp.MustCompile(`\{\s*$[\s\S]{500,}?\}`)
	complexConditionRegex  = regexp.MustCompile(`if\s*\([^)]*&&[^)]*&&[^)]*\)`)
)

// codeSmellPattern represents a compiled code smell pattern
type codeSmellPattern struct {
	regex       *regexp.Regexp
	type_       string
	severity    string
	description string
}

// analyzeFileForCodeSmells analyzes a single file for code smells
func (ca *CoverageAnalyzer) analyzeFileForCodeSmells(filePath string) ([]interfaces.CodeSmell, error) {
	var smells []interfaces.CodeSmell

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	// Pre-compiled code smell patterns for better performance
	smellPatterns := []codeSmellPattern{
		{technicalDebtRegex, "technical_debt", "medium", "TODO/FIXME comment found"},
		{debugCodeRegex, "debug_code", "low", "Debug statement found"},
		{magicNumberRegex, "magic_number", "low", "Magic number used"},
		{longParameterListRegex, "long_parameter_list", "medium", "Long parameter list"},
		{longMethodRegex, "long_method", "medium", "Method too long"},
		{complexConditionRegex, "complex_condition", "medium", "Complex conditional"},
	}

	for i, line := range lines {
		for _, pattern := range smellPatterns {
			if pattern.regex.MatchString(line) {
				smells = append(smells, interfaces.CodeSmell{
					Type:        pattern.type_,
					Severity:    pattern.severity,
					Description: pattern.description,
					File:        filePath,
					Line:        i + 1,
				})
			}
		}

		// Check line length
		if len(line) > 120 {
			smells = append(smells, interfaces.CodeSmell{
				Type:        "long_line",
				Severity:    "low",
				Description: fmt.Sprintf("Line too long (%d characters)", len(line)),
				File:        filePath,
				Line:        i + 1,
			})
		}
	}

	return smells, nil
}

// AnalyzeDuplications analyzes code for duplications
func (ca *CoverageAnalyzer) AnalyzeDuplications(path string) ([]interfaces.Duplication, error) {
	var duplications []interfaces.Duplication

	// This is a simplified implementation
	// In a real implementation, you would use more sophisticated algorithms
	// like suffix trees or rolling hashes to detect duplications

	fileContents := make(map[string][]string)

	// Read all source files
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ca.shouldSkipFile(filePath) || !ca.isSourceCodeFile(filePath) {
			return nil
		}

		// #nosec G304 - Audit tool legitimately reads files for analysis
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		fileContents[filePath] = lines
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Simple duplication detection (looking for identical blocks of 5+ lines)
	const minDuplicationLines = 5

	for file1, lines1 := range fileContents {
		for file2, lines2 := range fileContents {
			if file1 >= file2 { // Avoid duplicate comparisons
				continue
			}

			duplications = append(duplications, ca.findDuplicationsInFiles(file1, lines1, file2, lines2, minDuplicationLines)...)
		}
	}

	return duplications, nil
}

// findDuplicationsInFiles finds duplications between two files
func (ca *CoverageAnalyzer) findDuplicationsInFiles(file1 string, lines1 []string, file2 string, lines2 []string, minLines int) []interfaces.Duplication {
	var duplications []interfaces.Duplication

	// Simple sliding window approach
	for i := 0; i <= len(lines1)-minLines; i++ {
		for j := 0; j <= len(lines2)-minLines; j++ {
			matchLength := 0

			// Count consecutive matching lines
			for k := 0; i+k < len(lines1) && j+k < len(lines2); k++ {
				line1 := strings.TrimSpace(lines1[i+k])
				line2 := strings.TrimSpace(lines2[j+k])

				if line1 == line2 && line1 != "" {
					matchLength++
				} else {
					break
				}
			}

			if matchLength >= minLines {
				duplications = append(duplications, interfaces.Duplication{
					Files:      []string{file1, file2},
					Lines:      matchLength,
					Tokens:     matchLength * 10, // Rough estimate
					Percentage: float64(matchLength) / float64(len(lines1)) * 100,
				})
			}
		}
	}

	return duplications
}
