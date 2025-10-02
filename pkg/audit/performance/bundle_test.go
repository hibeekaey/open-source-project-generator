package performance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBundleAnalyzer(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	assert.NotNil(t, analyzer)
	assert.Equal(t, int64(5*1024*1024), analyzer.maxBundleSize)
	assert.Equal(t, 0.3, analyzer.compressionRatio)
	assert.Equal(t, int64(1024*1024), analyzer.assetSizeThreshold)
}

func TestBundleAnalyzer_AnalyzeBundleSize(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     func(string) error
		expectedAssets int
		expectedSize   int64
		expectError    bool
	}{
		{
			name: "analyze build directory",
			setupFiles: func(tempDir string) error {
				buildDir := filepath.Join(tempDir, "dist")
				if err := os.MkdirAll(buildDir, 0755); err != nil {
					return err
				}

				// Create test files
				files := map[string]string{
					"main.js":    "console.log('hello world');",
					"styles.css": "body { margin: 0; }",
					"image.png":  "fake image data",
				}

				for filename, content := range files {
					err := os.WriteFile(filepath.Join(buildDir, filename), []byte(content), 0644)
					if err != nil {
						return err
					}
				}
				return nil
			},
			expectedAssets: 3,
			expectedSize:   int64(len("console.log('hello world');") + len("body { margin: 0; }") + len("fake image data")),
			expectError:    false,
		},
		{
			name: "analyze source files when no build directory",
			setupFiles: func(tempDir string) error {
				files := map[string]string{
					"main.go":   "package main\n\nfunc main() {}",
					"utils.js":  "function utils() {}",
					"README.md": "# Test Project",
				}

				for filename, content := range files {
					err := os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
					if err != nil {
						return err
					}
				}
				return nil
			},
			expectedAssets: 2, // Only source code files (main.go, utils.js)
			expectedSize:   int64(len("package main\n\nfunc main() {}") + len("function utils() {}")),
			expectError:    false,
		},
		{
			name: "empty directory",
			setupFiles: func(tempDir string) error {
				return nil // No files
			},
			expectedAssets: 0,
			expectedSize:   0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "bundle-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Setup test files
			err = tt.setupFiles(tempDir)
			require.NoError(t, err)

			// Run analysis
			analyzer := NewBundleAnalyzer()
			result, err := analyzer.AnalyzeBundleSize(tempDir)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedAssets, len(result.Assets))
			assert.Equal(t, tt.expectedSize, result.TotalSize)
			assert.Equal(t, tt.expectedAssets, result.Summary.TotalAssets)

			// Verify gzipped size is calculated
			if result.TotalSize > 0 {
				assert.Greater(t, result.GzippedSize, int64(0))
				assert.LessOrEqual(t, result.GzippedSize, result.TotalSize)
			}
		})
	}
}

func TestBundleAnalyzer_GetAssetType(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		filename     string
		expectedType string
	}{
		{"main.js", "javascript"},
		{"app.jsx", "javascript"},
		{"component.ts", "javascript"},
		{"component.tsx", "javascript"},
		{"styles.css", "stylesheet"},
		{"styles.scss", "stylesheet"},
		{"styles.sass", "stylesheet"},
		{"styles.less", "stylesheet"},
		{"image.jpg", "image"},
		{"image.png", "image"},
		{"image.svg", "image"},
		{"font.woff", "font"},
		{"font.ttf", "font"},
		{"data.json", "data"},
		{"index.html", "markup"},
		{"unknown.xyz", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := analyzer.getAssetType(tt.filename)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestBundleAnalyzer_EstimateGzippedSize(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		filename      string
		size          int64
		expectedRatio float64
	}{
		{"main.js", 1000, 0.3},     // Text files compress well
		{"styles.css", 1000, 0.3},  // Text files compress well
		{"data.json", 1000, 0.3},   // Text files compress well
		{"image.jpg", 1000, 1.0},   // Already compressed
		{"image.png", 1000, 1.0},   // Already compressed
		{"unknown.bin", 1000, 0.6}, // Default compression
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := analyzer.estimateGzippedSize(tt.filename, tt.size)
			expected := int64(float64(tt.size) * tt.expectedRatio)
			assert.Equal(t, expected, result)
		})
	}
}

func TestBundleAnalyzer_ShouldSkipFile(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		filepath   string
		shouldSkip bool
	}{
		{"main.js", false},
		{"styles.css", false},
		{"image.png", false},
		{".hidden", true},
		{"app.exe", true},
		{"lib.dll", true},
		{"archive.zip", true},
		{"node_modules/package/index.js", true},
		{"vendor/lib.go", true},
		{".git/config", true},
		{"normal/file.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.filepath, func(t *testing.T) {
			result := analyzer.shouldSkipFile(tt.filepath)
			assert.Equal(t, tt.shouldSkip, result)
		})
	}
}

func TestBundleAnalyzer_IsSourceCodeFile(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		filepath     string
		isSourceCode bool
	}{
		{"main.go", true},
		{"app.js", true},
		{"component.tsx", true},
		{"styles.css", true},
		{"config.yaml", true},
		{"README.md", false},
		{"image.png", false},
		{"binary.exe", false},
		{"archive.zip", false},
	}

	for _, tt := range tests {
		t.Run(tt.filepath, func(t *testing.T) {
			result := analyzer.isSourceCodeFile(tt.filepath)
			assert.Equal(t, tt.isSourceCode, result)
		})
	}
}

