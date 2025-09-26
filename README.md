# Open Source Project Generator

A comprehensive CLI tool for generating production-ready project structures with modern technology stacks.

## Features

- **Multi-Platform Support**: Generate projects for frontend (Next.js), backend (Go), mobile (Android/iOS), and infrastructure
- **Modern Technology Stack**: Uses current stable versions - Go 1.21+, Node.js 20+, Next.js 13+, React 18+
- **Basic Infrastructure**: Includes Docker, Kubernetes, and Terraform configurations
- **Essential Documentation**: Generates README, CONTRIBUTING, and basic documentation
- **Interactive CLI**: Intuitive prompts for project configuration and component selection
- **Basic Validation**: Ensures generated project structures are correct

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/cuesoftinc/open-source-project-generator
cd open-source-project-generator

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

# Show version information
./bin/generator version
```

## Project Structure

```text
open-source-project-generator/
├── cmd/generator/          # Main application entry point
├── internal/               # Internal application components
│   ├── app/               # Core application logic
│   ├── config/            # Configuration management
│   └── container/         # Dependency injection
├── pkg/                   # Public packages
│   ├── cli/               # CLI interface
│   ├── filesystem/        # File system operations
│   ├── template/          # Template processing
│   ├── validation/        # Basic validation
│   ├── version/           # Version management
│   └── models/            # Data models
├── pkg/template/templates/ # Embedded project templates
│   ├── base/              # Base templates (embedded)
│   ├── frontend/          # Frontend templates (embedded)
│   ├── backend/           # Backend templates (embedded)
│   ├── mobile/            # Mobile templates (embedded)
│   └── infrastructure/    # Infrastructure templates (embedded)
└── scripts/               # Build and utility scripts
```

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Building

```bash
# Install dependencies
make install

# Build the application
make build

# Run tests
make test

# Run in development mode
make dev
```

### Available Commands

- `make build` - Build the generator binary
- `make test` - Run all tests
- `make clean` - Clean build artifacts
- `make run` - Build and run the generator
- `make dev` - Run in development mode
- `make lint` - Run code linter
- `make fmt` - Format code
- `make vet` - Run go vet

## Templates

The generator includes embedded templates for:

- **Frontend**: Next.js applications with React
- **Backend**: Go applications with Gin framework
- **Mobile**: Android (Kotlin) and iOS (Swift) applications
- **Infrastructure**: Docker, Kubernetes, and Terraform configurations

**Note**: All templates are embedded directly into the binary, making the generator completely self-contained and portable.

## Configuration

The generator uses a straightforward YAML configuration format:

```yaml
name: "my-project"
organization: "my-org"
description: "A sample project"
license: "MIT"
components:
  frontend:
    nextjs:
      app: true
      home: true
      admin: false
      shared: true
  backend:
    go_gin: true
  mobile:
    android: false
    ios: false
  infrastructure:
    docker: true
    kubernetes: false
    terraform: false
versions:
  node: "20.0.0"
  go: "1.21.0"
  packages:
    react: "18.2.0"
    next: "13.4.0"
    typescript: "5.0.0"
output_path: "./output/generated"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For questions and support, please open an issue on GitHub.
