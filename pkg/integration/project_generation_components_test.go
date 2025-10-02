package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestProjectGenerationComponents tests project generation with all component types
func TestProjectGenerationComponents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("backend_components_generation", func(t *testing.T) {
		testBackendComponentsGeneration(t, tempDir)
	})

	t.Run("frontend_components_generation", func(t *testing.T) {
		testFrontendComponentsGeneration(t, tempDir)
	})

	t.Run("mobile_components_generation", func(t *testing.T) {
		testMobileComponentsGeneration(t, tempDir)
	})

	t.Run("infrastructure_components_generation", func(t *testing.T) {
		testInfrastructureComponentsGeneration(t, tempDir)
	})

	t.Run("full_stack_project_generation", func(t *testing.T) {
		testFullStackProjectGeneration(t, tempDir)
	})

	t.Run("project_generators_integration", func(t *testing.T) {
		testProjectGeneratorsIntegration(t, tempDir)
	})
}

func testBackendComponentsGeneration(t *testing.T, tempDir string) {
	// Test Go Gin backend generation
	t.Run("go_gin_backend", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "go-gin-backend-test",
			Organization: "backend-org",
			Description:  "Test Go Gin backend project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "go-gin-backend"),
			Components: models.Components{
				Backend: models.BackendComponents{
					GoGin: true,
				},
			},
		}

		// Generate backend components
		backendGenerator := NewMockBackendGenerator()
		err := backendGenerator.GenerateGoGin(config)
		if err != nil {
			t.Errorf("Go Gin backend generation failed: %v", err)
		}

		// Verify backend structure
		expectedFiles := []string{
			"main.go",
			"go.mod",
			"go.sum",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"internal/middleware/cors.go",
			"internal/config/config.go",
			"pkg/api/routes.go",
			"Dockerfile",
			"docker-compose.yml",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected backend file %s to be generated", file)
			}
		}

		// Verify Go module content
		goModPath := filepath.Join(config.OutputPath, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			content, err := os.ReadFile(goModPath)
			if err != nil {
				t.Errorf("Failed to read go.mod: %v", err)
			} else {
				if !strings.Contains(string(content), config.Name) {
					t.Error("Expected go.mod to contain project name")
				}
			}
		}

		// Verify main.go content
		mainGoPath := filepath.Join(config.OutputPath, "main.go")
		if _, err := os.Stat(mainGoPath); err == nil {
			content, err := os.ReadFile(mainGoPath)
			if err != nil {
				t.Errorf("Failed to read main.go: %v", err)
			} else {
				if !strings.Contains(string(content), "gin") {
					t.Error("Expected main.go to contain Gin framework code")
				}
			}
		}
	})

	// Test Node.js Express backend generation
	t.Run("nodejs_express_backend", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "nodejs-express-backend-test",
			Organization: "backend-org",
			Description:  "Test Node.js Express backend project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "nodejs-express-backend"),
			Components: models.Components{
				Backend: models.BackendComponents{
					GoGin: true,
				},
			},
		}

		// Generate backend components
		backendGenerator := NewMockBackendGenerator()
		err := backendGenerator.GenerateNodeExpress(config)
		if err != nil {
			t.Errorf("Node.js Express backend generation failed: %v", err)
		}

		// Verify backend structure
		expectedFiles := []string{
			"package.json",
			"server.js",
			"app.js",
			"routes/index.js",
			"routes/api.js",
			"middleware/cors.js",
			"middleware/auth.js",
			"config/database.js",
			"models/index.js",
			"controllers/health.js",
			"Dockerfile",
			".env.example",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected backend file %s to be generated", file)
			}
		}

		// Verify package.json content
		packageJsonPath := filepath.Join(config.OutputPath, "package.json")
		if _, err := os.Stat(packageJsonPath); err == nil {
			content, err := os.ReadFile(packageJsonPath)
			if err != nil {
				t.Errorf("Failed to read package.json: %v", err)
			} else {
				if !strings.Contains(string(content), config.Name) {
					t.Error("Expected package.json to contain project name")
				}
				if !strings.Contains(string(content), "express") {
					t.Error("Expected package.json to contain Express dependency")
				}
			}
		}
	})
}

