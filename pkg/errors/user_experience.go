// Package errors provides user experience enhancements for error handling
package errors

import (
	"fmt"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// UserExperienceManager enhances error handling with user-focused improvements
type UserExperienceManager struct {
	config           *EnhancedErrorConfig
	logger           interfaces.Logger
	helpProvider     *ContextualHelpProvider
	quickFixProvider *QuickFixProvider
	docProvider      *DocumentationProvider
	stats            *UserExperienceStatistics
	mutex            sync.RWMutex
}

// UserExperienceInfo contains user experience enhancements
type UserExperienceInfo struct {
	ContextualHelp    string             `json:"contextual_help,omitempty"`
	QuickFixes        []QuickFix         `json:"quick_fixes,omitempty"`
	RelatedDocs       []Documentation    `json:"related_docs,omitempty"`
	SimilarIssues     []SimilarIssue     `json:"similar_issues,omitempty"`
	LearningResources []LearningResource `json:"learning_resources,omitempty"`
	NextSteps         []string           `json:"next_steps,omitempty"`
}

// QuickFix represents an actionable fix for an error
type QuickFix struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Command       string  `json:"command,omitempty"`
	Script        string  `json:"script,omitempty"`
	Automated     bool    `json:"automated"`
	Confidence    float64 `json:"confidence"` // 0.0 to 1.0
	EstimatedTime string  `json:"estimated_time,omitempty"`
}

// Documentation represents related documentation
type Documentation struct {
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Type        string  `json:"type"`      // "guide", "reference", "tutorial", "troubleshooting"
	Relevance   float64 `json:"relevance"` // 0.0 to 1.0
	Description string  `json:"description,omitempty"`
}

// SimilarIssue represents a similar issue that was resolved
type SimilarIssue struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Solution    string    `json:"solution"`
	URL         string    `json:"url,omitempty"`
	Similarity  float64   `json:"similarity"` // 0.0 to 1.0
	ResolvedAt  time.Time `json:"resolved_at"`
}

// LearningResource represents educational content
type LearningResource struct {
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	Type      string  `json:"type"`  // "video", "article", "course", "example"
	Level     string  `json:"level"` // "beginner", "intermediate", "advanced"
	Duration  string  `json:"duration,omitempty"`
	Relevance float64 `json:"relevance"` // 0.0 to 1.0
}

// UserExperienceStatistics tracks user experience metrics
type UserExperienceStatistics struct {
	TotalEnhancements     int            `json:"total_enhancements"`
	QuickFixesProvided    int            `json:"quick_fixes_provided"`
	QuickFixesExecuted    int            `json:"quick_fixes_executed"`
	QuickFixSuccessRate   float64        `json:"quick_fix_success_rate"`
	HelpRequestsByType    map[string]int `json:"help_requests_by_type"`
	DocumentationViews    int            `json:"documentation_views"`
	AverageResolutionTime time.Duration  `json:"average_resolution_time"`
	UserSatisfactionScore float64        `json:"user_satisfaction_score"`
}

// NewUserExperienceManager creates a new user experience manager
func NewUserExperienceManager(config *EnhancedErrorConfig, logger interfaces.Logger) *UserExperienceManager {
	return &UserExperienceManager{
		config:           config,
		logger:           logger,
		helpProvider:     NewContextualHelpProvider(),
		quickFixProvider: NewQuickFixProvider(),
		docProvider:      NewDocumentationProvider(),
		stats: &UserExperienceStatistics{
			HelpRequestsByType: make(map[string]int),
		},
	}
}

