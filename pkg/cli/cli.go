package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// CLI implements the CLIInterface for interactive user operations
type CLI struct {
	configManager interfaces.ConfigManager
	validator     interfaces.ValidationEngine
}

// NewCLI creates a new CLI instance
func NewCLI(configManager interfaces.ConfigManager, validator interfaces.ValidationEngine) *CLI {
	return &CLI{
		configManager: configManager,
		validator:     validator,
	}
}

// Run executes the CLI application
func (c *CLI) Run() error {
	// This method is implemented by the cobra command structure
	// The actual CLI execution is handled by the app package
	return nil
}

// PromptProjectDetails collects project configuration from user input
func (c *CLI) PromptProjectDetails() (*models.ProjectConfig, error) {
	fmt.Println("🚀 Welcome to the Open Source Template Generator!")
	fmt.Println("Let's configure your new project...")
	fmt.Println()

	config := &models.ProjectConfig{
		CustomVars:       make(map[string]string),
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}

	// Basic project information
	if err := c.promptBasicInfo(config); err != nil {
		return nil, fmt.Errorf("failed to collect basic project info: %w", err)
	}

	// Component selection
	components, err := c.SelectComponents()
	if err != nil {
		return nil, fmt.Errorf("failed to select components: %w", err)
	}

	if err := c.setSelectedComponents(config, components); err != nil {
		return nil, fmt.Errorf("failed to set selected components: %w", err)
	}

	// Output path
	if err := c.promptOutputPath(config); err != nil {
		return nil, fmt.Errorf("failed to set output path: %w", err)
	}

	// Load latest versions
	versions, err := c.configManager.GetLatestVersions()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not fetch latest versions, using defaults: %v\n", err)
		// Use default versions if fetching fails
		versions = c.getDefaultVersions()
	}
	config.Versions = versions

	return config, nil
}

// promptBasicInfo collects basic project information
func (c *CLI) promptBasicInfo(config *models.ProjectConfig) error {
	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Project name:",
				Help:    "Enter a name for your project (alphanumeric characters only)",
			},
			Validate: survey.Required,
		},
		{
			Name: "organization",
			Prompt: &survey.Input{
				Message: "Organization/Author:",
				Help:    "Your organization or personal name",
			},
			Validate: survey.Required,
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Project description:",
				Help:    "Brief description of what your project does",
			},
			Validate: survey.Required,
		},
		{
			Name: "license",
			Prompt: &survey.Select{
				Message: "Choose a license:",
				Options: []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"},
				Default: "MIT",
			},
		},
		{
			Name: "author",
			Prompt: &survey.Input{
				Message: "Author name (optional):",
				Help:    "Your full name for attribution",
			},
		},
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Author email (optional):",
				Help:    "Your email address for contact",
			},
		},
		{
			Name: "repository",
			Prompt: &survey.Input{
				Message: "Repository URL (optional):",
				Help:    "Git repository URL where the project will be hosted",
			},
		},
	}

	answers := struct {
		Name         string
		Organization string
		Description  string
		License      string
		Author       string
		Email        string
		Repository   string
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return fmt.Errorf("failed to collect basic info: %w", err)
	}

	config.Name = answers.Name
	config.Organization = answers.Organization
	config.Description = answers.Description
	config.License = answers.License
	config.Author = answers.Author
	config.Email = answers.Email
	config.Repository = answers.Repository

	return nil
}

