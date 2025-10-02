// Package tree provides tree formatting utilities for project previews.
package tree

import (
	"fmt"
)

// Formatter handles formatting tree-related data for display
type Formatter struct{}

// NewFormatter creates a new tree formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatBytes formats byte size in human readable format
func (f *Formatter) FormatBytes(bytes int64) string {
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

// CountDirectoriesRecursively counts directories in a tree
func (f *Formatter) CountDirectoriesRecursively(dir *DirectoryNode) int {
	if dir == nil {
		return 0
	}

	count := 1 // Count this directory
	for _, child := range dir.Children {
		count += f.CountDirectoriesRecursively(child)
	}
	return count
}

// CountFilesRecursively counts files in a tree
func (f *Formatter) CountFilesRecursively(dir *DirectoryNode) int {
	if dir == nil {
		return 0
	}

	count := len(dir.Files)
	for _, child := range dir.Children {
		count += f.CountFilesRecursively(child)
	}
	return count
}

// CalculateDirectorySize calculates total size of files in a directory tree
func (f *Formatter) CalculateDirectorySize(dir *DirectoryNode) int64 {
	if dir == nil {
		return 0
	}

	var size int64
	for _, file := range dir.Files {
		size += file.Size
	}
	for _, child := range dir.Children {
		size += f.CalculateDirectorySize(child)
	}
	return size
}

// HasChildDirectory checks if a directory already has a child with the given name
func (f *Formatter) HasChildDirectory(parent *DirectoryNode, name string) bool {
	for _, child := range parent.Children {
		if child.Name == name {
			return true
		}
	}
	return false
}
