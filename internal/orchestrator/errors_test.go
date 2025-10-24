package orchestrator

import (
	"errors"
	"testing"
)

func TestGenerationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *GenerationError
		expected string
	}{
		{
			name: "basic error without component",
			err: &GenerationError{
				Category: ErrCategoryToolNotFound,
				Message:  "tool not found",
			},
			expected: "[TOOL_NOT_FOUND] tool not found",
		},
		{
			name: "error with component",
			err: &GenerationError{
				Category:  ErrCategoryToolExecution,
				Message:   "execution failed",
				Component: "nextjs",
			},
			expected: "[TOOL_EXECUTION] Component 'nextjs': execution failed",
		},
		{
			name: "error with cause",
			err: &GenerationError{
				Category: ErrCategoryFileSystem,
				Message:  "file operation failed",
				Cause:    errors.New("permission denied"),
			},
			expected: "[FILE_SYSTEM] file operation failed (caused by: permission denied)",
		},
		{
			name: "error with component and cause",
			err: &GenerationError{
				Category:  ErrCategoryIntegration,
				Message:   "integration failed",
				Component: "go-backend",
				Cause:     errors.New("connection refused"),
			},
			expected: "[INTEGRATION] Component 'go-backend': integration failed (caused by: connection refused)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGenerationError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &GenerationError{
		Category: ErrCategoryToolNotFound,
		Message:  "tool not found",
		Cause:    cause,
	}

	unwrapped := err.Unwrap()
	if !errors.Is(unwrapped, cause) {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestGenerationError_WithSuggestions(t *testing.T) {
	err := &GenerationError{
		Category:    ErrCategoryToolNotFound,
		Message:     "tool not found",
		Suggestions: []string{},
	}

	// Return value intentionally ignored - testing side effect on err.Suggestions
	_ = err.WithSuggestions("suggestion 1", "suggestion 2")

	if len(err.Suggestions) != 2 {
		t.Errorf("WithSuggestions() resulted in %d suggestions, want 2", len(err.Suggestions))
	}

	if err.Suggestions[0] != "suggestion 1" {
		t.Errorf("First suggestion = %q, want %q", err.Suggestions[0], "suggestion 1")
	}

	if err.Suggestions[1] != "suggestion 2" {
		t.Errorf("Second suggestion = %q, want %q", err.Suggestions[1], "suggestion 2")
	}
}

func TestGenerationError_GetSuggestions(t *testing.T) {
	tests := []struct {
		name        string
		suggestions []string
		wantEmpty   bool
	}{
		{
			name:        "no suggestions",
			suggestions: []string{},
			wantEmpty:   true,
		},
		{
			name:        "with suggestions",
			suggestions: []string{"suggestion 1", "suggestion 2"},
			wantEmpty:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &GenerationError{
				Category:    ErrCategoryToolNotFound,
				Message:     "test error",
				Suggestions: tt.suggestions,
			}

			result := err.GetSuggestions()

			if tt.wantEmpty {
				if result != "" {
					t.Errorf("GetSuggestions() = %q, want empty string", result)
				}
			} else {
				if result == "" {
					t.Error("GetSuggestions() returned empty string, want non-empty")
				}
			}
		})
	}
}

func TestNewGenerationError(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewGenerationError(ErrCategoryToolNotFound, "test message", cause)

	if err.Category != ErrCategoryToolNotFound {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryToolNotFound)
	}

	if err.Message != "test message" {
		t.Errorf("Message = %q, want %q", err.Message, "test message")
	}

	if !errors.Is(err.Cause, cause) {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true for TOOL_NOT_FOUND")
	}

	if err.Suggestions == nil {
		t.Error("Suggestions is nil, want empty slice")
	}
}

