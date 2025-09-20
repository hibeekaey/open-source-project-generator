// Package template provides template management functionality for the
// Open Source Project Generator.
package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	yaml "gopkg.in/yaml.v3"
)

// Manager implements the TemplateManager interface for template operations.
type Manager struct {
	templateEngine interfaces.TemplateEngine
	cache          map[string]*models.TemplateInfo
	cacheTime      time.Time
	cacheTTL       time.Duration
}

// NewManager creates a new template manager instance.
func NewManager(templateEngine interfaces.TemplateEngine) interfaces.TemplateManager {
	return &Manager{
		templateEngine: templateEngine,
		cache:          make(map[string]*models.TemplateInfo),
		cacheTTL:       5 * time.Minute, // Cache templates for 5 minutes
	}
}

// ListTemplates lists available templates with optional filtering
func (m *Manager) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("🚫 couldn't find available templates: %w", err)
	}

	// Apply filters
	filtered := m.applyFilters(templates, filter)

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(filtered))
	for i, tmpl := range filtered {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// GetTemplateInfo retrieves detailed information about a specific template
func (m *Manager) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	for _, tmpl := range templates {
		if tmpl.Name == name {
			result := m.convertToInterfaceTemplateInfo(tmpl)
			return &result, nil
		}
	}

	return nil, fmt.Errorf("🚫 Template '%s' not found. Use 'generator list-templates' to see available options", name)
}

// SearchTemplates searches for templates by query string
func (m *Manager) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	query = strings.ToLower(query)
	var matches []*models.TemplateInfo

	for _, tmpl := range templates {
		// Search in name, display name, description, tags, and keywords
		if m.matchesQuery(tmpl, query) {
			matches = append(matches, tmpl)
		}
	}

	// Convert to interface type
	result := make([]interfaces.TemplateInfo, len(matches))
	for i, tmpl := range matches {
		result[i] = m.convertToInterfaceTemplateInfo(tmpl)
	}

	return result, nil
}

// ValidateTemplate validates a template structure and metadata
func (m *Manager) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &interfaces.TemplateValidationResult{
			Valid: false,
			Issues: []interfaces.ValidationIssue{
				{
					Type:     "error",
					Severity: "error",
					Message:  fmt.Sprintf("Template path does not exist: %s", path),
					Rule:     "path-exists",
					Fixable:  false,
				},
			},
		}, nil
	}

	var issues []models.ValidationIssue
	var warnings []models.ValidationIssue

	// Validate template structure
	structureIssues := m.validateTemplateStructure(path)
	issues = append(issues, structureIssues...)

	// Validate metadata if present
	metadataIssues := m.validateTemplateMetadataFile(path)
	issues = append(issues, metadataIssues...)

	// Validate template files
	templateIssues := m.validateTemplateFiles(path)
	issues = append(issues, templateIssues...)

	// Separate errors from warnings
	var errors []models.ValidationIssue
	for _, issue := range issues {
		if issue.Severity == "error" {
			errors = append(errors, issue)
		} else {
			warnings = append(warnings, issue)
		}
	}

	// Convert models.ValidationIssue to interfaces.ValidationIssue
	interfaceErrors := make([]interfaces.ValidationIssue, len(errors))
	for i, issue := range errors {
		interfaceErrors[i] = interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
			Metadata: issue.Metadata,
		}
	}

	interfaceWarnings := make([]interfaces.ValidationIssue, len(warnings))
	for i, issue := range warnings {
		interfaceWarnings[i] = interfaces.ValidationIssue{
			Type:     issue.Type,
			Severity: issue.Severity,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
			Metadata: issue.Metadata,
		}
	}

	return &interfaces.TemplateValidationResult{
		Valid:    len(errors) == 0,
		Issues:   interfaceErrors,
		Warnings: interfaceWarnings,
	}, nil
}

