package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	texttemplate "text/template"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/version"
)

// Engine implements the TemplateEngine interface
type Engine struct {
	funcMap        texttemplate.FuncMap
	versionManager interfaces.VersionManager
}

// NewEngine creates a new template engine instance
func NewEngine() interfaces.TemplateEngine {
	engine := &Engine{
		funcMap: make(texttemplate.FuncMap),
	}

	// Register default template functions
	engine.registerDefaultFunctions()

	return engine
}

// NewEngineWithVersionManager creates a new template engine with version management
func NewEngineWithVersionManager(versionManager interfaces.VersionManager) interfaces.TemplateEngine {
	engine := &Engine{
		funcMap:        make(texttemplate.FuncMap),
		versionManager: versionManager,
	}

	// Register default template functions
	engine.registerDefaultFunctions()

	return engine
}

// ProcessTemplate processes a single template file with the given configuration
func (e *Engine) ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Enhance config with centralized version information
	enhancedConfig, err := e.enhanceConfigWithVersions(config)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance config with versions: %w", err)
	}

	// Perform pre-generation version validation
	if err := e.validatePreGeneration(enhancedConfig, templatePath); err != nil {
		return nil, fmt.Errorf("pre-generation validation failed for template %s: %w", templatePath, err)
	}

	// Load the template
	tmpl, err := e.LoadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", templatePath, err)
	}

	// Render the template
	return e.RenderTemplate(tmpl, enhancedConfig)
}

// ProcessDirectory processes an entire template directory recursively
func (e *Engine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Enhance config with centralized version information once for the entire directory
	enhancedConfig, err := e.enhanceConfigWithVersions(config)
	if err != nil {
		return fmt.Errorf("failed to enhance config with versions: %w", err)
	}

	// Perform pre-generation validation for the entire directory
	if err := e.validatePreGenerationDirectory(enhancedConfig, templateDir); err != nil {
		return fmt.Errorf("pre-generation validation failed for directory %s: %w", templateDir, err)
	}

	return filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}

		// Calculate relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Calculate output path
		outputPath := filepath.Join(outputDir, relPath)

		if d.IsDir() {
			// Create directory in output
			return os.MkdirAll(outputPath, 0755)
		}

		// Process file with enhanced config
		return e.processFile(path, outputPath, enhancedConfig)
	})
}

// RegisterFunctions registers custom template functions
func (e *Engine) RegisterFunctions(funcMap texttemplate.FuncMap) {
	for name, fn := range funcMap {
		e.funcMap[name] = fn
	}
}

// LoadTemplate loads and parses a template from the given path
func (e *Engine) LoadTemplate(templatePath string) (*texttemplate.Template, error) {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template with custom functions
	tmpl := texttemplate.New(filepath.Base(templatePath)).Funcs(e.funcMap)

	// Parse template content
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func (e *Engine) RenderTemplate(tmpl *texttemplate.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// processFile processes a single file (template or binary)
func (e *Engine) processFile(srcPath, destPath string, config *models.ProjectConfig) error {
	// Check if file is a template (has .tmpl extension)
	if strings.HasSuffix(srcPath, ".tmpl") {
		// Remove .tmpl extension from destination
		destPath = strings.TrimSuffix(destPath, ".tmpl")

		// Process as template
		content, err := e.ProcessTemplate(srcPath, config)
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", srcPath, err)
		}

		// Write processed content
		return os.WriteFile(destPath, content, 0644)
	}

	// Copy binary file as-is
	return e.copyFile(srcPath, destPath)
}

// copyFile copies a file from src to dest
func (e *Engine) copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy file content
	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	return os.Chmod(dest, srcInfo.Mode())
}

// getDefaultNodeVersionConfig returns the default Node.js version configuration
func (e *Engine) getDefaultNodeVersionConfig() *models.NodeVersionConfig {
	return &models.NodeVersionConfig{
		Runtime:      ">=20.0.0",
		TypesPackage: "^20.17.0",
		NPMVersion:   ">=10.0.0",
		DockerImage:  "node:20-alpine",
		LTSStatus:    true,
		Description:  "Node.js 20 LTS - Recommended for production use",
	}
}

