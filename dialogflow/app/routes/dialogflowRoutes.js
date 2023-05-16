const express = require('express');
const dialogflowController = require('../controllers/dialogflowController');

const router = express.Router();

router.post('/detectIntent', dialogflowController.detectIntent);

module.exports = router;