package template

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	texttemplate "text/template"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/version"
)

// Engine implements the TemplateEngine interface
type Engine struct {
	textTemplate   *texttemplate.Template
	htmlTemplate   *htmltemplate.Template
	funcMap        texttemplate.FuncMap
	versionManager interfaces.VersionManager
}

// NewEngine creates a new template engine instance
func NewEngine() interfaces.TemplateEngine {
	engine := &Engine{
		funcMap: make(texttemplate.FuncMap),
	}

	// Register default template functions
	engine.registerDefaultFunctions()

	return engine
}

// NewEngineWithVersionManager creates a new template engine with version management
func NewEngineWithVersionManager(versionManager interfaces.VersionManager) interfaces.TemplateEngine {
	engine := &Engine{
		funcMap:        make(texttemplate.FuncMap),
		versionManager: versionManager,
	}

	// Register default template functions
	engine.registerDefaultFunctions()

	return engine
}

// ProcessTemplate processes a single template file with the given configuration
func (e *Engine) ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Enhance config with centralized version information
	enhancedConfig, err := e.enhanceConfigWithVersions(config)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance config with versions: %w", err)
	}

	// Load the template
	tmpl, err := e.LoadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", templatePath, err)
	}

	// Render the template
	return e.RenderTemplate(tmpl, enhancedConfig)
}

// ProcessDirectory processes an entire template directory recursively
func (e *Engine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Enhance config with centralized version information once for the entire directory
	enhancedConfig, err := e.enhanceConfigWithVersions(config)
	if err != nil {
		return fmt.Errorf("failed to enhance config with versions: %w", err)
	}

	return filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}

		// Calculate relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Calculate output path
		outputPath := filepath.Join(outputDir, relPath)

		if d.IsDir() {
			// Create directory in output
			return os.MkdirAll(outputPath, 0755)
		}

		// Process file with enhanced config
		return e.processFile(path, outputPath, enhancedConfig)
	})
}

// RegisterFunctions registers custom template functions
func (e *Engine) RegisterFunctions(funcMap texttemplate.FuncMap) {
	for name, fn := range funcMap {
		e.funcMap[name] = fn
	}
}

// LoadTemplate loads and parses a template from the given path
func (e *Engine) LoadTemplate(templatePath string) (*texttemplate.Template, error) {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template with custom functions
	tmpl := texttemplate.New(filepath.Base(templatePath)).Funcs(e.funcMap)

	// Parse template content
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func (e *Engine) RenderTemplate(tmpl *texttemplate.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// processFile processes a single file (template or binary)
func (e *Engine) processFile(srcPath, destPath string, config *models.ProjectConfig) error {
	// Check if file is a template (has .tmpl extension)
	if strings.HasSuffix(srcPath, ".tmpl") {
		// Remove .tmpl extension from destination
		destPath = strings.TrimSuffix(destPath, ".tmpl")

		// Process as template
		content, err := e.ProcessTemplate(srcPath, config)
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", srcPath, err)
		}

		// Write processed content
		return os.WriteFile(destPath, content, 0644)
	}

	// Copy binary file as-is
	return e.copyFile(srcPath, destPath)
}

// copyFile copies a file from src to dest
func (e *Engine) copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

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

	return os.Chmod(dest, srcInfo.Mode())
}

// enhanceConfigWithVersions merges centralized version information into the project config
func (e *Engine) enhanceConfigWithVersions(config *models.ProjectConfig) (*models.ProjectConfig, error) {
	// If no version manager is configured, return config as-is
	if e.versionManager == nil {
		return config, nil
	}

	// Create a copy of the config to avoid modifying the original
	enhancedConfig := *config
	if enhancedConfig.Versions == nil {
		enhancedConfig.Versions = &models.VersionConfig{
			Packages: make(map[string]string),
		}
	} else {
		// Deep copy the versions
		enhancedVersions := *enhancedConfig.Versions
		enhancedVersions.Packages = make(map[string]string)
		for k, v := range config.Versions.Packages {
			enhancedVersions.Packages[k] = v
		}
		enhancedConfig.Versions = &enhancedVersions
	}

	// Check if version manager supports storage (enhanced functionality)
	if managerWithStorage, ok := e.versionManager.(*version.Manager); ok {
		store, err := managerWithStorage.GetVersionStore()
		if err != nil {
			// Log warning but don't fail - use existing versions
			fmt.Printf("Warning: Could not load version store: %v\n", err)
			return &enhancedConfig, nil
		}

		// Merge language versions
		for name, info := range store.Languages {
			switch name {
			case "nodejs":
				enhancedConfig.Versions.Node = info.CurrentVersion
			case "go":
				enhancedConfig.Versions.Go = info.CurrentVersion
			case "kotlin":
				enhancedConfig.Versions.Kotlin = info.CurrentVersion
			case "swift":
				enhancedConfig.Versions.Swift = info.CurrentVersion
			}
		}

		// Merge framework versions
		for name, info := range store.Frameworks {
			switch name {
			case "nextjs":
				enhancedConfig.Versions.NextJS = info.CurrentVersion
			case "react":
				enhancedConfig.Versions.React = info.CurrentVersion
			}
		}

		// Merge package versions
		for name, info := range store.Packages {
			enhancedConfig.Versions.Packages[name] = info.CurrentVersion
		}

		enhancedConfig.Versions.UpdatedAt = store.LastUpdated
	} else {
		// Fallback to basic version manager functionality
		if nodeVersion, err := e.versionManager.GetLatestNodeVersion(); err == nil {
			enhancedConfig.Versions.Node = nodeVersion
		}
		if goVersion, err := e.versionManager.GetLatestGoVersion(); err == nil {
			enhancedConfig.Versions.Go = goVersion
		}
		if nextjsVersion, err := e.versionManager.GetLatestNPMPackage("next"); err == nil {
			enhancedConfig.Versions.NextJS = nextjsVersion
		}
		if reactVersion, err := e.versionManager.GetLatestNPMPackage("react"); err == nil {
			enhancedConfig.Versions.React = reactVersion
		}

		// Add common packages
		commonPackages := []string{
			"typescript", "tailwindcss", "eslint", "prettier", "jest",
			"@types/node", "@types/react", "autoprefixer", "postcss",
		}
		for _, pkg := range commonPackages {
			if version, err := e.versionManager.GetLatestNPMPackage(pkg); err == nil {
				enhancedConfig.Versions.Packages[pkg] = version
			}
		}
	}

	return &enhancedConfig, nil
}
