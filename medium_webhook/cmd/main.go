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
	// Interceptor of the request
	// whatsapp webhook send a JSON with the following structure:
	// data = {name,wa_id,message}

	// Get the query parameters
	queryParams := r.URL.Query()

	// Get the name of the user
	name := queryParams.Get("name")

	// Show in log what handle is working just now
	log.Printf("Handling request from %s", name)
}

// StartServer Function to start HTTP server
func StartServer(port string) error {
	// Whatsapp endpoint
	http.HandleFunc("/webhook/whatsapp", handleWhatsapp)
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
