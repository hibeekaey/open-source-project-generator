package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// DirectoryProcessor handles recursive template directory processing
type DirectoryProcessor struct {
	engine   *Engine
	metadata *MetadataParser
}

// NewDirectoryProcessor creates a new directory processor
func NewDirectoryProcessor(engine *Engine) *DirectoryProcessor {
	return &DirectoryProcessor{
		engine:   engine,
		metadata: NewMetadataParser(),
	}
}

// ProcessTemplateDirectory processes an entire template directory with conditional rendering
func (p *DirectoryProcessor) ProcessTemplateDirectory(templateDir, outputDir string, config *models.ProjectConfig) error {
	// Parse template metadata if available
	metadata, err := p.metadata.ParseMetadataFromDirectory(templateDir)
	if err != nil {
		return fmt.Errorf("failed to parse template metadata: %w", err)
	}

	// Process directory recursively
	return filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}

		// Skip metadata files and base templates
		if p.isMetadataFile(d.Name()) || p.isBaseTemplate(path, templateDir) {
			return nil
		}

		// Calculate relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Skip disabled template files
		if strings.HasSuffix(path, ".tmpl.disabled") {
			return nil
		}

		// Check if this path should be processed based on conditions
		shouldProcess, err := p.shouldProcessPath(relPath, metadata, config)
		if err != nil {
			return fmt.Errorf("failed to evaluate conditions for %s: %w", relPath, err)
		}

		if !shouldProcess {
			// Skip this path and its children if it's a directory
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate output path (process template variables in path names)
		outputPath, err := p.processPathTemplate(filepath.Join(outputDir, relPath), config)
		if err != nil {
			return fmt.Errorf("failed to process path template %s: %w", relPath, err)
		}

		if d.IsDir() {
			// Create directory in output with secure permissions
			return utils.SafeMkdirAll(outputPath)
		}

		// Process file
		return p.processFileWithInheritance(path, outputPath, config, metadata)
	})
}

// processFileWithInheritance processes a file with template inheritance support
func (p *DirectoryProcessor) processFileWithInheritance(srcPath, destPath string, config *models.ProjectConfig, _ *TemplateMetadata) error {
	// Check if file is a template
	if strings.HasSuffix(srcPath, ".tmpl") {
		// Remove .tmpl extension from destination
		destPath = strings.TrimSuffix(destPath, ".tmpl")

		// Process template with inheritance
		content, err := p.processTemplateWithInheritance(srcPath, config)
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", srcPath, err)
		}

		// Create destination directory if it doesn't exist
		destDir := filepath.Dir(destPath)
		if err := utils.SafeMkdirAll(destDir); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}

		// Write processed content with secure permissions
		return utils.SafeWriteFile(destPath, content)
	}

	// Handle binary files and assets
	return p.copyAsset(srcPath, destPath)
}

// processTemplateWithInheritance processes a template file with inheritance and partial includes
func (p *DirectoryProcessor) processTemplateWithInheritance(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Read template content with path validation
	content, err := utils.SafeReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	templateContent := string(content)

	// Process template inheritance (extends directive)
	if strings.Contains(templateContent, "{{/* extends ") {
		templateContent, err = p.processTemplateExtends(templateContent, templatePath, config)
		if err != nil {
			return nil, fmt.Errorf("failed to process template extends: %w", err)
		}
	}

	// Process partial includes
	templateContent, err = p.processTemplateIncludes(templateContent, templatePath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to process template includes: %w", err)
	}

	// Create and parse the final template
	tmpl := texttemplate.New(filepath.Base(templatePath)).Funcs(p.engine.funcMap)
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Render the template
	return p.engine.RenderTemplate(tmpl, config)
}

