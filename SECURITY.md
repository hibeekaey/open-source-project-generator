# Security

This document outlines the security measures implemented in the Open Source Project Generator (v1.3.0+).

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

### 4. Template Security Enhancements

**Issue**: Generated projects needed enhanced security measures.

**Fixes**:

- **JWT Security**: Implemented secure JWT token validation with algorithm verification
- **Database Security**: Added parameterized queries and input validation
- **Authentication Security**: Enhanced authentication middleware with rate limiting
- **File Permission Hardening**: Applied secure file permissions in generated projects
- **Input Validation**: Comprehensive input sanitization and validation
- **Error Handling**: Secure error responses that don't leak sensitive information

**Security Features in Generated Projects**:

- SQL injection prevention through parameterized queries
- XSS protection through input sanitization
- CSRF protection in web applications
- Secure password hashing with bcrypt
- Rate limiting on authentication endpoints
- Secure session management
- Comprehensive logging and monitoring

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

# Run govulncheck for Go vulnerabilities
govulncheck ./...

# Scan for secrets in templates
trufflehog filesystem ./pkg/template/templates/

# Validate generated projects
generator audit ./output --security --fail-on-high
```

## Reporting Security Issues

If you discover a security vulnerability, please report it privately to the maintainers rather than opening a public issue.

## Security Features in Generated Projects

The generator creates projects with built-in security features:

### Backend Security (Go Gin)

- **JWT Authentication**: Secure token validation with algorithm verification
- **Database Security**: Parameterized queries preventing SQL injection
- **Input Validation**: Comprehensive input sanitization and validation
- **Rate Limiting**: Protection against brute force attacks
- **Password Security**: bcrypt hashing with configurable cost factors
- **CORS Protection**: Configurable cross-origin resource sharing policies
- **Security Headers**: Automatic security header implementation

### Frontend Security (React/Next.js)

- **XSS Protection**: Input sanitization and CSP headers
- **CSRF Protection**: Token-based CSRF protection
- **Secure Authentication**: JWT token management with secure storage
- **Content Security Policy**: Comprehensive CSP implementation
- **Dependency Security**: Regular security updates and vulnerability scanning

### Mobile Security (Android/iOS)

- **Secure Storage**: Encrypted local storage for sensitive data
- **Certificate Pinning**: SSL certificate validation
- **Biometric Authentication**: Secure biometric authentication support
- **Network Security**: HTTPS enforcement and secure communication

### Infrastructure Security

- **Container Security**: Multi-stage builds with minimal attack surface
- **Kubernetes Security**: Security contexts and network policies
- **Terraform Security**: Infrastructure as code with security best practices
- **Monitoring**: Comprehensive security monitoring and alerting

## Security Checklist for Contributors

- [ ] Validate all file paths
- [ ] Use secure file permissions
- [ ] Handle all errors appropriately
- [ ] Sanitize user inputs
- [ ] Review templates for security issues
- [ ] Run security scanners before submitting PRs
- [ ] Test generated projects for security vulnerabilities
- [ ] Update security documentation for new features
