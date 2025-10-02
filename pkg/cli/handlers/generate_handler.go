// Package handlers provides command handler implementations for the CLI interface.
//
// This module contains the generate command handler which manages project generation
// workflows including interactive and non-interactive modes.
package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
)

// GenerateHandler handles the generate command execution.
//
// The GenerateHandler provides methods for:
//   - Project generation workflow orchestration
//   - Component-based project structure creation
//   - Template processing and file generation
//   - Pre and post-generation tasks
type GenerateHandler struct {
	cli             CLIInterface
	templateManager interfaces.TemplateManager
	configManager   interfaces.ConfigManager
	validator       interfaces.ValidationEngine
	logger          interfaces.Logger
}

// CLIInterface defines the CLI methods needed by the generate handler
type CLIInterface interface {
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	QuietOutput(format string, args ...interface{})
	WarningOutput(format string, args ...interface{})
	SuccessOutput(format string, args ...interface{})
	Error(text string) string
	Info(text string) string
	Warning(text string) string
	Success(text string) string
	Highlight(text string) string
	Dim(text string) string
	GetVersionManager() interfaces.VersionManager
}

// NewGenerateHandler creates a new generate handler instance.
func NewGenerateHandler(
	cli CLIInterface,
	templateManager interfaces.TemplateManager,
	configManager interfaces.ConfigManager,
	validator interfaces.ValidationEngine,
	logger interfaces.Logger,
) *GenerateHandler {
	return &GenerateHandler{
		cli:             cli,
		templateManager: templateManager,
		configManager:   configManager,
		validator:       validator,
		logger:          logger,
	}
}

// GenerateProjectFromComponents generates project structure based on selected components
func (gh *GenerateHandler) GenerateProjectFromComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	gh.cli.VerboseOutput("🏗️  Building your project structure...")

	// Use performance optimization if available
	if optimizedCLI, ok := gh.cli.(interface {
		OptimizeCommand(string, func(ctx context.Context) (interface{}, error)) (interface{}, error)
	}); ok {
		_, err := optimizedCLI.OptimizeCommand("generate", func(ctx context.Context) (interface{}, error) {
			return nil, gh.generateProjectInternal(config, outputPath, options)
		})
		return err
	}

	// Fallback to direct execution
	return gh.generateProjectInternal(config, outputPath, options)
}

// generateProjectInternal performs the actual project generation
func (gh *GenerateHandler) generateProjectInternal(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error {
	// Create the base project structure first
	if err := gh.generateBaseStructure(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to create project structure."),
			gh.cli.Info("Check output directory permissions and available disk space"))
	}

	// Process frontend components
	if err := gh.processFrontendComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to set up frontend components."),
			gh.cli.Info("Check if frontend templates are available and accessible"))
	}

	// Process backend components
	if err := gh.processBackendComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to set up backend components."),
			gh.cli.Info("Check if backend templates are available and accessible"))
	}

	// Process mobile components
	if err := gh.processMobileComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to set up mobile components."),
			gh.cli.Info("Check if mobile templates are available and accessible"))
	}

	// Process infrastructure components
	if err := gh.processInfrastructureComponents(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to set up infrastructure components."),
			gh.cli.Info("Check if infrastructure templates are available and accessible"))
	}

	return nil
}

// generateBaseStructure generates the base project structure
func (gh *GenerateHandler) generateBaseStructure(config *models.ProjectConfig, outputPath string) error {
	gh.cli.VerboseOutput("📋 Creating project foundation...")

	// Create the main project directories first
	dirs := []string{"Docs", "Scripts"}
	for _, dir := range dirs {
		dirPath := filepath.Join(outputPath, dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return fmt.Errorf("🚫 %s %s %s",
				gh.cli.Error("Unable to create directory"),
				gh.cli.Highlight(fmt.Sprintf("'%s'.", dir)),
				gh.cli.Info("Check permissions and available disk space"))
		}
	}

	// Process base template files directly from the embedded filesystem
	if err := gh.processBaseTemplateFiles(config, outputPath); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Failed to process base template files."),
			gh.cli.Info("Essential project files like README and LICENSE couldn't be created"))
	}

	// Process GitHub workflow templates (rename github -> .github)
	if err := gh.templateManager.ProcessTemplate("github", config, filepath.Join(outputPath, ".github")); err != nil {
		gh.cli.VerboseOutput("GitHub template not processed (optional): %v", err)
	}

	// Clean up duplicate github folder
	githubDir := filepath.Join(outputPath, "github")
	if _, err := os.Stat(githubDir); err == nil {
		gh.cli.VerboseOutput("🧹 Cleaning up duplicate folder structure...")
		if err := os.RemoveAll(githubDir); err != nil {
			gh.cli.VerboseOutput("⚠️  Could not clean up temporary files: %v", err)
		} else {
			gh.cli.VerboseOutput("✅ Project structure optimized")
		}
	}

	// Process Scripts templates
	if err := gh.templateManager.ProcessTemplate("scripts", config, filepath.Join(outputPath, "Scripts")); err != nil {
		gh.cli.VerboseOutput("Scripts template not processed (optional): %v", err)
	}

	return nil
}

