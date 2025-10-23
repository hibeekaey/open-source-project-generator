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

**Configuration Fields:**

| Field | Type | Required | Default | Description | Validation |
|-------|------|----------|---------|-------------|------------|
| `typescript` | boolean | No | `true` | Enable TypeScript | Must be boolean |
| `tailwind` | boolean | No | `true` | Include Tailwind CSS | Must be boolean |
| `app_router` | boolean | No | `true` | Use App Router (vs Pages Router) | Must be boolean |
| `eslint` | boolean | No | `true` | Include ESLint configuration | Must be boolean |
| `src_dir` | boolean | No | `false` | Use src/ directory structure | Must be boolean |

**Validation Rules:**

- All configuration values must be boolean (true/false)
- Component name must be a valid directory name
- No special validation required for boolean fields

**Example Configurations:**

```yaml
# TypeScript with Tailwind (recommended)
config:
  typescript: true
  tailwind: true
  app_router: true
  eslint: true

# JavaScript without Tailwind
config:
  typescript: false
  tailwind: false
  app_router: true
  eslint: true

# Pages Router (legacy)
config:
  typescript: true
  tailwind: true
  app_router: false
  eslint: true
```

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

**Configuration Fields:**

| Field | Type | Required | Default | Description | Validation |
|-------|------|----------|---------|-------------|------------|
| `module` | string | **Yes** | - | Go module path | Must be valid Go module path (e.g., github.com/org/project) |
| `framework` | string | No | `gin` | Web framework | Must be one of: `gin`, `echo`, `fiber` |
| `port` | integer | No | `8080` | HTTP server port | Must be 1-65535 |
| `cors_enabled` | boolean | No | `false` | Enable CORS middleware | Must be boolean |
| `auth_enabled` | boolean | No | `false` | Include authentication | Must be boolean |

**Validation Rules:**

- **module**: Required field, must follow Go module path format (domain/path)
  - Valid: `github.com/myorg/project`, `example.com/api`, `myproject`
  - Invalid: `my project`, `project/`, `/project`
- **framework**: Must be one of the supported frameworks
  - Supported: `gin`, `echo`, `fiber`
- **port**: Must be a valid port number
  - Range: 1-65535
  - Common ports: 8080, 8000, 3000, 9000
  - Avoid: 80, 443 (require root), 0 (invalid)

**Example Configurations:**

```yaml
# Minimal configuration
config:
  module: github.com/myorg/api-server

# Full configuration with Gin
config:
  module: github.com/myorg/api-server
  framework: gin
  port: 8080
  cors_enabled: true
  auth_enabled: true

# Echo framework
config:
  module: github.com/myorg/api-server
  framework: echo
  port: 9000
  cors_enabled: true

# Fiber framework (high performance)
config:
  module: github.com/myorg/api-server
  framework: fiber
  port: 3000
```

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

**Configuration Fields:**

| Field | Type | Required | Default | Description | Validation |
|-------|------|----------|---------|-------------|------------|
| `package` | string | **Yes** | - | Java/Kotlin package name | Must be valid Java package name |
| `min_sdk` | integer | No | `24` | Minimum Android SDK level | Must be valid API level (21-34) |
| `target_sdk` | integer | No | `34` | Target Android SDK level | Must be valid API level (21-34), >= min_sdk |
| `language` | string | No | `kotlin` | Programming language | Must be `kotlin` or `java` |
| `compose` | boolean | No | `true` | Use Jetpack Compose | Must be boolean |

**Validation Rules:**

- **package**: Required field, must follow Java package naming conventions
  - Format: lowercase letters, dots as separators, no special characters
  - Valid: `com.example.app`, `com.mycompany.myapp`, `org.project.mobile`
  - Invalid: `com.Example.App`, `com-example-app`, `example`, `com.example.`
  - Must have at least 2 segments (e.g., `com.example`)
