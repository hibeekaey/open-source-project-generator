// Package components provides component-specific preview functionality.
package components

import (
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview"
	"github.com/cuesoftinc/open-source-project-generator/pkg/ui/preview/tree"
)

// InfrastructurePreview handles infrastructure component preview generation
type InfrastructurePreview struct{}

// NewInfrastructurePreview creates a new infrastructure preview generator
func NewInfrastructurePreview() *InfrastructurePreview {
	return &InfrastructurePreview{}
}

// CreateDeployDirectory creates the Deploy/ directory structure for infrastructure templates
func (ip *InfrastructurePreview) CreateDeployDirectory(templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview) *tree.DirectoryNode {
	deployDir := &tree.DirectoryNode{
		Name:        "Deploy",
		Path:        "Deploy",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "infrastructure-templates",
		Description: "Deployment and infrastructure configuration",
	}

	// Create infrastructure subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]preview.TemplateSelection) bool
	}{
		{"docker", "Docker containers and compose files", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "docker") }},
		{"k8s", "Kubernetes manifests and configurations", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "kubernetes") }},
		{"terraform", "Infrastructure as Code with Terraform", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "terraform") }},
		{"monitoring", "Monitoring and observability tools", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "monitoring") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &tree.DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("Deploy", subdir.name),
				Children:    []*tree.DirectoryNode{},
				Files:       []*tree.FileNode{},
				Source:      "infrastructure-template",
				Description: subdir.description,
			}
			ip.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			deployDir.Children = append(deployDir.Children, childDir)
		}
	}

	return deployDir
}

// CreateMobileDirectory creates the Mobile/ directory structure for mobile templates
func (ip *InfrastructurePreview) CreateMobileDirectory(templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview) *tree.DirectoryNode {
	mobileDir := &tree.DirectoryNode{
		Name:        "Mobile",
		Path:        "Mobile",
		Children:    []*tree.DirectoryNode{},
		Files:       []*tree.FileNode{},
		Source:      "mobile-templates",
		Description: "Mobile applications for iOS and Android",
	}

	// Create platform-specific subdirectories
	subdirs := []struct {
		name        string
		description string
		condition   func([]preview.TemplateSelection) bool
	}{
		{"android", "Android application (Kotlin)", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "android-kotlin") }},
		{"ios", "iOS application (Swift)", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "ios-swift") }},
		{"shared", "Shared mobile resources and APIs", func(t []preview.TemplateSelection) bool { return ip.hasTemplateType(t, "shared") }},
	}

	for _, subdir := range subdirs {
		if subdir.condition(templates) {
			childDir := &tree.DirectoryNode{
				Name:        subdir.name,
				Path:        filepath.Join("Mobile", subdir.name),
				Children:    []*tree.DirectoryNode{},
				Files:       []*tree.FileNode{},
				Source:      "mobile-template",
				Description: subdir.description,
			}
			ip.populateDirectoryFromTemplates(childDir, templates, previews, subdir.name)
			mobileDir.Children = append(mobileDir.Children, childDir)
		}
	}

	return mobileDir
}

// hasTemplateType checks if any template matches the given type
func (ip *InfrastructurePreview) hasTemplateType(templates []preview.TemplateSelection, templateType string) bool {
	for _, template := range templates {
		if strings.Contains(strings.ToLower(template.Template.Name), templateType) {
			return true
		}
	}
	return false
}

// populateDirectoryFromTemplates populates a directory with files from matching templates
func (ip *InfrastructurePreview) populateDirectoryFromTemplates(dir *tree.DirectoryNode, templates []preview.TemplateSelection, previews map[string]*interfaces.TemplatePreview, targetSubdir string) {
	// Use the same logic as frontend but for infrastructure-specific templates
	frontend := NewFrontendPreview()
	frontend.populateDirectoryFromTemplates(dir, templates, previews, targetSubdir)
}
