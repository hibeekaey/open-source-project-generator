#!/bin/bash

# Package build script for creating distribution packages
# Supports DEB, RPM, and other package formats

set -e

# Configuration
APP_NAME="generator"
APP_DESCRIPTION="Open Source Project Generator - Create production-ready project structures"
APP_URL="https://github.com/cuesoftinc/open-source-project-generator"
MAINTAINER="Open Source Project Generator Team <team@example.com>"
LICENSE="MIT"
# Get version from git tags, fallback to "dev" if not in a git repo
# Can be overridden with VERSION environment variable
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}

# Detect architecture
SYSTEM_ARCH=$(uname -m)
case "$SYSTEM_ARCH" in
    x86_64)
        ARCH="amd64"
        RPM_ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        RPM_ARCH="aarch64"
        ;;
    armv7l)
        ARCH="armhf"
        RPM_ARCH="armv7hl"
        ;;
    *)
        ARCH="amd64"
        RPM_ARCH="x86_64"
        ;;
esac

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

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DIST_DIR="${PROJECT_ROOT}/dist"
PACKAGES_DIR="${PROJECT_ROOT}/packages"
BINARY_PATH="${DIST_DIR}/generator-linux-${ARCH}/generator"

# Cleanup trap
cleanup_on_error() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        print_error "Build failed with exit code $exit_code"
        print_status "Cleaning up temporary files..."
        # Keep packages directory for inspection but note the failure
        if [ -d "$PACKAGES_DIR" ]; then
            touch "${PACKAGES_DIR}/.build-failed"
        fi
    fi
}
trap cleanup_on_error EXIT

# Check dependencies
check_dependencies() {
    local package_type=$1
    
    case $package_type in
        "deb")
            if ! command -v dpkg-deb >/dev/null 2>&1; then
                print_error "dpkg-deb is required for building DEB packages"
                print_status "Install with: sudo apt-get install dpkg-dev"
                return 1
            fi
            ;;
        "rpm")
            if ! command -v rpmbuild >/dev/null 2>&1; then
                print_error "rpmbuild is required for building RPM packages"
                print_status "Install with: sudo apt-get install rpm (Ubuntu) or sudo yum install rpm-build (CentOS)"
                return 1
            fi
            ;;
        "arch")
            if ! command -v makepkg >/dev/null 2>&1; then
                print_error "makepkg is required for building Arch packages"
                print_status "This should be available on Arch Linux systems"
                return 1
            fi
            ;;
    esac
    
    return 0
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

# Prepare package directory
prepare_package_dir() {
    print_status "Preparing package directory..."
    
    # Validate version
    validate_version
    
    # Clean and create packages directory
    rm -rf "$PACKAGES_DIR"
    mkdir -p "$PACKAGES_DIR"
    
    # Check if binary exists
    if [ ! -f "$BINARY_PATH" ]; then
        print_error "Binary not found at $BINARY_PATH"
        print_status "Please run the build script first to create the binary"
        print_status "Expected location: $BINARY_PATH"
        exit 1
    fi
    
    print_status "Using binary: $BINARY_PATH"
    print_status "Architecture: $ARCH"
}