func testFrontendComponentsGeneration(t *testing.T, tempDir string) {
	// Test Next.js App generation
	t.Run("nextjs_app_frontend", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "nextjs-app-frontend-test",
			Organization: "frontend-org",
			Description:  "Test Next.js App frontend project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "nextjs-app-frontend"),
			Components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
		}

		// Generate frontend components
		frontendGenerator := NewMockFrontendGenerator()
		err := frontendGenerator.GenerateNextJSApp(config)
		if err != nil {
			t.Errorf("Next.js App frontend generation failed: %v", err)
		}

		// Verify frontend structure
		expectedFiles := []string{
			"package.json",
			"next.config.js",
			"tailwind.config.js",
			"tsconfig.json",
			"app/layout.tsx",
			"app/page.tsx",
			"app/globals.css",
			"components/ui/Button.tsx",
			"components/ui/Card.tsx",
			"lib/utils.ts",
			"public/favicon.ico",
			"public/next.svg",
			".env.local.example",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected frontend file %s to be generated", file)
			}
		}

		// Verify package.json content
		packageJsonPath := filepath.Join(config.OutputPath, "package.json")
		if _, err := os.Stat(packageJsonPath); err == nil {
			content, err := os.ReadFile(packageJsonPath)
			if err != nil {
				t.Errorf("Failed to read package.json: %v", err)
			} else {
				if !strings.Contains(string(content), "next") {
					t.Error("Expected package.json to contain Next.js dependency")
				}
				if !strings.Contains(string(content), "react") {
					t.Error("Expected package.json to contain React dependency")
				}
				if !strings.Contains(string(content), "tailwindcss") {
					t.Error("Expected package.json to contain Tailwind CSS dependency")
				}
			}
		}

		// Verify TypeScript configuration
		tsconfigPath := filepath.Join(config.OutputPath, "tsconfig.json")
		if _, err := os.Stat(tsconfigPath); err == nil {
			content, err := os.ReadFile(tsconfigPath)
			if err != nil {
				t.Errorf("Failed to read tsconfig.json: %v", err)
			} else {
				if !strings.Contains(string(content), "next") {
					t.Error("Expected tsconfig.json to contain Next.js configuration")
				}
			}
		}
	})

	// Test Next.js Admin Dashboard generation
	t.Run("nextjs_admin_frontend", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "nextjs-admin-frontend-test",
			Organization: "frontend-org",
			Description:  "Test Next.js Admin Dashboard frontend project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "nextjs-admin-frontend"),
			Components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						Admin: true,
					},
				},
			},
		}

		// Generate frontend components
		frontendGenerator := NewMockFrontendGenerator()
		err := frontendGenerator.GenerateNextJSAdmin(config)
		if err != nil {
			t.Errorf("Next.js Admin frontend generation failed: %v", err)
		}

		// Verify admin-specific structure
		expectedFiles := []string{
			"package.json",
			"app/dashboard/page.tsx",
			"app/dashboard/layout.tsx",
			"app/dashboard/users/page.tsx",
			"app/dashboard/settings/page.tsx",
			"components/dashboard/Sidebar.tsx",
			"components/dashboard/Header.tsx",
			"components/dashboard/DataTable.tsx",
			"components/charts/LineChart.tsx",
			"components/charts/BarChart.tsx",
			"lib/auth.ts",
			"lib/api.ts",
			"hooks/useAuth.ts",
			"hooks/useApi.ts",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected admin file %s to be generated", file)
			}
		}
	})

	// Test React Component Library generation
	t.Run("react_component_library", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "react-components-test",
			Organization: "frontend-org",
			Description:  "Test React Component Library project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "react-components"),
			Components: models.Components{
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
		}

		// Generate frontend components
		frontendGenerator := NewMockFrontendGenerator()
		err := frontendGenerator.GenerateReactComponents(config)
		if err != nil {
			t.Errorf("React Component Library generation failed: %v", err)
		}

		// Verify component library structure
		expectedFiles := []string{
			"package.json",
			"rollup.config.js",
			"tsconfig.json",
			"src/index.ts",
			"src/components/Button/Button.tsx",
			"src/components/Button/Button.stories.tsx",
			"src/components/Button/Button.test.tsx",
			"src/components/Card/Card.tsx",
			"src/components/Input/Input.tsx",
			"src/hooks/index.ts",
			"src/utils/index.ts",
			".storybook/main.js",
			".storybook/preview.js",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected component library file %s to be generated", file)
			}
		}
	})
}

func testMobileComponentsGeneration(t *testing.T, tempDir string) {
	// Test Android Kotlin generation
	t.Run("android_kotlin_mobile", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "android-kotlin-mobile-test",
			Organization: "mobile-org",
			Description:  "Test Android Kotlin mobile project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "android-kotlin-mobile"),
			Components: models.Components{
				Mobile: models.MobileComponents{
					Android: true,
				},
			},
		}

		// Generate mobile components
		mobileGenerator := NewMockMobileGenerator()
		err := mobileGenerator.GenerateAndroidKotlin(config)
		if err != nil {
			t.Errorf("Android Kotlin mobile generation failed: %v", err)
		}

		// Verify Android structure
		expectedFiles := []string{
			"build.gradle",
			"settings.gradle",
			"gradle.properties",
			"app/build.gradle",
			"app/src/main/AndroidManifest.xml",
			"app/src/main/java/com/example/MainActivity.kt",
			"app/src/main/java/com/example/ui/theme/Theme.kt",
			"app/src/main/java/com/example/ui/components/Button.kt",
			"app/src/main/res/layout/activity_main.xml",
			"app/src/main/res/values/strings.xml",
			"app/src/main/res/values/colors.xml",
			"app/src/test/java/com/example/ExampleUnitTest.kt",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected Android file %s to be generated", file)
			}
		}

		// Verify build.gradle content
		buildGradlePath := filepath.Join(config.OutputPath, "app/build.gradle")
		if _, err := os.Stat(buildGradlePath); err == nil {
			content, err := os.ReadFile(buildGradlePath)
			if err != nil {
				t.Errorf("Failed to read build.gradle: %v", err)
			} else {
				if !strings.Contains(string(content), "kotlin") {
					t.Error("Expected build.gradle to contain Kotlin configuration")
				}
			}
		}
	})

	// Test iOS Swift generation
	t.Run("ios_swift_mobile", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "ios-swift-mobile-test",
			Organization: "mobile-org",
			Description:  "Test iOS Swift mobile project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "ios-swift-mobile"),
			Components: models.Components{
				Mobile: models.MobileComponents{
					IOS: true,
				},
			},
		}

		// Generate mobile components
		mobileGenerator := NewMockMobileGenerator()
		err := mobileGenerator.GenerateiOSSwift(config)
		if err != nil {
			t.Errorf("iOS Swift mobile generation failed: %v", err)
		}

		// Verify iOS structure
		expectedFiles := []string{
			"Package.swift",
			"Sources/App/main.swift",
			"Sources/App/ContentView.swift",
			"Sources/App/Models/User.swift",
			"Sources/App/Views/HomeView.swift",
			"Sources/App/Views/Components/Button.swift",
			"Sources/App/Services/APIService.swift",
			"Sources/App/Utils/Extensions.swift",
			"Tests/AppTests/AppTests.swift",
			"Resources/Assets.xcassets/Contents.json",
			"Resources/Info.plist",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected iOS file %s to be generated", file)
			}
		}

		// Verify Package.swift content
		packageSwiftPath := filepath.Join(config.OutputPath, "Package.swift")
		if _, err := os.Stat(packageSwiftPath); err == nil {
			content, err := os.ReadFile(packageSwiftPath)
			if err != nil {
				t.Errorf("Failed to read Package.swift: %v", err)
			} else {
				if !strings.Contains(string(content), config.Name) {
					t.Error("Expected Package.swift to contain project name")
				}
			}
		}
	})

	// Test React Native generation
	t.Run("react_native_mobile", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "react-native-mobile-test",
			Organization: "mobile-org",
			Description:  "Test React Native mobile project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "react-native-mobile"),
			Components: models.Components{
				Mobile: models.MobileComponents{
					Android: true,
					IOS:     true,
				},
			},
		}

		// Generate mobile components
		mobileGenerator := NewMockMobileGenerator()
		err := mobileGenerator.GenerateAndroid(config)
		if err != nil {
			t.Errorf("Android mobile generation failed: %v", err)
		}

		// Verify React Native structure
		expectedFiles := []string{
			"package.json",
			"metro.config.js",
			"babel.config.js",
			"App.tsx",
			"index.js",
			"src/screens/HomeScreen.tsx",
			"src/components/Button.tsx",
			"src/navigation/AppNavigator.tsx",
			"src/services/api.ts",
			"src/utils/helpers.ts",
			"android/build.gradle",
			"ios/Podfile",
			"__tests__/App-test.tsx",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected React Native file %s to be generated", file)
			}
		}
	})
}

