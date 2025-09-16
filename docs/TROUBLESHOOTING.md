# Troubleshooting Guide

This guide helps you resolve common issues with the Open Source Template Generator installation and usage.

## Installation Issues

### Binary Not Found After Installation

**Problem**: `generator: command not found` after installation.

**Solutions**:

1. **Check PATH**:

   ```bash
   echo $PATH
   # Ensure /usr/local/bin or your installation directory is included
   ```

2. **Add to PATH manually**:

   ```bash
   # For bash
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   
   # For zsh
   echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Use full path**:

   ```bash
   /usr/local/bin/generator --version
   ```

4. **Check installation location**:

   ```bash
   which generator
   ls -la /usr/local/bin/generator
   ```

### Permission Denied Errors

**Problem**: `Permission denied` when running the generator.

**Solutions**:

1. **Make binary executable**:

   ```bash
   chmod +x /usr/local/bin/generator
   ```

2. **Check file permissions**:

   ```bash
   ls -la /usr/local/bin/generator
   # Should show: -rwxr-xr-x
   ```

3. **Reinstall with correct permissions**:

   ```bash
   sudo cp generator /usr/local/bin/
   sudo chmod +x /usr/local/bin/generator
   ```

### Package Installation Failures

**Problem**: Package installation fails with dependency errors.

**Solutions**:

1. **Update package lists**:

   ```bash
   # Debian/Ubuntu
   sudo apt update
   
   # Red Hat/CentOS
   sudo yum update
   ```

2. **Fix broken dependencies**:

   ```bash
   # Debian/Ubuntu
   sudo apt --fix-broken install
   
   # Red Hat/CentOS
   sudo yum check
   sudo yum update
   ```

3. **Install missing dependencies manually**:

   ```bash
   # Check what's missing
   dpkg -I generator_1.0.0_amd64.deb
   
   # Install dependencies
   sudo apt install <missing-packages>
   ```

### Download Failures

**Problem**: Cannot download releases or installation script fails.

**Solutions**:

1. **Check internet connection**:

   ```bash
   ping github.com
   ```

2. **Use alternative download method**:

   ```bash
   # If curl fails, try wget
   wget https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-linux-amd64.tar.gz
   ```

3. **Manual download**:
   - Visit the [releases page](https://github.com/cuesoftinc/open-source-project-generator/releases)
   - Download manually through browser
   - Extract and install manually

4. **Corporate firewall/proxy**:

   ```bash
   # Set proxy for curl
   export https_proxy=http://proxy.company.com:8080
   
   # Set proxy for wget
   export https_proxy=http://proxy.company.com:8080
   ```

## Runtime Issues

### Template Generation Failures

**Problem**: Generator fails to create project templates.

**Solutions**:

1. **Check output directory permissions**:

   ```bash
   ls -la /path/to/output/directory
   mkdir -p /path/to/output/directory
   ```

2. **Use different output directory**:

   ```bash
   generator generate --output ~/my-projects/new-project
   ```

3. **Run with verbose logging**:

   ```bash
   generator generate --log-level debug
   ```

4. **Check disk space**:

   ```bash
   df -h
   ```

### Network/API Errors

**Problem**: Cannot fetch latest package versions or templates.

**Solutions**:

1. **Check internet connectivity**:

   ```bash
   curl -I https://registry.npmjs.org/
   curl -I https://api.github.com/
   ```

2. **Use cached versions**:

   ```bash
   generator generate --offline
   ```

3. **Configure proxy settings**:

   ```bash
   export HTTP_PROXY=http://proxy.company.com:8080
   export HTTPS_PROXY=http://proxy.company.com:8080
   generator generate
   ```

4. **Use custom configuration**:

   ```yaml
   # ~/.config/generator/config.yaml
   network:
     timeout: 30s
     retries: 3
     proxy: "http://proxy.company.com:8080"
   ```

### Configuration Issues

**Problem**: Generator not using expected configuration.

**Solutions**:

1. **Check configuration file location**:

   ```bash
   generator config show
   generator config path
   ```

2. **Validate configuration syntax**:

   ```bash
   generator config validate
   ```

3. **Use explicit configuration**:

   ```bash
   generator generate --config /path/to/config.yaml
   ```

4. **Reset to defaults**:

   ```bash
   generator config reset
   ```

### Memory/Performance Issues

**Problem**: Generator runs slowly or uses too much memory.

**Solutions**:

1. **Check system resources**:

   ```bash
   free -h
   top
   ```

2. **Reduce concurrent operations**:

   ```yaml
   # config.yaml
   performance:
     max_concurrent_downloads: 2
     template_cache_size: 100MB
   ```

3. **Clear cache**:

   ```bash
   generator cache clear
   ```

4. **Use minimal templates**:

   ```bash
   generator generate --minimal
   ```

## Platform-Specific Issues

### macOS Issues

**Problem**: "generator cannot be opened because the developer cannot be verified"

**Solutions**:

1. **Allow in Security & Privacy**:
   - Go to System Preferences → Security & Privacy
   - Click "Allow Anyway" next to the generator message

2. **Remove quarantine attribute**:

   ```bash
   xattr -d com.apple.quarantine /usr/local/bin/generator
   ```

3. **Build from source**:

   ```bash
   git clone https://github.com/cuesoftinc/open-source-project-generator.git
   cd open-source-project-generator
   make build
   ```

**Problem**: Homebrew installation fails

**Solutions**:

1. **Update Homebrew**:

   ```bash
   brew update
   brew doctor
   ```

2. **Install manually**:

   ```bash
   curl -L https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-darwin-amd64.tar.gz | tar -xz
   mv generator-darwin-amd64/generator /usr/local/bin/
   ```

### Windows Issues

**Problem**: Windows Defender blocks the executable

**Solutions**:

1. **Add exception in Windows Defender**:
   - Open Windows Security
   - Go to Virus & threat protection
   - Add exclusion for the generator executable

2. **Download from trusted source**:
   - Only download from official GitHub releases
   - Verify checksums

**Problem**: PowerShell execution policy prevents running scripts

**Solutions**:

1. **Change execution policy**:

   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

2. **Bypass for single script**:

   ```powershell
   PowerShell -ExecutionPolicy Bypass -File install.ps1
   ```

**Problem**: PATH not updated after installation

**Solutions**:

1. **Restart terminal/PowerShell**

2. **Update PATH manually**:

   ```powershell
   $env:PATH += ";C:\Program Files\generator"
   ```

3. **Use system environment variables**:
   - Open System Properties → Advanced → Environment Variables
   - Edit PATH and add generator directory

### Linux Issues

**Problem**: `GLIBC` version errors

**Solutions**:

1. **Check GLIBC version**:

   ```bash
   ldd --version
   ```

2. **Use static binary** (if available):

   ```bash
   # Download static build
   wget https://github.com/cuesoftinc/open-source-project-generator/releases/latest/download/generator-linux-amd64-static.tar.gz
   ```

3. **Build from source**:

   ```bash
   git clone https://github.com/cuesoftinc/open-source-project-generator.git
   cd open-source-project-generator
   CGO_ENABLED=0 go build -o generator ./cmd/generator
   ```

**Problem**: Missing dependencies for templates

**Solutions**:

1. **Install Node.js** (for frontend templates):

   ```bash
   # Ubuntu/Debian
   curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
   sudo apt-get install -y nodejs
   
   # CentOS/RHEL
   curl -fsSL https://rpm.nodesource.com/setup_20.x | sudo bash -
   sudo yum install -y nodejs
   ```

2. **Install Docker** (for containerized templates):

   ```bash
   # Ubuntu/Debian
   sudo apt-get install docker.io
   
   # CentOS/RHEL
   sudo yum install docker
   ```

3. **Install Go** (for backend templates):

   ```bash
   # Download and install Go
   wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   ```

## Docker Issues

**Problem**: Docker container fails to run

**Solutions**:

1. **Check Docker installation**:

   ```bash
   docker --version
   docker run hello-world
   ```

2. **Use correct volume mounts**:

   ```bash
   docker run -v $(pwd):/workspace ghcr.io/cuesoftinc/open-source-project-generator:latest
   ```

3. **Check permissions**:

   ```bash
   # Ensure current user can access Docker
   sudo usermod -aG docker $USER
   # Logout and login again
   ```

## Template-Specific Issues

### Frontend Template Issues

**Problem**: Node.js version conflicts in generated projects

**Solutions**:

1. **Check Node.js version**:

   ```bash
   node --version
   # Should be 20.x or later
   ```

2. **Update Node.js**:

   ```bash
   # Using nvm
   nvm install 20
   nvm use 20
   
   # Or download from nodejs.org
   ```

3. **Use .nvmrc in generated projects**:

   ```bash
   cd generated-project/App/main
   nvm use
   npm install
   ```

**Problem**: Next.js build failures

**Solutions**:

1. **Clear Next.js cache**:

   ```bash
   rm -rf .next
   npm run build
   ```

2. **Update dependencies**:

   ```bash
   npm update
   npm audit fix
   ```

3. **Check TypeScript configuration**:

   ```bash
   npx tsc --noEmit
   ```

### Backend Template Issues

**Problem**: Go module resolution errors

**Solutions**:

1. **Initialize Go module**:

   ```bash
   cd generated-project/CommonServer
   go mod init your-module-name
   go mod tidy
   ```

2. **Update Go dependencies**:

   ```bash
   go get -u ./...
   go mod tidy
   ```

3. **Check Go version compatibility**:

   ```bash
   go version
   # Ensure compatibility with go.mod requirements
   ```

**Problem**: Database connection issues

**Solutions**:

1. **Check environment variables**:

   ```bash
   # Ensure database configuration is set
   cat .env
   ```

2. **Start database services**:

   ```bash
   docker-compose up -d postgres redis
   ```

3. **Run database migrations**:

   ```bash
   make migrate-up
   ```

### Mobile Template Issues

**Problem**: Android build failures

**Solutions**:

1. **Check Android SDK**:

   ```bash
   # Ensure Android SDK is installed and configured
   echo $ANDROID_HOME
   ```

2. **Update Gradle wrapper**:

   ```bash
   cd Mobile/android
   ./gradlew wrapper --gradle-version=8.5
   ```

3. **Clean and rebuild**:

   ```bash
   ./gradlew clean
   ./gradlew build
   ```

**Problem**: iOS build failures

**Solutions**:

1. **Update Xcode**:
   - Ensure Xcode is updated to latest version
   - Install command line tools: `xcode-select --install`

2. **Install CocoaPods dependencies**:

   ```bash
   cd Mobile/ios
   pod install --repo-update
   ```

3. **Clean build folder**:

   ```bash
   # In Xcode: Product → Clean Build Folder
   # Or via command line:
   xcodebuild clean -workspace YourApp.xcworkspace -scheme YourApp
   ```

### Infrastructure Template Issues

**Problem**: Docker build failures

**Solutions**:

1. **Check Docker version**:

   ```bash
   docker --version
   # Should be 24.x or later
   ```

2. **Clear Docker cache**:

   ```bash
   docker system prune -a
   ```

3. **Build with no cache**:

   ```bash
   docker build --no-cache -t your-app .
   ```

**Problem**: Kubernetes deployment issues

**Solutions**:

1. **Check cluster connectivity**:

   ```bash
   kubectl cluster-info
   kubectl get nodes
   ```

2. **Validate manifests**:

   ```bash
   kubectl apply --dry-run=client -f Deploy/k8s/
   ```

3. **Check resource quotas**:

   ```bash
   kubectl describe quota
   kubectl top nodes
   ```

## Build Issues

**Problem**: Build from source fails

**Solutions**:

1. **Check Go version**:

   ```bash
   go version
   # Should be 1.23 or later
   ```

2. **Update dependencies**:

   ```bash
   go mod download
   go mod tidy
   ```

3. **Clear module cache**:

   ```bash
   go clean -modcache
   go mod download
   ```

4. **Build with verbose output**:

   ```bash
   go build -v -o generator ./cmd/generator
   ```

## Getting Help

### Diagnostic Information

When reporting issues, please include:

```bash
# System information
uname -a
go version 2>/dev/null || echo "Go not installed"
docker --version 2>/dev/null || echo "Docker not installed"

