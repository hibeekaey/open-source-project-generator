# CLI Usage Guide

## Overview

The Open Source Project Generator is a comprehensive command-line tool that creates production-ready, enterprise-grade project structures following modern best practices and security standards. This guide covers all available commands, options, and advanced features.

## Installation

### From Source

```bash
git clone https://github.com/cuesoftinc/open-source-project-generator
cd open-source-project-generator
make setup
make build
```

### Using Go Install

```bash
go install github.com/cuesoftinc/open-source-template-generator/cmd/generator@latest
```

### Using Installation Script

```bash
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash
```

### Requirements

- **Go 1.21+** for building from source
- **Git** for cloning repositories
- **Make** for using the build system
- **Docker** (optional) for containerized usage
- **Node.js 20+** (optional) for frontend template validation

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

## Commands

### Core Commands

#### `generator generate`

Generates production-ready project structures from templates.

**Key Features:**

- Interactive project configuration with intelligent defaults
- Multiple generation modes (interactive, config file, template-specific, minimal)
- Comprehensive technology stack support
- Security-first configurations and best practices
- Offline mode with cached templates and versions

**Common Flags:**

- `--config, -c <file>`: Path to configuration file (YAML or JSON)
- `--output, -o <path>`: Output directory path
- `--template <name>`: Use specific template
- `--offline`: Use cached data without network requests
- `--minimal`: Generate minimal project structure
- `--dry-run`: Preview generation without creating files
- `--force`: Overwrite existing files
- `--non-interactive`: Run without user prompts (CI/CD mode)

**Examples:**

```bash
# Interactive generation (recommended for new users)
generator generate

# Generate from configuration file
generator generate --config project.yaml --output ./my-app

# Generate minimal Go API server
generator generate --minimal --template go-gin --output ./api

# Offline generation for air-gapped environments
generator generate --offline --template nextjs-app

# CI/CD automation
generator generate --non-interactive --config ci-config.yaml
```

#### `generator validate`

Validates project structure, configuration, and dependencies.

**Key Features:**

- Comprehensive project validation (structure, config, dependencies, security)
- Auto-fix capabilities for common issues
- Multiple report formats (text, JSON, HTML, markdown)
- CI/CD integration with proper exit codes
- Custom validation rules and severity filtering

**Common Flags:**

- `--fix`: Automatically fix common issues
- `--report`: Generate detailed validation report
- `--report-format <format>`: Report format (text, json, html, markdown)
- `--rules <rules>`: Specific validation rules to apply
- `--ignore-warnings`: Show only errors
- `--strict`: Use strict validation mode

**Examples:**

```bash
# Validate current directory
generator validate

# Validate and auto-fix issues
generator validate ./my-project --fix

# Generate HTML report
generator validate --report --report-format html --output-file report.html

# CI/CD validation
generator validate --non-interactive --output-format json
```

#### `generator audit`

Performs comprehensive security, quality, and compliance auditing.

**Key Features:**

- Security vulnerability scanning and policy compliance
- Code quality analysis and maintainability metrics
- License compliance checking and conflict detection
- Performance optimization recommendations
- Enterprise-grade reporting and scoring

**Common Flags:**

- `--security`: Perform security audit
- `--quality`: Perform code quality analysis
- `--licenses`: Perform license compliance check
- `--performance`: Perform performance analysis
- `--detailed`: Generate detailed audit report
- `--fail-on-high`: Fail if high severity issues found

**Examples:**

```bash
# Full comprehensive audit
generator audit

# Security-focused audit
generator audit --security --detailed

# Generate JSON report for automation
generator audit --output-format json --output-file audit.json

# CI/CD with failure conditions
generator audit --fail-on-high --min-score 7.5
```

#### `generator version`

Shows version information and checks for updates.

**Key Features:**

- Generator version and build information
- Latest package versions for all supported technologies
- Update checking with release notes
- Compatibility information and requirements
- Multiple output formats for automation

**Common Flags:**

- `--packages`: Show latest package versions
- `--check-updates`: Check for generator updates
- `--build-info`: Show detailed build information
- `--compatibility`: Show compatibility information

**Examples:**

```bash
# Show generator version
generator version

# Show all package versions
generator version --packages

# Check for updates with details
generator version --check-updates --verbose

# JSON output for automation
generator version --packages --output-format json
```

