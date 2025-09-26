# Open Source Project Generator Makefile

# Configuration
GITHUB_ACTOR ?= $(shell git config user.name || echo "unknown")
GITHUB_REPOSITORY_OWNER ?= $(shell git remote get-url origin | sed 's/.*[:/]\([^/]*\)\/[^/]*$$/\1/' || echo "cuesoftinc")
DOCKER_REGISTRY ?= ghcr.io
IMAGE_NAME ?= $(DOCKER_REGISTRY)/$(GITHUB_REPOSITORY_OWNER)/open-source-project-generator

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Clean version for Docker tags (replace / with -)
DOCKER_VERSION := $(shell echo $(VERSION) | sed 's/\//-/g')

.PHONY: help build test test-coverage clean run install dev lint fmt vet

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the generator binary
	@echo "Building generator..."
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@mkdir -p bin
	go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator ./cmd/generator

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

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
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.2; \
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
	@mkdir -p bin
	@echo "Building Linux AMD64..." && GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator-linux-amd64 ./cmd/generator &
	@echo "Building Darwin AMD64..." && GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator-darwin-amd64 ./cmd/generator &
	@echo "Building Darwin ARM64..." && GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator-darwin-arm64 ./cmd/generator &
	@echo "Building Windows AMD64..." && GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator-windows-amd64.exe ./cmd/generator &
	@wait
	@echo "All builds completed"

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

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE_NAME):$(DOCKER_VERSION) \
		-t $(IMAGE_NAME):latest .

docker-test: ## Test Docker image
	@echo "Testing Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	docker run --rm $(IMAGE_NAME):$(DOCKER_VERSION) version

docker-push: ## Push Docker image to registry
	@echo "Pushing Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	docker push $(IMAGE_NAME):$(DOCKER_VERSION)
	@echo "Pushing Docker image: $(IMAGE_NAME):latest"
	docker push $(IMAGE_NAME):latest

docker-login: ## Login to GitHub Container Registry
	@echo "Logging in to GitHub Container Registry as $(GITHUB_ACTOR)..."
	echo $(GITHUB_TOKEN) | docker login $(DOCKER_REGISTRY) -u $(GITHUB_ACTOR) --password-stdin

docker-info: ## Show Docker configuration
	@echo "Docker Configuration:"
	@echo "  Registry: $(DOCKER_REGISTRY)"
	@echo "  Repository Owner: $(GITHUB_REPOSITORY_OWNER)"
	@echo "  Image Name: $(IMAGE_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Docker Version: $(DOCKER_VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  GitHub Actor: $(GITHUB_ACTOR)"