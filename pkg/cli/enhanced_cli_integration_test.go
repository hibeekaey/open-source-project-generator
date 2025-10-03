package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestEnhancedCLIBehaviorIntegration tests the complete enhanced CLI behavior
func TestEnhancedCLIBehaviorIntegration(t *testing.T) {
	t.Run("enhanced_flag_validation", func(t *testing.T) {
		testEnhancedFlagValidation(t)
	})

	t.Run("enhanced_error_handling", func(t *testing.T) {
		testEnhancedErrorHandling(t)
	})

	t.Run("graceful_fallback_handling", func(t *testing.T) {
		testGracefulFallbackHandling(t)
	})
}

func testEnhancedFlagValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupFlags    func(*cobra.Command)
		expectError   bool
		errorContains string
	}{
		{
			name: "no_conflicts",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("verbose", false, "")
				_ = cmd.Flags().Set("verbose", "true")
			},
			expectError: false,
		},
		{
			name: "verbose_quiet_conflict",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("verbose", false, "")
				cmd.Flags().Bool("quiet", false, "")
				// Actually set the flags to true to trigger conflict
				_ = cmd.Flags().Set("verbose", "true")
				_ = cmd.Flags().Set("quiet", "true")
			},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "debug_quiet_conflict",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("debug", false, "")
				cmd.Flags().Bool("quiet", false, "")
				// Actually set the flags to true to trigger conflict
				_ = cmd.Flags().Set("debug", "true")
				_ = cmd.Flags().Set("quiet", "true")
			},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "verbose_debug_no_conflict",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("verbose", false, "")
				cmd.Flags().Bool("debug", false, "")
				// Set both flags - these should not conflict
				_ = cmd.Flags().Set("verbose", "true")
				_ = cmd.Flags().Set("debug", "true")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock logger with all necessary expectations
			mockLogger := &MockLogger{}
			mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
			mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
			mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()
			mockLogger.On("InfoWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()
			mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Return().Maybe()
			mockLogger.On("DebugWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()

			// Create CLI components
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			outputManager := NewOutputManager(false, false, false, mockLogger)
			flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

			// Create test command
			cmd := &cobra.Command{Use: "test"}
			flagHandler.SetupGlobalFlags(cmd)
			tt.setupFlags(cmd)

			// Test flag validation
			err := flagHandler.HandleGlobalFlags(cmd)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func testEnhancedErrorHandling(t *testing.T) {
	// Test that enhanced error messages provide comprehensive information
	mockLogger := &MockLogger{}
	mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
	mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
	mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()
	mockLogger.On("InfoWithFields", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return().Maybe()
	mockLogger.On("Info", mock.AnythingOfType("string"), mock.Anything).Return().Maybe()

	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	cmd := &cobra.Command{Use: "test"}
	flagHandler.SetupGlobalFlags(cmd)

	// Set conflicting flags
	cmd.Flags().Bool("verbose", true, "")
	cmd.Flags().Bool("quiet", true, "")

	err := flagHandler.HandleGlobalFlags(cmd)
	assert.Error(t, err)

	errorMsg := err.Error()

	// Verify enhanced error message components
	expectedComponents := []string{
		"Flag conflicts detected",
		"Conflict",
		"Suggestion",
	}

	for _, component := range expectedComponents {
		assert.Contains(t, errorMsg, component, "Error message should contain %s", component)
	}

	// Verify error message is comprehensive (not just a generic error)
	assert.Greater(t, len(errorMsg), 50, "Error message should be detailed")
}

func testGracefulFallbackHandling(t *testing.T) {
	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectRecoverable   bool
	}{
		{
			name:              "force_interactive_overrides",
			interactive:       true,
			forceInteractive:  true,
			expectRecoverable: true,
		},
		{
			name:                "force_noninteractive_overrides",
			nonInteractive:      true,
			forceNonInteractive: true,
			expectRecoverable:   true,
		},
		{
			name:              "non_recoverable_basic_conflict",
			interactive:       true,
			nonInteractive:    true,
			expectRecoverable: false,
		},
		{
			name:              "non_recoverable_multiple_conflicts",
			interactive:       true,
			nonInteractive:    true,
			forceInteractive:  true,
			expectRecoverable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &FlagHandler{}

			isRecoverable := fh.isRecoverableConflict(
				tt.nonInteractive,
				tt.interactive,
				tt.forceInteractive,
				tt.forceNonInteractive,
				tt.explicitMode,
			)

			assert.Equal(t, tt.expectRecoverable, isRecoverable,
				"Recoverability check failed for %s", tt.name)
		})
	}
}

// TestEnhancedCLIWorkflow tests complete CLI workflows with enhanced behavior
func TestEnhancedCLIWorkflow(t *testing.T) {
	t.Run("global_flags_with_enhanced_detection", func(t *testing.T) {
		testGlobalFlagsWithEnhancedDetection(t)
	})
}

func testGlobalFlagsWithEnhancedDetection(t *testing.T) {
	// Test that global flag conflicts are detected with enhanced messages
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
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	cmd := &cobra.Command{Use: "test"}
	flagHandler.SetupGlobalFlags(cmd)

	// Test multiple conflict scenarios
	conflictScenarios := []struct {
		name  string
		flags map[string]string
	}{
		{
			name: "verbose_quiet",
			flags: map[string]string{
				"verbose": "true",
				"quiet":   "true",
			},
		},
		{
			name: "debug_quiet",
			flags: map[string]string{
				"debug": "true",
				"quiet": "true",
			},
		},
	}

	for _, scenario := range conflictScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Reset command flags
			testCmd := &cobra.Command{Use: "test"}
			flagHandler.SetupGlobalFlags(testCmd)

			// Set conflicting flags
			for flag, value := range scenario.flags {
				err := testCmd.PersistentFlags().Set(flag, value)
				assert.NoError(t, err, "Failed to set flag %s", flag)
			}

			err := flagHandler.HandleGlobalFlags(testCmd)
			assert.Error(t, err, "Should detect conflict for %s", scenario.name)

			// Verify enhanced error message structure
			errorMsg := err.Error()
			assert.Contains(t, errorMsg, "Flag conflicts detected")
			assert.Greater(t, len(errorMsg), 30, "Error message should be detailed")
		})
	}
}

