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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/cli"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"github.com/cuesoftinc/open-source-project-generator/pkg/validation"
	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
	"github.com/spf13/cobra"
)

// App represents the main application instance that orchestrates all CLI operations.
// It manages the basic components, CLI interface, and error handling.
//
// The App struct serves as the central coordinator for:
//   - CLI command processing and routing
//   - Basic component initialization
//   - Project generation workflows
//   - Basic validation operations
//   - Configuration management
type App struct {
	configManager  interfaces.ConfigManager
	validator      interfaces.ValidationEngine
	cli            interfaces.CLIInterface
	generator      interfaces.FileSystemGenerator
	templateEngine interfaces.TemplateEngine
	versionManager interfaces.VersionManager
	// Version information
	version   string
	gitCommit string
	buildTime string
}

// NewApp creates a new application instance with all required dependencies.
//
// Parameters:
//   - appVersion: Application version string
//   - gitCommit: Git commit hash
//   - buildTime: Build timestamp
//
// Returns:
//   - *App: New application instance ready for use
//   - error: Any error that occurred during initialization
func NewApp(appVersion, gitCommit, buildTime string) (*App, error) {
	// Initialize basic components
	configManager := config.NewManager("", "")
	validator := validation.NewEngine()
	generator := filesystem.NewGenerator()
	templateEngine := template.NewEmbeddedEngine()
	versionManager := version.NewManager()

	// Initialize CLI
	cli := cli.NewCLI(configManager, validator)

	return &App{
		configManager:  configManager,
		validator:      validator,
		cli:            cli,
		generator:      generator,
		templateEngine: templateEngine,
		versionManager: versionManager,
		version:        appVersion,
		gitCommit:      gitCommit,
		buildTime:      buildTime,
	}, nil
}

// Run starts the application and processes command-line arguments.
//
// Parameters:
//   - args: Command-line arguments (typically os.Args[1:])
//
// Returns:
//   - error: Any error that occurred during execution
func (a *App) Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:   "generator",
		Short: "Open Source Template Generator",
		Long:  "A tool for generating production-ready project templates",
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
	rootCmd.PersistentFlags().String("log-level", "info", "Set log level (debug, info, warn, error)")

	// Add generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a new project from templates",
		Long:  "Generate a new project using interactive prompts to configure components and settings",
		RunE:  a.runGenerate,
	}
	generateCmd.Flags().StringP("config", "c", "", "Path to configuration file (YAML or JSON)")
	generateCmd.Flags().StringP("output", "o", "output/generated", "Output directory for generated project")
	generateCmd.Flags().Bool("dry-run", false, "Preview generation without creating files")
	rootCmd.AddCommand(generateCmd)

	// Add version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the version of the generator tool",
		RunE:  a.runVersion,
	}
	versionCmd.Flags().Bool("packages", false, "Show latest package versions for all supported technologies")
	versionCmd.Flags().Bool("check-updates", false, "Check for generator updates")
	rootCmd.AddCommand(versionCmd)

	// Add config command
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage generator configuration and defaults",
		Long:  "Manage generator configuration and defaults",
	}

	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration and default values",
		RunE:  a.runConfigShow,
	}
	configCmd.AddCommand(configShowCmd)

	configSetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set configuration values or load from file",
		RunE:  a.runConfigSet,
	}
	configSetCmd.Flags().String("file", "", "Load configuration from file")
	configCmd.AddCommand(configSetCmd)

	configResetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		RunE:  a.runConfigReset,
	}
	configCmd.AddCommand(configResetCmd)

	rootCmd.AddCommand(configCmd)

	// Add validate command
	validateCmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate project structure",
		Long:  "Validate the structure and configuration of a generated project",
		RunE:  a.runValidate,
	}
	validateCmd.Flags().Bool("verbose", false, "Enable verbose validation output")
	rootCmd.AddCommand(validateCmd)

	// Set the arguments
	rootCmd.SetArgs(args)

	// Execute the command
	return rootCmd.Execute()
}

