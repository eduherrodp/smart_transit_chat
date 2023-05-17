
"use strict";

// Access token for the bot
// Get access token from conf file ../config/tsconfig.json

// Imports dependencies and set up http server
const request = require("request"),
    express = require("express"),
    body_parser = require("body-parser"),
    axios = require("axios"),
    app = express().use(body_parser.json()); // creates express http server

// Sets server port and logs message on success
const PORT = require("../config/tsconfig.json").environment.PORT;
app.listen(PORT || 3000, () => console.log(`Webhook is listening on port ${PORT}`));

// Accepts POST requests at /webhook endpoint
app.post("/webhook", (req, res) => {
    // Parse the request body from the POST
    let body = req.body;

    // Check the Incoming webhook message
    console.log(JSON.stringify(req.body, null, 2));
});

// Accepts GET requests at the /webhook endpoint
app.get("/webhook", (req, res) => {
    const verify_token = require("../config/tsconfig.json").whatsapp.WHATSAPP_VERIFY_TOKEN;

    // Parse params from the webhook verification request
    let mode = req.query["hub.mode"];
    let token = req.query["hub.verify_token"];
    let challenge = req.query["hub.challenge"];

    // Check if a token and mode were sent
    if (mode && token) {
        // Check the mode and token sent are correct
        if (mode === "subscribe" && token === verify_token) {
            // Respond with 200 OK and challenge token from the request
            console.log("WEBHOOK_VERIFIED");
            res.status(200).send(challenge);
        } else {
            // Responds with '403 Forbidden' if verify tokens do not match
            res.sendStatus(403);
        }
    }
});