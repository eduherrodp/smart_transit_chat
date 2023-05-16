// dialogflowUtils.js
const dialogflow = require('@google-cloud/dialogflow');

const sessionClient = new dialogflow.SessionsClient();

async function detectIntent(projectId, sessionId, query, contexts, languageCode) {
    const sessionPath = sessionClient.projectAgentSessionPath(projectId, sessionId);

    const request = {
        session: sessionPath,
        queryInput: {
            text: {
                text: query,
                languageCode: languageCode,
            },
        },
        queryParams: {
            payload: {
                fields: {
                    actor: {
                        stringValue: '5ec2d85a-2586-4594-a230-19928f05b854',
                    },
                },
            },
        },
    };

    if (contexts && contexts.length > 0) {
        request.queryParams.contexts = contexts;
    }

    const [response] = await sessionClient.detectIntent(request);
    return response;
}

async function executeQueries(projectId, sessionId, queries, languageCode) {
    let context;
    let intentResponse;
    try {
        intentResponse = await detectIntent(projectId, sessionId, queries.join(' '), context, languageCode);
        context = intentResponse.queryResult.outputContexts;
    } catch (error) {
        console.log(error);
    }
    return intentResponse;
}

module.exports = {
    executeQueries,
};