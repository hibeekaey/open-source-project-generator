// Package interactive provides component selection UI for the CLI interface.
//
// This module contains the component selection UI which handles collecting
// component choices from user input through interactive prompts.
package interactive

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ComponentSelection handles interactive component selection.
//
// The ComponentSelection provides methods for:
//   - Collecting component selections from user
//   - Displaying available components with descriptions
//   - Validating component combinations
//   - Providing component recommendations
type ComponentSelection struct {
	logger interfaces.Logger
	output OutputInterface
}

// NewComponentSelection creates a new component selection instance.
func NewComponentSelection(logger interfaces.Logger, output OutputInterface) *ComponentSelection {
	return &ComponentSelection{
		logger: logger,
		output: output,
	}
}

// CollectComponentSelections collects component selections from user input.
func (cs *ComponentSelection) CollectComponentSelections(config *models.ProjectConfig) error {
	cs.output.QuietOutput("")
	cs.output.QuietOutput("🧩 %s", cs.output.Highlight("Component Selection"))
	cs.output.QuietOutput("%s", cs.output.Dim("Choose the components you want to include in your project"))
	cs.output.QuietOutput("")

	// Frontend components
	if err := cs.selectFrontendComponents(config); err != nil {
		return fmt.Errorf("failed to select frontend components: %w", err)
	}

	// Backend components
	if err := cs.selectBackendComponents(config); err != nil {
		return fmt.Errorf("failed to select backend components: %w", err)
	}

	// Mobile components
	if err := cs.selectMobileComponents(config); err != nil {
		return fmt.Errorf("failed to select mobile components: %w", err)
	}

	// Infrastructure components
	if err := cs.selectInfrastructureComponents(config); err != nil {
		return fmt.Errorf("failed to select infrastructure components: %w", err)
	}

	// Validate selections
	if err := cs.validateComponentSelections(config); err != nil {
		return fmt.Errorf("component validation failed: %w", err)
	}

	cs.output.QuietOutput("")
	cs.output.QuietOutput("✅ %s", cs.output.Success("Component selection completed!"))

	return nil
}

// selectFrontendComponents handles frontend component selection
func (cs *ComponentSelection) selectFrontendComponents(config *models.ProjectConfig) error {
	cs.output.QuietOutput("🎨 %s", cs.output.Info("Frontend Components (Next.js)"))
	cs.output.QuietOutput("")

	components := []struct {
		name        string
		description string
		field       *bool
	}{
		{"Main Application", "Full-featured Next.js application with routing and components", &config.Components.Frontend.NextJS.App},
		{"Landing Page", "Marketing/landing page with modern design", &config.Components.Frontend.NextJS.Home},
		{"Admin Dashboard", "Administrative interface with authentication", &config.Components.Frontend.NextJS.Admin},
		{"Shared Components", "Reusable component library for consistency", &config.Components.Frontend.NextJS.Shared},
	}

	for i, component := range components {
		selected, err := cs.promptYesNo(fmt.Sprintf("%d. %s", i+1, component.name), component.description, false)
		if err != nil {
			return err
		}
		*component.field = selected
	}

	return nil
}

// selectBackendComponents handles backend component selection
func (cs *ComponentSelection) selectBackendComponents(config *models.ProjectConfig) error {
	cs.output.QuietOutput("")
	cs.output.QuietOutput("⚙️  %s", cs.output.Info("Backend Components"))
	cs.output.QuietOutput("")

	selected, err := cs.promptYesNo("1. Go Gin API Server", "RESTful API server with Gin framework, middleware, and database integration", true)
	if err != nil {
		return err
	}
	config.Components.Backend.GoGin = selected

	return nil
}

// selectMobileComponents handles mobile component selection
func (cs *ComponentSelection) selectMobileComponents(config *models.ProjectConfig) error {
	cs.output.QuietOutput("")
	cs.output.QuietOutput("📱 %s", cs.output.Info("Mobile Components"))
	cs.output.QuietOutput("")

	// Android
	selected, err := cs.promptYesNo("1. Android Application", "Native Android app with Kotlin and modern architecture", false)
	if err != nil {
		return err
	}
	config.Components.Mobile.Android = selected

	// iOS
	selected, err = cs.promptYesNo("2. iOS Application", "Native iOS app with Swift and modern architecture", false)
	if err != nil {
		return err
	}
	config.Components.Mobile.IOS = selected

	return nil
}

// selectInfrastructureComponents handles infrastructure component selection
func (cs *ComponentSelection) selectInfrastructureComponents(config *models.ProjectConfig) error {
	cs.output.QuietOutput("")
	cs.output.QuietOutput("🚀 %s", cs.output.Info("Infrastructure Components"))
	cs.output.QuietOutput("")

	components := []struct {
		name        string
		description string
		field       *bool
	}{
		{"Docker", "Containerization with Docker and Docker Compose", &config.Components.Infrastructure.Docker},
		{"Kubernetes", "Kubernetes deployment manifests and configurations", &config.Components.Infrastructure.Kubernetes},
		{"Terraform", "Infrastructure as Code with Terraform modules", &config.Components.Infrastructure.Terraform},
	}

	for i, component := range components {
		selected, err := cs.promptYesNo(fmt.Sprintf("%d. %s", i+1, component.name), component.description, false)
		if err != nil {
			return err
		}
		*component.field = selected
	}

	return nil
}