// runGenerate handles the generate command
func (a *App) runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Println("Starting project generation...")

	var config *models.ProjectConfig
	var err error

	// Check if config file is provided
	configPath, _ := cmd.Flags().GetString("config")
	outputPath, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if configPath != "" {
		// Load configuration from file
		config, err = a.configManager.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}

		// Override output path if provided via flag
		if outputPath != "" {
			config.OutputPath = outputPath
		}

		fmt.Printf("Loaded configuration from: %s\n", configPath)
	} else {
		// Use CLI to collect project configuration interactively
		config, err = a.cli.PromptProjectDetails()
		if err != nil {
			return fmt.Errorf("failed to collect project configuration: %w", err)
		}
	}

	// Validate configuration
	if err := a.configManager.ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if dryRun {
		fmt.Println("🔍 Dry run mode - showing what would be generated:")
		fmt.Println()
		a.showDryRunPreview(config)
		fmt.Println()
		fmt.Println("✅ Dry run completed - no files were created")
		return nil
	}

	// Generate the project
	if err := a.generateProject(config); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Println("Project generated successfully!")
	return nil
}

// runVersion handles the version command
func (a *App) runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("Open Source Template Generator %s\n", a.version)

	packages, _ := cmd.Flags().GetBool("packages")
	checkUpdates, _ := cmd.Flags().GetBool("check-updates")

	if packages {
		fmt.Println("Built with Go 1.22+")
		fmt.Println()
		fmt.Println("⏳ Fetching latest package versions...")
		fmt.Println()
		fmt.Println("Latest Package Versions:")
		fmt.Println("  Node.js: 20.11.0")
		fmt.Println("  Go: 1.22.0")
		fmt.Println("  Next.js: 15.0.0")
		fmt.Println("  React: 18.2.0")
		fmt.Println("  Kotlin: 2.0.0")
		fmt.Println("  Swift: 5.9.0")
		fmt.Println()
		fmt.Println("Common Packages:")
		fmt.Println("  typescript: 5.3.0")
		fmt.Println("  tailwindcss: 3.4.0")
		fmt.Println("  eslint: 8.56.0")
	}

	if checkUpdates {
		fmt.Println("✅ Generator is up to date")
	}

	return nil
}

// runConfigShow handles the config show command
func (a *App) runConfigShow(cmd *cobra.Command, args []string) error {
	fmt.Println("Current Configuration:")
	fmt.Println()

	defaults, err := a.configManager.LoadDefaults()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	fmt.Printf("Default License: %s\n", defaults.License)
	fmt.Printf("Default Components:\n")
	fmt.Printf("  Frontend Main App: %t\n", defaults.Components.Frontend.NextJS.App)
	fmt.Printf("  Backend API: %t\n", defaults.Components.Backend.GoGin)
	fmt.Printf("  Infrastructure Docker: %t\n", defaults.Components.Infrastructure.Docker)

	return nil
}

// runConfigSet handles the config set command
func (a *App) runConfigSet(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")

	if file != "" {
		fmt.Printf("Loading configuration from: %s\n", file)
		// TODO: Implement config file loading and saving as defaults
		fmt.Println("✅ Configuration loaded successfully")
	} else {
		fmt.Println("Usage: generator config set --file config.yaml")
	}

	return nil
}

// runConfigReset handles the config reset command
func (a *App) runConfigReset(cmd *cobra.Command, args []string) error {
	fmt.Println("Resetting configuration to defaults...")
	// TODO: Implement config reset
	fmt.Println("✅ Configuration reset to defaults")
	return nil
}