func testInfrastructureComponentsGeneration(t *testing.T, tempDir string) {
	// Test Docker infrastructure generation
	t.Run("docker_infrastructure", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "docker-infrastructure-test",
			Organization: "infra-org",
			Description:  "Test Docker infrastructure project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "docker-infrastructure"),
			Components: models.Components{
				Infrastructure: models.InfrastructureComponents{
					Docker: true,
				},
			},
		}

		// Generate infrastructure components
		infraGenerator := NewMockInfrastructureGenerator()
		err := infraGenerator.GenerateDocker(config)
		if err != nil {
			t.Errorf("Docker infrastructure generation failed: %v", err)
		}

		// Verify Docker structure
		expectedFiles := []string{
			"Dockerfile",
			"docker-compose.yml",
			"docker-compose.prod.yml",
			"docker-compose.dev.yml",
			".dockerignore",
			"scripts/docker-build.sh",
			"scripts/docker-run.sh",
			"config/nginx.conf",
			"config/redis.conf",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected Docker file %s to be generated", file)
			}
		}

		// Verify docker-compose.yml content
		dockerComposePath := filepath.Join(config.OutputPath, "docker-compose.yml")
		if _, err := os.Stat(dockerComposePath); err == nil {
			content, err := os.ReadFile(dockerComposePath)
			if err != nil {
				t.Errorf("Failed to read docker-compose.yml: %v", err)
			} else {
				if !strings.Contains(string(content), "version:") {
					t.Error("Expected docker-compose.yml to contain version")
				}
				if !strings.Contains(string(content), "services:") {
					t.Error("Expected docker-compose.yml to contain services")
				}
			}
		}
	})

	// Test Kubernetes infrastructure generation
	t.Run("kubernetes_infrastructure", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "kubernetes-infrastructure-test",
			Organization: "infra-org",
			Description:  "Test Kubernetes infrastructure project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "kubernetes-infrastructure"),
			Components: models.Components{
				Infrastructure: models.InfrastructureComponents{
					Kubernetes: true,
				},
			},
		}

		// Generate infrastructure components
		infraGenerator := NewMockInfrastructureGenerator()
		err := infraGenerator.GenerateKubernetes(config)
		if err != nil {
			t.Errorf("Kubernetes infrastructure generation failed: %v", err)
		}

		// Verify Kubernetes structure
		expectedFiles := []string{
			"k8s/namespace.yaml",
			"k8s/deployment.yaml",
			"k8s/service.yaml",
			"k8s/ingress.yaml",
			"k8s/configmap.yaml",
			"k8s/secret.yaml",
			"k8s/hpa.yaml",
			"helm/Chart.yaml",
			"helm/values.yaml",
			"helm/templates/deployment.yaml",
			"helm/templates/service.yaml",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected Kubernetes file %s to be generated", file)
			}
		}

		// Verify deployment.yaml content
		deploymentPath := filepath.Join(config.OutputPath, "k8s/deployment.yaml")
		if _, err := os.Stat(deploymentPath); err == nil {
			content, err := os.ReadFile(deploymentPath)
			if err != nil {
				t.Errorf("Failed to read deployment.yaml: %v", err)
			} else {
				if !strings.Contains(string(content), "apiVersion:") {
					t.Error("Expected deployment.yaml to contain apiVersion")
				}
				if !strings.Contains(string(content), "kind: Deployment") {
					t.Error("Expected deployment.yaml to contain Deployment kind")
				}
			}
		}
	})

	// Test Terraform infrastructure generation
	t.Run("terraform_infrastructure", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "terraform-infrastructure-test",
			Organization: "infra-org",
			Description:  "Test Terraform infrastructure project",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "terraform-infrastructure"),
			Components: models.Components{
				Infrastructure: models.InfrastructureComponents{
					Terraform: true,
				},
			},
		}

		// Generate infrastructure components
		infraGenerator := NewMockInfrastructureGenerator()
		err := infraGenerator.GenerateTerraform(config)
		if err != nil {
			t.Errorf("Terraform infrastructure generation failed: %v", err)
		}

		// Verify Terraform structure
		expectedFiles := []string{
			"main.tf",
			"variables.tf",
			"outputs.tf",
			"versions.tf",
			"terraform.tfvars.example",
			"modules/vpc/main.tf",
			"modules/vpc/variables.tf",
			"modules/vpc/outputs.tf",
			"modules/compute/main.tf",
			"modules/database/main.tf",
			"environments/dev/main.tf",
			"environments/prod/main.tf",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected Terraform file %s to be generated", file)
			}
		}

		// Verify main.tf content
		mainTfPath := filepath.Join(config.OutputPath, "main.tf")
		if _, err := os.Stat(mainTfPath); err == nil {
			content, err := os.ReadFile(mainTfPath)
			if err != nil {
				t.Errorf("Failed to read main.tf: %v", err)
			} else {
				if !strings.Contains(string(content), "terraform") {
					t.Error("Expected main.tf to contain terraform configuration")
				}
			}
		}
	})
}