### Template Management

#### `generator list-templates`

Lists and discovers available project templates.

**Key Features:**

- Comprehensive template catalog with filtering
- Category and technology-based filtering
- Search capabilities across names and descriptions
- Detailed template information and compatibility

**Common Flags:**

- `--category <category>`: Filter by category (frontend, backend, mobile, infrastructure)
- `--technology <tech>`: Filter by technology (go, nodejs, react, etc.)
- `--search <query>`: Search templates by name or description
- `--tags <tags>`: Filter by tags
- `--detailed`: Show detailed template information

**Examples:**

```bash
# List all templates
generator list-templates

# List backend templates
generator list-templates --category backend

# Search for API templates
generator list-templates --search api

# Find React templates with TypeScript
generator list-templates --technology react --tags typescript
```

#### `generator template info`

Shows detailed information about specific templates.

**Examples:**

```bash
# Show template information
generator template info go-gin

# Show detailed information with all sections
generator template info nextjs-app --detailed --variables
```

#### `generator template validate`

Validates custom template structure and compliance.

**Examples:**

```bash
# Validate custom template
generator template validate ./my-custom-template

# Validate and auto-fix issues
generator template validate ./my-template --fix --detailed
```

### Configuration Management

#### `generator config show`

Displays current configuration values and sources.

**Examples:**

```bash
# Show all configuration
generator config show

# Show specific configuration key
generator config show default.license

# Show with source information
generator config show --sources --verbose
```

#### `generator config set`

Sets configuration values or loads from files.

**Examples:**

```bash
# Set individual values
generator config set default.license MIT
generator config set templates.path ./custom-templates

# Load from file
generator config set --file team-config.yaml --merge
```

#### `generator config edit`

Opens configuration files in editor for interactive editing.

**Examples:**

```bash
# Edit user configuration
generator config edit

# Edit with specific editor
generator config edit --editor code --backup
```

#### `generator config validate`

Validates configuration files and values.

**Examples:**

```bash
# Validate current configuration
generator config validate

# Validate specific file
generator config validate ./project-config.yaml --strict
```

#### `generator config export`

Exports configuration to shareable files.

**Examples:**

```bash
# Export to YAML
generator config export config.yaml

# Export as template
generator config export --template --format yaml project-template.yaml
```

### Update Management

#### `generator update`

Updates generator, templates, and package information.

**Key Features:**

- Generator binary updates with safety checks
- Template cache and package information updates
- Multiple update channels (stable, beta, alpha)
- Rollback support and compatibility checking

**Common Flags:**

- `--check`: Check for updates without installing
- `--install`: Install available updates
- `--templates`: Update template cache
- `--channel <channel>`: Update channel (stable, beta, alpha)
- `--force`: Force update even if risky

**Examples:**

```bash
# Check for updates
generator update --check

# Install updates safely
generator update --install --backup --verify

# Update templates only
generator update --templates

# Use beta channel
generator update --channel beta --install
```

### Cache Management

#### `generator cache`

Manages local cache for offline mode and performance.

**Subcommands:**

- `show`: Display cache statistics and information
- `clear`: Remove all cached data
- `clean`: Remove expired cache entries
- `validate`: Check cache integrity
- `repair`: Repair corrupted cache data

**Examples:**

```bash
# Show cache status
generator cache show

# Clean expired entries
generator cache clean

# Clear all cache
generator cache clear --force
```

### Logging and Debugging

#### `generator logs`

Views and analyzes application logs.

**Key Features:**

- Comprehensive log viewing with filtering
- Real-time log following (tail -f functionality)
- Multiple output formats and analysis capabilities
- Component and severity-based filtering

**Common Flags:**

- `--lines <n>`: Number of log lines to show
- `--level <level>`: Filter by log level
- `--component <component>`: Filter by component
- `--follow`: Follow logs in real-time
- `--since <time>`: Show logs since specific time

**Examples:**

```bash
# Show recent logs
generator logs

# Show error logs only
generator logs --level error --lines 100

# Follow logs in real-time
generator logs --follow --component template

# Show logs since specific time
generator logs --since "1h" --format json
```

## Global Flags

These flags can be used with any command:

**Logging and Output:**

