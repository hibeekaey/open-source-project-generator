// Package handlers provides workflow management for CLI operations.
//
// This module contains the workflow handler which manages the complete
// project generation workflow including validation, execution, and post-processing.
package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// WorkflowHandler manages the complete project generation workflow.
//
// The WorkflowHandler provides methods for:
//   - Complete generation workflow orchestration
//   - Pre and post-generation tasks
//   - Configuration loading and validation
//   - Output path determination and management
type WorkflowHandler struct {
	cli             CLIInterface
	generateHandler *GenerateHandler
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	logger          interfaces.Logger
}

// NewWorkflowHandler creates a new workflow handler instance.
func NewWorkflowHandler(
	cli CLIInterface,
	generateHandler *GenerateHandler,
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	logger interfaces.Logger,
) *WorkflowHandler {
	return &WorkflowHandler{
		cli:             cli,
		generateHandler: generateHandler,
		configManager:   configManager,
		validator:       validator,
		logger:          logger,
	}
}

// ExecuteGenerationWorkflow executes the project generation workflow
func (wh *WorkflowHandler) ExecuteGenerationWorkflow(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	wh.cli.VerboseOutput("🚀 Starting project generation for: %s", config.Name)

	// Validate configuration if not skipped
	if !options.SkipValidation {
		wh.cli.VerboseOutput("🔍 Validating project configuration...")
		if err := wh.validateGenerateConfiguration(config, options); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
		wh.cli.VerboseOutput("✅ Configuration validation passed")
	}

	// Set output path from options or config
	outputPath := wh.determineOutputPath(config, options)
	wh.cli.VerboseOutput("📁 Output directory: %s", outputPath)

	// Handle offline mode
	if options.Offline {
		wh.cli.VerboseOutput("📡 Running in offline mode - using cached templates and versions")
		// Note: Cache manager access would need to be added to interface
	}

	// Handle version updates
	if options.UpdateVersions && !options.Offline {
		wh.cli.VerboseOutput("📦 Fetching latest package versions...")
		if err := wh.generateHandler.UpdatePackageVersions(config); err != nil {
			wh.cli.WarningOutput("⚠️  Couldn't update package versions: %v", err)
		}
	}

	// Log CI environment information if detected
	ci := wh.detectCIEnvironment()
	if ci.IsCI {
		wh.cli.VerboseOutput("🤖 Detected CI environment: %s", ci.Provider)
		if ci.BuildID != "" {
			wh.cli.VerboseOutput("   Build ID: %s", ci.BuildID)
		}
		if ci.Branch != "" {
			wh.cli.VerboseOutput("   Branch: %s", ci.Branch)
		}
	}

	// Pre-generation checks
	if err := wh.performPreGenerationChecks(outputPath, options); err != nil {
		return fmt.Errorf("pre-generation checks failed: %w", err)
	}

	// Handle dry run mode
	if options.DryRun {
		wh.cli.QuietOutput("🔍 %s - would generate project %s in directory %s",
			wh.cli.Warning("Dry run mode"),
			wh.cli.Highlight(fmt.Sprintf("'%s'", config.Name)),
			wh.cli.Info(fmt.Sprintf("'%s'", outputPath)))
		wh.displayProjectSummary(config)
		return nil
	}

	// Execute the actual project generation
	wh.cli.VerboseOutput("🏗️  Generating project structure...")
	if err := wh.generateHandler.GenerateProjectFromComponents(config, outputPath, options); err != nil {
		return fmt.Errorf("project generation failed: %w", err)
	}

	// Post-generation tasks
	if err := wh.performPostGenerationTasks(config, outputPath, options); err != nil {
		wh.cli.WarningOutput("⚠️  Some post-generation tasks failed: %v", err)
	}

	wh.cli.SuccessOutput("🎉 Project %s %s!", wh.cli.Highlight(fmt.Sprintf("'%s'", config.Name)), wh.cli.Success("generated successfully"))
	wh.displayGenerationSummary(config, outputPath)

	return nil
}

