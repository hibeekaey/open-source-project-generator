package components

import (
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestInfrastructureGenerator_GenerateFiles(t *testing.T) {
	tests := []struct {
		name          string
		config        *models.ProjectConfig
		expectedFiles []string
		expectedError bool
	}{
		{
			name: "Generate all infrastructure components",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Infrastructure: models.InfrastructureComponents{
						Docker:     true,
						Kubernetes: true,
						Terraform:  true,
					},
				},
			},
			expectedFiles: []string{
				// Docker files
				"testproject/Deploy/docker/docker-compose.yml",
				"testproject/Deploy/docker/docker-compose.dev.yml",
				"testproject/Deploy/docker/docker-compose.prod.yml",
				"testproject/Deploy/docker/.dockerignore",
				"testproject/Deploy/docker/Dockerfile.frontend",
				"testproject/Deploy/docker/Dockerfile.backend",
				// Kubernetes files
				"testproject/Deploy/kubernetes/base/namespace.yaml",
				"testproject/Deploy/kubernetes/base/configmap.yaml",
				"testproject/Deploy/kubernetes/base/secret.yaml",
				"testproject/Deploy/kubernetes/base/backend-deployment.yaml",
				"testproject/Deploy/kubernetes/base/frontend-deployment.yaml",
				"testproject/Deploy/kubernetes/base/services.yaml",
				"testproject/Deploy/kubernetes/base/ingress.yaml",
				// Terraform files
				"testproject/Deploy/terraform/main.tf",
				"testproject/Deploy/terraform/variables.tf",
				"testproject/Deploy/terraform/outputs.tf",
				"testproject/Deploy/terraform/terraform.tfvars.example",
				"testproject/Deploy/terraform/modules/vpc/main.tf",
			},
			expectedError: false,
		},
		{
			name: "Generate only Docker",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Infrastructure: models.InfrastructureComponents{
						Docker: true,
					},
				},
			},
			expectedFiles: []string{
				"testproject/Deploy/docker/docker-compose.yml",
				"testproject/Deploy/docker/docker-compose.dev.yml",
				"testproject/Deploy/docker/docker-compose.prod.yml",
				"testproject/Deploy/docker/.dockerignore",
				"testproject/Deploy/docker/Dockerfile.frontend",
				"testproject/Deploy/docker/Dockerfile.backend",
			},
			expectedError: false,
		},
		{
			name:          "Nil config should return error",
			config:        nil,
			expectedFiles: []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := NewMockFileSystemOperations()
			ig := NewInfrastructureGenerator(mockFS)

			err := ig.GenerateFiles("testproject", tt.config)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that expected files were created
			for _, expectedFile := range tt.expectedFiles {
				if !mockFS.FileExists(expectedFile) {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}
		})
	}
}

