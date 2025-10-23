package cli

import (
	"errors"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// TestExitCodeHandler_DetermineExitCode_AllScenarios tests all exit code scenarios
func TestExitCodeHandler_DetermineExitCode_AllScenarios(t *testing.T) {
	log := logger.NewLogger()
	handler := NewExitCodeHandler(log)

	tests := []struct {
		name         string
		err          error
		expectedCode ExitCode
	}{
		// Success scenario
		{
			name:         "nil error returns success",
			err:          nil,
			expectedCode: ExitSuccess,
		},
		// User cancellation
		{
			name:         "user cancelled error",
			err:          ErrUserCancelled,
			expectedCode: ExitUserCancelled,
		},
		{
			name:         "wrapped user cancelled error",
			err:          errors.New("operation cancelled by user"),
			expectedCode: ExitGenerationFailed, // Plain error with "cancelled" text doesn't match ErrUserCancelled
		},
		// Configuration errors
		{
			name:         "validation failed error",
			err:          errors.New("validation failed: invalid field"),
			expectedCode: ExitConfigError,
		},
		{
			name:         "invalid configuration error",
			err:          errors.New("invalid configuration: missing required field"),
			expectedCode: ExitConfigError,
		},
		{
			name:         "config error message",
			err:          errors.New("config error: unable to parse"),
			expectedCode: ExitConfigError,
		},
		// Tool errors
		{
			name:         "tool not found error",
			err:          errors.New("tool not found: npx"),
			expectedCode: ExitToolsMissing,
		},
		{
			name:         "tool not available error",
			err:          errors.New("tool not available: go"),
			expectedCode: ExitToolsMissing,
		},
		{
			name:         "missing tool error",
			err:          errors.New("missing tool: gradle"),
			expectedCode: ExitToolsMissing,
		},
		{
			name:         "not whitelisted error",
			err:          errors.New("tool not whitelisted: unknown-tool"),
			expectedCode: ExitToolsMissing,
		},
		// Generation errors
		{
			name:         "generation failed error",
			err:          errors.New("generation failed: unable to create component"),
			expectedCode: ExitGenerationFailed,
		},
		{
			name:         "component generation error",
			err:          errors.New("component generation failed for nextjs"),
			expectedCode: ExitGenerationFailed,
		},
		{
			name:         "bootstrap failed error",
			err:          errors.New("bootstrap failed: command execution error"),
			expectedCode: ExitGenerationFailed,
		},
		{
			name:         "fallback failed error",
			err:          errors.New("fallback failed: template not found"),
			expectedCode: ExitGenerationFailed,
		},
		// File system errors
		{
			name:         "file system error",
			err:          errors.New("file system error: unable to write"),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "permission denied error",
			err:          errors.New("permission denied: cannot create directory"),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "no such file error",
			err:          errors.New("no such file or directory"),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "directory error",
			err:          errors.New("directory does not exist"),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "failed to create error",
			err:          errors.New("failed to create file"),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "failed to write error",
			err:          errors.New("failed to write to file"),
			expectedCode: ExitFileSystemError,
		},
		// Unknown errors default to generation failed
		{
			name:         "unknown error",
			err:          errors.New("something went wrong"),
			expectedCode: ExitGenerationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := handler.DetermineExitCode(tt.err)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

// TestExitCodeHandler_DetermineExitCode_CLIErrors tests CLIError handling
func TestExitCodeHandler_DetermineExitCode_CLIErrors(t *testing.T) {
	log := logger.NewLogger()
	handler := NewExitCodeHandler(log)

	tests := []struct {
		name         string
		cliErr       *CLIError
		expectedCode ExitCode
	}{
		{
			name:         "config error",
			cliErr:       NewConfigError("invalid config", nil),
			expectedCode: ExitConfigError,
		},
		{
			name:         "tool error",
			cliErr:       NewToolError("npx", "not found", nil),
			expectedCode: ExitToolsMissing,
		},
		{
			name:         "generation error",
			cliErr:       NewGenerationError("nextjs", "failed", nil),
			expectedCode: ExitGenerationFailed,
		},
		{
			name:         "filesystem error",
			cliErr:       NewFileSystemError("write", "/path/to/file", nil),
			expectedCode: ExitFileSystemError,
		},
		{
			name:         "user cancelled error",
			cliErr:       NewUserCancelledError("cancelled"),
			expectedCode: ExitUserCancelled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := handler.DetermineExitCode(tt.cliErr)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

// TestExitCodeHandler_DetermineExitCode_GenerationErrors tests GenerationError handling
func TestExitCodeHandler_DetermineExitCode_GenerationErrors(t *testing.T) {
	log := logger.NewLogger()
	handler := NewExitCodeHandler(log)

	tests := []struct {
		name         string
		genErr       *orchestrator.GenerationError
		expectedCode ExitCode
	}{
		{
			name: "invalid config category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryInvalidConfig,
				Message:  "invalid config",
			},
			expectedCode: ExitConfigError,
		},
		{
			name: "tool not found category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryToolNotFound,
				Message:  "tool not found",
			},
			expectedCode: ExitToolsMissing,
		},
		{
			name: "tool execution category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryToolExecution,
				Message:  "execution failed",
			},
			expectedCode: ExitGenerationFailed,
		},
		{
			name: "file system category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryFileSystem,
				Message:  "file system error",
			},
			expectedCode: ExitFileSystemError,
		},
		{
			name: "integration category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryIntegration,
				Message:  "integration failed",
			},
			expectedCode: ExitGenerationFailed,
		},
		{
			name: "validation category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryValidation,
				Message:  "validation failed",
			},
			expectedCode: ExitGenerationFailed,
		},
		{
			name: "timeout category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategoryTimeout,
				Message:  "operation timed out",
			},
			expectedCode: ExitGenerationFailed,
		},
		{
			name: "security category",
			genErr: &orchestrator.GenerationError{
				Category: orchestrator.ErrCategorySecurity,
				Message:  "security violation",
			},
			expectedCode: ExitConfigError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := handler.DetermineExitCode(tt.genErr)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

// TestExitCodeHandler_GetExitCodeDescription tests exit code descriptions
func TestExitCodeHandler_GetExitCodeDescription(t *testing.T) {
	log := logger.NewLogger()
	handler := NewExitCodeHandler(log)

	tests := []struct {
		code        ExitCode
		wantContain string
	}{
		{
			code:        ExitSuccess,
			wantContain: "successfully",
		},
		{
			code:        ExitConfigError,
			wantContain: "configuration",
		},
		{
			code:        ExitToolsMissing,
			wantContain: "tools",
		},
		{
			code:        ExitGenerationFailed,
			wantContain: "generate",
		},
		{
			code:        ExitFileSystemError,
			wantContain: "file system",
		},
		{
			code:        ExitUserCancelled,
			wantContain: "cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			description := handler.GetExitCodeDescription(tt.code)
			assert.NotEmpty(t, description)
			assert.Contains(t, description, tt.wantContain)
		})
	}
}

// TestCLIError_WithSuggestions tests adding suggestions to CLI errors
func TestCLIError_WithSuggestions(t *testing.T) {
	err := NewConfigError("invalid config", nil)

	// Initially should have default suggestions
	assert.NotEmpty(t, err.Suggestions)
	initialCount := len(err.Suggestions)

	// Add custom suggestions
	err = err.WithSuggestions("Custom suggestion 1", "Custom suggestion 2")

	// Should have more suggestions now
	assert.Len(t, err.Suggestions, initialCount+2)
	assert.Contains(t, err.Suggestions, "Custom suggestion 1")
	assert.Contains(t, err.Suggestions, "Custom suggestion 2")
}

// TestCLIError_WithExitCode tests setting exit code on CLI errors
func TestCLIError_WithExitCode(t *testing.T) {
	err := NewCLIError("test", "test error", nil)

	// Default exit code
	assert.Equal(t, ExitGenerationFailed, err.ExitCode)

	// Set custom exit code
	err = err.WithExitCode(ExitConfigError)
	assert.Equal(t, ExitConfigError, err.ExitCode)
}

// TestCLIError_GetSuggestions tests suggestion formatting
func TestCLIError_GetSuggestions(t *testing.T) {
	tests := []struct {
		name        string
		err         *CLIError
		wantEmpty   bool
		wantContain string
	}{
		{
			name:      "error with no suggestions",
			err:       &CLIError{Suggestions: []string{}},
			wantEmpty: true,
		},
		{
			name: "error with suggestions",
			err: &CLIError{
				Suggestions: []string{"Suggestion 1", "Suggestion 2"},
			},
			wantEmpty:   false,
			wantContain: "Next Steps:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := tt.err.GetSuggestions()

			if tt.wantEmpty {
				assert.Empty(t, suggestions)
			} else {
				assert.NotEmpty(t, suggestions)
				if tt.wantContain != "" {
					assert.Contains(t, suggestions, tt.wantContain)
				}
			}
		})
	}
}

// TestCLIError_Error tests error message formatting
func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         *CLIError
		wantContain []string
	}{
		{
			name: "error without cause",
			err: &CLIError{
				Category: "test",
				Message:  "test error",
				Cause:    nil,
			},
			wantContain: []string{"[test]", "test error"},
		},
		{
			name: "error with cause",
			err: &CLIError{
				Category: "test",
				Message:  "test error",
				Cause:    errors.New("underlying error"),
			},
			wantContain: []string{"[test]", "test error", "underlying error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, want := range tt.wantContain {
				assert.Contains(t, errMsg, want)
			}
		})
	}
}

