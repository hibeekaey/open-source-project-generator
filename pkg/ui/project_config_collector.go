// Package ui provides interactive user interface components for project configuration collection.
//
// This file implements the ProjectConfigCollector which handles interactive collection
// of project metadata including name, description, organization, author, and license
// with comprehensive validation and preview functionality.
package ui

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectConfigCollector handles interactive collection of project configuration
type ProjectConfigCollector struct {
	ui        interfaces.InteractiveUIInterface
	validator *ProjectConfigValidator
	defaults  *ProjectConfigDefaults
	logger    interfaces.Logger
}

// ProjectConfigValidator provides validation for project configuration fields
type ProjectConfigValidator struct {
	projectNameRegex *regexp.Regexp
	emailRegex       *regexp.Regexp
	urlRegex         *regexp.Regexp
}

// ProjectConfigDefaults provides default values for project configuration
type ProjectConfigDefaults struct {
	License      string
	Author       string
	Email        string
	Organization string
}

// NewProjectConfigCollector creates a new project configuration collector
func NewProjectConfigCollector(ui interfaces.InteractiveUIInterface, logger interfaces.Logger) *ProjectConfigCollector {
	validator := &ProjectConfigValidator{
		projectNameRegex: regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`),
		emailRegex:       regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		urlRegex:         regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`),
	}

	defaults := &ProjectConfigDefaults{
		License: "MIT",
		Author:  "",
		Email:   "",
	}

	return &ProjectConfigCollector{
		ui:        ui,
		validator: validator,
		defaults:  defaults,
		logger:    logger,
	}
}

// CollectProjectConfiguration interactively collects project configuration
func (pcc *ProjectConfigCollector) CollectProjectConfiguration(ctx context.Context) (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{}

	// Collect project name
	name, err := pcc.collectProjectName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect project name: %w", err)
	}
	config.Name = name

	// Collect project description
	description, err := pcc.collectProjectDescription(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect project description: %w", err)
	}
	config.Description = description

	// Collect organization
	organization, err := pcc.collectOrganization(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect organization: %w", err)
	}
	config.Organization = organization

	// Collect author
	author, err := pcc.collectAuthor(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect author: %w", err)
	}
	config.Author = author

	// Collect email
	email, err := pcc.collectEmail(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect email: %w", err)
	}
	config.Email = email

	// Collect license
	license, err := pcc.collectLicense(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect license: %w", err)
	}
	config.License = license

	// Collect repository URL (optional)
	repository, err := pcc.collectRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect repository: %w", err)
	}
	config.Repository = repository

	// Show metadata usage preview
	if err := pcc.showMetadataPreview(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to show metadata preview: %w", err)
	}

	return config, nil
}

// collectProjectName collects and validates the project name
func (pcc *ProjectConfigCollector) collectProjectName(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:      "Project Name",
		Description: "Enter a unique name for your project (alphanumeric, hyphens, and underscores allowed)",
		Required:    true,
		Validator:   pcc.validator.ValidateProjectName,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Project Name Guidelines:
• Must start with a letter
• Can contain letters, numbers, hyphens (-), and underscores (_)
• Should be descriptive and unique
• Will be used for directory names and package identifiers
• Examples: my-awesome-app, user_management_system, BlogEngine`,
		MaxLength: 50,
		MinLength: 2,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("project name collection cancelled")
	}

	return result.Value, nil
}

// collectProjectDescription collects the project description
func (pcc *ProjectConfigCollector) collectProjectDescription(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:      "Project Description",
		Description: "Provide a brief description of your project (optional)",
		Required:    false,
		Validator:   pcc.validator.ValidateDescription,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Project Description Guidelines:
• Brief summary of what your project does
• Will be used in README.md, package.json, and other metadata files
• Should be clear and concise (1-2 sentences)
• Examples: "A modern web application for task management", "RESTful API for user authentication"`,
		MaxLength: 200,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("project description collection cancelled")
	}

	return result.Value, nil
}

