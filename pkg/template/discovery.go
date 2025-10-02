// Package template provides template discovery and scanning functionality
package template

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	yaml "gopkg.in/yaml.v3"
)

// TemplateDiscovery handles template discovery and scanning operations
type TemplateDiscovery struct {
	embeddedFS    fs.FS
	externalPaths []string
	cache         map[string]*models.TemplateInfo
	cacheTime     time.Time
	cacheTTL      time.Duration
	mutex         sync.RWMutex
}

// NewTemplateDiscovery creates a new template discovery instance
func NewTemplateDiscovery(embeddedFS fs.FS) *TemplateDiscovery {
	return &TemplateDiscovery{
		embeddedFS:    embeddedFS,
		externalPaths: []string{},
		cache:         make(map[string]*models.TemplateInfo),
		cacheTTL:      5 * time.Minute, // Cache templates for 5 minutes
	}
}

// DiscoverTemplates discovers all available templates from embedded filesystem and external paths
func (td *TemplateDiscovery) DiscoverTemplates() ([]*models.TemplateInfo, error) {
	// Check cache first with read lock
	td.mutex.RLock()
	if time.Since(td.cacheTime) < td.cacheTTL && len(td.cache) > 0 {
		templates := make([]*models.TemplateInfo, 0, len(td.cache))
		for _, tmpl := range td.cache {
			templates = append(templates, tmpl)
		}
		td.mutex.RUnlock()
		return templates, nil
	}
	td.mutex.RUnlock()

	// Acquire write lock for cache update
	td.mutex.Lock()
	defer td.mutex.Unlock()

	// Double-check cache after acquiring write lock
	if time.Since(td.cacheTime) < td.cacheTTL && len(td.cache) > 0 {
		templates := make([]*models.TemplateInfo, 0, len(td.cache))
		for _, tmpl := range td.cache {
			templates = append(templates, tmpl)
		}
		return templates, nil
	}

	var templates []*models.TemplateInfo

	// Discover embedded templates
	embeddedTemplates, err := td.DiscoverEmbeddedTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover embedded templates: %w", err)
	}
	templates = append(templates, embeddedTemplates...)

	// TODO: Add external template discovery when needed
	// externalTemplates, err := td.discoverExternalTemplates()
	// templates = append(templates, externalTemplates...)

	// Update cache
	td.cache = make(map[string]*models.TemplateInfo)
	for _, tmpl := range templates {
		td.cache[tmpl.Name] = tmpl
	}
	td.cacheTime = time.Now()

	// Sort templates by name
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// DiscoverEmbeddedTemplates discovers templates from the embedded filesystem
func (td *TemplateDiscovery) DiscoverEmbeddedTemplates() ([]*models.TemplateInfo, error) {
	var templates []*models.TemplateInfo

	// Walk through embedded template directories
	err := fs.WalkDir(td.embeddedFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root templates directory
		if path == "templates" {
			return nil
		}

		// Only process directories that are direct children of category directories
		if d.IsDir() && td.isTemplateDirectory(path) {
			templateInfo, err := td.createTemplateInfoFromPath(path)
			if err != nil {
				// Log error but continue processing other templates
				fmt.Printf("⚠️  Failed to process template at %s: %v\n", path, err)
				return nil
			}
			templates = append(templates, templateInfo)
		}

		return nil
	})

	return templates, err
}

// FilterTemplates applies filtering criteria to templates
func (td *TemplateDiscovery) FilterTemplates(templates []*models.TemplateInfo, filter interfaces.TemplateFilter) []*models.TemplateInfo {
	var filtered []*models.TemplateInfo

	for _, tmpl := range templates {
		if td.matchesFilter(tmpl, filter) {
			filtered = append(filtered, tmpl)
		}
	}

	return filtered
}

// SearchTemplates searches for templates by query string
func (td *TemplateDiscovery) SearchTemplates(templates []*models.TemplateInfo, query string) []*models.TemplateInfo {
	// Return no templates for empty query
	if strings.TrimSpace(query) == "" {
		return []*models.TemplateInfo{}
	}

	query = strings.ToLower(query)
	var matches []*models.TemplateInfo

	for _, tmpl := range templates {
		// Search in name, display name, description, tags, and keywords
		if td.matchesQuery(tmpl, query) {
			matches = append(matches, tmpl)
		}
	}

	return matches
}

// GetTemplateByName finds a template by name from the given list
func (td *TemplateDiscovery) GetTemplateByName(templates []*models.TemplateInfo, name string) *models.TemplateInfo {
	for _, tmpl := range templates {
		if tmpl.Name == name {
			return tmpl
		}
	}
	return nil
}

// GetTemplatesByCategory filters templates by category
func (td *TemplateDiscovery) GetTemplatesByCategory(templates []*models.TemplateInfo, category string) []*models.TemplateInfo {
	var filtered []*models.TemplateInfo
	for _, tmpl := range templates {
		if strings.EqualFold(tmpl.Category, category) {
			filtered = append(filtered, tmpl)
		}
	}
	return filtered
}

