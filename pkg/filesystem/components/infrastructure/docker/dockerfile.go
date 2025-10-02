package docker

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// DockerfileGenerator handles Dockerfile generation
type DockerfileGenerator struct{}

// NewDockerfileGenerator creates a new Dockerfile generator
func NewDockerfileGenerator() *DockerfileGenerator {
	return &DockerfileGenerator{}
}

// GenerateFrontendDockerfile generates Dockerfile for frontend
func (dg *DockerfileGenerator) GenerateFrontendDockerfile(config *models.ProjectConfig) string {
	nodeVersion := "18"
	if config.Versions != nil && config.Versions.Node != "" {
		nodeVersion = config.Versions.Node
	}

	return fmt.Sprintf(`# %s Frontend Dockerfile
FROM node:%s-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json package-lock.json* ./
RUN npm ci --only=production

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build the application
RUN npm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]`, config.Name, nodeVersion)
}

// GenerateBackendDockerfile generates Dockerfile for backend
func (dg *DockerfileGenerator) GenerateBackendDockerfile(config *models.ProjectConfig) string {
	goVersion := "1.22"
	if config.Versions != nil && config.Versions.Go != "" {
		goVersion = config.Versions.Go
	}

	return fmt.Sprintf(`# %s Backend Dockerfile
FROM golang:%s-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]`, config.Name, goVersion)
}

// GenerateDockerIgnore generates .dockerignore content
func (dg *DockerfileGenerator) GenerateDockerIgnore(config *models.ProjectConfig) string {
	return `# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
.next/
*.exe
*.dll
*.so
*.dylib

# Environment files
.env
.env.local
.env.*.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Test coverage
coverage/
*.out

# Git
.git/
.gitignore

# Docker
Dockerfile*
docker-compose*
.dockerignore

# Documentation
README.md
docs/

# CI/CD
.github/
.gitlab-ci.yml`
}