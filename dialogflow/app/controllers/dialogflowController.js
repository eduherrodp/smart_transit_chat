const dialogflowUtils = require('../utils/dialogflowUtils');

async function detectIntent(req, res) {
    const { projectId, sessionId, query, languageCode } = req.body;
    // Dialogflow require the endpoint us-central1-dialogflow.googleapis.com
    const intentResponse = await dialogflowUtils.detectIntentText(projectId, 'us-central1','5ec2d85a-2586-4594-a230-19928f05b854',sessionId, query, languageCode);

    // Return the response
    res.json(intentResponse);
}

module.exports = {
    detectIntent,
};