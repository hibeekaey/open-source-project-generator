package processors

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateProcessor handles template processing and file generation
type TemplateProcessor struct {
	fsOps FileSystemOperationsInterface
}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor(fsOps FileSystemOperationsInterface) *TemplateProcessor {
	return &TemplateProcessor{
		fsOps: fsOps,
	}
}

// ProcessComponentFiles creates component-specific files
func (tp *TemplateProcessor) ProcessComponentFiles(projectPath string, config *models.ProjectConfig, frontendGen, backendGen, mobileGen, infraGen ComponentGenerator) error {
	// Generate frontend component files
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin || config.Components.Frontend.NextJS.Shared {
		if err := frontendGen.GenerateFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate frontend files: %w", err)
		}
	}

	// Generate backend component files
	if config.Components.Backend.GoGin {
		if err := backendGen.GenerateFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate backend files: %w", err)
		}
	}

	// Generate mobile component files
	if config.Components.Mobile.Android || config.Components.Mobile.IOS || config.Components.Mobile.Shared {
		if err := mobileGen.GenerateFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate mobile files: %w", err)
		}
	}

	// Generate infrastructure component files
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		if err := infraGen.GenerateFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate infrastructure files: %w", err)
		}
	}

	return nil
}

// ComponentGenerator interface for component generators
type ComponentGenerator interface {
	GenerateFiles(projectPath string, config *models.ProjectConfig) error
}
