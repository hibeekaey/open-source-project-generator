# API Reference

This document provides a developer-focused API reference for the Open Source Project Generator.

## Core Interfaces

### CLI Interface

The main CLI interface provides comprehensive command-line functionality:

```go
type CLIInterface interface {
    // Core operations
    Run(args []string) error
    PromptProjectDetails() (*models.ProjectConfig, error)
    SelectComponents() ([]string, error)
    ConfirmGeneration(*models.ProjectConfig) bool

    // Non-interactive operations
    GenerateFromConfig(path string, options GenerateOptions) error
    ValidateProject(path string, options ValidationOptions) (*ValidationResult, error)
    AuditProject(path string, options AuditOptions) (*AuditResult, error)

    // Template operations
    ListTemplates(filter TemplateFilter) ([]TemplateInfo, error)
    GetTemplateInfo(name string) (*TemplateInfo, error)
    ValidateTemplate(path string) (*TemplateValidationResult, error)

    // Configuration operations
    ShowConfig() error
    SetConfig(key, value string) error
    LoadConfiguration(sources []string) (*models.ProjectConfig, error)

    // Version operations
    ShowVersion(options VersionOptions) error
    CheckUpdates() (*UpdateInfo, error)
    GetPackageVersions() (map[string]string, error)
}
```

### Template Engine

Handles template processing and rendering:

```go
type TemplateEngine interface {
    ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error)
    ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error
    RegisterFunctions(funcMap template.FuncMap)
    LoadTemplate(templatePath string) (*template.Template, error)
    RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error)
}
```

### Configuration Manager

Handles configuration loading and validation:

```go
type ConfigManager interface {
    LoadConfig(path string) (*models.ProjectConfig, error)
    SaveConfig(config *models.ProjectConfig, path string) error
    LoadDefaults() (*models.ProjectConfig, error)
    GetSetting(key string) (interface{}, error)
    SetSetting(key string, value interface{}) error
    ValidateConfig(config *models.ProjectConfig) error
}
```

### Version Manager

Provides version management and update capabilities:

```go
type VersionManager interface {
    GetLatestNodeVersion() (string, error)
    GetLatestGoVersion() (string, error)
    GetLatestNPMPackage(packageName string) (string, error)
    GetLatestGoModule(moduleName string) (string, error)
    UpdateVersionsConfig() (*models.VersionConfig, error)
    CheckForUpdates() (*UpdateInfo, error)
    GetPackageVersions() (map[string]string, error)
}
```

### Validation Engine

Provides comprehensive project validation:

```go
type ValidationEngine interface {
    ValidateProject(path string) (*models.ValidationResult, error)
    ValidateProjectStructure(path string) (*StructureValidationResult, error)
    ValidateProjectDependencies(path string) (*DependencyValidationResult, error)
    ValidateConfiguration(config *models.ProjectConfig) (*ConfigValidationResult, error)
    FixValidationIssues(path string, issues []ValidationIssue) (*FixResult, error)
}
```

### Audit Engine

Provides security and quality auditing:

```go
type AuditEngine interface {
    AuditSecurity(path string) (*SecurityAuditResult, error)
    AuditCodeQuality(path string) (*QualityAuditResult, error)
    AuditLicenses(path string) (*LicenseAuditResult, error)
    AuditPerformance(path string) (*PerformanceAuditResult, error)
}
```

## Data Models

### Project Configuration

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
}
```

### Component Selection

```go
type Components struct {
    Frontend      FrontendComponents      `json:"frontend" yaml:"frontend"`
    Backend       BackendComponents       `json:"backend" yaml:"backend"`
    Mobile        MobileComponents        `json:"mobile" yaml:"mobile"`
    Infrastructure InfrastructureComponents `json:"infrastructure" yaml:"infrastructure"`
}

type FrontendComponents struct {
    MainApp         bool `json:"main_app" yaml:"main_app"`
    Home            bool `json:"home" yaml:"home"`
    Admin           bool `json:"admin" yaml:"admin"`
    SharedComponents bool `json:"shared_components" yaml:"shared_components"`
}

type BackendComponents struct {
    API      bool `json:"api" yaml:"api"`
    Auth     bool `json:"auth" yaml:"auth"`
    Database bool `json:"database" yaml:"database"`
}

type MobileComponents struct {
    Android bool `json:"android" yaml:"android"`
    iOS     bool `json:"ios" yaml:"ios"`
    Shared  bool `json:"shared" yaml:"shared"`
}

