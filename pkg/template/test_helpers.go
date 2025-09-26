package template

import (
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
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
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:    true,
					Home:   true,
					Admin:  true,
					Shared: true,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
				Shared:  true,
			},
			Infrastructure: models.InfrastructureComponents{
				Terraform:  true,
				Kubernetes: true,
				Docker:     true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "20.0.0",
			Go:   "1.22",
			Packages: map[string]string{
				"typescript": "5.0.0",
				"eslint":     "8.0.0",
			},
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