// GetTemplatesByTechnology filters templates by technology
func (td *TemplateDiscovery) GetTemplatesByTechnology(templates []*models.TemplateInfo, technology string) []*models.TemplateInfo {
	var filtered []*models.TemplateInfo
	for _, tmpl := range templates {
		if strings.EqualFold(tmpl.Technology, technology) {
			filtered = append(filtered, tmpl)
		}
	}
	return filtered
}

// ClearCache clears the template cache
func (td *TemplateDiscovery) ClearCache() {
	td.mutex.Lock()
	defer td.mutex.Unlock()
	td.cache = make(map[string]*models.TemplateInfo)
	td.cacheTime = time.Time{}
}

// RefreshCache refreshes the template cache
func (td *TemplateDiscovery) RefreshCache() error {
	td.ClearCache()
	_, err := td.DiscoverTemplates()
	return err
}

// AddExternalPath adds an external path for template discovery
func (td *TemplateDiscovery) AddExternalPath(path string) {
	td.externalPaths = append(td.externalPaths, path)
}

// RemoveExternalPath removes an external path from template discovery
func (td *TemplateDiscovery) RemoveExternalPath(path string) {
	for i, p := range td.externalPaths {
		if p == path {
			td.externalPaths = append(td.externalPaths[:i], td.externalPaths[i+1:]...)
			break
		}
	}
}

// GetExternalPaths returns the list of external paths
func (td *TemplateDiscovery) GetExternalPaths() []string {
	return append([]string{}, td.externalPaths...) // Return a copy
}

// isTemplateDirectory checks if a directory path represents a template
func (td *TemplateDiscovery) isTemplateDirectory(path string) bool {
	parts := strings.Split(path, "/")
	// Should be templates/category/template-name
	return len(parts) == 3 && parts[0] == "templates"
}

// createTemplateInfoFromPath creates template info from a filesystem path
func (td *TemplateDiscovery) createTemplateInfoFromPath(path string) (*models.TemplateInfo, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid template path structure: %s", path)
	}

	category := parts[1]
	templateName := parts[2]

	// Create basic template info
	templateInfo := &models.TemplateInfo{
		Name:        templateName,
		DisplayName: td.formatDisplayName(templateName),
		Category:    category,
		Path:        path,
		Source:      "embedded",
		Version:     "1.0.0", // Default version
	}

	// Try to load metadata from template.yaml or template.yml
	metadata, err := td.loadTemplateMetadata(path)
	if err == nil {
		templateInfo.Metadata = *metadata
		templateInfo.DisplayName = metadata.DisplayName
		templateInfo.Description = metadata.Description
		templateInfo.Technology = metadata.Technology
		templateInfo.Tags = metadata.Tags
		templateInfo.Dependencies = metadata.Dependencies
		templateInfo.Version = metadata.Version
	} else {
		// Set defaults based on path analysis
		templateInfo.Description = fmt.Sprintf("%s template for %s projects",
			td.formatDisplayName(templateName), category)
		templateInfo.Technology = td.inferTechnology(templateName)
		templateInfo.Tags = td.inferTags(templateName, category)
	}

	// Calculate template size and file count
	size, fileCount, err := td.calculateTemplateStats(path)
	if err == nil {
		templateInfo.Size = size
		templateInfo.FileCount = fileCount
	}

	templateInfo.LastModified = time.Now() // For embedded templates, use current time

	return templateInfo, nil
}

// matchesFilter checks if a template matches the given filter
func (td *TemplateDiscovery) matchesFilter(tmpl *models.TemplateInfo, filter interfaces.TemplateFilter) bool {
	// Category filter
	if filter.Category != "" && !strings.EqualFold(tmpl.Category, filter.Category) {
		return false
	}

	// Technology filter
	if filter.Technology != "" && !strings.EqualFold(tmpl.Technology, filter.Technology) {
		return false
	}

	// Tags filter (template must have all specified tags)
	if len(filter.Tags) > 0 {
		templateTags := make(map[string]bool)
		for _, tag := range tmpl.Tags {
			templateTags[strings.ToLower(tag)] = true
		}

		for _, filterTag := range filter.Tags {
			if !templateTags[strings.ToLower(filterTag)] {
				return false
			}
		}
	}

	// Version filters (simplified - would need proper semver comparison)
	if filter.MinVersion != "" && tmpl.Version < filter.MinVersion {
		return false
	}
	if filter.MaxVersion != "" && tmpl.Version > filter.MaxVersion {
		return false
	}

	return true
}

