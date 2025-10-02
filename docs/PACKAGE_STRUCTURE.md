# Package Structure Documentation

## Overview

This document describes the refactored package structure of the Open Source Project Generator following the code splitting initiative. The codebase has been modularized to improve maintainability, testability, and extensibility while preserving all existing functionality.

## Package Organization

### Core Architecture

The project follows a layered architecture with clear separation of concerns:

```text
pkg/
├── interfaces/        # Core interfaces and contracts
├── models/           # Data structures and configuration models
├── cli/              # Command-line interface (modularized)
├── template/         # Template processing system (modularized)
├── audit/            # Security and quality auditing (modularized)
├── validation/       # Project validation engine (modularized)
├── filesystem/       # File system operations (modularized)
├── cache/            # Caching system (modularized)
├── version/          # Version management
├── security/         # Security utilities
├── ui/               # Interactive user interface
├── errors/           # Error handling and categorization
├── utils/            # Utility functions
└── constants/        # Application constants
```

## Package Details

### CLI Package (`pkg/cli/`)

**Purpose**: Modularized command-line interface implementation

**Structure**:
```text
pkg/cli/
├── cli.go              # Main CLI struct and coordination (~500 lines)
├── commands.go         # Command registration and setup
├── handlers.go         # Command execution handlers
├── output.go           # Output formatting and colors
├── flags.go            # Flag management and validation
├── interactive.go      # Interactive mode management
├── validation.go       # CLI-specific validation
├── commands/           # Command-specific implementations
│   ├── generate.go     # Generate command
│   ├── validate.go     # Validate command
│   ├── audit.go        # Audit command
│   └── template.go     # Template management commands
├── handlers/           # Modular command handlers
│   ├── generate_handler.go  # Project generation workflow
│   └── workflow_handler.go  # Complete workflow management
├── interactive/        # Interactive UI components
│   ├── project_setup.go     # Project configuration UI
│   └── component_selection.go # Component selection UI
├── utils/              # CLI utilities
│   ├── formatters.go   # Output formatting utilities
│   └── helpers.go      # Common helper functions
└── validation/         # CLI-specific validation
    └── input_validator.go   # Input validation logic
```

**Key Components**:
- **CLI**: Main coordinator struct that orchestrates all CLI functionality
- **OutputManager**: Handles all output formatting, colors, and verbosity levels
- **FlagHandler**: Manages command-line flags and validation
- **CommandRegistry**: Registers and sets up all available commands
- **GenerateHandler**: Manages project generation workflows and component processing
- **WorkflowHandler**: Orchestrates complete generation workflows including validation and post-processing
- **InteractiveManager**: Manages interactive vs non-interactive modes
- **CLIValidator**: Validates CLI-specific configurations and inputs

### Audit Package (`pkg/audit/`)

**Purpose**: Modularized security and quality auditing system

**Structure**:
```text
pkg/audit/
├── engine.go           # Main audit orchestration (~300 lines)
├── rules.go            # Audit rule management
├── result.go           # Result aggregation and reporting
├── security/           # Security audit modules
│   ├── scanner.go      # Security scanning coordination
│   ├── secrets.go      # Secret detection
│   └── dependencies.go # Dependency vulnerability checking
├── quality/            # Code quality modules
│   ├── complexity.go   # Complexity analysis
│   └── coverage.go     # Coverage measurement
├── license/            # License compliance
│   └── compatibility.go # License compatibility checking
└── performance/        # Performance analysis
    ├── bundle.go       # Bundle size analysis
    └── metrics.go      # Performance metrics
```

**Key Components**:
- **Engine**: Orchestrates all audit types and generates comprehensive reports
- **RuleManager**: Manages audit rules, filtering, and validation
- **ResultAggregator**: Aggregates results from different audit types
- **SecurityScanner**: Coordinates security-related audits
- **QualityAnalyzer**: Handles code quality measurements
- **LicenseChecker**: Validates license compatibility
- **PerformanceAnalyzer**: Measures performance metrics

### Template Package (`pkg/template/`)

**Purpose**: Modularized template processing and management system

