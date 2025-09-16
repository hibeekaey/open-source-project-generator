# Security

This document outlines the security measures implemented in the Open Source Template Generator.

## Security Fixes Applied

### 1. Path Traversal Prevention (G304 - CWE-22)

**Issue**: File operations using user-controlled paths could allow directory traversal attacks.

**Fix**:

- Created `pkg/utils/security.go` with path validation functions
- Implemented `ValidatePath()` to detect and prevent path traversal attempts
- Added `SafeReadFile()`, `SafeWriteFile()`, `SafeMkdirAll()`, and `SafeOpenFile()` functions
- Updated config manager to use secure file operations

**Functions**:

```go
// Validates paths and prevents directory traversal
func ValidatePath(path string, allowedBasePaths ...string) error

// Safe file operations with path validation
func SafeReadFile(path string, allowedBasePaths ...string) ([]byte, error)
func SafeWriteFile(path string, data []byte, allowedBasePaths ...string) error
func SafeMkdirAll(path string, allowedBasePaths ...string) error
func SafeOpenFile(path string, flag int, perm os.FileMode, allowedBasePaths ...string) (*os.File, error)
```

### 2. File Permission Hardening (G301/G302/G306 - CWE-276)

**Issue**: Files and directories were created with overly permissive permissions.

**Changes**:

- Directory permissions: `0755` → `0750` (removed world read/execute)
- File permissions: `0644` → `0600` (removed group/world read)
- Log file permissions: `0644` → `0600`

**Files Updated**:

- `internal/config/manager.go`
- `internal/app/logger.go`
- `internal/app/app.go`
- `pkg/filesystem/generator.go`
- `pkg/cli/cli.go`

### 3. Error Handling (G104 - CWE-703)

**Issue**: Some function calls that could return errors were not being handled.

**Status**: Identified in generated code (not core generator code). These should be addressed in templates.

## Security Best Practices

### File Operations

1. Always validate file paths before operations
2. Use restrictive file permissions (0600 for files, 0750 for directories)
3. Sanitize user input that affects file paths
4. Use absolute paths when possible

### Path Validation

```go
// Example usage
if err := utils.ValidatePath(userPath, "/allowed/base/path"); err != nil {
    return fmt.Errorf("invalid path: %w", err)
}
```

### Secure File Creation

```go
// Use secure utilities instead of direct os calls
content, err := utils.SafeReadFile(path, allowedBasePaths...)
err = utils.SafeWriteFile(path, data, allowedBasePaths...)
err = utils.SafeMkdirAll(path, allowedBasePaths...)
```

## Remaining Security Considerations

### Template Security

- Templates should be reviewed for security issues
- Generated code should follow security best practices
- Consider implementing template sandboxing

### Input Validation

- Validate all user inputs (project names, paths, configurations)
- Sanitize template variables
- Implement proper error handling

### Dependency Security

- Regularly update dependencies
- Use dependency scanning tools
- Pin dependency versions

## Security Testing

Run security scans regularly:

```bash
# Run gosec security scanner
gosec ./...

# Run dependency vulnerability scanner
go list -json -deps ./... | nancy sleuth

# Run static analysis
staticcheck ./...
```

## Reporting Security Issues

If you discover a security vulnerability, please report it privately to the maintainers rather than opening a public issue.

## Security Checklist for Contributors

- [ ] Validate all file paths
- [ ] Use secure file permissions
- [ ] Handle all errors appropriately
- [ ] Sanitize user inputs
- [ ] Review templates for security issues
- [ ] Run security scanners before submitting PRs
