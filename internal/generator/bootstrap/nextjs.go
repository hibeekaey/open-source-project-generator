package bootstrap

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// NextJSExecutor handles Next.js project generation using create-next-app
type NextJSExecutor struct {
	*BaseExecutor
}

// NewNextJSExecutor creates a new Next.js executor
func NewNextJSExecutor() *NextJSExecutor {
	return &NextJSExecutor{
		BaseExecutor: NewBaseExecutor("npx"),
	}
}

// Execute generates a Next.js project using create-next-app
func (ne *NextJSExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Build Next.js specific command
	args, err := ne.buildNextJSCommand(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build Next.js command: %w", err)
	}

	// Create a new spec with the built arguments
	execSpec := &BootstrapSpec{
		ComponentType: spec.ComponentType,
		TargetDir:     spec.TargetDir,
		Config:        spec.Config,
		Flags:         args,
		Timeout:       spec.Timeout,
	}

	// Execute using base executor
	result, err := ne.BaseExecutor.Execute(ctx, execSpec)
	if err != nil {
		return result, fmt.Errorf("Next.js generation failed: %w", err)
	}

	// Set output directory to the generated project
	if projectName, ok := spec.Config["name"].(string); ok {
		result.OutputDir = filepath.Join(spec.TargetDir, projectName)
	}

	return result, nil
}

// SupportsComponent checks if this executor supports the given component type
func (ne *NextJSExecutor) SupportsComponent(componentType string) bool {
	return componentType == "nextjs" || componentType == "next" || componentType == "frontend"
}

// GetDefaultFlags returns default flags for Next.js generation
func (ne *NextJSExecutor) GetDefaultFlags(componentType string) []string {
	if !ne.SupportsComponent(componentType) {
		return []string{}
	}

	return []string{
		"create-next-app@latest",
		"--typescript",
		"--tailwind",
		"--app",
		"--eslint",
		"--no-git",
	}
}

// ValidateConfig validates component-specific configuration
func (ne *NextJSExecutor) ValidateConfig(config map[string]interface{}) error {
	// Validate name
	if name, ok := config["name"].(string); !ok || name == "" {
		return fmt.Errorf("name is required and must be a string")
	}

	// Validate boolean fields
	boolFields := []string{"typescript", "tailwind", "app_router", "eslint", "src_dir"}
	for _, field := range boolFields {
		if val, exists := config[field]; exists {
			if _, ok := val.(bool); !ok {
				return fmt.Errorf("%s must be a boolean value", field)
			}
		}
	}

	// Validate import_alias if provided
	if alias, exists := config["import_alias"]; exists {
		if _, ok := alias.(string); !ok {
			return fmt.Errorf("import_alias must be a string")
		}
	}

	return nil
}

// buildNextJSCommand builds the create-next-app command with appropriate flags
func (ne *NextJSExecutor) buildNextJSCommand(spec *BootstrapSpec) ([]string, error) {
	args := []string{"create-next-app@latest"}

	// Get project name from config
	projectName, ok := spec.Config["name"].(string)
	if !ok || projectName == "" {
		return nil, fmt.Errorf("project name is required in config")
	}
	args = append(args, projectName)

	// Add TypeScript flag (default: true)
	useTypeScript := true
	if ts, ok := spec.Config["typescript"].(bool); ok {
		useTypeScript = ts
	}
	if useTypeScript {
		args = append(args, "--typescript")
	} else {
		args = append(args, "--javascript")
	}

	// Add Tailwind CSS flag (default: true)
	useTailwind := true
	if tw, ok := spec.Config["tailwind"].(bool); ok {
		useTailwind = tw
	}
	if useTailwind {
		args = append(args, "--tailwind")
	} else {
		args = append(args, "--no-tailwind")
	}

	// Add App Router flag (default: true)
	useAppRouter := true
	if app, ok := spec.Config["app_router"].(bool); ok {
		useAppRouter = app
	}
	if useAppRouter {
		args = append(args, "--app")
	} else {
		args = append(args, "--no-app")
	}

	// Add ESLint flag (default: true)
	useESLint := true
	if eslint, ok := spec.Config["eslint"].(bool); ok {
		useESLint = eslint
	}
	if useESLint {
		args = append(args, "--eslint")
	} else {
		args = append(args, "--no-eslint")
	}

	// Add src directory flag (default: false)
	useSrcDir := false
	if src, ok := spec.Config["src_dir"].(bool); ok {
		useSrcDir = src
	}
	if useSrcDir {
		args = append(args, "--src-dir")
	} else {
		args = append(args, "--no-src-dir")
	}

	// Add import alias flag (default: @/*)
	if alias, ok := spec.Config["import_alias"].(string); ok && alias != "" {
		args = append(args, "--import-alias", alias)
	}

	// Disable git initialization (we'll handle it at project level)
	args = append(args, "--no-git")

	// Add any additional custom flags
	if len(spec.Flags) > 0 {
		// Validate flags for security
		if err := ne.ValidateFlags(spec.Flags); err != nil {
			return nil, fmt.Errorf("invalid flags: %w", err)
		}
		args = append(args, spec.Flags...)
	}

	return args, nil
}

// GetManualSteps returns manual steps required after Next.js generation
func (ne *NextJSExecutor) GetManualSteps(spec *BootstrapSpec) []string {
	steps := []string{
		"Navigate to the project directory",
		"Run 'npm install' to install dependencies (if not already done)",
		"Run 'npm run dev' to start the development server",
		"Open http://localhost:3000 in your browser",
	}

	// Add custom steps based on configuration
	if useTailwind, ok := spec.Config["tailwind"].(bool); ok && useTailwind {
		steps = append(steps, "Customize Tailwind configuration in tailwind.config.js")
	}

	return steps
}
