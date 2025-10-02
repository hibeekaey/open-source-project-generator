package cli

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Integration tests for InteractiveManager functionality
// These tests verify the integration between CLI and InteractiveManager

func TestInteractiveManager_Integration_NonInteractiveMode(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Create CLI with non-interactive mode enabled
	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	interactiveManager := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	// Set up CLI with all components
	cli.outputManager = outputManager
	cli.flagHandler = flagHandler
	cli.interactiveManager = interactiveManager

	t.Run("PromptProjectDetails should error in non-interactive mode", func(t *testing.T) {
		config, err := cli.PromptProjectDetails()
		assert.Nil(t, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Interactive prompts not available in non-interactive mode")
	})

	t.Run("ConfirmGeneration should auto-confirm in non-interactive mode", func(t *testing.T) {
		config := &models.ProjectConfig{Name: "test-project"}
		result := cli.ConfirmGeneration(config)
		assert.True(t, result)
	})

	t.Run("PromptAdvancedOptions should return defaults in non-interactive mode", func(t *testing.T) {
		options, err := cli.PromptAdvancedOptions()
		require.NoError(t, err)
		require.NotNil(t, options)

		// Verify default values
		assert.True(t, options.EnableSecurityScanning)
		assert.True(t, options.EnableQualityChecks)
		assert.False(t, options.EnablePerformanceOptimization)
		assert.True(t, options.GenerateDocumentation)
		assert.True(t, options.EnableCICD)
		assert.Equal(t, []string{"github-actions"}, options.CICDProviders)
		assert.False(t, options.EnableMonitoring)
	})

	t.Run("ConfirmAdvancedGeneration should auto-confirm in non-interactive mode", func(t *testing.T) {
		config := &models.ProjectConfig{Name: "test-project"}
		options := &interfaces.AdvancedOptions{
			EnableSecurityScanning: true,
			EnableQualityChecks:    true,
		}
		result := cli.ConfirmAdvancedGeneration(config, options)
		assert.True(t, result)
	})

	t.Run("SelectTemplateInteractively should error in non-interactive mode", func(t *testing.T) {
		filter := interfaces.TemplateFilter{}
		template, err := cli.SelectTemplateInteractively(filter)
		assert.Nil(t, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "interactive template selection not available in non-interactive mode")
	})
}

func TestInteractiveManager_Integration_ComponentInitialization(t *testing.T) {
	// Test that the InteractiveManager is properly initialized in a CLI instance
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Create CLI manually to test component initialization
	cli := &CLI{rootCmd: &cobra.Command{}}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	interactiveManager := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	// Set up CLI with all components
	cli.outputManager = outputManager
	cli.flagHandler = flagHandler
	cli.interactiveManager = interactiveManager

	// Verify that InteractiveManager is properly initialized
	assert.NotNil(t, cli.interactiveManager)
	assert.NotNil(t, cli.interactiveManager.cli)
	assert.NotNil(t, cli.interactiveManager.outputManager)
	assert.NotNil(t, cli.interactiveManager.flagHandler)
	assert.NotNil(t, cli.interactiveManager.logger)

	// Verify that the interactive manager is properly wired
	assert.Equal(t, cli, cli.interactiveManager.cli)
	assert.Equal(t, cli.outputManager, cli.interactiveManager.outputManager)
	assert.Equal(t, cli.flagHandler, cli.interactiveManager.flagHandler)
	assert.Equal(t, mockLogger, cli.interactiveManager.logger)
}

func TestInteractiveManager_Integration_MethodDelegation(t *testing.T) {
	// Test that CLI methods properly delegate to InteractiveManager
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Warn", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything, mock.Anything).Maybe()

	// Create CLI with non-interactive mode for predictable testing
	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	interactiveManager := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	// Set up CLI with all components
	cli.outputManager = outputManager
	cli.flagHandler = flagHandler
	cli.interactiveManager = interactiveManager

	t.Run("All interactive methods should delegate properly", func(t *testing.T) {
		// Test PromptProjectDetails delegation
		config, err := cli.PromptProjectDetails()
		assert.Nil(t, config)
		assert.Error(t, err)

		// Test ConfirmGeneration delegation
		testConfig := &models.ProjectConfig{Name: "test"}
		result := cli.ConfirmGeneration(testConfig)
		assert.True(t, result)

		// Test PromptAdvancedOptions delegation
		options, err := cli.PromptAdvancedOptions()
		assert.NoError(t, err)
		assert.NotNil(t, options)

		// Test ConfirmAdvancedGeneration delegation
		result = cli.ConfirmAdvancedGeneration(testConfig, options)
		assert.True(t, result)

		// Test SelectTemplateInteractively delegation
		filter := interfaces.TemplateFilter{}
		template, err := cli.SelectTemplateInteractively(filter)
		assert.Nil(t, template)
		assert.Error(t, err)
	})
}

// Integration tests focus on the InteractiveManager component integration
// Mock implementations are defined in other test files to avoid conflicts
