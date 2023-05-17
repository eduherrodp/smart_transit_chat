package webhook

import (
	"log"
	"net/http"
)

// ReceivedMessage represents the structure of the message received from WhatsApp
type ReceivedMessage struct {
	Message struct {
		Text struct {
			Body string `json:"body"`
		} `json:"text"`
	} `json:"message"`
}

// HandleWebhook handles the webhook verification
// We need get hub.mode, hub.verify_token and hub.challenge
// from the query parameters of the request
// and return hub.challenge back to Facebook
// and check if hub.verify_token is equal to the verifyToken

func HandleWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {

	// We need check if the request method came from Webhook Verify or from Whatsapp Message

	// If the request method is GET, we need to verify the webhook
	if r.Method == http.MethodGet { // Verify Webhook
		verifyWebhook(w, r, verifyToken)
	} else if r.Method == http.MethodPost { // Receive Message
		receiveMessage(w, r)
	} else {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// verifyWebhook verifies the webhook
func verifyWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {

	// Get hub.mode, hub.verify_token and hub.challenge from the query parameters of the request
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	// Check if hub.mode is equal to "subscribe"
	// and hub.verify_token is equal to verifyToken
	// and hub.challenge is not empty
	if mode == "subscribe" && token == verifyToken && challenge != "" {
		// Return hub.challenge back to Facebook
		w.Write([]byte(challenge))
	} else {
		log.Println("Error verifying webhook")
		http.Error(w, "Error verifying webhook", http.StatusBadRequest)
	}
}

// receiveMessage receives the message from WhatsApp
func receiveMessage(w http.ResponseWriter, r *http.Request) {

	// No decoding needed, just print the request body
	// because the request body is already in JSON format

	// Log the message received
	log.Printf("Message received: %s", r.Body)

	// Return a 200 OK status to WhatsApp
	w.WriteHeader(http.StatusOK)
}
