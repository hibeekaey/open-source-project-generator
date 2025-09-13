# Templates Directory

This directory contains all the template files used by the Open Source Template Generator. All templates use the latest stable versions of their respective technologies and follow modern best practices.

## Structure

```
templates/
├── base/                    # Core project files and documentation
│   ├── .github/            # GitHub Actions workflows and templates
│   ├── docs/               # Documentation templates (BUILD_SYSTEM, DEPLOYMENT, etc.)
│   └── scripts/            # Build and deployment scripts
├── frontend/               # Frontend application templates (Node.js 20+, Next.js 15+)
│   ├── nextjs-app/         # Main application template with React 19+
│   ├── nextjs-home/        # Landing page template with modern design
│   ├── nextjs-admin/       # Admin dashboard with comprehensive UI components
│   └── shared-components/  # Reusable component library
├── backend/                # Backend service templates (Go 1.24+)
│   └── go-gin/             # Go + Gin API server with JWT, GORM, Redis
├── mobile/                 # Mobile application templates
│   ├── android-kotlin/     # Android Kotlin 2.0+ with Jetpack Compose
│   ├── ios-swift/          # iOS Swift 5.9+ with SwiftUI
│   └── shared/             # Shared mobile resources and API specs
├── infrastructure/         # Infrastructure templates (latest versions)
│   ├── terraform/          # Terraform 1.6+ configurations
│   ├── kubernetes/         # Kubernetes 1.28+ manifests with security policies
│   └── docker/             # Docker 24+ configurations with multi-stage builds
└── README.md               # This file
```

## Template Processing

Templates use Go's text/template syntax with custom functions for:

- **Variable substitution**: Project name, organization, description, etc.
- **Conditional rendering**: Based on selected components and features
- **Version management**: Automatic latest version injection for dependencies
- **String manipulation**: Case conversion, sanitization, and formatting
- **Security defaults**: Secure configurations and best practices

## Template Features

### Frontend Templates

- **Next.js 15+** with App Router and TypeScript 5.7+
- **React 19+** with latest hooks and concurrent features
- **Tailwind CSS 3.4+** with modern design system
- **Comprehensive testing** with Jest 29+ and Testing Library
- **ESLint 9+** and Prettier 3+ for code quality
- **Performance optimization** with proper bundling and caching

### Backend Templates

- **Go 1.24+** with latest language features and performance improvements
- **Gin framework** with middleware for logging, CORS, and security
- **GORM** for database operations with PostgreSQL support
- **JWT authentication** with secure token handling
- **Redis integration** for caching and session management
- **Comprehensive testing** with testify and integration tests
- **API documentation** with Swagger/OpenAPI 3.0

### Mobile Templates

- **Android**: Kotlin 2.0+ with Jetpack Compose and Material Design 3
- **iOS**: Swift 5.9+ with SwiftUI and modern iOS architecture patterns
- **Security**: Proper keychain/keystore usage and network security
- **Testing**: Unit and UI tests with proper mocking
- **CI/CD**: Automated building and testing workflows

### Infrastructure Templates

- **Docker**: Multi-stage builds with security scanning and non-root users
- **Kubernetes**: Proper resource limits, security policies, and health checks
- **Terraform**: Multi-cloud support with proper state management
- **Monitoring**: Prometheus and Grafana configurations
- **Security**: Network policies, RBAC, and secret management

## Version Management

All templates automatically use the latest stable versions:

- **Go**: 1.24+ (latest stable)
- **Node.js**: 20+ LTS
- **Next.js**: 15+ (latest stable)
- **React**: 19+ (latest stable)
- **Kotlin**: 2.0+ (latest stable)
- **Swift**: 5.9+ (latest stable)
- **Docker**: 24+ (latest stable)
- **Kubernetes**: 1.28+ (latest stable)
- **Terraform**: 1.6+ (latest stable)

## Adding New Templates

1. **Create template structure** in the appropriate directory following existing patterns
2. **Use proper template syntax** with Go text/template and custom functions
3. **Include version variables** for automatic dependency management
4. **Add comprehensive documentation** within the template
5. **Include security best practices** and proper configurations
6. **Add validation logic** for template-specific requirements
7. **Test template generation** with various configurations
8. **Update this README** with new template information

## Template Validation

All templates include:

- **Syntax validation** for generated code
- **Dependency verification** for all package versions
- **Security scanning** for vulnerabilities and best practices
- **Build verification** to ensure generated projects compile/build successfully
- **Test validation** to ensure all generated tests pass
