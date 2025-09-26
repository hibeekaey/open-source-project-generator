// Package errors provides a unified error handling system for the CLI application
package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ErrorHandler provides a unified interface for comprehensive error handling
type ErrorHandler struct {
	logger           *ErrorLogger
	reporter         *ErrorReporter
	categorizer      *ErrorCategorizer
	recoveryManager  *RecoveryManager
	contextCollector *ErrorContextCollector
	suggestionGen    *ContextualSuggestionGenerator
	recoveryHistory  *RecoveryHistory
	config           *ErrorHandlerConfig
	mutex            sync.RWMutex
}

// ErrorHandlerConfig contains configuration for the error handler
type ErrorHandlerConfig struct {
	LogLevel        LogLevel     `json:"log_level"`
	LogFormat       string       `json:"log_format"` // "text" or "json"
	LogPath         string       `json:"log_path"`
	ReportPath      string       `json:"report_path"`
	ReportFormat    ReportFormat `json:"report_format"`
	MaxReports      int          `json:"max_reports"`
	EnableRecovery  bool         `json:"enable_recovery"`
	EnableReporting bool         `json:"enable_reporting"`
	VerboseMode     bool         `json:"verbose_mode"`
	QuietMode       bool         `json:"quiet_mode"`
	AutoGenerateID  bool         `json:"auto_generate_id"`
}

// DefaultErrorHandlerConfig returns a default configuration
func DefaultErrorHandlerConfig() *ErrorHandlerConfig {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".generator", "logs")
	reportDir := filepath.Join(homeDir, ".generator", "reports")

	return &ErrorHandlerConfig{
		LogLevel:        LogLevelInfo,
		LogFormat:       "text",
		LogPath:         filepath.Join(logDir, "generator.log"),
		ReportPath:      reportDir,
		ReportFormat:    ReportFormatJSON,
		MaxReports:      50,
		EnableRecovery:  true,
		EnableReporting: true,
		VerboseMode:     false,
		QuietMode:       false,
		AutoGenerateID:  true,
	}
}

// NewErrorHandler creates a new comprehensive error handler
func NewErrorHandler(config *ErrorHandlerConfig) (*ErrorHandler, error) {
	if config == nil {
		config = DefaultErrorHandlerConfig()
	}

	// Create logger
	jsonFormat := config.LogFormat == "json"
	logger, err := NewErrorLogger(config.LogPath, config.LogLevel, jsonFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to create error logger: %w", err)
	}

	// Create reporter
	var reporter *ErrorReporter
	if config.EnableReporting {
		reporter = NewErrorReporter(logger, config.ReportPath, config.MaxReports, config.ReportFormat)
	}

	// Create categorizer
	categorizer := NewErrorCategorizer()

	// Create recovery manager
	var recoveryManager *RecoveryManager
	if config.EnableRecovery {
		loggerAdapter := NewLoggerAdapter(logger)
		recoveryManager = NewRecoveryManager(loggerAdapter)
	}

	// Create other components
	contextCollector := NewErrorContextCollector()
	suggestionGen := NewContextualSuggestionGenerator()
	recoveryHistory := NewRecoveryHistory()

	return &ErrorHandler{
		logger:           logger,
		reporter:         reporter,
		categorizer:      categorizer,
		recoveryManager:  recoveryManager,
		contextCollector: contextCollector,
		suggestionGen:    suggestionGen,
		recoveryHistory:  recoveryHistory,
		config:           config,
	}, nil
}

// Close closes the error handler and releases resources
func (eh *ErrorHandler) Close() error {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	if eh.logger != nil {
		return eh.logger.Close()
	}
	return nil
}