// SelectComponents allows user to choose which components to include
func (c *CLI) SelectComponents() ([]string, error) {
	fmt.Println()
	fmt.Println("📦 Component Selection")
	fmt.Println("Choose which components to include in your project:")

	var selectedComponents []string

	// Frontend components
	frontendComponents := []string{
		"frontend.main_app - Main Next.js application",
		"frontend.home - Landing page application",
		"frontend.admin - Admin dashboard application",
	}

	var selectedFrontend []string
	frontendPrompt := &survey.MultiSelect{
		Message: "Frontend Applications:",
		Options: frontendComponents,
		Help:    "Select frontend applications to include. Main app provides core functionality.",
	}

	if err := survey.AskOne(frontendPrompt, &selectedFrontend); err != nil {
		return nil, fmt.Errorf("failed to select frontend components: %w", err)
	}

	selectedComponents = append(selectedComponents, selectedFrontend...)

	// Backend components
	backendComponents := []string{
		"backend.api - Go API server with Gin framework",
	}

	var selectedBackend []string
	backendPrompt := &survey.MultiSelect{
		Message: "Backend Services:",
		Options: backendComponents,
		Help:    "Select backend services to include. API server provides REST endpoints.",
	}

	if err := survey.AskOne(backendPrompt, &selectedBackend); err != nil {
		return nil, fmt.Errorf("failed to select backend components: %w", err)
	}

	selectedComponents = append(selectedComponents, selectedBackend...)

	// Mobile components
	mobileComponents := []string{
		"mobile.android - Android Kotlin application",
		"mobile.ios - iOS Swift application",
	}

	var selectedMobile []string
	mobilePrompt := &survey.MultiSelect{
		Message: "Mobile Applications:",
		Options: mobileComponents,
		Help:    "Select mobile applications to include. Requires backend API for full functionality.",
	}

	if err := survey.AskOne(mobilePrompt, &selectedMobile); err != nil {
		return nil, fmt.Errorf("failed to select mobile components: %w", err)
	}

	selectedComponents = append(selectedComponents, selectedMobile...)

	// Infrastructure components
	infraComponents := []string{
		"infrastructure.docker - Docker configurations",
		"infrastructure.kubernetes - Kubernetes manifests",
		"infrastructure.terraform - Terraform configurations",
	}

	var selectedInfra []string
	infraPrompt := &survey.MultiSelect{
		Message: "Infrastructure:",
		Options: infraComponents,
		Help:    "Select infrastructure components. Docker is recommended for all projects.",
	}

	if err := survey.AskOne(infraPrompt, &selectedInfra); err != nil {
		return nil, fmt.Errorf("failed to select infrastructure components: %w", err)
	}

	selectedComponents = append(selectedComponents, selectedInfra...)

	// Validate component dependencies
	if err := c.validateComponentDependencies(selectedComponents); err != nil {
		return nil, err
	}

	if len(selectedComponents) == 0 {
		return nil, fmt.Errorf("at least one component must be selected")
	}

	return selectedComponents, nil
}

// setSelectedComponents converts the selected component strings to the config structure
func (c *CLI) setSelectedComponents(config *models.ProjectConfig, selected []string) error {
	components := models.Components{}

	for _, component := range selected {
		parts := strings.Split(component, " - ") // Remove description part
		if len(parts) == 0 {
			continue
		}

		componentPath := parts[0]
		switch componentPath {
		case "frontend.main_app":
			components.Frontend.MainApp = true
		case "frontend.home":
			components.Frontend.Home = true
		case "frontend.admin":
			components.Frontend.Admin = true
		case "backend.api":
			components.Backend.API = true
		case "mobile.android":
			components.Mobile.Android = true
		case "mobile.ios":
			components.Mobile.IOS = true
		case "infrastructure.docker":
			components.Infrastructure.Docker = true
		case "infrastructure.kubernetes":
			components.Infrastructure.Kubernetes = true
		case "infrastructure.terraform":
			components.Infrastructure.Terraform = true
		}
	}

	config.Components = components
	return nil
}

