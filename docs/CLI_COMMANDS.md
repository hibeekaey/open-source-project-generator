# CLI Commands Reference

Complete reference for all commands in the Open Source Project Generator.

## Command Overview

| Command | Description |
|---------|-------------|
| `generate` | Generate a new project from configuration |
| `check-tools` | Check availability of bootstrap tools |
| `init-config` | Generate a configuration template |
| `cache-tools` | Manage tool cache for offline use |
| `version` | Print version information |

## Global Flags

These flags work with all commands:

```bash
--verbose, -v    # Enable verbose logging
--help, -h       # Show help information
```

---

## generate

Generate a new project using the specified configuration file.

### generate Synopsis

```bash
generator generate [flags]
```

### generate Description

The generate command orchestrates bootstrap tools (like `create-next-app`, `go mod init`) to create project components, then maps them to a standardized directory structure and integrates them together.

### generate Flags

```bash
--config, -c FILE          # Path to configuration file (YAML/JSON)
--output, -o DIR           # Output directory (overrides config)
--verbose, -v              # Enable verbose logging
--dry-run                  # Preview without creating files
--no-external-tools        # Force fallback generation
--offline                  # Force offline mode
--backup                   # Create backup before overwriting (default: true)
--force                    # Force overwrite existing directory
--interactive, -i          # Interactive mode (not yet implemented)
```

### generate Examples

```bash
# Generate from config file
generator generate --config project.yaml

# With custom output directory
generator generate --config project.yaml --output ./my-project

# Dry run (preview)
generator generate --config project.yaml --dry-run

# Force fallback generation
generator generate --config project.yaml --no-external-tools

# Offline mode
generator generate --config project.yaml --offline

# Verbose output
generator generate --config project.yaml --verbose

# Force overwrite
generator generate --config project.yaml --force
```

### generate Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Configuration error |
| 2 | Tool not found (and no fallback) |
| 3 | Tool execution failed |
| 4 | File system error |
| 5 | Validation error |

---

## check-tools

Check availability of required bootstrap tools and display installation instructions for missing tools.

### check-tools Synopsis

```bash
generator check-tools [tool-names...] [flags]
```

### check-tools Description

Validates your environment before project generation by checking which bootstrap tools are available on your system. Shows installation instructions for any missing tools.

### check-tools Flags

```bash
--verbose, -v    # Enable verbose logging
```

### check-tools Examples

```bash
# Check all registered tools
generator check-tools

# Check specific tools
generator check-tools npx go

# With verbose output
generator check-tools --verbose
```

### check-tools Output Format

```text
Checking 4 tools...

✓ Available Tools:
  ✓ npx (version: 10.2.3)
    Supports: [nextjs]
  ✓ go (version: 1.21.0)
    Supports: [go-backend]

✗ Missing Tools:
  ✗ gradle
    Required for: [android]
    
    Install instructions:
    macOS: brew install gradle
    Ubuntu: sudo apt install gradle

2 of 4 tools available
```

### check-tools Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All tools available |
| 1 | Some tools missing |

---

## init-config

Generate a configuration template file with all available options documented.

### init-config Synopsis

```bash
generator init-config [output-file] [flags]
```

### init-config Description

Creates a YAML configuration file that you can customize for your project. The template includes examples for all supported component types and integration options.

### init-config Flags

```bash
--minimal          # Generate minimal configuration (only required fields)
--example TYPE     # Generate example configuration (fullstack, frontend, backend, mobile)
```

### init-config Examples

```bash
# Generate default configuration
generator init-config

# Generate with custom filename
generator init-config my-project.yaml

# Generate minimal configuration
generator init-config --minimal

# Generate example configurations
generator init-config --example fullstack
generator init-config --example frontend
generator init-config --example backend
generator init-config --example mobile
```

### init-config Example Types

| Type | Description |
|------|-------------|
| `fullstack` | Next.js frontend + Go backend + Docker |
| `frontend` | Next.js frontend only |
| `backend` | Go backend only |
| `mobile` | Android + iOS apps |

---

## cache-tools

Manage the tool availability cache for offline project generation.

### cache-tools Synopsis

```bash
generator cache-tools [flags]
```

### cache-tools Description

Manages cached tool availability information to speed up tool checks and enable offline operation. The cache stores tool availability and version information.

