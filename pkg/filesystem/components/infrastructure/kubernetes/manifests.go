package kubernetes

import (
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ManifestGenerator handles Kubernetes manifest generation
type ManifestGenerator struct{}

// NewManifestGenerator creates a new Kubernetes manifest generator
func NewManifestGenerator() *ManifestGenerator {
	return &ManifestGenerator{}
}

// GenerateNamespace generates namespace.yaml content
func (mg *ManifestGenerator) GenerateNamespace(config *models.ProjectConfig) string {
	return fmt.Sprintf(`apiVersion: v1
kind: Namespace
metadata:
  name: %s
  labels:
    app: %s
    environment: production`, config.Name, config.Name)
}

// GenerateConfigMap generates configmap.yaml content
func (mg *ManifestGenerator) GenerateConfigMap(config *models.ProjectConfig) string {
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

// GenerateSecret generates secret.yaml content
func (mg *ManifestGenerator) GenerateSecret(config *models.ProjectConfig) string {
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

// GenerateBackendDeployment generates backend deployment
func (mg *ManifestGenerator) GenerateBackendDeployment(config *models.ProjectConfig) string {
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

// GenerateFrontendDeployment generates frontend deployment
func (mg *ManifestGenerator) GenerateFrontendDeployment(config *models.ProjectConfig) string {
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

// GenerateServices generates services.yaml content
func (mg *ManifestGenerator) GenerateServices(config *models.ProjectConfig) string {
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

// GenerateIngress generates ingress.yaml content
func (mg *ManifestGenerator) GenerateIngress(config *models.ProjectConfig) string {
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
