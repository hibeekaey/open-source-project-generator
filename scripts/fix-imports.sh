#!/bin/bash

# Script to fix import organization and remove unused imports

echo "=== Fixing Import Organization ==="
echo "Timestamp: $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

FIXED=0
ERRORS=0

report_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    ((ERRORS++))
}

report_success() {
    echo -e "${GREEN}[FIXED]${NC} $1"
    ((FIXED++))
}

report_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

echo "1. Running goimports to fix import organization..."

# Check if goimports is installed
GOIMPORTS_PATH=$(which goimports 2>/dev/null || echo "$HOME/go/bin/goimports")
if [ ! -f "$GOIMPORTS_PATH" ]; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    GOIMPORTS_PATH="$HOME/go/bin/goimports"
fi

echo "Using goimports at: $GOIMPORTS_PATH"

# Run goimports on all Go files
go_files=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -type f)

file_count=0
for file in $go_files; do
    ((file_count++))
    echo "Processing file $file_count: $file"
    
    # Create backup
    cp "$file" "$file.bak"
    
    # Run goimports
    if "$GOIMPORTS_PATH" -w "$file" 2>/dev/null; then
        # Check if file was modified
        if ! cmp -s "$file" "$file.bak"; then
            report_success "Fixed imports in $file"
        fi
    else
        report_error "Failed to fix imports in $file"
        # Restore backup
        mv "$file.bak" "$file"
        continue
    fi
    
    # Remove backup if successful
    rm "$file.bak"
done

echo "Processed $file_count Go files"

echo
echo "2. Running go mod tidy to clean up module dependencies..."

if go mod tidy; then
    report_success "Cleaned up go.mod and go.sum"
else
    report_error "Failed to run go mod tidy"
fi

echo
echo "3. Checking for remaining unused imports..."

# Run go vet to check for issues (excluding scripts directory)
go vet $(go list ./... | grep -v "/scripts") && report_success "No vet issues found" || report_error "Go vet found issues"

echo
echo "4. Running golangci-lint to check for import issues..."

# Check if golangci-lint is available
if command -v golangci-lint &> /dev/null; then
    golangci-lint run --disable-all --enable=goimports,unused,ineffassign && report_success "No linting issues found" || report_error "Linting found issues"
else
    report_info "golangci-lint not available, skipping advanced checks"
fi

echo
echo "=== Import Fix Summary ==="
echo "Files fixed: $FIXED"
echo "Errors: $ERRORS"

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✓ Import organization completed successfully!${NC}"
    exit 0
else
    echo -e "${RED}✗ Import organization completed with $ERRORS errors${NC}"
    exit 1
fi