package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/config/collectors"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/config/formatters"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/config/validators"
)

// InteractiveConfigurationManager handles interactive configuration operations
type InteractiveConfigurationManager struct {
	ui                    interfaces.InteractiveUIInterface
	persistence           *config.ConfigurationPersistence
	logger                interfaces.Logger
	inputCollector        *collectors.InputCollector
	confirmationCollector *collectors.ConfirmationCollector
	validator             *validators.ConfigValidator
	displayFormatter      *formatters.DisplayFormatter
	fileFormatter         *formatters.FileFormatter
	templateConverter     *formatters.TemplateConverter
}

// ConfigurationLoadOptions contains options for loading configurations
type ConfigurationLoadOptions struct {
	AllowModification bool
	ShowPreview       bool
	ConfirmLoad       bool
}

// LoadedConfiguration represents a loaded configuration with metadata
type LoadedConfiguration struct {
	Name               string                         `json:"name"`
	ProjectConfig      *models.ProjectConfig          `json:"project_config"`
	SelectedTemplates  []interfaces.TemplateSelection `json:"selected_templates"`
	GenerationSettings *config.GenerationSettings     `json:"generation_settings"`
	UserPreferences    *config.UserPreferences        `json:"user_preferences"`
	LoadedAt           time.Time                      `json:"loaded_at"`
	AllowModification  bool                           `json:"allow_modification"`
}

// NewInteractiveConfigurationManager creates a new interactive configuration manager
func NewInteractiveConfigurationManager(
	ui interfaces.InteractiveUIInterface,
	configDir string,
	logger interfaces.Logger,
) *InteractiveConfigurationManager {
	persistence := config.NewConfigurationPersistence(configDir, logger)

	return &InteractiveConfigurationManager{
		ui:                    ui,
		persistence:           persistence,
		logger:                logger,
		inputCollector:        collectors.NewInputCollector(ui),
		confirmationCollector: collectors.NewConfirmationCollector(ui),
		validator:             validators.NewConfigValidator(),
		displayFormatter:      formatters.NewDisplayFormatter(ui),
		fileFormatter:         formatters.NewFileFormatter(),
		templateConverter:     formatters.NewTemplateConverter(),
	}
}

