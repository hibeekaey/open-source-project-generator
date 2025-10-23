# API Reference

Developer API documentation for the Open Source Project Generator.

## Table of Contents

- [Core Interfaces](#core-interfaces)
- [Key Types](#key-types)
- [Models](#models)
- [Usage Examples](#usage-examples)
- [Error Handling](#error-handling)

---

## Core Interfaces

### ProjectCoordinatorInterface

Main interface for project generation orchestration.

**Location:** `pkg/interfaces/coordinator.go`

```go
type ProjectCoordinatorInterface interface {
    Generate(ctx context.Context, config interface{}) (interface{}, error)
    DryRun(ctx context.Context, config interface{}) (interface{}, error)
    Validate(config interface{}) error
}
```

**Methods:**

- `Generate` - Generate a complete project from configuration
- `DryRun` - Preview what would be generated without creating files
- `Validate` - Validate configuration before generation

**Usage:**

```go
coordinator := orchestrator.NewProjectCoordinator(logger)

result, err := coordinator.Generate(ctx, config)
if err != nil {
    log.Fatal(err)
}

generationResult := result.(*models.GenerationResult)
```

### BootstrapExecutorInterface

Interface for bootstrap tool executors.

**Location:** `pkg/interfaces/executor.go`

```go
type BootstrapExecutorInterface interface {
    Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error)
    SupportsComponent(componentType string) bool
    GetDefaultFlags(componentType string) []string
}
```

**Methods:**

- `Execute` - Execute the bootstrap tool with given specification
- `SupportsComponent` - Check if executor supports a component type
- `GetDefaultFlags` - Get default flags for a component type

### StructureMapperInterface

Interface for mapping generated structures.

**Location:** `pkg/interfaces/mapper.go`

```go
type StructureMapperInterface interface {
    Map(ctx context.Context, source, target, componentType string) error
    UpdateImportPaths(targetPath, componentType string) error
}
```

**Methods:**

- `Map` - Map generated files to target structure
- `UpdateImportPaths` - Update import paths in generated files

---

## Key Types

### ProjectConfig

Project configuration structure.

**Location:** `pkg/models/project.go`

```go
type ProjectConfig struct {
    Name        string              `yaml:"name" json:"name"`
    Description string              `yaml:"description" json:"description"`
    OutputDir   string              `yaml:"output_dir" json:"output_dir"`
    Author      string              `yaml:"author,omitempty" json:"author,omitempty"`
    Email       string              `yaml:"email,omitempty" json:"email,omitempty"`
    License     string              `yaml:"license,omitempty" json:"license,omitempty"`
    Repository  string              `yaml:"repository,omitempty" json:"repository,omitempty"`
    Components  []ComponentConfig   `yaml:"components" json:"components"`
    Integration IntegrationConfig   `yaml:"integration" json:"integration"`
    Options     ProjectOptions      `yaml:"options" json:"options"`
}
```

**Fields:**

- `Name` - Project name (required)
- `Description` - Project description
- `OutputDir` - Output directory path (required)
- `Author` - Project author
- `Email` - Author email
- `License` - License type (e.g., "MIT", "Apache-2.0")
- `Repository` - Repository URL
- `Components` - List of components to generate
- `Integration` - Integration settings
- `Options` - Generation options

### ComponentConfig

Component configuration structure.

```go
type ComponentConfig struct {
    Type    string                 `yaml:"type" json:"type"`
    Name    string                 `yaml:"name" json:"name"`
    Enabled bool                   `yaml:"enabled" json:"enabled"`
    Config  map[string]interface{} `yaml:"config" json:"config"`
}
```

**Fields:**

- `Type` - Component type (e.g., "nextjs", "go-backend")
- `Name` - Component name
- `Enabled` - Whether component is enabled
- `Config` - Component-specific configuration

### IntegrationConfig

Integration configuration structure.

```go
type IntegrationConfig struct {
    GenerateDockerCompose bool              `yaml:"generate_docker_compose" json:"generate_docker_compose"`
    GenerateScripts       bool              `yaml:"generate_scripts" json:"generate_scripts"`
    APIEndpoints          map[string]string `yaml:"api_endpoints" json:"api_endpoints"`
    SharedEnvironment     map[string]string `yaml:"shared_environment" json:"shared_environment"`
}
```

**Fields:**

- `GenerateDockerCompose` - Generate Docker Compose file
- `GenerateScripts` - Generate build/run scripts
- `APIEndpoints` - API endpoint configuration
- `SharedEnvironment` - Shared environment variables

### ProjectOptions

Generation options structure.

```go
type ProjectOptions struct {
    UseExternalTools bool `yaml:"use_external_tools" json:"use_external_tools"`
    DryRun           bool `yaml:"dry_run" json:"dry_run"`
    Verbose          bool `yaml:"verbose" json:"verbose"`
    CreateBackup     bool `yaml:"create_backup" json:"create_backup"`
    ForceOverwrite   bool `yaml:"force_overwrite" json:"force_overwrite"`
}
```

**Fields:**

- `UseExternalTools` - Use bootstrap tools or force fallback
- `DryRun` - Preview mode (don't create files)
- `Verbose` - Enable verbose logging
- `CreateBackup` - Create backup before overwriting
- `ForceOverwrite` - Force overwrite existing directory

---

## Models

### GenerationResult

Result of project generation.

```go
type GenerationResult struct {
    Success     bool               `json:"success"`
    ProjectRoot string             `json:"project_root"`
    Components  []*ComponentResult `json:"components"`
    Duration    time.Duration      `json:"duration"`
    Errors      []error            `json:"errors,omitempty"`
    Warnings    []string           `json:"warnings,omitempty"`
    DryRun      bool               `json:"dry_run"`
    LogFile     string             `json:"log_file,omitempty"`
}
```

**Fields:**

- `Success` - Whether generation succeeded
- `ProjectRoot` - Root directory of generated project
- `Components` - Results for each component
- `Duration` - Total generation time
- `Errors` - List of errors encountered
- `Warnings` - List of warnings
- `DryRun` - Whether this was a dry run
- `LogFile` - Path to log file

### ComponentResult

Result of component generation.

```go
type ComponentResult struct {
    Name        string        `json:"name"`
    Type        string        `json:"type"`
    Success     bool          `json:"success"`
    Method      string        `json:"method"`
    ToolUsed    string        `json:"tool_used"`
    OutputPath  string        `json:"output_path"`
    Duration    time.Duration `json:"duration"`
    Error       error         `json:"error,omitempty"`
    Warnings    []string      `json:"warnings,omitempty"`
    ManualSteps []string      `json:"manual_steps,omitempty"`
}
```

**Fields:**

- `Name` - Component name
- `Type` - Component type
- `Success` - Whether generation succeeded
- `Method` - Generation method ("bootstrap" or "fallback")
- `ToolUsed` - Tool that was used
- `OutputPath` - Path to generated component
- `Duration` - Generation time
- `Error` - Error if generation failed
- `Warnings` - List of warnings
- `ManualSteps` - Manual steps required (for fallback)

### ToolCheckResult

Result of tool availability check.

```go
type ToolCheckResult struct {
    AllAvailable bool                  `json:"all_available"`
    Tools        map[string]*ToolInfo  `json:"tools"`
    Missing      []string              `json:"missing"`
    Outdated     []string              `json:"outdated"`
    CheckedAt    time.Time             `json:"checked_at"`
}
```

**Fields:**

- `AllAvailable` - Whether all required tools are available
- `Tools` - Information about each tool
- `Missing` - List of missing tools
- `Outdated` - List of outdated tools
- `CheckedAt` - When check was performed

### ToolInfo

Information about a tool.

```go
type ToolInfo struct {
    Available        bool   `json:"available"`
    InstalledVersion string `json:"installed_version,omitempty"`
    MinVersion       string `json:"min_version,omitempty"`
    Path             string `json:"path,omitempty"`
}
```

**Fields:**

- `Available` - Whether tool is available
- `InstalledVersion` - Installed version
- `MinVersion` - Minimum required version
- `Path` - Path to tool executable

---

## Usage Examples

### Basic Project Generation

```go
package main

import (
    "context"
    "log"
    
    "github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
    "github.com/cuesoftinc/open-source-project-generator/pkg/logger"
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func main() {
    // Create logger
    log := logger.NewLogger()
    
    // Create configuration
    config := &models.ProjectConfig{
        Name:      "my-project",
        OutputDir: "./my-project",
        Components: []models.ComponentConfig{
            {
                Type:    "nextjs",
                Name:    "web-app",
                Enabled: true,
                Config: map[string]interface{}{
                    "typescript": true,
                    "tailwind":   true,
                },
            },
        },
        Integration: models.IntegrationConfig{
            GenerateDockerCompose: true,
            GenerateScripts:       true,
        },
        Options: models.ProjectOptions{
            UseExternalTools: true,
            CreateBackup:     true,
        },
    }
    
    // Create coordinator
    coordinator := orchestrator.NewProjectCoordinator(log)
    
    // Generate project
    result, err := coordinator.Generate(context.Background(), config)
    if err != nil {
        log.Fatal("Generation failed:", err)
    }
    
    // Check result
    genResult := result.(*models.GenerationResult)
    if genResult.Success {
        log.Info("Project generated successfully!")
        log.Info("Location:", genResult.ProjectRoot)
    } else {
        log.Error("Generation failed with errors:", genResult.Errors)
    }
}
```

### Dry Run

```go
// Preview what would be generated
result, err := coordinator.DryRun(context.Background(), config)
if err != nil {
    log.Fatal(err)
}

previewResult := result.(*models.PreviewResult)
log.Info("Would generate", len(previewResult.Components), "components")
for _, comp := range previewResult.Components {
    log.Info("Component:", comp.Name, "Type:", comp.Type)
}
```

### Check Tool Availability

```go
// Create tool discovery
toolDiscovery := orchestrator.NewToolDiscovery(log)

// Check specific tools
tools := []string{"npx", "go", "gradle"}
result, err := toolDiscovery.CheckRequirements(tools)
if err != nil {
    log.Fatal(err)
}

toolResult := result.(*models.ToolCheckResult)
if toolResult.AllAvailable {
    log.Info("All tools available!")
} else {
    log.Warn("Missing tools:", toolResult.Missing)
}
```

### Custom Executor

```go
package bootstrap

import (
    "context"
    "github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

type CustomExecutor struct {
    *BaseExecutor
}

func NewCustomExecutor(log *logger.Logger) *CustomExecutor {
    return &CustomExecutor{
        BaseExecutor: &BaseExecutor{logger: log},
    }
}

func (ce *CustomExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    // Custom implementation
    ce.logger.Info("Executing custom tool")
    
    // Build command
    spec.Tool = "custom-tool"
    spec.Flags = []string{"create", spec.Config["name"].(string)}
    
    // Execute
    return ce.BaseExecutor.Execute(ctx, spec)
}

func (ce *CustomExecutor) SupportsComponent(componentType string) bool {
    return componentType == "custom"
}

func (ce *CustomExecutor) GetDefaultFlags(componentType string) []string {
    return []string{"--default"}
}
```

---

## Error Handling

### Error Types

The generator uses typed errors for better error handling:

```go
import "github.com/cuesoftinc/open-source-project-generator/pkg/errors"

// Configuration errors
errors.NewConfigError("invalid configuration", err)

// Tool errors
errors.NewToolError("tool not found", err)

// File system errors
errors.NewFileSystemError("permission denied", err)

// Validation errors
errors.NewValidationError("invalid input", err)

// Security errors
errors.NewSecurityError("path traversal detected", err)
```

### Error Handling Example

```go
result, err := coordinator.Generate(ctx, config)
if err != nil {
    switch e := err.(type) {
    case *errors.ConfigError:
        log.Error("Configuration error:", e)
        // Handle configuration error
    case *errors.ToolError:
        log.Error("Tool error:", e)
        // Handle tool error
    case *errors.FileSystemError:
        log.Error("File system error:", e)
        // Handle file system error
    default:
        log.Error("Unknown error:", e)
    }
    return
}
```

### Checking Specific Errors

```go
import "errors"

if errors.Is(err, ErrToolNotFound) {
    log.Warn("Tool not found, using fallback")
    // Use fallback generator
}

if errors.Is(err, ErrConfigInvalid) {
    log.Error("Invalid configuration")
    // Show validation errors
}
```

---

## Logger Interface

### Logger Methods

```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Fatal(msg string, args ...interface{})
    SetLevel(level LogLevel)
}
```

### Logger Usage

```go
log := logger.NewLogger()

// Set log level
log.SetLevel(logger.DebugLevel)

// Log messages
log.Debug("Debug message", "key", "value")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message", "error", err)
```

---

## Configuration Validation

### Validator Interface

```go
type Validator interface {
    Validate(config *models.ProjectConfig) error
    ApplyDefaults(config *models.ProjectConfig) error
}
```

### Validator Usage

```go
validator := config.NewValidator()

// Validate configuration
if err := validator.Validate(config); err != nil {
    log.Fatal("Validation failed:", err)
}

// Apply defaults
if err := validator.ApplyDefaults(config); err != nil {
    log.Fatal("Failed to apply defaults:", err)
}
```

---

## Security

### Input Sanitization

```go
import "github.com/cuesoftinc/open-source-project-generator/pkg/security"

sanitizer := security.NewSanitizer()

// Sanitize path
cleanPath, err := sanitizer.SanitizePath(userPath)
if err != nil {
    log.Fatal("Invalid path:", err)
}

// Sanitize project name
cleanName, err := sanitizer.SanitizeName(userName)
if err != nil {
    log.Fatal("Invalid name:", err)
}
```

### Path Validation

```go
// Validate path is safe
if err := security.ValidatePath(path); err != nil {
    log.Fatal("Unsafe path:", err)
}

// Check for path traversal
if security.HasPathTraversal(path) {
    log.Fatal("Path traversal detected")
}
```

---

## See Also

- [Architecture](ARCHITECTURE.md) - System architecture
- [Adding Tools](ADDING_TOOLS.md) - Adding new bootstrap tools
- [Getting Started](GETTING_STARTED.md) - Installation and usage
- [Configuration Guide](CONFIGURATION.md) - Configuration options
