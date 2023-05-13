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
	"encoding/json"
	"example.com/m/v2/internal"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	// Rest API implementation for retrieving documents from the "log" collection
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var l internal.Log
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing request body: %v", err)
			}
		}(r.Body)

		// Create an instance of the FirestoreClientFactory
		factory := internal.FirestoreClientFactory{}

		// Create a context with a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create the document in Firestore
		_, err = factory.CreateDocument(ctx, l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		if err != nil {
			return
		}
	})

	port := ":8080"
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err.Error())
	}
}
