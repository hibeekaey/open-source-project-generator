package formatters

import (
	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// TemplateConverter handles conversion between UI and saved template formats
type TemplateConverter struct{}

// NewTemplateConverter creates a new template converter
func NewTemplateConverter() *TemplateConverter {
	return &TemplateConverter{}
}

// ConvertToSaved converts UI template selections to saved format
func (tc *TemplateConverter) ConvertToSaved(selections []interfaces.TemplateSelection) []config.TemplateSelection {
	result := make([]config.TemplateSelection, len(selections))
	for i, sel := range selections {
		result[i] = config.TemplateSelection{
			TemplateName: sel.Template.Name,
			Category:     sel.Template.Category,
			Technology:   sel.Template.Technology,
			Version:      sel.Template.Version,
			Selected:     sel.Selected,
			Options:      sel.Options,
		}
	}
	return result
}

// ConvertFromSaved converts saved template selections to UI format
func (tc *TemplateConverter) ConvertFromSaved(saved []config.TemplateSelection) []interfaces.TemplateSelection {
	result := make([]interfaces.TemplateSelection, len(saved))
	for i, sel := range saved {
		result[i] = interfaces.TemplateSelection{
			Template: interfaces.TemplateInfo{
				Name:       sel.TemplateName,
				Category:   sel.Category,
				Technology: sel.Technology,
				Version:    sel.Version,
			},
			Selected: sel.Selected,
			Options:  sel.Options,
		}
	}
	return result
}
