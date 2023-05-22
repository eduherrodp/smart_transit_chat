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
    let location1;
    let location2;
    if (response.queryResult.match.parameters != null) {
        location1 = response.queryResult.match.parameters.fields.location1.structValue.fields.original.stringValue
        console.log("location1: ", response.queryResult.match.parameters.fields.location1.structValue.fields.original.stringValue)
    }
    let agentResponse;
    for (const message of response.queryResult.responseMessages) {
        if (message.text) {
            // Save the agent response
            agentResponse = message.text.text[0];
            console.log(`Agent Response: ${message.text.text}`);
        }
    }

    // Prepare the data to be sent to the medium webhook
    let data;
    let header;
    if (location1 != null) {
        data = {
            'AgentResponse': agentResponse,
            'SessionID': sessionId,
            'DestinationLocation': location1,
        };
        header = {
            'Content-Type': 'application/json',
            'X-Origin': 'dialogflow',
            'X-Intent': 'Destination Location',
        }
    } else if (location2 != null) {
        data = {
            'AgentResponse': agentResponse,
            'SessionID': sessionId,
            'OriginLocation': location2,
        };
        header = {
            'Content-Type': 'application/json',
            'X-Origin': 'dialogflow',
            'X-Intent': 'Origin Location',
        }
    } else {
        data = {
            'AgentResponse': agentResponse,
            'SessionID': sessionId,
        };
        header = {
            'Content-Type': 'application/json',
            'X-Origin': 'dialogflow',
            'X-Intent': 'Default Welcome Intent',
        }
    }

    await mediumWebhook(data, header);

    // Just need to return the state code to the client
    return "[dialogflow]: Received\n";
}

// mediumWebhook function sends the response to the medium webhook
async function mediumWebhook(data, header) {
    // Check if data has the Destination Location field
    const mediumWebhookURL = 'http://localhost:3000/webhook';

    try {
        const response = await axios.post(mediumWebhookURL, data, {
            headers: header,
        });

        console.log(response.data);
    } catch (error) {
        console.error("Error sending request to Medium Webhook: ", error);
    }
}

module.exports = {
    detectIntentText,
};
