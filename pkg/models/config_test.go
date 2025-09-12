package models

import (
	"encoding/json"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v3"
)

func TestProjectConfigValidation(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name        string
		config      *ProjectConfig
		expectValid bool
		expectError string
	}{
		{
			name: "valid configuration",
			config: &ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "A test project for validation",
				License:      "MIT",
				Author:       "Test Author",
				Email:        "test@example.com",
				Repository:   "https://github.com/test-org/test-project",
				Components: Components{
					Frontend: FrontendComponents{MainApp: true},
					Backend:  BackendComponents{API: true},
					Mobile:   MobileComponents{Android: false, IOS: false},
					Infrastructure: InfrastructureComponents{
						Docker:     true,
						Kubernetes: false,
						Terraform:  false,
					},
				},
				Versions: &VersionConfig{
					Node:   "20.0.0",
					Go:     "1.22.0",
					NextJS: "15.0.0",
					React:  "18.0.0",
				},
				OutputPath:       "./output",
				GeneratedAt:      time.Now(),
				GeneratorVersion: "1.0.0",
			},
			expectValid: true,
		},
		{
			name: "missing required fields",
			config: &ProjectConfig{
				Name: "", // Missing required field
				// Missing other required fields
			},
			expectValid: false,
			expectError: "Name is required",
		},
		{
			name: "invalid email",
			config: &ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "A test project for validation",
				License:      "MIT",
				Email:        "invalid-email", // Invalid email
				Components: Components{
					Frontend: FrontendComponents{MainApp: true},
					Backend:  BackendComponents{},
					Mobile:   MobileComponents{},
					Infrastructure: InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			},
			expectValid: false,
			expectError: "Email must be a valid email address",
		},
		{
			name: "invalid license",
			config: &ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "A test project for validation",
				License:      "INVALID", // Invalid license
				Components: Components{
					Frontend: FrontendComponents{MainApp: true},
					Backend:  BackendComponents{},
					Mobile:   MobileComponents{},
					Infrastructure: InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			},
			expectValid: false,
			expectError: "License must be one of",
		},
		{
			name: "no components selected",
			config: &ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
				Description:  "A test project for validation",
				License:      "MIT",
				Components: Components{
					Frontend:       FrontendComponents{},
					Backend:        BackendComponents{},
					Mobile:         MobileComponents{},
					Infrastructure: InfrastructureComponents{},
				},
				Versions: &VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			},
			expectValid: false,
			expectError: "At least one component must be selected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateProjectConfig(tt.config)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if !tt.expectValid && tt.expectError != "" {
				found := false
				for _, err := range result.Errors {
					if err.Message == tt.expectError ||
						(len(tt.expectError) > 0 && len(err.Message) > 0 &&
							err.Message[:len(tt.expectError)] == tt.expectError) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s', got errors: %v", tt.expectError, result.Errors)
				}
			}
		})
	}
}

func TestVersionConfigValidation(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name        string
		config      *VersionConfig
		expectValid bool
		expectError string
	}{
		{
			name: "valid version config",
			config: &VersionConfig{
				Node:   "20.0.0",
				Go:     "1.22.0",
				Kotlin: "2.0.0",
				Swift:  "5.9.0",
				NextJS: "15.0.0",
				React:  "18.0.0",
				Packages: map[string]string{
					"express": "4.18.0",
					"lodash":  "4.17.21",
				},
				UpdatedAt: time.Now(),
			},
			expectValid: true,
		},
		{
			name: "invalid semantic versions",
			config: &VersionConfig{
				Node: "invalid-version", // Invalid semver
				Go:   "1.22.0",
			},
			expectValid: false,
			expectError: "Node must be a valid semantic version",
		},
		{
			name:   "missing required versions",
			config: &VersionConfig{
				// Missing required Node and Go versions
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateVersionConfig(tt.config)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if !tt.expectValid && tt.expectError != "" {
				found := false
				for _, err := range result.Errors {
					if err.Message == tt.expectError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error '%s', got errors: %v", tt.expectError, result.Errors)
				}
			}
		})
	}
}

