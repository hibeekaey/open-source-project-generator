package constants

import (
	"testing"
)

func TestPackageManagerConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"NPM package manager", PackageManagerNPM, "npm"},
		{"Yarn package manager", PackageManagerYarn, "yarn"},
		{"PNPM package manager", PackageManagerPNPM, "pnpm"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestLanguageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"JavaScript language", LanguageJavaScript, "javascript"},
		{"TypeScript language", LanguageTypeScript, "typescript"},
		{"Node.js runtime", LanguageNodeJS, "nodejs"},
		{"Go language", LanguageGo, "go"},
		{"Python language", LanguagePython, "python"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestVersionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Latest version", VersionLatest, "latest"},
		{"Unknown version", VersionUnknown, "unknown"},
		{"Development version", VersionDev, "dev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestSeverityConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Critical severity", SeverityCritical, "critical"},
		{"High severity", SeverityHigh, "high"},
		{"Medium severity", SeverityMedium, "medium"},
		{"Low severity", SeverityLow, "low"},
		{"Info severity", SeverityInfo, "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Success status", StatusSuccess, "success"},
		{"Failed status", StatusFailed, "failed"},
		{"Pending status", StatusPending, "pending"},
		{"Skipped status", StatusSkipped, "skipped"},
		{"In progress status", StatusInProgress, "in_progress"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestFileFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"JSON format", FormatJSON, "json"},
		{"YAML format", FormatYAML, "yaml"},
		{"YML format", FormatYML, "yml"},
		{"TOML format", FormatTOML, "toml"},
		{"XML format", FormatXML, "xml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestFileTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Package file type", FileTypePackage, "package"},
		{"Config file type", FileTypeConfig, "config"},
		{"Template file type", FileTypeTemplate, "template"},
		{"Documentation file type", FileTypeDocumentation, "documentation"},
		{"Script file type", FileTypeScript, "script"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestTemplateTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Frontend template", TemplateFrontend, "frontend"},
		{"Backend template", TemplateBackend, "backend"},
		{"Mobile template", TemplateMobile, "mobile"},
		{"Infrastructure template", TemplateInfrastructure, "infrastructure"},
		{"Base template", TemplateBase, "base"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestUISymbolConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Success symbol", SymbolSuccess, "✅"},
		{"Failure symbol", SymbolFailure, "❌"},
		{"Warning symbol", SymbolWarning, "⚠️"},
		{"Info symbol", SymbolInfo, "ℹ️"},
		{"Check symbol", SymbolCheck, "✓"},
		{"Cross symbol", SymbolCross, "✗"},
		{"Bullet symbol", SymbolBullet, "•"},
		{"Arrow symbol", SymbolArrow, "→"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestTemplatePathConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Template base dir", TemplateBaseDir, "templates"},
		{"Embedded template dir", EmbeddedTemplateDir, "templates"},
		{"Embedded template base path", EmbeddedTemplateBasePath, "pkg/template/templates"},
		{"Template relative base path", TemplateRelativeBasePath, "templates"},
		{"Template config dir", TemplateConfigDir, "config"},
		{"Template defaults file", TemplateDefaultsFile, "defaults.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestValidationConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Consistent validation", ValidationConsistent, "consistent"},
		{"Inconsistent validation", ValidationInconsistent, "inconsistent"},
		{"Required validation", ValidationRequired, "required"},
		{"Optional validation", ValidationOptional, "optional"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestStringTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"String type", StringType, "string"},
		{"Number type", NumberType, "number"},
		{"Boolean type", BooleanType, "boolean"},
		{"Object type", ObjectType, "object"},
		{"Array type", ArrayType, "array"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestAdditionalConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"High priority", PriorityHigh, "high"},
		{"Critical priority", PriorityCritical, "critical"},
		{"Present status", StatusPresent, "present"},
		{"Passed status", StatusPassed, "passed"},
		{"Package.json file", FilePackageJSON, "package.json"},
		{"Go.mod file", FileGoMod, "go.mod"},
		{"Dockerfile", FileDockerfile, "Dockerfile"},
		{"Language type", TypeLanguage, "language"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestMessageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Failed message", MessageFailed, "failed"},
		{"Success message", MessageSuccess, "success"},
		{"Processing message", MessageProcessing, "processing"},
		{"Completed message", MessageCompleted, "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestBuildModeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Development build mode", BuildModeDevelopment, "development"},
		{"Production build mode", BuildModeProduction, "production"},
		{"Test build mode", BuildModeTest, "test"},
		{"Staging build mode", BuildModeStaging, "staging"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}

// TestConstantsUniqueness ensures that constants within each category are unique
func TestConstantsUniqueness(t *testing.T) {
	// Test package manager constants uniqueness
	packageManagers := []string{PackageManagerNPM, PackageManagerYarn, PackageManagerPNPM}
	if !areUnique(packageManagers) {
		t.Error("Package manager constants are not unique")
	}

	// Test severity constants uniqueness
	severities := []string{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	if !areUnique(severities) {
		t.Error("Severity constants are not unique")
	}

	// Test status constants uniqueness
	statuses := []string{StatusSuccess, StatusFailed, StatusPending, StatusSkipped, StatusInProgress}
	if !areUnique(statuses) {
		t.Error("Status constants are not unique")
	}

	// Test file format constants uniqueness
	formats := []string{FormatJSON, FormatYAML, FormatYML, FormatTOML, FormatXML}
	if !areUnique(formats) {
		t.Error("File format constants are not unique")
	}
}

// Helper function to check if all strings in a slice are unique
func areUnique(strings []string) bool {
	seen := make(map[string]bool)
	for _, str := range strings {
		if seen[str] {
			return false
		}
		seen[str] = true
	}
	return true
}

// TestConstantsNotEmpty ensures that no constants are empty strings
func TestConstantsNotEmpty(t *testing.T) {
	constants := []struct {
		name  string
		value string
	}{
		{"PackageManagerNPM", PackageManagerNPM},
		{"LanguageJavaScript", LanguageJavaScript},
		{"VersionLatest", VersionLatest},
		{"SeverityCritical", SeverityCritical},
		{"StatusSuccess", StatusSuccess},
		{"FormatJSON", FormatJSON},
		{"FileTypePackage", FileTypePackage},
		{"TemplateFrontend", TemplateFrontend},
		{"SymbolSuccess", SymbolSuccess},
		{"TemplateBaseDir", TemplateBaseDir},
		{"ValidationConsistent", ValidationConsistent},
		{"StringType", StringType},
		{"MessageSuccess", MessageSuccess},
		{"BuildModeDevelopment", BuildModeDevelopment},
	}

	for _, constant := range constants {
		t.Run(constant.name, func(t *testing.T) {
			if constant.value == "" {
				t.Errorf("Constant %s should not be empty", constant.name)
			}
		})
	}
}
