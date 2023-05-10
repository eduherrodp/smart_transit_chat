package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"example.com/m/v2/config"
	"log"
	"sync"
)

// This is a singleton instance of firestore.Client
var (
	client     *firestore.Client
	clientOnce sync.Once
)

// CreateClient creates a Firestore client using the provided
// project ID and returns the client
func CreateClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// FirestoreClientFactory is a factory for creating Firestore clients
type FirestoreClientFactory struct {
	once sync.Once
}

// GetClient returns a singleton instance of firestore.Client
// by creating the client if it doesn't exist
func (f *FirestoreClientFactory) GetClient() (*firestore.Client, error) {
	var err error
	clientOnce.Do(func() {
		// Get project ID from config file
		projectID := config.GetProjectID()

		// Get a Firestore client
		ctx := context.Background()
		client, err = CreateClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	})

	return client, err
}
