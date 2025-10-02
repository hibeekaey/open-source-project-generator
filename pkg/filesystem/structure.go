package filesystem

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ProjectStructure defines the standard project directory structure
type ProjectStructure struct {
	// Root directories
	RootDirs []string
	// Component-specific directories
	FrontendDirs []string
	BackendDirs  []string
	MobileDirs   []string
	InfraDirs    []string
	// Common directories
	CommonDirs []string
}

// StructureManager manages project structure operations
type StructureManager struct {
	structure *ProjectStructure
	fsGen     *Generator
}

// NewStructureManager creates a new structure manager
func NewStructureManager() *StructureManager {
	return &StructureManager{
		structure: GetStandardProjectStructure(),
		fsGen:     NewGenerator().(*Generator),
	}
}

// NewDryRunStructureManager creates a new structure manager in dry-run mode
func NewDryRunStructureManager() *StructureManager {
	return &StructureManager{
		structure: GetStandardProjectStructure(),
		fsGen:     NewDryRunGenerator().(*Generator),
	}
}

// GetStandardProjectStructure returns the standard project directory structure
func GetStandardProjectStructure() *ProjectStructure {
	return &ProjectStructure{
		RootDirs: []string{
			"docs",
			"scripts",
			".github/workflows",
			".github/ISSUE_TEMPLATE",
			".github/PULL_REQUEST_TEMPLATE",
		},
		FrontendDirs: []string{
			"App/src/components/ui",
			"App/src/components/forms",
			"App/src/hooks",
			"App/src/context",
			"App/src/lib",
			"App/src/types",
			"App/public",
			"Home/src/components",
			"Home/src/sections",
			"Home/public",
			"Admin/src/components",
			"Admin/src/pages",
			"Admin/src/hooks",
			"Admin/public",
		},
		BackendDirs: []string{
			"CommonServer/cmd",
			"CommonServer/internal/controllers",
			"CommonServer/internal/models",
			"CommonServer/internal/services",
			"CommonServer/internal/middleware",
			"CommonServer/internal/repository",
			"CommonServer/internal/config",
			"CommonServer/pkg/auth",
			"CommonServer/pkg/database",
			"CommonServer/pkg/utils",
			"CommonServer/migrations",
			"CommonServer/tests",
		},
		MobileDirs: []string{
			"Mobile/Android/app/src/main/java",
			"Mobile/Android/app/src/main/res",
			"Mobile/Android/app/src/test/java",
			"Mobile/iOS/Sources",
			"Mobile/iOS/Resources",
			"Mobile/iOS/Tests",
			"Mobile/Shared/api",
			"Mobile/Shared/assets",
		},
		InfraDirs: []string{
			"Deploy/terraform/modules",
			"Deploy/terraform/environments/staging",
			"Deploy/terraform/environments/production",
			"Deploy/kubernetes/base",
			"Deploy/kubernetes/overlays/staging",
			"Deploy/kubernetes/overlays/production",
			"Deploy/docker",
			"Deploy/helm",
		},
		CommonDirs: []string{
			"Tests/integration",
			"Tests/e2e",
			"Tests/performance",
		},
	}
}

// GetStructure returns the current project structure
func (sm *StructureManager) GetStructure() *ProjectStructure {
	return sm.structure
}

// SetStructure sets a custom project structure
func (sm *StructureManager) SetStructure(structure *ProjectStructure) error {
	if structure == nil {
		return fmt.Errorf("project structure cannot be nil")
	}

	if err := sm.ValidateStructure(structure); err != nil {
		return fmt.Errorf("invalid project structure: %w", err)
	}

	sm.structure = structure
	return nil
}

