package main

import (
	"flag"
	"fmt"
	"github.com/Alexduanran/Distributed-System-MP2/msg"
	"github.com/Alexduanran/Distributed-System-MP2/tcp"
	"net"
)

func main() {
	// read user input Port
	var port string
	flag.StringVar(&port, "Port", "9999", "To specify port")
	flag.Parse()

	// a map that stores and maps all usernames to their respective connections
	users := make(map[string]net.Conn)

	// func handleConnection(net.Conn, map[net.Conn]struct{})
	// handleConnection handles messages received from the client connection
	handleConnection := func(conn net.Conn, connections map[net.Conn]struct{}) {
		// receives the first message from client that contains the username
		var message msg.Message
		tcp.UnicastReceive(conn, &message)

		// username already existed
		// notify client that the username is already taken
		if _, ok := users[message.Chat.From]; ok {
			tcp.UnicastSend(conn, msg.Message{"TAKEN", msg.Chat{message.Chat.From, "", ""}})
			return
		}

		// add new username and its corresponding connection to the map users
		users[message.Chat.From] = conn
		// add the new connections to the map connections
		connections[conn] = struct{}{}

		// keep waiting for new messages from this client
		for {
			tcp.UnicastReceive(conn, &message)
			switch message.Except {
			// receives an "EXIT" message
			// deleted the client from the maps connections and users, close the connection, and exit the function
			case "EXIT":
				delete(connections, conn)
				delete(users, message.Chat.From)
				conn.Close()
				return
			// receives a normal message
			// transmits the messages to the intended receiver
			default:
				receiver, ok := users[message.Chat.To]
				// to-client not connected
				// notify client that to-client does not exist
				if !ok {
					tcp.UnicastSend(conn, msg.Message{"NIL", msg.Chat{message.Chat.To, "", ""}})
					continue
				}
				fmt.Printf("Transmitting message '%v' from %v to %v\n", message.Chat.Content, message.Chat.From, message.Chat.To)
				tcp.UnicastSend(receiver, msg.Message{"", message.Chat})
			}
		}
	}

	tcp.MultiThreadedServer("127.0.0.1", port, handleConnection)
}


