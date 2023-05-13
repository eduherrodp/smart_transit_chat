package main

import (
	"example.com/m/v2/internal"
	"log"
)

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

import (
	"context"
	"time"
)

func main() {
	// Create an instance of the FirestoreClientFactory
	factory := internal.FirestoreClientFactory{}
	// Create a new log document
	newLog := internal.Log{
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
