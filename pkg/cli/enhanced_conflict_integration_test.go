package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEnhancedFlagConflictDetectionIntegration tests the complete enhanced flag conflict detection system
func TestEnhancedFlagConflictDetectionIntegration(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		flagState      map[string]bool
		expectConflict bool
		conflictType   string
	}{
		{
			name:        "no_conflicts_verbose_only",
			description: "Single verbose flag should not cause conflicts",
			flagState: map[string]bool{
				"--verbose": true,
			},
			expectConflict: false,
		},
		{
			name:        "verbose_quiet_conflict",
			description: "Verbose and quiet flags should conflict",
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   true,
			},
			expectConflict: true,
			conflictType:   "output_mode",
		},
		{
			name:        "debug_quiet_conflict",
			description: "Debug and quiet flags should conflict",
			flagState: map[string]bool{
				"--debug": true,
				"--quiet": true,
			},
			expectConflict: true,
			conflictType:   "output_mode",
		},
		{
			name:        "interactive_noninteractive_conflict",
			description: "Interactive and non-interactive modes should conflict",
			flagState: map[string]bool{
				"--interactive":     true,
				"--non-interactive": true,
			},
			expectConflict: true,
			conflictType:   "generation_mode",
		},
		{
			name:        "force_flags_conflict",
			description: "Force interactive and force non-interactive should conflict",
			flagState: map[string]bool{
				"--force-interactive":     true,
				"--force-non-interactive": true,
			},
			expectConflict: true,
			conflictType:   "generation_mode",
		},
		{
			name:        "interactive_with_mode_conflict",
			description: "Interactive flag with explicit mode should conflict",
			flagState: map[string]bool{
				"--interactive": true,
				"--mode":        true,
			},
			expectConflict: true,
			conflictType:   "generation_mode",
		},
		{
			name:        "multiple_conflicts",
			description: "Multiple conflicts should be detected",
			flagState: map[string]bool{
				"--verbose":         true,
				"--quiet":           true,
				"--interactive":     true,
				"--non-interactive": true,
			},
			expectConflict: true,
			conflictType:   "multiple",
		},
		{
			name:        "compatible_flags",
			description: "Compatible flags should not conflict",
			flagState: map[string]bool{
				"--verbose":     true,
				"--interactive": true,
				"--force":       true,
			},
			expectConflict: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create conflict matrix
			matrix := NewFlagConflictMatrix()
			fh := &FlagHandler{}

			// Check for conflicts using the enhanced system
			hasConflict := false
			conflictCount := 0

			for _, rule := range matrix.Rules {
				if fh.checkConflictRule(rule, tt.flagState) {
					hasConflict = true
					conflictCount++
				}
			}

			// Verify expectations
			assert.Equal(t, tt.expectConflict, hasConflict,
				"Test %s: %s", tt.name, tt.description)

			if tt.expectConflict {
				assert.Greater(t, conflictCount, 0,
					"Should have detected at least one conflict for %s", tt.name)
			} else {
				assert.Equal(t, 0, conflictCount,
					"Should not have detected any conflicts for %s", tt.name)
			}
		})
	}
}

// TestConflictMatrixComprehensiveness ensures the conflict matrix covers all expected scenarios
func TestConflictMatrixComprehensiveness(t *testing.T) {
	matrix := NewFlagConflictMatrix()

	// Define expected conflict categories
	expectedCategories := map[string][]string{
		"output_modes": {
			"--verbose,--quiet",
			"--debug,--quiet",
		},
		"generation_modes": {
			"--interactive,--non-interactive",
			"--force-interactive,--force-non-interactive",
			"--interactive,--force-non-interactive",
			"--non-interactive,--force-interactive",
		},
		"mode_specifications": {
			"--interactive,--mode",
			"--non-interactive,--mode",
			"--force-interactive,--mode",
			"--force-non-interactive,--mode",
		},
	}

	// Verify all expected conflicts are covered
	for category, expectedConflicts := range expectedCategories {
		t.Run(category, func(t *testing.T) {
			for _, expectedConflict := range expectedConflicts {
				found := false
				expectedFlags := strings.Split(expectedConflict, ",")

				for _, rule := range matrix.Rules {
					if containsAllFlags(rule.Flags, expectedFlags) {
						found = true

						// Verify rule quality
						assert.NotEmpty(t, rule.Description,
							"Rule for %v should have description", expectedFlags)
						assert.NotEmpty(t, rule.Suggestion,
							"Rule for %v should have suggestion", expectedFlags)
						assert.NotEmpty(t, rule.Examples,
							"Rule for %v should have examples", expectedFlags)
						assert.Contains(t, []string{"error", "warning", "info"}, rule.Severity,
							"Rule for %v should have valid severity", expectedFlags)
						break
					}
				}

				assert.True(t, found,
					"Expected conflict rule for %s not found in category %s",
					expectedConflict, category)
			}
		})
	}
}

// TestGracefulFallbackHandling tests the graceful fallback functionality
func TestGracefulFallbackHandling(t *testing.T) {
	tests := []struct {
		name                string
		nonInteractive      bool
		interactive         bool
		forceInteractive    bool
		forceNonInteractive bool
		explicitMode        string
		expectRecoverable   bool
		expectedResolution  string
	}{
		{
			name:               "force_interactive_overrides",
			interactive:        true,
			forceInteractive:   true,
			expectRecoverable:  true,
			expectedResolution: "interactive",
		},
		{
			name:                "force_noninteractive_overrides",
			nonInteractive:      true,
			forceNonInteractive: true,
			expectRecoverable:   true,
			expectedResolution:  "non-interactive",
		},
		{
			name:              "non_recoverable_multiple_conflicts",
			interactive:       true,
			nonInteractive:    true,
			forceInteractive:  true,
			expectRecoverable: false,
		},
		{
			name:              "non_recoverable_basic_conflict",
			interactive:       true,
			nonInteractive:    true,
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

// TestEnhancedErrorMessages verifies that enhanced error messages are properly formatted
func TestEnhancedErrorMessages(t *testing.T) {
	// This test verifies the structure and content of enhanced error messages
	// without requiring the full CLI dependency

	matrix := NewFlagConflictMatrix()

	// Verify each rule has comprehensive error information
	for i, rule := range matrix.Rules {
		t.Run(fmt.Sprintf("Rule_%d", i), func(t *testing.T) {
			// Check description quality
			assert.NotEmpty(t, rule.Description)
			assert.True(t, len(rule.Description) > 10,
				"Description should be descriptive: %s", rule.Description)

			// Check suggestion quality
			assert.NotEmpty(t, rule.Suggestion)
			assert.True(t, len(rule.Suggestion) > 20,
				"Suggestion should be helpful: %s", rule.Suggestion)

			// Check examples
			assert.NotEmpty(t, rule.Examples)
			assert.True(t, len(rule.Examples) >= 1,
				"Should have at least one example")

			// Verify examples contain actual flag usage
			for _, example := range rule.Examples {
				assert.True(t, strings.Contains(example, "--") || strings.Contains(example, "generator"),
					"Example should contain flag usage or command: %s", example)
			}

			// Check severity
			assert.Contains(t, []string{"error", "warning", "info"}, rule.Severity)
		})
	}
}
