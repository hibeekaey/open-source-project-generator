# Cleanup Infrastructure

This package provides comprehensive cleanup infrastructure and analysis tools for the Open Source Template Generator project. It implements a systematic approach to code quality analysis, backup management, and validation to ensure safe and effective cleanup operations.

## Overview

The cleanup infrastructure consists of four main components:

1. **Cleanup Manager** - Orchestrates the entire cleanup process
2. **Code Analyzer** - Analyzes Go source code for various issues
3. **Backup Manager** - Creates and manages file backups
4. **Validation Framework** - Ensures project integrity throughout the process

## Components

### Cleanup Manager (`manager.go`)

The central coordinator that manages the cleanup process:

```go
manager, err := cleanup.NewManager(projectRoot, config)
if err != nil {
    log.Fatal(err)
}
defer manager.Shutdown()

// Initialize infrastructure
if err := manager.Initialize(); err != nil {
    log.Fatal(err)
}

// Analyze project
analysis, err := manager.AnalyzeProject()
if err != nil {
    log.Fatal(err)
}
```

**Key Features:**
- Project integrity validation
- Comprehensive analysis orchestration
- Backup creation and management
- Validation checkpoints
- Graceful error handling and recovery

### Code Analyzer (`analyzer.go`)

Provides utilities for analyzing Go source code:

```go
analyzer := cleanup.NewCodeAnalyzer()

// Find TODO/FIXME comments
todos, err := analyzer.AnalyzeTODOComments(rootDir)

// Identify duplicate code
duplicates, err := analyzer.FindDuplicateCode(rootDir)

// Find unused code
unused, err := analyzer.IdentifyUnusedCode(rootDir)

// Validate import organization
issues, err := analyzer.ValidateImportOrganization(rootDir)
```

**Analysis Capabilities:**
- TODO/FIXME/HACK comment detection with priority and category classification
- Duplicate code block identification
- Unused function, variable, and import detection
- Import organization validation
- AST-based code analysis

### Backup Manager (`backup.go`)

Handles safe backup and restoration of files:

```go
backupMgr := cleanup.NewBackupManager(backupDir)

// Create backup
backup, err := backupMgr.CreateBackup(files)

// Restore if needed
err = backupMgr.RestoreBackup(backup)

// Cleanup old backups
err = backupMgr.CleanupOldBackups(maxAge)
```

**Backup Features:**
- Incremental file backups with checksums
- Directory structure preservation
- Automatic cleanup of old backups
- Integrity verification
- Rollback capabilities

### Validation Framework (`validator.go`)

Ensures no functionality is lost during cleanup:

```go
validator := cleanup.NewValidationFramework(projectRoot)

// Check project integrity
err := validator.EnsureProjectIntegrity()

// Create validation checkpoint
checkpoint, err := validator.CreateValidationCheckpoint()

// Validate after changes
result, err := validator.ValidateAfterChanges(changedFiles)

// Compare results
differences := validator.CompareValidationResults(before, after)
```

**Validation Features:**
- Project structure integrity checks
- Build validation
- Test suite execution
- Go module consistency validation
- Code quality checks (go vet, gofmt)
- Validation result comparison

## Configuration

The cleanup infrastructure is highly configurable:

```go
config := &cleanup.Config{
    BackupDir:       ".cleanup-backups",
    DryRun:          false,
    Verbose:         true,
    SkipPatterns:    []string{"vendor/", ".git/"},
    PreserveTODOs:   []string{"IMPORTANT:", "CRITICAL:"},
    MaxBackupAge:    7 * 24 * time.Hour,
    ValidationLevel: cleanup.ValidationStandard,
}
```

### Configuration Options

- **BackupDir**: Directory for storing backups
- **DryRun**: Perform analysis without making changes
- **Verbose**: Enable detailed logging
- **SkipPatterns**: File/directory patterns to skip during analysis
- **PreserveTODOs**: TODO patterns that should not be automatically resolved
- **MaxBackupAge**: Maximum age for backup retention
- **ValidationLevel**: Validation strictness (Basic, Standard, Strict)

## CLI Usage

The cleanup infrastructure includes a CLI tool:

```bash
# Initialize cleanup infrastructure
./cleanup-tool -action=init -verbose

# Analyze project for issues
./cleanup-tool -action=analyze -verbose

# Validate project integrity
./cleanup-tool -action=validate -verbose

# Perform dry run analysis
./cleanup-tool -action=analyze -dry-run -verbose
```

### CLI Options

- `-root`: Project root directory (default: current directory)
- `-action`: Action to perform (init, analyze, validate)
- `-dry-run`: Perform analysis without making changes
- `-verbose`: Enable verbose output
- `-backup-dir`: Custom backup directory

