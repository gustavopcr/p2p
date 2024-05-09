const net = require('net');

// Define the TCP server parameters
const HOST = '172.25.30.133'; // Change this to your server's IP address if needed
const PORT = 6060; // Change this to your server's port

// Create a TCP client
const client = new net.Socket();

// Connect to the server
client.connect(PORT, HOST, () => {
    console.log(`Connected to ${HOST}:${PORT}`);

    // Send data to the server
    const dataToSend = "HorasExtras.pdf";
    client.write(dataToSend);
});

// Listen for data from the server
client.on('data', (data) => {
    console.log(`Received data from server: ${data}`);
});

// Listen for the end of the stream (EOF)
client.on('end', () => {
    console.log('Server has closed the connection');
});

