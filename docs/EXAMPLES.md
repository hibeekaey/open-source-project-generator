# Usage Examples

This document provides practical examples of using the Open Source Template Generator for different types of projects.

## Example 1: Full-Stack Web Application

Generate a complete web application with frontend, backend, and infrastructure:

```bash
# Interactive generation
generator generate

# When prompted, select:
# - Project name: "awesome-webapp"
# - Organization: "mycompany"
# - Description: "A full-stack web application"
# - License: "MIT"
# - Components:
#   ✓ Frontend Main App
#   ✓ Frontend Home (landing page)
#   ✓ Backend API
#   ✓ Infrastructure Docker
#   ✓ Infrastructure Kubernetes
```

**Generated Structure:**

```
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

**Getting Started:**

```bash
cd awesome-webapp
make setup
make dev
```

## Example 2: Mobile-First Application

Generate a mobile application with backend API:

```bash
generator generate --config mobile-config.yaml
```

**mobile-config.yaml:**

```yaml
name: "mobile-app"
organization: "startup"
description: "Cross-platform mobile application"
license: "MIT"
author: "Jane Developer"
email: "jane@startup.com"

components:
  frontend:
    main_app: false
    home: false
    admin: false
  backend:
    api: true
  mobile:
    android: true
    ios: true
  infrastructure:
    docker: true
    kubernetes: false
    terraform: false

output_path: "./mobile-app"
```

**Generated Structure:**

```
mobile-app/
├── CommonServer/          # Go API server
├── Mobile/
│   ├── android/           # Kotlin Android app
│   ├── ios/               # Swift iOS app
│   └── shared/            # Shared resources
├── Deploy/docker/         # Docker configurations
└── .github/workflows/     # Mobile CI/CD
```

## Example 3: Enterprise SaaS Platform

Generate a comprehensive enterprise platform:

```bash
generator generate --dry-run  # Preview first
generator generate --config enterprise-config.yaml
```

**enterprise-config.yaml:**

```yaml
name: "enterprise-platform"
organization: "enterprise-corp"
description: "Enterprise SaaS platform with admin dashboard"
license: "Apache-2.0"
author: "Enterprise Team"
email: "team@enterprise-corp.com"
repository: "https://github.com/enterprise-corp/platform"

components:
  frontend:
    main_app: true
    home: true
    admin: true
  backend:
    api: true
  mobile:
    android: true
    ios: true
  infrastructure:
    docker: true
    kubernetes: true
    terraform: true

output_path: "./enterprise-platform"
```

**Generated Structure:**

```
enterprise-platform/
├── App/
│   ├── main/              # Main SaaS application
│   ├── home/              # Marketing website
│   ├── admin/             # Admin dashboard
│   └── shared/            # Shared components
├── CommonServer/          # Go API server
├── Mobile/
│   ├── android/           # Android app
│   ├── ios/               # iOS app
│   └── shared/            # Mobile shared resources
├── Deploy/
│   ├── docker/            # Docker configurations
│   ├── k8s/               # Kubernetes manifests
│   └── terraform/         # Infrastructure as code
├── Docs/                  # Documentation
├── Scripts/               # Automation scripts
└── .github/workflows/     # Complete CI/CD
```

## Example 4: Open Source Library

Generate a simple library project:

```bash
generator generate
# Select minimal components:
# - Frontend Main App (for documentation/demo)
# - Infrastructure Docker (for development)
```

## Example 5: Microservices Architecture

Generate multiple related services:

```bash
# Generate API Gateway
generator generate --config api-gateway.yaml

# Generate User Service
generator generate --config user-service.yaml

# Generate Notification Service
generator generate --config notification-service.yaml
```

**api-gateway.yaml:**

```yaml
name: "api-gateway"
organization: "microservices-corp"
description: "API Gateway for microservices architecture"
license: "MIT"

components:
  frontend:
    main_app: false
    home: false
    admin: false
  backend:
    api: true
  mobile:
    android: false
    ios: false
  infrastructure:
    docker: true
    kubernetes: true
    terraform: false

output_path: "./services/api-gateway"
```

## Example 6: Validation and Testing

After generating any project, validate its structure:

```bash
# Basic validation
generator validate ./my-project

# Detailed validation with verbose output
generator validate --verbose ./my-project

# Validate current directory
cd my-project
generator validate
```

**Example Validation Output:**

```
⏳ Validating project at ./my-project...
✅ Project validation completed successfully

Validation Warnings:
  ⚠️  Dependencies: Found 2 package.json files - ensure dependency versions are compatible

Validation Summary:
  Valid: true
  Errors: 0
  Warnings: 1
```

## Example 7: Configuration Management

Manage generator configuration:

```bash
# Show current configuration
generator config show

# Load configuration from file
generator config set --file my-defaults.yaml

# Reset to defaults
generator config reset
```

## Example 8: Version Management

Check latest package versions:

```bash
# Show generator version
generator version

# Show latest package versions
generator version --packages
```

**Example Output:**

```
Open Source Template Generator v1.0.0
Built with Go 1.22+

⏳ Fetching latest package versions...

Latest Package Versions:
  Node.js: 20.11.0
  Go: 1.22.0
  Next.js: 15.0.0
  React: 18.2.0
  Kotlin: 2.0.0
  Swift: 5.9.0

Common Packages:
  typescript: 5.3.0
  tailwindcss: 3.4.0
  eslint: 8.56.0
```

## Example 9: Development Workflow

Complete development workflow after generation:

```bash
# Generate project
generator generate --config my-config.yaml

# Navigate to project
cd my-awesome-project

# Setup development environment
make setup

# Start development servers
make dev

# In another terminal, run tests
make test

# Build for production
make build

# Deploy with Docker
make docker-build
make docker-up

# Deploy to Kubernetes (if configured)
make k8s-deploy
```

## Example 10: Troubleshooting

Debug issues with verbose logging:

```bash
# Generate with verbose output
generator generate --verbose

# Generate with debug logging
generator --log-level debug generate

# Quiet mode (errors only)
generator --quiet generate

# Check logs
tail -f ~/.cache/template-generator/logs/generator-$(date +%Y-%m-%d).log
```

## Configuration Templates

### Minimal Configuration

```yaml
name: "simple-project"
organization: "developer"
description: "Simple project"
license: "MIT"
components:
  frontend:
    main_app: true
  infrastructure:
    docker: true
```

### Full-Featured Configuration

```yaml
name: "full-project"
organization: "company"
description: "Full-featured project with all components"
license: "Apache-2.0"
author: "Development Team"
email: "dev@company.com"
repository: "https://github.com/company/full-project"

components:
  frontend:
    main_app: true
    home: true
    admin: true
  backend:
    api: true
  mobile:
    android: true
    ios: true
  infrastructure:
    docker: true
    kubernetes: true
    terraform: true

custom_vars:
  database_type: "postgresql"
  cache_type: "redis"
  monitoring: "prometheus"

output_path: "./full-project"
```

## Best Practices

1. **Always use dry-run first** for complex configurations
2. **Validate generated projects** before committing
3. **Use configuration files** for repeatable generations
4. **Check latest versions** regularly with `generator version --packages`
5. **Keep configurations in version control** for team consistency
6. **Use verbose logging** when troubleshooting issues

These examples demonstrate the flexibility and power of the Open Source Template Generator for various project types and development workflows.
