# Changelog

All notable changes to the Open Source Project Generator will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - TBD - Tool-Orchestration Architecture

### Breaking Changes

- Old template-based generation system moved to `deprecated/` directory
- Configuration schema updated to component-based format with `components`, `integration`, and `options` fields
- Removed `--template` flag (replaced with component-based configuration)
- Internal package structure reorganized:
  - Deprecated: `deprecated/pkg/template/`, `deprecated/pkg/cli/`, `deprecated/internal/app/`
  - New: `internal/orchestrator/`, `internal/generator/bootstrap/`, `internal/generator/fallback/`, `internal/generator/mapper/`, `pkg/logger/`

### Added

- Tool-orchestration system using industry-standard CLI tools (`create-next-app`, `go mod init`)
- Bootstrap tool executors for Next.js, Go, Android, and iOS with real-time output streaming
- Fallback generators for when external tools are unavailable
- Tool discovery and validation system with version checking and OS-specific installation instructions
- Enhanced logging with file output, timestamps, structured context, and colored terminal output
- Backup and rollback system with timestamped backups in `.backups/` directory
- Security enhancements: input sanitization, path validation, configuration value sanitization
- Offline support with tool caching and TTL-based cache invalidation
- Parallel component generation with configurable worker pool
- Integration features: Docker Compose generation, shared `.env` files, build scripts, comprehensive documentation
- New CLI commands: `check-tools`, `init-config`
- New CLI flags: `--no-external-tools`, `--stream-output`

### Changed

- Project generation now delegates to official framework tools instead of maintaining templates
- Configuration format updated to support multiple components per project
- Error handling improved with categorization, recovery strategies, and retry logic

### Migration Notes

- Existing projects generated with v1.x will continue to work
- New projects should use component-based configuration format
- Install required external tools (Node.js/npm, Go 1.25+) or use `--no-external-tools` flag
- Use `generator init-config` to generate new configuration templates
- Android/iOS generation requires native tooling or will use minimal fallback structure

## [1.5.0] - 2025-10-02 - Comprehensive Code Quality and Security Enhancement

### Added

- **Comprehensive Code Quality Fixes**: Complete resolution of all code quality issues
  - Fixed 88 linting issues across the entire codebase (errcheck, staticcheck, unused code, spelling, go vet)
  - Implemented consistent error handling patterns throughout all packages
  - Added comprehensive bounds checking for all numeric operations
  - Enhanced code documentation and inline comments
- **Complete Security Hardening**: Zero security vulnerabilities achieved
  - Fixed all 55 security issues identified by gosec security scanner
  - Implemented restrictive file permissions (0600 for files, 0750 for directories)
  - Added comprehensive path validation and sanitization for all file operations
  - Enhanced error handling security to prevent information leakage
- **Test Suite Stabilization**: 100% test reliability achieved
  - Fixed all failing test suites including CLI interactive, filesystem, integration, and security tests
  - Implemented proper mocking for non-interactive environments
  - Added comprehensive synchronization mechanisms for concurrent operations
  - Enhanced test isolation and cleanup procedures
- **Enhanced Error Handling**: Comprehensive error management improvements
  - Added proper error handling for all critical code paths
  - Implemented context-aware error messages and recovery mechanisms
  - Enhanced error categorization and logging throughout the codebase
  - Added security-focused error sanitization
- **Code Quality Improvements**: Significant maintainability enhancements
  - Removed unused code and optimized implementations
  - Refactored dependency injection container for better efficiency
  - Implemented code deduplication and created reusable utility functions
  - Enhanced performance optimizations and reduced memory allocations

### Changed

- **Container Architecture**: Streamlined dependency injection container
  - Removed unused fields and components for better performance
  - Refactored container interfaces for improved clarity
  - Enhanced dependency management and initialization
- **Security Posture**: Elevated security standards across all components
  - Updated file creation permissions to be more restrictive
  - Enhanced input validation and sanitization throughout
  - Improved secure file operations with comprehensive path checking
- **Test Infrastructure**: Enhanced testing reliability and coverage
  - Improved test environment detection and handling
  - Enhanced integration test stability with proper synchronization
  - Added comprehensive edge case and error condition testing
- **Performance Optimization**: Improved efficiency across all operations
  - Optimized memory usage calculations and algorithm efficiency
  - Enhanced caching mechanisms and reduced unnecessary allocations
  - Improved concurrent operation handling and thread safety

### Fixed

- **Critical Test Failures**: Resolved all test suite failures
  - CLI interactive tests now handle non-interactive environments properly
  - Filesystem generator tests properly validate nil configurations
  - Integration tests no longer have race conditions or concurrency issues
  - Security tests handle timing and environment variations correctly