// SaveConfigurationInteractively prompts the user to save a configuration
func (icm *InteractiveConfigurationManager) SaveConfigurationInteractively(
	ctx context.Context,
	projectConfig *models.ProjectConfig,
	selectedTemplates []interfaces.TemplateSelection,
) (string, error) {
	// Ask if user wants to save configuration
	confirmed, err := icm.confirmationCollector.ConfirmSave(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get save confirmation: %w", err)
	}

	if !confirmed {
		return "", nil // User chose not to save
	}

	// Get configuration details
	saveOptions, err := icm.inputCollector.CollectSaveOptions(ctx, icm.validator.ValidateConfigurationName)
	if err != nil {
		return "", fmt.Errorf("failed to collect save options: %w", err)
	}

	if saveOptions == nil {
		return "", nil // User cancelled
	}

	// Check if configuration already exists
	if icm.persistence.ConfigurationExists(saveOptions.Name) && !saveOptions.Overwrite {
		overwrite, err := icm.confirmationCollector.ConfirmOverwrite(ctx, saveOptions.Name)
		if err != nil {
			return "", fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !overwrite {
			return "", nil // User chose not to overwrite
		}
	}

	// Create saved configuration
	savedConfig := &config.SavedConfiguration{
		Name:              saveOptions.Name,
		Description:       saveOptions.Description,
		Tags:              saveOptions.Tags,
		ProjectConfig:     projectConfig,
		SelectedTemplates: icm.templateConverter.ConvertToSaved(selectedTemplates),
		GenerationSettings: &config.GenerationSettings{
			IncludeExamples:   true,
			IncludeTests:      true,
			IncludeDocs:       true,
			UpdateVersions:    false,
			MinimalMode:       false,
			BackupExisting:    true,
			OverwriteExisting: false,
		},
		UserPreferences: &config.UserPreferences{
			DefaultLicense:      projectConfig.License,
			DefaultAuthor:       projectConfig.Author,
			DefaultEmail:        projectConfig.Email,
			DefaultOrganization: projectConfig.Organization,
			PreferredFormat:     "yaml",
		},
	}

	// Save configuration
	if err := icm.persistence.SaveConfiguration(saveOptions.Name, savedConfig); err != nil {
		return "", fmt.Errorf("failed to save configuration: %w", err)
	}

	// Show success message
	icm.displayFormatter.ShowSaveSuccess(ctx, saveOptions.Name)

	return saveOptions.Name, nil
}

// LoadConfigurationInteractively prompts the user to load a saved configuration
func (icm *InteractiveConfigurationManager) LoadConfigurationInteractively(
	ctx context.Context,
	options *ConfigurationLoadOptions,
) (*LoadedConfiguration, error) {
	if options == nil {
		options = &ConfigurationLoadOptions{
			AllowModification: true,
			ShowPreview:       true,
			ConfirmLoad:       true,
		}
	}

	// List available configurations
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		icm.displayFormatter.ShowNoConfigurationsMessage(ctx)
		return nil, nil
	}

	// Show configuration selection menu
	selectedConfig, err := icm.inputCollector.SelectConfiguration(ctx, configs)
	if err != nil {
		return nil, fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil, nil // User cancelled selection
	}

	// Show configuration preview if requested
	if options.ShowPreview {
		if err := icm.displayFormatter.ShowConfigurationPreview(ctx, selectedConfig); err != nil {
			return nil, fmt.Errorf("failed to show configuration preview: %w", err)
		}
	}

	// Confirm loading if requested
	if options.ConfirmLoad {
		confirmed, err := icm.confirmationCollector.ConfirmLoad(ctx, selectedConfig.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm load: %w", err)
		}
		if !confirmed {
			return nil, nil // User chose not to load
		}
	}

	// Convert to loaded configuration
	loadedConfig := &LoadedConfiguration{
		Name:               selectedConfig.Name,
		ProjectConfig:      selectedConfig.ProjectConfig,
		SelectedTemplates:  icm.templateConverter.ConvertFromSaved(selectedConfig.SelectedTemplates),
		GenerationSettings: selectedConfig.GenerationSettings,
		UserPreferences:    selectedConfig.UserPreferences,
		LoadedAt:           time.Now(),
		AllowModification:  options.AllowModification,
	}

	icm.logger.InfoWithFields("Configuration loaded", map[string]interface{}{
		"config_name":  selectedConfig.Name,
		"project_name": selectedConfig.ProjectConfig.Name,
	})

	return loadedConfig, nil
}

// ManageConfigurationsInteractively provides an interactive configuration management interface
func (icm *InteractiveConfigurationManager) ManageConfigurationsInteractively(ctx context.Context) error {
	for {
		action, err := icm.inputCollector.ShowManagementMenu(ctx)
		if err != nil {
			return fmt.Errorf("failed to show management menu: %w", err)
		}

		if action == "quit" {
			break
		}

		switch action {
		case "list":
			if err := icm.listConfigurationsInteractively(ctx); err != nil {
				icm.displayFormatter.ShowError(ctx, "Failed to list configurations", err)
			}
		case "view":
			if err := icm.viewConfigurationInteractively(ctx); err != nil {
				icm.displayFormatter.ShowError(ctx, "Failed to view configuration", err)
			}
		case "delete":
			if err := icm.deleteConfigurationInteractively(ctx); err != nil {
				icm.displayFormatter.ShowError(ctx, "Failed to delete configuration", err)
			}
		case "export":
			if err := icm.exportConfigurationInteractively(ctx); err != nil {
				icm.displayFormatter.ShowError(ctx, "Failed to export configuration", err)
			}
		case "import":
			if err := icm.importConfigurationInteractively(ctx); err != nil {
				icm.displayFormatter.ShowError(ctx, "Failed to import configuration", err)
			}
		}
	}

	return nil
}

// Private helper methods

// listConfigurationsInteractively shows a list of all configurations
func (icm *InteractiveConfigurationManager) listConfigurationsInteractively(ctx context.Context) error {
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	return icm.displayFormatter.ShowConfigurationsList(ctx, configs)
}