func TestTemplateMetadataValidation(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name        string
		metadata    *TemplateMetadata
		expectValid bool
		expectError string
	}{
		{
			name: "valid template metadata",
			metadata: &TemplateMetadata{
				Name:         "test-template",
				Description:  "A test template for validation",
				Version:      "1.0.0",
				Author:       "Test Author",
				Dependencies: []string{"dep1", "dep2"},
				Variables: []TemplateVar{
					{
						Name:        "project_name",
						Type:        "string",
						Default:     "my-project",
						Description: "The project name",
						Required:    true,
					},
				},
				Conditions: []TemplateCondition{
					{
						Name:      "has_frontend",
						Component: "frontend.main_app",
						Operator:  "eq",
						Value:     true,
					},
				},
				Tags:                []string{"web", "api"},
				MinGeneratorVersion: "1.0.0",
			},
			expectValid: true,
		},
		{
			name:     "missing required fields",
			metadata: &TemplateMetadata{
				// Missing required fields
			},
			expectValid: false,
		},
		{
			name: "invalid variable type",
			metadata: &TemplateMetadata{
				Name:        "test-template",
				Description: "A test template",
				Version:     "1.0.0",
				Variables: []TemplateVar{
					{
						Name:        "test_var",
						Type:        "invalid_type", // Invalid type
						Description: "A test variable",
						Required:    true,
					},
				},
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateMetadata(tt.metadata)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
				if !result.Valid {
					t.Logf("Validation errors: %v", result.Errors)
				}
			}
		})
	}
}

func TestTemplateVarValidation(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name        string
		variable    *TemplateVar
		expectValid bool
		expectError string
	}{
		{
			name: "valid string variable",
			variable: &TemplateVar{
				Name:        "project_name",
				Type:        "string",
				Default:     "my-project",
				Description: "The project name",
				Required:    true,
			},
			expectValid: true,
		},
		{
			name: "valid int variable",
			variable: &TemplateVar{
				Name:        "port",
				Type:        "int",
				Default:     8080,
				Description: "The server port",
				Required:    false,
			},
			expectValid: true,
		},
		{
			name: "type mismatch",
			variable: &TemplateVar{
				Name:        "port",
				Type:        "int",
				Default:     "8080", // String instead of int
				Description: "The server port",
				Required:    false,
			},
			expectValid: false,
			expectError: "Default value type does not match declared type 'int'",
		},
		{
			name: "invalid variable name",
			variable: &TemplateVar{
				Name:        "invalid name!", // Invalid characters
				Type:        "string",
				Description: "A test variable",
				Required:    true,
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateTemplateVar(tt.variable)

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
				if !result.Valid {
					t.Logf("Validation errors: %v", result.Errors)
				}
			}

			if !tt.expectValid && tt.expectError != "" {
				found := false
				for _, err := range result.Errors {
					if err.Message == tt.expectError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error '%s', got errors: %v", tt.expectError, result.Errors)
				}
			}
		})
	}
}

