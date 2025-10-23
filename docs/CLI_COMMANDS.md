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
--interactive, -i          # Interactive mode (guided configuration wizard)
--stream-output            # Stream real-time output from bootstrap tools
--no-rollback              # Skip automatic rollback on failure
```

### generate Examples

```bash
# Generate from config file
generator generate --config project.yaml

# Interactive mode (guided wizard)
generator generate --interactive

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

# Stream real-time output from bootstrap tools
generator generate --config project.yaml --stream-output --verbose

# Force overwrite
generator generate --config project.yaml --force

# Skip automatic rollback on failure
generator generate --config project.yaml --no-rollback
```

### generate Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Project generated successfully |
| 1 | Configuration Error | Invalid or missing configuration |
| 2 | Tools Missing | Required tools not found and no fallback available |
| 3 | Generation Failed | Component generation or tool execution failed |
| 4 | File System Error | Cannot read/write files or directories |
| 5 | User Cancelled | User cancelled the operation (e.g., in interactive mode) |

**Exit Code Details:**

- **Exit Code 0 (Success)**: All components generated successfully, project structure validated
- **Exit Code 1 (Configuration Error)**: Configuration file not found, invalid YAML/JSON syntax, missing required fields, or invalid component configuration
- **Exit Code 2 (Tools Missing)**: Required bootstrap tools (npx, go, gradle, xcodebuild) not found in PATH and no fallback generator available
- **Exit Code 3 (Generation Failed)**: Bootstrap tool execution failed, component generation error, or integration step failed
- **Exit Code 4 (File System Error)**: Cannot create output directory, permission denied, disk full, or file operation failed
- **Exit Code 5 (User Cancelled)**: User cancelled operation in interactive mode or via Ctrl+C

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
--example TYPE     # Generate example configuration (fullstack, frontend, backend, mobile, microservice)
--force            # Force overwrite existing configuration file
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
generator init-config --example microservice

# Force overwrite existing file
generator init-config --force
generator init-config my-project.yaml --example fullstack --force
```

### init-config Example Types

| Type | Description | Components |
|------|-------------|------------|
| `fullstack` | Full-stack web application | Next.js frontend + Go backend + Docker |
| `frontend` | Frontend-only application | Next.js with TypeScript and Tailwind |
| `backend` | Backend-only API service | Go backend with Gin framework |
| `mobile` | Mobile applications | Android (Kotlin) + iOS (Swift) + API backend |
| `microservice` | Microservice architecture | Go backend optimized for microservices |

All generated templates include inline comments explaining each configuration option.

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
--stats          # Show cache statistics (default if no flags)
--save           # Save current tool availability to cache
--clear          # Clear the tool cache
--info           # Show cache information and location
--validate       # Validate cache integrity and report issues
--refresh        # Re-check all cached tools and update status
--export FILE    # Export cache data to portable format
--import FILE    # Import cache data from file
--verbose, -v    # Enable verbose logging
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

# Validate cache integrity
generator cache-tools --validate

# Refresh all cached tools
generator cache-tools --refresh

# Export cache for sharing
generator cache-tools --export tools-cache.json

# Import cache from file
generator cache-tools --import tools-cache.json

# Verbose output
generator cache-tools --validate --verbose
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

### cache-tools Validation Output

```text
CACHE VALIDATION REPORT

Status: Valid
Total Entries: 4
Corrupted Entries: 0
Expired Entries: 1

Expired:
  - gradle (last checked: 2024-01-10T10:00:00Z)

Warnings:
  - Cache is older than 7 days, consider refreshing

Checked At: 2024-01-15T10:30:00Z
```

### cache-tools Export Format

The export format is a portable JSON file that can be shared across machines:

```json
{
  "version": "1.0",
  "exported_at": "2024-01-15T10:30:00Z",
  "platform": "darwin",
  "entries": {
    "npx": {
      "available": true,
      "version": "10.2.3",
      "last_checked": "2024-01-15T10:00:00Z"
    },
    "go": {
      "available": true,
      "version": "1.21.0",
      "last_checked": "2024-01-15T10:00:00Z"
    }
  }
}
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

## Interactive Mode

The interactive mode provides a guided wizard for configuring your project without manually creating a configuration file.

### Starting Interactive Mode

```bash
generator generate --interactive
# or
generator generate -i
```

### Interactive Mode Workflow

The wizard guides you through these steps:

**1. Project Information**
- Project name
- Project description
- Output directory
- Author (optional)
- License (optional)

**2. Component Selection**
- Multi-select menu of available component types
- Descriptions for each component
- Select one or more components to include

**3. Component Configuration**
- For each selected component, configure specific options
- TypeScript, Tailwind, App Router for Next.js
- Module path, framework, port for Go backend
- Package name, SDK levels for Android
- Bundle ID, deployment target for iOS

**4. Integration Options**
- Docker Compose generation
- Build scripts generation
- API endpoint configuration
- Shared environment variables

**5. Confirmation**
- Review complete configuration summary
- Confirm to proceed or cancel to restart

### Interactive Mode Features

- **Input Validation**: Real-time validation of all inputs
- **Smart Defaults**: Sensible defaults for all options
- **Cancellation**: Press Ctrl+C at any prompt to cancel
- **Error Recovery**: Clear error messages with re-prompting
- **Configuration Preview**: See complete configuration before generation

### Interactive Mode Examples

```bash
# Start interactive wizard
generator generate --interactive

