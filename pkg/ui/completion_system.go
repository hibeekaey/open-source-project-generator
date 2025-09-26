// Package ui provides completion system for interactive CLI generation.
//
// This file implements the CompletionSystem which provides completion summaries,
// next steps guidance, and final project information after successful generation.
package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// CompletionSystem manages completion summaries and next steps
type CompletionSystem struct {
	ui         interfaces.InteractiveUIInterface
	helpSystem *HelpSystem
	logger     interfaces.Logger
	config     *CompletionConfig
}

// CompletionConfig defines configuration for the completion system
type CompletionConfig struct {
	ShowGeneratedFiles    bool `json:"show_generated_files"`
	ShowNextSteps         bool `json:"show_next_steps"`
	ShowCommands          bool `json:"show_commands"`
	ShowTroubleshooting   bool `json:"show_troubleshooting"`
	ShowAdditionalInfo    bool `json:"show_additional_info"`
	EnableInteractiveHelp bool `json:"enable_interactive_help"`
}

// GenerationResult contains the results of project generation
type GenerationResult struct {
	ProjectConfig     *models.ProjectConfig `json:"project_config"`
	OutputDirectory   string                `json:"output_directory"`
	SelectedTemplates []TemplateSelection   `json:"selected_templates"`
	GeneratedFiles    []GeneratedFileInfo   `json:"generated_files"`
	TotalFiles        int                   `json:"total_files"`
	TotalSize         int64                 `json:"total_size"`
	Duration          string                `json:"duration"`
	Errors            []string              `json:"errors,omitempty"`
	Warnings          []string              `json:"warnings,omitempty"`
}

// GeneratedFileInfo contains information about a generated file
type GeneratedFileInfo struct {
	Path         string `json:"path"`
	RelativePath string `json:"relative_path"`
	Type         string `json:"type"` // "file", "directory", "symlink"
	Size         int64  `json:"size"`
	Template     string `json:"template"`
	Category     string `json:"category"`
}

// NewCompletionSystem creates a new completion system
func NewCompletionSystem(ui interfaces.InteractiveUIInterface, helpSystem *HelpSystem, logger interfaces.Logger, config *CompletionConfig) *CompletionSystem {
	if config == nil {
		config = &CompletionConfig{
			ShowGeneratedFiles:    true,
			ShowNextSteps:         true,
			ShowCommands:          true,
			ShowTroubleshooting:   true,
			ShowAdditionalInfo:    true,
			EnableInteractiveHelp: true,
		}
	}

	return &CompletionSystem{
		ui:         ui,
		helpSystem: helpSystem,
		logger:     logger,
		config:     config,
	}
}