# Build DEB package
build_deb() {
    print_status "Building DEB package..."
    
    local deb_dir="${PACKAGES_DIR}/deb"
    local control_dir="${deb_dir}/DEBIAN"
    local bin_dir="${deb_dir}/usr/local/bin"
    local doc_dir="${deb_dir}/usr/share/doc/${APP_NAME}"
    local man_dir="${deb_dir}/usr/share/man/man1"
    
    # Create directory structure
    mkdir -p "$control_dir" "$bin_dir" "$doc_dir" "$man_dir"
    
    # Copy binary
    cp "$BINARY_PATH" "$bin_dir/"
    chmod 755 "${bin_dir}/generator"
    
    # Copy documentation
    cp "${PROJECT_ROOT}/README.md" "$doc_dir/"
    [ -f "${PROJECT_ROOT}/LICENSE" ] && cp "${PROJECT_ROOT}/LICENSE" "$doc_dir/"
    
    # Create control file
    cat > "${control_dir}/control" << EOF
Package: ${APP_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: ${APP_DESCRIPTION}
 The Open Source Project Generator creates production-ready project structures
 with modern best practices, including frontend applications, backend services,
 mobile applications, infrastructure code, and comprehensive documentation.
 .
 Features:
  - Interactive project generation with component selection
  - Support for Next.js, Go, Android, and iOS applications
  - Infrastructure as code templates (Docker, Kubernetes, Terraform)
  - Complete CI/CD workflows with GitHub Actions
  - Cross-platform support and package management integration
Homepage: ${APP_URL}
EOF

    # Create postinst script
    cat > "${control_dir}/postinst" << 'EOF'
#!/bin/bash
set -e

# Create symlink for global access
if [ ! -L "/usr/bin/generator" ]; then
    ln -sf /usr/local/bin/generator /usr/bin/generator
fi

echo "Generator installed successfully!"
echo "Run 'generator --help' to get started."
EOF

    # Create prerm script
    cat > "${control_dir}/prerm" << 'EOF'
#!/bin/bash
set -e

# Remove symlink
if [ -L "/usr/bin/generator" ]; then
    rm -f /usr/bin/generator
fi
EOF

    # Make scripts executable
    chmod 755 "${control_dir}/postinst" "${control_dir}/prerm"
    
    # Create man page
    create_man_page "${man_dir}/generator.1"
    gzip "${man_dir}/generator.1"
    
    # Build package
    dpkg-deb --build "$deb_dir" "${PACKAGES_DIR}/${APP_NAME}_${VERSION}_${ARCH}.deb"
    
    print_success "DEB package created: ${PACKAGES_DIR}/${APP_NAME}_${VERSION}_${ARCH}.deb"
}

# Build RPM package
build_rpm() {
    print_status "Building RPM package..."
    
    local rpm_dir="${PACKAGES_DIR}/rpm"
    local spec_file="${rpm_dir}/${APP_NAME}.spec"
    local sources_dir="${rpm_dir}/SOURCES"
    
    # Create RPM build directories
    mkdir -p "${rpm_dir}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS,BUILDROOT}
    
    # Copy source files to SOURCES
    mkdir -p "${sources_dir}/${APP_NAME}-${VERSION}"
    cp "$BINARY_PATH" "${sources_dir}/${APP_NAME}-${VERSION}/generator"
    cp "${PROJECT_ROOT}/README.md" "${sources_dir}/${APP_NAME}-${VERSION}/"
    [ -f "${PROJECT_ROOT}/LICENSE" ] && cp "${PROJECT_ROOT}/LICENSE" "${sources_dir}/${APP_NAME}-${VERSION}/"
    
    # Create man page in sources
    create_man_page "${sources_dir}/${APP_NAME}-${VERSION}/generator.1"
    gzip "${sources_dir}/${APP_NAME}-${VERSION}/generator.1"
    
    # Create tarball
    cd "${sources_dir}"
    tar -czf "${APP_NAME}-${VERSION}.tar.gz" "${APP_NAME}-${VERSION}"
    rm -rf "${APP_NAME}-${VERSION}"
    cd - >/dev/null
    
    # Create spec file
    cat > "$spec_file" << EOF
Name:           ${APP_NAME}
Version:        ${VERSION}
Release:        1%{?dist}
Summary:        ${APP_DESCRIPTION}

License:        ${LICENSE}
URL:            ${APP_URL}
Source0:        %{name}-%{version}.tar.gz
BuildArch:      ${RPM_ARCH}

%description
The Open Source Project Generator creates production-ready project structures
with modern best practices, including frontend applications, backend services,
mobile applications, infrastructure code, and comprehensive documentation.

Features:
- Interactive project generation with component selection
- Support for Next.js, Go, Android, and iOS applications
- Infrastructure as code templates (Docker, Kubernetes, Terraform)
- Complete CI/CD workflows with GitHub Actions
- Cross-platform support and package management integration

%prep
%setup -q

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/usr/share/doc/%{name}
mkdir -p %{buildroot}/usr/share/man/man1

install -m 755 generator %{buildroot}/usr/local/bin/
install -m 644 README.md %{buildroot}/usr/share/doc/%{name}/
[ -f LICENSE ] && install -m 644 LICENSE %{buildroot}/usr/share/doc/%{name}/
install -m 644 generator.1.gz %{buildroot}/usr/share/man/man1/

%files
%defattr(-,root,root,-)
/usr/local/bin/generator
%doc /usr/share/doc/%{name}/README.md
%{_mandir}/man1/generator.1.gz
%if 0%{?rhel} >= 7 || 0%{?fedora}
%license /usr/share/doc/%{name}/LICENSE
%else
%doc /usr/share/doc/%{name}/LICENSE
%endif

%post
# Create symlink for global access
if [ ! -L "/usr/bin/generator" ]; then
    ln -sf /usr/local/bin/generator /usr/bin/generator
fi

echo "Generator installed successfully!"
echo "Run 'generator --help' to get started."

%preun
# Remove symlink
if [ -L "/usr/bin/generator" ]; then
    rm -f /usr/bin/generator
fi

%changelog
* $(date +'%a %b %d %Y') ${MAINTAINER} - ${VERSION}-1
- Initial package release
EOF

    # Build RPM
    rpmbuild --define "_topdir ${rpm_dir}" \
             --define "_builddir ${rpm_dir}/BUILD" \
             --define "_buildrootdir ${rpm_dir}/BUILDROOT" \
             --define "_rpmdir ${rpm_dir}/RPMS" \
             --define "_srcrpmdir ${rpm_dir}/SRPMS" \
             --define "_specdir ${rpm_dir}/SPECS" \
             --define "_sourcedir ${rpm_dir}/SOURCES" \
             -bb "$spec_file"
    
    # Copy RPM to packages directory
    find "${rpm_dir}/RPMS" -name "*.rpm" -exec cp {} "${PACKAGES_DIR}/" \;
    
    print_success "RPM package created in ${PACKAGES_DIR}/"
}

# Build Arch Linux package
build_arch() {
    print_status "Building Arch Linux package..."
    
    local arch_dir="${PACKAGES_DIR}/arch"
    local pkgbuild_file="${arch_dir}/PKGBUILD"
    
    mkdir -p "$arch_dir"
    
    # Determine Arch Linux architecture
    local arch_pkg_arch
    case "$ARCH" in
        amd64) arch_pkg_arch="x86_64" ;;
        arm64) arch_pkg_arch="aarch64" ;;
        armhf) arch_pkg_arch="armv7h" ;;
        *) arch_pkg_arch="x86_64" ;;
    esac
    
    # Calculate SHA256 checksum if binary exists
    local sha256sum_value="SKIP"
    if [ -f "$BINARY_PATH" ]; then
        sha256sum_value=$(sha256sum "$BINARY_PATH" | cut -d' ' -f1)
        print_status "Calculated SHA256: $sha256sum_value"
    fi
    
    # Create PKGBUILD
    cat > "$pkgbuild_file" << EOF
