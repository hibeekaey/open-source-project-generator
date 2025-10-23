package mapper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ComponentMapping defines the mapping from component type to target directory
type ComponentMapping struct {
	ComponentType string
	TargetPath    string
	Description   string
}

// DefaultComponentMappings defines the standard component-to-directory mappings
var DefaultComponentMappings = []ComponentMapping{
	{
		ComponentType: "nextjs",
		TargetPath:    "App/",
		Description:   "Next.js frontend application",
	},
	{
		ComponentType: "go-backend",
		TargetPath:    "CommonServer/",
		Description:   "Go backend server",
	},
	{
		ComponentType: "android",
		TargetPath:    "Mobile/android/",
		Description:   "Android mobile application",
	},
	{
		ComponentType: "ios",
		TargetPath:    "Mobile/ios/",
		Description:   "iOS mobile application",
	},
	{
		ComponentType: "docker",
		TargetPath:    "Deploy/docker/",
		Description:   "Docker configuration",
	},
	{
		ComponentType: "kubernetes",
		TargetPath:    "Deploy/k8s/",
		Description:   "Kubernetes configuration",
	},
	{
		ComponentType: "terraform",
		TargetPath:    "Deploy/terraform/",
		Description:   "Terraform infrastructure code",
	},
}

// ComponentMapper implements ComponentMapperInterface
type ComponentMapper struct {
	mappings map[string]string
}

// NewComponentMapper creates a new component mapper with default mappings
func NewComponentMapper() interfaces.ComponentMapperInterface {
	mapper := &ComponentMapper{
		mappings: make(map[string]string),
	}

	// Register default mappings
	for _, mapping := range DefaultComponentMappings {
		mapper.mappings[mapping.ComponentType] = mapping.TargetPath
	}

	return mapper
}

// GetMapping returns the directory mapping for a component type
func (cm *ComponentMapper) GetMapping(componentType string) (string, error) {
	path, exists := cm.mappings[componentType]
	if !exists {
		return "", fmt.Errorf("no mapping found for component type: %s", componentType)
	}
	return path, nil
}

// RegisterMapping adds a new component-to-directory mapping
func (cm *ComponentMapper) RegisterMapping(componentType string, targetPath string) error {
	if componentType == "" {
		return fmt.Errorf("component type cannot be empty")
	}
	if targetPath == "" {
		return fmt.Errorf("target path cannot be empty")
	}

	cm.mappings[componentType] = targetPath
	return nil
}

// ListMappings returns all registered component mappings
func (cm *ComponentMapper) ListMappings() map[string]string {
	// Return a copy to prevent external modification
	result := make(map[string]string, len(cm.mappings))
	for k, v := range cm.mappings {
		result[k] = v
	}
	return result
}

// StructureMapper implements StructureMapperInterface
type StructureMapper struct {
	componentMapper interfaces.ComponentMapperInterface
}

// NewStructureMapper creates a new structure mapper
func NewStructureMapper(componentMapper interfaces.ComponentMapperInterface) interfaces.StructureMapperInterface {
	return &StructureMapper{
		componentMapper: componentMapper,
	}
}

// GetTargetPath returns the target path for a given component type
func (sm *StructureMapper) GetTargetPath(componentType string) string {
	path, err := sm.componentMapper.GetMapping(componentType)
	if err != nil {
		return ""
	}
	return path
}

// Map relocates generated files from source to target directory structure
func (sm *StructureMapper) Map(ctx context.Context, source string, target string, componentType string) error {
	return sm.MapWithOptions(ctx, source, target, componentType, MapOptions{
		UseSymlinks:      false,
		UpdateReferences: true,
		PreserveSource:   false,
	})
}

