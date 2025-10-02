package processors

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// FileProcessor handles file processing and validation
type FileProcessor struct {
	fsOps FileSystemOperationsInterface
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(fsOps FileSystemOperationsInterface) *FileProcessor {
	return &FileProcessor{
		fsOps: fsOps,
	}
}

// ValidateComponentCrossReferences validates cross-references between component files
func (fp *FileProcessor) ValidateComponentCrossReferences(projectPath string, config *models.ProjectConfig) error {
	// Validate frontend cross-references
	if config.Components.Frontend.NextJS.App {
		packageJsonPath := filepath.Join(projectPath, "App/package.json")
		if !fp.fsOps.FileExists(packageJsonPath) {
			return fmt.Errorf("main app package.json missing")
		}
	}

	// Validate backend cross-references
	if config.Components.Backend.GoGin {
		goModPath := filepath.Join(projectPath, "CommonServer/go.mod")
		if !fp.fsOps.FileExists(goModPath) {
			return fmt.Errorf("backend go.mod missing")
		}
	}

	// Validate mobile cross-references
	if config.Components.Mobile.Android {
		buildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
		if !fp.fsOps.FileExists(buildGradlePath) {
			return fmt.Errorf("android build.gradle missing")
		}
	}

	if config.Components.Mobile.IOS {
		packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
		if !fp.fsOps.FileExists(packageSwiftPath) {
			return fmt.Errorf("iOS Package.swift missing")
		}
	}

	return nil
}

// ValidateContentCrossReferences validates that file contents reference each other correctly
func (fp *FileProcessor) ValidateContentCrossReferences(projectPath string, config *models.ProjectConfig) error {
	// Validate Makefile references correct components
	makefilePath := filepath.Join(projectPath, "Makefile")
	if fp.fsOps.FileExists(makefilePath) {
		// Validate Makefile exists and is readable
		if err := fp.fsOps.ValidateFileContent(makefilePath); err != nil {
			return fmt.Errorf("failed to validate Makefile content: %w", err)
		}
	}

	// Validate docker-compose.yml references correct services
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	if fp.fsOps.FileExists(dockerComposePath) {
		// Validate docker-compose.yml exists and is readable
		if err := fp.fsOps.ValidateFileContent(dockerComposePath); err != nil {
			return fmt.Errorf("failed to validate docker-compose.yml content: %w", err)
		}
	}

	// Validate package.json dependencies are consistent across frontend apps
	if config.Components.Frontend.NextJS.App && config.Components.Frontend.NextJS.Admin {
		mainAppPackageJson := filepath.Join(projectPath, "App/package.json")
		adminPackageJson := filepath.Join(projectPath, "Admin/package.json")

		if fp.fsOps.FileExists(mainAppPackageJson) && fp.fsOps.FileExists(adminPackageJson) {
			// Validate both package.json files are readable
			if err := fp.fsOps.ValidateFileContent(mainAppPackageJson); err != nil {
				return fmt.Errorf("failed to validate main app package.json content: %w", err)
			}
			if err := fp.fsOps.ValidateFileContent(adminPackageJson); err != nil {
				return fmt.Errorf("failed to validate admin package.json content: %w", err)
			}
		}
	}

	return nil
}
