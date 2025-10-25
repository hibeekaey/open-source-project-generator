package input

import (
	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/mapper"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

type ComponentSelection struct {
	Folders []string
	Apps    models.Apps
}

func ReadComponentSelection() (*ComponentSelection, error) {
	selectedComponents, err := MultiSelect("Select components to create:", mapper.ComponentOptions)
	if err != nil {
		return nil, err
	}

	selection := &ComponentSelection{
		Folders: []string{},
		Apps:    models.Apps{},
	}

	for _, component := range selectedComponents {
		folder := mapper.ComponentToFolder(component)
		if folder != "" {
			selection.Folders = append(selection.Folders, folder)
		}

		if component == "frontend" {
			apps, err := MultiSelect("Select frontend apps to create:", constants.Apps.Frontend)
			if err != nil {
				return nil, err
			}
			selection.Apps.Frontend = apps
		}
	}

	return selection, nil
}
