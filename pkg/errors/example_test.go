package errors_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/errors"
)

// Example demonstrates basic usage of the error handling system
func ExampleErrorHandler() {
	// Create a temporary directory for testing
	tempDir, _ := os.MkdirTemp("", "error-handler-test")
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Configure error handler
	config := &errors.ErrorHandlerConfig{
		LogLevel:        errors.LogLevelInfo,
		LogFormat:       "text",
		LogPath:         filepath.Join(tempDir, "test.log"),
		ReportPath:      filepath.Join(tempDir, "reports"),
		ReportFormat:    errors.ReportFormatText,
		MaxReports:      10,
		EnableRecovery:  true,
		EnableReporting: true,
		VerboseMode:     false,
		QuietMode:       true, // Quiet for test
	}

	// Create error handler
	handler, err := errors.NewErrorHandler(config)
	if err != nil {
		fmt.Printf("Failed to create error handler: %v\n", err)
		return
	}
	defer func() { _ = handler.Close() }()

	// Create a validation error
	validationErr := errors.NewValidationError(
		"Invalid project name format",
		"project_name",
		"invalid-name!",
	)

	// Handle the error
	result := handler.HandleError(validationErr, map[string]interface{}{
		"operation": "project_validation",
		"user_id":   "test-user",
	})

	// Check results
	fmt.Printf("Error handled successfully: %t\n", result.Error != nil)
	fmt.Printf("Error type: %s\n", result.Error.Type)
	fmt.Printf("Error severity: %s\n", result.Error.Severity)
	fmt.Printf("Error recoverable: %t\n", result.Error.Recoverable)
	fmt.Printf("Suggestions provided: %t\n", len(result.Suggestions) > 0)

	// Output:
	// Error handled successfully: true
	// Error type: validation
	// Error severity: medium
	// Error recoverable: true
	// Suggestions provided: true
}

// TestErrorHandlerBasicFunctionality tests basic error handler functionality
func TestErrorHandlerBasicFunctionality(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "error-handler-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Configure error handler
	config := &errors.ErrorHandlerConfig{
		LogLevel:        errors.LogLevelDebug,
		LogFormat:       "json",
		LogPath:         filepath.Join(tempDir, "test.log"),
		ReportPath:      filepath.Join(tempDir, "reports"),
		ReportFormat:    errors.ReportFormatJSON,
		MaxReports:      5,
		EnableRecovery:  true,
		EnableReporting: true,
		VerboseMode:     true,
		QuietMode:       true, // Quiet for test
	}

	// Create error handler
	handler, err := errors.NewErrorHandler(config)
	if err != nil {
		t.Fatalf("Failed to create error handler: %v", err)
	}
	defer func() { _ = handler.Close() }()

	// Test different error types
	testCases := []struct {
		name     string
		err      *errors.CLIError
		context  map[string]interface{}
		expected struct {
			errorType   string
			severity    errors.Severity
			recoverable bool
		}
	}{
		{
			name: "Validation Error",
			err: errors.NewValidationError(
				"Invalid configuration value",
				"timeout",
				-1,
			),
			context: map[string]interface{}{
				"operation": "config_validation",
			},
			expected: struct {
				errorType   string
				severity    errors.Severity
				recoverable bool
			}{
				errorType:   errors.ErrorTypeValidation,
				severity:    errors.SeverityMedium,
				recoverable: true,
			},
		},
		{
			name: "Network Error",
			err: errors.NewNetworkError(
				"Failed to connect to remote server",
				"https://api.example.com",
				fmt.Errorf("connection timeout"),
			),
			context: map[string]interface{}{
				"operation": "api_call",
				"retry":     1,
			},
			expected: struct {
				errorType   string
				severity    errors.Severity
				recoverable bool
			}{
				errorType:   errors.ErrorTypeNetwork,
				severity:    errors.SeverityMedium,
				recoverable: true,
			},
		},
		{
			name: "Security Error",
			err: errors.NewSecurityError(
				"Potential security vulnerability detected",
				"path_traversal",
				"high",
				fmt.Errorf("unsafe path detected"),
			),
			context: map[string]interface{}{
				"operation": "file_access",
				"path":      "../../../etc/passwd",
			},
			expected: struct {
				errorType   string
				severity    errors.Severity
				recoverable bool
			}{
				errorType:   errors.ErrorTypeSecurity,
				severity:    errors.SeverityCritical,
				recoverable: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Handle the error
			result := handler.HandleError(tc.err, tc.context)

			// Verify results
			if result.Error == nil {
				t.Fatal("Expected error to be present in result")
			}

			if result.Error.Type != tc.expected.errorType {
				t.Errorf("Expected error type %s, got %s", tc.expected.errorType, result.Error.Type)
			}

			if result.Error.Severity != tc.expected.severity {
				t.Errorf("Expected severity %s, got %s", tc.expected.severity, result.Error.Severity)
			}

			if result.Error.Recoverable != tc.expected.recoverable {
				t.Errorf("Expected recoverable %t, got %t", tc.expected.recoverable, result.Error.Recoverable)
			}

			// Verify suggestions are provided
			if len(result.Suggestions) == 0 {
				t.Error("Expected suggestions to be provided")
			}

			// Verify category is assigned
			if result.Category == nil {
				t.Error("Expected error category to be assigned")
			}

			// Verify exit code is set
			if result.ExitCode == 0 {
				t.Error("Expected non-zero exit code for error")
			}
		})
	}

	// Test statistics
	stats := handler.GetStatistics()
	if stats.TotalErrors != len(testCases) {
		t.Errorf("Expected %d total errors, got %d", len(testCases), stats.TotalErrors)
	}

	// Test recovery history
	history := handler.GetRecoveryHistory()
	if history == nil {
		t.Error("Expected recovery history to be available")
	}

	// Test analysis report
	report := handler.GenerateAnalysisReport()
	if report == nil {
		t.Error("Expected analysis report to be generated")
		return
	}

	if report.TotalErrors != len(testCases) {
		t.Errorf("Expected %d total errors in report, got %d", len(testCases), report.TotalErrors)
	}
}

