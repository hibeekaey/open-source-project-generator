# Distribution Guide

This document describes the distribution and release process for the Open Source Project Generator.

## Overview

The Open Source Project Generator supports multiple distribution methods:

- **Binary releases**: Cross-platform binaries for direct download
- **Package managers**: Native packages for Linux distributions (DEB, RPM, Arch)
- **Container images**: Docker images via GitHub Container Registry (ghcr.io)
- **Docker Compose**: Multi-profile development and deployment environments
- **Source builds**: Build from source code with Go 1.25+
- **Go install**: Direct installation via Go toolchain

## Quick Start

```bash
# Using Docker Compose (recommended for development)
docker compose --profile production run --rm generator generate --help

# Using Go install
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest

# Using install script
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash
```

## Build System

### Cross-Platform Builds

The build system creates binaries for multiple platforms using the Makefile:

```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Output: bin/generator-{os}-{arch}
# - bin/generator-linux-amd64
# - bin/generator-darwin-amd64
# - bin/generator-darwin-arm64
# - bin/generator-windows-amd64.exe
```

**Supported Platforms**:

- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

**Build Requirements**:

- Go 1.25.0 or later
- Make
- Git (for version information)

### Package Building

Create native packages for different distributions:

```bash
# Build all package types
make package-all

# Individual package types (via scripts/build-packages.sh)
./scripts/build-packages.sh deb    # Debian/Ubuntu
./scripts/build-packages.sh rpm    # Red Hat/CentOS/Fedora
./scripts/build-packages.sh arch   # Arch Linux
./scripts/build-packages.sh all    # All packages

# Output: packages/deb/, packages/rpm/, packages/arch/
```

**Package Building Requirements**:

- Docker (uses Dockerfile.build)
- Or native tools: dpkg-dev, rpm, fakeroot, debhelper

### Docker Images

Build container images using Make or Docker Compose:

```bash
# Production image (multi-stage, minimal - 39.2 MB)
make docker-build
# Output: ghcr.io/cuesoftinc/open-source-project-generator:latest

# Test the production image
make docker-test

# Development image (with all dev tools)
docker compose --profile development build

# Build image (for creating packages)
docker compose --profile build build
```

**Docker Image Types**:

| Image | Base | Size | Purpose | User |
|-------|------|------|---------|------|
| Dockerfile | alpine:3.19 | ~39 MB | Production runtime | generator (UID 1001) |
| Dockerfile.dev | golang:1.25-alpine | ~500 MB | Development | developer (UID 1001) |
| Dockerfile.build | ubuntu:24.04 | ~1.5 GB | Package building | builder (UID 1001) |

### Docker Compose Profiles

Use Docker Compose for streamlined workflows:

```bash
# Production - Generate projects
docker compose --profile production run --rm generator generate

# Development - Interactive shell with hot reload
docker compose --profile development run --rm generator-dev bash

# Testing - Run tests
docker compose --profile testing up generator-test
docker compose --profile testing up generator-test-coverage
docker compose --profile testing up generator-test-integration

# Build - Create binaries and packages
docker compose --profile build up generator-build
docker compose --profile build up generator-build-all
docker compose --profile build up generator-package-all

# Linting - Code quality checks
docker compose --profile lint up generator-lint
docker compose --profile lint up generator-fmt

# Security - Security scanning
docker compose --profile security up generator-security
docker compose --profile security up generator-gosec
docker compose --profile security up generator-govulncheck
docker compose --profile security up generator-staticcheck
```

See [docker-compose.yml](docker-compose.yml) for complete configuration.

## Release Process

### Automated Releases (GitHub Actions)

Releases are automated through GitHub Actions workflows:

1. **Tag Creation**: Push a version tag (e.g., `v1.0.0`)

   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **Build Pipeline**: Automatically:
   - Builds all platforms (Linux, macOS, Windows)
   - Creates distribution packages (DEB, RPM, Arch)
   - Runs tests and security scans
   - Generates checksums

3. **Release Creation**: Creates GitHub release with all artifacts

4. **Container Push**: Pushes images to GitHub Container Registry (ghcr.io)

5. **Notifications**: Sends notifications on completion

### Manual Release

For manual releases or testing:

```bash
# 1. Run full validation and build
make pre-release

# This runs: clean, build, test, test-coverage, lint, security-scan, build-all

# 2. Prepare release artifacts (includes packages)
make release-prepare

# This runs: test, lint, security-scan, dist, package-all

# 3. Create checksums
cd dist && sha256sum * > checksums.txt

# 4. Create GitHub release manually
# Upload files from dist/ and packages/ directories
```

