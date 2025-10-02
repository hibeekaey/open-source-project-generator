// Package tree provides tree rendering functionality for project previews.
package tree

import (
	"context"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Renderer handles rendering tree structures for display
type Renderer struct {
	ui interfaces.InteractiveUIInterface
}

// NewRenderer creates a new tree renderer
func NewRenderer(ui interfaces.InteractiveUIInterface) *Renderer {
	return &Renderer{
		ui: ui,
	}
}

// RenderDirectoryTree displays a directory tree structure
func (r *Renderer) RenderDirectoryTree(ctx context.Context, root *DirectoryNode, title string) error {
	if root == nil {
		return fmt.Errorf("no directory tree to display")
	}

	treeConfig := interfaces.TreeConfig{
		Title:      title,
		Root:       r.convertToTreeNode(root),
		Expandable: true,
		ShowIcons:  true,
		MaxDepth:   10,
	}

	return r.ui.ShowTree(ctx, treeConfig)
}

// RenderTemplatePreview displays a single template preview tree
func (r *Renderer) RenderTemplatePreview(ctx context.Context, template interfaces.TemplateInfo, preview *interfaces.TemplatePreview) error {
	builder := NewBuilder()

	// Build tree structure for display
	treeConfig := interfaces.TreeConfig{
		Title:     fmt.Sprintf("Template Preview: %s", template.DisplayName),
		Root:      builder.BuildPreviewTree(preview),
		ShowIcons: true,
	}

	return r.ui.ShowTree(ctx, treeConfig)
}

// RenderCombinedPreview displays a combined template preview tree
func (r *Renderer) RenderCombinedPreview(ctx context.Context, root *DirectoryNode, title string) error {
	treeConfig := interfaces.TreeConfig{
		Title:     title,
		Root:      r.buildCombinedPreviewTree(root),
		ShowIcons: true,
	}

	return r.ui.ShowTree(ctx, treeConfig)
}

// convertToTreeNode converts DirectoryNode to TreeNode for display
func (r *Renderer) convertToTreeNode(dir *DirectoryNode) interfaces.TreeNode {
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
		childNode := r.convertToTreeNode(child)
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

// buildCombinedPreviewTree builds a tree structure for combined preview
func (r *Renderer) buildCombinedPreviewTree(root *DirectoryNode) interfaces.TreeNode {
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
		childNode := r.buildCombinedPreviewTree(child)
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