// LoadConfigFromFile loads configuration from a file
func (wh *WorkflowHandler) LoadConfigFromFile(configPath string) (*models.ProjectConfig, error) {
	if wh.configManager == nil {
		return nil, fmt.Errorf("configuration manager not initialized")
	}

	config, err := wh.configManager.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from file: %w", err)
	}

	wh.cli.VerboseOutput("✅ Configuration loaded from file: %s", configPath)
	return config, nil
}

// LoadConfigFromEnvironment loads configuration from environment variables
func (wh *WorkflowHandler) LoadConfigFromEnvironment() (*models.ProjectConfig, error) {
	wh.cli.VerboseOutput("🌍 Loading configuration from environment variables...")

	// Load environment configuration
	envConfig, err := wh.loadEnvironmentConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment configuration: %w", err)
	}

	// Convert environment config to project config
	config, err := wh.convertEnvironmentConfigToProjectConfig(envConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert environment configuration: %w", err)
	}

	// Validate required fields for non-interactive mode
	if config.Name == "" {
		return nil, fmt.Errorf("project name is required for non-interactive mode (set GENERATOR_PROJECT_NAME)")
	}

	wh.cli.VerboseOutput("✅ Configuration loaded from environment variables")
	return config, nil
}

// determineOutputPath determines the output path based on config and options
func (wh *WorkflowHandler) determineOutputPath(config *models.ProjectConfig, options interfaces.GenerateOptions) string {
	outputPath := options.OutputPath
	if outputPath == "" {
		outputPath = config.OutputPath
	}
	if outputPath == "" {
		outputPath = "./output/generated"
	}

	// Always append project name to the output path
	return filepath.Join(outputPath, config.Name)
}

// displayProjectSummary displays a summary of what would be generated
func (wh *WorkflowHandler) displayProjectSummary(config *models.ProjectConfig) {
	wh.cli.QuietOutput("\n%s", wh.cli.Highlight("📋 Project Summary:"))
	wh.cli.QuietOutput("%s", wh.cli.Dim("=================="))
	wh.cli.QuietOutput("Name: %s", wh.cli.Success(config.Name))
	if config.Organization != "" {
		wh.cli.QuietOutput("Organization: %s", wh.cli.Info(config.Organization))
	}
	if config.Description != "" {
		wh.cli.QuietOutput("Description: %s", wh.cli.Dim(config.Description))
	}
	wh.cli.QuietOutput("License: %s", wh.cli.Info(config.License))

	wh.cli.QuietOutput("\n%s", wh.cli.Highlight("🧩 Components:"))
	if wh.generateHandler.hasFrontendComponents(config) {
		wh.cli.QuietOutput("  %s %s", wh.cli.Success("✅"), wh.cli.Info("Frontend (Next.js)"))
	}
	if wh.generateHandler.hasBackendComponents(config) {
		wh.cli.QuietOutput("  %s %s", wh.cli.Success("✅"), wh.cli.Info("Backend (Go Gin)"))
	}
	if wh.generateHandler.hasMobileComponents(config) {
		wh.cli.QuietOutput("  %s %s", wh.cli.Success("✅"), wh.cli.Info("Mobile"))
	}
	if wh.generateHandler.hasInfrastructureComponents(config) {
		wh.cli.QuietOutput("  %s %s", wh.cli.Success("✅"), wh.cli.Info("Infrastructure"))
	}
}

// displayGenerationSummary displays a summary after successful generation
func (wh *WorkflowHandler) displayGenerationSummary(config *models.ProjectConfig, outputPath string) {
	wh.cli.QuietOutput("\n%s", wh.cli.Highlight("📊 Generation Summary:"))
	wh.cli.QuietOutput("%s", wh.cli.Dim("====================="))
	wh.cli.QuietOutput("Project: %s", wh.cli.Success(config.Name))
	wh.cli.QuietOutput("Location: %s", wh.cli.Info(outputPath))
	wh.cli.QuietOutput("Components generated: %s", wh.cli.Success(fmt.Sprintf("%d", wh.countSelectedComponents(config))))

	wh.cli.QuietOutput("\n%s", wh.cli.Highlight("🚀 Next Steps:"))
	wh.cli.QuietOutput("%s. Navigate to your project: %s", wh.cli.Info("1"), wh.cli.Highlight(fmt.Sprintf("cd %s", outputPath)))
	wh.cli.QuietOutput("%s. Review the generated %s for setup instructions", wh.cli.Info("2"), wh.cli.Highlight("README.md"))
	wh.cli.QuietOutput("%s. Install dependencies and start development", wh.cli.Info("3"))
}

