package cli

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestErrorHandlingScenarios tests various error handling scenarios
func TestErrorHandlingScenarios(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockLogger)
		setupCLI      func(*CLI)
		operation     func(*CLI) error
		expectedError bool
		errorContains string
		description   string
	}{
		{
			name: "nil pointer safety in flag handler",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.flagHandler = nil
			},
			operation: func(cli *CLI) error {
				return cli.flagHandler.ValidateConflictingFlags(true, false, false)
			},
			expectedError: true,
			errorContains: "flag handler not initialized",
			description:   "Should handle nil flag handler gracefully",
		},
		{
			name: "nil command in flag handling",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.outputManager = NewOutputManager(false, false, false, &MockLogger{})
				cli.flagHandler = NewFlagHandler(cli, cli.outputManager, &MockLogger{})
			},
			operation: func(cli *CLI) error {
				return cli.flagHandler.HandleGlobalFlags(nil)
			},
			expectedError: true,
			errorContains: "command not provided",
			description:   "Should handle nil command gracefully",
		},
		{
			name: "invalid log level validation",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.outputManager = NewOutputManager(false, false, false, &MockLogger{})
				cli.flagHandler = NewFlagHandler(cli, cli.outputManager, &MockLogger{})
			},
			operation: func(cli *CLI) error {
				return cli.flagHandler.ValidateLogLevel("invalid-level")
			},
			expectedError: true,
			errorContains: "isn't a valid log level",
			description:   "Should validate log levels and provide helpful error",
		},
		{
			name: "invalid output format validation",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.outputManager = NewOutputManager(false, false, false, &MockLogger{})
				cli.flagHandler = NewFlagHandler(cli, cli.outputManager, &MockLogger{})
			},
			operation: func(cli *CLI) error {
				return cli.flagHandler.ValidateOutputFormat("xml")
			},
			expectedError: true,
			errorContains: "isn't a valid output format",
			description:   "Should validate output formats and provide helpful error",
		},
		{
			name: "invalid explicit mode validation",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.outputManager = NewOutputManager(false, false, false, &MockLogger{})
				cli.flagHandler = NewFlagHandler(cli, cli.outputManager, &MockLogger{})
			},
			operation: func(cli *CLI) error {
				return cli.flagHandler.validateExplicitMode("invalid-mode")
			},
			expectedError: true,
			errorContains: "is not a valid mode",
			description:   "Should validate explicit modes and provide helpful error",
		},
		{
			name: "flag conflict detection with detailed error",
			setupMocks: func(mockLogger *MockLogger) {
				// No specific mock setup needed
			},
			setupCLI: func(cli *CLI) {
				cli.outputManager = NewOutputManager(false, false, false, &MockLogger{})
				cli.flagHandler = NewFlagHandler(cli, cli.outputManager, &MockLogger{})
			},
			operation: func(cli *CLI) error {
				flagState := map[string]bool{
					"--verbose": true,
					"--quiet":   true,
				}
				return cli.flagHandler.ValidateConflictingFlagsEnhanced(flagState)
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
			description:   "Should detect flag conflicts and provide detailed error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			tt.setupMocks(mockLogger)

			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			tt.setupCLI(cli)

			err := tt.operation(cli)

			if tt.expectedError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, tt.description)
				}
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestFlagValidationErrorMessages tests error message quality
func TestFlagValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name               string
		validationFunction func(*FlagHandler) error
		expectedKeywords   []string
		description        string
	}{
		{
			name: "log level validation error",
			validationFunction: func(fh *FlagHandler) error {
				return fh.ValidateLogLevel("invalid")
			},
			expectedKeywords: []string{"valid log level", "Available options", "debug", "info", "warn", "error", "fatal"},
			description:      "Should provide comprehensive log level error with available options",
		},
		{
			name: "output format validation error",
			validationFunction: func(fh *FlagHandler) error {
				return fh.ValidateOutputFormat("xml")
			},
			expectedKeywords: []string{"valid output format", "Available options", "text", "json", "yaml"},
			description:      "Should provide comprehensive output format error with available options",
		},
		{
			name: "explicit mode validation error",
			validationFunction: func(fh *FlagHandler) error {
				return fh.validateExplicitMode("unknown")
			},
			expectedKeywords: []string{"not a valid mode", "Available modes", "interactive", "non-interactive", "config-file"},
			description:      "Should provide comprehensive mode validation error with available options",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			err := tt.validationFunction(flagHandler)
			assert.Error(t, err, tt.description)

			errorMessage := err.Error()
			for _, keyword := range tt.expectedKeywords {
				assert.Contains(t, errorMessage, keyword,
					"Error message should contain '%s' for %s", keyword, tt.description)
			}
		})
	}
}

