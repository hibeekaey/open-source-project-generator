package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCLI is a mock implementation of CLI for testing
type MockCLI struct {
	mock.Mock
}

// Mock version command for testing - removed to avoid redeclaration

func TestNewCommandRegistry(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Verify registry is created correctly
	assert.NotNil(t, registry)
	assert.Equal(t, mockCLI, registry.cli)
}

func TestCommandRegistry_RegisterAllCommands(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Register all commands
	registry.RegisterAllCommands()

	// Verify that commands were added to the root command
	commands := mockCLI.rootCmd.Commands()
	assert.True(t, len(commands) > 0, "Commands should be registered")

	// Check for specific commands
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Use] = true
	}

	expectedCommands := []string{
		"generate [flags]",
		"validate [path] [flags]",
		"audit [path] [flags]",
		"version [flags]",
		"config",
		"list-templates [flags]",
		"template",
		"update [flags]",
		"cache <command> [flags]",
		"logs [flags]",
	}

	for _, expectedCmd := range expectedCommands {
		assert.True(t, commandNames[expectedCmd], "Command '%s' should be registered", expectedCmd)
	}
}

func TestCommandRegistry_setupGenerateCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup generate command
	registry.setupGenerateCommand()

	// Find the generate command
	var generateCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "generate [flags]" {
			generateCmd = cmd
			break
		}
	}

	// Verify generate command was created
	assert.NotNil(t, generateCmd, "Generate command should be created")
	assert.Equal(t, "generate [flags]", generateCmd.Use)
	assert.Equal(t, "Generate a new project from templates with modern best practices", generateCmd.Short)

	// Verify flags are set up
	flags := generateCmd.Flags()
	assert.NotNil(t, flags.Lookup("config"), "Should have config flag")
	assert.NotNil(t, flags.Lookup("output"), "Should have output flag")
	assert.NotNil(t, flags.Lookup("dry-run"), "Should have dry-run flag")
	assert.NotNil(t, flags.Lookup("offline"), "Should have offline flag")
	assert.NotNil(t, flags.Lookup("template"), "Should have template flag")
	assert.NotNil(t, flags.Lookup("force"), "Should have force flag")
}

func TestCommandRegistry_setupValidateCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup validate command
	registry.setupValidateCommand()

	// Find the validate command
	var validateCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "validate [path] [flags]" {
			validateCmd = cmd
			break
		}
	}

	// Verify validate command was created
	assert.NotNil(t, validateCmd, "Validate command should be created")
	assert.Equal(t, "validate [path] [flags]", validateCmd.Use)
	assert.Equal(t, "Validate project structure, configuration, and dependencies", validateCmd.Short)

	// Verify flags are set up
	flags := validateCmd.Flags()
	assert.NotNil(t, flags.Lookup("fix"), "Should have fix flag")
	assert.NotNil(t, flags.Lookup("report"), "Should have report flag")
	assert.NotNil(t, flags.Lookup("report-format"), "Should have report-format flag")
	assert.NotNil(t, flags.Lookup("strict"), "Should have strict flag")
}

func TestCommandRegistry_setupAuditCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup audit command
	registry.setupAuditCommand()

	// Find the audit command
	var auditCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "audit [path] [flags]" {
			auditCmd = cmd
			break
		}
	}

	// Verify audit command was created
	assert.NotNil(t, auditCmd, "Audit command should be created")
	assert.Equal(t, "audit [path] [flags]", auditCmd.Use)
	assert.Equal(t, "Comprehensive security, quality, and compliance auditing", auditCmd.Short)

	// Verify flags are set up
	flags := auditCmd.Flags()
	assert.NotNil(t, flags.Lookup("security"), "Should have security flag")
	assert.NotNil(t, flags.Lookup("quality"), "Should have quality flag")
	assert.NotNil(t, flags.Lookup("licenses"), "Should have licenses flag")
	assert.NotNil(t, flags.Lookup("performance"), "Should have performance flag")
	assert.NotNil(t, flags.Lookup("detailed"), "Should have detailed flag")
}

func TestCommandRegistry_setupConfigCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup config command
	registry.setupConfigCommand()

	// Find the config command
	var configCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "config" {
			configCmd = cmd
			break
		}
	}

	// Verify config command was created
	assert.NotNil(t, configCmd, "Config command should be created")
	assert.Equal(t, "config", configCmd.Use)
	assert.Equal(t, "Manage saved project configurations", configCmd.Short)

	// Verify subcommands are set up
	subcommands := configCmd.Commands()
	assert.True(t, len(subcommands) >= 4, "Should have at least 4 subcommands")

	subcommandNames := make(map[string]bool)
	for _, subcmd := range subcommands {
		subcommandNames[subcmd.Use] = true
	}

	expectedSubcommands := []string{"list [flags]", "view <config-name> [flags]", "delete <config-name> [flags]", "export [config-name-or-file] [flags]", "import [flags]", "manage"}
	foundCount := 0
	for _, expectedSubcmd := range expectedSubcommands {
		if subcommandNames[expectedSubcmd] {
			foundCount++
		}
	}
	assert.True(t, foundCount >= 3, "Should have at least 3 expected subcommands, found %d", foundCount)
}