// ValidateStructure validates a project structure definition
func (sm *StructureManager) ValidateStructure(structure *ProjectStructure) error {
	if structure == nil {
		return fmt.Errorf("project structure cannot be nil")
	}

	// Validate that at least root directories are defined
	if len(structure.RootDirs) == 0 {
		return fmt.Errorf("project structure must define at least one root directory")
	}

	// Validate directory paths don't contain invalid characters
	allDirs := append(structure.RootDirs, structure.CommonDirs...)
	allDirs = append(allDirs, structure.FrontendDirs...)
	allDirs = append(allDirs, structure.BackendDirs...)
	allDirs = append(allDirs, structure.MobileDirs...)
	allDirs = append(allDirs, structure.InfraDirs...)

	for _, dir := range allDirs {
		if dir == "" {
			return fmt.Errorf("directory path cannot be empty")
		}

		// Check for invalid path characters
		if filepath.IsAbs(dir) {
			return fmt.Errorf("directory paths must be relative: %s", dir)
		}

		// Check for path traversal attempts
		cleanPath := filepath.Clean(dir)
		if cleanPath != dir || strings.Contains(dir, "..") || strings.HasPrefix(dir, "/") {
			return fmt.Errorf("directory path contains invalid elements: %s", dir)
		}
	}

	return nil
}

// CreateDirectories creates directories based on the project structure and configuration
func (sm *StructureManager) CreateDirectories(projectPath string, config *models.ProjectConfig) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Create root directories (always created)
	if err := sm.createDirectories(projectPath, sm.structure.RootDirs); err != nil {
		return fmt.Errorf("failed to create root directories: %w", err)
	}

	// Create common directories (always created)
	if err := sm.createDirectories(projectPath, sm.structure.CommonDirs); err != nil {
		return fmt.Errorf("failed to create common directories: %w", err)
	}

	// Create component-specific directories based on configuration
	if err := sm.createComponentDirectories(projectPath, config); err != nil {
		return fmt.Errorf("failed to create component directories: %w", err)
	}

	return nil
}

// ValidateProjectStructure validates that a generated project structure matches expectations
func (sm *StructureManager) ValidateProjectStructure(projectPath string, config *models.ProjectConfig) error {
	if projectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}

	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Validate root directories exist
	for _, dir := range sm.structure.RootDirs {
		dirPath := filepath.Join(projectPath, dir)
		if !sm.fsGen.FileExists(dirPath) {
			return fmt.Errorf("required root directory missing: %s", dir)
		}
	}

	// Validate common directories exist
	for _, dir := range sm.structure.CommonDirs {
		dirPath := filepath.Join(projectPath, dir)
		if !sm.fsGen.FileExists(dirPath) {
			return fmt.Errorf("required common directory missing: %s", dir)
		}
	}

	// Validate component-specific directories based on configuration
	if err := sm.validateComponentDirectories(projectPath, config); err != nil {
		return fmt.Errorf("component directory validation failed: %w", err)
	}

	return nil
}

// CustomizeStructure allows customization of the project structure
func (sm *StructureManager) CustomizeStructure(customization *StructureCustomization) error {
	if customization == nil {
		return fmt.Errorf("structure customization cannot be nil")
	}

	// Apply additional directories
	if len(customization.AdditionalRootDirs) > 0 {
		sm.structure.RootDirs = append(sm.structure.RootDirs, customization.AdditionalRootDirs...)
	}

	if len(customization.AdditionalFrontendDirs) > 0 {
		sm.structure.FrontendDirs = append(sm.structure.FrontendDirs, customization.AdditionalFrontendDirs...)
	}

	if len(customization.AdditionalBackendDirs) > 0 {
		sm.structure.BackendDirs = append(sm.structure.BackendDirs, customization.AdditionalBackendDirs...)
	}

	if len(customization.AdditionalMobileDirs) > 0 {
		sm.structure.MobileDirs = append(sm.structure.MobileDirs, customization.AdditionalMobileDirs...)
	}

	if len(customization.AdditionalInfraDirs) > 0 {
		sm.structure.InfraDirs = append(sm.structure.InfraDirs, customization.AdditionalInfraDirs...)
	}

	// Remove excluded directories
	sm.structure.RootDirs = sm.filterDirectories(sm.structure.RootDirs, customization.ExcludedDirs)
	sm.structure.FrontendDirs = sm.filterDirectories(sm.structure.FrontendDirs, customization.ExcludedDirs)
	sm.structure.BackendDirs = sm.filterDirectories(sm.structure.BackendDirs, customization.ExcludedDirs)
	sm.structure.MobileDirs = sm.filterDirectories(sm.structure.MobileDirs, customization.ExcludedDirs)
	sm.structure.InfraDirs = sm.filterDirectories(sm.structure.InfraDirs, customization.ExcludedDirs)
	sm.structure.CommonDirs = sm.filterDirectories(sm.structure.CommonDirs, customization.ExcludedDirs)

	// Validate the customized structure
	return sm.ValidateStructure(sm.structure)
}

