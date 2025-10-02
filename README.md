# Open Source Project Generator

A comprehensive command-line tool that generates production-ready, enterprise-grade project structures following modern best practices and security standards.

## 🚀 Quick Start

### Installation

```bash
# Quick install (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash

# Using Go
go install github.com/cuesoftinc/open-source-template-generator/cmd/generator@latest

# Using Docker
docker run -it --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate
```

### Generate Your First Project

```bash
# Interactive mode (recommended for beginners)
generator generate

# Using a configuration file
generator generate --config project.yaml --output ./my-project

# Non-interactive mode (CI/CD)
GENERATOR_PROJECT_NAME=myapp generator generate --non-interactive
```

## ✨ Features

- **🎯 Interactive Project Configuration** - Guided setup with intelligent prompts
- **🏗️ Multi-Stack Support** - Frontend (Next.js, React), Backend (Go, Node.js), Mobile (Android, iOS), Infrastructure (Docker, K8s)
- **🔒 Security-First** - Built-in security best practices and vulnerability scanning
- **⚡ Offline Mode** - Generate projects without internet connectivity
- **🤖 CI/CD Ready** - Non-interactive mode for automation and pipelines
- **📦 Template Management** - Custom templates and validation
- **🔍 Project Validation** - Comprehensive validation and auditing
- **📊 Quality Assurance** - Code quality analysis and compliance checking
- **🧩 Modular Architecture** - Clean, maintainable codebase with focused components
- **🚀 High Performance** - Optimized for speed with comprehensive test coverage

## 🏗️ Generated Project Structure

The generator creates a standardized, modern project structure:

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
├── Scripts/               # Build and deployment automation
├── .github/workflows/     # CI/CD pipelines
├── Makefile              # Build system
└── docker-compose.yml     # Development environment
```

## 📖 Usage Examples

### Full-Stack Web Application

```bash
# Interactive generation
generator generate
# Select: Frontend (Next.js), Backend (Go), Infrastructure (Docker + K8s)

# Or using configuration
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
EOF

generator generate --config fullstack-config.yaml --output ./awesome-webapp
```

### Mobile Application

```bash
# Mobile-first project
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

### CI/CD Pipeline Integration

```yaml
# .github/workflows/generate-and-deploy.yml
name: Generate and Deploy
on: [push]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Generate Project
        env:
          GENERATOR_PROJECT_NAME: ${{ github.event.repository.name }}
          GENERATOR_BACKEND: true
          GENERATOR_FRONTEND: true
        run: |
          generator generate --non-interactive
          
      - name: Validate Project
        run: |
          generator validate ./output --output-format json
          
      - name: Security Audit
        run: |
          generator audit ./output --security --fail-on-high
```

## 🛠️ Commands

### Core Commands

```bash
# Generate projects
generator generate                    # Interactive mode
generator generate --config file.yaml # From configuration
generator generate --non-interactive  # CI/CD mode

# Validate projects
generator validate ./my-project       # Basic validation
generator validate --fix             # Auto-fix issues
generator validate --report          # Generate detailed report

# Audit projects
generator audit ./my-project          # Security and quality audit
generator audit --security           # Security-focused audit
generator audit --quality            # Code quality analysis

# Template management
generator list-templates             # List available templates
generator template info go-gin       # Template details
generator template validate ./custom  # Validate custom templates

# Configuration
generator config show                # Show current configuration
generator config set key value       # Set configuration values
generator config export file.yaml    # Export configuration

# Version management
generator version                    # Show version info
generator version --packages         # Show package versions
generator update --check             # Check for updates
```

### Global Options

```bash
--verbose, -v          # Verbose output (conflicts with --quiet)
--quiet, -q            # Quiet mode (conflicts with --verbose, --debug)
--debug, -d            # Debug mode (conflicts with --quiet)
--non-interactive      # Non-interactive mode (conflicts with --interactive)
--interactive          # Interactive mode (conflicts with --non-interactive)
--output-format json   # JSON output for automation
--log-level debug      # Set log level
```

### Enhanced Flag Validation

The CLI includes intelligent flag conflict detection:

```bash
# ✅ Valid combinations
generator generate --verbose --interactive
generator generate --quiet --non-interactive
generator version --json --packages

# ❌ Invalid combinations (will show helpful error messages)
generator generate --verbose --quiet          # Output mode conflict
generator generate --interactive --non-interactive  # Mode conflict
generator generate --debug --quiet            # Debug conflicts with quiet

# 💡 The CLI provides specific suggestions for resolving conflicts
generator generate --verbose --quiet
# Error: Verbose and quiet modes are mutually exclusive
# Suggestion: Choose either verbose output for detailed information OR quiet mode for minimal output
# Examples: --verbose, --quiet, --debug (implies verbose)
```

