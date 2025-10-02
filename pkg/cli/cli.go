// Package cli provides comprehensive command-line interface functionality for the
// Open Source Project Generator.
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/handlers"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/interactive"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/utils"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/performance"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui"
	"github.com/cuesoftinc/open-source-project-generator/pkg/workflow"
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

	// Workflow management
	workflowManager interfaces.WorkflowManager
	generator       interfaces.FileSystemGenerator

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

	// Enhanced error handling
	errorIntegration *ErrorIntegration

	// Performance optimization
	optimizer         *performance.CommandOptimizer
	lazyLoaderManager *performance.LazyLoaderManager

	// System monitoring
	systemMonitor *performance.SystemMonitor
	dashboard     *performance.Dashboard

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

	// Create filesystem generator
	generator := filesystem.NewGenerator()

	// Create workflow manager
	workflowManager := workflow.NewManager(
		generator,
		templateManager,
		validator,
		auditEngine,
		cacheManager,
		configManager,
		logger,
	)

	cli := &CLI{
		configManager: configManager, validator: validator, templateManager: templateManager,
		cacheManager: cacheManager, versionManager: versionManager, auditEngine: auditEngine,
		logger: logger, interactiveUI: interactiveUI, generatorVersion: version,
		gitCommit: gitCommit, buildTime: buildTime,
		workflowManager: workflowManager, generator: generator,
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

	// Initialize performance optimization
	cli.optimizer = performance.NewCommandOptimizer(cacheManager, logger)
	cli.lazyLoaderManager = performance.NewLazyLoaderManager()

	// Initialize system monitoring
	cli.systemMonitor = performance.NewSystemMonitor(cacheManager, logger)
	cli.dashboard = performance.NewDashboard(cli.systemMonitor, cli.optimizer.GetMetrics())

	// Register health checkers
	cli.registerHealthCheckers()

	// Register lazy loaders for common operations
	cli.registerLazyLoaders()

	cli.setupCommands()

	// Initialize enhanced error handling after commands are set up
	if errorIntegration, err := NewErrorIntegration(cli, logger); err == nil {
		cli.errorIntegration = errorIntegration
	} else {
		// Log error but continue with basic error handling
		if logger != nil {
			logger.Warn("Failed to initialize enhanced error handling: %v", err)
		}
	}

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
func (c *CLI) GetExitCode() int                               { return c.exitCode }
func (c *CLI) SetExitCode(code int)                           { c.exitCode = code }
func (c *CLI) GetVersionManager() interfaces.VersionManager   { return c.versionManager }
func (c *CLI) GetWorkflowManager() interfaces.WorkflowManager { return c.workflowManager }
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

// Cache management methods
func (c *CLI) ShowCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	stats, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Display cache information
	c.QuietOutput("📦 Cache Status")
	c.QuietOutput("==============")
	c.QuietOutput("Location: %s", c.cacheManager.GetLocation())
	c.QuietOutput("Total entries: %d", stats.TotalEntries)
	c.QuietOutput("Total size: %s", c.formatBytes(stats.TotalSize))
	c.QuietOutput("Hit rate: %.2f%%", stats.HitRate*100)
	c.QuietOutput("Last accessed: %s", stats.LastAccessed.Format(time.RFC3339))
	c.QuietOutput("Last modified: %s", stats.LastModified.Format(time.RFC3339))
	c.QuietOutput("Created at: %s", stats.CreatedAt.Format(time.RFC3339))

	// Show offline mode status
	if c.cacheManager.IsOfflineMode() {
		c.QuietOutput("Offline mode: %s", c.Success("enabled"))
	} else {
		c.QuietOutput("Offline mode: %s", c.Dim("disabled"))
	}

	return nil
}

func (c *CLI) ClearCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.Clear()
}

func (c *CLI) CleanCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.Clean()
}

func (c *CLI) ValidateCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.ValidateCache()
}

func (c *CLI) RepairCache() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.RepairCache()
}

func (c *CLI) EnableOfflineMode() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	err := c.cacheManager.EnableOfflineMode()
	if err != nil {
		return fmt.Errorf("failed to enable offline mode: %w", err)
	}

	c.QuietOutput("✅ %s", c.Success("Offline mode enabled"))
	return nil
}

