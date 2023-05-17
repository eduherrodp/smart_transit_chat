package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

type Profile struct {
	Name string `json:"name"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Context struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type Text struct {
	Body string `json:"body"`
}

type Message struct {
	Context   Context `json:"context"`
	From      string  `json:"from"`
	ID        string  `json:"id"`
	Text      Text    `json:"text"`
	Timestamp string  `json:"timestamp"`
	Type      string  `json:"type"`
}

type ReceivedMessage struct {
	Contacts []Contact `json:"contacts"`
	Messages []Message `json:"messages"`
}

// HandleWebhook handles the webhook verification
// We need to get hub.mode, hub.verify_token, and hub.challenge
// from the query parameters of the request
// and return hub.challenge back to Facebook
// and check if hub.verify_token is equal to the verifyToken

func HandleWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {

	// We need to check if the request method came from Webhook Verify or from WhatsApp Message

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

	// Get hub.mode, hub.verify_token, and hub.challenge from the query parameters of the request
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
		log.Println("Error decoding JSON data:", err)
		http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
		return
	}

	// Print the received message
	log.Println(receivedMessage)

	// Return a response
	w.Write([]byte("Message received"))
}
