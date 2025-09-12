#!/bin/bash

# Codebase Audit Script
# This script performs comprehensive audit of the template generator codebase

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
AUDIT_DIR="${PROJECT_ROOT}/audit-results"
LOG_FILE="${AUDIT_DIR}/audit.log"
REPORT_FILE="${AUDIT_DIR}/audit-report.json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "${LOG_FILE}"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" | tee -a "${LOG_FILE}"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "${LOG_FILE}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "${LOG_FILE}"
}

# Initialize audit environment
init_audit() {
    log_info "Initializing audit environment..."
    
    # Create audit results directory
    mkdir -p "${AUDIT_DIR}"
    
    # Initialize log file
    echo "Audit started at $(date)" > "${LOG_FILE}"
    
    # Initialize report structure
    cat > "${REPORT_FILE}" << 'EOF'
{
    "audit_info": {
        "timestamp": "",
        "version": "1.0.0",
        "project_root": ""
    },
    "phases": [],
    "summary": {
        "total_issues": 0,
        "critical_issues": 0,
        "files_analyzed": 0,
        "tests_passing": 0,
        "tests_failing": 0,
        "coverage_percent": 0.0
    },
    "issues": [],
    "recommendations": []
}
EOF
    
    # Update report with current info
    update_report_info
    
    log_success "Audit environment initialized"
}

# Update report with basic info
update_report_info() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # Use jq if available, otherwise use sed
    if command -v jq >/dev/null 2>&1; then
        jq --arg ts "$timestamp" --arg root "$PROJECT_ROOT" \
           '.audit_info.timestamp = $ts | .audit_info.project_root = $root' \
           "${REPORT_FILE}" > "${REPORT_FILE}.tmp" && mv "${REPORT_FILE}.tmp" "${REPORT_FILE}"
    else
        sed -i.bak "s|\"timestamp\": \"\"|\"timestamp\": \"$timestamp\"|" "${REPORT_FILE}"
        sed -i.bak "s|\"project_root\": \"\"|\"project_root\": \"$PROJECT_ROOT\"|" "${REPORT_FILE}"
        rm -f "${REPORT_FILE}.bak"
    fi
}

# Add issue to report
add_issue() {
    local type="$1"
    local severity="$2"
    local file="$3"
    local line="${4:-0}"
    local description="$5"
    local suggestion="${6:-}"
    
    local issue_json=$(cat << EOF
{
    "type": "$type",
    "severity": "$severity", 
    "file": "$file",
    "line": $line,
    "description": "$description",
    "suggestion": "$suggestion"
}
EOF
)
    
    if command -v jq >/dev/null 2>&1; then
        echo "$issue_json" | jq -c '.' >> "${AUDIT_DIR}/issues.jsonl"
    else
        echo "$issue_json" >> "${AUDIT_DIR}/issues.jsonl"
    fi
}

# Check if required tools are installed
check_tools() {
    log_info "Checking required tools..."
    
    local missing_tools=()
    
    # Check Go
    if ! command -v go >/dev/null 2>&1; then
        missing_tools+=("go")
    fi
    
    # Check golangci-lint
    if ! command -v golangci-lint >/dev/null 2>&1; then
        log_warn "golangci-lint not found, will install it"
    fi
    
    # Check other tools
    for tool in git make docker; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_tools+=("$tool")
        fi
    done
    
    if [ ${#missing_tools[@]} -gt 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        return 1
    fi
    
    log_success "All required tools are available"
}

# Install golangci-lint if not present
install_golangci_lint() {
    if ! command -v golangci-lint >/dev/null 2>&1; then
        log_info "Installing golangci-lint..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
        log_success "golangci-lint installed"
    else
        log_info "golangci-lint already installed"
    fi
}

# Run structural analysis
run_structural_analysis() {
    log_info "Running structural analysis..."
    
    # Check Go project structure
    "${SCRIPT_DIR}/audit-structure.sh" 2>&1 | tee -a "${LOG_FILE}"
    
    log_success "Structural analysis completed"
}

# Run dependency analysis
run_dependency_analysis() {
    log_info "Running dependency analysis..."
    
    # Check Go dependencies
    "${SCRIPT_DIR}/audit-dependencies.sh" 2>&1 | tee -a "${LOG_FILE}"
    
    log_success "Dependency analysis completed"
}

# Run code quality analysis
run_code_quality_analysis() {
    log_info "Running code quality analysis..."
    
    # Run linting tools
    "${SCRIPT_DIR}/audit-quality.sh" 2>&1 | tee -a "${LOG_FILE}"
    
    log_success "Code quality analysis completed"
}

# Generate final report
generate_report() {
    log_info "Generating final audit report..."
    
    # Consolidate issues
    if [ -f "${AUDIT_DIR}/issues.jsonl" ]; then
        local issue_count=$(wc -l < "${AUDIT_DIR}/issues.jsonl" 2>/dev/null || echo "0")
        log_info "Found $issue_count issues"
    fi
    
    # Create summary report
    cat > "${AUDIT_DIR}/summary.md" << EOF
# Audit Summary Report

Generated on: $(date)

## Overview

This report contains the results of the comprehensive codebase audit.

## Files Analyzed

- Go source files: $(find "${PROJECT_ROOT}" -name "*.go" -not -path "*/vendor/*" | wc -l)
- Template files: $(find "${PROJECT_ROOT}/templates" -name "*.tmpl" | wc -l)
- Test files: $(find "${PROJECT_ROOT}" -name "*_test.go" | wc -l)

## Issues Found

$(if [ -f "${AUDIT_DIR}/issues.jsonl" ]; then
    echo "Total issues: $(wc -l < "${AUDIT_DIR}/issues.jsonl")"
else
    echo "No issues file found"
fi)

## Recommendations

See detailed report in audit-report.json for specific recommendations.

EOF
    
    log_success "Audit report generated at ${AUDIT_DIR}/summary.md"
}

# Main audit function
run_audit() {
    log_info "Starting comprehensive codebase audit..."
    
    # Initialize
    init_audit
    
    # Check tools
    check_tools || exit 1
    
    # Install missing tools
    install_golangci_lint
    
    # Run analysis phases
    run_structural_analysis
    run_dependency_analysis  
    run_code_quality_analysis
    
    # Generate report
    generate_report
    
    log_success "Audit completed successfully!"
    log_info "Results available in: ${AUDIT_DIR}"
}

# Help function
show_help() {
    cat << EOF
Usage: $0 [OPTION]

Comprehensive codebase audit script for the template generator project.

Options:
    --structure     Run only structural analysis
    --dependencies  Run only dependency analysis
    --quality       Run only code quality analysis
    --help          Show this help message

Examples:
    $0                    # Run full audit
    $0 --structure        # Run only structural analysis
    $0 --quality          # Run only code quality analysis

EOF
}

# Main script logic
main() {
    case "${1:-}" in
        --structure)
            init_audit
            check_tools || exit 1
            run_structural_analysis
            ;;
        --dependencies)
            init_audit
            check_tools || exit 1
            run_dependency_analysis
            ;;
        --quality)
            init_audit
            check_tools || exit 1
            install_golangci_lint
            run_code_quality_analysis
            ;;
        --help)
            show_help
            ;;
        "")
            run_audit
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"