// promptOutputPath asks for the output directory
func (c *CLI) promptOutputPath(config *models.ProjectConfig) error {
	defaultPath := fmt.Sprintf("./%s", config.Name)

	outputPrompt := &survey.Input{
		Message: "Output directory:",
		Default: defaultPath,
		Help:    "Directory where the project will be generated",
	}

	var outputPath string
	if err := survey.AskOne(outputPrompt, &outputPath, survey.WithValidator(survey.Required)); err != nil {
		return fmt.Errorf("failed to get output path: %w", err)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	config.OutputPath = absPath
	return nil
}

// ConfirmGeneration shows configuration preview and asks for confirmation
func (c *CLI) ConfirmGeneration(config *models.ProjectConfig) bool {
	fmt.Println()
	fmt.Println("📋 Configuration Summary")
	fmt.Println("========================")
	fmt.Printf("Project Name: %s\n", config.Name)
	fmt.Printf("Organization: %s\n", config.Organization)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("License: %s\n", config.License)
	if config.Author != "" {
		fmt.Printf("Author: %s\n", config.Author)
	}
	if config.Email != "" {
		fmt.Printf("Email: %s\n", config.Email)
	}
	if config.Repository != "" {
		fmt.Printf("Repository: %s\n", config.Repository)
	}
	fmt.Printf("Output Path: %s\n", config.OutputPath)

	fmt.Println("\nSelected Components:")
	c.printSelectedComponents(config.Components)

	fmt.Println("\nPackage Versions:")
	if config.Versions != nil {
		fmt.Printf("  Node.js: %s\n", config.Versions.Node)
		fmt.Printf("  Go: %s\n", config.Versions.Go)
		if config.Versions.NextJS != "" {
			fmt.Printf("  Next.js: %s\n", config.Versions.NextJS)
		}
		if config.Versions.React != "" {
			fmt.Printf("  React: %s\n", config.Versions.React)
		}
		if config.Versions.Kotlin != "" {
			fmt.Printf("  Kotlin: %s\n", config.Versions.Kotlin)
		}
		if config.Versions.Swift != "" {
			fmt.Printf("  Swift: %s\n", config.Versions.Swift)
		}
	}

	fmt.Println()

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Proceed with project generation?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		fmt.Printf("Error getting confirmation: %v\n", err)
		return false
	}

	return confirm
}

// printSelectedComponents displays the selected components in a readable format
func (c *CLI) printSelectedComponents(components models.Components) {
	if components.Frontend.MainApp || components.Frontend.Home || components.Frontend.Admin {
		fmt.Println("  Frontend:")
		if components.Frontend.MainApp {
			fmt.Println("    ✓ Main Application (Next.js)")
		}
		if components.Frontend.Home {
			fmt.Println("    ✓ Landing Page")
		}
		if components.Frontend.Admin {
			fmt.Println("    ✓ Admin Dashboard")
		}
	}

	if components.Backend.API {
		fmt.Println("  Backend:")
		fmt.Println("    ✓ API Server (Go + Gin)")
	}

	if components.Mobile.Android || components.Mobile.IOS {
		fmt.Println("  Mobile:")
		if components.Mobile.Android {
			fmt.Println("    ✓ Android App (Kotlin)")
		}
		if components.Mobile.IOS {
			fmt.Println("    ✓ iOS App (Swift)")
		}
	}

	if components.Infrastructure.Docker || components.Infrastructure.Kubernetes || components.Infrastructure.Terraform {
		fmt.Println("  Infrastructure:")
		if components.Infrastructure.Docker {
			fmt.Println("    ✓ Docker Configurations")
		}
		if components.Infrastructure.Kubernetes {
			fmt.Println("    ✓ Kubernetes Manifests")
		}
		if components.Infrastructure.Terraform {
			fmt.Println("    ✓ Terraform Configurations")
		}
	}
}

// getDefaultVersions returns default version configuration
func (c *CLI) getDefaultVersions() *models.VersionConfig {
	return &models.VersionConfig{
		Node:      "20.11.0",
		Go:        "1.22.0",
		Kotlin:    "2.0.0",
		Swift:     "5.9.0",
		NextJS:    "15.0.0",
		React:     "18.2.0",
		Packages:  make(map[string]string),
		UpdatedAt: time.Now(),
	}
}

// ShowProgress displays a progress indicator with the given message
func (c *CLI) ShowProgress(message string) {
	fmt.Printf("⏳ %s...\n", message)
}