// collectOrganization collects the organization name
func (pcc *ProjectConfigCollector) collectOrganization(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:       "Organization",
		Description:  "Enter your organization or company name (optional)",
		DefaultValue: pcc.defaults.Organization,
		Required:     false,
		Validator:    pcc.validator.ValidateOrganization,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Organization Guidelines:
• Your company, team, or personal brand name
• Will be used in copyright notices and package metadata
• Can be left empty for personal projects
• Examples: "Acme Corp", "My Development Team", "John Doe"`,
		MaxLength: 100,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("organization collection cancelled")
	}

	return result.Value, nil
}

// collectAuthor collects the author name
func (pcc *ProjectConfigCollector) collectAuthor(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:       "Author Name",
		Description:  "Enter the primary author's name",
		DefaultValue: pcc.defaults.Author,
		Required:     true,
		Validator:    pcc.validator.ValidateAuthor,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Author Guidelines:
• Your full name or preferred professional name
• Will be used in copyright notices, package metadata, and documentation
• Should be consistent across your projects
• Examples: "John Doe", "Jane Smith", "Alex Johnson"`,
		MaxLength: 100,
		MinLength: 2,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("author collection cancelled")
	}

	return result.Value, nil
}

// collectEmail collects the author's email
func (pcc *ProjectConfigCollector) collectEmail(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:       "Email Address",
		Description:  "Enter your email address (optional)",
		DefaultValue: pcc.defaults.Email,
		Required:     false,
		Validator:    pcc.validator.ValidateEmail,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText: `Email Guidelines:
• Your professional or project-related email address
• Will be used in package metadata and contact information
• Should be a valid email format
• Can be left empty if you prefer not to include it
• Examples: "john@example.com", "jane.smith@company.com"`,
		MaxLength: 100,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("email collection cancelled")
	}

	return result.Value, nil
}

