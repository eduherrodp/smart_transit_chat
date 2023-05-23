const express = require('express');
const dialogflowController = require('../controllers/dialogflowController');

const router = express.Router();

router.post('/dialogflow', dialogflowController.detectIntent);

module.exports = router;