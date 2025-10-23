# Scripts Directory

This directory contains utility scripts for the Open Source Project Generator.

## Quick Reference

| Script | Purpose | Usage |
|--------|---------|-------|
| `build.sh` | Cross-platform binary compilation | `./scripts/build.sh [platform]` |
| `build-packages.sh` | Build distribution packages | `./scripts/build-packages.sh [deb\|rpm\|arch\|all]` |
| `check-latest-versions.sh` | Check for outdated dependencies | `./scripts/check-latest-versions.sh [--quiet\|--json]` |
| `ci-test.sh` | Run comprehensive test suite | `./scripts/ci-test.sh` |
| `get-version.sh` | Get version information | `./scripts/get-version.sh [version\|commit\|all]` |
| `install.sh` | Install the generator binary | `./scripts/install.sh [--version VERSION]` |
| `run_performance_benchmarks.sh` | Run performance benchmarks | `./scripts/run_performance_benchmarks.sh` |
| `validate-setup.sh` | Validate project setup | `./scripts/validate-setup.sh` |

---

## Detailed Documentation

### build.sh

Cross-platform binary compilation script that builds the generator for multiple operating systems and architectures.

**Usage:**

```bash
# Build for all platforms
./scripts/build.sh

# Build for specific platform
./scripts/build.sh linux-amd64

# Build with specific version
VERSION=1.2.0 ./scripts/build.sh

# Clean before building
./scripts/build.sh clean
```

**Supported Platforms:**

- `linux-amd64`, `linux-arm64`
- `darwin-amd64`, `darwin-arm64` (macOS)
- `windows-amd64`, `windows-arm64`

**Environment Variables:**

- `VERSION` - Set version number (default: from git tags or "dev")
- `GIT_COMMIT` - Git commit hash (auto-detected)
- `BUILD_TIME` - Build timestamp (auto-generated)

**Output:** `dist/` directory with compiled binaries

**Features:**

- ✅ Cross-compilation for 6 platforms
- ✅ Version injection via ldflags
- ✅ Parallel builds for speed
- ✅ Build verification
- ✅ Colored output for status

---

### build-packages.sh

Build distribution packages for various Linux distributions and package managers.

**Usage:**

```bash
# Build DEB package (Debian/Ubuntu)
./scripts/build-packages.sh deb

# Build RPM package (Red Hat/CentOS/Fedora)
./scripts/build-packages.sh rpm

# Build Arch Linux package
./scripts/build-packages.sh arch

# Build all package types
./scripts/build-packages.sh all

# Clean old packages before building
./scripts/build-packages.sh --clean all

# Build with specific version
VERSION=1.2.0 ./scripts/build-packages.sh deb
```

**Package Types:**

- `deb` - Debian/Ubuntu (.deb)
- `rpm` - Red Hat/CentOS/Fedora (.rpm)
- `arch` - Arch Linux (PKGBUILD)

**Requirements:**

- `dpkg-deb` (for DEB packages)
- `rpmbuild` (for RPM packages)
- `makepkg` (for Arch packages)

**Environment Variables:**

- `VERSION` - Package version (default: 1.0.0)

**Output:** `packages/` directory with distribution packages

**Features:**

- ✅ Multi-format package generation
- ✅ Automatic dependency detection
- ✅ Graceful fallback if tools missing
- ✅ Package verification
- ✅ Checksums generation

---

### check-latest-versions.sh

Check for the latest versions of all dependencies used in the project. Automatically checks 24 packages across frontend, backend, mobile, and infrastructure.

**Usage:**

```bash
# Standard output with colors
./scripts/check-latest-versions.sh

# Only show outdated packages
./scripts/check-latest-versions.sh --quiet

# JSON output for automation
./scripts/check-latest-versions.sh --json

# Show help
./scripts/check-latest-versions.sh --help
```

**What It Checks (24 packages):**

- **Frontend (5):** Next.js, React, React-DOM, TypeScript, create-next-app
- **Backend (4):** Gin, Gin CORS, Echo, Fiber
- **Android (7):** SDK API Level, Kotlin, Gradle, AndroidX libraries
- **iOS (3):** Swift, Xcode, iOS Deployment Target
- **Infrastructure (3):** Terraform, Kubernetes, Go toolchain
- **Docker (3):** Alpine, Golang, Ubuntu base images