// promptYesNo prompts for a yes/no answer with a description
func (cs *ComponentSelection) promptYesNo(prompt, description string, defaultValue bool) (bool, error) {

	cs.output.QuietOutput("   %s", cs.output.Info(prompt))
	cs.output.QuietOutput("   %s", cs.output.Dim(description))

	for {
		if defaultValue {
			fmt.Printf("   Include? (Y/n): ")
		} else {
			fmt.Printf("   Include? (y/N): ")
		}

		response, err := cs.readInput()
		if err != nil {
			return false, fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" {
			return defaultValue, nil
		}

		if response == "y" || response == "yes" {
			return true, nil
		}

		if response == "n" || response == "no" {
			return false, nil
		}

		cs.output.WarningOutput("   ⚠️  Please enter 'y' for yes or 'n' for no")
	}
}

// readInput reads a line of input from stdin (same as in project_setup.go)
func (cs *ComponentSelection) readInput() (string, error) {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return input, nil
}

// validateComponentSelections validates the selected components
func (cs *ComponentSelection) validateComponentSelections(config *models.ProjectConfig) error {
	hasAnyComponent := cs.hasAnyComponentSelected(config)

	if !hasAnyComponent {
		cs.output.WarningOutput("⚠️  %s", cs.output.Warning("No components selected!"))
		cs.output.QuietOutput("   %s", cs.output.Dim("You must select at least one component to generate a project"))

		// Offer to select default components
		useDefaults, err := cs.promptYesNo("Use default components", "Include main Next.js app and Go Gin API server", true)
		if err != nil {
			return err
		}

		if useDefaults {
			config.Components.Frontend.NextJS.App = true
			config.Components.Backend.GoGin = true
			cs.output.QuietOutput("✅ %s", cs.output.Success("Default components selected"))
		} else {
			return fmt.Errorf("at least one component must be selected")
		}
	}

	// Validate component combinations
	if err := cs.validateComponentCombinations(config); err != nil {
		return err
	}

	return nil
}

// hasAnyComponentSelected checks if any components are selected
func (cs *ComponentSelection) hasAnyComponentSelected(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared ||
		config.Components.Backend.GoGin ||
		config.Components.Mobile.Android ||
		config.Components.Mobile.IOS ||
		config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

// validateComponentCombinations validates component combinations
func (cs *ComponentSelection) validateComponentCombinations(config *models.ProjectConfig) error {
	// Warn if shared components are selected without other frontend components
	if config.Components.Frontend.NextJS.Shared &&
		!config.Components.Frontend.NextJS.App &&
		!config.Components.Frontend.NextJS.Home &&
		!config.Components.Frontend.NextJS.Admin {
		cs.output.WarningOutput("💡 %s", cs.output.Warning("Shared components are most useful with other frontend components"))
	}

	// Warn if Kubernetes is selected without Docker
	if config.Components.Infrastructure.Kubernetes && !config.Components.Infrastructure.Docker {
		cs.output.WarningOutput("💡 %s", cs.output.Warning("Kubernetes typically works best with Docker containers"))

		includeDocker, err := cs.promptYesNo("Include Docker", "Add Docker support for Kubernetes deployment", true)
		if err != nil {
			return err
		}

		if includeDocker {
			config.Components.Infrastructure.Docker = true
			cs.output.QuietOutput("✅ %s", cs.output.Success("Docker support added"))
		}
	}

	return nil
}

// ShowComponentSummary displays a summary of selected components
func (cs *ComponentSelection) ShowComponentSummary(config *models.ProjectConfig) {
	cs.output.QuietOutput("")
	cs.output.QuietOutput("📋 %s", cs.output.Highlight("Selected Components:"))
	cs.output.QuietOutput("%s", cs.output.Dim("==================="))

	// Frontend components
	if cs.hasFrontendComponents(config) {
		cs.output.QuietOutput("🎨 %s", cs.output.Info("Frontend:"))
		if config.Components.Frontend.NextJS.App {
			cs.output.QuietOutput("   ✅ Main Next.js Application")
		}
		if config.Components.Frontend.NextJS.Home {
			cs.output.QuietOutput("   ✅ Landing Page")
		}
		if config.Components.Frontend.NextJS.Admin {
			cs.output.QuietOutput("   ✅ Admin Dashboard")
		}
		if config.Components.Frontend.NextJS.Shared {
			cs.output.QuietOutput("   ✅ Shared Components")
		}
	}

	// Backend components
	if config.Components.Backend.GoGin {
		cs.output.QuietOutput("⚙️  %s", cs.output.Info("Backend:"))
		cs.output.QuietOutput("   ✅ Go Gin API Server")
	}

	// Mobile components
	if cs.hasMobileComponents(config) {
		cs.output.QuietOutput("📱 %s", cs.output.Info("Mobile:"))
		if config.Components.Mobile.Android {
			cs.output.QuietOutput("   ✅ Android Application")
		}
		if config.Components.Mobile.IOS {
			cs.output.QuietOutput("   ✅ iOS Application")
		}
	}

	// Infrastructure components
	if cs.hasInfrastructureComponents(config) {
		cs.output.QuietOutput("🚀 %s", cs.output.Info("Infrastructure:"))
		if config.Components.Infrastructure.Docker {
			cs.output.QuietOutput("   ✅ Docker")
		}
		if config.Components.Infrastructure.Kubernetes {
			cs.output.QuietOutput("   ✅ Kubernetes")
		}
		if config.Components.Infrastructure.Terraform {
			cs.output.QuietOutput("   ✅ Terraform")
		}
	}
}

// Helper methods to check component categories
func (cs *ComponentSelection) hasFrontendComponents(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared
}

func (cs *ComponentSelection) hasMobileComponents(config *models.ProjectConfig) bool {
	return config.Components.Mobile.Android || config.Components.Mobile.IOS
}

func (cs *ComponentSelection) hasInfrastructureComponents(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}