### Version Management

Version information is embedded during build using ldflags:

```bash
# Build with version information
VERSION=1.0.0 \
GIT_COMMIT=$(git rev-parse --short HEAD) \
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
make build

# Version is embedded in binary
./bin/generator version

# Docker builds also accept version arguments
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
  -t generator:1.0.0 .
```

### Pre-Release Checklist

Before creating a release:

- [ ] Update version in relevant files
- [ ] Run `make pre-release` to validate everything
- [ ] Update CHANGELOG.md with release notes
- [ ] Test on all target platforms
- [ ] Verify Docker images build successfully
- [ ] Check security scan results
- [ ] Update documentation if needed

## Distribution Channels

### GitHub Releases

Primary distribution method:

- **URL**: `https://github.com/cuesoftinc/open-source-project-generator/releases`
- **Assets**: Binaries, packages, checksums
- **Automation**: Fully automated via GitHub Actions

#### Asset Naming Convention

The following assets are automatically generated for each release:

**Binary Archives:**

- `generator-linux-amd64.tar.gz` - Linux 64-bit
- `generator-linux-arm64.tar.gz` - Linux ARM64  
- `generator-darwin-amd64.tar.gz` - macOS Intel
- `generator-darwin-arm64.tar.gz` - macOS Apple Silicon
- `generator-windows-amd64.zip` - Windows 64-bit
- `generator-freebsd-amd64.tar.gz` - FreeBSD 64-bit

**Package Files:**

- `generator_VERSION_amd64.deb` - Debian/Ubuntu package
- `generator-VERSION-1.x86_64.rpm` - Red Hat/CentOS package

**Additional Files:**

- `checksums.txt` - SHA256 checksums for all assets

> **Note**: Replace `VERSION` with the actual release version (e.g., `1.0.0`)

### Package Repositories

#### Debian/Ubuntu (APT)

```bash
# Install from release
wget https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator_VERSION_amd64.deb
sudo dpkg -i generator_VERSION_amd64.deb
```

#### Red Hat/CentOS/Fedora (YUM/DNF)

```bash
# Install from release
sudo yum install https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-VERSION-1.x86_64.rpm
```

#### Arch Linux (AUR)

```bash
# Install from AUR (when available)
yay -S generator
```

### Container Registry

#### GitHub Container Registry (Primary)

All Docker images are published to GitHub Container Registry (ghcr.io). The build system uses dynamic organization names based on the repository owner.

**Official Images (cuesoftinc):**

```bash
# Pull latest production image
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# Pull specific version
docker pull ghcr.io/cuesoftinc/open-source-project-generator:v1.0.0

# Pull development image
docker pull ghcr.io/cuesoftinc/open-source-project-generator:dev

# Pull build image
docker pull ghcr.io/cuesoftinc/open-source-project-generator:build
```

**Fork Images (automatically adapts to your GitHub username):**

```bash
# Pull from your fork (replace 'your-username')
docker pull ghcr.io/your-username/open-source-project-generator:latest

# Pull specific version from your fork
docker pull ghcr.io/your-username/open-source-project-generator:v1.0.0
```

**Dynamic Configuration:**

The build system automatically detects the repository owner and builds images accordingly:

- Environment variable: `GITHUB_REPOSITORY_OWNER` (default: cuesoftinc)
- Docker registry: `DOCKER_REGISTRY` (default: ghcr.io)
- Configured in: Makefile, docker-compose.yml, env.example

When you fork the repository and run GitHub Actions, images will be published to your own container registry namespace.

**Image Tags:**

- `latest` - Latest stable release (production image)
- `v1.0.0` - Specific version (production image)
- `dev` - Development image with all tools
- `build` - Build image for creating packages

### Package Managers

#### Homebrew (macOS/Linux)

**Status**: Planned for future release

```bash
# Install via Homebrew (coming soon)
brew tap cuesoftinc/tap
brew install generator
```

To add Homebrew support, we need to:

1. Create a Homebrew formula
2. Submit to homebrew-core or create a tap
3. Maintain formula updates

#### Chocolatey (Windows)

**Status**: Planned for future release

```powershell
# Install via Chocolatey (coming soon)
choco install generator
```

To add Chocolatey support, we need to:

1. Create a Chocolatey package (.nuspec)
2. Submit to Chocolatey community repository
3. Maintain package updates

#### Scoop (Windows)

