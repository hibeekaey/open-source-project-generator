package commands

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGenerateCLI is a mock implementation of GenerateCLI for testing.
type MockGenerateCLI struct {
	mock.Mock
}

func (m *MockGenerateCLI) ValidateGenerateOptions(options interfaces.GenerateOptions) error {
	args := m.Called(options)
	return args.Error(0)
}

func (m *MockGenerateCLI) DetectGenerationMode(configPath string, nonInteractive, interactive bool, explicitMode string) string {
	args := m.Called(configPath, nonInteractive, interactive, explicitMode)
	return args.String(0)
}

func (m *MockGenerateCLI) RouteToGenerationMethod(mode, configPath string, options interfaces.GenerateOptions) error {
	args := m.Called(mode, configPath, options)
	return args.Error(0)
}

func (m *MockGenerateCLI) ApplyModeOverrides(nonInteractive, interactive, forceInteractive, forceNonInteractive bool, explicitMode string) (bool, bool) {
	args := m.Called(nonInteractive, interactive, forceInteractive, forceNonInteractive, explicitMode)
	return args.Bool(0), args.Bool(1)
}

func (m *MockGenerateCLI) VerboseOutput(format string, args ...interface{}) {
	callArgs := []interface{}{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *MockGenerateCLI) DebugOutput(format string, args ...interface{}) {
	callArgs := []interface{}{format}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *MockGenerateCLI) Error(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func (m *MockGenerateCLI) Info(text string) string {
	args := m.Called(text)
	return args.String(0)
}

func TestNewGenerateCommand(t *testing.T) {
	mockCLI := &MockGenerateCLI{}
	cmd := NewGenerateCommand(mockCLI)

	assert.NotNil(t, cmd)
	assert.Equal(t, mockCLI, cmd.cli)
}

func TestGenerateCommand_Execute(t *testing.T) {
	tests := []struct {
		name          string
		setupFlags    func(*cobra.Command)
		setupMocks    func(*MockGenerateCLI)
		expectedError bool
		errorContains string
	}{
		{
			name: "successful execution with default flags",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().String("config", "", "")
				cmd.Flags().String("output", ".", "")
				cmd.Flags().Bool("dry-run", false, "")
				cmd.Flags().Bool("offline", false, "")
				cmd.Flags().Bool("minimal", false, "")
				cmd.Flags().String("template", "", "")
				cmd.Flags().Bool("update-versions", false, "")
				cmd.Flags().Bool("force", false, "")
				cmd.Flags().Bool("skip-validation", false, "")
				cmd.Flags().Bool("backup-existing", false, "")
				cmd.Flags().Bool("include-examples", false, "")
				cmd.Flags().Bool("non-interactive", false, "")
				cmd.Flags().StringSlice("exclude", []string{}, "")
				cmd.Flags().StringSlice("include-only", []string{}, "")
				cmd.Flags().Bool("interactive", false, "")
				cmd.Flags().String("preset", "", "")
				cmd.Flags().Bool("force-interactive", false, "")
				cmd.Flags().Bool("force-non-interactive", false, "")
				cmd.Flags().String("mode", "", "")
			},
			setupMocks: func(m *MockGenerateCLI) {
				m.On("ApplyModeOverrides", false, false, false, false, "").Return(false, false)
				m.On("VerboseOutput", "🔍 Validating your configuration...").Return()
				m.On("ValidateGenerateOptions", mock.AnythingOfType("interfaces.GenerateOptions")).Return(nil)
				m.On("DetectGenerationMode", "", false, false, "").Return("interactive")
				m.On("VerboseOutput", "🎯 Using %s mode for project generation", "interactive").Return()
				m.On("RouteToGenerationMethod", "interactive", "", mock.AnythingOfType("interfaces.GenerateOptions")).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "validation error",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().String("config", "", "")
				cmd.Flags().String("output", ".", "")
				cmd.Flags().Bool("dry-run", false, "")
				cmd.Flags().Bool("offline", false, "")
				cmd.Flags().Bool("minimal", false, "")
				cmd.Flags().String("template", "", "")
				cmd.Flags().Bool("update-versions", false, "")
				cmd.Flags().Bool("force", false, "")
				cmd.Flags().Bool("skip-validation", false, "")
				cmd.Flags().Bool("backup-existing", false, "")
				cmd.Flags().Bool("include-examples", false, "")
				cmd.Flags().Bool("non-interactive", false, "")
				cmd.Flags().StringSlice("exclude", []string{}, "")
				cmd.Flags().StringSlice("include-only", []string{}, "")
				cmd.Flags().Bool("interactive", false, "")
				cmd.Flags().String("preset", "", "")
				cmd.Flags().Bool("force-interactive", false, "")
				cmd.Flags().Bool("force-non-interactive", false, "")
				cmd.Flags().String("mode", "", "")
			},
			setupMocks: func(m *MockGenerateCLI) {
				m.On("ApplyModeOverrides", false, false, false, false, "").Return(false, false)
				m.On("VerboseOutput", "🔍 Validating your configuration...").Return()
				m.On("ValidateGenerateOptions", mock.AnythingOfType("interfaces.GenerateOptions")).Return(assert.AnError)
				m.On("Error", "Configuration validation failed.").Return("Configuration validation failed.")
				m.On("Info", "Please check your settings and try again").Return("Please check your settings and try again")
			},
			expectedError: true,
			errorContains: "Configuration validation failed",
		},
		{
			name: "conflicting mode flags",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().String("config", "", "")
				cmd.Flags().String("output", ".", "")
				cmd.Flags().Bool("dry-run", false, "")
				cmd.Flags().Bool("offline", false, "")
				cmd.Flags().Bool("minimal", false, "")
				cmd.Flags().String("template", "", "")
				cmd.Flags().Bool("update-versions", false, "")
				cmd.Flags().Bool("force", false, "")
				cmd.Flags().Bool("skip-validation", false, "")
				cmd.Flags().Bool("backup-existing", false, "")
				cmd.Flags().Bool("include-examples", false, "")
				cmd.Flags().Bool("non-interactive", true, "")
				cmd.Flags().StringSlice("exclude", []string{}, "")
				cmd.Flags().StringSlice("include-only", []string{}, "")
				cmd.Flags().Bool("interactive", true, "")
				cmd.Flags().String("preset", "", "")
				cmd.Flags().Bool("force-interactive", false, "")
				cmd.Flags().Bool("force-non-interactive", false, "")
				cmd.Flags().String("mode", "", "")
			},
			setupMocks: func(m *MockGenerateCLI) {
				m.On("Error", "Flag conflicts detected").Return("Flag conflicts detected")
				m.On("Info", "Conflict").Return("Conflict")
				m.On("Info", "#1").Return("#1")
				m.On("Info", "Conflicting flags").Return("Conflicting flags")
				m.On("Info", "Suggestion").Return("Suggestion")
				m.On("Info", "Examples").Return("Examples")
			},
			expectedError: true,
			errorContains: "Flag conflicts detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockGenerateCLI{}
			tt.setupMocks(mockCLI)

			cmd := &cobra.Command{}
			tt.setupFlags(cmd)

			generateCmd := NewGenerateCommand(mockCLI)
			err := generateCmd.Execute(cmd, []string{})

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockCLI.AssertExpectations(t)
		})
	}
}

