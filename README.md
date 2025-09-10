# Open Source Template Generator

A comprehensive CLI tool for generating production-ready, enterprise-grade open source project structures following modern best practices.

## Features

- **Multi-Platform Support**: Generate projects for frontend (Next.js), backend (Go), mobile (Android/iOS), and infrastructure
- **Latest Versions**: Automatically fetches and uses the latest stable versions of packages and frameworks
- **Complete CI/CD**: Includes GitHub Actions workflows, security scanning, and deployment configurations
- **Infrastructure as Code**: Terraform, Kubernetes, and Docker configurations included
- **Comprehensive Documentation**: Generates README, CONTRIBUTING, SECURITY, and API documentation
- **Interactive CLI**: User-friendly prompts for project configuration and component selection

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

- Go 1.22 or later
- Make

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
```

## Architecture

The generator follows a clean architecture pattern with dependency injection:

- **CLI Interface**: Handles user interaction and command processing
- **Template Engine**: Processes templates with variable substitution
- **Configuration Manager**: Manages project configuration and validation
- **File System Generator**: Creates directory structures and files
- **Version Manager**: Fetches latest package versions from registries
- **Validation Engine**: Validates generated project structures

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.