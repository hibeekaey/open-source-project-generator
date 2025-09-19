// Package ui provides template preview and combination logic for interactive CLI generation.
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

// TemplatePreviewManager handles template preview and combination logic
type TemplatePreviewManager struct {
	ui              interfaces.InteractiveUIInterface
	templateManager interfaces.TemplateManager
	logger          interfaces.Logger
}

// CombinedTemplatePreview represents a preview of multiple templates combined
type CombinedTemplatePreview struct {
	Templates        []TemplateSelection
	ProjectStructure *ProjectStructure
	FileConflicts    []FileConflict
	Dependencies     []string
	EstimatedSize    int64
	TotalFiles       int
	Warnings         []string
}

// ProjectStructure represents the combined project directory structure
type ProjectStructure struct {
	Root        *DirectoryNode
	Directories map[string]*DirectoryNode
	Files       map[string]*FileNode
}

// DirectoryNode represents a directory in the project structure
type DirectoryNode struct {
	Name        string
	Path        string
	Children    []*DirectoryNode
	Files       []*FileNode
	Source      string // which template created this directory
	Description string
}

// FileNode represents a file in the project structure
type FileNode struct {
	Name        string
	Path        string
	Size        int64
	Source      string // which template created this file
	Templated   bool   // whether this is a template file
	Executable  bool
	Description string
}

// FileConflict represents a conflict between templates
type FileConflict struct {
	Path       string
	Templates  []string
	Severity   string // "error", "warning", "info"
	Message    string
	Resolvable bool
}

// NewTemplatePreviewManager creates a new template preview manager
func NewTemplatePreviewManager(ui interfaces.InteractiveUIInterface, templateManager interfaces.TemplateManager, logger interfaces.Logger) *TemplatePreviewManager {
	return &TemplatePreviewManager{
		ui:              ui,
		templateManager: templateManager,
		logger:          logger,
	}
}

// PreviewIndividualTemplate shows a preview of a single template
func (tpm *TemplatePreviewManager) PreviewIndividualTemplate(ctx context.Context, template interfaces.TemplateInfo, config *models.ProjectConfig) error {
	// Get template preview from template manager
	preview, err := tpm.templateManager.PreviewTemplate(template.Name, config)
	if err != nil {
		return fmt.Errorf("failed to generate template preview: %w", err)
	}

	// Display preview information
	return tpm.displayTemplatePreview(ctx, template, preview)
}

// PreviewCombinedTemplates shows a preview of multiple templates combined
func (tpm *TemplatePreviewManager) PreviewCombinedTemplates(ctx context.Context, selections []TemplateSelection, config *models.ProjectConfig) (*CombinedTemplatePreview, error) {
	if len(selections) == 0 {
		return nil, fmt.Errorf("no templates selected for preview")
	}

	// Generate individual previews
	var individualPreviews []*interfaces.TemplatePreview
	for _, selection := range selections {
		preview, err := tpm.templateManager.PreviewTemplate(selection.Template.Name, config)
		if err != nil {
			tpm.logger.WarnWithFields("Failed to preview template", map[string]interface{}{
				"template": selection.Template.Name,
				"error":    err.Error(),
			})
			continue
		}
		individualPreviews = append(individualPreviews, preview)
	}

	// Combine previews
	combinedPreview, err := tpm.combineTemplatePreviews(selections, individualPreviews)
	if err != nil {
		return nil, fmt.Errorf("failed to combine template previews: %w", err)
	}

	// Display combined preview
	err = tpm.displayCombinedPreview(ctx, combinedPreview)
	if err != nil {
		return nil, fmt.Errorf("failed to display combined preview: %w", err)
	}

	return combinedPreview, nil
}

