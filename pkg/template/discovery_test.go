package template

import (
	"embed"
	"strings"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock embedded filesystem for testing
//
//go:embed testdata
var testEmbeddedFS embed.FS

func TestNewTemplateDiscovery(t *testing.T) {
	td := NewTemplateDiscovery(testEmbeddedFS)

	assert.NotNil(t, td)
	assert.NotNil(t, td.embeddedFS)
	assert.NotNil(t, td.cache)
	assert.Equal(t, 5*time.Minute, td.cacheTTL)
	assert.Empty(t, td.externalPaths)
}

func TestTemplateDiscovery_DiscoverTemplates(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, templates)

	// Verify templates are sorted by name
	for i := 1; i < len(templates); i++ {
		assert.True(t, templates[i-1].Name <= templates[i].Name,
			"Templates should be sorted by name")
	}

	// Verify cache is populated
	assert.NotEmpty(t, td.cache)
	assert.False(t, td.cacheTime.IsZero())
}

func TestTemplateDiscovery_DiscoverTemplates_Cache(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	// First call
	templates1, err := td.DiscoverTemplates()
	require.NoError(t, err)
	cacheTime1 := td.cacheTime

	// Second call should use cache
	templates2, err := td.DiscoverTemplates()
	require.NoError(t, err)
	cacheTime2 := td.cacheTime

	assert.Equal(t, len(templates1), len(templates2))
	assert.Equal(t, cacheTime1, cacheTime2) // Cache time should not change
}

func TestTemplateDiscovery_DiscoverEmbeddedTemplates(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	templates, err := td.DiscoverEmbeddedTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, templates)

	// Verify all templates have required fields
	for _, tmpl := range templates {
		assert.NotEmpty(t, tmpl.Name, "Template name should not be empty")
		assert.NotEmpty(t, tmpl.Category, "Template category should not be empty")
		assert.NotEmpty(t, tmpl.DisplayName, "Template display name should not be empty")
		assert.Equal(t, "embedded", tmpl.Source, "Template source should be embedded")
		assert.NotEmpty(t, tmpl.Path, "Template path should not be empty")
	}
}

func TestTemplateDiscovery_FilterTemplates(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)
	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)

	tests := []struct {
		name     string
		filter   interfaces.TemplateFilter
		validate func(t *testing.T, filtered []*models.TemplateInfo)
	}{
		{
			name: "filter by category",
			filter: interfaces.TemplateFilter{
				Category: "backend",
			},
			validate: func(t *testing.T, filtered []*models.TemplateInfo) {
				for _, tmpl := range filtered {
					assert.Equal(t, "backend", tmpl.Category)
				}
			},
		},
		{
			name: "filter by technology",
			filter: interfaces.TemplateFilter{
				Technology: "Go",
			},
			validate: func(t *testing.T, filtered []*models.TemplateInfo) {
				for _, tmpl := range filtered {
					// Technology comparison should be case insensitive
					assert.True(t, strings.EqualFold(tmpl.Technology, "Go"))
				}
			},
		},
		{
			name: "filter by tags",
			filter: interfaces.TemplateFilter{
				Tags: []string{"backend"},
			},
			validate: func(t *testing.T, filtered []*models.TemplateInfo) {
				for _, tmpl := range filtered {
					assert.Contains(t, tmpl.Tags, "backend")
				}
			},
		},
		{
			name:   "empty filter returns all",
			filter: interfaces.TemplateFilter{},
			validate: func(t *testing.T, filtered []*models.TemplateInfo) {
				assert.Equal(t, len(templates), len(filtered))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := td.FilterTemplates(templates, tt.filter)
			tt.validate(t, filtered)
		})
	}
}