// collectLicense collects the project license
func (pcc *ProjectConfigCollector) collectLicense(ctx context.Context) (string, error) {
	licenses := []interfaces.MenuOption{
		{Label: "MIT", Description: "Permissive license allowing commercial use", Value: "MIT"},
		{Label: "Apache 2.0", Description: "Permissive license with patent protection", Value: "Apache-2.0"},
		{Label: "GPL v3", Description: "Copyleft license requiring source disclosure", Value: "GPL-3.0"},
		{Label: "BSD 3-Clause", Description: "Permissive license with attribution requirement", Value: "BSD-3-Clause"},
		{Label: "ISC", Description: "Simple permissive license", Value: "ISC"},
		{Label: "Unlicense", Description: "Public domain dedication", Value: "Unlicense"},
		{Label: "Custom", Description: "Specify a custom license", Value: "custom"},
		{Label: "None", Description: "No license (all rights reserved)", Value: ""},
	}

	config := interfaces.MenuConfig{
		Title:       "Select License",
		Description: "Choose a license for your project",
		Options:     licenses,
		DefaultItem: 0, // MIT as default
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `License Selection Guidelines:
• MIT: Most permissive, allows commercial use with minimal restrictions
• Apache 2.0: Similar to MIT but includes patent protection
• GPL v3: Requires derivative works to be open source
• BSD 3-Clause: Permissive with attribution requirement
• ISC: Simple and permissive, similar to MIT
• Unlicense: Releases code to public domain
• Custom: Specify your own license identifier
• None: All rights reserved (not recommended for open source)`,
	}

	result, err := pcc.ui.ShowMenu(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "select" {
		return "", fmt.Errorf("license selection cancelled")
	}

	selectedLicense := result.SelectedValue.(string)

	// Handle custom license
	if selectedLicense == "custom" {
		customConfig := interfaces.TextPromptConfig{
			Prompt:      "Custom License",
			Description: "Enter your custom license identifier (e.g., 'Proprietary', 'Custom-1.0')",
			Required:    true,
			Validator:   pcc.validator.ValidateLicense,
			AllowBack:   true,
			AllowQuit:   true,
			ShowHelp:    true,
			HelpText: `Custom License Guidelines:
• Use standard SPDX identifiers when possible
• For proprietary licenses, use "Proprietary" or your company name
• For custom licenses, include version if applicable
• Examples: "Proprietary", "Acme-Corp-1.0", "Custom-MIT-Variant"`,
			MaxLength: 50,
		}

		customResult, err := pcc.ui.PromptText(ctx, customConfig)
		if err != nil {
			return "", err
		}

		if customResult.Cancelled || customResult.Action != "submit" {
			return "", fmt.Errorf("custom license collection cancelled")
		}

		return customResult.Value, nil
	}

	return selectedLicense, nil
}

// collectRepository collects the repository URL
func (pcc *ProjectConfigCollector) collectRepository(ctx context.Context) (string, error) {
	config := interfaces.TextPromptConfig{
		Prompt:      "Repository URL",
		Description: "Enter the repository URL (optional, e.g., https://github.com/user/repo)",
		Required:    false,
		Validator:   pcc.validator.ValidateRepository,
		AllowBack:   true,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText: `Repository URL Guidelines:
• Full URL to your project's repository
• Will be used in package metadata and documentation
• Should be publicly accessible if you want others to contribute
• Examples: "https://github.com/user/repo", "https://gitlab.com/user/project"`,
		MaxLength: 200,
	}

	result, err := pcc.ui.PromptText(ctx, config)
	if err != nil {
		return "", err
	}

	if result.Cancelled || result.Action != "submit" {
		return "", fmt.Errorf("repository collection cancelled")
	}

	return result.Value, nil
}

// showMetadataPreview shows how the collected metadata will be used
func (pcc *ProjectConfigCollector) showMetadataPreview(ctx context.Context, config *models.ProjectConfig) error {
	previewText := pcc.generateMetadataPreview(config)

	// Show preview in a table format
	headers := []string{"File/Location", "Usage", "Value"}
	rows := [][]string{
		{"README.md", "Project title", config.Name},
		{"README.md", "Description", config.Description},
		{"package.json", "Name field", strings.ToLower(strings.ReplaceAll(config.Name, " ", "-"))},
		{"package.json", "Author field", fmt.Sprintf("%s <%s>", config.Author, config.Email)},
		{"LICENSE", "License type", config.License},
		{"Copyright notices", "Organization", config.Organization},
		{"package.json", "Repository", config.Repository},
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Metadata Usage Preview",
		Headers: headers,
		Rows:    rows,
	}

	if err := pcc.ui.ShowTable(ctx, tableConfig); err != nil {
		return fmt.Errorf("failed to show metadata preview table: %w", err)
	}

	// Show detailed preview text
	fmt.Println("\nDetailed Preview:")
	fmt.Println(previewText)

	// Confirm with user
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Does this metadata look correct",
		Description:  "Review the information above and confirm if it's accurate",
		DefaultValue: true,
		AllowBack:    true,
		AllowQuit:    true,
		ShowHelp:     true,
		HelpText:     "If the information is incorrect, select 'No' to go back and modify your entries.",
	}

	result, err := pcc.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if result.Cancelled || result.Action != "confirm" || !result.Confirmed {
		return fmt.Errorf("metadata preview not confirmed")
	}

	return nil
}