# Maintainer: ${MAINTAINER}
pkgname=${APP_NAME}
pkgver=${VERSION}
pkgrel=1
pkgdesc="${APP_DESCRIPTION}"
arch=('${arch_pkg_arch}')
url="${APP_URL}"
license=('${LICENSE}')
depends=()
makedepends=()
source=("generator::\${url}/releases/download/\${pkgver}/generator-linux-${ARCH}.tar.gz")
sha256sums=('${sha256sum_value}')

package() {
    install -Dm755 "\${srcdir}/generator-linux-${ARCH}/generator" "\${pkgdir}/usr/local/bin/generator"
    install -Dm644 "\${srcdir}/generator-linux-${ARCH}/README.md" "\${pkgdir}/usr/share/doc/\${pkgname}/README.md"
    
    # Install man page
    install -Dm644 /dev/stdin "\${pkgdir}/usr/share/man/man1/generator.1.gz" << 'MANPAGE'
$(create_man_page_content | gzip -c | base64)
MANPAGE
    
    # Decode and decompress the man page
    base64 -d "\${pkgdir}/usr/share/man/man1/generator.1.gz" | gunzip > "\${pkgdir}/usr/share/man/man1/generator.1.tmp"
    gzip -c "\${pkgdir}/usr/share/man/man1/generator.1.tmp" > "\${pkgdir}/usr/share/man/man1/generator.1.gz"
    rm "\${pkgdir}/usr/share/man/man1/generator.1.tmp"
}
EOF

    # Create .SRCINFO
    cd "$arch_dir"
    if command -v makepkg >/dev/null 2>&1; then
        makepkg --printsrcinfo > .SRCINFO 2>/dev/null || true
        print_success "Arch Linux PKGBUILD created in ${arch_dir}/"
        print_status "To build: cd ${arch_dir} && makepkg -si"
    else
        print_warning "makepkg not available, created PKGBUILD only"
    fi
    cd - >/dev/null
}

