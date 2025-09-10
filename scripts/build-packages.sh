#!/bin/bash

# Package build script for creating distribution packages
# Supports DEB, RPM, and other package formats

set -e

# Configuration
APP_NAME="generator"
APP_DESCRIPTION="Open Source Template Generator - Create production-ready project structures"
APP_URL="https://github.com/open-source-template-generator/generator"
MAINTAINER="Open Source Template Generator Team <team@example.com>"
LICENSE="MIT"
VERSION=${VERSION:-"1.0.0"}
ARCH="amd64"

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DIST_DIR="${PROJECT_ROOT}/dist"
PACKAGES_DIR="${PROJECT_ROOT}/packages"
BINARY_PATH="${DIST_DIR}/generator-linux-amd64/generator"

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

# Check dependencies
check_dependencies() {
    local package_type=$1
    
    case $package_type in
        "deb")
            if ! command -v dpkg-deb >/dev/null 2>&1; then
                print_error "dpkg-deb is required for building DEB packages"
                print_status "Install with: sudo apt-get install dpkg-dev"
                exit 1
            fi
            ;;
        "rpm")
            if ! command -v rpmbuild >/dev/null 2>&1; then
                print_error "rpmbuild is required for building RPM packages"
                print_status "Install with: sudo apt-get install rpm (Ubuntu) or sudo yum install rpm-build (CentOS)"
                exit 1
            fi
            ;;
        "arch")
            if ! command -v makepkg >/dev/null 2>&1; then
                print_error "makepkg is required for building Arch packages"
                print_status "This should be available on Arch Linux systems"
                exit 1
            fi
            ;;
    esac
}

