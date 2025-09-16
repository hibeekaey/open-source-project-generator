package interfaces

import (
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateEngine defines the contract for template processing operations.
//
// The TemplateEngine interface provides comprehensive template processing capabilities
// for the Open Source Project Generator. It handles:
//   - Single template file processing with variable substitution
//   - Recursive directory processing for complete project generation
//   - Custom function registration for extended template functionality
//   - Version management integration for automatic dependency updates
//
// Implementations should provide:
//   - Robust error handling and validation
//   - Security considerations for template processing
//   - Integration with version management systems
type TemplateEngine interface {
	// ProcessTemplate processes a single template file with the given configuration.
	//
	// This method:
	//   - Loads the template from the specified path
	//   - Applies the project configuration as template variables
	//   - Performs variable substitution and conditional rendering
	//   - Returns the processed content as bytes
	//
	// Parameters:
	//   - templatePath: Absolute or relative path to the template file
	//   - config: Project configuration containing variables and settings
	//
	// Returns:
	//   - []byte: Processed template content ready for writing to file
	//   - error: Any error that occurred during processing
	//
	// The method validates the template syntax, applies security checks,
	// and ensures all required variables are available in the configuration.
	ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error)

	// ProcessDirectory processes an entire template directory recursively.
	//
	// This method:
	//   - Walks through all files and subdirectories in the template directory
	//   - Processes template files (.tmpl extension) using ProcessTemplate
	//   - Copies non-template files directly to the output directory
	//   - Maintains directory structure in the output location
	//   - Handles file permissions and metadata preservation
	//
	// Parameters:
	//   - templateDir: Path to the root template directory
	//   - outputDir: Path where processed files should be written
	//   - config: Project configuration for template processing
	//
	// Returns:
	//   - error: Any error that occurred during directory processing
	//
	// The method creates the output directory structure as needed and
	// provides progress feedback for long-running operations.
	ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error

	// RegisterFunctions registers custom template functions for use in templates.
	//
	// This method allows extending template functionality with custom functions
	// that can be called from within template files. Common use cases include:
	//   - String manipulation and formatting functions
	//   - Date and time formatting
	//   - Mathematical operations
	//   - Custom validation and transformation logic
	//
	// Parameters:
	//   - funcMap: Map of function names to function implementations
	//
	// The functions should be thread-safe and handle edge cases gracefully.
	// Function names should not conflict with built-in template functions.
	RegisterFunctions(funcMap template.FuncMap)

	// LoadTemplate loads and parses a template from the given path.
	//
	// This method:
	//   - Reads the template file from disk
	//   - Parses the template syntax and validates correctness
	//   - Applies registered custom functions
	//   - Returns a parsed template ready for rendering
	//
	// Parameters:
	//   - templatePath: Path to the template file to load
	//
	// Returns:
	//   - *template.Template: Parsed template ready for execution
	//   - error: Any error that occurred during loading or parsing
	//
	// The method validates template syntax and provides detailed error
	// messages for debugging template issues.
	LoadTemplate(templatePath string) (*template.Template, error)

	// RenderTemplate renders a template with the provided data.
	//
	// This method:
	//   - Executes the parsed template with the given data context
	//   - Performs variable substitution and function calls
	//   - Handles conditional rendering and loops
	//   - Returns the rendered content as bytes
	//
	// Parameters:
	//   - tmpl: Parsed template to render (from LoadTemplate)
	//   - data: Data context for template variable substitution
	//
	// Returns:
	//   - []byte: Rendered template content
	//   - error: Any error that occurred during rendering
	//
	// The method provides detailed error messages including line numbers
	// for template execution errors and validates data types.
	RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error)
}
