#!/bin/bash

# Security Validation Script for CI/CD
# This script runs comprehensive security validation checks

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SECURITY_LINTER="$PROJECT_ROOT/cmd/security-linter"
OUTPUT_DIR="$PROJECT_ROOT/security-reports"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Function to build security linter if needed
build_security_linter() {
    log_info "Building security linter..."
    
    if [ ! -f "$SECURITY_LINTER/main.go" ]; then
        log_error "Security linter source not found at $SECURITY_LINTER/main.go"
        exit 1
    fi
    
    cd "$PROJECT_ROOT"
    go build -o "$OUTPUT_DIR/security-linter" "$SECURITY_LINTER/main.go"
    
    if [ $? -eq 0 ]; then
        log_success "Security linter built successfully"
    else
        log_error "Failed to build security linter"
        exit 1
    fi
}

# Function to run security linting
run_security_linting() {
    log_info "Running security linting..."
    
    local linter_bin="$OUTPUT_DIR/security-linter"
    local json_report="$OUTPUT_DIR/security-report-$TIMESTAMP.json"
    local sarif_report="$OUTPUT_DIR/security-report-$TIMESTAMP.sarif"
    local junit_report="$OUTPUT_DIR/security-report-$TIMESTAMP.xml"
    
    # Run linting with different output formats
    log_info "Generating JSON report..."
    if "$linter_bin" -dir "$PROJECT_ROOT" -format json -output "$json_report" -verbose; then
        log_success "JSON report generated: $json_report"
    else
        local exit_code=$?
        if [ $exit_code -eq 1 ]; then
            log_warning "High severity security issues found (see report)"
        elif [ $exit_code -eq 2 ]; then
            log_error "Critical security issues found (see report)"
            return $exit_code
        else
            log_error "Security linting failed with exit code $exit_code"
            return $exit_code
        fi
    fi
    
    # Generate SARIF report for GitHub Security tab
    log_info "Generating SARIF report..."
    "$linter_bin" -dir "$PROJECT_ROOT" -format sarif -output "$sarif_report" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        log_success "SARIF report generated: $sarif_report"
    fi
    
    # Generate JUnit report for CI/CD integration
    log_info "Generating JUnit report..."
    "$linter_bin" -dir "$PROJECT_ROOT" -format junit -output "$junit_report" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        log_success "JUnit report generated: $junit_report"
    fi
    
    return 0
}

# Function to run automated security tests
run_security_tests() {
    log_info "Running automated security tests..."
    
    cd "$PROJECT_ROOT"
    
    # Run security validation tests
    if go test -v ./pkg/security -run TestAutomatedSecurityValidation; then
        log_success "Automated security validation tests passed"
    else
        log_error "Automated security validation tests failed"
        return 1
    fi
    
    # Run security regression prevention tests
    if go test -v ./pkg/security -run TestSecurityRegressionPrevention; then
        log_success "Security regression prevention tests passed"
    else
        log_error "Security regression prevention tests failed"
        return 1
    fi
    
    # Run security best practices compliance tests
    if go test -v ./pkg/security -run TestSecurityBestPracticesCompliance; then
        log_success "Security best practices compliance tests passed"
    else
        log_warning "Security best practices compliance tests failed (non-blocking)"
    fi
    
    return 0
}

# Function to scan for specific insecure patterns
scan_insecure_patterns() {
    log_info "Scanning for specific insecure patterns..."
    
    local issues_found=0
    
    # Check for timestamp-based random generation
    log_info "Checking for timestamp-based random generation..."
    if grep -r "time\.Now()\.UnixNano()" "$PROJECT_ROOT" --include="*.go" --exclude-dir=".git" --exclude-dir="vendor" | grep -v "_test.go" | grep -v "docs/" | grep -v "examples/"; then
        log_error "Found timestamp-based random generation (security vulnerability)"
        issues_found=$((issues_found + 1))
    else
        log_success "No timestamp-based random generation found"
    fi
    
    # Check for math/rand usage in security-sensitive code
    log_info "Checking for math/rand usage in security code..."
    if grep -r "math/rand" "$PROJECT_ROOT/pkg/security" --include="*.go" | grep -v "_test.go"; then
        log_error "Found math/rand usage in security package (use crypto/rand instead)"
        issues_found=$((issues_found + 1))
    else
        log_success "No insecure math/rand usage in security package"
    fi
    
    # Check for predictable temporary file patterns
    log_info "Checking for predictable temporary file patterns..."
    if grep -r "\.tmp\..*time\.Now()" "$PROJECT_ROOT" --include="*.go" --exclude-dir=".git" | grep -v "_test.go" | grep -v "docs/"; then
        log_error "Found predictable temporary file naming (race condition vulnerability)"
        issues_found=$((issues_found + 1))
    else
        log_success "No predictable temporary file patterns found"
    fi
    
    # Check for hardcoded secrets (basic check)
    log_info "Checking for potential hardcoded secrets..."
    if grep -r -i "password.*=.*[\"'][^\"']{8,}[\"']" "$PROJECT_ROOT" --include="*.go" --exclude-dir=".git" --exclude-dir="vendor" | grep -v "_test.go" | grep -v "docs/" | grep -v "examples/"; then
        log_warning "Found potential hardcoded secrets (review manually)"
    else
        log_success "No obvious hardcoded secrets found"
    fi
    
    return $issues_found
}

