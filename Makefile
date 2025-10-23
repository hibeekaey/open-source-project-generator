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

.PHONY: help build test clean run install dev fmt vet lint security-scan \
        validate audit dist package release docker-build docker-test docker-push \
        docker-login docker-info ci check benchmark version validate-setup \
        check-versions update-versions

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build
build: ## Build the generator binary
	@echo "Building generator..."
	@echo "  Version: $(VERSION)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@mkdir -p bin
	@go build -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -s -w" -trimpath -o bin/generator ./cmd/generator
	@echo "✓ Build completed: bin/generator"

# Testing
test: ## Run tests with coverage (use TEST_FLAGS for options)
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out $(TEST_FLAGS) ./...
	@go tool cover -html=coverage.out -o coverage.html 2>/dev/null || true
	@echo "✓ Tests completed. Coverage report: coverage.html"

# Clean build artifacts
clean: ## Clean all build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ output/ dist/ packages/ test-reports/ benchmark_results/
	@rm -f coverage.out coverage.html performance.test results.sarif checksums.txt gosec-report.txt
	@rm -f *.tar.gz *.zip *.deb *.rpm *.pkg.tar.zst
	@echo "✓ Clean completed"

# Run
run: build ## Build and run the generator
	@./bin/generator

# Dependencies
install: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

# Development
dev: ## Run in development mode
	@echo "Starting development mode..."
	@go run ./cmd/generator

# Code Quality
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (installs if needed)
	@echo "Running linter..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest; \
	fi
	@golangci-lint run

# Security
security-scan: ## Run all security scanners (installs tools if needed)
	@echo "Running security scans..."
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi
	@gosec -no-fail -fmt=sarif -out=results.sarif ./... 2>/dev/null || true
	@echo "  ✓ gosec completed (results.sarif)"
	@govulncheck ./...
	@staticcheck ./...
	@echo "Security scans completed. SARIF report: results.sarif"

# Version Management
check-versions: ## Check for latest versions of dependencies
	@./scripts/check-latest-versions.sh

update-versions: ## Update versions in configs/versions.yaml
	@./scripts/update-versions.sh --auto-update

# Project Validation
validate: build ## Validate a project (Usage: make validate PROJECT=./path)
	@if [ -z "$(PROJECT)" ]; then \
		echo "Usage: make validate PROJECT=./path/to/project"; \
		exit 1; \
	fi
	@./bin/generator validate $(PROJECT)

audit: build ## Audit a project (Usage: make audit PROJECT=./path)
	@if [ -z "$(PROJECT)" ]; then \
		echo "Usage: make audit PROJECT=./path/to/project"; \
		exit 1; \
	fi
	@./bin/generator audit $(PROJECT) --security --quality

# Distribution & Packaging
dist: ## Build cross-platform binaries
	@echo "Building distribution binaries..."
	@./scripts/build.sh

package: dist ## Build distribution packages (DEB, RPM, Arch)
	@echo "Building distribution packages..."
	@./scripts/build-packages.sh all

release: test lint security-scan dist package ## Prepare full release
	@echo "✓ Release artifacts ready in dist/ and packages/"

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE_NAME):$(DOCKER_VERSION) \
		-t $(IMAGE_NAME):latest .

docker-test: docker-build ## Test Docker image
	@echo "Testing Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	@docker run --rm $(IMAGE_NAME):$(DOCKER_VERSION) version

docker-push: docker-build ## Push Docker image to registry
	@echo "Pushing Docker image: $(IMAGE_NAME):$(DOCKER_VERSION)"
	@docker push $(IMAGE_NAME):$(DOCKER_VERSION)
	@echo "Pushing Docker image: $(IMAGE_NAME):latest"
	@docker push $(IMAGE_NAME):latest

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

# CI/CD
ci: ## Run full CI pipeline
	@echo "Running CI pipeline..."
	@./scripts/ci-test.sh

check: fmt vet lint test ## Run all checks (pre-commit)
	@echo "✓ All checks passed"

# Benchmarks
benchmark: ## Run benchmarks (use BENCH_FLAGS for options)
	@echo "Running benchmarks..."
	@./scripts/run_performance_benchmarks.sh $(BENCH_FLAGS)

# Utilities
version: ## Show version information
	@./scripts/get-version.sh all

validate-setup: ## Validate project setup
	@echo "Validating project setup..."
	@./scripts/validate-setup.sh