func TestGenerateCommand_validateModeFlags(t *testing.T) {
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
			name:          "no flags set",
			expectedError: false,
		},
		{
			name:           "only non-interactive",
			nonInteractive: true,
			expectedError:  false,
		},
		{
			name:          "only interactive",
			interactive:   true,
			expectedError: false,
		},
		{
			name:          "only explicit mode",
			explicitMode:  "interactive",
			expectedError: false,
		},
		{
			name:           "conflicting non-interactive and interactive",
			nonInteractive: true,
			interactive:    true,
			expectedError:  true,
			errorContains:  "ERROR",
		},
		{
			name:                "conflicting force flags",
			forceInteractive:    true,
			forceNonInteractive: true,
			expectedError:       true,
			errorContains:       "ERROR",
		},
		{
			name:           "conflicting non-interactive and explicit mode",
			nonInteractive: true,
			explicitMode:   "interactive",
			expectedError:  true,
			errorContains:  "ERROR",
		},
		{
			name:          "invalid explicit mode",
			explicitMode:  "invalid-mode",
			expectedError: true,
			errorContains: "ERROR",
		},
		{
			name:                "interactive with force-non-interactive conflict",
			interactive:         true,
			forceNonInteractive: true,
			expectedError:       true,
			errorContains:       "ERROR",
		},
		{
			name:             "non-interactive with force-interactive conflict",
			nonInteractive:   true,
			forceInteractive: true,
			expectedError:    true,
			errorContains:    "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockGenerateCLI{}

			// Set up mock expectations for error cases
			if tt.expectedError {
				mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
				mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")
			}

			generateCmd := NewGenerateCommand(mockCLI)

			err := generateCmd.validateModeFlags(
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

			if tt.expectedError {
				mockCLI.AssertExpectations(t)
			}
		})
	}
}

