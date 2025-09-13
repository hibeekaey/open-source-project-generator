package cleanup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CleanupUtils provides additional utility functions for cleanup operations
type CleanupUtils struct{}

// NewCleanupUtils creates a new cleanup utilities instance
func NewCleanupUtils() *CleanupUtils {
	return &CleanupUtils{}
}

// ScanProjectFiles recursively scans for Go files in a project
func (cu *CleanupUtils) ScanProjectFiles(rootDir string, skipPatterns []string) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Check skip patterns
		for _, pattern := range skipPatterns {
			if strings.Contains(path, pattern) {
				return nil
			}
		}

		goFiles = append(goFiles, path)
		return nil
	})

	return goFiles, err
}

// GenerateCleanupReport creates a comprehensive cleanup report
func (cu *CleanupUtils) GenerateCleanupReport(analysis *ProjectAnalysis, result *CleanupResult) string {
	var report strings.Builder

	report.WriteString("# Cleanup Report\n\n")
	report.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339)))
	report.WriteString(fmt.Sprintf("Project: %s\n\n", analysis.ProjectRoot))

	// Summary section
	report.WriteString("## Summary\n\n")
	if result != nil {
		report.WriteString(fmt.Sprintf("- Files Modified: %d\n", len(result.FilesModified)))
		report.WriteString(fmt.Sprintf("- Issues Fixed: %d\n", len(result.IssuesFixed)))
		report.WriteString(fmt.Sprintf("- Issues Remaining: %d\n", len(result.IssuesRemaining)))
		report.WriteString(fmt.Sprintf("- Duration: %v\n", result.Duration))
		report.WriteString(fmt.Sprintf("- Success: %t\n\n", result.Success))
	}

	// Analysis results
	report.WriteString("## Analysis Results\n\n")
	report.WriteString(analysis.GetSummary())
	report.WriteString("\n")

	// TODO/FIXME Details
	if len(analysis.TODOs) > 0 {
		report.WriteString("### TODO/FIXME Comments\n\n")

		// Group by category
		categories := make(map[Category][]TODOItem)
		for _, todo := range analysis.TODOs {
			categories[todo.Category] = append(categories[todo.Category], todo)
		}

		for category, todos := range categories {
			categoryName := cu.getCategoryName(category)
			report.WriteString(fmt.Sprintf("#### %s (%d items)\n\n", categoryName, len(todos)))

			for _, todo := range todos {
				priority := cu.getPriorityName(todo.Priority)
				report.WriteString(fmt.Sprintf("- **%s** [%s] `%s:%d` - %s\n",
					todo.Type, priority, todo.File, todo.Line, todo.Message))
			}
			report.WriteString("\n")
		}
	}

	// Duplicate Code Details
	if len(analysis.Duplicates) > 0 {
		report.WriteString("### Duplicate Code Blocks\n\n")
		for i, dup := range analysis.Duplicates {
			report.WriteString(fmt.Sprintf("#### Block %d (Similarity: %.1f%%)\n\n",
				i+1, dup.Similarity*100))
			report.WriteString("Files:\n")
			for _, file := range dup.Files {
				report.WriteString(fmt.Sprintf("- %s\n", file))
			}
			report.WriteString(fmt.Sprintf("\nSuggestion: %s\n\n", dup.Suggestion))
		}
	}

	// Unused Code Details
	if len(analysis.UnusedCode) > 0 {
		report.WriteString("### Unused Code Items\n\n")

		// Group by type
		types := make(map[string][]UnusedCodeItem)
		for _, unused := range analysis.UnusedCode {
			types[unused.Type] = append(types[unused.Type], unused)
		}

		for itemType, items := range types {
			report.WriteString(fmt.Sprintf("#### %s (%d items)\n\n",
				strings.Title(itemType), len(items)))

			for _, item := range items {
				report.WriteString(fmt.Sprintf("- `%s:%d` - %s (%s)\n",
					item.File, item.Line, item.Name, item.Reason))
			}
			report.WriteString("\n")
		}
	}

	// Import Issues
	if len(analysis.ImportIssues) > 0 {
		report.WriteString("### Import Organization Issues\n\n")
		for _, issue := range analysis.ImportIssues {
			report.WriteString(fmt.Sprintf("- `%s:%d` - %s: %s\n",
				issue.File, issue.Line, issue.Import, issue.Suggestion))
		}
		report.WriteString("\n")
	}

	// Recommendations
	report.WriteString("## Recommendations\n\n")
	recommendations := cu.generateRecommendations(analysis)
	for _, rec := range recommendations {
		report.WriteString(fmt.Sprintf("- %s\n", rec))
	}

	return report.String()
}

