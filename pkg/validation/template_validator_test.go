package validation

import (
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/template"
)

func TestNewTemplateValidator(t *testing.T) {
	validator := NewTemplateValidator()

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}

	if validator.standardConfigs == nil {
		t.Error("Expected standardConfigs to be initialized")
	}

	// Check if embedded filesystem is available
	if validator.embeddedFS == nil {
		t.Error("Expected embedded filesystem to be available")
	}

	if !validator.useEmbedded {
		t.Error("Expected useEmbedded to be true when embedded filesystem is available")
	}
}

func TestNewTemplateValidatorWithFS(t *testing.T) {
	embeddedFS := template.GetEmbeddedFS()
	validator := NewTemplateValidatorWithFS(embeddedFS)

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}

	if validator.embeddedFS != embeddedFS {
		t.Error("Expected embedded filesystem to match provided filesystem")
	}

	if !validator.useEmbedded {
		t.Error("Expected useEmbedded to be true when filesystem is provided")
	}
}

func TestValidateTemplateConsistency_EmbeddedMode(t *testing.T) {
	validator := NewTemplateValidator()

	// Test with embedded filesystem
	result, err := validator.ValidateTemplateConsistency("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Should have some validation results
	if result.Summary == "" {
		t.Error("Expected validation summary to be set")
	}
}

func TestValidateEmbeddedTemplateStructure(t *testing.T) {
	validator := NewTemplateValidator()

	result, err := validator.ValidateEmbeddedTemplateStructure()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Should complete validation
	if result.Summary == "" {
		t.Error("Expected validation summary to be set")
	}
}

func TestValidateAllEmbeddedTemplates(t *testing.T) {
	validator := NewTemplateValidator()

	result, err := validator.ValidateAllEmbeddedTemplates()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	// Should complete comprehensive validation
	if result.Summary == "" {
		t.Error("Expected validation summary to be set")
	}
}

func TestGetValidationSummary(t *testing.T) {
	validator := NewTemplateValidator()

	summary := validator.GetValidationSummary()

	if summary == nil {
		t.Fatal("Expected validation summary, got nil")
	}

	// Check required fields
	if _, exists := summary["embedded_filesystem_available"]; !exists {
		t.Error("Expected embedded_filesystem_available field in summary")
	}

	if _, exists := summary["validation_methods"]; !exists {
		t.Error("Expected validation_methods field in summary")
	}
}

func TestValidateTemplateSyntax(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid template",
			content: "Hello {{.Name}}!",
			wantErr: false,
		},
		{
			name:    "unclosed range",
			content: "{{range .Items}} item {{.Name}}",
			wantErr: true,
		},
		{
			name:    "suspicious spacing",
			content: "{{. Name}}",
			wantErr: true,
		},
		{
			name:    "valid range",
			content: "{{range .Items}} item {{.Name}} {{end}}",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateTemplateSyntax(tt.content, "test.tmpl")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTemplateSyntax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTemplateSecurity(t *testing.T) {
	validator := NewTemplateValidator()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "safe template",
			content: "Hello {{.Name}}!",
			wantErr: false,
		},
		{
			name:    "dangerous env access",
			content: "{{.Env.SECRET}}",
			wantErr: true,
		},
		{
			name:    "hardcoded password",
			content: "password=secret123",
			wantErr: true,
		},
		{
			name:    "system access",
			content: "{{.System.Command}}",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateTemplateSecurity(tt.content, "test.tmpl")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTemplateSecurity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidSemver(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"1.0.0", true},
		{"1.2.3", true},
		{"1.0.0-alpha", true},
		{"1.0.0+build", true},
		{"1.0", false},
		{"1", false},
		{"", false},
		{"1.0.0.0", false},
		{"a.b.c", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := isValidSemver(tt.version); got != tt.want {
				t.Errorf("isValidSemver(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
