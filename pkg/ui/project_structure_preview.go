// Package ui provides project structure preview functionality for interactive CLI generation.
package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectStructurePreviewGenerator handles generation and display of project structure previews
type ProjectStructurePreviewGenerator struct {
	ui              interfaces.InteractiveUIInterface
	templateManager interfaces.TemplateManager
	logger          interfaces.Logger
}

// ProjectStructurePreview represents a complete project structure preview
type ProjectStructurePreview struct {
	ProjectName       string                     `json:"project_name"`
	OutputDirectory   string                     `json:"output_directory"`
	SelectedTemplates []TemplateSelection        `json:"selected_templates"`
	DirectoryTree     *StandardDirectoryTree     `json:"directory_tree"`
	ComponentMapping  []ComponentTemplateMapping `json:"component_mapping"`
	Summary           ProjectSummary             `json:"summary"`
	Warnings          []string                   `json:"warnings"`
	Conflicts         []FileConflict             `json:"conflicts"`
}

// StandardDirectoryTree represents the standardized project directory structure
type StandardDirectoryTree struct {
	Root         *DirectoryNode `json:"root"`
	App          *DirectoryNode `json:"app,omitempty"`
	CommonServer *DirectoryNode `json:"common_server,omitempty"`
	Mobile       *DirectoryNode `json:"mobile,omitempty"`
	Deploy       *DirectoryNode `json:"deploy,omitempty"`
	Docs         *DirectoryNode `json:"docs"`
	Scripts      *DirectoryNode `json:"scripts"`
	GitHub       *DirectoryNode `json:"github"`
}

// ComponentTemplateMapping shows which templates contribute to each component
type ComponentTemplateMapping struct {
	Component     string   `json:"component"`
	Directory     string   `json:"directory"`
	Templates     []string `json:"templates"`
	Description   string   `json:"description"`
	FileCount     int      `json:"file_count"`
	EstimatedSize int64    `json:"estimated_size"`
}

// ProjectSummary provides overall project statistics
type ProjectSummary struct {
	TotalDirectories int      `json:"total_directories"`
	TotalFiles       int      `json:"total_files"`
	EstimatedSize    int64    `json:"estimated_size"`
	Technologies     []string `json:"technologies"`
	Categories       []string `json:"categories"`
	Dependencies     []string `json:"dependencies"`
}

