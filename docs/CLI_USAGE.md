# CLI Usage Guide

## Overview

The Open Source Template Generator is a command-line tool that creates production-ready, enterprise-grade project structures following modern best practices. This guide covers all available commands and options.

## Installation

### From Source

```bash
git clone https://github.com/your-org/open-source-template-generator
cd open-source-template-generator
make setup
make build
```

### Using Go Install

```bash
go install github.com/open-source-template-generator/cmd/generator@latest
```

### Using Installation Script

```bash
curl -sSL https://raw.githubusercontent.com/your-org/open-source-template-generator/main/scripts/install.sh | bash
```

### Requirements

- **Go 1.23+** for building from source
- **Git** for cloning repositories
- **Make** for using the build system
- **Docker** (optional) for containerized usage

## Basic Usage

### Generate a New Project

The most common use case is generating a new project interactively:

```bash
generator generate
```

This will:

1. Prompt you for project details (name, organization, description, etc.)
2. Let you select which components to include
3. Show a configuration preview
4. Generate the complete project structure

### Generate with Configuration File

You can also generate projects using a pre-defined configuration file:

```bash
generator generate --config project-config.yaml --output ./my-project
```

### Dry Run Mode

Preview what would be generated without creating any files:

```bash
generator generate --dry-run
```

## Commands

### `generator generate`

Generates a new project from templates.

**Flags:**

- `--config, -c <file>`: Path to configuration file (YAML or JSON)
- `--output, -o <path>`: Output directory path
- `--dry-run`: Preview generation without creating files

**Examples:**

```bash
# Interactive generation
generator generate

# Generate with config file
generator generate --config my-config.yaml

# Generate to specific directory
generator generate --output /path/to/project

# Preview generation
generator generate --dry-run
```

### `generator version`

Shows version information for the generator and available packages.

**Flags:**

- `--packages`: Show latest package versions for all supported technologies
- `--check-updates`: Check for generator updates

**Examples:**

```bash
# Show generator version
generator version

# Show generator version and latest package versions
generator version --packages

# Check for updates
generator version --check-updates
```

### `generator config`

Manages generator configuration and defaults.

#### `generator config show`

Shows current configuration and default values.

```bash
generator config show
```

#### `generator config set`

Sets configuration values or loads from file.

```bash
# Set individual value (future feature)
generator config set license MIT

# Load from file
generator config set --file config.yaml
```

#### `generator config reset`

Resets configuration to defaults.

```bash
generator config reset
```

## Global Flags

These flags can be used with any command:

- `--verbose, -v`: Enable verbose logging
- `--quiet, -q`: Suppress non-error output
- `--log-level <level>`: Set log level (debug, info, warn, error)

## Configuration File Format

Configuration files can be in YAML or JSON format:

### YAML Example

```yaml
name: "my-awesome-project"
organization: "myorg"
description: "An awesome open source project"
license: "MIT"
author: "John Doe"
email: "john@example.com"
repository: "https://github.com/myorg/my-awesome-project"

components:
  frontend:
    main_app: true
    home: true
    admin: false
  backend:
    api: true
  mobile:
    android: true
    ios: true
  infrastructure:
    docker: true
    kubernetes: true
    terraform: false

output_path: "./my-awesome-project"
```

### JSON Example

```json
{
  "name": "my-awesome-project",
  "organization": "myorg",
  "description": "An awesome open source project",
  "license": "MIT",
  "components": {
    "frontend": {
      "main_app": true,
      "home": true,
      "admin": false
    },
    "backend": {
      "api": true
    },
    "mobile": {
      "android": true,
      "ios": true
    },
    "infrastructure": {
      "docker": true,
      "kubernetes": true,
      "terraform": false
    }
  },
  "output_path": "./my-awesome-project"
}
```

## Generated Project Structure

The generator creates a standardized, modern project structure following best practices:

```
my-awesome-project/
├── App/                    # Frontend applications (Next.js 15+, React 19+)
│   ├── main/              # Main application with TypeScript and Tailwind CSS
│   ├── home/              # Landing page optimized for performance
│   ├── admin/             # Admin dashboard with comprehensive UI components
│   └── shared-components/ # Reusable component library
├── CommonServer/          # Backend API server (Go 1.23+)
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

### Backend Components (Go 1.23+)

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

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure you have write permissions to the output directory
2. **Network Timeouts**: Use cached versions if network is unavailable
3. **Template Errors**: Check that template files are not corrupted
4. **Validation Failures**: Review configuration file format and required fields

### Getting Help

- Use `generator --help` for command overview
- Use `generator <command> --help` for command-specific help
- Check log files in `~/.cache/template-generator/logs/`
- Use `--verbose` flag for detailed output

## Examples

### Generate a Full-Stack Project

```bash
generator generate
# Select: Frontend Main App, Backend API, Infrastructure Docker
```

### Generate a Mobile-First Project

```bash
generator generate
# Select: Mobile Android, Mobile iOS, Backend API, Infrastructure Docker + K8s
```

### Generate from Configuration

```bash
# Create config.yaml with your preferences
generator generate --config config.yaml --output ./my-project
```

### Validate Generated Project

```bash
cd my-project
generator validate --verbose
```

This comprehensive CLI provides everything needed to generate, configure, and validate modern open source projects with best practices built-in.
