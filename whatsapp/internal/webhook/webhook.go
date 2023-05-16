package webhook

import (
	"encoding/json"
	"io/ioutil"
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

// HandleMessage Now if the webhook is verified, we can start receiving messages from WhatsApp.
func HandleMessage(w http.ResponseWriter, r *http.Request) {
	// Get the body of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
		return
	}

	// Parse the body into a struct
	var receivedMessage ReceivedMessage
	err = json.Unmarshal(body, &receivedMessage)
	if err != nil {
		log.Println("Error parsing request body into struct")
		return
	}

	// Get the message from the struct
	message := receivedMessage.Message.Text.Body
	log.Printf("Message received: %s", message)

	// TODO: Process the message
	// ...

}
