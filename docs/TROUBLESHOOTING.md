# Troubleshooting Guide

Common issues and solutions for the Open Source Project Generator.

## Table of Contents

- [Tool Discovery Issues](#tool-discovery-issues)
- [Tool Execution Failures](#tool-execution-failures)
- [Configuration Problems](#configuration-problems)
- [Generation Errors](#generation-errors)
- [Offline Mode Issues](#offline-mode-issues)
- [Performance Problems](#performance-problems)
- [Debug Techniques](#debug-techniques)

---

## Tool Discovery Issues

### Tool Not Found

**Problem:** `Error: Tool 'npx' not found in PATH`

**Diagnosis:**

```bash
# Check if tool is installed
which npx
echo $PATH

# Check tool availability
generator check-tools
```

**Solutions:**

1. **Install the missing tool:**

   ```bash
   # For npx (Node.js)
   # macOS
   brew install node

   # Ubuntu/Debian
   sudo apt install nodejs npm

   # Windows
   # Download from https://nodejs.org/

   # Verify installation
   npx --version
   generator check-tools
   ```

2. **Add tool to PATH:**

   ```bash
   # Find tool location
   find /usr -name npx 2>/dev/null

   # Add to PATH (bash/zsh)
   export PATH="/usr/local/bin:$PATH"

   # Make permanent
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Use fallback generation:**

   ```bash
   # Force fallback generators
   generator generate --config project.yaml --no-external-tools
   ```

### Tool Version Too Old

**Problem:** `Warning: Tool 'go' version 1.18.0 is below minimum required version 1.21.0`

**Solutions:**

1. **Update the tool:**

   ```bash
   # Update Go
   # macOS
   brew upgrade go

   # Ubuntu/Debian
   sudo add-apt-repository ppa:longsleep/golang-backports
   sudo apt update
   sudo apt install golang-go

   # Verify version
   go version
   generator check-tools
   ```

2. **Install specific version:**

   ```bash
   # Using version managers
   # For Node.js (nvm)
   nvm install 20
   nvm use 20

   # For Go (gvm)
   gvm install go1.21
   gvm use go1.21
   ```

### Multiple Tool Versions Conflict

**Problem:** `Error: Multiple versions of 'node' found in PATH`

**Solutions:**

1. **Use version manager:**

   ```bash
   # Install nvm for Node.js
   curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

   # Set default version
   nvm alias default 20
   nvm use default
   ```

2. **Clean up PATH:**

   ```bash
   # Check PATH
   echo $PATH | tr ':' '\n'

   # Remove duplicate entries
   export PATH=$(echo "$PATH" | awk -v RS=':' -v ORS=":" '!a[$1]++')
   ```

---

## Tool Execution Failures

### Bootstrap Tool Fails to Execute

**Problem:** `Error: Tool execution failed: npx create-next-app`

**Diagnosis:**

```bash
# Run with verbose output
generator generate --config project.yaml --verbose

# Test tool manually
npx create-next-app@latest test-app --typescript

# Check logs
cat ~/.cache/generator/logs/generator.log
```

**Solutions:**

1. **Check network connectivity:**

   ```bash
   # Test internet connection
   ping registry.npmjs.org

   # Use offline mode if needed
   generator generate --offline
   ```

2. **Clear tool cache:**

   ```bash
   # Clear npm cache
   npm cache clean --force

   # Clear generator cache
   generator cache-tools --clear

   # Retry generation
   generator generate --config project.yaml
   ```

3. **Check disk space:**

   ```bash
   # Check available space
   df -h

   # Clean up if needed
   npm cache clean --force
   go clean -cache -modcache -testcache
   ```

4. **Use fallback generation:**

   ```bash
   # Skip external tools
   generator generate --config project.yaml --no-external-tools
   ```

### Tool Hangs or Times Out

**Problem:** `Error: Tool execution timeout after 5 minutes`

**Solutions:**

1. **Check for prompts:**

   ```bash
   # Some tools may be waiting for input
   # Run tool manually to see prompts
   npx create-next-app@latest test-app

   # Ensure non-interactive mode
   generator generate --config project.yaml
   ```

2. **Check system resources:**

   ```bash
   # Monitor during execution
   top
   htop

   # Check for resource constraints
   free -h
   df -h
   ```

3. **Increase timeout (if supported):**

   ```bash
   # Set timeout environment variable
   export GENERATOR_TOOL_TIMEOUT=600  # 10 minutes

   # Retry generation
   generator generate --config project.yaml
   ```

### Permission Denied During Tool Execution

**Problem:** `Error: EACCES: permission denied, mkdir '/usr/local/lib/node_modules'`

**Solutions:**

1. **Fix npm permissions:**

   ```bash
   # Change npm default directory
   mkdir ~/.npm-global
   npm config set prefix '~/.npm-global'
   export PATH=~/.npm-global/bin:$PATH

   # Add to shell profile
   echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
   ```

2. **Fix directory permissions:**

   ```bash
   # Change output directory ownership
   sudo chown -R $USER:$USER ./output-directory

   # Or use different output directory
   generator generate --output ~/projects/my-project
   ```

3. **Use sudo (not recommended):**

   ```bash
   # Only if absolutely necessary
   sudo generator generate --config project.yaml
   ```

---

## Configuration Problems

### Configuration File Not Found

**Problem:** `Error: configuration file not found`

**Solutions:**

```bash
# Check file exists
ls -la ./project-config.yaml

# Use absolute path
generator generate --config /full/path/to/config.yaml

# Create default configuration
generator init-config project.yaml
```

### Invalid Configuration Format

**Problem:** `Error: failed to parse configuration file`

**Solutions:**

```bash
# Validate YAML syntax
# Use online validator or yamllint
yamllint project.yaml

# Check for common YAML issues:
# - Proper indentation (spaces, not tabs)
# - Correct quotes around strings
# - Valid boolean values (true/false, not True/False)

# Generate valid template
generator init-config template.yaml

# Compare with your configuration
diff template.yaml project.yaml
```

### Missing Required Fields

**Problem:** `Error: Invalid configuration: missing required field 'name'`

**Solutions:**

```bash
# Check required fields
generator init-config --full > template.yaml

# Compare with your configuration
diff template.yaml project.yaml

# Add missing fields
# Minimum required:
# - name
# - output_dir
# - components (at least one enabled)
```

### Component Configuration Invalid

**Problem:** `Error: Invalid component configuration: unknown component type 'react-app'`

**Solutions:**

```bash
# Valid component types:
# - nextjs
# - go-backend
# - android
# - ios

# Fix component type in configuration
# Change 'react-app' to 'nextjs'
```

---

## Generation Errors

### Output Directory Issues

**Problem:** `Error: Cannot create or write to output directory`

**Solutions:**

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

### Generation Timeout

**Problem:** Generation process times out or hangs

**Solutions:**

```bash
# Use offline mode
generator generate --offline

# Use fallback generators
generator generate --no-external-tools

# Check system resources
top
df -h
free -h
```

### Partial Generation

**Problem:** Some components generated, others failed

**Diagnosis:**

```bash
# Check which components failed
generator generate --config project.yaml --verbose

# Review logs
cat ~/.cache/generator/logs/generator.log
```

**Solutions:**

```bash
# Disable failed components temporarily
# Edit config to set enabled: false for failing components

# Generate successfully working components first
generator generate --config working-components.yaml

# Debug failed components separately
generator generate --config failed-component.yaml --verbose
```

### File Conflicts

**Problem:** `Error: File already exists`

**Solutions:**

```bash
# Create backup before overwriting
generator generate --config project.yaml --backup

# Force overwrite
generator generate --config project.yaml --force

# Use different output directory
generator generate --config project.yaml --output ./new-project
```

---

## Offline Mode Issues

### Cannot Generate in Offline Mode

**Problem:** `Error: Tool 'npx' requires network access and is not cached`

**Solutions:**

1. **Cache tools before going offline:**

   ```bash
   # While online, cache all tools
   generator cache-tools --save

   # Verify cache
   generator cache-tools --stats

   # Then go offline
   generator generate --offline
   ```

2. **Use fallback generators:**

   ```bash
   # Fallback generators work offline
   generator generate --config project.yaml --no-external-tools
   ```

3. **Check cache location:**

   ```bash
   # Verify cache exists
   generator cache-tools --info

   # If cache is missing, recache while online
   generator cache-tools --save
   ```

### Cached Tools Expired

**Problem:** `Warning: Cached tool 'npx' is outdated`

**Solutions:**

```bash
# Update cache (while online)
generator cache-tools --clear
generator cache-tools --save

# Verify update
generator cache-tools --stats
```

### Network Detection Issues

**Problem:** Generator thinks it's offline when online

**Solutions:**

```bash
# Force online mode
unset GENERATOR_OFFLINE

# Check network connectivity
ping google.com

# Regenerate
generator generate --config project.yaml
```

---

## Performance Problems

### Slow Generation

**Problem:** Project generation is slow

**Solutions:**

```bash
# Use offline mode
generator generate --offline

# Use fallback generators (faster)
generator generate --no-external-tools

# Check system resources
top
free -h
df -h

# Close other applications
```

### Memory Issues

**Problem:** Out of memory errors during generation

**Solutions:**

```bash
# Check available memory
free -h

# Close other applications

# Generate components separately
# Split configuration into smaller files

# Increase swap space (if needed)
sudo swapon --show
```

### Cache Issues

**Problem:** Cache corruption or performance issues

**Solutions:**

```bash
# Check cache status
generator cache-tools --stats

# Clear cache
generator cache-tools --clear

# Rebuild cache
generator cache-tools --save
```

---

## Debug Techniques

### Enable Debug Mode

```bash
# Full debug output
generator generate --config project.yaml --debug --verbose

# Save debug output to file
generator generate --config project.yaml --debug 2>&1 | tee debug.log
```

### Check Log Files

```bash
# View recent logs
tail -f ~/.cache/generator/logs/generator.log

# Search for errors
grep -i error ~/.cache/generator/logs/generator.log

# View specific component logs
grep "Next.js" ~/.cache/generator/logs/generator.log
```

### Test Components Individually

```bash
# Test tool availability
generator check-tools --verbose

# Test configuration
generator generate --config project.yaml --dry-run --verbose

# Test specific component
# Create minimal config with one component
cat > test.yaml << EOF
name: "test"
output_dir: "./test"
components:
  - type: nextjs
    name: test-app
    enabled: true
    config:
      typescript: true
EOF

generator generate --config test.yaml --verbose
```

### Verify Tool Execution

```bash
# Run tools manually to isolate issues
npx create-next-app@latest test-app --typescript
go mod init github.com/test/app

# Compare with generator output
generator generate --config project.yaml --verbose
```

### Check System Resources

```bash
# Monitor during generation
# Terminal 1: Run generator
generator generate --config project.yaml --verbose

# Terminal 2: Monitor resources
watch -n 1 'ps aux | grep generator'
watch -n 1 'df -h'
watch -n 1 'free -h'
```

### Validate Generated Project

```bash
# Check directory structure
tree ./my-project -L 2

# Verify files were created
find ./my-project -type f | wc -l

# Check for errors in generated files
cd ./my-project/CommonServer
go vet ./...

cd ./my-project/App
npm run build
```

---

## Common Error Messages

| Error Message | Likely Cause | Solution |
|---------------|--------------|----------|
| `Tool 'X' not found` | Tool not installed or not in PATH | Install tool or add to PATH |
| `Tool execution failed` | Network issue, permissions, or tool error | Check logs, verify tool works manually |
| `Invalid configuration` | YAML/JSON syntax error | Validate configuration file |
| `Permission denied` | Insufficient permissions | Fix file/directory permissions |
| `Structure mapping failed` | Generated files in unexpected location | Check debug logs, verify component names |
| `Offline mode failed` | Tools not cached | Cache tools while online |
| `Timeout` | Tool taking too long | Check network, increase timeout |
| `Out of memory` | Insufficient RAM | Close applications, increase swap |

---

## Getting Help

### Built-in Help

```bash
# General help
generator --help

# Command-specific help
generator generate --help
generator check-tools --help

# Show examples
generator <command> --help
```

### Diagnostic Information

```bash
# Generate diagnostic report
generator version

# Check tool availability
generator check-tools --verbose

# Check cache status
generator cache-tools --stats

# View logs
cat ~/.cache/generator/logs/generator.log
```

### Community Support

- **GitHub Issues**: [Report bugs](https://github.com/cuesoftinc/open-source-project-generator/issues)
- **Discussions**: [Community discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- **Documentation**: [Online documentation](https://github.com/cuesoftinc/open-source-project-generator)
- **Email**: <support@cuesoft.io>

### Reporting Issues

When reporting issues, include:

1. **Generator version**: `generator version`
2. **Tool availability**: `generator check-tools`
3. **Configuration**: Your configuration file (sanitized)
4. **Error logs**: Relevant error messages
5. **Steps to reproduce**: Detailed steps
6. **Expected behavior**: What you expected
7. **Actual behavior**: What actually happened

### Example Issue Report

```bash
# Collect diagnostic information
generator version > issue-report.txt
generator check-tools >> issue-report.txt
echo "--- Configuration ---" >> issue-report.txt
cat project.yaml >> issue-report.txt
echo "--- Error Log ---" >> issue-report.txt
tail -100 ~/.cache/generator/logs/generator.log >> issue-report.txt

# Attach to GitHub issue
```

---

## Quick Fixes

### Reset to Defaults

```bash
# Clear cache
generator cache-tools --clear

# Regenerate configuration
generator init-config project.yaml

# Try again
generator generate --config project.yaml
```

### Force Fallback

```bash
# Skip external tools entirely
generator generate --config project.yaml --no-external-tools
```

### Use Verbose Mode

```bash
# See detailed output
generator generate --config project.yaml --verbose --debug
```

### Test with Dry Run

```bash
# Preview without creating files
generator generate --config project.yaml --dry-run
```

---

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [CLI Commands](CLI_COMMANDS.md) - Command reference
- [Configuration Guide](CONFIGURATION.md) - Configuration options
- [Examples](EXAMPLES.md) - Example configurations
- [Architecture](ARCHITECTURE.md) - How it works
