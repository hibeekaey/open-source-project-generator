#!/bin/bash

# CI Test Script
# This script runs the basic test suite for CI/CD pipelines

set -e

echo "🧪 Running basic test suite..."
echo "ℹ️  Running essential tests for core functionality"
echo ""

# Run basic test suite
go test -v -timeout=5m ./...

echo ""
echo "✅ Basic test suite completed successfully!"
echo "📊 All essential tests passed"
echo ""
echo "Note: This includes only essential tests:"
echo "  - Core template engine tests"
echo "  - Basic CLI functionality tests"
echo "  - Essential file operations tests"
echo "  - Configuration management tests"
echo ""
echo "The same command runs locally and in CI: go test -v ./..."