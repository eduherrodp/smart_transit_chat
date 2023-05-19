const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
function handleWebhook(req, res) {
    const { body } = req;

    const name = body.entry[0].changes[0].value.contacts[0].profile.name;
    const wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;
    const message = body.entry[0].changes[0].value.messages[0].text.body;
    // Set the time of the system
    const time = new Date().toLocaleString();

    console.log(time,"|> [Incoming message from whatsapp from]: ", wa_id)
    console.log(name, "says: ", message);


    const response = { name, wa_id, message };

    // Notificar a whatsapp que se recibió el mensaje y cerrar el request
    res.status(200).send("EVENT_RECEIVED");

    // Sent to medium webhook
    medium_webhook(response);
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
