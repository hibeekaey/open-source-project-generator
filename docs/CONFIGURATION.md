# Configuration Guide

Complete guide to configuring the Open Source Project Generator.

## Configuration File Format

The generator uses YAML configuration files:

```yaml
# Project metadata
name: "project-name"
description: "Project description"
output_dir: "./output-directory"

# Components to generate
components:
  - type: component-type
    name: component-name
    enabled: true
    config:
      # Component-specific configuration

# Integration settings
integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
  shared_environment:
    KEY: value

# Generation options
options:
  use_external_tools: true
  dry_run: false
  verbose: false
  create_backup: true
  force_overwrite: false
```

## Project Metadata

### Required Fields

```yaml
name: "my-project"           # Project name (required)
output_dir: "./my-project"   # Output directory (required)
```

### Optional Fields

```yaml
description: "Project description"
author: "Your Name"
email: "your@email.com"
license: "MIT"
repository: "https://github.com/user/repo"
```

## Components

Components define what to generate. Each component has:

```yaml
components:
  - type: component-type      # Component type (required)
    name: component-name      # Component name (required)
    enabled: true             # Enable/disable (required)
    config:                   # Component-specific config
      key: value
```

### Supported Component Types

#### Next.js Frontend (`nextjs`)

```yaml
- type: nextjs
  name: web-app
  enabled: true
  config:
    typescript: true          # Use TypeScript (default: true)
    tailwind: true            # Include Tailwind CSS (default: true)
    app_router: true          # Use App Router (default: true)
    eslint: true              # Include ESLint (default: true)
    src_dir: false            # Use src/ directory (default: false)
```

**Bootstrap Tool:** `npx create-next-app`

#### Go Backend (`go-backend`)

```yaml
- type: go-backend
  name: api-server
  enabled: true
  config:
    module: github.com/org/project  # Go module name (required)
    framework: gin                   # Web framework (default: gin)
    port: 8080                       # Server port (default: 8080)
    cors_enabled: true               # Enable CORS (default: false)
    auth_enabled: false              # Include auth (default: false)
```

**Bootstrap Tool:** `go mod init`

#### Android App (`android`)

```yaml
- type: android
  name: mobile-android
  enabled: true
  config:
    package: com.example.app  # Package name (required)
    min_sdk: 24               # Minimum SDK (default: 24)
    target_sdk: 34            # Target SDK (default: 34)
    language: kotlin          # Language (default: kotlin)
    compose: true             # Use Jetpack Compose (default: true)
```

**Bootstrap Tool:** `gradle` (fallback available)

#### iOS App (`ios`)

```yaml
- type: ios
  name: mobile-ios
  enabled: true
  config:
    bundle_id: com.example.app      # Bundle ID (required)
    deployment_target: "15.0"       # Min iOS version (default: "15.0")
    language: swift                 # Language (default: swift)
    swiftui: true                   # Use SwiftUI (default: true)
```

**Bootstrap Tool:** `xcodebuild` (fallback available)

## Integration Settings

Configure how components integrate:

```yaml
integration:
  # Generate Docker Compose file
  generate_docker_compose: true
  
  # Generate build/run scripts
  generate_scripts: true
  
  # API endpoint configuration
  api_endpoints:
    backend: http://localhost:8080
    frontend: http://localhost:3000
  
  # Shared environment variables
  shared_environment:
    NODE_ENV: development
    API_URL: http://localhost:8080
    LOG_LEVEL: info
```

### Docker Compose Generation

When `generate_docker_compose: true`:

- Creates `docker-compose.yml` in project root
- Includes services for all enabled components
- Configures networking between services
- Sets up volume mounts for development

### Script Generation

When `generate_scripts: true`:

- Creates `Makefile` with common commands
- Generates `scripts/` directory with helper scripts
- Includes build, run, test, and deploy scripts

## Generation Options

Control how generation behaves:

```yaml
options:
  # Use external bootstrap tools (npx, go, etc.)
  use_external_tools: true
  
  # Preview mode (don't create files)
  dry_run: false
  
  # Verbose output
  verbose: false
  
  # Create backup before overwriting
  create_backup: true
  
  # Force overwrite existing directory
  force_overwrite: false
```

