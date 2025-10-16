// Package components provides component-specific preview functionality.
package components

import (
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// FrontendPreview handles frontend component preview generation
type FrontendPreview struct{}

// NewFrontendPreview creates a new frontend preview generator
func NewFrontendPreview() *FrontendPreview {
	return &FrontendPreview{}
}

// CreateAppDirectory creates the App/ directory structure for frontend templates
func (fp *FrontendPreview) CreateAppDirectory(templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview) *tree.DirectoryNode {
	appDir := &tree.DirectoryNode{
		Name:        "App",
		Path:        "App",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "frontend-templates",
		Description: "Frontend applications and components",
	}

	// Create standard frontend subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]preview.TemplateSelection) bool
	}{
		{"main", "Main application frontend", func(t []preview.TemplateSelection) bool { return fp.hasTemplateType(t, "nextjs-app") }},
		{"home", "Home/landing page frontend", func(t []preview.TemplateSelection) bool { return fp.hasTemplateType(t, "nextjs-home") }},
		{"admin", "Admin dashboard frontend", func(t []preview.TemplateSelection) bool { return fp.hasTemplateType(t, "nextjs-admin") }},
		{"shared-components", "Shared UI components", func(t []preview.TemplateSelection) bool { return fp.hasTemplateType(t, "shared-components") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &tree.DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("App", subdir.name),
				Children:    []*tree.DirectoryNode{},
				Files:       []*tree.FileNode{},
				Source:      "frontend-template",
				Description: subdir.description,
			}
			fp.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			appDir.Children = append(appDir.Children, childDir)
		}
	}

	return appDir
}

// hasTemplateType checks if any template matches the given type
func (fp *FrontendPreview) hasTemplateType(templates []preview.TemplateSelection, templateType string) bool {
	for _, template := range templates {
		if strings.Contains(strings.ToLower(template.Template.Name), templateType) {
			return true
		}
	}
	return false
}

// populateDirectoryFromTemplates populates a directory with files from matching templates
func (fp *FrontendPreview) populateDirectoryFromTemplates(dir *tree.DirectoryNode, templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview, targetSubdir string) {
	formatter := tree.NewFormatter()

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
						if !formatter.HasChildDirectory(dir, filepath.Base(relativePath)) {
							childDir := &tree.DirectoryNode{
								Name:        filepath.Base(relativePath),
								Path:        filepath.Join(dir.Path, filepath.Base(relativePath)),
								Children:    []*tree.DirectoryNode{},
								Files:       []*tree.FileNode{},
								Source:      template.Template.Name,
								Description: "Generated from " + template.Template.DisplayName + " template",
							}
							dir.Children = append(dir.Children, childDir)
						}
					} else {
						// Add file
						fileNode := &tree.FileNode{
							Name:        filepath.Base(relativePath),
							Path:        filepath.Join(dir.Path, relativePath),
							Size:        file.Size,
							Source:      template.Template.Name,
							Templated:   file.Templated,
							Executable:  file.Executable,
							Description: "Generated from " + template.Template.DisplayName + " template",
						}
						dir.Files = append(dir.Files, fileNode)
					}
				}
			}
		}
	}
}
