// Package metadata provides template metadata parsing and validation functionality.
package metadata

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

// MetadataParser handles parsing of template metadata from various sources.
type MetadataParser struct {
	embeddedFS fs.FS
}

// NewMetadataParser creates a new metadata parser instance.
func NewMetadataParser(embeddedFS fs.FS) *MetadataParser {
	return &MetadataParser{
		embeddedFS: embeddedFS,
	}
}

// LoadTemplateMetadata loads metadata from template.yaml or template.yml files.
// It first tries template.yaml, then template.yml as fallback.
func (p *MetadataParser) LoadTemplateMetadata(templatePath string) (*models.TemplateMetadata, error) {
	// Try template.yaml first, then template.yml
	metadataFiles := []string{"template.yaml", "template.yml"}

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)

		// Try embedded filesystem first if available
		if p.embeddedFS != nil {
			if content, err := fs.ReadFile(p.embeddedFS, metadataPath); err == nil {
				return p.ParseTemplateYAML(content, filepath.Base(templatePath))
			}
		}

		// Try regular filesystem
		// #nosec G304 - metadataPath is constructed from validated templatePath and fixed filename
		if content, err := os.ReadFile(metadataPath); err == nil {
			return p.ParseTemplateYAML(content, filepath.Base(templatePath))
		}
	}

	return nil, fmt.Errorf("no metadata file found")
}

// ParseTemplateYAML parses template.yaml content into models.TemplateMetadata.
func (p *MetadataParser) ParseTemplateYAML(content []byte, templateName string) (*models.TemplateMetadata, error) {
	// Define a structure that matches the template.yaml format
	type TemplateYAML struct {
		Name         string   `yaml:"name"`
		DisplayName  string   `yaml:"display_name"`
		Description  string   `yaml:"description"`
		Category     string   `yaml:"category"`
		Technology   string   `yaml:"technology"`
		Version      string   `yaml:"version"`
		Tags         []string `yaml:"tags"`
		Dependencies []string `yaml:"dependencies"`
		Metadata     struct {
			Author      string            `yaml:"author"`
			License     string            `yaml:"license"`
			Repository  string            `yaml:"repository"`
			Homepage    string            `yaml:"homepage"`
			Keywords    []string          `yaml:"keywords"`
			Maintainers []string          `yaml:"maintainers"`
			Created     time.Time         `yaml:"created"`
			Updated     time.Time         `yaml:"updated"`
			Variables   map[string]string `yaml:"variables"`
			MinVersion  string            `yaml:"min_version"`
			MaxVersion  string            `yaml:"max_version"`
		} `yaml:"metadata"`
	}

	var yamlData TemplateYAML
	if err := yaml.Unmarshal(content, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse template YAML: %w", err)
	}

	// Convert to models.TemplateMetadata
	metadata := &models.TemplateMetadata{
		Name:         yamlData.Name,
		DisplayName:  yamlData.DisplayName,
		Description:  yamlData.Description,
		Version:      yamlData.Version,
		Author:       yamlData.Metadata.Author,
		License:      yamlData.Metadata.License,
		Category:     yamlData.Category,
		Technology:   yamlData.Technology,
		Tags:         yamlData.Tags,
		Dependencies: yamlData.Dependencies,
		CreatedAt:    yamlData.Metadata.Created,
		UpdatedAt:    yamlData.Metadata.Updated,
		Homepage:     yamlData.Metadata.Homepage,
		Repository:   yamlData.Metadata.Repository,
		Keywords:     yamlData.Metadata.Keywords,
		MinVersion:   yamlData.Metadata.MinVersion,
		MaxVersion:   yamlData.Metadata.MaxVersion,
		Variables:    make(map[string]models.TemplateVar),
	}

	// Convert variables from simple string map to TemplateVar map
	for name, description := range yamlData.Metadata.Variables {
		metadata.Variables[name] = models.TemplateVar{
			Name:        name,
			Type:        "string", // Default type
			Description: description,
			Required:    false, // Default to not required
		}
	}

	// Set defaults if not provided
	if metadata.Name == "" {
		metadata.Name = templateName
	}
	if metadata.DisplayName == "" {
		metadata.DisplayName = p.formatDisplayName(templateName)
	}
	if metadata.Version == "" {
		metadata.Version = "1.0.0"
	}
	if metadata.License == "" {
		metadata.License = "MIT"
	}
	if metadata.Author == "" {
		metadata.Author = "Open Source Project Generator"
	}

	// Initialize slices if nil to prevent issues
	if metadata.Tags == nil {
		metadata.Tags = []string{}
	}
	if metadata.Dependencies == nil {
		metadata.Dependencies = []string{}
	}
	if metadata.Keywords == nil {
		metadata.Keywords = []string{}
	}

	return metadata, nil
}

// formatDisplayName formats a template name for display.
// Converts kebab-case to Title Case.
func (p *MetadataParser) formatDisplayName(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// LoadMetadataFromFile loads metadata from a specific file path.
func (p *MetadataParser) LoadMetadataFromFile(filePath string) (*models.TemplateMetadata, error) {
	// Validate file path for security
	if err := utils.ValidatePath(filePath); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// #nosec G304 - filePath is validated above using ValidatePath
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file %s: %w", filePath, err)
	}

	templateName := filepath.Base(filepath.Dir(filePath))
	return p.ParseTemplateYAML(content, templateName)
}

// LoadMetadataFromEmbedded loads metadata from embedded filesystem.
func (p *MetadataParser) LoadMetadataFromEmbedded(templatePath string) (*models.TemplateMetadata, error) {
	if p.embeddedFS == nil {
		return nil, fmt.Errorf("embedded filesystem not available")
	}

	metadataFiles := []string{"template.yaml", "template.yml"}

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if content, err := fs.ReadFile(p.embeddedFS, metadataPath); err == nil {
			return p.ParseTemplateYAML(content, filepath.Base(templatePath))
		}
	}

	return nil, fmt.Errorf("no metadata file found in embedded template %s", templatePath)
}

// ParseMetadataContent parses raw metadata content with explicit template name.
func (p *MetadataParser) ParseMetadataContent(content []byte, templateName string) (*models.TemplateMetadata, error) {
	return p.ParseTemplateYAML(content, templateName)
}