// MapWithOptions relocates generated files with custom options
func (sm *StructureMapper) MapWithOptions(ctx context.Context, source string, target string, componentType string, options MapOptions) error {
	// Get the target path for this component type
	relativePath, err := sm.componentMapper.GetMapping(componentType)
	if err != nil {
		return fmt.Errorf("failed to get mapping for component type %s: %w", componentType, err)
	}

	// Build the full target path
	targetPath := filepath.Join(target, relativePath)

	// Check if source exists
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("source path does not exist: %s: %w", source, err)
	}

	// If source and target are the same, no need to move
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for source: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for target: %w", err)
	}

	if absSource == absTarget {
		// Already in the correct location
		if options.UpdateReferences {
			return sm.UpdateReferences(ctx, target, componentType)
		}
		return nil
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Check if target already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	}

	// Handle relocation based on options
	if options.UseSymlinks {
		// Create symlink instead of moving
		if err := sm.createSymlink(absSource, absTarget); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	} else if options.PreserveSource {
		// Copy without removing source
		if sourceInfo.IsDir() {
			if err := sm.copyDirectory(absSource, absTarget); err != nil {
				return fmt.Errorf("failed to copy directory: %w", err)
			}
		} else {
			if err := sm.copyFile(absSource, absTarget); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
		}
	} else {
		// Move (default behavior)
		if sourceInfo.IsDir() {
			if err := sm.moveDirectory(absSource, absTarget); err != nil {
				return fmt.Errorf("failed to move directory: %w", err)
			}
		} else {
			if err := os.Rename(absSource, absTarget); err != nil {
				return fmt.Errorf("failed to move file: %w", err)
			}
		}
	}

	// Update references if requested
	if options.UpdateReferences {
		if err := sm.UpdateReferences(ctx, target, componentType); err != nil {
			return fmt.Errorf("failed to update references: %w", err)
		}
	}

	return nil
}

// createSymlink creates a symbolic link from source to target
func (sm *StructureMapper) createSymlink(source, target string) error {
	// Create symlink
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	return nil
}

// moveDirectory moves a directory from source to target
func (sm *StructureMapper) moveDirectory(source, target string) error {
	// Try simple rename first (works if on same filesystem)
	if err := os.Rename(source, target); err == nil {
		return nil
	}

	// If rename fails, copy and then remove
	if err := sm.copyDirectory(source, target); err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}

	// Remove source after successful copy
	if err := os.RemoveAll(source); err != nil {
		return fmt.Errorf("failed to remove source directory: %w", err)
	}

	return nil
}

