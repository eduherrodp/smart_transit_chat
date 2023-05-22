const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
const https = require("https");
const {post} = require("axios");

function handleWebhook(req, res) {
    const { body } = req;
    const name = body.entry[0].changes[0].value.contacts[0].profile.name;
    const wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;
    const message = body.entry[0].changes[0].value.messages[0].text.body;
    const time = new Date().toLocaleString();

    const response = { name, wa_id, message };

    res.status(200).send("EVENT_RECEIVED");

    mediumWebhook(response).then(r => console.log(r));
    console.log(time, "|> [Incoming message]: ", wa_id + ":", name, "|> [Message]: ", message);
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

async function mediumWebhook(response) {
    const { name, wa_id, message } = response;
    const data = {
        name,
        wa_id,
        message,
    };

    const mediumWebhookURL = 'http://localhost:3000/webhook';

    try {
        const axiosResponse = await post(mediumWebhookURL, data, {
            headers: {
                'Content-Type': 'application/json',
                'X-Origin': 'whatsapp'
            },
        });

        console.log(axiosResponse.data);
    } catch (error) {
        console.error("Error sending request to Medium Webhook: ", error);
    }
}

// Send message to WhatsApp user using Facebook Graph API
async function sendMessage(req, res) {
    const { body } = req;
    const { wa_id, message } = body;
    const phone_number = "522471749048"


    const data = {
        messaging_product: "whatsapp",
        recipient: {
            id: phone_number
        },
        message: {
            text: message
        }
    };

    const options = {
        hostname: "graph.facebook.com",
        path: `/v16.0/me/messages?access_token=YOUR_ACCESS_TOKEN`,
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    };

    const httpRequest = https.request(options, (response) => {
        response.on("data", (d) => {
            process.stdout.write(d);
        });
    });

    httpRequest.on("error", (error) => {
        console.error(error);
    });

    httpRequest.write(JSON.stringify(data));
    httpRequest.end();
}


module.exports = {
    handleWebhook,
    verifyWebhook,
    sendMessage,
};
