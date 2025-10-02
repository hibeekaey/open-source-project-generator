// Package ui provides project structure preview management functionality.
package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/components"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/display"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// PreviewManager handles generation and display of project structure previews
type PreviewManager struct {
	ui              interfaces.InteractiveUIInterface
	templateManager interfaces.TemplateManager
	logger          interfaces.Logger

	// Component generators
	frontendPreview       *components.FrontendPreview
	backendPreview        *components.BackendPreview
	infrastructurePreview *components.InfrastructurePreview

	// Display handlers
	interactiveDisplay *display.Interactive
	consoleDisplay     *display.Console

	// Tree utilities
	treeBuilder   *tree.Builder
	treeFormatter *tree.Formatter
}

// NewPreviewManager creates a new project structure preview manager
func NewPreviewManager(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *PreviewManager {
	return &PreviewManager{
		ui:                    ui,
		templateManager:       templateManager,
		logger:                logger,
		frontendPreview:       components.NewFrontendPreview(),
		backendPreview:        components.NewBackendPreview(),
		infrastructurePreview: components.NewInfrastructurePreview(),
		interactiveDisplay:    display.NewInteractive(ui),
		consoleDisplay:        display.NewConsole(),
		treeBuilder:           tree.NewBuilder(),
		treeFormatter:         tree.NewFormatter(),
	}
}

// GenerateProjectStructurePreview creates a comprehensive project structure preview
func (pm *PreviewManager) GenerateProjectStructurePreview(ctx context.Context, config *models.ProjectConfig, selections []interfaces.TemplateSelection, outputDir string) (*display.ProjectStructurePreview, error) {
	if config == nil {
		return nil, fmt.Errorf("project config cannot be nil")
	}

	if len(selections) == 0 {
		return nil, fmt.Errorf("no templates selected")
	}

	pm.logger.InfoWithFields("Generating project structure preview", map[string]interface{}{
		"project_name":     config.Name,
		"template_count":   len(selections),
		"output_directory": outputDir,
	})

	// Convert selections to interface{} for display compatibility
	var displaySelections []interface{}
	for _, sel := range selections {
		displaySelections = append(displaySelections, sel)
	}

	// Initialize preview structure
	preview := &display.ProjectStructurePreview{
		ProjectName:       config.Name,
		OutputDirectory:   outputDir,
		SelectedTemplates: displaySelections,
		DirectoryTree:     &display.StandardDirectoryTree{},
		ComponentMapping:  []display.ComponentTemplateMapping{},
		Summary:           display.ProjectSummary{},
		Warnings:          []string{},
		Conflicts:         []display.FileConflict{},
	}

	// Generate individual template previews
	templatePreviews, err := pm.generateTemplatePreviews(ctx, selections, config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template previews: %w", err)
	}

	// Build standardized directory structure
	err = pm.buildStandardDirectoryStructure(preview, selections, templatePreviews, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build directory structure: %w", err)
	}

	// Generate component mappings
	pm.generateComponentMappings(preview, selections, templatePreviews)

	// Calculate project summary
	pm.calculateProjectSummary(preview, selections, templatePreviews)

	// Detect conflicts and generate warnings
	pm.detectConflictsAndWarnings(preview, selections, templatePreviews)

	return preview, nil
}

// DisplayProjectStructurePreview shows the project structure preview to the user
func (pm *PreviewManager) DisplayProjectStructurePreview(ctx context.Context, preview *display.ProjectStructurePreview) error {
	return pm.interactiveDisplay.DisplayProjectStructurePreview(ctx, preview)
}

// generateTemplatePreviews generates previews for all selected templates
func (pm *PreviewManager) generateTemplatePreviews(ctx context.Context, selections []interfaces.TemplateSelection, config *models.ProjectConfig) (map[string]*interfaces.TemplatePreview, error) {
	previews := make(map[string]*interfaces.TemplatePreview)

	for _, selection := range selections {
		preview, err := pm.templateManager.PreviewTemplate(selection.Template.Name, config)
		if err != nil {
			pm.logger.WarnWithFields("Failed to preview template", map[string]interface{}{
				"template": selection.Template.Name,
				"error":    err.Error(),
			})
			continue
		}
		previews[selection.Template.Name] = preview
	}

	return previews, nil
}

// buildStandardDirectoryStructure creates the standardized directory structure
func (pm *PreviewManager) buildStandardDirectoryStructure(preview *display.ProjectStructurePreview, selections []interfaces.TemplateSelection, templatePreviews map[string]*interfaces.TemplatePreview, config *models.ProjectConfig) error {
	// Create root directory
	preview.DirectoryTree.Root = &tree.DirectoryNode{
		Name:        config.Name,
		Path:        "",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "project-root",
		Description: "Project root directory",
	}

	// Analyze selected templates to determine which directories to create
	templatesByCategory := pm.categorizeTemplates(selections)

	// Create component directories using specialized generators
	if frontendTemplates, exists := templatesByCategory["frontend"]; exists && len(frontendTemplates) > 0 {
		convertedTemplates := pm.convertToPreviewTemplateSelection(frontendTemplates)
		appDir := pm.frontendPreview.CreateAppDirectory(convertedTemplates, templatePreviews)
		preview.DirectoryTree.App = appDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, appDir)
	}

	if backendTemplates, exists := templatesByCategory["backend"]; exists && len(backendTemplates) > 0 {
		convertedTemplates := pm.convertToPreviewTemplateSelection(backendTemplates)
		serverDir := pm.backendPreview.CreateCommonServerDirectory(convertedTemplates, templatePreviews)
		preview.DirectoryTree.CommonServer = serverDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, serverDir)
	}

	if mobileTemplates, exists := templatesByCategory["mobile"]; exists && len(mobileTemplates) > 0 {
		convertedTemplates := pm.convertToPreviewTemplateSelection(mobileTemplates)
		mobileDir := pm.infrastructurePreview.CreateMobileDirectory(convertedTemplates, templatePreviews)
		preview.DirectoryTree.Mobile = mobileDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, mobileDir)
	}

	if infraTemplates, exists := templatesByCategory["infrastructure"]; exists && len(infraTemplates) > 0 {
		convertedTemplates := pm.convertToPreviewTemplateSelection(infraTemplates)
		deployDir := pm.infrastructurePreview.CreateDeployDirectory(convertedTemplates, templatePreviews)
		preview.DirectoryTree.Deploy = deployDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, deployDir)
	}

	// Always create standard directories
	pm.createStandardDirectories(preview, templatePreviews)

	// Sort all directory children
	pm.treeBuilder.SortDirectoryTree(preview.DirectoryTree.Root)

	return nil
}

