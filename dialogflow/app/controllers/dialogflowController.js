const dialogflowUtils = require('../utils/dialogflowUtils');

async function detectIntent(req, res) {
    const { projectId, sessionId, query, languageCode } = req.body;
    const intentResponse = await dialogflowUtils.executeQueries(projectId, sessionId, [query], languageCode);
    res.status(200).json({fulfillmentText: intentResponse.queryResult.fulfillmentText});
}

module.exports = {
    detectIntent,
};