func testFullStackProjectGeneration(t *testing.T, tempDir string) {
	// Test complete full-stack project generation
	config := &models.ProjectConfig{
		Name:         "fullstack-project-test",
		Organization: "fullstack-org",
		Description:  "Test full-stack project with all components",
		License:      "MIT",
		OutputPath:   filepath.Join(tempDir, "fullstack-project"),
		Components: models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Admin: true,
				},
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
			},
		},
	}

	// Generate complete project
	projectGenerator := NewMockProjectGenerator()
	err := projectGenerator.GenerateProject(config)
	if err != nil {
		t.Errorf("Full-stack project generation failed: %v", err)
	}

	// Verify project structure
	expectedDirs := []string{
		"backend",
		"frontend/app",
		"frontend/admin",
		"mobile",
		"infrastructure/docker",
		"infrastructure/k8s",
		"docs",
		"scripts",
	}

	for _, dir := range expectedDirs {
		dirPath := filepath.Join(config.OutputPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to be generated", dir)
		}
	}

	// Verify root-level files
	expectedRootFiles := []string{
		"README.md",
		"LICENSE",
		"Makefile",
		"docker-compose.yml",
		".gitignore",
		".env.example",
	}

	for _, file := range expectedRootFiles {
		filePath := filepath.Join(config.OutputPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected root file %s to be generated", file)
		}
	}

	// Verify README.md content
	readmePath := filepath.Join(config.OutputPath, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		content, err := os.ReadFile(readmePath)
		if err != nil {
			t.Errorf("Failed to read README.md: %v", err)
		} else {
			if !strings.Contains(string(content), config.Name) {
				t.Error("Expected README.md to contain project name")
			}
			if !strings.Contains(string(content), config.Description) {
				t.Error("Expected README.md to contain project description")
			}
		}
	}

	// Verify Makefile content
	makefilePath := filepath.Join(config.OutputPath, "Makefile")
	if _, err := os.Stat(makefilePath); err == nil {
		content, err := os.ReadFile(makefilePath)
		if err != nil {
			t.Errorf("Failed to read Makefile: %v", err)
		} else {
			if !strings.Contains(string(content), "build") {
				t.Error("Expected Makefile to contain build target")
			}
			if !strings.Contains(string(content), "test") {
				t.Error("Expected Makefile to contain test target")
			}
		}
	}
}

func testProjectGeneratorsIntegration(t *testing.T, tempDir string) {
	// Test structure generator
	t.Run("structure_generator", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "structure-generator-test",
			Organization: "generator-org",
			OutputPath:   filepath.Join(tempDir, "structure-generator"),
		}

		structureGenerator := NewMockStructureGenerator()
		err := structureGenerator.GenerateStructure(config)
		if err != nil {
			t.Errorf("Structure generator failed: %v", err)
		}

		// Verify basic structure was created
		expectedDirs := []string{
			"src",
			"tests",
			"docs",
			"scripts",
		}

		for _, dir := range expectedDirs {
			dirPath := filepath.Join(config.OutputPath, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Errorf("Expected structure directory %s to be generated", dir)
			}
		}
	})

	// Test template generator
	t.Run("template_generator", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "template-generator-test",
			Organization: "generator-org",
			OutputPath:   filepath.Join(tempDir, "template-generator"),
		}

		templateGenerator := NewMockTemplateGenerator()
		err := templateGenerator.GenerateTemplates(config)
		if err != nil {
			t.Errorf("Template generator failed: %v", err)
		}

		// Verify templates were processed
		expectedFiles := []string{
			"README.md",
			"LICENSE",
			".gitignore",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected template file %s to be generated", file)
			}
		}
	})

	// Test configuration generator
	t.Run("configuration_generator", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "config-generator-test",
			Organization: "generator-org",
			OutputPath:   filepath.Join(tempDir, "config-generator"),
		}

		configGenerator := NewMockConfigurationGenerator()
		err := configGenerator.GenerateConfiguration(config)
		if err != nil {
			t.Errorf("Configuration generator failed: %v", err)
		}

		// Verify configuration files were created
		expectedFiles := []string{
			"config/app.yaml",
			"config/database.yaml",
			".env.example",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected config file %s to be generated", file)
			}
		}
	})

	// Test documentation generator
	t.Run("documentation_generator", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "docs-generator-test",
			Organization: "generator-org",
			Description:  "Test project for documentation generation",
			OutputPath:   filepath.Join(tempDir, "docs-generator"),
		}

		docsGenerator := NewMockDocumentationGenerator()
		err := docsGenerator.GenerateDocumentation(config)
		if err != nil {
			t.Errorf("Documentation generator failed: %v", err)
		}

		// Verify documentation files were created
		expectedFiles := []string{
			"docs/README.md",
			"docs/API.md",
			"docs/CONTRIBUTING.md",
			"docs/DEPLOYMENT.md",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected documentation file %s to be generated", file)
			}
		}

		// Verify documentation content
		apiDocsPath := filepath.Join(config.OutputPath, "docs/API.md")
		if _, err := os.Stat(apiDocsPath); err == nil {
			content, err := os.ReadFile(apiDocsPath)
			if err != nil {
				t.Errorf("Failed to read API.md: %v", err)
			} else {
				if !strings.Contains(string(content), "API") {
					t.Error("Expected API.md to contain API documentation")
				}
			}
		}
	})

	// Test CI/CD generator
	t.Run("cicd_generator", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "cicd-generator-test",
			Organization: "generator-org",
			OutputPath:   filepath.Join(tempDir, "cicd-generator"),
		}

		cicdGenerator := NewMockCICDGenerator()
		err := cicdGenerator.GenerateCICD(config)
		if err != nil {
			t.Errorf("CI/CD generator failed: %v", err)
		}

		// Verify CI/CD files were created
		expectedFiles := []string{
			".github/workflows/ci.yml",
			".github/workflows/cd.yml",
			".github/workflows/test.yml",
			".gitlab-ci.yml",
			"Jenkinsfile",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(config.OutputPath, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected CI/CD file %s to be generated", file)
			}
		}

		// Verify GitHub Actions workflow content
		ciWorkflowPath := filepath.Join(config.OutputPath, ".github/workflows/ci.yml")
		if _, err := os.Stat(ciWorkflowPath); err == nil {
			content, err := os.ReadFile(ciWorkflowPath)
			if err != nil {
				t.Errorf("Failed to read ci.yml: %v", err)
			} else {
				if !strings.Contains(string(content), "name:") {
					t.Error("Expected ci.yml to contain workflow name")
				}
				if !strings.Contains(string(content), "on:") {
					t.Error("Expected ci.yml to contain workflow triggers")
				}
			}
		}
	})
}

