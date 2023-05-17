package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config represents the configuration struct that is used
// to store configuration values from the YAML config file
type Config struct {
	Port        string `mapstructure:"port"`
	VerifyToken string `mapstructure:"verify_token"`
}

// GetConfig  reads the configuration values from the YAML
// config file and returns the project ID
func GetConfig() *Config {
	// Set the file name of the configurations file
	viper.SetConfigName("config_yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config file: %v", err)
	}

	return &config
}
