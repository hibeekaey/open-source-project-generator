// Package utils provides helper utilities for the CLI interface.
//
// This module contains common helper functions for CLI operations including
// environment detection, mode detection, and utility functions.
package utils

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Helper provides common utility functions for the CLI.
//
// The Helper provides methods for:
//   - Environment detection and analysis
//   - Mode detection (interactive, CI, etc.)
//   - String manipulation utilities
//   - System information gathering
type Helper struct{}

// NewHelper creates a new helper instance.
func NewHelper() *Helper {
	return &Helper{}
}

// DetectNonInteractiveMode detects if the CLI is running in non-interactive mode
func (h *Helper) DetectNonInteractiveMode(cmd *cobra.Command) bool {
	// Check if explicitly set via flag
	if cmd != nil {
		if nonInteractive, err := cmd.Flags().GetBool("non-interactive"); err == nil && nonInteractive {
			return true
		}
	}

	// Check for CI environment variables
	ciEnvVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"GITHUB_ACTIONS",
		"GITLAB_CI",
		"JENKINS_URL",
		"CIRCLECI",
		"TRAVIS",
		"BUILDKITE",
		"DRONE",
		"TEAMCITY_VERSION",
	}

	for _, envVar := range ciEnvVars {
		if value := os.Getenv(envVar); value != "" && value != "false" && value != "0" {
			return true
		}
	}

	// Check if stdin is not a terminal (piped input)
	if !h.isTerminal() {
		return true
	}

	return false
}

// isTerminal checks if stdin is connected to a terminal
func (h *Helper) isTerminal() bool {
	// This is a simplified check - in a real implementation,
	// you would use platform-specific APIs

	// Check if we're running in a known non-terminal environment
	if os.Getenv("TERM") == "" {
		return false
	}

	// Check for common non-interactive indicators
	if os.Getenv("DEBIAN_FRONTEND") == "noninteractive" {
		return false
	}

	return true
}

// GetSystemInfo returns system information
func (h *Helper) GetSystemInfo() map[string]string {
	return map[string]string{
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"go_version":  runtime.Version(),
		"num_cpu":     strconv.Itoa(runtime.NumCPU()),
		"hostname":    h.getHostname(),
		"user":        h.getCurrentUser(),
		"working_dir": h.getWorkingDirectory(),
	}
}

// getHostname returns the system hostname
func (h *Helper) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getCurrentUser returns the current user
func (h *Helper) getCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

// getWorkingDirectory returns the current working directory
func (h *Helper) getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}

// SanitizeInput sanitizes user input by removing control characters
func (h *Helper) SanitizeInput(input string) string {
	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ParseBoolFromString parses a boolean value from string with flexible formats
func (h *Helper) ParseBoolFromString(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))

	switch s {
	case "true", "yes", "y", "1", "on", "enable", "enabled":
		return true, nil
	case "false", "no", "n", "0", "off", "disable", "disabled":
		return false, nil
	default:
		return strconv.ParseBool(s)
	}
}

// GenerateTimestamp generates a timestamp string
func (h *Helper) GenerateTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
}

// GenerateShortTimestamp generates a short timestamp string
func (h *Helper) GenerateShortTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// SplitAndTrim splits a string by delimiter and trims whitespace
func (h *Helper) SplitAndTrim(s, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ContainsString checks if a slice contains a string (case-sensitive)
func (h *Helper) ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsStringIgnoreCase checks if a slice contains a string (case-insensitive)
func (h *Helper) ContainsStringIgnoreCase(slice []string, item string) bool {
	lowerItem := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == lowerItem {
			return true
		}
	}
	return false
}

// RemoveDuplicateStrings removes duplicate strings from a slice
func (h *Helper) RemoveDuplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// JoinNonEmpty joins non-empty strings with a separator
func (h *Helper) JoinNonEmpty(parts []string, separator string) string {
	nonEmpty := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, separator)
}

// ExpandPath expands ~ in file paths to the user's home directory
func (h *Helper) ExpandPath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path // Return original path if we can't get home dir
	}

	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return strings.Replace(path, "~", homeDir, 1)
	}

	return path
}

// GetEnvWithDefault gets an environment variable with a default value
func (h *Helper) GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsBool gets an environment variable as a boolean
func (h *Helper) GetEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := h.ParseBoolFromString(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// GetEnvAsInt gets an environment variable as an integer
func (h *Helper) GetEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// CreateErrorWithSuggestions creates an error message with suggestions
func (h *Helper) CreateErrorWithSuggestions(message string, suggestions []string) string {
	if len(suggestions) == 0 {
		return message
	}

	var result strings.Builder
	result.WriteString(message)
	result.WriteString("\n\nSuggestions:")

	for _, suggestion := range suggestions {
		result.WriteString("\n  • ")
		result.WriteString(suggestion)
	}

	return result.String()
}

// ValidateRequiredEnvVars validates that required environment variables are set
func (h *Helper) ValidateRequiredEnvVars(required []string) []string {
	var missing []string

	for _, envVar := range required {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}

	return missing
}

// MaskSensitiveValue masks sensitive values for logging
func (h *Helper) MaskSensitiveValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}

	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// IsValidURL performs basic URL validation
func (h *Helper) IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// NormalizeLineEndings normalizes line endings to Unix format
func (h *Helper) NormalizeLineEndings(text string) string {
	// Replace Windows line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Replace old Mac line endings
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}
