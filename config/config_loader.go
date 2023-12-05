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

		// Override with environment variables if they are set
		if port := os.Getenv("SERVER_PORT"); port != "" {
			config.ServerConfig.Port = port
		}
		if user := os.Getenv("DB_USER"); user != "" {
			config.Database.User = user
		}
		if password := os.Getenv("DB_PASSWORD"); password != "" {
			config.Database.Password = password
		}
		if host := os.Getenv("DB_HOST"); host != "" {
			config.Database.Host = host
		}
		if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
			config.Database.Port = dbPort
		}
		if dbName := os.Getenv("DB_NAME"); dbName != "" {
			config.Database.Name = dbName
		}
	})
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	return &config
}
