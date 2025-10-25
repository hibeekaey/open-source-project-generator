package config

import (
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"gopkg.in/yaml.v3"
)

type FrontendConfig struct {
	Enabled bool     `yaml:"enabled"`
	Apps    []string `yaml:"apps"`
}

type Components struct {
	Frontend interface{} `yaml:"frontend"`
	Backend  bool        `yaml:"backend"`
	Mobile   bool        `yaml:"mobile"`
	Deploy   bool        `yaml:"deploy"`
	Docs     bool        `yaml:"docs"`
	Scripts  bool        `yaml:"scripts"`
	Github   bool        `yaml:"github"`
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
	Apps         []string
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
		Apps:         []string{},
	}

	switch frontendValue := cfg.Components.Frontend.(type) {
	case []interface{}:
		if len(frontendValue) > 0 {
			project.Folders = append(project.Folders, constants.FolderApp)
			for _, app := range frontendValue {
				if appStr, ok := app.(string); ok {
					project.Apps = append(project.Apps, appStr)
				}
			}
		}
	case bool:
		if frontendValue {
			project.Folders = append(project.Folders, constants.FolderApp)
		}
	}
	if cfg.Components.Backend {
		project.Folders = append(project.Folders, constants.FolderCommonServer)
	}
	if cfg.Components.Mobile {
		project.Folders = append(project.Folders, constants.FolderMobile)
	}
	if cfg.Components.Deploy {
		project.Folders = append(project.Folders, constants.FolderDeploy)
	}
	if cfg.Components.Docs {
		project.Folders = append(project.Folders, constants.FolderDocs)
	}
	if cfg.Components.Scripts {
		project.Folders = append(project.Folders, constants.FolderScripts)
	}
	if cfg.Components.Github {
		project.Folders = append(project.Folders, constants.FolderGithub)
	}

	return project, nil
}
