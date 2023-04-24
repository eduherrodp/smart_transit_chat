package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/spf13/viper"
)

type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `mapstructure:"project_id"`
}

// GetProjectID gets the project ID from the config file
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

// CreateClient creates a Firestore client
func CreateClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	// Get project ID from config file
	projectID := GetProjectID()

	// Get a Firestore client
	ctx := context.Background()
	client, err := CreateClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer func(client *firestore.Client) {
		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}(client)
}