func TestProjectConfigSerialization(t *testing.T) {
	config := &ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "A test project for serialization",
		License:      "MIT",
		Author:       "Test Author",
		Email:        "test@example.com",
		Repository:   "https://github.com/test-org/test-project",
		Components: Components{
			Frontend: FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   true,
			},
			Backend: BackendComponents{
				API: true,
			},
			Mobile: MobileComponents{
				Android: true,
				IOS:     false,
			},
			Infrastructure: InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  false,
			},
		},
		Versions: &VersionConfig{
			Node:   "20.0.0",
			Go:     "1.22.0",
			Kotlin: "2.0.0",
			NextJS: "15.0.0",
			React:  "18.0.0",
			Packages: map[string]string{
				"express": "4.18.0",
				"lodash":  "4.17.21",
			},
			UpdatedAt: time.Now(),
		},
		CustomVars: map[string]string{
			"custom_var1": "value1",
			"custom_var2": "value2",
		},
		OutputPath:       "./output",
		GeneratedAt:      time.Now(),
		GeneratorVersion: "1.0.0",
	}

	t.Run("YAML serialization", func(t *testing.T) {
		// Test YAML marshaling
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			t.Fatalf("Failed to marshal to YAML: %v", err)
		}

		// Test YAML unmarshaling
		var unmarshaledConfig ProjectConfig
		err = yaml.Unmarshal(yamlData, &unmarshaledConfig)
		if err != nil {
			t.Fatalf("Failed to unmarshal from YAML: %v", err)
		}

		// Verify key fields
		if unmarshaledConfig.Name != config.Name {
			t.Errorf("Expected name %s, got %s", config.Name, unmarshaledConfig.Name)
		}
		if unmarshaledConfig.Organization != config.Organization {
			t.Errorf("Expected organization %s, got %s", config.Organization, unmarshaledConfig.Organization)
		}
		if unmarshaledConfig.Components.Frontend.MainApp != config.Components.Frontend.MainApp {
			t.Errorf("Expected MainApp %v, got %v", config.Components.Frontend.MainApp, unmarshaledConfig.Components.Frontend.MainApp)
		}
	})

	t.Run("JSON serialization", func(t *testing.T) {
		// Test JSON marshaling
		jsonData, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Failed to marshal to JSON: %v", err)
		}

		// Test JSON unmarshaling
		var unmarshaledConfig ProjectConfig
		err = json.Unmarshal(jsonData, &unmarshaledConfig)
		if err != nil {
			t.Fatalf("Failed to unmarshal from JSON: %v", err)
		}

		// Verify key fields
		if unmarshaledConfig.Name != config.Name {
			t.Errorf("Expected name %s, got %s", config.Name, unmarshaledConfig.Name)
		}
		if unmarshaledConfig.Organization != config.Organization {
			t.Errorf("Expected organization %s, got %s", config.Organization, unmarshaledConfig.Organization)
		}
		if unmarshaledConfig.Components.Backend.API != config.Components.Backend.API {
			t.Errorf("Expected API %v, got %v", config.Components.Backend.API, unmarshaledConfig.Components.Backend.API)
		}
	})
}

func TestCustomValidationFunctions(t *testing.T) {
	validator := NewConfigValidator()

	t.Run("semver validation", func(t *testing.T) {
		validVersions := []string{
			"1.0.0",
			"v1.0.0",
			"1.0.0-alpha",
			"1.0.0-alpha.1",
			"1.0.0+build.1",
			"1.0.0-alpha+build.1",
		}

		invalidVersions := []string{
			"1.0",
			"1.0.0.0",
			"invalid",
			"1.0.0-",
			"1.0.0+",
		}

		for _, version := range validVersions {
			config := &VersionConfig{
				Node: version,
				Go:   "1.22.0",
			}
			result := validator.ValidateVersionConfig(config)
			if !result.Valid {
				t.Errorf("Expected version %s to be valid, got errors: %v", version, result.Errors)
			}
		}

		for _, version := range invalidVersions {
			config := &VersionConfig{
				Node: version,
				Go:   "1.22.0",
			}
			result := validator.ValidateVersionConfig(config)
			if result.Valid {
				t.Errorf("Expected version %s to be invalid", version)
			}
		}
	})

	t.Run("alphanum validation", func(t *testing.T) {
		validNames := []string{
			"test-project",
			"test_project",
			"testproject",
			"test123",
			"123test",
		}

		invalidNames := []string{
			"test project", // space
			"test.project", // dot
			"test@project", // special character
			"test/project", // slash
		}

		for _, name := range validNames {
			config := &ProjectConfig{
				Name:         name,
				Organization: "test-org",
				Description:  "A test project",
				License:      "MIT",
				Components: Components{
					Frontend: FrontendComponents{MainApp: true},
					Backend:  BackendComponents{},
					Mobile:   MobileComponents{},
					Infrastructure: InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			}
			result := validator.ValidateProjectConfig(config)
			if !result.Valid {
				t.Errorf("Expected name %s to be valid, got errors: %v", name, result.Errors)
			}
		}

		for _, name := range invalidNames {
			config := &ProjectConfig{
				Name:         name,
				Organization: "test-org",
				Description:  "A test project",
				License:      "MIT",
				Components: Components{
					Frontend: FrontendComponents{MainApp: true},
					Backend:  BackendComponents{},
					Mobile:   MobileComponents{},
					Infrastructure: InfrastructureComponents{
						Docker: true,
					},
				},
				Versions: &VersionConfig{
					Node: "20.0.0",
					Go:   "1.22.0",
				},
				OutputPath: "./output",
			}
			result := validator.ValidateProjectConfig(config)
			if result.Valid {
				t.Errorf("Expected name %s to be invalid", name)
			}
		}
	})
}
