# CLI Usage Guide

## Overview

The Open Source Template Generator is a command-line tool that creates production-ready, enterprise-grade project structures following modern best practices. This guide covers all available commands and options.

## Installation

### From Source
```bash
git clone <repository-url>
cd open-source-template-generator
make build
```

### Using Go Install
```bash
go install github.com/open-source-template-generator/cmd/generator@latest
```

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

### `generator validate`

Validates a generated project structure and configuration.

**Usage:**
```bash
generator validate [project-path]
```

**Flags:**
- `--verbose, -v`: Show detailed validation output

**Examples:**
```bash
# Validate current directory
generator validate

# Validate specific project
generator validate /path/to/project

# Verbose validation
generator validate --verbose /path/to/project
```

### `generator version`

Shows version information for the generator and available packages.

**Flags:**
- `--packages`: Show latest package versions

**Examples:**
```bash
# Show generator version
generator version

# Show generator version and latest package versions
generator version --packages
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

## Project Structure

The generator creates a standardized project structure:

```
my-awesome-project/
├── App/                    # Frontend applications
│   ├── main/              # Main Next.js application
│   ├── home/              # Landing page
│   └── admin/             # Admin dashboard
├── CommonServer/          # Backend API server (Go)
├── Mobile/                # Mobile applications
│   ├── android/           # Android Kotlin app
│   ├── ios/               # iOS Swift app
│   └── shared/            # Shared mobile resources
├── Deploy/                # Infrastructure configurations
│   ├── docker/            # Docker configurations
│   ├── k8s/               # Kubernetes manifests
│   └── terraform/         # Terraform configurations
├── Docs/                  # Documentation
├── Scripts/               # Build and deployment scripts
├── .github/               # CI/CD workflows
├── Makefile               # Build system
├── docker-compose.yml     # Development environment
├── README.md              # Project documentation
├── CONTRIBUTING.md        # Contribution guidelines
└── LICENSE                # Project license
```

## Component Selection

### Frontend Components
- **Main App**: Core Next.js application with TypeScript and Tailwind CSS
- **Home**: Landing page application optimized for marketing
- **Admin**: Admin dashboard with forms, tables, and data management

### Backend Components
- **API**: Go API server with Gin framework, GORM, JWT auth, and Redis

### Mobile Components
- **Android**: Kotlin application with Jetpack Compose and Material Design 3
- **iOS**: Swift application with SwiftUI and proper MVVM architecture

### Infrastructure Components
- **Docker**: Multi-stage Dockerfiles and Docker Compose configurations
- **Kubernetes**: Complete K8s manifests with proper resource limits
- **Terraform**: Infrastructure as code for multi-cloud deployment

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