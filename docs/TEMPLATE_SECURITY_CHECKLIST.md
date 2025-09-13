# Template Security Checklist

## Overview

This checklist ensures that all template files follow security best practices and are free from common vulnerabilities. Use this checklist during template development, code reviews, and security audits.

## Pre-Development Checklist

### Planning Phase

- [ ] Security requirements identified and documented
- [ ] Threat model created for the template type
- [ ] Security controls planned for implementation
- [ ] Environment-specific security configurations defined

### Design Phase

- [ ] Security architecture reviewed
- [ ] Authentication and authorization flows designed
- [ ] Input validation strategy defined
- [ ] Error handling approach planned

## Development Checklist

### CORS Security

- [ ] **No 'null' values**: Never set `Access-Control-Allow-Origin` to 'null'
- [ ] **Explicit origins**: Use explicit origin validation, avoid wildcards in production
- [ ] **Omit headers**: For disallowed origins, omit CORS headers entirely
- [ ] **Credential handling**: Only allow credentials for trusted origins
- [ ] **Method restrictions**: Limit allowed HTTP methods to necessary ones
- [ ] **Header restrictions**: Explicitly define allowed headers
- [ ] **Preflight handling**: Proper OPTIONS request handling implemented

```go
// ✅ Correct CORS implementation
if isAllowedOrigin(origin) {
    c.Header("Access-Control-Allow-Origin", origin)
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
// Don't set headers for disallowed origins
```

### HTTP Security Headers

- [ ] **X-Content-Type-Options**: Set to 'nosniff'
- [ ] **X-Frame-Options**: Set to 'DENY' or 'SAMEORIGIN' as appropriate
- [ ] **X-XSS-Protection**: Set to '1; mode=block'
- [ ] **Referrer-Policy**: Set to 'strict-origin-when-cross-origin'
- [ ] **Content-Security-Policy**: Implemented with restrictive defaults
- [ ] **Strict-Transport-Security**: Enabled for HTTPS endpoints
- [ ] **Permissions-Policy**: Configured to restrict unnecessary features

```go
// ✅ Security headers middleware
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        if c.Request.TLS != nil {
            c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        c.Next()
    }
}
```

### Authentication Security

- [ ] **JWT algorithm validation**: Explicitly validate signing algorithm
- [ ] **Strong secrets**: Use cryptographically secure, randomly generated secrets
- [ ] **Token expiration**: Appropriate TTL (15 minutes for access tokens)
- [ ] **Refresh token security**: Secure refresh token implementation
- [ ] **Session security**: HttpOnly, Secure, SameSite cookie attributes
- [ ] **Password hashing**: Use bcrypt with cost factor ≥ 12
- [ ] **Rate limiting**: Implement authentication rate limiting

```go
// ✅ Secure JWT validation
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(secret), nil
})
```

### Database Security

- [ ] **Parameterized queries**: All database queries use parameters
- [ ] **No string concatenation**: Never build queries with string concatenation
- [ ] **Input validation**: All database inputs validated before use
- [ ] **Connection security**: SSL/TLS enabled for database connections
- [ ] **Least privilege**: Database user has minimal required permissions
- [ ] **Connection pooling**: Secure connection pool configuration

```go
// ✅ Secure database query
query := "SELECT * FROM users WHERE email = ? AND status = ?"
row := db.QueryRow(query, email, "active")
```

### Input Validation and Sanitization

- [ ] **Server-side validation**: All inputs validated on server side
- [ ] **Validation library**: Use established validation library (e.g., validator/v10)
- [ ] **Whitelist approach**: Use whitelist validation where possible
- [ ] **HTML escaping**: All HTML output properly escaped
- [ ] **SQL injection prevention**: Parameterized queries used exclusively
- [ ] **File upload security**: Secure file upload handling if applicable
- [ ] **Size limits**: Appropriate input size limits implemented

```go
// ✅ Input validation structure
type UserInput struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Age      int    `json:"age" validate:"min=0,max=120"`
}
```

### Error Handling Security

- [ ] **Generic error messages**: No sensitive information in user-facing errors
- [ ] **Detailed logging**: Full error details logged for debugging
- [ ] **Error codes**: Consistent error code structure
- [ ] **Stack trace protection**: No stack traces in production responses
- [ ] **Security event logging**: Security-relevant events logged
- [ ] **Log sanitization**: Sensitive data redacted from logs

```go
// ✅ Secure error handling
if err != nil {
    logger.Error("Database operation failed", "error", err, "user_id", userID)
    c.JSON(500, gin.H{"error": "Internal server error", "code": "DB_ERROR"})
    return
}
```

### Configuration Security