// Mock implementations for project generation testing

// Mock Component Generators
type MockBackendGenerator struct{}

func NewMockBackendGenerator() *MockBackendGenerator {
	return &MockBackendGenerator{}
}

func (m *MockBackendGenerator) GenerateGoGin(config *models.ProjectConfig) error {
	// Create mock Go Gin project structure
	return m.createMockFiles(config.OutputPath, map[string]string{
		"main.go":                     "package main\n\nfunc main() {\n\t// " + config.Name + " application\n}",
		"go.mod":                      "module " + config.Name + "\n\ngo 1.21",
		"go.sum":                      "// Go dependencies",
		"cmd/server/main.go":          "package main\n\nfunc main() {\n\t// Server entry point\n}",
		"internal/handlers/health.go": "package handlers\n\n// Health check handler",
		"internal/middleware/cors.go": "package middleware\n\n// CORS middleware",
		"internal/config/config.go":   "package config\n\n// Configuration",
		"pkg/api/routes.go":           "package api\n\n// API routes",
		"Dockerfile":                  "FROM golang:1.21\n\nWORKDIR /app",
		"docker-compose.yml":          "version: '3.8'\nservices:\n  app:\n    build: .",
	})
}

func (m *MockBackendGenerator) GenerateNodeExpress(config *models.ProjectConfig) error {
	// Create mock Node.js Express project structure
	return m.createMockFiles(config.OutputPath, map[string]string{
		"package.json":          `{"name": "` + config.Name + `", "version": "1.0.0", "dependencies": {"express": "^4.18.0"}}`,
		"server.js":             "const express = require('express');\nconst app = express();",
		"app.js":                "// Express application",
		"routes/index.js":       "// Index routes",
		"routes/api.js":         "// API routes",
		"middleware/cors.js":    "// CORS middleware",
		"middleware/auth.js":    "// Authentication middleware",
		"config/database.js":    "// Database configuration",
		"models/index.js":       "// Data models",
		"controllers/health.js": "// Health controller",
		"Dockerfile":            "FROM node:18\n\nWORKDIR /app",
		".env.example":          "NODE_ENV=development\nPORT=3000",
	})
}

func (m *MockBackendGenerator) createMockFiles(basePath string, files map[string]string) error {
	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockFrontendGenerator struct{}

func NewMockFrontendGenerator() *MockFrontendGenerator {
	return &MockFrontendGenerator{}
}

func (m *MockFrontendGenerator) GenerateNextJSApp(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"package.json":             `{"name": "` + config.Name + `", "dependencies": {"next": "^14.0.0", "react": "^18.0.0", "tailwindcss": "^3.0.0"}}`,
		"next.config.js":           "/** @type {import('next').NextConfig} */\nmodule.exports = {}",
		"tailwind.config.js":       "module.exports = { content: ['./app/**/*.{js,ts,jsx,tsx}'] }",
		"tsconfig.json":            `{"compilerOptions": {"target": "es5", "lib": ["dom"]}, "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx"]}`,
		"app/layout.tsx":           "export default function RootLayout({ children }) { return <html><body>{children}</body></html> }",
		"app/page.tsx":             "export default function Home() { return <div>Welcome to " + config.Name + "</div> }",
		"app/globals.css":          "@tailwind base;\n@tailwind components;\n@tailwind utilities;",
		"components/ui/Button.tsx": "export function Button() { return <button>Button</button> }",
		"components/ui/Card.tsx":   "export function Card() { return <div>Card</div> }",
		"lib/utils.ts":             "export function cn(...classes: string[]) { return classes.join(' ') }",
		"public/favicon.ico":       "// Favicon",
		"public/next.svg":          "// Next.js logo",
		".env.local.example":       "NEXT_PUBLIC_API_URL=http://localhost:3000",
	})
}