// HandleError is the main entry point for handling errors
func (eh *ErrorHandler) HandleError(err error, context map[string]interface{}) *ErrorHandlingResult {
	if err == nil {
		return &ErrorHandlingResult{Success: true}
	}

	startTime := time.Now()

	// Convert to CLIError if needed
	var cliErr *CLIError
	if ce, ok := err.(*CLIError); ok {
		cliErr = ce
	} else {
		cliErr = eh.convertToCLIError(err, context)
	}

	// Enhance error with context and suggestions
	eh.enhanceError(cliErr, context)

	// Record error for statistics
	eh.categorizer.RecordError(cliErr)

	// Log the error
	if eh.logger != nil {
		eh.logger.LogError(cliErr)
	}

	// Attempt recovery if enabled
	var recoveryResult *RecoveryResult
	var recoveryAttempt *RecoveryAttempt
	if eh.config.EnableRecovery && eh.recoveryManager != nil && cliErr.IsRecoverable() {
		recoveryResult, _ = eh.recoveryManager.AttemptRecovery(cliErr)

		if recoveryResult != nil {
			duration := time.Since(startTime)
			recoveryAttempt = &RecoveryAttempt{
				Error:     cliErr,
				Strategy:  "auto",
				Result:    recoveryResult,
				StartTime: startTime,
				Duration:  duration,
			}

			// Record recovery attempt
			eh.recoveryHistory.RecordAttempt(cliErr, "auto", recoveryResult, duration)

			// Log recovery attempt
			if eh.logger != nil {
				eh.logger.LogRecoveryAttempt(recoveryAttempt)
			}
		}
	}

	// Generate error report if enabled
	var reportPath string
	if eh.config.EnableReporting && eh.reporter != nil {
		if report, reportErr := eh.reporter.GenerateReport(cliErr, recoveryAttempt, context); reportErr == nil {
			if path, saveErr := eh.reporter.SaveReport(report); saveErr == nil {
				reportPath = path
			}
		}
	}

	// Create result
	result := &ErrorHandlingResult{
		Success:     recoveryResult != nil && recoveryResult.Success,
		Error:       cliErr,
		Recovery:    recoveryResult,
		ReportPath:  reportPath,
		Duration:    time.Since(startTime),
		Suggestions: cliErr.Suggestions,
		Category:    eh.categorizer.CategorizeError(cliErr),
		ExitCode:    cliErr.ExitCode(),
	}

	// Output error information if not in quiet mode
	if !eh.config.QuietMode {
		eh.outputError(cliErr, result)
	}

	return result
}

// ErrorHandlingResult represents the result of error handling
type ErrorHandlingResult struct {
	Success     bool            `json:"success"`
	Error       *CLIError       `json:"error"`
	Recovery    *RecoveryResult `json:"recovery,omitempty"`
	ReportPath  string          `json:"report_path,omitempty"`
	Duration    time.Duration   `json:"duration"`
	Suggestions []string        `json:"suggestions"`
	Category    *ErrorCategory  `json:"category"`
	ExitCode    int             `json:"exit_code"`
}

// convertToCLIError converts a regular error to a CLIError
func (eh *ErrorHandler) convertToCLIError(err error, context map[string]interface{}) *CLIError {
	// Try to determine error type from error message
	errorType := eh.determineErrorType(err.Error(), context)

	// Create CLIError
	cliErr := NewCLIError(errorType, err.Error(), eh.getExitCodeForType(errorType))
	cliErr = cliErr.WithCause(err)

	return cliErr
}

// determineErrorType attempts to determine error type from message and context
func (eh *ErrorHandler) determineErrorType(message string, context map[string]interface{}) string {
	message = strings.ToLower(message)

	// Check context for hints
	if context != nil {
		if operation, ok := context["operation"].(string); ok {
			switch operation {
			case "validate", "validation":
				return ErrorTypeValidation
			case "config", "configuration":
				return ErrorTypeConfiguration
			case "template":
				return ErrorTypeTemplate
			case "network", "download", "fetch":
				return ErrorTypeNetwork
			case "file", "filesystem", "directory":
				return ErrorTypeFileSystem
			case "cache":
				return ErrorTypeCache
			case "version":
				return ErrorTypeVersion
			case "audit":
				return ErrorTypeAudit
			case "generate", "generation":
				return ErrorTypeGeneration
			}
		}
	}

	// Check message content
	switch {
	case strings.Contains(message, "permission denied"):
		return ErrorTypePermission
	case strings.Contains(message, "not found"):
		return ErrorTypeFileSystem
	case strings.Contains(message, "connection"):
		return ErrorTypeNetwork
	case strings.Contains(message, "invalid"):
		return ErrorTypeValidation
	case strings.Contains(message, "config"):
		return ErrorTypeConfiguration
	case strings.Contains(message, "template"):
		return ErrorTypeTemplate
	case strings.Contains(message, "cache"):
		return ErrorTypeCache
	case strings.Contains(message, "version"):
		return ErrorTypeVersion
	case strings.Contains(message, "security"):
		return ErrorTypeSecurity
	case strings.Contains(message, "dependency"):
		return ErrorTypeDependency
	default:
		return ErrorTypeInternal
	}
}

