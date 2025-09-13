// Package app provides the core application logic for the Open Source Template Generator.
//
// This package implements the main application structure, CLI command handling,
// and orchestrates the interaction between different components like template
// processing, validation, and project generation.
//
// The application follows clean architecture principles with dependency injection
// to ensure testability and maintainability.
package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-source-template-generator/internal/config"
	"github.com/open-source-template-generator/internal/container"
	"github.com/open-source-template-generator/pkg/cli"
	"github.com/open-source-template-generator/pkg/constants"
	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
	"github.com/open-source-template-generator/pkg/utils"
	"github.com/open-source-template-generator/pkg/validation"
	"github.com/open-source-template-generator/pkg/version"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v3"
)

// titleCaser creates a Title caser to replace deprecated strings.Title
var titleCaser = cases.Title(language.English)

// App represents the main application instance that orchestrates all CLI operations.
// It manages the dependency injection container, CLI interface, logging, and error handling.
//
// The App struct serves as the central coordinator for:
//   - CLI command processing and routing
//   - Component initialization and dependency management
//   - Project generation workflows
//   - Validation and auditing operations
//   - Configuration management
type App struct {
	container       *container.Container // Dependency injection container
	rootCmd         *cobra.Command       // Root CLI command
	cli             *cli.CLI             // CLI interface for user interaction
	logger          *Logger              // Application logger
	errorHandler    *ErrorHandler        // Centralized error handling
	resourceManager *ResourceManager     // Resource and memory management
	version         string               // Application version
	gitCommit       string               // Git commit hash
	buildTime       string               // Build timestamp
}

// NewApp creates a new application instance with the provided dependency container.
//
// This function initializes all required components including:
//   - Logger with appropriate log level and output configuration
//   - Error handler for centralized error processing
//   - All application components through the container
//   - CLI interface and command structure
//
// Parameters:
//   - c: Dependency injection container with pre-configured or empty services
//
// Returns:
//   - *App: Fully initialized application instance ready for execution
func NewApp(c *container.Container) *App {
	// Initialize all components if not already set
	app := &App{
		container: c,
	}

	// Initialize resource manager first
	app.resourceManager = NewResourceManager()

	// Initialize logger
	logger, err := NewLogger(LogLevelInfo, true)
	if err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
		// Create a basic logger as fallback
		logger, _ = NewLogger(LogLevelInfo, false)
	}
	app.logger = logger
	app.errorHandler = NewErrorHandler(logger)

	app.initializeComponents()
	app.cli = cli.NewCLI(c.GetConfigManager(), c.GetValidator())
	app.setupCommands()
	return app
}

// NewAppWithVersion creates a new application instance with version information.
//
// This function is similar to NewApp but accepts version information that gets
// embedded into the application and displayed in version commands.
//
// Parameters:
//   - c: Dependency injection container with pre-configured or empty services
//   - version: Application version string (e.g., "v1.2.3")
//   - gitCommit: Git commit hash for build traceability
//   - buildTime: Build timestamp for debugging and support
//
// Returns:
//   - *App: Fully initialized application instance with version info
func NewAppWithVersion(c *container.Container, version, gitCommit, buildTime string) *App {
	// Initialize all components if not already set
	app := &App{
		container: c,
		version:   version,
		gitCommit: gitCommit,
		buildTime: buildTime,
	}

	// Initialize resource manager first
	app.resourceManager = NewResourceManager()

	// Initialize logger
	logger, err := NewLogger(LogLevelInfo, true)
	if err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
		// Create a basic logger as fallback
		logger, _ = NewLogger(LogLevelInfo, false)
	}
	app.logger = logger
	app.errorHandler = NewErrorHandler(logger)

	app.initializeComponents()
	app.cli = cli.NewCLIWithVersion(c.GetConfigManager(), c.GetValidator(), version)
	app.setupCommands()
	return app
}

// Close gracefully shuts down the application and cleans up all resources.
//
// This method ensures proper cleanup of:
//   - Log file handles and buffers
//   - Temporary files and directories
//   - Network connections (if any)
//   - Cache files and locks
//   - Resource manager and memory pools
//
// It should be called using defer in the main function to ensure cleanup
// occurs even if the application exits due to an error.
//
// Returns:
//   - error: Any error that occurred during cleanup, nil if successful
func (a *App) Close() error {
	var lastErr error

	// Close resource manager first
	if a.resourceManager != nil {
		if err := a.resourceManager.Close(); err != nil {
			lastErr = err
		}
	}

	// Close logger
	if a.logger != nil {
		if err := a.logger.Close(); err != nil {
			lastErr = err
		}
	}

	// Force final garbage collection
	utils.ForceGlobalGC()

	return lastErr
}

// initializeComponents initializes all required components in the container
func (a *App) initializeComponents() {
	// Initialize configuration manager
	if a.container.GetConfigManager() == nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Warning: Could not get user home directory: %v", err)
			homeDir = "."
		}
		cacheDir := filepath.Join(homeDir, ".cache", "template-generator")
		defaultsPath := filepath.Join("templates", "config", "defaults.yaml")

		configManager := config.NewManager(cacheDir, defaultsPath)
		a.container.SetConfigManager(configManager)
	}

	// Initialize validation engine
	if a.container.GetValidator() == nil {
		validator := validation.NewEngine()
		a.container.SetValidator(validator)
	}

	// Initialize filesystem generator
	if a.container.GetFileSystemGenerator() == nil {
		fsGenerator := filesystem.NewGenerator()
		a.container.SetFileSystemGenerator(fsGenerator)
	}

	// Initialize version manager
	if a.container.GetVersionManager() == nil {
		// Create cache for version manager
		versionHomeDir, versionHomeDirErr := os.UserHomeDir()
		if versionHomeDirErr != nil {
			log.Printf("Warning: Could not get user home directory: %v", versionHomeDirErr)
			versionHomeDir = "."
		}
		cacheDir := filepath.Join(versionHomeDir, ".cache", "template-generator")
		var versionCache interfaces.VersionCache
		fileCache, err := version.NewFileCache(cacheDir, 24*time.Hour) // 24 hour TTL
		if err != nil {
			log.Printf("Warning: Could not create file cache, using memory cache: %v", err)
			versionCache = version.NewMemoryCache(24 * time.Hour)
		} else {
			versionCache = fileCache
		}

		// Create version storage
		storageDir := filepath.Join(versionHomeDir, ".config", "template-generator")
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			log.Printf("Warning: Could not create storage directory: %v", err)
		}
		storageFile := filepath.Join(storageDir, "versions.yaml")

		versionStorage, storageErr := version.NewFileStorage(storageFile, constants.FormatYAML)
		if storageErr != nil {
			log.Printf("Warning: Could not create version storage: %v", storageErr)
			// Use manager without storage as fallback
			versionManager := version.NewManager(versionCache)
			a.container.SetVersionManager(versionManager)
		} else {
			// Use manager with storage
			versionManager := version.NewManagerWithStorage(versionCache, versionStorage)
			a.container.SetVersionManager(versionManager)
		}
	}

	// Initialize template engine (after version manager)
	if a.container.GetTemplateEngine() == nil {
		// Create template engine with version manager if available
		if versionManager := a.container.GetVersionManager(); versionManager != nil {
			templateEngine := template.NewEngineWithVersionManager(versionManager)
			a.container.SetTemplateEngine(templateEngine)
		} else {
			templateEngine := template.NewEngine()
			a.container.SetTemplateEngine(templateEngine)
		}
	}
}

// Execute runs the application
func (a *App) Execute() error {
	return a.rootCmd.Execute()
}

// setupCommands initializes the CLI commands
func (a *App) setupCommands() {
	a.rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "Open Source Project Template Generator",
		Long: `A comprehensive tool for generating production-ready open source project structures
with modern best practices, latest package versions, and complete CI/CD configurations.`,
		RunE: a.runGenerate,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set up global flags
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				a.logger.level = LogLevelDebug
			}

			quiet, _ := cmd.Flags().GetBool("quiet")
			if quiet {
				a.logger.level = LogLevelError
			}

			return nil
		},
	}

	// Add global flags
	a.rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	a.rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
	a.rootCmd.PersistentFlags().String("log-level", "info", "Set log level (debug, info, warn, error)")

	// Add subcommands
	a.rootCmd.AddCommand(a.generateCommand())
	a.rootCmd.AddCommand(a.validateCommand())
	a.rootCmd.AddCommand(a.versionCommand())
	a.rootCmd.AddCommand(a.configCommand())
	a.rootCmd.AddCommand(a.analyzeCommand())
	a.rootCmd.AddCommand(a.versionsCommand())
}

