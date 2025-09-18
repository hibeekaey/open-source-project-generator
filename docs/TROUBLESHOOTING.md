# Troubleshooting Guide

This guide provides solutions for common issues encountered when using the Open Source Project Generator.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Generation Problems](#generation-problems)
- [Validation Errors](#validation-errors)
- [Audit Failures](#audit-failures)
- [Configuration Issues](#configuration-issues)
- [Template Problems](#template-problems)
- [Network and Connectivity](#network-and-connectivity)
- [Performance Issues](#performance-issues)
- [Cache Problems](#cache-problems)
- [Update Issues](#update-issues)
- [Logging and Debugging](#logging-and-debugging)

## Installation Issues

### Generator Not Found After Installation

**Problem**: Command `generator` not found after installation.

**Solutions**:

```bash
# Check if binary is in PATH
which generator

# Add to PATH if installed in custom location
export PATH=$PATH:/path/to/generator

# Reinstall using Go
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest

# Verify installation
generator version
```

### Permission Denied During Installation

**Problem**: Permission errors when installing or running generator.

**Solutions**:

```bash
# Install to user directory
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest

# Or use sudo for system-wide installation (not recommended)
sudo make install

# Fix permissions for existing installation
sudo chown -R $USER:$USER ~/.cache/template-generator
```

### Build Errors from Source

**Problem**: Compilation errors when building from source.

**Solutions**:

```bash
# Ensure Go version compatibility
go version  # Should be 1.21+

# Clean and rebuild
make clean
make build

# Check dependencies
go mod tidy
go mod verify

# Build with verbose output
make build VERBOSE=1
```

## Generation Problems

### Project Generation Fails

**Problem**: Project generation stops with errors.

**Diagnosis**:

```bash
# Run with verbose output
generator generate --verbose --debug

# Check logs
generator logs --level error --lines 50

# Validate configuration first
generator config validate
```

**Common Solutions**:

```bash
# Fix output directory permissions
mkdir -p ./my-project
chmod 755 ./my-project

# Use different output directory
generator generate --output ~/projects/my-project

# Skip validation if blocking
generator generate --skip-validation

# Use minimal generation
generator generate --minimal
```

### Template Processing Errors

**Problem**: Template files fail to process correctly.

**Solutions**:

```bash
# Check template information
generator template info go-gin --detailed

# Validate template
generator template validate ./custom-template --fix

# Use different template
generator list-templates --category backend
generator generate --template nextjs-app

# Clear template cache
generator cache clear
generator update --templates
```

### Configuration File Errors

**Problem**: Configuration file is rejected or causes errors.

**Solutions**:

```bash
# Validate configuration syntax
generator config validate ./my-config.yaml

# Check configuration format
file ./my-config.yaml

# Use configuration template
generator config export --template config-template.yaml

# Show configuration schema
generator config show --schema
```

## Validation Errors

### Project Validation Fails

**Problem**: Generated or existing project fails validation.

**Diagnosis**:

```bash
# Run detailed validation
generator validate --verbose --detailed

# Check specific validation rules
generator validate --rules structure,dependencies

# Show available fixes
generator validate --show-fixes
```

**Solutions**:

```bash
# Auto-fix common issues
generator validate --fix --backup

# Fix specific categories
generator validate --fix --rules structure,formatting

# Use less strict validation
generator validate --ignore-warnings

# Skip problematic rules
generator validate --exclude-rules documentation,examples
```

### Dependency Validation Issues

**Problem**: Dependency validation fails or reports conflicts.

**Solutions**:

```bash
# Update to latest versions
generator generate --update-versions

# Check version compatibility
generator version --packages --compatibility

# Use offline mode with cached versions
generator generate --offline

# Override specific versions in configuration
# Add to config.yaml:
versions:
  go: "1.21"
  node: "20"
  react: "19"
```

## Audit Failures

### Security Audit Fails

**Problem**: Security audit reports high-severity issues.

**Solutions**:

```bash
# Run detailed security audit
generator audit --security --detailed --verbose

# Check specific security categories
generator audit --security --fail-on-medium

# Update dependencies for security fixes
generator update --security-only

# Review and fix security issues manually
generator audit --security --output-format html --output-file security-report.html
```

### Quality Audit Issues

**Problem**: Code quality audit reports poor scores.

**Solutions**:

```bash
# Run quality-focused audit
generator audit --quality --detailed

# Fix validation issues first
generator validate --fix

# Use stricter generation options
generator generate --include-examples --update-versions

# Review quality recommendations
generator audit --quality --recommendations
```

## Configuration Issues

### Configuration Not Loading

**Problem**: Configuration files are not being loaded or recognized.

**Diagnosis**:

```bash
# Show configuration sources
generator config show --sources

# Check configuration hierarchy
generator config show --verbose

# Validate configuration file
generator config validate ./config.yaml
```

**Solutions**:

```bash
# Use absolute path
generator generate --config /full/path/to/config.yaml

# Check file format
file ./config.yaml

# Convert between formats
generator config export --format json config.json

# Reset to defaults
generator config reset
```

### Environment Variables Not Working

**Problem**: Environment variables are not being recognized.

**Solutions**:

```bash
# Check environment variable format
env | grep GENERATOR_

# Use correct variable names
export GENERATOR_PROJECT_NAME="myapp"
export GENERATOR_TEMPLATE="go-gin"

# Verify non-interactive mode
generator generate --non-interactive --verbose

# Show effective configuration
generator config show --sources
```

## Template Problems

### Template Not Found

**Problem**: Specified template cannot be found or loaded.

**Solutions**:

```bash
# List available templates
generator list-templates

# Search for templates
generator list-templates --search api

# Update template cache
generator update --templates

# Check template path for custom templates
generator template validate ./my-template
```

### Custom Template Issues

**Problem**: Custom template validation or processing fails.

**Solutions**:

```bash
# Validate template structure
generator template validate ./my-template --detailed

# Fix template issues automatically
generator template validate ./my-template --fix

# Check template metadata
generator template info ./my-template --variables

# Use template debugging
generator generate --template ./my-template --debug --dry-run
```

## Network and Connectivity

### Network Timeouts

**Problem**: Network requests timeout or fail.

**Solutions**:

```bash
# Use offline mode
generator generate --offline

# Increase timeout (if supported)
export GENERATOR_TIMEOUT=300

# Use cached data
generator cache show
generator generate --offline --template go-gin

# Populate cache when network is available
generator update --templates --packages
```

### Registry Access Issues

**Problem**: Cannot access package registries (npm, Go modules, etc.).

**Solutions**:

```bash
# Check network connectivity
ping registry.npmjs.org
ping proxy.golang.org

# Use offline mode
generator cache offline enable
generator generate --offline

# Configure proxy if needed
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080

# Use cached versions
generator version --packages --cached-only
```

## Performance Issues

### Slow Generation

**Problem**: Project generation takes too long.

**Solutions**:

```bash
# Use minimal generation
generator generate --minimal

# Skip version updates
generator generate --no-update-versions

# Use cached data
generator generate --offline

# Skip validation
generator generate --skip-validation

# Check cache performance
generator cache show
generator cache clean
```

### High Memory Usage

**Problem**: Generator uses excessive memory.

**Solutions**:

```bash
# Use minimal generation
generator generate --minimal --exclude examples,docs

# Clear cache
generator cache clear

# Restart generator process
# (if running as daemon or service)

# Check system resources
free -h
df -h ~/.cache/template-generator
```

## Cache Problems

### Cache Corruption

**Problem**: Cache appears corrupted or invalid.

**Solutions**:

```bash
# Validate cache integrity
generator cache validate

# Repair cache
generator cache repair

# Clear and rebuild cache
generator cache clear
generator update --templates --packages

# Check cache location and permissions
generator cache show
ls -la ~/.cache/template-generator
```

### Cache Size Issues

**Problem**: Cache grows too large or fills disk.

**Solutions**:

```bash
# Check cache size
generator cache show

# Clean expired entries
generator cache clean

# Clear specific cache types
generator cache clear --templates-only
generator cache clear --versions-only

# Configure cache limits (if supported)
generator config set cache.max_size 1GB
generator config set cache.ttl 7d
```

## Update Issues

### Update Failures

**Problem**: Generator or template updates fail.

**Solutions**:

```bash
# Check for updates manually
generator update --check --verbose

# Force update
generator update --install --force

# Update specific components
generator update --templates --packages

# Check update channel
generator update --channel stable --check

# Rollback if needed
generator update --rollback --list
```

### Version Conflicts

**Problem**: Version conflicts after updates.

**Solutions**:

```bash
# Check compatibility
generator version --compatibility

# Use specific version
generator update --version v2.1.0

# Reset to stable channel
generator update --channel stable

# Clear cache after update
generator cache clear
generator update --templates
```

## Logging and Debugging

### Enable Debug Logging

```bash
# Enable debug mode
generator generate --debug --verbose

# Show debug logs
generator logs --level debug --lines 100

# Follow logs in real-time
generator logs --follow --level debug

# Save logs to file
generator logs --format json > debug.log
```

### Log Analysis

```bash
# Show error logs only
generator logs --level error --since "1h"

# Filter by component
generator logs --component template --level warn

# Show performance metrics
generator logs --level debug --search "performance"

# Export logs for analysis
generator logs --format csv --since "24h" > analysis.csv
```

### Common Log Patterns

**Template Processing Errors**:

```
ERROR template: failed to process template file.tmpl: template syntax error
```

Solution: Check template syntax and validate template files.

**Network Timeout Errors**:

```
ERROR version: timeout fetching package version: context deadline exceeded
```

Solution: Use offline mode or check network connectivity.

**Permission Errors**:

```
ERROR filesystem: permission denied creating directory /path/to/output
```

Solution: Check directory permissions or use different output path.

**Configuration Errors**:

```
ERROR config: invalid configuration value for key 'license': must be one of [MIT, Apache-2.0, GPL-3.0]
```

Solution: Use valid configuration values or check configuration schema.

## Getting Additional Help

### Built-in Help

```bash
# Command help
generator --help
generator <command> --help

# Show examples
generator generate --help | grep -A 10 "Examples:"
```

### Diagnostic Information

```bash
# System information
generator version --build-info --compatibility

# Configuration status
generator config show --sources --verbose

# Cache status
generator cache show --detailed

# Recent logs
generator logs --level error --lines 20
```

### Community Support

- **GitHub Issues**: Report bugs and get help
- **Documentation**: Check online docs for updates
- **Examples**: Review working examples and configurations
- **Community Forums**: Ask questions and share solutions

### Creating Bug Reports

When reporting issues, include:

1. **Generator version**: `generator version --build-info`
2. **Command used**: Full command with flags
3. **Configuration**: Sanitized configuration file
4. **Error output**: Complete error messages
5. **Logs**: Relevant log entries (`generator logs --level error`)
6. **Environment**: OS, Go version, network setup
7. **Steps to reproduce**: Minimal reproduction case

```bash
# Generate diagnostic bundle
generator logs --level error --since "1h" > error.log
generator config show --sources > config-info.txt
generator version --build-info > version-info.txt
generator cache show > cache-info.txt

# Include these files in your bug report
```
