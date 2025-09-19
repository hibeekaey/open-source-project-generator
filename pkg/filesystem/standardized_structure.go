package filesystem

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
)

// StandardizedStructureGenerator handles creation of standardized directory structures
// according to the interactive CLI generation requirements using the template system
type StandardizedStructureGenerator struct {
	fsGen           *Generator
	templateManager interfaces.TemplateManager
	templateEngine  interfaces.TemplateEngine
}

// Ensure StandardizedStructureGenerator implements the interface
var _ interfaces.StandardizedStructureGenerator = (*StandardizedStructureGenerator)(nil)

// NewStandardizedStructureGenerator creates a new standardized structure generator
func NewStandardizedStructureGenerator() *StandardizedStructureGenerator {
	engine := template.NewEmbeddedEngine()
	manager := template.NewManager(engine)

	return &StandardizedStructureGenerator{
		fsGen:           NewGenerator().(*Generator),
		templateManager: manager,
		templateEngine:  engine,
	}
}

// NewDryRunStandardizedStructureGenerator creates a new standardized structure generator in dry-run mode
func NewDryRunStandardizedStructureGenerator() *StandardizedStructureGenerator {
	engine := template.NewEmbeddedEngine()
	manager := template.NewManager(engine)

	return &StandardizedStructureGenerator{
		fsGen:           NewDryRunGenerator().(*Generator),
		templateManager: manager,
		templateEngine:  engine,
	}
}

// GenerateStandardizedStructure creates the complete standardized project structure using templates
func (ssg *StandardizedStructureGenerator) GenerateStandardizedStructure(config *models.ProjectConfig, outputPath string) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Create the root project directory
	projectPath := filepath.Join(outputPath, config.Name)
	if err := ssg.fsGen.EnsureDirectory(projectPath); err != nil {
		return fmt.Errorf("failed to create project root directory: %w", err)
	}

	// Always process base templates first for common project structure
	if err := ssg.processBaseTemplate(projectPath, config); err != nil {
		return fmt.Errorf("failed to process base template: %w", err)
	}

	if err := ssg.processTemplate("github", filepath.Join(projectPath, ".github"), config); err != nil {
		return fmt.Errorf("failed to process github template: %w", err)
	}

	if err := ssg.processTemplate("scripts", filepath.Join(projectPath, "scripts"), config); err != nil {
		return fmt.Errorf("failed to process scripts template: %w", err)
	}

	// Process frontend templates based on selected components
	if err := ssg.processFrontendTemplates(projectPath, config); err != nil {
		return fmt.Errorf("failed to process frontend templates: %w", err)
	}

	// Process backend template if selected
	if config.Components.Backend.GoGin {
		if err := ssg.processTemplate("go-gin", filepath.Join(projectPath, "CommonServer"), config); err != nil {
			return fmt.Errorf("failed to process backend template: %w", err)
		}
	}

	// Process mobile templates if selected
	if err := ssg.processMobileTemplates(projectPath, config); err != nil {
		return fmt.Errorf("failed to process mobile templates: %w", err)
	}

	// Process infrastructure templates if selected
	if err := ssg.processInfrastructureTemplates(projectPath, config); err != nil {
		return fmt.Errorf("failed to process infrastructure templates: %w", err)
	}

	return nil
}

// CreateFrontendDirectoryStructure creates App/ directory with main/, home/, admin/, shared-components/ subdirectories
// Implements conditional directory creation based on selected frontend templates
// Adds proper file structure for Next.js, React, and TypeScript components
func (ssg *StandardizedStructureGenerator) CreateFrontendDirectoryStructure(projectPath string, config *models.ProjectConfig) error {
	return ssg.processFrontendTemplates(projectPath, config)
}

// processFrontendTemplates processes frontend templates based on configuration
func (ssg *StandardizedStructureGenerator) processFrontendTemplates(projectPath string, config *models.ProjectConfig) error {
	// Process main App template if selected
	if config.Components.Frontend.NextJS.App {
		appPath := filepath.Join(projectPath, "App")
		if err := ssg.processTemplate("nextjs-app", appPath, config); err != nil {
			return fmt.Errorf("failed to process nextjs-app template: %w", err)
		}
	}

	// Process Home template if selected
	if config.Components.Frontend.NextJS.Home {
		homePath := filepath.Join(projectPath, "Home")
		if err := ssg.processTemplate("nextjs-home", homePath, config); err != nil {
			return fmt.Errorf("failed to process nextjs-home template: %w", err)
		}
	}

	// Process Admin template if selected
	if config.Components.Frontend.NextJS.Admin {
		adminPath := filepath.Join(projectPath, "Admin")
		if err := ssg.processTemplate("nextjs-admin", adminPath, config); err != nil {
			return fmt.Errorf("failed to process nextjs-admin template: %w", err)
		}
	}

	// Process shared components template if selected
	if config.Components.Frontend.NextJS.Shared {
		sharedPath := filepath.Join(projectPath, "shared-components")
		if err := ssg.processTemplate("shared-components", sharedPath, config); err != nil {
			return fmt.Errorf("failed to process shared-components template: %w", err)
		}
	}

	return nil
}

// hasFrontendComponents checks if any frontend components are selected
func (ssg *StandardizedStructureGenerator) hasFrontendComponents(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared
}

// hasMobileComponents checks if any mobile components are selected
func (ssg *StandardizedStructureGenerator) hasMobileComponents(config *models.ProjectConfig) bool {
	return config.Components.Mobile.Android || config.Components.Mobile.IOS
}

