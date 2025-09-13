# Open Source Template Generator Makefile

.PHONY: help build test clean run install dev lint fmt vet

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the generator binary
	@echo "Building generator..."
	go build -o bin/generator ./cmd/generator

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Run CI-friendly tests
test-ci: ## Run tests suitable for CI/CD pipelines
	@echo "Running CI test suite..."
	@echo "ℹ️  Excluding resource-intensive and flaky tests"
	go test -tags=ci -timeout=5m ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run: build ## Build and run the generator
	./bin/generator

# Install dependencies
install: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Development mode (with auto-rebuild)
dev: ## Run in development mode
	@echo "Starting development mode..."
	go run ./cmd/generator

# Install golangci-lint
install-lint: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2; \
	else \
		echo "golangci-lint already installed"; \
	fi

# Lint the code
lint: ## Run golangci-lint
	@echo "Running linter..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Installing..."; \
		$(MAKE) install-lint; \
	fi
	golangci-lint run

# Format the code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Vet the code
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Setup development environment
setup: ## Setup development environment
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	@echo "Development environment ready!"

# Build for multiple platforms
build-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o bin/generator-linux-amd64 ./cmd/generator
	GOOS=darwin GOARCH=amd64 go build -o bin/generator-darwin-amd64 ./cmd/generator
	GOOS=darwin GOARCH=arm64 go build -o bin/generator-darwin-arm64 ./cmd/generator
	GOOS=windows GOARCH=amd64 go build -o bin/generator-windows-amd64.exe ./cmd/generator

# Distribution targets
dist: ## Build distribution packages
	@echo "Building distribution packages..."
	./scripts/build.sh

dist-clean: ## Clean distribution artifacts
	@echo "Cleaning distribution artifacts..."
	rm -rf dist/ packages/

# Package building
package-deb: ## Build DEB package
	@echo "Building DEB package..."
	./scripts/build-packages.sh deb

package-rpm: ## Build RPM package
	@echo "Building RPM package..."
	./scripts/build-packages.sh rpm

package-arch: ## Build Arch Linux package
	@echo "Building Arch Linux package..."
	./scripts/build-packages.sh arch

package-all: ## Build all packages
	@echo "Building all packages..."
	./scripts/build-packages.sh all

# Release preparation
release-prepare: dist package-all ## Prepare release artifacts
	@echo "Preparing release artifacts..."
	@echo "Distribution files created in dist/"
	@echo "Package files created in packages/"

# Installation testing
test-install: ## Test installation script
	@echo "Testing installation script..."
	bash -n scripts/install.sh
	@echo "Installation script syntax is valid"

# Audit targets
audit: ## Run comprehensive codebase audit
	@echo "Running comprehensive codebase audit..."
	./scripts/audit.sh

audit-structure: ## Run structural analysis only
	@echo "Running structural analysis..."
	./scripts/audit.sh --structure

audit-dependencies: ## Run dependency analysis only
	@echo "Running dependency analysis..."
	./scripts/audit.sh --dependencies

audit-quality: ## Run code quality analysis only
	@echo "Running code quality analysis..."
	./scripts/audit.sh --quality

audit-clean: ## Clean audit results
	@echo "Cleaning audit results..."
	rm -rf audit-results/

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t generator:latest .

docker-test: ## Test Docker image
	@echo "Testing Docker image..."
	docker run --rm generator:latest version