func TestInfrastructureGenerator_generateDockerCompose(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateDockerCompose(config)

	expectedElements := []string{
		"version: '3.8'",
		"services:",
		"# testapp Backend Service",
		"backend:",
		"container_name: testapp-backend",
		"ports:",
		"- \"8080:8080\"",
		"# testapp Frontend Service",
		"frontend:",
		"container_name: testapp-frontend",
		"- \"3000:3000\"",
		"# PostgreSQL Database",
		"postgres:",
		"image: postgres:15-alpine",
		"POSTGRES_DB=testapp_db",
		"# Redis Cache",
		"redis:",
		"image: redis:7-alpine",
		"networks:",
		"testapp-network:",
		"volumes:",
		"postgres_data:",
		"redis_data:",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("docker-compose.yml should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateDockerComposeDev(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateDockerComposeDev(config)

	expectedElements := []string{
		"version: '3.8'",
		"services:",
		"backend:",
		"ENVIRONMENT=development",
		"LOG_LEVEL=debug",
		"HOT_RELOAD=true",
		"command: [\"go\", \"run\", \"main.go\"]",
		"frontend:",
		"NODE_ENV=development",
		"command: [\"npm\", \"run\", \"dev\"]",
		"postgres:",
		"POSTGRES_DB=testapp_dev_db",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("docker-compose.dev.yml should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateDockerComposeProd(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateDockerComposeProd(config)

	expectedElements := []string{
		"version: '3.8'",
		"services:",
		"backend:",
		"image: testapp-backend:latest",
		"ENVIRONMENT=production",
		"deploy:",
		"replicas: 2",
		"resources:",
		"limits:",
		"memory: 512M",
		"frontend:",
		"image: testapp-frontend:latest",
		"NODE_ENV=production",
		"postgres:",
		"POSTGRES_DB=testapp_prod_db",
		"nginx:",
		"image: nginx:alpine",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("docker-compose.prod.yml should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateFrontendDockerfile(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Node: "18",
		},
	}

	content := ig.generateFrontendDockerfile(config)

	expectedElements := []string{
		"# testapp Frontend Dockerfile",
		"FROM node:18-alpine AS base",
		"FROM base AS deps",
		"WORKDIR /app",
		"COPY package.json package-lock.json* ./",
		"RUN npm ci --only=production",
		"FROM base AS builder",
		"RUN npm run build",
		"FROM base AS runner",
		"ENV NODE_ENV production",
		"RUN addgroup --system --gid 1001 nodejs",
		"RUN adduser --system --uid 1001 nextjs",
		"USER nextjs",
		"EXPOSE 3000",
		"CMD [\"node\", \"server.js\"]",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Frontend Dockerfile should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateBackendDockerfile(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
		Versions: &models.VersionConfig{
			Go: "1.22",
		},
	}

	content := ig.generateBackendDockerfile(config)

	expectedElements := []string{
		"# testapp Backend Dockerfile",
		"FROM golang:1.22-alpine AS builder",
		"WORKDIR /app",
		"RUN apk add --no-cache git ca-certificates",
		"COPY go.mod go.sum ./",
		"RUN go mod download",
		"RUN CGO_ENABLED=0 GOOS=linux go build",
		"FROM alpine:latest",
		"RUN apk --no-cache add ca-certificates tzdata",
		"RUN adduser -D -s /bin/sh appuser",
		"USER appuser",
		"EXPOSE 8080",
		"HEALTHCHECK",
		"CMD [\"./main\"]",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Backend Dockerfile should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateKubernetesNamespace(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateKubernetesNamespace(config)

	expectedElements := []string{
		"apiVersion: v1",
		"kind: Namespace",
		"metadata:",
		"name: testapp",
		"labels:",
		"app: testapp",
		"environment: production",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Kubernetes namespace should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateKubernetesConfigMap(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateKubernetesConfigMap(config)

	expectedElements := []string{
		"apiVersion: v1",
		"kind: ConfigMap",
		"metadata:",
		"name: testapp-config",
		"namespace: testapp",
		"data:",
		"ENVIRONMENT: \"production\"",
		"LOG_LEVEL: \"info\"",
		"PORT: \"8080\"",
		"CORS_ORIGINS: \"https://testapp.com\"",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Kubernetes ConfigMap should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateKubernetesBackendDeployment(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateKubernetesBackendDeployment(config)

	expectedElements := []string{
		"apiVersion: apps/v1",
		"kind: Deployment",
		"metadata:",
		"name: testapp-backend",
		"namespace: testapp",
		"spec:",
		"replicas: 2",
		"selector:",
		"matchLabels:",
		"app: testapp-backend",
		"template:",
		"containers:",
		"- name: backend",
		"image: testapp-backend:latest",
		"ports:",
		"- containerPort: 8080",
		"livenessProbe:",
		"httpGet:",
		"path: /health",
		"readinessProbe:",
		"resources:",
		"requests:",
		"memory: \"256Mi\"",
		"cpu: \"250m\"",
		"limits:",
		"memory: \"512Mi\"",
		"cpu: \"500m\"",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Kubernetes backend deployment should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateTerraformMain(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateTerraformMain(config)

	expectedElements := []string{
		"# testapp Infrastructure",
		"terraform {",
		"required_version = \">= 1.0\"",
		"required_providers {",
		"aws = {",
		"source  = \"hashicorp/aws\"",
		"version = \"~> 5.0\"",
		"backend \"s3\" {",
		"provider \"aws\" {",
		"region = var.aws_region",
		"default_tags {",
		"Project     = var.project_name",
		"Environment = var.environment",
		"ManagedBy   = \"terraform\"",
		"module \"vpc\" {",
		"source = \"./modules/vpc\"",
		"resource \"aws_ecs_cluster\" \"main\" {",
		"resource \"aws_lb\" \"main\" {",
		"resource \"aws_security_group\" \"alb\" {",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Terraform main.tf should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateTerraformVariables(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateTerraformVariables(config)

	expectedElements := []string{
		"variable \"aws_region\" {",
		"description = \"AWS region\"",
		"type        = string",
		"default     = \"us-west-2\"",
		"variable \"project_name\" {",
		"default     = \"testapp\"",
		"variable \"environment\" {",
		"variable \"vpc_cidr\" {",
		"variable \"availability_zones\" {",
		"variable \"instance_type\" {",
		"variable \"min_capacity\" {",
		"variable \"max_capacity\" {",
		"variable \"desired_capacity\" {",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Terraform variables.tf should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateTerraformOutputs(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateTerraformOutputs(config)

	expectedElements := []string{
		"output \"vpc_id\" {",
		"description = \"ID of the VPC\"",
		"value       = module.vpc.vpc_id",
		"output \"public_subnet_ids\" {",
		"output \"private_subnet_ids\" {",
		"output \"load_balancer_dns_name\" {",
		"output \"load_balancer_zone_id\" {",
		"output \"ecs_cluster_name\" {",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Terraform outputs.tf should contain %s", element)
		}
	}
}

func TestInfrastructureGenerator_generateTerraformVPCModule(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	ig := NewInfrastructureGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := ig.generateTerraformVPCModule(config)

	expectedElements := []string{
		"variable \"project_name\" {",
		"variable \"environment\" {",
		"variable \"vpc_cidr\" {",
		"resource \"aws_vpc\" \"main\" {",
		"cidr_block           = var.vpc_cidr",
		"enable_dns_hostnames = true",
		"resource \"aws_internet_gateway\" \"main\" {",
		"resource \"aws_subnet\" \"public\" {",
		"count = 2",
		"map_public_ip_on_launch = true",
		"resource \"aws_subnet\" \"private\" {",
		"resource \"aws_route_table\" \"public\" {",
		"resource \"aws_route_table_association\" \"public\" {",
		"data \"aws_availability_zones\" \"available\" {",
		"output \"vpc_id\" {",
		"output \"public_subnet_ids\" {",
		"output \"private_subnet_ids\" {",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Terraform VPC module should contain %s", element)
		}
	}
}
