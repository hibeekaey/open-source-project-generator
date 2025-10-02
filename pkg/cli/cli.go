// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/handlers"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/interactive"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/utils"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
	"github.com/spf13/cobra"
)

// CLI implements the CLIInterface for comprehensive CLI operations.
type CLI struct {
	// Core dependencies
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	templateManager interfaces.TemplateManager
	cacheManager    interfaces.CacheManager
	versionManager  interfaces.VersionManager
	auditEngine     interfaces.AuditEngine
	logger          interfaces.Logger
	interactiveUI   interfaces.InteractiveUIInterface

	// Legacy CLI components (to be refactored)
	interactiveFlowManager *InteractiveFlowManager
	interactiveManager     *InteractiveManager
	outputManager          *OutputManager
	flagHandler            *FlagHandler
	commandRegistry        *CommandRegistry
	commandHandlers        *CommandHandlers
	cliValidator           *CLIValidator

	// New extracted components
	generateHandler    *handlers.GenerateHandler
	workflowHandler    *handlers.WorkflowHandler
	projectSetup       *interactive.ProjectSetup
	componentSelection *interactive.ComponentSelection
	inputValidator     *validation.InputValidator
	formatter          *utils.Formatter
	helper             *utils.Helper

	// CLI state
	generatorVersion string
	rootCmd          *cobra.Command
	exitCode         int
	gitCommit        string
	buildTime        string

	// Mode flags (for backward compatibility)
	verboseMode bool
	quietMode   bool
	debugMode   bool
}

// NewCLI creates a new CLI instance with all required dependencies.
func NewCLI(
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	templateManager interfaces.TemplateManager,
	cacheManager interfaces.CacheManager,
	versionManager interfaces.VersionManager,
	auditEngine interfaces.AuditEngine,
	logger interfaces.Logger,
	version string,
	gitCommit string,
	buildTime string,
) interfaces.CLIInterface {
	uiConfig := &ui.UIConfig{
		EnableColors: true, EnableUnicode: true, PageSize: 10, Timeout: 30 * time.Minute,
		AutoSave: true, ShowBreadcrumbs: true, ShowShortcuts: true, ConfirmOnQuit: true,
	}
	interactiveUI := ui.NewInteractiveUI(logger, uiConfig)

	cli := &CLI{
		configManager: configManager, validator: validator, templateManager: templateManager,
		cacheManager: cacheManager, versionManager: versionManager, auditEngine: auditEngine,
		logger: logger, interactiveUI: interactiveUI, generatorVersion: version,
		gitCommit: gitCommit, buildTime: buildTime,
	}

	// Initialize legacy components
	cli.outputManager = NewOutputManager(false, false, false, logger)
	cli.flagHandler = NewFlagHandler(cli, cli.outputManager, logger)
	cli.interactiveManager = NewInteractiveManager(cli, cli.outputManager, cli.flagHandler, logger)
	cli.interactiveFlowManager = NewInteractiveFlowManager(cli, templateManager, configManager, validator, logger, interactiveUI)
	cli.cliValidator = NewCLIValidator(cli, cli.outputManager, logger)
	cli.commandHandlers = NewCommandHandlers(cli)
	cli.commandRegistry = NewCommandRegistry(cli)

	// Initialize new extracted components
	outputAdapter := NewOutputAdapter(cli.outputManager)
	cli.generateHandler = handlers.NewGenerateHandler(cli, templateManager, configManager, validator, logger)
	cli.workflowHandler = handlers.NewWorkflowHandler(cli, cli.generateHandler, configManager, validator, logger)
	cli.projectSetup = interactive.NewProjectSetup(logger, outputAdapter)
	cli.componentSelection = interactive.NewComponentSelection(logger, outputAdapter)
	cli.inputValidator = validation.NewInputValidator(logger, outputAdapter)
	cli.formatter = utils.NewFormatter()
	cli.helper = utils.NewHelper()

	cli.setupCommands()
	return cli
}

// setupCommands initializes all CLI commands and their flags
func (c *CLI) setupCommands() {
	c.rootCmd = &cobra.Command{
		Use: "generator", Short: "Open Source Project Generator - Create production-ready projects with modern best practices",
		Long:         `Generate production-ready projects with modern best practices. Supports Go, Next.js, React, Android, iOS, Docker, Kubernetes, and Terraform.`,
		SilenceUsage: true, SilenceErrors: true,
	}
	c.flagHandler.SetupGlobalFlags(c.rootCmd)
	c.commandRegistry.RegisterAllCommands()
}

// Run executes the CLI application with the provided arguments
func (c *CLI) Run(args []string) error {
	c.rootCmd.SetArgs(args)
	c.rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error { return c.flagHandler.HandleGlobalFlags(cmd) }
	err := c.rootCmd.Execute()
	if err != nil {
		exitCode := c.handleError(err, "generator", args)
		if c.helper.DetectNonInteractiveMode(c.rootCmd) {
			os.Exit(exitCode)
		}
	}
	return err
}

