package cli

import (
	"testing"
)

func TestDetectGenerationMode(t *testing.T) {
	// Create a minimal CLI instance for testing
	cli := &CLI{}

	// Check if we're running in CI environment
	isCI := cli.detectCIEnvironment().IsCI

	tests := []struct {
		name           string
		configPath     string
		nonInteractive bool
		interactive    bool
		explicitMode   string
		expected       string
	}{
		{
			name:       "config file mode",
			configPath: "config.yaml",
			expected:   "config-file",
		},
		{
			name:           "explicit non-interactive",
			nonInteractive: true,
			expected:       "non-interactive",
		},
		{
			name:        "explicit interactive",
			interactive: true,
			expected: func() string {
				if isCI {
					return "non-interactive" // CI detection overrides explicit interactive flag
				}
				return "interactive"
			}(),
		},
		{
			name:         "explicit mode interactive",
			explicitMode: "interactive",
			expected:     "interactive",
		},
		{
			name:         "explicit mode non-interactive",
			explicitMode: "non-interactive",
			expected:     "non-interactive",
		},
		{
			name:         "explicit mode config-file",
			explicitMode: "config-file",
			expected:     "config-file",
		},
		{
			name:     "default to interactive",
			expected: func() string {
				if isCI {
					return "non-interactive"
				}
				return "interactive"
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cli.detectGenerationMode(tt.configPath, tt.nonInteractive, tt.interactive, tt.explicitMode)
			if result != tt.expected {
				t.Errorf("detectGenerationMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateModeFlags(t *testing.T) {
	cli := &CLI{}

	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectError         bool
	}{
		{
			name:        "no conflicts",
			interactive: true,
			expectError: false,
		},
		{
			name:           "single non-interactive",
			nonInteractive: true,
			expectError:    false,
		},
		{
			name:             "single force interactive",
			forceInteractive: true,
			expectError:      false,
		},
		{
			name:         "single explicit mode",
			explicitMode: "interactive",
			expectError:  false,
		},
		{
			name:             "conflict: non-interactive and force interactive",
			nonInteractive:   true,
			forceInteractive: true,
			expectError:      true,
		},
		{
			name:                "conflict: force interactive and force non-interactive",
			forceInteractive:    true,
			forceNonInteractive: true,
			expectError:         true,
		},
		{
			name:           "conflict: non-interactive and explicit mode",
			nonInteractive: true,
			explicitMode:   "interactive",
			expectError:    true,
		},
		{
			name:         "invalid explicit mode",
			explicitMode: "invalid-mode",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cli.validateModeFlags(tt.nonInteractive, tt.interactive, tt.forceInteractive, tt.forceNonInteractive, tt.explicitMode)
			if (err != nil) != tt.expectError {
				t.Errorf("validateModeFlags() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestApplyModeOverrides(t *testing.T) {
	cli := &CLI{}

	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectedNI          bool
		expectedI           bool
	}{
		{
			name:        "no overrides",
			interactive: true,
			expectedNI:  false,
			expectedI:   true,
		},
		{
			name:             "force interactive",
			forceInteractive: true,
			expectedNI:       false,
			expectedI:        true,
		},
		{
			name:                "force non-interactive",
			forceNonInteractive: true,
			expectedNI:          true,
			expectedI:           false,
		},
		{
			name:         "explicit mode interactive",
			explicitMode: "interactive",
			expectedNI:   false,
			expectedI:    true,
		},
		{
			name:         "explicit mode non-interactive",
			explicitMode: "non-interactive",
			expectedNI:   true,
			expectedI:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultNI, resultI := cli.applyModeOverrides(tt.nonInteractive, tt.interactive, tt.forceInteractive, tt.forceNonInteractive, tt.explicitMode)
			if resultNI != tt.expectedNI || resultI != tt.expectedI {
				t.Errorf("applyModeOverrides() = (%v, %v), want (%v, %v)", resultNI, resultI, tt.expectedNI, tt.expectedI)
			}
		})
	}
}

func TestValidateAndNormalizeMode(t *testing.T) {
	cli := &CLI{}

	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{
			name:     "interactive",
			mode:     "interactive",
			expected: "interactive",
		},
		{
			name:     "interactive short",
			mode:     "i",
			expected: "interactive",
		},
		{
			name:     "non-interactive",
			mode:     "non-interactive",
			expected: "non-interactive",
		},
		{
			name:     "non-interactive variants",
			mode:     "noninteractive",
			expected: "non-interactive",
		},
		{
			name:     "config-file",
			mode:     "config-file",
			expected: "config-file",
		},
		{
			name:     "config variants",
			mode:     "config",
			expected: "config-file",
		},
		{
			name:     "invalid mode",
			mode:     "invalid",
			expected: "interactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cli.validateAndNormalizeMode(tt.mode)
			if result != tt.expected {
				t.Errorf("validateAndNormalizeMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}