// TestContextualErrorBuilder tests the contextual error builder
func TestContextualErrorBuilder(t *testing.T) {
	// Build a complex error with context
	err := errors.NewContextualError(
		errors.ErrorTypeGeneration,
		"Failed to generate project files",
		errors.ExitCodeGenerationFailed,
	).
		WithOperation("project_generation").
		WithComponent("template_processor").
		WithFile("template.go", 42).
		WithDetail("template_name", "go-gin-api").
		WithDetail("output_dir", "/tmp/test-project").
		WithSuggestion("Check template syntax and variables").
		WithSuggestion("Verify output directory permissions").
		WithSeverity(errors.SeverityHigh).
		WithRecoverable(true).
		Build()

	// Verify error properties
	if err.Type != errors.ErrorTypeGeneration {
		t.Errorf("Expected error type %s, got %s", errors.ErrorTypeGeneration, err.Type)
	}

	if err.Code != errors.ExitCodeGenerationFailed {
		t.Errorf("Expected exit code %d, got %d", errors.ExitCodeGenerationFailed, err.Code)
	}

	if err.Severity != errors.SeverityHigh {
		t.Errorf("Expected severity %s, got %s", errors.SeverityHigh, err.Severity)
	}

	if !err.Recoverable {
		t.Error("Expected error to be recoverable")
	}

	// Verify context
	if err.Context == nil {
		t.Fatal("Expected context to be present")
	}

	if err.Context.Operation != "project_generation" {
		t.Errorf("Expected operation 'project_generation', got %s", err.Context.Operation)
	}

	if err.Context.Component != "template_processor" {
		t.Errorf("Expected component 'template_processor', got %s", err.Context.Component)
	}

	if err.Context.File != "template.go" {
		t.Errorf("Expected file 'template.go', got %s", err.Context.File)
	}

	if err.Context.Line != 42 {
		t.Errorf("Expected line 42, got %d", err.Context.Line)
	}

	// Verify details
	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	if err.Details["template_name"] != "go-gin-api" {
		t.Errorf("Expected template_name 'go-gin-api', got %v", err.Details["template_name"])
	}

	// Verify suggestions
	if len(err.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(err.Suggestions))
	}
}
