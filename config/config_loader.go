package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// LoadConfig loads the configuration from a specified YAML file path.
func LoadConfig(filePath string) (Config, error) {
	var config Config

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