# Create man page
create_man_page() {
    local man_file=$1
    local current_year=$(date +%Y)
    local man_date=$(date +"%B %Y")
    
    cat > "$man_file" << EOF
.TH GENERATOR 1 "$man_date" "Generator ${VERSION}" "User Commands"
.SH NAME
generator \- Open Source Project Generator
.SH SYNOPSIS
.B generator
[\fIOPTION\fR]... [\fICOMMAND\fR]
.SH DESCRIPTION
The Open Source Project Generator creates production-ready project structures with modern best practices, including frontend applications, backend services, mobile applications, infrastructure code, and comprehensive documentation.
.SH COMMANDS
.TP
.B generate
Generate a new project interactively
.TP
.B validate \fIPATH\fR
Validate a generated project structure
.TP
.B version
Show version information and latest package versions
.SH OPTIONS
.TP
.B \-h, \-\-help
Show help message
.TP
.B \-v, \-\-version
Show version information
.TP
.B \-\-config \fIFILE\fR
Use custom configuration file
.TP
.B \-\-output \fIDIR\fR
Set output directory for generated project
.TP
.B \-\-dry\-run
Preview generation without creating files
.SH EXAMPLES
.TP
Generate a new project interactively:
.B generator generate
.TP
Validate a project:
.B generator validate ./my-project
.TP
Show version and package information:
.B generator version --packages
.SH FILES
.TP
.I ~/.config/generator/config.yaml
User configuration file
.TP
.I ~/.cache/generator/
Cache directory for package versions
.SH AUTHOR
Open Source Project Generator Team
.SH REPORTING BUGS
Report bugs at: ${APP_URL}/issues
.SH COPYRIGHT
This is free software; see the source for copying conditions.
EOF
}

# Create man page content for PKGBUILD
create_man_page_content() {
    local current_year=$(date +%Y)
    local man_date=$(date +"%B %Y")
    
    cat << EOF
.TH GENERATOR 1 "$man_date" "Generator ${VERSION}" "User Commands"
.SH NAME
generator \- Open Source Project Generator
.SH SYNOPSIS
.B generator
[\fIOPTION\fR]... [\fICOMMAND\fR]
.SH DESCRIPTION
The Open Source Project Generator creates production-ready project structures with modern best practices, including frontend applications, backend services, mobile applications, infrastructure code, and comprehensive documentation.
.SH COMMANDS
.TP
.B generate
Generate a new project interactively
.TP
.B validate \fIPATH\fR
Validate a generated project structure
.TP
.B version
Show version information and latest package versions
.SH OPTIONS
.TP
.B \-h, \-\-help
Show help message
.TP
.B \-v, \-\-version
Show version information
.TP
.B \-\-config \fIFILE\fR
Use custom configuration file
.TP
.B \-\-output \fIDIR\fR
Set output directory for generated project
.TP
.B \-\-dry\-run
Preview generation without creating files
.SH EXAMPLES
.TP
Generate a new project interactively:
.B generator generate
.TP
Validate a project:
.B generator validate ./my-project
.TP
Show version and package information:
.B generator version --packages
.SH FILES
.TP
.I ~/.config/generator/config.yaml
User configuration file
.TP
.I ~/.cache/generator/
Cache directory for package versions
.SH AUTHOR
Open Source Project Generator Team
.SH REPORTING BUGS
Report bugs at: ${APP_URL}/issues
.SH COPYRIGHT
This is free software; see the source for copying conditions.
EOF
}