**Features:**

- ✅ 95% automation rate (24/25 packages)
- ✅ Color-coded output (green = up-to-date, red = outdated)
- ✅ JSON output for CI/CD integration
- ✅ Quiet mode to only show outdated packages
- ✅ Smart version comparison logic
- ✅ Filters out beta/alpha versions
- ✅ GitHub API integration with token support

**Requirements:**

- `npm` (for Node.js packages)
- `go` (for Go modules)
- `curl` (for API checks)
- `jq` (optional, for JSON processing and summary)

**Environment Variables:**

- `GITHUB_TOKEN` - Optional GitHub token for higher rate limits (5000/hour vs 60/hour)

**Example Output:**

```bash
=== Checking Latest Versions ===

### Node.js & Frontend ###
Checking create-next-app... ✓ Latest: 16.0.0 (Current: 16.0.0)
Checking next... ✓ Latest: 16.0.0 (Current: 16.0.0)
Checking react... ✓ Latest: 19.2.0 (Current: 19.2.0)

### Go Backend Frameworks ###
Checking github.com/gin-gonic/gin... ✓ Latest: v1.11.0 (Current: v1.11.0)

### Android Build Tools ###
Checking Gradle... ✓ Latest: 9.1.0 (Current: 9.1.0)
Checking Kotlin... ✓ Latest: v2.2.21 (Current: 2.2.21)
Checking Android SDK... ✓ Latest API: 36 (Current: 36)

### iOS ###
Checking Swift... ✓ Latest: swift-6.2-RELEASE (Current: swift-6.2-RELEASE)
Checking Xcode... ✓ Latest: 26.0.1 (Current: 26.0.1)
Checking iOS Deployment Target... ✓ Latest: 26.0 (Current: 26.0)

=== Summary ===
✓ 24 package(s) up-to-date
```

**CI/CD Integration:**

```yaml
# GitHub Actions example
- name: Check for outdated dependencies
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    ./scripts/check-latest-versions.sh --json > versions.json
    outdated=$(jq '[.[] | select(.status == "outdated")] | length' versions.json)
    if [ "$outdated" -gt 0 ]; then
      echo "::warning::$outdated package(s) are outdated"
      jq '[.[] | select(.status == "outdated")]' versions.json
    fi
```

**Makefile Integration:**

```makefile
.PHONY: check-versions
check-versions:
 @./scripts/check-latest-versions.sh

.PHONY: check-versions-quiet
check-versions-quiet:
 @./scripts/check-latest-versions.sh --quiet
```

---

### ci-test.sh

Comprehensive test suite runner for CI/CD pipelines with coverage reporting, race detection, and parallel execution.

**Usage:**

```bash
# Run full test suite
./scripts/ci-test.sh

# Run with custom timeout
TEST_TIMEOUT=15m ./scripts/ci-test.sh

# Run with coverage threshold
COVERAGE_THRESHOLD=80 ./scripts/ci-test.sh

# Run with more parallel tests
PARALLEL_TESTS=8 ./scripts/ci-test.sh
```

**Environment Variables:**

- `TEST_TIMEOUT` - Test timeout (default: 10m)
- `COVERAGE_THRESHOLD` - Minimum coverage percentage (default: 0)
- `PARALLEL_TESTS` - Number of parallel test processes (default: 4)

**Output:** `test-reports/` directory with coverage and test results

**Features:**

- ✅ Parallel test execution
- ✅ Race condition detection
- ✅ Coverage reporting (HTML + text)
- ✅ Failed test tracking
- ✅ JUnit XML output for CI
- ✅ Colored output

---

### get-version.sh

Get version information from git tags and repository metadata.

**Usage:**

```bash
# Get version string
./scripts/get-version.sh
./scripts/get-version.sh version

# Get git commit hash
./scripts/get-version.sh commit

# Get build timestamp
./scripts/get-version.sh build-time

# Get git branch name
./scripts/get-version.sh branch

# Get all version information
./scripts/get-version.sh all

# Get Go build ldflags
./scripts/get-version.sh ldflags

# Validate version format
./scripts/get-version.sh validate
```

**Features:**

- ✅ Semantic versioning support
- ✅ Git tag detection
- ✅ Fallback to "dev" if not in git repo
- ✅ Build metadata generation
- ✅ ldflags generation for Go builds

