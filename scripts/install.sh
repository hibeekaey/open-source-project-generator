#!/bin/bash

# Installation script for Open Source Project Generator
# This script automatically detects the platform and installs the appropriate binary

set -e

# Configuration
REPO_OWNER="cuesoftinc"
REPO_NAME="open-source-project-generator"
BINARY_NAME="generator"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="${HOME}/.local/bin"
MAX_RETRIES=3
RETRY_DELAY=2

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
cleanup() {
    local exit_code=$?
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
    if [ $exit_code -ne 0 ]; then
        print_error "Installation failed"
    fi
}
trap cleanup EXIT

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
        armv7l|armv7)
            ARCH="arm"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            print_status "Supported architectures: x86_64, arm64, armv7, i386"
            exit 1
            ;;
    esac
    
    print_status "Detected platform: $OS/$ARCH"
}

# Get latest release version
get_latest_version() {
    print_status "Fetching latest release information..."
    
    local api_url="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
    local response
    
    if command -v curl >/dev/null 2>&1; then
        response=$(curl -s -w "\n%{http_code}" "$api_url")
        local http_code=$(echo "$response" | tail -n1)
        local body=$(echo "$response" | sed '$d')
        
        if [ "$http_code" = "403" ]; then
            print_error "GitHub API rate limit exceeded"
            print_status "Please try again later or specify a version with --version"
            exit 1
        elif [ "$http_code" != "200" ]; then
            print_error "Failed to fetch release information (HTTP $http_code)"
            exit 1
        fi
        
        VERSION=$(echo "$body" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        response=$(wget -qO- "$api_url" 2>&1)
        if echo "$response" | grep -q "403"; then
            print_error "GitHub API rate limit exceeded"
            print_status "Please try again later or specify a version with --version"
            exit 1
        fi
        VERSION=$(echo "$response" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    if [ -z "$VERSION" ]; then
        print_error "Failed to fetch latest version"
        print_status "Please specify a version manually with --version"
        exit 1
    fi
    
    print_status "Latest version: $VERSION"
}

# Check if package manager installation is available
check_package_manager() {
    case $OS in
        linux)
            # Prefer newer package managers
            if command -v dnf >/dev/null 2>&1; then
                PACKAGE_MANAGER="dnf"
                return 0
            elif command -v apt-get >/dev/null 2>&1; then
                PACKAGE_MANAGER="apt"
                return 0
            elif command -v yum >/dev/null 2>&1; then
                PACKAGE_MANAGER="yum"
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
            package_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${BINARY_NAME}_${VERSION}_amd64.deb"
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
            package_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${BINARY_NAME}-${VERSION}-1.x86_64.rpm"
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

# Download file with retries
download_file() {
    local url=$1
    local output=$2
    local retries=0
    
    while [ $retries -lt $MAX_RETRIES ]; do
        if command -v curl >/dev/null 2>&1; then
            if curl -L --fail --progress-bar -o "$output" "$url"; then
                return 0
            fi
        elif command -v wget >/dev/null 2>&1; then
            if wget --show-progress -O "$output" "$url"; then
                return 0
            fi
        fi
        
        retries=$((retries + 1))
        if [ $retries -lt $MAX_RETRIES ]; then
            print_warning "Download failed, retrying in ${RETRY_DELAY}s... (attempt $retries/$MAX_RETRIES)"
            sleep $RETRY_DELAY
        fi
    done
    
    return 1
}

# Verify checksum
verify_checksum() {
    local file=$1
    local expected_checksum=$2
    
    if [ -z "$expected_checksum" ]; then
        print_warning "No checksum provided, skipping verification"
        return 0
    fi
    
    print_status "Verifying checksum..."
    
    local actual_checksum
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum=$(shasum -a 256 "$file" | cut -d' ' -f1)
    else
        print_warning "No checksum utility available, skipping verification"
        return 0
    fi
    
    if [ "$actual_checksum" = "$expected_checksum" ]; then
        print_success "Checksum verified"
        return 0
    else
        print_error "Checksum mismatch!"
        print_error "Expected: $expected_checksum"
        print_error "Got: $actual_checksum"
        return 1
    fi
}

# Download and install binary manually
install_binary() {
    local archive_name="${BINARY_NAME}-${OS}-${ARCH}"
    local archive_ext=".tar.gz"
    local download_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/${archive_name}${archive_ext}"
    local checksums_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/checksums.txt"
    
    TEMP_DIR=$(mktemp -d)
    
    print_status "Downloading $archive_name$archive_ext..."
    
    # Download archive with retries
    if ! download_file "$download_url" "$TEMP_DIR/${archive_name}${archive_ext}"; then
        print_error "Failed to download after $MAX_RETRIES attempts"
        print_status "URL: $download_url"
        exit 1
    fi
    
    # Download and verify checksum
    if download_file "$checksums_url" "$TEMP_DIR/checksums.txt" 2>/dev/null; then
        local expected_checksum=$(grep "${archive_name}${archive_ext}" "$TEMP_DIR/checksums.txt" | cut -d' ' -f1)
        if ! verify_checksum "$TEMP_DIR/${archive_name}${archive_ext}" "$expected_checksum"; then
            print_error "Checksum verification failed"
            exit 1
        fi
    else
        print_warning "Could not download checksums, skipping verification"
    fi
    
    # Extract archive
    print_status "Extracting archive..."
    if ! tar -xzf "$TEMP_DIR/${archive_name}${archive_ext}" -C "$TEMP_DIR"; then
        print_error "Failed to extract archive"
        exit 1
    fi
    
    # Verify binary exists
    if [ ! -f "$TEMP_DIR/$archive_name/$BINARY_NAME" ]; then
        print_error "Binary not found in archive"
        exit 1
    fi
    
    # Determine installation directory
    local install_dir
    local needs_sudo=false
    
    if [ -w "$INSTALL_DIR" ] || [ "$(id -u)" = "0" ]; then
        install_dir="$INSTALL_DIR"
        if [ "$(id -u)" != "0" ] && [ ! -w "$INSTALL_DIR" ]; then
            needs_sudo=true
            print_warning "Installing to $install_dir requires sudo privileges"
            print_status "You may be prompted for your password"
        fi
    else
        install_dir="$USER_INSTALL_DIR"
        mkdir -p "$install_dir"
        print_status "Installing to user directory: $install_dir"
    fi
    
    # Install binary
    print_status "Installing to $install_dir..."
    
    if [ "$needs_sudo" = true ]; then
        if ! sudo cp "$TEMP_DIR/$archive_name/$BINARY_NAME" "$install_dir/"; then
            print_error "Failed to copy binary (permission denied)"
            exit 1
        fi
        sudo chmod +x "$install_dir/$BINARY_NAME"
    else
        if ! cp "$TEMP_DIR/$archive_name/$BINARY_NAME" "$install_dir/"; then
            print_error "Failed to copy binary"
            exit 1
        fi
        chmod +x "$install_dir/$BINARY_NAME"
    fi
    
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
            if [ -f "$HOME/.bash_profile" ]; then
                shell_config="$HOME/.bash_profile"
            else
                shell_config="$HOME/.bashrc"
            fi
            ;;
        zsh)
            shell_config="$HOME/.zshrc"
            ;;
        fish)
            shell_config="$HOME/.config/fish/config.fish"
            mkdir -p "$(dirname "$shell_config")"
            ;;
        *)
            shell_config="$HOME/.profile"
            ;;
    esac
    
    # Check if PATH already contains the directory
    if echo ":$PATH:" | grep -q ":$USER_INSTALL_DIR:"; then
        print_status "$USER_INSTALL_DIR is already in PATH"
        return
    fi
    
    # Check if shell config already has the PATH entry
    if [ -f "$shell_config" ] && grep -q "$USER_INSTALL_DIR" "$shell_config"; then
        print_status "PATH entry already exists in $shell_config"
        print_warning "Please restart your shell or run: source $shell_config"
        return
    fi
    
    print_status "Adding $USER_INSTALL_DIR to PATH in $shell_config"
    
    case $current_shell in
        fish)
            echo "" >> "$shell_config"
            echo "# Added by generator installer" >> "$shell_config"
            echo "set -gx PATH $USER_INSTALL_DIR \$PATH" >> "$shell_config"
            ;;
        *)
            echo "" >> "$shell_config"
            echo "# Added by generator installer" >> "$shell_config"
            echo "export PATH=\"$USER_INSTALL_DIR:\$PATH\"" >> "$shell_config"
            ;;
    esac
    
    print_warning "Please restart your shell or run: source $shell_config"
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
    cat << EOF