// NewProjectStructurePreviewGenerator creates a new project structure preview generator
func NewProjectStructurePreviewGenerator(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *ProjectStructurePreviewGenerator {
	return &ProjectStructurePreviewGenerator{
		ui:              ui,
		templateManager: templateManager,
		logger:          logger,
	}
}

// GenerateProjectStructurePreview creates a comprehensive project structure preview
func (pspg *ProjectStructurePreviewGenerator) GenerateProjectStructurePreview(ctx context.Context, config *models.ProjectConfig, selections []TemplateSelection, outputDir string) (*ProjectStructurePreview, error) {
	if config == nil {
		return nil, fmt.Errorf("project config cannot be nil")
	}

	if len(selections) == 0 {
		return nil, fmt.Errorf("no templates selected")
	}

	pspg.logger.InfoWithFields("Generating project structure preview", map[string]interface{}{
		"project_name":     config.Name,
		"template_count":   len(selections),
		"output_directory": outputDir,
	})

	// Initialize preview structure
	preview := &ProjectStructurePreview{
		ProjectName:       config.Name,
		OutputDirectory:   outputDir,
		SelectedTemplates: selections,
		DirectoryTree:     &StandardDirectoryTree{},
		ComponentMapping:  []ComponentTemplateMapping{},
		Summary:           ProjectSummary{},
		Warnings:          []string{},
		Conflicts:         []FileConflict{},
	}

	// Generate individual template previews
	templatePreviews, err := pspg.generateTemplatePreviews(ctx, selections, config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template previews: %w", err)
	}

	// Build standardized directory structure
	err = pspg.buildStandardDirectoryStructure(preview, templatePreviews, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build directory structure: %w", err)
	}

	// Generate component mappings
	pspg.generateComponentMappings(preview, templatePreviews)

	// Calculate project summary
	pspg.calculateProjectSummary(preview, templatePreviews)

	// Detect conflicts and generate warnings
	pspg.detectConflictsAndWarnings(preview, templatePreviews)

	return preview, nil
}

// DisplayProjectStructurePreview shows the project structure preview to the user
func (pspg *ProjectStructurePreviewGenerator) DisplayProjectStructurePreview(ctx context.Context, preview *ProjectStructurePreview) error {
	if preview == nil {
		return fmt.Errorf("preview cannot be nil")
	}

	// Display project header
	pspg.displayProjectHeader(preview)

	// Display directory structure tree
	err := pspg.displayDirectoryTree(ctx, preview)
	if err != nil {
		return fmt.Errorf("failed to display directory tree: %w", err)
	}

	// Display component mappings
	pspg.displayComponentMappings(preview)

	// Display project summary
	pspg.displayProjectSummary(preview)

	// Display warnings and conflicts
	pspg.displayWarningsAndConflicts(preview)

	return nil
}

// generateTemplatePreviews generates previews for all selected templates
func (pspg *ProjectStructurePreviewGenerator) generateTemplatePreviews(ctx context.Context, selections []TemplateSelection, config *models.ProjectConfig) (map[string]*interfaces.TemplatePreview, error) {
	previews := make(map[string]*interfaces.TemplatePreview)

	for _, selection := range selections {
		preview, err := pspg.templateManager.PreviewTemplate(selection.Template.Name, config)
		if err != nil {
			pspg.logger.WarnWithFields("Failed to preview template", map[string]interface{}{
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
func (pspg *ProjectStructurePreviewGenerator) buildStandardDirectoryStructure(preview *ProjectStructurePreview, templatePreviews map[string]*interfaces.TemplatePreview, config *models.ProjectConfig) error {
	// Create root directory
	preview.DirectoryTree.Root = &DirectoryNode{
		Name:        config.Name,
		Path:        "",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "project-root",
		Description: "Project root directory",
	}

	// Analyze selected templates to determine which directories to create
	templatesByCategory := pspg.categorizeTemplates(preview.SelectedTemplates)

	// Create App/ directory for frontend templates
	if frontendTemplates, exists := templatesByCategory["frontend"]; exists && len(frontendTemplates) > 0 {
		appDir := pspg.createAppDirectory(frontendTemplates, templatePreviews)
		preview.DirectoryTree.App = appDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, appDir)
	}

	// Create CommonServer/ directory for backend templates
	if backendTemplates, exists := templatesByCategory["backend"]; exists && len(backendTemplates) > 0 {
		serverDir := pspg.createCommonServerDirectory(backendTemplates, templatePreviews)
		preview.DirectoryTree.CommonServer = serverDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, serverDir)
	}

	// Create Mobile/ directory for mobile templates
	if mobileTemplates, exists := templatesByCategory["mobile"]; exists && len(mobileTemplates) > 0 {
		mobileDir := pspg.createMobileDirectory(mobileTemplates, templatePreviews)
		preview.DirectoryTree.Mobile = mobileDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, mobileDir)
	}

	// Create Deploy/ directory for infrastructure templates
	if infraTemplates, exists := templatesByCategory["infrastructure"]; exists && len(infraTemplates) > 0 {
		deployDir := pspg.createDeployDirectory(infraTemplates, templatePreviews)
		preview.DirectoryTree.Deploy = deployDir
		preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, deployDir)
	}

	// Always create standard directories
	docsDir := pspg.createDocsDirectory(templatePreviews)
	preview.DirectoryTree.Docs = docsDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, docsDir)

	scriptsDir := pspg.createScriptsDirectory(templatePreviews)
	preview.DirectoryTree.Scripts = scriptsDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, scriptsDir)

	githubDir := pspg.createGitHubDirectory(templatePreviews)
	preview.DirectoryTree.GitHub = githubDir
	preview.DirectoryTree.Root.Children = append(preview.DirectoryTree.Root.Children, githubDir)

	// Add common project files to root
	pspg.addCommonProjectFiles(preview.DirectoryTree.Root, templatePreviews)

	// Sort all directory children
	pspg.sortDirectoryTree(preview.DirectoryTree.Root)

	return nil
}