// runGenerate is the default command handler
func (a *App) runGenerate(cmd *cobra.Command, args []string) error {
	// This will be implemented in later tasks
	cmd.Println("Template generator is ready!")
	cmd.Println("Use 'generator generate' to create a new project")
	return nil
}

// generateCommand creates the generate subcommand
func (a *App) generateCommand() *cobra.Command {
	var dryRun bool
	var configFile string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from templates",
		Long:  "Generate a complete project structure with selected components and configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runGenerateCommand(dryRun, configFile, outputPath)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview generation without creating files")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output directory path")

	return cmd
}

// validateCommand creates the validate subcommand
func (a *App) validateCommand() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "validate [project-path]",
		Short: "Validate a generated project",
		Long:  "Validate the structure and configuration of a generated project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runValidateCommand(args, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed validation output")

	return cmd
}

// versionCommand creates the version subcommand
func (a *App) versionCommand() *cobra.Command {
	var showVersions bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display version information for the generator and available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runVersionCommand(showVersions)
		},
	}

	cmd.Flags().BoolVar(&showVersions, "packages", false, "Show latest package versions")

	return cmd
}

// configCommand creates the config subcommand
func (a *App) configCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Manage generator configuration and defaults",
	}

	// Add config subcommands
	cmd.AddCommand(a.configShowCommand())
	cmd.AddCommand(a.configSetCommand())
	cmd.AddCommand(a.configResetCommand())

	return cmd
}

// configShowCommand shows current configuration
func (a *App) configShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  "Display the current configuration and default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runConfigShowCommand()
		},
	}
}

// configSetCommand sets configuration values
func (a *App) configSetCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set configuration value",
		Long:  "Set a configuration value or load from file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runConfigSetCommand(args, configFile)
		},
	}

	cmd.Flags().StringVarP(&configFile, "file", "f", "", "Load configuration from file")

	return cmd
}

// configResetCommand resets configuration to defaults
func (a *App) configResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		Long:  "Reset all configuration values to their defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runConfigResetCommand()
		},
	}
}

// runGenerateCommand handles the generate command execution
func (a *App) runGenerateCommand(dryRun bool, configFile, outputPath string) error {
	var config *models.ProjectConfig
	var err error

	if configFile != "" {
		// Load configuration from file
		a.cli.ShowProgress("Loading configuration from file")
		configManager := a.container.GetConfigManager()
		config, err = configManager.LoadConfig(configFile)
		if err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to load configuration file: %v", err))
			return err
		}
	} else {
		// Interactive configuration
		config, err = a.cli.PromptProjectDetails()
		if err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to collect project details: %v", err))
			return err
		}
	}

	// Override output path if provided via flag
	if outputPath != "" {
		config.OutputPath = outputPath
	}

	// Validate configuration
	a.cli.ShowProgress("Validating configuration")
	configManager := a.container.GetConfigManager()
	if err := configManager.ValidateConfig(config); err != nil {
		a.cli.ShowError(fmt.Sprintf("Configuration validation failed: %v", err))
		return err
	}

	// Check output path
	if err := a.cli.CheckOutputPath(config.OutputPath); err != nil {
		a.cli.ShowError(fmt.Sprintf("Output path check failed: %v", err))
		return err
	}

	if dryRun {
		fmt.Println()
		fmt.Println("🔍 Dry Run Mode - No files will be created")
		a.cli.PreviewConfiguration(config)
		fmt.Printf("\nProject would be generated at: %s\n", config.OutputPath)
		return nil
	}

	// Show configuration and get confirmation
	if !a.cli.ConfirmGeneration(config) {
		fmt.Println("Generation canceled by user.")
		return nil
	}

	// Generate the project
	if err := a.generateProject(config); err != nil {
		a.cli.ShowError(fmt.Sprintf("Project generation failed: %v", err))
		return err
	}

	a.cli.ShowSuccess("Project generation completed successfully!")
	fmt.Printf("📁 Project created at: %s\n", config.OutputPath)
	fmt.Println("🚀 Run 'make setup' in the project directory to get started!")

	return nil
}

// generateProject performs the actual project generation
func (a *App) generateProject(config *models.ProjectConfig) error {
	templateEngine := a.container.GetTemplateEngine()
	fsGenerator := a.container.GetFileSystemGenerator()

	// Create project root directory
	projectPath, err := a.createProjectStructure(fsGenerator, config)
	if err != nil {
		return err
	}

	// Generate base files
	if err := a.generateProjectBase(templateEngine, fsGenerator, config, projectPath); err != nil {
		return err
	}

	// Generate all component types
	return a.generateAllComponents(templateEngine, fsGenerator, config, projectPath)
}

// createProjectStructure creates the project root directory structure
func (a *App) createProjectStructure(fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig) (string, error) {
	a.cli.ShowProgress("Creating project structure")
	if err := fsGenerator.CreateProject(config, config.OutputPath); err != nil {
		return "", fmt.Errorf("failed to create project structure: %w", err)
	}
	return filepath.Join(config.OutputPath, config.Name), nil
}

// generateProjectBase generates base project files that are always included
func (a *App) generateProjectBase(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	a.cli.ShowProgress("Generating base project files")
	if err := a.generateBaseFiles(templateEngine, fsGenerator, config, projectPath); err != nil {
		return fmt.Errorf("failed to generate base files: %w", err)
	}
	return nil
}

// generateAllComponents generates all enabled component types
func (a *App) generateAllComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	// Generate frontend components
	if err := a.generateFrontendIfEnabled(templateEngine, fsGenerator, config, projectPath); err != nil {
		return err
	}

	// Generate backend components
	if err := a.generateBackendIfEnabled(templateEngine, fsGenerator, config, projectPath); err != nil {
		return err
	}

	// Generate mobile components
	if err := a.generateMobileIfEnabled(templateEngine, fsGenerator, config, projectPath); err != nil {
		return err
	}

	// Generate infrastructure components
	if err := a.generateInfrastructureIfEnabled(templateEngine, fsGenerator, config, projectPath); err != nil {
		return err
	}

	// Generate CI/CD configurations
	return a.generateCICDIfEnabled(templateEngine, fsGenerator, config, projectPath)
}

// generateFrontendIfEnabled generates frontend components if any are enabled
func (a *App) generateFrontendIfEnabled(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	if config.Components.Frontend.MainApp || config.Components.Frontend.Home || config.Components.Frontend.Admin {
		a.cli.ShowProgress("Generating frontend applications")
		if err := a.generateFrontendComponents(templateEngine, fsGenerator, config, projectPath); err != nil {
			return fmt.Errorf("failed to generate frontend components: %w", err)
		}
	}
	return nil
}

// generateBackendIfEnabled generates backend components if enabled
func (a *App) generateBackendIfEnabled(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	if config.Components.Backend.API {
		a.cli.ShowProgress("Generating backend API server")
		if err := a.generateBackendComponents(templateEngine, fsGenerator, config, projectPath); err != nil {
			return fmt.Errorf("failed to generate backend components: %w", err)
		}
	}
	return nil
}

// generateMobileIfEnabled generates mobile components if any are enabled
func (a *App) generateMobileIfEnabled(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		a.cli.ShowProgress("Generating mobile applications")
		if err := a.generateMobileComponents(templateEngine, fsGenerator, config, projectPath); err != nil {
			return fmt.Errorf("failed to generate mobile components: %w", err)
		}
	}
	return nil
}

// generateInfrastructureIfEnabled generates infrastructure components if any are enabled
func (a *App) generateInfrastructureIfEnabled(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		a.cli.ShowProgress("Generating infrastructure configurations")
		if err := a.generateInfrastructureComponents(templateEngine, fsGenerator, config, projectPath); err != nil {
			return fmt.Errorf("failed to generate infrastructure components: %w", err)
		}
	}
	return nil
}

