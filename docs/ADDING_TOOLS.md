# Adding Bootstrap Tools

Step-by-step guide for adding support for new bootstrap tools to the Open Source Project Generator.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Step-by-Step Guide](#step-by-step-guide)
- [Example: Adding Vite Support](#example-adding-vite-support)
- [Testing](#testing)
- [Best Practices](#best-practices)

---

## Overview

Adding a new bootstrap tool involves:

1. Defining tool metadata (name, version requirements, installation instructions)
2. Creating an executor to run the tool
3. Registering the executor with the coordinator
4. Adding structure mapping
5. Creating a fallback generator (optional)
6. Writing tests
7. Updating documentation

---

## Prerequisites

Before adding a new bootstrap tool, ensure you have:

- Go 1.25+ installed
- Familiarity with the tool you're adding
- Understanding of the tool's CLI interface
- Access to the tool for testing
- Read the [Architecture](ARCHITECTURE.md) documentation

---

## Step-by-Step Guide

### Step 1: Define Tool Metadata

Add your tool to the registry in `internal/orchestrator/tool_discovery.go`:

```go
func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        tools: map[string]*ToolMetadata{
            // ... existing tools ...
            
            "your-tool": {
                Name:            "your-tool",
                Command:         "your-tool",           // Command to execute
                VersionFlag:     "--version",           // Flag to get version
                MinVersion:      "1.0.0",              // Minimum required version
                ComponentTypes:  []string{"your-component"},
                InstallDocs: map[string]string{
                    "darwin":  "brew install your-tool",
                    "linux":   "sudo apt install your-tool",
                    "windows": "https://your-tool.com/install",
                },
                FallbackAvailable: true,               // Whether fallback exists
            },
        },
    }
}
```

**Field Descriptions:**

| Field | Description |
|-------|-------------|
| `Name` | Tool identifier (used in configuration) |
| `Command` | Actual command to execute |
| `VersionFlag` | Flag to check tool version |
| `MinVersion` | Minimum version required (semver format) |
| `ComponentTypes` | Component types this tool supports |
| `InstallDocs` | OS-specific installation instructions |
| `FallbackAvailable` | Whether a fallback generator exists |

### Step 2: Create Executor

Create a new file `internal/generator/bootstrap/yourtool.go`:

```go
package bootstrap

import (
    "context"
    "fmt"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/logger"
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

type YourToolExecutor struct {
    *BaseExecutor
}

func NewYourToolExecutor(log *logger.Logger) *YourToolExecutor {
    return &YourToolExecutor{
        BaseExecutor: &BaseExecutor{logger: log},
    }
}

func (yte *YourToolExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    // Extract configuration
    name := spec.Config["name"].(string)
    
    // Build command flags
    flags := []string{
        "create",
        name,
        "--flag1", spec.Config["option1"].(string),
        "--flag2", spec.Config["option2"].(string),
    }
    
    // Set tool and flags
    spec.Tool = "your-tool"
    spec.Flags = flags
    
    // Execute using base executor
    return yte.BaseExecutor.Execute(ctx, spec)
}

func (yte *YourToolExecutor) SupportsComponent(componentType string) bool {
    return componentType == "your-tool"
}

func (yte *YourToolExecutor) GetDefaultFlags(componentType string) []string {
    return []string{"--default-flag"}
}
```

**Key Points:**

- Extend `BaseExecutor` for common functionality
- Extract configuration from `spec.Config`
- Build appropriate command flags
- Call `BaseExecutor.Execute()` to run the tool
- Implement `SupportsComponent()` to identify supported types
- Provide default flags via `GetDefaultFlags()`

### Step 3: Register Executor

Register your executor in the coordinator. Create or update the executor registry:

```go
// In internal/orchestrator/coordinator.go or a new registry file

func (pc *ProjectCoordinator) getExecutor(componentType string) (BootstrapExecutor, error) {
    executors := map[string]BootstrapExecutor{
        "nextjs":     bootstrap.NewNextJSExecutor(pc.logger),
        "go-backend": bootstrap.NewGoExecutor(pc.logger),
        "android":    bootstrap.NewAndroidExecutor(pc.logger),
        "ios":        bootstrap.NewiOSExecutor(pc.logger),
        "your-tool":  bootstrap.NewYourToolExecutor(pc.logger), // Add here
    }
    
    executor, ok := executors[componentType]
    if !ok {
        return nil, fmt.Errorf("no executor found for component type: %s", componentType)
    }
    
    return executor, nil
}
```

### Step 4: Add Structure Mapping

Update `internal/generator/mapper/structure.go` to map your component's output:

```go
var ComponentMappings = map[string]string{
    "nextjs":      "App/",
    "go-backend":  "CommonServer/",
    "android":     "Mobile/android/",
    "ios":         "Mobile/ios/",
    "your-tool":   "YourTool/", // Add mapping
}
```

**Mapping Guidelines:**

- Use consistent naming conventions
- Consider project structure standards
- Group related components together

### Step 5: Add Configuration Schema

Update configuration validation to support your component type:

```go
// In internal/config/validator.go or schema.go

var ValidComponentTypes = []string{
    "nextjs",
    "go-backend",
    "android",
    "ios",
    "your-tool", // Add here
}

func (v *Validator) validateComponent(component *models.ComponentConfig) error {
    // Validate component type
    if !contains(ValidComponentTypes, component.Type) {
        return fmt.Errorf("invalid component type: %s", component.Type)
    }
    
    // Add type-specific validation
    if component.Type == "your-tool" {
        if _, ok := component.Config["required_field"]; !ok {
            return fmt.Errorf("your-tool requires 'required_field' in config")
        }
    }
    
    return nil
}
```

### Step 6: Create Fallback Generator (Optional)

If your tool might not be available, create a fallback generator:

```go
// internal/generator/fallback/yourtool.go
package fallback

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

type YourToolFallbackGenerator struct {
    logger *logger.Logger
}

func NewYourToolFallbackGenerator(log *logger.Logger) *YourToolFallbackGenerator {
    return &YourToolFallbackGenerator{logger: log}
}

func (ytfg *YourToolFallbackGenerator) Generate(ctx context.Context, spec *FallbackSpec) (*GenerationResult, error) {
    // Create minimal structure
    if err := ytfg.createDirectory(spec.TargetDir); err != nil {
        return nil, err
    }
    
    // Generate essential files
    files := map[string]string{
        "README.md":    ytfg.generateReadme(spec),
        "config.yaml":  ytfg.generateConfig(spec),
        "main.go":      ytfg.generateMain(spec),
    }
    
    for filename, content := range files {
        path := filepath.Join(spec.TargetDir, filename)
        if err := os.WriteFile(path, []byte(content), 0644); err != nil {
            return nil, err
        }
    }
    
    return &GenerationResult{
        Success: true,
        Method:  "fallback",
        ManualSteps: []string{
            "Install your-tool: https://your-tool.com/install",
            "Run: your-tool init",
            "Configure: edit config.yaml",
        },
    }, nil
}

func (ytfg *YourToolFallbackGenerator) generateReadme(spec *FallbackSpec) string {
    return fmt.Sprintf(`# %s

This project was generated using fallback mode because your-tool was not available.

## Setup Instructions

1. Install your-tool: https://your-tool.com/install
2. Run: your-tool init
3. Configure: edit config.yaml

## Next Steps

- Read the documentation
- Configure your project
- Start development
`, spec.Name)
}
```

Register the fallback generator:

```go
// In internal/generator/fallback/registry.go

func DefaultRegistry() *Registry {
    return &Registry{
        generators: map[string]Generator{
            "android":    NewAndroidFallbackGenerator(),
            "ios":        NewiOSFallbackGenerator(),
            "your-tool":  NewYourToolFallbackGenerator(), // Add here
        },
    }
}
```

### Step 7: Write Tests

Create `internal/generator/bootstrap/yourtool_test.go`:

```go
package bootstrap

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

func TestYourToolExecutor_Execute(t *testing.T) {
    executor := NewYourToolExecutor(logger.NewLogger())
    
    spec := &BootstrapSpec{
        ComponentType: "your-tool",
        TargetDir:     t.TempDir(),
        Config: map[string]interface{}{
            "name":    "test-project",
            "option1": "value1",
            "option2": "value2",
        },
    }
    
    result, err := executor.Execute(context.Background(), spec)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.True(t, result.Success)
    assert.Equal(t, "your-tool", result.ToolUsed)
}

func TestYourToolExecutor_SupportsComponent(t *testing.T) {
    executor := NewYourToolExecutor(logger.NewLogger())
    
    assert.True(t, executor.SupportsComponent("your-tool"))
    assert.False(t, executor.SupportsComponent("other-tool"))
}

func TestYourToolExecutor_GetDefaultFlags(t *testing.T) {
    executor := NewYourToolExecutor(logger.NewLogger())
    
    flags := executor.GetDefaultFlags("your-tool")
    
    assert.NotEmpty(t, flags)
    assert.Contains(t, flags, "--default-flag")
}
```

Add integration test:

```go
func TestYourToolExecutor_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Check if tool is available
    if _, err := exec.LookPath("your-tool"); err != nil {
        t.Skip("your-tool not available")
    }
    
    executor := NewYourToolExecutor(logger.NewLogger())
    tmpDir := t.TempDir()
    
    spec := &BootstrapSpec{
        ComponentType: "your-tool",
        TargetDir:     tmpDir,
        Config: map[string]interface{}{
            "name": "integration-test",
        },
    }
    
    result, err := executor.Execute(context.Background(), spec)
    
    assert.NoError(t, err)
    assert.True(t, result.Success)
    
    // Verify generated files
    assert.FileExists(t, filepath.Join(tmpDir, "expected-file.txt"))
}
```

### Step 8: Update Documentation

Update the following documentation files:

**1. Configuration Guide** (`docs/CONFIGURATION.md`):

```markdown
#### Your Tool (`your-tool`)

\`\`\`yaml
- type: your-tool
  name: your-app
  enabled: true
  config:
    option1: value1          # Description (default: default1)
    option2: value2          # Description (default: default2)
\`\`\`

**Bootstrap Tool:** `your-tool`
```

**2. Getting Started** (`docs/GETTING_STARTED.md`):

```markdown
### Your Tool

**Tool:** `your-tool`

\`\`\`bash
# Install Your Tool
# macOS
brew install your-tool

# Ubuntu/Debian
sudo apt install your-tool
\`\`\`

**What it does:** Creates projects using your-tool CLI

**Fallback:** If tools unavailable, generates minimal structure
```

**3. Examples** (`docs/EXAMPLES.md`):

Add example configuration using your tool.

---

## Example: Adding Vite Support

Let's walk through a complete example of adding Vite support.

### 1. Tool Metadata

```go
// internal/orchestrator/tool_discovery.go
"vite": {
    Name:            "vite",
    Command:         "npm",
    VersionFlag:     "create vite@6.0.0 -- --version",
    MinVersion:      "4.0.0",
    ComponentTypes:  []string{"vite"},
    InstallDocs: map[string]string{
        "darwin":  "brew install node",
        "linux":   "sudo apt install nodejs npm",
        "windows": "https://nodejs.org/",
    },
    FallbackAvailable: false,
},
```

### 2. Executor Implementation

```go
// internal/generator/bootstrap/vite.go
package bootstrap

import (
    "context"
    "github.com/cuesoftinc/open-source-project-generator/pkg/logger"
)

type ViteExecutor struct {
    *BaseExecutor
}

func NewViteExecutor(log *logger.Logger) *ViteExecutor {
    return &ViteExecutor{
        BaseExecutor: &BaseExecutor{logger: log},
    }
}

func (ve *ViteExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    name := spec.Config["name"].(string)
    template := spec.Config["template"].(string)
    
    flags := []string{
        "create",
        "vite@6.0.0",
        name,
        "--",
        "--template", template,
    }
    
    spec.Tool = "npm"
    spec.Flags = flags
    
    return ve.BaseExecutor.Execute(ctx, spec)
}

func (ve *ViteExecutor) SupportsComponent(componentType string) bool {
    return componentType == "vite"
}

func (ve *ViteExecutor) GetDefaultFlags(componentType string) []string {
    return []string{"--template", "react-ts"}
}
```

### 3. Registration

```go
// In coordinator or registry
"vite": bootstrap.NewViteExecutor(pc.logger),
```

### 4. Structure Mapping

```go
// internal/generator/mapper/structure.go
var ComponentMappings = map[string]string{
    // ... existing mappings ...
    "vite": "App/vite/",
}
```

### 5. Configuration Schema

```go
// internal/config/validator.go
var ValidComponentTypes = []string{
    "nextjs", "go-backend", "android", "ios", "vite",
}

func (v *Validator) validateComponent(component *models.ComponentConfig) error {
    // ... existing validation ...
    
    if component.Type == "vite" {
        if _, ok := component.Config["template"]; !ok {
            component.Config["template"] = "react-ts" // Default template
        }
    }
    
    return nil
}
```

### 6. Tests

```go
// internal/generator/bootstrap/vite_test.go
func TestViteExecutor_Execute(t *testing.T) {
    executor := NewViteExecutor(logger.NewLogger())
    
    spec := &BootstrapSpec{
        ComponentType: "vite",
        TargetDir:     t.TempDir(),
        Config: map[string]interface{}{
            "name":     "test-vite-app",
            "template": "react-ts",
        },
    }
    
    result, err := executor.Execute(context.Background(), spec)
    
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Equal(t, "npm", result.ToolUsed)
}
```

### 7. Configuration Example

```yaml
# Example configuration
name: "vite-project"
output_dir: "./vite-project"

components:
  - type: vite
    name: web-app
    enabled: true
    config:
      template: react-ts  # react, react-ts, vue, vue-ts, etc.

integration:
  generate_docker_compose: true
```

---

## Testing

### Unit Tests

```bash
# Run unit tests
go test ./internal/generator/bootstrap/yourtool_test.go -v

# Run with coverage
go test ./internal/generator/bootstrap/yourtool_test.go -cover
```

### Integration Tests

```bash
# Run integration tests (requires tool installed)
go test ./internal/generator/bootstrap/yourtool_test.go -v -run Integration

# Skip integration tests
go test ./internal/generator/bootstrap/yourtool_test.go -v -short
```

### Manual Testing

```bash
# Build the CLI
make build

# Create test configuration
cat > test-config.yaml << EOF
name: "test-project"
output_dir: "./test-output"
components:
  - type: your-tool
    name: test-app
    enabled: true
    config:
      option1: value1
EOF

# Test generation
./bin/generator generate --config test-config.yaml --verbose

# Verify output
ls -la ./test-output
```

---

## Best Practices

### 1. Error Handling

Always provide clear error messages:

```go
func (yte *YourToolExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    // Validate required configuration
    name, ok := spec.Config["name"].(string)
    if !ok || name == "" {
        return nil, fmt.Errorf("your-tool requires 'name' in configuration")
    }
    
    // Execute with error context
    result, err := yte.BaseExecutor.Execute(ctx, spec)
    if err != nil {
        return nil, fmt.Errorf("failed to execute your-tool: %w", err)
    }
    
    return result, nil
}
```

### 2. Configuration Validation

Validate configuration early:

```go
func (v *Validator) validateYourToolComponent(component *models.ComponentConfig) error {
    // Check required fields
    if _, ok := component.Config["required_field"]; !ok {
        return fmt.Errorf("your-tool requires 'required_field'")
    }
    
    // Validate field values
    if val, ok := component.Config["option"].(string); ok {
        validOptions := []string{"opt1", "opt2", "opt3"}
        if !contains(validOptions, val) {
            return fmt.Errorf("invalid option: %s (must be one of: %v)", val, validOptions)
        }
    }
    
    return nil
}
```

### 3. Logging

Add comprehensive logging:

```go
func (yte *YourToolExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    yte.logger.Info("Executing your-tool",
        "component", spec.ComponentType,
        "target", spec.TargetDir)
    
    result, err := yte.BaseExecutor.Execute(ctx, spec)
    
    if err != nil {
        yte.logger.Error("Your-tool execution failed",
            "error", err,
            "stderr", result.Stderr)
        return result, err
    }
    
    yte.logger.Info("Your-tool execution completed",
        "duration", result.Duration,
        "success", result.Success)
    
    return result, nil
}
```

### 4. Default Values

Provide sensible defaults:

```go
func (yte *YourToolExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*ExecutionResult, error) {
    // Set defaults for optional fields
    if _, ok := spec.Config["option1"]; !ok {
        spec.Config["option1"] = "default-value"
    }
    
    if _, ok := spec.Config["option2"]; !ok {
        spec.Config["option2"] = true
    }
    
    // Continue with execution
    // ...
}
```

### 5. Documentation

Document configuration options:

```go
// YourToolExecutor executes the your-tool CLI to generate projects.
//
// Required Configuration:
//   - name: Project name
//   - template: Project template
//
// Optional Configuration:
//   - option1: Description (default: "default-value")
//   - option2: Description (default: true)
//
// Example:
//   config:
//     name: my-project
//     template: basic
//     option1: custom-value
type YourToolExecutor struct {
    *BaseExecutor
}
```

---

## Checklist

Before submitting your changes:

- [ ] Tool metadata defined in `tool_discovery.go`
- [ ] Executor implemented in `internal/generator/bootstrap/`
- [ ] Executor registered with coordinator
- [ ] Structure mapping added
- [ ] Configuration validation added
- [ ] Fallback generator created (if applicable)
- [ ] Unit tests written
- [ ] Integration tests written
- [ ] Documentation updated
- [ ] Manual testing completed
- [ ] Code follows project conventions
- [ ] All tests pass

---

## See Also

- [Architecture](ARCHITECTURE.md) - System architecture
- [Getting Started](GETTING_STARTED.md) - Installation and usage
- [Configuration Guide](CONFIGURATION.md) - Configuration options
- [API Reference](API_REFERENCE.md) - Developer API
