#!/bin/bash

# Dead Code Analysis and Removal Script
# This script runs comprehensive dead code analysis and provides options for removal

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ANALYSIS_SCRIPT="$SCRIPT_DIR/dead_code_analysis.go"
REPORT_FILE="$PROJECT_ROOT/dead_code_analysis_report.md"
BACKUP_DIR="$PROJECT_ROOT/.dead_code_backups"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Function to create backup
create_backup() {
    local file="$1"
    local backup_file="$BACKUP_DIR/$(basename "$file").backup.$(date +%Y%m%d_%H%M%S)"
    
    mkdir -p "$BACKUP_DIR"
    cp "$file" "$backup_file"
    print_status "Created backup: $backup_file"
}

# Function to remove unused imports from a file
remove_unused_imports() {
    local file="$1"
    print_status "Removing unused imports from $file"
    
    # Use goimports to remove unused imports
    if command -v goimports >/dev/null 2>&1; then
        create_backup "$file"
        goimports -w "$file"
        print_success "Removed unused imports from $file"
    else
        print_warning "goimports not found. Please install it: go install golang.org/x/tools/cmd/goimports@latest"
    fi
}

# Function to remove commented code blocks
remove_commented_code() {
    local file="$1"
    local start_line="$2"
    local end_line="$3"
    
    print_status "Removing commented code block from $file (lines $start_line-$end_line)"
    create_backup "$file"
    
    # Use sed to remove the lines
    sed -i.tmp "${start_line},${end_line}d" "$file"
    rm -f "$file.tmp"
    
    print_success "Removed commented code block from $file"
}

# Function to remove unused functions/methods
remove_unused_function() {
    local file="$1"
    local function_name="$2"
    local line_number="$3"
    
    print_status "Removing unused function '$function_name' from $file (around line $line_number)"
    create_backup "$file"
    
    # This is a simplified removal - in practice, you'd want more sophisticated parsing
    # For now, we'll just mark it for manual review
    print_warning "Function removal requires manual review: $function_name in $file:$line_number"
}

# Function to run the dead code analysis
run_analysis() {
    print_status "Running dead code analysis on $PROJECT_ROOT"
    
    cd "$PROJECT_ROOT"
    go run "$ANALYSIS_SCRIPT" "$PROJECT_ROOT"
    
    if [ -f "$REPORT_FILE" ]; then
        print_success "Analysis complete. Report generated: $REPORT_FILE"
        return 0
    else
        print_error "Analysis failed or report not generated"
        return 1
    fi
}

# Function to process unused imports
process_unused_imports() {
    print_status "Processing unused imports..."
    
    # Extract unused imports from the report
    if [ -f "$REPORT_FILE" ]; then
        # Parse the report and remove unused imports
        grep -A 100 "## Unused Imports" "$REPORT_FILE" | grep "^- " | sed 's/.*in \([^:]*\):[0-9]*.*/\1/' | sort -u | while read -r file; do
            if [ -f "$file" ]; then
                remove_unused_imports "$file"
            fi
        done
    fi
}

# Function to process commented code blocks
process_commented_code() {
    print_status "Processing commented-out code blocks..."
    
    if [ -f "$REPORT_FILE" ]; then
        # Parse the report and remove commented code blocks
        grep -A 100 "## Commented-Out Code Blocks" "$REPORT_FILE" | grep "^- " | grep "confidence: 0\.[89]" | while IFS= read -r line; do
            file=$(echo "$line" | sed 's/^- \([^:]*\):.*/\1/')
            start_line=$(echo "$line" | sed 's/.*:\([0-9]*\)-[0-9]*.*/\1/')
            end_line=$(echo "$line" | sed 's/.*:[0-9]*-\([0-9]*\).*/\1/')
            if [ -f "$file" ]; then
                remove_commented_code "$file" "$start_line" "$end_line"
            fi
        done
    fi
}

# Function to run gofmt on all Go files
format_code() {
    print_status "Formatting Go code..."
    find "$PROJECT_ROOT" -name "*.go" -not -path "*/vendor/*" -not -path "*/.git/*" | while read -r file; do
        gofmt -w "$file"
    done
    print_success "Code formatting complete"
}

# Function to run tests to ensure nothing is broken
run_tests() {
    print_status "Running tests to verify changes..."
    cd "$PROJECT_ROOT"
    
    if make test >/dev/null 2>&1; then
        print_success "All tests passed"
        return 0
    else
        print_error "Tests failed. Please review the changes."
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -a, --analyze-only     Run analysis only, don't remove anything"
    echo "  -i, --imports-only     Remove unused imports only"
    echo "  -c, --comments-only    Remove commented code blocks only"
    echo "  -f, --full             Run full analysis and removal (default)"
    echo "  -t, --test             Run tests after changes"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --analyze-only      # Just run analysis and generate report"
    echo "  $0 --imports-only      # Remove unused imports only"
    echo "  $0 --full --test       # Full cleanup with test validation"
}

# Main execution
main() {
    local analyze_only=false
    local imports_only=false
    local comments_only=false
    local run_tests_after=false
    local full_cleanup=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -a|--analyze-only)
                analyze_only=true
                full_cleanup=false
                shift
                ;;
            -i|--imports-only)
                imports_only=true
                full_cleanup=false
                shift
                ;;
            -c|--comments-only)
                comments_only=true
                full_cleanup=false
                shift
                ;;
            -f|--full)
                full_cleanup=true
                shift
                ;;
            -t|--test)
                run_tests_after=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    print_status "Starting dead code analysis and removal process..."
    
    # Always run analysis first
    if ! run_analysis; then
        print_error "Analysis failed. Exiting."
        exit 1
    fi
    
    # If analyze-only, just show the report and exit
    if [ "$analyze_only" = true ]; then
        print_status "Analysis complete. Review the report at: $REPORT_FILE"
        exit 0
    fi
    
    # Process based on options
    if [ "$imports_only" = true ]; then
        process_unused_imports
    elif [ "$comments_only" = true ]; then
        process_commented_code
    elif [ "$full_cleanup" = true ]; then
        process_unused_imports
        process_commented_code
        format_code
    fi
    
    # Run tests if requested
    if [ "$run_tests_after" = true ]; then
        if ! run_tests; then
            print_error "Tests failed after cleanup. Check the backups in $BACKUP_DIR"
            exit 1
        fi
    fi
    
    print_success "Dead code analysis and cleanup complete!"
    print_status "Backups are available in: $BACKUP_DIR"
    print_status "Report is available at: $REPORT_FILE"
}

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Run main function with all arguments
main "$@"