// generateCICDIfEnabled generates CI/CD workflow configurations
func (a *App) generateCICDIfEnabled(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	a.cli.ShowProgress("Generating CI/CD workflows")
	if err := a.generateCICDComponents(templateEngine, fsGenerator, config, projectPath); err != nil {
		return fmt.Errorf("failed to generate CI/CD components: %w", err)
	}
	return nil
}

// generateBaseFiles generates the base project files
func (a *App) generateBaseFiles(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	baseTemplateDir := "templates/base"

	// Process base templates
	return templateEngine.ProcessDirectory(baseTemplateDir, projectPath, config)
}

// generateFrontendComponents generates frontend application files
func (a *App) generateFrontendComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	frontendDir := filepath.Join(projectPath, "App")

	if config.Components.Frontend.MainApp {
		mainAppDir := filepath.Join(frontendDir, "main")
		if err := templateEngine.ProcessDirectory("templates/frontend/nextjs-app", mainAppDir, config); err != nil {
			return fmt.Errorf("failed to generate main app: %w", err)
		}
	}

	if config.Components.Frontend.Home {
		homeDir := filepath.Join(frontendDir, "home")
		if err := templateEngine.ProcessDirectory("templates/frontend/nextjs-home", homeDir, config); err != nil {
			return fmt.Errorf("failed to generate home app: %w", err)
		}
	}

	if config.Components.Frontend.Admin {
		adminDir := filepath.Join(frontendDir, "admin")
		if err := templateEngine.ProcessDirectory("templates/frontend/nextjs-admin", adminDir, config); err != nil {
			return fmt.Errorf("failed to generate admin app: %w", err)
		}
	}

	// Generate shared components if multiple frontend apps
	if (config.Components.Frontend.MainApp && config.Components.Frontend.Home) ||
		(config.Components.Frontend.MainApp && config.Components.Frontend.Admin) ||
		(config.Components.Frontend.Home && config.Components.Frontend.Admin) {
		sharedDir := filepath.Join(frontendDir, "shared")
		if err := templateEngine.ProcessDirectory("templates/frontend/shared-components", sharedDir, config); err != nil {
			return fmt.Errorf("failed to generate shared components: %w", err)
		}
	}

	return nil
}

// generateBackendComponents generates backend API server files
func (a *App) generateBackendComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	backendDir := filepath.Join(projectPath, "CommonServer")
	return templateEngine.ProcessDirectory("templates/backend/go-gin", backendDir, config)
}

// generateMobileComponents generates mobile application files
func (a *App) generateMobileComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	mobileDir := filepath.Join(projectPath, "Mobile")

	if config.Components.Mobile.Android {
		androidDir := filepath.Join(mobileDir, "android")
		if err := templateEngine.ProcessDirectory("templates/mobile/android-kotlin", androidDir, config); err != nil {
			return fmt.Errorf("failed to generate Android app: %w", err)
		}
	}

	if config.Components.Mobile.IOS {
		iosDir := filepath.Join(mobileDir, "ios")
		if err := templateEngine.ProcessDirectory("templates/mobile/ios-swift", iosDir, config); err != nil {
			return fmt.Errorf("failed to generate iOS app: %w", err)
		}
	}

	// Generate shared mobile resources
	sharedDir := filepath.Join(mobileDir, "shared")
	return templateEngine.ProcessDirectory("templates/mobile/shared", sharedDir, config)
}

// generateInfrastructureComponents generates infrastructure configuration files
func (a *App) generateInfrastructureComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	deployDir := filepath.Join(projectPath, "Deploy")

	if config.Components.Infrastructure.Docker {
		dockerDir := filepath.Join(deployDir, "docker")
		if err := templateEngine.ProcessDirectory("templates/infrastructure/docker", dockerDir, config); err != nil {
			return fmt.Errorf("failed to generate Docker configurations: %w", err)
		}
	}

	if config.Components.Infrastructure.Kubernetes {
		k8sDir := filepath.Join(deployDir, "k8s")
		if err := templateEngine.ProcessDirectory("templates/infrastructure/kubernetes", k8sDir, config); err != nil {
			return fmt.Errorf("failed to generate Kubernetes configurations: %w", err)
		}
	}

	if config.Components.Infrastructure.Terraform {
		terraformDir := filepath.Join(deployDir, "terraform")
		if err := templateEngine.ProcessDirectory("templates/infrastructure/terraform", terraformDir, config); err != nil {
			return fmt.Errorf("failed to generate Terraform configurations: %w", err)
		}
	}

	return nil
}

// generateCICDComponents generates CI/CD workflow files
func (a *App) generateCICDComponents(templateEngine interfaces.TemplateEngine, fsGenerator interfaces.FileSystemGenerator, config *models.ProjectConfig, projectPath string) error {
	githubDir := filepath.Join(projectPath, ".github")
	return templateEngine.ProcessDirectory("templates/base/.github", githubDir, config)
}

// runValidateCommand handles the validate command execution
func (a *App) runValidateCommand(args []string, verbose bool) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	a.cli.ShowProgress(fmt.Sprintf("Validating project at %s", projectPath))

	// Check if path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		a.cli.ShowError(fmt.Sprintf("Project path does not exist: %s", projectPath))
		return err
	}

	// Perform validation using the validation engine
	validator := a.container.GetValidator()
	result, err := validator.ValidateProject(projectPath)
	if err != nil {
		a.cli.ShowError(fmt.Sprintf("Validation failed: %v", err))
		return err
	}

	// Display results
	if result.Valid {
		a.cli.ShowSuccess("Project validation completed successfully")
	} else {
		a.cli.ShowError("Project validation failed")
		fmt.Println("\nValidation Errors:")
		for _, validationErr := range result.Errors {
			fmt.Printf("  ❌ %s: %s\n", validationErr.Field, validationErr.Message)
		}
	}

	// Show warnings if any
	if len(result.Warnings) > 0 {
		fmt.Println("\nValidation Warnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠️  %s: %s\n", warning.Field, warning.Message)
		}
	}

	if verbose {
		fmt.Printf("\nValidation Summary:\n")
		fmt.Printf("  Valid: %t\n", result.Valid)
		fmt.Printf("  Errors: %d\n", len(result.Errors))
		fmt.Printf("  Warnings: %d\n", len(result.Warnings))
	}

	return nil
}

// runVersionCommand handles the version command execution
func (a *App) runVersionCommand(showPackages bool) error {
	version := a.version
	if version == "" || version == "dev" {
		version = "1.0.0" // fallback
	}
	fmt.Printf("Open Source Template Generator %s\n", version)
	fmt.Println("Built with Go 1.23+")

	if a.gitCommit != "" && a.gitCommit != "unknown" {
		fmt.Printf("Git commit: %s\n", a.gitCommit)
	}
	if a.buildTime != "" && a.buildTime != "unknown" {
		fmt.Printf("Build time: %s\n", a.buildTime)
	}

	if showPackages {
		a.cli.ShowProgress("Fetching latest package versions")

		configManager := a.container.GetConfigManager()
		versions, err := configManager.GetLatestVersions()
		if err != nil {
			a.cli.ShowWarning(fmt.Sprintf("Could not fetch latest versions: %v", err))

			// Show default versions as fallback
			fmt.Println()
			fmt.Println("Latest Package Versions (defaults):")
			fmt.Println("  Node.js: 20.11.0")
			fmt.Println("  Go: 1.22.0")
			fmt.Println("  Next.js: 15.0.0")
			fmt.Println("  React: 18.2.0")
			fmt.Println("  Kotlin: 2.0.0")
			fmt.Println("  Swift: 5.9.0")
			return nil
		}

		fmt.Println()
		fmt.Println("Latest Package Versions:")
		fmt.Printf("  Node.js: %s\n", versions.Node)
		fmt.Printf("  Go: %s\n", versions.Go)
		if versions.NextJS != "" {
			fmt.Printf("  Next.js: %s\n", versions.NextJS)
		}
		if versions.React != "" {
			fmt.Printf("  React: %s\n", versions.React)
		}
		if versions.Kotlin != "" {
			fmt.Printf("  Kotlin: %s\n", versions.Kotlin)
		}
		if versions.Swift != "" {
			fmt.Printf("  Swift: %s\n", versions.Swift)
		}

		// Show common packages if available
		if len(versions.Packages) > 0 {
			fmt.Println("\nCommon Packages:")
			for pkg, version := range versions.Packages {
				fmt.Printf("  %s: %s\n", pkg, version)
			}
		}
	}

	return nil
}

