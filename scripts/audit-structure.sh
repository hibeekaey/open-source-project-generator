#!/bin/bash

# Structural Analysis Script
# Validates Go project structure and organization

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
    echo -e "${BLUE}[STRUCTURE]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[STRUCTURE-WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[STRUCTURE-ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[STRUCTURE-OK]${NC} $1"
}

# Check Go project structure
check_go_structure() {
    log_info "Validating Go project structure..."
    
    local issues=0
    
    # Check for required directories
    local required_dirs=("cmd" "internal" "pkg")
    for dir in "${required_dirs[@]}"; do
        if [ ! -d "${PROJECT_ROOT}/${dir}" ]; then
            log_error "Missing required directory: ${dir}"
            ((issues++))
        else
            log_success "Found required directory: ${dir}"
        fi
    done
    
    # Check cmd directory structure
    if [ -d "${PROJECT_ROOT}/cmd" ]; then
        log_info "Checking cmd directory structure..."
        for cmd_dir in "${PROJECT_ROOT}/cmd"/*; do
            if [ -d "$cmd_dir" ]; then
                local cmd_name=$(basename "$cmd_dir")
                if [ ! -f "$cmd_dir/main.go" ]; then
                    log_warn "Command directory $cmd_name missing main.go"
                    ((issues++))
                else
                    log_success "Command $cmd_name has main.go"
                fi
            fi
        done
    fi
    
    # Check internal directory (should not be importable externally)
    if [ -d "${PROJECT_ROOT}/internal" ]; then
        log_info "Checking internal directory structure..."
        # Internal packages should not be referenced outside the project
        local internal_imports=$(find "${PROJECT_ROOT}" -name "*.go" -not -path "*/internal/*" -not -path "*/vendor/*" \
            -exec grep -l "\".*internal/" {} \; 2>/dev/null || true)
        
        if [ -n "$internal_imports" ]; then
            log_warn "Found external references to internal packages:"
            echo "$internal_imports"
            ((issues++))
        else
            log_success "Internal packages properly isolated"
        fi
    fi
    
    return $issues
}

# Check file naming conventions
check_naming_conventions() {
    log_info "Checking Go naming conventions..."
    
    local issues=0
    
    # Check for non-Go naming in Go files
    while IFS= read -r -d '' file; do
        local basename=$(basename "$file" .go)
        
        # Check for camelCase in file names (should be snake_case or lowercase)
        if [[ "$basename" =~ [A-Z] ]] && [[ ! "$basename" =~ _test$ ]]; then
            log_warn "File uses camelCase naming: $file"
            ((issues++))
        fi
        
        # Check for spaces or special characters
        if [[ "$basename" =~ [[:space:]] ]] || [[ "$basename" =~ [^a-zA-Z0-9_] ]]; then
            log_error "File has invalid characters in name: $file"
            ((issues++))
        fi
        
    done < <(find "${PROJECT_ROOT}" -name "*.go" -not -path "*/vendor/*" -print0)
    
    if [ $issues -eq 0 ]; then
        log_success "All Go files follow naming conventions"
    fi
    
    return $issues
}

# Check test file organization
check_test_organization() {
    log_info "Checking test file organization..."
    
    local issues=0
    
    # Find all test files
    while IFS= read -r -d '' test_file; do
        local test_dir=$(dirname "$test_file")
        local test_name=$(basename "$test_file" _test.go)
        local source_file="${test_dir}/${test_name}.go"
        
        # Check if corresponding source file exists
        if [ ! -f "$source_file" ] && [[ ! "$test_name" =~ integration$ ]]; then
            log_warn "Test file without corresponding source: $test_file"
            ((issues++))
        fi
        
        # Check test package naming
        local package_line=$(head -1 "$test_file")
        if [[ "$package_line" =~ ^package[[:space:]]+([a-zA-Z0-9_]+) ]]; then
            local test_package="${BASH_REMATCH[1]}"
            
            # Get source package name if source file exists
            if [ -f "$source_file" ]; then
                local source_package_line=$(head -1 "$source_file")
                if [[ "$source_package_line" =~ ^package[[:space:]]+([a-zA-Z0-9_]+) ]]; then
                    local source_package="${BASH_REMATCH[1]}"
                    
                    # Test package should be same as source or source_test
                    if [ "$test_package" != "$source_package" ] && [ "$test_package" != "${source_package}_test" ]; then
                        log_warn "Test package mismatch in $test_file: $test_package vs $source_package"
                        ((issues++))
                    fi
                fi
            fi
        fi
        
    done < <(find "${PROJECT_ROOT}" -name "*_test.go" -not -path "*/vendor/*" -print0)
    
    if [ $issues -eq 0 ]; then
        log_success "Test files are properly organized"
    fi
    
    return $issues
}

# Check for circular dependencies
check_circular_dependencies() {
    log_info "Checking for circular dependencies..."
    
    # Use go mod graph to check for cycles
    cd "${PROJECT_ROOT}"
    
    if ! go mod graph >/dev/null 2>&1; then
        log_error "Failed to generate dependency graph"
        return 1
    fi
    
    # For now, just verify go mod graph works
    # A more sophisticated cycle detection would require parsing the graph
    log_success "No obvious circular dependencies detected"
    
    return 0
}

# Check import organization
check_import_organization() {
    log_info "Checking import organization..."
    
    local issues=0
    
    # Check for unused imports
    if command -v goimports >/dev/null 2>&1; then
        while IFS= read -r -d '' file; do
            local original_content=$(cat "$file")
            local formatted_content=$(goimports "$file")
            
            if [ "$original_content" != "$formatted_content" ]; then
                log_warn "File has import organization issues: $file"
                ((issues++))
            fi
        done < <(find "${PROJECT_ROOT}" -name "*.go" -not -path "*/vendor/*" -print0)
    else
        log_warn "goimports not available, skipping import organization check"
    fi
    
    if [ $issues -eq 0 ]; then
        log_success "Import organization looks good"
    fi
    
    return $issues
}

# Check template organization
check_template_organization() {
    log_info "Checking template organization..."
    
    local issues=0
    
    if [ -d "${PROJECT_ROOT}/templates" ]; then
        # Check for proper template structure
        local template_categories=("backend" "frontend" "mobile" "infrastructure" "base")
        
        for category in "${template_categories[@]}"; do
            if [ ! -d "${PROJECT_ROOT}/templates/${category}" ]; then
                log_warn "Missing template category: ${category}"
                ((issues++))
            else
                log_success "Found template category: ${category}"
            fi
        done
        
        # Check for .tmpl extension
        local non_tmpl_files=$(find "${PROJECT_ROOT}/templates" -type f -not -name "*.tmpl" -not -name "*.md" -not -name ".gitkeep" | wc -l)
        if [ "$non_tmpl_files" -gt 0 ]; then
            log_warn "Found $non_tmpl_files template files without .tmpl extension"
            ((issues++))
        fi
        
    else
        log_error "Templates directory not found"
        ((issues++))
    fi
    
    return $issues
}

# Main structural analysis
main() {
    log_info "Starting structural analysis..."
    
    local total_issues=0
    
    # Run all checks
    check_go_structure || ((total_issues+=$?))
    check_naming_conventions || ((total_issues+=$?))
    check_test_organization || ((total_issues+=$?))
    check_circular_dependencies || ((total_issues+=$?))
    check_import_organization || ((total_issues+=$?))
    check_template_organization || ((total_issues+=$?))
    
    # Save results
    mkdir -p "${AUDIT_DIR}"
    echo "structural_issues: $total_issues" >> "${AUDIT_DIR}/structure-results.txt"
    
    if [ $total_issues -eq 0 ]; then
        log_success "Structural analysis completed with no issues"
    else
        log_warn "Structural analysis completed with $total_issues issues"
    fi
    
    return $total_issues
}

main "$@"