# Installation Guide

This guide provides comprehensive installation instructions for the Open Source Template Generator across different platforms and package managers.

## Quick Install

### Linux and macOS

```bash
curl -sSL https://raw.githubusercontent.com/open-source-template-generator/generator/main/scripts/install.sh | bash
```

### Windows

Download the Windows binary from the [releases page](https://github.com/open-source-template-generator/generator/releases) and follow the installation instructions included in the archive.

## Package Manager Installation

### Debian/Ubuntu (APT)

```bash
# Download and install the DEB package
wget https://github.com/open-source-template-generator/generator/releases/latest/download/generator_1.0.0_amd64.deb
sudo dpkg -i generator_1.0.0_amd64.deb

# Fix dependencies if needed
sudo apt-get install -f
```

### Red Hat/CentOS/Fedora (YUM/DNF)

```bash
# Using YUM
sudo yum install https://github.com/open-source-template-generator/generator/releases/latest/download/generator-1.0.0-1.x86_64.rpm

# Using DNF
sudo dnf install https://github.com/open-source-template-generator/generator/releases/latest/download/generator-1.0.0-1.x86_64.rpm
```

### Arch Linux (AUR)

```bash
# Using yay
yay -S generator

# Using paru
paru -S generator

# Manual installation
git clone https://aur.archlinux.org/generator.git
cd generator
makepkg -si
```

### macOS (Homebrew)

```bash
# Add our tap (coming soon)
brew tap open-source-template-generator/tap
brew install generator
```

### Windows (Chocolatey)

```powershell
# Coming soon
choco install generator
```

### Windows (Scoop)

```powershell
# Coming soon
scoop bucket add generator https://github.com/open-source-template-generator/scoop-bucket
scoop install generator
```

## Manual Installation

### Download Pre-built Binaries

1. Visit the [releases page](https://github.com/open-source-template-generator/generator/releases)
2. Download the appropriate archive for your platform:
   - Linux: `generator-linux-amd64.tar.gz`
   - macOS (Intel): `generator-darwin-amd64.tar.gz`
   - macOS (Apple Silicon): `generator-darwin-arm64.tar.gz`
   - Windows: `generator-windows-amd64.zip`
   - FreeBSD: `generator-freebsd-amd64.tar.gz`

### Extract and Install

#### Linux/macOS/FreeBSD

```bash
# Extract the archive
tar -xzf generator-linux-amd64.tar.gz

# Move to installation directory
sudo mv generator-linux-amd64/generator /usr/local/bin/

# Make executable (if needed)
sudo chmod +x /usr/local/bin/generator

# Verify installation
generator --version
```

#### Windows

1. Extract the ZIP file to a directory (e.g., `C:\Program Files\generator`)
2. Add the directory to your PATH environment variable:
   - Open System Properties → Advanced → Environment Variables
   - Edit the PATH variable and add the generator directory
   - Or run: `setx PATH "%PATH%;C:\Program Files\generator"`
3. Open a new command prompt and run: `generator --help`

### User Installation (No Admin Rights)

If you don't have administrator privileges, you can install to your user directory:

#### Linux/macOS

```bash
# Create user bin directory
mkdir -p ~/.local/bin

# Extract and move binary
tar -xzf generator-linux-amd64.tar.gz
mv generator-linux-amd64/generator ~/.local/bin/

# Add to PATH (add to ~/.bashrc, ~/.zshrc, or ~/.profile)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc

# Reload shell configuration
source ~/.bashrc

# Verify installation
generator --version
```

#### Windows

```powershell
# Create user bin directory
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin"

# Extract and move binary (assuming you've extracted to Downloads)
Move-Item "$env:USERPROFILE\Downloads\generator-windows-amd64\generator.exe" "$env:USERPROFILE\bin\"

# Add to user PATH
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$userPath;$env:USERPROFILE\bin", "User")

# Restart PowerShell and verify
generator --version
```

## Docker Installation

### Configuration

The Docker images are published to GitHub Container Registry with dynamic organization names.
The image name automatically adapts to the repository owner:

- For `cuesoftinc` organization: `ghcr.io/cuesoftinc/open-source-project-generator:latest`
- For your fork: `ghcr.io/yourusername/open-source-project-generator:latest`

### Pull from GitHub Container Registry

```bash
# For cuesoftinc organization
docker pull ghcr.io/cuesoftinc/open-source-project-generator:latest

# For your own fork (replace 'yourusername' with your GitHub username)
docker pull ghcr.io/yourusername/open-source-project-generator:latest
```

### Run in Container

```bash
# Interactive mode (using cuesoftinc organization)
docker run -it --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate

# One-time generation (using cuesoftinc organization)
docker run --rm -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest generate --config /workspace/config.yaml

# Using your own fork (replace 'yourusername')
docker run -it --rm -v $(pwd):/workspace ghcr.io/yourusername/open-source-project-generator:latest generate
```

### Environment Variables

For development and customization, you can configure the following environment variables:

```bash
# Set your GitHub organization/username
export GITHUB_REPOSITORY_OWNER=yourusername

# Set Docker registry (default: ghcr.io)  
export DOCKER_REGISTRY=ghcr.io

# For authentication when pushing images
export GITHUB_ACTOR=yourusername
export GITHUB_TOKEN=your_github_token
```

Copy `env.example` to `.env` and customize for your setup:

```bash
cp env.example .env
# Edit .env with your configuration
```

### Build from Source

```bash
git clone https://github.com/open-source-template-generator/generator.git
cd generator
docker build -t generator:local .
```

## Build from Source

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional, for using Makefile)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/open-source-template-generator/generator.git
cd generator

# Install dependencies
go mod download

# Build the binary
go build -o bin/generator ./cmd/generator

# Or use Make
make build

# Install to system (optional)
sudo cp bin/generator /usr/local/bin/
```

### Cross-compilation

```bash
# Build for all platforms
make build-all

# Or build specific platform
GOOS=linux GOARCH=amd64 go build -o generator-linux-amd64 ./cmd/generator
```

## Verification

After installation, verify that the generator is working correctly:

```bash
# Check version
generator --version

# Show help
generator --help

# Test generation (dry run)
generator generate --dry-run

# Check available templates
generator list-templates
```

## Configuration

### Global Configuration

The generator looks for configuration files in the following locations:

1. `~/.config/generator/config.yaml` (Linux/macOS)
2. `%APPDATA%\generator\config.yaml` (Windows)
3. `./generator.yaml` (current directory)

### Environment Variables

- `GENERATOR_CONFIG`: Path to custom configuration file
- `GENERATOR_CACHE_DIR`: Custom cache directory
- `GENERATOR_LOG_LEVEL`: Log level (debug, info, warn, error)
- `GENERATOR_NO_COLOR`: Disable colored output

### Example Configuration

```yaml
# ~/.config/generator/config.yaml
default_output_dir: "~/projects"
default_license: "MIT"
default_organization: "myorg"

# Template preferences
templates:
  frontend: "nextjs-app"
  backend: "go-gin"
  mobile: "android-kotlin"

# Version preferences
versions:
  node: "20"
  go: "1.24"
  
# Cache settings
cache:
  ttl: "24h"
  enabled: true
```

## Updating

### Package Manager Updates

If installed via package manager, use the same package manager to update:

```bash
# APT
sudo apt update && sudo apt upgrade generator

# YUM/DNF
sudo yum update generator
# or
sudo dnf update generator

# Homebrew
brew update && brew upgrade generator

# AUR
yay -Syu generator
```

### Manual Updates

1. Download the latest release
2. Replace the existing binary
3. Verify the new version

### Automatic Updates

The generator can check for updates automatically:

```bash
# Check for updates
generator update --check

# Update to latest version
generator update --install

# Enable automatic update checks
generator config set auto_update_check true
```

## Uninstallation

### Package Manager

```bash
# APT
sudo apt remove generator

# YUM/DNF
sudo yum remove generator
# or
sudo dnf remove generator

# Homebrew
brew uninstall generator

# AUR
yay -R generator
```

### Manual Uninstallation

```bash
# Remove binary
sudo rm /usr/local/bin/generator

# Remove configuration (optional)
rm -rf ~/.config/generator

# Remove cache (optional)
rm -rf ~/.cache/generator
```

## Platform-Specific Notes

### Linux

- The generator requires `glibc 2.17+` (most modern distributions)
- For older distributions, build from source
- Some templates may require additional tools (Docker, Node.js, etc.)

### macOS

- Supports macOS 10.15+ (Catalina and later)
- On first run, you may need to allow the binary in Security & Privacy settings
- For Apple Silicon Macs, use the `arm64` version for better performance

### Windows

- Requires Windows 10 or later
- Windows Defender may flag the binary initially (false positive)
- PowerShell 5.1+ recommended for best experience
- Git Bash or WSL recommended for shell scripts

### FreeBSD

- Supports FreeBSD 12.0+
- May require additional packages for some templates
- Build from source for optimal compatibility

## Troubleshooting

For installation issues, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

## Support

- 📖 [Documentation](https://github.com/open-source-template-generator/generator/wiki)
- 🐛 [Issue Tracker](https://github.com/open-source-template-generator/generator/issues)
- 💬 [Discussions](https://github.com/open-source-template-generator/generator/discussions)
- 📧 [Email Support](mailto:support@generator.dev)