- **Security Vulnerabilities**: Complete security issue resolution
  - Fixed integer overflow vulnerabilities in memory calculations
  - Resolved path traversal vulnerabilities with comprehensive validation
  - Enhanced file permission security across all file operations
  - Improved error handling to prevent information disclosure
- **Code Quality Issues**: Comprehensive linting compliance achieved
  - Fixed all errcheck issues with proper error handling
  - Resolved staticcheck violations and code quality improvements
  - Removed all unused code and variables throughout codebase
  - Corrected spelling issues and improved code clarity
- **Performance Issues**: Optimized resource usage and efficiency
  - Fixed memory usage calculation overflows
  - Improved algorithm efficiency and reduced allocations
  - Enhanced concurrent operation performance and safety

### Security

- **Zero Security Issues**: Complete gosec compliance achieved
  - All 55 previously identified security issues resolved
  - Enhanced security patterns implemented across all templates
  - Comprehensive input validation and path sanitization
  - Secure file operations with appropriate permissions
- **Security Best Practices**: Enhanced security implementation
  - Restrictive file permissions for all generated and internal files
  - Comprehensive path validation against traversal attacks
  - Secure error handling preventing information leakage
  - Enhanced input sanitization throughout the application

### Testing

- **100% Test Reliability**: Complete test suite stabilization
  - All test suites now pass consistently in all environments
  - Enhanced test coverage for edge cases and error conditions
  - Improved integration testing with proper isolation
  - Comprehensive security testing and validation
- **Quality Assurance**: Enhanced testing infrastructure
  - Added regression tests for all fixed issues
  - Implemented comprehensive validation testing
  - Enhanced performance testing and benchmarking
  - Improved cross-platform compatibility testing

### Documentation

- **Comprehensive Documentation Updates**: Enhanced project documentation
  - Created detailed code quality fixes summary documentation
  - Updated troubleshooting guide with new solutions and fixes
  - Enhanced API documentation and code comments
  - Added migration guide for the quality improvements
- **Quality Metrics**: Documented quality improvements and metrics
  - Complete test failure resolution (0 failures from ~20)
  - Complete linting compliance (0 issues from 88)
  - Complete security compliance (0 issues from 55)
  - Enhanced code coverage and quality metrics

### Migration Notes

- **No Breaking Changes**: Full backward compatibility maintained
  - All existing functionality preserved without changes
  - API compatibility maintained across all interfaces
  - Configuration format unchanged for existing users
  - Command-line interface remains consistent
- **Immediate Benefits**: Enhanced reliability and security
  - Improved error handling and recovery mechanisms
  - Better performance and resource efficiency
  - Enhanced security posture with no user action required
  - More reliable operation across all environments

## [1.4.0] - 2025-09-26 - Enhanced Interactive UI & Advanced Features

### Added

- **Enhanced Interactive UI System**: Complete overhaul of the user interface with advanced navigation and help systems
  - Navigation system with breadcrumbs and step counter
  - Context-sensitive help system with comprehensive documentation
  - Enhanced project structure preview with interactive navigation
  - Template preview system with detailed component information
  - Scrollable menu system for better user experience
  - Completion system for intelligent input suggestions
- **Advanced Audit Engine**: Comprehensive project auditing and quality analysis
  - Security pattern detection and validation
  - Code quality analysis with customizable rules
  - Compliance checking for generated projects
  - Automated security vulnerability scanning
  - Performance metrics and optimization recommendations
- **Intelligent Cache Management**: Advanced caching system for improved performance
  - Offline mode support with intelligent cache management
  - Version-aware caching with automatic invalidation
  - Cache compression and optimization
  - Metrics collection and performance monitoring
  - Automatic cleanup and maintenance
- **Enhanced Validation Engine**: Advanced validation system with auto-fix capabilities
  - Template structure validation with detailed error reporting
  - Configuration validation with intelligent suggestions
  - Auto-fix capabilities for common issues
  - Integration testing with comprehensive coverage
  - Performance benchmarking and optimization
- **Security & Safety Features**: Comprehensive security improvements
  - Input sanitization and validation
  - Secure file operations with path validation
  - Backup management with version control
  - Dry-run mode for safe testing
  - Concurrent operation safety
- **Advanced Template Management**: Enhanced template processing and management
  - Template import detection with enhanced algorithms
  - Template security pattern validation
  - Comprehensive template testing suite
  - Template fix automation for common issues
  - Enhanced template scanning and processing
- **Version & Update Management**: Intelligent version management system
  - GitHub integration for version checking
  - Automatic update notifications
  - Version compatibility validation
  - Dependency version management
  - Release note generation and management
