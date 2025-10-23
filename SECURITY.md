# Security

This document outlines the security measures implemented in the Open Source Project Generator (v1.3.0+).

## Security Fixes Applied

### 1. Path Traversal Prevention (G304 - CWE-22)

**Issue**: File operations using user-controlled paths could allow directory traversal attacks.

**Fix**:

- Created `pkg/security/` package with path validation and sanitization functions
- Implemented `SanitizePath()` to detect and prevent path traversal attempts
- Added secure file operation wrappers with automatic path validation
- Updated config manager and filesystem operations to use security package
- All security errors use categorized error types from `pkg/errors/`

**Functions**:

```go
// Sanitizes and validates paths, prevents directory traversal
func SanitizePath(path string) (string, error)

// Validates paths against allowed base paths
func ValidatePath(path string, allowedBasePaths ...string) error

// Secure file operations with automatic path validation
func SafeReadFile(path string, allowedBasePaths ...string) ([]byte, error)
func SafeWriteFile(path string, data []byte, allowedBasePaths ...string) error
func SafeMkdirAll(path string, allowedBasePaths ...string) error
```

### 2. File Permission Hardening (G301/G302/G306 - CWE-276)

**Issue**: Files and directories were created with overly permissive permissions.

**Changes**:

- Directory permissions: `0755` → `0750` (removed world read/execute)
- File permissions: `0644` → `0600` (removed group/world read)
- Log file permissions: `0644` → `0600`

**Files Updated**:

- `pkg/security/` - Security operations package
- `internal/config/manager.go` - Config file operations
- `internal/app/logger.go` - Log file creation
- `internal/app/app.go` - Application initialization
- `pkg/filesystem/` - File system operations
- `pkg/cli/` - CLI command handlers

### 3. Error Handling (G104 - CWE-703)

**Issue**: Some function calls that could return errors were not being handled.

**Fix**:

- All security operations return categorized errors from `pkg/errors/`
- Security errors use `errors.NewSecurityError()` for consistent handling
- Validation errors use `errors.NewValidationError()` for input issues
- Config errors use `errors.NewConfigError()` for configuration problems

**Status**: Core generator code uses proper error handling. Templates updated to follow same patterns.

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

1. Always sanitize file paths using `pkg/security` before operations
2. Use restrictive file permissions (0600 for files, 0750 for directories)
3. Validate user input that affects file paths
4. Use absolute paths when possible
5. Return categorized errors from `pkg/errors/` for proper error handling

### Path Sanitization

```go
import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/security"
    "github.com/cuesoftinc/open-source-project-generator/pkg/errors"
)

// Always sanitize user input before file operations
sanitized, err := security.SanitizePath(userPath)
if err != nil {
    return errors.NewSecurityError("invalid path", err)
}
```

### Path Validation

```go
// Validate against allowed base paths
if err := security.ValidatePath(sanitized, "/allowed/base/path"); err != nil {
    return errors.NewSecurityError("path traversal detected", err)
}
```

### Secure File Operations

```go
// Use security package functions instead of direct os calls
content, err := security.SafeReadFile(path, allowedBasePaths...)
if err != nil {
    return errors.NewSecurityError("failed to read file", err)
}

err = security.SafeWriteFile(path, data, allowedBasePaths...)
if err != nil {
    return errors.NewSecurityError("failed to write file", err)
}

err = security.SafeMkdirAll(path, allowedBasePaths...)
if err != nil {
    return errors.NewSecurityError("failed to create directory", err)
}
```

### Architecture Integration

Security operations follow the project's dependency injection pattern:

```go
// Components receive security dependencies via interfaces
type FileSystemManager struct {
    security security.Interface
    validator validation.Interface
}

func NewFileSystemManager(sec security.Interface, val validation.Interface) *FileSystemManager {
    return &FileSystemManager{
        security: sec,
        validator: val,
    }
}
```

## Remaining Security Considerations

### Tool Execution Security

- Only whitelisted tools and flags are executed
- All tool commands are validated before execution
- Command injection prevention through parameter validation
- Tool output is sanitized before processing
- Timeout protection prevents hanging processes

### Input Validation

- All user inputs validated through `internal/config/` validators
- Path inputs sanitized via `pkg/security/SanitizePath()`
- Configuration fields validated before tool execution
- Categorized error handling via `pkg/cli` error types
- Early validation with fail-fast approach

### Dependency Security

- Regularly update dependencies using `go get -u`
- Use dependency scanning tools (nancy, govulncheck)
- Pin dependency versions in `go.mod`
- Tool version requirements defined in `internal/orchestrator/tool_discovery.go`
- Offline mode supported via cached tool metadata

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

# Scan for secrets in code
trufflehog filesystem ./

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

- [ ] Sanitize all file paths using `pkg/security/SanitizePath()`
- [ ] Use secure file permissions (0600 for files, 0750 for directories)
- [ ] Return categorized errors from `pkg/cli` error types
- [ ] Validate user inputs through `internal/config/` validators
- [ ] Use dependency injection with interfaces from `pkg/interfaces/`
- [ ] Review tool execution code for command injection vulnerabilities
- [ ] Run security scanners before submitting PRs (`make security-scan`)
- [ ] Test generated projects for security vulnerabilities
- [ ] Update security documentation for new features
- [ ] Follow clean architecture patterns (presentation → business logic → infrastructure)
- [ ] Never use direct `os` package calls; use `pkg/filesystem/` and `pkg/security/` abstractions
