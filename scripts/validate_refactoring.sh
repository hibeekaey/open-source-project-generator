#!/bin/bash

# Refactoring Validation Script
# This script validates that the refactoring maintains all functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Track validation results
VALIDATION_ERRORS=0
VALIDATION_WARNINGS=0

# Function to record validation failure
validation_error() {
    log_error "$1"
    VALIDATION_ERRORS=$((VALIDATION_ERRORS + 1))
}

# Function to record validation warning
validation_warning() {
    log_warning "$1"
    VALIDATION_WARNINGS=$((VALIDATION_WARNINGS + 1))
}

# Check if we're in the right directory
check_project_root() {
    if [[ ! -f "go.mod" ]] || [[ ! -d "pkg" ]]; then
        validation_error "This script must be run from the project root directory"
        exit 1
    fi
}

# Validate Go code compilation
validate_compilation() {
    log_info "Validating Go code compilation..."
    
    if ! go build ./...; then
        validation_error "Go code compilation failed"
        return 1
    fi
    
    log_success "Go code compiles successfully"
}

# Validate all tests pass
validate_tests() {
    log_info "Running all tests..."
    
    if ! go test ./... -v; then
        validation_error "Some tests are failing"
        return 1
    fi
    
    log_success "All tests pass"
}

# Validate test coverage
validate_test_coverage() {
    log_info "Checking test coverage..."
    
    # Generate coverage report
    go test ./... -coverprofile=coverage.out -covermode=atomic
    
    # Get coverage percentage
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    
    if (( $(echo "$coverage < 70" | bc -l) )); then
        validation_warning "Test coverage is below 70% ($coverage%)"
    else
        log_success "Test coverage is adequate ($coverage%)"
    fi
    
    # Clean up
    rm -f coverage.out
}

# Validate linting
validate_linting() {
    log_info "Running linting checks..."
    
    # Check if golangci-lint is available
    if ! command -v golangci-lint &> /dev/null; then
        validation_warning "golangci-lint not found, skipping lint checks"
        return 0
    fi
    
    if ! golangci-lint run; then
        validation_error "Linting issues found"
        return 1
    fi
    
    log_success "Linting checks pass"
}

# Validate code formatting
validate_formatting() {
    log_info "Checking code formatting..."
    
    # Check gofmt
    unformatted=$(gofmt -l .)
    if [[ -n "$unformatted" ]]; then
        validation_error "Code formatting issues found in: $unformatted"
        return 1
    fi
    
    # Check goimports if available
    if command -v goimports &> /dev/null; then
        unformatted_imports=$(goimports -l .)
        if [[ -n "$unformatted_imports" ]]; then
            validation_warning "Import formatting issues found in: $unformatted_imports"
        fi
    fi
    
    log_success "Code formatting is correct"
}

# Validate go vet
validate_vet() {
    log_info "Running go vet..."
    
    if ! go vet ./...; then
        validation_error "go vet found issues"
        return 1
    fi
    
    log_success "go vet checks pass"
}

# Validate module dependencies
validate_dependencies() {
    log_info "Validating module dependencies..."
    
    # Check if go mod tidy would make changes
    cp go.mod go.mod.backup
    cp go.sum go.sum.backup
    
    go mod tidy
    
    if ! diff -q go.mod go.mod.backup > /dev/null || ! diff -q go.sum go.sum.backup > /dev/null; then
        validation_warning "go mod tidy made changes to dependencies"
        # Restore original files
        mv go.mod.backup go.mod
        mv go.sum.backup go.sum
    else
        log_success "Module dependencies are clean"
        rm -f go.mod.backup go.sum.backup
    fi
}

# Validate CLI functionality
validate_cli_functionality() {
    log_info "Validating CLI functionality..."
    
    # Build the CLI binary
    if ! go build -o test_generator ./cmd/generator; then
        validation_error "Failed to build CLI binary"
        return 1
    fi
    
    # Test basic CLI commands
    commands=(
        "--help"
        "version"
        "list-templates"
        "validate --help"
        "audit --help"
        "generate --help"
    )
    
    for cmd in "${commands[@]}"; do
        log_info "Testing command: generator $cmd"
        if ! ./test_generator $cmd > /dev/null 2>&1; then
            validation_error "CLI command failed: generator $cmd"
        fi
    done
    
    # Clean up
    rm -f test_generator
    
    log_success "CLI functionality validated"
}

# Validate file structure
validate_file_structure() {
    log_info "Validating file structure..."
    
    # Check that no single file is too large
    large_files=$(find pkg/ -name "*.go" -exec wc -l {} + | awk '$1 > 1000 {print $2 " (" $1 " lines)"}')
    
    if [[ -n "$large_files" ]]; then
        validation_warning "Large files found (>1000 lines):"
        echo "$large_files"
    fi
    
    # Check for proper package organization
    expected_dirs=(
        "pkg/cli"
        "pkg/audit"
        "pkg/template"
        "pkg/validation"
        "pkg/cache"
        "pkg/filesystem"
    )
    
    for dir in "${expected_dirs[@]}"; do
        if [[ ! -d "$dir" ]]; then
            validation_error "Expected directory not found: $dir"
        fi
    done
    
    log_success "File structure validation completed"
}

