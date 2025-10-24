package orchestrator

import (
	"errors"
	"testing"
)

// TestErrorWithContextAndSuggestions demonstrates how errors include context and suggestions
func TestErrorWithContextAndSuggestions(t *testing.T) {
	tests := []struct {
		name            string
		createError     func() *GenerationError
		wantComponent   string
		wantSuggestions int
		wantRecoverable bool
	}{
		{
			name: "tool not found with context",
			createError: func() *GenerationError {
				return NewToolNotFoundError("npx", "nextjs-app")
			},
			wantComponent:   "nextjs-app",
			wantSuggestions: 3,
			wantRecoverable: true,
		},
		{
			name: "tool execution error with context",
			createError: func() *GenerationError {
				return NewToolExecutionError("go", "backend-api", errors.New("exit code 1"))
			},
			wantComponent:   "backend-api",
			wantSuggestions: 4,
			wantRecoverable: true,
		},
		{
			name: "config validation error with suggestions",
			createError: func() *GenerationError {
				return NewConfigValidationError("output_dir", "path must be absolute")
			},
			wantComponent:   "",
			wantSuggestions: 3,
			wantRecoverable: false,
		},
		{
			name: "security error with context",
			createError: func() *GenerationError {
				return NewSecurityError("path traversal detected in '../../../etc/passwd'", nil)
			},
			wantComponent:   "",
			wantSuggestions: 3,
			wantRecoverable: false,
		},
		{
			name: "timeout error with component context",
			createError: func() *GenerationError {
				return NewTimeoutError("npm install", "frontend")
			},
			wantComponent:   "frontend",
			wantSuggestions: 3,
			wantRecoverable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()

			// Verify component context
			if err.Component != tt.wantComponent {
				t.Errorf("Component = %q, want %q", err.Component, tt.wantComponent)
			}

			// Verify suggestions
			if len(err.Suggestions) != tt.wantSuggestions {
				t.Errorf("Suggestions count = %d, want %d", len(err.Suggestions), tt.wantSuggestions)
			}

			// Verify recoverable flag
			if err.Recoverable != tt.wantRecoverable {
				t.Errorf("Recoverable = %v, want %v", err.Recoverable, tt.wantRecoverable)
			}

			// Verify error message includes context
			errMsg := err.Error()
			if errMsg == "" {
				t.Error("Error message is empty")
			}

			// Verify suggestions are formatted
			suggestions := err.GetSuggestions()
			if tt.wantSuggestions > 0 && suggestions == "" {
				t.Error("GetSuggestions() returned empty string, want formatted suggestions")
			}
		})
	}
}

// TestErrorContextInRecovery demonstrates how error context is used in recovery decisions
func TestErrorContextInRecovery(t *testing.T) {
	tests := []struct {
		name           string
		err            *GenerationError
		ctx            *ErrorContext
		expectRetry    bool
		expectFallback bool
	}{
		{
			name: "first attempt tool execution - should retry",
			err:  NewToolExecutionError("npx", "nextjs", errors.New("failed")),
			ctx: &ErrorContext{
				Operation:     "bootstrap",
				Component:     "nextjs",
				Phase:         "generation",
				AttemptNumber: 1,
				CanRetry:      true,
				CanFallback:   true,
			},
			expectRetry:    true,
			expectFallback: false,
		},
		{
			name: "second attempt tool execution - should fallback",
			err:  NewToolExecutionError("npx", "nextjs", errors.New("failed")),
			ctx: &ErrorContext{
				Operation:     "bootstrap",
				Component:     "nextjs",
				Phase:         "generation",
				AttemptNumber: 2,
				CanRetry:      true,
				CanFallback:   true,
			},
			expectRetry:    false,
			expectFallback: true,
		},
		{
			name: "tool not found - should fallback immediately",
			err:  NewToolNotFoundError("gradle", "android"),
			ctx: &ErrorContext{
				Operation:     "bootstrap",
				Component:     "android",
				Phase:         "generation",
				AttemptNumber: 1,
				CanRetry:      true,
				CanFallback:   true,
			},
			expectRetry:    false,
			expectFallback: true,
		},
		{
			name: "config error - no retry or fallback",
			err:  NewConfigValidationError("name", "invalid characters"),
			ctx: &ErrorContext{
				Operation:     "validation",
				Component:     "",
				Phase:         "pre-generation",
				AttemptNumber: 1,
				CanRetry:      true,
				CanFallback:   true,
			},
			expectRetry:    false,
			expectFallback: false,
		},
		{
			name: "file system error - should retry",
			err:  NewFileSystemError("create", "/path/to/file", errors.New("permission denied")),
			ctx: &ErrorContext{
				Operation:     "file_creation",
				Component:     "backend",
				Phase:         "generation",
				AttemptNumber: 1,
				CanRetry:      true,
				CanFallback:   false,
			},
			expectRetry:    true,
			expectFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldRetry := ShouldRetry(tt.err, tt.ctx)
			if shouldRetry != tt.expectRetry {
				t.Errorf("ShouldRetry() = %v, want %v", shouldRetry, tt.expectRetry)
			}

			shouldFallback := ShouldFallback(tt.err, tt.ctx)
			if shouldFallback != tt.expectFallback {
				t.Errorf("ShouldFallback() = %v, want %v", shouldFallback, tt.expectFallback)
			}
		})
	}
}

