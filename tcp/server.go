package tcp

import (
	"bufio"
	"fmt"
	"github.com/Alexduanran/Distributed-System-MP2/msg"
	"net"
	"os"
)

// MultiThreadedServer creates a server running on the given port that can handles multiple connections at the same time
// runs the given function handleConnection when connection with a client is established
func MultiThreadedServer(ip string, port string, handleConnection func(net.Conn, map[net.Conn]struct{})) {
	ln, err := net.Listen("tcp", ip + ":" +port)
	checkError(err, "Listening error")

	// closing the server
	defer fmt.Println("Server closed")

	// connections stores all connections to the server
	connections := make(map[net.Conn]struct{})

	// quit channel listens for signal to quit the process
	// listen channel listens for new connection being connected
	quit := make(chan struct{})
	listen := make(chan net.Conn)

	// a thread that waits for user inputting "EXIT"
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "EXIT" {
				fmt.Println("Exiting program...")
				// send exiting message to all connected clients
				for conn, _ := range connections {
					UnicastSend(conn, msg.Message{"EXIT", msg.Chat{"","",""}})
				}
				quit <- struct{}{} // tell the main process to quit
			}
		}
	}()

	// a separate thread to listen for incoming connections
	// if a connection is received, notify the listen channel and break out of select to handle connection
	fmt.Println("Server listening...")
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			listen <- conn
		}
	}()

	for {
		// blocks until either a quit signal or a listen signal is received
		// if quit is received, close the listener and exit the server
		// else if listen is received, handle the new connection
		select {
		case <- quit:
			ln.Close()
			return
		case conn := <- listen:
			go handleConnection(conn, connections)
		}
	}
}