// viewConfigurationInteractively shows detailed view of a configuration
func (icm *InteractiveConfigurationManager) viewConfigurationInteractively(ctx context.Context) error {
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		return icm.displayFormatter.ShowNoConfigurationsMessage(ctx)
	}

	selectedConfig, err := icm.inputCollector.SelectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	return icm.displayFormatter.ShowDetailedConfigurationView(ctx, selectedConfig)
}

// deleteConfigurationInteractively handles interactive configuration deletion
func (icm *InteractiveConfigurationManager) deleteConfigurationInteractively(ctx context.Context) error {
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		return icm.displayFormatter.ShowNoConfigurationsMessage(ctx)
	}

	selectedConfig, err := icm.inputCollector.SelectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	confirmed, err := icm.confirmationCollector.ConfirmDelete(ctx, selectedConfig.Name)
	if err != nil {
		return fmt.Errorf("failed to get deletion confirmation: %w", err)
	}

	if !confirmed {
		return nil // User chose not to delete
	}

	if err := icm.persistence.DeleteConfiguration(selectedConfig.Name); err != nil {
		return fmt.Errorf("failed to delete configuration: %w", err)
	}

	icm.displayFormatter.ShowDeleteSuccess(ctx, selectedConfig.Name)
	return nil
}

// exportConfigurationInteractively handles interactive configuration export
func (icm *InteractiveConfigurationManager) exportConfigurationInteractively(ctx context.Context) error {
	configs, err := icm.persistence.ListConfigurations(nil)
	if err != nil {
		return fmt.Errorf("failed to list configurations: %w", err)
	}

	if len(configs) == 0 {
		return icm.displayFormatter.ShowNoConfigurationsMessage(ctx)
	}

	selectedConfig, err := icm.inputCollector.SelectConfiguration(ctx, configs)
	if err != nil {
		return fmt.Errorf("failed to select configuration: %w", err)
	}

	if selectedConfig == nil {
		return nil // User cancelled
	}

	format, err := icm.inputCollector.CollectExportFormat(ctx)
	if err != nil {
		return fmt.Errorf("failed to select export format: %w", err)
	}

	if format == "" {
		return nil // User cancelled
	}

	defaultFilename := icm.fileFormatter.GenerateDefaultFilename(selectedConfig.Name, format)
	exportPath, err := icm.inputCollector.CollectExportPath(ctx, defaultFilename)
	if err != nil {
		return fmt.Errorf("failed to get export path: %w", err)
	}

	if exportPath == "" {
		return nil // User cancelled
	}

	data, err := icm.persistence.ExportConfiguration(selectedConfig.Name, format)
	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}

	if err := icm.fileFormatter.WriteExportFile(exportPath, data); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	icm.displayFormatter.ShowExportSuccess(ctx, selectedConfig.Name, exportPath)
	return nil
}

// importConfigurationInteractively handles interactive configuration import
func (icm *InteractiveConfigurationManager) importConfigurationInteractively(ctx context.Context) error {
	importPath, err := icm.inputCollector.CollectImportPath(ctx)
	if err != nil {
		return fmt.Errorf("failed to get import path: %w", err)
	}

	if importPath == "" {
		return nil // User cancelled
	}

	format := icm.fileFormatter.DetermineFormatFromPath(importPath)
	data, err := icm.fileFormatter.ReadImportFile(importPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	configName, err := icm.inputCollector.CollectConfigurationName(ctx, icm.validator.ValidateConfigurationName)
	if err != nil {
		return fmt.Errorf("failed to get configuration name: %w", err)
	}

	if configName == "" {
		return nil // User cancelled
	}

	if icm.persistence.ConfigurationExists(configName) {
		overwrite, err := icm.confirmationCollector.ConfirmOverwrite(ctx, configName)
		if err != nil {
			return fmt.Errorf("failed to confirm overwrite: %w", err)
		}
		if !overwrite {
			return nil // User chose not to overwrite
		}
	}

	if err := icm.persistence.ImportConfiguration(configName, data, format); err != nil {
		return fmt.Errorf("failed to import configuration: %w", err)
	}

	icm.displayFormatter.ShowImportSuccess(ctx, configName)
	return nil
}
