# Security Fixer

The Security Fixer is an automated tool that applies security fixes to template files in the codebase. It builds upon the Security Scanner to not only detect security issues but also automatically fix them.

## Features

- **Automated CORS Fixes**: Removes dangerous `Access-Control-Allow-Origin: null` headers and replaces wildcard origins with proper validation
- **Security Headers**: Automatically adds comprehensive security headers (X-Content-Type-Options, X-Frame-Options, X-XSS-Protection)
- **JWT Security**: Fixes JWT 'none' algorithm vulnerabilities and adds token expiration recommendations
- **SQL Injection Prevention**: Identifies and provides fixes for SQL injection vulnerabilities
- **Information Leakage Prevention**: Fixes detailed error messages and debug information exposure
- **Cookie Security**: Adds secure flags to cookie configurations

## Usage

### Basic Usage

```bash
# Fix all security issues in the templates directory
./security-fixer -dir templates

# Run in dry-run mode to see what would be fixed without making changes
./security-fixer -dir templates -dry-run

# Fix only specific types of issues
./security-fixer -dir templates -fix-type cors
./security-fixer -dir templates -fix-type headers
./security-fixer -dir templates -fix-type auth
./security-fixer -dir templates -fix-type sql

# Enable verbose output
./security-fixer -dir templates -verbose

# Disable backup creation
./security-fixer -dir templates -backup=false
```

### Command Line Options

- `-dir`: Directory containing template files to fix (default: "templates")
- `-dry-run`: Show what would be fixed without making changes (default: false)
- `-verbose`: Enable verbose output (default: false)
- `-fix-type`: Type of fixes to apply: all, cors, headers, auth, sql (default: "all")
- `-backup`: Create backup files before applying fixes (default: true)

## Supported Fix Types

### CORS Fixes (`-fix-type cors`)

1. **Null Origin Fix**: Removes `Access-Control-Allow-Origin: null` headers
   - **Before**: `c.Header("Access-Control-Allow-Origin", "null")`
   - **After**: Comments out the line with security explanation

2. **Wildcard Origin Fix**: Replaces wildcard CORS with proper validation
   - **Before**: `c.Header("Access-Control-Allow-Origin", "*")`
   - **After**: Adds origin validation logic

### Security Headers (`-fix-type headers`)

1. **Content-Type Options**: Adds X-Content-Type-Options header
   - **Added**: `c.Header("X-Content-Type-Options", "nosniff")`

2. **Comprehensive Headers**: Adds multiple security headers
   - **Added**: X-Frame-Options, X-XSS-Protection, X-Content-Type-Options

### Authentication Fixes (`-fix-type auth`)

1. **JWT None Algorithm**: Replaces 'none' algorithm with secure alternative
   - **Before**: `algorithm: "none"`
   - **After**: `algorithm: "HS256"`

2. **JWT Expiration**: Adds token expiration recommendations
   - **Added**: Comments about setting appropriate expiration times

3. **Cookie Security**: Adds secure cookie flags
   - **Added**: Comments about HttpOnly and Secure flags

### SQL Injection Fixes (`-fix-type sql`)

1. **String Concatenation**: Identifies SQL queries using string concatenation
   - **Added**: Comments recommending parameterized queries

2. **Variable Interpolation**: Identifies SQL queries with variable interpolation
   - **Added**: Comments recommending placeholder usage

## Examples

### Fix CORS Issues Only

```bash
./security-fixer -dir templates -fix-type cors -verbose
```

### Dry Run with Full Report

```bash
./security-fixer -dir templates -dry-run -verbose
```

### Fix All Issues with Backups

```bash
./security-fixer -dir templates -backup -verbose
```

## Output

The tool provides detailed output showing:

- Files processed
- Issues found and fixed
- Backup files created
- Summary statistics

### Example Output

```
Security Fix Report
==================

Fixed Issues (5):
------------------
✓ templates/backend/cors.go.tmpl (Line 14): Fixed CORS null origin vulnerability
  Fix: Removed Access-Control-Allow-Origin header for null origins
✓ templates/backend/auth.go.tmpl (Line 25): Fixed JWT none algorithm vulnerability
  Fix: Replaced 'none' algorithm with secure HS256
✓ templates/backend/handlers.go.tmpl (Line 45): Added comprehensive security headers
  Fix: Added X-Frame-Options, X-XSS-Protection, and X-Content-Type-Options headers

Summary: 5 issues fixed, 0 errors
Backup files created: 3
```

## Integration

The Security Fixer can be integrated into CI/CD pipelines to automatically fix security issues:

```yaml
# GitHub Actions example
- name: Fix Security Issues
  run: |
    go build -o security-fixer ./cmd/security-fixer
    ./security-fixer -dir templates -verbose
```

## Safety Features

- **Backup Creation**: Automatically creates timestamped backup files before applying fixes
- **Dry Run Mode**: Preview changes without modifying files
- **Selective Fixing**: Apply only specific types of fixes
- **Template Preservation**: Maintains template functionality and variables
- **Error Handling**: Continues processing other files if one fails

## Building

```bash
go build -o security-fixer ./cmd/security-fixer
```

## Testing

The Security Fixer includes comprehensive tests:

```bash
go test ./pkg/security -v -run TestFixer
```

## Related Tools

- **Security Scanner** (`cmd/security-scanner`): Detects security issues without fixing them
- **Import Analyzer** (`cmd/import-analyzer`): Analyzes import dependencies
- **Standards Tool** (`cmd/standards`): Validates coding standards