func (c *CLI) DisableOfflineMode() error {
	if c.cacheManager == nil {
		return fmt.Errorf("cache manager not initialized")
	}

	err := c.cacheManager.DisableOfflineMode()
	if err != nil {
		return fmt.Errorf("failed to disable offline mode: %w", err)
	}

	c.QuietOutput("✅ %s", c.Success("Offline mode disabled"))
	return nil
}

func (c *CLI) GetCacheStats() (*interfaces.CacheStats, error) {
	if c.cacheManager == nil {
		return nil, fmt.Errorf("cache manager not initialized")
	}

	return c.cacheManager.GetStats()
}

func (c *CLI) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Audit methods
func (c *CLI) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	if c.auditEngine == nil {
		return nil, fmt.Errorf("🚫 %s %s",
			c.Error("Audit engine not initialized."),
			c.Info("This is an internal error - please report this issue"))
	}

	c.VerboseOutput("🔍 Starting comprehensive project audit for: %s", path)

	// Create audit result
	result := &interfaces.AuditResult{
		ProjectPath:     path,
		AuditTime:       time.Now(),
		OverallScore:    0.0,
		Recommendations: []string{},
	}

	var totalScore float64
	var auditCount int

	// Security audit
	if options.Security {
		c.VerboseOutput("🔒 Running security audit...")
		securityResult, err := c.auditEngine.AuditSecurity(path)
		if err != nil {
			c.WarningOutput("⚠️  Security audit failed: %v", err)
		} else {
			result.Security = securityResult
			totalScore += securityResult.Score
			auditCount++
			c.VerboseOutput("✅ Security audit completed (Score: %.1f/100)", securityResult.Score)
		}
	}

	// Quality audit
	if options.Quality {
		c.VerboseOutput("✨ Running quality audit...")
		qualityResult, err := c.auditEngine.AuditCodeQuality(path)
		if err != nil {
			c.WarningOutput("⚠️  Quality audit failed: %v", err)
		} else {
			result.Quality = qualityResult
			totalScore += qualityResult.Score
			auditCount++
			c.VerboseOutput("✅ Quality audit completed (Score: %.1f/100)", qualityResult.Score)
		}
	}

	// License audit
	if options.Licenses {
		c.VerboseOutput("📄 Running license audit...")
		licenseResult, err := c.auditEngine.AuditLicenses(path)
		if err != nil {
			c.WarningOutput("⚠️  License audit failed: %v", err)
		} else {
			result.Licenses = licenseResult
			totalScore += licenseResult.Score
			auditCount++
			c.VerboseOutput("✅ License audit completed (Score: %.1f/100)", licenseResult.Score)
		}
	}

	// Performance audit
	if options.Performance {
		c.VerboseOutput("⚡ Running performance audit...")
		performanceResult, err := c.auditEngine.AuditPerformance(path)
		if err != nil {
			c.WarningOutput("⚠️  Performance audit failed: %v", err)
		} else {
			result.Performance = performanceResult
			totalScore += performanceResult.Score
			auditCount++
			c.VerboseOutput("✅ Performance audit completed (Score: %.1f/100)", performanceResult.Score)
		}
	}

	// Calculate overall score
	if auditCount > 0 {
		result.OverallScore = totalScore / float64(auditCount)
	}

	// Generate overall recommendations
	result.Recommendations = c.generateAuditRecommendations(result)

	c.VerboseOutput("🎯 Audit completed with overall score: %.1f/100", result.OverallScore)

	return result, nil
}
func (c *CLI) AuditProjectAdvanced(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	if options == nil {
		options = &interfaces.AuditOptions{
			Security:     true,
			Quality:      true,
			Licenses:     true,
			Performance:  true,
			OutputFormat: "text",
			Detailed:     false,
		}
	}
	return c.AuditProject(path, *options)
}

