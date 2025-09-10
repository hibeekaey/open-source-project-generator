# Open Source Project Design Specification

## Table of Contents

1. [Purpose and Intent](#purpose-and-intent)
2. [Directory Structure](#recommended-directory-structure)
3. [Technology Stack](#technology-stack-patterns)
4. [Package Versions](#current-package-versions)
5. [Configuration Files](#key-configuration-files)
6. [Container Configurations](#container-configurations)
7. [GitHub Actions](#github-actions-workflows)
8. [Deployment Configurations](#frontend-deployment-configurations)
9. [Documentation Templates](#essential-documentation-templates)
10. [Benefits](#benefits-of-this-structure)
11. [Implementation Guide](#implementation-checklist)
12. [Best Practices](#best-practices-compliance)
13. [Customization](#customization-notes)

## Purpose and Intent

This document serves as a **comprehensive design specification** for creating production-ready open source projects. It establishes a standardized, enterprise-grade foundation that can be used as a template for new projects or as a migration target for existing codebases.

### Design Goals

This specification is intended to be used by:

- **Development teams** starting new open source projects
- **Organizations** standardizing their project structure across repositories
- **Maintainers** upgrading existing projects to modern best practices
- **Contributors** understanding project organization and development workflows

### Scope and Coverage

This design specification provides:

1. **Complete Directory Structure** - A battle-tested layout supporting multi-service architecture
2. **Build System Design** - Make-based automation for cross-platform development
3. **Deployment Architecture** - Production-ready configurations for modern platforms
4. **Security Framework** - Comprehensive security policies and automated scanning
5. **Testing Strategy** - Multi-level testing across all project components
6. **Documentation Standards** - Professional open source governance and documentation

### Architecture Philosophy

The design follows a **modern multi-service architecture** pattern derived from analysis of production applications, emphasizing:

- **Monorepo structure** with clear service boundaries
- **Platform-specific deployment** optimized for each environment
- **Security-first approach** with automated vulnerability detection
- **Developer experience** through comprehensive tooling and automation
- **Production readiness** from day one with enterprise-grade practices
- **Open source compliance** meeting all community standards and best practices

## Recommended Directory Structure

```bash
project-name/
├── README.md                          # Main project overview
├── LICENSE                           # Open source license
├── CHANGELOG.md                      # Version history and release notes
├── CONTRIBUTING.md                   # Contribution guidelines
├── SECURITY.md                       # Security policy and vulnerability reporting
├── CODEOWNERS                        # Code ownership and review assignments
├── VERSION                           # Current version file
├── .gitignore                        # Git ignore rules
├── .gitattributes                   # Git attributes for consistent line endings
├── .editorconfig                    # Editor configuration for consistency
├── docker-compose.yml               # Local development orchestration
├── docker-compose.prod.yml          # Production docker compose
├── Makefile                         # Primary build system and automation
├── .env.example                     # Environment variables template
├── .dockerignore                    # Docker ignore rules
├── 
├── App/                             # Frontend application(s)
│   ├── Dockerfile                   # Frontend container configuration
│   ├── package.json                # Frontend dependencies
│   ├── package-lock.json           # Locked dependency versions
│   ├── next.config.js              # Next.js configuration with static export
│   ├── tsconfig.json               # TypeScript configuration
│   ├── tailwind.config.ts          # Tailwind CSS configuration
│   ├── postcss.config.js           # PostCSS configuration
│   ├── jest.config.js              # Jest testing configuration
│   ├── Makefile                    # Frontend build automation
│   ├── README.md                   # Frontend-specific documentation
│   ├── .env.example               # Frontend environment variables
│   ├── vercel.json                 # Vercel deployment configuration
│   ├── app.yaml                    # DigitalOcean App Platform configuration
│   ├── public/                     # Static assets
│   │   ├── favicon.ico
│   │   ├── logo.png
│   │   └── images/
│   ├── src/                        # Source code
│   │   ├── app/                    # Next.js app router (or main app)
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   ├── globals.css
│   │   │   └── [features]/
│   │   ├── components/             # Reusable UI components
│   │   │   ├── ui/                 # Base UI components
│   │   │   ├── layouts/            # Layout components
│   │   │   └── [feature-name]/     # Feature-specific components
│   │   ├── hooks/                  # Custom React hooks
│   │   ├── context/                # React context providers
│   │   ├── lib/                    # Utility libraries and configurations
│   │   ├── types/                  # TypeScript type definitions
│   │   ├── utils/                  # Helper functions
│   │   ├── styles/                 # Global styles and themes
│   │   ├── assets/                 # Local assets (images, icons)
│   │   │   ├── images/
│   │   │   ├── icons/
│   │   │   └── logos/
│   │   └── API/                    # API client code
│   │       ├── types.ts
│   │       ├── axiosSetup.ts
│   │       └── APIS/
│   └── __tests__/                  # Frontend tests
│       └── index.test.tsx
│
├── CommonServer/                    # Backend API server
│   ├── Dockerfile                  # Backend container configuration
│   ├── Dockerfile.dev              # Development container configuration
│   ├── go.mod                      # Go module definition
│   ├── go.sum                      # Go dependency checksums
│   ├── main.go                     # Application entry point
│   ├── Makefile                    # Build and development commands
│   ├── README.md                   # Backend-specific documentation
│   ├── .env.example               # Backend environment variables
│   ├── k8s/                        # Kubernetes manifests
│   │   ├── namespace.yaml          # Kubernetes namespace
│   │   ├── deployment.yaml         # Application deployment
│   │   ├── service.yaml            # Service definition
│   │   ├── configmap.yaml          # Configuration management
│   │   ├── secret.yaml.example     # Secrets template
│   │   └── hpa.yaml                # Horizontal Pod Autoscaler
│   ├── configs/                    # Configuration management
│   │   ├── database.go
│   │   └── env.go
│   ├── controllers/                # Request handlers
│   │   ├── userController.go
│   │   └── [entity]Controller.go
│   ├── models/                     # Data models
│   │   ├── userModel.go
│   │   └── [entity]Model.go
│   ├── routes/                     # API route definitions
│   │   ├── routes.go
│   │   └── [entity]Router.go
│   ├── middleware/                 # HTTP middleware
│   │   ├── authMiddleware.go
│   │   └── corsMiddleware.go
│   ├── services/                   # Business logic
│   │   ├── userServices.go
│   │   └── main_test.go
│   ├── repository/                 # Data access layer
│   │   └── [entity].go
│   ├── utils/                      # Utility functions
│   │   ├── jwt.go
│   │   ├── helpers.go
│   │   └── mail.go
│   ├── database/                   # Database configuration
│   │   └── database.go
│   ├── proto/                      # Protocol Buffer definitions (if using gRPC)
│   │   ├── user.proto
│   │   ├── user.pb.go
│   │   └── user_grpc.pb.go
│   └── seeds/                      # Database seeders
│       └── seed.go
│
├── Home/                           # Landing page/marketing site
│   ├── Dockerfile                  # Home container configuration
│   ├── package.json                # Home dependencies
│   ├── package-lock.json           # Locked dependency versions
│   ├── next.config.js              # Next.js configuration with static export
│   ├── tsconfig.json               # TypeScript configuration
│   ├── tailwind.config.ts          # Tailwind CSS configuration
│   ├── postcss.config.js           # PostCSS configuration
│   ├── jest.config.js              # Jest testing configuration
│   ├── Makefile                    # Home build automation
│   ├── README.md                   # Home-specific documentation
│   ├── .env.example               # Home environment variables
│   ├── vercel.json                 # Vercel deployment configuration
│   ├── app.yaml                    # DigitalOcean App Platform configuration
│   ├── public/                     # Static assets
│   │   ├── favicon.ico
│   │   ├── logo.png
│   │   └── images/
│   ├── src/                        # Source code
│   │   ├── app/                    # Next.js app router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   └── globals.css
│   │   ├── components/             # Reusable UI components
│   │   │   ├── ui/                 # Base UI components
│   │   │   ├── home/               # Home-specific components
│   │   │   └── shared/             # Shared components
│   │   ├── lib/                    # Utility libraries
│   │   ├── types/                  # TypeScript type definitions
│   │   ├── utils/                  # Helper functions
│   │   ├── styles/                 # Global styles
│   │   └── assets/                 # Local assets
│   │       ├── images/
│   │       ├── icons/
│   │       └── logos/
│   └── __tests__/                  # Home tests
│       └── index.test.tsx
│
├── Admin/                          # Admin dashboard
│   ├── Dockerfile                  # Admin container configuration
│   ├── package.json                # Admin dependencies
│   ├── package-lock.json           # Locked dependency versions
│   ├── next.config.js              # Next.js configuration
│   ├── tsconfig.json               # TypeScript configuration
│   ├── tailwind.config.ts          # Tailwind CSS configuration
│   ├── postcss.config.js           # PostCSS configuration
│   ├── jest.config.js              # Jest testing configuration
│   ├── Makefile                    # Admin build automation
│   ├── README.md                   # Admin-specific documentation
│   ├── .env.example               # Admin environment variables
│   ├── vercel.json                 # Vercel deployment configuration
│   ├── app.yaml                    # DigitalOcean App Platform configuration
│   ├── public/                     # Static assets
│   │   ├── favicon.ico
│   │   ├── logo.png
│   │   └── images/
│   ├── src/                        # Source code
│   │   ├── app/                    # Next.js app router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx
│   │   │   ├── globals.css
│   │   │   ├── dashboard/          # Dashboard pages
│   │   │   ├── users/              # User management
│   │   │   └── settings/           # Admin settings
│   │   ├── components/             # Reusable UI components
│   │   │   ├── ui/                 # Base UI components
│   │   │   ├── admin/              # Admin-specific components
│   │   │   ├── forms/              # Form components
│   │   │   └── tables/             # Data table components
│   │   ├── lib/                    # Utility libraries
│   │   ├── types/                  # TypeScript type definitions
│   │   ├── utils/                  # Helper functions
│   │   ├── styles/                 # Global styles
│   │   ├── hooks/                  # Custom React hooks
│   │   ├── context/                # React context providers
│   │   ├── API/                    # API client code
│   │   │   ├── types.ts
│   │   │   ├── axiosSetup.ts
│   │   │   └── admin/              # Admin API endpoints
│   │   └── assets/                 # Local assets
│   │       ├── images/
│   │       ├── icons/
│   │       └── logos/
│   └── __tests__/                  # Admin tests
│       └── index.test.tsx
│
├── Mobile/                         # Native mobile applications
│   ├── README.md                  # Mobile development documentation
│   ├── shared/                    # Shared resources between platforms
│   │   ├── assets/               # Common assets (images, fonts)
│   │   │   ├── images/
│   │   │   ├── icons/
│   │   │   └── fonts/
│   │   ├── api/                  # Shared API specifications
│   │   │   ├── swagger.yaml
│   │   │   └── endpoints.md
│   │   └── design/               # Design system resources
│   │       ├── colors.md
│   │       ├── typography.md
│   │       └── components.md
│   ├── android/                  # Android native application
│   │   ├── app/
│   │   │   ├── build.gradle
│   │   │   ├── proguard-rules.pro
│   │   │   └── src/
│   │   │       ├── main/
│   │   │       │   ├── AndroidManifest.xml
│   │   │       │   ├── kotlin/com/company/projectname/
│   │   │       │   │   ├── MainActivity.kt
│   │   │       │   │   ├── models/
│   │   │       │   │   ├── activities/
│   │   │       │   │   ├── fragments/
│   │   │       │   │   ├── adapters/
│   │   │       │   │   ├── services/
│   │   │       │   │   ├── utils/
│   │   │       │   │   └── api/
│   │   │       │   └── res/
│   │   │       │       ├── layout/
│   │   │       │       ├── values/
│   │   │       │       ├── drawable/
│   │   │       │       └── mipmap/
│   │   │       ├── test/          # Unit tests
│   │   │       └── androidTest/   # Instrumentation tests
│   │   ├── build.gradle           # Project-level build configuration
│   │   ├── gradle.properties      # Gradle properties
│   │   ├── settings.gradle        # Project settings
│   │   ├── gradle/
│   │   │   └── wrapper/
│   │   └── README.md              # Android-specific documentation
│   └── ios/                      # iOS native application
│       ├── ProjectName/
│       │   ├── AppDelegate.swift
│       │   ├── SceneDelegate.swift
│       │   ├── ViewController.swift
│       │   ├── Info.plist
│       │   ├── Models/
│       │   ├── Views/
│       │   │   ├── Storyboards/
│       │   │   │   ├── Main.storyboard
│       │   │   │   └── LaunchScreen.storyboard
│       │   │   └── XIBs/
│       │   ├── Controllers/
│       │   ├── Services/
│       │   ├── Utils/
│       │   ├── API/
│       │   └── Resources/
│       │       ├── Assets.xcassets/
│       │       ├── Fonts/
│       │       └── Localizable.strings
│       ├── ProjectName.xcodeproj/
│       │   └── project.pbxproj
│       ├── ProjectNameTests/       # Unit tests
│       ├── ProjectNameUITests/     # UI tests
│       ├── Podfile                # CocoaPods dependencies
│       ├── Podfile.lock           # Locked dependency versions
│       └── README.md              # iOS-specific documentation
│
├── Deploy/                         # Infrastructure as Code
│   ├── main.tf                    # Main Terraform configuration
│   ├── providers.tf               # Terraform providers
│   ├── variables.tf               # Terraform variables
│   └── Chart/                     # Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── _helpers.tpl
│           ├── deployment.yaml
│           └── service.yaml
│
├── Docs/                          # Documentation
│   ├── README.md                  # Documentation overview
│   ├── self_hosting.md           # Self-hosting guide
│   ├── assets/                   # Documentation images
│   │   └── [screenshots]
│   ├── Backend/                  # API documentation
│   │   ├── ApiDocumentation.md
│   │   ├── User.md
│   │   └── [Entity].md
│   └── frontend/                 # Frontend documentation
│       ├── overview.md
│       ├── auth.md
│       └── [feature].md
│
├── Tests/                         # Integration and E2E tests
│   ├── e2e/                      # End-to-end tests
│   │   ├── Dockerfile           # E2E test environment
│   │   ├── package.json
│   │   └── tests/
│   ├── integration/             # Integration tests
│   │   ├── api/                 # API integration tests
│   │   └── mobile/              # Mobile integration tests
│   ├── performance/             # Performance tests
│   │   ├── load/               # Load testing
│   │   └── stress/             # Stress testing
│   └── README.md               # Testing documentation
│
├── Scripts/                       # Build and utility scripts
│   ├── setup.sh                  # Development environment setup
│   ├── build.sh                  # Build script
│   ├── deploy.sh                 # Deployment script
│   ├── test.sh                   # Testing script
│   ├── lint.sh                   # Linting script
│   └── clean.sh                  # Cleanup script
│
├── .github/                      # GitHub-specific files
│   ├── workflows/                # GitHub Actions
│   │   ├── ci.yml               # Continuous Integration
│   │   ├── cd.yml               # Continuous Deployment
│   │   ├── test.yml             # Automated testing
│   │   ├── security.yml         # Security scanning
│   │   ├── release.yml          # Release automation
│   │   ├── dependency-review.yml # Dependency security review
│   │   └── mobile.yml           # Mobile app CI/CD
│   ├── ISSUE_TEMPLATE/          # Issue templates
│   │   ├── bug_report.md
│   │   ├── feature_request.md
│   │   ├── question.md
│   │   └── config.yml           # Issue template config
│   ├── PULL_REQUEST_TEMPLATE.md # PR template
│   ├── SECURITY.md              # Security policy
│   └── dependabot.yml           # Dependency updates
│
└── .devcontainer/               # Development container configuration
    ├── devcontainer.json        # VS Code dev container config
    └── Dockerfile               # Development environment
```

## Technology Stack Patterns

### Frontend Applications

- **Framework**: Next.js 15+ with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS + Styled Components
- **State Management**: React Context API
- **Testing**: Jest + React Testing Library
- **Animation**: Framer Motion
- **HTTP Client**: Axios

### Backend Services

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP) / gRPC
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Documentation**: Protocol Buffers (for gRPC)

### Mobile Applications

#### Android

- **Language**: Kotlin 2.0+
- **Architecture**: MVVM with Android Architecture Components
- **Dependency Injection**: Dagger Hilt
- **Networking**: Retrofit + OkHttp 4.12+
- **Database**: Room
- **UI**: Jetpack Compose + Material Design 3
- **Testing**: JUnit 5 + Espresso
- **Build**: Gradle 8.5+ with Version Catalogs

#### iOS

- **Language**: Swift 5.9+
- **Architecture**: MVVM with Combine
- **Dependency Injection**: Swinject
- **Networking**: URLSession + Alamofire 5.8+
- **Database**: SwiftData + Core Data
- **UI**: SwiftUI + UIKit
- **Testing**: XCTest

### Infrastructure

- **Containerization**: Docker 24+
- **Orchestration**: Kubernetes 1.28+
- **Infrastructure**: Terraform 1.6+
- **Package Management**: Helm 3.14+
- **CI/CD**: GitHub Actions
- **Database**: PostgreSQL 16+
- **Cache**: Redis 7+
- **Monitoring**: Prometheus + Grafana

## Current Package Versions

### Frontend Dependencies (package.json example)

```json
{
  "dependencies": {
    "next": "^15.0.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "typescript": "^5.3.0",
    "@next/font": "^15.0.0",
    "tailwindcss": "^3.4.0",
    "framer-motion": "^11.0.0",
    "axios": "^1.6.0",
    "@headlessui/react": "^1.7.0",
    "@heroicons/react": "^2.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.0",
    "@types/node": "^20.0.0",
    "eslint": "^8.57.0",
    "eslint-config-next": "^15.0.0",
    "prettier": "^3.2.0",
    "jest": "^29.7.0",
    "@testing-library/react": "^14.1.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0"
  }
}
```

### Backend Dependencies (go.mod example)

```go
module github.com/your-org/project-name

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/postgres v1.5.4
    github.com/redis/go-redis/v9 v9.3.0
    github.com/joho/godotenv v1.5.1
    github.com/google/uuid v1.5.0
    golang.org/x/crypto v0.17.0
    github.com/stretchr/testify v1.8.4
    google.golang.org/grpc v1.60.0
    google.golang.org/protobuf v1.31.0
)
```

### Android Dependencies (build.gradle example)

```kotlin
// App-level build.gradle (Kotlin DSL)
dependencies {
    implementation("androidx.core:core-ktx:1.12.0")
    implementation("androidx.lifecycle:lifecycle-runtime-ktx:2.7.0")
    implementation("androidx.activity:activity-compose:1.8.2")
    
    // Compose BOM
    implementation(platform("androidx.compose:compose-bom:2023.10.01"))
    implementation("androidx.compose.ui:ui")
    implementation("androidx.compose.ui:ui-tooling-preview")
    implementation("androidx.compose.material3:material3")
    
    // Navigation
    implementation("androidx.navigation:navigation-compose:2.7.6")
    
    // ViewModel
    implementation("androidx.lifecycle:lifecycle-viewmodel-compose:2.7.0")
    
    // Hilt
    implementation("com.google.dagger:hilt-android:2.48.1")
    kapt("com.google.dagger:hilt-compiler:2.48.1")
    
    // Networking
    implementation("com.squareup.retrofit2:retrofit:2.9.0")
    implementation("com.squareup.okhttp3:okhttp:4.12.0")
    implementation("com.squareup.retrofit2:converter-gson:2.9.0")
    
    // Room
    implementation("androidx.room:room-runtime:2.6.1")
    implementation("androidx.room:room-ktx:2.6.1")
    kapt("androidx.room:room-compiler:2.6.1")
    
    // Testing
    testImplementation("junit:junit:4.13.2")
    testImplementation("org.mockito:mockito-core:5.8.0")
    androidTestImplementation("androidx.test.ext:junit:1.1.5")
    androidTestImplementation("androidx.test.espresso:espresso-core:3.5.1")
    androidTestImplementation("androidx.compose.ui:ui-test-junit4")
}
```

### iOS Dependencies (Package.swift example)

```swift
// Package.swift
let package = Package(
    name: "ProjectName",
    platforms: [
        .iOS(.v15)
    ],
    dependencies: [
        .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.8.0"),
        .package(url: "https://github.com/Swinject/Swinject.git", from: "2.8.0"),
        .package(url: "https://github.com/realm/SwiftLint.git", from: "0.54.0"),
        .package(url: "https://github.com/kishikawakatsumi/KeychainAccess.git", from: "4.2.0")
    ]
)
```

## Key Configuration Files

### Root Level Files

#### Makefile (Root)

```makefile
# Project Configuration
PROJECT_NAME := project-name
VERSION := $(shell cat VERSION || echo "1.0.0")
REGISTRY := your-registry.com
IMAGE_TAG := $(REGISTRY)/$(PROJECT_NAME)

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

.PHONY: help setup clean build test lint deploy dev logs \
 build-app build-home build-admin build-server build-mobile \
 test-app test-home test-admin test-server test-mobile test-integration test-e2e \
 docker-build docker-push deploy-staging deploy-production \
 deps-update security-scan docs release dev-bg

## help: Show this help message
help:
 @echo "$(CYAN)Available commands:$(NC)"
 @sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## setup: Set up development environment
setup:
 @echo "$(CYAN)Setting up development environment...$(NC)"
 @if [ ! -f .env ]; then cp .env.example .env; fi
 @chmod +x Scripts/*.sh
 @./Scripts/setup.sh
 @echo "$(GREEN)Development environment setup complete$(NC)"

## dev: Start development environment
dev:
 @echo "$(CYAN)Starting development environment...$(NC)"
 @docker-compose up --build

## dev-bg: Start development environment in background
dev-bg:
 @echo "$(CYAN)Starting development environment in background...$(NC)"
 @docker-compose up -d --build

## build: Build all components
build: build-app build-home build-admin build-server build-mobile
 @echo "$(GREEN)All components built successfully$(NC)"

## build-app: Build main frontend application
build-app:
 @echo "$(CYAN)Building main frontend application...$(NC)"
 @cd App && make build

## build-home: Build landing page/marketing site
build-home:
 @echo "$(CYAN)Building landing page...$(NC)"
 @cd Home && make build

## build-admin: Build admin dashboard
build-admin:
 @echo "$(CYAN)Building admin dashboard...$(NC)"
 @cd Admin && make build

## build-server: Build backend server
build-server:
 @echo "$(CYAN)Building backend server...$(NC)"
 @cd CommonServer && make build

## build-mobile: Build mobile applications
build-mobile:
 @echo "$(CYAN)Building mobile applications...$(NC)"
 @if [ -d "Mobile/android" ]; then \
  cd Mobile/android && ./gradlew assembleRelease; \
 else \
  echo "$(YELLOW)Android project not found, skipping...$(NC)"; \
 fi
 @if [ -d "Mobile/ios" ]; then \
  cd Mobile/ios && xcodebuild -scheme ProjectName -configuration Release; \
 else \
  echo "$(YELLOW)iOS project not found, skipping...$(NC)"; \
 fi

## test: Run all tests
test: test-app test-home test-admin test-server test-mobile test-integration
 @echo "$(GREEN)All tests completed$(NC)"

## test-app: Run main frontend tests
test-app:
 @echo "$(CYAN)Running main frontend tests...$(NC)"
 @cd App && make test

## test-home: Run landing page tests
test-home:
 @echo "$(CYAN)Running landing page tests...$(NC)"
 @cd Home && make test

## test-admin: Run admin dashboard tests
test-admin:
 @echo "$(CYAN)Running admin dashboard tests...$(NC)"
 @cd Admin && make test

## test-server: Run backend tests
test-server:
 @echo "$(CYAN)Running backend tests...$(NC)"
 @cd CommonServer && make test

## test-mobile: Run mobile tests
test-mobile:
 @echo "$(CYAN)Running mobile tests...$(NC)"
 @if [ -d "Mobile/android" ]; then \
  cd Mobile/android && ./gradlew test; \
 fi
 @if [ -d "Mobile/ios" ]; then \
  cd Mobile/ios && xcodebuild test -scheme ProjectName -destination 'platform=iOS Simulator,name=iPhone 15'; \
 fi

## test-integration: Run integration tests
test-integration:
 @echo "$(CYAN)Running integration tests...$(NC)"
 @if [ -d "Tests" ]; then \
  cd Tests && make test; \
 else \
  echo "$(YELLOW)Integration tests not found, skipping...$(NC)"; \
 fi

## test-e2e: Run end-to-end tests
test-e2e:
 @echo "$(CYAN)Running end-to-end tests...$(NC)"
 @if [ -d "Tests/e2e" ]; then \
  cd Tests/e2e && npm test; \
 else \
  echo "$(YELLOW)E2E tests not found, skipping...$(NC)"; \
 fi

## lint: Run linting on all components
lint:
 @echo "$(CYAN)Running linting...$(NC)"
 @./Scripts/lint.sh

## clean: Clean build artifacts
clean:
 @echo "$(CYAN)Cleaning build artifacts...$(NC)"
 @./Scripts/clean.sh

## docker-build: Build Docker images
docker-build:
 @echo "$(CYAN)Building Docker images...$(NC)"
 @docker build -t $(IMAGE_TAG)-app:$(VERSION) ./App
 @docker build -t $(IMAGE_TAG)-home:$(VERSION) ./Home
 @docker build -t $(IMAGE_TAG)-admin:$(VERSION) ./Admin
 @docker build -t $(IMAGE_TAG)-server:$(VERSION) ./CommonServer

## docker-push: Push Docker images to registry
docker-push: docker-build
 @echo "$(CYAN)Pushing Docker images...$(NC)"
 @docker push $(IMAGE_TAG)-app:$(VERSION)
 @docker push $(IMAGE_TAG)-home:$(VERSION)
 @docker push $(IMAGE_TAG)-admin:$(VERSION)
 @docker push $(IMAGE_TAG)-server:$(VERSION)

## deploy-staging: Deploy to staging environment
deploy-staging:
 @echo "$(CYAN)Deploying to staging...$(NC)"
 @cd Deploy && terraform workspace select staging || terraform workspace new staging
 @cd Deploy && terraform apply -var="environment=staging" -var="image_tag=$(VERSION)"

## deploy-production: Deploy to production environment
deploy-production:
 @echo "$(YELLOW)Deploying to production...$(NC)"
 @read -p "Are you sure you want to deploy to production? [y/N] " confirm && \
 if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
  cd Deploy && terraform workspace select production || terraform workspace new production; \
  cd Deploy && terraform apply -var="environment=production" -var="image_tag=$(VERSION)"; \
 else \
  echo "$(RED)Deployment cancelled$(NC)"; \
 fi

## logs: View development logs
logs:
 @docker-compose logs -f

## deps-update: Update dependencies
deps-update:
 @echo "$(CYAN)Updating dependencies...$(NC)"
 @cd App && npm update
 @cd Home && npm update
 @cd Admin && npm update
 @cd CommonServer && go mod tidy && go get -u ./...
 @if [ -d "Mobile/android" ]; then cd Mobile/android && ./gradlew dependencies --refresh-dependencies; fi
 @if [ -d "Mobile/ios" ]; then cd Mobile/ios && pod update; fi

## security-scan: Run security scans
security-scan:
 @echo "$(CYAN)Running security scans...$(NC)"
 @cd App && npm audit
 @cd Home && npm audit
 @cd Admin && npm audit
 @cd CommonServer && gosec ./... || echo "$(YELLOW)gosec not installed, skipping Go security scan$(NC)"
 @docker run --rm -v $(PWD):/workspace securecodewarrior/docker-security-scanner || echo "$(YELLOW)Docker security scanner failed$(NC)"

## docs: Generate documentation
docs:
 @echo "$(CYAN)Generating documentation...$(NC)"
 @cd CommonServer && godoc -http=:6060 &
 @cd App && npm run docs || echo "$(YELLOW)No docs script in App$(NC)"
 @cd Home && npm run docs || echo "$(YELLOW)No docs script in Home$(NC)"
 @cd Admin && npm run docs || echo "$(YELLOW)No docs script in Admin$(NC)"

## release: Create a new release
release:
 @echo "$(CYAN)Creating release $(VERSION)...$(NC)"
 @git tag -a v$(VERSION) -m "Release v$(VERSION)"
 @git push origin v$(VERSION)
 @echo "$(GREEN)Release v$(VERSION) created$(NC)"
```

#### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: 
      context: ./App
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    volumes:
      - ./App:/usr/app
      - /usr/app/node_modules
      - /usr/app/.next
    depends_on:
      - server
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  home:
    build: 
      context: ./Home
      dockerfile: Dockerfile
    ports:
      - "3001:3000"
    environment:
      - NODE_ENV=development
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    volumes:
      - ./Home:/usr/app
      - /usr/app/node_modules
      - /usr/app/.next
    depends_on:
      - server
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  admin:
    build: 
      context: ./Admin
      dockerfile: Dockerfile
    ports:
      - "3002:3000"
    environment:
      - NODE_ENV=development
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    volumes:
      - ./Admin:/usr/app
      - /usr/app/node_modules
      - /usr/app/.next
    depends_on:
      - server
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  server:
    build: 
      context: ./CommonServer
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - GO_ENV=development
      - DB_HOST=database
      - DB_PORT=5432
      - DB_NAME=projectname
      - DB_USER=admin
      - DB_PASSWORD=password
      - REDIS_URL=redis:6379
      - JWT_SECRET=your-jwt-secret-key
    volumes:
      - ./CommonServer:/app
    depends_on:
      database:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  database:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=projectname
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=password
      - POSTGRES_INITDB_ARGS=--encoding=UTF-8 --lc-collate=en_US.UTF-8 --lc-ctype=en_US.UTF-8
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./CommonServer/database/init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U admin -d projectname"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

  # API Documentation
  api-docs:
    image: swaggerapi/swagger-ui
    ports:
      - "8081:8080"
    environment:
      - SWAGGER_JSON=/app/swagger.yaml
      - API_URL=http://localhost:8080
    volumes:
      - ./Mobile/shared/api:/app
      - ./Docs/Backend:/usr/share/nginx/html/docs
    depends_on:
      - server

  # Testing services
  test-db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=projectname_test
      - POSTGRES_USER=test
      - POSTGRES_PASSWORD=test
    ports:
      - "5433:5432"
    volumes:
      - test_db_data:/var/lib/postgresql/data
    profiles: ["test"]

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
  test_db_data:
    driver: local

networks:
  default:
    name: projectname-network
    driver: bridge
```

### Container Configurations

#### Frontend Dockerfile Example

```dockerfile
# Use latest Node.js LTS
FROM node:20-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on package manager
COPY package.json package-lock.json* ./
RUN npm ci --only=production

# Rebuild source code when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build

# Production image
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000
ENV PORT 3000

CMD ["node", "server.js"]
```

#### Backend Dockerfile Example

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

CMD ["./main"]
```

#### Essential Configuration Files

##### .gitignore Template

```gitignore
# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Production builds
build/
dist/
out/
.next/
.nuxt/

# Environment variables
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Runtime data
pids/
*.pid
*.seed
*.pid.lock

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
vendor/

# Mobile
*.ipa
*.apk
*.aab
build/
.gradle/
local.properties

# Testing
coverage/
.nyc_output/
junit.xml

# Docker
.dockerignore

# Terraform
*.tfstate
*.tfstate.*
.terraform/
.terraform.lock.hcl
```

##### .editorconfig Template

```ini
# EditorConfig is awesome: https://EditorConfig.org

root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
indent_style = space
indent_size = 2

[*.{js,jsx,ts,tsx,json,css,scss,md,yml,yaml}]
indent_size = 2

[*.{go,py}]
indent_size = 4

[*.{java,kt,swift}]
indent_size = 4

[Makefile]
indent_style = tab
```

##### SECURITY.md Template

```markdown
# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

Please report (suspected) security vulnerabilities to **[security@project-name.com](mailto:security@project-name.com)**. You will receive a response from us within 48 hours. If the issue is confirmed, we will release a patch as soon as possible depending on complexity but historically within a few days.

## Security Measures

### Code Security
- All dependencies are regularly updated
- Security scanning with npm audit, gosec, and other tools
- Input validation and sanitization
- Parameterized queries for database operations
- JWT tokens for authentication
- HTTPS/TLS encryption in production

### Infrastructure Security
- Container security scanning
- Kubernetes security policies
- Network segmentation
- Regular security updates
- Monitoring and alerting

### Development Security
- Branch protection rules
- Required pull request reviews
- Automated security checks in CI/CD
- Dependency vulnerability scanning
- Code quality gates
```

#### GitHub Actions Workflows

##### CI Workflow (.github/workflows/ci.yml)

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test-frontend:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [20.x]
        app: [App, Home, Admin]
    
    steps:
    - uses: actions/checkout@v4
    - name: Setup Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
        cache: 'npm'
        cache-dependency-path: ${{ matrix.app }}/package-lock.json
    
    - name: Install dependencies
      run: cd ${{ matrix.app }} && npm ci
    
    - name: Run tests
      run: cd ${{ matrix.app }} && npm test
    
    - name: Run lint
      run: cd ${{ matrix.app }} && npm run lint
    
    - name: Build
      run: cd ${{ matrix.app }} && npm run build

  test-backend:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
    - uses: actions/checkout@v4
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    
    - name: Build
      run: cd CommonServer && go build -v ./...
    
    - name: Test
      run: cd CommonServer && go test -v ./...
      env:
        DATABASE_URL: postgres://postgres:postgres@localhost:5432/test?sslmode=disable

  test-mobile:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Android
      if: hashFiles('Mobile/android/**') != ''
      uses: android-actions/setup-android@v2
    
    - name: Build Android
      if: hashFiles('Mobile/android/**') != ''
      run: cd Mobile/android && ./gradlew build
    
    - name: Test Android
      if: hashFiles('Mobile/android/**') != ''
      run: cd Mobile/android && ./gradlew test
```

##### Security Workflow (.github/workflows/security.yml)

```yaml
name: Security

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 2 * * 1' # Weekly on Monday

permissions:
  contents: read
  security-events: write

jobs:
  dependency-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Run npm audit (App)
      run: cd App && npm audit --audit-level high
      continue-on-error: true
    
    - name: Run npm audit (Home)
      run: cd Home && npm audit --audit-level high
      continue-on-error: true
    
    - name: Run npm audit (Admin)
      run: cd Admin && npm audit --audit-level high
      continue-on-error: true

  code-scanning:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Initialize CodeQL
      uses: github/codeql-action/init@v2
      with:
        languages: go, javascript
    
    - name: Autobuild
      uses: github/codeql-action/autobuild@v2
    
    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v2

  container-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Build images
      run: |
        docker build -t app:latest ./App
        docker build -t server:latest ./CommonServer
    
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: 'app:latest'
        format: 'sarif'
        output: 'trivy-results.sarif'
    
    - name: Upload Trivy scan results
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
```

##### Dependabot Configuration (.github/dependabot.yml)

```yaml
version: 2
updates:
  # Frontend dependencies
  - package-ecosystem: "npm"
    directory: "/App"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "team-frontend"
    assignees:
      - "lead-developer"

  - package-ecosystem: "npm"
    directory: "/Home"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5

  - package-ecosystem: "npm"
    directory: "/Admin"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5

  # Backend dependencies
  - package-ecosystem: "gomod"
    directory: "/CommonServer"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "team-backend"

  # Mobile dependencies
  - package-ecosystem: "gradle"
    directory: "/Mobile/android"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 3

  # Docker dependencies
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 3

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 3
```

#### Frontend Deployment Configurations

##### vercel.json

```json
{
  "version": 2,
  "builds": [
    {
      "src": "package.json",
      "use": "@vercel/static-build",
      "config": {
        "distDir": "out"
      }
    }
  ],
  "routes": [
    {
      "src": "/api/(.*)",
      "dest": "https://your-api-domain.com/api/$1"
    },
    {
      "src": "/(.*)",
      "dest": "/$1"
    }
  ],
  "env": {
    "NEXT_PUBLIC_API_URL": "https://your-api-domain.com"
  },
  "build": {
    "env": {
      "NEXT_PUBLIC_API_URL": "https://your-api-domain.com"
    }
  }
}
```

##### DigitalOcean App Platform (app.yaml)

```yaml
name: project-name-frontend
services:
- name: frontend
  source_dir: /
  github:
    repo: your-org/project-name
    branch: main
  run_command: npm start
  build_command: npm run build
  environment_slug: node-js
  instance_count: 1
  instance_size_slug: basic-xxs
  routes:
  - path: /
  envs:
  - key: NEXT_PUBLIC_API_URL
    value: https://your-api-domain.com
  - key: NODE_ENV
    value: production
```

## Essential Documentation Templates

### README.md Template

```markdown
# Project Name

Brief description of what the project does.

[![CI](https://github.com/your-org/project-name/workflows/CI/badge.svg)](https://github.com/your-org/project-name/actions)
[![Security](https://github.com/your-org/project-name/workflows/Security/badge.svg)](https://github.com/your-org/project-name/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features
- ✨ Feature 1
- 🚀 Feature 2  
- 🔒 Feature 3

## Architecture

This is a modern full-stack application with:
- **Frontend**: Next.js 15+ with static export (deployed on Vercel/DigitalOcean Apps)
- **Backend**: Go 1.22+ API server (deployed on Kubernetes/DigitalOcean)
- **Mobile**: Native Android (Kotlin 2.0+) and iOS (Swift 5.9+) applications
- **Database**: PostgreSQL 16+ with Redis 7+ caching
- **Infrastructure**: Docker 24+ containers with Kubernetes 1.28+ orchestration

## Quick Start

### Prerequisites

**Required:**
- Make
- Docker & Docker Compose
- Git

**For Backend Development:**
- Go 1.22+

**For Frontend Development:**
- Node.js 20+

**For Mobile Development:**

*Android:*
- Android Studio
- Android SDK (API level 24+)
- Java 17+

*iOS:*
- Xcode 15+
- iOS 15.0+
- CocoaPods

### Installation

1. **Clone and setup**
   ```bash
   git clone https://github.com/your-org/project-name.git
   cd project-name
   make setup
   ```

2. **Start development environment**

   ```bash
   make dev
   ```

3. **Access the applications**

- Main App: <http://localhost:3000>
- Landing Page: <http://localhost:3001>
- Admin Dashboard: <http://localhost:3002>
- Backend API: <http://localhost:8080>
- API Documentation: <http://localhost:8081>

### Development Commands

```bash
# Show all available commands
make help

# Development
make dev                 # Start development environment
make dev-bg             # Start in background

# Building
make build              # Build all components
make build-app          # Build main app
make build-home         # Build landing page
make build-admin        # Build admin dashboard
make build-server       # Build backend only
make build-mobile       # Build mobile apps

# Testing
make test               # Run all tests
make test-app          # Main app tests
make test-home         # Landing page tests
make test-admin        # Admin dashboard tests
make test-server       # Backend tests
make test-mobile       # Mobile tests
make test-e2e          # End-to-end tests

# Quality
make lint              # Run linting
make security-scan     # Security scanning

# Deployment
make docker-build      # Build Docker images
make deploy-staging    # Deploy to staging
make deploy-production # Deploy to production

# Utilities
make clean             # Clean build artifacts
make logs             # View logs
make deps-update      # Update dependencies
```

### Mobile Development

**Android Development:**

```bash
cd Mobile/android
./gradlew assembleDebug
```

**iOS Development:**

```bash
cd Mobile/ios
pod install
open ProjectName.xcworkspace
```

## Deployment

### Frontend Applications (Static)

- **Main App**: Deployed to Vercel/DO Apps from `./App`
- **Landing Page**: Deployed to Vercel/DO Apps from `./Home`
- **Admin Dashboard**: Deployed to Vercel/DO Apps from `./Admin`
- Each has individual `vercel.json` and `app.yaml` configurations

### Backend (Kubernetes)

```bash
# Deploy to staging
make deploy-staging

# Deploy to production (with confirmation)
make deploy-production
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=projectname
DB_USER=admin
DB_PASSWORD=your-password

# API
JWT_SECRET=your-jwt-secret
API_URL=http://localhost:8080

# Redis
REDIS_URL=redis://localhost:6379
```

## Testing

Run the full test suite:

```bash
make test
```

Individual test suites:

```bash
make test-app          # Main app unit tests
make test-home         # Landing page unit tests
make test-admin        # Admin dashboard unit tests
make test-server       # Backend unit tests  
make test-mobile       # Mobile unit tests
make test-integration  # Integration tests
make test-e2e         # End-to-end tests
```

## Documentation

- [API Documentation](./Docs/Backend/)
- [Main App Documentation](./App/README.md)
- [Landing Page Documentation](./Home/README.md)
- [Admin Dashboard Documentation](./Admin/README.md)
- [Mobile Documentation](./Mobile/README.md)
- [Deployment Guide](./Docs/deployment.md)
- [Self-Hosting Guide](./Docs/self_hosting.md)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run tests: `make test`
5. Run linting: `make lint`
6. Commit your changes: `git commit -am 'Add feature'`
7. Push to the branch: `git push origin feature-name`
8. Submit a pull request

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Support

- 📖 [Documentation](./Docs/)
- 🐛 [Issue Tracker](https://github.com/your-org/project-name/issues)
- 💬 [Discussions](https://github.com/your-org/project-name/discussions)

```md

### CONTRIBUTING.md Template

```markdown
# Contributing Guidelines

Thank you for your interest in contributing! This guide will help you get started.

## Quick Start

1. **Fork and clone**
   ```bash
   git clone https://github.com/your-username/project-name.git
   cd project-name
   ```

2. **Set up development environment**

   ```bash
   make setup
   ```

3. **Start development**

   ```bash
   make dev
   ```

4. **Run tests**

   ```bash
   make test
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

- Write code following our style guidelines
- Add tests for new features
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run specific test suites
make test-app          # Main app tests
make test-home         # Landing page tests
make test-admin        # Admin dashboard tests
make test-server       # Backend tests
make test-mobile       # Mobile tests
make test-integration  # Integration tests

# Run linting
make lint

# Run security scan
make security-scan
```

### 4. Commit Your Changes

```bash
git add .
git commit -m "feat: add your feature description"
```

We use [Conventional Commits](https://conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation updates
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

## Code Standards

### General Guidelines

- Write clear, readable code
- Add comments for complex logic
- Follow existing patterns and conventions
- Ensure all tests pass
- Maintain backwards compatibility when possible

### Frontend (Next.js/React)

- **Language**: TypeScript only
- **Styling**: Tailwind CSS
- **Components**: Functional components with hooks
- **State**: React Context for global state
- **Testing**: Jest + React Testing Library
- **Linting**: ESLint + Prettier

```bash
# Frontend commands (App, Home, Admin)
cd App                 # Main application
make test              # Run tests
make lint              # Lint code
make build             # Build for production

cd Home                # Landing page
make test              # Run tests
make lint              # Lint code
make build             # Build for production

cd Admin               # Admin dashboard
make test              # Run tests
make lint              # Lint code
make build             # Build for production
```

### Backend (Go)

- **Style**: Follow Go conventions
- **Architecture**: Clean architecture patterns
- **Database**: GORM for ORM
- **Testing**: Standard Go testing + testify
- **Linting**: golangci-lint

```bash
# Backend commands  
cd CommonServer
make test              # Run tests
make lint              # Lint code
make build             # Build binary
make dev               # Start dev server
```

### Mobile Development - Standards

#### Android (Kotlin)

- **Architecture**: MVVM + Android Architecture Components
- **Dependency Injection**: Dagger Hilt
- **Networking**: Retrofit + OkHttp 4.12+
- **Database**: Room
- **UI**: Jetpack Compose + Material Design 3
- **Testing**: JUnit 5 + Espresso
- **Linting**: ktlint
- **Build**: Gradle 8.5+ with Version Catalogs

```bash
cd Mobile/android
./gradlew test         # Unit tests
./gradlew connectedAndroidTest  # Instrumentation tests
./gradlew ktlintCheck  # Lint code
```

#### iOS (Swift)

- **Architecture**: MVVM + Combine
- **Dependency Injection**: Swinject
- **Networking**: URLSession + Alamofire 5.8+
- **Database**: SwiftData + Core Data
- **UI**: SwiftUI + UIKit
- **Testing**: XCTest
- **Linting**: SwiftLint

```bash
cd Mobile/ios
xcodebuild test -scheme ProjectName -destination 'platform=iOS Simulator,name=iPhone 15'
swiftlint              # Lint code
```

## Testing Requirements

### Coverage Standards

- Minimum 80% code coverage for all components
- 100% coverage for critical business logic
- All public APIs must have tests

### Test Types

1. **Unit Tests**: Test individual functions/components
2. **Integration Tests**: Test component interactions  
3. **End-to-End Tests**: Test complete user workflows
4. **Performance Tests**: Ensure acceptable performance

### Writing Tests

```bash
# Test naming convention
describe('ComponentName', () => {
  it('should do something when condition', () => {
    // Test implementation
  });
});
```

## Documentation - Standards

### Code Documentation

- Document all public APIs
- Add JSDoc/GoDoc comments
- Include usage examples
- Document complex algorithms

### User Documentation

- Update README.md for new features
- Add API documentation
- Create user guides for major features
- Update deployment guides

## Performance Guidelines

### Frontend

- Use Next.js static optimization
- Implement proper image optimization
- Use lazy loading for components
- Minimize bundle size

### Backend  

- Implement proper caching strategies
- Use database indexing effectively
- Profile and optimize critical paths
- Handle graceful degradation

### Mobile

- Optimize for different device capabilities
- Implement proper image caching
- Use lazy loading appropriately
- Monitor app performance metrics

## Security Guidelines

- Never commit secrets or credentials
- Sanitize all user inputs
- Use parameterized queries
- Implement proper authentication/authorization
- Follow OWASP security guidelines

## API Guidelines

### RESTful Conventions

- Use proper HTTP methods (GET, POST, PUT, DELETE)
- Use meaningful URLs
- Return appropriate status codes
- Include proper error messages

### Mobile API Integration

- Use shared API specifications in `/Mobile/shared/api/`
- Handle offline scenarios gracefully
- Implement proper error handling
- Support both platforms consistently

## Design System

### Consistency

- Follow shared design system in `/Mobile/shared/design/`
- Use platform-appropriate UI components
- Maintain consistency across platforms
- Support both light and dark modes

### Accessibility

- Implement proper ARIA labels
- Ensure keyboard navigation works
- Test with screen readers
- Maintain good color contrast ratios

## Pull Request Process

### Before Submitting

- [ ] All tests pass locally
- [ ] Code is properly formatted
- [ ] Documentation is updated
- [ ] Commit messages follow convention
- [ ] No merge conflicts

### PR Requirements

- Clear description of changes
- Link to related issues
- Screenshots for UI changes
- Performance impact assessment
- Breaking changes documented

### Review Process

1. Automated checks must pass
2. At least one maintainer review required
3. All feedback addressed
4. Final approval from code owner

## Getting Help

- 📖 [Documentation](./Docs/)
- 💬 [Discussions](https://github.com/your-org/project-name/discussions)
- 🐛 [Issues](https://github.com/your-org/project-name/issues)
- 📧 Email: <maintainers@project-name.com>

## Recognition

Contributors are recognized in:

- [CONTRIBUTORS.md](./CONTRIBUTORS.md)
- Release notes
- Annual contributor highlights

Thank you for contributing! 🎉

```md

## Benefits of This Structure

### **Production-Ready Architecture**
1. **Make-Based Build System**: Standardized, cross-platform build automation
2. **Container-First Design**: Docker containers with Kubernetes deployment
3. **Static Frontend Deployment**: Optimized for Vercel and DigitalOcean Apps
4. **Microservices Architecture**: Scalable backend with proper service separation
5. **Native Mobile Support**: Platform-specific optimizations for Android/iOS

### **Developer Experience**
6. **Single Command Setup**: `make setup` gets everything running
7. **Comprehensive Testing**: Unit, integration, E2E, and performance tests
8. **Development Containers**: VS Code dev containers for consistent environments
9. **Hot Reloading**: Fast development cycles across all platforms
10. **Automated Quality**: Linting, security scanning, and dependency updates

### **Open Source Excellence**
11. **Complete Governance**: CONTRIBUTING.md, SECURITY.md, CODEOWNERS
12. **Professional Documentation**: API docs, user guides, deployment guides
13. **CI/CD Ready**: GitHub Actions for testing, security, and deployment
14. **Dependency Management**: Automated updates with Dependabot
15. **Security First**: Built-in security scanning and best practices

### **Deployment Optimization**
16. **Multi-Environment**: Staging and production with Terraform workspaces
17. **Platform-Specific**: Vercel for frontend, Kubernetes for backend
18. **Monitoring Ready**: Health checks, logging, and observability built-in
19. **Scalable Infrastructure**: HPA, load balancing, and Redis caching
20. **Mobile Distribution**: Ready for App Store and Play Store deployment

## Implementation Checklist

### Initial Setup
- [ ] Clone template and rename project directories
- [ ] Update all `project-name` references to your actual project name
- [ ] Configure `.env.example` with your environment variables
- [ ] Set up your container registry in `Makefile` (REGISTRY variable)
- [ ] Update `CODEOWNERS` with your team members
- [ ] Configure GitHub repository settings (branch protection, secrets)

### Security Configuration
- [ ] Generate strong JWT secrets for all environments
- [ ] Configure database credentials securely
- [ ] Set up vulnerability scanning in CI/CD
- [ ] Enable Dependabot security updates
- [ ] Configure branch protection rules
- [ ] Set up secret scanning

### Deployment Setup
- [ ] Configure Vercel/DigitalOcean Apps for frontend applications
- [ ] Set up Kubernetes cluster for backend deployment
- [ ] Configure Terraform backend for state management
- [ ] Set up monitoring and logging
- [ ] Configure domain names and SSL certificates
- [ ] Set up database backups

### Development Environment
- [ ] Test `make setup` command works correctly
- [ ] Verify all services start with `make dev`
- [ ] Confirm all tests pass with `make test`
- [ ] Validate linting with `make lint`
- [ ] Test Docker builds with `make docker-build`

### Documentation
- [ ] Customize README.md with project-specific information
- [ ] Update CONTRIBUTING.md with your contribution guidelines
- [ ] Create API documentation
- [ ] Add architecture diagrams
- [ ] Document deployment procedures

## Best Practices Compliance

### ✅ Open Source Standards
- **LICENSE**: MIT license included
- **CONTRIBUTING.md**: Comprehensive contribution guidelines
- **SECURITY.md**: Security policy and vulnerability reporting
- **CODEOWNERS**: Code ownership and review assignments
- **Issue Templates**: Bug reports, feature requests, questions
- **Pull Request Template**: Standardized PR process

### ✅ Development Excellence
- **Make-based Build System**: Standardized cross-platform automation
- **Comprehensive Testing**: Unit, integration, E2E, and performance tests
- **Security Scanning**: Automated vulnerability detection
- **Code Quality**: Linting, formatting, and quality gates
- **Dependency Management**: Automated updates with security review

### ✅ Production Readiness
- **Container Security**: Multi-stage builds, non-root users, vulnerability scanning
- **Kubernetes Deployment**: HPA, health checks, proper resource limits
- **Static Frontend**: Optimized for CDN deployment
- **Database Management**: Migrations, backups, connection pooling
- **Monitoring**: Health checks, logging, metrics collection

### ✅ Modern Technology Stack
- **Latest Versions**: Go 1.22+, Node.js 20+, Next.js 15+
- **Mobile Native**: Kotlin 2.0+, Swift 5.9+, latest frameworks
- **Infrastructure**: Docker 24+, Kubernetes 1.28+, Terraform 1.6+
- **Security**: Latest security practices and tools

## Customization Notes

- Remove unused directories based on your project type (mobile, admin, etc.)
- Adjust technology stack as needed while maintaining version consistency
- Modify configuration files for your specific requirements
- Add project-specific documentation sections
- Customize CI/CD workflows for your deployment strategy
- Update package dependencies to latest stable versions
- Configure monitoring and alerting for your infrastructure

This template provides a **production-ready foundation** for open source projects while maintaining the flexibility to adapt to different use cases and technologies. It follows industry best practices and includes everything needed for professional software development and deployment.

## Summary

This design specification delivers a comprehensive, enterprise-grade foundation for open source projects that includes:

- **🏗️ Complete Architecture**: Multi-service monorepo with clear separation of concerns
- **🔧 Modern Tooling**: Make-based automation, latest technology versions, comprehensive testing
- **🚀 Production Ready**: Container security, Kubernetes deployment, monitoring, and observability
- **🔒 Security First**: Automated vulnerability scanning, security policies, and best practices
- **📚 Professional Standards**: Complete open source governance, documentation, and contribution guidelines
- **⚡ Developer Experience**: One-command setup, comprehensive automation, and clear workflows

Whether you're starting a new project or modernizing an existing one, this specification provides the blueprint for building maintainable, scalable, and secure open source software that meets enterprise standards while remaining accessible to the community.