// ShowSuccess displays a success message
func (c *CLI) ShowSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// ShowError displays an error message
func (c *CLI) ShowError(message string) {
	fmt.Printf("❌ %s\n", message)
}

// ShowWarning displays a warning message
func (c *CLI) ShowWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// validateComponentDependencies checks for component dependencies and provides warnings
func (c *CLI) validateComponentDependencies(selectedComponents []string) error {
	return c.validateComponentDependenciesWithPrompt(selectedComponents, true)
}

// validateComponentDependenciesWithPrompt allows controlling whether to prompt for warnings
func (c *CLI) validateComponentDependenciesWithPrompt(selectedComponents []string, promptForWarnings bool) error {
	componentMap := make(map[string]bool)
	for _, comp := range selectedComponents {
		parts := strings.Split(comp, " - ")
		if len(parts) > 0 {
			componentMap[parts[0]] = true
		}
	}

	var warnings []string
	var errors []string

	// Check mobile app dependencies
	if (componentMap["mobile.android"] || componentMap["mobile.ios"]) && !componentMap["backend.api"] {
		warnings = append(warnings, "Mobile applications work best with a backend API for data management")
	}

	// Check Kubernetes without Docker
	if componentMap["infrastructure.kubernetes"] && !componentMap["infrastructure.docker"] {
		warnings = append(warnings, "Kubernetes deployments typically require Docker containers")
	}

	// Check if no main components are selected
	hasMainComponent := componentMap["frontend.main_app"] || componentMap["backend.api"] ||
		componentMap["mobile.android"] || componentMap["mobile.ios"]

	if !hasMainComponent {
		errors = append(errors, "At least one main component (frontend app, backend API, or mobile app) must be selected")
	}

	// Display warnings
	if len(warnings) > 0 {
		fmt.Println()
		fmt.Println("⚠️  Dependency Warnings:")
		for _, warning := range warnings {
			fmt.Printf("  • %s\n", warning)
		}

		if promptForWarnings {
			var proceed bool
			proceedPrompt := &survey.Confirm{
				Message: "Continue with these warnings?",
				Default: true,
			}

			if err := survey.AskOne(proceedPrompt, &proceed); err != nil {
				return fmt.Errorf("failed to get proceed confirmation: %w", err)
			}

			if !proceed {
				return fmt.Errorf("component selection cancelled by user")
			}
		}
	}

	// Display errors
	if len(errors) > 0 {
		fmt.Println()
		fmt.Println("❌ Dependency Errors:")
		for _, err := range errors {
			fmt.Printf("  • %s\n", err)
		}
		return fmt.Errorf("component dependency validation failed")
	}

	return nil
}

// PreviewConfiguration shows a detailed preview of what will be generated
func (c *CLI) PreviewConfiguration(config *models.ProjectConfig) {
	fmt.Println()
	fmt.Println("🔍 Generation Preview")
	fmt.Println("====================")

	fmt.Println("Directory Structure:")
	c.showDirectoryStructure(config.Components)

	fmt.Println()
	fmt.Println("Key Files to be Generated:")
	c.showKeyFiles(config.Components)

	fmt.Println()
	fmt.Println("Build Commands Available:")
	c.showBuildCommands(config.Components)
}

