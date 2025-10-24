package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed versions.yaml
var versionsData []byte

type Versions struct {
	Frontend struct {
		NextJS struct {
			Version string `yaml:"version"`
		} `yaml:"nextjs"`
	} `yaml:"frontend"`
}

func LoadVersions() (*Versions, error) {
	var versions Versions
	if err := yaml.Unmarshal(versionsData, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}
