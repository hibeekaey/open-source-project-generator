// Package template provides template management functionality for the
// Open Source Project Generator.
package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template/metadata"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template/processor"
)

// Manager implements the TemplateManager interface for template operations.
// It coordinates between template discovery, caching, validation, and processing components.
type Manager struct {
	// Core components
	discovery        *TemplateDiscovery
	cache            *TemplateCache
	validator        *TemplateValidator
	processingEngine *processor.ProcessingEngine

	// Legacy template engine for backward compatibility
	templateEngine interfaces.TemplateEngine

	// Metadata components
	metadataParser    *metadata.MetadataParser
	metadataValidator *metadata.MetadataValidator
}

// NewManager creates a new template manager instance with coordinated components.
func NewManager(templateEngine interfaces.TemplateEngine) interfaces.TemplateManager {
	// Create core components
	discovery := NewTemplateDiscovery(embeddedTemplates)
	validator := NewTemplateValidator()
	processingEngine := processor.NewProcessingEngine()

	// Create metadata components
	metadataParser := metadata.NewMetadataParser(embeddedTemplates)
	metadataValidator := metadata.NewMetadataValidator(metadataParser)

	manager := &Manager{
		discovery:         discovery,
		validator:         validator,
		processingEngine:  processingEngine,
		templateEngine:    templateEngine,
		metadataParser:    metadataParser,
		metadataValidator: metadataValidator,
	}

	// Update validator with manager reference for enhanced validation
	manager.validator = NewTemplateValidatorWithManager(manager)

	// Create template cache with discovery function
	manager.cache = NewTemplateCache(nil, manager.discovery.DiscoverTemplates)

	return manager
}

// ListTemplates lists available templates with optional filtering
func (m *Manager) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	// Get templates from cache or discover them
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("🚫 couldn't find available templates: %w", err)
	}

	// Apply filters using discovery component
	filtered := m.discovery.FilterTemplates(templates, filter)

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(filtered))
	for i, tmpl := range filtered {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// GetTemplateInfo retrieves detailed information about a specific template
func (m *Manager) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	// Get templates from cache or discover them
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Use discovery component to find template by name
	tmpl := m.discovery.GetTemplateByName(templates, name)
	if tmpl == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	result := m.convertToInterfaceTemplateInfo(tmpl)
	return &result, nil
}

// SearchTemplates searches for templates by query string
func (m *Manager) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	// Get templates from cache or discover them
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Use discovery component to search templates
	matches := m.discovery.SearchTemplates(templates, query)

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(matches))
	for i, tmpl := range matches {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// ValidateTemplate validates a template structure and metadata
func (m *Manager) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Use the validator component to validate the template
	return m.validator.ValidateTemplate(path)
}

// ValidateTemplateStructure validates template structure
func (m *Manager) ValidateTemplateStructure(template *interfaces.TemplateInfo) error {
	// Use the validator component to validate template structure
	return m.validator.ValidateTemplateStructure(template)
}

// ProcessTemplate processes a template with the given configuration
func (m *Manager) ProcessTemplate(templateName string, config *models.ProjectConfig, outputPath string) error {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, templateName)
	if templateInfo == nil {
		return fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", templateName)
	}

	// Create a copy of config for template-specific processing
	processConfig := *config

	// For Android templates, add package-specific fields while keeping originals
	if templateName == "android-kotlin" {
		// Add Android-specific package naming fields
		processConfig.AndroidPackageOrg = strings.ToLower(config.Organization)
		processConfig.AndroidPackageName = strings.ToLower(config.Name)
	}

	// For embedded templates, use the embedded template engine (backward compatibility)
	if templateInfo.Source == "embedded" {
		templatePath := fmt.Sprintf("templates/%s/%s", templateInfo.Category, templateName)
		return m.templateEngine.ProcessDirectory(templatePath, outputPath, &processConfig)
	}

	// For file-based templates, use the processing engine
	return m.processingEngine.ProcessDirectory(templateInfo.Path, outputPath, &processConfig)
}

