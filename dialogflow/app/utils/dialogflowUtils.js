const { SessionsClient } = require('@google-cloud/dialogflow-cx');
const axios = require('axios');

// Endpoint configuration variables

const client = new SessionsClient({ apiEndpoint: 'us-central1-dialogflow.googleapis.com' });

async function detectIntentText(projectId, location, agentId, sessionId, query, languageCode) {
    const sessionPath = client.projectLocationAgentSessionPath(
        projectId,
        location,
        agentId,
        sessionId
    );

    console.info(sessionPath);

    const detectIntentRequest = {
        session: sessionPath,
        queryInput: {
            text: {
                text: query,
            },
            languageCode: languageCode,
        },
    };

    const [response] = await client.detectIntent(detectIntentRequest);
    console.log(`Query Text: ${response.queryResult.text}`);

    for (const message of response.queryResult.responseMessages) {
        if (message.text) {
            console.log(`Agent Response: ${message.text.text}`);
        }
    }

    // Prepare the data to be sent to the medium webhook
    let data;
    if (response.queryResult.match.intent.displayName === 'Destination Location') {
        data = {
            'AgentResponse': response.queryResult.text,
            'SessionID': sessionId,
            'DestinationLocation': response.queryResult.match.parameters.fields.location1.structValue.fields.original.stringValue,
        };
    } else {
        data = {
            'AgentResponse': response.queryResult.text,
            'SessionID': sessionId,
        };
    }

    await mediumWebhook(data);

    // Just need to return the state code to the client
    return "[dialogflow]: Received\n";
}

// mediumWebhook function sends the response to the medium webhook
async function mediumWebhook(data) {
    // Check if data has the Destination Location field
    const mediumWebhookURL = 'http://localhost:3000/webhook';

    const payload = {
        'AgentResponse': data.AgentResponse,
        'SessionID': data.SessionID,
        'DestinationLocation': data.DestinationLocation,
    };

    try {
        const response = await axios.post(mediumWebhookURL, payload, {
            headers: {
                'Content-Type': 'application/json',
                'X-Origin': 'dialogflow',
            },
        });

        console.log(response.data);
    } catch (error) {
        console.error("Error sending request to Medium Webhook: ", error);
    }
}

module.exports = {
    detectIntentText,
};