// EnhanceUserExperience enhances error handling with user experience improvements
func (uxm *UserExperienceManager) EnhanceUserExperience(err *CLIError, diagnostics *DiagnosticInfo) *UserExperienceInfo {
	if err == nil {
		return nil
	}

	uxm.mutex.Lock()
	uxm.stats.TotalEnhancements++
	uxm.stats.HelpRequestsByType[err.Type]++
	uxm.mutex.Unlock()

	info := &UserExperienceInfo{}

	// Generate contextual help
	if help := uxm.helpProvider.GetContextualHelp(err, diagnostics); help != "" {
		info.ContextualHelp = help
	}

	// Generate quick fixes
	if fixes := uxm.quickFixProvider.GetQuickFixes(err, diagnostics); len(fixes) > 0 {
		info.QuickFixes = fixes
		uxm.mutex.Lock()
		uxm.stats.QuickFixesProvided += len(fixes)
		uxm.mutex.Unlock()
	}

	// Get related documentation
	if docs := uxm.docProvider.GetRelatedDocumentation(err, diagnostics); len(docs) > 0 {
		info.RelatedDocs = docs
	}

	// Find similar issues
	if similar := uxm.findSimilarIssues(err); len(similar) > 0 {
		info.SimilarIssues = similar
	}

	// Get learning resources
	if resources := uxm.getLearningResources(err); len(resources) > 0 {
		info.LearningResources = resources
	}

	// Generate next steps
	if steps := uxm.generateNextSteps(err, info); len(steps) > 0 {
		info.NextSteps = steps
	}

	return info
}

// findSimilarIssues finds similar issues that have been resolved
func (uxm *UserExperienceManager) findSimilarIssues(err *CLIError) []SimilarIssue {
	// This would typically query a knowledge base or issue tracker
	// For now, return some common similar issues based on error type

	var issues []SimilarIssue

	switch err.Type {
	case ErrorTypeValidation:
		issues = append(issues, SimilarIssue{
			Title:       "Project structure validation failed",
			Description: "Similar validation error with missing required files",
			Solution:    "Run 'generator validate --fix' to automatically fix common issues",
			Similarity:  0.8,
			ResolvedAt:  time.Now().Add(-24 * time.Hour),
		})
	case ErrorTypeConfiguration:
		issues = append(issues, SimilarIssue{
			Title:       "Configuration file syntax error",
			Description: "YAML syntax error in configuration file",
			Solution:    "Check YAML indentation and syntax using an online validator",
			Similarity:  0.7,
			ResolvedAt:  time.Now().Add(-48 * time.Hour),
		})
	case ErrorTypeNetwork:
		issues = append(issues, SimilarIssue{
			Title:       "Network connectivity issue",
			Description: "Unable to download templates due to network error",
			Solution:    "Use --offline flag to work with cached templates",
			Similarity:  0.9,
			ResolvedAt:  time.Now().Add(-12 * time.Hour),
		})
	}

	return issues
}

// getLearningResources gets relevant learning resources
func (uxm *UserExperienceManager) getLearningResources(err *CLIError) []LearningResource {
	var resources []LearningResource

	switch err.Type {
	case ErrorTypeValidation:
		resources = append(resources, LearningResource{
			Title:     "Project Structure Best Practices",
			URL:       "https://docs.generator.dev/guides/project-structure",
			Type:      "guide",
			Level:     "beginner",
			Duration:  "10 min",
			Relevance: 0.9,
		})
	case ErrorTypeConfiguration:
		resources = append(resources, LearningResource{
			Title:     "Configuration File Reference",
			URL:       "https://docs.generator.dev/reference/configuration",
			Type:      "reference",
			Level:     "intermediate",
			Relevance: 0.8,
		})
	case ErrorTypeTemplate:
		resources = append(resources, LearningResource{
			Title:     "Working with Templates",
			URL:       "https://docs.generator.dev/tutorials/templates",
			Type:      "tutorial",
			Level:     "beginner",
			Duration:  "15 min",
			Relevance: 0.9,
		})
	}

	return resources
}