// processTemplateExtends handles template inheritance with extends directive
func (p *DirectoryProcessor) processTemplateExtends(content, templatePath string, config *models.ProjectConfig) (string, error) {
	// Find extends directive: {{/* extends "base.tmpl" */}}
	lines := strings.Split(content, "\n")
	var baseTemplate string
	var childContent strings.Builder
	var blocks = make(map[string]string)

	// Parse child template for extends directive and named blocks
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{{/* extends ") && strings.HasSuffix(trimmed, " */}}") {
			// Extract base template name
			start := strings.Index(trimmed, `"`) + 1
			end := strings.LastIndex(trimmed, `"`)
			if start > 0 && end > start {
				baseTemplate = trimmed[start:end]
			}
		} else if strings.Contains(line, "{{/* ") && strings.Contains(line, " */}}") {
			// Check for named block definitions: {{/* blockname */}}content{{/* /blockname */}}
			blockStart := strings.Index(line, "{{/* ")
			blockEnd := strings.Index(line, " */}}")
			if blockStart >= 0 && blockEnd > blockStart {
				blockName := strings.TrimSpace(line[blockStart+5 : blockEnd])

				// Look for block end marker
				blockEndMarker := "{{/* /" + blockName + " */}}"
				blockContent := strings.Builder{}

				// Extract content after the block start marker
				afterMarker := line[blockEnd+4:]
				if endIdx := strings.Index(afterMarker, blockEndMarker); endIdx >= 0 {
					// Single line block
					blockContent.WriteString(afterMarker[:endIdx])
					blocks[blockName] = blockContent.String()
				} else {
					// Multi-line block - collect until end marker
					blockContent.WriteString(afterMarker)
					if i < len(lines)-1 {
						blockContent.WriteString("\n")
					}

					for j := i + 1; j < len(lines); j++ {
						if strings.Contains(lines[j], blockEndMarker) {
							// Found end marker
							endIdx := strings.Index(lines[j], blockEndMarker)
							blockContent.WriteString(lines[j][:endIdx])
							blocks[blockName] = blockContent.String()
							// Note: i will be updated by the loop to continue processing
							break
						} else {
							blockContent.WriteString(lines[j])
							if j < len(lines)-1 {
								blockContent.WriteString("\n")
							}
						}
					}
				}
			}
		} else if i > 0 || !strings.HasPrefix(trimmed, "{{/* extends ") {
			// Add non-extends, non-block lines to child content
			if !strings.Contains(line, "{{/* ") || !strings.Contains(line, " */}}") {
				childContent.WriteString(line)
				if i < len(lines)-1 {
					childContent.WriteString("\n")
				}
			}
		}
	}

	if baseTemplate == "" {
		return content, nil // No extends directive found
	}

	// Resolve base template path - make it relative to template directory root
	templateDir := filepath.Dir(templatePath)
	for strings.Contains(templateDir, string(filepath.Separator)) && !strings.HasSuffix(templateDir, "template") {
		templateDir = filepath.Dir(templateDir)
	}

	var baseTemplatePath string
	if filepath.IsAbs(baseTemplate) {
		baseTemplatePath = baseTemplate
	} else {
		// Try relative to current template directory first
		baseTemplatePath = filepath.Join(filepath.Dir(templatePath), baseTemplate)
		if _, err := os.Stat(baseTemplatePath); os.IsNotExist(err) {
			// Try relative to template root directory
			baseTemplatePath = filepath.Join(templateDir, baseTemplate)
		}
	}

	// Read base template with path validation
	baseContent, err := utils.SafeReadFile(baseTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read base template %s: %w", baseTemplatePath, err)
	}

	// Process base template recursively if it also extends another template
	baseStr := string(baseContent)
	if strings.Contains(baseStr, "{{/* extends ") {
		baseStr, err = p.processTemplateExtends(baseStr, baseTemplatePath, config)
		if err != nil {
			return "", fmt.Errorf("failed to process base template inheritance: %w", err)
		}
	}

	// Replace placeholders in base template
	// First replace named blocks
	for blockName, blockContent := range blocks {
		placeholder := "{{/* " + blockName + " */}}"
		baseStr = strings.ReplaceAll(baseStr, placeholder, blockContent)
	}

	// Then replace generic content placeholder
	contentPlaceholder := "{{/* content */}}"
	if strings.Contains(baseStr, contentPlaceholder) {
		baseStr = strings.ReplaceAll(baseStr, contentPlaceholder, childContent.String())
	} else if childContent.Len() > 0 {
		// If no content placeholder, append child content
		baseStr += "\n" + childContent.String()
	}

	return baseStr, nil
}