// ShowCompletionSummary displays a comprehensive completion summary
func (cs *CompletionSystem) ShowCompletionSummary(ctx context.Context, result *GenerationResult) error {
	summary := cs.buildCompletionSummary(result)

	// Show the completion summary using the help system
	if err := cs.helpSystem.ShowCompletionSummary(ctx, summary); err != nil {
		return fmt.Errorf("failed to show completion summary: %w", err)
	}

	// Offer interactive help if enabled
	if cs.config.EnableInteractiveHelp {
		if err := cs.offerInteractiveHelp(ctx, result); err != nil {
			cs.logger.WarnWithFields("Failed to offer interactive help", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return nil
}

// buildCompletionSummary builds a completion summary from generation results
func (cs *CompletionSystem) buildCompletionSummary(result *GenerationResult) *CompletionSummary {
	summary := &CompletionSummary{
		Title:       "🎉 Project Generation Complete!",
		Description: fmt.Sprintf("Successfully generated '%s' project", result.ProjectConfig.Name),
		Metadata: map[string]interface{}{
			"project_name":     result.ProjectConfig.Name,
			"output_directory": result.OutputDirectory,
			"duration":         result.Duration,
			"total_files":      result.TotalFiles,
			"total_size":       cs.formatSize(result.TotalSize),
		},
	}

	// Add generated items
	if cs.config.ShowGeneratedFiles {
		summary.GeneratedItems = cs.buildGeneratedItems(result)
	}

	// Add next steps
	if cs.config.ShowNextSteps {
		summary.NextSteps = cs.buildNextSteps(result)
	}

	// Add additional information
	if cs.config.ShowAdditionalInfo {
		summary.AdditionalInfo = cs.buildAdditionalInfo(result)
	}

	return summary
}

// buildGeneratedItems builds the list of generated items
func (cs *CompletionSystem) buildGeneratedItems(result *GenerationResult) []GeneratedItem {
	items := []GeneratedItem{}

	// Add main project directory
	items = append(items, GeneratedItem{
		Type:        "directory",
		Name:        result.ProjectConfig.Name,
		Path:        result.OutputDirectory,
		Description: "Main project directory",
	})

	// Group files by category
	categories := make(map[string][]GeneratedFileInfo)
	for _, file := range result.GeneratedFiles {
		category := file.Category
		if category == "" {
			category = cs.categorizeFile(file.RelativePath)
		}
		categories[category] = append(categories[category], file)
	}

	// Add items by category
	for category, files := range categories {
		if len(files) > 5 {
			// Summarize large categories
			totalSize := int64(0)
			for _, file := range files {
				totalSize += file.Size
			}
			items = append(items, GeneratedItem{
				Type:        "category",
				Name:        fmt.Sprintf("%s (%d files)", category, len(files)),
				Path:        filepath.Join(result.OutputDirectory, cs.getCategoryPath(category)),
				Description: fmt.Sprintf("%s files and configurations", category),
				Size:        cs.formatSize(totalSize),
			})
		} else {
			// List individual files for small categories
			for _, file := range files {
				items = append(items, GeneratedItem{
					Type:        file.Type,
					Name:        filepath.Base(file.RelativePath),
					Path:        file.Path,
					Description: fmt.Sprintf("%s file", category),
					Size:        cs.formatSize(file.Size),
				})
			}
		}
	}

	return items
}

// buildNextSteps builds the list of next steps
func (cs *CompletionSystem) buildNextSteps(result *GenerationResult) []NextStep {
	steps := []NextStep{}

	// Add template-specific next steps
	for _, template := range result.SelectedTemplates {
		if !template.Selected {
			continue
		}

		templateSteps := cs.getTemplateNextSteps(template.Template, result)
		steps = append(steps, templateSteps...)
	}

	// Add general next steps
	generalSteps := cs.getGeneralNextSteps(result)
	steps = append(steps, generalSteps...)

	return steps
}

// getTemplateNextSteps returns next steps specific to a template
func (cs *CompletionSystem) getTemplateNextSteps(template interfaces.TemplateInfo, result *GenerationResult) []NextStep {
	steps := []NextStep{}

	switch template.Category {
	case "backend":
		if template.Technology == "go" {
			steps = append(steps, NextStep{
				Title:       "Install Go Dependencies",
				Description: "Download and install Go module dependencies",
				Command:     "go mod tidy",
				Optional:    false,
			})
			steps = append(steps, NextStep{
				Title:       "Build the Application",
				Description: "Compile the Go application",
				Command:     "go build ./cmd/server",
				Optional:    false,
			})
			steps = append(steps, NextStep{
				Title:       "Run the Server",
				Description: "Start the development server",
				Command:     "go run ./cmd/server",
				Optional:    true,
			})
		}

	case "frontend":
		if strings.Contains(template.Technology, "node") || strings.Contains(template.Name, "nextjs") {
			steps = append(steps, NextStep{
				Title:       "Install Node Dependencies",
				Description: "Install npm packages and dependencies",
				Command:     "npm install",
				Optional:    false,
			})
			steps = append(steps, NextStep{
				Title:       "Start Development Server",
				Description: "Start the development server with hot reload",
				Command:     "npm run dev",
				Optional:    true,
			})
		}

	case "mobile":
		steps = append(steps, NextStep{
			Title:       "Setup Mobile Development Environment",
			Description: "Configure your mobile development tools and SDKs",
			Optional:    false,
		})

	case "infrastructure":
		steps = append(steps, NextStep{
			Title:       "Review Infrastructure Configuration",
			Description: "Review and customize infrastructure settings before deployment",
			Optional:    false,
		})
	}

	return steps
}

// getGeneralNextSteps returns general next steps for any project
func (cs *CompletionSystem) getGeneralNextSteps(result *GenerationResult) []NextStep {
	return []NextStep{
		{
			Title:       "Navigate to Project Directory",
			Description: "Change to the generated project directory",
			Command:     fmt.Sprintf("cd %s", result.OutputDirectory),
			Optional:    false,
		},
		{
			Title:       "Initialize Git Repository",
			Description: "Initialize version control for your project",
			Command:     "git init && git add . && git commit -m \"Initial commit\"",
			Optional:    true,
		},
		{
			Title:       "Review Documentation",
			Description: "Read the generated README.md and documentation files",
			Optional:    true,
		},
		{
			Title:       "Customize Configuration",
			Description: "Review and customize configuration files for your needs",
			Optional:    true,
		},
		{
			Title:       "Set Up CI/CD",
			Description: "Configure continuous integration and deployment pipelines",
			Optional:    true,
		},
	}
}

// buildAdditionalInfo builds additional information for the summary
func (cs *CompletionSystem) buildAdditionalInfo(result *GenerationResult) []string {
	info := []string{}

	// Add generation statistics
	info = append(info, fmt.Sprintf("Generated %d files in %s", result.TotalFiles, result.Duration))
	info = append(info, fmt.Sprintf("Total project size: %s", cs.formatSize(result.TotalSize)))

	// Add template information
	templateNames := []string{}
	for _, template := range result.SelectedTemplates {
		if template.Selected {
			templateNames = append(templateNames, template.Template.DisplayName)
		}
	}
	if len(templateNames) > 0 {
		info = append(info, fmt.Sprintf("Templates used: %s", strings.Join(templateNames, ", ")))
	}

	// Add warnings if any
	if len(result.Warnings) > 0 {
		info = append(info, fmt.Sprintf("⚠️  %d warnings occurred during generation", len(result.Warnings)))
	}

	// Add error information if any
	if len(result.Errors) > 0 {
		info = append(info, fmt.Sprintf("❌ %d errors occurred during generation", len(result.Errors)))
	}

	// Add project-specific information
	if result.ProjectConfig.License != "" {
		info = append(info, fmt.Sprintf("License: %s", result.ProjectConfig.License))
	}

	if result.ProjectConfig.Author != "" {
		info = append(info, fmt.Sprintf("Author: %s", result.ProjectConfig.Author))
	}

	return info
}

// offerInteractiveHelp offers interactive help options after completion
func (cs *CompletionSystem) offerInteractiveHelp(ctx context.Context, result *GenerationResult) error {
	confirmConfig := interfaces.ConfirmConfig{
		Prompt:       "Interactive Help",
		Description:  "Would you like to see detailed help for getting started with your project?",
		DefaultValue: false,
		AllowBack:    false,
		AllowQuit:    false,
		ShowHelp:     false,
	}

	confirmResult, err := cs.ui.PromptConfirm(ctx, confirmConfig)
	if err != nil {
		return fmt.Errorf("failed to get help confirmation: %w", err)
	}

	if !confirmResult.Confirmed || confirmResult.Action != "confirm" {
		return nil
	}

	// Show detailed help menu
	return cs.showDetailedHelpMenu(ctx, result)
}

// showDetailedHelpMenu shows a menu of detailed help options
func (cs *CompletionSystem) showDetailedHelpMenu(ctx context.Context, result *GenerationResult) error {
	menuConfig := interfaces.MenuConfig{
		Title:       "Project Help Menu",
		Description: "Choose a topic to get detailed help and guidance",
		Options: []interfaces.MenuOption{
			{
				Label:       "🚀 Getting Started Guide",
				Description: "Step-by-step guide to get your project running",
				Value:       "getting_started",
				Icon:        "🚀",
			},
			{
				Label:       "📁 Project Structure",
				Description: "Understand the generated project structure",
				Value:       "project_structure",
				Icon:        "📁",
			},
			{
				Label:       "⚙️ Configuration Guide",
				Description: "Learn how to configure your project",
				Value:       "configuration",
				Icon:        "⚙️",
			},
			{
				Label:       "🔧 Development Workflow",
				Description: "Best practices for development and testing",
				Value:       "development",
				Icon:        "🔧",
			},
			{
				Label:       "🚀 Deployment Guide",
				Description: "How to deploy your project to production",
				Value:       "deployment",
				Icon:        "🚀",
			},
			{
				Label:       "❓ Troubleshooting",
				Description: "Common issues and solutions",
				Value:       "troubleshooting",
				Icon:        "❓",
			},
			{
				Label:       "✅ Done",
				Description: "Exit help menu",
				Value:       "done",
				Icon:        "✅",
			},
		},
		AllowBack: false,
		AllowQuit: true,
		ShowHelp:  false,
	}

	for {
		menuResult, err := cs.ui.ShowMenu(ctx, menuConfig)
		if err != nil {
			return fmt.Errorf("failed to show help menu: %w", err)
		}

		if menuResult.Cancelled || menuResult.SelectedValue == "done" {
			break
		}

		if err := cs.showSpecificHelp(ctx, menuResult.SelectedValue.(string), result); err != nil {
			cs.logger.WarnWithFields("Failed to show specific help", map[string]interface{}{
				"topic": menuResult.SelectedValue,
				"error": err.Error(),
			})
		}
	}

	return nil
}

// showSpecificHelp shows help for a specific topic
func (cs *CompletionSystem) showSpecificHelp(ctx context.Context, topic string, result *GenerationResult) error {
	switch topic {
	case "getting_started":
		return cs.helpSystem.ShowContextHelp(ctx, "getting_started", cs.buildGettingStartedHelp(result))
	case "project_structure":
		return cs.helpSystem.ShowContextHelp(ctx, "project_structure", cs.buildProjectStructureHelp(result))
	case "configuration":
		return cs.helpSystem.ShowContextHelp(ctx, "configuration", cs.buildConfigurationHelp(result))
	case "development":
		return cs.helpSystem.ShowContextHelp(ctx, "development", cs.buildDevelopmentHelp(result))
	case "deployment":
		return cs.helpSystem.ShowContextHelp(ctx, "deployment", cs.buildDeploymentHelp(result))
	case "troubleshooting":
		return cs.helpSystem.ShowContextHelp(ctx, "troubleshooting", cs.buildTroubleshootingHelp(result))
	default:
		return cs.helpSystem.ShowContextHelp(ctx, "general")
	}
}

// Helper methods for building help content

func (cs *CompletionSystem) buildGettingStartedHelp(result *GenerationResult) string {
	return fmt.Sprintf(`Getting Started with %s:

1. Navigate to your project directory:
   cd %s

2. Follow the next steps shown in the completion summary
3. Read the README.md file for project-specific instructions
4. Check the documentation in the Docs/ directory

Your project is ready to use!`, result.ProjectConfig.Name, result.OutputDirectory)
}

func (cs *CompletionSystem) buildProjectStructureHelp(result *GenerationResult) string {
	return `Project Structure Overview:

The generated project follows a standardized structure:
• App/ - Frontend applications and components
• CommonServer/ - Backend API and server code
• Mobile/ - Mobile applications (iOS/Android)
• Deploy/ - Infrastructure and deployment configurations
• Docs/ - Project documentation
• Scripts/ - Build and utility scripts
• .github/ - GitHub workflows and templates

Each directory contains specific components based on your template selections.`
}

func (cs *CompletionSystem) buildConfigurationHelp(result *GenerationResult) string {
	return `Configuration Guide:

Key configuration files in your project:
• Environment files (.env, .env.example) - Environment variables
• Package files (package.json, go.mod) - Dependencies
• Config files (config.yaml, settings.json) - Application settings
• Docker files (Dockerfile, docker-compose.yml) - Containerization
• CI/CD files (.github/workflows/) - Automation

Review and customize these files according to your needs.`
}

func (cs *CompletionSystem) buildDevelopmentHelp(result *GenerationResult) string {
	return `Development Workflow:

1. Set up your development environment
2. Install dependencies (npm install, go mod tidy, etc.)
3. Configure environment variables
4. Run tests to ensure everything works
5. Start development servers
6. Make changes and test iteratively
7. Use version control (git) for tracking changes

Check the Scripts/ directory for helpful development scripts.`
}

func (cs *CompletionSystem) buildDeploymentHelp(result *GenerationResult) string {
	return `Deployment Guide:

Your project includes deployment configurations:
• Docker configurations for containerization
• Kubernetes manifests for orchestration
• Terraform scripts for infrastructure
• CI/CD workflows for automation

Review the Deploy/ directory for deployment options and customize according to your target environment.`
}

func (cs *CompletionSystem) buildTroubleshootingHelp(result *GenerationResult) string {
	troubleshooting := `Common Issues and Solutions:

1. Dependencies not installing:
   - Check your internet connection
   - Verify package manager is installed
   - Clear package manager cache

2. Build failures:
   - Check for syntax errors
   - Verify all dependencies are installed
   - Check environment variables

3. Server not starting:
   - Check port availability
   - Verify configuration files
   - Check logs for error messages`

	if len(result.Errors) > 0 {
		troubleshooting += "\n\nGeneration Errors:\n"
		for _, err := range result.Errors {
			troubleshooting += fmt.Sprintf("• %s\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		troubleshooting += "\n\nGeneration Warnings:\n"
		for _, warning := range result.Warnings {
			troubleshooting += fmt.Sprintf("• %s\n", warning)
		}
	}

	return troubleshooting
}

// Utility methods

func (cs *CompletionSystem) categorizeFile(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// Categorize by directory
	if strings.Contains(dir, "App") || strings.Contains(dir, "frontend") {
		return "Frontend"
	}
	if strings.Contains(dir, "CommonServer") || strings.Contains(dir, "backend") {
		return "Backend"
	}
	if strings.Contains(dir, "Mobile") || strings.Contains(dir, "mobile") {
		return "Mobile"
	}
	if strings.Contains(dir, "Deploy") || strings.Contains(dir, "infrastructure") {
		return "Infrastructure"
	}
	if strings.Contains(dir, "Docs") || strings.Contains(dir, "docs") {
		return "Documentation"
	}
	if strings.Contains(dir, "Scripts") || strings.Contains(dir, "scripts") {
		return "Scripts"
	}
	if strings.Contains(dir, ".github") {
		return "CI/CD"
	}

	// Categorize by file extension
	ext := filepath.Ext(base)
	switch ext {
	case ".go":
		return "Backend"
	case ".js", ".ts", ".jsx", ".tsx", ".vue", ".svelte":
		return "Frontend"
	case ".swift", ".kt", ".java":
		return "Mobile"
	case ".yml", ".yaml", ".json", ".toml":
		return "Configuration"
	case ".md", ".txt", ".rst":
		return "Documentation"
	case ".sh", ".bat", ".ps1":
		return "Scripts"
	case ".dockerfile", "Dockerfile":
		return "Infrastructure"
	default:
		return "Other"
	}
}

func (cs *CompletionSystem) getCategoryPath(category string) string {
	switch category {
	case "Frontend":
		return "App"
	case "Backend":
		return "CommonServer"
	case "Mobile":
		return "Mobile"
	case "Infrastructure":
		return "Deploy"
	case "Documentation":
		return "Docs"
	case "Scripts":
		return "Scripts"
	case "CI/CD":
		return ".github"
	default:
		return ""
	}
}

func (cs *CompletionSystem) formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