// categorizeTemplates groups templates by their category
func (pspg *ProjectStructurePreviewGenerator) categorizeTemplates(selections []TemplateSelection) map[string][]TemplateSelection {
	categories := make(map[string][]TemplateSelection)

	for _, selection := range selections {
		category := strings.ToLower(selection.Template.Category)
		categories[category] = append(categories[category], selection)
	}

	return categories
}

// createAppDirectory creates the App/ directory structure for frontend templates
func (pspg *ProjectStructurePreviewGenerator) createAppDirectory(templates []TemplateSelection, previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	appDir := &DirectoryNode{
		Name:        "App",
		Path:        "App",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "frontend-templates",
		Description: "Frontend applications and components",
	}

	// Create standard frontend subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]TemplateSelection) bool
	}{
		{"main", "Main application frontend", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "nextjs-app") }},
		{"home", "Home/landing page frontend", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "nextjs-home") }},
		{"admin", "Admin dashboard frontend", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "nextjs-admin") }},
		{"shared-components", "Shared UI components", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "shared-components") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("App", subdir.name),
				Children:    []*DirectoryNode{},
				Files:       []*FileNode{},
				Source:      "frontend-template",
				Description: subdir.description,
			}
			pspg.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			appDir.Children = append(appDir.Children, childDir)
		}
	}

	return appDir
}

// createCommonServerDirectory creates the CommonServer/ directory structure for backend templates
func (pspg *ProjectStructurePreviewGenerator) createCommonServerDirectory(templates []TemplateSelection, previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	serverDir := &DirectoryNode{
		Name:        "CommonServer",
		Path:        "CommonServer",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "backend-templates",
		Description: "Backend API and server components",
	}

	// Create standard Go project structure
	subdirs := []struct {
		name        string
		description string
	}{
		{"cmd", "Command-line applications and entry points"},
		{"internal", "Internal application code"},
		{"pkg", "Public library code"},
		{"migrations", "Database migration files"},
		{"docs", "API documentation and specifications"},
	}

	for _, subdir := range subdirs {
		childDir := &DirectoryNode{
			Name:        subdir.name,
			Path:        filepath.Join("CommonServer", subdir.name),
			Children:    []*DirectoryNode{},
			Files:       []*FileNode{},
			Source:      "backend-template",
			Description: subdir.description,
		}
		pspg.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
		serverDir.Children = append(serverDir.Children, childDir)
	}

	return serverDir
}

// createMobileDirectory creates the Mobile/ directory structure for mobile templates
func (pspg *ProjectStructurePreviewGenerator) createMobileDirectory(templates []TemplateSelection, previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	mobileDir := &DirectoryNode{
		Name:        "Mobile",
		Path:        "Mobile",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "mobile-templates",
		Description: "Mobile applications for iOS and Android",
	}

	// Create platform-specific subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]TemplateSelection) bool
	}{
		{"android", "Android application (Kotlin)", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "android-kotlin") }},
		{"ios", "iOS application (Swift)", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "ios-swift") }},
		{"shared", "Shared mobile resources and APIs", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "shared") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("Mobile", subdir.name),
				Children:    []*DirectoryNode{},
				Files:       []*FileNode{},
				Source:      "mobile-template",
				Description: subdir.description,
			}
			pspg.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			mobileDir.Children = append(mobileDir.Children, childDir)
		}
	}

	return mobileDir
}

