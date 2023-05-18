package main

import (
	"log"
	"net/http"
)

// Webhook to communication with Medium API
// Process to intercept all the request of different services

// 1. Receive the request from the service Whatsapp webhook that will be listening on the port 1024
// 2. When the request is received, the webhook will send the make a request to dialogflow-cx service with GET method and the query parameters
// 3. Dialogflow-cx will return the response with the response of the intent, dialogflow webhook will be listening on the port 3000)
// 4. The webhook will send the response to the whatsapp webhook with POST method and the body of the response
// 5. The whatsapp webhook will send the response to the user
// 6. When the user receive the response

func handleWhatsapp(w http.ResponseWriter, r *http.Request) {
	// Return a simple Hello World message
	// Get the query parameters
	message := r.URL.Query().Get("message")
	// If message is empty, return an error
	if message == "" {
		log.Printf("Message is empty")
		http.Error(w, "Message is empty", http.StatusBadRequest)
		return
	} else {
		log.Printf("Message: %s", message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`Hello from server, I received your message`))
		if err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
	}
}

// StartServer Function to start HTTP server
func StartServer(port string) error {
	// Whatsapp endpoint
	http.HandleFunc("/whatsapp", handleWhatsapp)
	// Start HTTP server
	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(port, nil)
}

func main() {
	// This webhook will be listening on the port 3000
	// Start the HTTP server
	err := StartServer(":3000")
	if err != nil {
		panic(err)
	}
}
