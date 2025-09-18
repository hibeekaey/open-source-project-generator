# API Documentation

This document provides comprehensive API documentation for the Open Source Project Generator.

## Table of Contents

- [Overview](#overview)
- [Core Interfaces](#core-interfaces)
- [Package Structure](#package-structure)
- [CLI Interface](#cli-interface)
- [Template Engine](#template-engine)
- [Version Management](#version-management)
- [Configuration Management](#configuration-management)
- [File System Operations](#file-system-operations)
- [Usage Examples](#usage-examples)
- [Error Handling](#error-handling)
- [Performance Considerations](#performance-considerations)
- [Security Considerations](#security-considerations)

## Overview

The Open Source Project Generator provides a comprehensive API for generating production-ready project structures. The API is organized into several core interfaces that handle different aspects of project generation:

- **CLI Interface**: Interactive command-line operations
- **Template Engine**: Template processing and rendering
- **Config Manager**: Configuration loading and validation
- **File System Generator**: File and directory operations
- **Version Manager**: Package version management and updates

## Core Interfaces

### CLI Interface

The `CLIInterface` provides comprehensive command-line functionality for all generator operations including interactive configuration, validation, auditing, and automation support.

```go
type CLIInterface interface {
    // Core operations
    Run(args []string) error

    // Interactive operations
    PromptProjectDetails() (*models.ProjectConfig, error)
    SelectComponents() ([]string, error)
    ConfirmGeneration(*models.ProjectConfig) bool

    // Advanced interactive operations
    PromptAdvancedOptions() (*AdvancedOptions, error)
    ConfirmAdvancedGeneration(*models.ProjectConfig, *AdvancedOptions) bool
    SelectTemplateInteractively(filter TemplateFilter) (*TemplateInfo, error)

    // Non-interactive operations
    GenerateFromConfig(path string, options GenerateOptions) error
    ValidateProject(path string, options ValidationOptions) (*ValidationResult, error)
    AuditProject(path string, options AuditOptions) (*AuditResult, error)

    // Advanced non-interactive operations
    GenerateWithAdvancedOptions(config *models.ProjectConfig, options *AdvancedOptions) error
    ValidateProjectAdvanced(path string, options *ValidationOptions) (*ValidationResult, error)
    AuditProjectAdvanced(path string, options *AuditOptions) (*AuditResult, error)

    // Template operations
    ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
    GetTemplateInfo(name string) (*TemplateInfo, error)
    ValidateTemplate(path string) (*TemplateValidationResult, error)

    // Template management operations
    SearchTemplates(query string) ([]TemplateInfo, error)
    GetTemplateMetadata(name string) (*TemplateMetadata, error)
    GetTemplateDependencies(name string) ([]string, error)
    ValidateCustomTemplate(path string) (*TemplateValidationResult, error)

    // Configuration operations
    ShowConfig() error
    SetConfig(key, value string) error
    EditConfig() error
    ValidateConfig() error
    ExportConfig(path string) error

    // Configuration management operations
    LoadConfiguration(sources []string) (*models.ProjectConfig, error)
    MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error)
    ValidateConfigurationSchema(config *models.ProjectConfig) error
    GetConfigurationSources() ([]ConfigSource, error)

    // Version and update operations
    ShowVersion(options VersionOptions) error
    CheckUpdates() (*UpdateInfo, error)
    InstallUpdates() error

    // Advanced version operations
    GetPackageVersions() (map[string]string, error)
    GetLatestPackageVersions() (map[string]string, error)
    CheckCompatibility(path string) (*CompatibilityResult, error)

    // Cache operations
    ShowCache() error
    ClearCache() error
    CleanCache() error

    // Cache management operations
    GetCacheStats() (*CacheStats, error)
    ValidateCache() error
    RepairCache() error
    EnableOfflineMode() error
    DisableOfflineMode() error

    // Utility operations
    ShowLogs() error

    // Logging and debugging operations
    SetLogLevel(level string) error
    GetLogLevel() string
    ShowRecentLogs(lines int, level string) error
    GetLogFileLocations() ([]string, error)

    // Automation and integration operations
    RunNonInteractive(config *models.ProjectConfig, options *AdvancedOptions) error
    GenerateReport(reportType string, format string, outputFile string) error
    GetExitCode() int
    SetExitCode(code int)
}
```

**Key Features:**

- Complete command-line interface with all documented commands
- Interactive and non-interactive modes for different use cases
- Comprehensive validation and auditing capabilities
- Advanced configuration management with multiple sources
- Template discovery and management
- Update management with safety checks
- Cache management for offline operation
- Logging and debugging support

**Usage:**

```go
// Create CLI with all dependencies
cli := cli.NewCLI(
    configManager,
    validator,
    templateManager,
    cacheManager,
    versionManager,
    auditEngine,
    logger,
    version,
)

// Run with command-line arguments
err := cli.Run(os.Args[1:])
if err != nil {
    log.Fatal(err)
}

// Or use specific operations
config, err := cli.PromptProjectDetails()
if err != nil {
    return fmt.Errorf("failed to collect project details: %w", err)
}

// Validate project
result, err := cli.ValidateProject("./my-project", ValidationOptions{
    Fix:          true,
    Report:       true,
    ReportFormat: "html",
})
```

### TemplateEngine Interface

The `TemplateEngine` handles template processing, rendering, and directory operations.

```go
type TemplateEngine interface {
    ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error)
    ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error
    RegisterFunctions(funcMap template.FuncMap)
    LoadTemplate(templatePath string) (*template.Template, error)
    RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error)
}
```

**Key Features:**

- Single template file processing with variable substitution
- Recursive directory processing for complete project generation
- Custom function registration for extended template functionality
- Template caching and performance optimization
- Version management integration for automatic dependency updates

**Usage:**

```go
engine := template.NewEngine()

// Process a single template
content, err := engine.ProcessTemplate("templates/package.json.tmpl", config)
if err != nil {
    return fmt.Errorf("failed to process template: %w", err)
}

// Process entire directory
err = engine.ProcessDirectory("templates/frontend", "output/frontend", config)
if err != nil {
    return fmt.Errorf("failed to process directory: %w", err)
}
```

### Version Manager

The `VersionManager` provides comprehensive version management including update capabilities, compatibility checking, and caching.

```go
type VersionManager interface {
    // Basic version operations
    GetLatestNodeVersion() (string, error)
    GetLatestGoVersion() (string, error)
    GetLatestNPMPackage(packageName string) (string, error)
    GetLatestGoModule(moduleName string) (string, error)
    UpdateVersionsConfig() (*models.VersionConfig, error)
    GetLatestGitHubRelease(owner, repo string) (string, error)
    GetVersionHistory(packageName string) ([]string, error)

    // Enhanced version information
    GetCurrentVersion() string
    GetLatestVersion() (*VersionInfo, error)
    GetAllPackageVersions() (map[string]string, error)
    GetLatestPackageVersions() (map[string]string, error)
    GetDetailedVersionHistory() ([]VersionInfo, error)

    // Update management
    CheckForUpdates() (*UpdateInfo, error)
    DownloadUpdate(version string) error
    InstallUpdate(version string) error
    GetUpdateChannel() string
    SetUpdateChannel(channel string) error

    // Version compatibility
    CheckCompatibility(projectPath string) (*CompatibilityResult, error)
    GetSupportedVersions() (map[string][]string, error)
    ValidateVersionRequirements(requirements map[string]string) (*VersionValidationResult, error)

    // Version caching
    CacheVersionInfo(info *VersionInfo) error
    GetCachedVersionInfo() (*VersionInfo, error)
    RefreshVersionCache() error
    ClearVersionCache() error

    // Release information
    GetReleaseNotes(version string) (*ReleaseNotes, error)
    GetChangeLog(fromVersion, toVersion string) (*ChangeLog, error)
    GetSecurityAdvisories(version string) ([]SecurityAdvisory, error)

    // Package management
    GetPackageInfo(packageName string) (*PackageInfo, error)
    GetPackageVersions(packageName string) ([]string, error)
    CheckPackageUpdates(packages map[string]string) (map[string]PackageUpdate, error)

    // Version configuration
    SetVersionConfig(config *VersionConfig) error
    GetVersionConfig() (*VersionConfig, error)
    SetAutoUpdate(enabled bool) error
    SetUpdateNotifications(enabled bool) error
}
```

### Configuration Manager

The `ConfigManager` handles comprehensive configuration management with multiple sources and validation.

```go
type ConfigManager interface {
    // Configuration loading and saving
    LoadConfig(path string) (*models.ProjectConfig, error)
    SaveConfig(config *models.ProjectConfig, path string) error
    LoadDefaults() (*models.ProjectConfig, error)
    
    // Settings management
    GetSetting(key string) (interface{}, error)
    SetSetting(key string, value interface{}) error
    ValidateSettings() error
    
    // Configuration sources
    LoadFromFile(path string) (*models.ProjectConfig, error)
    LoadFromEnvironment() (*models.ProjectConfig, error)
    MergeConfigurations(configs ...*models.ProjectConfig) *models.ProjectConfig
    
    // Configuration validation
    ValidateConfig(config *models.ProjectConfig) error
    GetConfigSchema() *ConfigSchema
}
```

### Template Manager

The `TemplateManager` provides comprehensive template management including discovery, validation, and processing.

```go
type TemplateManager interface {
    // Template discovery
    ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
    GetTemplateInfo(name string) (*TemplateInfo, error)
    SearchTemplates(query string) ([]TemplateInfo, error)
    
    // Template validation
    ValidateTemplate(path string) (*TemplateValidationResult, error)
    ValidateTemplateStructure(template *TemplateInfo) error
    
    // Template processing
    ProcessTemplate(templateName string, config *models.ProjectConfig, outputPath string) error
    ProcessCustomTemplate(templatePath string, config *models.ProjectConfig, outputPath string) error
    
    // Template metadata
    GetTemplateMetadata(name string) (*TemplateMetadata, error)
    GetTemplateDependencies(name string) ([]string, error)
    GetTemplateCompatibility(name string) (*CompatibilityInfo, error)
}
```

### Validation Engine

The `ValidationEngine` provides comprehensive project validation with auto-fix capabilities.

```go
type ValidationEngine interface {
    // Basic validation methods
    ValidateProject(path string) (*models.ValidationResult, error)
    ValidatePackageJSON(path string) error
    ValidateGoMod(path string) error
    ValidateDockerfile(path string) error
    ValidateYAML(path string) error
    ValidateJSON(path string) error
    ValidateTemplate(path string) error

    // Enhanced project validation
    ValidateProjectStructure(path string) (*StructureValidationResult, error)
    ValidateProjectDependencies(path string) (*DependencyValidationResult, error)
    ValidateProjectSecurity(path string) (*SecurityValidationResult, error)
    ValidateProjectQuality(path string) (*QualityValidationResult, error)

    // Configuration validation
    ValidateConfiguration(config *models.ProjectConfig) (*ConfigValidationResult, error)
    ValidateConfigurationSchema(config any, schema *ConfigSchema) error
    ValidateConfigurationValues(config *models.ProjectConfig) (*ConfigValidationResult, error)

    // Template validation (enhanced versions)
    ValidateTemplateAdvanced(path string) (*TemplateValidationResult, error)
    ValidateTemplateMetadata(metadata *TemplateMetadata) error
    ValidateTemplateStructure(path string) (*StructureValidationResult, error)
    ValidateTemplateVariables(variables map[string]TemplateVariable) error

    // Validation options
    SetValidationRules(rules []ValidationRule) error
    GetValidationRules() []ValidationRule
    AddValidationRule(rule ValidationRule) error
    RemoveValidationRule(ruleID string) error

    // Auto-fix capabilities
    FixValidationIssues(path string, issues []ValidationIssue) (*FixResult, error)
    GetFixableIssues(issues []ValidationIssue) []ValidationIssue
    PreviewFixes(path string, issues []ValidationIssue) (*FixPreview, error)
    ApplyFix(path string, fix Fix) error

    // Validation reporting
    GenerateValidationReport(result *ValidationResult, format string) ([]byte, error)
    GetValidationSummary(results []*ValidationResult) (*ValidationSummary, error)
}
```

### Audit Engine

The `AuditEngine` provides comprehensive security, quality, and compliance auditing.

```go
type AuditEngine interface {
    // Security auditing
    AuditSecurity(path string) (*SecurityAuditResult, error)
    ScanVulnerabilities(path string) (*VulnerabilityReport, error)
    CheckSecurityPolicies(path string) (*PolicyComplianceResult, error)
    
    // Quality auditing
    AuditCodeQuality(path string) (*QualityAuditResult, error)
    CheckBestPractices(path string) (*BestPracticesResult, error)
    AnalyzeDependencies(path string) (*DependencyAnalysisResult, error)
    
    // License auditing
    AuditLicenses(path string) (*LicenseAuditResult, error)
    CheckLicenseCompatibility(path string) (*LicenseCompatibilityResult, error)
    
    // Performance auditing
    AuditPerformance(path string) (*PerformanceAuditResult, error)
    AnalyzeBundleSize(path string) (*BundleAnalysisResult, error)
}
```

### Cache Manager

The `CacheManager` provides comprehensive cache management for offline operation and performance optimization.

```go
type CacheManager interface {
    // Cache operations
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
    Clean() error
    
    // Cache information
    GetStats() (*CacheStats, error)
    GetSize() (int64, error)
    GetLocation() string
    
    // Cache validation
    ValidateCache() error
    RepairCache() error
    
    // Offline support
    EnableOfflineMode() error
    DisableOfflineMode() error
    IsOfflineMode() bool
}
```

**Key Features:**

- Multi-registry support (npm, Go modules, GitHub releases)
- Version caching for performance
- Automatic version updates
- Security vulnerability checking
- Version compatibility validation

**Usage:**

```go
versionManager := version.NewManager()

// Get latest Node.js version
nodeVersion, err := versionManager.GetLatestNodeVersion()
if err != nil {
    return fmt.Errorf("failed to get Node.js version: %w", err)
}

// Get latest NPM package version
reactVersion, err := versionManager.GetLatestNPMPackage("react")
if err != nil {
    return fmt.Errorf("failed to get React version: %w", err)
}
```

## Package Structure

The API is organized into the following packages:

### pkg/interfaces

Contains all interface definitions for dependency injection and testing.

### pkg/models

Contains comprehensive data models and configuration structures for all generator operations.

```go
type ProjectConfig struct {
    Name         string            `json:"name" yaml:"name"`
    Organization string            `json:"organization" yaml:"organization"`
    Description  string            `json:"description" yaml:"description"`
    License      string            `json:"license" yaml:"license"`
    Author       string            `json:"author,omitempty" yaml:"author,omitempty"`
    Email        string            `json:"email,omitempty" yaml:"email,omitempty"`
    Repository   string            `json:"repository,omitempty" yaml:"repository,omitempty"`
    OutputPath   string            `json:"output_path" yaml:"output_path"`
    Components   Components        `json:"components" yaml:"components"`
    Versions     *VersionConfig    `json:"versions,omitempty" yaml:"versions,omitempty"`
    CustomVars   map[string]interface{} `json:"custom_vars,omitempty" yaml:"custom_vars,omitempty"`
    GenerateOptions GenerateOptions `json:"generate_options,omitempty" yaml:"generate_options,omitempty"`
    GeneratedAt     time.Time       `json:"generated_at,omitempty" yaml:"generated_at,omitempty"`
    GeneratorVersion string         `json:"generator_version,omitempty" yaml:"generator_version,omitempty"`
}

type GenerateOptions struct {
    Force           bool     `json:"force" yaml:"force"`
    Minimal         bool     `json:"minimal" yaml:"minimal"`
    Offline         bool     `json:"offline" yaml:"offline"`
    UpdateVersions  bool     `json:"update_versions" yaml:"update_versions"`
    SkipValidation  bool     `json:"skip_validation" yaml:"skip_validation"`
    BackupExisting  bool     `json:"backup_existing" yaml:"backup_existing"`
    IncludeExamples bool     `json:"include_examples" yaml:"include_examples"`
    Templates       []string `json:"templates" yaml:"templates"`
    DryRun          bool     `json:"dry_run" yaml:"dry_run"`
    NonInteractive  bool     `json:"non_interactive" yaml:"non_interactive"`
}

type ValidationOptions struct {
    Verbose        bool     `json:"verbose"`
    Fix            bool     `json:"fix"`
    Report         bool     `json:"report"`
    ReportFormat   string   `json:"report_format"`
    Rules          []string `json:"rules"`
    IgnoreWarnings bool     `json:"ignore_warnings"`
    OutputFile     string   `json:"output_file"`
    Strict         bool     `json:"strict"`
    ShowFixes      bool     `json:"show_fixes"`
}

type AuditOptions struct {
    Security     bool     `json:"security"`
    Quality      bool     `json:"quality"`
    Licenses     bool     `json:"licenses"`
    Performance  bool     `json:"performance"`
    OutputFormat string   `json:"output_format"`
    OutputFile   string   `json:"output_file"`
    Detailed     bool     `json:"detailed"`
    FailOnHigh   bool     `json:"fail_on_high"`
    FailOnMedium bool     `json:"fail_on_medium"`
    MinScore     float64  `json:"min_score"`
}

type ValidationResult struct {
    Valid           bool               `json:"valid"`
    Issues          []ValidationIssue  `json:"issues"`
    Warnings        []ValidationIssue  `json:"warnings"`
    Summary         ValidationSummary  `json:"summary"`
    FixSuggestions  []FixSuggestion    `json:"fix_suggestions"`
}

type AuditResult struct {
    ProjectPath     string                  `json:"project_path"`
    AuditTime       time.Time              `json:"audit_time"`
    Security        *SecurityAuditResult    `json:"security,omitempty"`
    Quality         *QualityAuditResult     `json:"quality,omitempty"`
    Licenses        *LicenseAuditResult     `json:"licenses,omitempty"`
    Performance     *PerformanceAuditResult `json:"performance,omitempty"`
    OverallScore    float64                `json:"overall_score"`
    Recommendations []string               `json:"recommendations"`
}
```

### pkg/cli

Implements the CLI interface for interactive operations.

### pkg/template

Implements the template engine for processing templates.

### pkg/version

Implements version management functionality.

### pkg/filesystem

Implements file system operations.

## CLI Interface

### Interactive Project Configuration

The CLI interface provides a user-friendly way to collect project configuration:

```go
// Create CLI instance
cli := cli.NewCLI(configManager, validator)

// Collect project details interactively
config, err := cli.PromptProjectDetails()
if err != nil {
    log.Fatal(err)
}

// Select components
components, err := cli.SelectComponents()
if err != nil {
    log.Fatal(err)
}

// Confirm generation
if cli.ConfirmGeneration(config) {
    // Proceed with project generation
}
```

### Component Selection

The CLI supports selection of various project components:

- **Frontend Applications**: Next.js main app, landing page, admin dashboard
- **Backend Services**: Go API server with Gin framework
- **Mobile Applications**: Android Kotlin app, iOS Swift app
- **Infrastructure**: Docker configurations, Kubernetes manifests, Terraform

### Configuration Preview

Before generation, the CLI provides a comprehensive preview of:

- Project metadata (name, organization, description)
- Selected components and their dependencies
- Package versions that will be used
- Directory structure that will be created
- Available build commands

## Template Engine

### Template Processing

The template engine supports Go's text/template syntax with custom functions:

```go
// Register custom functions
engine.RegisterFunctions(template.FuncMap{
    "toUpper": strings.ToUpper,
    "toLower": strings.ToLower,
    "replace": strings.ReplaceAll,
})

// Process template with custom functions
content, err := engine.ProcessTemplate("template.tmpl", config)
```

### Custom Template Functions

Built-in template functions include:

- String manipulation (case conversion, replacement)
- Date and time formatting
- Mathematical operations
- Version formatting and comparison
- Security-related functions

### Template Variables

Templates have access to the complete project configuration:

```go
// In templates:
{{.Name}}                    // Project name
{{.Organization}}            // Organization name
{{nodeVersion .}}            // Node.js version
{{nextjsVersion .}}          // Next.js version
{{.Components.Frontend.NextJS.App}} // Component selection
```

### Conditional Rendering

Templates support conditional rendering based on selected components:

```go
{{if .Components.Frontend.NextJS.App}}
// Frontend-specific configuration
{{end}}

{{if .Components.Backend.GoGin}}
// Backend-specific configuration
{{end}}
```

## Version Management

### Registry Support

The version manager supports multiple package registries:

- **NPM Registry**: Node.js packages and frameworks
- **Go Module Registry**: Go packages and modules
- **GitHub Releases**: GitHub-hosted projects
- **Custom Registries**: Extensible registry support

### Version Caching

Version information is cached for performance:

```go
// Cache version for future use
err := versionManager.CacheVersion("react", "18.2.0")

// Retrieve cached version
version, found := versionManager.GetCachedVersion("react")
if found {
    fmt.Printf("Cached React version: %s\n", version)
}
```

### Version Updates

Automatic version updates keep projects current:

```go
// Update all versions
versionConfig, err := versionManager.UpdateVersionsConfig()
if err != nil {
    return err
}

// Access updated versions
fmt.Printf("Latest Node.js: %s\n", versionConfig.Node)
fmt.Printf("Latest React: %s\n", versionConfig.React)
```

## Configuration Management

### Configuration Loading

The config manager handles configuration loading and validation:

```go
configManager := config.NewManager()

// Load default configuration
defaultConfig, err := configManager.LoadDefaults()

// Load configuration from file
config, err := configManager.LoadConfig("project.yaml")

// Merge configurations
mergedConfig := configManager.MergeConfigs(defaultConfig, config)
```

### Configuration Validation

All configurations are validated before use:

```go
// Validate configuration
err := configManager.ValidateConfig(config)
if err != nil {
    return fmt.Errorf("invalid configuration: %w", err)
}
```

### Configuration Persistence

Configurations can be saved for reuse:

```go
// Save configuration to file
err := configManager.SaveConfig(config, "project.yaml")
if err != nil {
    return fmt.Errorf("failed to save config: %w", err)
}
```

## File System Operations

### Project Creation

The file system generator creates complete project structures:

```go
fsGenerator := filesystem.NewGenerator()

// Create entire project
err := fsGenerator.CreateProject(config, outputPath)
if err != nil {
    return fmt.Errorf("failed to create project: %w", err)
}
```

### Directory Operations

Individual directory operations are supported:

```go
// Create directory with proper permissions
err := fsGenerator.CreateDirectory("/path/to/dir")

// Ensure directory exists
err := fsGenerator.EnsureDirectory("/path/to/dir")

// Check if file exists
exists := fsGenerator.FileExists("/path/to/file")
```

### File Operations

File operations include writing and copying:

```go
// Write file with content and permissions
err := fsGenerator.WriteFile("file.txt", content, 0644)

// Copy assets from source to destination
err := fsGenerator.CopyAssets("src/assets", "dest/assets")

// Create symbolic link
err := fsGenerator.CreateSymlink("target", "link")
```

## Usage Examples

### Complete Project Generation

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/cli"
    "github.com/cuesoftinc/open-source-project-generator/pkg/config"
    "github.com/cuesoftinc/open-source-project-generator/pkg/template"
    "github.com/cuesoftinc/open-source-project-generator/pkg/validation"
    "github.com/cuesoftinc/open-source-project-generator/pkg/version"
    "github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
)

func main() {
    // Initialize components
    configManager := config.NewManager()
    validator := validation.NewEngine()
    versionManager := version.NewManager()
    templateEngine := template.NewEngineWithVersionManager(versionManager)
    fsGenerator := filesystem.NewGenerator()
    cliInterface := cli.NewCLI(configManager, validator)
    
    // Collect project configuration
    projectConfig, err := cliInterface.PromptProjectDetails()
    if err != nil {
        log.Fatalf("Failed to collect project details: %v", err)
    }
    
    // Validate configuration
    if err := configManager.ValidateConfig(projectConfig); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    // Confirm generation
    if !cliInterface.ConfirmGeneration(projectConfig) {
        fmt.Println("Generation canceled by user")
        return
    }
    
    // Generate project
    if err := templateEngine.ProcessDirectory("templates", projectConfig.OutputPath, projectConfig); err != nil {
        log.Fatalf("Failed to generate project: %v", err)
    }
    
    // Validate generated project
    result, err := validator.ValidateProject(projectConfig.OutputPath)
    if err != nil {
        log.Fatalf("Failed to validate project: %v", err)
    }
    
    if result.Valid {
        fmt.Println("✅ Project generated and validated successfully!")
    } else {
        fmt.Printf("⚠️ Project generated with %d validation issues\n", len(result.Issues))
    }
}
```

### Custom Template Function Registration

```go
// Register custom template functions
templateEngine.RegisterFunctions(template.FuncMap{
    "formatVersion": func(version string) string {
        // Custom version formatting logic
        return "^" + version
    },
    "generateSecret": func() string {
        // Generate secure random string
        return generateRandomString(32)
    },
    "formatDate": func(format string) string {
        return time.Now().Format(format)
    },
})
```

### Version Management Integration

```go
// Create version manager with caching
versionManager := version.NewManagerWithCache(cacheDir)

// Get latest versions for multiple packages
packages := []string{"react", "next", "typescript", "tailwindcss"}
versions := make(map[string]string)

for _, pkg := range packages {
    version, err := versionManager.GetLatestNPMPackage(pkg)
    if err != nil {
        log.Printf("Warning: Could not get version for %s: %v", pkg, err)
        continue
    }
    versions[pkg] = version
}

// Update project configuration with latest versions
projectConfig.Versions.Packages = versions
```

## Error Handling

### Error Types

The API uses structured error handling with specific error types:

```go
// Configuration errors
type ConfigError struct {
    Field   string
    Value   string
    Message string
}

// Template errors
type TemplateError struct {
    Template string
    Line     int
    Message  string
}

// Validation errors
type ValidationError struct {
    File     string
    Rule     string
    Severity string
    Message  string
}
```

### Error Handling Patterns

```go
// Handle specific error types
if configErr, ok := err.(*ConfigError); ok {
    fmt.Printf("Configuration error in field %s: %s\n", configErr.Field, configErr.Message)
}

// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to process template %s: %w", templatePath, err)
}

// Validate and handle multiple errors
if result, err := validator.ValidateProject(path); err != nil {
    return err
} else if !result.Valid {
    for _, issue := range result.Issues {
        if issue.Severity == "error" {
            return fmt.Errorf("validation error: %s", issue.Message)
        }
    }
}
```

### Recovery and Rollback

The API provides recovery mechanisms for failed operations:

```go
// Create backup before operations
backup, err := fsGenerator.CreateBackup(projectPath)
if err != nil {
    return fmt.Errorf("failed to create backup: %w", err)
}

// Attempt operation with rollback on failure
if err := templateEngine.ProcessDirectory(templateDir, outputDir, config); err != nil {
    // Rollback on failure
    if rollbackErr := fsGenerator.RestoreBackup(backup); rollbackErr != nil {
        return fmt.Errorf("operation failed and rollback failed: %v (original error: %w)", rollbackErr, err)
    }
    return fmt.Errorf("operation failed and was rolled back: %w", err)
}
```

## Performance Considerations

### Caching

The API implements multiple levels of caching:

- **Template Caching**: Parsed templates are cached for reuse
- **Version Caching**: Package versions are cached to reduce registry calls
- **Render Caching**: Rendered template results are cached

### Parallel Processing

Where appropriate, operations are parallelized:

```go
// Process multiple templates concurrently
var wg sync.WaitGroup
for _, template := range templates {
    wg.Add(1)
    go func(tmpl string) {
        defer wg.Done()
        processTemplate(tmpl, config)
    }(template)
}
wg.Wait()
```

### Memory Management

The API is designed for efficient memory usage:

- Streaming file operations for large files
- Proper resource cleanup with defer statements
- Memory-efficient template processing

## Security Considerations

### Input Validation

All user input is validated:

```go
// Validate project name
if !isValidProjectName(config.Name) {
    return fmt.Errorf("invalid project name: %s", config.Name)
}

// Sanitize file paths
outputPath := filepath.Clean(config.OutputPath)
if !isWithinAllowedPath(outputPath) {
    return fmt.Errorf("output path not allowed: %s", outputPath)
}
```

### Template Security

Template processing includes security measures:

- Restricted template functions
- Path traversal prevention
- Input sanitization
- Secure defaults

### Dependency Security

Version management includes security scanning:

```go
// Check for security vulnerabilities
vulnerabilities, err := versionManager.CheckSecurity("package-name", "1.0.0")
if err != nil {
    return fmt.Errorf("security check failed: %w", err)
}

if len(vulnerabilities) > 0 {
    for _, vuln := range vulnerabilities {
        fmt.Printf("Security vulnerability: %s (severity: %s)\n", vuln.Description, vuln.Severity)
    }
}
```

---

This API documentation provides comprehensive coverage of the Open Source Project Generator's public interfaces and usage patterns. For additional examples and implementation details, refer to the source code and test files in the respective packages.
