package main

import (
	"fmt"
	"net"
	"strings"
)

const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "GET"
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
	case SET:
		key := args[0]
		val := args[1]
		db.Lock()
		db.memory[key] = val
		db.Unlock()
		conn.Write([]byte("+OK\r\n"))
	case GET:
		key := args[0]
		db.Lock()
		v, ok := db.memory[key]
		db.Unlock()
		if ok {
			length := len(v)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", length, v)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}
	}
}
