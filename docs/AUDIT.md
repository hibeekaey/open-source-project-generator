# Codebase Audit Documentation

This document describes the comprehensive audit infrastructure for the open-source template generator project.

## Overview

The audit system performs systematic analysis of the codebase to ensure:

- Proper Go project structure and conventions
- Clean dependencies without security vulnerabilities
- High code quality and maintainability
- Comprehensive test coverage
- Up-to-date templates and dependencies

## Quick Start

### Run Full Audit

```bash
# Run comprehensive audit
make audit

# Or directly
./scripts/audit.sh
```

### Run Specific Analysis

```bash
# Structural analysis only
make audit-structure
./scripts/audit.sh --structure

# Dependency analysis only  
make audit-dependencies
./scripts/audit.sh --dependencies

# Code quality analysis only
make audit-quality
./scripts/audit.sh --quality
```

## Audit Components

### 1. Structural Analysis (`audit-structure.sh`)

Validates project organization and Go conventions:

- **Go Project Structure**: Validates `cmd/`, `internal/`, `pkg/` directories
- **Naming Conventions**: Checks file and package naming
- **Test Organization**: Ensures tests are properly co-located
- **Import Organization**: Validates import structure and detects circular dependencies
- **Template Organization**: Checks template directory structure

**Requirements Addressed**: 1.1, 1.2, 1.3, 1.4, 1.5

### 2. Dependency Analysis (`audit-dependencies.sh`)

Analyzes project dependencies:

- **Go Dependencies**: Checks for unused dependencies in `go.mod`
- **Security Scanning**: Uses `govulncheck` to find vulnerabilities
- **Template Dependencies**: Validates `package.json.tmpl` files
- **Docker Dependencies**: Checks Dockerfile base images
- **Version Consistency**: Ensures consistent Go versions across templates

**Requirements Addressed**: 1.1, 1.2, 1.3, 1.4, 1.5

### 3. Code Quality Analysis (`audit-quality.sh`)

Comprehensive code quality checks:

- **Formatting**: Uses `gofmt` to check code formatting
- **Static Analysis**: Runs `go vet` for static analysis
- **Linting**: Uses `golangci-lint` with comprehensive rules
- **Unused Code**: Detects unused functions, variables, imports
- **Complexity**: Measures cyclomatic complexity with `gocyclo`
- **Test Coverage**: Generates coverage reports and validates coverage levels

**Requirements Addressed**: 1.1, 1.2, 1.3, 1.4, 1.5

## Configuration

### golangci-lint Configuration

The audit uses `.golangci.yml` for comprehensive linting with:

- 40+ enabled linters
- Project-specific rules and exclusions
- Performance and security checks
- Code style and formatting validation

### Audit Configuration

`config/audit-config.yml` contains settings for:

- Analysis thresholds (coverage, complexity)
- Tool configurations
- Exclusion patterns
- Reporting preferences

## Output and Reports

### Audit Results Directory

All audit results are stored in `audit-results/`:

```
audit-results/
├── audit.log                    # Comprehensive audit log
├── audit-report.json           # Machine-readable report
├── summary.md                  # Human-readable summary
├── structure-results.txt       # Structural analysis results
├── dependency-results.txt      # Dependency analysis results
├── quality-results.txt         # Quality analysis results
├── dependency-report.md        # Detailed dependency report
├── quality-report.md           # Detailed quality report
├── unused-go-deps.txt          # List of unused Go dependencies
├── unformatted-files.txt       # Files needing formatting
├── golangci-lint-output.txt    # Linting issues
├── coverage.out                # Test coverage data
└── coverage-report.txt         # Coverage summary
```

### Report Formats

#### Summary Report (`summary.md`)

- High-level overview
- Issue counts by category
- Key recommendations

#### Detailed Reports

- `dependency-report.md`: Dependency analysis details
- `quality-report.md`: Code quality analysis details
- `audit-report.json`: Machine-readable structured data

## Tools and Dependencies

### Required Tools

- **Go**: 1.22+ for running the project
- **git**: Version control operations
- **make**: Build automation
- **docker**: Container analysis (optional)

### Automatically Installed Tools

The audit scripts automatically install these tools if missing:

- **golangci-lint**: Comprehensive Go linting
- **govulncheck**: Security vulnerability scanning
- **unused**: Unused code detection
- **gocyclo**: Cyclomatic complexity analysis

### Optional Tools

- **jq**: JSON processing (improves report generation)
- **goimports**: Import organization checking

## Integration

### Makefile Targets

```bash
make audit              # Full audit
make audit-structure    # Structural analysis
make audit-dependencies # Dependency analysis  
make audit-quality      # Code quality analysis
make audit-clean        # Clean audit results
```

### CI/CD Integration

Add to GitHub Actions workflow:

```yaml
- name: Run Codebase Audit
  run: |
    make audit
    
- name: Upload Audit Results
  uses: actions/upload-artifact@v3
  with:
    name: audit-results
    path: audit-results/
```

## Interpreting Results

### Exit Codes

- `0`: No issues found
- `>0`: Number of issues found (non-zero exit code)

### Issue Severity

- **ERROR**: Critical issues that must be fixed
- **WARN**: Issues that should be addressed
- **INFO**: Informational messages

### Common Issues and Solutions

#### Structural Issues

- **Misplaced files**: Move to appropriate directories
- **Naming violations**: Rename files to follow Go conventions
- **Import cycles**: Refactor to break circular dependencies

#### Dependency Issues

- **Unused dependencies**: Remove from `go.mod`
- **Security vulnerabilities**: Update to patched versions
- **Version inconsistencies**: Align versions across templates

#### Quality Issues

- **Formatting**: Run `gofmt -w .`
- **Linting violations**: Address specific golangci-lint issues
- **Low coverage**: Add tests for uncovered code
- **High complexity**: Refactor complex functions

## Best Practices

### Regular Auditing

- Run audit before major releases
- Include in CI/CD pipeline
- Address issues promptly
- Monitor trends over time

### Maintaining Quality

- Keep dependencies updated
- Follow Go conventions consistently
- Maintain high test coverage (>70%)
- Address TODO/FIXME comments regularly

### Template Maintenance

- Update template dependencies regularly
- Ensure security best practices in templates
- Validate generated projects work correctly
- Keep documentation synchronized

## Troubleshooting

### Common Problems

#### Tool Installation Failures

```bash
# Manually install tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

#### Permission Issues

```bash
# Make scripts executable
chmod +x scripts/audit*.sh
```

#### Large Codebases

```bash
# Run analysis in parts
./scripts/audit.sh --structure
./scripts/audit.sh --dependencies  
./scripts/audit.sh --quality
```

### Getting Help

- Check audit logs in `audit-results/audit.log`
- Review specific tool outputs in `audit-results/`
- Run individual scripts for debugging
- Use `--help` flag for usage information

## Contributing

When contributing to the audit infrastructure:

1. Test changes on the full codebase
2. Update documentation for new features
3. Ensure backward compatibility
4. Add appropriate error handling
5. Follow the existing code style

## Future Enhancements

Planned improvements:

- Integration with external security scanners
- Performance benchmarking
- Automated fix suggestions
- Historical trend analysis
- Custom rule definitions
- Integration with code review tools