// ValidateCleanupResult validates the results of a cleanup operation
func (cu *CleanupUtils) ValidateCleanupResult(result *CleanupResult) []string {
	var issues []string

	if result == nil {
		return []string{"Cleanup result is nil"}
	}

	// Check for basic consistency
	if !result.Success && len(result.IssuesFixed) > 0 {
		issues = append(issues, "Result marked as failed but shows fixed issues")
	}

	if result.ValidationResult != nil && !result.ValidationResult.Success && result.Success {
		issues = append(issues, "Cleanup marked as successful but validation failed")
	}

	// Check for reasonable duration
	if result.Duration < 0 {
		issues = append(issues, "Invalid negative duration")
	}

	if result.Duration > 24*time.Hour {
		issues = append(issues, "Unusually long cleanup duration (>24 hours)")
	}

	// Check file modification consistency
	if len(result.FilesModified) == 0 && len(result.IssuesFixed) > 0 {
		issues = append(issues, "Issues marked as fixed but no files were modified")
	}

	return issues
}

// EstimateCleanupTime estimates how long cleanup will take based on analysis
func (cu *CleanupUtils) EstimateCleanupTime(analysis *ProjectAnalysis) time.Duration {
	if analysis == nil {
		return 0
	}

	// Base time for setup and validation
	baseTime := 30 * time.Second

	// Time per TODO (varies by priority)
	todoTime := time.Duration(0)
	for _, todo := range analysis.TODOs {
		switch todo.Priority {
		case PriorityCritical:
			todoTime += 5 * time.Minute
		case PriorityHigh:
			todoTime += 2 * time.Minute
		case PriorityMedium:
			todoTime += 1 * time.Minute
		case PriorityLow:
			todoTime += 30 * time.Second
		}
	}

	// Time for duplicate code resolution
	duplicateTime := time.Duration(len(analysis.Duplicates)) * 3 * time.Minute

	// Time for unused code removal
	unusedTime := time.Duration(len(analysis.UnusedCode)) * 10 * time.Second

	// Time for import organization
	importTime := time.Duration(len(analysis.ImportIssues)) * 5 * time.Second

	total := baseTime + todoTime + duplicateTime + unusedTime + importTime

	// Add buffer for validation and testing
	return total + (total / 4) // 25% buffer
}

// CreateCleanupPlan creates a structured plan for cleanup operations
func (cu *CleanupUtils) CreateCleanupPlan(analysis *ProjectAnalysis) *CleanupPlan {
	plan := &CleanupPlan{
		ProjectRoot: analysis.ProjectRoot,
		CreatedAt:   time.Now(),
		Phases:      []CleanupPhase{},
	}

	// Phase 1: High-priority security issues
	securityTodos := cu.filterTodosByCategory(analysis.TODOs, CategorySecurity)
	if len(securityTodos) > 0 {
		plan.Phases = append(plan.Phases, CleanupPhase{
			Name:          "Security Issues",
			Description:   "Address security-related TODOs and vulnerabilities",
			Priority:      1,
			EstimatedTime: time.Duration(len(securityTodos)) * 3 * time.Minute,
			Tasks:         cu.createTasksFromTodos(securityTodos),
		})
	}

	// Phase 2: Remove unused code
	if len(analysis.UnusedCode) > 0 {
		plan.Phases = append(plan.Phases, CleanupPhase{
			Name:          "Unused Code Removal",
			Description:   "Remove unused functions, variables, and imports",
			Priority:      2,
			EstimatedTime: time.Duration(len(analysis.UnusedCode)) * 10 * time.Second,
			Tasks:         cu.createTasksFromUnused(analysis.UnusedCode),
		})
	}

	// Phase 3: Organize imports
	if len(analysis.ImportIssues) > 0 {
		plan.Phases = append(plan.Phases, CleanupPhase{
			Name:          "Import Organization",
			Description:   "Organize and clean up import statements",
			Priority:      3,
			EstimatedTime: time.Duration(len(analysis.ImportIssues)) * 5 * time.Second,
			Tasks:         cu.createTasksFromImports(analysis.ImportIssues),
		})
	}

	// Phase 4: Address remaining TODOs
	remainingTodos := cu.filterTodosExcludeCategory(analysis.TODOs, CategorySecurity)
	if len(remainingTodos) > 0 {
		plan.Phases = append(plan.Phases, CleanupPhase{
			Name:          "Remaining TODOs",
			Description:   "Address remaining TODO/FIXME comments",
			Priority:      4,
			EstimatedTime: time.Duration(len(remainingTodos)) * 1 * time.Minute,
			Tasks:         cu.createTasksFromTodos(remainingTodos),
		})
	}

	// Phase 5: Duplicate code consolidation
	if len(analysis.Duplicates) > 0 {
		plan.Phases = append(plan.Phases, CleanupPhase{
			Name:          "Duplicate Code",
			Description:   "Consolidate duplicate code blocks",
			Priority:      5,
			EstimatedTime: time.Duration(len(analysis.Duplicates)) * 3 * time.Minute,
			Tasks:         cu.createTasksFromDuplicates(analysis.Duplicates),
		})
	}

	return plan
}

// Helper methods

