package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCLIValidator(t *testing.T) {
	mockCLI := &CLI{}
	mockOutputManager := &OutputManager{}
	mockLogger := &MockLogger{}

	validator := NewCLIValidator(mockCLI, mockOutputManager, mockLogger)

	assert.NotNil(t, validator)
	assert.Equal(t, mockCLI, validator.cli)
	assert.Equal(t, mockOutputManager, validator.outputManager)
	assert.Equal(t, mockLogger, validator.logger)
}

func TestCLIValidator_IsValidProjectName(t *testing.T) {
	validator := &CLIValidator{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid simple name", "myproject", true},
		{"valid with hyphens", "my-project", true},
		{"valid with underscores", "my_project", true},
		{"valid with numbers", "project123", true},
		{"valid mixed", "my-project_123", true},
		{"empty string", "", false},
		{"with spaces", "my project", false},
		{"with special chars", "my@project", false},
		{"with dots", "my.project", false},
		{"too long", strings.Repeat("a", 101), false},
		{"single char", "a", true},
		{"max length", strings.Repeat("a", 100), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValidProjectName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCLIValidator_IsValidLicense(t *testing.T) {
	validator := &CLIValidator{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"MIT", "MIT", true},
		{"Apache-2.0", "Apache-2.0", true},
		{"GPL-3.0", "GPL-3.0", true},
		{"BSD-3-Clause", "BSD-3-Clause", true},
		{"case insensitive MIT", "mit", true},
		{"case insensitive Apache", "apache-2.0", true},
		{"invalid license", "INVALID", false},
		{"empty string", "", false},
		{"custom license", "MyCustomLicense", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValidLicense(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCLIValidator_IsValidTemplateName(t *testing.T) {
	validator := &CLIValidator{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid simple name", "template", true},
		{"valid with hyphens", "my-template", true},
		{"valid with underscores", "my_template", true},
		{"valid with dots", "my.template", true},
		{"valid with numbers", "template123", true},
		{"valid mixed", "my-template_123.v1", true},
		{"empty string", "", false},
		{"with spaces", "my template", false},
		{"with special chars", "my@template", false},
		{"with slashes", "my/template", false},
		{"too long", strings.Repeat("a", 51), false},
		{"single char", "a", true},
		{"max length", strings.Repeat("a", 50), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.IsValidTemplateName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCLIValidator_ValidateGenerateConfiguration(t *testing.T) {
	// For these tests, we'll test the validation logic directly
	// without mocking the CLI error creation methods
	validator := &CLIValidator{}

	t.Run("valid configuration", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:    "valid-project",
			License: "MIT",
		}

		// Test the validation logic directly
		assert.True(t, validator.IsValidProjectName(config.Name))
		assert.True(t, validator.IsValidLicense(config.License))
	})

	t.Run("invalid project name", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name: "invalid@name",
		}

		assert.False(t, validator.IsValidProjectName(config.Name))
	})

	t.Run("invalid license", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:    "valid-project",
			License: "INVALID",
		}

		assert.True(t, validator.IsValidProjectName(config.Name))
		assert.False(t, validator.IsValidLicense(config.License))
	})
}

func TestCLIValidator_ValidateGenerateOptions(t *testing.T) {
	validator := &CLIValidator{cli: &CLI{}}

	t.Run("valid options", func(t *testing.T) {
		options := interfaces.GenerateOptions{
			OutputPath: "./my-project",
			Templates:  []string{"go-gin", "nextjs-app"},
		}

		err := validator.ValidateGenerateOptions(options)
		assert.NoError(t, err)
	})

	t.Run("invalid output path characters", func(t *testing.T) {
		options := interfaces.GenerateOptions{
			OutputPath: "invalid<path>",
		}

		err := validator.ValidateGenerateOptions(options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("empty template name", func(t *testing.T) {
		options := interfaces.GenerateOptions{
			Templates: []string{"valid-template", ""},
		}

		err := validator.ValidateGenerateOptions(options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("invalid template name", func(t *testing.T) {
		options := interfaces.GenerateOptions{
			Templates: []string{"invalid@template"},
		}

		err := validator.ValidateGenerateOptions(options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("conflicting offline and update versions", func(t *testing.T) {
		options := interfaces.GenerateOptions{
			Offline:        true,
			UpdateVersions: true,
		}

		err := validator.ValidateGenerateOptions(options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func TestCLIValidator_CheckWritePermissions(t *testing.T) {
	validator := &CLIValidator{cli: &CLI{}}

	t.Run("writable directory", func(t *testing.T) {
		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "cli-validator-test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		err = validator.CheckWritePermissions(tempDir)
		assert.NoError(t, err)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		nonExistentDir := "/non/existent/directory"
		err := validator.CheckWritePermissions(nonExistentDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no write permission")
	})

	t.Run("read-only directory", func(t *testing.T) {
		// Create a temporary directory and make it read-only
		tempDir, err := os.MkdirTemp("", "cli-validator-test-readonly")
		require.NoError(t, err)
		defer func() {
			// Restore write permissions before cleanup
			os.Chmod(tempDir, 0755)
			os.RemoveAll(tempDir)
		}()

		// Make directory read-only
		err = os.Chmod(tempDir, 0444)
		require.NoError(t, err)

		err = validator.CheckWritePermissions(tempDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no write permission")
	})
}

func TestCLIValidator_PerformPreGenerationChecks(t *testing.T) {
	// Create a properly initialized CLI instance for testing
	mockLogger := &MockLogger{}
	cli := &CLI{}
	cli.outputManager = NewOutputManager(false, false, false, mockLogger)
	validator := &CLIValidator{cli: cli}

	t.Run("new directory creation", func(t *testing.T) {
		// Create a temporary parent directory
		parentDir, err := os.MkdirTemp("", "cli-validator-parent")
		require.NoError(t, err)
		defer os.RemoveAll(parentDir)

		// Test creating a new subdirectory
		newDir := filepath.Join(parentDir, "new-project")
		options := interfaces.GenerateOptions{
			NonInteractive: true,
		}

		err = validator.PerformPreGenerationChecks(newDir, options)
		assert.NoError(t, err)

		// Verify directory was created
		_, err = os.Stat(newDir)
		assert.NoError(t, err)
	})

	t.Run("existing directory with force", func(t *testing.T) {
		// Create a temporary directory with some content
		tempDir, err := os.MkdirTemp("", "cli-validator-existing")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create a file in the directory
		testFile := filepath.Join(tempDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)

		options := interfaces.GenerateOptions{
			Force:          true,
			NonInteractive: true,
		}

		err = validator.PerformPreGenerationChecks(tempDir, options)
		assert.NoError(t, err)

		// Verify the directory still exists but is empty
		_, err = os.Stat(tempDir)
		assert.NoError(t, err)

		// Verify the test file was removed
		_, err = os.Stat(testFile)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("permission denied on parent directory", func(t *testing.T) {
		// This test is platform-specific and may not work on all systems
		// We'll test with a clearly invalid path
		invalidPath := "/root/invalid/path/that/should/not/exist"
		options := interfaces.GenerateOptions{
			NonInteractive: true,
		}

		err := validator.PerformPreGenerationChecks(invalidPath, options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unable to create output directory")
	})
}

// Benchmark tests for performance validation
func BenchmarkCLIValidator_IsValidProjectName(b *testing.B) {
	validator := &CLIValidator{}
	testName := "my-valid-project-name"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValidProjectName(testName)
	}
}

func BenchmarkCLIValidator_IsValidTemplateName(b *testing.B) {
	validator := &CLIValidator{}
	testName := "my-valid-template.name"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValidTemplateName(testName)
	}
}

func BenchmarkCLIValidator_IsValidLicense(b *testing.B) {
	validator := &CLIValidator{}
	testLicense := "MIT"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.IsValidLicense(testLicense)
	}
}
