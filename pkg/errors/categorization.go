// Package errors provides error categorization and severity management
package errors

import (
	"fmt"
	"strings"
	"time"
)

// ErrorCategory represents a category of errors for better organization
type ErrorCategory struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
	Severity    Severity `json:"default_severity"`
	Recoverable bool     `json:"default_recoverable"`
}

// ErrorCategorizer manages error categorization and analysis
type ErrorCategorizer struct {
	categories map[string]*ErrorCategory
	stats      *ErrorStatistics
}

// ErrorStatistics tracks error statistics for analysis
type ErrorStatistics struct {
	TotalErrors      int            `json:"total_errors"`
	ErrorsByType     map[string]int `json:"errors_by_type"`
	ErrorsByCategory map[string]int `json:"errors_by_category"`
	ErrorsBySeverity map[string]int `json:"errors_by_severity"`
	ErrorsByHour     map[string]int `json:"errors_by_hour"`
	RecentErrors     []*CLIError    `json:"recent_errors"`
	FirstError       *time.Time     `json:"first_error,omitempty"`
	LastError        *time.Time     `json:"last_error,omitempty"`
	RecoveryRate     float64        `json:"recovery_rate"`
	CommonPatterns   []ErrorPattern `json:"common_patterns"`
}

// ErrorPattern represents a common error pattern
type ErrorPattern struct {
	Pattern     string   `json:"pattern"`
	Count       int      `json:"count"`
	Percentage  float64  `json:"percentage"`
	Suggestions []string `json:"suggestions"`
}

