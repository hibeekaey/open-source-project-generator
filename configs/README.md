# Configuration Examples

This directory contains example configuration files for the Open Source Project Generator. These examples demonstrate different project types and configuration options.

## Available Examples

### Quick Start

- **`minimal.yaml`** - Minimal configuration with only required fields
  - Demonstrates bare minimum configuration
  - All optional settings use defaults
  - Best starting point for new users

### Full-Stack Applications

- **`fullstack-complete.yaml`** - Complete full-stack application with all components
  - Next.js frontend
  - Go backend API
  - Android and iOS mobile apps
  - Docker Compose integration
  - Build scripts and documentation

- **`mobile-app-with-backend.yaml`** - Mobile apps with shared backend
  - Android (Kotlin) and iOS (Swift) apps
  - Go backend API
  - No web frontend
  - Ideal for mobile-first applications

### Single Component Projects

- **`frontend-only.yaml`** - Standalone frontend application
  - Next.js application only
  - Ideal for JAMstack, static sites, or connecting to existing APIs
  - No backend or mobile components

- **`backend-only.yaml`** - Standalone backend API
  - Go API server with Gin framework
  - Docker Compose for services (database, cache)
  - Ideal for microservices or API-first development

- **`mobile-only.yaml`** - Native mobile applications
  - Android (Kotlin) and iOS (Swift) apps
  - Connects to external APIs
  - No web frontend or backend

### Advanced Examples

- **`advanced-options.yaml`** - Comprehensive configuration reference
  - Demonstrates all available configuration options
  - Detailed comments explaining each option
  - Use as reference documentation

- **`performance-optimized.yaml`** - Performance-focused configuration
  - Shows performance-specific options
  - Parallel generation settings
  - Optimized for large projects

## Usage

Generate a project using any example configuration:

```bash
# Using a specific configuration file
project-generator generate --config configs/fullstack-complete.yaml

# Preview without creating files (dry run)
project-generator generate --config configs/minimal.yaml --dry-run

# Verbose output for debugging
project-generator generate --config configs/backend-only.yaml --verbose
```

## Configuration File Structure

All configuration files follow this structure:

```yaml
# Project metadata
name: project-name                    # Required: Project name
description: Project description      # Optional: Project description
output_dir: ./output-path            # Required: Where to generate project

# Component definitions
components:
  - type: component-type             # Required: Component type
    name: component-name             # Required: Component name
    enabled: true                    # Required: Enable/disable component
    config:                          # Required: Component-specific config
      # Component-specific options

# Integration configuration
integration:
  generate_docker_compose: true      # Optional: Generate Docker Compose
  generate_scripts: true             # Optional: Generate build scripts
  api_endpoints:                     # Optional: API endpoint configuration
    backend: http://localhost:8080
  shared_environment:                # Optional: Shared environment variables
    KEY: value

# Generation options
options:
  use_external_tools: true           # Optional: Use external CLI tools
  dry_run: false                     # Optional: Preview without creating files
  verbose: false                     # Optional: Detailed output
  create_backup: true                # Optional: Backup existing files
  force_overwrite: false             # Optional: Overwrite without prompting
```

## Component Types

### Frontend Components

#### `nextjs` - Next.js Web Application

```yaml
- type: nextjs
  name: web-app
  enabled: true
  config:
    name: web-app              # Required: Project name
    typescript: true           # Optional: Use TypeScript (default: true)
    tailwind: true             # Optional: Include Tailwind CSS (default: true)
    app_router: true           # Optional: Use App Router (default: true)
    eslint: true               # Optional: Include ESLint (default: true)
```

**Generated with:** `npx create-next-app@latest` (when `use_external_tools: true`)

### Backend Components

#### `go-backend` - Go API Server

```yaml
- type: go-backend
  name: api-server
  enabled: true
  config:
    name: api-server                        # Required: Server name
    module: github.com/user/project         # Required: Go module path
    framework: gin                          # Optional: Framework (default: gin)
    port: 8080                              # Optional: Server port (default: 8080)
```

**Generated with:** `go mod init` + `go get` (when `use_external_tools: true`)

### Mobile Components

#### `android` - Android Application

```yaml
- type: android
  name: mobile-android
  enabled: true
  config:
    name: mobile-android                    # Required: App name
    package: com.example.app                # Required: Package name
    min_sdk: 24                             # Optional: Min SDK (default: 24)
    target_sdk: 34                          # Optional: Target SDK (default: 34)
    language: kotlin                        # Optional: Language (default: kotlin)
```

**Generated with:** Gradle/Android Studio CLI (when available) or fallback templates

#### `ios` - iOS Application

```yaml
- type: ios
  name: mobile-ios
  enabled: true
  config:
    name: mobile-ios                        # Required: App name
    bundle_id: com.example.app              # Required: Bundle identifier
    deployment_target: "15.0"               # Optional: Min iOS version (default: "15.0")
    language: swift                         # Optional: Language (default: swift)
```