// ValidateTemplateStructure validates template structure
func (m *Manager) ValidateTemplateStructure(template *interfaces.TemplateInfo) error {
	// Validate required fields
	if template.Name == "" {
		return fmt.Errorf("🚫 template name is required")
	}
	if template.Category == "" {
		return fmt.Errorf("🚫 template category is required")
	}
	if template.Version == "" {
		return fmt.Errorf("🚫 template version is required")
	}

	// Validate name format (should be kebab-case)
	if !m.isValidTemplateName(template.Name) {
		return fmt.Errorf("🚫 template name must be in kebab-case format")
	}

	// Validate category
	validCategories := []string{"backend", "frontend", "mobile", "infrastructure", "base"}
	if !m.contains(validCategories, template.Category) {
		return fmt.Errorf("🚫 invalid category: %s. Valid categories: %v", template.Category, validCategories)
	}

	return nil
}

// ProcessTemplate processes a template with the given configuration
func (m *Manager) ProcessTemplate(templateName string, config *models.ProjectConfig, outputPath string) error {
	// Get template info from internal models (which has Source and Path)
	templates, err := m.discoverTemplates()
	if err != nil {
		return fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == templateName {
			templateInfo = tmpl
			break
		}
	}

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

	// For embedded templates, use the embedded template engine
	if templateInfo.Source == "embedded" {
		templatePath := fmt.Sprintf("templates/%s/%s", templateInfo.Category, templateName)
		return m.templateEngine.ProcessDirectory(templatePath, outputPath, &processConfig)
	}

	// For file-based templates, process the directory directly
	return m.templateEngine.ProcessDirectory(templateInfo.Path, outputPath, &processConfig)
}

// ProcessCustomTemplate processes a custom template from a path
func (m *Manager) ProcessCustomTemplate(templatePath string, config *models.ProjectConfig, outputPath string) error {
	// Validate template first
	validationResult, err := m.ValidateTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to validate template: %w", err)
	}

	if !validationResult.Valid {
		return fmt.Errorf("template validation failed: %d errors found", len(validationResult.Issues))
	}

	// Process the template directory
	return m.templateEngine.ProcessDirectory(templatePath, outputPath, config)
}

// GetTemplateMetadata retrieves template metadata
func (m *Manager) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	// Get template info from internal models
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == name {
			templateInfo = tmpl
			break
		}
	}

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
	templateInfo, err := m.GetTemplateInfo(name)
	if err != nil {
		return nil, err
	}

	return templateInfo.Dependencies, nil
}

// GetTemplateCompatibility retrieves template compatibility information
func (m *Manager) GetTemplateCompatibility(name string) (*interfaces.CompatibilityInfo, error) {
	// Get template info from internal models (which has MinVersion and MaxVersion)
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == name {
			templateInfo = tmpl
			break
		}
	}

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
	filter := interfaces.TemplateFilter{
		Category: category,
	}
	return m.ListTemplates(filter)
}

// GetTemplatesByTechnology gets templates by technology
func (m *Manager) GetTemplatesByTechnology(technology string) ([]interfaces.TemplateInfo, error) {
	filter := interfaces.TemplateFilter{
		Technology: technology,
	}
	return m.ListTemplates(filter)
}

// ValidateTemplateMetadata validates template metadata
func (m *Manager) ValidateTemplateMetadata(metadata *interfaces.TemplateMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	if metadata.Author == "" {
		return fmt.Errorf("metadata author is required")
	}
	if metadata.License == "" {
		return fmt.Errorf("metadata license is required")
	}

	return nil
}

// ValidateCustomTemplate validates custom template
func (m *Manager) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	// Use the same validation logic as ValidateTemplate
	return m.ValidateTemplate(path)
}

// PreviewTemplate previews template processing
func (m *Manager) PreviewTemplate(templateName string, config *models.ProjectConfig) (*interfaces.TemplatePreview, error) {
	// Get template info from internal models
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == templateName {
			templateInfo = tmpl
			break
		}
	}

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
	// Get template info from internal models
	templates, err := m.discoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == name {
			templateInfo = tmpl
			break
		}
	}

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
	// Get template info from internal models
	templates, err := m.discoverTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to discover templates: %w", err)
	}

	var templateInfo *models.TemplateInfo
	for _, tmpl := range templates {
		if tmpl.Name == name {
			templateInfo = tmpl
			break
		}
	}

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
	// Template is already cached in memory, so just verify it exists
	_, err := m.GetTemplateInfo(name)
	return err
}

