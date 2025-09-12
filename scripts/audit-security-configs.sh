#!/bin/bash

# Script to validate and improve security configurations
set -e

echo "=== Validating and improving security configurations ==="

# Function to check security headers in templates
check_security_headers() {
    echo "Checking for security headers in templates..."
    
    SECURITY_HEADERS=(
        "X-Content-Type-Options"
        "X-Frame-Options"
        "X-XSS-Protection"
        "Strict-Transport-Security"
        "Content-Security-Policy"
        "Referrer-Policy"
        "Permissions-Policy"
    )
    
    MISSING_HEADERS=()
    
    # Check Next.js config files for security headers
    NEXTJS_CONFIGS=$(find templates -name "next.config.js.tmpl" -type f)
    
    for config in $NEXTJS_CONFIGS; do
        echo "Checking $config for security headers..."
        
        for header in "${SECURITY_HEADERS[@]}"; do
            if ! grep -q "$header" "$config" 2>/dev/null; then
                MISSING_HEADERS+=("$header missing in $config")
            fi
        done
    done
    
    # Check for security middleware in backend templates
    MIDDLEWARE_FILES=$(find templates -path "*/middleware/*" -name "*.go.tmpl" -type f)
    
    for middleware in $MIDDLEWARE_FILES; do
        echo "Checking $middleware for security middleware..."
        
        # Check for CORS configuration
        if ! grep -q -i "cors" "$middleware" 2>/dev/null; then
            MISSING_HEADERS+=("CORS middleware missing in $middleware")
        fi
        
        # Check for security headers middleware
        if ! grep -q -i "security\|header" "$middleware" 2>/dev/null; then
            MISSING_HEADERS+=("Security headers middleware missing in $middleware")
        fi
    done
    
    return ${#MISSING_HEADERS[@]}
}

# Function to check authentication and authorization patterns
check_auth_patterns() {
    echo "Checking authentication and authorization patterns..."
    
    AUTH_ISSUES=()
    
    # Check for JWT secret configuration
    JWT_FILES=$(find templates -name "*.tmpl" -type f -exec grep -l -i "jwt" {} \; 2>/dev/null)
    
    for file in $JWT_FILES; do
        echo "Checking $file for JWT security..."
        
        # Check for proper JWT secret handling
        if grep -q -E "jwt.*secret.*=.*['\"][^'\"]{1,20}['\"]" "$file" 2>/dev/null; then
            AUTH_ISSUES+=("Potentially weak JWT secret in $file")
        fi
        
        # Check for JWT expiration
        if ! grep -q -i "exp\|expir" "$file" 2>/dev/null; then
            AUTH_ISSUES+=("JWT expiration not configured in $file")
        fi
    done
    
    # Check for password hashing
    PASSWORD_FILES=$(find templates -name "*.tmpl" -type f -exec grep -l -i "password" {} \; 2>/dev/null)
    
    for file in $PASSWORD_FILES; do
        echo "Checking $file for password security..."
        
        # Check for password hashing
        if grep -q -i "password" "$file" && ! grep -q -E "(bcrypt|scrypt|argon2|hash)" "$file" 2>/dev/null; then
            # Skip test files and documentation
            if [[ ! "$file" =~ (test|doc|readme) ]]; then
                AUTH_ISSUES+=("Password hashing not found in $file")
            fi
        fi
    done
    
    # Check for rate limiting
    RATE_LIMIT_FILES=$(find templates -name "*.tmpl" -type f -exec grep -l -i "rate\|limit" {} \; 2>/dev/null)
    
    if [ ${#RATE_LIMIT_FILES[@]} -eq 0 ]; then
        AUTH_ISSUES+=("Rate limiting not implemented in templates")
    fi
    
    return ${#AUTH_ISSUES[@]}
}

# Function to check Docker security configurations
check_docker_security() {
    echo "Checking Docker security configurations..."
    
    DOCKER_ISSUES=()
    
    DOCKERFILES=$(find templates -name "Dockerfile*.tmpl" -type f)
    
    for dockerfile in $DOCKERFILES; do
        echo "Checking $dockerfile for security best practices..."
        
        # Check for non-root user
        if ! grep -q "USER" "$dockerfile" 2>/dev/null; then
            DOCKER_ISSUES+=("Non-root user not configured in $dockerfile")
        fi
        
        # Check for HEALTHCHECK
        if ! grep -q "HEALTHCHECK" "$dockerfile" 2>/dev/null; then
            DOCKER_ISSUES+=("Health check not configured in $dockerfile")
        fi
        
        # Check for minimal base images
        if grep -q "FROM.*:latest" "$dockerfile" 2>/dev/null; then
            DOCKER_ISSUES+=("Using 'latest' tag in $dockerfile (use specific versions)")
        fi
        
        # Check for security updates
        if ! grep -q -E "(apk.*upgrade|apt.*upgrade|yum.*update)" "$dockerfile" 2>/dev/null; then
            DOCKER_ISSUES+=("Security updates not applied in $dockerfile")
        fi
    done
    
    return ${#DOCKER_ISSUES[@]}
}

# Function to check Kubernetes security configurations
check_k8s_security() {
    echo "Checking Kubernetes security configurations..."
    
    K8S_ISSUES=()
    
    K8S_FILES=$(find templates -name "*.yaml.tmpl" -o -name "*.yml.tmpl" | grep -E "(k8s|kubernetes)" 2>/dev/null)
    
    for k8s_file in $K8S_FILES; do
        echo "Checking $k8s_file for security configurations..."
        
        # Check for security context
        if grep -q "kind: Deployment" "$k8s_file" && ! grep -q "securityContext" "$k8s_file" 2>/dev/null; then
            K8S_ISSUES+=("Security context not configured in $k8s_file")
        fi
        
        # Check for resource limits
        if grep -q "kind: Deployment" "$k8s_file" && ! grep -q "resources:" "$k8s_file" 2>/dev/null; then
            K8S_ISSUES+=("Resource limits not configured in $k8s_file")
        fi
        
        # Check for network policies
        if ! find templates -name "*.tmpl" -exec grep -l "NetworkPolicy" {} \; 2>/dev/null | grep -q .; then
            K8S_ISSUES+=("Network policies not found in Kubernetes templates")
        fi
        
        # Check for pod security standards
        if grep -q "kind: Deployment" "$k8s_file" && ! grep -q -E "(runAsNonRoot|readOnlyRootFilesystem)" "$k8s_file" 2>/dev/null; then
            K8S_ISSUES+=("Pod security standards not configured in $k8s_file")
        fi
    done
    
    return ${#K8S_ISSUES[@]}
}

# Function to check database security configurations
check_database_security() {
    echo "Checking database security configurations..."
    
    DB_ISSUES=()
    
    # Check for SSL/TLS configuration
    DB_CONFIG_FILES=$(find templates -name "*.tmpl" -type f -exec grep -l -i "database\|postgres\|mysql\|mongo" {} \; 2>/dev/null)
    
    for config in $DB_CONFIG_FILES; do
        echo "Checking $config for database security..."
        
        # Check for SSL configuration
        if grep -q -i "database\|postgres" "$config" && ! grep -q -i "ssl\|tls" "$config" 2>/dev/null; then
            # Skip test files
            if [[ ! "$config" =~ test ]]; then
                DB_ISSUES+=("SSL/TLS not configured for database in $config")
            fi
        fi
        
        # Check for connection pooling
        if grep -q -i "database" "$config" && ! grep -q -i "pool\|connection" "$config" 2>/dev/null; then
            DB_ISSUES+=("Connection pooling not configured in $config")
        fi
    done
    
    return ${#DB_ISSUES[@]}
}

# Function to check API security configurations
check_api_security() {
    echo "Checking API security configurations..."
    
    API_ISSUES=()
    
    # Check for input validation
    API_FILES=$(find templates -name "*.tmpl" -type f -exec grep -l -i "api\|controller\|handler" {} \; 2>/dev/null)
    
    for api_file in $API_FILES; do
        echo "Checking $api_file for API security..."
        
        # Check for input validation
        if grep -q -i "controller\|handler" "$api_file" && ! grep -q -i "validat\|sanitiz" "$api_file" 2>/dev/null; then
            # Skip test files
            if [[ ! "$api_file" =~ test ]]; then
                API_ISSUES+=("Input validation not found in $api_file")
            fi
        fi
        
        # Check for error handling
        if grep -q -i "controller\|handler" "$api_file" && ! grep -q -i "error\|exception" "$api_file" 2>/dev/null; then
            API_ISSUES+=("Error handling not found in $api_file")
        fi
    done
    
    return ${#API_ISSUES[@]}
}

# Function to generate security recommendations
generate_recommendations() {
    echo "Generating security recommendations..."
    
    cat > security-recommendations.md << 'EOF'
# Security Configuration Recommendations

## Security Headers
- Implement Content Security Policy (CSP) headers
- Add X-Frame-Options to prevent clickjacking
- Set X-Content-Type-Options to prevent MIME sniffing
- Configure Strict-Transport-Security for HTTPS enforcement

## Authentication & Authorization
- Use strong JWT secrets (minimum 256 bits)
- Implement JWT token expiration
- Use bcrypt or Argon2 for password hashing
- Implement rate limiting for authentication endpoints

## Docker Security
- Always run containers as non-root users
- Use specific image tags instead of 'latest'
- Implement health checks
- Apply security updates in base images

## Kubernetes Security
- Configure security contexts for pods
- Set resource limits and requests
- Implement network policies
- Use pod security standards (runAsNonRoot, readOnlyRootFilesystem)

## Database Security
- Enable SSL/TLS for database connections
- Configure connection pooling
- Use parameterized queries to prevent SQL injection
- Implement database access controls

## API Security
- Validate and sanitize all input
- Implement proper error handling
- Use HTTPS for all API communications
- Implement API rate limiting

## General Recommendations
- Regular security audits and dependency updates
- Implement logging and monitoring
- Use secrets management systems
- Follow principle of least privilege
EOF

    echo "Security recommendations saved to security-recommendations.md"
}

# Main execution
echo "Starting security configuration audit..."
echo

TOTAL_ISSUES=0

echo "1. Checking security headers..."
check_security_headers
HEADER_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + HEADER_ISSUES))
echo

echo "2. Checking authentication and authorization patterns..."
check_auth_patterns
AUTH_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + AUTH_ISSUES))
echo

echo "3. Checking Docker security configurations..."
check_docker_security
DOCKER_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + DOCKER_ISSUES))
echo

echo "4. Checking Kubernetes security configurations..."
check_k8s_security
K8S_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + K8S_ISSUES))
echo

echo "5. Checking database security configurations..."
check_database_security
DB_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + DB_ISSUES))
echo

echo "6. Checking API security configurations..."
check_api_security
API_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + API_ISSUES))
echo

echo "7. Generating security recommendations..."
generate_recommendations
echo

# Summary
echo "========================================="
echo "Security Configuration Audit Summary"
echo "========================================="
echo "Total security issues found: $TOTAL_ISSUES"
echo "Security headers issues: $HEADER_ISSUES"
echo "Authentication/Authorization issues: $AUTH_ISSUES"
echo "Docker security issues: $DOCKER_ISSUES"
echo "Kubernetes security issues: $K8S_ISSUES"
echo "Database security issues: $DB_ISSUES"
echo "API security issues: $API_ISSUES"

if [ $TOTAL_ISSUES -eq 0 ]; then
    echo "✅ Security configurations look good"
else
    echo "⚠️  Security improvements recommended"
    echo "📋 Check security-recommendations.md for detailed guidance"
fi

echo "=== Security configuration audit complete ==="

exit 0