- `--verbose, -v`: Enable verbose logging with detailed operation information
- `--quiet, -q`: Suppress non-error output (quiet mode)
- `--debug, -d`: Enable debug logging with performance metrics
- `--log-level <level>`: Set log level (debug, info, warn, error, fatal)
- `--log-json`: Output logs in JSON format
- `--log-caller`: Include caller information in logs

**Automation and Integration:**

- `--non-interactive`: Run in non-interactive mode (no prompts)
- `--output-format <format>`: Output format (text, json, yaml)

**Examples:**

```bash
# Verbose output for debugging
generator generate --verbose

# Quiet mode for automation
generator validate --quiet --output-format json

# Debug mode with performance metrics
generator audit --debug --detailed

# Non-interactive mode for CI/CD
generator generate --non-interactive --config ci.yaml
```

## Advanced Features

### Offline Mode

The generator supports complete offline operation using cached templates and package information:

```bash
# Enable offline mode globally
generator cache offline enable

# Generate project offline
generator generate --offline --template go-gin

# Check offline mode status
generator cache offline status
```

### Non-Interactive Mode (CI/CD)

Perfect for automation, CI/CD pipelines, and scripted environments:

```bash
# Set configuration via environment variables
export GENERATOR_PROJECT_NAME="my-api"
export GENERATOR_TEMPLATE="go-gin"
export GENERATOR_OUTPUT_PATH="./build"

# Generate without prompts
generator generate --non-interactive

# Validate with JSON output
generator validate --non-interactive --output-format json

# Audit with failure conditions
generator audit --non-interactive --fail-on-high --min-score 8.0
```

### Configuration Management

Comprehensive configuration system with multiple sources and formats:

```bash
# View configuration hierarchy
generator config show --sources

# Set team defaults
generator config set default.organization "MyCompany"
generator config set default.license "MIT"

# Export shareable configuration
generator config export --template team-defaults.yaml

# Validate configuration
generator config validate --strict
```

### Template Customization

Create and validate custom templates:

```bash
# Validate custom template
generator template validate ./my-custom-template --detailed

# Show template information
generator template info my-template --variables --dependencies

# List templates with filtering
generator list-templates --category backend --technology go
```

### Security and Compliance

Enterprise-grade security and compliance features:

```bash
# Security-focused audit
generator audit --security --detailed --fail-on-medium

# Validate with security rules
generator validate --rules security,compliance --strict

# Check for vulnerabilities
generator version --packages --check-security
```

## Configuration File Format

Configuration files support YAML, JSON, and TOML formats with environment variable substitution:

### YAML Configuration Example

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

### JSON Configuration Example

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

### Environment Variable Configuration

```bash
# Project configuration
export GENERATOR_PROJECT_NAME="my-api-server"
export GENERATOR_ORGANIZATION="mycompany"
export GENERATOR_DESCRIPTION="Production API server"
export GENERATOR_LICENSE="MIT"

# Template and output
export GENERATOR_TEMPLATE="go-gin"
export GENERATOR_OUTPUT_PATH="./api-server"

# Generation options
export GENERATOR_FORCE="false"
export GENERATOR_MINIMAL="false"
export GENERATOR_OFFLINE="true"
export GENERATOR_UPDATE_VERSIONS="false"

# Component selection
export GENERATOR_COMPONENTS_BACKEND_API="true"
export GENERATOR_COMPONENTS_BACKEND_AUTH="true"
export GENERATOR_COMPONENTS_INFRASTRUCTURE_DOCKER="true"

# Generate using environment variables
generator generate --non-interactive
```

## Generated Project Structure

The generator creates a standardized, modern project structure following best practices:

