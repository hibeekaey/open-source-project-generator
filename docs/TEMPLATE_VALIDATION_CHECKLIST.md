# Template Validation Checklist

Use this checklist when reviewing template changes or adding new templates.

## Pre-Commit Validation

### Import Validation

- [ ] All used functions have corresponding import statements
- [ ] Standard library imports are grouped first and alphabetically ordered
- [ ] Third-party imports are grouped second and alphabetically ordered  
- [ ] Local imports are grouped last and alphabetically ordered
- [ ] No unused imports remain in the template
- [ ] Template variables in imports are properly formatted (e.g., `{{.ModuleName}}`)

### Function Usage Validation

- [ ] `time.Now()`, `time.Duration` → `"time"` import present
- [ ] `fmt.Sprintf()`, `fmt.Errorf()` → `"fmt"` import present
- [ ] `strings.Contains()`, `strings.Split()` → `"strings"` import present
- [ ] `strconv.Atoi()`, `strconv.Itoa()` → `"strconv"` import present
- [ ] `context.Context`, `context.Background()` → `"context"` import present
- [ ] `log.Printf()`, `log.Fatal()` → `"log"` import present
- [ ] `http.StatusOK`, `http.Request` → `"net/http"` import present
- [ ] `json.Marshal()`, `json.Unmarshal()` → `"encoding/json"` import present
- [ ] `errors.New()`, `errors.Is()` → `"errors"` import present
- [ ] `os.Getenv()`, `os.Exit()` → `"os"` import present

### Template Syntax Validation

- [ ] Template variables are properly escaped (e.g., `{{.Name}}`, `{{.PackageName}}`)
- [ ] Conditional blocks are properly closed (`{{if}}...{{end}}`)
- [ ] Template comments are used where appropriate (`{{/* comment */}}`)
- [ ] No hardcoded values that should be template variables

## Compilation Testing

### Automated Testing

- [ ] Run import detection: `go run scripts/validate-templates/main.go --check-imports`
- [ ] Generate test project: `go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-validation`
- [ ] Compile generated code: `cd test-validation && go mod tidy && go build ./...`
- [ ] Run template-specific tests if available

### Manual Testing

- [ ] Test template with minimal configuration
- [ ] Test template with full configuration options
- [ ] Verify generated code follows Go conventions
- [ ] Check that all template variables are properly substituted

## Code Quality

### Go Best Practices

- [ ] Package names follow Go conventions (lowercase, no underscores)
- [ ] Function names are properly capitalized (exported vs unexported)
- [ ] Error handling follows Go idioms
- [ ] Interface definitions are minimal and focused
- [ ] Struct field tags are properly formatted

### Security Considerations

#### CORS Security

- [ ] **No 'null' CORS headers**: Never set `Access-Control-Allow-Origin` to 'null'
- [ ] **Explicit origin validation**: Use explicit origin lists, avoid wildcards in production
- [ ] **Proper CORS omission**: For disallowed origins, omit CORS headers entirely
- [ ] **Credential handling**: Only allow credentials for trusted origins
- [ ] **Method restrictions**: Limit allowed HTTP methods to necessary ones
- [ ] **Header restrictions**: Explicitly define allowed headers

#### HTTP Security Headers

- [ ] **X-Content-Type-Options**: Set to 'nosniff' to prevent MIME sniffing
- [ ] **X-Frame-Options**: Set to 'DENY' or 'SAMEORIGIN' to prevent clickjacking
- [ ] **X-XSS-Protection**: Set to '1; mode=block' for legacy browser protection
- [ ] **Strict-Transport-Security**: Enabled for HTTPS endpoints with appropriate max-age
- [ ] **Content-Security-Policy**: Implemented with restrictive defaults
- [ ] **Referrer-Policy**: Set to 'strict-origin-when-cross-origin'
- [ ] **Permissions-Policy**: Configured to restrict unnecessary browser features

#### Authentication & Authorization

- [ ] **JWT algorithm validation**: Explicitly validate signing algorithms (prevent 'none' attacks)
- [ ] **Strong secrets**: Use cryptographically secure, randomly generated secrets
- [ ] **Token expiration**: Appropriate TTL (15 minutes for access tokens, 7 days for refresh)
- [ ] **Session security**: HttpOnly, Secure, SameSite cookie attributes properly set
- [ ] **Password hashing**: Use bcrypt with cost factor ≥ 12
- [ ] **Rate limiting**: Authentication endpoints have rate limiting protection
- [ ] **Session management**: Proper session regeneration to prevent fixation attacks

#### Database Security

