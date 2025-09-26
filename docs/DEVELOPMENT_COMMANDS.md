# Development Commands Guide

This document provides a comprehensive guide to the development commands available in the project, including testing, linting, and security scanning.

## Quick Reference

### Essential Commands

```bash
# Run all checks (recommended for development)
make check

# Run complete CI pipeline
make ci

# Quick development checks
make check-fast
```

### Individual Commands

```bash
# Testing
make test          # Full test suite with race detection and coverage
make test-fast     # Fast tests without race detection
make test-race     # Tests with race detection only
make test-coverage # Tests with HTML coverage report

# Linting
make lint          # Run golangci-lint with comprehensive checks
make lint-fix      # Run golangci-lint with auto-fix

# Security
make gosec         # Run gosec security scanner
make gosec-verbose # Run gosec with detailed SARIF output

# Combined
make check         # test + lint + gosec
make ci            # install + test + lint + gosec + build
```

## Testing Commands

### `make test`

**Purpose**: Run comprehensive test suite with race detection and coverage
**Features**:

- Race condition detection
- Code coverage analysis
- Verbose output
- All packages tested

**Output**:

- Coverage report: `coverage.out`
- Test results with timing

**Usage**:

```bash
make test
```

### `make test-fast`

**Purpose**: Run tests quickly without race detection
**Features**:

- Fast execution
- No race detection (faster)
- Verbose output
- All packages tested

**Usage**:

```bash
make test-fast
```

### `make test-race`

**Purpose**: Run tests with race detection only
**Features**:

- Race condition detection
- No coverage analysis
- Verbose output
- All packages tested

**Usage**:

```bash
make test-race
```

### `make test-coverage`

**Purpose**: Generate comprehensive coverage reports
**Features**:

- Code coverage analysis
- HTML coverage report
- Coverage file: `coverage.html`
- JSON coverage file: `coverage.out`

**Usage**:

```bash
make test-coverage
# Open coverage.html in browser to view detailed coverage
```

## Linting Commands

### `make lint`

**Purpose**: Run comprehensive code analysis with golangci-lint
**Features**:

- 20+ linters enabled
- Comprehensive code quality checks
- Security analysis
- Performance suggestions
- Style enforcement

**Linters Included**:

- `errcheck` - Error handling verification
- `govet` - Go vet analysis
- `staticcheck` - Static analysis
- `gosec` - Security analysis
- `gocritic` - Code quality analysis
- `revive` - Go linting
- `gosimple` - Code simplification
- And many more...

**Usage**:

```bash
make lint
```

### `make lint-fix`

**Purpose**: Run linter with automatic fixes
**Features**:

- Auto-fixable issues resolved
- Remaining issues reported
- Safe automatic corrections

**Usage**:

```bash
make lint-fix
```

## Security Commands

### `make gosec`

**Purpose**: Run comprehensive security analysis
**Features**:

- Security vulnerability detection
- JSON report output
- Comprehensive security checks
- Report file: `gosec-report.json`

**Security Checks**:

- SQL injection vulnerabilities
- Command injection risks
- Cryptographic issues
- File permission problems
- Network security issues
- And more...

**Usage**:

```bash
make gosec
# Check gosec-report.json for detailed security analysis
```

### `make gosec-verbose`

**Purpose**: Run security analysis with detailed SARIF output
**Features**:

- SARIF format output (industry standard)
- Detailed security analysis
- IDE integration support
- Report file: `gosec-report.sarif`

**Usage**:

```bash
make gosec-verbose
# Use gosec-report.sarif with security tools or IDEs
```

## Combined Commands

### `make check`

**Purpose**: Run comprehensive quality checks
**Includes**:

- `make test` (full test suite with race detection and coverage)
- `make lint` (comprehensive linting)
- `make gosec` (security analysis)

**Usage**:

```bash
make check
```

### `make ci`

**Purpose**: Complete CI/CD pipeline simulation
**Includes**:

- `make install` (dependency installation)
- `make test` (comprehensive testing)
- `make lint` (code quality analysis)
- `make gosec` (security scanning)
- `make build` (application build)

**Usage**:

```bash
make ci
```

### `make check-fast`

**Purpose**: Quick development checks
**Includes**:

- `make test-fast` (fast tests without race detection)
- `make lint` (comprehensive linting)

**Usage**:

```bash
make check-fast
```

## Installation Commands

### `make install-lint`