// generateNextSteps generates actionable next steps
func (uxm *UserExperienceManager) generateNextSteps(err *CLIError, info *UserExperienceInfo) []string {
	var steps []string

	// Add quick fix steps
	if len(info.QuickFixes) > 0 {
		for i, fix := range info.QuickFixes {
			if i < 3 { // Limit to top 3 quick fixes
				if fix.Command != "" {
					steps = append(steps, fmt.Sprintf("Try: %s", fix.Command))
				} else {
					steps = append(steps, fix.Description)
				}
			}
		}
	}

	// Add type-specific steps
	switch err.Type {
	case ErrorTypeValidation:
		steps = append(steps, "Run validation with --verbose for detailed information")
		steps = append(steps, "Check project structure against template requirements")
	case ErrorTypeConfiguration:
		steps = append(steps, "Validate configuration syntax")
		steps = append(steps, "Review configuration documentation")
	case ErrorTypeNetwork:
		steps = append(steps, "Check internet connectivity")
		steps = append(steps, "Consider using offline mode")
	case ErrorTypeFileSystem:
		steps = append(steps, "Check file permissions")
		steps = append(steps, "Verify disk space availability")
	}

	// Add documentation step if available
	if len(info.RelatedDocs) > 0 {
		steps = append(steps, fmt.Sprintf("Read: %s", info.RelatedDocs[0].Title))
	}

	return steps
}

// SetInteractiveMode sets interactive mode for user experience
func (uxm *UserExperienceManager) SetInteractiveMode(interactive bool) {
	uxm.mutex.Lock()
	defer uxm.mutex.Unlock()
	uxm.config.InteractiveMode = interactive
}

// GetStatistics returns user experience statistics
func (uxm *UserExperienceManager) GetStatistics() *UserExperienceStatistics {
	uxm.mutex.RLock()
	defer uxm.mutex.RUnlock()

	// Calculate success rate
	if uxm.stats.QuickFixesProvided > 0 {
		uxm.stats.QuickFixSuccessRate = float64(uxm.stats.QuickFixesExecuted) / float64(uxm.stats.QuickFixesProvided)
	}

	// Return a copy
	stats := *uxm.stats
	return &stats
}

// ContextualHelpProvider provides contextual help for errors
type ContextualHelpProvider struct{}

// NewContextualHelpProvider creates a new contextual help provider
func NewContextualHelpProvider() *ContextualHelpProvider {
	return &ContextualHelpProvider{}
}

// GetContextualHelp generates contextual help for an error
func (chp *ContextualHelpProvider) GetContextualHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	if err == nil {
		return ""
	}

	// Generate help based on error type and context
	switch err.Type {
	case ErrorTypeValidation:
		return chp.getValidationHelp(err, diagnostics)
	case ErrorTypeConfiguration:
		return chp.getConfigurationHelp(err, diagnostics)
	case ErrorTypeTemplate:
		return chp.getTemplateHelp(err, diagnostics)
	case ErrorTypeNetwork:
		return chp.getNetworkHelp(err, diagnostics)
	case ErrorTypeFileSystem:
		return chp.getFileSystemHelp(err, diagnostics)
	default:
		return chp.getGenericHelp(err, diagnostics)
	}
}

// getValidationHelp provides help for validation errors
func (chp *ContextualHelpProvider) getValidationHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	help := "Validation ensures your project meets quality standards and follows best practices."

	if field, ok := err.Details["field"].(string); ok {
		help += fmt.Sprintf(" The issue is with the '%s' field.", field)
	}

	help += " Use --fix flag to automatically resolve common validation issues."
	return help
}

// getConfigurationHelp provides help for configuration errors
func (chp *ContextualHelpProvider) getConfigurationHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	help := "Configuration files define how your project should be generated."

	if configPath, ok := err.Details["config_path"].(string); ok {
		help += fmt.Sprintf(" Check the syntax and structure of '%s'.", configPath)
	}

	help += " Ensure proper YAML/JSON formatting and required fields are present."
	return help
}

// getTemplateHelp provides help for template errors
func (chp *ContextualHelpProvider) getTemplateHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	help := "Templates provide the structure and files for your project."

	if templateName, ok := err.Details["template_name"].(string); ok {
		help += fmt.Sprintf(" The template '%s' may not exist or be accessible.", templateName)
	}

	help += " Use 'generator list-templates' to see available options."
	return help
}

