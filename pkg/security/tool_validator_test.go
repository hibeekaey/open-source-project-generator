package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolValidator_ValidateToolCommand(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		name      string
		toolName  string
		args      []string
		shouldErr bool
	}{
		{
			name:      "valid npx command",
			toolName:  "npx",
			args:      []string{"create-next-app", "--typescript", "--tailwind"},
			shouldErr: false,
		},
		{
			name:      "valid go command",
			toolName:  "go",
			args:      []string{"mod", "init", "github.com/user/project"},
			shouldErr: false,
		},
		{
			name:      "invalid tool",
			toolName:  "rm",
			args:      []string{"-rf", "/"},
			shouldErr: true,
		},
		{
			name:      "command injection attempt",
			toolName:  "npx",
			args:      []string{"create-next-app", "; rm -rf /"},
			shouldErr: true,
		},
		{
			name:      "too many arguments",
			toolName:  "npx",
			args:      make([]string, 100),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateToolCommand(tt.toolName, tt.args)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToolValidator_ValidateToolFlags(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		name      string
		flags     []string
		shouldErr bool
	}{
		{
			name:      "valid flags",
			flags:     []string{"--typescript", "--tailwind", "--app"},
			shouldErr: false,
		},
		{
			name:      "flag with semicolon",
			flags:     []string{"--flag; rm -rf /"},
			shouldErr: true,
		},
		{
			name:      "flag with pipe",
			flags:     []string{"--flag | cat /etc/passwd"},
			shouldErr: true,
		},
		{
			name:      "flag with backtick",
			flags:     []string{"--flag `whoami`"},
			shouldErr: true,
		},
		{
			name:      "flag with path traversal",
			flags:     []string{"--path", "../../../etc/passwd"},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateToolFlags(tt.flags)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToolValidator_IsToolWhitelisted(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		toolName    string
		whitelisted bool
	}{
		{"npx", true},
		{"go", true},
		{"gradle", true},
		{"xcodebuild", true},
		{"swift", true},
		{"docker", true},
		{"terraform", true},
		{"rm", false},
		{"curl", false},
		{"wget", false},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			result := validator.IsToolWhitelisted(tt.toolName)
			assert.Equal(t, tt.whitelisted, result)
		})
	}
}

func TestToolValidator_ValidateCommandString(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		name      string
		command   string
		shouldErr bool
	}{
		{
			name:      "clean command",
			command:   "npx create-next-app my-app",
			shouldErr: false,
		},
		{
			name:      "command with semicolon",
			command:   "npx create-next-app; rm -rf /",
			shouldErr: true,
		},
		{
			name:      "command with pipe",
			command:   "npx create-next-app | cat",
			shouldErr: true,
		},
		{
			name:      "command with redirection",
			command:   "npx create-next-app > /dev/null",
			shouldErr: true,
		},
		{
			name:      "command chaining",
			command:   "npx create-next-app && npm install",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCommandString(tt.command)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToolValidator_SanitizeCommandOutput(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "output with home path",
			input:    "Installing to /home/john/project",
			expected: "Installing to /home/[user]/project",
		},
		{
			name:     "output with Users path",
			input:    "Installing to /Users/jane/project",
			expected: "Installing to /Users/[user]/project",
		},
		{
			name:     "clean output",
			input:    "Installation complete",
			expected: "Installation complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizeCommandOutput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToolValidator_AddToolWhitelist(t *testing.T) {
	validator := NewToolValidator()

	whitelist := &ToolWhitelist{
		Command:      "custom-tool",
		AllowedFlags: []string{"--flag1", "--flag2"},
		MaxFlagCount: 10,
	}

	err := validator.AddToolWhitelist("custom-tool", whitelist)
	require.NoError(t, err)

	assert.True(t, validator.IsToolWhitelisted("custom-tool"))
}

func TestToolValidator_ValidateToolPath(t *testing.T) {
	validator := NewToolValidator()

	tests := []struct {
		name      string
		path      string
		shouldErr bool
	}{
		{
			name:      "valid path",
			path:      "/usr/bin/npx",
			shouldErr: false,
		},
		{
			name:      "path with traversal",
			path:      "/usr/bin/../../../etc/passwd",
			shouldErr: true,
		},
		{
			name:      "suspicious tmp path",
			path:      "/tmp/malicious-tool",
			shouldErr: true,
		},
		{
			name:      "suspicious dev path",
			path:      "/dev/null",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateToolPath(tt.path)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
