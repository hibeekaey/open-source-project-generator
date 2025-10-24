# Troubleshooting Guide

Common issues and solutions for the Open Source Project Generator.

## Table of Contents

- [Interactive Mode Issues](#interactive-mode-issues)
- [Tool Discovery Issues](#tool-discovery-issues)
- [Tool Execution Failures](#tool-execution-failures)
- [Configuration Problems](#configuration-problems)
- [Generation Errors](#generation-errors)
- [Cache Management Issues](#cache-management-issues)
- [Offline Mode Issues](#offline-mode-issues)
- [Performance Problems](#performance-problems)
- [Exit Code Reference](#exit-code-reference)
- [Debug Techniques](#debug-techniques)

---

## Interactive Mode Issues

### Interactive Mode Won't Start

**Problem:** `Error: Interactive mode not supported in this terminal`

**Diagnosis:**

```bash
# Check terminal type
echo $TERM

# Check if stdin is a terminal
test -t 0 && echo "Terminal" || echo "Not a terminal"
```

**Solutions:**

1. **Use a proper terminal:**

   ```bash
   # Try different terminal emulators
   # macOS: iTerm2, Terminal.app
   # Linux: gnome-terminal, konsole, xterm
   # Windows: Windows Terminal, PowerShell, Git Bash
   ```

2. **Use configuration file instead:**

   ```bash
   # Generate configuration template
   generator init-config project.yaml

   # Edit the file
   vim project.yaml

   # Generate from config
   generator generate --config project.yaml
   ```

3. **Check SSH connection:**

   ```bash
   # If using SSH, ensure terminal forwarding
   ssh -t user@host generator generate --interactive
   ```

### Input Not Recognized

**Problem:** Keyboard input doesn't work in interactive mode

**Solutions:**

1. **Check terminal compatibility:**

   ```bash
   # Ensure terminal supports ANSI escape codes
   # Try a different terminal emulator
   ```

2. **Use arrow keys for navigation:**

   - ↑/↓: Move between options
   - Space: Select/deselect
   - Enter: Confirm

3. **Disable terminal multiplexers temporarily:**

   ```bash
   # Exit tmux/screen if having issues
   exit  # from tmux/screen
   generator generate --interactive
   ```

### Validation Errors in Interactive Mode

**Problem:** Can't get past validation for certain fields

**Diagnosis:**

```bash
# Run with verbose to see validation rules
generator generate --interactive --verbose
```

**Solutions:**

1. **Follow validation rules:**

   - **Go module**: `github.com/org/project` format
   - **Android package**: `com.example.app` (lowercase, dots)
   - **iOS bundle ID**: `com.example.app` (reverse domain)
   - **Port**: 1-65535

2. **See validation examples:**

   ```bash
   # Check documentation
   cat docs/CONFIGURATION.md | grep -A 10 "Validation Rules"
   ```

3. **Use configuration file for complex cases:**

   ```bash
   generator init-config --example fullstack
   # Edit and customize
   generator generate --config project.yaml
   ```

### Interactive Mode Hangs

**Problem:** Interactive mode freezes or hangs

**Solutions:**

1. **Check for background processes:**

   ```bash
   # Kill any stuck generator processes
   pkill -9 generator
   ```

2. **Clear terminal state:**

   ```bash
   # Reset terminal
   reset

   # Or
   stty sane
   ```

3. **Use timeout:**

   ```bash
   # Set timeout for interactive mode
   timeout 300 generator generate --interactive
   ```

### Want to Go Back to Previous Step

**Problem:** Made a mistake and want to change previous answer

**Solutions:**

1. **Cancel and restart:**

   ```bash
   # Press Ctrl+C to cancel
   # Start over
   generator generate --interactive
   ```

2. **Use configuration file for more control:**

   ```bash
   # Generate config from interactive mode
   generator generate --interactive
   # Config saved to .generator/generated-config.yaml

   # Edit and regenerate
   vim .generator/generated-config.yaml
   generator generate --config .generator/generated-config.yaml
   ```

### Interactive Mode Cancelled Unexpectedly

**Problem:** `Error: Operation cancelled by user (exit code 5)`

**Diagnosis:**

- Check if you pressed Ctrl+C
- Check if terminal was closed
- Check if SSH connection dropped

**Solutions:**

1. **Resume with configuration file:**

   ```bash
   # If partial config was saved
   generator generate --config .generator/generated-config.yaml
   ```

2. **Start fresh:**

   ```bash
   generator generate --interactive
   ```

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

**Problem:** `Warning: Tool 'go' version 1.18.0 is below minimum required version 1.25.0`

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
npx create-next-app@16.0.0 test-app --typescript

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
   npx create-next-app@16.0.0 test-app

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

## Cache Management Issues

### Cache Corruption

**Problem:** `Error: Cache validation failed: corrupted entries detected`

**Diagnosis:**

```bash
# Validate cache
generator cache-tools --validate

# Check cache file
ls -la ~/.cache/generator/tools.cache
cat ~/.cache/generator/tools.cache
```

**Solutions:**

1. **Clear and rebuild cache:**

   ```bash
   # Clear corrupted cache
   generator cache-tools --clear

   # Rebuild cache
   generator cache-tools --save

   # Verify
   generator cache-tools --validate
   ```

2. **Delete cache file manually:**

   ```bash
   # Remove cache file
   rm ~/.cache/generator/tools.cache

   # Regenerate
   generator cache-tools --save
   ```

3. **Check file permissions:**

   ```bash
   # Fix permissions
   chmod 600 ~/.cache/generator/tools.cache
   ```

### Cache Expired Entries

**Problem:** `Warning: Cache contains expired entries`

**Diagnosis:**

```bash
# Check cache statistics
generator cache-tools --stats

# Validate cache
generator cache-tools --validate
```

**Solutions:**

1. **Refresh cache:**

   ```bash
   # Re-check all tools and update cache
   generator cache-tools --refresh

   # Verify
   generator cache-tools --stats
   ```

2. **Clear and rebuild:**

   ```bash
   generator cache-tools --clear
   generator cache-tools --save
   ```

### Cache Import/Export Issues

**Problem:** `Error: Failed to import cache: invalid format`

**Diagnosis:**

```bash
# Check export file format
cat tools-cache.json

# Validate JSON
python -m json.tool tools-cache.json
```

**Solutions:**

1. **Verify export format:**

   ```bash
   # Export should be valid JSON
   generator cache-tools --export tools-cache.json

   # Check format
   cat tools-cache.json | jq .
   ```

2. **Check version compatibility:**

   ```bash
   # Export includes version field
   # Ensure versions match
   cat tools-cache.json | jq .version
   ```

3. **Re-export from source:**

   ```bash
   # On source machine
   generator cache-tools --export tools-cache.json

   # Transfer to target machine
   scp tools-cache.json user@target:~/

   # On target machine
   generator cache-tools --import ~/tools-cache.json
   ```

### Cache Not Being Used

**Problem:** Generator not using cached tool information

**Diagnosis:**

```bash
# Check cache location
generator cache-tools --info

# Check cache statistics
generator cache-tools --stats

# Run with verbose
generator generate --config project.yaml --verbose
```

**Solutions:**

1. **Verify cache exists:**

   ```bash
   # Check cache file
   ls -la ~/.cache/generator/tools.cache

   # If missing, create it
   generator cache-tools --save
   ```

2. **Check cache TTL:**

   ```bash
   # Cache may be expired
   generator cache-tools --stats

   # Refresh if needed
   generator cache-tools --refresh
   ```

3. **Force cache usage:**

   ```bash
   # Use offline mode to force cache usage
   generator generate --config project.yaml --offline
   ```

### Cache Location Issues

**Problem:** Can't find or access cache directory

**Solutions:**

1. **Check cache location:**

   ```bash
   # Show cache information
   generator cache-tools --info

   # Default locations:
   # Linux: ~/.cache/generator/
   # macOS: ~/Library/Caches/generator/
   # Windows: %LOCALAPPDATA%\generator\cache\
   ```

2. **Create cache directory:**

   ```bash
   # Linux/macOS
   mkdir -p ~/.cache/generator

   # Set permissions
   chmod 700 ~/.cache/generator
   ```

3. **Use custom cache location:**

   ```bash
   # Set environment variable
   export GENERATOR_CACHE_DIR=~/my-cache
   generator cache-tools --save
   ```

### Cache Statistics Not Showing

**Problem:** `generator cache-tools --stats` shows no data

**Solutions:**

```bash
# Initialize cache
generator cache-tools --save

# Verify
generator cache-tools --stats

# If still empty, check tools
generator check-tools
generator cache-tools --save
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

## Exit Code Reference

Understanding exit codes helps diagnose issues and handle errors in scripts.

### Exit Code Summary

| Code | Name | Description | When It Occurs |
|------|------|-------------|----------------|
| 0 | Success | Operation completed successfully | Project generated without errors |
| 1 | Configuration Error | Invalid or missing configuration | Config file issues, validation failures |
| 2 | Tools Missing | Required tools not available | npx, go, gradle, xcodebuild not found |
| 3 | Generation Failed | Component generation failed | Tool execution error, network issue |
| 4 | File System Error | Cannot access files/directories | Permission denied, disk full |
| 5 | User Cancelled | User cancelled the operation | Ctrl+C pressed, cancelled in interactive mode |

### Exit Code 0: Success

**Meaning:** Project generated successfully

**What Happened:**

- All components generated without errors
- Project structure validated
- Integration files created
- No warnings or errors

**Example:**

```bash
generator generate --config project.yaml
echo $?  # Output: 0
```

### Exit Code 1: Configuration Error

**Meaning:** Invalid or missing configuration

**Common Causes:**

- Configuration file not found
- Invalid YAML/JSON syntax
- Missing required fields (name, output_dir)
- Invalid component type
- Missing required component configuration

**Examples:**

```bash
# Missing config file
generator generate --config missing.yaml
# Exit code: 1

# Invalid YAML syntax
generator generate --config invalid.yaml
# Exit code: 1

# Missing required field
# config.yaml: missing 'name' field
generator generate --config config.yaml
# Exit code: 1
```

**How to Fix:**

```bash
# Validate configuration
generator generate --config project.yaml --dry-run

# Generate valid template
generator init-config template.yaml

# Check syntax
yamllint project.yaml
```

### Exit Code 2: Tools Missing

**Meaning:** Required bootstrap tools not found

**Common Causes:**

- npx not installed (for Next.js)
- go not installed (for Go backend)
- gradle not installed (for Android)
- xcodebuild not installed (for iOS)
- Tool not in PATH

**Examples:**

```bash
# npx not found
generator generate --config nextjs-project.yaml
# Exit code: 2

# go not found
generator generate --config go-project.yaml
# Exit code: 2
```

**How to Fix:**

```bash
# Check which tools are missing
generator check-tools

# Install missing tools
# macOS
brew install node go gradle

# Ubuntu
sudo apt install nodejs npm golang-go gradle

# Or use fallback generation
generator generate --config project.yaml --no-external-tools
```

### Exit Code 3: Generation Failed

**Meaning:** Component generation or tool execution failed

**Common Causes:**

- Tool execution error
- Network connectivity issues
- Invalid tool options
- Tool version incompatibility
- Disk space issues during generation

**Examples:**

```bash
# Tool execution failed
generator generate --config project.yaml
# Exit code: 3

# Network error during npm install
generator generate --config nextjs-project.yaml
# Exit code: 3
```

**How to Fix:**

```bash
# Check logs
cat ~/.cache/generator/logs/generator.log

# Run with verbose output
generator generate --config project.yaml --verbose

# Check network connectivity
ping registry.npmjs.org

# Try with fallback
generator generate --config project.yaml --no-external-tools

# Check disk space
df -h
```

### Exit Code 4: File System Error

**Meaning:** Cannot read/write files or directories

**Common Causes:**

- Permission denied
- Disk full
- Invalid output path
- Read-only file system
- File already exists (without --force)

**Examples:**

```bash
# Permission denied
generator generate --config project.yaml --output /root/project
# Exit code: 4

# Disk full
generator generate --config project.yaml
# Exit code: 4

# Output directory exists
generator generate --config project.yaml --output ./existing-dir
# Exit code: 4
```

**How to Fix:**

```bash
# Check permissions
ls -la ./output-directory

# Fix permissions
chmod 755 ./output-directory

# Check disk space
df -h

# Use different output directory
generator generate --config project.yaml --output ~/projects/my-project

# Force overwrite
generator generate --config project.yaml --force
```

### Exit Code 5: User Cancelled

**Meaning:** User cancelled the operation

**Common Causes:**

- Pressed Ctrl+C during generation
- Cancelled in interactive mode
- Declined confirmation prompt
- SSH connection dropped

**Examples:**

```bash
# User pressed Ctrl+C
generator generate --interactive
^C
# Exit code: 5

# User declined confirmation
generator generate --config project.yaml
? Overwrite existing directory? No
# Exit code: 5
```

**How to Fix:**

```bash
# This is intentional cancellation
# To proceed, run the command again and confirm

# Or use --force to skip confirmation
generator generate --config project.yaml --force
```

### Using Exit Codes in Scripts

**Bash Script:**

```bash
#!/bin/bash

generator generate --config project.yaml
EXIT_CODE=$?

case $EXIT_CODE in
  0)
    echo "✓ Success"
    ;;
  1)
    echo "✗ Configuration error"
    generator init-config template.yaml
    exit 1
    ;;
  2)
    echo "✗ Tools missing"
    generator check-tools
    exit 2
    ;;
  3)
    echo "✗ Generation failed"
    cat ~/.cache/generator/logs/generator.log
    exit 3
    ;;
  4)
    echo "✗ File system error"
    df -h
    exit 4
    ;;
  5)
    echo "✗ Cancelled by user"
    exit 5
    ;;
esac
```

**Makefile:**

```makefile
.PHONY: generate
generate:
 @generator generate --config project.yaml || \
 (EXIT_CODE=$$?; \
  if [ $$EXIT_CODE -eq 2 ]; then \
    echo "Tools missing, trying fallback..."; \
    generator generate --config project.yaml --no-external-tools; \
  else \
    exit $$EXIT_CODE; \
  fi)
```

**CI/CD (GitHub Actions):**

```yaml
- name: Generate Project
  id: generate
  run: generator generate --config .generator/project.yaml
  continue-on-error: true

- name: Handle Failure
  if: steps.generate.outcome == 'failure'
  run: |
    EXIT_CODE=${{ steps.generate.outputs.exit_code }}
    if [ "$EXIT_CODE" = "2" ]; then
      echo "Tools missing, using fallback"
      generator generate --config .generator/project.yaml --no-external-tools
    else
      echo "Generation failed with exit code $EXIT_CODE"
      exit $EXIT_CODE
    fi
```

### Exit Code Best Practices

1. **Always check exit codes** in automated scripts
2. **Log exit codes** for debugging
3. **Handle each exit code appropriately**
4. **Provide fallback options** for recoverable errors
5. **Display helpful messages** based on exit code
6. **Use exit codes in CI/CD** for proper error handling

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
npx create-next-app@16.0.0 test-app --typescript
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

| Error Message | Exit Code | Likely Cause | Solution |
|---------------|-----------|--------------|----------|
| `Tool 'X' not found` | 2 | Tool not installed or not in PATH | Install tool or add to PATH |
| `Tool execution failed` | 3 | Network issue, permissions, or tool error | Check logs, verify tool works manually |
| `Invalid configuration` | 1 | YAML/JSON syntax error | Validate configuration file |
| `Permission denied` | 4 | Insufficient permissions | Fix file/directory permissions |
| `Structure mapping failed` | 3 | Generated files in unexpected location | Check debug logs, verify component names |
| `Offline mode failed` | 2 | Tools not cached | Cache tools while online |
| `Timeout` | 3 | Tool taking too long | Check network, increase timeout |
| `Out of memory` | 4 | Insufficient RAM | Close applications, increase swap |
| `Operation cancelled` | 5 | User pressed Ctrl+C | Restart operation |
| `Cache validation failed` | 4 | Corrupted cache | Clear and rebuild cache |
| `Interactive mode not supported` | 1 | Terminal incompatibility | Use configuration file instead |

### Error Message Details

#### Configuration Errors (Exit Code 1)

```text
Error: Invalid configuration: missing required field 'name'
Error: Invalid configuration: no components enabled
Error: Invalid component type: 'invalid-type'
Error: Invalid component config: missing required field 'module'
Error: Configuration file not found: project.yaml
Error: Failed to parse configuration file: invalid YAML syntax
```

**Fix:** Validate and correct configuration file

#### Tool Errors (Exit Code 2)

```text
Error: Tool 'npx' not found in PATH
Error: Tool 'go' not found in PATH
Error: Tool 'gradle' not found in PATH
Error: Required tools missing and no fallback available
```

**Fix:** Install missing tools or use fallback generation

#### Generation Errors (Exit Code 3)

```text
Error: Tool execution failed: npx create-next-app
Error: Component generation failed: go-backend
Error: Network error: cannot reach registry.npmjs.org
Error: Tool timeout after 5 minutes
Error: Integration step failed
```

**Fix:** Check logs, verify network, try fallback

#### File System Errors (Exit Code 4)

```text
Error: Permission denied: cannot create directory
Error: Disk full: no space left on device
Error: Cannot write to output directory
Error: File already exists (use --force to overwrite)
Error: Cache validation failed: corrupted entries
```

**Fix:** Check permissions, disk space, use --force

#### User Cancellation (Exit Code 5)

```text
Error: Operation cancelled by user
Error: User declined confirmation
Error: Interactive mode cancelled
```

**Fix:** Restart operation and confirm

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
