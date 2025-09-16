// Package template provides template processing capabilities for the Open Source Template Generator.
//
// This package implements the TemplateEngine interface and provides:
//   - Template file loading and parsing
//   - Template rendering with variable substitution
//   - Directory processing for complete project generation
//   - Custom function registration for extended template functionality
//
// The template engine ensures that generated projects are correctly processed
// and all template variables are properly substituted.
package template

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/utils"
)

// Engine implements the TemplateEngine interface
type Engine struct {
	funcMap template.FuncMap
}

// NewEngine creates a new template engine
func NewEngine() interfaces.TemplateEngine {
	engine := &Engine{
		funcMap: make(template.FuncMap),
	}
	engine.registerDefaultFunctions()
	return engine
}

// NewEngineWithVersionManager creates a new template engine with version manager
func NewEngineWithVersionManager(versionManager interfaces.VersionManager) interfaces.TemplateEngine {
	engine := &Engine{
		funcMap: make(template.FuncMap),
	}
	engine.registerDefaultFunctions()
	return engine
}

// ProcessTemplate processes a single template file with the given configuration
func (e *Engine) ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Load the template
	tmpl, err := e.LoadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// Render the template
	content, err := e.RenderTemplate(tmpl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return content, nil
}

// ProcessDirectory processes an entire template directory recursively
func (e *Engine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Walk through template directory
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
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

		// Calculate output path
		outputPath := filepath.Join(outputDir, relPath)

		// Remove .tmpl extension if present
		outputPath = strings.TrimSuffix(outputPath, ".tmpl")

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(outputPath, info.Mode())
		}

		// Process file
		if strings.HasSuffix(path, ".tmpl") {
			// Process template file
			content, err := e.ProcessTemplate(path, config)
			if err != nil {
				return fmt.Errorf("failed to process template %s: %w", path, err)
			}

			// Write processed content
			return os.WriteFile(outputPath, content, info.Mode())
		} else {
			// Copy non-template file
			return e.copyFile(path, outputPath, info.Mode())
		}
	})
}

// RegisterFunctions registers custom template functions
func (e *Engine) RegisterFunctions(funcMap template.FuncMap) {
	for name, fn := range funcMap {
		e.funcMap[name] = fn
	}
}

// LoadTemplate loads and parses a template from the given path
func (e *Engine) LoadTemplate(templatePath string) (*template.Template, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(templatePath); err != nil {
		return nil, fmt.Errorf("invalid template path: %w", err)
	}

	content, err := utils.SafeReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func (e *Engine) RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// Helper function to copy a file
func (e *Engine) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}
