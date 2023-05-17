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

	// Verify that the request came from Facebook
	if r.URL.Query().Get("hub.mode") == "subscribe" && r.URL.Query().Get("hub.verify_token") == verifyToken {
		// Respond with hub.challenge
		log.Println("Webhook verified")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Query().Get("hub.challenge")))
	} else {
		log.Println("Webhook not verified")
		w.WriteHeader(http.StatusForbidden)
	}
}