# Generator information
generator --version
generator config show

# Environment
echo "PATH: $PATH"
echo "GOPATH: $GOPATH"
echo "HOME: $HOME"
```

### Log Files

Enable debug logging for detailed troubleshooting:

```bash
# Set log level
export GENERATOR_LOG_LEVEL=debug

# Run with verbose output
generator generate --log-level debug --verbose

# Check log files
ls -la ~/.cache/generator/logs/
```

### Common Log Locations

- Linux: `~/.cache/generator/logs/`
- macOS: `~/Library/Caches/generator/logs/`
- Windows: `%LOCALAPPDATA%\generator\logs\`

### Support Channels

1. **GitHub Issues**: [Report bugs and feature requests](https://github.com/cuesoftinc/open-source-project-generator/issues)
2. **Discussions**: [Community support](https://github.com/cuesoftinc/open-source-project-generator/discussions)
3. **Documentation**: [Wiki and guides](https://github.com/cuesoftinc/open-source-project-generator/wiki)
4. **Email**: [Direct support](mailto:support@generator.dev)

### Before Reporting Issues

1. **Search existing issues**: Check if the problem has been reported
2. **Try latest version**: Update to the latest release
3. **Minimal reproduction**: Provide steps to reproduce the issue
4. **Include logs**: Attach relevant log files and error messages
5. **System details**: Include OS, architecture, and version information

## FAQ

### Q: Can I use the generator offline?

A: Yes, use the `--offline` flag to use cached templates and versions.

### Q: How do I update templates?

A: Templates are updated with each release. Update the generator to get the latest templates.

### Q: Can I create custom templates?

A: Yes, see the [Template Development Guide](TEMPLATE_DEVELOPMENT.md) for details.

### Q: Is the generator safe to use in corporate environments?

A: Yes, the generator doesn't send data externally except for fetching package versions. Use `--offline` mode for complete isolation.

### Q: How do I contribute to the project?

A: See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.
