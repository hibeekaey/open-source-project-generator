#!/bin/bash

# Installation script for Open Source Project Generator
# This script automatically detects the platform and installs the appropriate binary

set -e

# Configuration
REPO_OWNER="cuesoftinc"
REPO_NAME="open-source-project-generator"
BINARY_NAME="generator"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/bin"

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

# Detect platform
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        freebsd*)
            OS="freebsd"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    print_status "Detected platform: $OS/$ARCH"
}

# Get latest release version
get_latest_version() {
    print_status "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi
    
    print_status "Latest version: $VERSION"
}

# Check if package manager installation is available
check_package_manager() {
    case $OS in
        linux)
            if command -v apt-get >/dev/null 2>&1; then
                PACKAGE_MANAGER="apt"
                return 0
            elif command -v yum >/dev/null 2>&1; then
                PACKAGE_MANAGER="yum"
                return 0
            elif command -v dnf >/dev/null 2>&1; then
                PACKAGE_MANAGER="dnf"
                return 0
            elif command -v pacman >/dev/null 2>&1; then
                PACKAGE_MANAGER="pacman"
                return 0
            fi
            ;;
        darwin)
            if command -v brew >/dev/null 2>&1; then
                PACKAGE_MANAGER="brew"
                return 0
            fi
            ;;
    esac
    
    return 1
}

# Install via package manager
install_via_package_manager() {
    local package_url
    
    case $PACKAGE_MANAGER in
        apt)
            package_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${BINARY_NAME}_${VERSION#v}_amd64.deb"
            print_status "Installing via APT..."
            
            # Download package
            if command -v curl >/dev/null 2>&1; then
                curl -L -o "/tmp/${BINARY_NAME}.deb" "$package_url"
            else
                wget -O "/tmp/${BINARY_NAME}.deb" "$package_url"
            fi
            
            # Install package
            sudo dpkg -i "/tmp/${BINARY_NAME}.deb" || sudo apt-get install -f -y
            rm -f "/tmp/${BINARY_NAME}.deb"
            ;;
            
        yum|dnf)
            package_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${BINARY_NAME}-${VERSION#v}-1.x86_64.rpm"
            print_status "Installing via $PACKAGE_MANAGER..."
            
            # Install package directly from URL
            sudo $PACKAGE_MANAGER install -y "$package_url"
            ;;
            
        pacman)
            print_warning "Arch Linux package installation not yet supported via this script"
            print_status "Please install manually or use the AUR package when available"
            return 1
            ;;
            
        brew)
            print_warning "Homebrew formula not yet available"
            print_status "Installing manually..."
            return 1
            ;;
    esac
}

# Download and install binary manually
install_binary() {
    local archive_name="${BINARY_NAME}-${OS}-${ARCH}"
    local archive_ext=".tar.gz"
    local download_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${archive_name}${archive_ext}"
    local temp_dir=$(mktemp -d)
    
    print_status "Downloading $archive_name$archive_ext..."
    
    # Download archive
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_dir/${archive_name}${archive_ext}" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_dir/${archive_name}${archive_ext}" "$download_url"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi
    
    # Extract archive
    print_status "Extracting archive..."
    tar -xzf "$temp_dir/${archive_name}${archive_ext}" -C "$temp_dir"
    
    # Determine installation directory
    local install_dir
    local needs_sudo=false
    
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" = "0" ]; then
        install_dir="$INSTALL_DIR"
        if [ "$(id -u)" != "0" ] && [ ! -w "$INSTALL_DIR" ]; then
            needs_sudo=true
        fi
    else
        install_dir="$USER_INSTALL_DIR"
        mkdir -p "$install_dir"
        print_warning "Installing to user directory: $install_dir"
        print_status "Make sure $install_dir is in your PATH"
    fi
    
    # Install binary
    print_status "Installing to $install_dir..."
    
    if [ "$needs_sudo" = true ]; then
        sudo cp "$temp_dir/$archive_name/$BINARY_NAME" "$install_dir/"
        sudo chmod +x "$install_dir/$BINARY_NAME"
    else
        cp "$temp_dir/$archive_name/$BINARY_NAME" "$install_dir/"
        chmod +x "$install_dir/$BINARY_NAME"
    fi
    
    # Clean up
    rm -rf "$temp_dir"
    
    # Update PATH if needed
    if [ "$install_dir" = "$USER_INSTALL_DIR" ]; then
        update_path
    fi
}

