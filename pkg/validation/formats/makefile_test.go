package formats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMakefileValidator(t *testing.T) {
	validator := NewMakefileValidator()
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.commonTargets)
	assert.NotNil(t, validator.shellCommands)
}

func TestMakefileValidator_ValidateMakefile(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name           string
		content        string
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid Makefile",
			content: `CC = gcc
CFLAGS = -Wall -Wextra

all: build

build:
	$(CC) $(CFLAGS) -o app main.c

clean:
	rm -f app

test:
	./app --test

.PHONY: all clean test`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // enhanced validation produces more warnings
		},
		{
			name: "Makefile with spaces instead of tabs",
			content: `all:
    echo "This uses spaces"`,
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 4, // enhanced validation produces more warnings
		},
		{
			name: "Makefile without default target",
			content: `.PHONY: clean
clean:
	rm -f *.o`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // enhanced validation may have different behavior
		},
		{
			name: "Makefile with dangerous commands",
			content: `clean:
	rm -rf /
	chmod 777 /tmp`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // 2 dangerous commands + missing targets
		},
		{
			name: "Makefile with hardcoded paths",
			content: `INSTALL_DIR = /usr/local/bin

install:
	cp app $(INSTALL_DIR)`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // enhanced validation produces more warnings
		},
		{
			name: "Makefile with bash features",
			content: `test:
	[[ -f app ]] && ./app`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // enhanced validation produces more warnings
		},
		{
			name: "Makefile with missing .PHONY declarations",
			content: `all: build test

build:
	go build -o app

test:
	go test ./...

clean:
	rm -f app`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 6, // missing help target + 4 missing .PHONY declarations + should use variables
		},
		{
			name: "Makefile with long lines",
			content: `build:
	go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse HEAD)" -o app ./cmd/main.go`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 7, // long line + missing targets + missing .PHONY + should use variables
		},
		{
			name: "Makefile with trailing whitespace",
			content: `build:   
	go build -o app   `,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 7, // trailing whitespace + missing targets + missing .PHONY + should use variables
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "Makefile")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result, err := validator.ValidateMakefile(filePath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestMakefileValidator_ParseMakefile(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name              string
		content           string
		expectedTargets   int
		expectedVariables int
	}{
		{
			name: "simple Makefile",
			content: `CC = gcc
CFLAGS = -Wall

all: build

build:
	$(CC) $(CFLAGS) -o app main.c

clean:
	rm -f app`,
			expectedTargets:   3, // all, build, clean
			expectedVariables: 2, // CC, CFLAGS
		},
		{
			name: "Makefile with comments",
			content: `# Compiler settings
CC = gcc
CFLAGS = -Wall

# Default target
all: build

# Build the application
build:
	$(CC) $(CFLAGS) -o app main.c`,
			expectedTargets:   2, // all, build
			expectedVariables: 2, // CC, CFLAGS
		},
		{
			name: "Makefile with .PHONY",
			content: `.PHONY: all clean

all: build

build:
	gcc -o app main.c

clean:
	rm -f app`,
			expectedTargets:   4, // enhanced parsing may include .PHONY as a target
			expectedVariables: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targets, variables := validator.parseMakefile(tt.content)

			assert.Equal(t, tt.expectedTargets, len(targets))
			assert.Equal(t, tt.expectedVariables, len(variables))
		})
	}
}