**Structure**:
```text
pkg/template/
├── manager.go          # Template coordination (~400 lines)
├── discovery.go        # Template discovery and scanning
├── cache.go            # Template caching system
├── validation.go       # Template validation
├── processor/          # Template processing engine
│   └── engine.go       # Template compilation and execution
├── metadata/           # Template metadata handling
│   ├── parser.go       # Metadata parsing
│   └── validator.go    # Metadata validation
└── templates/          # Template files (unchanged)
    ├── base/
    ├── frontend/
    ├── backend/
    ├── mobile/
    └── infrastructure/
```

**Key Components**:
- **Manager**: Coordinates all template operations
- **TemplateDiscovery**: Discovers and scans available templates
- **TemplateCache**: Manages template caching and invalidation
- **TemplateValidator**: Validates template structure and syntax
- **ProcessingEngine**: Compiles and executes templates
- **MetadataParser**: Parses template metadata files
- **MetadataValidator**: Validates template metadata

### Validation Package (`pkg/validation/`)

**Purpose**: Modularized project validation system

**Structure**:
```text
pkg/validation/
├── engine.go           # Validation orchestration
├── config_validator.go # Configuration validation coordination
├── schemas.go          # Schema management
├── formats/            # Format-specific validators
│   ├── json.go         # JSON validation
│   ├── yaml.go         # YAML validation
│   ├── env.go          # Environment file validation
│   ├── docker.go       # Dockerfile validation
│   └── makefile.go     # Makefile validation
└── [other files...]    # Additional validation components
```

**Key Components**:
- **Engine**: Orchestrates all validation types
- **ConfigValidator**: Coordinates configuration validation
- **SchemaManager**: Manages validation schemas and rules
- **Format Validators**: Handle specific file format validation

### Filesystem Package (`pkg/filesystem/`)

**Purpose**: Modularized file system operations and project generation

**Structure**:
```text
pkg/filesystem/
├── project_generator.go # Main project generation coordination (~300 lines)
├── structure.go        # Project structure management
├── operations.go       # File system operations
├── components/         # Component-specific generators
│   ├── frontend.go     # Frontend component generation
│   ├── backend.go      # Backend component generation
│   ├── mobile.go       # Mobile component generation
│   └── infrastructure.go # Infrastructure component generation (~400 lines)
├── generators/         # Specialized generators
│   ├── structure.go    # Directory structure generation
│   ├── templates.go    # Template processing
│   ├── configuration.go # Configuration file generation
│   ├── documentation.go # Documentation generation
│   └── cicd.go         # CI/CD pipeline generation
├── processors/         # File processing engines
│   ├── template_processor.go # Template processing logic
│   ├── file_processor.go     # File operation logic
│   └── interfaces.go         # Processor interfaces
└── operations/         # File operation implementations
    ├── creator.go      # File creation operations
    ├── validator.go    # File validation operations
    └── interfaces.go   # Operation interfaces
```

**Key Components**:
- **ProjectGenerator**: Coordinates project generation workflow and component orchestration
- **StructureManager**: Manages project structure definitions and directory layouts
- **ComponentGenerators**: Generate component-specific files (frontend, backend, mobile, infrastructure)
- **TemplateProcessor**: Handles template compilation and variable substitution
- **FileProcessor**: Manages safe file operations with validation and error handling
- **OperationValidator**: Validates file operations and ensures data integrity

### UI Package (`pkg/ui/`)

**Purpose**: Modularized interactive user interface components

**Structure**:
```text
pkg/ui/
├── interactive.go          # Main interactive UI coordination
├── config_manager.go       # Configuration management (~300 lines)
├── project_structure_preview.go # Project preview (~300 lines)
├── template_selector.go    # Template selection interface
├── validation.go           # UI input validation
├── preview/                # Preview system components
│   ├── types.go           # Preview type definitions
│   ├── tree/              # Tree visualization
│   │   ├── builder.go     # Tree structure building
│   │   ├── renderer.go    # Tree rendering
│   │   └── formatter.go   # Tree formatting
│   ├── components/        # Component previews
│   │   ├── frontend.go    # Frontend preview
│   │   ├── backend.go     # Backend preview
│   │   └── infrastructure.go # Infrastructure preview
│   └── display/           # Display logic
│       ├── console.go     # Console display
│       └── interactive.go # Interactive display
└── config/                # Configuration UI components
    ├── collectors/        # Input collectors
    │   ├── project_info.go # Project information collection
    │   ├── components.go   # Component selection
    │   └── advanced.go     # Advanced options
    ├── validators/        # Input validators
    │   ├── project_validator.go # Project validation
    │   └── component_validator.go # Component validation
    └── formatters/        # Output formatters
        ├── summary.go     # Configuration summary
        └── export.go      # Configuration export
```

