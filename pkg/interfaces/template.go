package interfaces

import (
	"text/template"

	"github.com/open-source-template-generator/pkg/models"
)

// TemplateEngine defines the contract for template processing operations
type TemplateEngine interface {
	// ProcessTemplate processes a single template file with the given configuration
	ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error)

	// ProcessDirectory processes an entire template directory recursively
	ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error

	// RegisterFunctions registers custom template functions
	RegisterFunctions(funcMap template.FuncMap)

	// LoadTemplate loads and parses a template from the given path
	LoadTemplate(templatePath string) (*template.Template, error)

	// RenderTemplate renders a template with the provided data
	RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error)
}