// TestCLIError_Unwrap tests error unwrapping
func TestCLIError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &CLIError{
		Category: "test",
		Message:  "test error",
		Cause:    cause,
	}

	unwrapped := err.Unwrap()
	assert.Equal(t, cause, unwrapped)
}

// TestNewConfigError tests config error creation
func TestNewConfigError(t *testing.T) {
	err := NewConfigError("invalid field", nil)

	assert.Equal(t, "configuration", err.Category)
	assert.Equal(t, ExitConfigError, err.ExitCode)
	assert.NotEmpty(t, err.Suggestions)
	assert.Contains(t, err.Error(), "invalid field")
}

// TestNewToolError tests tool error creation
func TestNewToolError(t *testing.T) {
	err := NewToolError("npx", "not found", nil)

	assert.Equal(t, "tool", err.Category)
	assert.Equal(t, ExitToolsMissing, err.ExitCode)
	assert.NotEmpty(t, err.Suggestions)
	assert.Contains(t, err.Error(), "npx")
	assert.Contains(t, err.Error(), "not found")
}

// TestNewGenerationError tests generation error creation
func TestNewGenerationError(t *testing.T) {
	err := NewGenerationError("nextjs", "failed to generate", nil)

	assert.Equal(t, "generation", err.Category)
	assert.Equal(t, ExitGenerationFailed, err.ExitCode)
	assert.NotEmpty(t, err.Suggestions)
	assert.Contains(t, err.Error(), "nextjs")
	assert.Contains(t, err.Error(), "failed to generate")
}

