package formatters

import (
	"context"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// DisplayFormatter handles formatting and displaying configuration data
type DisplayFormatter struct {
	ui interfaces.InteractiveUIInterface
}

// NewDisplayFormatter creates a new display formatter
func NewDisplayFormatter(ui interfaces.InteractiveUIInterface) *DisplayFormatter {
	return &DisplayFormatter{
		ui: ui,
	}
}

// ShowConfigurationPreview displays a preview of the selected configuration
func (df *DisplayFormatter) ShowConfigurationPreview(ctx context.Context, config *config.SavedConfiguration) error {
	// Create preview table
	headers := []string{"Setting", "Value"}
	rows := [][]string{
		{"Configuration Name", config.Name},
		{"Project Name", config.ProjectConfig.Name},
		{"Organization", config.ProjectConfig.Organization},
		{"Author", config.ProjectConfig.Author},
		{"License", config.ProjectConfig.License},
		{"Created", config.CreatedAt.Format("2006-01-02 15:04:05")},
		{"Updated", config.UpdatedAt.Format("2006-01-02 15:04:05")},
	}

	if config.Description != "" {
		rows = append(rows, []string{"Description", config.Description})
	}

	if len(config.Tags) > 0 {
		rows = append(rows, []string{"Tags", strings.Join(config.Tags, ", ")})
	}

	// Add template information
	templateNames := make([]string, 0, len(config.SelectedTemplates))
	for _, template := range config.SelectedTemplates {
		if template.Selected {
			templateNames = append(templateNames, template.TemplateName)
		}
	}
	if len(templateNames) > 0 {
		rows = append(rows, []string{"Templates", strings.Join(templateNames, ", ")})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Configuration Preview",
		Headers: headers,
		Rows:    rows,
	}

	return df.ui.ShowTable(ctx, tableConfig)
}

// ShowConfigurationsList displays a list of all configurations
func (df *DisplayFormatter) ShowConfigurationsList(ctx context.Context, configs []*config.SavedConfiguration) error {
	if len(configs) == 0 {
		return df.ShowNoConfigurationsMessage(ctx)
	}

	// Create table with configuration information
	headers := []string{"Name", "Project", "Templates", "Created", "Updated"}
	rows := make([][]string, 0, len(configs))

	for _, config := range configs {
		templateCount := 0
		for _, template := range config.SelectedTemplates {
			if template.Selected {
				templateCount++
			}
		}

		rows = append(rows, []string{
			config.Name,
			config.ProjectConfig.Name,
			fmt.Sprintf("%d", templateCount),
			config.CreatedAt.Format("2006-01-02"),
			config.UpdatedAt.Format("2006-01-02"),
		})
	}

	tableConfig := interfaces.TableConfig{
		Title:      "Saved Configurations",
		Headers:    headers,
		Rows:       rows,
		Pagination: true,
		PageSize:   10,
		Sortable:   true,
	}

	return df.ui.ShowTable(ctx, tableConfig)
}

// ShowDetailedConfigurationView shows a detailed view of a configuration
func (df *DisplayFormatter) ShowDetailedConfigurationView(ctx context.Context, config *config.SavedConfiguration) error {
	// Show basic information
	if err := df.ShowConfigurationPreview(ctx, config); err != nil {
		return fmt.Errorf("failed to show configuration preview: %w", err)
	}

	// Show template details
	if err := df.ShowTemplateDetails(ctx, config); err != nil {
		return fmt.Errorf("failed to show template details: %w", err)
	}

	// Show generation settings if available
	if config.GenerationSettings != nil {
		if err := df.ShowGenerationSettings(ctx, config.GenerationSettings); err != nil {
			return fmt.Errorf("failed to show generation settings: %w", err)
		}
	}

	return nil
}

// ShowTemplateDetails shows detailed template information
func (df *DisplayFormatter) ShowTemplateDetails(ctx context.Context, config *config.SavedConfiguration) error {
	if len(config.SelectedTemplates) == 0 {
		return nil
	}

	headers := []string{"Template", "Category", "Technology", "Version", "Selected"}
	rows := make([][]string, 0, len(config.SelectedTemplates))

	for _, template := range config.SelectedTemplates {
		selected := "No"
		if template.Selected {
			selected = "Yes"
		}

		rows = append(rows, []string{
			template.TemplateName,
			template.Category,
			template.Technology,
			template.Version,
			selected,
		})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Template Details",
		Headers: headers,
		Rows:    rows,
	}

	return df.ui.ShowTable(ctx, tableConfig)
}

// ShowGenerationSettings shows generation settings information
func (df *DisplayFormatter) ShowGenerationSettings(ctx context.Context, settings *config.GenerationSettings) error {
	headers := []string{"Setting", "Value"}
	rows := [][]string{
		{"Include Examples", fmt.Sprintf("%t", settings.IncludeExamples)},
		{"Include Tests", fmt.Sprintf("%t", settings.IncludeTests)},
		{"Include Documentation", fmt.Sprintf("%t", settings.IncludeDocs)},
		{"Update Versions", fmt.Sprintf("%t", settings.UpdateVersions)},
		{"Minimal Mode", fmt.Sprintf("%t", settings.MinimalMode)},
		{"Backup Existing", fmt.Sprintf("%t", settings.BackupExisting)},
		{"Overwrite Existing", fmt.Sprintf("%t", settings.OverwriteExisting)},
	}

	if len(settings.ExcludePatterns) > 0 {
		rows = append(rows, []string{"Exclude Patterns", strings.Join(settings.ExcludePatterns, ", ")})
	}

	if len(settings.IncludeOnlyPaths) > 0 {
		rows = append(rows, []string{"Include Only Paths", strings.Join(settings.IncludeOnlyPaths, ", ")})
	}

	tableConfig := interfaces.TableConfig{
		Title:   "Generation Settings",
		Headers: headers,
		Rows:    rows,
	}

	return df.ui.ShowTable(ctx, tableConfig)
}

// ShowSaveSuccess displays a success message for saving configuration
func (df *DisplayFormatter) ShowSaveSuccess(ctx context.Context, name string) {
	// This would show a success message - placeholder implementation
	fmt.Printf("Configuration '%s' saved successfully!\n", name)
}

// ShowNoConfigurationsMessage displays a message when no configurations are found
func (df *DisplayFormatter) ShowNoConfigurationsMessage(ctx context.Context) error {
	// This would show a message about no configurations - placeholder implementation
	fmt.Println("No saved configurations found.")
	return nil
}

// ShowError displays an error message
func (df *DisplayFormatter) ShowError(ctx context.Context, title string, err error) {
	// This would show an error message - placeholder implementation
	fmt.Printf("Error: %s - %v\n", title, err)
}

// ShowDeleteSuccess displays a success message for deleting configuration
func (df *DisplayFormatter) ShowDeleteSuccess(ctx context.Context, name string) {
	fmt.Printf("Configuration '%s' deleted successfully.\n", name)
}

// ShowExportSuccess displays a success message for exporting configuration
func (df *DisplayFormatter) ShowExportSuccess(ctx context.Context, name, path string) {
	fmt.Printf("Configuration '%s' exported to '%s' successfully.\n", name, path)
}

// ShowImportSuccess displays a success message for importing configuration
func (df *DisplayFormatter) ShowImportSuccess(ctx context.Context, name string) {
	fmt.Printf("Configuration imported as '%s' successfully.\n", name)
}