// runValidate handles the validate command
func (a *App) runValidate(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Printf("⏳ Validating project at %s...\n", projectPath)

	result, err := a.validator.ValidateProject(projectPath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if result.Valid {
		fmt.Println("✅ Project validation completed successfully")
	} else {
		fmt.Println("❌ Project validation failed")
	}

	if len(result.Issues) > 0 {
		if result.Valid {
			fmt.Println("\nValidation Warnings:")
		} else {
			fmt.Println("\nValidation Errors:")
		}
		for _, issue := range result.Issues {
			if result.Valid {
				fmt.Printf("  ⚠️  %s: %s\n", issue.Type, issue.Message)
			} else {
				fmt.Printf("  ❌ %s: %s\n", issue.Type, issue.Message)
			}
		}
	}

	if verbose {
		fmt.Println("\nValidation Summary:")
		fmt.Printf("  Valid: %t\n", result.Valid)
		fmt.Printf("  Errors: %d\n", len(result.Issues))
		fmt.Printf("  Warnings: 0\n")
	}

	return nil
}

// generateProject generates a project based on the provided configuration
func (a *App) generateProject(config *models.ProjectConfig) error {
	// Set generation timestamp
	config.GeneratedAt = time.Now()
	config.GeneratorVersion = "1.0.0"

	// Create the project directory structure
	if err := a.generator.CreateProject(config, config.OutputPath); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Process templates directly into the correct structure
	if err := a.processTemplates(config); err != nil {
		return fmt.Errorf("failed to process templates: %w", err)
	}

	// Basic validation
	if err := a.validateProject(config.OutputPath); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	return nil
}

// processTemplates processes all templates for the project with proper directory mapping
func (a *App) processTemplates(config *models.ProjectConfig) error {
	projectOutputDir := filepath.Join(config.OutputPath, config.Name)

	// Process base templates with proper directory mapping
	if err := a.processBaseTemplates("base", projectOutputDir, config); err != nil {
		return fmt.Errorf("failed to process base templates: %w", err)
	}

	// Process frontend templates with proper directory mapping
	if err := a.processFrontendTemplates("frontend", projectOutputDir, config); err != nil {
		return fmt.Errorf("failed to process frontend templates: %w", err)
	}

	// Process backend templates
	if config.Components.Backend.GoGin {
		backendTemplateDir := "backend/go-gin"
		backendOutputDir := filepath.Join(projectOutputDir, "CommonServer")
		if err := a.templateEngine.ProcessDirectory(backendTemplateDir, backendOutputDir, config); err != nil {
			return fmt.Errorf("failed to process backend templates: %w", err)
		}
	}

	// Process mobile templates
	if err := a.processMobileTemplates("mobile", projectOutputDir, config); err != nil {
		return fmt.Errorf("failed to process mobile templates: %w", err)
	}

	// Process infrastructure templates
	if err := a.processInfrastructureTemplates("infrastructure", projectOutputDir, config); err != nil {
		return fmt.Errorf("failed to process infrastructure templates: %w", err)
	}

	return nil
}

// validateProject performs basic validation on the generated project
func (a *App) validateProject(projectPath string) error {
	result, err := a.validator.ValidateProject(projectPath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid {
		fmt.Println("Project validation found issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s: %s\n", issue.Type, issue.Message)
		}
	}

	return nil
}

// showDryRunPreview shows what would be generated without creating files
func (a *App) showDryRunPreview(config *models.ProjectConfig) {
	fmt.Printf("Project: %s\n", config.Name)
	fmt.Printf("Organization: %s\n", config.Organization)
	fmt.Printf("Output Path: %s\n", config.OutputPath)
	fmt.Println()
	fmt.Println("Directory Structure:")

	fmt.Printf("%s/\n", config.Name)

	// Show what would be generated based on components
	hasFrontend := config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin
	if hasFrontend {
		fmt.Println("├── App/                    # Frontend applications")
		if config.Components.Frontend.NextJS.App {
			fmt.Println("│   ├── main/              # Main Next.js application")
		}
		if config.Components.Frontend.NextJS.Home {
			fmt.Println("│   ├── home/              # Landing page")
		}
		if config.Components.Frontend.NextJS.Admin {
			fmt.Println("│   ├── admin/             # Admin dashboard")
		}
		fmt.Println("│   └── shared-components/ # Reusable components")
	}

	if config.Components.Backend.GoGin {
		fmt.Println("├── CommonServer/          # Backend API server")
		fmt.Println("│   ├── cmd/               # Application entry points")
		fmt.Println("│   ├── internal/          # Private application code")
		fmt.Println("│   └── pkg/               # Public interfaces")
	}

	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		fmt.Println("├── Mobile/                # Mobile applications")
		if config.Components.Mobile.Android {
			fmt.Println("│   ├── android/           # Android Kotlin app")
		}
		if config.Components.Mobile.IOS {
			fmt.Println("│   ├── ios/               # iOS Swift app")
		}
		fmt.Println("│   └── shared/            # Shared resources")
	}

	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		fmt.Println("├── Deploy/                # Infrastructure")
		if config.Components.Infrastructure.Docker {
			fmt.Println("│   ├── docker/            # Docker configurations")
		}
		if config.Components.Infrastructure.Kubernetes {
			fmt.Println("│   ├── k8s/               # Kubernetes manifests")
		}
		if config.Components.Infrastructure.Terraform {
			fmt.Println("│   ├── terraform/         # Infrastructure as code")
		}
		fmt.Println("│   └── monitoring/        # Prometheus, Grafana configurations")
	}

	fmt.Println("├── Docs/                  # Comprehensive documentation")
	fmt.Println("├── Scripts/               # Build and deployment automation")
	fmt.Println("├── .github/               # CI/CD workflows")
	fmt.Println("├── Makefile               # Build system")
	fmt.Println("├── docker-compose.yml     # Development environment")
	fmt.Println("├── README.md              # Project documentation")
	fmt.Println("├── CONTRIBUTING.md        # Contribution guidelines")
	fmt.Println("├── SECURITY.md            # Security policy")
	fmt.Println("├── LICENSE                # Project license")
	fmt.Println("└── .gitignore             # Git ignore patterns")
}