// createDeployDirectory creates the Deploy/ directory structure for infrastructure templates
func (pspg *ProjectStructurePreviewGenerator) createDeployDirectory(templates []TemplateSelection, previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	deployDir := &DirectoryNode{
		Name:        "Deploy",
		Path:        "Deploy",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "infrastructure-templates",
		Description: "Deployment and infrastructure configuration",
	}

	// Create infrastructure subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]TemplateSelection) bool
	}{
		{"docker", "Docker containers and compose files", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "docker") }},
		{"k8s", "Kubernetes manifests and configurations", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "kubernetes") }},
		{"terraform", "Infrastructure as Code with Terraform", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "terraform") }},
		{"monitoring", "Monitoring and observability tools", func(t []TemplateSelection) bool { return pspg.hasTemplateType(t, "monitoring") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("Deploy", subdir.name),
				Children:    []*DirectoryNode{},
				Files:       []*FileNode{},
				Source:      "infrastructure-template",
				Description: subdir.description,
			}
			pspg.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			deployDir.Children = append(deployDir.Children, childDir)
		}
	}

	return deployDir
}

// createDocsDirectory creates the Docs/ directory (always present)
func (pspg *ProjectStructurePreviewGenerator) createDocsDirectory(previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	docsDir := &DirectoryNode{
		Name:        "Docs",
		Path:        "Docs",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "base-template",
		Description: "Project documentation and guides",
	}

	// Add standard documentation files
	standardDocs := []string{
		"README.md",
		"API.md",
		"DEPLOYMENT.md",
		"DEVELOPMENT.md",
		"ARCHITECTURE.md",
	}

	for _, docFile := range standardDocs {
		fileNode := &FileNode{
			Name:        docFile,
			Path:        filepath.Join("Docs", docFile),
			Size:        2048, // Estimated size
			Source:      "base-template",
			Templated:   true,
			Executable:  false,
			Description: fmt.Sprintf("Project %s documentation", strings.TrimSuffix(docFile, ".md")),
		}
		docsDir.Files = append(docsDir.Files, fileNode)
	}

	return docsDir
}

// createScriptsDirectory creates the Scripts/ directory (always present)
func (pspg *ProjectStructurePreviewGenerator) createScriptsDirectory(previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	scriptsDir := &DirectoryNode{
		Name:        "Scripts",
		Path:        "Scripts",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "base-template",
		Description: "Build, deployment, and utility scripts",
	}

	// Add standard script files
	standardScripts := []string{
		"build.sh",
		"test.sh",
		"deploy.sh",
		"setup.sh",
		"clean.sh",
	}

	for _, scriptFile := range standardScripts {
		fileNode := &FileNode{
			Name:        scriptFile,
			Path:        filepath.Join("Scripts", scriptFile),
			Size:        1024, // Estimated size
			Source:      "base-template",
			Templated:   true,
			Executable:  true,
			Description: fmt.Sprintf("%s script", strings.TrimSuffix(scriptFile, ".sh")),
		}
		scriptsDir.Files = append(scriptsDir.Files, fileNode)
	}

	return scriptsDir
}

// createGitHubDirectory creates the .github/ directory (always present)
func (pspg *ProjectStructurePreviewGenerator) createGitHubDirectory(previews map[string]*interfaces.TemplatePreview) *DirectoryNode {
	githubDir := &DirectoryNode{
		Name:        ".github",
		Path:        ".github",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "base-template",
		Description: "GitHub workflows and repository configuration",
	}

	// Create workflows subdirectory
	workflowsDir := &DirectoryNode{
		Name:        "workflows",
		Path:        ".github/workflows",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "base-template",
		Description: "GitHub Actions CI/CD workflows",
	}

	// Add standard workflow files
	workflows := []string{
		"ci.yml",
		"cd.yml",
		"security.yml",
		"release.yml",
	}

	for _, workflow := range workflows {
		fileNode := &FileNode{
			Name:        workflow,
			Path:        filepath.Join(".github/workflows", workflow),
			Size:        3072, // Estimated size
			Source:      "base-template",
			Templated:   true,
			Executable:  false,
			Description: fmt.Sprintf("GitHub Actions %s workflow", strings.TrimSuffix(workflow, ".yml")),
		}
		workflowsDir.Files = append(workflowsDir.Files, fileNode)
	}

	githubDir.Children = append(githubDir.Children, workflowsDir)

	// Add other GitHub files
	githubFiles := []string{
		"PULL_REQUEST_TEMPLATE.md",
		"CODEOWNERS",
		"dependabot.yml",
	}

	for _, ghFile := range githubFiles {
		fileNode := &FileNode{
			Name:        ghFile,
			Path:        filepath.Join(".github", ghFile),
			Size:        1024, // Estimated size
			Source:      "base-template",
			Templated:   true,
			Executable:  false,
			Description: fmt.Sprintf("GitHub %s configuration", ghFile),
		}
		githubDir.Files = append(githubDir.Files, fileNode)
	}

	return githubDir
}

