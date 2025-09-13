package cleanup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// TODOScanner provides enhanced TODO/FIXME comment analysis
type TODOScanner struct {
	config *TODOScanConfig
}

// TODOScanConfig holds configuration for TODO scanning
type TODOScanConfig struct {
	SkipPatterns    []string
	IncludePatterns []string
	CustomKeywords  []string
	OutputFormat    string // "text", "json", "markdown"
}

// TODOReport represents a comprehensive report of all TODO items
type TODOReport struct {
	Timestamp         time.Time
	ProjectRoot       string
	TotalFiles        int
	FilesScanned      int
	TODOs             []TODOItem
	Summary           *TODOSummary
	CategoryBreakdown map[Category][]TODOItem
	PriorityBreakdown map[Priority][]TODOItem
	FileBreakdown     map[string][]TODOItem
}

// TODOSummary provides summary statistics
type TODOSummary struct {
	TotalTODOs       int
	SecurityTODOs    int
	PerformanceTODOs int
	FeatureTODOs     int
	BugTODOs         int
	CriticalTODOs    int
	HighTODOs        int
	MediumTODOs      int
	LowTODOs         int
}

// NewTODOScanner creates a new TODO scanner
func NewTODOScanner(config *TODOScanConfig) *TODOScanner {
	if config == nil {
		config = DefaultTODOScanConfig()
	}
	return &TODOScanner{config: config}
}

// DefaultTODOScanConfig returns default configuration
func DefaultTODOScanConfig() *TODOScanConfig {
	return &TODOScanConfig{
		SkipPatterns:    []string{"vendor/", ".git/", "node_modules/", ".cleanup-backups/"},
		IncludePatterns: []string{"*.go", "*.md", "*.yaml", "*.yml", "*.json"},
		CustomKeywords:  []string{"TODO", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},
		OutputFormat:    "markdown",
	}
}