func TestTemplateDiscovery_SearchTemplates(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)
	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)

	tests := []struct {
		name     string
		query    string
		validate func(t *testing.T, results []*models.TemplateInfo)
	}{
		{
			name:  "search by name",
			query: "go",
			validate: func(t *testing.T, results []*models.TemplateInfo) {
				for _, tmpl := range results {
					found := false
					if containsIgnoreCase(tmpl.Name, "go") ||
						containsIgnoreCase(tmpl.DisplayName, "go") ||
						containsIgnoreCase(tmpl.Description, "go") ||
						containsIgnoreCase(tmpl.Technology, "go") {
						found = true
					}
					for _, tag := range tmpl.Tags {
						if containsIgnoreCase(tag, "go") {
							found = true
							break
						}
					}
					for _, keyword := range tmpl.Metadata.Keywords {
						if containsIgnoreCase(keyword, "go") {
							found = true
							break
						}
					}
					assert.True(t, found, "Template should match search query")
				}
			},
		},
		{
			name:  "search by technology",
			query: "next",
			validate: func(t *testing.T, results []*models.TemplateInfo) {
				assert.NotEmpty(t, results, "Should find Next.js templates")
			},
		},
		{
			name:  "empty query returns no results",
			query: "",
			validate: func(t *testing.T, results []*models.TemplateInfo) {
				assert.Empty(t, results, "Empty query should return no results")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := td.SearchTemplates(templates, tt.query)
			tt.validate(t, results)
		})
	}
}

func TestTemplateDiscovery_GetTemplateByName(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)
	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	// Test finding existing template
	firstTemplate := templates[0]
	found := td.GetTemplateByName(templates, firstTemplate.Name)
	assert.NotNil(t, found)
	assert.Equal(t, firstTemplate.Name, found.Name)

	// Test finding non-existing template
	notFound := td.GetTemplateByName(templates, "non-existing-template")
	assert.Nil(t, notFound)
}

func TestTemplateDiscovery_GetTemplatesByCategory(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)
	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)

	// Get all unique categories
	categories := make(map[string]bool)
	for _, tmpl := range templates {
		categories[tmpl.Category] = true
	}

	// Test each category
	for category := range categories {
		filtered := td.GetTemplatesByCategory(templates, category)
		for _, tmpl := range filtered {
			assert.Equal(t, category, tmpl.Category)
		}
	}
}

func TestTemplateDiscovery_GetTemplatesByTechnology(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)
	templates, err := td.DiscoverTemplates()
	require.NoError(t, err)

	// Get all unique technologies
	technologies := make(map[string]bool)
	for _, tmpl := range templates {
		if tmpl.Technology != "" {
			technologies[tmpl.Technology] = true
		}
	}

	// Test each technology
	for technology := range technologies {
		filtered := td.GetTemplatesByTechnology(templates, technology)
		for _, tmpl := range filtered {
			assert.Equal(t, technology, tmpl.Technology)
		}
	}
}

func TestTemplateDiscovery_ClearCache(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	// Populate cache
	_, err := td.DiscoverTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, td.cache)
	assert.False(t, td.cacheTime.IsZero())

	// Clear cache
	td.ClearCache()
	assert.Empty(t, td.cache)
	assert.True(t, td.cacheTime.IsZero())
}

func TestTemplateDiscovery_RefreshCache(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	// Populate cache
	_, err := td.DiscoverTemplates()
	require.NoError(t, err)
	originalCacheTime := td.cacheTime

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Refresh cache
	err = td.RefreshCache()
	require.NoError(t, err)
	assert.True(t, td.cacheTime.After(originalCacheTime))
}

func TestTemplateDiscovery_ExternalPaths(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	// Initially empty
	assert.Empty(t, td.GetExternalPaths())

	// Add paths
	td.AddExternalPath("/path/to/templates1")
	td.AddExternalPath("/path/to/templates2")

	paths := td.GetExternalPaths()
	assert.Len(t, paths, 2)
	assert.Contains(t, paths, "/path/to/templates1")
	assert.Contains(t, paths, "/path/to/templates2")

	// Remove path
	td.RemoveExternalPath("/path/to/templates1")
	paths = td.GetExternalPaths()
	assert.Len(t, paths, 1)
	assert.Contains(t, paths, "/path/to/templates2")
	assert.NotContains(t, paths, "/path/to/templates1")
}