// getExitCodeForType returns the appropriate exit code for an error type
func (eh *ErrorHandler) getExitCodeForType(errorType string) int {
	switch errorType {
	case ErrorTypeValidation:
		return ExitCodeValidationFailed
	case ErrorTypeConfiguration:
		return ExitCodeConfigurationInvalid
	case ErrorTypeTemplate:
		return ExitCodeTemplateNotFound
	case ErrorTypeNetwork:
		return ExitCodeNetworkError
	case ErrorTypeFileSystem:
		return ExitCodeFileSystemError
	case ErrorTypePermission:
		return ExitCodePermissionDenied
	case ErrorTypeCache:
		return ExitCodeCacheError
	case ErrorTypeVersion:
		return ExitCodeVersionError
	case ErrorTypeAudit:
		return ExitCodeAuditFailed
	case ErrorTypeGeneration:
		return ExitCodeGenerationFailed
	case ErrorTypeDependency:
		return ExitCodeDependencyError
	case ErrorTypeSecurity:
		return ExitCodeSecurityError
	case ErrorTypeUser:
		return ExitCodeUserError
	default:
		return ExitCodeInternalError
	}
}

// enhanceError enhances an error with context and suggestions
func (eh *ErrorHandler) enhanceError(err *CLIError, context map[string]interface{}) {
	// Add context if not present
	if err.Context == nil {
		err.Context = eh.contextCollector.CollectContext()
	}

	// Add additional context from parameters
	for key, value := range context {
		err = err.WithDetails(key, value)
	}

	// Generate additional suggestions
	if eh.suggestionGen != nil {
		suggestions := eh.suggestionGen.GenerateSuggestions(err)
		for _, suggestion := range suggestions {
			// Avoid duplicates
			found := false
			for _, existing := range err.Suggestions {
				if existing == suggestion {
					found = true
					break
				}
			}
			if !found {
				err.Suggestions = append(err.Suggestions, suggestion)
			}
		}
	}
}

// outputError outputs error information to the user
func (eh *ErrorHandler) outputError(err *CLIError, result *ErrorHandlingResult) {
	// Format error for user display
	errorMsg := FormatErrorForUser(err, eh.config.VerboseMode)

	// Output to stderr
	fmt.Fprintln(os.Stderr, errorMsg)

	// Show recovery information if available
	if result.Recovery != nil {
		if result.Recovery.Success {
			fmt.Fprintf(os.Stderr, "\n✅ Recovery successful: %s\n", result.Recovery.Message)
		} else {
			fmt.Fprintf(os.Stderr, "\n❌ Recovery failed: %s\n", result.Recovery.Message)
		}

		if len(result.Recovery.Suggestions) > 0 {
			fmt.Fprintln(os.Stderr, "\nRecovery suggestions:")
			for _, suggestion := range result.Recovery.Suggestions {
				fmt.Fprintf(os.Stderr, "  • %s\n", suggestion)
			}
		}
	}

	// Show report path if available
	if result.ReportPath != "" && eh.config.VerboseMode {
		fmt.Fprintf(os.Stderr, "\n📄 Error report saved to: %s\n", result.ReportPath)
	}
}

// GetStatistics returns error statistics
func (eh *ErrorHandler) GetStatistics() *ErrorStatistics {
	return eh.categorizer.GetStatistics()
}

// GetRecoveryHistory returns recovery history
func (eh *ErrorHandler) GetRecoveryHistory() *RecoveryHistory {
	return eh.recoveryHistory
}

// GenerateAnalysisReport generates a comprehensive error analysis report
func (eh *ErrorHandler) GenerateAnalysisReport() *ErrorAnalysisReport {
	return eh.categorizer.GenerateErrorReport()
}

// SetVerboseMode sets verbose mode
func (eh *ErrorHandler) SetVerboseMode(verbose bool) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.config.VerboseMode = verbose
}

