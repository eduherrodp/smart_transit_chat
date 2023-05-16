package webhook

import (
	"log"
	"net/http"
)

// HandleWebhook handles the webhook verification
// We need get hub.mode, hub.verify_token and hub.challenge
// from the query parameters of the request
// and return hub.challenge back to Facebook
// and check if hub.verify_token is equal to the verifyToken

func HandleWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {
	// Get query params
	query := r.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")

	// Check if mode and token are in the query params
	if mode != "" && token != "" {
		// Check if mode and token are correct
		if mode == "subscribe" && token == verifyToken {
			log.Println("Webhook verified!")
			// Return challenge back to Facebook
			_, err := w.Write([]byte(challenge))
			if err != nil {
				log.Println("Error writing challenge back to Facebook")
				return
			}
			return
		}
	}

	log.Println("Webhook not verified!")
	w.WriteHeader(http.StatusBadRequest)
}
