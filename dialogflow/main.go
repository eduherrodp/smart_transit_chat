/*
* @copyright Copyright (C) 2023 José Eduardo Hernández Rodríguez
* @license MIT License
* All rights reserved.
* This source code is licensed under the MIT License found in the
* LICENSE file in the root directory of this source tree
*
* Author: José Eduardo Hernández Rodríguez
* Email: eduher.rodp@gmail.com
* Date: May 04, 2023
*
* TODO: Description:
 */

package main

import (
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"net/http"
)

// Request represents the request body of a Dialogflow webhook call
type Request struct {
	ResponseID  string                 `json:"responseId"`
	QueryResult dialogflow.QueryResult `json:"queryResult"`
}

// Response represents the response body of a Dialogflow webhook call
type Response struct {
	FulfillmentText string `json:"fulfillmentText"`
}

// Credentials section

// Config represents the configuration struct that is used
// to store configuration values from the YAML config file
type Config struct {
	// ProjectID is the ID of the project to run the sample
	ProjectID string `mapstructure:"project_id"`
}

// GetProjectID reads the configuration values from the YAML
// config file and returns the project ID
func GetProjectID() string {
	// Load configuration from file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Unmarshal configuration into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config file: %v", err)
	}

	log.Printf("Project ID: %s", config.ProjectID)
	return config.ProjectID
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Read request body
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Create Dialogflow client
	projectID := GetProjectID()
	sessionID := req.QueryResult.SessionID
	languageCode := req.QueryResult.LanguageCode
	creds, err := google.FindDefaultCredentials(ctx, dialogflow.DefaultAuthScopes()...)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	opts := []option.ClientOption{
		option.WithCredentials(creds),
		option.WithEndpoint("dialogflow.googleapis.com:443"),
		option.WithGRPCDialOption(
			internaloption.WithDefaultCallOptions(
				internaloption.WithDefaultTimeouts(30),
			),
		),
	}
	sessionClient, err := dialogflow.NewSessionsClient(ctx, opts...)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer sessionClient.Close()

	// Create session path
	sessionPath := sessionClient.SessionPath(projectID, sessionID)

	// Send request to Dialogflow
	text := req.QueryResult.QueryText
	request := dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}
	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		log.Fatalf("Failed to detect intent: %v", err)
	}

	// Create response
	// Extract fulfillment text from response
	var fulfillmentText string
	if response.GetQueryResult().GetFulfillmentText() != "" {
		fulfillmentText = response.GetQueryResult().GetFulfillmentText()
	} else {
		fulfillmentText = "Sorry, I couldn't understand what you said."
	}

	// Create webhook response
	webhookResponse := Response{
		FulfillmentText: fulfillmentText,
	}

	// Encode response as JSON and write to HTTP response
	if err := json.NewEncoder(w).Encode(webhookResponse); err != nil {
		log.Printf("Error encoding webhook response: %v", err)
		return
	}
}

func main() {
	// Start HTTP server
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
