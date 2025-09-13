#!/bin/bash

# CI Test Script
# This script runs the unified test suite for CI/CD pipelines
# All tests have been optimized and consolidated for reliable CI execution

set -e

echo "🧪 Running unified test suite..."
echo "ℹ️  Using optimized and consolidated test execution"
echo "ℹ️  All tests now run with improved performance and reliability"
echo ""

# Run unified test suite
go test -v -timeout=10m ./...

echo ""
echo "✅ Unified test suite completed successfully!"
echo "📊 All tests passed with consolidated execution"
echo ""
echo "Note: This now includes all tests in a unified execution:"
echo "  - Security validation tests (optimized for CI)"
echo "  - Template compilation tests (with mocked dependencies)"
echo "  - Template edge case tests (optimized for performance)"
echo "  - Integration tests (with improved reliability)"
echo ""
echo "The same command runs locally and in CI: go test -v ./..."