// processBaseTemplateFiles processes base template files directly
func (gh *GenerateHandler) processBaseTemplateFiles(config *models.ProjectConfig, outputPath string) error {
	// Create an embedded template engine to process the base directory directly
	embeddedEngine := template.NewEmbeddedEngine()

	// Process the base template directory directly
	return embeddedEngine.ProcessDirectory("templates/base", outputPath, config)
}

// processFrontendComponents processes frontend components
func (gh *GenerateHandler) processFrontendComponents(config *models.ProjectConfig, outputPath string) error {
	if !gh.hasFrontendComponents(config) {
		gh.cli.VerboseOutput("No frontend components selected, skipping")
		return nil
	}

	gh.cli.VerboseOutput("🎨 Setting up frontend applications...")

	// Create App directory structure
	appDir := filepath.Join(outputPath, "App")
	if err := os.MkdirAll(appDir, 0750); err != nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Unable to create App directory."),
			gh.cli.Info("Check output directory permissions and available disk space"))
	}

	// Process Next.js components based on configuration
	if config.Components.Frontend.NextJS.App {
		gh.cli.VerboseOutput("   ✨ Creating main Next.js application")
		mainAppPath := filepath.Join(appDir, "main")
		if err := gh.templateManager.ProcessTemplate("nextjs-app", config, mainAppPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Next.js app template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	if config.Components.Frontend.NextJS.Home {
		gh.cli.VerboseOutput("   🏠 Creating landing page application")
		homePath := filepath.Join(appDir, "home")
		if err := gh.templateManager.ProcessTemplate("nextjs-home", config, homePath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Next.js home template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	if config.Components.Frontend.NextJS.Admin {
		gh.cli.VerboseOutput("   👑 Creating admin dashboard")
		adminPath := filepath.Join(appDir, "admin")
		if err := gh.templateManager.ProcessTemplate("nextjs-admin", config, adminPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Next.js admin template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	if config.Components.Frontend.NextJS.Shared {
		gh.cli.VerboseOutput("📦 Creating shared component library...")
		sharedPath := filepath.Join(appDir, "shared-components")
		if err := gh.templateManager.ProcessTemplate("shared-components", config, sharedPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process shared components template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	return nil
}

// processBackendComponents processes backend components
func (gh *GenerateHandler) processBackendComponents(config *models.ProjectConfig, outputPath string) error {
	if !gh.hasBackendComponents(config) {
		gh.cli.VerboseOutput("No backend components selected, skipping")
		return nil
	}

	gh.cli.VerboseOutput("⚙️  Setting up backend services...")

	// Create CommonServer directory
	serverDir := filepath.Join(outputPath, "CommonServer")
	if err := os.MkdirAll(serverDir, 0750); err != nil {
		return fmt.Errorf("failed to create CommonServer directory: %w", err)
	}

	// Process Go Gin backend
	if config.Components.Backend.GoGin {
		gh.cli.VerboseOutput("   🔧 Creating Go API server")
		if err := gh.templateManager.ProcessTemplate("go-gin", config, serverDir); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Go Gin template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	return nil
}

// processMobileComponents processes mobile components
func (gh *GenerateHandler) processMobileComponents(config *models.ProjectConfig, outputPath string) error {
	if !gh.hasMobileComponents(config) {
		gh.cli.VerboseOutput("No mobile components selected, skipping")
		return nil
	}

	gh.cli.VerboseOutput("📱 Setting up mobile applications...")

	// Create Mobile directory
	mobileDir := filepath.Join(outputPath, "Mobile")
	if err := os.MkdirAll(mobileDir, 0750); err != nil {
		return fmt.Errorf("failed to create Mobile directory: %w", err)
	}

	// Process Android components
	if config.Components.Mobile.Android {
		gh.cli.VerboseOutput("   🤖 Creating Android application")
		androidPath := filepath.Join(mobileDir, "android")
		if err := gh.templateManager.ProcessTemplate("android-kotlin", config, androidPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Android Kotlin template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	// Process iOS components
	if config.Components.Mobile.IOS {
		gh.cli.VerboseOutput("   🍎 Creating iOS application")
		iosPath := filepath.Join(mobileDir, "ios")
		if err := gh.templateManager.ProcessTemplate("ios-swift", config, iosPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process iOS Swift template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	// Process shared mobile components
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		gh.cli.VerboseOutput("🔗 Creating shared mobile resources...")
		sharedPath := filepath.Join(mobileDir, "shared")
		if err := gh.templateManager.ProcessTemplate("shared", config, sharedPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process mobile shared template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	return nil
}

// processInfrastructureComponents processes infrastructure components
func (gh *GenerateHandler) processInfrastructureComponents(config *models.ProjectConfig, outputPath string) error {
	if !gh.hasInfrastructureComponents(config) {
		gh.cli.VerboseOutput("No infrastructure components selected, skipping")
		return nil
	}

	gh.cli.VerboseOutput("🚀 Setting up deployment infrastructure...")

	// Create Deploy directory
	deployDir := filepath.Join(outputPath, "Deploy")
	if err := os.MkdirAll(deployDir, 0750); err != nil {
		return fmt.Errorf("failed to create Deploy directory: %w", err)
	}

	// Process Docker components
	if config.Components.Infrastructure.Docker {
		gh.cli.VerboseOutput("   🐳 Setting up Docker containers")
		dockerPath := filepath.Join(deployDir, "docker")
		if err := gh.templateManager.ProcessTemplate("docker", config, dockerPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Docker template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	// Process Kubernetes components
	if config.Components.Infrastructure.Kubernetes {
		gh.cli.VerboseOutput("   ☸️  Setting up Kubernetes deployment")
		k8sPath := filepath.Join(deployDir, "k8s")
		if err := gh.templateManager.ProcessTemplate("kubernetes", config, k8sPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Kubernetes template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	// Process Terraform components
	if config.Components.Infrastructure.Terraform {
		gh.cli.VerboseOutput("   🏗️  Setting up Terraform infrastructure")
		terraformPath := filepath.Join(deployDir, "terraform")
		if err := gh.templateManager.ProcessTemplate("terraform", config, terraformPath); err != nil {
			return fmt.Errorf("🚫 %s %s",
				gh.cli.Error("Failed to process Terraform template."),
				gh.cli.Info("Check if the template files are accessible and valid"))
		}
	}

	return nil
}

// Helper methods to check if components are selected
func (gh *GenerateHandler) hasFrontendComponents(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared
}

func (gh *GenerateHandler) hasBackendComponents(config *models.ProjectConfig) bool {
	return config.Components.Backend.GoGin
}

func (gh *GenerateHandler) hasMobileComponents(config *models.ProjectConfig) bool {
	return config.Components.Mobile.Android || config.Components.Mobile.IOS
}

func (gh *GenerateHandler) hasInfrastructureComponents(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

// UpdatePackageVersions updates package versions in the configuration
func (gh *GenerateHandler) UpdatePackageVersions(config *models.ProjectConfig) error {
	versionManager := gh.cli.GetVersionManager()
	if versionManager == nil {
		return fmt.Errorf("🚫 %s %s",
			gh.cli.Error("Version manager not initialized."),
			gh.cli.Info("This is an internal error - please report this issue"))
	}

	gh.cli.VerboseOutput("Fetching latest package versions...")

	// This would fetch latest versions and update the config
	// For now, we'll just log that we would do this
	gh.cli.VerboseOutput("Would update package versions for project type based on configuration")

	return nil
}