func TestCommandRegistry_setupTemplateCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup template command
	registry.setupTemplateCommand()

	// Find the template command
	var templateCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "template" {
			templateCmd = cmd
			break
		}
	}

	// Verify template command was created
	assert.NotNil(t, templateCmd, "Template command should be created")
	assert.Equal(t, "template", templateCmd.Use)
	assert.Equal(t, "Template management operations", templateCmd.Short)

	// Verify subcommands are set up
	subcommands := templateCmd.Commands()
	assert.True(t, len(subcommands) >= 2, "Should have at least 2 subcommands")

	subcommandNames := make(map[string]bool)
	for _, subcmd := range subcommands {
		subcommandNames[subcmd.Use] = true
	}

	expectedSubcommands := []string{"info <template-name> [flags]", "validate <template-path> [flags]"}
	for _, expectedSubcmd := range expectedSubcommands {
		assert.True(t, subcommandNames[expectedSubcmd], "Subcommand '%s' should be registered", expectedSubcmd)
	}
}

func TestCommandRegistry_setupCacheCommand(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Setup cache command
	registry.setupCacheCommand()

	// Find the cache command
	var cacheCmd *cobra.Command
	for _, cmd := range mockCLI.rootCmd.Commands() {
		if cmd.Use == "cache <command> [flags]" {
			cacheCmd = cmd
			break
		}
	}

	// Verify cache command was created
	assert.NotNil(t, cacheCmd, "Cache command should be created")
	assert.Equal(t, "cache <command> [flags]", cacheCmd.Use)
	assert.Equal(t, "Manage local cache for offline mode and performance", cacheCmd.Short)

	// Verify subcommands are set up
	subcommands := cacheCmd.Commands()
	assert.True(t, len(subcommands) >= 5, "Should have at least 5 subcommands")

	subcommandNames := make(map[string]bool)
	for _, subcmd := range subcommands {
		subcommandNames[subcmd.Use] = true
	}

	expectedSubcommands := []string{"show", "clear", "clean", "validate", "repair", "offline"}
	for _, expectedSubcmd := range expectedSubcommands {
		assert.True(t, subcommandNames[expectedSubcmd], "Subcommand '%s' should be registered", expectedSubcmd)
	}
}

func TestCommandRegistry_CommandFlagValidation(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Register all commands
	registry.RegisterAllCommands()

	// Test that all commands have proper flag setup
	commands := mockCLI.rootCmd.Commands()
	for _, cmd := range commands {
		// Verify that each command has a proper Use field
		assert.NotEmpty(t, cmd.Use, "Command should have a Use field")
		assert.NotEmpty(t, cmd.Short, "Command should have a Short description")

		// Verify that commands with arguments have proper Args validation
		if cmd.Use == "template info <template-name> [flags]" ||
			cmd.Use == "template validate <template-path> [flags]" ||
			cmd.Use == "config save <name> [flags]" ||
			cmd.Use == "config load <name> [flags]" ||
			cmd.Use == "config delete <name> [flags]" {
			assert.NotNil(t, cmd.Args, "Commands with required arguments should have Args validation")
		}
	}
}

func TestCommandRegistry_CommandExamples(t *testing.T) {
	// Create a mock CLI with root command
	mockCLI := &CLI{
		rootCmd: &cobra.Command{
			Use: "generator",
		},
	}

	// Create command registry
	registry := NewCommandRegistry(mockCLI)

	// Register all commands
	registry.RegisterAllCommands()

	// Test that main commands have examples
	commands := mockCLI.rootCmd.Commands()
	for _, cmd := range commands {
		if cmd.Use == "generate [flags]" ||
			cmd.Use == "validate [path] [flags]" ||
			cmd.Use == "audit [path] [flags]" ||
			cmd.Use == "list-templates [flags]" ||
			cmd.Use == "update [flags]" ||
			cmd.Use == "logs [flags]" {
			assert.NotEmpty(t, cmd.Example, "Main commands should have examples: %s", cmd.Use)
		}
	}
}

// Benchmark tests for command registration performance
func BenchmarkCommandRegistry_RegisterAllCommands(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mockCLI := &CLI{
			rootCmd: &cobra.Command{
				Use: "generator",
			},
		}
		registry := NewCommandRegistry(mockCLI)
		registry.RegisterAllCommands()
	}
}

func BenchmarkCommandRegistry_SetupGenerateCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mockCLI := &CLI{
			rootCmd: &cobra.Command{
				Use: "generator",
			},
		}
		registry := NewCommandRegistry(mockCLI)
		registry.setupGenerateCommand()
	}
}
