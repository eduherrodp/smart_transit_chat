const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;
const https = require("https");
const {post} = require("axios");

function handleWebhook(req, res) {
    const { body } = req;

    let name, wa_id, message;

    if (body.entry && body.entry[0].changes && body.entry[0].changes[0].value.contacts && body.entry[0].changes[0].value.messages) {
        name = body.entry[0].changes[0].value.contacts[0].profile.name;
        wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;
        message = body.entry[0].changes[0].value.messages[0].text.body;

        let time = new Date().toLocaleString();

        const response = { name, wa_id, message };

        res.sendStatus(200);

        mediumWebhook(response).then(r => null);
        console.log(time, "|> [Incoming message]: ", wa_id + ":", name, "|> [Message]: ", message);
    } else {
        res.sendStatus(400);
    }
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
    const phone_number = wa_id.slice(0, 2) + wa_id.slice(3);

    const data = {
        messaging_product: "whatsapp",
        to: phone_number,
        type: "text",
        text: {
            preview_url: false,
            body: message,
        },
    };

    const options = {
        hostname: "graph.facebook.com",
        path: `/v16.0/101271482969769/messages`,
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer `,
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
