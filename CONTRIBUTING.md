# Contributing to Open Source Project Generator

Thank you for your interest in contributing to the Open Source Project Generator! We welcome contributions from everyone and are grateful for every pull request, bug report, and feature suggestion.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Issue Reporting](#issue-reporting)
- [Security Issues](#security-issues)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Ways to Contribute

- 🐛 **Bug Reports**: Help us identify and fix bugs
- ✨ **Feature Requests**: Suggest new features or improvements
- 📝 **Documentation**: Improve our documentation
- 🧪 **Testing**: Add or improve tests
- 💻 **Code**: Submit bug fixes or new features
- 🔧 **Tools**: Add support for new bootstrap tools

### Before You Start

1. **Check existing issues**: Look for existing issues or discussions about your idea
2. **Create an issue**: For significant changes, create an issue first to discuss the approach
3. **Fork the repository**: Create your own fork to work on
4. **Read the guidelines**: Familiarize yourself with our coding standards and processes

## Development Setup

### Prerequisites

Ensure you have the required tools installed:

- **Go**: 1.25.0 or later
- **Make**: Build automation tool
- **Git**: Latest stable version
- **Docker**: Latest stable version (optional, for testing)

### Initial Setup

1. **Fork and Clone**

   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/open-source-project-generator.git
   cd open-source-project-generator
   
   # Add upstream remote
   git remote add upstream https://github.com/cuesoftinc/open-source-project-generator.git
   ```

2. **Environment Setup**

   ```bash
   # Install dependencies
   go mod download
   
   # Build the generator
   make build
   ```

3. **Verify Installation**

   ```bash
   # Run tests to verify setup
   make test
   
   # Run the generator
   ./bin/generator --help
   ```

### Development Workflow

```bash
# Keep your fork updated
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ... edit files ...

# Run all checks (recommended before committing)
make check

# This runs: fmt, vet, lint, test

# Or run individual checks
make test            # Run tests
make lint            # Run linter
make security-scan   # Run security scans

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name

# Create a Pull Request on GitHub
```

### Docker Compose Development Workflow

Use Docker Compose for a containerized development environment:

```bash
# Start development environment
docker compose --profile development up -d generator-dev

# Enter interactive shell
docker compose --profile development run --rm generator-dev bash

# Inside the container, you can:
make build
make test
make check

# Run tests in containers
docker compose --profile testing up generator-test

# Run security scans in containers
docker compose --profile security up generator-security

# Run linting in containers
docker compose --profile lint up generator-lint

# Stop all services
docker compose down
```

See [docker-compose.yml](docker-compose.yml) for all available profiles and services.

## Contributing Guidelines

### Types of Contributions

#### 🐛 Bug Fixes

- Fix existing functionality that doesn't work as expected
- Include tests that verify the fix
- Update documentation if necessary

#### ✨ New Features

- Add new functionality to the generator
- Discuss significant features in an issue first
- Include comprehensive tests
- Update documentation and examples

#### 📝 Documentation Contributions

- Improve existing documentation
- Add missing documentation
- Fix typos and grammar
- Add examples and tutorials

#### 🔧 Bootstrap Tools

- Add support for new CLI tools (Vite, Angular, etc.)
- Improve existing tool executors
- Add fallback generators for tools
- Update tool version requirements

#### 🧪 Tests

- Add missing test coverage
- Improve existing tests
- Add integration or end-to-end tests
- Performance and load testing

### Contribution Process

1. **Issue First**: For significant changes, create an issue to discuss the approach
2. **Small PRs**: Keep pull requests focused and small when possible
3. **Clear Description**: Provide clear descriptions of what your PR does
4. **Tests Required**: All code changes must include appropriate tests
5. **Documentation**: Update documentation for any user-facing changes

## Pull Request Process

### Before Submitting

- [ ] **All Checks Pass**: Run all checks (`make check`)
- [ ] **Tests Pass**: Ensure all tests pass locally (`make test`)
- [ ] **Security Scans**: Run security scans if needed (`make security-scan`)
- [ ] **Documentation**: Update relevant documentation
- [ ] **Tool Testing**: Test tool executors if bootstrap code was modified
- [ ] **Rebase**: Rebase your branch on the latest main branch
- [ ] **Docker**: Test Docker builds if Dockerfiles were modified (`make docker-build`)

### PR Checklist

```markdown
## Pull Request Checklist

- [ ] I have read the [Contributing Guidelines](CONTRIBUTING.md)
- [ ] My code follows the project's coding standards
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] All new and existing tests pass
- [ ] I have updated the documentation accordingly
- [ ] My commits are properly formatted and descriptive
- [ ] I have rebased my branch on the latest main branch
- [ ] I have tested tool executors if bootstrap code was modified

## Description

Brief description of changes...

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Bootstrap tool update
- [ ] Performance improvement
- [ ] Code refactoring

## Testing

Describe the tests you ran and how to reproduce them...

## Tool Executor Testing (if applicable)

If you modified tool executors, describe how you tested them:
- [ ] Generated sample projects with modified executors
- [ ] Verified generated projects build successfully
- [ ] Tested with different component combinations
- [ ] Tested fallback generation when tool is unavailable
```

### Review Process

1. **Automated Checks**: All PRs must pass automated CI checks
2. **Code Review**: At least one maintainer must review and approve
3. **Tool Testing**: Tool executor changes are tested with sample generation
4. **Merge**: Approved PRs are merged by maintainers

## Coding Standards

### General Principles

- **Consistency**: Follow existing code patterns and conventions
- **Readability**: Write code that is easy to read and understand
- **Simplicity**: Prefer simple solutions over complex ones
- **Performance**: Consider performance implications of your changes
- **Security**: Follow security best practices

### Go Standards

#### Code Style

```go
// Use clear, descriptive names
type BootstrapExecutor struct {
    logger *logger.Logger
}

// Document public functions
// Execute runs the bootstrap tool with the provided configuration
func (e *TemplateEngine) ProcessTemplate(templatePath string, config *models.ProjectConfig) error {
    if err := e.validateTemplate(templatePath); err != nil {
        return fmt.Errorf("template validation failed: %w", err)
    }
    
    // Process template logic here
    return nil
}
```

#### Naming Conventions

- **Packages**: lowercase (`template`, `validation`)
- **Types**: PascalCase (`TemplateEngine`, `ProjectConfig`)
- **Functions**: PascalCase for public (`ProcessTemplate`), camelCase for private (`validateTemplate`)
- **Variables**: camelCase (`templatePath`, `configManager`)
- **Constants**: PascalCase (`DefaultTimeout`, `MaxRetries`)

#### File Organization

The project follows a tool-orchestration architecture with clear separation of concerns:

```text
cmd/                    # Command-line applications
└── generator/          # Main generator CLI
    └── main.go        # Main entry point with Cobra commands

internal/               # Private application code
├── config/            # Configuration parsing and validation
│   ├── parser.go      # YAML/JSON parsing
│   ├── validator.go   # Configuration validation
│   ├── schema.go      # Configuration schema
│   ├── nextjs_validator.go    # Next.js validation
│   ├── go_validator.go        # Go validation
│   ├── android_validator.go   # Android validation
│   └── ios_validator.go       # iOS validation
├── generator/         # Component generators
│   ├── bootstrap/     # Bootstrap tool executors
│   │   ├── executor.go    # Base executor
│   │   ├── nextjs.go      # Next.js executor (uses create-next-app)
│   │   ├── golang.go      # Go executor (uses go mod init)
│   │   ├── android.go     # Android executor (uses Gradle)
│   │   └── ios.go         # iOS executor (uses Xcode)
│   ├── fallback/      # Fallback generators
│   │   ├── generator.go   # Generator interface and registry
│   │   ├── android.go     # Android fallback
│   │   ├── ios.go         # iOS fallback
│   │   └── templates/     # Minimal templates
│   └── mapper/        # Structure mapping
│       └── structure.go   # Maps generated output to standard layout
└── orchestrator/      # Project generation orchestration
    ├── coordinator.go         # Main coordinator
    ├── tool_discovery.go      # Tool detection and version checking
    ├── executor_registry.go   # Executor registry
    ├── integration.go         # Component integration
    ├── rollback.go           # Rollback on failure
    ├── progress.go           # Progress tracking
    └── cache/                # Tool metadata caching
        └── manager.go        # Cache management

pkg/                   # Public interfaces and libraries
├── cli/              # CLI utilities and error types
│   ├── exit_codes.go      # Exit code definitions
│   ├── diagnostics.go     # Error diagnostics
│   ├── suggestion_engine.go # Error suggestions
│   └── interactive/       # Interactive prompts
│       └── prompter.go    # User prompts
├── filesystem/       # File system operations
│   ├── operations.go     # File operations
│   └── backup.go         # Backup and restore
├── interfaces/       # Core interfaces
│   ├── executor.go       # Executor interface
│   ├── generator.go      # Generator interface
│   └── mapper.go         # Mapper interface
├── logger/           # Logging infrastructure
│   ├── logger.go         # Logger implementation
│   └── formatter.go      # Output formatting
├── models/           # Data structures
│   ├── project.go        # Project configuration
│   ├── result.go         # Generation results
│   └── tool.go           # Tool metadata
├── security/         # Security operations
│   ├── sanitizer.go      # Path sanitization
│   ├── scanner.go        # Security scanning
│   ├── validator.go      # Input validation
│   └── tool_validator.go # Tool command validation
└── testhelpers/      # Test utilities
    └── environment.go    # Test environment setup
├── utils/            # Utility functions
└── constants/        # Application constants
```

### Template Standards

#### Template Structure

- Use Go text/template syntax
- Include proper variable substitution
- Add conditional rendering for components
- Follow consistent naming patterns

#### Template Validation

- All templates must be syntactically valid
- Generated projects must build successfully
- Include proper error handling
- Use secure defaults

#### Template Maintenance

For detailed guidelines on maintaining template files, see:

- **[Template Development Guide](docs/TEMPLATE_DEVELOPMENT.md)** - Comprehensive guide for template development
- **[API Reference](docs/API_REFERENCE.md)** - Template functions and variables reference
- **[Configuration Guide](docs/CONFIGURATION.md)** - Template configuration options

Key requirements for template changes:

- All used functions must have corresponding import statements
- Follow Go import organization conventions (standard library, third-party, local)
- Run validation tools before committing: `go run cmd/generator/main.go template validate`
- Test template compilation with sample data: `go run cmd/generator/main.go generate --config configs/minimal.yaml --output test-validation`

### Code Documentation Standards

- **Code Comments**: Explain why, not what
- **Function Documentation**: Document all public functions
- **README Updates**: Update README for user-facing changes
- **Template Documentation**: Document template variables and usage
- **Security Documentation**: Update SECURITY.md for security-related changes
- **Distribution Documentation**: Update DISTRIBUTION.md for build/release changes
- **Docker Documentation**: Update docker-compose.yml comments for service changes

### Git Standards

#### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `template`: Template changes

**Examples:**

```text
feat(cli): add interactive component selection

fix(template): handle null values in package.json generation

docs(readme): update installation instructions

template(frontend): update Next.js to version 15

test(validation): add integration tests for project validation
```

## Testing Guidelines

### Test Requirements

- **Unit Tests**: All new functions must have unit tests
- **Integration Tests**: Add integration tests for new features
- **Template Tests**: Test template generation and validation
- **Coverage**: Maintain or improve test coverage

### Testing Standards

#### Unit Testing

```go
func TestTemplateEngine_ProcessTemplate(t *testing.T) {
    // Setup
    engine := template.NewEngine()
    config := &models.ProjectConfig{
        Name:         "test-project",
        Organization: "test-org",
    }
    
    // Execute
    result, err := executor.Execute(ctx, spec)
    
    // Assert
    assert.NoError(t, err)
}
```

#### Integration Testing

```go
func TestProjectGeneration_FullWorkflow(t *testing.T) {
    // Test complete project generation workflow
    tempDir := t.TempDir()
    
    config := &models.ProjectConfig{
        Name:       "integration-test",
        OutputPath: tempDir,
        Components: models.Components{
            Frontend: models.FrontendComponents{MainApp: true},
            Backend:  models.BackendComponents{API: true},
        },
    }
    
    // Generate project
    err := generateProject(config)
    assert.NoError(t, err)
    
    // Verify generated files exist
    assert.FileExists(t, filepath.Join(tempDir, "integration-test", "README.md"))
    assert.FileExists(t, filepath.Join(tempDir, "integration-test", "Makefile"))
}
```

### Running Tests

```bash
# Run all tests with coverage
make test

# Run tests with specific flags
make test TEST_FLAGS="-v -race"

# Run tests with integration tags
make test TEST_FLAGS="-tags=integration"
go test ./internal/generator/...

# Run all security scans (gosec, govulncheck, staticcheck)
make security-scan

# Run linting
make lint

# Format code
make fmt

# Run go vet
make vet

# Using Docker Compose for testing
docker compose --profile testing up generator-test
docker compose --profile testing up generator-test-coverage
docker compose --profile testing up generator-test-integration

# Validate templates (recommended before committing template changes)
./bin/generator template validate

# Generate test project to verify templates
./bin/generator generate --config configs/minimal.yaml --output test-validation
```

## Documentation

### Types of Documentation

1. **Code Documentation**: Inline comments and function documentation
2. **CLI Documentation**: Command help text and usage examples
3. **User Documentation**: README, guides, and tutorials
4. **Template Documentation**: Template usage and variable reference
5. **Developer Documentation**: Architecture decisions, setup guides

### Documentation Standards

- **Clear and Concise**: Write clear, easy-to-understand documentation
- **Examples**: Include practical examples and code snippets
- **Up-to-Date**: Keep documentation synchronized with code changes
- **Accessible**: Use inclusive language and consider accessibility

## Issue Reporting

### Bug Reports

When reporting bugs, please include:

1. **Clear Title**: Descriptive title summarizing the issue
2. **Environment**: OS, Go version, generator version
3. **Steps to Reproduce**: Detailed steps to reproduce the issue
4. **Expected Behavior**: What you expected to happen
5. **Actual Behavior**: What actually happened
6. **Configuration**: Project configuration used (if applicable)
7. **Logs**: Relevant error messages or logs

### Feature Requests

When requesting features, please include:

1. **Problem Statement**: What problem does this solve?
2. **Proposed Solution**: How should this feature work?
3. **Alternatives**: What alternatives have you considered?
4. **Use Cases**: Specific use cases for this feature
5. **Templates**: Which templates would be affected?

## Security Issues

**Do not report security vulnerabilities through public GitHub issues.**

Instead, please report them responsibly:

- **Security Advisory**: Use GitHub's private vulnerability reporting
- **Details**: See [SECURITY.md](SECURITY.md) for complete security policy and reporting guidelines

### Security Best Practices for Contributors

When contributing code, follow these security practices:

- **Path Sanitization**: Always use `pkg/security/SanitizePath()` for user-provided paths
- **Categorized Errors**: Use error types from `pkg/cli` package
- **No Code Execution**: Never execute user-provided code
- **Input Validation**: Validate all user inputs through `internal/config/` validators
- **Tool Execution**: Validate all tool commands before execution
- **File Permissions**: Use restrictive permissions (0600 for files, 0750 for directories)
- **Security Scanning**: Run `make security-scan` before submitting PRs

See [SECURITY.md](SECURITY.md) for detailed security guidelines.

## Community

### Communication Channels

- **GitHub Discussions**: For general questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Email**: For private communications

### Community Guidelines

- **Be Respectful**: Treat everyone with respect and kindness
- **Be Inclusive**: Welcome people of all backgrounds and experience levels
- **Be Constructive**: Provide helpful feedback and suggestions
- **Be Patient**: Remember that everyone is learning and growing

## Modular Development Guidelines

### Working with the Refactored Structure

The codebase has been refactored into a modular architecture. Understanding this structure is crucial for effective contributions.

#### Adding New Features

1. **Identify the Package**: Determine which package the feature belongs to based on the modular structure
2. **Check Interfaces**: Ensure the feature fits existing interfaces in `pkg/interfaces/` or create new ones
3. **Follow Patterns**: Use existing patterns for similar functionality within the same package
4. **Maintain Modularity**: Keep components focused and avoid cross-cutting concerns
5. **Add Tests**: Include comprehensive tests for new components
6. **Update Documentation**: Document new functionality and interfaces

#### Modifying Existing Features

1. **Locate Components**: Use the modular package structure to find relevant code quickly
2. **Check Dependencies**: Understand component dependencies through interfaces in `pkg/interfaces/`
3. **Respect Boundaries**: Ensure changes don't violate component boundaries
4. **Test Changes**: Ensure changes don't break existing functionality
5. **Update Tests**: Modify tests to reflect changes, including component-specific tests
6. **Validate Integration**: Run integration tests to ensure system coherence

#### Package-Specific Development

**CLI Development** (`cmd/generator/` and `pkg/cli/`):

- **Main Entry**: Main CLI logic in `cmd/generator/main.go`
- **Error Types**: Custom error types in `pkg/cli/exit_codes.go`
- **Diagnostics**: Error diagnostics in `pkg/cli/diagnostics.go`
- **Suggestions**: Error suggestions in `pkg/cli/suggestion_engine.go`
- **Interactive**: Interactive prompts in `pkg/cli/interactive/`

**Bootstrap Tool Development** (`internal/generator/bootstrap/`):

- **Executors**: Tool executors in `internal/generator/bootstrap/` (~200 lines max per executor)
- **Base Executor**: Common functionality in `executor.go`
- **Tool-Specific**: Component executors (nextjs.go, golang.go, android.go, ios.go)
- **Testing**: Comprehensive tests for each executor

**Fallback Generator Development** (`internal/generator/fallback/`):

- **Generators**: Fallback generators for when tools are unavailable
- **Android**: Android fallback in `android.go`
- **iOS**: iOS fallback in `ios.go`
- **Templates**: Minimal templates in `templates/` directory

**Configuration Development** (`internal/config/`):

- **Parser**: Configuration parsing in `parser.go`
- **Validator**: Configuration validation in `validator.go`
- **Component Validators**: Component-specific validation (nextjs_validator.go, go_validator.go, etc.)
- **Schema**: Configuration schema in `schema.go`

#### Best Practices for Modular Development

- **Keep Files Small**: Target maximum 1,000 lines per file (strictly enforced)
- **Single Responsibility**: Each file should have one clear, focused purpose
- **Interface First**: Design interfaces in `pkg/interfaces/` before implementations
- **Component Isolation**: Ensure components can be tested and developed independently
- **Clear Dependencies**: Use dependency injection through interfaces
- **Test Coverage**: Maintain high test coverage for all components
- **Documentation**: Keep documentation up-to-date with changes
- **Package Cohesion**: Keep related functionality within the same package
- **Minimal Coupling**: Minimize dependencies between packages

For detailed information about the package structure, see the [Package Structure Documentation](docs/PACKAGE_STRUCTURE.md) and [Migration Guide](docs/MIGRATION_GUIDE.md).

## Development Resources

### Project Documentation

- **[README.md](README.md)**: Project overview and quick start
- **[SECURITY.md](SECURITY.md)**: Security practices and reporting
- **[DISTRIBUTION.md](DISTRIBUTION.md)**: Distribution and release process
- **[env.example](env.example)**: Environment variable reference
- **[docker-compose.yml](docker-compose.yml)**: Docker orchestration with profiles

### Useful Links

- **Go Documentation**: [https://golang.org/doc/](https://golang.org/doc/)
- **Cobra CLI**: [https://cobra.dev/](https://cobra.dev/)
- **Go Templates**: [https://pkg.go.dev/text/template](https://pkg.go.dev/text/template)
- **Testing**: [https://golang.org/doc/tutorial/add-a-test](https://golang.org/doc/tutorial/add-a-test)
- **Docker Compose**: [https://docs.docker.com/compose/](https://docs.docker.com/compose/)
- **Conventional Commits**: [https://www.conventionalcommits.org/](https://www.conventionalcommits.org/)

### Project Structure

Understanding the project structure helps with contributions:

**Core Directories:**

```text
.
├── cmd/                    # Command-line applications
│   └── generator/          # Main generator CLI (main.go)
├── internal/               # Private application code
│   ├── config/            # Configuration parsing and validation
│   ├── generator/         # Component generators
│   │   ├── bootstrap/     # Bootstrap tool executors (nextjs, go, android, ios)
│   │   ├── fallback/      # Fallback generators when tools unavailable
│   │   └── mapper/        # Structure mapping
│   └── orchestrator/      # Project generation orchestration
│       └── cache/         # Tool metadata caching
├── pkg/                   # Public interfaces and libraries
│   ├── cli/              # CLI utilities and error types
│   │   └── interactive/  # Interactive prompts
│   ├── filesystem/       # File operations
│   ├── interfaces/       # Core interfaces
│   ├── logger/           # Logging infrastructure
│   ├── models/           # Data structures
│   ├── security/         # Security operations
│   └── testhelpers/      # Test utilities
├── configs/              # Example configuration files
├── docs/                 # Documentation
├── scripts/              # Build and utility scripts
└── .github/              # GitHub workflows and templates
```

- **`docs/`**: Documentation files
- **`output/`**: Generated project output

**Docker Files:**

- **`Dockerfile`**: Production image (alpine:3.19, ~39 MB, UID 1001)
- **`Dockerfile.dev`**: Development image (golang:1.25-alpine, ~500 MB, UID 1001)
- **`Dockerfile.build`**: Build image (ubuntu:24.04, ~1.5 GB, UID 1001)
- **`docker-compose.yml`**: Multi-profile orchestration (production, development, testing, build, lint, security)

**Configuration Files:**

- **`Makefile`**: Build automation and commands
- **`go.mod`**: Go dependencies (Go 1.25.0)
- **`env.example`**: Environment variable reference

**Documentation:**

- **`README.md`**: Project overview and quick start
- **`CONTRIBUTING.md`**: This file - contribution guidelines
- **`SECURITY.md`**: Security practices and reporting
- **`DISTRIBUTION.md`**: Distribution and release process
- **`LICENSE`**: MIT License

**Important Notes:**

- All Docker containers use **UID 1001** for consistency
- Use `pkg/security/` for path sanitization
- Use `pkg/cli` error types for categorized error handling
- Follow the modular architecture patterns

## FAQ

### Common Questions

**Q: How do I add support for a new tool?**
A: Create a new executor in `internal/generator/bootstrap/`, register it in the tool discovery system, and add tests. See `docs/ADDING_TOOLS.md` for detailed instructions.

**Q: How do I test tool executor changes?**
A: Run `make test` for unit tests, then test template generation manually with `./bin/generator generate` or use Docker Compose: `docker compose --profile testing up generator-test`.

**Q: What should I work on as a first contribution?**
A: Look for issues labeled `good first issue` or `help wanted`. Template improvements, documentation updates, and test additions are often good starting points.

**Q: How do I update package versions in templates?**
A: Update the version variables in template files and test that generated projects build successfully with the new versions. Run `make test` and generate a test project to verify.

**Q: Can I add support for a new technology stack?**
A: Yes! Create an issue first to discuss the approach, then add the necessary templates and update the CLI to support the new stack.

**Q: How do I work with Docker for development?**
A: Use `docker compose --profile development run --rm generator-dev bash` to get an interactive shell with all development tools. All containers use UID 1001.

**Q: What security practices should I follow?**
A: Always use `pkg/security/SanitizePath()` for user paths, return categorized errors from `pkg/cli` error types, and run `make security-scan` before submitting PRs. See [SECURITY.md](SECURITY.md) for details.

**Q: How do I run CI checks locally?**
A: Run `make check` for quick checks (fmt, vet, lint, test) or `make ci` for the full CI pipeline. For release validation, run `make release`.

**Q: What's the difference between the three Dockerfiles?**
A: `Dockerfile` is for production (39 MB), `Dockerfile.dev` is for development with all tools (500 MB), and `Dockerfile.build` is for creating packages (1.5 GB). All use UID 1001.

## Thank You

Thank you for contributing to the Open Source Project Generator! Your contributions help developers worldwide create better projects with modern best practices.

## Quick Command Reference

### Essential Commands

```bash
# Development
make build              # Build for current platform
make test               # Run tests with coverage
make check              # Run all checks (fmt, vet, lint, test)
make fmt                # Format code
make dev                # Run in development mode

# Quality Assurance
make lint               # Run linter (auto-installs if needed)
make security-scan      # Run all security scans (auto-installs tools)
make ci                 # Run full CI pipeline
make release            # Full release validation

# Docker Compose
docker compose --profile development run --rm generator-dev bash
docker compose --profile testing up generator-test
docker compose --profile lint up generator-lint
docker compose --profile security up generator-security

# Testing
make test                                    # Tests with coverage
make test TEST_FLAGS="-v -race"              # With race detector
make test TEST_FLAGS="-tags=integration"     # Integration tests

# Building
make dist               # Build cross-platform binaries
make package            # Build distribution packages (DEB, RPM, Arch)
make release            # Full release (test, lint, security-scan, dist, package)
make docker-build       # Build production Docker image
make docker-test        # Test Docker image

# Utilities
make version            # Show version information
make validate-setup     # Validate project setup
make clean              # Clean all build artifacts (binaries, packages, reports, archives)
```

---

**Questions?** Feel free to reach out through GitHub Discussions or create an issue.

**Need Help?** Check out our [README](README.md), [SECURITY.md](SECURITY.md), or [DISTRIBUTION.md](DISTRIBUTION.md) for more information.