- **Enhanced Error Handling**: Comprehensive error management system
  - Categorized error handling with detailed context
  - Error recovery mechanisms
  - Security-focused error handling
  - Logging integration with structured output
  - User-friendly error messages and suggestions

### Changed

- **CLI Interface**: Complete redesign of command-line interface
  - Enhanced automation capabilities with intelligent workflows
  - Improved interactive flow management
  - Better mode detection and switching
  - Enhanced configuration commands
  - Streamlined user experience with better error handling
- **Template Processing**: Advanced template processing engine
  - Enhanced embedded template system
  - Improved template function processing
  - Better template security validation
  - Optimized template scanning and loading
- **Configuration Management**: Enhanced configuration system
  - Improved configuration persistence
  - Better configuration validation
  - Enhanced configuration testing
  - Streamlined configuration management
- **Filesystem Operations**: Advanced filesystem management
  - Enhanced project generation with standardized structure
  - Improved filesystem security
  - Better project structure validation
  - Optimized file operations

### Fixed

- **Template Generation**: Fixed multiple template generation issues
  - Android Kotlin template structure improvements
  - iOS Swift template enhancements
  - Backend Go template optimizations
  - Frontend Next.js template improvements
  - Infrastructure template fixes
- **Security Vulnerabilities**: Resolved all security issues
  - Path traversal protection
  - Input validation improvements
  - Secure file permission handling
  - Enhanced error handling security
- **Performance Issues**: Optimized performance across all components
  - Cache performance improvements
  - Template processing optimization
  - Memory usage optimization
  - Concurrent operation improvements

### Security

- **Zero Security Issues**: Maintained complete gosec compliance
  - Enhanced security patterns in all templates
  - Improved input validation and sanitization
  - Secure file operations with proper permissions
  - Enhanced error handling security
- **Security Best Practices**: All generated projects follow enhanced security practices
  - Secure default configurations
  - Input validation and sanitization
  - Path traversal protection
  - Secure file permissions

### Testing

- **Comprehensive Testing Suite**: Enhanced testing coverage and quality
  - Integration testing for all major components
  - Performance benchmarking and optimization
  - Security testing and validation
  - Cross-platform compatibility testing
  - User experience testing
- **Quality Assurance**: Improved code quality and reliability
  - Enhanced error handling testing
  - Template validation testing
  - Configuration testing
  - Performance testing

### Documentation

- **Enhanced Documentation**: Comprehensive documentation updates
  - API reference documentation
  - Configuration guide improvements
  - Getting started guide enhancements
  - Template development documentation
  - Troubleshooting guide updates
- **User Experience**: Improved user experience documentation
  - Interactive UI documentation
  - Configuration examples
  - Best practices guide
  - Migration guides

## [1.3.0] - 2025-09-16 - Android Generation Fix & Configuration System

### Added

- **Comprehensive Configuration System**: Complete set of configuration examples for all use cases
  - `config-full-usage.yaml/json` - Demonstrates all available features and components
  - `config-minimal.yaml` - Minimal configuration for simple projects
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
- **Enhanced Version Configuration**: Extended package version support
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
- **Configuration Validation**: Enhanced configuration validation and error handling
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

- **Default Version Management**: Updated to latest stable versions for all package dependencies
  - React: `^19.2.0`, Next.js: `16.0.0`, Go: `1.25.0`, Node.js: `20.11.0`
  - Kotlin: `2.1.0`, Swift: `6.0` with automatic version resolution
  - Android: SDK 35, Gradle 8.11.1, AndroidX Core KTX 1.17.0
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
- **Security**: All generated projects now follow enhanced security practices automatically

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
- Enhanced golangci-lint configuration with comprehensive rule coverage
- Modern text casing support using `golang.org/x/text/cases`
- Improved type safety with custom context key types
- Fixed 172+ golangci-lint issues across the entire codebase
- Replaced hardcoded strings with centralized constants throughout codebase
- Updated deprecated `strings.Title` usage to modern `cases.Title` implementation

#### Platform and Dependency Updates

- **BREAKING:** Updated Go version requirement to 1.25+
- Updated Node.js templates to use version 20.11.0+
- Updated Next.js templates to version 16.0.0
- Updated React templates to version 19.2.0
- Updated TypeScript templates to version 5.9.3
- Updated Kotlin templates to version 2.1.0
- Updated Swift templates to version 6.0
- Updated Android SDK to API level 35
- Updated Gradle to version 8.11.1
- Updated AndroidX libraries to latest stable versions
- Updated Go backend frameworks (Gin 1.11.0, Echo 4.13.4, Fiber 2.52.9)
- Updated Docker base images to alpine:3.19, golang:1.25-alpine
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