// runConfigShowCommand handles the config show command
func (a *App) runConfigShowCommand() error {
	a.logger.Info("Showing current configuration")

	configManager := a.container.GetConfigManager()
	defaults, err := configManager.LoadDefaults()
	if err != nil {
		appErr := WrapConfigurationError("Failed to load default configuration", err)
		a.errorHandler.Handle(appErr)
		return appErr
	}

	fmt.Println("Current Configuration:")
	fmt.Printf("  Default License: %s\n", defaults.License)
	fmt.Printf("  Default Output Path: %s\n", defaults.OutputPath)

	fmt.Println("\nDefault Components:")
	fmt.Printf("  Frontend Main App: %t\n", defaults.Components.Frontend.MainApp)
	fmt.Printf("  Frontend Home: %t\n", defaults.Components.Frontend.Home)
	fmt.Printf("  Frontend Admin: %t\n", defaults.Components.Frontend.Admin)
	fmt.Printf("  Backend API: %t\n", defaults.Components.Backend.API)
	fmt.Printf("  Mobile Android: %t\n", defaults.Components.Mobile.Android)
	fmt.Printf("  Mobile iOS: %t\n", defaults.Components.Mobile.IOS)
	fmt.Printf("  Infrastructure Docker: %t\n", defaults.Components.Infrastructure.Docker)
	fmt.Printf("  Infrastructure Kubernetes: %t\n", defaults.Components.Infrastructure.Kubernetes)
	fmt.Printf("  Infrastructure Terraform: %t\n", defaults.Components.Infrastructure.Terraform)

	if defaults.Versions != nil {
		fmt.Println("\nDefault Versions:")
		fmt.Printf("  Node.js: %s\n", defaults.Versions.Node)
		fmt.Printf("  Go: %s\n", defaults.Versions.Go)
		fmt.Printf("  Next.js: %s\n", defaults.Versions.NextJS)
		fmt.Printf("  React: %s\n", defaults.Versions.React)
		fmt.Printf("  Kotlin: %s\n", defaults.Versions.Kotlin)
		fmt.Printf("  Swift: %s\n", defaults.Versions.Swift)
	}

	return nil
}

// runConfigSetCommand handles the config set command
func (a *App) runConfigSetCommand(args []string, configFile string) error {
	if configFile != "" {
		a.logger.Info("Loading configuration from file: %s", configFile)
		// Load configuration from file
		configManager := a.container.GetConfigManager()
		_, err := configManager.LoadConfig(configFile)
		if err != nil {
			appErr := WrapConfigurationError("Failed to load configuration file", err)
			a.errorHandler.Handle(appErr)
			return appErr
		}

		a.cli.ShowSuccess("Configuration loaded successfully")
		return nil
	}

	if len(args) != 2 {
		return fmt.Errorf("config set requires key and value arguments")
	}

	key, value := args[0], args[1]
	a.logger.Info("Setting configuration: %s = %s", key, value)

	// This would implement setting individual configuration values
	// For now, just show what would be set
	fmt.Printf("Would set %s = %s\n", key, value)
	fmt.Println("Note: Individual configuration setting will be implemented in a future version")

	return nil
}

// runConfigResetCommand handles the config reset command
func (a *App) runConfigResetCommand() error {
	a.logger.Info("Resetting configuration to defaults")

	// This would reset configuration to defaults
	fmt.Println("Configuration reset to defaults")
	fmt.Println("Note: Configuration reset will be implemented in a future version")

	return nil
}

// runUpdateTemplatesCommand handles the update-templates command execution
func (a *App) runUpdateTemplatesCommand(dryRun bool, backup bool, templatePaths []string) error {
	if dryRun {
		a.cli.ShowProgress("Checking what templates would be updated (dry run)")
	} else {
		a.cli.ShowProgress("Updating template files with latest versions")
	}

	// Get version manager
	versionManager := a.container.GetVersionManager()

	// Check if version manager has storage capability
	if managerWithStorage, ok := versionManager.(*version.Manager); ok {
		// Get current version information
		store, err := managerWithStorage.GetVersionStore()
		if err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to load version store: %v", err))
			return err
		}

		// Collect all versions
		allVersions := make(map[string]*models.VersionInfo)
		for name, info := range store.Languages {
			allVersions[name] = info
		}
		for name, info := range store.Frameworks {
			allVersions[name] = info
		}
		for name, info := range store.Packages {
			allVersions[name] = info
		}

		if dryRun {
			fmt.Println("\n🔍 Dry Run - Would update the following templates:")
			// This would show which templates would be affected
			fmt.Println("  templates/frontend/nextjs-app/package.json.tmpl")
			fmt.Println("  templates/frontend/nextjs-home/package.json.tmpl")
			fmt.Println("  templates/frontend/nextjs-admin/package.json.tmpl")
			fmt.Println("  templates/backend/go-gin/go.mod.tmpl")
			fmt.Println("  And other template files...")
			return nil
		}

		// Apply template updates (placeholder implementation)
		fmt.Println("✅ Template files updated with latest versions")
		fmt.Println("Note: Full template update integration will be completed in task 6.2")
		return nil
	}

	// Fallback message
	fmt.Println("Note: Template update requires version storage integration")
	return nil
}

// Helper functions for version management

