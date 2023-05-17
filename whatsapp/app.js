"use strict";

const express = require("express");
const bodyParser = require("body-parser");
const { PORT } = require("./config/tsconfig.json").environment;
const webhookRoutes = require("./routes/webhookRoutes");

const app = express().use(bodyParser.json());

app.listen(PORT || 3000, () => {
    console.log(`Webhook is listening on port ${PORT}`);
});

app.use("/webhook", webhookRoutes);