// Core interface methods
func (c *CLI) GetExitCode() int                             { return c.exitCode }
func (c *CLI) SetExitCode(code int)                         { c.exitCode = code }
func (c *CLI) GetVersionManager() interfaces.VersionManager { return c.versionManager }
func (c *CLI) GetBuildInfo() (version, gitCommit, buildTime string) {
	return c.generatorVersion, c.gitCommit, c.buildTime
}

// Output methods - delegate to OutputManager
func (c *CLI) VerboseOutput(format string, args ...interface{}) {
	c.outputManager.VerboseOutput(format, args...)
}
func (c *CLI) DebugOutput(format string, args ...interface{}) {
	c.outputManager.DebugOutput(format, args...)
}
func (c *CLI) QuietOutput(format string, args ...interface{}) {
	c.outputManager.QuietOutput(format, args...)
}
func (c *CLI) ErrorOutput(format string, args ...interface{}) {
	c.outputManager.ErrorOutput(format, args...)
}
func (c *CLI) WarningOutput(format string, args ...interface{}) {
	c.outputManager.WarningOutput(format, args...)
}
func (c *CLI) SuccessOutput(format string, args ...interface{}) {
	c.outputManager.SuccessOutput(format, args...)
}
func (c *CLI) IsQuietMode() bool { return c.outputManager.IsQuietMode() }

// Color methods - delegate to ColorManager
func (c *CLI) Error(text string) string     { return c.outputManager.GetColorManager().Error(text) }
func (c *CLI) Warning(text string) string   { return c.outputManager.GetColorManager().Warning(text) }
func (c *CLI) Info(text string) string      { return c.outputManager.GetColorManager().Info(text) }
func (c *CLI) Success(text string) string   { return c.outputManager.GetColorManager().Success(text) }
func (c *CLI) Highlight(text string) string { return c.outputManager.GetColorManager().Highlight(text) }
func (c *CLI) Dim(text string) string       { return c.outputManager.GetColorManager().Dim(text) }

// Interactive methods - delegate to extracted components
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	return c.projectSetup.CollectProjectDetails()
}

func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	confirmed, err := c.projectSetup.ConfirmProjectDetails(config)
	if err != nil {
		c.logger.Error("Failed to get user confirmation", "error", err)
		return false
	}
	return confirmed
}

// Generation methods - delegate to handlers
func (c *CLI) GenerateFromConfig(configPath string, options interfaces.GenerateOptions) error {
	var config *models.ProjectConfig
	var err error
	if configPath != "" {
		config, err = c.workflowHandler.LoadConfigFromFile(configPath)
	} else {
		config, err = c.workflowHandler.LoadConfigFromEnvironment()
	}
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	return c.workflowHandler.ExecuteGenerationWorkflow(config, options)
}

func (c *CLI) GenerateWithAdvancedOptions(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	generateOptions := interfaces.GenerateOptions{
		OutputPath: options.OutputPath, Force: options.Force, DryRun: options.DryRun,
		Offline: options.Offline, UpdateVersions: options.UpdateVersions, SkipValidation: options.SkipValidation,
		Minimal: options.Minimal, NonInteractive: options.NonInteractive,
	}
	return c.workflowHandler.ExecuteGenerationWorkflow(config, generateOptions)
}

// Validation methods
func (c *CLI) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	if c.validator == nil {
		return nil, fmt.Errorf("validation engine not initialized")
	}
	return &interfaces.ValidationResult{Valid: true, Issues: []interfaces.ValidationIssue{}, Warnings: []interfaces.ValidationIssue{}}, nil
}

// Interface compatibility methods
func (c *CLI) CreateAuditError(message string, score float64) error {
	structuredErr := c.createAuditError(message, score)
	return fmt.Errorf("%s", structuredErr.Error())
}
func (c *CLI) OutputMachineReadable(data interface{}, format string) error {
	output, err := c.formatter.FormatMachineReadable(data, format)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return nil
}
func (c *CLI) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
	c.outputManager.PerformanceOutput(operation, duration, metrics)
}
func (c *CLI) StartOperationWithOutput(operation string, description string) *interfaces.OperationContext {
	return c.outputManager.StartOperationWithOutput(operation, description)
}
func (c *CLI) FinishOperationWithOutput(ctx *interfaces.OperationContext, operation string, description string) {
	c.outputManager.FinishOperationWithOutput(ctx, operation, description)
}
func (c *CLI) FinishOperationWithError(ctx *interfaces.OperationContext, operation string, err error) {
	c.outputManager.FinishOperationWithError(ctx, operation, err)
}

// Audit methods
func (c *CLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	return nil, fmt.Errorf("not implemented")
}
func (c *CLI) AuditProjectAdvanced(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	return nil, fmt.Errorf("not implemented")
}
func (c *CLI) ValidateProjectAdvanced(path string, options *interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	return nil, fmt.Errorf("not implemented")
}

// validateAndNormalizeMode validates and normalizes the mode string
func (c *CLI) validateAndNormalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "interactive", "i":
		return "interactive"
	case "non-interactive", "noninteractive", "ni", "auto":
		return "non-interactive"
	case "config-file", "config", "file", "cf":
		return "config-file"
	default:
		c.WarningOutput("Unknown mode '%s', defaulting to interactive", mode)
		return "interactive"
	}
}
