# Non-Interactive Mode and Automation Guide

This guide explains how to use the Open Source Project Generator in non-interactive mode for automation, CI/CD pipelines, and scripting.

## Overview

The generator supports non-interactive mode for automation scenarios where user input is not available or desired. This mode is automatically detected in CI/CD environments or can be explicitly enabled.

## Automatic Detection

Non-interactive mode is automatically enabled when:

- Running in a CI/CD environment (GitHub Actions, GitLab CI, Jenkins, etc.)
- Input is piped (not from a terminal)
- The `--non-interactive` flag is used
- The `GENERATOR_NON_INTERACTIVE=true` environment variable is set

## Environment Variables

Configure the generator using environment variables for non-interactive operation:

### Project Configuration

```bash
# Required
export GENERATOR_PROJECT_NAME="my-awesome-project"

# Optional
export GENERATOR_PROJECT_ORGANIZATION="my-org"
export GENERATOR_PROJECT_DESCRIPTION="An awesome project"
export GENERATOR_PROJECT_LICENSE="MIT"
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
```

### Component Selection

```bash
# Enable/disable components
export GENERATOR_FRONTEND=true
export GENERATOR_BACKEND=true
export GENERATOR_MOBILE=false
export GENERATOR_INFRASTRUCTURE=true

# Technology selection
export GENERATOR_FRONTEND_TECH="nextjs-app"
export GENERATOR_BACKEND_TECH="go-gin"
export GENERATOR_MOBILE_TECH="android-kotlin"
export GENERATOR_INFRASTRUCTURE_TECH="docker"
```

### CLI Behavior

```bash
# Output and logging
export GENERATOR_NON_INTERACTIVE=true
export GENERATOR_OUTPUT_FORMAT="json"      # json, yaml, text
export GENERATOR_LOG_LEVEL="info"          # debug, info, warn, error
export GENERATOR_VERBOSE=false
export GENERATOR_QUIET=false
```

## Machine-Readable Output

Use JSON or YAML output for parsing results in scripts:

```bash
# JSON output
generator list-templates --output-format json

# YAML output  
generator validate ./project --output-format yaml

# Pipe to jq for processing
generator audit ./project --output-format json | jq '.overall_score'
```

## Exit Codes

The generator returns specific exit codes for automation:

- `0` - Success
- `1` - General error
- `2` - Validation failed
- `3` - Configuration invalid
- `4` - Template not found
- `5` - Network error
- `6` - File system error
- `7` - Permission denied
- `8` - Cache error
- `9` - Version error
- `10` - Audit failed
- `11` - Generation failed
- `99` - Internal error

## CI/CD Examples

### GitHub Actions

```yaml
name: Generate Project
on: [push]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Generate Project
        env:
          GENERATOR_PROJECT_NAME: ${{ github.event.repository.name }}
          GENERATOR_PROJECT_ORGANIZATION: ${{ github.repository_owner }}
          GENERATOR_BACKEND: true
          GENERATOR_FRONTEND: true
          GENERATOR_OUTPUT_FORMAT: json
        run: |
          generator generate --non-interactive
          
      - name: Validate Generated Project
        run: |
          generator validate ./output --output-format json --report
          
      - name: Audit Security
        run: |
          generator audit ./output --security --fail-on-high --output-format json
```

### GitLab CI

```yaml
generate_project:
  stage: generate
  variables:
    GENERATOR_PROJECT_NAME: "my-project"
    GENERATOR_BACKEND: "true"
    GENERATOR_FRONTEND: "true"
    GENERATOR_OUTPUT_FORMAT: "json"
  script:
    - generator generate --non-interactive
    - generator validate ./output --output-format json
    - generator audit ./output --output-format json
  artifacts:
    reports:
      junit: validation-report.xml
    paths:
      - output/
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any
    
    environment {
        GENERATOR_PROJECT_NAME = "${env.JOB_NAME}"
        GENERATOR_BACKEND = "true"
        GENERATOR_OUTPUT_FORMAT = "json"
    }
    
    stages {
        stage('Generate') {
            steps {
                sh 'generator generate --non-interactive'
            }
        }
        
        stage('Validate') {
            steps {
                sh 'generator validate ./output --output-format json --report'
            }
        }
        
        stage('Audit') {
            steps {
                sh '''
                    generator audit ./output \\
                        --output-format json \\
                        --fail-on-high \\
                        --min-score 7.0
                '''
            }
        }
    }
    
    post {
        always {
            archiveArtifacts artifacts: 'output/**', fingerprint: true
        }
    }
}
```

