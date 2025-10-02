package cli

import (
	"os"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger is a mock implementation of the Logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Fatal(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *MockLogger) ErrorWithError(msg string, err error, fields map[string]interface{}) {
	m.Called(msg, err, fields)
}

func (m *MockLogger) StartOperation(operation string, fields map[string]interface{}) *interfaces.OperationContext {
	args := m.Called(operation, fields)
	return args.Get(0).(*interfaces.OperationContext)
}

func (m *MockLogger) LogOperationStart(operation string, fields map[string]interface{}) {
	m.Called(operation, fields)
}

func (m *MockLogger) LogOperationSuccess(operation string, duration time.Duration, fields map[string]interface{}) {
	m.Called(operation, duration, fields)
}

func (m *MockLogger) LogOperationError(operation string, err error, fields map[string]interface{}) {
	m.Called(operation, err, fields)
}

func (m *MockLogger) FinishOperation(ctx *interfaces.OperationContext, additionalFields map[string]interface{}) {
	m.Called(ctx, additionalFields)
}

func (m *MockLogger) FinishOperationWithError(ctx *interfaces.OperationContext, err error, additionalFields map[string]interface{}) {
	m.Called(ctx, err, additionalFields)
}

func (m *MockLogger) LogPerformanceMetrics(operation string, metrics map[string]interface{}) {
	m.Called(operation, metrics)
}

func (m *MockLogger) LogMemoryUsage(operation string) {
	m.Called(operation)
}

func (m *MockLogger) SetLevel(level int) {
	m.Called(level)
}

func (m *MockLogger) GetLevel() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockLogger) SetJSONOutput(enabled bool) {
	m.Called(enabled)
}

func (m *MockLogger) SetCallerInfo(enabled bool) {
	m.Called(enabled)
}

func (m *MockLogger) IsDebugEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) IsInfoEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLogger) WithComponent(component string) interfaces.Logger {
	args := m.Called(component)
	return args.Get(0).(interfaces.Logger)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.LoggerContext {
	args := m.Called(fields)
	return args.Get(0).(interfaces.LoggerContext)
}

func (m *MockLogger) GetLogDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLogger) GetRecentEntries(limit int) []interfaces.LogEntry {
	args := m.Called(limit)
	return args.Get(0).([]interfaces.LogEntry)
}

func (m *MockLogger) FilterEntries(level string, component string, since time.Time, limit int) []interfaces.LogEntry {
	args := m.Called(level, component, since, limit)
	return args.Get(0).([]interfaces.LogEntry)
}

func (m *MockLogger) GetLogFiles() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockLogger) ReadLogFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockLogger) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewFlagHandler(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{}
	outputManager := NewOutputManager(false, false, false, mockLogger)

	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	assert.NotNil(t, flagHandler)
	assert.Equal(t, cli, flagHandler.cli)
	assert.Equal(t, outputManager, flagHandler.outputManager)
	assert.Equal(t, mockLogger, flagHandler.logger)
}

func TestSetupGlobalFlags(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	rootCmd := &cobra.Command{
		Use: "test",
	}

	flagHandler.SetupGlobalFlags(rootCmd)

	// Check that all expected flags are added
	flags := []string{"verbose", "quiet", "debug", "log-level", "log-json", "log-caller", "non-interactive", "output-format"}
	for _, flagName := range flags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag %s should be added", flagName)
	}

	// Check flag defaults
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	assert.Equal(t, "false", verboseFlag.DefValue)

	logLevelFlag := rootCmd.PersistentFlags().Lookup("log-level")
	assert.Equal(t, "info", logLevelFlag.DefValue)

	outputFormatFlag := rootCmd.PersistentFlags().Lookup("output-format")
	assert.Equal(t, "text", outputFormatFlag.DefValue)
}