Installation Script

Usage: $0 [OPTIONS]

Description:
  Automatically detects the platform and installs the generator binary.
  Supports package managers (apt, yum, dnf) and direct binary installation.

Options:
  --version VERSION   Install specific version (e.g., v1.0.0)
  --force-binary      Force binary installation (skip package managers)
  --user              Install to user directory (~/.local/bin)
  -h, --help          Show this help message

Environment Variables:
  INSTALL_DIR              Custom installation directory
  GENERATOR_SKIP_VERIFY    Skip installation verification

Examples:
  $0                        # Install latest version
  $0 --version v1.0.0       # Install specific version
  $0 --force-binary         # Skip package manager, install binary
  $0 --user                 # Install to user directory

Requirements:
  - curl or wget
  - tar

Output:
  Binary installed to: /usr/local/bin/ or ~/.local/bin/
EOF
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
install_generator() {
    echo "Open Source Project Generator Installation Script"
    echo "================================================="
    echo ""
    
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
    echo "To uninstall:"
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        echo "  sudo rm $INSTALL_DIR/$BINARY_NAME"
    elif [ -f "$USER_INSTALL_DIR/$BINARY_NAME" ]; then
        echo "  rm $USER_INSTALL_DIR/$BINARY_NAME"
    fi
    echo ""
}

# Parse command line arguments
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
        "help"|"-h"|"--help")
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

# Run installation
install_generator