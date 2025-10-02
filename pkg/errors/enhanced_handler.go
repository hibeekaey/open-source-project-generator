// Package errors provides enhanced error handling with comprehensive user experience
package errors

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// EnhancedErrorHandler provides comprehensive error handling with user experience focus
type EnhancedErrorHandler struct {
	baseHandler        *ErrorHandler
	diagnostics        *DiagnosticCollector
	userExperience     *UserExperienceManager
	performanceTracker *PerformanceTracker
	config             *EnhancedErrorConfig
	logger             interfaces.Logger
	mutex              sync.RWMutex
}

// EnhancedErrorConfig contains configuration for enhanced error handling
type EnhancedErrorConfig struct {
	*ErrorHandlerConfig

	// User experience settings
	ShowSuggestions     bool `json:"show_suggestions"`
	ShowRecoveryOptions bool `json:"show_recovery_options"`
	ShowDiagnostics     bool `json:"show_diagnostics"`
	InteractiveMode     bool `json:"interactive_mode"`

	// Performance monitoring
	EnablePerformanceTracking bool                     `json:"enable_performance_tracking"`
	PerformanceThresholds     map[string]time.Duration `json:"performance_thresholds"`

	// Diagnostic settings
	CollectSystemInfo  bool `json:"collect_system_info"`
	CollectEnvironment bool `json:"collect_environment"`
	CollectStackTraces bool `json:"collect_stack_traces"`

	// Localization
	Language string `json:"language"`
	Locale   string `json:"locale"`
}

// DefaultEnhancedErrorConfig returns default configuration for enhanced error handling
func DefaultEnhancedErrorConfig() *EnhancedErrorConfig {
	baseConfig := DefaultErrorHandlerConfig()

	return &EnhancedErrorConfig{
		ErrorHandlerConfig:        baseConfig,
		ShowSuggestions:           true,
		ShowRecoveryOptions:       true,
		ShowDiagnostics:           false,
		InteractiveMode:           true,
		EnablePerformanceTracking: true,
		PerformanceThresholds: map[string]time.Duration{
			"command_execution": 5 * time.Second,
			"file_operation":    1 * time.Second,
			"network_request":   10 * time.Second,
			"validation":        2 * time.Second,
		},
		CollectSystemInfo:  true,
		CollectEnvironment: true,
		CollectStackTraces: true,
		Language:           "en",
		Locale:             "en_US",
	}
}

// NewEnhancedErrorHandler creates a new enhanced error handler
func NewEnhancedErrorHandler(config *EnhancedErrorConfig, logger interfaces.Logger) (*EnhancedErrorHandler, error) {
	if config == nil {
		config = DefaultEnhancedErrorConfig()
	}

	// Create base error handler
	baseHandler, err := NewErrorHandler(config.ErrorHandlerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create base error handler: %w", err)
	}

	// Create diagnostic collector
	diagnostics := NewDiagnosticCollector(config, logger)

	// Create user experience manager
	userExperience := NewUserExperienceManager(config, logger)

	// Create performance tracker
	performanceTracker := NewPerformanceTracker(config, logger)

	return &EnhancedErrorHandler{
		baseHandler:        baseHandler,
		diagnostics:        diagnostics,
		userExperience:     userExperience,
		performanceTracker: performanceTracker,
		config:             config,
		logger:             logger,
	}, nil
}

// HandleError provides enhanced error handling with comprehensive user experience
func (eh *EnhancedErrorHandler) HandleError(ctx context.Context, err error, operation string, context map[string]interface{}) *EnhancedErrorResult {
	if err == nil {
		return &EnhancedErrorResult{Success: true}
	}

	startTime := time.Now()

	// Collect diagnostic information
	diagnosticInfo := eh.diagnostics.CollectDiagnostics(ctx, err, operation, context)

	// Handle error with base handler
	baseResult := eh.baseHandler.HandleError(err, context)

	// Enhance with user experience improvements
	userExperience := eh.userExperience.EnhanceUserExperience(baseResult.Error, diagnosticInfo)

	// Track performance
	duration := time.Since(startTime)
	performanceInfo := eh.performanceTracker.TrackOperation(operation, duration, context)

	// Create enhanced result
	result := &EnhancedErrorResult{
		Success:        baseResult.Success,
		Error:          baseResult.Error,
		Recovery:       baseResult.Recovery,
		ReportPath:     baseResult.ReportPath,
		Duration:       baseResult.Duration,
		Suggestions:    baseResult.Suggestions,
		Category:       baseResult.Category,
		ExitCode:       baseResult.ExitCode,
		Diagnostics:    diagnosticInfo,
		UserExperience: userExperience,
		Performance:    performanceInfo,
	}

	// Display enhanced error information
	eh.displayEnhancedError(result)

	return result
}