## Docker Usage

```dockerfile
FROM generator:latest

# Set environment variables
ENV GENERATOR_PROJECT_NAME=my-project
ENV GENERATOR_BACKEND=true
ENV GENERATOR_OUTPUT_FORMAT=json
ENV GENERATOR_NON_INTERACTIVE=true

# Generate project
RUN generator generate --output /app

WORKDIR /app
```

## Bash Scripting

```bash
#!/bin/bash

set -e  # Exit on error

# Configuration
PROJECT_NAME="my-new-project"
OUTPUT_DIR="./generated"

# Set environment variables
export GENERATOR_PROJECT_NAME="$PROJECT_NAME"
export GENERATOR_BACKEND=true
export GENERATOR_FRONTEND=true
export GENERATOR_OUTPUT_FORMAT=json
export GENERATOR_NON_INTERACTIVE=true

# Generate project
echo "Generating project: $PROJECT_NAME"
if generator generate --output "$OUTPUT_DIR"; then
    echo "✓ Project generated successfully"
else
    echo "✗ Project generation failed"
    exit 1
fi

# Validate project
echo "Validating project structure..."
if generator validate "$OUTPUT_DIR" --output-format json > validation.json; then
    echo "✓ Project validation passed"
else
    echo "✗ Project validation failed"
    cat validation.json
    exit 2
fi

# Audit project
echo "Auditing project security and quality..."
if generator audit "$OUTPUT_DIR" --output-format json --min-score 7.0 > audit.json; then
    echo "✓ Project audit passed"
    SCORE=$(jq -r '.overall_score' audit.json)
    echo "Overall score: $SCORE"
else
    echo "✗ Project audit failed"
    cat audit.json
    exit 10
fi

echo "✓ All checks passed! Project ready for use."
```

## Error Handling

Handle errors gracefully in scripts:

```bash
#!/bin/bash

# Function to handle errors
handle_error() {
    local exit_code=$1
    local command=$2
    
    case $exit_code in
        2)
            echo "Validation failed. Check project structure."
            ;;
        3)
            echo "Configuration invalid. Check environment variables."
            ;;
        4)
            echo "Template not found. Check template name."
            ;;
        10)
            echo "Audit failed. Security or quality issues found."
            ;;
        *)
            echo "Command failed with exit code: $exit_code"
            ;;
    esac
}

# Run command with error handling
run_with_error_handling() {
    local command="$1"
    
    if ! $command; then
        local exit_code=$?
        handle_error $exit_code "$command"
        exit $exit_code
    fi
}

# Usage
run_with_error_handling "generator generate --non-interactive"
run_with_error_handling "generator validate ./output"
run_with_error_handling "generator audit ./output --min-score 8.0"
```

## Best Practices

1. **Always set required environment variables** before running in non-interactive mode
2. **Use specific exit code handling** for different error scenarios
3. **Enable JSON/YAML output** for parsing results in scripts
4. **Set appropriate log levels** for debugging vs. production
5. **Use --dry-run** to preview operations before execution
6. **Cache templates and versions** for faster CI/CD runs
7. **Validate configuration** before generation in CI pipelines
8. **Set minimum audit scores** for quality gates
9. **Archive generated artifacts** for debugging and deployment
10. **Use structured error output** for better error handling

## Troubleshooting

### Common Issues

1. **Missing required environment variables**

   ```bash
   # Error: project name is required in non-interactive mode
   export GENERATOR_PROJECT_NAME="my-project"
   ```

2. **Permission errors**

   ```bash
   # Ensure output directory is writable
   mkdir -p ./output
   chmod 755 ./output
   ```

3. **Network issues in CI**

   ```bash
   # Use offline mode with cached data
   export GENERATOR_OFFLINE=true
   ```

4. **Template not found**

   ```bash
   # List available templates
   generator list-templates --output-format json
   ```

### Debug Mode

Enable debug output for troubleshooting:

```bash
export GENERATOR_LOG_LEVEL=debug
export GENERATOR_VERBOSE=true
generator generate --non-interactive --debug
```

### Validation Issues

Check validation details:

```bash
generator validate ./project --verbose --show-fixes --output-format json
```

This comprehensive guide should help you integrate the generator into your automation workflows effectively.
