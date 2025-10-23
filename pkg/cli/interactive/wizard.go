package interactive

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
)

// InteractiveWizard manages the interactive configuration flow
type InteractiveWizard struct {
	prompter  Prompter
	logger    *logger.Logger
	sanitizer *security.Sanitizer
}

// NewInteractiveWizard creates a new interactive wizard
func NewInteractiveWizard(log *logger.Logger) *InteractiveWizard {
	return &InteractiveWizard{
		prompter:  NewCLIPrompter(),
		logger:    log,
		sanitizer: security.NewSanitizer(),
	}
}

// ProjectInfo holds basic project information
type ProjectInfo struct {
	Name        string
	Description string
	OutputDir   string
}

// Run executes the interactive configuration wizard
func (iw *InteractiveWizard) Run(ctx context.Context) (*models.ProjectConfig, error) {
	iw.logger.Info("Starting interactive configuration wizard")

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("  Open Source Project Generator - Interactive Mode")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("\nThis wizard will guide you through creating a project configuration.")
	fmt.Println("Press Ctrl+C at any time to cancel.")

	// Step 1: Collect project information
	projectInfo, err := iw.CollectProjectInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to collect project info: %w", err)
	}

	// Step 2: Select components
	componentTypes, err := iw.SelectComponents()
	if err != nil {
		return nil, fmt.Errorf("failed to select components: %w", err)
	}

	// Step 3: Configure each component
	components := make([]models.ComponentConfig, 0, len(componentTypes))
	for _, componentType := range componentTypes {
		config, err := iw.ConfigureComponent(componentType)
		if err != nil {
			return nil, fmt.Errorf("failed to configure %s component: %w", componentType, err)
		}
		components = append(components, config)
	}

	// Step 4: Configure integration options
	integration, err := iw.ConfigureIntegration()
	if err != nil {
		return nil, fmt.Errorf("failed to configure integration: %w", err)
	}

	// Build the project configuration
	projectConfig := &models.ProjectConfig{
		Name:        projectInfo.Name,
		Description: projectInfo.Description,
		OutputDir:   projectInfo.OutputDir,
		Components:  components,
		Integration: integration,
		Options: models.ProjectOptions{
			UseExternalTools: true,
			CreateBackup:     true,
			Verbose:          false,
			DryRun:           false,
			ForceOverwrite:   false,
		},
	}

	// Step 5: Confirm configuration
	confirmed, err := iw.ConfirmConfiguration(projectConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm configuration: %w", err)
	}

	if !confirmed {
		return nil, fmt.Errorf("configuration cancelled by user")
	}

	iw.logger.Info("Interactive configuration completed successfully")
	return projectConfig, nil
}