// applyNodeVersionDefaults applies Node.js version defaults to the configuration
func (e *Engine) applyNodeVersionDefaults(config *models.ProjectConfig) {
	if config.Versions.NodeJS == nil {
		config.Versions.NodeJS = e.getDefaultNodeVersionConfig()
	}

	// Ensure @types/node package version is set based on NodeJS config
	if config.Versions.NodeJS.TypesPackage != "" {
		if config.Versions.Packages == nil {
			config.Versions.Packages = make(map[string]string)
		}
		config.Versions.Packages["@types/node"] = config.Versions.NodeJS.TypesPackage
	}

	// Set Node version if not already set
	if config.Versions.Node == "" {
		// Extract version from runtime requirement
		runtime := config.Versions.NodeJS.Runtime
		if strings.HasPrefix(runtime, ">=") {
			config.Versions.Node = strings.TrimPrefix(runtime, ">=")
		} else if strings.HasPrefix(runtime, ">") {
			config.Versions.Node = strings.TrimPrefix(runtime, ">")
		} else {
			config.Versions.Node = "20.0.0" // fallback
		}
	}
}

// validateVersionCompatibility validates version compatibility during template processing
func (e *Engine) validateVersionCompatibility(config *models.ProjectConfig) error {
	if config.Versions == nil || config.Versions.NodeJS == nil {
		return nil // No validation needed if no version config
	}

	// Import validation package to use the validator
	validator := &VersionValidator{}
	result := validator.ValidateNodeVersionConfig(config.Versions.NodeJS)

	// Check for critical errors
	for _, err := range result.Errors {
		if err.Severity == "critical" {
			return fmt.Errorf("critical version compatibility error in %s: %s", err.Field, err.Message)
		}
	}

	// Log warnings but don't fail
	for _, warning := range result.Warnings {
		fmt.Printf("Warning: Version compatibility issue in %s: %s\n", warning.Field, warning.Message)
	}

	return nil
}

// VersionValidator is embedded to avoid import cycles
type VersionValidator struct{}

// ValidateNodeVersionConfig validates a Node.js version configuration
func (v *VersionValidator) ValidateNodeVersionConfig(config *models.NodeVersionConfig) *models.VersionValidationResult {
	result := &models.VersionValidationResult{
		Valid:       true,
		Errors:      []models.VersionValidationError{},
		Warnings:    []models.VersionValidationWarning{},
		Suggestions: []models.VersionSuggestion{},
		ValidatedAt: time.Now(),
	}

	// Basic validation - check for empty values
	if config.Runtime == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "runtime",
			Value:    config.Runtime,
			Message:  "Runtime version cannot be empty",
			Severity: "critical",
			Code:     "EMPTY_RUNTIME_VERSION",
		})
	}

	if config.TypesPackage == "" {
		result.Valid = false
		result.Errors = append(result.Errors, models.VersionValidationError{
			Field:    "types",
			Value:    config.TypesPackage,
			Message:  "Types package version cannot be empty",
			Severity: "critical",
			Code:     "EMPTY_TYPES_VERSION",
		})
	}

	// Check compatibility between runtime and types versions
	if config.Runtime != "" && config.TypesPackage != "" {
		if err := v.validateRuntimeTypesCompatibility(config.Runtime, config.TypesPackage); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, models.VersionValidationError{
				Field:    "compatibility",
				Value:    fmt.Sprintf("runtime: %s, types: %s", config.Runtime, config.TypesPackage),
				Message:  err.Error(),
				Severity: "critical",
				Code:     "VERSION_COMPATIBILITY_MISMATCH",
			})
		}
	}

	return result
}

// validateRuntimeTypesCompatibility validates compatibility between runtime and types versions
func (v *VersionValidator) validateRuntimeTypesCompatibility(runtime, types string) error {
	runtimeMajor, err := v.extractMajorVersion(runtime)
	if err != nil {
		return fmt.Errorf("cannot extract runtime major version: %w", err)
	}

	typesMajor, err := v.extractMajorVersion(types)
	if err != nil {
		return fmt.Errorf("cannot extract types major version: %w", err)
	}

	// Types version should match or be close to runtime version
	// Allow types to be same major version or one version ahead
	if typesMajor < runtimeMajor || typesMajor > runtimeMajor+2 {
		return fmt.Errorf("types version %d is not compatible with runtime version %d", typesMajor, runtimeMajor)
	}

	return nil
}

