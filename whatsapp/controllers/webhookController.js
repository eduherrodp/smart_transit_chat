const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
function handleWebhook(req, res) {
    const { body } = req;
    // console.log(JSON.stringify(body, null, 2));

    console.log(body);


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
