package processor

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProcessingEngine(t *testing.T) {
	engine := NewProcessingEngine()
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.funcMap)

	// Verify some default functions are registered
	assert.Contains(t, engine.funcMap, "lower")
	assert.Contains(t, engine.funcMap, "upper")
	assert.Contains(t, engine.funcMap, "camelCase")
	assert.Contains(t, engine.funcMap, "nodeVersion")
}

func TestProcessingEngine_RegisterFunctions(t *testing.T) {
	engine := NewProcessingEngine()

	customFuncs := template.FuncMap{
		"customFunc": func() string { return "custom" },
		"testFunc":   func(s string) string { return "test_" + s },
	}

	engine.RegisterFunctions(customFuncs)

	assert.Contains(t, engine.funcMap, "customFunc")
	assert.Contains(t, engine.funcMap, "testFunc")
}

func TestProcessingEngine_LoadTemplate(t *testing.T) {
	engine := NewProcessingEngine()

	// Create a temporary template file
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")
	templateContent := "Hello {{.Name}}!"

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Load the template
	tmpl, err := engine.LoadTemplate(templatePath)
	require.NoError(t, err)
	assert.NotNil(t, tmpl)
	assert.Equal(t, "test.tmpl", tmpl.Name())
}

func TestProcessingEngine_RenderTemplate(t *testing.T) {
	engine := NewProcessingEngine()

	// Create a simple template
	tmpl, err := template.New("test").Funcs(engine.funcMap).Parse("Hello {{.Name}}!")
	require.NoError(t, err)

	// Test data
	data := struct {
		Name string
	}{
		Name: "World",
	}

	// Render the template
	result, err := engine.RenderTemplate(tmpl, data)
	require.NoError(t, err)
	assert.Equal(t, "Hello World!", string(result))
}

func TestProcessingEngine_ProcessTemplate(t *testing.T) {
	engine := NewProcessingEngine()

	// Create a temporary template file
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, "test.tmpl")
	templateContent := "Project: {{.Name}}\nOrganization: {{.Organization}}"

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(t, err)

	// Create test config
	config := &models.ProjectConfig{
		Name:         "TestProject",
		Organization: "TestOrg",
	}

	// Process the template
	result, err := engine.ProcessTemplate(templatePath, config)
	require.NoError(t, err)

	expected := "Project: TestProject\nOrganization: TestOrg"
	assert.Equal(t, expected, string(result))
}

