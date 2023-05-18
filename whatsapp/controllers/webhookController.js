const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
function handleWebhook(req, res) {
    const { body } = req;
    console.log(JSON.stringify(body, null, 2));
    res.sendStatus(200);
    medium_webhook(req)
}

function verifyWebhook(req, res) {
    const { "hub.mode": mode, "hub.verify_token": token, "hub.challenge": challenge } = req.query;

    if (mode && token) {
        if (mode === "subscribe" && token === WHATSAPP_VERIFY_TOKEN) {
            console.log("WEBHOOK_VERIFIED");
            res.status(200).send(challenge);
        } else {
            res.sendStatus(403);
        }
    }
}

// Function to  construct a message to send a medium_webhook
function medium_webhook(body) {
    // Extract the name, wa_id, timestamp, and message from the body of the request
    const { name, wa_id, timestamp, message } = body;
    console.log("name: " + name);
    console.log("wa_id: " + wa_id);
    console.log("timestamp: " + timestamp);
    console.log("message: " + message);
}

module.exports = {
    handleWebhook,
    verifyWebhook,
};
