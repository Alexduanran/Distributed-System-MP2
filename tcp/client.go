package tcp

import (
	"fmt"
	"net"
)

// Connect returns the connection of a client if connection is successfully established with server running on port
func Connect(ip string, port string) (net.Conn) {
	conn, err := net.Dial("tcp", ip + ":"+port)
	checkError(err, "Connection error")
	fmt.Println("Connection established...")

	return conn
}