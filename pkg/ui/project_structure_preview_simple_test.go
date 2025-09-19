package ui

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

func TestProjectStructurePreviewGenerator_Basic(t *testing.T) {
	// Create a minimal test without mocks to verify the structure
	selections := []TemplateSelection{
		{
			Template: interfaces.TemplateInfo{
				Name:        "nextjs-app",
				DisplayName: "Next.js Application",
				Category:    "frontend",
				Technology:  "React",
			},
			Selected: true,
		},
	}

	// Test the categorizeTemplates function
	generator := &ProjectStructurePreviewGenerator{}
	categories := generator.categorizeTemplates(selections)

	if len(categories) == 0 {
		t.Error("Expected categories to be populated")
	}

	if _, exists := categories["frontend"]; !exists {
		t.Error("Expected frontend category to exist")
	}

	if len(categories["frontend"]) != 1 {
		t.Errorf("Expected 1 frontend template, got %d", len(categories["frontend"]))
	}
}

func TestStandardDirectoryTree_Structure(t *testing.T) {
	// Test the directory tree structure
	tree := &StandardDirectoryTree{}

	// Verify the structure can be created
	// Note: tree is created with &StandardDirectoryTree{} so it cannot be nil
	if tree.Root == nil {
		// Initialize root if needed
		tree.Root = &DirectoryNode{}
	}

	// Test directory node creation
	node := &DirectoryNode{
		Name:        "test",
		Path:        "test",
		Children:    []*DirectoryNode{},
		Files:       []*FileNode{},
		Source:      "test-source",
		Description: "test description",
	}

	if node.Name != "test" {
		t.Errorf("Expected node name 'test', got '%s'", node.Name)
	}

	if node.Path != "test" {
		t.Errorf("Expected node path 'test', got '%s'", node.Path)
	}
}

func TestComponentTemplateMapping_Structure(t *testing.T) {
	// Test the component mapping structure
	mapping := ComponentTemplateMapping{
		Component:     "Frontend Applications",
		Directory:     "App/",
		Templates:     []string{"nextjs-app"},
		Description:   "Frontend components",
		FileCount:     10,
		EstimatedSize: 1024,
	}

	if mapping.Component != "Frontend Applications" {
		t.Errorf("Expected component 'Frontend Applications', got '%s'", mapping.Component)
	}

	if mapping.Directory != "App/" {
		t.Errorf("Expected directory 'App/', got '%s'", mapping.Directory)
	}

	if len(mapping.Templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(mapping.Templates))
	}

	if mapping.Templates[0] != "nextjs-app" {
		t.Errorf("Expected template 'nextjs-app', got '%s'", mapping.Templates[0])
	}
}

func TestProjectSummary_Structure(t *testing.T) {
	// Test the project summary structure
	summary := ProjectSummary{
		TotalDirectories: 5,
		TotalFiles:       20,
		EstimatedSize:    1024 * 1024,
		Technologies:     []string{"React", "Go"},
		Categories:       []string{"frontend", "backend"},
		Dependencies:     []string{"dep1", "dep2"},
	}

	if summary.TotalDirectories != 5 {
		t.Errorf("Expected 5 directories, got %d", summary.TotalDirectories)
	}

	if summary.TotalFiles != 20 {
		t.Errorf("Expected 20 files, got %d", summary.TotalFiles)
	}

	if len(summary.Technologies) != 2 {
		t.Errorf("Expected 2 technologies, got %d", len(summary.Technologies))
	}

	if len(summary.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(summary.Categories))
	}
}
