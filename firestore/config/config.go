package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config represents the configuration struct that is used
// to store configuration values from the YAML config file
type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `firestore:"project_id"`
}

// GetProjectID reads the configuration values from the YAML
// config file and returns the project ID
func GetProjectID() string {
	// Load configuration from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config file: %v", err)
	}

	log.Printf("Project ID: %s", config.ProjectID)
	return config.ProjectID
}