type InfrastructureComponents struct {
    Docker      bool `json:"docker" yaml:"docker"`
    Kubernetes  bool `json:"kubernetes" yaml:"kubernetes"`
    Terraform   bool `json:"terraform" yaml:"terraform"`
    Monitoring  bool `json:"monitoring" yaml:"monitoring"`
}
```

### Generation Options

```go
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
```

### Validation Options

```go
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
```

### Audit Options

```go
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
```

## Usage Examples

### Basic CLI Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/cli"
    "github.com/cuesoftinc/open-source-project-generator/pkg/config"
    "github.com/cuesoftinc/open-source-project-generator/pkg/validation"
)

func main() {
    // Initialize components
    configManager := config.NewManager()
    validator := validation.NewEngine()
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
}
```

### Template Processing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/template"
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func main() {
    // Create template engine
    engine := template.NewEngine()
    
    // Register custom functions
    engine.RegisterFunctions(template.FuncMap{
        "toUpper": strings.ToUpper,
        "toLower": strings.ToLower,
        "replace": strings.ReplaceAll,
    })
    
    // Process template
    config := &models.ProjectConfig{
        Name: "my-project",
        Organization: "my-org",
    }
    
    content, err := engine.ProcessTemplate("templates/package.json.tmpl", config)
    if err != nil {
        log.Fatalf("Failed to process template: %v", err)
    }
    
    fmt.Printf("Generated content: %s\n", string(content))
}
```

### Configuration Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/config"
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func main() {
    // Create config manager
    configManager := config.NewManager()
    
    // Load default configuration
    defaultConfig, err := configManager.LoadDefaults()
    if err != nil {
        log.Fatalf("Failed to load defaults: %v", err)
    }
    
    // Load configuration from file
    config, err := configManager.LoadConfig("project.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Merge configurations
    mergedConfig := configManager.MergeConfigs(defaultConfig, config)
    
    // Validate configuration
    if err := configManager.ValidateConfig(mergedConfig); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    fmt.Printf("Configuration loaded: %s\n", mergedConfig.Name)
}
```

### Version Management

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/version"
)

func main() {
    // Create version manager
    versionManager := version.NewManager()
    
    // Get latest Node.js version
    nodeVersion, err := versionManager.GetLatestNodeVersion()
    if err != nil {
        log.Fatalf("Failed to get Node.js version: %v", err)
    }
    
    // Get latest NPM package version
    reactVersion, err := versionManager.GetLatestNPMPackage("react")
    if err != nil {
        log.Fatalf("Failed to get React version: %v", err)
    }
    
    fmt.Printf("Latest Node.js: %s\n", nodeVersion)
    fmt.Printf("Latest React: %s\n", reactVersion)
}
```

### Project Validation

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/validation"
)

func main() {
    // Create validation engine
    validator := validation.NewEngine()
    
    // Validate project
    result, err := validator.ValidateProject("./my-project")
    if err != nil {
        log.Fatalf("Failed to validate project: %v", err)
    }
    
    if result.Valid {
        fmt.Println("✅ Project validation passed")
    } else {
        fmt.Printf("⚠️ Project validation failed with %d issues\n", len(result.Issues))
        for _, issue := range result.Issues {
            fmt.Printf("- %s: %s\n", issue.Severity, issue.Message)
        }
    }
}
```

### Security Auditing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/audit"
)

func main() {
    // Create audit engine
    auditEngine := audit.NewEngine()
    
    // Perform security audit
    result, err := auditEngine.AuditSecurity("./my-project")
    if err != nil {
        log.Fatalf("Failed to audit project: %v", err)
    }
    
    fmt.Printf("Security Score: %.1f/10\n", result.Score)
    fmt.Printf("Vulnerabilities Found: %d\n", len(result.Vulnerabilities))
    
    for _, vuln := range result.Vulnerabilities {
        fmt.Printf("- %s (Severity: %s)\n", vuln.Description, vuln.Severity)
    }
}
```

## Error Handling

### Error Types

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

## Testing

### Unit Tests

```go
func TestTemplateProcessing(t *testing.T) {
    engine := template.NewEngine()
    config := &models.ProjectConfig{
        Name: "test-project",
    }
    
    content, err := engine.ProcessTemplate("test.tmpl", config)
    assert.NoError(t, err)
    assert.Contains(t, string(content), "test-project")
}
```

### Integration Tests

```go
func TestProjectGeneration(t *testing.T) {
    // Create temporary directory
    tempDir := t.TempDir()
    
    // Generate project
    err := generateProject(tempDir, testConfig)
    assert.NoError(t, err)
    
    // Validate generated files
    assert.FileExists(t, filepath.Join(tempDir, "package.json"))
    assert.FileExists(t, filepath.Join(tempDir, "README.md"))
}
```

This API reference provides the essential information needed for developers to integrate with and extend the Open Source Project Generator.
