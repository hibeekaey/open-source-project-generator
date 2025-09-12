#!/bin/bash

# Script to scan for hardcoded secrets and sensitive data
set -e

echo "=== Scanning for hardcoded secrets and sensitive data ==="

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to scan for common secret patterns
scan_secrets() {
    echo "Scanning for common secret patterns..."
    
    # Common secret patterns
    PATTERNS=(
        "password\s*=\s*['\"][^'\"]{3,}['\"]"
        "api[_-]?key\s*=\s*['\"][^'\"]{10,}['\"]"
        "secret[_-]?key\s*=\s*['\"][^'\"]{10,}['\"]"
        "access[_-]?token\s*=\s*['\"][^'\"]{10,}['\"]"
        "auth[_-]?token\s*=\s*['\"][^'\"]{10,}['\"]"
        "private[_-]?key\s*=\s*['\"][^'\"]{10,}['\"]"
        "database[_-]?url\s*=\s*['\"][^'\"]{10,}['\"]"
        "connection[_-]?string\s*=\s*['\"][^'\"]{10,}['\"]"
        "jwt[_-]?secret\s*=\s*['\"][^'\"]{10,}['\"]"
        "encryption[_-]?key\s*=\s*['\"][^'\"]{10,}['\"]"
        "-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----"
        "-----BEGIN\s+CERTIFICATE-----"
        "sk-[a-zA-Z0-9]{48}"
        "xoxb-[0-9]{11}-[0-9]{11}-[a-zA-Z0-9]{24}"
        "ghp_[a-zA-Z0-9]{36}"
        "gho_[a-zA-Z0-9]{36}"
        "ghu_[a-zA-Z0-9]{36}"
        "ghs_[a-zA-Z0-9]{36}"
        "ghr_[a-zA-Z0-9]{36}"
        "AKIA[0-9A-Z]{16}"
        "[0-9a-f]{32}"
        "[0-9a-f]{40}"
        "[0-9a-f]{64}"
    )
    
    FOUND_SECRETS=()
    
    for pattern in "${PATTERNS[@]}"; do
        echo "Checking pattern: $pattern"
        
        # Search in Go files, excluding test files, documentation examples, and environment variable usage
        results=$(grep -r -i -E "$pattern" --include="*.go" --include="*.yaml" --include="*.yml" --include="*.json" --include="*.env*" --include="*.conf" --include="*.config" --include="*.tmpl" . 2>/dev/null | grep -v -E "(test\.go|_test\.go|\.md\.tmpl|example|placeholder|your-|shasum)" | grep -v '\$' | grep -v '{{' || true)
        
        if [ -n "$results" ]; then
            echo "⚠️  Potential secret found with pattern: $pattern"
            echo "$results"
            FOUND_SECRETS+=("Pattern '$pattern' found in codebase")
        fi
    done
    
    return ${#FOUND_SECRETS[@]}
}

# Function to scan for environment variable patterns in templates
scan_env_patterns() {
    echo "Scanning template files for proper environment variable usage..."
    
    # Look for hardcoded values that should be environment variables
    HARDCODED_ISSUES=()
    
    # Check for hardcoded database connections
    if grep -r -i "host.*localhost" --include="*.tmpl" templates/ 2>/dev/null; then
        HARDCODED_ISSUES+=("Found hardcoded localhost references in templates")
    fi
    
    # Check for hardcoded ports (except common defaults)
    if grep -r -E ":[0-9]{4,5}" --include="*.tmpl" templates/ | grep -v -E ":(3000|8080|5432|3306|6379|27017|9200|5672)" 2>/dev/null; then
        HARDCODED_ISSUES+=("Found potentially hardcoded ports in templates")
    fi
    
    # Check for hardcoded URLs
    if grep -r -E "https?://[^{]" --include="*.tmpl" templates/ 2>/dev/null; then
        HARDCODED_ISSUES+=("Found hardcoded URLs in templates")
    fi
    
    # Check for proper environment variable usage
    echo "Checking for proper environment variable patterns..."
    
    # Good patterns: {{.Env.VAR_NAME}} or ${VAR_NAME} or process.env.VAR_NAME
    GOOD_PATTERNS=(
        "{{\.Env\.[A-Z_]+}}"
        "\${[A-Z_]+}"
        "process\.env\.[A-Z_]+"
        "os\.Getenv\(['\"][A-Z_]+['\"]\)"
    )
    
    TEMPLATE_FILES=$(find templates -name "*.tmpl" -type f)
    
    for file in $TEMPLATE_FILES; do
        echo "Checking $file for environment variable usage..."
        
        # Check if file contains configuration that should use env vars
        if grep -q -E "(password|secret|key|token|url|host|port)" "$file" 2>/dev/null; then
            # Check if it uses proper env var patterns
            has_good_pattern=false
            for pattern in "${GOOD_PATTERNS[@]}"; do
                if grep -q -E "$pattern" "$file" 2>/dev/null; then
                    has_good_pattern=true
                    break
                fi
            done
            
            if [ "$has_good_pattern" = false ]; then
                echo "⚠️  $file may contain configuration that should use environment variables"
            else
                echo "✅ $file uses proper environment variable patterns"
            fi
        fi
    done
    
    return ${#HARDCODED_ISSUES[@]}
}

# Function to check for sensitive files
check_sensitive_files() {
    echo "Checking for sensitive files that shouldn't be committed..."
    
    SENSITIVE_FILES=(
        "*.pem"
        "*.key"
        "*.p12"
        "*.pfx"
        "*.jks"
        "*.keystore"
        "*.crt"
        "*.cer"
        "*.der"
        ".env"
        ".env.local"
        ".env.production"
        "config.json"
        "secrets.yaml"
        "credentials.json"
        "service-account.json"
        "*.credentials"
    )
    
    FOUND_FILES=()
    
    for pattern in "${SENSITIVE_FILES[@]}"; do
        if find . -name "$pattern" -not -path "./.git/*" -not -path "./node_modules/*" -not -path "./vendor/*" 2>/dev/null | grep -q .; then
            echo "⚠️  Found potentially sensitive files matching pattern: $pattern"
            find . -name "$pattern" -not -path "./.git/*" -not -path "./node_modules/*" -not -path "./vendor/*" 2>/dev/null
            FOUND_FILES+=("$pattern")
        fi
    done
    
    return ${#FOUND_FILES[@]}
}

# Function to check .gitignore for sensitive patterns
check_gitignore() {
    echo "Checking .gitignore for sensitive file patterns..."
    
    REQUIRED_PATTERNS=(
        "*.env"
        "*.key"
        "*.pem"
        "*.p12"
        "*.pfx"
        "*.jks"
        "*.keystore"
        "config/secrets*"
        "credentials*"
        ".DS_Store"
    )
    
    MISSING_PATTERNS=()
    
    if [ -f ".gitignore" ]; then
        for pattern in "${REQUIRED_PATTERNS[@]}"; do
            if ! grep -q "$pattern" .gitignore; then
                MISSING_PATTERNS+=("$pattern")
            fi
        done
        
        if [ ${#MISSING_PATTERNS[@]} -gt 0 ]; then
            echo "⚠️  Missing patterns in .gitignore:"
            printf '%s\n' "${MISSING_PATTERNS[@]}"
        else
            echo "✅ .gitignore contains appropriate sensitive file patterns"
        fi
    else
        echo "⚠️  No .gitignore file found"
        MISSING_PATTERNS=("${REQUIRED_PATTERNS[@]}")
    fi
    
    return ${#MISSING_PATTERNS[@]}
}

# Function to scan with truffleHog if available
scan_with_trufflehog() {
    if command_exists trufflehog; then
        echo "Running TruffleHog scan..."
        trufflehog filesystem . --only-verified --json > /tmp/trufflehog-results.json 2>/dev/null || true
        
        if [ -s /tmp/trufflehog-results.json ]; then
            echo "⚠️  TruffleHog found potential secrets:"
            cat /tmp/trufflehog-results.json
            rm -f /tmp/trufflehog-results.json
            return 1
        else
            echo "✅ TruffleHog found no verified secrets"
            rm -f /tmp/trufflehog-results.json
            return 0
        fi
    else
        echo "TruffleHog not available. Install with: brew install trufflehog (macOS)"
        return 0
    fi
}

# Main execution
echo "Starting comprehensive secrets audit..."
echo

# Run all checks
TOTAL_ISSUES=0

echo "1. Scanning for secret patterns..."
scan_secrets
SECRET_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + SECRET_ISSUES))
echo

echo "2. Checking environment variable usage in templates..."
scan_env_patterns
ENV_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + ENV_ISSUES))
echo

echo "3. Checking for sensitive files..."
check_sensitive_files
FILE_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + FILE_ISSUES))
echo

echo "4. Checking .gitignore configuration..."
check_gitignore
GITIGNORE_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + GITIGNORE_ISSUES))
echo

echo "5. Running TruffleHog scan..."
scan_with_trufflehog
TRUFFLEHOG_ISSUES=$?
TOTAL_ISSUES=$((TOTAL_ISSUES + TRUFFLEHOG_ISSUES))
echo

# Summary
echo "========================================="
echo "Secrets and Sensitive Data Audit Summary"
echo "========================================="
echo "Total issues found: $TOTAL_ISSUES"

if [ $TOTAL_ISSUES -eq 0 ]; then
    echo "✅ No secrets or sensitive data issues found"
else
    echo "⚠️  Issues found that require attention"
    echo
    echo "Recommendations:"
    echo "1. Remove any hardcoded secrets from the codebase"
    echo "2. Use environment variables for all sensitive configuration"
    echo "3. Update .gitignore to exclude sensitive files"
    echo "4. Use template variables for configuration in template files"
    echo "5. Consider using secret management tools for production"
fi

echo "=== Secrets audit complete ==="

exit $TOTAL_ISSUES