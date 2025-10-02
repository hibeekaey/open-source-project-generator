#!/bin/bash

# Comprehensive Dead Code Analysis Script
# This script identifies potentially unused code in the Go project

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

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Create output directory
mkdir -p analysis_output

log_info "Starting comprehensive dead code analysis..."

# 1. Find potentially unused exported functions
log_info "Analyzing exported functions..."

cat > analysis_output/unused_exported_functions.txt << 'EOF'
# Potentially Unused Exported Functions
# These functions are exported but may not be used outside their package

EOF

# Find all exported functions and check their usage
find pkg/ -name "*.go" -exec grep -l "^func [A-Z]" {} \; | while read -r file; do
    grep "^func [A-Z][a-zA-Z0-9_]*(" "$file" | while read -r line; do
        func_name=$(echo "$line" | sed 's/^func \([A-Z][a-zA-Z0-9_]*\).*/\1/')
        
        # Count usage across the codebase (excluding the definition file)
        usage_count=$(find . -name "*.go" -not -path "./.git/*" -exec grep -l "\b$func_name\b" {} \; | grep -v "$file" | wc -l)
        
        if [ "$usage_count" -eq 0 ]; then
            echo "UNUSED: $func_name in $file" >> analysis_output/unused_exported_functions.txt
        fi
    done
done

# 2. Find potentially unused struct methods
log_info "Analyzing struct methods..."

cat > analysis_output/unused_methods.txt << 'EOF'
# Potentially Unused Methods
# These methods may not be called anywhere in the codebase

EOF

find pkg/ -name "*.go" -exec grep -l "^func ([^)]*) [A-Z]" {} \; | while read -r file; do
    grep "^func ([^)]*) [A-Z][a-zA-Z0-9_]*(" "$file" | while read -r line; do
        method_name=$(echo "$line" | sed 's/^func ([^)]*) \([A-Z][a-zA-Z0-9_]*\).*/\1/')
        
        # Skip common interface methods
        if [[ "$method_name" =~ ^(Error|String|GoString|MarshalJSON|UnmarshalJSON)$ ]]; then
            continue
        fi
        
        # Count usage across the codebase (excluding the definition file)
        usage_count=$(find . -name "*.go" -not -path "./.git/*" -exec grep -l "\b$method_name\b" {} \; | grep -v "$file" | wc -l)
        
        if [ "$usage_count" -eq 0 ]; then
            echo "UNUSED METHOD: $method_name in $file" >> analysis_output/unused_methods.txt
        fi
    done
done

# 3. Find unused imports
log_info "Analyzing unused imports..."

cat > analysis_output/unused_imports.txt << 'EOF'
# Files with Potentially Unused Imports
# Run 'go mod tidy' and 'goimports' to clean these up

EOF

find pkg/ cmd/ internal/ -name "*.go" | while read -r file; do
    # Use go list to check for unused imports
    if ! go list -f '{{.ImportPath}}' "$file" >/dev/null 2>&1; then
        continue
    fi
    
    # Check if goimports would make changes
    if command -v goimports >/dev/null 2>&1; then
        if ! goimports -l "$file" | grep -q "^$"; then
            echo "HAS UNUSED IMPORTS: $file" >> analysis_output/unused_imports.txt
        fi
    fi
done

# 4. Find unused constants and variables
log_info "Analyzing unused constants and variables..."

cat > analysis_output/unused_constants_vars.txt << 'EOF'
# Potentially Unused Constants and Variables
# These may be defined but not used

EOF

find pkg/ -name "*.go" | while read -r file; do
    # Find exported constants
    grep "^const [A-Z]" "$file" | while read -r line; do
        const_name=$(echo "$line" | sed 's/^const \([A-Z][a-zA-Z0-9_]*\).*/\1/')
        
        usage_count=$(find . -name "*.go" -not -path "./.git/*" -exec grep -l "\b$const_name\b" {} \; | grep -v "$file" | wc -l)
        
        if [ "$usage_count" -eq 0 ]; then
            echo "UNUSED CONST: $const_name in $file" >> analysis_output/unused_constants_vars.txt
        fi
    done
    
    # Find exported variables
    grep "^var [A-Z]" "$file" | while read -r line; do
        var_name=$(echo "$line" | sed 's/^var \([A-Z][a-zA-Z0-9_]*\).*/\1/')
        
        usage_count=$(find . -name "*.go" -not -path "./.git/*" -exec grep -l "\b$var_name\b" {} \; | grep -v "$file" | wc -l)
        
        if [ "$usage_count" -eq 0 ]; then
            echo "UNUSED VAR: $var_name in $file" >> analysis_output/unused_constants_vars.txt
        fi
    done
done

# 5. Find unused struct types
log_info "Analyzing unused struct types..."

cat > analysis_output/unused_structs.txt << 'EOF'
# Potentially Unused Struct Types
# These structs may be defined but not instantiated

EOF

find pkg/ -name "*.go" | while read -r file; do
    grep "^type [A-Z][a-zA-Z0-9_]* struct" "$file" | while read -r line; do
        struct_name=$(echo "$line" | sed 's/^type \([A-Z][a-zA-Z0-9_]*\) struct.*/\1/')
        
        # Count usage (instantiation, type assertions, etc.)
        usage_count=$(find . -name "*.go" -not -path "./.git/*" -exec grep -l "\b$struct_name\b" {} \; | grep -v "$file" | wc -l)
        
        if [ "$usage_count" -eq 0 ]; then
            echo "UNUSED STRUCT: $struct_name in $file" >> analysis_output/unused_structs.txt
        fi
    done