// SeverityLevel provides detailed severity information
type SeverityLevel struct {
	Level       Severity `json:"level"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Icon        string   `json:"icon"`
	Priority    int      `json:"priority"`
}

// NewErrorCategorizer creates a new error categorizer with default categories
func NewErrorCategorizer() *ErrorCategorizer {
	ec := &ErrorCategorizer{
		categories: make(map[string]*ErrorCategory),
		stats: &ErrorStatistics{
			ErrorsByType:     make(map[string]int),
			ErrorsByCategory: make(map[string]int),
			ErrorsBySeverity: make(map[string]int),
			ErrorsByHour:     make(map[string]int),
			RecentErrors:     make([]*CLIError, 0),
			CommonPatterns:   make([]ErrorPattern, 0),
		},
	}

	// Register default categories
	ec.registerDefaultCategories()

	return ec
}

// registerDefaultCategories registers the default error categories
func (ec *ErrorCategorizer) registerDefaultCategories() {
	categories := []*ErrorCategory{
		{
			Name:        "User Input",
			Description: "Errors related to user input and interaction",
			Types:       []string{ErrorTypeUser, ErrorTypeValidation},
			Severity:    SeverityLow,
			Recoverable: true,
		},
		{
			Name:        "Configuration",
			Description: "Errors related to configuration files and settings",
			Types:       []string{ErrorTypeConfiguration},
			Severity:    SeverityMedium,
			Recoverable: true,
		},
		{
			Name:        "Templates",
			Description: "Errors related to template processing and management",
			Types:       []string{ErrorTypeTemplate},
			Severity:    SeverityMedium,
			Recoverable: true,
		},
		{
			Name:        "File System",
			Description: "Errors related to file and directory operations",
			Types:       []string{ErrorTypeFileSystem, ErrorTypePermission},
			Severity:    SeverityHigh,
			Recoverable: false,
		},
		{
			Name:        "Network",
			Description: "Errors related to network operations and connectivity",
			Types:       []string{ErrorTypeNetwork},
			Severity:    SeverityMedium,
			Recoverable: true,
		},
		{
			Name:        "System",
			Description: "Errors related to system resources and operations",
			Types:       []string{ErrorTypeCache, ErrorTypeVersion},
			Severity:    SeverityMedium,
			Recoverable: true,
		},
		{
			Name:        "Security",
			Description: "Errors related to security and vulnerabilities",
			Types:       []string{ErrorTypeSecurity},
			Severity:    SeverityCritical,
			Recoverable: false,
		},
		{
			Name:        "Generation",
			Description: "Errors related to project generation and processing",
			Types:       []string{ErrorTypeGeneration, ErrorTypeAudit},
			Severity:    SeverityHigh,
			Recoverable: true,
		},
		{
			Name:        "Dependencies",
			Description: "Errors related to dependency management and compatibility",
			Types:       []string{ErrorTypeDependency},
			Severity:    SeverityMedium,
			Recoverable: true,
		},
		{
			Name:        "Internal",
			Description: "Internal system errors and unexpected conditions",
			Types:       []string{ErrorTypeInternal},
			Severity:    SeverityCritical,
			Recoverable: false,
		},
	}

	for _, category := range categories {
		ec.RegisterCategory(category)
	}
}

// RegisterCategory registers a new error category
func (ec *ErrorCategorizer) RegisterCategory(category *ErrorCategory) {
	ec.categories[category.Name] = category
}

// CategorizeError categorizes an error and returns the category
func (ec *ErrorCategorizer) CategorizeError(err *CLIError) *ErrorCategory {
	for _, category := range ec.categories {
		for _, errorType := range category.Types {
			if err.Type == errorType {
				return category
			}
		}
	}

	// Return default category if no match found
	return &ErrorCategory{
		Name:        "Unknown",
		Description: "Uncategorized error",
		Types:       []string{err.Type},
		Severity:    SeverityMedium,
		Recoverable: false,
	}
}

// RecordError records an error for statistical analysis
func (ec *ErrorCategorizer) RecordError(err *CLIError) {
	if err == nil {
		return
	}

	// Update basic statistics
	ec.stats.TotalErrors++
	ec.stats.ErrorsByType[err.Type]++
	ec.stats.ErrorsBySeverity[string(err.Severity)]++

	// Update category statistics
	category := ec.CategorizeError(err)
	ec.stats.ErrorsByCategory[category.Name]++

	// Update hourly statistics
	hour := err.Timestamp.Format("2006-01-02 15")
	ec.stats.ErrorsByHour[hour]++

	// Update time tracking
	if ec.stats.FirstError == nil {
		ec.stats.FirstError = &err.Timestamp
	}
	ec.stats.LastError = &err.Timestamp

	// Add to recent errors (keep last 100)
	ec.stats.RecentErrors = append(ec.stats.RecentErrors, err)
	if len(ec.stats.RecentErrors) > 100 {
		ec.stats.RecentErrors = ec.stats.RecentErrors[1:]
	}

	// Update patterns
	ec.updateErrorPatterns(err)
}

// updateErrorPatterns updates common error patterns
func (ec *ErrorCategorizer) updateErrorPatterns(err *CLIError) {
	// Extract pattern from error message
	pattern := ec.extractPattern(err.Message)

	// Find existing pattern or create new one
	var found *ErrorPattern
	for i := range ec.stats.CommonPatterns {
		if ec.stats.CommonPatterns[i].Pattern == pattern {
			found = &ec.stats.CommonPatterns[i]
			break
		}
	}

	if found != nil {
		found.Count++
	} else {
		newPattern := ErrorPattern{
			Pattern:     pattern,
			Count:       1,
			Suggestions: ec.generatePatternSuggestions(pattern),
		}
		ec.stats.CommonPatterns = append(ec.stats.CommonPatterns, newPattern)
	}

	// Update percentages
	ec.updatePatternPercentages()
}

// extractPattern extracts a pattern from an error message
func (ec *ErrorCategorizer) extractPattern(message string) string {
	// Normalize the message by removing specific details
	pattern := message

	// Replace file paths with placeholder
	pattern = strings.ReplaceAll(pattern, `\`, "/")
	if strings.Contains(pattern, "/") {
		parts := strings.Split(pattern, " ")
		for i, part := range parts {
			if strings.Contains(part, "/") {
				parts[i] = "<path>"
			}
		}
		pattern = strings.Join(parts, " ")
	}

	// Replace numbers with placeholder
	words := strings.Fields(pattern)
	for i, word := range words {
		if isNumeric(word) {
			words[i] = "<number>"
		}
	}
	pattern = strings.Join(words, " ")

	// Replace quoted strings with placeholder
	if strings.Contains(pattern, `"`) {
		pattern = strings.ReplaceAll(pattern, `"`, `"<string>"`)
	}

	return pattern
}

// isNumeric checks if a string is numeric
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return len(s) > 0
}

