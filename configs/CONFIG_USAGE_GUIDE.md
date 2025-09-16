# Configuration Usage Guide

This guide explains how to use the configuration files for the Open Source Project Generator to achieve full usage of all available features.

## Configuration Files Overview

### Available Configuration Examples

1. **`config-full-usage.yaml`** - Demonstrates all available features and components
2. **`config-full-usage.json`** - Same as above in JSON format
3. **`config-minimal.yaml`** - Minimal configuration for simple projects
4. **`config-frontend-only.yaml`** - Frontend-focused applications
5. **`config-mobile-focused.yaml`** - Mobile applications with backend API
6. **`config-enterprise.yaml`** - Enterprise-grade full-stack platform

## Using Configuration Files

### Command Line Usage

```bash
# Generate project using YAML configuration
./generator generate --config config-full-usage.yaml

# Generate project using JSON configuration
./generator generate --config config-full-usage.json

# Override output directory
./generator generate --config config-enterprise.yaml --output /path/to/output

# Preview what would be generated (dry run)
./generator generate --config config-full-usage.yaml --dry-run
```

### Interactive Mode

```bash
# Start interactive configuration (no config file needed)
./generator generate

# The CLI will prompt you for all configuration options
```

## Configuration Structure

### Basic Project Information

```yaml
name: "project-name"              # Required: Project name
organization: "org-name"          # Required: Organization/company name
description: "Project description" # Optional: Project description
license: "MIT"                    # Optional: License type (MIT, Apache-2.0, GPL-3.0, BSD-3-Clause)
author: "Your Name"               # Optional: Author name
email: "your@email.com"           # Optional: Contact email
repository: "https://github.com/..." # Optional: Repository URL
output_path: "output/path"        # Required: Where to generate the project
```

### Component Selection

#### Frontend Components

```yaml
components:
  frontend:
    nextjs:
      app: true     # Main Next.js application
      home: true    # Landing page application
      admin: true   # Admin dashboard
      shared: true  # Shared components library
```

**Generated Structure:**

- `App/main/` - Main Next.js application
- `App/home/` - Landing page application  
- `App/admin/` - Admin dashboard
- `App/shared-components/` - Reusable components

#### Backend Components

```yaml
components:
  backend:
    go_gin: true    # Go Gin REST API server
```

**Generated Structure:**

- `CommonServer/` - Go Gin backend with clean architecture

#### Mobile Components

```yaml
components:
  mobile:
    android: true   # Android Kotlin application
    ios: true       # iOS Swift application
```

**Generated Structure:**

- `Mobile/android/` - Android Kotlin application
- `Mobile/ios/` - iOS Swift application
- `Mobile/shared/` - Shared mobile resources

#### Infrastructure Components

```yaml
components:
  infrastructure:
    docker: true      # Docker containerization
    kubernetes: true  # Kubernetes deployment manifests
    terraform: true   # Infrastructure as Code
```

**Generated Structure:**

- `Deploy/docker/` - Docker configurations
- `Deploy/k8s/` - Kubernetes manifests
- `Deploy/terraform/` - Terraform infrastructure code
- `Deploy/monitoring/` - Monitoring configurations

### Version Configuration

```yaml
versions:
  # Language Runtime Versions
  node: "20.11.0"
  go: "1.22.0"
  
  # Package Versions
  packages:
    "react": "18.2.0"
    "next": "15.0.0"
    "typescript": "5.3.0"
    # ... more packages
```

## Generated Project Structure

When using full configuration, the generator creates:

```text
my-awesome-project/
├── App/                          # Frontend Applications
│   ├── main/                     # Main Next.js app
│   ├── home/                     # Landing page
│   ├── admin/                    # Admin dashboard
│   └── shared-components/        # Shared components
├── CommonServer/                 # Backend API
│   ├── cmd/                      # Application entry points
│   ├── internal/                 # Private application code
│   └── pkg/                      # Public interfaces
├── Mobile/                       # Mobile Applications
│   ├── android/                  # Android Kotlin app
│   ├── ios/                      # iOS Swift app
│   └── shared/                   # Shared resources
├── Deploy/                       # Infrastructure
│   ├── docker/                   # Docker configurations
│   ├── k8s/                      # Kubernetes manifests
│   ├── terraform/                # Infrastructure as code
│   └── monitoring/               # Monitoring setup
├── Docs/                         # Documentation
├── Scripts/                      # Build automation
├── .github/                      # CI/CD workflows
├── Makefile                      # Build system
├── docker-compose.yml            # Development environment
├── README.md                     # Project documentation
├── CONTRIBUTING.md               # Contribution guidelines
├── SECURITY.md                   # Security policy
├── LICENSE                       # Project license
└── .gitignore                    # Git ignore patterns
```

## Use Case Examples

### 1. Full-Stack Web Application

Use `config-full-usage.yaml` or `config-enterprise.yaml` for:

- Web applications with multiple frontends
- REST API backend
- Complete infrastructure setup
- CI/CD pipelines

### 2. Frontend-Only Project

Use `config-frontend-only.yaml` for:

- Static sites
- Single-page applications
- Frontend microservices
- Documentation sites

### 3. Mobile Application

Use `config-mobile-focused.yaml` for:

- Cross-platform mobile apps
- Mobile-first applications
- API-driven mobile solutions

### 4. Microservice

Use `config-minimal.yaml` for:

- Single-purpose services
- API-only applications
- Lightweight deployments

## Advanced Configuration

### Custom Package Versions

```yaml
versions:
  packages:
    # Override specific package versions
    "react": "18.3.0-beta"
    "custom-package": "1.0.0"
```

### Environment-Specific Configurations

Create multiple config files for different environments:

```bash
# Development
./generator generate --config config-dev.yaml

# Staging  
./generator generate --config config-staging.yaml

# Production
./generator generate --config config-prod.yaml
```

## Validation and Testing

### Validate Configuration

```bash
# Validate generated project
./generator validate /path/to/generated/project

# Verbose validation output
./generator validate /path/to/generated/project --verbose
```

### Configuration Management

```bash
# Show current configuration defaults
./generator config show

# Set configuration from file
./generator config set --file config-full-usage.yaml

# Reset to defaults
./generator config reset
```

## Best Practices

1. **Start with Examples**: Use provided example configurations as templates
2. **Version Pinning**: Specify exact versions for reproducible builds
3. **Component Selection**: Only enable components you actually need
4. **Validation**: Always validate generated projects before deployment
5. **Documentation**: Document any custom configuration choices
6. **Version Control**: Store configuration files in version control

## Troubleshooting

### Common Issues

1. **Invalid Configuration**: Use `--dry-run` to preview before generation
2. **Missing Components**: Check that all required components are enabled
3. **Version Conflicts**: Verify package version compatibility
4. **Path Issues**: Ensure output paths are writable and valid

### Getting Help

```bash
# Show version and help
./generator version

# Show command help
./generator generate --help
./generator config --help
./generator validate --help
```

## Configuration Schema

For IDE support and validation, the configuration follows this schema:

- **Required Fields**: `name`, `organization`, `output_path`
- **Optional Fields**: All other fields have sensible defaults
- **Supported Formats**: YAML (`.yaml`, `.yml`) and JSON (`.json`)
- **Validation**: Automatic validation on load with helpful error messages

This comprehensive configuration system allows you to generate everything from simple prototypes to enterprise-grade applications with full infrastructure setup.
