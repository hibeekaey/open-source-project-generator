// Package tree provides tree structure building functionality for project previews.
package tree

import (
	"path/filepath"
	"sort"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Builder handles building tree structures from directory nodes
type Builder struct{}

// NewBuilder creates a new tree builder
func NewBuilder() *Builder {
	return &Builder{}
}

// BuildDirectoryTree builds a hierarchical directory tree from a flat structure
func (b *Builder) BuildDirectoryTree(directories map[string]*DirectoryNode, files map[string]*FileNode) *DirectoryNode {
	root := &DirectoryNode{
		Name:     "project",
		Path:     "",
		Children: []*DirectoryNode{},
		Files:    []*FileNode{},
	}

	// Build parent-child relationships
	for path, dir := range directories {
		if path == "" || path == "." {
			continue
		}

		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = ""
		}

		if parentPath == "" {
			root.Children = append(root.Children, dir)
		} else if parent, exists := directories[parentPath]; exists {
			parent.Children = append(parent.Children, dir)
		}
	}

	// Add files to their parent directories
	for path, file := range files {
		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = ""
		}

		if parentPath == "" {
			root.Files = append(root.Files, file)
		} else if parent, exists := directories[parentPath]; exists {
			parent.Files = append(parent.Files, file)
		}
	}

	// Sort children and files
	b.SortDirectoryTree(root)

	return root
}

// SortDirectoryTree recursively sorts directory contents
func (b *Builder) SortDirectoryTree(node *DirectoryNode) {
	// Sort directories first, then files
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Name < node.Children[j].Name
	})

	sort.Slice(node.Files, func(i, j int) bool {
		return node.Files[i].Name < node.Files[j].Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		b.SortDirectoryTree(child)
	}
}

// BuildPreviewTree builds a tree structure for single template preview
func (b *Builder) BuildPreviewTree(preview *interfaces.TemplatePreview) interfaces.TreeNode {
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
	root.Children = b.buildTreeNodes("", dirMap)

	return root
}

// buildTreeNodes recursively builds tree nodes from directory map
func (b *Builder) buildTreeNodes(currentDir string, dirMap map[string][]interfaces.TemplatePreviewFile) []interfaces.TreeNode {
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