**Status**: Planned for future release

```powershell
# Install via Scoop (coming soon)
scoop bucket add cuesoftinc https://github.com/cuesoftinc/scoop-bucket
scoop install generator
```

To add Scoop support, we need to:

1. Create a Scoop manifest
2. Create a Scoop bucket repository
3. Maintain manifest updates

**Current Workaround**: Use the install script or download binaries directly from GitHub releases.

## Installation Methods

### Quick Install Script

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash

# With options
curl -sSL https://raw.githubusercontent.com/cuesoftinc/open-source-project-generator/main/scripts/install.sh | bash -s -- --version v1.3.0
```

### Go Install Method

```bash
# Install latest version
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest

# Install specific version
go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@v1.3.0
```

### Manual Installation

1. Download appropriate binary from releases
2. Extract archive
3. Move binary to PATH
4. Make executable (Unix systems)

### Package Manager Installation

Use native package managers when available:

```bash
# Debian/Ubuntu
sudo apt install generator

# Red Hat/CentOS
sudo yum install generator

# Arch Linux
pacman -S generator

# macOS
brew install generator

# Windows
choco install generator
```

## Build Configuration

### Environment Variables

See [env.example](env.example) for a complete list of environment variables.

**Build Variables:**

- `VERSION`: Version to embed in binary (default: dev)
- `GIT_COMMIT`: Git commit hash (default: unknown)
- `BUILD_TIME`: Build timestamp (default: unknown)
- `GOOS`: Target operating system
- `GOARCH`: Target architecture
- `CGO_ENABLED`: Enable/disable CGO (default: 0)

**Docker Variables:**

- `DOCKER_REGISTRY`: Docker registry (default: ghcr.io)
- `GITHUB_REPOSITORY_OWNER`: Repository owner (default: cuesoftinc)
- `GITHUB_ACTOR`: GitHub username for authentication
- `GITHUB_TOKEN`: GitHub token for pushing images

**Generator Variables:**

- `GENERATOR_LOG_LEVEL`: Log level (debug, info, warn, error)
- `GENERATOR_CONFIG_DIR`: Configuration directory
- `GENERATOR_CACHE_DIR`: Cache directory
- `GENERATOR_OUTPUT_PATH`: Output directory for generated projects

**Example Configuration:**

```bash
# Copy example environment file
cp env.example .env

# Edit with your values
vim .env

# Source for current session
source .env
```

### Build Flags

```bash
# Production build with optimizations
go build -ldflags="-w -s" -o generator ./cmd/generator

# Debug build with symbols
go build -gcflags="all=-N -l" -o generator ./cmd/generator

# Static build (Linux)
CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o generator ./cmd/generator
```

### Cross-Compilation

The Makefile handles cross-compilation automatically:

```bash
# Build for all platforms at once
make build-all

# Manual cross-compilation examples
GOOS=linux GOARCH=amd64 go build -o bin/generator-linux-amd64 ./cmd/generator
GOOS=linux GOARCH=arm64 go build -o bin/generator-linux-arm64 ./cmd/generator
GOOS=darwin GOARCH=amd64 go build -o bin/generator-darwin-amd64 ./cmd/generator
GOOS=darwin GOARCH=arm64 go build -o bin/generator-darwin-arm64 ./cmd/generator
GOOS=windows GOARCH=amd64 go build -o bin/generator-windows-amd64.exe ./cmd/generator

# With version information
VERSION=1.0.0 \
GIT_COMMIT=$(git rev-parse --short HEAD) \
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" \
  -o bin/generator-linux-amd64 ./cmd/generator
```

**Platform Support:**

| OS | Architecture | Binary Name | Status |
|----|--------------|-------------|--------|
| Linux | amd64 | generator-linux-amd64 | ✅ Supported |
| Linux | arm64 | generator-linux-arm64 | ✅ Supported |
| macOS | amd64 | generator-darwin-amd64 | ✅ Supported |
| macOS | arm64 | generator-darwin-arm64 | ✅ Supported |
| Windows | amd64 | generator-windows-amd64.exe | ✅ Supported |

## Quality Assurance

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run integration tests
make test-integration

# Test installation script syntax
make test-install

# Using Docker Compose
docker compose --profile testing up generator-test
docker compose --profile testing up generator-test-coverage
docker compose --profile testing up generator-test-integration
```

### Validation