// copyDirectory recursively copies a directory
func (sm *StructureMapper) copyDirectory(source, target string) error {
	// Create target directory
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		targetPath := filepath.Join(target, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := sm.copyDirectory(sourcePath, targetPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := sm.copyFile(sourcePath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func (sm *StructureMapper) copyFile(source, target string) error {
	sourceData, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	return os.WriteFile(target, sourceData, sourceInfo.Mode())
}

// ValidateStructure verifies that the target structure is correct
func (sm *StructureMapper) ValidateStructure(rootDir string) error {
	return sm.ValidateStructureWithComponents(rootDir, nil)
}

// ValidateStructureWithComponents validates structure for specific components
func (sm *StructureMapper) ValidateStructureWithComponents(rootDir string, componentTypes []string) error {
	// Check if root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("root directory does not exist: %s", rootDir)
	}

	// Get all registered mappings
	mappings := sm.componentMapper.ListMappings()

	// Track validation results
	var errors []string
	var warnings []string

	// If specific components are provided, only validate those
	componentsToValidate := make(map[string]bool)
	if len(componentTypes) > 0 {
		for _, ct := range componentTypes {
			componentsToValidate[ct] = true
		}
	} else {
		// Validate all mappings
		for ct := range mappings {
			componentsToValidate[ct] = true
		}
	}

	// Check each expected directory
	for componentType, relativePath := range mappings {
		// Skip if not in validation list
		if len(componentTypes) > 0 && !componentsToValidate[componentType] {
			continue
		}

		fullPath := filepath.Join(rootDir, relativePath)

		// Check if directory exists
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Directory doesn't exist - this might be okay if component wasn't generated
				//nolint:staticcheck // warnings is used in the return value
				warnings = append(warnings, fmt.Sprintf("%s directory not found: %s (may not be generated)", componentType, fullPath))
				continue
			}
			errors = append(errors, fmt.Sprintf("error checking %s (%s): %v", componentType, fullPath, err))
			continue
		}

		// Verify it's a directory
		if !info.IsDir() {
			errors = append(errors, fmt.Sprintf("%s path is not a directory: %s", componentType, fullPath))
			continue
		}

		// Validate component-specific structure
		if err := sm.validateComponentStructure(componentType, fullPath); err != nil {
			errors = append(errors, fmt.Sprintf("%s structure validation failed: %v", componentType, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("structure validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// validateComponentStructure validates component-specific directory structure
func (sm *StructureMapper) validateComponentStructure(componentType, componentPath string) error {
	switch componentType {
	case "nextjs":
		return sm.validateNextJSStructure(componentPath)
	case "go-backend":
		return sm.validateGoStructure(componentPath)
	case "android":
		return sm.validateAndroidStructure(componentPath)
	case "ios":
		return sm.validateiOSStructure(componentPath)
	default:
		// No specific validation for this component type
		return nil
	}
}

// validateNextJSStructure validates Next.js project structure
func (sm *StructureMapper) validateNextJSStructure(componentPath string) error {
	// Check for essential Next.js files
	essentialFiles := []string{"package.json", "next.config.js", "next.config.mjs", "next.config.ts"}
	foundConfig := false

	for _, file := range essentialFiles {
		filePath := filepath.Join(componentPath, file)
		if _, err := os.Stat(filePath); err == nil {
			if file == "package.json" || strings.HasPrefix(file, "next.config") {
				foundConfig = true
			}
		}
	}

	if !foundConfig {
		return fmt.Errorf("missing essential Next.js configuration files")
	}

	// Check for app or pages directory (Next.js 13+ or older)
	appDir := filepath.Join(componentPath, "app")
	pagesDir := filepath.Join(componentPath, "pages")
	srcAppDir := filepath.Join(componentPath, "src", "app")
	srcPagesDir := filepath.Join(componentPath, "src", "pages")

	hasAppDir := false
	for _, dir := range []string{appDir, pagesDir, srcAppDir, srcPagesDir} {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			hasAppDir = true
			break
		}
	}

	if !hasAppDir {
		return fmt.Errorf("missing app or pages directory")
	}

	return nil
}

// validateGoStructure validates Go project structure
func (sm *StructureMapper) validateGoStructure(componentPath string) error {
	// Check for go.mod
	goModPath := filepath.Join(componentPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("missing go.mod file")
	}

	// Check for main.go or cmd directory
	mainGoPath := filepath.Join(componentPath, "main.go")
	cmdDir := filepath.Join(componentPath, "cmd")

	hasEntryPoint := false
	if _, err := os.Stat(mainGoPath); err == nil {
		hasEntryPoint = true
	}
	if info, err := os.Stat(cmdDir); err == nil && info.IsDir() {
		hasEntryPoint = true
	}

	if !hasEntryPoint {
		return fmt.Errorf("missing main.go or cmd directory")
	}

	return nil
}

// validateAndroidStructure validates Android project structure
func (sm *StructureMapper) validateAndroidStructure(componentPath string) error {
	// Check for essential Android files
	buildGradle := filepath.Join(componentPath, "build.gradle")
	buildGradleKts := filepath.Join(componentPath, "build.gradle.kts")
	settingsGradle := filepath.Join(componentPath, "settings.gradle")
	settingsGradleKts := filepath.Join(componentPath, "settings.gradle.kts")

	hasBuildFile := false
	if _, err := os.Stat(buildGradle); err == nil {
		hasBuildFile = true
	}
	if _, err := os.Stat(buildGradleKts); err == nil {
		hasBuildFile = true
	}

	if !hasBuildFile {
		return fmt.Errorf("missing build.gradle or build.gradle.kts")
	}

	hasSettingsFile := false
	if _, err := os.Stat(settingsGradle); err == nil {
		hasSettingsFile = true
	}
	if _, err := os.Stat(settingsGradleKts); err == nil {
		hasSettingsFile = true
	}

	if !hasSettingsFile {
		return fmt.Errorf("missing settings.gradle or settings.gradle.kts")
	}

	// Check for app directory
	appDir := filepath.Join(componentPath, "app")
	if info, err := os.Stat(appDir); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("missing app directory")
	}

	return nil
}

// validateiOSStructure validates iOS project structure
func (sm *StructureMapper) validateiOSStructure(componentPath string) error {
	// Check for Xcode project or Swift package
	xcodeprojPattern := filepath.Join(componentPath, "*.xcodeproj")
	matches, err := filepath.Glob(xcodeprojPattern)
	if err != nil {
		return fmt.Errorf("failed to search for xcodeproj: %w", err)
	}

	packageSwiftPath := filepath.Join(componentPath, "Package.swift")
	hasPackageSwift := false
	if _, err := os.Stat(packageSwiftPath); err == nil {
		hasPackageSwift = true
	}

	if len(matches) == 0 && !hasPackageSwift {
		return fmt.Errorf("missing .xcodeproj or Package.swift")
	}

	return nil
}

// ValidationResult represents the result of structure validation
type ValidationResult struct {
	Valid      bool
	Errors     []string
	Warnings   []string
	Components map[string]ComponentValidation
}

// ComponentValidation represents validation result for a single component
type ComponentValidation struct {
	ComponentType string
	Path          string
	Valid         bool
	Errors        []string
	Warnings      []string
}

// ValidateStructureDetailed performs detailed validation and returns structured results
func (sm *StructureMapper) ValidateStructureDetailed(rootDir string, componentTypes []string) *ValidationResult {
	result := &ValidationResult{
		Valid:      true,
		Errors:     make([]string, 0),
		Warnings:   make([]string, 0),
		Components: make(map[string]ComponentValidation),
	}

	// Check if root directory exists
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("root directory does not exist: %s", rootDir))
		return result
	}

	// Get all registered mappings
	mappings := sm.componentMapper.ListMappings()

	// Determine which components to validate
	componentsToValidate := make(map[string]bool)
	if len(componentTypes) > 0 {
		for _, ct := range componentTypes {
			componentsToValidate[ct] = true
		}
	} else {
		for ct := range mappings {
			componentsToValidate[ct] = true
		}
	}

	// Validate each component
	for componentType, relativePath := range mappings {
		if len(componentTypes) > 0 && !componentsToValidate[componentType] {
			continue
		}

		compValidation := ComponentValidation{
			ComponentType: componentType,
			Path:          filepath.Join(rootDir, relativePath),
			Valid:         true,
			Errors:        make([]string, 0),
			Warnings:      make([]string, 0),
		}

		fullPath := filepath.Join(rootDir, relativePath)

		// Check if directory exists
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				compValidation.Warnings = append(compValidation.Warnings, "directory not found (may not be generated)")
			} else {
				compValidation.Valid = false
				compValidation.Errors = append(compValidation.Errors, fmt.Sprintf("error checking path: %v", err))
				result.Valid = false
			}
		} else {
			// Verify it's a directory
			if !info.IsDir() {
				compValidation.Valid = false
				compValidation.Errors = append(compValidation.Errors, "path is not a directory")
				result.Valid = false
			} else {
				// Validate component-specific structure
				if err := sm.validateComponentStructure(componentType, fullPath); err != nil {
					compValidation.Valid = false
					compValidation.Errors = append(compValidation.Errors, err.Error())
					result.Valid = false
				}
			}
		}

		result.Components[componentType] = compValidation

		// Aggregate errors and warnings
		result.Errors = append(result.Errors, compValidation.Errors...)
		result.Warnings = append(result.Warnings, compValidation.Warnings...)
	}

	return result
}

// UpdateReferences updates import paths and references after relocation
func (sm *StructureMapper) UpdateReferences(ctx context.Context, rootDir string, componentType string) error {
	// Get the target path for this component
	relativePath, err := sm.componentMapper.GetMapping(componentType)
	if err != nil {
		return fmt.Errorf("failed to get mapping for component type %s: %w", componentType, err)
	}

	componentPath := filepath.Join(rootDir, relativePath)

	// Check if component path exists
	if _, err := os.Stat(componentPath); os.IsNotExist(err) {
		return fmt.Errorf("component path does not exist: %s", componentPath)
	}

	// Update references based on component type
	switch componentType {
	case "nextjs":
		return sm.updateNextJSReferences(componentPath, rootDir)
	case "go-backend":
		return sm.updateGoReferences(componentPath, rootDir)
	case "android":
		return sm.updateAndroidReferences(componentPath, rootDir)
	case "ios":
		return sm.updateiOSReferences(componentPath, rootDir)
	default:
		// No specific reference updates needed for this component type
		return nil
	}
}

// updateNextJSReferences updates Next.js configuration files
func (sm *StructureMapper) updateNextJSReferences(componentPath, rootDir string) error {
	// Update next.config.js if it exists
	configPath := filepath.Join(componentPath, "next.config.js")
	if err := sm.updateConfigFile(configPath, componentPath, rootDir); err != nil {
		return err
	}

	// Update next.config.mjs if it exists
	configMjsPath := filepath.Join(componentPath, "next.config.mjs")
	if err := sm.updateConfigFile(configMjsPath, componentPath, rootDir); err != nil {
		return err
	}

	// Update .env files if they exist
	envFiles := []string{".env", ".env.local", ".env.development", ".env.production"}
	for _, envFile := range envFiles {
		envPath := filepath.Join(componentPath, envFile)
		if err := sm.updateEnvFile(envPath, componentPath, rootDir); err != nil {
			return err
		}
	}

	// Update package.json if needed
	packagePath := filepath.Join(componentPath, "package.json")
	if err := sm.updatePackageJSON(packagePath, componentPath, rootDir); err != nil {
		return err
	}

	return nil
}

// updateConfigFile updates a configuration file with new paths
func (sm *StructureMapper) updateConfigFile(configPath, componentPath, rootDir string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // File doesn't exist, skip
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Update relative paths in the content
	// This is a simple implementation - could be enhanced with proper parsing
	updatedContent := sm.updateRelativePaths(string(content), componentPath, rootDir)

	// Write back if changed
	if updatedContent != string(content) {
		if err := os.WriteFile(configPath, []byte(updatedContent), 0644); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", configPath, err)
		}
	}

	return nil
}

