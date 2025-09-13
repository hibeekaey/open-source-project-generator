# Template Security Guide

## Overview

This guide provides comprehensive security guidelines for developing, maintaining, and using templates in the project generator. It covers security best practices, common vulnerabilities, and secure coding patterns to ensure generated applications are secure by default.

## Security Principles

### 1. Secure by Default

All templates should implement security best practices by default, requiring developers to explicitly opt-out of security features rather than opt-in.

### 2. Defense in Depth

Implement multiple layers of security controls to protect against various attack vectors.

### 3. Principle of Least Privilege

Grant only the minimum permissions necessary for functionality to work correctly.

### 4. Fail Securely

When security controls fail, the system should fail to a secure state rather than an insecure one.

## CORS Security

### The Problem

Cross-Origin Resource Sharing (CORS) misconfigurations can lead to security vulnerabilities. The most critical issue is setting `Access-Control-Allow-Origin` to 'null' for disallowed origins.

### Secure Implementation

#### ❌ Insecure Pattern

```go
// NEVER do this - setting to 'null' can be bypassed
if !isAllowedOrigin(origin) {
    c.Header("Access-Control-Allow-Origin", "null")
}
```

#### ✅ Secure Pattern

```go
// Correct approach - omit header for disallowed origins
if isAllowedOrigin(origin) {
    c.Header("Access-Control-Allow-Origin", origin)
    // Only set other CORS headers for allowed origins
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
    c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
// For disallowed origins, don't set any CORS headers
```

### CORS Configuration Best Practices

1. **Explicit Origin Validation**: Never use wildcards in production
2. **Credential Handling**: Only allow credentials for trusted origins
3. **Method Restrictions**: Limit allowed HTTP methods to what's necessary
4. **Header Restrictions**: Explicitly list allowed headers

```go
type CORSConfig struct {
    AllowedOrigins   []string `yaml:"allowed_origins"`
    AllowedMethods   []string `yaml:"allowed_methods"`
    AllowedHeaders   []string `yaml:"allowed_headers"`
    AllowCredentials bool     `yaml:"allow_credentials"`
    MaxAge          int      `yaml:"max_age"`
}

func (c *CORSConfig) IsOriginAllowed(origin string) bool {
    for _, allowed := range c.AllowedOrigins {
        if allowed == origin {
            return true
        }
    }
    return false
}
```

## HTTP Security Headers

### Essential Security Headers

All HTTP responses should include these security headers:

#### X-Content-Type-Options

Prevents MIME type sniffing attacks.

```go
c.Header("X-Content-Type-Options", "nosniff")
```

#### X-Frame-Options

Protects against clickjacking attacks.

```go
// For API endpoints
c.Header("X-Frame-Options", "DENY")

// For web applications that need to be embedded
c.Header("X-Frame-Options", "SAMEORIGIN")
```

#### X-XSS-Protection

Enables browser XSS filtering (legacy browsers).

```go
c.Header("X-XSS-Protection", "1; mode=block")
```

#### Strict-Transport-Security (HSTS)

Forces HTTPS connections.

```go
// Only set for HTTPS connections
if c.Request.TLS != nil {
    c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}
```

#### Content-Security-Policy (CSP)

Prevents XSS and data injection attacks.

```go
// Restrictive default policy
csp := "default-src 'self'; " +
       "script-src 'self' 'unsafe-inline'; " +
       "style-src 'self' 'unsafe-inline'; " +
       "img-src 'self' data: https:; " +
       "font-src 'self'; " +
       "connect-src 'self'; " +
       "frame-ancestors 'none'"
c.Header("Content-Security-Policy", csp)
```

### Security Headers Middleware Template

```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Basic security headers
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // HSTS for HTTPS
        if c.Request.TLS != nil {
            c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        // CSP for web applications
        if isWebEndpoint(c.Request.URL.Path) {
            c.Header("Content-Security-Policy", getCSPPolicy())
        }
        
        c.Next()
    }
}
```

## Authentication Security

### JWT Token Security

#### Secure JWT Configuration

```go
type JWTConfig struct {
    SigningMethod   string        `yaml:"signing_method"`   // Use RS256 or HS256
    AccessTokenTTL  time.Duration `yaml:"access_token_ttl"` // 15 minutes
    RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"` // 7 days
    Issuer          string        `yaml:"issuer"`
    Audience        string        `yaml:"audience"`
}
```

#### ❌ Insecure JWT Patterns

```go
// NEVER allow 'none' algorithm
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // This is vulnerable - doesn't validate algorithm
    return []byte(secret), nil
})

// NEVER use weak secrets
secret := "secret123"
```

#### ✅ Secure JWT Patterns

```go
// Always validate the signing method
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Validate the algorithm
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(secret), nil
})

