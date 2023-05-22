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

func (s WhatsappStrategy) ProcessResponse(responseData []byte) (string, error) {
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
		return "Error realizando la solicitud", err
	}

	// Agregar el header X-Origin
	request.Header.Add("X-Origin", "whatsapp")

	// Enviar la solicitud al webhook de Dialogflow
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "Error realizando la solicitud", err
	}

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "Error leyendo la respuesta", err
	}

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return "Error analizando la respuesta", err
	}
	// Devolver la respuesta de Dialogflow
	return "Respuesta de Dialogflow", nil
}

// DialogflowStrategy Implementación de la estrategia para el servicio de Dialogflow
type DialogflowStrategy struct {
	Data map[string]interface{}
}

func (s DialogflowStrategy) ProcessResponse(responseData []byte) (string, error) {
	// Lógica específica para procesar la respuesta del servicio de Dialogflow
	// y devolver una representación adecuada
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
		log.Println("[" + r.Header.Get("X-Origin") + "]: " + requestData["message"].(string) + " | " + requestData["wa_id"].(string))
	case "dialogflow":
		strategy = DialogflowStrategy{
			Data: requestData,
		}
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