// ProcessCustomTemplate processes a custom template from a path
func (m *Manager) ProcessCustomTemplate(templatePath string, config *models.ProjectConfig, outputPath string) error {
	// Validate template first using validator component
	validationResult, err := m.validator.ValidateTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to validate template: %w", err)
	}

	if !validationResult.Valid {
		return fmt.Errorf("template validation failed: %d errors found", len(validationResult.Issues))
	}

	// Process the template directory using processing engine
	return m.processingEngine.ProcessDirectory(templatePath, outputPath, config)
}

// GetTemplateMetadata retrieves template metadata
func (m *Manager) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, name)
	if templateInfo == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	// Convert to interface metadata
	metadata := &interfaces.TemplateMetadata{
		Author:     templateInfo.Metadata.Author,
		License:    templateInfo.Metadata.License,
		Repository: templateInfo.Metadata.Repository,
		Homepage:   templateInfo.Metadata.Homepage,
		Keywords:   templateInfo.Metadata.Keywords,
		Created:    templateInfo.Metadata.CreatedAt,
		Updated:    templateInfo.Metadata.UpdatedAt,
	}

	return metadata, nil
}

// GetTemplateDependencies retrieves template dependencies
func (m *Manager) GetTemplateDependencies(name string) ([]string, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, name)
	if templateInfo == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	return templateInfo.Dependencies, nil
}

// GetTemplateCompatibility retrieves template compatibility information
func (m *Manager) GetTemplateCompatibility(name string) (*interfaces.CompatibilityInfo, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, name)
	if templateInfo == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	// Create compatibility info based on template metadata
	compatibility := &interfaces.CompatibilityInfo{
		MinGeneratorVersion: templateInfo.Metadata.MinVersion,
		MaxGeneratorVersion: templateInfo.Metadata.MaxVersion,
		SupportedPlatforms:  []string{"linux", "darwin", "windows"}, // Default platforms
		RequiredFeatures:    []string{},
		Dependencies:        []interfaces.TemplateDependency{},
		Conflicts:           []string{},
	}

	// Convert dependencies to TemplateDependency format
	for _, dep := range templateInfo.Dependencies {
		compatibility.Dependencies = append(compatibility.Dependencies, interfaces.TemplateDependency{
			Name:        dep,
			Version:     "latest",
			Type:        "template",
			Required:    true,
			Description: fmt.Sprintf("Required dependency: %s", dep),
		})
	}

	// Set required features based on template technology
	switch strings.ToLower(templateInfo.Technology) {
	case "go", "golang":
		compatibility.RequiredFeatures = append(compatibility.RequiredFeatures, "go-runtime")
	case "node.js", "nextjs", "react":
		compatibility.RequiredFeatures = append(compatibility.RequiredFeatures, "node-runtime")
	case "python":
		compatibility.RequiredFeatures = append(compatibility.RequiredFeatures, "python-runtime")
	case "docker":
		compatibility.RequiredFeatures = append(compatibility.RequiredFeatures, "docker-runtime")
	}

	return compatibility, nil
}

