// Package ui provides metadata processing for template incorporation.
//
// This file implements the MetadataProcessor which handles incorporation
// of project metadata into template processing and file generation.
package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MetadataProcessor handles incorporation of metadata into templates
type MetadataProcessor struct {
	logger interfaces.Logger
}

// TemplateMetadata represents metadata available to templates
type TemplateMetadata struct {
	// Project information
	ProjectName        string `json:"project_name"`
	ProjectDescription string `json:"project_description"`
	ProjectSlug        string `json:"project_slug"`  // URL-friendly version
	ProjectTitle       string `json:"project_title"` // Display-friendly version

	// Author information
	Author       string `json:"author"`
	Email        string `json:"email"`
	Organization string `json:"organization"`

	// Legal information
	License       string `json:"license"`
	Copyright     string `json:"copyright"`
	CopyrightYear string `json:"copyright_year"`

	// Repository information
	Repository    string `json:"repository"`
	RepositoryURL string `json:"repository_url"`
	ModuleName    string `json:"module_name"`  // For Go modules
	PackageName   string `json:"package_name"` // For npm packages

	// Generation metadata
	GeneratedAt      string `json:"generated_at"`
	GeneratorVersion string `json:"generator_version"`

	// Computed values
	AuthorWithEmail string `json:"author_with_email"`
	LicenseFile     string `json:"license_file"`
	HasRepository   bool   `json:"has_repository"`
	HasOrganization bool   `json:"has_organization"`
}

// FileMetadata represents metadata for specific file types
type FileMetadata struct {
	// Common metadata
	*TemplateMetadata

	// File-specific metadata
	FileName     string `json:"file_name"`
	FileType     string `json:"file_type"`
	RelativePath string `json:"relative_path"`

	// Language-specific metadata
	GoPackage    string `json:"go_package,omitempty"`
	JSModule     string `json:"js_module,omitempty"`
	PythonModule string `json:"python_module,omitempty"`
}

// NewMetadataProcessor creates a new metadata processor
func NewMetadataProcessor(logger interfaces.Logger) *MetadataProcessor {
	return &MetadataProcessor{
		logger: logger,
	}
}

// ProcessMetadata converts project configuration to template metadata
func (mp *MetadataProcessor) ProcessMetadata(config *models.ProjectConfig) (*TemplateMetadata, error) {
	if config == nil {
		return nil, fmt.Errorf("project configuration cannot be nil")
	}

	metadata := &TemplateMetadata{
		ProjectName:        config.Name,
		ProjectDescription: config.Description,
		Author:             config.Author,
		Email:              config.Email,
		Organization:       config.Organization,
		License:            config.License,
		Repository:         config.Repository,
		GeneratorVersion:   config.GeneratorVersion,
	}

	// Process computed values
	mp.processComputedValues(metadata)

	// Process timestamps
	mp.processTimestamps(metadata, config)

	// Process repository information
	mp.processRepositoryInfo(metadata)

	// Process legal information
	mp.processLegalInfo(metadata)

	// Process naming conventions
	mp.processNamingConventions(metadata)

	return metadata, nil
}

// processComputedValues computes derived values from basic metadata
func (mp *MetadataProcessor) processComputedValues(metadata *TemplateMetadata) {
	// Project slug (URL-friendly name)
	metadata.ProjectSlug = mp.createSlug(metadata.ProjectName)

	// Project title (display-friendly name)
	metadata.ProjectTitle = mp.createTitle(metadata.ProjectName)

	// Author with email
	if metadata.Author != "" && metadata.Email != "" {
		metadata.AuthorWithEmail = fmt.Sprintf("%s <%s>", metadata.Author, metadata.Email)
	} else {
		metadata.AuthorWithEmail = metadata.Author
	}

	// Boolean flags
	metadata.HasRepository = metadata.Repository != ""
	metadata.HasOrganization = metadata.Organization != ""

	// License file name
	if metadata.License != "" {
		metadata.LicenseFile = "LICENSE"
	}
}

