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
	dialogflowWebhookURL := "http://localhost:3001/dialogflow"

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

	// Do not return the response to the client because it is not necessary
	return "Request sent to Dialogflow: " + response.Status, nil
}

// DialogflowStrategy Implementación de la estrategia para el servicio de Dialogflow
type DialogflowStrategy struct {
	Data map[string]interface{}
}

func (s DialogflowStrategy) ProcessResponse([]byte) (string, error) {

	// Enviar la respuesta a Whatsapp o a Google Maps utilizando el webhook
	whatsappWebhookURL := "http://localhost:3001/whatsapp"
	//googleMapsWebhookURL := "http://localhost:3001/google-maps"

	// Si queryResult->parameters->fields->location1->structValue->fields->original->stringValue
	// es diferente de vacío entonces enviar la respuesta a google maps, de lo contrario enviar la respuesta a whatsapp

	// Construir los datos de la solicitud al webhook de Whatsapp
	//message is the response from Dialogflow:
	// queryResult->responseMessages[0]->text->text[0]
	// wa_id is the phone number of the user that sent the message and that is stored in the session id

	requestBody := map[string]interface{}{
		"message": s.Data["queryResult"].(map[string]interface{})["responseMessages"].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["text"].([]interface{})[0],
		"wa_id":   s.Data["sessionId"],
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

	log.Println("[" + "whatsapp" + "]: " + requestBody["message"].(string) + " | " + requestBody["wa_id"].(string))
	log.Printf("Response from Whatsapp: %s", response.Status)

	return "Respuesta de Dialogflow", nil
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
	var strategy ResponseStrategy

	// The service fly in base 64
	switch r.Header.Get("X-Origin") {
	// the value of the header is on base 64
	case "whatsapp":
		strategy = WhatsappStrategy{
			Data: requestData,
		}
		// Print the request data
		log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["message"].(string) + " | " + requestData["wa_id"].(string) + " | " + requestData["wa_name"].(string))
	case "dialogflow":
		strategy = DialogflowStrategy{
			Data: requestData,
		}
		// Print the request data
		log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["query"].(string) + " | " + requestData["sessionId"].(string))
	default:
		http.Error(w, "Servicio no soportado", http.StatusBadRequest)
		return
	}

	responseData := []byte("Respuesta del servicio")

	// Procesar la respuesta con la estrategia seleccionada
	result, err := strategy.ProcessResponse(responseData)
	if err != nil {
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

	// Enviar la respuesta al cliente
	w.Write([]byte(result))
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
		return
	} else {
		log.Println("Servidor iniciado en http://localhost:8080")
	}
}
