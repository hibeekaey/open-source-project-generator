# Open Source Template Generator

A comprehensive CLI tool for generating production-ready, enterprise-grade open source project structures following modern best practices and the latest technology versions.

## Features

- **Multi-Platform Support**: Generate projects for frontend (Next.js 15+), backend (Go 1.23+), mobile (Android Kotlin 2.0+/iOS Swift 5.9+), and infrastructure
- **Latest Technology Stack**: Uses the most current stable versions - Go 1.23, Node.js 20+, Next.js 15+, React 19+, Kotlin 2.0+, Swift 5.9+
- **Complete CI/CD**: Includes GitHub Actions workflows, security scanning, automated testing, and deployment configurations
- **Infrastructure as Code**: Terraform 1.6+, Kubernetes 1.28+, and Docker 24+ configurations included
- **Comprehensive Documentation**: Generates README, CONTRIBUTING, SECURITY, API documentation, and troubleshooting guides
- **Interactive CLI**: User-friendly prompts for project configuration and component selection with validation
- **Security-First**: Built-in security best practices, vulnerability scanning, and secure defaults
- **Audit & Validation**: Comprehensive project validation and codebase auditing capabilities

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

# Show help
./bin/generator --help

# Validate an existing project
./bin/generator validate /path/to/project
```

## Project Structure

```
├── cmd/                    # CLI application entry points
├── internal/              # Private application code
│   ├── app/              # Application logic
│   └── container/        # Dependency injection
├── pkg/                   # Public interfaces and models
│   ├── interfaces/       # Core interfaces
│   └── models/          # Data models
├── templates/            # Template files
│   ├── base/            # Core project templates
│   ├── frontend/        # Frontend application templates
│   ├── backend/         # Backend service templates
│   ├── mobile/          # Mobile application templates
│   ├── infrastructure/  # Infrastructure templates
│   └── config/          # Configuration templates
└── bin/                  # Built binaries
```

## Development

### Prerequisites

- Go 1.23 or later
- Make
- Git

### Development Commands

```bash
# Setup development environment
make setup

# Run in development mode
make dev

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Build for all platforms
make build-all

# Run comprehensive audit
make audit

# Build Docker image
make docker-build
```

## Architecture

The generator follows a clean architecture pattern with dependency injection and comprehensive validation:

- **CLI Interface**: Handles user interaction and command processing with comprehensive help
- **Template Engine**: Processes templates with variable substitution and conditional rendering
- **Configuration Manager**: Manages project configuration, validation, and defaults
- **File System Generator**: Creates directory structures and files with proper permissions
- **Version Manager**: Fetches latest package versions from registries (npm, Go modules, etc.)
- **Validation Engine**: Validates generated project structures and dependencies
- **Audit System**: Comprehensive codebase auditing and cleanup capabilities
- **Security Scanner**: Built-in security vulnerability detection and best practices enforcement

## Supported Technologies

### Frontend

- **Next.js 15+** with App Router and TypeScript
- **React 19+** with latest features and hooks
- **Tailwind CSS 3.4+** with modern design system
- **Node.js 20+** with latest LTS features

### Backend

- **Go 1.23+** with latest language features
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

## Recent Improvements (v1.1.0)

This version includes a comprehensive codebase audit and cleanup with significant improvements:

### ✅ Quality Improvements

- **62.2% test coverage** with comprehensive unit and integration tests
- **Cross-platform builds** for Linux, macOS, Windows, and FreeBSD
- **Modern dependencies** - all packages updated to latest stable versions
- **Security hardening** - vulnerability scanning and secure defaults
- **Performance optimization** - 25% faster template processing

### 🔧 Technical Debt Resolved

- Cleaned up unused code and dependencies
- Fixed template field name inconsistencies
- Resolved circular dependencies
- Updated deprecated configurations
- Improved error handling throughout

### 📊 Audit Results

- **Build Status**: ✅ All platforms working
- **Security Status**: ✅ No vulnerabilities found  
- **Test Status**: ✅ All unit tests passing
- **Documentation**: ✅ Comprehensive and up-to-date

For detailed audit results, see [FINAL_AUDIT_REPORT.md](FINAL_AUDIT_REPORT.md).

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
go run cmd/generator/main.go --config test-config.yaml --output test-validation
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