- **min_sdk**: Must be a valid Android API level
  - Minimum recommended: 21 (Android 5.0 Lollipop)
  - Common values: 21, 23, 24, 26, 28, 29, 30, 31, 33, 34
  - Must be <= target_sdk
- **target_sdk**: Must be a valid Android API level
  - Should be latest stable version (currently 34)
  - Must be >= min_sdk
- **language**: Must be a supported language
  - `kotlin` (recommended, modern)
  - `java` (legacy support)

**Example Configurations:**

```yaml
# Modern Android app (recommended)
config:
  package: com.mycompany.myapp
  min_sdk: 24
  target_sdk: 34
  language: kotlin
  compose: true

# Legacy Android app
config:
  package: com.mycompany.legacyapp
  min_sdk: 21
  target_sdk: 33
  language: java
  compose: false

# Minimum configuration
config:
  package: com.example.app
```

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

**Configuration Fields:**

| Field | Type | Required | Default | Description | Validation |
|-------|------|----------|---------|-------------|------------|
| `bundle_id` | string | **Yes** | - | iOS bundle identifier | Must be valid bundle ID (reverse domain) |
| `deployment_target` | string | No | `"15.0"` | Minimum iOS version | Must be valid iOS version string |
| `language` | string | No | `swift` | Programming language | Must be `swift` or `objective-c` |
| `swiftui` | boolean | No | `true` | Use SwiftUI framework | Must be boolean |

**Validation Rules:**

- **bundle_id**: Required field, must follow Apple bundle identifier format
  - Format: reverse domain notation, lowercase, dots as separators
  - Valid: `com.example.app`, `com.mycompany.myapp`, `org.project.mobile`
  - Invalid: `com.Example.App`, `com-example-app`, `example`, `com.example.`
  - Must have at least 2 segments (e.g., `com.example`)
  - Only alphanumeric characters and dots allowed
- **deployment_target**: Must be a valid iOS version
  - Format: Major.Minor (e.g., "15.0", "16.0", "17.0")
  - Minimum recommended: "15.0" (iOS 15)
  - Common values: "15.0", "16.0", "17.0", "18.0"
  - Must be quoted as string in YAML
- **language**: Must be a supported language
  - `swift` (recommended, modern)
  - `objective-c` (legacy support)

**Example Configurations:**

```yaml
# Modern iOS app (recommended)
config:
  bundle_id: com.mycompany.myapp
  deployment_target: "15.0"
  language: swift
  swiftui: true

# Legacy iOS app
config:
  bundle_id: com.mycompany.legacyapp
  deployment_target: "13.0"
  language: objective-c
  swiftui: false

# Latest iOS version
config:
  bundle_id: com.example.app
  deployment_target: "17.0"
  language: swift
  swiftui: true

# Minimum configuration
config:
  bundle_id: com.example.app
```

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

The generator validates configuration before generation to catch errors early.

### Validation Process

1. **Syntax Validation**: YAML/JSON syntax must be valid
2. **Schema Validation**: Required fields must be present
3. **Type Validation**: Field types must match expected types
4. **Component Validation**: Component-specific rules must pass
5. **Integration Validation**: Integration settings must be valid

### Project-Level Validation Rules

- **name**: Required, must be non-empty string, valid directory name
- **output_dir**: Required, must be valid path, writable location
- **components**: Required, at least one component must be enabled
- **description**: Optional, string
- **author**: Optional, string
- **license**: Optional, string

### Component-Level Validation Rules

#### Next.js Component

- **type**: Must be `nextjs`
- **name**: Required, valid directory name
- **enabled**: Required, boolean
- **config.typescript**: Optional, boolean (default: true)
- **config.tailwind**: Optional, boolean (default: true)
- **config.app_router**: Optional, boolean (default: true)
- **config.eslint**: Optional, boolean (default: true)
- **config.src_dir**: Optional, boolean (default: false)

#### Go Backend Component

