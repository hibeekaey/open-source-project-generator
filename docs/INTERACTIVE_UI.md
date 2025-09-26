# Interactive UI Framework Integration

The Open Source Project Generator now includes a comprehensive interactive UI framework that provides a user-friendly, guided experience for project configuration and generation.

## Overview

The interactive UI framework has been fully integrated with the `generate` command to provide:

- **Guided Project Configuration**: Step-by-step collection of project details with validation
- **Component Selection**: Multi-select interface for choosing project components
- **Real-time Validation**: Input validation with helpful error messages and recovery options
- **Progress Tracking**: Visual progress indicators during project generation
- **Context-sensitive Help**: Comprehensive help system available throughout the process

## Features

### 🎯 Interactive Project Configuration

The interactive mode guides users through project setup with intelligent prompts:

```bash
generator generate --interactive
```

**Configuration Steps:**

1. **Basic Details**: Project name, organization, description, author
2. **License Selection**: Choose from popular open source licenses
3. **Component Selection**: Select project components with search and filtering
4. **Confirmation**: Review configuration before generation
5. **Generation**: Real-time progress tracking with detailed logs

### ⌨️ Keyboard Navigation

Consistent keyboard shortcuts work across all interactive components:

- **Navigation**: `↑/↓` or `j/k` to move up/down
- **Selection**: `Enter` to select, `Space` to toggle
- **Search**: `/` to start searching (where available)
- **Actions**: `h` for help, `b` to go back, `q` to quit
- **Numbers**: Type `1-9` to jump to specific options

### 🔍 Component Selection

Advanced multi-select interface for choosing project components:

- **Available Components**:
  - Go Gin API (RESTful backend)
  - Next.js Frontend (React-based)
  - PostgreSQL Database
  - Redis Cache
  - Docker Configuration
  - Kubernetes Manifests
  - CI/CD Pipeline
  - Monitoring Setup

- **Features**:
  - Search and filter components by name or tags
  - Minimum/maximum selection validation
  - Category-based organization
  - Detailed descriptions for each component

### ✅ Input Validation

Comprehensive validation with user-friendly error handling:

- **Project Name**: Validates format, length, and character restrictions
- **Email**: RFC-compliant email validation (optional fields)
- **URLs**: HTTP/HTTPS URL validation with protocol checking
- **Versions**: Semantic version validation

**Validation Features:**

- Real-time validation feedback
- Helpful error messages with suggestions
- Recovery options for common mistakes
- Input sanitization and correction suggestions

### 📊 Progress Tracking

Real-time progress tracking during project generation:

- **Visual Progress Bar**: Shows completion percentage
- **Step Indicators**: Current step with description
- **ETA Calculation**: Estimated time to completion
- **Activity Log**: Recent operations and status messages
- **Cancellation Support**: Ability to cancel long-running operations

### 🆘 Error Handling

Comprehensive error handling with recovery options:

- **Validation Errors**: Clear messages with correction suggestions
- **Recovery Actions**: Automated fixes for common issues
- **Safety Indicators**: Safe vs. cautious recovery options
- **Context Help**: Detailed help for error resolution

## Usage Examples

### Basic Interactive Generation

```bash
# Start interactive project generation
generator generate --interactive

# Interactive with specific output directory
generator generate --interactive --output ./my-project
```

### Advanced Options

```bash
# Interactive with template pre-selection
generator generate --interactive --template go-gin

# Interactive with dry-run mode
generator generate --interactive --dry-run

# Interactive with verbose output
generator generate --interactive --verbose
```

### Non-Interactive Mode

The generator also supports non-interactive mode for automation:

```bash
# Generate from configuration file
generator generate --config project.yaml

# Non-interactive with environment variables
GENERATOR_PROJECT_NAME=myapp generator generate --non-interactive
```

## Configuration

The interactive UI can be customized through configuration:

```yaml
# .kiro/config.yaml
ui:
  enable_colors: true
  enable_unicode: true
  page_size: 10
  timeout: 30m
  auto_save: true
  show_breadcrumbs: true
  show_shortcuts: true
  confirm_on_quit: true
```

## Architecture

### Components

The interactive UI framework consists of several key components:

1. **Interactive UI Interface** (`pkg/interfaces/interactive.go`)
   - Defines contracts for all interactive components
   - Session management and state persistence
   - Progress tracking and error handling

2. **Core Implementation** (`pkg/ui/interactive.go`)
   - Menu navigation and selection
   - Text input with validation
   - Multi-select and checkbox interfaces

3. **Display Components** (`pkg/ui/display.go`)
   - Table display with pagination and sorting
   - Tree structure visualization
   - Progress tracking implementation

4. **Validation Framework** (`pkg/ui/validation.go`)
   - Input validation with chaining
   - Common validators (email, URL, version, etc.)
   - Error recovery mechanisms

5. **Error Handling** (`pkg/ui/error_handling.go`)
   - User-friendly error display
   - Recovery options with safety indicators
   - Context-sensitive help system

### Integration Points

The interactive UI is integrated with the CLI through:

- **CLI Struct**: Includes `interactiveUI` field for UI operations
- **Generate Command**: Comprehensive with interactive mode support
- **Session Management**: Automatic session lifecycle management
- **Progress Tracking**: Real-time feedback during generation

## Testing

Comprehensive test coverage includes:

- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end workflow testing
- **Mock Framework**: Complete mock implementations for testing
- **Benchmarks**: Performance testing for UI operations

Run tests with:

```bash
# Run all UI tests
go test ./pkg/ui/...

# Run CLI integration tests
go test ./pkg/cli/... -run TestInteractive

# Run with verbose output
go test -v ./pkg/ui/... ./pkg/cli/...
```

## Best Practices

### For Users

1. **Use Interactive Mode**: Recommended for first-time users and complex configurations
2. **Leverage Search**: Use `/` to search in multi-select components
3. **Read Help**: Press `h` for context-sensitive help
4. **Review Configuration**: Always review the summary before generation

### For Developers

1. **Validation First**: Always validate input before processing
2. **Error Recovery**: Provide helpful recovery options for errors
3. **Progress Feedback**: Show progress for long-running operations
4. **Consistent Navigation**: Use standard keyboard shortcuts

## Troubleshooting

### Common Issues

**Issue**: Interactive mode not working in CI/CD
**Solution**: Use `--non-interactive` flag or configuration files

**Issue**: Colors not displaying correctly
**Solution**: Set `TERM` environment variable or disable colors in config

**Issue**: Keyboard shortcuts not working
**Solution**: Ensure terminal supports the required key combinations

### Debug Mode

Enable debug mode for detailed logging:

```bash
generator generate --interactive --debug
```

This provides:

- Detailed operation logs
- Performance metrics
- UI state information
- Error stack traces

## Future Enhancements

Planned improvements include:

- **Theme Support**: Customizable color schemes and themes
- **Plugin System**: Extensible UI components
- **Advanced Search**: Fuzzy search and filtering
- **Accessibility**: Screen reader support and keyboard-only navigation
- **Internationalization**: Multi-language support

## Contributing

To contribute to the interactive UI framework:

1. **Follow Patterns**: Use existing UI patterns and interfaces
2. **Add Tests**: Include comprehensive test coverage
3. **Document Changes**: Update documentation for new features
4. **Validate Input**: Always include proper input validation
5. **Handle Errors**: Provide user-friendly error messages

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.
