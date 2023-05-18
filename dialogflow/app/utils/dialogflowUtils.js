// dialogflowUtils.js
const { SessionsClient } = require('@google-cloud/dialogflow-cx');

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

    const request = {
        session: sessionPath,
        queryParams: {
            disableWebhook: true,
        },
        queryInput: {
            text: {
                text: query,
            },
            languageCode: languageCode,
        },
    };
    const [response] = await client.detectIntent(request);
    console.log(`Detect Intent Request: ${request.queryParams.disableWebhook}`);
    // Show intent match
    console.log(`Detected Intent: ${response.queryResult.intent}`);
    // Show what is received from dialogflow
    console.log(`Query Text: ${response.queryResult.text}`);


    for (const message of response.queryResult.responseMessages) {
        if (message.text) {
            console.log(`Agent Response: ${message.text.text}`);
        }
    }

    return response;
}



module.exports = {
    detectIntentText,
};