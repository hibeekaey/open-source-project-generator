# Version Management Guide

This document explains how version information is managed throughout the Open Source Project Generator project.

## Overview

The project uses **semantic versioning** (SemVer) with git tags as the source of truth. Version information is automatically derived from git tags and injected into binaries at build time.

## Version Format

### Semantic Versioning

```text
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

Examples:

- `1.0.0` - Release version
- `1.0.0-beta.1` - Pre-release version
- `1.0.0-rc.1` - Release candidate
- `v1.2.0-27-gbc317ef` - Development version (27 commits after v1.2.0)
- `dev` - Development version (no tags)

### Version Components

1. **MAJOR** - Incompatible API changes
2. **MINOR** - Backwards-compatible functionality additions
3. **PATCH** - Backwards-compatible bug fixes
4. **PRERELEASE** - Optional pre-release identifier (alpha, beta, rc)
5. **BUILD** - Optional build metadata

## Version Sources

### 1. Git Tags (Primary Source)

The version is automatically derived from git tags:

```bash
# Get current version
git describe --tags --always --dirty

# Examples:
# v1.0.0              - On a tagged commit
# v1.0.0-5-g1234567   - 5 commits after v1.0.0
# v1.0.0-dirty        - Uncommitted changes
# dev                 - No tags found
```

**Important**: The version from git tags is used exactly as returned. If your tags have a `v` prefix, the version will include it. If not, it won't. This ensures consistency across all build artifacts.

### 2. Environment Variable (Override)

You can override the version using the `VERSION` environment variable:

```bash
VERSION=1.2.3 make build
VERSION=v2.0.0-beta.1 ./scripts/build.sh
```

**Note**: When you set VERSION manually, it will be used exactly as provided without any modification. Choose your format consistently.

### 3. CI/CD Workflows

In GitHub Actions, version is determined by:

- **Tag pushes**: Uses the tag name (e.g., `v1.0.0`)
- **Manual dispatch**: Uses the input version
- **Branch pushes**: Uses git describe

## Version Injection

### Go Binary

Version information is injected at build time using `-ldflags`:

```bash
go build -ldflags "\
  -X main.Version=${VERSION} \
  -X main.GitCommit=${GIT_COMMIT} \
  -X main.BuildTime=${BUILD_TIME}" \
  ./cmd/generator
```

### Variables in main.go

```go
var (
    Version   = "dev"        // Overridden at build time
    GitCommit = "unknown"    // Overridden at build time
    BuildTime = "unknown"    // Overridden at build time
)
```

## Version Helper Script

The `scripts/get-version.sh` script provides consistent version information:

```bash
# Get version
./scripts/get-version.sh version

# Get commit hash
./scripts/get-version.sh commit

# Get build time
./scripts/get-version.sh build-time

# Get all info
./scripts/get-version.sh all

# Get ldflags for go build
./scripts/get-version.sh ldflags

# Validate version format
./scripts/get-version.sh validate
```

## Build Tools Integration

### Makefile

The Makefile automatically derives version from git:

```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
```

Usage:

```bash
make build                    # Uses git-derived version
VERSION=1.2.3 make build      # Override version
```

### Build Scripts

All build scripts (`build.sh`, `build-packages.sh`) use the same logic:

```bash
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
```

**Critical**: Scripts use VERSION as a single source of truth. The value is never modified by stripping or adding prefixes. This ensures that:

- Package names use the exact version you specify
- Binary version output matches your input
- Docker tags are consistent with your version
- All artifacts have the same version string

### Docker

Docker images are tagged with the version:

```bash
make docker-build             # Tags with git-derived version
VERSION=1.2.3 make docker-build  # Override version
```

## CI/CD Workflows

### Release Workflow

Triggered by tag pushes or manual dispatch:

```yaml
- name: Get version
  id: version
  run: |
    if [ "${{ github.event_name }}" = "workflow_dispatch" ]; then
      echo "version=${{ github.event.inputs.version }}" >> $GITHUB_OUTPUT
    else
      echo "version=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT
    fi

- name: Build binary
  env:
    VERSION: ${{ steps.version.outputs.version }}
  run: |
    LDFLAGS="-X main.Version=${VERSION} ..."
    go build -ldflags "${LDFLAGS}" ./cmd/generator
```

### Docker Workflow

Uses metadata action for version tagging:

```yaml
- name: Extract metadata
  id: meta
  uses: docker/metadata-action@v5
  with:
    tags: |
      type=semver,pattern={{version}}
      type=semver,pattern={{major}}.{{minor}}
      type=raw,value=latest,enable={{is_default_branch}}
```

## Creating a Release

### 1. Update Version

Ensure CHANGELOG.md is updated with the new version.

### 2. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push tag to trigger release workflow
git push origin v1.0.0
```

### 3. Automated Release

The GitHub Actions release workflow will:

1. Build binaries for all platforms
2. Create distribution packages (DEB, RPM, Arch)
3. Generate checksums
4. Create GitHub release with artifacts
5. Build and push Docker images

### 4. Manual Release (if needed)

```bash
# Build locally with specific version
VERSION=v1.0.0 make dist
VERSION=v1.0.0 make package

# Or use scripts directly
VERSION=v1.0.0 ./scripts/build.sh
VERSION=v1.0.0 ./scripts/build-packages.sh all
```