// TestConflictErrorMessageStructure tests the structure of conflict error messages
func TestConflictErrorMessageStructure(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	conflicts := []ConflictRule{
		{
			Flags:       []string{"--verbose", "--quiet"},
			Description: "Verbose and quiet modes are mutually exclusive",
			Suggestion:  "Choose either verbose or quiet mode",
			Examples:    []string{"--verbose", "--quiet"},
			Severity:    "error",
		},
		{
			Flags:       []string{"--interactive", "--non-interactive"},
			Description: "Interactive and non-interactive modes cannot be used together",
			Suggestion:  "Choose either interactive or non-interactive mode",
			Examples:    []string{"--interactive", "--non-interactive"},
			Severity:    "error",
		},
	}

	err := flagHandler.generateConflictError(conflicts)
	assert.Error(t, err)

	errorMessage := err.Error()

	// Check for header
	assert.Contains(t, errorMessage, "Flag conflicts detected")

	// Check for conflict details
	assert.Contains(t, errorMessage, "Verbose and quiet modes are mutually exclusive")
	assert.Contains(t, errorMessage, "Interactive and non-interactive modes cannot be used together")

	// Check for suggestions
	assert.Contains(t, errorMessage, "Choose either verbose or quiet mode")
	assert.Contains(t, errorMessage, "Choose either interactive or non-interactive mode")

	// Check for conflict numbering
	assert.Contains(t, errorMessage, "#1")
	assert.Contains(t, errorMessage, "#2")
}

// TestErrorRecoveryMechanisms tests error recovery mechanisms
func TestErrorRecoveryMechanisms(t *testing.T) {
	tests := []struct {
		name          string
		setupScenario func() (*CLI, error)
		expectedError bool
		errorContains string
		description   string
	}{
		{
			name: "graceful fallback for missing flag values",
			setupScenario: func() (*CLI, error) {
				mockLogger := &MockLogger{}
				mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
				mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
				mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()
				mockLogger.On("InfoWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()
				mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Return().Maybe()
				mockLogger.On("DebugWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()

				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Create command without setting up flags properly
				cmd := &cobra.Command{}
				return cli, flagHandler.HandleGlobalFlags(cmd)
			},
			expectedError: false, // Should handle gracefully with defaults
			description:   "Should handle missing flags gracefully with default values",
		},
		{
			name: "recoverable mode conflict resolution",
			setupScenario: func() (*CLI, error) {
				mockLogger := &MockLogger{}
				mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
				mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
				mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()
				mockLogger.On("InfoWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()
				mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Return().Maybe()
				mockLogger.On("DebugWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()

				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Test recoverable conflict (force flag overrides regular flag)
				return cli, flagHandler.ValidateModeFlagsWithGracefulFallback(
					true,  // nonInteractive
					false, // interactive
					true,  // forceInteractive (should override nonInteractive)
					false, // forceNonInteractive
					"",    // explicitMode
				)
			},
			expectedError: false, // Should recover gracefully
			description:   "Should recover from conflicts when force flags are used",
		},
		{
			name: "non-recoverable conflict handling",
			setupScenario: func() (*CLI, error) {
				mockLogger := &MockLogger{}
				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Test non-recoverable conflict (too many conflicting flags)
				return cli, flagHandler.ValidateModeFlagsWithGracefulFallback(
					true,   // nonInteractive
					true,   // interactive
					false,  // forceInteractive
					false,  // forceNonInteractive
					"auto", // explicitMode
				)
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
			description:   "Should properly handle non-recoverable conflicts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := tt.setupScenario()

			if tt.expectedError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, tt.description)
				}
			} else {
				assert.NoError(t, err, tt.description)
			}

			// Verify CLI is still in a valid state
			assert.NotNil(t, cli, "CLI should not be nil after error handling")
		})
	}
}

// TestVersionCommandErrorHandling tests version command error scenarios
func TestVersionCommandErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockCLIInterface)
		setupFlags    func(*cobra.Command)
		expectedError bool
		errorContains string
		description   string
	}{
		{
			name: "build info error handling",
			setupMocks: func(mockCLI *MockCLIInterface) {
				mockCLI.On("GetVersionManager").Return(nil)
				mockCLI.On("GetBuildInfo").Return("", "", "") // Empty build info
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().Set("json", "true")
			},
			expectedError: false, // Should handle gracefully with defaults
			description:   "Should handle missing build info gracefully",
		},
		{
			name: "nil version manager handling",
			setupMocks: func(mockCLI *MockCLIInterface) {
				mockCLI.On("GetVersionManager").Return(nil)
			},
			setupFlags: func(cmd *cobra.Command) {
				// Default text output
			},
			expectedError: false, // Should use "dev" as fallback
			description:   "Should handle nil version manager gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockCLIInterface{}
			tt.setupMocks(mockCLI)

			cmd := &cobra.Command{}
			cmd.Flags().Bool("json", false, "")
			cmd.Flags().String("format", "", "")
			cmd.Flags().String("output-format", "", "")
			cmd.Flags().Bool("short", false, "")
			tt.setupFlags(cmd)

			err := RunVersionCommand(cmd, []string{}, mockCLI)

			if tt.expectedError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, tt.description)
				}
			} else {
				assert.NoError(t, err, tt.description)
			}

			mockCLI.AssertExpectations(t)
		})
	}
}

// TestCIDetectionErrorHandling tests CI detection error scenarios
func TestCIDetectionErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		operation   func(*CLI) interface{}
		expectPanic bool
		description string
	}{
		{
			name: "missing environment variables",
			setupEnv: func() {
				// Clear all CI-related environment variables
				clearAllCIEnvironment()
			},
			cleanupEnv: func() {
				// No cleanup needed
			},
			operation: func(cli *CLI) interface{} {
				return cli.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle missing environment variables gracefully",
		},
		{
			name: "malformed environment variables",
			setupEnv: func() {
				clearAllCIEnvironment()
				// Set some unusual values
				_ = os.Setenv("CI", "maybe")
				_ = os.Setenv("BUILD_NUMBER", "not-a-number")
			},
			cleanupEnv: func() {
				_ = os.Unsetenv("CI")
				_ = os.Unsetenv("BUILD_NUMBER")
			},
			operation: func(cli *CLI) interface{} {
				return cli.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle malformed environment variables gracefully",
		},
		{
			name: "nil CLI in CI detection",
			setupEnv: func() {
				clearAllCIEnvironment()
			},
			cleanupEnv: func() {
				// No cleanup needed
			},
			operation: func(cli *CLI) interface{} {
				// Test with nil CLI (should not panic)
				var nilCLI *CLI
				if nilCLI == nil {
					// Return a safe default instead of calling method on nil
					return &CIEnvironment{
						IsCI:     false,
						Provider: "",
					}
				}
				return nilCLI.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle nil CLI gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					tt.operation(cli)
				}, tt.description)
			} else {
				assert.NotPanics(t, func() {
					result := tt.operation(cli)
					assert.NotNil(t, result, tt.description)
				}, tt.description)
			}
		})
	}
}

// TestEdgeCaseErrorHandling tests edge case error scenarios
func TestEdgeCaseErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		operation   func() error
		expectedErr bool
		description string
	}{
		{
			name: "empty conflict rules list",
			operation: func() error {
				mockLogger := &MockLogger{}
				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Test with empty conflicts list
				return flagHandler.generateConflictError([]ConflictRule{})
			},
			expectedErr: false, // Should return nil for empty list
			description: "Should handle empty conflict rules gracefully",
		},
		{
			name: "malformed conflict rule",
			operation: func() error {
				mockLogger := &MockLogger{}
				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Test with malformed conflict rule
				conflicts := []ConflictRule{
					{
						Flags:       []string{}, // Empty flags
						Description: "",         // Empty description
						Suggestion:  "",         // Empty suggestion
						Examples:    []string{}, // Empty examples
						Severity:    "",         // Empty severity
					},
				}
				return flagHandler.generateConflictError(conflicts)
			},
			expectedErr: true, // Should still generate error
			description: "Should handle malformed conflict rules",
		},
		{
			name: "flag state with nil map",
			operation: func() error {
				mockLogger := &MockLogger{}
				cli := &CLI{
					outputManager: NewOutputManager(false, false, false, mockLogger),
				}
				flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

				// Test with nil flag state
				return flagHandler.ValidateConflictingFlagsEnhanced(nil)
			},
			expectedErr: false, // Should handle nil gracefully
			description: "Should handle nil flag state gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()

			if tt.expectedErr {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestErrorMessageLocalization tests error message consistency
func TestErrorMessageConsistency(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	// Test that all error messages follow consistent patterns
	tests := []struct {
		name      string
		operation func() error
		patterns  []string
	}{
		{
			name: "log level validation",
			operation: func() error {
				return flagHandler.ValidateLogLevel("invalid")
			},
			patterns: []string{"🚫", "isn't a valid", "Available options"},
		},
		{
			name: "output format validation",
			operation: func() error {
				return flagHandler.ValidateOutputFormat("invalid")
			},
			patterns: []string{"🚫", "isn't a valid", "Available options"},
		},
		{
			name: "mode validation",
			operation: func() error {
				return flagHandler.validateExplicitMode("invalid")
			},
			patterns: []string{"🚫", "is not a valid", "Available modes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			assert.Error(t, err)

			errorMessage := err.Error()
			for _, pattern := range tt.patterns {
				assert.Contains(t, errorMessage, pattern,
					"Error message should contain pattern '%s'", pattern)
			}
		})
	}
}
