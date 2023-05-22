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

    mediumWebhook(response);
    console.log(time, "|> [Incoming message]: ", wa_id + ":", name, "|> [Message]: ", message)
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

function mediumWebhook(res) {
    const { name, wa_id, message } = res;
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
        res.on("data", (d) => {
            process.stdout.write(d);
        });
    });

    req.on("error", (error) => {
        console.error(error);
    });

    req.write(JSON.stringify(data));
    req.end();
}

// Send message to whatsapp to the user with the data provided by the medium webhook
function sendMessage(req, res) {
    // The request to send the message to whatsapp has the next structure
    // curl -i -X POST \
    // https://graph.facebook.com/v16.0/105954558954427/messages \
    //     -H 'Authorization: Bearer EAAFl...' \
    // -H 'Content-Type: application/json' \
    // -d '{ "messaging_product": "whatsapp", "to": "15555555555", "type": "template", "template": { "name": "hello_world", "language": { "code": "en_US" } } }'

    const { body } = req;
    const { wa_id, message } = body;
    const data = {
        messaging_product: "whatsapp",
        to: wa_id,
        type: "text",
        text: {
            preview_url: false,
            body: message,
        },
    }

    const options = {
        hostname: "graph.facebook.com/v16.0/"+wa_id/+"messages",
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer EAAx1iTx7xK4BAPZBMrDoURgceb1Yb0MZBb4egtVeDgMXC8Y2jXTXqARfpIgAR48SQh8hLhZAecZBmZBd0WTmjxIxHoiJGWiOPqnoP39FloTayKRNK4PrIwZBnt4chG20fQSZBCJfduw8V4ZCUDKoThbi0LABpShJ94q4QOpqbpj7LddHUve4mY6gcpWZBMuyLBbXiawE0UalIQXxNZBRYrsVfNgmx6rsY5PmUZD"
        },
        data: JSON.stringify(data),
    }

    const req = http.request(options, (res) => {
        res.on("data", (d) => {
            process.stdout.write(d);
        });
    });

    req.on("error", (error) => {
        console.error(error);
    });

    req.write(JSON.stringify(data));
    req.end();
}


module.exports = {
    handleWebhook,
    verifyWebhook,
    sendMessage,
};
