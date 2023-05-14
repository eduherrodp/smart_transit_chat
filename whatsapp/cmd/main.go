package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	methodGet  = "GET"
	methodPost = "POST"
)

const (
	// This is a prob
	verifyToken = "a1b2c3d4"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Printf("HomePage: Request: %+v\n", r)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case methodGet:
		verifyTokenHandler(w, r)
	case methodPost:
		messageReceivedHandler(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("hub.verify_token")
	if token == verifyToken {
		write, err := w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		if err != nil {
			log.Println("VerifyTokenHandler: Error responding to verification request: ", err)
			return
		}
		log.Println("VerifyTokenHandler: Wrote ", write, " bytes to response.")
	} else {
		http.Error(w, "Invalid verify token", http.StatusUnauthorized)
		log.Println("VerifyTokenHandler: Invalid verification request.")
	}
}

func messageReceivedHandler(w http.ResponseWriter, r *http.Request) {
	// Process the incoming message here
	fmt.Println("Incoming message: ", r.Body)
	w.WriteHeader(http.StatusOK)
	log.Println("MessageReceivedHandler: Responded to message.")
}
