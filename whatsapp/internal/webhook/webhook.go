package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

// ReceivedMessage represents the structure of the message received from WhatsApp
type ReceivedMessage struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
}

// ReceivedWebhook represents the structure of the webhook payload received from WhatsApp
type ReceivedWebhook struct {
	Field string `json:"field"`
	Value struct {
		MessagingProduct string            `json:"messaging_product"`
		Metadata         Metadata          `json:"metadata"`
		Contacts         []Contact         `json:"contacts"`
		Messages         []ReceivedMessage `json:"messages"`
	} `json:"value"`
}

// Metadata represents the metadata structure within the webhook payload
type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// Contact represents the contact structure within the webhook payload
type Contact struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	WaID string `json:"wa_id"`
}

// HandleWebhook handles the webhook verification and message reception
func HandleWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {
	// If the request method is GET, verify the webhook
	if r.Method == http.MethodGet {
		verifyWebhook(w, r, verifyToken)
	} else if r.Method == http.MethodPost {
		receiveMessage(w, r)
	} else {
		log.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// verifyWebhook verifies the webhook
func verifyWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {
	// Get the query parameters
	query := r.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")

	// Check if the query parameters are present
	if mode == "" || token == "" {
		log.Println("Missing query parameters")
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	// Check if the mode is subscribed and the token is correct
	if mode == "subscribe" && token == verifyToken {
		log.Println("Webhook verified")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(challenge))
		if err != nil {
			log.Println("Error writing response")
			return
		}
		return
	}
}

// receiveMessage receives the message from WhatsApp
func receiveMessage(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body
	var receivedWebhook ReceivedWebhook
	err := json.NewDecoder(r.Body).Decode(&receivedWebhook)
	if err != nil {
		log.Println("Error decoding JSON body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Access the first message in the webhook payload
	if len(receivedWebhook.Value.Messages) > 0 {
		message := receivedWebhook.Value.Messages[0]

		// Access the fields of the message
		from := message.From
		messageID := message.ID
		timestamp := message.Timestamp
		messageType := message.Type
		messageBody := message.Text.Body

		// Log the message received
		log.Printf("From: %s", from)
		log.Printf("Message ID: %s", messageID)
		log.Printf("Timestamp: %s", timestamp)
		log.Printf("Message Type: %s", messageType)
		log.Printf("Message Body: %s", messageBody)
	}

	// Return a 200 OK status to WhatsApp
	w.WriteHeader(http.StatusOK)
}
