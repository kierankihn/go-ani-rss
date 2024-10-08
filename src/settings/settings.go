package settings

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ItemConfig struct {
	Name     string `yaml:"name"`
	Filter   string `yaml:"filter"`
	Path     string `yaml:"path"`
	Season   int    `yaml:"season"`
	Progress int    `yaml:"progress"`
}

type AppConfig struct {
	ItemConfigs []ItemConfig `yaml:"downloads"`
}

var (
	Config AppConfig
)

func ParserSettings() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	return nil
}

func SaveSettings() error {
	data, err := yaml.Marshal(&Config)
	if err != nil {
		return err
	}

	err = os.WriteFile("config.yaml", data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
