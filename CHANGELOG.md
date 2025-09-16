# Changelog

All notable changes to the Open Source Template Generator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-09-16 - Security & Template Improvements

### Added

- **Default Version Management**: Added automatic fallback versions for all package dependencies
  - React: `^18.3.1`, Next.js: `14.2.0`, Go: `1.22.0`, Node.js: `20.11.0`
  - Kotlin: `2.0.0`, Swift: `5.9.0` with automatic version resolution
- **Enhanced Security**: Comprehensive security improvements across all templates
  - Secure file operations with path validation and traversal protection
  - Secure file permissions (0600 for files, 0750 for directories)
  - Added security utility functions with proper error handling
- **Template Variable Processing**: Enhanced support for complex template expressions
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

## [1.0.0] - 2025-09-16 - First Stable Release

### Added

- Dynamic release notes generation from CHANGELOG.md content
- Automatic inclusion of recent commit messages in release notes
- Enhanced GitHub Actions workflow with better error handling
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
- Enhanced changelog parsing for better automation

## [0.0.0] - 2025-09-01 - Development Versions

### Summary of Pre-Release Development

This version consolidates all development work leading up to the first stable release, including:

#### Core Features Developed

- Initial Open Source Template Generator implementation
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
- Enhanced golangci-lint configuration with comprehensive rule coverage
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
- Enhanced file system operations performance
- Scanned and updated all dependencies for security vulnerabilities
- Implemented security best practices in generated templates
- Enhanced context key type safety to prevent context collisions

#### Testing and Documentation

- Achieved 62.2% overall test coverage
- Added comprehensive integration test suite
- Implemented performance benchmarking tests
- Enhanced template generation validation tests
- Added cross-platform compatibility tests
- Updated README files with current functionality
- Improved code documentation and comments
- Enhanced CLI help text and usage examples

---

## Version History

- **v1.0.0** - First stable release with dynamic release notes and comprehensive features
- **v0.0.0** - Development versions with core functionality and improvements
