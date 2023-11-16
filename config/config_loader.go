package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// LoadConfig carrega a configuração do arquivo config.yaml
func LoadConfig() (Config, error) {
	var config Config

	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