// SetQuietMode sets quiet mode
func (eh *ErrorHandler) SetQuietMode(quiet bool) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.config.QuietMode = quiet
}

// SetLogLevel sets the logging level
func (eh *ErrorHandler) SetLogLevel(level LogLevel) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()
	eh.config.LogLevel = level
	if eh.logger != nil {
		eh.logger.SetLevel(level)
	}
}

// GetLogPath returns the path to the log file
func (eh *ErrorHandler) GetLogPath() string {
	return eh.config.LogPath
}

// GetReportPath returns the path to the report directory
func (eh *ErrorHandler) GetReportPath() string {
	return eh.config.ReportPath
}

// Convenience methods for creating specific error types

// NewValidationErrorHandler creates a validation error and handles it
func (eh *ErrorHandler) NewValidationErrorHandler(message, field string, value interface{}) *ErrorHandlingResult {
	err := NewValidationError(message, field, value)
	return eh.HandleError(err, map[string]interface{}{
		"operation": "validation",
		"field":     field,
		"value":     value,
	})
}

// NewConfigurationErrorHandler creates a configuration error and handles it
func (eh *ErrorHandler) NewConfigurationErrorHandler(message, configPath string, cause error) *ErrorHandlingResult {
	err := NewConfigurationError(message, configPath, cause)
	return eh.HandleError(err, map[string]interface{}{
		"operation":   "configuration",
		"config_path": configPath,
	})
}

// NewTemplateErrorHandler creates a template error and handles it
func (eh *ErrorHandler) NewTemplateErrorHandler(message, templateName string, cause error) *ErrorHandlingResult {
	err := NewTemplateError(message, templateName, cause)
	return eh.HandleError(err, map[string]interface{}{
		"operation":     "template",
		"template_name": templateName,
	})
}

// NewNetworkErrorHandler creates a network error and handles it
func (eh *ErrorHandler) NewNetworkErrorHandler(message, url string, cause error) *ErrorHandlingResult {
	err := NewNetworkError(message, url, cause)
	return eh.HandleError(err, map[string]interface{}{
		"operation": "network",
		"url":       url,
	})
}

// NewFileSystemErrorHandler creates a filesystem error and handles it
func (eh *ErrorHandler) NewFileSystemErrorHandler(message, path, operation string, cause error) *ErrorHandlingResult {
	err := NewFileSystemError(message, path, operation, cause)
	return eh.HandleError(err, map[string]interface{}{
		"operation": "filesystem",
		"path":      path,
		"fs_op":     operation,
	})
}

// NewGenerationErrorHandler creates a generation error and handles it
func (eh *ErrorHandler) NewGenerationErrorHandler(message, component, operation string, cause error) *ErrorHandlingResult {
	err := NewGenerationError(message, component, operation, cause)
	return eh.HandleError(err, map[string]interface{}{
		"operation": "generation",
		"component": component,
		"gen_op":    operation,
	})
}

// Global error handler instance
var globalErrorHandler *ErrorHandler
var globalErrorHandlerOnce sync.Once

// InitializeGlobalErrorHandler initializes the global error handler
func InitializeGlobalErrorHandler(config *ErrorHandlerConfig) error {
	var err error
	globalErrorHandlerOnce.Do(func() {
		globalErrorHandler, err = NewErrorHandler(config)
	})
	return err
}

// GetGlobalErrorHandler returns the global error handler
func GetGlobalErrorHandler() *ErrorHandler {
	return globalErrorHandler
}

// HandleGlobalError handles an error using the global error handler
func HandleGlobalError(err error, context map[string]interface{}) *ErrorHandlingResult {
	if globalErrorHandler == nil {
		// Fallback: create a minimal error handler
		config := DefaultErrorHandlerConfig()
		config.EnableReporting = false // Disable reporting for fallback
		if handler, handlerErr := NewErrorHandler(config); handlerErr == nil {
			return handler.HandleError(err, context)
		}

		// Ultimate fallback: just print the error
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return &ErrorHandlingResult{
			Success:  false,
			ExitCode: ExitCodeGeneral,
		}
	}

	return globalErrorHandler.HandleError(err, context)
}
