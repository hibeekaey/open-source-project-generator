// Package preview provides shared types for project preview functionality.
package preview

import "github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"

// TemplateSelection represents a selected template with options
type TemplateSelection struct {
	Template interfaces.TemplateInfo
	Selected bool
}
