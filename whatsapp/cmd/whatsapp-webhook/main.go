package main

import (
	"example.com/m/v2/config"
	"example.com/m/v2/internal/webhook"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	// VerifyToken is the token used to verify the webhook and the verification will be process with GET method
	// Messages will be received with POST method

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhook.HandleWebhook(w, r, cfg.VerifyToken)
	})

	addr := ":" + cfg.Port
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
