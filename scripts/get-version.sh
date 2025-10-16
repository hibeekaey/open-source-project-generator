#!/bin/bash

# Version Helper Script
# This script provides consistent version information across all build tools
# It follows semantic versioning and integrates with git tags

set -e

# Get version from git tags with fallback
get_version() {
    local version=""
    
    # Check if VERSION is set in environment (override)
    if [ -n "${VERSION:-}" ]; then
        echo "$VERSION"
        return 0
    fi
    
    # Try to get version from git tags
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # Check if we're on a tagged commit
        if git describe --exact-match --tags HEAD 2>/dev/null; then
            version=$(git describe --exact-match --tags HEAD 2>/dev/null)
        else
            # Get the latest tag and add commit info
            version=$(git describe --tags --always --dirty 2>/dev/null || echo "")
        fi
    fi
    
    # Fallback to "dev" if no version found
    if [ -z "$version" ]; then
        version="dev"
    fi
    
    echo "$version"
}

# Get git commit hash
get_commit() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git rev-parse --short HEAD 2>/dev/null || echo "unknown"
    else
        echo "unknown"
    fi
}

# Get build time
get_build_time() {
    date -u '+%Y-%m-%d_%H:%M:%S'
}

# Get git branch
get_branch() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown"
    else
        echo "unknown"
    fi
}

# Validate version format (semantic versioning)
validate_version() {
    local version=$1
    
    # Allow "dev" as a special case
    if [ "$version" = "dev" ]; then
        return 0
    fi
    
    # Check semantic versioning format: X.Y.Z or vX.Y.Z or X.Y.Z-prerelease
    if [[ "$version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.+-]+)?$ ]]; then
        return 0
    fi
    
    # Allow git describe format: vX.Y.Z-N-gHASH or X.Y.Z-N-gHASH
    if [[ "$version" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-g[0-9a-f]+(-dirty)?$ ]]; then
        return 0
    fi
    
    return 1
}

# Show usage
show_usage() {
    cat << EOF
Version Helper Script

Usage: $0 [COMMAND]

Commands:
  version         Get version string (default)
  commit          Get git commit hash
  build-time      Get build timestamp
  branch          Get git branch name
  validate        Validate version format
  all             Get all version information
  ldflags         Get Go build ldflags for version injection

Examples:
  $0                      # Get version
  $0 version              # Get version
  $0 commit               # Get commit hash
  $0 all                  # Get all info
  $0 ldflags              # Get ldflags for go build

Environment Variables:
  VERSION         Override version (useful for CI/CD)

Output:
  Prints requested information to stdout
EOF
}

# Main command handler
case "${1:-version}" in
    "version")
        get_version
        ;;
    "commit")
        get_commit
        ;;
    "build-time")
        get_build_time
        ;;
    "branch")
        get_branch
        ;;
    "validate")
        version=$(get_version)
        if validate_version "$version"; then
            echo "✓ Valid version: $version"
            exit 0
        else
            echo "✗ Invalid version format: $version"
            exit 1
        fi
        ;;
    "all")
        echo "Version:    $(get_version)"
        echo "Commit:     $(get_commit)"
        echo "Branch:     $(get_branch)"
        echo "Build Time: $(get_build_time)"
        ;;
    "ldflags")
        version=$(get_version)
        commit=$(get_commit)
        build_time=$(get_build_time)
        echo "-X main.Version=${version} -X main.GitCommit=${commit} -X main.BuildTime=${build_time}"
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        echo "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac
