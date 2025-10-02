// Package display provides interactive display functionality for project previews.
package display

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// Interactive handles interactive display of project previews
type Interactive struct {
	ui       interfaces.InteractiveUIInterface
	renderer *tree.Renderer
	console  *Console
}

// NewInteractive creates a new interactive display handler
func NewInteractive(ui interfaces.InteractiveUIInterface) *Interactive {
	return &Interactive{
		ui:       ui,
		renderer: tree.NewRenderer(ui),
		console:  NewConsole(),
	}
}

// DisplayProjectStructurePreview shows the complete project structure preview
func (i *Interactive) DisplayProjectStructurePreview(ctx context.Context, preview *ProjectStructurePreview) error {
	if preview == nil {
		return fmt.Errorf("preview cannot be nil")
	}

	// Display project header
	i.console.DisplayProjectHeader(preview)

	// Display directory structure tree
	err := i.DisplayDirectoryTree(ctx, preview)
	if err != nil {
		return fmt.Errorf("failed to display directory tree: %w", err)
	}

	// Display component mappings
	i.console.DisplayComponentMappings(preview)

	// Display project summary
	i.console.DisplayProjectSummary(preview)

	// Display warnings and conflicts
	i.console.DisplayWarningsAndConflicts(preview)

	return nil
}

// DisplayDirectoryTree shows the directory structure as a tree
func (i *Interactive) DisplayDirectoryTree(ctx context.Context, preview *ProjectStructurePreview) error {
	if preview.DirectoryTree.Root == nil {
		return fmt.Errorf("no directory tree to display")
	}

	return i.renderer.RenderDirectoryTree(ctx, preview.DirectoryTree.Root, "Project Directory Structure")
}

// DisplayTemplatePreview displays a single template preview
func (i *Interactive) DisplayTemplatePreview(ctx context.Context, template interfaces.TemplateInfo, preview *interfaces.TemplatePreview) error {
	err := i.renderer.RenderTemplatePreview(ctx, template, preview)
	if err != nil {
		return fmt.Errorf("failed to display template preview: %w", err)
	}

	// Show summary information
	i.console.DisplayTemplateSummary(&TemplatePreview{
		Summary: struct {
			TotalFiles       int   `json:"total_files"`
			TotalDirectories int   `json:"total_directories"`
			TotalSize        int64 `json:"total_size"`
			TemplatedFiles   int   `json:"templated_files"`
			ExecutableFiles  int   `json:"executable_files"`
		}{
			TotalFiles:       preview.Summary.TotalFiles,
			TotalDirectories: preview.Summary.TotalDirectories,
			TotalSize:        preview.Summary.TotalSize,
			TemplatedFiles:   preview.Summary.TemplatedFiles,
			ExecutableFiles:  preview.Summary.ExecutableFiles,
		},
	})

	return nil
}

// DisplayCombinedPreview displays a combined template preview
func (i *Interactive) DisplayCombinedPreview(ctx context.Context, combined *CombinedTemplatePreview) error {
	// Show summary
	i.console.DisplayTemplateSummary(combined)

	return nil
}

// CombinedTemplatePreview represents a preview of multiple templates combined
type CombinedTemplatePreview struct {
	Templates     []interface{}  `json:"templates"`
	FileConflicts []FileConflict `json:"file_conflicts"`
	Dependencies  []string       `json:"dependencies"`
	EstimatedSize int64          `json:"estimated_size"`
	TotalFiles    int            `json:"total_files"`
	Warnings      []string       `json:"warnings"`
}