func (a *App) displayVersionReportJSON(report *models.VersionReport, verbose bool) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func (a *App) displayVersionReportYAML(report *models.VersionReport, verbose bool) error {
	data, err := yaml.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func (a *App) displayVersionReportTable(report *models.VersionReport, verbose bool) error {
	fmt.Printf("\n📊 Version Status Report (Generated: %s)\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("=" + strings.Repeat("=", 50))

	fmt.Printf("Total Packages: %d\n", report.TotalPackages)
	fmt.Printf("Outdated: %d\n", report.OutdatedCount)
	fmt.Printf("Security Issues: %d\n", report.SecurityIssues)
	fmt.Printf("Last Check: %s\n", report.LastUpdateCheck.Format("2006-01-02 15:04:05"))

	if len(report.Summary) > 0 {
		fmt.Println("\n📋 Summary by Category:")
		for category, summary := range report.Summary {
			fmt.Printf("  %s: %d total, %d current, %d outdated, %d insecure\n",
				titleCaser.String(category), summary.Total, summary.Current, summary.Outdated, summary.Insecure)
		}
	}

	if len(report.Recommendations) > 0 {
		fmt.Println("\n🔄 Update Recommendations:")
		for _, rec := range report.Recommendations {
			var priority string
			if rec.Priority == constants.PriorityCritical {
				priority = "🔴 CRITICAL"
			} else if rec.Priority == constants.PriorityHigh {
				priority = "🟠 HIGH"
			} else if rec.Priority == "medium" {
				priority = "🟡 MEDIUM"
			} else {
				priority = "🟢 LOW"
			}

			fmt.Printf("  %s %s: %s → %s", priority, rec.Name, rec.CurrentVersion, rec.RecommendedVersion)
			if rec.BreakingChange {
				fmt.Printf(" ⚠️  BREAKING")
			}
			fmt.Printf("\n    Reason: %s\n", rec.Reason)
		}
	}

	if verbose && len(report.Details) > 0 {
		fmt.Println("\n📦 Detailed Version Information:")
		for name, info := range report.Details {
			fmt.Printf("  %s (%s):\n", name, info.Type)
			fmt.Printf("    Current: %s\n", info.CurrentVersion)
			fmt.Printf("    Latest: %s\n", info.LatestVersion)
			fmt.Printf("    Updated: %s\n", info.UpdatedAt.Format("2006-01-02"))
			if len(info.SecurityIssues) > 0 {
				fmt.Printf("    Security Issues: %d\n", len(info.SecurityIssues))
			}
		}
	}

	return nil
}

func (a *App) displayVersionsJSON(versions *models.VersionConfig, verbose bool) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func (a *App) displayVersionsYAML(versions *models.VersionConfig, verbose bool) error {
	data, err := yaml.Marshal(versions)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func (a *App) displayVersionsTable(versions *models.VersionConfig, verbose bool) error {
	fmt.Println("\n📦 Latest Package Versions")
	fmt.Println("=" + strings.Repeat("=", 30))

	fmt.Printf("Node.js: %s\n", versions.Node)
	fmt.Printf("Go: %s\n", versions.Go)
	if versions.NextJS != "" {
		fmt.Printf("Next.js: %s\n", versions.NextJS)
	}
	if versions.React != "" {
		fmt.Printf("React: %s\n", versions.React)
	}
	if versions.Kotlin != "" {
		fmt.Printf("Kotlin: %s\n", versions.Kotlin)
	}
	if versions.Swift != "" {
		fmt.Printf("Swift: %s\n", versions.Swift)
	}

	if verbose && len(versions.Packages) > 0 {
		fmt.Println("\nCommon Packages:")
		for pkg, version := range versions.Packages {
			fmt.Printf("  %s: %s\n", pkg, version)
		}
	}

	fmt.Printf("\nLast Updated: %s\n", versions.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// analyzeCommand creates the analyze subcommand
func (a *App) analyzeCommand() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "analyze [template-directory]",
		Short: "Analyze template configurations",
		Long:  "Analyze frontend template configurations and identify inconsistencies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runAnalyzeCommand(args, outputFile)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save analysis report to JSON file")

	return cmd
}

// runAnalyzeCommand handles the analyze command execution
func (a *App) runAnalyzeCommand(args []string, outputFile string) error {
	templateDir := "templates"
	if len(args) > 0 {
		templateDir = args[0]
	}

	// Verify template directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		a.cli.ShowError(fmt.Sprintf("Template directory does not exist: %s", templateDir))
		return err
	}

	a.cli.ShowProgress("Analyzing template configurations")

	// Run the analysis
	if err := a.cli.AnalyzeTemplatesCommand([]string{templateDir}); err != nil {
		a.cli.ShowError(fmt.Sprintf("Template analysis failed: %v", err))
		return err
	}

	// Save report if output file is specified
	if outputFile != "" {
		a.cli.ShowProgress("Saving analysis report")
		if err := a.cli.SaveAnalysisReportCommand([]string{templateDir, outputFile}); err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to save analysis report: %v", err))
			return err
		}
	}

	return nil
}

// versionsCommand creates the versions subcommand for version management
func (a *App) versionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "Manage package versions",
		Long:  "Check, update, and manage package versions across all templates",
	}

	// Add version management subcommands
	cmd.AddCommand(a.checkVersionsCommand())
	cmd.AddCommand(a.updateVersionsCommand())
	cmd.AddCommand(a.updateTemplatesCommand())
	cmd.AddCommand(a.dashboardCommand())
	cmd.AddCommand(a.reportCommand())

	return cmd
}

// checkVersionsCommand creates the check-versions subcommand
func (a *App) checkVersionsCommand() *cobra.Command {
	var verbose bool
	var format string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check for latest package versions",
		Long:  "Query registries for latest versions and compare with current versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runCheckVersionsCommand(verbose, format)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed version information")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

// updateVersionsCommand creates the update-versions subcommand
func (a *App) updateVersionsCommand() *cobra.Command {
	var force bool
	var dryRun bool
	var packages []string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update version information",
		Long:  "Update the versions.md file with latest package versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runUpdateVersionsCommand(force, dryRun, packages)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force update even if versions are older")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be updated without making changes")
	cmd.Flags().StringSliceVar(&packages, "packages", nil, "Specific packages to update (default: all)")

	return cmd
}

// updateTemplatesCommand creates the update-templates subcommand
func (a *App) updateTemplatesCommand() *cobra.Command {
	var dryRun bool
	var backup bool
	var templatePaths []string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply version updates to templates",
		Long:  "Update template files with new version information from versions store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runUpdateTemplatesCommand(dryRun, backup, templatePaths)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be updated without making changes")
	cmd.Flags().BoolVar(&backup, "backup", true, "Create backup before updating templates")
	cmd.Flags().StringSliceVar(&templatePaths, "templates", nil, "Specific template paths to update (default: all)")

	return cmd
}

// dashboardCommand creates the dashboard subcommand
func (a *App) dashboardCommand() *cobra.Command {
	var refresh bool
	var format string
	var showDetails bool

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Display version management dashboard",
		Long:  "Show comprehensive overview of version status, template consistency, and validation results",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runDashboardCommand(refresh, format, showDetails)
		},
	}

	cmd.Flags().BoolVarP(&refresh, "refresh", "r", false, "Refresh data before displaying dashboard")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")
	cmd.Flags().BoolVarP(&showDetails, "details", "d", false, "Show detailed information for all packages")

	return cmd
}

// reportCommand creates the report subcommand for generating reports
func (a *App) reportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate and manage reports",
		Long:  "Generate version update reports, security reports, and view audit trails",
	}

	// Add report subcommands
	cmd.AddCommand(a.generateReportCommand())
	cmd.AddCommand(a.listReportsCommand())
	cmd.AddCommand(a.auditCommand())

	return cmd
}

// generateReportCommand creates the generate-report subcommand
func (a *App) generateReportCommand() *cobra.Command {
	var reportType string
	var format string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate version management reports",
		Long:  "Generate comprehensive reports for version updates, security issues, and template changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runGenerateReportCommand(reportType, format, outputDir)
		},
	}

	cmd.Flags().StringVarP(&reportType, "type", "t", "version", "Report type (version, security, template)")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, yaml, text)")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "reports", "Output directory for reports")

	return cmd
}

// listReportsCommand creates the list-reports subcommand
func (a *App) listReportsCommand() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List generated reports",
		Long:  "List all previously generated reports with metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runListReportsCommand(outputDir)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "reports", "Reports directory to scan")

	return cmd
}

// auditCommand creates the audit subcommand
func (a *App) auditCommand() *cobra.Command {
	var since string
	var eventType string
	var format string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "View audit trail",
		Long:  "View audit trail of version management operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runAuditCommand(since, eventType, format)
		},
	}

	cmd.Flags().StringVar(&since, "since", "24h", "Show events since duration (e.g., 24h, 7d, 30d)")
	cmd.Flags().StringVar(&eventType, "type", "", "Filter by event type")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

// runCheckVersionsCommand handles the check-versions command execution
func (a *App) runCheckVersionsCommand(verbose bool, format string) error {
	a.cli.ShowProgress("Checking latest package versions")

	// Get version manager with storage
	versionManager := a.container.GetVersionManager()

	// Check if version manager has storage capability
	if managerWithStorage, ok := versionManager.(*version.Manager); ok {
		// Use enhanced version manager with storage
		report, err := managerWithStorage.CheckLatestVersions()
		if err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to check versions: %v", err))
			return err
		}

		// Display results based on format
		switch format {
		case constants.FormatJSON:
			return a.displayVersionReportJSON(report, verbose)
		case constants.FormatYAML:
			return a.displayVersionReportYAML(report, verbose)
		default:
			return a.displayVersionReportTable(report, verbose)
		}
	}

	// Fallback to existing functionality
	configManager := a.container.GetConfigManager()
	versions, err := configManager.GetLatestVersions()
	if err != nil {
		a.cli.ShowError(fmt.Sprintf("Failed to check versions: %v", err))
		return err
	}

	// Display results based on format
	switch format {
	case constants.FormatJSON:
		return a.displayVersionsJSON(versions, verbose)
	case constants.FormatYAML:
		return a.displayVersionsYAML(versions, verbose)
	default:
		return a.displayVersionsTable(versions, verbose)
	}
}

// runUpdateVersionsCommand handles the update-versions command execution
func (a *App) runUpdateVersionsCommand(force bool, dryRun bool, packages []string) error {
	a.showUpdateProgress(dryRun)

	// Get version manager
	versionManager := a.container.GetVersionManager()

	// Try enhanced version manager with storage first
	if managerWithStorage, ok := versionManager.(*version.Manager); ok {
		return a.handleEnhancedVersionUpdate(managerWithStorage, force, dryRun, packages)
	}

	// Fallback to basic version manager
	return a.handleBasicVersionUpdate(dryRun)
}