done

# 6. Find commented-out code blocks
log_info "Analyzing commented-out code..."

cat > analysis_output/commented_code.txt << 'EOF'
# Files with Potentially Commented-Out Code
# These may contain large blocks of commented code that can be removed

EOF

find pkg/ cmd/ internal/ -name "*.go" | while read -r file; do
    # Look for files with many consecutive comment lines (potential commented code)
    comment_blocks=$(grep -n "^\s*//.*" "$file" | awk -F: '{print $1}' | uniq -c | awk '$1 > 5 {print $2}')
    
    if [ -n "$comment_blocks" ]; then
        echo "LARGE COMMENT BLOCKS: $file (lines: $comment_blocks)" >> analysis_output/commented_code.txt
    fi
done

# 7. Find duplicate function implementations
log_info "Analyzing duplicate functions..."

cat > analysis_output/duplicate_functions.txt << 'EOF'
# Potentially Duplicate Functions
# These functions may have similar implementations

EOF

# Create a temporary file to store function signatures
temp_file=$(mktemp)

find pkg/ -name "*.go" -exec grep -H "^func " {} \; | while read -r line; do
    file=$(echo "$line" | cut -d: -f1)
    func_sig=$(echo "$line" | cut -d: -f2- | sed 's/^func //' | sed 's/{.*//')
    echo "$func_sig|$file" >> "$temp_file"
done

# Find functions with similar names
sort "$temp_file" | while read -r entry; do
    func_sig=$(echo "$entry" | cut -d'|' -f1)
    file=$(echo "$entry" | cut -d'|' -f2)
    func_name=$(echo "$func_sig" | sed 's/(.*//')
    
    # Look for similar function names
    similar_count=$(grep -c "^$func_name" "$temp_file" || true)
    if [ "$similar_count" -gt 1 ]; then
        echo "SIMILAR FUNCTION: $func_name in $file" >> analysis_output/duplicate_functions.txt
    fi
done

rm -f "$temp_file"

# 8. Generate summary report
log_info "Generating summary report..."

cat > analysis_output/summary.md << 'EOF'
# Dead Code Analysis Summary

This analysis identifies potentially unused code in the project. Please review each item carefully as some may be:
- Interface implementations that are used via interfaces
- Functions used in tests or examples
- Code intended for future use
- False positives due to reflection or dynamic usage

## Files Generated:
- `unused_exported_functions.txt` - Exported functions not used outside their package
- `unused_methods.txt` - Methods that may not be called
- `unused_imports.txt` - Files with unused imports
- `unused_constants_vars.txt` - Unused constants and variables
- `unused_structs.txt` - Struct types that may not be instantiated
- `commented_code.txt` - Files with large comment blocks
- `duplicate_functions.txt` - Functions with similar names/implementations

## Recommended Actions:

1. **Review unused exports**: Consider making them unexported if they're only used internally
2. **Clean up imports**: Run `goimports -w .` to remove unused imports
3. **Remove dead code**: Carefully remove confirmed unused code
4. **Refactor duplicates**: Consolidate similar functions where appropriate
5. **Clean comments**: Remove large blocks of commented-out code

## Verification Steps:

Before removing any code:
1. Run all tests: `go test ./...`
2. Check for interface implementations
3. Verify the code isn't used via reflection
4. Ensure it's not part of a public API that external users might depend on

EOF

# Count findings
unused_funcs=$(grep -c "UNUSED:" analysis_output/unused_exported_functions.txt || echo "0")
unused_methods=$(grep -c "UNUSED METHOD:" analysis_output/unused_methods.txt || echo "0")
unused_imports=$(grep -c "HAS UNUSED IMPORTS:" analysis_output/unused_imports.txt || echo "0")
unused_consts=$(grep -c "UNUSED" analysis_output/unused_constants_vars.txt || echo "0")
unused_structs=$(grep -c "UNUSED STRUCT:" analysis_output/unused_structs.txt || echo "0")
comment_blocks=$(grep -c "LARGE COMMENT BLOCKS:" analysis_output/commented_code.txt || echo "0")
duplicate_funcs=$(grep -c "SIMILAR FUNCTION:" analysis_output/duplicate_functions.txt || echo "0")

cat >> analysis_output/summary.md << EOF

## Statistics:
- Potentially unused exported functions: $unused_funcs
- Potentially unused methods: $unused_methods
- Files with unused imports: $unused_imports
- Unused constants/variables: $unused_consts
- Unused struct types: $unused_structs
- Files with large comment blocks: $comment_blocks
- Functions with similar names: $duplicate_funcs

Generated on: $(date)
EOF

log_success "Dead code analysis completed!"
log_info "Results saved in analysis_output/ directory"
log_info "Review analysis_output/summary.md for an overview"

# Display summary
echo ""
echo "=== ANALYSIS SUMMARY ==="
echo "Potentially unused exported functions: $unused_funcs"
echo "Potentially unused methods: $unused_methods"
echo "Files with unused imports: $unused_imports"
echo "Unused constants/variables: $unused_consts"
echo "Unused struct types: $unused_structs"
echo "Files with large comment blocks: $comment_blocks"
echo "Functions with similar names: $duplicate_funcs"
echo ""
echo "Review the files in analysis_output/ for detailed findings."