// countSelectedComponents counts the number of selected components
func (wh *WorkflowHandler) countSelectedComponents(config *models.ProjectConfig) int {
	count := 0
	if wh.generateHandler.hasFrontendComponents(config) {
		count++
	}
	if wh.generateHandler.hasBackendComponents(config) {
		count++
	}
	if wh.generateHandler.hasMobileComponents(config) {
		count++
	}
	if wh.generateHandler.hasInfrastructureComponents(config) {
		count++
	}
	return count
}

// performPostGenerationTasks performs tasks after project generation
func (wh *WorkflowHandler) performPostGenerationTasks(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	wh.cli.VerboseOutput("🔧 Running post-generation tasks...")

	// Initialize git repository if not in minimal mode
	if !options.Minimal {
		if err := wh.initializeGitRepository(outputPath); err != nil {
			wh.cli.VerboseOutput("⚠️  Git initialization skipped: %v", err)
		}
	}

	// Set file permissions
	if err := wh.setFilePermissions(outputPath); err != nil {
		wh.cli.VerboseOutput("⚠️  File permission setup skipped: %v", err)
	}

	return nil
}

// initializeGitRepository initializes a git repository in the output directory
func (wh *WorkflowHandler) initializeGitRepository(outputPath string) error {
	wh.cli.VerboseOutput("📝 Initializing git repository...")

	cmd := exec.Command("git", "init")
	cmd.Dir = outputPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	wh.cli.VerboseOutput("✅ Git repository initialized")
	return nil
}

// setFilePermissions sets appropriate file permissions for generated files
func (wh *WorkflowHandler) setFilePermissions(outputPath string) error {
	wh.cli.VerboseOutput("🔒 Setting file permissions...")

	// Make script files executable
	scriptsDir := filepath.Join(outputPath, "Scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		err := filepath.Walk(scriptsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(path, ".sh") || strings.HasSuffix(path, ".py")) {
				return os.Chmod(path, 0600)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to set script permissions: %w", err)
		}
	}

	wh.cli.VerboseOutput("✅ File permissions set")
	return nil
}

// CI environment detection
type CIEnvironment struct {
	IsCI     bool
	Provider string
	BuildID  string
	Branch   string
}

// detectCIEnvironment detects if running in a CI environment
func (wh *WorkflowHandler) detectCIEnvironment() CIEnvironment {
	ci := CIEnvironment{}

	// Check common CI environment variables
	if os.Getenv("CI") == "true" || os.Getenv("CONTINUOUS_INTEGRATION") == "true" {
		ci.IsCI = true
	}

	// Detect specific CI providers
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ci.Provider = "GitHub Actions"
		ci.BuildID = os.Getenv("GITHUB_RUN_ID")
		ci.Branch = os.Getenv("GITHUB_REF_NAME")
	} else if os.Getenv("GITLAB_CI") == "true" {
		ci.Provider = "GitLab CI"
		ci.BuildID = os.Getenv("CI_PIPELINE_ID")
		ci.Branch = os.Getenv("CI_COMMIT_REF_NAME")
	} else if os.Getenv("JENKINS_URL") != "" {
		ci.Provider = "Jenkins"
		ci.BuildID = os.Getenv("BUILD_ID")
		ci.Branch = os.Getenv("GIT_BRANCH")
	} else if os.Getenv("CIRCLECI") == "true" {
		ci.Provider = "CircleCI"
		ci.BuildID = os.Getenv("CIRCLE_BUILD_NUM")
		ci.Branch = os.Getenv("CIRCLE_BRANCH")
	}

	return ci
}

// Environment configuration loading
type EnvironmentConfig struct {
	ProjectName  string
	Organization string
	Description  string
	License      string
	OutputPath   string
	Components   map[string]bool
}

