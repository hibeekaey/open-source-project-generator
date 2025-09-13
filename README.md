# Open Source Template Generator

A focused CLI tool for generating production-ready, enterprise-grade open source project structures following modern best practices and the latest technology versions.

## Features

- **Multi-Platform Support**: Generate projects for frontend (Next.js 15+), backend (Go 1.24+), mobile (Android Kotlin 2.0+/iOS Swift 5.9+), and infrastructure
- **Latest Technology Stack**: Uses the most current stable versions - Go 1.24, Node.js 20+, Next.js 15+, React 19+, Kotlin 2.0+, Swift 5.9+
- **Complete CI/CD**: Includes GitHub Actions workflows, security scanning, automated testing, and deployment configurations
- **Infrastructure as Code**: Terraform 1.6+, Kubernetes 1.28+, and Docker 24+ configurations included
- **Comprehensive Documentation**: Generates README, CONTRIBUTING, SECURITY, API documentation, and troubleshooting guides
- **Interactive CLI**: User-friendly prompts for project configuration and component selection with validation
- **Security-First**: Built-in security best practices and secure defaults in generated templates
- **Project Validation**: Basic project validation to ensure generated structures are correct

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/open-source-template-generator
cd open-source-template-generator

# Install dependencies and build
make setup
make build
```

### Usage

```bash
# Generate a new project interactively
./bin/generator generate

# Generate with specific configuration
./bin/generator generate --config config/test-configs/test-config.yaml --output my-project

# Preview generation without creating files
./bin/generator generate --dry-run

# Show help
./bin/generator --help

# Show version information
./bin/generator version
```

## Project Structure

```
├── cmd/
│   └── generator/           # Main template generation command
├── internal/
│   ├── app/                # Core application logic
│   ├── config/             # Configuration management
│   └── container/          # Dependency injection
├── pkg/
│   ├── cli/                # Command-line interface
│   ├── filesystem/         # File operations
│   ├── interfaces/         # Core interfaces
│   ├── models/             # Data models
│   ├── template/           # Template processing engine
│   ├── validation/         # Validation logic
│   └── version/            # Version management
├── templates/              # Template files
│   ├── base/              # Core project templates
│   ├── frontend/          # Frontend application templates
│   ├── backend/           # Backend service templates
│   ├── mobile/            # Mobile application templates
│   └── infrastructure/    # Infrastructure templates
├── config/                 # Configuration files
├── docs/                   # Documentation
├── scripts/                # Build and validation scripts
└── bin/                    # Built binaries
```

## Development

### Prerequisites

- Go 1.24 or later
- Make
- Git

### Development Commands

```bash
# Setup development environment
make setup

# Run in development mode
make dev

# Run tests (unified test suite)
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Build for all platforms
make build-all

# Clean up project
make clean

# Build Docker image
make docker-build
```

## Architecture

The generator follows a clean architecture pattern with dependency injection and focused functionality:

- **CLI Interface**: Handles user interaction and command processing for template generation
- **Template Engine**: Processes templates with variable substitution and conditional rendering
- **Configuration Manager**: Manages project configuration, validation, and defaults
- **File System Generator**: Creates directory structures and files with proper permissions
- **Version Manager**: Fetches latest package versions from registries (npm, Go modules, etc.)
- **Validation Engine**: Validates generated project structures

## Supported Technologies

### Frontend

- **Next.js 15+** with App Router and TypeScript
- **React 19+** with latest features and hooks
- **Tailwind CSS 3.4+** with modern design system
- **Node.js 20+** with latest LTS features

### Backend

- **Go 1.24+** with latest language features
- **Gin Framework** for high-performance APIs
- **GORM** for database operations
- **JWT Authentication** with secure defaults
- **Redis** for caching and sessions

### Mobile

- **Android**: Kotlin 2.0+ with Jetpack Compose and Material Design 3
- **iOS**: Swift 5.9+ with SwiftUI and modern iOS patterns
- **Shared**: Common API specifications and design systems

### Infrastructure

- **Docker 24+** with multi-stage builds and security best practices
- **Kubernetes 1.28+** with proper resource management and security policies
- **Terraform 1.6+** for infrastructure as code
- **GitHub Actions** for CI/CD with comprehensive workflows

## Key Features

The generator provides a streamlined, focused experience for creating modern, production-ready projects with:

- **Single Purpose**: Dedicated solely to template generation with no auxiliary tools
- **Clean Architecture**: Well-organized code structure following Go best practices
- **Modern Technology Stack**: Uses the latest stable versions of popular frameworks and tools
- **Comprehensive Templates**: Production-ready templates with proper configuration and documentation
- **Simple CLI Interface**: Easy-to-use command-line interface for quick project generation
- **Lightweight**: Minimal overhead with fast build times and efficient resource usage

## Template Development

### Template Maintenance

For developers working with template files, comprehensive guidelines are available:

- **[Template Maintenance Guidelines](docs/TEMPLATE_MAINTENANCE.md)** - Complete guide for maintaining Go template files
- **[Quick Reference](docs/TEMPLATE_QUICK_REFERENCE.md)** - Essential commands and patterns for daily development

Key points for template development:

- All template files must have proper import statements for used packages
- Follow Go import organization conventions (standard library, third-party, local)
- Use validation tools before committing template changes
- Test template compilation with sample data

```bash
# Validate templates before committing
go run scripts/validate-templates/main.go --check-imports

# Generate test project to verify templates
go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-validation
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License.
