package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Versions struct {
	Frontend struct {
		NextJS struct {
			Version string `yaml:"version"`
		} `yaml:"nextjs"`
	} `yaml:"frontend"`
}

func LoadVersions(configPath string) (*Versions, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var versions Versions
	if err := yaml.Unmarshal(data, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}