func TestMakefileValidator_ValidateTargetName(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name        string
		targetName  string
		expectError bool
	}{
		{
			name:        "valid lowercase target",
			targetName:  "build",
			expectError: false,
		},
		{
			name:        "valid target with hyphens",
			targetName:  "build-all",
			expectError: false,
		},
		{
			name:        "valid target with underscores",
			targetName:  "build_all",
			expectError: false,
		},
		{
			name:        "valid target with dots",
			targetName:  "build.o",
			expectError: false,
		},
		{
			name:        "empty target name",
			targetName:  "",
			expectError: true,
		},
		{
			name:        "target with spaces",
			targetName:  "build all",
			expectError: true,
		},
		{
			name:        "uppercase target",
			targetName:  "BUILD",
			expectError: true,
		},
		{
			name:        "target with invalid characters",
			targetName:  "build@all",
			expectError: true,
		},
		{
			name:        "special target (allowed)",
			targetName:  ".PHONY",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateTargetName(tt.targetName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMakefileValidator_ValidateVariableName(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name         string
		variableName string
		expectError  bool
	}{
		{
			name:         "valid uppercase variable",
			variableName: "CC",
			expectError:  false,
		},
		{
			name:         "valid variable with underscores",
			variableName: "CFLAGS_DEBUG",
			expectError:  false,
		},
		{
			name:         "valid variable with numbers",
			variableName: "VERSION_2",
			expectError:  false,
		},
		{
			name:         "empty variable name",
			variableName: "",
			expectError:  true,
		},
		{
			name:         "variable with spaces",
			variableName: "C FLAGS",
			expectError:  true,
		},
		{
			name:         "lowercase variable",
			variableName: "cflags",
			expectError:  true,
		},
		{
			name:         "variable with hyphens",
			variableName: "C-FLAGS",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateVariableName(tt.variableName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMakefileValidator_ValidateTarget(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name           string
		target         MakefileTarget
		expectWarnings int
	}{
		{
			name: "valid target with commands",
			target: MakefileTarget{
				Name:     "build",
				Commands: []string{"gcc -o app main.c"},
				LineNum:  1,
			},
			expectWarnings: 2, // enhanced validation may suggest improvements
		},
		{
			name: "empty target",
			target: MakefileTarget{
				Name:         "empty",
				Commands:     []string{},
				Dependencies: []string{},
				LineNum:      1,
			},
			expectWarnings: 1,
		},
		{
			name: "target with invalid name",
			target: MakefileTarget{
				Name:     "Build All",
				Commands: []string{"echo test"},
				LineNum:  1,
			},
			expectWarnings: 1,
		},
		{
			name: "phony target that should be declared",
			target: MakefileTarget{
				Name:     "clean",
				Commands: []string{"rm -f *.o"},
				LineNum:  1,
				IsPhony:  false,
			},
			expectWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			validator.validateTarget(tt.target, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestMakefileValidator_ValidateVariable(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name           string
		variable       MakefileVariable
		expectWarnings int
	}{
		{
			name: "valid variable",
			variable: MakefileVariable{
				Name:    "CC",
				Value:   "gcc",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 0,
		},
		{
			name: "variable with invalid name",
			variable: MakefileVariable{
				Name:    "cc",
				Value:   "gcc",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 2, // enhanced validation may produce additional warnings
		},
		{
			name: "empty variable",
			variable: MakefileVariable{
				Name:    "EMPTY",
				Value:   "",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 1,
		},
		{
			name: "variable with hardcoded path",
			variable: MakefileVariable{
				Name:    "INSTALL_DIR",
				Value:   "/usr/local/bin",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 1,
		},
		{
			name: "recursive assignment with self-reference",
			variable: MakefileVariable{
				Name:    "PATH",
				Value:   "$(PATH):/usr/local/bin",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 2, // enhanced validation may produce additional warnings
		},
		{
			name: "variable with Windows path",
			variable: MakefileVariable{
				Name:    "INSTALL_DIR",
				Value:   "C:\\Program Files\\MyApp",
				LineNum: 1,
				Type:    "=",
			},
			expectWarnings: 1,
		},
		{
			name: "variable with immediate assignment",
			variable: MakefileVariable{
				Name:    "BUILD_TIME",
				Value:   "$(shell date)",
				LineNum: 1,
				Type:    ":=",
			},
			expectWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			validator.validateVariable(tt.variable, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestMakefileValidator_ValidateCommand(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name           string
		targetName     string
		command        string
		expectWarnings int
	}{
		{
			name:           "safe command",
			targetName:     "build",
			command:        "gcc -o app main.c",
			expectWarnings: 1, // enhanced validation may suggest using variables
		},
		{
			name:           "dangerous command",
			targetName:     "clean",
			command:        "rm -rf /",
			expectWarnings: 1,
		},
		{
			name:           "command that should use variables",
			targetName:     "build",
			command:        "gcc -o app main.c",
			expectWarnings: 1, // should use $(CC)
		},
		{
			name:           "command needing error handling",
			targetName:     "download",
			command:        "curl -o file.tar.gz http://example.com/file.tar.gz",
			expectWarnings: 1,
		},
		{
			name:           "bash-specific command",
			targetName:     "test",
			command:        "[[ -f app ]] && ./app",
			expectWarnings: 1,
		},
		{
			name:           "docker command needing error handling",
			targetName:     "docker-build",
			command:        "docker build -t myapp .",
			expectWarnings: 2, // should use variables + needs error handling
		},
		{
			name:           "npm install command",
			targetName:     "install",
			command:        "npm install",
			expectWarnings: 2, // should use variables + needs error handling
		},
		{
			name:           "git command needing error handling",
			targetName:     "update",
			command:        "git pull origin main",
			expectWarnings: 1, // needs error handling
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &interfaces.ConfigValidationResult{
				Valid:    true,
				Errors:   []interfaces.ConfigValidationError{},
				Warnings: []interfaces.ConfigValidationError{},
				Summary: interfaces.ConfigValidationSummary{
					TotalProperties: 0,
					ValidProperties: 0,
					ErrorCount:      0,
					WarningCount:    0,
					MissingRequired: 0,
				},
			}

			validator.validateCommand(tt.targetName, tt.command, 0, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestMakefileValidator_ShouldBePhony(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name     string
		target   MakefileTarget
		expected bool
	}{
		{
			name: "clean target should be phony",
			target: MakefileTarget{
				Name: "clean",
			},
			expected: true,
		},
		{
			name: "all target should be phony",
			target: MakefileTarget{
				Name: "all",
			},
			expected: true,
		},
		{
			name: "test target should be phony",
			target: MakefileTarget{
				Name: "test",
			},
			expected: true,
		},
		{
			name: "file target should not be phony",
			target: MakefileTarget{
				Name: "app.o",
			},
			expected: false,
		},
		{
			name: "custom target should not be phony",
			target: MakefileTarget{
				Name: "custom",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.shouldBePhony(tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMakefileValidator_ContainsHardcodedPaths(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "relative path",
			value:    "./bin/app",
			expected: false,
		},
		{
			name:     "hardcoded /usr/local path",
			value:    "/usr/local/bin",
			expected: true,
		},
		{
			name:     "hardcoded /opt path",
			value:    "/opt/myapp",
			expected: true,
		},
		{
			name:     "Windows path",
			value:    "C:\\Program Files\\MyApp",
			expected: true,
		},
		{
			name:     "home directory path",
			value:    "/home/user/app",
			expected: true,
		},
		{
			name:     "variable reference",
			value:    "$(PREFIX)/bin",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.containsHardcodedPaths(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMakefileValidator_LooksLikeCommand(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "echo command",
			line:     "    echo 'Hello World'",
			expected: true,
		},
		{
			name:     "gcc command",
			line:     "    gcc -o app main.c",
			expected: false, // gcc not in the basic patterns
		},
		{
			name:     "go command",
			line:     "    go build -o app",
			expected: true,
		},
		{
			name:     "command with pipe",
			line:     "    cat file | grep pattern",
			expected: true,
		},
		{
			name:     "command with &&",
			line:     "    mkdir -p bin && cp app bin/",
			expected: true,
		},
		{
			name:     "variable assignment",
			line:     "    VAR = value",
			expected: false,
		},
		{
			name:     "comment",
			line:     "    # This is a comment",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.looksLikeCommand(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMakefileValidator_IsVariableAssignment(t *testing.T) {
	validator := NewMakefileValidator()

	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "simple assignment",
			line:     "CC = gcc",
			expected: true,
		},
		{
			name:     "immediate assignment",
			line:     "CC := gcc",
			expected: true,
		},
		{
			name:     "append assignment",
			line:     "CFLAGS += -Wall",
			expected: true,
		},
		{
			name:     "conditional assignment",
			line:     "CC ?= gcc",
			expected: true,
		},
		{
			name:     "target definition",
			line:     "all: build",
			expected: false,
		},
		{
			name:     "command line",
			line:     "\techo 'Hello'",
			expected: false,
		},
		{
			name:     "comment",
			line:     "# This is a comment",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isVariableAssignment(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMakefileValidator_FileReadError(t *testing.T) {
	validator := NewMakefileValidator()

	// Test with non-existent file
	result, err := validator.ValidateMakefile("/non/existent/Makefile")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "read_error", result.Errors[0].Type)
}
