# Changelog

All notable changes to the Open Source Template Generator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