### Option Details

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `use_external_tools` | boolean | `true` | Use bootstrap tools or force fallback |
| `dry_run` | boolean | `false` | Preview without creating files |
| `verbose` | boolean | `false` | Enable verbose logging |
| `create_backup` | boolean | `true` | Backup existing files |
| `force_overwrite` | boolean | `false` | Overwrite without prompting |

## Complete Examples

### Minimal Configuration

```yaml
name: "simple-project"
output_dir: "./simple-project"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true

integration:
  generate_docker_compose: false
  generate_scripts: true

options:
  use_external_tools: true
```

### Full-Stack Configuration

```yaml
name: "fullstack-app"
description: "Full-stack web application"
output_dir: "./fullstack-app"
author: "Your Name"
license: "MIT"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true
      eslint: true

  - type: go-backend
    name: api-server
    enabled: true
    config:
      module: github.com/myorg/fullstack-app
      framework: gin
      port: 8080
      cors_enabled: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
    frontend: http://localhost:3000
  shared_environment:
    NODE_ENV: development
    API_URL: http://localhost:8080
    DATABASE_URL: postgres://localhost:5432/app

options:
  use_external_tools: true
  create_backup: true
  verbose: false
```

### Mobile App Configuration

```yaml
name: "mobile-app"
description: "Cross-platform mobile application"
output_dir: "./mobile-app"

components:
  - type: android
    name: android-app
    enabled: true
    config:
      package: com.example.mobileapp
      min_sdk: 24
      target_sdk: 34
      language: kotlin
      compose: true

  - type: ios
    name: ios-app
    enabled: true
    config:
      bundle_id: com.example.mobileapp
      deployment_target: "15.0"
      language: swift
      swiftui: true

  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/mobile-app
      framework: gin
      port: 8080

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
  shared_environment:
    API_URL: http://localhost:8080
    API_TIMEOUT: "30"

options:
  use_external_tools: true
  create_backup: true
```

## Environment Variables

Override configuration with environment variables:

```bash
# Project settings
export GENERATOR_PROJECT_NAME="my-project"
export GENERATOR_OUTPUT_DIR="./output"

# Generation options
export GENERATOR_USE_EXTERNAL_TOOLS=true
export GENERATOR_VERBOSE=true
export GENERATOR_DRY_RUN=false
export GENERATOR_FORCE=false
export GENERATOR_CREATE_BACKUP=true
```

Environment variables take precedence over configuration file values.

## Configuration Validation

The generator validates configuration before generation:

### Validation Rules

- Project name is required and must be valid
- Output directory must be writable
- At least one component must be enabled
- Component types must be supported
- Component configurations must be valid
- Integration settings must be valid

### Validation Errors

Common validation errors:

```text
Error: Invalid configuration: missing required field 'name'
Error: Invalid configuration: no components enabled
Error: Invalid component type: 'invalid-type'
Error: Invalid component config: missing required field 'module'
```

## Best Practices

### 1. Use Descriptive Names

```yaml
name: "ecommerce-platform"  # Good
name: "proj1"               # Bad
```

### 2. Enable Only What You Need

```yaml
components:
  - type: nextjs
    enabled: true    # Only enable what you'll use
  - type: android
    enabled: false   # Disable unused components
```

### 3. Set Appropriate Defaults

```yaml
integration:
  shared_environment:
    NODE_ENV: development    # Set sensible defaults
    LOG_LEVEL: info
```

### 4. Version Control Your Configs

```bash
mkdir -p .generator
generator init-config .generator/project.yaml
git add .generator/
```

### 5. Use Comments

```yaml
components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true  # Always use TypeScript for type safety
      tailwind: true    # Tailwind for rapid UI development
```

## Troubleshooting

### Configuration Not Found

```bash
# Check file exists
ls -la project.yaml

# Use absolute path
generator generate --config /full/path/to/project.yaml
```

### Invalid YAML Syntax

```bash
# Validate YAML syntax
# Use online validator or yamllint
yamllint project.yaml
```

### Component Configuration Errors

```bash
# Generate a valid template
generator init-config template.yaml

# Compare with your configuration
diff template.yaml project.yaml
```

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [CLI Commands](CLI_COMMANDS.md) - Command reference
- [Examples](EXAMPLES.md) - Example configurations
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
