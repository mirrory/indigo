const express = require('express')
const app = express()
const port = 3000

// to do set session guids

app.use(express.static('ui'))

app.get('/', (req, res) => {res.send('Aloha!')})

app.get('/indigo', (req, res) => {
	res.sendFile('indigo.html', { root: 'ui/html' });
});

app.listen(port, () => {console.log(`Indigo running on port ${port}`)})
