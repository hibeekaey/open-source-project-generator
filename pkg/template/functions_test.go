package template

import (
	"testing"
	"time"

	"github.com/open-source-template-generator/pkg/models"
)

func TestStringManipulationFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(string) string
		input    string
		expected string
	}{
		{"toCamelCase", toCamelCase, "hello-world", "helloWorld"},
		{"toCamelCase", toCamelCase, "hello_world", "helloWorld"},
		{"toCamelCase", toCamelCase, "HelloWorld", "helloWorld"},
		{"toCamelCase", toCamelCase, "HELLO_WORLD", "helloWorld"},
		{"toSnakeCase", toSnakeCase, "HelloWorld", "hello_world"},
		{"toSnakeCase", toSnakeCase, "helloWorld", "hello_world"},
		{"toSnakeCase", toSnakeCase, "hello-world", "hello-world"},
		{"toKebabCase", toKebabCase, "HelloWorld", "hello-world"},
		{"toKebabCase", toKebabCase, "helloWorld", "hello-world"},
		{"toKebabCase", toKebabCase, "hello_world", "hello_world"},
		{"toPascalCase", toPascalCase, "hello-world", "HelloWorld"},
		{"toPascalCase", toPascalCase, "hello_world", "HelloWorld"},
		{"toPascalCase", toPascalCase, "helloWorld", "HelloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.input, func(t *testing.T) {
			result := tt.function(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestVersionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{"getSemverMajor", "1.2.3", "1"},
		{"getSemverMajor", "v1.2.3", "1"},
		{"getSemverMajor", "10.0.0-beta.1", "10"},
		{"getSemverMinor", "1.2.3", "2"},
		{"getSemverMinor", "v1.2.3", "2"},
		{"getSemverMinor", "1.15.0", "15"},
		{"getSemverPatch", "1.2.3", "3"},
		{"getSemverPatch", "v1.2.3", "3"},
		{"getSemverPatch", "1.2.10-beta.1", "10"},
		{"getSemverPatch", "1.2.5+build.123", "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.version, func(t *testing.T) {
			var result string
			switch tt.name {
			case "getSemverMajor":
				result = getSemverMajor(tt.version)
			case "getSemverMinor":
				result = getSemverMinor(tt.version)
			case "getSemverPatch":
				result = getSemverPatch(tt.version)
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestVersionComparison(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"v1.0.0", "v1.0.0", 0},
		{"1.2.3", "1.2.4", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			result := compareSemver(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestVersionPrefixFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"addVersionPrefix", "1.0.0", "v1.0.0"},
		{"addVersionPrefix", "v1.0.0", "v1.0.0"},
		{"stripVersionPrefix", "v1.0.0", "1.0.0"},
		{"stripVersionPrefix", "1.0.0", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.input, func(t *testing.T) {
			var result string
			switch tt.name {
			case "addVersionPrefix":
				result = addVersionPrefix(tt.input)
			case "stripVersionPrefix":
				result = stripVersionPrefix(tt.input)
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	config := &models.ProjectConfig{
		Versions: &models.VersionConfig{
			Packages: map[string]string{
				"react": "18.2.0",
				"next":  "14.0.0",
			},
		},
	}

	tests := []struct {
		packageName string
		expected    string
	}{
		{"react", "18.2.0"},
		{"next", "14.0.0"},
		{"unknown", "latest"},
	}

	for _, tt := range tests {
		t.Run(tt.packageName, func(t *testing.T) {
			result := getLatestVersion(config, tt.packageName)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}

	// Test with nil versions
	configNil := &models.ProjectConfig{}
	result := getLatestVersion(configNil, "react")
	if result != "latest" {
		t.Errorf("Expected 'latest' for nil versions, got %s", result)
	}
}

func TestConditionalFunctions(t *testing.T) {
	tests := []struct {
		name      string
		function  func(bool, interface{}, interface{}) interface{}
		condition bool
		trueVal   interface{}
		falseVal  interface{}
		expected  interface{}
	}{
		{"templateIf_true", templateIf, true, "yes", "no", "yes"},
		{"templateIf_false", templateIf, false, "yes", "no", "no"},
		{"templateIfNot_true", templateIfNot, true, "yes", "no", "no"},
		{"templateIfNot_false", templateIfNot, false, "yes", "no", "yes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.condition, tt.trueVal, tt.falseVal)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLogicalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(bool, bool) bool
		a        bool
		b        bool
		expected bool
	}{
		{"templateAnd_true_true", templateAnd, true, true, true},
		{"templateAnd_true_false", templateAnd, true, false, false},
		{"templateAnd_false_true", templateAnd, false, true, false},
		{"templateAnd_false_false", templateAnd, false, false, false},
		{"templateOr_true_true", templateOr, true, true, true},
		{"templateOr_true_false", templateOr, true, false, true},
		{"templateOr_false_true", templateOr, false, true, true},
		{"templateOr_false_false", templateOr, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

	// Test templateNot
	if templateNot(true) != false {
		t.Error("templateNot(true) should return false")
	}
	if templateNot(false) != true {
		t.Error("templateNot(false) should return true")
	}
}

func TestComparisonFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func(interface{}, interface{}) bool
		a        interface{}
		b        interface{}
		expected bool
	}{
		{"templateEq_equal", templateEq, "hello", "hello", true},
		{"templateEq_not_equal", templateEq, "hello", "world", false},
		{"templateNe_equal", templateNe, "hello", "hello", false},
		{"templateNe_not_equal", templateNe, "hello", "world", true},
		{"templateLt_less", templateLt, "a", "b", true},
		{"templateLt_greater", templateLt, "b", "a", false},
		{"templateGt_greater", templateGt, "b", "a", true},
		{"templateGt_less", templateGt, "a", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEmptyFunctions(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		isEmpty  bool
		nonEmpty bool
	}{
		{"nil", nil, true, false},
		{"empty_string", "", true, false},
		{"non_empty_string", "hello", false, true},
		{"empty_slice", []interface{}{}, true, false},
		{"non_empty_slice", []interface{}{"item"}, false, true},
		{"empty_map", map[string]interface{}{}, true, false},
		{"non_empty_map", map[string]interface{}{"key": "value"}, false, true},
		{"number", 42, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isEmpty(tt.value) != tt.isEmpty {
				t.Errorf("isEmpty(%v) expected %v, got %v", tt.value, tt.isEmpty, isEmpty(tt.value))
			}
			if isNonEmpty(tt.value) != tt.nonEmpty {
				t.Errorf("isNonEmpty(%v) expected %v, got %v", tt.value, tt.nonEmpty, isNonEmpty(tt.value))
			}
		})
	}
}

func TestComponentCheckingFunctions(t *testing.T) {
	config := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    false,
				Admin:   true,
			},
			Backend: models.BackendComponents{
				API: true,
			},
			Mobile: models.MobileComponents{
				Android: false,
				IOS:     true,
			},
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: false,
				Terraform:  true,
			},
		},
	}

	tests := []struct {
		name     string
		function func(*models.ProjectConfig) bool
		expected bool
	}{
		{"hasFrontendComponent", hasFrontendComponent, true},
		{"hasBackendComponent", hasBackendComponent, true},
		{"hasMobileComponent", hasMobileComponent, true},
		{"hasInfrastructureComponent", hasInfrastructureComponent, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(config)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

	// Test hasComponent function
	componentTests := []struct {
		componentType string
		componentName string
		expected      bool
	}{
		{"frontend", "main_app", true},
		{"frontend", "home", false},
		{"frontend", "admin", true},
		{"backend", "api", true},
		{"mobile", "android", false},
		{"mobile", "ios", true},
		{"infrastructure", "docker", true},
		{"infrastructure", "kubernetes", false},
		{"infrastructure", "terraform", true},
		{"unknown", "unknown", false},
	}

	for _, tt := range componentTests {
		t.Run("hasComponent_"+tt.componentType+"_"+tt.componentName, func(t *testing.T) {
			result := hasComponent(config, tt.componentType, tt.componentName)
			if result != tt.expected {
				t.Errorf("hasComponent(%s, %s) expected %v, got %v",
					tt.componentType, tt.componentName, tt.expected, result)
			}
		})
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test arithmetic functions
	if add(5, 3) != 8 {
		t.Error("add(5, 3) should return 8")
	}
	if sub(5, 3) != 2 {
		t.Error("sub(5, 3) should return 2")
	}
	if mul(5, 3) != 15 {
		t.Error("mul(5, 3) should return 15")
	}
	if div(6, 3) != 2 {
		t.Error("div(6, 3) should return 2")
	}
	if div(5, 0) != 0 {
		t.Error("div(5, 0) should return 0")
	}
	if mod(7, 3) != 1 {
		t.Error("mod(7, 3) should return 1")
	}
	if mod(5, 0) != 0 {
		t.Error("mod(5, 0) should return 0")
	}

	// Test defaultValue function
	if defaultValue("", "default") != "default" {
		t.Error("defaultValue with empty string should return default")
	}
	if defaultValue("value", "default") != "value" {
		t.Error("defaultValue with non-empty string should return value")
	}

	// Test coalesce function
	if coalesce("", nil, "first") != "first" {
		t.Error("coalesce should return first non-empty value")
	}
	if coalesce("value", "other") != "value" {
		t.Error("coalesce should return first value if not empty")
	}

	// Test quote functions
	if quote("hello") != `"hello"` {
		t.Error("quote should add double quotes")
	}
	if singleQuote("hello") != "'hello'" {
		t.Error("singleQuote should add single quotes")
	}
	if singleQuote("it's") != "'it\\'s'" {
		t.Error("singleQuote should escape single quotes")
	}

	// Test indent functions
	text := "line1\nline2"
	indented := indent(2, text)
	expected := "  line1\n  line2"
	if indented != expected {
		t.Errorf("indent(2, %q) expected %q, got %q", text, expected, indented)
	}

	nindented := nindent(2, text)
	expectedN := "\n  line1\n  line2"
	if nindented != expectedN {
		t.Errorf("nindent(2, %q) expected %q, got %q", text, expectedN, nindented)
	}

	// Test formatTime function
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	formatted := formatTime(testTime, "2006-01-02 15:04:05")
	expected = "2023-12-25 15:30:45"
	if formatted != expected {
		t.Errorf("formatTime expected %s, got %s", expected, formatted)
	}
}

func TestEdgeCases(t *testing.T) {
	// Test empty string cases
	if toCamelCase("") != "" {
		t.Error("toCamelCase with empty string should return empty string")
	}
	if toPascalCase("") != "" {
		t.Error("toPascalCase with empty string should return empty string")
	}

	// Test single character cases
	if toCamelCase("a") != "a" {
		t.Error("toCamelCase with single character should return lowercase")
	}
	if toPascalCase("a") != "A" {
		t.Error("toPascalCase with single character should return uppercase")
	}

	// Test version edge cases
	if getSemverMajor("") != "0" {
		t.Error("getSemverMajor with empty string should return '0'")
	}
	if getSemverMinor("1") != "0" {
		t.Error("getSemverMinor with single number should return '0'")
	}
	if getSemverPatch("1.2") != "0" {
		t.Error("getSemverPatch with major.minor should return '0'")
	}
}
