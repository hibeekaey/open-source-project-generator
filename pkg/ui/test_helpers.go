package ui

import (
	"os"
	"testing"
)

// skipIfNotInteractive skips the test if not running in an interactive environment
func skipIfNotInteractive(t *testing.T) {
	// Skip if CI environment variable is set
	if os.Getenv("CI") != "" {
		t.Skip("Skipping interactive test in CI environment")
	}

	// Skip if not running with a TTY
	if !isTerminal() {
		t.Skip("Skipping interactive test - not running in a terminal")
	}
}

// isTerminal checks if we're running in a terminal
func isTerminal() bool {
	if fileInfo, _ := os.Stdin.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