func (m *MockFrontendGenerator) GenerateNextJSAdmin(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"package.json":                       `{"name": "` + config.Name + `", "dependencies": {"next": "^14.0.0", "react": "^18.0.0"}}`,
		"app/dashboard/page.tsx":             "export default function Dashboard() { return <div>Dashboard</div> }",
		"app/dashboard/layout.tsx":           "export default function DashboardLayout({ children }) { return <div>{children}</div> }",
		"app/dashboard/users/page.tsx":       "export default function Users() { return <div>Users</div> }",
		"app/dashboard/settings/page.tsx":    "export default function Settings() { return <div>Settings</div> }",
		"components/dashboard/Sidebar.tsx":   "export function Sidebar() { return <nav>Sidebar</nav> }",
		"components/dashboard/Header.tsx":    "export function Header() { return <header>Header</header> }",
		"components/dashboard/DataTable.tsx": "export function DataTable() { return <table>Data Table</table> }",
		"components/charts/LineChart.tsx":    "export function LineChart() { return <div>Line Chart</div> }",
		"components/charts/BarChart.tsx":     "export function BarChart() { return <div>Bar Chart</div> }",
		"lib/auth.ts":                        "export function authenticate() { return true }",
		"lib/api.ts":                         "export function apiCall() { return fetch('/api') }",
		"hooks/useAuth.ts":                   "export function useAuth() { return { user: null } }",
		"hooks/useApi.ts":                    "export function useApi() { return { data: null } }",
	})
}

func (m *MockFrontendGenerator) GenerateReactComponents(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"package.json":                             `{"name": "` + config.Name + `", "dependencies": {"react": "^18.0.0"}}`,
		"rollup.config.js":                         "export default { input: 'src/index.ts', output: { file: 'dist/index.js' } }",
		"tsconfig.json":                            `{"compilerOptions": {"target": "es5", "lib": ["dom"]}}`,
		"src/index.ts":                             "export * from './components'",
		"src/components/Button/Button.tsx":         "export function Button() { return <button>Button</button> }",
		"src/components/Button/Button.stories.tsx": "export default { title: 'Button' }",
		"src/components/Button/Button.test.tsx":    "test('Button renders', () => { expect(true).toBe(true) })",
		"src/components/Card/Card.tsx":             "export function Card() { return <div>Card</div> }",
		"src/components/Input/Input.tsx":           "export function Input() { return <input /> }",
		"src/hooks/index.ts":                       "export * from './useAuth'",
		"src/utils/index.ts":                       "export * from './helpers'",
		".storybook/main.js":                       "module.exports = { stories: ['../src/**/*.stories.@(js|jsx|ts|tsx)'] }",
		".storybook/preview.js":                    "export const parameters = { actions: { argTypesRegex: '^on[A-Z].*' } }",
	})
}

func (m *MockFrontendGenerator) createMockFiles(basePath string, files map[string]string) error {
	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockMobileGenerator struct{}

func NewMockMobileGenerator() *MockMobileGenerator {
	return &MockMobileGenerator{}
}

func (m *MockMobileGenerator) GenerateAndroidKotlin(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"build.gradle":                                          "// Top-level build file",
		"settings.gradle":                                       "rootProject.name = '" + config.Name + "'",
		"gradle.properties":                                     "android.useAndroidX=true",
		"app/build.gradle":                                      "apply plugin: 'com.android.application'\napply plugin: 'kotlin-android'",
		"app/src/main/AndroidManifest.xml":                      `<?xml version="1.0" encoding="utf-8"?><manifest package="com.example">`,
		"app/src/main/java/com/example/MainActivity.kt":         "class MainActivity : AppCompatActivity()",
		"app/src/main/java/com/example/ui/theme/Theme.kt":       "// Theme definitions",
		"app/src/main/java/com/example/ui/components/Button.kt": "// Button component",
		"app/src/main/res/layout/activity_main.xml":             `<?xml version="1.0" encoding="utf-8"?><LinearLayout />`,
		"app/src/main/res/values/strings.xml":                   `<?xml version="1.0" encoding="utf-8"?><resources><string name="app_name">` + config.Name + `</string></resources>`,
		"app/src/main/res/values/colors.xml":                    `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
		"app/src/test/java/com/example/ExampleUnitTest.kt":      "class ExampleUnitTest { @Test fun addition_isCorrect() { assertEquals(4, 2 + 2) } }",
	})
}

func (m *MockMobileGenerator) GenerateiOSSwift(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"Package.swift":                             "// swift-tools-version:5.5\nimport PackageDescription\n\nlet package = Package(name: \"" + config.Name + "\")",
		"Sources/App/main.swift":                    "import SwiftUI\n\n@main\nstruct " + config.Name + "App: App { var body: some Scene { WindowGroup { ContentView() } } }",
		"Sources/App/ContentView.swift":             "import SwiftUI\n\nstruct ContentView: View { var body: some View { Text(\"Hello, " + config.Name + "!\") } }",
		"Sources/App/Models/User.swift":             "struct User { let id: String; let name: String }",
		"Sources/App/Views/HomeView.swift":          "import SwiftUI\n\nstruct HomeView: View { var body: some View { Text(\"Home\") } }",
		"Sources/App/Views/Components/Button.swift": "import SwiftUI\n\nstruct CustomButton: View { var body: some View { Button(\"Tap\") {} } }",
		"Sources/App/Services/APIService.swift":     "class APIService { func fetchData() async throws -> Data { return Data() } }",
		"Sources/App/Utils/Extensions.swift":        "extension String { var trimmed: String { return self.trimmingCharacters(in: .whitespacesAndNewlines) } }",
		"Tests/AppTests/AppTests.swift":             "import XCTest\n\nfinal class AppTests: XCTestCase { func testExample() { XCTAssertEqual(2 + 2, 4) } }",
		"Resources/Assets.xcassets/Contents.json":   `{"info": {"version": 1, "author": "xcode"}}`,
		"Resources/Info.plist":                      `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"><plist version="1.0"><dict></dict></plist>`,
	})
}

