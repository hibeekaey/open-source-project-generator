# Multi-stage Dockerfile for Open Source Template Generator

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION="dev"
ARG GIT_COMMIT="unknown"
ARG BUILD_TIME="unknown"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o generator ./cmd/generator

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    git \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 generator && \
    adduser -D -s /bin/bash -u 1000 -G generator generator

# Create necessary directories
RUN mkdir -p /workspace /home/generator/.config/generator /home/generator/.cache/generator && \
    chown -R generator:generator /workspace /home/generator

# Copy binary from builder stage
COPY --from=builder /app/generator /usr/local/bin/generator

# Copy templates and configuration
COPY --chown=generator:generator templates/ /usr/share/generator/templates/
COPY --chown=generator:generator docs/ /usr/share/generator/docs/

# Set permissions
RUN chmod +x /usr/local/bin/generator

# Switch to non-root user
USER generator

# Set working directory
WORKDIR /workspace

# Set environment variables
ENV GENERATOR_TEMPLATES_DIR=/usr/share/generator/templates
ENV GENERATOR_DOCS_DIR=/usr/share/generator/docs
ENV GENERATOR_CONFIG_DIR=/home/generator/.config/generator
ENV GENERATOR_CACHE_DIR=/home/generator/.cache/generator

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD generator version || exit 1

# Default command
ENTRYPOINT ["generator"]
CMD ["--help"]

# Labels for metadata
LABEL maintainer="Open Source Template Generator Team <team@example.com>"
LABEL description="Open Source Template Generator - Create production-ready project structures"
ARG VERSION="dev"
LABEL version="${VERSION}"
LABEL org.opencontainers.image.title="Open Source Template Generator"
LABEL org.opencontainers.image.description="Create production-ready project structures with modern best practices"
LABEL org.opencontainers.image.url="https://github.com/cuesoftinc/open-source-project-generator"
LABEL org.opencontainers.image.source="https://github.com/cuesoftinc/open-source-project-generator"
LABEL org.opencontainers.image.vendor="Open Source Template Generator Team"
LABEL org.opencontainers.image.licenses="MIT"