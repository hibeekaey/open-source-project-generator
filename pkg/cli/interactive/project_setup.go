// Package interactive provides interactive UI components for the CLI interface.
//
// This module contains the project setup UI which handles collecting
// project configuration from user input through interactive prompts.
package interactive

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// Pre-compiled regular expressions for interactive validation
var (
	interactiveProjectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$`)
)

// ProjectSetup handles interactive project configuration collection.
//
// The ProjectSetup provides methods for:
//   - Collecting basic project information
//   - Validating user input
//   - Providing helpful prompts and suggestions
//   - Handling user input errors gracefully
type ProjectSetup struct {
	logger interfaces.Logger
	output OutputInterface
}

// OutputInterface defines the output methods needed by project setup
type OutputInterface interface {
	QuietOutput(format string, args ...interface{})
	VerboseOutput(format string, args ...interface{})
	WarningOutput(format string, args ...interface{})
	Error(text string) string
	Info(text string) string
	Warning(text string) string
	Success(text string) string
	Highlight(text string) string
	Dim(text string) string
}

// NewProjectSetup creates a new project setup instance.
func NewProjectSetup(logger interfaces.Logger, output OutputInterface) *ProjectSetup {
	return &ProjectSetup{
		logger: logger,
		output: output,
	}
}

// CollectProjectDetails collects basic project configuration from user input.
func (ps *ProjectSetup) CollectProjectDetails() (*models.ProjectConfig, error) {
	ps.output.QuietOutput("🎯 %s", ps.output.Highlight("Let's set up your new project!"))
	ps.output.QuietOutput("")

	config := &models.ProjectConfig{}

	// Collect project name
	name, err := ps.promptProjectName()
	if err != nil {
		return nil, fmt.Errorf("failed to get project name: %w", err)
	}
	config.Name = name

	// Collect organization (optional)
	organization, err := ps.promptOrganization()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	config.Organization = organization

	// Collect description (optional)
	description, err := ps.promptDescription()
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}
	config.Description = description

	// Collect license
	license, err := ps.promptLicense()
	if err != nil {
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	config.License = license

	ps.output.QuietOutput("")
	ps.output.QuietOutput("✅ %s", ps.output.Success("Basic project information collected!"))

	return config, nil
}

// promptProjectName prompts for and validates the project name
func (ps *ProjectSetup) promptProjectName() (string, error) {
	for {
		ps.output.QuietOutput("📝 %s", ps.output.Info("Project name (required):"))
		ps.output.QuietOutput("   %s", ps.output.Dim("This will be used as the directory name and in configuration files"))
		fmt.Print("   > ")

		name, err := ps.readInput()
		if err != nil {
			return "", err
		}

		name = strings.TrimSpace(name)
		if name == "" {
			ps.output.WarningOutput("⚠️  Project name cannot be empty. Please try again.")
			continue
		}

		if !ps.isValidProjectName(name) {
			ps.output.WarningOutput("⚠️  %s", ps.output.Warning("Invalid project name."))
			ps.output.QuietOutput("   %s", ps.output.Dim("Project names should contain only letters, numbers, hyphens, and underscores"))
			ps.output.QuietOutput("   %s", ps.output.Dim("Examples: my-project, awesome_app, project123"))
			continue
		}

		return name, nil
	}
}

// promptOrganization prompts for the organization name (optional)
func (ps *ProjectSetup) promptOrganization() (string, error) {
	ps.output.QuietOutput("🏢 %s", ps.output.Info("Organization name (optional):"))
	ps.output.QuietOutput("   %s", ps.output.Dim("Your company or organization name (press Enter to skip)"))
	fmt.Print("   > ")

	organization, err := ps.readInput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(organization), nil
}

// promptDescription prompts for the project description (optional)
func (ps *ProjectSetup) promptDescription() (string, error) {
	ps.output.QuietOutput("📄 %s", ps.output.Info("Project description (optional):"))
	ps.output.QuietOutput("   %s", ps.output.Dim("A brief description of your project (press Enter to skip)"))
	fmt.Print("   > ")

	description, err := ps.readInput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(description), nil
}

// promptLicense prompts for and validates the license
func (ps *ProjectSetup) promptLicense() (string, error) {
	ps.output.QuietOutput("⚖️  %s", ps.output.Info("License (default: MIT):"))
	ps.output.QuietOutput("   %s", ps.output.Dim("Common options: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause"))
	ps.output.QuietOutput("   %s", ps.output.Dim("Press Enter for MIT license"))
	fmt.Print("   > ")

	license, err := ps.readInput()
	if err != nil {
		return "", err
	}

	license = strings.TrimSpace(license)
	if license == "" {
		license = "MIT"
	}

	if !ps.isValidLicense(license) {
		ps.output.WarningOutput("⚠️  %s", ps.output.Warning("Unknown license identifier."))
		ps.output.QuietOutput("   %s", ps.output.Dim("Using the provided license anyway. Make sure it's a valid SPDX identifier."))
	}

	return license, nil
}

// readInput reads a line of input from stdin
func (ps *ProjectSetup) readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSuffix(input, "\n"), nil
}

// isValidProjectName validates project name format
func (ps *ProjectSetup) isValidProjectName(name string) bool {
	// Project names should contain only letters, numbers, hyphens, and underscores
	// Should not start or end with hyphens or underscores
	// Should be between 1 and 100 characters
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Check for valid characters using pre-compiled regex
	return interactiveProjectNameRegex.MatchString(name)
}

// isValidLicense validates license identifier
func (ps *ProjectSetup) isValidLicense(license string) bool {
	// Common SPDX license identifiers
	commonLicenses := map[string]bool{
		"MIT":          true,
		"Apache-2.0":   true,
		"GPL-3.0":      true,
		"GPL-2.0":      true,
		"BSD-3-Clause": true,
		"BSD-2-Clause": true,
		"ISC":          true,
		"MPL-2.0":      true,
		"LGPL-3.0":     true,
		"LGPL-2.1":     true,
		"Unlicense":    true,
		"CC0-1.0":      true,
		"AGPL-3.0":     true,
		"Proprietary":  true,
	}

	return commonLicenses[license]
}

// ShowProjectSummary displays a summary of the collected project information
func (ps *ProjectSetup) ShowProjectSummary(config *models.ProjectConfig) {
	ps.output.QuietOutput("")
	ps.output.QuietOutput("📋 %s", ps.output.Highlight("Project Summary:"))
	ps.output.QuietOutput("%s", ps.output.Dim("================"))
	ps.output.QuietOutput("Name: %s", ps.output.Success(config.Name))

	if config.Organization != "" {
		ps.output.QuietOutput("Organization: %s", ps.output.Info(config.Organization))
	}

	if config.Description != "" {
		ps.output.QuietOutput("Description: %s", ps.output.Dim(config.Description))
	}

	ps.output.QuietOutput("License: %s", ps.output.Info(config.License))
}

// ConfirmProjectDetails asks the user to confirm the project details
func (ps *ProjectSetup) ConfirmProjectDetails(config *models.ProjectConfig) (bool, error) {
	ps.ShowProjectSummary(config)
	ps.output.QuietOutput("")
	ps.output.QuietOutput("❓ %s", ps.output.Info("Does this look correct? (Y/n):"))
	fmt.Print("   > ")

	response, err := ps.readInput()
	if err != nil {
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes", nil
}

// PromptForCorrections allows the user to correct project details
func (ps *ProjectSetup) PromptForCorrections(config *models.ProjectConfig) (*models.ProjectConfig, error) {
	ps.output.QuietOutput("")
	ps.output.QuietOutput("🔧 %s", ps.output.Info("What would you like to change?"))
	ps.output.QuietOutput("1. Project name")
	ps.output.QuietOutput("2. Organization")
	ps.output.QuietOutput("3. Description")
	ps.output.QuietOutput("4. License")
	ps.output.QuietOutput("5. Start over")
	ps.output.QuietOutput("")
	fmt.Print("Enter your choice (1-5): ")

	choice, err := ps.readInput()
	if err != nil {
		return nil, fmt.Errorf("failed to read choice: %w", err)
	}

	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		name, err := ps.promptProjectName()
		if err != nil {
			return nil, err
		}
		config.Name = name
	case "2":
		organization, err := ps.promptOrganization()
		if err != nil {
			return nil, err
		}
		config.Organization = organization
	case "3":
		description, err := ps.promptDescription()
		if err != nil {
			return nil, err
		}
		config.Description = description
	case "4":
		license, err := ps.promptLicense()
		if err != nil {
			return nil, err
		}
		config.License = license
	case "5":
		return ps.CollectProjectDetails()
	default:
		ps.output.WarningOutput("⚠️  Invalid choice. Please try again.")
		return ps.PromptForCorrections(config)
	}

	return config, nil
}
