package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFlagConflictMatrix(t *testing.T) {
	matrix := NewFlagConflictMatrix()

	assert.NotNil(t, matrix)
	assert.NotEmpty(t, matrix.Rules)

	// Verify we have the expected conflict rules
	expectedConflicts := [][]string{
		{"--verbose", "--quiet"},
		{"--debug", "--quiet"},
		{"--interactive", "--non-interactive"},
		{"--force-interactive", "--force-non-interactive"},
		{"--interactive", "--force-non-interactive"},
		{"--non-interactive", "--force-interactive"},
		{"--interactive", "--mode"},
		{"--non-interactive", "--mode"},
		{"--force-interactive", "--mode"},
		{"--force-non-interactive", "--mode"},
	}

	// Check that all expected conflicts are present
	for _, expectedFlags := range expectedConflicts {
		found := false
		for _, rule := range matrix.Rules {
			if containsAllFlags(rule.Flags, expectedFlags) {
				found = true
				// Verify rule has required fields
				assert.NotEmpty(t, rule.Description, "Rule for %v should have description", expectedFlags)
				assert.NotEmpty(t, rule.Suggestion, "Rule for %v should have suggestion", expectedFlags)
				assert.NotEmpty(t, rule.Examples, "Rule for %v should have examples", expectedFlags)
				assert.NotEmpty(t, rule.Severity, "Rule for %v should have severity", expectedFlags)
				break
			}
		}
		assert.True(t, found, "Expected conflict rule for flags %v not found", expectedFlags)
	}
}

func TestConflictRuleStructure(t *testing.T) {
	matrix := NewFlagConflictMatrix()

	for i, rule := range matrix.Rules {
		t.Run(fmt.Sprintf("Rule_%d", i), func(t *testing.T) {
			// Each rule should have at least 2 flags
			assert.GreaterOrEqual(t, len(rule.Flags), 2, "Rule should have at least 2 conflicting flags")

			// Each rule should have all required fields
			assert.NotEmpty(t, rule.Description, "Rule should have description")
			assert.NotEmpty(t, rule.Suggestion, "Rule should have suggestion")
			assert.NotEmpty(t, rule.Examples, "Rule should have examples")
			assert.Contains(t, []string{"error", "warning", "info"}, rule.Severity, "Rule should have valid severity")

			// Flags should start with --
			for _, flag := range rule.Flags {
				assert.True(t, strings.HasPrefix(flag, "--"), "Flag %s should start with --", flag)
			}
		})
	}
}

func TestFlagHandler_checkConflictRule_Isolated(t *testing.T) {
	tests := []struct {
		name           string
		rule           ConflictRule
		flagState      map[string]bool
		expectConflict bool
	}{
		{
			name: "no conflict - only one flag active",
			rule: ConflictRule{
				Flags: []string{"--verbose", "--quiet"},
			},
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   false,
			},
			expectConflict: false,
		},
		{
			name: "conflict detected - both flags active",
			rule: ConflictRule{
				Flags: []string{"--verbose", "--quiet"},
			},
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   true,
			},
			expectConflict: true,
		},
		{
			name: "no conflict - no flags active",
			rule: ConflictRule{
				Flags: []string{"--interactive", "--non-interactive"},
			},
			flagState: map[string]bool{
				"--interactive":     false,
				"--non-interactive": false,
			},
			expectConflict: false,
		},
		{
			name: "conflict with mode flag",
			rule: ConflictRule{
				Flags: []string{"--interactive", "--mode"},
			},
			flagState: map[string]bool{
				"--interactive": true,
				"--mode":        true,
			},
			expectConflict: true,
		},
		{
			name: "three-way conflict",
			rule: ConflictRule{
				Flags: []string{"--interactive", "--non-interactive", "--mode"},
			},
			flagState: map[string]bool{
				"--interactive":     true,
				"--non-interactive": true,
				"--mode":            true,
			},
			expectConflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &FlagHandler{}

			result := fh.checkConflictRule(tt.rule, tt.flagState)

			assert.Equal(t, tt.expectConflict, result)
		})
	}
}

func TestFlagHandler_isFlagActive_Isolated(t *testing.T) {
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
		{
			name: "flag not in state",
			flag: "--unknown",
			flagState: map[string]bool{
				"--verbose": true,
			},
			expected: false,
		},
		{
			name: "mode flag without value set",
			flag: "--mode=interactive",
			flagState: map[string]bool{
				"--mode": false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fh := &FlagHandler{}

			result := fh.isFlagActive(tt.flag, tt.flagState)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConflictDetectionLogic(t *testing.T) {
	matrix := NewFlagConflictMatrix()
	fh := &FlagHandler{}

	// Test various flag combinations
	testCases := []struct {
		name           string
		flagState      map[string]bool
		expectConflict bool
		description    string
	}{
		{
			name: "verbose and quiet conflict",
			flagState: map[string]bool{
				"--verbose": true,
				"--quiet":   true,
			},
			expectConflict: true,
			description:    "Should detect verbose/quiet conflict",
		},
		{
			name: "debug and quiet conflict",
			flagState: map[string]bool{
				"--debug": true,
				"--quiet": true,
			},
			expectConflict: true,
			description:    "Should detect debug/quiet conflict",
		},
		{
			name: "interactive mode conflicts",
			flagState: map[string]bool{
				"--interactive":     true,
				"--non-interactive": true,
			},
			expectConflict: true,
			description:    "Should detect interactive mode conflicts",
		},
		{
			name: "force flag conflicts",
			flagState: map[string]bool{
				"--force-interactive":     true,
				"--force-non-interactive": true,
			},
			expectConflict: true,
			description:    "Should detect force flag conflicts",
		},
		{
			name: "mode flag conflicts",
			flagState: map[string]bool{
				"--interactive": true,
				"--mode":        true,
			},
			expectConflict: true,
			description:    "Should detect mode flag conflicts",
		},
		{
			name: "no conflicts",
			flagState: map[string]bool{
				"--verbose":     true,
				"--interactive": true,
			},
			expectConflict: false,
			description:    "Should not detect conflicts for compatible flags",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasConflict := false

			// Check each rule against the flag state
			for _, rule := range matrix.Rules {
				if fh.checkConflictRule(rule, tc.flagState) {
					hasConflict = true
					break
				}
			}

			assert.Equal(t, tc.expectConflict, hasConflict, tc.description)
		})
	}
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
