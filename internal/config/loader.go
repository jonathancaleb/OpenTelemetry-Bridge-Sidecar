package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(filePath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filePath)

	// Also check for environment variables
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		// If config file doesn't exist, we can still use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &cfg, nil
}