// Helper methods continue in the next part...
// addCommonProjectFiles adds common project files to the root directory
func (pspg *ProjectStructurePreviewGenerator) addCommonProjectFiles(rootDir *DirectoryNode, previews map[string]*interfaces.TemplatePreview) {
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
		fileNode := &FileNode{
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

// hasTemplateType checks if any template matches the given type
func (pspg *ProjectStructurePreviewGenerator) hasTemplateType(templates []TemplateSelection, templateType string) bool {
	for _, template := range templates {
		if strings.Contains(strings.ToLower(template.Template.Name), templateType) {
			return true
		}
	}
	return false
}

// populateDirectoryFromTemplates populates a directory with files from matching templates
func (pspg *ProjectStructurePreviewGenerator) populateDirectoryFromTemplates(dir *DirectoryNode, templates []TemplateSelection, previews map[string]*interfaces.TemplatePreview, targetSubdir string) {
	for _, template := range templates {
		if preview, exists := previews[template.Template.Name]; exists {
			for _, file := range preview.Files {
				// Check if file belongs to this subdirectory
				if strings.HasPrefix(file.Path, targetSubdir+"/") ||
					(targetSubdir == "" && !strings.Contains(file.Path, "/")) {

					relativePath := strings.TrimPrefix(file.Path, targetSubdir+"/")
					if relativePath == "" {
						relativePath = file.Path
					}

					if file.Type == "directory" {
						// Add subdirectory if not exists
						if !pspg.hasChildDirectory(dir, filepath.Base(relativePath)) {
							childDir := &DirectoryNode{
								Name:        filepath.Base(relativePath),
								Path:        filepath.Join(dir.Path, filepath.Base(relativePath)),
								Children:    []*DirectoryNode{},
								Files:       []*FileNode{},
								Source:      template.Template.Name,
								Description: fmt.Sprintf("Generated from %s template", template.Template.DisplayName),
							}
							dir.Children = append(dir.Children, childDir)
						}
					} else {
						// Add file
						fileNode := &FileNode{
							Name:        filepath.Base(relativePath),
							Path:        filepath.Join(dir.Path, relativePath),
							Size:        file.Size,
							Source:      template.Template.Name,
							Templated:   file.Templated,
							Executable:  file.Executable,
							Description: fmt.Sprintf("Generated from %s template", template.Template.DisplayName),
						}
						dir.Files = append(dir.Files, fileNode)
					}
				}
			}
		}
	}
}

// hasChildDirectory checks if a directory already has a child with the given name
func (pspg *ProjectStructurePreviewGenerator) hasChildDirectory(parent *DirectoryNode, name string) bool {
	for _, child := range parent.Children {
		if child.Name == name {
			return true
		}
	}
	return false
}

// sortDirectoryTree recursively sorts directory contents
func (pspg *ProjectStructurePreviewGenerator) sortDirectoryTree(node *DirectoryNode) {
	// Sort directories first, then files
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Name < node.Children[j].Name
	})

	sort.Slice(node.Files, func(i, j int) bool {
		return node.Files[i].Name < node.Files[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		pspg.sortDirectoryTree(child)
	}
}

// generateComponentMappings creates mappings between components and templates
func (pspg *ProjectStructurePreviewGenerator) generateComponentMappings(preview *ProjectStructurePreview, templatePreviews map[string]*interfaces.TemplatePreview) {
	mappings := []ComponentTemplateMapping{}

	// Map each major directory to its contributing templates
	if preview.DirectoryTree.App != nil {
		mapping := pspg.createComponentMapping("Frontend Applications", "App/", preview.SelectedTemplates, templatePreviews, "frontend")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.CommonServer != nil {
		mapping := pspg.createComponentMapping("Backend API", "CommonServer/", preview.SelectedTemplates, templatePreviews, "backend")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.Mobile != nil {
		mapping := pspg.createComponentMapping("Mobile Applications", "Mobile/", preview.SelectedTemplates, templatePreviews, "mobile")
		mappings = append(mappings, mapping)
	}

	if preview.DirectoryTree.Deploy != nil {
		mapping := pspg.createComponentMapping("Infrastructure & Deployment", "Deploy/", preview.SelectedTemplates, templatePreviews, "infrastructure")
		mappings = append(mappings, mapping)
	}

	// Always include standard components
	mappings = append(mappings, ComponentTemplateMapping{
		Component:     "Documentation",
		Directory:     "Docs/",
		Templates:     []string{"base-template"},
		Description:   "Project documentation and guides",
		FileCount:     len(preview.DirectoryTree.Docs.Files),
		EstimatedSize: pspg.calculateDirectorySize(preview.DirectoryTree.Docs),
	})

	mappings = append(mappings, ComponentTemplateMapping{
		Component:     "Build Scripts",
		Directory:     "Scripts/",
		Templates:     []string{"base-template"},
		Description:   "Build, test, and deployment automation",
		FileCount:     len(preview.DirectoryTree.Scripts.Files),
		EstimatedSize: pspg.calculateDirectorySize(preview.DirectoryTree.Scripts),
	})

	mappings = append(mappings, ComponentTemplateMapping{
		Component:     "CI/CD Workflows",
		Directory:     ".github/",
		Templates:     []string{"base-template"},
		Description:   "GitHub Actions workflows and repository configuration",
		FileCount:     pspg.countFilesRecursively(preview.DirectoryTree.GitHub),
		EstimatedSize: pspg.calculateDirectorySize(preview.DirectoryTree.GitHub),
	})

	preview.ComponentMapping = mappings
}

// createComponentMapping creates a component mapping for a specific category
func (pspg *ProjectStructurePreviewGenerator) createComponentMapping(component, directory string, selections []TemplateSelection, previews map[string]*interfaces.TemplatePreview, category string) ComponentTemplateMapping {
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

	description := fmt.Sprintf("%s components generated from %s templates", component, category)

	return ComponentTemplateMapping{
		Component:     component,
		Directory:     directory,
		Templates:     templates,
		Description:   description,
		FileCount:     fileCount,
		EstimatedSize: estimatedSize,
	}
}

// calculateProjectSummary calculates overall project statistics
func (pspg *ProjectStructurePreviewGenerator) calculateProjectSummary(preview *ProjectStructurePreview, templatePreviews map[string]*interfaces.TemplatePreview) {
	summary := ProjectSummary{
		Technologies: []string{},
		Categories:   []string{},
		Dependencies: []string{},
	}

	// Count directories and files recursively
	summary.TotalDirectories = pspg.countDirectoriesRecursively(preview.DirectoryTree.Root)
	summary.TotalFiles = pspg.countFilesRecursively(preview.DirectoryTree.Root)
	summary.EstimatedSize = pspg.calculateDirectorySize(preview.DirectoryTree.Root)

	// Collect unique technologies, categories, and dependencies
	techMap := make(map[string]bool)
	catMap := make(map[string]bool)
	depMap := make(map[string]bool)

	for _, selection := range preview.SelectedTemplates {
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

	// Convert maps to slices
	for tech := range techMap {
		summary.Technologies = append(summary.Technologies, tech)
	}
	for cat := range catMap {
		summary.Categories = append(summary.Categories, cat)
	}
	for dep := range depMap {
		summary.Dependencies = append(summary.Dependencies, dep)
	}

	// Sort for consistent output
	sort.Strings(summary.Technologies)
	sort.Strings(summary.Categories)
	sort.Strings(summary.Dependencies)

	preview.Summary = summary
}

// detectConflictsAndWarnings identifies potential issues with the project structure
func (pspg *ProjectStructurePreviewGenerator) detectConflictsAndWarnings(preview *ProjectStructurePreview, templatePreviews map[string]*interfaces.TemplatePreview) {
	var warnings []string
	var conflicts []FileConflict

	// Check for large project size
	if preview.Summary.EstimatedSize > 500*1024*1024 { // 500MB
		warnings = append(warnings, fmt.Sprintf("Large project size: %s", formatBytes(preview.Summary.EstimatedSize)))
	}

	// Check for many files
	if preview.Summary.TotalFiles > 2000 {
		warnings = append(warnings, fmt.Sprintf("Large number of files: %d", preview.Summary.TotalFiles))
	}

	// Check for missing dependencies
	selectedTemplates := make(map[string]bool)
	for _, selection := range preview.SelectedTemplates {
		selectedTemplates[selection.Template.Name] = true
	}

	for _, selection := range preview.SelectedTemplates {
		for _, dep := range selection.Template.Dependencies {
			if !selectedTemplates[dep] {
				warnings = append(warnings, fmt.Sprintf("Template '%s' requires dependency '%s' which is not selected", selection.Template.DisplayName, dep))
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

			conflicts = append(conflicts, FileConflict{
				Path:       path,
				Templates:  sources,
				Severity:   severity,
				Message:    message,
				Resolvable: severity != "error",
			})
		}
	}

	preview.Warnings = warnings
	preview.Conflicts = conflicts
}

// Display methods

// displayProjectHeader shows the project header information
func (pspg *ProjectStructurePreviewGenerator) displayProjectHeader(preview *ProjectStructurePreview) {
	fmt.Printf("\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("                    PROJECT STRUCTURE PREVIEW\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("\n")
	fmt.Printf("Project Name:     %s\n", preview.ProjectName)
	fmt.Printf("Output Directory: %s\n", preview.OutputDirectory)
	fmt.Printf("Templates:        %d selected\n", len(preview.SelectedTemplates))
	fmt.Printf("\n")
}

// displayDirectoryTree shows the directory structure as a tree
func (pspg *ProjectStructurePreviewGenerator) displayDirectoryTree(ctx context.Context, preview *ProjectStructurePreview) error {
	if preview.DirectoryTree.Root == nil {
		return fmt.Errorf("no directory tree to display")
	}

	treeConfig := interfaces.TreeConfig{
		Title:      "Project Directory Structure",
		Root:       pspg.convertToTreeNode(preview.DirectoryTree.Root),
		Expandable: true,
		ShowIcons:  true,
		MaxDepth:   10,
	}

	return pspg.ui.ShowTree(ctx, treeConfig)
}

// convertToTreeNode converts DirectoryNode to TreeNode for display
func (pspg *ProjectStructurePreviewGenerator) convertToTreeNode(dir *DirectoryNode) interfaces.TreeNode {
	node := interfaces.TreeNode{
		Label:      dir.Name,
		Icon:       "📁",
		Expanded:   true,
		Selectable: false,
		Children:   []interfaces.TreeNode{},
		Metadata: map[string]interface{}{
			"path":        dir.Path,
			"source":      dir.Source,
			"description": dir.Description,
		},
	}

	// Add child directories
	for _, child := range dir.Children {
		childNode := pspg.convertToTreeNode(child)
		node.Children = append(node.Children, childNode)
	}

	// Add files
	for _, file := range dir.Files {
		icon := "📄"
		if file.Executable {
			icon = "⚙️"
		} else if file.Templated {
			icon = "📝"
		}

		fileNode := interfaces.TreeNode{
			Label:      file.Name,
			Icon:       icon,
			Selectable: false,
			Metadata: map[string]interface{}{
				"path":        file.Path,
				"size":        file.Size,
				"source":      file.Source,
				"templated":   file.Templated,
				"executable":  file.Executable,
				"description": file.Description,
			},
		}
		node.Children = append(node.Children, fileNode)
	}

	return node
}

// displayComponentMappings shows which templates contribute to each component
func (pspg *ProjectStructurePreviewGenerator) displayComponentMappings(preview *ProjectStructurePreview) {
	fmt.Printf("\n")
	fmt.Printf("Component-Template Mapping:\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n")

	for _, mapping := range preview.ComponentMapping {
		fmt.Printf("\n📦 %s (%s)\n", mapping.Component, mapping.Directory)
		fmt.Printf("   Description: %s\n", mapping.Description)
		fmt.Printf("   Files: %d | Size: %s\n", mapping.FileCount, formatBytes(mapping.EstimatedSize))
		fmt.Printf("   Templates: %s\n", strings.Join(mapping.Templates, ", "))
	}
}

// displayProjectSummary shows overall project statistics
func (pspg *ProjectStructurePreviewGenerator) displayProjectSummary(preview *ProjectStructurePreview) {
	fmt.Printf("\n")
	fmt.Printf("Project Summary:\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n")
	fmt.Printf("Directories:   %d\n", preview.Summary.TotalDirectories)
	fmt.Printf("Files:         %d\n", preview.Summary.TotalFiles)
	fmt.Printf("Estimated Size: %s\n", formatBytes(preview.Summary.EstimatedSize))

	if len(preview.Summary.Technologies) > 0 {
		fmt.Printf("Technologies:  %s\n", strings.Join(preview.Summary.Technologies, ", "))
	}

	if len(preview.Summary.Categories) > 0 {
		fmt.Printf("Categories:    %s\n", strings.Join(preview.Summary.Categories, ", "))
	}

	if len(preview.Summary.Dependencies) > 0 {
		fmt.Printf("Dependencies:  %s\n", strings.Join(preview.Summary.Dependencies, ", "))
	}
}

// displayWarningsAndConflicts shows warnings and conflicts
func (pspg *ProjectStructurePreviewGenerator) displayWarningsAndConflicts(preview *ProjectStructurePreview) {
	if len(preview.Warnings) > 0 {
		fmt.Printf("\n")
		fmt.Printf("⚠️  Warnings:\n")
		fmt.Printf("───────────────────────────────────────────────────────────────\n")
		for _, warning := range preview.Warnings {
			fmt.Printf("   • %s\n", warning)
		}
	}

	if len(preview.Conflicts) > 0 {
		fmt.Printf("\n")
		fmt.Printf("⚠️  File Conflicts:\n")
		fmt.Printf("───────────────────────────────────────────────────────────────\n")
		for _, conflict := range preview.Conflicts {
			icon := "⚠️"
			if conflict.Severity == "error" {
				icon = "❌"
			}
			fmt.Printf("   %s %s: %s\n", icon, conflict.Path, conflict.Message)
		}
	}

	if len(preview.Warnings) == 0 && len(preview.Conflicts) == 0 {
		fmt.Printf("\n")
		fmt.Printf("✅ No warnings or conflicts detected.\n")
	}
}

// Utility methods

// countDirectoriesRecursively counts directories in a tree
func (pspg *ProjectStructurePreviewGenerator) countDirectoriesRecursively(dir *DirectoryNode) int {
	if dir == nil {
		return 0
	}

	count := 1 // Count this directory
	for _, child := range dir.Children {
		count += pspg.countDirectoriesRecursively(child)
	}
	return count
}

// countFilesRecursively counts files in a tree
func (pspg *ProjectStructurePreviewGenerator) countFilesRecursively(dir *DirectoryNode) int {
	if dir == nil {
		return 0
	}

	count := len(dir.Files)
	for _, child := range dir.Children {
		count += pspg.countFilesRecursively(child)
	}
	return count
}

// calculateDirectorySize calculates total size of files in a directory tree
func (pspg *ProjectStructurePreviewGenerator) calculateDirectorySize(dir *DirectoryNode) int64 {
	if dir == nil {
		return 0
	}

	var size int64
	for _, file := range dir.Files {
		size += file.Size
	}
	for _, child := range dir.Children {
		size += pspg.calculateDirectorySize(child)
	}
	return size
}
