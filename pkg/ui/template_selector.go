// Package ui provides template selection components for interactive CLI generation.
package ui

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TemplateSelector provides interactive template selection functionality
type TemplateSelector struct {
	ui              interfaces.InteractiveUIInterface
	templateManager interfaces.TemplateManager
	logger          interfaces.Logger
}

// TemplateCategory represents a category of templates
type TemplateCategory struct {
	Name        string
	DisplayName string
	Description string
	Icon        string
	Templates   []interfaces.TemplateInfo
}

// TemplateSelection is an alias for the interfaces version
type TemplateSelection = interfaces.TemplateSelection

// NewTemplateSelector creates a new template selector instance
func NewTemplateSelector(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *TemplateSelector {
	return &TemplateSelector{
		ui:              ui,
		templateManager: templateManager,
		logger:          logger,
	}
}

// SelectTemplatesInteractively provides an interactive template selection interface
func (ts *TemplateSelector) SelectTemplatesInteractively(ctx context.Context) ([]TemplateSelection, error) {
	// Get all available templates
	allTemplates, err := ts.templateManager.ListTemplates(interfaces.TemplateFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	if len(allTemplates) == 0 {
		return nil, fmt.Errorf("no templates available")
	}

	// Organize templates by category
	categories := ts.organizeTemplatesByCategory(allTemplates)

	// Show category selection first
	selectedCategory, err := ts.selectCategory(ctx, categories)
	if err != nil {
		return nil, fmt.Errorf("failed to select category: %w", err)
	}

	if selectedCategory == nil {
		return []TemplateSelection{}, nil // User cancelled
	}

	// Show templates in selected category
	selections, err := ts.selectTemplatesInCategory(ctx, *selectedCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to select templates: %w", err)
	}

	return selections, nil
}

// organizeTemplatesByCategory groups templates by their category
func (ts *TemplateSelector) organizeTemplatesByCategory(templates []interfaces.TemplateInfo) []TemplateCategory {
	categoryMap := make(map[string][]interfaces.TemplateInfo)

	// Group templates by category
	for _, template := range templates {
		category := template.Category
		if category == "" {
			category = "other"
		}
		categoryMap[category] = append(categoryMap[category], template)
	}

	// Convert to structured categories
	var categories []TemplateCategory

	// Define category order and metadata
	categoryOrder := []struct {
		key         string
		displayName string
		description string
		icon        string
	}{
		{"frontend", "Frontend", "Web applications and user interfaces", "🌐"},
		{"backend", "Backend", "Server-side applications and APIs", "⚙️"},
		{"mobile", "Mobile", "Mobile applications for iOS and Android", "📱"},
		{"infrastructure", "Infrastructure", "Deployment and infrastructure as code", "🏗️"},
		{"base", "Base", "Core project files and documentation", "📄"},
		{"other", "Other", "Miscellaneous templates", "📦"},
	}

	for _, catInfo := range categoryOrder {
		if templates, exists := categoryMap[catInfo.key]; exists {
			// Sort templates within category by name
			sort.Slice(templates, func(i, j int) bool {
				return templates[i].Name < templates[j].Name
			})

			categories = append(categories, TemplateCategory{
				Name:        catInfo.key,
				DisplayName: catInfo.displayName,
				Description: catInfo.description,
				Icon:        catInfo.icon,
				Templates:   templates,
			})
		}
	}

	return categories
}

// selectCategory shows category selection menu
func (ts *TemplateSelector) selectCategory(ctx context.Context, categories []TemplateCategory) (*TemplateCategory, error) {
	if len(categories) == 0 {
		return nil, fmt.Errorf("no template categories available")
	}

	// If only one category, select it automatically
	if len(categories) == 1 {
		return &categories[0], nil
	}

	// Build menu options
	var options []interfaces.MenuOption
	for i, category := range categories {
		description := fmt.Sprintf("%s (%d templates)", category.Description, len(category.Templates))

		options = append(options, interfaces.MenuOption{
			Label:       category.DisplayName,
			Description: description,
			Value:       i,
			Icon:        category.Icon,
		})
	}

	// Add option to select from all categories
	options = append(options, interfaces.MenuOption{
		Label:       "All Categories",
		Description: "Browse templates from all categories",
		Value:       -1,
		Icon:        "📋",
	})

	menuConfig := interfaces.MenuConfig{
		Title:       "Select Template Category",
		Description: "Choose a category to browse available templates",
		Options:     options,
		AllowBack:   false,
		AllowQuit:   true,
		ShowHelp:    true,
		HelpText:    "Select a template category to see available templates. Each category contains templates for specific types of projects.",
	}

	result, err := ts.ui.ShowMenu(ctx, menuConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to show category menu: %w", err)
	}

	if result.Cancelled {
		return nil, nil
	}

	categoryIndex, ok := result.SelectedValue.(int)
	if !ok {
		return nil, fmt.Errorf("invalid category selection")
	}

	// Handle "All Categories" selection
	if categoryIndex == -1 {
		// Create a combined category with all templates
		var allTemplates []interfaces.TemplateInfo
		for _, category := range categories {
			allTemplates = append(allTemplates, category.Templates...)
		}

		// Sort all templates by name
		sort.Slice(allTemplates, func(i, j int) bool {
			return allTemplates[i].Name < allTemplates[j].Name
		})

		return &TemplateCategory{
			Name:        "all",
			DisplayName: "All Templates",
			Description: "All available templates",
			Icon:        "📋",
			Templates:   allTemplates,
		}, nil
	}

	if categoryIndex < 0 || categoryIndex >= len(categories) {
		return nil, fmt.Errorf("invalid category index: %d", categoryIndex)
	}

	return &categories[categoryIndex], nil
}

// selectTemplatesInCategory shows template selection within a category
func (ts *TemplateSelector) selectTemplatesInCategory(ctx context.Context, category TemplateCategory) ([]TemplateSelection, error) {
	if len(category.Templates) == 0 {
		return nil, fmt.Errorf("no templates available in category: %s", category.DisplayName)
	}

	// Build multi-select options
	var options []interfaces.SelectOption
	for _, template := range category.Templates {
		// Create detailed description
		description := template.Description
		if template.Technology != "" {
			description = fmt.Sprintf("%s (Technology: %s)", description, template.Technology)
		}
		if len(template.Tags) > 0 {
			description = fmt.Sprintf("%s | Tags: %s", description, strings.Join(template.Tags, ", "))
		}

		// Determine icon based on technology or category
		icon := ts.getTemplateIcon(template)

		options = append(options, interfaces.SelectOption{
			Label:       template.DisplayName,
			Description: description,
			Value:       template.Name,
			Icon:        icon,
			Category:    template.Category,
			Tags:        template.Tags,
			Metadata: map[string]interface{}{
				"template": template,
			},
		})
	}

	multiSelectConfig := interfaces.MultiSelectConfig{
		Title:         fmt.Sprintf("Select Templates - %s", category.DisplayName),
		Description:   fmt.Sprintf("Choose one or more templates from the %s category", category.DisplayName),
		Options:       options,
		MinSelection:  1,
		MaxSelection:  0, // No limit
		AllowBack:     true,
		AllowQuit:     true,
		ShowHelp:      true,
		SearchEnabled: true,
		HelpText:      "Select templates to include in your project. You can select multiple templates that will be combined. Use search (/) to filter templates.",
	}

	result, err := ts.ui.ShowMultiSelect(ctx, multiSelectConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to show template selection: %w", err)
	}

	if result.Cancelled {
		return []TemplateSelection{}, nil
	}

	// Convert results to TemplateSelection
	var selections []TemplateSelection
	for i, selectedIndex := range result.SelectedIndices {
		if selectedIndex >= 0 && selectedIndex < len(category.Templates) {
			template := category.Templates[selectedIndex]

			selections = append(selections, TemplateSelection{
				Template: template,
				Selected: true,
				Options:  make(map[string]interface{}),
			})

			if ts.logger != nil {
				ts.logger.InfoWithFields("Template selected", map[string]interface{}{
					"template":  template.Name,
					"category":  template.Category,
					"selection": i + 1,
				})
			}
		}
	}

	return selections, nil
}

// getTemplateIcon returns an appropriate icon for a template based on its technology
func (ts *TemplateSelector) getTemplateIcon(template interfaces.TemplateInfo) string {
	technology := strings.ToLower(template.Technology)
	category := strings.ToLower(template.Category)

	// Technology-specific icons
	switch {
	case strings.Contains(technology, "go"):
		return "🐹"
	case strings.Contains(technology, "next") || strings.Contains(technology, "react"):
		return "⚛️"
	case strings.Contains(technology, "node"):
		return "🟢"
	case strings.Contains(technology, "python"):
		return "🐍"
	case strings.Contains(technology, "java"):
		return "☕"
	case strings.Contains(technology, "kotlin"):
		return "🤖"
	case strings.Contains(technology, "swift"):
		return "🍎"
	case strings.Contains(technology, "docker"):
		return "🐳"
	case strings.Contains(technology, "terraform"):
		return "🏗️"
	case strings.Contains(technology, "kubernetes"):
		return "☸️"
	}

	// Category-specific icons as fallback
	switch category {
	case "frontend":
		return "🌐"
	case "backend":
		return "⚙️"
	case "mobile":
		return "📱"
	case "infrastructure":
		return "🏗️"
	case "base":
		return "📄"
	default:
		return "📦"
	}
}

// ShowTemplateInfo displays detailed information about a template
func (ts *TemplateSelector) ShowTemplateInfo(ctx context.Context, template interfaces.TemplateInfo) error {
	// Get additional template metadata
	metadata, err := ts.templateManager.GetTemplateMetadata(template.Name)
	if err != nil {
		ts.logger.WarnWithFields("Failed to get template metadata", map[string]interface{}{
			"template": template.Name,
			"error":    err.Error(),
		})
	}

	// Build information display
	var info strings.Builder
	info.WriteString(fmt.Sprintf("Template: %s\n", template.DisplayName))
	info.WriteString(fmt.Sprintf("Name: %s\n", template.Name))
	info.WriteString(fmt.Sprintf("Category: %s\n", template.Category))
	info.WriteString(fmt.Sprintf("Version: %s\n", template.Version))

	if template.Technology != "" {
		info.WriteString(fmt.Sprintf("Technology: %s\n", template.Technology))
	}

	if template.Description != "" {
		info.WriteString(fmt.Sprintf("Description: %s\n", template.Description))
	}

	if len(template.Tags) > 0 {
		info.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(template.Tags, ", ")))
	}

	if len(template.Dependencies) > 0 {
		info.WriteString(fmt.Sprintf("Dependencies: %s\n", strings.Join(template.Dependencies, ", ")))
	}

	if metadata != nil {
		if metadata.Author != "" {
			info.WriteString(fmt.Sprintf("Author: %s\n", metadata.Author))
		}
		if metadata.License != "" {
			info.WriteString(fmt.Sprintf("License: %s\n", metadata.License))
		}
		if len(metadata.Keywords) > 0 {
			info.WriteString(fmt.Sprintf("Keywords: %s\n", strings.Join(metadata.Keywords, ", ")))
		}
	}

	// Show as a simple text display (we'll implement a proper info dialog later)
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("TEMPLATE INFORMATION")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Print(info.String())
	fmt.Println(strings.Repeat("=", 60))
	fmt.Print("Press Enter to continue...")

	// Simple input reading for now
	var input string
	_, _ = fmt.Scanln(&input)

	return nil
}

// ValidateTemplateSelections validates that selected templates are compatible
func (ts *TemplateSelector) ValidateTemplateSelections(selections []TemplateSelection) error {
	if len(selections) == 0 {
		return fmt.Errorf("no templates selected")
	}

	// Check for conflicts between selected templates
	categories := make(map[string]int)
	for _, selection := range selections {
		categories[selection.Template.Category]++
	}

	// Validate category combinations
	if categories["frontend"] > 1 {
		return fmt.Errorf("multiple frontend templates selected - only one frontend template is allowed")
	}

	if categories["backend"] > 1 {
		return fmt.Errorf("multiple backend templates selected - only one backend template is allowed")
	}

	// Check template dependencies
	for _, selection := range selections {
		for _, dependency := range selection.Template.Dependencies {
			found := false
			for _, other := range selections {
				if other.Template.Name == dependency {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("template '%s' requires dependency '%s' which is not selected",
					selection.Template.Name, dependency)
			}
		}
	}

	return nil
}