**Key Components**:
- **InteractiveUI**: Main coordinator for user interactions and workflow management
- **ConfigManager**: Manages configuration collection, validation, and formatting
- **PreviewManager**: Handles project structure previews and component visualization
- **TemplateSelector**: Provides template selection and filtering interface
- **TreeBuilder**: Constructs hierarchical project structure representations
- **ComponentCollectors**: Gather user input for different project components
- **ValidationEngine**: Validates user input and provides feedback

### Cache Package (`pkg/cache/`)

**Purpose**: Modularized caching system

**Structure**:
```text
pkg/cache/
├── manager.go          # Cache coordination (~300 lines)
├── storage.go          # Cache storage operations
├── operations.go       # Cache operations coordination
├── validation.go       # Cache validation
├── cleanup.go          # Cache cleanup
├── operations/         # Modular cache operations
│   ├── operations.go   # Operation coordination and callbacks
│   ├── get.go          # Get operations with hit/miss tracking
│   ├── set.go          # Set operations with eviction policies
│   ├── delete.go       # Delete operations and cache clearing
│   └── cleanup.go      # Cleanup and compaction operations
├── storage/            # Storage backend implementations
│   ├── interface.go    # Storage interface definitions
│   ├── filesystem.go   # File system storage backend
│   └── memory.go       # In-memory storage backend
├── metrics/            # Cache metrics and reporting
│   ├── collector.go    # Metrics collection
│   └── reporter.go     # Metrics reporting
└── validation/         # Cache validation
    └── validator.go    # Cache integrity validation
```

**Key Components**:
- **Manager**: Coordinates all cache operations and manages storage backends
- **CacheOperations**: Orchestrates get, set, delete, and cleanup operations with callbacks
- **StorageBackends**: Pluggable storage implementations (filesystem, memory)
- **MetricsCollector**: Tracks cache performance and usage statistics
- **CacheValidator**: Validates cache integrity and handles corruption recovery
- **CleanupManager**: Handles cache cleanup, compaction, and expired entry removal

## Interface-Based Design

### Core Interfaces (`pkg/interfaces/`)

All major components implement well-defined interfaces to enable:
- **Dependency Injection**: Components receive dependencies through interfaces
- **Testing**: Easy mocking and unit testing
- **Extensibility**: New implementations can be added without changing existing code
- **Modularity**: Clear contracts between components

Key interface categories:
- **CLI Interfaces**: Command handling, output management, interaction
- **Template Interfaces**: Template processing, discovery, validation
- **Audit Interfaces**: Security scanning, quality analysis, reporting
- **Validation Interfaces**: Configuration validation, format checking
- **Filesystem Interfaces**: File operations, project generation
- **Cache Interfaces**: Caching operations, storage management

## Data Models (`pkg/models/`)

Centralized data structures used across the application:
- **Configuration Models**: Project configuration, component settings
- **Template Models**: Template metadata, processing context
- **Audit Models**: Audit results, rules, reports
- **Validation Models**: Validation results, error reporting
- **Error Models**: Structured error handling

## Benefits of the New Structure

### Maintainability
- **Single Responsibility**: Each file has a clear, focused purpose
- **Smaller Files**: No file exceeds 1,000 lines (target achieved for most packages)
- **Clear Dependencies**: Interface-based design makes dependencies explicit
- **Modular Testing**: Components can be tested in isolation

### Extensibility
- **Plugin Architecture**: New audit types, validators, or generators can be added easily
- **Interface Compliance**: New implementations just need to satisfy interfaces
- **Component Isolation**: Changes to one component don't affect others

