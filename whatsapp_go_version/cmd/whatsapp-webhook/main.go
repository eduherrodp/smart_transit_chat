package main

import (
	"example.com/m/v2/config"
	"example.com/m/v2/internal/webhook"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhook.VerifyTokenHandler(w, r, cfg.VerifyToken)
	})

	http.HandleFunc("/message", webhook.MessageReceivedHandler)

	addr := ":" + cfg.Port
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}