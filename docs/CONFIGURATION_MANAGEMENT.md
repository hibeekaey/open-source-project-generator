# Configuration Management

The Open Source Project Generator provides comprehensive configuration management capabilities that allow you to save, load, and reuse project generation settings. This eliminates the need to repeatedly enter the same information for similar projects.

## Overview

Configuration management enables you to:

- **Save configurations** after interactive project setup
- **Load saved configurations** to skip repetitive setup steps
- **Manage configurations** through CLI commands or interactive interface
- **Export/import configurations** for sharing or backup
- **Organize configurations** with tags and descriptions

## Configuration Storage

Configurations are stored as YAML files in your user configuration directory:

- **Linux/macOS**: `~/.generator/configs/`
- **Windows**: `%USERPROFILE%\.generator\configs\`

Each configuration includes:

- Project metadata (name, description, author, license, etc.)
- Selected templates and their options
- Generation settings (include examples, tests, documentation, etc.)
- User preferences and defaults

## Basic Usage

### Saving Configurations

#### During Interactive Generation

```bash
# The system will automatically offer to save your configuration
generator generate
# ... complete interactive setup ...
# At the end: "Would you like to save this configuration for future use?"
```

#### Using CLI Flags

```bash
# Save configuration with a specific name during generation
generator generate --save-config my-api-project
```

### Loading Configurations

#### Interactive Loading

```bash
# Load a saved configuration interactively
generator generate --load-config my-api-project
```

#### List Available Configurations

```bash
# Show all saved configurations
generator config list

# Search configurations
generator config list --search "api"

# Filter by tags
generator config list --tags backend,go
```

### Managing Configurations

#### View Configuration Details

```bash
# View detailed information about a configuration
generator config view my-api-project
```

#### Delete Configurations

```bash
# Delete a configuration (with confirmation)
generator config delete old-project

# Delete without confirmation
generator config delete old-project --force
```

#### Interactive Management

```bash
# Launch interactive configuration management interface
generator config manage
```

## Advanced Features

### Export and Import

#### Export Configuration

```bash
# Export to YAML file
generator config export my-config --output my-config.yaml

# Export to JSON file
generator config export my-config --format json --output my-config.json

# Export to stdout
generator config export my-config --format yaml
```

#### Import Configuration

```bash
# Import from file
generator config import --file my-config.yaml --name imported-config

# Import with automatic name detection
generator config import --file my-config.yaml

# Overwrite existing configuration
generator config import --file my-config.yaml --name existing-config --force
```

### Configuration Filtering and Sorting

```bash
# List configurations sorted by creation date
generator config list --sort-by created_at --sort-order desc

# Search in names and descriptions
generator config list --search "frontend"

# Filter by multiple tags
generator config list --tags frontend,react,typescript

# Limit results with pagination
generator config list --limit 10 --offset 20
```

### JSON Output for Automation

```bash
# Get configuration list in JSON format
generator config list --output-format json

# View configuration details in JSON
generator config view my-config --output-format json
```

## Configuration Structure

### Saved Configuration Format

```yaml
name: my-api-project
description: RESTful API server with authentication
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
version: 1.0.0
tags:
  - backend
  - api
  - go

project_config:
  name: my-awesome-api
  organization: my-company
  author: John Doe
  email: john@example.com
  license: MIT
  description: A RESTful API server with JWT authentication
  repository: https://github.com/my-company/my-awesome-api

selected_templates:
  - template_name: go-gin
    category: backend
    technology: go
    version: 1.0.0
    selected: true
    options:
      database: postgresql
      auth: jwt
  - template_name: docker
    category: infrastructure
    technology: docker
    version: 1.0.0
    selected: true

generation_settings:
  include_examples: true
  include_tests: true
  include_docs: true
  update_versions: false
  minimal_mode: false
  backup_existing: true
  overwrite_existing: false

user_preferences:
  default_license: MIT
  default_author: John Doe
  default_email: john@example.com
  default_organization: my-company
  preferred_format: yaml
```

## Interactive Configuration Management

The interactive management interface provides a menu-driven approach to configuration management:

```bash
generator config manage
```

This launches an interactive interface with options to:

1. **List Configurations** - View all saved configurations in a table
2. **View Configuration** - See detailed information about a specific configuration
3. **Delete Configuration** - Remove configurations you no longer need
4. **Export Configuration** - Save configurations to files for backup or sharing
5. **Import Configuration** - Load configurations from files

### Navigation

- Use **arrow keys** to navigate menus
- Press **Enter** to select options
- Press **b** to go back
- Press **q** to quit
- Press **h** for help

## Integration with Project Generation

### Loading During Generation

When you load a saved configuration, you can:

1. **Use as-is** - Generate the project with the exact saved settings
2. **Modify before use** - Update project details while keeping template selections
3. **Preview changes** - See what will be generated before proceeding

### Workflow Example

```bash
# 1. Load a saved configuration
generator generate --load-config my-api-template

