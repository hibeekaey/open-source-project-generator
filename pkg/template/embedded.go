package template

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

//go:embed templates
var embeddedTemplates embed.FS

// EmbeddedEngine implements the TemplateEngine interface using embedded templates
type EmbeddedEngine struct {
	funcMap template.FuncMap
	fs      fs.FS
}

// NewEmbeddedEngine creates a new template engine with embedded templates
func NewEmbeddedEngine() interfaces.TemplateEngine {
	engine := &EmbeddedEngine{
		funcMap: make(template.FuncMap),
		fs:      embeddedTemplates,
	}
	engine.registerDefaultFunctions()
	return engine
}

// ProcessTemplate processes a single template file with the given configuration
func (e *EmbeddedEngine) ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Normalize path for embedded filesystem
	templatePath = strings.TrimPrefix(templatePath, "./")
	if !strings.HasPrefix(templatePath, "templates/") {
		templatePath = "templates/" + templatePath
	}

	// Load the template from embedded filesystem
	tmpl, err := e.LoadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded template: %w", err)
	}

	// Render the template
	content, err := e.RenderTemplate(tmpl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return content, nil
}

// ProcessDirectory processes an entire template directory recursively from embedded filesystem
func (e *EmbeddedEngine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Create output directory if it doesn't exist with secure permissions
	if err := utils.SafeMkdirAll(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Normalize template directory path
	templateDir = strings.TrimPrefix(templateDir, "./")
	if !strings.HasPrefix(templateDir, "templates/") {
		templateDir = "templates/" + templateDir
	}

	// Walk through embedded template directory
	return fs.WalkDir(e.fs, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip if it's the template directory itself
		if relPath == "." {
			return nil
		}

		// Skip disabled template files
		if strings.HasSuffix(path, ".tmpl.disabled") {
			return nil
		}

		// Calculate output path with template variable processing
		outputPath := filepath.Join(outputDir, relPath)

		// Process template variables in path names
		outputPath, err = e.processPathTemplate(outputPath, config)
		if err != nil {
			return fmt.Errorf("failed to process path template %s: %w", relPath, err)
		}

		// Remove .tmpl extension if present
		outputPath = strings.TrimSuffix(outputPath, ".tmpl")

		// Handle special cases for files that should start with dots
		outputPath = e.restoreHiddenFileName(outputPath)

		if d.IsDir() {
			// Create directory with secure permissions
			return utils.SafeMkdirAll(outputPath)
		}

		// Process file
		if strings.HasSuffix(path, ".tmpl") {
			// Process template file
			content, err := e.ProcessTemplate(path, config)
			if err != nil {
				return fmt.Errorf("failed to process embedded template %s: %w", path, err)
			}

			// Write processed content with secure permissions
			return utils.SafeWriteFile(outputPath, content)
		} else {
			// Copy non-template file from embedded filesystem
			return e.copyEmbeddedFile(path, outputPath)
		}
	})
}

// RegisterFunctions registers custom template functions
func (e *EmbeddedEngine) RegisterFunctions(funcMap template.FuncMap) {
	for name, fn := range funcMap {
		e.funcMap[name] = fn
	}
}

// LoadTemplate loads and parses a template from the embedded filesystem
func (e *EmbeddedEngine) LoadTemplate(templatePath string) (*template.Template, error) {
	content, err := fs.ReadFile(e.fs, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded template file: %w", err)
	}

	// Create template with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func (e *EmbeddedEngine) RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return []byte(buf.String()), nil
}

// Helper function to copy a file from embedded filesystem
func (e *EmbeddedEngine) copyEmbeddedFile(src, dst string) error {
	content, err := fs.ReadFile(e.fs, src)
	if err != nil {
		return fmt.Errorf("failed to read embedded file: %w", err)
	}

	// Restore hidden file name for output
	dst = e.restoreHiddenFileName(dst)

	// Create directory if it doesn't exist with secure permissions
	if err := utils.SafeMkdirAll(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return utils.SafeWriteFile(dst, content)
}

// processPathTemplate processes template variables in path names
func (e *EmbeddedEngine) processPathTemplate(path string, config *models.ProjectConfig) (string, error) {
	// Simple path template processing - replace common variables
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

// registerDefaultFunctions registers the default template functions
func (e *EmbeddedEngine) registerDefaultFunctions() {
	// Use the same comprehensive function map from functions.go
	engine := &Engine{}
	engine.registerDefaultFunctions()
	e.funcMap = engine.funcMap
}

// restoreHiddenFileName restores the original hidden file names
func (e *EmbeddedEngine) restoreHiddenFileName(outputPath string) string {
	dir := filepath.Dir(outputPath)
	filename := filepath.Base(outputPath)

	// Map of renamed files back to their original names
	hiddenFiles := map[string]string{
		"gitignore":         ".gitignore",
		"prettierrc":        ".prettierrc",
		"eslintrc.json":     ".eslintrc.json",
		"env.local.example": ".env.local.example",
		"env.example":       ".env.example",
		"env.test":          ".env.test",
		"golangci.yml":      ".golangci.yml",
		"dockerignore":      ".dockerignore",
		"gitkeep":           ".gitkeep",
	}

	if originalName, exists := hiddenFiles[filename]; exists {
		return filepath.Join(dir, originalName)
	}

	return outputPath
}