// ScanProject scans the entire project for TODO comments
func (ts *TODOScanner) ScanProject(rootDir string) (*TODOReport, error) {
	report := &TODOReport{
		Timestamp:         time.Now(),
		ProjectRoot:       rootDir,
		TODOs:             []TODOItem{},
		CategoryBreakdown: make(map[Category][]TODOItem),
		PriorityBreakdown: make(map[Priority][]TODOItem),
		FileBreakdown:     make(map[string][]TODOItem),
	}

	// Create regex pattern for TODO detection
	keywords := strings.Join(ts.config.CustomKeywords, "|")
	todoRegex := regexp.MustCompile(fmt.Sprintf(`(?i)(%s)[\s:]*(.*)`, keywords))

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and files that don't match patterns
		if info.IsDir() || ts.shouldSkipFile(path) {
			return nil
		}

		report.TotalFiles++

		// Only scan text files
		if !ts.isTextFile(path) {
			return nil
		}

		report.FilesScanned++

		todos, err := ts.scanFile(path, todoRegex)
		if err != nil {
			return fmt.Errorf("failed to scan file %s: %w", path, err)
		}

		report.TODOs = append(report.TODOs, todos...)

		// Update breakdowns
		for _, todo := range todos {
			report.CategoryBreakdown[todo.Category] = append(report.CategoryBreakdown[todo.Category], todo)
			report.PriorityBreakdown[todo.Priority] = append(report.PriorityBreakdown[todo.Priority], todo)
			report.FileBreakdown[todo.File] = append(report.FileBreakdown[todo.File], todo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Generate summary
	report.Summary = ts.generateSummary(report.TODOs)

	// Sort TODOs by priority and then by file
	sort.Slice(report.TODOs, func(i, j int) bool {
		if report.TODOs[i].Priority != report.TODOs[j].Priority {
			return report.TODOs[i].Priority > report.TODOs[j].Priority // Higher priority first
		}
		return report.TODOs[i].File < report.TODOs[j].File
	})

	return report, nil
}

// scanFile scans a single file for TODO comments
func (ts *TODOScanner) scanFile(filePath string, todoRegex *regexp.Regexp) ([]TODOItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var todos []TODOItem
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if matches := todoRegex.FindStringSubmatch(line); matches != nil {
			todo := TODOItem{
				File:     filePath,
				Line:     lineNum,
				Type:     strings.ToUpper(matches[1]),
				Message:  strings.TrimSpace(matches[2]),
				Context:  strings.TrimSpace(line),
				Priority: ts.determinePriority(matches[1], matches[2], line),
				Category: ts.determineCategory(matches[2], line, filePath),
			}
			todos = append(todos, todo)
		}
	}

	return todos, scanner.Err()
}

// shouldSkipFile determines if a file should be skipped
func (ts *TODOScanner) shouldSkipFile(path string) bool {
	for _, pattern := range ts.config.SkipPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

// isTextFile determines if a file is a text file that should be scanned
func (ts *TODOScanner) isTextFile(path string) bool {
	for _, pattern := range ts.config.IncludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// determinePriority determines the priority of a TODO item
func (ts *TODOScanner) determinePriority(todoType, message, context string) Priority {
	message = strings.ToLower(message)
	context = strings.ToLower(context)
	todoType = strings.ToLower(todoType)

	// Critical priority indicators
	if todoType == "fixme" || todoType == "bug" ||
		strings.Contains(message, "security") ||
		strings.Contains(message, "vulnerability") ||
		strings.Contains(message, "critical") ||
		strings.Contains(message, "urgent") ||
		strings.Contains(context, "security") {
		return PriorityCritical
	}

	// High priority indicators
	if todoType == "hack" ||
		strings.Contains(message, "performance") ||
		strings.Contains(message, "memory leak") ||
		strings.Contains(message, "deadlock") ||
		strings.Contains(message, "race condition") ||
		strings.Contains(message, "important") {
		return PriorityHigh
	}

	// Medium priority indicators
	if strings.Contains(message, "refactor") ||
		strings.Contains(message, "cleanup") ||
		strings.Contains(message, "optimize") ||
		strings.Contains(message, "improve") ||
		todoType == "optimize" {
		return PriorityMedium
	}

	// Default to low priority
	return PriorityLow
}

// determineCategory determines the category of a TODO item
func (ts *TODOScanner) determineCategory(message, context, filePath string) Category {
	message = strings.ToLower(message)
	context = strings.ToLower(context)
	filePath = strings.ToLower(filePath)

	// Security category
	if strings.Contains(message, "security") ||
		strings.Contains(message, "vulnerability") ||
		strings.Contains(message, "auth") ||
		strings.Contains(message, "encrypt") ||
		strings.Contains(message, "secure") ||
		strings.Contains(context, "security") ||
		strings.Contains(filePath, "security") {
		return CategorySecurity
	}

	// Performance category
	if strings.Contains(message, "performance") ||
		strings.Contains(message, "optimize") ||
		strings.Contains(message, "speed") ||
		strings.Contains(message, "memory") ||
		strings.Contains(message, "cache") ||
		strings.Contains(message, "slow") {
		return CategoryPerformance
	}

	// Documentation category
	if strings.Contains(message, "doc") ||
		strings.Contains(message, "comment") ||
		strings.Contains(message, "documentation") ||
		strings.Contains(filePath, "doc") ||
		strings.HasSuffix(filePath, ".md") {
		return CategoryDocumentation
	}

	// Bug category
	if strings.Contains(message, "bug") ||
		strings.Contains(message, "fix") ||
		strings.Contains(message, "error") ||
		strings.Contains(message, "issue") ||
		strings.Contains(message, "broken") {
		return CategoryBug
	}

	// Refactor category
	if strings.Contains(message, "refactor") ||
		strings.Contains(message, "cleanup") ||
		strings.Contains(message, "reorganize") ||
		strings.Contains(message, "restructure") {
		return CategoryRefactor
	}

	// Default to feature
	return CategoryFeature
}

// generateSummary generates summary statistics
func (ts *TODOScanner) generateSummary(todos []TODOItem) *TODOSummary {
	summary := &TODOSummary{}

	for _, todo := range todos {
		summary.TotalTODOs++

		// Count by category
		switch todo.Category {
		case CategorySecurity:
			summary.SecurityTODOs++
		case CategoryPerformance:
			summary.PerformanceTODOs++
		case CategoryFeature:
			summary.FeatureTODOs++
		case CategoryBug:
			summary.BugTODOs++
		}

		// Count by priority
		switch todo.Priority {
		case PriorityCritical:
			summary.CriticalTODOs++
		case PriorityHigh:
			summary.HighTODOs++
		case PriorityMedium:
			summary.MediumTODOs++
		case PriorityLow:
			summary.LowTODOs++
		}
	}

	return summary
}

// GenerateReport generates a formatted report
func (ts *TODOScanner) GenerateReport(report *TODOReport) (string, error) {
	switch ts.config.OutputFormat {
	case "json":
		return ts.generateJSONReport(report)
	case "text":
		return ts.generateTextReport(report)
	default:
		return ts.generateMarkdownReport(report)
	}
}

// generateMarkdownReport generates a markdown report
func (ts *TODOScanner) generateMarkdownReport(report *TODOReport) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# TODO/FIXME Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Project Root:** %s\n", report.ProjectRoot))
	sb.WriteString(fmt.Sprintf("**Files Scanned:** %d/%d\n\n", report.FilesScanned, report.TotalFiles))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Total TODOs:** %d\n", report.Summary.TotalTODOs))
	sb.WriteString(fmt.Sprintf("- **Critical:** %d\n", report.Summary.CriticalTODOs))
	sb.WriteString(fmt.Sprintf("- **High Priority:** %d\n", report.Summary.HighTODOs))
	sb.WriteString(fmt.Sprintf("- **Medium Priority:** %d\n", report.Summary.MediumTODOs))
	sb.WriteString(fmt.Sprintf("- **Low Priority:** %d\n", report.Summary.LowTODOs))
	sb.WriteString("\n")

	// Category breakdown
	sb.WriteString("### By Category\n\n")
	sb.WriteString(fmt.Sprintf("- **Security:** %d\n", report.Summary.SecurityTODOs))
	sb.WriteString(fmt.Sprintf("- **Performance:** %d\n", report.Summary.PerformanceTODOs))
	sb.WriteString(fmt.Sprintf("- **Features:** %d\n", report.Summary.FeatureTODOs))
	sb.WriteString(fmt.Sprintf("- **Bugs:** %d\n", report.Summary.BugTODOs))
	sb.WriteString("\n")

	// Detailed breakdown by priority
	priorities := []Priority{PriorityCritical, PriorityHigh, PriorityMedium, PriorityLow}
	priorityNames := map[Priority]string{
		PriorityCritical: "Critical",
		PriorityHigh:     "High",
		PriorityMedium:   "Medium",
		PriorityLow:      "Low",
	}

	for _, priority := range priorities {
		todos := report.PriorityBreakdown[priority]
		if len(todos) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("## %s Priority TODOs (%d)\n\n", priorityNames[priority], len(todos)))

		for _, todo := range todos {
			sb.WriteString(fmt.Sprintf("### %s:%d - %s\n", todo.File, todo.Line, todo.Type))
			sb.WriteString(fmt.Sprintf("**Message:** %s\n\n", todo.Message))
			sb.WriteString(fmt.Sprintf("**Context:** `%s`\n\n", todo.Context))
			sb.WriteString(fmt.Sprintf("**Category:** %s\n\n", ts.categoryToString(todo.Category)))
			sb.WriteString("---\n\n")
		}
	}

	// Files with most TODOs
	sb.WriteString("## Files with Most TODOs\n\n")
	type fileCount struct {
		file  string
		count int
	}

	var fileCounts []fileCount
	for file, todos := range report.FileBreakdown {
		fileCounts = append(fileCounts, fileCount{file: file, count: len(todos)})
	}

	sort.Slice(fileCounts, func(i, j int) bool {
		return fileCounts[i].count > fileCounts[j].count
	})

	for i, fc := range fileCounts {
		if i >= 10 { // Show top 10
			break
		}
		sb.WriteString(fmt.Sprintf("- **%s:** %d TODOs\n", fc.file, fc.count))
	}

	return sb.String(), nil
}

// generateTextReport generates a plain text report
func (ts *TODOScanner) generateTextReport(report *TODOReport) (string, error) {
	var sb strings.Builder

	sb.WriteString("TODO/FIXME Analysis Report\n")
	sb.WriteString("==========================\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("Project Root: %s\n", report.ProjectRoot))
	sb.WriteString(fmt.Sprintf("Files Scanned: %d/%d\n\n", report.FilesScanned, report.TotalFiles))

	sb.WriteString("SUMMARY\n")
	sb.WriteString("-------\n")
	sb.WriteString(fmt.Sprintf("Total TODOs: %d\n", report.Summary.TotalTODOs))
	sb.WriteString(fmt.Sprintf("Critical: %d, High: %d, Medium: %d, Low: %d\n\n",
		report.Summary.CriticalTODOs, report.Summary.HighTODOs,
		report.Summary.MediumTODOs, report.Summary.LowTODOs))

	for _, todo := range report.TODOs {
		sb.WriteString(fmt.Sprintf("[%s] %s:%d - %s\n",
			ts.priorityToString(todo.Priority), todo.File, todo.Line, todo.Type))
		sb.WriteString(fmt.Sprintf("    %s\n", todo.Message))
		sb.WriteString(fmt.Sprintf("    Context: %s\n\n", todo.Context))
	}

	return sb.String(), nil
}

// generateJSONReport generates a JSON report
func (ts *TODOScanner) generateJSONReport(report *TODOReport) (string, error) {
	// For simplicity, we'll create a basic JSON structure
	// In a real implementation, you'd use json.Marshal
	return fmt.Sprintf(`{
  "timestamp": "%s",
  "project_root": "%s",
  "total_files": %d,
  "files_scanned": %d,
  "summary": {
    "total_todos": %d,
    "critical": %d,
    "high": %d,
    "medium": %d,
    "low": %d,
    "security": %d,
    "performance": %d,
    "features": %d,
    "bugs": %d
  },
  "todo_count": %d
}`,
		report.Timestamp.Format(time.RFC3339),
		report.ProjectRoot,
		report.TotalFiles,
		report.FilesScanned,
		report.Summary.TotalTODOs,
		report.Summary.CriticalTODOs,
		report.Summary.HighTODOs,
		report.Summary.MediumTODOs,
		report.Summary.LowTODOs,
		report.Summary.SecurityTODOs,
		report.Summary.PerformanceTODOs,
		report.Summary.FeatureTODOs,
		report.Summary.BugTODOs,
		len(report.TODOs)), nil
}

// Helper methods for string conversion
func (ts *TODOScanner) priorityToString(p Priority) string {
	switch p {
	case PriorityCritical:
		return "CRITICAL"
	case PriorityHigh:
		return "HIGH"
	case PriorityMedium:
		return "MEDIUM"
	case PriorityLow:
		return "LOW"
	default:
		return "UNKNOWN"
	}
}

func (ts *TODOScanner) categoryToString(c Category) string {
	switch c {
	case CategorySecurity:
		return "Security"
	case CategoryPerformance:
		return "Performance"
	case CategoryFeature:
		return "Feature"
	case CategoryBug:
		return "Bug"
	case CategoryDocumentation:
		return "Documentation"
	case CategoryRefactor:
		return "Refactor"
	default:
		return "Unknown"
	}
}
