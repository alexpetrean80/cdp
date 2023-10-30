package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Editor string `yaml:"editor"`
	Source struct {
		Dirs           []string `yaml:"dirs"`
		ProjectMarkers []string `yaml:"project_markers"`
	} `yaml:"source"`
}

func New() (*Config, error) {
	conf := &Config{}
	file, err := os.Open(fmt.Sprintf("%s/.config/cdp/config.yaml", os.Getenv("HOME")))
	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(file).Decode(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c Config) Dirs() []string {
	homeDir := os.Getenv("HOME")
	var res []string
	for _, dir := range c.Source.Dirs {
		res = append(res, fmt.Sprintf("%s/%s", homeDir, dir))
	}

	return res
}

func (c Config) Markers() []string {
	return c.Source.ProjectMarkers
}