func TestBundleAnalyzer_GenerateRecommendations(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		name                    string
		result                  *interfaces.BundleAnalysisResult
		expectedRecommendations int
	}{
		{
			name: "large bundle size",
			result: &interfaces.BundleAnalysisResult{
				TotalSize:   10 * 1024 * 1024, // 10MB - exceeds 5MB limit
				GzippedSize: 3 * 1024 * 1024,
				Assets: []interfaces.BundleAsset{
					{Name: "main.js", Size: 10 * 1024 * 1024, Type: "javascript"},
				},
				Summary: interfaces.BundleAnalysisSummary{},
			},
			expectedRecommendations: 1, // Bundle size recommendation
		},
		{
			name: "poor compression ratio",
			result: &interfaces.BundleAnalysisResult{
				TotalSize:   1024 * 1024, // 1MB
				GzippedSize: 900 * 1024,  // 900KB - poor compression (90%)
				Assets:      []interfaces.BundleAsset{},
				Summary: interfaces.BundleAnalysisSummary{
					CompressionRatio: 0.9, // Set compression ratio directly
				},
			},
			expectedRecommendations: 1, // Compression recommendation
		},
		{
			name: "large individual asset",
			result: &interfaces.BundleAnalysisResult{
				TotalSize:   2 * 1024 * 1024,
				GzippedSize: 1024 * 1024,
				Assets: []interfaces.BundleAsset{
					{Name: "large-image.png", Size: 2 * 1024 * 1024, Type: "image"}, // 2MB asset
				},
				Summary: interfaces.BundleAnalysisSummary{},
			},
			expectedRecommendations: 1, // Large asset recommendation
		},
		{
			name: "javascript heavy bundle",
			result: &interfaces.BundleAnalysisResult{
				TotalSize:   2 * 1024 * 1024,
				GzippedSize: 1024 * 1024,
				Assets: []interfaces.BundleAsset{
					{Name: "main.js", Size: 1200 * 1024, Type: "javascript"}, // >50% of bundle
					{Name: "styles.css", Size: 800 * 1024, Type: "stylesheet"},
				},
				Summary: interfaces.BundleAnalysisSummary{},
			},
			expectedRecommendations: 1, // JavaScript heavy recommendation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer.generateRecommendations(tt.result)
			assert.GreaterOrEqual(t, len(tt.result.Summary.Recommendations), tt.expectedRecommendations)
		})
	}
}

func TestBundleAnalyzer_AnalyzePerformanceIssues(t *testing.T) {
	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "perf-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a large source file
	largeContent := make([]byte, 150*1024) // 150KB
	for i := range largeContent {
		largeContent[i] = 'a'
	}
	err = os.WriteFile(filepath.Join(tempDir, "large.js"), largeContent, 0644)
	require.NoError(t, err)

	// Create a large asset file
	largeAsset := make([]byte, 2*1024*1024) // 2MB
	err = os.WriteFile(filepath.Join(tempDir, "large.png"), largeAsset, 0644)
	require.NoError(t, err)

	// Create a file with performance anti-patterns
	codeWithIssues := `
function badCode() {
    for (let i = 0; i < 100; i++) {
        for (let j = 0; j < 100; j++) {
            for (let k = 0; k < 100; k++) {
                console.log(i, j, k);
            }
        }
    }
    
    document.getElementById('test').innerHTML = 'test';
    console.log('debug statement');
}
`
	err = os.WriteFile(filepath.Join(tempDir, "bad-code.js"), []byte(codeWithIssues), 0644)
	require.NoError(t, err)

	analyzer := NewBundleAnalyzer()
	issues, err := analyzer.AnalyzePerformanceIssues(tempDir)

	require.NoError(t, err)
	assert.Greater(t, len(issues), 0)

	// Check for expected issue types
	issueTypes := make(map[string]bool)
	for _, issue := range issues {
		issueTypes[issue.Type] = true
	}

	assert.True(t, issueTypes["large_file"], "Should detect large source file")
	assert.True(t, issueTypes["large_asset"], "Should detect large asset file")
	assert.True(t, issueTypes["nested_loops"], "Should detect nested loops")
	assert.True(t, issueTypes["dom_manipulation"], "Should detect DOM manipulation")
}

func TestBundleAnalyzer_IsAssetFile(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	tests := []struct {
		filepath string
		isAsset  bool
	}{
		{"image.jpg", true},
		{"image.png", true},
		{"video.mp4", true},
		{"document.pdf", true},
		{"archive.zip", true},
		{"main.js", false},
		{"styles.css", false},
		{"README.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.filepath, func(t *testing.T) {
			result := analyzer.isAssetFile(tt.filepath)
			assert.Equal(t, tt.isAsset, result)
		})
	}
}

func TestBundleAnalyzer_CalculateSummary(t *testing.T) {
	analyzer := NewBundleAnalyzer()

	result := &interfaces.BundleAnalysisResult{
		TotalSize:   1000,
		GzippedSize: 300,
		Assets: []interfaces.BundleAsset{
			{Name: "main.js", Size: 600, Type: "javascript"},
			{Name: "styles.css", Size: 400, Type: "stylesheet"},
		},
		Chunks: []interfaces.BundleChunk{
			{Name: "main", Size: 600},
		},
		Summary: interfaces.BundleAnalysisSummary{},
	}

	analyzer.calculateSummary(result)

	assert.Equal(t, 2, result.Summary.TotalAssets)
	assert.Equal(t, 1, result.Summary.TotalChunks)
	assert.Equal(t, 0.3, result.Summary.CompressionRatio)
	assert.Equal(t, "main.js", result.Summary.LargestAsset)

	// Check percentages
	assert.Equal(t, 60.0, result.Assets[0].Percentage) // 600/1000 * 100
	assert.Equal(t, 40.0, result.Assets[1].Percentage) // 400/1000 * 100
}