## Analysis Results

### TODO Analysis

The analyzer categorizes TODO comments by:

**Types**: TODO, FIXME, HACK, XXX, BUG, NOTE

**Priorities**:
- **Critical**: Security-related issues
- **High**: FIXME, BUG, security mentions
- **Medium**: HACK, performance mentions
- **Low**: General TODOs and notes

**Categories**:
- Security
- Performance
- Feature
- Bug
- Documentation
- Refactor

### Code Quality Issues

The analyzer identifies:

- **Duplicate Code**: Similar function implementations across files
- **Unused Code**: Unused functions, variables, types, and imports
- **Import Issues**: Incorrect import organization and unused imports
- **Naming Inconsistencies**: Violations of Go naming conventions

## Safety Features

### Backup System

- **Automatic Backups**: Created before any modifications
- **Checksum Verification**: Ensures backup integrity
- **Incremental Storage**: Efficient storage using file differences
- **Easy Restoration**: Simple rollback to previous state

### Validation Framework

- **Pre-modification Validation**: Ensures project is in good state
- **Post-modification Validation**: Verifies no functionality lost
- **Continuous Monitoring**: Tracks changes throughout process
- **Rollback Triggers**: Automatic rollback on validation failures

### Error Handling

- **Graceful Degradation**: Continues operation when possible
- **Detailed Error Reporting**: Comprehensive error context
- **Recovery Mechanisms**: Automatic recovery from common issues
- **User Guidance**: Actionable suggestions for error resolution

## Integration Testing

The package includes comprehensive integration tests:

```bash
# Run all cleanup tests
go test ./internal/cleanup/... -v

# Run integration tests only
go test ./internal/cleanup/ -run Integration -v

# Run specific component tests
go test ./internal/cleanup/ -run TestCodeAnalyzer -v
```

## Best Practices

### Before Using

1. **Commit Changes**: Ensure all changes are committed to version control
2. **Run Tests**: Verify all tests pass before cleanup
3. **Review Configuration**: Adjust settings for your project needs
4. **Start with Dry Run**: Always test with `-dry-run` first

### During Cleanup

1. **Monitor Progress**: Use verbose mode to track operations
2. **Validate Frequently**: Check validation results at each step
3. **Backup Critical Files**: Ensure important files are backed up
4. **Test Incrementally**: Validate after each major change

### After Cleanup

1. **Run Full Test Suite**: Ensure all functionality preserved
2. **Review Changes**: Manually review all modifications
3. **Update Documentation**: Reflect any structural changes
4. **Clean Up Backups**: Remove old backups when satisfied

## Error Recovery

If issues occur during cleanup:

1. **Check Logs**: Review verbose output for error details
2. **Restore Backup**: Use backup manager to restore files
3. **Validate State**: Run validation to check project integrity
4. **Report Issues**: Document any bugs or unexpected behavior

## Performance Considerations

- **Large Projects**: Use skip patterns to exclude large directories
- **Memory Usage**: Monitor memory usage on very large codebases
- **Parallel Processing**: Some operations can be parallelized
- **Incremental Analysis**: Focus on specific directories when possible

## Security Considerations

- **Backup Security**: Backups contain full file contents
- **Temporary Files**: Cleanup temporary files after operations
- **File Permissions**: Preserve original file permissions
- **Sensitive Data**: Be careful with files containing secrets

## Future Enhancements

Planned improvements include:

- **Parallel Processing**: Concurrent analysis for better performance
- **Advanced Duplicate Detection**: More sophisticated similarity algorithms
- **Custom Rules**: User-defined analysis and cleanup rules
- **IDE Integration**: Direct integration with development environments
- **Metrics Dashboard**: Web interface for cleanup metrics and progress

## Contributing

When contributing to the cleanup infrastructure:

1. **Add Tests**: Include comprehensive tests for new features
2. **Update Documentation**: Keep this README current
3. **Follow Patterns**: Maintain consistency with existing code
4. **Consider Safety**: Ensure new features don't compromise safety
5. **Performance Impact**: Consider impact on large projects

## Troubleshooting

### Common Issues

**"Essential file missing"**: Ensure project has required Go project structure
**"Build failed"**: Fix compilation errors before running cleanup
**"Tests failed"**: Resolve test failures before proceeding
**"Permission denied"**: Check file permissions and backup directory access

### Debug Mode

Enable debug logging for detailed troubleshooting:

```go
config.Verbose = true
config.ValidationLevel = ValidationStrict
```

### Support

For issues or questions:

1. Check the integration tests for usage examples
2. Review the CLI tool implementation
3. Examine the test files for expected behavior
4. Create detailed issue reports with logs and context