```text
my-awesome-project/
├── App/                    # Frontend applications (Next.js 15+, React 19+)
│   ├── main/              # Main application with TypeScript and Tailwind CSS
│   ├── home/              # Landing page optimized for performance
│   ├── admin/             # Admin dashboard with comprehensive UI components
│   └── shared-components/ # Reusable component library
├── CommonServer/          # Backend API server (Go 1.25+)
│   ├── cmd/               # Application entry points
│   ├── internal/          # Private application code
│   ├── pkg/               # Public interfaces and utilities
│   ├── migrations/        # Database migrations
│   └── docs/              # API documentation (Swagger/OpenAPI)
├── Mobile/                # Mobile applications
│   ├── android/           # Android Kotlin 2.0+ with Jetpack Compose
│   ├── ios/               # iOS Swift 5.9+ with SwiftUI
│   └── shared/            # Shared resources, API specs, design system
├── Deploy/                # Infrastructure configurations (latest versions)
│   ├── docker/            # Docker 24+ with multi-stage builds
│   ├── k8s/               # Kubernetes 1.28+ with security policies
│   ├── terraform/         # Terraform 1.6+ for infrastructure as code
│   └── monitoring/        # Prometheus, Grafana configurations
├── Docs/                  # Comprehensive documentation
│   ├── API.md             # API documentation
│   ├── DEPLOYMENT.md      # Deployment guide
│   ├── SECURITY_GUIDE.md  # Security best practices
│   └── USER_GUIDE.md      # User documentation
├── Scripts/               # Build and deployment automation
│   ├── build.sh           # Build scripts for all components
│   ├── deploy.sh          # Deployment automation
│   ├── test.sh            # Testing automation
│   └── setup.sh           # Development environment setup
├── .github/               # CI/CD workflows and templates
│   ├── workflows/         # GitHub Actions workflows
│   ├── ISSUE_TEMPLATE/    # Issue templates
│   └── PULL_REQUEST_TEMPLATE.md
├── Makefile               # Comprehensive build system
├── docker-compose.yml     # Development environment
├── README.md              # Project documentation
├── CONTRIBUTING.md        # Contribution guidelines
├── SECURITY.md            # Security policy
├── LICENSE                # Project license
└── .gitignore             # Git ignore patterns
```

## Component Selection

### Frontend Components (Node.js 20+, Next.js 15+, React 19+)

- **Main App**: Core Next.js application with TypeScript 5.7+, Tailwind CSS 3.4+, and comprehensive testing
- **Home**: Landing page optimized for performance and SEO with modern design patterns
- **Admin**: Admin dashboard with forms, tables, data management, and advanced UI components
- **Shared Components**: Reusable component library with proper TypeScript definitions

### Backend Components (Go 1.25+)

- **API Server**: Go API server with Gin framework, GORM, JWT authentication, Redis caching
- **Database**: PostgreSQL integration with migrations and proper connection pooling
- **Authentication**: JWT-based auth with refresh tokens and secure session management
- **Documentation**: Automatic Swagger/OpenAPI 3.0 documentation generation

### Mobile Components

- **Android**: Kotlin 2.0+ with Jetpack Compose, Material Design 3, and modern architecture
- **iOS**: Swift 5.9+ with SwiftUI, proper MVVM architecture, and iOS best practices
- **Shared Resources**: Common API specifications, design system, and shared assets

### Infrastructure Components (Latest Versions)

- **Docker**: Multi-stage Dockerfiles with security scanning and non-root users (Docker 24+)
- **Kubernetes**: Complete K8s manifests with resource limits, security policies, and health checks (K8s 1.28+)
- **Terraform**: Infrastructure as code for multi-cloud deployment with proper state management (Terraform 1.6+)
- **Monitoring**: Prometheus and Grafana configurations for comprehensive observability
- **CI/CD**: GitHub Actions workflows with security scanning, testing, and automated deployment

## Build System

The generated projects include a comprehensive Makefile with common commands:

```bash
# Setup development environment
make setup

# Start development servers
make dev

# Run tests
make test

# Build all components
make build

# Clean build artifacts
make clean

# Docker operations (if Docker component selected)
make docker-build
make docker-up

# Kubernetes operations (if K8s component selected)
make k8s-deploy
```

## Error Handling

The CLI provides comprehensive error handling with different error types:

- **Validation Errors**: Configuration or input validation issues
- **Template Errors**: Problems processing template files
- **File System Errors**: Issues creating directories or files
- **Network Errors**: Problems fetching package versions
- **Configuration Errors**: Issues with configuration files
- **Generation Errors**: Problems during project generation

Use `--verbose` flag to see detailed error information and stack traces.

## Logging

The generator creates log files in `~/.cache/template-generator/logs/` for debugging purposes. Log levels can be controlled with the `--log-level` flag.

## Best Practices

### Project Generation

1. **Start with Interactive Mode**: Use `generator generate` for your first project to understand options
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

## Troubleshooting

### Common Issues and Solutions

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

#### Performance Issues

