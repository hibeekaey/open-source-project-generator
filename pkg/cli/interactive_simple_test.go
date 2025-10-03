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

// Simple tests for InteractiveManager that focus on core functionality
// without dealing with terminal/CI detection complexities

func TestNewInteractiveManager_Simple(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)

	// Test creation
	im := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	assert.NotNil(t, im)
	assert.Equal(t, cli, im.cli)
	assert.Equal(t, outputManager, im.outputManager)
	assert.Equal(t, flagHandler, im.flagHandler)
	assert.Equal(t, mockLogger, im.logger)
}

func TestInteractiveManager_ConfirmGeneration_NonInteractiveMode_Simple(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	_ = cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	im := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	config := &models.ProjectConfig{Name: "test"}

	// Test non-interactive mode (should auto-confirm)
	result := im.ConfirmGeneration(config)
	assert.True(t, result)
}

func TestInteractiveManager_PromptAdvancedOptions_NonInteractiveMode_Simple(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	_ = cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	im := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	// Test non-interactive mode (should return defaults)
	options, err := im.PromptAdvancedOptions()

	require.NoError(t, err)
	require.NotNil(t, options)
	assert.True(t, options.EnableSecurityScanning)
	assert.True(t, options.EnableQualityChecks)
	assert.False(t, options.EnablePerformanceOptimization)
	assert.True(t, options.GenerateDocumentation)
	assert.True(t, options.EnableCICD)
	assert.Equal(t, []string{"github-actions"}, options.CICDProviders)
	assert.False(t, options.EnableMonitoring)
}

func TestInteractiveManager_ConfirmAdvancedGeneration_NonInteractiveMode_Simple(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	_ = cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	im := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	config := &models.ProjectConfig{Name: "test"}
	options := &interfaces.AdvancedOptions{
		EnableSecurityScanning: true,
		EnableQualityChecks:    true,
	}

	// Test non-interactive mode (should auto-confirm)
	result := im.ConfirmAdvancedGeneration(config, options)
	assert.True(t, result)
}

func TestInteractiveManager_SelectTemplateInteractively_NonInteractiveMode_Simple(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	_ = cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	im := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	filter := interfaces.TemplateFilter{}

	// Test non-interactive mode
	template, err := im.SelectTemplateInteractively(filter)

	assert.Nil(t, template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interactive template selection not available in non-interactive mode")
}

// Test the delegation methods in CLI
func TestCLI_InteractiveMethodDelegation(t *testing.T) {
	// Create mock dependencies
	mockLogger := &MockLogger{}
	mockLogger.On("Info", mock.Anything, mock.Anything).Maybe()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Maybe()

	cli := &CLI{rootCmd: &cobra.Command{}}
	cli.rootCmd.PersistentFlags().Bool("non-interactive", true, "")
	_ = cli.rootCmd.ParseFlags([]string{"--non-interactive"})

	outputManager := NewOutputManager(false, false, false, mockLogger)
	flagHandler := NewFlagHandler(cli, outputManager, mockLogger)
	interactiveManager := NewInteractiveManager(cli, outputManager, flagHandler, mockLogger)

	// Set up CLI with interactive manager
	cli.outputManager = outputManager
	cli.flagHandler = flagHandler
	cli.interactiveManager = interactiveManager

	// Test PromptProjectDetails delegation
	_, err := cli.PromptProjectDetails()
	assert.Error(t, err) // Should error in non-interactive mode
	assert.Contains(t, err.Error(), "interactive prompts not available in non-interactive mode")

	// Test ConfirmGeneration delegation
	config := &models.ProjectConfig{Name: "test"}
	result := cli.ConfirmGeneration(config)
	assert.True(t, result) // Should auto-confirm in non-interactive mode

	// Test PromptAdvancedOptions delegation
	options, err := cli.PromptAdvancedOptions()
	assert.NoError(t, err)
	assert.NotNil(t, options)
	assert.True(t, options.EnableSecurityScanning)

	// Test ConfirmAdvancedGeneration delegation
	result = cli.ConfirmAdvancedGeneration(config, options)
	assert.True(t, result) // Should auto-confirm in non-interactive mode

	// Test SelectTemplateInteractively delegation
	filter := interfaces.TemplateFilter{}
	template, err := cli.SelectTemplateInteractively(filter)
	assert.Nil(t, template)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interactive template selection not available in non-interactive mode")
}
