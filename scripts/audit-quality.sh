#!/bin/bash

# Code Quality Analysis Script
# Runs linting tools and checks code quality

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
    echo -e "${BLUE}[QUALITY]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[QUALITY-WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[QUALITY-ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[QUALITY-OK]${NC} $1"
}

# Run gofmt check
run_gofmt() {
    log_info "Running gofmt check..."
    
    cd "${PROJECT_ROOT}"
    
    local issues=0
    local unformatted_files=()
    
    # Check all Go files for formatting
    while IFS= read -r -d '' file; do
        if ! gofmt -l "$file" | grep -q .; then
            continue
        else
            unformatted_files+=("$file")
            ((issues++))
        fi
    done < <(find . -name "*.go" -not -path "./vendor/*" -print0)
    
    if [ ${#unformatted_files[@]} -gt 0 ]; then
        log_warn "Found ${#unformatted_files[@]} unformatted Go files:"
        printf '%s\n' "${unformatted_files[@]}" | head -10
        if [ ${#unformatted_files[@]} -gt 10 ]; then
            log_info "... and $((${#unformatted_files[@]} - 10)) more files"
        fi
        
        # Save to file
        printf '%s\n' "${unformatted_files[@]}" > "${AUDIT_DIR}/unformatted-files.txt"
    else
        log_success "All Go files are properly formatted"
    fi
    
    return $issues
}

# Run go vet
run_go_vet() {
    log_info "Running go vet..."
    
    cd "${PROJECT_ROOT}"
    
    local vet_output="${AUDIT_DIR}/go-vet-output.txt"
    
    if go vet ./... 2> "$vet_output"; then
        log_success "go vet passed with no issues"
        return 0
    else
        log_error "go vet found issues:"
        cat "$vet_output"
        return 1
    fi
}

# Run golangci-lint
run_golangci_lint() {
    log_info "Running golangci-lint..."
    
    cd "${PROJECT_ROOT}"
    
    # Check if golangci-lint is available
    if ! command -v golangci-lint >/dev/null 2>&1; then
        log_error "golangci-lint not found"
        return 1
    fi
    
    local lint_output="${AUDIT_DIR}/golangci-lint-output.txt"
    
    # Run golangci-lint with configuration
    if golangci-lint run --out-format=tab ./... > "$lint_output" 2>&1; then
        log_success "golangci-lint passed with no issues"
        return 0
    else
        local issue_count=$(wc -l < "$lint_output" 2>/dev/null || echo "0")
        log_warn "golangci-lint found $issue_count issues:"
        head -20 "$lint_output"
        if [ "$issue_count" -gt 20 ]; then
            log_info "... and $((issue_count - 20)) more issues (see $lint_output)"
        fi
        return 1
    fi
}

# Check for unused code
check_unused_code() {
    log_info "Checking for unused code..."
    
    cd "${PROJECT_ROOT}"
    
    local issues=0
    
    # Install and run go-unused if available
    if command -v unused >/dev/null 2>&1; then
        local unused_output="${AUDIT_DIR}/unused-code.txt"
        
        if unused ./... > "$unused_output" 2>&1; then
            if [ -s "$unused_output" ]; then
                local unused_count=$(wc -l < "$unused_output")
                log_warn "Found $unused_count unused code elements:"
                head -10 "$unused_output"
                if [ "$unused_count" -gt 10 ]; then
                    log_info "... and $((unused_count - 10)) more (see $unused_output)"
                fi
                ((issues++))
            else
                log_success "No unused code found"
            fi
        else
            log_warn "unused tool failed to run"
        fi
    else
        log_warn "unused tool not available, installing..."
        if go install honnef.co/go/tools/cmd/unused@latest; then
            # Retry after installation
            check_unused_code
            return $?
        else
            log_warn "Could not install unused tool, skipping unused code check"
        fi
    fi
    
    return $issues
}

# Check code complexity
check_code_complexity() {
    log_info "Checking code complexity..."
    
    cd "${PROJECT_ROOT}"
    
    local issues=0
    
    # Use gocyclo if available
    if command -v gocyclo >/dev/null 2>&1; then
        local complexity_output="${AUDIT_DIR}/complexity.txt"
        
        # Check for functions with cyclomatic complexity > 10
        if gocyclo -over 10 . > "$complexity_output" 2>&1; then
            if [ -s "$complexity_output" ]; then
                local complex_count=$(wc -l < "$complexity_output")
                log_warn "Found $complex_count functions with high complexity (>10):"
                head -5 "$complexity_output"
                ((issues++))
            else
                log_success "No overly complex functions found"
            fi
        else
            log_warn "gocyclo failed to run"
        fi
    else
        log_warn "gocyclo not available, installing..."
        if go install github.com/fzipp/gocyclo/cmd/gocyclo@latest; then
            # Retry after installation
            check_code_complexity
            return $?
        else
            log_warn "Could not install gocyclo, skipping complexity check"
        fi
    fi
    
    return $issues
}

# Check for TODO/FIXME comments
check_todo_comments() {
    log_info "Checking for TODO/FIXME comments..."
    
    cd "${PROJECT_ROOT}"
    
    local todo_output="${AUDIT_DIR}/todo-comments.txt"
    
    # Find TODO, FIXME, HACK comments
    grep -rn --include="*.go" -E "(TODO|FIXME|HACK|XXX)" . > "$todo_output" 2>/dev/null || true
    
    if [ -s "$todo_output" ]; then
        local todo_count=$(wc -l < "$todo_output")
        log_info "Found $todo_count TODO/FIXME comments:"
        head -10 "$todo_output"
        if [ "$todo_count" -gt 10 ]; then
            log_info "... and $((todo_count - 10)) more (see $todo_output)"
        fi
    else
        log_success "No TODO/FIXME comments found"
    fi
    
    return 0
}

# Check test coverage
check_test_coverage() {
    log_info "Checking test coverage..."
    
    cd "${PROJECT_ROOT}"
    
    local coverage_output="${AUDIT_DIR}/coverage.out"
    local coverage_report="${AUDIT_DIR}/coverage-report.txt"
    
    # Run tests with coverage
    if go test -coverprofile="$coverage_output" ./...; then
        # Generate coverage report
        go tool cover -func="$coverage_output" > "$coverage_report"
        
        # Extract total coverage
        local total_coverage=$(tail -1 "$coverage_report" | awk '{print $3}')
        log_info "Total test coverage: $total_coverage"
        
        # Check if coverage is acceptable (>= 70%)
        local coverage_percent=$(echo "$total_coverage" | sed 's/%//')
        if (( $(echo "$coverage_percent >= 70" | bc -l) )); then
            log_success "Test coverage is acceptable: $total_coverage"
        else
            log_warn "Test coverage is low: $total_coverage (target: >= 70%)"
            return 1
        fi
    else
        log_error "Tests failed, cannot generate coverage report"
        return 1
    fi
    
    return 0
}

# Generate quality report
generate_quality_report() {
    log_info "Generating code quality report..."
    
    mkdir -p "${AUDIT_DIR}"
    
    cat > "${AUDIT_DIR}/quality-report.md" << EOF
# Code Quality Analysis Report

Generated on: $(date)

## Summary

### Files Analyzed
- Go source files: $(find "${PROJECT_ROOT}" -name "*.go" -not -path "*/vendor/*" | wc -l)
- Test files: $(find "${PROJECT_ROOT}" -name "*_test.go" | wc -l)

### Tools Used
- gofmt: Code formatting
- go vet: Static analysis
- golangci-lint: Comprehensive linting
- unused: Unused code detection
- gocyclo: Complexity analysis

### Results

#### Formatting Issues
$(if [ -f "${AUDIT_DIR}/unformatted-files.txt" ]; then
    echo "Unformatted files: $(wc -l < "${AUDIT_DIR}/unformatted-files.txt")"
else
    echo "No formatting issues found"
fi)

#### Linting Issues
$(if [ -f "${AUDIT_DIR}/golangci-lint-output.txt" ]; then
    echo "Linting issues: $(wc -l < "${AUDIT_DIR}/golangci-lint-output.txt")"
else
    echo "No linting issues found"
fi)

#### Test Coverage
$(if [ -f "${AUDIT_DIR}/coverage-report.txt" ]; then
    tail -1 "${AUDIT_DIR}/coverage-report.txt"
else
    echo "Coverage report not available"
fi)

#### TODO Comments
$(if [ -f "${AUDIT_DIR}/todo-comments.txt" ]; then
    echo "TODO/FIXME comments: $(wc -l < "${AUDIT_DIR}/todo-comments.txt")"
else
    echo "No TODO/FIXME comments found"
fi)

## Recommendations

1. Fix all formatting issues with: \`gofmt -w .\`
2. Address linting issues reported by golangci-lint
3. Remove unused code elements
4. Improve test coverage where needed
5. Address TODO/FIXME comments

EOF

    log_success "Quality report generated at ${AUDIT_DIR}/quality-report.md"
}

# Main quality analysis
main() {
    log_info "Starting code quality analysis..."
    
    local total_issues=0
    
    # Create audit directory
    mkdir -p "${AUDIT_DIR}"
    
    # Run all checks
    run_gofmt || ((total_issues+=$?))
    run_go_vet || ((total_issues+=$?))
    run_golangci_lint || ((total_issues+=$?))
    check_unused_code || ((total_issues+=$?))
    check_code_complexity || ((total_issues+=$?))
    check_todo_comments || ((total_issues+=$?))
    check_test_coverage || ((total_issues+=$?))
    
    # Generate report
    generate_quality_report
    
    # Save results
    echo "quality_issues: $total_issues" >> "${AUDIT_DIR}/quality-results.txt"
    
    if [ $total_issues -eq 0 ]; then
        log_success "Code quality analysis completed with no issues"
    else
        log_warn "Code quality analysis completed with $total_issues issues"
    fi
    
    return $total_issues
}

main "$@"