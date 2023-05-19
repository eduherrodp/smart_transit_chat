const { WHATSAPP_VERIFY_TOKEN } = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;

function handleWebhook(req, res) {
    const { body } = req;

    // Verificar si el cambio es un mensaje
    if (body.entry && body.entry[0].changes && body.entry[0].changes[0].value.messages) {
        const name = body.entry[0].changes[0].value.contacts[0].profile.name;
        const wa_id = body.entry[0].changes[0].value.contacts[0].wa_id;

        console.log(name);
        console.log(wa_id);

        // Realizar otras operaciones necesarias con los datos del mensaje

        // Responder al webhook para indicar que se ha procesado correctamente
        res.sendStatus(200);
    } else {
        // No se recibió un mensaje, responder con un código de estado 200
        res.sendStatus(200);
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

// Function to construct a message to send a medium_webhook
function medium_webhook(body) {
    // Tu lógica para construir una respuesta al webhook
}

module.exports = {
    handleWebhook,
    verifyWebhook,
};