// showDirectoryStructure displays the directory structure that will be created
func (c *CLI) showDirectoryStructure(components models.Components) {
	fmt.Printf("  %s/\n", "project-root")

	if components.Frontend.MainApp || components.Frontend.Home || components.Frontend.Admin {
		fmt.Println("  ├── App/")
		if components.Frontend.MainApp {
			fmt.Println("  │   ├── main/          # Main Next.js application")
		}
		if components.Frontend.Home {
			fmt.Println("  │   ├── home/          # Landing page")
		}
		if components.Frontend.Admin {
			fmt.Println("  │   └── admin/         # Admin dashboard")
		}
	}

	if components.Backend.API {
		fmt.Println("  ├── CommonServer/      # Go API server")
		fmt.Println("  │   ├── cmd/")
		fmt.Println("  │   ├── internal/")
		fmt.Println("  │   └── pkg/")
	}

	if components.Mobile.Android || components.Mobile.IOS {
		fmt.Println("  ├── Mobile/")
		if components.Mobile.Android {
			fmt.Println("  │   ├── android/       # Kotlin Android app")
		}
		if components.Mobile.IOS {
			fmt.Println("  │   └── ios/           # Swift iOS app")
		}
	}

	if components.Infrastructure.Docker || components.Infrastructure.Kubernetes || components.Infrastructure.Terraform {
		fmt.Println("  ├── Deploy/")
		if components.Infrastructure.Docker {
			fmt.Println("  │   ├── docker/        # Docker configurations")
		}
		if components.Infrastructure.Kubernetes {
			fmt.Println("  │   ├── k8s/           # Kubernetes manifests")
		}
		if components.Infrastructure.Terraform {
			fmt.Println("  │   └── terraform/     # Infrastructure as code")
		}
	}

	fmt.Println("  ├── Docs/              # Documentation")
	fmt.Println("  ├── Scripts/           # Build and deployment scripts")
	fmt.Println("  ├── .github/           # CI/CD workflows")
	fmt.Println("  └── Makefile           # Build system")
}

// showKeyFiles displays key files that will be generated
func (c *CLI) showKeyFiles(components models.Components) {
	files := []string{
		"README.md - Project documentation",
		"CONTRIBUTING.md - Contribution guidelines",
		"LICENSE - Project license",
		"Makefile - Build system",
		"docker-compose.yml - Development environment",
		".gitignore - Git ignore rules",
		".github/workflows/ - CI/CD pipelines",
	}

	if components.Frontend.MainApp || components.Frontend.Home || components.Frontend.Admin {
		files = append(files, "package.json - Node.js dependencies")
		files = append(files, "next.config.js - Next.js configuration")
		files = append(files, "tailwind.config.js - Tailwind CSS setup")
	}

	if components.Backend.API {
		files = append(files, "go.mod - Go module definition")
		files = append(files, "main.go - API server entry point")
	}

	if components.Mobile.Android {
		files = append(files, "build.gradle - Android build configuration")
	}

	if components.Mobile.IOS {
		files = append(files, "Podfile - iOS dependencies")
	}

	for _, file := range files {
		fmt.Printf("  • %s\n", file)
	}
}

// showBuildCommands displays available build commands
func (c *CLI) showBuildCommands(components models.Components) {
	commands := []string{
		"make setup - Initialize development environment",
		"make dev - Start development servers",
		"make test - Run all tests",
		"make build - Build all components",
		"make clean - Clean build artifacts",
	}

	if components.Infrastructure.Docker {
		commands = append(commands, "make docker-build - Build Docker images")
		commands = append(commands, "make docker-up - Start with Docker Compose")
	}

	if components.Infrastructure.Kubernetes {
		commands = append(commands, "make k8s-deploy - Deploy to Kubernetes")
	}

	for _, cmd := range commands {
		fmt.Printf("  • %s\n", cmd)
	}
}

// CheckOutputPath validates that the output path is suitable for generation
func (c *CLI) CheckOutputPath(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); err == nil {
		// Path exists, check if it's empty
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("cannot read directory %s: %w", path, err)
		}

		if len(entries) > 0 {
			var overwrite bool
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Directory %s is not empty. Overwrite?", path),
				Default: false,
			}

			if err := survey.AskOne(prompt, &overwrite); err != nil {
				return fmt.Errorf("failed to get overwrite confirmation: %w", err)
			}

			if !overwrite {
				return fmt.Errorf("generation cancelled by user")
			}
		}
	}

	return nil
}

// ExecuteCommand executes CLI commands based on the provided arguments
func (c *CLI) ExecuteCommand(command string, args []string) error {
	switch command {
	case "analyze-templates":
		return c.AnalyzeTemplatesCommand(args)
	case "save-analysis-report":
		return c.SaveAnalysisReportCommand(args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}