// generatePatternSuggestions generates suggestions for a pattern
func (ec *ErrorCategorizer) generatePatternSuggestions(pattern string) []string {
	var suggestions []string

	// Pattern-based suggestions
	switch {
	case strings.Contains(pattern, "not found"):
		suggestions = append(suggestions, []string{
			"Check if the resource exists",
			"Verify spelling and case sensitivity",
			"Ensure required dependencies are installed",
		}...)
	case strings.Contains(pattern, "permission denied"):
		suggestions = append(suggestions, []string{
			"Check file/directory permissions",
			"Run with appropriate user privileges",
			"Verify ownership of target files",
		}...)
	case strings.Contains(pattern, "connection"):
		suggestions = append(suggestions, []string{
			"Check network connectivity",
			"Verify service availability",
			"Check firewall settings",
		}...)
	case strings.Contains(pattern, "invalid"):
		suggestions = append(suggestions, []string{
			"Check input format and syntax",
			"Verify configuration values",
			"Consult documentation for valid options",
		}...)
	case strings.Contains(pattern, "timeout"):
		suggestions = append(suggestions, []string{
			"Increase timeout settings",
			"Check network latency",
			"Retry the operation",
		}...)
	default:
		suggestions = append(suggestions, []string{
			"Check error details for specific guidance",
			"Consult documentation for troubleshooting",
			"Use --verbose flag for more information",
		}...)
	}

	return suggestions
}

// updatePatternPercentages updates the percentage for each pattern
func (ec *ErrorCategorizer) updatePatternPercentages() {
	total := ec.stats.TotalErrors
	if total == 0 {
		return
	}

	for i := range ec.stats.CommonPatterns {
		ec.stats.CommonPatterns[i].Percentage = float64(ec.stats.CommonPatterns[i].Count) / float64(total) * 100
	}
}

// GetStatistics returns current error statistics
func (ec *ErrorCategorizer) GetStatistics() *ErrorStatistics {
	// Calculate recovery rate
	if ec.stats.TotalErrors > 0 {
		recoverable := 0
		for _, err := range ec.stats.RecentErrors {
			if err.Recoverable {
				recoverable++
			}
		}
		ec.stats.RecoveryRate = float64(recoverable) / float64(len(ec.stats.RecentErrors)) * 100
	}

	return ec.stats
}

// GetCategories returns all registered categories
func (ec *ErrorCategorizer) GetCategories() map[string]*ErrorCategory {
	return ec.categories
}

// GetSeverityLevels returns detailed information about severity levels
func (ec *ErrorCategorizer) GetSeverityLevels() []SeverityLevel {
	return []SeverityLevel{
		{
			Level:       SeverityLow,
			Name:        "Low",
			Description: "Minor issues that don't prevent operation",
			Color:       "#28a745", // Green
			Icon:        "ℹ️",
			Priority:    1,
		},
		{
			Level:       SeverityMedium,
			Name:        "Medium",
			Description: "Issues that may affect functionality",
			Color:       "#ffc107", // Yellow
			Icon:        "⚠️",
			Priority:    2,
		},
		{
			Level:       SeverityHigh,
			Name:        "High",
			Description: "Serious issues that prevent normal operation",
			Color:       "#fd7e14", // Orange
			Icon:        "🚨",
			Priority:    3,
		},
		{
			Level:       SeverityCritical,
			Name:        "Critical",
			Description: "Critical issues requiring immediate attention",
			Color:       "#dc3545", // Red
			Icon:        "🔥",
			Priority:    4,
		},
	}
}

// GenerateErrorReport generates a comprehensive error analysis report
func (ec *ErrorCategorizer) GenerateErrorReport() *ErrorAnalysisReport {
	stats := ec.GetStatistics()

	report := &ErrorAnalysisReport{
		GeneratedAt:       time.Now(),
		TotalErrors:       stats.TotalErrors,
		RecoveryRate:      stats.RecoveryRate,
		Categories:        make([]CategoryAnalysis, 0),
		SeverityBreakdown: make([]SeverityAnalysis, 0),
		TopPatterns:       ec.getTopPatterns(5),
		Recommendations:   ec.generateRecommendations(),
		TimeRange:         ec.getTimeRange(),
	}

	// Category analysis
	for name, count := range stats.ErrorsByCategory {
		if category, exists := ec.categories[name]; exists {
			analysis := CategoryAnalysis{
				Category:   category,
				Count:      count,
				Percentage: float64(count) / float64(stats.TotalErrors) * 100,
				Trend:      ec.calculateCategoryTrend(name),
			}
			report.Categories = append(report.Categories, analysis)
		}
	}

	// Severity analysis
	severityLevels := ec.GetSeverityLevels()
	for _, level := range severityLevels {
		count := stats.ErrorsBySeverity[string(level.Level)]
		analysis := SeverityAnalysis{
			Severity:   level,
			Count:      count,
			Percentage: float64(count) / float64(stats.TotalErrors) * 100,
		}
		report.SeverityBreakdown = append(report.SeverityBreakdown, analysis)
	}

	return report
}

