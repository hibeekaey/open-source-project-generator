package filesystem

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem/components"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem/generators"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem/operations"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem/processors"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectGenerator handles complete project structure generation
type ProjectGenerator struct {
	fsGen                   *Generator
	structureManager        *StructureManager
	fsOps                   *FileSystemOperations
	frontendGenerator       *components.FrontendGenerator
	backendGenerator        *components.BackendGenerator
	mobileGenerator         *components.MobileGenerator
	infrastructureGenerator *components.InfrastructureGenerator

	// New focused generators
	structureGen *generators.StructureGenerator
	templateGen  *generators.TemplateGenerator
	configGen    *generators.ConfigurationGenerator
	docGen       *generators.DocumentationGenerator
	cicdGen      *generators.CICDGenerator

	// Processors
	templateProcessor *processors.TemplateProcessor
	fileProcessor     *processors.FileProcessor

	// Operations
	creator   *operations.Creator
	validator *operations.Validator
}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator() *ProjectGenerator {
	fsOps := NewFileSystemOperations()
	return &ProjectGenerator{
		fsGen:                   NewGenerator().(*Generator),
		structureManager:        NewStructureManager(),
		fsOps:                   fsOps,
		frontendGenerator:       components.NewFrontendGenerator(fsOps),
		backendGenerator:        components.NewBackendGenerator(fsOps),
		mobileGenerator:         components.NewMobileGenerator(fsOps),
		infrastructureGenerator: components.NewInfrastructureGenerator(fsOps),

		// Initialize new components
		structureGen:      generators.NewStructureGenerator(fsOps),
		templateGen:       generators.NewTemplateGenerator(fsOps),
		configGen:         generators.NewConfigurationGenerator(fsOps),
		docGen:            generators.NewDocumentationGenerator(fsOps),
		cicdGen:           generators.NewCICDGenerator(fsOps),
		templateProcessor: processors.NewTemplateProcessor(fsOps),
		fileProcessor:     processors.NewFileProcessor(fsOps),
		creator:           operations.NewCreator(fsOps),
		validator:         operations.NewValidator(fsOps),
	}
}

// NewDryRunProjectGenerator creates a new project generator in dry-run mode
func NewDryRunProjectGenerator() *ProjectGenerator {
	fsOps := NewDryRunFileSystemOperations()
	return &ProjectGenerator{
		fsGen:                   NewDryRunGenerator().(*Generator),
		structureManager:        NewDryRunStructureManager(),
		fsOps:                   fsOps,
		frontendGenerator:       components.NewFrontendGenerator(fsOps),
		backendGenerator:        components.NewBackendGenerator(fsOps),
		mobileGenerator:         components.NewMobileGenerator(fsOps),
		infrastructureGenerator: components.NewInfrastructureGenerator(fsOps),

		// Initialize new components for dry-run
		structureGen:      generators.NewStructureGenerator(fsOps),
		templateGen:       generators.NewTemplateGenerator(fsOps),
		configGen:         generators.NewConfigurationGenerator(fsOps),
		docGen:            generators.NewDocumentationGenerator(fsOps),
		cicdGen:           generators.NewCICDGenerator(fsOps),
		templateProcessor: processors.NewTemplateProcessor(fsOps),
		fileProcessor:     processors.NewFileProcessor(fsOps),
		creator:           operations.NewCreator(fsOps),
		validator:         operations.NewValidator(fsOps),
	}
}

// GenerateProjectStructure creates the complete project directory structure
func (pg *ProjectGenerator) GenerateProjectStructure(config *models.ProjectConfig, outputPath string) error {
	// Use the new structure generator
	if err := pg.structureGen.GenerateProjectStructure(config, outputPath); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Create the root project directory using FileSystemOperations
	projectPath, err := pg.creator.CreateProjectRoot(config, outputPath)
	if err != nil {
		return fmt.Errorf("failed to create project root: %w", err)
	}

	// Use StructureManager to create directories
	if err := pg.structureManager.CreateDirectories(projectPath, config); err != nil {
		return fmt.Errorf("failed to create project directories: %w", err)
	}

	return nil
}

