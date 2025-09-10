package models

// TemplateMetadata represents metadata for a template
type TemplateMetadata struct {
	Name                string              `yaml:"name" json:"name" validate:"required,min=1,max=100"`
	Description         string              `yaml:"description" json:"description" validate:"required,max=500"`
	Version             string              `yaml:"version" json:"version" validate:"required,semver"`
	Author              string              `yaml:"author" json:"author" validate:"omitempty,max=100"`
	Dependencies        []string            `yaml:"dependencies" json:"dependencies" validate:"omitempty,dive,required"`
	Variables           []TemplateVar       `yaml:"variables" json:"variables" validate:"omitempty,dive"`
	Conditions          []TemplateCondition `yaml:"conditions" json:"conditions" validate:"omitempty,dive"`
	Tags                []string            `yaml:"tags" json:"tags" validate:"omitempty,dive,required"`
	MinGeneratorVersion string              `yaml:"min_generator_version" json:"min_generator_version" validate:"omitempty,semver"`
}

// TemplateVar represents a template variable definition
type TemplateVar struct {
	Name        string      `yaml:"name" json:"name" validate:"required,min=1,max=50,alphanum"`
	Type        string      `yaml:"type" json:"type" validate:"required,oneof=string int bool float array object"`
	Default     interface{} `yaml:"default" json:"default"`
	Description string      `yaml:"description" json:"description" validate:"required,max=200"`
	Required    bool        `yaml:"required" json:"required"`
	Options     []string    `yaml:"options" json:"options" validate:"omitempty,dive,required"`
	Pattern     string      `yaml:"pattern" json:"pattern" validate:"omitempty"`
}

// TemplateCondition represents a conditional template rendering rule
type TemplateCondition struct {
	Name      string      `yaml:"name" json:"name" validate:"required,min=1,max=50"`
	Component string      `yaml:"component" json:"component" validate:"required"`
	Operator  string      `yaml:"operator" json:"operator" validate:"required,oneof=eq ne gt lt gte lte in nin"`
	Value     interface{} `yaml:"value" json:"value" validate:"required"`
}