// GetTemplatesByCategory gets templates by category
func (m *Manager) GetTemplatesByCategory(category string) ([]interfaces.TemplateInfo, error) {
	// Get templates from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Use discovery component to filter by category
	filtered := m.discovery.GetTemplatesByCategory(templates, category)

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(filtered))
	for i, tmpl := range filtered {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// GetTemplatesByTechnology gets templates by technology
func (m *Manager) GetTemplatesByTechnology(technology string) ([]interfaces.TemplateInfo, error) {
	// Get templates from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Use discovery component to filter by technology
	filtered := m.discovery.GetTemplatesByTechnology(templates, technology)

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(filtered))
	for i, tmpl := range filtered {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// ValidateTemplateMetadata validates template metadata
func (m *Manager) ValidateTemplateMetadata(metadata *interfaces.TemplateMetadata) error {
	// Use the validator component to validate metadata
	return m.validator.ValidateTemplateMetadata(metadata)
}

// ValidateCustomTemplate validates custom template
func (m *Manager) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Use the validator component to validate custom template
	return m.validator.ValidateCustomTemplate(path)
}

// PreviewTemplate previews template processing
func (m *Manager) PreviewTemplate(templateName string, config *models.ProjectConfig) (*interfaces.TemplatePreview, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, templateName)
	if templateInfo == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", templateName)
	}

	preview := &interfaces.TemplatePreview{
		TemplateName: templateName,
		OutputPath:   config.OutputPath,
		Files:        []interfaces.TemplatePreviewFile{},
		Variables:    make(map[string]any),
		Summary: interfaces.TemplatePreviewSummary{
			TotalFiles:       0,
			TotalDirectories: 0,
			TotalSize:        0,
			TemplatedFiles:   0,
			ExecutableFiles:  0,
		},
	}

	// Add config variables to preview
	preview.Variables["Name"] = config.Name
	preview.Variables["Organization"] = config.Organization
	preview.Variables["Description"] = config.Description
	preview.Variables["License"] = config.License

	// For embedded templates, walk through the embedded filesystem
	if templateInfo.Source == "embedded" {
		templatePath := fmt.Sprintf("templates/%s/%s", templateInfo.Category, templateName)
		err = m.previewEmbeddedTemplate(templatePath, preview)
	} else {
		// For file-based templates, walk through the file system
		err = m.previewFileTemplate(templateInfo.Path, preview)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate preview: %w", err)
	}

	return preview, nil
}

// previewEmbeddedTemplate generates preview for embedded template
func (m *Manager) previewEmbeddedTemplate(templatePath string, preview *interfaces.TemplatePreview) error {
	return fs.WalkDir(embeddedTemplates, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory
		if path == templatePath {
			return nil
		}

		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			preview.Summary.TotalDirectories++
			preview.Files = append(preview.Files, interfaces.TemplatePreviewFile{
				Path: relPath,
				Type: "directory",
			})
		} else {
			preview.Summary.TotalFiles++

			info, err := d.Info()
			if err != nil {
				return err
			}

			size := info.Size()
			preview.Summary.TotalSize += size

			isTemplated := strings.HasSuffix(path, ".tmpl")
			if isTemplated {
				preview.Summary.TemplatedFiles++
			}

			// Check if file would be executable (simplified check)
			isExecutable := strings.Contains(relPath, "script") || strings.HasSuffix(relPath, ".sh")
			if isExecutable {
				preview.Summary.ExecutableFiles++
			}

			preview.Files = append(preview.Files, interfaces.TemplatePreviewFile{
				Path:       relPath,
				Type:       "file",
				Size:       size,
				Templated:  isTemplated,
				Executable: isExecutable,
			})
		}

		return nil
	})
}

// previewFileTemplate generates preview for file-based template
func (m *Manager) previewFileTemplate(templatePath string, preview *interfaces.TemplatePreview) error {
	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory
		if path == templatePath {
			return nil
		}

		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			preview.Summary.TotalDirectories++
			preview.Files = append(preview.Files, interfaces.TemplatePreviewFile{
				Path: relPath,
				Type: "directory",
			})
		} else {
			preview.Summary.TotalFiles++

			size := info.Size()
			preview.Summary.TotalSize += size

			isTemplated := strings.HasSuffix(path, ".tmpl")
			if isTemplated {
				preview.Summary.TemplatedFiles++
			}

			// Check if file would be executable
			isExecutable := (info.Mode() & 0111) != 0
			if isExecutable {
				preview.Summary.ExecutableFiles++
			}

			preview.Files = append(preview.Files, interfaces.TemplatePreviewFile{
				Path:       relPath,
				Type:       "file",
				Size:       size,
				Templated:  isTemplated,
				Executable: isExecutable,
			})
		}

		return nil
	})
}