# Show usage
show_usage() {
    cat << EOF
Package Build Script

Usage: $0 [OPTIONS] PACKAGE_TYPE

Description:
  Creates distribution packages (DEB, RPM, Arch) for the generator binary.

Package Types:
  deb     Build Debian/Ubuntu package (.deb)
  rpm     Build Red Hat/CentOS package (.rpm)
  arch    Build Arch Linux package (PKGBUILD)
  all     Build all package types

Options:
  --clean         Remove old packages before building
  -h, --help      Show this help message

Environment Variables:
  VERSION         Set the version number (default: 1.0.0)
                  Must follow semantic versioning (e.g., 1.0.0 or v1.0.0)

Examples:
  $0 deb                    # Build DEB package
  VERSION=1.2.0 $0 rpm      # Build RPM with specific version
  $0 --clean all            # Clean and build all packages

Requirements:
  - dpkg-deb (for DEB packages)
  - rpmbuild (for RPM packages)
  - makepkg (for Arch packages)

Output:
  Packages: packages/
EOF
}

# Build packages
build_packages() {
    local package_type=$1
    
    print_status "Building packages for ${APP_NAME} ${VERSION}"
    
    prepare_package_dir
    
    case $package_type in
        "deb")
            check_dependencies "deb"
            build_deb
            ;;
        "rpm")
            check_dependencies "rpm"
            build_rpm
            ;;
        "arch")
            # For Arch, always create PKGBUILD even if makepkg isn't available
            build_arch
            ;;
        "all")
            local failed_builds=()
            
            # Try to build all packages, continue on failure
            if check_dependencies "deb" 2>/dev/null; then
                if ! build_deb; then
                    failed_builds+=("deb")
                fi
            else
                print_warning "Skipping DEB build - dependencies not met"
            fi
            
            if check_dependencies "rpm" 2>/dev/null; then
                if ! build_rpm; then
                    failed_builds+=("rpm")
                fi
            else
                print_warning "Skipping RPM build - dependencies not met"
            fi
            
            # Always create Arch PKGBUILD (doesn't require makepkg to create the file)
            if ! build_arch; then
                failed_builds+=("arch")
            fi
            
            if [ ${#failed_builds[@]} -gt 0 ]; then
                print_warning "Some builds failed: ${failed_builds[*]}"
            fi
            ;;
    esac
    
    print_success "Package build completed!"
    print_status "Packages created in: ${PACKAGES_DIR}/"
    ls -lh "$PACKAGES_DIR/" 2>/dev/null || true
}

# Handle command line arguments
case "${1:-}" in
    --clean)
        shift
        if [ -z "$1" ]; then
            print_error "Package type required after --clean"
            show_usage
            exit 1
        fi
        print_status "Cleaning old packages..."
        rm -rf "$PACKAGES_DIR"
        build_packages "$1"
        ;;
    deb|rpm|arch|all)
        build_packages "$1"
        ;;
    "help"|"-h"|"--help")
        show_usage
        exit 0
        ;;
    *)
        print_error "Unknown option or missing package type: ${1:-}"
        show_usage
        exit 1
        ;;
esac