---

### install.sh

Installation script that automatically detects the platform and installs the appropriate binary.

**Usage:**

```bash
# Install latest version
./scripts/install.sh

# Install specific version
./scripts/install.sh --version v1.2.0

# Force binary installation (skip package managers)
./scripts/install.sh --force-binary

# Install to user directory
./scripts/install.sh --user

# Show help
./scripts/install.sh --help
```

**Options:**

- `--version VERSION` - Install specific version
- `--force-binary` - Skip package managers, install binary directly
- `--user` - Install to user directory (~/.local/bin)

**Features:**

- ✅ Automatic platform detection
- ✅ Package manager integration (apt, yum, pacman)
- ✅ Binary fallback
- ✅ Retry logic for network failures
- ✅ Checksum verification
- ✅ PATH configuration

---

### run_performance_benchmarks.sh

Run comprehensive performance benchmarks and compare results against baseline.

**Usage:**

```bash
# Run all benchmarks
./scripts/run_performance_benchmarks.sh

# Run with custom settings
BENCHMARK_COUNT=10 ./scripts/run_performance_benchmarks.sh

# Run with custom timeout
BENCHMARK_TIMEOUT=45m ./scripts/run_performance_benchmarks.sh

# Set regression threshold
REGRESSION_THRESHOLD=15 ./scripts/run_performance_benchmarks.sh
```

**Environment Variables:**

- `BENCHMARK_TIMEOUT` - Benchmark timeout (default: 30m)
- `BENCHMARK_COUNT` - Number of benchmark runs (default: 5)
- `REGRESSION_THRESHOLD` - Percentage threshold for regression detection (default: 10)

**Output:** `benchmark_results/` directory with:

- Baseline results
- Current results
- Performance comparison
- Markdown report
- CPU/memory profiles

**Features:**

- ✅ Baseline comparison
- ✅ Regression detection
- ✅ CPU and memory profiling
- ✅ Statistical analysis
- ✅ Markdown report generation
- ✅ benchstat integration

---

### validate-setup.sh

Validate project setup, configuration, and development environment.

**Usage:**

```bash
# Run all validations
./scripts/validate-setup.sh

# Show help
./scripts/validate-setup.sh --help
```

**What It Validates:**

- ✅ Required tools (Go, npm, Docker)
- ✅ Go version compatibility
- ✅ Script executability
- ✅ Makefile syntax
- ✅ Dockerfile syntax
- ✅ CI/CD workflow files
- ✅ Configuration files
- ✅ Directory structure

**Features:**

- ✅ Comprehensive validation checks
- ✅ Error and warning tracking
- ✅ Colored output
- ✅ Summary report
- ✅ Exit code based on results

---

## Adding New Scripts

When adding new scripts:

1. **Make executable:** `chmod +x scripts/your-script.sh`
2. **Add shebang:** `#!/bin/bash`
3. **Add usage documentation** in comments at the top
4. **Add error handling:** `set -e`
5. **Use consistent output:**
   - Use colored output functions (print_status, print_success, print_error)
   - Provide `--help` option
   - Support both interactive and non-interactive modes
6. **Document in this README** following the format above
7. **Add to Makefile** for easy access (optional)
8. **Test on multiple platforms** (Linux, macOS)

## Best Practices

### Error Handling

```bash
set -e  # Exit on errors
trap cleanup EXIT  # Cleanup on exit
```

### Colored Output

```bash
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}
```

### Help Option

```bash
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_usage
    exit 0
fi
```

### Environment Variables

```bash
VERSION=${VERSION:-"1.0.0"}  # Default value
```

### JSON Output

```bash
if [ "$JSON_OUTPUT" = true ]; then
    echo '{"status": "success"}'
fi
```

## Common Issues

### Script Not Executable

```bash
chmod +x scripts/script-name.sh
```

### Wrong Line Endings (Windows)

```bash
dos2unix scripts/script-name.sh
```

### Missing Dependencies

Run `./scripts/validate-setup.sh` to check for missing tools.

## Integration with Makefile

All scripts can be accessed via Makefile targets:

```bash
make build              # Run build.sh
make test               # Run ci-test.sh
make check-versions     # Run check-latest-versions.sh
make validate           # Run validate-setup.sh
```

See the main `Makefile` for all available targets.
