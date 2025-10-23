# Version Management

This document explains the centralized version management system.

## Overview

All dependency versions are now managed in a single source of truth: `configs/versions.yaml`. This eliminates the need to update versions in multiple places across the codebase.

## Structure

The `versions.yaml` file contains version information for:

- **Frontend**: Next.js, React, React-DOM, TypeScript
- **Backend**: Go, Gin, Echo, Fiber frameworks
- **Android**: SDK, Kotlin, Gradle, AndroidX libraries
- **iOS**: Swift, Xcode, deployment targets
- **Docker**: Base images (Alpine, Go, Ubuntu)
- **Infrastructure**: Terraform, Kubernetes

## Usage

### In Code

Import and use the versions package:

```go
import "github.com/cuesoftinc/open-source-project-generator/pkg/versions"

// Get the version configuration
versionConfig, err := versions.Get()
if err != nil {
    // Handle error or use fallback
}

// Access specific versions
nextjsVersion := versionConfig.Frontend.NextJS.Version
ginVersion := versionConfig.Backend.Frameworks.Gin.Version
kotlinVersion := versionConfig.Android.Kotlin.Version
```

### Checking for Updates

To check which versions are outdated:

```bash
# Using Makefile (recommended)
make check-versions

# Or directly
./scripts/check-latest-versions.sh

# Only show outdated packages
./scripts/check-latest-versions.sh --quiet

# JSON output for automation
./scripts/check-latest-versions.sh --json
```

The script reads current versions from `configs/versions.yaml` and compares them with the latest available versions from external registries.

### Updating Versions

To update the `versions.yaml` file with the latest versions:

```bash
# Automatic mode (updates without prompting)
make update-versions

# Dry run (shows what would be updated)
./scripts/update-versions.sh --dry-run
```

**Note**: The update script requires `yq` to be installed:

- macOS: `brew install yq`
- Linux: `snap install yq`
- Or download from: <https://github.com/mikefarah/yq/releases>

If `yq` is not available, you can manually update `configs/versions.yaml` based on the output from `make check-versions`.

## How It Works

1. **Centralized Config**: `configs/versions.yaml` stores all version numbers
2. **Go Package**: `pkg/versions` provides a type-safe API to access versions
3. **Lazy Loading**: Versions are loaded once and cached for performance
4. **Fallback Values**: Code includes fallback versions if config can't be loaded
5. **Update Script**: `scripts/update-versions.sh` fetches latest versions and updates the YAML file

## Benefits

- **Single Source of Truth**: Update versions in one place
- **Type Safety**: Compile-time checks for version access
- **Easy Updates**: Automated scripts to check and update versions
- **Consistency**: All components use the same version numbers
- **Maintainability**: No need to search and replace across multiple files

## Adding New Versions

To add a new dependency version:

1. Update `configs/versions.yaml` with the new dependency
2. Update the `Config` struct in `pkg/versions/versions.go`
3. Add version checking logic to `scripts/update-versions.sh`
4. Use the version in your code via `versions.Get()`

## Example: Adding a New Framework

```yaml
# In configs/versions.yaml
backend:
  frameworks:
    chi:
      version: "v5.0.0"
      package: "github.com/go-chi/chi/v5"
```

```go
// In pkg/versions/versions.go
type BackendVersions struct {
    // ... existing fields ...
    Frameworks struct {
        // ... existing frameworks ...
        Chi struct {
            Version string `yaml:"version"`
            Package string `yaml:"package"`
        } `yaml:"chi"`
    } `yaml:"frameworks"`
}
```

```go
// In your code
versionConfig, _ := versions.Get()
chiVersion := versionConfig.Backend.Frameworks.Chi.Version
```

## Troubleshooting

### Version Config Not Found

If you see errors about `versions.yaml` not being found:

1. Ensure you're running from the project root
2. Check that `configs/versions.yaml` exists
3. Use `versions.FindConfigPath()` to search for the file

### Outdated Versions

If `make check-versions` shows outdated versions:

1. Run `make update-versions` to update automatically
2. Or manually edit `configs/versions.yaml`
3. Test the changes with `make test`

### Adding to CI/CD

To ensure versions stay up-to-date in CI/CD:

```yaml
# In your CI workflow
- name: Check versions
  run: make check-versions

- name: Update versions (optional)
  run: make update-versions
```

## Related Files

- `configs/versions.yaml` - Version configuration
- `pkg/versions/versions.go` - Go package for accessing versions
- `pkg/versions/versions_test.go` - Tests
- `scripts/update-versions.sh` - Update script
- `scripts/check-latest-versions.sh` - Version checking script