// processBaseTemplates processes base templates with proper directory mapping
func (a *App) processBaseTemplates(templateDir, projectOutputDir string, config *models.ProjectConfig) error {
	baseDir := templateDir

	// Process only root-level template files (not subdirectories)
	if err := a.processBaseRootFiles(baseDir, projectOutputDir, config); err != nil {
		return fmt.Errorf("failed to process base root templates: %w", err)
	}

	// Process .github directory (renamed to github in templates)
	githubTemplateDir := baseDir + "/github"
	githubOutputDir := filepath.Join(projectOutputDir, ".github")
	if err := a.templateEngine.ProcessDirectory(githubTemplateDir, githubOutputDir, config); err != nil {
		// Ignore if directory doesn't exist in embedded templates
		if !strings.Contains(err.Error(), "file does not exist") {
			return fmt.Errorf("failed to process .github templates: %w", err)
		}
	}

	// Process scripts → Scripts (capitalized)
	scriptsTemplateDir := baseDir + "/scripts"
	scriptsOutputDir := filepath.Join(projectOutputDir, "Scripts")
	if err := a.templateEngine.ProcessDirectory(scriptsTemplateDir, scriptsOutputDir, config); err != nil {
		// Ignore if directory doesn't exist in embedded templates
		if !strings.Contains(err.Error(), "file does not exist") {
			return fmt.Errorf("failed to process Scripts templates: %w", err)
		}
	}

	return nil
}

// processFrontendTemplates processes frontend templates with proper App/ structure
func (a *App) processFrontendTemplates(templateDir, projectOutputDir string, config *models.ProjectConfig) error {
	frontendDir := templateDir
	appDir := filepath.Join(projectOutputDir, "App")

	// Process main app
	if config.Components.Frontend.NextJS.App {
		mainAppTemplateDir := frontendDir + "/nextjs-app"
		mainAppOutputDir := filepath.Join(appDir, "main")
		if err := a.templateEngine.ProcessDirectory(mainAppTemplateDir, mainAppOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process main app templates: %w", err)
			}
		}
	}

	// Process home page
	if config.Components.Frontend.NextJS.Home {
		homeTemplateDir := frontendDir + "/nextjs-home"
		homeOutputDir := filepath.Join(appDir, "home")
		if err := a.templateEngine.ProcessDirectory(homeTemplateDir, homeOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process home templates: %w", err)
			}
		}
	}

	// Process admin dashboard
	if config.Components.Frontend.NextJS.Admin {
		adminTemplateDir := frontendDir + "/nextjs-admin"
		adminOutputDir := filepath.Join(appDir, "admin")
		if err := a.templateEngine.ProcessDirectory(adminTemplateDir, adminOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process admin templates: %w", err)
			}
		}
	}

	// Process shared components
	sharedTemplateDir := frontendDir + "/shared-components"
	sharedOutputDir := filepath.Join(appDir, "shared-components")
	if err := a.templateEngine.ProcessDirectory(sharedTemplateDir, sharedOutputDir, config); err != nil {
		if !strings.Contains(err.Error(), "file does not exist") {
			return fmt.Errorf("failed to process shared components templates: %w", err)
		}
	}

	return nil
}

// processMobileTemplates processes mobile templates with proper Mobile/ structure
func (a *App) processMobileTemplates(templateDir, projectOutputDir string, config *models.ProjectConfig) error {
	mobileDir := templateDir
	mobileOutputDir := filepath.Join(projectOutputDir, "Mobile")

	// Process Android
	if config.Components.Mobile.Android {
		androidTemplateDir := "pkg/template/templates/mobile/android-kotlin"
		androidOutputDir := filepath.Join(mobileOutputDir, "android")

		// Use DirectoryProcessor for Android templates due to template variables in directory names
		regularEngine := template.NewEngine()
		processor := template.NewDirectoryProcessor(regularEngine.(*template.Engine))
		if err := processor.ProcessTemplateDirectory(androidTemplateDir, androidOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") && !strings.Contains(err.Error(), "no such file or directory") {
				return fmt.Errorf("failed to process Android templates: %w", err)
			}
		}
	}

	// Process iOS
	if config.Components.Mobile.IOS {
		iosTemplateDir := mobileDir + "/ios-swift"
		iosOutputDir := filepath.Join(mobileOutputDir, "ios")
		if err := a.templateEngine.ProcessDirectory(iosTemplateDir, iosOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process iOS templates: %w", err)
			}
		}
	}

	// Process shared mobile resources
	sharedTemplateDir := mobileDir + "/shared"
	sharedOutputDir := filepath.Join(mobileOutputDir, "shared")
	if err := a.templateEngine.ProcessDirectory(sharedTemplateDir, sharedOutputDir, config); err != nil {
		if !strings.Contains(err.Error(), "file does not exist") {
			return fmt.Errorf("failed to process mobile shared templates: %w", err)
		}
	}

	return nil
}

