package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Webhook to communication with Medium API
// Process to intercept all the request of different services

// 1. Receive the request from the service Whatsapp webhook that will be listening on the port 1024
// 2. When the request is received, the webhook will send the make a request to dialogflow-cx service with GET method and the query parameters
// 3. Dialogflow-cx will return the response with the response of the intent, dialogflow webhook will be listening on the port 3000)
// 4. The webhook will send the response to the whatsapp webhook with POST method and the body of the response
// 5. The whatsapp webhook will send the response to the user
// 6. When the user receive the response

func handleWhatsapp(w http.ResponseWriter, r *http.Request) {
	// Interceptar la solicitud
	// El cuerpo de la solicitud contiene los datos enviados desde el webhook de Node.js

	// Parsear el cuerpo de la solicitud JSON en una estructura de datos
	var data struct {
		Name    string `json:"name"`
		WaID    string `json:"wa_id"`
		Message string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error al parsear la solicitud", http.StatusBadRequest)
		return
	}

	// Realizar operaciones con los datos recibidos, por ejemplo, enviar una respuesta al cliente

	// Construir la respuesta
	response := "Mensaje recibido"

	// Enviar la respuesta al webhook de Node.js
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	log.Printf("Mensaje recibido: %s", data.Message)
}

// StartServer Function to start HTTP server
func StartServer(port string) error {
	// Whatsapp endpoint
	http.HandleFunc("/webhook/whatsapp", handleWhatsapp)
	// Start HTTP server
	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(port, nil)
}

func main() {
	// This webhook will be listening on the port 3000
	// Start the HTTP server
	err := StartServer(":3000")
	if err != nil {
		panic(err)
	}
}