// GetTemplateVariables gets template variables
func (m *Manager) GetTemplateVariables(name string) (map[string]interfaces.TemplateVariable, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, name)
	if templateInfo == nil {
		return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	variables := make(map[string]interfaces.TemplateVariable)

	// Convert from models.TemplateVar to interfaces.TemplateVariable
	for varName, templateVar := range templateInfo.Metadata.Variables {
		variable := interfaces.TemplateVariable{
			Name:        templateVar.Name,
			Type:        templateVar.Type,
			Description: templateVar.Description,
			Default:     templateVar.Default,
			Required:    templateVar.Required,
		}

		// Add validation if pattern or enum is specified
		if templateVar.Pattern != "" || len(templateVar.Enum) > 0 {
			variable.Validation = &interfaces.VariableValidation{
				Pattern: templateVar.Pattern,
			}
		}

		variables[varName] = variable
	}

	// Add default variables if none are defined
	if len(variables) == 0 {
		variables = m.getDefaultTemplateVariables()
	}

	return variables, nil
}

// getDefaultTemplateVariables returns default template variables
func (m *Manager) getDefaultTemplateVariables() map[string]interfaces.TemplateVariable {
	return map[string]interfaces.TemplateVariable{
		"Name": {
			Name:        "Name",
			Type:        "string",
			Description: "Project name",
			Required:    true,
		},
		"Organization": {
			Name:        "Organization",
			Type:        "string",
			Description: "Organization or author name",
			Required:    false,
		},
		"Description": {
			Name:        "Description",
			Type:        "string",
			Description: "Project description",
			Required:    false,
		},
		"License": {
			Name:        "License",
			Type:        "string",
			Description: "Project license",
			Default:     "MIT",
			Required:    false,
		},
	}
}

// InstallTemplate installs a template
func (m *Manager) InstallTemplate(source string, name string) error {
	// For now, return not implemented as this requires external template support
	return fmt.Errorf("template installation from external sources not yet implemented")
}

// UninstallTemplate uninstalls a template
func (m *Manager) UninstallTemplate(name string) error {
	// For now, return not implemented as this requires external template support
	return fmt.Errorf("template uninstallation not yet implemented for embedded templates")
}

// UpdateTemplate updates a template
func (m *Manager) UpdateTemplate(name string) error {
	// For now, return not implemented as this requires external template support
	return fmt.Errorf("template updates not yet implemented for embedded templates")
}

// GetTemplateLocation gets template location
func (m *Manager) GetTemplateLocation(name string) (string, error) {
	// Get template info from cache or discover
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return "", fmt.Errorf("failed to discover templates: %w", err)
	}

	// Find the template using discovery component
	templateInfo := m.discovery.GetTemplateByName(templates, name)
	if templateInfo == nil {
		return "", fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
	}

	if templateInfo.Source == "embedded" {
		return fmt.Sprintf("embedded:templates/%s/%s", templateInfo.Category, name), nil
	}

	return templateInfo.Path, nil
}

// CacheTemplate caches a template
func (m *Manager) CacheTemplate(name string) error {
	// Use GetOrRefresh to ensure template is cached
	_, err := m.cache.GetOrRefresh(name)
	return err
}

// GetCachedTemplates gets cached templates
func (m *Manager) GetCachedTemplates() ([]interfaces.TemplateInfo, error) {
	// Get all templates from cache or refresh if needed
	templates, err := m.cache.GetAllOrRefresh()
	if err != nil {
		return nil, err
	}

	// Convert to interface type
	return m.cache.ConvertToInterfaceTemplateInfos(templates), nil
}

// ClearTemplateCache clears template cache
func (m *Manager) ClearTemplateCache() error {
	m.cache.Clear()
	return nil
}

// RefreshTemplateCache refreshes template cache
func (m *Manager) RefreshTemplateCache() error {
	return m.cache.Refresh()
}

