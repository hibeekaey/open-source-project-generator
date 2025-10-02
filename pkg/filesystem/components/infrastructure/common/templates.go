package common

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateManager handles common template operations for infrastructure components
type TemplateManager struct{}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{}
}

// ValidateConfig validates common infrastructure configuration
func (tm *TemplateManager) ValidateConfig(config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}
	return nil
}

// GetProjectName returns the sanitized project name for infrastructure use
func (tm *TemplateManager) GetProjectName(config *models.ProjectConfig) string {
	if config == nil {
		return "unknown"
	}
	return config.Name
}