// StructureCustomization defines customizations to the project structure
type StructureCustomization struct {
	AdditionalRootDirs     []string
	AdditionalFrontendDirs []string
	AdditionalBackendDirs  []string
	AdditionalMobileDirs   []string
	AdditionalInfraDirs    []string
	ExcludedDirs           []string
}

// createDirectories creates a list of directories under the project path
func (sm *StructureManager) createDirectories(projectPath string, dirs []string) error {
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := sm.fsGen.EnsureDirectory(dirPath); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// createComponentDirectories creates directories based on selected components
func (sm *StructureManager) createComponentDirectories(projectPath string, config *models.ProjectConfig) error {
	// Create frontend directories if any frontend component is selected
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		if err := sm.createDirectories(projectPath, sm.structure.FrontendDirs); err != nil {
			return fmt.Errorf("failed to create frontend directories: %w", err)
		}
	}

	// Create backend directories if backend is selected
	if config.Components.Backend.GoGin {
		if err := sm.createDirectories(projectPath, sm.structure.BackendDirs); err != nil {
			return fmt.Errorf("failed to create backend directories: %w", err)
		}
	}

	// Create mobile directories if any mobile component is selected
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		if err := sm.createDirectories(projectPath, sm.structure.MobileDirs); err != nil {
			return fmt.Errorf("failed to create mobile directories: %w", err)
		}
	}

	// Create infrastructure directories if any infrastructure component is selected
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		if err := sm.createDirectories(projectPath, sm.structure.InfraDirs); err != nil {
			return fmt.Errorf("failed to create infrastructure directories: %w", err)
		}
	}

	return nil
}

// validateComponentDirectories validates component-specific directories
func (sm *StructureManager) validateComponentDirectories(projectPath string, config *models.ProjectConfig) error {
	// Validate frontend directories
	if config.Components.Frontend.NextJS.App || config.Components.Frontend.NextJS.Home || config.Components.Frontend.NextJS.Admin {
		for _, dir := range sm.structure.FrontendDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !sm.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required frontend directory missing: %s", dir)
			}
		}
	}

	// Validate backend directories
	if config.Components.Backend.GoGin {
		for _, dir := range sm.structure.BackendDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !sm.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required backend directory missing: %s", dir)
			}
		}
	}

	// Validate mobile directories
	if config.Components.Mobile.Android || config.Components.Mobile.IOS {
		for _, dir := range sm.structure.MobileDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !sm.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required mobile directory missing: %s", dir)
			}
		}
	}

	// Validate infrastructure directories
	if config.Components.Infrastructure.Terraform || config.Components.Infrastructure.Kubernetes || config.Components.Infrastructure.Docker {
		for _, dir := range sm.structure.InfraDirs {
			dirPath := filepath.Join(projectPath, dir)
			if !sm.fsGen.FileExists(dirPath) {
				return fmt.Errorf("required infrastructure directory missing: %s", dir)
			}
		}
	}

	return nil
}

// filterDirectories removes excluded directories from a list
func (sm *StructureManager) filterDirectories(dirs []string, excluded []string) []string {
	if len(excluded) == 0 {
		return dirs
	}

	excludedMap := make(map[string]bool)
	for _, dir := range excluded {
		excludedMap[dir] = true
	}

	filtered := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if !excludedMap[dir] {
			filtered = append(filtered, dir)
		}
	}

	return filtered
}
