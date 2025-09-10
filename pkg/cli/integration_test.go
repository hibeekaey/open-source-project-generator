package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

// TestCLIWorkflow tests the complete CLI workflow
func TestCLIWorkflow(t *testing.T) {
	// Create mock dependencies
	configManager := &MockConfigManager{}
	validator := &MockValidationEngine{}

	cli := NewCLI(configManager, validator)

	// Test component dependency validation
	t.Run("ValidateComponentDependencies", func(t *testing.T) {
		testCases := []struct {
			name        string
			components  []string
			expectError bool
		}{
			{
				name: "valid frontend and backend",
				components: []string{
					"frontend.main_app - Main Next.js application",
					"backend.api - Go API server with Gin framework",
				},
				expectError: false,
			},
			{
				name: "mobile without backend (warning but allowed)",
				components: []string{
					"mobile.android - Android Kotlin application",
					"infrastructure.docker - Docker configurations",
				},
				expectError: false, // Should show warning but not error
			},
			{
				name: "no main components",
				components: []string{
					"infrastructure.docker - Docker configurations",
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Use the non-prompting version for tests
				err := cli.validateComponentDependenciesWithPrompt(tc.components, false)
				if tc.expectError && err == nil {
					t.Error("Expected error but got none")
				}
				if !tc.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			})
		}
	})

	// Test configuration preview
	t.Run("PreviewConfiguration", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "test-project",
			Organization: "test-org",
			Components: models.Components{
				Frontend: models.FrontendComponents{
					MainApp: true,
					Admin:   true,
				},
				Backend: models.BackendComponents{
					API: true,
				},
				Infrastructure: models.InfrastructureComponents{
					Docker:     true,
					Kubernetes: true,
				},
			},
		}

		// This should not panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PreviewConfiguration panicked: %v", r)
			}
		}()

		cli.PreviewConfiguration(config)
	})

	// Test output path validation
	t.Run("CheckOutputPath", func(t *testing.T) {
		// Test with temporary directory
		tempDir := t.TempDir()

		// Empty directory should be fine
		err := cli.CheckOutputPath(tempDir)
		if err != nil {
			t.Errorf("CheckOutputPath failed for empty directory: %v", err)
		}

		// Create a file in the directory
		testFile := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Non-empty directory should prompt for confirmation
		// In test environment, this will fail due to no interactive input
		err = cli.CheckOutputPath(tempDir)
		if err == nil {
			t.Error("Expected error for non-empty directory in non-interactive mode")
		}
	})
}

