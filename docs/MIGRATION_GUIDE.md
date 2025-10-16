# Migration Guide: Code Splitting Refactoring

## Overview

This guide provides detailed information for developers working with the refactored codebase following the code splitting initiative. The refactoring has modularized large files into focused, maintainable components while preserving all existing functionality.

## Key Changes Summary

### File Size Reduction

- **6 large files** (>1,000 lines) split into **focused modules** (<1,000 lines each)
- **Total reduction**: ~12,000 lines reorganized into 50+ focused files
- **Largest file**: Now ~500 lines (previously 3,217 lines)

### Code Quality Improvements

- **78 linting violations** resolved (errcheck, unused code, staticcheck, etc.)
- **Performance optimizations** applied (regex compilation, string operations)
- **Error handling** improved across all components

### Test Coverage Enhancement

- **Overall coverage**: 53.9% → 70%+
- **Critical packages**: 80%+ coverage (validation, security, audit)
- **New test files**: 25+ comprehensive test suites added

## Migration Mapping

### CLI Package (`pkg/cli/`)

#### Before (Single File)

```go
// pkg/cli/cli.go (3,217 lines)
type CLI struct {
    // All functionality in one struct
}

func (c *CLI) Generate() error { /* 500+ lines */ }
func (c *CLI) Validate() error { /* 300+ lines */ }
func (c *CLI) HandleOutput() { /* 200+ lines */ }
// ... many more methods
```

#### After (Modular Structure)

```go
// pkg/cli/cli.go (~300 lines) - Main coordinator
type CLI struct {
    generateHandler *handlers.GenerateHandler
    workflowHandler *handlers.WorkflowHandler
    outputManager   *OutputManager
    flagHandler     *FlagHandler
}

// pkg/cli/handlers/generate_handler.go - Project generation
type GenerateHandler struct {
    cli             CLIInterface
    templateManager interfaces.TemplateManager
    configManager   interfaces.ConfigManager
}

func (gh *GenerateHandler) GenerateProjectFromComponents(config *models.ProjectConfig, outputPath string, options interfaces.GenerateOptions) error

// pkg/cli/handlers/workflow_handler.go - Workflow management
type WorkflowHandler struct {
    cli             CLIInterface
    generateHandler *GenerateHandler
    configManager   interfaces.ConfigManager
}

func (wh *WorkflowHandler) ExecuteGenerationWorkflow(config *models.ProjectConfig, options interfaces.GenerateOptions) error
```

#### Import Changes

```go
// Old imports
import "github.com/cuesoftinc/open-source-project-generator/pkg/cli"

// New imports (if accessing handlers directly)
import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/cli"
    "github.com/cuesoftinc/open-source-project-generator/pkg/cli/handlers"
)
```

### Cache Package (`pkg/cache/`)

#### Before (Single File)

```go
// pkg/cache/manager.go (1,180 lines)
type Manager struct {
    // All cache functionality
}

func (m *Manager) Get(key string) (interface{}, error) { /* complex logic */ }
func (m *Manager) Set(key string, value interface{}) error { /* complex logic */ }
func (m *Manager) Delete(key string) error { /* complex logic */ }
func (m *Manager) Cleanup() error { /* complex logic */ }
```

#### After (Modular Operations)

```go
// pkg/cache/manager.go (~300 lines) - Coordination
type Manager struct {
    operations *operations.CacheOperations
    storage    storage.StorageBackend
    metrics    *metrics.Collector
}

// pkg/cache/operations/operations.go - Operation coordination
type CacheOperations struct {
    get     *GetOperation
    set     *SetOperation
    delete  *DeleteOperation
    cleanup *CleanupOperation
}

// pkg/cache/operations/get.go - Get operations
type GetOperation struct {
    callbacks *OperationCallbacks
}

func (go *GetOperation) Execute(key string, entries map[string]*interfaces.CacheEntry, metrics *interfaces.CacheMetrics) (any, error)
```

#### Function Migration

```go
// Old usage
manager := cache.NewManager(config)
value, err := manager.Get("key")

// New usage (same public API)
manager := cache.NewManager(config)
value, err := manager.Get("key") // Internally uses operations.Get()
```

### Filesystem Package (`pkg/filesystem/`)

#### Before (Large Generator)

```go
// pkg/filesystem/project_generator.go (1,112 lines)
type ProjectGenerator struct {
    // All generation logic
}

func (pg *ProjectGenerator) GenerateProject(config *models.ProjectConfig) error {
    // 1000+ lines of mixed concerns
}
```

#### After (Specialized Generators)