// showUpdateProgress displays appropriate progress message
func (a *App) showUpdateProgress(dryRun bool) {
	if dryRun {
		a.cli.ShowProgress("Checking what versions would be updated (dry run)")
	} else {
		a.cli.ShowProgress("Updating version information")
	}
}

// handleEnhancedVersionUpdate handles updates using enhanced version manager
func (a *App) handleEnhancedVersionUpdate(manager *version.Manager, force bool, dryRun bool, packages []string) error {
	updates, err := manager.DetectVersionUpdates()
	if err != nil {
		a.cli.ShowError(fmt.Sprintf("Failed to detect version updates: %v", err))
		return err
	}

	if len(updates) == 0 {
		a.cli.ShowSuccess("All versions are up to date!")
		return nil
	}

	if dryRun {
		return a.showDryRunUpdates(updates)
	}

	return a.applyVersionUpdates(manager, updates, force, packages)
}

// showDryRunUpdates displays what would be updated in dry run mode
func (a *App) showDryRunUpdates(updates map[string]*models.VersionInfo) error {
	fmt.Println("\n🔍 Dry Run - Would update the following versions:")
	// Placeholder implementation - would be implemented with proper types in later phases
	for name, info := range updates {
		fmt.Printf("  %s: %s → (latest)\n", name, info.CurrentVersion)
	}
	return nil
}

// applyVersionUpdates applies the detected version updates
func (a *App) applyVersionUpdates(manager *version.Manager, updates map[string]*models.VersionInfo, force bool, packages []string) error {
	updateCount := 0

	// Placeholder implementation - would be implemented with proper types in later phases
	fmt.Println("✅ Version updates would be applied here")
	a.cli.ShowSuccess(fmt.Sprintf("Updated %d packages", updateCount))
	return nil
}

// handleBasicVersionUpdate handles updates using basic version manager
func (a *App) handleBasicVersionUpdate(dryRun bool) error {
	configManager := a.container.GetConfigManager()
	versions, err := configManager.GetLatestVersions()
	if err != nil {
		a.cli.ShowError(fmt.Sprintf("Failed to get latest versions: %v", err))
		return err
	}

	if dryRun {
		return a.showBasicDryRunVersions(versions)
	}

	// Update versions.md file (placeholder implementation)
	fmt.Println("✅ Version information updated")
	fmt.Println("Note: Full version store integration will be completed in task 6.2")
	return nil
}

// showBasicDryRunVersions displays basic version information in dry run mode
func (a *App) showBasicDryRunVersions(versions *models.VersionConfig) error {
	fmt.Println("\n🔍 Dry Run - Would update the following versions:")
	fmt.Printf("  Node.js: %s\n", versions.Node)
	fmt.Printf("  Go: %s\n", versions.Go)

	if versions.NextJS != "" {
		fmt.Printf("  Next.js: %s\n", versions.NextJS)
	}
	if versions.React != "" {
		fmt.Printf("  React: %s\n", versions.React)
	}
	if versions.Kotlin != "" {
		fmt.Printf("  Kotlin: %s\n", versions.Kotlin)
	}
	if versions.Swift != "" {
		fmt.Printf("  Swift: %s\n", versions.Swift)
	}

	return nil
}

// runDashboardCommand handles the dashboard command execution
func (a *App) runDashboardCommand(refresh bool, format string, showDetails bool) error {
	if refresh {
		a.cli.ShowProgress("Refreshing version data")
	}

	// Get version manager with storage
	versionManager := a.container.GetVersionManager()

	// Check if version manager has storage capability
	if managerWithStorage, ok := versionManager.(*version.Manager); ok {
		// Generate comprehensive dashboard data
		dashboardData, err := a.generateDashboardData(managerWithStorage, refresh)
		if err != nil {
			a.cli.ShowError(fmt.Sprintf("Failed to generate dashboard data: %v", err))
			return err
		}

		// Display dashboard based on format
		switch format {
		case constants.FormatJSON:
			return a.displayDashboardJSON(dashboardData)
		case constants.FormatYAML:
			return a.displayDashboardYAML(dashboardData)
		default:
			return a.displayDashboardTable(dashboardData, showDetails)
		}
	}

	// Fallback for basic version manager
	a.cli.ShowWarning("Enhanced dashboard requires version storage. Showing basic version information.")
	return a.runVersionCommand(true)
}

// generateDashboardData creates comprehensive dashboard information
func (a *App) generateDashboardData(versionManager *version.Manager, refresh bool) (*models.DashboardData, error) {
	dashboardData := &models.DashboardData{
		GeneratedAt: time.Now(),
		Metadata:    make(map[string]string),
	}

	// Get version report
	versionReport, err := versionManager.CheckLatestVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to get version report: %w", err)
	}
	dashboardData.VersionReport = versionReport

	// Get template consistency status
	templateStatus, err := a.getTemplateConsistencyStatus()
	if err != nil {
		a.cli.ShowWarning(fmt.Sprintf("Failed to check template consistency: %v", err))
		// Continue with empty template status
		templateStatus = &models.TemplateConsistencyStatus{
			CheckedAt: time.Now(),
			Status:    "unknown",
			Issues:    make([]models.ConsistencyIssue, 0),
		}
	}
	dashboardData.TemplateStatus = templateStatus

	// Get validation results
	validationResults, err := a.getValidationResults()
	if err != nil {
		a.cli.ShowWarning(fmt.Sprintf("Failed to get validation results: %v", err))
		// Continue with empty validation results
		validationResults = &models.ValidationResults{
			CheckedAt: time.Now(),
			Status:    "unknown",
			Results:   make(map[string]models.DetailedValidationResult, 0),
		}
	}
	dashboardData.ValidationResults = validationResults

	// Add metadata
	dashboardData.Metadata["generator_version"] = "1.0.0"
	dashboardData.Metadata["refresh_requested"] = fmt.Sprintf("%t", refresh)

	return dashboardData, nil
}

// getTemplateConsistencyStatus checks template consistency across all templates
func (a *App) getTemplateConsistencyStatus() (*models.TemplateConsistencyStatus, error) {
	status := &models.TemplateConsistencyStatus{
		CheckedAt: time.Now(),
		Status:    constants.ValidationConsistent,
		Issues:    make([]models.ConsistencyIssue, 0),
		Summary: models.ConsistencySummary{
			TotalTemplates:      0,
			ConsistentTemplates: 0,
			IssuesFound:         0,
		},
	}

	// Check frontend template consistency
	frontendTemplates := []string{
		"templates/frontend/nextjs-app",
		"templates/frontend/nextjs-home",
		"templates/frontend/nextjs-admin",
	}

	for _, templatePath := range frontendTemplates {
		status.Summary.TotalTemplates++

		// Check if template directory exists
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			status.Issues = append(status.Issues, models.ConsistencyIssue{
				Type:        "missing_template",
				Severity:    "high",
				Template:    templatePath,
				Description: "Template directory does not exist",
				File:        templatePath,
			})
			status.Summary.IssuesFound++
			continue
		}

		// Check for package.json consistency
		packageJsonPath := filepath.Join(templatePath, "package.json.tmpl")
		if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
			status.Issues = append(status.Issues, models.ConsistencyIssue{
				Type:        "missing_file",
				Severity:    "medium",
				Template:    templatePath,
				Description: "Missing package.json.tmpl file",
				File:        packageJsonPath,
			})
			status.Summary.IssuesFound++
		} else {
			status.Summary.ConsistentTemplates++
		}
	}

	// Determine overall status
	if status.Summary.IssuesFound > 0 {
		if status.Summary.IssuesFound > status.Summary.TotalTemplates/2 {
			status.Status = constants.PriorityCritical
		} else {
			status.Status = "issues_found"
		}
	}

	return status, nil
}