// combineTemplatePreviews combines multiple template previews into one
func (tpm *TemplatePreviewManager) combineTemplatePreviews(selections []TemplateSelection, previews []*interfaces.TemplatePreview) (*CombinedTemplatePreview, error) {
	combined := &CombinedTemplatePreview{
		Templates: selections,
		ProjectStructure: &ProjectStructure{
			Directories: make(map[string]*DirectoryNode),
			Files:       make(map[string]*FileNode),
		},
		FileConflicts: []FileConflict{},
		Dependencies:  []string{},
		Warnings:      []string{},
	}

	// Track file paths to detect conflicts
	fileSources := make(map[string][]string)

	// Process each template preview
	for i, preview := range previews {
		templateName := selections[i].Template.Name

		// Add to total counts
		combined.EstimatedSize += preview.Summary.TotalSize
		combined.TotalFiles += preview.Summary.TotalFiles

		// Collect dependencies
		for _, dep := range selections[i].Template.Dependencies {
			if !contains(combined.Dependencies, dep) {
				combined.Dependencies = append(combined.Dependencies, dep)
			}
		}

		// Process files and directories
		for _, file := range preview.Files {
			normalizedPath := filepath.Clean(file.Path)

			// Track which templates create this file
			fileSources[normalizedPath] = append(fileSources[normalizedPath], templateName)

			if file.Type == "directory" {
				tpm.addDirectoryToStructure(combined.ProjectStructure, normalizedPath, templateName, "")
			} else {
				tpm.addFileToStructure(combined.ProjectStructure, normalizedPath, templateName, file)
			}
		}
	}

	// Detect and analyze conflicts
	combined.FileConflicts = tpm.detectFileConflicts(fileSources)

	// Generate warnings
	combined.Warnings = tpm.generateWarnings(combined)

	// Create root directory structure
	combined.ProjectStructure.Root = tpm.buildDirectoryTree(combined.ProjectStructure)

	return combined, nil
}

// addDirectoryToStructure adds a directory to the project structure
func (tpm *TemplatePreviewManager) addDirectoryToStructure(structure *ProjectStructure, path, source, description string) {
	if _, exists := structure.Directories[path]; !exists {
		structure.Directories[path] = &DirectoryNode{
			Name:        filepath.Base(path),
			Path:        path,
			Source:      source,
			Description: description,
			Children:    []*DirectoryNode{},
			Files:       []*FileNode{},
		}
	}
}

// addFileToStructure adds a file to the project structure
func (tpm *TemplatePreviewManager) addFileToStructure(structure *ProjectStructure, path, source string, previewFile interfaces.TemplatePreviewFile) {
	if _, exists := structure.Files[path]; !exists {
		structure.Files[path] = &FileNode{
			Name:        filepath.Base(path),
			Path:        path,
			Size:        previewFile.Size,
			Source:      source,
			Templated:   previewFile.Templated,
			Executable:  previewFile.Executable,
			Description: "",
		}
	}
}

// detectFileConflicts identifies conflicts between templates
func (tpm *TemplatePreviewManager) detectFileConflicts(fileSources map[string][]string) []FileConflict {
	var conflicts []FileConflict

	for path, sources := range fileSources {
		if len(sources) > 1 {
			severity := "warning"
			resolvable := true
			message := fmt.Sprintf("File created by multiple templates: %s", strings.Join(sources, ", "))

			// Determine severity based on file type
			if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts") {
				severity = "error"
				resolvable = false
				message = "Code file conflict - manual resolution required"
			} else if strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".txt") {
				severity = "info"
				message = "Documentation file will be merged"
			}

			conflicts = append(conflicts, FileConflict{
				Path:       path,
				Templates:  sources,
				Severity:   severity,
				Message:    message,
				Resolvable: resolvable,
			})
		}
	}

	// Sort conflicts by severity (errors first)
	sort.Slice(conflicts, func(i, j int) bool {
		severityOrder := map[string]int{"error": 0, "warning": 1, "info": 2}
		return severityOrder[conflicts[i].Severity] < severityOrder[conflicts[j].Severity]
	})

	return conflicts
}

