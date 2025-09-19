# Interactive Directory Selection

This document describes the interactive directory selection feature implemented for the CLI generator.

## Overview

The interactive directory selection feature allows users to interactively choose and validate output directories for project generation. It provides a user-friendly interface with conflict resolution, validation, and backup capabilities.

## Features

### 1. Interactive Directory Path Input

- **Default Path**: Suggests `output/generated` as the default output directory
- **Custom Path Input**: Allows users to enter custom paths with real-time validation
- **Path Normalization**: Automatically converts relative paths to absolute paths
- **Home Directory Expansion**: Supports `~/` notation for home directory paths

### 2. Path Validation

- **Format Validation**: Ensures paths use valid characters and format
- **Length Limits**: Enforces reasonable path length limits (max 500 characters)
- **Reserved Names**: Prevents use of system reserved names (CON, PRN, etc.)
- **Parent Directory Check**: Validates that parent directories exist and are writable

### 3. Directory Conflict Resolution

When the target directory already exists and contains files, the system provides multiple resolution options:

#### **Overwrite with Backup**

- Creates a complete backup of existing files before overwriting
- Backup location: `<original-name>.backup.<timestamp>`
- Requires explicit user confirmation
- Safest option for replacing existing projects

#### **Merge Files**

- Keeps existing files and adds new ones alongside
- Overwrites files with the same names (no backup for individual files)
- May cause conflicts with existing project structure
- Requires user confirmation

#### **Choose Different Directory**

- Returns to directory selection to pick a different location
- Allows users to avoid conflicts entirely
- Recommended approach for new projects

#### **Cancel Generation**

- Stops the generation process entirely
- No changes are made to the file system

### 4. Safety Features

- **Backup Creation**: Automatic backup before destructive operations
- **Confirmation Dialogs**: Explicit confirmation required for risky operations
- **Directory Contents Display**: Shows existing files to help users make informed decisions
- **Error Recovery**: Graceful error handling with helpful suggestions

## Usage

### Interactive Mode

When running the generator in interactive mode, directory selection is automatically included in the workflow:

```bash
generator generate
```

The system will:

1. Collect project configuration
2. **Prompt for output directory** (this feature)
3. Show project structure preview
4. Generate the project

### Directory Selection Flow

1. **Path Input**: User enters desired output directory path
2. **Validation**: System validates path format and accessibility
3. **Conflict Check**: System checks if directory exists and contains files
4. **Resolution**: If conflicts exist, user chooses resolution strategy
5. **Confirmation**: User confirms destructive operations
6. **Preparation**: System creates directories and backups as needed

## Implementation Details

### Core Components

#### `DirectorySelector`

- Main class handling directory selection workflow
- Located in `pkg/ui/directory_selector.go`
- Integrates with interactive UI framework

#### `DirectoryValidator`

- Handles path validation and safety checks
- Validates format, length, reserved names, and permissions
- Provides detailed error messages with suggestions

#### `DirectorySelectionResult`

- Contains the result of directory selection process
- Includes selected path, conflict resolution, and backup information

### Integration Points

#### CLI Integration

- Integrated into `pkg/cli/cli.go` interactive generation workflow
- Called between project configuration and generation steps
- Uses `runInteractiveDirectorySelection()` method

#### UI Framework

- Uses existing interactive UI components (menus, prompts, confirmations)
- Consistent with other interactive features
- Supports navigation (back, quit, help)

## Error Handling

### Validation Errors

- **Invalid Characters**: Clear error messages with valid character suggestions
- **Path Too Long**: Suggests shorter alternatives
- **Reserved Names**: Explains why names are reserved and suggests alternatives
- **Permission Issues**: Provides guidance on fixing permission problems

### Recovery Options

- **Automatic Suggestions**: System provides helpful suggestions for common issues
- **Retry Capability**: Users can retry after fixing issues
- **Graceful Degradation**: Falls back to safe defaults when possible

## Testing

### Unit Tests

- **Path Validation**: Tests for all validation scenarios
- **Directory Creation**: Verifies directory creation functionality
- **Backup Operations**: Tests backup creation and restoration
- **Conflict Resolution**: Tests all conflict resolution paths

### Integration Tests

- **CLI Integration**: Tests integration with CLI generate command
- **End-to-End Workflows**: Tests complete directory selection workflows
- **Error Scenarios**: Tests error handling and recovery

### Test Coverage

- Located in `pkg/ui/directory_selector_unit_test.go`
- Covers all major functionality and edge cases
- Uses temporary directories for safe testing

## Configuration

### Default Settings

- **Default Path**: `output/generated`
- **Max Path Length**: 500 characters
- **Backup Naming**: `<original>.backup.<pid>`

### Customization

The directory selector can be customized through:

- Default path configuration
- Validation rules
- UI behavior settings

## Security Considerations

### Path Security

- **Directory Traversal Prevention**: Validates against `../` attacks
- **Permission Validation**: Ensures write permissions before proceeding
- **Safe Operations**: Uses atomic operations where possible

### Backup Security

- **Complete Backups**: Preserves all file permissions and metadata
- **Unique Names**: Uses process ID to ensure unique backup names
- **Cleanup**: Provides mechanisms for backup cleanup

## Future Enhancements

### Planned Features

- **Template-Aware Validation**: Validate paths based on selected templates
- **Batch Operations**: Support for generating multiple projects
- **Advanced Backup Options**: Configurable backup strategies
- **Integration with Version Control**: Git integration for backup management

### Performance Optimizations

- **Async Operations**: Background backup creation
- **Progress Tracking**: Detailed progress for large operations
- **Caching**: Cache validation results for repeated operations

## Troubleshooting

### Common Issues

#### "Permission Denied" Errors

- **Cause**: Insufficient write permissions on target directory
- **Solution**: Choose a directory where you have write access, or run with appropriate permissions

#### "Directory Not Empty" Warnings

- **Cause**: Target directory contains existing files
- **Solution**: Choose a conflict resolution option (overwrite, merge, or different directory)

#### "Invalid Path" Errors

- **Cause**: Path contains invalid characters or exceeds length limits
- **Solution**: Use only valid characters (letters, numbers, hyphens, underscores, slashes)

### Debug Information

Enable verbose mode for detailed directory selection information:

```bash
generator generate --verbose
```

This will show:

- Path validation steps
- Directory existence checks
- Backup creation progress
- Conflict resolution choices
