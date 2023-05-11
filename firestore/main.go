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

package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"time"
	"github.com/eduherrodp/smart_transit_chat/firestore/internal"
)
type Log struct {
	ID        int64     `firestore:"id,omitempty"`
	Input     string    `firestore:"input,omitempty"`
	Output    string    `firestore:"output,omitempty"`
	Timestamp time.Time `firestore:"timestamp,omitempty"`
	UserID    string    `firestore:"user_id,omitempty"`
}

// CreateDocument creates a new document in the "log" collection with the given data
func (f *FirestoreClientFactory) CreateDocument(ctx context.Context, data Log) (*firestore.DocumentRef, error) {
	client, err := f.GetClient()
	if err != nil {
		return nil, err
	}

	docRef, _, err := client.Collection("log").Add(ctx, data)
	if err != nil {
		return nil, err
	}

	return docRef, nil
}

// GetDocument retrieves a document with the given ID from the "log" collection
func (f *FirestoreClientFactory) GetDocument(ctx context.Context, id string) (*Log, error) {
	client, err := f.GetClient()
	if err != nil {
		return nil, err
	}

	docRef := client.Collection("log").Doc(fmt.Sprintf("%v", id))
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve document: %v", err)
	}

	var logData Log
	if err := doc.DataTo(&logData); err != nil {
		return nil, fmt.Errorf("failed to convert document data: %v", err)
	}

	return &logData, nil
}

// UpdateDocument updates a document with the given ID in the "log" collection with the given data
func (f *FirestoreClientFactory) UpdateDocument(ctx context.Context, id int64, data Log) error {
	client, err := f.GetClient()
	if err != nil {
		return err
	}

	docRef := client.Collection("log").Doc(fmt.Sprintf("%d", id))
	_, err = docRef.Set(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDocument deletes a document with the given ID from the "log" collection
func (f *FirestoreClientFactory) DeleteDocument(ctx context.Context, id int64) error {
	client, err := f.GetClient()
	if err != nil {
		return err
	}

	docRef := client.Collection("log").Doc(fmt.Sprintf("%d", id))
	_, err = docRef.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
	"time"
)

func main() {
	// Create an instance of the FirestoreClientFactory
	factory := FirestoreClientFactory{}
	// Create a new log document
	newLog := Log{
		ID:        2,
		Input:     "example input from main",
		Output:    "example output from main",
		Timestamp: time.Now(),
		UserID:    "example_user from main",
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the document in Firestore
	docRef, err := factory.CreateDocument(ctx, newLog)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	log.Printf("Created document with ID: %s", docRef.ID)

	// Retrieve the document from Firestore
	retrievedLog, err := factory.GetDocument(ctx, docRef.ID)
	if err != nil {
		// Show id of document that failed to retrieve
		log.Print(newLog.ID)
		log.Fatalf("Failed to retrieve document: %v", err)
	}

	log.Printf("Retrieved document with ID %d: %+v", retrievedLog.ID, retrievedLog)

	// Update the document in Firestore
	retrievedLog.Output = "updated output"
	err = factory.UpdateDocument(ctx, retrievedLog.ID, *retrievedLog)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}

	log.Printf("Updated document with ID %d", retrievedLog.ID)
}
