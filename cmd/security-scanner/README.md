# Security Scanner

A Go-based security scanner utility that analyzes template files for security vulnerabilities and best practice violations.

## Features

- **CORS Vulnerability Detection**: Identifies insecure CORS configurations, including the critical 'null' origin issue
- **Missing Security Headers**: Detects when security headers should be added to HTTP responses
- **Authentication Issues**: Finds weak JWT configurations and insecure authentication patterns
- **SQL Injection Risks**: Identifies potential SQL injection vulnerabilities from string concatenation
- **Information Leakage**: Detects debug information and detailed error messages that could leak sensitive data
- **Comprehensive Reporting**: Provides detailed reports with severity levels and fix recommendations

## Installation

```bash
go build -o security-scanner ./cmd/security-scanner
```

## Usage

### Basic Scan

Scan the templates directory and output results to console:

```bash
./security-scanner -dir templates
```

### Verbose Output

Enable verbose logging to see scan progress:

```bash
./security-scanner -dir templates -verbose
```

### JSON Report

Generate a JSON report file:

```bash
./security-scanner -dir templates -output security-report.json
```

### Custom Directory

Scan a different directory:

```bash
./security-scanner -dir /path/to/templates -verbose
```

## Command Line Options

- `-dir`: Directory to scan for template files (default: "templates")
- `-output`: Output file for security report in JSON format (optional)
- `-verbose`: Enable verbose output showing scan progress

## Exit Codes

- `0`: No critical security issues found
- `1`: Critical security issues detected

## Security Issue Types

### CORS Vulnerabilities

- **Critical**: Setting Access-Control-Allow-Origin to 'null'
- **High**: Using wildcard (*) with credentials enabled
- **Medium**: Overly permissive wildcard CORS policies

### Authentication Issues

- **Critical**: JWT 'none' algorithm usage
- **High**: Weak or default JWT secrets
- **Medium**: Missing token expiration times
- **Low**: Cookie configurations without security flags

### SQL Injection Risks

- **Critical**: String concatenation in SQL queries
- **High**: Direct variable interpolation in SQL

### Information Leakage

- **Medium**: Detailed error messages exposing internal information
- **Medium**: Debug information enabled in production

### Missing Security Headers

- **Low**: HTTP responses without proper security headers

## Example Output

```
Security Scan Report
===================

File: templates/backend/go-gin/internal/middleware/cors.go.tmpl (Line 56)
Type: cors_vulnerability
Severity: critical
Description: Setting Access-Control-Allow-Origin to 'null' can allow bypass attacks
Recommendation: Omit the Access-Control-Allow-Origin header entirely for disallowed origins instead of setting it to 'null'

Summary: 138 issues found
Critical: 1, High: 13, Medium: 26, Low: 98
```

## Integration

The security scanner can be integrated into CI/CD pipelines to automatically detect security issues in template files:

```yaml
# GitHub Actions example
- name: Run Security Scanner
  run: |
    go build -o security-scanner ./cmd/security-scanner
    ./security-scanner -dir templates -output security-report.json
    
- name: Upload Security Report
  uses: actions/upload-artifact@v3
  with:
    name: security-report
    path: security-report.json
```

## Supported File Types

The scanner analyzes files with the following extensions:

- `.tmpl` - Template files
- `.go` - Go source files
- `.js` - JavaScript files
- `.ts` - TypeScript files
- `.yaml`, `.yml` - YAML configuration files
- `.json` - JSON configuration files

## Adding Custom Patterns

Security patterns are defined in `pkg/security/patterns.go`. To add new security checks:

1. Define a new `SecurityPattern` with appropriate regex
2. Set the issue type, severity, and recommendations
3. Add test cases in `pkg/security/patterns_test.go`

## Requirements Addressed

This security scanner addresses the following requirements from the template security improvements specification:

- **Requirement 3.1**: Automated detection of security vulnerabilities in template files
- **Requirement 3.2**: Pattern matching for CORS vulnerabilities and missing security headers
- **Requirement 3.3**: Reporting system with severity levels and fix recommendations
- **Requirement 3.4**: Foundation for security validation and regression testing
