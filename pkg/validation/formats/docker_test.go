package formats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDockerValidator(t *testing.T) {
	validator := NewDockerValidator()
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.securityRules)
}

func TestDockerValidator_ValidateDockerfile(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		content        string
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid Dockerfile",
			content: `FROM node:16-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production && rm -rf /var/lib/apt/lists/*
COPY . .
USER node
EXPOSE 3000
HEALTHCHECK --interval=30s CMD curl -f http://localhost:3000/health || exit 1
CMD ["npm", "start"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // rm -rf command detected as potentially dangerous
		},
		{
			name: "Dockerfile without FROM",
			content: `WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 2, // missing USER + missing HEALTHCHECK
		},
		{
			name: "Dockerfile with latest tag",
			content: `FROM node:latest
WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // latest tag + missing USER + missing HEALTHCHECK
		},
		{
			name: "Dockerfile with root user",
			content: `FROM node:16
USER root
WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 2, // root user + missing healthcheck
		},
		{
			name: "Dockerfile with privileged operations",
			content: `FROM ubuntu:20.04
RUN apt-get update && apt-get install -y curl
RUN chmod 777 /tmp
RUN rm -rf /var/lib/apt/lists/*
CMD ["bash"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 6, // apt cleanup, chmod 777, rm -rf command, missing USER, missing healthcheck, apt without cleanup
		},
		{
			name: "Dockerfile with security issues",
			content: `FROM node:16
ENV API_KEY=sk_test_1234567890abcdef
RUN curl -o script.sh http://example.com/script.sh && bash script.sh
EXPOSE 21
USER root
CMD ["npm", "start"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // env secret, unverified download, insecure port, root user, missing healthcheck
		},
		{
			name: "Dockerfile with performance issues",
			content: `FROM node:16
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git
RUN apt-get install -y vim
RUN apt-get install -y wget
COPY package.json .
COPY package-lock.json .
COPY src/ ./src/
COPY public/ ./public/
COPY config/ ./config/
CMD ["npm", "start"]`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 8, // too many RUNs, too many COPYs, missing USER, missing healthcheck, multiple apt without cleanup warnings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "Dockerfile")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result, err := validator.ValidateDockerfile(filePath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ParseDockerfile(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name             string
		content          string
		expectedCommands []string
	}{
		{
			name: "simple Dockerfile",
			content: `FROM node:16
WORKDIR /app
COPY . .
CMD ["npm", "start"]`,
			expectedCommands: []string{"FROM", "WORKDIR", "COPY", "CMD"},
		},
		{
			name: "Dockerfile with comments and empty lines",
			content: `# Base image
FROM node:16

# Set working directory
WORKDIR /app

# Copy files
COPY . .

# Start application
CMD ["npm", "start"]`,
			expectedCommands: []string{"FROM", "WORKDIR", "COPY", "CMD"},
		},
		{
			name: "Dockerfile with line continuations",
			content: `FROM node:16
RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*
CMD ["npm", "start"]`,
			expectedCommands: []string{"FROM", "RUN", "APT-GET", "RM", "CMD"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instructions := validator.parseDockerfile(tt.content)

			assert.Equal(t, len(tt.expectedCommands), len(instructions))
			for i, expectedCmd := range tt.expectedCommands {
				assert.Equal(t, expectedCmd, instructions[i].Command)
			}
		})
	}
}

