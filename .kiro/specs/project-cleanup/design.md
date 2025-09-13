# Design Document

## Overview

This design outlines a focused cleanup strategy for the Open Source Template Generator Go project to streamline it to only the core generator functionality. The cleanup will systematically remove all unnecessary commands, packages, configurations, documentation, and scripts while preserving the essential template generation capabilities.

The project currently has multiple commands in the cmd directory, but only the generator command is essential. The cleanup will remove all other commands and their supporting infrastructure, resulting in a simplified, focused codebase that maintains the core template generation functionality.

## Architecture

### Cleanup Strategy Architecture

The cleanup process follows a systematic removal approach:

1. **Analysis Layer**: Identify dependencies between components and the generator
2. **Impact Assessment Layer**: Determine which components can be safely removed
3. **Removal Layer**: Systematically remove unnecessary components
4. **Validation Layer**: Ensure generator functionality remains intact
5. **Cleanup Layer**: Remove orphaned dependencies and update configurations

### Component Analysis

Based on the codebase examination, the following components will be evaluated:

#### Commands to Remove
- `cmd/cleanup/` - Project cleanup utilities (not core generator function)
- `cmd/consolidate-duplicates/` - Code consolidation tools (not core generator function)
- `cmd/duplicate-scanner/` - Duplicate detection tools (not core generator function)
- `cmd/import-analyzer/` - Import analysis tools (not core generator function)
- `cmd/import-detector/` - Import detection tools (not core generator function)
- `cmd/remove-unused-code/` - Code cleanup tools (not core generator function)
- `cmd/security-fixer/` - Security fixing tools (not core generator function)
- `cmd/security-linter/` - Security linting tools (not core generator function)
- `cmd/security-scanner/` - Security scanning tools (not core generator function)
- `cmd/standards/` - Standards enforcement tools (not core generator function)
- `cmd/todo-resolver/` - TODO resolution tools (not core generator function)
- `cmd/todo-scanner/` - TODO scanning tools (not core generator function)
- `cmd/unused-code-scanner/` - Unused code detection tools (not core generator function)

#### Commands to Keep
- `cmd/generator/` - Core template generation functionality (ESSENTIAL)

#### Supporting Components to Evaluate
- `internal/` packages - Keep only those essential for generator
- `pkg/` packages - Keep only those essential for template generation
- `config/` files - Keep only generator-related configurations
- `docs/` files - Keep only generator-related documentation
- `scripts/` files - Keep only generator build/development scripts

## Components and Interfaces

### 1. Dependency Analyzer

```go
type DependencyAnalyzer interface {
    AnalyzeCommandDependencies(cmdPath string) (*DependencyReport, error)
    FindPackageUsage(pkgPath string) ([]string, error)
    IdentifyOrphanedComponents() ([]string, error)
    ValidateGeneratorDependencies() error
}
```

**Implementation Strategy:**

- Parse Go files to identify import relationships
- Build dependency graph to understand component relationships
- Identify which internal and pkg components are used by generator
- Find components that become orphaned after command removal

### 2. Component Remover

```go
type ComponentRemover interface {
    RemoveCommand(cmdPath string) error
    RemovePackage(pkgPath string) error
    RemoveConfigFiles(patterns []string) error
    RemoveDocumentation(patterns []string) error
    RemoveScripts(patterns []string) error
}
```

**Implementation Strategy:**

- Safely remove directories and files
- Update import statements in remaining files
- Clean up references in configuration files
- Update build scripts and documentation

### 3. Generator Validator

```go
type GeneratorValidator interface {
    ValidateGeneratorFunctionality() error
    TestTemplateGeneration() error
    VerifyBuildProcess() error
    CheckRequiredDependencies() error
}
```

**Implementation Strategy:**

- Run generator with test inputs to ensure functionality
- Verify all template types can still be generated
- Ensure build process works for simplified project
- Validate that all required dependencies are still present

### 4. Configuration Updater

```go
type ConfigurationUpdater interface {
    UpdateBuildConfigs() error
    UpdateDocumentation() error
    CleanupGoMod() error
    UpdateCIConfigs() error
}
```

**Implementation Strategy:**

- Update Makefile and build scripts to remove references to deleted commands
- Update README and documentation to reflect simplified project
- Run go mod tidy to remove unused dependencies
- Update CI/CD configurations if they exist

## Data Models

### Component Analysis Models

```go
type DependencyReport struct {
    Component        string
    DirectDependents []string
    IndirectDependents []string
    IsEssential      bool
    RemovalImpact    string
}

type ComponentInfo struct {
    Path            string
    Type            string // command, package, config, doc, script
    IsEssential     bool
    Dependencies    []string
    Dependents      []string
    RemovalSafe     bool
}

type GeneratorValidation struct {
    FunctionalityIntact bool
    BuildSuccessful     bool
    TestsPassing        bool
    TemplatesWorking    []string
    Issues              []string
}
```

### Cleanup Configuration Models

```go
type CleanupConfig struct {
    PreservePatterns []string
    BackupEnabled    bool
    DryRun          bool
    Verbose         bool
    ValidateAfter   bool
}

type CleanupResult struct {
    RemovedCommands     []string
    RemovedPackages     []string
    RemovedConfigs      []string
    RemovedDocs         []string
    RemovedScripts      []string
    UpdatedFiles        []string
    ValidationResults   *GeneratorValidation
    Summary             *CleanupSummary
}

type CleanupSummary struct {
    TotalFilesRemoved   int
    TotalFilesModified  int
    SpaceSaved          int64
    GeneratorStatus     string
    Issues              []string
}
```

## Error Handling

### Cleanup Error Management

