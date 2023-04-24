package main

import (
	"context"
	"github.com/spf13/viper"
	"log"

	"cloud.google.com/go/firestore"
)

type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `mapstructure:"project_id"`
}

// Code based on the documentation: https://github.com/spf13/viper
func getProjectID() string {
	// Load configuration from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Printf("ProjectID: %s", config.ProjectID)
	return config.ProjectID
}

// Create a client to connect to Firestore, we will be parsed the projectID from config file
func createClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Principal function
func main() {
	// Get projectID from config file
	projectID := getProjectID()

	// Get a Firestore client
	ctx := context.Background()
	client, err := createClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer func(client *firestore.Client) {
		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}(client)
}
