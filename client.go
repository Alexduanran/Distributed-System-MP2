package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Alexduanran/Distributed-System-MP2/msg"
	"github.com/Alexduanran/Distributed-System-MP2/tcp"
	"net"
	"os"
)

func main() {
	// parse user inputs
	var address, port, username string

	flag.StringVar(&address, "Address", "127.0.0.1", "To specify host address")
	flag.StringVar(&port, "Port", "9999", "To specify port")
	flag.StringVar(&username, "Username", "", "To specify username")
	flag.Parse()

	// build connection with server and send username to server
	conn := tcp.Connect(address, port)
	defer fmt.Println("Client closed")
	tcp.UnicastSend(conn, msg.Message{"JOIN", msg.Chat{"", username, ""}})

	// quitScan channel listens for signal to stop the scanning process
	// quitReceive channel listens for signal to stop handleMessages
	// scan channel listens for signal when user inputs a new command
	quitScan := make(chan struct{})
	quitReceive := make(chan struct{})
	scan := make(chan struct{})

	// listener that waits for messages from the server
	go handleMessages(conn, quitScan, quitReceive)

	// listen for user command
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// a separate thread that blocks scanning for user input
		// if a new input is scanned, notify through the scan channel
		go func() {
			scanner.Scan()
			scan <- struct{}{}
		}()

		// blocks until either a quitScan signal or a scan signal is received
		// if quitScan is received, breaks out of the function
		// else if scan is received, continue to handle the new user input
		select {
		case <- quitScan:
			return
		case <- scan:
			break
		}
		input := scanner.Text()

		// if user inputs "EXIT",
		// notify server about connection closing and close the client process
		if input == "EXIT" {
			tcp.UnicastSend(conn, msg.Message{"EXIT", msg.Chat{"", username, ""}})
			quitReceive <- struct{}{}
			return

		// if user inputs "NEW"
		// compose a new message and send it to server
		} else if input == "NEW" {
			newChat := msg.Chat{"", username, ""}

			fmt.Print("Send to: ")
			scanner.Scan()
			newChat.To = scanner.Text()
			fmt.Print("--> ")
			scanner.Scan()
			newChat.Content = scanner.Text()

			tcp.UnicastSend(conn, msg.Message{"", newChat})
			fmt.Printf(">>> Sent message '%v' to %v\n", newChat.Content, newChat.To)
		}
	}
}

// handleMessages handle the messages received from the server through the connection conn
// takes in quitScan and quitReceive to close the two processes
func handleMessages(conn net.Conn, quitScan chan struct{}, quitReceive chan struct{}) {
	// receive channel listens for signal when a new message is received from the server
	receive := make(chan struct{})

	for {
		// a separate thread that waits for messages from the server
		// if a new message is received, notify through the receive channel
		var message msg.Message
		go func() {
			tcp.UnicastReceive(conn, &message)
			receive <- struct{}{}
		}()

		// blocks until either a quitReceive signal or a receive signal is received
		// if quitReceive is received, breaks out of the function
		// else if receive is received, continue to handle the message received
		select {
		case <-quitReceive:
			return
		case <-receive:
			break
		}

		// a separate thread that handles the message received from the server
		go func(msg msg.Message) {
			switch message.Except {
			// termination message received from the server
			case "EXIT":
				fmt.Println("Termination message received from server. Exiting")
				quitReceive <- struct{}{}
				quitScan <- struct{}{}
			// the username the user chooses already exists in the server's database
			case "TAKEN":
				fmt.Printf("The username '%v' is already taken. Please enter another username.\n", message.Chat.To)
				quitReceive <- struct{}{}
				quitScan <- struct{}{}
			// the username the user sends the message to does not exist in the server's database
			case "NIL":
				fmt.Printf("The username '%v' does not exist. Please enter the correct username.\n", message.Chat.To)
			// receives a message from another client
			default:
				fmt.Printf("<<< Received a message from %v: '%v'\n", message.Chat.From, message.Chat.Content)
			}
		}(message)
	}
}