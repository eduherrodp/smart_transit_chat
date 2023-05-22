const express = require ('express');
const dialogflowRoutes = require('./routes/dialogflowRoutes');

const app = express();

app.use(express.json());

app.use('/', dialogflowRoutes);

const port = 3002;
app.listen(port, () => console.log('Server running on port ', port));