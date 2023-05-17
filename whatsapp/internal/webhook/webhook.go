package webhook

import (
	"encoding/json"
	"log"
	"net/http"
)

// ReceivedWebhook represents the structure of the webhook payload received from WhatsApp
type ReceivedWebhook struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				// Add specific fields from the Webhook payload
				// specific Webhooks payload
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
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

	// Log the message received
	log.Printf("Webhook received: %+v\n", receivedWebhook)

	// Access specific fields from the webhook payload
	// specific Webhooks payload
}