func TestDockerValidator_ValidateFromInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "specific version tag",
			instruction: DockerInstruction{
				Command: "FROM",
				Args:    "node:16.14.0",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "latest tag",
			instruction: DockerInstruction{
				Command: "FROM",
				Args:    "node:latest",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "no tag (implies latest)",
			instruction: DockerInstruction{
				Command: "FROM",
				Args:    "node",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "insecure registry",
			instruction: DockerInstruction{
				Command: "FROM",
				Args:    "http://registry.example.com/node:16",
				LineNum: 1,
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

			validator.validateFromInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateRunInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "apt-get with cleanup",
			instruction: DockerInstruction{
				Command: "RUN",
				Args:    "apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*",
				LineNum: 1,
			},
			expectWarnings: 2, // rm -rf command detected + apt-get without cleanup for first part
		},
		{
			name: "apt-get without cleanup",
			instruction: DockerInstruction{
				Command: "RUN",
				Args:    "apt-get update && apt-get install -y curl",
				LineNum: 1,
			},
			expectWarnings: 2, // apt-get without cleanup for both commands
		},
		{
			name: "yum without cleanup",
			instruction: DockerInstruction{
				Command: "RUN",
				Args:    "yum install -y curl",
				LineNum: 1,
			},
			expectWarnings: 2, // yum without cleanup + enhanced validation
		},
		{
			name: "dangerous command",
			instruction: DockerInstruction{
				Command: "RUN",
				Args:    "chmod 777 /tmp",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "unverified download",
			instruction: DockerInstruction{
				Command: "RUN",
				Args:    "curl -o script.sh http://example.com/script.sh",
				LineNum: 1,
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

			validator.validateRunInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateCopyAddInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "COPY specific files",
			instruction: DockerInstruction{
				Command: "COPY",
				Args:    "package*.json ./",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "ADD for simple copy",
			instruction: DockerInstruction{
				Command: "ADD",
				Args:    "package.json ./",
				LineNum: 1,
			},
			expectWarnings: 1, // prefer COPY
		},
		{
			name: "ADD for URL (valid use)",
			instruction: DockerInstruction{
				Command: "ADD",
				Args:    "http://example.com/file.tar.gz /tmp/",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "COPY entire context",
			instruction: DockerInstruction{
				Command: "COPY",
				Args:    ". /app",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "COPY with absolute source path",
			instruction: DockerInstruction{
				Command: "COPY",
				Args:    "/host/path /container/path",
				LineNum: 1,
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

			validator.validateCopyAddInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateUserInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "non-root user",
			instruction: DockerInstruction{
				Command: "USER",
				Args:    "node",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "root user by name",
			instruction: DockerInstruction{
				Command: "USER",
				Args:    "root",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "root user by UID",
			instruction: DockerInstruction{
				Command: "USER",
				Args:    "0",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "numeric UID",
			instruction: DockerInstruction{
				Command: "USER",
				Args:    "1000",
				LineNum: 1,
			},
			expectWarnings: 1, // numeric UID warning
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

			validator.validateUserInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateExposeInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "safe port",
			instruction: DockerInstruction{
				Command: "EXPOSE",
				Args:    "8080",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "HTTP port",
			instruction: DockerInstruction{
				Command: "EXPOSE",
				Args:    "80",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "FTP port",
			instruction: DockerInstruction{
				Command: "EXPOSE",
				Args:    "21",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "multiple ports with some insecure",
			instruction: DockerInstruction{
				Command: "EXPOSE",
				Args:    "80 443 8080",
				LineNum: 1,
			},
			expectWarnings: 1, // only port 80 is flagged
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

			validator.validateExposeInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateEnvInstruction(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instruction    DockerInstruction
		expectWarnings int
	}{
		{
			name: "safe environment variable",
			instruction: DockerInstruction{
				Command: "ENV",
				Args:    "NODE_ENV=production",
				LineNum: 1,
			},
			expectWarnings: 0,
		},
		{
			name: "potential secret",
			instruction: DockerInstruction{
				Command: "ENV",
				Args:    "API_KEY=sk_test_1234567890abcdef",
				LineNum: 1,
			},
			expectWarnings: 1,
		},
		{
			name: "password in env",
			instruction: DockerInstruction{
				Command: "ENV",
				Args:    "DB_PASSWORD=mysecretpassword",
				LineNum: 1,
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

			validator.validateEnvInstruction(tt.instruction, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidateSecurityPractices(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instructions   []DockerInstruction
		expectWarnings int
	}{
		{
			name: "has USER and HEALTHCHECK",
			instructions: []DockerInstruction{
				{Command: "FROM", Args: "node:16"},
				{Command: "USER", Args: "node"},
				{Command: "HEALTHCHECK", Args: "CMD curl -f http://localhost:3000/health"},
			},
			expectWarnings: 0,
		},
		{
			name: "missing USER",
			instructions: []DockerInstruction{
				{Command: "FROM", Args: "node:16"},
				{Command: "HEALTHCHECK", Args: "CMD curl -f http://localhost:3000/health"},
			},
			expectWarnings: 1,
		},
		{
			name: "missing HEALTHCHECK",
			instructions: []DockerInstruction{
				{Command: "FROM", Args: "node:16"},
				{Command: "USER", Args: "node"},
			},
			expectWarnings: 1,
		},
		{
			name: "missing both USER and HEALTHCHECK",
			instructions: []DockerInstruction{
				{Command: "FROM", Args: "node:16"},
				{Command: "WORKDIR", Args: "/app"},
			},
			expectWarnings: 2,
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

			validator.validateSecurityPractices(tt.instructions, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ValidatePerformancePractices(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name           string
		instructions   []DockerInstruction
		expectWarnings int
	}{
		{
			name: "reasonable number of instructions",
			instructions: []DockerInstruction{
				{Command: "FROM", Args: "node:16"},
				{Command: "RUN", Args: "apt-get update"},
				{Command: "COPY", Args: "package.json ."},
			},
			expectWarnings: 0,
		},
		{
			name: "too many RUN instructions",
			instructions: func() []DockerInstruction {
				instructions := []DockerInstruction{{Command: "FROM", Args: "node:16"}}
				for i := 0; i < 11; i++ {
					instructions = append(instructions, DockerInstruction{Command: "RUN", Args: "echo test"})
				}
				return instructions
			}(),
			expectWarnings: 1,
		},
		{
			name: "too many COPY instructions",
			instructions: func() []DockerInstruction {
				instructions := []DockerInstruction{{Command: "FROM", Args: "node:16"}}
				for i := 0; i < 6; i++ {
					instructions = append(instructions, DockerInstruction{Command: "COPY", Args: "file ."})
				}
				return instructions
			}(),
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

			validator.validatePerformancePractices(tt.instructions, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestDockerValidator_ContainsPotentialSecret(t *testing.T) {
	validator := NewDockerValidator()

	tests := []struct {
		name     string
		args     string
		expected bool
	}{
		{
			name:     "no secrets",
			args:     "NODE_ENV=production PORT=3000",
			expected: false,
		},
		{
			name:     "API key",
			args:     "API_KEY=sk_test_1234567890",
			expected: true,
		},
		{
			name:     "password",
			args:     "DB_PASSWORD=mysecret",
			expected: true,
		},
		{
			name:     "token",
			args:     "AUTH_TOKEN=abc123def456",
			expected: true,
		},
		{
			name:     "just variable name (no value)",
			args:     "API_KEY PASSWORD TOKEN",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.containsPotentialSecret(tt.args)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDockerValidator_FileReadError(t *testing.T) {
	validator := NewDockerValidator()

	// Test with non-existent file
	result, err := validator.ValidateDockerfile("/non/existent/Dockerfile")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "read_error", result.Errors[0].Type)
}