- [ ] **Environment variables**: Secrets loaded from environment, not config files
- [ ] **Environment-specific configs**: Different security settings per environment
- [ ] **Secure defaults**: All configurations have secure default values
- [ ] **Configuration validation**: Config values validated at startup
- [ ] **Secret rotation**: Support for secret rotation without downtime
- [ ] **No hardcoded secrets**: No secrets in source code

```yaml
# ✅ Secure configuration structure
security:
  cors:
    allowed_origins: ["${ALLOWED_ORIGINS}"]
  jwt:
    secret: "${JWT_SECRET}"  # From environment
    ttl: "15m"
```

### Logging Security

- [ ] **Sensitive data redaction**: Passwords, tokens, secrets redacted
- [ ] **Security event logging**: Authentication, authorization events logged
- [ ] **Log integrity**: Logs protected from tampering
- [ ] **Log retention**: Appropriate log retention policies
- [ ] **Structured logging**: Use structured logging format
- [ ] **Log level configuration**: Appropriate log levels for different environments

```go
// ✅ Secure logging with redaction
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "login",
    "ip":      clientIP,
    // password field automatically redacted by hook
}).Info("User login attempt")
```

## Code Review Checklist

### Security Anti-patterns to Look For

- [ ] **CORS 'null' values**: Check for `Access-Control-Allow-Origin: null`
- [ ] **String concatenation in SQL**: Look for query building with `+` or `fmt.Sprintf`
- [ ] **Hardcoded secrets**: Search for hardcoded passwords, API keys, tokens
- [ ] **Missing input validation**: Ensure all user inputs are validated
- [ ] **Information leakage**: Check error messages for sensitive information
- [ ] **Insecure defaults**: Verify all defaults are secure
- [ ] **Missing security headers**: Ensure all responses include security headers
- [ ] **Weak authentication**: Check for insecure JWT or session handling

### Security Controls to Verify

- [ ] **Authentication middleware**: Applied to all protected routes
- [ ] **Authorization checks**: Proper permission verification
- [ ] **Input sanitization**: All inputs properly sanitized
- [ ] **Output encoding**: All outputs properly encoded
- [ ] **Rate limiting**: Implemented where appropriate
- [ ] **HTTPS enforcement**: HTTPS required for sensitive operations
- [ ] **Security headers**: Comprehensive security headers implemented

### Configuration Review

- [ ] **Environment separation**: Different configs for dev/staging/prod
- [ ] **Secret management**: Secrets properly externalized
- [ ] **Security settings**: All security features enabled
- [ ] **Default values**: Secure defaults for all settings
- [ ] **Validation**: Configuration validation implemented

## Testing Checklist

### Automated Security Tests

- [ ] **CORS tests**: Test various origin scenarios
- [ ] **Security header tests**: Verify all headers present and correct
- [ ] **Authentication tests**: Test token validation and expiration
- [ ] **Input validation tests**: Test with malicious inputs
- [ ] **SQL injection tests**: Verify parameterized query protection
- [ ] **XSS prevention tests**: Test output encoding
- [ ] **Error handling tests**: Verify no information leakage

### Manual Security Testing

- [ ] **Penetration testing**: Basic security testing performed
- [ ] **CORS bypass attempts**: Test for CORS bypass techniques
- [ ] **Authentication bypass**: Test authentication mechanisms
- [ ] **Input fuzzing**: Test with various malicious inputs
- [ ] **Error message analysis**: Review all error responses
- [ ] **Configuration testing**: Test with various configurations

### Security Regression Tests

- [ ] **Previous vulnerabilities**: Tests for previously fixed issues
- [ ] **Common vulnerabilities**: Tests for OWASP Top 10
- [ ] **Framework-specific issues**: Tests for framework-specific vulnerabilities
- [ ] **Dependency vulnerabilities**: Tests for known dependency issues

## Deployment Checklist

### Pre-deployment Security

- [ ] **Security scan**: Automated security scanning completed
- [ ] **Dependency audit**: All dependencies scanned for vulnerabilities
- [ ] **Configuration review**: Production configuration reviewed
- [ ] **Secret verification**: All secrets properly configured
- [ ] **SSL/TLS verification**: HTTPS properly configured
- [ ] **Security headers verification**: All headers properly set

### Post-deployment Verification

- [ ] **Security headers check**: Verify headers in production
- [ ] **CORS verification**: Test CORS configuration in production
- [ ] **Authentication testing**: Verify authentication works correctly
- [ ] **Error handling verification**: Ensure no information leakage
- [ ] **Logging verification**: Confirm security events are logged
- [ ] **Monitoring setup**: Security monitoring configured

## Maintenance Checklist

### Regular Security Reviews