// GetCachedTemplates gets cached templates
func (m *Manager) GetCachedTemplates() ([]interfaces.TemplateInfo, error) {
	// Return all templates as they are cached in memory
	return m.ListTemplates(interfaces.TemplateFilter{})
}

// ClearTemplateCache clears template cache
func (m *Manager) ClearTemplateCache() error {
	m.cache = make(map[string]*models.TemplateInfo)
	m.cacheTime = time.Time{}
	return nil
}

// RefreshTemplateCache refreshes template cache
func (m *Manager) RefreshTemplateCache() error {
	m.cache = make(map[string]*models.TemplateInfo)
	m.cacheTime = time.Time{}
	_, err := m.discoverTemplates()
	return err
}

// discoverTemplates discovers all available templates from embedded filesystem
func (m *Manager) discoverTemplates() ([]*models.TemplateInfo, error) {
	// Check cache first
	if time.Since(m.cacheTime) < m.cacheTTL && len(m.cache) > 0 {
		templates := make([]*models.TemplateInfo, 0, len(m.cache))
		for _, tmpl := range m.cache {
			templates = append(templates, tmpl)
		}
		return templates, nil
	}

	var templates []*models.TemplateInfo

	// Discover embedded templates
	embeddedTemplates, err := m.discoverEmbeddedTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover embedded templates: %w", err)
	}
	templates = append(templates, embeddedTemplates...)

	// Update cache
	m.cache = make(map[string]*models.TemplateInfo)
	for _, tmpl := range templates {
		m.cache[tmpl.Name] = tmpl
	}
	m.cacheTime = time.Now()

	// Sort templates by name
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// discoverEmbeddedTemplates discovers templates from the embedded filesystem
func (m *Manager) discoverEmbeddedTemplates() ([]*models.TemplateInfo, error) {
	var templates []*models.TemplateInfo

	// Walk through embedded template directories
	err := fs.WalkDir(embeddedTemplates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root templates directory
		if path == "templates" {
			return nil
		}

		// Only process directories that are direct children of category directories
		if d.IsDir() && m.isTemplateDirectory(path) {
			templateInfo, err := m.createTemplateInfoFromPath(path)
			if err != nil {
				// Log error but continue processing other templates
				fmt.Printf("⚠️  Failed to process template at %s: %v\n", path, err)
				return nil
			}
			templates = append(templates, templateInfo)
		}

		return nil
	})

	return templates, err
}

// isTemplateDirectory checks if a directory path represents a template
func (m *Manager) isTemplateDirectory(path string) bool {
	parts := strings.Split(path, "/")
	// Should be templates/category/template-name
	return len(parts) == 3 && parts[0] == "templates"
}

// createTemplateInfoFromPath creates template info from a filesystem path
func (m *Manager) createTemplateInfoFromPath(path string) (*models.TemplateInfo, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid template path structure: %s", path)
	}

	category := parts[1]
	templateName := parts[2]

	// Create basic template info
	templateInfo := &models.TemplateInfo{
		Name:        templateName,
		DisplayName: m.formatDisplayName(templateName),
		Category:    category,
		Path:        path,
		Source:      "embedded",
		Version:     "1.0.0", // Default version
	}

	// Try to load metadata from template.yaml or template.yml
	metadata, err := m.loadTemplateMetadata(path)
	if err == nil {
		templateInfo.Metadata = *metadata
		templateInfo.DisplayName = metadata.DisplayName
		templateInfo.Description = metadata.Description
		templateInfo.Technology = metadata.Technology
		templateInfo.Tags = metadata.Tags
		templateInfo.Dependencies = metadata.Dependencies
		templateInfo.Version = metadata.Version
	} else {
		// Set defaults based on path analysis
		templateInfo.Description = fmt.Sprintf("%s template for %s projects",
			m.formatDisplayName(templateName), category)
		templateInfo.Technology = m.inferTechnology(templateName)
		templateInfo.Tags = m.inferTags(templateName, category)
	}

	// Calculate template size and file count
	size, fileCount, err := m.calculateTemplateStats(path)
	if err == nil {
		templateInfo.Size = size
		templateInfo.FileCount = fileCount
	}

	templateInfo.LastModified = time.Now() // For embedded templates, use current time

	return templateInfo, nil
}