# 2. System shows configuration preview
# 3. Option to modify settings (project name, output directory, etc.)
# 4. Preview project structure
# 5. Confirm and generate
```

## Best Practices

### Naming Conventions

- Use descriptive names: `go-api-with-auth` instead of `config1`
- Include technology stack: `nextjs-frontend`, `go-backend`
- Use consistent naming patterns within your team

### Organization with Tags

```bash
# Tag by technology
--tags go,backend,api

# Tag by project type
--tags frontend,admin,dashboard

# Tag by team or purpose
--tags team-alpha,prototype,production
```

### Configuration Maintenance

- **Regular cleanup**: Remove outdated configurations
- **Export important configs**: Backup configurations you want to keep long-term
- **Share team configs**: Export and share configurations within your team
- **Version control**: Store exported configurations in your project repositories

## Automation and CI/CD

### Non-Interactive Usage

```bash
# Generate project using saved configuration in CI/CD
generator generate --load-config production-api --non-interactive --output ./generated-project
```

### Batch Operations

```bash
# Export all configurations for backup
for config in $(generator config list --output-format json | jq -r '.[].name'); do
  generator config export "$config" --output "backup/${config}.yaml"
done
```

### Environment-Specific Configurations

Create different configurations for different environments:

- `my-project-dev` - Development settings with examples and verbose logging
- `my-project-prod` - Production settings with minimal output and security focus
- `my-project-test` - Testing configuration with comprehensive test suites

## Troubleshooting

### Common Issues

#### Configuration Not Found

```bash
Error: configuration 'my-config' not found
```

**Solution**: Check available configurations with `generator config list`

#### Permission Errors

```bash
Error: failed to create config directory: permission denied
```

**Solution**: Ensure you have write permissions to your home directory

#### Invalid Configuration Format

```bash
Error: failed to unmarshal configuration: yaml: line 5: mapping values are not allowed in this context
```

**Solution**: Check YAML syntax in imported configuration files

### Debug Mode

Use debug mode to troubleshoot configuration issues:

```bash
generator config list --debug
generator config view my-config --verbose
```

### Configuration Directory

Check your configuration directory:

```bash
# Linux/macOS
ls -la ~/.generator/configs/

# Windows
dir %USERPROFILE%\.generator\configs\
```

## Security Considerations

### Sensitive Information

- **Don't store secrets** in configurations (API keys, passwords, etc.)
- **Review before sharing** - configurations may contain email addresses or organization names
- **Use environment variables** for sensitive data in generated projects

### File Permissions

Configuration files are created with restricted permissions (0644) to prevent unauthorized access.

### Backup and Recovery

- **Regular exports**: Export important configurations periodically
- **Version control**: Store configuration exports in version control
- **Team sharing**: Use secure channels to share configuration files

## API Reference

### Configuration Persistence

The configuration system provides programmatic access through the `ConfigurationPersistence` interface:

```go
// Create persistence manager
persistence := config.NewConfigurationPersistence(configDir, logger)

// Save configuration
err := persistence.SaveConfiguration(name, config)

// Load configuration
config, err := persistence.LoadConfiguration(name)

// List configurations
configs, err := persistence.ListConfigurations(options)

// Delete configuration
err := persistence.DeleteConfiguration(name)
```

### Interactive Management

The interactive configuration manager provides UI-based configuration management:

```go
// Create interactive manager
manager := ui.NewInteractiveConfigurationManager(ui, configDir, logger)

// Save configuration interactively
name, err := manager.SaveConfigurationInteractively(ctx, projectConfig, templates)

// Load configuration interactively
config, err := manager.LoadConfigurationInteractively(ctx, options)

// Manage configurations interactively
err := manager.ManageConfigurationsInteractively(ctx)
```

## Examples

### Complete Workflow Example

```bash
# 1. Create a new project interactively
generator generate
# ... complete setup for a Go API project ...
# Save as: "go-api-template"

# 2. Later, create similar project
generator generate --load-config go-api-template
# Modify project name and output directory
# Generate new project with same template selections

# 3. Share configuration with team
generator config export go-api-template --output team-configs/go-api.yaml

# 4. Team member imports configuration
generator config import --file team-configs/go-api.yaml --name team-go-api

# 5. Use in CI/CD pipeline
generator generate --load-config team-go-api --non-interactive --output ./api-service
```

### Configuration Templates for Different Project Types

#### Backend API Configuration

```yaml
name: backend-api-template
description: Standard backend API with authentication and database
tags: [backend, api, go, postgresql]
# ... rest of configuration
```

#### Frontend Application Configuration

```yaml
name: frontend-app-template
description: React frontend with TypeScript and Tailwind CSS
tags: [frontend, react, typescript, tailwind]
# ... rest of configuration
```

#### Full-Stack Configuration

```yaml
name: fullstack-template
description: Complete full-stack application with backend and frontend
tags: [fullstack, go, react, postgresql]
# ... rest of configuration
```

This configuration management system provides a powerful way to standardize and accelerate project generation while maintaining flexibility for customization.
