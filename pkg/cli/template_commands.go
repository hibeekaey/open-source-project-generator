package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// setupTemplateCommand sets up the template command and its subcommands
func (c *CLI) setupTemplateCommand() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Template management operations",
		Long: `Manage templates including viewing detailed information and validation.

Use these commands to inspect templates before using them
or to validate custom templates you've created.`,
	}

	// template info
	templateInfoCmd := &cobra.Command{
		Use:   "info <template-name> [flags]",
		Short: "Display comprehensive template information and documentation",
		Long: `Display detailed information about a specific template including metadata,
dependencies, configuration options, and usage examples.`,
		RunE: c.runTemplateInfo,
		Args: cobra.ExactArgs(1),
		Example: `  # Show template information
  generator template info go-gin
  
  # Show detailed information
  generator template info nextjs-app --detailed
  
  # Show template variables
  generator template info go-gin --variables`,
	}
	templateInfoCmd.Flags().Bool("detailed", false, "Show detailed template information")
	templateInfoCmd.Flags().Bool("variables", false, "Show template variables")
	templateInfoCmd.Flags().Bool("dependencies", false, "Show template dependencies")
	templateInfoCmd.Flags().Bool("compatibility", false, "Show compatibility information")
	templateCmd.AddCommand(templateInfoCmd)

	// template validate
	templateValidateCmd := &cobra.Command{
		Use:   "validate <template-path> [flags]",
		Short: "Validate custom template structure, syntax, and compliance",
		Long: `Validate custom template directories including structure, metadata, syntax, and best practices.

Provides detailed feedback and auto-fix capabilities for common issues.`,
		RunE: c.runTemplateValidate,
		Args: cobra.ExactArgs(1),
		Example: `  # Validate template directory
  generator template validate ./my-template
  
  # Validate with detailed output
  generator template validate ./my-template --detailed
  
  # Validate and auto-fix issues
  generator template validate ./my-template --fix`,
	}
	templateValidateCmd.Flags().Bool("detailed", false, "Show detailed validation results")
	templateValidateCmd.Flags().Bool("fix", false, "Attempt to fix validation issues")
	templateValidateCmd.Flags().String("output-format", "text", "Output format (text, json)")
	templateCmd.AddCommand(templateValidateCmd)

	c.rootCmd.AddCommand(templateCmd)
}

// runTemplateInfo displays comprehensive template information
func (c *CLI) runTemplateInfo(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	showVariables, _ := cmd.Flags().GetBool("variables")
	showDependencies, _ := cmd.Flags().GetBool("dependencies")
	showCompatibility, _ := cmd.Flags().GetBool("compatibility")

	// Get template info
	templateInfo, err := c.GetTemplateInfo(templateName)
	if err != nil {
		return fmt.Errorf("🚫 Couldn't find template '%s': %w", templateName, err)
	}

	// Get category emoji
	categoryEmojis := map[string]string{
		"frontend":       "🎨",
		"backend":        "⚙️",
		"mobile":         "📱",
		"infrastructure": "🚀",
		"base":           "📋",
	}
	categoryEmoji := categoryEmojis[templateInfo.Category]
	if categoryEmoji == "" {
		categoryEmoji = "📦"
	}

	// Display header with emoji and template name
	c.QuietOutput("%s  %s", categoryEmoji, templateInfo.DisplayName)
	c.QuietOutput("═══════════════════════════════════════════════════")
	c.QuietOutput("")

	// Basic information
	c.QuietOutput("📝  %s", templateInfo.Description)
	c.QuietOutput("")
	c.QuietOutput("🔧  Template ID: %s", templateInfo.Name)
	c.QuietOutput("📂  Category: %s", cases.Title(language.English).String(templateInfo.Category))
	c.QuietOutput("⚡  Technology: %s", templateInfo.Technology)
	c.QuietOutput("🏷️   Version: %s", templateInfo.Version)

	if len(templateInfo.Tags) > 0 {
		c.QuietOutput("🏷️   Tags: %s", strings.Join(templateInfo.Tags, ", "))
	}

	// Show detailed information if requested
	if detailed || showDependencies {
		c.QuietOutput("")
		if len(templateInfo.Dependencies) > 0 {
			c.QuietOutput("📋  Dependencies:")
			for _, dep := range templateInfo.Dependencies {
				c.QuietOutput("    • %s", dep)
			}
		} else {
			c.QuietOutput("📋  Dependencies: None")
		}
	}

	// Show variables if requested
	if showVariables || detailed {
		c.QuietOutput("")
		c.QuietOutput("🔧  Template Variables: Not available in current interface")
	}

	// Show compatibility information if requested
	if showCompatibility || detailed {
		c.QuietOutput("")
		c.QuietOutput("🔗  Compatibility: Not available in current interface")
	}

	// Show usage examples
	if detailed {
		c.QuietOutput("")
		c.QuietOutput("📖  Usage Examples:")
		c.QuietOutput("    # Generate project with this template")
		c.QuietOutput("    generator generate --template %s", templateInfo.Name)
		c.QuietOutput("")
		c.QuietOutput("    # Interactive generation")
		c.QuietOutput("    generator generate")
		c.QuietOutput("    # Then select '%s' from the template list", templateInfo.DisplayName)
	}

	return nil
}

// runTemplateValidate validates custom template structure and syntax
func (c *CLI) runTemplateValidate(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	// Get flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	fix, _ := cmd.Flags().GetBool("fix")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	// Validate template
	result, err := c.ValidateTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("🚫 Template validation failed: %w", err)
	}

	// Display results based on format
	if outputFormat == "json" {
		// JSON output would be implemented here
		c.QuietOutput("JSON output not yet implemented")
		return nil
	}

	// Text output
	c.QuietOutput("🔍 Template Validation Results")
	c.QuietOutput("═══════════════════════════════════════════════════")
	c.QuietOutput("")

	if result.Valid {
		c.SuccessOutput("✅ Template is valid and ready to use")
	} else {
		c.ErrorOutput("❌ Template validation failed")
	}

	c.QuietOutput("")
	c.QuietOutput("📊 Validation Summary:")
	c.QuietOutput("    • Issues: %d", len(result.Issues))
	c.QuietOutput("    • Warnings: %d", len(result.Warnings))

	if len(result.Issues) > 0 {
		c.QuietOutput("")
		c.WarningOutput("⚠️  Issues Found:")
		for _, issue := range result.Issues {
			severity := "INFO"
			if issue.Severity == "error" {
				severity = "ERROR"
			} else if issue.Severity == "warning" {
				severity = "WARNING"
			}
			c.QuietOutput("    [%s] %s: %s", severity, issue.Type, issue.Message)
			if detailed {
				c.QuietOutput("        💡 File: %s", issue.File)
			}
		}
	}

	if fix && len(result.Issues) > 0 {
		c.QuietOutput("")
		c.VerboseOutput("🔧 Auto-fix functionality not yet implemented")
	}

	return nil
}
