package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestEnhancedFlagValidation tests the enhanced flag validation logic
func TestEnhancedFlagValidation(t *testing.T) {
	tests := []struct {
		name          string
		flagState     map[string]bool
		expectedError bool
		errorContains string
	}{
		{
			name: "no conflicts",
			flagState: map[string]bool{
				"--verbose": true,
				"--debug":   false,
				"--quiet":   false,
			},
			expectedError: false,
		},
		{
			name: "verbose and quiet conflict",
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   true,
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "debug and quiet conflict",
			flagState: map[string]bool{
				"--debug": true,
				"--quiet": true,
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "interactive and non-interactive conflict",
			flagState: map[string]bool{
				"--interactive":     true,
				"--non-interactive": true,
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "force flags conflict",
			flagState: map[string]bool{
				"--force-interactive":     true,
				"--force-non-interactive": true,
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
		{
			name: "mode flag with interactive conflict",
			flagState: map[string]bool{
				"--interactive": true,
				"--mode":        true,
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			err := flagHandler.ValidateConflictingFlagsEnhanced(tt.flagState)

			if tt.expectedError {
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

// TestFlagConflictMatrix tests the conflict matrix functionality
func TestFlagConflictMatrix(t *testing.T) {
	matrix := NewFlagConflictMatrix()

	assert.NotNil(t, matrix)
	assert.NotEmpty(t, matrix.Rules)

	// Verify all rules have required fields
	for i, rule := range matrix.Rules {
		assert.NotEmpty(t, rule.Flags, "Rule %d should have flags", i)
		assert.NotEmpty(t, rule.Description, "Rule %d should have description", i)
		assert.NotEmpty(t, rule.Suggestion, "Rule %d should have suggestion", i)
		assert.NotEmpty(t, rule.Severity, "Rule %d should have severity", i)
		assert.Contains(t, []string{"error", "warning", "info"}, rule.Severity, "Rule %d should have valid severity", i)
	}
}

// TestCheckConflictRule tests individual conflict rule checking
func TestCheckConflictRule(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	tests := []struct {
		name      string
		rule      ConflictRule
		flagState map[string]bool
		expected  bool
	}{
		{
			name: "conflict detected",
			rule: ConflictRule{
				Flags: []string{"--verbose", "--quiet"},
			},
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   true,
			},
			expected: true,
		},
		{
			name: "no conflict - only one flag active",
			rule: ConflictRule{
				Flags: []string{"--verbose", "--quiet"},
			},
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   false,
			},
			expected: false,
		},
		{
			name: "no conflict - no flags active",
			rule: ConflictRule{
				Flags: []string{"--verbose", "--quiet"},
			},
			flagState: map[string]bool{
				"--verbose": false,
				"--quiet":   false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flagHandler.checkConflictRule(tt.rule, tt.flagState)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsFlagActive tests flag activity detection
func TestIsFlagActive(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	tests := []struct {
		name      string
		flag      string
		flagState map[string]bool
		expected  bool
	}{
		{
			name: "direct flag active",
			flag: "--verbose",
			flagState: map[string]bool{
				"--verbose": true,
			},
			expected: true,
		},
		{
			name: "direct flag inactive",
			flag: "--verbose",
			flagState: map[string]bool{
				"--verbose": false,
			},
			expected: false,
		},
		{
			name: "mode flag with value",
			flag: "--mode=interactive",
			flagState: map[string]bool{
				"--mode": true,
			},
			expected: true,
		},
		{
			name: "output format flag with value",
			flag: "--output-format=json",
			flagState: map[string]bool{
				"--output-format": true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flagHandler.isFlagActive(tt.flag, tt.flagState)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGenerateConflictError tests conflict error message generation
func TestGenerateConflictError(t *testing.T) {
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
	}

	err := flagHandler.generateConflictError(conflicts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Flag conflicts detected")
	assert.Contains(t, err.Error(), "Verbose and quiet modes are mutually exclusive")
}

// TestCollectFlagState tests flag state collection
func TestCollectFlagState(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	cmd := &cobra.Command{}
	cmd.Flags().Bool("verbose", true, "")
	cmd.Flags().Bool("quiet", false, "")
	cmd.Flags().Bool("debug", true, "")
	cmd.Flags().String("mode", "interactive", "")

	flagState := flagHandler.collectFlagState(cmd)

	assert.True(t, flagState["--verbose"])
	assert.False(t, flagState["--quiet"])
	assert.True(t, flagState["--debug"])
	assert.True(t, flagState["--mode"])
	assert.True(t, flagState["--mode=interactive"])
}

// TestValidateAllFlags tests comprehensive flag validation
func TestValidateAllFlags(t *testing.T) {
	tests := []struct {
		name          string
		setupFlags    func(*cobra.Command)
		expectedError bool
		errorContains string
	}{
		{
			name: "no conflicts",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("verbose", true, "")
				cmd.Flags().Bool("quiet", false, "")
			},
			expectedError: false,
		},
		{
			name: "verbose and quiet conflict",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("verbose", true, "")
				cmd.Flags().Bool("quiet", true, "")
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			cmd := &cobra.Command{}
			tt.setupFlags(cmd)

			err := flagHandler.ValidateAllFlags(cmd)

			if tt.expectedError {
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

// TestValidateModeFlagsWithGracefulFallback tests graceful fallback handling
func TestValidateModeFlagsWithGracefulFallback(t *testing.T) {
	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectedError       bool
		errorContains       string
	}{
		{
			name:          "no conflicts",
			interactive:   true,
			expectedError: false,
		},
		{
			name:             "recoverable conflict - force interactive wins",
			nonInteractive:   true,
			forceInteractive: true,
			expectedError:    false,
		},
		{
			name:                "recoverable conflict - force non-interactive wins",
			interactive:         true,
			forceNonInteractive: true,
			expectedError:       false,
		},
		{
			name:           "non-recoverable conflict",
			interactive:    true,
			nonInteractive: true,
			explicitMode:   "config",
			expectedError:  true,
			errorContains:  "Flag conflicts detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			err := flagHandler.ValidateModeFlagsWithGracefulFallback(
				tt.nonInteractive,
				tt.interactive,
				tt.forceInteractive,
				tt.forceNonInteractive,
				tt.explicitMode,
			)

			if tt.expectedError {
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

// TestIsRecoverableConflict tests recoverable conflict detection
func TestIsRecoverableConflict(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expected            bool
	}{
		{
			name:             "recoverable - force interactive with non-interactive",
			nonInteractive:   true,
			forceInteractive: true,
			expected:         true,
		},
		{
			name:                "recoverable - force non-interactive with interactive",
			interactive:         true,
			forceNonInteractive: true,
			expected:            true,
		},
		{
			name:           "not recoverable - no force flags",
			interactive:    true,
			nonInteractive: true,
			expected:       false,
		},
		{
			name:           "not recoverable - too many flags",
			interactive:    true,
			nonInteractive: true,
			explicitMode:   "config",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flagHandler.isRecoverableConflict(
				tt.nonInteractive,
				tt.interactive,
				tt.forceInteractive,
				tt.forceNonInteractive,
				tt.explicitMode,
			)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHandleRecoverableConflict tests recoverable conflict handling
func TestHandleRecoverableConflict(t *testing.T) {
	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectedError       bool
	}{
		{
			name:             "resolve force interactive conflict",
			nonInteractive:   true,
			forceInteractive: true,
			expectedError:    false,
		},
		{
			name:                "resolve force non-interactive conflict",
			interactive:         true,
			forceNonInteractive: true,
			expectedError:       false,
		},
		{
			name:          "no conflict to resolve",
			interactive:   true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			err := flagHandler.handleRecoverableConflict(
				tt.nonInteractive,
				tt.interactive,
				tt.forceInteractive,
				tt.forceNonInteractive,
				tt.explicitMode,
			)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetModeFromFlags tests mode extraction from flags
func TestGetModeFromFlags(t *testing.T) {
	tests := []struct {
		name          string
		setupFlags    func(*cobra.Command)
		expectedMode  GenerationMode
		expectedError bool
		errorContains string
	}{
		{
			name: "interactive mode",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("interactive", true, "")
				cmd.Flags().Bool("non-interactive", false, "")
			},
			expectedMode:  ModeInteractive,
			expectedError: false,
		},
		{
			name: "non-interactive mode",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("interactive", false, "")
				cmd.Flags().Bool("non-interactive", true, "")
			},
			expectedMode:  ModeNonInteractive,
			expectedError: false,
		},
		{
			name: "explicit mode",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().String("mode", "config-file", "")
			},
			expectedMode:  ModeConfig,
			expectedError: false,
		},
		{
			name: "conflicting flags",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Bool("interactive", true, "")
				cmd.Flags().Bool("non-interactive", true, "")
			},
			expectedMode:  ModeAuto,
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}
			flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

			cmd := &cobra.Command{}
			tt.setupFlags(cmd)

			mode, err := flagHandler.GetModeFromFlags(cmd)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMode, mode)
			}
		})
	}
}

// TestParseExplicitMode tests explicit mode parsing
func TestParseExplicitMode(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

	tests := []struct {
		name     string
		mode     string
		expected GenerationMode
	}{
		{
			name:     "interactive",
			mode:     "interactive",
			expected: ModeInteractive,
		},
		{
			name:     "interactive short",
			mode:     "i",
			expected: ModeInteractive,
		},
		{
			name:     "non-interactive",
			mode:     "non-interactive",
			expected: ModeNonInteractive,
		},
		{
			name:     "non-interactive short",
			mode:     "ni",
			expected: ModeNonInteractive,
		},
		{
			name:     "config-file",
			mode:     "config-file",
			expected: ModeConfig,
		},
		{
			name:     "config short",
			mode:     "cf",
			expected: ModeConfig,
		},
		{
			name:     "auto mode",
			mode:     "auto",
			expected: ModeNonInteractive,
		},
		{
			name:     "case insensitive",
			mode:     "INTERACTIVE",
			expected: ModeInteractive,
		},
		{
			name:     "with whitespace",
			mode:     "  interactive  ",
			expected: ModeInteractive,
		},
		{
			name:     "unknown mode",
			mode:     "unknown",
			expected: ModeAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flagHandler.parseExplicitMode(tt.mode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNullSafetyChecks tests null safety in flag handling
func TestNullSafetyChecks(t *testing.T) {
	t.Run("nil flag handler", func(t *testing.T) {
		var flagHandler *FlagHandler
		err := flagHandler.ValidateConflictingFlags(true, false, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "flag handler not initialized")
	})

	t.Run("nil command", func(t *testing.T) {
		mockLogger := &MockLogger{}
		cli := &CLI{
			outputManager: NewOutputManager(false, false, false, mockLogger),
		}
		flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

		err := flagHandler.HandleGlobalFlags(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "command not provided")
	})

	t.Run("nil command in GetModeFromFlags", func(t *testing.T) {
		mockLogger := &MockLogger{}
		cli := &CLI{
			outputManager: NewOutputManager(false, false, false, mockLogger),
		}
		flagHandler := NewFlagHandler(cli, cli.outputManager, mockLogger)

		mode, err := flagHandler.GetModeFromFlags(nil)
		assert.Error(t, err)
		assert.Equal(t, ModeAuto, mode)
		assert.Contains(t, err.Error(), "command not provided")
	})
}
