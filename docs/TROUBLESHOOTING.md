# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the Open Source Project Generator.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Problems](#configuration-problems)
- [Generation Errors](#generation-errors)
- [Validation Issues](#validation-issues)
- [Template Problems](#template-problems)
- [Performance Issues](#performance-issues)
- [Debugging](#debugging)
- [Getting Help](#getting-help)

## Installation Issues

### Permission Denied Errors

**Problem**: Cannot install or run the generator due to permission errors.

**Solutions**:

```bash
# Check current permissions
ls -la /usr/local/bin/generator

# Fix permissions
sudo chmod +x /usr/local/bin/generator

# Install to user directory instead
mkdir -p ~/.local/bin
cp generator ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

### Binary Not Found

**Problem**: `generator: command not found` after installation.

**Solutions**:

```bash
# Check if binary exists
which generator
ls -la /usr/local/bin/generator

# Add to PATH
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Or use full path
/usr/local/bin/generator --version
```

### Go Version Compatibility

**Problem**: Build fails due to Go version incompatibility.

**Solutions**:

```bash
# Check Go version
go version

# Update Go (if needed)
# Using g (Go version manager)
g install latest
g use latest

# Or download from https://golang.org/dl/
```

### Docker Issues

**Problem**: Docker container fails to start or run.

**Solutions**:

```bash
# Check Docker is running
docker --version
docker ps

# Pull latest image
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# Run with proper volume mounting
docker run -it --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate
```

## Configuration Problems

### Configuration File Not Found

**Problem**: `Error: configuration file not found`.

**Solutions**:

```bash
# Check file exists
ls -la ./project-config.yaml

# Use absolute path
generator generate --config /full/path/to/config.yaml

# Create default configuration
generator config export default-config.yaml
```

### Invalid Configuration Format

**Problem**: `Error: failed to parse configuration file`.

**Solutions**:

```bash
# Validate YAML syntax
generator config validate ./config.yaml

# Check for common YAML issues
# - Proper indentation (spaces, not tabs)
# - Correct quotes around strings
# - Valid boolean values (true/false, not True/False)

# Use online YAML validator
# https://www.yamllint.com/
```

### Environment Variable Issues

**Problem**: Environment variables not being recognized.

**Solutions**:

```bash
# Check environment variables
env | grep GENERATOR

# Set required variables
export GENERATOR_PROJECT_NAME="my-project"
export GENERATOR_NON_INTERACTIVE=true

# Use .env file
echo "GENERATOR_PROJECT_NAME=my-project" > .env
source .env
```

### Configuration Conflicts

**Problem**: Configuration values conflict or override unexpectedly.

**Solutions**:

```bash
# Show configuration sources
generator config show --sources

# Check precedence order
# 1. Command-line arguments
# 2. Environment variables
# 3. Configuration files
# 4. Default values

# Override specific values
generator generate --config base.yaml --force --minimal
```

## Generation Errors

### Output Directory Issues

**Problem**: Cannot create or write to output directory.

**Solutions**:

```bash
# Check directory permissions
ls -la ./output-directory

# Create directory with proper permissions
mkdir -p ./my-project
chmod 755 ./my-project

# Use different output directory
generator generate --output ~/projects/my-project

# Check disk space
df -h
```

### Template Not Found

**Problem**: `Error: template not found`.

**Solutions**:

```bash
# List available templates
generator list-templates

# Search for specific template
generator list-templates --search "go"

# Check template name spelling
generator template info go-gin

# Use full template path
generator generate --template ./custom-templates/my-template
```

### Generation Timeout

**Problem**: Generation process times out or hangs.

**Solutions**:

```bash
# Use offline mode
generator generate --offline

# Increase timeout
export GENERATOR_TIMEOUT="60m"
generator generate

# Use minimal generation
generator generate --minimal

# Check system resources
top
df -h
```

### Network Issues

**Problem**: Cannot fetch package versions or templates.

**Solutions**:

```bash
# Use offline mode
generator generate --offline

# Check network connectivity
ping github.com
curl -I https://registry.npmjs.org

# Use proxy (if needed)
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080

# Disable SSL verification (not recommended)
export GENERATOR_INSECURE=true
```

## Validation Issues

### Validation Failures

**Problem**: Generated project fails validation.

**Solutions**:

```bash
# Run validation with detailed output
generator validate ./my-project --verbose

# Auto-fix common issues
generator validate ./my-project --fix

# Check specific validation rules
generator validate ./my-project --rules structure,dependencies

# Show fix suggestions
generator validate ./my-project --show-fixes
```

### Import Errors

**Problem**: Generated Go code has import errors.

**Solutions**:

```bash
# Check Go module initialization
cd ./my-project
go mod init my-project
go mod tidy

# Fix import paths
goimports -w .

# Check for missing dependencies
go mod download
go build ./...
```

### Dependency Issues

**Problem**: Package dependencies are outdated or incompatible.

**Solutions**:

```bash
# Update dependencies
generator generate --update-versions

# Check for security vulnerabilities
generator audit --security

# Use specific versions
# In configuration file:
versions:
  react: "18.2.0"
  typescript: "4.9.5"
```

## Template Problems

### Template Syntax Errors

**Problem**: Template compilation fails due to syntax errors.

**Solutions**:

```bash
# Validate template syntax
generator template validate ./my-template

# Check for common syntax issues
# - Missing {{end}} for {{if}} statements
# - Incorrect variable references
# - Invalid function calls

# Test template with sample data
generator template test ./my-template --data sample-data.yaml
```

### Template Variable Issues

**Problem**: Template variables are undefined or incorrect.

**Solutions**:

```bash
# Check template metadata
generator template info my-template

# Validate variable definitions
generator template validate ./my-template --check-variables

# Use template with explicit variables
generator generate --template my-template --var project_name="MyProject"
```

### Custom Template Issues

**Problem**: Custom template doesn't work as expected.

**Solutions**:

```bash
# Validate custom template
generator template validate ./custom-template --detailed

# Check template structure
ls -la ./custom-template/
cat ./custom-template/metadata.yaml

# Test with minimal configuration
generator generate --template ./custom-template --minimal
```

## Performance Issues

### Slow Generation

**Problem**: Project generation is slow.

**Solutions**:

```bash
# Use offline mode
generator generate --offline

# Enable caching
generator cache enable

# Use minimal generation
generator generate --minimal

# Check system resources
top
free -h
df -h
```

### Memory Issues

**Problem**: Out of memory errors during generation.

**Solutions**:

```bash
# Check available memory
free -h

# Use minimal generation
generator generate --minimal

# Clear cache
generator cache clear

# Increase swap space (if needed)
sudo swapon --show
```

### Cache Issues

**Problem**: Cache corruption or performance issues.

**Solutions**:

```bash
# Check cache status
generator cache show

# Clean expired cache
generator cache clean

# Clear all cache
generator cache clear --force

# Rebuild cache
generator update --templates --packages
```

## Debugging

### Enable Debug Mode

```bash
# Debug mode with detailed logging
generator generate --debug --verbose

# Set debug log level
export GENERATOR_LOG_LEVEL=debug
generator generate

# Show debug information
generator version --build-info
```

### Log Analysis

```bash
# View recent logs
generator logs --lines 100

# Filter by log level
generator logs --level error

# Follow logs in real-time
generator logs --follow

# Show log file locations
generator logs --locations
```

### System Information

```bash
# Show system information
generator version --system-info

# Check dependencies
generator version --dependencies

# Verify installation
generator version --verify
```

### Performance Profiling

```bash
# Enable performance metrics
generator generate --debug --profile

# Check generation time
time generator generate

# Monitor resource usage
generator generate --debug &
top -p $!
```

## Getting Help

### Built-in Help

```bash
# General help
generator --help

# Command-specific help
generator generate --help
generator validate --help
generator audit --help

# Show examples
generator <command> --help | grep -A 20 "Examples:"
```

### Diagnostic Information

```bash
# Generate diagnostic report
generator version --diagnostic > diagnostic.txt

# Include system information
generator version --system-info >> diagnostic.txt

# Include configuration
generator config show >> diagnostic.txt
```

### Community Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/cuesoftinc/open-source-project-generator/issues)
- **Discussions**: [Community discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- **Documentation**: [Online documentation](https://github.com/cuesoftinc/open-source-project-generator/wiki)
- **Email Support**: [support@generator.dev](mailto:support@generator.dev)

### Reporting Issues

When reporting issues, include:

1. **Generator version**: `generator version`
2. **System information**: `generator version --system-info`
3. **Configuration**: `generator config show`
4. **Error logs**: `generator logs --level error --lines 50`
5. **Steps to reproduce**: Detailed steps to reproduce the issue
6. **Expected behavior**: What you expected to happen
7. **Actual behavior**: What actually happened

### Example Issue Report

```bash
# Generate complete diagnostic information
generator version --diagnostic > issue-report.txt
generator config show >> issue-report.txt
generator logs --level error --lines 100 >> issue-report.txt

# Attach to GitHub issue
```

## Common Solutions

### Quick Fixes

```bash
# Reset to defaults
generator config reset

# Clear cache and rebuild
generator cache clear
generator update --templates --packages

# Use offline mode
generator generate --offline

# Use minimal generation
generator generate --minimal

# Force regeneration
generator generate --force
```

### Environment Reset

```bash
# Remove all configuration
rm -rf ~/.generator
rm -rf ~/.cache/generator

# Reinstall
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash

# Verify installation
generator --version
```

### Configuration Reset

```bash
# Reset configuration to defaults
generator config reset

# Remove custom configuration
rm -f ~/.config/generator/config.yaml
rm -f ./generator.yaml

# Start fresh
generator generate
```

This troubleshooting guide should help you resolve most common issues with the Open Source Project Generator. If you continue to experience problems, please refer to the community support channels listed above.
