## gcloud prerequisites

This directory contains the files needed to deploy the chatbot to Google Cloud Platform. The following steps are needed to deploy the chatbot:

1. Create a new project in Google Cloud Platform.
2. Enable the following APIs:
    - Dialogflow API
    - Cloud Firestore API
    - Identity and Access Management (IAM) API
3. Clone the repository

### Cloud Firestore API

Cloud Firestore API is used to store the user's preferences. To enable the API, follow these steps:

1. Go to the [Cloud Firestore API](https://console.cloud.google.com/apis/library/firestore.googleapis.com) page and enable the API.
2. Access to Google Cloud Platform console and go to the [Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts).
3. Create the credentials needed to access the API and download the JSON file into ```gcloud```. 

#### Create a service account key (SAK)

A service account key is needed to access the Cloud Firestore API. To create a service account key you have two options (you can choose the one you prefer) in our case we will use the first option:

##### Option 1: Create a service account key using the [algorithm](https://cloud.google.com/iam/docs/keys-create-delete#iam-service-account-keys-delete-go) provided by Google and modified by us.

1. Create a directory called **_private_** in ```gcloud``` path.
2. Execute the following command to create the service account key:

```bash 
go run main.go
```

The **_SAK_keyfile.json_** will be created in the directory private.

##### Option 2: Create a service account key using the google cloud console.

Just execute the following command:

```bash

gcloud iam service-accounts keys create private/SAK_keyfile.json --iam-account [YOUR_SERVICE_ACCOUNT_NAME]@[YOUR_PROJECT_ID].iam.gserviceaccount.com

```
[YOUR_SERVICE_ACCOUNT_NAME]@[YOUR_PROJECT_ID].iam.gserviceaccount.com is the email address of the service account. You can find it in the [Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) page.

## Authors

- [Eduardo Hernandez](https://github.com/eduherrodp)