// TestErrorFormattingWithContext demonstrates error formatting with full context
func TestErrorFormattingWithContext(t *testing.T) {
	tests := []struct {
		name     string
		err      *GenerationError
		contains []string
	}{
		{
			name: "tool not found error formatting",
			err:  NewToolNotFoundError("npx", "nextjs-frontend"),
			contains: []string{
				"TOOL_NOT_FOUND",
				"nextjs-frontend",
				"npx",
				"Suggestions:",
				"Install",
			},
		},
		{
			name: "tool execution error with cause",
			err:  NewToolExecutionError("go", "api-server", errors.New("exit status 1: module not found")),
			contains: []string{
				"TOOL_EXECUTION",
				"api-server",
				"go",
				"exit status 1",
				"Suggestions:",
			},
		},
		{
			name: "custom error with additional suggestions",
			err: NewGenerationError(
				ErrCategoryIntegration,
				"failed to connect components",
				errors.New("network timeout"),
			).WithSuggestions(
				"Check network connectivity",
				"Verify component endpoints",
				"Review firewall settings",
			),
			contains: []string{
				"INTEGRATION",
				"failed to connect components",
				"network timeout",
				"Suggestions:",
				"Check network connectivity",
				"Verify component endpoints",
				"Review firewall settings",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatError(tt.err)

			for _, substr := range tt.contains {
				if !containsSubstring(formatted, substr) {
					t.Errorf("Formatted error does not contain %q\nGot: %s", substr, formatted)
				}
			}
		})
	}
}

// TestRecoveryStrategyWithContext demonstrates recovery strategies with context
func TestRecoveryStrategyWithContext(t *testing.T) {
	tests := []struct {
		name            string
		err             *GenerationError
		wantDescription bool
		wantActions     bool
		minActionCount  int
	}{
		{
			name:            "tool not found strategy",
			err:             NewToolNotFoundError("npx", "nextjs"),
			wantDescription: true,
			wantActions:     true,
			minActionCount:  2,
		},
		{
			name:            "tool execution strategy",
			err:             NewToolExecutionError("go", "backend", errors.New("failed")),
			wantDescription: true,
			wantActions:     true,
			minActionCount:  2,
		},
		{
			name:            "file system strategy",
			err:             NewFileSystemError("write", "/path", errors.New("permission denied")),
			wantDescription: true,
			wantActions:     true,
			minActionCount:  2,
		},
		{
			name:            "timeout strategy",
			err:             NewTimeoutError("npm install", "frontend"),
			wantDescription: true,
			wantActions:     true,
			minActionCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetRecoveryStrategy(tt.err)

			if strategy == nil {
				t.Fatal("GetRecoveryStrategy() returned nil")
			}

			if tt.wantDescription && strategy.Description == "" {
				t.Error("Strategy description is empty")
			}

			if tt.wantActions && len(strategy.Actions) < tt.minActionCount {
				t.Errorf("Strategy has %d actions, want at least %d", len(strategy.Actions), tt.minActionCount)
			}

			// Verify each action has a description
			for i, action := range strategy.Actions {
				if action.Description == "" {
					t.Errorf("Action %d has empty description", i)
				}
			}
		})
	}
}