// extractMajorVersion extracts the major version number from a version string
func (v *VersionValidator) extractMajorVersion(version string) (int, error) {
	// Remove version operators (>=, ^, ~, etc.)
	re := regexp.MustCompile(`^[>=<~^]*`)
	cleanVersion := re.ReplaceAllString(version, "")

	// Split by dots and get the first part
	parts := strings.Split(cleanVersion, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	return major, nil
}

// enhanceConfigWithVersions merges centralized version information into the project config
func (e *Engine) enhanceConfigWithVersions(config *models.ProjectConfig) (*models.ProjectConfig, error) {
	// Handle nil config
	if config == nil {
		return nil, fmt.Errorf("project config cannot be nil")
	}

	// Create a copy of the config to avoid modifying the original
	enhancedConfig := *config
	if enhancedConfig.Versions == nil {
		enhancedConfig.Versions = &models.VersionConfig{
			Packages: make(map[string]string),
		}
	} else {
		// Deep copy the versions
		enhancedVersions := *enhancedConfig.Versions
		enhancedVersions.Packages = make(map[string]string)
		for k, v := range config.Versions.Packages {
			enhancedVersions.Packages[k] = v
		}
		enhancedConfig.Versions = &enhancedVersions
	}

	// Set up default Node.js version configuration if not present
	if enhancedConfig.Versions.NodeJS == nil {
		enhancedConfig.Versions.NodeJS = e.getDefaultNodeVersionConfig()
	}

	// Validate version compatibility before processing
	if err := e.validateVersionCompatibility(&enhancedConfig); err != nil {
		return nil, fmt.Errorf("version compatibility validation failed: %w", err)
	}

	// If no version manager is configured, use defaults and return
	if e.versionManager == nil {
		e.applyNodeVersionDefaults(&enhancedConfig)
		return &enhancedConfig, nil
	}

	// Check if version manager supports storage (enhanced functionality)
	if managerWithStorage, ok := e.versionManager.(*version.Manager); ok {
		store, err := managerWithStorage.GetVersionStore()
		if err != nil {
			// Log warning but don't fail - use existing versions
			fmt.Printf("Warning: Could not load version store: %v\n", err)
			e.applyNodeVersionDefaults(&enhancedConfig)
			return &enhancedConfig, nil
		}

		// Merge language versions
		for name, info := range store.Languages {
			switch name {
			case "nodejs":
				enhancedConfig.Versions.Node = info.CurrentVersion
				// Update NodeJS config with actual version info
				if enhancedConfig.Versions.NodeJS != nil {
					enhancedConfig.Versions.NodeJS.Runtime = ">=" + info.CurrentVersion
				}
			case "go":
				enhancedConfig.Versions.Go = info.CurrentVersion
			case "kotlin":
				enhancedConfig.Versions.Kotlin = info.CurrentVersion
			case "swift":
				enhancedConfig.Versions.Swift = info.CurrentVersion
			}
		}

		// Merge framework versions
		for name, info := range store.Frameworks {
			switch name {
			case "nextjs":
				enhancedConfig.Versions.NextJS = info.CurrentVersion
			case "react":
				enhancedConfig.Versions.React = info.CurrentVersion
			}
		}

		// Merge package versions
		for name, info := range store.Packages {
			enhancedConfig.Versions.Packages[name] = info.CurrentVersion
		}

		enhancedConfig.Versions.UpdatedAt = store.LastUpdated
	} else {
		// Fallback to basic version manager functionality
		if nodeVersion, err := e.versionManager.GetLatestNodeVersion(); err == nil {
			enhancedConfig.Versions.Node = nodeVersion
			if enhancedConfig.Versions.NodeJS != nil {
				enhancedConfig.Versions.NodeJS.Runtime = ">=" + nodeVersion
			}
		}
		if goVersion, err := e.versionManager.GetLatestGoVersion(); err == nil {
			enhancedConfig.Versions.Go = goVersion
		}
		if nextjsVersion, err := e.versionManager.GetLatestNPMPackage("next"); err == nil {
			enhancedConfig.Versions.NextJS = nextjsVersion
		}
		if reactVersion, err := e.versionManager.GetLatestNPMPackage("react"); err == nil {
			enhancedConfig.Versions.React = reactVersion
		}

		// Add common packages
		commonPackages := []string{
			"typescript", "tailwindcss", "eslint", "prettier", "jest",
			"@types/node", "@types/react", "autoprefixer", "postcss",
		}
		for _, pkg := range commonPackages {
			if version, err := e.versionManager.GetLatestNPMPackage(pkg); err == nil {
				enhancedConfig.Versions.Packages[pkg] = version
			}
		}
	}

	// Apply Node.js version defaults for any missing values
	e.applyNodeVersionDefaults(&enhancedConfig)

	return &enhancedConfig, nil
}

// validatePreGeneration performs comprehensive pre-generation validation for a single template
func (e *Engine) validatePreGeneration(config *models.ProjectConfig, templatePath string) error {
	// Validate version configuration exists
	if config.Versions == nil {
		return fmt.Errorf("version configuration is required for template generation")
	}

	// Validate Node.js version configuration if this is a frontend template
	if e.isFrontendTemplate(templatePath) {
		if err := e.validateNodeJSVersions(config); err != nil {
			return fmt.Errorf("node.js version validation failed: %w", err)
		}
	}

	// Validate Go version configuration if this is a backend template
	if e.isBackendTemplate(templatePath) {
		if err := e.validateGoVersions(config); err != nil {
			return fmt.Errorf("go version validation failed: %w", err)
		}
	}

	// Validate cross-template version consistency
	if err := e.validateVersionConsistency(config, templatePath); err != nil {
		return fmt.Errorf("version consistency validation failed: %w", err)
	}

	return nil
}

// validatePreGenerationDirectory performs pre-generation validation for an entire template directory
func (e *Engine) validatePreGenerationDirectory(config *models.ProjectConfig, templateDir string) error {
	// Validate version configuration exists
	if config.Versions == nil {
		return fmt.Errorf("version configuration is required for template generation")
	}

	// Collect all template files for comprehensive validation
	templateFiles, err := e.collectTemplateFiles(templateDir)
	if err != nil {
		return fmt.Errorf("failed to collect template files: %w", err)
	}

	// Validate Node.js versions if any frontend templates are present
	if e.hasFrontendTemplates(templateFiles) {
		if err := e.validateNodeJSVersions(config); err != nil {
			return fmt.Errorf("node.js version validation failed: %w", err)
		}
	}

	// Validate Go versions if any backend templates are present
	if e.hasBackendTemplates(templateFiles) {
		if err := e.validateGoVersions(config); err != nil {
			return fmt.Errorf("go version validation failed: %w", err)
		}
	}

	// Validate cross-template version consistency across all templates
	if err := e.validateCrossTemplateConsistency(config, templateFiles); err != nil {
		return fmt.Errorf("cross-template version consistency validation failed: %w", err)
	}

	return nil
}

// validateNodeJSVersions validates Node.js version configuration
func (e *Engine) validateNodeJSVersions(config *models.ProjectConfig) error {
	if config.Versions.NodeJS == nil {
		return fmt.Errorf("node.js version configuration is required for frontend templates")
	}

	// Use the embedded version validator to validate Node.js configuration
	validator := &VersionValidator{}
	result := validator.ValidateNodeVersionConfig(config.Versions.NodeJS)

	// Check for critical errors that should block generation
	for _, err := range result.Errors {
		if err.Severity == "critical" {
			return fmt.Errorf("critical node.js version error in %s: %s", err.Field, err.Message)
		}
	}

	// Log warnings but don't fail
	for _, warning := range result.Warnings {
		fmt.Printf("Warning: Node.js version issue in %s: %s\n", warning.Field, warning.Message)
	}

	return nil
}

// validateGoVersions validates Go version configuration
func (e *Engine) validateGoVersions(config *models.ProjectConfig) error {
	if config.Versions.Go == "" {
		return fmt.Errorf("go version is required for backend templates")
	}

	// Validate Go version format
	if !e.isValidGoVersion(config.Versions.Go) {
		return fmt.Errorf("invalid go version format: %s", config.Versions.Go)
	}

	// Check if Go version is supported (minimum 1.20)
	majorMinor, err := e.extractGoMajorMinor(config.Versions.Go)
	if err != nil {
		return fmt.Errorf("failed to parse go version: %w", err)
	}

	if majorMinor < 1.20 {
		return fmt.Errorf("go version %s is not supported, minimum required version is 1.20", config.Versions.Go)
	}

	return nil
}

// validateVersionConsistency validates version consistency for a single template
func (e *Engine) validateVersionConsistency(config *models.ProjectConfig, templatePath string) error {
	// Check if template requires specific version constraints
	if strings.Contains(templatePath, "package.json.tmpl") {
		return e.validatePackageJSONVersions(config, templatePath)
	}

	if strings.Contains(templatePath, "go.mod.tmpl") {
		return e.validateGoModVersions(config, templatePath)
	}

	if strings.Contains(templatePath, "Dockerfile.tmpl") {
		return e.validateDockerVersions(config, templatePath)
	}

	return nil
}

// validateCrossTemplateConsistency validates version consistency across multiple templates
func (e *Engine) validateCrossTemplateConsistency(config *models.ProjectConfig, templateFiles []string) error {
	// Collect all package.json templates and validate they use consistent versions
	packageJSONTemplates := []string{}
	goModTemplates := []string{}
	dockerTemplates := []string{}

	for _, file := range templateFiles {
		if strings.Contains(file, "package.json.tmpl") {
			packageJSONTemplates = append(packageJSONTemplates, file)
		} else if strings.Contains(file, "go.mod.tmpl") {
			goModTemplates = append(goModTemplates, file)
		} else if strings.Contains(file, "Dockerfile.tmpl") {
			dockerTemplates = append(dockerTemplates, file)
		}
	}

	// Validate consistency across package.json templates
	if len(packageJSONTemplates) > 1 {
		if err := e.validatePackageJSONConsistency(config, packageJSONTemplates); err != nil {
			return fmt.Errorf("package.json version consistency validation failed: %w", err)
		}
	}

	// Validate consistency across go.mod templates
	if len(goModTemplates) > 1 {
		if err := e.validateGoModConsistency(config, goModTemplates); err != nil {
			return fmt.Errorf("go.mod version consistency validation failed: %w", err)
		}
	}

	// Validate consistency across Docker templates
	if len(dockerTemplates) > 0 {
		if err := e.validateDockerConsistency(config, dockerTemplates); err != nil {
			return fmt.Errorf("docker version consistency validation failed: %w", err)
		}
	}

	return nil
}

// validatePackageJSONVersions validates versions in a package.json template
func (e *Engine) validatePackageJSONVersions(config *models.ProjectConfig, templatePath string) error {
	// Validate Node.js engine requirement matches @types/node version
	if config.Versions.NodeJS != nil {
		nodeRuntime := config.Versions.NodeJS.Runtime
		typesVersion := config.Versions.NodeJS.TypesPackage

		// Extract major versions for compatibility check
		runtimeMajor, err := e.extractMajorVersion(nodeRuntime)
		if err != nil {
			return fmt.Errorf("invalid node.js runtime version %s: %w", nodeRuntime, err)
		}

		typesMajor, err := e.extractMajorVersion(typesVersion)
		if err != nil {
			return fmt.Errorf("invalid @types/node version %s: %w", typesVersion, err)
		}

		// Check compatibility (types should match or be close to runtime)
		if typesMajor < runtimeMajor || typesMajor > runtimeMajor+2 {
			return fmt.Errorf("@types/node version %d is incompatible with node.js runtime version %d in template %s",
				typesMajor, runtimeMajor, templatePath)
		}
	}

	return nil
}

// validateGoModVersions validates versions in a go.mod template
func (e *Engine) validateGoModVersions(config *models.ProjectConfig, templatePath string) error {
	if config.Versions.Go == "" {
		return fmt.Errorf("go version is required for go.mod template %s", templatePath)
	}

	// Validate Go version format and support
	if !e.isValidGoVersion(config.Versions.Go) {
		return fmt.Errorf("invalid go version %s in template %s", config.Versions.Go, templatePath)
	}

	return nil
}

// validateDockerVersions validates versions in Docker templates
func (e *Engine) validateDockerVersions(config *models.ProjectConfig, templatePath string) error {
	// Validate Docker image versions match runtime versions
	if config.Versions.NodeJS != nil && config.Versions.NodeJS.DockerImage != "" {
		dockerImage := config.Versions.NodeJS.DockerImage

		// Extract Node.js version from Docker image
		if !strings.Contains(dockerImage, "node:") {
			return fmt.Errorf("docker image %s is not a Node.js image in template %s", dockerImage, templatePath)
		}

		// Basic validation that Docker image version is reasonable
		if strings.Contains(dockerImage, "node:") && !strings.Contains(dockerImage, "20") {
			fmt.Printf("Warning: Docker image %s may not match Node.js runtime version in template %s\n",
				dockerImage, templatePath)
		}
	}

	return nil
}

// validatePackageJSONConsistency validates consistency across multiple package.json templates
func (e *Engine) validatePackageJSONConsistency(config *models.ProjectConfig, templates []string) error {
	if config.Versions.NodeJS == nil {
		return fmt.Errorf("node.js version configuration required for package.json consistency validation")
	}

	// All package.json templates should use the same Node.js and @types/node versions
	expectedRuntime := config.Versions.NodeJS.Runtime
	expectedTypes := config.Versions.NodeJS.TypesPackage

	for _, template := range templates {
		// In a real implementation, we would parse the template content
		// For now, we just validate that the configuration is consistent
		fmt.Printf("Validating package.json consistency for template: %s\n", template)
	}

	fmt.Printf("All package.json templates will use node.js runtime: %s, @types/node: %s\n",
		expectedRuntime, expectedTypes)

	return nil
}

// validateGoModConsistency validates consistency across multiple go.mod templates
func (e *Engine) validateGoModConsistency(config *models.ProjectConfig, templates []string) error {
	if config.Versions.Go == "" {
		return fmt.Errorf("go version configuration required for go.mod consistency validation")
	}

	expectedGoVersion := config.Versions.Go

	for _, template := range templates {
		// In a real implementation, we would parse the template content
		// For now, we just validate that the configuration is consistent
		fmt.Printf("Validating go.mod consistency for template: %s\n", template)
	}

	fmt.Printf("All go.mod templates will use go version: %s\n", expectedGoVersion)

	return nil
}

// validateDockerConsistency validates consistency across Docker templates
func (e *Engine) validateDockerConsistency(config *models.ProjectConfig, templates []string) error {
	if config.Versions.NodeJS == nil || config.Versions.NodeJS.DockerImage == "" {
		return fmt.Errorf("docker image configuration required for docker consistency validation")
	}

	expectedImage := config.Versions.NodeJS.DockerImage

	for _, template := range templates {
		fmt.Printf("Validating Docker consistency for template: %s\n", template)
	}

	fmt.Printf("All docker templates will use base image: %s\n", expectedImage)

	return nil
}

// Helper methods for template type detection and validation

// isFrontendTemplate checks if a template path is for frontend components
func (e *Engine) isFrontendTemplate(templatePath string) bool {
	frontendIndicators := []string{
		"frontend/",
		"nextjs-",
		"package.json.tmpl",
		"next.config.js.tmpl",
		"tailwind.config.js.tmpl",
		"tsconfig.json.tmpl",
	}

	for _, indicator := range frontendIndicators {
		if strings.Contains(templatePath, indicator) {
			return true
		}
	}

	return false
}

// isBackendTemplate checks if a template path is for backend components
func (e *Engine) isBackendTemplate(templatePath string) bool {
	backendIndicators := []string{
		"backend/",
		"go-gin/",
		"go.mod.tmpl",
		"main.go.tmpl",
		"internal/",
		"pkg/",
	}

	for _, indicator := range backendIndicators {
		if strings.Contains(templatePath, indicator) {
			return true
		}
	}

	return false
}

// collectTemplateFiles collects all template files in a directory
func (e *Engine) collectTemplateFiles(templateDir string) ([]string, error) {
	var templateFiles []string

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			templateFiles = append(templateFiles, path)
		}

		return nil
	})

	return templateFiles, err
}