// getNetworkHelp provides help for network errors
func (chp *ContextualHelpProvider) getNetworkHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	help := "Network connectivity is required to download templates and check for updates."

	if url, ok := err.Details["url"].(string); ok {
		help += fmt.Sprintf(" Unable to reach '%s'.", url)
	}

	help += " Check your internet connection or use --offline mode with cached data."
	return help
}

// getFileSystemHelp provides help for filesystem errors
func (chp *ContextualHelpProvider) getFileSystemHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	help := "File system operations require proper permissions and available disk space."

	if path, ok := err.Details["path"].(string); ok {
		help += fmt.Sprintf(" Issue with path: '%s'.", path)
	}

	help += " Check file permissions, disk space, and that the path is accessible."
	return help
}

// getGenericHelp provides generic help for unknown error types
func (chp *ContextualHelpProvider) getGenericHelp(err *CLIError, diagnostics *DiagnosticInfo) string {
	return "An unexpected error occurred. Use --verbose for more details or --help for command usage information."
}

// QuickFixProvider provides automated fixes for common errors
type QuickFixProvider struct{}

// NewQuickFixProvider creates a new quick fix provider
func NewQuickFixProvider() *QuickFixProvider {
	return &QuickFixProvider{}
}

// GetQuickFixes generates quick fixes for an error
func (qfp *QuickFixProvider) GetQuickFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	if err == nil {
		return nil
	}

	var fixes []QuickFix

	switch err.Type {
	case ErrorTypeValidation:
		fixes = append(fixes, qfp.getValidationFixes(err, diagnostics)...)
	case ErrorTypeConfiguration:
		fixes = append(fixes, qfp.getConfigurationFixes(err, diagnostics)...)
	case ErrorTypeTemplate:
		fixes = append(fixes, qfp.getTemplateFixes(err, diagnostics)...)
	case ErrorTypeNetwork:
		fixes = append(fixes, qfp.getNetworkFixes(err, diagnostics)...)
	case ErrorTypeFileSystem:
		fixes = append(fixes, qfp.getFileSystemFixes(err, diagnostics)...)
	}

	return fixes
}

// getValidationFixes provides fixes for validation errors
func (qfp *QuickFixProvider) getValidationFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	return []QuickFix{
		{
			ID:            "auto-fix-validation",
			Title:         "Auto-fix validation issues",
			Description:   "Automatically fix common validation problems",
			Command:       "generator validate --fix",
			Automated:     true,
			Confidence:    0.8,
			EstimatedTime: "30 seconds",
		},
		{
			ID:            "verbose-validation",
			Title:         "Get detailed validation report",
			Description:   "Run validation with detailed output to see all issues",
			Command:       "generator validate --verbose",
			Automated:     false,
			Confidence:    1.0,
			EstimatedTime: "1 minute",
		},
	}
}

// getConfigurationFixes provides fixes for configuration errors
func (qfp *QuickFixProvider) getConfigurationFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	fixes := []QuickFix{
		{
			ID:            "validate-config",
			Title:         "Validate configuration",
			Description:   "Check configuration file syntax and structure",
			Command:       "generator config validate",
			Automated:     false,
			Confidence:    0.9,
			EstimatedTime: "30 seconds",
		},
	}

	if configPath, ok := err.Details["config_path"].(string); ok {
		fixes = append(fixes, QuickFix{
			ID:            "edit-config",
			Title:         "Edit configuration file",
			Description:   fmt.Sprintf("Open configuration file for editing: %s", configPath),
			Command:       fmt.Sprintf("$EDITOR %s", configPath),
			Automated:     false,
			Confidence:    0.7,
			EstimatedTime: "5 minutes",
		})
	}

	return fixes
}