// GenerateComponentFiles creates component-specific files based on user selection
func (pg *ProjectGenerator) GenerateComponentFiles(config *models.ProjectConfig, outputPath string) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	projectPath := filepath.Join(outputPath, config.Name)

	// Generate root configuration files using template generator
	if err := pg.templateGen.GenerateRootFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate root files: %w", err)
	}

	// Generate component-specific files using processors
	if err := pg.templateProcessor.ProcessComponentFiles(projectPath, config,
		pg.frontendGenerator, pg.backendGenerator, pg.mobileGenerator, pg.infrastructureGenerator); err != nil {
		return fmt.Errorf("failed to generate component files: %w", err)
	}

	// Generate configuration files
	if err := pg.generateConfigurationFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate configuration files: %w", err)
	}

	// Generate documentation files
	if err := pg.docGen.GenerateDocumentationFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate documentation files: %w", err)
	}

	// Generate CI/CD files
	if err := pg.cicdGen.GenerateCICDFiles(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate CI/CD files: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates the generated project structure
func (pg *ProjectGenerator) ValidateProjectStructure(projectPath string, config *models.ProjectConfig) error {
	// Use the new structure generator for validation
	if err := pg.structureGen.ValidateProjectStructure(projectPath, config); err != nil {
		return fmt.Errorf("structure validation failed: %w", err)
	}

	// Use StructureManager to validate project structure
	return pg.structureManager.ValidateProjectStructure(projectPath, config)
}

// ValidateCrossReferences validates cross-references between generated files
func (pg *ProjectGenerator) ValidateCrossReferences(projectPath string, config *models.ProjectConfig) error {
	// Validate project root using validator
	if err := pg.validator.ValidateProjectRoot(projectPath, config); err != nil {
		return fmt.Errorf("project root validation failed: %w", err)
	}

	// Validate that required configuration files exist
	requiredFiles := []string{
		"Makefile",
		"README.md",
		"docker-compose.yml",
		".gitignore",
	}

	if err := pg.validator.ValidateRequiredFiles(projectPath, requiredFiles); err != nil {
		return fmt.Errorf("required root files validation failed: %w", err)
	}

	// Validate component-specific cross-references using file processor
	if err := pg.fileProcessor.ValidateComponentCrossReferences(projectPath, config); err != nil {
		return fmt.Errorf("component cross-reference validation failed: %w", err)
	}

	// Validate content cross-references using file processor
	if err := pg.fileProcessor.ValidateContentCrossReferences(projectPath, config); err != nil {
		return fmt.Errorf("content cross-reference validation failed: %w", err)
	}

	return nil
}

// generateConfigurationFiles generates configuration files for all components
func (pg *ProjectGenerator) generateConfigurationFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate frontend configuration files
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		if err := pg.configGen.GenerateFrontendFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate frontend configuration files: %w", err)
		}
	}

	// Generate backend configuration files
	if config.Components.Backend.GoGin {
		if err := pg.configGen.GenerateBackendFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate backend configuration files: %w", err)
		}
	}

	// Generate mobile configuration files
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		if err := pg.configGen.GenerateMobileFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate mobile configuration files: %w", err)
		}
	}

	// Generate infrastructure configuration files
	if config.Components.Infrastructure.Terraform {
		if err := pg.configGen.GenerateInfrastructureFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate infrastructure configuration files: %w", err)
		}
	}

	return nil
}

// Legacy methods for backward compatibility with tests
// These methods were refactored into the generators package but are kept for test compatibility

// generateMakefileContent generates Makefile content
func (pg *ProjectGenerator) generateMakefileContent(config *models.ProjectConfig) string {
	return fmt.Sprintf("# Makefile for %s\n\nhelp:\n\t@echo \"Available targets:\"\n\nsetup:\n\t@echo \"Setting up %s\"\n\ndev:\n\t@echo \"Starting development server\"\n\ntest:\n\t@echo \"Running tests\"\n\nbuild:\n\t@echo \"Building %s\"", config.Name, config.Name, config.Name)
}

// generateReadmeContent generates README content
func (pg *ProjectGenerator) generateReadmeContent(config *models.ProjectConfig) string {
	return fmt.Sprintf("# %s\n\n%s\n\nOrganization: %s\n\n## Getting Started\n\nWelcome to %s!\n\n## Project Structure\n\nThis project follows standard conventions.", config.Name, config.Description, config.Organization, config.Name)
}

// generateDockerComposeContent generates docker-compose content
func (pg *ProjectGenerator) generateDockerComposeContent(config *models.ProjectConfig) string {
	return fmt.Sprintf("version: '3.8'\n\nservices:\n  %s:\n    build: .\n    ports:\n      - \"3000:3000\"", config.Name)
}

// generateGitignoreContent generates .gitignore content
func (pg *ProjectGenerator) generateGitignoreContent(config *models.ProjectConfig) string {
	return "node_modules/\n*.log\n.env\n.DS_Store\n*.tmp\n*.temp\ndist/\nbuild/"
}
