package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	PING = "PING"
	ECHO = "ECHO"
)

func processCommand(command []string, conn net.Conn) {
	cmd := strings.ToUpper(command[0])
	args := command[1:]
	switch cmd {
	case PING:
		conn.Write([]byte("+PONG\r\n"))
	case ECHO:
		length := len(args[0])
		val := args[0]
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", length, val)))
	}
}
