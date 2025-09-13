# API Documentation

This document provides comprehensive API documentation for the Open Source Template Generator.

## Table of Contents

- [Overview](#overview)
- [Core Interfaces](#core-interfaces)
- [Package Structure](#package-structure)
- [CLI Interface](#cli-interface)
- [Template Engine](#template-engine)
- [Validation Engine](#validation-engine)
- [Version Management](#version-management)
- [Configuration Management](#configuration-management)
- [File System Operations](#file-system-operations)
- [Usage Examples](#usage-examples)
- [Error Handling](#error-handling)

## Overview

The Open Source Template Generator provides a comprehensive API for generating production-ready project structures. The API is organized into several core interfaces that handle different aspects of project generation:

- **CLI Interface**: Interactive command-line operations
- **Template Engine**: Template processing and rendering
- **Config Manager**: Configuration loading and validation
- **File System Generator**: File and directory operations
- **Version Manager**: Package version management and updates

## Core Interfaces

### CLIInterface

The `CLIInterface` provides interactive command-line functionality for collecting project configuration and managing user interactions.

```go
type CLIInterface interface {
    Run() error
    PromptProjectDetails() (*models.ProjectConfig, error)
    SelectComponents() ([]string, error)
    ConfirmGeneration(*models.ProjectConfig) bool
}
```

**Key Features:**

- Interactive project configuration collection
- Component selection with dependency validation
- Configuration preview and confirmation
- Progress indication and user feedback

**Usage:**

```go
cli := cli.NewCLI(configManager, validator)
config, err := cli.PromptProjectDetails()
if err != nil {
    return fmt.Errorf("failed to collect project details: %w", err)
}

if cli.ConfirmGeneration(config) {
    // Proceed with generation
}
```

### TemplateEngine

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

### VersionManager

The `VersionManager` handles package version management and updates from various registries.

```go
type VersionManager interface {
    GetLatestNodeVersion() (string, error)
    GetLatestGoVersion() (string, error)
    GetLatestNPMPackage(packageName string) (string, error)
    GetLatestGoModule(moduleName string) (string, error)
    UpdateVersionsConfig() (*models.VersionConfig, error)
    GetLatestGitHubRelease(owner, repo string) (string, error)
    CacheVersion(key, version string) error
    GetCachedVersion(key string) (string, bool)
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

Contains data models and configuration structures.

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
    CustomVars   map[string]string `json:"custom_vars,omitempty" yaml:"custom_vars,omitempty"`
    // ... additional fields
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
{{.Versions.Node}}          // Node.js version
{{.Versions.NextJS}}        // Next.js version
{{.Components.Frontend.MainApp}} // Component selection
```

### Conditional Rendering

Templates support conditional rendering based on selected components:

```go
{{if .Components.Frontend.MainApp}}
// Frontend-specific configuration
{{end}}

{{if .Components.Backend.API}}
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
    
    "github.com/open-source-template-generator/pkg/cli"
    "github.com/open-source-template-generator/pkg/config"
    "github.com/open-source-template-generator/pkg/template"
    "github.com/open-source-template-generator/pkg/validation"
    "github.com/open-source-template-generator/pkg/version"
    "github.com/open-source-template-generator/pkg/filesystem"
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
        fmt.Println("Generation cancelled by user")
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

### Custom Template Functions

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

This API documentation provides comprehensive coverage of the Open Source Template Generator's public interfaces and usage patterns. For additional examples and implementation details, refer to the source code and test files in the respective packages.