// generateMetadataPreview generates a preview of how metadata will be used
func (pcc *ProjectConfigCollector) generateMetadataPreview(config *models.ProjectConfig) string {
	var preview strings.Builder

	preview.WriteString("Your project metadata will be used in the following ways:\n\n")

	// README.md preview
	preview.WriteString("📄 README.md:\n")
	preview.WriteString(fmt.Sprintf("# %s\n\n", config.Name))
	if config.Description != "" {
		preview.WriteString(fmt.Sprintf("%s\n\n", config.Description))
	}
	preview.WriteString("## Installation\n...\n\n")

	// package.json preview (for frontend projects)
	preview.WriteString("📦 package.json:\n")
	preview.WriteString("{\n")
	preview.WriteString(fmt.Sprintf("  \"name\": \"%s\",\n", strings.ToLower(strings.ReplaceAll(config.Name, " ", "-"))))
	if config.Description != "" {
		preview.WriteString(fmt.Sprintf("  \"description\": \"%s\",\n", config.Description))
	}
	preview.WriteString("  \"version\": \"1.0.0\",\n")
	if config.Author != "" {
		authorField := config.Author
		if config.Email != "" {
			authorField = fmt.Sprintf("%s <%s>", config.Author, config.Email)
		}
		preview.WriteString(fmt.Sprintf("  \"author\": \"%s\",\n", authorField))
	}
	if config.License != "" {
		preview.WriteString(fmt.Sprintf("  \"license\": \"%s\",\n", config.License))
	}
	if config.Repository != "" {
		preview.WriteString(fmt.Sprintf("  \"repository\": \"%s\",\n", config.Repository))
	}
	preview.WriteString("  ...\n}\n\n")

	// License file preview
	if config.License != "" {
		preview.WriteString("📜 LICENSE:\n")
		preview.WriteString(fmt.Sprintf("%s License\n\n", config.License))
		if config.Organization != "" {
			preview.WriteString(fmt.Sprintf("Copyright (c) 2024 %s\n\n", config.Organization))
		} else if config.Author != "" {
			preview.WriteString(fmt.Sprintf("Copyright (c) 2024 %s\n\n", config.Author))
		}
		preview.WriteString("[License text will be generated based on selected license]\n\n")
	}

	// Go module preview (for backend projects)
	preview.WriteString("🐹 go.mod:\n")
	moduleName := strings.ToLower(strings.ReplaceAll(config.Name, " ", "-"))
	if config.Repository != "" {
		// Extract module name from repository URL
		if strings.Contains(config.Repository, "github.com") {
			parts := strings.Split(config.Repository, "/")
			if len(parts) >= 2 {
				moduleName = fmt.Sprintf("github.com/%s/%s", parts[len(parts)-2], parts[len(parts)-1])
			}
		}
	}
	preview.WriteString(fmt.Sprintf("module %s\n\n", moduleName))
	preview.WriteString("go 1.21\n\n")

	return preview.String()
}

// Validation methods for ProjectConfigValidator

// ValidateProjectName validates the project name according to naming conventions
func (v *ProjectConfigValidator) ValidateProjectName(name string) error {
	if name == "" {
		return interfaces.NewValidationError("name", name, "Project name is required", "required").
			WithSuggestions("Enter a descriptive name for your project")
	}

	if len(name) < 2 {
		return interfaces.NewValidationError("name", name, "Project name must be at least 2 characters long", "min_length").
			WithSuggestions("Use a longer, more descriptive name")
	}

	if len(name) > 50 {
		return interfaces.NewValidationError("name", name, "Project name must be at most 50 characters long", "max_length").
			WithSuggestions("Use a shorter, more concise name")
	}

	if !v.projectNameRegex.MatchString(name) {
		return interfaces.NewValidationError("name", name, "Project name must start with a letter and contain only letters, numbers, hyphens, and underscores", "invalid_format").
			WithSuggestions(
				"Start with a letter (a-z, A-Z)",
				"Use only letters, numbers, hyphens (-), and underscores (_)",
				"Examples: my-project, UserManager, blog_engine",
			)
	}

	// Check for reserved names
	reservedNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	lowerName := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			return interfaces.NewValidationError("name", name, fmt.Sprintf("'%s' is a reserved name and cannot be used", name), "reserved_name").
				WithSuggestions("Choose a different name that doesn't conflict with system reserved words")
		}
	}

	return nil
}

// ValidateDescription validates the project description
func (v *ProjectConfigValidator) ValidateDescription(description string) error {
	if len(description) > 200 {
		return interfaces.NewValidationError("description", description, "Description must be at most 200 characters long", "max_length").
			WithSuggestions("Keep the description concise and focused on the main purpose")
	}

	// Check for common issues
	if strings.Contains(description, "\n") {
		return interfaces.NewValidationError("description", description, "Description should be a single line", "multiline").
			WithSuggestions("Use a brief, single-line description")
	}

	return nil
}

// ValidateOrganization validates the organization name
func (v *ProjectConfigValidator) ValidateOrganization(organization string) error {
	if len(organization) > 100 {
		return interfaces.NewValidationError("organization", organization, "Organization name must be at most 100 characters long", "max_length").
			WithSuggestions("Use a shorter organization name or abbreviation")
	}

	return nil
}

