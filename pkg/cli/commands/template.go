// Package commands provides individual command implementations for the CLI interface.
//
// This module contains the TemplateCommands struct and its associated functionality,
// extracted from the main CLI handlers to improve modularity and maintainability.
package commands

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TemplateCLI defines the CLI methods needed by TemplateCommands.
type TemplateCLI interface {
	// Template operations
	ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error)
	GetTemplateInfo(name string) (*interfaces.TemplateInfo, error)
	ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error)

	// Output methods
	VerboseOutput(format string, args ...interface{})
	DebugOutput(format string, args ...interface{})
	QuietOutput(format string, args ...interface{})
	ErrorOutput(format string, args ...interface{})
	WarningOutput(format string, args ...interface{})
	SuccessOutput(format string, args ...interface{})

	// Color formatting methods
	Error(text string) string
	Warning(text string) string
	Info(text string) string
	Success(text string) string
	Highlight(text string) string
	Dim(text string) string

	// Utility methods
	IsQuietMode() bool
	OutputMachineReadable(data interface{}, format string) error

	// Error creation methods
	CreateTemplateError(message string, templateName string) error
	OutputSuccess(message string, data interface{}, operation string, args []string) error

	// Flag handling
	IsNonInteractiveMode(cmd *cobra.Command) bool
}

// TemplateCommands handles template management command functionality.
//
// The TemplateCommands provides centralized template command execution including:
//   - Template listing with filtering and search
//   - Template information display
//   - Template validation and verification
//   - Result formatting and output
type TemplateCommands struct {
	cli TemplateCLI
}

// NewTemplateCommands creates a new TemplateCommands instance.
func NewTemplateCommands(cli TemplateCLI) *TemplateCommands {
	return &TemplateCommands{
		cli: cli,
	}
}

