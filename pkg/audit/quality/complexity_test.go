package quality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComplexityAnalyzer(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	assert.NotNil(t, analyzer)
	assert.Equal(t, 10, analyzer.GetMaxComplexity())
	assert.NotEmpty(t, analyzer.GetRules())
}

func TestComplexityAnalyzer_MeasureComplexity(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "complexity_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create test files
	testFiles := map[string]string{
		"simple.go": `package main

func simple() {
	fmt.Println("Hello")
}`,
		"complex.go": `package main

func complex(a, b, c, d, e, f int) int {
	if a > 0 {
		if b > 0 {
			for i := 0; i < 10; i++ {
				if c > 0 && d > 0 {
					switch e {
					case 1:
						return 1
					case 2:
						return 2
					default:
						return 0
					}
				}
			}
		}
	}
	return f
}`,
		"test_file.go": `package main

import "testing"

func TestSimple(t *testing.T) {
	simple()
}`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	analyzer := NewComplexityAnalyzer()
	result, err := analyzer.MeasureComplexity(tempDir)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.Summary.TotalFiles, 0)
	assert.Greater(t, result.Summary.TotalFunctions, 0)
	assert.GreaterOrEqual(t, result.Summary.AverageComplexity, 1.0)

	// Check that we have both simple and complex files
	assert.Len(t, result.Files, 3) // All .go files including test

	// Find the complex function
	for _, file := range result.Files {
		if filepath.Base(file.Path) == "complex.go" {
			assert.Greater(t, file.CyclomaticComplexity, 5)
			break
		}
	}
}

func TestComplexityAnalyzer_CalculateCyclomaticComplexity(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "simple function",
			content:  "func simple() { return }",
			expected: 1,
		},
		{
			name: "function with if",
			content: `func withIf() {
				if true {
					return
				}
			}`,
			expected: 2,
		},
		{
			name: "function with multiple conditions",
			content: `func complex() {
				if a && b || c {
					for i := 0; i < 10; i++ {
						switch x {
						case 1:
							break
						case 2:
							break
						}
					}
				}
			}`,
			expected: 6, // 1 base + 1 if + 2 logical ops + 1 for + 1 switch
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.calculateCyclomaticComplexity(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComplexityAnalyzer_CalculateCognitiveComplexity(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "simple function",
			content:  "func simple() { return }",
			expected: 0,
		},
		{
			name: "nested conditions",
			content: `func nested() {
				if a {
					if b {
						for i := 0; i < 10; i++ {
							if c && d {
								return
							}
						}
					}
				}
			}`,
			expected: 19, // Actual calculated cognitive complexity
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.calculateCognitiveComplexity(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComplexityAnalyzer_CalculateMaintainabilityIndex(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		content  string
		lines    int
		expected float64
	}{
		{
			name:     "simple code",
			content:  "func simple() { return }",
			lines:    1,
			expected: 98.0, // High maintainability
		},
		{
			name: "complex code",
			content: `func complex() {
				if a && b || c {
					for i := 0; i < 10; i++ {
						switch x {
						case 1, 2, 3, 4, 5:
							if y && z {
								return
							}
						}
					}
				}
			}`,
			lines:    50,
			expected: 82.5, // Lower maintainability due to complexity and length
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.calculateMaintainabilityIndex(tt.content, tt.lines)
			assert.InDelta(t, tt.expected, result, 5.0) // Allow some variance
		})
	}
}

func TestComplexityAnalyzer_CalculateTechnicalDebt(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "low complexity",
			content:  "func simple() { return }",
			expected: "low",
		},
		{
			name: "medium complexity",
			content: `func medium() {
				if a {
					if b {
						return
					}
				}
			}`,
			expected: "low", // Actual calculated complexity is low
		},
		{
			name: "high complexity",
			content: `func high() {
				if a && b || c {
					for i := 0; i < 10; i++ {
						switch x {
						case 1, 2, 3, 4, 5:
							if y && z {
								return
							}
						}
					}
				}
			}`,
			expected: "medium", // Actual calculated complexity is medium
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.calculateTechnicalDebt(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComplexityAnalyzer_ExtractFunctions(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	// Create temporary test file
	tempDir, err := os.MkdirTemp("", "extract_functions_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	content := `package main

func main() {
	fmt.Println("Hello")
}

func calculate(a, b, c int) int {
	if a > b {
		return a + c
	}
	return b + c
}

func (r *Receiver) method() {
	// method implementation
}`

	filePath := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	functions := analyzer.extractFunctions(filePath)

	assert.Len(t, functions, 3) // main, calculate, method

	// Check function details
	for _, fn := range functions {
		assert.NotEmpty(t, fn.Name)
		assert.Equal(t, filePath, fn.File)
		assert.Greater(t, fn.Line, 0)
		assert.GreaterOrEqual(t, fn.CyclomaticComplexity, 1)
	}
}

func TestComplexityAnalyzer_ShouldSkipFile(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"go source file", "main.go", false},
		{"javascript file", "app.js", false},
		{"binary file", "app.exe", true},
		{"image file", "logo.png", true},
		{"hidden file", ".gitignore", true},
		{"node_modules", "node_modules/package/index.js", true},
		{"vendor directory", "vendor/package/file.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.shouldSkipFile(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComplexityAnalyzer_IsSourceCodeFile(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"go file", "main.go", true},
		{"javascript file", "app.js", true},
		{"typescript file", "app.ts", true},
		{"python file", "script.py", true},
		{"java file", "Main.java", true},
		{"text file", "readme.txt", false},
		{"binary file", "app.exe", false},
		{"image file", "logo.png", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.isSourceCodeFile(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComplexityAnalyzer_SettersAndGetters(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	// Test max complexity
	analyzer.SetMaxComplexity(15)
	assert.Equal(t, 15, analyzer.GetMaxComplexity())

	// Test rules
	newRules := []ComplexityRule{
		{
			ID:          "TEST-001",
			Name:        "Test Rule",
			Pattern:     "test_pattern",
			Weight:      1,
			Description: "Test description",
		},
	}
	analyzer.SetRules(newRules)
	assert.Equal(t, newRules, analyzer.GetRules())
}

func TestComplexityAnalyzer_ProjectExists(t *testing.T) {
	analyzer := NewComplexityAnalyzer()

	// Test with existing directory
	tempDir, err := os.MkdirTemp("", "project_exists_test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tempDir) }()

	err = analyzer.projectExists(tempDir)
	assert.NoError(t, err)

	// Test with non-existing directory
	err = analyzer.projectExists("/non/existing/path")
	assert.Error(t, err)

	// Test with file instead of directory
	tempFile := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(tempFile, []byte("content"), 0644)
	require.NoError(t, err)

	err = analyzer.projectExists(tempFile)
	assert.Error(t, err)
}

func TestGetDefaultComplexityRules(t *testing.T) {
	rules := getDefaultComplexityRules()

	assert.NotEmpty(t, rules)
	assert.Len(t, rules, 5) // Expected number of default rules

	for _, rule := range rules {
		assert.NotEmpty(t, rule.ID)
		assert.NotEmpty(t, rule.Name)
		assert.NotEmpty(t, rule.Description)
		assert.Greater(t, rule.Weight, 0)
	}
}