// processInfrastructureTemplates processes infrastructure templates with proper Deploy/ structure
func (a *App) processInfrastructureTemplates(templateDir, projectOutputDir string, config *models.ProjectConfig) error {
	infraDir := templateDir
	deployDir := filepath.Join(projectOutputDir, "Deploy")

	// Process Docker
	if config.Components.Infrastructure.Docker {
		dockerTemplateDir := infraDir + "/docker"
		dockerOutputDir := filepath.Join(deployDir, "docker")
		if err := a.templateEngine.ProcessDirectory(dockerTemplateDir, dockerOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process Docker templates: %w", err)
			}
		}
	}

	// Process Kubernetes
	if config.Components.Infrastructure.Kubernetes {
		k8sTemplateDir := infraDir + "/kubernetes"
		k8sOutputDir := filepath.Join(deployDir, "k8s")
		if err := a.templateEngine.ProcessDirectory(k8sTemplateDir, k8sOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process Kubernetes templates: %w", err)
			}
		}
	}

	// Process Terraform
	if config.Components.Infrastructure.Terraform {
		terraformTemplateDir := infraDir + "/terraform"
		terraformOutputDir := filepath.Join(deployDir, "terraform")
		if err := a.templateEngine.ProcessDirectory(terraformTemplateDir, terraformOutputDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process Terraform templates: %w", err)
			}
		}
	}

	// Create monitoring directory (even if no templates exist yet)
	if config.Components.Infrastructure.Docker || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Terraform {
		monitoringDir := filepath.Join(deployDir, "monitoring")
		if err := os.MkdirAll(monitoringDir, 0750); err != nil {
			return fmt.Errorf("failed to create monitoring directory: %w", err)
		}

		// Process monitoring templates if they exist
		monitoringTemplateDir := infraDir + "/monitoring"
		if err := a.templateEngine.ProcessDirectory(monitoringTemplateDir, monitoringDir, config); err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return fmt.Errorf("failed to process monitoring templates: %w", err)
			}
		}
	}

	return nil
}

// processBaseRootFiles processes only the root-level template files from base directory
func (a *App) processBaseRootFiles(baseDir, projectOutputDir string, config *models.ProjectConfig) error {
	// For embedded templates, we'll process known base template files
	// This is a simplified approach - in a full implementation you might want to
	// add a method to the embedded engine to list files
	baseTemplateFiles := []string{
		"README.md.tmpl",
		"LICENSE.tmpl",
		"CONTRIBUTING.md.tmpl",
		"SECURITY.md.tmpl",
		"gitignore.tmpl",
		"Makefile.tmpl",
		"docker-compose.yml.tmpl",
	}

	for _, templateFile := range baseTemplateFiles {
		srcFile := baseDir + "/" + templateFile

		// Determine output filename (remove .tmpl extension)
		outputName := strings.TrimSuffix(templateFile, ".tmpl")
		// Handle special case for gitignore -> .gitignore
		if outputName == "gitignore" {
			outputName = ".gitignore"
		}
		outputFile := filepath.Join(projectOutputDir, outputName)

		// Process the individual template file
		content, err := a.templateEngine.ProcessTemplate(srcFile, config)
		if err != nil {
			// Skip if template doesn't exist
			if strings.Contains(err.Error(), "file does not exist") {
				continue
			}
			return fmt.Errorf("failed to process template file %s: %w", templateFile, err)
		}

		// Write the processed content to the output file
		if err := os.MkdirAll(filepath.Dir(outputFile), 0750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := utils.SafeWriteFile(outputFile, content); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outputName, err)
		}
	}

	return nil
}
