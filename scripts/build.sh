#!/bin/bash

# Build script for cross-platform binary compilation
# This script builds the template generator for multiple platforms

set -e

# Configuration
APP_NAME="generator"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Get version from git tags, fallback to "dev" if not in a git repo
# Can be overridden with VERSION environment variable
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_DIR="dist"
MAIN_PACKAGE="./cmd/generator"

# Track build results
declare -a SUCCESSFUL_BUILDS=()
declare -a FAILED_BUILDS=()

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Cleanup trap
cleanup_on_error() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        print_error "Build failed with exit code $exit_code"
        if [ ${#FAILED_BUILDS[@]} -gt 0 ]; then
            print_error "Failed builds: ${FAILED_BUILDS[*]}"
        fi
    fi
}
trap cleanup_on_error EXIT

# Clean previous builds
clean_build() {
    print_status "Cleaning previous builds..."
    if [ -d "${BUILD_DIR}" ]; then
        # Try to remove contents but not the directory itself (in case it's a volume mount)
        rm -rf ${BUILD_DIR}/* 2>/dev/null || true
    fi
    mkdir -p ${BUILD_DIR}
}

# Validate version format
validate_version() {
    # Accept version as-is, whether it has 'v' prefix or not
    if ! [[ "$VERSION" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.+-]+)?$ ]] && [ "$VERSION" != "dev" ]; then
        print_error "Invalid version format: $VERSION"
        print_status "Version must follow semantic versioning (e.g., 1.0.0 or v1.0.0)"
        exit 1
    fi
}

# Get git commit hash for build info
get_build_info() {
    local git_commit="unknown"
    local git_branch="unknown"
    
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git_commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
        git_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
    fi
    
    GIT_COMMIT=$git_commit
    GIT_BRANCH=$git_branch
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    
    LDFLAGS="-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"
    export LDFLAGS
    
    print_status "Build info: version=${VERSION}, commit=${GIT_COMMIT}, branch=${GIT_BRANCH}"
}

# Build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    local platform_name="${os}/${arch}"
    
    local output_name="${APP_NAME}"
    if [ -n "$ext" ]; then
        output_name="${output_name}${ext}"
    fi
    
    local platform_dir="${BUILD_DIR}/${APP_NAME}-${os}-${arch}"
    local output_path="${platform_dir}/${output_name}"
    
    print_status "Building for ${platform_name}..."
    
    # Create platform directory
    mkdir -p "${platform_dir}"
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0
    
    # Build the binary with optimizations
    if go build -ldflags "${LDFLAGS} -s -w" -trimpath -o "${output_path}" ${MAIN_PACKAGE} 2>&1; then
        # Verify the binary was created and is executable
        if [ ! -f "${output_path}" ]; then
            print_error "Binary not created for ${platform_name}"
            FAILED_BUILDS+=("${platform_name}")
            return 1
        fi
        
        # Get binary size
        local size=$(du -h "${output_path}" | cut -f1)
        print_success "Built ${platform_name} (${size})"
        
        # Copy additional files if they exist
        [ -f "README.md" ] && cp README.md "${platform_dir}/" || print_warning "README.md not found"
        [ -f "LICENSE" ] && cp LICENSE "${platform_dir}/" || print_warning "LICENSE not found"
        
        # Create platform-specific installation instructions
        create_install_instructions "$os" "${platform_dir}"
        
        # Create archive
        if create_archive "$os" "$arch"; then
            SUCCESSFUL_BUILDS+=("${platform_name}")
        else
            print_warning "Archive creation failed for ${platform_name}"
            SUCCESSFUL_BUILDS+=("${platform_name} (no archive)")
        fi
    else
        print_error "Failed to build for ${platform_name}"
        FAILED_BUILDS+=("${platform_name}")
        return 1
    fi
}

# Create installation instructions for each platform
create_install_instructions() {
    local os=$1
    local dir=$2
    
    case $os in
        "windows")
            cat > "${dir}/INSTALL.txt" << 'EOF'
Installation Instructions for Windows
====================================

1. Extract the archive to a directory (e.g., C:\Program Files\generator)
2. Add the directory to your PATH environment variable:
   - Open System Properties > Advanced > Environment Variables
   - Edit the PATH variable and add the generator directory
   - Or run: setx PATH "%PATH%;C:\Program Files\generator"
3. Open a new command prompt and run: generator --help

Alternative: Place generator.exe in any directory that's already in your PATH.

Note: The generator is completely self-contained with all templates embedded.
EOF
            ;;
        "darwin")
            cat > "${dir}/INSTALL.txt" << 'EOF'
Installation Instructions for macOS
==================================

Option 1: System Installation
1. Extract the archive
2. Copy to system directory: sudo cp generator /usr/local/bin/
3. Make executable: sudo chmod +x /usr/local/bin/generator
4. Run: generator --help

Option 2: User Installation
1. Create ~/bin directory: mkdir -p ~/bin
2. Copy generator to ~/bin/
3. Make executable: chmod +x ~/bin/generator
4. Add to PATH: echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
5. Reload shell: source ~/.zshrc
6. Run: generator --help

Note: The generator is completely self-contained with all templates embedded.
EOF
            ;;
        "linux")
            cat > "${dir}/INSTALL.txt" << 'EOF'
Installation Instructions for Linux
==================================

Option 1: System-wide Installation (requires sudo)
1. Extract the archive
2. Copy to system directory: sudo cp generator /usr/local/bin/
3. Make executable: sudo chmod +x /usr/local/bin/generator
4. Run: generator --help

Option 2: User Installation
1. Create ~/bin directory: mkdir -p ~/bin
2. Copy generator to ~/bin/
3. Make executable: chmod +x ~/bin/generator
4. Add to PATH: echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
5. Reload shell: source ~/.bashrc
6. Run: generator --help

Option 3: Using package managers
- For Debian/Ubuntu: Create a .deb package (see build-packages.sh)
- For Red Hat/CentOS: Create an .rpm package (see build-packages.sh)
- For Arch Linux: Create a PKGBUILD (see build-packages.sh)

Note: The generator is completely self-contained with all templates embedded.
EOF
            ;;
    esac
}

# Create archive for distribution
create_archive() {
    local os=$1
    local arch=$2
    local dir_name="${APP_NAME}-${os}-${arch}"
    local current_dir=$(pwd)
    
    print_status "Creating archive for ${os}/${arch}..."
    
    cd "${BUILD_DIR}" || return 1
    
    case $os in
        "windows")
            if command -v zip >/dev/null 2>&1; then
                if zip -q -r "${dir_name}.zip" "${dir_name}/"; then
                    print_success "Created ${dir_name}.zip"
                    cd "${current_dir}"
                    return 0
                else
                    print_error "Failed to create zip archive"
                    cd "${current_dir}"
                    return 1
                fi
            else
                print_warning "zip not available, skipping archive creation for Windows"
                cd "${current_dir}"
                return 1
            fi
            ;;
        *)
            if tar -czf "${dir_name}.tar.gz" "${dir_name}/"; then
                print_success "Created ${dir_name}.tar.gz"
                cd "${current_dir}"
                return 0
            else
                print_error "Failed to create tar.gz archive"
                cd "${current_dir}"
                return 1
            fi
            ;;
    esac
}

# Generate checksums
generate_checksums() {
    print_status "Generating checksums..."
    
    local current_dir=$(pwd)
    cd "${BUILD_DIR}" || return 1
    
    # Find all archives
    local archives=$(find . -maxdepth 1 \( -name "*.tar.gz" -o -name "*.zip" \) 2>/dev/null)
    
    if [ -z "$archives" ]; then
        print_warning "No archives found for checksum generation"
        cd "${current_dir}"
        return 1
    fi
    
    # Generate SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        find . -maxdepth 1 \( -name "*.tar.gz" -o -name "*.zip" \) -exec sha256sum {} \; > checksums.txt
        print_success "Generated checksums.txt using sha256sum"
    elif command -v shasum >/dev/null 2>&1; then
        find . -maxdepth 1 \( -name "*.tar.gz" -o -name "*.zip" \) -exec shasum -a 256 {} \; > checksums.txt
        print_success "Generated checksums.txt using shasum"
    else
        print_warning "No checksum utility available (sha256sum or shasum)"
        cd "${current_dir}"
        return 1
    fi
    
    cd "${current_dir}"
}

# Main build function
main() {
    print_status "Starting cross-platform build for ${APP_NAME} ${VERSION}"
    
    # Validate version
    validate_version
    
    # Check if Go is installed
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version (require 1.25+)
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Using Go version: ${GO_VERSION}"
    
    # Verify we're in a Go module
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. This script must be run from the project root."
        exit 1
    fi
    
    # Verify main package exists
    if [ ! -d "cmd/generator" ]; then
        print_error "Main package not found at cmd/generator"
        exit 1
    fi
    
    # Clean and prepare
    clean_build
    get_build_info
    
    # Build for different platforms
    print_status "Building binaries for multiple platforms..."
    
    # Linux (most common architectures)
    build_platform "linux" "amd64" "" || true
    build_platform "linux" "arm64" "" || true
    build_platform "linux" "arm" "" || true  # ARMv7
    
    # macOS (Intel and Apple Silicon)
    build_platform "darwin" "amd64" "" || true
    build_platform "darwin" "arm64" "" || true
    
    # Windows (64-bit and 32-bit)
    build_platform "windows" "amd64" ".exe" || true
    build_platform "windows" "386" ".exe" || true
    
    # FreeBSD
    build_platform "freebsd" "amd64" "" || true
    
    # Generate checksums
    generate_checksums || print_warning "Checksum generation failed"
    
    # Show build summary
    echo ""
    print_status "========================================="
    print_status "Build Summary"
    print_status "========================================="
    
    if [ ${#SUCCESSFUL_BUILDS[@]} -gt 0 ]; then
        print_success "Successful builds (${#SUCCESSFUL_BUILDS[@]}):"
        for build in "${SUCCESSFUL_BUILDS[@]}"; do
            echo "  ✓ $build"
        done
    fi
    
    if [ ${#FAILED_BUILDS[@]} -gt 0 ]; then
        echo ""
        print_error "Failed builds (${#FAILED_BUILDS[@]}):"
        for build in "${FAILED_BUILDS[@]}"; do
            echo "  ✗ $build"
        done
        echo ""
        print_warning "Some builds failed, but continuing..."
    fi
    
    echo ""
    print_status "Build artifacts in ${BUILD_DIR}/:"
    ls -lh ${BUILD_DIR}/*.tar.gz ${BUILD_DIR}/*.zip 2>/dev/null || print_warning "No archives created"
    
    if [ -f "${BUILD_DIR}/checksums.txt" ]; then
        echo ""
        print_status "Checksums:"
        cat "${BUILD_DIR}/checksums.txt"
    fi
    
    echo ""
    print_status "Total size of all artifacts:"
    du -sh ${BUILD_DIR}/ 2>/dev/null || true
    
    # Exit with error if all builds failed
    if [ ${#SUCCESSFUL_BUILDS[@]} -eq 0 ]; then
        print_error "All builds failed!"
        exit 1
    fi
    
    print_success "Build completed with ${#SUCCESSFUL_BUILDS[@]} successful platform(s)!"
}

# Show usage
show_usage() {
    cat << EOF
Cross-Platform Build Script

Usage: $0 [COMMAND]

Description:
  Builds the generator binary for multiple platforms and architectures.

Commands:
  (none)          Build all platforms (default)
  clean           Clean build directory
  help            Show this help message

Environment Variables:
  VERSION         Set the version number (default: 1.0.0)
                  Must follow semantic versioning (e.g., 1.0.0 or v1.0.0)

Supported Platforms:
  - Linux: amd64, arm64, arm (ARMv7)
  - macOS: amd64 (Intel), arm64 (Apple Silicon)
  - Windows: amd64, 386
  - FreeBSD: amd64

Examples:
  $0                        # Build all platforms
  VERSION=1.2.0 $0          # Build with specific version
  $0 clean                  # Clean build directory

Requirements:
  - Go 1.25+

Output:
  Binaries: dist/generator-{os}-{arch}/
  Archives: dist/generator-{os}-{arch}.tar.gz (or .zip for Windows)
  Checksums: dist/checksums.txt
EOF
}

# Handle command line arguments
case "${1:-}" in
    "clean")
        clean_build
        print_success "Build directory cleaned"
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        main
        ;;
esac