**Purpose**: Install golangci-lint
**Features**:

- Automatic installation if not present
- Latest stable version
- Proper PATH configuration

### `make install-gosec`

**Purpose**: Install gosec security scanner
**Features**:

- Automatic installation if not present
- Latest version from GitHub
- Go module installation

## Development Workflow

### Daily Development

```bash
# Quick check during development
make check-fast

# Before committing
make check

# Before pushing
make ci
```

### Pre-commit Checklist

1. `make check-fast` - Quick validation
2. `make check` - Comprehensive validation
3. Fix any issues reported
4. `make ci` - Final validation

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Run comprehensive checks
  run: make ci

# GitLab CI example
script:
  - make ci

# Jenkins example
sh 'make ci'
```

## Configuration Files

### `.golangci.yml`

Comprehensive linting configuration with:

- 20+ enabled linters
- Custom rules and settings
- Performance optimizations
- Security exclusions

### Linter Settings

- **gocritic**: Code quality analysis
- **revive**: Go linting rules
- **gomnd**: Magic number detection
- **gosec**: Security analysis
- **govet**: Go vet analysis

## Troubleshooting

### Common Issues

#### Tests Failing

```bash
# Run specific package
go test -v ./pkg/cli/...

# Run with more verbose output
go test -v -race ./...

# Check for race conditions
go test -race ./...
```

#### Linting Issues

```bash
# Run with verbose output
golangci-lint run --verbose

# Check specific linter
golangci-lint run --enable=gosec

# Auto-fix issues
make lint-fix
```

#### Security Issues

```bash
# Run with verbose output
gosec -fmt json -out report.json ./...

# Check specific rules
gosec -include=G204,G304 ./...
```

### Performance Issues

#### Slow Tests

```bash
# Use fast tests during development
make test-fast

# Run specific packages
go test -v ./pkg/cli/...
```

#### Slow Linting

```bash
# Run with fewer linters
golangci-lint run --disable-all --enable=govet,errcheck

# Skip specific directories
golangci-lint run --skip-dirs=vendor,testdata
```

## Best Practices

### Development

1. Use `make check-fast` during active development
2. Use `make check` before committing
3. Use `make ci` before pushing
4. Fix linting issues immediately
5. Address security issues promptly

### CI/CD

1. Use `make ci` in CI pipelines
2. Set appropriate timeouts (15+ minutes)
3. Cache Go modules for faster builds
4. Monitor performance metrics
5. Generate and store reports

### Code Quality

1. Keep test coverage above 80%
2. Address all linting issues
3. Fix security vulnerabilities
4. Use race detection in CI
5. Regular dependency updates

## Report Files

### Generated Files

- `coverage.out` - Go coverage data
- `coverage.html` - HTML coverage report
- `gosec-report.json` - Security analysis (JSON)
- `gosec-report.sarif` - Security analysis (SARIF)

### Cleaning

```bash
# Clean all reports
make clean

# Clean specific files
rm -f coverage.out coverage.html
rm -f gosec-report.json gosec-report.sarif
```

## Integration with IDEs

### VS Code

- Install Go extension
- Configure to use project's golangci-lint
- Enable gosec integration
- Use coverage reports

### GoLand/IntelliJ

- Configure external tools
- Set up file watchers
- Enable security scanning
- Configure coverage display

### Vim/Neovim

- Use vim-go plugin
- Configure ALE for linting
- Set up gosec integration
- Enable coverage highlighting

## Advanced Usage

### Custom Linting

```bash
# Run specific linters
golangci-lint run --enable=gosec,errcheck

# Exclude specific rules
golangci-lint run --disable=gocyclo

# Custom configuration
golangci-lint run --config=.golangci-custom.yml
```

### Custom Security Scanning

```bash
# Include specific rules
gosec -include=G204,G304 ./...

# Exclude specific rules
gosec -exclude=G204 ./...

# Custom severity
gosec -severity=high ./...
```

### Custom Testing

```bash
# Run with specific tags
go test -v -tags=integration ./...

# Run with timeout
go test -v -timeout=10m ./...

# Run with specific packages
go test -v ./pkg/cli/... ./pkg/config/...
```

## Support

For issues with the development commands:

1. Check this documentation
2. Verify Go version (1.25+ required)
3. Ensure all tools are installed
4. Check for configuration issues
5. Review error messages carefully

For questions or issues, please refer to the main project documentation or create an issue in the repository.