# Update PATH in shell configuration
update_path() {
    local shell_config=""
    local current_shell=$(basename "$SHELL")
    
    case $current_shell in
        bash)
            shell_config="$HOME/.bashrc"
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            ;;
        fish)
            shell_config="$HOME/.config/fish/config.fish"
            ;;
        *)
            shell_config="$HOME/.profile"
            ;;
    esac
    
    # Check if PATH update is needed
    if ! echo "$PATH" | grep -q "$USER_INSTALL_DIR"; then
        print_status "Adding $USER_INSTALL_DIR to PATH in $shell_config"
        
        case $current_shell in
            fish)
                echo "set -gx PATH $USER_INSTALL_DIR \$PATH" >> "$shell_config"
                ;;
            *)
                echo "export PATH=\"$USER_INSTALL_DIR:\$PATH\"" >> "$shell_config"
                ;;
        esac
        
        print_warning "Please restart your shell or run: source $shell_config"
    fi
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version=$($BINARY_NAME --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation successful!"
        print_status "Installed version: $installed_version"
        print_status "Run '$BINARY_NAME --help' to get started"
    else
        print_error "Installation verification failed"
        print_status "The binary may not be in your PATH"
        print_status "Try running: export PATH=\"$USER_INSTALL_DIR:\$PATH\""
        exit 1
    fi
}

# Show usage
show_usage() {
    echo "Open Source Project Generator Installation Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --version VERSION    Install specific version (e.g., v1.0.0)"
    echo "  --force-binary       Force binary installation (skip package managers)"
    echo "  --user               Install to user directory (~/.local/bin)"
    echo "  --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                           # Install latest version"
    echo "  $0 --version v1.0.0          # Install specific version"
    echo "  $0 --force-binary            # Skip package manager, install binary"
    echo "  $0 --user                    # Install to user directory"
    echo ""
    echo "Environment variables:"
    echo "  INSTALL_DIR              Custom installation directory"
    echo "  GENERATOR_SKIP_VERIFY    Skip installation verification"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --force-binary)
                FORCE_BINARY=true
                shift
                ;;
            --user)
                INSTALL_DIR="$USER_INSTALL_DIR"
                shift
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Check system requirements
check_requirements() {
    print_status "Checking system requirements..."
    
    # Check for required tools
    if ! command -v tar >/dev/null 2>&1; then
        print_error "tar is required but not installed"
        exit 1
    fi
    
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        print_error "Either curl or wget is required but neither is installed"
        exit 1
    fi
    
    print_success "System requirements satisfied"
}

# Main installation function
main() {
    echo "Open Source Project Generator Installation Script"
    echo "================================================="
    echo ""
    
    # Parse arguments
    parse_args "$@"
    
    # Check requirements
    check_requirements
    
    # Detect platform
    detect_platform
    
    # Get version if not specified
    if [ -z "$VERSION" ]; then
        get_latest_version
    fi
    
    # Try package manager installation first (unless forced to use binary)
    if [ "$FORCE_BINARY" != true ] && check_package_manager; then
        print_status "Package manager detected: $PACKAGE_MANAGER"
        
        if install_via_package_manager; then
            print_success "Installed via package manager"
        else
            print_warning "Package manager installation failed, falling back to binary installation"
            install_binary
        fi
    else
        print_status "Installing binary directly"
        install_binary
    fi
    
    # Verify installation unless skipped
    if [ "$GENERATOR_SKIP_VERIFY" != true ]; then
        verify_installation
    fi
    
    echo ""
    print_success "Installation completed!"
    echo ""
    echo "Next steps:"
    echo "  1. Run 'generator --help' to see available commands"
    echo "  2. Run 'generator generate' to create your first project"
    echo "  3. Visit https://github.com/$REPO_OWNER/$REPO_NAME for documentation"
    echo ""
}

# Run main function
main "$@"