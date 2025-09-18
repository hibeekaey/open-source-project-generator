// Package errors provides comprehensive error handling for the CLI application.
//
// This package implements a complete error handling system that includes:
//
// # Error Types and Categorization
//
// The package defines various error types for different categories of failures:
//   - ErrorTypeValidation: Input validation and project structure errors
//   - ErrorTypeConfiguration: Configuration file and settings errors
//   - ErrorTypeTemplate: Template processing and management errors
//   - ErrorTypeNetwork: Network connectivity and API errors
//   - ErrorTypeFileSystem: File and directory operation errors
//   - ErrorTypePermission: Permission and access control errors
//   - ErrorTypeCache: Cache management and storage errors
//   - ErrorTypeVersion: Version checking and update errors
//   - ErrorTypeAudit: Security and quality audit errors
//   - ErrorTypeGeneration: Project generation and processing errors
//   - ErrorTypeDependency: Dependency management and compatibility errors
//   - ErrorTypeSecurity: Security vulnerabilities and policy violations
//   - ErrorTypeUser: User input and interaction errors
//   - ErrorTypeInternal: Internal system errors and unexpected conditions
//
// # Error Severity Levels
//
// Errors are classified by severity to help prioritize resolution:
//   - SeverityLow: Minor issues that don't prevent operation
//   - SeverityMedium: Issues that may affect functionality
//   - SeverityHigh: Serious issues that prevent normal operation
//   - SeverityCritical: Critical issues requiring immediate attention
//
// # Error Recovery
//
// The package includes automatic error recovery mechanisms:
//   - NetworkRecoveryStrategy: Handles network failures with offline fallback
//   - FileSystemRecoveryStrategy: Creates missing directories and handles permissions
//   - CacheRecoveryStrategy: Clears corrupted cache data
//   - ConfigurationRecoveryStrategy: Creates default configurations
//   - TemplateRecoveryStrategy: Suggests alternative templates
//   - ValidationRecoveryStrategy: Provides valid value suggestions
//
// # Contextual Error Information
//
// Errors include rich contextual information:
//   - Command and arguments that caused the error
//   - Working directory and environment details
//   - CI/CD environment detection
//   - File and line number information
//   - Operation and component context
//
// # Error Logging and Reporting
//
// Comprehensive logging and reporting capabilities:
//   - Structured logging in text or JSON format
//   - Automatic error report generation
//   - Error statistics and pattern analysis
//   - Recovery attempt tracking
//   - Multiple output formats (text, JSON, HTML)
//
// # Usage Examples
//
// Basic error handling:
//
//	handler, err := errors.NewErrorHandler(errors.DefaultErrorHandlerConfig())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer handler.Close()
//
//	// Handle a validation error
//	result := handler.NewValidationErrorHandler(
//	    "Invalid project name",
//	    "name",
//	    "invalid-name!",
//	)
//
//	if !result.Success {
//	    os.Exit(result.ExitCode)
//	}
//
// Creating custom errors with context:
//
//	err := errors.NewContextualError(errors.ErrorTypeTemplate, "Template not found", errors.ExitCodeTemplateNotFound).
//	    WithOperation("generate").
//	    WithComponent("template-manager").
//	    WithDetail("template_name", "invalid-template").
//	    WithSuggestion("Use 'generator list-templates' to see available templates").
//	    Build()
//
//	result := handler.HandleError(err, map[string]interface{}{
//	    "operation": "generate",
//	    "user_input": "invalid-template",
//	})
//
// Global error handling:
//
//	// Initialize global handler
//	errors.InitializeGlobalErrorHandler(errors.DefaultErrorHandlerConfig())
//
//	// Use global handler
//	result := errors.HandleGlobalError(someError, map[string]interface{}{
//	    "operation": "validate",
//	})
//
// # Error Statistics and Analysis
//
// The package tracks error statistics for analysis:
//
//	stats := handler.GetStatistics()
//	fmt.Printf("Total errors: %d\n", stats.TotalErrors)
//	fmt.Printf("Recovery rate: %.1f%%\n", stats.RecoveryRate)
//
//	report := handler.GenerateAnalysisReport()
//	// Use report for insights and improvements
//
// # Integration with CLI
//
// The error handling system integrates seamlessly with the CLI:
//   - Automatic context collection from command execution
//   - CI/CD environment detection and adaptation
//   - Machine-readable output for automation
//   - Verbose and quiet mode support
//   - Exit code management for scripts
//
// # Configuration
//
// The error handler can be configured for different environments:
//
//	config := &errors.ErrorHandlerConfig{
//	    LogLevel:        errors.LogLevelDebug,
//	    LogFormat:       "json",
//	    EnableRecovery:  true,
//	    EnableReporting: true,
//	    VerboseMode:     true,
//	}
//
//	handler, err := errors.NewErrorHandler(config)
//
// # Thread Safety
//
// All components in this package are thread-safe and can be used
// concurrently from multiple goroutines.
package errors
