// Package components provides component-specific preview functionality.
package components

import (
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// BackendPreview handles backend component preview generation
type BackendPreview struct{}

// NewBackendPreview creates a new backend preview generator
func NewBackendPreview() *BackendPreview {
	return &BackendPreview{}
}

// CreateCommonServerDirectory creates the CommonServer/ directory structure for backend templates
func (bp *BackendPreview) CreateCommonServerDirectory(templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview) *tree.DirectoryNode {
	serverDir := &tree.DirectoryNode{
		Name:        "CommonServer",
		Path:        "CommonServer",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
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
		childDir := &tree.DirectoryNode{
			Name:        subdir.name,
			Path:        filepath.Join("CommonServer", subdir.name),
			Children:    []*tree.DirectoryNode{},
			Files:       []*tree.FileNode{},
			Source:      "backend-template",
			Description: subdir.description,
		}
		bp.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
		serverDir.Children = append(serverDir.Children, childDir)
	}

	return serverDir
}

// populateDirectoryFromTemplates populates a directory with files from matching templates
func (bp *BackendPreview) populateDirectoryFromTemplates(dir *tree.DirectoryNode, templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview, targetSubdir string) {
	// Use the same logic as frontend but for backend-specific templates
	frontend := NewFrontendPreview()
	frontend.populateDirectoryFromTemplates(dir, templates, previews, targetSubdir)
}
