package webhook

import (
	"io/ioutil"
	"log"
	"net/http"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request, verifyToken string) {
	if r.URL.Query().Get("hub.verify_token") == verifyToken {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		if err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
	} else {
		http.Error(w, "Invalid verify token", http.StatusUnauthorized)
	}
}

func HandleMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Incoming message: %s", string(body))
	w.WriteHeader(http.StatusOK)
}
