# Architecture

System architecture and design of the Open Source Project Generator.

## Table of Contents

- [Overview](#overview)
- [High-Level Architecture](#high-level-architecture)
- [Core Components](#core-components)
- [Generation Workflow](#generation-workflow)
- [Design Decisions](#design-decisions)
- [Extension Points](#extension-points)

---

## Overview

The Open Source Project Generator uses a **tool-orchestration architecture** that delegates project creation to industry-standard CLI tools rather than maintaining templates manually.

### Key Principles

1. **Delegate to Experts** - Use official framework CLIs (`create-next-app`, `go mod init`, etc.)
2. **Graceful Degradation** - Provide fallback generators when tools unavailable
3. **Standardized Output** - Map all outputs to consistent directory structure
4. **Offline Support** - Cache tool availability for offline operation
5. **Security First** - Sanitize all inputs, validate all operations

---

## High-Level Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (cobra commands, flags, user interaction)                  │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                   Orchestration Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Project    │  │    Tool      │  │  Integration │     │
│  │  Coordinator │  │  Discovery   │  │   Manager    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                   Generation Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Bootstrap   │  │   Fallback   │  │   Structure  │     │
│  │   Executor   │  │  Generator   │  │   Mapper     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                  Infrastructure Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  File System │  │   Security   │  │    Logger    │     │
│  │  Operations  │  │  Validator   │  │              │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## Core Components

### 1. CLI Layer

**Location:** `cmd/generator/main.go`

**Responsibilities:**

- Parse command-line arguments
- Handle user interaction
- Display results and errors
- Manage process lifecycle

**Key Commands:**

- `generate` - Generate projects
- `check-tools` - Validate tool availability
- `init-config` - Create configuration templates
- `cache-tools` - Manage tool cache
- `version` - Show version information

### 2. Project Coordinator

**Location:** `internal/orchestrator/coordinator.go`

**Responsibilities:**

- Orchestrate complete generation workflow
- Coordinate between all components
- Handle errors and rollback
- Manage generation state

**Key Methods:**

```go
func (pc *ProjectCoordinator) Generate(ctx context.Context, config *models.ProjectConfig) (*models.GenerationResult, error)
func (pc *ProjectCoordinator) DryRun(ctx context.Context, config *models.ProjectConfig) (*models.PreviewResult, error)
func (pc *ProjectCoordinator) Validate(config *models.ProjectConfig) error
```

**Generation Steps:**

1. Validate configuration
2. Apply defaults
3. Sanitize inputs
4. Discover available tools
5. Create backup (if requested)
6. Prepare output directory
7. Generate each component
8. Map to target structure
9. Integrate components
10. Validate generated project
11. Return results

### 3. Tool Discovery

**Location:** `internal/orchestrator/tool_discovery.go`

**Responsibilities:**

- Detect available bootstrap tools
- Verify version requirements
- Provide installation instructions
- Cache tool availability

**Tool Registry:**

```go
type ToolMetadata struct {
    Name              string
    Command           string
    VersionFlag       string
    MinVersion        string
    ComponentTypes    []string
    InstallDocs       map[string]string
    FallbackAvailable bool
}
```

**Registered Tools:**

- `npx` - For Next.js projects
- `go` - For Go backend projects
- `gradle` - For Android projects
- `xcodebuild` - For iOS projects

**Key Methods:**

```go
func (td *ToolDiscovery) IsAvailable(toolName string) (bool, error)
func (td *ToolDiscovery) GetVersion(toolName string) (string, error)
func (td *ToolDiscovery) CheckRequirements(tools []string) (*models.ToolCheckResult, error)
```

### 4. Bootstrap Executors

**Location:** `internal/generator/bootstrap/`

**Responsibilities:**

- Execute external CLI tools
- Build commands with appropriate flags
- Capture output
- Handle errors and timeouts

**Base Executor:**

```go
type BaseExecutor struct {
    logger *logger.Logger
}

func (be *BaseExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error)
```

**Component-Specific Executors:**

- `NextJSExecutor` - Runs `npx create-next-app`
- `GoExecutor` - Runs `go mod init` and sets up Gin
- `AndroidExecutor` - Runs Gradle commands
- `iOSExecutor` - Runs Xcode commands

### 5. Fallback Generators

**Location:** `internal/generator/fallback/`

**Responsibilities:**

- Generate minimal project structures
- Provide setup instructions
- Support offline generation
- Work when tools unavailable

**Fallback Registry:**

```go
type Registry struct {
    generators map[string]Generator
}

func (r *Registry) Get(componentType string) (Generator, bool)
```

**Available Fallbacks:**

- Android - Minimal Gradle project structure
- iOS - Minimal Xcode project structure

### 6. Structure Mapper

**Location:** `internal/generator/mapper/`

**Responsibilities:**

- Map generated files to target structure
- Update import paths
- Create symlinks
- Validate final structure

**Component Mappings:**

```go
var ComponentMappings = map[string]string{
    "nextjs":      "App/",
    "go-backend":  "CommonServer/",
    "android":     "Mobile/android/",
    "ios":         "Mobile/ios/",
}
```

**Key Methods:**

```go
func (sm *StructureMapper) Map(ctx context.Context, source, target, componentType string) error
func (sm *StructureMapper) UpdateImportPaths(targetPath, componentType string) error
```

### 7. Integration Manager

**Location:** `internal/orchestrator/integration.go`

**Responsibilities:**

- Generate Docker Compose files
- Create build scripts
- Configure environment variables
- Set up component communication

**Key Methods:**

```go
func (im *IntegrationManager) Integrate(ctx context.Context, results []*models.ComponentResult) error
func (im *IntegrationManager) GenerateDockerCompose(components []*models.Component) (string, error)
func (im *IntegrationManager) GenerateScripts(projectRoot string) error
```

### 8. Offline Detector

**Location:** `internal/orchestrator/offline.go`

**Responsibilities:**

- Detect network connectivity
- Manage offline mode
- Provide offline status messages
- Force offline mode when requested

**Key Methods:**

```go
func (od *OfflineDetector) IsOffline() bool
func (od *OfflineDetector) ForceOffline(offline bool)
func (od *OfflineDetector) GetOfflineMessage() string
```

### 9. Tool Cache

**Location:** `internal/orchestrator/tool_cache.go`

**Responsibilities:**

- Cache tool availability
- Store version information
- Manage cache expiration
- Support offline operation

**Cache Structure:**

```go
type CachedTool struct {
    Available        bool
    Version          string
    Path             string
    LastChecked      time.Time
}
```

### 10. Rollback Manager

**Location:** `internal/orchestrator/rollback.go`

**Responsibilities:**

- Track generation operations
- Restore backups on failure
- Clean up partial generations
- Provide rollback points

**Key Methods:**

```go
func (rm *RollbackManager) CreateCheckpoint(name string) error
func (rm *RollbackManager) Rollback() error
func (rm *RollbackManager) Commit() error
```

---

## Generation Workflow

### Complete Generation Flow

```text
1. User Input
   ↓
2. Configuration Loading/Validation
   ↓
3. Tool Discovery
   ↓
4. Backup Creation (optional)
   ↓
5. Component Generation
   │
   ├─→ Bootstrap Tool Available?
   │   ├─→ Yes: Execute Bootstrap Tool
   │   └─→ No: Use Fallback Generator
   │
   ↓
6. Structure Mapping
   ↓
7. Integration
   │
   ├─→ Docker Compose Generation
   ├─→ Script Generation
   └─→ Environment Configuration
   ↓
8. Validation
   ↓
9. Result Return
```

### Detailed Component Generation

```text
For each component:

1. Check Tool Availability
   ├─→ Tool Available
   │   ├─→ Build Command
   │   ├─→ Execute Tool
   │   ├─→ Capture Output
   │   └─→ Handle Errors
   │
   └─→ Tool Unavailable
       ├─→ Check Fallback Available
       ├─→ Generate Minimal Structure
       └─→ Provide Setup Instructions

2. Map to Target Structure
   ├─→ Move Files
   ├─→ Update Imports
   └─→ Create Symlinks

3. Record Result
   ├─→ Success/Failure
   ├─→ Duration
   ├─→ Warnings
   └─→ Manual Steps
```

---

## Design Decisions

### 1. Tool Orchestration vs Templates

**Decision:** Delegate to external CLI tools instead of maintaining templates.

**Rationale:**

- Framework authors maintain their own CLIs
- Always get latest versions without manual updates
- Reduces maintenance burden
- Leverages community expertise

**Trade-offs:**

- Requires external tools to be installed
- Less control over generated structure
- Need fallback generators for unavailable tools

### 2. Standardized Directory Structure

**Decision:** Map all generated outputs to consistent structure.

**Rationale:**

- Consistent project layout across all projects
- Easy navigation and understanding
- Simplifies integration between components
- Familiar structure for team members

**Implementation:**

- Structure Mapper relocates files after generation
- Updates import paths and references
- Creates symlinks where appropriate

### 3. Graceful Degradation with Fallback Generators

**Decision:** Provide fallback generators when bootstrap tools unavailable.

**Rationale:**

- Ensures all component types can be generated
- Works in restricted environments
- Supports offline development
- Provides minimal viable structure

**Implementation:**

- Fallback generators create minimal project structure
- Include detailed setup instructions
- Use embedded templates (Go embed)
- Clearly indicate manual steps required

### 4. Offline Support

**Decision:** Support offline project generation after initial setup.

**Rationale:**

- Enable development in restricted networks
- Improve reliability
- Faster generation with cached tools
- Support air-gapped environments

**Implementation:**

- Tool caching system
- Fallback to cached tools
- Offline mode detection
- Clear offline indicators

### 5. Security-First Approach

**Decision:** Validate and sanitize all inputs before execution.

**Rationale:**

- Prevent command injection
- Avoid path traversal attacks
- Whitelist allowed tools and flags
- Protect user systems

**Implementation:**

- Input sanitization via `pkg/security`
- Tool whitelisting in tool registry
- Flag validation before execution
- Timeout and resource limits

---

## Extension Points

### Adding New Bootstrap Tools

1. **Define Tool Metadata** in `tool_discovery.go`:

   ```go
   "your-tool": {
       Name:              "your-tool",
       Command:           "your-tool",
       VersionFlag:       "--version",
       MinVersion:        "1.0.0",
       ComponentTypes:    []string{"your-component"},
       InstallDocs:       map[string]string{...},
       FallbackAvailable: true,
   }
   ```

2. **Create Executor** in `internal/generator/bootstrap/`:

   ```go
   type YourToolExecutor struct {
       *BaseExecutor
   }

   func (yte *YourToolExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
       // Implementation
   }
   ```

3. **Register Executor** in coordinator

4. **Add Structure Mapping** in `mapper/structure.go`

5. **Create Fallback Generator** (optional) in `internal/generator/fallback/`

### Adding New Component Types

1. Define component type in configuration schema
2. Create executor or fallback generator
3. Add structure mapping
4. Update documentation

### Adding New Integration Features

1. Extend Integration Manager
2. Add configuration options
3. Implement generation logic
4. Update documentation

---

## Data Flow

### Configuration Flow

```text
User Input (YAML/JSON)
  ↓
Configuration Parser
  ↓
Validator (apply defaults, validate)
  ↓
Sanitizer (security checks)
  ↓
Project Coordinator
```

### Tool Discovery Flow

```text
Tool Registry
  ↓
Tool Discovery
  ├─→ Check PATH
  ├─→ Verify Version
  └─→ Cache Result
  ↓
Tool Check Result
```

### Generation Flow

```text
Project Config
  ↓
Component Loop
  ├─→ Tool Discovery
  ├─→ Bootstrap Executor OR Fallback Generator
  ├─→ Structure Mapper
  └─→ Component Result
  ↓
Integration Manager
  ├─→ Docker Compose
  ├─→ Scripts
  └─→ Environment
  ↓
Generation Result
```

---

## Error Handling

### Error Categories

1. **Configuration Errors** - Invalid configuration
2. **Tool Errors** - Tool not found or execution failed
3. **File System Errors** - Permission denied, disk full
4. **Validation Errors** - Invalid inputs or outputs
5. **Network Errors** - Offline mode issues

### Error Recovery

1. **Rollback** - Restore from backup on failure
2. **Partial Success** - Continue with successful components
3. **Graceful Degradation** - Use fallback generators
4. **Clear Messages** - Provide actionable error messages

---

## Performance Considerations

### Optimization Strategies

1. **Parallel Generation** - Generate components concurrently
2. **Tool Caching** - Cache tool availability checks
3. **Incremental Updates** - Only regenerate changed components
4. **Lazy Loading** - Load components on demand

### Resource Management

1. **Timeout Limits** - Prevent hanging operations
2. **Memory Limits** - Control resource usage
3. **Disk Space Checks** - Verify available space
4. **Process Cleanup** - Clean up child processes

---

## Security Architecture

### Security Layers

1. **Input Sanitization** - Clean all user inputs
2. **Path Validation** - Prevent path traversal
3. **Tool Whitelisting** - Only allow known tools
4. **Flag Validation** - Validate all command flags
5. **Resource Limits** - Prevent resource exhaustion

### Security Components

- **Sanitizer** (`pkg/security/sanitizer.go`) - Input sanitization
- **Validator** (`internal/config/validator.go`) - Configuration validation
- **Tool Registry** - Whitelist of allowed tools
- **Rollback Manager** - Recovery from failures

---

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and usage
- [CLI Commands](CLI_COMMANDS.md) - Command reference
- [Adding Tools](ADDING_TOOLS.md) - Extending the system
- [API Reference](API_REFERENCE.md) - Developer API