// TestNewFileSystemError tests filesystem error creation
func TestNewFileSystemError(t *testing.T) {
	err := NewFileSystemError("write", "/path/to/file", nil)

	assert.Equal(t, "filesystem", err.Category)
	assert.Equal(t, ExitFileSystemError, err.ExitCode)
	assert.NotEmpty(t, err.Suggestions)
	assert.Contains(t, err.Error(), "write")
	assert.Contains(t, err.Error(), "/path/to/file")
}

// TestNewUserCancelledError tests user cancelled error creation
func TestNewUserCancelledError(t *testing.T) {
	err := NewUserCancelledError("operation cancelled")

	assert.Equal(t, "user", err.Category)
	assert.Equal(t, ExitUserCancelled, err.ExitCode)
	assert.Contains(t, err.Error(), "operation cancelled")
}

// TestExitCode_String tests exit code string representation
func (ec ExitCode) String() string {
	switch ec {
	case ExitSuccess:
		return "ExitSuccess"
	case ExitConfigError:
		return "ExitConfigError"
	case ExitToolsMissing:
		return "ExitToolsMissing"
	case ExitGenerationFailed:
		return "ExitGenerationFailed"
	case ExitFileSystemError:
		return "ExitFileSystemError"
	case ExitUserCancelled:
		return "ExitUserCancelled"
	default:
		return "Unknown"
	}
}

// TestExitCodeHandler_Integration tests complete exit code flow
func TestExitCodeHandler_Integration(t *testing.T) {
	log := logger.NewLogger()
	handler := NewExitCodeHandler(log)

	// Test complete flow: create error -> determine code -> get description
	tests := []struct {
		name         string
		createError  func() error
		expectedCode ExitCode
	}{
		{
			name: "config validation flow",
			createError: func() error {
				return NewConfigError("missing required field", nil).
					WithSuggestions("Add the missing field to your config")
			},
			expectedCode: ExitConfigError,
		},
		{
			name: "tool missing flow",
			createError: func() error {
				return NewToolError("npx", "not found in PATH", nil)
			},
			expectedCode: ExitToolsMissing,
		},
		{
			name: "generation failure flow",
			createError: func() error {
				return NewGenerationError("nextjs", "bootstrap command failed", errors.New("exit status 1"))
			},
			expectedCode: ExitGenerationFailed,
		},
		{
			name: "filesystem error flow",
			createError: func() error {
				return NewFileSystemError("create", "/protected/path", errors.New("permission denied"))
			},
			expectedCode: ExitFileSystemError,
		},
		{
			name: "user cancellation flow",
			createError: func() error {
				return NewUserCancelledError("user pressed Ctrl+C")
			},
			expectedCode: ExitUserCancelled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError()

			// Determine exit code
			code := handler.DetermineExitCode(err)
			assert.Equal(t, tt.expectedCode, code)

			// Get description
			description := handler.GetExitCodeDescription(code)
			assert.NotEmpty(t, description)

			// Verify error has suggestions if it's a CLIError
			var cliErr *CLIError
			if errors.As(err, &cliErr) {
				suggestions := cliErr.GetSuggestions()
				if len(cliErr.Suggestions) > 0 {
					assert.NotEmpty(t, suggestions)
					assert.Contains(t, suggestions, "Next Steps:")
				}
			}
		})
	}
}