```bash
# Validate packages
dpkg -I packages/deb/generator_VERSION_amd64.deb
rpm -qip packages/rpm/generator-VERSION-1.x86_64.rpm

# Test binaries
./bin/generator-linux-amd64 version
./bin/generator-darwin-amd64 version
./bin/generator-darwin-arm64 version
./bin/generator-windows-amd64.exe version

# Validate Docker images
docker run --rm generator:latest version
docker run --rm --entrypoint sh generator:latest -c "id"  # Check user is 1001
```

### Security Scanning

```bash
# Run all security scans
make security-scan

# Individual security tools
make gosec           # Run gosec security scanner
make govulncheck     # Check for Go vulnerabilities
make staticcheck     # Run static analysis

# Using Docker Compose
docker compose --profile security up generator-security
docker compose --profile security up generator-gosec
docker compose --profile security up generator-govulncheck
docker compose --profile security up generator-staticcheck

# Check dependencies for updates
go list -m -u all

# Verify checksums
sha256sum -c dist/checksums.txt
```

### Linting and Code Quality

```bash
# Run linter
make lint

# Format code
make fmt

# Run go vet
make vet

# Using Docker Compose
docker compose --profile lint up generator-lint
docker compose --profile lint up generator-fmt
```

## Monitoring and Analytics

### Download Statistics

Monitor release downloads via GitHub API:

```bash
# Get download stats
curl -s https://api.github.com/repos/cuesoftinc/open-source-project-generator/releases | jq '.[].assets[].download_count'
```

### Usage Analytics

Optional usage analytics (opt-in):

- Installation method tracking
- Platform distribution
- Feature usage statistics
- Error reporting

## CI/CD Pipeline

### GitHub Actions Workflows

The project uses GitHub Actions for continuous integration and deployment:

**Build and Test Workflow:**

- Triggered on: Push, Pull Request
- Runs: `make test`, `make lint`, `make security-scan`
- Tests on: Linux, macOS, Windows
- Go versions: 1.25+

**Release Workflow:**

- Triggered on: Tag push (v*)
- Builds: All platforms and packages
- Creates: GitHub release with artifacts
- Pushes: Docker images to ghcr.io
- Generates: Checksums and release notes

**Security Workflow:**

- Triggered on: Schedule (daily), Pull Request
- Runs: CodeQL analysis, dependency scanning
- Tools: gosec, govulncheck, staticcheck
- Reports: Security advisories

### Local CI/CD Testing

Test the full CI/CD pipeline locally:

```bash
# Run full CI pipeline
make ci

# This runs: lint, test, security-scan

# Run pre-commit checks
make pre-commit

# This runs: fmt, vet, lint, test

# Run full pre-release validation
make pre-release

# This runs: clean, build, test, test-coverage, lint, security-scan, build-all
```

## Maintenance

### Regular Tasks

1. **Dependency Updates**: Update Go modules monthly

   ```bash
   go get -u ./...
   go mod tidy
   make test
   ```

2. **Security Patches**: Apply security updates immediately

   ```bash
   make security-scan
   make govulncheck
   ```

3. **Platform Testing**: Test on all supported platforms

   ```bash
   make build-all
   # Test each binary
   ```

4. **Documentation**: Keep installation guides updated
   - Update version numbers
   - Test installation commands
   - Update screenshots/examples

5. **Docker Images**: Rebuild and test regularly

   ```bash
   make docker-build
   make docker-test
   ```

### Automation

- **Dependabot**: Automated dependency updates (configured in `.github/dependabot.yml`)
- **Security Scanning**: CodeQL and vulnerability scanning (daily)
- **Build Testing**: Continuous integration on all platforms (on push/PR)
- **Release Automation**: Automated releases on tag push
- **Docker Builds**: Automated image builds and pushes to ghcr.io

## Troubleshooting

### Build Issues

```bash
# Clean all build artifacts
make clean

# Clean Go caches
go clean -cache -modcache -testcache

# Rebuild dependencies
go mod download
go mod tidy
go mod verify

# Verbose build
go build -v -x ./cmd/generator

# Check Go environment
go env

# Verify Go version (requires 1.25.0+)
go version
```

### Docker Build Issues

```bash
# Clean Docker build cache
docker builder prune -a

# Build without cache
docker build --no-cache -f Dockerfile -t generator:test .

# Check Docker Compose configuration
docker compose config --quiet

# View Docker Compose service configuration
docker compose --profile production config

# Test specific Dockerfile
docker build -f Dockerfile.dev -t generator-dev:test .
docker build -f Dockerfile.build -t generator-build:test .
```

