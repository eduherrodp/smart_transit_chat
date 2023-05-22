// routes/webhookRoutes.js

const express = require("express");
const router = express.Router();
const webhookController = require("../controllers/webhookController");

router.post("/", webhookController.handleWebhook);
router.get("/", webhookController.verifyWebhook);
router.post("/", webhookController.sendMessage);
module.exports = router;