// TestComponentSelection tests the component selection logic
func TestComponentSelection(t *testing.T) {
	cli := &CLI{}

	testCases := []struct {
		name       string
		components []string
		expected   models.Components
	}{
		{
			name: "full stack selection",
			components: []string{
				"frontend.main_app - Main Next.js application",
				"frontend.admin - Admin dashboard application",
				"backend.api - Go API server with Gin framework",
				"mobile.android - Android Kotlin application",
				"mobile.ios - iOS Swift application",
				"infrastructure.docker - Docker configurations",
				"infrastructure.kubernetes - Kubernetes manifests",
				"infrastructure.terraform - Terraform configurations",
			},
			expected: models.Components{
				Frontend: models.FrontendComponents{
					MainApp: true,
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
					Docker:     true,
					Kubernetes: true,
					Terraform:  true,
				},
			},
		},
		{
			name: "minimal selection",
			components: []string{
				"frontend.main_app - Main Next.js application",
			},
			expected: models.Components{
				Frontend: models.FrontendComponents{
					MainApp: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &models.ProjectConfig{}

			err := cli.setSelectedComponents(config, tc.components)
			if err != nil {
				t.Fatalf("setSelectedComponents failed: %v", err)
			}

			// Verify all components are set correctly
			if config.Components.Frontend.MainApp != tc.expected.Frontend.MainApp {
				t.Errorf("Frontend.MainApp: expected %v, got %v", tc.expected.Frontend.MainApp, config.Components.Frontend.MainApp)
			}
			if config.Components.Frontend.Home != tc.expected.Frontend.Home {
				t.Errorf("Frontend.Home: expected %v, got %v", tc.expected.Frontend.Home, config.Components.Frontend.Home)
			}
			if config.Components.Frontend.Admin != tc.expected.Frontend.Admin {
				t.Errorf("Frontend.Admin: expected %v, got %v", tc.expected.Frontend.Admin, config.Components.Frontend.Admin)
			}
			if config.Components.Backend.API != tc.expected.Backend.API {
				t.Errorf("Backend.API: expected %v, got %v", tc.expected.Backend.API, config.Components.Backend.API)
			}
			if config.Components.Mobile.Android != tc.expected.Mobile.Android {
				t.Errorf("Mobile.Android: expected %v, got %v", tc.expected.Mobile.Android, config.Components.Mobile.Android)
			}
			if config.Components.Mobile.IOS != tc.expected.Mobile.IOS {
				t.Errorf("Mobile.IOS: expected %v, got %v", tc.expected.Mobile.IOS, config.Components.Mobile.IOS)
			}
			if config.Components.Infrastructure.Docker != tc.expected.Infrastructure.Docker {
				t.Errorf("Infrastructure.Docker: expected %v, got %v", tc.expected.Infrastructure.Docker, config.Components.Infrastructure.Docker)
			}
			if config.Components.Infrastructure.Kubernetes != tc.expected.Infrastructure.Kubernetes {
				t.Errorf("Infrastructure.Kubernetes: expected %v, got %v", tc.expected.Infrastructure.Kubernetes, config.Components.Infrastructure.Kubernetes)
			}
			if config.Components.Infrastructure.Terraform != tc.expected.Infrastructure.Terraform {
				t.Errorf("Infrastructure.Terraform: expected %v, got %v", tc.expected.Infrastructure.Terraform, config.Components.Infrastructure.Terraform)
			}
		})
	}
}

// TestDefaultVersions tests the default version configuration
func TestDefaultVersions(t *testing.T) {
	cli := &CLI{}
	versions := cli.getDefaultVersions()

	// Verify all required versions are present
	requiredVersions := map[string]string{
		"Node":   versions.Node,
		"Go":     versions.Go,
		"Kotlin": versions.Kotlin,
		"Swift":  versions.Swift,
		"NextJS": versions.NextJS,
		"React":  versions.React,
	}

	for name, version := range requiredVersions {
		if version == "" {
			t.Errorf("%s version should not be empty", name)
		}
	}

	// Verify UpdatedAt is recent
	if time.Since(versions.UpdatedAt) > time.Minute {
		t.Error("UpdatedAt should be recent")
	}

	// Verify Packages map is initialized
	if versions.Packages == nil {
		t.Error("Packages map should be initialized")
	}
}

// TestProgressMethods tests the progress indication methods
func TestProgressMethods(t *testing.T) {
	cli := &CLI{}

	// These methods should not panic and should handle various inputs
	testMessages := []string{
		"Simple message",
		"Message with special characters: !@#$%^&*()",
		"",
		"Very long message that might wrap around the terminal and should still be handled gracefully without causing any issues",
	}

	for _, msg := range testMessages {
		t.Run("ShowProgress", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShowProgress panicked with message '%s': %v", msg, r)
				}
			}()
			cli.ShowProgress(msg)
		})

		t.Run("ShowSuccess", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShowSuccess panicked with message '%s': %v", msg, r)
				}
			}()
			cli.ShowSuccess(msg)
		})

		t.Run("ShowError", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShowError panicked with message '%s': %v", msg, r)
				}
			}()
			cli.ShowError(msg)
		})

		t.Run("ShowWarning", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShowWarning panicked with message '%s': %v", msg, r)
				}
			}()
			cli.ShowWarning(msg)
		})
	}
}

// BenchmarkComponentSelection benchmarks the component selection performance
func BenchmarkComponentSelection(b *testing.B) {
	cli := &CLI{}
	config := &models.ProjectConfig{}
	components := []string{
		"frontend.main_app - Main Next.js application",
		"frontend.home - Landing page application",
		"frontend.admin - Admin dashboard application",
		"backend.api - Go API server with Gin framework",
		"mobile.android - Android Kotlin application",
		"mobile.ios - iOS Swift application",
		"infrastructure.docker - Docker configurations",
		"infrastructure.kubernetes - Kubernetes manifests",
		"infrastructure.terraform - Terraform configurations",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Components = models.Components{} // Reset
		err := cli.setSelectedComponents(config, components)
		if err != nil {
			b.Fatalf("setSelectedComponents failed: %v", err)
		}
	}
}
