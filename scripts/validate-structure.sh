#!/bin/bash

# Go Project Structure Validation Script
# This script validates that the Go project follows standard conventions

set -e

echo "=== Go Project Structure Validation ==="
echo "Timestamp: $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
ISSUES=0
WARNINGS=0

# Function to report issues
report_issue() {
    echo -e "${RED}[ERROR]${NC} $1"
    ((ISSUES++))
}

report_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
    ((WARNINGS++))
}

report_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

echo "1. Validating directory structure..."

# Check for required directories
required_dirs=("cmd" "internal" "pkg")
for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        report_success "Directory $dir exists"
    else
        report_issue "Required directory $dir is missing"
    fi
done

# Check cmd directory structure
echo
echo "2. Validating cmd/ directory..."
if [ -d "cmd" ]; then
    cmd_dirs=$(find cmd -mindepth 1 -maxdepth 1 -type d)
    if [ -z "$cmd_dirs" ]; then
        report_warning "cmd/ directory is empty"
    else
        for cmd_dir in $cmd_dirs; do
            cmd_name=$(basename "$cmd_dir")
            # Check if main.go exists
            if [ -f "$cmd_dir/main.go" ]; then
                report_success "Command $cmd_name has main.go"
            else
                report_issue "Command $cmd_name missing main.go"
            fi
            
            # Check naming convention (should be lowercase, no underscores)
            if [[ "$cmd_name" =~ ^[a-z][a-z0-9-]*$ ]]; then
                report_success "Command $cmd_name follows naming convention"
            else
                report_issue "Command $cmd_name violates naming convention (should be lowercase, hyphens allowed)"
            fi
        done
    fi
fi

# Check internal directory structure
echo
echo "3. Validating internal/ directory..."
if [ -d "internal" ]; then
    internal_dirs=$(find internal -mindepth 1 -maxdepth 1 -type d)
    if [ -z "$internal_dirs" ]; then
        report_warning "internal/ directory is empty"
    else
        for internal_dir in $internal_dirs; do
            dir_name=$(basename "$internal_dir")
            # Check naming convention
            if [[ "$dir_name" =~ ^[a-z][a-z0-9]*$ ]]; then
                report_success "Internal package $dir_name follows naming convention"
            else
                report_issue "Internal package $dir_name violates naming convention (should be lowercase, no underscores/hyphens)"
            fi
        done
    fi
fi

# Check pkg directory structure
echo
echo "4. Validating pkg/ directory..."
if [ -d "pkg" ]; then
    pkg_dirs=$(find pkg -mindepth 1 -maxdepth 1 -type d)
    if [ -z "$pkg_dirs" ]; then
        report_warning "pkg/ directory is empty"
    else
        for pkg_dir in $pkg_dirs; do
            dir_name=$(basename "$pkg_dir")
            # Check naming convention
            if [[ "$dir_name" =~ ^[a-z][a-z0-9]*$ ]]; then
                report_success "Public package $dir_name follows naming convention"
            else
                report_issue "Public package $dir_name violates naming convention (should be lowercase, no underscores/hyphens)"
            fi
        done
    fi
fi

echo
echo "5. Checking for misplaced files..."

# Check for Go files in root directory (should be minimal)
root_go_files=$(find . -maxdepth 1 -name "*.go" -type f)
if [ -n "$root_go_files" ]; then
    for file in $root_go_files; do
        report_warning "Go file in root directory: $file (consider moving to appropriate package)"
    done
else
    report_success "No Go files in root directory"
fi

# Check for test files co-location
echo
echo "6. Validating test file organization..."
test_files=$(find . -name "*_test.go" -type f)
for test_file in $test_files; do
    test_dir=$(dirname "$test_file")
    test_basename=$(basename "$test_file" _test.go)
    source_file="$test_dir/$test_basename.go"
    
    if [ -f "$source_file" ]; then
        report_success "Test file $test_file is co-located with source"
    else
        # Check if there are any .go files in the same directory
        go_files_in_dir=$(find "$test_dir" -maxdepth 1 -name "*.go" ! -name "*_test.go" -type f)
        if [ -z "$go_files_in_dir" ]; then
            report_warning "Test file $test_file has no corresponding Go files in same directory"
        fi
    fi
done

echo
echo "7. Checking package naming conventions..."

# Find all Go packages and check their names
go_files=$(find . -name "*.go" -type f ! -path "./vendor/*")
for go_file in $go_files; do
    # Extract package name from file
    package_name=$(head -n 10 "$go_file" | grep -E "^package " | head -n 1 | awk '{print $2}')
    if [ -n "$package_name" ]; then
        # Check if package name follows conventions
        if [[ "$package_name" =~ ^[a-z][a-z0-9]*$ ]] || [ "$package_name" = "main" ]; then
            # Package name is valid
            continue
        else
            report_issue "Package name '$package_name' in $go_file violates naming convention"
        fi
    fi
done

echo
echo "=== Validation Summary ==="
echo "Issues found: $ISSUES"
echo "Warnings: $WARNINGS"

if [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✓ Project structure validation passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Project structure validation failed with $ISSUES issues${NC}"
    exit 1
fi