```bash
# Check cache status and clean if needed
generator cache show
generator cache clean

# Use minimal generation for faster results
generator generate --minimal --skip-validation

# Enable debug mode to identify bottlenecks
generator generate --debug --verbose
```

### Getting Help and Support

#### Built-in Help

```bash
# General help
generator --help

# Command-specific help
generator generate --help
generator validate --help

# Show examples and usage patterns
generator <command> --help | grep -A 20 "Examples:"
```

#### Debugging and Diagnostics

```bash
# View recent logs
generator logs --level error --lines 50

# Debug mode with performance metrics
generator generate --debug --verbose

# Show configuration and sources
generator config show --sources --verbose

# Check system information
generator version --build-info --compatibility
```

#### Log Files and Diagnostics

- **Log Location**: `~/.cache/template-generator/logs/`
- **Configuration**: `~/.generator/config.yaml`
- **Cache Location**: `~/.cache/template-generator/cache/`

#### Community and Support

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Check online documentation for updates
- **Examples**: Review example configurations and templates
- **Community**: Join community discussions and forums

## Complete Examples

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
    monitoring: true
generate_options:
  update_versions: true
  include_examples: true
EOF

generator generate --config fullstack-config.yaml --output ./awesome-webapp
```

### Microservices API Backend

Generate a microservices-ready API backend:

```bash
generator generate --template go-gin --minimal --output ./api-service
cd ./api-service
generator validate --fix
generator audit --security --quality
```

### Mobile Application with Backend

Generate mobile applications with shared backend:

```bash
# Configuration for mobile-first project
cat > mobile-config.yaml << EOF
name: "mobile-app"
components:
  mobile:
    android: true
    ios: true
    shared: true
  backend:
    api: true
    auth: true
  infrastructure:
    docker: true
EOF

generator generate --config mobile-config.yaml
```

### CI/CD Pipeline Integration

Complete CI/CD pipeline example:

```yaml
# .github/workflows/generate-and-deploy.yml
name: Generate and Deploy
on:
  push:
    branches: [main]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install Generator
        run: |
          curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash
      
      - name: Generate Project
        run: |
          generator generate --non-interactive --config .generator/production.yaml --output ./build
      
      - name: Validate Project
        run: |
          generator validate ./build --non-interactive --output-format json > validation-results.json
      
      - name: Security Audit
        run: |
          generator audit ./build --security --fail-on-high --non-interactive
      
      - name: Deploy
        if: success()
        run: |
          # Deploy generated project
          cd ./build && make deploy
```

### Team Configuration Management

Set up team-wide defaults and standards:

```bash
# Create team configuration template
generator config export --template team-defaults.yaml

# Customize for team standards
cat >> team-defaults.yaml << EOF
default:
  organization: "MyCompany"
  license: "MIT"
  author: "MyCompany Team"
templates:
  preferred:
    - go-gin
    - nextjs-app
validation:
  strict: true
  rules:
    - security
    - quality
    - compliance
EOF

# Team members can use shared configuration
generator config set --file team-defaults.yaml --merge
generator generate --config project-specific.yaml
```

### Custom Template Development

Create and validate custom templates:

```bash
# Create custom template directory
mkdir -p ./my-custom-template/{templates,metadata}

# Validate template structure
generator template validate ./my-custom-template --detailed --fix

# Test template generation
generator generate --template ./my-custom-template --dry-run

# Share template with team
generator template validate ./my-custom-template --report --output-format html
```

### Offline Development Environment

Set up for air-gapped or offline development:

```bash
# Populate cache with all templates and versions
generator update --templates --packages --force

# Enable offline mode
generator cache offline enable

# Verify offline capability
generator generate --offline --template go-gin --dry-run

# Check cache status
generator cache show
```

### Monitoring and Maintenance

Regular maintenance and monitoring:

```bash
# Weekly maintenance script
#!/bin/bash
echo "Checking for updates..."
generator update --check

echo "Cleaning cache..."
generator cache clean

echo "Validating projects..."
find ~/projects -name ".generator" -type d | while read project; do
    echo "Validating $(dirname $project)..."
    generator validate "$(dirname $project)" --summary-only
done

echo "Security audit..."
generator audit ~/projects/production-app --security --detailed
```

This comprehensive CLI provides everything needed to generate, configure, validate, and maintain modern open source projects with enterprise-grade best practices built-in.
