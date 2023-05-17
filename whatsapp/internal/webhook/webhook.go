package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

// ReceivedMessage represents the structure of the message received from WhatsApp
type ReceivedMessage struct {
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
}

// ReceivedWebhook represents the structure of the webhook payload received from WhatsApp
type ReceivedWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string          `json:"messaging_product"`
				Metadata         Metadata        `json:"metadata"`
				Message          ReceivedMessage `json:"message"` // Add this field to store the message
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// Metadata represents the metadata structure within the webhook payload
type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
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

	// Access the first entry in the webhook payload
	if len(receivedWebhook.Entry) > 0 {
		entry := receivedWebhook.Entry[0]

		// Access the first change in the entry
		if len(entry.Changes) > 0 {
			change := entry.Changes[0]

			// Access the value field in the change
			value := change.Value

			// Access the messaging_product and metadata fields
			messagingProduct := value.MessagingProduct
			metadata := value.Metadata

			// Access the specific fields related to the message
			// specific Webhooks payload

			// Log the message received
			log.Printf("Messaging Product: %s", messagingProduct)
			log.Printf("Display Phone Number: %s", metadata.DisplayPhoneNumber)
			log.Printf("Phone Number ID: %s", metadata.PhoneNumberID)
		}
	}

	// Return a 200 OK status to WhatsApp
	w.WriteHeader(http.StatusOK)
}
