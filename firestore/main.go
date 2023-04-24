package main

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/spf13/viper"
)

type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `mapstructure:"project_id"`
}

var (
	client *firestore.Client
	once   sync.Once
)

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

// GetFirestoreClient returns a singleton instance of firestore.Client
func GetFirestoreClient() (*firestore.Client, error) {
	var err error

	once.Do(func() {
		// Get project ID from config file
		projectID := GetProjectID()

		// Get a Firestore client
		ctx := context.Background()
		client, err = CreateClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	})

	return client, err
}

func main() {
	// Get a Firestore client
	firestoreClient, err := GetFirestoreClient()
	if err != nil {
		log.Fatalf("Failed to get Firestore client: %v", err)
	}
	iter := firestoreClient.Collection("log").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
