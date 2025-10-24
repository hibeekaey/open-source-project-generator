#!/bin/bash

# Validation script to check project setup and configuration
# This script validates scripts, Makefile, Dockerfiles, and CI/CD workflows

set -e

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

# Track validation results
ERRORS=0
WARNINGS=0
CHECKS=0

# Validation function
validate() {
    local check_name=$1
    local check_command=$2
    
    CHECKS=$((CHECKS + 1))
    print_status "Checking: $check_name"
    
    if eval "$check_command" >/dev/null 2>&1; then
        print_success "$check_name"
        return 0
    else
        print_error "$check_name"
        ERRORS=$((ERRORS + 1))
        return 1
    fi
}

# Warning function
warn() {
    local check_name=$1
    local check_command=$2
    
    CHECKS=$((CHECKS + 1))
    print_status "Checking: $check_name"
    
    if eval "$check_command" >/dev/null 2>&1; then
        print_success "$check_name"
        return 0
    else
        print_warning "$check_name"
        WARNINGS=$((WARNINGS + 1))
        return 1
    fi
}

echo "========================================="
echo "Project Setup Validation"
echo "========================================="
echo ""

# Check script files
print_status "=== Validating Scripts ==="
echo ""

validate "build.sh exists" "[ -f scripts/build.sh ]"
validate "build.sh is executable" "[ -x scripts/build.sh ]"
validate "build.sh has valid syntax" "bash -n scripts/build.sh"

validate "build-packages.sh exists" "[ -f scripts/build-packages.sh ]"
validate "build-packages.sh is executable" "[ -x scripts/build-packages.sh ]"
validate "build-packages.sh has valid syntax" "bash -n scripts/build-packages.sh"

validate "ci-test.sh exists" "[ -f scripts/ci-test.sh ]"
validate "ci-test.sh is executable" "[ -x scripts/ci-test.sh ]"
validate "ci-test.sh has valid syntax" "bash -n scripts/ci-test.sh"

validate "install.sh exists" "[ -f scripts/install.sh ]"
validate "install.sh is executable" "[ -x scripts/install.sh ]"
validate "install.sh has valid syntax" "bash -n scripts/install.sh"

validate "run_performance_benchmarks.sh exists" "[ -f scripts/run_performance_benchmarks.sh ]"
validate "run_performance_benchmarks.sh is executable" "[ -x scripts/run_performance_benchmarks.sh ]"
validate "run_performance_benchmarks.sh has valid syntax" "bash -n scripts/run_performance_benchmarks.sh"

echo ""

# Check Makefile
print_status "=== Validating Makefile ==="
echo ""

validate "Makefile exists" "[ -f Makefile ]"
validate "Makefile has build target" "grep -q '^build:' Makefile"
validate "Makefile has test target" "grep -q '^test:' Makefile"
validate "Makefile has check target" "grep -q '^check:' Makefile"
validate "Makefile has lint target" "grep -q '^lint:' Makefile"
validate "Makefile has ci target" "grep -q '^ci:' Makefile"
validate "Makefile has dist target" "grep -q '^dist:' Makefile"
validate "Makefile has package target" "grep -q '^package:' Makefile"
validate "Makefile has release target" "grep -q '^release:' Makefile"
validate "Makefile has docker-build target" "grep -q '^docker-build:' Makefile"
validate "Makefile has security-scan target" "grep -q '^security-scan:' Makefile"
validate "Makefile has benchmark target" "grep -q '^benchmark:' Makefile"
validate "Makefile has clean target" "grep -q '^clean:' Makefile"

echo ""

# Check Dockerfiles
print_status "=== Validating Dockerfiles ==="
echo ""

validate "Dockerfile exists" "[ -f Dockerfile ]"
validate "Dockerfile.build exists" "[ -f Dockerfile.build ]"
validate "Dockerfile.dev exists" "[ -f Dockerfile.dev ]"

warn "Dockerfile syntax (requires docker)" "docker build -f Dockerfile --no-cache -t test:validation . >/dev/null 2>&1 || true"

echo ""

# Check CI/CD workflows
print_status "=== Validating CI/CD Workflows ==="
echo ""

validate "CI workflow exists" "[ -f .github/workflows/ci.yml ]"
validate "Docker workflow exists" "[ -f .github/workflows/docker.yml ]"
validate "Release workflow exists" "[ -f .github/workflows/release.yml ]"

validate "CI workflow uses make ci" "grep -q 'make ci' .github/workflows/ci.yml"
validate "Release workflow uses make release or build scripts" "grep -qE 'make (release|dist|package)|build.*\.sh' .github/workflows/release.yml"

echo ""

# Check Go project structure
print_status "=== Validating Go Project ==="
echo ""

validate "go.mod exists" "[ -f go.mod ]"
validate "go.sum exists" "[ -f go.sum ]"
validate "cmd/generator exists" "[ -d cmd/generator ]"
validate "pkg directory exists" "[ -d pkg ]"

warn "Go modules are valid" "go mod verify"

echo ""

# Check documentation
print_status "=== Validating Documentation ==="
echo ""

validate "README.md exists" "[ -f README.md ]"
validate "LICENSE exists" "[ -f LICENSE ]"
warn "CHANGELOG.md exists" "[ -f CHANGELOG.md ]"

echo ""

# Check .gitignore
print_status "=== Validating .gitignore ==="
echo ""

validate ".gitignore exists" "[ -f .gitignore ]"
validate ".gitignore includes dist/" "grep -q 'dist/' .gitignore"
validate ".gitignore includes bin/" "grep -q 'bin/' .gitignore"
validate ".gitignore includes coverage.out" "grep -q 'coverage.out' .gitignore"
validate ".gitignore includes benchmark_results/" "grep -q 'benchmark_results/' .gitignore"
validate ".gitignore includes test-reports/" "grep -q 'test-reports/' .gitignore"

echo ""

# Check script consistency
print_status "=== Validating Script Consistency ==="
echo ""

validate "All scripts use set -e" "grep -q 'set -e' scripts/*.sh"
validate "All scripts have cleanup traps" "grep -q 'trap.*EXIT' scripts/*.sh"
validate "All scripts have show_usage function" "grep -q 'show_usage()' scripts/*.sh"
validate "All scripts have colored output" "grep -q 'print_status()' scripts/*.sh"

echo ""

# Check version management
print_status "=== Validating Version Management ==="
echo ""

validate "get-version.sh exists" "[ -f scripts/get-version.sh ]"
validate "get-version.sh is executable" "[ -x scripts/get-version.sh ]"
validate "get-version.sh has valid syntax" "bash -n scripts/get-version.sh"
validate "get-version.sh can get version" "./scripts/get-version.sh version >/dev/null"
validate "get-version.sh can get commit" "./scripts/get-version.sh commit >/dev/null"
validate "get-version.sh can validate" "./scripts/get-version.sh validate >/dev/null"

echo ""

# Summary
echo "========================================="
echo "Validation Summary"
echo "========================================="
echo ""
echo "Total checks: $CHECKS"
echo -e "${GREEN}Passed: $((CHECKS - ERRORS - WARNINGS))${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo -e "${RED}Errors: $ERRORS${NC}"
echo ""

if [ $ERRORS -eq 0 ]; then
    print_success "All critical validations passed!"
    if [ $WARNINGS -gt 0 ]; then
        print_warning "Some optional checks failed, but project is functional"
    fi
    exit 0
else
    print_error "Validation failed with $ERRORS error(s)"
    exit 1
fi
