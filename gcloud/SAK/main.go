package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iam/v1"
	"io/ioutil"
	"path/filepath"
)

// createKey creates a service account key.
func createKey() (*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	// Read the JSON file
	jsonFilePath := filepath.Join("private/", "sanguine-tome-381917-fd857170bedd.json")
	jsonFile, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %v", err)
	}

	// Parse the JSON file
	var jsonData map[string]interface{}
	err = json.Unmarshal(jsonFile, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	// Get the client email from the JSON data
	clientEmail, ok := jsonData["client_email"].(string)
	if !ok {
		return nil, fmt.Errorf("client_email not found in JSON data")
	}

	resource := "projects/-/serviceAccounts/" + clientEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}

	return key, nil
}

func main() {
	// Call to function createKey to generate the key
	key, err := createKey()
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}

	// Convert the key to JSON format
	keyJSON, err := json.Marshal(key)
	if err != nil {
		fmt.Printf("Error converting key to JSON: %v\n", err)
		return
	}

	// Write the key to a file named SAK_keyfile.json in the current directory
	err = ioutil.WriteFile("private/SAK_keyfile.json", keyJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing key to file: %v\n", err)
		return
	}

	fmt.Println("Key created successfully and saved to SAK_keyfile.json")
}
