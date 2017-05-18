var net = require('net');

var HOST = 'localhost';
var PORT = 6969;

// Create a server instance, and chain the listen function to it
// The function passed to net.createServer() becomes the event handler for the 'connection' event
// The sock object the callback function receives UNIQUE for each connection
net.createServer(function (sock) {

    // We have a connection - a socket object is assigned to the connection automatically
    console.log('CONNECTED: ' + sock.remoteAddress + ':' + sock.remotePort);

    // Add a 'data' event handler to this instance of socket
    sock.on('data', function (data) {

        console.log('DATA ' + sock.remoteAddress + ': ' + data);
        json = JSON.parse(data);
        console.log("name = ", json.name);
        console.log("age  = ", json.age);
        // Write the data back to the socket, the client will receive it as data from the server
        sock.write('You said "' + data + '"');
        sock.write('goodbye');
    });

    // Add a 'close' event handler to this instance of socket
    sock.on('close', function (data) {
        console.log('CLOSED: ' + sock.remoteAddress + ' ' + sock.remotePort);
    });

}).listen(PORT, HOST);

console.log('Server listening on ' + HOST + ':' + PORT);


// console.log("Hello World");

// var http = require('http');

// // Configure our HTTP server to respond with Hello World to all requests.
// var serverhttp = http.createServer(function (request, response) {
//   response.writeHead(200, {"Content-Type": "text/plain"});
//   response.end("Hello World\n");
// });

// // Listen on port 8000, IP defaults to 127.0.0.1
// serverhttp.listen(8000);

// // Put a friendly message on the terminal
// console.log("Server running at http://127.0.0.1:8000/");



/*
var net = require('net')
const serversock = net.createServer((socket) => {
  socket.end('goodbye\n');
});
serversock.on('error', (err) => {
  // handle errors here
  throw err;
});

// grab a random port.
serversock.listen(() => {
  console.log('opened server on', serversock.address());
});

*/
/*
var net = require('net');

var server = net.createServer(function(socket) {
	socket.write('Echo server\r\n');
	socket.pipe(socket);
});

server.listen(1337, '127.0.0.1');*/