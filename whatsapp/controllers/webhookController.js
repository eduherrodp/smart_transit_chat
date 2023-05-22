const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
const http = require("http");

function handleWebhook(req, res) {
    const { body } = req;

    const name = body.entry[0].changes[0].value.contacts[0].profile.name;
    const wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;
    const message = body.entry[0].changes[0].value.messages[0].text.body;
    const time = new Date().toLocaleString();

    const response = { name, wa_id, message };

    res.status(200).send("EVENT_RECEIVED");

    console.log(time, "|> [Incoming message]: ", wa_id + ":", name, "|> [Message]: ", message + medium_webhook(response))
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

function medium_webhook(response)  {
    var statusCode;
    const { name, wa_id, message } = response;
    const data = {
        name,
        wa_id,
        message,
    };

    // Send this data to medium webhook, medium is listening on port 3000, and has a handle in /webhook
    // Do not use fetch
    const options = {
        hostname: "localhost",
        port: 3000,
        path: "/webhook",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-Origin": "whatsapp"
        },
    };

    const req = http.request(options, (res) => {
        statusCode = res.statusCode;
        res.on("data", (d) => {
            process.stdout.write(d);
        });
    });

    req.on("error", (error) => {
        console.error(error);
    });

    req.write(JSON.stringify(data));
    req.end();

    return " |> [webhook]: " + statusCode;
}

module.exports = {
    handleWebhook,
    verifyWebhook,
};
