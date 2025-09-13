package template

import (
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// createTestProjectConfig creates test data for template generation
func createTestProjectConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:         "testproject",
		Organization: "testorg",
		Description:  "A test project for template validation",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/testorg/testproject",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    true,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Node:   "20.0.0",
			Go:     "1.22",
			Kotlin: "1.9.0",
			Swift:  "5.9",
			NextJS: "14.0.0",
			React:  "18.0.0",
			Packages: map[string]string{
				"typescript": "5.0.0",
				"eslint":     "8.0.0",
			},
			UpdatedAt: time.Now(),
		},
		CustomVars: map[string]string{
			"DATABASE_URL": "postgresql://localhost:5432/testdb",
			"REDIS_URL":    "redis://localhost:6379",
		},
		OutputPath:       "output",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}
}

// createCompilationTestData creates comprehensive test data for template compilation
func createCompilationTestData() *models.ProjectConfig {
	return createTestProjectConfig()
}