// generateWarnings generates warnings for the combined preview
func (tpm *TemplatePreviewManager) generateWarnings(combined *CombinedTemplatePreview) []string {
	var warnings []string

	// Check for missing dependencies
	selectedTemplates := make(map[string]bool)
	for _, selection := range combined.Templates {
		selectedTemplates[selection.Template.Name] = true
	}

	for _, dep := range combined.Dependencies {
		if !selectedTemplates[dep] {
			warnings = append(warnings, fmt.Sprintf("Dependency '%s' is required but not selected", dep))
		}
	}

	// Check for large project size
	if combined.EstimatedSize > 100*1024*1024 { // 100MB
		warnings = append(warnings, fmt.Sprintf("Large project size: %s", formatBytes(combined.EstimatedSize)))
	}

	// Check for many files
	if combined.TotalFiles > 1000 {
		warnings = append(warnings, fmt.Sprintf("Large number of files: %d", combined.TotalFiles))
	}

	// Check for conflicts
	errorConflicts := 0
	for _, conflict := range combined.FileConflicts {
		if conflict.Severity == "error" {
			errorConflicts++
		}
	}
	if errorConflicts > 0 {
		warnings = append(warnings, fmt.Sprintf("%d file conflicts require manual resolution", errorConflicts))
	}

	return warnings
}

// buildDirectoryTree builds a hierarchical directory tree
func (tpm *TemplatePreviewManager) buildDirectoryTree(structure *ProjectStructure) *DirectoryNode {
	root := &DirectoryNode{
		Name:     "project",
		Path:     "",
		Children: []*DirectoryNode{},
		Files:    []*FileNode{},
	}

	// Build parent-child relationships
	for path, dir := range structure.Directories {
		if path == "" || path == "." {
			continue
		}

		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = ""
		}

		if parentPath == "" {
			root.Children = append(root.Children, dir)
		} else if parent, exists := structure.Directories[parentPath]; exists {
			parent.Children = append(parent.Children, dir)
		}
	}

	// Add files to their parent directories
	for path, file := range structure.Files {
		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = ""
		}

		if parentPath == "" {
			root.Files = append(root.Files, file)
		} else if parent, exists := structure.Directories[parentPath]; exists {
			parent.Files = append(parent.Files, file)
		}
	}

	// Sort children and files
	tpm.sortDirectoryTree(root)

	return root
}