func TestNewToolNotFoundError(t *testing.T) {
	err := NewToolNotFoundError("npx", "nextjs")

	if err.Category != ErrCategoryToolNotFound {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryToolNotFound)
	}

	if err.Component != "nextjs" {
		t.Errorf("Component = %q, want %q", err.Component, "nextjs")
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewToolExecutionError(t *testing.T) {
	cause := errors.New("execution failed")
	err := NewToolExecutionError("npx", "nextjs", cause)

	if err.Category != ErrCategoryToolExecution {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryToolExecution)
	}

	if err.Component != "nextjs" {
		t.Errorf("Component = %q, want %q", err.Component, "nextjs")
	}

	if !errors.Is(err.Cause, cause) {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewConfigValidationError(t *testing.T) {
	err := NewConfigValidationError("name", "must not be empty")

	if err.Category != ErrCategoryInvalidConfig {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryInvalidConfig)
	}

	if err.Recoverable {
		t.Error("Recoverable = true, want false for config errors")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewFileSystemError(t *testing.T) {
	cause := errors.New("permission denied")
	err := NewFileSystemError("create", "/path/to/file", cause)

	if err.Category != ErrCategoryFileSystem {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryFileSystem)
	}

	if !errors.Is(err.Cause, cause) {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewSecurityError(t *testing.T) {
	cause := errors.New("path traversal detected")
	err := NewSecurityError("security violation", cause)

	if err.Category != ErrCategorySecurity {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategorySecurity)
	}

	if err.Recoverable {
		t.Error("Recoverable = true, want false for security errors")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewIntegrationError(t *testing.T) {
	cause := errors.New("connection failed")
	err := NewIntegrationError("integration failed", cause)

	if err.Category != ErrCategoryIntegration {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryIntegration)
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewValidationError(t *testing.T) {
	cause := errors.New("validation failed")
	err := NewValidationError("structure validation failed", cause)

	if err.Category != ErrCategoryValidation {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryValidation)
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestNewTimeoutError(t *testing.T) {
	err := NewTimeoutError("bootstrap execution", "nextjs")

	if err.Category != ErrCategoryTimeout {
		t.Errorf("Category = %v, want %v", err.Category, ErrCategoryTimeout)
	}

	if err.Component != "nextjs" {
		t.Errorf("Component = %q, want %q", err.Component, "nextjs")
	}

	if !err.Recoverable {
		t.Error("Recoverable = false, want true")
	}

	if len(err.Suggestions) == 0 {
		t.Error("Suggestions is empty, want at least one suggestion")
	}
}

func TestIsRecoverable(t *testing.T) {
	tests := []struct {
		name        string
		category    ErrorCategory
		recoverable bool
	}{
		{"tool not found", ErrCategoryToolNotFound, true},
		{"tool execution", ErrCategoryToolExecution, true},
		{"invalid config", ErrCategoryInvalidConfig, false},
		{"file system", ErrCategoryFileSystem, true},
		{"security", ErrCategorySecurity, false},
		{"integration", ErrCategoryIntegration, true},
		{"validation", ErrCategoryValidation, true},
		{"timeout", ErrCategoryTimeout, true},
		{"unknown", ErrCategoryUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRecoverable(tt.category)
			if result != tt.recoverable {
				t.Errorf("isRecoverable(%v) = %v, want %v", tt.category, result, tt.recoverable)
			}
		})
	}
}

func TestGetRecoveryStrategy(t *testing.T) {
	tests := []struct {
		name       string
		err        *GenerationError
		wantNil    bool
		minActions int
	}{
		{
			name: "tool not found",
			err: &GenerationError{
				Category: ErrCategoryToolNotFound,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "tool execution",
			err: &GenerationError{
				Category: ErrCategoryToolExecution,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "file system",
			err: &GenerationError{
				Category: ErrCategoryFileSystem,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "integration",
			err: &GenerationError{
				Category: ErrCategoryIntegration,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "validation",
			err: &GenerationError{
				Category: ErrCategoryValidation,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "timeout",
			err: &GenerationError{
				Category: ErrCategoryTimeout,
			},
			wantNil:    false,
			minActions: 1,
		},
		{
			name: "unknown category",
			err: &GenerationError{
				Category: ErrCategoryUnknown,
			},
			wantNil:    false,
			minActions: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetRecoveryStrategy(tt.err)

			if tt.wantNil {
				if strategy != nil {
					t.Error("GetRecoveryStrategy() returned non-nil, want nil")
				}
				return
			}

			if strategy == nil {
				t.Fatal("GetRecoveryStrategy() returned nil, want non-nil")
			}

			if strategy.Category != tt.err.Category {
				t.Errorf("Strategy.Category = %v, want %v", strategy.Category, tt.err.Category)
			}

			if len(strategy.Actions) < tt.minActions {
				t.Errorf("Strategy has %d actions, want at least %d", len(strategy.Actions), tt.minActions)
			}
		})
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name        string
		err         *GenerationError
		ctx         *ErrorContext
		shouldRetry bool
	}{
		{
			name: "retry tool execution on first attempt",
			err: &GenerationError{
				Category:    ErrCategoryToolExecution,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      true,
			},
			shouldRetry: true,
		},
		{
			name: "no retry after max attempts",
			err: &GenerationError{
				Category:    ErrCategoryToolExecution,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 2,
				CanRetry:      true,
			},
			shouldRetry: false,
		},
		{
			name: "no retry when context disallows",
			err: &GenerationError{
				Category:    ErrCategoryToolExecution,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      false,
			},
			shouldRetry: false,
		},
		{
			name: "no retry for non-recoverable errors",
			err: &GenerationError{
				Category:    ErrCategoryInvalidConfig,
				Recoverable: false,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      true,
			},
			shouldRetry: false,
		},
		{
			name: "retry file system errors",
			err: &GenerationError{
				Category:    ErrCategoryFileSystem,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      true,
			},
			shouldRetry: true,
		},
		{
			name: "retry timeout errors",
			err: &GenerationError{
				Category:    ErrCategoryTimeout,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      true,
			},
			shouldRetry: true,
		},
		{
			name: "no retry for integration errors",
			err: &GenerationError{
				Category:    ErrCategoryIntegration,
				Recoverable: true,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanRetry:      true,
			},
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldRetry(tt.err, tt.ctx)
			if result != tt.shouldRetry {
				t.Errorf("ShouldRetry() = %v, want %v", result, tt.shouldRetry)
			}
		})
	}
}

func TestShouldFallback(t *testing.T) {
	tests := []struct {
		name           string
		err            *GenerationError
		ctx            *ErrorContext
		shouldFallback bool
	}{
		{
			name: "fallback for tool not found",
			err: &GenerationError{
				Category: ErrCategoryToolNotFound,
			},
			ctx: &ErrorContext{
				CanFallback: true,
			},
			shouldFallback: true,
		},
		{
			name: "fallback for tool execution after retry",
			err: &GenerationError{
				Category: ErrCategoryToolExecution,
			},
			ctx: &ErrorContext{
				AttemptNumber: 2,
				CanFallback:   true,
			},
			shouldFallback: true,
		},
		{
			name: "no fallback for tool execution on first attempt",
			err: &GenerationError{
				Category: ErrCategoryToolExecution,
			},
			ctx: &ErrorContext{
				AttemptNumber: 1,
				CanFallback:   true,
			},
			shouldFallback: false,
		},
		{
			name: "no fallback when context disallows",
			err: &GenerationError{
				Category: ErrCategoryToolNotFound,
			},
			ctx: &ErrorContext{
				CanFallback: false,
			},
			shouldFallback: false,
		},
		{
			name: "no fallback for other error types",
			err: &GenerationError{
				Category: ErrCategoryFileSystem,
			},
			ctx: &ErrorContext{
				CanFallback: true,
			},
			shouldFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldFallback(tt.err, tt.ctx)
			if result != tt.shouldFallback {
				t.Errorf("ShouldFallback() = %v, want %v", result, tt.shouldFallback)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains []string
	}{
		{
			name: "generation error with suggestions",
			err: &GenerationError{
				Category:    ErrCategoryToolNotFound,
				Message:     "tool not found",
				Suggestions: []string{"Install the tool", "Use fallback"},
			},
			contains: []string{"TOOL_NOT_FOUND", "tool not found", "Suggestions:", "Install the tool", "Use fallback"},
		},
		{
			name: "generation error without suggestions",
			err: &GenerationError{
				Category: ErrCategoryFileSystem,
				Message:  "file operation failed",
			},
			contains: []string{"FILE_SYSTEM", "file operation failed"},
		},
		{
			name:     "regular error",
			err:      errors.New("regular error"),
			contains: []string{"regular error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatError(tt.err)

			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("FormatError() result does not contain %q\nGot: %s", substr, result)
				}
			}
		})
	}
}

func TestAggregateErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		wantNil  bool
		contains []string
	}{
		{
			name:    "no errors",
			errors:  []error{},
			wantNil: true,
		},
		{
			name:     "single error",
			errors:   []error{errors.New("error 1")},
			wantNil:  false,
			contains: []string{"error 1"},
		},
		{
			name: "multiple errors",
			errors: []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			},
			wantNil:  false,
			contains: []string{"multiple errors occurred", "error 1", "error 2", "error 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AggregateErrors(tt.errors)

			if tt.wantNil {
				if result != nil {
					t.Errorf("AggregateErrors() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Fatal("AggregateErrors() returned nil, want non-nil")
			}

			resultStr := result.Error()
			for _, substr := range tt.contains {
				if !contains(resultStr, substr) {
					t.Errorf("AggregateErrors() result does not contain %q\nGot: %s", substr, resultStr)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