// EnhancedErrorResult represents the result of enhanced error handling
type EnhancedErrorResult struct {
	Success        bool                `json:"success"`
	Error          *CLIError           `json:"error"`
	Recovery       *RecoveryResult     `json:"recovery,omitempty"`
	ReportPath     string              `json:"report_path,omitempty"`
	Duration       time.Duration       `json:"duration"`
	Suggestions    []string            `json:"suggestions"`
	Category       *ErrorCategory      `json:"category"`
	ExitCode       int                 `json:"exit_code"`
	Diagnostics    *DiagnosticInfo     `json:"diagnostics,omitempty"`
	UserExperience *UserExperienceInfo `json:"user_experience,omitempty"`
	Performance    *PerformanceInfo    `json:"performance,omitempty"`
}

// displayEnhancedError displays comprehensive error information to the user
func (eh *EnhancedErrorHandler) displayEnhancedError(result *EnhancedErrorResult) {
	if eh.config.QuietMode {
		return
	}

	// Display main error
	eh.displayMainError(result.Error)

	// Display suggestions if enabled
	if eh.config.ShowSuggestions && len(result.Suggestions) > 0 {
		eh.displaySuggestions(result.Suggestions)
	}

	// Display recovery options if enabled
	if eh.config.ShowRecoveryOptions && result.Recovery != nil {
		eh.displayRecoveryInfo(result.Recovery)
	}

	// Display diagnostics if enabled and in verbose mode
	if eh.config.ShowDiagnostics && eh.config.VerboseMode && result.Diagnostics != nil {
		eh.displayDiagnostics(result.Diagnostics)
	}

	// Display performance information if enabled and in debug mode
	if eh.config.EnablePerformanceTracking && result.Performance != nil {
		eh.displayPerformanceInfo(result.Performance)
	}

	// Display user experience enhancements
	if result.UserExperience != nil {
		eh.displayUserExperienceInfo(result.UserExperience)
	}
}

// displayMainError displays the main error message with formatting
func (eh *EnhancedErrorHandler) displayMainError(err *CLIError) {
	if err == nil {
		return
	}

	// Format error message based on severity
	var icon string
	switch err.Severity {
	case SeverityLow:
		icon = "ℹ️"
	case SeverityMedium:
		icon = "⚠️"
	case SeverityHigh:
		icon = "🚨"
	case SeverityCritical:
		icon = "🔥"
	default:
		icon = "❌"
	}

	fmt.Fprintf(os.Stderr, "\n%s %s\n", icon, err.Message)

	// Show error type and code in verbose mode
	if eh.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "   Type: %s (Code: %d)\n", err.Type, err.Code)
		fmt.Fprintf(os.Stderr, "   Severity: %s\n", err.Severity)

		if err.Context != nil && err.Context.Operation != "" {
			fmt.Fprintf(os.Stderr, "   Operation: %s\n", err.Context.Operation)
		}
	}
}

// displaySuggestions displays actionable suggestions to the user
func (eh *EnhancedErrorHandler) displaySuggestions(suggestions []string) {
	if len(suggestions) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "\n💡 Suggestions:\n")
	for i, suggestion := range suggestions {
		if i < 5 { // Limit to top 5 suggestions to avoid overwhelming
			fmt.Fprintf(os.Stderr, "   %d. %s\n", i+1, suggestion)
		}
	}

	if len(suggestions) > 5 {
		fmt.Fprintf(os.Stderr, "   ... and %d more suggestions (use --verbose for all)\n", len(suggestions)-5)
	}
}

