package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIFlagConflicts tests CLI flag conflict detection in integration scenarios
func TestCLIFlagConflicts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Build the CLI binary for testing
	binaryPath := buildCLIBinaryForTest(t)
	defer func() { _ = os.Remove(binaryPath) }()

	t.Run("generate_command_flag_conflicts", func(t *testing.T) {
		testGenerateCommandFlagConflicts(t, binaryPath)
	})

	t.Run("global_flag_conflicts", func(t *testing.T) {
		testGlobalFlagConflicts(t, binaryPath)
	})

	t.Run("enhanced_error_messages", func(t *testing.T) {
		testEnhancedErrorMessages(t, binaryPath)
	})
}

func testGenerateCommandFlagConflicts(t *testing.T, binaryPath string) {
	tests := []struct {
		name          string
		args          []string
		expectError   bool
		errorContains string
	}{
		{
			name:          "interactive_and_non_interactive_conflict",
			args:          []string{"generate", "--interactive", "--non-interactive"},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name:          "force_interactive_and_force_non_interactive_conflict",
			args:          []string{"generate", "--force-interactive", "--force-non-interactive"},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name:          "interactive_with_mode_conflict",
			args:          []string{"generate", "--interactive", "--mode=non-interactive"},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name:        "valid_single_mode_flag",
			args:        []string{"generate", "--non-interactive", "--help"},
			expectError: false,
		},
		{
			name:        "valid_explicit_mode",
			args:        []string{"generate", "--mode=interactive", "--help"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runCLICommandWithOutput(t, binaryPath, tt.args...)

			if tt.expectError {
				if exitCode == 0 {
					t.Errorf("Expected non-zero exit code for conflicting flags, got 0")
				}

				errorOutput := stderr
				if errorOutput == "" {
					errorOutput = stdout
				}

				if !strings.Contains(errorOutput, tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, errorOutput)
				}
			} else {
				if exitCode != 0 && !strings.Contains(stdout, "Usage:") {
					t.Errorf("Expected success or help output, got exit code %d. Stderr: %s", exitCode, stderr)
				}
			}
		})
	}
}

func testGlobalFlagConflicts(t *testing.T, binaryPath string) {
	tests := []struct {
		name          string
		args          []string
		expectError   bool
		errorContains string
	}{
		{
			name:          "verbose_and_quiet_conflict",
			args:          []string{"--verbose", "--quiet", "version"},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name:          "debug_and_quiet_conflict",
			args:          []string{"--debug", "--quiet", "version"},
			expectError:   true,
			errorContains: "Flag conflicts detected",
		},
		{
			name:        "verbose_only",
			args:        []string{"--verbose", "version"},
			expectError: false,
		},
		{
			name:        "debug_only",
			args:        []string{"--debug", "version"},
			expectError: false,
		},
		{
			name:        "quiet_only",
			args:        []string{"--quiet", "version"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runCLICommandWithOutput(t, binaryPath, tt.args...)

			if tt.expectError {
				if exitCode == 0 {
					t.Errorf("Expected non-zero exit code for conflicting flags, got 0")
				}

				errorOutput := stderr
				if errorOutput == "" {
					errorOutput = stdout
				}

				if !strings.Contains(errorOutput, tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.errorContains, errorOutput)
				}
			} else {
				// For non-error cases, we expect either success or help output
				if exitCode != 0 && !strings.Contains(stdout, "version") && !strings.Contains(stderr, "version") {
					t.Errorf("Expected success or version output, got exit code %d. Stderr: %s", exitCode, stderr)
				}
			}
		})
	}
}

func testEnhancedErrorMessages(t *testing.T, binaryPath string) {
	// Test that enhanced error messages provide helpful information
	stdout, stderr, exitCode := runCLICommandWithOutput(t, binaryPath, "generate", "--interactive", "--non-interactive")

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for conflicting flags")
	}

	errorOutput := stderr
	if errorOutput == "" {
		errorOutput = stdout
	}

	// Check for enhanced error message components
	expectedComponents := []string{
		"Flag conflicts detected",
		"Conflict",
		"Suggestion",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(errorOutput, component) {
			t.Errorf("Expected error message to contain '%s', got: %s", component, errorOutput)
		}
	}

	// Verify that the error message is helpful and not just a generic error
	if len(errorOutput) < 50 {
		t.Errorf("Expected detailed error message, got short message: %s", errorOutput)
	}
}

// Helper functions for CLI testing
func buildCLIBinaryForTest(t *testing.T) string {
	t.Helper()

	// Create temporary binary
	binaryPath := filepath.Join(t.TempDir(), "generator")

	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/generator")
	cmd.Dir = getProjectRootForTest()

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, output)
	}

	return binaryPath
}

func getProjectRootForTest() string {
	// Get the project root directory
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "."
}

func runCLICommandWithOutput(t *testing.T, binaryPath string, args ...string) (string, string, int) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			t.Logf("Command execution error: %v", err)
			exitCode = 1
		}
	}

	return stdout.String(), stderr.String(), exitCode
}