// processTimestamps handles timestamp-related metadata
func (mp *MetadataProcessor) processTimestamps(metadata *TemplateMetadata, config *models.ProjectConfig) {
	now := time.Now()

	// Use generation time from config if available, otherwise use current time
	if !config.GeneratedAt.IsZero() {
		metadata.GeneratedAt = config.GeneratedAt.Format(time.RFC3339)
		metadata.CopyrightYear = config.GeneratedAt.Format("2006")
	} else {
		metadata.GeneratedAt = now.Format(time.RFC3339)
		metadata.CopyrightYear = now.Format("2006")
	}
}

// processRepositoryInfo processes repository-related metadata
func (mp *MetadataProcessor) processRepositoryInfo(metadata *TemplateMetadata) {
	if metadata.Repository == "" {
		return
	}

	// Ensure repository URL is properly formatted
	metadata.RepositoryURL = metadata.Repository
	if !strings.HasPrefix(metadata.Repository, "http") {
		metadata.RepositoryURL = "https://" + metadata.Repository
	}

	// Extract module name for Go projects
	metadata.ModuleName = mp.extractModuleName(metadata.Repository, metadata.ProjectSlug)

	// Extract package name for npm projects
	metadata.PackageName = metadata.ProjectSlug
}

// processLegalInfo processes legal and copyright information
func (mp *MetadataProcessor) processLegalInfo(metadata *TemplateMetadata) {
	// Generate copyright notice
	copyrightHolder := metadata.Organization
	if copyrightHolder == "" {
		copyrightHolder = metadata.Author
	}

	if copyrightHolder != "" {
		metadata.Copyright = fmt.Sprintf("Copyright (c) %s %s", metadata.CopyrightYear, copyrightHolder)
	}
}

// processNamingConventions processes various naming conventions
func (mp *MetadataProcessor) processNamingConventions(metadata *TemplateMetadata) {
	// Additional naming conventions can be added here
	// For example: PascalCase, camelCase, UPPER_CASE versions
}

// ProcessFileMetadata creates file-specific metadata
func (mp *MetadataProcessor) ProcessFileMetadata(templateMetadata *TemplateMetadata, filePath string) (*FileMetadata, error) {
	if templateMetadata == nil {
		return nil, fmt.Errorf("template metadata cannot be nil")
	}

	fileMetadata := &FileMetadata{
		TemplateMetadata: templateMetadata,
		FileName:         filepath.Base(filePath),
		RelativePath:     filePath,
	}

	// Determine file type
	fileMetadata.FileType = mp.determineFileType(filePath)

	// Add language-specific metadata
	mp.addLanguageSpecificMetadata(fileMetadata)

	return fileMetadata, nil
}

// createSlug creates a URL-friendly slug from a project name
func (mp *MetadataProcessor) createSlug(name string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// createTitle creates a display-friendly title from a project name
func (mp *MetadataProcessor) createTitle(name string) string {
	// If the name is already properly formatted, return as-is
	if strings.Contains(name, " ") {
		return name
	}

	// Convert camelCase or PascalCase to Title Case
	var result strings.Builder
	for i, char := range name {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune(' ')
		}
		if i == 0 {
			result.WriteRune(char)
		} else {
			result.WriteRune(char)
		}
	}

	title := result.String()

	// Replace hyphens and underscores with spaces
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// extractModuleName extracts a Go module name from repository URL
func (mp *MetadataProcessor) extractModuleName(repository, projectSlug string) string {
	if repository == "" {
		return projectSlug
	}

	// Handle GitHub URLs
	if strings.Contains(repository, "github.com") {
		// Extract from URL: https://github.com/user/repo -> github.com/user/repo
		if strings.HasPrefix(repository, "https://github.com/") {
			return strings.TrimPrefix(repository, "https://")
		}
		if strings.HasPrefix(repository, "http://github.com/") {
			return strings.TrimPrefix(repository, "http://")
		}
		if strings.Contains(repository, "github.com/") {
			parts := strings.Split(repository, "github.com/")
			if len(parts) > 1 {
				return "github.com/" + parts[1]
			}
		}
	}

	// Handle other Git hosting services
	for _, host := range []string{"gitlab.com", "bitbucket.org", "codeberg.org"} {
		if strings.Contains(repository, host) {
			if strings.HasPrefix(repository, "https://"+host+"/") {
				return strings.TrimPrefix(repository, "https://")
			}
			if strings.HasPrefix(repository, "http://"+host+"/") {
				return strings.TrimPrefix(repository, "http://")
			}
		}
	}

	// Fallback to project slug
	return projectSlug
}

// determineFileType determines the file type based on extension
func (mp *MetadataProcessor) determineFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	typeMap := map[string]string{
		".go":   "go",
		".js":   "javascript",
		".ts":   "typescript",
		".py":   "python",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".rs":   "rust",
		".php":  "php",
		".rb":   "ruby",
		".md":   "markdown",
		".yml":  "yaml",
		".yaml": "yaml",
		".json": "json",
		".xml":  "xml",
		".html": "html",
		".css":  "css",
		".scss": "scss",
		".sass": "sass",
		".sql":  "sql",
		".sh":   "shell",
		".bat":  "batch",
		".ps1":  "powershell",
	}

	if fileType, exists := typeMap[ext]; exists {
		return fileType
	}

	return "text"
}

