package main

import "net"

const (
	PING = "PING"
	ECHO = "ECHO"
)

func processCommand(command []string, conn net.Conn) {
	cmd := command[0]
	args := command[1:]
	switch cmd {
	case PING:
		conn.Write([]byte("+PONG\r\n"))
	case ECHO:
		conn.Write([]byte(args[0]))
	}
}
