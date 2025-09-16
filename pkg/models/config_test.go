package models

import (
	"encoding/json"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestProjectConfigSerialization(t *testing.T) {
	config := &ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project description",
		Email:        "test@example.com",
		License:      "MIT",
		OutputPath:   "./test-output",
		Components: Components{
			Frontend: FrontendComponents{
				NextJS: NextJSComponents{
					App: true,
				},
			},
			Backend: BackendComponents{
				GoGin: true,
			},
		},
		Versions: &VersionConfig{
			Node: "20.11.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"next":  "14.0.0",
				"react": "18.2.0",
			},
		},
	}

	// Test YAML serialization
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("YAML serialization failed: %v", err)
	}

	var yamlConfig ProjectConfig
	err = yaml.Unmarshal(yamlData, &yamlConfig)
	if err != nil {
		t.Fatalf("YAML deserialization failed: %v", err)
	}

	if yamlConfig.Name != config.Name {
		t.Errorf("YAML serialization: expected name %s, got %s", config.Name, yamlConfig.Name)
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var jsonConfig ProjectConfig
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if jsonConfig.Name != config.Name {
		t.Errorf("JSON serialization: expected name %s, got %s", config.Name, jsonConfig.Name)
	}
}

func TestVersionConfigSerialization(t *testing.T) {
	versionConfig := &VersionConfig{
		Node: "20.11.0",
		Go:   "1.21.0",
		Packages: map[string]string{
			"next":       "14.0.0",
			"react":      "18.2.0",
			"typescript": "^5.0.0",
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(versionConfig)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var jsonConfig VersionConfig
	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if jsonConfig.Node != versionConfig.Node {
		t.Errorf("Expected Node %s, got %s", versionConfig.Node, jsonConfig.Node)
	}

	if jsonConfig.Go != versionConfig.Go {
		t.Errorf("Expected Go %s, got %s", versionConfig.Go, jsonConfig.Go)
	}

	if len(jsonConfig.Packages) != len(versionConfig.Packages) {
		t.Errorf("Expected %d packages, got %d", len(versionConfig.Packages), len(jsonConfig.Packages))
	}
}

func TestTemplateMetadataSerialization(t *testing.T) {
	metadata := &TemplateMetadata{
		Name:        "test-template",
		Description: "Test template description",
		Version:     "1.0.0",
		Author:      "Test Author",
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}

	var jsonMetadata TemplateMetadata
	err = json.Unmarshal(jsonData, &jsonMetadata)
	if err != nil {
		t.Fatalf("JSON deserialization failed: %v", err)
	}

	if jsonMetadata.Name != metadata.Name {
		t.Errorf("Expected name %s, got %s", metadata.Name, jsonMetadata.Name)
	}
}