### cache-tools Flags

```bash
--stats      # Show cache statistics (default if no flags)
--save       # Save current tool availability to cache
--clear      # Clear the tool cache
--info       # Show cache information and location
--verbose, -v # Enable verbose logging
```

### cache-tools Examples

```bash
# View cache statistics
generator cache-tools
generator cache-tools --stats

# Save current tool availability
generator cache-tools --save

# Clear cache
generator cache-tools --clear

# Show cache information
generator cache-tools --info
```

### cache-tools Cache Statistics Output

```text
TOOL CACHE STATISTICS

Total Entries: 4
Available Tools: 2
Unavailable Tools: 2
Expired Entries: 0

Cache File: ~/.cache/generator/tools.cache
Cache TTL: 24h0m0s
Last Saved: 2024-01-15T10:30:00Z
Time Since Check: 2h15m30s
```

### cache-tools What Gets Cached

**Cached:**

- Tool availability (whether tool is in PATH)
- Tool versions
- Last check timestamp

**Not Cached:**

- Tool executables themselves
- Tool dependencies
- Network-based resources

---

## version

Print version information about the generator.

### version Synopsis

```bash
generator version
```

### version Description

Displays version information including build time and git commit.

### version Output

```text
Open Source Project Generator
Version: 1.0.0
Built: 2024-01-15T10:00:00Z
Commit: abc123def456
```

---

## Environment Variables

Override configuration with environment variables:

### Project Configuration

```bash
GENERATOR_PROJECT_NAME        # Project name
GENERATOR_OUTPUT_DIR          # Output directory
GENERATOR_CONFIG_FILE         # Configuration file path
```

### Generation Options

```bash
GENERATOR_USE_EXTERNAL_TOOLS  # Use bootstrap tools (true/false)
GENERATOR_OFFLINE             # Offline mode (true/false)
GENERATOR_DRY_RUN             # Dry run mode (true/false)
GENERATOR_VERBOSE             # Verbose output (true/false)
GENERATOR_FORCE               # Force overwrite (true/false)
GENERATOR_CREATE_BACKUP       # Create backup (true/false)
```

### Example Usage

```bash
# Set environment variables
export GENERATOR_PROJECT_NAME="my-app"
export GENERATOR_VERBOSE=true

# Generate with environment variables
generator generate --config project.yaml
```

---

## Common Workflows

### First-Time Setup

```bash
# 1. Check available tools
generator check-tools

# 2. Install missing tools (if needed)
brew install node go  # macOS

# 3. Create configuration
generator init-config my-project.yaml

# 4. Edit configuration
vim my-project.yaml

# 5. Generate project
generator generate --config my-project.yaml
```

### Offline Preparation

```bash
# While online:
# 1. Check and cache tools
generator check-tools
generator cache-tools --save

# Later, offline:
# 2. Generate project
generator generate --config project.yaml --offline
```

### CI/CD Pipeline

```bash
# In your CI/CD script:
# 1. Check tools
generator check-tools || exit 1

# 2. Generate project
generator generate \
  --config .generator/project.yaml \
  --output ./generated \
  --verbose

# 3. Validate generated project
cd generated && make test
```

### Debugging Issues

```bash
# 1. Enable verbose output
generator generate --config project.yaml --verbose

# 2. Try dry run
generator generate --config project.yaml --dry-run

# 3. Check tool availability
generator check-tools --verbose

# 4. Try with fallback
generator generate --config project.yaml --no-external-tools
```

---

## Tips and Tricks

### Preview Before Generating

Always use `--dry-run` to preview what will be generated:

```bash
generator generate --config project.yaml --dry-run
```

### Check Tools First

Before generating, verify tools are available:

```bash
generator check-tools
```

### Use Verbose Mode for Debugging

When troubleshooting, enable verbose output:

```bash
generator generate --config project.yaml --verbose
```

### Cache for Offline Use

Prepare for offline work by caching tool information:

```bash
generator cache-tools --save
```

### Version Control Your Configs

Save configuration files in version control:

```bash
mkdir -p .generator
generator init-config .generator/project.yaml
git add .generator/
```

---

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [Configuration Guide](CONFIGURATION.md) - Configuration file format
- [Examples](EXAMPLES.md) - Example configurations
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
- [Architecture](ARCHITECTURE.md) - How it works