### Package Issues

```bash
# Validate package structure
dpkg-deb --contents packages/deb/generator_VERSION_amd64.deb
rpm2cpio packages/rpm/generator-VERSION-1.x86_64.rpm | cpio -tv

# Test package installation in Docker
docker run --rm -v $(pwd)/packages:/packages ubuntu:24.04 bash -c \
  "apt update && dpkg -i /packages/deb/generator_VERSION_amd64.deb && generator version"

docker run --rm -v $(pwd)/packages:/packages fedora:latest bash -c \
  "dnf install -y /packages/rpm/generator-VERSION-1.x86_64.rpm && generator version"

# Check package dependencies
dpkg -I packages/deb/generator_VERSION_amd64.deb | grep Depends
rpm -qpR packages/rpm/generator-VERSION-1.x86_64.rpm
```

### Distribution Issues

```bash
# Test download URLs
curl -I https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-linux-amd64.tar.gz

# Verify checksums
curl -sL https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/checksums.txt | sha256sum -c

# Test Docker image pull
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# Check Docker image layers
docker history ghcr.io/cuesoftinc/open-source-project-generator:latest

# Inspect Docker image
docker inspect ghcr.io/cuesoftinc/open-source-project-generator:latest
```

### Permission Issues

All Docker containers run as non-root user (UID 1001):

```bash
# Check user in container
docker run --rm --entrypoint sh generator:latest -c "id"
# Should show: uid=1001(generator) gid=1001(generator)

# Fix volume permissions on host
sudo chown -R 1001:1001 ./output ./config

# Or run with current user
docker run --rm --user $(id -u):$(id -g) generator:latest version
```

## Contributing

### Adding New Platforms

1. Update build scripts with new platform
2. Add platform-specific installation instructions
3. Test on target platform
4. Update documentation

### Adding Package Managers

1. Create package definition (formula, spec, etc.)
2. Submit to package repository
3. Update installation documentation
4. Add to CI/CD pipeline

### Container Registries

1. Configure registry credentials
2. Update build pipeline
3. Test image deployment
4. Update documentation

## Quick Reference

### Common Commands

```bash
# Build
make build              # Build for current platform
make build-all          # Build for all platforms
make clean              # Clean build artifacts

# Test
make test               # Run tests
make test-coverage      # Run tests with coverage
make lint               # Run linter
make security-scan      # Run security scans

# Package
make package-all        # Build all packages

# Docker
make docker-build       # Build production image
make docker-test        # Test Docker image

# CI/CD
make ci                 # Run CI pipeline
make pre-commit         # Run pre-commit checks
make pre-release        # Full pre-release validation

# Docker Compose
docker compose --profile production run --rm generator
docker compose --profile development run --rm generator-dev bash
docker compose --profile testing up generator-test
docker compose --profile build up generator-build-all
docker compose --profile lint up generator-lint
docker compose --profile security up generator-security
```

### File Reference

- **Makefile** - Build automation and commands
- **Dockerfile** - Production image (alpine:3.19, 39 MB)
- **Dockerfile.dev** - Development image (golang:1.25-alpine, ~500 MB)
- **Dockerfile.build** - Build image (ubuntu:24.04, ~1.5 GB)
- **docker-compose.yml** - Multi-profile orchestration
- **env.example** - Environment variable reference
- **go.mod** - Go dependencies (Go 1.25.0)
- **SECURITY.md** - Security practices and scanning
- **README.md** - Project overview and quick start

### Architecture Summary

**User IDs**: All Docker containers use UID 1001 for consistency

**Image Sizes**:

- Production: ~39 MB (static binary, Alpine)
- Development: ~500 MB (with dev tools)
- Build: ~1.5 GB (with package tools)

**Supported Platforms**:

- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

**Package Formats**:

- DEB (Debian/Ubuntu)
- RPM (Red Hat/CentOS/Fedora)
- Arch Linux

## Support

For distribution-related issues:

- 📖 [README](README.md) - Project overview and quick start
- 📖 [Installation Guide](docs/INSTALLATION.md) - Detailed installation instructions
- 🔒 [Security Guide](SECURITY.md) - Security practices and scanning
- 🐛 [Issue Tracker](https://github.com/cuesoftinc/open-source-project-generator/issues)
- 💬 [Discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- 📧 [Email Support](mailto:support@cuesoft.io)

## Contributing to Distribution

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Adding new platforms
- Creating package definitions
- Submitting to package repositories
- Improving build automation
- Testing distribution methods