- **type**: Must be `go-backend`
- **name**: Required, valid directory name
- **enabled**: Required, boolean
- **config.module**: **Required**, valid Go module path
  - Must contain at least one slash or be a simple name
  - Examples: `github.com/org/project`, `example.com/api`, `myproject`
- **config.framework**: Optional, must be `gin`, `echo`, or `fiber` (default: gin)
- **config.port**: Optional, integer 1-65535 (default: 8080)
- **config.cors_enabled**: Optional, boolean (default: false)
- **config.auth_enabled**: Optional, boolean (default: false)

#### Android Component

- **type**: Must be `android`
- **name**: Required, valid directory name
- **enabled**: Required, boolean
- **config.package**: **Required**, valid Java package name
  - Must be lowercase, dots as separators
  - At least 2 segments (e.g., `com.example`)
  - Only alphanumeric and dots allowed
- **config.min_sdk**: Optional, integer 21-34 (default: 24)
- **config.target_sdk**: Optional, integer 21-34, >= min_sdk (default: 34)
- **config.language**: Optional, must be `kotlin` or `java` (default: kotlin)
- **config.compose**: Optional, boolean (default: true)

#### iOS Component

- **type**: Must be `ios`
- **name**: Required, valid directory name
- **enabled**: Required, boolean
- **config.bundle_id**: **Required**, valid bundle identifier
  - Reverse domain notation (e.g., `com.example.app`)
  - At least 2 segments
  - Only alphanumeric and dots allowed
- **config.deployment_target**: Optional, valid iOS version string (default: "15.0")
- **config.language**: Optional, must be `swift` or `objective-c` (default: swift)
- **config.swiftui**: Optional, boolean (default: true)

### Validation Errors

Common validation errors and how to fix them:

#### Missing Required Fields

```text
Error: Invalid configuration: missing required field 'name'
```

**Fix**: Add the required field to your configuration:

```yaml
name: "my-project"  # Add this
output_dir: "./my-project"
```

#### Invalid Component Type

```text
Error: Invalid component type: 'react-app'
```

**Fix**: Use a supported component type:

```yaml
components:
  - type: nextjs  # Change from 'react-app' to 'nextjs'
    name: web-app
    enabled: true
```

#### Missing Component Configuration

```text
Error: Invalid component config: missing required field 'module'
```

**Fix**: Add the required configuration field:

```yaml
- type: go-backend
  name: api
  enabled: true
  config:
    module: github.com/myorg/project  # Add this required field
```

#### Invalid Package Name

```text
Error: Invalid Android package name: 'Com.Example.App'
```

**Fix**: Use lowercase package name:

```yaml
config:
  package: com.example.app  # Must be lowercase
```

#### Invalid Port Number

```text
Error: Invalid port number: 70000 (must be 1-65535)
```

**Fix**: Use a valid port number:

```yaml
config:
  port: 8080  # Must be 1-65535
```

#### Invalid Bundle ID

```text
Error: Invalid iOS bundle ID: 'example'
```

**Fix**: Use reverse domain notation:

```yaml
config:
  bundle_id: com.example.app  # Must have at least 2 segments
```

#### SDK Version Mismatch

```text
Error: target_sdk (30) must be >= min_sdk (33)
```

**Fix**: Ensure target_sdk >= min_sdk:

```yaml
config:
  min_sdk: 24
  target_sdk: 34  # Must be >= min_sdk
```

### Validation Best Practices

1. **Use init-config**: Generate a valid template to start from
   ```bash
   generator init-config --example fullstack
   ```

2. **Test with dry-run**: Validate before generating
   ```bash
   generator generate --config project.yaml --dry-run
   ```

3. **Check syntax**: Validate YAML syntax
   ```bash
   yamllint project.yaml
   ```

4. **Use verbose mode**: See detailed validation errors
   ```bash
   generator generate --config project.yaml --verbose
   ```

5. **Validate incrementally**: Add components one at a time and validate

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