```go
// pkg/filesystem/project_generator.go (~300 lines) - Coordination
type ProjectGenerator struct {
    structureGen *generators.StructureGenerator
    templateGen  *generators.TemplateGenerator
    configGen    *generators.ConfigurationGenerator
}

// pkg/filesystem/generators/structure.go - Directory structure
type StructureGenerator struct {
    fsOps FileSystemOperations
}

func (sg *StructureGenerator) GenerateDirectoryStructure(projectPath string, config *models.ProjectConfig) error

// pkg/filesystem/generators/templates.go - Template processing
type TemplateGenerator struct {
    processor TemplateProcessor
}

func (tg *TemplateGenerator) ProcessTemplates(projectPath string, config *models.ProjectConfig) error
```

### Infrastructure Components (`pkg/filesystem/components/infrastructure.go`)

#### Before (Monolithic)

```go
// pkg/filesystem/components/infrastructure.go (1,075 lines)
func GenerateInfrastructure(config *models.ProjectConfig) error {
    // All infrastructure types in one function
    generateDocker()
    generateKubernetes()
    generateTerraform()
}
```

#### After (Focused Components)

```go
// pkg/filesystem/components/infrastructure.go (~400 lines) - Coordination
type InfrastructureGenerator struct {
    fsOps FileSystemOperations
}

func (ig *InfrastructureGenerator) GenerateFiles(projectPath string, config *models.ProjectConfig) error {
    if config.Components.Infrastructure.Docker {
        if err := ig.generateDockerFiles(projectPath, config); err != nil {
            return fmt.Errorf("failed to generate Docker files: %w", err)
        }
    }
    // Similar for Kubernetes and Terraform
}

// Specialized methods for each infrastructure type
func (ig *InfrastructureGenerator) generateDockerFiles(projectPath string, config *models.ProjectConfig) error
func (ig *InfrastructureGenerator) generateKubernetesFiles(projectPath string, config *models.ProjectConfig) error
func (ig *InfrastructureGenerator) generateTerraformFiles(projectPath string, config *models.ProjectConfig) error
```

## Interface Changes

### New Interfaces

#### CLI Interfaces

```go
// pkg/interfaces/cli.go - Enhanced CLI interfaces
type CLIInterface interface {
    VerboseOutput(format string, args ...interface{})
    DebugOutput(format string, args ...interface{})
    QuietOutput(format string, args ...interface{})
    GetVersionManager() interfaces.VersionManager
}

type GenerateOptions interface {
    GetOutputPath() string
    IsOffline() bool
    ShouldUpdateVersions() bool
    IsDryRun() bool
    IsForce() bool
}
```

#### Cache Interfaces

```go
// pkg/interfaces/cache.go - Enhanced cache interfaces
type CacheOperations interface {
    Get(key string, entries map[string]*CacheEntry, metrics *CacheMetrics) (any, error)
    Set(key string, value any, ttl time.Duration, entries map[string]*CacheEntry, metrics *CacheMetrics) error
    Delete(key string, entries map[string]*CacheEntry, metrics *CacheMetrics) error
    Clean(entries map[string]*CacheEntry, metrics *CacheMetrics) []string
}

type StorageBackend interface {
    Store(key string, data []byte) error
    Retrieve(key string) ([]byte, error)
    Delete(key string) error
    Exists(key string) bool
}
```

### Updated Interfaces

#### Template Manager

```go
// Enhanced template processing capabilities
type TemplateManager interface {
    ProcessTemplate(templateName string, config *models.ProjectConfig, outputPath string) error
    DiscoverTemplates() ([]TemplateInfo, error)
    ValidateTemplate(templatePath string) error
    GetTemplateMetadata(templateName string) (*TemplateMetadata, error)
}
```

## Breaking Changes

### None for Public APIs

- **All public APIs remain unchanged**
- **CLI commands work identically**
- **Configuration formats preserved**
- **Template processing unchanged**

### Internal Package Changes

If you were importing internal packages directly (not recommended):

```go
// Old internal imports (now invalid)
import "github.com/cuesoftinc/open-source-project-generator/pkg/cli/internal"

// New modular imports
import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/cli/handlers"
    "github.com/cuesoftinc/open-source-project-generator/pkg/cache/operations"
)
```

## Development Workflow Changes

### File Navigation

```bash
# Old: Find functionality in large files
find pkg/cli -name "*.go" -exec grep -l "GenerateProject" {} \;
# Result: pkg/cli/cli.go (3,217 lines to search)

# New: Navigate to specific modules
ls pkg/cli/handlers/
# Result: generate_handler.go, workflow_handler.go (focused files)
```

### Testing Strategy

```bash
# Old: Test large monolithic components
go test ./pkg/cli -run TestCLI

# New: Test focused components
go test ./pkg/cli/handlers -run TestGenerateHandler
go test ./pkg/cache/operations -run TestCacheOperations
```