// ExecuteList handles the list-templates command execution.
func (tc *TemplateCommands) ExecuteList(cmd *cobra.Command, args []string) error {
	// Get flags
	category, _ := cmd.Flags().GetString("category")
	technology, _ := cmd.Flags().GetString("technology")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	search, _ := cmd.Flags().GetString("search")
	detailed, _ := cmd.Flags().GetBool("detailed")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = tc.cli.IsNonInteractiveMode(cmd)
	}

	// Create template filter
	filter := interfaces.TemplateFilter{
		Category:   category,
		Technology: technology,
		Tags:       tags,
	}

	var templates []interfaces.TemplateInfo
	var err error

	// Search or list templates
	if search != "" {
		// For now, use ListTemplates and filter by search term
		// This would be enhanced when SearchTemplates is fully implemented
		allTemplates, err := tc.cli.ListTemplates(filter)
		if err != nil {
			return tc.cli.CreateTemplateError("failed to search templates", search)
		}

		// Simple search filtering
		templates = []interfaces.TemplateInfo{}
		for _, template := range allTemplates {
			if strings.Contains(strings.ToLower(template.Name), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(template.Description), strings.ToLower(search)) {
				templates = append(templates, template)
			}
		}
	} else {
		templates, err = tc.cli.ListTemplates(filter)
	}

	if err != nil {
		return tc.cli.CreateTemplateError("failed to list templates", "")
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
		"filter":    filter,
		"search":    search,
		"detailed":  detailed,
	}

	// Output in machine-readable format if requested
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return tc.cli.OutputMachineReadable(responseData, outputFormat)
	}

	// Human-readable output
	if len(templates) == 0 {
		tc.cli.QuietOutput("🔍 No templates found matching your criteria")
		tc.cli.QuietOutput("💡 Try a different category or search term")
		return tc.cli.OutputSuccess("No templates found", responseData, "list-templates", []string{})
	}

	if !tc.cli.IsQuietMode() {
		// Group templates by category for better organization
		categories := make(map[string][]interfaces.TemplateInfo)
		for _, template := range templates {
			categories[template.Category] = append(categories[template.Category], template)
		}

		tc.cli.QuietOutput("📦 Available Templates (%d found)", len(templates))
		tc.cli.QuietOutput("")

		// Display templates grouped by category
		categoryOrder := []string{"frontend", "backend", "mobile", "infrastructure", "base"}
		categoryEmojis := map[string]string{
			"frontend":       "🎨",
			"backend":        "⚙️ ",
			"mobile":         "📱",
			"infrastructure": "🚀",
			"base":           "📋",
		}

		for _, cat := range categoryOrder {
			if templates, exists := categories[cat]; exists {
				emoji := categoryEmojis[cat]
				if emoji == "" {
					emoji = "📦"
				}
				tc.cli.QuietOutput("%s  %s Templates:", emoji, cases.Title(language.English).String(cat))

				for _, template := range templates {
					if detailed {
						tc.cli.QuietOutput("  • %s (%s)", template.DisplayName, template.Name)
						tc.cli.QuietOutput("    %s", template.Description)
						if len(template.Tags) > 0 {
							tc.cli.QuietOutput("    🏷️   %s", strings.Join(template.Tags, ", "))
						}
						if template.Metadata.Author != "" {
							tc.cli.VerboseOutput("    👤  %s", template.Metadata.Author)
						}
						if len(template.Dependencies) > 0 {
							tc.cli.VerboseOutput("    📋  Dependencies: %s", strings.Join(template.Dependencies, ", "))
						}
					} else {
						tc.cli.QuietOutput("  • %s - %s", template.DisplayName, template.Description)
					}
				}
				tc.cli.QuietOutput("")
			}
		}

		// Display any templates from categories not in our predefined list
		for cat, templates := range categories {
			found := false
			for _, knownCat := range categoryOrder {
				if cat == knownCat {
					found = true
					break
				}
			}
			if !found {
				tc.cli.QuietOutput("📦  %s Templates:", cases.Title(language.English).String(cat))
				for _, template := range templates {
					if detailed {
						tc.cli.QuietOutput("  • %s (%s)", template.DisplayName, template.Name)
						tc.cli.QuietOutput("    %s", template.Description)
						if len(template.Tags) > 0 {
							tc.cli.QuietOutput("    🏷️   %s", strings.Join(template.Tags, ", "))
						}
					} else {
						tc.cli.QuietOutput("  • %s - %s", template.DisplayName, template.Description)
					}
				}
				tc.cli.QuietOutput("")
			}
		}

		tc.cli.QuietOutput("💡 Use --detailed for more information")
		tc.cli.QuietOutput("🔍 Use --search <term> to find specific templates")
	}

	return tc.cli.OutputSuccess(fmt.Sprintf("Listed %d templates", len(templates)), responseData, "list-templates", []string{})
}

// ExecuteInfo handles the template info command execution.
func (tc *TemplateCommands) ExecuteInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	variables, _ := cmd.Flags().GetBool("variables")
	dependencies, _ := cmd.Flags().GetBool("dependencies")
	compatibility, _ := cmd.Flags().GetBool("compatibility")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = tc.cli.IsNonInteractiveMode(cmd)
	}

	// Get template information
	templateInfo, err := tc.cli.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			tc.cli.Error("Template not found or inaccessible."),
			tc.cli.Info("Use 'generator list-templates' to see available templates"))
	}

	// Output in machine-readable format if requested
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return tc.cli.OutputMachineReadable(templateInfo, outputFormat)
	}

	// Human-readable output
	tc.displayTemplateInfo(templateInfo, detailed, variables, dependencies, compatibility)

	return nil
}

// ExecuteValidate handles the template validate command execution.
func (tc *TemplateCommands) ExecuteValidate(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	fix, _ := cmd.Flags().GetBool("fix")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Get global flags
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
	globalOutputFormat, _ := cmd.Flags().GetString("output-format")

	// Auto-detect non-interactive mode
	if !nonInteractive {
		nonInteractive = tc.cli.IsNonInteractiveMode(cmd)
	}

	// Use global output format if command-specific format not set
	if outputFormat == "text" && globalOutputFormat != "text" {
		outputFormat = globalOutputFormat
	}

	// Validate template
	result, err := tc.cli.ValidateTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("🚫 %s %s",
			tc.cli.Error("Template validation failed."),
			tc.cli.Info("Check that the template path exists and is accessible"))
	}

	// Output in machine-readable format if requested
	if nonInteractive && (outputFormat == "json" || outputFormat == "yaml") {
		return tc.cli.OutputMachineReadable(result, outputFormat)
	}

	// Human-readable output
	tc.displayValidationResults(result, templatePath, detailed, fix)

	// Return appropriate exit code
	if !result.Valid {
		return fmt.Errorf("🚫 %s %s",
			tc.cli.Error("Template validation failed."),
			tc.cli.Info("Please fix the issues and try again"))
	}

	return nil
}