func TestTemplateDiscovery_isTemplateDirectory(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	tests := []struct {
		path     string
		expected bool
	}{
		{"templates/backend/go-gin", true},
		{"templates/frontend/nextjs-app", true},
		{"templates", false},
		{"templates/backend", false},
		{"some/other/path", false},
		{"templates/backend/go-gin/subdir", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := td.isTemplateDirectory(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateDiscovery_formatDisplayName(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	tests := []struct {
		input    string
		expected string
	}{
		{"go-gin", "Go Gin"},
		{"nextjs-app", "Nextjs App"},
		{"android-kotlin", "Android Kotlin"},
		{"single", "Single"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := td.formatDisplayName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateDiscovery_inferTechnology(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	tests := []struct {
		templateName string
		expected     string
	}{
		{"go-gin", "Go"},
		{"nextjs-app", "Next.js"},
		{"react-component", "React"},
		{"node-api", "Node.js"},
		{"python-flask", "Python"},
		{"java-spring", "Java"},
		{"android-kotlin", "Kotlin"},
		{"ios-swift", "Swift"},
		{"docker-compose", "Docker"},
		{"kubernetes-deployment", "Kubernetes"},
		{"terraform-aws", "Terraform"},
		{"unknown-template", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.templateName, func(t *testing.T) {
			result := td.inferTechnology(tt.templateName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateDiscovery_inferTags(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	tests := []struct {
		templateName string
		category     string
		expectedTags []string
	}{
		{
			templateName: "go-gin",
			category:     "backend",
			expectedTags: []string{"backend", "go", "backend", "api"},
		},
		{
			templateName: "nextjs-admin",
			category:     "frontend",
			expectedTags: []string{"frontend", "nextjs", "react", "frontend", "web", "admin", "dashboard"},
		},
		{
			templateName: "android-kotlin",
			category:     "mobile",
			expectedTags: []string{"mobile", "android", "mobile", "kotlin", "android"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.templateName, func(t *testing.T) {
			result := td.inferTags(tt.templateName, tt.category)

			// Check that all expected tags are present (allowing for duplicates)
			for _, expectedTag := range tt.expectedTags {
				assert.Contains(t, result, expectedTag)
			}
		})
	}
}

func TestTemplateDiscovery_matchesFilter(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	template := &models.TemplateInfo{
		Name:       "test-template",
		Category:   "backend",
		Technology: "Go",
		Version:    "1.2.0",
		Tags:       []string{"api", "backend", "go"},
	}

	tests := []struct {
		name     string
		filter   interfaces.TemplateFilter
		expected bool
	}{
		{
			name:     "empty filter matches all",
			filter:   interfaces.TemplateFilter{},
			expected: true,
		},
		{
			name:     "category match",
			filter:   interfaces.TemplateFilter{Category: "backend"},
			expected: true,
		},
		{
			name:     "category no match",
			filter:   interfaces.TemplateFilter{Category: "frontend"},
			expected: false,
		},
		{
			name:     "technology match",
			filter:   interfaces.TemplateFilter{Technology: "Go"},
			expected: true,
		},
		{
			name:     "technology no match",
			filter:   interfaces.TemplateFilter{Technology: "Python"},
			expected: false,
		},
		{
			name:     "tags match",
			filter:   interfaces.TemplateFilter{Tags: []string{"api"}},
			expected: true,
		},
		{
			name:     "tags no match",
			filter:   interfaces.TemplateFilter{Tags: []string{"frontend"}},
			expected: false,
		},
		{
			name:     "version range match",
			filter:   interfaces.TemplateFilter{MinVersion: "1.0.0", MaxVersion: "2.0.0"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := td.matchesFilter(template, tt.filter)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateDiscovery_matchesQuery(t *testing.T) {
	td := NewTemplateDiscovery(embeddedTemplates)

	template := &models.TemplateInfo{
		Name:        "go-gin-api",
		DisplayName: "Go Gin API",
		Description: "A REST API template using Go and Gin framework",
		Technology:  "Go",
		Tags:        []string{"api", "backend", "go"},
		Metadata: models.TemplateMetadata{
			Keywords: []string{"rest", "microservice"},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "match in name",
			query:    "gin",
			expected: true,
		},
		{
			name:     "match in display name",
			query:    "api",
			expected: true,
		},
		{
			name:     "match in description",
			query:    "rest",
			expected: true,
		},
		{
			name:     "match in technology",
			query:    "go",
			expected: true,
		},
		{
			name:     "match in tags",
			query:    "backend",
			expected: true,
		},
		{
			name:     "match in keywords",
			query:    "microservice",
			expected: true,
		},
		{
			name:     "no match",
			query:    "python",
			expected: false,
		},
		{
			name:     "case insensitive match",
			query:    "gin",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := td.matchesQuery(template, tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function for case-insensitive string contains check
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