// Use strong, randomly generated secrets
secret := generateSecureSecret(32) // 256-bit secret
```

### Session Security

#### Secure Session Configuration

```go
store := sessions.NewCookieStore([]byte(sessionSecret))
store.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   3600, // 1 hour
    HttpOnly: true, // Prevent XSS access
    Secure:   true, // HTTPS only
    SameSite: http.SameSiteStrictMode, // CSRF protection
}
```

### Password Security

#### Secure Password Hashing

```go
import "golang.org/x/crypto/bcrypt"

// Hash password with appropriate cost
func HashPassword(password string) (string, error) {
    // Use cost of 12 or higher for production
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

// Verify password
func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

## Database Security

### SQL Injection Prevention

#### ❌ Vulnerable Patterns

```go
// NEVER use string concatenation for queries
query := "SELECT * FROM users WHERE id = " + userID
```

#### ✅ Secure Patterns

```go
// Always use parameterized queries
query := "SELECT * FROM users WHERE id = ?"
row := db.QueryRow(query, userID)

// Or with named parameters
query := "SELECT * FROM users WHERE email = :email AND status = :status"
rows, err := db.NamedQuery(query, map[string]interface{}{
    "email":  email,
    "status": "active",
})
```

### Database Connection Security

```go
type DatabaseConfig struct {
    Host            string `yaml:"host"`
    Port            int    `yaml:"port"`
    Database        string `yaml:"database"`
    Username        string `yaml:"username"`
    Password        string `yaml:"password"`
    SSLMode         string `yaml:"ssl_mode"`         // require, verify-full
    MaxConnections  int    `yaml:"max_connections"`
    ConnMaxLifetime string `yaml:"conn_max_lifetime"`
}

func (cfg *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
}
```

## Input Validation and Sanitization

### Input Validation Patterns

```go
import (
    "github.com/go-playground/validator/v10"
    "html"
    "strings"
)

type UserInput struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Age      int    `json:"age" validate:"min=0,max=120"`
}

func ValidateInput(input interface{}) error {
    validate := validator.New()
    return validate.Struct(input)
}

// Sanitize HTML input
func SanitizeHTML(input string) string {
    return html.EscapeString(strings.TrimSpace(input))
}
```

## Error Handling Security

### Secure Error Responses

#### ❌ Information Leakage

```go
// NEVER expose internal errors to users
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()}) // May leak sensitive info
}
```

#### ✅ Secure Error Handling

```go
// Log detailed errors internally, return generic messages to users
if err != nil {
    logger.Error("Database connection failed", "error", err, "user_id", userID)
    c.JSON(500, gin.H{"error": "Internal server error"})
    return
}

// For validation errors, be specific but safe
if validationErr != nil {
    c.JSON(400, gin.H{"error": "Invalid input format"})
    return
}
```

### Error Response Structure

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}

// Production error handler
func HandleError(c *gin.Context, err error, userMessage string) {
    // Log full error details
    logger.Error("Request failed", 
        "error", err,
        "path", c.Request.URL.Path,
        "method", c.Request.Method,
        "user_agent", c.Request.UserAgent(),
    )
    
    // Return safe error to user
    c.JSON(500, ErrorResponse{
        Error: userMessage,
        Code:  "INTERNAL_ERROR",
    })
}
```

## Configuration Security

### Environment-Specific Security

#### Development Configuration

```yaml
# configs/development.yaml
security:
  cors:
    allowed_origins: ["http://localhost:3000", "http://localhost:8080"]
    allow_credentials: true
  headers:
    csp_enabled: false  # More permissive for development
  auth:
    jwt_ttl: "1h"      # Longer for development convenience
```

#### Production Configuration

```yaml
# configs/production.yaml
security:
  cors:
    allowed_origins: ["https://yourdomain.com"]
    allow_credentials: true
  headers:
    csp_enabled: true
    hsts_enabled: true
  auth:
    jwt_ttl: "15m"     # Short-lived tokens
    require_https: true
```

### Secret Management

```go
import "os"

type SecurityConfig struct {
    JWTSecret     string `yaml:"-"` // Never in config files
    DatabaseURL   string `yaml:"-"` // Load from environment
    APIKeys       map[string]string `yaml:"-"`
}