func (m *MockMobileGenerator) GenerateAndroid(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"package.json":                    `{"name": "` + config.Name + `", "dependencies": {"react-native": "^0.72.0"}}`,
		"metro.config.js":                 "module.exports = { transformer: { getTransformOptions: async () => ({ transform: { experimentalImportSupport: false, inlineRequires: true } }) } }",
		"babel.config.js":                 "module.exports = { presets: ['module:metro-react-native-babel-preset'] }",
		"App.tsx":                         "import React from 'react';\nimport { Text, View } from 'react-native';\n\nfunction App() { return <View><Text>Welcome to " + config.Name + "</Text></View> }\n\nexport default App;",
		"index.js":                        "import { AppRegistry } from 'react-native';\nimport App from './App';\n\nAppRegistry.registerComponent('" + config.Name + "', () => App);",
		"src/screens/HomeScreen.tsx":      "import React from 'react';\nimport { View, Text } from 'react-native';\n\nexport function HomeScreen() { return <View><Text>Home</Text></View> }",
		"src/components/Button.tsx":       "import React from 'react';\nimport { TouchableOpacity, Text } from 'react-native';\n\nexport function Button() { return <TouchableOpacity><Text>Button</Text></TouchableOpacity> }",
		"src/navigation/AppNavigator.tsx": "import React from 'react';\nimport { NavigationContainer } from '@react-navigation/native';\n\nexport function AppNavigator() { return <NavigationContainer></NavigationContainer> }",
		"src/services/api.ts":             "export async function fetchData() { return fetch('/api/data') }",
		"src/utils/helpers.ts":            "export function formatDate(date: Date) { return date.toISOString() }",
		"android/build.gradle":            "// Top-level build file for Android",
		"ios/Podfile":                     "platform :ios, '11.0'\nuse_react_native!",
		"__tests__/App-test.tsx":          "import 'react-native';\nimport React from 'react';\nimport App from '../App';\nimport renderer from 'react-test-renderer';\n\ntest('renders correctly', () => { renderer.create(<App />); });",
	})
}

func (m *MockMobileGenerator) createMockFiles(basePath string, files map[string]string) error {
	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockInfrastructureGenerator struct{}

func NewMockInfrastructureGenerator() *MockInfrastructureGenerator {
	return &MockInfrastructureGenerator{}
}

func (m *MockInfrastructureGenerator) GenerateDocker(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"Dockerfile":              "FROM node:18\n\nWORKDIR /app\nCOPY . .\nRUN npm install\nEXPOSE 3000\nCMD [\"npm\", \"start\"]",
		"docker-compose.yml":      "version: '3.8'\nservices:\n  app:\n    build: .\n    ports:\n      - \"3000:3000\"",
		"docker-compose.prod.yml": "version: '3.8'\nservices:\n  app:\n    build: .\n    environment:\n      - NODE_ENV=production",
		"docker-compose.dev.yml":  "version: '3.8'\nservices:\n  app:\n    build: .\n    volumes:\n      - .:/app",
		".dockerignore":           "node_modules\n.git\n.env",
		"scripts/docker-build.sh": "#!/bin/bash\ndocker build -t " + config.Name + " .",
		"scripts/docker-run.sh":   "#!/bin/bash\ndocker run -p 3000:3000 " + config.Name,
		"config/nginx.conf":       "server {\n  listen 80;\n  location / {\n    proxy_pass http://app:3000;\n  }\n}",
		"config/redis.conf":       "bind 127.0.0.1\nport 6379",
	})
}

func (m *MockInfrastructureGenerator) GenerateKubernetes(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"k8s/namespace.yaml":             "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: " + config.Name,
		"k8s/deployment.yaml":            "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: " + config.Name + "\nspec:\n  replicas: 3",
		"k8s/service.yaml":               "apiVersion: v1\nkind: Service\nmetadata:\n  name: " + config.Name + "\nspec:\n  ports:\n  - port: 80",
		"k8s/ingress.yaml":               "apiVersion: networking.k8s.io/v1\nkind: Ingress\nmetadata:\n  name: " + config.Name,
		"k8s/configmap.yaml":             "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: " + config.Name + "-config",
		"k8s/secret.yaml":                "apiVersion: v1\nkind: Secret\nmetadata:\n  name: " + config.Name + "-secret",
		"k8s/hpa.yaml":                   "apiVersion: autoscaling/v2\nkind: HorizontalPodAutoscaler\nmetadata:\n  name: " + config.Name + "-hpa",
		"helm/Chart.yaml":                "apiVersion: v2\nname: " + config.Name + "\nversion: 0.1.0",
		"helm/values.yaml":               "replicaCount: 3\nimage:\n  repository: " + config.Name + "\n  tag: latest",
		"helm/templates/deployment.yaml": "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: {{ .Values.name }}",
		"helm/templates/service.yaml":    "apiVersion: v1\nkind: Service\nmetadata:\n  name: {{ .Values.name }}",
	})
}