**Generated with:** Xcode CLI (when available) or fallback templates

## Integration Options

### Docker Compose Generation

```yaml
integration:
  generate_docker_compose: true
```

Generates a `docker-compose.yml` file that orchestrates all enabled components for local development.

### Build Scripts

```yaml
integration:
  generate_scripts: true
```

Generates build and run scripts for all components:
- `build.sh` / `build.bat` - Build all components
- `run.sh` / `run.bat` - Run all components
- `dev.sh` / `dev.bat` - Development mode

### API Endpoints

```yaml
integration:
  api_endpoints:
    backend: http://localhost:8080
    frontend: http://localhost:3000
```

Configures how components connect to each other. Frontend and mobile apps will use these endpoints.

### Shared Environment Variables

```yaml
integration:
  shared_environment:
    NODE_ENV: development
    LOG_LEVEL: info
    API_TIMEOUT: "30"
```

Environment variables shared across all components. Generated in `.env` files.

## Generation Options

### Use External Tools

```yaml
options:
  use_external_tools: true
```

- `true`: Use external CLI tools (npx, go, gradle) when available
- `false`: Always use fallback template generation

### Dry Run

```yaml
options:
  dry_run: true
```

Preview what would be generated without creating any files.

### Verbose Output

```yaml
options:
  verbose: true
```

Show detailed output during generation, including:
- Tool detection results
- Command execution details
- File operations
- Validation results

### Create Backup

```yaml
options:
  create_backup: true
```

Create backups of existing directories before overwriting.

### Force Overwrite

```yaml
options:
  force_overwrite: true
```

Overwrite existing directories without prompting for confirmation.

## Validation Rules

### Project Name

- Required field
- 1-100 characters
- Alphanumeric, dash, and underscore only
- Cannot start or end with dash or underscore

### Output Directory

- Required field
- Valid file system path
- No path traversal (`..`)
- No invalid characters (`*`, `?`, `"`, `<`, `>`, `|`)

### Go Module Path

- Required for `go-backend` components
- Must be a valid Go module path (e.g., `github.com/user/project`)
- Must contain at least one `/`

### Android Package Name

- Required for `android` components
- Must be a valid Java package name (e.g., `com.example.app`)
- At least two segments separated by `.`
- Each segment starts with a letter
- Only letters, digits, and underscores

### iOS Bundle Identifier

- Required for `ios` components
- Must be a valid bundle identifier (e.g., `com.example.app`)
- At least two segments separated by `.`
- Each segment starts with a letter
- Only letters, digits, and hyphens

### Port Number

- Optional for `go-backend` components
- Must be between 1 and 65535
- Default: 8080

## Creating Custom Configurations

1. Start with an example that's closest to your needs
2. Copy and modify the configuration file
3. Update project metadata (name, description, output_dir)
4. Enable/disable components as needed
5. Configure component-specific options
6. Set integration and generation options
7. Validate your configuration:

```bash
project-generator generate --config your-config.yaml --dry-run
```

## Tips and Best Practices

### Component Naming

- Use descriptive names that indicate the component's purpose
- Use kebab-case for consistency (e.g., `web-app`, `api-server`)
- Keep names short but meaningful

### Module Paths

- Use your actual repository path for Go modules
- Example: `github.com/username/project-name`
- This ensures proper import paths in generated code

### Package Names

- Follow platform conventions:
  - Android: Reverse domain notation (e.g., `com.company.app`)
  - iOS: Reverse domain notation (e.g., `com.company.app`)
- Use lowercase for Android packages
- Avoid reserved keywords

### Environment Variables

- Use `NEXT_PUBLIC_` prefix for Next.js client-side variables
- Keep sensitive values out of configuration files
- Use `.env.example` files for documentation

### Development Workflow

1. Start with `dry_run: true` to preview
2. Review the preview output
3. Set `dry_run: false` to generate
4. Use `verbose: true` for troubleshooting
5. Keep `create_backup: true` for safety

## Troubleshooting

### "Component type not supported"

Ensure the component type is one of:
- `nextjs`
- `go-backend`
- `android`
- `ios`

### "Required config field missing"

Check that all required fields are present for each component type. See component documentation above.

### "Validation failed"

Review validation rules for the specific field mentioned in the error message.

### "Tool not found"

If external tools are not available:
- Install the required tools (npx, go, gradle, xcodebuild)
- Or set `use_external_tools: false` to use fallback generation

## Additional Resources

- [Getting Started Guide](../docs/GETTING_STARTED.md)
- [Configuration Reference](../docs/CONFIGURATION.md)
- [Template Development](../docs/TEMPLATE_DEVELOPMENT.md)
- [Troubleshooting](../docs/TROUBLESHOOTING.md)

## Support

For issues or questions:
- Check the [Troubleshooting Guide](../docs/TROUBLESHOOTING.md)
- Review existing [GitHub Issues](https://github.com/cuesoftinc/open-source-project-generator/issues)
- Create a new issue with your configuration file and error output