// matchesQuery checks if a template matches a search query
func (td *TemplateDiscovery) matchesQuery(tmpl *models.TemplateInfo, query string) bool {
	// Search in name
	if strings.Contains(strings.ToLower(tmpl.Name), query) {
		return true
	}

	// Search in display name
	if strings.Contains(strings.ToLower(tmpl.DisplayName), query) {
		return true
	}

	// Search in description
	if strings.Contains(strings.ToLower(tmpl.Description), query) {
		return true
	}

	// Search in technology
	if strings.Contains(strings.ToLower(tmpl.Technology), query) {
		return true
	}

	// Search in tags
	for _, tag := range tmpl.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	// Search in keywords
	for _, keyword := range tmpl.Metadata.Keywords {
		if strings.Contains(strings.ToLower(keyword), query) {
			return true
		}
	}

	return false
}

// loadTemplateMetadata loads metadata from template.yaml or template.yml
func (td *TemplateDiscovery) loadTemplateMetadata(templatePath string) (*models.TemplateMetadata, error) {
	// Try template.yaml first, then template.yml
	metadataFiles := []string{"template.yaml", "template.yml"}

	for _, filename := range metadataFiles {
		metadataPath := filepath.Join(templatePath, filename)
		if content, err := fs.ReadFile(td.embeddedFS, metadataPath); err == nil {
			return td.parseTemplateYAML(content, filepath.Base(templatePath))
		}
	}

	return nil, fmt.Errorf("no metadata file found")
}

// parseTemplateYAML parses template.yaml content into models.TemplateMetadata
func (td *TemplateDiscovery) parseTemplateYAML(content []byte, templateName string) (*models.TemplateMetadata, error) {
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
		metadata.DisplayName = td.formatDisplayName(templateName)
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

// calculateTemplateStats calculates size and file count for a template
func (td *TemplateDiscovery) calculateTemplateStats(templatePath string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := fs.WalkDir(td.embeddedFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			fileCount++
			if info, err := d.Info(); err == nil {
				totalSize += info.Size()
			}
		}

		return nil
	})

	return totalSize, fileCount, err
}

// formatDisplayName formats a template name for display
func (td *TemplateDiscovery) formatDisplayName(name string) string {
	// Convert kebab-case to Title Case
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// inferTechnology infers technology from template name
func (td *TemplateDiscovery) inferTechnology(templateName string) string {
	name := strings.ToLower(templateName)

	if strings.Contains(name, "go") || strings.Contains(name, "gin") {
		return "Go"
	}
	if strings.Contains(name, "nextjs") || strings.Contains(name, "next") {
		return "Next.js"
	}
	if strings.Contains(name, "react") {
		return "React"
	}
	if strings.Contains(name, "node") {
		return "Node.js"
	}
	if strings.Contains(name, "python") {
		return "Python"
	}
	if strings.Contains(name, "java") {
		return "Java"
	}
	if strings.Contains(name, "kotlin") {
		return "Kotlin"
	}
	if strings.Contains(name, "swift") {
		return "Swift"
	}
	if strings.Contains(name, "docker") {
		return "Docker"
	}
	if strings.Contains(name, "kubernetes") || strings.Contains(name, "k8s") {
		return "Kubernetes"
	}
	if strings.Contains(name, "terraform") {
		return "Terraform"
	}

	return "Unknown"
}

// inferTags infers tags from template name and category
func (td *TemplateDiscovery) inferTags(templateName, category string) []string {
	var tags []string
	name := strings.ToLower(templateName)

	// Add category as a tag
	tags = append(tags, category)

	// Add technology-specific tags
	if strings.Contains(name, "go") || strings.Contains(name, "gin") {
		tags = append(tags, "go", "backend", "api")
	}
	if strings.Contains(name, "nextjs") || strings.Contains(name, "next") {
		tags = append(tags, "nextjs", "react", "frontend", "web")
	}
	if strings.Contains(name, "react") {
		tags = append(tags, "react", "frontend", "web")
	}
	if strings.Contains(name, "admin") {
		tags = append(tags, "admin", "dashboard")
	}
	if strings.Contains(name, "home") {
		tags = append(tags, "landing", "homepage")
	}
	if strings.Contains(name, "mobile") {
		tags = append(tags, "mobile", "app")
	}
	if strings.Contains(name, "android") {
		tags = append(tags, "android", "mobile")
	}
	if strings.Contains(name, "ios") {
		tags = append(tags, "ios", "mobile")
	}
	if strings.Contains(name, "kotlin") {
		tags = append(tags, "kotlin", "android")
	}
	if strings.Contains(name, "swift") {
		tags = append(tags, "swift", "ios")
	}
	if strings.Contains(name, "docker") {
		tags = append(tags, "docker", "containerization")
	}
	if strings.Contains(name, "kubernetes") || strings.Contains(name, "k8s") {
		tags = append(tags, "kubernetes", "orchestration")
	}
	if strings.Contains(name, "terraform") {
		tags = append(tags, "terraform", "iac")
	}

	return tags
}