// displayRecoveryInfo displays recovery information
func (eh *EnhancedErrorHandler) displayRecoveryInfo(recovery *RecoveryResult) {
	if recovery == nil {
		return
	}

	if recovery.Success {
		fmt.Fprintf(os.Stderr, "\n✅ Recovery: %s\n", recovery.Message)
	} else {
		fmt.Fprintf(os.Stderr, "\n❌ Recovery failed: %s\n", recovery.Message)

		if len(recovery.Suggestions) > 0 {
			fmt.Fprintf(os.Stderr, "\n🔧 Recovery suggestions:\n")
			for _, suggestion := range recovery.Suggestions {
				fmt.Fprintf(os.Stderr, "   • %s\n", suggestion)
			}
		}
	}
}

// displayDiagnostics displays diagnostic information
func (eh *EnhancedErrorHandler) displayDiagnostics(diagnostics *DiagnosticInfo) {
	if diagnostics == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\n🔍 Diagnostics:\n")

	if diagnostics.SystemInfo != nil {
		fmt.Fprintf(os.Stderr, "   OS: %s %s\n", diagnostics.SystemInfo.OS, diagnostics.SystemInfo.Arch)
		fmt.Fprintf(os.Stderr, "   Memory: %s\n", formatBytes(diagnostics.SystemInfo.MemoryUsage))
	}

	if diagnostics.EnvironmentInfo != nil {
		fmt.Fprintf(os.Stderr, "   Working Dir: %s\n", diagnostics.EnvironmentInfo.WorkingDir)

		if diagnostics.EnvironmentInfo.CI != nil && diagnostics.EnvironmentInfo.CI.IsCI {
			fmt.Fprintf(os.Stderr, "   CI Environment: %s\n", diagnostics.EnvironmentInfo.CI.Provider)
		}
	}

	if len(diagnostics.RelevantLogs) > 0 {
		fmt.Fprintf(os.Stderr, "   Recent logs: %d entries\n", len(diagnostics.RelevantLogs))
	}
}

// displayPerformanceInfo displays performance information
func (eh *EnhancedErrorHandler) displayPerformanceInfo(performance *PerformanceInfo) {
	if performance == nil {
		return
	}

	if performance.IsSlowOperation {
		fmt.Fprintf(os.Stderr, "\n⚡ Performance: Operation took %v (threshold: %v)\n",
			performance.Duration, performance.Threshold)

		if len(performance.Suggestions) > 0 {
			fmt.Fprintf(os.Stderr, "   Performance suggestions:\n")
			for _, suggestion := range performance.Suggestions {
				fmt.Fprintf(os.Stderr, "   • %s\n", suggestion)
			}
		}
	}
}

// displayUserExperienceInfo displays user experience enhancements
func (eh *EnhancedErrorHandler) displayUserExperienceInfo(ux *UserExperienceInfo) {
	if ux == nil {
		return
	}

	// Display contextual help if available
	if ux.ContextualHelp != "" {
		fmt.Fprintf(os.Stderr, "\n📖 Help: %s\n", ux.ContextualHelp)
	}

	// Display quick fixes if available
	if len(ux.QuickFixes) > 0 {
		fmt.Fprintf(os.Stderr, "\n🔧 Quick fixes:\n")
		for i, fix := range ux.QuickFixes {
			fmt.Fprintf(os.Stderr, "   %d. %s\n", i+1, fix.Description)
			if eh.config.VerboseMode && fix.Command != "" {
				fmt.Fprintf(os.Stderr, "      Command: %s\n", fix.Command)
			}
		}
	}

	// Display related documentation if available
	if len(ux.RelatedDocs) > 0 && eh.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "\n📚 Related documentation:\n")
		for _, doc := range ux.RelatedDocs {
			fmt.Fprintf(os.Stderr, "   • %s: %s\n", doc.Title, doc.URL)
		}
	}
}

