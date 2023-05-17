package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

type Contact struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaID string `json:"wa_id"`
}

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
	Type string `json:"type"`
}

type ReceivedMessage struct {
	Contacts []Contact `json:"contacts"`
	Messages []Message `json:"messages"`
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

	// Decode the JSON data
	var receivedMessage ReceivedMessage
	if err := json.NewDecoder(r.Body).Decode(&receivedMessage); err != nil {
		log.Println("Error decoding JSON data: ", err)
		http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
		return
	}

	// Print the received message
	log.Println(receivedMessage)

	// Return a response
	w.Write([]byte("Message received"))
}