func TestValidateConflictingFlags(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	tests := []struct {
		name     string
		verbose  bool
		quiet    bool
		debug    bool
		hasError bool
	}{
		{
			name:     "no conflicts",
			verbose:  false,
			quiet:    false,
			debug:    false,
			hasError: false,
		},
		{
			name:     "verbose only",
			verbose:  true,
			quiet:    false,
			debug:    false,
			hasError: false,
		},
		{
			name:     "debug only",
			verbose:  false,
			quiet:    false,
			debug:    true,
			hasError: false,
		},
		{
			name:     "verbose and quiet conflict",
			verbose:  true,
			quiet:    true,
			debug:    false,
			hasError: true,
		},
		{
			name:     "debug and quiet conflict",
			verbose:  false,
			quiet:    true,
			debug:    true,
			hasError: true,
		},
		{
			name:     "verbose and debug together (no conflict)",
			verbose:  true,
			quiet:    false,
			debug:    true,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flagHandler.ValidateConflictingFlags(tt.verbose, tt.quiet, tt.debug)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLogLevel(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	tests := []struct {
		name     string
		logLevel string
		hasError bool
	}{
		{
			name:     "valid debug level",
			logLevel: "debug",
			hasError: false,
		},
		{
			name:     "valid info level",
			logLevel: "info",
			hasError: false,
		},
		{
			name:     "valid warn level",
			logLevel: "warn",
			hasError: false,
		},
		{
			name:     "valid error level",
			logLevel: "error",
			hasError: false,
		},
		{
			name:     "valid fatal level",
			logLevel: "fatal",
			hasError: false,
		},
		{
			name:     "invalid level",
			logLevel: "invalid",
			hasError: true,
		},
		{
			name:     "empty level",
			logLevel: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flagHandler.ValidateLogLevel(tt.logLevel)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputFormat(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	tests := []struct {
		name         string
		outputFormat string
		hasError     bool
	}{
		{
			name:         "valid text format",
			outputFormat: "text",
			hasError:     false,
		},
		{
			name:         "valid json format",
			outputFormat: "json",
			hasError:     false,
		},
		{
			name:         "valid yaml format",
			outputFormat: "yaml",
			hasError:     false,
		},
		{
			name:         "invalid format",
			outputFormat: "xml",
			hasError:     true,
		},
		{
			name:         "empty format",
			outputFormat: "",
			hasError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flagHandler.ValidateOutputFormat(tt.outputFormat)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFlagHandlerValidateModeFlags(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		hasError            bool
	}{
		{
			name:                "no flags set",
			nonInteractive:      false,
			interactive:         false,
			forceInteractive:    false,
			forceNonInteractive: false,
			explicitMode:        "",
			hasError:            false,
		},
		{
			name:                "only non-interactive",
			nonInteractive:      true,
			interactive:         false,
			forceInteractive:    false,
			forceNonInteractive: false,
			explicitMode:        "",
			hasError:            false,
		},
		{
			name:                "only force-interactive",
			nonInteractive:      false,
			interactive:         false,
			forceInteractive:    true,
			forceNonInteractive: false,
			explicitMode:        "",
			hasError:            false,
		},
		{
			name:                "conflicting flags",
			nonInteractive:      true,
			interactive:         false,
			forceInteractive:    true,
			forceNonInteractive: false,
			explicitMode:        "",
			hasError:            true,
		},
		{
			name:                "valid explicit mode",
			nonInteractive:      false,
			interactive:         false,
			forceInteractive:    false,
			forceNonInteractive: false,
			explicitMode:        "interactive",
			hasError:            false,
		},
		{
			name:                "invalid explicit mode",
			nonInteractive:      false,
			interactive:         false,
			forceInteractive:    false,
			forceNonInteractive: false,
			explicitMode:        "invalid",
			hasError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flagHandler.ValidateModeFlags(tt.nonInteractive, tt.interactive, tt.forceInteractive, tt.forceNonInteractive, tt.explicitMode)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseBoolEnv(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "true value",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "1 value",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "yes value",
			envValue:     "yes",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "false value",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "0 value",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "empty value uses default",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "invalid value uses default",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				os.Setenv("TEST_VAR", tt.envValue)
				defer os.Unsetenv("TEST_VAR")
			}

			result := flagHandler.parseBoolEnv("TEST_VAR", tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleGlobalFlags(t *testing.T) {
	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	// Create a test command with flags
	cmd := &cobra.Command{
		Use: "test",
	}
	flagHandler.SetupGlobalFlags(cmd)

	// Set up mock expectations
	mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return()
	mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return()
	mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return()

	// Test with default values
	err := flagHandler.HandleGlobalFlags(cmd)
	assert.NoError(t, err)

	// Verify mock calls
	mockLogger.AssertExpectations(t)
}

func TestHandleGlobalFlagsWithConflicts(t *testing.T) {
	mockLogger := &MockLogger{}
	mockLogger.On("SetLevel", mock.AnythingOfType("int")).Return()
	mockLogger.On("SetJSONOutput", mock.AnythingOfType("bool")).Return()
	mockLogger.On("SetCallerInfo", mock.AnythingOfType("bool")).Return()
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	// Create a test command with flags
	cmd := &cobra.Command{
		Use: "test",
	}
	flagHandler.SetupGlobalFlags(cmd)

	// Set conflicting flags
	cmd.PersistentFlags().Set("verbose", "true")
	cmd.PersistentFlags().Set("quiet", "true")

	// Test should return error due to conflicting flags
	err := flagHandler.HandleGlobalFlags(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Flag conflicts detected")
}