```go
type CleanupError struct {
    Type        ErrorType
    Component   string
    Operation   string
    Message     string
    Suggestion  string
    Recoverable bool
}

type ErrorType string

const (
    ErrorTypeDependency   ErrorType = "dependency"
    ErrorTypeRemoval      ErrorType = "removal"
    ErrorTypeValidation   ErrorType = "validation"
    ErrorTypeConfiguration ErrorType = "configuration"
)
```

**Error Handling Strategy:**

- Validate dependencies before removal to prevent breaking generator
- Create backups before any destructive operations
- Provide clear error messages with recovery suggestions
- Implement rollback mechanisms for failed operations

### Safety Mechanisms

- Dependency validation before component removal
- Automatic backup creation before modifications
- Generator functionality validation after each major change
- Rollback capability if generator functionality is compromised

## Testing Strategy

### Generator Functionality Validation

1. **Pre-Cleanup Validation**
   - Document current generator functionality
   - Test all template types and generation scenarios
   - Record baseline performance metrics
   - Capture current build and test results

2. **Incremental Validation**
   - Test generator after each component removal
   - Validate that all template types still work
   - Ensure build process remains functional
   - Check that essential dependencies are preserved

3. **Post-Cleanup Validation**
   - Comprehensive generator functionality testing
   - Performance comparison with baseline
   - Integration testing with all template types
   - Documentation accuracy validation

### Validation Test Suite

```go
type GeneratorTestSuite struct {
    TemplateTests     []TemplateTest
    BuildTests        []BuildTest
    IntegrationTests  []IntegrationTest
    PerformanceTests  []PerformanceTest
}

type TemplateTest struct {
    TemplateType    string
    InputParams     map[string]interface{}
    ExpectedOutput  string
    ValidationFunc  func(output string) error
}
```

**Testing Approach:**

- Automated testing after each component removal
- Rollback if generator functionality is compromised
- Comprehensive validation before finalizing cleanup
- Performance regression detection

## Implementation Phases

### Phase 1: Analysis and Dependency Mapping

1. Analyze all commands and their dependencies
2. Map internal and pkg package usage by generator
3. Identify configuration files used by generator
4. Document current generator functionality baseline

### Phase 2: Command Removal

1. Remove non-essential commands from cmd directory
2. Validate generator still builds and runs after each removal
3. Update any references to removed commands in documentation
4. Clean up command-specific configurations

### Phase 3: Package Cleanup

1. Remove internal packages not used by generator
2. Remove pkg packages not essential for template generation
3. Update import statements in remaining files
4. Validate generator functionality after package removals

### Phase 4: Configuration and Documentation Cleanup

1. Remove configuration files not needed by generator
2. Remove documentation for removed features
3. Update main README to reflect simplified project
4. Clean up scripts not needed for generator development

### Phase 5: Dependency and Build Cleanup

1. Run go mod tidy to remove unused dependencies
2. Update build scripts and Makefile
3. Update CI/CD configurations if present
4. Remove orphaned test files

### Phase 6: Final Validation and Documentation

1. Comprehensive generator functionality testing
2. Performance validation
3. Update project documentation
4. Generate cleanup report

## Security Considerations

### Generator Security Preservation

Since we're removing security-related commands and packages, we need to ensure:

1. **Template Security**
   - Preserve any security features in generated templates
   - Ensure template generation doesn't introduce vulnerabilities
   - Maintain secure defaults in template configurations

2. **Generator Security**
   - Validate that generator itself remains secure
   - Ensure file operations are safe
   - Preserve input validation for template parameters

**Security Impact Assessment:**

- Most security packages (security-scanner, security-linter, etc.) are tools, not core generator functionality
- Template security features should be preserved in template files
- Generator core should maintain basic security practices

### Security Validation After Cleanup

- Verify generated templates maintain security configurations
- Ensure generator doesn't create insecure file permissions
- Validate that template parameters are properly sanitized

## Performance Impact

### Expected Performance Improvements

1. **Reduced Binary Size**
   - Removing unused commands will reduce compiled binary size
   - Fewer dependencies will reduce build time
   - Simplified codebase will improve maintainability

2. **Faster Build Times**
   - Fewer packages to compile
   - Reduced dependency graph complexity
   - Simplified test suite execution

3. **Generator Performance**
   - Core generator performance should remain unchanged
   - May see slight improvements due to reduced memory footprint
   - Template processing speed should be unaffected

### Performance Monitoring

```go
type PerformanceMetrics struct {
    BinarySize        int64
    BuildTime         time.Duration
    GenerationTime    time.Duration
    MemoryUsage      int64
    DependencyCount   int
}
```

## Cleanup Monitoring and Logging

### Cleanup Progress Tracking

1. **Component Removal Logging**
   - Log each component being analyzed for removal
   - Track dependencies and impact assessment
   - Record removal decisions and rationale

2. **Validation Logging**
   - Log generator functionality tests after each change
   - Track build success/failure after modifications
   - Record any issues discovered during cleanup

3. **Summary Reporting**
   - Generate detailed cleanup report
   - Document space savings and simplification achieved
   - Provide before/after comparison metrics

## Backup and Recovery Strategy

### Backup Strategy

```go
type BackupManager interface {
    CreateFullBackup() (*Backup, error)
    CreateIncrementalBackup(components []string) (*Backup, error)
    RestoreBackup(backup *Backup) error
    ValidateBackup(backup *Backup) error
}
```

### Recovery Mechanisms

1. **Pre-Cleanup Backup**
   - Create complete project backup before starting
   - Include git state and all files
   - Validate backup integrity

2. **Incremental Safety**
   - Test generator after each major component removal
   - Automatic rollback if generator functionality breaks
   - Preserve working state at each validation checkpoint

3. **Recovery Validation**
   - Verify generator works correctly after any rollback
   - Ensure all template types still generate properly
   - Validate build process remains functional
