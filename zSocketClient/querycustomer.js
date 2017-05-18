var net = require('net');

var HOST = 'localhost';
var PORT = 6970;

var client = new net.Socket();
client.connect(PORT, HOST, function () {

    console.log('CONNECTED TO: ' + HOST + ':' + PORT);
    // Write a message to the socket as soon as the client is connected, the server will receive it as message from the client 

    var obj = {
						type: 'querycustomer',
						cardid: '777',
    }
    client.write(JSON.stringify(obj));

});

// Add a 'data' event handler for the client socket
// data is what the server sent to this socket
client.on('data', function (data) {

    console.log('DATA: ' + data);
    // Close the client socket completely
    if (data == 'goodbye')
        client.destroy();

});

// Add a 'close' event handler for the client socket
client.on('close', function () {
    console.log('Connection closed');
});



/*
var net = require('net');

var client = new net.Socket();
client.connect(1337, '127.0.0.1', function() {
	console.log('Connected');
	client.write('Hello, server! Love, Client.');
});

client.on('data', function(data) {
	console.log('Received: ' + data);
	client.destroy(); // kill client after server's response
});

client.on('close', function() {
	console.log('Connection closed');
});
*/