# Interactive with streaming output
generator generate --interactive --stream-output

# Interactive with verbose logging
generator generate --interactive --verbose
```

### Interactive Mode Exit Codes

Interactive mode uses the same exit codes as non-interactive mode:

- **0**: Project generated successfully
- **1**: Configuration validation failed
- **2**: Required tools missing
- **3**: Generation failed
- **4**: File system error
- **5**: User cancelled operation

---

## Common Workflows

### Interactive Project Creation

```bash
# 1. Start interactive mode
generator generate --interactive

# 2. Follow the prompts:
#    - Enter project name: "my-awesome-app"
#    - Enter description: "A full-stack web application"
#    - Select output directory: "./my-awesome-app"
#    - Select components: [nextjs, go-backend]
#    - Configure Next.js: TypeScript=yes, Tailwind=yes
#    - Configure Go backend: Module path, Port=8080
#    - Enable Docker Compose: yes
#    - Confirm and generate

# 3. Project is generated automatically
cd my-awesome-app
```

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

# 5. Stream real-time output
generator generate --config project.yaml --stream-output --verbose
```

### Streaming Output

View real-time output from bootstrap tools during generation:

```bash
# Enable streaming output
generator generate --config project.yaml --stream-output

# With verbose mode for detailed output
generator generate --config project.yaml --stream-output --verbose
```

**Streaming Output Features:**

- Real-time display of tool output (npx, go, gradle, etc.)
- Component name prefixing for clarity
- Both stdout and stderr in verbose mode
- Immediate error display
- Progress indicators for non-streaming operations

**Example Output:**

```text
[nextjs] Creating Next.js application...
[nextjs] ✓ Creating project directory
[nextjs] ✓ Installing dependencies
[nextjs] ✓ Initializing git repository
[nextjs] Success! Created web-app

[go-backend] Initializing Go module...
[go-backend] go: creating new go.mod: module github.com/myorg/project
[go-backend] go: to add module requirements and sums:
[go-backend]     go mod tidy
[go-backend] Success! Initialized Go module
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

## Exit Codes Reference

All commands return specific exit codes for different scenarios. Use these in CI/CD pipelines and scripts.

### Exit Code Summary

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| 0 | Success | Operation completed successfully | - |
| 1 | Configuration Error | Invalid or missing configuration | Missing config file, invalid YAML, missing required fields |
| 2 | Tools Missing | Required tools not available | npx, go, gradle, or xcodebuild not in PATH |
| 3 | Generation Failed | Component generation failed | Tool execution error, network issue, invalid options |
| 4 | File System Error | Cannot access files/directories | Permission denied, disk full, invalid path |
| 5 | User Cancelled | User cancelled the operation | Ctrl+C pressed, cancelled in interactive mode |

### Using Exit Codes in Scripts

**Bash Script Example:**

```bash
#!/bin/bash

# Generate project
generator generate --config project.yaml

# Check exit code
case $? in
  0)
    echo "✓ Project generated successfully"
    cd my-project && make build
    ;;
  1)
    echo "✗ Configuration error - check your config file"
    exit 1
    ;;
  2)
    echo "✗ Missing tools - install required tools"
    generator check-tools
    exit 2
    ;;
  3)
    echo "✗ Generation failed - check logs"
    cat ~/.cache/generator/logs/generator.log
    exit 3
    ;;
  4)
    echo "✗ File system error - check permissions"
    exit 4
    ;;
  5)
    echo "✗ Operation cancelled by user"
    exit 5
    ;;
esac
```

**CI/CD Pipeline Example:**

```yaml
# GitHub Actions
- name: Generate Project
  run: generator generate --config .generator/project.yaml
  continue-on-error: false

- name: Handle Generation Failure
  if: failure()
  run: |
    echo "Generation failed with exit code $?"
    generator check-tools
    cat ~/.cache/generator/logs/generator.log
```

**Makefile Example:**

```makefile
.PHONY: generate
generate:
	@generator generate --config project.yaml || \
	(echo "Generation failed with exit code $$?"; exit 1)

.PHONY: generate-or-fallback
generate-or-fallback:
	@generator generate --config project.yaml || \
	(echo "Trying with fallback..."; \
	 generator generate --config project.yaml --no-external-tools)
```

### Exit Code Handling Best Practices

1. **Always check exit codes** in automated scripts
2. **Log exit codes** for debugging
3. **Provide fallback options** for non-zero exits
4. **Use specific error handling** for each exit code
5. **Display helpful messages** based on exit code

---

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [Configuration Guide](CONFIGURATION.md) - Configuration file format
- [Examples](EXAMPLES.md) - Example configurations
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
- [Architecture](ARCHITECTURE.md) - How it works