// ValidateAuthor validates the author name
func (v *ProjectConfigValidator) ValidateAuthor(author string) error {
	if author == "" {
		return interfaces.NewValidationError("author", author, "Author name is required", "required").
			WithSuggestions("Enter your full name or preferred professional name")
	}

	if len(author) < 2 {
		return interfaces.NewValidationError("author", author, "Author name must be at least 2 characters long", "min_length").
			WithSuggestions("Enter your full name")
	}

	if len(author) > 100 {
		return interfaces.NewValidationError("author", author, "Author name must be at most 100 characters long", "max_length").
			WithSuggestions("Use a shorter version of your name")
	}

	// Check for suspicious patterns
	if strings.Contains(author, "@") {
		return interfaces.NewValidationError("author", author, "Author name should not contain email addresses", "invalid_format").
			WithSuggestions("Enter just your name, email will be collected separately")
	}

	return nil
}

// ValidateEmail validates the email address format
func (v *ProjectConfigValidator) ValidateEmail(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	if len(email) > 100 {
		return interfaces.NewValidationError("email", email, "Email address must be at most 100 characters long", "max_length").
			WithSuggestions("Use a shorter email address")
	}

	if !v.emailRegex.MatchString(email) {
		return interfaces.NewValidationError("email", email, "Invalid email address format", "invalid_format").
			WithSuggestions(
				"Use a valid email format: user@domain.com",
				"Check for typos in the domain name",
				"Ensure the email has both @ symbol and domain extension",
			)
	}

	return nil
}

// ValidateLicense validates the license identifier
func (v *ProjectConfigValidator) ValidateLicense(license string) error {
	if license == "" {
		return nil // License can be empty (all rights reserved)
	}

	if len(license) > 50 {
		return interfaces.NewValidationError("license", license, "License identifier must be at most 50 characters long", "max_length").
			WithSuggestions("Use a standard SPDX license identifier or shorter custom identifier")
	}

	// Check for common license identifiers (case-insensitive)
	commonLicenses := []string{
		"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "ISC", "Unlicense",
		"GPL-2.0", "LGPL-3.0", "LGPL-2.1", "BSD-2-Clause", "MPL-2.0",
		"Proprietary", "Custom",
	}

	lowerLicense := strings.ToLower(license)
	for _, common := range commonLicenses {
		if strings.ToLower(common) == lowerLicense {
			return nil // Valid common license
		}
	}

	// If not a common license, just warn but don't fail validation
	return nil
}

// ValidateRepository validates the repository URL
func (v *ProjectConfigValidator) ValidateRepository(repository string) error {
	if repository == "" {
		return nil // Repository is optional
	}

	if len(repository) > 200 {
		return interfaces.NewValidationError("repository", repository, "Repository URL must be at most 200 characters long", "max_length").
			WithSuggestions("Use a shorter URL or check for extra characters")
	}

	if !v.urlRegex.MatchString(repository) {
		return interfaces.NewValidationError("repository", repository, "Invalid repository URL format", "invalid_format").
			WithSuggestions(
				"Use a complete URL starting with http:// or https://",
				"Examples: https://github.com/user/repo, https://gitlab.com/user/project",
				"Ensure the URL is publicly accessible",
			)
	}

	// Check for common repository hosting services
	supportedHosts := []string{"github.com", "gitlab.com", "bitbucket.org", "codeberg.org", "sourceforge.net"}
	isSupported := false
	for _, host := range supportedHosts {
		if strings.Contains(strings.ToLower(repository), host) {
			isSupported = true
			break
		}
	}

	if !isSupported {
		// Not an error, but could provide a suggestion
		return nil
	}

	return nil
}

// SetDefaults sets default values for the configuration collector
func (pcc *ProjectConfigCollector) SetDefaults(defaults *ProjectConfigDefaults) {
	if defaults != nil {
		pcc.defaults = defaults
	}
}

// GetDefaults returns the current default values
func (pcc *ProjectConfigCollector) GetDefaults() *ProjectConfigDefaults {
	return pcc.defaults
}

// LoadDefaultsFromEnvironment loads default values from environment variables
func (pcc *ProjectConfigCollector) LoadDefaultsFromEnvironment() {
	// This could load from environment variables like GIT_AUTHOR_NAME, GIT_AUTHOR_EMAIL, etc.
	// For now, we'll keep it simple and just provide a placeholder
	if pcc.defaults == nil {
		pcc.defaults = &ProjectConfigDefaults{}
	}
}
