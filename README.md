# Distributed-System-MP2

MP2 for Distributed System Spring 2021

A simple chat room application that supports only private message

## How To Run
### Server
In one terminal, start the server
```bash
go run server.go -Port [port#]
```
where [port#] should be replaced with a port number of your choice.

In the command line, input
```bash
EXIT
```
to exit the server process as well as all the open client processes.
### Clients
In separate terminals, start the client processes
```bash
go run client.go -Address [address] -Port [port#] -Username [name]
```
where [address] and [port#] should be the host address and port number of the server process, and [username] should be a unique username of your choice. 

Once the client has successfully connected to the server, in the command line, input 
```bash
NEW
```
and follow the prompts to send a private message to other users.

Otherwise, enter
```bash
EXIT
```
to exit the client process. 

## Package Design
### MSG
```go
type Chat struct {
  To string       // username of message receipient
  From string     // username of message sender
  Content string  // content of message
}
```
```go
type Message struct {
  Except string  // signal for special circumstances
  Chat Chat      // direct message
}
```
When `Except` is 
  * `""` – a normal message from client to server or server to client
  * `"EXIT"` - an exit signal signifying a process is exiting
  * `"JOIN"` - the message a client sends to the server right after it connects to inform the server of its username
  * `"TAKEN"` - a message from server notifying a client that its username is already taken by someone else
  * `"NIL"` - a message from server notifying a client that its target recepient does not exist
  
### TCP
Seperates TCP logic from the main process.
Supports building servers, connecting clients, and sending/encoding/decoding messages using ` "encoding/gob" ` in between. 
#### tcp/server.go
```go
// quit channel listens for signal to quit the process
// listen channel listens for new connection being connected
quit := make(chan struct{})
listen := make(chan net.Conn)
```
- a main thread that blocks until a *quit* signal is received from the quit channel or a new client is connected
- a separate thread that waits for an "EXIT" command from the user; if the command is received, notify the main thread through the **quit** channel
- a separate thread that listens for incoming connections; if a new connection is received, notify the main thread through the **listen** channel

### Main
#### Server.go
```go
// a map that stores and maps all usernames to their respective connections
users := make(map[string]net.Conn)
```
```go
// connections stores all connections to the server
connections := make(map[net.Conn]struct{})
```
- A username and its connection is deleted once its client is exited and the server receives its exit message
- If errors are encountered (i.e. to-client username does not exist), notify the client; else, transmit the message to the target recepient

#### Client.go
```go
// quitScan channel listens for signal to stop the scanning process
// scan channel listens for signal when user inputs a new command
quitScan := make(chan struct{})
scan := make(chan struct{})
```
- a main thread that blocks until a *quit* signal is received from the quitScan channel or the user has input a new message
- a separate thread that waits for the user input and notify the main thread through the **scan** channel
- a separate thread that waits for messages from the server, and close the main thread through the **quitScan** channel if an "EXIT" message is received.
