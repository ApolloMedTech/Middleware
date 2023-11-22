package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

// global variable to store the configuration
var (
	config Config
	once   sync.Once
)

// LoadConfig loads the configuration from a specified YAML file path.
func LoadConfig(filePath string) {
	once.Do(func() {
		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			logrus.Fatalf("Failed to read config file: %s", err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			logrus.Fatalf("Failed to unmarshal config: %s", err)
		}
	})
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	return &config
}