// displayTemplateInfo displays template information in human-readable format.
func (tc *TemplateCommands) displayTemplateInfo(templateInfo *interfaces.TemplateInfo, detailed, variables, dependencies, compatibility bool) {
	// Category emojis for better visual organization
	categoryEmojis := map[string]string{
		"frontend":       "🎨",
		"backend":        "⚙️ ",
		"mobile":         "📱",
		"infrastructure": "🚀",
		"base":           "📋",
	}
	categoryEmoji := categoryEmojis[templateInfo.Category]
	if categoryEmoji == "" {
		categoryEmoji = "📦"
	}

	// Display header with emoji and template name
	tc.cli.QuietOutput("%s  %s", categoryEmoji, templateInfo.DisplayName)
	tc.cli.QuietOutput("═══════════════════════════════════════════════════")
	tc.cli.QuietOutput("")

	// Basic information
	tc.cli.QuietOutput("📝  %s", templateInfo.Description)
	tc.cli.QuietOutput("")
	tc.cli.QuietOutput("🔧  Template ID: %s", templateInfo.Name)
	tc.cli.QuietOutput("📂  Category: %s", cases.Title(language.English).String(templateInfo.Category))
	tc.cli.QuietOutput("⚡  Technology: %s", templateInfo.Technology)
	tc.cli.QuietOutput("🏷️   Version: %s", templateInfo.Version)

	if len(templateInfo.Tags) > 0 {
		tc.cli.QuietOutput("🏷️   Tags: %s", strings.Join(templateInfo.Tags, ", "))
	}

	// Show variables if requested or in detailed mode
	if detailed || variables {
		tc.cli.QuietOutput("")
		if len(templateInfo.Metadata.Variables) > 0 {
			tc.cli.QuietOutput("📝  Template Variables:")
			for key, description := range templateInfo.Metadata.Variables {
				tc.cli.QuietOutput("    • %s: %s", key, description)
			}
		}
	}

	// Show dependencies if requested or in detailed mode
	if detailed || dependencies {
		tc.cli.QuietOutput("")
		if len(templateInfo.Dependencies) > 0 {
			tc.cli.QuietOutput("📋  Dependencies:")
			for _, dep := range templateInfo.Dependencies {
				tc.cli.QuietOutput("    • %s", dep)
			}
		}
	}

	// Show detailed metadata if requested
	if detailed {
		tc.cli.QuietOutput("")
		tc.cli.QuietOutput("👤  Author: %s", templateInfo.Metadata.Author)
		tc.cli.QuietOutput("📄  License: %s", templateInfo.Metadata.License)
		if templateInfo.Metadata.Repository != "" {
			tc.cli.QuietOutput("🔗  Repository: %s", templateInfo.Metadata.Repository)
		}
		if templateInfo.Metadata.Homepage != "" {
			tc.cli.QuietOutput("🌐  Homepage: %s", templateInfo.Metadata.Homepage)
		}
		if len(templateInfo.Metadata.Keywords) > 0 {
			tc.cli.QuietOutput("🔍  Keywords: %s", strings.Join(templateInfo.Metadata.Keywords, ", "))
		}
		if !templateInfo.Metadata.Created.IsZero() {
			tc.cli.QuietOutput("📅  Created: %s", templateInfo.Metadata.Created.Format("2006-01-02"))
		}
		if !templateInfo.Metadata.Updated.IsZero() {
			tc.cli.QuietOutput("🔄  Updated: %s", templateInfo.Metadata.Updated.Format("2006-01-02"))
		}
	}

	// Show compatibility information if requested
	if compatibility {
		tc.cli.QuietOutput("")
		tc.cli.QuietOutput("🔍  Compatibility:")
		tc.cli.QuietOutput("    Generator version: %s+", "latest") // This would be replaced with actual version
		tc.cli.QuietOutput("    Template version: %s", templateInfo.Version)
	}
}

