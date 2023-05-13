package httpserver

import (
	"context"
	"encoding/json"
	"example.com/m/v2/internal"
	"io"
	"log"
	"net/http"
	"time"
)

// handleCreateDocument for drive HTTP request of create a document on Firestore
func handleCreateDocument(w http.ResponseWriter, r *http.Request) {
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

	// Create a context with a timeout of  seconds
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
}

// StartServer Function to start HTTP server
func StartServer(port string) error {
	// Registry of routes
	http.HandleFunc("/log", handleCreateDocument)

	// Start HTTP server
	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(port, nil)
}
