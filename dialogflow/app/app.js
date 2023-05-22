const express = require ('express');
const dialogflowRoutes = require('./routes/dialogflowRoutes');

const app = express();

app.use(express.json());

app.use('/', dialogflowRoutes);

app.listen(3000, () => console.log('Server running on port 3000'));