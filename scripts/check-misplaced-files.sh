#!/bin/bash

# Script to check for misplaced files and directories
set -e

echo "=== Checking for Misplaced Files and Directories ==="
echo "Timestamp: $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ISSUES=0
WARNINGS=0

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

echo "1. Checking test file co-location..."

# Check if test files are co-located with their source files
test_files=$(find . -name "*_test.go" -type f)
for test_file in $test_files; do
    test_dir=$(dirname "$test_file")
    test_basename=$(basename "$test_file" _test.go)
    source_file="$test_dir/$test_basename.go"
    
    if [ -f "$source_file" ]; then
        report_success "Test file $test_file is co-located with source"
    else
        # Check if it's an integration test (acceptable to not have direct source file)
        if [[ "$test_file" == *"integration"* ]] || [[ "$test_file" == *"e2e"* ]]; then
            report_success "Integration/E2E test file $test_file (no direct source file needed)"
        else
            # Check if there are any .go files in the same directory
            go_files_in_dir=$(find "$test_dir" -maxdepth 1 -name "*.go" ! -name "*_test.go" -type f | wc -l)
            if [ "$go_files_in_dir" -eq 0 ]; then
                report_warning "Test file $test_file has no corresponding Go files in same directory"
            else
                report_success "Test file $test_file is in directory with other Go files"
            fi
        fi
    fi
done

echo
echo "2. Checking for files that should be in internal/..."

# Check for implementation details in pkg/ that should be internal
echo "Checking pkg/ packages for potential internal candidates..."

# Look for packages that might be implementation details
pkg_dirs=$(find pkg -mindepth 1 -maxdepth 1 -type d)
for pkg_dir in $pkg_dirs; do
    pkg_name=$(basename "$pkg_dir")
    
    # Simple check: count Go files in the package
    go_files_count=$(find "$pkg_dir" -name "*.go" ! -name "*_test.go" -type f | wc -l)
    test_files_count=$(find "$pkg_dir" -name "*_test.go" -type f | wc -l)
    
    if [ "$go_files_count" -gt 0 ]; then
        report_success "Package $pkg_name has $go_files_count Go files"
    elif [ "$test_files_count" -gt 0 ]; then
        # Integration/test-only packages are acceptable
        if [[ "$pkg_name" == *"integration"* ]] || [[ "$pkg_name" == *"test"* ]]; then
            report_success "Package $pkg_name is test-only package with $test_files_count test files"
        else
            report_warning "Package $pkg_name has no Go files, only $test_files_count test files"
        fi
    else
        report_warning "Package $pkg_name has no Go files"
    fi
done

echo
echo "3. Checking directory naming conventions..."

# Check all directories for naming convention violations
all_dirs=$(find . -type d -not -path "./.git/*" -not -path "./vendor/*" -not -path "./.kiro/*")
for dir in $all_dirs; do
    dir_name=$(basename "$dir")
    
    # Skip special directories
    if [[ "$dir_name" == "." ]] || [[ "$dir_name" == ".git" ]] || [[ "$dir_name" == ".github" ]] || [[ "$dir_name" == ".vscode" ]] || [[ "$dir_name" == ".kiro" ]]; then
        continue
    fi
    
    # Check naming convention (should be lowercase, no underscores for Go packages)
    if [[ "$dir" == *"/pkg/"* ]] || [[ "$dir" == *"/internal/"* ]] || [[ "$dir" == *"/cmd/"* ]]; then
        if [[ "$dir_name" =~ ^[a-z][a-z0-9]*$ ]]; then
            report_success "Directory $dir follows Go naming convention"
        else
            if [[ "$dir_name" =~ [_-] ]]; then
                report_issue "Directory $dir violates Go naming convention (contains underscores/hyphens)"
            else
                report_warning "Directory $dir may violate Go naming convention"
            fi
        fi
    fi
done

echo
echo "4. Checking for orphaned files..."

# Check for Go files in unexpected locations
unexpected_go_files=$(find . -name "*.go" -not -path "./pkg/*" -not -path "./internal/*" -not -path "./cmd/*" -not -path "./test/*" -not -path "./vendor/*" -not -path "./scripts/*" -type f)

if [ -n "$unexpected_go_files" ]; then
    for file in $unexpected_go_files; do
        report_warning "Go file in unexpected location: $file"
    done
else
    report_success "No Go files in unexpected locations"
fi

echo
echo "=== File Organization Summary ==="
echo "Issues found: $ISSUES"
echo "Warnings: $WARNINGS"

if [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✓ File organization check passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ File organization check failed with $ISSUES issues${NC}"
    exit 1
fi