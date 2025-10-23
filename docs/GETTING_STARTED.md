# Getting Started

This guide will help you install and use the Open Source Project Generator to create production-ready projects.

## What is This Tool?

The Open Source Project Generator uses a **tool-orchestration architecture** that delegates project creation to industry-standard CLI tools like `create-next-app`, `go mod init`, and others. This means:

- ✅ **Always up-to-date** - No manual template maintenance
- ✅ **Industry-standard** - Uses official framework CLIs
- ✅ **Graceful fallback** - Works even when tools are unavailable
- ✅ **Offline support** - Can work without internet after initial setup

## Installation

### Quick Install

**Linux/macOS:**

```bash
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash
```

**Using Go:**

```bash
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest
```

**Using Docker:**

```bash
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/cuesoftinc/open-source-project-generator
cd open-source-project-generator

# Build
make build

# Install (optional)
sudo cp bin/generator /usr/local/bin/
```

## Quick Start

### 1. Check Your Environment

Before generating projects, check which bootstrap tools are available:

```bash
generator check-tools
```

**Example output:**

```text
Checking 4 tools...

✓ Available Tools:
  ✓ npx (version: 10.2.3)
    Supports: [nextjs]
  ✓ go (version: 1.21.0)
    Supports: [go-backend]

✗ Missing Tools:
  ✗ gradle
    Required for: [android]
    
    Install instructions:
    macOS: brew install gradle
    Ubuntu: sudo apt install gradle

2 of 4 tools available
```

### 2. Create a Configuration File

Generate a configuration template:

```bash
generator init-config my-project.yaml
```

This creates a YAML file you can customize. Example:

```yaml
name: "my-awesome-project"
description: "A full-stack web application"
output_dir: "./my-awesome-project"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true

  - type: go-backend
    name: api-server
    enabled: true
    config:
      module: github.com/myorg/my-awesome-project
      framework: gin
      port: 8080

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080

options:
  use_external_tools: true
  create_backup: true
  verbose: false
```

### 3. Generate Your Project

```bash
generator generate --config my-project.yaml
```

**What happens:**

1. Validates your configuration
2. Checks which bootstrap tools are available
3. Executes tools (e.g., `npx create-next-app`, `go mod init`)
4. Maps outputs to standardized structure
5. Generates integration files (Docker Compose, scripts, etc.)

### 4. Explore the Generated Project

```bash
cd my-awesome-project
ls -la
```

**Generated structure:**

```text
my-awesome-project/
├── App/                    # Next.js frontend
│   ├── app/
│   ├── public/
│   ├── package.json
│   └── next.config.js
├── CommonServer/          # Go backend
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── go.mod
│   └── main.go
├── Deploy/
│   └── docker/
│       ├── Dockerfile.frontend
│       └── Dockerfile.backend
├── docker-compose.yml
├── Makefile
└── README.md
```

## Tool Requirements

The generator uses external bootstrap tools to create projects. Here's what you need:

### Frontend (Next.js)

**Tool:** `npx` (comes with Node.js)

```bash
# Install Node.js
# macOS
brew install node

# Ubuntu/Debian
sudo apt install nodejs npm

# Windows
# Download from https://nodejs.org/
```

**What it does:** Runs `npx create-next-app@latest` with your specified options

### Backend (Go)

**Tool:** `go`

```bash
# Install Go
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Windows
# Download from https://golang.org/dl/
```

**What it does:** Runs `go mod init` and sets up a Gin-based API server

### Mobile (Android)

**Tool:** `gradle` (or Android Studio)

```bash
# Install Gradle
# macOS
brew install gradle

# Ubuntu/Debian
sudo apt install gradle
```

**Fallback:** If Gradle is unavailable, generates minimal structure with setup instructions

### Mobile (iOS)

**Tool:** `xcodebuild` (comes with Xcode, macOS only)

```bash
# Install Xcode Command Line Tools
xcode-select --install
```

**Fallback:** If unavailable, generates minimal structure with setup instructions

## Common Commands

### Generate Projects

```bash
# From configuration file
generator generate --config project.yaml

# With custom output directory
generator generate --config project.yaml --output ./my-project

# Dry run (preview without creating files)
generator generate --config project.yaml --dry-run

# Force fallback (don't use external tools)
generator generate --config project.yaml --no-external-tools

# Verbose output
generator generate --config project.yaml --verbose
```

### Check Tools

```bash
# Check all registered tools
generator check-tools

# Check with verbose output
generator check-tools --verbose

# Check specific tools
generator check-tools npx go
```

