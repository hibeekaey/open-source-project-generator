# Unified Generation Workflow

## Overview

The Open Source Project Generator now uses a comprehensive generation workflow that ensures both config mode and non-interactive mode follow the same execution path. This provides consistency, maintainability, and a better user experience.

## Architecture

### Unified Workflow Design

Both config mode and non-interactive mode now use the same core workflow:

1. **Configuration Loading** - Different sources, same format
2. **Configuration Validation** - Unified validation logic
3. **Pre-generation Checks** - Consistent safety checks
4. **Project Generation** - Same generation engine
5. **Post-generation Tasks** - Consistent cleanup and setup

### Mode-Specific Differences

The only difference between modes is how the configuration is initialized:

- **Config Mode**: Loads configuration from YAML/JSON file
- **Non-Interactive Mode**: Loads configuration from environment variables

After configuration loading, both modes use the identical `executeUnifiedGenerationWorkflow` function.

## Implementation Details

### Key Functions

#### `executeUnifiedGenerationWorkflow(config, options)`

The core workflow function that:

- Validates configuration
- Determines output path
- Handles offline/online modes
- Performs pre-generation checks
- Executes project generation
- Runs post-generation tasks
- Displays results

#### `loadConfigFromFile(configPath)`

Loads and parses configuration from file for config mode.

#### `loadConfigFromEnvironment()`

Loads configuration from environment variables for non-interactive mode.

### Workflow Steps

1. **Mode Detection**: Determines whether to use config file or environment variables
2. **Configuration Loading**: Loads configuration using appropriate method
3. **Option Merging**: Merges environment variables with command-line options
4. **Unified Execution**: Calls `executeUnifiedGenerationWorkflow`

## Benefits

### Consistency

- Both modes generate identical project structures
- Same validation rules apply to both modes
- Consistent error handling and messaging

### Maintainability

- Single code path for core generation logic
- Easier to add new features that work in both modes
- Reduced code duplication

### User Experience

- Predictable behavior across modes
- Same output format and messaging
- Consistent performance characteristics

## Usage Examples

### Config Mode

```bash
# Create config file
cat > project-config.yaml << EOF
name: "my-project"
organization: "my-org"
components:
  frontend:
    nextjs:
      app: true
  backend:
    go_gin: true
  infrastructure:
    docker: true
EOF

# Generate project
generator generate --config project-config.yaml
```

### Non-Interactive Mode

```bash
# Set environment variables
export GENERATOR_PROJECT_NAME="my-project"
export GENERATOR_FRONTEND=true
export GENERATOR_BACKEND=true
export GENERATOR_INFRASTRUCTURE=true

# Generate project
generator generate --non-interactive
```

Both commands will produce identical results using the same workflow.

## Environment Variables

For non-interactive mode, the following environment variables are supported:

### Project Configuration

- `GENERATOR_PROJECT_NAME` - Project name (required)
- `GENERATOR_PROJECT_ORGANIZATION` - Organization name
- `GENERATOR_PROJECT_DESCRIPTION` - Project description
- `GENERATOR_PROJECT_LICENSE` - License type
- `GENERATOR_OUTPUT_PATH` - Output directory

### Component Selection

- `GENERATOR_FRONTEND` - Enable frontend components
- `GENERATOR_BACKEND` - Enable backend components
- `GENERATOR_MOBILE` - Enable mobile components
- `GENERATOR_INFRASTRUCTURE` - Enable infrastructure components

### Technology Selection

- `GENERATOR_FRONTEND_TECH` - Frontend technology (nextjs-app)
- `GENERATOR_BACKEND_TECH` - Backend technology (go-gin)
- `GENERATOR_MOBILE_TECH` - Mobile technology (android, ios)
- `GENERATOR_INFRASTRUCTURE_TECH` - Infrastructure technology (kubernetes, terraform)

### Generation Options

- `GENERATOR_FORCE` - Overwrite existing files
- `GENERATOR_MINIMAL` - Generate minimal structure
- `GENERATOR_OFFLINE` - Use offline mode
- `GENERATOR_UPDATE_VERSIONS` - Update to latest versions
- `GENERATOR_SKIP_VALIDATION` - Skip validation
- `GENERATOR_BACKUP_EXISTING` - Backup existing files
- `GENERATOR_INCLUDE_EXAMPLES` - Include example code

## CI/CD Integration

The comprehensive workflow automatically detects CI/CD environments and adjusts behavior accordingly:

- Auto-enables non-interactive mode in CI environments
- Provides structured logging for automation
- Supports machine-readable output formats
- Handles exit codes appropriately for automation

### Supported CI/CD Platforms

- GitHub Actions
- GitLab CI
- Jenkins
- Travis CI
- CircleCI
- Azure DevOps
- Bitbucket Pipelines
- AWS CodeBuild

## Error Handling

The comprehensive workflow provides consistent error handling:

- Clear error messages with context
- Appropriate exit codes for automation
- Detailed logging in verbose mode
- Suggestions for common issues

## Future Enhancements

The comprehensive workflow enables future enhancements:

- Template caching and offline mode improvements
- Comprehensive validation and linting
- Plugin system for custom components
- Advanced configuration management
- Performance optimizations

## Migration Notes

Existing code using the old separate workflows will continue to work as the old methods now redirect to the comprehensive workflow. However, new development should use the comprehensive approach directly.
