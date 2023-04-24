/*
 * @copyright Copyright (C) 2023 José Eduardo Hernández Rodríguez
 * @license MIT License
 * All rights reserved.
 * This source code is licensed under the MIT License found in the
 * LICENSE file in the root directory of this source tree
 *
 * Author: José Eduardo Hernández Rodríguez
 * Email: eduher.rodp@gmail.com
 * Date: April 23, 2023
 *
 * Description: This file contains the implementation of a factory for creating
 * Firestore clients. along with a main function that retrieves documents from
 * the "log" collection using the singleton instance of firestore.Client.
 */

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

// Config represents the configuration struct that is used
// to store configuration values from the YAML config file
type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `mapstructure:"project_id"`
}

// This is a singleton instance of firestore.Client
var (
	client *firestore.Client
	once   sync.Once
)

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

	f.once.Do(func() {
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

// main function that creates a factory for Firestore clients
// and retrieves documents from the "log" collection
func main() {
	// Create a factory for Firestore clients
	factory := FirestoreClientFactory{}

	// Get a Firestore client
	firestoreClient, err := factory.GetClient()
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
