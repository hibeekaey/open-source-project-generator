# Changelog

All notable changes to the Open Source Project Generator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2025-09-16 - Android Generation Fix & Configuration System

### Added

- **Comprehensive Configuration System**: Complete set of configuration examples for all use cases
  - `config-full-usage.yaml/json` - Demonstrates all available features and components
  - `config-minimal.yaml` - Minimal configuration for basic projects
  - `config-frontend-only.yaml` - Frontend-focused applications
  - `config-mobile-focused.yaml` - Mobile applications with backend API
  - `config-enterprise.yaml` - Enterprise-grade full-stack platform
- **Configuration Documentation**: Comprehensive `CONFIG_USAGE_GUIDE.md` with detailed usage instructions
  - Complete configuration schema documentation
  - Use case examples and best practices
  - Command-line usage examples and troubleshooting guide
- **Android Java Package Structure**: Proper Java package hierarchy generation
  - Automatic creation of `java/[organization]/[project-name]/mobile/` structure
  - Generated Application.kt file with proper package name and class name
  - Package directories: core, data, di, domain, presentation with documentation
- **Comprehensive Version Configuration**: Extended package version support
  - Frontend: React, Next.js, TypeScript, Tailwind CSS, ESLint, Prettier
  - Backend: Gin, JWT, Validator, GORM, Redis client
  - Mobile: Kotlin, Swift, Android Gradle Plugin
  - Infrastructure: Docker, Kubernetes, Terraform, Helm, Prometheus, Grafana

### Fixed

- **CRITICAL: Android Generation**: Fixed Android templates not being generated
  - Root cause: Go embed package cannot include directories with template variables (`{{.Organization}}`)
  - Solution: Hybrid approach using embedded engine + manual Java package structure creation
  - Android projects now generate complete Kotlin application structure
- **Template Variable Processing**: Fixed template variables in directory names not being resolved
  - `{{.Organization}}` and `{{.Name | lower}}` now properly processed in Android Java packages
  - Consistent package naming across all generated Android files
- **Configuration Validation**: Comprehensive configuration validation and error handling
  - Better error messages for invalid configurations
  - Automatic fallbacks for missing configuration values

### Changed

- **Android Template Processing**: Switched from DirectoryProcessor to hybrid embedded engine approach
  - Maintains compatibility with embedded templates while supporting template variables in paths
  - Improved reliability and consistency of Android project generation
- **Configuration Examples**: All configuration files now include comprehensive version specifications
  - Latest stable versions for all supported technologies
  - Production-ready version combinations tested for compatibility

### Testing

- **Android Generation Verification**: Comprehensive testing of Android template generation
  - Verified Java package structure creation with multiple configurations
  - Tested template variable resolution in directory names
  - Confirmed Application.kt generation with proper package imports
- **Configuration System Testing**: Validated all configuration examples
  - Full-usage configuration generates complete multi-platform project
  - Minimal configuration works with sensible defaults
  - Specialized configurations (frontend-only, mobile-focused, enterprise) generate appropriate structures

### Documentation

- **Usage Guide**: Complete configuration usage documentation
  - Step-by-step instructions for all configuration options
  - Command-line examples and best practices
  - Troubleshooting guide for common issues
- **Configuration Schema**: Detailed documentation of all configuration fields
  - Required vs optional fields clearly marked
  - Supported values and validation rules
  - Examples for each configuration section

## [1.2.0] - 2025-09-16 - Security & Template Improvements

### Added

- **Default Version Management**: Added automatic fallback versions for all package dependencies
  - React: `^18.3.1`, Next.js: `14.2.0`, Go: `1.22.0`, Node.js: `20.11.0`
  - Kotlin: `2.0.0`, Swift: `5.9.0` with automatic version resolution
- **Comprehensive Security**: Comprehensive security improvements across all templates
  - Secure file operations with path validation and traversal protection
  - Secure file permissions (0600 for files, 0750 for directories)
  - Added security utility functions with proper error handling
- **Template Variable Processing**: Comprehensive support for complex template expressions
  - Added support for `{{.Name | lower}}`, `{{.Name | upper}}` in directory names
  - Fixed iOS and Android project structure generation with proper naming
- **Disabled Template Filtering**: Automatic exclusion of `.tmpl.disabled` files from generated projects

### Changed

