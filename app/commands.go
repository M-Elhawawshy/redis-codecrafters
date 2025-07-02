package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "GET"
	PX   = "PX"
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

		var expiry time.Time
		if len(args) > 3 && strings.ToUpper(args[2]) == PX {
			ms, err := strconv.Atoi(args[3])
			if err != nil {
				conn.Write([]byte("-ERR invalid expire time\r\n"))
				return
			}
			expiry = time.Now().Add(time.Millisecond * time.Duration(ms))
		}

		db.Lock()
		db.memory[key] = ValueWithExpiry{
			Value:     val,
			ExpiresAt: expiry,
		}
		db.Unlock()

		conn.Write([]byte("+OK\r\n"))
	case GET:
		key := args[0]
		db.Lock()
		v, ok := db.memory[key]
		db.Unlock()
		if ok && (v.ExpiresAt.IsZero() || v.ExpiresAt.After(time.Now())) {
			length := len(v.Value)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", length, v.Value)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}
	}
}