// hasInfrastructureComponents checks if any infrastructure components are selected
func (ssg *StandardizedStructureGenerator) hasInfrastructureComponents(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

// CreateBackendDirectoryStructure creates CommonServer/ directory with cmd/, internal/, pkg/, migrations/, docs/ structure
// Implements Go project structure with proper package organization
// Adds API documentation and migration file generation
func (ssg *StandardizedStructureGenerator) CreateBackendDirectoryStructure(projectPath string, config *models.ProjectConfig) error {
	if config.Components.Backend.GoGin {
		backendPath := filepath.Join(projectPath, "CommonServer")
		return ssg.processTemplate("go-gin", backendPath, config)
	}
	return nil
}

// CreateMobileDirectoryStructure creates Mobile/ directory with android/, ios/, shared/ subdirectories
// Adds platform-specific project structures for Kotlin and Swift
// Implements shared resources and API specification handling
func (ssg *StandardizedStructureGenerator) CreateMobileDirectoryStructure(projectPath string, config *models.ProjectConfig) error {
	return ssg.processMobileTemplates(projectPath, config)
}

// processMobileTemplates processes mobile templates based on configuration
func (ssg *StandardizedStructureGenerator) processMobileTemplates(projectPath string, config *models.ProjectConfig) error {
	// Process Android template if selected
	if config.Components.Mobile.Android {
		androidPath := filepath.Join(projectPath, "Mobile", "android")
		if err := ssg.processTemplate("android-kotlin", androidPath, config); err != nil {
			return fmt.Errorf("failed to process android-kotlin template: %w", err)
		}
	}

	// Process iOS template if selected
	if config.Components.Mobile.IOS {
		iosPath := filepath.Join(projectPath, "Mobile", "ios")
		if err := ssg.processTemplate("ios-swift", iosPath, config); err != nil {
			return fmt.Errorf("failed to process ios-swift template: %w", err)
		}
	}

	// Process shared mobile template if any mobile component is selected
	if ssg.hasMobileComponents(config) {
		sharedPath := filepath.Join(projectPath, "Mobile", "shared")
		if err := ssg.processTemplate("shared", sharedPath, config); err != nil {
			return fmt.Errorf("failed to process mobile shared template: %w", err)
		}
	}

	return nil
}

// CreateInfrastructureDirectoryStructure creates Deploy/ directory with docker/, k8s/, terraform/, monitoring/ subdirectories
// Adds configuration files for Docker, Kubernetes, Terraform, and monitoring tools
// Implements infrastructure template processing and file generation
func (ssg *StandardizedStructureGenerator) CreateInfrastructureDirectoryStructure(projectPath string, config *models.ProjectConfig) error {
	return ssg.processInfrastructureTemplates(projectPath, config)
}

// processInfrastructureTemplates processes infrastructure templates based on configuration
func (ssg *StandardizedStructureGenerator) processInfrastructureTemplates(projectPath string, config *models.ProjectConfig) error {
	deployPath := filepath.Join(projectPath, "Deploy")

	// Process Docker template if selected
	if config.Components.Infrastructure.Docker {
		dockerPath := filepath.Join(deployPath, "docker")
		if err := ssg.processTemplate("docker", dockerPath, config); err != nil {
			return fmt.Errorf("failed to process docker template: %w", err)
		}
	}

	// Process Kubernetes template if selected
	if config.Components.Infrastructure.Kubernetes {
		k8sPath := filepath.Join(deployPath, "k8s")
		if err := ssg.processTemplate("kubernetes", k8sPath, config); err != nil {
			return fmt.Errorf("failed to process kubernetes template: %w", err)
		}
	}

	// Process Terraform template if selected
	if config.Components.Infrastructure.Terraform {
		terraformPath := filepath.Join(deployPath, "terraform")
		if err := ssg.processTemplate("terraform", terraformPath, config); err != nil {
			return fmt.Errorf("failed to process terraform template: %w", err)
		}
	}

	return nil
}

// CreateCommonDirectoryStructure always creates Docs/, Scripts/, and .github/ directories with appropriate content
// Generates standard project files (README.md, CONTRIBUTING.md, LICENSE, .gitignore, Makefile)
// Implements CI/CD workflow generation for GitHub Actions
func (ssg *StandardizedStructureGenerator) CreateCommonDirectoryStructure(projectPath string, config *models.ProjectConfig) error {
	// The base template already handles common directory structure
	// This method is kept for interface compatibility
	return nil
}

// GenerateStandardProjectFiles creates standard project files with appropriate content
func (ssg *StandardizedStructureGenerator) GenerateStandardProjectFiles(projectPath string, config *models.ProjectConfig) error {
	// Standard project files are generated by the base template
	// This method is kept for interface compatibility and can be used for additional file generation
	return nil
}

// processTemplate processes a template using the template manager
func (ssg *StandardizedStructureGenerator) processTemplate(templateName, outputPath string, config *models.ProjectConfig) error {
	return ssg.templateManager.ProcessTemplate(templateName, config, outputPath)
}

// processBaseTemplate processes the base template directory directly using the template engine
func (ssg *StandardizedStructureGenerator) processBaseTemplate(outputPath string, config *models.ProjectConfig) error {
	// Process the base template directory directly
	return ssg.templateEngine.ProcessDirectory("templates/base", outputPath, config)
}

// Helper methods for checking component selections are defined above