// ErrorAnalysisReport represents a comprehensive error analysis report
type ErrorAnalysisReport struct {
	GeneratedAt       time.Time          `json:"generated_at"`
	TotalErrors       int                `json:"total_errors"`
	RecoveryRate      float64            `json:"recovery_rate"`
	Categories        []CategoryAnalysis `json:"categories"`
	SeverityBreakdown []SeverityAnalysis `json:"severity_breakdown"`
	TopPatterns       []ErrorPattern     `json:"top_patterns"`
	Recommendations   []string           `json:"recommendations"`
	TimeRange         *TimeRange         `json:"time_range"`
}

// CategoryAnalysis represents analysis for a specific category
type CategoryAnalysis struct {
	Category   *ErrorCategory `json:"category"`
	Count      int            `json:"count"`
	Percentage float64        `json:"percentage"`
	Trend      string         `json:"trend"` // "increasing", "decreasing", "stable"
}

// SeverityAnalysis represents analysis for a specific severity level
type SeverityAnalysis struct {
	Severity   SeverityLevel `json:"severity"`
	Count      int           `json:"count"`
	Percentage float64       `json:"percentage"`
}

// TimeRange represents the time range of the analysis
type TimeRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

// getTopPatterns returns the top N error patterns
func (ec *ErrorCategorizer) getTopPatterns(n int) []ErrorPattern {
	patterns := make([]ErrorPattern, len(ec.stats.CommonPatterns))
	copy(patterns, ec.stats.CommonPatterns)

	// Sort by count (descending)
	for i := 0; i < len(patterns)-1; i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[i].Count < patterns[j].Count {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}

	// Return top N
	if len(patterns) > n {
		patterns = patterns[:n]
	}

	return patterns
}

// calculateCategoryTrend calculates the trend for a category
func (ec *ErrorCategorizer) calculateCategoryTrend(categoryName string) string {
	// This would implement trend calculation based on historical data
	// For now, return "stable" as a placeholder
	return "stable"
}

// generateRecommendations generates recommendations based on error patterns
func (ec *ErrorCategorizer) generateRecommendations() []string {
	var recommendations []string

	stats := ec.stats

	// High error rate recommendations
	if stats.TotalErrors > 100 {
		recommendations = append(recommendations, "Consider reviewing system configuration due to high error rate")
	}

	// Low recovery rate recommendations
	if stats.RecoveryRate < 50 {
		recommendations = append(recommendations, "Improve error recovery mechanisms to increase recovery rate")
	}

	// Category-specific recommendations
	for category, count := range stats.ErrorsByCategory {
		percentage := float64(count) / float64(stats.TotalErrors) * 100

		if percentage > 30 {
			switch category {
			case "Configuration":
				recommendations = append(recommendations, "High configuration errors - review configuration validation and documentation")
			case "Network":
				recommendations = append(recommendations, "High network errors - implement better offline mode and retry mechanisms")
			case "File System":
				recommendations = append(recommendations, "High file system errors - review file permission handling and error messages")
			case "User Input":
				recommendations = append(recommendations, "High user input errors - improve input validation and user guidance")
			}
		}
	}

	// Pattern-specific recommendations
	for _, pattern := range ec.getTopPatterns(3) {
		if pattern.Percentage > 20 {
			recommendations = append(recommendations, fmt.Sprintf("Address common pattern: %s (%.1f%% of errors)", pattern.Pattern, pattern.Percentage))
		}
	}

	return recommendations
}

// getTimeRange returns the time range of recorded errors
func (ec *ErrorCategorizer) getTimeRange() *TimeRange {
	return &TimeRange{
		Start: ec.stats.FirstError,
		End:   ec.stats.LastError,
	}
}

// Reset resets all error statistics
func (ec *ErrorCategorizer) Reset() {
	ec.stats = &ErrorStatistics{
		ErrorsByType:     make(map[string]int),
		ErrorsByCategory: make(map[string]int),
		ErrorsBySeverity: make(map[string]int),
		ErrorsByHour:     make(map[string]int),
		RecentErrors:     make([]*CLIError, 0),
		CommonPatterns:   make([]ErrorPattern, 0),
	}
}
