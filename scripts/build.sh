#!/bin/bash

# Build script for cross-platform binary compilation
# This script builds the template generator for multiple platforms

set -e

# Configuration
APP_NAME="generator"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="dist"
MAIN_PACKAGE="./cmd/generator"

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

# Clean previous builds
clean_build() {
    print_status "Cleaning previous builds..."
    if [ -d "${BUILD_DIR}" ]; then
        # Try to remove contents but not the directory itself (in case it's a volume mount)
        rm -rf ${BUILD_DIR}/* 2>/dev/null || true
    fi
    mkdir -p ${BUILD_DIR}
}

# Get git commit hash for build info
get_build_info() {
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    
    LDFLAGS="-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"
    export LDFLAGS
}

# Build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    
    local output_name="${APP_NAME}"
    if [ ! -z "$ext" ]; then
        output_name="${output_name}${ext}"
    fi
    
    local output_path="${BUILD_DIR}/${APP_NAME}-${os}-${arch}/${output_name}"
    
    print_status "Building for ${os}/${arch}..."
    
    # Create platform directory
    mkdir -p "${BUILD_DIR}/${APP_NAME}-${os}-${arch}"
    
    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0
    
    # Build the binary with optimizations
    if go build -ldflags "${LDFLAGS} -s -w" -trimpath -o "${output_path}" ${MAIN_PACKAGE}; then
        print_success "Built ${output_path}"
        
        # Copy additional files
        cp README.md "${BUILD_DIR}/${APP_NAME}-${os}-${arch}/"
        cp LICENSE "${BUILD_DIR}/${APP_NAME}-${os}-${arch}/" 2>/dev/null || true
        
        # Create platform-specific installation instructions
        create_install_instructions "$os" "${BUILD_DIR}/${APP_NAME}-${os}-${arch}"
        
        # Create archive
        create_archive "$os" "$arch"
    else
        print_error "Failed to build for ${os}/${arch}"
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
    
    print_status "Creating archive for ${os}/${arch}..."
    
    cd ${BUILD_DIR}
    
    case $os in
        "windows")
            if command -v zip >/dev/null 2>&1; then
                zip -r "${dir_name}.zip" "${dir_name}/"
                print_success "Created ${dir_name}.zip"
            else
                print_warning "zip not available, skipping archive creation for Windows"
            fi
            ;;
        *)
            tar -czf "${dir_name}.tar.gz" "${dir_name}/"
            print_success "Created ${dir_name}.tar.gz"
            ;;
    esac
    
    cd ..
}

# Generate checksums
generate_checksums() {
    print_status "Generating checksums..."
    
    cd ${BUILD_DIR}
    
    # Generate SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt || true
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || true
    else
        print_warning "No checksum utility available"
    fi
    
    if [ -f checksums.txt ]; then
        print_success "Generated checksums.txt"
    fi
    
    cd ..
}

# Main build function
main() {
    print_status "Starting cross-platform build for ${APP_NAME} v${VERSION}"
    
    # Check if Go is installed
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Using Go version: ${GO_VERSION}"
    
    # Clean and prepare
    clean_build
    get_build_info
    
    # Build for different platforms
    print_status "Building binaries for multiple platforms..."
    
    # Linux
    build_platform "linux" "amd64" ""
    build_platform "linux" "arm64" ""
    build_platform "linux" "386" ""
    
    # macOS
    build_platform "darwin" "amd64" ""
    build_platform "darwin" "arm64" ""
    
    # Windows
    build_platform "windows" "amd64" ".exe"
    build_platform "windows" "386" ".exe"
    
    # FreeBSD
    build_platform "freebsd" "amd64" ""
    
    # Generate checksums
    generate_checksums
    
    # Show build summary
    print_success "Build completed successfully!"
    print_status "Build artifacts:"
    ls -la ${BUILD_DIR}/
    
    print_status "Total size of all artifacts:"
    du -sh ${BUILD_DIR}/
}

# Handle command line arguments
case "${1:-}" in
    "clean")
        clean_build
        print_success "Build directory cleaned"
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [clean|help]"
        echo ""
        echo "Commands:"
        echo "  clean    Clean build directory"
        echo "  help     Show this help message"
        echo ""
        echo "Environment variables:"
        echo "  VERSION  Set the version number (default: 1.0.0)"
        echo ""
        echo "Examples:"
        echo "  $0                    # Build all platforms"
        echo "  VERSION=1.2.0 $0     # Build with specific version"
        echo "  $0 clean             # Clean build directory"
        ;;
    *)
        main
        ;;
esac