### Developer Experience
- **Easier Navigation**: Developers can quickly find relevant code
- **Parallel Development**: Multiple developers can work on different components
- **Clear Ownership**: Each package has a clear purpose and maintainer
- **Better IDE Support**: Smaller files improve IDE performance and navigation

## Migration from Old Structure

### Major Changes

1. **CLI Package Split**:
   - `pkg/cli/cli.go` (5,450 lines) → Multiple focused files (~500 lines each)
   - Command-specific logic moved to `pkg/cli/commands/`
   - Output and formatting logic moved to `pkg/cli/output.go`

2. **Audit Engine Split**:
   - `pkg/audit/engine.go` (3,029 lines) → Orchestration (~300 lines) + specialized modules
   - Security audits moved to `pkg/audit/security/`
   - Quality audits moved to `pkg/audit/quality/`

3. **Template Manager Split**:
   - `pkg/template/manager.go` (1,548 lines) → Coordination (~400 lines) + specialized modules
   - Processing logic moved to `pkg/template/processor/`
   - Metadata handling moved to `pkg/template/metadata/`

4. **Validation Split**:
   - Format-specific validation moved to `pkg/validation/formats/`
   - Schema management centralized in `pkg/validation/schemas.go`

5. **Filesystem Split**:
   - Component-specific generation moved to `pkg/filesystem/components/`
   - File operations isolated in `pkg/filesystem/operations.go`

### Backward Compatibility

- **Public APIs**: All public interfaces remain unchanged
- **CLI Commands**: All commands work identically
- **Configuration**: All configuration options preserved
- **Templates**: Template processing unchanged
- **Functionality**: 100% feature parity maintained

## Development Guidelines

### Adding New Features

1. **Identify the Package**: Determine which package the feature belongs to
2. **Check Interfaces**: Ensure the feature fits existing interfaces or create new ones
3. **Follow Patterns**: Use existing patterns for similar functionality
4. **Add Tests**: Include comprehensive tests for new components
5. **Update Documentation**: Document new functionality and interfaces

### Modifying Existing Features

1. **Locate Components**: Use the package structure to find relevant code
2. **Check Dependencies**: Understand component dependencies through interfaces
3. **Test Changes**: Ensure changes don't break existing functionality
4. **Update Tests**: Modify tests to reflect changes
5. **Validate Integration**: Run integration tests to ensure system coherence

### Best Practices

- **Keep Files Small**: Target maximum 1,000 lines per file
- **Single Responsibility**: Each file should have one clear purpose
- **Interface First**: Design interfaces before implementations
- **Test Coverage**: Maintain high test coverage for all components
- **Documentation**: Keep documentation up-to-date with changes

## Performance Considerations

### Benefits of Modularization

- **Faster Compilation**: Smaller files compile faster
- **Better Caching**: Go's build cache works better with smaller files
- **Parallel Processing**: Multiple files can be processed in parallel
- **Reduced Memory**: Smaller working sets during development

### Potential Concerns

- **Import Overhead**: More files mean more imports (minimal impact)
- **Interface Calls**: Interface calls have slight overhead (negligible in practice)
- **Package Initialization**: More packages mean more initialization (minimal impact)

## Testing Strategy

### Unit Testing
- Each component has focused unit tests
- Interfaces enable easy mocking
- Smaller files make tests more targeted

### Integration Testing
- Tests verify components work together correctly
- Interface contracts are validated
- End-to-end workflows are tested

### Performance Testing
- Benchmarks ensure no performance regressions
- Memory usage is monitored
- Build times are tracked

## Future Improvements

### Potential Enhancements

1. **Plugin System**: Further modularize with plugin architecture
2. **Microservices**: Split into separate services if needed
3. **API Extraction**: Extract core functionality to reusable APIs
4. **Configuration DSL**: Create domain-specific language for configuration

### Monitoring

- **File Size Monitoring**: Automated checks for file size limits
- **Dependency Analysis**: Regular analysis of component dependencies
- **Performance Tracking**: Continuous monitoring of performance metrics
- **Test Coverage**: Automated coverage reporting and enforcement

---

This modular structure provides a solid foundation for future development while maintaining all existing functionality and improving the developer experience significantly.