// loadEnvironmentConfig loads configuration from environment variables
func (wh *WorkflowHandler) loadEnvironmentConfig() (*EnvironmentConfig, error) {
	config := &EnvironmentConfig{
		ProjectName:  os.Getenv("GENERATOR_PROJECT_NAME"),
		Organization: os.Getenv("GENERATOR_ORGANIZATION"),
		Description:  os.Getenv("GENERATOR_DESCRIPTION"),
		License:      os.Getenv("GENERATOR_LICENSE"),
		OutputPath:   os.Getenv("GENERATOR_OUTPUT_PATH"),
		Components:   make(map[string]bool),
	}

	// Load component selections
	components := []string{
		"GENERATOR_FRONTEND_NEXTJS_APP",
		"GENERATOR_FRONTEND_NEXTJS_HOME",
		"GENERATOR_FRONTEND_NEXTJS_ADMIN",
		"GENERATOR_FRONTEND_NEXTJS_SHARED",
		"GENERATOR_BACKEND_GO_GIN",
		"GENERATOR_MOBILE_ANDROID",
		"GENERATOR_MOBILE_IOS",
		"GENERATOR_INFRASTRUCTURE_DOCKER",
		"GENERATOR_INFRASTRUCTURE_KUBERNETES",
		"GENERATOR_INFRASTRUCTURE_TERRAFORM",
	}

	for _, component := range components {
		value := os.Getenv(component)
		config.Components[component] = value == "true" || value == "1"
	}

	return config, nil
}

// convertEnvironmentConfigToProjectConfig converts environment config to project config
func (wh *WorkflowHandler) convertEnvironmentConfigToProjectConfig(envConfig *EnvironmentConfig) (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{
		Name:         envConfig.ProjectName,
		Organization: envConfig.Organization,
		Description:  envConfig.Description,
		License:      envConfig.License,
		OutputPath:   envConfig.OutputPath,
	}

	// Set default license if not specified
	if config.License == "" {
		config.License = "MIT"
	}

	// Convert component selections
	config.Components.Frontend.NextJS.App = envConfig.Components["GENERATOR_FRONTEND_NEXTJS_APP"]
	config.Components.Frontend.NextJS.Home = envConfig.Components["GENERATOR_FRONTEND_NEXTJS_HOME"]
	config.Components.Frontend.NextJS.Admin = envConfig.Components["GENERATOR_FRONTEND_NEXTJS_ADMIN"]
	config.Components.Frontend.NextJS.Shared = envConfig.Components["GENERATOR_FRONTEND_NEXTJS_SHARED"]
	config.Components.Backend.GoGin = envConfig.Components["GENERATOR_BACKEND_GO_GIN"]
	config.Components.Mobile.Android = envConfig.Components["GENERATOR_MOBILE_ANDROID"]
	config.Components.Mobile.IOS = envConfig.Components["GENERATOR_MOBILE_IOS"]
	config.Components.Infrastructure.Docker = envConfig.Components["GENERATOR_INFRASTRUCTURE_DOCKER"]
	config.Components.Infrastructure.Kubernetes = envConfig.Components["GENERATOR_INFRASTRUCTURE_KUBERNETES"]
	config.Components.Infrastructure.Terraform = envConfig.Components["GENERATOR_INFRASTRUCTURE_TERRAFORM"]

	return config, nil
}

// Validation methods (these would typically delegate to a validator)
func (wh *WorkflowHandler) validateGenerateConfiguration(config *models.ProjectConfig, options interfaces.GenerateOptions) error {
	// Basic validation - in a real implementation this would use the validator
	if config.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

func (wh *WorkflowHandler) performPreGenerationChecks(outputPath string, options interfaces.GenerateOptions) error {
	// Check if output directory exists and handle accordingly
	if _, err := os.Stat(outputPath); err == nil {
		if !options.Force {
			return fmt.Errorf("output directory already exists: %s (use --force to overwrite)", outputPath)
		}
		wh.cli.WarningOutput("⚠️  Output directory exists and will be overwritten")
	}

	// Create parent directory if it doesn't exist
	parentDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(parentDir, 0750); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	return nil
}