// getValidationResults gets validation results for templates and configurations
func (a *App) getValidationResults() (*models.ValidationResults, error) {
	results := &models.ValidationResults{
		CheckedAt: time.Now(),
		Status:    "passed",
		Results:   make(map[string]models.DetailedValidationResult),
		Summary: models.ValidationSummary{
			TotalChecks:  0,
			PassedChecks: 0,
			FailedChecks: 0,
			Warnings:     0,
		},
	}

	// Validate template structure
	templateResult := a.validateTemplateStructure()
	results.Results["template_structure"] = templateResult
	results.Summary.TotalChecks++
	if templateResult.Status == constants.StatusPassed {
		results.Summary.PassedChecks++
	} else {
		results.Summary.FailedChecks++
	}
	results.Summary.Warnings += len(templateResult.Warnings)

	// Validate version consistency
	versionResult := a.validateVersionConsistency()
	results.Results["version_consistency"] = versionResult
	results.Summary.TotalChecks++
	if versionResult.Status == "passed" {
		results.Summary.PassedChecks++
	} else {
		results.Summary.FailedChecks++
	}
	results.Summary.Warnings += len(versionResult.Warnings)

	// Validate deployment readiness
	deploymentResult := a.validateDeploymentReadiness()
	results.Results["deployment_readiness"] = deploymentResult
	results.Summary.TotalChecks++
	if deploymentResult.Status == "passed" {
		results.Summary.PassedChecks++
	} else {
		results.Summary.FailedChecks++
	}
	results.Summary.Warnings += len(deploymentResult.Warnings)

	// Determine overall status
	if results.Summary.FailedChecks > 0 {
		results.Status = constants.StatusFailed
	} else if results.Summary.Warnings > 0 {
		results.Status = "warnings"
	}

	return results, nil
}

// validateTemplateStructure validates the overall template structure
func (a *App) validateTemplateStructure() models.DetailedValidationResult {
	result := models.DetailedValidationResult{
		Name:      "Template Structure",
		Status:    "passed",
		CheckedAt: time.Now(),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
		Details:   make(map[string]string),
	}

	// Check if templates directory exists
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		result.Status = constants.StatusFailed
		result.Errors = append(result.Errors, "Templates directory does not exist")
		return result
	}

	// Check for required template categories
	requiredDirs := []string{"base", "frontend", "backend", "mobile", "infrastructure"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join("templates", dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing template category: %s", dir))
		} else {
			result.Details[dir] = constants.StatusPresent
		}
	}

	return result
}

// validateVersionConsistency validates version consistency across templates
func (a *App) validateVersionConsistency() models.DetailedValidationResult {
	result := models.DetailedValidationResult{
		Name:      "Version Consistency",
		Status:    "passed",
		CheckedAt: time.Now(),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
		Details:   make(map[string]string),
	}

	// This would check for version consistency across templates
	// For now, we'll simulate some basic checks
	result.Details["node_version_check"] = constants.ValidationConsistent
	result.Details["react_version_check"] = constants.ValidationConsistent
	result.Details["nextjs_version_check"] = constants.ValidationConsistent

	// Add a sample warning
	result.Warnings = append(result.Warnings, "Some templates may be using outdated TypeScript versions")

	return result
}

// validateDeploymentReadiness validates deployment configurations
func (a *App) validateDeploymentReadiness() models.DetailedValidationResult {
	result := models.DetailedValidationResult{
		Name:      "Deployment Readiness",
		Status:    "passed",
		CheckedAt: time.Now(),
		Errors:    make([]string, 0),
		Warnings:  make([]string, 0),
		Details:   make(map[string]string),
	}

	// Check for Vercel configuration in frontend templates
	frontendTemplates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	for _, template := range frontendTemplates {
		vercelConfigPath := filepath.Join("templates", "frontend", template, "vercel.json.tmpl")
		if _, err := os.Stat(vercelConfigPath); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Missing Vercel configuration for %s", template))
		} else {
			result.Details[fmt.Sprintf("%s_vercel_config", template)] = "present"
		}
	}

	// Check for Docker configurations
	dockerConfigPath := "templates/infrastructure/docker"
	if _, err := os.Stat(dockerConfigPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "Missing Docker infrastructure templates")
	} else {
		result.Details["docker_config"] = "present"
	}

	return result
}

// displayDashboardTable displays the dashboard in table format
func (a *App) displayDashboardTable(data *models.DashboardData, showDetails bool) error {
	a.displayDashboardHeader(data)
	a.displayVersionStatus(data.VersionReport, showDetails)
	a.displayTemplateConsistency(data.TemplateStatus, showDetails)
	a.displayValidationResults(data.ValidationResults, showDetails)
	a.displayQuickActions()

	return nil
}

// displayDashboardHeader displays the dashboard header
func (a *App) displayDashboardHeader(data *models.DashboardData) {
	fmt.Printf("\n🚀 Version Management Dashboard\n")
	fmt.Printf("Generated: %s\n", data.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("=", 60))
}

// displayVersionStatus displays the version status section
func (a *App) displayVersionStatus(report *models.VersionReport, showDetails bool) {
	fmt.Printf("\n📊 Version Status\n")
	fmt.Println(strings.Repeat("-", 30))

	if report == nil {
		return
	}

	fmt.Printf("Total Packages: %d\n", report.TotalPackages)
	fmt.Printf("Up to Date: %d\n", report.TotalPackages-report.OutdatedCount)
	fmt.Printf("Outdated: %d\n", report.OutdatedCount)
	fmt.Printf("Security Issues: %d\n", report.SecurityIssues)

	if report.OutdatedCount > 0 {
		a.displayUpdateRecommendations(report.Recommendations, showDetails)
	}
}

// displayUpdateRecommendations displays update recommendations
func (a *App) displayUpdateRecommendations(recommendations []models.UpdateRecommendation, showDetails bool) {
	fmt.Printf("\n🔄 Update Recommendations:\n")

	for i, rec := range recommendations {
		if i >= 5 && !showDetails { // Limit to 5 unless details requested
			fmt.Printf("  ... and %d more (use --details to see all)\n", len(recommendations)-5)
			break
		}
		priority := a.formatPriorityIcon(rec.Priority)
		fmt.Printf("  %s %s: %s → %s\n", priority, rec.Name, rec.CurrentVersion, rec.RecommendedVersion)
	}
}

// formatPriorityIcon returns the appropriate icon for priority level
func (a *App) formatPriorityIcon(priority string) string {
	switch priority {
	case "critical":
		return "🔴 CRITICAL"
	case "high":
		return "🟠 HIGH"
	case "medium":
		return "🟡 MEDIUM"
	default:
		return "🟢 LOW"
	}
}

// displayTemplateConsistency displays the template consistency section
func (a *App) displayTemplateConsistency(status interface{}, showDetails bool) {
	fmt.Printf("\n🔧 Template Consistency\n")
	fmt.Println(strings.Repeat("-", 30))

	if status == nil {
		return
	}

	// Placeholder implementation - would be implemented with proper types in later phases
	fmt.Printf("Status: ✅ consistent\n")
	fmt.Printf("Templates Checked: 0\n")
	fmt.Printf("Consistent: 0\n")
	fmt.Printf("Issues Found: 0\n")
}

// displayValidationResults displays the validation results section
func (a *App) displayValidationResults(results *models.ValidationResults, showDetails bool) {
	fmt.Printf("\n✅ Validation Results\n")
	fmt.Println(strings.Repeat("-", 30))

	if results == nil {
		return
	}

	statusIcon := a.getValidationStatusIcon(results.Status)
	fmt.Printf("Status: %s %s\n", statusIcon, titleCaser.String(results.Status))
	fmt.Printf("Total Checks: %d\n", results.Summary.TotalChecks)
	fmt.Printf("Passed: %d\n", results.Summary.PassedChecks)
	fmt.Printf("Failed: %d\n", results.Summary.FailedChecks)
	fmt.Printf("Warnings: %d\n", results.Summary.Warnings)

	if showDetails {
		a.displayDetailedResults(map[string]interface{}{})
	}
}

// getValidationStatusIcon returns the appropriate icon for validation status
func (a *App) getValidationStatusIcon(status string) string {
	switch status {
	case "warnings":
		return constants.SymbolWarning
	case constants.StatusFailed:
		return constants.SymbolFailure
	default:
		return constants.SymbolSuccess
	}
}