// displayValidationResults displays template validation results in human-readable format.
func (tc *TemplateCommands) displayValidationResults(result *interfaces.TemplateValidationResult, templatePath string, detailed, fix bool) {
	tc.cli.QuietOutput("🔍 Template Validation Results")
	tc.cli.QuietOutput("==============================")
	tc.cli.QuietOutput("Template: %s", templatePath)

	if result.Valid {
		tc.cli.QuietOutput("Status: %s", tc.cli.Success("✅ Valid"))
	} else {
		tc.cli.QuietOutput("Status: %s", tc.cli.Error("❌ Invalid"))
	}

	tc.cli.QuietOutput("Files checked: %d", result.Summary.TotalFiles)
	tc.cli.QuietOutput("Issues: %d", len(result.Issues))
	tc.cli.QuietOutput("Warnings: %d", len(result.Warnings))

	if len(result.Issues) > 0 {
		tc.cli.QuietOutput("\n🚨 Issues:")
		for _, issue := range result.Issues {
			tc.cli.QuietOutput("  - %s: %s", issue.Severity, issue.Message)
			if detailed && issue.File != "" {
				tc.cli.QuietOutput("    File: %s:%d:%d", issue.File, issue.Line, issue.Column)
			}
		}
	}

	if len(result.Warnings) > 0 {
		tc.cli.QuietOutput("\n⚠️  Warnings:")
		for _, warning := range result.Warnings {
			tc.cli.QuietOutput("  - %s: %s", warning.Severity, warning.Message)
			if detailed && warning.File != "" {
				tc.cli.QuietOutput("    File: %s:%d:%d", warning.File, warning.Line, warning.Column)
			}
		}
	}

	if fix && !result.Valid {
		tc.cli.QuietOutput("\n🔧 Attempting to fix issues...")
		// In a real implementation, this would attempt to fix the issues
		tc.cli.QuietOutput("Auto-fix functionality would be implemented here")
		tc.cli.QuietOutput("This would attempt to resolve common template issues automatically")
	}

	// Additional suggestions could be added here based on validation results
	if !result.Valid && len(result.Issues) > 0 {
		tc.cli.QuietOutput("\n💡 Suggestions:")
		tc.cli.QuietOutput("  - Review the issues above and fix them")
		tc.cli.QuietOutput("  - Check template documentation for best practices")
		if fix {
			tc.cli.QuietOutput("  - Use --fix flag to attempt automatic fixes")
		}
	}
}

// SetupListFlags sets up the list-templates command flags.
func (tc *TemplateCommands) SetupListFlags(cmd *cobra.Command) {
	cmd.Flags().String("category", "", "Filter by category (frontend, backend, mobile, infrastructure)")
	cmd.Flags().String("technology", "", "Filter by technology (go, nodejs, react, etc.)")
	cmd.Flags().StringSlice("tags", []string{}, "Filter by tags")
	cmd.Flags().String("search", "", "Search templates by name or description")
	cmd.Flags().Bool("detailed", false, "Show detailed template information")
}

// SetupInfoFlags sets up the template info command flags.
func (tc *TemplateCommands) SetupInfoFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("detailed", false, "Show detailed template information")
	cmd.Flags().Bool("variables", false, "Show template variables")
	cmd.Flags().Bool("dependencies", false, "Show template dependencies")
	cmd.Flags().Bool("compatibility", false, "Show compatibility information")
}

// SetupValidateFlags sets up the template validate command flags.
func (tc *TemplateCommands) SetupValidateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("detailed", false, "Show detailed validation results")
	cmd.Flags().Bool("fix", false, "Attempt to fix validation issues")
	cmd.Flags().String("output-format", "text", "Output format (text, json)")
}
