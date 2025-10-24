#!/bin/bash

# CI Test Script
# This script runs the comprehensive test suite for CI/CD pipelines
# with coverage reporting, race detection, and linting

set -e

# Configuration
TEST_TIMEOUT=${TEST_TIMEOUT:-"10m"}
COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD:-"0"}  # Minimum coverage percentage (0 = no minimum)
PARALLEL_TESTS=${PARALLEL_TESTS:-"4"}
REPORTS_DIR="test-reports"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Cleanup function
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        print_error "Tests failed with exit code $exit_code"
        if [ -f "${REPORTS_DIR}/failed-tests.txt" ]; then
            echo ""
            print_error "Failed tests:"
            cat "${REPORTS_DIR}/failed-tests.txt"
        fi
    fi
}
trap cleanup EXIT

# Print environment information
print_env_info() {
    echo "========================================="
    echo "CI Test Environment"
    echo "========================================="
    print_status "Go version: $(go version)"
    print_status "OS: $(uname -s)"
    print_status "Architecture: $(uname -m)"
    print_status "Working directory: $(pwd)"
    print_status "Test timeout: ${TEST_TIMEOUT}"
    print_status "Parallel tests: ${PARALLEL_TESTS}"
    
    if [ -n "${CI}" ]; then
        print_status "Running in CI environment"
    fi
    
    echo ""
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check Go installation
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check if we're in a Go module
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. This script must be run from the project root."
        exit 1
    fi
    
    # Create reports directory
    mkdir -p "${REPORTS_DIR}"
    
    print_success "Prerequisites check passed"
    echo ""
}

# Download dependencies
download_dependencies() {
    print_status "Downloading dependencies..."
    
    if go mod download; then
        print_success "Dependencies downloaded"
    else
        print_error "Failed to download dependencies"
        exit 1
    fi
    
    echo ""
}

# Verify code compiles
verify_build() {
    print_status "Verifying code compiles..."
    
    if go build -v ./...; then
        print_success "Code compiles successfully"
    else
        print_error "Code compilation failed"
        exit 1
    fi
    
    echo ""
}

# Run linting
run_lint() {
    print_status "Running linters..."
    
    # Check if golangci-lint is available
    if command -v golangci-lint >/dev/null 2>&1; then
        if golangci-lint run --timeout=5m ./...; then
            print_success "Linting passed"
        else
            print_warning "Linting found issues (non-blocking)"
            # Don't fail on lint errors in CI, just warn
        fi
    else
        print_warning "golangci-lint not found, skipping linting"
        print_status "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
    fi
    
    echo ""
}

# Run tests with coverage
run_tests() {
    print_status "Running test suite with coverage..."
    
    local coverage_file="${REPORTS_DIR}/coverage.out"
    local coverage_html="${REPORTS_DIR}/coverage.html"
    local test_output="${REPORTS_DIR}/test-output.txt"
    
    # Run tests with coverage and race detection
    if go test -v -race -timeout="${TEST_TIMEOUT}" -parallel="${PARALLEL_TESTS}" \
        -coverprofile="${coverage_file}" -covermode=atomic \
        ./... 2>&1 | tee "${test_output}"; then
        print_success "All tests passed"
    else
        # Extract failed tests
        grep -E "^--- FAIL:" "${test_output}" > "${REPORTS_DIR}/failed-tests.txt" 2>/dev/null || true
        print_error "Some tests failed"
        return 1
    fi
    
    echo ""
}

# Generate coverage report
generate_coverage_report() {
    print_status "Generating coverage report..."
    
    local coverage_file="${REPORTS_DIR}/coverage.out"
    local coverage_html="${REPORTS_DIR}/coverage.html"
    
    if [ ! -f "${coverage_file}" ]; then
        print_warning "Coverage file not found, skipping coverage report"
        return
    fi
    
    # Generate HTML coverage report
    go tool cover -html="${coverage_file}" -o="${coverage_html}"
    print_success "Coverage HTML report: ${coverage_html}"
    
    # Calculate coverage percentage
    local coverage_percent=$(go tool cover -func="${coverage_file}" | grep total | awk '{print $3}' | sed 's/%//')
    
    echo ""
    print_status "========================================="
    print_status "Coverage Summary"
    print_status "========================================="
    go tool cover -func="${coverage_file}" | tail -10
    echo ""
    print_status "Total coverage: ${coverage_percent}%"
    
    # Check coverage threshold
    if [ -n "${COVERAGE_THRESHOLD}" ] && [ "${COVERAGE_THRESHOLD}" != "0" ]; then
        if (( $(echo "${coverage_percent} < ${COVERAGE_THRESHOLD}" | bc -l) )); then
            print_error "Coverage ${coverage_percent}% is below threshold ${COVERAGE_THRESHOLD}%"
            return 1
        else
            print_success "Coverage meets threshold (${COVERAGE_THRESHOLD}%)"
        fi
    fi
    
    echo ""
}

# Run benchmarks (optional)
run_benchmarks() {
    if [ "${RUN_BENCHMARKS}" = "true" ]; then
        print_status "Running benchmarks..."
        
        local bench_output="${REPORTS_DIR}/benchmarks.txt"
        
        if go test -bench=. -benchmem -benchtime=1s -run=^$ ./... > "${bench_output}" 2>&1; then
            print_success "Benchmarks completed"
            echo ""
            print_status "Benchmark results:"
            cat "${bench_output}"
        else
            print_warning "Benchmarks failed (non-blocking)"
        fi
        
        echo ""
    fi
}

# Main function
main() {
    echo "========================================="
    echo "🧪 CI Test Suite"
    echo "========================================="
    echo ""
    
    print_env_info
    check_prerequisites
    download_dependencies
    verify_build
    run_lint
    run_tests
    generate_coverage_report
    run_benchmarks
    
    echo ""
    echo "========================================="
    print_success "✅ All CI checks passed!"
    echo "========================================="
    echo ""
    print_status "Test reports available in: ${REPORTS_DIR}/"
    ls -lh "${REPORTS_DIR}/" 2>/dev/null || true
    echo ""
    print_status "Summary:"
    echo "  ✓ Code compilation successful"
    echo "  ✓ All tests passed"
    echo "  ✓ Race detection passed"
    echo "  ✓ Coverage report generated"
    echo ""
    print_status "The same command runs locally and in CI:"
    echo "  go test -v -race -coverprofile=coverage.out ./..."
    echo ""
}

# Show usage
show_usage() {
    cat << EOF
CI Test Suite

Usage: $0 [OPTIONS]

Description:
  Runs comprehensive test suite for CI/CD pipelines with coverage reporting,
  race detection, and linting.

Options:
  -h, --help      Show this help message

Environment Variables:
  TEST_TIMEOUT         Test timeout duration (default: 10m)
  COVERAGE_THRESHOLD   Minimum coverage percentage (default: 0)
  PARALLEL_TESTS       Number of parallel test processes (default: 4)
  RUN_BENCHMARKS       Set to 'true' to run benchmarks (default: false)

Examples:
  $0                                # Run full test suite
  TEST_TIMEOUT=15m $0               # Run with 15 minute timeout
  COVERAGE_THRESHOLD=80 $0          # Require 80% coverage
  RUN_BENCHMARKS=true $0            # Include benchmarks

Requirements:
  - Go 1.25+
  - golangci-lint (optional, for linting)

Output:
  Test reports: test-reports/
  Coverage: test-reports/coverage.html
EOF
}

# Handle command line arguments
case "${1:-}" in
    "help"|"-h"|"--help")
        show_usage
        exit 0
        ;;
    *)
        main
        ;;
esac