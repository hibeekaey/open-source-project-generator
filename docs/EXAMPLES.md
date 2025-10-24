# Configuration Examples

Real-world configuration examples for common project types.

## Table of Contents

- [Full-Stack Web Application](#full-stack-web-application)
- [Frontend-Only Application](#frontend-only-application)
- [Backend-Only API](#backend-only-api)
- [Mobile Application](#mobile-application)
- [Microservice](#microservice)
- [Multi-Component Platform](#multi-component-platform)
- [CI/CD Integration](#cicd-integration)

---

## Full-Stack Web Application

A complete web application with Next.js frontend and Go backend.

### Full-Stack Configuration

```yaml
name: "ecommerce-platform"
description: "E-commerce platform with Next.js and Go"
output_dir: "./ecommerce-platform"
author: "Your Team"
license: "MIT"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true
      eslint: true

  - type: go-backend
    name: api-server
    enabled: true
    config:
      module: github.com/myorg/ecommerce-platform
      framework: gin
      port: 8080
      cors_enabled: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
    frontend: http://localhost:3000
  shared_environment:
    NODE_ENV: development
    API_URL: http://localhost:8080
    DATABASE_URL: postgres://localhost:5432/ecommerce
    REDIS_URL: redis://localhost:6379

options:
  use_external_tools: true
  create_backup: true
  verbose: false
```

### Full-Stack Generate

```bash
generator generate --config ecommerce-platform.yaml
```

### Full-Stack Generated Structure

```text
ecommerce-platform/
├── App/                    # Next.js frontend
│   ├── app/
│   ├── components/
│   ├── public/
│   ├── package.json
│   └── next.config.js
├── CommonServer/          # Go backend
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── go.mod
│   └── main.go
├── Deploy/
│   └── docker/
├── docker-compose.yml
├── Makefile
└── README.md
```

### Full-Stack Usage

```bash
cd ecommerce-platform

# Start with Docker Compose
docker-compose up

# Or start individually
cd App && npm run dev
cd CommonServer && go run cmd/main.go
```

---

## Frontend-Only Application

A standalone Next.js application without backend.

### Frontend-Only Configuration

```yaml
name: "marketing-website"
description: "Marketing website with Next.js"
output_dir: "./marketing-website"

components:
  - type: nextjs
    name: website
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true
      eslint: true
      src_dir: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: https://api.example.com
  shared_environment:
    NEXT_PUBLIC_API_URL: https://api.example.com
    NEXT_PUBLIC_SITE_NAME: "Marketing Website"

options:
  use_external_tools: true
  create_backup: true
```

### Frontend-Only Generate

```bash
generator generate --config marketing-website.yaml
```

### Frontend-Only Usage

```bash
cd marketing-website/App
npm install
npm run dev
```

---

## Backend-Only API

A standalone Go API server.

### Backend-Only Configuration

```yaml
name: "user-service"
description: "User management microservice"
output_dir: "./user-service"

components:
  - type: go-backend
    name: user-api
    enabled: true
    config:
      module: github.com/myorg/user-service
      framework: gin
      port: 8081
      cors_enabled: true
      auth_enabled: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    database: postgres://localhost:5432/users
    cache: redis://localhost:6379
  shared_environment:
    SERVICE_NAME: user-service
    SERVICE_PORT: "8081"
    DB_HOST: postgres
    DB_PORT: "5432"
    DB_NAME: users
    REDIS_HOST: redis
    REDIS_PORT: "6379"
    LOG_LEVEL: info

options:
  use_external_tools: true
  create_backup: true
```

### Backend-Only Generate

```bash
generator generate --config user-service.yaml
```

### Backend-Only Usage

```bash
cd user-service/CommonServer
go mod tidy
go run cmd/main.go
```

---

## Mobile Application

Cross-platform mobile app with backend API.

### Mobile Configuration

```yaml
name: "fitness-tracker"
description: "Fitness tracking mobile application"
output_dir: "./fitness-tracker"

components:
  - type: android
    name: android-app
    enabled: true
    config:
      package: com.example.fitnesstracker
      min_sdk: 24
      target_sdk: 34
      language: kotlin
      compose: true

  - type: ios
    name: ios-app
    enabled: true
    config:
      bundle_id: com.example.fitnesstracker
      deployment_target: "15.0"
      language: swift
      swiftui: true

  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/fitness-tracker
      framework: gin
      port: 8080
      cors_enabled: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    backend: http://localhost:8080
    backend_prod: https://api.fitnesstracker.com
  shared_environment:
    API_URL: http://localhost:8080
    API_TIMEOUT: "30"
    ENABLE_LOGGING: "true"

options:
  use_external_tools: true
  create_backup: true
```

### Mobile Generate

```bash
generator generate --config fitness-tracker.yaml
```

### Mobile Generated Structure

```text
fitness-tracker/
├── Mobile/
│   ├── android/           # Android Kotlin app
│   └── ios/               # iOS Swift app
├── CommonServer/          # Go backend
├── Deploy/
│   └── docker/
├── docker-compose.yml
└── README.md
```

### Mobile Usage

```bash
# Start backend
cd fitness-tracker/CommonServer
go run cmd/main.go

# Android development
cd fitness-tracker/Mobile/android
./gradlew assembleDebug

# iOS development (macOS only)
cd fitness-tracker/Mobile/ios
open FitnessTracker.xcodeproj
```

---

## Microservice

A microservice with Docker and Kubernetes support.

### Microservice Configuration

```yaml
name: "payment-service"
description: "Payment processing microservice"
output_dir: "./payment-service"

components:
  - type: go-backend
    name: payment-api
    enabled: true
    config:
      module: github.com/myorg/payment-service
      framework: gin
      port: 8082
      cors_enabled: true
      auth_enabled: true

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    database: postgres://localhost:5432/payments
    cache: redis://localhost:6379
    message_queue: amqp://localhost:5672
  shared_environment:
    SERVICE_NAME: payment-service
    SERVICE_PORT: "8082"
    DB_HOST: postgres
    DB_PORT: "5432"
    DB_NAME: payments
    REDIS_HOST: redis
    REDIS_PORT: "6379"
    RABBITMQ_HOST: rabbitmq
    RABBITMQ_PORT: "5672"
    LOG_LEVEL: info
    METRICS_ENABLED: "true"

options:
  use_external_tools: true
  create_backup: true
  verbose: false
```

### Microservice Generate

```bash
generator generate --config payment-service.yaml
```

### Microservice Deployment

```bash
cd payment-service

# Local development
docker-compose up

# Build Docker image
docker build -t payment-service:latest -f Deploy/docker/Dockerfile .

# Deploy to Kubernetes (if K8s manifests generated)
kubectl apply -f Deploy/k8s/
```

---

## Multi-Component Platform

A comprehensive platform with multiple services.

### Multi-Component Configuration

```yaml
name: "saas-platform"
description: "Multi-tenant SaaS platform"
output_dir: "./saas-platform"

components:
  # Frontend applications
  - type: nextjs
    name: customer-portal
    enabled: true
    config:
      typescript: true
      tailwind: true
      app_router: true

  # Backend services
  - type: go-backend
    name: auth-service
    enabled: true
    config:
      module: github.com/myorg/saas-platform/auth
      framework: gin
      port: 8081

  - type: go-backend
    name: billing-service
    enabled: true
    config:
      module: github.com/myorg/saas-platform/billing
      framework: gin
      port: 8082

  - type: go-backend
    name: notification-service
    enabled: true
    config:
      module: github.com/myorg/saas-platform/notification
      framework: gin
      port: 8083

integration:
  generate_docker_compose: true
  generate_scripts: true
  api_endpoints:
    auth: http://localhost:8081
    billing: http://localhost:8082
    notification: http://localhost:8083
    frontend: http://localhost:3000
  shared_environment:
    NODE_ENV: development
    DATABASE_URL: postgres://localhost:5432/saas
    REDIS_URL: redis://localhost:6379
    JWT_SECRET: your-secret-key
    STRIPE_API_KEY: your-stripe-key

options:
  use_external_tools: true
  create_backup: true
```

### Multi-Component Generate

```bash
generator generate --config saas-platform.yaml
```

---

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/generate-and-test.yml
name: Generate and Test

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  generate:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Install Generator
        run: |
          go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@latest
      
      - name: Check Tools
        run: generator check-tools
      
      - name: Generate Project
        run: |
          generator generate \
            --config .generator/project.yaml \
            --output ./generated \
            --verbose
      
      - name: Test Frontend
        run: |
          cd generated/App
          npm install
          npm run build
      
      - name: Test Backend
        run: |
          cd generated/CommonServer
          go mod tidy
          go test ./...
          go build ./cmd/...
      
      - name: Upload Artifacts
        uses: actions/upload-artifact@v3
        with:
          name: generated-project
          path: generated/
```

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - generate
  - test
  - build

variables:
  GENERATOR_VERSION: "latest"

generate:
  stage: generate
  image: golang:1.21
  before_script:
    - apt-get update && apt-get install -y nodejs npm
    - go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@${GENERATOR_VERSION}
  script:
    - generator check-tools
    - generator generate --config .generator/project.yaml --output ./generated
  artifacts:
    paths:
      - generated/
    expire_in: 1 hour

test-frontend:
  stage: test
  image: node:20
  dependencies:
    - generate
  script:
    - cd generated/App
    - npm install
    - npm run build
    - npm test

test-backend:
  stage: test
  image: golang:1.21
  dependencies:
    - generate
  script:
    - cd generated/CommonServer
    - go mod tidy
    - go test ./...
    - go build ./cmd/...

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  dependencies:
    - generate
  script:
    - cd generated
    - docker-compose build
```

### Jenkins Pipeline

```groovy
// Jenkinsfile
pipeline {
    agent any
    
    environment {
        GENERATOR_VERSION = 'latest'
    }
    
    stages {
        stage('Setup') {
            steps {
                sh '''
                    go install github.com/cuesoftinc/open-source-project-generator/cmd/generator@${GENERATOR_VERSION}
                '''
            }
        }
        
        stage('Check Tools') {
            steps {
                sh 'generator check-tools'
            }
        }
        
        stage('Generate') {
            steps {
                sh '''
                    generator generate \
                        --config .generator/project.yaml \
                        --output ./generated \
                        --verbose
                '''
            }
        }
        
        stage('Test') {
            parallel {
                stage('Test Frontend') {
                    steps {
                        dir('generated/App') {
                            sh 'npm install'
                            sh 'npm run build'
                            sh 'npm test'
                        }
                    }
                }
                
                stage('Test Backend') {
                    steps {
                        dir('generated/CommonServer') {
                            sh 'go mod tidy'
                            sh 'go test ./...'
                            sh 'go build ./cmd/...'
                        }
                    }
                }
            }
        }
        
        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'generated/**/*', fingerprint: true
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}
```

---

## Environment-Specific Configurations

### Development

```yaml
name: "myapp-dev"
description: "Development environment"
output_dir: "./myapp-dev"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true

  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/myapp
      framework: gin
      port: 8080

integration:
  generate_docker_compose: true
  shared_environment:
    NODE_ENV: development
    GO_ENV: development
    LOG_LEVEL: debug
    ENABLE_DEBUG: "true"

options:
  use_external_tools: true
  verbose: true
```

### Production

```yaml
name: "myapp-prod"
description: "Production environment"
output_dir: "./myapp-prod"

components:
  - type: nextjs
    name: web-app
    enabled: true
    config:
      typescript: true
      tailwind: true

  - type: go-backend
    name: api
    enabled: true
    config:
      module: github.com/myorg/myapp
      framework: gin
      port: 8080

integration:
  generate_docker_compose: true
  shared_environment:
    NODE_ENV: production
    GO_ENV: production
    LOG_LEVEL: info
    ENABLE_DEBUG: "false"

options:
  use_external_tools: true
  verbose: false
```

---

## Tips for Configuration

### 1. Start with Templates

```bash
# Generate a template
generator init-config --example fullstack

# Customize it
vim project.yaml

# Generate
generator generate --config project.yaml
```

### 2. Use Dry Run

```bash
# Preview before generating
generator generate --config project.yaml --dry-run
```

### 3. Version Control Configs

```bash
mkdir -p .generator
generator init-config .generator/dev.yaml
generator init-config .generator/prod.yaml
git add .generator/
```

### 4. Environment Variables

```bash
# Override specific values
export GENERATOR_VERBOSE=true
generator generate --config project.yaml
```

### 5. Incremental Development

```bash
# Start with minimal
generator init-config --minimal > minimal.yaml
generator generate --config minimal.yaml

# Add components as needed
# Edit minimal.yaml to add more components
generator generate --config minimal.yaml
```

---

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [Configuration Guide](CONFIGURATION.md) - Detailed configuration options
- [CLI Commands](CLI_COMMANDS.md) - Command reference
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
