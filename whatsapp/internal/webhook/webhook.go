package webhook

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func VerifyTokenHandler(w http.ResponseWriter, r *http.Request, verifyToken string) {
	if r.URL.Query().Get("hub.verify_token") == verifyToken {
		fprintf, err := fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
		if err != nil {
			return
		}
		log.Printf("VerifyTokenHandler: Wrote %d bytes to response.", fprintf)
	} else {
		http.Error(w, "Invalid verify token", http.StatusUnauthorized)
	}
}

func MessageReceivedHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Incoming message: %s", string(body))
	w.WriteHeader(http.StatusOK)
}
