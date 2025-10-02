package components

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// InfrastructureGenerator handles infrastructure component file generation
type InfrastructureGenerator struct {
	fsOps FileSystemOperations
}

// NewInfrastructureGenerator creates a new infrastructure generator
func NewInfrastructureGenerator(fsOps FileSystemOperations) *InfrastructureGenerator {
	return &InfrastructureGenerator{
		fsOps: fsOps,
	}
}

// GenerateFiles creates infrastructure component files based on configuration
func (ig *InfrastructureGenerator) GenerateFiles(projectPath string, config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Generate Docker files if selected
	if config.Components.Infrastructure.Docker {
		if err := ig.generateDockerFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Docker files: %w", err)
		}
	}

	// Generate Kubernetes files if selected
	if config.Components.Infrastructure.Kubernetes {
		if err := ig.generateKubernetesFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Kubernetes files: %w", err)
		}
	}

	// Generate Terraform files if selected
	if config.Components.Infrastructure.Terraform {
		if err := ig.generateTerraformFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Terraform files: %w", err)
		}
	}

	return nil
}

// generateDockerFiles creates Docker configuration files
func (ig *InfrastructureGenerator) generateDockerFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate docker-compose.yml
	dockerComposeContent := ig.generateDockerCompose(config)
	dockerComposePath := filepath.Join(projectPath, "Deploy/docker/docker-compose.yml")
	if err := ig.fsOps.WriteFile(dockerComposePath, []byte(dockerComposeContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Generate docker-compose.dev.yml
	dockerComposeDevContent := ig.generateDockerComposeDev(config)
	dockerComposeDevPath := filepath.Join(projectPath, "Deploy/docker/docker-compose.dev.yml")
	if err := ig.fsOps.WriteFile(dockerComposeDevPath, []byte(dockerComposeDevContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.dev.yml: %w", err)
	}

	// Generate docker-compose.prod.yml
	dockerComposeProdContent := ig.generateDockerComposeProd(config)
	dockerComposeProdPath := filepath.Join(projectPath, "Deploy/docker/docker-compose.prod.yml")
	if err := ig.fsOps.WriteFile(dockerComposeProdPath, []byte(dockerComposeProdContent), 0644); err != nil {
		return fmt.Errorf("failed to create docker-compose.prod.yml: %w", err)
	}

	// Generate .dockerignore
	dockerIgnoreContent := ig.generateDockerIgnore(config)
	dockerIgnorePath := filepath.Join(projectPath, "Deploy/docker/.dockerignore")
	if err := ig.fsOps.WriteFile(dockerIgnorePath, []byte(dockerIgnoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .dockerignore: %w", err)
	}

	// Generate Dockerfile for frontend
	frontendDockerfileContent := ig.generateFrontendDockerfile(config)
	frontendDockerfilePath := filepath.Join(projectPath, "Deploy/docker/Dockerfile.frontend")
	if err := ig.fsOps.WriteFile(frontendDockerfilePath, []byte(frontendDockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile.frontend: %w", err)
	}

	// Generate Dockerfile for backend
	backendDockerfileContent := ig.generateBackendDockerfile(config)
	backendDockerfilePath := filepath.Join(projectPath, "Deploy/docker/Dockerfile.backend")
	if err := ig.fsOps.WriteFile(backendDockerfilePath, []byte(backendDockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile.backend: %w", err)
	}

	return nil
}

// generateKubernetesFiles creates Kubernetes configuration files
func (ig *InfrastructureGenerator) generateKubernetesFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate namespace.yaml
	namespaceContent := ig.generateKubernetesNamespace(config)
	namespacePath := filepath.Join(projectPath, "Deploy/kubernetes/base/namespace.yaml")
	if err := ig.fsOps.WriteFile(namespacePath, []byte(namespaceContent), 0644); err != nil {
		return fmt.Errorf("failed to create namespace.yaml: %w", err)
	}

	// Generate configmap.yaml
	configMapContent := ig.generateKubernetesConfigMap(config)
	configMapPath := filepath.Join(projectPath, "Deploy/kubernetes/base/configmap.yaml")
	if err := ig.fsOps.WriteFile(configMapPath, []byte(configMapContent), 0644); err != nil {
		return fmt.Errorf("failed to create configmap.yaml: %w", err)
	}

	// Generate secret.yaml
	secretContent := ig.generateKubernetesSecret(config)
	secretPath := filepath.Join(projectPath, "Deploy/kubernetes/base/secret.yaml")
	if err := ig.fsOps.WriteFile(secretPath, []byte(secretContent), 0644); err != nil {
		return fmt.Errorf("failed to create secret.yaml: %w", err)
	}

	// Generate backend deployment
	backendDeploymentContent := ig.generateKubernetesBackendDeployment(config)
	backendDeploymentPath := filepath.Join(projectPath, "Deploy/kubernetes/base/backend-deployment.yaml")
	if err := ig.fsOps.WriteFile(backendDeploymentPath, []byte(backendDeploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to create backend-deployment.yaml: %w", err)
	}

	// Generate frontend deployment
	frontendDeploymentContent := ig.generateKubernetesFrontendDeployment(config)
	frontendDeploymentPath := filepath.Join(projectPath, "Deploy/kubernetes/base/frontend-deployment.yaml")
	if err := ig.fsOps.WriteFile(frontendDeploymentPath, []byte(frontendDeploymentContent), 0644); err != nil {
		return fmt.Errorf("failed to create frontend-deployment.yaml: %w", err)
	}

	// Generate services
	servicesContent := ig.generateKubernetesServices(config)
	servicesPath := filepath.Join(projectPath, "Deploy/kubernetes/base/services.yaml")
	if err := ig.fsOps.WriteFile(servicesPath, []byte(servicesContent), 0644); err != nil {
		return fmt.Errorf("failed to create services.yaml: %w", err)
	}

	// Generate ingress
	ingressContent := ig.generateKubernetesIngress(config)
	ingressPath := filepath.Join(projectPath, "Deploy/kubernetes/base/ingress.yaml")
	if err := ig.fsOps.WriteFile(ingressPath, []byte(ingressContent), 0644); err != nil {
		return fmt.Errorf("failed to create ingress.yaml: %w", err)
	}

	return nil
}

// generateTerraformFiles creates Terraform configuration files
func (ig *InfrastructureGenerator) generateTerraformFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate main.tf
	mainTfContent := ig.generateTerraformMain(config)
	mainTfPath := filepath.Join(projectPath, "Deploy/terraform/main.tf")
	if err := ig.fsOps.WriteFile(mainTfPath, []byte(mainTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.tf: %w", err)
	}

	// Generate variables.tf
	variablesTfContent := ig.generateTerraformVariables(config)
	variablesTfPath := filepath.Join(projectPath, "Deploy/terraform/variables.tf")
	if err := ig.fsOps.WriteFile(variablesTfPath, []byte(variablesTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create variables.tf: %w", err)
	}

	// Generate outputs.tf
	outputsTfContent := ig.generateTerraformOutputs(config)
	outputsTfPath := filepath.Join(projectPath, "Deploy/terraform/outputs.tf")
	if err := ig.fsOps.WriteFile(outputsTfPath, []byte(outputsTfContent), 0644); err != nil {
		return fmt.Errorf("failed to create outputs.tf: %w", err)
	}

	// Generate terraform.tfvars.example
	tfVarsExampleContent := ig.generateTerraformTfVarsExample(config)
	tfVarsExamplePath := filepath.Join(projectPath, "Deploy/terraform/terraform.tfvars.example")
	if err := ig.fsOps.WriteFile(tfVarsExamplePath, []byte(tfVarsExampleContent), 0644); err != nil {
		return fmt.Errorf("failed to create terraform.tfvars.example: %w", err)
	}

	// Generate modules/vpc/main.tf
	vpcModuleContent := ig.generateTerraformVPCModule(config)
	vpcModulePath := filepath.Join(projectPath, "Deploy/terraform/modules/vpc/main.tf")
	if err := ig.fsOps.WriteFile(vpcModulePath, []byte(vpcModuleContent), 0644); err != nil {
		return fmt.Errorf("failed to create modules/vpc/main.tf: %w", err)
	}

	return nil
}

// generateDockerCompose generates docker-compose.yml content
func (ig *InfrastructureGenerator) generateDockerCompose(config *models.ProjectConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  # %s Backend Service
  backend:
    build:
      context: ../../CommonServer
      dockerfile: Dockerfile
    container_name: %s-backend
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=development
      - DATABASE_URL=postgres://postgres:password@postgres:5432/%s_db?sslmode=disable
    depends_on:
      - postgres
    networks:
      - %s-network
    volumes:
      - ../../CommonServer:/app
    restart: unless-stopped

  # %s Frontend Service
  frontend:
    build:
      context: ../../App
      dockerfile: ../Deploy/docker/Dockerfile.frontend
    container_name: %s-frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
    depends_on:
      - backend
    networks:
      - %s-network
    volumes:
      - ../../App:/app
      - /app/node_modules
    restart: unless-stopped

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: %s-postgres
    environment:
      - POSTGRES_DB=%s_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - %s-network
    restart: unless-stopped

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: %s-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - %s-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  %s-network:
    driver: bridge`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateDockerComposeDev generates docker-compose.dev.yml content
func (ig *InfrastructureGenerator) generateDockerComposeDev(config *models.ProjectConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  backend:
    build:
      context: ../../CommonServer
      dockerfile: Dockerfile.dev
    environment:
      - ENVIRONMENT=development
      - LOG_LEVEL=debug
      - HOT_RELOAD=true
    volumes:
      - ../../CommonServer:/app
      - /app/bin
    command: ["go", "run", "main.go"]

  frontend:
    build:
      context: ../../App
      dockerfile: ../Deploy/docker/Dockerfile.frontend.dev
    environment:
      - NODE_ENV=development
      - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
    volumes:
      - ../../App:/app
      - /app/node_modules
      - /app/.next
    command: ["npm", "run", "dev"]

  postgres:
    environment:
      - POSTGRES_DB=%s_dev_db
    volumes:
      - ./dev-init.sql:/docker-entrypoint-initdb.d/init.sql`, config.Name)
}

// generateDockerComposeProd generates docker-compose.prod.yml content
func (ig *InfrastructureGenerator) generateDockerComposeProd(config *models.ProjectConfig) string {
	return fmt.Sprintf(`version: '3.8'

services:
  backend:
    image: %s-backend:latest
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

  frontend:
    image: %s-frontend:latest
    environment:
      - NODE_ENV=production
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M

  postgres:
    environment:
      - POSTGRES_DB=%s_prod_db
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - backend`, config.Name, config.Name, config.Name)
}

// generateDockerIgnore generates .dockerignore content
func (ig *InfrastructureGenerator) generateDockerIgnore(config *models.ProjectConfig) string {
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

// generateFrontendDockerfile generates Dockerfile for frontend
func (ig *InfrastructureGenerator) generateFrontendDockerfile(config *models.ProjectConfig) string {
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

// generateBackendDockerfile generates Dockerfile for backend
func (ig *InfrastructureGenerator) generateBackendDockerfile(config *models.ProjectConfig) string {
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

// generateKubernetesNamespace generates namespace.yaml content
func (ig *InfrastructureGenerator) generateKubernetesNamespace(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    app: %s
    environment: production`, config.Name, config.Name)
}

// generateKubernetesConfigMap generates configmap.yaml content
func (ig *InfrastructureGenerator) generateKubernetesConfigMap(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: %s-config
  namespace: %s
data:
  ENVIRONMENT: "production"
  LOG_LEVEL: "info"
  PORT: "8080"
  CORS_ORIGINS: "https://%s.com"`, config.Name, config.Name, config.Name)
}

// generateKubernetesSecret generates secret.yaml content
func (ig *InfrastructureGenerator) generateKubernetesSecret(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: %s-secrets
  namespace: %s
type: Opaque
data:
  # Base64 encoded values
  # Use: echo -n "your-secret" | base64
  DATABASE_URL: ""
  JWT_SECRET: ""`, config.Name, config.Name)
}

// generateKubernetesBackendDeployment generates backend deployment
func (ig *InfrastructureGenerator) generateKubernetesBackendDeployment(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s-backend
  namespace: %s
  labels:
    app: %s-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: %s-backend
  template:
    metadata:
      labels:
        app: %s-backend
    spec:
      containers:
      - name: backend
        image: %s-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          valueFrom:
            configMapKeyRef:
              name: %s-config
              key: ENVIRONMENT
        - name: PORT
          valueFrom:
            configMapKeyRef:
              name: %s-config
              key: PORT
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: %s-secrets
              key: DATABASE_URL
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: %s-secrets
              key: JWT_SECRET
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateKubernetesFrontendDeployment generates frontend deployment
func (ig *InfrastructureGenerator) generateKubernetesFrontendDeployment(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s-frontend
  namespace: %s
  labels:
    app: %s-frontend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: %s-frontend
  template:
    metadata:
      labels:
        app: %s-frontend
    spec:
      containers:
      - name: frontend
        image: %s-frontend:latest
        ports:
        - containerPort: 3000
        env:
        - name: NODE_ENV
          value: "production"
        - name: NEXT_PUBLIC_API_URL
          value: "https://api.%s.com/api/v1"
        livenessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateKubernetesServices generates services.yaml content
func (ig *InfrastructureGenerator) generateKubernetesServices(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s-backend-service
  namespace: %s
spec:
  selector:
    app: %s-backend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
apiVersion: v1
kind: Service
metadata:
  name: %s-frontend-service
  namespace: %s
spec:
  selector:
    app: %s-frontend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3000
  type: ClusterIP`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateKubernetesIngress generates ingress.yaml content
func (ig *InfrastructureGenerator) generateKubernetesIngress(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s-ingress
  namespace: %s
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - %s.com
    - api.%s.com
    secretName: %s-tls
  rules:
  - host: %s.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s-frontend-service
            port:
              number: 80
  - host: api.%s.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: %s-backend-service
            port:
              number: 80`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateTerraformMain generates main.tf content
func (ig *InfrastructureGenerator) generateTerraformMain(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Infrastructure
# Generated by Open Source Project Generator

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    # Configure your S3 backend here
    # bucket = "%s-terraform-state"
    # key    = "terraform.tfstate"
    # region = "us-west-2"
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"
  
  project_name = var.project_name
  environment  = var.environment
  vpc_cidr     = var.vpc_cidr
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "%%{var.project_name}-%%{var.environment}"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "%%{var.project_name}-%%{var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = module.vpc.public_subnet_ids

  enable_deletion_protection = var.environment == "production"
}

# Security Group for ALB
resource "aws_security_group" "alb" {
  name_prefix = "%%{var.project_name}-%%{var.environment}-alb-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}`, config.Name, config.Name)
}

// generateTerraformVariables generates variables.tf content
func (ig *InfrastructureGenerator) generateTerraformVariables(config *models.ProjectConfig) string {
	return fmt.Sprintf(`variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "%s"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b"]
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "min_capacity" {
  description = "Minimum number of instances"
  type        = number
  default     = 1
}

variable "max_capacity" {
  description = "Maximum number of instances"
  type        = number
  default     = 3
}

variable "desired_capacity" {
  description = "Desired number of instances"
  type        = number
  default     = 2
}`, config.Name)
}

// generateTerraformOutputs generates outputs.tf content
func (ig *InfrastructureGenerator) generateTerraformOutputs(config *models.ProjectConfig) string {
	return `output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

output "load_balancer_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "load_balancer_zone_id" {
  description = "Zone ID of the load balancer"
  value       = aws_lb.main.zone_id
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}`
}

// generateTerraformTfVarsExample generates terraform.tfvars.example content
func (ig *InfrastructureGenerator) generateTerraformTfVarsExample(config *models.ProjectConfig) string {
	return fmt.Sprintf(`# %s Terraform Variables Example
# Copy this file to terraform.tfvars and update the values

aws_region = "us-west-2"
project_name = "%s"
environment = "dev"

# VPC Configuration
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-west-2a", "us-west-2b"]

# Instance Configuration
instance_type = "t3.micro"
min_capacity = 1
max_capacity = 3
desired_capacity = 2`, config.Name, config.Name)
}

// generateTerraformVPCModule generates VPC module content
func (ig *InfrastructureGenerator) generateTerraformVPCModule(config *models.ProjectConfig) string {
	return `variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.project_name}-${var.environment}-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-${var.environment}-igw"
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count = 2

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-${var.environment}-public-${count.index + 1}"
    Type = "public"
  }
}

# Private Subnets
resource "aws_subnet" "private" {
  count = 2

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 10)
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = {
    Name = "${var.project_name}-${var.environment}-private-${count.index + 1}"
    Type = "private"
  }
}

# Route Table for Public Subnets
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-public-rt"
  }
}

# Route Table Association for Public Subnets
resource "aws_route_table_association" "public" {
  count = length(aws_subnet.public)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Data source for availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

# Outputs
output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_ids" {
  value = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  value = aws_subnet.private[*].id
}

output "internet_gateway_id" {
  value = aws_internet_gateway.main.id
}`
}