// TestCLIIntegrationWithRealCommands tests CLI integration with actual command execution
func TestCLIIntegrationWithRealCommands(t *testing.T) {
	t.Run("help_command_integration", func(t *testing.T) {
		testHelpCommandIntegration(t)
	})

	t.Run("version_command_integration", func(t *testing.T) {
		testVersionCommandIntegration(t)
	})
}

func testHelpCommandIntegration(t *testing.T) {
	// Create a minimal CLI for testing help command
	mockLogger := &MockLogger{}
	mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
	mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
	mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()

	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
		rootCmd: &cobra.Command{
			Use:   "generator",
			Short: "Test CLI",
		},
	}

	// Test help command execution
	err := cli.Run([]string{"--help"})

	// Help command should not return an error (it exits with code 0)
	assert.NoError(t, err)
}

func testVersionCommandIntegration(t *testing.T) {
	// Test version command with enhanced CLI
	mockLogger := &MockLogger{}
	mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return().Maybe()
	mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return().Maybe()
	mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return().Maybe()

	cli := &CLI{
		outputManager:    NewOutputManager(false, false, false, mockLogger),
		generatorVersion: "test-version",
		gitCommit:        "test-commit",
		buildTime:        "test-build-time",
		rootCmd: &cobra.Command{
			Use:   "generator",
			Short: "Test CLI",
		},
	}

	// Initialize flag handler
	cli.flagHandler = NewFlagHandler(cli, cli.outputManager, mockLogger)

	// Add a simple version command for testing
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("Version: %s\n", cli.generatorVersion)
			return nil
		},
	}
	cli.rootCmd.AddCommand(versionCmd)

	// Test version command execution
	err := cli.Run([]string{"version"})
	assert.NoError(t, err)
}