## Version Verification

### Check Binary Version

```bash
# After building
./bin/generator version

# Output:
# Generator version v1.0.0
# Git commit: abc1234
# Build time: 2025-10-16_12:00:00
```

### Check Docker Image Version

```bash
docker run --rm generator:latest version
```

### Check Package Version

```bash
# DEB package
dpkg -I packages/generator_1.0.0_amd64.deb | grep Version

# RPM package
rpm -qip packages/generator-1.0.0-1.x86_64.rpm | grep Version
```

## Version in Different Contexts

### Development

```bash
# Local development (uses git describe)
make build
./bin/generator version
# Output: v1.2.0-27-gbc317ef-dirty
```

### CI/CD (Branch)

```bash
# Automatic from git
make ci
# Version: v1.2.0-27-gbc317ef
```

### CI/CD (Tag)

```bash
# From tag name
git tag v1.0.0
git push origin v1.0.0
# Version: v1.0.0
```

### Manual Override

```bash
# Explicit version
VERSION=2.0.0-beta.1 make build
./bin/generator version
# Output: 2.0.0-beta.1
```

## Version Consistency Principle

**Single Source of Truth**: The VERSION variable is treated as the authoritative version string throughout the entire build system. It is:

- ✅ **Never modified** - No stripping of "v" prefix, no adding prefixes
- ✅ **Used exactly as provided** - Whether from git tags or environment variable
- ✅ **Consistent everywhere** - Same version in binaries, packages, Docker tags, and documentation
- ✅ **User-controlled** - You decide the format (with or without "v" prefix)

This approach eliminates version inconsistencies and ensures that what you specify is what you get in all build artifacts.

## Best Practices

### 1. Always Use Git Tags

- Create annotated tags for releases: `git tag -a v1.0.0 -m "Release 1.0.0"`
- Follow semantic versioning strictly
- Never delete or modify published tags

### 2. Version Format Consistency

- **Recommended**: Use `v` prefix for tags: `v1.0.0`
- **Be consistent**: Choose one format and stick with it across all tags
- **Trust the system**: VERSION is used exactly as you provide it
- **No surprises**: The version in your binary will match your tag exactly

### 3. Pre-release Versions

- Use standard identifiers: `alpha`, `beta`, `rc`
- Examples: `v1.0.0-alpha.1`, `v1.0.0-beta.2`, `v1.0.0-rc.1`

### 4. Development Versions

- Between releases: `v1.0.0-N-gHASH` (automatic from git describe)
- No tags: `dev` (fallback)

### 5. CI/CD Integration

- Let git tags drive releases
- Use manual dispatch only for hotfixes or special cases
- Always validate version format before release

## Troubleshooting

### Version Shows "dev"

**Cause**: No git tags found

**Solution**:

```bash
# Check if tags exist
git tag -l

# Fetch tags from remote
git fetch --tags

# Create initial tag if needed
git tag -a v0.1.0 -m "Initial version"
```

### Version Shows "dirty"

**Cause**: Uncommitted changes in working directory

**Solution**:

```bash
# Check status
git status

# Commit or stash changes
git add .
git commit -m "Your changes"

# Or clean working directory
git stash
```

### Wrong Version in Binary

**Cause**: Version not passed to build

**Solution**:

```bash
# Use Makefile (handles version automatically)
make build

# Or pass version explicitly
VERSION=v1.0.0 go build -ldflags "-X main.Version=v1.0.0" ./cmd/generator
```

### CI/CD Version Mismatch

**Cause**: Tag not pushed or workflow not triggered

**Solution**:

```bash
# Ensure tag is pushed
git push origin v1.0.0

# Check workflow runs in GitHub Actions
# Verify tag format matches workflow triggers
```

## Examples

### Local Development Build

```bash
# Build with automatic version
make build

# Check version
./bin/generator version
# Output: v1.2.0-27-gbc317ef-dirty
```

### Release Build

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release 1.0.0"
git push origin v1.0.0

# GitHub Actions automatically:
# - Builds binaries with VERSION=v1.0.0
# - Creates packages
# - Publishes release
```

### Custom Version Build

```bash
# Build with custom version
VERSION=2.0.0-beta.1 make build

# Build packages with custom version
VERSION=2.0.0-beta.1 ./scripts/build-packages.sh all

# Build Docker image with custom version
VERSION=2.0.0-beta.1 make docker-build
```

### Verify Version Consistency

```bash
# Check all version sources
./scripts/get-version.sh all

# Validate version format
./scripts/get-version.sh validate

# Get ldflags for manual build
LDFLAGS=$(./scripts/get-version.sh ldflags)
go build -ldflags "$LDFLAGS" ./cmd/generator
```

## Summary

- **Source of Truth**: Git tags
- **Format**: Semantic versioning (MAJOR.MINOR.PATCH)
- **Injection**: Build-time via `-ldflags`
- **Override**: `VERSION` environment variable
- **Automation**: Makefile and scripts handle automatically
- **CI/CD**: Triggered by tag pushes
- **Helper**: `scripts/get-version.sh` for consistency

The version management system ensures that every build, whether local or in CI/CD, has accurate and traceable version information.