### Enhanced JSON Output

The version command now provides structured JSON output with comprehensive information:

```bash
# Enhanced JSON output with build metadata
generator version --json
{
  "version": "1.0.0",
  "gitCommit": "abc123def",
  "buildTime": "2024-01-01T00:00:00Z",
  "goVersion": "go1.21.0",
  "platform": "darwin",
  "architecture": "arm64"
}

# Perfect for CI/CD automation and scripting
generator version --json --packages | jq '.version'
```

## ⚙️ Configuration

### Environment Variables

```bash
# Project configuration
export GENERATOR_PROJECT_NAME="my-project"
export GENERATOR_ORGANIZATION="my-org"
export GENERATOR_DESCRIPTION="My awesome project"
export GENERATOR_LICENSE="MIT"

# Generation options
export GENERATOR_FORCE=true
export GENERATOR_MINIMAL=false
export GENERATOR_OFFLINE=false
export GENERATOR_UPDATE_VERSIONS=true

# Component selection
export GENERATOR_FRONTEND=true
export GENERATOR_BACKEND=true
export GENERATOR_MOBILE=false
export GENERATOR_INFRASTRUCTURE=true
```

### Configuration File

```yaml
# project-config.yaml
name: "my-awesome-project"
organization: "mycompany"
description: "An awesome open source project"
license: "MIT"
author: "John Doe"
email: "john@example.com"
repository: "https://github.com/mycompany/my-awesome-project"

components:
  frontend:
    main_app: true
    home: true
    admin: false
  backend:
    api: true
    auth: true
  mobile:
    android: true
    ios: true
  infrastructure:
    docker: true
    kubernetes: true
    terraform: false

generate_options:
  force: false
  minimal: false
  offline: false
  update_versions: true
  include_examples: true

output_path: "./my-awesome-project"
```

## 🔧 Development

### Prerequisites

- Go 1.25+
- Git
- Make (optional)
- Docker (optional)

### Build from Source

```bash
# Clone repository
git clone https://github.com/cuesoftinc/open-source-project-generator
cd open-source-project-generator

# Install dependencies
go mod download

# Build binary
make build

# Run tests
make test

# Generate cross-platform binaries
make build-all
```

### Development Workflow

```bash
# Setup development environment
make setup

# Run tests
make test

# Build and test
make build && ./bin/generator --version

# Run with debug logging
./bin/generator generate --debug --verbose
```

### Architecture Highlights

The codebase follows a **modular architecture** with clear separation of concerns:

- **📦 Focused Packages**: Each package has a single, clear responsibility
- **🔌 Interface-Based Design**: Components communicate through well-defined interfaces
- **🧪 Comprehensive Testing**: 70%+ test coverage with focused unit and integration tests
- **⚡ Performance Optimized**: Regex compilation, string operations, and memory usage optimized
- **🛡️ Zero Linting Violations**: Clean code with comprehensive error handling

**Key Architectural Benefits**:
- **Maintainability**: No file exceeds 1,000 lines, making code easy to navigate and modify
- **Testability**: Interface-based design enables easy mocking and isolated testing
- **Extensibility**: New features can be added without modifying existing components
- **Performance**: Smaller files compile faster and enable better IDE performance

## 📚 Documentation

- **[Getting Started Guide](docs/GETTING_STARTED.md)** - Complete installation and usage guide
- **[Configuration Guide](docs/CONFIGURATION.md)** - Configuration management and customization
- **[Template Development](docs/TEMPLATE_DEVELOPMENT.md)** - Creating and maintaining templates
- **[API Reference](docs/API_REFERENCE.md)** - Developer API documentation
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Package Structure](docs/PACKAGE_STRUCTURE.md)** - Modular architecture documentation
- **[Migration Guide](docs/MIGRATION_GUIDE.md)** - Code splitting refactoring guide

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Fork and clone the repository
git clone https://github.com/your-username/open-source-project-generator
cd open-source-project-generator

# Create a feature branch
git checkout -b feature/amazing-feature

# Make your changes and test
make test

# Commit your changes
git commit -m "Add amazing feature"

# Push to your fork
git push origin feature/amazing-feature

# Create a Pull Request
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📖 [Documentation](https://github.com/cuesoftinc/open-source-project-generator/wiki)
- 🐛 [Issue Tracker](https://github.com/cuesoftinc/open-source-project-generator/issues)
- 💬 [Discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- 📧 [Email Support](mailto:support@generator.dev)

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [Cobra](https://github.com/spf13/cobra) for CLI
- Inspired by modern development practices
- Community feedback and contributions

---

**Ready to generate your next project?** Start with `generator generate` and follow the interactive prompts!