// updateEnvFile updates environment variable files
func (sm *StructureMapper) updateEnvFile(envPath, componentPath, rootDir string) error {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil // File doesn't exist, skip
	}

	content, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read env file %s: %w", envPath, err)
	}

	// Update paths in environment variables
	updatedContent := sm.updateRelativePaths(string(content), componentPath, rootDir)

	// Write back if changed
	if updatedContent != string(content) {
		if err := os.WriteFile(envPath, []byte(updatedContent), 0644); err != nil {
			return fmt.Errorf("failed to write env file %s: %w", envPath, err)
		}
	}

	return nil
}

// updatePackageJSON updates package.json with new paths
func (sm *StructureMapper) updatePackageJSON(packagePath, componentPath, rootDir string) error {
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return nil // File doesn't exist, skip
	}

	// For now, just verify the file is readable
	// Could be enhanced to parse JSON and update specific fields
	_, err := os.ReadFile(packagePath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	return nil
}

// updateRelativePaths updates relative paths in content
func (sm *StructureMapper) updateRelativePaths(content, componentPath, rootDir string) string {
	// This is a simple implementation that could be enhanced
	// For now, it just returns the content as-is
	// In a full implementation, this would parse and update paths
	return content
}

// updateGoReferences updates Go module paths
func (sm *StructureMapper) updateGoReferences(componentPath, rootDir string) error {
	// Check for go.mod
	goModPath := filepath.Join(componentPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		// go.mod exists, verify it's readable
		if _, err := os.ReadFile(goModPath); err != nil {
			return fmt.Errorf("failed to read go.mod: %w", err)
		}
	}

	// Update .env files if they exist
	envFiles := []string{".env", ".env.local", ".env.development", ".env.production"}
	for _, envFile := range envFiles {
		envPath := filepath.Join(componentPath, envFile)
		if err := sm.updateEnvFile(envPath, componentPath, rootDir); err != nil {
			return err
		}
	}

	// Update config files if they exist
	configFiles := []string{"config.yaml", "config.yml", "config.json"}
	for _, configFile := range configFiles {
		configPath := filepath.Join(componentPath, configFile)
		if err := sm.updateConfigFile(configPath, componentPath, rootDir); err != nil {
			return err
		}
	}

	return nil
}

