package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	PING   = "PING"
	ECHO   = "ECHO"
	SET    = "SET"
	GET    = "GET"
	PX     = "PX"
	CONFIG = "CONFIG"
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

		app.db.Lock()
		app.db.memory[key] = ValueWithExpiry{
			Value:     val,
			ExpiresAt: expiry,
		}
		app.db.Unlock()
		conn.Write([]byte("+OK\r\n"))
	case GET:
		key := args[0]
		app.db.RLock()
		v, ok := app.db.memory[key]
		app.db.RUnlock()
		if ok && (v.ExpiresAt.IsZero() || v.ExpiresAt.After(time.Now())) {
			length := len(v.Value)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", length, v.Value)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}
	case CONFIG:
		subcommand := args[0]
		switch strings.ToUpper(subcommand) {
		case "SET":
		case "GET":
			key := args[1]
			app.config.RLock()
			val := app.config.settings[key]
			app.config.RUnlock()
			lenKey := len(key)
			lenVal := len(val)
			resp := fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", lenKey, key, lenVal, val)
			conn.Write([]byte(resp))
		}
	}
}