func LoadSecurityConfig() *SecurityConfig {
    return &SecurityConfig{
        JWTSecret:   os.Getenv("JWT_SECRET"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
        APIKeys: map[string]string{
            "stripe": os.Getenv("STRIPE_API_KEY"),
            "sendgrid": os.Getenv("SENDGRID_API_KEY"),
        },
    }
}
```

## Logging Security

### Secure Logging Practices

```go
import "github.com/sirupsen/logrus"

// Configure secure logging
func SetupSecureLogging() *logrus.Logger {
    logger := logrus.New()
    
    // Don't log sensitive fields
    logger.AddHook(&SensitiveDataHook{
        SensitiveFields: []string{"password", "token", "secret", "key"},
    })
    
    return logger
}

// Custom hook to redact sensitive data
type SensitiveDataHook struct {
    SensitiveFields []string
}

func (hook *SensitiveDataHook) Fire(entry *logrus.Entry) error {
    for _, field := range hook.SensitiveFields {
        if _, exists := entry.Data[field]; exists {
            entry.Data[field] = "[REDACTED]"
        }
    }
    return nil
}
```

### Security Event Logging

```go
// Log security events for monitoring
func LogSecurityEvent(event string, details map[string]interface{}) {
    logger.WithFields(logrus.Fields{
        "event_type": "security",
        "event":      event,
        "timestamp":  time.Now().UTC(),
        "details":    details,
    }).Info("Security event")
}

// Usage examples
LogSecurityEvent("failed_login", map[string]interface{}{
    "username": username,
    "ip":       clientIP,
    "attempts": attemptCount,
})

LogSecurityEvent("cors_violation", map[string]interface{}{
    "origin":    origin,
    "endpoint":  endpoint,
    "blocked":   true,
})
```

## Template Development Guidelines

### Security Review Checklist

Before submitting template changes, ensure:

1. **CORS Configuration**
   - [ ] No 'null' values for Access-Control-Allow-Origin
   - [ ] Explicit origin validation
   - [ ] Appropriate credential handling

2. **Security Headers**
   - [ ] All required security headers present
   - [ ] CSP policy appropriate for application type
   - [ ] HSTS enabled for HTTPS endpoints

3. **Authentication**
   - [ ] Secure JWT configuration
   - [ ] Proper algorithm validation
   - [ ] Appropriate token expiration times

4. **Input Validation**
   - [ ] All user inputs validated
   - [ ] Parameterized database queries
   - [ ] HTML output properly escaped

5. **Error Handling**
   - [ ] No sensitive information in error responses
   - [ ] Proper logging of security events
   - [ ] Generic error messages for users

6. **Configuration**
   - [ ] Secrets loaded from environment variables
   - [ ] Environment-specific security settings
   - [ ] Secure defaults for all configurations

### Code Review Guidelines

When reviewing template code:

1. **Look for Security Anti-patterns**
   - String concatenation in SQL queries
   - Hardcoded secrets or credentials
   - Missing input validation
   - Information leakage in error messages

2. **Verify Security Controls**
   - CORS headers properly configured
   - Authentication middleware applied to protected routes
   - Input validation on all user-provided data
   - Proper error handling throughout

3. **Check Configuration Security**
   - Environment-specific settings
   - Secure defaults
   - Proper secret management

## Common Vulnerabilities and Fixes

### 1. CORS Bypass

**Issue**: Setting Access-Control-Allow-Origin to 'null'
**Fix**: Omit header for disallowed origins

### 2. Missing Security Headers

**Issue**: Responses lack security headers
**Fix**: Implement comprehensive security headers middleware

### 3. JWT Algorithm Confusion

**Issue**: Not validating JWT signing algorithm
**Fix**: Explicitly validate expected algorithm

### 4. SQL Injection

**Issue**: String concatenation in queries
**Fix**: Use parameterized queries exclusively

### 5. Information Leakage

**Issue**: Detailed error messages to users
**Fix**: Generic user messages, detailed logging

### 6. Weak Session Configuration

**Issue**: Insecure cookie settings
**Fix**: HttpOnly, Secure, SameSite attributes

## Security Testing

### Automated Security Tests

```go
// Test CORS security
func TestCORSSecurity(t *testing.T) {
    tests := []struct {
        origin   string
        expected bool
    }{
        {"https://trusted.com", true},
        {"https://malicious.com", false},
        {"null", false},
        {"", false},
    }
    
    for _, test := range tests {
        allowed := isOriginAllowed(test.origin)
        assert.Equal(t, test.expected, allowed)
    }
}

// Test security headers
func TestSecurityHeaders(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    SecurityHeadersMiddleware()(c)
    
    assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
    assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
    assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
}
```

### Manual Security Testing

1. **CORS Testing**
   - Test with various origins
   - Verify preflight requests
   - Check credential handling

2. **Authentication Testing**
   - Test token validation
   - Verify expiration handling
   - Check algorithm validation

3. **Input Validation Testing**
   - Test with malicious inputs
   - Verify SQL injection protection
   - Check XSS prevention

## Conclusion

Security is not optional - it must be built into every template from the ground up. By following these guidelines and implementing the security patterns outlined in this document, we can ensure that all generated applications are secure by default and protected against common vulnerabilities.

Remember: Security is an ongoing process, not a one-time implementation. Regularly review and update security configurations as new threats emerge and best practices evolve.
