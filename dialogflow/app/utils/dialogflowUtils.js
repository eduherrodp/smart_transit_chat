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
    console.log(`Query Text: ${response.queryResult.text}`);


    for (const message of response.queryResult.responseMessages) {
        if (message.text) {
            console.log(`Agent Response: ${message.text.text}`);
        }
    }

    // We need to return the following:
    // Agent Response: response.queryResult.text
    // Session ID: sessionId

    // If the display name intent is "Destination Location" then
    // we need to return the following:
    // Destination Location: queryResult.match.parameters.fields.location1.structValue.fields.original.stringValue

    // Construct the response

    let data;
    if (response.queryResult.match.intent.displayName === 'Destination Location') {
        data = {
            'Agent Response': response.queryResult.text,
            'Session ID': sessionId,
            'Destination Location': response.queryResult.match.parameters.fields.location1.structValue.fields.original.stringValue,
        };
    } else {
        data = {
            'Agent Response': response.queryResult.text,
            'Session ID': sessionId,
        };
    }

    return data;
}

module.exports = {
    detectIntentText,
};