func TestGenerateFlagConflictDetector_ValidateModeFlags(t *testing.T) {
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
			name:           "basic conflict",
			interactive:    true,
			nonInteractive: true,
			expectedError:  true,
			errorContains:  "ERROR",
		},
		{
			name:                "force flags conflict",
			forceInteractive:    true,
			forceNonInteractive: true,
			expectedError:       true,
			errorContains:       "ERROR",
		},
		{
			name:          "mode with interactive conflict",
			interactive:   true,
			explicitMode:  "non-interactive",
			expectedError: true,
			errorContains: "ERROR",
		},
		{
			name:          "valid mode variations",
			explicitMode:  "i",
			expectedError: false,
		},
		{
			name:          "valid mode variations - ni",
			explicitMode:  "ni",
			expectedError: false,
		},
		{
			name:          "valid mode variations - auto",
			explicitMode:  "auto",
			expectedError: false,
		},
		{
			name:          "invalid mode",
			explicitMode:  "unknown-mode",
			expectedError: true,
			errorContains: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockGenerateCLI{}

			if tt.expectedError {
				mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
				mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")
			}

			detector := NewGenerateFlagConflictDetector(mockCLI)

			err := detector.ValidateModeFlags(
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

			if tt.expectedError {
				mockCLI.AssertExpectations(t)
			}
		})
	}
}

func TestGenerateFlagConflictDetector_getGenerateConflictRules(t *testing.T) {
	mockCLI := &MockGenerateCLI{}
	detector := NewGenerateFlagConflictDetector(mockCLI)

	rules := detector.getGenerateConflictRules()

	assert.NotEmpty(t, rules)

	// Verify we have the expected conflict rules
	expectedConflicts := [][]string{
		{"--interactive", "--non-interactive"},
		{"--force-interactive", "--force-non-interactive"},
		{"--interactive", "--force-non-interactive"},
		{"--non-interactive", "--force-interactive"},
		{"--interactive", "--mode"},
		{"--non-interactive", "--mode"},
		{"--force-interactive", "--mode"},
		{"--force-non-interactive", "--mode"},
	}

	for _, expectedFlags := range expectedConflicts {
		found := false
		for _, rule := range rules {
			if containsAllFlags(rule.Flags, expectedFlags) {
				found = true
				// Verify rule has required fields
				assert.NotEmpty(t, rule.Description)
				assert.NotEmpty(t, rule.Suggestion)
				assert.NotEmpty(t, rule.Examples)
				assert.Equal(t, "error", rule.Severity)
				break
			}
		}
		assert.True(t, found, "Expected conflict rule for flags %v not found", expectedFlags)
	}
}

