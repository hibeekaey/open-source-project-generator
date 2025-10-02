package docker

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ComposeGenerator handles Docker Compose file generation
type ComposeGenerator struct{}

// NewComposeGenerator creates a new Docker Compose generator
func NewComposeGenerator() *ComposeGenerator {
	return &ComposeGenerator{}
}

// GenerateDockerCompose generates docker-compose.yml content
func (cg *ComposeGenerator) GenerateDockerCompose(config *models.ProjectConfig) string {
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

// GenerateDockerComposeDev generates docker-compose.dev.yml content
func (cg *ComposeGenerator) GenerateDockerComposeDev(config *models.ProjectConfig) string {
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

// GenerateDockerComposeProd generates docker-compose.prod.yml content
func (cg *ComposeGenerator) GenerateDockerComposeProd(config *models.ProjectConfig) string {
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
