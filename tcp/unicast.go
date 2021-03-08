package tcp

import (
	"encoding/gob"
	"net"

	"github.com/Alexduanran/Distributed-System-MP2/msg"
)

// UnicastSend transmits the given message msg to the given connection conn
func UnicastSend(conn net.Conn, msg msg.Message) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(msg)
	checkError(err, "Encode error")
}

// UnicastReceive receives message from connection conn and stores it in msg
func UnicastReceive(conn net.Conn, msg *msg.Message) {
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(msg)
	checkError(err, "Decode error")
}