func TestGenerateFlagConflictDetector_validateExplicitMode(t *testing.T) {
	tests := []struct {
		name          string
		mode          string
		expectedError bool
		errorContains string
	}{
		{
			name:          "valid mode - interactive",
			mode:          "interactive",
			expectedError: false,
		},
		{
			name:          "valid mode - non-interactive",
			mode:          "non-interactive",
			expectedError: false,
		},
		{
			name:          "valid mode - config-file",
			mode:          "config-file",
			expectedError: false,
		},
		{
			name:          "valid variation - i",
			mode:          "i",
			expectedError: false,
		},
		{
			name:          "valid variation - ni",
			mode:          "ni",
			expectedError: false,
		},
		{
			name:          "valid variation - auto",
			mode:          "auto",
			expectedError: false,
		},
		{
			name:          "valid variation - config",
			mode:          "config",
			expectedError: false,
		},
		{
			name:          "valid variation - file",
			mode:          "file",
			expectedError: false,
		},
		{
			name:          "valid variation - cf",
			mode:          "cf",
			expectedError: false,
		},
		{
			name:          "case insensitive",
			mode:          "INTERACTIVE",
			expectedError: false,
		},
		{
			name:          "with whitespace",
			mode:          "  interactive  ",
			expectedError: false,
		},
		{
			name:          "invalid mode",
			mode:          "unknown",
			expectedError: true,
			errorContains: "ERROR",
		},
		{
			name:          "empty mode",
			mode:          "",
			expectedError: true,
			errorContains: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockGenerateCLI{}

			if tt.expectedError {
				mockCLI.On("Error", mock.AnythingOfType("string")).Return("ERROR")
				mockCLI.On("Info", mock.AnythingOfType("string")).Return("INFO")
			}

			detector := NewGenerateFlagConflictDetector(mockCLI)

			err := detector.validateExplicitMode(tt.mode)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedError {
				mockCLI.AssertExpectations(t)
			}
		})
	}
}

func TestGenerateCommand_createGenerateOptions(t *testing.T) {
	mockCLI := &MockGenerateCLI{}
	generateCmd := NewGenerateCommand(mockCLI)

	options := generateCmd.createGenerateOptions(
		true,            // force
		false,           // minimal
		true,            // offline
		false,           // updateVersions
		true,            // skipValidation
		false,           // backupExisting
		true,            // includeExamples
		"/test/path",    // outputPath
		false,           // dryRun
		true,            // nonInteractive
		"test-template", // template
	)

	assert.True(t, options.Force)
	assert.False(t, options.Minimal)
	assert.True(t, options.Offline)
	assert.False(t, options.UpdateVersions)
	assert.True(t, options.SkipValidation)
	assert.False(t, options.BackupExisting)
	assert.True(t, options.IncludeExamples)
	assert.Equal(t, "/test/path", options.OutputPath)
	assert.False(t, options.DryRun)
	assert.True(t, options.NonInteractive)
	assert.Equal(t, []string{"test-template"}, options.Templates)
}

func TestGenerateCommand_SetupFlags(t *testing.T) {
	mockCLI := &MockGenerateCLI{}
	generateCmd := NewGenerateCommand(mockCLI)

	cmd := &cobra.Command{}
	generateCmd.SetupFlags(cmd)

	// Test that all expected flags are present
	expectedFlags := []string{
		"config", "output", "dry-run", "offline", "minimal", "template",
		"update-versions", "force", "skip-validation", "backup-existing",
		"include-examples", "interactive", "force-interactive",
		"force-non-interactive", "mode", "exclude", "include-only", "preset",
	}

	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag %s should be present", flagName)
	}

	// Test flag shortcuts
	assert.NotNil(t, cmd.Flags().ShorthandLookup("c"), "Config flag should have shorthand 'c'")
	assert.NotNil(t, cmd.Flags().ShorthandLookup("o"), "Output flag should have shorthand 'o'")
	assert.NotNil(t, cmd.Flags().ShorthandLookup("n"), "Dry-run flag should have shorthand 'n'")
	assert.NotNil(t, cmd.Flags().ShorthandLookup("t"), "Template flag should have shorthand 't'")
	assert.NotNil(t, cmd.Flags().ShorthandLookup("f"), "Force flag should have shorthand 'f'")
}

// Helper function to check if a rule contains all expected flags
func containsAllFlags(ruleFlags, expectedFlags []string) bool {
	if len(ruleFlags) != len(expectedFlags) {
		return false
	}

	for _, expected := range expectedFlags {
		found := false
		for _, ruleFlag := range ruleFlags {
			if ruleFlag == expected {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