// CollectProjectInfo prompts for basic project information
func (iw *InteractiveWizard) CollectProjectInfo() (*ProjectInfo, error) {
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Step 1: Project Information")
	fmt.Println(strings.Repeat("-", 70) + "\n")

	// Project name
	name, err := InputWithValidation(
		iw.prompter,
		"Project name",
		"my-project",
		func(input string) error {
			sanitized, err := iw.sanitizer.SanitizeProjectName(input)
			if err != nil {
				return err
			}
			if sanitized != input {
				return fmt.Errorf("project name will be sanitized to: %s", sanitized)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Sanitize the name
	name, err = iw.sanitizer.SanitizeProjectName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize project name: %w", err)
	}

	// Project description
	description, err := iw.prompter.Input("Project description", "A new project")
	if err != nil {
		return nil, err
	}

	// Sanitize description
	description, err = iw.sanitizer.SanitizeString(description)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize description: %w", err)
	}

	// Output directory
	defaultOutputDir := "./" + name
	outputDir, err := InputWithValidation(
		iw.prompter,
		"Output directory",
		defaultOutputDir,
		func(input string) error {
			_, err := iw.sanitizer.SanitizePath(input)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	// Sanitize path
	outputDir, err = iw.sanitizer.SanitizePath(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize output directory: %w", err)
	}

	// Convert to absolute path
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	fmt.Println()
	return &ProjectInfo{
		Name:        name,
		Description: description,
		OutputDir:   outputDir,
	}, nil
}

// SelectComponents presents component selection menu
func (iw *InteractiveWizard) SelectComponents() ([]string, error) {
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Step 2: Component Selection")
	fmt.Println(strings.Repeat("-", 70) + "\n")

	componentOptions := []string{
		"nextjs - Next.js frontend with TypeScript and Tailwind",
		"go-backend - Go backend with Gin framework",
		"android - Android mobile app with Kotlin",
		"ios - iOS mobile app with Swift",
	}

	selected, err := iw.prompter.MultiSelect(
		"Select components to include in your project:",
		componentOptions,
	)
	if err != nil {
		return nil, err
	}

	// Extract component types from selections
	componentTypes := make([]string, 0, len(selected))
	for _, selection := range selected {
		// Extract the component type (first word before the dash)
		parts := strings.SplitN(selection, " ", 2)
		if len(parts) > 0 {
			componentTypes = append(componentTypes, parts[0])
		}
	}

	fmt.Printf("\n✓ Selected %d component(s)\n\n", len(componentTypes))
	return componentTypes, nil
}

// ConfigureComponent prompts for component-specific configuration
func (iw *InteractiveWizard) ConfigureComponent(componentType string) (models.ComponentConfig, error) {
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Step 3: Configure %s Component\n", componentType)
	fmt.Println(strings.Repeat("-", 70) + "\n")

	config := models.ComponentConfig{
		Type:    componentType,
		Enabled: true,
		Config:  make(map[string]interface{}),
	}

	// Component name
	defaultName := componentType
	name, err := InputWithValidation(
		iw.prompter,
		fmt.Sprintf("Component name for %s", componentType),
		defaultName,
		ValidateProjectName,
	)
	if err != nil {
		return config, err
	}
	config.Name = name
	config.Config["name"] = name

	// Component-specific configuration
	switch componentType {
	case "nextjs":
		if err := iw.configureNextJS(&config); err != nil {
			return config, err
		}
	case "go-backend":
		if err := iw.configureGoBackend(&config); err != nil {
			return config, err
		}
	case "android":
		if err := iw.configureAndroid(&config); err != nil {
			return config, err
		}
	case "ios":
		if err := iw.configureIOS(&config); err != nil {
			return config, err
		}
	default:
		return config, fmt.Errorf("unsupported component type: %s", componentType)
	}

	fmt.Println()
	return config, nil
}

// ConfirmConfiguration displays summary and asks for confirmation
func (iw *InteractiveWizard) ConfirmConfiguration(config *models.ProjectConfig) (bool, error) {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Configuration Summary")
	fmt.Println(strings.Repeat("=", 70) + "\n")

	fmt.Printf("Project Name: %s\n", config.Name)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("Output Directory: %s\n", config.OutputDir)
	fmt.Printf("\nComponents (%d):\n", len(config.Components))

	for i, comp := range config.Components {
		fmt.Printf("  %d. %s (%s)\n", i+1, comp.Name, comp.Type)

		// Display key configuration options
		switch comp.Type {
		case "nextjs":
			fmt.Printf("     - TypeScript: %v\n", comp.Config["typescript"])
			fmt.Printf("     - Tailwind: %v\n", comp.Config["tailwind"])
			fmt.Printf("     - App Router: %v\n", comp.Config["app_router"])
		case "go-backend":
			fmt.Printf("     - Module: %v\n", comp.Config["module"])
			fmt.Printf("     - Framework: %v\n", comp.Config["framework"])
			fmt.Printf("     - Port: %v\n", comp.Config["port"])
		case "android":
			fmt.Printf("     - Package: %v\n", comp.Config["package"])
			fmt.Printf("     - Min SDK: %v\n", comp.Config["min_sdk"])
			fmt.Printf("     - Target SDK: %v\n", comp.Config["target_sdk"])
		case "ios":
			fmt.Printf("     - Bundle ID: %v\n", comp.Config["bundle_id"])
			fmt.Printf("     - Deployment Target: %v\n", comp.Config["deployment_target"])
		}
	}

	fmt.Printf("\nIntegration:\n")
	fmt.Printf("  - Docker Compose: %v\n", config.Integration.GenerateDockerCompose)
	fmt.Printf("  - Build Scripts: %v\n", config.Integration.GenerateScripts)

	fmt.Println("\n" + strings.Repeat("=", 70))

	confirmed, err := iw.prompter.Confirm("Proceed with this configuration?", true)
	if err != nil {
		return false, err
	}

	return confirmed, nil
}

// ConfigureIntegration prompts for integration options
func (iw *InteractiveWizard) ConfigureIntegration() (models.IntegrationConfig, error) {
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Step 4: Integration Options")
	fmt.Println(strings.Repeat("-", 70) + "\n")

	integration := models.IntegrationConfig{
		APIEndpoints:      make(map[string]string),
		SharedEnvironment: make(map[string]string),
	}

	// Docker Compose
	generateDocker, err := iw.prompter.Confirm("Generate Docker Compose configuration?", true)
	if err != nil {
		return integration, err
	}
	integration.GenerateDockerCompose = generateDocker

	// Build scripts
	generateScripts, err := iw.prompter.Confirm("Generate build and deployment scripts?", true)
	if err != nil {
		return integration, err
	}
	integration.GenerateScripts = generateScripts

	fmt.Println()
	return integration, nil
}
