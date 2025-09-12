#!/bin/bash

# Dependency Analysis Script
# Checks for unused dependencies and security vulnerabilities

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
AUDIT_DIR="${PROJECT_ROOT}/audit-results"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[DEPS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[DEPS-WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[DEPS-ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[DEPS-OK]${NC} $1"
}

# Check Go dependencies
check_go_dependencies() {
    log_info "Analyzing Go dependencies..."
    
    cd "${PROJECT_ROOT}"
    
    # Check if go.mod exists
    if [ ! -f "go.mod" ]; then
        log_error "go.mod not found"
        return 1
    fi
    
    local issues=0
    
    # Run go mod tidy to check for unused dependencies
    log_info "Running go mod tidy..."
    if ! go mod tidy; then
        log_error "go mod tidy failed"
        ((issues++))
    else
        log_success "go mod tidy completed successfully"
    fi
    
    # Check for unused dependencies using go mod why
    log_info "Checking for unused dependencies..."
    local unused_deps=()
    
    # Get all dependencies
    while IFS= read -r dep; do
        if [ -n "$dep" ]; then
            # Check if dependency is used
            if ! go mod why "$dep" | grep -q "^#"; then
                log_info "Dependency $dep appears to be used"
            else
                log_warn "Potentially unused dependency: $dep"
                unused_deps+=("$dep")
                ((issues++))
            fi
        fi
    done < <(go list -m all | tail -n +2 | cut -d' ' -f1)
    
    # Save unused dependencies
    if [ ${#unused_deps[@]} -gt 0 ]; then
        printf '%s\n' "${unused_deps[@]}" > "${AUDIT_DIR}/unused-go-deps.txt"
        log_warn "Found ${#unused_deps[@]} potentially unused Go dependencies"
    else
        log_success "No unused Go dependencies found"
    fi
    
    return $issues
}

# Check for security vulnerabilities in Go dependencies
check_go_security() {
    log_info "Checking Go dependencies for security vulnerabilities..."
    
    cd "${PROJECT_ROOT}"
    
    local issues=0
    
    # Use go list to check for known vulnerabilities
    if command -v govulncheck >/dev/null 2>&1; then
        log_info "Running govulncheck..."
        if ! govulncheck ./...; then
            log_error "Security vulnerabilities found in Go dependencies"
            ((issues++))
        else
            log_success "No security vulnerabilities found in Go dependencies"
        fi
    else
        log_warn "govulncheck not available, installing..."
        if go install golang.org/x/vuln/cmd/govulncheck@latest; then
            if ! govulncheck ./...; then
                log_error "Security vulnerabilities found in Go dependencies"
                ((issues++))
            else
                log_success "No security vulnerabilities found in Go dependencies"
            fi
        else
            log_warn "Could not install govulncheck, skipping vulnerability check"
        fi
    fi
    
    return $issues
}

# Check template dependencies (package.json files)
check_template_dependencies() {
    log_info "Analyzing template dependencies..."
    
    local issues=0
    
    # Find all package.json.tmpl files
    while IFS= read -r -d '' package_file; do
        log_info "Checking template: $package_file"
        
        # Check if file is valid JSON (ignoring template variables)
        local temp_file=$(mktemp)
        
        # Replace common template variables with dummy values for JSON validation
        sed -e 's/{{[^}]*}}/dummy/g' "$package_file" > "$temp_file"
        
        if command -v jq >/dev/null 2>&1; then
            if ! jq empty "$temp_file" 2>/dev/null; then
                log_warn "Invalid JSON structure in $package_file"
                ((issues++))
            else
                log_success "Valid JSON structure in $package_file"
            fi
        else
            log_warn "jq not available, skipping JSON validation"
        fi
        
        rm -f "$temp_file"
        
    done < <(find "${PROJECT_ROOT}/templates" -name "package.json.tmpl" -print0 2>/dev/null || true)
    
    return $issues
}

# Check Docker dependencies
check_docker_dependencies() {
    log_info "Analyzing Docker dependencies..."
    
    local issues=0
    
    # Find all Dockerfile templates
    while IFS= read -r -d '' dockerfile; do
        log_info "Checking Dockerfile: $dockerfile"
        
        # Check for outdated base images
        while IFS= read -r line; do
            if [[ "$line" =~ ^FROM[[:space:]]+([^[:space:]]+) ]]; then
                local image="${BASH_REMATCH[1]}"
                
                # Check for specific outdated patterns
                case "$image" in
                    *:18|*:16|*:14)
                        log_warn "Potentially outdated base image in $dockerfile: $image"
                        ((issues++))
                        ;;
                    *:latest)
                        log_warn "Using 'latest' tag in $dockerfile: $image (consider pinning version)"
                        ((issues++))
                        ;;
                    *)
                        log_info "Base image in $dockerfile: $image"
                        ;;
                esac
            fi
        done < "$dockerfile"
        
    done < <(find "${PROJECT_ROOT}/templates" -name "Dockerfile*.tmpl" -print0 2>/dev/null || true)
    
    return $issues
}

# Check Go version consistency
check_go_version_consistency() {
    log_info "Checking Go version consistency..."
    
    local issues=0
    
    # Get Go version from go.mod
    local go_mod_version=""
    if [ -f "${PROJECT_ROOT}/go.mod" ]; then
        go_mod_version=$(grep "^go " "${PROJECT_ROOT}/go.mod" | cut -d' ' -f2)
        log_info "Go version in go.mod: $go_mod_version"
    fi
    
    # Check Go version in template go.mod files
    while IFS= read -r -d '' template_go_mod; do
        local template_version=$(grep "^go " "$template_go_mod" | cut -d' ' -f2 || echo "")
        
        if [ -n "$template_version" ] && [ "$template_version" != "$go_mod_version" ]; then
            log_warn "Go version mismatch in $template_go_mod: $template_version (expected: $go_mod_version)"
            ((issues++))
        fi
        
    done < <(find "${PROJECT_ROOT}/templates" -name "go.mod.tmpl" -print0 2>/dev/null || true)
    
    return $issues
}

# Generate dependency report
generate_dependency_report() {
    log_info "Generating dependency report..."
    
    mkdir -p "${AUDIT_DIR}"
    
    cat > "${AUDIT_DIR}/dependency-report.md" << EOF
# Dependency Analysis Report

Generated on: $(date)

## Go Dependencies

### Main Dependencies
$(cd "${PROJECT_ROOT}" && go list -m all | head -20)

### Dependency Count
- Total dependencies: $(cd "${PROJECT_ROOT}" && go list -m all | wc -l)
- Direct dependencies: $(cd "${PROJECT_ROOT}" && go list -m -f '{{if not .Indirect}}{{.Path}}{{end}}' all | grep -v "^$" | wc -l)

## Template Dependencies

### Package.json Templates
$(find "${PROJECT_ROOT}/templates" -name "package.json.tmpl" | wc -l) package.json template files found

### Dockerfile Templates  
$(find "${PROJECT_ROOT}/templates" -name "Dockerfile*.tmpl" | wc -l) Dockerfile template files found

## Recommendations

- Keep dependencies up to date
- Remove unused dependencies
- Pin Docker image versions
- Regular security audits

EOF

    log_success "Dependency report generated at ${AUDIT_DIR}/dependency-report.md"
}

# Main dependency analysis
main() {
    log_info "Starting dependency analysis..."
    
    local total_issues=0
    
    # Run all checks
    check_go_dependencies || ((total_issues+=$?))
    check_go_security || ((total_issues+=$?))
    check_template_dependencies || ((total_issues+=$?))
    check_docker_dependencies || ((total_issues+=$?))
    check_go_version_consistency || ((total_issues+=$?))
    
    # Generate report
    generate_dependency_report
    
    # Save results
    mkdir -p "${AUDIT_DIR}"
    echo "dependency_issues: $total_issues" >> "${AUDIT_DIR}/dependency-results.txt"
    
    if [ $total_issues -eq 0 ]; then
        log_success "Dependency analysis completed with no issues"
    else
        log_warn "Dependency analysis completed with $total_issues issues"
    fi
    
    return $total_issues
}

main "$@"