- [ ] **Monthly security review**: Regular review of security configurations
- [ ] **Dependency updates**: Regular updates of security-related dependencies
- [ ] **Vulnerability scanning**: Regular automated vulnerability scans
- [ ] **Log analysis**: Regular review of security logs
- [ ] **Configuration drift**: Check for configuration changes

### Security Updates

- [ ] **Security patches**: Timely application of security patches
- [ ] **Framework updates**: Keep frameworks updated for security fixes
- [ ] **Dependency updates**: Regular updates of all dependencies
- [ ] **Configuration updates**: Update security configurations as needed

### Incident Response

- [ ] **Incident procedures**: Security incident response procedures defined
- [ ] **Contact information**: Security team contact information available
- [ ] **Escalation procedures**: Clear escalation procedures defined
- [ ] **Recovery procedures**: Security incident recovery procedures defined

## Template-Specific Checklists

### Backend API Templates

- [ ] **API authentication**: Proper API authentication implemented
- [ ] **Rate limiting**: API rate limiting configured
- [ ] **Input validation**: All API inputs validated
- [ ] **Output sanitization**: All API outputs sanitized
- [ ] **CORS configuration**: API CORS properly configured
- [ ] **Security headers**: API security headers implemented

### Frontend Templates

- [ ] **CSP configuration**: Content Security Policy properly configured
- [ ] **XSS prevention**: XSS prevention measures implemented
- [ ] **Secure communication**: HTTPS enforced for all communications
- [ ] **Authentication handling**: Secure token storage and handling
- [ ] **Input sanitization**: All user inputs sanitized

### Database Templates

- [ ] **Connection security**: Secure database connections
- [ ] **Query parameterization**: All queries properly parameterized
- [ ] **Access controls**: Proper database access controls
- [ ] **Encryption**: Data encryption at rest and in transit
- [ ] **Backup security**: Secure backup procedures

### Infrastructure Templates

- [ ] **Network security**: Proper network security configurations
- [ ] **Access controls**: Infrastructure access controls implemented
- [ ] **Monitoring**: Security monitoring configured
- [ ] **Logging**: Comprehensive security logging
- [ ] **Backup and recovery**: Secure backup and recovery procedures

## Compliance Checklist

### OWASP Top 10 Compliance

- [ ] **A01 - Broken Access Control**: Proper access controls implemented
- [ ] **A02 - Cryptographic Failures**: Strong cryptography used
- [ ] **A03 - Injection**: Injection attacks prevented
- [ ] **A04 - Insecure Design**: Secure design principles followed
- [ ] **A05 - Security Misconfiguration**: Secure configurations implemented
- [ ] **A06 - Vulnerable Components**: Dependencies regularly updated
- [ ] **A07 - Authentication Failures**: Strong authentication implemented
- [ ] **A08 - Software Integrity Failures**: Software integrity verified
- [ ] **A09 - Logging Failures**: Comprehensive logging implemented
- [ ] **A10 - Server-Side Request Forgery**: SSRF attacks prevented

### Industry Standards

- [ ] **NIST Cybersecurity Framework**: Relevant controls implemented
- [ ] **ISO 27001**: Information security management practices followed
- [ ] **SOC 2**: Relevant security controls implemented
- [ ] **PCI DSS**: Payment card security requirements met (if applicable)
- [ ] **GDPR**: Data protection requirements met (if applicable)
- [ ] **HIPAA**: Healthcare security requirements met (if applicable)

## Documentation Requirements

### Security Documentation

- [ ] **Security architecture**: Security architecture documented
- [ ] **Threat model**: Threat model documented and updated
- [ ] **Security controls**: All security controls documented
- [ ] **Configuration guide**: Security configuration guide available
- [ ] **Incident procedures**: Security incident procedures documented

### User Documentation

- [ ] **Security guide**: User security guide available
- [ ] **Configuration examples**: Secure configuration examples provided
- [ ] **Best practices**: Security best practices documented
- [ ] **Troubleshooting**: Security troubleshooting guide available

## Sign-off Requirements

### Development Sign-off

- [ ] **Developer review**: Developer has reviewed all security requirements
- [ ] **Code review**: Security-focused code review completed
- [ ] **Testing**: All security tests passing
- [ ] **Documentation**: Security documentation complete

### Security Team Sign-off

- [ ] **Security review**: Security team review completed
- [ ] **Penetration testing**: Security testing completed
- [ ] **Compliance review**: Compliance requirements verified
- [ ] **Risk assessment**: Security risk assessment completed

### Final Approval

- [ ] **All checklist items**: All applicable checklist items completed
- [ ] **Documentation**: All required documentation complete
- [ ] **Testing**: All security tests passing
- [ ] **Approvals**: All required approvals obtained

---

**Note**: This checklist should be customized based on specific template types and organizational requirements. Not all items may be applicable to every template, but all applicable items should be completed before deployment.
