// Package display provides display functionality for project previews.
package display

import (
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// Console handles console-based display of project previews
type Console struct {
	formatter *tree.Formatter
}

// NewConsole creates a new console display handler
func NewConsole() *Console {
	return &Console{
		formatter: tree.NewFormatter(),
	}
}

// DisplayProjectHeader shows the project header information
func (c *Console) DisplayProjectHeader(preview *ProjectStructurePreview) {
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

// DisplayComponentMappings shows which templates contribute to each component
func (c *Console) DisplayComponentMappings(preview *ProjectStructurePreview) {
	fmt.Printf("\n")
	fmt.Printf("Component-Template Mapping:\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n")

	for _, mapping := range preview.ComponentMapping {
		fmt.Printf("\n📦 %s (%s)\n", mapping.Component, mapping.Directory)
		fmt.Printf("   Description: %s\n", mapping.Description)
		fmt.Printf("   Files: %d | Size: %s\n", mapping.FileCount, c.formatter.FormatBytes(mapping.EstimatedSize))
		fmt.Printf("   Templates: %s\n", strings.Join(mapping.Templates, ", "))
	}
}

// DisplayProjectSummary shows overall project statistics
func (c *Console) DisplayProjectSummary(preview *ProjectStructurePreview) {
	fmt.Printf("\n")
	fmt.Printf("Project Summary:\n")
	fmt.Printf("───────────────────────────────────────────────────────────────\n")
	fmt.Printf("Directories:   %d\n", preview.Summary.TotalDirectories)
	fmt.Printf("Files:         %d\n", preview.Summary.TotalFiles)
	fmt.Printf("Estimated Size: %s\n", c.formatter.FormatBytes(preview.Summary.EstimatedSize))

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

// DisplayWarningsAndConflicts shows warnings and conflicts
func (c *Console) DisplayWarningsAndConflicts(preview *ProjectStructurePreview) {
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

// DisplayTemplateSummary shows summary information for a single template
func (c *Console) DisplayTemplateSummary(preview interface{}) {
	// Type assertion to handle different preview types
	switch p := preview.(type) {
	case *TemplatePreview:
		fmt.Printf("\nTemplate Summary:\n")
		fmt.Printf("  Files: %d\n", p.Summary.TotalFiles)
		fmt.Printf("  Directories: %d\n", p.Summary.TotalDirectories)
		fmt.Printf("  Size: %s\n", c.formatter.FormatBytes(p.Summary.TotalSize))
		fmt.Printf("  Templated Files: %d\n", p.Summary.TemplatedFiles)
		fmt.Printf("  Executable Files: %d\n", p.Summary.ExecutableFiles)
	case *CombinedTemplatePreview:
		fmt.Printf("\nProject Summary:\n")
		fmt.Printf("  Templates: %d\n", len(p.Templates))
		fmt.Printf("  Total Files: %d\n", p.TotalFiles)
		fmt.Printf("  Estimated Size: %s\n", c.formatter.FormatBytes(p.EstimatedSize))

		// Show selected templates
		fmt.Printf("\nSelected Templates:\n")
		fmt.Printf("  Templates: %d\n", len(p.Templates))

		// Show conflicts if any
		if len(p.FileConflicts) > 0 {
			fmt.Printf("\nFile Conflicts:\n")
			for _, conflict := range p.FileConflicts {
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
		if len(p.Warnings) > 0 {
			fmt.Printf("\nWarnings:\n")
			for _, warning := range p.Warnings {
				fmt.Printf("  ⚠️ %s\n", warning)
			}
		}
	}
}

// ProjectStructurePreview represents a complete project structure preview
type ProjectStructurePreview struct {
	ProjectName       string                     `json:"project_name"`
	OutputDirectory   string                     `json:"output_directory"`
	SelectedTemplates []interface{}              `json:"selected_templates"`
	DirectoryTree     *StandardDirectoryTree     `json:"directory_tree"`
	ComponentMapping  []ComponentTemplateMapping `json:"component_mapping"`
	Summary           ProjectSummary             `json:"summary"`
	Warnings          []string                   `json:"warnings"`
	Conflicts         []FileConflict             `json:"conflicts"`
}

// StandardDirectoryTree represents the standardized project directory structure
type StandardDirectoryTree struct {
	Root         *tree.DirectoryNode `json:"root"`
	App          *tree.DirectoryNode `json:"app,omitempty"`
	CommonServer *tree.DirectoryNode `json:"common_server,omitempty"`
	Mobile       *tree.DirectoryNode `json:"mobile,omitempty"`
	Deploy       *tree.DirectoryNode `json:"deploy,omitempty"`
	Docs         *tree.DirectoryNode `json:"docs"`
	Scripts      *tree.DirectoryNode `json:"scripts"`
	GitHub       *tree.DirectoryNode `json:"github"`
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

// FileConflict represents a conflict between templates
type FileConflict struct {
	Path       string   `json:"path"`
	Templates  []string `json:"templates"`
	Severity   string   `json:"severity"` // "error", "warning", "info"
	Message    string   `json:"message"`
	Resolvable bool     `json:"resolvable"`
}

// TemplatePreview represents a preview of a single template
type TemplatePreview struct {
	Summary struct {
		TotalFiles       int   `json:"total_files"`
		TotalDirectories int   `json:"total_directories"`
		TotalSize        int64 `json:"total_size"`
		TemplatedFiles   int   `json:"templated_files"`
		ExecutableFiles  int   `json:"executable_files"`
	} `json:"summary"`
}