// loadTemplateMetadata loads metadata from template.yaml or template.yml
func (m *Manager) loadTemplateMetadata(templatePath string) (*models.TemplateMetadata, error) {
	// Try template.yaml first, then template.yml
	metadataFiles := []string{"template.yaml", "template.yml"}

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if content, err := fs.ReadFile(embeddedTemplates, metadataPath); err == nil {
			return m.parseTemplateYAML(content, filepath.Base(templatePath))
		}
	}

	return nil, fmt.Errorf("no metadata file found")
}

// parseTemplateYAML parses template.yaml content into models.TemplateMetadata
func (m *Manager) parseTemplateYAML(content []byte, templateName string) (*models.TemplateMetadata, error) {
	// Define a structure that matches the template.yaml format
	type TemplateYAML struct {
		Name         string   `yaml:"name"`
		DisplayName  string   `yaml:"display_name"`
		Description  string   `yaml:"description"`
		Category     string   `yaml:"category"`
		Technology   string   `yaml:"technology"`
		Version      string   `yaml:"version"`
		Tags         []string `yaml:"tags"`
		Dependencies []string `yaml:"dependencies"`
		Metadata     struct {
			Author      string            `yaml:"author"`
			License     string            `yaml:"license"`
			Repository  string            `yaml:"repository"`
			Homepage    string            `yaml:"homepage"`
			Keywords    []string          `yaml:"keywords"`
			Maintainers []string          `yaml:"maintainers"`
			Created     time.Time         `yaml:"created"`
			Updated     time.Time         `yaml:"updated"`
			Variables   map[string]string `yaml:"variables"`
		} `yaml:"metadata"`
	}

	var yamlData TemplateYAML
	if err := yaml.Unmarshal(content, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse template YAML: %w", err)
	}

	// Convert to models.TemplateMetadata
	metadata := &models.TemplateMetadata{
		Name:         yamlData.Name,
		DisplayName:  yamlData.DisplayName,
		Description:  yamlData.Description,
		Version:      yamlData.Version,
		Author:       yamlData.Metadata.Author,
		License:      yamlData.Metadata.License,
		Category:     yamlData.Category,
		Technology:   yamlData.Technology,
		Tags:         yamlData.Tags,
		Dependencies: yamlData.Dependencies,
		CreatedAt:    yamlData.Metadata.Created,
		UpdatedAt:    yamlData.Metadata.Updated,
		Homepage:     yamlData.Metadata.Homepage,
		Repository:   yamlData.Metadata.Repository,
		Keywords:     yamlData.Metadata.Keywords,
		Variables:    make(map[string]models.TemplateVar),
	}

	// Convert variables from simple string map to TemplateVar map
	for name, description := range yamlData.Metadata.Variables {
		metadata.Variables[name] = models.TemplateVar{
			Name:        name,
			Type:        "string", // Default type
			Description: description,
			Required:    false, // Default to not required
		}
	}

	// Set defaults if not provided
	if metadata.Name == "" {
		metadata.Name = templateName
	}
	if metadata.DisplayName == "" {
		metadata.DisplayName = m.formatDisplayName(templateName)
	}
	if metadata.Version == "" {
		metadata.Version = "1.0.0"
	}
	if metadata.License == "" {
		metadata.License = "MIT"
	}
	if metadata.Author == "" {
		metadata.Author = "Open Source Project Generator"
	}

	// Initialize slices if nil to prevent issues
	if metadata.Tags == nil {
		metadata.Tags = []string{}
	}
	if metadata.Dependencies == nil {
		metadata.Dependencies = []string{}
	}
	if metadata.Keywords == nil {
		metadata.Keywords = []string{}
	}

	return metadata, nil
}

// calculateTemplateStats calculates size and file count for a template
func (m *Manager) calculateTemplateStats(templatePath string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := fs.WalkDir(embeddedTemplates, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			fileCount++
			if info, err := d.Info(); err == nil {
				totalSize += info.Size()
			}
		}

		return nil
	})

	return totalSize, fileCount, err
}