# Prepare package directory
prepare_package_dir() {
    print_status "Preparing package directory..."
    
    # Clean and create packages directory
    rm -rf "$PACKAGES_DIR"
    mkdir -p "$PACKAGES_DIR"
    
    # Check if binary exists
    if [ ! -f "$BINARY_PATH" ]; then
        print_error "Binary not found at $BINARY_PATH"
        print_status "Please run the build script first to create the binary"
        exit 1
    fi
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
Version: ${VERSION#v}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: ${APP_DESCRIPTION}
 The Open Source Template Generator creates production-ready project structures
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

# Add /usr/local/bin to PATH if not already present
if ! echo "$PATH" | grep -q "/usr/local/bin"; then
    echo "Adding /usr/local/bin to PATH..."
    echo 'export PATH="/usr/local/bin:$PATH"' >> /etc/environment
fi

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
    dpkg-deb --build "$deb_dir" "${PACKAGES_DIR}/${APP_NAME}_${VERSION#v}_${ARCH}.deb"
    
    print_success "DEB package created: ${PACKAGES_DIR}/${APP_NAME}_${VERSION#v}_${ARCH}.deb"
}

# Build RPM package
build_rpm() {
    print_status "Building RPM package..."
    
    local rpm_dir="${PACKAGES_DIR}/rpm"
    local spec_file="${rpm_dir}/${APP_NAME}.spec"
    
    # Create RPM build directories
    mkdir -p "${rpm_dir}"/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
    mkdir -p "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/local/bin"
    mkdir -p "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/doc/${APP_NAME}"
    mkdir -p "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/man/man1"
    
    # Copy binary and documentation
    cp "$BINARY_PATH" "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/local/bin/"
    cp "${PROJECT_ROOT}/README.md" "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/doc/${APP_NAME}/"
    [ -f "${PROJECT_ROOT}/LICENSE" ] && cp "${PROJECT_ROOT}/LICENSE" "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/doc/${APP_NAME}/"
    
    # Create man page
    create_man_page "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/man/man1/generator.1"
    gzip "${rpm_dir}/BUILDROOT/${APP_NAME}-${VERSION#v}-1.x86_64/usr/share/man/man1/generator.1"
    
    # Create spec file
    cat > "$spec_file" << EOF
Name:           ${APP_NAME}
Version:        ${VERSION#v}
Release:        1%{?dist}
Summary:        ${APP_DESCRIPTION}

License:        ${LICENSE}
URL:            ${APP_URL}
BuildArch:      x86_64

%description
The Open Source Template Generator creates production-ready project structures
with modern best practices, including frontend applications, backend services,
mobile applications, infrastructure code, and comprehensive documentation.

Features:
- Interactive project generation with component selection
- Support for Next.js, Go, Android, and iOS applications
- Infrastructure as code templates (Docker, Kubernetes, Terraform)
- Complete CI/CD workflows with GitHub Actions
- Cross-platform support and package management integration

%files
/usr/local/bin/generator
/usr/share/doc/${APP_NAME}/README.md
/usr/share/man/man1/generator.1.gz
%if 0%{?rhel} >= 7 || 0%{?fedora}
%license /usr/share/doc/${APP_NAME}/LICENSE
%endif

%post
# Add /usr/local/bin to PATH if not already present
if ! echo "\$PATH" | grep -q "/usr/local/bin"; then
    echo 'export PATH="/usr/local/bin:\$PATH"' >> /etc/environment
fi

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
* $(date +'%a %b %d %Y') ${MAINTAINER} - ${VERSION#v}-1
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
    
    # Create PKGBUILD
    cat > "$pkgbuild_file" << EOF
# Maintainer: ${MAINTAINER}
pkgname=${APP_NAME}
pkgver=${VERSION#v}
pkgrel=1
pkgdesc="${APP_DESCRIPTION}"
arch=('x86_64')
url="${APP_URL}"
license=('${LICENSE}')
depends=()
makedepends=()
source=("generator::\${url}/releases/download/v\${pkgver}/generator-linux-amd64.tar.gz")
sha256sums=('SKIP')

package() {
    install -Dm755 "\${srcdir}/generator-linux-amd64/generator" "\${pkgdir}/usr/local/bin/generator"
    install -Dm644 "\${srcdir}/generator-linux-amd64/README.md" "\${pkgdir}/usr/share/doc/\${pkgname}/README.md"
    
    # Create man page
    install -Dm644 /dev/stdin "\${pkgdir}/usr/share/man/man1/generator.1" << 'MANPAGE'
$(create_man_page_content)
MANPAGE
}
EOF

    # Create .SRCINFO
    cd "$arch_dir"
    if command -v makepkg >/dev/null 2>&1; then
        makepkg --printsrcinfo > .SRCINFO
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
    
    cat > "$man_file" << 'EOF'
.TH GENERATOR 1 "2024" "Generator 1.0" "User Commands"
.SH NAME
generator \- Open Source Template Generator
.SH SYNOPSIS
.B generator
[\fIOPTION\fR]... [\fICOMMAND\fR]
.SH DESCRIPTION
The Open Source Template Generator creates production-ready project structures with modern best practices, including frontend applications, backend services, mobile applications, infrastructure code, and comprehensive documentation.
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
Open Source Template Generator Team
.SH REPORTING BUGS
Report bugs at: https://github.com/open-source-template-generator/generator/issues
.SH COPYRIGHT
This is free software; see the source for copying conditions.
EOF
}

# Create man page content for PKGBUILD
create_man_page_content() {
    cat << 'EOF'
.TH GENERATOR 1 "2024" "Generator 1.0" "User Commands"
.SH NAME
generator \- Open Source Template Generator
.SH SYNOPSIS
.B generator
[\fIOPTION\fR]... [\fICOMMAND\fR]
.SH DESCRIPTION
The Open Source Template Generator creates production-ready project structures with modern best practices.
.SH COMMANDS
.TP
.B generate
Generate a new project interactively
.TP
.B validate \fIPATH\fR
Validate a generated project structure
.TP
.B version
Show version information
.SH OPTIONS
.TP
.B \-h, \-\-help
Show help message
.TP
.B \-v, \-\-version
Show version information
.SH AUTHOR
Open Source Template Generator Team
EOF
}

# Show usage
show_usage() {
    echo "Usage: $0 [PACKAGE_TYPE]"
    echo ""
    echo "Package types:"
    echo "  deb     Build Debian/Ubuntu package (.deb)"
    echo "  rpm     Build Red Hat/CentOS package (.rpm)"
    echo "  arch    Build Arch Linux package (PKGBUILD)"
    echo "  all     Build all package types"
    echo ""
    echo "Environment variables:"
    echo "  VERSION    Set the version number (default: 1.0.0)"
    echo ""
    echo "Examples:"
    echo "  $0 deb                    # Build DEB package"
    echo "  VERSION=1.2.0 $0 rpm     # Build RPM with specific version"
    echo "  $0 all                    # Build all packages"
}

# Main function
main() {
    local package_type=${1:-""}
    
    if [ -z "$package_type" ]; then
        show_usage
        exit 1
    fi
    
    print_status "Building packages for ${APP_NAME} v${VERSION}"
    
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
            check_dependencies "arch"
            build_arch
            ;;
        "all")
            check_dependencies "deb"
            check_dependencies "rpm"
            build_deb
            build_rpm
            build_arch
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_error "Unknown package type: $package_type"
            show_usage
            exit 1
            ;;
    esac
    
    print_success "Package build completed!"
    print_status "Packages created in: ${PACKAGES_DIR}/"
    ls -la "$PACKAGES_DIR/"
}

# Run main function
main "$@"