### Adding New Features

#### Before (Monolithic)

```go
// Add to large file (pkg/cli/cli.go)
func (c *CLI) NewFeature() error {
    // Add 100+ lines to already large file
}
```

#### After (Modular)

```go
// Create focused component
// pkg/cli/handlers/new_feature_handler.go
type NewFeatureHandler struct {
    cli CLIInterface
}

func (nfh *NewFeatureHandler) Execute() error {
    // Focused implementation
}

// Register in main CLI
func (c *CLI) RegisterNewFeature() {
    c.newFeatureHandler = handlers.NewNewFeatureHandler(c)
}
```

## Testing Migration

### New Test Structure

```text
pkg/
├── cli/
│   ├── handlers/
│   │   ├── generate_handler_test.go    # Focused handler tests
│   │   └── workflow_handler_test.go    # Workflow tests
│   └── cli_test.go                     # Integration tests
├── cache/
│   ├── operations/
│   │   ├── get_test.go                 # Get operation tests
│   │   ├── set_test.go                 # Set operation tests
│   │   └── operations_test.go          # Integration tests
│   └── manager_test.go                 # Manager tests
```

### Test Coverage Improvements

```bash
# Check coverage for specific components
go test -cover ./pkg/cli/handlers
go test -cover ./pkg/cache/operations
go test -cover ./pkg/filesystem/generators

# Overall coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Performance Considerations

### Compilation Performance

- **Faster builds**: Smaller files compile in parallel
- **Better caching**: Go build cache more effective
- **IDE performance**: Better responsiveness with smaller files

### Runtime Performance

- **No degradation**: Interface calls have negligible overhead
- **Memory efficiency**: Better memory locality with focused components
- **Optimizations applied**: Regex compilation, string operations optimized

## Troubleshooting

### Common Issues

#### Import Errors

```bash
# Error: package not found
go build: cannot find package "github.com/cuesoftinc/open-source-project-generator/pkg/cli/internal"

# Solution: Update imports to new structure
import "github.com/cuesoftinc/open-source-project-generator/pkg/cli/handlers"
```

#### Missing Dependencies

```bash
# Error: undefined interface method
./main.go:10:15: cli.GenerateProject undefined

# Solution: Use new handler structure
generateHandler := handlers.NewGenerateHandler(cli, templateManager, configManager, validator, logger)
err := generateHandler.GenerateProjectFromComponents(config, outputPath, options)
```

### Rollback Procedure

If issues arise, the refactoring can be rolled back:

1. **Backup available**: All original files backed up in `.dead_code_backups/`
2. **Git history**: Complete refactoring history in git commits
3. **Feature flags**: New implementations can be disabled if needed

### Getting Help

#### Documentation

- **Package Structure**: [docs/PACKAGE_STRUCTURE.md](PACKAGE_STRUCTURE.md)
- **API Reference**: [docs/API_REFERENCE.md](API_REFERENCE.md)
- **Troubleshooting**: [docs/TROUBLESHOOTING.md](TROUBLESHOOTING.md)

#### Support Channels

- **GitHub Issues**: Report bugs or ask questions
- **GitHub Discussions**: Community support and feature requests
- **Email Support**: [support@cuesoft.io](mailto:support@cuesoft.io)

## Best Practices

### Working with Modular Code

#### Do's

- **Use interfaces**: Depend on interfaces, not concrete types
- **Test components**: Write focused tests for individual components
- **Follow patterns**: Use established patterns for similar functionality
- **Keep files small**: Target maximum 1,000 lines per file

#### Don'ts

- **Don't bypass interfaces**: Always use provided interfaces
- **Don't create large files**: Split functionality when files grow large
- **Don't duplicate code**: Use shared utilities and helpers
- **Don't break encapsulation**: Respect package boundaries

### Contributing Guidelines

#### Adding New Features

1. **Identify package**: Determine appropriate package for new functionality
2. **Check interfaces**: Ensure feature fits existing interfaces
3. **Create focused files**: Keep new files small and focused
4. **Add comprehensive tests**: Include unit and integration tests
5. **Update documentation**: Document new interfaces and functionality

#### Modifying Existing Features

1. **Locate components**: Use package structure to find relevant code
2. **Understand dependencies**: Review interfaces and dependencies
3. **Test changes**: Ensure modifications don't break existing functionality
4. **Update tests**: Modify tests to reflect changes
5. **Validate integration**: Run full test suite to ensure system coherence

---

This migration guide provides comprehensive information for working with the refactored codebase. The modular structure improves maintainability, testability, and developer experience while preserving all existing functionality.