// TestAggregateErrorsWithContext demonstrates aggregating multiple errors with context
func TestAggregateErrorsWithContext(t *testing.T) {
	errors := []error{
		NewToolNotFoundError("npx", "frontend"),
		NewToolExecutionError("go", "backend", errors.New("build failed")),
		NewFileSystemError("create", "/output/dir", errors.New("permission denied")),
	}

	aggregated := AggregateErrors(errors)

	if aggregated == nil {
		t.Fatal("AggregateErrors() returned nil")
	}

	errMsg := aggregated.Error()

	// Verify all component contexts are included
	expectedSubstrings := []string{
		"multiple errors occurred",
		"frontend",
		"backend",
		"TOOL_NOT_FOUND",
		"TOOL_EXECUTION",
		"FILE_SYSTEM",
	}

	for _, substr := range expectedSubstrings {
		if !containsSubstring(errMsg, substr) {
			t.Errorf("Aggregated error does not contain %q\nGot: %s", substr, errMsg)
		}
	}
}

// TestErrorContextFields verifies all error context fields are properly set
func TestErrorContextFields(t *testing.T) {
	ctx := &ErrorContext{
		Operation:     "component_generation",
		Component:     "nextjs-app",
		Phase:         "bootstrap",
		AttemptNumber: 1,
		CanRetry:      true,
		CanFallback:   true,
	}

	// Verify all fields
	if ctx.Operation != "component_generation" {
		t.Errorf("Operation = %q, want %q", ctx.Operation, "component_generation")
	}

	if ctx.Component != "nextjs-app" {
		t.Errorf("Component = %q, want %q", ctx.Component, "nextjs-app")
	}

	if ctx.Phase != "bootstrap" {
		t.Errorf("Phase = %q, want %q", ctx.Phase, "bootstrap")
	}

	if ctx.AttemptNumber != 1 {
		t.Errorf("AttemptNumber = %d, want %d", ctx.AttemptNumber, 1)
	}

	if !ctx.CanRetry {
		t.Error("CanRetry = false, want true")
	}

	if !ctx.CanFallback {
		t.Error("CanFallback = false, want true")
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || containsSubstring(s[1:], substr))))
}

// TestErrorUnwrapping verifies error unwrapping works correctly
func TestErrorUnwrapping(t *testing.T) {
	cause := errors.New("underlying cause")
	err := NewToolExecutionError("npx", "nextjs", cause)

	unwrapped := errors.Unwrap(err)
	if !errors.Is(unwrapped, cause) {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}

	// Test with errors.Is
	if !errors.Is(err, cause) {
		t.Error("errors.Is() returned false, want true")
	}
}

// TestErrorCategorization verifies all error categories are properly categorized
func TestErrorCategorization(t *testing.T) {
	tests := []struct {
		name     string
		err      *GenerationError
		category ErrorCategory
	}{
		{
			name:     "tool not found",
			err:      NewToolNotFoundError("npx", "nextjs"),
			category: ErrCategoryToolNotFound,
		},
		{
			name:     "tool execution",
			err:      NewToolExecutionError("go", "backend", nil),
			category: ErrCategoryToolExecution,
		},
		{
			name:     "config validation",
			err:      NewConfigValidationError("name", "invalid"),
			category: ErrCategoryInvalidConfig,
		},
		{
			name:     "file system",
			err:      NewFileSystemError("create", "/path", nil),
			category: ErrCategoryFileSystem,
		},
		{
			name:     "security",
			err:      NewSecurityError("security violation", nil),
			category: ErrCategorySecurity,
		},
		{
			name:     "integration",
			err:      NewIntegrationError("integration failed", nil),
			category: ErrCategoryIntegration,
		},
		{
			name:     "validation",
			err:      NewValidationError("validation failed", nil),
			category: ErrCategoryValidation,
		},
		{
			name:     "timeout",
			err:      NewTimeoutError("operation", "component"),
			category: ErrCategoryTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Category != tt.category {
				t.Errorf("Category = %v, want %v", tt.err.Category, tt.category)
			}
		})
	}
}