func (m *MockInfrastructureGenerator) GenerateTerraform(config *models.ProjectConfig) error {
	return m.createMockFiles(config.OutputPath, map[string]string{
		"main.tf":                   "terraform {\n  required_version = \">= 1.0\"\n}\n\nprovider \"aws\" {\n  region = var.aws_region\n}",
		"variables.tf":              "variable \"aws_region\" {\n  description = \"AWS region\"\n  type        = string\n  default     = \"us-west-2\"\n}",
		"outputs.tf":                "output \"vpc_id\" {\n  description = \"VPC ID\"\n  value       = module.vpc.vpc_id\n}",
		"versions.tf":               "terraform {\n  required_providers {\n    aws = {\n      source  = \"hashicorp/aws\"\n      version = \"~> 5.0\"\n    }\n  }\n}",
		"terraform.tfvars.example":  "aws_region = \"us-west-2\"\nproject_name = \"" + config.Name + "\"",
		"modules/vpc/main.tf":       "resource \"aws_vpc\" \"main\" {\n  cidr_block = \"10.0.0.0/16\"\n}",
		"modules/vpc/variables.tf":  "variable \"cidr_block\" {\n  description = \"CIDR block for VPC\"\n  type        = string\n}",
		"modules/vpc/outputs.tf":    "output \"vpc_id\" {\n  value = aws_vpc.main.id\n}",
		"modules/compute/main.tf":   "resource \"aws_instance\" \"app\" {\n  ami           = \"ami-12345678\"\n  instance_type = \"t3.micro\"\n}",
		"modules/database/main.tf":  "resource \"aws_db_instance\" \"main\" {\n  engine = \"postgres\"\n}",
		"environments/dev/main.tf":  "module \"vpc\" {\n  source = \"../../modules/vpc\"\n}",
		"environments/prod/main.tf": "module \"vpc\" {\n  source = \"../../modules/vpc\"\n}",
	})
}

func (m *MockInfrastructureGenerator) createMockFiles(basePath string, files map[string]string) error {
	for filePath, content := range files {
		fullPath := filepath.Join(basePath, filePath)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockProjectGenerator struct{}

func NewMockProjectGenerator() *MockProjectGenerator {
	return &MockProjectGenerator{}
}

func (m *MockProjectGenerator) GenerateProject(config *models.ProjectConfig) error {
	// Create basic project structure
	dirs := []string{
		"backend",
		"frontend/app",
		"frontend/admin",
		"mobile",
		"infrastructure/docker",
		"infrastructure/k8s",
		"docs",
		"scripts",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(config.OutputPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	// Create root-level files
	rootFiles := map[string]string{
		"README.md":          "# " + config.Name + "\n\n" + config.Description,
		"LICENSE":            config.License + " License\n\nCopyright (c) 2024 " + config.Organization,
		"Makefile":           "build:\n\techo \"Building " + config.Name + "\"\n\ntest:\n\techo \"Testing " + config.Name + "\"",
		"docker-compose.yml": "version: '3.8'\nservices:\n  app:\n    build: .",
		".gitignore":         "node_modules/\n.env\n*.log",
		".env.example":       "NODE_ENV=development\nAPI_URL=http://localhost:3000",
	}

	for fileName, content := range rootFiles {
		filePath := filepath.Join(config.OutputPath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// Mock Generators
type MockStructureGenerator struct{}

func NewMockStructureGenerator() *MockStructureGenerator {
	return &MockStructureGenerator{}
}

func (m *MockStructureGenerator) GenerateStructure(config *models.ProjectConfig) error {
	dirs := []string{"src", "tests", "docs", "scripts"}
	for _, dir := range dirs {
		dirPath := filepath.Join(config.OutputPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	return nil
}

type MockTemplateGenerator struct{}

func NewMockTemplateGenerator() *MockTemplateGenerator {
	return &MockTemplateGenerator{}
}

func (m *MockTemplateGenerator) GenerateTemplates(config *models.ProjectConfig) error {
	files := map[string]string{
		"README.md":  "# " + config.Name,
		"LICENSE":    config.License + " License",
		".gitignore": "node_modules/\n.env",
	}

	for fileName, content := range files {
		filePath := filepath.Join(config.OutputPath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockConfigurationGenerator struct{}

func NewMockConfigurationGenerator() *MockConfigurationGenerator {
	return &MockConfigurationGenerator{}
}

func (m *MockConfigurationGenerator) GenerateConfiguration(config *models.ProjectConfig) error {
	// Create config directory
	configDir := filepath.Join(config.OutputPath, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	files := map[string]string{
		"config/app.yaml":      "app:\n  name: " + config.Name,
		"config/database.yaml": "database:\n  host: localhost",
		".env.example":         "NODE_ENV=development",
	}

	for fileName, content := range files {
		filePath := filepath.Join(config.OutputPath, fileName)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockDocumentationGenerator struct{}

func NewMockDocumentationGenerator() *MockDocumentationGenerator {
	return &MockDocumentationGenerator{}
}

func (m *MockDocumentationGenerator) GenerateDocumentation(config *models.ProjectConfig) error {
	// Create docs directory
	docsDir := filepath.Join(config.OutputPath, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	files := map[string]string{
		"docs/README.md":       "# " + config.Name + " Documentation",
		"docs/API.md":          "# API Documentation\n\nAPI endpoints for " + config.Name,
		"docs/CONTRIBUTING.md": "# Contributing Guide",
		"docs/DEPLOYMENT.md":   "# Deployment Guide",
	}

	for fileName, content := range files {
		filePath := filepath.Join(config.OutputPath, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

type MockCICDGenerator struct{}

func NewMockCICDGenerator() *MockCICDGenerator {
	return &MockCICDGenerator{}
}

func (m *MockCICDGenerator) GenerateCICD(config *models.ProjectConfig) error {
	// Create .github/workflows directory
	workflowsDir := filepath.Join(config.OutputPath, ".github/workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return err
	}

	files := map[string]string{
		".github/workflows/ci.yml":   "name: CI\non:\n  push:\n    branches: [main]",
		".github/workflows/cd.yml":   "name: CD\non:\n  push:\n    tags: ['v*']",
		".github/workflows/test.yml": "name: Test\non:\n  pull_request:",
		".gitlab-ci.yml":             "stages:\n  - test\n  - build\n  - deploy",
		"Jenkinsfile":                "pipeline {\n  agent any\n  stages {\n    stage('Build') {\n      steps {\n        echo 'Building...'\n      }\n    }\n  }\n}",
	}

	for fileName, content := range files {
		filePath := filepath.Join(config.OutputPath, fileName)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}
