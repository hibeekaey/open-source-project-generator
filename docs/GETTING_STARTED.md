# Getting Started Guide

This comprehensive guide will help you install, configure, and use the Open Source Project Generator to create production-ready projects.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [Advanced Features](#advanced-features)
- [Troubleshooting](#troubleshooting)

## Installation

### Quick Install (Recommended)

#### Linux and macOS

```bash
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash
```

#### Windows

Download the Windows binary from the [releases page](https://github.com/cuesoftinc/open-source-template-generator/releases) and follow the installation instructions included in the archive.

### Package Manager Installation

#### Debian/Ubuntu (APT)

```bash
wget https://github.com/cuesoftinc/open-source-template-generator/releases/latest/download/generator_VERSION_amd64.deb
sudo dpkg -i generator_VERSION_amd64.deb
sudo apt-get install -f
```

#### Red Hat/CentOS/Fedora (YUM/DNF)

```bash
# Using YUM
sudo yum install https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-VERSION-1.x86_64.rpm

# Using DNF
sudo dnf install https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-VERSION-1.x86_64.rpm
```

#### macOS (Homebrew)

```bash
brew tap cuesoftinc/tap
brew install generator
```

#### Windows (Chocolatey)

```powershell
choco install generator
```

### Manual Installation

#### Download Pre-built Binaries

1. Visit the [releases page](https://github.com/cuesoftinc/open-source-project-generator/releases)
2. Download the appropriate archive for your platform:
   - Linux: `generator-linux-amd64.tar.gz`
   - macOS (Intel): `generator-darwin-amd64.tar.gz`
   - macOS (Apple Silicon): `generator-darwin-arm64.tar.gz`
   - Windows: `generator-windows-amd64.zip`

#### Extract and Install

**Linux/macOS/FreeBSD:**

```bash
# Extract the archive
tar -xzf generator-linux-amd64.tar.gz

# Move to installation directory
sudo mv generator-linux-amd64/generator /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/generator

# Verify installation
generator --version
```

**Windows:**

1. Extract the ZIP file to a directory (e.g., `C:\Program Files\generator`)
2. Add the directory to your PATH environment variable
3. Open a new command prompt and run: `generator --help`

### Docker Installation

```bash
# Pull from GitHub Container Registry
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# Run interactively
docker run -it --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate

# One-time generation
docker run --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate --config /workspace/config.yaml
```

### Build from Source

#### Prerequisites

- Go 1.25 or later
- Git
- Make (optional)

#### Clone and Build

```bash
# Clone the repository
git clone https://github.com/cuesoftinc/open-source-project-generator.git
cd open-source-project-generator

# Install dependencies
go mod download

# Build the binary
go build -o bin/generator ./cmd/generator

# Or use Make
make build

# Install to system (optional)
sudo cp bin/generator /usr/local/bin/
```

## Quick Start

### Interactive Project Generation

The most common use case is generating a new project interactively:

```bash
generator generate
```

This will:

1. Prompt you for project details (name, organization, description, etc.)
2. Let you select which components to include (frontend, backend, mobile, infrastructure)
3. Show a comprehensive configuration preview
4. Generate the complete project structure with best practices

### Generate with Configuration File

You can also generate projects using a pre-defined configuration file:

```bash
generator generate --config project-config.yaml --output ./my-project
```

### Non-Interactive Mode (CI/CD)

For automation and CI/CD pipelines:

```bash
# Using environment variables
GENERATOR_PROJECT_NAME=myapp GENERATOR_TEMPLATE=go-gin generator generate --non-interactive

# Using configuration file
generator generate --config ci-config.yaml --non-interactive --output ./build
```

### Dry Run Mode

Preview what would be generated without creating any files:

```bash
generator generate --dry-run --config project.yaml
```

### Offline Mode

Generate projects using cached templates and versions:

```bash
generator generate --offline --template nextjs-app
```

## Configuration

### Environment Variables

Configure the generator using environment variables for non-interactive operation:

#### Project Configuration

```bash
# Required
export GENERATOR_PROJECT_NAME="my-awesome-project"

# Optional
export GENERATOR_PROJECT_ORGANIZATION="my-org"
export GENERATOR_PROJECT_DESCRIPTION="An awesome project"
export GENERATOR_PROJECT_LICENSE="MIT"
export GENERATOR_OUTPUT_PATH="./output"
```

#### Generation Options

```bash
# Generation behavior
export GENERATOR_FORCE=true                # Overwrite existing files
export GENERATOR_MINIMAL=false             # Generate full project structure
export GENERATOR_OFFLINE=false             # Use cached data only
export GENERATOR_UPDATE_VERSIONS=true      # Fetch latest package versions
export GENERATOR_SKIP_VALIDATION=false     # Skip configuration validation
export GENERATOR_BACKUP_EXISTING=true      # Backup existing files
export GENERATOR_INCLUDE_EXAMPLES=true     # Include example code
export GENERATOR_TEMPLATE="go-gin"         # Specific template to use
```

#### Component Selection

```bash
# Enable/disable components
export GENERATOR_FRONTEND=true
export GENERATOR_BACKEND=true
export GENERATOR_MOBILE=false
export GENERATOR_INFRASTRUCTURE=true

# Technology selection
export GENERATOR_FRONTEND_TECH="nextjs-app"
export GENERATOR_BACKEND_TECH="go-gin"
export GENERATOR_MOBILE_TECH="android-kotlin"
export GENERATOR_INFRASTRUCTURE_TECH="docker"
```

### Configuration File Format

Configuration files support YAML, JSON, and TOML formats with environment variable substitution:

#### YAML Configuration Example

```yaml
# Project metadata
name: "my-awesome-project"
organization: "myorg"
description: "An awesome open source project"
license: "MIT"
author: "John Doe"
email: "john@example.com"
repository: "https://github.com/myorg/my-awesome-project"

# Component selection
components:
  frontend:
    main_app: true
    home: true
    admin: false
    shared_components: true
  backend:
    api: true
    auth: true
    database: true
  mobile:
    android: true
    ios: true
    shared: true
  infrastructure:
    docker: true
    kubernetes: true
    terraform: false
    monitoring: true

# Generation options
generate_options:
  force: false
  minimal: false
  offline: false
  update_versions: true
  skip_validation: false
  backup_existing: true
  include_examples: true

# Version overrides
versions:
  go: "1.21"
  node: "20"
  react: "19"
  typescript: "5"

# Custom variables
custom_vars:
  database_name: "myapp_db"
  api_port: "8080"
  frontend_port: "3000"

# Output configuration
output_path: "./my-awesome-project"
```

#### JSON Configuration Example

```json
{
  "name": "my-awesome-project",
  "organization": "myorg",
  "description": "An awesome open source project",
  "license": "MIT",
  "author": "John Doe",
  "email": "john@example.com",
  "repository": "https://github.com/myorg/my-awesome-project",
  "components": {
    "frontend": {
      "main_app": true,
      "home": true,
      "admin": false,
      "shared_components": true
    },
    "backend": {
      "api": true,
      "auth": true,
      "database": true
    },
    "mobile": {
      "android": true,
      "ios": true,
      "shared": true
    },
    "infrastructure": {
      "docker": true,
      "kubernetes": true,
      "terraform": false,
      "monitoring": true
    }
  },
  "generate_options": {
    "force": false,
    "minimal": false,
    "offline": false,
    "update_versions": true,
    "skip_validation": false,
    "backup_existing": true,
    "include_examples": true
  },
  "versions": {
    "go": "1.21",
    "node": "20",
    "react": "19",
    "typescript": "5"
  },
  "custom_vars": {
    "database_name": "myapp_db",
    "api_port": "8080",
    "frontend_port": "3000"
  },
  "output_path": "./my-awesome-project"
}
```

## Usage Examples

### Full-Stack Web Application

Generate a complete web application with frontend, backend, and infrastructure:

```bash
# Interactive generation
generator generate
# Select: Frontend (Next.js + React), Backend (Go + Gin), Infrastructure (Docker + K8s)

# Or using configuration file
cat > fullstack-config.yaml << EOF
name: "awesome-webapp"
organization: "mycompany"
description: "Full-stack web application"
license: "MIT"
components:
  frontend:
    main_app: true
    admin: true
  backend:
    api: true
    auth: true
  infrastructure:
    docker: true
    kubernetes: true
generate_options:
  update_versions: true
  include_examples: true
EOF

generator generate --config fullstack-config.yaml --output ./awesome-webapp
```

**Generated Structure:**

```text
awesome-webapp/
├── App/
│   ├── main/              # Next.js main application
│   └── home/              # Landing page
├── CommonServer/          # Go API server
├── Deploy/
│   ├── docker/            # Docker configurations
│   └── k8s/               # Kubernetes manifests
├── .github/workflows/     # CI/CD pipelines
├── Makefile
└── docker-compose.yml
```

### Mobile Application

Generate a mobile application with backend API:

```bash
# Configuration for mobile-first project
cat > mobile-config.yaml << EOF
name: "mobile-app"
components:
  mobile:
    android: true
    ios: true
  backend:
    api: true
  infrastructure:
    docker: true
EOF

generator generate --config mobile-config.yaml
```

**Generated Structure:**

```text
mobile-app/
├── CommonServer/          # Go API server
├── Mobile/
│   ├── android/           # Kotlin Android app
│   ├── ios/               # Swift iOS app
│   └── shared/            # Shared resources
├── Deploy/docker/         # Docker configurations
└── .github/workflows/     # Mobile CI/CD
```

### Microservices Architecture

Generate multiple related services:

```bash
# Generate API Gateway
generator generate --config api-gateway.yaml

# Generate User Service
generator generate --config user-service.yaml

# Generate Notification Service
generator generate --config notification-service.yaml
```

### Enterprise SaaS Platform

Generate a comprehensive enterprise platform:

```bash
generator generate --dry-run  # Preview first
generator generate --config enterprise-config.yaml
```

## Advanced Features

### Interactive UI Framework

The generator includes a comprehensive interactive UI framework:

```bash
# Interactive mode with advanced features
generator generate --interactive

# Interactive with specific output directory
generator generate --interactive --output ./my-project

# Interactive with template pre-selection
generator generate --interactive --template go-gin
```

**Features:**

- Guided project configuration with intelligent prompts
- Multi-select interface for choosing project components
- Real-time validation with helpful error messages
- Progress tracking during project generation
- Context-sensitive help system

### Template Management

```bash
# List available templates
generator list-templates

# Search templates
generator list-templates --search "api"

# Show template details
generator template info go-gin

# Validate custom templates
generator template validate ./my-custom-template
```

### Project Validation and Auditing

```bash
# Basic validation
generator validate ./my-project

# Detailed validation with auto-fix
generator validate --fix --report

# Security audit
generator audit --security --detailed

# Quality analysis
generator audit --quality --min-score 8.0
```

### Configuration Management

```bash
# Show current configuration
generator config show

# Set configuration values
generator config set default.license MIT
generator config set default.organization "MyCompany"

# Export configuration
generator config export team-defaults.yaml

# Import configuration
generator config import --file team-defaults.yaml
```

### Version Management

```bash
# Show version information
generator version

# Show latest package versions
generator version --packages

# Check for updates
generator update --check

# Install updates
generator update --install
```

### Cache Management

```bash
# Show cache status
generator cache show

# Clean expired entries
generator cache clean

# Clear all cache
generator cache clear --force

# Enable offline mode
generator cache offline enable
```

## Troubleshooting

### Common Issues

#### Permission Denied Errors

```bash
# Check output directory permissions
ls -la ./output-directory

# Create directory with proper permissions
mkdir -p ./my-project && chmod 755 ./my-project

# Use different output directory
generator generate --output ~/projects/my-project
```

#### Network and Connectivity Issues

```bash
# Use offline mode
generator generate --offline --template go-gin

# Check cache status
generator cache show

# Populate cache for offline use
generator update --templates --packages
```

#### Template and Configuration Errors

```bash
# Validate configuration before generation
generator config validate ./my-config.yaml

# Check template information
generator template info go-gin --detailed

# Validate custom templates
generator template validate ./my-custom-template --fix
```

#### Validation and Audit Failures

```bash
# Run validation with detailed output
generator validate --verbose --detailed

# Fix common issues automatically
generator validate --fix --backup

# Check specific validation rules
generator validate --rules structure,dependencies --show-fixes
```

### Debug Mode

Enable debug output for troubleshooting:

```bash
# Debug mode with performance metrics
generator generate --debug --verbose

# Check logs
generator logs --level error --lines 50

# Show configuration and sources
generator config show --sources --verbose
```

### Getting Help

```bash
# General help
generator --help

# Command-specific help
generator generate --help
generator validate --help

# Show examples and usage patterns
generator <command> --help | grep -A 20 "Examples:"
```

### Log Files and Diagnostics

- **Log Location**: `~/.cache/template-generator/logs/`
- **Configuration**: `~/.generator/config.yaml`
- **Cache Location**: `~/.cache/template-generator/cache/`

## Best Practices

### Project Generation

1. **Start with Interactive Mode**: Use `generator generate` for your first project
2. **Use Configuration Files**: Save and version control your project configurations
3. **Enable Validation**: Always validate generated projects with `generator validate`
4. **Regular Audits**: Run `generator audit` periodically to check security and quality
5. **Keep Updated**: Regularly check for updates with `generator update --check`

### CI/CD Integration

```bash
# Example CI/CD pipeline step
- name: Generate Project
  run: |
    generator generate --non-interactive --config .generator/ci-config.yaml
    generator validate --non-interactive --output-format json > validation-results.json
    generator audit --non-interactive --fail-on-high --min-score 8.0
```

### Team Collaboration

1. **Shared Configuration**: Use `generator config export` to create team templates
2. **Custom Templates**: Validate custom templates with `generator template validate`
3. **Documentation**: Include generator configuration in project documentation
4. **Standards**: Establish team standards for component selection and configuration

### Security Considerations

1. **Regular Audits**: Use `generator audit --security` for security scanning
2. **Dependency Updates**: Enable `--update-versions` for latest secure versions
3. **Validation**: Use strict validation with `generator validate --strict`
4. **Offline Mode**: Use offline mode in secure environments

## Next Steps

- **[Configuration Guide](CONFIGURATION.md)** - Advanced configuration options
- **[Template Development](TEMPLATE_DEVELOPMENT.md)** - Creating custom templates
- **[API Reference](API_REFERENCE.md)** - Developer API documentation
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions

---

**Ready to generate your next project?** Start with `generator generate` and follow the interactive prompts!