func TestProcessingEngine_ProcessDirectory(t *testing.T) {
	engine := NewProcessingEngine()

	// Create temporary directories
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	// Create template structure
	err := os.MkdirAll(filepath.Join(templateDir, "subdir"), 0755)
	require.NoError(t, err)

	// Create template files
	templateFiles := map[string]string{
		"README.md.tmpl":        "# {{.Name}}\n{{.Description}}",
		"package.json.tmpl":     `{"name": "{{.Name | lower}}", "version": "1.0.0"}`,
		"subdir/config.go.tmpl": "package main\n// Project: {{.Name}}",
		"static.txt":            "This is a static file",
	}

	for path, content := range templateFiles {
		fullPath := filepath.Join(templateDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create test config
	config := &models.ProjectConfig{
		Name:        "MyProject",
		Description: "A test project",
	}

	// Process the directory
	err = engine.ProcessDirectory(templateDir, outputDir, config)
	require.NoError(t, err)

	// Verify output files
	readmeContent, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	require.NoError(t, err)
	assert.Equal(t, "# MyProject\nA test project", string(readmeContent))

	packageContent, err := os.ReadFile(filepath.Join(outputDir, "package.json"))
	require.NoError(t, err)
	assert.Equal(t, `{"name": "myproject", "version": "1.0.0"}`, string(packageContent))

	configContent, err := os.ReadFile(filepath.Join(outputDir, "subdir", "config.go"))
	require.NoError(t, err)
	assert.Equal(t, "package main\n// Project: MyProject", string(configContent))

	staticContent, err := os.ReadFile(filepath.Join(outputDir, "static.txt"))
	require.NoError(t, err)
	assert.Equal(t, "This is a static file", string(staticContent))
}

func TestProcessingEngine_ProcessDirectory_SkipsDisabledFiles(t *testing.T) {
	engine := NewProcessingEngine()

	// Create temporary directories
	templateDir := t.TempDir()
	outputDir := t.TempDir()

	// Create template files including disabled ones
	templateFiles := map[string]string{
		"enabled.txt.tmpl":           "Enabled: {{.Name}}",
		"disabled.txt.tmpl.disabled": "Disabled: {{.Name}}",
	}

	for path, content := range templateFiles {
		fullPath := filepath.Join(templateDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create test config
	config := &models.ProjectConfig{
		Name: "TestProject",
	}

	// Process the directory
	err := engine.ProcessDirectory(templateDir, outputDir, config)
	require.NoError(t, err)

	// Verify enabled file exists
	_, err = os.Stat(filepath.Join(outputDir, "enabled.txt"))
	assert.NoError(t, err)

	// Verify disabled file does not exist
	_, err = os.Stat(filepath.Join(outputDir, "disabled.txt"))
	assert.True(t, os.IsNotExist(err))
}

// Test string manipulation functions
func TestStringFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function string
		input    string
		expected string
	}{
		{"camelCase", "camelCase", "hello-world", "helloWorld"},
		{"camelCase existing", "camelCase", "helloWorld", "helloWorld"},
		{"snakeCase", "snakeCase", "HelloWorld", "hello_world"},
		{"kebabCase", "kebabCase", "HelloWorld", "hello-world"},
		{"pascalCase", "pascalCase", "hello-world", "HelloWorld"},
		{"pascalCase from camel", "pascalCase", "helloWorld", "HelloWorld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewProcessingEngine()

			tmplStr := "{{" + tt.function + " .Input}}"
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tmplStr)
			require.NoError(t, err)

			data := map[string]interface{}{
				"Input": tt.input,
			}

			result, err := engine.RenderTemplate(tmpl, data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test version functions
func TestVersionFunctions(t *testing.T) {
	config := &models.ProjectConfig{
		Versions: &models.VersionConfig{
			Node: "18.17.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"next":  "13.4.0",
				"react": "18.2.0",
			},
		},
	}

	tests := []struct {
		name     string
		function string
		expected string
	}{
		{"nodeVersion", "nodeVersion", "18.17.0"},
		{"goVersion", "goVersion", "1.21.0"},
		{"nextjsVersion", "nextjsVersion", "13.4.0"},
		{"reactVersion", "reactVersion", "18.2.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewProcessingEngine()

			tmplStr := "{{" + tt.function + " .}}"
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tmplStr)
			require.NoError(t, err)

			result, err := engine.RenderTemplate(tmpl, config)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test semver functions
func TestSemverFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function string
		version  string
		expected string
	}{
		{"semverMajor", "semverMajor", "1.2.3", "1"},
		{"semverMajor with v", "semverMajor", "v1.2.3", "1"},
		{"semverMinor", "semverMinor", "1.2.3", "2"},
		{"semverPatch", "semverPatch", "1.2.3", "3"},
		{"semverPatch with prerelease", "semverPatch", "1.2.3-alpha.1", "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewProcessingEngine()

			tmplStr := "{{" + tt.function + " .Version}}"
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tmplStr)
			require.NoError(t, err)

			data := map[string]string{"Version": tt.version}
			result, err := engine.RenderTemplate(tmpl, data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test component checking functions
func TestComponentFunctions(t *testing.T) {
	config := &models.ProjectConfig{
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Admin: false,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     false,
			},
		},
	}

	tests := []struct {
		name     string
		function string
		expected string
	}{
		{"hasFrontend", "hasFrontend", "true"},
		{"hasBackend", "hasBackend", "true"},
		{"hasMobile", "hasMobile", "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewProcessingEngine()

			tmplStr := "{{" + tt.function + " .}}"
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tmplStr)
			require.NoError(t, err)

			result, err := engine.RenderTemplate(tmpl, config)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test conditional functions
func TestConditionalFunctions(t *testing.T) {
	engine := NewProcessingEngine()

	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{"if true", "{{if .Condition}}yes{{else}}no{{end}}", map[string]bool{"Condition": true}, "yes"},
		{"if false", "{{if .Condition}}yes{{else}}no{{end}}", map[string]bool{"Condition": false}, "no"},
		{"empty string", "{{empty .Value}}", map[string]string{"Value": ""}, "true"},
		{"nonempty string", "{{nonempty .Value}}", map[string]string{"Value": "test"}, "true"},
		{"eq true", "{{eq .A .B}}", map[string]string{"A": "test", "B": "test"}, "true"},
		{"eq false", "{{eq .A .B}}", map[string]string{"A": "test", "B": "other"}, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tt.template)
			require.NoError(t, err)

			result, err := engine.RenderTemplate(tmpl, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test utility functions
func TestUtilityFunctions(t *testing.T) {
	engine := NewProcessingEngine()

	tests := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{"add", "{{add .A .B}}", map[string]int{"A": 5, "B": 3}, "8"},
		{"sub", "{{sub .A .B}}", map[string]int{"A": 5, "B": 3}, "2"},
		{"mul", "{{mul .A .B}}", map[string]int{"A": 5, "B": 3}, "15"},
		{"div", "{{div .A .B}}", map[string]int{"A": 15, "B": 3}, "5"},
		{"quote", "{{quote .Value}}", map[string]string{"Value": "hello world"}, `"hello world"`},
		{"default empty", "{{default .Empty \"fallback\"}}", map[string]string{"Empty": ""}, "fallback"},
		{"default nonempty", "{{default .Value \"fallback\"}}", map[string]string{"Value": "actual"}, "actual"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tt.template)
			require.NoError(t, err)

			result, err := engine.RenderTemplate(tmpl, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test GitHub Actions functions
func TestGitHubActionsFunctions(t *testing.T) {
	engine := NewProcessingEngine()

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"secrets", "{{secrets \"API_KEY\"}}", "${{ secrets.API_KEY }}"},
		{"matrix", "{{matrix \"os\"}}", "${{ matrix.os }}"},
		{"github", "{{github \"ref\"}}", "${{ github.ref }}"},
		{"env", "{{env \"NODE_ENV\"}}", "${{ env.NODE_ENV }}"},
		{"env with default", "{{env \"NODE_ENV\" \"development\"}}", "${{ env.NODE_ENV || 'development' }}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(tt.template)
			require.NoError(t, err)

			result, err := engine.RenderTemplate(tmpl, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Test complex template processing
func TestComplexTemplateProcessing(t *testing.T) {
	engine := NewProcessingEngine()

	templateContent := `
# {{.Name | title}}

{{if hasFrontend . -}}
## Frontend Components
- Next.js App: {{.Components.Frontend.NextJS.App}}
- Admin Panel: {{.Components.Frontend.NextJS.Admin}}
{{- end}}

{{if hasBackend . -}}
## Backend
- Go/Gin API: {{.Components.Backend.GoGin}}
{{- end}}

## Versions
- Node: {{nodeVersion .}}
- Go: {{goVersion .}}
- Next.js: {{nextjsVersion .}}

## Package Info
- Package Name: {{.Name | kebabCase}}
- Camel Case: {{.Name | camelCase}}
- Snake Case: {{.Name | snakeCase}}
`

	config := &models.ProjectConfig{
		Name: "MyAwesomeProject",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				NextJS: models.NextJSComponents{
					App:   true,
					Admin: false,
				},
			},
			Backend: models.BackendComponents{
				GoGin: true,
			},
		},
		Versions: &models.VersionConfig{
			Node: "18.17.0",
			Go:   "1.21.0",
			Packages: map[string]string{
				"next": "13.4.0",
			},
		},
	}

	tmpl, err := template.New("test").Funcs(engine.funcMap).Parse(templateContent)
	require.NoError(t, err)

	result, err := engine.RenderTemplate(tmpl, config)
	require.NoError(t, err)

	resultStr := string(result)

	// Verify key parts of the output
	assert.Contains(t, resultStr, "# Myawesomeproject")
	assert.Contains(t, resultStr, "## Frontend Components")
	assert.Contains(t, resultStr, "Next.js App: true")
	assert.Contains(t, resultStr, "Admin Panel: false")
	assert.Contains(t, resultStr, "## Backend")
	assert.Contains(t, resultStr, "Go/Gin API: true")
	assert.Contains(t, resultStr, "Node: 18.17.0")
	assert.Contains(t, resultStr, "Go: 1.21.0")
	assert.Contains(t, resultStr, "Next.js: 13.4.0")
	assert.Contains(t, resultStr, "Package Name: my-awesome-project")
	assert.Contains(t, resultStr, "Camel Case: myAwesomeProject")
	assert.Contains(t, resultStr, "Snake Case: my_awesome_project")
}

// Benchmark tests
func BenchmarkProcessingEngine_ProcessTemplate(b *testing.B) {
	engine := NewProcessingEngine()

	// Create a temporary template file
	tmpDir := b.TempDir()
	templatePath := filepath.Join(tmpDir, "bench.tmpl")
	templateContent := "Project: {{.Name | title}}\nDescription: {{.Description}}\nPackage: {{.Name | kebabCase}}"

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	require.NoError(b, err)

	config := &models.ProjectConfig{
		Name:        "BenchmarkProject",
		Description: "A project for benchmarking template processing",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.ProcessTemplate(templatePath, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessingEngine_RenderTemplate(b *testing.B) {
	engine := NewProcessingEngine()

	tmpl, err := template.New("bench").Funcs(engine.funcMap).Parse("{{.Name | title}} - {{.Description}}")
	require.NoError(b, err)

	data := map[string]string{
		"Name":        "BenchmarkProject",
		"Description": "A project for benchmarking",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.RenderTemplate(tmpl, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test error cases
func TestProcessingEngine_ErrorCases(t *testing.T) {
	engine := NewProcessingEngine()

	t.Run("LoadTemplate with invalid path", func(t *testing.T) {
		_, err := engine.LoadTemplate("/nonexistent/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read template file")
	})

	t.Run("ProcessTemplate with invalid template", func(t *testing.T) {
		tmpDir := t.TempDir()
		templatePath := filepath.Join(tmpDir, "invalid.tmpl")
		invalidContent := "{{.Name" // Missing closing brace

		err := os.WriteFile(templatePath, []byte(invalidContent), 0644)
		require.NoError(t, err)

		config := &models.ProjectConfig{Name: "Test"}
		_, err = engine.ProcessTemplate(templatePath, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load template")
	})

	t.Run("RenderTemplate with invalid template syntax", func(t *testing.T) {
		_, err := template.New("test").Funcs(engine.funcMap).Parse("{{.Name | invalidFunction}}")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "function \"invalidFunction\" not defined")
	})

	t.Run("ProcessDirectory with invalid template dir", func(t *testing.T) {
		outputDir := t.TempDir()
		err := engine.ProcessDirectory("/nonexistent", outputDir, &models.ProjectConfig{})
		assert.Error(t, err)
	})
}