// generateAuditRecommendations generates overall recommendations based on audit results
func (c *CLI) generateAuditRecommendations(result *interfaces.AuditResult) []string {
	recommendations := []string{}

	// Security recommendations
	if result.Security != nil && result.Security.Score < 80 {
		if len(result.Security.Vulnerabilities) > 0 {
			recommendations = append(recommendations, "Update dependencies to fix security vulnerabilities")
		}
		if len(result.Security.PolicyViolations) > 0 {
			recommendations = append(recommendations, "Address security policy violations")
		}
		if result.Security.Score < 50 {
			recommendations = append(recommendations, "Consider implementing additional security measures")
		}
	}

	// Quality recommendations
	if result.Quality != nil && result.Quality.Score < 80 {
		if len(result.Quality.CodeSmells) > 0 {
			recommendations = append(recommendations, "Refactor code to address quality issues")
		}
		if result.Quality.TestCoverage < 70 {
			recommendations = append(recommendations, "Increase test coverage to improve code quality")
		}
		if result.Quality.Score < 50 {
			recommendations = append(recommendations, "Consider code review and refactoring")
		}
	}

	// License recommendations
	if result.Licenses != nil && !result.Licenses.Compatible {
		recommendations = append(recommendations, "Review and resolve license compatibility issues")
		if len(result.Licenses.Conflicts) > 0 {
			recommendations = append(recommendations, "Address license conflicts in dependencies")
		}
	}

	// Performance recommendations
	if result.Performance != nil && result.Performance.Score < 80 {
		if result.Performance.BundleSize > 1024*1024 { // 1MB
			recommendations = append(recommendations, "Optimize bundle size to improve performance")
		}
		if len(result.Performance.Issues) > 0 {
			recommendations = append(recommendations, "Address performance issues identified in the audit")
		}
	}

	// Overall recommendations
	if result.OverallScore < 70 {
		recommendations = append(recommendations, "Consider running regular audits to maintain project health")
	}

	return recommendations
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

// Enhanced error handling and diagnostics methods

// ShowLogs displays recent log entries
func (c *CLI) ShowLogs() error {
	if c.errorIntegration == nil {
		return fmt.Errorf("enhanced error handling not initialized")
	}

	stats := c.errorIntegration.GetStatistics()
	if stats == nil {
		c.QuietOutput("No logging statistics available")
		return nil
	}

	c.QuietOutput("📊 Logging Statistics")
	c.QuietOutput("===================")
	c.QuietOutput("Total Errors: %d", stats.TotalErrors)
	c.QuietOutput("Recovery Rate: %.1f%%", stats.RecoveryRate)

	if len(stats.ErrorsByCategory) > 0 {
		c.QuietOutput("\nError Categories:")
		for category, count := range stats.ErrorsByCategory {
			percentage := float64(count) / float64(stats.TotalErrors) * 100
			c.QuietOutput("  %s: %d (%.1f%%)", category, count, percentage)
		}
	}

	return nil
}

// SetLogLevel sets the logging level
func (c *CLI) SetLogLevel(level string) error {
	if c.logger == nil {
		return fmt.Errorf("logger not initialized")
	}

	var logLevel int
	switch strings.ToLower(level) {
	case "debug":
		logLevel = 0
	case "info":
		logLevel = 1
	case "warn", "warning":
		logLevel = 2
	case "error":
		logLevel = 3
	case "fatal":
		logLevel = 4
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	c.logger.SetLevel(logLevel)
	c.VerboseOutput("Log level set to: %s", level)
	return nil
}

// GetLogLevel returns the current logging level
func (c *CLI) GetLogLevel() string {
	if c.logger == nil {
		return "unknown"
	}

	if c.logger.IsDebugEnabled() {
		return "debug"
	}
	if c.logger.IsInfoEnabled() {
		return "info"
	}
	return "warn"
}

// EnableVerboseLogging enables verbose logging with detailed output
func (c *CLI) EnableVerboseLogging() {
	c.verboseMode = true
	c.outputManager.SetVerboseMode(true)

	if c.errorIntegration != nil {
		c.errorIntegration.SetVerboseMode(true)
	}

	if c.logger != nil {
		c.logger.SetLevel(0) // Debug level
	}

	c.VerboseOutput("Verbose logging enabled")
}

// EnableDebugMode enables debug mode with comprehensive tracing
func (c *CLI) EnableDebugMode() {
	c.debugMode = true
	c.verboseMode = true
	c.outputManager.SetDebugMode(true)
	c.outputManager.SetVerboseMode(true)

	if c.errorIntegration != nil {
		c.errorIntegration.SetVerboseMode(true)
	}

	if c.logger != nil {
		c.logger.SetLevel(0) // Debug level
		c.logger.SetCallerInfo(true)
	}

	c.DebugOutput("Debug mode enabled with comprehensive tracing")
}

// GenerateDiagnosticReport generates a comprehensive diagnostic report
func (c *CLI) GenerateDiagnosticReport() string {
	var report strings.Builder

	report.WriteString("CLI Diagnostic Report\n")
	report.WriteString("====================\n\n")

	// Basic CLI information
	report.WriteString(fmt.Sprintf("Version: %s\n", c.generatorVersion))
	report.WriteString(fmt.Sprintf("Git Commit: %s\n", c.gitCommit))
	report.WriteString(fmt.Sprintf("Build Time: %s\n", c.buildTime))
	report.WriteString(fmt.Sprintf("Verbose Mode: %t\n", c.verboseMode))
	report.WriteString(fmt.Sprintf("Quiet Mode: %t\n", c.quietMode))
	report.WriteString(fmt.Sprintf("Debug Mode: %t\n", c.debugMode))

	// Error handling statistics
	if c.errorIntegration != nil {
		report.WriteString("\n")
		report.WriteString(c.errorIntegration.GenerateErrorReport())
	}

	// Component status
	report.WriteString("\nComponent Status:\n")
	report.WriteString(fmt.Sprintf("  Config Manager: %s\n", c.getComponentStatus(c.configManager)))
	report.WriteString(fmt.Sprintf("  Validator: %s\n", c.getComponentStatus(c.validator)))
	report.WriteString(fmt.Sprintf("  Template Manager: %s\n", c.getComponentStatus(c.templateManager)))
	report.WriteString(fmt.Sprintf("  Cache Manager: %s\n", c.getComponentStatus(c.cacheManager)))
	report.WriteString(fmt.Sprintf("  Version Manager: %s\n", c.getComponentStatus(c.versionManager)))
	report.WriteString(fmt.Sprintf("  Audit Engine: %s\n", c.getComponentStatus(c.auditEngine)))
	report.WriteString(fmt.Sprintf("  Logger: %s\n", c.getComponentStatus(c.logger)))

	return report.String()
}

// getComponentStatus returns the status of a component
func (c *CLI) getComponentStatus(component interface{}) string {
	if component == nil {
		return c.Error("Not Initialized")
	}
	return c.Success("Initialized")
}

// HandleEnhancedError handles an error with enhanced error handling
func (c *CLI) HandleEnhancedError(err error, operation string, context map[string]interface{}) int {
	if c.errorIntegration != nil {
		return c.errorIntegration.HandleError(err, operation, context)
	}

	// Fallback to basic error handling
	c.ErrorOutput("Error in %s: %v", operation, err)
	return 1
}

// CreateWorkflowError creates a workflow-specific error with context
func (c *CLI) CreateWorkflowError(message string, workflowType string) error {
	return fmt.Errorf("workflow error [%s]: %s", workflowType, message)
}

// registerHealthCheckers registers health checkers for system monitoring
func (c *CLI) registerHealthCheckers() {
	if c.systemMonitor == nil {
		return
	}

	// Register template health checker
	if c.templateManager != nil {
		templateChecker := &performance.TemplateHealthChecker{}
		c.systemMonitor.RegisterHealthChecker("templates", templateChecker)
	}

	// Register config health checker
	if c.configManager != nil {
		configChecker := &performance.ConfigHealthChecker{}
		c.systemMonitor.RegisterHealthChecker("config", configChecker)
	}
}

// registerLazyLoaders registers lazy loaders for expensive operations
func (c *CLI) registerLazyLoaders() {
	// Template list lazy loader
	templateLoader := performance.NewTemplateLazyLoader(c.templateManager, c.cacheManager)
	c.lazyLoaderManager.Register("templates", templateLoader)

	// Version info lazy loader
	versionLoader := performance.NewVersionLazyLoader(c.versionManager, c.cacheManager)
	c.lazyLoaderManager.Register("version", versionLoader)

	// Register with optimizer
	c.optimizer.RegisterLazyLoader("templates", templateLoader)
	c.optimizer.RegisterLazyLoader("version", versionLoader)
}

// OptimizeCommand wraps command execution with performance optimizations
func (c *CLI) OptimizeCommand(commandName string, operation func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx := context.Background()
	result, metrics, err := c.optimizer.OptimizeCommand(ctx, commandName, operation)

	// Log performance metrics if verbose mode is enabled
	if c.verboseMode && metrics != nil {
		c.VerboseOutput("⚡ Command '%s' completed in %v", commandName, metrics.Duration)
		if metrics.CacheHits > 0 {
			c.VerboseOutput("📊 Cache hits: %d, Cache misses: %d", metrics.CacheHits, metrics.CacheMisses)
		}
		if len(metrics.Optimizations) > 0 {
			c.VerboseOutput("🔧 Optimizations used: %v", metrics.Optimizations)
		}
	}

	return result, err
}

// GetPerformanceReport returns a performance report for all commands
func (c *CLI) GetPerformanceReport() *performance.MetricsReport {
	if c.optimizer == nil {
		return nil
	}
	return c.optimizer.GetPerformanceReport()
}

// ClearPerformanceCache clears performance-related cache entries
func (c *CLI) ClearPerformanceCache() error {
	if c.optimizer == nil {
		return fmt.Errorf("performance optimizer not initialized")
	}
	return c.optimizer.ClearCache()
}

// StartSystemMonitoring starts system health monitoring
func (c *CLI) StartSystemMonitoring() error {
	if c.systemMonitor == nil {
		return fmt.Errorf("system monitor not initialized")
	}

	ctx := context.Background()
	return c.systemMonitor.StartMonitoring(ctx)
}

// StopSystemMonitoring stops system health monitoring
func (c *CLI) StopSystemMonitoring() error {
	if c.systemMonitor == nil {
		return fmt.Errorf("system monitor not initialized")
	}
	return c.systemMonitor.StopMonitoring()
}

// GetSystemHealth returns current system health status
func (c *CLI) GetSystemHealth() *performance.SystemHealthSnapshot {
	if c.systemMonitor == nil {
		return nil
	}

	ctx := context.Background()
	return c.systemMonitor.GetCurrentHealth(ctx)
}

// GetDashboardReport returns a dashboard report
func (c *CLI) GetDashboardReport(format string) (string, error) {
	if c.dashboard == nil {
		return "", fmt.Errorf("dashboard not initialized")
	}

	switch format {
	case "json":
		return c.dashboard.GenerateJSONReport()
	case "text", "":
		return c.dashboard.GenerateTextReport(), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// RecordError records an error for diagnostics
func (c *CLI) RecordError(component, operation, errorType, message, severity string, context map[string]interface{}) {
	if c.systemMonitor != nil {
		c.systemMonitor.RecordError(component, operation, errorType, message, severity, context)
	}
}

// RecordPerformance records a performance measurement
func (c *CLI) RecordPerformance(operation string, duration time.Duration, memoryUsed int64, success bool, metadata map[string]interface{}) {
	if c.systemMonitor != nil {
		c.systemMonitor.RecordPerformance(operation, duration, memoryUsed, success, metadata)
	}
}

// Close closes the CLI and releases resources
func (c *CLI) Close() error {
	var errors []error

	// Close error integration
	if c.errorIntegration != nil {
		if err := c.errorIntegration.Close(); err != nil {
			errors = append(errors, fmt.Errorf("error integration close failed: %w", err))
		}
	}

	// Close other components as needed
	// (Add other component cleanup here)

	if len(errors) > 0 {
		return fmt.Errorf("errors during CLI close: %v", errors)
	}

	return nil
}
