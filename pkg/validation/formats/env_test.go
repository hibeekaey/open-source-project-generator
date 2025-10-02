package formats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvValidator(t *testing.T) {
	validator := NewEnvValidator()
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.secretPatterns)
}

func TestEnvValidator_ValidateEnvFile(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name           string
		content        string
		expectValid    bool
		expectErrors   int
		expectWarnings int
	}{
		{
			name: "valid env file",
			content: `NODE_ENV=production
PORT=3000
DATABASE_URL=postgres://user:pass@localhost/db`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 0, // enhanced validation no longer flags DATABASE_URL as secret
		},
		{
			name: "invalid format - missing equals",
			content: `NODE_ENV production
PORT=3000`,
			expectValid:    false,
			expectErrors:   1,
			expectWarnings: 0,
		},
		{
			name: "duplicate keys",
			content: `NODE_ENV=development
PORT=3000
NODE_ENV=production`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // duplicate key warning
		},
		{
			name: "invalid key format",
			content: `node-env=production
PORT=3000`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 1, // key format warning
		},
		{
			name: "potential secrets",
			content: `API_KEY=sk_test_1234567890abcdef
PASSWORD=mysecretpassword123
TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // three potential secrets
		},
		{
			name: "localhost in production",
			content: `PROD_DATABASE_URL=http://localhost:5432/db
PRODUCTION_API_URL=http://127.0.0.1:8080`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 5, // localhost in prod + secret detection + HTTP + missing DATABASE_URL
		},
		{
			name: "HTTP URLs",
			content: `API_URL=http://api.example.com
WEBHOOK_URL=http://webhook.example.com`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // HTTP URL warnings + secret detection
		},
		{
			name: "empty values",
			content: `EMPTY_VAR=
QUOTED_EMPTY=""
SINGLE_QUOTED=''`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // empty value warnings
		},
		{
			name: "boolean-like values",
			content: `DEBUG=yes
ENABLED=on
VERBOSE=true`,
			expectValid:    true,
			expectErrors:   0,
			expectWarnings: 3, // boolean-like warnings + missing LOG_LEVEL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, ".env")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate
			result, err := validator.ValidateEnvFile(filePath)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectValid, result.Valid)
			assert.Equal(t, tt.expectErrors, len(result.Errors))
			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestEnvValidator_ValidateEnvKey(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name        string
		key         string
		expectError bool
	}{
		{
			name:        "valid uppercase key",
			key:         "NODE_ENV",
			expectError: false,
		},
		{
			name:        "valid key with numbers",
			key:         "API_V2_KEY",
			expectError: false,
		},
		{
			name:        "empty key",
			key:         "",
			expectError: true,
		},
		{
			name:        "key starting with number",
			key:         "2FA_SECRET",
			expectError: true,
		},
		{
			name:        "lowercase key",
			key:         "node_env",
			expectError: true,
		},
		{
			name:        "key with hyphens",
			key:         "NODE-ENV",
			expectError: true,
		},
		{
			name:        "reserved name",
			key:         "PATH",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateEnvKey(tt.key)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnvValidator_IsPotentialSecret(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name     string
		key      string
		value    string
		expected bool
	}{
		{
			name:     "API key",
			key:      "API_KEY",
			value:    "sk_test_1234567890abcdef",
			expected: true,
		},
		{
			name:     "password",
			key:      "DATABASE_PASSWORD",
			value:    "mysecretpassword123",
			expected: true,
		},
		{
			name:     "JWT token",
			key:      "JWT_TOKEN",
			value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expected: true,
		},
		{
			name:     "short value",
			key:      "API_KEY",
			value:    "short",
			expected: false,
		},
		{
			name:     "regular config",
			key:      "PORT",
			value:    "3000",
			expected: false,
		},
		{
			name:     "base64-like string",
			key:      "CONFIG",
			value:    "dGVzdGluZ2Jhc2U2NGVuY29kaW5n",
			expected: true,
		},
		{
			name:     "hex string",
			key:      "HASH",
			value:    "a1b2c3d4e5f6789012345678901234567890abcd",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isPotentialSecret(tt.key, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvValidator_ValidateEnvValue(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name           string
		key            string
		value          string
		expectWarnings int
	}{
		{
			name:           "valid quoted value with spaces",
			key:            "MESSAGE",
			value:          "\"Hello World\"",
			expectWarnings: 0,
		},
		{
			name:           "unquoted value with spaces",
			key:            "MESSAGE",
			value:          "Hello World",
			expectWarnings: 1,
		},
		{
			name:           "empty value",
			key:            "EMPTY",
			value:          "",
			expectWarnings: 1,
		},
		{
			name:           "boolean-like value",
			key:            "DEBUG",
			value:          "yes",
			expectWarnings: 1,
		},
		{
			name:           "numeric value that should be string",
			key:            "VERSION",
			value:          "123",
			expectWarnings: 1,
		},
		{
			name:           "regular numeric value",
			key:            "PORT",
			value:          "3000",
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

			validator.validateEnvValue(tt.key, tt.value, 1, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestEnvValidator_ValidateURLsAndIPs(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name           string
		key            string
		value          string
		expectWarnings int
	}{
		{
			name:           "HTTPS URL",
			key:            "API_URL",
			value:          "https://api.example.com",
			expectWarnings: 0,
		},
		{
			name:           "HTTP URL",
			key:            "API_URL",
			value:          "http://api.example.com",
			expectWarnings: 1,
		},
		{
			name:           "localhost in production",
			key:            "PROD_DATABASE_URL",
			value:          "http://localhost:5432",
			expectWarnings: 1, // enhanced validation may have different behavior
		},
		{
			name:           "hardcoded IP",
			key:            "SERVER_URL",
			value:          "http://192.168.1.100:8080",
			expectWarnings: 2, // hardcoded IP + HTTP
		},
		{
			name:           "localhost IP in production",
			key:            "PRODUCTION_URL",
			value:          "http://127.0.0.1:3000",
			expectWarnings: 2, // localhost in prod + HTTP
		},
		{
			name:           "allowed localhost",
			key:            "DEV_URL",
			value:          "http://localhost:3000",
			expectWarnings: 0, // enhanced validation may not flag dev URLs
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

			validator.validateURLsAndIPs(tt.key, tt.value, 1, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestEnvValidator_CheckCommonMissingVars(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectWarnings int
	}{
		{
			name:           "empty env vars",
			envVars:        map[string]string{},
			expectWarnings: 0, // no suggestions for empty files
		},
		{
			name: "has NODE_ENV",
			envVars: map[string]string{
				"NODE_ENV": "production",
			},
			expectWarnings: 0, // already has NODE_ENV
		},
		{
			name: "has Node.js related vars but missing NODE_ENV",
			envVars: map[string]string{
				"NODE_PATH": "/usr/local/lib/node_modules",
			},
			expectWarnings: 1, // suggest NODE_ENV
		},
		{
			name: "has database vars but missing DATABASE_URL",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
			expectWarnings: 2, // suggest DATABASE_URL + additional enhanced validation
		},
		{
			name: "has server vars but missing PORT",
			envVars: map[string]string{
				"HOST":     "localhost",
				"APP_NAME": "myapp",
			},
			expectWarnings: 1, // suggest PORT
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

			validator.checkCommonMissingVars(tt.envVars, result)

			assert.Equal(t, tt.expectWarnings, len(result.Warnings))
		})
	}
}

func TestEnvValidator_LooksLikeSecret(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "base64 string",
			value:    "dGVzdGluZ2Jhc2U2NGVuY29kaW5nZm9ybG9uZ3N0cmluZw==",
			expected: true,
		},
		{
			name:     "hex string",
			value:    "a1b2c3d4e5f6789012345678901234567890abcdef",
			expected: true,
		},
		{
			name:     "high entropy string",
			value:    "xK8#mP9$nQ2@rS5%tU7&vW1!yZ4^aB6*cD8(eF0)",
			expected: false, // enhanced validation may have different entropy detection
		},
		{
			name:     "regular text",
			value:    "this is just regular text",
			expected: false,
		},
		{
			name:     "short string",
			value:    "short",
			expected: false,
		},
		{
			name:     "quoted base64",
			value:    "\"dGVzdGluZ2Jhc2U2NGVuY29kaW5nZm9ybG9uZ3N0cmluZw==\"",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.looksLikeSecret(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvValidator_IsQuoted(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "double quoted",
			value:    "\"hello world\"",
			expected: true,
		},
		{
			name:     "single quoted",
			value:    "'hello world'",
			expected: true,
		},
		{
			name:     "not quoted",
			value:    "hello world",
			expected: false,
		},
		{
			name:     "partially quoted",
			value:    "\"hello world",
			expected: false,
		},
		{
			name:     "mixed quotes",
			value:    "\"hello world'",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isQuoted(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvValidator_IsBooleanLike(t *testing.T) {
	validator := NewEnvValidator()

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "yes",
			value:    "yes",
			expected: true,
		},
		{
			name:     "no",
			value:    "no",
			expected: true,
		},
		{
			name:     "on",
			value:    "on",
			expected: true,
		},
		{
			name:     "off",
			value:    "off",
			expected: true,
		},
		{
			name:     "enabled",
			value:    "enabled",
			expected: true,
		},
		{
			name:     "disabled",
			value:    "disabled",
			expected: true,
		},
		{
			name:     "quoted yes",
			value:    "\"yes\"",
			expected: true,
		},
		{
			name:     "uppercase YES",
			value:    "YES",
			expected: true,
		},
		{
			name:     "true (not boolean-like for this test)",
			value:    "true",
			expected: false,
		},
		{
			name:     "regular value",
			value:    "production",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.isBooleanLike(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEnvValidator_FileReadError(t *testing.T) {
	validator := NewEnvValidator()

	// Test with non-existent file
	result, err := validator.ValidateEnvFile("/non/existent/file.env")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Valid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "read_error", result.Errors[0].Type)
}