# Validate no circular imports
validate_imports() {
    log_info "Checking for circular imports..."
    
    # This is a basic check - Go compiler would catch circular imports anyway
    if ! go list ./...; then
        validation_error "Import issues detected"
        return 1
    fi
    
    log_success "No circular import issues found"
}

# Validate performance (basic check)
validate_performance() {
    log_info "Running basic performance validation..."
    
    # Build and time a simple operation
    if ! go build -o test_generator ./cmd/generator; then
        validation_error "Failed to build for performance test"
        return 1
    fi
    
    # Time the help command (should be fast)
    start_time=$(date +%s%N)
    ./test_generator --help > /dev/null 2>&1
    end_time=$(date +%s%N)
    
    duration=$(( (end_time - start_time) / 1000000 )) # Convert to milliseconds
    
    if (( duration > 1000 )); then
        validation_warning "CLI startup time is slow: ${duration}ms"
    else
        log_success "CLI startup time is acceptable: ${duration}ms"
    fi
    
    # Clean up
    rm -f test_generator
}

# Generate validation report
generate_report() {
    log_info "Generating validation report..."
    
    cat > validation_report.md << EOF
# Refactoring Validation Report

Generated on: $(date)

## Summary
- **Validation Errors**: $VALIDATION_ERRORS
- **Validation Warnings**: $VALIDATION_WARNINGS

## Test Results
EOF

    if [[ $VALIDATION_ERRORS -eq 0 ]]; then
        cat >> validation_report.md << EOF
✅ **All critical validations passed**

The refactoring has been successful and maintains all functionality.
EOF
    else
        cat >> validation_report.md << EOF
❌ **$VALIDATION_ERRORS critical validation(s) failed**

Please review and fix the issues before proceeding.
EOF
    fi

    if [[ $VALIDATION_WARNINGS -gt 0 ]]; then
        cat >> validation_report.md << EOF

⚠️ **$VALIDATION_WARNINGS warning(s) found**

These should be reviewed but don't block the refactoring.
EOF
    fi

    cat >> validation_report.md << EOF

## Validation Steps Performed
- [x] Go code compilation
- [x] All tests execution
- [x] Test coverage analysis
- [x] Linting checks
- [x] Code formatting validation
- [x] Go vet analysis
- [x] Module dependencies check
- [x] CLI functionality testing
- [x] File structure validation
- [x] Import validation
- [x] Basic performance check

## Next Steps
EOF

    if [[ $VALIDATION_ERRORS -eq 0 ]]; then
        cat >> validation_report.md << EOF
1. Review any warnings and address if necessary
2. Update documentation if needed
3. Commit the refactored code
4. Deploy and monitor for any issues
EOF
    else
        cat >> validation_report.md << EOF
1. **Fix all validation errors before proceeding**
2. Re-run this validation script
3. Address any remaining warnings
4. Update documentation
EOF
    fi

    log_success "Validation report generated: validation_report.md"
}

# Main validation function
run_validation() {
    log_info "Starting comprehensive refactoring validation..."
    
    check_project_root
    
    # Run all validation steps
    validate_compilation
    validate_tests
    validate_test_coverage
    validate_linting
    validate_formatting
    validate_vet
    validate_dependencies
    validate_cli_functionality
    validate_file_structure
    validate_imports
    validate_performance
    
    generate_report
    
    # Final summary
    echo ""
    echo "=== VALIDATION SUMMARY ==="
    if [[ $VALIDATION_ERRORS -eq 0 ]]; then
        log_success "All critical validations passed! ✅"
        if [[ $VALIDATION_WARNINGS -gt 0 ]]; then
            log_warning "$VALIDATION_WARNINGS warnings found - review recommended"
        fi
        echo ""
        echo "The refactoring has been successful and maintains all functionality."
        exit 0
    else
        log_error "$VALIDATION_ERRORS critical validation(s) failed! ❌"
        if [[ $VALIDATION_WARNINGS -gt 0 ]]; then
            log_warning "$VALIDATION_WARNINGS warnings also found"
        fi
        echo ""
        echo "Please fix the issues and re-run validation."
        exit 1
    fi
}

# Run validation based on command line argument
case "${1:-all}" in
    "compile")
        check_project_root
        validate_compilation
        ;;
    "test")
        check_project_root
        validate_tests
        ;;
    "lint")
        check_project_root
        validate_linting
        ;;
    "format")
        check_project_root
        validate_formatting
        ;;
    "cli")
        check_project_root
        validate_cli_functionality
        ;;
    "all")
        run_validation
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [validation_type]"
        echo "Validation types:"
        echo "  compile  - Check Go compilation"
        echo "  test     - Run all tests"
        echo "  lint     - Run linting checks"
        echo "  format   - Check code formatting"
        echo "  cli      - Test CLI functionality"
        echo "  all      - Run all validations (default)"
        echo "  help     - Show this help"
        ;;
    *)
        log_error "Unknown validation type: $1"
        log_info "Use '$0 help' for usage information"
        exit 1
        ;;
esac