// Close closes the enhanced error handler and releases resources
func (eh *EnhancedErrorHandler) Close() error {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	var errors []error

	if eh.baseHandler != nil {
		if err := eh.baseHandler.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if eh.diagnostics != nil {
		if err := eh.diagnostics.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if eh.performanceTracker != nil {
		if err := eh.performanceTracker.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing enhanced error handler: %v", errors)
	}

	return nil
}

// GetStatistics returns comprehensive error statistics
func (eh *EnhancedErrorHandler) GetStatistics() *EnhancedErrorStatistics {
	baseStats := eh.baseHandler.GetStatistics()

	return &EnhancedErrorStatistics{
		ErrorStatistics:     baseStats,
		DiagnosticStats:     eh.diagnostics.GetStatistics(),
		PerformanceStats:    eh.performanceTracker.GetStatistics(),
		UserExperienceStats: eh.userExperience.GetStatistics(),
	}
}

// EnhancedErrorStatistics contains comprehensive error statistics
type EnhancedErrorStatistics struct {
	*ErrorStatistics
	DiagnosticStats     *DiagnosticStatistics     `json:"diagnostic_stats"`
	PerformanceStats    *PerformanceStatistics    `json:"performance_stats"`
	UserExperienceStats *UserExperienceStatistics `json:"user_experience_stats"`
}

// SetVerboseMode sets verbose mode for enhanced error handling
func (eh *EnhancedErrorHandler) SetVerboseMode(verbose bool) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	eh.config.VerboseMode = verbose
	eh.baseHandler.SetVerboseMode(verbose)

	if eh.diagnostics != nil {
		eh.diagnostics.SetVerboseMode(verbose)
	}
}

// SetQuietMode sets quiet mode for enhanced error handling
func (eh *EnhancedErrorHandler) SetQuietMode(quiet bool) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	eh.config.QuietMode = quiet
	eh.baseHandler.SetQuietMode(quiet)
}

// SetInteractiveMode sets interactive mode for enhanced error handling
func (eh *EnhancedErrorHandler) SetInteractiveMode(interactive bool) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	eh.config.InteractiveMode = interactive

	if eh.userExperience != nil {
		eh.userExperience.SetInteractiveMode(interactive)
	}
}

// Helper function to format bytes
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Convenience methods for creating enhanced errors

// NewEnhancedValidationError creates an enhanced validation error
func (eh *EnhancedErrorHandler) NewEnhancedValidationError(ctx context.Context, message, field string, value interface{}) *EnhancedErrorResult {
	err := NewValidationError(message, field, value)
	return eh.HandleError(ctx, err, "validation", map[string]interface{}{
		"field": field,
		"value": value,
	})
}

// NewEnhancedConfigurationError creates an enhanced configuration error
func (eh *EnhancedErrorHandler) NewEnhancedConfigurationError(ctx context.Context, message, configPath string, cause error) *EnhancedErrorResult {
	err := NewConfigurationError(message, configPath, cause)
	return eh.HandleError(ctx, err, "configuration", map[string]interface{}{
		"config_path": configPath,
	})
}

// NewEnhancedNetworkError creates an enhanced network error
func (eh *EnhancedErrorHandler) NewEnhancedNetworkError(ctx context.Context, message, url string, cause error) *EnhancedErrorResult {
	err := NewNetworkError(message, url, cause)
	return eh.HandleError(ctx, err, "network", map[string]interface{}{
		"url": url,
	})
}

// NewEnhancedFileSystemError creates an enhanced filesystem error
func (eh *EnhancedErrorHandler) NewEnhancedFileSystemError(ctx context.Context, message, path, operation string, cause error) *EnhancedErrorResult {
	err := NewFileSystemError(message, path, operation, cause)
	return eh.HandleError(ctx, err, "filesystem", map[string]interface{}{
		"path":      path,
		"operation": operation,
	})
}

// Global enhanced error handler instance
var globalEnhancedErrorHandler *EnhancedErrorHandler
var globalEnhancedErrorHandlerOnce sync.Once

// InitializeGlobalEnhancedErrorHandler initializes the global enhanced error handler
func InitializeGlobalEnhancedErrorHandler(config *EnhancedErrorConfig, logger interfaces.Logger) error {
	var err error
	globalEnhancedErrorHandlerOnce.Do(func() {
		globalEnhancedErrorHandler, err = NewEnhancedErrorHandler(config, logger)
	})
	return err
}

// GetGlobalEnhancedErrorHandler returns the global enhanced error handler
func GetGlobalEnhancedErrorHandler() *EnhancedErrorHandler {
	return globalEnhancedErrorHandler
}

// HandleGlobalEnhancedError handles an error using the global enhanced error handler
func HandleGlobalEnhancedError(ctx context.Context, err error, operation string, context map[string]interface{}) *EnhancedErrorResult {
	if globalEnhancedErrorHandler == nil {
		// Fallback to basic error handling
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return &EnhancedErrorResult{
			Success:  false,
			ExitCode: ExitCodeGeneral,
		}
	}

	return globalEnhancedErrorHandler.HandleError(ctx, err, operation, context)
}
