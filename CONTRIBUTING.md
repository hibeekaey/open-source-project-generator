# Contributing to Open Source Template Generator

Thank you for your interest in contributing to the Open Source Template Generator! We welcome contributions from everyone and are grateful for every pull request, bug report, and feature suggestion.

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
- 🎨 **Templates**: Add or improve project templates

### Before You Start

1. **Check existing issues**: Look for existing issues or discussions about your idea
2. **Create an issue**: For significant changes, create an issue first to discuss the approach
3. **Fork the repository**: Create your own fork to work on
4. **Read the guidelines**: Familiarize yourself with our coding standards and processes

## Development Setup

### Prerequisites

Ensure you have the required tools installed:

- **Go**: 1.24 or later
- **Make**: Build automation tool
- **Git**: Latest stable version
- **Docker**: Latest stable version (optional, for testing)

### Initial Setup

1. **Fork and Clone**

   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/open-source-template-generator.git
   cd open-source-template-generator
   
   # Add upstream remote
   git remote add upstream https://github.com/your-org/open-source-template-generator.git
   ```

2. **Environment Setup**

   ```bash
   # Set up development environment
   make setup
   
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

# Test your changes
make test
make lint

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name

# Create a Pull Request on GitHub
```

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

#### 📝 Documentation

- Improve existing documentation
- Add missing documentation
- Fix typos and grammar
- Add examples and tutorials

#### 🎨 Templates

- Add new project templates
- Improve existing templates
- Update templates to latest versions
- Add template validation

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

- [ ] **Tests Pass**: Ensure all tests pass locally
- [ ] **Linting**: Code passes all linting checks
- [ ] **Documentation**: Update relevant documentation
- [ ] **Templates**: Test template generation if templates were modified
- [ ] **Rebase**: Rebase your branch on the latest main branch

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
- [ ] I have tested template generation if templates were modified

## Description

Brief description of changes...

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Template update
- [ ] Performance improvement
- [ ] Code refactoring

## Testing

Describe the tests you ran and how to reproduce them...

## Template Testing (if applicable)

If you modified templates, describe how you tested them:
- [ ] Generated sample projects with modified templates
- [ ] Verified generated projects build successfully
- [ ] Tested with different component combinations
```

### Review Process

1. **Automated Checks**: All PRs must pass automated CI checks
2. **Code Review**: At least one maintainer must review and approve
3. **Template Testing**: Template changes are tested with sample generation
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
type TemplateEngine struct {
    versionManager interfaces.VersionManager
    logger         *slog.Logger
}

// Document public functions
// ProcessTemplate processes a template file with the provided configuration
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

```
cmd/                    # Command-line applications
├── generator/          # Main generator CLI
└── standards/          # Standards validation CLI

internal/               # Private application code
├── app/               # Application logic
├── config/            # Configuration management
└── container/         # Dependency injection

pkg/                   # Public interfaces and libraries
├── cli/               # CLI interface
├── filesystem/        # File system operations
├── interfaces/        # Core interfaces
├── models/           # Data models
├── template/         # Template processing
├── validation/       # Validation engine
└── version/          # Version management

templates/             # Template files
├── base/             # Base project templates
├── frontend/         # Frontend templates
├── backend/          # Backend templates
├── mobile/           # Mobile templates
└── infrastructure/   # Infrastructure templates
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

- **[Template Maintenance Guidelines](docs/TEMPLATE_MAINTENANCE.md)** - Comprehensive guide for template development
- **[Template Quick Reference](docs/TEMPLATE_QUICK_REFERENCE.md)** - Essential commands and patterns
- **[Template Validation Checklist](docs/TEMPLATE_VALIDATION_CHECKLIST.md)** - Pre-commit validation checklist

Key requirements for template changes:

- All used functions must have corresponding import statements
- Follow Go import organization conventions (standard library, third-party, local)
- Run validation tools before committing: `go run scripts/validate-templates/main.go --check-imports`
- Test template compilation with sample data: `go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-validation`

### Documentation Standards

- **Code Comments**: Explain why, not what
- **Function Documentation**: Document all public functions
- **README Updates**: Update README for user-facing changes
- **Template Documentation**: Document template variables and usage

### Git Standards

#### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
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

```
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
    err := engine.ProcessTemplate("templates/test.tmpl", config)
    
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
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test packages
go test ./pkg/template/...

# Run integration tests
go test ./test/integration/...

# Validate templates (recommended before committing template changes)
go run scripts/validate-templates/main.go --check-imports

# Generate test project to verify templates
go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-validation
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

- **Email**: <security@your-org.com>
- **Security Advisory**: Use GitHub's private vulnerability reporting
- **Details**: See our Security Policy for full details

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

## Development Resources

### Useful Links

- **Go Documentation**: [https://golang.org/doc/](https://golang.org/doc/)
- **Cobra CLI**: [https://cobra.dev/](https://cobra.dev/)
- **Go Templates**: [https://pkg.go.dev/text/template](https://pkg.go.dev/text/template)
- **Testing**: [https://golang.org/doc/tutorial/add-a-test](https://golang.org/doc/tutorial/add-a-test)

### Project Structure

Understanding the project structure helps with contributions:

- **`cmd/`**: Command-line applications
- **`internal/`**: Private application code
- **`pkg/`**: Public interfaces and libraries
- **`templates/`**: Project templates
- **`test/`**: Integration tests
- **`scripts/`**: Build and utility scripts
- **`docs/`**: Documentation files

## FAQ

### Common Questions

**Q: How do I add a new template?**
A: Create the template files in the appropriate `templates/` subdirectory, following existing patterns. Include proper variable substitution and test the template generation.

**Q: How do I test template changes?**
A: Use `make test` to run unit tests, validate templates with `go run scripts/validate-templates/main.go --check-imports`, then test template generation manually with `./bin/generator generate --dry-run`.

**Q: What should I work on as a first contribution?**
A: Look for issues labeled `good first issue` or `help wanted`. Template improvements are often good starting points.

**Q: How do I update package versions in templates?**
A: Update the version variables in template files and test that generated projects build successfully with the new versions.

**Q: Can I add support for a new technology stack?**
A: Yes! Create an issue first to discuss the approach, then add the necessary templates and update the CLI to support the new stack.

## Thank You

Thank you for contributing to the Open Source Template Generator! Your contributions help developers worldwide create better projects with modern best practices.

---

**Questions?** Feel free to reach out through GitHub Discussions or create an issue.

*Last updated: December 2024*
</text>