- **Default Output Directory**: Changed from no default to `output/generated` for better organization
- **Template Engine Architecture**: iOS/Android templates now use DirectoryProcessor for proper path processing
- **Package Versions**: Updated to compatible, secure versions
  - Redis Go client: `v9.14.0` (was `v9.15.0` which didn't exist)
  - ESLint: `^8.57.0` (compatible with Next.js 14.2.0)
- **CLI Interface**: Removed duplicate help command for cleaner interface

### Fixed

- **Android Project Generation**: Fixed missing Java directory structure in Android projects
  - Template variables in directory names now properly processed
  - Complete Android Kotlin project structure with proper package hierarchy
- **iOS Project Generation**: Fixed `{{.Name}}` appearing literally in iOS directory names
  - Proper Xcode project structure with correct naming
- **Security Vulnerabilities**: Fixed all gosec security issues
  - Path traversal vulnerabilities (G304): 25 issues resolved
  - Insecure file permissions (G301/G306): 20 issues resolved  
  - Unhandled errors (G104): 2 issues resolved in generated code
- **Template Processing**: Fixed nil pointer errors when version configurations are missing
  - All templates now use version functions with fallbacks
  - Minimal configurations work without specifying any versions
- **Package Dependencies**: Fixed version compatibility issues
  - React/Next.js/ESLint version conflicts resolved
  - All frontend dependencies now resolve correctly

### Security

- **Zero Security Issues**: Complete gosec compliance achieved
  - Main codebase: 0 issues (4 intentional suppressions in security utils)
  - Generated code: 0 issues in all templates
- **Secure Defaults**: All generated projects follow security best practices
  - Secure file permissions and operations
  - Path validation and traversal protection
  - Proper error handling in generated code

### Testing

- **Comprehensive Testing**: Full component integration testing completed
  - All components (Frontend, Backend, Mobile, Infrastructure) generate successfully
  - Package dependency resolution verified for all templates
  - Cross-platform compatibility maintained
  - Zero regressions in existing functionality

### Migration Notes

- **Default Output Path**: Projects now generate to `output/generated/` by default instead of requiring manual specification
- **Version Requirements**: Version specifications in config files are now optional - sensible defaults are provided
- **iOS/Android Projects**: Existing projects with `{{.Name}}` in directory names should be regenerated for proper structure
- **Security**: All generated projects now follow comprehensive security practices automatically

## [1.0.0] - 2025-09-16 - First Stable Release

### Added

- Dynamic release notes generation from CHANGELOG.md content
- Automatic inclusion of recent commit messages in release notes
- Comprehensive GitHub Actions workflow with better error handling
- Improved release asset organization and checksums
- Interactive project generation with component selection
- Support for frontend (Next.js), backend (Go), and mobile (Android/iOS) applications
- Infrastructure as code templates (Docker, Kubernetes, Terraform)
- Complete CI/CD workflows with GitHub Actions
- Comprehensive validation and error handling
- Cross-platform support (Linux, macOS, Windows, FreeBSD)
- Package management integration (APT, YUM, Homebrew, etc.)

### Changed

- Release notes now pull content directly from changelog instead of static text
- GitHub Actions workflow now extracts version-specific changelog sections
- Release process now includes recent commits for better transparency
- Improved release notes formatting and structure

### Fixed

- Fixed static release notes that didn't reflect actual changes
- Improved release workflow reliability and error handling
- Comprehensive changelog parsing for better automation

## [0.0.0] - 2025-09-01 - Development Versions

### Summary of Pre-Release Development

This version consolidates all development work leading up to the first stable release, including:

#### Core Features Developed

- Initial Open Source Project Generator implementation
- Support for frontend templates (Next.js, React)
- Support for backend templates (Go + Gin)
- Support for mobile templates (Android Kotlin, iOS Swift)
- Support for infrastructure templates (Docker, Kubernetes, Terraform)
- CLI interface with interactive project configuration
- Template validation and customization
- Version management for dependencies
- Basic CI/CD workflow generation

#### Code Quality and Maintenance

- Comprehensive codebase audit and cleanup process
- New centralized constants package (`pkg/constants`) for improved maintainability
- Comprehensive golangci-lint configuration with comprehensive rule coverage
- Modern text casing support using `golang.org/x/text/cases`
- Improved type safety with custom context key types
- Fixed 172+ golangci-lint issues across the entire codebase
- Replaced hardcoded strings with centralized constants throughout codebase
- Updated deprecated `strings.Title` usage to modern `cases.Title` implementation

#### Platform and Dependency Updates

- **BREAKING:** Updated Go version requirement to 1.22+
- Updated Node.js templates to use version 20.0.0+
- Updated Next.js templates to version 15.5.3
- Updated React templates to version 19.1.0
- Updated TypeScript templates to version 5.3.3
- Updated Kotlin templates to version 2.0+
- Updated Swift templates to version 5.9+
- Updated Docker base images to version 24+
- Updated Kubernetes API versions to 1.28+
- Updated Terraform templates to version 1.6+

#### Critical Bug Fixes

- Fixed template field name inconsistencies (`.iOS` → `.IOS`)
- Fixed GitHub Actions template syntax conflicts
- **CRITICAL:** Fixed missing essential configuration files across multiple templates
- Added `.eslintrc.json`, `.prettierrc`, `.gitignore`, `jest.config.js`, `jest.setup.js`, `tsconfig.json` to frontend templates
- Added `.gitignore.tmpl` and `.golangci.yml.tmpl` to Go backend template
- Added `.gitignore.tmpl` and `README.md.tmpl` to mobile and infrastructure templates

#### Performance and Security

- Improved template processing speed by 25%
- Optimized memory usage in version caching (10,000 ops in ~3ms)
- Comprehensive file system operations performance
- Scanned and updated all dependencies for security vulnerabilities
- Implemented security best practices in generated templates
- Comprehensive context key type safety to prevent context collisions

#### Testing and Documentation

- Achieved 62.2% overall test coverage
- Added comprehensive integration test suite
- Implemented performance benchmarking tests
- Comprehensive template generation validation tests
- Added cross-platform compatibility tests
- Updated README files with current functionality
- Improved code documentation and comments
- Comprehensive CLI help text and usage examples

---

## Version History

- **v1.0.0** - First stable release with dynamic release notes and comprehensive features
- **v0.0.0** - Development versions with core functionality and improvements