// createStandardDirectories creates standard project directories
func (pm *PreviewManager) createStandardDirectories(preview *display.ProjectStructurePreview, templatePreviews map[string]*interfaces.TemplatePreview) {
	// Create standard directories that are always present
	docsDir := pm.createDocsDirectory()
	preview.DirectoryTree.Docs = docsDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, docsDir)

	scriptsDir := pm.createScriptsDirectory()
	preview.DirectoryTree.Scripts = scriptsDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, scriptsDir)

	githubDir := pm.createGitHubDirectory()
	preview.DirectoryTree.GitHub = githubDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, githubDir)

	// Add common project files to root
	pm.addCommonProjectFiles(preview.DirectoryTree.Root)
}

// categorizeTemplates groups templates by their category
func (pm *PreviewManager) categorizeTemplates(selections []interfaces.TemplateSelection) map[string][]interfaces.TemplateSelection {
	categories := make(map[string][]interfaces.TemplateSelection)

	for _, selection := range selections {
		category := strings.ToLower(selection.Template.Category)
		categories[category] = append(categories[category], selection)
	}

	return categories
}

// generateComponentMappings creates mappings between components and templates
func (pm *PreviewManager) generateComponentMappings(preview *display.ProjectStructurePreview, selections []interfaces.TemplateSelection, templatePreviews map[string]*interfaces.TemplatePreview) {
	mappings := []display.ComponentTemplateMapping{}

	// Map each major directory to its contributing templates
	if preview.DirectoryTree.App != nil {
		mapping := pm.createComponentMapping("Frontend Applications", "App/", selections, templatePreviews, "frontend")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.CommonServer != nil {
		mapping := pm.createComponentMapping("Backend API", "CommonServer/", selections, templatePreviews, "backend")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.Mobile != nil {
		mapping := pm.createComponentMapping("Mobile Applications", "Mobile/", selections, templatePreviews, "mobile")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.Deploy != nil {
		mapping := pm.createComponentMapping("Infrastructure & Deployment", "Deploy/", selections, templatePreviews, "infrastructure")
		mappings = append(mappings, mapping)
	}

	// Always include standard components
	pm.addStandardComponentMappings(&mappings, preview)

	preview.ComponentMapping = mappings
}

// calculateProjectSummary calculates overall project statistics
func (pm *PreviewManager) calculateProjectSummary(preview *display.ProjectStructurePreview, selections []interfaces.TemplateSelection, templatePreviews map[string]*interfaces.TemplatePreview) {
	summary := display.ProjectSummary{
		Technologies: []string{},
		Categories:   []string{},
		Dependencies: []string{},
	}

	// Count directories and files recursively
	summary.TotalDirectories = pm.treeFormatter.CountDirectoriesRecursively(preview.DirectoryTree.Root)
	summary.TotalFiles = pm.treeFormatter.CountFilesRecursively(preview.DirectoryTree.Root)
	summary.EstimatedSize = pm.treeFormatter.CalculateDirectorySize(preview.DirectoryTree.Root)

	// Collect unique technologies, categories, and dependencies
	pm.collectProjectMetadata(&summary, selections)

	preview.Summary = summary
}

// detectConflictsAndWarnings identifies potential issues with the project structure
func (pm *PreviewManager) detectConflictsAndWarnings(preview *display.ProjectStructurePreview, selections []interfaces.TemplateSelection, templatePreviews map[string]*interfaces.TemplatePreview) {
	var warnings []string
	var conflicts []display.FileConflict

	// Check for large project size
	if preview.Summary.EstimatedSize > 500*1024*1024 { // 500MB
		warnings = append(warnings, fmt.Sprintf("Large project size: %s", pm.treeFormatter.FormatBytes(preview.Summary.EstimatedSize)))
	}

	// Check for many files
	if preview.Summary.TotalFiles > 2000 {
		warnings = append(warnings, fmt.Sprintf("Large number of files: %d", preview.Summary.TotalFiles))
	}

	// Check for missing dependencies and file conflicts
	pm.checkDependenciesAndConflicts(&warnings, &conflicts, selections, templatePreviews)

	preview.Warnings = warnings
	preview.Conflicts = conflicts
}

// convertToPreviewTemplateSelection converts TemplateSelection to preview.TemplateSelection
func (pm *PreviewManager) convertToPreviewTemplateSelection(selections []interfaces.TemplateSelection) []preview.TemplateSelection {
	var converted []preview.TemplateSelection
	for _, sel := range selections {
		converted = append(converted, preview.TemplateSelection{
			Template: sel.Template,
			Selected: sel.Selected,
		})
	}
	return converted
}

// Helper methods for creating standard directories and mappings
func (pm *PreviewManager) createDocsDirectory() *tree.DirectoryNode {
	docsDir := &tree.DirectoryNode{
		Name:        "Docs",
		Path:        "Docs",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "base-template",
		Description: "Project documentation and guides",
	}

	// Add standard documentation files
	standardDocs := []string{"README.md", "API.md", "DEPLOYMENT.md", "DEVELOPMENT.md", "ARCHITECTURE.md"}
	for _, docFile := range standardDocs {
		fileNode := &tree.FileNode{
			Name:        docFile,
			Path:        filepath.Join("Docs", docFile),
			Size:        2048,
			Source:      "base-template",
			Templated:   true,
			Executable:  false,
			Description: fmt.Sprintf("Project %s documentation", strings.TrimSuffix(docFile, ".md")),
		}
		docsDir.Files = append(docsDir.Files, fileNode)
	}

	return docsDir
}

func (pm *PreviewManager) createScriptsDirectory() *tree.DirectoryNode {
	scriptsDir := &tree.DirectoryNode{
		Name:        "Scripts",
		Path:        "Scripts",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "base-template",
		Description: "Build, deployment, and utility scripts",
	}

	// Add standard script files
	standardScripts := []string{"build.sh", "test.sh", "deploy.sh", "setup.sh", "clean.sh"}
	for _, scriptFile := range standardScripts {
		fileNode := &tree.FileNode{
			Name:        scriptFile,
			Path:        filepath.Join("Scripts", scriptFile),
			Size:        1024,
			Source:      "base-template",
			Templated:   true,
			Executable:  true,
			Description: fmt.Sprintf("%s script", strings.TrimSuffix(scriptFile, ".sh")),
		}
		scriptsDir.Files = append(scriptsDir.Files, fileNode)
	}

	return scriptsDir
}

func (pm *PreviewManager) createGitHubDirectory() *tree.DirectoryNode {
	githubDir := &tree.DirectoryNode{
		Name:        ".github",
		Path:        ".github",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "base-template",
		Description: "GitHub workflows and repository configuration",
	}

	// Create workflows subdirectory with standard files
	workflowsDir := &tree.DirectoryNode{
		Name:        "workflows",
		Path:        ".github/workflows",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "base-template",
		Description: "GitHub Actions CI/CD workflows",
	}

	workflows := []string{"ci.yml", "cd.yml", "security.yml", "release.yml"}
	for _, workflow := range workflows {
		fileNode := &tree.FileNode{
			Name:        workflow,
			Path:        filepath.Join(".github/workflows", workflow),
			Size:        3072,
			Source:      "base-template",
			Templated:   true,
			Executable:  false,
			Description: fmt.Sprintf("GitHub Actions %s workflow", strings.TrimSuffix(workflow, ".yml")),
		}
		workflowsDir.Files = append(workflowsDir.Files, fileNode)
	}

	githubDir.Children = append(githubDir.Children, workflowsDir)

	return githubDir
}

func (pm *PreviewManager) addCommonProjectFiles(rootDir *tree.DirectoryNode) {
	commonFiles := []struct {
		name        string
		description string
		size        int64
		executable  bool
	}{
		{"README.md", "Project overview and setup instructions", 4096, false},
		{"CONTRIBUTING.md", "Contribution guidelines", 2048, false},
		{"LICENSE", "Project license", 1024, false},
		{".gitignore", "Git ignore patterns", 512, false},
		{"Makefile", "Build automation", 2048, false},
		{"docker-compose.yml", "Docker services configuration", 3072, false},
	}

	for _, file := range commonFiles {
		fileNode := &tree.FileNode{
			Name:        file.name,
			Path:        file.name,
			Size:        file.size,
			Source:      "base-template",
			Templated:   true,
			Executable:  file.executable,
			Description: file.description,
		}
		rootDir.Files = append(rootDir.Files, fileNode)
	}
}

func (pm *PreviewManager) createComponentMapping(component, directory string, selections []interfaces.TemplateSelection, previews map[string]*interfaces.TemplatePreview, category string) display.ComponentTemplateMapping {
	var templates []string
	var fileCount int
	var estimatedSize int64

	for _, selection := range selections {
		if strings.ToLower(selection.Template.Category) == category {
			templates = append(templates, selection.Template.DisplayName)
			if preview, exists := previews[selection.Template.Name]; exists {
				fileCount += preview.Summary.TotalFiles
				estimatedSize += preview.Summary.TotalSize
			}
		}
	}

	return display.ComponentTemplateMapping{
		Component:     component,
		Directory:     directory,
		Templates:     templates,
		Description:   fmt.Sprintf("%s components generated from %s templates", component, category),
		FileCount:     fileCount,
		EstimatedSize: estimatedSize,
	}
}

func (pm *PreviewManager) addStandardComponentMappings(mappings *[]display.ComponentTemplateMapping, preview *display.ProjectStructurePreview) {
	*mappings = append(*mappings, display.ComponentTemplateMapping{
		Component:     "Documentation",
		Directory:     "Docs/",
		Templates:     []string{"base-template"},
		Description:   "Project documentation and guides",
		FileCount:     len(preview.DirectoryTree.Docs.Files),
		EstimatedSize: pm.treeFormatter.CalculateDirectorySize(preview.DirectoryTree.Docs),
	})

	*mappings = append(*mappings, display.ComponentTemplateMapping{
		Component:     "Build Scripts",
		Directory:     "Scripts/",
		Templates:     []string{"base-template"},
		Description:   "Build, test, and deployment automation",
		FileCount:     len(preview.DirectoryTree.Scripts.Files),
		EstimatedSize: pm.treeFormatter.CalculateDirectorySize(preview.DirectoryTree.Scripts),
	})

	*mappings = append(*mappings, display.ComponentTemplateMapping{
		Component:     "CI/CD Workflows",
		Directory:     ".github/",
		Templates:     []string{"base-template"},
		Description:   "GitHub Actions workflows and repository configuration",
		FileCount:     pm.treeFormatter.CountFilesRecursively(preview.DirectoryTree.GitHub),
		EstimatedSize: pm.treeFormatter.CalculateDirectorySize(preview.DirectoryTree.GitHub),
	})
}

func (pm *PreviewManager) collectProjectMetadata(summary *display.ProjectSummary, selections []interfaces.TemplateSelection) {
	techMap := make(map[string]bool)
	catMap := make(map[string]bool)
	depMap := make(map[string]bool)

	for _, selection := range selections {
		if selection.Template.Technology != "" {
			techMap[selection.Template.Technology] = true
		}
		if selection.Template.Category != "" {
			catMap[selection.Template.Category] = true
		}
		for _, dep := range selection.Template.Dependencies {
			depMap[dep] = true
		}
	}

	// Convert maps to slices and sort
	for tech := range techMap {
		summary.Technologies = append(summary.Technologies, tech)
	}
	for cat := range catMap {
		summary.Categories = append(summary.Categories, cat)
	}
	for dep := range depMap {
		summary.Dependencies = append(summary.Dependencies, dep)
	}

	sort.Strings(summary.Technologies)
	sort.Strings(summary.Categories)
	sort.Strings(summary.Dependencies)
}

func (pm *PreviewManager) checkDependenciesAndConflicts(warnings *[]string, conflicts *[]display.FileConflict, selections []interfaces.TemplateSelection, templatePreviews map[string]*interfaces.TemplatePreview) {
	// Check for missing dependencies
	selectedTemplates := make(map[string]bool)
	for _, selection := range selections {
		selectedTemplates[selection.Template.Name] = true
	}

	for _, selection := range selections {
		for _, dep := range selection.Template.Dependencies {
			if !selectedTemplates[dep] {
				*warnings = append(*warnings, fmt.Sprintf("Template '%s' requires dependency '%s' which is not selected", selection.Template.DisplayName, dep))
			}
		}
	}

	// Check for potential file conflicts (simplified)
	fileMap := make(map[string][]string)
	for templateName, templatePreview := range templatePreviews {
		for _, file := range templatePreview.Files {
			if file.Type == "file" {
				fileMap[file.Path] = append(fileMap[file.Path], templateName)
			}
		}
	}

	for path, sources := range fileMap {
		if len(sources) > 1 {
			severity := "warning"
			message := fmt.Sprintf("File may be created by multiple templates: %s", strings.Join(sources, ", "))

			if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts") {
				severity = "error"
				message = "Code file conflict - manual resolution may be required"
			}

			*conflicts = append(*conflicts, display.FileConflict{
				Path:       path,
				Templates:  sources,
				Severity:   severity,
				Message:    message,
				Resolvable: severity != "error",
			})
		}
	}
}