func (cu *CleanupUtils) getCategoryName(category Category) string {
	switch category {
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
		return "Other"
	}
}

func (cu *CleanupUtils) getPriorityName(priority Priority) string {
	switch priority {
	case PriorityCritical:
		return "Critical"
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

func (cu *CleanupUtils) generateRecommendations(analysis *ProjectAnalysis) []string {
	var recommendations []string

	// Security recommendations
	securityCount := len(cu.filterTodosByCategory(analysis.TODOs, CategorySecurity))
	if securityCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Address %d security-related TODOs immediately", securityCount))
	}

	// Performance recommendations
	perfCount := len(cu.filterTodosByCategory(analysis.TODOs, CategoryPerformance))
	if perfCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Review %d performance-related TODOs for optimization opportunities", perfCount))
	}

	// Code quality recommendations
	if len(analysis.UnusedCode) > 10 {
		recommendations = append(recommendations,
			"Consider running automated unused code removal")
	}

	if len(analysis.Duplicates) > 5 {
		recommendations = append(recommendations,
			"Significant code duplication detected - consider refactoring")
	}

	if len(analysis.ImportIssues) > 0 {
		recommendations = append(recommendations,
			"Run goimports to automatically fix import organization")
	}

	// General recommendations
	totalIssues := len(analysis.TODOs) + len(analysis.UnusedCode) + len(analysis.Duplicates) + len(analysis.ImportIssues)
	if totalIssues > 50 {
		recommendations = append(recommendations,
			"Consider implementing automated code quality checks in CI/CD pipeline")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Code quality looks good! Consider regular cleanup maintenance.")
	}

	return recommendations
}

func (cu *CleanupUtils) filterTodosByCategory(todos []TODOItem, category Category) []TODOItem {
	var filtered []TODOItem
	for _, todo := range todos {
		if todo.Category == category {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func (cu *CleanupUtils) filterTodosExcludeCategory(todos []TODOItem, category Category) []TODOItem {
	var filtered []TODOItem
	for _, todo := range todos {
		if todo.Category != category {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func (cu *CleanupUtils) createTasksFromTodos(todos []TODOItem) []CleanupTask {
	var tasks []CleanupTask
	for _, todo := range todos {
		tasks = append(tasks, CleanupTask{
			Type:        "todo",
			Description: fmt.Sprintf("Address %s: %s", todo.Type, todo.Message),
			File:        todo.File,
			Line:        todo.Line,
			Priority:    todo.Priority,
		})
	}
	return tasks
}

func (cu *CleanupUtils) createTasksFromUnused(unused []UnusedCodeItem) []CleanupTask {
	var tasks []CleanupTask
	for _, item := range unused {
		tasks = append(tasks, CleanupTask{
			Type:        "unused",
			Description: fmt.Sprintf("Remove unused %s: %s", item.Type, item.Name),
			File:        item.File,
			Line:        item.Line,
			Priority:    PriorityMedium,
		})
	}
	return tasks
}

func (cu *CleanupUtils) createTasksFromImports(issues []ImportIssue) []CleanupTask {
	var tasks []CleanupTask
	for _, issue := range issues {
		tasks = append(tasks, CleanupTask{
			Type:        "import",
			Description: fmt.Sprintf("Fix import organization: %s", issue.Suggestion),
			File:        issue.File,
			Line:        issue.Line,
			Priority:    PriorityLow,
		})
	}
	return tasks
}

func (cu *CleanupUtils) createTasksFromDuplicates(duplicates []DuplicateCodeBlock) []CleanupTask {
	var tasks []CleanupTask
	for i, dup := range duplicates {
		tasks = append(tasks, CleanupTask{
			Type:        "duplicate",
			Description: fmt.Sprintf("Consolidate duplicate code block %d: %s", i+1, dup.Suggestion),
			File:        strings.Join(dup.Files, ", "),
			Line:        0,
			Priority:    PriorityMedium,
		})
	}
	return tasks
}

// CleanupPlan represents a structured plan for cleanup operations
type CleanupPlan struct {
	ProjectRoot string
	CreatedAt   time.Time
	Phases      []CleanupPhase
}

// CleanupPhase represents a phase in the cleanup plan
type CleanupPhase struct {
	Name          string
	Description   string
	Priority      int
	EstimatedTime time.Duration
	Tasks         []CleanupTask
}

// CleanupTask represents an individual cleanup task
type CleanupTask struct {
	Type        string
	Description string
	File        string
	Line        int
	Priority    Priority
}

// GetTotalEstimatedTime returns the total estimated time for all phases
func (cp *CleanupPlan) GetTotalEstimatedTime() time.Duration {
	var total time.Duration
	for _, phase := range cp.Phases {
		total += phase.EstimatedTime
	}
	return total
}

// GetTaskCount returns the total number of tasks across all phases
func (cp *CleanupPlan) GetTaskCount() int {
	var count int
	for _, phase := range cp.Phases {
		count += len(phase.Tasks)
	}
	return count
}