// sortDirectoryTree recursively sorts directory contents
func (tpm *TemplatePreviewManager) sortDirectoryTree(node *DirectoryNode) {
	// Sort directories
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Name < node.Children[j].Name
	})

	// Sort files
	sort.Slice(node.Files, func(i, j int) bool {
		return node.Files[i].Name < node.Files[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		tpm.sortDirectoryTree(child)
	}
}

// displayTemplatePreview displays a single template preview
func (tpm *TemplatePreviewManager) displayTemplatePreview(ctx context.Context, template interfaces.TemplateInfo, preview *interfaces.TemplatePreview) error {
	// Build tree structure for display
	treeConfig := interfaces.TreeConfig{
		Title:     fmt.Sprintf("Template Preview: %s", template.DisplayName),
		Root:      tpm.buildPreviewTree(preview),
		ShowIcons: true,
	}

	err := tpm.ui.ShowTree(ctx, treeConfig)
	if err != nil {
		return fmt.Errorf("failed to display template preview: %w", err)
	}

	// Show summary information
	fmt.Printf("\nTemplate Summary:\n")
	fmt.Printf("  Files: %d\n", preview.Summary.TotalFiles)
	fmt.Printf("  Directories: %d\n", preview.Summary.TotalDirectories)
	fmt.Printf("  Size: %s\n", formatBytes(preview.Summary.TotalSize))
	fmt.Printf("  Templated Files: %d\n", preview.Summary.TemplatedFiles)
	fmt.Printf("  Executable Files: %d\n", preview.Summary.ExecutableFiles)

	return nil
}

// displayCombinedPreview displays a combined template preview
func (tpm *TemplatePreviewManager) displayCombinedPreview(ctx context.Context, combined *CombinedTemplatePreview) error {
	// Show project structure
	treeConfig := interfaces.TreeConfig{
		Title:     "Combined Project Structure",
		Root:      tpm.buildCombinedPreviewTree(combined.ProjectStructure.Root),
		ShowIcons: true,
	}

	err := tpm.ui.ShowTree(ctx, treeConfig)
	if err != nil {
		return fmt.Errorf("failed to display combined preview: %w", err)
	}

	// Show summary
	fmt.Printf("\nProject Summary:\n")
	fmt.Printf("  Templates: %d\n", len(combined.Templates))
	fmt.Printf("  Total Files: %d\n", combined.TotalFiles)
	fmt.Printf("  Estimated Size: %s\n", formatBytes(combined.EstimatedSize))

	// Show selected templates
	fmt.Printf("\nSelected Templates:\n")
	for _, selection := range combined.Templates {
		fmt.Printf("  • %s (%s)\n", selection.Template.DisplayName, selection.Template.Category)
	}

	// Show conflicts if any
	if len(combined.FileConflicts) > 0 {
		fmt.Printf("\nFile Conflicts:\n")
		for _, conflict := range combined.FileConflicts {
			icon := "⚠️"
			switch conflict.Severity {
			case "error":
				icon = "❌"
			case "info":
				icon = "ℹ️"
			}
			fmt.Printf("  %s %s: %s\n", icon, conflict.Path, conflict.Message)
		}
	}

	// Show warnings if any
	if len(combined.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, warning := range combined.Warnings {
			fmt.Printf("  ⚠️ %s\n", warning)
		}
	}

	return nil
}

// buildPreviewTree builds a tree structure for single template preview
func (tpm *TemplatePreviewManager) buildPreviewTree(preview *interfaces.TemplatePreview) interfaces.TreeNode {
	root := interfaces.TreeNode{
		Label:    preview.TemplateName,
		Icon:     "📁",
		Expanded: true,
		Children: []interfaces.TreeNode{},
	}

	// Group files by directory
	dirMap := make(map[string][]interfaces.TemplatePreviewFile)
	for _, file := range preview.Files {
		dir := filepath.Dir(file.Path)
		if dir == "." {
			dir = ""
		}
		dirMap[dir] = append(dirMap[dir], file)
	}

	// Build tree recursively
	root.Children = tpm.buildTreeNodes("", dirMap)

	return root
}

// buildCombinedPreviewTree builds a tree structure for combined preview
func (tpm *TemplatePreviewManager) buildCombinedPreviewTree(root *DirectoryNode) interfaces.TreeNode {
	if root == nil {
		return interfaces.TreeNode{Label: "Empty Project"}
	}

	node := interfaces.TreeNode{
		Label:    root.Name,
		Icon:     "📁",
		Expanded: true,
		Children: []interfaces.TreeNode{},
	}

	// Add child directories
	for _, child := range root.Children {
		childNode := tpm.buildCombinedPreviewTree(child)
		node.Children = append(node.Children, childNode)
	}

	// Add files
	for _, file := range root.Files {
		icon := "📄"
		if file.Executable {
			icon = "⚙️"
		} else if file.Templated {
			icon = "📝"
		}

		fileNode := interfaces.TreeNode{
			Label: file.Name,
			Icon:  icon,
			Metadata: map[string]interface{}{
				"size":   file.Size,
				"source": file.Source,
			},
		}
		node.Children = append(node.Children, fileNode)
	}

	return node
}

// buildTreeNodes recursively builds tree nodes from directory map
func (tpm *TemplatePreviewManager) buildTreeNodes(currentDir string, dirMap map[string][]interfaces.TemplatePreviewFile) []interfaces.TreeNode {
	var nodes []interfaces.TreeNode

	// Add files in current directory
	if files, exists := dirMap[currentDir]; exists {
		for _, file := range files {
			if file.Type == "file" {
				icon := "📄"
				if file.Executable {
					icon = "⚙️"
				} else if file.Templated {
					icon = "📝"
				}

				nodes = append(nodes, interfaces.TreeNode{
					Label: file.Path,
					Icon:  icon,
					Metadata: map[string]interface{}{
						"size": file.Size,
						"type": file.Type,
					},
				})
			}
		}
	}

	return nodes
}

// Utility functions

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// formatBytes formats byte size in human readable format
func formatBytes(bytes int64) string {
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
