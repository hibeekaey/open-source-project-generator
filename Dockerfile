# Build stage
FROM golang:1.25-alpine AS builder

# Build arguments
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

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

# Build the binary with version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -extldflags '-static' -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o generator ./cmd/generator

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN adduser -D -s /bin/sh generator

# Copy binary from builder stage
COPY --from=builder /app/generator /usr/local/bin/generator

# Make binary executable
RUN chmod +x /usr/local/bin/generator

# Switch to non-root user
USER generator

# Set working directory
WORKDIR /workspace

# Set entrypoint
ENTRYPOINT ["generator"]
CMD ["--help"]