// formatDisplayName formats a template name for display
func (m *Manager) formatDisplayName(name string) string {
	// Convert kebab-case to Title Case
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// inferTechnology infers technology from template name
func (m *Manager) inferTechnology(templateName string) string {
	name := strings.ToLower(templateName)

	if strings.Contains(name, "go") || strings.Contains(name, "gin") {
		return "Go"
	}
	if strings.Contains(name, "nextjs") || strings.Contains(name, "next") {
		return "Next.js"
	}
	if strings.Contains(name, "react") {
		return "React"
	}
	if strings.Contains(name, "vue") {
		return "Vue.js"
	}
	if strings.Contains(name, "angular") {
		return "Angular"
	}
	if strings.Contains(name, "node") {
		return "Node.js"
	}
	if strings.Contains(name, "python") || strings.Contains(name, "django") || strings.Contains(name, "flask") {
		return "Python"
	}
	if strings.Contains(name, "java") || strings.Contains(name, "spring") {
		return "Java"
	}
	if strings.Contains(name, "kotlin") {
		return "Kotlin"
	}
	if strings.Contains(name, "swift") || strings.Contains(name, "ios") {
		return "Swift"
	}
	if strings.Contains(name, "terraform") {
		return "Terraform"
	}
	if strings.Contains(name, "docker") || strings.Contains(name, "kubernetes") {
		return "Docker"
	}

	return "Unknown"
}

// inferTags infers tags from template name and category
func (m *Manager) inferTags(templateName, category string) []string {
	var tags []string

	// Add category as a tag
	tags = append(tags, category)

	name := strings.ToLower(templateName)

	// Technology tags
	if strings.Contains(name, "go") {
		tags = append(tags, "golang", "backend")
	}
	if strings.Contains(name, "gin") {
		tags = append(tags, "gin", "web-framework", "api")
	}
	if strings.Contains(name, "nextjs") {
		tags = append(tags, "nextjs", "react", "frontend", "ssr")
	}
	if strings.Contains(name, "admin") {
		tags = append(tags, "admin", "dashboard")
	}
	if strings.Contains(name, "app") {
		tags = append(tags, "application", "web-app")
	}
	if strings.Contains(name, "home") {
		tags = append(tags, "landing-page", "website")
	}
	if strings.Contains(name, "android") {
		tags = append(tags, "android", "mobile")
	}
	if strings.Contains(name, "ios") {
		tags = append(tags, "ios", "mobile")
	}
	if strings.Contains(name, "kotlin") {
		tags = append(tags, "kotlin", "android")
	}
	if strings.Contains(name, "swift") {
		tags = append(tags, "swift", "ios")
	}
	if strings.Contains(name, "docker") {
		tags = append(tags, "docker", "containerization")
	}
	if strings.Contains(name, "kubernetes") {
		tags = append(tags, "kubernetes", "k8s", "orchestration")
	}
	if strings.Contains(name, "terraform") {
		tags = append(tags, "terraform", "iac", "infrastructure")
	}

	return tags
}

// applyFilters applies filtering criteria to templates
func (m *Manager) applyFilters(templates []*models.TemplateInfo, filter interfaces.TemplateFilter) []*models.TemplateInfo {
	var filtered []*models.TemplateInfo

	for _, tmpl := range templates {
		if m.matchesFilter(tmpl, filter) {
			filtered = append(filtered, tmpl)
		}
	}

	return filtered
}

// matchesFilter checks if a template matches the given filter
func (m *Manager) matchesFilter(tmpl *models.TemplateInfo, filter interfaces.TemplateFilter) bool {
	// Category filter
	if filter.Category != "" && !strings.EqualFold(tmpl.Category, filter.Category) {
		return false
	}

	// Technology filter
	if filter.Technology != "" && !strings.EqualFold(tmpl.Technology, filter.Technology) {
		return false
	}

	// Tags filter (template must have all specified tags)
	if len(filter.Tags) > 0 {
		templateTags := make(map[string]bool)
		for _, tag := range tmpl.Tags {
			templateTags[strings.ToLower(tag)] = true
		}

		for _, filterTag := range filter.Tags {
			if !templateTags[strings.ToLower(filterTag)] {
				return false
			}
		}
	}

	// Version filters (simplified - would need proper semver comparison)
	if filter.MinVersion != "" && tmpl.Version < filter.MinVersion {
		return false
	}
	if filter.MaxVersion != "" && tmpl.Version > filter.MaxVersion {
		return false
	}

	return true
}

// matchesQuery checks if a template matches a search query
func (m *Manager) matchesQuery(tmpl *models.TemplateInfo, query string) bool {
	// Search in name
	if strings.Contains(strings.ToLower(tmpl.Name), query) {
		return true
	}

	// Search in display name
	if strings.Contains(strings.ToLower(tmpl.DisplayName), query) {
		return true
	}

	// Search in description
	if strings.Contains(strings.ToLower(tmpl.Description), query) {
		return true
	}

	// Search in technology
	if strings.Contains(strings.ToLower(tmpl.Technology), query) {
		return true
	}

	// Search in tags
	for _, tag := range tmpl.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	// Search in keywords
	for _, keyword := range tmpl.Metadata.Keywords {
		if strings.Contains(strings.ToLower(keyword), query) {
			return true
		}
	}

	return false
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

// validateTemplateStructure validates the basic structure of a template directory
func (m *Manager) validateTemplateStructure(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check if it's a directory
	info, err := os.Stat(templatePath)
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot access template path: %v", err),
			Rule:     "path-accessible",
			Fixable:  false,
		})
		return issues
	}

	if !info.IsDir() {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Template path must be a directory",
			Rule:     "is-directory",
			Fixable:  false,
		})
		return issues
	}

	// Check for required files/directories
	requiredItems := []string{
		// At least one template file should exist
	}

	for _, item := range requiredItems {
		itemPath := filepath.Join(templatePath, item)
		if _, err := os.Stat(itemPath); os.IsNotExist(err) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Recommended item missing: %s", item),
				Rule:     "recommended-structure",
				Fixable:  true,
			})
		}
	}

	// Check for template files
	hasTemplateFiles, err := m.hasTemplateFiles(templatePath)
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Error checking template files: %v", err),
			Rule:     "template-files-check",
			Fixable:  false,
		})
	} else if !hasTemplateFiles {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "No template files (.tmpl) found in template directory",
			Rule:     "has-template-files",
			Fixable:  false,
		})
	}

	return issues
}