### Configuration Templates

```bash
# Generate default configuration
generator init-config

# Generate minimal configuration
generator init-config --minimal

# Generate example configuration
generator init-config --example fullstack
generator init-config --example frontend
generator init-config --example backend
generator init-config --example mobile
```

### Cache Management

```bash
# View cache statistics
generator cache-tools --stats

# Save current tool availability
generator cache-tools --save

# Clear cache
generator cache-tools --clear

# Show cache information
generator cache-tools --info
```

## Configuration File Format

The configuration file defines your project structure:

```yaml
# Project metadata
name: "project-name"
description: "Project description"
output_dir: "./output-directory"

# Components to generate
components:
  - type: nextjs              # Component type
    name: web-app             # Component name
    enabled: true             # Enable/disable
    config:                   # Component-specific config
      typescript: true
      tailwind: true

  - type: go-backend
    name: api-server
    enabled: true
    config:
      module: github.com/org/project
      framework: gin
      port: 8080

# Integration settings
integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
  shared_environment:
    NODE_ENV: development
    API_URL: http://localhost:8080

# Generation options
options:
  use_external_tools: true    # Use bootstrap tools
  dry_run: false              # Preview mode
  verbose: false              # Verbose output
  create_backup: true         # Backup existing files
  force_overwrite: false      # Force overwrite
```

### Supported Component Types

| Type | Description | Tool Required | Fallback |
|------|-------------|---------------|----------|
| `nextjs` | Next.js frontend | `npx` | No |
| `go-backend` | Go API server | `go` | No |
| `android` | Android app | `gradle` | Yes |
| `ios` | iOS app | `xcodebuild` | Yes |

## Examples

### Full-Stack Web App

```yaml
name: "fullstack-app"
output_dir: "./fullstack-app"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true

  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/fullstack-app
      framework: gin

integration:
  generate_docker_compose: true
  api_endpoints:
    backend: http://localhost:8080
```

### Frontend Only

```yaml
name: "frontend-app"
output_dir: "./frontend-app"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true

integration:
  generate_docker_compose: false
```

### Backend Only

```yaml
name: "api-service"
output_dir: "./api-service"

components:
  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/api-service
      framework: gin
      port: 8080

integration:
  generate_docker_compose: true
```

## Offline Mode

The generator supports offline operation:

### 1. Cache Tools (While Online)

```bash
generator cache-tools --save
```

This caches tool availability information.

### 2. Generate Offline

```bash
generator generate --config project.yaml --offline
```

**Note:** External bootstrap tools themselves must be installed before going offline. The cache only stores availability information.

## Troubleshooting

### Tool Not Found

**Problem:** `Error: Tool 'npx' not found`

**Solution:**

```bash
# Check which tools are missing
generator check-tools

# Install the missing tool
brew install node  # macOS
sudo apt install nodejs npm  # Ubuntu
```

### Tool Execution Failed

**Problem:** `Error: Tool execution failed`

**Solution:**

```bash
# Run with verbose output
generator generate --config project.yaml --verbose

# Check tool version
npx --version

# Try with fallback
generator generate --config project.yaml --no-external-tools
```

### Configuration Errors

**Problem:** `Error: Invalid configuration`

**Solution:**

```bash
# Generate a valid template
generator init-config template.yaml

# Validate with dry run
generator generate --config my-config.yaml --dry-run
```

### Permission Denied

**Problem:** `Error: Permission denied`

**Solution:**

```bash
# Fix permissions
chmod +x /usr/local/bin/generator

# Or install to user directory
mkdir -p ~/bin
cp generator ~/bin/
export PATH="$HOME/bin:$PATH"
```

## Next Steps

- **[CLI Commands](CLI_COMMANDS.md)** - Complete command reference
- **[Configuration Guide](CONFIGURATION.md)** - Detailed configuration options
- **[Examples](EXAMPLES.md)** - More example configurations
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions
- **[Architecture](ARCHITECTURE.md)** - How the system works

## Best Practices

1. **Check tools first** - Run `generator check-tools` before generating
2. **Use configuration files** - Save and version control your configs
3. **Test with dry run** - Preview with `--dry-run` before generating
4. **Cache for offline** - Use `generator cache-tools --save` for offline work
5. **Enable verbose mode** - Use `--verbose` when debugging issues

---

**Ready to generate your project?** Run `generator check-tools` to verify your environment, then `generator generate --config your-config.yaml`!
