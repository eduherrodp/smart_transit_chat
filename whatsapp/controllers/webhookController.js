const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
const fetch = require("node-fetch");

function handleWebhook(req, res) {
    const { body } = req;

    const name = body.entry[0].changes[0].value.contacts[0].profile.name;
    const wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;
    const message = body.entry[0].changes[0].value.messages[0].text.body;
    const time = new Date().toLocaleString();

    console.log(time, "|> [Incoming message from whatsapp from]: ", wa_id)
    console.log(name, "says: ", message);

    const response = { name, wa_id, message };

    res.status(200).send("EVENT_RECEIVED");

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

function medium_webhook(response) {
    const { name, wa_id, message } = response;
    const data = {
        name,
        wa_id,
        message,
    };

    const options = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    };

    try {
        fetch("https://www.smarttransit.online/webhook/whatsapp", options)
            .then(res => res.json())
            .then(json => console.log(json))
            .catch(err => console.log(err));
    } catch (error) {
        console.log(error);
    }
}

module.exports = {
    handleWebhook,
    verifyWebhook,
};
