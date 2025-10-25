package config

import (
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/mapper"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"gopkg.in/yaml.v3"
)

type Components struct {
	Frontend any  `yaml:"frontend"`
	Backend  bool `yaml:"backend"`
	Mobile   bool `yaml:"mobile"`
	Deploy   bool `yaml:"deploy"`
	Docs     bool `yaml:"docs"`
	Scripts  bool `yaml:"scripts"`
	Github   bool `yaml:"github"`
}

type ProjectConfig struct {
	ProjectName  string     `yaml:"project_name"`
	OutputFolder string     `yaml:"output_folder"`
	Components   Components `yaml:"components"`
}

type Project struct {
	ProjectName  string
	OutputFolder string
	Folders      []string
	Apps         models.Apps
}

func LoadProject(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	project := &Project{
		ProjectName:  cfg.ProjectName,
		OutputFolder: cfg.OutputFolder,
		Folders:      []string{},
		Apps:         models.Apps{},
	}

	switch frontendValue := cfg.Components.Frontend.(type) {
	case []any:
		if len(frontendValue) > 0 {
			project.Folders = append(project.Folders, mapper.ComponentToFolder("frontend"))
			for _, app := range frontendValue {
				if appStr, ok := app.(string); ok {
					project.Apps.Frontend = append(project.Apps.Frontend, appStr)
				}
			}
		}
	case bool:
		if frontendValue {
			project.Folders = append(project.Folders, mapper.ComponentToFolder("frontend"))
			project.Apps.Frontend = constants.Apps.Frontend
		}
	}
	if cfg.Components.Backend {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("backend"))
	}
	if cfg.Components.Mobile {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("mobile"))
	}
	if cfg.Components.Deploy {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("deploy"))
	}
	if cfg.Components.Docs {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("docs"))
	}
	if cfg.Components.Scripts {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("scripts"))
	}
	if cfg.Components.Github {
		project.Folders = append(project.Folders, mapper.ComponentToFolder("github"))
	}

	return project, nil
}