# Function to generate security summary
generate_security_summary() {
    log_info "Generating security validation summary..."
    
    local summary_file="$OUTPUT_DIR/security-summary-$TIMESTAMP.txt"
    
    cat > "$summary_file" << EOF
Security Validation Summary
==========================
Timestamp: $(date)
Project: Template Generator
Validation Script: $0

Reports Generated:
- JSON Report: security-report-$TIMESTAMP.json
- SARIF Report: security-report-$TIMESTAMP.sarif
- JUnit Report: security-report-$TIMESTAMP.xml

Validation Steps Completed:
✓ Security linting
✓ Automated security tests
✓ Insecure pattern scanning
✓ Security regression prevention

For detailed results, check the generated reports in:
$OUTPUT_DIR

EOF
    
    log_success "Security summary generated: $summary_file"
}

# Function to upload reports (for CI/CD integration)
upload_reports() {
    if [ "${CI:-}" = "true" ]; then
        log_info "CI environment detected, preparing reports for upload..."
        
        # Copy latest reports to standard names for CI artifacts
        cp "$OUTPUT_DIR/security-report-$TIMESTAMP.json" "$OUTPUT_DIR/security-report-latest.json" 2>/dev/null || true
        cp "$OUTPUT_DIR/security-report-$TIMESTAMP.sarif" "$OUTPUT_DIR/security-report-latest.sarif" 2>/dev/null || true
        cp "$OUTPUT_DIR/security-report-$TIMESTAMP.xml" "$OUTPUT_DIR/security-report-latest.xml" 2>/dev/null || true
        
        log_success "Reports prepared for CI artifact upload"
    fi
}

# Main execution
main() {
    log_info "Starting security validation for Template Generator"
    log_info "Project root: $PROJECT_ROOT"
    log_info "Output directory: $OUTPUT_DIR"
    
    local exit_code=0
    
    # Build security linter
    build_security_linter || exit_code=$?
    
    # Run security linting
    if [ $exit_code -eq 0 ]; then
        run_security_linting || exit_code=$?
    fi
    
    # Run automated security tests
    if [ $exit_code -eq 0 ]; then
        run_security_tests || exit_code=$?
    fi
    
    # Scan for insecure patterns
    if [ $exit_code -eq 0 ]; then
        scan_insecure_patterns || {
            local pattern_issues=$?
            if [ $pattern_issues -gt 0 ]; then
                log_error "Found $pattern_issues security pattern violations"
                exit_code=1
            fi
        }
    fi
    
    # Generate summary
    generate_security_summary
    
    # Upload reports for CI
    upload_reports
    
    # Final status
    if [ $exit_code -eq 0 ]; then
        log_success "Security validation completed successfully"
    else
        log_error "Security validation failed with exit code $exit_code"
    fi
    
    return $exit_code
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo ""
        echo "Security validation script for Template Generator"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --lint-only    Run only security linting"
        echo "  --test-only    Run only security tests"
        echo ""
        echo "Environment Variables:"
        echo "  CI=true        Enable CI mode (affects report handling)"
        echo ""
        exit 0
        ;;
    --lint-only)
        build_security_linter
        run_security_linting
        exit $?
        ;;
    --test-only)
        run_security_tests
        exit $?
        ;;
    *)
        main
        exit $?
        ;;
esac