// getTemplateFixes provides fixes for template errors
func (qfp *QuickFixProvider) getTemplateFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	return []QuickFix{
		{
			ID:            "list-templates",
			Title:         "List available templates",
			Description:   "Show all available templates",
			Command:       "generator list-templates",
			Automated:     false,
			Confidence:    1.0,
			EstimatedTime: "10 seconds",
		},
		{
			ID:            "search-templates",
			Title:         "Search for similar templates",
			Description:   "Find templates with similar names",
			Command:       "generator list-templates --search",
			Automated:     false,
			Confidence:    0.8,
			EstimatedTime: "15 seconds",
		},
	}
}

// getNetworkFixes provides fixes for network errors
func (qfp *QuickFixProvider) getNetworkFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	return []QuickFix{
		{
			ID:            "enable-offline",
			Title:         "Enable offline mode",
			Description:   "Use cached data instead of downloading",
			Command:       "generator --offline",
			Automated:     true,
			Confidence:    0.9,
			EstimatedTime: "immediate",
		},
		{
			ID:            "check-connectivity",
			Title:         "Test network connectivity",
			Description:   "Check if you can reach the internet",
			Command:       "ping -c 3 8.8.8.8",
			Automated:     false,
			Confidence:    0.7,
			EstimatedTime: "10 seconds",
		},
	}
}

// getFileSystemFixes provides fixes for filesystem errors
func (qfp *QuickFixProvider) getFileSystemFixes(err *CLIError, diagnostics *DiagnosticInfo) []QuickFix {
	fixes := []QuickFix{
		{
			ID:            "check-permissions",
			Title:         "Check file permissions",
			Description:   "Verify you have the necessary permissions",
			Command:       "ls -la",
			Automated:     false,
			Confidence:    0.8,
			EstimatedTime: "30 seconds",
		},
	}

	if path, ok := err.Details["path"].(string); ok {
		fixes = append(fixes, QuickFix{
			ID:            "create-directory",
			Title:         "Create missing directory",
			Description:   fmt.Sprintf("Create the directory: %s", path),
			Command:       fmt.Sprintf("mkdir -p %s", path),
			Automated:     true,
			Confidence:    0.9,
			EstimatedTime: "immediate",
		})
	}

	return fixes
}

// DocumentationProvider provides relevant documentation links
type DocumentationProvider struct{}

// NewDocumentationProvider creates a new documentation provider
func NewDocumentationProvider() *DocumentationProvider {
	return &DocumentationProvider{}
}

// GetRelatedDocumentation gets documentation related to an error
func (dp *DocumentationProvider) GetRelatedDocumentation(err *CLIError, diagnostics *DiagnosticInfo) []Documentation {
	if err == nil {
		return nil
	}

	var docs []Documentation

	switch err.Type {
	case ErrorTypeValidation:
		docs = append(docs, Documentation{
			Title:       "Project Validation Guide",
			URL:         "https://docs.generator.dev/guides/validation",
			Type:        "guide",
			Relevance:   0.9,
			Description: "Learn about project validation and how to fix common issues",
		})
	case ErrorTypeConfiguration:
		docs = append(docs, Documentation{
			Title:       "Configuration Reference",
			URL:         "https://docs.generator.dev/reference/configuration",
			Type:        "reference",
			Relevance:   0.9,
			Description: "Complete reference for configuration file options",
		})
	case ErrorTypeTemplate:
		docs = append(docs, Documentation{
			Title:       "Template System Overview",
			URL:         "https://docs.generator.dev/concepts/templates",
			Type:        "guide",
			Relevance:   0.8,
			Description: "Understanding how templates work and how to use them",
		})
	case ErrorTypeNetwork:
		docs = append(docs, Documentation{
			Title:       "Offline Mode and Caching",
			URL:         "https://docs.generator.dev/guides/offline-mode",
			Type:        "guide",
			Relevance:   0.9,
			Description: "Working without internet connectivity using cached data",
		})
	}

	// Add general troubleshooting guide
	docs = append(docs, Documentation{
		Title:       "Troubleshooting Guide",
		URL:         "https://docs.generator.dev/troubleshooting",
		Type:        "troubleshooting",
		Relevance:   0.7,
		Description: "Common issues and their solutions",
	})

	return docs
}