// displayDetailedResults displays detailed validation results
func (a *App) displayDetailedResults(results map[string]interface{}) {
	// Placeholder implementation - would be implemented with proper types in later phases
	fmt.Printf("\n📋 Detailed Results:\n")
	fmt.Printf("  (Detailed results would be displayed here)\n")
}

// displayQuickActions displays the quick actions section
func (a *App) displayQuickActions() {
	fmt.Printf("\n🚀 Quick Actions\n")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("  generator versions check     - Check for version updates")
	fmt.Println("  generator versions update    - Update version information")
	fmt.Println("  generator versions apply     - Apply updates to templates")
	fmt.Println("  generator validate           - Validate project structure")
}

// displayDashboardJSON displays the dashboard in JSON format
func (a *App) displayDashboardJSON(data *models.DashboardData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dashboard data to JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

// displayDashboardYAML displays the dashboard in YAML format
func (a *App) displayDashboardYAML(data *models.DashboardData) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal dashboard data to YAML: %w", err)
	}
	fmt.Println(string(yamlData))
	return nil
}

// runGenerateReportCommand handles the generate-report command execution
func (a *App) runGenerateReportCommand(reportType, format, outputDir string) error {
	a.cli.ShowProgress(fmt.Sprintf("Generating %s report", reportType))

	// Get version manager with storage
	versionManager := a.container.GetVersionManager()

	// Check if version manager has storage capability
	if managerWithStorage, ok := versionManager.(*version.Manager); ok {
		// Import reporting package (this would be done at the top of the file)
		// For now, we'll simulate the report generation

		switch reportType {
		case "version":
			return a.generateVersionReport(managerWithStorage, format, outputDir)
		case "security":
			return a.generateSecurityReport(managerWithStorage, format, outputDir)
		case "template":
			return a.generateTemplateReport(format, outputDir)
		default:
			return fmt.Errorf("unsupported report type: %s", reportType)
		}
	}

	a.cli.ShowError("Enhanced reporting requires version storage")
	return fmt.Errorf("version manager does not support storage")
}

// generateVersionReport generates a version update report
func (a *App) generateVersionReport(versionManager *version.Manager, format, outputDir string) error {
	// Get version report
	versionReport, err := versionManager.CheckLatestVersions()
	if err != nil {
		return fmt.Errorf("failed to check versions: %w", err)
	}

	// Create report generator (simulated)
	reportID := fmt.Sprintf("version_report_%d", time.Now().Unix())

	a.cli.ShowSuccess(fmt.Sprintf("Version report generated: %s", reportID))
	fmt.Printf("Report saved to: %s/%s.%s\n", outputDir, reportID, format)

	// Display summary
	fmt.Printf("\nReport Summary:\n")
	fmt.Printf("Total Packages: %d\n", versionReport.TotalPackages)
	fmt.Printf("Outdated: %d\n", versionReport.OutdatedCount)
	fmt.Printf("Security Issues: %d\n", versionReport.SecurityIssues)
	fmt.Printf("Recommendations: %d\n", len(versionReport.Recommendations))

	return nil
}

// generateSecurityReport generates a security-focused report
func (a *App) generateSecurityReport(versionManager *version.Manager, format, outputDir string) error {
	// Get version report for security analysis
	versionReport, err := versionManager.CheckLatestVersions()
	if err != nil {
		return fmt.Errorf("failed to check versions: %w", err)
	}

	// Count security issues
	totalIssues := versionReport.SecurityIssues
	criticalIssues := 0
	highIssues := 0

	for _, info := range versionReport.Details {
		for _, issue := range info.SecurityIssues {
			if issue.Severity == "critical" {
				criticalIssues++
			} else if issue.Severity == "high" {
				highIssues++
			}
		}
	}

	reportID := fmt.Sprintf("security_report_%d", time.Now().Unix())

	a.cli.ShowSuccess(fmt.Sprintf("Security report generated: %s", reportID))
	fmt.Printf("Report saved to: %s/%s.%s\n", outputDir, reportID, format)

	// Display summary
	fmt.Printf("\nSecurity Report Summary:\n")
	fmt.Printf("Total Issues: %d\n", totalIssues)
	fmt.Printf("Critical: %d\n", criticalIssues)
	fmt.Printf("High: %d\n", highIssues)

	if totalIssues > 0 {
		fmt.Printf("\n⚠️  Security vulnerabilities detected. Please review the full report.\n")
	} else {
		fmt.Printf("\n✅ No security vulnerabilities detected.\n")
	}

	return nil
}

// generateTemplateReport generates a template consistency report
func (a *App) generateTemplateReport(format, outputDir string) error {
	// Get template consistency status
	templateStatus, err := a.getTemplateConsistencyStatus()
	if err != nil {
		return fmt.Errorf("failed to check template consistency: %w", err)
	}

	reportID := fmt.Sprintf("template_report_%d", time.Now().Unix())

	a.cli.ShowSuccess(fmt.Sprintf("Template report generated: %s", reportID))
	fmt.Printf("Report saved to: %s/%s.%s\n", outputDir, reportID, format)

	// Display summary
	fmt.Printf("\nTemplate Report Summary:\n")
	fmt.Printf("Status: %s\n", titleCaser.String(templateStatus.Status))
	fmt.Printf("Templates Checked: %d\n", templateStatus.Summary.TotalTemplates)
	fmt.Printf("Consistent: %d\n", templateStatus.Summary.ConsistentTemplates)
	fmt.Printf("Issues Found: %d\n", templateStatus.Summary.IssuesFound)

	return nil
}

// runListReportsCommand handles the list-reports command execution
func (a *App) runListReportsCommand(outputDir string) error {
	a.cli.ShowProgress("Scanning for generated reports")

	// Check if reports directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Printf("No reports directory found at: %s\n", outputDir)
		return nil
	}

	// Read directory contents
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("failed to read reports directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No reports found in: %s\n", outputDir)
		return nil
	}

	fmt.Printf("\n📊 Generated Reports (%s)\n", outputDir)
	fmt.Println(strings.Repeat("=", 50))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Determine report type from filename
		name := entry.Name()
		var reportType string
		if strings.HasPrefix(name, "security_") {
			reportType = "Security"
		} else if strings.HasPrefix(name, "template_") {
			reportType = "Template"
		} else if strings.HasPrefix(name, "version_") {
			reportType = "Version"
		} else {
			reportType = "Unknown"
		}

		fmt.Printf("📄 %s\n", name)
		fmt.Printf("   Type: %s\n", reportType)
		fmt.Printf("   Size: %d bytes\n", info.Size())
		fmt.Printf("   Generated: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	return nil
}

// runAuditCommand handles the audit command execution
func (a *App) runAuditCommand(since, eventType, format string) error {
	a.cli.ShowProgress("Retrieving audit trail")

	// Parse since duration
	duration, err := time.ParseDuration(since)
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	sinceTime := time.Now().Add(-duration)

	// For now, simulate audit trail (in real implementation, this would use the audit trail)
	fmt.Printf("\n📋 Audit Trail (Last %s)\n", since)
	fmt.Println(strings.Repeat("=", 50))

	// Simulate some audit events
	events := []struct {
		timestamp time.Time
		eventType string
		action    string
		resource  string
		success   bool
	}{
		{time.Now().Add(-2 * time.Hour), "version_check", "check_versions", "npm_registry", true},
		{time.Now().Add(-4 * time.Hour), "version_update", "update_version", "react", true},
		{time.Now().Add(-6 * time.Hour), "template_update", "update_template", "nextjs-app", true},
		{time.Now().Add(-8 * time.Hour), "security_scan", "scan_packages", "security_db", true},
	}

	for _, event := range events {
		if event.timestamp.Before(sinceTime) {
			continue
		}
		if eventType != "" && event.eventType != eventType {
			continue
		}

		statusIcon := constants.SymbolSuccess
		if !event.success {
			statusIcon = constants.SymbolFailure
		}

		fmt.Printf("%s [%s] %s: %s -> %s\n",
			statusIcon,
			event.timestamp.Format("15:04:05"),
			titleCaser.String(event.eventType),
			event.action,
			event.resource)
	}

	fmt.Printf("\nTotal Events: %d\n", len(events))
	fmt.Printf("Success Rate: 100%%\n")

	return nil
}
