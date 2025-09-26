# Configuration Guide

This guide covers all aspects of configuring the Open Source Project Generator, from basic settings to advanced customization.

## Table of Contents

- [Configuration Sources](#configuration-sources)
- [Configuration File Format](#configuration-file-format)
- [Environment Variables](#environment-variables)
- [Global Settings](#global-settings)
- [Template Configuration](#template-configuration)
- [Validation Rules](#validation-rules)
- [Advanced Configuration](#advanced-configuration)

## Configuration Sources

The generator loads configuration from multiple sources in order of precedence:

1. **Command-line arguments** (highest priority)
2. **Environment variables**
3. **Configuration files**
4. **Default values** (lowest priority)

### Configuration File Locations

The generator looks for configuration files in the following locations:

1. `./generator.yaml` (current directory)
2. `~/.config/generator/config.yaml` (Linux/macOS)
3. `%APPDATA%\generator\config.yaml` (Windows)
4. `/etc/generator/config.yaml` (system-wide)

### Environment Variable Override

You can override any configuration using environment variables:

```bash
# Override configuration file settings
export GENERATOR_PROJECT_NAME="my-project"
export GENERATOR_FORCE=true
export GENERATOR_OFFLINE=true
```

## Configuration File Format

### YAML Configuration

```yaml
# Global settings
default_output_dir: "~/projects"
default_license: "MIT"
default_organization: "myorg"

# Template preferences
templates:
  frontend: "nextjs-app"
  backend: "go-gin"
  mobile: "android-kotlin"

# Version preferences
versions:
  node: "20"
  go: "1.25"
  react: "19"
  typescript: "5"

# Cache settings
cache:
  ttl: "24h"
  enabled: true
  offline_mode: false

# Validation settings
validation:
  strict: false
  auto_fix: true
  rules:
    - security
    - quality
    - compliance

# UI settings
ui:
  enable_colors: true
  enable_unicode: true
  page_size: 10
  timeout: "30m"
  auto_save: true
  show_breadcrumbs: true
  show_shortcuts: true
  confirm_on_quit: true

# Logging settings
logging:
  level: "info"
  format: "text"
  file: "~/.cache/generator/logs/generator.log"
  max_size: "100MB"
  max_files: 5
  compress: true
```

### JSON Configuration

```json
{
  "default_output_dir": "~/projects",
  "default_license": "MIT",
  "default_organization": "myorg",
  "templates": {
    "frontend": "nextjs-app",
    "backend": "go-gin",
    "mobile": "android-kotlin"
  },
  "versions": {
    "node": "20",
    "go": "1.25",
    "react": "19",
    "typescript": "5"
  },
  "cache": {
    "ttl": "24h",
    "enabled": true,
    "offline_mode": false
  },
  "validation": {
    "strict": false,
    "auto_fix": true,
    "rules": ["security", "quality", "compliance"]
  },
  "ui": {
    "enable_colors": true,
    "enable_unicode": true,
    "page_size": 10,
    "timeout": "30m",
    "auto_save": true,
    "show_breadcrumbs": true,
    "show_shortcuts": true,
    "confirm_on_quit": true
  },
  "logging": {
    "level": "info",
    "format": "text",
    "file": "~/.cache/generator/logs/generator.log",
    "max_size": "100MB",
    "max_files": 5,
    "compress": true
  }
}
```

## Environment Variables

### Project Configuration

```bash
# Required for non-interactive mode
export GENERATOR_PROJECT_NAME="my-awesome-project"

# Optional project metadata
export GENERATOR_PROJECT_ORGANIZATION="my-org"
export GENERATOR_PROJECT_DESCRIPTION="An awesome project"
export GENERATOR_PROJECT_LICENSE="MIT"
export GENERATOR_PROJECT_AUTHOR="John Doe"
export GENERATOR_PROJECT_EMAIL="john@example.com"
export GENERATOR_PROJECT_REPOSITORY="https://github.com/my-org/my-awesome-project"
export GENERATOR_OUTPUT_PATH="./output"
```

### Generation Options

```bash
# Generation behavior
export GENERATOR_FORCE=true                # Overwrite existing files
export GENERATOR_MINIMAL=false             # Generate full project structure
export GENERATOR_OFFLINE=false             # Use cached data only
export GENERATOR_UPDATE_VERSIONS=true      # Fetch latest package versions
export GENERATOR_SKIP_VALIDATION=false     # Skip configuration validation
export GENERATOR_BACKUP_EXISTING=true      # Backup existing files
export GENERATOR_INCLUDE_EXAMPLES=true     # Include example code
export GENERATOR_TEMPLATE="go-gin"         # Specific template to use
export GENERATOR_DRY_RUN=false            # Preview without creating files
export GENERATOR_NON_INTERACTIVE=false     # Non-interactive mode
```

### Component Selection

```bash
# Enable/disable components
export GENERATOR_FRONTEND=true
export GENERATOR_BACKEND=true
export GENERATOR_MOBILE=false
export GENERATOR_INFRASTRUCTURE=true

# Frontend components
export GENERATOR_FRONTEND_MAIN_APP=true
export GENERATOR_FRONTEND_HOME=true
export GENERATOR_FRONTEND_ADMIN=false
export GENERATOR_FRONTEND_SHARED_COMPONENTS=true

# Backend components
export GENERATOR_BACKEND_API=true
export GENERATOR_BACKEND_AUTH=true
export GENERATOR_BACKEND_DATABASE=true

# Mobile components
export GENERATOR_MOBILE_ANDROID=true
export GENERATOR_MOBILE_IOS=true
export GENERATOR_MOBILE_SHARED=true

# Infrastructure components
export GENERATOR_INFRASTRUCTURE_DOCKER=true
export GENERATOR_INFRASTRUCTURE_KUBERNETES=true
export GENERATOR_INFRASTRUCTURE_TERRAFORM=false
export GENERATOR_INFRASTRUCTURE_MONITORING=true
```

### CLI Behavior

```bash
# Output and logging
export GENERATOR_NON_INTERACTIVE=true
export GENERATOR_OUTPUT_FORMAT="json"      # json, yaml, text
export GENERATOR_LOG_LEVEL="info"          # debug, info, warn, error
export GENERATOR_VERBOSE=false
export GENERATOR_QUIET=false
export GENERATOR_DEBUG=false

# UI settings
export GENERATOR_UI_COLORS=true
export GENERATOR_UI_UNICODE=true
export GENERATOR_UI_PAGE_SIZE=10
export GENERATOR_UI_TIMEOUT="30m"
export GENERATOR_UI_AUTO_SAVE=true
```

### Cache and Offline Mode

```bash
# Cache settings
export GENERATOR_CACHE_ENABLED=true
export GENERATOR_CACHE_TTL="24h"
export GENERATOR_CACHE_OFFLINE_MODE=false

# Offline mode
export GENERATOR_OFFLINE=true
export GENERATOR_UPDATE_TEMPLATES=false
export GENERATOR_UPDATE_PACKAGES=false
```

## Global Settings

### Default Values

```yaml
# Default project settings
defaults:
  license: "MIT"
  organization: "myorg"
  author: "Developer"
  email: "dev@example.com"
  repository: "https://github.com/myorg/project"

# Default generation options
default_generation:
  force: false
  minimal: false
  offline: false
  update_versions: true
  skip_validation: false
  backup_existing: true
  include_examples: true

# Default component selection
default_components:
  frontend:
    main_app: true
    home: false
    admin: false
    shared_components: true
  backend:
    api: true
    auth: false
    database: false
  mobile:
    android: false
    ios: false
    shared: false
  infrastructure:
    docker: true
    kubernetes: false
    terraform: false
    monitoring: false
```

### Template Preferences

```yaml
# Preferred templates for each category
template_preferences:
  frontend:
    default: "nextjs-app"
    alternatives:
      - "react-app"
      - "vue-app"
      - "angular-app"
  backend:
    default: "go-gin"
    alternatives:
      - "node-express"
      - "python-fastapi"
      - "java-spring"
  mobile:
    default: "android-kotlin"
    alternatives:
      - "ios-swift"
      - "react-native"
      - "flutter"
  infrastructure:
    default: "docker"
    alternatives:
      - "kubernetes"
      - "terraform"
      - "helm"
```

### Version Preferences

```yaml
# Preferred versions for each technology
version_preferences:
  node: "20"
  go: "1.25"
  react: "19"
  typescript: "5"
  nextjs: "15"
  tailwindcss: "3.4"
  kotlin: "2.0"
  swift: "5.9"
  docker: "24"
  kubernetes: "1.28"
  terraform: "1.6"
```

## Template Configuration

### Custom Templates

```yaml
# Custom template configuration
custom_templates:
  enabled: true
  directories:
    - "~/templates"
    - "./custom-templates"
  validation:
    strict: true
    auto_fix: true
  caching:
    enabled: true
    ttl: "1h"
```

### Template Variables

```yaml
# Global template variables
template_variables:
  company_name: "MyCompany"
  company_url: "https://mycompany.com"
  support_email: "support@mycompany.com"
  default_port: "8080"
  default_database: "postgresql"
  default_cache: "redis"
```

### Template Functions

```yaml
# Custom template functions
template_functions:
  enabled: true
  functions:
    - name: "formatVersion"
      description: "Format version with caret prefix"
    - name: "generateSecret"
      description: "Generate secure random string"
    - name: "formatDate"
      description: "Format current date"
```

## Validation Rules

### Validation Configuration

```yaml
# Validation settings
validation:
  enabled: true
  strict: false
  auto_fix: true
  fail_on_warnings: false
  
  # Validation rules
  rules:
    - name: "security"
      enabled: true
      severity: "error"
      auto_fix: false
    - name: "quality"
      enabled: true
      severity: "warning"
      auto_fix: true
    - name: "compliance"
      enabled: true
      severity: "warning"
      auto_fix: false
    - name: "structure"
      enabled: true
      severity: "error"
      auto_fix: true
    - name: "dependencies"
      enabled: true
      severity: "warning"
      auto_fix: true

  # Custom validation rules
  custom_rules:
    - name: "company_standards"
      pattern: ".*"
      message: "Must follow company coding standards"
      severity: "warning"
```

### Audit Configuration

```yaml
# Audit settings
audit:
  enabled: true
  security:
    enabled: true
    fail_on_high: true
    fail_on_medium: false
    min_score: 7.0
  quality:
    enabled: true
    min_score: 8.0
    check_complexity: true
    check_duplication: true
  licenses:
    enabled: true
    allowed_licenses:
      - "MIT"
      - "Apache-2.0"
      - "BSD-3-Clause"
    forbidden_licenses:
      - "GPL-3.0"
      - "AGPL-3.0"
  performance:
    enabled: true
    check_bundle_size: true
    max_bundle_size: "1MB"
```

## Advanced Configuration

### Plugin System

```yaml
# Plugin configuration
plugins:
  enabled: true
  directories:
    - "~/.generator/plugins"
    - "./plugins"
  auto_load: true
  plugins:
    - name: "custom-validator"
      enabled: true
      config:
        rules: ["custom-rule-1", "custom-rule-2"]
    - name: "template-processor"
      enabled: true
      config:
        processors: ["pre-processor", "post-processor"]
```

### Integration Settings

```yaml
# External service integration
integrations:
  github:
    enabled: true
    token: "${GITHUB_TOKEN}"
    api_url: "https://api.github.com"
  npm:
    enabled: true
    registry: "https://registry.npmjs.org"
  docker:
    enabled: true
    registry: "docker.io"
  kubernetes:
    enabled: true
    context: "default"
```

### Performance Settings

```yaml
# Performance configuration
performance:
  parallel_processing: true
  max_workers: 4
  memory_limit: "2GB"
  timeout: "30m"
  
  # Caching
  cache:
    enabled: true
    ttl: "24h"
    max_size: "1GB"
    compression: true
  
  # Template processing
  template_processing:
    cache_templates: true
    parallel_rendering: true
    max_concurrent: 8
```

### Security Settings

```yaml
# Security configuration
security:
  # Input validation
  input_validation:
    enabled: true
    max_length: 1000
    allowed_chars: "a-zA-Z0-9-_."
    sanitize_input: true
  
  # Path security
  path_security:
    enabled: true
    allowed_directories:
      - "~/projects"
      - "./output"
    forbidden_patterns:
      - "../"
      - "/etc/"
      - "/root/"
  
  # Template security
  template_security:
    enabled: true
    restricted_functions:
      - "exec"
      - "system"
      - "eval"
    allowed_functions:
      - "toUpper"
      - "toLower"
      - "replace"
```

## Configuration Management

### Configuration Commands

```bash
# Show current configuration
generator config show

# Show specific configuration key
generator config show default.license

# Set configuration values
generator config set default.license MIT
generator config set default.organization "MyCompany"

# Load configuration from file
generator config set --file team-config.yaml --merge

# Export configuration
generator config export my-config.yaml

# Import configuration
generator config import --file my-config.yaml

# Reset to defaults
generator config reset

# Validate configuration
generator config validate
```

### Configuration Validation

```bash
# Validate current configuration
generator config validate

# Validate specific file
generator config validate ./my-config.yaml --strict

# Check configuration syntax
generator config validate --syntax-only

# Validate with custom schema
generator config validate --schema ./custom-schema.json
```

### Configuration Templates

```bash
# Export configuration template
generator config export --template team-defaults.yaml

# Create configuration from template
generator config create --template team-defaults.yaml --name my-project

# List available templates
generator config templates

# Show template details
generator config template info team-defaults
```

## Best Practices

### Configuration Organization

1. **Use environment variables** for sensitive data (tokens, passwords)
2. **Use configuration files** for project-specific settings
3. **Use global configuration** for user preferences
4. **Version control** configuration files for team consistency
5. **Document** custom configurations and their purposes

### Security Considerations

1. **Never commit** sensitive configuration data
2. **Use environment variables** for secrets
3. **Validate** all configuration inputs
4. **Restrict** file system access in configuration
5. **Use** secure defaults for all settings

### Performance Optimization

1. **Enable caching** for frequently used data
2. **Use offline mode** when possible
3. **Configure** appropriate timeouts
4. **Limit** parallel processing based on system resources
5. **Monitor** configuration performance impact

### Team Collaboration

1. **Share** configuration templates
2. **Document** configuration changes
3. **Use** consistent naming conventions
4. **Validate** configurations before sharing
5. **Maintain** configuration version history

This comprehensive configuration guide provides all the information needed to customize and optimize the Open Source Project Generator for your specific needs.
