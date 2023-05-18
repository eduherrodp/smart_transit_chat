const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
function handleWebhook(req, res) {
    const { body } = req;
    console.log(JSON.stringify(body, null, 2));
    res.sendStatus(200);
    // Extract the name, wa_id, timestamp, message from the body of the request
    // position of the name in the body: entry->changes->contacts->profile->name
    // position of the wa_id in the body: entry->changes->contacts->profile->wa_id
    // position of the timestamp in the body: entry->changes->messages->timestamp
    // position of the message in the body: entry->changes->messages->text->body.

    // Print type of body
    console.log("Type of body: " + typeof body);

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

}

module.exports = {
    handleWebhook,
    verifyWebhook,
};
