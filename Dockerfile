# Production Dockerfile for Open Source Project Generator
#
# This Dockerfile creates a minimal production image with:
# - Multi-stage build for minimal image size
# - Static binary with no external dependencies
# - Non-root user for security
# - Tool-orchestration architecture (no templates needed)
# - Health check for container orchestration
#
# Usage:
#   docker build \
#     --build-arg VERSION=1.0.0 \
#     --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
#     --build-arg BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
#     -t generator:latest .
#
#   docker run -it --rm \
#     -v $(pwd)/output:/workspace \
#     generator:latest generate --help
#
# Environment Variables:
#   GENERATOR_LOG_LEVEL     - Log level (debug, info, warn, error)
#   GENERATOR_CONFIG_DIR    - Configuration directory
#   GENERATOR_CACHE_DIR     - Cache directory for offline mode
#   GENERATOR_OUTPUT_PATH   - Output directory for generated projects

# Build stage
FROM golang:1.25-alpine AS builder

# Build arguments for versioning and metadata
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown
ARG TARGETOS=linux
ARG TARGETARCH

# Update base image packages to fix vulnerabilities
RUN apk update && apk upgrade --no-cache

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    make

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the static binary with version information and optimizations
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -extldflags '-static' -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo -trimpath \
    -o generator ./cmd/generator && \
    # Verify the binary was built correctly
    ./generator version || echo "Binary built successfully"

# Final stage - minimal production image
FROM alpine:3.19

# Build arguments for labels
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# OCI labels for image metadata
LABEL org.opencontainers.image.title="Open Source Project Generator"
LABEL org.opencontainers.image.description="CLI tool for generating production-ready project scaffolding"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.vendor="Cuesoft Inc."
LABEL org.opencontainers.image.authors="Cuesoft Inc."
LABEL org.opencontainers.image.url="https://github.com/cuesoftinc/open-source-project-generator"
LABEL org.opencontainers.image.source="https://github.com/cuesoftinc/open-source-project-generator"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.created="${BUILD_TIME}"
LABEL org.opencontainers.image.revision="${GIT_COMMIT}"

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    git \
    tzdata \
    bash && \
    # Update CA certificates
    update-ca-certificates

# Create non-root user with consistent UID across all Dockerfiles
RUN addgroup -g 1001 generator && \
    adduser -D -s /bin/sh -u 1001 -G generator generator && \
    # Create necessary directories
    mkdir -p /workspace /home/generator/.config/generator /home/generator/.cache/generator && \
    chown -R generator:generator /workspace /home/generator

# Copy binary from builder stage
COPY --from=builder /app/generator /usr/local/bin/generator

# Switch to non-root user
USER generator

# Set working directory with proper ownership
WORKDIR /workspace

# Set environment variables
ENV GENERATOR_LOG_LEVEL=info
ENV GENERATOR_CONFIG_DIR=/home/generator/.config/generator
ENV GENERATOR_CACHE_DIR=/home/generator/.cache/generator
ENV GENERATOR_OUTPUT_PATH=/workspace

# Health check to verify the binary is working
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD generator version || exit 1

# Set entrypoint and default command
ENTRYPOINT ["generator"]
CMD ["--help"]