// convertToInterfaceTemplateInfo converts models.TemplateInfo to interfaces.TemplateInfo
func (m *Manager) convertToInterfaceTemplateInfo(tmpl *models.TemplateInfo) interfaces.TemplateInfo {
	return interfaces.TemplateInfo{
		Name:         tmpl.Name,
		DisplayName:  tmpl.DisplayName,
		Description:  tmpl.Description,
		Category:     tmpl.Category,
		Technology:   tmpl.Technology,
		Version:      tmpl.Version,
		Tags:         tmpl.Tags,
		Dependencies: tmpl.Dependencies,
		Metadata: interfaces.TemplateMetadata{
			Author:     tmpl.Metadata.Author,
			License:    tmpl.Metadata.License,
			Repository: tmpl.Metadata.Repository,
			Homepage:   tmpl.Metadata.Homepage,
			Keywords:   tmpl.Metadata.Keywords,
			Created:    tmpl.Metadata.CreatedAt,
			Updated:    tmpl.Metadata.UpdatedAt,
		},
	}
}

// ValidateTemplateStructureAdvanced validates template structure with advanced checks
func (m *Manager) ValidateTemplateStructureAdvanced(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	return m.validator.ValidateTemplateStructureAdvanced(templateInfo)
}

// ValidateTemplateDependencies validates template dependencies
func (m *Manager) ValidateTemplateDependencies(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	return m.validator.ValidateTemplateDependencies(templateInfo)
}

// ValidateTemplateCompatibility validates template compatibility
func (m *Manager) ValidateTemplateCompatibility(templateInfo *interfaces.TemplateInfo) (*interfaces.TemplateValidationResult, error) {
	return m.validator.ValidateTemplateCompatibility(templateInfo)
}

// ValidateTemplateMetadataAdvanced validates template metadata with advanced checks
func (m *Manager) ValidateTemplateMetadataAdvanced(metadata *interfaces.TemplateMetadata) (*interfaces.TemplateValidationResult, error) {
	return m.validator.ValidateTemplateMetadataAdvanced(metadata)
}

// ValidateTemplateComprehensive performs comprehensive template validation
func (m *Manager) ValidateTemplateComprehensive(templateName string) (*interfaces.TemplateValidationResult, error) {
	// Get template info
	templateInfo, err := m.GetTemplateInfo(templateName)
	if err != nil {
		return &interfaces.TemplateValidationResult{
			Valid: false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: "error",
					Message:  fmt.Sprintf("Template '%s' not found: %v", templateName, err),
					Rule:     "template-not-found",
					Fixable:  false,
				},
			},
		}, nil
	}

	var allIssues []interfaces.ValidationIssue
	var allWarnings []interfaces.ValidationIssue

	// Validate structure
	structureResult, err := m.ValidateTemplateStructureAdvanced(templateInfo)
	if err != nil {
		return nil, fmt.Errorf("structure validation failed: %w", err)
	}
	allIssues = append(allIssues, structureResult.Issues...)
	allWarnings = append(allWarnings, structureResult.Warnings...)

	// Validate dependencies
	depResult, err := m.ValidateTemplateDependencies(templateInfo)
	if err != nil {
		return nil, fmt.Errorf("dependency validation failed: %w", err)
	}
	allIssues = append(allIssues, depResult.Issues...)
	allWarnings = append(allWarnings, depResult.Warnings...)

	// Validate compatibility
	compatResult, err := m.ValidateTemplateCompatibility(templateInfo)
	if err != nil {
		return nil, fmt.Errorf("compatibility validation failed: %w", err)
	}
	allIssues = append(allIssues, compatResult.Issues...)
	allWarnings = append(allWarnings, compatResult.Warnings...)

	// Validate metadata
	metadataResult, err := m.ValidateTemplateMetadataAdvanced(&templateInfo.Metadata)
	if err != nil {
		return nil, fmt.Errorf("metadata validation failed: %w", err)
	}
	allIssues = append(allIssues, metadataResult.Issues...)
	allWarnings = append(allWarnings, metadataResult.Warnings...)

	return &interfaces.TemplateValidationResult{
		Valid:    len(allIssues) == 0,
		Issues:   allIssues,
		Warnings: allWarnings,
		Summary: interfaces.ValidationSummary{
			ErrorCount:   len(allIssues),
			WarningCount: len(allWarnings),
			FixableCount: 0,
		},
	}, nil
}