// validateTemplateMetadataFile validates template metadata file if present
func (m *Manager) validateTemplateMetadataFile(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check for metadata files
	metadataFiles := []string{"template.yaml", "template.yml", "metadata.yaml", "metadata.yml"}
	var foundMetadata bool

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if _, err := os.Stat(metadataPath); err == nil {
			foundMetadata = true
			// Validate metadata file content
			if validationIssues := m.validateMetadataFileContent(metadataPath); len(validationIssues) > 0 {
				issues = append(issues, validationIssues...)
			}
			break
		}
	}

	if !foundMetadata {
		issues = append(issues, models.ValidationIssue{
			Type:     "warning",
			Severity: "warning",
			Message:  "No metadata file found (template.yaml, template.yml, metadata.yaml, or metadata.yml)",
			Rule:     "has-metadata",
			Fixable:  true,
		})
	}

	return issues
}

// validateTemplateFiles validates individual template files
func (m *Manager) validateTemplateFiles(templatePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check template files
		if strings.HasSuffix(path, ".tmpl") {
			if fileIssues := m.validateTemplateFile(path); len(fileIssues) > 0 {
				issues = append(issues, fileIssues...)
			}
		}

		return nil
	})

	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Error walking template directory: %v", err),
			Rule:     "directory-walk",
			Fixable:  false,
		})
	}

	return issues
}

