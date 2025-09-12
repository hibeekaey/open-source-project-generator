#!/bin/bash

# Script to audit npm dependencies in template files
set -e

echo "=== Auditing npm dependencies in template files ==="

# Create temporary directory for auditing
TEMP_DIR=$(mktemp -d)
echo "Using temporary directory: $TEMP_DIR"

# Function to create package.json from template and audit it
audit_template() {
    local template_file="$1"
    local template_name=$(basename "$template_file" .tmpl)
    local temp_package_dir="$TEMP_DIR/$template_name"
    
    echo "Auditing $template_file..."
    
    # Create directory for this template
    mkdir -p "$temp_package_dir"
    
    # Create a basic package.json by replacing template variables with defaults
    sed -e 's/{{\.Versions\.NextJS}}/15.5.3/g' \
        -e 's/{{\.Versions\.React}}/19.0.0/g' \
        -e 's/{{\.Organization | lower}}/example/g' \
        -e 's/{{\.Name | lower}}/example-app/g' \
        -e 's/{{\.Name}}/ExampleApp/g' \
        -e 's/{{\.Description}}/Example Application/g' \
        -e 's/{{\.Author}}/Example Author/g' \
        -e 's/{{\.Email}}/author@example.com/g' \
        -e 's/{{\.License}}/MIT/g' \
        -e 's/{{if \.Repository}}//g' \
        -e 's/{{end}}//g' \
        -e 's/{{\.Repository}}/https:\/\/github.com\/example\/repo/g' \
        "$template_file" > "$temp_package_dir/package.json"
    
    # Run npm audit in the temporary directory
    cd "$temp_package_dir"
    
    # Install dependencies to create package-lock.json
    echo "Installing dependencies for $template_name..."
    npm install --package-lock-only --silent 2>/dev/null || true
    
    # Run audit
    echo "Running audit for $template_name..."
    if npm audit --audit-level=moderate 2>/dev/null; then
        echo "✅ No vulnerabilities found in $template_name"
    else
        echo "⚠️  Vulnerabilities found in $template_name"
        npm audit --audit-level=moderate || true
    fi
    
    cd - > /dev/null
}

# Find and audit all package.json template files
find templates -name "package.json.tmpl" | while read template_file; do
    audit_template "$template_file"
done

# Cleanup
rm -rf "$TEMP_DIR"
echo "=== npm audit complete ==="