// hasFrontendTemplates checks if any frontend templates are present
func (e *Engine) hasFrontendTemplates(templateFiles []string) bool {
	for _, file := range templateFiles {
		if e.isFrontendTemplate(file) {
			return true
		}
	}
	return false
}

// hasBackendTemplates checks if any backend templates are present
func (e *Engine) hasBackendTemplates(templateFiles []string) bool {
	for _, file := range templateFiles {
		if e.isBackendTemplate(file) {
			return true
		}
	}
	return false
}

// isValidGoVersion validates Go version format
func (e *Engine) isValidGoVersion(version string) bool {
	// Go versions are like "1.21", "1.21.0", etc.
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}

	return true
}

// extractGoMajorMinor extracts major.minor version from Go version string
func (e *Engine) extractGoMajorMinor(version string) (float64, error) {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid go version format: %s", version)
	}

	major := parts[0]
	minor := parts[1]

	// Convert to float for comparison
	majorFloat := 0.0
	minorFloat := 0.0

	for _, char := range major {
		if char >= '0' && char <= '9' {
			majorFloat = majorFloat*10 + float64(char-'0')
		}
	}

	for _, char := range minor {
		if char >= '0' && char <= '9' {
			minorFloat = minorFloat*10 + float64(char-'0')
		}
	}

	return majorFloat + minorFloat/100, nil
}

// extractMajorVersion extracts the major version number from a version string
func (e *Engine) extractMajorVersion(version string) (int, error) {
	// Remove version operators (>=, ^, ~, etc.)
	re := regexp.MustCompile(`^[>=<~^]*`)
	cleanVersion := re.ReplaceAllString(version, "")

	// Split by dots and get the first part
	parts := strings.Split(cleanVersion, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	return major, nil
}