// validateTemplateFile validates a single template file
func (m *Manager) validateTemplateFile(filePath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Validate file path to prevent path traversal attacks
	if err := m.validateFilePath(filePath); err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid file path: %v", err),
			File:     filePath,
			Rule:     "path-validation",
			Fixable:  false,
		})
		return issues
	}

	// Read file content
	content, err := os.ReadFile(filePath) // #nosec G304 - path is validated above
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot read template file: %v", err),
			File:     filePath,
			Rule:     "file-readable",
			Fixable:  false,
		})
		return issues
	}

	// Basic template syntax validation
	contentStr := string(content)

	// Check for unmatched template delimiters
	openCount := strings.Count(contentStr, "{{")
	closeCount := strings.Count(contentStr, "}}")

	if openCount != closeCount {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  "Unmatched template delimiters {{ and }}",
			File:     filePath,
			Rule:     "template-syntax",
			Fixable:  false,
		})
	}

	// Check for common template variables
	commonVars := []string{"{{.Name}}", "{{.Organization}}", "{{.Description}}"}
	hasVars := false
	for _, variable := range commonVars {
		if strings.Contains(contentStr, variable) {
			hasVars = true
			break
		}
	}

	if !hasVars && openCount > 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "info",
			Severity: "info",
			Message:  "Template file contains template syntax but no common variables",
			File:     filePath,
			Rule:     "has-common-vars",
			Fixable:  false,
		})
	}

	return issues
}

// validateMetadataFileContent validates the content of a metadata file
func (m *Manager) validateMetadataFileContent(metadataPath string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Validate file path to prevent path traversal attacks
	if err := m.validateFilePath(metadataPath); err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid metadata file path: %v", err),
			File:     metadataPath,
			Rule:     "path-validation",
			Fixable:  false,
		})
		return issues
	}

	// Read metadata file
	content, err := os.ReadFile(metadataPath) // #nosec G304 - path is validated above
	if err != nil {
		issues = append(issues, models.ValidationIssue{
			Type:     "error",
			Severity: "error",
			Message:  fmt.Sprintf("Cannot read metadata file: %v", err),
			File:     metadataPath,
			Rule:     "metadata-readable",
			Fixable:  false,
		})
		return issues
	}

	// Basic YAML syntax check (simplified)
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for basic YAML key-value format
		if !strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  "Line does not appear to be valid YAML",
				File:     metadataPath,
				Line:     i + 1,
				Rule:     "yaml-syntax",
				Fixable:  false,
			})
		}
	}

	// Check for required metadata fields
	requiredFields := []string{"name:", "description:", "version:", "author:"}
	for _, field := range requiredFields {
		if !strings.Contains(contentStr, field) {
			issues = append(issues, models.ValidationIssue{
				Type:     "warning",
				Severity: "warning",
				Message:  fmt.Sprintf("Missing recommended field: %s", strings.TrimSuffix(field, ":")),
				File:     metadataPath,
				Rule:     "required-fields",
				Fixable:  true,
			})
		}
	}

	return issues
}

// hasTemplateFiles checks if directory contains template files
func (m *Manager) hasTemplateFiles(templatePath string) (bool, error) {
	hasFiles := false

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			hasFiles = true
			return filepath.SkipDir // Stop walking once we find a template file
		}

		return nil
	})

	return hasFiles, err
}

// isValidTemplateName checks if template name follows kebab-case convention
func (m *Manager) isValidTemplateName(name string) bool {
	// Basic kebab-case validation: lowercase letters, numbers, and hyphens only
	for _, char := range name {
		if char < 'a' || char > 'z' {
			if char < '0' || char > '9' {
				if char != '-' {
					return false
				}
			}
		}
	}

	// Should not start or end with hyphen
	return !strings.HasPrefix(name, "-") && !strings.HasSuffix(name, "-")
}

// contains checks if slice contains string
func (m *Manager) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// validateFilePath validates file path to prevent path traversal attacks
func (m *Manager) validateFilePath(filePath string) error {
	// Clean the path to resolve any .. or . elements
	cleanPath := filepath.Clean(filePath)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in file path")
	}

	// Ensure path is absolute or relative to current directory
	if filepath.IsAbs(cleanPath) {
		// For absolute paths, ensure they don't access system directories
		systemDirs := []string{"/etc", "/proc", "/sys", "/dev", "/root"}
		for _, sysDir := range systemDirs {
			if strings.HasPrefix(cleanPath, sysDir) {
				return fmt.Errorf("access to system directory not allowed")
			}
		}
	}

	return nil
}
