package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// ResponseStrategy Define una interfaz común para las estrategias
type ResponseStrategy interface {
	ProcessResponse(responseData []byte) (string, error)
}

// WhatsappStrategy Implementación de la estrategia para el servicio de Whatsapp
type WhatsappStrategy struct {
	Data map[string]interface{}
}

func (s WhatsappStrategy) ProcessResponse([]byte) (string, error) {
	// Enviar la respuesta a Dialogflow utilizando el webhook
	dialogflowWebhookURL := "http://localhost:3002/dialogflow"

	// Construir los datos de la solicitud al webhook de Dialogflow
	requestBody := map[string]interface{}{
		"projectId":    "sanguine-tome-381917",
		"sessionId":    s.Data["wa_id"],
		"query":        s.Data["message"],
		"languageCode": "es",
	}

	// Convertir los datos de la solicitud a JSON
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Crear la solicitud HTTP
	request, err := http.NewRequest("POST", dialogflowWebhookURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "Cannot send request to Dialogflow", err
	}

	// Establecer el header de la solicitud
	request.Header.Set("Content-Type", "application/json")

	// Crear el cliente HTTP
	client := &http.Client{}

	// Enviar la solicitud al webhook de Dialogflow
	response, err := client.Do(request)
	if err != nil {
		return "Cannot send request to Dialogflow", err
	}

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Cannot read response from Dialogflow", err
	}

	// Decodificar el cuerpo de la respuesta en una estructura
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "Cannot parse response from Dialogflow", err
	}

	// Construct the response to show in medium_webhook
	// The response is in the format:
	// queryResult->responseMessages[0]->text->text[0]
	log.Printf("[dialogflow]: %s", responseData["AgentResponse"])
	return "res", nil
}

// DialogflowStrategy Implementación de la estrategia para el servicio de Dialogflow
type DialogflowStrategy struct {
	Data map[string]interface{}
}

func (s DialogflowStrategy) ProcessResponse([]byte) (string, error) {
	// Enviar la respuesta a Whatsapp o a Google Maps utilizando el webhook
	whatsappWebhookURL := "http://localhost:1024/webhook/send-message"
	//googleMapsWebhookURL := "http://localhost:3003/google-maps"

	// Construir los datos de la solicitud al webhook de Whatsapp
	requestBody := map[string]interface{}{
		"wa_id":   s.Data["SessionID"],
		"message": s.Data["AgentResponse"],
	}

	// Convertir los datos de la solicitud a JSON
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Crear la solicitud HTTP
	request, err := http.NewRequest("POST", whatsappWebhookURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "Cannot send request to Whatsapp", err
	}

	// Establecer el header de la solicitud
	request.Header.Set("Content-Type", "application/json")

	// Crear el cliente HTTP
	client := &http.Client{}

	// Enviar la solicitud al webhook de Whatsapp
	response, err := client.Do(request)
	if err != nil {
		return "Cannot send request to Whatsapp", err
	}

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Cannot read response from Whatsapp", err
	}

	// Decodificar el cuerpo de la respuesta en una estructura
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "Cannot parse response from Whatsapp", err
	}

	return "", nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la solicitud en una estructura
	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error al analizar los datos", http.StatusBadRequest)
		return
	}

	// Determinar la estrategia basada en el header de la petición X-Origin-Service
	var responseStrategy ResponseStrategy

	// The service fly in base 64
	switch r.Header.Get("X-Origin") {
	// the value of the header is on base 64
	case "whatsapp":
		responseStrategy = WhatsappStrategy{
			Data: requestData,
		}
		// Print the request data
		log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["message"].(string) + " | " + requestData["wa_id"].(string) + " | " + requestData["name"].(string))
	case "dialogflow":
		responseStrategy = DialogflowStrategy{
			Data: requestData,
		}
		// Print the request data
		if r.Header.Get("X-Intent") == "Destination Location" {
			log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["AgentResponse"].(string) + " | " + requestData["SessionID"].(string) + " | " + requestData["DestinationLocation"].(string))
		} else if r.Header.Get("X-Intent") == "Origin Location" {
			log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["AgentResponse"].(string) + " | " + requestData["SessionID"].(string) + " | " + requestData["OriginLocation"].(string))
			// Una vez que se obtiene la ubicación de origen, se envía la solicitud al servicio de Google Maps para obtener la ruta entre el origen y el destino
			// Construir los datos de la solicitud al webhook de Google Maps
			requestBody := map[string]interface{}{
				"address":     requestData["OriginLocation"],
				"destination": requestData["DestinationLocation"],
			}

			// Convertir los datos de la solicitud a JSON
			requestData, err := json.Marshal(requestBody)
			if err != nil {
				http.Error(w, "Error al analizar los datos", http.StatusBadRequest)
				return
			}

			// Crear la solicitud HTTP
			request, err := http.NewRequest("GET", "http://localhost:3003/google-maps", bytes.NewBuffer(requestData))
			if err != nil {
				http.Error(w, "Error al crear la solicitud", http.StatusBadRequest)
				return
			}

			// Establecer el header de la solicitud
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Origin", "dialogflow")

			// Crear el cliente HTTP
			client := &http.Client{}

			// Enviar la solicitud al webhook de Google Maps
			_, err = client.Do(request)
			if err != nil {
				http.Error(w, "Error al enviar la solicitud", http.StatusBadRequest)
				return
			}

		} else {
			log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["AgentResponse"].(string) + " | " + requestData["SessionID"].(string))
		}
	default:
		http.Error(w, "Servicio no soportado", http.StatusBadRequest)
		return
	}

	// Procesar la respuesta
	_, err = responseStrategy.ProcessResponse(body)
	if err != nil {
		return
	}

}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
		return
	}
	log.Println("Webhook is running on port ", 3000)
}
