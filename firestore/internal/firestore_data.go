package internal

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"time"
)

// Log Define a struct to represent a log document
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
