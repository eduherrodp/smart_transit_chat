package main

import (
	"example.com/m/v2/config"
	"example.com/m/v2/httpserver"
	"log"
)

func main() {
	cfg := config.GetConfig()

	srv := httpserver.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