// addLanguageSpecificMetadata adds language-specific metadata
func (mp *MetadataProcessor) addLanguageSpecificMetadata(fileMetadata *FileMetadata) {
	switch fileMetadata.FileType {
	case "go":
		// For Go files, determine package name from path
		fileMetadata.GoPackage = mp.determineGoPackage(fileMetadata.RelativePath)

	case "javascript", "typescript":
		// For JS/TS files, use the module name
		fileMetadata.JSModule = fileMetadata.PackageName

	case "python":
		// For Python files, determine module name
		fileMetadata.PythonModule = mp.determinePythonModule(fileMetadata.RelativePath)
	}
}

// determineGoPackage determines the Go package name from file path
func (mp *MetadataProcessor) determineGoPackage(filePath string) string {
	dir := filepath.Dir(filePath)

	// Handle special cases
	if dir == "." || dir == "" {
		return "main"
	}

	// Use the directory name as package name
	packageName := filepath.Base(dir)

	// Handle nested packages
	if packageName == "internal" || packageName == "pkg" {
		// Look at parent directory
		parent := filepath.Dir(dir)
		if parent != "." && parent != "" {
			packageName = filepath.Base(parent)
		}
	}

	return packageName
}

// determinePythonModule determines the Python module name from file path
func (mp *MetadataProcessor) determinePythonModule(filePath string) string {
	// Convert file path to Python module notation
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return ""
	}

	// Replace path separators with dots
	module := strings.ReplaceAll(dir, string(filepath.Separator), ".")

	return module
}

// ValidateMetadata validates the processed metadata
func (mp *MetadataProcessor) ValidateMetadata(metadata *TemplateMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}

	// Validate required fields
	if metadata.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if metadata.Author == "" {
		return fmt.Errorf("author is required")
	}

	// Validate computed fields
	if metadata.ProjectSlug == "" {
		return fmt.Errorf("project slug could not be generated")
	}

	if metadata.CopyrightYear == "" {
		return fmt.Errorf("copyright year could not be determined")
	}

	return nil
}

// GetMetadataForTemplate returns metadata formatted for template usage
func (mp *MetadataProcessor) GetMetadataForTemplate(metadata *TemplateMetadata) map[string]interface{} {
	result := make(map[string]interface{})

	// Convert struct to map for template usage
	result["Name"] = metadata.ProjectName
	result["Description"] = metadata.ProjectDescription
	result["Slug"] = metadata.ProjectSlug
	result["Title"] = metadata.ProjectTitle
	result["Author"] = metadata.Author
	result["Email"] = metadata.Email
	result["Organization"] = metadata.Organization
	result["License"] = metadata.License
	result["Copyright"] = metadata.Copyright
	result["CopyrightYear"] = metadata.CopyrightYear
	result["Repository"] = metadata.Repository
	result["RepositoryURL"] = metadata.RepositoryURL
	result["ModuleName"] = metadata.ModuleName
	result["PackageName"] = metadata.PackageName
	result["GeneratedAt"] = metadata.GeneratedAt
	result["GeneratorVersion"] = metadata.GeneratorVersion
	result["AuthorWithEmail"] = metadata.AuthorWithEmail
	result["LicenseFile"] = metadata.LicenseFile
	result["HasRepository"] = metadata.HasRepository
	result["HasOrganization"] = metadata.HasOrganization

	return result
}
