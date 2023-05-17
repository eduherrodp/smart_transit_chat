package httpserver

import (
	"example.com/m/v2/config"
	"example.com/m/v2/internal/webhook"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.homePage)
	http.HandleFunc("/webhook", s.webhookHandler)

	log.Printf("Starting server on :%s", s.cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), returnnil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) homePage(w http.ResponseWriter, _ *http.Request) {
	fprint, err := fmt.Fprint(w, "Hello, welcome to WhatsApp Webhook tester!")
	if err != nil {
		log.Println("HomePage: Error writing response: ", err)
		return
	}
	log.Printf("HomePage: Wrote %d bytes to response.", fprint)
}

func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		webhook.VerifyTokenHandler(w, r, s.cfg.VerifyToken)
	} else if r.Method == "POST" {
		webhook.MessageReceivedHandler(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