// processTemplateIncludes handles partial template includes
func (p *DirectoryProcessor) processTemplateIncludes(content, templatePath string, _ *models.ProjectConfig) (string, error) {
	// Process include directives: {{/* include "partial.tmpl" */}}
	for {
		start := strings.Index(content, "{{/* include ")
		if start == -1 {
			break // No more includes
		}

		end := strings.Index(content[start:], " */}}")
		if end == -1 {
			return "", fmt.Errorf("malformed include directive in template")
		}
		end += start + 4 // Include the " */}}" part

		// Extract include directive
		includeDirective := content[start:end]

		// Extract partial template name
		quotStart := strings.Index(includeDirective, `"`) + 1
		quotEnd := strings.LastIndex(includeDirective, `"`)
		if quotStart <= 0 || quotEnd <= quotStart {
			return "", fmt.Errorf("malformed include directive: %s", includeDirective)
		}

		partialName := includeDirective[quotStart:quotEnd]

		// Resolve partial template path - make it relative to template directory root
		templateDir := filepath.Dir(templatePath)
		for strings.Contains(templateDir, string(filepath.Separator)) && !strings.HasSuffix(templateDir, "template") {
			templateDir = filepath.Dir(templateDir)
		}

		var partialPath string
		if filepath.IsAbs(partialName) {
			partialPath = partialName
		} else {
			// Try relative to current template directory first
			partialPath = filepath.Join(filepath.Dir(templatePath), partialName)
			if _, err := os.Stat(partialPath); os.IsNotExist(err) {
				// Try relative to template root directory
				partialPath = filepath.Join(templateDir, partialName)
			}
		}

		// Read partial template with path validation
		partialContent, err := utils.SafeReadFile(partialPath)
		if err != nil {
			return "", fmt.Errorf("failed to read partial template %s: %w", partialPath, err)
		}

		// Replace include directive with partial content
		content = content[:start] + string(partialContent) + content[end:]
	}

	return content, nil
}

// copyAsset copies binary files and assets
func (p *DirectoryProcessor) copyAsset(srcPath, destPath string) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := utils.SafeMkdirAll(destDir); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file with path validation
	srcFile, err := utils.SafeOpen(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Create destination file with secure permissions
	destFile, err := utils.SafeCreate(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	// Copy file content
	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	return os.Chmod(destPath, srcInfo.Mode())
}

// shouldProcessPath determines if a path should be processed based on conditions
func (p *DirectoryProcessor) shouldProcessPath(relPath string, metadata *TemplateMetadata, config *models.ProjectConfig) (bool, error) {
	// Check file-specific conditions
	for _, file := range metadata.Files {
		if file.Source == relPath || strings.HasSuffix(relPath, file.Source) {
			return p.metadata.EvaluateConditions(file.Conditions, config)
		}
	}

	// Check global conditions
	return p.metadata.EvaluateConditions(metadata.Conditions, config)
}

// processPathTemplate processes template variables in path names
func (p *DirectoryProcessor) processPathTemplate(path string, config *models.ProjectConfig) (string, error) {
	// Basic path template processing - replace common variables
	result := path

	// Replace common template variables in paths
	replacements := map[string]string{
		"{{.Name}}":                 config.Name,
		"{{.Name | lower}}":         strings.ToLower(config.Name),
		"{{.Name | upper}}":         strings.ToUpper(config.Name),
		"{{kebabCase .Name}}":       toKebabCase(config.Name),
		"{{snakeCase .Name}}":       toSnakeCase(config.Name),
		"{{.Organization}}":         config.Organization,
		"{{.Organization | lower}}": strings.ToLower(config.Organization),
		"{{.Organization | upper}}": strings.ToUpper(config.Organization),
	}

	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result, nil
}

// isMetadataFile checks if a file is a template metadata file
func (p *DirectoryProcessor) isMetadataFile(filename string) bool {
	metadataFiles := []string{
		"template.yaml",
		"template.yml",
		"metadata.yaml",
		"metadata.yml",
		".template.yaml",
		".template.yml",
	}

	for _, metaFile := range metadataFiles {
		if filename == metaFile {
			return true
		}
	}

	return false
}

// isBaseTemplate checks if a file is a base template that should not be copied to output
func (p *DirectoryProcessor) isBaseTemplate(filePath, templateDir string) bool {
	// Check if this file is referenced as a base template by other files
	relPath, err := filepath.Rel(templateDir, filePath)
	if err != nil {
		return false
	}

	// Common base template patterns
	basePatterns := []string{
		"base.",
		"layout.",
		"layouts/",
		"partials/",
		"_",
	}

	for _, pattern := range basePatterns {
		if strings.Contains(relPath, pattern) {
			return true
		}
	}

	return false
}