- [ ] **Parameterized queries**: All database queries use parameters, no string concatenation
- [ ] **Input validation**: All database inputs validated before use
- [ ] **Connection security**: SSL/TLS enabled for database connections
- [ ] **Least privilege**: Database user has minimal required permissions
- [ ] **Query limits**: Appropriate LIMIT clauses to prevent resource exhaustion
- [ ] **Error handling**: Database errors don't leak sensitive information

#### Input Validation & Sanitization

- [ ] **Server-side validation**: All inputs validated on server side using validation library
- [ ] **Whitelist validation**: Use whitelist approach where possible
- [ ] **HTML escaping**: All HTML output properly escaped to prevent XSS
- [ ] **Size limits**: Appropriate input size limits implemented
- [ ] **File upload security**: Secure file upload handling with type/size restrictions
- [ ] **SQL injection prevention**: Parameterized queries used exclusively
- [ ] **Command injection prevention**: No direct command execution with user input

#### Error Handling Security

- [ ] **Generic error messages**: No sensitive information in user-facing errors
- [ ] **Detailed logging**: Full error details logged internally for debugging
- [ ] **Consistent error responses**: Uniform error response structure
- [ ] **Stack trace protection**: No stack traces in production responses
- [ ] **Security event logging**: Authentication/authorization failures logged
- [ ] **Log sanitization**: Sensitive data (passwords, tokens) redacted from logs

#### Configuration Security

- [ ] **Environment variables**: Secrets loaded from environment, not config files
- [ ] **No hardcoded secrets**: No secrets, API keys, or passwords in source code
- [ ] **Environment-specific configs**: Different security settings per environment
- [ ] **Secure defaults**: All configurations have secure default values
- [ ] **Configuration validation**: Config values validated at startup
- [ ] **Secret rotation support**: Architecture supports secret rotation without downtime

#### Logging Security

- [ ] **Sensitive data redaction**: Passwords, tokens, secrets automatically redacted
- [ ] **Security event logging**: Authentication, CORS violations, failed requests logged
- [ ] **Structured logging**: Use structured logging format for security analysis
- [ ] **Log integrity**: Logs protected from tampering where possible
- [ ] **Appropriate log levels**: Security events logged at appropriate levels

## Documentation

### Code Documentation

- [ ] Package documentation is present and accurate
- [ ] Exported functions have proper documentation comments
- [ ] Complex logic includes inline comments
- [ ] TODO comments include context and assignee

### Template Documentation

- [ ] Template purpose is documented
- [ ] Required template variables are listed
- [ ] Dependencies and prerequisites are noted
- [ ] Usage examples are provided where helpful

## Integration Testing

### Multi-Template Testing

- [ ] Related templates work together correctly
- [ ] Cross-references between templates are valid
- [ ] Shared components are properly imported
- [ ] Configuration consistency across templates

### Platform Testing

- [ ] Templates work on Linux
- [ ] Templates work on macOS  
- [ ] Templates work on Windows (if applicable)
- [ ] Docker builds succeed with generated code

## Performance Considerations

### Template Processing

- [ ] Template compilation is reasonably fast
- [ ] No unnecessary complexity in template logic
- [ ] File generation doesn't create excessive files
- [ ] Memory usage is reasonable during generation

### Generated Code Performance

- [ ] No obvious performance anti-patterns
- [ ] Database queries are efficient
- [ ] HTTP handlers don't block unnecessarily
- [ ] Resource cleanup is properly handled

## Maintenance

### Future Maintainability

- [ ] Template structure is clear and logical
- [ ] Changes can be made without breaking existing functionality
- [ ] Template variables are well-named and documented
- [ ] Dependencies are clearly identified

### Version Compatibility

- [ ] Go version requirements are documented
- [ ] Third-party package versions are appropriate
- [ ] Backward compatibility is considered
- [ ] Migration path is clear for breaking changes

## Final Verification

### End-to-End Testing

- [ ] Generate complete project using modified templates
- [ ] Run all tests in generated project
- [ ] Build and run generated application
- [ ] Verify all features work as expected

### Documentation Updates

- [ ] Update TEMPLATE_MAINTENANCE.md if new patterns are introduced
- [ ] Update TEMPLATE_QUICK_REFERENCE.md with new function mappings
- [ ] Update this checklist if new validation steps are needed
- [ ] Commit documentation changes with template changes

## Sign-off

- [ ] All checklist items completed
- [ ] Automated tests pass
- [ ] Manual testing completed
- [ ] Documentation updated
- [ ] Ready for code review

**Reviewer**: ________________  
**Date**: ________________  
**Template(s)**: ________________
