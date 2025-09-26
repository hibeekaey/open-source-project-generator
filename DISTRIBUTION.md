# Distribution Guide

This document describes the distribution and release process for the Open Source Project Generator.

## Overview

The Open Source Project Generator (v1.3.0+) supports multiple distribution methods:

- **Binary releases**: Cross-platform binaries for direct download
- **Package managers**: Native packages for Linux distributions
- **Container images**: Docker images for containerized environments
- **Source builds**: Build from source code
- **Go install**: Direct installation via Go toolchain

## Build System

### Cross-Platform Builds

The build system creates binaries for multiple platforms:

```bash
# Build all platforms
make dist

# Or use the build script directly
./scripts/build.sh
```

**Supported Platforms**:

- Linux: amd64, arm64, 386
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64, 386
- FreeBSD: amd64

**Build Requirements**:

- Go 1.25.0 or later
- Cross-compilation support for all target platforms

### Package Building

Create native packages for different distributions:

```bash
# Build all packages
make package-all

# Build specific package types
make package-deb    # Debian/Ubuntu
make package-rpm    # Red Hat/CentOS/Fedora
make package-arch   # Arch Linux
```

### Docker Images

Build container images:

```bash
# Production image
make docker-build

# Development image
docker-compose --profile development build

# Build image
docker-compose --profile build build
```

## Release Process

### Automated Releases

Releases are automated through GitHub Actions:

1. **Tag Creation**: Push a version tag (e.g., `v1.0.0`)
2. **Build Pipeline**: Automatically builds all platforms and packages
3. **Release Creation**: Creates GitHub release with all artifacts
4. **Docker Push**: Pushes images to Docker Hub
5. **Notifications**: Sends notifications on completion

### Manual Release

For manual releases:

```bash
# 1. Prepare release artifacts
make release-prepare

# 2. Create checksums
cd dist && sha256sum * > checksums.txt

# 3. Create GitHub release manually
# Upload files from dist/ and packages/ directories
```

### Version Management

Version information is embedded during build:

```bash
# Set version during build
VERSION=v1.3.0 ./scripts/build.sh

# Version is embedded in binary
./bin/generator --version
```

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

The Docker images use dynamic organization names based on the repository owner.

**Official Images (cuesoftinc):**

```bash
# Pull latest image
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# Pull specific version
docker pull ghcr.io/cuesoftinc/open-source-project-generator:v1.3.0
```

**Fork Images (automatically adapts to your GitHub username):**

```bash
# Pull from your fork (replace 'your-username')
docker pull ghcr.io/your-username/open-source-project-generator:latest

# Pull specific version from your fork
docker pull ghcr.io/your-username/open-source-project-generator:v1.3.0
```

**Dynamic Configuration:**
The build system automatically detects the repository owner and builds images accordingly. When you fork the repository and run GitHub Actions, images will be published to your own container registry namespace.

### Package Managers

#### Homebrew (macOS/Linux)

```bash
# Install via Homebrew (when available)
brew install generator
```

#### Chocolatey (Windows)

```powershell
# Install via Chocolatey (when available)
choco install generator
```

#### Scoop (Windows)

```powershell
# Install via Scoop (when available)
scoop install generator
```

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

- `VERSION`: Version to embed in binary
- `GOOS`: Target operating system
- `GOARCH`: Target architecture
- `CGO_ENABLED`: Enable/disable CGO (default: 0)

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

```bash
# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o generator-linux-arm64 ./cmd/generator

# Windows
GOOS=windows GOARCH=amd64 go build -o generator-windows-amd64.exe ./cmd/generator

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o generator-darwin-arm64 ./cmd/generator
```

## Quality Assurance

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Integration tests
make test-integration

# Test installation script
make test-install
```

### Validation

```bash
# Validate packages
dpkg -I packages/generator_VERSION_amd64.deb
rpm -qip packages/generator-VERSION-1.x86_64.rpm

# Test binaries
./dist/generator-linux-amd64/generator --version
./dist/generator-darwin-amd64/generator --version
```

### Security

```bash
# Scan for vulnerabilities
govulncheck ./...

# Check dependencies
go list -m -u all

# Verify checksums
sha256sum -c dist/checksums.txt
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

## Maintenance

### Regular Tasks

1. **Dependency Updates**: Update Go modules monthly
2. **Security Patches**: Apply security updates immediately
3. **Platform Testing**: Test on all supported platforms
4. **Documentation**: Keep installation guides updated

### Automation

- **Dependabot**: Automated dependency updates
- **Security Scanning**: CodeQL and vulnerability scanning
- **Build Testing**: Continuous integration on all platforms
- **Release Automation**: Automated releases on tag push

## Troubleshooting

### Build Issues

```bash
# Clean build cache
go clean -cache -modcache

# Rebuild dependencies
go mod download
go mod tidy

# Verbose build
go build -v -x ./cmd/generator
```

### Package Issues

```bash
# Validate package structure
dpkg-deb --contents generator_VERSION_amd64.deb
rpm2cpio generator-VERSION-1.x86_64.rpm | cpio -tv

# Test package installation
docker run --rm -v $(pwd):/packages ubuntu:22.04 bash -c "apt update && dpkg -i /packages/generator_VERSION_amd64.deb"
```

### Distribution Issues

```bash
# Test download URLs
curl -I https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-linux-amd64.tar.gz

# Verify checksums
curl -sL https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/checksums.txt | sha256sum -c
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

## Support

For distribution-related issues:

- 📖 [Installation Guide](docs/INSTALLATION.md)
- 🐛 [Issue Tracker](https://github.com/cuesoftinc/open-source-project-generator/issues)
- 💬 [Discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- 📧 [Email Support](mailto:support@generator.dev)