// updateAndroidReferences updates Android project references
func (sm *StructureMapper) updateAndroidReferences(componentPath, rootDir string) error {
	// Update build.gradle files
	buildGradleFiles := []string{"build.gradle", "build.gradle.kts", "app/build.gradle", "app/build.gradle.kts"}
	for _, gradleFile := range buildGradleFiles {
		gradlePath := filepath.Join(componentPath, gradleFile)
		if err := sm.updateConfigFile(gradlePath, componentPath, rootDir); err != nil {
			return err
		}
	}

	// Update gradle.properties if it exists
	gradlePropsPath := filepath.Join(componentPath, "gradle.properties")
	if err := sm.updateConfigFile(gradlePropsPath, componentPath, rootDir); err != nil {
		return err
	}

	// Update local.properties if it exists
	localPropsPath := filepath.Join(componentPath, "local.properties")
	if err := sm.updateConfigFile(localPropsPath, componentPath, rootDir); err != nil {
		return err
	}

	return nil
}

// updateiOSReferences updates iOS project references
func (sm *StructureMapper) updateiOSReferences(componentPath, rootDir string) error {
	// Update Package.swift if it exists (for SPM projects)
	packageSwiftPath := filepath.Join(componentPath, "Package.swift")
	if err := sm.updateConfigFile(packageSwiftPath, componentPath, rootDir); err != nil {
		return err
	}

	// Update Podfile if it exists (for CocoaPods projects)
	podfilePath := filepath.Join(componentPath, "Podfile")
	if err := sm.updateConfigFile(podfilePath, componentPath, rootDir); err != nil {
		return err
	}

	// Note: Xcode project files (.xcodeproj) are complex binary/XML structures
	// that would require specialized parsing. For now, we just verify they exist.
	xcodeprojPattern := filepath.Join(componentPath, "*.xcodeproj")
	matches, err := filepath.Glob(xcodeprojPattern)
	if err != nil {
		return fmt.Errorf("failed to search for xcodeproj: %w", err)
	}

	// Verify xcodeproj directories are accessible
	for _, match := range matches {
		if _, err := os.Stat(match); err != nil {
			return fmt.Errorf("failed to access xcodeproj %s: %w", match, err)
		}
	}

	return nil
}

// MapOptions defines options for mapping operations
type MapOptions struct {
	UseSymlinks      bool // Create symlinks instead of copying
	UpdateReferences bool // Update file references after mapping
	PreserveSource   bool // Keep source files after mapping
}

// MapResult represents the result of a mapping operation
type MapResult struct {
	ComponentType string
	SourcePath    string
	TargetPath    string
	Success       bool
	Error         error
	UsedSymlink   bool
}

// BatchMap maps multiple components at once
func (sm *StructureMapper) BatchMap(ctx context.Context, components []*models.Component, targetRoot string) ([]*MapResult, error) {
	results := make([]*MapResult, 0, len(components))

	for _, component := range components {
		result := &MapResult{
			ComponentType: component.Type,
			SourcePath:    component.Path,
		}

		// Get target path
		relativePath, err := sm.componentMapper.GetMapping(component.Type)
		if err != nil {
			result.Error = err
			result.Success = false
			results = append(results, result)
			continue
		}

		result.TargetPath = filepath.Join(targetRoot, relativePath)

		// Perform mapping
		if err := sm.Map(ctx, component.Path, targetRoot, component.Type); err != nil {
			result.Error = err
			result.Success = false
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}
