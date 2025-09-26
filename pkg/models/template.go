package models

import "time"

// TemplateMetadata represents comprehensive metadata for a template
type TemplateMetadata struct {
	Name         string                 `yaml:"name" json:"name"`
	DisplayName  string                 `yaml:"display_name" json:"display_name"`
	Description  string                 `yaml:"description" json:"description"`
	Version      string                 `yaml:"version" json:"version"`
	Author       string                 `yaml:"author" json:"author"`
	License      string                 `yaml:"license" json:"license"`
	Category     string                 `yaml:"category" json:"category"`
	Technology   string                 `yaml:"technology" json:"technology"`
	Tags         []string               `yaml:"tags" json:"tags"`
	Dependencies []string               `yaml:"dependencies" json:"dependencies"`
	Variables    map[string]TemplateVar `yaml:"variables" json:"variables"`
	CreatedAt    time.Time              `yaml:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `yaml:"updated_at" json:"updated_at"`
	Homepage     string                 `yaml:"homepage" json:"homepage"`
	Repository   string                 `yaml:"repository" json:"repository"`
	Keywords     []string               `yaml:"keywords" json:"keywords"`
	MinVersion   string                 `yaml:"min_version" json:"min_version"`
	MaxVersion   string                 `yaml:"max_version" json:"max_version"`
}

// TemplateVar represents a template variable definition with validation
type TemplateVar struct {
	Name        string      `yaml:"name" json:"name"`
	Type        string      `yaml:"type" json:"type"`
	Default     interface{} `yaml:"default" json:"default"`
	Description string      `yaml:"description" json:"description"`
	Required    bool        `yaml:"required" json:"required"`
	Pattern     string      `yaml:"pattern" json:"pattern"`
	Enum        []string    `yaml:"enum" json:"enum"`
	MinLength   int         `yaml:"min_length" json:"min_length"`
	MaxLength   int         `yaml:"max_length" json:"max_length"`
}

// TemplateInfo represents complete information about a template
type TemplateInfo struct {
	Name         string           `json:"name"`
	DisplayName  string           `json:"display_name"`
	Description  string           `json:"description"`
	Category     string           `json:"category"`
	Technology   string           `json:"technology"`
	Version      string           `json:"version"`
	Tags         []string         `json:"tags"`
	Dependencies []string         `json:"dependencies"`
	Metadata     TemplateMetadata `json:"metadata"`
	Path         string           `json:"path"`
	Source       string           `json:"source"` // embedded, file, git, etc.
	Size         int64            `json:"size"`
	FileCount    int              `json:"file_count"`
	LastModified time.Time        `json:"last_modified"`
}

// TemplateFilter defines filtering options for template listing
type TemplateFilter struct {
	Category   string   `json:"category,omitempty"`
	Technology string   `json:"technology,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	MinVersion string   `json:"min_version,omitempty"`
	MaxVersion string   `json:"max_version,omitempty"`
	Source     string   `json:"source,omitempty"`
	Keywords   []string `json:"keywords,omitempty"`
	Author     string   `json:"author,omitempty"`
}

// TemplateValidationResult contains the result of template validation
type TemplateValidationResult struct {
	TemplateName string                    `json:"template_name"`
	Valid        bool                      `json:"valid"`
	Issues       []ValidationIssue         `json:"issues"`
	Warnings     []ValidationIssue         `json:"warnings"`
	Summary      TemplateValidationSummary `json:"summary"`
}

// TemplateValidationSummary provides a summary of template validation results
type TemplateValidationSummary struct {
	TotalIssues  int `json:"total_issues"`
	ErrorCount   int `json:"error_count"`
	WarningCount int `json:"warning_count"`
	InfoCount    int `json:"info_count"`
	FixableCount int `json:"fixable_count"`
}
