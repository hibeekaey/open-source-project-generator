#!/bin/bash

# Template Validation Script
# This script runs the template validation tool

set -e

echo "🔧 Building template validator..."
cd "$(dirname "$0")"

# Build the validator
go build -o validator .

echo "🔍 Running template validation..."

# Run validation on the templates directory
./validator ../../templates

echo "✅ Validation complete!"